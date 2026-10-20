// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder_test

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

	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The tree of the scripted language, and the packages and modules that this
// spec alone reads.
const (
	scriptedTree         = "testdata/scripted"
	aliasedPackage       = "example.com/shop/aliased"
	instancedPackage     = "example.com/shop/instanced"
	parameterizedPackage = "example.com/shop/parameterized"
	badPackage           = "example.com/shop/bad"
	ticketModule         = "ticket"
	badModule            = "bad"
	shopPackage          = "shop"
)

// The declarations under the directive that this spec alone reads, and the
// fields of the findings.
const (
	profileName  = "Profile"
	envelopeName = "Envelope"
	hiddenName   = "Hidden"
	policyName   = "Policy"
	mirrorName   = "Mirror"
	pageName     = "Page"
	pagesName    = "Pages"
	roledName    = "Roled"
	keyedName    = "Keyed"
	derivedName  = "Derived"
	heirName     = "Heir"
	counterName  = "Counter"
	binName      = "Bin"
	ticketName   = "Ticket"
	voucherName  = "Voucher"
	couponName   = "Coupon"
	ledgerName   = "Ledger"
	looseName    = "Loose"
	namedName    = "Named"
	typoName     = "Typo"
	blankName    = "Blank"
	sealedName   = "Sealed"
	shelfName    = "Shelf"
	gaugeName    = "Gauge"
	secondsName  = "Seconds"
	slotsName    = "Slots"
	levelName    = "level"
	tierName     = "Tier"
	deployName   = "Deployment"
	nameField    = "Name"
)

// builder records the model of a builder of each declaration under its
// directive and the checks of each builder, in Go and in TypeScript, and
// passes the plugin suite.
func TestBuilder(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the plugin suite over the Go package catalog", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := builder.New()
				return p, fixture(tb, golang.Target, goTree, p, catalogPackage).Fixture
			})
		})

		t.Run("passes the plugin suite over the TypeScript module catalog", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := builder.New()
				return p, fixture(tb, typescript.Target, tsTree, p, catalogModule).Fixture
			})
		})

		t.Run("tests each builder from outside its package", func(t *testing.T) {
			t.Parallel()

			gentest.External(t, builder.New(), builder.TestTag)
		})

		members := []struct {
			name   string
			target plugin.Target
			tree   string
			decl   string
			want   string
			msg    string
		}{
			{
				name: "records each field of Address that code outside the package assigns", target: golang.Target,
				tree: goTree, decl: addressName, want: "Street City Floor Lines Labels Tags Raw Grid Timeout",
				msg: "the unexported field note is no member of Address",
			},
			{
				name:   "records each field of the class Address that code outside the module assigns",
				target: typescript.Target, tree: tsTree, decl: addressName, want: "street city floor lines labels tags",
				msg: "the static, the readonly, the private and the hard private fields are no members of Address",
			},
			{
				name: "records each embedded field and no field that an embed promotes", target: golang.Target,
				tree: goTree, decl: profileName, want: "Role Audit Active",
				msg: "Profile has the embeds Role and Audit and its field in the order of the source, and none of " +
					"the fields that the embeds promote",
			},
			{
				name: "records an embedded interface and an embedded defined type", target: golang.Target,
				tree: goTree, decl: envelopeName, want: "Notifier Codes",
				msg: "the embed of the interface error, which the workspace does not declare, is no member",
			},
			{
				name: "leaves out a field with the tag that leaves it out", target: golang.Target, tree: goTree,
				decl: hiddenName, want: nameField, msg: "the tag builder of Secret has the value -",
			},
			{
				name: "records the field that the class Derived inherits from its base", target: typescript.Target,
				tree: tsTree, decl: derivedName, want: "active createdBy", msg: "the field of the base is a member",
			},
			{
				name:   "records the fields of the bases of a class without fields in the order of the member walk",
				target: typescript.Target, tree: tsTree, decl: heirName, want: "active createdBy",
				msg: "the field of the nearer base comes first, also where its declaration follows the farther base",
			},
		}
		for _, tt := range members {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				pkg := catalogPackage
				if tt.target == typescript.Target {
					pkg = catalogModule
				}
				got := generate(t, tt.target, tt.tree, pkg)
				assert.Contains(t, got.models, tt.decl, "builder records the model of "+tt.decl)
				assert.Equal(t, fields(got.models[tt.decl].Members), tt.want, tt.msg)
			})
		}

		forms := []struct {
			name      string
			target    plugin.Target
			tree      string
			pkg       string
			decl      string
			give      string
			wantForm  builder.Form
			wantTypes []string
			wantCopy  bool
		}{
			{
				name: "records a field of text as a scalar", target: golang.Target, tree: goTree, pkg: catalogPackage,
				decl: addressName, give: "Street", wantForm: builder.FormScalar,
			},
			{
				name: "records a pointer as an optional of the type that it points at", target: golang.Target,
				tree: goTree, pkg: catalogPackage, decl: addressName, give: "Floor", wantForm: builder.FormOptional,
				wantTypes: []string{"int"},
			},
			{
				name: "records a slice as a list of its element", target: golang.Target, tree: goTree,
				pkg: catalogPackage, decl: addressName, give: "Lines", wantForm: builder.FormList,
				wantTypes: []string{"string"}, wantCopy: true,
			},
			{
				name: "records a map as a map of its key and its value", target: golang.Target, tree: goTree,
				pkg: catalogPackage, decl: addressName, give: "Labels", wantForm: builder.FormMap,
				wantTypes: []string{"string", "string"}, wantCopy: true,
			},
			{
				name: "records the bytes of a slice of bytes", target: golang.Target, tree: goTree, pkg: catalogPackage,
				decl: addressName, give: "Raw", wantForm: builder.FormBytes, wantCopy: true,
			},
			{
				name: "records an array of a fixed length as a scalar", target: golang.Target, tree: goTree,
				pkg: catalogPackage, decl: addressName, give: "Grid", wantForm: builder.FormScalar,
			},
			{
				name: "records a duration as a scalar", target: golang.Target, tree: goTree, pkg: catalogPackage,
				decl: addressName, give: "Timeout", wantForm: builder.FormScalar,
			},
			{
				name: "records a slice through an alias as a list of its element", target: golang.Target, tree: goTree,
				pkg: catalogPackage, decl: orderName, give: "Remarks", wantForm: builder.FormList,
				wantTypes: []string{"string"}, wantCopy: true,
			},
			{
				name: "records a field of a defined type as a scalar", target: golang.Target, tree: goTree,
				pkg: catalogPackage, decl: orderName, give: "Codes", wantForm: builder.FormScalar,
			},
			{
				name: "records a struct embedded by pointer as a scalar", target: golang.Target, tree: goTree,
				pkg: catalogPackage, decl: profileName, give: "Audit", wantForm: builder.FormScalar,
			},
			{
				name: "records a property with a question mark as an optional of its type", target: typescript.Target,
				tree: tsTree, pkg: catalogModule, decl: addressName, give: "floor", wantForm: builder.FormOptional,
				wantTypes: []string{"number"},
			},
			{
				name: "records an array as a list of its element", target: typescript.Target, tree: tsTree,
				pkg: catalogModule, decl: addressName, give: "lines", wantForm: builder.FormList,
				wantTypes: []string{"string"}, wantCopy: true,
			},
			{
				name: "records a Map as a map of its key and its value", target: typescript.Target, tree: tsTree,
				pkg: catalogModule, decl: addressName, give: "labels", wantForm: builder.FormMap,
				wantTypes: []string{"string", "string"}, wantCopy: true,
			},
			{
				name: "records no copy of a ReadonlyArray", target: typescript.Target, tree: tsTree, pkg: catalogModule,
				decl: binName, give: "codes", wantForm: builder.FormList, wantTypes: []string{"string"},
			},
			{
				name: "records a set with a question mark as an optional of the set", target: typescript.Target,
				tree: tsTree, pkg: catalogModule, decl: binName, give: "marks", wantForm: builder.FormOptional,
				wantTypes: []string{"Set"},
			},
			{
				name: "records an array through an alias as a list of its element", target: typescript.Target,
				tree: tsTree, pkg: catalogModule, decl: binName, give: "aliases", wantForm: builder.FormList,
				wantTypes: []string{"string"}, wantCopy: true,
			},
			{
				name: "records a Uint8Array as bytes", target: typescript.Target, tree: tsTree, pkg: catalogModule,
				decl: binName, give: "raw", wantForm: builder.FormBytes, wantCopy: true,
			},
		}
		for _, tt := range forms {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := member(t, generate(t, tt.target, tt.tree, tt.pkg).models[tt.decl], tt.give)
				expect.Equal(t, got.Form, tt.wantForm, "the form of "+tt.give+" follows the shape of its type")
				expect.Equal(t, children(got), tt.wantTypes, "the member records the children of its type")
				expect.Equal(t, got.Copied, tt.wantCopy, "a clone copies a list, a map, a set and bytes alone")
			})
		}

		t.Run("records the type that a struct embeds by pointer", func(t *testing.T) {
			t.Parallel()

			got := member(t, generate(t, golang.Target, goTree, catalogPackage).models[profileName], "Audit")
			assert.NotNil(t, got.Type, "the embedded field has a type")
			expect.Equal(t, got.Type.Form, symbol.FormOptional, "Profile embeds a pointer to Audit")
		})

		t.Run("records a property without a type annotation as a scalar without a type", func(t *testing.T) {
			t.Parallel()

			got := member(t, generate(t, typescript.Target, tsTree, catalogModule).models[counterName], "count")
			expect.Equal(t, got.Form, builder.FormScalar, "a property without a type is a scalar")
			expect.Nil(t, got.Type, "the template of TypeScript spells the type of count from the class")
		})

		t.Run("records nothing for a struct of a language that builder has no part for", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, scriptedTree, shopPackage)
			expect.Empty(t, got.models, "builder records no model of the scripted struct Widget")
			expect.Empty(t, got.findings, "builder reports nothing at Widget")
		})

		t.Run("records the value of a generic builder at its own type parameters", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, genericPackage).models[pageName]
			assert.Length(t, got.Value.Args, 1, "the value of the builder of Page has one type argument")
			expect.Equal(t, got.Value.Args[0].Spelling, "T", "the type argument is the type parameter T")
		})

		seeds := []struct {
			name      string
			pkg       string
			give      string
			wantSeed  builder.Seed
			wantStart string
		}{
			{
				name: "seeds a struct with defaults through its defaults", pkg: catalogPackage, give: addressName,
				wantSeed: builder.SeedDefaults, wantStart: "v.Street v.City v.Timeout",
			},
			{
				name: "assigns no default of a field that an embed promotes", pkg: catalogPackage, give: profileName,
				wantSeed: builder.SeedDefaults, wantStart: "v.Active",
			},
			{
				name: "seeds a struct without a seed or a default with the zero value", pkg: catalogPackage,
				give: orderName, wantSeed: builder.SeedZero, wantStart: "",
			},
			{
				name: "seeds a struct through the function of the param defaults", pkg: seededPackage,
				give: policyName, wantSeed: builder.SeedCompanion, wantStart: "v",
			},
		}
		for _, tt := range seeds {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, golang.Target, goTree, tt.pkg)
				expect.Equal(t, got.models[tt.give].Seed, tt.wantSeed, "the seed of "+tt.give)
				expect.Equal(t, gentest.Assigned(got.starts[tt.give].Body.Stmts), tt.wantStart,
					"the start body declares the seed and assigns each default in the order of the fields")
			})
		}

		t.Run("reads the default of a text field as a reference where it names a declaration", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, seededPackage)
			assert.Empty(t, got.findings, "builder reports nothing in the package seeded")
			start := got.starts[deployName]
			values := make([]string, 0, len(start.Body.Stmts))
			for _, s := range start.Body.Stmts {
				v := s.Value.Val
				value := v.Kind.String() + " " + v.Callee.Name
				if v.Kind == emit.ValueLiteral {
					value = v.Literal.String() + " " + v.Text
				}
				values = append(values, value)
			}
			want := []string{
				"string true", "name Region", "call zone", "name GOOS", "string web", "string a b", "string true",
			}
			assert.Equal(t, values, want,
				"a name and a call in scope and a quoted string are what the rules read, and other text is the string")
		})

		t.Run("seeds through the function that the param defaults refers to", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, seededPackage)
			assert.Equal(t, seedOf(t, got.starts[policyName]), "defaultPolicy",
				"the value is the result of the function of the directive")
		})

		t.Run("assigns the default of a field that a base declares", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, catalogModule)
			expect.Equal(t, got.models[derivedName].Seed, builder.SeedDefaults, "createdBy seeds Derived")
			assert.Equal(t, gentest.Assigned(got.starts[derivedName].Body.Stmts), "v.createdBy",
				"the start body assigns the default of the field of Base")
		})

		t.Run("assigns the default of each field of the struct of an alias", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, aliasedPackage, catalogPackage)
			assert.Equal(t, gentest.Assigned(got.starts[mirrorName].Body.Stmts), "v.Name",
				"the start body of Mirror assigns the default of the field Name of catalog.Role")
		})

		companions := []struct {
			name   string
			target plugin.Target
			tree   string
			pkg    string
			give   string
		}{
			{
				name:   "seeds nothing through a function with the name of the companion that takes a parameter",
				target: golang.Target, tree: goTree, pkg: companionsPackage, give: "Taken",
			},
			{
				name:   "seeds nothing through a function with the name of the companion that has two results",
				target: golang.Target, tree: goTree, pkg: companionsPackage, give: "Paired",
			},
			{
				name:   "seeds nothing through a function with the name of the companion that has no result",
				target: golang.Target, tree: goTree, pkg: companionsPackage, give: "Bare",
			},
			{
				name:   "seeds nothing through a function with the name of the companion that returns another type",
				target: golang.Target, tree: goTree, pkg: companionsPackage, give: "Other",
			},
			{
				name:   "seeds nothing through a function with the name of the companion that has no return type",
				target: typescript.Target, tree: tsTree, pkg: catalogModule, give: counterName,
			},
			{
				name:   "seeds nothing through a function whose name joins the struct and another word",
				target: golang.Target, tree: goTree, pkg: seededPackage, give: tierName,
			},
			{
				name:   "seeds nothing through a function whose name joins the class and another word",
				target: typescript.Target, tree: tsTree, pkg: catalogModule, give: tierName,
			},
		}
		for _, tt := range companions {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, tt.target, tt.tree, tt.pkg)
				assert.Contains(t, got.models, tt.give, "builder records the model of "+tt.give)
				expect.Equal(t, got.models[tt.give].Seed, builder.SeedZero, "no companion seeds "+tt.give)
			})
		}

		checks := []struct {
			name         string
			target       plugin.Target
			tree         string
			pkg          string
			decl         string
			wantSeeds    string
			wantDefaults string
			wantSetters  string
			wantClones   string
		}{
			{
				name:         "records the checks of the constructors, of the setters and of the clones",
				target:       golang.Target,
				tree:         goTree,
				pkg:          catalogPackage,
				decl:         addressName,
				wantSeeds:    "Street City Grid Timeout",
				wantDefaults: "Street City Timeout",
				wantSetters:  "Street City Floor Lines Labels Tags Raw Grid Timeout",
				wantClones:   "Lines Labels Tags Raw",
			},
			{
				name:         "records the checks of the constructors, of the setters and of the clones of the class Address",
				target:       typescript.Target,
				tree:         tsTree,
				pkg:          catalogModule,
				decl:         addressName,
				wantSeeds:    "street city",
				wantDefaults: "street",
				wantSetters:  "street city floor lines labels tags",
				wantClones:   "lines labels tags",
			},
			{
				name:        "records no checks of an embedded field",
				target:      golang.Target,
				tree:        goTree,
				pkg:         catalogPackage,
				decl:        profileName,
				wantSeeds:   "Active",
				wantSetters: "Active",
			},
		}
		for _, tt := range checks {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, tt.target, tt.tree, tt.pkg).checks[tt.decl]
				expect.True(t, got.Instantiable, "the checks of "+tt.decl+" name concrete types")
				expect.Equal(t, strings.Join(got.Seeds, " "), tt.wantSeeds,
					"the check of the constructor from a value compares each scalar field with a sample")
				expect.Equal(t, strings.Join(got.Defaults, " "), tt.wantDefaults,
					"each default other than the zero value has a check")
				expect.Equal(t, checked(got.Setters), tt.wantSetters,
					"each setter whose argument has a sample has a check")
				expect.Equal(t, checked(got.Clones), tt.wantClones,
					"each copied member with two samples has a check of a clone")
			})
		}

		probes := []struct {
			name       string
			target     plugin.Target
			tree       string
			pkg        string
			decl       string
			wantProbes string
			wantFills  string
		}{
			{
				name:   "records a probe of each method of Address that the checks of the setters leave out",
				target: golang.Target, tree: goTree, pkg: catalogPackage, decl: addressName,
				wantProbes: "append AppendLines, entry WithLabelsEntry, entries WithLabelsEntries, " +
					"entry WithTagsEntry, entries WithTagsEntries, text WithRawString",
				wantFills: "Lines Labels Tags Raw",
			},
			{
				name:   "records a probe of each method of the class Address that the checks of the setters leave out",
				target: typescript.Target, tree: tsTree, pkg: catalogModule, decl: addressName,
				wantProbes: "append appendLines, entry withLabelsEntry, entries withLabelsEntries, " +
					"entry withTagsEntry, entries withTagsEntries",
				wantFills: "lines labels tags",
			},
			{
				name:   "records a probe of each method of a member without a sample",
				target: golang.Target, tree: goTree, pkg: catalogPackage, decl: looseName,
				wantProbes: "with WithItems, append AppendItems, with WithIndex, entry WithIndexEntry, " +
					"entries WithIndexEntries, with WithKeys, entry WithKeysEntry, entries WithKeysEntries, " +
					"with WithExtra, with WithHook",
				wantFills: "Items Index Keys",
			},
			{
				name:   "records a probe of each method of a property without a sample",
				target: typescript.Target, tree: tsTree, pkg: catalogModule, decl: looseName,
				wantProbes: "with withItems, append appendItems, with withIndex, entry withIndexEntry, " +
					"entries withIndexEntries, with withKeys, entry withKeysEntry, entries withKeysEntries, " +
					"with withExtra, with withValue",
				wantFills: "items index keys",
			},
			{
				name:   "records no probe of a setter of bytes from text in a language without one",
				target: typescript.Target, tree: tsTree, pkg: catalogModule, decl: binName,
				wantProbes: "append appendCodes, append appendAliases",
				wantFills:  "aliases raw",
			},
			{
				name:   "records a probe of the setter of each embedded field",
				target: golang.Target, tree: goTree, pkg: catalogPackage, decl: profileName,
				wantProbes: "with WithRole, with WithAudit",
			},
			{
				name:   "records a probe of the setter of each member of a struct whose members are embedded fields",
				target: golang.Target, tree: goTree, pkg: catalogPackage, decl: envelopeName,
				wantProbes: "with WithNotifier, with WithCodes",
			},
			{
				name:   "records a probe of the setter of a property without a type",
				target: typescript.Target, tree: tsTree, pkg: catalogModule, decl: counterName,
				wantProbes: "with withCount",
			},
			{
				name:   "records a probe of the setter of a property of a class that shadows a global class",
				target: typescript.Target, tree: tsTree, pkg: shadowModule, decl: shadowedName,
				wantProbes: "with withTags",
			},
		}
		for _, tt := range probes {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, tt.target, tt.tree, tt.pkg)
				assert.Contains(t, got.checks, tt.decl, "builder records the checks of "+tt.decl)
				expect.Equal(t, probed(got.checks[tt.decl].Probes), tt.wantProbes,
					"each method that no check of a setter calls has a probe, in the order of the members")
				expect.Equal(t, checkedFields(got.checks[tt.decl].Fills), tt.wantFills,
					"the check of a clone fills each member that a clone copies")
			})
		}

		t.Run("binds the sample of the setter of a member to each of its probes", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage)
			expect.Equal(t, sampled(got.tests[addressName], got.checks[addressName].Probes),
				[]string{"setterLines", "setterLabels", "setterLabels", "setterTags", "setterTags", "setterRaw"},
				"a probe of a member with a sample binds the sample of the check of its setter")
			expect.Equal(t, sampled(got.tests[addressName], got.checks[addressName].Fills),
				[]string{"setterLines", "setterLabels", "setterTags", "setterRaw"}, "a fill binds the same sample")
			expect.Equal(t, sampled(got.tests[looseName], got.checks[looseName].Probes), []string{},
				"a probe of a member without a sample binds nothing")
		})

		t.Run("records the types of a probe at the witnesses of the type parameters", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, genericPackage).checks[pageName]
			assert.Length(t, got.Probes, 1, "the list Items of Page has the probe of its appender")
			assert.Length(t, got.Probes[0].Type.Elems, 1, "the probe types the list by its element")
			expect.Equal(t, got.Probes[0].Type.Elems[0].Spelling, "int", "the probe types the list at the witness int")
			expect.Equal(t, got.Probes[0].Elem.Spelling, "int", "the probe spells the element at the witness int")
			assert.Length(t, got.Fills, 1, "the check of a clone fills Items")
			expect.Equal(t, got.Fills[0].Elem.Spelling, "int", "the fill spells the element at the witness int")
		})

		t.Run("records the key and the value of a probe of a map", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage).checks[looseName]
			at := slices.IndexFunc(got.Probes, func(p builder.Probe) bool { return p.Method == "WithIndexEntry" })
			assert.NotEqual(t, at, -1, "Loose has the probe of WithIndexEntry")
			expect.Equal(t, got.Probes[at].Key.Spelling, "string", "the probe spells the key of the map")
			expect.Equal(t, got.Probes[at].Value.Spelling, "Notifier", "the probe spells the value of the map")
			expect.Equal(t, got.Probes[at].With, "WithIndex", "the probe knows the setter that replaces the map")
			expect.Equal(t, got.Probes[at].Form, builder.FormMap, "the probe knows the form of the member")
		})

		t.Run("records the setter of each check", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage).checks[addressName]
			expect.Equal(t, setters(got.Setters),
				"WithStreet WithCity WithFloor WithLines WithLabels WithTags WithRaw WithGrid WithTimeout",
				"the check of each setter calls the setter of its field")
			expect.Equal(t, appenders(got.Clones), []string{"AppendLines", "", "", ""},
				"the check of a clone of a list knows its appender, and no other form has one")
		})

		t.Run("seeds the check of the constructor from a value with a sample of each scalar field", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage).tests[addressName]
			assert.Equal(t, gentest.Assigned(seed(t, got).Body.Stmts),
				"seed.Street seed.City seed.Grid seed.Timeout", "the seed assigns each scalar field with a sample")
		})

		t.Run("binds each default to the parameter of its check", func(t *testing.T) {
			t.Parallel()

			test := generate(t, golang.Target, goTree, catalogPackage).tests[addressName]
			defaults := make([]*emit.Method, 0, len(test.Types.Items())-1)
			for _, d := range test.Types.Items()[1:] {
				defaults = append(defaults, d.(*emit.Method))
			}
			assert.Equal(t, wanted(defaults), []string{"string", "string", "time.Duration"},
				"the check of each default of Address takes a value of the type of its field")
			expect.Equal(t, names(defaults), "defaultStreet defaultCity defaultTimeout",
				"the name of the check of a default ends with its field")
			for _, d := range defaults {
				expect.Equal(t, gentest.Assigned(d.Body.Stmts), "want", "the check of a default binds it to want")
			}
		})

		t.Run("binds a sample of the argument of each setter to want", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage)
			methods := bound(got.tests[addressName], got.checks[addressName].Setters, method)
			assert.Equal(
				t,
				wanted(methods),
				[]string{
					"string", "string", "int", "[]string", "map[string]string", "map[string]struct{}", "[]byte",
					"[2]int", "time.Duration",
				},
				"the check of a setter takes the whole list and the type that a pointer points at",
			)
			expect.Equal(
				t,
				names(methods),
				"setterStreet setterCity setterFloor setterLines setterLabels setterTags setterRaw setterGrid setterTimeout",
				"the name of the check of a setter ends with its field",
			)
			for _, m := range methods {
				expect.Equal(t, gentest.Assigned(m.Body.Stmts), "want", "the check of a setter binds a sample to want")
			}
		})

		t.Run("binds the sample and the alternate of each scalar setter to first and want", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage)
			methods := bound(got.tests[addressName], got.checks[addressName].Setters, replacement)
			assert.Equal(t, names(methods), "replaceStreet replaceCity replaceGrid replaceTimeout",
				"each scalar setter has the check of a replacement, and the optional Floor has none")
			for _, m := range methods {
				expect.Equal(t, gentest.Assigned(m.Body.Stmts), "first want",
					"the check of a replacement binds the sample to first and the alternate to want")
			}
		})

		t.Run("binds the sample and the alternate of each copied member to want and other", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage)
			methods := bound(got.tests[addressName], got.checks[addressName].Clones, method)
			assert.Equal(t, names(methods), "cloneLines cloneLabels cloneTags cloneRaw",
				"each copied member has the check of a clone")
			for _, m := range methods {
				expect.Equal(t, gentest.Assigned(m.Body.Stmts), "want other",
					"the check of a clone binds the sample to want and the alternate to other")
			}
		})

		t.Run("records no check of a replacement or of a clone of a member without an alternate", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, catalogPackage)
			checks := got.checks[orderName]
			expect.Equal(t, checked(checks.Setters), "Remarks Codes Flags Marker Marks",
				"each member with a sample has the check of its setter")
			expect.Equal(t, names(bound(got.tests[orderName], checks.Setters, replacement)), "replaceCodes",
				"the scalar Marker has no alternate, so it has no check of a replacement")
			expect.Equal(t, checked(checks.Clones), "Remarks Flags",
				"the sequence Marks has no alternate, so it has no check of a clone")
		})

		t.Run("binds the whole array of an optional to the check of its setter", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, catalogModule)
			methods := bound(got.tests[binName], got.checks[binName].Setters, method)
			assert.Equal(t, wanted(methods), []string{"ReadonlyArray", "Set", "string[]", "Lines", "Uint8Array"},
				"the setters of the optional marks and notes take the whole set and the whole array")
		})

		t.Run("records the checks of a generic struct at the witness of its type parameter", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, genericPackage)
			checks, recorded := got.checks[pageName]
			assert.True(t, recorded, "builder records the checks of Page")
			assert.Length(t, checks.Args, 1, "Page is checked at one witness")
			expect.Equal(t, checks.Args[0].Spelling, "int", "the witness of T is int")
			expect.Equal(t, strings.Join(checks.Defaults, " "), "Next", "the default of Next has a check")
		})

		t.Run("records no checks of a generic struct whose type parameter has no witness", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, genericPackage)
			expect.Contains(t, got.models, keyedName, "builder records the model of Keyed")
			assert.Contains(t, got.checks, keyedName, "builder records the note of Keyed")
			expect.False(t, got.checks[keyedName].Instantiable, "the checks of Keyed name no concrete type")
			expect.Empty(t, got.checks[keyedName].Setters, "Keyed has no checks")
		})

		t.Run("records the builder of an alias with the fields of its struct", func(t *testing.T) {
			t.Parallel()

			got, recorded := generate(t, golang.Target, goTree, aliasedPackage, catalogPackage).models[mirrorName]
			assert.True(t, recorded, "builder records the model of Mirror")
			expect.Equal(t, got.Value.Spelling, mirrorName, "the value of the builder is the alias")
			expect.Equal(t, fields(got.Members), nameField, "the builder sets the fields of catalog.Role")
		})

		t.Run("records the builder of a class whose constructor parameters are all optional", func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, ticketModule)
			assert.Contains(t, got.models, voucherName, "builder records the model of Voucher")
			expect.Equal(t, fields(got.models[voucherName].Members), "code note",
				"the builder sets each field of Voucher")
		})

		t.Run("names each builder after the option suffix", func(t *testing.T) {
			t.Parallel()

			p := builder.New()
			p.(plugin.OptionsProvider).Options().(*builder.Options).Suffix = "factory"
			got := run(t, p, golang.Target, goTree, catalogPackage)
			expect.Equal(t, got.models[addressName].New, "NewAddressFactory",
				"the constructor opens with New and the name of the builder AddressFactory")
			expect.Equal(t, got.checks[addressName].New, "NewAddressFactory", "the checks call the same constructor")
		})

		suffixes := []struct {
			name string
			give string
		}{
			{name: "names each builder after the word builder where the option suffix is empty", give: ""},
			{name: "names each builder after the word builder by default", give: builder.Word},
		}
		for _, tt := range suffixes {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := builder.New()
				options := p.(plugin.OptionsProvider).Options().(*builder.Options)
				expect.Equal(t, options.Suffix, builder.Word, "the option suffix is builder by default")
				options.Suffix = tt.give
				got := run(t, p, golang.Target, goTree, catalogPackage)
				expect.Equal(t, got.models[addressName].New, "NewAddressBuilder", "the builder is AddressBuilder")
			})
		}

		companionWords := []struct {
			name   string
			target plugin.Target
			tree   string
			pkg    string
			want   string
		}{
			{
				name:   "seeds a struct through the function whose name joins the struct and the option companion-word",
				target: golang.Target, tree: goTree, pkg: seededPackage, want: "TierSeed",
			},
			{
				name:   "seeds a class through the function whose name joins the class and the option companion-word",
				target: typescript.Target, tree: tsTree, pkg: catalogModule, want: "tierSeed",
			},
		}
		for _, tt := range companionWords {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := builder.New()
				p.(plugin.OptionsProvider).Options().(*builder.Options).CompanionWord = "seed"
				got := run(t, p, tt.target, tt.tree, tt.pkg)
				expect.Equal(t, got.models[tierName].Seed, builder.SeedCompanion, "the function "+tt.want+" seeds Tier")
				assert.Equal(t, seedOf(t, got.starts[tierName]), tt.want, "the value is the result of "+tt.want)
			})
		}

		t.Run("seeds through AccountDefaults where the option companion-word is empty", func(t *testing.T) {
			t.Parallel()

			p := builder.New()
			options := p.(plugin.OptionsProvider).Options().(*builder.Options)
			expect.Equal(t, options.CompanionWord, builder.CompanionWord,
				"the option companion-word is defaults by default")
			options.CompanionWord = ""
			got := run(t, p, golang.Target, goTree, seededPackage)
			assert.Equal(t, seedOf(t, got.starts[accountName]), "AccountDefaults",
				"the value is the result of AccountDefaults")
		})

		t.Run("records no checks under the option no-tests", func(t *testing.T) {
			t.Parallel()

			p := builder.New()
			p.(plugin.OptionsProvider).Options().(*builder.Options).NoTests = true
			got := run(t, p, golang.Target, goTree, catalogPackage)
			expect.Contains(t, got.models, addressName, "builder records the model of Address")
			expect.Empty(t, got.checks, "builder records no checks")
		})

		t.Run("reports BadTag at a field whose tag has another value, and keeps its setter", func(t *testing.T) {
			t.Parallel()

			p := builder.New()
			f := fixture(t, golang.Target, goTree, p, badPackage)
			res := f.Generate(t, p)
			assert.NoError(t, res.Err, "the generator runs whole")
			got := gentest.Finding(t, res, f.Find(t, symbol.KindField, typoName, nameField))
			expect.Equal(t, got.Code, builder.BadTag, "the finding has the code BadTag")
			expect.Equal(t, got.Msg,
				`the tag builder of Name has the value "skip", and only "-" leaves a field out of the builder`,
				"the message contains the field and the value of its tag")
			subject := func(m builder.Model) string { return m.Subject }
			models := gentest.Records(t, res, builderTemplate, "", subject)
			assert.Contains(t, models, typoName, "builder records the model of Typo")
			expect.Equal(t, fields(models[typoName].Data.Members), nameField, "the field Name keeps its setter")
		})

		findings := []struct {
			name     string
			target   plugin.Target
			tree     string
			pkgs     []string
			wantCode diag.Code
			wantAt   symbol.Kind
			wantHost string
			wantName string
			wantMsg  string
		}{
			{
				name: "reports BadDefault at a field whose default is no value of its type", target: golang.Target,
				tree: goTree, pkgs: []string{badPackage}, wantCode: builder.BadDefault,
				wantAt: symbol.KindField, wantHost: shelfName, wantName: slotsName,
				wantMsg: `the default "many" of Slots is no value of its type int`,
			},
			{
				name: "reports BadDefault at a field of a struct whose default names nothing", target: golang.Target,
				tree: goTree, pkgs: []string{badPackage}, wantCode: builder.BadDefault,
				wantAt: symbol.KindField, wantHost: "Board", wantName: "Origin",
				wantMsg: `the default "origin" of Origin is no value of its type Cell`,
			},
			{
				name:   "reports BadDefault at the default of a property without a type annotation",
				target: typescript.Target, tree: tsTree, pkgs: []string{badModule}, wantCode: builder.BadDefault,
				wantAt: symbol.KindField, wantHost: gaugeName, wantName: levelName,
				wantMsg: `level does not declare a type, so builder cannot check its default "5"`,
			},
			{
				name: "reports NotStruct at an alias of a type that is no struct", target: golang.Target, tree: goTree,
				pkgs: []string{aliasedPackage, catalogPackage}, wantCode: builder.NotStruct,
				wantAt: symbol.KindAlias, wantName: secondsName,
				wantMsg: "Seconds is an alias of int, which is no struct of the workspace",
			},
			{
				name: "reports NotStruct at an alias of an instance of a generic struct", target: golang.Target,
				tree: goTree, pkgs: []string{instancedPackage, genericPackage}, wantCode: builder.NotStruct,
				wantAt: symbol.KindAlias, wantName: pagesName,
				wantMsg: "the alias Pages or its struct Page has type parameters, and a builder sets the fields " +
					"of a struct without them",
			},
			{
				name: "reports NotStruct at a generic alias of a struct without type parameters", target: golang.Target,
				tree: goTree, pkgs: []string{parameterizedPackage, catalogPackage}, wantCode: builder.NotStruct,
				wantAt: symbol.KindAlias, wantName: roledName,
				wantMsg: "the alias Roled or its struct Role has type parameters, and a builder sets the fields " +
					"of a struct without them",
			},
			{
				name: "reports Unconstructed at a class whose constructor takes a parameter", target: typescript.Target,
				tree: tsTree, pkgs: []string{ticketModule}, wantCode: builder.Unconstructed,
				wantAt: symbol.KindStruct, wantName: ticketName,
				wantMsg: "each public constructor of Ticket takes a parameter that a caller must pass, " +
					"and a builder constructs Ticket without arguments",
			},
			{
				name:   "reports Unconstructed at a class whose constructor takes a required parameter after a default",
				target: typescript.Target, tree: tsTree, pkgs: []string{ticketModule}, wantCode: builder.Unconstructed,
				wantAt: symbol.KindStruct, wantName: couponName,
				wantMsg: "each public constructor of Coupon takes a parameter that a caller must pass, " +
					"and a builder constructs Coupon without arguments",
			},
			{
				name: "reports Collision at a struct whose two fields have setters of one name", target: golang.Target,
				tree: goTree, pkgs: []string{badPackage}, wantCode: builder.Collision,
				wantAt: symbol.KindStruct, wantName: ledgerName,
				wantMsg: "the builder of Ledger would declare WithLabelsEntry for Labels and for LabelsEntry, " +
					"so Ledger has no builder",
			},
			{
				name:   "reports Collision at a class whose two properties have setters of one name",
				target: typescript.Target, tree: tsTree, pkgs: []string{badModule}, wantCode: builder.Collision,
				wantAt: symbol.KindStruct, wantName: namedName,
				wantMsg: "the builder of Named would declare withName for name and for Name, so Named has no builder",
			},
			{
				name: "reports Empty at a struct without a field that a builder can set", target: golang.Target,
				tree: goTree, pkgs: []string{badPackage}, wantCode: builder.Empty, wantAt: symbol.KindStruct,
				wantName: blankName, wantMsg: "Blank has no field that a builder can set, so it has no builder",
			},
			{
				name: "reports Empty at a class without a property that a builder can set", target: typescript.Target,
				tree: tsTree, pkgs: []string{badModule}, wantCode: builder.Empty, wantAt: symbol.KindStruct,
				wantName: sealedName, wantMsg: "Sealed has no field that a builder can set, so it has no builder",
			},
		}
		for _, tt := range findings {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := builder.New()
				f := fixture(t, tt.target, tt.tree, p, tt.pkgs...)
				res := f.Generate(t, p)
				assert.NoError(t, res.Err, "the generator runs whole")
				at := f.Find(t, tt.wantAt, tt.wantHost, tt.wantName)
				got := gentest.Finding(t, res, at)
				expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of its case")
				expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the declaration and the reason")
				subject := at
				if tt.wantHost != "" {
					subject = f.Find(t, symbol.KindStruct, "", tt.wantHost)
				}
				gentest.Unrecorded(t, res, subject)
			})
		}
	})
}

// seed returns the function of the record of the checks that seeds each
// scalar member. It stops the test where the record has no seed.
func seed(t *testing.T, test *emit.Struct) *emit.Method {
	t.Helper()

	assert.NotNil(t, test, "builder records the checks")
	assert.NotEmpty(t, test.Types.Items(), "the record of the checks has the seed")
	return test.Types.Items()[0].(*emit.Method)
}

// checked returns the fields of the checks, joined by spaces.
func checked(cs []builder.Check) string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.Field)
	}
	return strings.Join(out, " ")
}

// probed returns the kind and the method of each probe, joined by commas.
func probed(ps []builder.Probe) string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, string(p.Kind)+" "+p.Method)
	}
	return strings.Join(out, ", ")
}

// checkedFields returns the fields of the probes, joined by spaces.
func checkedFields(ps []builder.Probe) string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Field)
	}
	return strings.Join(out, " ")
}

// sampled returns the name of the method of the record test that binds the
// sample of each probe. It skips a probe without a sample.
func sampled(test *emit.Struct, ps []builder.Probe) []string {
	methods := test.Methods.Items()
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		if p.Sample >= 0 {
			out = append(out, methods[p.Sample].Name)
		}
	}
	return out
}

// setters returns the setters that the checks call, joined by spaces.
func setters(cs []builder.Check) string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.Setter)
	}
	return strings.Join(out, " ")
}

// appenders returns the appenders of the checks. A check of a member without
// an appender has the empty name.
func appenders(cs []builder.Check) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.Appender)
	}
	return out
}

// method returns the index of the method of the check c that binds want.
func method(c builder.Check) int { return c.Method }

// replacement returns the index of the method of the check of a replacement
// of the check c, and -1 where c has none.
func replacement(c builder.Check) int { return c.Replace }

// bound returns the methods of the record test at the indexes that at
// returns for the checks, in the order of the checks. It skips a negative
// index, which a check without such a method has.
func bound(test *emit.Struct, cs []builder.Check, at func(builder.Check) int) []*emit.Method {
	methods := test.Methods.Items()
	out := make([]*emit.Method, 0, len(cs))
	for _, c := range cs {
		if i := at(c); i >= 0 {
			out = append(out, methods[i])
		}
	}
	return out
}

// children returns the spellings of the element, the key and the value of a
// member that it records, and nil for a member without them.
func children(m builder.Member) []string {
	var out []string
	for _, ref := range []*emit.TypeRef{m.Elem, m.Key, m.Value} {
		if ref != nil {
			out = append(out, ref.Spelling)
		}
	}
	return out
}

// names returns the names of the methods, joined by spaces.
func names(methods []*emit.Method) string {
	out := make([]string, 0, len(methods))
	for _, m := range methods {
		out = append(out, m.Name)
	}
	return strings.Join(out, " ")
}

// wanted returns the spelling of the type of the first parameter of each
// method.
func wanted(methods []*emit.Method) []string {
	out := make([]string, 0, len(methods))
	for _, m := range methods {
		out = append(out, m.Params.Items()[0].Type.Spelling)
	}
	return out
}
