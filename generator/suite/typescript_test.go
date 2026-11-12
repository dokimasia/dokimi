// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	"go.dokimi.dev/eidos/lang/typescript"
	tsbackend "go.dokimi.dev/eidos/lang/typescript/backend"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The modules of the TypeScript tree: the store under the directive suite,
// its subject in memory, an authored check of the store, a function under the
// directive suite, the sink of a value of each form of a shape, the modules
// of a second and a third item, a declaration of one callable, and a store
// whose reads never go back.
const (
	storeModule     = "store"
	memoryModule    = "memory"
	checksModule    = "checks"
	slugModule      = "slug"
	shapesModule    = "shapes"
	otherModule     = "other"
	thirdModule     = "third"
	singleModule    = "single"
	versionsModule  = "versions"
	badChecksModule = "badchecks"
	twinsModule     = "twins"
)

// The part of TypeScript declares the target typescript with the template of
// each family, and its templates write the suites and the conformance tests
// in the conventions of TypeScript. The record of a TypeScript suite passes
// the signal of a probe and the samples as they are, and states the shape of
// each type that a generated binding draws.
func TestTypeScript(t *testing.T) {
	t.Parallel()

	got := generate(t, tsTree, storeModule, memoryModule, slugModule, checksModule)
	store := got.suites["Store"]

	t.Run("takes the statement of an authored check from the summary of its TSDoc", func(t *testing.T) {
		t.Parallel()

		assert.Length(t, store.Checks, 2, "Store has the authored checks keepsTheItem and missesTheUnknown")
		expect.Equal(t, store.Checks[0].Statement, "get returns the item that put stored",
			"the statement is the summary of the check after Checks that")
	})

	t.Run("takes a property check whose second parameter is the case of the property engine", func(t *testing.T) {
		t.Parallel()

		at := slices.IndexFunc(store.Checks, func(c suite.Check) bool { return c.Segment == "misses-the-unknown" })
		assert.True(t, at >= 0, "Store has the authored check missesTheUnknown")
		expect.True(t, store.Checks[at].Property, "missesTheUnknown takes a prop.Case after its handle")
	})

	t.Run("declares the template of each family for the target typescript", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, suite.New(), typescript.Target, suiteTemplate, conformTemplate)
	})

	t.Run("passes the signal of the probe to a parameter of the type AbortSignal", func(t *testing.T) {
		t.Parallel()

		d := derivedNamed(t, &store, "get", "cancel")
		assert.NotNil(t, d.Args[0].Closure, "the probe takes a closure of get")
		expect.Equal(t, d.Args[0].Closure.Args[1], suite.Expr{Word: "signal"},
			"the closure passes its signal to the parameter signal of get")
	})

	t.Run("passes the sample of an input without a conversion", func(t *testing.T) {
		t.Parallel()

		d := derivedNamed(t, &store, "get", "smoke")
		assert.NotNil(t, d.Args[0].Closure, "the probe takes a closure of get")
		sample := d.Args[0].Closure.Args[0].Value
		assert.NotNil(t, sample, "the closure passes a sample of id")
		expect.Equal(t, sample.Kind, emit.ValueLiteral, "a TypeScript literal is a value of each alias of its type")
	})

	t.Run("refuses zero-on-error of each callable that returns a value through returns-beside-error",
		func(t *testing.T) {
			t.Parallel()

			at := slices.IndexFunc(store.Derivation.Refusals, func(r derive.Refusal) bool {
				return r.Callable == "get" && r.Claim == spec.ZeroOnError
			})
			assert.True(t, at >= 0, "the suite refuses zero-on-error of get")
			expect.Equal(t, store.Derivation.Refusals[at].Gate, spec.ReturnsBesideError,
				"a TypeScript call that fails throws and returns no value")
		})

	t.Run("names the alias of the options after the type Options of the runtime", func(t *testing.T) {
		t.Parallel()

		expect.Equal(t, store.Names.Option, "StoreOptions", "the alias of the options is StoreOptions")
	})

	t.Run("states the shape file of each type that a generated binding draws", func(t *testing.T) {
		t.Parallel()

		got := make([][2]string, 0, len(store.Shapes))
		for _, s := range store.Shapes {
			got = append(got, [2]string{s.Name, s.File})
		}
		expect.Equal(t, got, [][2]string{
			{"storeIdShape", `{"shape":"string"}`},
			{"storeItemShape", `{"definitions":{"Item":{"fields":[["id",{"shape":"string"}],` +
				`["count",{"shape":"float","width":64}],["label",{"of":{"shape":"string"},"shape":"optional"}]],` +
				`"shape":"record"}},"name":"Item","shape":"ref"}`},
		}, "the shape of Item is a record of its fields, and its optional field is an optional")
	})

	t.Run("states the generator of each pool and each drawn input of an action", func(t *testing.T) {
		t.Parallel()

		q := generate(t, tsTree, versionsModule).suites["Versions"].Sequences
		assert.NotNil(t, q, "the check of reads-never-go-back runs generated sequences")
		got := make([]string, 0, len(q.Pools))
		for _, p := range q.Pools {
			got = append(got, p.Generator)
		}
		for _, st := range q.Steps {
			for _, p := range st.Params {
				if p.Arg >= 0 && p.Pool == "" {
					got = append(got, p.Generator)
				}
			}
		}
		expect.Equal(t, got, []string{"versionsKeyShape", "versionsRShape", "versionsRShape"},
			"the pool of the keys of version and the records of put and replace take the generators of their types")
	})

	t.Run("gives the shape definitions of two suites of one module distinct names", func(t *testing.T) {
		t.Parallel()

		twins := generate(t, tsTree, twinsModule).suites
		var got []string
		for _, decl := range []string{"Left", "Right"} {
			for _, sh := range twins[decl].Shapes {
				got = append(got, sh.Name)
			}
		}
		expect.Equal(t, got, []string{"leftKeyShape", "rightKeyShape"},
			"each definition of a key joins the name of its declaration")
	})

	t.Run("numbers a shape definition whose name a type in a drawn type has", func(t *testing.T) {
		t.Parallel()

		sink := generate(t, tsTree, shapesModule, storeModule, otherModule, thirdModule).suites["Sink"]
		var got []string
		for _, sh := range sink.Shapes {
			if strings.HasPrefix(sh.Name, "sinkPShape") {
				got = append(got, sh.Name)
			}
		}
		expect.Equal(t, got, []string{"sinkPShape0", "sinkPShape1", "sinkPShape2"},
			"the file imports the class sinkPShape for Box<sinkPShape>, so each definition of a draw p takes a number")
	})

	{
		tests := []struct {
			name    string
			give    string
			wantMsg string
		}{
			{
				name: "reports BadCheck at a check whose second parameter is a case of another module",
				give: "takesAnotherCase",
				wantMsg: "the second parameter of the check takesAnotherCase has the type Case, and a property " +
					"check takes a prop.Case after the handle of its check",
			},
			{
				name: "reports BadCheck at a check whose second parameter is another type of the property engine",
				give: "takesAnOption",
				wantMsg: "the second parameter of the check takesAnOption has the type prop.Option, and a property " +
					"check takes a prop.Case after the handle of its check",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				f := fixture(t, tsTree, p, storeModule, badChecksModule)
				stamps, _ := phases(t, f, p)
				fn := f.Find(t, symbol.KindFunction, "", tt.give).Position()
				var got []diag.Diag
				for d := range stamps.Sink.All() {
					// The directive of each check is on the line above its function.
					if d.Pos.File == fn.File && d.Pos.Line == fn.Line-1 {
						got = append(got, d)
					}
				}
				assert.Length(t, got, 1, "suite reports one finding at the directive of "+tt.give)
				expect.Equal(t, []any{got[0].Code, got[0].Msg}, []any{suite.BadCheck, tt.wantMsg},
					"the finding has the code BadCheck and a message with the check and the reason")
			})
		}
	}

	{
		sink := generate(t, tsTree, shapesModule, storeModule, otherModule, thirdModule).suites["Sink"]
		item := `"Item":{"fields":[["id",{"shape":"string"}],["count",{"shape":"float","width":64}],` +
			`["label",{"of":{"shape":"string"},"shape":"optional"}]],"shape":"record"}`
		tests := []struct {
			name      string
			giveField string
			want      string
		}{
			{name: "states a boolean as a bool", giveField: "flag", want: `{"shape":"bool"}`},
			{name: "states a Uint8Array as bytes", giveField: "bytes", want: `{"shape":"bytes"}`},
			{
				name: "states an array as a list of its elements", giveField: "texts",
				want: `{"of":{"shape":"string"},"shape":"list"}`,
			},
			{
				name: "states an Array of numbers as a list of floats", giveField: "numbers",
				want: `{"of":{"shape":"float","width":64},"shape":"list"}`,
			},
			{
				name: "states a ReadonlyArray of texts as a list of strings", giveField: "frozen",
				want: `{"of":{"shape":"string"},"shape":"list"}`,
			},
			{
				name: "states a bigint as a signed int of 64 bits", giveField: "count",
				want: `{"shape":"int","signed":true,"width":64}`,
			},
			{
				name:      "states an enumeration without values as the literal of the positions of its variants",
				giveField: "level",
				want:      `{"shape":"literal","values":[{"type":"float","value":0},{"type":"float","value":1}]}`,
			},
			{
				name:      "states an enumeration of strings as the literal of its texts",
				giveField: "mode",
				want:      `{"shape":"literal","values":[{"type":"string","value":"read"},{"type":"string","value":"write"}]}`,
			},
			{
				name:      "states an enumeration of numbers as the literal of each number and of the number after it",
				giveField: "status",
				want:      `{"shape":"literal","values":[{"type":"float","value":1},{"type":"float","value":2}]}`,
			},
			{
				name: "refers to the registered name of an enumeration of computed numbers", giveField: "flags",
				want: `{"name":"Flags","shape":"ref"}`,
			},
			{
				name: "refers to the registered name of an enumeration of an infinite number", giveField: "limit",
				want: `{"name":"Limit","shape":"ref"}`,
			},
			{
				name:      "refers to the registered name of an enumeration of a number that is not a number",
				giveField: "odd",
				want:      `{"name":"Odd","shape":"ref"}`,
			},
			{name: "states an alias as its target", giveField: "key", want: `{"shape":"string"}`},
			{
				name: "states a recursive interface as a definition that refers to itself", giveField: "tree",
				want: `{"definitions":{"Tree":{"fields":[["value",{"shape":"string"}],` +
					`["children",{"of":{"name":"Tree","shape":"ref"},"shape":"list"}]],"shape":"record"}},` +
					`"name":"Tree","shape":"ref"}`,
			},
			{
				name: "refers to the registered name of an interface with a method", giveField: "source",
				want: `{"name":"Source","shape":"ref"}`,
			},
			{
				name: "refers to the registered name of a union", giveField: "choice",
				want: `{"name":"string | number","shape":"ref"}`,
			},
			{
				name:      "numbers the definition of a second record with the name of the first",
				giveField: "pair",
				want: `{"definitions":{` + item + `,"Item2":{"fields":[["name",{"shape":"string"}]],"shape":"record"},` +
					`"Pair":{"fields":[["left",{"name":"Item","shape":"ref"}],["right",{"name":"Item2","shape":"ref"}]],` +
					`"shape":"record"}},"name":"Pair","shape":"ref"}`,
			},
			{
				name: "states a Map as a map of its keys to its values", giveField: "lookup",
				want: `{"key":{"shape":"string"},"of":{"shape":"float","width":64},"shape":"map"}`,
			},
			{
				name: "refers to the registered name of a Record with its type arguments", giveField: "table",
				want: `{"name":"Record<string, boolean>","shape":"ref"}`,
			},
			{
				name:      "states a generic interface as a definition of the fields at its type arguments",
				giveField: "box",
				want: `{"definitions":{"Box<number>":{"fields":[["value",{"shape":"float","width":64}]],` +
					`"shape":"record"}},"name":"Box<number>","shape":"ref"}`,
			},
			{
				name: "states a class as a definition of its fields", giveField: "point",
				want: `{"definitions":{"Point":{"fields":[["x",{"shape":"float","width":64}],` +
					`["y",{"shape":"float","width":64}]],"shape":"record"}},"name":"Point","shape":"ref"}`,
			},
			{
				name: "refers to the registered name of a type parameter of a generic method", giveField: "any",
				want: `{"name":"T","shape":"ref"}`,
			},
			{
				name:      "states an enumeration of a negative number as the literal of the number and of the number after it",
				giveField: "signed",
				want:      `{"shape":"literal","values":[{"type":"float","value":-1},{"type":"float","value":0}]}`,
			},
			{
				name: "states a rest parameter as the shape of one of its elements", giveField: "many",
				want: `{"shape":"string"}`,
			},
			{name: "states an optional parameter as its type", giveField: "maybe", want: `{"shape":"string"}`},
			{
				name: "refers to the registered name of an alias of a function type", giveField: "step",
				want: `{"name":"Step","shape":"ref"}`,
			},
			{
				name: "states an alias of an interface as the definition of the interface", giveField: "node",
				want: `{"definitions":{"Tree":{"fields":[["value",{"shape":"string"}],` +
					`["children",{"of":{"name":"Tree","shape":"ref"},"shape":"list"}]],"shape":"record"}},` +
					`"name":"Tree","shape":"ref"}`,
			},
			{
				name: "states an alias of an array as a list of its elements", giveField: "names",
				want: `{"of":{"shape":"string"},"shape":"list"}`,
			},
			{
				name: "states a union with undefined as the shape of its other type", giveField: "unset",
				want: `{"shape":"string"}`,
			},
			{
				name:      "numbers the definition of a third record with the name of the first",
				giveField: "triple",
				want: `{"definitions":{` + item + `,"Item2":{"fields":[["name",{"shape":"string"}]],"shape":"record"},` +
					`"Item3":{"fields":[["flag",{"shape":"bool"}]],"shape":"record"},` +
					`"Triple":{"fields":[["a",{"name":"Item","shape":"ref"}],["b",{"name":"Item2","shape":"ref"}],` +
					`["c",{"name":"Item3","shape":"ref"}]],"shape":"record"}},"name":"Triple","shape":"ref"}`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				expect.Equal(t, called(t, sink, tt.giveField).Params[0].Shape, tt.want,
					"the shape file of the parameter of "+tt.giveField)
			})
		}
	}

	t.Run("starts the draws of an alias at its authored sample", func(t *testing.T) {
		t.Parallel()

		sink := generate(t, tsTree, shapesModule, storeModule, otherModule, thirdModule).suites["Sink"]
		param := called(t, sink, "tag").Params[0]
		assert.NotNil(t, param.Start, "the draws of Tag have a starting value")
		expect.Equal(t, []any{param.Start.Kind, param.Start.Callee.Name}, []any{emit.ValueCall, "newTag"},
			"the starting value of Tag is the call of its authored sample")
		key := called(t, sink, "key").Params[0]
		expect.Nil(t, key.Start, "the draws of Key, an alias without a sample or a default, have no starting value")
	})

	t.Run("lists the required fields of the subject type that a planted defect forwards", func(t *testing.T) {
		t.Parallel()

		expect.Equal(t, store.Fields, []string{"zone"},
			"Store has the required field zone beside its optional field region and its methods")
	})

	t.Run("states the subject of the records of a conformance test with its module", func(t *testing.T) {
		t.Parallel()

		test, recorded := got.tests["TestMemoryConformsToStore"]
		assert.True(t, recorded, "the subject Memory has a conformance test")
		expect.Equal(t, test.Data.Qualified, "memory.Memory", "the subject joins its module and its name")
	})

	t.Run("writes the golden tree typescript", func(t *testing.T) {
		t.Parallel()

		p := suite.New()
		f := fixture(t, tsTree, p, storeModule, memoryModule, checksModule, slugModule, shapesModule, otherModule,
			thirdModule, singleModule, versionsModule)
		_, res := phases(t, f, p)
		golden.MatchTree(t, "typescript", gentest.Rendered(t, f, res, p, tsbackend.New()), golden.ShouldUpdate())
	})
}
