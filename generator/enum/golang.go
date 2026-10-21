// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"embed"
	"math/big"
	"slices"
	"strconv"
	"strings"

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

// The names of the members of the surface of Go that do not depend on the
// enumeration: the method String of fmt.Stringer, the method IsValid, and the
// two methods of encoding.TextMarshaler and encoding.TextUnmarshaler. The
// name of the list of the variants ends with golangValues.
const (
	golangRender = "String"
	golangValid  = "IsValid"
	golangEncode = "MarshalText"
	golangDecode = "UnmarshalText"
	golangValues = "Values"
)

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

// Options returns the template tree of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{sdk.TemplatesOf(golangTemplates, golang.Target)}
}

// names returns the names of the members of the surface of the enumeration
// named subject, such as ParseStatus, StatusValues and ErrUnknownStatus for
// Status and the words parse and unknown. The rules of Go name the error.
func (golangPart) names(b *rules.Bound, subject, parse, sentinel string) Names {
	unknown, _ := b.SentinelName(naming.Pascal(sentinel) + subject)
	return Names{
		Render: golangRender, Parse: naming.Pascal(parse) + subject, Values: subject + golangValues,
		Valid: golangValid, Encode: golangEncode, Decode: golangDecode, Sentinel: unknown,
	}
}

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
