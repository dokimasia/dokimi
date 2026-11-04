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
			conform.ExcuseSuite[sequenced.Counter](
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
			conform.ExcuseSuite[sequenced.Counter](
				"the test applies the check of its option alone",
			),
		)
	})
}

// counterSuite is the suite of Counter that the tests of the planted
// defect of a ceiling apply.
var counterSuite = sequenced.CounterSuite()

// counterIdle returns the planted defect of a ceiling of Counter with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func counterIdle(h *sequenced.CounterHandle) sequenced.Counter {
	return counterSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Counter](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Counter](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Counter](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Counter](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Counter](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Counter](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Counter](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Notes]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Notes]("the test applies the check of its option alone"),
		)
	})
}

// notesSuite is the suite of Notes that the tests of the planted defect of
// a ceiling apply.
var notesSuite = sequenced.NotesSuite()

// notesIdle returns the planted defect of a ceiling of Notes with a spend
// of nothing. Each call of the defect returns its zero results at once.
func notesIdle(h *sequenced.NotesHandle) sequenced.Notes {
	return notesSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Notes](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Set]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Set]("the test applies the check of its option alone"),
		)
	})
}

// setSuite is the suite of Set that the tests of the planted defect of a
// ceiling apply.
var setSuite = sequenced.SetSuite()

// setIdle returns the planted defect of a ceiling of Set with a spend of
// nothing. Each call of the defect returns its zero results at once.
func setIdle(h *sequenced.SetHandle) sequenced.Set {
	return setSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Set](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Set](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Set](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Set](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Set](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Set](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.TallyMethods](
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
			conform.ExcuseSuite[sequenced.TallyMethods](
				"the test applies the check of its option alone",
			),
		)
	})
}

// tallySuite is the suite of Tally that the tests of the planted defect of
// a ceiling apply.
var tallySuite = sequenced.TallySuite()

// tallyIdle returns the planted defect of a ceiling of Tally with a spend
// of nothing. Each call of the defect returns its zero results at once.
func tallyIdle(h *sequenced.TallyHandle) sequenced.TallyMethods {
	return tallySuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.TallyMethods](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.TallyMethods](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.TallyMethods](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.TallyMethods](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Pile]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Pile]("the test applies the check of its option alone"),
		)
	})
}

// pileSuite is the suite of Pile that the tests of the planted defect of a
// ceiling apply.
var pileSuite = sequenced.PileSuite()

// pileIdle returns the planted defect of a ceiling of Pile with a spend of
// nothing. Each call of the defect returns its zero results at once.
func pileIdle(h *sequenced.PileHandle) sequenced.Pile {
	return pileSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Pile](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Pile](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Pile](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Pile](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Pile](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Run]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Run]("the test applies the check of its option alone"),
		)
	})
}

// runSuite is the suite of Run that the tests of the planted defect of a
// ceiling apply.
var runSuite = sequenced.RunSuite()

// runIdle returns the planted defect of a ceiling of Run with a spend of
// nothing. Each call of the defect returns its zero results at once.
func runIdle(h *sequenced.RunHandle) sequenced.Run {
	return runSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Run](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Run](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Run](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Run](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Run](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Run](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Log]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Log]("the test applies the check of its option alone"),
		)
	})
}

// logSuite is the suite of Log that the tests of the planted defect of a
// ceiling apply.
var logSuite = sequenced.LogSuite()

// logIdle returns the planted defect of a ceiling of Log with a spend of
// nothing. Each call of the defect returns its zero results at once.
func logIdle(h *sequenced.LogHandle) sequenced.Log {
	return logSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Log](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Log](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Feed]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Feed]("the test applies the check of its option alone"),
		)
	})
}

// feedSuite is the suite of Feed that the tests of the planted defect of a
// ceiling apply.
var feedSuite = sequenced.FeedSuite()

// feedIdle returns the planted defect of a ceiling of Feed with a spend of
// nothing. Each call of the defect returns its zero results at once.
func feedIdle(h *sequenced.FeedHandle) sequenced.Feed {
	return feedSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Feed](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Feed](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Feed](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Feed](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Feed](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Register](
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
			conform.ExcuseSuite[sequenced.Register](
				"the test applies the check of its option alone",
			),
		)
	})
}

// registerSuite is the suite of Register that the tests of the planted
// defect of a ceiling apply.
var registerSuite = sequenced.RegisterSuite()

// registerIdle returns the planted defect of a ceiling of Register with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func registerIdle(h *sequenced.RegisterHandle) sequenced.Register {
	return registerSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Register](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Register](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Register](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.DoubleFunc](
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
			conform.ExcuseSuite[sequenced.DoubleFunc](
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
			conform.ExcuseSuite[sequenced.Cached]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Cached]("the test applies the check of its option alone"),
		)
	})
}

// cachedSuite is the suite of Cached that the tests of the planted defect
// of a ceiling apply.
var cachedSuite = sequenced.CachedSuite()

// cachedIdle returns the planted defect of a ceiling of Cached with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func cachedIdle(h *sequenced.CachedHandle) sequenced.Cached {
	return cachedSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Cached](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Cached](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Cached](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Cached](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Cached](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Index]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Index]("the test applies the check of its option alone"),
		)
	})
}

// indexSuite is the suite of Index that the tests of the planted defect of
// a ceiling apply.
var indexSuite = sequenced.IndexSuite()

// indexIdle returns the planted defect of a ceiling of Index with a spend
// of nothing. Each call of the defect returns its zero results at once.
func indexIdle(h *sequenced.IndexHandle) sequenced.Index {
	return indexSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Index](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Index](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Index](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Index](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Index](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Index](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Stamped](
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
			conform.ExcuseSuite[sequenced.Stamped](
				"the test applies the check of its option alone",
			),
		)
	})
}

// stampedSuite is the suite of Stamped that the tests of the planted
// defect of a ceiling apply.
var stampedSuite = sequenced.StampedSuite()

// stampedIdle returns the planted defect of a ceiling of Stamped with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func stampedIdle(h *sequenced.StampedHandle) sequenced.Stamped {
	return stampedSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Stamped](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stamped](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stamped](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Versioned](
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
			conform.ExcuseSuite[sequenced.Versioned](
				"the test applies the check of its option alone",
			),
		)
	})
}

// versionedSuite is the suite of Versioned that the tests of the planted
// defect of a ceiling apply.
var versionedSuite = sequenced.VersionedSuite()

// versionedIdle returns the planted defect of a ceiling of Versioned with
// a spend of nothing. Each call of the defect returns its zero results at
// once.
func versionedIdle(h *sequenced.VersionedHandle) sequenced.Versioned {
	return versionedSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Versioned](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Versioned](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Versioned](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Accounts](
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
			conform.ExcuseSuite[sequenced.Accounts](
				"the test applies the check of its option alone",
			),
		)
	})
}

// accountsSuite is the suite of Accounts that the tests of the planted
// defect of a ceiling apply.
var accountsSuite = sequenced.AccountsSuite()

// accountsIdle returns the planted defect of a ceiling of Accounts with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func accountsIdle(h *sequenced.AccountsHandle) sequenced.Accounts {
	return accountsSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Accounts](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Accounts](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Accounts](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Accounts](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Accounts](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Accounts](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Accounts](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Journal](
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
			conform.ExcuseSuite[sequenced.Journal](
				"the test applies the check of its option alone",
			),
		)
	})
}

// journalSuite is the suite of Journal that the tests of the planted
// defect of a ceiling apply.
var journalSuite = sequenced.JournalSuite()

// journalIdle returns the planted defect of a ceiling of Journal with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func journalIdle(h *sequenced.JournalHandle) sequenced.Journal {
	return journalSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Journal](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Journal](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Journal](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Stock]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Stock]("the test applies the check of its option alone"),
		)
	})
}

// stockSuite is the suite of Stock that the tests of the planted defect of
// a ceiling apply.
var stockSuite = sequenced.StockSuite()

// stockIdle returns the planted defect of a ceiling of Stock with a spend
// of nothing. Each call of the defect returns its zero results at once.
func stockIdle(h *sequenced.StockHandle) sequenced.Stock {
	return stockSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Stock](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Outbox]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Outbox]("the test applies the check of its option alone"),
		)
	})
}

// outboxSuite is the suite of Outbox that the tests of the planted defect
// of a ceiling apply.
var outboxSuite = sequenced.OutboxSuite()

// outboxIdle returns the planted defect of a ceiling of Outbox with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func outboxIdle(h *sequenced.OutboxHandle) sequenced.Outbox {
	return outboxSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Outbox](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Broker]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sequenced.Broker]("the test applies the check of its option alone"),
		)
	})
}

// brokerSuite is the suite of Broker that the tests of the planted defect
// of a ceiling apply.
var brokerSuite = sequenced.BrokerSuite()

// brokerIdle returns the planted defect of a ceiling of Broker with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func brokerIdle(h *sequenced.BrokerHandle) sequenced.Broker {
	return brokerSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Broker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Broker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Broker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Broker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Broker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Broker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Broker](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sequenced.Watches](
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
			conform.ExcuseSuite[sequenced.Watches](
				"the test applies the check of its option alone",
			),
		)
	})
}

// watchesSuite is the suite of Watches that the tests of the planted
// defect of a ceiling apply.
var watchesSuite = sequenced.WatchesSuite()

// watchesIdle returns the planted defect of a ceiling of Watches with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func watchesIdle(h *sequenced.WatchesHandle) sequenced.Watches {
	return watchesSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sequenced.Watches](
			"the planted defect of a ceiling returns the zero results of each call",
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
