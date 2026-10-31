package hidden_test

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/hidden"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/stub"
)

// hostsStubFindArgs returns the samples of the arguments of Find.
func hostsStubFindArgs(tb assert.TB) (context.Context, string) {
	tb.Helper()
	a0 := assert.Context(tb)
	var a1 string
	a1 = "Find-a"
	return a0, a1
}

// hostsStubFind calls Find of s with the samples of its arguments.
func hostsStubFind(tb assert.TB, s *hidden.HostsStub) (hidden.Host, error) {
	tb.Helper()
	a0, a1 := hostsStubFindArgs(tb)
	return s.Find(a0, a1)
}

// hostsStubFindPins returns the samples of the results of Find, with a new
// error as its error result.
func hostsStubFindPins() (hidden.Host, error) {
	var want0 hidden.Host
	want1 := errors.New("Find-returns")
	return want0, want1
}

// hostsStubFindSubject returns the subject of [stub.Behaviour] for Find.
// The subject calls a new double that reports to tb.
func hostsStubFindSubject(tb assert.TB) stub.Subject[hidden.HostsFindCall, hidden.HostsFindReturn] {
	tb.Helper()
	s := hidden.NewHostsStub(tb)
	return stub.Subject[hidden.HostsFindCall, hidden.HostsFindReturn]{
		Stub: s.OnFind.MethodStub,
		Call: func() {
			_, _ = hostsStubFind(tb, s)
		},
		Result: func() hidden.HostsFindReturn {
			got0, got1 := hostsStubFind(tb, s)
			return hidden.HostsFindReturn{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnFind.Func(func(_ context.Context, _ string) (hidden.Host, error) {
				mark()
				var z0 hidden.Host
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := hostsStubFind(tb, s)
			return r1
		},
	}
}

// hostsStubInstance is the instance of a hidden.HostsStub that the checks
// of the constructors and the options drive.
type hostsStubInstance = stub.Instance[hidden.HostsFindCall]

// hostsStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Find alone.
func hostsStubInstanceOf(tb assert.TB, s *hidden.HostsStub) hostsStubInstance {
	tb.Helper()
	return hostsStubInstance{
		Stub: s.OnFind.MethodStub,
		Call: func() {
			_, _ = hostsStubFind(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestHostsStub checks the methods, the constructors and the options of
// the double HostsStub.
func TestHostsStub(t *testing.T) {
	t.Parallel()

	t.Run("Find", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Find", hostsStubFindSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := hidden.NewHostsStub(t)
			_, _ = hostsStubFind(t, s)
			got := s.OnFind.AssertCalledOnce(t, "Find records the call")
			a0, a1 := hostsStubFindArgs(t)
			expect.Equal(t, got.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(t, got.Name, a1, "the call records Name")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := hidden.NewHostsStub(t)
			var seen []hidden.HostsFindCall
			s.OnFind.OnRecord(func(c hidden.HostsFindCall) { seen = append(seen, c) })
			_, _ = hostsStubFind(t, s)
			_, _ = hostsStubFind(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Find")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) hostsStubInstance {
			return hostsStubInstanceOf(tb, hidden.NewHostsStub(tb))
		})

		t.Run("clears the recorded calls of Find", func(t *testing.T) {
			t.Parallel()
			s := hidden.NewHostsStub(t)
			_, _ = hostsStubFind(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnFind.CallCount(), 0, "ResetCalls clears the calls of Find")
		})

		t.Run("keeps the results of Returns of Find", func(t *testing.T) {
			t.Parallel()
			s := hidden.NewHostsStub(t)
			want0, want1 := hostsStubFindPins()
			s.OnFind.Returns(want0, want1)
			s.ResetCalls()
			_, got1 := hostsStubFind(t, s)
			assert.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Find")
		})
	})

	t.Run("NewHostsStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) hostsStubInstance {
				return hostsStubInstanceOf(tb, hidden.NewHostsStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewHostsStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Find with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-failing")
			s := hidden.NewHostsStubFailing(t, want)
			_, r1 := hostsStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of NewHostsStubFailing")
		})
	})

	t.Run("HostsStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) hostsStubInstance {
			return hostsStubInstanceOf(tb, hidden.NewHostsStub(tb, hidden.HostsStubStrict()))
		})
	})

	t.Run("HostsStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Find to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := hidden.NewHostsStub(t)
			s := hidden.NewHostsStub(t, hidden.HostsStubDelegateTo(inner))
			_, _ = hostsStubFind(t, s)
			inner.OnFind.AssertCalledOnce(t, "Find calls the implementation")
		})

		t.Run("fails Find with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-delegate")
			inner := hidden.NewHostsStub(t)
			inner.OnFind.Faults(want, 1)
			s := hidden.NewHostsStub(t, hidden.HostsStubDelegateTo(inner))
			_, r1 := hostsStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of the implementation")
		})
	})

	t.Run("HostsStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) hostsStubInstance {
			return hostsStubInstanceOf(tb, hidden.NewHostsStub(tb, hidden.HostsStubWithClock(clk)))
		})
	})

	t.Run("HostsStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) hostsStubInstance {
			return hostsStubInstanceOf(
				tb,
				hidden.NewHostsStub(tb, hidden.HostsStubWithRandSource(src)),
			)
		})
	})

	t.Run("HostsStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) hostsStubInstance {
			return hostsStubInstanceOf(tb, hidden.NewHostsStub(tb, hidden.HostsStubBenchMode()))
		})
	})

	t.Run("HostsStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) hostsStubInstance {
				return hostsStubInstanceOf(
					tb,
					hidden.NewHostsStub(
						tb,
						hidden.HostsStubWithClock(clk),
						hidden.HostsStubWithLatency(d),
					),
				)
			},
		)
	})

	t.Run("HostsStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Find with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-failing")
			s := hidden.NewHostsStub(t, hidden.HostsStubFailing(want))
			_, r1 := hostsStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of HostsStubFailing")
		})
	})

	t.Run("WithHostsFind", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Find", func(t *testing.T) {
			t.Parallel()
			called := false
			s := hidden.NewHostsStub(
				t,
				hidden.WithHostsFind(func(_ context.Context, _ string) (hidden.Host, error) {
					called = true
					var z0 hidden.Host
					var z1 error
					return z0, z1
				}),
			)
			_, _ = hostsStubFind(t, s)
			assert.True(t, called, "WithHostsFind sets the function of Find")
		})
	})
}

// TestHostsFindStub checks the methods of HostsFindStub, which configures
// Find.
func TestHostsFindStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Find", func(t *testing.T) {
			t.Parallel()
			s := hidden.NewHostsStub(t)
			want0, want1 := hostsStubFindPins()
			s.OnFind.Returns(want0, want1)
			_, got1 := hostsStubFind(t, s)
			assert.ErrorIs(t, got1, want1, "Find returns the error that Returns sets")
		})

		t.Run("configures Find for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := hidden.NewHostsStub(r)
			s.OnFind.Strict()
			var want0 hidden.Host
			var want1 error
			s.OnFind.Returns(want0, want1)
			_, _ = hostsStubFind(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Find")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Find to the function", func(t *testing.T) {
			t.Parallel()
			s := hidden.NewHostsStub(t)
			var got0 context.Context
			var got1 string
			s.OnFind.Func(func(p0 context.Context, p1 string) (hidden.Host, error) {
				got0, got1 = p0, p1
				var z0 hidden.Host
				var z1 error
				return z0, z1
			})
			_, _ = hostsStubFind(t, s)
			a0, a1 := hostsStubFindArgs(t)
			expect.Equal(t, got0, a0, "the function receives Ctx", expect.ByIdentity())
			expect.Equal(t, got1, a1, "the function receives Name")
		})
	})
}

// FinderStub is generic, and its type parameters have no witnesses that a
// test outside its package can refer to, so the file has no checks of it.
// A test in the package of a consumer instantiates it.
