// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"cmp"
	"errors"
	"reflect"
	"slices"

	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/internal/lang"
)

// ID is the name of the plugin.
const ID plugin.ID = "builder"

// Version is the version of the plugin. The fingerprint of a composition
// includes it.
const Version = "0.1.0"

// BuilderDirective puts a struct, a class or an alias of one under builder.
// DefaultsParam is the param of the directive that refers to a function in
// scope. The result of that function seeds the value of each new builder.
const (
	BuilderDirective directive.Name     = "builder"
	DefaultsParam    directive.ParamKey = "defaults"
)

// KeyName is the key of the fact that contains the name of the builder of a
// struct, a class or an alias under [BuilderDirective], such as
// AddressBuilder. KeySeed is the key of the fact that contains the name of the
// constructor from the seed: the function NewAddressBuilder in Go, and the
// static method create of the builder class in TypeScript. builder stamps
// both facts on each declaration that it writes a builder of, so another
// generator calls the builder without the naming of builder.
const (
	KeyName meta.KeyName = "builder.name"
	KeySeed meta.KeyName = "builder.seed"
)

// SkipTag is the key of the tag that leaves a field out of the builder, and
// SkipValue is the one value of the tag that does so, as in builder:"-".
const (
	SkipTag   = "builder"
	SkipValue = "-"
)

// Word is the word of the family of the builders and of the family of their
// tests, and the default of the option suffix. TestTag is the tag of the
// family of the tests.
const (
	Word    = "builder"
	TestTag = sdk.TagTest
)

// CompanionWord is the default of the option companion-word: the word that
// the name of a companion function joins to the name of its declaration, as
// in AddressDefaults.
const CompanionWord = "defaults"

// The tree of each target contains the file template of the builders and
// the file template of their tests.
const (
	builderTemplate = "builder.tmpl"
	testTemplate    = "builder.test.tmpl"
)

// The bodies of a record assign valueLocal to the value of a new builder,
// seedLocal to the seed of a check, firstLocal to the value that a check of
// a replacement writes first, wantLocal to the value that a check compares
// and otherLocal to the value that a check of a clone writes through the
// clone. selector joins a local and a field.
const (
	valueLocal = "v"
	seedLocal  = "seed"
	firstLocal = "first"
	wantLocal  = "want"
	otherLocal = "other"
	selector   = "."
)

// The names of the methods of the record of the checks open with these
// words: the check of a default, of a setter, of a replacement and of a
// clone. Each name ends with the name of its field, so that no two methods
// of a record share a name.
const (
	defaultCheck = "default"
	setterCheck  = "setter"
	replaceCheck = "replace"
	cloneCheck   = "clone"
)

// Options are the options of builder. Suffix is the word that the name of
// each builder joins to the name of its declaration, and CompanionWord is the
// word that the name of a companion function joins to it. An empty word
// takes its default, [Word] or [CompanionWord]. NoTests leaves out the
// checks of every builder.
type Options struct {
	Suffix        string `opt:"suffix" doc:"the word that the name of each builder joins to the name of its declaration"`
	CompanionWord string `opt:"companion-word" doc:"the word that the name of a companion function joins to the name of its declaration"`
	NoTests       bool   `opt:"no-tests" doc:"leaves out the checks of every builder"`
}

// Form is the form of the type of a member. The setters of the member and
// its copy in a clone depend on the form.
type Form string

// The forms of the type of a member. A list, a map, bytes and an optional
// are forms of the canonical shapes of eidos. A set is a type that the part
// of a language writes as a set. Every other type is a scalar, and a builder
// sets a scalar as a whole.
const (
	FormScalar   Form = "scalar"
	FormList     Form = "list"
	FormMap      Form = "map"
	FormSet      Form = "set"
	FormBytes    Form = "bytes"
	FormOptional Form = "optional"
)

// Seed is the source of the value of a new builder.
type Seed string

// The seeds of a new builder. A companion function seeds the value, and the
// declared defaults then apply over it. Without a companion, the declared
// defaults apply over the zero value of the language. Without either, the
// builder starts from the zero value.
const (
	SeedCompanion Seed = "companion"
	SeedDefaults  Seed = "defaults"
	SeedZero      Seed = "zero"
)

// Model is the data of the record of one builder, and the file template of
// the builders reads it. Subject is the name of the declaration under the
// directive. New and From are the names of the constructor from the seed and
// of the constructor from a value. In a language whose constructors are
// functions of the package, the record declares a function of each name
// beside the struct of the builder, and the template writes both functions
// at the struct. Value is the type of the value of the builder at its own
// type parameters. Struct is the struct that the builder constructs. It is
// Value for a struct and the struct of the alias for an alias. A language
// such as TypeScript has no value of an alias. Seed is the source of the
// value of a new builder, and the start body of the record declares or
// assigns that value. Members are the fields that the builder sets.
type Model struct {
	Subject string
	New     string
	From    string
	Value   *emit.TypeRef
	Struct  *emit.TypeRef
	Seed    Seed
	Members []Member
}

// Member is one field that a builder sets. Field is the name of the field in
// the source, Form is its form and Type is its type. Type is nil for a field
// without a declared type, such as a TypeScript property without a type
// annotation, and the form of such a member is scalar. Elem is the element of
// a list, the type that an optional contains and the key of a set. Key and
// Value are the key and the value of a map. Names are the names of the
// methods that the builder declares for the member. Copied reports that a
// clone of the builder copies the member.
type Member struct {
	Type   *emit.TypeRef
	Elem   *emit.TypeRef
	Key    *emit.TypeRef
	Value  *emit.TypeRef
	Names  Names
	Field  string
	Form   Form
	Copied bool
}

// Names are the names of the methods that a builder declares for one
// member, in the conventions of the language of the builder. With replaces
// the member. Append appends elements to a list. Entry and Entries set the
// entries of a map or add keys to a set. Text sets bytes from text. A name is
// empty where the builder declares no such method for the form of the
// member.
type Names struct {
	With    string
	Append  string
	Entry   string
	Entries string
	Text    string
}

// Checks is the data of the record of the checks of one builder, and the
// file template of the tests reads it. Subject is the name of the
// declaration under the directive. New and From are the names of the two
// constructors. In a language whose constructors are functions of the
// package, the record of the checks has a field of each constructor, in that
// order, whose type refers to the function. The layout of eidos qualifies
// that reference where the checks are in another package than the builder.
// Value is the type of the value at the witnesses of its type parameters,
// and Args are the witnesses. Struct is the struct that the checks
// construct. For an alias it is the struct of the alias. Instantiable
// reports that the checks can name the types that they run at. It is false
// for a generic declaration whose type parameters have no witnesses, which
// has no checks.
//
// Seeds are the fields that the first function of the record seeds, which
// the check of the constructor from a value compares one by one. Defaults
// are the fields whose defaults have checks, in the order of the functions
// of the record after the first. Setters are the checks of the setters, and
// Clones are the checks that a clone copies a member. Probes are the checks
// of the methods of the members that the checks of the setters leave out.
// Fills are the members that a clone copies, which the check of a clone sets
// before it clones the builder.
type Checks struct {
	Subject      string
	New          string
	From         string
	Value        *emit.TypeRef
	Struct       *emit.TypeRef
	Args         []*emit.TypeRef
	Seeds        []string
	Defaults     []string
	Setters      []Check
	Clones       []Check
	Probes       []Probe
	Fills        []Probe
	Instantiable bool
}

// Kind is the method of a member that a probe checks.
type Kind string

// A probe calls the method of its kind.
//
//   - KindWith is the setter that replaces a member.
//   - KindAppend is the appender of a list.
//   - KindEntry is the setter of one entry of a map or of one key of a set.
//   - KindEntries is the setter that merges a map or adds a list of keys.
//   - KindText is the setter of bytes from text.
const (
	KindWith    Kind = "with"
	KindAppend  Kind = "append"
	KindEntry   Kind = "entry"
	KindEntries Kind = "entries"
	KindText    Kind = "text"
)

// Probe is the check of one method of a member that the checks of the
// setters leave out: the setter of a member without a sample, the appender
// of a list, the setters of the entries of a map and of the keys of a set,
// and the setter of bytes from text. The fill of a member in the check of a
// clone is a probe of the kind with.
//
// Type is the type of the member at the witnesses of the type parameters of
// the subject, and Elem, Key and Value are its parts, as a [Member] records
// them. Field is the name of the member and Form is its form. Kind is the
// kind of the method, Method is its name and With is the name of the setter
// that replaces the member. Sample is the index of the method of the record
// that binds a sample of the argument of the setter to want, and -1 where the
// member has no sample. The template of each language writes a probe with
// the sample, or with a value that the language spells without one.
type Probe struct {
	Type   *emit.TypeRef
	Elem   *emit.TypeRef
	Key    *emit.TypeRef
	Value  *emit.TypeRef
	Field  string
	Form   Form
	Kind   Kind
	Method string
	With   string
	Sample int
}

// Check is one check of a member. Type is the type of the value that the
// check writes through the setter, Field is the name of the member and Form
// is its form. Setter and Appender are the names of the setter that replaces
// the member and of the setter that appends to a list. Method is the index
// of the method of the record whose body binds want. Replace is the index of
// the method whose body binds first and want, and -1 where the member has no
// check of a replacement. That check calls the setter with first and then
// with want, and expects want.
type Check struct {
	Type     *emit.TypeRef
	Field    string
	Form     Form
	Setter   string
	Appender string
	Method   int
	Replace  int
}

// part is the part of builder for one language. It declares the target with
// its template tree and its helpers. It also takes each step that the facts
// of eidos leave to the language.
type part interface {
	lang.Part
	// builder returns the name of the builder of the declaration named name,
	// which joins the name and word.
	builder(name, word string) string
	// companion returns the name of the companion function of the
	// declaration named name, which joins the name and word.
	companion(name, word string) string
	// constructors returns the names of the constructor from the seed and of
	// the constructor from a value of the builder named builder.
	constructors(builder string) (fromSeed, fromValue string)
	// functions reports whether the constructors of a builder are functions
	// of its package, and false where they are members of the builder.
	functions() bool
	// set returns the key of the type ref, and false where the language does
	// not write ref as a set.
	set(ref *node.TypeRef) (*node.TypeRef, bool)
	// names returns the names of the methods of the member named field. The
	// form of the member decides which methods it has.
	names(field string, form Form) Names
	// selectable reports whether code of the language writes a field named
	// name after a dot, as a builder reads and assigns each member.
	selectable(name string) bool
}

// parts are the parts of builder, one for each language that it serves.
var parts = []part{golangPart{}, typescriptPart{}}

// subject is a declaration under the directive builder. origin, name and
// visibility belong to the struct or the alias that the builder is named
// after. decl is the struct whose fields the builder sets, and params are
// the type parameters of the subject.
type subject struct {
	origin     symbol.Identity
	name       string
	decl       *node.Struct
	params     []*node.TypeParam
	visibility symbol.Visibility
}

// member is one field that a builder sets: a field with its declaration, or
// an embedded field, whose field is nil. id and ref are the identity and the
// type that the source declares.
type member struct {
	field *node.Field
	ref   *node.TypeRef
	name  string
	id    symbol.Identity
}

// generator is one instance of builder: its options, and the handles of its
// keys. The composition writes the options and registers the keys before the
// first phase call, and the handlers read them afterwards.
//
// # Concurrency
//
// The workers of a phase call read the options and the handles concurrently,
// and no handler writes them.
type generator struct {
	options *Options
	name    meta.Key[string]
	seed    meta.Key[string]
}

// invocation is the state of one builder. It contains the options of the
// instance, the match of the subject, the part of the language of the
// subject, the emitter, the subject and the rules of the language. The
// emitter is nil in the annotate half.
type invocation struct {
	options *Options
	m       sdk.SubjectMatch
	p       part
	e       *sdk.Emitter
	sub     *subject
	b       rules.Bound
}

// schema is the directive of builder.
var schema = directive.Schema{
	Name: BuilderDirective,
	Params: []directive.ParamSpec{{
		Key: DefaultsParam, Type: directive.TypeReference, Resolution: directive.ResolveCallableInScope,
		Doc: "the function whose result seeds the value of each new builder",
	}},
	Doc: "generates a builder of the declaration, with a setter for each field, and the checks of the builder",
}

// New returns builder, a plugin with an annotate half and a generate half.
// The annotate half stamps [KeyName] and [KeySeed] on each struct or class
// under [BuilderDirective] and on each alias of one under it, where builder
// can write the builder of the declaration. It reports each declaration that
// has no builder. The generate half runs on each declaration with the facts.
// It records a [Model] into the family [Word], under the name that the
// annotate half stamped. The template of the family writes the builder in the
// conventions of the language of the plan. builder records the [Checks] of
// each builder into the family [Word] under the tag [TestTag]. The template of
// that family writes the checks of the builder. The checks test the builder
// through the exported API of its package, so a Go file of the family is in
// the external test package of its directory. builder also records the Checks
// of a generic declaration whose type parameters have no witnesses, and the
// template writes a note in place of its checks. The [Options] name the
// builders and the companion functions, and can leave the checks out.
//
// The members of a builder are the public fields of the instance that are
// not immutable. They include the embedded fields, which a builder sets
// whole, and the fields that a base declares. A field that an embed promotes
// is no member, and a field with the tag [SkipTag] of the value [SkipValue]
// is no member either. The form of a member is the canonical shape of its
// type, and the part of the language reports whether a type is a set. The
// part also names the constructors and the methods of each member, and leaves
// out a field whose name its language does not write after a dot.
//
// The function of [DefaultsParam] seeds the value of a new builder. Without
// the param, the companion function of the declaration seeds it, and the
// part of the language names that function. Each default that the plugin
// defaults stamps on a member then applies over the seed, also on a member
// that a base declares.
// builder requires the capability of defaults for these stamps.
//
// builder serves Go and TypeScript, the languages of its parts. It leaves a
// declaration of another language alone. Its annotate half reports
// [BadDefault] at a field whose default is no value of its type, and [BadTag]
// at a field whose tag [SkipTag] has another value than [SkipValue]. It
// reports [NotStruct] at an alias whose target is not a struct of the
// workspace, and at an alias that has type parameters or refers to a generic
// struct. It reports [Unconstructed] at a declaration whose public
// constructors each take a parameter that a caller must pass, [Empty] at a
// declaration without a member, and [Collision] at a declaration whose
// builder would have two methods of one name. builder stamps and records
// nothing for a declaration that it reports, except after [BadTag], which
// leaves the field in the builder. Each call returns a new plugin.
func New() sdk.Dual {
	g := &generator{options: &Options{Suffix: Word, CompanionWord: CompanionWord}}
	b := sdk.NewPlugin(ID).
		Version(Version).
		Requires(defaults.Capability).
		Options(g.options).
		Keys(g.register).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: builderTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, Word: Word, Template: testTemplate, External: true})
	lang.Register(b, parts)
	return b.Handle(sdk.Directive(schema,
		sdk.OnStruct(lang.Served(parts, g.concludeStruct)), sdk.OnAlias(lang.Served(parts, g.concludeAlias)),
		sdk.Where(sdk.HasKey(meta.Named[string](KeyName)),
			sdk.OnStruct(lang.Served(parts, g.generateStruct)), sdk.OnAlias(lang.Served(parts, g.generateAlias)),
		),
	)).Dual()
}

// register claims the namespace of builder and registers [KeyName] and
// [KeySeed].
//
// Error modes: the error of the claim, such as a namespace that another
// registrant claimed, and the errors of the two registrations, such as a key
// that the registrant registered with another type, joined.
func (g *generator) register(r *meta.Registry) error {
	errClaim := r.ClaimNamespace(KeyName.Namespace())
	kinds := []symbol.Kind{symbol.KindStruct, symbol.KindAlias}
	var errName, errSeed error
	g.name, errName = meta.Register[string](r, meta.KeySpec{
		Name: KeyName, Kinds: kinds, Doc: "contains the name of the builder of a struct, a class or an alias",
	})
	g.seed, errSeed = meta.Register[string](r, meta.KeySpec{
		Name: KeySeed, Kinds: kinds, Doc: "contains the name of the constructor of the builder from the seed",
	})
	return errors.Join(errClaim, errName, errSeed)
}

// concludeStruct stamps the facts of the builder of a struct under the
// directive.
func (g *generator) concludeStruct(m *sdk.StructMatch, p part, st *sdk.Stamper) error {
	s := m.Struct
	g.conclude(m, p, st, &subject{origin: s.ID, name: s.Name, decl: s, params: s.TypeParams, visibility: s.Visibility})
	return nil
}

// concludeAlias stamps the facts of the builder of an alias of a struct under
// the directive. The builder is named after the alias and sets the fields of
// the struct. concludeAlias reports [NotStruct] at an alias whose target is
// not a struct of the workspace. It also reports [NotStruct] at a generic
// alias and at an alias of a generic struct, because the builder of an alias
// sets the fields of a struct without type parameters. The builder of the
// generic struct builds each instance of the struct. An alias of Go and of
// TypeScript always has a target.
func (g *generator) concludeAlias(m *sdk.AliasMatch, p part, st *sdk.Stamper) error {
	a := m.Alias
	decl, _ := m.Reader().Lookup(a.Target.Target)
	target, isStruct := decl.(*node.Struct)
	switch {
	case !isStruct:
		m.Errorf(NotStruct, "%s is an alias of %s, which is no struct of the workspace", a.Name, a.Target.Spelling)
		return nil
	case len(a.TypeParams) > 0 || len(target.TypeParams) > 0:
		m.Errorf(NotStruct, "the alias %s or its struct %s has type parameters, and a builder sets the fields "+
			"of a struct without them", a.Name, target.Name)
		return nil
	}
	g.conclude(m, p, st, &subject{origin: a.ID, name: a.Name, decl: target, visibility: a.Visibility})
	return nil
}

// conclude stamps [KeyName] and [KeySeed] on sub where builder can write the
// builder of sub. p is the part of the language of sub. conclude reports
// [Unconstructed] at a declaration whose public constructors each take a
// parameter that a caller must pass, [Empty] at a declaration without a
// member, [BadDefault] at a field whose default is no value of its type, and
// [Collision] at a declaration whose builder would have two methods of one
// name. It stamps nothing on such a declaration. conclude also reports
// [BadTag] at each member whose tag has another value than [SkipValue], and
// such a member keeps its setter.
func (g *generator) conclude(m sdk.SubjectMatch, p part, st *sdk.Stamper, sub *subject) {
	in := &invocation{options: g.options, m: m, p: p, sub: sub, b: m.Rules()}
	if !in.constructible() {
		m.Errorf(Unconstructed, "each public constructor of %s takes a parameter that a caller must pass, "+
			"and a builder constructs %s without arguments", sub.decl.Name, sub.name)
		return
	}
	members := in.members()
	in.mistagged(members)
	if len(members) == 0 {
		m.Errorf(Empty, "%s has no field that a builder can set, so it has no builder", sub.name)
		return
	}
	if _, _, _, valid := in.start(members); !valid {
		return
	}
	models := make([]Member, 0, len(members))
	for i := range members {
		models = append(models, in.modelOf(&members[i]))
	}
	if name, first, second, collides := collision(models); collides {
		m.Errorf(Collision, "the builder of %s would declare %s for %s and for %s, so %s has no builder",
			sub.name, name, first, second, sub.name)
		return
	}
	name := p.builder(sub.name, cmp.Or(g.options.Suffix, Word))
	fromSeed, _ := p.constructors(name)
	sdk.Stamp(st, g.name, name)
	sdk.Stamp(st, g.seed, fromSeed)
}

// generateStruct records the builder of a struct with the facts.
func (g *generator) generateStruct(m *sdk.StructMatch, p part, e *sdk.Emitter) error {
	s := m.Struct
	g.write(m, p, e, &subject{origin: s.ID, name: s.Name, decl: s, params: s.TypeParams, visibility: s.Visibility})
	return nil
}

// generateAlias records the builder of an alias with the facts. The annotate
// half stamps the facts on an alias of a struct alone.
func (g *generator) generateAlias(m *sdk.AliasMatch, p part, e *sdk.Emitter) error {
	a := m.Alias
	decl, _ := m.Reader().Lookup(a.Target.Target)
	target, _ := decl.(*node.Struct)
	g.write(m, p, e, &subject{origin: a.ID, name: a.Name, decl: target, visibility: a.Visibility})
	return nil
}

// write records the model of the builder of sub into the family of the
// builders, under the name of [KeyName] on sub. It records the checks of the
// builder into the family of the tests, unless the option no-tests is set. p
// is the part of the language of sub. The annotate half stamps the facts on a
// declaration whose builder write can record alone, so write reports nothing.
func (g *generator) write(m sdk.SubjectMatch, p part, e *sdk.Emitter, sub *subject) {
	in := &invocation{options: g.options, m: m, p: p, e: e, sub: sub, b: m.Rules()}
	members := in.members()
	start, seed, values, _ := in.start(members)
	name, _ := sdk.Fact(m, g.name)
	builder := &emit.Struct{
		Origin:     sub.origin,
		Name:       name,
		Visibility: sub.visibility,
		TypeParams: sdk.MirrorTypeParams(sub.params),
	}
	fromSeed, fromValue := p.constructors(builder.Name)
	model := Model{
		Subject: sub.name, New: fromSeed, From: fromValue,
		Value: in.valueOf(sdk.TypeArgs(builder.TypeParams)), Seed: seed,
	}
	model.Struct = in.structOf(model.Value)
	for i := range members {
		model.Members = append(model.Members, in.modelOf(&members[i]))
	}
	builder.Types.Append(start)
	record := []symbol.Symbol{builder}
	if p.functions() {
		record = append(record, in.function(fromSeed), in.function(fromValue))
	}
	e.File().Record(model, record...)
	if g.options.NoTests {
		return
	}
	test, checks := in.tests(builder.Name, members, &model, values)
	if p.functions() {
		test.Fields.Append(referent(fromSeed), referent(fromValue))
	}
	e.File(TestTag).Record(checks, test)
}

// function returns the declaration of the function named name, which the
// template of the builders writes beside the builder. The declaration lets
// the layout of eidos check the name in the package of the builder and
// qualify a reference to it.
func (in *invocation) function(name string) *emit.Function {
	return &emit.Function{Origin: in.sub.origin, Name: name, Visibility: in.sub.visibility}
}

// referent returns a field of the record of the checks whose type refers to
// the function named name, which the record of the builder declares.
func referent(name string) *emit.Field {
	return &emit.Field{Name: name, Type: &emit.TypeRef{Spelling: name}}
}

// start returns the start body of the record of a builder, the seed of the
// builder and the values of the declared defaults by field. The body
// declares the value through the seed function where one exists, and then
// assigns each declared default. start reports [BadDefault] and false at a
// field whose default it cannot read as a value of the type of the field.
func (in *invocation) start(members []member) (*emit.Method, Seed, map[string]emit.Value, bool) {
	sub := in.sub
	start := &emit.Method{Origin: sub.origin, Name: valueLocal, Visibility: symbol.VisibilityPackage}
	seed := SeedZero
	if call, seeded := in.seed(); seeded {
		seed = SeedCompanion
		start.Body.Stmts = append(start.Body.Stmts, emit.Stmt{
			Kind: emit.StmtAssign, Names: []string{valueLocal}, Declare: true,
			Value: emit.Expr{Kind: emit.ExprCall, Fn: &call},
		})
	}
	values := map[string]emit.Value{}
	scope := in.b.View().ScopeOf(sub.decl.ID)
	for i := range members {
		f := members[i].field
		if f == nil {
			continue
		}
		text, declared := sdk.FactOf(in.m, f.ID, meta.Named[string](defaults.KeyValue))
		if !declared {
			continue
		}
		if f.Type == nil {
			in.m.ErrorfAt(BadDefault, f.Pos, "%s does not declare a type, so builder cannot check its default %q",
				f.Name, text)
			return nil, "", nil, false
		}
		value, lifted := in.lift(&scope, f.Type, text)
		if !lifted {
			in.m.ErrorfAt(BadDefault, f.Pos, "the default %q of %s is no value of its type %s",
				text, f.Name, f.Type.Spelling)
			return nil, "", nil, false
		}
		values[f.Name] = value
		start.Body.Stmts = append(start.Body.Stmts, assigned(valueLocal+selector+f.Name, &value))
		if seed == SeedZero {
			seed = SeedDefaults
		}
	}
	return start, seed, values, true
}

// tests returns the record of the checks of the builder named name, with its
// data. The first function of the record seeds each scalar member that has a
// sample, and each further function binds the default of one field to want.
// Each method binds the values of one check: want for the check of a setter,
// first and want for the check of a replacement, and want and other for the
// check of a clone. A method has a parameter of the name and the type of
// each local that its body binds.
//
// A check writes a sample of the argument of the setter of its member: the
// whole list for a list, and the value that the optional contains for an
// optional. A scalar member with an alternate sample also has the check of a
// replacement, which writes the sample and then the alternate. A member that
// a clone copies has the check of the clone where it has an alternate
// sample. That check sets the sample on the original, writes the alternate
// into the clone in place, and compares the original. An embedded field and
// a member without a sample have none of these checks.
//
// Each member has a probe of each method that these checks leave out, and
// each member that a clone copies is a fill of the check of a clone. A probe
// and a fill of a member with a sample bind the method of the check of its
// setter. tests checks a generic declaration at the witnesses of its type
// parameters. For a declaration whose type parameters have no witnesses, it
// returns checks without a check of a member, and Instantiable is false.
func (in *invocation) tests(
	name string, members []member, model *Model, values map[string]emit.Value,
) (*emit.Struct, Checks) {
	sub := in.sub
	checks := Checks{Subject: sub.name, New: model.New, From: model.From}
	test := &emit.Struct{Origin: sub.origin, Name: name, Visibility: sub.visibility}
	witnesses := in.b.Witnesses(sub.params)
	if len(sub.params) > 0 && witnesses == nil {
		return test, checks
	}
	checks.Instantiable = true
	for _, w := range witnesses {
		checks.Args = append(checks.Args, rules.EmitRef(w))
	}
	checks.Value = in.valueOf(checks.Args)
	checks.Struct = in.structOf(checks.Value)
	seed := in.method(seedLocal)
	test.Types.Append(seed)
	for i := range members {
		mem := &members[i]
		ref := rules.Substitute(mem.ref, sub.params, witnesses)
		at := in.memberOf(mem, ref)
		sampled := -1
		if mem.field != nil {
			sampled = in.checkField(test, seed, &checks, mem, &at, ref, values)
		}
		checks.Probes = append(checks.Probes, probes(&at, sampled)...)
		if at.Copied {
			checks.Fills = append(checks.Fills, probe(&at, KindWith, at.Names.With, sampled))
		}
	}
	return test, checks
}

// checkField appends the checks of the field of mem to the record test and
// to checks, and returns the index of the method that binds the sample of the
// check of its setter, or -1 where the field has no sample. at is the member
// at the witnesses, and ref is its type there. seed is the first function of
// the record, and values are the declared defaults by field.
func (in *invocation) checkField(
	test *emit.Struct, seed *emit.Method, checks *Checks, mem *member, at *Member, ref *node.TypeRef,
	values map[string]emit.Value,
) int {
	if value, declared := values[mem.name]; declared {
		// ZeroValue returns the zero value of emit for a type without a
		// zero value, and no default equals it.
		if zero, _ := in.b.ZeroValue(ref); !reflect.DeepEqual(value, zero) {
			check := in.method(in.e.JoinName(mem.name, defaultCheck), wantLocal)
			in.bind(check, rules.EmitRef(ref), assigned(wantLocal, &value))
			test.Types.Append(check)
			checks.Defaults = append(checks.Defaults, mem.name)
		}
	}
	argument := in.argumentOf(at.Form, mem.field, ref)
	sample, alternate := in.b.SamplesOf(mem.id, argument, mem.name)
	if !sample.OK() {
		return -1
	}
	if at.Form == FormScalar {
		seed.Body.Stmts = append(seed.Body.Stmts, assigned(seedLocal+selector+mem.name, &sample.Value))
		checks.Seeds = append(checks.Seeds, mem.name)
	}
	typ := rules.EmitRef(argument)
	check := Check{
		Type: typ, Field: mem.name, Form: at.Form, Setter: at.Names.With, Appender: at.Names.Append,
		Method: test.Methods.Len(), Replace: -1,
	}
	test.Methods.Append(in.bind(in.method(in.e.JoinName(mem.name, setterCheck), wantLocal), typ,
		assigned(wantLocal, &sample.Value)))
	if at.Form == FormScalar && alternate.OK() {
		check.Replace = test.Methods.Len()
		test.Methods.Append(in.bind(in.method(in.e.JoinName(mem.name, replaceCheck), firstLocal, wantLocal), typ,
			assigned(firstLocal, &sample.Value), assigned(wantLocal, &alternate.Value)))
	}
	checks.Setters = append(checks.Setters, check)
	if at.Copied && alternate.OK() {
		clone := check
		clone.Method, clone.Replace = test.Methods.Len(), -1
		test.Methods.Append(in.bind(in.method(in.e.JoinName(mem.name, cloneCheck), wantLocal, otherLocal), typ,
			assigned(wantLocal, &sample.Value), assigned(otherLocal, &alternate.Value)))
		checks.Clones = append(checks.Clones, clone)
	}
	return check.Method
}

// probes returns the probes of the methods of the member at that the checks
// of the setters leave out. sampled is the index of the method that binds the
// sample of the check of its setter, and -1 where the member has none. A
// member without a sample has the probe of its setter. A list has the probe
// of its appender, a map and a set have the probes of the setters of their
// entries or keys, and bytes have the probe of the setter from text where
// the language declares one.
func probes(at *Member, sampled int) []Probe {
	var out []Probe
	if sampled < 0 {
		out = append(out, probe(at, KindWith, at.Names.With, sampled))
	}
	switch at.Form {
	case FormList:
		out = append(out, probe(at, KindAppend, at.Names.Append, sampled))
	case FormMap, FormSet:
		out = append(out, probe(at, KindEntry, at.Names.Entry, sampled),
			probe(at, KindEntries, at.Names.Entries, sampled))
	case FormBytes:
		if at.Names.Text != "" {
			out = append(out, probe(at, KindText, at.Names.Text, sampled))
		}
	default:
		// A scalar and an optional have the setter alone.
	}
	return out
}

// probe returns the probe of the method named method of the member at. kind
// is the kind of that method. sampled is the index of the method that binds
// the sample of the check of its setter, and -1 where the member has none.
func probe(at *Member, kind Kind, method string, sampled int) Probe {
	return Probe{
		Type: at.Type, Elem: at.Elem, Key: at.Key, Value: at.Value, Field: at.Field, Form: at.Form,
		Kind: kind, Method: method, With: at.Names.With, Sample: sampled,
	}
}

// method returns a method of the record of the checks named name, with a
// parameter for each local of locals.
func (in *invocation) method(name string, locals ...string) *emit.Method {
	m := &emit.Method{Origin: in.sub.origin, Name: name, Visibility: symbol.VisibilityPackage}
	for _, local := range locals {
		m.Params.Append(&emit.Param{Name: local})
	}
	return m
}

// bind gives each parameter of the method m the type typ, appends the
// statements to the body of m, and returns m.
func (*invocation) bind(m *emit.Method, typ *emit.TypeRef, stmts ...emit.Stmt) *emit.Method {
	for _, p := range m.Params.Items() {
		p.Type = typ
	}
	m.Body.Stmts = append(m.Body.Stmts, stmts...)
	return m
}

// assigned returns the statement that assigns value to name.
func assigned(name string, value *emit.Value) emit.Stmt {
	return emit.Stmt{Kind: emit.StmtAssign, Names: []string{name}, Value: emit.ValueExpr(*value)}
}

// constructible reports whether a builder can make a value of the struct of
// the subject without arguments. In a language without constructors, every
// type has such a value. In a language with constructors, the struct needs a
// public constructor whose parameters are all optional. A parameter with a
// default, a parameter with the optional mark of its language and a rest
// parameter are optional. An inherited or an implicit constructor counts.
func (in *invocation) constructible() bool {
	c, constructs := in.b.Source().(rules.ConstructRules)
	if !constructs {
		return true
	}
	for _, ctor := range c.Constructors(in.sub.decl, in.b.View()) {
		if !slices.ContainsFunc(ctor.Params, func(p rules.ParamView) bool { return !p.Optional }) {
			return true
		}
	}
	return false
}

// seed returns the call of the function that seeds the value of a new
// builder. That function is the function of [DefaultsParam], or else the
// companion function of the declaration. The name of a companion joins the
// name of the subject and the option companion-word. A companion is a
// function without parameters in the package of the subject, and its one
// result is the value of the builder. seed reports false where neither
// function exists.
func (in *invocation) seed() (emit.Expr, bool) {
	if param, named := in.m.Directive().Params[DefaultsParam]; named {
		return emit.RefExpr(param.Target), true
	}
	origin := in.sub.origin
	name := in.p.companion(in.sub.name, cmp.Or(in.options.CompanionWord, CompanionWord))
	id := symbol.Identity{Lang: origin.Lang, Package: origin.Package, Name: name, Kind: symbol.KindFunction}
	decl, _ := in.m.Reader().Lookup(id)
	fn, isFunction := decl.(*node.Function)
	if !isFunction || len(fn.Params) > 0 || len(fn.Returns) != 1 ||
		fn.Returns[0].Type == nil || fn.Returns[0].Type.Target != origin {

		return emit.Expr{}, false
	}
	return emit.RefExpr(id), true
}

// members returns the fields that the builder sets. They are the fields that
// the struct declares and its embedded fields, in the order of the source,
// and then the fields that a base declares, in the order of the member walk
// of the language. A field that an embed promotes is no member, because the
// builder sets the embedded field whole. A field that code outside the
// package cannot assign on a value is no member, and neither is a field with
// the tag [SkipTag] of the value [SkipValue]. members leaves out a field
// whose name the language does not write after a dot, such as the quoted
// TypeScript property content-type and the computed name
// [Symbol.toStringTag].
func (in *invocation) members() []member {
	set, _ := in.b.MembersOf(in.sub.decl)
	// The walk lists the members that the struct declares first, its fields
	// before its embedded fields.
	declared := set.Members
	if at := slices.IndexFunc(declared, func(m rules.Member) bool { return m.Depth > 0 }); at >= 0 {
		declared = declared[:at]
	}
	slices.SortStableFunc(declared, func(a, b rules.Member) int {
		return a.Symbol.Position().Compare(b.Symbol.Position())
	})
	out := make([]member, 0, len(set.Members))
	for i := range set.Members {
		mem := &set.Members[i]
		if in.promoted(mem) {
			continue
		}
		switch s := mem.Symbol.(type) {
		case *node.Field:
			if assignable(s) && in.p.selectable(s.Name) && !in.skipped(s) {
				out = append(out, member{field: s, ref: s.Type, name: s.Name, id: s.ID})
			}
		case *node.Embed:
			if in.exposed(s) {
				out = append(out, member{ref: s.Ref, name: s.ID.Name, id: s.ID})
			}
		}
	}
	return out
}

// exposed reports whether code outside the package sets the embedded field
// e. The field has the name of its type, so code outside the package sets it
// where the declaration of that type is public. An embed of a type outside
// the view, such as a predeclared interface, is no member.
func (in *invocation) exposed(e *node.Embed) bool {
	decl, _ := in.m.Reader().Lookup(contained(e.Ref).Target)
	switch d := decl.(type) {
	case *node.Struct:
		return d.Visibility == symbol.VisibilityPublic
	case *node.Interface:
		return d.Visibility == symbol.VisibilityPublic
	case *node.Alias:
		return d.Visibility == symbol.VisibilityPublic
	default:
		return false
	}
}

// promoted reports whether an embed of the struct of the subject promotes
// the member mem. A member below the struct arrives through an embed or
// through a base. A language of embeds, such as Go, has no bases, and a
// language of bases, such as TypeScript, has no embeds. So a member below a
// struct with embeds arrives through an embed.
func (in *invocation) promoted(mem *rules.Member) bool {
	return mem.Depth > 0 && len(in.sub.decl.Embeds) > 0
}

// skipped reports whether the field f has the tag [SkipTag] of the value
// [SkipValue]. A field without the tag has the empty value.
func (in *invocation) skipped(f *node.Field) bool {
	text, _ := in.b.Tag(f, SkipTag)
	return text == SkipValue
}

// mistagged reports [BadTag] at the field of each of members with the tag
// [SkipTag]. The members leave out each field whose tag has the value
// [SkipValue], so the tag of a member has another value, and the builder
// keeps the field. An embedded member has no field, and the rules report no
// tag of a nil field.
func (in *invocation) mistagged(members []member) {
	for i := range members {
		f := members[i].field
		if text, tagged := in.b.Tag(f, SkipTag); tagged {
			in.m.ErrorfAt(BadTag, f.Pos, "the tag %s of %s has the value %q, and only %q leaves a field out of "+
				"the builder", SkipTag, f.Name, text, SkipValue)
		}
	}
}

// assignable reports whether code outside the package can assign the field
// f on a value. Such a field is public, belongs to the instance and is not
// immutable. A field with a hard private name is private.
func assignable(f *node.Field) bool {
	return f.Visibility == symbol.VisibilityPublic && f.Level != symbol.LevelType &&
		f.Mutability != symbol.MutabilityImmutable
}

// modelOf returns the model of the member mem at the type that its source
// declares.
func (in *invocation) modelOf(mem *member) Member {
	return in.memberOf(mem, mem.ref)
}

// memberOf returns the model of the member mem whose type is typ, such as
// the type that its source declares or that type at the witnesses of the type
// parameters of the subject. An embedded field is a scalar, which the builder
// sets whole. A field that may be absent, such as a TypeScript property with
// a question mark, is an optional of its type. For any other field, the part
// of the language reports whether the type is a set. Otherwise the form of
// the member is the form of the canonical shape of the type: a list, a map,
// bytes, an optional or a scalar. The type follows each transparent alias,
// so an alias of a list is a list. A Go defined type is a reference to its
// declaration, and so a scalar. A clone copies each list, map, set and bytes
// value that is not immutable. The part of the language names the methods
// of the member from its field and its form.
func (in *invocation) memberOf(mem *member, typ *node.TypeRef) Member {
	out := Member{Field: mem.name, Form: FormScalar, Type: rules.EmitRef(typ)}
	if f := mem.field; f != nil {
		ref := rules.Unaliased(typ, in.b.View())
		shape := in.b.TypeOf(ref)
		children := childrenOf(ref)
		key, isSet := in.p.set(ref)
		switch {
		case f.Optional:
			out.Form, out.Elem = FormOptional, rules.EmitRef(typ)
		case isSet:
			out.Form, out.Elem = FormSet, rules.EmitRef(key)
		case shape.Form == symbol.FormBytes:
			out.Form = FormBytes
		case shape.Form == symbol.FormList:
			out.Form, out.Elem = FormList, rules.EmitRef(children[0])
		case shape.Form == symbol.FormMap:
			out.Form, out.Key, out.Value = FormMap, rules.EmitRef(children[0]), rules.EmitRef(children[1])
		case shape.Form == symbol.FormOptional:
			out.Form, out.Elem = FormOptional, rules.EmitRef(children[0])
		}
		switch out.Form {
		case FormList, FormMap, FormSet, FormBytes:
			out.Copied = shape.Mutability != symbol.MutabilityImmutable
		default:
			// A scalar and an optional are values that the builder sets whole.
		}
	}
	out.Names = in.p.names(mem.name, out.Form)
	return out
}

// collision returns the first name that two methods of the members share,
// with the fields of the two members, and false where each method has a name
// of its own. A name is empty where the form of a member has no such method.
func collision(members []Member) (name, first, second string, collides bool) {
	seen := map[string]string{}
	for i := range members {
		n := &members[i].Names
		for _, method := range []string{n.With, n.Append, n.Entry, n.Entries, n.Text} {
			if method == "" {
				continue
			}
			if field, taken := seen[method]; taken {
				return method, field, members[i].Field, true
			}
			seen[method] = members[i].Field
		}
	}
	return "", "", "", false
}

// argumentOf returns the type of the value that a check writes through the
// setter of the field f of the form form. ref is the type of f. The setter
// of an optional of the type of the field, such as a Go pointer, takes the
// value that the optional contains. The setter of every other form takes a
// value of the type of the field, and the setter of a list takes the
// elements of that value.
func (in *invocation) argumentOf(form Form, f *node.Field, ref *node.TypeRef) *node.TypeRef {
	if form == FormOptional && !f.Optional {
		return childrenOf(rules.Unaliased(ref, in.b.View()))[0]
	}
	return ref
}

// lift returns the value of the text of the default of a field of the type
// ref. The rules read the text as a value in the language of the field, in
// the scope of the declaration. Such a value is a literal, or the name of a
// declaration or a call in a language that reads them. For a field of text,
// a text that the rules do not read is the content of the string, as the
// text of an authored sample is. lift reports false for a text that is no
// value of the type.
func (in *invocation) lift(scope *rules.Scope, ref *node.TypeRef, text string) (emit.Value, bool) {
	value, read := in.b.LiteralFor(*scope, ref, text)
	if read || !in.textual(ref) {
		return value, read
	}
	return emit.Literal(emit.LiteralString, text), true
}

// textual reports whether the type ref is text. A type that refers to a
// declared type, such as type Code string of Go, is text where the target
// of that declaration is text.
func (in *invocation) textual(ref *node.TypeRef) bool {
	shape := in.b.TypeOf(ref)
	for shape.Form == symbol.FormReference {
		decl, _ := in.m.Reader().Lookup(shape.Ref)
		alias, declared := decl.(*node.Alias)
		if !declared {
			return false
		}
		shape = in.b.TypeOf(alias.Target)
	}
	return shape.Form == symbol.FormText
}

// valueOf returns the type of the value of the builder with the type
// arguments args.
func (in *invocation) valueOf(args []*emit.TypeRef) *emit.TypeRef {
	return &emit.TypeRef{Spelling: in.sub.name, Target: in.sub.origin, Args: args}
}

// structOf returns the struct that the builder constructs, where value is
// the type of the value of the builder. It is value for a struct and the
// struct of the alias for an alias. The builder of an alias has no type
// parameters.
func (in *invocation) structOf(value *emit.TypeRef) *emit.TypeRef {
	if in.sub.origin == in.sub.decl.ID {
		return value
	}
	return &emit.TypeRef{Spelling: in.sub.decl.Name, Target: in.sub.decl.ID}
}

// childrenOf returns the references that the fold of ref makes the children
// of its shape. They are the elements of a structural reference, or the type
// arguments of a named reference, such as string in Array<string>.
func childrenOf(ref *node.TypeRef) []*node.TypeRef {
	if ref == nil {
		return nil
	}
	if len(ref.Elems) > 0 {
		return ref.Elems
	}
	return ref.Args
}

// contained returns the reference that the optional reference ref contains,
// and ref itself where it is no optional. An optional has one element.
func contained(ref *node.TypeRef) *node.TypeRef {
	if ref.Form == symbol.FormOptional {
		return ref.Elems[0]
	}
	return ref
}
