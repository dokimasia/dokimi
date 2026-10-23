// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"cmp"
	"errors"
	"slices"
	"strings"

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
const ID plugin.ID = "enum"

// Version is the version of the plugin. The fingerprint of a composition
// includes it.
const Version = "0.1.0"

// EnumDirective puts an enumeration under enum. MethodsParam is the param of
// the directive whose value MethodsOff leaves out the surface and keeps the
// checks, as in methods=off. TextParam is the param whose value TextComment
// takes the text of each variant from its trailing comment, as in
// text=comment.
const (
	EnumDirective directive.Name     = "enum"
	MethodsParam  directive.ParamKey = "methods"
	MethodsOff                       = "off"
	TextParam     directive.ParamKey = "text"
	TextComment                      = "comment"
)

// Word is the word of the family of the surfaces and of the family of their
// checks. TestTag is the tag of the family of the checks.
const (
	Word    = "enum"
	TestTag = sdk.TagTest
)

// DefaultParseFunction, DefaultValuesFunction, DefaultParseError and
// DefaultValidMethod are the defaults of the options parse-function,
// values-function, parse-error and valid-method. For the enumeration Status
// they render ParseStatus, StatusValues, ErrUnknownStatus and IsValid in Go.
const (
	DefaultParseFunction  NameTemplate   = "Parse{{.Enum}}"
	DefaultValuesFunction NameTemplate   = "{{.Enum}}Values"
	DefaultParseError     NameTemplate   = "Unknown{{.Enum}}"
	DefaultValidMethod    MethodTemplate = "IsValid"
)

// The names of the fields of the record of the checks that refer to the
// parse function, to the list of the variants and to the error of the parse
// function, in a language whose checks refer to them through their package.
const (
	parseField  = "parse"
	valuesField = "values"
	errorField  = "error"
)

// The tree of each target contains the file template of the surfaces and
// the file template of their checks.
const (
	enumTemplate = "enum.tmpl"
	testTemplate = "enum.test.tmpl"
)

// unknownMark is the character that a text of no variant repeats.
const unknownMark = "?"

// Options are the options of enum. ParseFunction, ValuesFunction and
// ParseError are the name templates of the parse function, of the function
// that lists the variants and of the error of the parse function.
// ValidMethod is the method template of the validity method of a Go
// enumeration. TypeScript names its type guard after the enum, as isStatus,
// and reads no method template. An empty template takes its default. NoTests
// leaves out the checks of every enumeration.
type Options struct {
	ParseFunction  NameTemplate   `opt:"parse-function" doc:"The name template of the function that parses a text into a variant. The default Parse{{.Enum}} renders ParseStatus in Go and parseStatus in TypeScript."`
	ValuesFunction NameTemplate   `opt:"values-function" doc:"The name template of the function that lists the variants. The default {{.Enum}}Values renders StatusValues in Go and statusValues in TypeScript."`
	ParseError     NameTemplate   `opt:"parse-error" doc:"The name template of the error that the parse function returns for a text that matches no variant. The default Unknown{{.Enum}} renders ErrUnknownStatus in Go and UnknownStatusError in TypeScript."`
	ValidMethod    MethodTemplate `opt:"valid-method" doc:"The method template of the method of a Go enumeration that reports whether a value is a variant. The default is IsValid."`
	NoTests        bool           `opt:"no-tests" doc:"When true, the generator writes no checks."`
}

// baseNames are the names that the templates of the options render for one
// enumeration, before the part of a language writes each name in its case
// convention. parse names the parse function, values the function that lists
// the variants, err the error of the parse function and valid the validity
// method.
type baseNames struct {
	parse, values, err, valid string
}

// Names are the names of the members of the surface of one enumeration, in
// the conventions of the language of the enumeration. Render returns the text
// of a value, and Parse returns the value of a text or the error Sentinel.
// Values returns the variants, and Valid reports whether a value is a
// variant. Encode and Decode are the pair that writes a value as text and
// reads it back, and they are empty in a language without such a pair.
type Names struct {
	Render   string
	Parse    string
	Values   string
	Valid    string
	Encode   string
	Decode   string
	Sentinel string
}

// Variant is one variant of an enumeration. Name is the name of the variant.
// Member is the reference to the variant in the language of the
// enumeration, which the template writes after the enumeration where the
// language refers to a variant through it, as in Status.Draft of TypeScript.
// Value is the name of the declaration of the variant as a value, which the
// builtin value of a template writes with the import of its package. Text is
// the text of the variant as a string literal of the language. Stated
// reports that the declaration states the text: the value of a variant of an
// enumeration of text, or the trailing comment of a variant under
// text=comment.
type Variant struct {
	Value  emit.Value
	Name   string
	Member string
	Text   string
	Stated bool
}

// Members are the members of the surface that a record declares or that a
// check calls. Error is the error of the parse function.
type Members struct {
	Render bool
	Parse  bool
	Error  bool
	Values bool
	Valid  bool
	Encode bool
	Decode bool
}

// member is one member of the surface of an enumeration.
type member uint8

// The members of a surface, in the order in which declares looks for them.
const (
	memberRender member = iota
	memberParse
	memberError
	memberValues
	memberValid
	memberEncode
	memberDecode
)

// members are the members of a surface in their order.
var members = []member{memberRender, memberParse, memberError, memberValues, memberValid, memberEncode, memberDecode}

// labels name each member in the message of a [Conflict].
var labels = [...]string{
	memberRender: "render",
	memberParse:  "parse function",
	memberError:  "error of the parse function",
	memberValues: "list of the variants",
	memberValid:  "validity method",
	memberEncode: "encode",
	memberDecode: "decode",
}

// kind is the kind of declaration that a member of a surface is: a function of
// the package, a method of the enumeration, a variable of the package, or a
// class of the module.
type kind uint8

// The kinds of the members of a surface.
const (
	kindFunction kind = iota + 1
	kindMethod
	kindVariable
	kindClass
)

// shape is the shape of an input or of a value of a member of a surface: any
// value, a text, a truth value, bytes, a variant of the enumeration, or a list
// of its variants.
type shape uint8

// The shapes of the inputs and the values of the members of a surface.
const (
	shapeAny shape = iota
	shapeText
	shapeBool
	shapeBytes
	shapeVariant
	shapeVariants
)

// form is what a declaration under the name of a member of a surface is: its
// kind, the shapes of its inputs and of its values, and whether it returns an
// error. text describes the form in the message of a [Conflict].
type form struct {
	text   string
	inputs []shape
	values []shape
	kind   kind
	fails  bool
}

// fits reports whether the declaration d has the form f. enum is the
// enumeration of the surface, whose variants a variant shape takes. b are the
// rules of the language of d.
func (f *form) fits(b *rules.Bound, d node.Declaration, enum *symbol.Identity) bool {
	switch d.(type) {
	case *node.Variable:
		return f.kind == kindVariable
	case *node.Struct:
		return f.kind == kindClass
	case *node.Function:
		if f.kind != kindFunction {
			return false
		}
	case *node.Method:
		// declares looks a method up for a member of the form of a method
		// alone, so f is the form of a method.
	default:
		return false
	}
	// A check calls a member with its inputs alone, so every parameter
	// counts, a context included.
	view, _ := b.CallableOf(d)
	inputs := make([]*node.TypeRef, 0, len(view.Params))
	for _, p := range view.Params {
		inputs = append(inputs, p.Ref)
	}
	values := make([]*node.TypeRef, 0, len(view.Returns))
	fails := false
	for _, r := range view.Returns {
		if r.Role == rules.ReturnError {
			fails = true
			continue
		}
		values = append(values, r.Ref)
	}
	return fails == f.fails && shaped(b, inputs, f.inputs, enum) && shaped(b, values, f.values, enum)
}

// shaped reports whether the types refs have the shapes want, in order.
func shaped(b *rules.Bound, refs []*node.TypeRef, want []shape, enum *symbol.Identity) bool {
	if len(refs) != len(want) {
		return false
	}
	for i, ref := range refs {
		if t := b.TypeOf(ref); !want[i].admits(&t, enum) {
			return false
		}
	}
	return true
}

// admits reports whether the type t has the shape s. A variant shape admits a
// reference to the enumeration enum, the one form whose referent is enum. A
// list has one element type.
func (s shape) admits(t *rules.TypeShape, enum *symbol.Identity) bool {
	switch s {
	case shapeText:
		return t.Form == symbol.FormText
	case shapeBool:
		return t.Form == symbol.FormBool
	case shapeBytes:
		return t.Form == symbol.FormBytes
	case shapeVariant:
		return t.Ref == *enum
	case shapeVariants:
		return t.Form == symbol.FormList && shapeVariant.admits(&t.Elems[0], enum)
	default:
		return true
	}
}

// Model is the data of the record of the surface of one enumeration, and the
// file template of the surfaces reads it. Subject is the name of the
// enumeration and Package the name of its package, which opens the message of
// Sentinel. Type is the type of the enumeration. Fallback is the type that a
// value outside the variants converts to before the render prints it with
// the verb Format, and it is nil for an enumeration of text, which renders
// such a value as itself. Generate are the members that the record declares,
// which are the members that the enumeration and its package do not declare.
// ErrorRef refers to the error of the parse function that the package
// declares, and it is nil where the surface writes the error.
//
// Table is the table of the texts of an enumeration of integers, which the
// render slices. It is nil where the part of the language writes no table,
// such as for an enumeration of another kind. Converts reports that the text
// of each variant is its value, so the render converts a value to text.
//
// The record of a surface names the struct of the enumeration. In a language
// whose checks refer to the surface through its package, the record also
// declares the parse function, the list of the variants and the error of the
// parse function that the surface writes beside that struct, and the template
// writes them at the struct.
type Model struct {
	Type     *emit.TypeRef
	Fallback *emit.TypeRef
	ErrorRef *emit.TypeRef
	Table    *Table
	Subject  string
	Package  string
	Format   string
	Names    Names
	Variants []Variant
	Generate Members
	Textual  bool
	Converts bool
}

// Table is the table of the texts of an enumeration of integers, from which
// the render slices the text of a value. Values are the variants with their
// values in the order of the declaration, which a check of the values at
// compile time reads. Runs are the runs of consecutive values in increasing
// order. Indexed reports that a run has two values or more, so the render
// declares the index of that run. Map reports that the runs are too many for
// a switch, so the render looks the text of a value up in a map. Name is then
// the texts of the variants in the order of their values in one string
// literal, and Entries are the variants in that order with the bounds of
// their texts in Name.
type Table struct {
	Name    string
	Values  []Entry
	Runs    []Run
	Entries []Entry
	Indexed bool
	Map     bool
}

// Entry is one variant of a table: its name, its value as a literal, and the
// bounds of its text in the name of the table.
type Entry struct {
	Name  string
	Value string
	Start int
	End   int
}

// Run is one run of consecutive values of a table. Name is the texts of its
// variants in one string literal, and Index the offset of each text in Name,
// followed by the length of Name. Bits is the width of the unsigned integers
// of the index. Low and High are the first and the last value, as literals.
// Open reports that the run starts at the zero of an unsigned type, so a
// value of the run has no lower bound to check. Single reports that the run
// has one value, whose text is all of Name.
type Run struct {
	Name   string
	Low    string
	High   string
	Index  []int
	Bits   int
	Open   bool
	Single bool
}

// Checks is the data of the record of the checks of one enumeration, and the
// file template of the checks reads it. Zero is the variant whose value is
// the zero value of its type, and it is nil where no variant has that value.
// Unknown is a text that no variant has, as a string literal of the
// language. OutOfRange is a literal of the language of a value outside the
// variants, which the template converts to the enumeration, and it is empty
// where the rules derive none. Calls are the members that a check calls: a
// member that the record of the surface declares, or that the enumeration or
// its package declares. A check calls the parse function where a render
// exists, and the error where it calls the parse function. Generate are the
// members that the record of the surface declares. Stated reports that a
// variant states its text, and the checks pin the text of each such variant.
//
// In a language whose checks refer to the parse function, the list of the
// variants and the error through their package, the record of the checks has
// a field of each that a check calls, named parse, values and error. The
// type of the field refers to the declaration of the surface by its name, and
// to a declaration of the package by its identity. The layout of eidos
// qualifies each reference where the checks are in another package than the
// surface.
type Checks struct {
	Type       *emit.TypeRef
	Zero       *Variant
	Subject    string
	Unknown    string
	OutOfRange string
	Names      Names
	Variants   []Variant
	Calls      Members
	Generate   Members
	Stated     bool
}

// part is the part of enum for one language. It declares the target with its
// template tree and its helpers, and it names the members of the surface.
type part interface {
	lang.Part
	// templated returns the name of the enumeration named subject that the
	// name templates of the options render, before names writes each name in
	// the case convention of the language.
	templated(subject string) string
	// names returns the names of the members of the surface of the
	// enumeration named subject in the case conventions of the language. n
	// are the names that the templates of the options render for the
	// enumeration, and b are the rules of its language.
	names(b *rules.Bound, subject string, n baseNames) Names
	// form returns the form of a declaration under the name of the member x,
	// and the zero form for a member that the language does not have.
	form(x member) form
	// fallback returns the type that a value outside the variants of e
	// converts to before the render prints it, and the verb that prints it.
	// It returns nil for an enumeration of text, and for a language whose
	// render prints such a value without a conversion.
	fallback(b *rules.Bound, e *node.Enum, textual bool) (*emit.TypeRef, string)
	// member returns the reference to the variant named name, which the
	// template writes after the enumeration where the language refers to a
	// variant through it.
	member(name string) string
	// literal returns the literal of the language of the value v: a string
	// in quotes, and a number or a truth value as its text.
	literal(v *emit.Value) string
	// qualified reports whether a check refers to a function or a variable
	// of the surface through the package of the surface, and false where a
	// check imports each of them from the module of the surface.
	qualified() bool
	// table returns the table of the texts of the enumeration e, whose
	// variants have the texts of variants and the identities of ids. It
	// returns nil where the language writes the render of e without a table.
	// m is the match of e, and b are the rules of its language.
	table(m sdk.Matcher, b *rules.Bound, e *node.Enum, variants []rules.VariantText,
		ids map[string]symbol.Identity) *Table
}

// parts are the parts of enum, one for each language that it serves.
var parts = []part{golangPart{}, typescriptPart{}}

// generator is one instance of enum. The composition writes its options
// before the first phase call, and the handlers read them afterwards.
//
// # Concurrency
//
// The workers of a phase call read the options concurrently, and no handler
// writes them.
type generator struct {
	options *Options
}

// schema is the directive of enum.
var schema = directive.Schema{
	Name: EnumDirective,
	Params: []directive.ParamSpec{
		{
			Key: MethodsParam, Type: directive.TypeString, Choices: []string{MethodsOff},
			Doc: "off leaves out the surface and keeps the checks",
		},
		{
			Key: TextParam, Type: directive.TypeString, Choices: []string{TextComment},
			Doc: "comment takes the text of each variant from its trailing comment",
		},
	},
	Doc: "generates the text surface of the enumeration and the checks of the surface",
}

// New returns enum. enum records a [Model] into the family [Word] for each
// enumeration under [EnumDirective], and the template of the family writes
// the surface of the enumeration in the conventions of the language of the
// plan. The surface renders a value as text, parses a text back, lists the
// variants and reports whether a value is a variant. In Go it also encodes a
// value as text and decodes it. enum records the [Checks] of each exported
// enumeration into the family [Word] under the tag [TestTag]. The checks
// test the enumeration through the exported API of its package, so a Go file
// of the family is in the external test package of its directory. An
// enumeration that its package does not export has no checks, and the
// package tests it through its own API.
//
// The surface writes each member that the enumeration or its package does not
// declare under its name. It writes a parse function only beside a render that
// it writes, because a declared render can write other texts, and its parse
// function wraps the error that the package declares. A decode needs a parse
// function, and an encode needs a decode. The value off of [MethodsParam]
// leaves out the whole surface and keeps the checks.
//
// The text of a variant of an enumeration of text is its value, and the text
// of any other variant is its name, as the rules of its language derive them.
// Under the value TextComment of [TextParam], the trailing comment of a
// variant is its text, as in Alpha Grade = iota // alpha. The checks pin each
// text that a value or a comment states. enum serves Go and TypeScript, the
// languages of its parts. It reports [Empty] at an enumeration without a
// variant, and [Duplicate] at an enumeration whose two variants have one value
// or one text. It reports [Conflict] at a declaration of the package that has
// the name of a member and another form. enum records nothing for an
// enumeration that it reports. Each call returns a new plugin.
func New() plugin.Generator {
	g := &generator{options: &Options{
		ParseFunction: DefaultParseFunction, ValuesFunction: DefaultValuesFunction, ParseError: DefaultParseError,
		ValidMethod: DefaultValidMethod,
	}}
	b := sdk.NewPlugin(ID).
		Version(Version).
		Options(g.options).
		Output(plugin.Output{Per: plugin.PerSource, Word: Word, Template: enumTemplate}).
		Output(plugin.Output{Tag: TestTag, Per: plugin.PerSource, Word: Word, Template: testTemplate, External: true})
	lang.Register(b, parts)
	return b.Handle(sdk.Directive(schema, sdk.OnEnum(lang.Served(parts, g.generate)))).Generator()
}

// generate records the surface and the checks of an enumeration under the
// directive.
func (g *generator) generate(m *sdk.EnumMatch, p part, e *sdk.Emitter) error {
	enum := m.Enum
	b := m.Rules()
	//nolint:forcetypeassert // the rules of Go and of TypeScript, the languages of the parts, project enumerations
	info := b.Source().(rules.EnumRules).EnumOf(enum, b.View())
	if len(info.Variants) == 0 {
		m.Errorf(Empty, "%s has no variant, so it has no surface", enum.Name)
		return nil
	}
	textual := info.Form == rules.EnumValue
	stated, commented := stateTexts(m, info.Variants, textual)
	if first, second, shared := sharedText(info.Variants); shared {
		m.Errorf(Duplicate, "the variants %s and %s of %s have one text, so it has no surface",
			first, second, enum.Name)
		return nil
	}
	if info.Duplicate != "" {
		m.Errorf(Duplicate, "a variant of %s after %s has the value of %s, so it has no surface",
			enum.Name, info.Duplicate, info.Duplicate)
		return nil
	}
	base, err := g.baseNames(p.templated(enum.Name))
	if err != nil {
		return err
	}
	names := p.names(&b, enum.Name, base)
	pkg, _ := m.Reader().PackageOf(enum.ID.PackageIdentity())
	found, fits := declares(m, p, &b, pkg, &names)
	if !fits {
		return nil
	}
	has := func(x member) bool { _, declared := found[x]; return declared }
	declared := Members{
		Render: has(memberRender), Parse: has(memberParse), Error: has(memberError), Values: has(memberValues),
		Valid: has(memberValid), Encode: has(memberEncode), Decode: has(memberDecode),
	}
	generated := surface(m, &names, &declared)
	typ := &emit.TypeRef{Spelling: enum.Name, Target: enum.ID}
	fallback, format := p.fallback(&b, enum, textual)
	ids := make(map[string]symbol.Identity, len(enum.Variants))
	for _, v := range enum.Variants {
		ids[v.Name] = v.ID
	}
	variants := make([]Variant, 0, len(info.Variants))
	longest := 0
	for i := range info.Variants {
		v := &info.Variants[i]
		variants = append(variants, Variant{
			Value: emit.Name(ids[v.Name]), Name: v.Name, Member: p.member(v.Name), Text: p.literal(&v.Text),
			Stated: stated[i],
		})
		longest = max(longest, len(v.Text.Text))
	}
	model := Model{
		Type: typ, Fallback: fallback, Table: p.table(m, &b, enum, info.Variants, ids), Subject: enum.Name,
		Package: pkg.Name, Format: format, Names: names, Variants: variants, Generate: generated,
		Textual: textual, Converts: textual && !commented,
	}
	if d, declaredError := found[memberError]; declaredError {
		model.ErrorRef = &emit.TypeRef{Spelling: names.Sentinel, Target: d.Identity()}
	}
	if generated != (Members{}) {
		record := []symbol.Symbol{&emit.Struct{Origin: enum.ID, Name: enum.Name, Visibility: enum.Visibility}}
		// The record declares each function and variable of the surface that
		// the checks refer to through its package.
		if p.qualified() {
			if generated.Parse {
				record = append(record,
					&emit.Function{Origin: enum.ID, Name: names.Parse, Visibility: enum.Visibility})
			}
			if generated.Values {
				record = append(record,
					&emit.Function{Origin: enum.ID, Name: names.Values, Visibility: enum.Visibility})
			}
			if generated.Error {
				record = append(record,
					&emit.Variable{Origin: enum.ID, Name: names.Sentinel, Visibility: enum.Visibility})
			}
		}
		e.File().Record(model, record...)
	}
	// The checks test the enumeration from the external test package, which
	// cannot refer to an enumeration that its package does not export.
	if g.options.NoTests || enum.Visibility != symbol.VisibilityPublic {
		return nil
	}
	render := generated.Render || declared.Render
	parse := render && (generated.Parse || declared.Parse)
	unknown := emit.Literal(emit.LiteralString, strings.Repeat(unknownMark, longest+1))
	checks := Checks{
		Type: typ, Subject: enum.Name, Names: names, Variants: variants, Generate: generated,
		Stated: slices.Contains(stated, true), Unknown: p.literal(&unknown),
		Calls: Members{
			Render: render, Parse: parse, Error: parse && (generated.Error || declared.Error),
			Values: generated.Values || declared.Values, Valid: generated.Valid || declared.Valid,
			Encode: generated.Encode || declared.Encode, Decode: generated.Encode && generated.Decode,
		},
	}
	if inner := info.OutOfRange.Inner; inner != nil {
		checks.OutOfRange = p.literal(inner)
	}
	if at := slices.IndexFunc(variants, func(v Variant) bool { return v.Name == info.Zero }); at >= 0 {
		checks.Zero = &variants[at]
	}
	test := &emit.Struct{Origin: enum.ID, Name: enum.Name, Visibility: enum.Visibility}
	if p.qualified() {
		refer := func(field string, x member, name string, called bool) {
			if !called {
				return
			}
			// A generated member is a declaration of the record of the
			// enumeration, which an emitted reference finds by its origin.
			ref := &emit.TypeRef{Spelling: name, Target: enum.ID, Emitted: true}
			if d, declaredHere := found[x]; declaredHere {
				ref.Target, ref.Emitted = d.Identity(), false
			}
			test.Fields.Append(&emit.Field{Name: field, Type: ref})
		}
		refer(parseField, memberParse, names.Parse, checks.Calls.Parse)
		refer(valuesField, memberValues, names.Values, checks.Calls.Values)
		refer(errorField, memberError, names.Sentinel, checks.Calls.Error)
	}
	e.File(TestTag).Record(checks, test)
	return nil
}

// baseNames returns the names that the templates of the options of g render
// for the enumeration named enum. An empty option takes its default. It
// returns the error of each template that does not render, such as a
// template that a program sets without its UnmarshalText.
func (g *generator) baseNames(enum string) (baseNames, error) {
	o := g.options
	parse, parseErr := render(string(cmp.Or(o.ParseFunction, DefaultParseFunction)), enum)
	values, valuesErr := render(string(cmp.Or(o.ValuesFunction, DefaultValuesFunction)), enum)
	unknown, unknownErr := render(string(cmp.Or(o.ParseError, DefaultParseError)), enum)
	valid, validErr := render(string(cmp.Or(o.ValidMethod, DefaultValidMethod)), enum)
	return baseNames{parse: parse, values: values, err: unknown, valid: valid},
		errors.Join(parseErr, valuesErr, unknownErr, validErr)
}

// stateTexts sets the text of each variant with a trailing comment under
// text=comment to the comment without the spaces around it. A variant without
// a comment keeps the text that the rules derive. stateTexts reports for each
// variant whether its declaration states its text. A variant of an
// enumeration of text states its value, which the rules derive as its text.
// stateTexts also reports whether a comment states the text of a variant.
func stateTexts(m *sdk.EnumMatch, variants []rules.VariantText, textual bool) (stated []bool, commented bool) {
	_, byComment := m.Directive().Params[TextParam]
	comments := make(map[string]string, len(m.Enum.Variants))
	for _, v := range m.Enum.Variants {
		comments[v.Name] = strings.TrimSpace(v.Comment)
	}
	stated = make([]bool, len(variants))
	for i := range variants {
		stated[i] = textual
		if comment := comments[variants[i].Name]; byComment && comment != "" {
			variants[i].Text = emit.Literal(emit.LiteralString, comment)
			stated[i] = true
			commented = true
		}
	}
	return stated, commented
}

// sharedText returns the first variant whose text a later variant repeats,
// with that later variant, and false where each variant has a text of its
// own.
func sharedText(variants []rules.VariantText) (first, second string, shared bool) {
	texts := map[string]string{}
	for i := range variants {
		v := &variants[i]
		if earlier, seen := texts[v.Text.Text]; seen {
			return earlier, v.Name, true
		}
		texts[v.Text.Text] = v.Name
	}
	return "", "", false
}

// declares returns the declaration of each member of the surface that the
// enumeration of m or its package pkg declares under the name that n states.
// A member of the form of a method is a method of the enumeration, and any
// other member is a function, a variable, a constant or a type of the package.
// A member that the language does not have has an empty name, which no
// declaration has. declares reports [Conflict] at each declaration of another
// form than the form of its member in the part p, and then returns false. b
// are the rules of the language of the enumeration.
func declares(
	m *sdk.EnumMatch, p part, b *rules.Bound, pkg *node.Package, n *Names,
) (map[member]node.Declaration, bool) {
	scope := map[string]node.Declaration{}
	for _, f := range pkg.Files {
		for _, s := range f.Decls {
			//nolint:forcetypeassert // every kind of the node model of eidos is a Declaration
			if d := s.(node.Declaration); d.Identity().Owner == "" {
				scope[d.Identity().Name] = d
			}
		}
	}
	methods := make(map[string]node.Declaration, len(m.Enum.Methods))
	for _, x := range m.Enum.Methods {
		methods[x.Name] = x
	}
	names := map[member]string{
		memberRender: n.Render, memberParse: n.Parse, memberError: n.Sentinel, memberValues: n.Values,
		memberValid: n.Valid, memberEncode: n.Encode, memberDecode: n.Decode,
	}
	found := map[member]node.Declaration{}
	fits := true
	for _, x := range members {
		f := p.form(x)
		in := scope
		if f.kind == kindMethod {
			in = methods
		}
		d, declared := in[names[x]]
		if !declared {
			continue
		}
		if !f.fits(b, d, &m.Enum.ID) {
			m.ErrorfAt(Conflict, d.Position(), "%s has the name of the %s of %s and another form, so %s has no "+
				"surface. The %s is %s", names[x], labels[x], m.Enum.Name, m.Enum.Name, labels[x], f.text)
			fits = false
			continue
		}
		found[x] = d
	}
	return found, fits
}

// surface returns the members that the record of the enumeration of m
// declares. It leaves out each member that is already declared, and every
// member under methods=off. It writes a parse function only beside a render
// that it writes, and the error of the parse function only beside a parse
// function that it writes. A language without a pair that writes a value as
// text has no encode and no decode. A decode needs a parse function, and an
// encode needs a decode.
func surface(m *sdk.EnumMatch, n *Names, declared *Members) Members {
	if _, off := m.Directive().Params[MethodsParam]; off {
		return Members{}
	}
	out := Members{Render: !declared.Render, Values: !declared.Values, Valid: !declared.Valid}
	out.Parse = out.Render && !declared.Parse
	out.Error = out.Parse && !declared.Error
	if n.Encode == "" {
		return out
	}
	out.Decode = !declared.Decode && (out.Parse || declared.Parse)
	out.Encode = !declared.Encode && (out.Decode || declared.Decode)
	return out
}
