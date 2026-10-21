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
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees and the packages that this spec alone reads, and the registrant
// that claims the namespace of suite before suite registers its keys.
const (
	scriptedTree    = "testdata/scripted"
	opaquePackage   = "example.com/shop/opaque"
	boxPackage      = "example.com/shop/box"
	badPackage      = "example.com/shop/bad"
	twinPackage     = "example.com/shop/twin"
	shopPackage     = "shop"
	otherRegistrant = "other"
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

		{
			tests := []struct {
				name      string
				giveClaim bool
				giveKey   meta.KeyName
			}{
				{name: "returns from Keys the error of a namespace that another registrant claimed", giveClaim: true},
				{
					name:    "returns from Keys the error of a key of the declaration registered with another type",
					giveKey: suite.KeyDeclaration,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					r := meta.NewRegistry()
					own := r.For(string(suite.ID))
					namespace := suite.KeyDeclaration.Namespace()
					if tt.giveClaim {
						assert.NoError(t, r.For(otherRegistrant).ClaimNamespace(namespace),
							"another registrant claims the namespace of suite")
					} else {
						assert.NoError(t, own.ClaimNamespace(namespace), "suite claims its namespace")
						spec := meta.KeySpec{Name: tt.giveKey, Doc: "a value of another type"}
						_, err := meta.Register[int64](own, spec)
						assert.NoError(t, err, "a key of another type takes the name")
					}
					keys, provides := suite.New().(plugin.KeyProvider)
					assert.True(t, provides, "suite registers keys")
					assert.HasError(t, keys.Keys(own), "suite cannot register its keys")
				})
			}
		}

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

		t.Run("records the falsification test of each suite", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage)
			got, found := recorded(t, f, p).tests["TestStoreSuite"]
			assert.True(t, found, "suite records the falsification test of the suite of Store")
			expect.Equal(t, got.Data.Declaration, "Store", "the test falsifies the suite of Store")
			expect.Equal(t, got.Data.Subject, "", "the falsification test applies the suite to no subject")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 1, "the record refers to the function of the suite")
			store := f.Find(t, symbol.KindInterface, "", "Store").Identity()
			expect.Equal(t, *refs[0].Type, emit.TypeRef{Spelling: "StoreSuite", Target: store, Emitted: true},
				"the record refers to the function that the suite of Store emits")
		})

		t.Run("records the subject type and the callables of an interface", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage)
			got := recorded(t, f, p).suites["Store"]
			assert.NotNil(t, got.Subject, "the suite has a type of subject")
			expect.Equal(t, got.Subject.Spelling, "Store", "the subject type is the interface")
			expect.Equal(t, got.Subject.Target, f.Find(t, symbol.KindInterface, "", "Store").Identity(),
				"the subject type refers to the interface")
			expect.Equal(t, got.Callables, []suite.Callable{
				{Name: "Store.Get", Field: "Get", Type: "StoreGet"},
				{Name: "Store.Put", Field: "Put", Type: "StorePut"},
			}, "each method is a callable in the order of the member walk")
			expect.Nil(t, got.Methods, "an interface needs no interface of its methods")
			expect.Nil(t, got.Func, "an interface has no function type")
		})

		t.Run("records the names of the declarations of the suite", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, storePackage).suites["Store"]
			expect.Equal(t, got.Names, suite.Names{
				Checks: "StoreChecks", Suite: "StoreSuite", Run: "RunStoreSuite",
				RunFromHandle: "RunStoreSuiteFromHandle", Handle: "StoreHandle", Check: "StoreCheck",
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
				{Name: "Catalog.Find", Field: "Find", Type: "CatalogFind"},
				{Name: "Catalog.Add", Field: "Add", Type: "CatalogAdd"},
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
			expect.Equal(t, got.Callables, []suite.Callable{{Name: "Slug", Type: "SlugCallable"}},
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
						"StoreGet", "StorePut", "StoreSuite", "RunStoreSuite", "RunStoreSuiteFromHandle",
						"StoreHandle", "StoreCheck", "StoreOption",
					},
				},
				{
					name: "declares the interface of the methods of a struct", givePkg: catalogPackage, give: "Catalog",
					wantDecl: []string{
						"CatalogMethods", "CatalogFind", "CatalogAdd", "CatalogSuite", "RunCatalogSuite",
						"RunCatalogSuiteFromHandle", "CatalogHandle", "CatalogCheck", "CatalogOption",
					},
				},
				{
					name: "declares the function type of a function", givePkg: slugPackage, give: "Slug",
					wantDecl: []string{
						"SlugFunc", "SlugCallable", "SlugSuite", "RunSlugSuite", "RunSlugSuiteFromHandle",
						"SlugHandle", "SlugCheck", "SlugOption",
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
					wantTest: "TestStoreConformsToStore", wantRun: "RunStoreSuite", wantNew: "New",
				},
				{
					name: "records the conformance test of a subject whose factory takes the handle", give: "Clocked",
					wantTest: "TestClockedConformsToStore", wantRun: "RunStoreSuiteFromHandle", wantNew: "NewClocked",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, storePackage, memoryPackage)
					got, found := recorded(t, f, p).tests[tt.wantTest]
					assert.True(t, found, "suite records the conformance test of "+tt.give)
					expect.Equal(t, got.Data.Subject, tt.give, "the test applies the suite to the subject")
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
					_, res := phases(t, f, p)
					at := f.Find(t, tt.giveKind, "", tt.give)
					got := gentest.Finding(t, res, at)
					expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of its case")
					expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the declaration and the reason")
					gentest.Unrecorded(t, res, at)
				})
			}
		}

		t.Run("binds each check of a consumer to the callable of its directive", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage)
			got := recorded(t, f, p).suites["Store"]
			expect.Equal(t, got.Checks, []suite.Check{
				{
					Apply: ref(t, f, "KeepsTheItem"), Defect: ref(t, f, "forgetful"), Field: "Get", Type: "StoreGet",
					Method: "KeepsTheItem", Segment: "keeps-the-item",
					Statement: "Get returns the item that Put stored",
				},
				{
					Apply: ref(t, f, "RefusesACancelledPut"), Defect: ref(t, f, "careless"), Field: "Put",
					Type: "StorePut", Method: "RefusesACancelledPut", Segment: "refuses-a-cancelled-put",
					Statement: "Put reports a cancelled context",
				},
			}, "each check binds to its callable, in the order of the callables")
		})

		t.Run("binds a check of a consumer to a function", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, slugPackage)
			got := recorded(t, f, p).suites["Slug"]
			expect.Equal(t, got.Checks, []suite.Check{
				{
					Apply: ref(t, f, "KeepsTheOrder"), Defect: ref(t, f, "dropsWords"), Type: "SlugCallable",
					Method: "KeepsTheOrder", Segment: "keeps-the-order",
					Statement: "Slug keeps the words of the title in their order",
				},
				{
					Apply: ref(t, f, "KeepsTheWords"), Defect: ref(t, f, "dropsWords"), Type: "SlugCallable",
					Method: "KeepsTheWords", Segment: "keeps-the-words", Statement: "Slug keeps each word of the title",
				},
			}, "each check binds to the function, with or without the name of its package")
		})

		t.Run("binds a check of a consumer to a method of a struct", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, catalogPackage)
			got := recorded(t, f, p).suites["Catalog"]
			expect.Equal(t, got.Checks, []suite.Check{{
				Apply: ref(t, f, "FindsTheAddedProduct"), Defect: ref(t, f, "amnesiac"), Field: "Find",
				Type: "CatalogFind", Method: "FindsTheAddedProduct", Segment: "finds-the-added-product",
				Statement: "Find returns the product that Add listed",
			}}, "the check binds to the exported method of the struct")
		})

		t.Run("binds a check whose callable opens with the name of its package", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage, twinPackage, badPackage)
			_, res := phases(t, f, p)
			byPackage := gentest.Records(t, res, suiteTemplate, "",
				func(s suite.Suite) string { return s.Subject.Target.Package })
			expect.Equal(t, byPackage[twinPackage].Data.Checks, []suite.Check{{
				Apply: ref(t, f, "KeepsTheName"), Defect: ref(t, f, "broken"), Field: "Get", Type: "StoreGet",
				Method: "KeepsTheName", Segment: "keeps-the-name", Statement: "Get of the twin store returns the name",
			}}, "the check binds to the Store of the package twin alone")
		})

		t.Run("binds no check that it reports", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage, twinPackage, badPackage)
			_, res := phases(t, f, p)
			byPackage := gentest.Records(t, res, suiteTemplate, "",
				func(s suite.Suite) string { return s.Subject.Target.Package })
			checks := byPackage[storePackage].Data.Checks
			got := make([]string, 0, len(checks))
			for _, c := range checks {
				got = append(got, c.Method)
			}
			expect.Equal(t, got, []string{"RefusesACancelledPut"},
				"the suite of store.Store binds the one check of store.Store that suite does not report")
		})

		{
			tests := []struct {
				name     string
				give     string
				wantCode diag.Code
				wantMsg  string
			}{
				{
					name: "reports UnknownCallable at a check of a declaration without the directive suite",
					give: "GetsNothing", wantCode: suite.UnknownCallable,
					wantMsg: "no declaration under the directive suite has the callable Lister.List, " +
						"so GetsNothing is no check",
				},
				{
					name: "reports UnknownCallable at a check of a callable that its declaration lacks",
					give: "FetchesTheItem", wantCode: suite.UnknownCallable,
					wantMsg: "store.Store has no callable Fetch, so FetchesTheItem is no check. " +
						"The callables of store.Store are Get and Put",
				},
				{
					name: "reports UnknownCallable at a check of a declaration with one other callable",
					give: "PutsTheName", wantCode: suite.UnknownCallable,
					wantMsg: "twin.Store has no callable Put, so PutsTheName is no check. " +
						"The callable of twin.Store is Get",
				},
				{
					name: "reports UnknownCallable at a check of a declaration without callables",
					give: "FindsNothing", wantCode: suite.UnknownCallable,
					wantMsg: "Empty has no callable Find, so FindsNothing is no check",
				},
				{
					name: "reports UnknownCallable at a check of a method of a function",
					give: "TrimsTheText", wantCode: suite.UnknownCallable,
					wantMsg: "no declaration under the directive suite has the callable Normalize.Trim, " +
						"so TrimsTheText is no check",
				},
				{
					name: "reports UnknownCallable at a check of a declaration without a method",
					give: "NamesNoCallable", wantCode: suite.UnknownCallable,
					wantMsg: "no declaration under the directive suite has the callable Empty, " +
						"so NamesNoCallable is no check",
				},
				{
					name: "reports UnknownCallable at a check of a name that two declarations have",
					give: "KeepsAnItem", wantCode: suite.UnknownCallable,
					wantMsg: "2 declarations under the directive suite have the name Store, " +
						"so KeepsAnItem is no check. Qualify the callable with the name of its package, " +
						"as store.Store.Get or twin.Store.Get",
				},
				{
					name: "reports BadCheck at a check without a statement", give: "Undocumented",
					wantCode: suite.BadCheck,
					wantMsg: `the doc comment of Undocumented does not open with "Undocumented checks that", ` +
						"so the check has no statement",
				},
				{
					name: "reports BadCheck at a check of two parameters", give: "TakesTwo", wantCode: suite.BadCheck,
					wantMsg: "the check TakesTwo takes 2 parameters, and a check takes the handle of its check alone",
				},
				{
					name: "reports BadCheck at a check with a result", give: "ReturnsOne", wantCode: suite.BadCheck,
					wantMsg: "the check ReturnsOne returns a result, and a check returns nothing",
				},
				{
					name: "reports BadCheck at a defect without the handle", give: "BuildsNothing",
					wantCode: suite.BadCheck,
					wantMsg: "the defect shapeless of BuildsNothing takes 0 parameters, " +
						"and a defect takes the handle of its check alone",
				},
				{
					name: "reports BadCheck at a defect of two results", give: "BuildsTwo", wantCode: suite.BadCheck,
					wantMsg: "the defect paired of BuildsTwo returns 2 results, " +
						"and a defect returns the planted defect alone",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, storePackage, twinPackage, badPackage)
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
					expect.Equal(t, got[0].Code, tt.wantCode, "the finding has the code of its case")
					expect.Equal(t, got[0].Msg, tt.wantMsg, "the message contains the check and the reason")
					key, registered := meta.Lookup[string](f.Keys, suite.KeySegment)
					assert.True(t, registered, "suite registers the key of the segment")
					_, stamped := meta.Get(f.Facts, f.Find(t, symbol.KindFunction, "", tt.give).Identity(), key)
					expect.False(t, stamped, "suite stamps nothing on a check that it reports")
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
