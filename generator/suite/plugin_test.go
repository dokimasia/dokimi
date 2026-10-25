// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"fmt"
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
	"go.dokimi.dev/dokimi/internal/weave"
)

// The trees and the packages that this spec alone reads, and the registrant
// that claims the namespace of suite before suite registers its keys.
const (
	scriptedTree    = "testdata/scripted"
	opaquePackage   = "example.com/shop/opaque"
	boxPackage      = "example.com/shop/box"
	badPackage      = "example.com/shop/bad"
	twinPackage     = "example.com/shop/twin"
	reservedPackage = "example.com/shop/reserved"
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

		t.Run("passes the plugin suite over two Go packages", func(t *testing.T) {
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

		t.Run("provides its capability in the annotator role alone", func(t *testing.T) {
			t.Parallel()

			ordered, orders := suite.New().(plugin.CapabilityProvider)
			assert.True(t, orders, "suite declares its capabilities")
			expect.Equal(t, ordered.Provides(plugin.RoleAnnotator), []plugin.Capability{weave.Suite},
				"the annotate half provides the capability after which stub runs")
			expect.Empty(t, ordered.Provides(plugin.RoleGenerator),
				"the generate half provides no capability, so it runs after stub")
		})

		doubles := []struct {
			name     string
			giveKind symbol.Kind
			give     string
			want     bool
		}{
			{
				name: "stamps the fact of a double on an interface", giveKind: symbol.KindInterface, give: "Store",
				want: true,
			},
			{name: "stamps the fact of a double on a struct", giveKind: symbol.KindStruct, give: "Catalog", want: true},
			{name: "stamps no fact of a double on a function", giveKind: symbol.KindFunction, give: "Slug"},
		}
		for _, tt := range doubles {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				f := fixture(t, goTree, p, storePackage, catalogPackage, slugPackage)
				phases(t, f, p)
				k, registered := meta.Lookup[bool](f.Keys, weave.KeyDouble)
				assert.True(t, registered, "suite registers the key of the double")
				got, _ := meta.Get(f.Facts, f.Find(t, tt.giveKind, "", tt.give).Identity(), k)
				expect.Equal(t, got, tt.want, "the fact asks stub for the double of the declaration")
			})
		}

		t.Run("declares that the directive excuse requires the directive subject", func(t *testing.T) {
			t.Parallel()

			declared, declares := suite.New().(plugin.DirectiveProvider)
			assert.True(t, declares, "suite declares its directives")
			schemas := declared.Directives()
			i := slices.IndexFunc(schemas, func(s directive.Schema) bool { return s.Name == suite.ExcuseDirective })
			assert.True(t, i >= 0, "suite declares the directive excuse")
			expect.Equal(t, schemas[i].Requires, []directive.Name{suite.SubjectDirective},
				"an excusal applies to the conformance test of a subject")
		})

		t.Run("writes its specs into the package with the suffix _test", func(t *testing.T) {
			t.Parallel()

			gentest.External(t, suite.New(), suite.TestTag)
		})

		t.Run("records the falsification of each suite", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage)
			got, found := recorded(t, f, p).tests["TestStoreSuite"]
			assert.True(t, found, "suite records the falsification test of the suite of Store")
			expect.Equal(t, got.Data.Declaration, "Store", "the test falsifies the suite of Store")
			expect.Equal(t, got.Data.Subject, "", "the falsification test applies the suite to no subject")
			expect.False(t, got.Data.Function, "the suite of an interface applies to the subjects of a factory")
			expect.Equal(t, got.Data.Builds, "builds-the-subject", "the check of a factory has the segment Builds")
			expect.Equal(t, got.Data.Receives, "receives-the-function",
				"the check of the form Run of a function has the segment Receives")
			alone := "the test applies the check of its option alone"
			expect.Equal(t, got.Data.Excusals, []suite.Excusal{
				{Field: "Get", Method: "KeepsTheItem", Reason: alone},
				{Field: "Put", Method: "RefusesACancelledPut", Reason: alone},
			}, "the test excuses each check of the suite where it applies a check of its own")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			store := f.Find(t, symbol.KindInterface, "", "Store").Identity()
			emitted := func(name string) emit.TypeRef { return emit.TypeRef{Spelling: name, Target: store, Emitted: true} }
			refs := make([]emit.TypeRef, 0, got.Decl.Fields.Len())
			for _, field := range got.Decl.Fields.Items() {
				refs = append(refs, *field.Type)
			}
			expect.Equal(t, refs, []emit.TypeRef{
				emitted("StoreSuite"),
				emitted("RunStoreSuite"),
				emitted("RunStoreSuiteFromHandle"),
				emitted("StoreHandle"),
				emitted("StoreCheck"),
				{Spelling: "Store", Target: store},
				emitted("StoreChecks"),
				emitted("FuzzStoreSuite"),
				emitted("RunStoreBudgets"),
				emitted("BenchStoreSuite"),
			}, "the record refers to the suite, its two applications, its aliases, the subject type, the callables, "+
				"the fuzz form, the application of the allocation ceilings and the benchmark")
		})

		t.Run("records the parallel benchmark of a suite with a concurrent callable", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, meterPackage)
			got := recorded(t, f, p)
			expect.Equal(t, got.suites["Meter"].Names.BenchParallel, "BenchMeterSuiteParallel",
				"the suite declares the parallel benchmark of its concurrent callable Count")
			test, found := got.tests["TestMeterSuite"]
			assert.True(t, found, "suite records the falsification test of the suite of Meter")
			refs := test.Decl.Fields.Items()
			assert.Length(t, refs, 11, "the record has a field of each name of the test")
			meter := f.Find(t, symbol.KindInterface, "", "Meter").Identity()
			parallel := emit.TypeRef{Spelling: "BenchMeterSuiteParallel", Target: meter, Emitted: true}
			expect.Equal(t, *refs[10].Type, parallel,
				"the record refers to the parallel benchmark that the suite of Meter emits")
		})

		{
			tests := []struct {
				name         string
				givePkg      string
				giveKind     symbol.Kind
				give         string
				wantSubject  string
				wantFields   int
				wantFunction bool
			}{
				{
					name:    "records the interface of the methods of a struct as the type of its subjects",
					givePkg: catalogPackage, giveKind: symbol.KindStruct, give: "Catalog",
					wantSubject: "CatalogMethods", wantFields: 10,
				},
				{
					name:    "records the function type of a function as the type of its subjects",
					givePkg: slugPackage, giveKind: symbol.KindFunction, give: "Slug",
					wantSubject: "SlugFunc", wantFields: 11, wantFunction: true,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, tt.givePkg)
					got, found := recorded(t, f, p).tests["Test"+tt.give+"Suite"]
					assert.True(t, found, "suite records the test of the suite of "+tt.give)
					expect.Equal(t, got.Data.Function, tt.wantFunction, "the test applies the function itself")
					assert.NotNil(t, got.Decl, "the record of the test is a struct")
					refs := got.Decl.Fields.Items()
					assert.Length(t, refs, tt.wantFields, "the record has a field of each name of the test")
					origin := f.Find(t, tt.giveKind, "", tt.give).Identity()
					want := emit.TypeRef{Spelling: tt.wantSubject, Target: origin, Emitted: true}
					expect.Equal(t, *refs[5].Type, want, "the record refers to the subject type that the suite emits")
				})
			}
		}

		t.Run("records the function of the suite of a function", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, slugPackage)
			got, found := recorded(t, f, p).tests["TestSlugSuite"]
			assert.True(t, found, "suite records the test of the suite of Slug")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 11, "the record has a field of each name of the test")
			slug := f.Find(t, symbol.KindFunction, "", "Slug").Identity()
			expect.Equal(t, *refs[7].Type, emit.TypeRef{Spelling: "Slug", Target: slug},
				"the record refers to the function that the form Run applies")
		})

		t.Run("records the fuzz target of the suite of a function", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, slugPackage)
			got, found := recorded(t, f, p).tests["TestSlugSuite"]
			assert.True(t, found, "suite records the test of the suite of Slug")
			expect.Equal(t, got.Data.Targets, []suite.Target{{Name: "FuzzSlug"}},
				"the function is the one callable of its suite and the one subject of its fuzz target")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 11, "the record has a field of each name of the test")
			slug := f.Find(t, symbol.KindFunction, "", "Slug").Identity()
			expect.Equal(t, *refs[8].Type, emit.TypeRef{Spelling: "FuzzSlugSuite", Target: slug, Emitted: true},
				"the record refers to the fuzz form that the suite of Slug emits")
		})

		t.Run("records a fuzz target of each callable of the suite of an interface", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, storePackage).tests["TestStoreSuite"]
			assert.True(t, found, "suite records the test of the suite of Store")
			expect.Equal(t, got.Data.Targets, []suite.Target{
				{Name: "FuzzStoreSuiteGet", Field: "Get"},
				{Name: "FuzzStoreSuitePut", Field: "Put"},
			}, "each target fuzzes the planted defect of a ceiling, in the order of the callables")
		})

		{
			tests := []struct {
				name     string
				givePkgs []string
				give     string
				want     []suite.Identity
			}{
				{
					name: "records the identity of each check of the consumer", givePkgs: []string{storePackage},
					give: "TestStoreSuite",
					want: []suite.Identity{
						{Field: "Get", Method: "KeepsTheItem", ID: "Store.Get/keeps-the-item"},
						{Field: "Put", Method: "RefusesACancelledPut", ID: "Store.Put/refuses-a-cancelled-put"},
					},
				},
				{
					name: "records the identity of each ceiling of a budget", givePkgs: []string{meterPackage},
					give: "TestMeterSuite",
					want: []suite.Identity{
						{Field: "Count", Method: "AllocationBudget", ID: "Meter.Count/allocation-budget"},
						{Field: "Count", Method: "ByteBudget", ID: "Meter.Count/byte-budget"},
						{Field: "Count", Method: "LatencyBudget", ID: "Meter.Count/latency-budget"},
						{Field: "Count", Method: "MeanLatencyBudget", ID: "Meter.Count/mean-latency-budget"},
						{Field: "Record", Method: "LatencyBudget", ID: "Meter.Record/latency-budget"},
					},
				},
				{
					name: "records the identity of each check of a function", givePkgs: []string{slugPackage},
					give: "TestSlugSuite",
					want: []suite.Identity{
						{Method: "KeepsTheOrder", ID: "Slug/keeps-the-order"},
						{Method: "KeepsTheWords", ID: "Slug/keeps-the-words"},
					},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got, found := generate(t, goTree, tt.givePkgs...).tests[tt.give]
					assert.True(t, found, "suite records the test "+tt.give)
					expect.Equal(t, got.Data.Identities, tt.want,
						"the test compares each identity with the method of the names of the callables")
				})
			}
		}

		t.Run("records no idle test for the suite of an interface without callables", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, pathsPackage, inputsPackage).tests["TestIdleSuite"]
			assert.True(t, found, "suite records the test of the suite of Idle")
			expect.Equal(t, got.Data.Budgets, "", "a suite without callables has no ceiling to apply")
			expect.Equal(t, got.Data.Benchmark, "", "a suite without callables has no callable to measure")
			expect.Empty(t, got.Data.Targets, "a suite without callables has no callable to fuzz")
		})

		{
			tests := []struct {
				name       string
				givePkgs   []string
				give       string
				wantLevels map[string]suite.Level
			}{
				{
					name:     "records the 0 ceiling at each level of a suite without a budget",
					givePkgs: []string{meterPackage}, give: "TestMeterSuite",
					wantLevels: map[string]suite.Level{
						"Count":  {Field: "Count", Concurrent: true},
						"Record": {Field: "Record"},
						"Merge":  {Field: "Merge", Ceiled: true},
						"Total":  {Field: "Total", Ceiled: true},
					},
				},
				{
					name:     "records no 0 ceiling at a level of a suite whose calls build an input",
					givePkgs: []string{pathsPackage, inputsPackage}, give: "TestPipelineSuite",
					wantLevels: map[string]suite.Level{
						"Pump":  {Field: "Pump"},
						"Fold":  {Field: "Fold"},
						"Label": {Field: "Label", Ceiled: true},
					},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got, found := generate(t, goTree, tt.givePkgs...).tests[tt.give]
					assert.True(t, found, "suite records the test "+tt.give)
					levels := make(map[string]suite.Level, len(tt.wantLevels))
					for _, l := range got.Data.Levels {
						if _, want := tt.wantLevels[l.Field]; want {
							levels[l.Field] = l
						}
					}
					expect.Equal(t, levels, tt.wantLevels,
						"a level without a budget whose calls build no input takes the 0 ceiling of allocations")
				})
			}
		}

		t.Run("records the callables of an interface with its subject type", func(t *testing.T) {
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
				RunFromHandle: "RunStoreSuiteFromHandle", Fuzz: "FuzzStoreSuite", Budgets: "RunStoreBudgets",
				Bench: "BenchStoreSuite", Handle: "StoreHandle", Check: "StoreCheck", Option: "StoreOption",
				Subjects: "StoreSubjects", Over: "StoreOver", NewOver: "NewStoreOver", Embedded: "StoreOverEmbedded",
			}, "each name joins the name of the declaration to the word of its part, and a suite without a "+
				"concurrent callable has no parallel benchmark")
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
						"StoreGet", "StorePut", "StoreGetCalls", "StorePutCalls", "StoreOver", "StoreOverEmbedded",
						"NewStoreOver",
						"StoreSuite", "RunStoreSuite", "RunStoreSuiteFromHandle", "FuzzStoreSuite", "RunStoreBudgets",
						"BenchStoreSuite", "StoreHandle", "StoreCheck", "StoreOption", "StoreSubjects",
					},
				},
				{
					name:    "declares the parallel benchmark of a suite with a concurrent callable",
					givePkg: meterPackage, give: "Meter",
					wantDecl: []string{
						"MeterCount", "MeterRecord", "MeterMerge", "MeterTotal", "MeterCountCalls",
						"MeterRecordCalls", "MeterMergeCalls", "MeterTotalCalls", "MeterOver", "MeterOverEmbedded",
						"NewMeterOver",
						"MeterSuite", "RunMeterSuite", "RunMeterSuiteFromHandle", "FuzzMeterSuite", "RunMeterBudgets",
						"BenchMeterSuite", "BenchMeterSuiteParallel", "MeterHandle", "MeterCheck", "MeterOption",
						"MeterSubjects",
					},
				},
				{
					name: "declares the interface of the methods of a struct", givePkg: catalogPackage, give: "Catalog",
					wantDecl: []string{
						"CatalogMethods", "CatalogFind", "CatalogAdd", "CatalogFindCalls", "CatalogAddCalls",
						"CatalogOver", "CatalogOverEmbedded", "NewCatalogOver", "CatalogSuite", "RunCatalogSuite",
						"RunCatalogSuiteFromHandle", "FuzzCatalogSuite", "RunCatalogBudgets", "BenchCatalogSuite",
						"CatalogHandle", "CatalogCheck", "CatalogOption", "CatalogSubjects",
					},
				},
				{
					name: "declares the function type of a function", givePkg: slugPackage, give: "Slug",
					wantDecl: []string{
						"SlugFunc", "SlugCallable", "SlugCalls", "NewSlugOver", "SlugSuite", "RunSlugSuite",
						"RunSlugSuiteFromHandle", "FuzzSlugSuite", "RunSlugBudgets", "BenchSlugSuite", "SlugHandle",
						"SlugCheck", "SlugOption", "SlugSubjects",
					},
				},
				{
					name: "declares each generator and each implementation that the calls take", givePkg: pathsPackage,
					give: "Route",
					wantDecl: []string{
						"RouteFunc", "RouteCallable", "RouteCalls", "RouteItemValues", "RouteDepotInput",
						"RouteDepotInputEmbedded", "NewRouteOver", "RouteSuite", "RunRouteSuite",
						"RunRouteSuiteFromHandle", "FuzzRouteSuite", "RunRouteBudgets", "BenchRouteSuite",
						"RouteHandle", "RouteCheck", "RouteOption", "RouteSubjects",
					},
				},
				{
					name: "declares no planted defect of an interface without callables", givePkg: pathsPackage,
					give: "Idle",
					wantDecl: []string{
						"IdleSuite", "RunIdleSuite", "RunIdleSuiteFromHandle", "FuzzIdleSuite", "RunIdleBudgets",
						"BenchIdleSuite", "IdleHandle", "IdleCheck", "IdleOption", "IdleSubjects",
					},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					expect.Equal(t, generate(t, goTree, tt.givePkg, inputsPackage).beside[tt.give], tt.wantDecl,
						"the record declares each name that the template writes beside the names of the callables")
				})
			}
		}

		{
			tests := []struct {
				name           string
				give           string
				wantTest       string
				wantRun        string
				wantNew        string
				wantFromHandle bool
			}{
				{
					name: "records the conformance of a subject whose factory takes nothing", give: "Store",
					wantTest: "TestStoreConformsToStore", wantRun: "RunStoreSuite", wantNew: "New",
				},
				{
					name: "records the conformance of a subject whose factory takes the handle", give: "Clocked",
					wantTest: "TestClockedConformsToStore", wantRun: "RunStoreSuiteFromHandle", wantNew: "NewClocked",
					wantFromHandle: true,
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
					expect.Equal(t, got.Data.FromHandle, tt.wantFromHandle, "the factory takes the handle of a check")
					assert.NotNil(t, got.Decl, "the record of the test is a struct")
					refs := got.Decl.Fields.Items()
					assert.Length(t, refs, 10, "the record refers to the application, the factory, the names of the "+
						"callables, the fuzz form, the handle, the subject type of the factory, the application of the "+
						"allocation ceilings, the benchmark, the option and the subject type of the suite")
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

		t.Run("records the conformance of a subject with an exported factory in the package with the suffix _test",
			func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				f := fixture(t, goTree, p, storePackage, memoryPackage)
				got := recorded(t, f, p)
				subject := f.Find(t, symbol.KindStruct, "", "Store").Identity()
				expect.False(t, got.inside[subject], "the external test package calls an exported factory")
				_, found := got.tests["TestStoreConformsToStore"]
				expect.True(t, found, "suite records the conformance test of Store")
			})

		{
			tests := []struct {
				name string
				give string
			}{
				{
					name: "records the conformance of a subject with a factory of its external test package",
					give: "Store",
				},
				{
					name: "records the conformance of a subject with a factory of its external test package by its " +
						"import path",
					give: "Vault",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, storePackage, speccedPackage, speccedTestPackage)
					got := recorded(t, f, p)
					subject := f.Find(t, symbol.KindStruct, "", tt.give).Identity()
					expect.False(t, got.inside[subject], "the external test package calls a factory of its own")
					test, found := got.tests["Test"+tt.give+"ConformsToStore"]
					assert.True(t, found, "suite records the conformance test of "+tt.give)
					factory := test.Decl.Fields.Items()[1].Type
					expect.Equal(t, factory.Target.Package, speccedTestPackage,
						"the record refers to the factory of the external test package")
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
					name:     "reports BadFactory at a subject whose factory its package does not export",
					givePkgs: []string{storePackage, badPackage}, giveKind: symbol.KindStruct, give: "Hidden",
					wantCode: suite.BadFactory,
					wantMsg: "the factory newHidden is not exported, so Hidden has no conformance test. Export it, " +
						"or declare it in the external test package",
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

		t.Run("records a fuzz target of each callable of a subject", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, storePackage, memoryPackage).tests["TestStoreConformsToStore"]
			assert.True(t, found, "suite records the conformance test of Store")
			expect.Equal(t, got.Data.Targets, []suite.Target{
				{Name: "FuzzStoreConformsToStoreGet", Field: "Get"},
				{Name: "FuzzStoreConformsToStorePut", Field: "Put"},
			}, "each target joins the subject, the declaration and the callable, in the order of the callables")
		})

		t.Run("records the references of the fuzz targets of a subject", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage, memoryPackage)
			got, found := recorded(t, f, p).tests["TestStoreConformsToStore"]
			assert.True(t, found, "suite records the conformance test of Store")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 10, "the record has a field of each name of the test")
			store := f.Find(t, symbol.KindInterface, "", "Store").Identity()
			emitted := func(name string) emit.TypeRef { return emit.TypeRef{Spelling: name, Target: store, Emitted: true} }
			expect.Equal(t, []emit.TypeRef{*refs[2].Type, *refs[3].Type, *refs[4].Type},
				[]emit.TypeRef{emitted("StoreChecks"), emitted("FuzzStoreSuite"), emitted("StoreHandle")},
				"the record refers to the names of the callables, the fuzz form and the handle of the suite")
			expect.Equal(t, refs[5].Type.Spelling, "*Store", "the record refers to the subject type of the factory")
		})

		t.Run("records the references of the entry points of the budgets of a subject", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage, memoryPackage)
			got, found := recorded(t, f, p).tests["TestStoreConformsToStore"]
			assert.True(t, found, "suite records the conformance test of Store")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 10, "the record has a field of each name of the test")
			store := f.Find(t, symbol.KindInterface, "", "Store").Identity()
			emitted := func(name string) emit.TypeRef { return emit.TypeRef{Spelling: name, Target: store, Emitted: true} }
			expect.Equal(t, []emit.TypeRef{*refs[6].Type, *refs[7].Type, *refs[8].Type, *refs[9].Type},
				[]emit.TypeRef{
					emitted("RunStoreBudgets"),
					emitted("BenchStoreSuite"),
					emitted("StoreOption"),
					{Spelling: "Store", Target: store},
				},
				"the record refers to the application of the allocation ceilings, the benchmark, the option and "+
					"the subject type of the suite")
		})

		t.Run("records the parallel benchmark of a subject with a concurrent callable", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, meterPackage, meteredPackage)
			got, found := recorded(t, f, p).tests["TestMeterConformsToMeter"]
			assert.True(t, found, "suite records the conformance test of Meter")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 11, "the record has a field of each name of the test")
			meter := f.Find(t, symbol.KindInterface, "", "Meter").Identity()
			parallel := emit.TypeRef{Spelling: "BenchMeterSuiteParallel", Target: meter, Emitted: true}
			expect.Equal(t, *refs[10].Type, parallel,
				"the record refers to the parallel benchmark that the suite of Meter emits")
		})

		{
			tests := []struct {
				name     string
				givePkgs []string
				give     string
				want     string
				wantArgs []string
			}{
				{
					name:     "records the type of the subjects of the suite of a struct",
					givePkgs: []string{catalogPackage}, give: "TestCatalogConformsToCatalog", want: "CatalogMethods",
				},
				{
					name:     "records the type of the subjects of the suite of a generic interface at its witnesses",
					givePkgs: []string{genericPackage, cachedPackage}, give: "TestCacheConformsToCache", want: "Cache",
					wantArgs: []string{"int", "int"},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got, found := generate(t, goTree, tt.givePkgs...).tests[tt.give]
					assert.True(t, found, "suite records the test "+tt.give)
					refs := got.Decl.Fields.Items()
					assert.Length(t, refs, 10, "the record has a field of each name of the test")
					subject := refs[9].Type
					args := make([]string, 0, len(subject.Args))
					for _, a := range subject.Args {
						args = append(args, a.Spelling)
					}
					expect.Equal(t, subject.Spelling, tt.want, "the record refers to the type of the subjects")
					expect.Equal(t, args, tt.wantArgs, "the type takes the witnesses of the suite",
						assert.EquateEmpty())
				})
			}
		}

		t.Run("records the number of cases that a subject states", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, storePackage, memoryPackage).tests["TestClockedConformsToStore"]
			assert.True(t, found, "suite records the conformance test of Clocked")
			expect.Equal(t, got.Data.Cases, int64(20), "the test passes the cases to each property run")
		})

		t.Run("records the references of the entry points of the budgets of a function", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, meterPackage)
			got, found := recorded(t, f, p).tests["TestWeighSuite"]
			assert.True(t, found, "suite records the test of the suite of Weigh")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 12, "the record has a field of each name of the test")
			weigh := f.Find(t, symbol.KindFunction, "", "Weigh").Identity()
			emitted := func(name string) emit.TypeRef { return emit.TypeRef{Spelling: name, Target: weigh, Emitted: true} }
			want := []emit.TypeRef{
				emitted("RunWeighBudgets"),
				emitted("BenchWeighSuite"),
				emitted("BenchWeighSuiteParallel"),
			}
			expect.Equal(t, []emit.TypeRef{*refs[9].Type, *refs[10].Type, *refs[11].Type}, want,
				"the record refers to the application of the allocation ceilings and the two benchmarks of the suite "+
					"of a concurrent function")
		})

		t.Run("records each excusal of a subject in the order of its directives", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, storePackage, memoryPackage)
			got, found := recorded(t, f, p).tests["TestForgetfulConformsToStore"]
			assert.True(t, found, "suite records the conformance test of Forgetful")
			expect.Equal(t, got.Data.Excusals, []suite.Excusal{
				{Field: "Get", Method: "KeepsTheItem", Reason: "the store drops each item"},
				{
					Field: "Put", Method: "RefusesACancelledPut",
					Reason: "the store drops each item before it reads the context",
				},
			}, "each excusal names the callable and the method of the identity of its check")
			assert.NotNil(t, got.Decl, "the record of the test is a struct")
			refs := got.Decl.Fields.Items()
			assert.Length(t, refs, 10, "the record has a field of each name of the test")
			store := f.Find(t, symbol.KindInterface, "", "Store").Identity()
			expect.Equal(t, *refs[2].Type, emit.TypeRef{Spelling: "StoreChecks", Target: store, Emitted: true},
				"the record refers to the names of the checks that the suite of Store emits")
		})

		t.Run("records a conformance test for each directive subject of a struct", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, dualPackage).tests
			reader, found := got["TestBothConformsToReader"]
			assert.True(t, found, "suite records the conformance test of Both and Reader")
			writer, found := got["TestBothConformsToWriter"]
			assert.True(t, found, "suite records the conformance test of Both and Writer")
			expect.Equal(t, []string{reader.Data.Declaration, writer.Data.Declaration}, []string{"Reader", "Writer"},
				"each test applies the suite of the declaration of its directive")
			expect.Empty(t, reader.Data.Excusals, "the subject of Reader leaves the excusal of Writer.Write")
			expect.Equal(t, writer.Data.Excusals, []suite.Excusal{
				{Field: "Write", Method: "AllocationBudget", Reason: "the map allocates on each write"},
			}, "the subject of Writer applies the excusal of Writer.Write")
		})

		t.Run("records the test of a subject beside a subject whose excusal names no check", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, dualPackage).tests
			_, reader := got["TestMixedConformsToReader"]
			expect.False(t, reader, "the subject of Reader excuses a check that its suite does not have")
			writer, found := got["TestMixedConformsToWriter"]
			assert.True(t, found, "suite records the conformance test of Mixed and Writer")
			expect.Equal(t, writer.Data.Excusals, []suite.Excusal{
				{Field: "Write", Method: "AllocationBudget", Reason: "the map allocates on each write"},
			}, "the subject of Writer applies its own excusal")
		})

		t.Run("names each test of a subject after the struct and the declaration", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, dualPackage).tests["TestBothConformsToReader"]
			assert.True(t, found, "suite records the conformance test of Both and Reader")
			expect.Equal(t, []string{got.Data.Allocs, got.Data.Benchmark}, []string{
				"TestBothConformsToReaderAllocs", "BenchmarkBothConformsToReader",
			}, "the names of the allocation test and the benchmark join the struct and the declaration")
			expect.Equal(t, got.Data.Targets, []suite.Target{{Name: "FuzzBothConformsToReaderRead", Field: "Read"}},
				"the fuzz target joins the struct, the declaration and the callable")
		})

		t.Run("reports RepeatedSubject at a second directive subject of one declaration", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, dualPackage)
			stamps, res := phases(t, f, p)
			at := f.Find(t, symbol.KindStruct, "", "Twice").Position()
			var got []diag.Diag
			for _, d := range slices.Concat(slices.Collect(stamps.Sink.All()), slices.Collect(res.Sink.All())) {
				if d.Pos.File == at.File && d.Code == suite.RepeatedSubject {
					got = append(got, d)
				}
			}
			assert.Length(t, got, 1, "suite reports the second directive subject of Twice once")
			expect.Equal(t, got[0].Msg, fmt.Sprintf("Twice has a directive subject of Reader at line %d already, so "+
				"this one has no conformance test", at.Line-2), "the message names the line of the first directive")
			expect.Equal(t, got[0].Pos.Line, at.Line-1, "the finding is at the second directive")
			_, found := generate(t, goTree, dualPackage).tests["TestTwiceConformsToReader"]
			expect.True(t, found, "suite records the conformance test of the first directive")
		})

		{
			tests := []struct {
				name      string
				give      string
				wantMsg   string
				wantTests []string
			}{
				{
					name: "reports UnknownExcusal once for a struct with two directives subject", give: "Unruly",
					wantMsg: "Lister.List/keeps-the-order is no check of the suite of Reader, so Unruly has no " +
						"conformance test. The suite of Reader has no check of Lister.List",
				},
				{
					name: "reports UnknownExcusal at an excusal of a declaration that is not a subject of the struct",
					give: "Stray",
					wantMsg: "Writer.Write/allocation-budget is no check of the suite of Reader, so Stray has no " +
						"conformance test. The suite of Reader has no check of Writer.Write",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, dualPackage)
					_, res := phases(t, f, p)
					subject := f.Find(t, symbol.KindStruct, "", tt.give)
					var got []string
					for d := range res.Sink.All() {
						if d.Code == suite.UnknownExcusal && d.Pos.Line == subject.Position().Line-1 {
							got = append(got, d.Msg)
						}
					}
					expect.Equal(t, got, []string{tt.wantMsg}, "the first directive subject reports the excusal")
					gentest.Unrecorded(t, res, subject)
				})
			}
		}

		t.Run("records the excusal of a ceiling of the budget of the declaration", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, meterPackage, meteredPackage).tests["TestPlainConformsToMeter"]
			assert.True(t, found, "suite records the conformance test of Plain")
			expect.Equal(t, got.Data.Excusals, []suite.Excusal{
				{Field: "Count", Method: "LatencyBudget", Reason: "the meter scans every sample"},
			}, "the excusal names the method that returns the identity of the ceiling")
		})

		t.Run("reports UnknownExcusal with the ceilings of the callable", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, meterPackage, badBudgetPackage)
			_, res := phases(t, f, p)
			subject := f.Find(t, symbol.KindStruct, "", "Unexcused")
			at := subject.Position()
			var got []string
			for d := range res.Sink.All() {
				// The directive excuse of the subject is on the line above its struct.
				if d.Pos.File == at.File && d.Pos.Line == at.Line-1 && d.Code == suite.UnknownExcusal {
					got = append(got, d.Msg)
				}
			}
			expect.Equal(t, got, []string{"Meter.Count/keeps-the-count is no check of the suite of Meter, so " +
				"Unexcused has no conformance test. The checks of Meter.Count are Meter.Count/allocation-budget, " +
				"Meter.Count/byte-budget, Meter.Count/latency-budget and Meter.Count/mean-latency-budget"},
				"the message lists the ceilings of the budget of the callable")
			gentest.Unrecorded(t, res, subject)
		})

		{
			tests := []struct {
				name            string
				givePkgs        []string
				give            string
				wantSuiteAllocs string
				wantBudgets     string
				wantAllocs      string
				wantBenchmark   string
				wantParallel    string
			}{
				{
					name:     "names the falsification of the allocation ceilings of a suite",
					givePkgs: []string{meterPackage}, give: "TestMeterSuite", wantSuiteAllocs: "TestMeterSuiteAllocs",
					wantBudgets: "TestMeterSuiteBudgets", wantBenchmark: "BenchmarkMeterSuite",
					wantParallel: "BenchmarkMeterSuiteParallel",
				},
				{
					name:     "names no test of the allocation ceilings of a suite without one",
					givePkgs: []string{storePackage}, give: "TestStoreSuite", wantBudgets: "TestStoreSuiteBudgets",
					wantBenchmark: "BenchmarkStoreSuite",
				},
				{
					name: "names the tests of the budget of a function", givePkgs: []string{meterPackage},
					give: "TestWeighSuite", wantSuiteAllocs: "TestWeighSuiteAllocs", wantAllocs: "TestWeighAllocs",
					wantBenchmark: "BenchmarkWeigh", wantParallel: "BenchmarkWeighParallel",
				},
				{
					name: "names the benchmark of a function without concurrent calls", givePkgs: []string{slugPackage},
					give: "TestSlugSuite", wantBenchmark: "BenchmarkSlug",
				},
				{
					name:     "names the tests of the budgets of a subject",
					givePkgs: []string{meterPackage, meteredPackage}, give: "TestMeterConformsToMeter",
					wantAllocs: "TestMeterConformsToMeterAllocs", wantBenchmark: "BenchmarkMeterConformsToMeter",
					wantParallel: "BenchmarkMeterConformsToMeterParallel",
				},
				{
					name:     "names the test of the allocation ceilings of a subject with a budget of its own",
					givePkgs: []string{storePackage, memoryPackage}, give: "TestClockedConformsToStore",
					wantAllocs: "TestClockedConformsToStoreAllocs", wantBenchmark: "BenchmarkClockedConformsToStore",
				},
				{
					name:     "names no test of the allocation ceilings of a subject without one",
					givePkgs: []string{storePackage, memoryPackage}, give: "TestStoreConformsToStore",
					wantBenchmark: "BenchmarkStoreConformsToStore",
				},
				{
					name:     "names no test of the allocation ceilings of a subject with a p99 of its own alone",
					givePkgs: []string{pathsPackage, inputsPackage}, give: "TestMemoryConformsToRouter",
					wantBenchmark: "BenchmarkMemoryConformsToRouter",
					wantParallel:  "BenchmarkMemoryConformsToRouterParallel",
				},
				{
					name:     "names the falsification of the checks that count allocations of a suite",
					givePkgs: []string{ledgerPackage}, give: "TestLedgerSuite",
					wantSuiteAllocs: "TestLedgerSuiteAllocs", wantBudgets: "TestLedgerSuiteBudgets",
					wantBenchmark: "BenchmarkLedgerSuite",
				},
				{
					name:     "names the test of the checks that count allocations of a subject",
					givePkgs: []string{ledgerPackage}, give: "TestBookConformsToLedger",
					wantAllocs: "TestBookConformsToLedgerAllocs", wantBenchmark: "BenchmarkBookConformsToLedger",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got, found := generate(t, goTree, tt.givePkgs...).tests[tt.give]
					assert.True(t, found, "suite records the test "+tt.give)
					expect.Equal(t, got.Data.SuiteAllocs, tt.wantSuiteAllocs,
						"the test of a suite falsifies the allocation ceilings of its calls")
					expect.Equal(t, got.Data.Budgets, tt.wantBudgets,
						"the test of the suite of an interface applies its ceilings to the planted defect of a ceiling")
					expect.Equal(t, got.Data.Allocs, tt.wantAllocs,
						"the test applies the allocation ceilings that apply")
					expect.Equal(t, got.Data.Benchmark, tt.wantBenchmark,
						"the test measures each callable of its subject")
					expect.Equal(t, got.Data.Parallel, tt.wantParallel,
						"the test measures each concurrent callable of its subject on goroutines at once")
				})
			}
		}

		t.Run("records a level of the benchmarks for each callable of a subject", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, meterPackage, meteredPackage).tests["TestMeterConformsToMeter"]
			assert.True(t, found, "suite records the conformance test of Meter")
			expect.Equal(t, got.Data.Levels, []suite.Level{
				{Field: "Count", Concurrent: true}, {Field: "Record"}, {Field: "Merge"}, {Field: "Total"},
			}, "each level is a callable, and a concurrent callable is a level of the parallel benchmark as well")
		})

		{
			tests := []struct {
				name    string
				give    string
				wantMsg string
			}{
				{
					name: "reports UnknownExcusal at an excusal of a callable with checks", give: "Priceless",
					wantMsg: "Store.Get/keeps-the-price is no check of the suite of Store, so Priceless has no " +
						"conformance test. The checks of Store.Get are Store.Get/keeps-an-item and " +
						"Store.Get/keeps-the-item",
				},
				{
					name: "reports UnknownExcusal at an excusal of a callable with one check", give: "Misplaced",
					wantMsg: "Store.Put/keeps-the-item is no check of the suite of Store, so Misplaced has no " +
						"conformance test. The check of Store.Put is Store.Put/refuses-a-cancelled-put",
				},
				{
					name: "reports UnknownExcusal at an excusal of a callable without checks", give: "Fetchless",
					wantMsg: "Store.Fetch/finds-the-item is no check of the suite of Store, so Fetchless has no " +
						"conformance test. The suite of Store has no check of Store.Fetch",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := suite.New()
					f := fixture(t, goTree, p, storePackage, badPackage)
					_, res := phases(t, f, p)
					subject := f.Find(t, symbol.KindStruct, "", tt.give)
					at := subject.Position()
					var got []diag.Diag
					for d := range res.Sink.All() {
						// The directive excuse of each subject is on the line above its struct.
						if d.Pos.File == at.File && d.Pos.Line == at.Line-1 {
							got = append(got, d)
						}
					}
					assert.Length(t, got, 1, "suite reports one finding at the excusal of "+tt.give)
					expect.Equal(t, got[0].Code, suite.UnknownExcusal, "the finding has the code UnknownExcusal")
					expect.Equal(t, got[0].Msg, tt.wantMsg, "the message contains the checks of the callable")
					gentest.Unrecorded(t, res, subject)
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

		t.Run("binds a property check of a consumer to the callable of its directive", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, sessionPackage)
			got := recorded(t, f, p).suites["Sessions"]
			expect.Equal(t, got.Checks, []suite.Check{{
				Apply: ref(t, f, "KeepsTheLastPut"), Defect: ref(t, f, "keepsTheFirstPut"), Field: "Get",
				Type: "SessionsGet", Method: "KeepsTheLastPut", Segment: "keeps-the-last-put",
				Statement: "Get returns the session of the last Put of its ID", Property: true,
			}}, "a check that takes a case of the property engine is a property check")
		})

		t.Run("binds a check that counts allocations to the callable of its directive", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, ledgerPackage)
			got := recorded(t, f, p).suites["Ledger"]
			expect.Equal(t, got.Checks, []suite.Check{
				{
					Apply: ref(t, f, "CreditsInPlace"), Defect: ref(t, f, "copiesBalances"), Field: "Credit",
					Type: "LedgerCredit", Method: "CreditsInPlace", Segment: "credits-in-place",
					Statement: "Credit allocates nothing once the account has a balance", Counted: true,
				},
				{
					Apply: ref(t, f, "KeepsTheSum"), Defect: ref(t, f, "dropsCredits"), Field: "Balance",
					Type: "LedgerBalance", Method: "KeepsTheSum", Segment: "keeps-the-sum",
					Statement: "Balance returns the sum of the credits of the account",
				},
			}, "the param counted marks the check that counts allocations alone")
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

		t.Run("binds a check of a method with the segment that the falsification adds", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			f := fixture(t, goTree, p, reservedPackage)
			_, res := phases(t, f, p)
			byPackage := gentest.Records(t, res, suiteTemplate, "",
				func(s suite.Suite) string { return s.Subject.Target.Package })
			checks := byPackage[reservedPackage].Data.Checks
			assert.Length(t, checks, 1, "the suite of Shelf binds its one check")
			expect.Equal(t, checks[0].Segment, "builds-the-subject",
				"the identity of the check names the callable Get, so it differs from the check of the test")
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
					name: "reports BadCheck at a check whose second parameter is not a case", give: "TakesTwo",
					wantCode: suite.BadCheck,
					wantMsg: "the second parameter of the check TakesTwo has the type any, and a property check " +
						"takes a *prop.Case after the handle of its check",
				},
				{
					name: "reports BadCheck at a check of three parameters", give: "TakesThree",
					wantCode: suite.BadCheck,
					wantMsg: "the check TakesThree takes 3 parameters, and a check takes the handle of its check, " +
						"with a *prop.Case after it for a property check",
				},
				{
					name: "reports BadCheck at a check whose second parameter is a list of cases", give: "TakesCases",
					wantCode: suite.BadCheck,
					wantMsg: "the second parameter of the check TakesCases has the type []prop.Case, and a property " +
						"check takes a *prop.Case after the handle of its check",
				},
				{
					name: "reports BadCheck at a check whose second parameter is a case of another package",
					give: "TakesAnotherCase", wantCode: suite.BadCheck,
					wantMsg: "the second parameter of the check TakesAnotherCase has the type *Case, and a property " +
						"check takes a *prop.Case after the handle of its check",
				},
				{
					name: "reports BadCheck at a check whose second parameter is another type of the property engine",
					give: "TakesAnOption", wantCode: suite.BadCheck,
					wantMsg: "the second parameter of the check TakesAnOption has the type *prop.Option, and a " +
						"property check takes a *prop.Case after the handle of its check",
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
				{
					name:     "reports BadCheck at a check of a function with the segment of the check of a factory",
					give:     "BuildsTheSubject",
					wantCode: suite.BadCheck,
					wantMsg: "the test of the suite of Normalize adds the check Normalize/builds-the-subject " +
						"to the forms of its application, so BuildsTheSubject is no check. Rename the function",
				},
				{
					name:     "reports BadCheck at a check of a function with the segment of the check of the form Run",
					give:     "ReceivesTheFunction",
					wantCode: suite.BadCheck,
					wantMsg: "the test of the suite of Normalize adds the check Normalize/receives-the-function " +
						"to the forms of its application, so ReceivesTheFunction is no check. Rename the function",
				},
				{
					name:     "reports BadCheck at a check with the segment of a ceiling of a budget",
					give:     "LatencyBudget",
					wantCode: suite.BadCheck,
					wantMsg: "a budget of Store.Get adds the check Store.Get/latency-budget, so LatencyBudget is no " +
						"check. Rename the function",
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
