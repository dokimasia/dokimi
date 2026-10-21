package edge

import (
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// EdgeSkipCall records one invocation of Edge.Skip.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeSkipCall struct {
	A0 string
	E  int
}

// EdgeReadCall records one invocation of Edge.Read.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeReadCall struct {
	R   int
	Err error
}

// EdgeKeepCall records one invocation of Edge.Keep.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeKeepCall struct {
	R      int
	Result int
}

// EdgeUseCall records one invocation of Edge.Use.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeUseCall struct {
	S   string
	Err error
}

// EdgeRecordCall records one invocation of Edge.Record.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeRecordCall struct {
	Call string
	N    int
	Err  error
}

// EdgeClashCall records one invocation of Edge.Clash.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeClashCall struct {
	S   string
	Err error
}

// EdgeShiftCall records one invocation of Edge.Shift.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeShiftCall struct {
	A00 string
	A0  int
}

// EdgeFetchCall records one invocation of Edge.Fetch.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeFetchCall struct {
	A00 string
	A0  int
	Err error
}

// EdgeWideCall records one invocation of Edge.Wide.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgeWideCall struct {
	A    int
	A10  int
	A1   int
	B    int
	C    int
	D    int
	E    int
	F    int
	G    int
	H    int
	A100 int
}

// EdgePairCall records one invocation of Edge.Pair.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type EdgePairCall struct {
	AB int
	Ab int
}

// --- Per-method configuration ---

// EdgeSkipStub controls how the double answers Skip and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeSkipStub struct {
	*stub.MethodStub[EdgeSkipCall]

	fn func(string, int)
}

// Func supplies a body for Skip, for when the answer depends on the
// arguments. Nothing takes precedence over it: Skip answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *EdgeSkipStub) Func(fn func(string, int)) *EdgeSkipStub {
	s.fn = fn
	return s
}

// EdgeReadStub controls how the double answers Read and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeReadStub struct {
	*stub.MethodStub[EdgeReadCall]

	fn       func() (int, error)
	fallback *EdgeReadReturn
}

// EdgeReadReturn holds the fixed answer configured through Returns.
type EdgeReadReturn struct {
	R   int
	Err error
}

// Returns pins a fixed result for every call to Read. A Func override
// takes precedence over it, and so does an injected fault.
func (s *EdgeReadStub) Returns(r int, err error) *EdgeReadStub {
	s.fallback = &EdgeReadReturn{R: r, Err: err}
	return s
}

// Func supplies a body for Read, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *EdgeReadStub) Func(fn func() (int, error)) *EdgeReadStub {
	s.fn = fn
	return s
}

// EdgeKeepStub controls how the double answers Keep and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeKeepStub struct {
	*stub.MethodStub[EdgeKeepCall]

	fn       func(int) int
	fallback *EdgeKeepReturn
}

// EdgeKeepReturn holds the fixed answer configured through Returns.
type EdgeKeepReturn struct {
	Result int
}

// Returns pins a fixed result for every call to Keep. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Keep answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *EdgeKeepStub) Returns(result int) *EdgeKeepStub {
	s.fallback = &EdgeKeepReturn{Result: result}
	return s
}

// Func supplies a body for Keep, for when the answer depends on the
// arguments. Nothing takes precedence over it: Keep answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *EdgeKeepStub) Func(fn func(int) int) *EdgeKeepStub {
	s.fn = fn
	return s
}

// EdgeUseStub controls how the double answers Use and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeUseStub struct {
	*stub.MethodStub[EdgeUseCall]

	fn       func(string) error
	fallback *EdgeUseReturn
}

// EdgeUseReturn holds the fixed answer configured through Returns.
type EdgeUseReturn struct {
	Err error
}

// Returns pins a fixed result for every call to Use. A Func override takes
// precedence over it, and so does an injected fault.
func (s *EdgeUseStub) Returns(err error) *EdgeUseStub {
	s.fallback = &EdgeUseReturn{Err: err}
	return s
}

// Func supplies a body for Use, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *EdgeUseStub) Func(fn func(string) error) *EdgeUseStub {
	s.fn = fn
	return s
}

// EdgeRecordStub controls how the double answers Record and records what
// it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeRecordStub struct {
	*stub.MethodStub[EdgeRecordCall]

	fn       func(string) (int, error)
	fallback *EdgeRecordReturn
}

// EdgeRecordReturn holds the fixed answer configured through Returns.
type EdgeRecordReturn struct {
	N   int
	Err error
}

// Returns pins a fixed result for every call to Record. A Func override
// takes precedence over it, and so does an injected fault.
func (s *EdgeRecordStub) Returns(n int, err error) *EdgeRecordStub {
	s.fallback = &EdgeRecordReturn{N: n, Err: err}
	return s
}

// Func supplies a body for Record, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *EdgeRecordStub) Func(fn func(string) (int, error)) *EdgeRecordStub {
	s.fn = fn
	return s
}

// EdgeClashStub controls how the double answers Clash and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeClashStub struct {
	*stub.MethodStub[EdgeClashCall]

	fn       func() (string, error)
	fallback *EdgeClashReturn
}

// EdgeClashReturn holds the fixed answer configured through Returns.
type EdgeClashReturn struct {
	S   string
	Err error
}

// Returns pins a fixed result for every call to Clash. A Func override
// takes precedence over it, and so does an injected fault.
func (s *EdgeClashStub) Returns(s0 string, err error) *EdgeClashStub {
	s.fallback = &EdgeClashReturn{S: s0, Err: err}
	return s
}

// Func supplies a body for Clash, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *EdgeClashStub) Func(fn func() (string, error)) *EdgeClashStub {
	s.fn = fn
	return s
}

// EdgeShiftStub controls how the double answers Shift and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeShiftStub struct {
	*stub.MethodStub[EdgeShiftCall]

	fn func(string, int)
}

// Func supplies a body for Shift, for when the answer depends on the
// arguments. Nothing takes precedence over it: Shift answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *EdgeShiftStub) Func(fn func(string, int)) *EdgeShiftStub {
	s.fn = fn
	return s
}

// EdgeFetchStub controls how the double answers Fetch and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeFetchStub struct {
	*stub.MethodStub[EdgeFetchCall]

	fn       func(string) (int, error)
	fallback *EdgeFetchReturn
}

// EdgeFetchReturn holds the fixed answer configured through Returns.
type EdgeFetchReturn struct {
	A0  int
	Err error
}

// Returns pins a fixed result for every call to Fetch. A Func override
// takes precedence over it, and so does an injected fault.
func (s *EdgeFetchStub) Returns(a0 int, err error) *EdgeFetchStub {
	s.fallback = &EdgeFetchReturn{A0: a0, Err: err}
	return s
}

// Func supplies a body for Fetch, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *EdgeFetchStub) Func(fn func(string) (int, error)) *EdgeFetchStub {
	s.fn = fn
	return s
}

// EdgeWideStub controls how the double answers Wide and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgeWideStub struct {
	*stub.MethodStub[EdgeWideCall]

	fn func(int, int, int, int, int, int, int, int, int, int, int)
}

// Func supplies a body for Wide, for when the answer depends on the
// arguments. Nothing takes precedence over it: Wide answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *EdgeWideStub) Func(fn func(int, int, int, int, int, int, int, int, int, int, int)) *EdgeWideStub {
	s.fn = fn
	return s
}

// EdgePairStub controls how the double answers Pair and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type EdgePairStub struct {
	*stub.MethodStub[EdgePairCall]

	fn       func() (int, int)
	fallback *EdgePairReturn
}

// EdgePairReturn holds the fixed answer configured through Returns.
type EdgePairReturn struct {
	AB int
	Ab int
}

// Returns pins a fixed result for every call to Pair. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Pair answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *EdgePairStub) Returns(ab int, ab0 int) *EdgePairStub {
	s.fallback = &EdgePairReturn{AB: ab, Ab: ab0}
	return s
}

// Func supplies a body for Pair, for when the answer depends on the
// arguments. Nothing takes precedence over it: Pair answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *EdgePairStub) Func(fn func() (int, int)) *EdgePairStub {
	s.fn = fn
	return s
}

// --- EdgeStub ---

// The witness of the major version of the runtime of doubles.
var _ = stub.CompatV1

// EdgeStubOption configures a [EdgeStub] at construction time.
type EdgeStubOption func(*EdgeStub)

// EdgeStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func EdgeStubStrict() EdgeStubOption {
	return func(s *EdgeStub) { s.strict = true }
}

// EdgeStubDelegateTo wraps a real implementation: every method forwards to
// impl and is recorded on the double on the way through.
func EdgeStubDelegateTo(impl Edge) EdgeStubOption {
	return func(s *EdgeStub) {
		s.OnSkip.Func(impl.Skip)
		s.OnRead.Func(impl.Read)
		s.OnKeep.Func(impl.Keep)
		s.OnUse.Func(impl.Use)
		s.OnRecord.Func(impl.Record)
		s.OnClash.Func(impl.Clash)
		s.OnShift.Func(impl.Shift)
		s.OnFetch.Func(impl.Fetch)
		s.OnWide.Func(impl.Wide)
		s.OnPair.Func(impl.Pair)
	}
}

// EdgeStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func EdgeStubWithClock(clk assert.Clock) EdgeStubOption {
	return func(s *EdgeStub) { s.group.SetClock(clk) }
}

// EdgeStubWithRandSource makes probabilistic fault injection reproducible.
// A failure nobody can replay is not a finding.
func EdgeStubWithRandSource(src rand.Source) EdgeStubOption {
	return func(s *EdgeStub) { s.group.SetRandSource(src) }
}

// EdgeStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func EdgeStubBenchMode() EdgeStubOption {
	return func(s *EdgeStub) { s.group.BenchMode() }
}

// EdgeStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with EdgeStubWithClock, or d elapses in wall time.
func EdgeStubWithLatency(d time.Duration) EdgeStubOption {
	return func(s *EdgeStub) { s.group.SetLatency(d) }
}

// EdgeStubFailing arms every method that can fail with err, so a caller's
// error path is one option rather than one Faults call per method. A
// method answering no error is left alone.
func EdgeStubFailing(err error) EdgeStubOption {
	return func(s *EdgeStub) {
		s.OnRead.Faults(err, 1)
		s.OnUse.Faults(err, 1)
		s.OnRecord.Faults(err, 1)
		s.OnClash.Faults(err, 1)
		s.OnFetch.Faults(err, 1)
	}
}

// WithEdgeSkip sets Skip's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithEdgeSkip(fn func(string, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnSkip.Func(fn) }
}

// WithEdgeRead sets Read's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithEdgeRead(fn func() (int, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnRead.Func(fn) }
}

// WithEdgeKeep sets Keep's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithEdgeKeep(fn func(int) int) EdgeStubOption {
	return func(s *EdgeStub) { s.OnKeep.Func(fn) }
}

// WithEdgeUse sets Use's body at construction time, for the common case of
// configuring one method and taking the defaults for the rest.
func WithEdgeUse(fn func(string) error) EdgeStubOption {
	return func(s *EdgeStub) { s.OnUse.Func(fn) }
}

// WithEdgeRecord sets Record's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithEdgeRecord(fn func(string) (int, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnRecord.Func(fn) }
}

// WithEdgeClash sets Clash's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithEdgeClash(fn func() (string, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnClash.Func(fn) }
}

// WithEdgeShift sets Shift's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithEdgeShift(fn func(string, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnShift.Func(fn) }
}

// WithEdgeFetch sets Fetch's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithEdgeFetch(fn func(string) (int, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnFetch.Func(fn) }
}

// WithEdgeWide sets Wide's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithEdgeWide(fn func(int, int, int, int, int, int, int, int, int, int, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnWide.Func(fn) }
}

// WithEdgePair sets Pair's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithEdgePair(fn func() (int, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnPair.Func(fn) }
}

// EdgeStub is a recording test double for Edge.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type EdgeStub struct {
	OnSkip   *EdgeSkipStub
	OnRead   *EdgeReadStub
	OnKeep   *EdgeKeepStub
	OnUse    *EdgeUseStub
	OnRecord *EdgeRecordStub
	OnClash  *EdgeClashStub
	OnShift  *EdgeShiftStub
	OnFetch  *EdgeFetchStub
	OnWide   *EdgeWideStub
	OnPair   *EdgePairStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that EdgeStub satisfies Edge.
var _ Edge = (*EdgeStub)(nil)

// NewEdgeStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewEdgeStub(tb assert.TB, opts ...EdgeStubOption) *EdgeStub {
	s := &EdgeStub{
		OnSkip:   &EdgeSkipStub{MethodStub: stub.NewMethodStub[EdgeSkipCall](tb, "Edge.Skip")},
		OnRead:   &EdgeReadStub{MethodStub: stub.NewMethodStub[EdgeReadCall](tb, "Edge.Read")},
		OnKeep:   &EdgeKeepStub{MethodStub: stub.NewMethodStub[EdgeKeepCall](tb, "Edge.Keep")},
		OnUse:    &EdgeUseStub{MethodStub: stub.NewMethodStub[EdgeUseCall](tb, "Edge.Use")},
		OnRecord: &EdgeRecordStub{MethodStub: stub.NewMethodStub[EdgeRecordCall](tb, "Edge.Record")},
		OnClash:  &EdgeClashStub{MethodStub: stub.NewMethodStub[EdgeClashCall](tb, "Edge.Clash")},
		OnShift:  &EdgeShiftStub{MethodStub: stub.NewMethodStub[EdgeShiftCall](tb, "Edge.Shift")},
		OnFetch:  &EdgeFetchStub{MethodStub: stub.NewMethodStub[EdgeFetchCall](tb, "Edge.Fetch")},
		OnWide:   &EdgeWideStub{MethodStub: stub.NewMethodStub[EdgeWideCall](tb, "Edge.Wide")},
		OnPair:   &EdgePairStub{MethodStub: stub.NewMethodStub[EdgePairCall](tb, "Edge.Pair")},
	}
	s.group = stub.NewGroup(
		s.OnSkip.MethodStub,
		s.OnRead.MethodStub,
		s.OnKeep.MethodStub,
		s.OnUse.MethodStub,
		s.OnRecord.MethodStub,
		s.OnClash.MethodStub,
		s.OnShift.MethodStub,
		s.OnFetch.MethodStub,
		s.OnWide.MethodStub,
		s.OnPair.MethodStub,
	)
	for _, opt := range opts {
		opt(s)
	}
	if s.strict {
		s.group.StrictAll()
	}
	s.group.Bind(tb)
	return s
}

// NewEdgeStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewEdgeStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...EdgeStubOption) *EdgeStub {
	return NewEdgeStub(tb, append([]EdgeStubOption{
		EdgeStubWithClock(clk),
		EdgeStubWithLatency(d),
	}, opts...)...)
}

// NewEdgeStubFailing returns a double whose every fallible method answers
// err — [EdgeStubFailing] applied at construction, for the test whose
// whole subject is the caller's error path. Further options apply after
// it.
func NewEdgeStubFailing(tb assert.TB, err error, opts ...EdgeStubOption) *EdgeStub {
	return NewEdgeStub(tb, append([]EdgeStubOption{EdgeStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *EdgeStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *EdgeSkipStub) invokeVoid(a0 string, e int) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(a0, e)
		return struct{}{}
	}
}

// Skip records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Skip(a0 string, e int) {
	call := EdgeSkipCall{A0: a0, E: e}
	stub.Answer(s.OnSkip.MethodStub, &call, stub.Arms[EdgeSkipCall, struct{}]{
		Invoke: s.OnSkip.invokeVoid(a0, e),
	})
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *EdgeReadStub) invoke() func() EdgeReadReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeReadReturn {
		r0, r1 := s.fn()
		return EdgeReadReturn{R: r0, Err: r1}
	}
}

// Read records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Read() (int, error) {
	call := EdgeReadCall{}
	r := stub.Answer(s.OnRead.MethodStub, &call, stub.Arms[EdgeReadCall, EdgeReadReturn]{
		Invoke:   s.OnRead.invoke(),
		Fallback: s.OnRead.fallback,
		Fault:    func(err error) EdgeReadReturn { return EdgeReadReturn{Err: err} },
		Stamp: func(c *EdgeReadCall, r EdgeReadReturn) {
			c.R = r.R
			c.Err = r.Err
		},
	})
	return r.R, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *EdgeKeepStub) invoke(r int) func() EdgeKeepReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeKeepReturn {
		r0 := s.fn(r)
		return EdgeKeepReturn{Result: r0}
	}
}

// Keep records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Keep(r int) int {
	call := EdgeKeepCall{R: r}
	r1 := stub.Answer(s.OnKeep.MethodStub, &call, stub.Arms[EdgeKeepCall, EdgeKeepReturn]{
		Invoke:   s.OnKeep.invoke(r),
		Fallback: s.OnKeep.fallback,
		Stamp: func(c *EdgeKeepCall, r EdgeKeepReturn) {
			c.Result = r.Result
		},
	})
	return r1.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s0 *EdgeUseStub) invoke(s string) func() EdgeUseReturn {
	if s0.fn == nil {
		return nil
	}
	return func() EdgeUseReturn {
		r0 := s0.fn(s)
		return EdgeUseReturn{Err: r0}
	}
}

// Use records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s0 *EdgeStub) Use(s string) error {
	call := EdgeUseCall{S: s}
	r := stub.Answer(s0.OnUse.MethodStub, &call, stub.Arms[EdgeUseCall, EdgeUseReturn]{
		Invoke:   s0.OnUse.invoke(s),
		Fallback: s0.OnUse.fallback,
		Fault:    func(err error) EdgeUseReturn { return EdgeUseReturn{Err: err} },
		Stamp: func(c *EdgeUseCall, r EdgeUseReturn) {
			c.Err = r.Err
		},
	})
	return r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *EdgeRecordStub) invoke(call string) func() EdgeRecordReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeRecordReturn {
		n, err := s.fn(call)
		return EdgeRecordReturn{N: n, Err: err}
	}
}

// Record records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Record(call string) (n int, err error) {
	call0 := EdgeRecordCall{Call: call}
	r := stub.Answer(s.OnRecord.MethodStub, &call0, stub.Arms[EdgeRecordCall, EdgeRecordReturn]{
		Invoke:   s.OnRecord.invoke(call),
		Fallback: s.OnRecord.fallback,
		Fault:    func(err error) EdgeRecordReturn { return EdgeRecordReturn{Err: err} },
		Stamp: func(c *EdgeRecordCall, r EdgeRecordReturn) {
			c.N = r.N
			c.Err = r.Err
		},
	})
	return r.N, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *EdgeClashStub) invoke() func() EdgeClashReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeClashReturn {
		r0, r1 := s.fn()
		return EdgeClashReturn{S: r0, Err: r1}
	}
}

// Clash records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Clash() (string, error) {
	call := EdgeClashCall{}
	r := stub.Answer(s.OnClash.MethodStub, &call, stub.Arms[EdgeClashCall, EdgeClashReturn]{
		Invoke:   s.OnClash.invoke(),
		Fallback: s.OnClash.fallback,
		Fault:    func(err error) EdgeClashReturn { return EdgeClashReturn{Err: err} },
		Stamp: func(c *EdgeClashCall, r EdgeClashReturn) {
			c.S = r.S
			c.Err = r.Err
		},
	})
	return r.S, r.Err
}

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *EdgeShiftStub) invokeVoid(a00 string, a0 int) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(a00, a0)
		return struct{}{}
	}
}

// Shift records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Shift(a00 string, a0 int) {
	call := EdgeShiftCall{A00: a00, A0: a0}
	stub.Answer(s.OnShift.MethodStub, &call, stub.Arms[EdgeShiftCall, struct{}]{
		Invoke: s.OnShift.invokeVoid(a00, a0),
	})
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *EdgeFetchStub) invoke(a00 string) func() EdgeFetchReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeFetchReturn {
		a0, err := s.fn(a00)
		return EdgeFetchReturn{A0: a0, Err: err}
	}
}

// Fetch records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Fetch(a00 string) (a0 int, err error) {
	call := EdgeFetchCall{A00: a00}
	r := stub.Answer(s.OnFetch.MethodStub, &call, stub.Arms[EdgeFetchCall, EdgeFetchReturn]{
		Invoke:   s.OnFetch.invoke(a00),
		Fallback: s.OnFetch.fallback,
		Fault:    func(err error) EdgeFetchReturn { return EdgeFetchReturn{Err: err} },
		Stamp: func(c *EdgeFetchCall, r EdgeFetchReturn) {
			c.A0 = r.A0
			c.Err = r.Err
		},
	})
	return r.A0, r.Err
}

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *EdgeWideStub) invokeVoid(a int, a10 int, a1 int, b int, c int, d int, e int, f int, g int, h int, a100 int) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(a, a10, a1, b, c, d, e, f, g, h, a100)
		return struct{}{}
	}
}

// Wide records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Wide(a int, a10 int, a1 int, b int, c int, d int, e int, f int, g int, h int, a100 int) {
	call := EdgeWideCall{A: a, A10: a10, A1: a1, B: b, C: c, D: d, E: e, F: f, G: g, H: h, A100: a100}
	stub.Answer(s.OnWide.MethodStub, &call, stub.Arms[EdgeWideCall, struct{}]{
		Invoke: s.OnWide.invokeVoid(a, a10, a1, b, c, d, e, f, g, h, a100),
	})
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *EdgePairStub) invoke() func() EdgePairReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgePairReturn {
		a_b, ab := s.fn()
		return EdgePairReturn{AB: a_b, Ab: ab}
	}
}

// Pair records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *EdgeStub) Pair() (aB int, ab int) {
	call := EdgePairCall{}
	r := stub.Answer(s.OnPair.MethodStub, &call, stub.Arms[EdgePairCall, EdgePairReturn]{
		Invoke:   s.OnPair.invoke(),
		Fallback: s.OnPair.fallback,
		Stamp: func(c *EdgePairCall, r EdgePairReturn) {
			c.AB = r.AB
			c.Ab = r.Ab
		},
	})
	return r.AB, r.Ab
}
