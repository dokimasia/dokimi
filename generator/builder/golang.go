// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"embed"
	"go/token"
	"strings"

	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// golangEmptyStruct is the spelling of the empty struct, whose values make a
// map a set.
const golangEmptyStruct = "struct{}"

// The name of each constructor of a builder opens with golangNew and the
// name of the builder, such as NewAddressBuilder. The name of each
// constructor of an unexported builder opens with golangNewUnexported and the
// name of the builder in Pascal case, such as newConfigBuilder. The name of
// the constructor from a value ends with golangFrom.
const (
	golangNew           = "New"
	golangNewUnexported = "new"
	golangFrom          = "From"
)

// The setter of a member opens with golangWith, and the appender of a list
// opens with golangAppend. The names of the setters of the entries of a map
// or of the keys of a set end with golangEntry and golangEntries. The name
// of the setter of bytes from text ends with golangText.
const (
	golangWith    = "With"
	golangAppend  = "Append"
	golangEntry   = "Entry"
	golangEntries = "Entries"
	golangText    = "String"
)

// golangTemplates contains the file templates of the target golang.
//
//go:embed templates/golang
var golangTemplates embed.FS

// golangPart is the part of builder for Go.
type golangPart struct{}

// Lang returns the language Go.
func (golangPart) Lang() symbol.Lang { return golang.Lang }

// Target returns the target golang.
func (golangPart) Target() plugin.Target { return golang.Target }

// Options returns the template tree of the target golang.
func (golangPart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{sdk.TemplatesOf(golangTemplates, golang.Target)}
}

// builder returns the name of the builder of the declaration named name, such
// as AddressBuilder for Address and the word builder. The word joins the name
// in Pascal case with the initialisms of Go.
func (golangPart) builder(name, word string) string { return name + naming.Pascal(word) }

// companion returns the name of the companion function of the declaration
// named name, such as AddressDefaults for Address and the word defaults. The
// word joins the name in Pascal case with the initialisms of Go.
func (golangPart) companion(name, word string) string { return name + naming.Pascal(word) }

// constructors returns NewAddressBuilder and NewAddressBuilderFrom for the
// builder AddressBuilder. A constructor of Go opens with New and the type
// that it returns, so it never has the name of a constructor of the
// declaration itself, such as NewAddress. The constructors of an unexported
// builder are unexported, such as newConfigBuilder and newConfigBuilderFrom
// for configBuilder.
func (golangPart) constructors(builder string) (fromSeed, fromValue string) {
	fromSeed = golangNew + builder
	if !token.IsExported(builder) {
		fromSeed = golangNewUnexported + naming.Pascal(builder)
	}
	return fromSeed, fromSeed + golangFrom
}

// functions reports true, because a constructor of Go is a function of the
// package of its type.
func (golangPart) functions() bool { return true }

// names returns the names of the methods of the member named field, such as
// WithLines and AppendLines for the list Lines. The form of the member
// decides which methods it has. The name of a member of Go is exported, so
// each name joins it as it is.
func (golangPart) names(field string, form Form) Names {
	n := Names{With: golangWith + field}
	switch form {
	case FormList:
		n.Append = golangAppend + field
	case FormMap, FormSet:
		n.Entry, n.Entries = n.With+golangEntry, n.With+golangEntries
	case FormBytes:
		n.Text = n.With + golangText
	default:
		// A scalar and an optional have a setter alone.
	}
	return n
}

// selectable reports true, because the name of a field of Go is an
// identifier, which a selector writes after a dot.
func (golangPart) selectable(string) bool { return true }

// referable reports true, because the file of a builder is in the package of
// its struct.
func (golangPart) referable(symbol.Visibility) bool { return true }

// set returns the key of a map whose values are the empty struct, the set
// of Go, and false for every other type. A map of Go has its key and its
// value as its elements, and the spelling of the empty struct may contain
// spaces.
func (golangPart) set(ref *node.TypeRef) (*node.TypeRef, bool) {
	if ref.Form != symbol.FormMap || strings.Join(strings.Fields(ref.Elems[1].Spelling), "") != golangEmptyStruct {
		return nil, false
	}
	return ref.Elems[0], true
}
