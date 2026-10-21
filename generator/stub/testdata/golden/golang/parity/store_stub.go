package parity

import (
	"context"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/fault"
	"go.dokimi.dev/stub"
)

// StoreGetCall records one invocation of Store.Get.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreGetCall struct {
	Ctx  context.Context
	ID   string
	Item string
	Err  error
}

// StoreListCall records one invocation of Store.List.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreListCall struct {
	Result []string
	Err    error
}

// StoreCloseCall records one invocation of Store.Close.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreCloseCall struct {
}

// --- Per-method configuration ---

// StoreGetStub controls how the double answers Get and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreGetStub struct {
	*stub.MethodStub[StoreGetCall]

	fn       func(context.Context, string) (string, error)
	fallback *StoreGetReturn
}

// StoreGetReturn holds the fixed answer configured through Returns.
type StoreGetReturn struct {
	Item string
	Err  error
}

// Returns pins a fixed result for every call to Get. A Func override takes
// precedence over it, and so does an injected fault.
func (s *StoreGetStub) Returns(item string, err error) *StoreGetStub {
	s.fallback = &StoreGetReturn{Item: item, Err: err}
	return s
}

// Func supplies a body for Get, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreGetStub) Func(fn func(context.Context, string) (string, error)) *StoreGetStub {
	s.fn = fn
	return s
}

// FaultNotFound makes the next Get call fail with ErrNotFound; the calls
// after it succeed.
//
// One-shot, because that is what a test asserting on this error's handling
// wants: poison a single call, then let the scenario continue. Faults is
// the periodic form for the rest.
func (s *StoreGetStub) FaultNotFound() *StoreGetStub {
	s.FaultsOnce(ErrNotFound)
	return s
}

// FaultGone makes the next Get call fail with ErrGone; the calls after it
// succeed.
//
// One-shot, because that is what a test asserting on this error's handling
// wants: poison a single call, then let the scenario continue. Faults is
// the periodic form for the rest.
func (s *StoreGetStub) FaultGone() *StoreGetStub {
	s.FaultsOnce(ErrGone)
	return s
}

// FaultForPartition makes Get fail every nth call whose ID matches key,
// leaving other partitions untouched.
//
// Partitioned failure is what proves isolation: a test asserting that one
// tenant's outage does not reach another needs the fault to be selective.
func (s *StoreGetStub) FaultForPartition(key string, err error, n int) *StoreGetStub {
	s.FaultsWhen(func(c StoreGetCall) bool { return c.ID == key }, err, n)
	return s
}

// FaultForOtherPartitions inverts FaultForPartition: every nth call whose
// ID does not match key fails, and the named partition is left healthy.
func (s *StoreGetStub) FaultForOtherPartitions(key string, err error, n int) *StoreGetStub {
	s.FaultsWhen(func(c StoreGetCall) bool { return c.ID != key }, err, n)
	return s
}

// StoreListStub controls how the double answers List and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreListStub struct {
	*stub.MethodStub[StoreListCall]

	fn       func() ([]string, error)
	fallback *StoreListReturn
}

// StoreListReturn holds the fixed answer configured through Returns.
type StoreListReturn struct {
	Result []string
	Err    error
}

// Returns pins a fixed result for every call to List. A Func override
// takes precedence over it, and so does an injected fault.
func (s *StoreListStub) Returns(result []string, err error) *StoreListStub {
	s.fallback = &StoreListReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for List, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreListStub) Func(fn func() ([]string, error)) *StoreListStub {
	s.fn = fn
	return s
}

// RetrySchedule makes List fail with err on every call before call number
// 3, which succeeds along with every call after it.
func (s *StoreListStub) RetrySchedule(err error) *StoreListStub {
	s.SetFault(fault.NewRetryFault[StoreListCall](err, 3))
	return s
}

// StoreCloseStub controls how the double answers Close and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreCloseStub struct {
	*stub.MethodStub[StoreCloseCall]

	fn func()
}

// Func supplies a body for Close, for when the answer depends on the
// arguments. Nothing takes precedence over it: Close answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *StoreCloseStub) Func(fn func()) *StoreCloseStub {
	s.fn = fn
	return s
}

// --- StoreStub ---

// The witness of the major version of the runtime of doubles.
var _ = stub.CompatV1

// The witness of the major version of the fault strategies.
var _ = fault.CompatV1

// StoreStubOption configures a [StoreStub] at construction time.
type StoreStubOption func(*StoreStub)

// StoreStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func StoreStubStrict() StoreStubOption {
	return func(s *StoreStub) { s.strict = true }
}

// StoreStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func StoreStubDelegateTo(impl Store) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Func(impl.Get)
		s.OnList.Func(impl.List)
		s.OnClose.Func(impl.Close)
	}
}

// StoreStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func StoreStubWithClock(clk assert.Clock) StoreStubOption {
	return func(s *StoreStub) { s.group.SetClock(clk) }
}

// StoreStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func StoreStubWithRandSource(src rand.Source) StoreStubOption {
	return func(s *StoreStub) { s.group.SetRandSource(src) }
}

// StoreStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func StoreStubBenchMode() StoreStubOption {
	return func(s *StoreStub) { s.group.BenchMode() }
}

// StoreStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with StoreStubWithClock, or d elapses in wall time.
func StoreStubWithLatency(d time.Duration) StoreStubOption {
	return func(s *StoreStub) { s.group.SetLatency(d) }
}

// StoreStubFailing arms every method that can fail with err, so a caller's
// error path is one option rather than one Faults call per method. A
// method answering no error is left alone.
func StoreStubFailing(err error) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Faults(err, 1)
		s.OnList.Faults(err, 1)
	}
}

// WithStoreGet sets Get's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreGet(fn func(context.Context, string) (string, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnGet.Func(fn) }
}

// WithStoreList sets List's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreList(fn func() ([]string, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnList.Func(fn) }
}

// WithStoreClose sets Close's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithStoreClose(fn func()) StoreStubOption {
	return func(s *StoreStub) { s.OnClose.Func(fn) }
}

// StoreStub is a recording test double for Store.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type StoreStub struct {
	OnGet   *StoreGetStub
	OnList  *StoreListStub
	OnClose *StoreCloseStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that StoreStub satisfies Store.
var _ Store = (*StoreStub)(nil)

// NewStoreStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewStoreStub(tb assert.TB, opts ...StoreStubOption) *StoreStub {
	s := &StoreStub{
		OnGet:   &StoreGetStub{MethodStub: stub.NewMethodStub[StoreGetCall](tb, "Store.Get")},
		OnList:  &StoreListStub{MethodStub: stub.NewMethodStub[StoreListCall](tb, "Store.List")},
		OnClose: &StoreCloseStub{MethodStub: stub.NewMethodStub[StoreCloseCall](tb, "Store.Close")},
	}
	s.group = stub.NewGroup(
		s.OnGet.MethodStub,
		s.OnList.MethodStub,
		s.OnClose.MethodStub,
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

// NewStoreStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewStoreStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...StoreStubOption) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{
		StoreStubWithClock(clk),
		StoreStubWithLatency(d),
	}, opts...)...)
}

// NewStoreStubFailing returns a double whose every fallible method answers
// err — [StoreStubFailing] applied at construction, for the test whose
// whole subject is the caller's error path. Further options apply after
// it.
func NewStoreStubFailing(tb assert.TB, err error, opts ...StoreStubOption) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{StoreStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *StoreStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreGetStub) invoke(ctx context.Context, id string) func() StoreGetReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreGetReturn {
		item, err := s.fn(ctx, id)
		return StoreGetReturn{Item: item, Err: err}
	}
}

// Get records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Get(ctx context.Context, id string) (item string, err error) {
	call := StoreGetCall{Ctx: ctx, ID: id}
	r := stub.Answer(s.OnGet.MethodStub, &call, stub.Arms[StoreGetCall, StoreGetReturn]{
		Invoke:   s.OnGet.invoke(ctx, id),
		Fallback: s.OnGet.fallback,
		Fault:    func(err error) StoreGetReturn { return StoreGetReturn{Err: err} },
		Stamp: func(c *StoreGetCall, r StoreGetReturn) {
			c.Item = r.Item
			c.Err = r.Err
		},
	})
	return r.Item, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreListStub) invoke() func() StoreListReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreListReturn {
		r0, r1 := s.fn()
		return StoreListReturn{Result: r0, Err: r1}
	}
}

// List records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) List() ([]string, error) {
	call := StoreListCall{}
	r := stub.Answer(s.OnList.MethodStub, &call, stub.Arms[StoreListCall, StoreListReturn]{
		Invoke:   s.OnList.invoke(),
		Fallback: s.OnList.fallback,
		Fault:    func(err error) StoreListReturn { return StoreListReturn{Err: err} },
		Stamp: func(c *StoreListCall, r StoreListReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *StoreCloseStub) invokeVoid() func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn()
		return struct{}{}
	}
}

// Close records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Close() {
	call := StoreCloseCall{}
	stub.Answer(s.OnClose.MethodStub, &call, stub.Arms[StoreCloseCall, struct{}]{
		Invoke: s.OnClose.invokeVoid(),
	})
}
