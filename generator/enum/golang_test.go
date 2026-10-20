// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum_test

import (
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"
	"go.dokimi.dev/eidos/sdk/toolchain"

	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The part of Go declares the target golang with the template of each
// family, and names the members of the surface of an enumeration. It prints a
// value outside the variants through the type under the enumeration, and
// quotes the text of each variant. The checks refer to the functions and the
// variable of the surface through its package.
func TestGolang(t *testing.T) {
	t.Parallel()

	t.Run("declares the template of each family for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, enum.New(), golang.Target, enumTemplate, testTemplate)
	})

	t.Run("writes the surfaces and the checks of the golden tree golang", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "golang", rendered(t), golden.ShouldUpdate())
	})

	t.Run("writes Go that vets and whose tests pass", func(t *testing.T) {
		t.Parallel()

		a := gotesting.New()
		if !toolchain.Require(t, a) {
			return
		}
		project := gentest.Project(rendered(t), goTree, "example.com/shop")
		gotesting.AssertVets(t.Context(), t, project)
		toolchain.AssertTestsPass(t.Context(), t, a, project)
	})

	t.Run("names the methods, the functions and the error of the surface", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage).models[statusName].Names
		assert.Equal(t, got, enum.Names{
			Render: "String", Parse: "ParseStatus", Values: "StatusValues", Valid: "IsValid",
			Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnknownStatus",
		}, "the surface of Status has the names of Go")
	})

	fallbacks := []struct {
		name       string
		pkg        string
		give       string
		wantType   string
		wantFormat string
	}{
		{
			name: "prints a number through its type with the verb d", pkg: statusPackage, give: statusName,
			wantType: "int", wantFormat: "%d",
		},
		{
			name: "prints a float through its type with the verb g", pkg: quietPackage, give: levelName,
			wantType: "float64", wantFormat: "%g",
		},
		{
			name: "prints a truth value through its type with the verb t", pkg: quietPackage, give: "Flag",
			wantType: "bool", wantFormat: "%t",
		},
	}
	for _, tt := range fallbacks {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, tt.pkg).models[tt.give]
			assert.NotNil(t, got.Fallback, tt.give+" converts a value outside the variants")
			expect.Equal(t, got.Fallback.Spelling, tt.wantType, "the value converts to the type under "+tt.give)
			expect.Equal(t, got.Format, tt.wantFormat, "the verb prints a value of that type")
		})
	}

	t.Run("prints a value of text outside the variants without a conversion", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, regionPackage).models[regionName]
		expect.Nil(t, got.Fallback, "Region renders a value outside the variants as its text")
	})

	t.Run("writes a variant as the constant of its package", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage).models[statusName]
		assert.NotEmpty(t, got.Variants, "Status has variants")
		expect.Equal(t, got.Variants[0].Member, "Draft", "a template writes Draft alone")
	})

	t.Run("declares the package members of the surface beside the enumeration", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage)
		assert.Equal(t, got.beside[statusName], []string{"ParseStatus", "StatusValues", "ErrUnknownStatus"},
			"the record of the surface declares the parse function, the list of the variants and the error")
	})

	t.Run("refers to each package member of the surface from the record of the checks", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, statusPackage)
		assert.Equal(t, gentest.Referents(got.tests[statusName]),
			[]string{"ParseStatus", "StatusValues", "ErrUnknownStatus"},
			"the fields of the record of the checks refer to the declarations of the surface in their order")
	})

	t.Run("declares nothing beside an enumeration with its own parse function", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, handwrittenPackage)
		expect.Empty(t, got.beside[weekdayName], "the surface of Weekday has no parse function of its own")
	})

	t.Run("refers to nothing from the checks of an enumeration with its own parse function", func(t *testing.T) {
		t.Parallel()

		got := generate(t, golang.Target, goTree, handwrittenPackage)
		assert.NotNil(t, got.tests[weekdayName], "Weekday has checks")
		expect.Empty(t, gentest.Referents(got.tests[weekdayName]), "the checks of Weekday call no parse function")
	})
}

// rendered runs enum over the packages of the Go tree whose enumerations
// have surfaces, and returns the files that the Go backend renders from
// them.
func rendered(t *testing.T) fstest.MapFS {
	t.Helper()

	p := enum.New()
	res := fixture(t, golang.Target, goTree, p, statusPackage, "example.com/shop/priority", regionPackage,
		"example.com/shop/split", "example.com/shop/handwritten", quietPackage, commentedPackage).Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	return gentest.Rendered(t, res, p, gobackend.New())
}
