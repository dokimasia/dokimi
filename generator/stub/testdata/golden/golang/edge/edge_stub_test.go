package edge

import (
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/dokimi-go/stub"
)

// edgeStubSkipSubject binds Skip into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubSkipSubject(tb assert.TB) stub.Subject[EdgeSkipCall, struct{}] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeSkipCall, struct{}]{
		Stub: s.OnSkip.MethodStub,
		Call: func() {
			var a0 string
			var a1 int
			a0 = "Skip-a"
			a1 = 42
			s.Skip(a0, a1)
		},
		Override: func(mark func()) {
			s.OnSkip.Func(func(_ string, _ int) {
				mark()
			})
		},
	}
}

// TestEdgeStubSkip pins how Skip answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubSkip(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Skip", edgeStubSkipSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var a0 string
		var a1 int
		a0 = "Skip-a"
		a1 = 42
		s.Skip(a0, a1)
		got := s.OnSkip.AssertCalledOnce(t, "Skip must record the call")
		assert.Equal(t, got.A0, a0, "the recorded call carries A0")
		assert.Equal(t, got.E, a1, "the recorded call carries E")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeSkipCall
		s.OnSkip.OnRecord(func(c EdgeSkipCall) { seen = append(seen, c) })
		var a0 string
		var a1 int
		a0 = "Skip-a"
		a1 = 42
		s.Skip(a0, a1)
		s.Skip(a0, a1)
		assert.Length(t, seen, 2, "OnRecord must fire once per Skip call")
	})

	t.Run("wires WithEdgeSkip at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeSkip(func(_ string, _ int) {
			called = true
		}))
		var a0 string
		var a1 int
		a0 = "Skip-a"
		a1 = 42
		s.Skip(a0, a1)
		assert.True(t, called, "WithEdgeSkip must install the override")
	})
}

// edgeStubReadSubject binds Read into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubReadSubject(tb assert.TB) stub.Subject[EdgeReadCall, EdgeReadReturn] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeReadCall, EdgeReadReturn]{
		Stub: s.OnRead.MethodStub,
		Call: func() {
			_, _ = s.Read()
		},
		Result: func() EdgeReadReturn {
			got0, got1 := s.Read()
			return EdgeReadReturn{R: got0, Err: got1}
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
			_, r1 := s.Read()
			return r1
		},
	}
}

// TestEdgeStubRead pins how Read answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubRead(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Read", edgeStubReadSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 error
		want0 = 42
		s.OnRead.Returns(want0, want1)
		got0, _ := s.Read()
		assert.Equal(t, got0, want0, "Read must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		_, _ = s.Read()
		s.OnRead.AssertCalledOnce(t, "Read must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeReadCall
		s.OnRead.OnRecord(func(c EdgeReadCall) { seen = append(seen, c) })
		_, _ = s.Read()
		_, _ = s.Read()
		assert.Length(t, seen, 2, "OnRecord must fire once per Read call")
	})

	t.Run("wires WithEdgeRead at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeRead(func() (int, error) {
			called = true
			var z0 int
			var z1 error
			return z0, z1
		}))
		_, _ = s.Read()
		assert.True(t, called, "WithEdgeRead must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Read-failing")
		s := NewEdgeStubFailing(t, want)
		_, r1 := s.Read()
		assert.ErrorIs(t, r1, want, "Read must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 error
		want0 = 42
		s.OnRead.Returns(want0, want1)
		_, _ = s.Read()
		s.ResetCalls()
		got0, _ := s.Read()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// edgeStubKeepSubject binds Keep into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubKeepSubject(tb assert.TB) stub.Subject[EdgeKeepCall, EdgeKeepReturn] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeKeepCall, EdgeKeepReturn]{
		Stub: s.OnKeep.MethodStub,
		Call: func() {
			var a0 int
			a0 = 42
			_ = s.Keep(a0)
		},
		Result: func() EdgeKeepReturn {
			var a0 int
			a0 = 42
			got0 := s.Keep(a0)
			return EdgeKeepReturn{Result: got0}
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

// TestEdgeStubKeep pins how Keep answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubKeep(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Keep", edgeStubKeepSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		want0 = 42
		s.OnKeep.Returns(want0)
		var a0 int
		a0 = 42
		got0 := s.Keep(a0)
		assert.Equal(t, got0, want0, "Keep must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var a0 int
		a0 = 42
		_ = s.Keep(a0)
		got := s.OnKeep.AssertCalledOnce(t, "Keep must record the call")
		assert.Equal(t, got.R, a0, "the recorded call carries R")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeKeepCall
		s.OnKeep.OnRecord(func(c EdgeKeepCall) { seen = append(seen, c) })
		var a0 int
		a0 = 42
		_ = s.Keep(a0)
		_ = s.Keep(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Keep call")
	})

	t.Run("wires WithEdgeKeep at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeKeep(func(_ int) int {
			called = true
			var z0 int
			return z0
		}))
		var a0 int
		a0 = 42
		_ = s.Keep(a0)
		assert.True(t, called, "WithEdgeKeep must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		want0 = 42
		s.OnKeep.Returns(want0)
		var a0 int
		a0 = 42
		_ = s.Keep(a0)
		s.ResetCalls()
		got0 := s.Keep(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// edgeStubUseSubject binds Use into the shape [stub.Behaviour] drives: how
// to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubUseSubject(tb assert.TB) stub.Subject[EdgeUseCall, EdgeUseReturn] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeUseCall, EdgeUseReturn]{
		Stub: s.OnUse.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Use-a"
			_ = s.Use(a0)
		},
		Result: func() EdgeUseReturn {
			var a0 string
			a0 = "Use-a"
			got0 := s.Use(a0)
			return EdgeUseReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnUse.Func(func(_ string) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			var a0 string
			a0 = "Use-a"
			r0 := s.Use(a0)
			return r0
		},
	}
}

// TestEdgeStubUse pins how Use answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubUse(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Use", edgeStubUseSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var a0 string
		a0 = "Use-a"
		_ = s.Use(a0)
		got := s.OnUse.AssertCalledOnce(t, "Use must record the call")
		assert.Equal(t, got.S, a0, "the recorded call carries S")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeUseCall
		s.OnUse.OnRecord(func(c EdgeUseCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Use-a"
		_ = s.Use(a0)
		_ = s.Use(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Use call")
	})

	t.Run("wires WithEdgeUse at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeUse(func(_ string) error {
			called = true
			var z0 error
			return z0
		}))
		var a0 string
		a0 = "Use-a"
		_ = s.Use(a0)
		assert.True(t, called, "WithEdgeUse must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Use-failing")
		s := NewEdgeStubFailing(t, want)
		var a0 string
		a0 = "Use-a"
		r0 := s.Use(a0)
		assert.ErrorIs(t, r0, want, "Use must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 error
		s.OnUse.Returns(want0)
		var a0 string
		a0 = "Use-a"
		_ = s.Use(a0)
		s.ResetCalls()
		_ = s.Use(a0)
		s.OnUse.AssertCalledOnce(t, "Use must record the call that followed the reset")
	})
}

// edgeStubRecordSubject binds Record into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubRecordSubject(tb assert.TB) stub.Subject[EdgeRecordCall, EdgeRecordReturn] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeRecordCall, EdgeRecordReturn]{
		Stub: s.OnRecord.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Record-a"
			_, _ = s.Record(a0)
		},
		Result: func() EdgeRecordReturn {
			var a0 string
			a0 = "Record-a"
			got0, got1 := s.Record(a0)
			return EdgeRecordReturn{N: got0, Err: got1}
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
			var a0 string
			a0 = "Record-a"
			_, err := s.Record(a0)
			return err
		},
	}
}

// TestEdgeStubRecord pins how Record answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubRecord(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Record", edgeStubRecordSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 error
		want0 = 42
		s.OnRecord.Returns(want0, want1)
		var a0 string
		a0 = "Record-a"
		got0, _ := s.Record(a0)
		assert.Equal(t, got0, want0, "Record must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var a0 string
		a0 = "Record-a"
		_, _ = s.Record(a0)
		got := s.OnRecord.AssertCalledOnce(t, "Record must record the call")
		assert.Equal(t, got.Call, a0, "the recorded call carries Call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeRecordCall
		s.OnRecord.OnRecord(func(c EdgeRecordCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Record-a"
		_, _ = s.Record(a0)
		_, _ = s.Record(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Record call")
	})

	t.Run("wires WithEdgeRecord at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeRecord(func(_ string) (int, error) {
			called = true
			var z0 int
			var z1 error
			return z0, z1
		}))
		var a0 string
		a0 = "Record-a"
		_, _ = s.Record(a0)
		assert.True(t, called, "WithEdgeRecord must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Record-failing")
		s := NewEdgeStubFailing(t, want)
		var a0 string
		a0 = "Record-a"
		_, err := s.Record(a0)
		assert.ErrorIs(t, err, want, "Record must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 error
		want0 = 42
		s.OnRecord.Returns(want0, want1)
		var a0 string
		a0 = "Record-a"
		_, _ = s.Record(a0)
		s.ResetCalls()
		got0, _ := s.Record(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// edgeStubClashSubject binds Clash into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubClashSubject(tb assert.TB) stub.Subject[EdgeClashCall, EdgeClashReturn] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeClashCall, EdgeClashReturn]{
		Stub: s.OnClash.MethodStub,
		Call: func() {
			_, _ = s.Clash()
		},
		Result: func() EdgeClashReturn {
			got0, got1 := s.Clash()
			return EdgeClashReturn{S: got0, Err: got1}
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
			_, r1 := s.Clash()
			return r1
		},
	}
}

// TestEdgeStubClash pins how Clash answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubClash(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Clash", edgeStubClashSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 string
		var want1 error
		want0 = "Clash-a"
		s.OnClash.Returns(want0, want1)
		got0, _ := s.Clash()
		assert.Equal(t, got0, want0, "Clash must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		_, _ = s.Clash()
		s.OnClash.AssertCalledOnce(t, "Clash must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeClashCall
		s.OnClash.OnRecord(func(c EdgeClashCall) { seen = append(seen, c) })
		_, _ = s.Clash()
		_, _ = s.Clash()
		assert.Length(t, seen, 2, "OnRecord must fire once per Clash call")
	})

	t.Run("wires WithEdgeClash at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeClash(func() (string, error) {
			called = true
			var z0 string
			var z1 error
			return z0, z1
		}))
		_, _ = s.Clash()
		assert.True(t, called, "WithEdgeClash must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Clash-failing")
		s := NewEdgeStubFailing(t, want)
		_, r1 := s.Clash()
		assert.ErrorIs(t, r1, want, "Clash must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 string
		var want1 error
		want0 = "Clash-a"
		s.OnClash.Returns(want0, want1)
		_, _ = s.Clash()
		s.ResetCalls()
		got0, _ := s.Clash()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// edgeStubShiftSubject binds Shift into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubShiftSubject(tb assert.TB) stub.Subject[EdgeShiftCall, struct{}] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeShiftCall, struct{}]{
		Stub: s.OnShift.MethodStub,
		Call: func() {
			var a0 string
			var a1 int
			a0 = "Shift-a"
			a1 = 42
			s.Shift(a0, a1)
		},
		Override: func(mark func()) {
			s.OnShift.Func(func(_ string, _ int) {
				mark()
			})
		},
	}
}

// TestEdgeStubShift pins how Shift answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubShift(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Shift", edgeStubShiftSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var a0 string
		var a1 int
		a0 = "Shift-a"
		a1 = 42
		s.Shift(a0, a1)
		got := s.OnShift.AssertCalledOnce(t, "Shift must record the call")
		assert.Equal(t, got.A00, a0, "the recorded call carries A00")
		assert.Equal(t, got.A0, a1, "the recorded call carries A0")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeShiftCall
		s.OnShift.OnRecord(func(c EdgeShiftCall) { seen = append(seen, c) })
		var a0 string
		var a1 int
		a0 = "Shift-a"
		a1 = 42
		s.Shift(a0, a1)
		s.Shift(a0, a1)
		assert.Length(t, seen, 2, "OnRecord must fire once per Shift call")
	})

	t.Run("wires WithEdgeShift at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeShift(func(_ string, _ int) {
			called = true
		}))
		var a0 string
		var a1 int
		a0 = "Shift-a"
		a1 = 42
		s.Shift(a0, a1)
		assert.True(t, called, "WithEdgeShift must install the override")
	})
}

// edgeStubFetchSubject binds Fetch into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubFetchSubject(tb assert.TB) stub.Subject[EdgeFetchCall, EdgeFetchReturn] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeFetchCall, EdgeFetchReturn]{
		Stub: s.OnFetch.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Fetch-a"
			_, _ = s.Fetch(a0)
		},
		Result: func() EdgeFetchReturn {
			var a0 string
			a0 = "Fetch-a"
			got0, got1 := s.Fetch(a0)
			return EdgeFetchReturn{A0: got0, Err: got1}
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
			var a0 string
			a0 = "Fetch-a"
			_, err := s.Fetch(a0)
			return err
		},
	}
}

// TestEdgeStubFetch pins how Fetch answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubFetch(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Fetch", edgeStubFetchSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 error
		want0 = 42
		s.OnFetch.Returns(want0, want1)
		var a0 string
		a0 = "Fetch-a"
		got0, _ := s.Fetch(a0)
		assert.Equal(t, got0, want0, "Fetch must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var a0 string
		a0 = "Fetch-a"
		_, _ = s.Fetch(a0)
		got := s.OnFetch.AssertCalledOnce(t, "Fetch must record the call")
		assert.Equal(t, got.A00, a0, "the recorded call carries A00")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeFetchCall
		s.OnFetch.OnRecord(func(c EdgeFetchCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Fetch-a"
		_, _ = s.Fetch(a0)
		_, _ = s.Fetch(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Fetch call")
	})

	t.Run("wires WithEdgeFetch at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeFetch(func(_ string) (int, error) {
			called = true
			var z0 int
			var z1 error
			return z0, z1
		}))
		var a0 string
		a0 = "Fetch-a"
		_, _ = s.Fetch(a0)
		assert.True(t, called, "WithEdgeFetch must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Fetch-failing")
		s := NewEdgeStubFailing(t, want)
		var a0 string
		a0 = "Fetch-a"
		_, err := s.Fetch(a0)
		assert.ErrorIs(t, err, want, "Fetch must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 error
		want0 = 42
		s.OnFetch.Returns(want0, want1)
		var a0 string
		a0 = "Fetch-a"
		_, _ = s.Fetch(a0)
		s.ResetCalls()
		got0, _ := s.Fetch(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// edgeStubWideSubject binds Wide into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubWideSubject(tb assert.TB) stub.Subject[EdgeWideCall, struct{}] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgeWideCall, struct{}]{
		Stub: s.OnWide.MethodStub,
		Call: func() {
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
			s.Wide(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
		},
		Override: func(mark func()) {
			s.OnWide.Func(func(_ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int) {
				mark()
			})
		},
	}
}

// TestEdgeStubWide pins how Wide answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubWide(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Wide", edgeStubWideSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
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
		s.Wide(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
		got := s.OnWide.AssertCalledOnce(t, "Wide must record the call")
		assert.Equal(t, got.A, a0, "the recorded call carries A")
		assert.Equal(t, got.A10, a1, "the recorded call carries A10")
		assert.Equal(t, got.A1, a2, "the recorded call carries A1")
		assert.Equal(t, got.B, a3, "the recorded call carries B")
		assert.Equal(t, got.C, a4, "the recorded call carries C")
		assert.Equal(t, got.D, a5, "the recorded call carries D")
		assert.Equal(t, got.E, a6, "the recorded call carries E")
		assert.Equal(t, got.F, a7, "the recorded call carries F")
		assert.Equal(t, got.G, a8, "the recorded call carries G")
		assert.Equal(t, got.H, a9, "the recorded call carries H")
		assert.Equal(t, got.A100, a10, "the recorded call carries A100")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgeWideCall
		s.OnWide.OnRecord(func(c EdgeWideCall) { seen = append(seen, c) })
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
		s.Wide(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
		s.Wide(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
		assert.Length(t, seen, 2, "OnRecord must fire once per Wide call")
	})

	t.Run("wires WithEdgeWide at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgeWide(func(_ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int, _ int) {
			called = true
		}))
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
		s.Wide(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
		assert.True(t, called, "WithEdgeWide must install the override")
	})
}

// edgeStubPairSubject binds Pair into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func edgeStubPairSubject(tb assert.TB) stub.Subject[EdgePairCall, EdgePairReturn] {
	tb.Helper()
	s := NewEdgeStub(tb)
	return stub.Subject[EdgePairCall, EdgePairReturn]{
		Stub: s.OnPair.MethodStub,
		Call: func() {
			_, _ = s.Pair()
		},
		Result: func() EdgePairReturn {
			got0, got1 := s.Pair()
			return EdgePairReturn{AB: got0, Ab: got1}
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

// TestEdgeStubPair pins how Pair answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestEdgeStubPair(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Pair", edgeStubPairSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 int
		want0 = 42
		want1 = 42
		s.OnPair.Returns(want0, want1)
		got0, got1 := s.Pair()
		assert.Equal(t, got0, want0, "Pair must answer with what Returns pinned")
		assert.Equal(t, got1, want1, "Pair must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		_, _ = s.Pair()
		s.OnPair.AssertCalledOnce(t, "Pair must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var seen []EdgePairCall
		s.OnPair.OnRecord(func(c EdgePairCall) { seen = append(seen, c) })
		_, _ = s.Pair()
		_, _ = s.Pair()
		assert.Length(t, seen, 2, "OnRecord must fire once per Pair call")
	})

	t.Run("wires WithEdgePair at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewEdgeStub(t, WithEdgePair(func() (int, int) {
			called = true
			var z0 int
			var z1 int
			return z0, z1
		}))
		_, _ = s.Pair()
		assert.True(t, called, "WithEdgePair must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewEdgeStub(t)
		var want0 int
		var want1 int
		want0 = 42
		want1 = 42
		s.OnPair.Returns(want0, want1)
		_, _ = s.Pair()
		s.ResetCalls()
		got0, got1 := s.Pair()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
		assert.Equal(t, got1, want1, "a reset must keep what Returns pinned")
	})
}

// edgeStubDouble describes how to build a EdgeStub under each option whose
// effect is the same whatever a method's signature.
//
// Skip stands in for the double as a whole: these checks assert only that
// an option reached it.
func edgeStubDouble() stub.Double[EdgeSkipCall] {
	instance := func(s *EdgeStub) stub.Instance[EdgeSkipCall] {
		return stub.Instance[EdgeSkipCall]{
			Stub: s.OnSkip.MethodStub,
			Call: func() {
				var a0 string
				var a1 int
				a0 = "Skip-a"
				a1 = 42
				s.Skip(a0, a1)
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[EdgeSkipCall]{
		New: func(tb assert.TB) stub.Instance[EdgeSkipCall] {
			return instance(NewEdgeStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[EdgeSkipCall] {
			return instance(NewEdgeStub(tb, EdgeStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[EdgeSkipCall] {
			return instance(NewEdgeStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[EdgeSkipCall] {
			return instance(NewEdgeStub(tb, EdgeStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[EdgeSkipCall] {
			return instance(NewEdgeStub(tb, EdgeStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[EdgeSkipCall] {
			return instance(NewEdgeStub(tb, EdgeStubStrict()))
		},
	}
}

// TestEdgeStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestEdgeStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, edgeStubDouble())
}

// TestEdgeStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestEdgeStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewEdgeStub(t)
	s := NewEdgeStub(t, EdgeStubDelegateTo(inner))

	t.Run("forwards Skip to the wrapped implementation", func(t *testing.T) {
		var a0 string
		var a1 int
		a0 = "Skip-a"
		a1 = 42
		s.Skip(a0, a1)
		inner.OnSkip.AssertCalledOnce(t, "Skip must reach the wrapped implementation")
	})

	t.Run("forwards Read to the wrapped implementation", func(t *testing.T) {
		_, _ = s.Read()
		inner.OnRead.AssertCalledOnce(t, "Read must reach the wrapped implementation")
	})

	t.Run("surfaces what Read answered", func(t *testing.T) {
		want := errors.New("Read-delegate")
		inner.OnRead.FaultsFor(time.Hour, want)
		_, r1 := s.Read()
		assert.ErrorIs(t, r1, want, "Read must surface the wrapped answer")
	})

	t.Run("forwards Keep to the wrapped implementation", func(t *testing.T) {
		var a0 int
		a0 = 42
		_ = s.Keep(a0)
		inner.OnKeep.AssertCalledOnce(t, "Keep must reach the wrapped implementation")
	})

	t.Run("forwards Use to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Use-a"
		_ = s.Use(a0)
		inner.OnUse.AssertCalledOnce(t, "Use must reach the wrapped implementation")
	})

	t.Run("surfaces what Use answered", func(t *testing.T) {
		want := errors.New("Use-delegate")
		inner.OnUse.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "Use-a"
		r0 := s.Use(a0)
		assert.ErrorIs(t, r0, want, "Use must surface the wrapped answer")
	})

	t.Run("forwards Record to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Record-a"
		_, _ = s.Record(a0)
		inner.OnRecord.AssertCalledOnce(t, "Record must reach the wrapped implementation")
	})

	t.Run("surfaces what Record answered", func(t *testing.T) {
		want := errors.New("Record-delegate")
		inner.OnRecord.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "Record-a"
		_, err := s.Record(a0)
		assert.ErrorIs(t, err, want, "Record must surface the wrapped answer")
	})

	t.Run("forwards Clash to the wrapped implementation", func(t *testing.T) {
		_, _ = s.Clash()
		inner.OnClash.AssertCalledOnce(t, "Clash must reach the wrapped implementation")
	})

	t.Run("surfaces what Clash answered", func(t *testing.T) {
		want := errors.New("Clash-delegate")
		inner.OnClash.FaultsFor(time.Hour, want)
		_, r1 := s.Clash()
		assert.ErrorIs(t, r1, want, "Clash must surface the wrapped answer")
	})

	t.Run("forwards Shift to the wrapped implementation", func(t *testing.T) {
		var a0 string
		var a1 int
		a0 = "Shift-a"
		a1 = 42
		s.Shift(a0, a1)
		inner.OnShift.AssertCalledOnce(t, "Shift must reach the wrapped implementation")
	})

	t.Run("forwards Fetch to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Fetch-a"
		_, _ = s.Fetch(a0)
		inner.OnFetch.AssertCalledOnce(t, "Fetch must reach the wrapped implementation")
	})

	t.Run("surfaces what Fetch answered", func(t *testing.T) {
		want := errors.New("Fetch-delegate")
		inner.OnFetch.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "Fetch-a"
		_, err := s.Fetch(a0)
		assert.ErrorIs(t, err, want, "Fetch must surface the wrapped answer")
	})

	t.Run("forwards Wide to the wrapped implementation", func(t *testing.T) {
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
		s.Wide(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
		inner.OnWide.AssertCalledOnce(t, "Wide must reach the wrapped implementation")
	})

	t.Run("forwards Pair to the wrapped implementation", func(t *testing.T) {
		_, _ = s.Pair()
		inner.OnPair.AssertCalledOnce(t, "Pair must reach the wrapped implementation")
	})
}
