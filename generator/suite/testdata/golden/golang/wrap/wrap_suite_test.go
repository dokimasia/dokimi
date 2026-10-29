package wrap

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// reachSubjectEmbedded is Reach under the name of the field that
// reachSubject embeds.
type reachSubjectEmbedded = Reach

// reachSubject is a value of Reach that only TestReachSuite builds.
type reachSubject struct{ reachSubjectEmbedded }

// TestReachSuite checks the suite of Reach and the two forms of its
// application.
func TestReachSuite(t *testing.T) {
	t.Parallel()

	t.Run("ReachSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ReachSuite())
	})

	t.Run("RunReachSuite", func(t *testing.T) {
		t.Parallel()
		var built *reachSubject
		RunReachSuite(
			t,
			func() *reachSubject {
				built = &reachSubject{}
				return built
			},
			conform.With(ReachCheck{
				ID:        "Reach/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ReachHandle) {
					got := h.Subject()
					assert.Equal(h, got, Reach(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ReachHandle) Reach { return &reachSubject{} },
			}),
			conform.Excuse(
				ReachChecks.Current.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ReachChecks.Current.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ReachChecks.Current.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunReachSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *reachSubject
		RunReachSuiteFromHandle(
			t,
			func(*ReachHandle) *reachSubject {
				built = &reachSubject{}
				return built
			},
			conform.With(ReachCheck{
				ID:        "Reach/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ReachHandle) {
					got := h.Subject()
					assert.Equal(h, got, Reach(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ReachHandle) Reach { return &reachSubject{} },
			}),
			conform.Excuse(
				ReachChecks.Current.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ReachChecks.Current.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ReachChecks.Current.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// reachIdle returns the planted defect of a ceiling of Reach with a spend
// of nothing, a subject that answers each call at once.
func reachIdle(h *ReachHandle) Reach {
	return ReachSuite().Calls[0].Over(h, func() {})
}

// TestReachSuiteBudgets applies the allocation ceilings of the suite of
// Reach to reachIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestReachSuiteBudgets(t *testing.T) {
	RunReachBudgets(
		t,
		reachIdle,
		conform.Ceilings(ReachChecks.Store.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ReachChecks.Current.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ReachChecks.Gear.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzReachSuiteStore registers the fuzz target of Reach.Store on
// reachIdle, with every check of the suite excused.
func FuzzReachSuiteStore(f *testing.F) {
	FuzzReachSuite(
		f,
		ReachChecks.Store.Callable,
		reachIdle,
		conform.Excuse(
			ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzReachSuiteCurrent registers the fuzz target of Reach.Current on
// reachIdle, with every check of the suite excused.
func FuzzReachSuiteCurrent(f *testing.F) {
	FuzzReachSuite(
		f,
		ReachChecks.Current.Callable,
		reachIdle,
		conform.Excuse(
			ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzReachSuiteGear registers the fuzz target of Reach.Gear on reachIdle,
// with every check of the suite excused.
func FuzzReachSuiteGear(f *testing.F) {
	FuzzReachSuite(
		f,
		ReachChecks.Gear.Callable,
		reachIdle,
		conform.Excuse(
			ReachChecks.Current.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ReachChecks.Current.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkReachSuite measures each callable of the suite of Reach on
// reachIdle.
func BenchmarkReachSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		BenchReachSuite(b, ReachChecks.Store.Callable, reachIdle)
	})
	b.Run("Current", func(b *testing.B) {
		BenchReachSuite(b, ReachChecks.Current.Callable, reachIdle)
	})
	b.Run("Gear", func(b *testing.B) {
		BenchReachSuite(b, ReachChecks.Gear.Callable, reachIdle)
	})
}

// tunedSubjectEmbedded is Tuned under the name of the field that
// tunedSubject embeds.
type tunedSubjectEmbedded = Tuned

// tunedSubject is a value of Tuned that only TestTunedSuite builds.
type tunedSubject struct{ tunedSubjectEmbedded }

// TestTunedSuite checks the suite of Tuned and the two forms of its
// application.
func TestTunedSuite(t *testing.T) {
	t.Parallel()

	t.Run("TunedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, TunedSuite())
	})

	t.Run("RunTunedSuite", func(t *testing.T) {
		t.Parallel()
		var built *tunedSubject
		RunTunedSuite(
			t,
			func() *tunedSubject {
				built = &tunedSubject{}
				return built
			},
			conform.With(TunedCheck{
				ID:        "Tuned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TunedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Tuned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TunedHandle) Tuned { return &tunedSubject{} },
			}),
		)
	})

	t.Run("RunTunedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *tunedSubject
		RunTunedSuiteFromHandle(
			t,
			func(*TunedHandle) *tunedSubject {
				built = &tunedSubject{}
				return built
			},
			conform.With(TunedCheck{
				ID:        "Tuned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TunedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Tuned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TunedHandle) Tuned { return &tunedSubject{} },
			}),
		)
	})
}

// tunedIdle returns the planted defect of a ceiling of Tuned with a spend
// of nothing, a subject that answers each call at once.
func tunedIdle(h *TunedHandle) Tuned {
	return TunedSuite().Calls[0].Over(h, func() {})
}

// TestTunedSuiteBudgets applies the allocation ceilings of the suite of
// Tuned to tunedIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestTunedSuiteBudgets(t *testing.T) {
	RunTunedBudgets(
		t,
		tunedIdle,
		conform.Ceilings(TunedChecks.Tune.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TunedChecks.Watch.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TunedChecks.Turn.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzTunedSuiteTune registers the fuzz target of Tuned.Tune on tunedIdle,
// with every check of the suite excused.
func FuzzTunedSuiteTune(f *testing.F) {
	FuzzTunedSuite(
		f,
		TunedChecks.Tune.Callable,
		tunedIdle,
	)
}

// FuzzTunedSuiteWatch registers the fuzz target of Tuned.Watch on
// tunedIdle, with every check of the suite excused.
func FuzzTunedSuiteWatch(f *testing.F) {
	FuzzTunedSuite(
		f,
		TunedChecks.Watch.Callable,
		tunedIdle,
	)
}

// FuzzTunedSuiteTurn registers the fuzz target of Tuned.Turn on tunedIdle,
// with every check of the suite excused.
func FuzzTunedSuiteTurn(f *testing.F) {
	FuzzTunedSuite(
		f,
		TunedChecks.Turn.Callable,
		tunedIdle,
	)
}

// BenchmarkTunedSuite measures each callable of the suite of Tuned on
// tunedIdle.
func BenchmarkTunedSuite(b *testing.B) {
	b.Run("Tune", func(b *testing.B) {
		BenchTunedSuite(b, TunedChecks.Tune.Callable, tunedIdle)
	})
	b.Run("Watch", func(b *testing.B) {
		BenchTunedSuite(b, TunedChecks.Watch.Callable, tunedIdle)
	})
	b.Run("Turn", func(b *testing.B) {
		BenchTunedSuite(b, TunedChecks.Turn.Callable, tunedIdle)
	})
}

// limitedSubjectEmbedded is Limited under the name of the field that
// limitedSubject embeds.
type limitedSubjectEmbedded = Limited

// limitedSubject is a value of Limited that only TestLimitedSuite builds.
type limitedSubject struct{ limitedSubjectEmbedded }

// TestLimitedSuite checks the suite of Limited and the two forms of its
// application.
func TestLimitedSuite(t *testing.T) {
	t.Parallel()

	t.Run("LimitedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, LimitedSuite())
	})

	t.Run("RunLimitedSuite", func(t *testing.T) {
		t.Parallel()
		var built *limitedSubject
		RunLimitedSuite(
			t,
			func() *limitedSubject {
				built = &limitedSubject{}
				return built
			},
			conform.With(LimitedCheck{
				ID:        "Limited/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LimitedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Limited(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LimitedHandle) Limited { return &limitedSubject{} },
			}),
			conform.Excuse(
				LimitedChecks.Limit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimitedChecks.Measure.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimitedChecks.Phase.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLimitedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *limitedSubject
		RunLimitedSuiteFromHandle(
			t,
			func(*LimitedHandle) *limitedSubject {
				built = &limitedSubject{}
				return built
			},
			conform.With(LimitedCheck{
				ID:        "Limited/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LimitedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Limited(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LimitedHandle) Limited { return &limitedSubject{} },
			}),
			conform.Excuse(
				LimitedChecks.Limit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimitedChecks.Measure.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimitedChecks.Phase.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// limitedIdle returns the planted defect of a ceiling of Limited with a
// spend of nothing, a subject that answers each call at once.
func limitedIdle(h *LimitedHandle) Limited {
	return LimitedSuite().Calls[0].Over(h, func() {})
}

// TestLimitedSuiteBudgets applies the allocation ceilings of the suite of
// Limited to limitedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLimitedSuiteBudgets(t *testing.T) {
	RunLimitedBudgets(
		t,
		limitedIdle,
		conform.Ceilings(LimitedChecks.Limit.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Grip.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Clamp.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Frame.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Mount.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Span.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Pitch.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Measure.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LimitedChecks.Phase.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteLimit registers the fuzz target of Limited.Limit on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteLimit(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Limit.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteGrip registers the fuzz target of Limited.Grip on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteGrip(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Grip.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteClamp registers the fuzz target of Limited.Clamp on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteClamp(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Clamp.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteFrame registers the fuzz target of Limited.Frame on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteFrame(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Frame.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteMount registers the fuzz target of Limited.Mount on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteMount(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Mount.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteSpan registers the fuzz target of Limited.Span on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteSpan(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Span.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuitePitch registers the fuzz target of Limited.Pitch on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuitePitch(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Pitch.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteMeasure registers the fuzz target of Limited.Measure on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteMeasure(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Measure.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuitePhase registers the fuzz target of Limited.Phase on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuitePhase(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Phase.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimitedSuiteSense registers the fuzz target of Limited.Sense on
// limitedIdle, with every check of the suite excused.
func FuzzLimitedSuiteSense(f *testing.F) {
	FuzzLimitedSuite(
		f,
		LimitedChecks.Sense.Callable,
		limitedIdle,
		conform.Excuse(
			LimitedChecks.Limit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Measure.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimitedChecks.Phase.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLimitedSuite measures each callable of the suite of Limited on
// limitedIdle.
func BenchmarkLimitedSuite(b *testing.B) {
	b.Run("Limit", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Limit.Callable, limitedIdle)
	})
	b.Run("Grip", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Grip.Callable, limitedIdle)
	})
	b.Run("Clamp", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Clamp.Callable, limitedIdle)
	})
	b.Run("Frame", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Frame.Callable, limitedIdle)
	})
	b.Run("Mount", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Mount.Callable, limitedIdle)
	})
	b.Run("Span", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Span.Callable, limitedIdle)
	})
	b.Run("Pitch", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Pitch.Callable, limitedIdle)
	})
	b.Run("Measure", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Measure.Callable, limitedIdle)
	})
	b.Run("Phase", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Phase.Callable, limitedIdle)
	})
	b.Run("Sense", func(b *testing.B) {
		BenchLimitedSuite(b, LimitedChecks.Sense.Callable, limitedIdle)
	})
}

// plainSubjectEmbedded is PlainMethods under the name of the field that
// plainSubject embeds.
type plainSubjectEmbedded = PlainMethods

// plainSubject is a value of Plain that only TestPlainSuite builds.
type plainSubject struct{ plainSubjectEmbedded }

// TestPlainSuite checks the suite of Plain and the two forms of its
// application.
func TestPlainSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlainSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PlainSuite())
	})

	t.Run("RunPlainSuite", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		RunPlainSuite(
			t,
			func() *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, PlainMethods(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlainHandle) PlainMethods { return &plainSubject{} },
			}),
			conform.Excuse(
				PlainChecks.Name.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Name.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPlainSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		RunPlainSuiteFromHandle(
			t,
			func(*PlainHandle) *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, PlainMethods(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlainHandle) PlainMethods { return &plainSubject{} },
			}),
			conform.Excuse(
				PlainChecks.Name.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Name.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// plainIdle returns the planted defect of a ceiling of Plain with a spend
// of nothing, a subject that answers each call at once.
func plainIdle(h *PlainHandle) PlainMethods {
	return PlainSuite().Calls[0].Over(h, func() {})
}

// TestPlainSuiteBudgets applies the allocation ceilings of the suite of
// Plain to plainIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPlainSuiteBudgets(t *testing.T) {
	RunPlainBudgets(
		t,
		plainIdle,
		conform.Ceilings(PlainChecks.Name.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			PlainChecks.Name.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Name.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteName registers the fuzz target of Plain.Name on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteName(f *testing.F) {
	FuzzPlainSuite(
		f,
		PlainChecks.Name.Callable,
		plainIdle,
		conform.Excuse(
			PlainChecks.Name.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Name.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPlainSuite measures each callable of the suite of Plain on
// plainIdle.
func BenchmarkPlainSuite(b *testing.B) {
	b.Run("Name", func(b *testing.B) {
		BenchPlainSuite(b, PlainChecks.Name.Callable, plainIdle)
	})
}
