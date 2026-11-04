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

func TestLaw(t *testing.T) {
	t.Parallel()

	t.Run("Laws", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each law of the spec with its rule, in the order of the spec", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(rules.Laws(), func(l rules.Law) spec.LawID { return l.ID }),
				keysOf(spec.Laws(), func(l spec.Law) spec.LawID { return l.ID }),
				"each law of the spec has a rule",
			)
			for _, l := range rules.Laws() {
				want, _ := spec.LawOf(l.ID)
				expect.Equal(t, l.Law, want, "law "+string(l.ID)+" has the runtime contract of the spec")
			}
		})

		t.Run("returns laws that each have a rule of laws.json", func(t *testing.T) {
			t.Parallel()
			ids := keysOf(spec.Laws(), func(l spec.Law) spec.LawID { return l.ID })
			slices.Sort(ids)

			assert.Equal(t, ruleIDs(t, "laws", "laws"), ids,
				"laws.json has a rule for each law of the spec and for no other law")
		})

		t.Run("returns laws whose depth is a machinery of the spec", func(t *testing.T) {
			t.Parallel()

			for _, l := range rules.Laws() {
				_, ok := spec.MachineryOf(l.Depth)
				expect.True(t, ok, "the depth of law "+string(l.ID)+" is a machinery of the spec")
			}
		})

		t.Run("returns laws that are distinct from other laws of the spec", func(t *testing.T) {
			t.Parallel()

			for _, l := range rules.Laws() {
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

		t.Run("returns the law with its statement and its depth", func(t *testing.T) {
			t.Parallel()

			got, ok := rules.LawOf(spec.NoValueBesideAnError)

			assert.True(t, ok, "the rules have the law")
			assert.Equal(t, got.Statement, "A call that fails returns the zero value of each value result.",
				"the law has the statement of the spec")
			assert.Equal(t, got.Depth, spec.NothingBeyondACaller, "the law has its depth")
		})

		t.Run("returns false for an identifier without a law", func(t *testing.T) {
			t.Parallel()

			_, ok := rules.LawOf("no-such-law")

			assert.False(t, ok, "the rules have no such law")
		})
	})
}
