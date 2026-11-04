// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rules_test

import (
	"maps"
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/plugin/shape"

	"go.dokimi.dev/dokimi/rules"
)

func TestReference(t *testing.T) {
	t.Parallel()

	t.Run("References", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the forms of references.json in sorted order", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(rules.References(), func(f rules.ReferenceForm) rules.ReferenceFormID { return f.ID }),
				ruleIDs(t, "references", "references"),
				"the forms are the keys of references.json",
			)
		})

		t.Run("returns forms with an effect for each role of their contract", func(t *testing.T) {
			t.Parallel()

			for _, f := range rules.References() {
				s, ok := shape.SpecOf(f.Contract)
				assert.True(t, ok && s.Form == shape.FormContract,
					"the shape catalog has the contract of form "+string(f.ID))
				roles := make([]string, 0, len(s.Roles))
				for _, r := range s.Roles {
					roles = append(roles, string(r.Name))
				}
				expect.Equal(t, slices.Sorted(maps.Keys(f.Roles)), slices.Sorted(slices.Values(roles)),
					"form "+string(f.ID)+" has an effect for each role of its contract")
			}
		})

		t.Run("returns forms whose effects act on the state of the form", func(t *testing.T) {
			t.Parallel()

			for _, f := range rules.References() {
				effects := slices.Collect(maps.Values(f.Roles))
				if f.Writers != "" {
					effects = append(effects, f.Writers)
				}
				slices.Sort(effects)
				want := []rules.Effect{rules.EffectAppend}
				if f.State == rules.StateMap {
					want = []rules.Effect{rules.EffectLoad, rules.EffectStore}
				}
				expect.Equal(t, slices.Compact(effects), want,
					"form "+string(f.ID)+" has the effects of its state "+string(f.State))
			}
		})
	})

	t.Run("ReferenceOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the form with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := rules.ReferenceOf("cache")

			assert.True(t, ok, "the rules have the form")
			expect.Equal(t, got, rules.ReferenceForm{
				ID: "cache", Contract: "cache", State: rules.StateMap,
				Roles:   map[string]rules.Effect{"cache": rules.EffectLoad, "backing": rules.EffectLoad},
				Writers: rules.EffectStore, Since: "1.1.0",
			}, "the form of cache loads through both roles and stores through its writers")
		})

		t.Run("returns false for an identifier without a form", func(t *testing.T) {
			t.Parallel()

			_, ok := rules.ReferenceOf("no-such-form")

			assert.False(t, ok, "the rules have no such form")
		})
	})
}
