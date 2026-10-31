package edge_test

import (
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/edge"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/stub"
)

// edgeStubSkipArgs returns the samples of the arguments of Skip.
func edgeStubSkipArgs(tb assert.TB) (string, int) {
	tb.Helper()
	var a0 string
	var a1 int
	a0 = "Skip-a"
	a1 = 42
	return a0, a1
}

// edgeStubSkip calls Skip of s with the samples of its arguments.
func edgeStubSkip(tb assert.TB, s *edge.EdgeStub) {
	tb.Helper()
	a0, a1 := edgeStubSkipArgs(tb)
	s.Skip(a0, a1)
}

// edgeStubRead calls Read of s with the samples of its arguments.
func edgeStubRead(tb assert.TB, s *edge.EdgeStub) (int, error) {
	tb.Helper()
	return s.Read()
}

// edgeStubReadPins returns the samples of the results of Read, with a new
// error as its error result.
func edgeStubReadPins() (int, error) {
	var want0 int
	want1 := errors.New("Read-returns")
	want0 = 42
	return want0, want1
}

// edgeStubKeepArgs returns the samples of the arguments of Keep.
func edgeStubKeepArgs(tb assert.TB) int {
	tb.Helper()
	var a0 int
	a0 = 42
	return a0
}

// edgeStubKeep calls Keep of s with the samples of its arguments.
func edgeStubKeep(tb assert.TB, s *edge.EdgeStub) int {
	tb.Helper()
	a0 := edgeStubKeepArgs(tb)
	return s.Keep(a0)
}

// edgeStubKeepPins returns the samples of the results of Keep.
func edgeStubKeepPins() int {
	var want0 int
	want0 = 42
	return want0
}

// edgeStubUseArgs returns the samples of the arguments of Use.
func edgeStubUseArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Use-a"
	return a0
}

// edgeStubUse calls Use of s with the samples of its arguments.
func edgeStubUse(tb assert.TB, s *edge.EdgeStub) error {
	tb.Helper()
	a0 := edgeStubUseArgs(tb)
	return s.Use(a0)
}

// edgeStubUsePins returns the samples of the results of Use, with a new
// error as its error result.
func edgeStubUsePins() error {
	want0 := errors.New("Use-returns")
	return want0
}

// edgeStubRecordArgs returns the samples of the arguments of Record.
func edgeStubRecordArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Record-a"
	return a0
}

// edgeStubRecord calls Record of s with the samples of its arguments.
func edgeStubRecord(tb assert.TB, s *edge.EdgeStub) (int, error) {
	tb.Helper()
	a0 := edgeStubRecordArgs(tb)
	return s.Record(a0)
}

// edgeStubRecordPins returns the samples of the results of Record, with a
// new error as its error result.
func edgeStubRecordPins() (int, error) {
	var want0 int
	want1 := errors.New("Record-returns")
	want0 = 42
	return want0, want1
}

// edgeStubClash calls Clash of s with the samples of its arguments.
func edgeStubClash(tb assert.TB, s *edge.EdgeStub) (string, error) {
	tb.Helper()
	return s.Clash()
}

// edgeStubClashPins returns the samples of the results of Clash, with a
// new error as its error result.
func edgeStubClashPins() (string, error) {
	var want0 string
	want1 := errors.New("Clash-returns")
	want0 = "Clash-a"
	return want0, want1
}

// edgeStubShiftArgs returns the samples of the arguments of Shift.
func edgeStubShiftArgs(tb assert.TB) (string, int) {
	tb.Helper()
	var a0 string
	var a1 int
	a0 = "Shift-a"
	a1 = 42
	return a0, a1
}

// edgeStubShift calls Shift of s with the samples of its arguments.
func edgeStubShift(tb assert.TB, s *edge.EdgeStub) {
	tb.Helper()
	a0, a1 := edgeStubShiftArgs(tb)
	s.Shift(a0, a1)
}

// edgeStubFetchArgs returns the samples of the arguments of Fetch.
func edgeStubFetchArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Fetch-a"
	return a0
}

// edgeStubFetch calls Fetch of s with the samples of its arguments.
func edgeStubFetch(tb assert.TB, s *edge.EdgeStub) (int, error) {
	tb.Helper()
	a0 := edgeStubFetchArgs(tb)
	return s.Fetch(a0)
}

// edgeStubFetchPins returns the samples of the results of Fetch, with a
// new error as its error result.
func edgeStubFetchPins() (int, error) {
	var want0 int
	want1 := errors.New("Fetch-returns")
	want0 = 42
	return want0, want1
}

// edgeStubWideArgs returns the samples of the arguments of Wide.
func edgeStubWideArgs(tb assert.TB) (int, int, int, int, int, int, int, int, int, int, int) {
	tb.Helper()
	var a0 int
	var a1 int
	var a2 int
	var a3 int
	var a4 int
	var a5 int
	var a6 int
	var a7 int
	var a8 int
	var a9 int
	var a10 int
	a0 = 42
	a1 = 42
	a2 = 42
	a3 = 42
	a4 = 42
	a5 = 42
	a6 = 42
	a7 = 42
	a8 = 42
	a9 = 42
	a10 = 42
	return a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10
}

// edgeStubWide calls Wide of s with the samples of its arguments.
func edgeStubWide(tb assert.TB, s *edge.EdgeStub) {
	tb.Helper()
	a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 := edgeStubWideArgs(tb)
	s.Wide(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
}

// edgeStubPair calls Pair of s with the samples of its arguments.
func edgeStubPair(tb assert.TB, s *edge.EdgeStub) (int, int) {
	tb.Helper()
	return s.Pair()
}

// edgeStubPairPins returns the samples of the results of Pair.
func edgeStubPairPins() (int, int) {
	var want0 int
	var want1 int
	want0 = 42
	want1 = 42
	return want0, want1
}

// edgeStubSkipSubject returns the subject of [stub.Behaviour] for Skip.
// The subject calls a new double that reports to tb.
func edgeStubSkipSubject(tb assert.TB) stub.Subject[edge.EdgeSkipCall, struct{}] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeSkipCall, struct{}]{
		Stub: s.OnSkip.MethodStub,
		Call: func() {
			edgeStubSkip(tb, s)
		},
		Override: func(mark func()) {
			s.OnSkip.Func(func(_ string, _ int) {
				mark()
			})
		},
	}
}

// edgeStubReadSubject returns the subject of [stub.Behaviour] for Read.
// The subject calls a new double that reports to tb.
func edgeStubReadSubject(tb assert.TB) stub.Subject[edge.EdgeReadCall, edge.EdgeReadReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeReadCall, edge.EdgeReadReturn]{
		Stub: s.OnRead.MethodStub,
		Call: func() {
			_, _ = edgeStubRead(tb, s)
		},
		Result: func() edge.EdgeReadReturn {
			got0, got1 := edgeStubRead(tb, s)
			return edge.EdgeReadReturn{R: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnRead.Func(func() (int, error) {
				mark()
				var z0 int
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := edgeStubRead(tb, s)
			return r1
		},
	}
}

// edgeStubKeepSubject returns the subject of [stub.Behaviour] for Keep.
// The subject calls a new double that reports to tb.
func edgeStubKeepSubject(tb assert.TB) stub.Subject[edge.EdgeKeepCall, edge.EdgeKeepReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeKeepCall, edge.EdgeKeepReturn]{
		Stub: s.OnKeep.MethodStub,
		Call: func() {
			_ = edgeStubKeep(tb, s)
		},
		Result: func() edge.EdgeKeepReturn {
			got0 := edgeStubKeep(tb, s)
			return edge.EdgeKeepReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnKeep.Func(func(_ int) int {
				mark()
				var z0 int
				return z0
			})
		},
	}
}

// edgeStubUseSubject returns the subject of [stub.Behaviour] for Use. The
// subject calls a new double that reports to tb.
func edgeStubUseSubject(tb assert.TB) stub.Subject[edge.EdgeUseCall, edge.EdgeUseReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeUseCall, edge.EdgeUseReturn]{
		Stub: s.OnUse.MethodStub,
		Call: func() {
			_ = edgeStubUse(tb, s)
		},
		Result: func() edge.EdgeUseReturn {
			got0 := edgeStubUse(tb, s)
			return edge.EdgeUseReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnUse.Func(func(_ string) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			r0 := edgeStubUse(tb, s)
			return r0
		},
	}
}

// edgeStubRecordSubject returns the subject of [stub.Behaviour] for
// Record. The subject calls a new double that reports to tb.
func edgeStubRecordSubject(tb assert.TB) stub.Subject[edge.EdgeRecordCall, edge.EdgeRecordReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeRecordCall, edge.EdgeRecordReturn]{
		Stub: s.OnRecord.MethodStub,
		Call: func() {
			_, _ = edgeStubRecord(tb, s)
		},
		Result: func() edge.EdgeRecordReturn {
			got0, got1 := edgeStubRecord(tb, s)
			return edge.EdgeRecordReturn{N: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnRecord.Func(func(_ string) (int, error) {
				mark()
				var z0 int
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, err := edgeStubRecord(tb, s)
			return err
		},
	}
}

// edgeStubClashSubject returns the subject of [stub.Behaviour] for Clash.
// The subject calls a new double that reports to tb.
func edgeStubClashSubject(tb assert.TB) stub.Subject[edge.EdgeClashCall, edge.EdgeClashReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeClashCall, edge.EdgeClashReturn]{
		Stub: s.OnClash.MethodStub,
		Call: func() {
			_, _ = edgeStubClash(tb, s)
		},
		Result: func() edge.EdgeClashReturn {
			got0, got1 := edgeStubClash(tb, s)
			return edge.EdgeClashReturn{S: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnClash.Func(func() (string, error) {
				mark()
				var z0 string
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := edgeStubClash(tb, s)
			return r1
		},
	}
}

// edgeStubShiftSubject returns the subject of [stub.Behaviour] for Shift.
// The subject calls a new double that reports to tb.
func edgeStubShiftSubject(tb assert.TB) stub.Subject[edge.EdgeShiftCall, struct{}] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeShiftCall, struct{}]{
		Stub: s.OnShift.MethodStub,
		Call: func() {
			edgeStubShift(tb, s)
		},
		Override: func(mark func()) {
			s.OnShift.Func(func(_ string, _ int) {
				mark()
			})
		},
	}
}

// edgeStubFetchSubject returns the subject of [stub.Behaviour] for Fetch.
// The subject calls a new double that reports to tb.
func edgeStubFetchSubject(tb assert.TB) stub.Subject[edge.EdgeFetchCall, edge.EdgeFetchReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeFetchCall, edge.EdgeFetchReturn]{
		Stub: s.OnFetch.MethodStub,
		Call: func() {
			_, _ = edgeStubFetch(tb, s)
		},
		Result: func() edge.EdgeFetchReturn {
			got0, got1 := edgeStubFetch(tb, s)
			return edge.EdgeFetchReturn{A0: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnFetch.Func(func(_ string) (int, error) {
				mark()
				var z0 int
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, err := edgeStubFetch(tb, s)
			return err
		},
	}
}

// edgeStubWideSubject returns the subject of [stub.Behaviour] for Wide.
// The subject calls a new double that reports to tb.
func edgeStubWideSubject(tb assert.TB) stub.Subject[edge.EdgeWideCall, struct{}] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeWideCall, struct{}]{
		Stub: s.OnWide.MethodStub,
		Call: func() {
			edgeStubWide(tb, s)
		},
		Override: func(mark func()) {
			s.OnWide.Func(func(_, _, _, _, _, _, _, _, _, _, _ int) {
				mark()
			})
		},
	}
}

// edgeStubPairSubject returns the subject of [stub.Behaviour] for Pair.
// The subject calls a new double that reports to tb.
func edgeStubPairSubject(tb assert.TB) stub.Subject[edge.EdgePairCall, edge.EdgePairReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgePairCall, edge.EdgePairReturn]{
		Stub: s.OnPair.MethodStub,
		Call: func() {
			_, _ = edgeStubPair(tb, s)
		},
		Result: func() edge.EdgePairReturn {
			got0, got1 := edgeStubPair(tb, s)
			return edge.EdgePairReturn{AB: got0, Ab: got1}
		},
		Override: func(mark func()) {
			s.OnPair.Func(func() (int, int) {
				mark()
				var z0 int
				var z1 int
				return z0, z1
			})
		},
	}
}

// edgeStubInstance is the instance of a edge.EdgeStub that the checks of
// the constructors and the options drive.
type edgeStubInstance = stub.Instance[edge.EdgeSkipCall]

// edgeStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Skip alone.
func edgeStubInstanceOf(tb assert.TB, s *edge.EdgeStub) edgeStubInstance {
	tb.Helper()
	return edgeStubInstance{
		Stub: s.OnSkip.MethodStub,
		Call: func() {
			edgeStubSkip(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestEdgeStub checks the methods, the constructors and the options of the
// double EdgeStub.
func TestEdgeStub(t *testing.T) {
	t.Parallel()

	t.Run("Skip", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Skip", edgeStubSkipSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			edgeStubSkip(t, s)
			got := s.OnSkip.AssertCalledOnce(t, "Skip records the call")
			a0, a1 := edgeStubSkipArgs(t)
			expect.Equal(t, got.A0, a0, "the call records A0")
			expect.Equal(t, got.E, a1, "the call records E")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeSkipCall
			s.OnSkip.OnRecord(func(c edge.EdgeSkipCall) { seen = append(seen, c) })
			edgeStubSkip(t, s)
			edgeStubSkip(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Skip")
		})
	})

	t.Run("Read", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Read", edgeStubReadSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeReadCall
			s.OnRead.OnRecord(func(c edge.EdgeReadCall) { seen = append(seen, c) })
			_, _ = edgeStubRead(t, s)
			_, _ = edgeStubRead(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Read")
		})
	})

	t.Run("Keep", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Keep", edgeStubKeepSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_ = edgeStubKeep(t, s)
			got := s.OnKeep.AssertCalledOnce(t, "Keep records the call")
			a0 := edgeStubKeepArgs(t)
			assert.Equal(t, got.R, a0, "the call records R")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeKeepCall
			s.OnKeep.OnRecord(func(c edge.EdgeKeepCall) { seen = append(seen, c) })
			_ = edgeStubKeep(t, s)
			_ = edgeStubKeep(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Keep")
		})
	})

	t.Run("Use", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Use", edgeStubUseSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_ = edgeStubUse(t, s)
			got := s.OnUse.AssertCalledOnce(t, "Use records the call")
			a0 := edgeStubUseArgs(t)
			assert.Equal(t, got.S, a0, "the call records S")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeUseCall
			s.OnUse.OnRecord(func(c edge.EdgeUseCall) { seen = append(seen, c) })
			_ = edgeStubUse(t, s)
			_ = edgeStubUse(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Use")
		})
	})

	t.Run("Record", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Record", edgeStubRecordSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_, _ = edgeStubRecord(t, s)
			got := s.OnRecord.AssertCalledOnce(t, "Record records the call")
			a0 := edgeStubRecordArgs(t)
			assert.Equal(t, got.Call, a0, "the call records Call")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeRecordCall
			s.OnRecord.OnRecord(func(c edge.EdgeRecordCall) { seen = append(seen, c) })
			_, _ = edgeStubRecord(t, s)
			_, _ = edgeStubRecord(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Record")
		})
	})

	t.Run("Clash", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Clash", edgeStubClashSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeClashCall
			s.OnClash.OnRecord(func(c edge.EdgeClashCall) { seen = append(seen, c) })
			_, _ = edgeStubClash(t, s)
			_, _ = edgeStubClash(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Clash")
		})
	})

	t.Run("Shift", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Shift", edgeStubShiftSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			edgeStubShift(t, s)
			got := s.OnShift.AssertCalledOnce(t, "Shift records the call")
			a0, a1 := edgeStubShiftArgs(t)
			expect.Equal(t, got.A00, a0, "the call records A00")
			expect.Equal(t, got.A0, a1, "the call records A0")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeShiftCall
			s.OnShift.OnRecord(func(c edge.EdgeShiftCall) { seen = append(seen, c) })
			edgeStubShift(t, s)
			edgeStubShift(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Shift")
		})
	})

	t.Run("Fetch", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Fetch", edgeStubFetchSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_, _ = edgeStubFetch(t, s)
			got := s.OnFetch.AssertCalledOnce(t, "Fetch records the call")
			a0 := edgeStubFetchArgs(t)
			assert.Equal(t, got.A00, a0, "the call records A00")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeFetchCall
			s.OnFetch.OnRecord(func(c edge.EdgeFetchCall) { seen = append(seen, c) })
			_, _ = edgeStubFetch(t, s)
			_, _ = edgeStubFetch(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Fetch")
		})
	})

	t.Run("Wide", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Wide", edgeStubWideSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			edgeStubWide(t, s)
			got := s.OnWide.AssertCalledOnce(t, "Wide records the call")
			a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 := edgeStubWideArgs(t)
			expect.Equal(t, got.A, a0, "the call records A")
			expect.Equal(t, got.A10, a1, "the call records A10")
			expect.Equal(t, got.A1, a2, "the call records A1")
			expect.Equal(t, got.B, a3, "the call records B")
			expect.Equal(t, got.C, a4, "the call records C")
			expect.Equal(t, got.D, a5, "the call records D")
			expect.Equal(t, got.E, a6, "the call records E")
			expect.Equal(t, got.F, a7, "the call records F")
			expect.Equal(t, got.G, a8, "the call records G")
			expect.Equal(t, got.H, a9, "the call records H")
			expect.Equal(t, got.A100, a10, "the call records A100")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgeWideCall
			s.OnWide.OnRecord(func(c edge.EdgeWideCall) { seen = append(seen, c) })
			edgeStubWide(t, s)
			edgeStubWide(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Wide")
		})
	})

	t.Run("Pair", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Pair", edgeStubPairSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var seen []edge.EdgePairCall
			s.OnPair.OnRecord(func(c edge.EdgePairCall) { seen = append(seen, c) })
			_, _ = edgeStubPair(t, s)
			_, _ = edgeStubPair(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Pair")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) edgeStubInstance {
			return edgeStubInstanceOf(tb, edge.NewEdgeStub(tb))
		})

		t.Run("clears the recorded calls of Skip", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			edgeStubSkip(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnSkip.CallCount(), 0, "ResetCalls clears the calls of Skip")
		})

		t.Run("clears the recorded calls of Read", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_, _ = edgeStubRead(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnRead.CallCount(), 0, "ResetCalls clears the calls of Read")
		})

		t.Run("keeps the results of Returns of Read", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubReadPins()
			s.OnRead.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := edgeStubRead(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the R of Read")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Read")
		})

		t.Run("clears the recorded calls of Keep", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_ = edgeStubKeep(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnKeep.CallCount(), 0, "ResetCalls clears the calls of Keep")
		})

		t.Run("keeps the results of Returns of Keep", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0 := edgeStubKeepPins()
			s.OnKeep.Returns(want0)
			s.ResetCalls()
			got0 := edgeStubKeep(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Keep")
		})

		t.Run("clears the recorded calls of Use", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_ = edgeStubUse(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnUse.CallCount(), 0, "ResetCalls clears the calls of Use")
		})

		t.Run("keeps the results of Returns of Use", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0 := edgeStubUsePins()
			s.OnUse.Returns(want0)
			s.ResetCalls()
			got0 := edgeStubUse(t, s)
			assert.ErrorIs(t, got0, want0, "ResetCalls keeps the error of Use")
		})

		t.Run("clears the recorded calls of Record", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_, _ = edgeStubRecord(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnRecord.CallCount(), 0, "ResetCalls clears the calls of Record")
		})

		t.Run("keeps the results of Returns of Record", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubRecordPins()
			s.OnRecord.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := edgeStubRecord(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the N of Record")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Record")
		})

		t.Run("clears the recorded calls of Clash", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_, _ = edgeStubClash(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnClash.CallCount(), 0, "ResetCalls clears the calls of Clash")
		})

		t.Run("keeps the results of Returns of Clash", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubClashPins()
			s.OnClash.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := edgeStubClash(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the S of Clash")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Clash")
		})

		t.Run("clears the recorded calls of Shift", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			edgeStubShift(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnShift.CallCount(), 0, "ResetCalls clears the calls of Shift")
		})

		t.Run("clears the recorded calls of Fetch", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_, _ = edgeStubFetch(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnFetch.CallCount(), 0, "ResetCalls clears the calls of Fetch")
		})

		t.Run("keeps the results of Returns of Fetch", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubFetchPins()
			s.OnFetch.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := edgeStubFetch(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the A0 of Fetch")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Fetch")
		})

		t.Run("clears the recorded calls of Wide", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			edgeStubWide(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnWide.CallCount(), 0, "ResetCalls clears the calls of Wide")
		})

		t.Run("clears the recorded calls of Pair", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			_, _ = edgeStubPair(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnPair.CallCount(), 0, "ResetCalls clears the calls of Pair")
		})

		t.Run("keeps the results of Returns of Pair", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubPairPins()
			s.OnPair.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := edgeStubPair(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the AB of Pair")
			expect.Equal(t, got1, want1, "ResetCalls keeps the Ab of Pair")
		})
	})

	t.Run("NewEdgeStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) edgeStubInstance {
				return edgeStubInstanceOf(tb, edge.NewEdgeStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewEdgeStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Read with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Read-failing")
			s := edge.NewEdgeStubFailing(t, want)
			_, r1 := edgeStubRead(t, s)
			assert.ErrorIs(t, r1, want, "Read fails with the error of NewEdgeStubFailing")
		})

		t.Run("fails Use with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Use-failing")
			s := edge.NewEdgeStubFailing(t, want)
			r0 := edgeStubUse(t, s)
			assert.ErrorIs(t, r0, want, "Use fails with the error of NewEdgeStubFailing")
		})

		t.Run("fails Record with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Record-failing")
			s := edge.NewEdgeStubFailing(t, want)
			_, err := edgeStubRecord(t, s)
			assert.ErrorIs(t, err, want, "Record fails with the error of NewEdgeStubFailing")
		})

		t.Run("fails Clash with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Clash-failing")
			s := edge.NewEdgeStubFailing(t, want)
			_, r1 := edgeStubClash(t, s)
			assert.ErrorIs(t, r1, want, "Clash fails with the error of NewEdgeStubFailing")
		})

		t.Run("fails Fetch with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Fetch-failing")
			s := edge.NewEdgeStubFailing(t, want)
			_, err := edgeStubFetch(t, s)
			assert.ErrorIs(t, err, want, "Fetch fails with the error of NewEdgeStubFailing")
		})
	})

	t.Run("EdgeStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) edgeStubInstance {
			return edgeStubInstanceOf(tb, edge.NewEdgeStub(tb, edge.EdgeStubStrict()))
		})
	})

	t.Run("EdgeStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Skip to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			edgeStubSkip(t, s)
			inner.OnSkip.AssertCalledOnce(t, "Skip calls the implementation")
		})

		t.Run("passes Read to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, _ = edgeStubRead(t, s)
			inner.OnRead.AssertCalledOnce(t, "Read calls the implementation")
		})

		t.Run("fails Read with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Read-delegate")
			inner := edge.NewEdgeStub(t)
			inner.OnRead.Faults(want, 1)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, r1 := edgeStubRead(t, s)
			assert.ErrorIs(t, r1, want, "Read fails with the error of the implementation")
		})

		t.Run("passes Keep to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_ = edgeStubKeep(t, s)
			inner.OnKeep.AssertCalledOnce(t, "Keep calls the implementation")
		})

		t.Run("passes Use to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_ = edgeStubUse(t, s)
			inner.OnUse.AssertCalledOnce(t, "Use calls the implementation")
		})

		t.Run("fails Use with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Use-delegate")
			inner := edge.NewEdgeStub(t)
			inner.OnUse.Faults(want, 1)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			r0 := edgeStubUse(t, s)
			assert.ErrorIs(t, r0, want, "Use fails with the error of the implementation")
		})

		t.Run("passes Record to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, _ = edgeStubRecord(t, s)
			inner.OnRecord.AssertCalledOnce(t, "Record calls the implementation")
		})

		t.Run("fails Record with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Record-delegate")
			inner := edge.NewEdgeStub(t)
			inner.OnRecord.Faults(want, 1)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, err := edgeStubRecord(t, s)
			assert.ErrorIs(t, err, want, "Record fails with the error of the implementation")
		})

		t.Run("passes Clash to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, _ = edgeStubClash(t, s)
			inner.OnClash.AssertCalledOnce(t, "Clash calls the implementation")
		})

		t.Run("fails Clash with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Clash-delegate")
			inner := edge.NewEdgeStub(t)
			inner.OnClash.Faults(want, 1)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, r1 := edgeStubClash(t, s)
			assert.ErrorIs(t, r1, want, "Clash fails with the error of the implementation")
		})

		t.Run("passes Shift to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			edgeStubShift(t, s)
			inner.OnShift.AssertCalledOnce(t, "Shift calls the implementation")
		})

		t.Run("passes Fetch to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, _ = edgeStubFetch(t, s)
			inner.OnFetch.AssertCalledOnce(t, "Fetch calls the implementation")
		})

		t.Run("fails Fetch with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Fetch-delegate")
			inner := edge.NewEdgeStub(t)
			inner.OnFetch.Faults(want, 1)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, err := edgeStubFetch(t, s)
			assert.ErrorIs(t, err, want, "Fetch fails with the error of the implementation")
		})

		t.Run("passes Wide to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			edgeStubWide(t, s)
			inner.OnWide.AssertCalledOnce(t, "Wide calls the implementation")
		})

		t.Run("passes Pair to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := edge.NewEdgeStub(t)
			s := edge.NewEdgeStub(t, edge.EdgeStubDelegateTo(inner))
			_, _ = edgeStubPair(t, s)
			inner.OnPair.AssertCalledOnce(t, "Pair calls the implementation")
		})
	})

	t.Run("EdgeStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) edgeStubInstance {
			return edgeStubInstanceOf(tb, edge.NewEdgeStub(tb, edge.EdgeStubWithClock(clk)))
		})
	})

	t.Run("EdgeStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) edgeStubInstance {
			return edgeStubInstanceOf(tb, edge.NewEdgeStub(tb, edge.EdgeStubWithRandSource(src)))
		})
	})

	t.Run("EdgeStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) edgeStubInstance {
			return edgeStubInstanceOf(tb, edge.NewEdgeStub(tb, edge.EdgeStubBenchMode()))
		})
	})

	t.Run("EdgeStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) edgeStubInstance {
				return edgeStubInstanceOf(
					tb,
					edge.NewEdgeStub(tb, edge.EdgeStubWithClock(clk), edge.EdgeStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("EdgeStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Read with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Read-failing")
			s := edge.NewEdgeStub(t, edge.EdgeStubFailing(want))
			_, r1 := edgeStubRead(t, s)
			assert.ErrorIs(t, r1, want, "Read fails with the error of EdgeStubFailing")
		})

		t.Run("fails Use with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Use-failing")
			s := edge.NewEdgeStub(t, edge.EdgeStubFailing(want))
			r0 := edgeStubUse(t, s)
			assert.ErrorIs(t, r0, want, "Use fails with the error of EdgeStubFailing")
		})

		t.Run("fails Record with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Record-failing")
			s := edge.NewEdgeStub(t, edge.EdgeStubFailing(want))
			_, err := edgeStubRecord(t, s)
			assert.ErrorIs(t, err, want, "Record fails with the error of EdgeStubFailing")
		})

		t.Run("fails Clash with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Clash-failing")
			s := edge.NewEdgeStub(t, edge.EdgeStubFailing(want))
			_, r1 := edgeStubClash(t, s)
			assert.ErrorIs(t, r1, want, "Clash fails with the error of EdgeStubFailing")
		})

		t.Run("fails Fetch with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Fetch-failing")
			s := edge.NewEdgeStub(t, edge.EdgeStubFailing(want))
			_, err := edgeStubFetch(t, s)
			assert.ErrorIs(t, err, want, "Fetch fails with the error of EdgeStubFailing")
		})
	})

	t.Run("WithEdgeSkip", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Skip", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeSkip(func(_ string, _ int) {
				called = true
			}))
			edgeStubSkip(t, s)
			assert.True(t, called, "WithEdgeSkip sets the function of Skip")
		})
	})

	t.Run("WithEdgeRead", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Read", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeRead(func() (int, error) {
				called = true
				var z0 int
				var z1 error
				return z0, z1
			}))
			_, _ = edgeStubRead(t, s)
			assert.True(t, called, "WithEdgeRead sets the function of Read")
		})
	})

	t.Run("WithEdgeKeep", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Keep", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeKeep(func(_ int) int {
				called = true
				var z0 int
				return z0
			}))
			_ = edgeStubKeep(t, s)
			assert.True(t, called, "WithEdgeKeep sets the function of Keep")
		})
	})

	t.Run("WithEdgeUse", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Use", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeUse(func(_ string) error {
				called = true
				var z0 error
				return z0
			}))
			_ = edgeStubUse(t, s)
			assert.True(t, called, "WithEdgeUse sets the function of Use")
		})
	})

	t.Run("WithEdgeRecord", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Record", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeRecord(func(_ string) (int, error) {
				called = true
				var z0 int
				var z1 error
				return z0, z1
			}))
			_, _ = edgeStubRecord(t, s)
			assert.True(t, called, "WithEdgeRecord sets the function of Record")
		})
	})

	t.Run("WithEdgeClash", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Clash", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeClash(func() (string, error) {
				called = true
				var z0 string
				var z1 error
				return z0, z1
			}))
			_, _ = edgeStubClash(t, s)
			assert.True(t, called, "WithEdgeClash sets the function of Clash")
		})
	})

	t.Run("WithEdgeShift", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Shift", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeShift(func(_ string, _ int) {
				called = true
			}))
			edgeStubShift(t, s)
			assert.True(t, called, "WithEdgeShift sets the function of Shift")
		})
	})

	t.Run("WithEdgeFetch", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Fetch", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeFetch(func(_ string) (int, error) {
				called = true
				var z0 int
				var z1 error
				return z0, z1
			}))
			_, _ = edgeStubFetch(t, s)
			assert.True(t, called, "WithEdgeFetch sets the function of Fetch")
		})
	})

	t.Run("WithEdgeWide", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Wide", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgeWide(func(_, _, _, _, _, _, _, _, _, _, _ int) {
				called = true
			}))
			edgeStubWide(t, s)
			assert.True(t, called, "WithEdgeWide sets the function of Wide")
		})
	})

	t.Run("WithEdgePair", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Pair", func(t *testing.T) {
			t.Parallel()
			called := false
			s := edge.NewEdgeStub(t, edge.WithEdgePair(func() (int, int) {
				called = true
				var z0 int
				var z1 int
				return z0, z1
			}))
			_, _ = edgeStubPair(t, s)
			assert.True(t, called, "WithEdgePair sets the function of Pair")
		})
	})
}

// TestEdgeSkipStub checks the methods of EdgeSkipStub, which configures
// Skip.
func TestEdgeSkipStub(t *testing.T) {
	t.Parallel()

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Skip to the function", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var got0 string
			var got1 int
			s.OnSkip.Func(func(p0 string, p1 int) {
				got0, got1 = p0, p1
			})
			edgeStubSkip(t, s)
			a0, a1 := edgeStubSkipArgs(t)
			expect.Equal(t, got0, a0, "the function receives A0")
			expect.Equal(t, got1, a1, "the function receives E")
		})
	})
}

// TestEdgeReadStub checks the methods of EdgeReadStub, which configures
// Read.
func TestEdgeReadStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Read", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubReadPins()
			s.OnRead.Returns(want0, want1)
			got0, got1 := edgeStubRead(t, s)
			expect.Equal(t, got0, want0, "Read returns the R that Returns sets")
			expect.ErrorIs(t, got1, want1, "Read returns the error that Returns sets")
		})

		t.Run("configures Read for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := edge.NewEdgeStub(r)
			s.OnRead.Strict()
			var want0 int
			var want1 error
			s.OnRead.Returns(want0, want1)
			_, _ = edgeStubRead(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Read")
		})
	})
}

// TestEdgeKeepStub checks the methods of EdgeKeepStub, which configures
// Keep.
func TestEdgeKeepStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Keep", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0 := edgeStubKeepPins()
			s.OnKeep.Returns(want0)
			got0 := edgeStubKeep(t, s)
			assert.Equal(t, got0, want0, "Keep returns the Result that Returns sets")
		})

		t.Run("configures Keep for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := edge.NewEdgeStub(r)
			s.OnKeep.Strict()
			var want0 int
			s.OnKeep.Returns(want0)
			_ = edgeStubKeep(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Keep")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Keep to the function", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var got0 int
			s.OnKeep.Func(func(p0 int) int {
				got0 = p0
				var z0 int
				return z0
			})
			_ = edgeStubKeep(t, s)
			a0 := edgeStubKeepArgs(t)
			assert.Equal(t, got0, a0, "the function receives R")
		})
	})
}

// TestEdgeUseStub checks the methods of EdgeUseStub, which configures Use.
func TestEdgeUseStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Use", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0 := edgeStubUsePins()
			s.OnUse.Returns(want0)
			got0 := edgeStubUse(t, s)
			assert.ErrorIs(t, got0, want0, "Use returns the error that Returns sets")
		})

		t.Run("configures Use for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := edge.NewEdgeStub(r)
			s.OnUse.Strict()
			var want0 error
			s.OnUse.Returns(want0)
			_ = edgeStubUse(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Use")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Use to the function", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var got0 string
			s.OnUse.Func(func(p0 string) error {
				got0 = p0
				var z0 error
				return z0
			})
			_ = edgeStubUse(t, s)
			a0 := edgeStubUseArgs(t)
			assert.Equal(t, got0, a0, "the function receives S")
		})
	})
}

// TestEdgeRecordStub checks the methods of EdgeRecordStub, which
// configures Record.
func TestEdgeRecordStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Record", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubRecordPins()
			s.OnRecord.Returns(want0, want1)
			got0, got1 := edgeStubRecord(t, s)
			expect.Equal(t, got0, want0, "Record returns the N that Returns sets")
			expect.ErrorIs(t, got1, want1, "Record returns the error that Returns sets")
		})

		t.Run("configures Record for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := edge.NewEdgeStub(r)
			s.OnRecord.Strict()
			var want0 int
			var want1 error
			s.OnRecord.Returns(want0, want1)
			_, _ = edgeStubRecord(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Record")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Record to the function", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var got0 string
			s.OnRecord.Func(func(p0 string) (int, error) {
				got0 = p0
				var z0 int
				var z1 error
				return z0, z1
			})
			_, _ = edgeStubRecord(t, s)
			a0 := edgeStubRecordArgs(t)
			assert.Equal(t, got0, a0, "the function receives Call")
		})
	})
}

// TestEdgeClashStub checks the methods of EdgeClashStub, which configures
// Clash.
func TestEdgeClashStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Clash", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubClashPins()
			s.OnClash.Returns(want0, want1)
			got0, got1 := edgeStubClash(t, s)
			expect.Equal(t, got0, want0, "Clash returns the S that Returns sets")
			expect.ErrorIs(t, got1, want1, "Clash returns the error that Returns sets")
		})

		t.Run("configures Clash for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := edge.NewEdgeStub(r)
			s.OnClash.Strict()
			var want0 string
			var want1 error
			s.OnClash.Returns(want0, want1)
			_, _ = edgeStubClash(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Clash")
		})
	})
}

// TestEdgeShiftStub checks the methods of EdgeShiftStub, which configures
// Shift.
func TestEdgeShiftStub(t *testing.T) {
	t.Parallel()

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Shift to the function", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var got0 string
			var got1 int
			s.OnShift.Func(func(p0 string, p1 int) {
				got0, got1 = p0, p1
			})
			edgeStubShift(t, s)
			a0, a1 := edgeStubShiftArgs(t)
			expect.Equal(t, got0, a0, "the function receives A00")
			expect.Equal(t, got1, a1, "the function receives A0")
		})
	})
}

// TestEdgeFetchStub checks the methods of EdgeFetchStub, which configures
// Fetch.
func TestEdgeFetchStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Fetch", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubFetchPins()
			s.OnFetch.Returns(want0, want1)
			got0, got1 := edgeStubFetch(t, s)
			expect.Equal(t, got0, want0, "Fetch returns the A0 that Returns sets")
			expect.ErrorIs(t, got1, want1, "Fetch returns the error that Returns sets")
		})

		t.Run("configures Fetch for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := edge.NewEdgeStub(r)
			s.OnFetch.Strict()
			var want0 int
			var want1 error
			s.OnFetch.Returns(want0, want1)
			_, _ = edgeStubFetch(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Fetch")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Fetch to the function", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var got0 string
			s.OnFetch.Func(func(p0 string) (int, error) {
				got0 = p0
				var z0 int
				var z1 error
				return z0, z1
			})
			_, _ = edgeStubFetch(t, s)
			a0 := edgeStubFetchArgs(t)
			assert.Equal(t, got0, a0, "the function receives A00")
		})
	})
}

// TestEdgeWideStub checks the methods of EdgeWideStub, which configures
// Wide.
func TestEdgeWideStub(t *testing.T) {
	t.Parallel()

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Wide to the function", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			var got0 int
			var got1 int
			var got2 int
			var got3 int
			var got4 int
			var got5 int
			var got6 int
			var got7 int
			var got8 int
			var got9 int
			var got10 int
			s.OnWide.Func(func(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10 int) {
				got0, got1, got2, got3, got4, got5, got6, got7, got8, got9, got10 = p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10
			})
			edgeStubWide(t, s)
			a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 := edgeStubWideArgs(t)
			expect.Equal(t, got0, a0, "the function receives A")
			expect.Equal(t, got1, a1, "the function receives A10")
			expect.Equal(t, got2, a2, "the function receives A1")
			expect.Equal(t, got3, a3, "the function receives B")
			expect.Equal(t, got4, a4, "the function receives C")
			expect.Equal(t, got5, a5, "the function receives D")
			expect.Equal(t, got6, a6, "the function receives E")
			expect.Equal(t, got7, a7, "the function receives F")
			expect.Equal(t, got8, a8, "the function receives G")
			expect.Equal(t, got9, a9, "the function receives H")
			expect.Equal(t, got10, a10, "the function receives A100")
		})
	})
}

// TestEdgePairStub checks the methods of EdgePairStub, which configures
// Pair.
func TestEdgePairStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Pair", func(t *testing.T) {
			t.Parallel()
			s := edge.NewEdgeStub(t)
			want0, want1 := edgeStubPairPins()
			s.OnPair.Returns(want0, want1)
			got0, got1 := edgeStubPair(t, s)
			expect.Equal(t, got0, want0, "Pair returns the AB that Returns sets")
			expect.Equal(t, got1, want1, "Pair returns the Ab that Returns sets")
		})

		t.Run("configures Pair for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := edge.NewEdgeStub(r)
			s.OnPair.Strict()
			var want0 int
			var want1 int
			s.OnPair.Returns(want0, want1)
			_, _ = edgeStubPair(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Pair")
		})
	})
}
