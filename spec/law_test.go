// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/spec"
)

func TestLaw(t *testing.T) {
	t.Parallel()

	t.Run("Laws", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the laws of laws.json in the order of the file", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, keysOf(spec.Laws(), func(l spec.Law) spec.LawID { return l.ID }), idsOf(t, "laws", "laws"),
				"the laws keep the order of laws.json")
		})

		t.Run("returns laws that the schema accepts", func(t *testing.T) {
			t.Parallel()

			validates(t, "laws")
		})

		t.Run("returns laws whose depth is a machinery of the spec", func(t *testing.T) {
			t.Parallel()

			for _, l := range spec.Laws() {
				_, ok := spec.MachineryOf(l.Depth)
				expect.True(t, ok, "the depth of law "+string(l.ID)+" is a machinery of the spec")
			}
		})

		t.Run("returns laws that are distinct from other laws of the spec", func(t *testing.T) {
			t.Parallel()

			for _, l := range spec.Laws() {
				for _, d := range l.DistinctFrom {
					_, ok := spec.LawOf(d.Law)
					expect.True(t, ok && d.Law != l.ID, "law "+string(l.ID)+" is distinct from another law of the spec")
					expect.True(
						t,
						d.Case != "",
						"law "+string(l.ID)+" states the case that separates it from "+string(d.Law),
					)
				}
			}
		})
	})

	t.Run("LawOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the law with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.LawOf(spec.NoValueBesideAnError)

			assert.True(t, ok, "the spec has the law")
			assert.Equal(t, got.Statement, "A call that fails returns the zero value of each value result.",
				"the law has its statement")
		})

		t.Run("returns false for an identifier without a law", func(t *testing.T) {
			t.Parallel()

			_, ok := spec.LawOf("no-such-law")

			assert.False(t, ok, "the spec has no such law")
		})
	})
}
