// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
	"go.dokimi.dev/dokimi/spec"
)

// refusedBy returns the refusal of the claim of the callable by the
// directive refuse with the reason why.
func refusedBy(claim spec.ClaimID, callable, key, why string) derive.Refusal {
	return derive.Refusal{
		Claim: claim, Callable: callable, Key: key, Gate: spec.Unrefused, Why: why,
		Remedy: "remove the directive refuse of the claim where " + callable + " makes it",
	}
}

func TestRefusal(t *testing.T) {
	t.Parallel()

	got := generate(t, goTree, refusedPackage)

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("records the refusal of each claim that a directive refuse of a method names", func(t *testing.T) {
			t.Parallel()

			clock, found := got.suites["HybridClock"]
			assert.True(t, found, "suite records the suite of HybridClock")
			expect.Equal(t, clock.Derivation.Refusals, []derive.Refusal{
				refusedBy(spec.Cancel, "Update", "hybrid-clock.update/cancel",
					"Update returns before it reads its context"),
				refusedBy(spec.ZeroOnError, "Update", "hybrid-clock.update/zero-on-error",
					"Update fails only for an instant beyond the max offset of the clock"),
			}, "the derivation refuses each claim with the reason of its directive")
			expect.False(t, slices.ContainsFunc(clock.Derived, func(d suite.Derived) bool {
				return d.Field == "Update" && (d.Claim == "ZeroOnError" || d.Claim == "Cancel")
			}), "the suite has no check of a refused claim")
		})

		t.Run("records no refusal of a callable without a directive refuse", func(t *testing.T) {
			t.Parallel()

			clock := got.suites["HybridClock"]
			assert.Length(t, clock.Facts, 2, "the clock has two callables")
			assert.Equal(t, clock.Facts[1].Callable, "MaxOffset", "the second callable is MaxOffset")
			expect.Nil(t, clock.Facts[1].Refused, "MaxOffset refuses no claim")
		})

		t.Run("records the refusal of a claim that a directive refuse of a function names", func(t *testing.T) {
			t.Parallel()

			title, found := got.suites["Title"]
			assert.True(t, found, "suite records the suite of Title")
			expect.Equal(t, title.Derivation.Refusals, []derive.Refusal{
				refusedBy(spec.Deterministic, "Title", "title/deterministic", "Title salts each name"),
			}, "the function refuses the claim with the reason of its directive")
		})

		t.Run("records the first refusal of a claim that two directives refuse name", func(t *testing.T) {
			t.Parallel()

			doubled, found := got.suites["Doubled"]
			assert.True(t, found, "suite records the suite of Doubled")
			expect.Equal(t, doubled.Derivation.Refusals, []derive.Refusal{
				refusedBy(spec.ZeroOnError, "Read", "doubled.read/zero-on-error", "Read never fails"),
			}, "the first directive refuses the claim")
		})

		t.Run("reports UnknownRefusal at each refusal that refuses no claim of a suite", func(t *testing.T) {
			t.Parallel()

			var messages []string
			for _, d := range got.findings {
				if d.Code == suite.UnknownRefusal {
					messages = append(messages, d.Msg)
				}
			}
			slices.Sort(messages)
			expect.Equal(t, messages, []string{
				"Doubled.Read refuses the claim zero-on-error twice, and a callable refuses each claim once",
				"Marker.Mark refuses the claim deterministic, which the derivation does not license for Mark, so " +
					"Marker has no suite",
				"Pair.Write refuses the claim zero-on-error, which the derivation does not license for Write, so " +
					"Pair has no suite. The claims of Write are smoke and cancel",
				"Shout refuses the claim cancel, which the derivation does not license for Shout, so Shout has no " +
					"suite. The claims of Shout are smoke and deterministic",
				"Spinner.Spin refuses the claim cancel, which the derivation does not license for Spin, so Spinner " +
					"has no suite. The claim of Spin is smoke",
				"Tally.Record refuses the claim deterministic, which the derivation does not license for Record, so " +
					"Tally has no suite. The claims of Record are smoke and cancel",
				"Unsuited has no directive suite, so the refusal of smoke by Unsuited.Get refuses no claim",
				"Untitled has no directive suite, so the refusal of deterministic refuses no claim",
				"Watch.Now refuses the claim zero-on-error, which the derivation does not license for Now, so Watch " +
					"has no suite. The claims of Now are smoke, deterministic, time-aware and steady",
			}, "suite reports each refusal that refuses no claim of a suite")
			for _, name := range []string{"Tally", "Watch", "Spinner", "Marker", "Pair", "Shout"} {
				_, found := got.suites[name]
				expect.False(t, found, "suite records no suite of "+name+", which refuses an unknown claim")
			}
		})
	})
}
