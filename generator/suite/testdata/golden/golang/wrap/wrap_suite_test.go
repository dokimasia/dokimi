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
			conform.ExcuseSuite[wrap.Reach]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[wrap.Reach]("the test applies the check of its option alone"),
		)
	})
}

// reachSuite is the suite of Reach that the tests of the planted defect of
// a ceiling apply.
var reachSuite = wrap.ReachSuite()

// reachIdle returns the planted defect of a ceiling of Reach with a spend
// of nothing. Each call of the defect returns its zero results at once.
func reachIdle(h *wrap.ReachHandle) wrap.Reach {
	return reachSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[wrap.Reach](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Reach](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Reach](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Reach](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[wrap.Tuned]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[wrap.Tuned]("the test applies the check of its option alone"),
		)
	})
}

// tunedSuite is the suite of Tuned that the tests of the planted defect of
// a ceiling apply.
var tunedSuite = wrap.TunedSuite()

// tunedIdle returns the planted defect of a ceiling of Tuned with a spend
// of nothing. Each call of the defect returns its zero results at once.
func tunedIdle(h *wrap.TunedHandle) wrap.Tuned {
	return tunedSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[wrap.Tuned](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTunedSuiteTune registers the fuzz target of Tuned.Tune on tunedIdle,
// with every check of the suite excused.
func FuzzTunedSuiteTune(f *testing.F) {
	wrap.FuzzTunedSuite(
		f,
		wrap.TunedChecks.Tune.Callable,
		tunedIdle,
		conform.ExcuseSuite[wrap.Tuned](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTunedSuiteWatch registers the fuzz target of Tuned.Watch on
// tunedIdle, with every check of the suite excused.
func FuzzTunedSuiteWatch(f *testing.F) {
	wrap.FuzzTunedSuite(
		f,
		wrap.TunedChecks.Watch.Callable,
		tunedIdle,
		conform.ExcuseSuite[wrap.Tuned](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTunedSuiteTurn registers the fuzz target of Tuned.Turn on tunedIdle,
// with every check of the suite excused.
func FuzzTunedSuiteTurn(f *testing.F) {
	wrap.FuzzTunedSuite(
		f,
		wrap.TunedChecks.Turn.Callable,
		tunedIdle,
		conform.ExcuseSuite[wrap.Tuned](
			"the planted defect of a ceiling returns the zero results of each call",
		),
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
			conform.ExcuseSuite[wrap.Limited]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[wrap.Limited]("the test applies the check of its option alone"),
		)
	})
}

// limitedSuite is the suite of Limited that the tests of the planted
// defect of a ceiling apply.
var limitedSuite = wrap.LimitedSuite()

// limitedIdle returns the planted defect of a ceiling of Limited with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func limitedIdle(h *wrap.LimitedHandle) wrap.Limited {
	return limitedSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.Limited](
			"the planted defect of a ceiling returns the zero results of each call",
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

// playedSubjectEmbedded is wrap.Played under the name of the field that
// playedSubject embeds.
type playedSubjectEmbedded = wrap.Played

// playedSubject is a value of Played that only TestPlayedSuite builds.
type playedSubject struct{ playedSubjectEmbedded }

// TestPlayedSuite checks the suite of Played and the two forms of its
// application.
func TestPlayedSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlayedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, wrap.PlayedSuite())
	})

	t.Run("RunPlayedSuite", func(t *testing.T) {
		t.Parallel()
		var built *playedSubject
		wrap.RunPlayedSuite(
			t,
			func() *playedSubject {
				built = &playedSubject{}
				return built
			},
			conform.With(wrap.PlayedCheck{
				ID:        "Played/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.PlayedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Played(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.PlayedHandle) wrap.Played { return &playedSubject{} },
			}),
			conform.ExcuseSuite[wrap.Played]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunPlayedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *playedSubject
		wrap.RunPlayedSuiteFromHandle(
			t,
			func(*wrap.PlayedHandle) *playedSubject {
				built = &playedSubject{}
				return built
			},
			conform.With(wrap.PlayedCheck{
				ID:        "Played/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.PlayedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Played(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.PlayedHandle) wrap.Played { return &playedSubject{} },
			}),
			conform.ExcuseSuite[wrap.Played]("the test applies the check of its option alone"),
		)
	})
}

// playedSuite is the suite of Played that the tests of the planted defect
// of a ceiling apply.
var playedSuite = wrap.PlayedSuite()

// playedIdle returns the planted defect of a ceiling of Played with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func playedIdle(h *wrap.PlayedHandle) wrap.Played {
	return playedSuite.Calls[0].Over(h, func() {})
}

// TestPlayedSuiteBudgets applies the allocation ceilings of the suite of
// Played to playedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestPlayedSuiteBudgets(t *testing.T) {
	wrap.RunPlayedBudgets(
		t,
		playedIdle,
		conform.Ceilings(wrap.PlayedChecks.Play.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			wrap.PlayedChecks.StateAt.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[wrap.Played](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzPlayedSuitePlay registers the fuzz target of Played.Play on
// playedIdle, with every check of the suite excused.
func FuzzPlayedSuitePlay(f *testing.F) {
	wrap.FuzzPlayedSuite(
		f,
		wrap.PlayedChecks.Play.Callable,
		playedIdle,
		conform.ExcuseSuite[wrap.Played](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzPlayedSuiteStateAt registers the fuzz target of Played.StateAt on
// playedIdle, with every check of the suite excused.
func FuzzPlayedSuiteStateAt(f *testing.F) {
	wrap.FuzzPlayedSuite(
		f,
		wrap.PlayedChecks.StateAt.Callable,
		playedIdle,
		conform.ExcuseSuite[wrap.Played](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkPlayedSuite measures each callable of the suite of Played on
// playedIdle.
func BenchmarkPlayedSuite(b *testing.B) {
	b.Run("Play", func(b *testing.B) {
		wrap.BenchPlayedSuite(b, wrap.PlayedChecks.Play.Callable, playedIdle)
	})
	b.Run("StateAt", func(b *testing.B) {
		wrap.BenchPlayedSuite(b, wrap.PlayedChecks.StateAt.Callable, playedIdle)
	})
}

// gaugedSubjectEmbedded is wrap.Gauged under the name of the field that
// gaugedSubject embeds.
type gaugedSubjectEmbedded = wrap.Gauged

// gaugedSubject is a value of Gauged that only TestGaugedSuite builds.
type gaugedSubject struct{ gaugedSubjectEmbedded }

// TestGaugedSuite checks the suite of Gauged and the two forms of its
// application.
func TestGaugedSuite(t *testing.T) {
	t.Parallel()

	t.Run("GaugedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, wrap.GaugedSuite())
	})

	t.Run("RunGaugedSuite", func(t *testing.T) {
		t.Parallel()
		var built *gaugedSubject
		wrap.RunGaugedSuite(
			t,
			func() *gaugedSubject {
				built = &gaugedSubject{}
				return built
			},
			conform.With(wrap.GaugedCheck{
				ID:        "Gauged/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.GaugedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Gauged(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.GaugedHandle) wrap.Gauged { return &gaugedSubject{} },
			}),
			conform.ExcuseSuite[wrap.Gauged]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunGaugedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *gaugedSubject
		wrap.RunGaugedSuiteFromHandle(
			t,
			func(*wrap.GaugedHandle) *gaugedSubject {
				built = &gaugedSubject{}
				return built
			},
			conform.With(wrap.GaugedCheck{
				ID:        "Gauged/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *wrap.GaugedHandle) {
					got := h.Subject()
					assert.Equal(h, got, wrap.Gauged(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*wrap.GaugedHandle) wrap.Gauged { return &gaugedSubject{} },
			}),
			conform.ExcuseSuite[wrap.Gauged]("the test applies the check of its option alone"),
		)
	})
}

// gaugedSuite is the suite of Gauged that the tests of the planted defect
// of a ceiling apply.
var gaugedSuite = wrap.GaugedSuite()

// gaugedIdle returns the planted defect of a ceiling of Gauged with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func gaugedIdle(h *wrap.GaugedHandle) wrap.Gauged {
	return gaugedSuite.Calls[0].Over(h, func() {})
}

// TestGaugedSuiteBudgets applies the allocation ceilings of the suite of
// Gauged to gaugedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestGaugedSuiteBudgets(t *testing.T) {
	wrap.RunGaugedBudgets(
		t,
		gaugedIdle,
		conform.Ceilings(wrap.GaugedChecks.Fill.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			wrap.GaugedChecks.LevelAt.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[wrap.Gauged](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzGaugedSuiteFill registers the fuzz target of Gauged.Fill on
// gaugedIdle, with every check of the suite excused.
func FuzzGaugedSuiteFill(f *testing.F) {
	wrap.FuzzGaugedSuite(
		f,
		wrap.GaugedChecks.Fill.Callable,
		gaugedIdle,
		conform.ExcuseSuite[wrap.Gauged](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzGaugedSuiteLevelAt registers the fuzz target of Gauged.LevelAt on
// gaugedIdle, with every check of the suite excused.
func FuzzGaugedSuiteLevelAt(f *testing.F) {
	wrap.FuzzGaugedSuite(
		f,
		wrap.GaugedChecks.LevelAt.Callable,
		gaugedIdle,
		conform.ExcuseSuite[wrap.Gauged](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkGaugedSuite measures each callable of the suite of Gauged on
// gaugedIdle.
func BenchmarkGaugedSuite(b *testing.B) {
	b.Run("Fill", func(b *testing.B) {
		wrap.BenchGaugedSuite(b, wrap.GaugedChecks.Fill.Callable, gaugedIdle)
	})
	b.Run("LevelAt", func(b *testing.B) {
		wrap.BenchGaugedSuite(b, wrap.GaugedChecks.LevelAt.Callable, gaugedIdle)
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
			conform.ExcuseSuite[wrap.PlainMethods](
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
			conform.ExcuseSuite[wrap.PlainMethods](
				"the test applies the check of its option alone",
			),
		)
	})
}

// plainSuite is the suite of Plain that the tests of the planted defect of
// a ceiling apply.
var plainSuite = wrap.PlainSuite()

// plainIdle returns the planted defect of a ceiling of Plain with a spend
// of nothing. Each call of the defect returns its zero results at once.
func plainIdle(h *wrap.PlainHandle) wrap.PlainMethods {
	return plainSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[wrap.PlainMethods](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[wrap.PlainMethods](
			"the planted defect of a ceiling returns the zero results of each call",
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
