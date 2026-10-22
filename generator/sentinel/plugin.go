// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel

import (
	"cmp"
	"slices"

	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/lang"
)

// ID is the name of the plugin.
const ID plugin.ID = "sentinel"

// Version is the version of the plugin. The fingerprint of a composition
// includes it.
const Version = "0.1.0"

// SentinelDirective puts a package under sentinel. PrefixParam is its param
// that states the opening of each message, and the value PrefixOff of the
// param leaves out the checks of the opening. NeighboursParam is its param
// that lists the packages whose error values the error values of the package
// do not match.
const (
	SentinelDirective directive.Name     = "sentinel"
	PrefixParam       directive.ParamKey = "prefix"
	NeighboursParam   directive.ParamKey = "no-overlap-with"
	PrefixOff                            = "off"
)

// PrefixSeparator joins the name of a package to the rest of each of its
// messages, as in "store: not found".
const PrefixSeparator = ": "

// Word is the word of the family of the checks, and TestTag is its tag.
const (
	Word    = "sentinel"
	TestTag = sdk.TagTest
)

// testTemplate is the file template of the checks in the tree of each
// target.
const testTemplate = "sentinel.test.tmpl"

// Options are the options of sentinel. Prefix is the opening of each message
// of a package whose directive states no prefix, in place of the name of the
// package. The separator of [PrefixSeparator] follows it.
type Options struct {
	Prefix string `opt:"prefix" doc:"the opening of each message in place of the name of its package"`
}

// Sentinel is one error value of a package: its name, and the reference to
// it that the template writes.
type Sentinel struct {
	Ref  *emit.TypeRef
	Name string
}

// Member is one field of an error type that a check writes. Value is the
// sample that the check writes, and Text reports that the message contains
// the text of the value unchanged, as it does for a field of text.
type Member struct {
	Value emit.Value
	Name  string
	Text  bool
}

// Kind is one error type of a package: its name and the reference to it.
// Pointer reports that a pointer to the type is the error. Cause is the
// field that keeps the cause that a check gives a value, and it is empty
// where the type has none. Compares reports that the type has the hook that
// compares it with a target, and Unwraps that it has the hook that returns
// its cause. Members are the fields of the type that a check writes.
//
// Method, CauseMethod and PeersMethod are indexes of the methods of the
// record, and each is -1 where the record has no such method:
//
//   - The body of Method declares a value of the type named got. In a
//     language without constructors, got is a composite with a sample in each
//     member, behind an address where Pointer is true, and Method is -1 for a
//     type without members. In a language with constructors, got is a
//     construction with a sample of each parameter.
//   - The body of CauseMethod declares got as the same construction with the
//     value named cause in the parameter of the cause.
//   - The body of PeersMethod declares others as the list of a construction
//     of each other error type of the package.
//
// A language without constructors has neither CauseMethod nor PeersMethod.
type Kind struct {
	Ref     *emit.TypeRef
	Name    string
	Cause   string
	Members []Member
	// value is how a check makes a value of the type.
	value       value
	Method      int
	CauseMethod int
	PeersMethod int
	Pointer     bool
	Compares    bool
	Unwraps     bool
}

// value is how a check makes a value of an error type. A language with
// constructors constructs the value from args, and from caused where the
// value takes the cause. caused contains the value named cause in the
// parameter of the cause, and it is nil where no parameter takes the cause.
type value struct {
	args       []emit.Value
	caused     []emit.Value
	constructs bool
}

// The names of the locals of a check: got is the value of an error type,
// cause is the cause that the check gives it, and others is the list of the
// other error types. causeSuffix and peersSuffix end the names of the
// methods of the record that make a value with the cause and the list of the
// others.
const (
	gotLocal    = "got"
	causeLocal  = "cause"
	othersLocal = "others"
	causeSuffix = "Cause"
	peersSuffix = "Peers"
)

// Neighbour is a package whose error values the error values of the package
// do not match: its path, its name and its error values.
type Neighbour struct {
	Path      string
	Name      string
	Sentinels []Sentinel
}

// Checks is the data of the record of the checks of one package, and the
// file template of the checks reads it. Subject is the name of the package
// and Prefix the opening of each of its messages, which is empty where the
// directive leaves out the checks of the opening. Sentinels are the error
// values of the package, Kinds its error types and Neighbours the packages
// of [NeighboursParam].
type Checks struct {
	Subject    string
	Prefix     string
	Sentinels  []Sentinel
	Kinds      []Kind
	Neighbours []Neighbour
}

// part is the part of sentinel for one language. It declares the target with
// its template tree and its helpers, and it states what the facts of eidos
// leave to the language.
type part interface {
	lang.Part
	// prefix returns the opening of each message of a package under the
	// directive without a prefix: base, which is the name of the package or
	// the option prefix, and the separator. A language without the
	// convention returns the empty prefix.
	prefix(base string) string
	// inherent reports whether a field of the name belongs to every error of
	// the language, such as the field message of TypeScript. A check writes
	// no such field as a member.
	inherent(name string) bool
}

// parts are the parts of sentinel, one for each language that it serves.
var parts = []part{golangPart{}, typescriptPart{}}

// generator is one instance of sentinel. The composition writes its options
// before the first phase call, and the handlers read them afterwards.
//
// # Concurrency
//
// The workers of a phase call read the options concurrently, and no handler
// writes them.
type generator struct {
	options *Options
}

// schema is the directive of sentinel.
var schema = directive.Schema{
	Name: SentinelDirective,
	Params: []directive.ParamSpec{
		{
			Key: PrefixParam, Type: directive.TypeString,
			Doc: "the opening of each message without its separator, and off leaves out its checks",
		},
		{
			Key: NeighboursParam, Type: directive.TypeList, ListOf: directive.TypeString,
			Doc: "the paths of the packages whose error values the error values of the package do not match",
		},
	},
	Doc: "generates the checks of the error values and of the error types of the package",
}

// New returns sentinel. sentinel records the [Checks] of each package under
// [SentinelDirective] into the family [Word] under the tag [TestTag]. The
// template of the family writes the checks in the conventions of the
// language of the plan. The checks test the package through its exported
// API, so a Go file of the family is in the external test package of its
// directory.
//
// The error values of a package are its variables whose names the rules of
// its language accept for an error value, such as ErrNotFound in Go. The
// error types are its public structs and classes without type parameters
// that the rules of its language call errors. The checks pin that each error
// value is set, has a message with the opening of [PrefixParam], shares no
// message with another, has no message that opens with another, and matches
// no other error value of the package or of a package of [NeighboursParam].
// They pin that each error type reports a message with the opening, matches
// no other error type, agrees with the standard comparison where it compares
// itself and has no cause, exposes the cause that it was given, and puts
// each field of text into its message.
//
// sentinel serves Go and TypeScript, the languages of its parts. It reports
// [NoErrors] at a package without an error value and without an error type,
// [OwnNeighbour] at a package that lists itself under [NeighboursParam], and
// [Unloaded] at a package that lists a package that the run did not load. It
// reports [Unconstructed] at an error type of a language with constructors
// that a check cannot construct. sentinel records nothing for a package that
// it reports under NoErrors. Each call returns a new plugin.
func New() plugin.Generator {
	g := &generator{options: &Options{}}
	b := sdk.NewPlugin(ID).
		Version(Version).
		Options(g.options).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, Word: Word, Template: testTemplate, External: true})
	lang.Register(b, parts)
	return b.Handle(sdk.Directive(schema, sdk.OnPackage(lang.Served(parts, g.generate)))).Generator()
}

// generate records the checks of a package under the directive into the
// source file of the first error value, or else of the first error type. The
// record takes the origin of that declaration, so the file of the checks is
// beside it.
func (g *generator) generate(m *sdk.PackageMatch, p part, e *sdk.Emitter) error {
	pkg := m.Package
	b := m.Rules()
	sentinels, origin := sentinelsOf(&b, pkg)
	kinds, first := kindsOf(m, &b, p)
	if len(sentinels) == 0 && len(kinds) == 0 {
		m.Errorf(NoErrors, "the package %s has no error value and no error type, so it has no checks", pkg.Name)
		return nil
	}
	if origin.IsZero() {
		origin = first
	}
	checks := Checks{
		Subject: pkg.Name, Prefix: g.prefixOf(m, p), Sentinels: sentinels, Kinds: kinds,
		Neighbours: neighboursOf(m),
	}
	record := &emit.Struct{Origin: origin, Name: pkg.Name, Visibility: symbol.VisibilityPublic}
	for i := range checks.Kinds {
		methodsOf(record, checks.Kinds, i)
	}
	e.FileOf(origin, TestTag).Record(checks, record)
	return nil
}

// methodsOf appends the methods of the error type kinds[i] to record, and
// sets their indexes on it, as [Kind] states them. A language without
// constructors gets the method that declares got alone.
func methodsOf(record *emit.Struct, kinds []Kind, i int) {
	k := &kinds[i]
	k.Method, k.CauseMethod, k.PeersMethod = -1, -1, -1
	if !k.value.constructs {
		if len(k.Members) == 0 {
			return
		}
		fields := make([]emit.ValueField, 0, len(k.Members))
		for j := range k.Members {
			fields = append(fields, emit.NamedField(k.Members[j].Name, k.Members[j].Value))
		}
		made := emit.Composite(k.Ref, fields...)
		if k.Pointer {
			made = emit.Address(made)
		}
		k.Method = declared(record, k.Name, gotLocal, &made)
		return
	}
	made := emit.Construction(k.Ref, k.value.args...)
	k.Method = declared(record, k.Name, gotLocal, &made)
	if k.value.caused != nil {
		caused := emit.Construction(k.Ref, k.value.caused...)
		k.CauseMethod = declared(record, k.Name+causeSuffix, gotLocal, &caused)
	}
	if len(kinds) < 2 {
		return
	}
	peers := make([]emit.ValueField, 0, len(kinds)-1)
	for j := range kinds {
		if j != i {
			peers = append(peers, emit.Element(emit.Construction(kinds[j].Ref, kinds[j].value.args...)))
		}
	}
	// The list takes the type of its first element. The list literal of
	// TypeScript does not write the type.
	list := &emit.TypeRef{Form: symbol.FormList, Elems: []*emit.TypeRef{peers[0].Value.Type}}
	others := emit.Composite(list, peers...)
	k.PeersMethod = declared(record, k.Name+peersSuffix, othersLocal, &others)
}

// declared appends to record a method named name whose body declares the
// local with the value v, and returns the index of the method.
func declared(record *emit.Struct, name, local string, v *emit.Value) int {
	method := &emit.Method{Origin: record.Origin, Name: name, Visibility: symbol.VisibilityPackage}
	method.Body.Stmts = []emit.Stmt{
		{Kind: emit.StmtAssign, Declare: true, Names: []string{local}, Value: emit.ValueExpr(*v)},
	}
	record.Methods.Append(method)
	return record.Methods.Len() - 1
}

// prefixOf returns the opening of each message of the package of m. The
// param prefix states it, and the value off leaves it empty. Without the
// param, the part of the language derives it from the option prefix, or
// else from the name of the package.
func (g *generator) prefixOf(m *sdk.PackageMatch, p part) string {
	if param, set := m.Directive().Params[PrefixParam]; set {
		if param.Str == PrefixOff {
			return ""
		}
		return param.Str + PrefixSeparator
	}
	return p.prefix(cmp.Or(g.options.Prefix, m.Package.Name))
}

// sentinelsOf returns the error values of pkg in the order of its files and
// declarations, with the identity of the first. A variable is an error
// value where the rules of its language accept its name for one.
func sentinelsOf(b *rules.Bound, pkg *node.Package) ([]Sentinel, symbol.Identity) {
	var out []Sentinel
	var first symbol.Identity
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			v, isVariable := d.(*node.Variable)
			if !isVariable {
				continue
			}
			if _, named := b.SentinelBase(v.Name); !named {
				continue
			}
			if first.IsZero() {
				first = v.ID
			}
			out = append(out, Sentinel{Name: v.Name, Ref: &emit.TypeRef{Spelling: v.Name, Target: v.ID}})
		}
	}
	return out, first
}

// kindsOf returns the public error types of the package of m sorted by name,
// with the identity of the first in the order of the files. The members of a
// type are its own public fields that a check can write, each with a sample.
//
// In a language without constructors, a check writes into the zero value of
// a type. The cause of a type is the field of its cause where the type has
// the hook that returns it, and no cause where an embed by pointer promotes
// that field, because a check assigns the cause through the embed.
//
// In a language with constructors, a check constructs a value of a type, as
// constructionOf states. The cause of a type is its field of the cause where
// its constructor has a parameter of the name of that field. kindsOf reports
// [Unconstructed] at a type without a constructor that a check can call, and
// leaves the type out.
func kindsOf(m *sdk.PackageMatch, b *rules.Bound, p part) ([]Kind, symbol.Identity) {
	ctors, constructs := b.Source().(rules.ConstructRules)
	var out []Kind
	var first symbol.Identity
	for _, f := range m.Package.Files {
		for _, d := range f.Decls {
			s, isStruct := d.(*node.Struct)
			if !isStruct || s.Visibility != symbol.VisibilityPublic || len(s.TypeParams) > 0 {
				continue
			}
			contract, isError := b.ErrorContract(s)
			if !isError {
				continue
			}
			k := Kind{
				Ref: &emit.TypeRef{Spelling: s.Name, Target: s.ID}, Name: s.Name, Members: membersOf(b, s, p),
				Pointer: contract.Pointer, Compares: contract.Is != nil, Unwraps: contract.Unwrap != nil,
				value: value{constructs: constructs},
			}
			if constructs {
				var sampled bool
				k.value.args, k.value.caused, sampled = constructionOf(b, ctors, s, contract.Cause)
				if !sampled {
					m.WarnfAt(Unconstructed, s.Position(), "the error type %s has no public constructor whose "+
						"parameters have samples, so it has no checks", s.Name)
					continue
				}
				if k.value.caused != nil {
					k.Cause = contract.Cause.Name
				}
			} else {
				k.Cause = causeOf(b, s, contract.Unwrap)
			}
			if first.IsZero() {
				first = s.ID
			}
			out = append(out, k)
		}
	}
	slices.SortFunc(out, func(a, b Kind) int { return cmp.Compare(a.Name, b.Name) })
	return out, first
}

// constructionOf returns the arguments of a construction of the error type
// s through the first public constructor of c whose parameters have
// samples, as argumentsOf states them. caused is the samples of the
// parameters before the parameter of the name of the field cause, followed
// by the value named cause. caused is nil where cause is nil, where no
// parameter has its name, and where a required parameter follows that
// parameter. sampled reports false where no constructor of s has samples.
func constructionOf(
	b *rules.Bound, c rules.ConstructRules, s *node.Struct, cause *node.Field,
) (args, caused []emit.Value, sampled bool) {
	for _, ctor := range c.Constructors(s, b.View()) {
		given, at, ok := argumentsOf(b, ctor.Params, cause)
		if !ok {
			continue
		}
		if at < 0 || slices.ContainsFunc(ctor.Params[at+1:], required) {
			return given, nil, true
		}
		return given, append(slices.Clone(given[:at]), emit.Raw(s.ID.Lang, causeLocal)), true
	}
	return nil, nil, false
}

// argumentsOf returns a sample of each parameter of params up to the first
// rest parameter, and up to the first optional parameter without a sample.
// It returns the position of the parameter of the name of the field cause,
// and -1 where no parameter before the end of the arguments has it. It
// reports false where a required parameter has no sample, and where a
// required parameter follows an optional parameter without a sample.
func argumentsOf(b *rules.Bound, params []rules.ParamView, cause *node.Field) ([]emit.Value, int, bool) {
	args := make([]emit.Value, 0, len(params))
	at := -1
	for i := range params {
		param := &params[i]
		if param.Variadic {
			break
		}
		if cause != nil && param.Name == cause.Name {
			at = i
		}
		// A parameter has no declaration of its own in the view, so its
		// sample is the sample of its type.
		sample, _ := b.SamplesOf(symbol.Identity{}, param.Ref, param.Name)
		if sample.OK() {
			args = append(args, sample.Value)
			continue
		}
		if !param.Optional || slices.ContainsFunc(params[i+1:], required) {
			return nil, 0, false
		}
		break
	}
	return args, at, true
}

// required reports whether a call passes an argument for the parameter p.
func required(p rules.ParamView) bool { return !p.Optional }

// causeOf returns the name of the field of the cause of the error type s,
// as causeField finds it. It returns the empty name where s has no cause, and
// where an embed by pointer promotes the field. The zero value of an embed
// by pointer is nil, so a check that assigns the cause through it would
// panic.
func causeOf(b *rules.Bound, s *node.Struct, unwrap *node.Method) string {
	set, _ := b.MembersOf(s)
	cause := causeField(b, &set, unwrap)
	for i := range set.Members {
		if set.Members[i].Symbol == cause {
			if set.Members[i].ViaOptional {
				return ""
			}
			return cause.Name
		}
	}
	return ""
}

// causeField returns the field that the hook unwrap of an error type returns
// as its cause, from the members of the type. The field is the cause of the
// contract of the type that declares the hook: the type itself, or the type
// of an embed that promotes the hook. causeField returns nil for a type
// without the hook.
func causeField(b *rules.Bound, set *rules.MemberSet, unwrap *node.Method) *node.Field {
	for i := range set.Members {
		if set.Members[i].Symbol == unwrap {
			owner, _ := b.View().Lookup(set.Members[i].Owner)
			declared, _ := b.ErrorContract(owner)
			return declared.Cause
		}
	}
	return nil
}

// membersOf returns the public fields that the error type s declares itself,
// each with a sample, in the order of the declaration. A promoted field is a
// member of the type that declares it, whose checks write it. A field that
// the part p calls inherent is no member. A field of text keeps its text in
// a message, so a check finds it there.
func membersOf(b *rules.Bound, s *node.Struct, p part) []Member {
	var out []Member
	for _, f := range s.Fields {
		if f.Visibility != symbol.VisibilityPublic || p.inherent(f.Name) {
			continue
		}
		sample, _ := b.SamplesOf(f.ID, f.Type, f.Name)
		if !sample.OK() {
			continue
		}
		out = append(out, Member{
			Value: sample.Value, Name: f.Name, Text: b.TypeOf(f.Type).Form == symbol.FormText,
		})
	}
	return out
}

// neighboursOf returns the packages of [NeighboursParam] with their error
// values, in the order of the param. It reports [OwnNeighbour] for the
// package of m, which it leaves out, and [Unloaded] for a package that the
// run did not load, which it keeps without error values.
func neighboursOf(m *sdk.PackageMatch) []Neighbour {
	list := m.Directive().Params[NeighboursParam].List
	own := m.Package.ID
	var out []Neighbour
	for i := range list {
		path := list[i].Str
		id := symbol.Identity{Lang: own.Lang, Package: path, Kind: symbol.KindPackage}
		if id.Package == own.Package {
			m.Errorf(OwnNeighbour, "the package %s lists itself under %s, and each error value matches itself",
				m.Package.Name, NeighboursParam)
			continue
		}
		pkg, loaded := m.Reader().PackageOf(id)
		if !loaded {
			m.Warnf(Unloaded, "the package %s lists %s under %s, and the run did not load it, "+
				"so the check of its error values compares with nothing", m.Package.Name, path, NeighboursParam)
			out = append(out, Neighbour{Path: path, Name: path})
			continue
		}
		b := m.RulesFor(own.Lang)
		sentinels, _ := sentinelsOf(&b, pkg)
		out = append(out, Neighbour{Path: path, Name: pkg.Name, Sentinels: sentinels})
	}
	return out
}
