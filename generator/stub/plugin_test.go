// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees, packages and modules that this spec alone reads, and the
// declarations under the directive that it alone reads.
const (
	scriptedTree      = "testdata/scripted"
	remotePackage     = "example.com/shop/remote"
	listenerPackage   = "example.com/shop/listener"
	emptyPackage      = "example.com/shop/empty"
	clashPackage      = "example.com/shop/clash"
	otherPackage      = "example.com/shop/other"
	clientPackage     = "example.com/shop/client"
	unfinishedPackage = "example.com/shop/unfinished"
	partitionPackage  = "example.com/shop/partition"
	helpersPackage    = "example.com/shop/helpers"
	suitedPackage     = "example.com/shop/suited"
	wrapPackage       = "example.com/shop/wrap"
	toolPackage       = "example.com/shop/lib/tool"
	deepPackage       = "example.com/shop/lib/internal/deep"
	shopPackage       = "shop"
	bagName           = "Bag"
	pairsName         = "Pairs"
	sessionName       = "Session"
	cursorName        = "Cursor"
	pointName         = "Point"
	boxName           = "Box"
	tableName         = "Table"
)

// stub records the model of a double of each declaration under its directive
// and of each handle that the declaration produces, in Go and in TypeScript,
// and passes the plugin suite.
func TestStub(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the plugin suite over the Go package catalog", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := stub.New()
				return p, fixture(tb, golang.Target, goTree, p, catalogPackage).Fixture
			})
		})

		t.Run("passes the plugin suite over the TypeScript module store", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := stub.New()
				return p, fixture(tb, typescript.Target, tsTree, p, storeModule).Fixture
			})
		})

		t.Run("runs after no capability", func(t *testing.T) {
			t.Parallel()

			ordered, orders := stub.New().(plugin.CapabilityProvider)
			assert.True(t, orders, "stub declares its capabilities")
			expect.Empty(t, ordered.After(), "stub runs in any order with the suite generator")
		})

		t.Run("records a double of a declaration under the directive suite alone only as a handle", func(t *testing.T) {
			t.Parallel()

			p := stub.New()
			f := gentest.New(t, &gentest.Spec{
				Dir: goTree, Packages: []string{suitedPackage}, Target: golang.Target, Plugins: []plugin.Plugin{p},
				Annotators: append(annotators(), suite.New()),
			})
			res := f.Generate(t, p)
			assert.NoError(t, res.Err, "the generator runs whole")
			models := gentest.List[stub.Model](res)
			subjects := make([]string, 0, len(models))
			for _, m := range models {
				subjects = append(subjects, m.Subject)
			}
			slices.Sort(subjects)
			expect.Equal(t, subjects, []string{"Journal", "Ledger", "Register"},
				"a declaration under stub has one double, a declaration under suite alone has one where a double "+
					"returns it, and the struct Cashbox under suite alone has none")
		})

		t.Run("leaves the import of a method to the layout where the plan previews no package", func(t *testing.T) {
			t.Parallel()

			p := stub.New()
			f := gentest.New(t, &gentest.Spec{
				Dir: goTree, Packages: []string{wrapPackage, toolPackage, deepPackage}, Target: golang.Target,
				Plugins: []plugin.Plugin{p}, Annotators: annotators(),
			})
			res := f.Generate(t, p)
			assert.NoError(t, res.Err, "the generator runs whole")
			subjects := make([]string, 0, 4)
			for _, m := range gentest.List[stub.Model](res) {
				subjects = append(subjects, m.Subject)
			}
			expect.Contains(t, subjects, "Wrapped", "stub records the double of Wrapped, and the layout refuses it")
		})

		t.Run("checks each double from outside its package", func(t *testing.T) {
			t.Parallel()

			gentest.External(t, stub.New(), stub.TestTag)
		})

		{
			builtPkgs := []string{builtPackage, billingPackage}
			builtModules := []string{builtModule, billingModule}
			goAccount := symbol.Identity{
				Lang: golang.Lang, Package: builtPackage, Name: "Account", Kind: symbol.KindStruct,
			}
			tsAccount := symbol.Identity{
				Lang: typescript.Lang, Package: builtModule, Name: "Account", Kind: symbol.KindStruct,
			}
			account := &stub.Default{Ref: "NewAccountBuilder", Target: goAccount}
			accountClass := &stub.Default{Ref: "AccountBuilder", Seed: "create", Target: tsAccount}
			tests := []struct {
				name          string
				giveTarget    plugin.Target
				giveTree      string
				givePkgs      []string
				give          string
				want          []*stub.Default
				wantDefaulted bool
			}{
				{
					name:       "records the default result of a value of a struct with a builder",
					giveTarget: golang.Target, giveTree: goTree, givePkgs: builtPkgs, give: "Get",
					want: []*stub.Default{account, nil}, wantDefaulted: true,
				},
				{
					name:       "records the default result of a pointer to a struct with a builder",
					giveTarget: golang.Target, giveTree: goTree, givePkgs: builtPkgs, give: "Find",
					want:          []*stub.Default{{Ref: "NewAccountBuilder", Target: goAccount, Optional: true}, nil},
					wantDefaulted: true,
				},
				{
					name:       "records no default result of a flag beside a struct with a builder",
					giveTarget: golang.Target, giveTree: goTree, givePkgs: builtPkgs, give: "Lookup",
					want: []*stub.Default{account, nil}, wantDefaulted: true,
				},
				{
					name:       "records no default result of a struct without a builder",
					giveTarget: golang.Target, giveTree: goTree, givePkgs: builtPkgs, give: "Note",
					want: []*stub.Default{nil},
				},
				{
					name:       "records no default result of a struct whose builder is in another package",
					giveTarget: golang.Target, giveTree: goTree, givePkgs: builtPkgs, give: "Bill",
					want: []*stub.Default{nil},
				},
				{
					name:       "records the default result of a promise of a class with a builder",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: builtModules, give: "get",
					want: []*stub.Default{accountClass}, wantDefaulted: true,
				},
				{
					name:       "records the default result of an optional class with a builder",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: builtModules, give: "find",
					want: []*stub.Default{
						{Ref: "AccountBuilder", Seed: "create", Target: tsAccount, Optional: true},
					},
					wantDefaulted: true,
				},
				{
					name:       "records no default result of a class without a builder",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: builtModules, give: "note",
					want: []*stub.Default{nil},
				},
				{
					name:       "records no default result of a class whose builder is in another module",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: builtModules, give: "bill",
					want: []*stub.Default{nil},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, tt.giveTarget, tt.giveTree, tt.givePkgs...).models[storeName], tt.give)
					defaults := make([]*stub.Default, 0, len(got.Returns))
					for _, r := range got.Returns {
						defaults = append(defaults, r.Default)
					}
					expect.Equal(t, defaults, tt.want, "each result of "+tt.give+" has the default result of its case")
					expect.Equal(t, got.Defaulted, tt.wantDefaulted,
						"the method has a default result where one of its results has one")
				})
			}
		}

		t.Run("records the type arguments of the default result of a generic struct", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, builtPackage, billingPackage).models[storeName], "List")
			assert.NotNil(t, got.Returns[0].Default, "the page of accounts has a default result")
			expect.Equal(t, got.Returns[0].Default.Ref, "NewPageBuilder",
				"the default result calls the constructor of the builder of Page")
			assert.Length(t, got.Returns[0].Default.Args, 1, "the constructor takes one type argument")
			expect.Equal(t, got.Returns[0].Default.Args[0].Spelling, "Account", "the type argument is Account")
		})

		{
			tests := []struct {
				name       string
				giveTarget plugin.Target
				giveTree   string
				givePkgs   []string
				want       []string
			}{
				{
					name:       "refers to each constructor of a default result from each generated file",
					giveTarget: golang.Target, giveTree: goTree, givePkgs: []string{builtPackage, billingPackage},
					want: []string{"NewAccountBuilder", "NewPageBuilder"},
				},
				{
					name:       "refers to each builder class of a default result from each generated file",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{builtModule, billingModule},
					want: []string{"AccountBuilder"},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, tt.giveTarget, tt.giveTree, tt.givePkgs...)
					expect.Equal(t, gentest.Referents(got.doubles[storeName]), tt.want,
						"the record of the double refers to each declaration that a default result calls, once")
					refs := gentest.Referents(got.tests[storeName])
					assert.True(t, len(refs) >= len(tt.want), "the record of the checks refers to the declarations")
					expect.Equal(t, refs[len(refs)-len(tt.want):], tt.want,
						"the record of the checks refers to the same declarations after the declarations of the double")
				})
			}
		}

		{
			tests := []struct {
				name        string
				giveTarget  plugin.Target
				giveTree    string
				givePkg     string
				giveSubject string
				want        string
				wantFrom    string
			}{
				{
					name:       "records a method of each method in the member walk of an interface",
					giveTarget: golang.Target, giveTree: goTree, givePkg: catalogPackage, giveSubject: storeName,
					want:     "Store.Get Store.Put Store.Len Store.Lookup Store.Log Store.Split Store.Close",
					wantFrom: "Close Closer",
				},
				{
					name:       "records a method of each method in the member walk of a TypeScript interface",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, giveSubject: storeName,
					want: "Store.get Store.put Store.size Store.log Store.close", wantFrom: "close Closer",
				},
				{
					name:       "records a method of each method in the method set of a struct",
					giveTarget: golang.Target, giveTree: goTree, givePkg: catalogPackage, giveSubject: clockName,
					want: "Clock.Now Clock.Zone Clock.Every",
				},
				{
					name:       "records a method of each public instance method of a class",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: clockModule, giveSubject: clockName,
					want: "Clock.now Clock.every Clock.stopper",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[tt.giveSubject]
					expect.Equal(t, qualified(got.Methods), tt.want,
						"the double has a method of each method of the walk")
					expect.Equal(t, contributed(got.Methods), tt.wantFrom,
						"a method that a contributor declares names the contributor")
				})
			}
		}

		t.Run("records no field of a Go struct as a property of its double", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage).models[clockName]
			expect.Empty(t, got.Properties, "a double of a Go struct declares no field of the struct")
		})

		t.Run("records each property of a class as a field of its double", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, clockModule).models[clockName]
			expect.Equal(t, properties(got.Properties), "zone", "the double declares the property zone")
		})

		t.Run("records a double of an interface with properties alone", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, pointModule).models[pointName]
			expect.Equal(t, properties(got.Properties), "x y", "the double of Point declares its properties")
		})

		t.Run("records each property of an interface as a field of its double", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, storeModule).models[storeName]
			assert.Equal(t, properties(got.Properties), "zone label", "the double declares each property of Store")
			expect.Equal(t, []bool{got.Properties[0].Optional, got.Properties[1].Optional}, []bool{false, true},
				"the property label may be absent")
			expect.Equal(t, got.Properties[0].Type.Spelling, "string", "the field has the type of its property")
		})

		{
			tests := []struct {
				name        string
				giveTarget  plugin.Target
				giveTree    string
				givePkg     string
				giveSubject string
				want        string
			}{
				{
					name: "records an interface as the value of its double", giveTarget: golang.Target,
					giveTree: goTree, givePkg: catalogPackage, giveSubject: storeName, want: storeName,
				},
				{
					name: "records a class as the value of its double", giveTarget: typescript.Target,
					giveTree: tsTree, givePkg: clockModule, giveSubject: clockName, want: clockName,
				},
				{
					name: "records no value of a Go struct as its double", giveTarget: golang.Target,
					giveTree: goTree, givePkg: catalogPackage, giveSubject: clockName,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[tt.giveSubject]
					expect.Equal(t, spellingOf(got.Value), tt.want,
						"the double is a value of its subject where the language admits it")
				})
			}
		}

		{
			tests := []struct {
				name        string
				givePkg     string
				giveSubject string
				want        string
			}{
				{
					name:    "records an interface as the type of the compile-time check of its double",
					givePkg: catalogPackage, giveSubject: storeName, want: storeName,
				},
				{
					name: "records no compile-time check of a double of a Go struct", givePkg: catalogPackage,
					giveSubject: clockName,
				},
				{
					name: "records the compile-time check of a generic double with witnesses", givePkg: genericPackage,
					giveSubject: repoName, want: repoName,
				},
				{
					name:    "records no compile-time check of a generic double without witnesses",
					givePkg: genericPackage, giveSubject: bagName,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, golang.Target, goTree, tt.givePkg).models[tt.giveSubject]
					expect.Equal(t, spellingOf(got.Guard), tt.want, "the check asserts a value of the subject")
				})
			}
		}

		{
			tests := []struct {
				name        string
				giveTarget  plugin.Target
				giveTree    string
				givePkg     string
				giveSubject string
				want        string
			}{
				{
					name: "records an interface as the type that DelegateTo wraps", giveTarget: golang.Target,
					giveTree: goTree, givePkg: catalogPackage, giveSubject: storeName, want: storeName,
				},
				{
					name:       "records no type that DelegateTo wraps for a Go struct",
					giveTarget: golang.Target, giveTree: goTree, givePkg: catalogPackage, giveSubject: clockName,
				},
				{
					name:       "records a TypeScript interface as the type that DelegateTo wraps",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, giveSubject: storeName,
					want: storeName,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[tt.giveSubject]
					expect.Equal(t, spellingOf(got.Delegate), tt.want,
						"DelegateTo wraps an implementation of the subject")
				})
			}
		}

		{
			tests := []struct {
				name        string
				givePkg     string
				giveSubject string
				want        bool
			}{
				{
					name: "records a subject without type parameters as not generic", givePkg: catalogPackage,
					giveSubject: storeName,
				},
				{
					name: "records a subject with a type parameter as generic", givePkg: genericPackage,
					giveSubject: repoName, want: true,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, golang.Target, goTree, tt.givePkg).models[tt.giveSubject]
					expect.Equal(t, got.Generic, tt.want, "the double is generic where its subject is")
				})
			}
		}

		t.Run("records the names that join the subject to the method", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], "Get")
			expect.Equal(t, []string{got.Qualified, got.Call, got.Stub, got.Return, got.Field},
				[]string{"Store.Get", "StoreGetCall", "StoreGetStub", "StoreGetReturn", "OnGet"},
				"the names join the subject, the method and a word")
		})

		t.Run("sets the option suffix to the word stub by default", func(t *testing.T) {
			t.Parallel()

			options := stub.New().(plugin.OptionsProvider).Options().(*stub.Options)
			expect.Equal(t, options.Suffix, stub.Word, "the option suffix is stub by default")
		})

		t.Run("names each double after the option suffix", func(t *testing.T) {
			t.Parallel()

			p := stub.New()
			p.(plugin.OptionsProvider).Options().(*stub.Options).Suffix = "double"
			got := run(t, p, golang.Target, goTree, catalogPackage)
			expect.Equal(t, got.doubles[storeName].Name, "StoreDouble", "the double joins the suffix")
			expect.Equal(t, got.models[storeName].Double, "StoreDouble", "the model names the double")
			expect.Equal(t, method(t, got.models[storeName], "Get").Stub, "StoreGetDouble",
				"the type that configures Get joins the suffix")
		})

		t.Run("names each double after the word stub where the option suffix is empty", func(t *testing.T) {
			t.Parallel()

			p := stub.New()
			p.(plugin.OptionsProvider).Options().(*stub.Options).Suffix = ""
			got := run(t, p, golang.Target, goTree, catalogPackage)
			expect.Equal(t, got.doubles[storeName].Name, "StoreStub", "the double joins the word stub")
		})

		{
			tests := []struct {
				name       string
				giveTarget plugin.Target
				giveTree   string
				givePkg    string
				want       string
			}{
				{
					name:       "records each Go method that returns an error as a method that can fail",
					giveTarget: golang.Target, giveTree: goTree, givePkg: catalogPackage,
					want: "Store.Get Store.Put Store.Split Store.Close",
				},
				{
					name:       "records each TypeScript method as a method that can fail",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule,
					want: "Store.get Store.put Store.size Store.log Store.close",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[storeName]
					expect.Equal(t, failing(got.Methods), tt.want, "a fault can fail the method")
				})
			}
		}

		{
			tests := []struct {
				name        string
				giveSubject string
				want        bool
			}{
				{name: "records a double with a method that can fail as fallible", giveSubject: storeName, want: true},
				{name: "records a double without a method that can fail as infallible", giveSubject: clockName},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, golang.Target, goTree, catalogPackage).models[tt.giveSubject]
					expect.Equal(t, got.Fallible, tt.want, "a double is fallible where a fault can fail a method")
				})
			}
		}

		t.Run("records the methods that answer asynchronously", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, storeModule).models[storeName]
			var async []string
			for _, m := range got.Methods {
				if m.Async {
					async = append(async, m.Name)
				}
			}
			expect.Equal(t, async, []string{"get", "close"}, "the methods that return a promise are async")
		})

		{
			tests := []struct {
				name        string
				givePkg     string
				giveSubject string
				give        string
				want        string
			}{
				{
					name: "names each parameter without a name by its position", givePkg: catalogPackage,
					giveSubject: storeName, give: "Put", want: "a0 a1",
				},
				{
					name: "names the parameter _ by its position", givePkg: edgePackage, giveSubject: edgeName,
					give: "Skip", want: "a0 e",
				},
				{
					name: "names the parameter _ apart from another parameter", givePkg: edgePackage,
					giveSubject: edgeName, give: "Shift", want: "a00 a0",
				},
				{
					name: "names the parameter _ apart from a result", givePkg: edgePackage, giveSubject: edgeName,
					give: "Fetch", want: "a00",
				},
				{
					name: "names two parameters _ apart from each other", givePkg: edgePackage, giveSubject: edgeName,
					give: "Wide", want: "a a10 a1 b c d e f g h a100",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, golang.Target, goTree, tt.givePkg).models[tt.giveSubject], tt.give)
					expect.Equal(t, params(got.Params), tt.want, "each parameter has a name of its own")
				})
			}
		}

		{
			tests := []struct {
				name       string
				giveTarget plugin.Target
				giveTree   string
				givePkg    string
				give       string
			}{
				{
					name:       "records no variadic Go parameter as optional",
					giveTarget: golang.Target, giveTree: goTree, givePkg: catalogPackage, give: "Log",
				},
				{
					name:       "records no variadic TypeScript parameter as optional",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, give: "log",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[storeName], tt.give)
					assert.Length(t, got.Params, 2, "the method has a format and a variadic parameter")
					expect.False(t, got.Params[1].Optional, "a call always has the list of the variadic arguments")
				})
			}
		}

		t.Run("marks the parameter that the rules of Go give the role of a context", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], "Get")
			assert.Length(t, got.Params, 2, "Get takes a context and an identity")
			expect.Equal(t, []bool{got.Params[0].Context, got.Params[1].Context}, []bool{true, false},
				"the context of Get is the one parameter with the role of a context")
		})

		{
			tests := []struct {
				name         string
				give         string
				wantReceiver string
				wantRecord   string
				wantAnswer   string
			}{
				{
					name: "names the receiver of a method s", give: "Read", wantReceiver: "s", wantRecord: "call",
					wantAnswer: "r",
				},
				{
					name: "names the receiver of a method with a parameter s another name", give: "Use",
					wantReceiver: "s0", wantRecord: "call", wantAnswer: "r",
				},
				{
					name: "names the local of the call of a method with a parameter call another name",
					give: "Record", wantReceiver: "s", wantRecord: "call0", wantAnswer: "r",
				},
				{
					name: "names the local of the answer of a method with a parameter r another name", give: "Keep",
					wantReceiver: "s", wantRecord: "call", wantAnswer: "r1",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, golang.Target, goTree, edgePackage).models[edgeName], tt.give)
					expect.Equal(t, []string{got.Receiver, got.Record, got.Answer},
						[]string{tt.wantReceiver, tt.wantRecord, tt.wantAnswer}, "no parameter has the name of a local")
				})
			}
		}

		{
			tests := []struct {
				name string
				give string
				want string
			}{
				{name: "names the local of the function of a method fn", give: "now", want: "fn"},
				{
					name: "names the local of the function of a method with a parameter fn another name",
					give: "every", want: "fn0",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, typescript.Target, tsTree, clockModule).models[clockName], tt.give)
					expect.Equal(t, got.Func, tt.want, "no parameter has the name of the local of the function")
				})
			}
		}

		{
			tests := []struct {
				name        string
				givePkg     string
				giveSubject string
				give        string
				want        bool
				wantLocals  string
			}{
				{
					name: "declares the names of results that each have a name", givePkg: catalogPackage,
					giveSubject: storeName, give: "Lookup", want: true, wantLocals: "item ok",
				},
				{
					name: "declares no names of results without names", givePkg: catalogPackage,
					giveSubject: storeName, give: "Split", wantLocals: "r0 r1 r2",
				},
				{
					name: "declares no names of results with a result named _", givePkg: edgePackage,
					giveSubject: edgeName, give: "Read", wantLocals: "r0 r1",
				},
				{
					name: "declares no names of results with a result named like the receiver", givePkg: edgePackage,
					giveSubject: edgeName, give: "Clash", wantLocals: "r0 r1",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, golang.Target, goTree, tt.givePkg)
					m := method(t, got.models[tt.giveSubject], tt.give)
					expect.Equal(t, m.Named, tt.want, "the double declares the names where Go compiles them")
					expect.Equal(t, resultLocals(m.Returns), tt.wantLocals, "the locals of the results")
					double := declared(t, got.doubles[tt.giveSubject], tt.give)
					expect.Equal(t, declaredResults(double), names(m.Returns),
						"the method of the double declares the names of the model")
				})
			}
		}

		{
			tests := []struct {
				name        string
				givePkg     string
				giveSubject string
				give        string
				want        string
			}{
				{
					name:    "names each parameter of the setter of the results after its field in camel case",
					givePkg: catalogPackage, giveSubject: storeName, give: "Lookup", want: "item ok",
				},
				{
					name:    "names a parameter of the setter of the results apart from its receiver",
					givePkg: edgePackage, giveSubject: edgeName, give: "Clash", want: "s0 err",
				},
				{
					name:    "names two parameters of the setter of the results apart from each other",
					givePkg: edgePackage, giveSubject: edgeName, give: "Pair", want: "ab ab0",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, golang.Target, goTree, tt.givePkg).models[tt.giveSubject], tt.give)
					expect.Equal(t, setterParams(got.Returns), tt.want,
						"each parameter of the setter has a name of its own")
				})
			}
		}

		{
			tests := []struct {
				name string
				give string
				want int
			}{
				{name: "records the position of the error result", give: "Split", want: 2},
				{name: "records no position of an error of a method that cannot fail", give: "Len", want: -1},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], tt.give)
					expect.Equal(t, got.Err, tt.want, "the position of the error among the results")
				})
			}
		}

		t.Run("marks the error result of a method as its error", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], "Split")
			expect.Equal(t, []bool{got.Returns[0].Error, got.Returns[1].Error, got.Returns[2].Error},
				[]bool{false, false, true}, "the third result of Split is its error")
		})

		{
			tests := []struct {
				name       string
				giveTarget plugin.Target
				giveTree   string
				givePkg    string
				give       string
				want       string
			}{
				{
					name:       "records the field of the failure of a call where the language throws its errors",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, give: "get", want: "error",
				},
				{
					name:       "records no field of the failure of a call where the language returns its errors",
					giveTarget: golang.Target, giveTree: goTree, givePkg: catalogPackage, give: "Get",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[storeName], tt.give)
					expect.Equal(t, got.Failure, tt.want, "the call records a thrown error in its field of the failure")
				})
			}
		}

		{
			tests := []struct {
				name     string
				give     string
				wantArgs string
				wantPins string
			}{
				{
					name: "takes no sample of the context of Get", give: "Get", wantArgs: "a1",
					wantPins: "want0",
				},
				{
					name: "samples the key of Lookup with its results", give: "Lookup", wantArgs: "a0",
					wantPins: "want0 want1",
				},
				{name: "leaves the variadic parameter of Log at its zero value", give: "Log", wantArgs: "a0"},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, golang.Target, goTree, catalogPackage)
					m := method(t, got.models[storeName], tt.give)
					expect.Equal(t, assigned(got.tests[storeName], m.Args), tt.wantArgs, "the body of the arguments")
					expect.Equal(t, assigned(got.tests[storeName], m.Pins), tt.wantPins, "the body of the results")
				})
			}
		}

		t.Run("records no samples of a method whose values have none", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], "Close")
			expect.Equal(t, []int{got.Args, got.Pins}, []int{-1, -1}, "Close has no argument and no sampled result")
			expect.False(t, got.Pinnable, "no result of Close has a sample")
		})

		t.Run("marks each result with a sample as pinned", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], "Get")
			expect.Equal(t, []bool{got.Returns[0].Pinned, got.Returns[1].Pinned}, []bool{true, false},
				"the item of Get has a sample and its error has none")
			expect.True(t, got.Pinnable, "a result of Get has a sample")
		})

		{
			tests := []struct {
				name        string
				giveSubject string
				want        string
			}{
				{
					name: "records the witness of a type parameter that the rules derive", giveSubject: repoName,
					want: "int",
				},
				{
					name: "records the authored witness of each type parameter", giveSubject: pairsName,
					want: "string int",
				},
				{name: "records no witnesses of a type parameter without a witness", giveSubject: bagName},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, golang.Target, goTree, genericPackage).models[tt.giveSubject]
					expect.Equal(t, spellings(got.Witnesses), tt.want, "the witnesses of the type parameters")
				})
			}
		}

		t.Run("records the compile-time check of a generic double at its witnesses", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, genericPackage).models[repoName]
			assert.NotNil(t, got.Guard, "a compile-time check asserts the double of Repo")
			expect.Equal(t, spellings(got.Guard.Args), "int", "the check asserts Repo at the witness")
		})

		t.Run("records the spec of a generic double without witnesses", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, genericPackage)
			expect.Contains(t, got.tests, bagName, "the tests of Bag state that the double has no checks")
		})

		{
			tests := []struct {
				name        string
				giveSubject string
				want        string
			}{
				{
					name:        "records the samples of the properties of a double without methods",
					giveSubject: pointName,
					want:        "double.x double.y",
				},
				{name: "records no samples of a property whose type mentions a type parameter", giveSubject: boxName},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := generate(t, typescript.Target, tsTree, pointModule)
					assert.Contains(t, got.tests, tt.giveSubject, "a double without methods has tests")
					expect.Equal(t, assigned(got.tests[tt.giveSubject], got.models[tt.giveSubject].Props), tt.want,
						"the body of the properties assigns each sample to the property of the double")
				})
			}
		}

		{
			tests := []struct {
				name           string
				giveTarget     plugin.Target
				giveTree       string
				givePkg        string
				giveSubject    string
				give           string
				wantTypeParams string
				wantArgs       string
				wantWitnesses  string
			}{
				{
					name:       "records the type parameters that the signature of a method mentions",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, giveSubject: tableName,
					give: "get", wantTypeParams: "K V", wantArgs: "K", wantWitnesses: "unknown",
				},
				{
					name:       "records no type parameters of a method whose signature mentions none",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, giveSubject: tableName,
					give: "size",
				},
				{
					name:       "records the type parameters that the result of a method mentions alone",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, giveSubject: repoName,
					give: "find", wantTypeParams: "T",
				},
				{
					name:       "records the witnesses of the type parameters that the parameters mention",
					giveTarget: golang.Target, giveTree: goTree, givePkg: genericPackage, giveSubject: repoName,
					give: "Put", wantTypeParams: "T", wantArgs: "T", wantWitnesses: "int",
				},
				{
					name:       "records no witnesses of a method of a double without witnesses",
					giveTarget: golang.Target, giveTree: goTree, givePkg: genericPackage, giveSubject: bagName,
					give: "Add", wantTypeParams: "K", wantArgs: "K",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					subject := generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[tt.giveSubject]
					got := method(t, subject, tt.give)
					expect.Equal(t, typeParams(got.TypeParams), tt.wantTypeParams,
						"the types of the method declare the type parameters that its signature mentions")
					expect.Equal(t, typeParams(got.ArgTypeParams), tt.wantArgs,
						"the samples of the arguments declare the type parameters that the parameters mention")
					expect.Equal(t, spellings(got.ArgWitnesses), tt.wantWitnesses,
						"the samples of the arguments take the witnesses of their type parameters")
				})
			}
		}

		t.Run("records no samples of an argument whose type mentions a type parameter", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, genericPackage).models[repoName], "Put")
			expect.Equal(t, got.Args, -1, "the authored sample of v is no value of T inside the generic checks")
		})

		t.Run("records no samples of a result whose type mentions a type parameter", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, genericPackage).models[repoName], "Get")
			expect.Equal(t, got.Pins, -1, "the result T of Get has no sample inside the generic checks")
		})

		t.Run("records a double of a handle that a method returns", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, handlesPackage).models
			cursor, recorded := got[cursorName]
			assert.True(t, recorded, "stub records the double of the handle Cursor")
			expect.Equal(t, cursor.Producer, "Producer", "the first producer of Cursor writes its double")
			expect.Equal(t, qualified(cursor.Methods), "Cursor.Next Cursor.Close", "the handle has its methods")
		})

		t.Run("records a double of a handle that a method of an embedded interface returns", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, handlesPackage).models
			log, recorded := got["Log"]
			assert.True(t, recorded, "stub records the double of the handle Log")
			expect.Equal(t, log.Producer, "Archive", "Archive returns Log through the interface Opener")
		})

		t.Run("records a double of a handle under the directive without a producer", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, handlesPackage).models
			expect.Equal(t, got["Subscription"].Producer, "",
				"the directive of Subscription writes its double, and no producer does")
		})

		t.Run("records no double of a handle without methods", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, handlesPackage).models
			expect.NotContains(t, got, "Empty", "an interface without methods has no double")
		})

		t.Run("reports nothing at a producer of a handle without methods", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, handlesPackage)
			expect.Empty(t, got.findings, "Nothing returns an interface without methods, which is no defect")
		})

		t.Run("records one double of a handle that several methods return", func(t *testing.T) {
			t.Parallel()

			p := stub.New()
			res := fixture(t, golang.Target, goTree, p, handlesPackage).Generate(t, p)
			assert.NoError(t, res.Err, "the generator runs whole")
			records := gentest.Records(t, res, stubTemplate, "", func(m stub.Model) string { return m.Subject })
			subjects := make([]string, 0, len(records))
			for _, r := range records {
				subjects = append(subjects, r.Data.Subject)
			}
			expect.Equal(t, count(subjects, cursorName), 1, "Open, Resume and Read return Cursor, which has one double")
		})

		t.Run("records no double of a handle of another package", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, handlesPackage, otherPackage).models
			expect.NotContains(t, got, "Thing", "the package other has the handle Thing")
		})

		t.Run("records the double of a handle at a producer in the package of the handle", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, handlesPackage, clientPackage).models
			expect.Equal(t, got[cursorName].Producer, "Producer", "Client returns Cursor from another package")
		})

		t.Run("records a double of a handle of a TypeScript interface", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, handlesModule).models[cursorName]
			expect.Equal(t, got.Producer, "Producer", "open returns a promise of the handle Cursor")
		})

		{
			tests := []struct {
				name       string
				giveTarget plugin.Target
				giveTree   string
				givePkg    string
				give       string
				want       bool
			}{
				{
					name: "records the mixin deprecated of a method", giveTarget: golang.Target, giveTree: goTree,
					givePkg: mixedPackage, give: "Old", want: true,
				},
				{
					name: "records no mixin deprecated of a method without it", giveTarget: golang.Target,
					giveTree: goTree, givePkg: mixedPackage, give: "Commit",
				},
				{
					name: "records the mixin deprecated of a TypeScript method", giveTarget: typescript.Target,
					giveTree: tsTree, givePkg: mixedModule, give: "old", want: true,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[sessionName], tt.give)
					expect.Equal(t, got.Deprecated, tt.want, "the method is deprecated where it has the mixin")
				})
			}
		}

		{
			tests := []struct {
				name       string
				giveTarget plugin.Target
				giveTree   string
				givePkg    string
				give       string
				want       string
			}{
				{
					name: "records the method that the mixin order-after names", giveTarget: golang.Target,
					giveTree: goTree, givePkg: mixedPackage, give: "Commit", want: "Prepare",
				},
				{
					name:       "records no method that a method without the mixin order-after follows",
					giveTarget: golang.Target, giveTree: goTree, givePkg: mixedPackage, give: "Prepare",
				},
				{
					name:       "records the method that the mixin order-after of a TypeScript method names",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: mixedModule, give: "commit",
					want: "prepare",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[sessionName], tt.give)
					expect.Equal(t, got.After, tt.want, "the method follows the method that its mixin names")
				})
			}
		}

		t.Run("records a double with a method that follows another as ordered", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, mixedPackage).models[sessionName]
			expect.True(t, got.Ordered, "Commit follows Prepare")
		})

		{
			tests := []struct {
				name          string
				giveTarget    plugin.Target
				giveTree      string
				givePkg       string
				give          string
				wantErrors    []string
				wantRetry     int
				wantPartition int
			}{
				{
					name:       "records the attempt of retry of the directive fault on a Go method",
					giveTarget: golang.Target, giveTree: goTree, givePkg: parityPackage, give: "List",
					wantRetry: 3, wantPartition: -1,
				},
				{
					name:       "records the helpers of the error values with the partition of a Go method",
					giveTarget: golang.Target, giveTree: goTree, givePkg: parityPackage, give: "Get",
					wantErrors: []string{
						"FaultNotFound arms example.com/shop/parity.ErrNotFound",
						"FaultGone arms example.com/shop/parity.ErrGone",
					},
					wantPartition: 1,
				},
				{
					name:       "records each param of the directive fault on a TypeScript method",
					giveTarget: typescript.Target, giveTree: tsTree, givePkg: storeModule, give: "get",
					wantErrors: []string{"faultNotFound arms store.notFound"}, wantRetry: 3,
				},
				{
					name:       "records no fault of a method without the directive fault",
					giveTarget: golang.Target, giveTree: goTree, givePkg: parityPackage, give: "Close",
					wantPartition: -1,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got := method(t, generate(t, tt.giveTarget, tt.giveTree, tt.givePkg).models[storeName], tt.give)
					var armed []string
					for _, e := range got.Errors {
						expect.Equal(t, e.Value.Kind, emit.ValueName, "the value of "+e.Name+" names its declaration")
						armed = append(armed, e.Helper+" arms "+e.Value.Callee.Package+"."+e.Value.Callee.Name)
					}
					expect.Equal(t, armed, tt.wantErrors, "the helper of each error value arms that value")
					expect.Equal(t, got.Retry, tt.wantRetry, "the attempt that the retry schedule lets succeed")
					expect.Equal(t, got.Partition, tt.wantPartition, "the position of the partition parameter")
				})
			}
		}

		{
			tests := []struct {
				name    string
				give    string
				wantMsg string
			}{
				{
					name: "reports Incomparable at a partition on a variadic parameter", give: "Ledger",
					wantMsg: "the partition parameter args of Log is variadic, so a partition fault cannot compare " +
						"its arguments with a key",
				},
				{
					name: "reports Incomparable at a partition whose type mentions a type parameter", give: "Box",
					wantMsg: "the partition parameter v of Put has a type that mentions a type parameter, so a " +
						"partition fault cannot compare its arguments with a key",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := stub.New()
					f := fixture(t, golang.Target, goTree, p, partitionPackage)
					res := f.Generate(t, p)
					assert.NoError(t, res.Err, "the generator runs whole")
					got := gentest.Finding(t, res, f.Find(t, symbol.KindInterface, "", tt.give))
					expect.Equal(t, got.Code, stub.Incomparable, "the finding has the code Incomparable")
					expect.Equal(t, got.Msg, tt.wantMsg, "the message names the parameter and the reason")
					models := generate(t, golang.Target, goTree, partitionPackage).models
					assert.Contains(t, models, tt.give, "stub records the double without the partition helpers")
					expect.Equal(t, models[tt.give].Methods[0].Partition, -1, "the method has no partition")
				})
			}
		}

		{
			tests := []struct {
				name       string
				give       string
				wantElem   string
				wantSecond string
				wantErrors bool
			}{
				{name: "records a sequence of values", give: "Names", wantElem: "string"},
				{
					name: "records a sequence of values with an error", give: "Rows", wantElem: "Item",
					wantSecond: "error", wantErrors: true,
				},
				{name: "records a sequence of pairs", give: "Counts", wantElem: "string", wantSecond: "int"},
				{name: "records a sequence through an alias of its type", give: "Tail", wantElem: "string"},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					table := generate(t, golang.Target, goTree, sequencePackage).models["Table"]
					got := method(t, table, tt.give).Sequence
					assert.NotNil(t, got, tt.give+" returns a sequence")
					expect.Equal(t, got.Elem.Spelling, tt.wantElem, "the element of the sequence")
					expect.Equal(t, spellingOf(got.Second), tt.wantSecond, "the second value of a pair")
					expect.Equal(t, got.Errors, tt.wantErrors, "the second value is an error")
				})
			}
		}

		t.Run("records a sequence of a TypeScript async iterable", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, streamModule).models["Feed"]
			rows := method(t, got, "rows").Sequence
			assert.NotNil(t, rows, "rows returns a sequence")
			expect.Equal(t, rows.Elem.Spelling, "number", "the element of rows is a number")
		})

		t.Run("records no sequence of a method that returns none", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, catalogPackage).models[storeName], "Get")
			expect.Nil(t, got.Sequence, "Get returns an item and an error")
		})

		t.Run("records no sequence of a stream without a type argument", func(t *testing.T) {
			t.Parallel()

			got := method(t, generate(t, golang.Target, goTree, unfinishedPackage).models["Feed"], "Names")
			expect.Nil(t, got.Sequence, "iter.Seq without a type argument has no element")
		})

		t.Run("records nothing for a struct of a language that stub has no part for", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, scriptedTree, shopPackage)
			expect.Empty(t, got.models, "stub records no double of the scripted struct Widget")
			expect.Empty(t, got.findings, "stub reports nothing at Widget")
		})

		{
			tests := []struct {
				name       string
				giveTarget plugin.Target
				giveTree   string
				givePkgs   []string
				wantCode   diag.Code
				wantAt     symbol.Kind
				wantName   string
				wantMsg    string
			}{
				{
					name: "reports Constraint at an interface with type-set elements", giveTarget: golang.Target,
					giveTree: goTree, givePkgs: []string{genericPackage}, wantCode: stub.Constraint,
					wantAt: symbol.KindInterface, wantName: "Keyable",
					wantMsg: "Keyable has type-set elements and only constrains type parameters, so it has no double",
				},
				{
					name:       "reports IncompleteMembers at an interface that embeds an interface of no loaded package",
					giveTarget: golang.Target,
					giveTree:   goTree,
					givePkgs:   []string{remotePackage},
					wantCode:   stub.IncompleteMembers,
					wantAt:     symbol.KindInterface,
					wantName:   "Remote",
					wantMsg:    "the members of Remote are incomplete, so it has no double",
				},
				{
					name:       "reports IncompleteMembers at an interface with an overloaded method",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"overloaded"},
					wantCode: stub.IncompleteMembers, wantAt: symbol.KindInterface, wantName: "Overloaded",
					wantMsg: "the members of Overloaded are incomplete, so it has no double",
				},
				{
					name:       "reports Collision at an interface whose double has two members of one name",
					giveTarget: golang.Target, giveTree: goTree, givePkgs: []string{listenerPackage},
					wantCode: stub.Collision, wantAt: symbol.KindInterface, wantName: "Listener",
					wantMsg: "the double of Listener would have two members named OnClose, so it has no double",
				},
				{
					name:       "reports Collision at an interface with a property of the name of the field of a method",
					giveTarget: typescript.Target,
					giveTree:   tsTree,
					givePkgs:   []string{"toggle"},
					wantCode:   stub.Collision,
					wantAt:     symbol.KindInterface,
					wantName:   "Toggle",
					wantMsg:    "the double of Toggle would have two members named onRun, so it has no double",
				},
				{
					name:       "reports Collision at an interface whose call has a parameter with the field of its error",
					giveTarget: golang.Target,
					giveTree:   goTree,
					givePkgs:   []string{clashPackage},
					wantCode:   stub.Collision,
					wantAt:     symbol.KindInterface,
					wantName:   "Clash",
					wantMsg:    "the double of Clash would have two members named ClashFailCall.Err, so it has no double",
				},
				{
					name:       "reports Collision at an interface whose stub of a method has two helpers of one name",
					giveTarget: golang.Target,
					giveTree:   goTree,
					givePkgs:   []string{helpersPackage},
					wantCode:   stub.Collision,
					wantAt:     symbol.KindInterface,
					wantName:   "Fetcher",
					wantMsg: "the double of Fetcher would have two members named FetcherFetchStub.FaultGone, " +
						"so it has no double",
				},
				{
					name:       "reports Collision at an interface whose call has a parameter of the name of its result",
					giveTarget: typescript.Target,
					giveTree:   tsTree,
					givePkgs:   []string{"clashed"},
					wantCode:   stub.Collision,
					wantAt:     symbol.KindInterface,
					wantName:   "Clashed",
					wantMsg: "the double of Clashed would have two members named ClashedSendCall.result, " +
						"so it has no double",
				},
				{
					name:       "reports Collision at an interface whose call has a parameter of the name of its failure",
					giveTarget: typescript.Target,
					giveTree:   tsTree,
					givePkgs:   []string{"clashed"},
					wantCode:   stub.Collision,
					wantAt:     symbol.KindInterface,
					wantName:   "Failed",
					wantMsg:    "the double of Failed would have two members named FailedFailCall.error, so it has no double",
				},
				{
					name: "reports Empty at an interface without methods", giveTarget: golang.Target,
					giveTree: goTree, givePkgs: []string{emptyPackage}, wantCode: stub.Empty,
					wantAt: symbol.KindInterface, wantName: "Nothing",
					wantMsg: "Nothing has no member that a double can implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at an interface with an index signature",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"indexed"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Indexed",
					wantMsg: "Indexed has an index signature, which a double cannot implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at an interface with a method with type parameters",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"picked"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Picked",
					wantMsg: "Picked has the method pick with type parameters, which a double cannot implement, " +
						"so it has no double",
				},
				{
					name:       "reports Unimplementable at an interface with an accessor",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"accessed"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Accessed",
					wantMsg: "Accessed has the accessor size, which a double cannot implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at an interface with a method of a computed name",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"computed"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Computed",
					wantMsg: "Computed has the method [Symbol.iterator] with a computed name, which a double " +
						"cannot implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at an interface with a property of a computed name",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"computed"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Tagged",
					wantMsg: "Tagged has the property [Symbol.toStringTag] with a computed name, which a double " +
						"cannot implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at an interface with a construct signature",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"constructed"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Constructed",
					wantMsg: "Constructed has a construct signature, which a double cannot implement, " +
						"so it has no double",
				},
				{
					name:       "reports Unimplementable at a class with a private member",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"sealed"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindStruct, wantName: "Sealed",
					wantMsg: "Sealed has the member count that is not public, which a double cannot implement, " +
						"so it has no double",
				},
				{
					name:       "reports Unimplementable at a class with a protected method",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"guarded"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindStruct, wantName: "Guarded",
					wantMsg: "Guarded has the member check that is not public, which a double cannot implement, " +
						"so it has no double",
				},
				{
					name:       "reports Unimplementable at a class with a method without a declared result type",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"undeclared"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindStruct, wantName: "Undeclared",
					wantMsg: "Undeclared has the method now without a declared result type, which a double cannot " +
						"implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at a class with a property without a declared type",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"undeclared"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindStruct, wantName: "Counted",
					wantMsg: "Counted has the property count without a declared type, which a double cannot " +
						"implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at an interface with a parameter without a declared type",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"undeclared"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Loose",
					wantMsg: "Loose has the parameter x of run without a declared type, which a double cannot " +
						"implement, so it has no double",
				},
				{
					name:       "reports Unimplementable at a struct with a method that refers to an internal package",
					giveTarget: golang.Target,
					giveTree:   goTree,
					givePkgs:   []string{wrapPackage, toolPackage, deepPackage},
					wantCode:   stub.Unimplementable,
					wantAt:     symbol.KindStruct,
					wantName:   "Wrapped",
					wantMsg: "Wrapped has the method Level, whose signature refers to deep.Level, which the file of " +
						"its double in the package example.com/shop/wrap cannot import, so it has no double: golang: " +
						"example.com/shop/lib/internal/deep is internal to example.com/shop/lib, and the package " +
						"example.com/shop/wrap is outside it",
				},
				{
					name:       "reports Unimplementable at an interface with a method that refers to an internal package",
					giveTarget: golang.Target,
					giveTree:   goTree,
					givePkgs:   []string{wrapPackage, toolPackage, deepPackage},
					wantCode:   stub.Unimplementable,
					wantAt:     symbol.KindInterface,
					wantName:   "Measured",
					wantMsg: "Measured has the method Level, whose signature refers to deep.Level, which the file of " +
						"its double in the package example.com/shop/wrap cannot import, so it has no double: golang: " +
						"example.com/shop/lib/internal/deep is internal to example.com/shop/lib, and the package " +
						"example.com/shop/wrap is outside it",
				},
				{
					name:       "reports Unimplementable at a struct with a parameter that refers to an internal package",
					giveTarget: golang.Target,
					giveTree:   goTree,
					givePkgs:   []string{wrapPackage, toolPackage, deepPackage},
					wantCode:   stub.Unimplementable,
					wantAt:     symbol.KindStruct,
					wantName:   "Adjusted",
					wantMsg: "Adjusted has the method Set, whose signature refers to deep.Level, which the file of " +
						"its double in the package example.com/shop/wrap cannot import, so it has no double: golang: " +
						"example.com/shop/lib/internal/deep is internal to example.com/shop/lib, and the package " +
						"example.com/shop/wrap is outside it",
				},
				{
					name:       "reports Unimplementable at an interface that is not exported",
					giveTarget: typescript.Target, giveTree: tsTree, givePkgs: []string{"hidden"},
					wantCode: stub.Unimplementable, wantAt: symbol.KindInterface, wantName: "Hidden",
					wantMsg: "Hidden is not exported, so the file of its double cannot refer to it",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := stub.New()
					f := fixture(t, tt.giveTarget, tt.giveTree, p, tt.givePkgs...)
					res := f.Generate(t, p)
					assert.NoError(t, res.Err, "the generator runs whole")
					at := f.Find(t, tt.wantAt, "", tt.wantName)
					got := gentest.Finding(t, res, at)
					expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of its case")
					expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the declaration and the reason")
					gentest.Unrecorded(t, res, at)
				})
			}
		}
	})
}

// qualified returns the qualified names of the methods, joined by spaces.
func qualified(ms []stub.Method) string {
	out := make([]string, 0, len(ms))
	for i := range ms {
		out = append(out, ms[i].Qualified)
	}
	return strings.Join(out, " ")
}

// contributed returns the name of each method that a contributor declares
// and the name of its contributor, joined by spaces.
func contributed(ms []stub.Method) string {
	var out []string
	for i := range ms {
		if ms[i].From != "" {
			out = append(out, ms[i].Name, ms[i].From)
		}
	}
	return strings.Join(out, " ")
}

// failing returns the qualified names of the methods that a fault can fail,
// joined by spaces.
func failing(ms []stub.Method) string {
	var out []string
	for i := range ms {
		if ms[i].Fails {
			out = append(out, ms[i].Qualified)
		}
	}
	return strings.Join(out, " ")
}

// properties returns the fields of the properties, joined by spaces.
func properties(vs []stub.Value) string {
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		out = append(out, v.Field)
	}
	return strings.Join(out, " ")
}

// params returns the names of the parameters, joined by spaces.
func params(ps []stub.Param) string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Name)
	}
	return strings.Join(out, " ")
}

// resultLocals returns the locals of the results, joined by spaces.
func resultLocals(rs []stub.Result) string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.Local)
	}
	return strings.Join(out, " ")
}

// setterParams returns the parameters of the setter of the results, joined
// by spaces.
func setterParams(rs []stub.Result) string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.Param)
	}
	return strings.Join(out, " ")
}

// names returns the declared names of the results, joined by spaces.
func names(rs []stub.Result) string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.Name)
	}
	return strings.Join(out, " ")
}

// assigned returns the names that the body of the method at the position i
// of the record test assigns, joined by spaces, and the empty string for the
// position -1.
func assigned(test *emit.Struct, i int) string {
	if i < 0 {
		return ""
	}
	return gentest.Assigned(test.Methods.Items()[i].Body.Stmts)
}

// spellings returns the spelling of each reference, joined by spaces.
func spellings(refs []*emit.TypeRef) string {
	out := make([]string, 0, len(refs))
	for _, r := range refs {
		out = append(out, r.Spelling)
	}
	return strings.Join(out, " ")
}

// typeParams returns the name of each type parameter, joined by spaces.
func typeParams(tps []*emit.TypeParam) string {
	out := make([]string, 0, len(tps))
	for _, tp := range tps {
		out = append(out, tp.Name)
	}
	return strings.Join(out, " ")
}

// spellingOf returns the spelling of the reference ref, and the empty string
// for no reference.
func spellingOf(ref *emit.TypeRef) string {
	if ref == nil {
		return ""
	}
	return ref.Spelling
}

// declaredResults returns the names of the results that the method m
// declares, joined by spaces.
func declaredResults(m *emit.Method) string {
	out := make([]string, 0, m.Returns.Len())
	for _, r := range m.Returns.Items() {
		out = append(out, r.Name)
	}
	return strings.Join(out, " ")
}

// count returns how many of names equal name.
func count(names []string, name string) int {
	n := 0
	for _, got := range names {
		if got == name {
			n++
		}
	}
	return n
}
