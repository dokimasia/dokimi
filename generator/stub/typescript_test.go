// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub_test

import (
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	"go.dokimi.dev/eidos/lang/typescript"
	tsbackend "go.dokimi.dev/eidos/lang/typescript/backend"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// finderModule is the module of the interface Finder, whose method has a word
// that the initialisms of Go contain.
const finderModule = "finder"

// The part of TypeScript declares the target typescript with the template of
// each family and its helpers. It returns the names of the member that
// configures a method and of the fields of a call in the case of TypeScript.
// A TypeScript method returns its one value as its result. The part declares
// nothing beside a double, because the checks import from its module.
func TestTypeScript(t *testing.T) {
	t.Parallel()

	templates, provides := stub.New().(plugin.TemplateProvider)
	assert.True(t, provides, "stub declares templates")

	t.Run("declares the template of each family for the target typescript", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, stub.New(), typescript.Target, stubTemplate, testTemplate)
	})

	t.Run("names the member that configures each method after the method with on", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, storeModule).models[storeName]
		fields := make([]string, 0, len(got.Methods))
		for i := range got.Methods {
			fields = append(fields, got.Methods[i].Field)
		}
		assert.Equal(t, fields, []string{"onGet", "onPut", "onSize", "onLog", "onClose"},
			"the member joins on and the name of its method")
	})

	t.Run("keeps the case of each word of the name of a method in the member that configures it", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, finderModule).models[finderName]
		assert.Equal(t, method(t, got, "findById").Field, "onFindById", "the member of findById is onFindById")
	})

	t.Run("writes no receiver on a method of a double", func(t *testing.T) {
		t.Parallel()

		get := declared(t, generate(t, typescript.Target, tsTree, storeModule).doubles[storeName], "get")
		expect.Nil(t, get.Receiver, "a TypeScript method reads its object through this")
	})

	t.Run("declares nothing beside a double", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, storeModule)
		expect.Empty(t, got.beside[storeName], "the checks import the double and its calls from its module")
		expect.Empty(t, gentest.Referents(got.tests[storeName]), "the record of the checks refers to no declaration")
	})

	t.Run("names the fields of a call after its parameters", func(t *testing.T) {
		t.Parallel()

		got := method(t, generate(t, typescript.Target, tsTree, storeModule).models[storeName], "put")
		assert.Equal(t, paramFields(got.Params), "item note", "a field of a call has the name of its parameter")
		expect.True(t, got.Params[1].Optional, "the field of the optional parameter note may be absent")
	})

	t.Run("records the variadic parameter of a call", func(t *testing.T) {
		t.Parallel()

		got := method(t, generate(t, typescript.Target, tsTree, storeModule).models[storeName], "log")
		assert.Equal(t, paramFields(got.Params), "format args", "the call has a field of each parameter")
		expect.True(t, got.Params[1].Variadic, "the field args collects the rest of the arguments")
	})

	t.Run("returns the one value of a method as its result", func(t *testing.T) {
		t.Parallel()

		got := method(t, generate(t, typescript.Target, tsTree, storeModule).models[storeName], "size")
		expect.Equal(t, got.Return, "", "a TypeScript method has no type of its results")
		expect.Equal(t, resultFields(got.Returns), "result", "the call records the value in the field result")
		expect.False(t, got.Named, "a TypeScript method declares no names of its result")
	})

	t.Run("records the awaited value of a method that returns a promise", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, storeModule).models[storeName]
		get := method(t, got, "get")
		assert.Length(t, get.Returns, 1, "get returns a value")
		expect.Equal(t, get.Returns[0].Type.Spelling, "Item", "the value of get is an item")
		expect.Empty(t, method(t, got, "close").Returns, "close returns a promise of no value")
	})

	t.Run("returns a reference to a type of a module without a target", func(t *testing.T) {
		t.Parallel()

		typeOf, declared := templates.TemplateFuncs(typescript.Target)["stubType"].(func(string, string) *emit.TypeRef)
		assert.True(t, declared, "the target typescript declares the helper stubType")
		expect.Equal(t, *typeOf("Seat", "@dokimi/assert"), emit.TypeRef{Spelling: "Seat", Package: "@dokimi/assert"},
			"the reference names the type Seat of the module @dokimi/assert")
	})

	{
		item := &emit.TypeRef{Spelling: "Item"}
		tests := []struct {
			name   string
			give   string
			want   emit.TypeRef
			wantIs string
		}{
			{
				name: "returns the union of a type with undefined", give: "stubOptional",
				want: emit.TypeRef{
					Form: symbol.FormUnion, Elems: []*emit.TypeRef{item, {Spelling: "undefined"}},
				},
				wantIs: "the union of Item and undefined",
			},
			{
				name: "returns a list of a type", give: "stubList",
				want: emit.TypeRef{Form: symbol.FormList, Elems: []*emit.TypeRef{item}}, wantIs: "a list of Item",
			},
			{
				name: "returns the union of a type with a promise of the type", give: "stubAwaitable",
				want: emit.TypeRef{
					Form:  symbol.FormUnion,
					Elems: []*emit.TypeRef{item, {Spelling: "PromiseLike", Args: []*emit.TypeRef{item}}},
				},
				wantIs: "the union of Item and PromiseLike of Item",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				helper, declared := templates.TemplateFuncs(typescript.Target)[tt.give].(func(*emit.TypeRef) *emit.TypeRef)
				assert.True(t, declared, "the target typescript declares the helper "+tt.give)
				expect.Equal(t, *helper(item), tt.want, "the reference is "+tt.wantIs)
			})
		}
	}

	t.Run("writes a name in Pascal case without the initialisms of Go", func(t *testing.T) {
		t.Parallel()

		pascal, declared := templates.TemplateFuncs(typescript.Target)["stubPascal"].(func(string) string)
		assert.True(t, declared, "the target typescript declares the helper stubPascal")
		expect.Equal(t, pascal("findById"), "FindById", "the helper keeps the case of Id")
	})

	t.Run("writes the golden tree typescript", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "typescript", renderedTypeScript(t), golden.ShouldUpdate())
	})

	t.Run("returns the module of the doubles of a spec in its directory", func(t *testing.T) {
		t.Parallel()

		module, declared := templates.TemplateFuncs(typescript.Target)["stubModule"].(func(string) string)
		assert.True(t, declared, "the target typescript declares the helper stubModule")
		expect.Equal(t, module("store.stub.test.ts"), "./store.stub", "the test file imports the module of its doubles")
	})

	t.Run("writes a name in camel case without the initialisms of Go", func(t *testing.T) {
		t.Parallel()

		camel, declared := templates.TemplateFuncs(typescript.Target)["stubCamel"].(func(string) string)
		assert.True(t, declared, "the target typescript declares the helper stubCamel")
		expect.Equal(t, camel("FindByIdStub"), "findByIdStub", "the helper keeps the case of Id")
	})

	t.Run("returns the reference of a double to the builder class of a default result", func(t *testing.T) {
		t.Parallel()

		ref, declared := templates.TemplateFuncs(typescript.Target)["stubRef"].(func(*emit.Struct, string) (*emit.TypeRef, error))
		assert.True(t, declared, "the target typescript declares the helper stubRef")
		got, err := ref(generate(t, typescript.Target, tsTree, builtModule, billingModule).doubles[storeName],
			"AccountBuilder")
		assert.NoError(t, err, "the record of the double of Store refers to AccountBuilder")
		account := fixture(t, typescript.Target, tsTree, stub.New(), builtModule).
			Find(t, symbol.KindStruct, "", "Account").Identity()
		expect.Equal(t, *got, emit.TypeRef{Spelling: "AccountBuilder", Target: account, Emitted: true},
			"the reference names the builder class that the record of the builder of Account declares")
	})
}

// renderedTypeScript runs stub over the modules of the TypeScript tree whose
// doubles render, and returns the files that the TypeScript backend renders
// from them.
func renderedTypeScript(t *testing.T) fstest.MapFS {
	t.Helper()

	p := stub.New()
	f, b := planned(t, typescript.Target, tsTree, p, storeModule, clockModule, handlesModule, streamModule,
		mixedModule, finderModule, pointModule, builtModule, billingModule)
	assert.NoError(t, f.Generate(t, b).Err, "builder runs whole")
	res := f.Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	return gentest.Rendered(t, f, res, p, tsbackend.New(), b)
}
