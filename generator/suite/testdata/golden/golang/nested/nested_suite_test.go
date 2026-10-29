package nested

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// courierSubjectEmbedded is Courier under the name of the field that
// courierSubject embeds.
type courierSubjectEmbedded = Courier

// courierSubject is a value of Courier that only TestCourierSuite builds.
type courierSubject struct{ courierSubjectEmbedded }

// TestCourierSuite checks the suite of Courier and the two forms of its
// application.
func TestCourierSuite(t *testing.T) {
	t.Parallel()

	t.Run("CourierSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, CourierSuite())
	})

	t.Run("RunCourierSuite", func(t *testing.T) {
		t.Parallel()
		var built *courierSubject
		RunCourierSuite(
			t,
			func() *courierSubject {
				built = &courierSubject{}
				return built
			},
			conform.With(CourierCheck{
				ID:        "Courier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CourierHandle) {
					got := h.Subject()
					assert.Equal(h, got, Courier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CourierHandle) Courier { return &courierSubject{} },
			}),
			conform.Excuse(
				CourierChecks.Ship.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CourierChecks.Ship.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CourierChecks.Send.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCourierSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *courierSubject
		RunCourierSuiteFromHandle(
			t,
			func(*CourierHandle) *courierSubject {
				built = &courierSubject{}
				return built
			},
			conform.With(CourierCheck{
				ID:        "Courier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CourierHandle) {
					got := h.Subject()
					assert.Equal(h, got, Courier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CourierHandle) Courier { return &courierSubject{} },
			}),
			conform.Excuse(
				CourierChecks.Ship.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CourierChecks.Ship.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CourierChecks.Send.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// courierIdle returns the planted defect of a ceiling of Courier with a
// spend of nothing, a subject that answers each call at once.
func courierIdle(h *CourierHandle) Courier {
	return CourierSuite().Calls[0].Over(h, func() {})
}

// TestCourierSuiteBudgets applies the allocation ceilings of the suite of
// Courier to courierIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCourierSuiteBudgets(t *testing.T) {
	RunCourierBudgets(
		t,
		courierIdle,
		conform.Ceilings(CourierChecks.Ship.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(CourierChecks.Send.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			CourierChecks.Ship.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CourierChecks.Ship.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CourierChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCourierSuiteShip registers the fuzz target of Courier.Ship on
// courierIdle, with every check of the suite excused.
func FuzzCourierSuiteShip(f *testing.F) {
	FuzzCourierSuite(
		f,
		CourierChecks.Ship.Callable,
		courierIdle,
		conform.Excuse(
			CourierChecks.Ship.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CourierChecks.Ship.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CourierChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCourierSuiteSend registers the fuzz target of Courier.Send on
// courierIdle, with every check of the suite excused.
func FuzzCourierSuiteSend(f *testing.F) {
	FuzzCourierSuite(
		f,
		CourierChecks.Send.Callable,
		courierIdle,
		conform.Excuse(
			CourierChecks.Ship.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CourierChecks.Ship.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CourierChecks.Send.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCourierSuite measures each callable of the suite of Courier on
// courierIdle.
func BenchmarkCourierSuite(b *testing.B) {
	b.Run("Ship", func(b *testing.B) {
		BenchCourierSuite(b, CourierChecks.Ship.Callable, courierIdle)
	})
	b.Run("Send", func(b *testing.B) {
		BenchCourierSuite(b, CourierChecks.Send.Callable, courierIdle)
	})
}

// trackerSubjectEmbedded is Tracker under the name of the field that
// trackerSubject embeds.
type trackerSubjectEmbedded = Tracker

// trackerSubject is a value of Tracker that only TestTrackerSuite builds.
type trackerSubject struct{ trackerSubjectEmbedded }

// TestTrackerSuite checks the suite of Tracker and the two forms of its
// application.
func TestTrackerSuite(t *testing.T) {
	t.Parallel()

	t.Run("TrackerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, TrackerSuite())
	})

	t.Run("RunTrackerSuite", func(t *testing.T) {
		t.Parallel()
		var built *trackerSubject
		RunTrackerSuite(
			t,
			func() *trackerSubject {
				built = &trackerSubject{}
				return built
			},
			conform.With(TrackerCheck{
				ID:        "Tracker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TrackerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Tracker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TrackerHandle) Tracker { return &trackerSubject{} },
			}),
		)
	})

	t.Run("RunTrackerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *trackerSubject
		RunTrackerSuiteFromHandle(
			t,
			func(*TrackerHandle) *trackerSubject {
				built = &trackerSubject{}
				return built
			},
			conform.With(TrackerCheck{
				ID:        "Tracker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TrackerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Tracker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TrackerHandle) Tracker { return &trackerSubject{} },
			}),
		)
	})
}

// trackerIdle returns the planted defect of a ceiling of Tracker with a
// spend of nothing, a subject that answers each call at once.
func trackerIdle(h *TrackerHandle) Tracker {
	return TrackerSuite().Calls[0].Over(h, func() {})
}

// TestTrackerSuiteBudgets applies the allocation ceilings of the suite of
// Tracker to trackerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestTrackerSuiteBudgets(t *testing.T) {
	RunTrackerBudgets(
		t,
		trackerIdle,
		conform.Ceilings(TrackerChecks.Walk.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TrackerChecks.Pick.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TrackerChecks.Play.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TrackerChecks.Link.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzTrackerSuiteRoute registers the fuzz target of Tracker.Route on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteRoute(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Route.Callable,
		trackerIdle,
	)
}

// FuzzTrackerSuiteFollow registers the fuzz target of Tracker.Follow on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteFollow(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Follow.Callable,
		trackerIdle,
	)
}

// FuzzTrackerSuiteWalk registers the fuzz target of Tracker.Walk on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteWalk(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Walk.Callable,
		trackerIdle,
	)
}

// FuzzTrackerSuiteVolley registers the fuzz target of Tracker.Volley on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteVolley(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Volley.Callable,
		trackerIdle,
	)
}

// FuzzTrackerSuiteTrack registers the fuzz target of Tracker.Track on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteTrack(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Track.Callable,
		trackerIdle,
	)
}

// FuzzTrackerSuitePick registers the fuzz target of Tracker.Pick on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuitePick(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Pick.Callable,
		trackerIdle,
	)
}

// FuzzTrackerSuitePlay registers the fuzz target of Tracker.Play on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuitePlay(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Play.Callable,
		trackerIdle,
	)
}

// FuzzTrackerSuiteLink registers the fuzz target of Tracker.Link on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteLink(f *testing.F) {
	FuzzTrackerSuite(
		f,
		TrackerChecks.Link.Callable,
		trackerIdle,
	)
}

// BenchmarkTrackerSuite measures each callable of the suite of Tracker on
// trackerIdle.
func BenchmarkTrackerSuite(b *testing.B) {
	b.Run("Route", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Route.Callable, trackerIdle)
	})
	b.Run("Follow", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Follow.Callable, trackerIdle)
	})
	b.Run("Walk", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Walk.Callable, trackerIdle)
	})
	b.Run("Volley", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Volley.Callable, trackerIdle)
	})
	b.Run("Track", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Track.Callable, trackerIdle)
	})
	b.Run("Pick", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Pick.Callable, trackerIdle)
	})
	b.Run("Play", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Play.Callable, trackerIdle)
	})
	b.Run("Link", func(b *testing.B) {
		BenchTrackerSuite(b, TrackerChecks.Link.Callable, trackerIdle)
	})
}

// TestHireSuite checks the suite of Hire and the two forms of its
// application.
func TestHireSuite(t *testing.T) {
	t.Parallel()

	t.Run("HireSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, HireSuite())
	})

	t.Run("Hire", func(t *testing.T) {
		t.Parallel()
		RunHireSuite(t, Hire)
	})

	t.Run("RunHireSuite", func(t *testing.T) {
		t.Parallel()
		RunHireSuite(
			t,
			Hire,
			conform.With(HireCheck{
				ID:        "Hire/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *HireHandle) {
					assert.Equal(h, h.Subject(), HireFunc(Hire), h.Statement())
				},
				Defect: func(*HireHandle) HireFunc { return nil },
			}),
		)
	})

	t.Run("RunHireSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		RunHireSuiteFromHandle(
			t,
			func(*HireHandle) HireFunc {
				return Hire
			},
			conform.With(HireCheck{
				ID:        "Hire/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *HireHandle) {
					assert.Equal(h, h.Subject(), HireFunc(Hire), h.Statement())
				},
				Defect: func(*HireHandle) HireFunc { return nil },
			}),
		)
	})
}

// BenchmarkHire measures Hire.
func BenchmarkHire(b *testing.B) {
	BenchHireSuite(b, Hire)
}

// FuzzHire fuzzes the property checks of Hire.
func FuzzHire(f *testing.F) {
	FuzzHireSuite(f, Hire)
}
