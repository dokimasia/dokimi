// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub

import (
	"embed"
	"strings"
	"text/template"

	"go.dokimi.dev/eidos/lang/naming"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// The member that configures a method joins typescriptOn and the name of the
// method. The one result of a method is in the field typescriptResult of its
// call, and the error that fails a call in the field typescriptFailure. The
// template of the tests imports the doubles from the module that moduleHelper
// returns: the name of its file without testSuffix, in its own directory.
// The templates spell a type of a module through typeHelper, which imports
// the type for the type checker alone, and write a name in Pascal case
// through pascalHelper. They build the type of a value that may be absent
// through optionalHelper, of a list through listHelper and of the result of
// a function of an asynchronous method through awaitableHelper. Such a
// function returns its result or a value of promiseLike.
const (
	typescriptOn      = "on"
	typescriptResult  = "result"
	typescriptFailure = "error"
	moduleHelper      = "stubModule"
	testSuffix        = ".test.ts"
	typeHelper        = "stubType"
	pascalHelper      = "stubPascal"
	optionalHelper    = "stubOptional"
	listHelper        = "stubList"
	awaitableHelper   = "stubAwaitable"
	promiseLike       = "PromiseLike"
)

// typescriptFault opens the name of the helper of the stub of a method that
// arms an error value for the next call, as in faultNotFound for notFound.
const typescriptFault = "fault"

// typescriptTemplates contains the file templates of the target typescript.
//
//go:embed templates/typescript
var typescriptTemplates embed.FS

// typescriptCaser cases the names that the part of TypeScript derives. It
// recognises no initialism, so each word keeps the case that its author
// wrote: the member that configures findById is onFindById.
var typescriptCaser = naming.New()

// typescriptPart is the part of stub for TypeScript.
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
			moduleHelper: module, typeHelper: typeOf, camelHelper: typescriptCaser.Camel,
			pascalHelper: typescriptCaser.Pascal, siteHelper: siteOf, optionalHelper: optional,
			listHelper: listOf, awaitableHelper: awaitable, refHelper: referentOf,
		}),
	}
}

// module returns the specifier of the module of the doubles that the test
// file named file tests: the name of the file without its suffix .test.ts,
// in the directory of the file. The specifier is wrong for a double that a
// config routes into another directory.
func module(file string) string { return "./" + strings.TrimSuffix(file, testSuffix) }

// typeOf returns the reference of the type name of the module specifier,
// such as the type Seat of @dokimi/assert. The reference has no target, so
// the speller imports it from its module as a type. An empty specifier names
// a global type, such as void, which the speller imports from no module.
func typeOf(name, specifier string) *emit.TypeRef {
	return &emit.TypeRef{Spelling: name, Package: specifier}
}

// optional returns the union of t and undefined, the type of a value of t
// that may be absent. The speller writes the union from its members, and
// puts a function type in parentheses.
func optional(t *emit.TypeRef) *emit.TypeRef {
	return &emit.TypeRef{Form: symbol.FormUnion, Elems: []*emit.TypeRef{t, {Spelling: typescript.TypeUndefined}}}
}

// listOf returns the list of t, the type of the arguments of a variadic
// parameter of t. The speller puts a union or a function type of t in
// parentheses.
func listOf(t *emit.TypeRef) *emit.TypeRef {
	return &emit.TypeRef{Form: symbol.FormList, Elems: []*emit.TypeRef{t}}
}

// awaitable returns the union of t and PromiseLike of t, the type of the
// result of a function that a test sets for an asynchronous method whose
// result has the type t.
func awaitable(t *emit.TypeRef) *emit.TypeRef {
	return &emit.TypeRef{
		Form: symbol.FormUnion, Elems: []*emit.TypeRef{t, {Spelling: promiseLike, Args: []*emit.TypeRef{t}}},
	}
}

// structural reports true, because TypeScript types a value by its members
// alone. A double of a class is a value of the class, and a double declares
// the properties of its subject.
func (typescriptPart) structural() bool { return true }

// referable reports whether the subject is exported, because the file of a
// double imports its subject from the module of the subject.
func (typescriptPart) referable(v symbol.Visibility) bool { return v == symbol.VisibilityPublic }

// receiver returns m unchanged, because a TypeScript method reads its object
// through this and declares no receiver.
func (typescriptPart) receiver(m *emit.Method, _ []*emit.TypeParam, _ string) *emit.Method { return m }

// constraint reports false, because every TypeScript interface is the type of
// its values.
func (typescriptPart) constraint(sdk.Matcher, *symbol.Identity) bool { return false }

// names returns the names of a TypeScript method of a double. The member that
// configures the method is on followed by the name of the method in Pascal
// case, such as onGet for get and onFindById for findById. A field of the call
// has the name of its parameter. A TypeScript method returns one value at
// most, so the call records its result in the field result, and the setter
// Returns takes it as result. The method declares no names of its result. A
// TypeScript method fails by a throw, so the call records the error of a
// fault or of a function in the field error.
func (typescriptPart) names(method string, params []string, _ []rules.ReturnView, _ map[string]bool) Names {
	return Names{
		Field:   typescriptOn + typescriptCaser.Pascal(method),
		Failure: typescriptFailure,
		Args:    params,
		Results: []string{typescriptResult},
		Setters: []string{typescriptResult},
	}
}

// delegate returns value, because a class types its values by its members
// alone.
func (typescriptPart) delegate(value *emit.TypeRef, _ bool) *emit.TypeRef { return value }

// declarations returns nil, because a check of TypeScript imports the double
// and the types of its calls from the module of the double.
func (typescriptPart) declarations(*emit.Struct, *Model) []symbol.Symbol { return nil }

// fault returns fault followed by base in Pascal case, such as faultNotFound
// for the error value notFound.
func (typescriptPart) fault(base string) string {
	return typescriptFault + typescriptCaser.Pascal(base)
}

// construction returns the builder class named name and its static method
// seed, such as AccountBuilder and create. A TypeScript double calls the
// method of the class.
func (typescriptPart) construction(name, seed string) (ref, method string) { return name, seed }
