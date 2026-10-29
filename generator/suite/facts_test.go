// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
)

// derivedOf returns the identities of the derived checks of the suite s.
func derivedOf(s *derive.Suite) []string {
	out := make([]string, 0, len(s.Checks))
	for _, c := range s.Checks {
		out = append(out, c.Identity)
	}
	return out
}

// factsOf returns the facts of the callable named callable in the suite s.
func factsOf(t *testing.T, s *suite.Suite, callable string) *derive.Facts {
	t.Helper()

	at := slices.IndexFunc(s.Facts, func(f derive.Facts) bool { return f.Callable == callable })
	assert.True(t, at >= 0, "the suite has the facts of "+callable)
	return &s.Facts[at]
}

// inputNamed returns the input named name of the facts f.
func inputNamed(t *testing.T, f *derive.Facts, name string) derive.Value {
	t.Helper()

	at := slices.IndexFunc(f.Inputs, func(v derive.Value) bool { return v.Name == name })
	assert.True(t, at >= 0, "the callable "+f.Callable+" has the input "+name)
	return f.Inputs[at]
}

// sampleTexts returns the texts of the two samples of v.
func sampleTexts(v derive.Value) [2]string {
	return [2]string{v.Samples[0].Value.Text, v.Samples[1].Value.Text}
}

func TestFacts(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("records the derived checks of each callable of an interface", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, storePackage).suites["Store"]

			assert.Equal(t, derivedOf(&got.Derivation), []string{
				"Store.Get/smoke", "Store.Get/cancel", "Store.Get/zero-on-error", "Store.Put/smoke", "Store.Put/cancel",
			}, "each callable takes the claims of every callable, of a context and of its shape")
		})

		t.Run("binds the claims of the classifications of a callable", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, derivedPackage).suites["Store"]

			derived := derivedOf(&got.Derivation)
			assert.True(t, slices.Contains(derived, "Store.Get/read-after-write"),
				"the mixin read-after-write with its write binds its claim")
			assert.True(t, slices.Contains(derived, "Store.Get/miss"), "the mixin not-found binds the claim miss")
			at := slices.IndexFunc(got.Derivation.Checks, func(c derive.Check) bool {
				return c.Identity == "Store.Get/read-after-write"
			})
			assert.Equal(t, got.Derivation.Checks[at].Pairs["write"].Partner.Callable, "Put",
				"the param write of the mixin names the partner")
		})

		t.Run("records the facts of a callable from the rules of its language", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, storePackage).suites["Store"]

			get := factsOf(t, &got, "Get")
			assert.Equal(t, get.Declaration, "Store", "the declaration of Get is Store")
			assert.Equal(t, get.ID.Name, "Get", "the identity of the facts is the identity of the method")
			assert.Equal(t, get.Signature.Errors, rules.ErrorsLastReturn, "Get returns its error last")
			id := inputNamed(t, get, "id")
			assert.Equal(t, id.Shape.Form, symbol.FormText, "the identity is text")
			assert.True(t, id.Samples[0].OK() && id.Samples[1].OK(), "eidos derives two samples of a string")
			assert.True(t, id.HasZero, "a string has a zero value")
			assert.Length(t, get.Results, 1, "Get returns one value beside its error")
			assert.Equal(t, get.Results[0].Fields, []derive.Field{
				{Name: "ID", Shape: get.Results[0].Fields[0].Shape},
				{Name: "Name", Shape: get.Results[0].Fields[1].Shape},
			}, "the item has the fields ID and Name")
			assert.Equal(t, get.Results[0].Fields[1].Shape.Form, symbol.FormText, "the field Name is text")
		})

		t.Run("records the checks of a function under the name of the function", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, slugPackage).suites["Slug"]

			assert.Equal(t, derivedOf(&got.Derivation), []string{"Slug/smoke", "Slug/deterministic"},
				"a computation without a context takes the claims of every callable and of its shape")
		})

		t.Run("records the derived checks of the methods of a struct", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, catalogPackage).suites["Catalog"]

			assert.True(t, slices.Contains(derivedOf(&got.Derivation), "Catalog.Find/smoke"),
				"the exported method Find takes the claim of every callable")
		})

		t.Run("takes the sample that a parameter states as the first sample", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, taggedPackage).suites["Hasher"]

			r := inputNamed(t, factsOf(t, &got, "CombineTagged"), "r")
			assert.Equal(t, r.Samples[0].Value.Text, "129", "the stated sample 0x81 comes first")
			assert.NotEqual(t, r.Samples[1].Value.Text, "129", "eidos derives a second sample that differs")
		})

		t.Run("takes the sample that a parameter of a function states as the first sample", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, taggedPackage).suites["Digest"]

			r := inputNamed(t, factsOf(t, &got, "Digest"), "r")
			assert.Equal(t, r.Samples[0].Value.Text, "3", "the stated sample 0x03 comes first")
		})

		t.Run("orders the stated sample, the declared default and the derived sample", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, derivedPackage).suites["Store"]

			count := factsOf(t, &got, "Count")
			assert.Equal(t, sampleTexts(inputNamed(t, count, "stated"))[0], "7",
				"the stated sample comes before the default")
			assert.Equal(t, sampleTexts(inputNamed(t, count, "plain"))[0], "3",
				"the default of Count comes before the derived sample")
			name := inputNamed(t, count, "name")
			assert.True(t, name.Samples[0].OK() && name.Samples[1].OK(), "eidos derives two samples of a Name")
			assert.NotEqual(t, name.Samples[0].Value, emit.Literal(emit.LiteralString, ""),
				"a type without a default takes a derived sample")
			assert.NotEqual(t, inputNamed(t, count, "refused").Samples[0].Refusal, rules.RefusedNone,
				"a stated sample that refuses keeps its refusal before the default")
			visit := inputNamed(t, count, "visit").Samples[0].Value
			assert.Equal(t, visit.Fields[0].Name, "Pages", "the struct takes its declared default")
			assert.Equal(t, visit.Fields[0].Value.Text, "5", "the field Pages takes its default 5")
		})

		t.Run("records the exported fields of the value of a pointer", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, derivedPackage).suites["Store"]

			find := factsOf(t, &got, "Find")
			names := make([]string, 0, len(find.Results[0].Fields))
			for _, f := range find.Results[0].Fields {
				names = append(names, f.Name)
			}
			assert.Equal(t, names, []string{"ID", "User"}, "the session has two exported fields")
		})

		t.Run("records the inputs of a generic interface at its witnesses", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, genericPackage).suites["Cache"]

			key := inputNamed(t, factsOf(t, &got, "Get"), "key")
			assert.Equal(t, key.Shape.Form, symbol.FormScalar, "the key takes the witness int")
		})

		t.Run("takes the declared default of a type as the first sample", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, pathsPackage).suites["Router"]

			hold := factsOf(t, &got, "Hold")
			assert.Equal(t, sampleTexts(inputNamed(t, hold, "score"))[0], "3", "the default of Score comes first")
		})
	})
}
