// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"embed"
	"fmt"
	"go/token"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"text/template"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// golangConstValue is the handle of the exact value of a constant, which the
// frontend of Go stamps where the package of the constant suffices to
// evaluate it.
var golangConstValue = meta.Named[string](golang.ConstValueKey)

// golangRuns is the largest number of runs of consecutive values that the
// render of Go switches over. The render of an enumeration with more runs
// looks the text of a value up in a map, as the tool stringer writes it.
const golangRuns = 10

// The widths of the unsigned integers of the index of a run of a table: 8
// bits for texts below 256 bytes, 16 bits for texts below 65,536 bytes, and
// 32 bits for longer texts.
const (
	golangByteBits  = 8
	golangShortBits = 16
	golangWideBits  = 32
)

// golangValue is one variant of an enumeration of integers: its value, its
// text, and its entry in the table.
type golangValue struct {
	n     *big.Int
	text  string
	entry Entry
}

// The names of the members of the surface of Go that the language fixes: the
// method String of fmt.Stringer, and the two methods of
// encoding.TextMarshaler and encoding.TextUnmarshaler.
const (
	golangRender = "String"
	golangEncode = "MarshalText"
	golangDecode = "UnmarshalText"
)

// golangFieldHelper is the helper that the template of the checks calls with
// the record of the checks and the name of a field.
const golangFieldHelper = "fieldType"

// golangForms are the forms of the members of the surface of Go, which a
// declaration under the name of a member has.
var golangForms = map[member]form{
	memberRender: {
		kind: kindMethod, values: []shape{shapeText},
		text: "a method of the enumeration without a parameter that returns a text",
	},
	memberParse: {
		kind: kindFunction, inputs: []shape{shapeText}, values: []shape{shapeVariant}, fails: true,
		text: "a function of a text that returns a variant and an error",
	},
	memberError: {kind: kindVariable, text: "a variable"},
	memberValues: {
		kind: kindFunction, values: []shape{shapeVariants},
		text: "a function without a parameter that returns a list of the enumeration",
	},
	memberValid: {
		kind: kindMethod, values: []shape{shapeBool},
		text: "a method of the enumeration without a parameter that returns a truth value",
	},
	memberEncode: {
		kind: kindMethod, values: []shape{shapeBytes}, fails: true,
		text: "a method of the enumeration without a parameter that returns bytes and an error",
	},
	memberDecode: {
		kind: kindMethod, inputs: []shape{shapeBytes}, fails: true,
		text: "a method of the enumeration of bytes that returns an error",
	},
}

// The verbs of fmt that print a value outside the variants after its
// conversion to the underlying type: a float, a truth value, and any other
// number.
const (
	golangFloatVerb  = "%g"
	golangBoolVerb   = "%t"
	golangNumberVerb = "%d"
)

// golangTemplates contains the file templates of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// golangPart is the part of enum for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// Options returns the template tree and the helper of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(golangTemplates, golang.Target),
		sdk.Funcs(template.FuncMap{golangFieldHelper: fieldType}),
	}
}

// fieldType returns the type of the field named name of the record d of the
// checks, and an error where d has no such field. A check spells the type to
// call the parse function or the list of the variants, or to match the error.
// The layout of eidos respells the name of a field in the case of its
// language, such as Values for values, so fieldType compares names without
// their case.
func fieldType(d *emit.Struct, name string) (*emit.TypeRef, error) {
	for _, f := range d.Fields.Items() {
		if strings.EqualFold(f.Name, name) {
			return f.Type, nil
		}
	}
	return nil, fmt.Errorf("enum: the record of the checks of %s has no field %q", d.Name, name)
}

// templated returns the name of an enumeration in PascalCase, as State for
// state, so a name template renders ParseState for an enumeration that its
// package does not export. names writes each rendered name in PascalCase
// again, so an exported name renders as it did before.
func (golangPart) templated(subject string) string { return naming.Pascal(subject) }

// names returns the names of the members of the surface of an enumeration in
// PascalCase, such as ParseStatus, StatusValues, IsValid and ErrUnknownStatus
// for Status and the default templates. The rules of Go name the error with
// the prefix Err. The functions and the error of an enumeration that its
// package does not export are not exported either, as parseState,
// stateValues and errUnknownState for state. Its methods keep their names.
func (golangPart) names(b *rules.Bound, subject string, n baseNames) Names {
	unknown, _ := b.SentinelName(n.err)
	out := Names{
		Render: golangRender, Parse: naming.Pascal(n.parse), Values: naming.Pascal(n.values),
		Valid: naming.Pascal(n.valid), Encode: golangEncode, Decode: golangDecode, Sentinel: unknown,
	}
	if !token.IsExported(subject) {
		out.Parse, out.Values, out.Sentinel = naming.Camel(out.Parse), naming.Camel(out.Values),
			naming.Camel(out.Sentinel)
	}
	return out
}

// form returns the form of a member of the surface of Go, which has every
// member.
func (golangPart) form(x member) form { return golangForms[x] }

// fallback returns the type that the enumeration e of Go declares under it,
// and the verb that prints a value of that type. It returns nil for an
// enumeration of text, whose render returns a value outside the variants as
// its text. A Go enumeration is a defined type, so it always declares the
// type under it.
func (golangPart) fallback(b *rules.Bound, e *node.Enum, textual bool) (*emit.TypeRef, string) {
	if textual {
		return nil, ""
	}
	shape := b.TypeOf(e.Underlying)
	switch {
	case shape.Form == symbol.FormBool:
		return rules.EmitRef(e.Underlying), golangBoolVerb
	case shape.Class == rules.ScalarFloat:
		return rules.EmitRef(e.Underlying), golangFloatVerb
	default:
		return rules.EmitRef(e.Underlying), golangNumberVerb
	}
}

// member returns name, because a variant of Go is a constant of the package
// of its enumeration.
func (golangPart) member(name string) string { return name }

// literal returns the literal of Go of the value v: a string in quotes, and
// a number or a truth value as its text.
func (golangPart) literal(v *emit.Value) string {
	if v.Literal == emit.LiteralString {
		return strconv.Quote(v.Text)
	}
	return v.Text
}

// qualified reports true, because a check of Go names a function of another
// package through the name of that package, as in catalog.ParseStatus.
func (golangPart) qualified() bool { return true }

// table returns the table of the texts of an enumeration of integers, from
// which the render of Go slices the text of a value in the form of the tool
// stringer. The value of each variant is the exact value that the frontend
// stamps. table returns nil for an enumeration of another kind, and where a
// variant has no exact value, such as a variant whose value reads a
// constant of another package.
func (golangPart) table(
	m sdk.Matcher, b *rules.Bound, e *node.Enum, variants []rules.VariantText, ids map[string]symbol.Identity,
) *Table {
	shape := b.TypeOf(e.Underlying)
	if shape.Form != symbol.FormScalar || shape.Class == rules.ScalarFloat {
		return nil
	}
	t := &Table{}
	values := make([]golangValue, 0, len(variants))
	for i := range variants {
		v := &variants[i]
		exact, _ := sdk.FactOf(m, ids[v.Name], golangConstValue)
		n, read := new(big.Int).SetString(exact, 10)
		if !read {
			return nil
		}
		entry := Entry{Name: v.Name, Value: exact}
		t.Values = append(t.Values, entry)
		values = append(values, golangValue{n: n, text: v.Text.Text, entry: entry})
	}
	slices.SortFunc(values, func(x, y golangValue) int { return x.n.Cmp(y.n) })
	one := big.NewInt(1)
	for rest := values; len(rest) > 0; {
		n := 1
		for n < len(rest) && new(big.Int).Sub(rest[n].n, rest[n-1].n).Cmp(one) == 0 {
			n++
		}
		run := golangRun(rest[:n], shape.Class == rules.ScalarUint)
		t.Runs = append(t.Runs, run)
		t.Indexed = t.Indexed || !run.Single
		rest = rest[n:]
	}
	if len(t.Runs) <= golangRuns {
		return t
	}
	t.Map = true
	var name strings.Builder
	for _, v := range values {
		entry := v.entry
		entry.Start = name.Len()
		name.WriteString(v.text)
		entry.End = name.Len()
		t.Entries = append(t.Entries, entry)
	}
	t.Name = strconv.Quote(name.String())
	return t
}

// golangRun returns the run of the consecutive values in increasing order:
// their texts in one string literal, the offset of each text, and the first
// and the last value. unsigned reports that the type of the values is
// unsigned.
func golangRun(values []golangValue, unsigned bool) Run {
	var name strings.Builder
	index := make([]int, 1, len(values)+1)
	for _, v := range values {
		name.WriteString(v.text)
		index = append(index, name.Len())
	}
	return Run{
		Name: strconv.Quote(name.String()), Low: values[0].entry.Value, High: values[len(values)-1].entry.Value,
		Index: index, Bits: golangWidth(name.Len()), Open: unsigned && values[0].n.Sign() == 0,
		Single: len(values) == 1,
	}
}

// golangWidth returns the width of the unsigned integers of an index of
// offsets into texts of size bytes, as the tool stringer sizes it.
func golangWidth(size int) int {
	switch {
	case size < 1<<golangByteBits:
		return golangByteBits
	case size < 1<<golangShortBits:
		return golangShortBits
	default:
		return golangWideBits
	}
}
