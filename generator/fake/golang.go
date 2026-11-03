// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"embed"
	"slices"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// golangComparable is the constraint of Go that permits only comparable
// types.
const golangComparable = "comparable"

// golangTemplates contains the file templates of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// golangPart is the part of fake for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// Options returns the template tree of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{sdk.TemplatesOf(golangTemplates, golang.Target)}
}

// receiver sets a pointer receiver named name on the method m of a fake, as
// Go declares a method that changes its value, and returns m. The receiver
// restates the type parameters tps of a generic fake.
func (golangPart) receiver(m *emit.Method, tps []*emit.TypeParam, name string) *emit.Method {
	m.Receives.Args = sdk.TypeArgs(tps)
	golang.PointerReceiver(m).Receiver.Name = name
	return m
}

// field returns the expression of the field named field of the parameter
// value, such as v.Key.
func (golangPart) field(value, field string) string { return value + "." + field }

// keyable returns the reason why a map of Go cannot key its values by the
// parameter key, and the empty string where it can. A slice and a map are not
// comparable, and neither is a type parameter of tps whose constraint is not
// comparable. The reference forms refuse a key of a function type, and the
// compiler checks the comparability of every other type.
func (golangPart) keyable(key *node.Param, tps []*node.TypeParam) string {
	at := slices.IndexFunc(tps, func(tp *node.TypeParam) bool { return tp.ID == key.Type.Target })
	switch form := key.Type.Form; {
	case at >= 0 && !slices.ContainsFunc(tps[at].Bounds, func(b *node.TypeRef) bool {
		return b.Spelling == golangComparable
	}):
		return "whose constraint permits types that a Go map cannot key"
	case form == symbol.FormList || form == symbol.FormMap:
		return "which a Go map cannot key"
	}
	return ""
}
