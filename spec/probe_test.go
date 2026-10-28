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

func TestProbe(t *testing.T) {
	t.Parallel()

	t.Run("Probes", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the probes of probes.json in the order of the file", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				keysOf(spec.Probes(), func(p spec.Probe) spec.ProbeID { return p.ID }),
				idsOf(t, "probes", "probes"),
				"the probes keep the order of probes.json",
			)
		})

		t.Run("returns probes that the schema accepts", func(t *testing.T) {
			t.Parallel()

			validates(t, "probes")
		})

		t.Run("returns probes whose reason codes and input kinds the spec defines", func(t *testing.T) {
			t.Parallel()
			kinds := []spec.InputKind{
				spec.KindCall, spec.KindValue, spec.KindError, spec.KindDuration, spec.KindCount, spec.KindProjection,
				spec.KindScope, spec.KindHook, spec.KindPredicate, spec.KindClock, spec.KindValues,
			}

			for _, p := range spec.Probes() {
				expect.True(t, len(p.Reasons) > 0, "probe "+string(p.ID)+" reports a reason code")
				for _, r := range p.Reasons {
					_, ok := spec.ReasonOf(r)
					expect.True(t, ok, "the spec has the reason code "+string(r)+" of probe "+string(p.ID))
				}
				for _, in := range p.Inputs {
					expect.True(
						t,
						slices.Contains(kinds, in.Kind),
						"input "+in.Name+" of probe "+string(p.ID)+" has a kind of the spec",
					)
				}
				expect.True(t, p.Vacuous != "", "probe "+string(p.ID)+" states when it is vacuous")
			}
		})

		t.Run("returns probes that a claim calls", func(t *testing.T) {
			t.Parallel()
			called := make([]spec.ProbeID, 0, len(spec.Claims()))
			for _, c := range spec.Claims() {
				called = append(called, c.Probe)
			}

			for _, p := range spec.Probes() {
				expect.True(t, slices.Contains(called, p.ID), "a claim calls the probe "+string(p.ID))
			}
		})
	})

	t.Run("ProbeOf", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the probe with the identifier", func(t *testing.T) {
			t.Parallel()

			got, ok := spec.ProbeOf(spec.ZeroBesideError)

			assert.True(t, ok, "the spec has the probe")
			assert.Equal(
				t,
				got.Inputs,
				[]spec.Input{{Name: "call", Kind: spec.KindCall}, {Name: "zero", Kind: spec.KindValue}},
				"the probe has its inputs in order",
			)
		})

		t.Run("returns false for an identifier without a probe", func(t *testing.T) {
			t.Parallel()

			_, ok := spec.ProbeOf("no-such-probe")

			assert.False(t, ok, "the spec has no such probe")
		})
	})
}
