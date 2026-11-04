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
			conform.ExcuseSuite[drawn.Tallies]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[drawn.Tallies]("the test applies the check of its option alone"),
		)
	})
}

// talliesSuite is the suite of Tallies that the tests of the planted
// defect of a ceiling apply.
var talliesSuite = drawn.TalliesSuite()

// talliesIdle returns the planted defect of a ceiling of Tallies with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func talliesIdle(h *drawn.TalliesHandle) drawn.Tallies {
	return talliesSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Tallies](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[drawn.Notes]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[drawn.Notes]("the test applies the check of its option alone"),
		)
	})
}

// notesSuite is the suite of Notes that the tests of the planted defect of
// a ceiling apply.
var notesSuite = drawn.NotesSuite()

// notesIdle returns the planted defect of a ceiling of Notes with a spend
// of nothing. Each call of the defect returns its zero results at once.
func notesIdle(h *drawn.NotesHandle) drawn.Notes {
	return notesSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[drawn.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[drawn.Hooked]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[drawn.Hooked]("the test applies the check of its option alone"),
		)
	})
}

// hookedSuite is the suite of Hooked that the tests of the planted defect
// of a ceiling apply.
var hookedSuite = drawn.HookedSuite()

// hookedIdle returns the planted defect of a ceiling of Hooked with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func hookedIdle(h *drawn.HookedHandle) drawn.Hooked {
	return hookedSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[drawn.Hooked](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Hooked](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Hooked](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[drawn.Lookups]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[drawn.Lookups]("the test applies the check of its option alone"),
		)
	})
}

// lookupsSuite is the suite of Lookups that the tests of the planted
// defect of a ceiling apply.
var lookupsSuite = drawn.LookupsSuite()

// lookupsIdle returns the planted defect of a ceiling of Lookups with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func lookupsIdle(h *drawn.LookupsHandle) drawn.Lookups {
	return lookupsSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Lookups](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[drawn.Labels]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[drawn.Labels]("the test applies the check of its option alone"),
		)
	})
}

// labelsSuite is the suite of Labels that the tests of the planted defect
// of a ceiling apply.
var labelsSuite = drawn.LabelsSuite()

// labelsIdle returns the planted defect of a ceiling of Labels with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func labelsIdle(h *drawn.LabelsHandle) drawn.Labels {
	return labelsSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[drawn.Labels](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Labels](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[drawn.Labels](
			"the planted defect of a ceiling returns the zero results of each call",
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
