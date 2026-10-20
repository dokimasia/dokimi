// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"embed"
	"strconv"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

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
