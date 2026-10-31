package wrap_test

import (
	"testing"

	"example.com/shop/wrap"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// reachSubjectEmbedded is wrap.Reach under the name of the field that
// reachSubject embeds.
type reachSubjectEmbedded = wrap.Reach

// reachSubject is a value of Reach that only TestReachSuite builds.
type reachSubject struct{ reachSubjectEmbedded }

// TestReachSuite checks the suite of Reach and the two forms of its
// application.
func TestReachSuite(t *testing.T) {
	t.Parallel()

	t.Run("ReachSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, wrap.ReachSuite())
	})

	t.Run("RunReachSuite", func(t *testing.T) {
		t.Parallel()
		var built *reachSubject
		wrap.RunReachSuite(
			t,
			func() *reachSubject {
				built = &reachSubject{}
				return built
			},
			conform.With(wrap.ReachCheck{
				ID:        "Reach/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.ReachHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Reach(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.ReachHandle) wrap.Reach { return &reachSubject{} },
			}),
			conform.Excuse(
				wrap.ReachChecks.Current.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.ReachChecks.Current.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.ReachChecks.Current.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunReachSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *reachSubject
		wrap.RunReachSuiteFromHandle(
			t,
			func(*wrap.ReachHandle) *reachSubject {
				built = &reachSubject{}
				return built
			},
			conform.With(wrap.ReachCheck{
				ID:        "Reach/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.ReachHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Reach(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.ReachHandle) wrap.Reach { return &reachSubject{} },
			}),
			conform.Excuse(
				wrap.ReachChecks.Current.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.ReachChecks.Current.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.ReachChecks.Current.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// reachIdle returns the planted defect of a ceiling of Reach with a spend
// of nothing, a subject that answers each call at once.
func reachIdle(h *wrap.ReachHandle) wrap.Reach {
	return wrap.ReachSuite().Calls[0].Over(h, func() {})
}

// TestReachSuiteBudgets applies the allocation ceilings of the suite of
// Reach to reachIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestReachSuiteBudgets(t *testing.T) {
	wrap.RunReachBudgets(
		t,
		reachIdle,
		conform.Ceilings(wrap.ReachChecks.Store.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.ReachChecks.Current.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.ReachChecks.Gear.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			wrap.ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzReachSuiteStore registers the fuzz target of Reach.Store on
// reachIdle, with every check of the suite excused.
func FuzzReachSuiteStore(f *testing.F) {
	wrap.FuzzReachSuite(
		f,
		wrap.ReachChecks.Store.Callable,
		reachIdle,
		conform.Excuse(
			wrap.ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzReachSuiteCurrent registers the fuzz target of Reach.Current on
// reachIdle, with every check of the suite excused.
func FuzzReachSuiteCurrent(f *testing.F) {
	wrap.FuzzReachSuite(
		f,
		wrap.ReachChecks.Current.Callable,
		reachIdle,
		conform.Excuse(
			wrap.ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzReachSuiteGear registers the fuzz target of Reach.Gear on reachIdle,
// with every check of the suite excused.
func FuzzReachSuiteGear(f *testing.F) {
	wrap.FuzzReachSuite(
		f,
		wrap.ReachChecks.Gear.Callable,
		reachIdle,
		conform.Excuse(
			wrap.ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkReachSuite measures each callable of the suite of Reach on
// reachIdle.
func BenchmarkReachSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		wrap.BenchReachSuite(b, wrap.ReachChecks.Store.Callable, reachIdle)
	})
	b.Run("Current", func(b *testing.B) {
		wrap.BenchReachSuite(b, wrap.ReachChecks.Current.Callable, reachIdle)
	})
	b.Run("Gear", func(b *testing.B) {
		wrap.BenchReachSuite(b, wrap.ReachChecks.Gear.Callable, reachIdle)
	})
}

// tunedSubjectEmbedded is wrap.Tuned under the name of the field that
// tunedSubject embeds.
type tunedSubjectEmbedded = wrap.Tuned

// tunedSubject is a value of Tuned that only TestTunedSuite builds.
type tunedSubject struct{ tunedSubjectEmbedded }

// TestTunedSuite checks the suite of Tuned and the two forms of its
// application.
func TestTunedSuite(t *testing.T) {
	t.Parallel()

	t.Run("TunedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, wrap.TunedSuite())
	})

	t.Run("RunTunedSuite", func(t *testing.T) {
		t.Parallel()
		var built *tunedSubject
		wrap.RunTunedSuite(
			t,
			func() *tunedSubject {
				built = &tunedSubject{}
				return built
			},
			conform.With(wrap.TunedCheck{
				ID:        "Tuned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.TunedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Tuned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.TunedHandle) wrap.Tuned { return &tunedSubject{} },
			}),
		)
	})

	t.Run("RunTunedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *tunedSubject
		wrap.RunTunedSuiteFromHandle(
			t,
			func(*wrap.TunedHandle) *tunedSubject {
				built = &tunedSubject{}
				return built
			},
			conform.With(wrap.TunedCheck{
				ID:        "Tuned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.TunedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Tuned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.TunedHandle) wrap.Tuned { return &tunedSubject{} },
			}),
		)
	})
}

// tunedIdle returns the planted defect of a ceiling of Tuned with a spend
// of nothing, a subject that answers each call at once.
func tunedIdle(h *wrap.TunedHandle) wrap.Tuned {
	return wrap.TunedSuite().Calls[0].Over(h, func() {})
}

// TestTunedSuiteBudgets applies the allocation ceilings of the suite of
// Tuned to tunedIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestTunedSuiteBudgets(t *testing.T) {
	wrap.RunTunedBudgets(
		t,
		tunedIdle,
		conform.Ceilings(wrap.TunedChecks.Tune.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.TunedChecks.Watch.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.TunedChecks.Turn.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzTunedSuiteTune registers the fuzz target of Tuned.Tune on tunedIdle,
// with every check of the suite excused.
func FuzzTunedSuiteTune(f *testing.F) {
	wrap.FuzzTunedSuite(
		f,
		wrap.TunedChecks.Tune.Callable,
		tunedIdle,
	)
}

// FuzzTunedSuiteWatch registers the fuzz target of Tuned.Watch on
// tunedIdle, with every check of the suite excused.
func FuzzTunedSuiteWatch(f *testing.F) {
	wrap.FuzzTunedSuite(
		f,
		wrap.TunedChecks.Watch.Callable,
		tunedIdle,
	)
}

// FuzzTunedSuiteTurn registers the fuzz target of Tuned.Turn on tunedIdle,
// with every check of the suite excused.
func FuzzTunedSuiteTurn(f *testing.F) {
	wrap.FuzzTunedSuite(
		f,
		wrap.TunedChecks.Turn.Callable,
		tunedIdle,
	)
}

// BenchmarkTunedSuite measures each callable of the suite of Tuned on
// tunedIdle.
func BenchmarkTunedSuite(b *testing.B) {
	b.Run("Tune", func(b *testing.B) {
		wrap.BenchTunedSuite(b, wrap.TunedChecks.Tune.Callable, tunedIdle)
	})
	b.Run("Watch", func(b *testing.B) {
		wrap.BenchTunedSuite(b, wrap.TunedChecks.Watch.Callable, tunedIdle)
	})
	b.Run("Turn", func(b *testing.B) {
		wrap.BenchTunedSuite(b, wrap.TunedChecks.Turn.Callable, tunedIdle)
	})
}

// limitedSubjectEmbedded is wrap.Limited under the name of the field that
// limitedSubject embeds.
type limitedSubjectEmbedded = wrap.Limited

// limitedSubject is a value of Limited that only TestLimitedSuite builds.
type limitedSubject struct{ limitedSubjectEmbedded }

// TestLimitedSuite checks the suite of Limited and the two forms of its
// application.
func TestLimitedSuite(t *testing.T) {
	t.Parallel()

	t.Run("LimitedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, wrap.LimitedSuite())
	})

	t.Run("RunLimitedSuite", func(t *testing.T) {
		t.Parallel()
		var built *limitedSubject
		wrap.RunLimitedSuite(
			t,
			func() *limitedSubject {
				built = &limitedSubject{}
				return built
			},
			conform.With(wrap.LimitedCheck{
				ID:        "Limited/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.LimitedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Limited(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.LimitedHandle) wrap.Limited { return &limitedSubject{} },
			}),
			conform.Excuse(
				wrap.LimitedChecks.Limit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.LimitedChecks.Measure.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.LimitedChecks.Phase.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLimitedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *limitedSubject
		wrap.RunLimitedSuiteFromHandle(
			t,
			func(*wrap.LimitedHandle) *limitedSubject {
				built = &limitedSubject{}
				return built
			},
			conform.With(wrap.LimitedCheck{
				ID:        "Limited/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.LimitedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Limited(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.LimitedHandle) wrap.Limited { return &limitedSubject{} },
			}),
			conform.Excuse(
				wrap.LimitedChecks.Limit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.LimitedChecks.Measure.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.LimitedChecks.Phase.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// limitedIdle returns the planted defect of a ceiling of Limited with a
// spend of nothing, a subject that answers each call at once.
func limitedIdle(h *wrap.LimitedHandle) wrap.Limited {
	return wrap.LimitedSuite().Calls[0].Over(h, func() {})
}

// TestLimitedSuiteBudgets applies the allocation ceilings of the suite of
// Limited to limitedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLimitedSuiteBudgets(t *testing.T) {
	wrap.RunLimitedBudgets(
		t,
		limitedIdle,
		conform.Ceilings(wrap.LimitedChecks.Limit.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.LimitedChecks.Grip.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.LimitedChecks.Clamp.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.LimitedChecks.Frame.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.LimitedChecks.Mount.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.LimitedChecks.Span.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(wrap.LimitedChecks.Pitch.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			wrap.LimitedChecks.Measure.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(wrap.LimitedChecks.Phase.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteLimit registers the fuzz target of Limited.Limit on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteLimit(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Limit.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteGrip registers the fuzz target of Limited.Grip on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteGrip(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Grip.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteClamp registers the fuzz target of Limited.Clamp on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteClamp(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Clamp.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteFrame registers the fuzz target of Limited.Frame on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteFrame(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Frame.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteMount registers the fuzz target of Limited.Mount on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteMount(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Mount.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteSpan registers the fuzz target of Limited.Span on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteSpan(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Span.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuitePitch registers the fuzz target of Limited.Pitch on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuitePitch(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Pitch.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteMeasure registers the fuzz target of Limited.Measure on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteMeasure(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Measure.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuitePhase registers the fuzz target of Limited.Phase on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuitePhase(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Phase.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteSense registers the fuzz target of Limited.Sense on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteSense(f *testing.F) {
	wrap.FuzzLimitedSuite(
		f,
		wrap.LimitedChecks.Sense.Callable,
		limitedIdle,
		conform.Excuse(
			wrap.LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLimitedSuite measures each callable of the suite of Limited on
// limitedIdle.
func BenchmarkLimitedSuite(b *testing.B) {
	b.Run("Limit", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Limit.Callable, limitedIdle)
	})
	b.Run("Grip", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Grip.Callable, limitedIdle)
	})
	b.Run("Clamp", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Clamp.Callable, limitedIdle)
	})
	b.Run("Frame", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Frame.Callable, limitedIdle)
	})
	b.Run("Mount", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Mount.Callable, limitedIdle)
	})
	b.Run("Span", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Span.Callable, limitedIdle)
	})
	b.Run("Pitch", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Pitch.Callable, limitedIdle)
	})
	b.Run("Measure", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Measure.Callable, limitedIdle)
	})
	b.Run("Phase", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Phase.Callable, limitedIdle)
	})
	b.Run("Sense", func(b *testing.B) {
		wrap.BenchLimitedSuite(b, wrap.LimitedChecks.Sense.Callable, limitedIdle)
	})
}

// plainSubjectEmbedded is wrap.PlainMethods under the name of the field
// that plainSubject embeds.
type plainSubjectEmbedded = wrap.PlainMethods

// plainSubject is a value of Plain that only TestPlainSuite builds.
type plainSubject struct{ plainSubjectEmbedded }

// TestPlainSuite checks the suite of Plain and the two forms of its
// application.
func TestPlainSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlainSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, wrap.PlainSuite())
	})

	t.Run("RunPlainSuite", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		wrap.RunPlainSuite(
			t,
			func() *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(wrap.PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.PlainHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						wrap.PlainMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*wrap.PlainHandle) wrap.PlainMethods { return &plainSubject{} },
			}),
			conform.Excuse(
				wrap.PlainChecks.Name.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.PlainChecks.Name.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPlainSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		wrap.RunPlainSuiteFromHandle(
			t,
			func(*wrap.PlainHandle) *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(wrap.PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.PlainHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						wrap.PlainMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*wrap.PlainHandle) wrap.PlainMethods { return &plainSubject{} },
			}),
			conform.Excuse(
				wrap.PlainChecks.Name.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				wrap.PlainChecks.Name.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// plainIdle returns the planted defect of a ceiling of Plain with a spend
// of nothing, a subject that answers each call at once.
func plainIdle(h *wrap.PlainHandle) wrap.PlainMethods {
	return wrap.PlainSuite().Calls[0].Over(h, func() {})
}

// TestPlainSuiteBudgets applies the allocation ceilings of the suite of
// Plain to plainIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPlainSuiteBudgets(t *testing.T) {
	wrap.RunPlainBudgets(
		t,
		plainIdle,
		conform.Ceilings(wrap.PlainChecks.Name.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			wrap.PlainChecks.Name.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.PlainChecks.Name.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteName registers the fuzz target of Plain.Name on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteName(f *testing.F) {
	wrap.FuzzPlainSuite(
		f,
		wrap.PlainChecks.Name.Callable,
		plainIdle,
		conform.Excuse(
			wrap.PlainChecks.Name.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			wrap.PlainChecks.Name.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPlainSuite measures each callable of the suite of Plain on
// plainIdle.
func BenchmarkPlainSuite(b *testing.B) {
	b.Run("Name", func(b *testing.B) {
		wrap.BenchPlainSuite(b, wrap.PlainChecks.Name.Callable, plainIdle)
	})
}
