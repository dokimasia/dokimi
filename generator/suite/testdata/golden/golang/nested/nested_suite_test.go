package nested_test

import (
	"testing"

	"example.com/shop/nested"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// courierSubjectEmbedded is nested.Courier under the name of the field
// that courierSubject embeds.
type courierSubjectEmbedded = nested.Courier

// courierSubject is a value of Courier that only TestCourierSuite builds.
type courierSubject struct{ courierSubjectEmbedded }

// TestCourierSuite checks the suite of Courier and the two forms of its
// application.
func TestCourierSuite(t *testing.T) {
	t.Parallel()

	t.Run("CourierSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, nested.CourierSuite())
	})

	t.Run("RunCourierSuite", func(t *testing.T) {
		t.Parallel()
		var built *courierSubject
		nested.RunCourierSuite(
			t,
			func() *courierSubject {
				built = &courierSubject{}
				return built
			},
			conform.With(nested.CourierCheck{
				ID:        "Courier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *nested.CourierHandle) {
					got := h.Subject()
					assert.Equal(h, got, nested.Courier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*nested.CourierHandle) nested.Courier { return &courierSubject{} },
			}),
			conform.ExcuseSuite[nested.Courier]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunCourierSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *courierSubject
		nested.RunCourierSuiteFromHandle(
			t,
			func(*nested.CourierHandle) *courierSubject {
				built = &courierSubject{}
				return built
			},
			conform.With(nested.CourierCheck{
				ID:        "Courier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *nested.CourierHandle) {
					got := h.Subject()
					assert.Equal(h, got, nested.Courier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*nested.CourierHandle) nested.Courier { return &courierSubject{} },
			}),
			conform.ExcuseSuite[nested.Courier]("the test applies the check of its option alone"),
		)
	})
}

// courierSuite is the suite of Courier that the tests of the planted
// defect of a ceiling apply.
var courierSuite = nested.CourierSuite()

// courierIdle returns the planted defect of a ceiling of Courier with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func courierIdle(h *nested.CourierHandle) nested.Courier {
	return courierSuite.Calls[0].Over(h, func() {})
}

// TestCourierSuiteBudgets applies the allocation ceilings of the suite of
// Courier to courierIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCourierSuiteBudgets(t *testing.T) {
	nested.RunCourierBudgets(
		t,
		courierIdle,
		conform.Ceilings(
			nested.CourierChecks.Ship.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			nested.CourierChecks.Send.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[nested.Courier](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzCourierSuiteShip registers the fuzz target of Courier.Ship on
// courierIdle, with every check of the suite excused.
func FuzzCourierSuiteShip(f *testing.F) {
	nested.FuzzCourierSuite(
		f,
		nested.CourierChecks.Ship.Callable,
		courierIdle,
		conform.ExcuseSuite[nested.Courier](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzCourierSuiteSend registers the fuzz target of Courier.Send on
// courierIdle, with every check of the suite excused.
func FuzzCourierSuiteSend(f *testing.F) {
	nested.FuzzCourierSuite(
		f,
		nested.CourierChecks.Send.Callable,
		courierIdle,
		conform.ExcuseSuite[nested.Courier](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkCourierSuite measures each callable of the suite of Courier on
// courierIdle.
func BenchmarkCourierSuite(b *testing.B) {
	b.Run("Ship", func(b *testing.B) {
		nested.BenchCourierSuite(b, nested.CourierChecks.Ship.Callable, courierIdle)
	})
	b.Run("Send", func(b *testing.B) {
		nested.BenchCourierSuite(b, nested.CourierChecks.Send.Callable, courierIdle)
	})
}

// trackerSubjectEmbedded is nested.Tracker under the name of the field
// that trackerSubject embeds.
type trackerSubjectEmbedded = nested.Tracker

// trackerSubject is a value of Tracker that only TestTrackerSuite builds.
type trackerSubject struct{ trackerSubjectEmbedded }

// TestTrackerSuite checks the suite of Tracker and the two forms of its
// application.
func TestTrackerSuite(t *testing.T) {
	t.Parallel()

	t.Run("TrackerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, nested.TrackerSuite())
	})

	t.Run("RunTrackerSuite", func(t *testing.T) {
		t.Parallel()
		var built *trackerSubject
		nested.RunTrackerSuite(
			t,
			func() *trackerSubject {
				built = &trackerSubject{}
				return built
			},
			conform.With(nested.TrackerCheck{
				ID:        "Tracker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *nested.TrackerHandle) {
					got := h.Subject()
					assert.Equal(h, got, nested.Tracker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*nested.TrackerHandle) nested.Tracker { return &trackerSubject{} },
			}),
			conform.ExcuseSuite[nested.Tracker]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunTrackerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *trackerSubject
		nested.RunTrackerSuiteFromHandle(
			t,
			func(*nested.TrackerHandle) *trackerSubject {
				built = &trackerSubject{}
				return built
			},
			conform.With(nested.TrackerCheck{
				ID:        "Tracker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *nested.TrackerHandle) {
					got := h.Subject()
					assert.Equal(h, got, nested.Tracker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*nested.TrackerHandle) nested.Tracker { return &trackerSubject{} },
			}),
			conform.ExcuseSuite[nested.Tracker]("the test applies the check of its option alone"),
		)
	})
}

// trackerSuite is the suite of Tracker that the tests of the planted
// defect of a ceiling apply.
var trackerSuite = nested.TrackerSuite()

// trackerIdle returns the planted defect of a ceiling of Tracker with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func trackerIdle(h *nested.TrackerHandle) nested.Tracker {
	return trackerSuite.Calls[0].Over(h, func() {})
}

// TestTrackerSuiteBudgets applies the allocation ceilings of the suite of
// Tracker to trackerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestTrackerSuiteBudgets(t *testing.T) {
	nested.RunTrackerBudgets(
		t,
		trackerIdle,
		conform.Ceilings(
			nested.TrackerChecks.Walk.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			nested.TrackerChecks.Pick.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			nested.TrackerChecks.Play.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			nested.TrackerChecks.Link.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuiteRoute registers the fuzz target of Tracker.Route on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteRoute(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Route.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuiteFollow registers the fuzz target of Tracker.Follow on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteFollow(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Follow.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuiteWalk registers the fuzz target of Tracker.Walk on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteWalk(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Walk.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuiteVolley registers the fuzz target of Tracker.Volley on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteVolley(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Volley.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuiteTrack registers the fuzz target of Tracker.Track on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteTrack(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Track.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuitePick registers the fuzz target of Tracker.Pick on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuitePick(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Pick.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuitePlay registers the fuzz target of Tracker.Play on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuitePlay(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Play.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuiteLink registers the fuzz target of Tracker.Link on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteLink(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Link.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTrackerSuiteTick registers the fuzz target of Tracker.Tick on
// trackerIdle, with every check of the suite excused.
func FuzzTrackerSuiteTick(f *testing.F) {
	nested.FuzzTrackerSuite(
		f,
		nested.TrackerChecks.Tick.Callable,
		trackerIdle,
		conform.ExcuseSuite[nested.Tracker](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkTrackerSuite measures each callable of the suite of Tracker on
// trackerIdle.
func BenchmarkTrackerSuite(b *testing.B) {
	b.Run("Route", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Route.Callable, trackerIdle)
	})
	b.Run("Follow", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Follow.Callable, trackerIdle)
	})
	b.Run("Walk", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Walk.Callable, trackerIdle)
	})
	b.Run("Volley", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Volley.Callable, trackerIdle)
	})
	b.Run("Track", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Track.Callable, trackerIdle)
	})
	b.Run("Pick", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Pick.Callable, trackerIdle)
	})
	b.Run("Play", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Play.Callable, trackerIdle)
	})
	b.Run("Link", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Link.Callable, trackerIdle)
	})
	b.Run("Tick", func(b *testing.B) {
		nested.BenchTrackerSuite(b, nested.TrackerChecks.Tick.Callable, trackerIdle)
	})
}

// TestHireSuite checks the suite of Hire and the two forms of its
// application.
func TestHireSuite(t *testing.T) {
	t.Parallel()

	t.Run("HireSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, nested.HireSuite())
	})

	t.Run("Hire", func(t *testing.T) {
		t.Parallel()
		nested.RunHireSuite(t, nested.Hire)
	})

	t.Run("RunHireSuite", func(t *testing.T) {
		t.Parallel()
		nested.RunHireSuite(
			t,
			nested.Hire,
			conform.With(nested.HireCheck{
				ID:        "Hire/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *nested.HireHandle) {
					assert.Equal(h, h.Subject(), nested.HireFunc(nested.Hire), h.Statement())
				},
				Defect: func(*nested.HireHandle) nested.HireFunc { return nil },
			}),
			conform.ExcuseSuite[nested.HireFunc]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunHireSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		nested.RunHireSuiteFromHandle(
			t,
			func(*nested.HireHandle) nested.HireFunc {
				return nested.Hire
			},
			conform.With(nested.HireCheck{
				ID:        "Hire/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *nested.HireHandle) {
					assert.Equal(h, h.Subject(), nested.HireFunc(nested.Hire), h.Statement())
				},
				Defect: func(*nested.HireHandle) nested.HireFunc { return nil },
			}),
			conform.ExcuseSuite[nested.HireFunc]("the test applies the check of its option alone"),
		)
	})
}

// BenchmarkHire measures Hire.
func BenchmarkHire(b *testing.B) {
	nested.BenchHireSuite(b, nested.Hire)
}

// FuzzHire fuzzes the property checks of Hire.
func FuzzHire(f *testing.F) {
	nested.FuzzHireSuite(f, nested.Hire)
}
