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
