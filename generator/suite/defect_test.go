// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/suite"
)

// defectNamed returns the planted defect name of the suite s.
func defectNamed(t *testing.T, s *suite.Suite, name string) *suite.Defect {
	t.Helper()

	at := slices.IndexFunc(s.Defects, func(d *suite.Defect) bool { return d.Name == name })
	assert.True(t, at >= 0, "the suite has the planted defect "+name)
	return s.Defects[at]
}

func TestDefect(t *testing.T) {
	t.Parallel()

	families := generate(t, goTree, familiesPackage).suites
	family := func(name string) *suite.Suite {
		s := families[name]
		return &s
	}

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name          string
			giveSuite     string
			give          string
			wantBody      string
			wantInput     int
			wantBudget    time.Duration
			wantDelegates bool
			wantCounts    bool
		}{
			{
				name: "writes a panic without a call of the subject", giveSuite: "Store", give: "StoreGetPanics",
				wantBody: "panic",
			},
			{
				name: "panics on the repeat after a first call of the subject", giveSuite: "Lifecycle",
				give: "LifecycleStopPanicsOnRepeat", wantBody: "panic-later", wantDelegates: true, wantCounts: true,
			},
			{
				name: "fails the repeat after a first call of the subject", giveSuite: "Lifecycle",
				give: "LifecyclePutRefusesRepeat", wantBody: "fail-later", wantDelegates: true, wantCounts: true,
			},
			{
				name: "succeeds on the repeat after a first call of the subject", giveSuite: "Guarded",
				give: "GuardedPutIgnoresThePredicate", wantBody: "zero-later", wantDelegates: true, wantCounts: true,
			},
			{
				name: "succeeds in place of the rollback of the subject", giveSuite: "Ledger",
				give: "LedgerRollbackSettlesTwice", wantBody: "zero",
			},
			{
				name: "fails every call without a call of the subject", giveSuite: "Lifecycle",
				give: "LifecycleSyncFailsAlways", wantBody: "fail",
			},
			{
				name: "fails with a sample beside the error", giveSuite: "Store", give: "StoreGetValueBeside",
				wantBody: "fail-beside",
			},
			{
				name: "returns a sample beside a false flag of the subject", giveSuite: "Store",
				give: "StoreLoadValueBesideAFalseFlag", wantBody: "beside-false", wantDelegates: true,
			},
			{
				name: "drifts from the zero value to a sample on its repeat", giveSuite: "Store",
				give: "StoreLoadDrifts", wantBody: "drift", wantCounts: true,
			},
			{
				name: "returns the first element of the subject twice", giveSuite: "Bag",
				give: "BagItemsDrainsTwice", wantBody: "drain-twice", wantDelegates: true,
			},
			{
				name: "reverses the drain of the subject on its repeat", giveSuite: "Bag",
				give: "BagItemsReordersDrain", wantBody: "reorder", wantDelegates: true, wantCounts: true,
			},
			{
				name: "reverses every replay of the subject", giveSuite: "Log", give: "LogReplayReplaysReversed",
				wantBody: "reverse", wantDelegates: true,
			},
			{
				name: "returns the first input of the call", giveSuite: "Renderer", give: "RendererRenderEchoes",
				wantBody: "echo", wantInput: 1,
			},
			{
				name: "reads the first page of the subject for every cursor", giveSuite: "Pager",
				give: "PagerFetchServesOnePage", wantBody: "first-page", wantInput: 1, wantDelegates: true,
			},
			{
				name: "wraps the handle that the subject opens", giveSuite: "Opener",
				give: "OpenerOpenOpensALiveHandle", wantBody: "live-handle", wantDelegates: true,
			},
			{
				name: "reports leadership on every call", giveSuite: "Election", give: "ElectionIsLeaderLeadsOn",
				wantBody: "lead",
			},
			{
				name: "waits twice the budget before the call of the subject", giveSuite: "Timers",
				give: "TimersSlowSlow", wantBody: "slow", wantBudget: 10 * time.Second, wantDelegates: true,
			},
			{
				name: "returns a sample on every call", giveSuite: "Cached", give: "CachedGetAnswersASample",
				wantBody: "sample",
			},
			{
				name: "returns success and stores nothing", giveSuite: "Store", give: "StorePutForgetsWrites",
				wantBody: "zero",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				d := defectNamed(t, family(tt.giveSuite), tt.give)
				expect.Equal(t, d.Body, tt.wantBody, "the defect has the body of its form")
				expect.Equal(t, d.Type, tt.give+"Defect", "the type of the defect adds the word Defect to its name")
				expect.Equal(t, d.Input, tt.wantInput, "the defect knows the position of the first input")
				expect.Equal(t, d.Budget, tt.wantBudget, "the defect waits twice the budget of the timeout")
				expect.Equal(t, d.Delegates, tt.wantDelegates, "the defect calls the subject where its body does")
				expect.Equal(t, d.Counts, tt.wantCounts, "the defect counts its calls where its body does")
			})
		}

		t.Run("names the partner that a form changes by the role of the partner", func(t *testing.T) {
			t.Parallel()

			d := defectNamed(t, family("Election"), "ElectionIsLeaderLeadsOn")
			expect.Equal(t, d.Method.Name, "IsLeader", "the form leads-on changes the partner in the role is-leader")
		})

		t.Run("writes the defect of a function as a function without a type", func(t *testing.T) {
			t.Parallel()

			slug := generate(t, goTree, slugPackage).suites["Slug"]
			d := defectNamed(t, &slug, "SlugPanics")
			expect.Equal(t, d.Type, "", "the defect of a function has no type")
		})

		t.Run("records the version field, the variadic argument and the stream of a defect", func(t *testing.T) {
			t.Parallel()

			log := defectNamed(t, family("Log"), "LogAppendPanics")
			expect.Equal(t, log.Version, "Seq", "the defect of Log has the version field of the mixin causal")
			find := defectNamed(t, family("Store"), "StoreFindPanics")
			expect.Equal(t, find.Args, []string{"p0", "p1..."}, "the defect passes the variadic keys with an ellipsis")
			all := defectNamed(t, family("Plain"), "PlainAllPanics")
			assert.Length(t, all.Results, 1, "the stream is the one result of All")
			expect.Equal(t, all.Results[0].Role, "value", "a stream is a value of the defect")
			expect.Equal(t, *all.Results[0].Sample, suite.Expr{Word: "nil"}, "the sample of a stream is nil")
		})

		t.Run("returns a sample that differs from the zero value", func(t *testing.T) {
			t.Parallel()

			d := defectNamed(t, family("Store"), "StoreLoadDrifts")
			assert.NotNil(t, d.Results[0].Sample, "the value of Load has a sample")
			expect.Equal(t, d.Results[1].Role, "flag", "the second result of Load is its flag")
			expect.Nil(t, d.Results[1].Sample, "a flag has no sample")
		})

		t.Run("records the method next of the live handle of an opened cursor", func(t *testing.T) {
			t.Parallel()

			live := defectNamed(t, family("Opener"), "OpenerOpenOpensALiveHandle")
			assert.NotNil(t, live.Handle, "the defect returns a live handle")
			expect.Equal(t, live.Handle.Method.Name, "Next", "the handle overrides the method that next-method names")
			expect.Length(t, live.Handle.Results, 3, "the handle returns the value, the flag and the error of Next")
			expect.Equal(t, live.Handle.Body, "zero", "the method of the handle returns its zero results")
		})
	})
}
