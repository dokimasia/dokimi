package built

import (
	"context"
	rand "math/rand/v2"
	"time"

	"example.com/shop/billing"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/dokimi-go/stub"
)

// StoreGetCall records one invocation of Store.Get.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreGetCall struct {
	Ctx    context.Context
	ID     string
	Result Account
	Err    error
}

// StoreFindCall records one invocation of Store.Find.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreFindCall struct {
	Ctx    context.Context
	ID     string
	Result *Account
	Err    error
}

// StoreLookupCall records one invocation of Store.Lookup.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreLookupCall struct {
	ID      string
	Result0 Account
	Result1 bool
}

// StoreListCall records one invocation of Store.List.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreListCall struct {
	Ctx    context.Context
	Result Page[Account]
	Err    error
}

// StoreNoteCall records one invocation of Store.Note.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreNoteCall struct {
	ID     string
	Result Note
}

// StoreBillCall records one invocation of Store.Bill.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreBillCall struct {
	ID     string
	Result billing.Invoice
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

	fn       func(context.Context, string) (Account, error)
	fallback *StoreGetReturn
}

// StoreGetReturn holds the fixed answer configured through Returns.
type StoreGetReturn struct {
	Result Account
	Err    error
}

// Returns pins a fixed result for every call to Get. A Func override takes
// precedence over it, and so does an injected fault.
func (s *StoreGetStub) Returns(result Account, err error) *StoreGetStub {
	s.fallback = &StoreGetReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for Get, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreGetStub) Func(fn func(context.Context, string) (Account, error)) *StoreGetStub {
	s.fn = fn
	return s
}

// StoreFindStub controls how the double answers Find and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreFindStub struct {
	*stub.MethodStub[StoreFindCall]

	fn       func(context.Context, string) (*Account, error)
	fallback *StoreFindReturn
}

// StoreFindReturn holds the fixed answer configured through Returns.
type StoreFindReturn struct {
	Result *Account
	Err    error
}

// Returns pins a fixed result for every call to Find. A Func override
// takes precedence over it, and so does an injected fault.
func (s *StoreFindStub) Returns(result *Account, err error) *StoreFindStub {
	s.fallback = &StoreFindReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for Find, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreFindStub) Func(fn func(context.Context, string) (*Account, error)) *StoreFindStub {
	s.fn = fn
	return s
}

// StoreLookupStub controls how the double answers Lookup and records what
// it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreLookupStub struct {
	*stub.MethodStub[StoreLookupCall]

	fn       func(string) (Account, bool)
	fallback *StoreLookupReturn
}

// StoreLookupReturn holds the fixed answer configured through Returns.
type StoreLookupReturn struct {
	Result0 Account
	Result1 bool
}

// Returns pins a fixed result for every call to Lookup. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Lookup answers no error, so the
// double is built with no fault arm and a fault armed through Faults fails
// the test.
func (s *StoreLookupStub) Returns(result0 Account, result1 bool) *StoreLookupStub {
	s.fallback = &StoreLookupReturn{Result0: result0, Result1: result1}
	return s
}

// Func supplies a body for Lookup, for when the answer depends on the
// arguments. Nothing takes precedence over it: Lookup answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *StoreLookupStub) Func(fn func(string) (Account, bool)) *StoreLookupStub {
	s.fn = fn
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

	fn       func(context.Context) (Page[Account], error)
	fallback *StoreListReturn
}

// StoreListReturn holds the fixed answer configured through Returns.
type StoreListReturn struct {
	Result Page[Account]
	Err    error
}

// Returns pins a fixed result for every call to List. A Func override
// takes precedence over it, and so does an injected fault.
func (s *StoreListStub) Returns(result Page[Account], err error) *StoreListStub {
	s.fallback = &StoreListReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for List, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreListStub) Func(fn func(context.Context) (Page[Account], error)) *StoreListStub {
	s.fn = fn
	return s
}

// StoreNoteStub controls how the double answers Note and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreNoteStub struct {
	*stub.MethodStub[StoreNoteCall]

	fn       func(string) Note
	fallback *StoreNoteReturn
}

// StoreNoteReturn holds the fixed answer configured through Returns.
type StoreNoteReturn struct {
	Result Note
}

// Returns pins a fixed result for every call to Note. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Note answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *StoreNoteStub) Returns(result Note) *StoreNoteStub {
	s.fallback = &StoreNoteReturn{Result: result}
	return s
}

// Func supplies a body for Note, for when the answer depends on the
// arguments. Nothing takes precedence over it: Note answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *StoreNoteStub) Func(fn func(string) Note) *StoreNoteStub {
	s.fn = fn
	return s
}

// StoreBillStub controls how the double answers Bill and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreBillStub struct {
	*stub.MethodStub[StoreBillCall]

	fn       func(string) billing.Invoice
	fallback *StoreBillReturn
}

// StoreBillReturn holds the fixed answer configured through Returns.
type StoreBillReturn struct {
	Result billing.Invoice
}

// Returns pins a fixed result for every call to Bill. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Bill answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *StoreBillStub) Returns(result billing.Invoice) *StoreBillStub {
	s.fallback = &StoreBillReturn{Result: result}
	return s
}

// Func supplies a body for Bill, for when the answer depends on the
// arguments. Nothing takes precedence over it: Bill answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *StoreBillStub) Func(fn func(string) billing.Invoice) *StoreBillStub {
	s.fn = fn
	return s
}

// --- StoreStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// StoreStubOption configures a [StoreStub] at construction time.
type StoreStubOption func(*StoreStub)

// StoreStubStrict makes an unconfigured method fail the test rather than
// return its default results.
func StoreStubStrict() StoreStubOption {
	return func(s *StoreStub) { s.strict = true }
}

// StoreStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func StoreStubDelegateTo(impl Store) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Func(impl.Get)
		s.OnFind.Func(impl.Find)
		s.OnLookup.Func(impl.Lookup)
		s.OnList.Func(impl.List)
		s.OnNote.Func(impl.Note)
		s.OnBill.Func(impl.Bill)
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
		s.OnFind.Faults(err, 1)
		s.OnList.Faults(err, 1)
	}
}

// WithStoreGet sets Get's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreGet(fn func(context.Context, string) (Account, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnGet.Func(fn) }
}

// WithStoreFind sets Find's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreFind(fn func(context.Context, string) (*Account, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnFind.Func(fn) }
}

// WithStoreLookup sets Lookup's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithStoreLookup(fn func(string) (Account, bool)) StoreStubOption {
	return func(s *StoreStub) { s.OnLookup.Func(fn) }
}

// WithStoreList sets List's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreList(fn func(context.Context) (Page[Account], error)) StoreStubOption {
	return func(s *StoreStub) { s.OnList.Func(fn) }
}

// WithStoreNote sets Note's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreNote(fn func(string) Note) StoreStubOption {
	return func(s *StoreStub) { s.OnNote.Func(fn) }
}

// WithStoreBill sets Bill's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreBill(fn func(string) billing.Invoice) StoreStubOption {
	return func(s *StoreStub) { s.OnBill.Func(fn) }
}

// StoreStub is a recording test double for Store.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method records the call and returns its default results: a new value
// from the builder of each result type that has one, and the zero value of
// each other result.
type StoreStub struct {
	OnGet    *StoreGetStub
	OnFind   *StoreFindStub
	OnLookup *StoreLookupStub
	OnList   *StoreListStub
	OnNote   *StoreNoteStub
	OnBill   *StoreBillStub

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
		OnGet:    &StoreGetStub{MethodStub: stub.NewMethodStub[StoreGetCall](tb, "Store.Get")},
		OnFind:   &StoreFindStub{MethodStub: stub.NewMethodStub[StoreFindCall](tb, "Store.Find")},
		OnLookup: &StoreLookupStub{MethodStub: stub.NewMethodStub[StoreLookupCall](tb, "Store.Lookup")},
		OnList:   &StoreListStub{MethodStub: stub.NewMethodStub[StoreListCall](tb, "Store.List")},
		OnNote:   &StoreNoteStub{MethodStub: stub.NewMethodStub[StoreNoteCall](tb, "Store.Note")},
		OnBill:   &StoreBillStub{MethodStub: stub.NewMethodStub[StoreBillCall](tb, "Store.Bill")},
	}
	s.group = stub.NewGroup(
		s.OnGet.MethodStub,
		s.OnFind.MethodStub,
		s.OnLookup.MethodStub,
		s.OnList.MethodStub,
		s.OnNote.MethodStub,
		s.OnBill.MethodStub,
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
		r0, r1 := s.fn(ctx, id)
		return StoreGetReturn{Result: r0, Err: r1}
	}
}

// Get records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the default results — is [stub.Answer]'s to decide, so every generated
// double resolves a call the same way and the ordering is tested once
// rather than restated per method.
func (s *StoreStub) Get(ctx context.Context, id string) (Account, error) {
	call := StoreGetCall{Ctx: ctx, ID: id}
	r := stub.Answer(s.OnGet.MethodStub, &call, stub.Arms[StoreGetCall, StoreGetReturn]{
		Invoke:   s.OnGet.invoke(ctx, id),
		Fallback: s.OnGet.fallback,
		Fault:    func(err error) StoreGetReturn { return StoreGetReturn{Err: err} },
		Default: func() StoreGetReturn {
			r0 := NewAccountBuilder().Build()
			return StoreGetReturn{Result: r0}
		},
		Stamp: func(c *StoreGetCall, r StoreGetReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreFindStub) invoke(ctx context.Context, id string) func() StoreFindReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreFindReturn {
		r0, r1 := s.fn(ctx, id)
		return StoreFindReturn{Result: r0, Err: r1}
	}
}

// Find records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the default results — is [stub.Answer]'s to decide, so every generated
// double resolves a call the same way and the ordering is tested once
// rather than restated per method.
func (s *StoreStub) Find(ctx context.Context, id string) (*Account, error) {
	call := StoreFindCall{Ctx: ctx, ID: id}
	r := stub.Answer(s.OnFind.MethodStub, &call, stub.Arms[StoreFindCall, StoreFindReturn]{
		Invoke:   s.OnFind.invoke(ctx, id),
		Fallback: s.OnFind.fallback,
		Fault:    func(err error) StoreFindReturn { return StoreFindReturn{Err: err} },
		Default: func() StoreFindReturn {
			r0 := NewAccountBuilder().Build()
			return StoreFindReturn{Result: &r0}
		},
		Stamp: func(c *StoreFindCall, r StoreFindReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreLookupStub) invoke(id string) func() StoreLookupReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreLookupReturn {
		r0, r1 := s.fn(id)
		return StoreLookupReturn{Result0: r0, Result1: r1}
	}
}

// Lookup records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the default results — is [stub.Answer]'s to decide, so every generated
// double resolves a call the same way and the ordering is tested once
// rather than restated per method.
func (s *StoreStub) Lookup(id string) (Account, bool) {
	call := StoreLookupCall{ID: id}
	r := stub.Answer(s.OnLookup.MethodStub, &call, stub.Arms[StoreLookupCall, StoreLookupReturn]{
		Invoke:   s.OnLookup.invoke(id),
		Fallback: s.OnLookup.fallback,
		Default: func() StoreLookupReturn {
			r0 := NewAccountBuilder().Build()
			return StoreLookupReturn{Result0: r0}
		},
		Stamp: func(c *StoreLookupCall, r StoreLookupReturn) {
			c.Result0 = r.Result0
			c.Result1 = r.Result1
		},
	})
	return r.Result0, r.Result1
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreListStub) invoke(ctx context.Context) func() StoreListReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreListReturn {
		r0, r1 := s.fn(ctx)
		return StoreListReturn{Result: r0, Err: r1}
	}
}

// List records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the default results — is [stub.Answer]'s to decide, so every generated
// double resolves a call the same way and the ordering is tested once
// rather than restated per method.
func (s *StoreStub) List(ctx context.Context) (Page[Account], error) {
	call := StoreListCall{Ctx: ctx}
	r := stub.Answer(s.OnList.MethodStub, &call, stub.Arms[StoreListCall, StoreListReturn]{
		Invoke:   s.OnList.invoke(ctx),
		Fallback: s.OnList.fallback,
		Fault:    func(err error) StoreListReturn { return StoreListReturn{Err: err} },
		Default: func() StoreListReturn {
			r0 := NewPageBuilder[Account]().Build()
			return StoreListReturn{Result: r0}
		},
		Stamp: func(c *StoreListCall, r StoreListReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreNoteStub) invoke(id string) func() StoreNoteReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreNoteReturn {
		r0 := s.fn(id)
		return StoreNoteReturn{Result: r0}
	}
}

// Note records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Note(id string) Note {
	call := StoreNoteCall{ID: id}
	r := stub.Answer(s.OnNote.MethodStub, &call, stub.Arms[StoreNoteCall, StoreNoteReturn]{
		Invoke:   s.OnNote.invoke(id),
		Fallback: s.OnNote.fallback,
		Stamp: func(c *StoreNoteCall, r StoreNoteReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreBillStub) invoke(id string) func() StoreBillReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreBillReturn {
		r0 := s.fn(id)
		return StoreBillReturn{Result: r0}
	}
}

// Bill records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Bill(id string) billing.Invoice {
	call := StoreBillCall{ID: id}
	r := stub.Answer(s.OnBill.MethodStub, &call, stub.Arms[StoreBillCall, StoreBillReturn]{
		Invoke:   s.OnBill.invoke(id),
		Fallback: s.OnBill.fallback,
		Stamp: func(c *StoreBillCall, r StoreBillReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}
