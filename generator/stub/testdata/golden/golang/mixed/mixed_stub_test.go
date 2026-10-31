package mixed_test

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/mixed"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
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

// sessionStubPrepareSubject returns the subject of [stub.Behaviour] for
// Prepare. The subject calls a new double that reports to tb.
func sessionStubPrepareSubject(
	tb assert.TB,
) stub.Subject[mixed.SessionPrepareCall, mixed.SessionPrepareReturn] {
	tb.Helper()
	s := mixed.NewSessionStub(tb)
	return stub.Subject[mixed.SessionPrepareCall, mixed.SessionPrepareReturn]{
		Stub: s.OnPrepare.MethodStub,
		Call: func() {
			_ = sessionStubPrepare(tb, s)
		},
		Result: func() mixed.SessionPrepareReturn {
			got0 := sessionStubPrepare(tb, s)
			return mixed.SessionPrepareReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnPrepare.Func(func(_ context.Context) error {
				mark()
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

// sessionStubCommitSubject returns the subject of [stub.Behaviour] for
// Commit. The subject calls a new double that reports to tb.
func sessionStubCommitSubject(
	tb assert.TB,
) stub.Subject[mixed.SessionCommitCall, mixed.SessionCommitReturn] {
	tb.Helper()
	s := mixed.NewSessionStub(tb)
	return stub.Subject[mixed.SessionCommitCall, mixed.SessionCommitReturn]{
		Stub: s.OnCommit.MethodStub,
		Call: func() {
			_ = sessionStubCommit(tb, s)
		},
		Result: func() mixed.SessionCommitReturn {
			got0 := sessionStubCommit(tb, s)
			return mixed.SessionCommitReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnCommit.Func(func(_ context.Context) error {
				mark()
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

// sessionStubOldSubject returns the subject of [stub.Behaviour] for Old.
// The subject calls a new double that reports to tb.
func sessionStubOldSubject(
	tb assert.TB,
) stub.Subject[mixed.SessionOldCall, mixed.SessionOldReturn] {
	tb.Helper()
	s := mixed.NewSessionStub(tb)
	return stub.Subject[mixed.SessionOldCall, mixed.SessionOldReturn]{
		Stub: s.OnOld.MethodStub,
		Call: func() {
			_, _ = sessionStubOld(tb, s)
		},
		Result: func() mixed.SessionOldReturn {
			got0, got1 := sessionStubOld(tb, s)
			return mixed.SessionOldReturn{Result: got0, Err: got1}
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
		Reset: s.ResetCalls,
	}
}

// TestSessionStub checks the methods, the constructors and the options of
// the double SessionStub.
func TestSessionStub(t *testing.T) {
	t.Parallel()

	t.Run("Prepare", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Prepare", sessionStubPrepareSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			_ = sessionStubPrepare(t, s)
			got := s.OnPrepare.AssertCalledOnce(t, "Prepare records the call")
			a0 := sessionStubPrepareArgs(t)
			assert.Equal(t, got.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			var seen []mixed.SessionPrepareCall
			s.OnPrepare.OnRecord(func(c mixed.SessionPrepareCall) { seen = append(seen, c) })
			_ = sessionStubPrepare(t, s)
			_ = sessionStubPrepare(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Prepare")
		})
	})

	t.Run("Commit", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Commit", sessionStubCommitSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			_ = sessionStubCommit(t, s)
			got := s.OnCommit.AssertCalledOnce(t, "Commit records the call")
			a0 := sessionStubCommitArgs(t)
			assert.Equal(t, got.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			var seen []mixed.SessionCommitCall
			s.OnCommit.OnRecord(func(c mixed.SessionCommitCall) { seen = append(seen, c) })
			_ = sessionStubCommit(t, s)
			_ = sessionStubCommit(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Commit")
		})
	})

	t.Run("Old", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Old", sessionStubOldSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			_, _ = sessionStubOld(t, s)
			got := s.OnOld.AssertCalledOnce(t, "Old records the call")
			a0 := sessionStubOldArgs(t)
			assert.Equal(t, got.Key, a0, "the call records Key")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			var seen []mixed.SessionOldCall
			s.OnOld.OnRecord(func(c mixed.SessionOldCall) { seen = append(seen, c) })
			_, _ = sessionStubOld(t, s)
			_, _ = sessionStubOld(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Old")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) sessionStubInstance {
			return sessionStubInstanceOf(tb, mixed.NewSessionStub(tb))
		})

		t.Run("clears the recorded calls of Prepare", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			_ = sessionStubPrepare(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnPrepare.CallCount(), 0, "ResetCalls clears the calls of Prepare")
		})

		t.Run("keeps the results of Returns of Prepare", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			want0 := sessionStubPreparePins()
			s.OnPrepare.Returns(want0)
			s.ResetCalls()
			got0 := sessionStubPrepare(t, s)
			assert.ErrorIs(t, got0, want0, "ResetCalls keeps the error of Prepare")
		})

		t.Run("clears the recorded calls of Commit", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			_ = sessionStubCommit(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnCommit.CallCount(), 0, "ResetCalls clears the calls of Commit")
		})

		t.Run("keeps the results of Returns of Commit", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			want0 := sessionStubCommitPins()
			s.OnCommit.Returns(want0)
			s.ResetCalls()
			got0 := sessionStubCommit(t, s)
			assert.ErrorIs(t, got0, want0, "ResetCalls keeps the error of Commit")
		})

		t.Run("clears the recorded calls of Old", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			_, _ = sessionStubOld(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnOld.CallCount(), 0, "ResetCalls clears the calls of Old")
		})

		t.Run("keeps the results of Returns of Old", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			want0, want1 := sessionStubOldPins()
			s.OnOld.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := sessionStubOld(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result of Old")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Old")
		})
	})

	t.Run("NewSessionStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) sessionStubInstance {
				return sessionStubInstanceOf(tb, mixed.NewSessionStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewSessionStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Prepare with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Prepare-failing")
			s := mixed.NewSessionStubFailing(t, want)
			r0 := sessionStubPrepare(t, s)
			assert.ErrorIs(t, r0, want, "Prepare fails with the error of NewSessionStubFailing")
		})

		t.Run("fails Commit with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Commit-failing")
			s := mixed.NewSessionStubFailing(t, want)
			r0 := sessionStubCommit(t, s)
			assert.ErrorIs(t, r0, want, "Commit fails with the error of NewSessionStubFailing")
		})

		t.Run("fails Old with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Old-failing")
			s := mixed.NewSessionStubFailing(t, want)
			_, r1 := sessionStubOld(t, s)
			assert.ErrorIs(t, r1, want, "Old fails with the error of NewSessionStubFailing")
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

// TestSessionPrepareStub checks the methods of SessionPrepareStub, which
// configures Prepare.
func TestSessionPrepareStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Prepare", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			want0 := sessionStubPreparePins()
			s.OnPrepare.Returns(want0)
			got0 := sessionStubPrepare(t, s)
			assert.ErrorIs(t, got0, want0, "Prepare returns the error that Returns sets")
		})

		t.Run("configures Prepare for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := mixed.NewSessionStub(r)
			s.OnPrepare.Strict()
			var want0 error
			s.OnPrepare.Returns(want0)
			_ = sessionStubPrepare(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Prepare")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Prepare to the function", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			var got0 context.Context
			s.OnPrepare.Func(func(p0 context.Context) error {
				got0 = p0
				var z0 error
				return z0
			})
			_ = sessionStubPrepare(t, s)
			a0 := sessionStubPrepareArgs(t)
			assert.Equal(t, got0, a0, "the function receives Ctx", assert.ByIdentity())
		})
	})
}

// TestSessionCommitStub checks the methods of SessionCommitStub, which
// configures Commit.
func TestSessionCommitStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Commit", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			want0 := sessionStubCommitPins()
			s.OnCommit.Returns(want0)
			got0 := sessionStubCommit(t, s)
			assert.ErrorIs(t, got0, want0, "Commit returns the error that Returns sets")
		})

		t.Run("configures Commit for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := mixed.NewSessionStub(r)
			s.OnCommit.Strict()
			var want0 error
			s.OnCommit.Returns(want0)
			_ = sessionStubCommit(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Commit")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Commit to the function", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			var got0 context.Context
			s.OnCommit.Func(func(p0 context.Context) error {
				got0 = p0
				var z0 error
				return z0
			})
			_ = sessionStubCommit(t, s)
			a0 := sessionStubCommitArgs(t)
			assert.Equal(t, got0, a0, "the function receives Ctx", assert.ByIdentity())
		})
	})
}

// TestSessionOldStub checks the methods of SessionOldStub, which
// configures Old.
func TestSessionOldStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Old", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			want0, want1 := sessionStubOldPins()
			s.OnOld.Returns(want0, want1)
			got0, got1 := sessionStubOld(t, s)
			expect.Equal(t, got0, want0, "Old returns the Result that Returns sets")
			expect.ErrorIs(t, got1, want1, "Old returns the error that Returns sets")
		})

		t.Run("configures Old for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := mixed.NewSessionStub(r)
			s.OnOld.Strict()
			var want0 string
			var want1 error
			s.OnOld.Returns(want0, want1)
			_, _ = sessionStubOld(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Old")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Old to the function", func(t *testing.T) {
			t.Parallel()
			s := mixed.NewSessionStub(t)
			var got0 string
			s.OnOld.Func(func(p0 string) (string, error) {
				got0 = p0
				var z0 string
				var z1 error
				return z0, z1
			})
			_, _ = sessionStubOld(t, s)
			a0 := sessionStubOldArgs(t)
			assert.Equal(t, got0, a0, "the function receives Key")
		})
	})
}
