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

func TestGate(t *testing.T) {
	t.Parallel()

	t.Run("Gates", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each gate of the spec with its rule, in the order of the spec", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(rules.Gates(), func(g rules.Gate) spec.GateID { return g.ID }),
				keysOf(spec.Gates(), func(g spec.Gate) spec.GateID { return g.ID }),
				"each gate of the spec has a rule",
			)
			for _, g := range rules.Gates() {
				want, _ := spec.GateOf(g.ID)
				expect.Equal(t, g.Gate, want, "gate "+string(g.ID)+" has the runtime contract of the spec")
			}
		})

		t.Run("returns gates that each have a rule of gates.json", func(t *testing.T) {
			t.Parallel()
			ids := keysOf(spec.Gates(), func(g spec.Gate) spec.GateID { return g.ID })
			slices.Sort(ids)

			assert.Equal(t, ruleIDs(t, "gates", "gates"), ids,
				"gates.json has a rule for each gate of the spec and for no other gate")
		})

		t.Run("returns gates with a reason and a remedy", func(t *testing.T) {
			t.Parallel()

			for _, g := range rules.Gates() {
				expect.True(t, g.Why != "" && g.Remedy != "", "gate "+string(g.ID)+" has a reason and a remedy")
			}
		})
	})

	t.Run("GateOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the gate with its reason", func(t *testing.T) {
			t.Parallel()

			got, ok := rules.GateOf(spec.ReturnsBesideError)

			assert.True(t, ok, "the rules have the gate")
			assert.Equal(t, got.Why, "{callable} returns no value when it fails", "the gate has its reason")
		})

		t.Run("returns false for an identifier without a gate", func(t *testing.T) {
			t.Parallel()

			_, ok := rules.GateOf("no-such-gate")

			assert.False(t, ok, "the rules have no such gate")
		})
	})
}
