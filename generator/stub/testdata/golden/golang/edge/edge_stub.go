package edge

import (
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// EdgeSkipCall records one call of Edge.Skip.
type EdgeSkipCall struct {
	A0 string
	E  int
}

// EdgeReadCall records one call of Edge.Read.
type EdgeReadCall struct {
	R   int
	Err error
}

// EdgeKeepCall records one call of Edge.Keep.
type EdgeKeepCall struct {
	R      int
	Result int
}

// EdgeUseCall records one call of Edge.Use.
type EdgeUseCall struct {
	S   string
	Err error
}

// EdgeRecordCall records one call of Edge.Record.
type EdgeRecordCall struct {
	Call string
	N    int
	Err  error
}

// EdgeClashCall records one call of Edge.Clash.
type EdgeClashCall struct {
	S   string
	Err error
}

// EdgeShiftCall records one call of Edge.Shift.
type EdgeShiftCall struct {
	A00 string
	A0  int
}

// EdgeFetchCall records one call of Edge.Fetch.
type EdgeFetchCall struct {
	A00 string
	A0  int
	Err error
}

// EdgeWideCall records one call of Edge.Wide.
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

// EdgePairCall records one call of Edge.Pair.
type EdgePairCall struct {
	AB int
	Ab int
}

// EdgeSkipStub configures Skip and records its calls. Skip returns no
// error, so a fault armed through it fails the test.
type EdgeSkipStub struct {
	*stub.MethodStub[EdgeSkipCall]

	fn func(string, int)
}

// Func sets the function that computes the results of Skip from its
// arguments.
func (s *EdgeSkipStub) Func(fn func(string, int)) *EdgeSkipStub {
	s.fn = fn
	return s
}

// EdgeReadStub configures Read and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type EdgeReadStub struct {
	*stub.MethodStub[EdgeReadCall]

	fn       func() (int, error)
	fallback *EdgeReadReturn
}

// EdgeReadReturn contains the fixed results of Read.
type EdgeReadReturn struct {
	R   int
	Err error
}

// Returns sets the results of every call of Read. A function and an
// injected fault take precedence over them.
func (s *EdgeReadStub) Returns(r int, err error) *EdgeReadStub {
	s.fallback = &EdgeReadReturn{R: r, Err: err}
	return s
}

// Func sets the function that computes the results of Read from its
// arguments. An injected fault takes precedence over it.
func (s *EdgeReadStub) Func(fn func() (int, error)) *EdgeReadStub {
	s.fn = fn
	return s
}

// EdgeKeepStub configures Keep and records its calls. Keep returns no
// error, so a fault armed through it fails the test.
type EdgeKeepStub struct {
	*stub.MethodStub[EdgeKeepCall]

	fn       func(int) int
	fallback *EdgeKeepReturn
}

// EdgeKeepReturn contains the fixed results of Keep.
type EdgeKeepReturn struct {
	Result int
}

// Returns sets the results of every call of Keep. A function takes
// precedence over them.
func (s *EdgeKeepStub) Returns(result int) *EdgeKeepStub {
	s.fallback = &EdgeKeepReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Keep from its
// arguments.
func (s *EdgeKeepStub) Func(fn func(int) int) *EdgeKeepStub {
	s.fn = fn
	return s
}

// EdgeUseStub configures Use and records its calls. A fault armed through
// it takes precedence over the function and the fixed results.
type EdgeUseStub struct {
	*stub.MethodStub[EdgeUseCall]

	fn       func(string) error
	fallback *EdgeUseReturn
}

// EdgeUseReturn contains the fixed results of Use.
type EdgeUseReturn struct {
	Err error
}

// Returns sets the results of every call of Use. A function and an
// injected fault take precedence over them.
func (s *EdgeUseStub) Returns(err error) *EdgeUseStub {
	s.fallback = &EdgeUseReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Use from its
// arguments. An injected fault takes precedence over it.
func (s *EdgeUseStub) Func(fn func(string) error) *EdgeUseStub {
	s.fn = fn
	return s
}

// EdgeRecordStub configures Record and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type EdgeRecordStub struct {
	*stub.MethodStub[EdgeRecordCall]

	fn       func(string) (int, error)
	fallback *EdgeRecordReturn
}

// EdgeRecordReturn contains the fixed results of Record.
type EdgeRecordReturn struct {
	N   int
	Err error
}

// Returns sets the results of every call of Record. A function and an
// injected fault take precedence over them.
func (s *EdgeRecordStub) Returns(n int, err error) *EdgeRecordStub {
	s.fallback = &EdgeRecordReturn{N: n, Err: err}
	return s
}

// Func sets the function that computes the results of Record from its
// arguments. An injected fault takes precedence over it.
func (s *EdgeRecordStub) Func(fn func(string) (int, error)) *EdgeRecordStub {
	s.fn = fn
	return s
}

// EdgeClashStub configures Clash and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type EdgeClashStub struct {
	*stub.MethodStub[EdgeClashCall]

	fn       func() (string, error)
	fallback *EdgeClashReturn
}

// EdgeClashReturn contains the fixed results of Clash.
type EdgeClashReturn struct {
	S   string
	Err error
}

// Returns sets the results of every call of Clash. A function and an
// injected fault take precedence over them.
func (s *EdgeClashStub) Returns(s0 string, err error) *EdgeClashStub {
	s.fallback = &EdgeClashReturn{S: s0, Err: err}
	return s
}

// Func sets the function that computes the results of Clash from its
// arguments. An injected fault takes precedence over it.
func (s *EdgeClashStub) Func(fn func() (string, error)) *EdgeClashStub {
	s.fn = fn
	return s
}

// EdgeShiftStub configures Shift and records its calls. Shift returns no
// error, so a fault armed through it fails the test.
type EdgeShiftStub struct {
	*stub.MethodStub[EdgeShiftCall]

	fn func(string, int)
}

// Func sets the function that computes the results of Shift from its
// arguments.
func (s *EdgeShiftStub) Func(fn func(string, int)) *EdgeShiftStub {
	s.fn = fn
	return s
}

// EdgeFetchStub configures Fetch and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type EdgeFetchStub struct {
	*stub.MethodStub[EdgeFetchCall]

	fn       func(string) (int, error)
	fallback *EdgeFetchReturn
}

// EdgeFetchReturn contains the fixed results of Fetch.
type EdgeFetchReturn struct {
	A0  int
	Err error
}

// Returns sets the results of every call of Fetch. A function and an
// injected fault take precedence over them.
func (s *EdgeFetchStub) Returns(a0 int, err error) *EdgeFetchStub {
	s.fallback = &EdgeFetchReturn{A0: a0, Err: err}
	return s
}

// Func sets the function that computes the results of Fetch from its
// arguments. An injected fault takes precedence over it.
func (s *EdgeFetchStub) Func(fn func(string) (int, error)) *EdgeFetchStub {
	s.fn = fn
	return s
}

// EdgeWideStub configures Wide and records its calls. Wide returns no
// error, so a fault armed through it fails the test.
type EdgeWideStub struct {
	*stub.MethodStub[EdgeWideCall]

	fn func(int, int, int, int, int, int, int, int, int, int, int)
}

// Func sets the function that computes the results of Wide from its
// arguments.
func (s *EdgeWideStub) Func(
	fn func(int, int, int, int, int, int, int, int, int, int, int),
) *EdgeWideStub {
	s.fn = fn
	return s
}

// EdgePairStub configures Pair and records its calls. Pair returns no
// error, so a fault armed through it fails the test.
type EdgePairStub struct {
	*stub.MethodStub[EdgePairCall]

	fn       func() (int, int)
	fallback *EdgePairReturn
}

// EdgePairReturn contains the fixed results of Pair.
type EdgePairReturn struct {
	AB int
	Ab int
}

// Returns sets the results of every call of Pair. A function takes
// precedence over them.
func (s *EdgePairStub) Returns(ab, ab0 int) *EdgePairStub {
	s.fallback = &EdgePairReturn{AB: ab, Ab: ab0}
	return s
}

// Func sets the function that computes the results of Pair from its
// arguments.
func (s *EdgePairStub) Func(fn func() (int, int)) *EdgePairStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// EdgeStubOption configures a [EdgeStub] at its construction.
type EdgeStubOption func(*EdgeStub)

// EdgeStubStrict returns the option that makes an unconfigured method fail
// the test instead of returning its zero value.
func EdgeStubStrict() EdgeStubOption {
	return func(s *EdgeStub) { s.strict = true }
}

// EdgeStubDelegateTo returns the option that passes each call to impl. The
// double still records the call.
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

// EdgeStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func EdgeStubWithClock(clk assert.Clock) EdgeStubOption {
	return func(s *EdgeStub) { s.group.SetClock(clk) }
}

// EdgeStubWithRandSource returns the option that sets the random source of
// the probabilistic faults to src. A seeded source repeats the faults of a
// run.
func EdgeStubWithRandSource(src rand.Source) EdgeStubOption {
	return func(s *EdgeStub) { s.group.SetRandSource(src) }
}

// EdgeStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func EdgeStubBenchMode() EdgeStubOption {
	return func(s *EdgeStub) { s.group.BenchMode() }
}

// EdgeStubWithLatency returns the option that makes every method sleep d
// before it returns. Without EdgeStubWithClock, d passes in real time.
func EdgeStubWithLatency(d time.Duration) EdgeStubOption {
	return func(s *EdgeStub) { s.group.SetLatency(d) }
}

// EdgeStubFailing returns the option that arms err on each method that can
// fail. Each call of such a method then fails with err.
func EdgeStubFailing(err error) EdgeStubOption {
	return func(s *EdgeStub) {
		s.OnRead.Faults(err, 1)
		s.OnUse.Faults(err, 1)
		s.OnRecord.Faults(err, 1)
		s.OnClash.Faults(err, 1)
		s.OnFetch.Faults(err, 1)
	}
}

// WithEdgeSkip returns the option that sets the function of Skip.
func WithEdgeSkip(fn func(string, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnSkip.Func(fn) }
}

// WithEdgeRead returns the option that sets the function of Read.
func WithEdgeRead(fn func() (int, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnRead.Func(fn) }
}

// WithEdgeKeep returns the option that sets the function of Keep.
func WithEdgeKeep(fn func(int) int) EdgeStubOption {
	return func(s *EdgeStub) { s.OnKeep.Func(fn) }
}

// WithEdgeUse returns the option that sets the function of Use.
func WithEdgeUse(fn func(string) error) EdgeStubOption {
	return func(s *EdgeStub) { s.OnUse.Func(fn) }
}

// WithEdgeRecord returns the option that sets the function of Record.
func WithEdgeRecord(fn func(string) (int, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnRecord.Func(fn) }
}

// WithEdgeClash returns the option that sets the function of Clash.
func WithEdgeClash(fn func() (string, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnClash.Func(fn) }
}

// WithEdgeShift returns the option that sets the function of Shift.
func WithEdgeShift(fn func(string, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnShift.Func(fn) }
}

// WithEdgeFetch returns the option that sets the function of Fetch.
func WithEdgeFetch(fn func(string) (int, error)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnFetch.Func(fn) }
}

// WithEdgeWide returns the option that sets the function of Wide.
func WithEdgeWide(fn func(int, int, int, int, int, int, int, int, int, int, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnWide.Func(fn) }
}

// WithEdgePair returns the option that sets the function of Pair.
func WithEdgePair(fn func() (int, int)) EdgeStubOption {
	return func(s *EdgeStub) { s.OnPair.Func(fn) }
}

// EdgeStub is a recording test double of Edge. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
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

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that EdgeStub implements Edge.
var _ Edge = (*EdgeStub)(nil)

// NewEdgeStub returns a double of Edge that reports its failures to tb and
// applies opts. The double checks the expected counts of calls at the end
// of the test. With a nil tb, the double reports nothing.
func NewEdgeStub(tb assert.TB, opts ...EdgeStubOption) *EdgeStub {
	s := &EdgeStub{
		OnSkip: &EdgeSkipStub{MethodStub: stub.NewMethodStub[EdgeSkipCall](tb, "Edge.Skip")},
		OnRead: &EdgeReadStub{MethodStub: stub.NewMethodStub[EdgeReadCall](tb, "Edge.Read")},
		OnKeep: &EdgeKeepStub{MethodStub: stub.NewMethodStub[EdgeKeepCall](tb, "Edge.Keep")},
		OnUse:  &EdgeUseStub{MethodStub: stub.NewMethodStub[EdgeUseCall](tb, "Edge.Use")},
		OnRecord: &EdgeRecordStub{
			MethodStub: stub.NewMethodStub[EdgeRecordCall](tb, "Edge.Record"),
		},
		OnClash: &EdgeClashStub{MethodStub: stub.NewMethodStub[EdgeClashCall](tb, "Edge.Clash")},
		OnShift: &EdgeShiftStub{MethodStub: stub.NewMethodStub[EdgeShiftCall](tb, "Edge.Shift")},
		OnFetch: &EdgeFetchStub{MethodStub: stub.NewMethodStub[EdgeFetchCall](tb, "Edge.Fetch")},
		OnWide:  &EdgeWideStub{MethodStub: stub.NewMethodStub[EdgeWideCall](tb, "Edge.Wide")},
		OnPair:  &EdgePairStub{MethodStub: stub.NewMethodStub[EdgePairCall](tb, "Edge.Pair")},
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

// NewEdgeStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewEdgeStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...EdgeStubOption,
) *EdgeStub {
	return NewEdgeStub(tb, append([]EdgeStubOption{
		EdgeStubWithClock(clk),
		EdgeStubWithLatency(d),
	}, opts...)...)
}

// NewEdgeStubFailing returns a double on which each call of a method that
// can fail fails with err. It applies opts after the failure.
func NewEdgeStubFailing(tb assert.TB, err error, opts ...EdgeStubOption) *EdgeStub {
	return NewEdgeStub(tb, append([]EdgeStubOption{EdgeStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *EdgeStub) ResetCalls() {
	s.group.Reset()
}

// invokeVoid returns a call of the function of Func with the arguments, and
// nil where no function is set. The empty struct is the result of a method
// without results.
func (s *EdgeSkipStub) invokeVoid(a0 string, e int) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(a0, e)
		return struct{}{}
	}
}

// Skip records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *EdgeStub) Skip(a0 string, e int) {
	call := EdgeSkipCall{A0: a0, E: e}
	stub.Answer(s.OnSkip.MethodStub, &call, stub.Arms[EdgeSkipCall, struct{}]{
		Invoke: s.OnSkip.invokeVoid(a0, e),
	})
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *EdgeReadStub) invoke() func() EdgeReadReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeReadReturn {
		r0, r1 := s.fn()
		return EdgeReadReturn{R: r0, Err: r1}
	}
}

// Read records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *EdgeKeepStub) invoke(r int) func() EdgeKeepReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeKeepReturn {
		r0 := s.fn(r)
		return EdgeKeepReturn{Result: r0}
	}
}

// Keep records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s0 *EdgeUseStub) invoke(s string) func() EdgeUseReturn {
	if s0.fn == nil {
		return nil
	}
	return func() EdgeUseReturn {
		r0 := s0.fn(s)
		return EdgeUseReturn{Err: r0}
	}
}

// Use records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *EdgeRecordStub) invoke(call string) func() EdgeRecordReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeRecordReturn {
		n, err := s.fn(call)
		return EdgeRecordReturn{N: n, Err: err}
	}
}

// Record records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *EdgeClashStub) invoke() func() EdgeClashReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeClashReturn {
		r0, r1 := s.fn()
		return EdgeClashReturn{S: r0, Err: r1}
	}
}

// Clash records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invokeVoid returns a call of the function of Func with the arguments, and
// nil where no function is set. The empty struct is the result of a method
// without results.
func (s *EdgeShiftStub) invokeVoid(a00 string, a0 int) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(a00, a0)
		return struct{}{}
	}
}

// Shift records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *EdgeStub) Shift(a00 string, a0 int) {
	call := EdgeShiftCall{A00: a00, A0: a0}
	stub.Answer(s.OnShift.MethodStub, &call, stub.Arms[EdgeShiftCall, struct{}]{
		Invoke: s.OnShift.invokeVoid(a00, a0),
	})
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *EdgeFetchStub) invoke(a00 string) func() EdgeFetchReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgeFetchReturn {
		a0, err := s.fn(a00)
		return EdgeFetchReturn{A0: a0, Err: err}
	}
}

// Fetch records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invokeVoid returns a call of the function of Func with the arguments, and
// nil where no function is set. The empty struct is the result of a method
// without results.
func (s *EdgeWideStub) invokeVoid(a, a10, a1, b, c, d, e, f, g, h, a100 int) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(a, a10, a1, b, c, d, e, f, g, h, a100)
		return struct{}{}
	}
}

// Wide records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *EdgeStub) Wide(a, a10, a1, b, c, d, e, f, g, h, a100 int) {
	call := EdgeWideCall{
		A:    a,
		A10:  a10,
		A1:   a1,
		B:    b,
		C:    c,
		D:    d,
		E:    e,
		F:    f,
		G:    g,
		H:    h,
		A100: a100,
	}
	stub.Answer(s.OnWide.MethodStub, &call, stub.Arms[EdgeWideCall, struct{}]{
		Invoke: s.OnWide.invokeVoid(a, a10, a1, b, c, d, e, f, g, h, a100),
	})
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *EdgePairStub) invoke() func() EdgePairReturn {
	if s.fn == nil {
		return nil
	}
	return func() EdgePairReturn {
		a_b, ab := s.fn()
		return EdgePairReturn{AB: a_b, Ab: ab}
	}
}

// Pair records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *EdgeStub) Pair() (aB, ab int) {
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
