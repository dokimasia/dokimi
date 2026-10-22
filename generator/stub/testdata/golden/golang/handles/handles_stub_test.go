package handles

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/other"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/stub"
)

// archiveStubSize calls Size of s with the samples of its arguments.
func archiveStubSize(tb assert.TB, s *ArchiveStub) int {
	tb.Helper()
	return s.Size()
}

// archiveStubSizePins returns the samples of the results of Size.
func archiveStubSizePins() int {
	var want0 int
	want0 = 42
	return want0
}

// archiveStubOpenLogArgs returns the samples of the arguments of OpenLog.
func archiveStubOpenLogArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "OpenLog-a"
	return a0
}

// archiveStubOpenLog calls OpenLog of s with the samples of its arguments.
func archiveStubOpenLog(tb assert.TB, s *ArchiveStub) (Log, error) {
	tb.Helper()
	a0 := archiveStubOpenLogArgs(tb)
	return s.OpenLog(a0)
}

// archiveStubOpenLogPins returns the samples of the results of OpenLog,
// with a new error as its error result.
func archiveStubOpenLogPins() (Log, error) {
	var want0 Log
	want1 := errors.New("OpenLog-returns")
	return want0, want1
}

// archiveStubSizeSubject returns the subject of [stub.Behaviour] for Size.
// The subject calls a new double that reports to tb.
func archiveStubSizeSubject(tb assert.TB) stub.Subject[ArchiveSizeCall, ArchiveSizeReturn] {
	tb.Helper()
	s := NewArchiveStub(tb)
	return stub.Subject[ArchiveSizeCall, ArchiveSizeReturn]{
		Stub: s.OnSize.MethodStub,
		Call: func() {
			_ = archiveStubSize(tb, s)
		},
		Result: func() ArchiveSizeReturn {
			got0 := archiveStubSize(tb, s)
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

// archiveStubOpenLogSubject returns the subject of [stub.Behaviour] for
// OpenLog. The subject calls a new double that reports to tb.
func archiveStubOpenLogSubject(
	tb assert.TB,
) stub.Subject[ArchiveOpenLogCall, ArchiveOpenLogReturn] {
	tb.Helper()
	s := NewArchiveStub(tb)
	return stub.Subject[ArchiveOpenLogCall, ArchiveOpenLogReturn]{
		Stub: s.OnOpenLog.MethodStub,
		Call: func() {
			_, _ = archiveStubOpenLog(tb, s)
		},
		Result: func() ArchiveOpenLogReturn {
			got0, got1 := archiveStubOpenLog(tb, s)
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
			_, r1 := archiveStubOpenLog(tb, s)
			return r1
		},
	}
}

// archiveStubInstance is the instance of a ArchiveStub that the checks of
// the constructors and the options drive.
type archiveStubInstance = stub.Instance[ArchiveSizeCall]

// archiveStubInstanceOf returns the instance of s. The options apply to
// the whole double, so the checks call Size alone.
func archiveStubInstanceOf(tb assert.TB, s *ArchiveStub) archiveStubInstance {
	tb.Helper()
	return archiveStubInstance{
		Stub: s.OnSize.MethodStub,
		Call: func() {
			_ = archiveStubSize(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestArchiveStub checks the methods, the constructors and the options of
// the double ArchiveStub.
func TestArchiveStub(t *testing.T) {
	t.Parallel()

	t.Run("Size", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Size", archiveStubSizeSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			var seen []ArchiveSizeCall
			s.OnSize.OnRecord(func(c ArchiveSizeCall) { seen = append(seen, c) })
			_ = archiveStubSize(t, s)
			_ = archiveStubSize(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Size")
		})
	})

	t.Run("OpenLog", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "OpenLog", archiveStubOpenLogSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			_, _ = archiveStubOpenLog(t, s)
			got := s.OnOpenLog.AssertCalledOnce(t, "OpenLog records the call")
			a0 := archiveStubOpenLogArgs(t)
			assert.Equal(t, got.Name, a0, "the call records Name")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			var seen []ArchiveOpenLogCall
			s.OnOpenLog.OnRecord(func(c ArchiveOpenLogCall) { seen = append(seen, c) })
			_, _ = archiveStubOpenLog(t, s)
			_, _ = archiveStubOpenLog(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of OpenLog")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) archiveStubInstance {
			return archiveStubInstanceOf(tb, NewArchiveStub(tb))
		})

		t.Run("clears the recorded calls of Size", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			_ = archiveStubSize(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnSize.CallCount(), 0, "ResetCalls clears the calls of Size")
		})

		t.Run("keeps the results of Returns of Size", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			want0 := archiveStubSizePins()
			s.OnSize.Returns(want0)
			s.ResetCalls()
			got0 := archiveStubSize(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Size")
		})

		t.Run("clears the recorded calls of OpenLog", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			_, _ = archiveStubOpenLog(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnOpenLog.CallCount(), 0, "ResetCalls clears the calls of OpenLog")
		})

		t.Run("keeps the results of Returns of OpenLog", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			want0, want1 := archiveStubOpenLogPins()
			s.OnOpenLog.Returns(want0, want1)
			s.ResetCalls()
			_, got1 := archiveStubOpenLog(t, s)
			assert.ErrorIs(t, got1, want1, "ResetCalls keeps the error of OpenLog")
		})
	})

	t.Run("NewArchiveStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) archiveStubInstance {
				return archiveStubInstanceOf(tb, NewArchiveStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewArchiveStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails OpenLog with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("OpenLog-failing")
			s := NewArchiveStubFailing(t, want)
			_, r1 := archiveStubOpenLog(t, s)
			assert.ErrorIs(t, r1, want, "OpenLog fails with the error of NewArchiveStubFailing")
		})
	})

	t.Run("ArchiveStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) archiveStubInstance {
			return archiveStubInstanceOf(tb, NewArchiveStub(tb, ArchiveStubStrict()))
		})
	})

	t.Run("ArchiveStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Size to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewArchiveStub(t)
			s := NewArchiveStub(t, ArchiveStubDelegateTo(inner))
			_ = archiveStubSize(t, s)
			inner.OnSize.AssertCalledOnce(t, "Size calls the implementation")
		})

		t.Run("passes OpenLog to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewArchiveStub(t)
			s := NewArchiveStub(t, ArchiveStubDelegateTo(inner))
			_, _ = archiveStubOpenLog(t, s)
			inner.OnOpenLog.AssertCalledOnce(t, "OpenLog calls the implementation")
		})

		t.Run("fails OpenLog with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("OpenLog-delegate")
			inner := NewArchiveStub(t)
			inner.OnOpenLog.Faults(want, 1)
			s := NewArchiveStub(t, ArchiveStubDelegateTo(inner))
			_, r1 := archiveStubOpenLog(t, s)
			assert.ErrorIs(t, r1, want, "OpenLog fails with the error of the implementation")
		})
	})

	t.Run("ArchiveStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) archiveStubInstance {
			return archiveStubInstanceOf(tb, NewArchiveStub(tb, ArchiveStubWithClock(clk)))
		})
	})

	t.Run("ArchiveStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) archiveStubInstance {
			return archiveStubInstanceOf(tb, NewArchiveStub(tb, ArchiveStubWithRandSource(src)))
		})
	})

	t.Run("ArchiveStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) archiveStubInstance {
			return archiveStubInstanceOf(tb, NewArchiveStub(tb, ArchiveStubBenchMode()))
		})
	})

	t.Run("ArchiveStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) archiveStubInstance {
				return archiveStubInstanceOf(
					tb,
					NewArchiveStub(tb, ArchiveStubWithClock(clk), ArchiveStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("ArchiveStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails OpenLog with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("OpenLog-failing")
			s := NewArchiveStub(t, ArchiveStubFailing(want))
			_, r1 := archiveStubOpenLog(t, s)
			assert.ErrorIs(t, r1, want, "OpenLog fails with the error of ArchiveStubFailing")
		})
	})

	t.Run("WithArchiveSize", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Size", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewArchiveStub(t, WithArchiveSize(func() int {
				called = true
				var z0 int
				return z0
			}))
			_ = archiveStubSize(t, s)
			assert.True(t, called, "WithArchiveSize sets the function of Size")
		})
	})

	t.Run("WithArchiveOpenLog", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of OpenLog", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewArchiveStub(t, WithArchiveOpenLog(func(_ string) (Log, error) {
				called = true
				var z0 Log
				var z1 error
				return z0, z1
			}))
			_, _ = archiveStubOpenLog(t, s)
			assert.True(t, called, "WithArchiveOpenLog sets the function of OpenLog")
		})
	})
}

// TestArchiveSizeStub checks the methods of ArchiveSizeStub, which
// configures Size.
func TestArchiveSizeStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Size", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			want0 := archiveStubSizePins()
			s.OnSize.Returns(want0)
			got0 := archiveStubSize(t, s)
			assert.Equal(t, got0, want0, "Size returns the Result that Returns sets")
		})

		t.Run("configures Size for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewArchiveStub(r)
			s.OnSize.Strict()
			var want0 int
			s.OnSize.Returns(want0)
			_ = archiveStubSize(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Size")
		})
	})
}

// TestArchiveOpenLogStub checks the methods of ArchiveOpenLogStub, which
// configures OpenLog.
func TestArchiveOpenLogStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of OpenLog", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			want0, want1 := archiveStubOpenLogPins()
			s.OnOpenLog.Returns(want0, want1)
			_, got1 := archiveStubOpenLog(t, s)
			assert.ErrorIs(t, got1, want1, "OpenLog returns the error that Returns sets")
		})

		t.Run("configures OpenLog for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewArchiveStub(r)
			s.OnOpenLog.Strict()
			var want0 Log
			var want1 error
			s.OnOpenLog.Returns(want0, want1)
			_, _ = archiveStubOpenLog(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of OpenLog")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of OpenLog to the function", func(t *testing.T) {
			t.Parallel()
			s := NewArchiveStub(t)
			var got0 string
			s.OnOpenLog.Func(func(p0 string) (Log, error) {
				got0 = p0
				var z0 Log
				var z1 error
				return z0, z1
			})
			_, _ = archiveStubOpenLog(t, s)
			a0 := archiveStubOpenLogArgs(t)
			assert.Equal(t, got0, a0, "the function receives Name")
		})
	})
}

// logStubAppendArgs returns the samples of the arguments of Append.
func logStubAppendArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Append-a"
	return a0
}

// logStubAppend calls Append of s with the samples of its arguments.
func logStubAppend(tb assert.TB, s *LogStub) error {
	tb.Helper()
	a0 := logStubAppendArgs(tb)
	return s.Append(a0)
}

// logStubAppendPins returns the samples of the results of Append, with a
// new error as its error result.
func logStubAppendPins() error {
	want0 := errors.New("Append-returns")
	return want0
}

// logStubAppendSubject returns the subject of [stub.Behaviour] for Append.
// The subject calls a new double that reports to tb.
func logStubAppendSubject(tb assert.TB) stub.Subject[LogAppendCall, LogAppendReturn] {
	tb.Helper()
	s := NewLogStub(tb)
	return stub.Subject[LogAppendCall, LogAppendReturn]{
		Stub: s.OnAppend.MethodStub,
		Call: func() {
			_ = logStubAppend(tb, s)
		},
		Result: func() LogAppendReturn {
			got0 := logStubAppend(tb, s)
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
			r0 := logStubAppend(tb, s)
			return r0
		},
	}
}

// logStubInstance is the instance of a LogStub that the checks of the
// constructors and the options drive.
type logStubInstance = stub.Instance[LogAppendCall]

// logStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Append alone.
func logStubInstanceOf(tb assert.TB, s *LogStub) logStubInstance {
	tb.Helper()
	return logStubInstance{
		Stub: s.OnAppend.MethodStub,
		Call: func() {
			_ = logStubAppend(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestLogStub checks the methods, the constructors and the options of the
// double LogStub.
func TestLogStub(t *testing.T) {
	t.Parallel()

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Append", logStubAppendSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewLogStub(t)
			_ = logStubAppend(t, s)
			got := s.OnAppend.AssertCalledOnce(t, "Append records the call")
			a0 := logStubAppendArgs(t)
			assert.Equal(t, got.Line, a0, "the call records Line")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewLogStub(t)
			var seen []LogAppendCall
			s.OnAppend.OnRecord(func(c LogAppendCall) { seen = append(seen, c) })
			_ = logStubAppend(t, s)
			_ = logStubAppend(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Append")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) logStubInstance {
			return logStubInstanceOf(tb, NewLogStub(tb))
		})

		t.Run("clears the recorded calls of Append", func(t *testing.T) {
			t.Parallel()
			s := NewLogStub(t)
			_ = logStubAppend(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnAppend.CallCount(), 0, "ResetCalls clears the calls of Append")
		})

		t.Run("keeps the results of Returns of Append", func(t *testing.T) {
			t.Parallel()
			s := NewLogStub(t)
			want0 := logStubAppendPins()
			s.OnAppend.Returns(want0)
			s.ResetCalls()
			got0 := logStubAppend(t, s)
			assert.ErrorIs(t, got0, want0, "ResetCalls keeps the error of Append")
		})
	})

	t.Run("NewLogStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) logStubInstance {
				return logStubInstanceOf(tb, NewLogStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewLogStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Append with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Append-failing")
			s := NewLogStubFailing(t, want)
			r0 := logStubAppend(t, s)
			assert.ErrorIs(t, r0, want, "Append fails with the error of NewLogStubFailing")
		})
	})

	t.Run("LogStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) logStubInstance {
			return logStubInstanceOf(tb, NewLogStub(tb, LogStubStrict()))
		})
	})

	t.Run("LogStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Append to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewLogStub(t)
			s := NewLogStub(t, LogStubDelegateTo(inner))
			_ = logStubAppend(t, s)
			inner.OnAppend.AssertCalledOnce(t, "Append calls the implementation")
		})

		t.Run("fails Append with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Append-delegate")
			inner := NewLogStub(t)
			inner.OnAppend.Faults(want, 1)
			s := NewLogStub(t, LogStubDelegateTo(inner))
			r0 := logStubAppend(t, s)
			assert.ErrorIs(t, r0, want, "Append fails with the error of the implementation")
		})
	})

	t.Run("LogStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) logStubInstance {
			return logStubInstanceOf(tb, NewLogStub(tb, LogStubWithClock(clk)))
		})
	})

	t.Run("LogStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) logStubInstance {
			return logStubInstanceOf(tb, NewLogStub(tb, LogStubWithRandSource(src)))
		})
	})

	t.Run("LogStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) logStubInstance {
			return logStubInstanceOf(tb, NewLogStub(tb, LogStubBenchMode()))
		})
	})

	t.Run("LogStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) logStubInstance {
				return logStubInstanceOf(
					tb,
					NewLogStub(tb, LogStubWithClock(clk), LogStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("LogStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Append with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Append-failing")
			s := NewLogStub(t, LogStubFailing(want))
			r0 := logStubAppend(t, s)
			assert.ErrorIs(t, r0, want, "Append fails with the error of LogStubFailing")
		})
	})

	t.Run("WithLogAppend", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Append", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewLogStub(t, WithLogAppend(func(_ string) error {
				called = true
				var z0 error
				return z0
			}))
			_ = logStubAppend(t, s)
			assert.True(t, called, "WithLogAppend sets the function of Append")
		})
	})
}

// TestLogAppendStub checks the methods of LogAppendStub, which configures
// Append.
func TestLogAppendStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Append", func(t *testing.T) {
			t.Parallel()
			s := NewLogStub(t)
			want0 := logStubAppendPins()
			s.OnAppend.Returns(want0)
			got0 := logStubAppend(t, s)
			assert.ErrorIs(t, got0, want0, "Append returns the error that Returns sets")
		})

		t.Run("configures Append for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewLogStub(r)
			s.OnAppend.Strict()
			var want0 error
			s.OnAppend.Returns(want0)
			_ = logStubAppend(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Append")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Append to the function", func(t *testing.T) {
			t.Parallel()
			s := NewLogStub(t)
			var got0 string
			s.OnAppend.Func(func(p0 string) error {
				got0 = p0
				var z0 error
				return z0
			})
			_ = logStubAppend(t, s)
			a0 := logStubAppendArgs(t)
			assert.Equal(t, got0, a0, "the function receives Line")
		})
	})
}

// producerStubOpenArgs returns the samples of the arguments of Open.
func producerStubOpenArgs(tb assert.TB) context.Context {
	tb.Helper()
	a0 := assert.Context(tb)
	return a0
}

// producerStubOpen calls Open of s with the samples of its arguments.
func producerStubOpen(tb assert.TB, s *ProducerStub) (Cursor, error) {
	tb.Helper()
	a0 := producerStubOpenArgs(tb)
	return s.Open(a0)
}

// producerStubOpenPins returns the samples of the results of Open, with a
// new error as its error result.
func producerStubOpenPins() (Cursor, error) {
	var want0 Cursor
	want1 := errors.New("Open-returns")
	return want0, want1
}

// producerStubResumeArgs returns the samples of the arguments of Resume.
func producerStubResumeArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Resume-a"
	return a0
}

// producerStubResume calls Resume of s with the samples of its arguments.
func producerStubResume(tb assert.TB, s *ProducerStub) Cursor {
	tb.Helper()
	a0 := producerStubResumeArgs(tb)
	return s.Resume(a0)
}

// producerStubWatch calls Watch of s with the samples of its arguments.
func producerStubWatch(tb assert.TB, s *ProducerStub) Subscription {
	tb.Helper()
	return s.Watch()
}

// producerStubNothing calls Nothing of s with the samples of its
// arguments.
func producerStubNothing(tb assert.TB, s *ProducerStub) Empty {
	tb.Helper()
	return s.Nothing()
}

// producerStubElsewhere calls Elsewhere of s with the samples of its
// arguments.
func producerStubElsewhere(tb assert.TB, s *ProducerStub) other.Thing {
	tb.Helper()
	return s.Elsewhere()
}

// producerStubOpenSubject returns the subject of [stub.Behaviour] for
// Open. The subject calls a new double that reports to tb.
func producerStubOpenSubject(tb assert.TB) stub.Subject[ProducerOpenCall, ProducerOpenReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerOpenCall, ProducerOpenReturn]{
		Stub: s.OnOpen.MethodStub,
		Call: func() {
			_, _ = producerStubOpen(tb, s)
		},
		Result: func() ProducerOpenReturn {
			got0, got1 := producerStubOpen(tb, s)
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
			_, r1 := producerStubOpen(tb, s)
			return r1
		},
	}
}

// producerStubResumeSubject returns the subject of [stub.Behaviour] for
// Resume. The subject calls a new double that reports to tb.
func producerStubResumeSubject(
	tb assert.TB,
) stub.Subject[ProducerResumeCall, ProducerResumeReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerResumeCall, ProducerResumeReturn]{
		Stub: s.OnResume.MethodStub,
		Call: func() {
			_ = producerStubResume(tb, s)
		},
		Result: func() ProducerResumeReturn {
			got0 := producerStubResume(tb, s)
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

// producerStubWatchSubject returns the subject of [stub.Behaviour] for
// Watch. The subject calls a new double that reports to tb.
func producerStubWatchSubject(tb assert.TB) stub.Subject[ProducerWatchCall, ProducerWatchReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerWatchCall, ProducerWatchReturn]{
		Stub: s.OnWatch.MethodStub,
		Call: func() {
			_ = producerStubWatch(tb, s)
		},
		Result: func() ProducerWatchReturn {
			got0 := producerStubWatch(tb, s)
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

// producerStubNothingSubject returns the subject of [stub.Behaviour] for
// Nothing. The subject calls a new double that reports to tb.
func producerStubNothingSubject(
	tb assert.TB,
) stub.Subject[ProducerNothingCall, ProducerNothingReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerNothingCall, ProducerNothingReturn]{
		Stub: s.OnNothing.MethodStub,
		Call: func() {
			_ = producerStubNothing(tb, s)
		},
		Result: func() ProducerNothingReturn {
			got0 := producerStubNothing(tb, s)
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

// producerStubElsewhereSubject returns the subject of [stub.Behaviour] for
// Elsewhere. The subject calls a new double that reports to tb.
func producerStubElsewhereSubject(
	tb assert.TB,
) stub.Subject[ProducerElsewhereCall, ProducerElsewhereReturn] {
	tb.Helper()
	s := NewProducerStub(tb)
	return stub.Subject[ProducerElsewhereCall, ProducerElsewhereReturn]{
		Stub: s.OnElsewhere.MethodStub,
		Call: func() {
			_ = producerStubElsewhere(tb, s)
		},
		Result: func() ProducerElsewhereReturn {
			got0 := producerStubElsewhere(tb, s)
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

// producerStubInstance is the instance of a ProducerStub that the checks
// of the constructors and the options drive.
type producerStubInstance = stub.Instance[ProducerOpenCall]

// producerStubInstanceOf returns the instance of s. The options apply to
// the whole double, so the checks call Open alone.
func producerStubInstanceOf(tb assert.TB, s *ProducerStub) producerStubInstance {
	tb.Helper()
	return producerStubInstance{
		Stub: s.OnOpen.MethodStub,
		Call: func() {
			_, _ = producerStubOpen(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestProducerStub checks the methods, the constructors and the options of
// the double ProducerStub.
func TestProducerStub(t *testing.T) {
	t.Parallel()

	t.Run("Open", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Open", producerStubOpenSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			_, _ = producerStubOpen(t, s)
			got := s.OnOpen.AssertCalledOnce(t, "Open records the call")
			a0 := producerStubOpenArgs(t)
			assert.Equal(t, got.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			var seen []ProducerOpenCall
			s.OnOpen.OnRecord(func(c ProducerOpenCall) { seen = append(seen, c) })
			_, _ = producerStubOpen(t, s)
			_, _ = producerStubOpen(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Open")
		})
	})

	t.Run("Resume", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Resume", producerStubResumeSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			_ = producerStubResume(t, s)
			got := s.OnResume.AssertCalledOnce(t, "Resume records the call")
			a0 := producerStubResumeArgs(t)
			assert.Equal(t, got.Token, a0, "the call records Token")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			var seen []ProducerResumeCall
			s.OnResume.OnRecord(func(c ProducerResumeCall) { seen = append(seen, c) })
			_ = producerStubResume(t, s)
			_ = producerStubResume(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Resume")
		})
	})

	t.Run("Watch", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Watch", producerStubWatchSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			var seen []ProducerWatchCall
			s.OnWatch.OnRecord(func(c ProducerWatchCall) { seen = append(seen, c) })
			_ = producerStubWatch(t, s)
			_ = producerStubWatch(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Watch")
		})
	})

	t.Run("Nothing", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Nothing", producerStubNothingSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			var seen []ProducerNothingCall
			s.OnNothing.OnRecord(func(c ProducerNothingCall) { seen = append(seen, c) })
			_ = producerStubNothing(t, s)
			_ = producerStubNothing(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Nothing")
		})
	})

	t.Run("Elsewhere", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Elsewhere", producerStubElsewhereSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			var seen []ProducerElsewhereCall
			s.OnElsewhere.OnRecord(func(c ProducerElsewhereCall) { seen = append(seen, c) })
			_ = producerStubElsewhere(t, s)
			_ = producerStubElsewhere(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Elsewhere")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) producerStubInstance {
			return producerStubInstanceOf(tb, NewProducerStub(tb))
		})

		t.Run("clears the recorded calls of Open", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			_, _ = producerStubOpen(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnOpen.CallCount(), 0, "ResetCalls clears the calls of Open")
		})

		t.Run("keeps the results of Returns of Open", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			want0, want1 := producerStubOpenPins()
			s.OnOpen.Returns(want0, want1)
			s.ResetCalls()
			_, got1 := producerStubOpen(t, s)
			assert.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Open")
		})

		t.Run("clears the recorded calls of Resume", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			_ = producerStubResume(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnResume.CallCount(), 0, "ResetCalls clears the calls of Resume")
		})

		t.Run("clears the recorded calls of Watch", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			_ = producerStubWatch(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnWatch.CallCount(), 0, "ResetCalls clears the calls of Watch")
		})

		t.Run("clears the recorded calls of Nothing", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			_ = producerStubNothing(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnNothing.CallCount(), 0, "ResetCalls clears the calls of Nothing")
		})

		t.Run("clears the recorded calls of Elsewhere", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			_ = producerStubElsewhere(t, s)
			s.ResetCalls()
			assert.Equal(
				t,
				s.OnElsewhere.CallCount(),
				0,
				"ResetCalls clears the calls of Elsewhere",
			)
		})
	})

	t.Run("NewProducerStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) producerStubInstance {
				return producerStubInstanceOf(tb, NewProducerStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewProducerStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Open with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Open-failing")
			s := NewProducerStubFailing(t, want)
			_, r1 := producerStubOpen(t, s)
			assert.ErrorIs(t, r1, want, "Open fails with the error of NewProducerStubFailing")
		})
	})

	t.Run("ProducerStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) producerStubInstance {
			return producerStubInstanceOf(tb, NewProducerStub(tb, ProducerStubStrict()))
		})
	})

	t.Run("ProducerStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Open to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewProducerStub(t)
			s := NewProducerStub(t, ProducerStubDelegateTo(inner))
			_, _ = producerStubOpen(t, s)
			inner.OnOpen.AssertCalledOnce(t, "Open calls the implementation")
		})

		t.Run("fails Open with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Open-delegate")
			inner := NewProducerStub(t)
			inner.OnOpen.Faults(want, 1)
			s := NewProducerStub(t, ProducerStubDelegateTo(inner))
			_, r1 := producerStubOpen(t, s)
			assert.ErrorIs(t, r1, want, "Open fails with the error of the implementation")
		})

		t.Run("passes Resume to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewProducerStub(t)
			s := NewProducerStub(t, ProducerStubDelegateTo(inner))
			_ = producerStubResume(t, s)
			inner.OnResume.AssertCalledOnce(t, "Resume calls the implementation")
		})

		t.Run("passes Watch to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewProducerStub(t)
			s := NewProducerStub(t, ProducerStubDelegateTo(inner))
			_ = producerStubWatch(t, s)
			inner.OnWatch.AssertCalledOnce(t, "Watch calls the implementation")
		})

		t.Run("passes Nothing to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewProducerStub(t)
			s := NewProducerStub(t, ProducerStubDelegateTo(inner))
			_ = producerStubNothing(t, s)
			inner.OnNothing.AssertCalledOnce(t, "Nothing calls the implementation")
		})

		t.Run("passes Elsewhere to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewProducerStub(t)
			s := NewProducerStub(t, ProducerStubDelegateTo(inner))
			_ = producerStubElsewhere(t, s)
			inner.OnElsewhere.AssertCalledOnce(t, "Elsewhere calls the implementation")
		})
	})

	t.Run("ProducerStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) producerStubInstance {
			return producerStubInstanceOf(tb, NewProducerStub(tb, ProducerStubWithClock(clk)))
		})
	})

	t.Run("ProducerStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) producerStubInstance {
			return producerStubInstanceOf(tb, NewProducerStub(tb, ProducerStubWithRandSource(src)))
		})
	})

	t.Run("ProducerStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) producerStubInstance {
			return producerStubInstanceOf(tb, NewProducerStub(tb, ProducerStubBenchMode()))
		})
	})

	t.Run("ProducerStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) producerStubInstance {
				return producerStubInstanceOf(
					tb,
					NewProducerStub(tb, ProducerStubWithClock(clk), ProducerStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("ProducerStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Open with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Open-failing")
			s := NewProducerStub(t, ProducerStubFailing(want))
			_, r1 := producerStubOpen(t, s)
			assert.ErrorIs(t, r1, want, "Open fails with the error of ProducerStubFailing")
		})
	})

	t.Run("WithProducerOpen", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Open", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewProducerStub(t, WithProducerOpen(func(_ context.Context) (Cursor, error) {
				called = true
				var z0 Cursor
				var z1 error
				return z0, z1
			}))
			_, _ = producerStubOpen(t, s)
			assert.True(t, called, "WithProducerOpen sets the function of Open")
		})
	})

	t.Run("WithProducerResume", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Resume", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewProducerStub(t, WithProducerResume(func(_ string) Cursor {
				called = true
				var z0 Cursor
				return z0
			}))
			_ = producerStubResume(t, s)
			assert.True(t, called, "WithProducerResume sets the function of Resume")
		})
	})

	t.Run("WithProducerWatch", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Watch", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewProducerStub(t, WithProducerWatch(func() Subscription {
				called = true
				var z0 Subscription
				return z0
			}))
			_ = producerStubWatch(t, s)
			assert.True(t, called, "WithProducerWatch sets the function of Watch")
		})
	})

	t.Run("WithProducerNothing", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Nothing", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewProducerStub(t, WithProducerNothing(func() Empty {
				called = true
				var z0 Empty
				return z0
			}))
			_ = producerStubNothing(t, s)
			assert.True(t, called, "WithProducerNothing sets the function of Nothing")
		})
	})

	t.Run("WithProducerElsewhere", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Elsewhere", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewProducerStub(t, WithProducerElsewhere(func() other.Thing {
				called = true
				var z0 other.Thing
				return z0
			}))
			_ = producerStubElsewhere(t, s)
			assert.True(t, called, "WithProducerElsewhere sets the function of Elsewhere")
		})
	})
}

// TestProducerOpenStub checks the methods of ProducerOpenStub, which
// configures Open.
func TestProducerOpenStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Open", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			want0, want1 := producerStubOpenPins()
			s.OnOpen.Returns(want0, want1)
			_, got1 := producerStubOpen(t, s)
			assert.ErrorIs(t, got1, want1, "Open returns the error that Returns sets")
		})

		t.Run("configures Open for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewProducerStub(r)
			s.OnOpen.Strict()
			var want0 Cursor
			var want1 error
			s.OnOpen.Returns(want0, want1)
			_, _ = producerStubOpen(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Open")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Open to the function", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			var got0 context.Context
			s.OnOpen.Func(func(p0 context.Context) (Cursor, error) {
				got0 = p0
				var z0 Cursor
				var z1 error
				return z0, z1
			})
			_, _ = producerStubOpen(t, s)
			a0 := producerStubOpenArgs(t)
			assert.Equal(t, got0, a0, "the function receives Ctx", assert.ByIdentity())
		})
	})
}

// TestProducerResumeStub checks the methods of ProducerResumeStub, which
// configures Resume.
func TestProducerResumeStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("configures Resume for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewProducerStub(r)
			s.OnResume.Strict()
			var want0 Cursor
			s.OnResume.Returns(want0)
			_ = producerStubResume(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Resume")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Resume to the function", func(t *testing.T) {
			t.Parallel()
			s := NewProducerStub(t)
			var got0 string
			s.OnResume.Func(func(p0 string) Cursor {
				got0 = p0
				var z0 Cursor
				return z0
			})
			_ = producerStubResume(t, s)
			a0 := producerStubResumeArgs(t)
			assert.Equal(t, got0, a0, "the function receives Token")
		})
	})
}

// TestProducerWatchStub checks the methods of ProducerWatchStub, which
// configures Watch.
func TestProducerWatchStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("configures Watch for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewProducerStub(r)
			s.OnWatch.Strict()
			var want0 Subscription
			s.OnWatch.Returns(want0)
			_ = producerStubWatch(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Watch")
		})
	})
}

// TestProducerNothingStub checks the methods of ProducerNothingStub, which
// configures Nothing.
func TestProducerNothingStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("configures Nothing for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewProducerStub(r)
			s.OnNothing.Strict()
			var want0 Empty
			s.OnNothing.Returns(want0)
			_ = producerStubNothing(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Nothing")
		})
	})
}

// TestProducerElsewhereStub checks the methods of ProducerElsewhereStub,
// which configures Elsewhere.
func TestProducerElsewhereStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("configures Elsewhere for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewProducerStub(r)
			s.OnElsewhere.Strict()
			var want0 other.Thing
			s.OnElsewhere.Returns(want0)
			_ = producerStubElsewhere(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Elsewhere")
		})
	})
}

// cursorStubNextArgs returns the samples of the arguments of Next.
func cursorStubNextArgs(tb assert.TB) context.Context {
	tb.Helper()
	a0 := assert.Context(tb)
	return a0
}

// cursorStubNext calls Next of s with the samples of its arguments.
func cursorStubNext(tb assert.TB, s *CursorStub) (string, bool, error) {
	tb.Helper()
	a0 := cursorStubNextArgs(tb)
	return s.Next(a0)
}

// cursorStubNextPins returns the samples of the results of Next, with a
// new error as its error result.
func cursorStubNextPins() (string, bool, error) {
	var want0 string
	var want1 bool
	want2 := errors.New("Next-returns")
	want0 = "Next-a"
	want1 = true
	return want0, want1, want2
}

// cursorStubCloseArgs returns the samples of the arguments of Close.
func cursorStubCloseArgs(tb assert.TB) context.Context {
	tb.Helper()
	a0 := assert.Context(tb)
	return a0
}

// cursorStubClose calls Close of s with the samples of its arguments.
func cursorStubClose(tb assert.TB, s *CursorStub) error {
	tb.Helper()
	a0 := cursorStubCloseArgs(tb)
	return s.Close(a0)
}

// cursorStubClosePins returns the samples of the results of Close, with a
// new error as its error result.
func cursorStubClosePins() error {
	want0 := errors.New("Close-returns")
	return want0
}

// cursorStubNextSubject returns the subject of [stub.Behaviour] for Next.
// The subject calls a new double that reports to tb.
func cursorStubNextSubject(tb assert.TB) stub.Subject[CursorNextCall, CursorNextReturn] {
	tb.Helper()
	s := NewCursorStub(tb)
	return stub.Subject[CursorNextCall, CursorNextReturn]{
		Stub: s.OnNext.MethodStub,
		Call: func() {
			_, _, _ = cursorStubNext(tb, s)
		},
		Result: func() CursorNextReturn {
			got0, got1, got2 := cursorStubNext(tb, s)
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
			_, _, r2 := cursorStubNext(tb, s)
			return r2
		},
	}
}

// cursorStubCloseSubject returns the subject of [stub.Behaviour] for
// Close. The subject calls a new double that reports to tb.
func cursorStubCloseSubject(tb assert.TB) stub.Subject[CursorCloseCall, CursorCloseReturn] {
	tb.Helper()
	s := NewCursorStub(tb)
	return stub.Subject[CursorCloseCall, CursorCloseReturn]{
		Stub: s.OnClose.MethodStub,
		Call: func() {
			_ = cursorStubClose(tb, s)
		},
		Result: func() CursorCloseReturn {
			got0 := cursorStubClose(tb, s)
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
			r0 := cursorStubClose(tb, s)
			return r0
		},
	}
}

// cursorStubInstance is the instance of a CursorStub that the checks of
// the constructors and the options drive.
type cursorStubInstance = stub.Instance[CursorNextCall]

// cursorStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Next alone.
func cursorStubInstanceOf(tb assert.TB, s *CursorStub) cursorStubInstance {
	tb.Helper()
	return cursorStubInstance{
		Stub: s.OnNext.MethodStub,
		Call: func() {
			_, _, _ = cursorStubNext(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestCursorStub checks the methods, the constructors and the options of
// the double CursorStub.
func TestCursorStub(t *testing.T) {
	t.Parallel()

	t.Run("Next", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Next", cursorStubNextSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			_, _, _ = cursorStubNext(t, s)
			got := s.OnNext.AssertCalledOnce(t, "Next records the call")
			a0 := cursorStubNextArgs(t)
			assert.Equal(t, got.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			var seen []CursorNextCall
			s.OnNext.OnRecord(func(c CursorNextCall) { seen = append(seen, c) })
			_, _, _ = cursorStubNext(t, s)
			_, _, _ = cursorStubNext(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Next")
		})
	})

	t.Run("Close", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Close", cursorStubCloseSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			_ = cursorStubClose(t, s)
			got := s.OnClose.AssertCalledOnce(t, "Close records the call")
			a0 := cursorStubCloseArgs(t)
			assert.Equal(t, got.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			var seen []CursorCloseCall
			s.OnClose.OnRecord(func(c CursorCloseCall) { seen = append(seen, c) })
			_ = cursorStubClose(t, s)
			_ = cursorStubClose(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Close")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) cursorStubInstance {
			return cursorStubInstanceOf(tb, NewCursorStub(tb))
		})

		t.Run("clears the recorded calls of Next", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			_, _, _ = cursorStubNext(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnNext.CallCount(), 0, "ResetCalls clears the calls of Next")
		})

		t.Run("keeps the results of Returns of Next", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			want0, want1, want2 := cursorStubNextPins()
			s.OnNext.Returns(want0, want1, want2)
			s.ResetCalls()
			got0, got1, got2 := cursorStubNext(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result0 of Next")
			expect.Equal(t, got1, want1, "ResetCalls keeps the Result1 of Next")
			expect.ErrorIs(t, got2, want2, "ResetCalls keeps the error of Next")
		})

		t.Run("clears the recorded calls of Close", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			_ = cursorStubClose(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnClose.CallCount(), 0, "ResetCalls clears the calls of Close")
		})

		t.Run("keeps the results of Returns of Close", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			want0 := cursorStubClosePins()
			s.OnClose.Returns(want0)
			s.ResetCalls()
			got0 := cursorStubClose(t, s)
			assert.ErrorIs(t, got0, want0, "ResetCalls keeps the error of Close")
		})
	})

	t.Run("NewCursorStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) cursorStubInstance {
				return cursorStubInstanceOf(tb, NewCursorStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewCursorStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Next with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Next-failing")
			s := NewCursorStubFailing(t, want)
			_, _, r2 := cursorStubNext(t, s)
			assert.ErrorIs(t, r2, want, "Next fails with the error of NewCursorStubFailing")
		})

		t.Run("fails Close with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Close-failing")
			s := NewCursorStubFailing(t, want)
			r0 := cursorStubClose(t, s)
			assert.ErrorIs(t, r0, want, "Close fails with the error of NewCursorStubFailing")
		})
	})

	t.Run("CursorStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) cursorStubInstance {
			return cursorStubInstanceOf(tb, NewCursorStub(tb, CursorStubStrict()))
		})
	})

	t.Run("CursorStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Next to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewCursorStub(t)
			s := NewCursorStub(t, CursorStubDelegateTo(inner))
			_, _, _ = cursorStubNext(t, s)
			inner.OnNext.AssertCalledOnce(t, "Next calls the implementation")
		})

		t.Run("fails Next with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Next-delegate")
			inner := NewCursorStub(t)
			inner.OnNext.Faults(want, 1)
			s := NewCursorStub(t, CursorStubDelegateTo(inner))
			_, _, r2 := cursorStubNext(t, s)
			assert.ErrorIs(t, r2, want, "Next fails with the error of the implementation")
		})

		t.Run("passes Close to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewCursorStub(t)
			s := NewCursorStub(t, CursorStubDelegateTo(inner))
			_ = cursorStubClose(t, s)
			inner.OnClose.AssertCalledOnce(t, "Close calls the implementation")
		})

		t.Run("fails Close with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Close-delegate")
			inner := NewCursorStub(t)
			inner.OnClose.Faults(want, 1)
			s := NewCursorStub(t, CursorStubDelegateTo(inner))
			r0 := cursorStubClose(t, s)
			assert.ErrorIs(t, r0, want, "Close fails with the error of the implementation")
		})
	})

	t.Run("CursorStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) cursorStubInstance {
			return cursorStubInstanceOf(tb, NewCursorStub(tb, CursorStubWithClock(clk)))
		})
	})

	t.Run("CursorStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) cursorStubInstance {
			return cursorStubInstanceOf(tb, NewCursorStub(tb, CursorStubWithRandSource(src)))
		})
	})

	t.Run("CursorStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) cursorStubInstance {
			return cursorStubInstanceOf(tb, NewCursorStub(tb, CursorStubBenchMode()))
		})
	})

	t.Run("CursorStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) cursorStubInstance {
				return cursorStubInstanceOf(
					tb,
					NewCursorStub(tb, CursorStubWithClock(clk), CursorStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("CursorStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Next with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Next-failing")
			s := NewCursorStub(t, CursorStubFailing(want))
			_, _, r2 := cursorStubNext(t, s)
			assert.ErrorIs(t, r2, want, "Next fails with the error of CursorStubFailing")
		})

		t.Run("fails Close with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Close-failing")
			s := NewCursorStub(t, CursorStubFailing(want))
			r0 := cursorStubClose(t, s)
			assert.ErrorIs(t, r0, want, "Close fails with the error of CursorStubFailing")
		})
	})

	t.Run("WithCursorNext", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Next", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewCursorStub(t, WithCursorNext(func(_ context.Context) (string, bool, error) {
				called = true
				var z0 string
				var z1 bool
				var z2 error
				return z0, z1, z2
			}))
			_, _, _ = cursorStubNext(t, s)
			assert.True(t, called, "WithCursorNext sets the function of Next")
		})
	})

	t.Run("WithCursorClose", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Close", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewCursorStub(t, WithCursorClose(func(_ context.Context) error {
				called = true
				var z0 error
				return z0
			}))
			_ = cursorStubClose(t, s)
			assert.True(t, called, "WithCursorClose sets the function of Close")
		})
	})
}

// TestCursorNextStub checks the methods of CursorNextStub, which
// configures Next.
func TestCursorNextStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Next", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			want0, want1, want2 := cursorStubNextPins()
			s.OnNext.Returns(want0, want1, want2)
			got0, got1, got2 := cursorStubNext(t, s)
			expect.Equal(t, got0, want0, "Next returns the Result0 that Returns sets")
			expect.Equal(t, got1, want1, "Next returns the Result1 that Returns sets")
			expect.ErrorIs(t, got2, want2, "Next returns the error that Returns sets")
		})

		t.Run("configures Next for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewCursorStub(r)
			s.OnNext.Strict()
			var want0 string
			var want1 bool
			var want2 error
			s.OnNext.Returns(want0, want1, want2)
			_, _, _ = cursorStubNext(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Next")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Next to the function", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			var got0 context.Context
			s.OnNext.Func(func(p0 context.Context) (string, bool, error) {
				got0 = p0
				var z0 string
				var z1 bool
				var z2 error
				return z0, z1, z2
			})
			_, _, _ = cursorStubNext(t, s)
			a0 := cursorStubNextArgs(t)
			assert.Equal(t, got0, a0, "the function receives Ctx", assert.ByIdentity())
		})
	})
}

// TestCursorCloseStub checks the methods of CursorCloseStub, which
// configures Close.
func TestCursorCloseStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Close", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			want0 := cursorStubClosePins()
			s.OnClose.Returns(want0)
			got0 := cursorStubClose(t, s)
			assert.ErrorIs(t, got0, want0, "Close returns the error that Returns sets")
		})

		t.Run("configures Close for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewCursorStub(r)
			s.OnClose.Strict()
			var want0 error
			s.OnClose.Returns(want0)
			_ = cursorStubClose(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Close")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Close to the function", func(t *testing.T) {
			t.Parallel()
			s := NewCursorStub(t)
			var got0 context.Context
			s.OnClose.Func(func(p0 context.Context) error {
				got0 = p0
				var z0 error
				return z0
			})
			_ = cursorStubClose(t, s)
			a0 := cursorStubCloseArgs(t)
			assert.Equal(t, got0, a0, "the function receives Ctx", assert.ByIdentity())
		})
	})
}

// readerStubRead calls Read of s with the samples of its arguments.
func readerStubRead(tb assert.TB, s *ReaderStub) Cursor {
	tb.Helper()
	return s.Read()
}

// readerStubReadSubject returns the subject of [stub.Behaviour] for Read.
// The subject calls a new double that reports to tb.
func readerStubReadSubject(tb assert.TB) stub.Subject[ReaderReadCall, ReaderReadReturn] {
	tb.Helper()
	s := NewReaderStub(tb)
	return stub.Subject[ReaderReadCall, ReaderReadReturn]{
		Stub: s.OnRead.MethodStub,
		Call: func() {
			_ = readerStubRead(tb, s)
		},
		Result: func() ReaderReadReturn {
			got0 := readerStubRead(tb, s)
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

// readerStubInstance is the instance of a ReaderStub that the checks of
// the constructors and the options drive.
type readerStubInstance = stub.Instance[ReaderReadCall]

// readerStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Read alone.
func readerStubInstanceOf(tb assert.TB, s *ReaderStub) readerStubInstance {
	tb.Helper()
	return readerStubInstance{
		Stub: s.OnRead.MethodStub,
		Call: func() {
			_ = readerStubRead(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestReaderStub checks the methods, the constructors and the options of
// the double ReaderStub.
func TestReaderStub(t *testing.T) {
	t.Parallel()

	t.Run("Read", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Read", readerStubReadSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewReaderStub(t)
			var seen []ReaderReadCall
			s.OnRead.OnRecord(func(c ReaderReadCall) { seen = append(seen, c) })
			_ = readerStubRead(t, s)
			_ = readerStubRead(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Read")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) readerStubInstance {
			return readerStubInstanceOf(tb, NewReaderStub(tb))
		})

		t.Run("clears the recorded calls of Read", func(t *testing.T) {
			t.Parallel()
			s := NewReaderStub(t)
			_ = readerStubRead(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnRead.CallCount(), 0, "ResetCalls clears the calls of Read")
		})
	})

	t.Run("NewReaderStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) readerStubInstance {
				return readerStubInstanceOf(tb, NewReaderStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("ReaderStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) readerStubInstance {
			return readerStubInstanceOf(tb, NewReaderStub(tb, ReaderStubStrict()))
		})
	})

	t.Run("ReaderStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Read to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewReaderStub(t)
			s := NewReaderStub(t, ReaderStubDelegateTo(inner))
			_ = readerStubRead(t, s)
			inner.OnRead.AssertCalledOnce(t, "Read calls the implementation")
		})
	})

	t.Run("ReaderStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) readerStubInstance {
			return readerStubInstanceOf(tb, NewReaderStub(tb, ReaderStubWithClock(clk)))
		})
	})

	t.Run("ReaderStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) readerStubInstance {
			return readerStubInstanceOf(tb, NewReaderStub(tb, ReaderStubWithRandSource(src)))
		})
	})

	t.Run("ReaderStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) readerStubInstance {
			return readerStubInstanceOf(tb, NewReaderStub(tb, ReaderStubBenchMode()))
		})
	})

	t.Run("ReaderStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) readerStubInstance {
				return readerStubInstanceOf(
					tb,
					NewReaderStub(tb, ReaderStubWithClock(clk), ReaderStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("WithReaderRead", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Read", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewReaderStub(t, WithReaderRead(func() Cursor {
				called = true
				var z0 Cursor
				return z0
			}))
			_ = readerStubRead(t, s)
			assert.True(t, called, "WithReaderRead sets the function of Read")
		})
	})
}

// TestReaderReadStub checks the methods of ReaderReadStub, which
// configures Read.
func TestReaderReadStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("configures Read for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewReaderStub(r)
			s.OnRead.Strict()
			var want0 Cursor
			s.OnRead.Returns(want0)
			_ = readerStubRead(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Read")
		})
	})
}

// subscriptionStubCancel calls Cancel of s with the samples of its
// arguments.
func subscriptionStubCancel(tb assert.TB, s *SubscriptionStub) {
	tb.Helper()
	s.Cancel()
}

// subscriptionStubCancelSubject returns the subject of [stub.Behaviour]
// for Cancel. The subject calls a new double that reports to tb.
func subscriptionStubCancelSubject(tb assert.TB) stub.Subject[SubscriptionCancelCall, struct{}] {
	tb.Helper()
	s := NewSubscriptionStub(tb)
	return stub.Subject[SubscriptionCancelCall, struct{}]{
		Stub: s.OnCancel.MethodStub,
		Call: func() {
			subscriptionStubCancel(tb, s)
		},
		Override: func(mark func()) {
			s.OnCancel.Func(func() {
				mark()
			})
		},
	}
}

// subscriptionStubInstance is the instance of a SubscriptionStub that the
// checks of the constructors and the options drive.
type subscriptionStubInstance = stub.Instance[SubscriptionCancelCall]

// subscriptionStubInstanceOf returns the instance of s. The options apply
// to the whole double, so the checks call Cancel alone.
func subscriptionStubInstanceOf(tb assert.TB, s *SubscriptionStub) subscriptionStubInstance {
	tb.Helper()
	return subscriptionStubInstance{
		Stub: s.OnCancel.MethodStub,
		Call: func() {
			subscriptionStubCancel(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestSubscriptionStub checks the methods, the constructors and the
// options of the double SubscriptionStub.
func TestSubscriptionStub(t *testing.T) {
	t.Parallel()

	t.Run("Cancel", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Cancel", subscriptionStubCancelSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewSubscriptionStub(t)
			var seen []SubscriptionCancelCall
			s.OnCancel.OnRecord(func(c SubscriptionCancelCall) { seen = append(seen, c) })
			subscriptionStubCancel(t, s)
			subscriptionStubCancel(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Cancel")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) subscriptionStubInstance {
			return subscriptionStubInstanceOf(tb, NewSubscriptionStub(tb))
		})

		t.Run("clears the recorded calls of Cancel", func(t *testing.T) {
			t.Parallel()
			s := NewSubscriptionStub(t)
			subscriptionStubCancel(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnCancel.CallCount(), 0, "ResetCalls clears the calls of Cancel")
		})
	})

	t.Run("NewSubscriptionStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) subscriptionStubInstance {
				return subscriptionStubInstanceOf(tb, NewSubscriptionStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("SubscriptionStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) subscriptionStubInstance {
			return subscriptionStubInstanceOf(tb, NewSubscriptionStub(tb, SubscriptionStubStrict()))
		})
	})

	t.Run("SubscriptionStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Cancel to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewSubscriptionStub(t)
			s := NewSubscriptionStub(t, SubscriptionStubDelegateTo(inner))
			subscriptionStubCancel(t, s)
			inner.OnCancel.AssertCalledOnce(t, "Cancel calls the implementation")
		})
	})

	t.Run("SubscriptionStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) subscriptionStubInstance {
			return subscriptionStubInstanceOf(
				tb,
				NewSubscriptionStub(tb, SubscriptionStubWithClock(clk)),
			)
		})
	})

	t.Run("SubscriptionStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) subscriptionStubInstance {
			return subscriptionStubInstanceOf(
				tb,
				NewSubscriptionStub(tb, SubscriptionStubWithRandSource(src)),
			)
		})
	})

	t.Run("SubscriptionStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) subscriptionStubInstance {
			return subscriptionStubInstanceOf(
				tb,
				NewSubscriptionStub(tb, SubscriptionStubBenchMode()),
			)
		})
	})

	t.Run("SubscriptionStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) subscriptionStubInstance {
				return subscriptionStubInstanceOf(
					tb,
					NewSubscriptionStub(
						tb,
						SubscriptionStubWithClock(clk),
						SubscriptionStubWithLatency(d),
					),
				)
			},
		)
	})

	t.Run("WithSubscriptionCancel", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Cancel", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewSubscriptionStub(t, WithSubscriptionCancel(func() {
				called = true
			}))
			subscriptionStubCancel(t, s)
			assert.True(t, called, "WithSubscriptionCancel sets the function of Cancel")
		})
	})
}
