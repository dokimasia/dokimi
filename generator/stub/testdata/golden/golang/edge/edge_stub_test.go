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

// edgeStubSkipSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Skip. The subject calls a new double that
// reports to tb.
func edgeStubSkipSubject(tb assert.TB) stub.Subject[edge.EdgeSkipCall, struct{}] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeSkipCall, struct{}]{
		Stub:  s.OnSkip.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			edgeStubSkip(tb, s)
		},
		Arguments: func(tb assert.TB, c edge.EdgeSkipCall) {
			tb.Helper()
			a0, a1 := edgeStubSkipArgs(tb)
			expect.Equal(tb, c.A0, a0, "the call records A0")
			expect.Equal(tb, c.E, a1, "the call records E")
		},
		Override: func(seen func(edge.EdgeSkipCall)) {
			s.OnSkip.Func(func(p0 string, p1 int) {
				seen(edge.EdgeSkipCall{A0: p0, E: p1})
			})
		},
	}
}

// edgeStubReadSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Read. The subject calls a new double that
// reports to tb.
func edgeStubReadSubject(tb assert.TB) stub.Subject[edge.EdgeReadCall, edge.EdgeReadReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeReadCall, edge.EdgeReadReturn]{
		Stub:  s.OnRead.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = edgeStubRead(tb, s)
		},
		Result: func() edge.EdgeReadReturn {
			got0, got1 := edgeStubRead(tb, s)
			return edge.EdgeReadReturn{R: got0, Err: got1}
		},
		Returns: func() edge.EdgeReadReturn {
			want0, want1 := edgeStubReadPins()
			s.OnRead.Returns(want0, want1)
			return edge.EdgeReadReturn{R: want0, Err: want1}
		},
		Override: func(seen func(edge.EdgeReadCall)) {
			s.OnRead.Func(func() (int, error) {
				seen(edge.EdgeReadCall{})
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

// edgeStubKeepSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Keep. The subject calls a new double that
// reports to tb.
func edgeStubKeepSubject(tb assert.TB) stub.Subject[edge.EdgeKeepCall, edge.EdgeKeepReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeKeepCall, edge.EdgeKeepReturn]{
		Stub:  s.OnKeep.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = edgeStubKeep(tb, s)
		},
		Arguments: func(tb assert.TB, c edge.EdgeKeepCall) {
			tb.Helper()
			a0 := edgeStubKeepArgs(tb)
			assert.Equal(tb, c.R, a0, "the call records R")
		},
		Result: func() edge.EdgeKeepReturn {
			got0 := edgeStubKeep(tb, s)
			return edge.EdgeKeepReturn{Result: got0}
		},
		Returns: func() edge.EdgeKeepReturn {
			want0 := edgeStubKeepPins()
			s.OnKeep.Returns(want0)
			return edge.EdgeKeepReturn{Result: want0}
		},
		Override: func(seen func(edge.EdgeKeepCall)) {
			s.OnKeep.Func(func(p0 int) int {
				seen(edge.EdgeKeepCall{R: p0})
				var z0 int
				return z0
			})
		},
	}
}

// edgeStubUseSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Use. The subject calls a new double that
// reports to tb.
func edgeStubUseSubject(tb assert.TB) stub.Subject[edge.EdgeUseCall, edge.EdgeUseReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeUseCall, edge.EdgeUseReturn]{
		Stub:  s.OnUse.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = edgeStubUse(tb, s)
		},
		Arguments: func(tb assert.TB, c edge.EdgeUseCall) {
			tb.Helper()
			a0 := edgeStubUseArgs(tb)
			assert.Equal(tb, c.S, a0, "the call records S")
		},
		Result: func() edge.EdgeUseReturn {
			got0 := edgeStubUse(tb, s)
			return edge.EdgeUseReturn{Err: got0}
		},
		Returns: func() edge.EdgeUseReturn {
			want0 := edgeStubUsePins()
			s.OnUse.Returns(want0)
			return edge.EdgeUseReturn{Err: want0}
		},
		Override: func(seen func(edge.EdgeUseCall)) {
			s.OnUse.Func(func(p0 string) error {
				seen(edge.EdgeUseCall{S: p0})
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

// edgeStubRecordSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Record. The subject calls a new double that
// reports to tb.
func edgeStubRecordSubject(tb assert.TB) stub.Subject[edge.EdgeRecordCall, edge.EdgeRecordReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeRecordCall, edge.EdgeRecordReturn]{
		Stub:  s.OnRecord.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = edgeStubRecord(tb, s)
		},
		Arguments: func(tb assert.TB, c edge.EdgeRecordCall) {
			tb.Helper()
			a0 := edgeStubRecordArgs(tb)
			assert.Equal(tb, c.Call, a0, "the call records Call")
		},
		Result: func() edge.EdgeRecordReturn {
			got0, got1 := edgeStubRecord(tb, s)
			return edge.EdgeRecordReturn{N: got0, Err: got1}
		},
		Returns: func() edge.EdgeRecordReturn {
			want0, want1 := edgeStubRecordPins()
			s.OnRecord.Returns(want0, want1)
			return edge.EdgeRecordReturn{N: want0, Err: want1}
		},
		Override: func(seen func(edge.EdgeRecordCall)) {
			s.OnRecord.Func(func(p0 string) (int, error) {
				seen(edge.EdgeRecordCall{Call: p0})
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

// edgeStubClashSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Clash. The subject calls a new double that
// reports to tb.
func edgeStubClashSubject(tb assert.TB) stub.Subject[edge.EdgeClashCall, edge.EdgeClashReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeClashCall, edge.EdgeClashReturn]{
		Stub:  s.OnClash.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = edgeStubClash(tb, s)
		},
		Result: func() edge.EdgeClashReturn {
			got0, got1 := edgeStubClash(tb, s)
			return edge.EdgeClashReturn{S: got0, Err: got1}
		},
		Returns: func() edge.EdgeClashReturn {
			want0, want1 := edgeStubClashPins()
			s.OnClash.Returns(want0, want1)
			return edge.EdgeClashReturn{S: want0, Err: want1}
		},
		Override: func(seen func(edge.EdgeClashCall)) {
			s.OnClash.Func(func() (string, error) {
				seen(edge.EdgeClashCall{})
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

// edgeStubShiftSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Shift. The subject calls a new double that
// reports to tb.
func edgeStubShiftSubject(tb assert.TB) stub.Subject[edge.EdgeShiftCall, struct{}] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeShiftCall, struct{}]{
		Stub:  s.OnShift.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			edgeStubShift(tb, s)
		},
		Arguments: func(tb assert.TB, c edge.EdgeShiftCall) {
			tb.Helper()
			a0, a1 := edgeStubShiftArgs(tb)
			expect.Equal(tb, c.A00, a0, "the call records A00")
			expect.Equal(tb, c.A0, a1, "the call records A0")
		},
		Override: func(seen func(edge.EdgeShiftCall)) {
			s.OnShift.Func(func(p0 string, p1 int) {
				seen(edge.EdgeShiftCall{A00: p0, A0: p1})
			})
		},
	}
}

// edgeStubFetchSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Fetch. The subject calls a new double that
// reports to tb.
func edgeStubFetchSubject(tb assert.TB) stub.Subject[edge.EdgeFetchCall, edge.EdgeFetchReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeFetchCall, edge.EdgeFetchReturn]{
		Stub:  s.OnFetch.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = edgeStubFetch(tb, s)
		},
		Arguments: func(tb assert.TB, c edge.EdgeFetchCall) {
			tb.Helper()
			a0 := edgeStubFetchArgs(tb)
			assert.Equal(tb, c.A00, a0, "the call records A00")
		},
		Result: func() edge.EdgeFetchReturn {
			got0, got1 := edgeStubFetch(tb, s)
			return edge.EdgeFetchReturn{A0: got0, Err: got1}
		},
		Returns: func() edge.EdgeFetchReturn {
			want0, want1 := edgeStubFetchPins()
			s.OnFetch.Returns(want0, want1)
			return edge.EdgeFetchReturn{A0: want0, Err: want1}
		},
		Override: func(seen func(edge.EdgeFetchCall)) {
			s.OnFetch.Func(func(p0 string) (int, error) {
				seen(edge.EdgeFetchCall{A00: p0})
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

// edgeStubWideSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Wide. The subject calls a new double that
// reports to tb.
func edgeStubWideSubject(tb assert.TB) stub.Subject[edge.EdgeWideCall, struct{}] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgeWideCall, struct{}]{
		Stub:  s.OnWide.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			edgeStubWide(tb, s)
		},
		Arguments: func(tb assert.TB, c edge.EdgeWideCall) {
			tb.Helper()
			a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 := edgeStubWideArgs(tb)
			expect.Equal(tb, c.A, a0, "the call records A")
			expect.Equal(tb, c.A10, a1, "the call records A10")
			expect.Equal(tb, c.A1, a2, "the call records A1")
			expect.Equal(tb, c.B, a3, "the call records B")
			expect.Equal(tb, c.C, a4, "the call records C")
			expect.Equal(tb, c.D, a5, "the call records D")
			expect.Equal(tb, c.E, a6, "the call records E")
			expect.Equal(tb, c.F, a7, "the call records F")
			expect.Equal(tb, c.G, a8, "the call records G")
			expect.Equal(tb, c.H, a9, "the call records H")
			expect.Equal(tb, c.A100, a10, "the call records A100")
		},
		Override: func(seen func(edge.EdgeWideCall)) {
			s.OnWide.Func(func(p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10 int) {
				seen(
					edge.EdgeWideCall{
						A:    p0,
						A10:  p1,
						A1:   p2,
						B:    p3,
						C:    p4,
						D:    p5,
						E:    p6,
						F:    p7,
						G:    p8,
						H:    p9,
						A100: p10,
					},
				)
			})
		},
	}
}

// edgeStubPairSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Pair. The subject calls a new double that
// reports to tb.
func edgeStubPairSubject(tb assert.TB) stub.Subject[edge.EdgePairCall, edge.EdgePairReturn] {
	tb.Helper()
	s := edge.NewEdgeStub(tb)
	return stub.Subject[edge.EdgePairCall, edge.EdgePairReturn]{
		Stub:  s.OnPair.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = edgeStubPair(tb, s)
		},
		Result: func() edge.EdgePairReturn {
			got0, got1 := edgeStubPair(tb, s)
			return edge.EdgePairReturn{AB: got0, Ab: got1}
		},
		Returns: func() edge.EdgePairReturn {
			want0, want1 := edgeStubPairPins()
			s.OnPair.Returns(want0, want1)
			return edge.EdgePairReturn{AB: want0, Ab: want1}
		},
		Override: func(seen func(edge.EdgePairCall)) {
			s.OnPair.Func(func() (int, int) {
				seen(edge.EdgePairCall{})
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
	}
}

// TestEdgeStub checks the methods, the constructors and the options of the
// double EdgeStub.
func TestEdgeStub(t *testing.T) {
	t.Parallel()

	t.Run("Skip", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Skip", edgeStubSkipSubject)
	})

	t.Run("Read", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Read", edgeStubReadSubject)
	})

	t.Run("Keep", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Keep", edgeStubKeepSubject)
	})

	t.Run("Use", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Use", edgeStubUseSubject)
	})

	t.Run("Record", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Record", edgeStubRecordSubject)
	})

	t.Run("Clash", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Clash", edgeStubClashSubject)
	})

	t.Run("Shift", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Shift", edgeStubShiftSubject)
	})

	t.Run("Fetch", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Fetch", edgeStubFetchSubject)
	})

	t.Run("Wide", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Wide", edgeStubWideSubject)
	})

	t.Run("Pair", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Pair", edgeStubPairSubject)
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("Skip", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Skip", edgeStubSkipSubject)
		})

		t.Run("Read", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Read", edgeStubReadSubject)
		})

		t.Run("Keep", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Keep", edgeStubKeepSubject)
		})

		t.Run("Use", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Use", edgeStubUseSubject)
		})

		t.Run("Record", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Record", edgeStubRecordSubject)
		})

		t.Run("Clash", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Clash", edgeStubClashSubject)
		})

		t.Run("Shift", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Shift", edgeStubShiftSubject)
		})

		t.Run("Fetch", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Fetch", edgeStubFetchSubject)
		})

		t.Run("Wide", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Wide", edgeStubWideSubject)
		})

		t.Run("Pair", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Pair", edgeStubPairSubject)
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
