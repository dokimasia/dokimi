// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"embed"
	"slices"
	"strings"
	"text/template"
	"unicode/utf8"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// typescriptCreate and typescriptFrom are the names of the static methods of
// a builder class that construct a builder from the seed and from a value.
const (
	typescriptCreate = "create"
	typescriptFrom   = "from"
)

// The setter of a member opens with typescriptWith, and the appender of a
// list opens with typescriptAppend. The names of the setters of the entries
// of a map or of the keys of a set end with typescriptEntry and
// typescriptEntries.
const (
	typescriptWith    = "with"
	typescriptAppend  = "append"
	typescriptEntry   = "Entry"
	typescriptEntries = "Entries"
)

// The helpers that the templates of the target typescript call: whether a
// map is an instance of the class Map rather than an object, the removal of a
// suffix from the name of a file, and the name of a function of cases.
const (
	mapClassHelper   = "mapclass"
	trimSuffixHelper = "trimSuffix"
	casesHelper      = "builderCases"
)

// typescriptCases ends the name of a function that declares the cases of one
// function of a builder.
const typescriptCases = "Cases"

// typescriptTemplates contains the file templates of the target typescript.
//
//go:embed templates/typescript
var typescriptTemplates embed.FS

// typescriptCaser cases the names that the part of TypeScript derives. It
// recognises no initialism, so each word keeps the case that its author
// wrote. The setter of userId is withUserId, and the companion function of
// HTTPClient is httpClientDefaults.
var typescriptCaser = naming.New()

// typescriptPart is the part of builder for TypeScript.
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
		sdk.Funcs(template.FuncMap{
			mapClassHelper: mapClass, trimSuffixHelper: strings.TrimSuffix, casesHelper: casesOf,
		}),
	}
}

// casesOf returns the name of the function that declares the cases of the
// function named function of the builder class named builder, such as
// addressBuilderWithStreetCases for withStreet of AddressBuilder.
func casesOf(builder, function string) string {
	return typescriptCaser.Camel(builder) + typescriptCaser.Pascal(function) + typescriptCases
}

// builder returns the name of the builder class of the class named name, such
// as AddressBuilder for Address and the word builder. The word joins the name
// in Pascal case.
func (typescriptPart) builder(name, word string) string { return name + typescriptCaser.Pascal(word) }

// companion returns the name of the companion function of the class named
// name in camel case, such as addressDefaults for Address and
// httpClientDefaults for HTTPClient with the word defaults.
func (typescriptPart) companion(name, word string) string {
	return typescriptCaser.Camel(name) + typescriptCaser.Pascal(word)
}

// constructors returns create and from, the static methods of every builder
// class. A caller names the class beside them, as in AddressBuilder.create().
func (typescriptPart) constructors(string) (fromSeed, fromValue string) {
	return typescriptCreate, typescriptFrom
}

// functions reports false, because the constructors of a builder class are
// its static methods.
func (typescriptPart) functions() bool { return false }

// names returns the names of the methods of the member named field in camel
// case, such as withLines and appendLines for the list lines. The form of
// the member decides which methods it has. TypeScript has no setter of bytes
// from text, because a caller encodes text with TextEncoder.
func (typescriptPart) names(field string, form Form) Names {
	pascal := typescriptCaser.Pascal(field)
	n := Names{With: typescriptWith + pascal}
	switch form {
	case FormList:
		n.Append = typescriptAppend + pascal
	case FormMap, FormSet:
		n.Entry, n.Entries = n.With+typescriptEntry, n.With+typescriptEntries
	default:
		// A scalar, bytes and an optional have a setter alone.
	}
	return n
}

// selectable reports whether name is an identifier name, which a property
// access of TypeScript writes after a dot. A reserved word is an identifier
// name, as in value.default. A quoted name such as content-type is not.
// utf8.DecodeRuneInString returns utf8.RuneError for the empty name, and
// RuneError starts no identifier.
func (typescriptPart) selectable(name string) bool {
	first, size := utf8.DecodeRuneInString(name)
	return typescript.IsIdentifierStart(first) &&
		!strings.ContainsFunc(name[size:], func(r rune) bool { return !typescript.IsIdentifierPart(r) })
}

// referable reports whether the class is exported, because the module of a
// builder imports its class from the module of the class.
func (typescriptPart) referable(v symbol.Visibility) bool { return v == symbol.VisibilityPublic }

// set returns the type argument of a reference to the global class Set, the
// set of TypeScript, and false for every other type. A property without a
// type annotation has a nil type. TypeScript requires the one type argument
// of Set in a type annotation.
func (typescriptPart) set(ref *node.TypeRef) (*node.TypeRef, bool) {
	if ref == nil || !global(&ref.Target, ref.Package, ref.Spelling, typescript.TypeSet) {
		return nil, false
	}
	return ref.Args[0], true
}

// mapClass reports whether the map ref is an instance of the global class
// Map or ReadonlyMap, whose entries a template sets through the class,
// rather than an object whose properties are the entries.
func mapClass(ref *emit.TypeRef) bool {
	return global(&ref.Target, ref.Package, ref.Spelling, typescript.TypeMap, typescript.TypeReadonlyMap)
}

// global reports whether a reference with the target, the package and the
// spelling refers to one of the global types of TypeScript in names. The
// library of TypeScript declares a global type, or the load leaves a
// reference to it without a target. The workspace and an import do not
// declare a global type.
func global(target *symbol.Identity, pkg, spelling string, names ...string) bool {
	return typescript.Global(*target, pkg) && slices.Contains(names, spelling)
}
