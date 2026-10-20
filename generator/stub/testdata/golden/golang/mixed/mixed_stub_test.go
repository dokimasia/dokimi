package mixed

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/dokimi-go/stub"
)

// sessionStubPrepareSubject binds Prepare into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func sessionStubPrepareSubject(tb assert.TB) stub.Subject[SessionPrepareCall, SessionPrepareReturn] {
	tb.Helper()
	s := NewSessionStub(tb)
	return stub.Subject[SessionPrepareCall, SessionPrepareReturn]{
		Stub: s.OnPrepare.MethodStub,
		Call: func() {
			var a0 context.Context
			_ = s.Prepare(a0)
		},
		Result: func() SessionPrepareReturn {
			var a0 context.Context
			got0 := s.Prepare(a0)
			return SessionPrepareReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnPrepare.Func(func(_ context.Context) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			var a0 context.Context
			r0 := s.Prepare(a0)
			return r0
		},
	}
}

// TestSessionStubPrepare pins how Prepare answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestSessionStubPrepare(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Prepare", sessionStubPrepareSubject)

	// No check that Prepare answers what Returns pinned: no result of its
	// signature admits a value this run can tell apart from the zero, and an
	// unconfigured double answers the zero anyway.
	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewSessionStub(t)
		var a0 context.Context
		_ = s.Prepare(a0)
		got := s.OnPrepare.AssertCalledOnce(t, "Prepare must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		// The hook fires synchronously as each call lands, which is what a
		// concurrency test observes progress with.
		s := NewSessionStub(t)
		var seen []SessionPrepareCall
		s.OnPrepare.OnRecord(func(c SessionPrepareCall) { seen = append(seen, c) })
		var a0 context.Context
		_ = s.Prepare(a0)
		_ = s.Prepare(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Prepare call")
	})

	t.Run("wires WithSessionPrepare at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewSessionStub(t, WithSessionPrepare(func(_ context.Context) error {
			called = true
			var z0 error
			return z0
		}))
		var a0 context.Context
		_ = s.Prepare(a0)
		assert.True(t, called, "WithSessionPrepare must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		// The preset arms every method that can fail, so it is checked once
		// per method. Built through the convenience constructor, which puts
		// the wrapper and the option it applies under the same check.
		want := errors.New("Prepare-failing")
		s := NewSessionStubFailing(t, want)
		var a0 context.Context
		r0 := s.Prepare(a0)
		assert.ErrorIs(t, r0, want, "Prepare must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		// No result of Prepare admits a distinguishable value, so what Returns
		// pinned cannot be told from the zero. What a reset can still be held
		// to is that it clears the log without taking the double with it:
		// configure, call, reset, call again, and exactly one call is
		// recorded.
		s := NewSessionStub(t)
		var want0 error
		s.OnPrepare.Returns(want0)
		var a0 context.Context
		_ = s.Prepare(a0)
		s.ResetCalls()
		_ = s.Prepare(a0)
		s.OnPrepare.AssertCalledOnce(t, "Prepare must record the call that followed the reset")
	})
}

// sessionStubCommitSubject binds Commit into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func sessionStubCommitSubject(tb assert.TB) stub.Subject[SessionCommitCall, SessionCommitReturn] {
	tb.Helper()
	s := NewSessionStub(tb)
	return stub.Subject[SessionCommitCall, SessionCommitReturn]{
		Stub: s.OnCommit.MethodStub,
		Call: func() {
			var a0 context.Context
			_ = s.Commit(a0)
		},
		Result: func() SessionCommitReturn {
			var a0 context.Context
			got0 := s.Commit(a0)
			return SessionCommitReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnCommit.Func(func(_ context.Context) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			var a0 context.Context
			r0 := s.Commit(a0)
			return r0
		},
	}
}

// TestSessionStubCommit pins how Commit answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestSessionStubCommit(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Commit", sessionStubCommitSubject)

	// No check that Commit answers what Returns pinned: no result of its
	// signature admits a value this run can tell apart from the zero, and an
	// unconfigured double answers the zero anyway.
	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewSessionStub(t)
		var a0 context.Context
		_ = s.Commit(a0)
		got := s.OnCommit.AssertCalledOnce(t, "Commit must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		// The hook fires synchronously as each call lands, which is what a
		// concurrency test observes progress with.
		s := NewSessionStub(t)
		var seen []SessionCommitCall
		s.OnCommit.OnRecord(func(c SessionCommitCall) { seen = append(seen, c) })
		var a0 context.Context
		_ = s.Commit(a0)
		_ = s.Commit(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Commit call")
	})

	t.Run("wires WithSessionCommit at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewSessionStub(t, WithSessionCommit(func(_ context.Context) error {
			called = true
			var z0 error
			return z0
		}))
		var a0 context.Context
		_ = s.Commit(a0)
		assert.True(t, called, "WithSessionCommit must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		// The preset arms every method that can fail, so it is checked once
		// per method. Built through the convenience constructor, which puts
		// the wrapper and the option it applies under the same check.
		want := errors.New("Commit-failing")
		s := NewSessionStubFailing(t, want)
		var a0 context.Context
		r0 := s.Commit(a0)
		assert.ErrorIs(t, r0, want, "Commit must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		// No result of Commit admits a distinguishable value, so what Returns
		// pinned cannot be told from the zero. What a reset can still be held
		// to is that it clears the log without taking the double with it:
		// configure, call, reset, call again, and exactly one call is
		// recorded.
		s := NewSessionStub(t)
		var want0 error
		s.OnCommit.Returns(want0)
		var a0 context.Context
		_ = s.Commit(a0)
		s.ResetCalls()
		_ = s.Commit(a0)
		s.OnCommit.AssertCalledOnce(t, "Commit must record the call that followed the reset")
	})

	t.Run("refuses Commit before Prepare in strict mode", func(t *testing.T) {
		t.Parallel()
		r := assert.NewRecorder()
		s := NewSessionStub(r, SessionStubStrict())
		var a0 context.Context
		_ = s.Commit(a0)
		assert.True(t, r.Failed(), "Commit before Prepare must be refused")
	})
}

// sessionStubOldSubject binds Old into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func sessionStubOldSubject(tb assert.TB) stub.Subject[SessionOldCall, SessionOldReturn] {
	tb.Helper()
	s := NewSessionStub(tb)
	return stub.Subject[SessionOldCall, SessionOldReturn]{
		Stub: s.OnOld.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Old-a"
			_, _ = s.Old(a0)
		},
		Result: func() SessionOldReturn {
			var a0 string
			a0 = "Old-a"
			got0, got1 := s.Old(a0)
			return SessionOldReturn{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnOld.Func(func(_ string) (string, error) {
				mark()
				var z0 string
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			var a0 string
			a0 = "Old-a"
			_, r1 := s.Old(a0)
			return r1
		},
	}
}

// TestSessionStubOld pins how Old answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestSessionStubOld(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Old", sessionStubOldSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewSessionStub(t)
		// Distinguishable from the zero wherever a value can be written: a
		// double that ignored Returns answers the zero too.
		var want0 string
		var want1 error
		want0 = "Old-a"
		s.OnOld.Returns(want0, want1)
		var a0 string
		a0 = "Old-a"
		got0, _ := s.Old(a0)
		assert.Equal(t, got0, want0, "Old must answer with what Returns pinned")
		// Blanked above: no literal exists for result 1 (error), so the
		// only value to pin there is the zero, which a double that ignored
		// Returns answers too.
	})
	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewSessionStub(t)
		var a0 string
		a0 = "Old-a"
		_, _ = s.Old(a0)
		got := s.OnOld.AssertCalledOnce(t, "Old must record the call")
		assert.Equal(t, got.Key, a0, "the recorded call carries Key")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		// The hook fires synchronously as each call lands, which is what a
		// concurrency test observes progress with.
		s := NewSessionStub(t)
		var seen []SessionOldCall
		s.OnOld.OnRecord(func(c SessionOldCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Old-a"
		_, _ = s.Old(a0)
		_, _ = s.Old(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Old call")
	})

	t.Run("wires WithSessionOld at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewSessionStub(t, WithSessionOld(func(_ string) (string, error) {
			called = true
			var z0 string
			var z1 error
			return z0, z1
		}))
		var a0 string
		a0 = "Old-a"
		_, _ = s.Old(a0)
		assert.True(t, called, "WithSessionOld must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		// The preset arms every method that can fail, so it is checked once
		// per method. Built through the convenience constructor, which puts
		// the wrapper and the option it applies under the same check.
		want := errors.New("Old-failing")
		s := NewSessionStubFailing(t, want)
		var a0 string
		a0 = "Old-a"
		_, r1 := s.Old(a0)
		assert.ErrorIs(t, r1, want, "Old must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		// A reset clears what happened, not what was configured. Pinned to a
		// distinguishable value: a reset that lost the configuration answers
		// the zero, and so does one that kept a configured zero.
		s := NewSessionStub(t)
		var want0 string
		var want1 error
		want0 = "Old-a"
		s.OnOld.Returns(want0, want1)
		var a0 string
		a0 = "Old-a"
		_, _ = s.Old(a0)
		s.ResetCalls()
		got0, _ := s.Old(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// sessionStubDouble describes how to build a SessionStub under each option
// whose effect is the same whatever a method's signature.
//
// Prepare stands in for the double as a whole: these checks assert only
// that an option reached it.
func sessionStubDouble() stub.Double[SessionPrepareCall] {
	instance := func(s *SessionStub) stub.Instance[SessionPrepareCall] {
		return stub.Instance[SessionPrepareCall]{
			Stub: s.OnPrepare.MethodStub,
			Call: func() {
				var a0 context.Context
				_ = s.Prepare(a0)
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[SessionPrepareCall]{
		New: func(tb assert.TB) stub.Instance[SessionPrepareCall] {
			return instance(NewSessionStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[SessionPrepareCall] {
			return instance(NewSessionStub(tb, SessionStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[SessionPrepareCall] {
			// Through the convenience constructor, which puts the wrapper
			// and the pair it bundles under the same check.
			return instance(NewSessionStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[SessionPrepareCall] {
			return instance(NewSessionStub(tb, SessionStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[SessionPrepareCall] {
			return instance(NewSessionStub(tb, SessionStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[SessionPrepareCall] {
			return instance(NewSessionStub(tb, SessionStubStrict()))
		},
	}
}

// TestSessionStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestSessionStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, sessionStubDouble())
}

// TestSessionStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestSessionStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewSessionStub(t)
	s := NewSessionStub(t, SessionStubDelegateTo(inner))

	t.Run("forwards Prepare to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		_ = s.Prepare(a0)
		inner.OnPrepare.AssertCalledOnce(t, "Prepare must reach the wrapped implementation")
	})

	t.Run("surfaces what Prepare answered", func(t *testing.T) {
		want := errors.New("Prepare-delegate")
		inner.OnPrepare.FaultsFor(time.Hour, want)
		var a0 context.Context
		r0 := s.Prepare(a0)
		assert.ErrorIs(t, r0, want, "Prepare must surface the wrapped answer")
	})

	t.Run("forwards Commit to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		_ = s.Commit(a0)
		inner.OnCommit.AssertCalledOnce(t, "Commit must reach the wrapped implementation")
	})

	t.Run("surfaces what Commit answered", func(t *testing.T) {
		want := errors.New("Commit-delegate")
		inner.OnCommit.FaultsFor(time.Hour, want)
		var a0 context.Context
		r0 := s.Commit(a0)
		assert.ErrorIs(t, r0, want, "Commit must surface the wrapped answer")
	})

	t.Run("forwards Old to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Old-a"
		_, _ = s.Old(a0)
		inner.OnOld.AssertCalledOnce(t, "Old must reach the wrapped implementation")
	})

	t.Run("surfaces what Old answered", func(t *testing.T) {
		want := errors.New("Old-delegate")
		inner.OnOld.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "Old-a"
		_, r1 := s.Old(a0)
		assert.ErrorIs(t, r1, want, "Old must surface the wrapped answer")
	})
}
