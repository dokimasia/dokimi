// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package derive_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

func TestGate(t *testing.T) {
	t.Parallel()

	t.Run("Passes", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			giveGate spec.GateID
			give     derive.Facts
			want     bool
		}{
			{
				name:     "passes failure for a callable with an error result",
				giveGate: spec.Failure,
				give:     reader(),
				want:     true,
			},
			{
				name:     "fails failure for a callable without an error model",
				giveGate: spec.Failure,
				give:     pure(),
				want:     false,
			},
			{
				name: "passes returns-beside-error for the error model last-return", giveGate: spec.ReturnsBesideError,
				give: reader(), want: true,
			},
			{
				name: "fails returns-beside-error for a thrown error", giveGate: spec.ReturnsBesideError,
				give: withErrors(reader(), rules.ErrorsThrown), want: false,
			},
			{
				name:     "passes cancellation for a callable with a context",
				giveGate: spec.Cancellation,
				give:     reader(),
				want:     true,
			},
			{
				name:     "fails cancellation for a callable without a context",
				giveGate: spec.Cancellation,
				give:     pure(),
				want:     false,
			},
			{
				name:     "passes one-value for a callable with one value",
				giveGate: spec.OneValue,
				give:     reader(),
				want:     true,
			},
			{
				name:     "fails one-value for a callable with two values",
				giveGate: spec.OneValue,
				give:     lookup(),
				want:     false,
			},
			{
				name:     "passes two-values for a callable with two values",
				giveGate: spec.TwoValues,
				give:     lookup(),
				want:     true,
			},
			{
				name:     "fails two-values for a callable with one value",
				giveGate: spec.TwoValues,
				give:     reader(),
				want:     false,
			},
			{name: "passes no-value for a callable without values", giveGate: spec.NoValue, give: writer(), want: true},
			{name: "fails no-value for a callable with a value", giveGate: spec.NoValue, give: reader(), want: false},
			{
				name:     "passes presence-flag for a callable with an ok-bool",
				giveGate: spec.PresenceFlag,
				give:     lookup(),
				want:     true,
			},
			{
				name:     "fails presence-flag for a callable without an ok-bool",
				giveGate: spec.PresenceFlag,
				give:     reader(),
				want:     false,
			},
			{
				name: "passes optional-value for an optional result", giveGate: spec.OptionalValue,
				give: withResult(reader(), symbol.FormOptional), want: true,
			},
			{
				name:     "fails optional-value for a struct result",
				giveGate: spec.OptionalValue,
				give:     reader(),
				want:     false,
			},
			{
				name: "passes list-value for a list result", giveGate: spec.ListValue,
				give: withResult(reader(), symbol.FormList), want: true,
			},
			{name: "passes list-value for a stream return", giveGate: spec.ListValue, give: stream(), want: true},
			{name: "fails list-value for a struct result", giveGate: spec.ListValue, give: reader(), want: false},
			{
				name: "passes ordered-value for a number result", giveGate: spec.OrderedValue,
				give: withResult(reader(), symbol.FormScalar), want: true,
			},
			{
				name: "passes ordered-value for a string result", giveGate: spec.OrderedValue,
				give: withResult(reader(), symbol.FormText), want: true,
			},
			{name: "fails ordered-value for a struct result", giveGate: spec.OrderedValue, give: reader(), want: false},
			{
				name:     "passes value-samples for two distinct samples",
				giveGate: spec.ValueSamples,
				give:     reader(),
				want:     true,
			},
			{
				name: "fails value-samples for two equal samples", giveGate: spec.ValueSamples,
				give: withSameSamples(reader()), want: false,
			},
			{
				name: "fails value-samples for a refused sample", giveGate: spec.ValueSamples,
				give: withRefusedSample(reader()), want: false,
			},
			{
				name: "passes optional-input for an optional input", giveGate: spec.OptionalInput,
				give: withInput(reader(), symbol.FormOptional), want: true,
			},
			{name: "fails optional-input for a text input", giveGate: spec.OptionalInput, give: reader(), want: false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, defined := derive.Passes(tt.giveGate, &tt.give)

				assert.True(t, defined, "the derivation has a condition of the gate")
				assert.Equal(t, got, tt.want, "the gate judges the facts")
			})
		}

		t.Run("returns false for a gate without a condition on the facts of a callable", func(t *testing.T) {
			t.Parallel()
			f := reader()

			_, defined := derive.Passes(spec.PartnerFound, &f)

			assert.False(t, defined, "the partner resolution judges the gate partner-found")
		})

		t.Run("has a condition for each gate of the spec that judges the facts of a callable", func(t *testing.T) {
			t.Parallel()
			f := reader()

			for _, g := range spec.Gates() {
				_, defined := derive.Passes(g.ID, &f)
				expect.Equal(t, defined, !slices.Contains(derive.OtherGates, g.ID),
					"the gate "+string(g.ID)+" has a condition, or another part of the derivation judges it")
			}
		})
	})
}
