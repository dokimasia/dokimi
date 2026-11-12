// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/symbol"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
)

func TestMachinery(t *testing.T) {
	t.Parallel()

	p := suite.New()
	f := fixture(t, goTree, p, faultsPackage)
	got := recorded(t, f, p)

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name        string
			giveSuite   string
			giveField   string
			giveSegment string
			wantProbe   spec.ProbeID
		}{
			{
				name: "derives a reported failure of each callable with an error", giveSuite: "Cache",
				giveField: "Get", giveSegment: "failure-is-reported", wantProbe: spec.ReportsAFailedDependency,
			},
			{
				name: "derives a failed call without effect of an atomic write", giveSuite: "Writer",
				giveField: "Put", giveSegment: "no-partial-write", wantProbe: spec.FailsWithoutEffect,
			},
			{
				name: "derives a retry under faults of a callable that retries", giveSuite: "Syncer",
				giveField: "Sync", giveSegment: "retry-succeeds-under-faults", wantProbe: spec.SucceedsAfterAFailure,
			},
			{
				name: "derives a rollback under faults of a transaction", giveSuite: "Account",
				giveField: "Run", giveSegment: "transaction-under-faults", wantProbe: spec.RollsBackUnderFaults,
			},
			{
				name: "derives a breaker that opens", giveSuite: "Breaker", giveField: "Do",
				giveSegment: "fails-fast", wantProbe: spec.RefusesWithoutTheDependency,
			},
			{
				name: "derives one instance for each key of a sticky callable", giveSuite: "Balancer",
				giveField: "Route", giveSegment: "keeps-its-instance", wantProbe: spec.UsesOneInstance,
			},
			{
				name: "derives a release on every path of a callable that acquires", giveSuite: "Copier",
				giveField: "Copy", giveSegment: "releases-on-every-path", wantProbe: spec.ReleasesWhatItAcquires,
			},
			{
				name: "derives a retry under faults of a function", giveSuite: "Fetch",
				giveSegment: "retry-succeeds-under-faults", wantProbe: spec.SucceedsAfterAFailure,
			},
			{
				name: "derives a reported failure of a method of a struct", giveSuite: "Mirror", giveField: "Get",
				giveSegment: "failure-is-reported", wantProbe: spec.ReportsAFailedDependency,
			},
			{
				name:      "derives the claims of a machinery that the declaration states beside an unknown one",
				giveSuite: "Gravity", giveField: "Get", giveSegment: "failure-is-reported",
				wantProbe: spec.ReportsAFailedDependency,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				s := got.suites[tt.giveSuite]
				d := derivedNamed(t, &s, tt.giveField, tt.giveSegment)
				expect.Equal(t, d.Probe, string(tt.wantProbe), "the check calls the probe of its claim")
				expect.Equal(t, d.Machineries,
					[]string{string(spec.CallThatFailsOnDemand), string(spec.NothingBeyondACaller)},
					"the check lists the machinery of failing calls first")
			})
		}

		t.Run("lists the classification of a claim of a machinery that is not stated", func(t *testing.T) {
			t.Parallel()

			s := got.suites["Undeclared"]
			at := slices.IndexFunc(s.Derivation.Unclaimed, func(u derive.Unclaimed) bool {
				return u.Callable == "Put" && u.Classification == "mixin:atomic"
			})
			assert.True(t, at >= 0, "the suite lists the mixin atomic of Put as unclaimed")
			u := s.Derivation.Unclaimed[at]
			expect.Equal(t, u.Stance, spec.Needs, "the stance names the missing machinery")
			expect.Equal(t, u.Note,
				"The claim no-partial-write needs the param machineries=[faults] of the directive suite.",
				"the note names the param that states the machinery")
			reported := slices.ContainsFunc(s.Derived, func(d suite.Derived) bool {
				return d.Segment == "failure-is-reported"
			})
			expect.False(t, reported, "a declaration without the machinery derives no reported failure")
		})

		t.Run("stamps the machineries that the declaration states", func(t *testing.T) {
			t.Parallel()

			key, registered := meta.Lookup[[]string](f.Keys, suite.KeyMachineries)
			assert.True(t, registered, "suite registers the key of the machineries")
			cache, stamped := meta.Get(f.Facts, f.Find(t, symbol.KindInterface, "", "Cache").Identity(), key)
			assert.True(t, stamped, "suite stamps the machineries of Cache")
			expect.Equal(t, cache, []string{string(spec.CallThatFailsOnDemand)}, "the param faults names failing calls")
			gravity, _ := meta.Get(f.Facts, f.Find(t, symbol.KindInterface, "", "Gravity").Identity(), key)
			expect.Equal(t, gravity, []string{string(spec.CallThatFailsOnDemand)},
				"the machineries of Gravity leave out the value that names no machinery")
			_, undeclared := meta.Get(f.Facts, f.Find(t, symbol.KindInterface, "", "Undeclared").Identity(), key)
			expect.False(t, undeclared, "suite stamps nothing on a declaration without the param")
		})

		t.Run("reports UnknownMachinery at a value that names no machinery of a consumer", func(t *testing.T) {
			t.Parallel()

			decl := f.Find(t, symbol.KindInterface, "", "Gravity").Position()
			var found []string
			for _, d := range got.findings {
				if d.Code == suite.UnknownMachinery {
					expect.Equal(t, d.Pos.Line, decl.Line-1, "the finding is at the directive suite of Gravity")
					found = append(found, d.Msg)
				}
			}
			expect.Equal(t, found, []string{
				"no consumer supplies the machinery gravity, so the suite of Gravity derives no claim of it. " +
					"The values of the param machineries are faults, reference, second and medium",
			}, "suite reports the one value that names no machinery")
		})
	})
}
