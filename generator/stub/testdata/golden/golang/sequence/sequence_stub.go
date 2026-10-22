package sequence

import (
	"iter"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// TableNamesCall records one call of Table.Names.
type TableNamesCall struct {
	Result iter.Seq[string]
}

// TableRowsCall records one call of Table.Rows.
type TableRowsCall struct {
	Limit  int
	Result iter.Seq2[Item, error]
}

// TableCountsCall records one call of Table.Counts.
type TableCountsCall struct {
	Result iter.Seq2[string, int]
}

// TableTailCall records one call of Table.Tail.
type TableTailCall struct {
	Result Lines
}

// TableNamesStub configures Names and records its calls. Names returns no
// error, so a fault armed through it fails the test.
type TableNamesStub struct {
	*stub.MethodStub[TableNamesCall]

	fn       func() iter.Seq[string]
	fallback *TableNamesReturn
}

// TableNamesReturn contains the fixed results of Names.
type TableNamesReturn struct {
	Result iter.Seq[string]
}

// Returns sets the results of every call of Names. A function takes
// precedence over them.
func (s *TableNamesStub) Returns(result iter.Seq[string]) *TableNamesStub {
	s.fallback = &TableNamesReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Names from its
// arguments.
func (s *TableNamesStub) Func(fn func() iter.Seq[string]) *TableNamesStub {
	s.fn = fn
	return s
}

// Yields makes Names return a sequence of items. The sequence stops where
// the caller breaks out of the loop.
func (s *TableNamesStub) Yields(items ...string) *TableNamesStub {
	s.fn = func() iter.Seq[string] {
		return func(yield func(string) bool) {
			for _, v := range items {
				if !yield(v) {
					return
				}
			}
		}
	}
	return s
}

// TableRowsStub configures Rows and records its calls. Rows returns no
// error, so a fault armed through it fails the test.
type TableRowsStub struct {
	*stub.MethodStub[TableRowsCall]

	fn       func(int) iter.Seq2[Item, error]
	fallback *TableRowsReturn
}

// TableRowsReturn contains the fixed results of Rows.
type TableRowsReturn struct {
	Result iter.Seq2[Item, error]
}

// Returns sets the results of every call of Rows. A function takes
// precedence over them.
func (s *TableRowsStub) Returns(result iter.Seq2[Item, error]) *TableRowsStub {
	s.fallback = &TableRowsReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Rows from its
// arguments.
func (s *TableRowsStub) Func(fn func(int) iter.Seq2[Item, error]) *TableRowsStub {
	s.fn = fn
	return s
}

// Yields makes Rows return a sequence of items. The sequence stops where
// the caller breaks out of the loop.
func (s *TableRowsStub) Yields(items ...Item) *TableRowsStub {
	s.fn = func(_ int) iter.Seq2[Item, error] {
		return func(yield func(Item, error) bool) {
			for _, v := range items {
				if !yield(v, nil) {
					return
				}
			}
		}
	}
	return s
}

// YieldsError makes Rows return a sequence of items that fails with err
// after the last item.
func (s *TableRowsStub) YieldsError(items []Item, err error) *TableRowsStub {
	s.fn = func(_ int) iter.Seq2[Item, error] {
		return func(yield func(Item, error) bool) {
			for _, v := range items {
				if !yield(v, nil) {
					return
				}
			}
			var zero Item
			yield(zero, err)
		}
	}
	return s
}

// TableCountsStub configures Counts and records its calls. Counts returns
// no error, so a fault armed through it fails the test.
type TableCountsStub struct {
	*stub.MethodStub[TableCountsCall]

	fn       func() iter.Seq2[string, int]
	fallback *TableCountsReturn
}

// TableCountsReturn contains the fixed results of Counts.
type TableCountsReturn struct {
	Result iter.Seq2[string, int]
}

// Returns sets the results of every call of Counts. A function takes
// precedence over them.
func (s *TableCountsStub) Returns(result iter.Seq2[string, int]) *TableCountsStub {
	s.fallback = &TableCountsReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Counts from its
// arguments.
func (s *TableCountsStub) Func(fn func() iter.Seq2[string, int]) *TableCountsStub {
	s.fn = fn
	return s
}

// Yields makes Counts return a sequence of items. The sequence stops where
// the caller breaks out of the loop.
func (s *TableCountsStub) Yields(items ...string) *TableCountsStub {
	s.fn = func() iter.Seq2[string, int] {
		return func(yield func(string, int) bool) {
			var second int
			for _, v := range items {
				if !yield(v, second) {
					return
				}
			}
		}
	}
	return s
}

// TableTailStub configures Tail and records its calls. Tail returns no
// error, so a fault armed through it fails the test.
type TableTailStub struct {
	*stub.MethodStub[TableTailCall]

	fn       func() Lines
	fallback *TableTailReturn
}

// TableTailReturn contains the fixed results of Tail.
type TableTailReturn struct {
	Result Lines
}

// Returns sets the results of every call of Tail. A function takes
// precedence over them.
func (s *TableTailStub) Returns(result Lines) *TableTailStub {
	s.fallback = &TableTailReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Tail from its
// arguments.
func (s *TableTailStub) Func(fn func() Lines) *TableTailStub {
	s.fn = fn
	return s
}

// Yields makes Tail return a sequence of items. The sequence stops where
// the caller breaks out of the loop.
func (s *TableTailStub) Yields(items ...string) *TableTailStub {
	s.fn = func() Lines {
		return func(yield func(string) bool) {
			for _, v := range items {
				if !yield(v) {
					return
				}
			}
		}
	}
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// TableStubOption configures a [TableStub] at its construction.
type TableStubOption func(*TableStub)

// TableStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func TableStubStrict() TableStubOption {
	return func(s *TableStub) { s.strict = true }
}

// TableStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func TableStubDelegateTo(impl Table) TableStubOption {
	return func(s *TableStub) {
		s.OnNames.Func(impl.Names)
		s.OnRows.Func(impl.Rows)
		s.OnCounts.Func(impl.Counts)
		s.OnTail.Func(impl.Tail)
	}
}

// TableStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func TableStubWithClock(clk assert.Clock) TableStubOption {
	return func(s *TableStub) { s.group.SetClock(clk) }
}

// TableStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func TableStubWithRandSource(src rand.Source) TableStubOption {
	return func(s *TableStub) { s.group.SetRandSource(src) }
}

// TableStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func TableStubBenchMode() TableStubOption {
	return func(s *TableStub) { s.group.BenchMode() }
}

// TableStubWithLatency returns the option that makes every method sleep d
// before it returns. Without TableStubWithClock, d passes in real time.
func TableStubWithLatency(d time.Duration) TableStubOption {
	return func(s *TableStub) { s.group.SetLatency(d) }
}

// WithTableNames returns the option that sets the function of Names.
func WithTableNames(fn func() iter.Seq[string]) TableStubOption {
	return func(s *TableStub) { s.OnNames.Func(fn) }
}

// WithTableRows returns the option that sets the function of Rows.
func WithTableRows(fn func(int) iter.Seq2[Item, error]) TableStubOption {
	return func(s *TableStub) { s.OnRows.Func(fn) }
}

// WithTableCounts returns the option that sets the function of Counts.
func WithTableCounts(fn func() iter.Seq2[string, int]) TableStubOption {
	return func(s *TableStub) { s.OnCounts.Func(fn) }
}

// WithTableTail returns the option that sets the function of Tail.
func WithTableTail(fn func() Lines) TableStubOption {
	return func(s *TableStub) { s.OnTail.Func(fn) }
}

// TableStub is a recording test double of Table. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type TableStub struct {
	OnNames  *TableNamesStub
	OnRows   *TableRowsStub
	OnCounts *TableCountsStub
	OnTail   *TableTailStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that TableStub implements Table.
var _ Table = (*TableStub)(nil)

// NewTableStub returns a double of Table that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewTableStub(tb assert.TB, opts ...TableStubOption) *TableStub {
	s := &TableStub{
		OnNames: &TableNamesStub{
			MethodStub: stub.NewMethodStub[TableNamesCall](tb, "Table.Names"),
		},
		OnRows: &TableRowsStub{MethodStub: stub.NewMethodStub[TableRowsCall](tb, "Table.Rows")},
		OnCounts: &TableCountsStub{
			MethodStub: stub.NewMethodStub[TableCountsCall](tb, "Table.Counts"),
		},
		OnTail: &TableTailStub{MethodStub: stub.NewMethodStub[TableTailCall](tb, "Table.Tail")},
	}
	s.group = stub.NewGroup(
		s.OnNames.MethodStub,
		s.OnRows.MethodStub,
		s.OnCounts.MethodStub,
		s.OnTail.MethodStub,
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

// NewTableStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewTableStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...TableStubOption,
) *TableStub {
	return NewTableStub(tb, append([]TableStubOption{
		TableStubWithClock(clk),
		TableStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *TableStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *TableNamesStub) invoke() func() TableNamesReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableNamesReturn {
		r0 := s.fn()
		return TableNamesReturn{Result: r0}
	}
}

// Names records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *TableStub) Names() iter.Seq[string] {
	call := TableNamesCall{}
	r := stub.Answer(s.OnNames.MethodStub, &call, stub.Arms[TableNamesCall, TableNamesReturn]{
		Invoke:   s.OnNames.invoke(),
		Fallback: s.OnNames.fallback,
		Stamp: func(c *TableNamesCall, r TableNamesReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *TableRowsStub) invoke(limit int) func() TableRowsReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableRowsReturn {
		r0 := s.fn(limit)
		return TableRowsReturn{Result: r0}
	}
}

// Rows records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *TableStub) Rows(limit int) iter.Seq2[Item, error] {
	call := TableRowsCall{Limit: limit}
	r := stub.Answer(s.OnRows.MethodStub, &call, stub.Arms[TableRowsCall, TableRowsReturn]{
		Invoke:   s.OnRows.invoke(limit),
		Fallback: s.OnRows.fallback,
		Stamp: func(c *TableRowsCall, r TableRowsReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *TableCountsStub) invoke() func() TableCountsReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableCountsReturn {
		r0 := s.fn()
		return TableCountsReturn{Result: r0}
	}
}

// Counts records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *TableStub) Counts() iter.Seq2[string, int] {
	call := TableCountsCall{}
	r := stub.Answer(s.OnCounts.MethodStub, &call, stub.Arms[TableCountsCall, TableCountsReturn]{
		Invoke:   s.OnCounts.invoke(),
		Fallback: s.OnCounts.fallback,
		Stamp: func(c *TableCountsCall, r TableCountsReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *TableTailStub) invoke() func() TableTailReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableTailReturn {
		r0 := s.fn()
		return TableTailReturn{Result: r0}
	}
}

// Tail records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *TableStub) Tail() Lines {
	call := TableTailCall{}
	r := stub.Answer(s.OnTail.MethodStub, &call, stub.Arms[TableTailCall, TableTailReturn]{
		Invoke:   s.OnTail.invoke(),
		Fallback: s.OnTail.fallback,
		Stamp: func(c *TableTailCall, r TableTailReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}
