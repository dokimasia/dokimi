// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder_test

import (
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/eidos/sdk/toolchain"

	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// shopModule is the module of the Go tree.
const shopModule = "example.com/shop"

// The part of Go declares the target golang with the template of each
// family, and writes a map whose values are the empty struct as a set. It
// names the constructors after the builder and declares them as functions of
// the package. It names the methods of each member, and derives the name of
// the companion function of a struct from the struct.
func TestGolang(t *testing.T) {
	t.Parallel()

	t.Run("declares the template of each family for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, builder.New(), golang.Target, builderTemplate, testTemplate)
	})

	t.Run("writes the golden tree golang", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "golang", rendered(t), golden.ShouldUpdate())
	})

	t.Run("writes Go that passes go vet then go test", func(t *testing.T) {
		t.Parallel()

		a := gotesting.New()
		if !toolchain.Require(t, a) {
			return
		}
		project := gentest.Project(rendered(t), goTree, shopModule)
		gotesting.AssertVets(t.Context(), t, project)
		toolchain.AssertTestsPass(t.Context(), t, a, project)
	})

	t.Run("writes a map whose values are the empty struct as a set of its key", func(t *testing.T) {
		t.Parallel()

		got := member(t, generate(t, golang.Target, goTree, catalogPackage).models[addressName], "Tags")
		assert.Equal(t, got.Form, builder.FormSet, "a map to struct{} is the set of Go")
		assert.Equal(t, got.Elem.Spelling, "string", "the set records its key")
		expect.True(t, got.Copied, "a clone copies a set")
	})

	t.Run("names the constructors after the builder", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, catalogPackage)
		expect.Equal(t, got.models[addressName].New, "NewAddressBuilder",
			"the constructor from the seed opens with New and the name of the builder")
		expect.Equal(t, got.models[addressName].From, "NewAddressBuilderFrom",
			"the constructor from a value ends with From")
		expect.Equal(t, got.checks[addressName].New, "NewAddressBuilder", "the checks call the same constructor")
	})

	t.Run("declares a function of each constructor beside the builder", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, catalogPackage)
		assert.Equal(t, got.functions[addressName], []string{"NewAddressBuilder", "NewAddressBuilderFrom"},
			"the record of the builder declares the constructor from the seed and then the one from a value")
	})

	t.Run("refers to each constructor from the record of the checks", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, catalogPackage)
		address := fixture(t, golang.Target, goTree, builder.New(), catalogPackage).
			Find(t, symbol.KindStruct, "", addressName).Identity()
		refs := make([]emit.TypeRef, 0, got.tests[addressName].Fields.Len())
		for _, f := range got.tests[addressName].Fields.Items() {
			refs = append(refs, *f.Type)
		}
		assert.Equal(t, refs, []emit.TypeRef{
			{Spelling: "NewAddressBuilder", Target: address, Emitted: true},
			{Spelling: "NewAddressBuilderFrom", Target: address, Emitted: true},
		}, "the record of the checks refers to the constructor from the seed and then to the other, which the "+
			"record of the builder of Address declares")
	})

	t.Run("writes a map whose empty struct has a space as a set of its key", func(t *testing.T) {
		t.Parallel()

		got := member(t, generate(t, golang.Target, goTree, catalogPackage).models[orderName], "Flags")
		assert.Equal(t, got.Form, builder.FormSet, "the spelling struct {} is the empty struct")
		expect.Equal(t, got.Elem.Spelling, "string", "the set records its key")
	})

	t.Run("writes a map whose values have data as a map", func(t *testing.T) {
		t.Parallel()

		got := member(t, generate(t, golang.Target, goTree, catalogPackage).models[addressName], "Labels")
		assert.Equal(t, got.Form, builder.FormMap, "a map to a string is no set")
	})

	names := []struct {
		name string
		give string
		want builder.Names
	}{
		{name: "names the setter of a scalar after its field", give: "Street", want: builder.Names{With: "WithStreet"}},
		{
			name: "names the methods of a list", give: "Lines",
			want: builder.Names{With: "WithLines", Append: "AppendLines"},
		},
		{
			name: "names the setters of a map", give: "Labels",
			want: builder.Names{With: "WithLabels", Entry: "WithLabelsEntry", Entries: "WithLabelsEntries"},
		},
		{
			name: "names the setters of a set", give: "Tags",
			want: builder.Names{With: "WithTags", Entry: "WithTagsEntry", Entries: "WithTagsEntries"},
		},
		{
			name: "names the setters of bytes", give: "Raw",
			want: builder.Names{With: "WithRaw", Text: "WithRawString"},
		},
	}
	for _, tt := range names {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := member(t, generate(t, golang.Target, goTree, catalogPackage).models[addressName], tt.give)
			assert.Equal(t, got.Names, tt.want, "the methods of "+tt.give+" have the names of Go")
		})
	}

	t.Run("seeds a struct through its companion with the suffix Defaults", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, seededPackage)
		expect.Equal(t, got.models[accountName].Seed, builder.SeedCompanion, "the companion seeds Account")
		assert.Equal(t, seedOf(t, got.starts[accountName]), "AccountDefaults",
			"the value is the result of AccountDefaults")
		expect.Equal(t, gentest.Assigned(got.starts[accountName].Body.Stmts[1:]), "v.Owner",
			"the start body then assigns the default of Owner")
	})
}

// rendered runs builder over the packages of the Go tree whose builders
// build, and returns the files that the Go backend renders from them.
func rendered(t *testing.T) fstest.MapFS {
	t.Helper()

	p := builder.New()
	f := fixture(t, golang.Target, goTree, p, catalogPackage, seededPackage, genericPackage, companionsPackage,
		hiddenPackage)
	_, res := phases(t, f, p)
	return gentest.Rendered(t, f, res, p, gobackend.New())
}
