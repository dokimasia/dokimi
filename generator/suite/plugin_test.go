// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees and the packages that this spec alone reads.
const (
	scriptedTree  = "testdata/scripted"
	opaquePackage = "example.com/shop/opaque"
	boxPackage    = "example.com/shop/box"
	badPackage    = "example.com/shop/bad"
	shopPackage   = "shop"
)

// suite records the suite of each declaration under its directive and the
// conformance test of each subject under its directive, and passes the plugin
// suite.
func TestSuite(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the plugin suite over the Go packages store and memory", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := suite.New()
				return p, fixture(tb, goTree, p, storePackage, memoryPackage).Fixture
			})
		})

		t.Run("runs after the generator of the doubles", func(t *testing.T) {
			t.Parallel()

			ordered, orders := suite.New().(plugin.CapabilityProvider)
			assert.True(t, orders, "suite declares its capabilities")
			expect.Equal(t, ordered.Requires(), []plugin.Capability{stub.Capability},
				"suite requires the capability of stub")
		})

		t.Run("tests each subject from inside its package", func(t *testing.T) {
			t.Parallel()

			families, declares := suite.New().(plugin.OutputProvider)
			assert.True(t, declares, "suite declares its families")
			at := slices.IndexFunc(families.Outputs(), func(o plugin.Output) bool { return o.Tag == suite.TestTag })
			assert.NotEqual(t, at, -1, "suite declares the family of the conformance tests")
			expect.False(t, families.Outputs()[at].External, "a conformance test calls a factory that is not exported")
		})

		t.Run("records the subject type and the callables of an interface", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage)
			got := recorded(t, f.Generate(t, p)).suites["Store"]
			assert.NotNil(t, got.Subject, "the suite has a type of subject")
			expect.Equal(t, got.Subject.Spelling, "Store", "the subject type is the interface")
			expect.Equal(t, got.Subject.Target, f.Find(t, symbol.KindInterface, "", "Store").Identity(),
				"the subject type refers to the interface")
			expect.Equal(t, got.Callables, []suite.Callable{
				{Name: "Store.Get", Field: "Get", Type: "storeGet"},
				{Name: "Store.Put", Field: "Put", Type: "storePut"},
			}, "each method is a callable in the order of the member walk")
			expect.Nil(t, got.Methods, "an interface needs no interface of its methods")
			expect.Nil(t, got.Func, "an interface has no function type")
		})

		t.Run("records the names of the declarations of the suite", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, storePackage).suites["Store"]
			expect.Equal(t, got.Names, suite.Names{
				Checks: "StoreChecks", Suite: "StoreSuite", Run: "RunStore", RunOn: "RunStoreOn",
				Control: "ControlStore", ControlOn: "ControlStoreOn", Handle: "StoreHandle", Check: "StoreCheck",
				Option: "StoreOption",
			}, "each name joins the name of the declaration to the word of its part")
		})

		t.Run("records the subject type of a generic interface at its witnesses", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, genericPackage).suites["Cache"]
			assert.NotNil(t, got.Subject, "the suite has a type of subject")
			args := make([]string, 0, len(got.Subject.Args))
			for _, a := range got.Subject.Args {
				args = append(args, a.Spelling)
			}
			expect.Equal(t, args, []string{"int", "int"}, "the subject type instantiates Cache at its witnesses")
		})

		t.Run("records the interface of the exported methods of a struct", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, catalogPackage).suites["Catalog"]
			assert.NotNil(t, got.Methods, "the suite has an interface of the methods")
			methods := make([]string, 0, got.Methods.Methods.Len())
			for _, m := range got.Methods.Methods.Items() {
				methods = append(methods, m.Name)
			}
			expect.Equal(t, methods, []string{"Find", "Add"}, "the interface has each exported method")
			expect.Equal(t, got.Names.Methods, "CatalogMethods", "the interface joins Methods to the struct")
			expect.Equal(t, got.Subject.Spelling, "CatalogMethods", "the subject type is the interface of the methods")
			expect.Equal(t, got.Callables, []suite.Callable{
				{Name: "Catalog.Find", Field: "Find", Type: "catalogFind"},
				{Name: "Catalog.Add", Field: "Add", Type: "catalogAdd"},
			}, "each exported method is a callable")
		})

		t.Run("records the function type of a function", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, slugPackage).suites["Slug"]
			assert.NotNil(t, got.Func, "the suite has a function type")
			params := make([]string, 0, got.Func.Params.Len())
			for _, p := range got.Func.Params.Items() {
				params = append(params, p.Name)
			}
			expect.Equal(t, params, []string{"title", "keep"}, "the function type has the parameters of Slug")
			expect.Equal(t, got.Names.Func, "SlugFunc", "the function type joins Func to the name of the function")
			expect.Equal(t, got.Subject.Spelling, "SlugFunc", "the subject type is the function type")
			expect.Equal(t, got.Callables, []suite.Callable{{Name: "Slug", Type: "slugCallable"}},
				"the function is its one callable")
		})

		{
			tests := []struct {
				name     string
				givePkg  string
				give     string
				wantDecl []string
			}{
				{
					name: "declares each name of the suite of an interface", givePkg: storePackage, give: "Store",
					wantDecl: []string{
						"storeGet", "storePut", "StoreSuite", "RunStore", "RunStoreOn", "ControlStore",
						"ControlStoreOn", "StoreHandle", "StoreCheck", "StoreOption",
					},
				},
				{
					name: "declares the interface of the methods of a struct", givePkg: catalogPackage, give: "Catalog",
					wantDecl: []string{
						"CatalogMethods", "catalogFind", "catalogAdd", "CatalogSuite", "RunCatalog", "RunCatalogOn",
						"ControlCatalog", "ControlCatalogOn", "CatalogHandle", "CatalogCheck", "CatalogOption",
					},
				},
				{
					name: "declares the function type of a function", givePkg: slugPackage, give: "Slug",
					wantDecl: []string{
						"SlugFunc", "slugCallable", "SlugSuite", "RunSlug", "RunSlugOn", "ControlSlug",
						"ControlSlugOn", "SlugHandle", "SlugCheck", "SlugOption",
					},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					expect.Equal(t, generate(t, goTree, tt.givePkg).beside[tt.give], tt.wantDecl,
						"the record declares each name that the template writes beside the names of the callables")
				})
			}
		}

		{
			tests := []struct {
				name     string
				give     string
				wantTest string
				wantRun  string
				wantNew  string
			}{
				{
					name: "records the conformance test of a subject whose factory takes nothing", give: "Store",
					wantTest: "TestStoreConformsToStore", wantRun: "RunStore", wantNew: "New",
				},
				{
					name: "records the conformance test of a subject whose factory takes the handle", give: "Clocked",
					wantTest: "TestClockedConformsToStore", wantRun: "RunStoreOn", wantNew: "NewClocked",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, storePackage, memoryPackage)
					got, found := recorded(t, f.Generate(t, p)).tests[tt.give]
					assert.True(t, found, "suite records the conformance test of "+tt.give)
					expect.Equal(t, got.Data.Test, tt.wantTest, "the test joins the subject and the declaration")
					expect.Equal(t, got.Data.Declaration, "Store", "the test applies the suite of Store")
					assert.NotNil(t, got.Decl, "the record of the test is a struct")
					refs := got.Decl.Fields.Items()
					assert.Length(t, refs, 2, "the record refers to the application and to the factory")
					store := f.Find(t, symbol.KindInterface, "", "Store").Identity()
					expect.Equal(t, *refs[0].Type, emit.TypeRef{Spelling: tt.wantRun, Target: store, Emitted: true},
						"the record refers to the application that the suite of Store emits")
					factory := refs[1].Type
					expect.Equal(t, factory.Spelling, tt.wantNew, "the record refers to the factory")
					expect.Equal(t, factory.Target.Package, memoryPackage, "the factory is in the package memory")
					expect.Equal(t, factory.Target.Kind, symbol.KindFunction, "the factory is a function")
					expect.Equal(t, factory.Target.Name, tt.wantNew, "the reference has the identity of the factory")
				})
			}
		}

		{
			tests := []struct {
				name     string
				givePkgs []string
				giveKind symbol.Kind
				give     string
				wantCode diag.Code
				wantMsg  string
			}{
				{
					name:     "reports Uninstantiable at a generic interface without witnesses",
					givePkgs: []string{opaquePackage}, giveKind: symbol.KindInterface, give: "Ranked",
					wantCode: suite.Uninstantiable,
					wantMsg:  "the type parameters of Ranked have no witnesses, so Ranked has no suite",
				},
				{
					name: "reports Uninstantiable at a generic struct", givePkgs: []string{boxPackage},
					giveKind: symbol.KindStruct, give: "Box", wantCode: suite.Uninstantiable,
					wantMsg: "Box is a generic struct, so it has no suite",
				},
				{
					name: "reports Uninstantiable at a generic function", givePkgs: []string{boxPackage},
					giveKind: symbol.KindFunction, give: "Mapped", wantCode: suite.Uninstantiable,
					wantMsg: "Mapped is a generic function, so it has no suite",
				},
				{
					name:     "reports NoSuite at a subject of a declaration without the directive suite",
					givePkgs: []string{storePackage, badPackage}, giveKind: symbol.KindStruct, give: "Plain",
					wantCode: suite.NoSuite,
					wantMsg:  "Lister has no directive suite, so Plain has no conformance test",
				},
				{
					name:     "reports BadFactory at a subject whose factory takes two parameters",
					givePkgs: []string{storePackage, badPackage}, giveKind: symbol.KindStruct, give: "Twice",
					wantCode: suite.BadFactory,
					wantMsg:  "the factory NewTwice takes 2 parameters, so Twice has no conformance test",
				},
				{
					name:     "reports BadFactory at a subject whose factory returns two results",
					givePkgs: []string{storePackage, badPackage}, giveKind: symbol.KindStruct, give: "Paired",
					wantCode: suite.BadFactory,
					wantMsg:  "the factory NewPaired returns 2 results, so Paired has no conformance test",
				},
				{
					name:     "reports BadFactory at a subject whose factory returns no result",
					givePkgs: []string{storePackage, badPackage}, giveKind: symbol.KindStruct, give: "Nothing",
					wantCode: suite.BadFactory,
					wantMsg:  "the factory NewNothing returns 0 results, so Nothing has no conformance test",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, tt.givePkgs...)
					res := f.Generate(t, p)
					assert.NoError(t, res.Err, "the generator runs whole")
					at := f.Find(t, tt.giveKind, "", tt.give)
					got := gentest.Finding(t, res, at)
					expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of its case")
					expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the declaration and the reason")
					gentest.Unrecorded(t, res, at)
				})
			}
		}

		t.Run("records nothing for a declaration of a language that suite has no part for", func(t *testing.T) {
			t.Parallel()

			got := generate(t, scriptedTree, shopPackage)
			expect.Empty(t, got.suites, "suite records no suite of the scripted struct Widget")
			expect.Empty(t, got.findings, "suite reports nothing at Widget")
		})
	})
}
