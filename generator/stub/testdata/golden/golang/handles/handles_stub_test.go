package handles

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/other"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/dokimi-go/stub"
)

// archiveStubSizeSubject binds Size into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func archiveStubSizeSubject(tb assert.TB) stub.Subject[ArchiveSizeCall, ArchiveSizeReturn] {
	tb.Helper()
	s := NewArchiveStub(tb)
	return stub.Subject[ArchiveSizeCall, ArchiveSizeReturn]{
		Stub: s.OnSize.MethodStub,
		Call: func() {
			_ = s.Size()
		},
		Result: func() ArchiveSizeReturn {
			got0 := s.Size()
			return ArchiveSizeReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnSize.Func(func() int {
				mark()
				var z0 int
				return z0
			})
		},
	}
}

// TestArchiveStubSize pins how Size answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestArchiveStubSize(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Size", archiveStubSizeSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewArchiveStub(t)
		var want0 int
		want0 = 42
		s.OnSize.Returns(want0)
		got0 := s.Size()
		assert.Equal(t, got0, want0, "Size must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewArchiveStub(t)
		_ = s.Size()
		s.OnSize.AssertCalledOnce(t, "Size must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewArchiveStub(t)
		var seen []ArchiveSizeCall
		s.OnSize.OnRecord(func(c ArchiveSizeCall) { seen = append(seen, c) })
		_ = s.Size()
		_ = s.Size()
		assert.Length(t, seen, 2, "OnRecord must fire once per Size call")
	})

	t.Run("wires WithArchiveSize at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewArchiveStub(t, WithArchiveSize(func() int {
			called = true
			var z0 int
			return z0
		}))
		_ = s.Size()
		assert.True(t, called, "WithArchiveSize must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewArchiveStub(t)
		var want0 int
		want0 = 42
		s.OnSize.Returns(want0)
		_ = s.Size()
		s.ResetCalls()
		got0 := s.Size()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// archiveStubOpenLogSubject binds OpenLog into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func archiveStubOpenLogSubject(tb assert.TB) stub.Subject[ArchiveOpenLogCall, ArchiveOpenLogReturn] {
	tb.Helper()
	s := NewArchiveStub(tb)
	return stub.Subject[ArchiveOpenLogCall, ArchiveOpenLogReturn]{
		Stub: s.OnOpenLog.MethodStub,
		Call: func() {
			var a0 string
			a0 = "OpenLog-a"
			_, _ = s.OpenLog(a0)
		},
		Result: func() ArchiveOpenLogReturn {
			var a0 string
			a0 = "OpenLog-a"
			got0, got1 := s.OpenLog(a0)
			return ArchiveOpenLogReturn{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnOpenLog.Func(func(_ string) (Log, error) {
				mark()
				var z0 Log
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			var a0 string
			a0 = "OpenLog-a"
			_, r1 := s.OpenLog(a0)
			return r1
		},
	}
}

// TestArchiveStubOpenLog pins how OpenLog answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestArchiveStubOpenLog(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "OpenLog", archiveStubOpenLogSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewArchiveStub(t)
		var a0 string
		a0 = "OpenLog-a"
		_, _ = s.OpenLog(a0)
		got := s.OnOpenLog.AssertCalledOnce(t, "OpenLog must record the call")
		assert.Equal(t, got.Name, a0, "the recorded call carries Name")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewArchiveStub(t)
		var seen []ArchiveOpenLogCall
		s.OnOpenLog.OnRecord(func(c ArchiveOpenLogCall) { seen = append(seen, c) })
		var a0 string
		a0 = "OpenLog-a"
		_, _ = s.OpenLog(a0)
		_, _ = s.OpenLog(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per OpenLog call")
	})

	t.Run("wires WithArchiveOpenLog at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewArchiveStub(t, WithArchiveOpenLog(func(_ string) (Log, error) {
			called = true
			var z0 Log
			var z1 error
			return z0, z1
		}))
		var a0 string
		a0 = "OpenLog-a"
		_, _ = s.OpenLog(a0)
		assert.True(t, called, "WithArchiveOpenLog must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("OpenLog-failing")
		s := NewArchiveStubFailing(t, want)
		var a0 string
		a0 = "OpenLog-a"
		_, r1 := s.OpenLog(a0)
		assert.ErrorIs(t, r1, want, "OpenLog must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewArchiveStub(t)
		var want0 Log
		var want1 error
		s.OnOpenLog.Returns(want0, want1)
		var a0 string
		a0 = "OpenLog-a"
		_, _ = s.OpenLog(a0)
		s.ResetCalls()
		_, _ = s.OpenLog(a0)
		s.OnOpenLog.AssertCalledOnce(t, "OpenLog must record the call that followed the reset")
	})
}

// archiveStubDouble describes how to build a ArchiveStub under each option
// whose effect is the same whatever a method's signature.
//
// Size stands in for the double as a whole: these checks assert only that
// an option reached it.
func archiveStubDouble() stub.Double[ArchiveSizeCall] {
	instance := func(s *ArchiveStub) stub.Instance[ArchiveSizeCall] {
		return stub.Instance[ArchiveSizeCall]{
			Stub: s.OnSize.MethodStub,
			Call: func() {
				_ = s.Size()
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[ArchiveSizeCall]{
		New: func(tb assert.TB) stub.Instance[ArchiveSizeCall] {
			return instance(NewArchiveStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[ArchiveSizeCall] {
			return instance(NewArchiveStub(tb, ArchiveStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[ArchiveSizeCall] {
			return instance(NewArchiveStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[ArchiveSizeCall] {
			return instance(NewArchiveStub(tb, ArchiveStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[ArchiveSizeCall] {
			return instance(NewArchiveStub(tb, ArchiveStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[ArchiveSizeCall] {
			return instance(NewArchiveStub(tb, ArchiveStubStrict()))
		},
	}
}

// TestArchiveStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestArchiveStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, archiveStubDouble())
}

// TestArchiveStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestArchiveStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewArchiveStub(t)
	s := NewArchiveStub(t, ArchiveStubDelegateTo(inner))

	t.Run("forwards Size to the wrapped implementation", func(t *testing.T) {
		_ = s.Size()
		inner.OnSize.AssertCalledOnce(t, "Size must reach the wrapped implementation")
	})

	t.Run("forwards OpenLog to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "OpenLog-a"
		_, _ = s.OpenLog(a0)
		inner.OnOpenLog.AssertCalledOnce(t, "OpenLog must reach the wrapped implementation")
	})

	t.Run("surfaces what OpenLog answered", func(t *testing.T) {
		want := errors.New("OpenLog-delegate")
		inner.OnOpenLog.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "OpenLog-a"
		_, r1 := s.OpenLog(a0)
		assert.ErrorIs(t, r1, want, "OpenLog must surface the wrapped answer")
	})
}

// logStubAppendSubject binds Append into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func logStubAppendSubject(tb assert.TB) stub.Subject[LogAppendCall, LogAppendReturn] {
	tb.Helper()
	s := NewLogStub(tb)
	return stub.Subject[LogAppendCall, LogAppendReturn]{
		Stub: s.OnAppend.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Append-a"
			_ = s.Append(a0)
		},
		Result: func() LogAppendReturn {
			var a0 string
			a0 = "Append-a"
			got0 := s.Append(a0)
			return LogAppendReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnAppend.Func(func(_ string) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			var a0 string
			a0 = "Append-a"
			r0 := s.Append(a0)
			return r0
		},
	}
}

// TestLogStubAppend pins how Append answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestLogStubAppend(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Append", logStubAppendSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewLogStub(t)
		var a0 string
		a0 = "Append-a"
		_ = s.Append(a0)
		got := s.OnAppend.AssertCalledOnce(t, "Append must record the call")
		assert.Equal(t, got.Line, a0, "the recorded call carries Line")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewLogStub(t)
		var seen []LogAppendCall
		s.OnAppend.OnRecord(func(c LogAppendCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Append-a"
		_ = s.Append(a0)
		_ = s.Append(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Append call")
	})

	t.Run("wires WithLogAppend at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewLogStub(t, WithLogAppend(func(_ string) error {
			called = true
			var z0 error
			return z0
		}))
		var a0 string
		a0 = "Append-a"
		_ = s.Append(a0)
		assert.True(t, called, "WithLogAppend must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Append-failing")
		s := NewLogStubFailing(t, want)
		var a0 string
		a0 = "Append-a"
		r0 := s.Append(a0)
		assert.ErrorIs(t, r0, want, "Append must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewLogStub(t)
		var want0 error
		s.OnAppend.Returns(want0)
		var a0 string
		a0 = "Append-a"
		_ = s.Append(a0)
		s.ResetCalls()
		_ = s.Append(a0)
		s.OnAppend.AssertCalledOnce(t, "Append must record the call that followed the reset")
	})
}

// logStubDouble describes how to build a LogStub under each option whose
// effect is the same whatever a method's signature.
//
// Append stands in for the double as a whole: these checks assert only
// that an option reached it.
func logStubDouble() stub.Double[LogAppendCall] {
	instance := func(s *LogStub) stub.Instance[LogAppendCall] {
		return stub.Instance[LogAppendCall]{
			Stub: s.OnAppend.MethodStub,
			Call: func() {
				var a0 string
				a0 = "Append-a"
				_ = s.Append(a0)
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[LogAppendCall]{
		New: func(tb assert.TB) stub.Instance[LogAppendCall] {
			return instance(NewLogStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[LogAppendCall] {
			return instance(NewLogStub(tb, LogStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[LogAppendCall] {
			return instance(NewLogStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[LogAppendCall] {
			return instance(NewLogStub(tb, LogStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[LogAppendCall] {
			return instance(NewLogStub(tb, LogStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[LogAppendCall] {
			return instance(NewLogStub(tb, LogStubStrict()))
		},
	}
}

// TestLogStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestLogStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, logStubDouble())
}

// TestLogStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestLogStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewLogStub(t)
	s := NewLogStub(t, LogStubDelegateTo(inner))

	t.Run("forwards Append to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Append-a"
		_ = s.Append(a0)
		inner.OnAppend.AssertCalledOnce(t, "Append must reach the wrapped implementation")
	})

	t.Run("surfaces what Append answered", func(t *testing.T) {
		want := errors.New("Append-delegate")
		inner.OnAppend.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "Append-a"
		r0 := s.Append(a0)
		assert.ErrorIs(t, r0, want, "Append must surface the wrapped answer")
	})
}

// producerStubOpenSubject binds Open into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func producerStubOpenSubject(tb assert.TB) stub.Subject[ProducerOpenCall, ProducerOpenReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerOpenCall, ProducerOpenReturn]{
		Stub: s.OnOpen.MethodStub,
		Call: func() {
			var a0 context.Context
			_, _ = s.Open(a0)
		},
		Result: func() ProducerOpenReturn {
			var a0 context.Context
			got0, got1 := s.Open(a0)
			return ProducerOpenReturn{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnOpen.Func(func(_ context.Context) (Cursor, error) {
				mark()
				var z0 Cursor
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			var a0 context.Context
			_, r1 := s.Open(a0)
			return r1
		},
	}
}

// TestProducerStubOpen pins how Open answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestProducerStubOpen(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Open", producerStubOpenSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var a0 context.Context
		_, _ = s.Open(a0)
		got := s.OnOpen.AssertCalledOnce(t, "Open must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var seen []ProducerOpenCall
		s.OnOpen.OnRecord(func(c ProducerOpenCall) { seen = append(seen, c) })
		var a0 context.Context
		_, _ = s.Open(a0)
		_, _ = s.Open(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Open call")
	})

	t.Run("wires WithProducerOpen at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewProducerStub(t, WithProducerOpen(func(_ context.Context) (Cursor, error) {
			called = true
			var z0 Cursor
			var z1 error
			return z0, z1
		}))
		var a0 context.Context
		_, _ = s.Open(a0)
		assert.True(t, called, "WithProducerOpen must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Open-failing")
		s := NewProducerStubFailing(t, want)
		var a0 context.Context
		_, r1 := s.Open(a0)
		assert.ErrorIs(t, r1, want, "Open must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var want0 Cursor
		var want1 error
		s.OnOpen.Returns(want0, want1)
		var a0 context.Context
		_, _ = s.Open(a0)
		s.ResetCalls()
		_, _ = s.Open(a0)
		s.OnOpen.AssertCalledOnce(t, "Open must record the call that followed the reset")
	})
}

// producerStubResumeSubject binds Resume into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func producerStubResumeSubject(tb assert.TB) stub.Subject[ProducerResumeCall, ProducerResumeReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerResumeCall, ProducerResumeReturn]{
		Stub: s.OnResume.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Resume-a"
			_ = s.Resume(a0)
		},
		Result: func() ProducerResumeReturn {
			var a0 string
			a0 = "Resume-a"
			got0 := s.Resume(a0)
			return ProducerResumeReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnResume.Func(func(_ string) Cursor {
				mark()
				var z0 Cursor
				return z0
			})
		},
	}
}

// TestProducerStubResume pins how Resume answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestProducerStubResume(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Resume", producerStubResumeSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var a0 string
		a0 = "Resume-a"
		_ = s.Resume(a0)
		got := s.OnResume.AssertCalledOnce(t, "Resume must record the call")
		assert.Equal(t, got.Token, a0, "the recorded call carries Token")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var seen []ProducerResumeCall
		s.OnResume.OnRecord(func(c ProducerResumeCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Resume-a"
		_ = s.Resume(a0)
		_ = s.Resume(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Resume call")
	})

	t.Run("wires WithProducerResume at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewProducerStub(t, WithProducerResume(func(_ string) Cursor {
			called = true
			var z0 Cursor
			return z0
		}))
		var a0 string
		a0 = "Resume-a"
		_ = s.Resume(a0)
		assert.True(t, called, "WithProducerResume must install the override")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var want0 Cursor
		s.OnResume.Returns(want0)
		var a0 string
		a0 = "Resume-a"
		_ = s.Resume(a0)
		s.ResetCalls()
		_ = s.Resume(a0)
		s.OnResume.AssertCalledOnce(t, "Resume must record the call that followed the reset")
	})
}

// producerStubWatchSubject binds Watch into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func producerStubWatchSubject(tb assert.TB) stub.Subject[ProducerWatchCall, ProducerWatchReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerWatchCall, ProducerWatchReturn]{
		Stub: s.OnWatch.MethodStub,
		Call: func() {
			_ = s.Watch()
		},
		Result: func() ProducerWatchReturn {
			got0 := s.Watch()
			return ProducerWatchReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnWatch.Func(func() Subscription {
				mark()
				var z0 Subscription
				return z0
			})
		},
	}
}

// TestProducerStubWatch pins how Watch answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestProducerStubWatch(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Watch", producerStubWatchSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		_ = s.Watch()
		s.OnWatch.AssertCalledOnce(t, "Watch must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var seen []ProducerWatchCall
		s.OnWatch.OnRecord(func(c ProducerWatchCall) { seen = append(seen, c) })
		_ = s.Watch()
		_ = s.Watch()
		assert.Length(t, seen, 2, "OnRecord must fire once per Watch call")
	})

	t.Run("wires WithProducerWatch at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewProducerStub(t, WithProducerWatch(func() Subscription {
			called = true
			var z0 Subscription
			return z0
		}))
		_ = s.Watch()
		assert.True(t, called, "WithProducerWatch must install the override")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var want0 Subscription
		s.OnWatch.Returns(want0)
		_ = s.Watch()
		s.ResetCalls()
		_ = s.Watch()
		s.OnWatch.AssertCalledOnce(t, "Watch must record the call that followed the reset")
	})
}

// producerStubNothingSubject binds Nothing into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func producerStubNothingSubject(tb assert.TB) stub.Subject[ProducerNothingCall, ProducerNothingReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerNothingCall, ProducerNothingReturn]{
		Stub: s.OnNothing.MethodStub,
		Call: func() {
			_ = s.Nothing()
		},
		Result: func() ProducerNothingReturn {
			got0 := s.Nothing()
			return ProducerNothingReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnNothing.Func(func() Empty {
				mark()
				var z0 Empty
				return z0
			})
		},
	}
}

// TestProducerStubNothing pins how Nothing answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestProducerStubNothing(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Nothing", producerStubNothingSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		_ = s.Nothing()
		s.OnNothing.AssertCalledOnce(t, "Nothing must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var seen []ProducerNothingCall
		s.OnNothing.OnRecord(func(c ProducerNothingCall) { seen = append(seen, c) })
		_ = s.Nothing()
		_ = s.Nothing()
		assert.Length(t, seen, 2, "OnRecord must fire once per Nothing call")
	})

	t.Run("wires WithProducerNothing at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewProducerStub(t, WithProducerNothing(func() Empty {
			called = true
			var z0 Empty
			return z0
		}))
		_ = s.Nothing()
		assert.True(t, called, "WithProducerNothing must install the override")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var want0 Empty
		s.OnNothing.Returns(want0)
		_ = s.Nothing()
		s.ResetCalls()
		_ = s.Nothing()
		s.OnNothing.AssertCalledOnce(t, "Nothing must record the call that followed the reset")
	})
}

// producerStubElsewhereSubject binds Elsewhere into the shape
// [stub.Behaviour] drives: how to call it, what it answers with, and how
// to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func producerStubElsewhereSubject(tb assert.TB) stub.Subject[ProducerElsewhereCall, ProducerElsewhereReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerElsewhereCall, ProducerElsewhereReturn]{
		Stub: s.OnElsewhere.MethodStub,
		Call: func() {
			_ = s.Elsewhere()
		},
		Result: func() ProducerElsewhereReturn {
			got0 := s.Elsewhere()
			return ProducerElsewhereReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnElsewhere.Func(func() other.Thing {
				mark()
				var z0 other.Thing
				return z0
			})
		},
	}
}

// TestProducerStubElsewhere pins how Elsewhere answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestProducerStubElsewhere(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Elsewhere", producerStubElsewhereSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		_ = s.Elsewhere()
		s.OnElsewhere.AssertCalledOnce(t, "Elsewhere must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var seen []ProducerElsewhereCall
		s.OnElsewhere.OnRecord(func(c ProducerElsewhereCall) { seen = append(seen, c) })
		_ = s.Elsewhere()
		_ = s.Elsewhere()
		assert.Length(t, seen, 2, "OnRecord must fire once per Elsewhere call")
	})

	t.Run("wires WithProducerElsewhere at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewProducerStub(t, WithProducerElsewhere(func() other.Thing {
			called = true
			var z0 other.Thing
			return z0
		}))
		_ = s.Elsewhere()
		assert.True(t, called, "WithProducerElsewhere must install the override")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewProducerStub(t)
		var want0 other.Thing
		s.OnElsewhere.Returns(want0)
		_ = s.Elsewhere()
		s.ResetCalls()
		_ = s.Elsewhere()
		s.OnElsewhere.AssertCalledOnce(t, "Elsewhere must record the call that followed the reset")
	})
}

// producerStubDouble describes how to build a ProducerStub under each
// option whose effect is the same whatever a method's signature.
//
// Open stands in for the double as a whole: these checks assert only that
// an option reached it.
func producerStubDouble() stub.Double[ProducerOpenCall] {
	instance := func(s *ProducerStub) stub.Instance[ProducerOpenCall] {
		return stub.Instance[ProducerOpenCall]{
			Stub: s.OnOpen.MethodStub,
			Call: func() {
				var a0 context.Context
				_, _ = s.Open(a0)
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[ProducerOpenCall]{
		New: func(tb assert.TB) stub.Instance[ProducerOpenCall] {
			return instance(NewProducerStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[ProducerOpenCall] {
			return instance(NewProducerStub(tb, ProducerStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[ProducerOpenCall] {
			return instance(NewProducerStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[ProducerOpenCall] {
			return instance(NewProducerStub(tb, ProducerStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[ProducerOpenCall] {
			return instance(NewProducerStub(tb, ProducerStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[ProducerOpenCall] {
			return instance(NewProducerStub(tb, ProducerStubStrict()))
		},
	}
}

// TestProducerStubOptions pins the settings that apply to the double
// rather than to one method: bench mode, the injected clock and random
// source, and what a reset clears.
func TestProducerStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, producerStubDouble())
}

// TestProducerStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestProducerStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewProducerStub(t)
	s := NewProducerStub(t, ProducerStubDelegateTo(inner))

	t.Run("forwards Open to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		_, _ = s.Open(a0)
		inner.OnOpen.AssertCalledOnce(t, "Open must reach the wrapped implementation")
	})

	t.Run("surfaces what Open answered", func(t *testing.T) {
		want := errors.New("Open-delegate")
		inner.OnOpen.FaultsFor(time.Hour, want)
		var a0 context.Context
		_, r1 := s.Open(a0)
		assert.ErrorIs(t, r1, want, "Open must surface the wrapped answer")
	})

	t.Run("forwards Resume to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Resume-a"
		_ = s.Resume(a0)
		inner.OnResume.AssertCalledOnce(t, "Resume must reach the wrapped implementation")
	})

	t.Run("forwards Watch to the wrapped implementation", func(t *testing.T) {
		_ = s.Watch()
		inner.OnWatch.AssertCalledOnce(t, "Watch must reach the wrapped implementation")
	})

	t.Run("forwards Nothing to the wrapped implementation", func(t *testing.T) {
		_ = s.Nothing()
		inner.OnNothing.AssertCalledOnce(t, "Nothing must reach the wrapped implementation")
	})

	t.Run("forwards Elsewhere to the wrapped implementation", func(t *testing.T) {
		_ = s.Elsewhere()
		inner.OnElsewhere.AssertCalledOnce(t, "Elsewhere must reach the wrapped implementation")
	})
}

// cursorStubNextSubject binds Next into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func cursorStubNextSubject(tb assert.TB) stub.Subject[CursorNextCall, CursorNextReturn] {
	tb.Helper()
	s := NewCursorStub(tb)
	return stub.Subject[CursorNextCall, CursorNextReturn]{
		Stub: s.OnNext.MethodStub,
		Call: func() {
			var a0 context.Context
			_, _, _ = s.Next(a0)
		},
		Result: func() CursorNextReturn {
			var a0 context.Context
			got0, got1, got2 := s.Next(a0)
			return CursorNextReturn{Result0: got0, Result1: got1, Err: got2}
		},
		Override: func(mark func()) {
			s.OnNext.Func(func(_ context.Context) (string, bool, error) {
				mark()
				var z0 string
				var z1 bool
				var z2 error
				return z0, z1, z2
			})
		},
		Fails: func() error {
			var a0 context.Context
			_, _, r2 := s.Next(a0)
			return r2
		},
	}
}

// TestCursorStubNext pins how Next answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestCursorStubNext(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Next", cursorStubNextSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewCursorStub(t)
		var want0 string
		var want1 bool
		var want2 error
		want0 = "Next-a"
		want1 = true
		s.OnNext.Returns(want0, want1, want2)
		var a0 context.Context
		got0, got1, _ := s.Next(a0)
		assert.Equal(t, got0, want0, "Next must answer with what Returns pinned")
		assert.Equal(t, got1, want1, "Next must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewCursorStub(t)
		var a0 context.Context
		_, _, _ = s.Next(a0)
		got := s.OnNext.AssertCalledOnce(t, "Next must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewCursorStub(t)
		var seen []CursorNextCall
		s.OnNext.OnRecord(func(c CursorNextCall) { seen = append(seen, c) })
		var a0 context.Context
		_, _, _ = s.Next(a0)
		_, _, _ = s.Next(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Next call")
	})

	t.Run("wires WithCursorNext at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewCursorStub(t, WithCursorNext(func(_ context.Context) (string, bool, error) {
			called = true
			var z0 string
			var z1 bool
			var z2 error
			return z0, z1, z2
		}))
		var a0 context.Context
		_, _, _ = s.Next(a0)
		assert.True(t, called, "WithCursorNext must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Next-failing")
		s := NewCursorStubFailing(t, want)
		var a0 context.Context
		_, _, r2 := s.Next(a0)
		assert.ErrorIs(t, r2, want, "Next must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewCursorStub(t)
		var want0 string
		var want1 bool
		var want2 error
		want0 = "Next-a"
		want1 = true
		s.OnNext.Returns(want0, want1, want2)
		var a0 context.Context
		_, _, _ = s.Next(a0)
		s.ResetCalls()
		got0, got1, _ := s.Next(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
		assert.Equal(t, got1, want1, "a reset must keep what Returns pinned")
	})
}

// cursorStubCloseSubject binds Close into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func cursorStubCloseSubject(tb assert.TB) stub.Subject[CursorCloseCall, CursorCloseReturn] {
	tb.Helper()
	s := NewCursorStub(tb)
	return stub.Subject[CursorCloseCall, CursorCloseReturn]{
		Stub: s.OnClose.MethodStub,
		Call: func() {
			var a0 context.Context
			_ = s.Close(a0)
		},
		Result: func() CursorCloseReturn {
			var a0 context.Context
			got0 := s.Close(a0)
			return CursorCloseReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnClose.Func(func(_ context.Context) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			var a0 context.Context
			r0 := s.Close(a0)
			return r0
		},
	}
}

// TestCursorStubClose pins how Close answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestCursorStubClose(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Close", cursorStubCloseSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewCursorStub(t)
		var a0 context.Context
		_ = s.Close(a0)
		got := s.OnClose.AssertCalledOnce(t, "Close must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewCursorStub(t)
		var seen []CursorCloseCall
		s.OnClose.OnRecord(func(c CursorCloseCall) { seen = append(seen, c) })
		var a0 context.Context
		_ = s.Close(a0)
		_ = s.Close(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Close call")
	})

	t.Run("wires WithCursorClose at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewCursorStub(t, WithCursorClose(func(_ context.Context) error {
			called = true
			var z0 error
			return z0
		}))
		var a0 context.Context
		_ = s.Close(a0)
		assert.True(t, called, "WithCursorClose must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Close-failing")
		s := NewCursorStubFailing(t, want)
		var a0 context.Context
		r0 := s.Close(a0)
		assert.ErrorIs(t, r0, want, "Close must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewCursorStub(t)
		var want0 error
		s.OnClose.Returns(want0)
		var a0 context.Context
		_ = s.Close(a0)
		s.ResetCalls()
		_ = s.Close(a0)
		s.OnClose.AssertCalledOnce(t, "Close must record the call that followed the reset")
	})
}

// cursorStubDouble describes how to build a CursorStub under each option
// whose effect is the same whatever a method's signature.
//
// Next stands in for the double as a whole: these checks assert only that
// an option reached it.
func cursorStubDouble() stub.Double[CursorNextCall] {
	instance := func(s *CursorStub) stub.Instance[CursorNextCall] {
		return stub.Instance[CursorNextCall]{
			Stub: s.OnNext.MethodStub,
			Call: func() {
				var a0 context.Context
				_, _, _ = s.Next(a0)
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[CursorNextCall]{
		New: func(tb assert.TB) stub.Instance[CursorNextCall] {
			return instance(NewCursorStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[CursorNextCall] {
			return instance(NewCursorStub(tb, CursorStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[CursorNextCall] {
			return instance(NewCursorStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[CursorNextCall] {
			return instance(NewCursorStub(tb, CursorStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[CursorNextCall] {
			return instance(NewCursorStub(tb, CursorStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[CursorNextCall] {
			return instance(NewCursorStub(tb, CursorStubStrict()))
		},
	}
}

// TestCursorStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestCursorStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, cursorStubDouble())
}

// TestCursorStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestCursorStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewCursorStub(t)
	s := NewCursorStub(t, CursorStubDelegateTo(inner))

	t.Run("forwards Next to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		_, _, _ = s.Next(a0)
		inner.OnNext.AssertCalledOnce(t, "Next must reach the wrapped implementation")
	})

	t.Run("surfaces what Next answered", func(t *testing.T) {
		want := errors.New("Next-delegate")
		inner.OnNext.FaultsFor(time.Hour, want)
		var a0 context.Context
		_, _, r2 := s.Next(a0)
		assert.ErrorIs(t, r2, want, "Next must surface the wrapped answer")
	})

	t.Run("forwards Close to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		_ = s.Close(a0)
		inner.OnClose.AssertCalledOnce(t, "Close must reach the wrapped implementation")
	})

	t.Run("surfaces what Close answered", func(t *testing.T) {
		want := errors.New("Close-delegate")
		inner.OnClose.FaultsFor(time.Hour, want)
		var a0 context.Context
		r0 := s.Close(a0)
		assert.ErrorIs(t, r0, want, "Close must surface the wrapped answer")
	})
}

// readerStubReadSubject binds Read into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func readerStubReadSubject(tb assert.TB) stub.Subject[ReaderReadCall, ReaderReadReturn] {
	tb.Helper()
	s := NewReaderStub(tb)
	return stub.Subject[ReaderReadCall, ReaderReadReturn]{
		Stub: s.OnRead.MethodStub,
		Call: func() {
			_ = s.Read()
		},
		Result: func() ReaderReadReturn {
			got0 := s.Read()
			return ReaderReadReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnRead.Func(func() Cursor {
				mark()
				var z0 Cursor
				return z0
			})
		},
	}
}

// TestReaderStubRead pins how Read answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestReaderStubRead(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Read", readerStubReadSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewReaderStub(t)
		_ = s.Read()
		s.OnRead.AssertCalledOnce(t, "Read must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewReaderStub(t)
		var seen []ReaderReadCall
		s.OnRead.OnRecord(func(c ReaderReadCall) { seen = append(seen, c) })
		_ = s.Read()
		_ = s.Read()
		assert.Length(t, seen, 2, "OnRecord must fire once per Read call")
	})

	t.Run("wires WithReaderRead at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewReaderStub(t, WithReaderRead(func() Cursor {
			called = true
			var z0 Cursor
			return z0
		}))
		_ = s.Read()
		assert.True(t, called, "WithReaderRead must install the override")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewReaderStub(t)
		var want0 Cursor
		s.OnRead.Returns(want0)
		_ = s.Read()
		s.ResetCalls()
		_ = s.Read()
		s.OnRead.AssertCalledOnce(t, "Read must record the call that followed the reset")
	})
}

// readerStubDouble describes how to build a ReaderStub under each option
// whose effect is the same whatever a method's signature.
//
// Read stands in for the double as a whole: these checks assert only that
// an option reached it.
func readerStubDouble() stub.Double[ReaderReadCall] {
	instance := func(s *ReaderStub) stub.Instance[ReaderReadCall] {
		return stub.Instance[ReaderReadCall]{
			Stub: s.OnRead.MethodStub,
			Call: func() {
				_ = s.Read()
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[ReaderReadCall]{
		New: func(tb assert.TB) stub.Instance[ReaderReadCall] {
			return instance(NewReaderStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[ReaderReadCall] {
			return instance(NewReaderStub(tb, ReaderStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[ReaderReadCall] {
			return instance(NewReaderStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[ReaderReadCall] {
			return instance(NewReaderStub(tb, ReaderStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[ReaderReadCall] {
			return instance(NewReaderStub(tb, ReaderStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[ReaderReadCall] {
			return instance(NewReaderStub(tb, ReaderStubStrict()))
		},
	}
}

// TestReaderStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestReaderStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, readerStubDouble())
}

// TestReaderStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestReaderStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewReaderStub(t)
	s := NewReaderStub(t, ReaderStubDelegateTo(inner))

	t.Run("forwards Read to the wrapped implementation", func(t *testing.T) {
		_ = s.Read()
		inner.OnRead.AssertCalledOnce(t, "Read must reach the wrapped implementation")
	})
}

// subscriptionStubCancelSubject binds Cancel into the shape
// [stub.Behaviour] drives: how to call it, what it answers with, and how
// to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func subscriptionStubCancelSubject(tb assert.TB) stub.Subject[SubscriptionCancelCall, struct{}] {
	tb.Helper()
	s := NewSubscriptionStub(tb)
	return stub.Subject[SubscriptionCancelCall, struct{}]{
		Stub: s.OnCancel.MethodStub,
		Call: func() {
			s.Cancel()
		},
		Override: func(mark func()) {
			s.OnCancel.Func(func() {
				mark()
			})
		},
	}
}

// TestSubscriptionStubCancel pins how Cancel answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestSubscriptionStubCancel(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Cancel", subscriptionStubCancelSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewSubscriptionStub(t)
		s.Cancel()
		s.OnCancel.AssertCalledOnce(t, "Cancel must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewSubscriptionStub(t)
		var seen []SubscriptionCancelCall
		s.OnCancel.OnRecord(func(c SubscriptionCancelCall) { seen = append(seen, c) })
		s.Cancel()
		s.Cancel()
		assert.Length(t, seen, 2, "OnRecord must fire once per Cancel call")
	})

	t.Run("wires WithSubscriptionCancel at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewSubscriptionStub(t, WithSubscriptionCancel(func() {
			called = true
		}))
		s.Cancel()
		assert.True(t, called, "WithSubscriptionCancel must install the override")
	})
}

// subscriptionStubDouble describes how to build a SubscriptionStub under
// each option whose effect is the same whatever a method's signature.
//
// Cancel stands in for the double as a whole: these checks assert only
// that an option reached it.
func subscriptionStubDouble() stub.Double[SubscriptionCancelCall] {
	instance := func(s *SubscriptionStub) stub.Instance[SubscriptionCancelCall] {
		return stub.Instance[SubscriptionCancelCall]{
			Stub: s.OnCancel.MethodStub,
			Call: func() {
				s.Cancel()
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[SubscriptionCancelCall]{
		New: func(tb assert.TB) stub.Instance[SubscriptionCancelCall] {
			return instance(NewSubscriptionStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[SubscriptionCancelCall] {
			return instance(NewSubscriptionStub(tb, SubscriptionStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[SubscriptionCancelCall] {
			return instance(NewSubscriptionStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[SubscriptionCancelCall] {
			return instance(NewSubscriptionStub(tb, SubscriptionStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[SubscriptionCancelCall] {
			return instance(NewSubscriptionStub(tb, SubscriptionStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[SubscriptionCancelCall] {
			return instance(NewSubscriptionStub(tb, SubscriptionStubStrict()))
		},
	}
}

// TestSubscriptionStubOptions pins the settings that apply to the double
// rather than to one method: bench mode, the injected clock and random
// source, and what a reset clears.
func TestSubscriptionStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, subscriptionStubDouble())
}

// TestSubscriptionStubDelegateTo pins that a wrapped implementation is
// called through, and recorded on the way past.
func TestSubscriptionStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewSubscriptionStub(t)
	s := NewSubscriptionStub(t, SubscriptionStubDelegateTo(inner))

	t.Run("forwards Cancel to the wrapped implementation", func(t *testing.T) {
		s.Cancel()
		inner.OnCancel.AssertCalledOnce(t, "Cancel must reach the wrapped implementation")
	})
}
