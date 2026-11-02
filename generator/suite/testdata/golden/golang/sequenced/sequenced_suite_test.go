package sequenced_test

import (
	"testing"

	"example.com/shop/sequenced"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// counterSubjectEmbedded is sequenced.Counter under the name of the field
// that counterSubject embeds.
type counterSubjectEmbedded = sequenced.Counter

// counterSubject is a value of Counter that only TestCounterSuite builds.
type counterSubject struct{ counterSubjectEmbedded }

// TestCounterSuite checks the suite of Counter and the two forms of its
// application.
func TestCounterSuite(t *testing.T) {
	t.Parallel()

	t.Run("CounterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.CounterSuite())
	})

	t.Run("RunCounterSuite", func(t *testing.T) {
		t.Parallel()
		var built *counterSubject
		sequenced.RunCounterSuite(
			t,
			func() *counterSubject {
				built = &counterSubject{}
				return built
			},
			conform.With(sequenced.CounterCheck{
				ID:        "Counter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.CounterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Counter(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.CounterHandle) sequenced.Counter { return &counterSubject{} },
			}),
			conform.Excuse(
				sequenced.CounterChecks.Inc.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Inc.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.Monotonic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.NeverGoesBack(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Grow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Grow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Size.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Size.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Size.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Throttle.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Throttle.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Load.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Load.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCounterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *counterSubject
		sequenced.RunCounterSuiteFromHandle(
			t,
			func(*sequenced.CounterHandle) *counterSubject {
				built = &counterSubject{}
				return built
			},
			conform.With(sequenced.CounterCheck{
				ID:        "Counter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.CounterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Counter(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.CounterHandle) sequenced.Counter { return &counterSubject{} },
			}),
			conform.Excuse(
				sequenced.CounterChecks.Inc.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Inc.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.Monotonic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Value.NeverGoesBack(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Grow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Grow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Size.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Size.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Size.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Throttle.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Throttle.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Load.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CounterChecks.Load.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// counterIdle returns the planted defect of a ceiling of Counter with a
// spend of nothing, a subject that answers each call at once.
func counterIdle(h *sequenced.CounterHandle) sequenced.Counter {
	return sequenced.CounterSuite().Calls[0].Over(h, func() {})
}

// TestCounterSuiteBudgets applies the allocation ceilings of the suite of
// Counter to counterIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCounterSuiteBudgets(t *testing.T) {
	sequenced.RunCounterBudgets(
		t,
		counterIdle,
		conform.Ceilings(
			sequenced.CounterChecks.Inc.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CounterChecks.Value.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CounterChecks.Grow.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CounterChecks.Size.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CounterChecks.Throttle.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CounterChecks.Load.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Monotonic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.NeverGoesBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCounterSuiteInc registers the fuzz target of Counter.Inc on
// counterIdle, with every check of the suite excused.
func FuzzCounterSuiteInc(f *testing.F) {
	sequenced.FuzzCounterSuite(
		f,
		sequenced.CounterChecks.Inc.Callable,
		counterIdle,
		conform.Excuse(
			sequenced.CounterChecks.Inc.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Monotonic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.NeverGoesBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCounterSuiteValue registers the fuzz target of Counter.Value on
// counterIdle, with every check of the suite excused.
func FuzzCounterSuiteValue(f *testing.F) {
	sequenced.FuzzCounterSuite(
		f,
		sequenced.CounterChecks.Value.Callable,
		counterIdle,
		conform.Excuse(
			sequenced.CounterChecks.Inc.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Monotonic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.NeverGoesBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCounterSuiteGrow registers the fuzz target of Counter.Grow on
// counterIdle, with every check of the suite excused.
func FuzzCounterSuiteGrow(f *testing.F) {
	sequenced.FuzzCounterSuite(
		f,
		sequenced.CounterChecks.Grow.Callable,
		counterIdle,
		conform.Excuse(
			sequenced.CounterChecks.Inc.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Monotonic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.NeverGoesBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCounterSuiteSize registers the fuzz target of Counter.Size on
// counterIdle, with every check of the suite excused.
func FuzzCounterSuiteSize(f *testing.F) {
	sequenced.FuzzCounterSuite(
		f,
		sequenced.CounterChecks.Size.Callable,
		counterIdle,
		conform.Excuse(
			sequenced.CounterChecks.Inc.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Monotonic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.NeverGoesBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCounterSuiteThrottle registers the fuzz target of Counter.Throttle
// on counterIdle, with every check of the suite excused.
func FuzzCounterSuiteThrottle(f *testing.F) {
	sequenced.FuzzCounterSuite(
		f,
		sequenced.CounterChecks.Throttle.Callable,
		counterIdle,
		conform.Excuse(
			sequenced.CounterChecks.Inc.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Monotonic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.NeverGoesBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCounterSuiteLoad registers the fuzz target of Counter.Load on
// counterIdle, with every check of the suite excused.
func FuzzCounterSuiteLoad(f *testing.F) {
	sequenced.FuzzCounterSuite(
		f,
		sequenced.CounterChecks.Load.Callable,
		counterIdle,
		conform.Excuse(
			sequenced.CounterChecks.Inc.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Inc.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.Monotonic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Value.NeverGoesBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Grow.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Size.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Throttle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CounterChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCounterSuite measures each callable of the suite of Counter on
// counterIdle.
func BenchmarkCounterSuite(b *testing.B) {
	b.Run("Inc", func(b *testing.B) {
		sequenced.BenchCounterSuite(b, sequenced.CounterChecks.Inc.Callable, counterIdle)
	})
	b.Run("Value", func(b *testing.B) {
		sequenced.BenchCounterSuite(b, sequenced.CounterChecks.Value.Callable, counterIdle)
	})
	b.Run("Grow", func(b *testing.B) {
		sequenced.BenchCounterSuite(b, sequenced.CounterChecks.Grow.Callable, counterIdle)
	})
	b.Run("Size", func(b *testing.B) {
		sequenced.BenchCounterSuite(b, sequenced.CounterChecks.Size.Callable, counterIdle)
	})
	b.Run("Throttle", func(b *testing.B) {
		sequenced.BenchCounterSuite(b, sequenced.CounterChecks.Throttle.Callable, counterIdle)
	})
	b.Run("Load", func(b *testing.B) {
		sequenced.BenchCounterSuite(b, sequenced.CounterChecks.Load.Callable, counterIdle)
	})
}

// notesSubjectEmbedded is sequenced.Notes under the name of the field that
// notesSubject embeds.
type notesSubjectEmbedded = sequenced.Notes

// notesSubject is a value of Notes that only TestNotesSuite builds.
type notesSubject struct{ notesSubjectEmbedded }

// TestNotesSuite checks the suite of Notes and the two forms of its
// application.
func TestNotesSuite(t *testing.T) {
	t.Parallel()

	t.Run("NotesSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.NotesSuite())
	})

	t.Run("RunNotesSuite", func(t *testing.T) {
		t.Parallel()
		var built *notesSubject
		sequenced.RunNotesSuite(
			t,
			func() *notesSubject {
				built = &notesSubject{}
				return built
			},
			conform.With(sequenced.NotesCheck{
				ID:        "Notes/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.NotesHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Notes(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.NotesHandle) sequenced.Notes { return &notesSubject{} },
			}),
			conform.Excuse(
				sequenced.NotesChecks.Push.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pop.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Text.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Text.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Text.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Head.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Head.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Head.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunNotesSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *notesSubject
		sequenced.RunNotesSuiteFromHandle(
			t,
			func(*sequenced.NotesHandle) *notesSubject {
				built = &notesSubject{}
				return built
			},
			conform.With(sequenced.NotesCheck{
				ID:        "Notes/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.NotesHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Notes(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.NotesHandle) sequenced.Notes { return &notesSubject{} },
			}),
			conform.Excuse(
				sequenced.NotesChecks.Push.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pop.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Text.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Text.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Text.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Pin.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Head.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Head.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.NotesChecks.Head.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// notesIdle returns the planted defect of a ceiling of Notes with a spend
// of nothing, a subject that answers each call at once.
func notesIdle(h *sequenced.NotesHandle) sequenced.Notes {
	return sequenced.NotesSuite().Calls[0].Over(h, func() {})
}

// TestNotesSuiteBudgets applies the allocation ceilings of the suite of
// Notes to notesIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestNotesSuiteBudgets(t *testing.T) {
	sequenced.RunNotesBudgets(
		t,
		notesIdle,
		conform.Ceilings(
			sequenced.NotesChecks.Push.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.NotesChecks.Pop.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.NotesChecks.Items.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.NotesChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.NotesChecks.Text.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.NotesChecks.Pin.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.NotesChecks.Head.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuitePush registers the fuzz target of Notes.Push on notesIdle,
// with every check of the suite excused.
func FuzzNotesSuitePush(f *testing.F) {
	sequenced.FuzzNotesSuite(
		f,
		sequenced.NotesChecks.Push.Callable,
		notesIdle,
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuitePop registers the fuzz target of Notes.Pop on notesIdle,
// with every check of the suite excused.
func FuzzNotesSuitePop(f *testing.F) {
	sequenced.FuzzNotesSuite(
		f,
		sequenced.NotesChecks.Pop.Callable,
		notesIdle,
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuiteItems registers the fuzz target of Notes.Items on
// notesIdle, with every check of the suite excused.
func FuzzNotesSuiteItems(f *testing.F) {
	sequenced.FuzzNotesSuite(
		f,
		sequenced.NotesChecks.Items.Callable,
		notesIdle,
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuiteAppend registers the fuzz target of Notes.Append on
// notesIdle, with every check of the suite excused.
func FuzzNotesSuiteAppend(f *testing.F) {
	sequenced.FuzzNotesSuite(
		f,
		sequenced.NotesChecks.Append.Callable,
		notesIdle,
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuiteText registers the fuzz target of Notes.Text on notesIdle,
// with every check of the suite excused.
func FuzzNotesSuiteText(f *testing.F) {
	sequenced.FuzzNotesSuite(
		f,
		sequenced.NotesChecks.Text.Callable,
		notesIdle,
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuitePin registers the fuzz target of Notes.Pin on notesIdle,
// with every check of the suite excused.
func FuzzNotesSuitePin(f *testing.F) {
	sequenced.FuzzNotesSuite(
		f,
		sequenced.NotesChecks.Pin.Callable,
		notesIdle,
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzNotesSuiteHead registers the fuzz target of Notes.Head on notesIdle,
// with every check of the suite excused.
func FuzzNotesSuiteHead(f *testing.F) {
	sequenced.FuzzNotesSuite(
		f,
		sequenced.NotesChecks.Head.Callable,
		notesIdle,
		conform.Excuse(
			sequenced.NotesChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Push.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pop.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Items.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Append.KeepsWithinItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Text.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Pin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.NotesChecks.Head.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkNotesSuite measures each callable of the suite of Notes on
// notesIdle.
func BenchmarkNotesSuite(b *testing.B) {
	b.Run("Push", func(b *testing.B) {
		sequenced.BenchNotesSuite(b, sequenced.NotesChecks.Push.Callable, notesIdle)
	})
	b.Run("Pop", func(b *testing.B) {
		sequenced.BenchNotesSuite(b, sequenced.NotesChecks.Pop.Callable, notesIdle)
	})
	b.Run("Items", func(b *testing.B) {
		sequenced.BenchNotesSuite(b, sequenced.NotesChecks.Items.Callable, notesIdle)
	})
	b.Run("Append", func(b *testing.B) {
		sequenced.BenchNotesSuite(b, sequenced.NotesChecks.Append.Callable, notesIdle)
	})
	b.Run("Text", func(b *testing.B) {
		sequenced.BenchNotesSuite(b, sequenced.NotesChecks.Text.Callable, notesIdle)
	})
	b.Run("Pin", func(b *testing.B) {
		sequenced.BenchNotesSuite(b, sequenced.NotesChecks.Pin.Callable, notesIdle)
	})
	b.Run("Head", func(b *testing.B) {
		sequenced.BenchNotesSuite(b, sequenced.NotesChecks.Head.Callable, notesIdle)
	})
}

// setSubjectEmbedded is sequenced.Set under the name of the field that
// setSubject embeds.
type setSubjectEmbedded = sequenced.Set

// setSubject is a value of Set that only TestSetSuite builds.
type setSubject struct{ setSubjectEmbedded }

// TestSetSuite checks the suite of Set and the two forms of its
// application.
func TestSetSuite(t *testing.T) {
	t.Parallel()

	t.Run("SetSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.SetSuite())
	})

	t.Run("RunSetSuite", func(t *testing.T) {
		t.Parallel()
		var built *setSubject
		sequenced.RunSetSuite(
			t,
			func() *setSubject {
				built = &setSubject{}
				return built
			},
			conform.With(sequenced.SetCheck{
				ID:        "Set/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.SetHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Set(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.SetHandle) sequenced.Set { return &setSubject{} },
			}),
			conform.Excuse(
				sequenced.SetChecks.Merge.MergesInAnyOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Join.MergesInAnyOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSetSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *setSubject
		sequenced.RunSetSuiteFromHandle(
			t,
			func(*sequenced.SetHandle) *setSubject {
				built = &setSubject{}
				return built
			},
			conform.With(sequenced.SetCheck{
				ID:        "Set/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.SetHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Set(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.SetHandle) sequenced.Set { return &setSubject{} },
			}),
			conform.Excuse(
				sequenced.SetChecks.Merge.MergesInAnyOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Join.MergesInAnyOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.SetChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// setIdle returns the planted defect of a ceiling of Set with a spend of
// nothing, a subject that answers each call at once.
func setIdle(h *sequenced.SetHandle) sequenced.Set {
	return sequenced.SetSuite().Calls[0].Over(h, func() {})
}

// TestSetSuiteBudgets applies the allocation ceilings of the suite of Set
// to setIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestSetSuiteBudgets(t *testing.T) {
	sequenced.RunSetBudgets(
		t,
		setIdle,
		conform.Ceilings(sequenced.SetChecks.Add.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			sequenced.SetChecks.Items.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.SetChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Join.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSetSuiteMerge registers the fuzz target of Set.Merge on setIdle,
// with every check of the suite excused.
func FuzzSetSuiteMerge(f *testing.F) {
	sequenced.FuzzSetSuite(
		f,
		sequenced.SetChecks.Merge.Callable,
		setIdle,
		conform.Excuse(
			sequenced.SetChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Join.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSetSuiteJoin registers the fuzz target of Set.Join on setIdle, with
// every check of the suite excused.
func FuzzSetSuiteJoin(f *testing.F) {
	sequenced.FuzzSetSuite(
		f,
		sequenced.SetChecks.Join.Callable,
		setIdle,
		conform.Excuse(
			sequenced.SetChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Join.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSetSuiteWeigh registers the fuzz target of Set.Weigh on setIdle,
// with every check of the suite excused.
func FuzzSetSuiteWeigh(f *testing.F) {
	sequenced.FuzzSetSuite(
		f,
		sequenced.SetChecks.Weigh.Callable,
		setIdle,
		conform.Excuse(
			sequenced.SetChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Join.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSetSuiteAdd registers the fuzz target of Set.Add on setIdle, with
// every check of the suite excused.
func FuzzSetSuiteAdd(f *testing.F) {
	sequenced.FuzzSetSuite(
		f,
		sequenced.SetChecks.Add.Callable,
		setIdle,
		conform.Excuse(
			sequenced.SetChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Join.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSetSuiteItems registers the fuzz target of Set.Items on setIdle,
// with every check of the suite excused.
func FuzzSetSuiteItems(f *testing.F) {
	sequenced.FuzzSetSuite(
		f,
		sequenced.SetChecks.Items.Callable,
		setIdle,
		conform.Excuse(
			sequenced.SetChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Join.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.SetChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSetSuite measures each callable of the suite of Set on setIdle.
func BenchmarkSetSuite(b *testing.B) {
	b.Run("Merge", func(b *testing.B) {
		sequenced.BenchSetSuite(b, sequenced.SetChecks.Merge.Callable, setIdle)
	})
	b.Run("Join", func(b *testing.B) {
		sequenced.BenchSetSuite(b, sequenced.SetChecks.Join.Callable, setIdle)
	})
	b.Run("Weigh", func(b *testing.B) {
		sequenced.BenchSetSuite(b, sequenced.SetChecks.Weigh.Callable, setIdle)
	})
	b.Run("Add", func(b *testing.B) {
		sequenced.BenchSetSuite(b, sequenced.SetChecks.Add.Callable, setIdle)
	})
	b.Run("Items", func(b *testing.B) {
		sequenced.BenchSetSuite(b, sequenced.SetChecks.Items.Callable, setIdle)
	})
}

// tallySubjectEmbedded is sequenced.TallyMethods under the name of the
// field that tallySubject embeds.
type tallySubjectEmbedded = sequenced.TallyMethods

// tallySubject is a value of Tally that only TestTallySuite builds.
type tallySubject struct{ tallySubjectEmbedded }

// TestTallySuite checks the suite of Tally and the two forms of its
// application.
func TestTallySuite(t *testing.T) {
	t.Parallel()

	t.Run("TallySuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.TallySuite())
	})

	t.Run("RunTallySuite", func(t *testing.T) {
		t.Parallel()
		var built *tallySubject
		sequenced.RunTallySuite(
			t,
			func() *tallySubject {
				built = &tallySubject{}
				return built
			},
			conform.With(sequenced.TallyCheck{
				ID:        "Tally/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.TallyHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.TallyMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.TallyHandle) sequenced.TallyMethods { return &tallySubject{} },
			}),
			conform.Excuse(
				sequenced.TallyChecks.Merge.MergesInAnyOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTallySuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *tallySubject
		sequenced.RunTallySuiteFromHandle(
			t,
			func(*sequenced.TallyHandle) *tallySubject {
				built = &tallySubject{}
				return built
			},
			conform.With(sequenced.TallyCheck{
				ID:        "Tally/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.TallyHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.TallyMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.TallyHandle) sequenced.TallyMethods { return &tallySubject{} },
			}),
			conform.Excuse(
				sequenced.TallyChecks.Merge.MergesInAnyOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.TallyChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// tallyIdle returns the planted defect of a ceiling of Tally with a spend
// of nothing, a subject that answers each call at once.
func tallyIdle(h *sequenced.TallyHandle) sequenced.TallyMethods {
	return sequenced.TallySuite().Calls[0].Over(h, func() {})
}

// TestTallySuiteBudgets applies the allocation ceilings of the suite of
// Tally to tallyIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestTallySuiteBudgets(t *testing.T) {
	sequenced.RunTallyBudgets(
		t,
		tallyIdle,
		conform.Ceilings(
			sequenced.TallyChecks.Merge.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.TallyChecks.Add.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.TallyChecks.Count.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.TallyChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTallySuiteMerge registers the fuzz target of Tally.Merge on
// tallyIdle, with every check of the suite excused.
func FuzzTallySuiteMerge(f *testing.F) {
	sequenced.FuzzTallySuite(
		f,
		sequenced.TallyChecks.Merge.Callable,
		tallyIdle,
		conform.Excuse(
			sequenced.TallyChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTallySuiteAdd registers the fuzz target of Tally.Add on tallyIdle,
// with every check of the suite excused.
func FuzzTallySuiteAdd(f *testing.F) {
	sequenced.FuzzTallySuite(
		f,
		sequenced.TallyChecks.Add.Callable,
		tallyIdle,
		conform.Excuse(
			sequenced.TallyChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTallySuiteCount registers the fuzz target of Tally.Count on
// tallyIdle, with every check of the suite excused.
func FuzzTallySuiteCount(f *testing.F) {
	sequenced.FuzzTallySuite(
		f,
		sequenced.TallyChecks.Count.Callable,
		tallyIdle,
		conform.Excuse(
			sequenced.TallyChecks.Merge.MergesInAnyOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.TallyChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkTallySuite measures each callable of the suite of Tally on
// tallyIdle.
func BenchmarkTallySuite(b *testing.B) {
	b.Run("Merge", func(b *testing.B) {
		sequenced.BenchTallySuite(b, sequenced.TallyChecks.Merge.Callable, tallyIdle)
	})
	b.Run("Add", func(b *testing.B) {
		sequenced.BenchTallySuite(b, sequenced.TallyChecks.Add.Callable, tallyIdle)
	})
	b.Run("Count", func(b *testing.B) {
		sequenced.BenchTallySuite(b, sequenced.TallyChecks.Count.Callable, tallyIdle)
	})
}

// pileSubjectEmbedded is sequenced.Pile under the name of the field that
// pileSubject embeds.
type pileSubjectEmbedded = sequenced.Pile

// pileSubject is a value of Pile that only TestPileSuite builds.
type pileSubject struct{ pileSubjectEmbedded }

// TestPileSuite checks the suite of Pile and the two forms of its
// application.
func TestPileSuite(t *testing.T) {
	t.Parallel()

	t.Run("PileSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.PileSuite())
	})

	t.Run("RunPileSuite", func(t *testing.T) {
		t.Parallel()
		var built *pileSubject
		sequenced.RunPileSuite(
			t,
			func() *pileSubject {
				built = &pileSubject{}
				return built
			},
			conform.With(sequenced.PileCheck{
				ID:        "Pile/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.PileHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Pile(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.PileHandle) sequenced.Pile { return &pileSubject{} },
			}),
			conform.Excuse(
				sequenced.PileChecks.Absorb.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Absorb.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Absorb.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Absorb.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPileSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *pileSubject
		sequenced.RunPileSuiteFromHandle(
			t,
			func(*sequenced.PileHandle) *pileSubject {
				built = &pileSubject{}
				return built
			},
			conform.With(sequenced.PileCheck{
				ID:        "Pile/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.PileHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Pile(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.PileHandle) sequenced.Pile { return &pileSubject{} },
			}),
			conform.Excuse(
				sequenced.PileChecks.Absorb.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Absorb.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Absorb.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Absorb.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.PileChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// pileIdle returns the planted defect of a ceiling of Pile with a spend of
// nothing, a subject that answers each call at once.
func pileIdle(h *sequenced.PileHandle) sequenced.Pile {
	return sequenced.PileSuite().Calls[0].Over(h, func() {})
}

// TestPileSuiteBudgets applies the allocation ceilings of the suite of
// Pile to pileIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestPileSuiteBudgets(t *testing.T) {
	sequenced.RunPileBudgets(
		t,
		pileIdle,
		conform.Ceilings(
			sequenced.PileChecks.Absorb.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(sequenced.PileChecks.Add.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			sequenced.PileChecks.Count.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPileSuiteAbsorb registers the fuzz target of Pile.Absorb on
// pileIdle, with every check of the suite excused.
func FuzzPileSuiteAbsorb(f *testing.F) {
	sequenced.FuzzPileSuite(
		f,
		sequenced.PileChecks.Absorb.Callable,
		pileIdle,
		conform.Excuse(
			sequenced.PileChecks.Absorb.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPileSuitePour registers the fuzz target of Pile.Pour on pileIdle,
// with every check of the suite excused.
func FuzzPileSuitePour(f *testing.F) {
	sequenced.FuzzPileSuite(
		f,
		sequenced.PileChecks.Pour.Callable,
		pileIdle,
		conform.Excuse(
			sequenced.PileChecks.Absorb.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPileSuiteAdd registers the fuzz target of Pile.Add on pileIdle, with
// every check of the suite excused.
func FuzzPileSuiteAdd(f *testing.F) {
	sequenced.FuzzPileSuite(
		f,
		sequenced.PileChecks.Add.Callable,
		pileIdle,
		conform.Excuse(
			sequenced.PileChecks.Absorb.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPileSuiteCount registers the fuzz target of Pile.Count on pileIdle,
// with every check of the suite excused.
func FuzzPileSuiteCount(f *testing.F) {
	sequenced.FuzzPileSuite(
		f,
		sequenced.PileChecks.Count.Callable,
		pileIdle,
		conform.Excuse(
			sequenced.PileChecks.Absorb.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Absorb.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.PileChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPileSuite measures each callable of the suite of Pile on
// pileIdle.
func BenchmarkPileSuite(b *testing.B) {
	b.Run("Absorb", func(b *testing.B) {
		sequenced.BenchPileSuite(b, sequenced.PileChecks.Absorb.Callable, pileIdle)
	})
	b.Run("Pour", func(b *testing.B) {
		sequenced.BenchPileSuite(b, sequenced.PileChecks.Pour.Callable, pileIdle)
	})
	b.Run("Add", func(b *testing.B) {
		sequenced.BenchPileSuite(b, sequenced.PileChecks.Add.Callable, pileIdle)
	})
	b.Run("Count", func(b *testing.B) {
		sequenced.BenchPileSuite(b, sequenced.PileChecks.Count.Callable, pileIdle)
	})
}

// runSubjectEmbedded is sequenced.Run under the name of the field that
// runSubject embeds.
type runSubjectEmbedded = sequenced.Run

// runSubject is a value of Run that only TestRunSuite builds.
type runSubject struct{ runSubjectEmbedded }

// TestRunSuite checks the suite of Run and the two forms of its
// application.
func TestRunSuite(t *testing.T) {
	t.Parallel()

	t.Run("RunSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.RunSuite())
	})

	t.Run("RunRunSuite", func(t *testing.T) {
		t.Parallel()
		var built *runSubject
		sequenced.RunRunSuite(
			t,
			func() *runSubject {
				built = &runSubject{}
				return built
			},
			conform.With(sequenced.RunCheck{
				ID:        "Run/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.RunHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Run(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.RunHandle) sequenced.Run { return &runSubject{} },
			}),
			conform.Excuse(
				sequenced.RunChecks.Advance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.FollowsItsGraph(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.FollowsItsGraph(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.FollowsItsGraph(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunRunSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *runSubject
		sequenced.RunRunSuiteFromHandle(
			t,
			func(*sequenced.RunHandle) *runSubject {
				built = &runSubject{}
				return built
			},
			conform.With(sequenced.RunCheck{
				ID:        "Run/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.RunHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Run(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.RunHandle) sequenced.Run { return &runSubject{} },
			}),
			conform.Excuse(
				sequenced.RunChecks.Advance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Advance.FollowsItsGraph(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Retry.FollowsItsGraph(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Restart.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Toggle.FollowsItsGraph(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RunChecks.Status.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// runIdle returns the planted defect of a ceiling of Run with a spend of
// nothing, a subject that answers each call at once.
func runIdle(h *sequenced.RunHandle) sequenced.Run {
	return sequenced.RunSuite().Calls[0].Over(h, func() {})
}

// TestRunSuiteBudgets applies the allocation ceilings of the suite of Run
// to runIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestRunSuiteBudgets(t *testing.T) {
	sequenced.RunRunBudgets(
		t,
		runIdle,
		conform.Ceilings(
			sequenced.RunChecks.Advance.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.RunChecks.Retry.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.RunChecks.Restart.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.RunChecks.Toggle.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.RunChecks.Status.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRunSuiteAdvance registers the fuzz target of Run.Advance on runIdle,
// with every check of the suite excused.
func FuzzRunSuiteAdvance(f *testing.F) {
	sequenced.FuzzRunSuite(
		f,
		sequenced.RunChecks.Advance.Callable,
		runIdle,
		conform.Excuse(
			sequenced.RunChecks.Advance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRunSuiteRetry registers the fuzz target of Run.Retry on runIdle,
// with every check of the suite excused.
func FuzzRunSuiteRetry(f *testing.F) {
	sequenced.FuzzRunSuite(
		f,
		sequenced.RunChecks.Retry.Callable,
		runIdle,
		conform.Excuse(
			sequenced.RunChecks.Advance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRunSuiteRestart registers the fuzz target of Run.Restart on runIdle,
// with every check of the suite excused.
func FuzzRunSuiteRestart(f *testing.F) {
	sequenced.FuzzRunSuite(
		f,
		sequenced.RunChecks.Restart.Callable,
		runIdle,
		conform.Excuse(
			sequenced.RunChecks.Advance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRunSuiteToggle registers the fuzz target of Run.Toggle on runIdle,
// with every check of the suite excused.
func FuzzRunSuiteToggle(f *testing.F) {
	sequenced.FuzzRunSuite(
		f,
		sequenced.RunChecks.Toggle.Callable,
		runIdle,
		conform.Excuse(
			sequenced.RunChecks.Advance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRunSuiteStatus registers the fuzz target of Run.Status on runIdle,
// with every check of the suite excused.
func FuzzRunSuiteStatus(f *testing.F) {
	sequenced.FuzzRunSuite(
		f,
		sequenced.RunChecks.Status.Callable,
		runIdle,
		conform.Excuse(
			sequenced.RunChecks.Advance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Advance.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Retry.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Restart.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Toggle.FollowsItsGraph(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RunChecks.Status.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkRunSuite measures each callable of the suite of Run on runIdle.
func BenchmarkRunSuite(b *testing.B) {
	b.Run("Advance", func(b *testing.B) {
		sequenced.BenchRunSuite(b, sequenced.RunChecks.Advance.Callable, runIdle)
	})
	b.Run("Retry", func(b *testing.B) {
		sequenced.BenchRunSuite(b, sequenced.RunChecks.Retry.Callable, runIdle)
	})
	b.Run("Restart", func(b *testing.B) {
		sequenced.BenchRunSuite(b, sequenced.RunChecks.Restart.Callable, runIdle)
	})
	b.Run("Toggle", func(b *testing.B) {
		sequenced.BenchRunSuite(b, sequenced.RunChecks.Toggle.Callable, runIdle)
	})
	b.Run("Status", func(b *testing.B) {
		sequenced.BenchRunSuite(b, sequenced.RunChecks.Status.Callable, runIdle)
	})
}

// logSubjectEmbedded is sequenced.Log under the name of the field that
// logSubject embeds.
type logSubjectEmbedded = sequenced.Log

// logSubject is a value of Log that only TestLogSuite builds.
type logSubject struct{ logSubjectEmbedded }

// TestLogSuite checks the suite of Log and the two forms of its
// application.
func TestLogSuite(t *testing.T) {
	t.Parallel()

	t.Run("LogSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.LogSuite())
	})

	t.Run("RunLogSuite", func(t *testing.T) {
		t.Parallel()
		var built *logSubject
		sequenced.RunLogSuite(
			t,
			func() *logSubject {
				built = &logSubject{}
				return built
			},
			conform.With(sequenced.LogCheck{
				ID:        "Log/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.LogHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Log(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.LogHandle) sequenced.Log { return &logSubject{} },
			}),
			conform.Excuse(
				sequenced.LogChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.Appender(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.AppenderGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.AppendsKeepAscending(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLogSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *logSubject
		sequenced.RunLogSuiteFromHandle(
			t,
			func(*sequenced.LogHandle) *logSubject {
				built = &logSubject{}
				return built
			},
			conform.With(sequenced.LogCheck{
				ID:        "Log/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.LogHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Log(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.LogHandle) sequenced.Log { return &logSubject{} },
			}),
			conform.Excuse(
				sequenced.LogChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.Appender(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.AppenderGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.LogChecks.Append.AppendsKeepAscending(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// logIdle returns the planted defect of a ceiling of Log with a spend of
// nothing, a subject that answers each call at once.
func logIdle(h *sequenced.LogHandle) sequenced.Log {
	return sequenced.LogSuite().Calls[0].Over(h, func() {})
}

// TestLogSuiteBudgets applies the allocation ceilings of the suite of Log
// to logIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestLogSuiteBudgets(t *testing.T) {
	sequenced.RunLogBudgets(
		t,
		logIdle,
		conform.Ceilings(
			sequenced.LogChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.Appender(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.AppenderGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.AppendsKeepAscending(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLogSuiteAppend registers the fuzz target of Log.Append on logIdle,
// with every check of the suite excused.
func FuzzLogSuiteAppend(f *testing.F) {
	sequenced.FuzzLogSuite(
		f,
		sequenced.LogChecks.Append.Callable,
		logIdle,
		conform.Excuse(
			sequenced.LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.Appender(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.AppenderGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.LogChecks.Append.AppendsKeepAscending(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLogSuite measures each callable of the suite of Log on logIdle.
func BenchmarkLogSuite(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		sequenced.BenchLogSuite(b, sequenced.LogChecks.Append.Callable, logIdle)
	})
}

// feedSubjectEmbedded is sequenced.Feed under the name of the field that
// feedSubject embeds.
type feedSubjectEmbedded = sequenced.Feed

// feedSubject is a value of Feed that only TestFeedSuite builds.
type feedSubject struct{ feedSubjectEmbedded }

// TestFeedSuite checks the suite of Feed and the two forms of its
// application.
func TestFeedSuite(t *testing.T) {
	t.Parallel()

	t.Run("FeedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.FeedSuite())
	})

	t.Run("RunFeedSuite", func(t *testing.T) {
		t.Parallel()
		var built *feedSubject
		sequenced.RunFeedSuite(
			t,
			func() *feedSubject {
				built = &feedSubject{}
				return built
			},
			conform.With(sequenced.FeedCheck{
				ID:        "Feed/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.FeedHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Feed(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.FeedHandle) sequenced.Feed { return &feedSubject{} },
			}),
			conform.Excuse(
				sequenced.FeedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.NoDuplicatesGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.YieldsEachOnce(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.NoDuplicatesGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.YieldsEachOnce(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.NoDuplicatesGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.YieldsEachOnce(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunFeedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *feedSubject
		sequenced.RunFeedSuiteFromHandle(
			t,
			func(*sequenced.FeedHandle) *feedSubject {
				built = &feedSubject{}
				return built
			},
			conform.With(sequenced.FeedCheck{
				ID:        "Feed/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.FeedHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Feed(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.FeedHandle) sequenced.Feed { return &feedSubject{} },
			}),
			conform.Excuse(
				sequenced.FeedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.NoDuplicatesGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.All.YieldsEachOnce(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.NoDuplicatesGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Each.YieldsEachOnce(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.NoDuplicatesGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.FeedChecks.Pairs.YieldsEachOnce(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// feedIdle returns the planted defect of a ceiling of Feed with a spend of
// nothing, a subject that answers each call at once.
func feedIdle(h *sequenced.FeedHandle) sequenced.Feed {
	return sequenced.FeedSuite().Calls[0].Over(h, func() {})
}

// TestFeedSuiteBudgets applies the allocation ceilings of the suite of
// Feed to feedIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestFeedSuiteBudgets(t *testing.T) {
	sequenced.RunFeedBudgets(
		t,
		feedIdle,
		conform.Ceilings(sequenced.FeedChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(sequenced.FeedChecks.All.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			sequenced.FeedChecks.Each.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.FeedChecks.Pairs.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzFeedSuitePut registers the fuzz target of Feed.Put on feedIdle, with
// every check of the suite excused.
func FuzzFeedSuitePut(f *testing.F) {
	sequenced.FuzzFeedSuite(
		f,
		sequenced.FeedChecks.Put.Callable,
		feedIdle,
		conform.Excuse(
			sequenced.FeedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzFeedSuiteAll registers the fuzz target of Feed.All on feedIdle, with
// every check of the suite excused.
func FuzzFeedSuiteAll(f *testing.F) {
	sequenced.FuzzFeedSuite(
		f,
		sequenced.FeedChecks.All.Callable,
		feedIdle,
		conform.Excuse(
			sequenced.FeedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzFeedSuiteEach registers the fuzz target of Feed.Each on feedIdle,
// with every check of the suite excused.
func FuzzFeedSuiteEach(f *testing.F) {
	sequenced.FuzzFeedSuite(
		f,
		sequenced.FeedChecks.Each.Callable,
		feedIdle,
		conform.Excuse(
			sequenced.FeedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzFeedSuitePairs registers the fuzz target of Feed.Pairs on feedIdle,
// with every check of the suite excused.
func FuzzFeedSuitePairs(f *testing.F) {
	sequenced.FuzzFeedSuite(
		f,
		sequenced.FeedChecks.Pairs.Callable,
		feedIdle,
		conform.Excuse(
			sequenced.FeedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.All.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Each.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.NoDuplicatesGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.FeedChecks.Pairs.YieldsEachOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkFeedSuite measures each callable of the suite of Feed on
// feedIdle.
func BenchmarkFeedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		sequenced.BenchFeedSuite(b, sequenced.FeedChecks.Put.Callable, feedIdle)
	})
	b.Run("All", func(b *testing.B) {
		sequenced.BenchFeedSuite(b, sequenced.FeedChecks.All.Callable, feedIdle)
	})
	b.Run("Each", func(b *testing.B) {
		sequenced.BenchFeedSuite(b, sequenced.FeedChecks.Each.Callable, feedIdle)
	})
	b.Run("Pairs", func(b *testing.B) {
		sequenced.BenchFeedSuite(b, sequenced.FeedChecks.Pairs.Callable, feedIdle)
	})
}

// registerSubjectEmbedded is sequenced.Register under the name of the
// field that registerSubject embeds.
type registerSubjectEmbedded = sequenced.Register

// registerSubject is a value of Register that only TestRegisterSuite
// builds.
type registerSubject struct{ registerSubjectEmbedded }

// TestRegisterSuite checks the suite of Register and the two forms of its
// application.
func TestRegisterSuite(t *testing.T) {
	t.Parallel()

	t.Run("RegisterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.RegisterSuite())
	})

	t.Run("RunRegisterSuite", func(t *testing.T) {
		t.Parallel()
		var built *registerSubject
		sequenced.RunRegisterSuite(
			t,
			func() *registerSubject {
				built = &registerSubject{}
				return built
			},
			conform.With(sequenced.RegisterCheck{
				ID:        "Register/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.RegisterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Register(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.RegisterHandle) sequenced.Register { return &registerSubject{} },
			}),
			conform.Excuse(
				sequenced.RegisterChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Set.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Set.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.AgreesWithTheReference(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunRegisterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *registerSubject
		sequenced.RunRegisterSuiteFromHandle(
			t,
			func(*sequenced.RegisterHandle) *registerSubject {
				built = &registerSubject{}
				return built
			},
			conform.With(sequenced.RegisterCheck{
				ID:        "Register/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.RegisterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Register(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.RegisterHandle) sequenced.Register { return &registerSubject{} },
			}),
			conform.Excuse(
				sequenced.RegisterChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Set.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Set.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.RegisterChecks.Get.AgreesWithTheReference(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// registerIdle returns the planted defect of a ceiling of Register with a
// spend of nothing, a subject that answers each call at once.
func registerIdle(h *sequenced.RegisterHandle) sequenced.Register {
	return sequenced.RegisterSuite().Calls[0].Over(h, func() {})
}

// TestRegisterSuiteBudgets applies the allocation ceilings of the suite of
// Register to registerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRegisterSuiteBudgets(t *testing.T) {
	sequenced.RunRegisterBudgets(
		t,
		registerIdle,
		conform.Ceilings(
			sequenced.RegisterChecks.Set.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.RegisterChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.AgreesWithTheReference(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRegisterSuiteSet registers the fuzz target of Register.Set on
// registerIdle, with every check of the suite excused.
func FuzzRegisterSuiteSet(f *testing.F) {
	sequenced.FuzzRegisterSuite(
		f,
		sequenced.RegisterChecks.Set.Callable,
		registerIdle,
		conform.Excuse(
			sequenced.RegisterChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.AgreesWithTheReference(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRegisterSuiteGet registers the fuzz target of Register.Get on
// registerIdle, with every check of the suite excused.
func FuzzRegisterSuiteGet(f *testing.F) {
	sequenced.FuzzRegisterSuite(
		f,
		sequenced.RegisterChecks.Get.Callable,
		registerIdle,
		conform.Excuse(
			sequenced.RegisterChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.RegisterChecks.Get.AgreesWithTheReference(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkRegisterSuite measures each callable of the suite of Register
// on registerIdle.
func BenchmarkRegisterSuite(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		sequenced.BenchRegisterSuite(b, sequenced.RegisterChecks.Set.Callable, registerIdle)
	})
	b.Run("Get", func(b *testing.B) {
		sequenced.BenchRegisterSuite(b, sequenced.RegisterChecks.Get.Callable, registerIdle)
	})
}

// TestDoubleSuite checks the suite of Double and the two forms of its
// application.
func TestDoubleSuite(t *testing.T) {
	t.Parallel()

	t.Run("DoubleSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.DoubleSuite())
	})

	t.Run("Double", func(t *testing.T) {
		t.Parallel()
		sequenced.RunDoubleSuite(t, sequenced.Double)
	})

	t.Run("RunDoubleSuite", func(t *testing.T) {
		t.Parallel()
		sequenced.RunDoubleSuite(
			t,
			sequenced.Double,
			conform.With(sequenced.DoubleCheck{
				ID:        "Double/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *sequenced.DoubleHandle) {
					assert.Equal(
						h,
						h.Subject(),
						sequenced.DoubleFunc(sequenced.Double),
						h.Statement(),
					)
				},
				Defect: func(*sequenced.DoubleHandle) sequenced.DoubleFunc { return nil },
			}),
			conform.Excuse(
				sequenced.DoubleChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.AgreesWithTheReference(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunDoubleSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		sequenced.RunDoubleSuiteFromHandle(
			t,
			func(*sequenced.DoubleHandle) sequenced.DoubleFunc {
				return sequenced.Double
			},
			conform.With(sequenced.DoubleCheck{
				ID:        "Double/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *sequenced.DoubleHandle) {
					assert.Equal(
						h,
						h.Subject(),
						sequenced.DoubleFunc(sequenced.Double),
						h.Statement(),
					)
				},
				Defect: func(*sequenced.DoubleHandle) sequenced.DoubleFunc { return nil },
			}),
			conform.Excuse(
				sequenced.DoubleChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.DoubleChecks.AgreesWithTheReference(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// BenchmarkDouble measures Double.
func BenchmarkDouble(b *testing.B) {
	sequenced.BenchDoubleSuite(b, sequenced.Double)
}

// FuzzDouble fuzzes the property checks of Double.
func FuzzDouble(f *testing.F) {
	sequenced.FuzzDoubleSuite(f, sequenced.Double)
}

// cachedSubjectEmbedded is sequenced.Cached under the name of the field
// that cachedSubject embeds.
type cachedSubjectEmbedded = sequenced.Cached

// cachedSubject is a value of Cached that only TestCachedSuite builds.
type cachedSubject struct{ cachedSubjectEmbedded }

// TestCachedSuite checks the suite of Cached and the two forms of its
// application.
func TestCachedSuite(t *testing.T) {
	t.Parallel()

	t.Run("CachedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.CachedSuite())
	})

	t.Run("RunCachedSuite", func(t *testing.T) {
		t.Parallel()
		var built *cachedSubject
		sequenced.RunCachedSuite(
			t,
			func() *cachedSubject {
				built = &cachedSubject{}
				return built
			},
			conform.With(sequenced.CachedCheck{
				ID:        "Cached/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.CachedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Cached(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.CachedHandle) sequenced.Cached { return &cachedSubject{} },
			}),
			conform.Excuse(
				sequenced.CachedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.CacheAgreesOverSequences(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.CacheAgreesOverSequences(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCachedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *cachedSubject
		sequenced.RunCachedSuiteFromHandle(
			t,
			func(*sequenced.CachedHandle) *cachedSubject {
				built = &cachedSubject{}
				return built
			},
			conform.With(sequenced.CachedCheck{
				ID:        "Cached/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.CachedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Cached(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.CachedHandle) sequenced.Cached { return &cachedSubject{} },
			}),
			conform.Excuse(
				sequenced.CachedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Get.CacheAgreesOverSequences(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Load.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Peek.CacheAgreesOverSequences(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.CachedChecks.Fetch.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// cachedIdle returns the planted defect of a ceiling of Cached with a
// spend of nothing, a subject that answers each call at once.
func cachedIdle(h *sequenced.CachedHandle) sequenced.Cached {
	return sequenced.CachedSuite().Calls[0].Over(h, func() {})
}

// TestCachedSuiteBudgets applies the allocation ceilings of the suite of
// Cached to cachedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCachedSuiteBudgets(t *testing.T) {
	sequenced.RunCachedBudgets(
		t,
		cachedIdle,
		conform.Ceilings(
			sequenced.CachedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CachedChecks.Load.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CachedChecks.Peek.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.CachedChecks.Fetch.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteGet registers the fuzz target of Cached.Get on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteGet(f *testing.F) {
	sequenced.FuzzCachedSuite(
		f,
		sequenced.CachedChecks.Get.Callable,
		cachedIdle,
		conform.Excuse(
			sequenced.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteLoad registers the fuzz target of Cached.Load on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteLoad(f *testing.F) {
	sequenced.FuzzCachedSuite(
		f,
		sequenced.CachedChecks.Load.Callable,
		cachedIdle,
		conform.Excuse(
			sequenced.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuitePeek registers the fuzz target of Cached.Peek on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuitePeek(f *testing.F) {
	sequenced.FuzzCachedSuite(
		f,
		sequenced.CachedChecks.Peek.Callable,
		cachedIdle,
		conform.Excuse(
			sequenced.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteFetch registers the fuzz target of Cached.Fetch on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteFetch(f *testing.F) {
	sequenced.FuzzCachedSuite(
		f,
		sequenced.CachedChecks.Fetch.Callable,
		cachedIdle,
		conform.Excuse(
			sequenced.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Get.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Load.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Peek.CacheAgreesOverSequences(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.CachedChecks.Fetch.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCachedSuite measures each callable of the suite of Cached on
// cachedIdle.
func BenchmarkCachedSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		sequenced.BenchCachedSuite(b, sequenced.CachedChecks.Get.Callable, cachedIdle)
	})
	b.Run("Load", func(b *testing.B) {
		sequenced.BenchCachedSuite(b, sequenced.CachedChecks.Load.Callable, cachedIdle)
	})
	b.Run("Peek", func(b *testing.B) {
		sequenced.BenchCachedSuite(b, sequenced.CachedChecks.Peek.Callable, cachedIdle)
	})
	b.Run("Fetch", func(b *testing.B) {
		sequenced.BenchCachedSuite(b, sequenced.CachedChecks.Fetch.Callable, cachedIdle)
	})
}

// indexSubjectEmbedded is sequenced.Index under the name of the field that
// indexSubject embeds.
type indexSubjectEmbedded = sequenced.Index

// indexSubject is a value of Index that only TestIndexSuite builds.
type indexSubject struct{ indexSubjectEmbedded }

// TestIndexSuite checks the suite of Index and the two forms of its
// application.
func TestIndexSuite(t *testing.T) {
	t.Parallel()

	t.Run("IndexSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.IndexSuite())
	})

	t.Run("RunIndexSuite", func(t *testing.T) {
		t.Parallel()
		var built *indexSubject
		sequenced.RunIndexSuite(
			t,
			func() *indexSubject {
				built = &indexSubject{}
				return built
			},
			conform.With(sequenced.IndexCheck{
				ID:        "Index/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.IndexHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Index(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.IndexHandle) sequenced.Index { return &indexSubject{} },
			}),
			conform.Excuse(
				sequenced.IndexChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunIndexSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *indexSubject
		sequenced.RunIndexSuiteFromHandle(
			t,
			func(*sequenced.IndexHandle) *indexSubject {
				built = &indexSubject{}
				return built
			},
			conform.With(sequenced.IndexCheck{
				ID:        "Index/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.IndexHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Index(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.IndexHandle) sequenced.Index { return &indexSubject{} },
			}),
			conform.Excuse(
				sequenced.IndexChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Get.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Find.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// indexIdle returns the planted defect of a ceiling of Index with a spend
// of nothing, a subject that answers each call at once.
func indexIdle(h *sequenced.IndexHandle) sequenced.Index {
	return sequenced.IndexSuite().Calls[0].Over(h, func() {})
}

// TestIndexSuiteBudgets applies the allocation ceilings of the suite of
// Index to indexIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestIndexSuiteBudgets(t *testing.T) {
	sequenced.RunIndexBudgets(
		t,
		indexIdle,
		conform.Ceilings(
			sequenced.IndexChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.IndexChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.IndexChecks.Find.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.IndexChecks.ByRev.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.IndexChecks.Lookup.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzIndexSuitePut registers the fuzz target of Index.Put on indexIdle,
// with every check of the suite excused.
func FuzzIndexSuitePut(f *testing.F) {
	sequenced.FuzzIndexSuite(
		f,
		sequenced.IndexChecks.Put.Callable,
		indexIdle,
		conform.Excuse(
			sequenced.IndexChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzIndexSuiteGet registers the fuzz target of Index.Get on indexIdle,
// with every check of the suite excused.
func FuzzIndexSuiteGet(f *testing.F) {
	sequenced.FuzzIndexSuite(
		f,
		sequenced.IndexChecks.Get.Callable,
		indexIdle,
		conform.Excuse(
			sequenced.IndexChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzIndexSuiteFind registers the fuzz target of Index.Find on indexIdle,
// with every check of the suite excused.
func FuzzIndexSuiteFind(f *testing.F) {
	sequenced.FuzzIndexSuite(
		f,
		sequenced.IndexChecks.Find.Callable,
		indexIdle,
		conform.Excuse(
			sequenced.IndexChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzIndexSuiteByRev registers the fuzz target of Index.ByRev on
// indexIdle, with every check of the suite excused.
func FuzzIndexSuiteByRev(f *testing.F) {
	sequenced.FuzzIndexSuite(
		f,
		sequenced.IndexChecks.ByRev.Callable,
		indexIdle,
		conform.Excuse(
			sequenced.IndexChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzIndexSuiteLookup registers the fuzz target of Index.Lookup on
// indexIdle, with every check of the suite excused.
func FuzzIndexSuiteLookup(f *testing.F) {
	sequenced.FuzzIndexSuite(
		f,
		sequenced.IndexChecks.Lookup.Callable,
		indexIdle,
		conform.Excuse(
			sequenced.IndexChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.ByRev.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadAfterWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.IndexChecks.Lookup.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkIndexSuite measures each callable of the suite of Index on
// indexIdle.
func BenchmarkIndexSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		sequenced.BenchIndexSuite(b, sequenced.IndexChecks.Put.Callable, indexIdle)
	})
	b.Run("Get", func(b *testing.B) {
		sequenced.BenchIndexSuite(b, sequenced.IndexChecks.Get.Callable, indexIdle)
	})
	b.Run("Find", func(b *testing.B) {
		sequenced.BenchIndexSuite(b, sequenced.IndexChecks.Find.Callable, indexIdle)
	})
	b.Run("ByRev", func(b *testing.B) {
		sequenced.BenchIndexSuite(b, sequenced.IndexChecks.ByRev.Callable, indexIdle)
	})
	b.Run("Lookup", func(b *testing.B) {
		sequenced.BenchIndexSuite(b, sequenced.IndexChecks.Lookup.Callable, indexIdle)
	})
}

// stampedSubjectEmbedded is sequenced.Stamped under the name of the field
// that stampedSubject embeds.
type stampedSubjectEmbedded = sequenced.Stamped

// stampedSubject is a value of Stamped that only TestStampedSuite builds.
type stampedSubject struct{ stampedSubjectEmbedded }

// TestStampedSuite checks the suite of Stamped and the two forms of its
// application.
func TestStampedSuite(t *testing.T) {
	t.Parallel()

	t.Run("StampedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.StampedSuite())
	})

	t.Run("RunStampedSuite", func(t *testing.T) {
		t.Parallel()
		var built *stampedSubject
		sequenced.RunStampedSuite(
			t,
			func() *stampedSubject {
				built = &stampedSubject{}
				return built
			},
			conform.With(sequenced.StampedCheck{
				ID:        "Stamped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.StampedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Stamped(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.StampedHandle) sequenced.Stamped { return &stampedSubject{} },
			}),
			conform.Excuse(
				sequenced.StampedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Store.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Store.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.MonotonicReads(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.MonotonicReadsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.ReadsNeverGoBack(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStampedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *stampedSubject
		sequenced.RunStampedSuiteFromHandle(
			t,
			func(*sequenced.StampedHandle) *stampedSubject {
				built = &stampedSubject{}
				return built
			},
			conform.With(sequenced.StampedCheck{
				ID:        "Stamped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.StampedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Stamped(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.StampedHandle) sequenced.Stamped { return &stampedSubject{} },
			}),
			conform.Excuse(
				sequenced.StampedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Store.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Store.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.MonotonicReads(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.MonotonicReadsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StampedChecks.Get.ReadsNeverGoBack(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// stampedIdle returns the planted defect of a ceiling of Stamped with a
// spend of nothing, a subject that answers each call at once.
func stampedIdle(h *sequenced.StampedHandle) sequenced.Stamped {
	return sequenced.StampedSuite().Calls[0].Over(h, func() {})
}

// TestStampedSuiteBudgets applies the allocation ceilings of the suite of
// Stamped to stampedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestStampedSuiteBudgets(t *testing.T) {
	sequenced.RunStampedBudgets(
		t,
		stampedIdle,
		conform.Ceilings(
			sequenced.StampedChecks.Store.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.StampedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.MonotonicReadsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ReadsNeverGoBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStampedSuiteStore registers the fuzz target of Stamped.Store on
// stampedIdle, with every check of the suite excused.
func FuzzStampedSuiteStore(f *testing.F) {
	sequenced.FuzzStampedSuite(
		f,
		sequenced.StampedChecks.Store.Callable,
		stampedIdle,
		conform.Excuse(
			sequenced.StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.MonotonicReadsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ReadsNeverGoBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStampedSuiteGet registers the fuzz target of Stamped.Get on
// stampedIdle, with every check of the suite excused.
func FuzzStampedSuiteGet(f *testing.F) {
	sequenced.FuzzStampedSuite(
		f,
		sequenced.StampedChecks.Get.Callable,
		stampedIdle,
		conform.Excuse(
			sequenced.StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Store.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.MonotonicReadsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StampedChecks.Get.ReadsNeverGoBack(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStampedSuite measures each callable of the suite of Stamped on
// stampedIdle.
func BenchmarkStampedSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		sequenced.BenchStampedSuite(b, sequenced.StampedChecks.Store.Callable, stampedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		sequenced.BenchStampedSuite(b, sequenced.StampedChecks.Get.Callable, stampedIdle)
	})
}

// versionedSubjectEmbedded is sequenced.Versioned under the name of the
// field that versionedSubject embeds.
type versionedSubjectEmbedded = sequenced.Versioned

// versionedSubject is a value of Versioned that only TestVersionedSuite
// builds.
type versionedSubject struct{ versionedSubjectEmbedded }

// TestVersionedSuite checks the suite of Versioned and the two forms of
// its application.
func TestVersionedSuite(t *testing.T) {
	t.Parallel()

	t.Run("VersionedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.VersionedSuite())
	})

	t.Run("RunVersionedSuite", func(t *testing.T) {
		t.Parallel()
		var built *versionedSubject
		sequenced.RunVersionedSuite(
			t,
			func() *versionedSubject {
				built = &versionedSubject{}
				return built
			},
			conform.With(sequenced.VersionedCheck{
				ID:        "Versioned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.VersionedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Versioned(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.VersionedHandle) sequenced.Versioned { return &versionedSubject{} },
			}),
			conform.Excuse(
				sequenced.VersionedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunVersionedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *versionedSubject
		sequenced.RunVersionedSuiteFromHandle(
			t,
			func(*sequenced.VersionedHandle) *versionedSubject {
				built = &versionedSubject{}
				return built
			},
			conform.With(sequenced.VersionedCheck{
				ID:        "Versioned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.VersionedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Versioned(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.VersionedHandle) sequenced.Versioned { return &versionedSubject{} },
			}),
			conform.Excuse(
				sequenced.VersionedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.VersionedChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// versionedIdle returns the planted defect of a ceiling of Versioned with
// a spend of nothing, a subject that answers each call at once.
func versionedIdle(h *sequenced.VersionedHandle) sequenced.Versioned {
	return sequenced.VersionedSuite().Calls[0].Over(h, func() {})
}

// TestVersionedSuiteBudgets applies the allocation ceilings of the suite
// of Versioned to versionedIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestVersionedSuiteBudgets(t *testing.T) {
	sequenced.RunVersionedBudgets(
		t,
		versionedIdle,
		conform.Ceilings(
			sequenced.VersionedChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.VersionedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVersionedSuitePut registers the fuzz target of Versioned.Put on
// versionedIdle, with every check of the suite excused.
func FuzzVersionedSuitePut(f *testing.F) {
	sequenced.FuzzVersionedSuite(
		f,
		sequenced.VersionedChecks.Put.Callable,
		versionedIdle,
		conform.Excuse(
			sequenced.VersionedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVersionedSuiteGet registers the fuzz target of Versioned.Get on
// versionedIdle, with every check of the suite excused.
func FuzzVersionedSuiteGet(f *testing.F) {
	sequenced.FuzzVersionedSuite(
		f,
		sequenced.VersionedChecks.Get.Callable,
		versionedIdle,
		conform.Excuse(
			sequenced.VersionedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.VersionedChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkVersionedSuite measures each callable of the suite of Versioned
// on versionedIdle.
func BenchmarkVersionedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		sequenced.BenchVersionedSuite(b, sequenced.VersionedChecks.Put.Callable, versionedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		sequenced.BenchVersionedSuite(b, sequenced.VersionedChecks.Get.Callable, versionedIdle)
	})
}

// accountsSubjectEmbedded is sequenced.Accounts under the name of the
// field that accountsSubject embeds.
type accountsSubjectEmbedded = sequenced.Accounts

// accountsSubject is a value of Accounts that only TestAccountsSuite
// builds.
type accountsSubject struct{ accountsSubjectEmbedded }

// TestAccountsSuite checks the suite of Accounts and the two forms of its
// application.
func TestAccountsSuite(t *testing.T) {
	t.Parallel()

	t.Run("AccountsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.AccountsSuite())
	})

	t.Run("RunAccountsSuite", func(t *testing.T) {
		t.Parallel()
		var built *accountsSubject
		sequenced.RunAccountsSuite(
			t,
			func() *accountsSubject {
				built = &accountsSubject{}
				return built
			},
			conform.With(sequenced.AccountsCheck{
				ID:        "Accounts/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.AccountsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Accounts(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.AccountsHandle) sequenced.Accounts { return &accountsSubject{} },
			}),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Total.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Total.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Total.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Names.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Names.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Names.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Audit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Audit.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunAccountsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *accountsSubject
		sequenced.RunAccountsSuiteFromHandle(
			t,
			func(*sequenced.AccountsHandle) *accountsSubject {
				built = &accountsSubject{}
				return built
			},
			conform.With(sequenced.AccountsCheck{
				ID:        "Accounts/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.AccountsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Accounts(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.AccountsHandle) sequenced.Accounts { return &accountsSubject{} },
			}),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Total.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Total.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Total.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Rename.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Names.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Names.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Names.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Fold.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Audit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.AccountsChecks.Audit.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// accountsIdle returns the planted defect of a ceiling of Accounts with a
// spend of nothing, a subject that answers each call at once.
func accountsIdle(h *sequenced.AccountsHandle) sequenced.Accounts {
	return sequenced.AccountsSuite().Calls[0].Over(h, func() {})
}

// TestAccountsSuiteBudgets applies the allocation ceilings of the suite of
// Accounts to accountsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestAccountsSuiteBudgets(t *testing.T) {
	sequenced.RunAccountsBudgets(
		t,
		accountsIdle,
		conform.Ceilings(
			sequenced.AccountsChecks.Transfer.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.AccountsChecks.Total.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.AccountsChecks.Rename.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.AccountsChecks.Names.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.AccountsChecks.Fold.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.AccountsChecks.Audit.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountsSuiteTransfer registers the fuzz target of Accounts.Transfer
// on accountsIdle, with every check of the suite excused.
func FuzzAccountsSuiteTransfer(f *testing.F) {
	sequenced.FuzzAccountsSuite(
		f,
		sequenced.AccountsChecks.Transfer.Callable,
		accountsIdle,
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountsSuiteTotal registers the fuzz target of Accounts.Total on
// accountsIdle, with every check of the suite excused.
func FuzzAccountsSuiteTotal(f *testing.F) {
	sequenced.FuzzAccountsSuite(
		f,
		sequenced.AccountsChecks.Total.Callable,
		accountsIdle,
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountsSuiteRename registers the fuzz target of Accounts.Rename on
// accountsIdle, with every check of the suite excused.
func FuzzAccountsSuiteRename(f *testing.F) {
	sequenced.FuzzAccountsSuite(
		f,
		sequenced.AccountsChecks.Rename.Callable,
		accountsIdle,
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountsSuiteNames registers the fuzz target of Accounts.Names on
// accountsIdle, with every check of the suite excused.
func FuzzAccountsSuiteNames(f *testing.F) {
	sequenced.FuzzAccountsSuite(
		f,
		sequenced.AccountsChecks.Names.Callable,
		accountsIdle,
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountsSuiteFold registers the fuzz target of Accounts.Fold on
// accountsIdle, with every check of the suite excused.
func FuzzAccountsSuiteFold(f *testing.F) {
	sequenced.FuzzAccountsSuite(
		f,
		sequenced.AccountsChecks.Fold.Callable,
		accountsIdle,
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountsSuiteAudit registers the fuzz target of Accounts.Audit on
// accountsIdle, with every check of the suite excused.
func FuzzAccountsSuiteAudit(f *testing.F) {
	sequenced.FuzzAccountsSuite(
		f,
		sequenced.AccountsChecks.Audit.Callable,
		accountsIdle,
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Transfer.ConservesTheQuantity(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Total.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Rename.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Names.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Fold.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.AccountsChecks.Audit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkAccountsSuite measures each callable of the suite of Accounts
// on accountsIdle.
func BenchmarkAccountsSuite(b *testing.B) {
	b.Run("Transfer", func(b *testing.B) {
		sequenced.BenchAccountsSuite(b, sequenced.AccountsChecks.Transfer.Callable, accountsIdle)
	})
	b.Run("Total", func(b *testing.B) {
		sequenced.BenchAccountsSuite(b, sequenced.AccountsChecks.Total.Callable, accountsIdle)
	})
	b.Run("Rename", func(b *testing.B) {
		sequenced.BenchAccountsSuite(b, sequenced.AccountsChecks.Rename.Callable, accountsIdle)
	})
	b.Run("Names", func(b *testing.B) {
		sequenced.BenchAccountsSuite(b, sequenced.AccountsChecks.Names.Callable, accountsIdle)
	})
	b.Run("Fold", func(b *testing.B) {
		sequenced.BenchAccountsSuite(b, sequenced.AccountsChecks.Fold.Callable, accountsIdle)
	})
	b.Run("Audit", func(b *testing.B) {
		sequenced.BenchAccountsSuite(b, sequenced.AccountsChecks.Audit.Callable, accountsIdle)
	})
}

// journalSubjectEmbedded is sequenced.Journal under the name of the field
// that journalSubject embeds.
type journalSubjectEmbedded = sequenced.Journal

// journalSubject is a value of Journal that only TestJournalSuite builds.
type journalSubject struct{ journalSubjectEmbedded }

// TestJournalSuite checks the suite of Journal and the two forms of its
// application.
func TestJournalSuite(t *testing.T) {
	t.Parallel()

	t.Run("JournalSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.JournalSuite())
	})

	t.Run("RunJournalSuite", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		sequenced.RunJournalSuite(
			t,
			func() *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(sequenced.JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.JournalHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Journal(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.JournalHandle) sequenced.Journal { return &journalSubject{} },
			}),
			conform.Excuse(
				sequenced.JournalChecks.Record.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.NoPartialWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.NoPartialWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.NoTraceOfAFailedCall(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunJournalSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		sequenced.RunJournalSuiteFromHandle(
			t,
			func(*sequenced.JournalHandle) *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(sequenced.JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.JournalHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Journal(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.JournalHandle) sequenced.Journal { return &journalSubject{} },
			}),
			conform.Excuse(
				sequenced.JournalChecks.Record.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.NoPartialWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.NoPartialWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Record.NoTraceOfAFailedCall(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.JournalChecks.Entry.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// journalIdle returns the planted defect of a ceiling of Journal with a
// spend of nothing, a subject that answers each call at once.
func journalIdle(h *sequenced.JournalHandle) sequenced.Journal {
	return sequenced.JournalSuite().Calls[0].Over(h, func() {})
}

// TestJournalSuiteBudgets applies the allocation ceilings of the suite of
// Journal to journalIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestJournalSuiteBudgets(t *testing.T) {
	sequenced.RunJournalBudgets(
		t,
		journalIdle,
		conform.Ceilings(
			sequenced.JournalChecks.Record.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.JournalChecks.Entry.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoPartialWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoPartialWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoTraceOfAFailedCall(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzJournalSuiteRecord registers the fuzz target of Journal.Record on
// journalIdle, with every check of the suite excused.
func FuzzJournalSuiteRecord(f *testing.F) {
	sequenced.FuzzJournalSuite(
		f,
		sequenced.JournalChecks.Record.Callable,
		journalIdle,
		conform.Excuse(
			sequenced.JournalChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoPartialWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoPartialWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoTraceOfAFailedCall(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzJournalSuiteEntry registers the fuzz target of Journal.Entry on
// journalIdle, with every check of the suite excused.
func FuzzJournalSuiteEntry(f *testing.F) {
	sequenced.FuzzJournalSuite(
		f,
		sequenced.JournalChecks.Entry.Callable,
		journalIdle,
		conform.Excuse(
			sequenced.JournalChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoPartialWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoPartialWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Record.NoTraceOfAFailedCall(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.JournalChecks.Entry.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkJournalSuite measures each callable of the suite of Journal on
// journalIdle.
func BenchmarkJournalSuite(b *testing.B) {
	b.Run("Record", func(b *testing.B) {
		sequenced.BenchJournalSuite(b, sequenced.JournalChecks.Record.Callable, journalIdle)
	})
	b.Run("Entry", func(b *testing.B) {
		sequenced.BenchJournalSuite(b, sequenced.JournalChecks.Entry.Callable, journalIdle)
	})
}

// stockSubjectEmbedded is sequenced.Stock under the name of the field that
// stockSubject embeds.
type stockSubjectEmbedded = sequenced.Stock

// stockSubject is a value of Stock that only TestStockSuite builds.
type stockSubject struct{ stockSubjectEmbedded }

// TestStockSuite checks the suite of Stock and the two forms of its
// application.
func TestStockSuite(t *testing.T) {
	t.Parallel()

	t.Run("StockSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.StockSuite())
	})

	t.Run("RunStockSuite", func(t *testing.T) {
		t.Parallel()
		var built *stockSubject
		sequenced.RunStockSuite(
			t,
			func() *stockSubject {
				built = &stockSubject{}
				return built
			},
			conform.With(sequenced.StockCheck{
				ID:        "Stock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.StockHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Stock(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.StockHandle) sequenced.Stock { return &stockSubject{} },
			}),
			conform.Excuse(
				sequenced.StockChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Items.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Items.IterationSeesMutations(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.IterationSeesMutations(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.IterationSeesMutations(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Sizes.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Snapshot.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Snapshot.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Snapshot.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStockSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *stockSubject
		sequenced.RunStockSuiteFromHandle(
			t,
			func(*sequenced.StockHandle) *stockSubject {
				built = &stockSubject{}
				return built
			},
			conform.With(sequenced.StockCheck{
				ID:        "Stock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.StockHandle) {
					got := h.Subject()
					assert.Equal(h, got, sequenced.Stock(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sequenced.StockHandle) sequenced.Stock { return &stockSubject{} },
			}),
			conform.Excuse(
				sequenced.StockChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Items.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Items.IterationSeesMutations(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Each.IterationSeesMutations(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Counts.IterationSeesMutations(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Sizes.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Snapshot.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Snapshot.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Snapshot.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.StockChecks.Remove.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// stockIdle returns the planted defect of a ceiling of Stock with a spend
// of nothing, a subject that answers each call at once.
func stockIdle(h *sequenced.StockHandle) sequenced.Stock {
	return sequenced.StockSuite().Calls[0].Over(h, func() {})
}

// TestStockSuiteBudgets applies the allocation ceilings of the suite of
// Stock to stockIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestStockSuiteBudgets(t *testing.T) {
	sequenced.RunStockBudgets(
		t,
		stockIdle,
		conform.Ceilings(
			sequenced.StockChecks.Items.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.StockChecks.Each.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.StockChecks.Counts.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.StockChecks.Sizes.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.StockChecks.Snapshot.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.StockChecks.Add.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.StockChecks.Remove.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStockSuiteItems registers the fuzz target of Stock.Items on
// stockIdle, with every check of the suite excused.
func FuzzStockSuiteItems(f *testing.F) {
	sequenced.FuzzStockSuite(
		f,
		sequenced.StockChecks.Items.Callable,
		stockIdle,
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStockSuiteEach registers the fuzz target of Stock.Each on stockIdle,
// with every check of the suite excused.
func FuzzStockSuiteEach(f *testing.F) {
	sequenced.FuzzStockSuite(
		f,
		sequenced.StockChecks.Each.Callable,
		stockIdle,
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStockSuiteCounts registers the fuzz target of Stock.Counts on
// stockIdle, with every check of the suite excused.
func FuzzStockSuiteCounts(f *testing.F) {
	sequenced.FuzzStockSuite(
		f,
		sequenced.StockChecks.Counts.Callable,
		stockIdle,
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStockSuiteSizes registers the fuzz target of Stock.Sizes on
// stockIdle, with every check of the suite excused.
func FuzzStockSuiteSizes(f *testing.F) {
	sequenced.FuzzStockSuite(
		f,
		sequenced.StockChecks.Sizes.Callable,
		stockIdle,
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStockSuiteSnapshot registers the fuzz target of Stock.Snapshot on
// stockIdle, with every check of the suite excused.
func FuzzStockSuiteSnapshot(f *testing.F) {
	sequenced.FuzzStockSuite(
		f,
		sequenced.StockChecks.Snapshot.Callable,
		stockIdle,
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStockSuiteAdd registers the fuzz target of Stock.Add on stockIdle,
// with every check of the suite excused.
func FuzzStockSuiteAdd(f *testing.F) {
	sequenced.FuzzStockSuite(
		f,
		sequenced.StockChecks.Add.Callable,
		stockIdle,
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStockSuiteRemove registers the fuzz target of Stock.Remove on
// stockIdle, with every check of the suite excused.
func FuzzStockSuiteRemove(f *testing.F) {
	sequenced.FuzzStockSuite(
		f,
		sequenced.StockChecks.Remove.Callable,
		stockIdle,
		conform.Excuse(
			sequenced.StockChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Items.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Each.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Counts.IterationSeesMutations(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.StockChecks.Remove.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStockSuite measures each callable of the suite of Stock on
// stockIdle.
func BenchmarkStockSuite(b *testing.B) {
	b.Run("Items", func(b *testing.B) {
		sequenced.BenchStockSuite(b, sequenced.StockChecks.Items.Callable, stockIdle)
	})
	b.Run("Each", func(b *testing.B) {
		sequenced.BenchStockSuite(b, sequenced.StockChecks.Each.Callable, stockIdle)
	})
	b.Run("Counts", func(b *testing.B) {
		sequenced.BenchStockSuite(b, sequenced.StockChecks.Counts.Callable, stockIdle)
	})
	b.Run("Sizes", func(b *testing.B) {
		sequenced.BenchStockSuite(b, sequenced.StockChecks.Sizes.Callable, stockIdle)
	})
	b.Run("Snapshot", func(b *testing.B) {
		sequenced.BenchStockSuite(b, sequenced.StockChecks.Snapshot.Callable, stockIdle)
	})
	b.Run("Add", func(b *testing.B) {
		sequenced.BenchStockSuite(b, sequenced.StockChecks.Add.Callable, stockIdle)
	})
	b.Run("Remove", func(b *testing.B) {
		sequenced.BenchStockSuite(b, sequenced.StockChecks.Remove.Callable, stockIdle)
	})
}

// outboxSubjectEmbedded is sequenced.Outbox under the name of the field
// that outboxSubject embeds.
type outboxSubjectEmbedded = sequenced.Outbox

// outboxSubject is a value of Outbox that only TestOutboxSuite builds.
type outboxSubject struct{ outboxSubjectEmbedded }

// TestOutboxSuite checks the suite of Outbox and the two forms of its
// application.
func TestOutboxSuite(t *testing.T) {
	t.Parallel()

	t.Run("OutboxSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.OutboxSuite())
	})

	t.Run("RunOutboxSuite", func(t *testing.T) {
		t.Parallel()
		var built *outboxSubject
		sequenced.RunOutboxSuite(
			t,
			func() *outboxSubject {
				built = &outboxSubject{}
				return built
			},
			conform.With(sequenced.OutboxCheck{
				ID:        "Outbox/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.OutboxHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Outbox(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.OutboxHandle) sequenced.Outbox { return &outboxSubject{} },
			}),
			conform.Excuse(
				sequenced.OutboxChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.DeliversEachAppend(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Subscribe.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Subscribe.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Notes.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Notes.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Logged.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Logged.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Logged.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Queue.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Queue.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Flush.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Flush.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Drain.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunOutboxSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *outboxSubject
		sequenced.RunOutboxSuiteFromHandle(
			t,
			func(*sequenced.OutboxHandle) *outboxSubject {
				built = &outboxSubject{}
				return built
			},
			conform.With(sequenced.OutboxCheck{
				ID:        "Outbox/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.OutboxHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Outbox(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.OutboxHandle) sequenced.Outbox { return &outboxSubject{} },
			}),
			conform.Excuse(
				sequenced.OutboxChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Append.DeliversEachAppend(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Subscribe.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Subscribe.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Note.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Notes.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Notes.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Log.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Logged.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Logged.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Logged.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Post.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Queue.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Queue.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Push.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Flush.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Flush.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Pour.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.OutboxChecks.Drain.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// outboxIdle returns the planted defect of a ceiling of Outbox with a
// spend of nothing, a subject that answers each call at once.
func outboxIdle(h *sequenced.OutboxHandle) sequenced.Outbox {
	return sequenced.OutboxSuite().Calls[0].Over(h, func() {})
}

// TestOutboxSuiteBudgets applies the allocation ceilings of the suite of
// Outbox to outboxIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestOutboxSuiteBudgets(t *testing.T) {
	sequenced.RunOutboxBudgets(
		t,
		outboxIdle,
		conform.Ceilings(
			sequenced.OutboxChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Subscribe.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Note.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Notes.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Log.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Logged.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Post.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Queue.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Push.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Flush.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Pour.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.OutboxChecks.Drain.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteAppend registers the fuzz target of Outbox.Append on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteAppend(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Append.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteSubscribe registers the fuzz target of Outbox.Subscribe
// on outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteSubscribe(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Subscribe.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteNote registers the fuzz target of Outbox.Note on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteNote(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Note.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteNotes registers the fuzz target of Outbox.Notes on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteNotes(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Notes.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteLog registers the fuzz target of Outbox.Log on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteLog(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Log.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteLogged registers the fuzz target of Outbox.Logged on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteLogged(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Logged.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuitePost registers the fuzz target of Outbox.Post on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuitePost(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Post.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteQueue registers the fuzz target of Outbox.Queue on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteQueue(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Queue.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuitePush registers the fuzz target of Outbox.Push on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuitePush(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Push.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteFlush registers the fuzz target of Outbox.Flush on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteFlush(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Flush.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuitePour registers the fuzz target of Outbox.Pour on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuitePour(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Pour.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOutboxSuiteDrain registers the fuzz target of Outbox.Drain on
// outboxIdle, with every check of the suite excused.
func FuzzOutboxSuiteDrain(f *testing.F) {
	sequenced.FuzzOutboxSuite(
		f,
		sequenced.OutboxChecks.Drain.Callable,
		outboxIdle,
		conform.Excuse(
			sequenced.OutboxChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Append.DeliversEachAppend(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Subscribe.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Note.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Notes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Log.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Logged.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Queue.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Push.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Flush.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Pour.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.OutboxChecks.Drain.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkOutboxSuite measures each callable of the suite of Outbox on
// outboxIdle.
func BenchmarkOutboxSuite(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Append.Callable, outboxIdle)
	})
	b.Run("Subscribe", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Subscribe.Callable, outboxIdle)
	})
	b.Run("Note", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Note.Callable, outboxIdle)
	})
	b.Run("Notes", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Notes.Callable, outboxIdle)
	})
	b.Run("Log", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Log.Callable, outboxIdle)
	})
	b.Run("Logged", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Logged.Callable, outboxIdle)
	})
	b.Run("Post", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Post.Callable, outboxIdle)
	})
	b.Run("Queue", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Queue.Callable, outboxIdle)
	})
	b.Run("Push", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Push.Callable, outboxIdle)
	})
	b.Run("Flush", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Flush.Callable, outboxIdle)
	})
	b.Run("Pour", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Pour.Callable, outboxIdle)
	})
	b.Run("Drain", func(b *testing.B) {
		sequenced.BenchOutboxSuite(b, sequenced.OutboxChecks.Drain.Callable, outboxIdle)
	})
}

// brokerSubjectEmbedded is sequenced.Broker under the name of the field
// that brokerSubject embeds.
type brokerSubjectEmbedded = sequenced.Broker

// brokerSubject is a value of Broker that only TestBrokerSuite builds.
type brokerSubject struct{ brokerSubjectEmbedded }

// TestBrokerSuite checks the suite of Broker and the two forms of its
// application.
func TestBrokerSuite(t *testing.T) {
	t.Parallel()

	t.Run("BrokerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.BrokerSuite())
	})

	t.Run("RunBrokerSuite", func(t *testing.T) {
		t.Parallel()
		var built *brokerSubject
		sequenced.RunBrokerSuite(
			t,
			func() *brokerSubject {
				built = &brokerSubject{}
				return built
			},
			conform.With(sequenced.BrokerCheck{
				ID:        "Broker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.BrokerHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Broker(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.BrokerHandle) sequenced.Broker { return &brokerSubject{} },
			}),
			conform.Excuse(
				sequenced.BrokerChecks.Send.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.DeliversByItsMode(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Listen.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Listen.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.DeliversByItsMode(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Tap.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBrokerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *brokerSubject
		sequenced.RunBrokerSuiteFromHandle(
			t,
			func(*sequenced.BrokerHandle) *brokerSubject {
				built = &brokerSubject{}
				return built
			},
			conform.With(sequenced.BrokerCheck{
				ID:        "Broker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.BrokerHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Broker(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.BrokerHandle) sequenced.Broker { return &brokerSubject{} },
			}),
			conform.Excuse(
				sequenced.BrokerChecks.Send.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Send.DeliversByItsMode(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Listen.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Listen.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Fire.DeliversByItsMode(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Refire.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.BrokerChecks.Tap.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// brokerIdle returns the planted defect of a ceiling of Broker with a
// spend of nothing, a subject that answers each call at once.
func brokerIdle(h *sequenced.BrokerHandle) sequenced.Broker {
	return sequenced.BrokerSuite().Calls[0].Over(h, func() {})
}

// TestBrokerSuiteBudgets applies the allocation ceilings of the suite of
// Broker to brokerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestBrokerSuiteBudgets(t *testing.T) {
	sequenced.RunBrokerBudgets(
		t,
		brokerIdle,
		conform.Ceilings(
			sequenced.BrokerChecks.Send.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.BrokerChecks.Resend.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.BrokerChecks.Listen.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.BrokerChecks.Fire.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.BrokerChecks.Refire.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.BrokerChecks.Tap.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Tap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBrokerSuiteSend registers the fuzz target of Broker.Send on
// brokerIdle, with every check of the suite excused.
func FuzzBrokerSuiteSend(f *testing.F) {
	sequenced.FuzzBrokerSuite(
		f,
		sequenced.BrokerChecks.Send.Callable,
		brokerIdle,
		conform.Excuse(
			sequenced.BrokerChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Tap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBrokerSuiteResend registers the fuzz target of Broker.Resend on
// brokerIdle, with every check of the suite excused.
func FuzzBrokerSuiteResend(f *testing.F) {
	sequenced.FuzzBrokerSuite(
		f,
		sequenced.BrokerChecks.Resend.Callable,
		brokerIdle,
		conform.Excuse(
			sequenced.BrokerChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Tap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBrokerSuiteListen registers the fuzz target of Broker.Listen on
// brokerIdle, with every check of the suite excused.
func FuzzBrokerSuiteListen(f *testing.F) {
	sequenced.FuzzBrokerSuite(
		f,
		sequenced.BrokerChecks.Listen.Callable,
		brokerIdle,
		conform.Excuse(
			sequenced.BrokerChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Tap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBrokerSuiteFire registers the fuzz target of Broker.Fire on
// brokerIdle, with every check of the suite excused.
func FuzzBrokerSuiteFire(f *testing.F) {
	sequenced.FuzzBrokerSuite(
		f,
		sequenced.BrokerChecks.Fire.Callable,
		brokerIdle,
		conform.Excuse(
			sequenced.BrokerChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Tap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBrokerSuiteRefire registers the fuzz target of Broker.Refire on
// brokerIdle, with every check of the suite excused.
func FuzzBrokerSuiteRefire(f *testing.F) {
	sequenced.FuzzBrokerSuite(
		f,
		sequenced.BrokerChecks.Refire.Callable,
		brokerIdle,
		conform.Excuse(
			sequenced.BrokerChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Tap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBrokerSuiteTap registers the fuzz target of Broker.Tap on
// brokerIdle, with every check of the suite excused.
func FuzzBrokerSuiteTap(f *testing.F) {
	sequenced.FuzzBrokerSuite(
		f,
		sequenced.BrokerChecks.Tap.Callable,
		brokerIdle,
		conform.Excuse(
			sequenced.BrokerChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Send.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Resend.RedeliversOnDemand(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Listen.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Fire.DeliversByItsMode(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Refire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.BrokerChecks.Tap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBrokerSuite measures each callable of the suite of Broker on
// brokerIdle.
func BenchmarkBrokerSuite(b *testing.B) {
	b.Run("Send", func(b *testing.B) {
		sequenced.BenchBrokerSuite(b, sequenced.BrokerChecks.Send.Callable, brokerIdle)
	})
	b.Run("Resend", func(b *testing.B) {
		sequenced.BenchBrokerSuite(b, sequenced.BrokerChecks.Resend.Callable, brokerIdle)
	})
	b.Run("Listen", func(b *testing.B) {
		sequenced.BenchBrokerSuite(b, sequenced.BrokerChecks.Listen.Callable, brokerIdle)
	})
	b.Run("Fire", func(b *testing.B) {
		sequenced.BenchBrokerSuite(b, sequenced.BrokerChecks.Fire.Callable, brokerIdle)
	})
	b.Run("Refire", func(b *testing.B) {
		sequenced.BenchBrokerSuite(b, sequenced.BrokerChecks.Refire.Callable, brokerIdle)
	})
	b.Run("Tap", func(b *testing.B) {
		sequenced.BenchBrokerSuite(b, sequenced.BrokerChecks.Tap.Callable, brokerIdle)
	})
}

// watchesSubjectEmbedded is sequenced.Watches under the name of the field
// that watchesSubject embeds.
type watchesSubjectEmbedded = sequenced.Watches

// watchesSubject is a value of Watches that only TestWatchesSuite builds.
type watchesSubject struct{ watchesSubjectEmbedded }

// TestWatchesSuite checks the suite of Watches and the two forms of its
// application.
func TestWatchesSuite(t *testing.T) {
	t.Parallel()

	t.Run("WatchesSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sequenced.WatchesSuite())
	})

	t.Run("RunWatchesSuite", func(t *testing.T) {
		t.Parallel()
		var built *watchesSubject
		sequenced.RunWatchesSuite(
			t,
			func() *watchesSubject {
				built = &watchesSubject{}
				return built
			},
			conform.With(sequenced.WatchesCheck{
				ID:        "Watches/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.WatchesHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Watches(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.WatchesHandle) sequenced.Watches { return &watchesSubject{} },
			}),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunWatchesSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *watchesSubject
		sequenced.RunWatchesSuiteFromHandle(
			t,
			func(*sequenced.WatchesHandle) *watchesSubject {
				built = &watchesSubject{}
				return built
			},
			conform.With(sequenced.WatchesCheck{
				ID:        "Watches/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sequenced.WatchesHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sequenced.Watches(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sequenced.WatchesHandle) sequenced.Watches { return &watchesSubject{} },
			}),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Set.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Tick.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Peeked.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Finds.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Found.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Reads.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Written.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancels.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Cancelled.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closes.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Closed.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Fields.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Filled.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Follows.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sequenced.WatchesChecks.Touch.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// watchesIdle returns the planted defect of a ceiling of Watches with a
// spend of nothing, a subject that answers each call at once.
func watchesIdle(h *sequenced.WatchesHandle) sequenced.Watches {
	return sequenced.WatchesSuite().Calls[0].Over(h, func() {})
}

// TestWatchesSuiteBudgets applies the allocation ceilings of the suite of
// Watches to watchesIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestWatchesSuiteBudgets(t *testing.T) {
	sequenced.RunWatchesBudgets(
		t,
		watchesIdle,
		conform.Ceilings(
			sequenced.WatchesChecks.Watch.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Set.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Ticks.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Tick.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Peeks.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Peeked.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Finds.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Found.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Reads.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Written.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Cancels.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Cancelled.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Closes.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Closed.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Fields.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Filled.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Follows.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sequenced.WatchesChecks.Touch.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteWatch registers the fuzz target of Watches.Watch on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteWatch(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Watch.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteSet registers the fuzz target of Watches.Set on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteSet(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Set.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteTicks registers the fuzz target of Watches.Ticks on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteTicks(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Ticks.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteTick registers the fuzz target of Watches.Tick on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteTick(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Tick.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuitePeeks registers the fuzz target of Watches.Peeks on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuitePeeks(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Peeks.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuitePeeked registers the fuzz target of Watches.Peeked on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuitePeeked(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Peeked.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteFinds registers the fuzz target of Watches.Finds on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteFinds(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Finds.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteFound registers the fuzz target of Watches.Found on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteFound(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Found.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteReads registers the fuzz target of Watches.Reads on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteReads(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Reads.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteWritten registers the fuzz target of Watches.Written on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteWritten(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Written.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteCancels registers the fuzz target of Watches.Cancels on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteCancels(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Cancels.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteCancelled registers the fuzz target of Watches.Cancelled
// on watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteCancelled(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Cancelled.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteCloses registers the fuzz target of Watches.Closes on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteCloses(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Closes.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteClosed registers the fuzz target of Watches.Closed on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteClosed(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Closed.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteFields registers the fuzz target of Watches.Fields on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteFields(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Fields.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteFilled registers the fuzz target of Watches.Filled on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteFilled(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Filled.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteFollows registers the fuzz target of Watches.Follows on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteFollows(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Follows.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWatchesSuiteTouch registers the fuzz target of Watches.Touch on
// watchesIdle, with every check of the suite excused.
func FuzzWatchesSuiteTouch(f *testing.F) {
	sequenced.FuzzWatchesSuite(
		f,
		sequenced.WatchesChecks.Touch.Callable,
		watchesIdle,
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Watch.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Set.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Ticks.WatchSeesEachTrigger(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Tick.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeks.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Peeked.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Finds.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Found.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Reads.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Written.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancels.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Cancelled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closes.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Closed.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Fields.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Filled.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Follows.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sequenced.WatchesChecks.Touch.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkWatchesSuite measures each callable of the suite of Watches on
// watchesIdle.
func BenchmarkWatchesSuite(b *testing.B) {
	b.Run("Watch", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Watch.Callable, watchesIdle)
	})
	b.Run("Set", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Set.Callable, watchesIdle)
	})
	b.Run("Ticks", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Ticks.Callable, watchesIdle)
	})
	b.Run("Tick", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Tick.Callable, watchesIdle)
	})
	b.Run("Peeks", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Peeks.Callable, watchesIdle)
	})
	b.Run("Peeked", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Peeked.Callable, watchesIdle)
	})
	b.Run("Finds", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Finds.Callable, watchesIdle)
	})
	b.Run("Found", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Found.Callable, watchesIdle)
	})
	b.Run("Reads", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Reads.Callable, watchesIdle)
	})
	b.Run("Written", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Written.Callable, watchesIdle)
	})
	b.Run("Cancels", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Cancels.Callable, watchesIdle)
	})
	b.Run("Cancelled", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Cancelled.Callable, watchesIdle)
	})
	b.Run("Closes", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Closes.Callable, watchesIdle)
	})
	b.Run("Closed", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Closed.Callable, watchesIdle)
	})
	b.Run("Fields", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Fields.Callable, watchesIdle)
	})
	b.Run("Filled", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Filled.Callable, watchesIdle)
	})
	b.Run("Follows", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Follows.Callable, watchesIdle)
	})
	b.Run("Touch", func(b *testing.B) {
		sequenced.BenchWatchesSuite(b, sequenced.WatchesChecks.Touch.Callable, watchesIdle)
	})
}
