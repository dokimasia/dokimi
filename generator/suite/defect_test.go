// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"maps"
	"slices"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"

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
	faults := generate(t, goTree, faultsPackage).suites
	maps.Copy(families, generate(t, goTree, drawnPackage).suites)
	maps.Copy(families, generate(t, goTree, classifiedPackage).suites)
	family := func(name string) *suite.Suite {
		s, found := families[name]
		if !found {
			s = faults[name]
		}
		return &s
	}

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name          string
			giveSuite     string
			give          string
			wantBody      string
			wantAcquire   string
			wantInput     int
			wantSecond    int
			wantBudget    time.Duration
			wantDelegates bool
			wantCounts    bool
			wantKeeps     bool
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
				name: "drops the last value of the batch that the subject reads", giveSuite: "Shelf",
				give: "ShelfGetManyDropsTheLastKey", wantBody: "drop-last", wantDelegates: true,
			},
			{
				name: "returns the empty sample of a stream that drains nothing", giveSuite: "Shelf",
				give: "ShelfAllDrainsNothing", wantBody: "sample",
			},
			{
				name: "returns the zero value of a list that drains nothing", giveSuite: "Shelf",
				give: "ShelfSnapshotDrainsNothing", wantBody: "zero",
			},
			{
				name: "fails each call where the call of the subject succeeds", giveSuite: "Lifecycle",
				give: "LifecycleSyncFailsAlways", wantBody: "fail-success", wantDelegates: true,
			},
			{
				name: "returns a sample beside the error of the subject", giveSuite: "Store",
				give: "StoreGetValueBeside", wantBody: "beside-error", wantDelegates: true,
			},
			{
				name: "calls the subject under a context that the caller cannot cancel", giveSuite: "Store",
				give: "StoreGetIgnoresCancel", wantBody: "uncancelled", wantDelegates: true,
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
				name: "returns the value of the write of the subject with the version 0", giveSuite: "Stamped",
				give: "StampedStoreStampsThePast", wantBody: "version-zero", wantDelegates: true,
			},
			{
				name: "returns a sample where the read of the subject succeeds", giveSuite: "Store",
				give: "StoreGetReadsAnotherValue", wantBody: "sample-success", wantDelegates: true,
			},
			{
				name: "succeeds in place of the read of the subject until a dependency fails", giveSuite: "Scope",
				give: "ScopeGetKeepsFailedWrites", wantBody: "zero-before-failure", wantDelegates: true,
				wantKeeps: true,
			},
			{
				name: "returns no error where a dependency of the subject fails", giveSuite: "Cache",
				give: "CacheGetSwallowsFailure", wantBody: "swallow", wantDelegates: true, wantKeeps: true,
			},
			{
				name: "fails each call once a dependency of the subject failed", giveSuite: "Syncer",
				give: "SyncerSyncGivesUpEarly", wantBody: "fail-after-failure", wantDelegates: true, wantKeeps: true,
			},
			{
				name: "reads a sample once a dependency of the subject failed", giveSuite: "Writer",
				give: "WriterGetHalfApplied", wantBody: "sample-after-failure", wantDelegates: true, wantKeeps: true,
			},
			{
				name: "reads a sample once a commit of the subject failed", giveSuite: "Account",
				give: "AccountGetCommitsOnFailure", wantBody: "sample-after-failure", wantDelegates: true,
				wantKeeps: true,
			},
			{
				name: "calls a fresh breaker for each call", giveSuite: "Breaker", give: "BreakerDoNeverOpens",
				wantBody: "fresh-each-call", wantKeeps: true,
			},
			{
				name: "calls a fresh balancer for each call", giveSuite: "Balancer",
				give: "BalancerRouteBalancesEachCall", wantBody: "fresh-each-call", wantKeeps: true,
			},
			{
				name: "acquires once more where a dependency of the subject fails", giveSuite: "Copier",
				give: "CopierCopyOpensWithoutClosing", wantBody: "open-on-failure", wantAcquire: "Open",
				wantDelegates: true,
			},
			{
				name: "returns the first input of a combination of its input type", giveSuite: "Tallies",
				give: "TalliesMergeIgnoresTheSecondInput", wantBody: "echo",
			},
			{
				name:      "passes the zero value as the second input of a combination of another type",
				giveSuite: "Tallies", give: "TalliesDistanceIgnoresTheSecondInput", wantBody: "zero-second",
				wantSecond: 1, wantDelegates: true,
			},
			{
				name: "returns its first input where that input is the value of its first call", giveSuite: "Tallies",
				give: "TalliesMergeKeepsItsFirstResult", wantBody: "keep-first", wantDelegates: true,
			},
			{
				name: "returns the read that another read pinned while the two overlap", giveSuite: "Snapshots",
				give: "SnapshotsGetRacyPin", wantBody: "shared-result", wantDelegates: true, wantKeeps: true,
			},
			{
				name: "returns the version of another read that overlaps it", giveSuite: "Stamped",
				give: "StampedGetRacyCausalRead", wantBody: "shared-result", wantDelegates: true, wantKeeps: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				d := defectNamed(t, family(tt.giveSuite), tt.give)
				expect.Equal(t, d.Body, tt.wantBody, "the defect has the body of its form")
				expect.Equal(t, d.Type, tt.give+"Defect", "the type of the defect adds the word Defect to its name")
				expect.Equal(t, d.Input, tt.wantInput, "the defect knows the position of the first input")
				expect.Equal(t, d.Second, tt.wantSecond, "the defect knows the position of the second input")
				expect.Equal(t, d.Budget, tt.wantBudget, "the defect waits twice the budget of the timeout")
				expect.Equal(t, d.Delegates, tt.wantDelegates, "the defect calls the subject where its body does")
				expect.Equal(t, d.Counts, tt.wantCounts, "the defect counts its calls where its body does")
				expect.Equal(t, d.Keeps, tt.wantKeeps, "the defect keeps its handle where its body reads it")
				expect.Equal(t, d.Acquire, tt.wantAcquire, "the defect acquires through the callable that open names")
			})
		}

		t.Run("acquires through the function that open names in the defect of a function", func(t *testing.T) {
			t.Parallel()

			load := faults["Load"]
			d := defectNamed(t, &load, "LoadOpensWithoutClosing")
			expect.Equal(t, d.Acquire, "Connect", "the defect of Load acquires through Connect")
			expect.Equal(t, d.Type, "", "the defect of a function has no type")
		})

		t.Run("returns the error of a cancelled context before the break of its callable", func(t *testing.T) {
			t.Parallel()

			s := family("Store")
			misreads := defectNamed(t, s, "StoreGetReadsAnotherValue")
			expect.Equal(t, misreads.Context, "p0", "the defect knows the parameter of the context")
			expect.True(t, misreads.Guards, "the defect returns the error of its context first")
			expect.False(t, defectNamed(t, s, "StoreGetPanics").Guards, "a panic needs no guard")
			expect.False(t, defectNamed(t, s, "StoreGetValueBeside").Guards,
				"a value beside an error reads the error of the subject")
			ignores := defectNamed(t, s, "StoreGetIgnoresCancel")
			expect.False(t, ignores.Guards, "the defect of cancel ignores the context")
			expect.Equal(t, ignores.Context, "p0", "the defect of cancel reads its context")
			expect.False(t, defectNamed(t, s, "StoreLoadDrifts").Guards,
				"a callable without an error has no error to return")
			expect.Equal(t, defectNamed(t, family("Plain"), "PlainSetPanics").Context, "",
				"a callable without a context has no context")
		})

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
			assert.Length(t, all.Results[0].Sample.Stream, 1, "the sample of a stream is an empty stream of one type")
			expect.Equal(t, spelled(all.Results[0].Sample.Stream[0]), "Value", "the stream yields values")
		})

		t.Run("returns a sample that differs from the zero value", func(t *testing.T) {
			t.Parallel()

			d := defectNamed(t, family("Store"), "StoreLoadDrifts")
			assert.NotNil(t, d.Results[0].Sample, "the value of Load has a sample")
			expect.Equal(t, d.Results[1].Role, "flag", "the second result of Load is its flag")
			expect.Nil(t, d.Results[1].Sample, "a flag has no sample")
		})

		samples := []struct {
			name     string
			give     string
			wantText string
		}{
			{
				name:     "returns the sample of a struct whose field differs from its zero value",
				give:     "Has",
				wantText: "true",
			},
			{
				name:     "returns the sample of a struct whose converted field differs from its zero value",
				give:     "State",
				wantText: "true",
			},
			{name: "returns the sample of a bool type that differs from false", give: "IsPresent", wantText: "true"},
			{name: "returns the alternate sample of a bool type that is true", give: "IsEnabled", wantText: "true"},
			{name: "returns the alternate sample of a list whose element is false", give: "Toggles", wantText: "false"},
			{name: "returns the sample of an int type that differs from 0", give: "Amount", wantText: "5"},
			{name: "returns the sample of a float type that differs from 0", give: "Ratio", wantText: "0.5"},
			{
				name:     "returns the sample of a string type that differs from the empty string",
				give:     "Code",
				wantText: "code",
			},
		}
		for _, tt := range samples {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				d := defectNamed(t, family("Lookups"), "Lookups"+tt.give+"ValueBeside")
				sample := d.Results[0].Sample
				assert.NotNil(t, sample, "the value of "+tt.give+" has a sample")
				assert.NotNil(t, sample.Value, "the sample is a value")
				inner := sample.Value
				for inner.Inner != nil || len(inner.Fields) > 0 {
					if inner.Inner == nil {
						inner = &inner.Fields[0].Value
						continue
					}
					inner = inner.Inner
				}
				expect.Equal(t, inner.Text, tt.wantText, "the defect returns the sample that is not the zero value")
			})
		}

		t.Run("races the callables of a declaration through the yield of its handle", func(t *testing.T) {
			t.Parallel()

			board := generate(t, goTree, locksPackage).suites["Board"]
			d := defectNamed(t, &board, "BoardRacyReadWrite")
			expect.Equal(t, []any{d.Type, d.Keeps}, []any{"BoardRacyReadWriteDefect", true},
				"the defect of the declaration keeps the handle whose yield it calls")
			overrides := make([]string, 0, len(d.Overrides))
			for _, o := range d.Overrides {
				overrides = append(overrides, o.Method.Name+" "+o.Body)
			}
			expect.Equal(t, overrides, []string{"Post mark", "Read marked-sample"},
				"the first callable with a value result returns its sample, and every other callable marks its call")
			read := d.Overrides[1]
			expect.Equal(t, []any{read.Context, read.Guards, read.Delegates}, []any{"p0", true, true},
				"the read returns the error of a cancelled context first, and calls the subject otherwise")
			assert.Length(t, read.Results, 2, "the read returns its value and its error")
			expect.NotNil(t, read.Results[0].Sample, "the value of the read has a sample")
		})

		t.Run("keeps the results of the subject in a field and yields before it returns them", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, locksPackage).suites["Shelf"]
			d := defectNamed(t, &shelf, "ShelfCountRacyRead")
			expect.Equal(t, []any{d.Body, d.Keeps, d.Delegates}, []any{"shared-result", true, true},
				"the defect calls the subject and keeps the handle whose yield it calls")
		})

		t.Run("marks each call and yields before it grants a call that finds another call marked", func(t *testing.T) {
			t.Parallel()

			locker := generate(t, goTree, locksPackage).suites["Locker"]
			d := defectNamed(t, &locker, "LockerLockCheckThenGrant")
			expect.Equal(t, []any{d.Body, d.Keeps, d.Delegates, d.Guards}, []any{"check-then-grant", true, true, true},
				"the defect keeps the handle whose yield it calls, and calls the subject one call at a time")
		})

		t.Run("calls a fresh subject of its handle for each call of one key", func(t *testing.T) {
			t.Parallel()

			flight := generate(t, goTree, locksPackage).suites["Flight"]
			d := defectNamed(t, &flight, "FlightDoComputesPerCall")
			expect.Equal(t, []any{d.Body, d.Keeps, d.Delegates}, []any{"fresh-each-call", true, false},
				"the defect keeps the handle whose fresh subject each call calls")
		})

		t.Run("records the method next of the live handle of an opened cursor", func(t *testing.T) {
			t.Parallel()

			live := defectNamed(t, family("Opener"), "OpenerOpenOpensALiveHandle")
			assert.NotNil(t, live.Handle, "the defect returns a live handle")
			expect.Equal(t, live.Handle.Method.Name, "Next", "the handle overrides the method that next-method names")
			expect.Length(t, live.Handle.Results, 3, "the handle returns the value, the flag and the error of Next")
			expect.Equal(t, live.Handle.Body, "zero", "the method of the handle returns its zero results")
		})

		vault := generate(t, goTree, vaultPackage).suites["Vault"]
		durables := []struct {
			name       string
			give       string
			wantType   string
			wantBody   string
			wantMedium string
		}{
			{
				name: "opens the defect of acknowledges-before-sync over an unsynced medium",
				give: "VaultPutAcknowledgesBeforeSync", wantMedium: "Unsynced",
			},
			{
				name: "opens the defect of hides-a-full-medium over an unfilled medium",
				give: "VaultPutHidesAFullMedium", wantMedium: "Unfilled",
			},
			{
				name:     "writes the zero value before the first value of trusts-a-partial-record",
				give:     "VaultPutTrustsAPartialRecord",
				wantType: "VaultPutTrustsAPartialRecordDefect",
				wantBody: "zero-first",
			},
			{
				name:     "swallows the declared error of a corrupted record in ignores-checksum",
				give:     "VaultGetIgnoresChecksum",
				wantType: "VaultGetIgnoresChecksumDefect",
				wantBody: "zero-on-sentinel",
			},
		}
		for _, tt := range durables {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				d := defectNamed(t, &vault, tt.give)
				expect.Equal(
					t,
					[]any{d.Type, d.Body, d.Medium, d.Door},
					[]any{tt.wantType, tt.wantBody, tt.wantMedium, tt.give + "Open"},
					"the defect of a durable claim opens over a medium through its door",
				)
			})
		}

		t.Run("replaces the written value of the first write with its zero value", func(t *testing.T) {
			t.Parallel()

			d := defectNamed(t, &vault, "VaultPutTrustsAPartialRecord")
			expect.Equal(t, []any{d.Input, d.Delegates}, []any{2, true},
				"the defect calls Put with the zero value of its third parameter first")
		})

		t.Run("swallows the error that the param corrupt of the mixin durable refers to", func(t *testing.T) {
			t.Parallel()

			d := defectNamed(t, &vault, "VaultGetIgnoresChecksum")
			assert.NotNil(t, d.Sentinel, "the defect has the declared error of a corrupted record")
			assert.NotNil(t, d.Sentinel.Value, "the declared error is a value")
			expect.Equal(
				t,
				[]any{d.Sentinel.Value.Kind, d.Sentinel.Value.Callee.Name},
				[]any{emit.ValueName, "ErrCorrupt"},
				"the declared error is the name ErrCorrupt",
			)
			expect.True(t, d.Delegates, "the defect calls Get of the subject")
		})
	})
}
