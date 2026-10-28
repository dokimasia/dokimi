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

func TestMachinery(t *testing.T) {
	t.Parallel()

	t.Run("Machineries", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the machineries of machineries.json in the order of the file", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, keysOf(spec.Machineries(), func(m spec.Machinery) spec.MachineryID { return m.ID }),
				idsOf(t, "machineries", "machineries"), "the machineries keep the order of machineries.json")
		})

		t.Run("returns machineries that the schema accepts", func(t *testing.T) {
			t.Parallel()

			validates(t, "machineries")
		})

		t.Run("returns a param of its own for each machinery of the consumer and for no other", func(t *testing.T) {
			t.Parallel()
			params := make([]string, 0, len(spec.Machineries()))
			for _, m := range spec.Machineries() {
				expect.True(t, slices.Contains([]spec.Supplier{spec.ByRuntime, spec.ByConsumer, spec.BySim}, m.By),
					"machinery "+string(m.ID)+" has a supplier of the spec")
				expect.Equal(
					t,
					m.Param != "",
					m.By == spec.ByConsumer,
					"machinery "+string(m.ID)+" has a param when the consumer supplies it",
				)
				expect.False(
					t,
					m.Param != "" && slices.Contains(params, m.Param),
					"the param of machinery "+string(m.ID)+" is its own",
				)
				params = append(params, m.Param)
			}
		})
	})

	t.Run("MachineryOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the machinery with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.MachineryOf(spec.ReferenceImplementation)

			assert.True(t, ok, "the spec has the machinery")
			assert.Equal(t, got.Param, "reference", "the machinery has the param of its directive")
		})

		t.Run("returns false for an identifier without a machinery", func(t *testing.T) {
			t.Parallel()

			_, ok := spec.MachineryOf("no-such-machinery")

			assert.False(t, ok, "the spec has no such machinery")
		})
	})
}
