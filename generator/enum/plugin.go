// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"cmp"
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

// ParseWord and SentinelWord are the defaults of the options parse-word and
// sentinel-word: the words that the names of the parse function and of its
// error join to the name of the enumeration, as in ParseStatus and
// ErrUnknownStatus.
const (
	ParseWord    = "parse"
	SentinelWord = "unknown"
)

// The tree of each target contains the file template of the surfaces and
// the file template of their checks.
const (
	enumTemplate = "enum.tmpl"
	testTemplate = "enum.test.tmpl"
)

// unknownMark is the character that a text of no variant repeats.
const unknownMark = "?"

// Options are the options of enum. ParseWord is the word of the name of each
// parse function, and SentinelWord is the word of the name of the error that
// a parse function returns for a text that no variant has. An empty word
// takes its default, [ParseWord] or [SentinelWord]. NoTests leaves out the
// checks of every enumeration.
type Options struct {
	ParseWord    string `opt:"parse-word" doc:"the word that the name of each parse function joins to the name of its enumeration"`
	SentinelWord string `opt:"sentinel-word" doc:"the word that the name of the error of a parse function joins to the name of its enumeration"`
	NoTests      bool   `opt:"no-tests" doc:"leaves out the checks of every enumeration"`
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
// language reaches a variant through it, as in Status.Draft of TypeScript.
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
// check calls.
type Members struct {
	Render bool
	Parse  bool
	Values bool
	Valid  bool
	Encode bool
	Decode bool
}

// Model is the data of the record of the surface of one enumeration, and the
// file template of the surfaces reads it. Subject is the name of the
// enumeration and Package the name of its package, which opens the message of
// Sentinel. Type is the type of the enumeration. Fallback is the type that a
// value outside the variants converts to before the render prints it with
// the verb Format, and it is nil for an enumeration of text, which renders
// such a value as itself. Generate are the members that the record declares,
// which are the members that the enumeration and its package do not declare.
//
// Table is the table of the texts of an enumeration of integers, which the
// render slices. It is nil where the part of the language writes no table,
// such as for an enumeration of another kind. Converts reports that the text
// of each variant is its value, so the render converts a value to text.
//
// The record of a surface with a parse function names the struct of the
// enumeration. In a language whose checks refer to the surface through its
// package, the record also declares the parse function, the list of the
// variants and the error of the parse function beside that struct, and the
// template writes them at the struct.
type Model struct {
	Type     *emit.TypeRef
	Fallback *emit.TypeRef
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
// member that the record of the surface declares, or that the enumeration
// declares itself. Generate are the members that the record of the surface
// declares. Stated reports that a variant states its text, and the checks
// pin the text of each such variant.
//
// Where the record of the surface declares the parse function, the list of
// the variants and the error, the record of the checks has a field of each,
// in that order, whose type refers to the declaration. The layout of eidos
// qualifies that reference where the checks are in another package than the
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
	// names returns the names of the members of the surface of the
	// enumeration named subject. parse and sentinel are the words of the
	// parse function and of its error. b are the rules of the language of
	// the enumeration.
	names(b *rules.Bound, subject, parse, sentinel string) Names
	// fallback returns the type that a value outside the variants of e
	// converts to before the render prints it, and the verb that prints it.
	// It returns nil for an enumeration of text, and for a language whose
	// render prints such a value without a conversion.
	fallback(b *rules.Bound, e *node.Enum, textual bool) (*emit.TypeRef, string)
	// member returns the reference to the variant named name, which the
	// template writes after the enumeration where the language reaches a
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
// value as text and decodes it. enum records the [Checks] of each
// enumeration into the family [Word] under the tag [TestTag]. The checks
// test the enumeration through the exported API of its package, so a Go file
// of the family is in the external test package of its directory.
//
// The surface leaves out each member that the enumeration or its package
// already declares under its name. A parse function and a list of the
// variants go with the render, because they are beside the type: an
// enumeration that renders itself parses itself. A decode needs a parse
// function, so without one the surface leaves out the encode too. The value
// off of [MethodsParam] leaves out the whole surface and keeps the checks.
//
// The text of a variant of an enumeration of text is its value, and the text
// of any other variant is its name, as the rules of its language derive them.
// Under the value TextComment of [TextParam], the trailing comment of a
// variant is its text, as in Alpha Grade = iota // alpha. The checks pin each
// text that a value or a comment states. enum serves Go and TypeScript, the
// languages of its parts. It reports
// [Empty] at an enumeration without a variant, and [Duplicate] at an
// enumeration whose two variants have one value or one text. enum records
// nothing for an enumeration that it reports. Each call returns a new plugin.
func New() plugin.Generator {
	g := &generator{options: &Options{ParseWord: ParseWord, SentinelWord: SentinelWord}}
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
	names := p.names(&b, enum.Name, cmp.Or(g.options.ParseWord, ParseWord),
		cmp.Or(g.options.SentinelWord, SentinelWord))
	pkg, _ := m.Reader().PackageOf(enum.ID.PackageIdentity())
	declared := declares(m, pkg, &names)
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
	// The record declares the functions and the variable of the surface that
	// the checks refer to through its package. They go with the parse
	// function.
	referred := p.qualified() && generated.Parse
	if generated != (Members{}) {
		record := []symbol.Symbol{&emit.Struct{Origin: enum.ID, Name: enum.Name, Visibility: enum.Visibility}}
		if referred {
			record = append(record,
				&emit.Function{Origin: enum.ID, Name: names.Parse, Visibility: enum.Visibility},
				&emit.Function{Origin: enum.ID, Name: names.Values, Visibility: enum.Visibility},
				&emit.Variable{Origin: enum.ID, Name: names.Sentinel, Visibility: enum.Visibility},
			)
		}
		e.File().Record(model, record...)
	}
	if g.options.NoTests {
		return nil
	}
	unknown := emit.Literal(emit.LiteralString, strings.Repeat(unknownMark, longest+1))
	checks := Checks{
		Type: typ, Subject: enum.Name, Names: names, Variants: variants, Generate: generated,
		Stated: slices.Contains(stated, true), Unknown: p.literal(&unknown),
		Calls: Members{
			Render: generated.Render || declared.Render, Parse: generated.Parse,
			Values: generated.Values, Valid: generated.Valid || declared.Valid,
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
	if referred {
		for _, name := range []string{names.Parse, names.Values, names.Sentinel} {
			test.Fields.Append(&emit.Field{Name: name, Type: &emit.TypeRef{Spelling: name}})
		}
	}
	e.File(TestTag).Record(checks, test)
	return nil
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

// declares reports the members of the surface that the enumeration of m or
// its package pkg already declares: a method of the enumeration, or a
// function of the package, under the name of the member. No method and no
// function has the empty name of a member that the language does not have.
// The list of the variants goes with the render, so declares does not look
// for it.
func declares(m *sdk.EnumMatch, pkg *node.Package, n *Names) Members {
	functions := map[string]bool{}
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			if fn, isFunction := d.(*node.Function); isFunction {
				functions[fn.Name] = true
			}
		}
	}
	declared := func(name string) bool {
		return functions[name] ||
			slices.ContainsFunc(m.Enum.Methods, func(x *node.Method) bool { return x.Name == name })
	}
	return Members{
		Render: declared(n.Render), Parse: declared(n.Parse), Valid: declared(n.Valid),
		Encode: declared(n.Encode), Decode: declared(n.Decode),
	}
}

// surface returns the members that the record of the enumeration of m
// declares. It leaves out each member that is already declared, and every
// member under methods=off. The parse function and the list of the variants
// go with the render. A language without a pair that writes a value as text
// has no encode and no decode. A decode needs a parse function, so without
// one the surface leaves out the encode too.
func surface(m *sdk.EnumMatch, n *Names, declared *Members) Members {
	if _, off := m.Directive().Params[MethodsParam]; off {
		return Members{}
	}
	out := Members{Render: !declared.Render, Valid: !declared.Valid}
	out.Parse, out.Values = out.Render, out.Render
	if n.Encode == "" {
		return out
	}
	out.Encode = !declared.Encode
	switch {
	case declared.Decode:
	case out.Parse || declared.Parse:
		out.Decode = true
	default:
		out.Encode = false
	}
	return out
}
