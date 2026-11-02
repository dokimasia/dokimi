package locks_test

import (
	"testing"

	"example.com/shop/locks"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// lockSubjectEmbedded is locks.Lock under the name of the field that
// lockSubject embeds.
type lockSubjectEmbedded = locks.Lock

// lockSubject is a value of Lock that only TestLockSuite builds.
type lockSubject struct{ lockSubjectEmbedded }

// TestLockSuite checks the suite of Lock and the two forms of its
// application.
func TestLockSuite(t *testing.T) {
	t.Parallel()

	t.Run("LockSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.LockSuite())
	})

	t.Run("RunLockSuite", func(t *testing.T) {
		t.Parallel()
		var built *lockSubject
		locks.RunLockSuite(
			t,
			func() *lockSubject {
				built = &lockSubject{}
				return built
			},
			conform.With(locks.LockCheck{
				ID:        "Lock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.LockHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Lock(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.LockHandle) locks.Lock { return &lockSubject{} },
			}),
			conform.Excuse(
				locks.LockChecks.Acquire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Release.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Release.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Held.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Held.Predicate(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Acquire.ReleaseEndsTheWait(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLockSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *lockSubject
		locks.RunLockSuiteFromHandle(
			t,
			func(*locks.LockHandle) *lockSubject {
				built = &lockSubject{}
				return built
			},
			conform.With(locks.LockCheck{
				ID:        "Lock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.LockHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Lock(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.LockHandle) locks.Lock { return &lockSubject{} },
			}),
			conform.Excuse(
				locks.LockChecks.Acquire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Release.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Release.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Held.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Held.Predicate(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockChecks.Acquire.ReleaseEndsTheWait(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("LockChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(locks.LockChecks.Acquire.ReleaseEndsTheWait()),
			"Lock.Acquire/release-ends-the-wait",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// lockIdle returns the planted defect of a ceiling of Lock with a spend of
// nothing, a subject that answers each call at once.
func lockIdle(h *locks.LockHandle) locks.Lock {
	return locks.LockSuite().Calls[0].Over(h, func() {})
}

// TestLockSuiteBudgets applies the allocation ceilings of the suite of
// Lock to lockIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestLockSuiteBudgets(t *testing.T) {
	locks.RunLockBudgets(
		t,
		lockIdle,
		conform.Ceilings(locks.LockChecks.Release.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(locks.LockChecks.Held.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			locks.LockChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Predicate(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockSuiteAcquire registers the fuzz target of Lock.Acquire on
// lockIdle, with every check of the suite excused.
func FuzzLockSuiteAcquire(f *testing.F) {
	locks.FuzzLockSuite(
		f,
		locks.LockChecks.Acquire.Callable,
		lockIdle,
		conform.Excuse(
			locks.LockChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Predicate(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockSuiteRelease registers the fuzz target of Lock.Release on
// lockIdle, with every check of the suite excused.
func FuzzLockSuiteRelease(f *testing.F) {
	locks.FuzzLockSuite(
		f,
		locks.LockChecks.Release.Callable,
		lockIdle,
		conform.Excuse(
			locks.LockChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Predicate(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockSuiteHeld registers the fuzz target of Lock.Held on lockIdle,
// with every check of the suite excused.
func FuzzLockSuiteHeld(f *testing.F) {
	locks.FuzzLockSuite(
		f,
		locks.LockChecks.Held.Callable,
		lockIdle,
		conform.Excuse(
			locks.LockChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Release.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Held.Predicate(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLockSuite measures each callable of the suite of Lock on
// lockIdle.
func BenchmarkLockSuite(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		locks.BenchLockSuite(b, locks.LockChecks.Release.Callable, lockIdle)
	})
	b.Run("Held", func(b *testing.B) {
		locks.BenchLockSuite(b, locks.LockChecks.Held.Callable, lockIdle)
	})
}

// BenchmarkLockSuiteParallel measures each callable of the suite of Lock
// that admits concurrent calls on lockIdle, on GOMAXPROCS goroutines at
// once.
func BenchmarkLockSuiteParallel(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		locks.BenchLockSuiteParallel(b, locks.LockChecks.Release.Callable, lockIdle)
	})
}

// boardSubjectEmbedded is locks.Board under the name of the field that
// boardSubject embeds.
type boardSubjectEmbedded = locks.Board

// boardSubject is a value of Board that only TestBoardSuite builds.
type boardSubject struct{ boardSubjectEmbedded }

// TestBoardSuite checks the suite of Board and the two forms of its
// application.
func TestBoardSuite(t *testing.T) {
	t.Parallel()

	t.Run("BoardSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.BoardSuite())
	})

	t.Run("RunBoardSuite", func(t *testing.T) {
		t.Parallel()
		var built *boardSubject
		locks.RunBoardSuite(
			t,
			func() *boardSubject {
				built = &boardSubject{}
				return built
			},
			conform.With(locks.BoardCheck{
				ID:        "Board/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.BoardHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Board(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.BoardHandle) locks.Board { return &boardSubject{} },
			}),
			conform.Excuse(
				locks.BoardChecks.Post.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Post.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Post.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Post.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.AgreesWithTheReference(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.Linearizable(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBoardSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *boardSubject
		locks.RunBoardSuiteFromHandle(
			t,
			func(*locks.BoardHandle) *boardSubject {
				built = &boardSubject{}
				return built
			},
			conform.With(locks.BoardCheck{
				ID:        "Board/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.BoardHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Board(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.BoardHandle) locks.Board { return &boardSubject{} },
			}),
			conform.Excuse(
				locks.BoardChecks.Post.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Post.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Post.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Post.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.AgreesWithTheReference(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.BoardChecks.Read.Linearizable(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// boardIdle returns the planted defect of a ceiling of Board with a spend
// of nothing, a subject that answers each call at once.
func boardIdle(h *locks.BoardHandle) locks.Board {
	return locks.BoardSuite().Calls[0].Over(h, func() {})
}

// TestBoardSuiteBudgets applies the allocation ceilings of the suite of
// Board to boardIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestBoardSuiteBudgets(t *testing.T) {
	locks.RunBoardBudgets(
		t,
		boardIdle,
		conform.Ceilings(locks.BoardChecks.Post.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(locks.BoardChecks.Read.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			locks.BoardChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.AgreesWithTheReference(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Linearizable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuitePost registers the fuzz target of Board.Post on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuitePost(f *testing.F) {
	locks.FuzzBoardSuite(
		f,
		locks.BoardChecks.Post.Callable,
		boardIdle,
		conform.Excuse(
			locks.BoardChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.AgreesWithTheReference(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Linearizable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteRead registers the fuzz target of Board.Read on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteRead(f *testing.F) {
	locks.FuzzBoardSuite(
		f,
		locks.BoardChecks.Read.Callable,
		boardIdle,
		conform.Excuse(
			locks.BoardChecks.Post.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Post.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.AgreesWithTheReference(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.BoardChecks.Read.Linearizable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBoardSuite measures each callable of the suite of Board on
// boardIdle.
func BenchmarkBoardSuite(b *testing.B) {
	b.Run("Post", func(b *testing.B) {
		locks.BenchBoardSuite(b, locks.BoardChecks.Post.Callable, boardIdle)
	})
	b.Run("Read", func(b *testing.B) {
		locks.BenchBoardSuite(b, locks.BoardChecks.Read.Callable, boardIdle)
	})
}

// BenchmarkBoardSuiteParallel measures each callable of the suite of Board
// that admits concurrent calls on boardIdle, on GOMAXPROCS goroutines at
// once.
func BenchmarkBoardSuiteParallel(b *testing.B) {
	b.Run("Post", func(b *testing.B) {
		locks.BenchBoardSuiteParallel(b, locks.BoardChecks.Post.Callable, boardIdle)
	})
	b.Run("Read", func(b *testing.B) {
		locks.BenchBoardSuiteParallel(b, locks.BoardChecks.Read.Callable, boardIdle)
	})
}

// shelfSubjectEmbedded is locks.Shelf under the name of the field that
// shelfSubject embeds.
type shelfSubjectEmbedded = locks.Shelf

// shelfSubject is a value of Shelf that only TestShelfSuite builds.
type shelfSubject struct{ shelfSubjectEmbedded }

// TestShelfSuite checks the suite of Shelf and the two forms of its
// application.
func TestShelfSuite(t *testing.T) {
	t.Parallel()

	t.Run("ShelfSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.ShelfSuite())
	})

	t.Run("RunShelfSuite", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		locks.RunShelfSuite(
			t,
			func() *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(locks.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.ShelfHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Shelf(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.ShelfHandle) locks.Shelf { return &shelfSubject{} },
			}),
			conform.Excuse(
				locks.ShelfChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.ReadsAgreeUnderConcurrency(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunShelfSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		locks.RunShelfSuiteFromHandle(
			t,
			func(*locks.ShelfHandle) *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(locks.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.ShelfHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Shelf(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.ShelfHandle) locks.Shelf { return &shelfSubject{} },
			}),
			conform.Excuse(
				locks.ShelfChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.ShelfChecks.Count.ReadsAgreeUnderConcurrency(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// shelfIdle returns the planted defect of a ceiling of Shelf with a spend
// of nothing, a subject that answers each call at once.
func shelfIdle(h *locks.ShelfHandle) locks.Shelf {
	return locks.ShelfSuite().Calls[0].Over(h, func() {})
}

// TestShelfSuiteBudgets applies the allocation ceilings of the suite of
// Shelf to shelfIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestShelfSuiteBudgets(t *testing.T) {
	locks.RunShelfBudgets(
		t,
		shelfIdle,
		conform.Ceilings(locks.ShelfChecks.Add.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(locks.ShelfChecks.Count.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			locks.ShelfChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ReadsAgreeUnderConcurrency(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteAdd registers the fuzz target of Shelf.Add on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteAdd(f *testing.F) {
	locks.FuzzShelfSuite(
		f,
		locks.ShelfChecks.Add.Callable,
		shelfIdle,
		conform.Excuse(
			locks.ShelfChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ReadsAgreeUnderConcurrency(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteCount registers the fuzz target of Shelf.Count on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteCount(f *testing.F) {
	locks.FuzzShelfSuite(
		f,
		locks.ShelfChecks.Count.Callable,
		shelfIdle,
		conform.Excuse(
			locks.ShelfChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.ShelfChecks.Count.ReadsAgreeUnderConcurrency(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkShelfSuite measures each callable of the suite of Shelf on
// shelfIdle.
func BenchmarkShelfSuite(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		locks.BenchShelfSuite(b, locks.ShelfChecks.Add.Callable, shelfIdle)
	})
	b.Run("Count", func(b *testing.B) {
		locks.BenchShelfSuite(b, locks.ShelfChecks.Count.Callable, shelfIdle)
	})
}

// BenchmarkShelfSuiteParallel measures each callable of the suite of Shelf
// that admits concurrent calls on shelfIdle, on GOMAXPROCS goroutines at
// once.
func BenchmarkShelfSuiteParallel(b *testing.B) {
	b.Run("Count", func(b *testing.B) {
		locks.BenchShelfSuiteParallel(b, locks.ShelfChecks.Count.Callable, shelfIdle)
	})
}

// flightSubjectEmbedded is locks.Flight under the name of the field that
// flightSubject embeds.
type flightSubjectEmbedded = locks.Flight

// flightSubject is a value of Flight that only TestFlightSuite builds.
type flightSubject struct{ flightSubjectEmbedded }

// TestFlightSuite checks the suite of Flight and the two forms of its
// application.
func TestFlightSuite(t *testing.T) {
	t.Parallel()

	t.Run("FlightSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.FlightSuite())
	})

	t.Run("RunFlightSuite", func(t *testing.T) {
		t.Parallel()
		var built *flightSubject
		locks.RunFlightSuite(
			t,
			func() *flightSubject {
				built = &flightSubject{}
				return built
			},
			conform.With(locks.FlightCheck{
				ID:        "Flight/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.FlightHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Flight(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.FlightHandle) locks.Flight { return &flightSubject{} },
			}),
			conform.Excuse(
				locks.FlightChecks.Do.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.ComputesOnce(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.ComputesOnceGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunFlightSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *flightSubject
		locks.RunFlightSuiteFromHandle(
			t,
			func(*locks.FlightHandle) *flightSubject {
				built = &flightSubject{}
				return built
			},
			conform.With(locks.FlightCheck{
				ID:        "Flight/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.FlightHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Flight(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.FlightHandle) locks.Flight { return &flightSubject{} },
			}),
			conform.Excuse(
				locks.FlightChecks.Do.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.ComputesOnce(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.FlightChecks.Do.ComputesOnceGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// flightIdle returns the planted defect of a ceiling of Flight with a
// spend of nothing, a subject that answers each call at once.
func flightIdle(h *locks.FlightHandle) locks.Flight {
	return locks.FlightSuite().Calls[0].Over(h, func() {})
}

// TestFlightSuiteBudgets applies the allocation ceilings of the suite of
// Flight to flightIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestFlightSuiteBudgets(t *testing.T) {
	locks.RunFlightBudgets(
		t,
		flightIdle,
		conform.Ceilings(locks.FlightChecks.Do.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			locks.FlightChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.ComputesOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.ComputesOnceGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzFlightSuiteDo registers the fuzz target of Flight.Do on flightIdle,
// with every check of the suite excused.
func FuzzFlightSuiteDo(f *testing.F) {
	locks.FuzzFlightSuite(
		f,
		locks.FlightChecks.Do.Callable,
		flightIdle,
		conform.Excuse(
			locks.FlightChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.ComputesOnce(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.FlightChecks.Do.ComputesOnceGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkFlightSuite measures each callable of the suite of Flight on
// flightIdle.
func BenchmarkFlightSuite(b *testing.B) {
	b.Run("Do", func(b *testing.B) {
		locks.BenchFlightSuite(b, locks.FlightChecks.Do.Callable, flightIdle)
	})
}

// lockerSubjectEmbedded is locks.Locker under the name of the field that
// lockerSubject embeds.
type lockerSubjectEmbedded = locks.Locker

// lockerSubject is a value of Locker that only TestLockerSuite builds.
type lockerSubject struct{ lockerSubjectEmbedded }

// TestLockerSuite checks the suite of Locker and the two forms of its
// application.
func TestLockerSuite(t *testing.T) {
	t.Parallel()

	t.Run("LockerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.LockerSuite())
	})

	t.Run("RunLockerSuite", func(t *testing.T) {
		t.Parallel()
		var built *lockerSubject
		locks.RunLockerSuite(
			t,
			func() *lockerSubject {
				built = &lockerSubject{}
				return built
			},
			conform.With(locks.LockerCheck{
				ID:        "Locker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.LockerHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Locker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.LockerHandle) locks.Locker { return &lockerSubject{} },
			}),
			conform.Excuse(
				locks.LockerChecks.Lock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.Lease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.LeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.ReleaseFreesTheLease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.ReleaseFreesTheLeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.OneHolderAtATime(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.OneHolderAtATimeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLockerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *lockerSubject
		locks.RunLockerSuiteFromHandle(
			t,
			func(*locks.LockerHandle) *lockerSubject {
				built = &lockerSubject{}
				return built
			},
			conform.With(locks.LockerCheck{
				ID:        "Locker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.LockerHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Locker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.LockerHandle) locks.Locker { return &lockerSubject{} },
			}),
			conform.Excuse(
				locks.LockerChecks.Lock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.Lease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.LeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.ReleaseFreesTheLease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.ReleaseFreesTheLeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.OneHolderAtATime(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Lock.OneHolderAtATimeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.LockerChecks.Unlock.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// lockerIdle returns the planted defect of a ceiling of Locker with a
// spend of nothing, a subject that answers each call at once.
func lockerIdle(h *locks.LockerHandle) locks.Locker {
	return locks.LockerSuite().Calls[0].Over(h, func() {})
}

// TestLockerSuiteBudgets applies the allocation ceilings of the suite of
// Locker to lockerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLockerSuiteBudgets(t *testing.T) {
	locks.RunLockerBudgets(
		t,
		lockerIdle,
		conform.Ceilings(locks.LockerChecks.Lock.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			locks.LockerChecks.Unlock.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			locks.LockerChecks.Lock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockerSuiteLock registers the fuzz target of Locker.Lock on
// lockerIdle, with every check of the suite excused.
func FuzzLockerSuiteLock(f *testing.F) {
	locks.FuzzLockerSuite(
		f,
		locks.LockerChecks.Lock.Callable,
		lockerIdle,
		conform.Excuse(
			locks.LockerChecks.Lock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockerSuiteUnlock registers the fuzz target of Locker.Unlock on
// lockerIdle, with every check of the suite excused.
func FuzzLockerSuiteUnlock(f *testing.F) {
	locks.FuzzLockerSuite(
		f,
		locks.LockerChecks.Unlock.Callable,
		lockerIdle,
		conform.Excuse(
			locks.LockerChecks.Lock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Lock.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.LockerChecks.Unlock.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLockerSuite measures each callable of the suite of Locker on
// lockerIdle.
func BenchmarkLockerSuite(b *testing.B) {
	b.Run("Lock", func(b *testing.B) {
		locks.BenchLockerSuite(b, locks.LockerChecks.Lock.Callable, lockerIdle)
	})
	b.Run("Unlock", func(b *testing.B) {
		locks.BenchLockerSuite(b, locks.LockerChecks.Unlock.Callable, lockerIdle)
	})
}

// gateSubjectEmbedded is locks.Gate under the name of the field that
// gateSubject embeds.
type gateSubjectEmbedded = locks.Gate

// gateSubject is a value of Gate that only TestGateSuite builds.
type gateSubject struct{ gateSubjectEmbedded }

// TestGateSuite checks the suite of Gate and the two forms of its
// application.
func TestGateSuite(t *testing.T) {
	t.Parallel()

	t.Run("GateSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.GateSuite())
	})

	t.Run("RunGateSuite", func(t *testing.T) {
		t.Parallel()
		var built *gateSubject
		locks.RunGateSuite(
			t,
			func() *gateSubject {
				built = &gateSubject{}
				return built
			},
			conform.With(locks.GateCheck{
				ID:        "Gate/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.GateHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Gate(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.GateHandle) locks.Gate { return &gateSubject{} },
			}),
			conform.Excuse(
				locks.GateChecks.Pass.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunGateSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *gateSubject
		locks.RunGateSuiteFromHandle(
			t,
			func(*locks.GateHandle) *gateSubject {
				built = &gateSubject{}
				return built
			},
			conform.With(locks.GateCheck{
				ID:        "Gate/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.GateHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Gate(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.GateHandle) locks.Gate { return &gateSubject{} },
			}),
			conform.Excuse(
				locks.GateChecks.Pass.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// gateIdle returns the planted defect of a ceiling of Gate with a spend of
// nothing, a subject that answers each call at once.
func gateIdle(h *locks.GateHandle) locks.Gate {
	return locks.GateSuite().Calls[0].Over(h, func() {})
}

// TestGateSuiteBudgets applies the allocation ceilings of the suite of
// Gate to gateIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestGateSuiteBudgets(t *testing.T) {
	locks.RunGateBudgets(
		t,
		gateIdle,
		conform.Excuse(
			locks.GateChecks.Pass.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzGateSuitePass registers the fuzz target of Gate.Pass on gateIdle,
// with every check of the suite excused.
func FuzzGateSuitePass(f *testing.F) {
	locks.FuzzGateSuite(
		f,
		locks.GateChecks.Pass.Callable,
		gateIdle,
		conform.Excuse(
			locks.GateChecks.Pass.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// TestAwaitSuite checks the suite of Await and the two forms of its
// application.
func TestAwaitSuite(t *testing.T) {
	t.Parallel()

	t.Run("AwaitSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.AwaitSuite())
	})

	t.Run("Await", func(t *testing.T) {
		t.Parallel()
		locks.RunAwaitSuite(t, locks.Await)
	})

	t.Run("RunAwaitSuite", func(t *testing.T) {
		t.Parallel()
		locks.RunAwaitSuite(
			t,
			locks.Await,
			conform.With(locks.AwaitCheck{
				ID:        "Await/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *locks.AwaitHandle) {
					assert.Equal(h, h.Subject(), locks.AwaitFunc(locks.Await), h.Statement())
				},
				Defect: func(*locks.AwaitHandle) locks.AwaitFunc { return nil },
			}),
			conform.Excuse(
				locks.AwaitChecks.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunAwaitSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		locks.RunAwaitSuiteFromHandle(
			t,
			func(*locks.AwaitHandle) locks.AwaitFunc {
				return locks.Await
			},
			conform.With(locks.AwaitCheck{
				ID:        "Await/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *locks.AwaitHandle) {
					assert.Equal(h, h.Subject(), locks.AwaitFunc(locks.Await), h.Statement())
				},
				Defect: func(*locks.AwaitHandle) locks.AwaitFunc { return nil },
			}),
			conform.Excuse(
				locks.AwaitChecks.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// FuzzAwait fuzzes the property checks of Await.
func FuzzAwait(f *testing.F) {
	locks.FuzzAwaitSuite(f, locks.Await)
}

// workerSubjectEmbedded is locks.Worker under the name of the field that
// workerSubject embeds.
type workerSubjectEmbedded = locks.Worker

// workerSubject is a value of Worker that only TestWorkerSuite builds.
type workerSubject struct{ workerSubjectEmbedded }

// TestWorkerSuite checks the suite of Worker and the two forms of its
// application.
func TestWorkerSuite(t *testing.T) {
	t.Parallel()

	t.Run("WorkerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.WorkerSuite())
	})

	t.Run("RunWorkerSuite", func(t *testing.T) {
		t.Parallel()
		var built *workerSubject
		locks.RunWorkerSuite(
			t,
			func() *workerSubject {
				built = &workerSubject{}
				return built
			},
			conform.With(locks.WorkerCheck{
				ID:        "Worker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.WorkerHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Worker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.WorkerHandle) locks.Worker { return &workerSubject{} },
			}),
			conform.Excuse(
				locks.WorkerChecks.Stop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.WorkerChecks.Stop.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunWorkerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *workerSubject
		locks.RunWorkerSuiteFromHandle(
			t,
			func(*locks.WorkerHandle) *workerSubject {
				built = &workerSubject{}
				return built
			},
			conform.With(locks.WorkerCheck{
				ID:        "Worker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *locks.WorkerHandle) {
					got := h.Subject()
					assert.Equal(h, got, locks.Worker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*locks.WorkerHandle) locks.Worker { return &workerSubject{} },
			}),
			conform.Excuse(
				locks.WorkerChecks.Stop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				locks.WorkerChecks.Stop.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// workerIdle returns the planted defect of a ceiling of Worker with a
// spend of nothing, a subject that answers each call at once.
func workerIdle(h *locks.WorkerHandle) locks.Worker {
	return locks.WorkerSuite().Calls[0].Over(h, func() {})
}

// TestWorkerSuiteBudgets applies the allocation ceilings of the suite of
// Worker to workerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestWorkerSuiteBudgets(t *testing.T) {
	locks.RunWorkerBudgets(
		t,
		workerIdle,
		conform.Ceilings(locks.WorkerChecks.Stop.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			locks.WorkerChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.WorkerChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWorkerSuiteRun registers the fuzz target of Worker.Run on
// workerIdle, with every check of the suite excused.
func FuzzWorkerSuiteRun(f *testing.F) {
	locks.FuzzWorkerSuite(
		f,
		locks.WorkerChecks.Run.Callable,
		workerIdle,
		conform.Excuse(
			locks.WorkerChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.WorkerChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWorkerSuiteStop registers the fuzz target of Worker.Stop on
// workerIdle, with every check of the suite excused.
func FuzzWorkerSuiteStop(f *testing.F) {
	locks.FuzzWorkerSuite(
		f,
		locks.WorkerChecks.Stop.Callable,
		workerIdle,
		conform.Excuse(
			locks.WorkerChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			locks.WorkerChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkWorkerSuite measures each callable of the suite of Worker on
// workerIdle.
func BenchmarkWorkerSuite(b *testing.B) {
	b.Run("Stop", func(b *testing.B) {
		locks.BenchWorkerSuite(b, locks.WorkerChecks.Stop.Callable, workerIdle)
	})
}

// TestServeSuite checks the suite of Serve and the two forms of its
// application.
func TestServeSuite(t *testing.T) {
	t.Parallel()

	t.Run("ServeSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, locks.ServeSuite())
	})

	t.Run("Serve", func(t *testing.T) {
		t.Parallel()
		locks.RunServeSuite(t, locks.Serve)
	})

	t.Run("RunServeSuite", func(t *testing.T) {
		t.Parallel()
		locks.RunServeSuite(
			t,
			locks.Serve,
			conform.With(locks.ServeCheck{
				ID:        "Serve/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *locks.ServeHandle) {
					assert.Equal(h, h.Subject(), locks.ServeFunc(locks.Serve), h.Statement())
				},
				Defect: func(*locks.ServeHandle) locks.ServeFunc { return nil },
			}),
		)
	})

	t.Run("RunServeSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		locks.RunServeSuiteFromHandle(
			t,
			func(*locks.ServeHandle) locks.ServeFunc {
				return locks.Serve
			},
			conform.With(locks.ServeCheck{
				ID:        "Serve/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *locks.ServeHandle) {
					assert.Equal(h, h.Subject(), locks.ServeFunc(locks.Serve), h.Statement())
				},
				Defect: func(*locks.ServeHandle) locks.ServeFunc { return nil },
			}),
		)
	})
}

// FuzzServe fuzzes the property checks of Serve.
func FuzzServe(f *testing.F) {
	locks.FuzzServeSuite(f, locks.Serve)
}

// TestMutexConformsToLock applies the checks of the suite of Lock to
// Mutex.
func TestMutexConformsToLock(t *testing.T) {
	t.Parallel()

	locks.RunLockSuite(t, locks.NewMutex)
}

// BenchmarkMutexConformsToLock measures each callable of the suite of Lock
// on Mutex.
func BenchmarkMutexConformsToLock(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		locks.BenchLockSuite(
			b,
			locks.LockChecks.Release.Callable,
			func(*locks.LockHandle) *locks.Mutex { return locks.NewMutex() },
		)
	})
	b.Run("Held", func(b *testing.B) {
		locks.BenchLockSuite(
			b,
			locks.LockChecks.Held.Callable,
			func(*locks.LockHandle) *locks.Mutex { return locks.NewMutex() },
		)
	})
}

// BenchmarkMutexConformsToLockParallel measures each callable of the suite
// of Lock that admits concurrent calls on Mutex, on GOMAXPROCS goroutines
// at once.
func BenchmarkMutexConformsToLockParallel(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		locks.BenchLockSuiteParallel(
			b,
			locks.LockChecks.Release.Callable,
			func(*locks.LockHandle) *locks.Mutex { return locks.NewMutex() },
		)
	})
}

// FuzzMutexConformsToLockAcquire fuzzes the property checks of
// Lock.Acquire on Mutex.
func FuzzMutexConformsToLockAcquire(f *testing.F) {
	locks.FuzzLockSuite(
		f,
		locks.LockChecks.Acquire.Callable,
		func(*locks.LockHandle) *locks.Mutex { return locks.NewMutex() },
	)
}

// FuzzMutexConformsToLockRelease fuzzes the property checks of
// Lock.Release on Mutex.
func FuzzMutexConformsToLockRelease(f *testing.F) {
	locks.FuzzLockSuite(
		f,
		locks.LockChecks.Release.Callable,
		func(*locks.LockHandle) *locks.Mutex { return locks.NewMutex() },
	)
}

// FuzzMutexConformsToLockHeld fuzzes the property checks of Lock.Held on
// Mutex.
func FuzzMutexConformsToLockHeld(f *testing.F) {
	locks.FuzzLockSuite(
		f,
		locks.LockChecks.Held.Callable,
		func(*locks.LockHandle) *locks.Mutex { return locks.NewMutex() },
	)
}

// TestTurnstileConformsToGate applies the checks of the suite of Gate to
// Turnstile.
func TestTurnstileConformsToGate(t *testing.T) {
	t.Parallel()

	locks.RunGateSuite(t, locks.NewTurnstile)
}

// FuzzTurnstileConformsToGatePass fuzzes the property checks of Gate.Pass
// on Turnstile.
func FuzzTurnstileConformsToGatePass(f *testing.F) {
	locks.FuzzGateSuite(
		f,
		locks.GateChecks.Pass.Callable,
		func(*locks.GateHandle) *locks.Turnstile { return locks.NewTurnstile() },
	)
}

// TestLoopConformsToWorker applies the checks of the suite of Worker to
// Loop.
func TestLoopConformsToWorker(t *testing.T) {
	t.Parallel()

	locks.RunWorkerSuite(t, locks.NewLoop)
}

// BenchmarkLoopConformsToWorker measures each callable of the suite of
// Worker on Loop.
func BenchmarkLoopConformsToWorker(b *testing.B) {
	b.Run("Stop", func(b *testing.B) {
		locks.BenchWorkerSuite(
			b,
			locks.WorkerChecks.Stop.Callable,
			func(*locks.WorkerHandle) *locks.Loop { return locks.NewLoop() },
		)
	})
}

// FuzzLoopConformsToWorkerRun fuzzes the property checks of Worker.Run on
// Loop.
func FuzzLoopConformsToWorkerRun(f *testing.F) {
	locks.FuzzWorkerSuite(
		f,
		locks.WorkerChecks.Run.Callable,
		func(*locks.WorkerHandle) *locks.Loop { return locks.NewLoop() },
	)
}

// FuzzLoopConformsToWorkerStop fuzzes the property checks of Worker.Stop
// on Loop.
func FuzzLoopConformsToWorkerStop(f *testing.F) {
	locks.FuzzWorkerSuite(
		f,
		locks.WorkerChecks.Stop.Callable,
		func(*locks.WorkerHandle) *locks.Loop { return locks.NewLoop() },
	)
}
