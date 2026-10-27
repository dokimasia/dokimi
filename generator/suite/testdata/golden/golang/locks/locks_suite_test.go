package locks

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// lockSubjectEmbedded is Lock under the name of the field that lockSubject
// embeds.
type lockSubjectEmbedded = Lock

// lockSubject is a value of Lock that only TestLockSuite builds.
type lockSubject struct{ lockSubjectEmbedded }

// TestLockSuite checks the suite of Lock and the two forms of its
// application.
func TestLockSuite(t *testing.T) {
	t.Parallel()

	t.Run("LockSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, LockSuite())
	})

	t.Run("RunLockSuite", func(t *testing.T) {
		t.Parallel()
		var built *lockSubject
		RunLockSuite(
			t,
			func() *lockSubject {
				built = &lockSubject{}
				return built
			},
			conform.With(LockCheck{
				ID:        "Lock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LockHandle) {
					got := h.Subject()
					assert.Equal(h, got, Lock(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LockHandle) Lock { return &lockSubject{} },
			}),
			conform.Excuse(
				LockChecks.Acquire.ReleaseEndsTheWait(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLockSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *lockSubject
		RunLockSuiteFromHandle(
			t,
			func(*LockHandle) *lockSubject {
				built = &lockSubject{}
				return built
			},
			conform.With(LockCheck{
				ID:        "Lock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LockHandle) {
					got := h.Subject()
					assert.Equal(h, got, Lock(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LockHandle) Lock { return &lockSubject{} },
			}),
			conform.Excuse(
				LockChecks.Acquire.ReleaseEndsTheWait(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("LockChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(LockChecks.Acquire.ReleaseEndsTheWait()),
			"Lock.Acquire/release-ends-the-wait",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// lockIdle returns the planted defect of a ceiling of Lock with a spend of
// nothing, a subject that answers each call at once.
func lockIdle(h *LockHandle) Lock {
	return LockSuite().Calls[0].Over(h, func() {})
}

// TestLockSuiteBudgets applies the allocation ceilings of the suite of
// Lock to lockIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestLockSuiteBudgets(t *testing.T) {
	RunLockBudgets(
		t,
		lockIdle,
		conform.Ceilings(LockChecks.Release.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LockChecks.Held.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockSuiteAcquire registers the fuzz target of Lock.Acquire on
// lockIdle, with every check of the suite excused.
func FuzzLockSuiteAcquire(f *testing.F) {
	FuzzLockSuite(
		f,
		LockChecks.Acquire.Callable,
		lockIdle,
		conform.Excuse(
			LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockSuiteRelease registers the fuzz target of Lock.Release on
// lockIdle, with every check of the suite excused.
func FuzzLockSuiteRelease(f *testing.F) {
	FuzzLockSuite(
		f,
		LockChecks.Release.Callable,
		lockIdle,
		conform.Excuse(
			LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockSuiteHeld registers the fuzz target of Lock.Held on lockIdle,
// with every check of the suite excused.
func FuzzLockSuiteHeld(f *testing.F) {
	FuzzLockSuite(
		f,
		LockChecks.Held.Callable,
		lockIdle,
		conform.Excuse(
			LockChecks.Acquire.ReleaseEndsTheWait(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLockSuite measures each callable of the suite of Lock on
// lockIdle.
func BenchmarkLockSuite(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		BenchLockSuite(b, LockChecks.Release.Callable, lockIdle)
	})
	b.Run("Held", func(b *testing.B) {
		BenchLockSuite(b, LockChecks.Held.Callable, lockIdle)
	})
}

// BenchmarkLockSuiteParallel measures each callable of the suite of Lock
// that admits concurrent calls on lockIdle, on GOMAXPROCS goroutines at
// once.
func BenchmarkLockSuiteParallel(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		BenchLockSuiteParallel(b, LockChecks.Release.Callable, lockIdle)
	})
}

// gateSubjectEmbedded is Gate under the name of the field that gateSubject
// embeds.
type gateSubjectEmbedded = Gate

// gateSubject is a value of Gate that only TestGateSuite builds.
type gateSubject struct{ gateSubjectEmbedded }

// TestGateSuite checks the suite of Gate and the two forms of its
// application.
func TestGateSuite(t *testing.T) {
	t.Parallel()

	t.Run("GateSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, GateSuite())
	})

	t.Run("RunGateSuite", func(t *testing.T) {
		t.Parallel()
		var built *gateSubject
		RunGateSuite(
			t,
			func() *gateSubject {
				built = &gateSubject{}
				return built
			},
			conform.With(GateCheck{
				ID:        "Gate/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *GateHandle) {
					got := h.Subject()
					assert.Equal(h, got, Gate(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*GateHandle) Gate { return &gateSubject{} },
			}),
		)
	})

	t.Run("RunGateSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *gateSubject
		RunGateSuiteFromHandle(
			t,
			func(*GateHandle) *gateSubject {
				built = &gateSubject{}
				return built
			},
			conform.With(GateCheck{
				ID:        "Gate/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *GateHandle) {
					got := h.Subject()
					assert.Equal(h, got, Gate(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*GateHandle) Gate { return &gateSubject{} },
			}),
		)
	})
}

// gateIdle returns the planted defect of a ceiling of Gate with a spend of
// nothing, a subject that answers each call at once.
func gateIdle(h *GateHandle) Gate {
	return GateSuite().Calls[0].Over(h, func() {})
}

// TestGateSuiteBudgets applies the allocation ceilings of the suite of
// Gate to gateIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestGateSuiteBudgets(t *testing.T) {
	RunGateBudgets(
		t,
		gateIdle,
	)
}

// FuzzGateSuitePass registers the fuzz target of Gate.Pass on gateIdle,
// with every check of the suite excused.
func FuzzGateSuitePass(f *testing.F) {
	FuzzGateSuite(
		f,
		GateChecks.Pass.Callable,
		gateIdle,
	)
}

// TestAwaitSuite checks the suite of Await and the two forms of its
// application.
func TestAwaitSuite(t *testing.T) {
	t.Parallel()

	t.Run("AwaitSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, AwaitSuite())
	})

	t.Run("Await", func(t *testing.T) {
		t.Parallel()
		RunAwaitSuite(t, Await)
	})

	t.Run("RunAwaitSuite", func(t *testing.T) {
		t.Parallel()
		RunAwaitSuite(
			t,
			Await,
			conform.With(AwaitCheck{
				ID:        "Await/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *AwaitHandle) {
					assert.Equal(h, h.Subject(), AwaitFunc(Await), h.Statement())
				},
				Defect: func(*AwaitHandle) AwaitFunc { return nil },
			}),
		)
	})

	t.Run("RunAwaitSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		RunAwaitSuiteFromHandle(
			t,
			func(*AwaitHandle) AwaitFunc {
				return Await
			},
			conform.With(AwaitCheck{
				ID:        "Await/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *AwaitHandle) {
					assert.Equal(h, h.Subject(), AwaitFunc(Await), h.Statement())
				},
				Defect: func(*AwaitHandle) AwaitFunc { return nil },
			}),
		)
	})
}

// FuzzAwait fuzzes the property checks of Await.
func FuzzAwait(f *testing.F) {
	FuzzAwaitSuite(f, Await)
}

// TestMutexConformsToLock applies the checks of the suite of Lock to
// Mutex.
func TestMutexConformsToLock(t *testing.T) {
	t.Parallel()

	RunLockSuite(t, NewMutex)
}

// BenchmarkMutexConformsToLock measures each callable of the suite of Lock
// on Mutex.
func BenchmarkMutexConformsToLock(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		BenchLockSuite(
			b,
			LockChecks.Release.Callable,
			func(*LockHandle) *Mutex { return NewMutex() },
		)
	})
	b.Run("Held", func(b *testing.B) {
		BenchLockSuite(b, LockChecks.Held.Callable, func(*LockHandle) *Mutex { return NewMutex() })
	})
}

// BenchmarkMutexConformsToLockParallel measures each callable of the suite
// of Lock that admits concurrent calls on Mutex, on GOMAXPROCS goroutines
// at once.
func BenchmarkMutexConformsToLockParallel(b *testing.B) {
	b.Run("Release", func(b *testing.B) {
		BenchLockSuiteParallel(
			b,
			LockChecks.Release.Callable,
			func(*LockHandle) *Mutex { return NewMutex() },
		)
	})
}

// FuzzMutexConformsToLockAcquire fuzzes the property checks of
// Lock.Acquire on Mutex.
func FuzzMutexConformsToLockAcquire(f *testing.F) {
	FuzzLockSuite(f, LockChecks.Acquire.Callable, func(*LockHandle) *Mutex { return NewMutex() })
}

// FuzzMutexConformsToLockRelease fuzzes the property checks of
// Lock.Release on Mutex.
func FuzzMutexConformsToLockRelease(f *testing.F) {
	FuzzLockSuite(f, LockChecks.Release.Callable, func(*LockHandle) *Mutex { return NewMutex() })
}

// FuzzMutexConformsToLockHeld fuzzes the property checks of Lock.Held on
// Mutex.
func FuzzMutexConformsToLockHeld(f *testing.F) {
	FuzzLockSuite(f, LockChecks.Held.Callable, func(*LockHandle) *Mutex { return NewMutex() })
}

// TestTurnstileConformsToGate applies the checks of the suite of Gate to
// Turnstile.
func TestTurnstileConformsToGate(t *testing.T) {
	t.Parallel()

	RunGateSuite(t, NewTurnstile)
}

// FuzzTurnstileConformsToGatePass fuzzes the property checks of Gate.Pass
// on Turnstile.
func FuzzTurnstileConformsToGatePass(f *testing.F) {
	FuzzGateSuite(
		f,
		GateChecks.Pass.Callable,
		func(*GateHandle) *Turnstile { return NewTurnstile() },
	)
}
