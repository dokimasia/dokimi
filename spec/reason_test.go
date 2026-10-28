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

func TestReason(t *testing.T) {
	t.Parallel()

	t.Run("Reasons", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the reason codes of reasons.json in the order of the file", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(spec.Reasons(), func(r spec.Reason) spec.ReasonID { return r.ID }),
				idsOf(t, "reasons", "reasons"),
				"the reason codes keep the order of reasons.json",
			)
		})

		t.Run("returns reason codes that the schema accepts", func(t *testing.T) {
			t.Parallel()

			validates(t, "reasons")
		})

		t.Run("returns reason codes that a probe reports", func(t *testing.T) {
			t.Parallel()
			var reported []spec.ReasonID
			for _, p := range spec.Probes() {
				reported = append(reported, p.Reasons...)
			}

			for _, r := range spec.Reasons() {
				expect.True(t, slices.Contains(reported, r.ID), "a probe reports the reason code "+string(r.ID))
			}
		})
	})

	t.Run("ReasonOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the reason code with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.ReasonOf(spec.ValueBesideError)

			assert.True(t, ok, "the spec has the reason code")
			assert.Equal(t, got.Meaning, "The call failed and returned a value other than the zero value.",
				"the reason code has its meaning")
		})

		t.Run("returns false for an identifier without a reason code", func(t *testing.T) {
			t.Parallel()

			_, ok := spec.ReasonOf("no-such-reason")

			assert.False(t, ok, "the spec has no such reason code")
		})
	})
}
