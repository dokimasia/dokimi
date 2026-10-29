// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/spec"
)

func TestDefect(t *testing.T) {
	t.Parallel()

	t.Run("DefectForms", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the forms of defects.json in the order of the file", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, keysOf(spec.DefectForms(), func(f spec.DefectForm) spec.DefectFormID { return f.ID }),
				idsOf(t, "defects", "forms"), "the forms keep the order of defects.json")
		})

		t.Run("returns forms that the schema accepts", func(t *testing.T) {
			t.Parallel()

			validates(t, "defects")
		})

		t.Run("returns forms that a claim plants with a sequence of at least one call", func(t *testing.T) {
			t.Parallel()
			planted := make([]spec.DefectFormID, 0, len(spec.Claims()))
			for _, c := range spec.Claims() {
				planted = append(planted, c.Defect.Form)
			}

			for _, f := range spec.DefectForms() {
				expect.True(t, slices.Contains(planted, f.ID), "a claim plants the form "+string(f.ID))
				expect.True(t, f.MinSteps >= 1, "the form "+string(f.ID)+" needs at least one call")
			}
		})
	})

	t.Run("DefectFormOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the form with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.DefectFormOf(spec.PutRefuses)

			assert.True(t, ok, "the spec has the form")
			assert.Equal(t, got.Changes, "partner:put", "the form changes the partner put")
		})

		t.Run("returns false for an identifier without a form", func(t *testing.T) {
			t.Parallel()

			_, ok := spec.DefectFormOf("no-such-form")

			assert.False(t, ok, "the spec has no such form")
		})
	})
}
