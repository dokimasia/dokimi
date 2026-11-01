package drawn_test

import (
	"testing"

	"example.com/shop/drawn"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// talliesSubjectEmbedded is drawn.Tallies under the name of the field that
// talliesSubject embeds.
type talliesSubjectEmbedded = drawn.Tallies

// talliesSubject is a value of Tallies that only TestTalliesSuite builds.
type talliesSubject struct{ talliesSubjectEmbedded }

// TestTalliesSuite checks the suite of Tallies and the two forms of its
// application.
func TestTalliesSuite(t *testing.T) {
	t.Parallel()

	t.Run("TalliesSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, drawn.TalliesSuite())
	})

	t.Run("RunTalliesSuite", func(t *testing.T) {
		t.Parallel()
		var built *talliesSubject
		drawn.RunTalliesSuite(
			t,
			func() *talliesSubject {
				built = &talliesSubject{}
				return built
			},
			conform.With(drawn.TalliesCheck{
				ID:        "Tallies/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.TalliesHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Tallies(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.TalliesHandle) drawn.Tallies { return &talliesSubject{} },
			}),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.DeterministicGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Commutative(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Associative(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.DeterministicGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.Commutative(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Subscribe.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Subscribe.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTalliesSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *talliesSubject
		drawn.RunTalliesSuiteFromHandle(
			t,
			func(*drawn.TalliesHandle) *talliesSubject {
				built = &talliesSubject{}
				return built
			},
			conform.With(drawn.TalliesCheck{
				ID:        "Tallies/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.TalliesHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Tallies(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.TalliesHandle) drawn.Tallies { return &talliesSubject{} },
			}),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.DeterministicGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Commutative(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Merge.Associative(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.DeterministicGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Distance.Commutative(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Each.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Save.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Limit.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Visit.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Schedule.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Link.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Subscribe.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.TalliesChecks.Subscribe.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// talliesIdle returns the planted defect of a ceiling of Tallies with a
// spend of nothing, a subject that answers each call at once.
func talliesIdle(h *drawn.TalliesHandle) drawn.Tallies {
	return drawn.TalliesSuite().Calls[0].Over(h, func() {})
}

// TestTalliesSuiteBudgets applies the allocation ceilings of the suite of
// Tallies to talliesIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestTalliesSuiteBudgets(t *testing.T) {
	drawn.RunTalliesBudgets(
		t,
		talliesIdle,
		conform.Ceilings(
			drawn.TalliesChecks.Merge.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.TalliesChecks.Distance.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(drawn.TalliesChecks.Each.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(drawn.TalliesChecks.Save.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			drawn.TalliesChecks.Limit.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.TalliesChecks.Count.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.TalliesChecks.Visit.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.TalliesChecks.Schedule.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(drawn.TalliesChecks.Link.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			drawn.TalliesChecks.Subscribe.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteMerge registers the fuzz target of Tallies.Merge on
// talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteMerge(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Merge.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteDistance registers the fuzz target of Tallies.Distance
// on talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteDistance(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Distance.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteEach registers the fuzz target of Tallies.Each on
// talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteEach(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Each.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteSave registers the fuzz target of Tallies.Save on
// talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteSave(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Save.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteLimit registers the fuzz target of Tallies.Limit on
// talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteLimit(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Limit.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteCount registers the fuzz target of Tallies.Count on
// talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteCount(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Count.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteVisit registers the fuzz target of Tallies.Visit on
// talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteVisit(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Visit.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteSchedule registers the fuzz target of Tallies.Schedule
// on talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteSchedule(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Schedule.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteLink registers the fuzz target of Tallies.Link on
// talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteLink(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Link.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTalliesSuiteSubscribe registers the fuzz target of Tallies.Subscribe
// on talliesIdle, with every check of the suite excused.
func FuzzTalliesSuiteSubscribe(f *testing.F) {
	drawn.FuzzTalliesSuite(
		f,
		drawn.TalliesChecks.Subscribe.Callable,
		talliesIdle,
		conform.Excuse(
			drawn.TalliesChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Merge.Associative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.DeterministicGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Distance.Commutative(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Each.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Save.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Limit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Visit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Schedule.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Link.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.TalliesChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkTalliesSuite measures each callable of the suite of Tallies on
// talliesIdle.
func BenchmarkTalliesSuite(b *testing.B) {
	b.Run("Merge", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Merge.Callable, talliesIdle)
	})
	b.Run("Distance", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Distance.Callable, talliesIdle)
	})
	b.Run("Each", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Each.Callable, talliesIdle)
	})
	b.Run("Save", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Save.Callable, talliesIdle)
	})
	b.Run("Limit", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Limit.Callable, talliesIdle)
	})
	b.Run("Count", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Count.Callable, talliesIdle)
	})
	b.Run("Visit", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Visit.Callable, talliesIdle)
	})
	b.Run("Schedule", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Schedule.Callable, talliesIdle)
	})
	b.Run("Link", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Link.Callable, talliesIdle)
	})
	b.Run("Subscribe", func(b *testing.B) {
		drawn.BenchTalliesSuite(b, drawn.TalliesChecks.Subscribe.Callable, talliesIdle)
	})
}

// notesSubjectEmbedded is drawn.Notes under the name of the field that
// notesSubject embeds.
type notesSubjectEmbedded = drawn.Notes

// notesSubject is a value of Notes that only TestNotesSuite builds.
type notesSubject struct{ notesSubjectEmbedded }

// TestNotesSuite checks the suite of Notes and the two forms of its
// application.
func TestNotesSuite(t *testing.T) {
	t.Parallel()

	t.Run("NotesSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, drawn.NotesSuite())
	})

	t.Run("RunNotesSuite", func(t *testing.T) {
		t.Parallel()
		var built *notesSubject
		drawn.RunNotesSuite(
			t,
			func() *notesSubject {
				built = &notesSubject{}
				return built
			},
			conform.With(drawn.NotesCheck{
				ID:        "Notes/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.NotesHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Notes(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.NotesHandle) drawn.Notes { return &notesSubject{} },
			}),
			conform.Excuse(
				drawn.NotesChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.UpdaterGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunNotesSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *notesSubject
		drawn.RunNotesSuiteFromHandle(
			t,
			func(*drawn.NotesHandle) *notesSubject {
				built = &notesSubject{}
				return built
			},
			conform.With(drawn.NotesCheck{
				ID:        "Notes/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.NotesHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Notes(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.NotesHandle) drawn.Notes { return &notesSubject{} },
			}),
			conform.Excuse(
				drawn.NotesChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Put.UpdaterGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.NotesChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// notesIdle returns the planted defect of a ceiling of Notes with a spend
// of nothing, a subject that answers each call at once.
func notesIdle(h *drawn.NotesHandle) drawn.Notes {
	return drawn.NotesSuite().Calls[0].Over(h, func() {})
}

// TestNotesSuiteBudgets applies the allocation ceilings of the suite of
// Notes to notesIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestNotesSuiteBudgets(t *testing.T) {
	drawn.RunNotesBudgets(
		t,
		notesIdle,
		conform.Ceilings(drawn.NotesChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(drawn.NotesChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			drawn.NotesChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.UpdaterGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuitePut registers the fuzz target of Notes.Put on notesIdle,
// with every check of the suite excused.
func FuzzNotesSuitePut(f *testing.F) {
	drawn.FuzzNotesSuite(
		f,
		drawn.NotesChecks.Put.Callable,
		notesIdle,
		conform.Excuse(
			drawn.NotesChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.UpdaterGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuiteGet registers the fuzz target of Notes.Get on notesIdle,
// with every check of the suite excused.
func FuzzNotesSuiteGet(f *testing.F) {
	drawn.FuzzNotesSuite(
		f,
		drawn.NotesChecks.Get.Callable,
		notesIdle,
		conform.Excuse(
			drawn.NotesChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Put.UpdaterGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.NotesChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkNotesSuite measures each callable of the suite of Notes on
// notesIdle.
func BenchmarkNotesSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		drawn.BenchNotesSuite(b, drawn.NotesChecks.Put.Callable, notesIdle)
	})
	b.Run("Get", func(b *testing.B) {
		drawn.BenchNotesSuite(b, drawn.NotesChecks.Get.Callable, notesIdle)
	})
}

// hookedSubjectEmbedded is drawn.Hooked under the name of the field that
// hookedSubject embeds.
type hookedSubjectEmbedded = drawn.Hooked

// hookedSubject is a value of Hooked that only TestHookedSuite builds.
type hookedSubject struct{ hookedSubjectEmbedded }

// TestHookedSuite checks the suite of Hooked and the two forms of its
// application.
func TestHookedSuite(t *testing.T) {
	t.Parallel()

	t.Run("HookedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, drawn.HookedSuite())
	})

	t.Run("RunHookedSuite", func(t *testing.T) {
		t.Parallel()
		var built *hookedSubject
		drawn.RunHookedSuite(
			t,
			func() *hookedSubject {
				built = &hookedSubject{}
				return built
			},
			conform.With(drawn.HookedCheck{
				ID:        "Hooked/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.HookedHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Hooked(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.HookedHandle) drawn.Hooked { return &hookedSubject{} },
			}),
			conform.Excuse(
				drawn.HookedChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.PartitionGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunHookedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *hookedSubject
		drawn.RunHookedSuiteFromHandle(
			t,
			func(*drawn.HookedHandle) *hookedSubject {
				built = &hookedSubject{}
				return built
			},
			conform.With(drawn.HookedCheck{
				ID:        "Hooked/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.HookedHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Hooked(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.HookedHandle) drawn.Hooked { return &hookedSubject{} },
			}),
			conform.Excuse(
				drawn.HookedChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Set.PartitionGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.HookedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// hookedIdle returns the planted defect of a ceiling of Hooked with a
// spend of nothing, a subject that answers each call at once.
func hookedIdle(h *drawn.HookedHandle) drawn.Hooked {
	return drawn.HookedSuite().Calls[0].Over(h, func() {})
}

// TestHookedSuiteBudgets applies the allocation ceilings of the suite of
// Hooked to hookedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestHookedSuiteBudgets(t *testing.T) {
	drawn.RunHookedBudgets(
		t,
		hookedIdle,
		conform.Ceilings(drawn.HookedChecks.Set.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(drawn.HookedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			drawn.HookedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.PartitionGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHookedSuiteSet registers the fuzz target of Hooked.Set on
// hookedIdle, with every check of the suite excused.
func FuzzHookedSuiteSet(f *testing.F) {
	drawn.FuzzHookedSuite(
		f,
		drawn.HookedChecks.Set.Callable,
		hookedIdle,
		conform.Excuse(
			drawn.HookedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.PartitionGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHookedSuiteGet registers the fuzz target of Hooked.Get on
// hookedIdle, with every check of the suite excused.
func FuzzHookedSuiteGet(f *testing.F) {
	drawn.FuzzHookedSuite(
		f,
		drawn.HookedChecks.Get.Callable,
		hookedIdle,
		conform.Excuse(
			drawn.HookedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Set.PartitionGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.HookedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkHookedSuite measures each callable of the suite of Hooked on
// hookedIdle.
func BenchmarkHookedSuite(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		drawn.BenchHookedSuite(b, drawn.HookedChecks.Set.Callable, hookedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		drawn.BenchHookedSuite(b, drawn.HookedChecks.Get.Callable, hookedIdle)
	})
}

// lookupsSubjectEmbedded is drawn.Lookups under the name of the field that
// lookupsSubject embeds.
type lookupsSubjectEmbedded = drawn.Lookups

// lookupsSubject is a value of Lookups that only TestLookupsSuite builds.
type lookupsSubject struct{ lookupsSubjectEmbedded }

// TestLookupsSuite checks the suite of Lookups and the two forms of its
// application.
func TestLookupsSuite(t *testing.T) {
	t.Parallel()

	t.Run("LookupsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, drawn.LookupsSuite())
	})

	t.Run("RunLookupsSuite", func(t *testing.T) {
		t.Parallel()
		var built *lookupsSubject
		drawn.RunLookupsSuite(
			t,
			func() *lookupsSubject {
				built = &lookupsSubject{}
				return built
			},
			conform.With(drawn.LookupsCheck{
				ID:        "Lookups/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.LookupsHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Lookups(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.LookupsHandle) drawn.Lookups { return &lookupsSubject{} },
			}),
			conform.Excuse(
				drawn.LookupsChecks.Has.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLookupsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *lookupsSubject
		drawn.RunLookupsSuiteFromHandle(
			t,
			func(*drawn.LookupsHandle) *lookupsSubject {
				built = &lookupsSubject{}
				return built
			},
			conform.With(drawn.LookupsCheck{
				ID:        "Lookups/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.LookupsHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Lookups(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.LookupsHandle) drawn.Lookups { return &lookupsSubject{} },
			}),
			conform.Excuse(
				drawn.LookupsChecks.Has.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// lookupsIdle returns the planted defect of a ceiling of Lookups with a
// spend of nothing, a subject that answers each call at once.
func lookupsIdle(h *drawn.LookupsHandle) drawn.Lookups {
	return drawn.LookupsSuite().Calls[0].Over(h, func() {})
}

// TestLookupsSuiteBudgets applies the allocation ceilings of the suite of
// Lookups to lookupsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLookupsSuiteBudgets(t *testing.T) {
	drawn.RunLookupsBudgets(
		t,
		lookupsIdle,
		conform.Ceilings(drawn.LookupsChecks.Has.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			drawn.LookupsChecks.IsPresent.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.LookupsChecks.IsEnabled.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.LookupsChecks.State.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.LookupsChecks.Toggles.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.LookupsChecks.Amount.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			drawn.LookupsChecks.Ratio.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(drawn.LookupsChecks.Code.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteHas registers the fuzz target of Lookups.Has on
// lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteHas(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.Has.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteIsPresent registers the fuzz target of Lookups.IsPresent
// on lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteIsPresent(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.IsPresent.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteIsEnabled registers the fuzz target of Lookups.IsEnabled
// on lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteIsEnabled(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.IsEnabled.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteState registers the fuzz target of Lookups.State on
// lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteState(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.State.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteToggles registers the fuzz target of Lookups.Toggles on
// lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteToggles(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.Toggles.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteAmount registers the fuzz target of Lookups.Amount on
// lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteAmount(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.Amount.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteRatio registers the fuzz target of Lookups.Ratio on
// lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteRatio(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.Ratio.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLookupsSuiteCode registers the fuzz target of Lookups.Code on
// lookupsIdle, with every check of the suite excused.
func FuzzLookupsSuiteCode(f *testing.F) {
	drawn.FuzzLookupsSuite(
		f,
		drawn.LookupsChecks.Code.Callable,
		lookupsIdle,
		conform.Excuse(
			drawn.LookupsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Has.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsPresent.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.IsEnabled.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.State.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Toggles.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Amount.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Ratio.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LookupsChecks.Code.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLookupsSuite measures each callable of the suite of Lookups on
// lookupsIdle.
func BenchmarkLookupsSuite(b *testing.B) {
	b.Run("Has", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.Has.Callable, lookupsIdle)
	})
	b.Run("IsPresent", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.IsPresent.Callable, lookupsIdle)
	})
	b.Run("IsEnabled", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.IsEnabled.Callable, lookupsIdle)
	})
	b.Run("State", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.State.Callable, lookupsIdle)
	})
	b.Run("Toggles", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.Toggles.Callable, lookupsIdle)
	})
	b.Run("Amount", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.Amount.Callable, lookupsIdle)
	})
	b.Run("Ratio", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.Ratio.Callable, lookupsIdle)
	})
	b.Run("Code", func(b *testing.B) {
		drawn.BenchLookupsSuite(b, drawn.LookupsChecks.Code.Callable, lookupsIdle)
	})
}

// labelsSubjectEmbedded is drawn.Labels under the name of the field that
// labelsSubject embeds.
type labelsSubjectEmbedded = drawn.Labels

// labelsSubject is a value of Labels that only TestLabelsSuite builds.
type labelsSubject struct{ labelsSubjectEmbedded }

// TestLabelsSuite checks the suite of Labels and the two forms of its
// application.
func TestLabelsSuite(t *testing.T) {
	t.Parallel()

	t.Run("LabelsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, drawn.LabelsSuite())
	})

	t.Run("RunLabelsSuite", func(t *testing.T) {
		t.Parallel()
		var built *labelsSubject
		drawn.RunLabelsSuite(
			t,
			func() *labelsSubject {
				built = &labelsSubject{}
				return built
			},
			conform.With(drawn.LabelsCheck{
				ID:        "Labels/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.LabelsHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Labels(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.LabelsHandle) drawn.Labels { return &labelsSubject{} },
			}),
			conform.Excuse(
				drawn.LabelsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLabelsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *labelsSubject
		drawn.RunLabelsSuiteFromHandle(
			t,
			func(*drawn.LabelsHandle) *labelsSubject {
				built = &labelsSubject{}
				return built
			},
			conform.With(drawn.LabelsCheck{
				ID:        "Labels/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *drawn.LabelsHandle) {
					got := h.Subject()
					assert.Equal(h, got, drawn.Labels(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*drawn.LabelsHandle) drawn.Labels { return &labelsSubject{} },
			}),
			conform.Excuse(
				drawn.LabelsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				drawn.LabelsChecks.Get.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// labelsIdle returns the planted defect of a ceiling of Labels with a
// spend of nothing, a subject that answers each call at once.
func labelsIdle(h *drawn.LabelsHandle) drawn.Labels {
	return drawn.LabelsSuite().Calls[0].Over(h, func() {})
}

// TestLabelsSuiteBudgets applies the allocation ceilings of the suite of
// Labels to labelsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLabelsSuiteBudgets(t *testing.T) {
	drawn.RunLabelsBudgets(
		t,
		labelsIdle,
		conform.Ceilings(drawn.LabelsChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(drawn.LabelsChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			drawn.LabelsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLabelsSuitePut registers the fuzz target of Labels.Put on
// labelsIdle, with every check of the suite excused.
func FuzzLabelsSuitePut(f *testing.F) {
	drawn.FuzzLabelsSuite(
		f,
		drawn.LabelsChecks.Put.Callable,
		labelsIdle,
		conform.Excuse(
			drawn.LabelsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLabelsSuiteGet registers the fuzz target of Labels.Get on
// labelsIdle, with every check of the suite excused.
func FuzzLabelsSuiteGet(f *testing.F) {
	drawn.FuzzLabelsSuite(
		f,
		drawn.LabelsChecks.Get.Callable,
		labelsIdle,
		conform.Excuse(
			drawn.LabelsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			drawn.LabelsChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLabelsSuite measures each callable of the suite of Labels on
// labelsIdle.
func BenchmarkLabelsSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		drawn.BenchLabelsSuite(b, drawn.LabelsChecks.Put.Callable, labelsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		drawn.BenchLabelsSuite(b, drawn.LabelsChecks.Get.Callable, labelsIdle)
	})
}
