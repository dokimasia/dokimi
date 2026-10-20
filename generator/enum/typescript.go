// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import (
	"bytes"
	"embed"
	"encoding/json"
	"strings"
	"text/template"
	"unicode/utf8"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// The parts of the names of the members of the surface of TypeScript: the
// render ends with typescriptRender, the list of the variants ends with
// typescriptValues, the check of a value opens with typescriptValid, and the
// error of the parse function ends with typescriptError.
const (
	typescriptRender = "ToString"
	typescriptValues = "Values"
	typescriptValid  = "is"
	typescriptError  = "Error"
)

// The helpers that the template of the checks calls with the name of its
// file: the module of the surface beside it, and the module of the
// enumeration. The name of the file of the checks ends with
// typescriptChecksSuffix, and the name of the module of the surface with
// typescriptSurfaceSuffix.
const (
	surfaceModuleHelper     = "surfaceModule"
	sourceModuleHelper      = "sourceModule"
	typescriptChecksSuffix  = ".test.ts"
	typescriptSurfaceSuffix = ".enum"
	typescriptRelative      = "./"
)

// typescriptTemplates contains the file templates of the target typescript.
//
//go:embed templates/typescript
var typescriptTemplates embed.FS

// typescriptCaser cases the names that the part of TypeScript derives. It
// recognises no initialism, so each word keeps the case that its author
// wrote.
var typescriptCaser = naming.New()

// typescriptPart is the part of enum for TypeScript.
type typescriptPart struct{}

// Lang returns the language TypeScript.
func (typescriptPart) Lang() symbol.Lang { return typescript.Lang }

// Target returns the target typescript.
func (typescriptPart) Target() plugin.Target { return typescript.Target }

// Options returns the template tree and the helpers of the target
// typescript.
func (typescriptPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{
		sdk.TemplatesOf(typescriptTemplates, typescript.Target),
		sdk.Funcs(template.FuncMap{surfaceModuleHelper: surfaceModule, sourceModuleHelper: sourceModule}),
	}
}

// surfaceModule returns the module of the surface beside the file of the
// checks named file, such as ./status.enum for status.enum.test.ts.
func surfaceModule(file string) string {
	return typescriptRelative + strings.TrimSuffix(file, typescriptChecksSuffix)
}

// sourceModule returns the module of the enumeration beside the file of the
// checks named file, such as ./status for status.enum.test.ts.
func sourceModule(file string) string {
	return strings.TrimSuffix(surfaceModule(file), typescriptSurfaceSuffix)
}

// names returns the names of the functions of the surface of the enum named
// subject, such as statusToString, parseStatus, statusValues, isStatus and
// UnknownStatusError for Status and the words parse and unknown. TypeScript
// has no pair that writes a value as text and reads it back, because JSON
// writes the value of a member.
func (typescriptPart) names(_ *rules.Bound, subject, parse, sentinel string) Names {
	camel := typescriptCaser.Camel(subject)
	pascal := typescriptCaser.Pascal(subject)
	return Names{
		Render: camel + typescriptRender, Parse: typescriptCaser.Camel(parse) + pascal,
		Values: camel + typescriptValues, Valid: typescriptValid + pascal,
		Sentinel: typescriptCaser.Pascal(sentinel) + pascal + typescriptError,
	}
}

// fallback returns nil, because the render of TypeScript prints a value
// outside the members through String.
func (typescriptPart) fallback(*rules.Bound, *node.Enum, bool) (*emit.TypeRef, string) {
	return nil, ""
}

// member returns the access of the member named name of an enum: after a dot
// where the name is an identifier, and in brackets otherwise.
func (p typescriptPart) member(name string) string {
	first, size := utf8.DecodeRuneInString(name)
	if typescript.IsIdentifierStart(first) &&
		!strings.ContainsFunc(name[size:], func(r rune) bool { return !typescript.IsIdentifierPart(r) }) {

		return "." + name
	}
	quoted := emit.Literal(emit.LiteralString, name)
	return "[" + p.literal(&quoted) + "]"
}

// qualified reports false, because a check of TypeScript imports each
// function of the surface from the module of the surface.
func (typescriptPart) qualified() bool { return false }

// literal returns the literal of TypeScript of the value v: a string in
// double quotes, which JSON writes as TypeScript reads it, and a number or a
// truth value as its text.
func (typescriptPart) literal(v *emit.Value) string {
	if v.Literal != emit.LiteralString {
		return v.Text
	}
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	// A string always encodes.
	_ = encoder.Encode(v.Text)
	return strings.TrimSuffix(out.String(), "\n")
}
