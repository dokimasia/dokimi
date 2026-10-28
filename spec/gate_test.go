// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/spec"
)

func TestGate(t *testing.T) {
	t.Parallel()

	t.Run("Gates", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the gates of gates.json in the order of the file", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(spec.Gates(), func(g spec.Gate) spec.GateID { return g.ID }),
				idsOf(t, "gates", "gates"),
				"the gates keep the order of gates.json",
			)
		})

		t.Run("returns gates that the schema accepts", func(t *testing.T) {
			t.Parallel()

			validates(t, "gates")
		})

		t.Run("returns gates with a reason and a remedy", func(t *testing.T) {
			t.Parallel()

			for _, g := range spec.Gates() {
				expect.True(t, g.Why != "" && g.Remedy != "", "gate "+string(g.ID)+" has a reason and a remedy")
			}
		})
	})

	t.Run("GateOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the gate with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.GateOf(spec.ReturnsBesideError)

			assert.True(t, ok, "the spec has the gate")
			assert.Equal(t, got.Why, "{callable} returns no value when it fails", "the gate has its reason")
		})

		t.Run("returns false for an identifier without a gate", func(t *testing.T) {
			t.Parallel()

			_, ok := spec.GateOf("no-such-gate")

			assert.False(t, ok, "the spec has no such gate")
		})
	})
}
