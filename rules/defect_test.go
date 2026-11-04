// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/rules"
)

func TestDefect(t *testing.T) {
	t.Parallel()

	t.Run("DefectForms", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each form of the spec with its rule, in the order of the spec", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(rules.DefectForms(), func(f rules.DefectForm) spec.DefectFormID { return f.ID }),
				keysOf(spec.DefectForms(), func(f spec.DefectForm) spec.DefectFormID { return f.ID }),
				"each form of the spec has a rule",
			)
			for _, f := range rules.DefectForms() {
				want, _ := spec.DefectFormOf(f.ID)
				expect.Equal(t, f.DefectForm, want, "form "+string(f.ID)+" has the runtime contract of the spec")
			}
		})

		t.Run("returns forms that each have a rule of defects.json", func(t *testing.T) {
			t.Parallel()
			ids := keysOf(spec.DefectForms(), func(f spec.DefectForm) spec.DefectFormID { return f.ID })
			slices.Sort(ids)

			assert.Equal(t, ruleIDs(t, "defects", "forms"), ids,
				"defects.json has a rule for each form of the spec and for no other form")
		})

		t.Run("returns forms that a sequence of at least one call shows", func(t *testing.T) {
			t.Parallel()

			for _, f := range rules.DefectForms() {
				expect.True(t, f.MinSteps >= 1, "the form "+string(f.ID)+" needs at least one call")
			}
		})
	})

	t.Run("DefectFormOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the form with the callable that it changes", func(t *testing.T) {
			t.Parallel()

			got, ok := rules.DefectFormOf(spec.PutRefuses)

			assert.True(t, ok, "the rules have the form")
			assert.Equal(t, got.Changes, "partner:put", "the form changes the partner put")
		})

		t.Run("returns false for an identifier without a form", func(t *testing.T) {
			t.Parallel()

			_, ok := rules.DefectFormOf("no-such-form")

			assert.False(t, ok, "the rules have no such form")
		})
	})
}
