package built

import (
	"context"
	rand "math/rand/v2"
	"time"

	"example.com/shop/billing"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// StoreGetCall records one call of Store.Get.
type StoreGetCall struct {
	Ctx    context.Context
	ID     string
	Result Account
	Err    error
}

// StoreFindCall records one call of Store.Find.
type StoreFindCall struct {
	Ctx    context.Context
	ID     string
	Result *Account
	Err    error
}

// StoreLookupCall records one call of Store.Lookup.
type StoreLookupCall struct {
	ID      string
	Result0 Account
	Result1 bool
}

// StoreListCall records one call of Store.List.
type StoreListCall struct {
	Ctx    context.Context
	Result Page[Account]
	Err    error
}

// StoreNoteCall records one call of Store.Note.
type StoreNoteCall struct {
	ID     string
	Result Note
}

// StoreBillCall records one call of Store.Bill.
type StoreBillCall struct {
	ID     string
	Result billing.Invoice
}

// StoreGetStub configures Get and records its calls. A fault armed through
// it takes precedence over the function and the fixed results.
type StoreGetStub struct {
	*stub.MethodStub[StoreGetCall]

	fn       func(context.Context, string) (Account, error)
	fallback *StoreGetReturn
}

// StoreGetReturn contains the fixed results of Get.
type StoreGetReturn struct {
	Result Account
	Err    error
}

// Returns sets the results of every call of Get. A function and an
// injected fault take precedence over them.
func (s *StoreGetStub) Returns(result Account, err error) *StoreGetStub {
	s.fallback = &StoreGetReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Get from its
// arguments. An injected fault takes precedence over it.
func (s *StoreGetStub) Func(fn func(context.Context, string) (Account, error)) *StoreGetStub {
	s.fn = fn
	return s
}

// StoreFindStub configures Find and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type StoreFindStub struct {
	*stub.MethodStub[StoreFindCall]

	fn       func(context.Context, string) (*Account, error)
	fallback *StoreFindReturn
}

// StoreFindReturn contains the fixed results of Find.
type StoreFindReturn struct {
	Result *Account
	Err    error
}

// Returns sets the results of every call of Find. A function and an
// injected fault take precedence over them.
func (s *StoreFindStub) Returns(result *Account, err error) *StoreFindStub {
	s.fallback = &StoreFindReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Find from its
// arguments. An injected fault takes precedence over it.
func (s *StoreFindStub) Func(fn func(context.Context, string) (*Account, error)) *StoreFindStub {
	s.fn = fn
	return s
}

// StoreLookupStub configures Lookup and records its calls. Lookup returns
// no error, so a fault armed through it fails the test.
type StoreLookupStub struct {
	*stub.MethodStub[StoreLookupCall]

	fn       func(string) (Account, bool)
	fallback *StoreLookupReturn
}

// StoreLookupReturn contains the fixed results of Lookup.
type StoreLookupReturn struct {
	Result0 Account
	Result1 bool
}

// Returns sets the results of every call of Lookup. A function takes
// precedence over them.
func (s *StoreLookupStub) Returns(result0 Account, result1 bool) *StoreLookupStub {
	s.fallback = &StoreLookupReturn{Result0: result0, Result1: result1}
	return s
}

// Func sets the function that computes the results of Lookup from its
// arguments.
func (s *StoreLookupStub) Func(fn func(string) (Account, bool)) *StoreLookupStub {
	s.fn = fn
	return s
}

// StoreListStub configures List and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type StoreListStub struct {
	*stub.MethodStub[StoreListCall]

	fn       func(context.Context) (Page[Account], error)
	fallback *StoreListReturn
}

// StoreListReturn contains the fixed results of List.
type StoreListReturn struct {
	Result Page[Account]
	Err    error
}

// Returns sets the results of every call of List. A function and an
// injected fault take precedence over them.
func (s *StoreListStub) Returns(result Page[Account], err error) *StoreListStub {
	s.fallback = &StoreListReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of List from its
// arguments. An injected fault takes precedence over it.
func (s *StoreListStub) Func(fn func(context.Context) (Page[Account], error)) *StoreListStub {
	s.fn = fn
	return s
}

// StoreNoteStub configures Note and records its calls. Note returns no
// error, so a fault armed through it fails the test.
type StoreNoteStub struct {
	*stub.MethodStub[StoreNoteCall]

	fn       func(string) Note
	fallback *StoreNoteReturn
}

// StoreNoteReturn contains the fixed results of Note.
type StoreNoteReturn struct {
	Result Note
}

// Returns sets the results of every call of Note. A function takes
// precedence over them.
func (s *StoreNoteStub) Returns(result Note) *StoreNoteStub {
	s.fallback = &StoreNoteReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Note from its
// arguments.
func (s *StoreNoteStub) Func(fn func(string) Note) *StoreNoteStub {
	s.fn = fn
	return s
}

// StoreBillStub configures Bill and records its calls. Bill returns no
// error, so a fault armed through it fails the test.
type StoreBillStub struct {
	*stub.MethodStub[StoreBillCall]

	fn       func(string) billing.Invoice
	fallback *StoreBillReturn
}

// StoreBillReturn contains the fixed results of Bill.
type StoreBillReturn struct {
	Result billing.Invoice
}

// Returns sets the results of every call of Bill. A function takes
// precedence over them.
func (s *StoreBillStub) Returns(result billing.Invoice) *StoreBillStub {
	s.fallback = &StoreBillReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Bill from its
// arguments.
func (s *StoreBillStub) Func(fn func(string) billing.Invoice) *StoreBillStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// StoreStubOption configures a [StoreStub] at its construction.
type StoreStubOption func(*StoreStub)

// StoreStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its default results.
func StoreStubStrict() StoreStubOption {
	return func(s *StoreStub) { s.strict = true }
}

// StoreStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
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

// StoreStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func StoreStubWithClock(clk assert.Clock) StoreStubOption {
	return func(s *StoreStub) { s.group.SetClock(clk) }
}

// StoreStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func StoreStubWithRandSource(src rand.Source) StoreStubOption {
	return func(s *StoreStub) { s.group.SetRandSource(src) }
}

// StoreStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func StoreStubBenchMode() StoreStubOption {
	return func(s *StoreStub) { s.group.BenchMode() }
}

// StoreStubWithLatency returns the option that makes every method sleep d
// before it returns. Without StoreStubWithClock, d passes in real time.
func StoreStubWithLatency(d time.Duration) StoreStubOption {
	return func(s *StoreStub) { s.group.SetLatency(d) }
}

// StoreStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func StoreStubFailing(err error) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Faults(err, 1)
		s.OnFind.Faults(err, 1)
		s.OnList.Faults(err, 1)
	}
}

// WithStoreGet returns the option that sets the function of Get.
func WithStoreGet(fn func(context.Context, string) (Account, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnGet.Func(fn) }
}

// WithStoreFind returns the option that sets the function of Find.
func WithStoreFind(fn func(context.Context, string) (*Account, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnFind.Func(fn) }
}

// WithStoreLookup returns the option that sets the function of Lookup.
func WithStoreLookup(fn func(string) (Account, bool)) StoreStubOption {
	return func(s *StoreStub) { s.OnLookup.Func(fn) }
}

// WithStoreList returns the option that sets the function of List.
func WithStoreList(fn func(context.Context) (Page[Account], error)) StoreStubOption {
	return func(s *StoreStub) { s.OnList.Func(fn) }
}

// WithStoreNote returns the option that sets the function of Note.
func WithStoreNote(fn func(string) Note) StoreStubOption {
	return func(s *StoreStub) { s.OnNote.Func(fn) }
}

// WithStoreBill returns the option that sets the function of Bill.
func WithStoreBill(fn func(string) billing.Invoice) StoreStubOption {
	return func(s *StoreStub) { s.OnBill.Func(fn) }
}

// StoreStub is a recording test double of Store. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns its default results. The default result of a type with a builder
// is a new value from the builder, and every other result is its zero
// value.
type StoreStub struct {
	OnGet    *StoreGetStub
	OnFind   *StoreFindStub
	OnLookup *StoreLookupStub
	OnList   *StoreListStub
	OnNote   *StoreNoteStub
	OnBill   *StoreBillStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that StoreStub implements Store.
var _ Store = (*StoreStub)(nil)

// NewStoreStub returns a double of Store that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewStoreStub(tb assert.TB, opts ...StoreStubOption) *StoreStub {
	s := &StoreStub{
		OnGet:  &StoreGetStub{MethodStub: stub.NewMethodStub[StoreGetCall](tb, "Store.Get")},
		OnFind: &StoreFindStub{MethodStub: stub.NewMethodStub[StoreFindCall](tb, "Store.Find")},
		OnLookup: &StoreLookupStub{
			MethodStub: stub.NewMethodStub[StoreLookupCall](tb, "Store.Lookup"),
		},
		OnList: &StoreListStub{MethodStub: stub.NewMethodStub[StoreListCall](tb, "Store.List")},
		OnNote: &StoreNoteStub{MethodStub: stub.NewMethodStub[StoreNoteCall](tb, "Store.Note")},
		OnBill: &StoreBillStub{MethodStub: stub.NewMethodStub[StoreBillCall](tb, "Store.Bill")},
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

// NewStoreStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewStoreStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...StoreStubOption,
) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{
		StoreStubWithClock(clk),
		StoreStubWithLatency(d),
	}, opts...)...)
}

// NewStoreStubFailing returns a double on which each call of a method that
// can fail fails with err. It applies opts after the failure.
func NewStoreStubFailing(tb assert.TB, err error, opts ...StoreStubOption) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{StoreStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *StoreStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreGetStub) invoke(ctx context.Context, id string) func() StoreGetReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreGetReturn {
		r0, r1 := s.fn(ctx, id)
		return StoreGetReturn{Result: r0, Err: r1}
	}
}

// Get records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the default results.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreFindStub) invoke(ctx context.Context, id string) func() StoreFindReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreFindReturn {
		r0, r1 := s.fn(ctx, id)
		return StoreFindReturn{Result: r0, Err: r1}
	}
}

// Find records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the default results.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreLookupStub) invoke(id string) func() StoreLookupReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreLookupReturn {
		r0, r1 := s.fn(id)
		return StoreLookupReturn{Result0: r0, Result1: r1}
	}
}

// Lookup records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the default results.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreListStub) invoke(ctx context.Context) func() StoreListReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreListReturn {
		r0, r1 := s.fn(ctx)
		return StoreListReturn{Result: r0, Err: r1}
	}
}

// List records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the default results.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreNoteStub) invoke(id string) func() StoreNoteReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreNoteReturn {
		r0 := s.fn(id)
		return StoreNoteReturn{Result: r0}
	}
}

// Note records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreBillStub) invoke(id string) func() StoreBillReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreBillReturn {
		r0 := s.fn(id)
		return StoreBillReturn{Result: r0}
	}
}

// Bill records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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
