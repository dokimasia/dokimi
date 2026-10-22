// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/directive"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/position"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The directive budget states the ceilings of a callable on its method or its
// function, or on a subject, and suite refuses a budget that no benchmark can
// apply.
func TestBudget(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("declares the directive budget for three kinds of declaration", func(t *testing.T) {
			t.Parallel()

			declared, declares := suite.New().(plugin.DirectiveProvider)
			assert.True(t, declares, "suite declares its directives")
			schemas := declared.Directives()
			i := slices.IndexFunc(schemas, func(s directive.Schema) bool { return s.Name == suite.BudgetDirective })
			assert.True(t, i >= 0, "suite declares the directive budget")
			expect.Equal(t, schemas[i].Kinds, []symbol.Kind{symbol.KindMethod, symbol.KindFunction, symbol.KindStruct},
				"a budget is on a callable or on a subject")
			expect.True(t, schemas[i].Repeatable, "a subject states a budget of each callable")
		})

		tests := []struct {
			name     string
			givePkg  string
			giveKind symbol.Kind
			giveOwn  string
			give     string
			want     ceilings
		}{
			{
				name: "stamps the ceilings of a budget on its method", givePkg: meterPackage,
				giveKind: symbol.KindMethod, giveOwn: "Meter", give: "Count",
				want: ceilings{
					allocs: 1, allocsSet: true, bytes: 96, bytesSet: true, latency: 2 * time.Microsecond,
					mean: time.Microsecond,
				},
			},
			{
				name: "stamps the ceilings of a budget on a method of a struct", givePkg: catalogPackage,
				giveKind: symbol.KindMethod, giveOwn: "Catalog", give: "Find", want: ceilings{allocsSet: true},
			},
			{
				name: "stamps the ceilings of a budget on a function", givePkg: meterPackage,
				giveKind: symbol.KindFunction, give: "Weigh",
				want: ceilings{allocsSet: true, latency: 100 * time.Nanosecond},
			},
			{
				name: "stamps no ceiling that a budget leaves out", givePkg: meterPackage, giveKind: symbol.KindMethod,
				giveOwn: "Meter", give: "Record",
				want: ceilings{latency: time.Millisecond, warmup: 500, warmupSet: true},
			},
			{
				name: "stamps the warm-up of a budget without a ceiling", givePkg: badBudgetPackage,
				giveKind: symbol.KindMethod, giveOwn: "Gauge", give: "Warm",
				want: ceilings{warmup: 10, warmupSet: true},
			},
			{
				name: "stamps nothing on a callable without a budget", givePkg: meterPackage,
				giveKind: symbol.KindMethod, giveOwn: "Meter", give: "Total",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				f := fixture(t, goTree, p, tt.givePkg)
				phases(t, f, p)
				expect.Equal(t, stamped(t, f, f.Find(t, tt.giveKind, tt.giveOwn, tt.give).Identity()), tt.want,
					"the facts contain the ceilings of the budget")
			})
		}

		findings := []struct {
			name     string
			giveKind symbol.Kind
			giveOwn  string
			give     string
			wantMsg  []string
		}{
			{
				name: "reports BadBudget at a budget without a ceiling", giveKind: symbol.KindMethod, giveOwn: "Gauge",
				give: "Empty",
				wantMsg: []string{
					"the budget of Gauge.Empty states no ceiling and no warm-up. State allocs, bytes, p99, mean or warmup",
				},
			},
			{
				name: "reports BadBudget at a warm-up below 0", giveKind: symbol.KindMethod, giveOwn: "Gauge",
				give:    "Cold",
				wantMsg: []string{"the budget of Gauge.Cold states warmup=-1, and a warm-up is 0 iterations or more"},
			},
			{
				name: "reports BadBudget at a negative count of allocations", giveKind: symbol.KindMethod,
				giveOwn: "Gauge", give: "Negative",
				wantMsg: []string{"the budget of Gauge.Negative states allocs=-1, and a ceiling of allocations is 0 " +
					"or more"},
			},
			{
				name: "reports BadBudget at a negative count of bytes", giveKind: symbol.KindMethod, giveOwn: "Gauge",
				give:    "Shrinking",
				wantMsg: []string{"the budget of Gauge.Shrinking states bytes=-2, and a ceiling of bytes is 0 or more"},
			},
			{
				name: "reports BadBudget at a p99 that is no duration", giveKind: symbol.KindMethod, giveOwn: "Gauge",
				give:    "Fast",
				wantMsg: []string{`the budget of Gauge.Fast states p99="fast", and a p99 is a duration such as 2us`},
			},
			{
				name: "reports BadBudget at a p99 of no time", giveKind: symbol.KindMethod, giveOwn: "Gauge",
				give:    "Instant",
				wantMsg: []string{`the budget of Gauge.Instant states p99="0s", and a p99 is above 0`},
			},
			{
				name: "reports BadBudget at a mean that is no duration", giveKind: symbol.KindMethod, giveOwn: "Gauge",
				give: "Steady",
				wantMsg: []string{
					`the budget of Gauge.Steady states mean="steady", and a mean is a duration such as 2us`,
				},
			},
			{
				name: "reports BadBudget at a mean without a p99", giveKind: symbol.KindMethod, giveOwn: "Gauge",
				give: "Average",
				wantMsg: []string{"the budget of Gauge.Average states a mean and no p99, and a mean alone passes a " +
					"tail that grows"},
			},
			{
				name: "reports BadBudget at a budget of a method that names a callable", giveKind: symbol.KindMethod,
				giveOwn: "Gauge", give: "Named",
				wantMsg: []string{"the budget of Gauge.Named names the callable Gauge.Named, and only the budget of " +
					"a subject names its callable"},
			},
			{
				name: "reports BadBudget at each of two budgets of one method", giveKind: symbol.KindMethod,
				giveOwn: "Gauge", give: "Twice",
				wantMsg: []string{
					"Gauge.Twice states 2 budgets, and a callable has one budget",
					"Gauge.Twice states 2 budgets, and a callable has one budget",
				},
			},
			{
				name: "reports BadBudget at a budget of a method that is not exported", giveKind: symbol.KindMethod,
				giveOwn: "Gauge", give: "hidden",
				wantMsg: []string{"Gauge.hidden is not exported, so its budget has no check"},
			},
			{
				name: "reports BadBudget at a budget of a declaration without a suite", giveKind: symbol.KindMethod,
				giveOwn: "Unsuited", give: "Get",
				wantMsg: []string{"Unsuited has no directive suite, so the budget of Unsuited.Get has no check"},
			},
			{
				name: "reports BadBudget at a budget of a function without a suite", giveKind: symbol.KindFunction,
				give:    "Loose",
				wantMsg: []string{"Loose has no directive suite, so the budget of Loose has no check"},
			},
		}
		for _, tt := range findings {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				f := fixture(t, goTree, p, meterPackage, badBudgetPackage)
				stamps, _ := phases(t, f, p)
				decl := f.Find(t, tt.giveKind, tt.giveOwn, tt.give)
				expect.Equal(t, messages(stamps.Sink.All(), decl.Position()), tt.wantMsg,
					"suite reports each budget that it refuses above the callable")
				expect.Equal(t, stamped(t, f, decl.Identity()), ceilings{}, "suite stamps nothing that it refuses")
			})
		}

		t.Run("records the ceilings of a subject for its conformance test", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, meterPackage, meteredPackage).tests["TestMeterConformsToMeter"]
			assert.True(t, found, "suite records the conformance test of Meter")
			expect.Equal(t, got.Data.Ceilings, []suite.Ceiling{{
				Field: "Merge", Budget: suite.Budget{Allocs: new(int64(0)), Latency: time.Millisecond},
			}}, "the test applies the budget of the subject to Merge")
		})

		t.Run("records the warm-up of a subject for its conformance test", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, meterPackage, meteredPackage).tests["TestWarmedConformsToMeter"]
			assert.True(t, found, "suite records the conformance test of Warmed")
			expect.Equal(t, got.Data.Ceilings, []suite.Ceiling{{
				Field: "Merge", Budget: suite.Budget{Warmup: new(int64(300))},
			}}, "the test applies the warm-up of the subject to the benchmarks of Merge")
		})

		t.Run("records no ceiling of a subject without a budget", func(t *testing.T) {
			t.Parallel()

			got, found := generate(t, goTree, meterPackage, meteredPackage).tests["TestPlainConformsToMeter"]
			assert.True(t, found, "suite records the conformance test of Plain")
			expect.Empty(t, got.Data.Ceilings, "the test applies the budgets of the declaration alone")
		})

		subjects := []struct {
			name    string
			give    string
			wantMsg []string
		}{
			{
				name: "reports BadBudget at a budget of a subject without a callable", give: "Unnamed",
				wantMsg: []string{"the budget of Unnamed names no callable of Meter, so it has no check. " +
					"The callables of Meter are Count, Record, Merge and Total"},
			},
			{
				name: "reports BadBudget at a budget of a callable that the declaration lacks", give: "Lost",
				wantMsg: []string{"Meter.Missing is no callable of Meter, so the budget of Lost has no check. " +
					"The callables of Meter are Count, Record, Merge and Total"},
			},
			{
				name: "reports BadBudget at a budget of a callable whose declaration states one", give: "Doubled",
				wantMsg: []string{"Doubled states a budget of Meter.Count, whose declaration states a budget. " +
					"Excuse the checks of the declaration instead"},
			},
			{
				name:    "reports BadBudget at a budget of a callable whose declaration states a warm-up alone",
				give:    "Rewarmed",
				wantMsg: []string{"Rewarmed states a budget of Gauge.Warm, whose declaration states a budget"},
			},
			{
				name: "reports BadBudget at each of two budgets of one callable of a subject", give: "Repeated",
				wantMsg: []string{
					"Repeated states 2 budgets of Meter.Total, and a callable has one budget",
					"Repeated states 2 budgets of Meter.Total, and a callable has one budget",
				},
			},
			{
				name: "reports BadBudget at a ceiling of a subject that is no duration", give: "Invalid",
				wantMsg: []string{
					`the budget of Meter.Total states p99="fast", and a p99 is a duration such as 2us`,
				},
			},
			{
				name: "reports BadBudget at a budget of a struct that is no subject", give: "Unbound",
				wantMsg: []string{"Unbound has no directive subject, so its budget has no check"},
			},
		}
		for _, tt := range subjects {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				f := fixture(t, goTree, p, meterPackage, badBudgetPackage)
				_, res := phases(t, f, p)
				subject := f.Find(t, symbol.KindStruct, "", tt.give)
				expect.Equal(t, messages(res.Sink.All(), subject.Position()), tt.wantMsg,
					"suite reports each budget that it refuses above the subject")
				gentest.Unrecorded(t, res, subject)
			})
		}
	})
}

// ceilings are the ceilings and the warm-up that the facts of suite contain for
// a callable. allocsSet, bytesSet and warmupSet report that a fact of the count
// exists.
type ceilings struct {
	allocs, bytes, warmup          int64
	latency, mean                  time.Duration
	allocsSet, bytesSet, warmupSet bool
}

// stamped returns the ceilings that the facts of the fixture f contain for the
// callable id.
func stamped(t *testing.T, f *gentest.Fixture, id symbol.Identity) ceilings {
	t.Helper()

	fact := func(key meta.KeyName) (int64, bool) {
		k, registered := meta.Lookup[int64](f.Keys, key)
		assert.True(t, registered, "suite registers the key "+string(key))
		return meta.Get(f.Facts, id, k)
	}
	var c ceilings
	c.allocs, c.allocsSet = fact(suite.KeyAllocs)
	c.bytes, c.bytesSet = fact(suite.KeyBytes)
	c.warmup, c.warmupSet = fact(suite.KeyWarmup)
	latency, _ := fact(suite.KeyLatency)
	mean, _ := fact(suite.KeyMean)
	c.latency, c.mean = time.Duration(latency), time.Duration(mean)
	return c
}

// messages returns the message of each finding of all on the lines of the
// directives above the declaration at, in the order of their lines. The
// directives of a declaration are on the three lines before it.
func messages(all func(yield func(diag.Diag) bool), at position.Pos) []string {
	var found []diag.Diag
	for d := range all {
		if d.Pos.File == at.File && d.Pos.Line < at.Line && d.Pos.Line >= at.Line-3 {
			found = append(found, d)
		}
	}
	slices.SortStableFunc(found, func(a, b diag.Diag) int { return a.Pos.Line - b.Pos.Line })
	out := make([]string, 0, len(found))
	for _, d := range found {
		out = append(out, d.Msg)
	}
	return out
}
