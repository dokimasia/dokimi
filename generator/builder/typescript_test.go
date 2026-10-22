// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder_test

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

	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The modules of the classes whose field tags has a type named Set: the
// class Frozen, whose Set an import binds, and the class Shadowed, whose Set
// the workspace declares.
const (
	immutableModule = "immutable"
	shadowModule    = "shadow"
	frozenName      = "Frozen"
	shadowedName    = "Shadowed"
)

// The classes of the TypeScript tree whose names contain an acronym:
// HTTPClient, whose companion is httpClientDefaults, MD5Hash, whose companion
// is md5HashDefaults, UUID, whose companion is uuidDefaults, and UserUrl,
// whose companion is userUrlDefaults.
const (
	httpClientName = "HTTPClient"
	md5HashName    = "MD5Hash"
	uuidName       = "UUID"
	userURLName    = "UserUrl"
)

// headerName is the class of the TypeScript tree with properties whose names
// are not identifiers.
const headerName = "Header"

// The part of TypeScript declares the target typescript with the template
// of each family and its helpers, and writes the class Set as a set. It
// names the static constructors and the methods of each member in camel
// case, and declares no function beside a builder class. It leaves out a
// property whose name is not an identifier, and derives the name of the
// companion function of a class in camel case.
func TestTypeScript(t *testing.T) {
	t.Parallel()

	templates, provides := builder.New().(plugin.TemplateProvider)
	assert.True(t, provides, "builder declares templates")

	t.Run("declares the template of each family for the target typescript", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, builder.New(), typescript.Target, builderTemplate, testTemplate)
	})

	t.Run("writes the golden tree typescript", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "typescript", renderedTypeScript(t), golden.ShouldUpdate())
	})

	t.Run("writes the class Set as a set of its type argument", func(t *testing.T) {
		t.Parallel()

		got := member(t, generate(t, typescript.Target, tsTree, catalogModule).models[addressName], "tags")
		assert.Equal(t, got.Form, builder.FormSet, "Set is the set of TypeScript")
		assert.Equal(t, got.Elem.Spelling, "string", "the set records its type argument")
		expect.True(t, got.Copied, "a clone copies a set")
	})

	t.Run("names the two static constructors", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, catalogModule).models[addressName]
		expect.Equal(t, got.New, "create", "AddressBuilder.create constructs a builder from the seed")
		expect.Equal(t, got.From, "from", "AddressBuilder.from constructs a builder from a value")
	})

	t.Run("declares the static constructors as members of the class alone", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, catalogModule)
		expect.Empty(t, got.functions[addressName], "the record of the builder declares no function")
		expect.Empty(t, gentest.Referents(got.tests[addressName]), "the record of the checks refers to no function")
	})

	sets := []struct {
		name string
		pkg  string
		give string
	}{
		{name: "writes a Set that an import binds as a scalar", pkg: immutableModule, give: frozenName},
		{name: "writes a Set that the workspace declares as a scalar", pkg: shadowModule, give: shadowedName},
	}
	for _, tt := range sets {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := member(t, generate(t, typescript.Target, tsTree, tt.pkg).models[tt.give], "tags")
			expect.Equal(t, got.Form, builder.FormScalar, "only the global class Set is the set of TypeScript")
		})
	}

	companions := []struct {
		name string
		give string
		want string
	}{
		{
			name: "seeds a class through the companion in the case of a function",
			give: accountName,
			want: "accountDefaults",
		},
		{name: "lowers the acronym that opens the name of the class", give: httpClientName, want: "httpClientDefaults"},
		{name: "lowers an acronym that a digit follows", give: md5HashName, want: "md5HashDefaults"},
		{name: "lowers a name of capitals alone", give: uuidName, want: "uuidDefaults"},
		{
			name: "keeps the case of a word of the name after the first",
			give: userURLName,
			want: "userUrlDefaults",
		},
	}
	for _, tt := range companions {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := generate(t, typescript.Target, tsTree, catalogModule)
			expect.Equal(t, got.models[tt.give].Seed, builder.SeedCompanion, "the companion seeds "+tt.give)
			assert.Equal(t, seedOf(t, got.starts[tt.give]), tt.want, "the value is the result of "+tt.want)
		})
	}

	names := []struct {
		name string
		decl string
		give string
		want builder.Names
	}{
		{
			name: "names the setter of a scalar after its property", decl: addressName, give: "street",
			want: builder.Names{With: "withStreet"},
		},
		{
			name: "names the methods of a list", decl: addressName, give: "lines",
			want: builder.Names{With: "withLines", Append: "appendLines"},
		},
		{
			name: "names the setters of a map", decl: addressName, give: "labels",
			want: builder.Names{With: "withLabels", Entry: "withLabelsEntry", Entries: "withLabelsEntries"},
		},
		{
			name: "names the setters of a set", decl: addressName, give: "tags",
			want: builder.Names{With: "withTags", Entry: "withTagsEntry", Entries: "withTagsEntries"},
		},
		{
			name: "keeps the case of each word of the property", decl: userURLName, give: "userId",
			want: builder.Names{With: "withUserId"},
		},
	}
	for _, tt := range names {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := member(t, generate(t, typescript.Target, tsTree, catalogModule).models[tt.decl], tt.give)
			assert.Equal(t, got.Names, tt.want, "the methods of "+tt.give+" have the names of TypeScript")
		})
	}

	t.Run("leaves out a property whose name is not an identifier", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, catalogModule).models[headerName]
		assert.Equal(t, fields(got.Members), "v2 $ref",
			"the builder assigns a property after a dot, where content-type, 0 and [Symbol.toStringTag] are no names")
	})

	funcs := templates.TemplateFuncs(typescript.Target)

	maps := []struct {
		name string
		give *emit.TypeRef
		want bool
	}{
		{name: "reports a Map as an instance of the class", give: &emit.TypeRef{Spelling: "Map"}, want: true},
		{
			name: "reports a Map of the library of TypeScript as an instance of the class",
			give: &emit.TypeRef{
				Spelling: "Map",
				Target: symbol.Identity{
					Lang: typescript.Lang, Package: plugin.StorePath(typescript.LibraryStore, "lib.es2022.full.d.ts"),
					Name: "Map",
				},
			},
			want: true,
		},
		{
			name: "reports a ReadonlyMap as an instance of the class",
			give: &emit.TypeRef{Spelling: "ReadonlyMap"},
			want: true,
		},
		{name: "reports a Record as an object", give: &emit.TypeRef{Spelling: "Record"}},
		{
			name: "reports a declaration of the workspace named Map as no global class",
			give: &emit.TypeRef{
				Spelling: "Map",
				Target:   symbol.Identity{Lang: typescript.Lang, Package: "geo", Name: "Map"},
			},
		},
		{
			name: "reports a Map that an import binds as no global class",
			give: &emit.TypeRef{Spelling: "Map", Package: "immutable"},
		},
	}
	for _, tt := range maps {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mapClass, declared := funcs["mapclass"].(func(*emit.TypeRef) bool)
			assert.True(t, declared, "the target typescript declares the helper mapclass")
			assert.Equal(t, mapClass(tt.give), tt.want, "mapclass tells the class Map from an object")
		})
	}

	t.Run("trims the suffix of the name of a spec", func(t *testing.T) {
		t.Parallel()

		trimSuffix, declared := funcs["trimSuffix"].(func(string, string) string)
		assert.True(t, declared, "the target typescript declares the helper trimSuffix")
		expect.Equal(t, trimSuffix("catalog.builder.test.ts", ".test.ts"), "catalog.builder",
			"the test file imports the module of its builders")
	})
}

// renderedTypeScript runs builder over the modules of the TypeScript tree
// whose builders build, and returns the files that the TypeScript backend
// renders from them.
func renderedTypeScript(t *testing.T) fstest.MapFS {
	t.Helper()

	p := builder.New()
	_, res := phases(t, fixture(t, typescript.Target, tsTree, p, catalogModule, immutableModule, shadowModule,
		aliasedModule), p)
	return gentest.Rendered(t, res, p, tsbackend.New())
}
