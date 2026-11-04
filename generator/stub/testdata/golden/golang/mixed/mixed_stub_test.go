package mixed_test

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/mixed"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// sessionStubPrepareArgs returns the samples of the arguments of Prepare.
func sessionStubPrepareArgs(tb assert.TB) context.Context {
	tb.Helper()
	a0 := assert.Context(tb)
	return a0
}

// sessionStubPrepare calls Prepare of s with the samples of its arguments.
func sessionStubPrepare(tb assert.TB, s *mixed.SessionStub) error {
	tb.Helper()
	a0 := sessionStubPrepareArgs(tb)
	return s.Prepare(a0)
}

// sessionStubPreparePins returns the samples of the results of Prepare,
// with a new error as its error result.
func sessionStubPreparePins() error {
	want0 := errors.New("Prepare-returns")
	return want0
}

// sessionStubCommitArgs returns the samples of the arguments of Commit.
func sessionStubCommitArgs(tb assert.TB) context.Context {
	tb.Helper()
	a0 := assert.Context(tb)
	return a0
}

// sessionStubCommit calls Commit of s with the samples of its arguments.
func sessionStubCommit(tb assert.TB, s *mixed.SessionStub) error {
	tb.Helper()
	a0 := sessionStubCommitArgs(tb)
	return s.Commit(a0)
}

// sessionStubCommitPins returns the samples of the results of Commit, with
// a new error as its error result.
func sessionStubCommitPins() error {
	want0 := errors.New("Commit-returns")
	return want0
}

// sessionStubOldArgs returns the samples of the arguments of Old.
func sessionStubOldArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Old-a"
	return a0
}

// sessionStubOld calls Old of s with the samples of its arguments.
func sessionStubOld(tb assert.TB, s *mixed.SessionStub) (string, error) {
	tb.Helper()
	a0 := sessionStubOldArgs(tb)
	return s.Old(a0)
}

// sessionStubOldPins returns the samples of the results of Old, with a new
// error as its error result.
func sessionStubOldPins() (string, error) {
	var want0 string
	want1 := errors.New("Old-returns")
	want0 = "Old-a"
	return want0, want1
}

// sessionStubPrepareSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Prepare. The subject calls a new double that
// reports to tb.
func sessionStubPrepareSubject(
	tb assert.TB,
) stub.Subject[mixed.SessionPrepareCall, mixed.SessionPrepareReturn] {
	tb.Helper()
	s := mixed.NewSessionStub(tb)
	return stub.Subject[mixed.SessionPrepareCall, mixed.SessionPrepareReturn]{
		Stub:  s.OnPrepare.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = sessionStubPrepare(tb, s)
		},
		Arguments: func(tb assert.TB, c mixed.SessionPrepareCall) {
			tb.Helper()
			a0 := sessionStubPrepareArgs(tb)
			assert.Equal(tb, c.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		},
		Result: func() mixed.SessionPrepareReturn {
			got0 := sessionStubPrepare(tb, s)
			return mixed.SessionPrepareReturn{Err: got0}
		},
		Returns: func() mixed.SessionPrepareReturn {
			want0 := sessionStubPreparePins()
			s.OnPrepare.Returns(want0)
			return mixed.SessionPrepareReturn{Err: want0}
		},
		Override: func(seen func(mixed.SessionPrepareCall)) {
			s.OnPrepare.Func(func(p0 context.Context) error {
				seen(mixed.SessionPrepareCall{Ctx: p0})
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			r0 := sessionStubPrepare(tb, s)
			return r0
		},
	}
}

// sessionStubCommitSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Commit. The subject calls a new double that
// reports to tb.
func sessionStubCommitSubject(
	tb assert.TB,
) stub.Subject[mixed.SessionCommitCall, mixed.SessionCommitReturn] {
	tb.Helper()
	s := mixed.NewSessionStub(tb)
	return stub.Subject[mixed.SessionCommitCall, mixed.SessionCommitReturn]{
		Stub:  s.OnCommit.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = sessionStubCommit(tb, s)
		},
		Arguments: func(tb assert.TB, c mixed.SessionCommitCall) {
			tb.Helper()
			a0 := sessionStubCommitArgs(tb)
			assert.Equal(tb, c.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		},
		Result: func() mixed.SessionCommitReturn {
			got0 := sessionStubCommit(tb, s)
			return mixed.SessionCommitReturn{Err: got0}
		},
		Returns: func() mixed.SessionCommitReturn {
			want0 := sessionStubCommitPins()
			s.OnCommit.Returns(want0)
			return mixed.SessionCommitReturn{Err: want0}
		},
		Override: func(seen func(mixed.SessionCommitCall)) {
			s.OnCommit.Func(func(p0 context.Context) error {
				seen(mixed.SessionCommitCall{Ctx: p0})
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			r0 := sessionStubCommit(tb, s)
			return r0
		},
	}
}

// sessionStubOldSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Old. The subject calls a new double that
// reports to tb.
func sessionStubOldSubject(
	tb assert.TB,
) stub.Subject[mixed.SessionOldCall, mixed.SessionOldReturn] {
	tb.Helper()
	s := mixed.NewSessionStub(tb)
	return stub.Subject[mixed.SessionOldCall, mixed.SessionOldReturn]{
		Stub:  s.OnOld.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = sessionStubOld(tb, s)
		},
		Arguments: func(tb assert.TB, c mixed.SessionOldCall) {
			tb.Helper()
			a0 := sessionStubOldArgs(tb)
			assert.Equal(tb, c.Key, a0, "the call records Key")
		},
		Result: func() mixed.SessionOldReturn {
			got0, got1 := sessionStubOld(tb, s)
			return mixed.SessionOldReturn{Result: got0, Err: got1}
		},
		Returns: func() mixed.SessionOldReturn {
			want0, want1 := sessionStubOldPins()
			s.OnOld.Returns(want0, want1)
			return mixed.SessionOldReturn{Result: want0, Err: want1}
		},
		Override: func(seen func(mixed.SessionOldCall)) {
			s.OnOld.Func(func(p0 string) (string, error) {
				seen(mixed.SessionOldCall{Key: p0})
				var z0 string
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := sessionStubOld(tb, s)
			return r1
		},
	}
}

// sessionStubInstance is the instance of a mixed.SessionStub that the
// checks of the constructors and the options drive.
type sessionStubInstance = stub.Instance[mixed.SessionPrepareCall]

// sessionStubInstanceOf returns the instance of s. The options apply to
// the whole double, so the checks call Prepare alone.
func sessionStubInstanceOf(tb assert.TB, s *mixed.SessionStub) sessionStubInstance {
	tb.Helper()
	return sessionStubInstance{
		Stub: s.OnPrepare.MethodStub,
		Call: func() {
			_ = sessionStubPrepare(tb, s)
		},
	}
}

// TestSessionStub checks the methods, the constructors and the options of
// the double SessionStub.
func TestSessionStub(t *testing.T) {
	t.Parallel()

	t.Run("Prepare", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Prepare", sessionStubPrepareSubject)
	})

	t.Run("Commit", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Commit", sessionStubCommitSubject)
	})

	t.Run("Old", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Old", sessionStubOldSubject)
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("Prepare", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Prepare", sessionStubPrepareSubject)
		})

		t.Run("Commit", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Commit", sessionStubCommitSubject)
		})

		t.Run("Old", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Old", sessionStubOldSubject)
		})
	})

	t.Run("SessionStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) sessionStubInstance {
			return sessionStubInstanceOf(tb, mixed.NewSessionStub(tb, mixed.SessionStubStrict()))
		})

		t.Run("refuses Commit before Prepare", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			fn := func(_ context.Context) error {
				var z0 error
				return z0
			}
			s := mixed.NewSessionStub(r, mixed.SessionStubStrict(), mixed.WithSessionCommit(fn))
			_ = sessionStubCommit(t, s)
			assert.True(t, r.Failed(), "a strict double refuses Commit before Prepare")
		})
	})

	t.Run("SessionStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Prepare to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := mixed.NewSessionStub(t)
			s := mixed.NewSessionStub(t, mixed.SessionStubDelegateTo(inner))
			_ = sessionStubPrepare(t, s)
			inner.OnPrepare.AssertCalledOnce(t, "Prepare calls the implementation")
		})

		t.Run("fails Prepare with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Prepare-delegate")
			inner := mixed.NewSessionStub(t)
			inner.OnPrepare.Faults(want, 1)
			s := mixed.NewSessionStub(t, mixed.SessionStubDelegateTo(inner))
			r0 := sessionStubPrepare(t, s)
			assert.ErrorIs(t, r0, want, "Prepare fails with the error of the implementation")
		})

		t.Run("passes Commit to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := mixed.NewSessionStub(t)
			s := mixed.NewSessionStub(t, mixed.SessionStubDelegateTo(inner))
			_ = sessionStubCommit(t, s)
			inner.OnCommit.AssertCalledOnce(t, "Commit calls the implementation")
		})

		t.Run("fails Commit with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Commit-delegate")
			inner := mixed.NewSessionStub(t)
			inner.OnCommit.Faults(want, 1)
			s := mixed.NewSessionStub(t, mixed.SessionStubDelegateTo(inner))
			r0 := sessionStubCommit(t, s)
			assert.ErrorIs(t, r0, want, "Commit fails with the error of the implementation")
		})

		t.Run("passes Old to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := mixed.NewSessionStub(t)
			s := mixed.NewSessionStub(t, mixed.SessionStubDelegateTo(inner))
			_, _ = sessionStubOld(t, s)
			inner.OnOld.AssertCalledOnce(t, "Old calls the implementation")
		})

		t.Run("fails Old with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Old-delegate")
			inner := mixed.NewSessionStub(t)
			inner.OnOld.Faults(want, 1)
			s := mixed.NewSessionStub(t, mixed.SessionStubDelegateTo(inner))
			_, r1 := sessionStubOld(t, s)
			assert.ErrorIs(t, r1, want, "Old fails with the error of the implementation")
		})
	})

	t.Run("SessionStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) sessionStubInstance {
			return sessionStubInstanceOf(
				tb,
				mixed.NewSessionStub(tb, mixed.SessionStubWithClock(clk)),
			)
		})
	})

	t.Run("SessionStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) sessionStubInstance {
			return sessionStubInstanceOf(
				tb,
				mixed.NewSessionStub(tb, mixed.SessionStubWithRandSource(src)),
			)
		})
	})

	t.Run("SessionStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) sessionStubInstance {
			return sessionStubInstanceOf(tb, mixed.NewSessionStub(tb, mixed.SessionStubBenchMode()))
		})
	})

	t.Run("SessionStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) sessionStubInstance {
				return sessionStubInstanceOf(
					tb,
					mixed.NewSessionStub(
						tb,
						mixed.SessionStubWithClock(clk),
						mixed.SessionStubWithLatency(d),
					),
				)
			},
		)
	})

	t.Run("SessionStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Prepare with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Prepare-failing")
			s := mixed.NewSessionStub(t, mixed.SessionStubFailing(want))
			r0 := sessionStubPrepare(t, s)
			assert.ErrorIs(t, r0, want, "Prepare fails with the error of SessionStubFailing")
		})

		t.Run("fails Commit with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Commit-failing")
			s := mixed.NewSessionStub(t, mixed.SessionStubFailing(want))
			r0 := sessionStubCommit(t, s)
			assert.ErrorIs(t, r0, want, "Commit fails with the error of SessionStubFailing")
		})

		t.Run("fails Old with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Old-failing")
			s := mixed.NewSessionStub(t, mixed.SessionStubFailing(want))
			_, r1 := sessionStubOld(t, s)
			assert.ErrorIs(t, r1, want, "Old fails with the error of SessionStubFailing")
		})
	})

	t.Run("WithSessionPrepare", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Prepare", func(t *testing.T) {
			t.Parallel()
			called := false
			s := mixed.NewSessionStub(t, mixed.WithSessionPrepare(func(_ context.Context) error {
				called = true
				var z0 error
				return z0
			}))
			_ = sessionStubPrepare(t, s)
			assert.True(t, called, "WithSessionPrepare sets the function of Prepare")
		})
	})

	t.Run("WithSessionCommit", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Commit", func(t *testing.T) {
			t.Parallel()
			called := false
			s := mixed.NewSessionStub(t, mixed.WithSessionCommit(func(_ context.Context) error {
				called = true
				var z0 error
				return z0
			}))
			_ = sessionStubCommit(t, s)
			assert.True(t, called, "WithSessionCommit sets the function of Commit")
		})
	})

	t.Run("WithSessionOld", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Old", func(t *testing.T) {
			t.Parallel()
			called := false
			s := mixed.NewSessionStub(t, mixed.WithSessionOld(func(_ string) (string, error) {
				called = true
				var z0 string
				var z1 error
				return z0, z1
			}))
			_, _ = sessionStubOld(t, s)
			assert.True(t, called, "WithSessionOld sets the function of Old")
		})
	})
}
