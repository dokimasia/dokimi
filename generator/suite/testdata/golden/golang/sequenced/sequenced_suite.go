package sequenced

import (
	"context"
	"iter"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Counter:
//
// Counter.Inc/smoke: Inc returns without a panic for derived inputs.
//
// Counter.Inc/cancel: Inc returns an error for a cancelled context.
//
// Counter.Inc/keeps-within-its-bounds: Inc keeps the use that Value reads
// within its bounds, over any sequence of calls.
//
// Counter.Value/smoke: Value returns without a panic for derived inputs.
//
// Counter.Value/cancel: Value returns an error for a cancelled context.
//
// Counter.Value/zero-on-error: Value returns the zero value beside an
// error.
//
// Counter.Value/monotonic: Value never returns less than it returned
// before.
//
// Counter.Value/never-goes-back: Value never returns less than a result
// that it returned before, over any sequence of calls.
//
// Counter.Grow/smoke: Grow returns without a panic for derived inputs.
//
// Counter.Grow/cancel: Grow returns an error for a cancelled context.
//
// Counter.Grow/keeps-within-its-bounds: Grow keeps the use that Size reads
// within its bounds, over any sequence of calls.
//
// Counter.Size/smoke: Size returns without a panic for derived inputs.
//
// Counter.Size/cancel: Size returns an error for a cancelled context.
//
// Counter.Size/zero-on-error: Size returns the zero value beside an error.
//
// Counter.Throttle/smoke: Throttle returns without a panic for derived
// inputs.
//
// Counter.Throttle/cancel: Throttle returns an error for a cancelled
// context.
//
// Counter.Load/smoke: Load returns without a panic for derived inputs.
//
// Counter.Load/cancel: Load returns an error for a cancelled context.
//
// Counter.Load/zero-on-error: Load returns the zero value beside an error.

// The derived checks of Counter that witness their laws at the shallow
// end:
//
// Counter.Inc/smoke
//
// Counter.Value/smoke
//
// Counter.Value/monotonic
//
// Counter.Grow/smoke
//
// Counter.Size/smoke
//
// Counter.Throttle/smoke
//
// Counter.Load/smoke

// The refused claims of Counter:
//
// Counter.Throttle/keeps-within-its-bounds: the reader that observe names
// for Throttle returns no number and no length. The remedy is to name with
// observe a callable that returns the use of Throttle as a number or a
// list.

// The unclaimed classifications of Counter:
//
// The shape:writer of Counter.Inc has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Counter.Grow has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The shape:lifecycle of Counter.Throttle has the stance silent. A
// lifecycle call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Counter writes.
type (
	CounterHandle = conform.Handle[Counter]
	CounterCheck  = conform.Check[Counter]
	CounterOption = conform.Option[Counter]
)

// CounterChecks contains the name of each callable of Counter, with the
// identities of its checks.
var CounterChecks = struct {
	Inc      counterInc
	Value    counterValue
	Grow     counterGrow
	Size     counterSize
	Throttle counterThrottle
	Load     counterLoad
}{
	Inc:      counterInc{"Counter.Inc"},
	Value:    counterValue{"Counter.Value"},
	Grow:     counterGrow{"Counter.Grow"},
	Size:     counterSize{"Counter.Size"},
	Throttle: counterThrottle{"Counter.Throttle"},
	Load:     counterLoad{"Counter.Load"},
}

// counterInc is the name of Counter.Inc, with the identities of its
// checks.
type counterInc struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Inc returns without a panic
// for derived inputs.
func (c counterInc) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Inc returns an error for a
// cancelled context.
func (c counterInc) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// KeepsWithinItsBounds returns the identity of the check that Inc keeps
// the use that Value reads within its bounds, over any sequence of calls.
func (c counterInc) KeepsWithinItsBounds() conform.ID[Counter] {
	return c.Check("keeps-within-its-bounds")
}

// counterValue is the name of Counter.Value, with the identities of its
// checks.
type counterValue struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Value returns without a
// panic for derived inputs.
func (c counterValue) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Value returns an error for
// a cancelled context.
func (c counterValue) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Value returns the
// zero value beside an error.
func (c counterValue) ZeroOnError() conform.ID[Counter] {
	return c.Check("zero-on-error")
}

// Monotonic returns the identity of the check that Value never returns
// less than it returned before.
func (c counterValue) Monotonic() conform.ID[Counter] {
	return c.Check("monotonic")
}

// NeverGoesBack returns the identity of the check that Value never returns
// less than a result that it returned before, over any sequence of calls.
func (c counterValue) NeverGoesBack() conform.ID[Counter] {
	return c.Check("never-goes-back")
}

// counterGrow is the name of Counter.Grow, with the identities of its
// checks.
type counterGrow struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Grow returns without a
// panic for derived inputs.
func (c counterGrow) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Grow returns an error for
// a cancelled context.
func (c counterGrow) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// KeepsWithinItsBounds returns the identity of the check that Grow keeps
// the use that Size reads within its bounds, over any sequence of calls.
func (c counterGrow) KeepsWithinItsBounds() conform.ID[Counter] {
	return c.Check("keeps-within-its-bounds")
}

// counterSize is the name of Counter.Size, with the identities of its
// checks.
type counterSize struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Size returns without a
// panic for derived inputs.
func (c counterSize) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Size returns an error for
// a cancelled context.
func (c counterSize) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Size returns the zero
// value beside an error.
func (c counterSize) ZeroOnError() conform.ID[Counter] {
	return c.Check("zero-on-error")
}

// counterThrottle is the name of Counter.Throttle, with the identities of
// its checks.
type counterThrottle struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Throttle returns without a
// panic for derived inputs.
func (c counterThrottle) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Throttle returns an error
// for a cancelled context.
func (c counterThrottle) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// counterLoad is the name of Counter.Load, with the identities of its
// checks.
type counterLoad struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Load returns without a
// panic for derived inputs.
func (c counterLoad) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Load returns an error for
// a cancelled context.
func (c counterLoad) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Load returns the zero
// value beside an error.
func (c counterLoad) ZeroOnError() conform.ID[Counter] {
	return c.Check("zero-on-error")
}

// CounterSuite returns the checks of Counter, and the calls and the
// budgets of its callables.
func CounterSuite() conform.Suite[Counter] {
	return conform.Suite[Counter]{
		Declaration: "Counter",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Counter]{
			{
				ID:          CounterChecks.Inc.Smoke(),
				Statement:   "Inc returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Inc(ctx, Option(func(*int) {}))
					})
				},
				Defect: counterIncPanics,
			},
			{
				ID:          CounterChecks.Inc.Cancel(),
				Statement:   "Inc returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Inc(ctx, Option(func(*int) {}))
					})
				},
				Defect: counterIncIgnoresCancel,
			},
			{
				ID:          CounterChecks.Inc.KeepsWithinItsBounds(),
				Statement:   "Inc keeps the use that Value reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CounterHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						counterActions,
						func(ctx context.Context) (int, error) {
							values, err := counterValueCall(ctx, h.Subject(), []any{})
							v, _ := values[0].(int)
							return int(v), err
						},
						0,
						3,
					)
				},
				Defect: counterValueGrowsPastTheLimitAfterThree,
			},
			{
				ID:          CounterChecks.Value.Smoke(),
				Statement:   "Value returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx)
						return err
					})
				},
				Defect: counterValuePanics,
			},
			{
				ID:          CounterChecks.Value.Cancel(),
				Statement:   "Value returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx)
						return err
					})
				},
				Defect: counterValueIgnoresCancel,
			},
			{
				ID:          CounterChecks.Value.ZeroOnError(),
				Statement:   "Value returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Value(ctx)
					}, int(0))
				},
				Defect: counterValueValueBeside,
			},
			{
				ID:          CounterChecks.Value.Monotonic(),
				Statement:   "Value never returns less than it returned before",
				Claim:       spec.Monotonic,
				Probe:       spec.NeverDecreases,
				Reason:      spec.Regressed,
				Law:         spec.ResultsNeverDecrease,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.NeverDecreases(h, func(ctx context.Context) (int, error) {
						return h.Subject().Value(ctx)
					})
				},
				Defect: counterValueDescends,
			},
			{
				ID:          CounterChecks.Value.NeverGoesBack(),
				Statement:   "Value never returns less than a result that it returned before, over any sequence of calls",
				Claim:       spec.NeverGoesBack,
				Probe:       spec.ResultsNeverFall,
				Reason:      spec.Regressed,
				Law:         spec.ResultsNeverDecrease,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Counter]{CounterChecks.Value.Monotonic()},
				Property: func(h *CounterHandle, c *prop.Case) {
					probe.ResultsNeverFall[Counter, int](h, c, counterActions, "Value")
				},
				Defect: counterValueGoesBackAfterThree,
			},
			{
				ID:          CounterChecks.Grow.Smoke(),
				Statement:   "Grow returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Grow(ctx)
					})
				},
				Defect: counterGrowPanics,
			},
			{
				ID:          CounterChecks.Grow.Cancel(),
				Statement:   "Grow returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Grow(ctx)
					})
				},
				Defect: counterGrowIgnoresCancel,
			},
			{
				ID:          CounterChecks.Grow.KeepsWithinItsBounds(),
				Statement:   "Grow keeps the use that Size reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CounterHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						counterActions,
						func(ctx context.Context) (int, error) {
							values, err := counterSizeCall(ctx, h.Subject(), []any{})
							v, _ := values[0].(uint64)
							return int(v), err
						},
						0,
						64,
					)
				},
				Defect: counterSizeGrowsPastTheLimitAfterThree,
			},
			{
				ID:          CounterChecks.Size.Smoke(),
				Statement:   "Size returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Size(ctx)
						return err
					})
				},
				Defect: counterSizePanics,
			},
			{
				ID:          CounterChecks.Size.Cancel(),
				Statement:   "Size returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Size(ctx)
						return err
					})
				},
				Defect: counterSizeIgnoresCancel,
			},
			{
				ID:          CounterChecks.Size.ZeroOnError(),
				Statement:   "Size returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (uint64, error) {
						return h.Subject().Size(ctx)
					}, uint64(0))
				},
				Defect: counterSizeValueBeside,
			},
			{
				ID:          CounterChecks.Throttle.Smoke(),
				Statement:   "Throttle returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Throttle(ctx)
					})
				},
				Defect: counterThrottlePanics,
			},
			{
				ID:          CounterChecks.Throttle.Cancel(),
				Statement:   "Throttle returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Throttle(ctx)
					})
				},
				Defect: counterThrottleIgnoresCancel,
			},
			{
				ID:          CounterChecks.Load.Smoke(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx)
						return err
					})
				},
				Defect: counterLoadPanics,
			},
			{
				ID:          CounterChecks.Load.Cancel(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx)
						return err
					})
				},
				Defect: counterLoadIgnoresCancel,
			},
			{
				ID:          CounterChecks.Load.ZeroOnError(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (float64, error) {
						return h.Subject().Load(ctx)
					}, float64(0))
				},
				Defect: counterLoadValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.KeepsWithinItsBounds,
				Callable: "Counter.Throttle",
				Gate:     spec.CountableUse,
				Why:      "the reader that observe names for Throttle returns no number and no length",
				Remedy:   "name with observe a callable that returns the use of Throttle as a number or a list",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Counter.Inc",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Counter.Grow",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Counter.Throttle",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Counter]{
			{
				Calls:    counterIncCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Inc.Callable,
			},
			{
				Calls:    counterValueCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Value.Callable,
			},
			{
				Calls:    counterGrowCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Grow.Callable,
			},
			{
				Calls:    counterSizeCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Size.Callable,
			},
			{
				Calls:    counterThrottleCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Throttle.Callable,
			},
			{
				Calls:    counterLoadCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Load.Callable,
			},
		},
	}
}

// RunCounterSuite applies the checks of Counter, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCounterSuite[T Counter](t *testing.T, newSubject func() T, opts ...CounterOption) {
	t.Helper()
	RunCounterSuiteFromHandle(t, func(*CounterHandle) T { return newSubject() }, opts...)
}

// RunCounterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCounterSuiteFromHandle[T Counter](
	t *testing.T,
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	t.Helper()
	conform.Run(t, CounterSuite(), counterSubjects(newSubject), opts...)
}

// FuzzCounterSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCounterSuite[T Counter](
	f *testing.F,
	callable conform.Callable[Counter],
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	f.Helper()
	conform.Fuzz(f, CounterSuite(), callable, counterSubjects(newSubject), opts...)
}

// RunCounterBudgets applies the allocation ceilings of Counter and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCounterBudgets[T Counter](
	t *testing.T,
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	t.Helper()
	conform.RunBudgets(t, CounterSuite(), counterSubjects(newSubject), opts...)
}

// BenchCounterSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchCounterSuite[T Counter](
	b conform.B,
	callable conform.Callable[Counter],
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	b.Helper()
	conform.Bench(b, CounterSuite(), callable, counterSubjects(newSubject), opts...)
}

// counterSubjects returns the factory of the subjects of Counter that
// calls newSubject with the handle of each check.
func counterSubjects[T Counter](
	newSubject func(h *CounterHandle) T,
) func(h *CounterHandle) Counter {
	return func(h *CounterHandle) Counter { return newSubject(h) }
}

// counterIncCalls generates the calls of Counter.Inc with the inputs that
// each case draws.
var counterIncCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		return subject.Inc(ctx)
	}
})

// counterValueCalls generates the calls of Counter.Value with the inputs
// that each case draws.
var counterValueCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		_, err := subject.Value(ctx)
		return err
	}
})

// counterGrowCalls generates the calls of Counter.Grow with the inputs
// that each case draws.
var counterGrowCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		return subject.Grow(ctx)
	}
})

// counterSizeCalls generates the calls of Counter.Size with the inputs
// that each case draws.
var counterSizeCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		_, err := subject.Size(ctx)
		return err
	}
})

// counterThrottleCalls generates the calls of Counter.Throttle with the
// inputs that each case draws.
var counterThrottleCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		return subject.Throttle(ctx)
	}
})

// counterLoadCalls generates the calls of Counter.Load with the inputs
// that each case draws.
var counterLoadCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		_, err := subject.Load(ctx)
		return err
	}
})

// counterOverEmbedded is Counter under the name of the field that
// counterOver embeds.
type counterOverEmbedded = Counter

// counterOver is a subject of Counter whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type counterOver struct {
	counterOverEmbedded
	spend func()
}

// Inc calls spend and returns its zero results.
func (o counterOver) Inc(context.Context, ...Option) (r0 error) {
	o.spend()
	return r0
}

// Value calls spend and returns its zero results.
func (o counterOver) Value(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// Grow calls spend and returns its zero results.
func (o counterOver) Grow(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Size calls spend and returns its zero results.
func (o counterOver) Size(context.Context) (r0 uint64, r1 error) {
	o.spend()
	return r0, r1
}

// Throttle calls spend and returns its zero results.
func (o counterOver) Throttle(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Load calls spend and returns its zero results.
func (o counterOver) Load(context.Context) (r0 float64, r1 error) {
	o.spend()
	return r0, r1
}

// newCounterOver returns the planted defect of a ceiling of Counter, whose
// callables call spend.
func newCounterOver(_ *CounterHandle, spend func()) Counter {
	return counterOver{spend: spend}
}

// counterIncPanicsDefect is the planted defect of a Counter whose Inc
// panics.
type counterIncPanicsDefect struct {
	counterOverEmbedded
}

// Inc is the override of the planted defect.
func (*counterIncPanicsDefect) Inc(p0 context.Context, p1 ...Option) (r0 error) {
	panic("the planted defect of a Counter whose Inc panics")
}

// counterIncPanics returns the planted defect of a Counter whose Inc
// panics over a fresh subject.
func counterIncPanics(h *CounterHandle) Counter {
	return &counterIncPanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterIncIgnoresCancelDefect is the planted defect of a Counter whose
// Inc ignores the context.
type counterIncIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Inc is the override of the planted defect.
func (d *counterIncIgnoresCancelDefect) Inc(p0 context.Context, p1 ...Option) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.counterOverEmbedded.Inc(p0, p1...)
}

// counterIncIgnoresCancel returns the planted defect of a Counter whose
// Inc ignores the context over a fresh subject.
func counterIncIgnoresCancel(h *CounterHandle) Counter {
	return &counterIncIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Counter whose Value reports a use above the limit from its fourth call
// on.
type counterValueGrowsPastTheLimitAfterThreeDefect struct {
	counterOverEmbedded
	calls int
}

// Value is the override of the planted defect.
func (d *counterValueGrowsPastTheLimitAfterThreeDefect) Value(
	p0 context.Context,
) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.counterOverEmbedded.Value(p0)
	d.calls++
	if d.calls > 3 {
		r0 = 3 + 1
	}
	return r0, r1
}

// counterValueGrowsPastTheLimitAfterThree returns the planted defect of a
// Counter whose Value reports a use above the limit from its fourth call
// on over a fresh subject.
func counterValueGrowsPastTheLimitAfterThree(h *CounterHandle) Counter {
	return &counterValueGrowsPastTheLimitAfterThreeDefect{counterOverEmbedded: h.Fresh()}
}

// counterValuePanicsDefect is the planted defect of a Counter whose Value
// panics.
type counterValuePanicsDefect struct {
	counterOverEmbedded
}

// Value is the override of the planted defect.
func (*counterValuePanicsDefect) Value(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Counter whose Value panics")
}

// counterValuePanics returns the planted defect of a Counter whose Value
// panics over a fresh subject.
func counterValuePanics(h *CounterHandle) Counter {
	return &counterValuePanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueIgnoresCancelDefect is the planted defect of a Counter whose
// Value ignores the context.
type counterValueIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Value is the override of the planted defect.
func (d *counterValueIgnoresCancelDefect) Value(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.counterOverEmbedded.Value(p0)
}

// counterValueIgnoresCancel returns the planted defect of a Counter whose
// Value ignores the context over a fresh subject.
func counterValueIgnoresCancel(h *CounterHandle) Counter {
	return &counterValueIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueValueBesideDefect is the planted defect of a Counter whose
// Value returns a value beside an error.
type counterValueValueBesideDefect struct {
	counterOverEmbedded
}

// Value is the override of the planted defect.
func (d *counterValueValueBesideDefect) Value(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.counterOverEmbedded.Value(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// counterValueValueBeside returns the planted defect of a Counter whose
// Value returns a value beside an error over a fresh subject.
func counterValueValueBeside(h *CounterHandle) Counter {
	return &counterValueValueBesideDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueDescendsDefect is the planted defect of a Counter whose
// Value returns less than its previous result.
type counterValueDescendsDefect struct {
	counterOverEmbedded
	calls int
}

// Value is the override of the planted defect.
func (d *counterValueDescendsDefect) Value(p0 context.Context) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls == 1 {
		r0 = int(7)
	}
	return r0, r1
}

// counterValueDescends returns the planted defect of a Counter whose Value
// returns less than its previous result over a fresh subject.
func counterValueDescends(h *CounterHandle) Counter {
	return &counterValueDescendsDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueGoesBackAfterThreeDefect is the planted defect of a Counter
// whose Value returns its first result again from its fourth call on.
type counterValueGoesBackAfterThreeDefect struct {
	counterOverEmbedded
	calls int
	first *int
}

// Value is the override of the planted defect.
func (d *counterValueGoesBackAfterThreeDefect) Value(p0 context.Context) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.counterOverEmbedded.Value(p0)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// counterValueGoesBackAfterThree returns the planted defect of a Counter
// whose Value returns its first result again from its fourth call on over
// a fresh subject.
func counterValueGoesBackAfterThree(h *CounterHandle) Counter {
	return &counterValueGoesBackAfterThreeDefect{counterOverEmbedded: h.Fresh()}
}

// counterGrowPanicsDefect is the planted defect of a Counter whose Grow
// panics.
type counterGrowPanicsDefect struct {
	counterOverEmbedded
}

// Grow is the override of the planted defect.
func (*counterGrowPanicsDefect) Grow(p0 context.Context) (r0 error) {
	panic("the planted defect of a Counter whose Grow panics")
}

// counterGrowPanics returns the planted defect of a Counter whose Grow
// panics over a fresh subject.
func counterGrowPanics(h *CounterHandle) Counter {
	return &counterGrowPanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterGrowIgnoresCancelDefect is the planted defect of a Counter whose
// Grow ignores the context.
type counterGrowIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Grow is the override of the planted defect.
func (d *counterGrowIgnoresCancelDefect) Grow(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.counterOverEmbedded.Grow(p0)
}

// counterGrowIgnoresCancel returns the planted defect of a Counter whose
// Grow ignores the context over a fresh subject.
func counterGrowIgnoresCancel(h *CounterHandle) Counter {
	return &counterGrowIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizeGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Counter whose Size reports a use above the limit from its fourth call
// on.
type counterSizeGrowsPastTheLimitAfterThreeDefect struct {
	counterOverEmbedded
	calls int
}

// Size is the override of the planted defect.
func (d *counterSizeGrowsPastTheLimitAfterThreeDefect) Size(
	p0 context.Context,
) (r0 uint64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.counterOverEmbedded.Size(p0)
	d.calls++
	if d.calls > 3 {
		r0 = 64 + 1
	}
	return r0, r1
}

// counterSizeGrowsPastTheLimitAfterThree returns the planted defect of a
// Counter whose Size reports a use above the limit from its fourth call on
// over a fresh subject.
func counterSizeGrowsPastTheLimitAfterThree(h *CounterHandle) Counter {
	return &counterSizeGrowsPastTheLimitAfterThreeDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizePanicsDefect is the planted defect of a Counter whose Size
// panics.
type counterSizePanicsDefect struct {
	counterOverEmbedded
}

// Size is the override of the planted defect.
func (*counterSizePanicsDefect) Size(p0 context.Context) (r0 uint64, r1 error) {
	panic("the planted defect of a Counter whose Size panics")
}

// counterSizePanics returns the planted defect of a Counter whose Size
// panics over a fresh subject.
func counterSizePanics(h *CounterHandle) Counter {
	return &counterSizePanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizeIgnoresCancelDefect is the planted defect of a Counter whose
// Size ignores the context.
type counterSizeIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Size is the override of the planted defect.
func (d *counterSizeIgnoresCancelDefect) Size(p0 context.Context) (r0 uint64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.counterOverEmbedded.Size(p0)
}

// counterSizeIgnoresCancel returns the planted defect of a Counter whose
// Size ignores the context over a fresh subject.
func counterSizeIgnoresCancel(h *CounterHandle) Counter {
	return &counterSizeIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizeValueBesideDefect is the planted defect of a Counter whose
// Size returns a value beside an error.
type counterSizeValueBesideDefect struct {
	counterOverEmbedded
}

// Size is the override of the planted defect.
func (d *counterSizeValueBesideDefect) Size(p0 context.Context) (r0 uint64, r1 error) {
	r0, r1 = d.counterOverEmbedded.Size(p0)
	if r1 != nil {
		r0 = uint64(7)
	}
	return r0, r1
}

// counterSizeValueBeside returns the planted defect of a Counter whose
// Size returns a value beside an error over a fresh subject.
func counterSizeValueBeside(h *CounterHandle) Counter {
	return &counterSizeValueBesideDefect{counterOverEmbedded: h.Fresh()}
}

// counterThrottlePanicsDefect is the planted defect of a Counter whose
// Throttle panics.
type counterThrottlePanicsDefect struct {
	counterOverEmbedded
}

// Throttle is the override of the planted defect.
func (*counterThrottlePanicsDefect) Throttle(p0 context.Context) (r0 error) {
	panic("the planted defect of a Counter whose Throttle panics")
}

// counterThrottlePanics returns the planted defect of a Counter whose
// Throttle panics over a fresh subject.
func counterThrottlePanics(h *CounterHandle) Counter {
	return &counterThrottlePanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterThrottleIgnoresCancelDefect is the planted defect of a Counter
// whose Throttle ignores the context.
type counterThrottleIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Throttle is the override of the planted defect.
func (d *counterThrottleIgnoresCancelDefect) Throttle(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.counterOverEmbedded.Throttle(p0)
}

// counterThrottleIgnoresCancel returns the planted defect of a Counter
// whose Throttle ignores the context over a fresh subject.
func counterThrottleIgnoresCancel(h *CounterHandle) Counter {
	return &counterThrottleIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterLoadPanicsDefect is the planted defect of a Counter whose Load
// panics.
type counterLoadPanicsDefect struct {
	counterOverEmbedded
}

// Load is the override of the planted defect.
func (*counterLoadPanicsDefect) Load(p0 context.Context) (r0 float64, r1 error) {
	panic("the planted defect of a Counter whose Load panics")
}

// counterLoadPanics returns the planted defect of a Counter whose Load
// panics over a fresh subject.
func counterLoadPanics(h *CounterHandle) Counter {
	return &counterLoadPanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterLoadIgnoresCancelDefect is the planted defect of a Counter whose
// Load ignores the context.
type counterLoadIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Load is the override of the planted defect.
func (d *counterLoadIgnoresCancelDefect) Load(p0 context.Context) (r0 float64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.counterOverEmbedded.Load(p0)
}

// counterLoadIgnoresCancel returns the planted defect of a Counter whose
// Load ignores the context over a fresh subject.
func counterLoadIgnoresCancel(h *CounterHandle) Counter {
	return &counterLoadIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterLoadValueBesideDefect is the planted defect of a Counter whose
// Load returns a value beside an error.
type counterLoadValueBesideDefect struct {
	counterOverEmbedded
}

// Load is the override of the planted defect.
func (d *counterLoadValueBesideDefect) Load(p0 context.Context) (r0 float64, r1 error) {
	r0, r1 = d.counterOverEmbedded.Load(p0)
	if r1 != nil {
		r0 = float64(2.5)
	}
	return r0, r1
}

// counterLoadValueBeside returns the planted defect of a Counter whose
// Load returns a value beside an error over a fresh subject.
func counterLoadValueBeside(h *CounterHandle) Counter {
	return &counterLoadValueBesideDefect{counterOverEmbedded: h.Fresh()}
}

// counterActions returns the actions of the generated sequences of Counter
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func counterActions(c *prop.Case, subject func() Counter) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Inc",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Inc", args)
				values, err := counterIncCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Value",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Value", args)
				values, err := counterValueCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Grow",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Grow", args)
				values, err := counterGrowCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Size",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Size", args)
				values, err := counterSizeCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Throttle",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Throttle", args)
				values, err := counterThrottleCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Load",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Load", args)
				values, err := counterLoadCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// counterIncCall calls Inc on subject with the arguments args, and returns
// the values of the call beside its error.
func counterIncCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	return nil, subject.Inc(ctx)
}

// counterValueCall calls Value on subject with the arguments args, and
// returns the values of the call beside its error.
func counterValueCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	r0, err := subject.Value(ctx)
	return []any{r0}, err
}

// counterGrowCall calls Grow on subject with the arguments args, and
// returns the values of the call beside its error.
func counterGrowCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	return nil, subject.Grow(ctx)
}

// counterSizeCall calls Size on subject with the arguments args, and
// returns the values of the call beside its error.
func counterSizeCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	r0, err := subject.Size(ctx)
	return []any{r0}, err
}

// counterThrottleCall calls Throttle on subject with the arguments args,
// and returns the values of the call beside its error.
func counterThrottleCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	return nil, subject.Throttle(ctx)
}

// counterLoadCall calls Load on subject with the arguments args, and
// returns the values of the call beside its error.
func counterLoadCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	r0, err := subject.Load(ctx)
	return []any{r0}, err
}

// The derived checks of Notes:
//
// Notes.Push/smoke: Push returns without a panic for derived inputs.
//
// Notes.Push/smoke#generated: Push returns without a panic for derived
// inputs.
//
// Notes.Push/cancel: Push returns an error for a cancelled context.
//
// Notes.Push/cancel#generated: Push returns an error for a cancelled
// context.
//
// Notes.Push/keeps-within-its-bounds: Push keeps the use that Items reads
// within its bounds, over any sequence of calls.
//
// Notes.Pop/smoke: Pop returns without a panic for derived inputs.
//
// Notes.Pop/cancel: Pop returns an error for a cancelled context.
//
// Notes.Items/smoke: Items returns without a panic for derived inputs.
//
// Notes.Items/smoke#generated: Items returns without a panic for derived
// inputs.
//
// Notes.Items/cancel: Items returns an error for a cancelled context.
//
// Notes.Items/cancel#generated: Items returns an error for a cancelled
// context.
//
// Notes.Items/zero-on-error: Items returns the zero value beside an error.
//
// Notes.Items/zero-on-error#generated: Items returns the zero value beside
// an error.
//
// Notes.Append/smoke: Append returns without a panic for derived inputs.
//
// Notes.Append/smoke#generated: Append returns without a panic for derived
// inputs.
//
// Notes.Append/cancel: Append returns an error for a cancelled context.
//
// Notes.Append/cancel#generated: Append returns an error for a cancelled
// context.
//
// Notes.Append/keeps-within-its-bounds: Append keeps the use that Text
// reads within its bounds, over any sequence of calls.
//
// Notes.Text/smoke: Text returns without a panic for derived inputs.
//
// Notes.Text/cancel: Text returns an error for a cancelled context.
//
// Notes.Text/zero-on-error: Text returns the zero value beside an error.
//
// Notes.Pin/smoke: Pin returns without a panic for derived inputs.
//
// Notes.Pin/smoke#generated: Pin returns without a panic for derived
// inputs.
//
// Notes.Pin/cancel: Pin returns an error for a cancelled context.
//
// Notes.Pin/cancel#generated: Pin returns an error for a cancelled
// context.
//
// Notes.Head/smoke: Head returns without a panic for derived inputs.
//
// Notes.Head/cancel: Head returns an error for a cancelled context.
//
// Notes.Head/zero-on-error: Head returns the zero value beside an error.

// The derived checks of Notes that witness their laws at the shallow end:
//
// Notes.Push/smoke
//
// Notes.Pop/smoke
//
// Notes.Items/smoke
//
// Notes.Append/smoke
//
// Notes.Text/smoke
//
// Notes.Pin/smoke
//
// Notes.Head/smoke

// The refused claims of Notes:
//
// Notes.Pop/keeps-within-its-bounds: the reader that observe names for Pop
// returns no number and no length. The remedy is to name with observe a
// callable that returns the use of Pop as a number or a list.
//
// Notes.Pin/keeps-within-its-bounds: the reader that observe names for Pin
// returns no number and no length. The remedy is to name with observe a
// callable that returns the use of Pin as a number or a list.

// The unclaimed classifications of Notes:
//
// The shape:writer of Notes.Push has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Notes.Pop has the stance silent. A lifecycle call
// returns nothing that a check could compare.
//
// The shape:writer of Notes.Append has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Notes.Pin has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Notes writes.
type (
	NotesHandle = conform.Handle[Notes]
	NotesCheck  = conform.Check[Notes]
	NotesOption = conform.Option[Notes]
)

// NotesChecks contains the name of each callable of Notes, with the
// identities of its checks.
var NotesChecks = struct {
	Push   notesPush
	Pop    notesPop
	Items  notesItems
	Append notesAppend
	Text   notesText
	Pin    notesPin
	Head   notesHead
}{
	Push:   notesPush{"Notes.Push"},
	Pop:    notesPop{"Notes.Pop"},
	Items:  notesItems{"Notes.Items"},
	Append: notesAppend{"Notes.Append"},
	Text:   notesText{"Notes.Text"},
	Pin:    notesPin{"Notes.Pin"},
	Head:   notesHead{"Notes.Head"},
}

// notesPush is the name of Notes.Push, with the identities of its checks.
type notesPush struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Push returns without a
// panic for derived inputs.
func (c notesPush) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Push returns
// without a panic for derived inputs.
func (c notesPush) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Push returns an error for
// a cancelled context.
func (c notesPush) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Push returns an
// error for a cancelled context.
func (c notesPush) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// KeepsWithinItsBounds returns the identity of the check that Push keeps
// the use that Items reads within its bounds, over any sequence of calls.
func (c notesPush) KeepsWithinItsBounds() conform.ID[Notes] {
	return c.Check("keeps-within-its-bounds")
}

// notesPop is the name of Notes.Pop, with the identities of its checks.
type notesPop struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Pop returns without a panic
// for derived inputs.
func (c notesPop) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Pop returns an error for a
// cancelled context.
func (c notesPop) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// notesItems is the name of Notes.Items, with the identities of its
// checks.
type notesItems struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Items returns without a
// panic for derived inputs.
func (c notesItems) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Items returns
// without a panic for derived inputs.
func (c notesItems) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Items returns an error for
// a cancelled context.
func (c notesItems) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Items returns an
// error for a cancelled context.
func (c notesItems) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Items returns the
// zero value beside an error.
func (c notesItems) ZeroOnError() conform.ID[Notes] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Items
// returns the zero value beside an error.
func (c notesItems) ZeroOnErrorGenerated() conform.ID[Notes] {
	return c.Check("zero-on-error#generated")
}

// notesAppend is the name of Notes.Append, with the identities of its
// checks.
type notesAppend struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c notesAppend) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Append returns
// without a panic for derived inputs.
func (c notesAppend) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c notesAppend) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Append returns an
// error for a cancelled context.
func (c notesAppend) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// KeepsWithinItsBounds returns the identity of the check that Append keeps
// the use that Text reads within its bounds, over any sequence of calls.
func (c notesAppend) KeepsWithinItsBounds() conform.ID[Notes] {
	return c.Check("keeps-within-its-bounds")
}

// notesText is the name of Notes.Text, with the identities of its checks.
type notesText struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Text returns without a
// panic for derived inputs.
func (c notesText) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Text returns an error for
// a cancelled context.
func (c notesText) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Text returns the zero
// value beside an error.
func (c notesText) ZeroOnError() conform.ID[Notes] {
	return c.Check("zero-on-error")
}

// notesPin is the name of Notes.Pin, with the identities of its checks.
type notesPin struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Pin returns without a panic
// for derived inputs.
func (c notesPin) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Pin returns
// without a panic for derived inputs.
func (c notesPin) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Pin returns an error for a
// cancelled context.
func (c notesPin) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Pin returns an
// error for a cancelled context.
func (c notesPin) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// notesHead is the name of Notes.Head, with the identities of its checks.
type notesHead struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Head returns without a
// panic for derived inputs.
func (c notesHead) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Head returns an error for
// a cancelled context.
func (c notesHead) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Head returns the zero
// value beside an error.
func (c notesHead) ZeroOnError() conform.ID[Notes] {
	return c.Check("zero-on-error")
}

// NotesSuite returns the checks of Notes, and the calls and the budgets of
// its callables.
func NotesSuite() conform.Suite[Notes] {
	return conform.Suite[Notes]{
		Declaration: "Notes",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Notes]{
			{
				ID:          NotesChecks.Push.Smoke(),
				Statement:   "Push returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, string("item-a"))
					})
				},
				Defect: notesPushPanics,
			},
			{
				ID:          NotesChecks.Push.SmokeGenerated(),
				Statement:   "Push returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, drawnItem)
					})
				},
				Defect: notesPushPanics,
			},
			{
				ID:          NotesChecks.Push.Cancel(),
				Statement:   "Push returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, string("item-a"))
					})
				},
				Defect: notesPushIgnoresCancel,
			},
			{
				ID:          NotesChecks.Push.CancelGenerated(),
				Statement:   "Push returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, drawnItem)
					})
				},
				Defect: notesPushIgnoresCancel,
			},
			{
				ID:          NotesChecks.Push.KeepsWithinItsBounds(),
				Statement:   "Push keeps the use that Items reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *NotesHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						notesActions,
						func(ctx context.Context) (int, error) {
							values, err := notesItemsCall(
								ctx,
								h.Subject(),
								[]any{string("prefix-a")},
							)
							v, _ := values[0].([]string)
							return len(v), err
						},
						0,
						4,
					)
				},
				Defect: notesItemsGrowsPastTheLimitAfterThree,
			},
			{
				ID:          NotesChecks.Pop.Smoke(),
				Statement:   "Pop returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pop(ctx)
					})
				},
				Defect: notesPopPanics,
			},
			{
				ID:          NotesChecks.Pop.Cancel(),
				Statement:   "Pop returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pop(ctx)
					})
				},
				Defect: notesPopIgnoresCancel,
			},
			{
				ID:          NotesChecks.Items.Smoke(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, string("prefix-a"))
						return err
					})
				},
				Defect: notesItemsPanics,
			},
			{
				ID:          NotesChecks.Items.SmokeGenerated(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnPrefix := c.Draw(prop.Of[string](), "prefix")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, drawnPrefix)
						return err
					})
				},
				Defect: notesItemsPanics,
			},
			{
				ID:          NotesChecks.Items.Cancel(),
				Statement:   "Items returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, string("prefix-a"))
						return err
					})
				},
				Defect: notesItemsIgnoresCancel,
			},
			{
				ID:          NotesChecks.Items.CancelGenerated(),
				Statement:   "Items returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnPrefix := c.Draw(prop.Of[string](), "prefix")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, drawnPrefix)
						return err
					})
				},
				Defect: notesItemsIgnoresCancel,
			},
			{
				ID:          NotesChecks.Items.ZeroOnError(),
				Statement:   "Items returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Items(ctx, string("prefix-a"))
					}, []string(nil))
				},
				Defect: notesItemsValueBeside,
			},
			{
				ID:          NotesChecks.Items.ZeroOnErrorGenerated(),
				Statement:   "Items returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnPrefix := c.Draw(prop.Of[string](), "prefix")
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Items(ctx, drawnPrefix)
					}, []string(nil))
				},
				Defect: notesItemsValueBeside,
			},
			{
				ID:          NotesChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, string("text-a"))
					})
				},
				Defect: notesAppendPanics,
			},
			{
				ID:          NotesChecks.Append.SmokeGenerated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnText := c.Draw(prop.Of[string](), "text")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnText)
					})
				},
				Defect: notesAppendPanics,
			},
			{
				ID:          NotesChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, string("text-a"))
					})
				},
				Defect: notesAppendIgnoresCancel,
			},
			{
				ID:          NotesChecks.Append.CancelGenerated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnText := c.Draw(prop.Of[string](), "text")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnText)
					})
				},
				Defect: notesAppendIgnoresCancel,
			},
			{
				ID:          NotesChecks.Append.KeepsWithinItsBounds(),
				Statement:   "Append keeps the use that Text reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *NotesHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						notesActions,
						func(ctx context.Context) (int, error) {
							values, err := notesTextCall(ctx, h.Subject(), []any{})
							v, _ := values[0].(string)
							return len(v), err
						},
						0,
						8,
					)
				},
				Defect: notesTextGrowsPastTheLimitAfterThree,
			},
			{
				ID:          NotesChecks.Text.Smoke(),
				Statement:   "Text returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Text(ctx)
						return err
					})
				},
				Defect: notesTextPanics,
			},
			{
				ID:          NotesChecks.Text.Cancel(),
				Statement:   "Text returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Text(ctx)
						return err
					})
				},
				Defect: notesTextIgnoresCancel,
			},
			{
				ID:          NotesChecks.Text.ZeroOnError(),
				Statement:   "Text returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Text(ctx)
					}, string(""))
				},
				Defect: notesTextValueBeside,
			},
			{
				ID:          NotesChecks.Pin.Smoke(),
				Statement:   "Pin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: notesPinPanics,
			},
			{
				ID:          NotesChecks.Pin.SmokeGenerated(),
				Statement:   "Pin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, drawnIt)
					})
				},
				Defect: notesPinPanics,
			},
			{
				ID:          NotesChecks.Pin.Cancel(),
				Statement:   "Pin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: notesPinIgnoresCancel,
			},
			{
				ID:          NotesChecks.Pin.CancelGenerated(),
				Statement:   "Pin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, drawnIt)
					})
				},
				Defect: notesPinIgnoresCancel,
			},
			{
				ID:          NotesChecks.Head.Smoke(),
				Statement:   "Head returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Head(ctx)
						return err
					})
				},
				Defect: notesHeadPanics,
			},
			{
				ID:          NotesChecks.Head.Cancel(),
				Statement:   "Head returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Head(ctx)
						return err
					})
				},
				Defect: notesHeadIgnoresCancel,
			},
			{
				ID:          NotesChecks.Head.ZeroOnError(),
				Statement:   "Head returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Head(ctx)
					}, Item{})
				},
				Defect: notesHeadValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.KeepsWithinItsBounds,
				Callable: "Notes.Pop",
				Gate:     spec.CountableUse,
				Why:      "the reader that observe names for Pop returns no number and no length",
				Remedy:   "name with observe a callable that returns the use of Pop as a number or a list",
			},
			{
				Claim:    spec.KeepsWithinItsBounds,
				Callable: "Notes.Pin",
				Gate:     spec.CountableUse,
				Why:      "the reader that observe names for Pin returns no number and no length",
				Remedy:   "name with observe a callable that returns the use of Pin as a number or a list",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Notes.Push",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Notes.Pop",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Notes.Append",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Notes.Pin",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Notes]{
			{
				Calls:    notesPushCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Push.Callable,
			},
			{
				Calls:    notesPopCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Pop.Callable,
			},
			{
				Calls:    notesItemsCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Items.Callable,
			},
			{
				Calls:    notesAppendCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Append.Callable,
			},
			{
				Calls:    notesTextCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Text.Callable,
			},
			{
				Calls:    notesPinCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Pin.Callable,
			},
			{
				Calls:    notesHeadCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Head.Callable,
			},
		},
	}
}

// RunNotesSuite applies the checks of Notes, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunNotesSuite[T Notes](t *testing.T, newSubject func() T, opts ...NotesOption) {
	t.Helper()
	RunNotesSuiteFromHandle(t, func(*NotesHandle) T { return newSubject() }, opts...)
}

// RunNotesSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunNotesSuiteFromHandle[T Notes](
	t *testing.T,
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	t.Helper()
	conform.Run(t, NotesSuite(), notesSubjects(newSubject), opts...)
}

// FuzzNotesSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzNotesSuite[T Notes](
	f *testing.F,
	callable conform.Callable[Notes],
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	f.Helper()
	conform.Fuzz(f, NotesSuite(), callable, notesSubjects(newSubject), opts...)
}

// RunNotesBudgets applies the allocation ceilings of Notes and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunNotesBudgets[T Notes](
	t *testing.T,
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	t.Helper()
	conform.RunBudgets(t, NotesSuite(), notesSubjects(newSubject), opts...)
}

// BenchNotesSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchNotesSuite[T Notes](
	b conform.B,
	callable conform.Callable[Notes],
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	b.Helper()
	conform.Bench(b, NotesSuite(), callable, notesSubjects(newSubject), opts...)
}

// notesSubjects returns the factory of the subjects of Notes that calls
// newSubject with the handle of each check.
func notesSubjects[T Notes](newSubject func(h *NotesHandle) T) func(h *NotesHandle) Notes {
	return func(h *NotesHandle) Notes { return newSubject(h) }
}

// notesPushCalls generates the calls of Notes.Push with the inputs that
// each case draws.
var notesPushCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Notes) error {
		return subject.Push(ctx, item)
	}
})

// notesPopCalls generates the calls of Notes.Pop with the inputs that each
// case draws.
var notesPopCalls = prop.Composite(func(*prop.Case) func(context.Context, Notes) error {
	return func(ctx context.Context, subject Notes) error {
		return subject.Pop(ctx)
	}
})

// notesItemsCalls generates the calls of Notes.Items with the inputs that
// each case draws.
var notesItemsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	prefix := c.Draw(prop.Of[string](), "prefix")
	return func(ctx context.Context, subject Notes) error {
		_, err := subject.Items(ctx, prefix)
		return err
	}
})

// notesAppendCalls generates the calls of Notes.Append with the inputs
// that each case draws.
var notesAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	text := c.Draw(prop.Of[string](), "text")
	return func(ctx context.Context, subject Notes) error {
		return subject.Append(ctx, text)
	}
})

// notesTextCalls generates the calls of Notes.Text with the inputs that
// each case draws.
var notesTextCalls = prop.Composite(func(*prop.Case) func(context.Context, Notes) error {
	return func(ctx context.Context, subject Notes) error {
		_, err := subject.Text(ctx)
		return err
	}
})

// notesPinCalls generates the calls of Notes.Pin with the inputs that each
// case draws.
var notesPinCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Notes) error {
		return subject.Pin(ctx, it)
	}
})

// notesHeadCalls generates the calls of Notes.Head with the inputs that
// each case draws.
var notesHeadCalls = prop.Composite(func(*prop.Case) func(context.Context, Notes) error {
	return func(ctx context.Context, subject Notes) error {
		_, err := subject.Head(ctx)
		return err
	}
})

// notesOverEmbedded is Notes under the name of the field that notesOver
// embeds.
type notesOverEmbedded = Notes

// notesOver is a subject of Notes whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type notesOver struct {
	notesOverEmbedded
	spend func()
}

// Push calls spend and returns its zero results.
func (o notesOver) Push(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Pop calls spend and returns its zero results.
func (o notesOver) Pop(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Items calls spend and returns its zero results.
func (o notesOver) Items(context.Context, string) (r0 []string, r1 error) {
	o.spend()
	return r0, r1
}

// Append calls spend and returns its zero results.
func (o notesOver) Append(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Text calls spend and returns its zero results.
func (o notesOver) Text(context.Context) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// Pin calls spend and returns its zero results.
func (o notesOver) Pin(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Head calls spend and returns its zero results.
func (o notesOver) Head(context.Context) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newNotesOver returns the planted defect of a ceiling of Notes, whose
// callables call spend.
func newNotesOver(_ *NotesHandle, spend func()) Notes {
	return notesOver{spend: spend}
}

// notesPushPanicsDefect is the planted defect of a Notes whose Push
// panics.
type notesPushPanicsDefect struct {
	notesOverEmbedded
}

// Push is the override of the planted defect.
func (*notesPushPanicsDefect) Push(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Notes whose Push panics")
}

// notesPushPanics returns the planted defect of a Notes whose Push panics
// over a fresh subject.
func notesPushPanics(h *NotesHandle) Notes {
	return &notesPushPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesPushIgnoresCancelDefect is the planted defect of a Notes whose Push
// ignores the context.
type notesPushIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Push is the override of the planted defect.
func (d *notesPushIgnoresCancelDefect) Push(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Push(p0, p1)
}

// notesPushIgnoresCancel returns the planted defect of a Notes whose Push
// ignores the context over a fresh subject.
func notesPushIgnoresCancel(h *NotesHandle) Notes {
	return &notesPushIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Notes whose Items reports a use above the limit from its fourth call on.
type notesItemsGrowsPastTheLimitAfterThreeDefect struct {
	notesOverEmbedded
	calls int
}

// Items is the override of the planted defect.
func (d *notesItemsGrowsPastTheLimitAfterThreeDefect) Items(
	p0 context.Context,
	p1 string,
) (r0 []string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.notesOverEmbedded.Items(p0, p1)
	d.calls++
	if d.calls > 3 {
		r0 = make([]string, 4+1)
	}
	return r0, r1
}

// notesItemsGrowsPastTheLimitAfterThree returns the planted defect of a
// Notes whose Items reports a use above the limit from its fourth call on
// over a fresh subject.
func notesItemsGrowsPastTheLimitAfterThree(h *NotesHandle) Notes {
	return &notesItemsGrowsPastTheLimitAfterThreeDefect{notesOverEmbedded: h.Fresh()}
}

// notesPopPanicsDefect is the planted defect of a Notes whose Pop panics.
type notesPopPanicsDefect struct {
	notesOverEmbedded
}

// Pop is the override of the planted defect.
func (*notesPopPanicsDefect) Pop(p0 context.Context) (r0 error) {
	panic("the planted defect of a Notes whose Pop panics")
}

// notesPopPanics returns the planted defect of a Notes whose Pop panics
// over a fresh subject.
func notesPopPanics(h *NotesHandle) Notes {
	return &notesPopPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesPopIgnoresCancelDefect is the planted defect of a Notes whose Pop
// ignores the context.
type notesPopIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Pop is the override of the planted defect.
func (d *notesPopIgnoresCancelDefect) Pop(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Pop(p0)
}

// notesPopIgnoresCancel returns the planted defect of a Notes whose Pop
// ignores the context over a fresh subject.
func notesPopIgnoresCancel(h *NotesHandle) Notes {
	return &notesPopIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsPanicsDefect is the planted defect of a Notes whose Items
// panics.
type notesItemsPanicsDefect struct {
	notesOverEmbedded
}

// Items is the override of the planted defect.
func (*notesItemsPanicsDefect) Items(p0 context.Context, p1 string) (r0 []string, r1 error) {
	panic("the planted defect of a Notes whose Items panics")
}

// notesItemsPanics returns the planted defect of a Notes whose Items
// panics over a fresh subject.
func notesItemsPanics(h *NotesHandle) Notes {
	return &notesItemsPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsIgnoresCancelDefect is the planted defect of a Notes whose
// Items ignores the context.
type notesItemsIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Items is the override of the planted defect.
func (d *notesItemsIgnoresCancelDefect) Items(
	p0 context.Context,
	p1 string,
) (r0 []string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.notesOverEmbedded.Items(p0, p1)
}

// notesItemsIgnoresCancel returns the planted defect of a Notes whose
// Items ignores the context over a fresh subject.
func notesItemsIgnoresCancel(h *NotesHandle) Notes {
	return &notesItemsIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsValueBesideDefect is the planted defect of a Notes whose Items
// returns a value beside an error.
type notesItemsValueBesideDefect struct {
	notesOverEmbedded
}

// Items is the override of the planted defect.
func (d *notesItemsValueBesideDefect) Items(p0 context.Context, p1 string) (r0 []string, r1 error) {
	r0, r1 = d.notesOverEmbedded.Items(p0, p1)
	if r1 != nil {
		r0 = []string{"sample-b"}
	}
	return r0, r1
}

// notesItemsValueBeside returns the planted defect of a Notes whose Items
// returns a value beside an error over a fresh subject.
func notesItemsValueBeside(h *NotesHandle) Notes {
	return &notesItemsValueBesideDefect{notesOverEmbedded: h.Fresh()}
}

// notesAppendPanicsDefect is the planted defect of a Notes whose Append
// panics.
type notesAppendPanicsDefect struct {
	notesOverEmbedded
}

// Append is the override of the planted defect.
func (*notesAppendPanicsDefect) Append(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Notes whose Append panics")
}

// notesAppendPanics returns the planted defect of a Notes whose Append
// panics over a fresh subject.
func notesAppendPanics(h *NotesHandle) Notes {
	return &notesAppendPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesAppendIgnoresCancelDefect is the planted defect of a Notes whose
// Append ignores the context.
type notesAppendIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Append is the override of the planted defect.
func (d *notesAppendIgnoresCancelDefect) Append(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Append(p0, p1)
}

// notesAppendIgnoresCancel returns the planted defect of a Notes whose
// Append ignores the context over a fresh subject.
func notesAppendIgnoresCancel(h *NotesHandle) Notes {
	return &notesAppendIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Notes whose Text reports a use above the limit from its fourth call on.
type notesTextGrowsPastTheLimitAfterThreeDefect struct {
	notesOverEmbedded
	calls int
}

// Text is the override of the planted defect.
func (d *notesTextGrowsPastTheLimitAfterThreeDefect) Text(
	p0 context.Context,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.notesOverEmbedded.Text(p0)
	d.calls++
	if d.calls > 3 {
		r0 = string(strings.Repeat("x", 8+1))
	}
	return r0, r1
}

// notesTextGrowsPastTheLimitAfterThree returns the planted defect of a
// Notes whose Text reports a use above the limit from its fourth call on
// over a fresh subject.
func notesTextGrowsPastTheLimitAfterThree(h *NotesHandle) Notes {
	return &notesTextGrowsPastTheLimitAfterThreeDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextPanicsDefect is the planted defect of a Notes whose Text
// panics.
type notesTextPanicsDefect struct {
	notesOverEmbedded
}

// Text is the override of the planted defect.
func (*notesTextPanicsDefect) Text(p0 context.Context) (r0 string, r1 error) {
	panic("the planted defect of a Notes whose Text panics")
}

// notesTextPanics returns the planted defect of a Notes whose Text panics
// over a fresh subject.
func notesTextPanics(h *NotesHandle) Notes {
	return &notesTextPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextIgnoresCancelDefect is the planted defect of a Notes whose Text
// ignores the context.
type notesTextIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Text is the override of the planted defect.
func (d *notesTextIgnoresCancelDefect) Text(p0 context.Context) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.notesOverEmbedded.Text(p0)
}

// notesTextIgnoresCancel returns the planted defect of a Notes whose Text
// ignores the context over a fresh subject.
func notesTextIgnoresCancel(h *NotesHandle) Notes {
	return &notesTextIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextValueBesideDefect is the planted defect of a Notes whose Text
// returns a value beside an error.
type notesTextValueBesideDefect struct {
	notesOverEmbedded
}

// Text is the override of the planted defect.
func (d *notesTextValueBesideDefect) Text(p0 context.Context) (r0 string, r1 error) {
	r0, r1 = d.notesOverEmbedded.Text(p0)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// notesTextValueBeside returns the planted defect of a Notes whose Text
// returns a value beside an error over a fresh subject.
func notesTextValueBeside(h *NotesHandle) Notes {
	return &notesTextValueBesideDefect{notesOverEmbedded: h.Fresh()}
}

// notesPinPanicsDefect is the planted defect of a Notes whose Pin panics.
type notesPinPanicsDefect struct {
	notesOverEmbedded
}

// Pin is the override of the planted defect.
func (*notesPinPanicsDefect) Pin(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Notes whose Pin panics")
}

// notesPinPanics returns the planted defect of a Notes whose Pin panics
// over a fresh subject.
func notesPinPanics(h *NotesHandle) Notes {
	return &notesPinPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesPinIgnoresCancelDefect is the planted defect of a Notes whose Pin
// ignores the context.
type notesPinIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Pin is the override of the planted defect.
func (d *notesPinIgnoresCancelDefect) Pin(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Pin(p0, p1)
}

// notesPinIgnoresCancel returns the planted defect of a Notes whose Pin
// ignores the context over a fresh subject.
func notesPinIgnoresCancel(h *NotesHandle) Notes {
	return &notesPinIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesHeadPanicsDefect is the planted defect of a Notes whose Head
// panics.
type notesHeadPanicsDefect struct {
	notesOverEmbedded
}

// Head is the override of the planted defect.
func (*notesHeadPanicsDefect) Head(p0 context.Context) (r0 Item, r1 error) {
	panic("the planted defect of a Notes whose Head panics")
}

// notesHeadPanics returns the planted defect of a Notes whose Head panics
// over a fresh subject.
func notesHeadPanics(h *NotesHandle) Notes {
	return &notesHeadPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesHeadIgnoresCancelDefect is the planted defect of a Notes whose Head
// ignores the context.
type notesHeadIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Head is the override of the planted defect.
func (d *notesHeadIgnoresCancelDefect) Head(p0 context.Context) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.notesOverEmbedded.Head(p0)
}

// notesHeadIgnoresCancel returns the planted defect of a Notes whose Head
// ignores the context over a fresh subject.
func notesHeadIgnoresCancel(h *NotesHandle) Notes {
	return &notesHeadIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesHeadValueBesideDefect is the planted defect of a Notes whose Head
// returns a value beside an error.
type notesHeadValueBesideDefect struct {
	notesOverEmbedded
}

// Head is the override of the planted defect.
func (d *notesHeadValueBesideDefect) Head(p0 context.Context) (r0 Item, r1 error) {
	r0, r1 = d.notesOverEmbedded.Head(p0)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// notesHeadValueBeside returns the planted defect of a Notes whose Head
// returns a value beside an error over a fresh subject.
func notesHeadValueBeside(h *NotesHandle) Notes {
	return &notesHeadValueBesideDefect{notesOverEmbedded: h.Fresh()}
}

// notesActions returns the actions of the generated sequences of Notes
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func notesActions(c *prop.Case, subject func() Notes) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Push",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.Of[string](), "push.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Push", args)
				values, err := notesPushCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Pop",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Pop", args)
				values, err := notesPopCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Items",
			Input: func(c *prop.Case, _ struct{}) any {
				prefix := c.Draw(prop.SampledFrom(keysString...), "items.prefix")
				return []any{prefix}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Items", args, args[0])
				values, err := notesItemsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Append",
			Input: func(c *prop.Case, _ struct{}) any {
				text := c.Draw(prop.Of[string](), "append.text")
				return []any{text}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Append", args)
				values, err := notesAppendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Text",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Text", args)
				values, err := notesTextCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Pin",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "pin.it")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Pin", args)
				values, err := notesPinCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Head",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Head", args)
				values, err := notesHeadCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// notesPushCall calls Push on subject with the arguments args, and returns
// the values of the call beside its error.
func notesPushCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Push(ctx, a0)
}

// notesPopCall calls Pop on subject with the arguments args, and returns
// the values of the call beside its error.
func notesPopCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	return nil, subject.Pop(ctx)
}

// notesItemsCall calls Items on subject with the arguments args, and
// returns the values of the call beside its error.
func notesItemsCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Items(ctx, a0)
	return []any{r0}, err
}

// notesAppendCall calls Append on subject with the arguments args, and
// returns the values of the call beside its error.
func notesAppendCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Append(ctx, a0)
}

// notesTextCall calls Text on subject with the arguments args, and returns
// the values of the call beside its error.
func notesTextCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	r0, err := subject.Text(ctx)
	return []any{r0}, err
}

// notesPinCall calls Pin on subject with the arguments args, and returns
// the values of the call beside its error.
func notesPinCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Pin(ctx, a0)
}

// notesHeadCall calls Head on subject with the arguments args, and returns
// the values of the call beside its error.
func notesHeadCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	r0, err := subject.Head(ctx)
	return []any{r0}, err
}

// The derived checks of Set:
//
// Set.Merge/merges-in-any-order: two replicas that Add made diverge read
// the same state through Items after Merge merges them in either order.
//
// Set.Join/merges-in-any-order: two replicas that Add made diverge read
// the same state through Items after Join merges them in either order.
//
// Set.Add/smoke: Add returns without a panic for derived inputs.
//
// Set.Add/smoke#generated: Add returns without a panic for derived inputs.
//
// Set.Add/cancel: Add returns an error for a cancelled context.
//
// Set.Add/cancel#generated: Add returns an error for a cancelled context.
//
// Set.Items/smoke: Items returns without a panic for derived inputs.
//
// Set.Items/cancel: Items returns an error for a cancelled context.
//
// Set.Items/zero-on-error: Items returns the zero value beside an error.

// The derived checks of Set that witness their laws at the shallow end:
//
// Set.Add/smoke
//
// Set.Items/smoke

// The refused claims of Set:
//
// Set.Merge/cancel: eidos derives no sample of the parameter peer of
// Merge. The remedy is to author a sample of the parameter peer with the
// directive sample.
//
// Set.Join/cancel: eidos derives no sample of the parameter other of Join.
// The remedy is to author a sample of the parameter other with the
// directive sample.
//
// Set.Weigh/cancel: eidos derives no sample of the parameter peer of
// Weigh. The remedy is to author a sample of the parameter peer with the
// directive sample.
//
// Set.Weigh/merges-in-any-order: Weigh takes no input that another replica
// of Set can be. The remedy is to give Weigh an input of an interface that
// Set implements.

// The unclaimed classifications of Set:
//
// The shape:writer of Set.Merge has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:writer of Set.Join has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Set.Weigh has the stance silent. A write
// is an input of the other checks, and its result is an error alone.
//
// The shape:writer of Set.Add has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Set writes.
type (
	SetHandle = conform.Handle[Set]
	SetCheck  = conform.Check[Set]
	SetOption = conform.Option[Set]
)

// SetChecks contains the name of each callable of Set, with the identities
// of its checks.
var SetChecks = struct {
	Merge setMerge
	Join  setJoin
	Weigh setWeigh
	Add   setAdd
	Items setItems
}{
	Merge: setMerge{"Set.Merge"},
	Join:  setJoin{"Set.Join"},
	Weigh: setWeigh{"Set.Weigh"},
	Add:   setAdd{"Set.Add"},
	Items: setItems{"Set.Items"},
}

// setMerge is the name of Set.Merge, with the identities of its checks.
type setMerge struct {
	conform.Callable[Set]
}

// MergesInAnyOrder returns the identity of the check that two replicas
// that Add made diverge read the same state through Items after Merge
// merges them in either order.
func (c setMerge) MergesInAnyOrder() conform.ID[Set] {
	return c.Check("merges-in-any-order")
}

// setJoin is the name of Set.Join, with the identities of its checks.
type setJoin struct {
	conform.Callable[Set]
}

// MergesInAnyOrder returns the identity of the check that two replicas
// that Add made diverge read the same state through Items after Join
// merges them in either order.
func (c setJoin) MergesInAnyOrder() conform.ID[Set] {
	return c.Check("merges-in-any-order")
}

// setWeigh is the name of Set.Weigh, with the identities of its checks.
type setWeigh struct {
	conform.Callable[Set]
}

// setAdd is the name of Set.Add, with the identities of its checks.
type setAdd struct {
	conform.Callable[Set]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c setAdd) Smoke() conform.ID[Set] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Add returns
// without a panic for derived inputs.
func (c setAdd) SmokeGenerated() conform.ID[Set] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c setAdd) Cancel() conform.ID[Set] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Add returns an
// error for a cancelled context.
func (c setAdd) CancelGenerated() conform.ID[Set] {
	return c.Check("cancel#generated")
}

// setItems is the name of Set.Items, with the identities of its checks.
type setItems struct {
	conform.Callable[Set]
}

// Smoke returns the identity of the check that Items returns without a
// panic for derived inputs.
func (c setItems) Smoke() conform.ID[Set] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Items returns an error for
// a cancelled context.
func (c setItems) Cancel() conform.ID[Set] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Items returns the
// zero value beside an error.
func (c setItems) ZeroOnError() conform.ID[Set] {
	return c.Check("zero-on-error")
}

// SetSuite returns the checks of Set, and the calls and the budgets of its
// callables.
func SetSuite() conform.Suite[Set] {
	return conform.Suite[Set]{
		Declaration: "Set",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Set]{
			{
				ID:          SetChecks.Merge.MergesInAnyOrder(),
				Statement:   "two replicas that Add made diverge read the same state through Items after Merge merges them in either order",
				Claim:       spec.MergesInAnyOrder,
				Probe:       spec.MergesCommute,
				Reason:      spec.OrderDependentMerge,
				Law:         spec.MergeOrderDoesNotMatter,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *SetHandle, c *prop.Case) {
					probe.MergesCommute(h, c, setActions, setReference, "Merge", "Add", "Items")
				},
				Defect: setMergeIgnoresTheOtherReplica,
			},
			{
				ID:          SetChecks.Join.MergesInAnyOrder(),
				Statement:   "two replicas that Add made diverge read the same state through Items after Join merges them in either order",
				Claim:       spec.MergesInAnyOrder,
				Probe:       spec.MergesCommute,
				Reason:      spec.OrderDependentMerge,
				Law:         spec.MergeOrderDoesNotMatter,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *SetHandle, c *prop.Case) {
					probe.MergesCommute(h, c, setActions, setReference, "Join", "Add", "Items")
				},
				Defect: setJoinIgnoresTheOtherReplica,
			},
			{
				ID:          SetChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SetHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("item-a"))
					})
				},
				Defect: setAddPanics,
			},
			{
				ID:          SetChecks.Add.SmokeGenerated(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *SetHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnItem)
					})
				},
				Defect: setAddPanics,
			},
			{
				ID:          SetChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SetHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("item-a"))
					})
				},
				Defect: setAddIgnoresCancel,
			},
			{
				ID:          SetChecks.Add.CancelGenerated(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *SetHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnItem)
					})
				},
				Defect: setAddIgnoresCancel,
			},
			{
				ID:          SetChecks.Items.Smoke(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SetHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx)
						return err
					})
				},
				Defect: setItemsPanics,
			},
			{
				ID:          SetChecks.Items.Cancel(),
				Statement:   "Items returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SetHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx)
						return err
					})
				},
				Defect: setItemsIgnoresCancel,
			},
			{
				ID:          SetChecks.Items.ZeroOnError(),
				Statement:   "Items returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SetHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Items(ctx)
					}, []string(nil))
				},
				Defect: setItemsValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Set.Merge",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter peer of Merge",
				Remedy:   "author a sample of the parameter peer with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Set.Join",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter other of Join",
				Remedy:   "author a sample of the parameter other with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Set.Weigh",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter peer of Weigh",
				Remedy:   "author a sample of the parameter peer with the directive sample",
			},
			{
				Claim:    spec.MergesInAnyOrder,
				Callable: "Set.Weigh",
				Gate:     spec.TakesAReplica,
				Why:      "Weigh takes no input that another replica of Set can be",
				Remedy:   "give Weigh an input of an interface that Set implements",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Set.Merge",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Set.Join",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Set.Weigh",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Set.Add",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Set]{
			{
				Calls:    setMergeCalls,
				Over:     newSetOver,
				Callable: SetChecks.Merge.Callable,
			},
			{
				Calls:    setJoinCalls,
				Over:     newSetOver,
				Callable: SetChecks.Join.Callable,
			},
			{
				Calls:    setWeighCalls,
				Over:     newSetOver,
				Callable: SetChecks.Weigh.Callable,
			},
			{
				Calls:    setAddCalls,
				Over:     newSetOver,
				Callable: SetChecks.Add.Callable,
			},
			{
				Calls:    setItemsCalls,
				Over:     newSetOver,
				Callable: SetChecks.Items.Callable,
			},
		},
	}
}

// RunSetSuite applies the checks of Set, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunSetSuite[T Set](t *testing.T, newSubject func() T, opts ...SetOption) {
	t.Helper()
	RunSetSuiteFromHandle(t, func(*SetHandle) T { return newSubject() }, opts...)
}

// RunSetSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunSetSuiteFromHandle[T Set](
	t *testing.T,
	newSubject func(h *SetHandle) T,
	opts ...SetOption,
) {
	t.Helper()
	conform.Run(t, SetSuite(), setSubjects(newSubject), opts...)
}

// FuzzSetSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzSetSuite[T Set](
	f *testing.F,
	callable conform.Callable[Set],
	newSubject func(h *SetHandle) T,
	opts ...SetOption,
) {
	f.Helper()
	conform.Fuzz(f, SetSuite(), callable, setSubjects(newSubject), opts...)
}

// RunSetBudgets applies the allocation ceilings of Set and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunSetBudgets[T Set](t *testing.T, newSubject func(h *SetHandle) T, opts ...SetOption) {
	t.Helper()
	conform.RunBudgets(t, SetSuite(), setSubjects(newSubject), opts...)
}

// BenchSetSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchSetSuite[T Set](
	b conform.B,
	callable conform.Callable[Set],
	newSubject func(h *SetHandle) T,
	opts ...SetOption,
) {
	b.Helper()
	conform.Bench(b, SetSuite(), callable, setSubjects(newSubject), opts...)
}

// setSubjects returns the factory of the subjects of Set that calls
// newSubject with the handle of each check.
func setSubjects[T Set](newSubject func(h *SetHandle) T) func(h *SetHandle) Set {
	return func(h *SetHandle) Set { return newSubject(h) }
}

// setMergeCalls generates the calls of Set.Merge with the inputs that each
// case draws.
var setMergeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	peerItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "peer-items-0")
	peerItemsResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(peerItemsResults0)),
			prop.MaxSize(len(peerItemsResults0)),
		),
		"peer-items-1",
	)
	var peerItemsAt atomic.Int64
	peer := &setReplicaInput{
		onItems: func(context.Context) (r0 []string, r1 error) {
			i := int(peerItemsAt.Add(1) - 1)
			if i >= len(peerItemsResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return peerItemsResults0[i],
				peerItemsResults1[i]
		},
	}
	return func(ctx context.Context, subject Set) error {
		return subject.Merge(ctx, peer)
	}
})

// setJoinCalls generates the calls of Set.Join with the inputs that each
// case draws.
var setJoinCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	otherMergeResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-merge",
	)
	var otherMergeAt atomic.Int64
	otherJoinResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-join",
	)
	var otherJoinAt atomic.Int64
	otherWeighResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-weigh",
	)
	var otherWeighAt atomic.Int64
	otherAddResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-add",
	)
	var otherAddAt atomic.Int64
	otherItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "other-items-0")
	otherItemsResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(otherItemsResults0)),
			prop.MaxSize(len(otherItemsResults0)),
		),
		"other-items-1",
	)
	var otherItemsAt atomic.Int64
	other := &setSetInput{
		onMerge: func(context.Context, Replica) (r0 error) {
			i := int(otherMergeAt.Add(1) - 1)
			if i >= len(otherMergeResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherMergeResults[i]
		},
		onJoin: func(context.Context, Set) (r0 error) {
			i := int(otherJoinAt.Add(1) - 1)
			if i >= len(otherJoinResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherJoinResults[i]
		},
		onWeigh: func(context.Context, Replica, int) (r0 error) {
			i := int(otherWeighAt.Add(1) - 1)
			if i >= len(otherWeighResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherWeighResults[i]
		},
		onAdd: func(context.Context, string) (r0 error) {
			i := int(otherAddAt.Add(1) - 1)
			if i >= len(otherAddResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherAddResults[i]
		},
		onItems: func(context.Context) (r0 []string, r1 error) {
			i := int(otherItemsAt.Add(1) - 1)
			if i >= len(otherItemsResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return otherItemsResults0[i],
				otherItemsResults1[i]
		},
	}
	return func(ctx context.Context, subject Set) error {
		return subject.Join(ctx, other)
	}
})

// setWeighCalls generates the calls of Set.Weigh with the inputs that each
// case draws.
var setWeighCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	peerItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "peer-items-0")
	peerItemsResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(peerItemsResults0)),
			prop.MaxSize(len(peerItemsResults0)),
		),
		"peer-items-1",
	)
	var peerItemsAt atomic.Int64
	peer := &setReplicaInput{
		onItems: func(context.Context) (r0 []string, r1 error) {
			i := int(peerItemsAt.Add(1) - 1)
			if i >= len(peerItemsResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return peerItemsResults0[i],
				peerItemsResults1[i]
		},
	}
	weight := c.Draw(prop.Of[int](), "weight")
	return func(ctx context.Context, subject Set) error {
		return subject.Weigh(ctx, peer, weight)
	}
})

// setAddCalls generates the calls of Set.Add with the inputs that each
// case draws.
var setAddCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Set) error {
		return subject.Add(ctx, item)
	}
})

// setItemsCalls generates the calls of Set.Items with the inputs that each
// case draws.
var setItemsCalls = prop.Composite(func(*prop.Case) func(context.Context, Set) error {
	return func(ctx context.Context, subject Set) error {
		_, err := subject.Items(ctx)
		return err
	}
})

// setReplicaInputEmbedded is Replica under the name of the field that
// setReplicaInput embeds.
type setReplicaInputEmbedded = Replica

// setReplicaInput implements Replica. Its methods return the results that
// the generators of the suite of Set draw for them, in turn.
type setReplicaInput struct {
	setReplicaInputEmbedded
	onItems func(a0 context.Context) ([]string, error)
}

// Items calls the function of Items.
func (in setReplicaInput) Items(a0 context.Context) ([]string, error) {
	return in.onItems(a0)
}

// setSetInputEmbedded is Set under the name of the field that setSetInput
// embeds.
type setSetInputEmbedded = Set

// setSetInput implements Set. Its methods return the results that the
// generators of the suite of Set draw for them, in turn.
type setSetInput struct {
	setSetInputEmbedded
	onMerge func(a0 context.Context, a1 Replica) error
	onJoin  func(a0 context.Context, a1 Set) error
	onWeigh func(a0 context.Context, a1 Replica, a2 int) error
	onAdd   func(a0 context.Context, a1 string) error
	onItems func(a0 context.Context) ([]string, error)
}

// Merge calls the function of Merge.
func (in setSetInput) Merge(a0 context.Context, a1 Replica) error {
	return in.onMerge(a0, a1)
}

// Join calls the function of Join.
func (in setSetInput) Join(a0 context.Context, a1 Set) error {
	return in.onJoin(a0, a1)
}

// Weigh calls the function of Weigh.
func (in setSetInput) Weigh(a0 context.Context, a1 Replica, a2 int) error {
	return in.onWeigh(a0, a1, a2)
}

// Add calls the function of Add.
func (in setSetInput) Add(a0 context.Context, a1 string) error {
	return in.onAdd(a0, a1)
}

// Items calls the function of Items.
func (in setSetInput) Items(a0 context.Context) ([]string, error) {
	return in.onItems(a0)
}

// setOverEmbedded is Set under the name of the field that setOver embeds.
type setOverEmbedded = Set

// setOver is a subject of Set whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type setOver struct {
	setOverEmbedded
	spend func()
}

// Merge calls spend and returns its zero results.
func (o setOver) Merge(context.Context, Replica) (r0 error) {
	o.spend()
	return r0
}

// Join calls spend and returns its zero results.
func (o setOver) Join(context.Context, Set) (r0 error) {
	o.spend()
	return r0
}

// Weigh calls spend and returns its zero results.
func (o setOver) Weigh(context.Context, Replica, int) (r0 error) {
	o.spend()
	return r0
}

// Add calls spend and returns its zero results.
func (o setOver) Add(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Items calls spend and returns its zero results.
func (o setOver) Items(context.Context) (r0 []string, r1 error) {
	o.spend()
	return r0, r1
}

// newSetOver returns the planted defect of a ceiling of Set, whose
// callables call spend.
func newSetOver(_ *SetHandle, spend func()) Set {
	return setOver{spend: spend}
}

// setMergeIgnoresTheOtherReplicaDefect is the planted defect of a Set
// whose Merge ignores the other replica.
type setMergeIgnoresTheOtherReplicaDefect struct {
	setOverEmbedded
}

// Merge is the override of the planted defect.
func (*setMergeIgnoresTheOtherReplicaDefect) Merge(p0 context.Context, p1 Replica) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// setMergeIgnoresTheOtherReplica returns the planted defect of a Set whose
// Merge ignores the other replica over a fresh subject.
func setMergeIgnoresTheOtherReplica(h *SetHandle) Set {
	return &setMergeIgnoresTheOtherReplicaDefect{setOverEmbedded: h.Fresh()}
}

// setJoinIgnoresTheOtherReplicaDefect is the planted defect of a Set whose
// Join ignores the other replica.
type setJoinIgnoresTheOtherReplicaDefect struct {
	setOverEmbedded
}

// Join is the override of the planted defect.
func (*setJoinIgnoresTheOtherReplicaDefect) Join(p0 context.Context, p1 Set) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// setJoinIgnoresTheOtherReplica returns the planted defect of a Set whose
// Join ignores the other replica over a fresh subject.
func setJoinIgnoresTheOtherReplica(h *SetHandle) Set {
	return &setJoinIgnoresTheOtherReplicaDefect{setOverEmbedded: h.Fresh()}
}

// setAddPanicsDefect is the planted defect of a Set whose Add panics.
type setAddPanicsDefect struct {
	setOverEmbedded
}

// Add is the override of the planted defect.
func (*setAddPanicsDefect) Add(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Set whose Add panics")
}

// setAddPanics returns the planted defect of a Set whose Add panics over a
// fresh subject.
func setAddPanics(h *SetHandle) Set {
	return &setAddPanicsDefect{setOverEmbedded: h.Fresh()}
}

// setAddIgnoresCancelDefect is the planted defect of a Set whose Add
// ignores the context.
type setAddIgnoresCancelDefect struct {
	setOverEmbedded
}

// Add is the override of the planted defect.
func (d *setAddIgnoresCancelDefect) Add(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.setOverEmbedded.Add(p0, p1)
}

// setAddIgnoresCancel returns the planted defect of a Set whose Add
// ignores the context over a fresh subject.
func setAddIgnoresCancel(h *SetHandle) Set {
	return &setAddIgnoresCancelDefect{setOverEmbedded: h.Fresh()}
}

// setItemsPanicsDefect is the planted defect of a Set whose Items panics.
type setItemsPanicsDefect struct {
	setOverEmbedded
}

// Items is the override of the planted defect.
func (*setItemsPanicsDefect) Items(p0 context.Context) (r0 []string, r1 error) {
	panic("the planted defect of a Set whose Items panics")
}

// setItemsPanics returns the planted defect of a Set whose Items panics
// over a fresh subject.
func setItemsPanics(h *SetHandle) Set {
	return &setItemsPanicsDefect{setOverEmbedded: h.Fresh()}
}

// setItemsIgnoresCancelDefect is the planted defect of a Set whose Items
// ignores the context.
type setItemsIgnoresCancelDefect struct {
	setOverEmbedded
}

// Items is the override of the planted defect.
func (d *setItemsIgnoresCancelDefect) Items(p0 context.Context) (r0 []string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.setOverEmbedded.Items(p0)
}

// setItemsIgnoresCancel returns the planted defect of a Set whose Items
// ignores the context over a fresh subject.
func setItemsIgnoresCancel(h *SetHandle) Set {
	return &setItemsIgnoresCancelDefect{setOverEmbedded: h.Fresh()}
}

// setItemsValueBesideDefect is the planted defect of a Set whose Items
// returns a value beside an error.
type setItemsValueBesideDefect struct {
	setOverEmbedded
}

// Items is the override of the planted defect.
func (d *setItemsValueBesideDefect) Items(p0 context.Context) (r0 []string, r1 error) {
	r0, r1 = d.setOverEmbedded.Items(p0)
	if r1 != nil {
		r0 = []string{"sample-b"}
	}
	return r0, r1
}

// setItemsValueBeside returns the planted defect of a Set whose Items
// returns a value beside an error over a fresh subject.
func setItemsValueBeside(h *SetHandle) Set {
	return &setItemsValueBesideDefect{setOverEmbedded: h.Fresh()}
}

// setActions returns the actions of the generated sequences of Set over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func setActions(c *prop.Case, subject func() Set) []stateful.Action[struct{}] {
	keysReplica := c.Draw(
		prop.List(prop.Of[Replica](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.replica",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Merge",
			Input: func(c *prop.Case, _ struct{}) any {
				peerItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "peer-items-0")
				peerItemsResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, conform.ErrGenerated),
						prop.MinSize(len(peerItemsResults0)),
						prop.MaxSize(len(peerItemsResults0)),
					),
					"peer-items-1",
				)
				var peerItemsAt atomic.Int64
				peer := &setReplicaInput{
					onItems: func(context.Context) (r0 []string, r1 error) {
						i := int(peerItemsAt.Add(1) - 1)
						if i >= len(peerItemsResults0) {
							r1 = conform.ErrGenerated
							return r0, r1
						}
						return peerItemsResults0[i],
							peerItemsResults1[i]
					},
				}
				return []any{peer}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Merge", args)
				values, err := setMergeCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Join",
			Input: func(c *prop.Case, _ struct{}) any {
				otherMergeResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-merge",
				)
				var otherMergeAt atomic.Int64
				otherJoinResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-join",
				)
				var otherJoinAt atomic.Int64
				otherWeighResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-weigh",
				)
				var otherWeighAt atomic.Int64
				otherAddResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-add",
				)
				var otherAddAt atomic.Int64
				otherItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "other-items-0")
				otherItemsResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, conform.ErrGenerated),
						prop.MinSize(len(otherItemsResults0)),
						prop.MaxSize(len(otherItemsResults0)),
					),
					"other-items-1",
				)
				var otherItemsAt atomic.Int64
				other := &setSetInput{
					onMerge: func(context.Context, Replica) (r0 error) {
						i := int(otherMergeAt.Add(1) - 1)
						if i >= len(otherMergeResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherMergeResults[i]
					},
					onJoin: func(context.Context, Set) (r0 error) {
						i := int(otherJoinAt.Add(1) - 1)
						if i >= len(otherJoinResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherJoinResults[i]
					},
					onWeigh: func(context.Context, Replica, int) (r0 error) {
						i := int(otherWeighAt.Add(1) - 1)
						if i >= len(otherWeighResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherWeighResults[i]
					},
					onAdd: func(context.Context, string) (r0 error) {
						i := int(otherAddAt.Add(1) - 1)
						if i >= len(otherAddResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherAddResults[i]
					},
					onItems: func(context.Context) (r0 []string, r1 error) {
						i := int(otherItemsAt.Add(1) - 1)
						if i >= len(otherItemsResults0) {
							r1 = conform.ErrGenerated
							return r0, r1
						}
						return otherItemsResults0[i],
							otherItemsResults1[i]
					},
				}
				return []any{other}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Join", args)
				values, err := setJoinCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Weigh",
			Input: func(c *prop.Case, _ struct{}) any {
				peer := c.Draw(prop.SampledFrom(keysReplica...), "weigh.peer")
				weight := c.Draw(prop.Of[int](), "weigh.weight")
				return []any{peer, weight}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Weigh", args, args[0])
				values, err := setWeighCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Add",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.Of[string](), "add.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Add", args)
				values, err := setAddCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Items",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Items", args)
				values, err := setItemsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// setMergeCall calls Merge on subject with the arguments args, and returns
// the values of the call beside its error.
func setMergeCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(Replica)
	return nil, subject.Merge(ctx, a0)
}

// setJoinCall calls Join on subject with the arguments args, and returns
// the values of the call beside its error.
func setJoinCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(Set)
	return nil, subject.Join(ctx, a0)
}

// setWeighCall calls Weigh on subject with the arguments args, and returns
// the values of the call beside its error.
func setWeighCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(Replica)
	a1, _ := args[1].(int)
	return nil, subject.Weigh(ctx, a0, a1)
}

// setAddCall calls Add on subject with the arguments args, and returns the
// values of the call beside its error.
func setAddCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Add(ctx, a0)
}

// setItemsCall calls Items on subject with the arguments args, and returns
// the values of the call beside its error.
func setItemsCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	r0, err := subject.Items(ctx)
	return []any{r0}, err
}

// setReference returns the function that applies an operation of a history
// to the reference r by the name of its callable.
func setReference(r Set) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Merge":
			values, err = setMergeCall(ctx, r, args)
		case "Join":
			values, err = setJoinCall(ctx, r, args)
		case "Weigh":
			values, err = setWeighCall(ctx, r, args)
		case "Add":
			values, err = setAddCall(ctx, r, args)
		case "Items":
			values, err = setItemsCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Tally:
//
// Tally.Merge/merges-in-any-order: two replicas that Add made diverge read
// the same state through Count after Merge merges them in either order.
//
// Tally.Add/smoke: Add returns without a panic for derived inputs.
//
// Tally.Add/cancel: Add returns an error for a cancelled context.
//
// Tally.Count/smoke: Count returns without a panic for derived inputs.
//
// Tally.Count/cancel: Count returns an error for a cancelled context.
//
// Tally.Count/zero-on-error: Count returns the zero value beside an error.

// The derived checks of Tally that witness their laws at the shallow end:
//
// Tally.Add/smoke
//
// Tally.Count/smoke

// The refused claims of Tally:
//
// Tally.Merge/cancel: eidos derives no sample of the parameter other of
// Merge. The remedy is to author a sample of the parameter other with the
// directive sample.

// The unclaimed classifications of Tally:
//
// The shape:writer of Tally.Merge has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Tally.Add has the stance silent. A lifecycle call
// returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Tally writes.
type (
	TallyHandle = conform.Handle[TallyMethods]
	TallyCheck  = conform.Check[TallyMethods]
	TallyOption = conform.Option[TallyMethods]
)

// TallyMethods is the interface of the exported methods of Tally. Each
// subject of the suite of Tally implements it.
type TallyMethods interface {
	Merge(ctx context.Context, other *Tally) error
	Add(ctx context.Context) error
	Count(ctx context.Context) (int, error)
}

// TallyChecks contains the name of each callable of Tally, with the
// identities of its checks.
var TallyChecks = struct {
	Merge tallyMerge
	Add   tallyAdd
	Count tallyCount
}{
	Merge: tallyMerge{"Tally.Merge"},
	Add:   tallyAdd{"Tally.Add"},
	Count: tallyCount{"Tally.Count"},
}

// tallyMerge is the name of Tally.Merge, with the identities of its
// checks.
type tallyMerge struct {
	conform.Callable[TallyMethods]
}

// MergesInAnyOrder returns the identity of the check that two replicas
// that Add made diverge read the same state through Count after Merge
// merges them in either order.
func (c tallyMerge) MergesInAnyOrder() conform.ID[TallyMethods] {
	return c.Check("merges-in-any-order")
}

// tallyAdd is the name of Tally.Add, with the identities of its checks.
type tallyAdd struct {
	conform.Callable[TallyMethods]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c tallyAdd) Smoke() conform.ID[TallyMethods] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c tallyAdd) Cancel() conform.ID[TallyMethods] {
	return c.Check("cancel")
}

// tallyCount is the name of Tally.Count, with the identities of its
// checks.
type tallyCount struct {
	conform.Callable[TallyMethods]
}

// Smoke returns the identity of the check that Count returns without a
// panic for derived inputs.
func (c tallyCount) Smoke() conform.ID[TallyMethods] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Count returns an error for
// a cancelled context.
func (c tallyCount) Cancel() conform.ID[TallyMethods] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Count returns the
// zero value beside an error.
func (c tallyCount) ZeroOnError() conform.ID[TallyMethods] {
	return c.Check("zero-on-error")
}

// TallySuite returns the checks of Tally, and the calls and the budgets of
// its callables.
func TallySuite() conform.Suite[TallyMethods] {
	return conform.Suite[TallyMethods]{
		Declaration: "Tally",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[TallyMethods]{
			{
				ID:          TallyChecks.Merge.MergesInAnyOrder(),
				Statement:   "two replicas that Add made diverge read the same state through Count after Merge merges them in either order",
				Claim:       spec.MergesInAnyOrder,
				Probe:       spec.MergesCommute,
				Reason:      spec.OrderDependentMerge,
				Law:         spec.MergeOrderDoesNotMatter,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *TallyHandle, c *prop.Case) {
					probe.MergesCommute(h, c, tallyActions, tallyReference, "Merge", "Add", "Count")
				},
				Defect: tallyMergeIgnoresTheOtherReplica,
			},
			{
				ID:          TallyChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TallyHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: tallyAddPanics,
			},
			{
				ID:          TallyChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TallyHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: tallyAddIgnoresCancel,
			},
			{
				ID:          TallyChecks.Count.Smoke(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TallyHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: tallyCountPanics,
			},
			{
				ID:          TallyChecks.Count.Cancel(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TallyHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: tallyCountIgnoresCancel,
			},
			{
				ID:          TallyChecks.Count.ZeroOnError(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TallyHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx)
					}, int(0))
				},
				Defect: tallyCountValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Tally.Merge",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter other of Merge",
				Remedy:   "author a sample of the parameter other with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Tally.Merge",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Tally.Add",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[TallyMethods]{
			{
				Calls:    tallyMergeCalls,
				Over:     newTallyOver,
				Callable: TallyChecks.Merge.Callable,
			},
			{
				Calls:    tallyAddCalls,
				Over:     newTallyOver,
				Callable: TallyChecks.Add.Callable,
			},
			{
				Calls:    tallyCountCalls,
				Over:     newTallyOver,
				Callable: TallyChecks.Count.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(tallyTallyValues),
		},
	}
}

// RunTallySuite applies the checks of Tally, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunTallySuite[T TallyMethods](t *testing.T, newSubject func() T, opts ...TallyOption) {
	t.Helper()
	RunTallySuiteFromHandle(t, func(*TallyHandle) T { return newSubject() }, opts...)
}

// RunTallySuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTallySuiteFromHandle[T TallyMethods](
	t *testing.T,
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	t.Helper()
	conform.Run(t, TallySuite(), tallySubjects(newSubject), opts...)
}

// FuzzTallySuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTallySuite[T TallyMethods](
	f *testing.F,
	callable conform.Callable[TallyMethods],
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	f.Helper()
	conform.Fuzz(f, TallySuite(), callable, tallySubjects(newSubject), opts...)
}

// RunTallyBudgets applies the allocation ceilings of Tally and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunTallyBudgets[T TallyMethods](
	t *testing.T,
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	t.Helper()
	conform.RunBudgets(t, TallySuite(), tallySubjects(newSubject), opts...)
}

// BenchTallySuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchTallySuite[T TallyMethods](
	b conform.B,
	callable conform.Callable[TallyMethods],
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	b.Helper()
	conform.Bench(b, TallySuite(), callable, tallySubjects(newSubject), opts...)
}

// tallySubjects returns the factory of the subjects of Tally that calls
// newSubject with the handle of each check.
func tallySubjects[T TallyMethods](
	newSubject func(h *TallyHandle) T,
) func(h *TallyHandle) TallyMethods {
	return func(h *TallyHandle) TallyMethods { return newSubject(h) }
}

// tallyMergeCalls generates the calls of Tally.Merge with the inputs that
// each case draws.
var tallyMergeCalls = prop.Composite(func(c *prop.Case) func(context.Context, TallyMethods) error {
	other := c.Draw(prop.Of[Tally](), "other")
	return func(ctx context.Context, subject TallyMethods) error {
		return subject.Merge(ctx, &other)
	}
})

// tallyAddCalls generates the calls of Tally.Add with the inputs that each
// case draws.
var tallyAddCalls = prop.Composite(func(*prop.Case) func(context.Context, TallyMethods) error {
	return func(ctx context.Context, subject TallyMethods) error {
		return subject.Add(ctx)
	}
})

// tallyCountCalls generates the calls of Tally.Count with the inputs that
// each case draws.
var tallyCountCalls = prop.Composite(func(*prop.Case) func(context.Context, TallyMethods) error {
	return func(ctx context.Context, subject TallyMethods) error {
		_, err := subject.Count(ctx)
		return err
	}
})

// tallyTallyValues generates the starting value of Tally, which has no
// field to change.
var tallyTallyValues = prop.Composite(func(*prop.Case) Tally {
	return Tally{}
})

// tallyOverEmbedded is TallyMethods under the name of the field that
// tallyOver embeds.
type tallyOverEmbedded = TallyMethods

// tallyOver is a subject of Tally whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type tallyOver struct {
	tallyOverEmbedded
	spend func()
}

// Merge calls spend and returns its zero results.
func (o tallyOver) Merge(context.Context, *Tally) (r0 error) {
	o.spend()
	return r0
}

// Add calls spend and returns its zero results.
func (o tallyOver) Add(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o tallyOver) Count(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newTallyOver returns the planted defect of a ceiling of Tally, whose
// callables call spend.
func newTallyOver(_ *TallyHandle, spend func()) TallyMethods {
	return tallyOver{spend: spend}
}

// tallyMergeIgnoresTheOtherReplicaDefect is the planted defect of a Tally
// whose Merge ignores the other replica.
type tallyMergeIgnoresTheOtherReplicaDefect struct {
	tallyOverEmbedded
}

// Merge is the override of the planted defect.
func (*tallyMergeIgnoresTheOtherReplicaDefect) Merge(p0 context.Context, p1 *Tally) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// tallyMergeIgnoresTheOtherReplica returns the planted defect of a Tally
// whose Merge ignores the other replica over a fresh subject.
func tallyMergeIgnoresTheOtherReplica(h *TallyHandle) TallyMethods {
	return &tallyMergeIgnoresTheOtherReplicaDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyAddPanicsDefect is the planted defect of a Tally whose Add panics.
type tallyAddPanicsDefect struct {
	tallyOverEmbedded
}

// Add is the override of the planted defect.
func (*tallyAddPanicsDefect) Add(p0 context.Context) (r0 error) {
	panic("the planted defect of a Tally whose Add panics")
}

// tallyAddPanics returns the planted defect of a Tally whose Add panics
// over a fresh subject.
func tallyAddPanics(h *TallyHandle) TallyMethods {
	return &tallyAddPanicsDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyAddIgnoresCancelDefect is the planted defect of a Tally whose Add
// ignores the context.
type tallyAddIgnoresCancelDefect struct {
	tallyOverEmbedded
}

// Add is the override of the planted defect.
func (d *tallyAddIgnoresCancelDefect) Add(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.tallyOverEmbedded.Add(p0)
}

// tallyAddIgnoresCancel returns the planted defect of a Tally whose Add
// ignores the context over a fresh subject.
func tallyAddIgnoresCancel(h *TallyHandle) TallyMethods {
	return &tallyAddIgnoresCancelDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyCountPanicsDefect is the planted defect of a Tally whose Count
// panics.
type tallyCountPanicsDefect struct {
	tallyOverEmbedded
}

// Count is the override of the planted defect.
func (*tallyCountPanicsDefect) Count(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Tally whose Count panics")
}

// tallyCountPanics returns the planted defect of a Tally whose Count
// panics over a fresh subject.
func tallyCountPanics(h *TallyHandle) TallyMethods {
	return &tallyCountPanicsDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyCountIgnoresCancelDefect is the planted defect of a Tally whose
// Count ignores the context.
type tallyCountIgnoresCancelDefect struct {
	tallyOverEmbedded
}

// Count is the override of the planted defect.
func (d *tallyCountIgnoresCancelDefect) Count(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.tallyOverEmbedded.Count(p0)
}

// tallyCountIgnoresCancel returns the planted defect of a Tally whose
// Count ignores the context over a fresh subject.
func tallyCountIgnoresCancel(h *TallyHandle) TallyMethods {
	return &tallyCountIgnoresCancelDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyCountValueBesideDefect is the planted defect of a Tally whose Count
// returns a value beside an error.
type tallyCountValueBesideDefect struct {
	tallyOverEmbedded
}

// Count is the override of the planted defect.
func (d *tallyCountValueBesideDefect) Count(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.tallyOverEmbedded.Count(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// tallyCountValueBeside returns the planted defect of a Tally whose Count
// returns a value beside an error over a fresh subject.
func tallyCountValueBeside(h *TallyHandle) TallyMethods {
	return &tallyCountValueBesideDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyActions returns the actions of the generated sequences of Tally
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func tallyActions(c *prop.Case, subject func() TallyMethods) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Merge",
			Input: func(c *prop.Case, _ struct{}) any {
				other := c.Draw(prop.Of[Tally](), "merge.other")
				return []any{&other}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Merge", args)
				values, err := tallyMergeCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Add",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Add", args)
				values, err := tallyAddCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Count",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Count", args)
				values, err := tallyCountCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// tallyMergeCall calls Merge on subject with the arguments args, and
// returns the values of the call beside its error.
func tallyMergeCall(ctx context.Context, subject TallyMethods, args []any) ([]any, error) {
	a0, _ := args[0].(*Tally)
	return nil, subject.Merge(ctx, a0)
}

// tallyAddCall calls Add on subject with the arguments args, and returns
// the values of the call beside its error.
func tallyAddCall(ctx context.Context, subject TallyMethods, args []any) ([]any, error) {
	return nil, subject.Add(ctx)
}

// tallyCountCall calls Count on subject with the arguments args, and
// returns the values of the call beside its error.
func tallyCountCall(ctx context.Context, subject TallyMethods, args []any) ([]any, error) {
	r0, err := subject.Count(ctx)
	return []any{r0}, err
}

// tallyReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func tallyReference(r TallyMethods) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Merge":
			values, err = tallyMergeCall(ctx, r, args)
		case "Add":
			values, err = tallyAddCall(ctx, r, args)
		case "Count":
			values, err = tallyCountCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Pile:
//
// Pile.Absorb/smoke: Absorb returns without a panic for derived inputs.
//
// Pile.Absorb/smoke#generated: Absorb returns without a panic for derived
// inputs.
//
// Pile.Absorb/cancel: Absorb returns an error for a cancelled context.
//
// Pile.Absorb/cancel#generated: Absorb returns an error for a cancelled
// context.
//
// Pile.Add/smoke: Add returns without a panic for derived inputs.
//
// Pile.Add/cancel: Add returns an error for a cancelled context.
//
// Pile.Count/smoke: Count returns without a panic for derived inputs.
//
// Pile.Count/cancel: Count returns an error for a cancelled context.
//
// Pile.Count/zero-on-error: Count returns the zero value beside an error.

// The derived checks of Pile that witness their laws at the shallow end:
//
// Pile.Absorb/smoke
//
// Pile.Add/smoke
//
// Pile.Count/smoke

// The refused claims of Pile:
//
// Pile.Pour/cancel: eidos derives no sample of the parameter into of Pour.
// The remedy is to author a sample of the parameter into with the
// directive sample.
//
// Pile.Absorb/merges-in-any-order: Absorb takes no input that another
// replica of Pile can be. The remedy is to give Absorb an input of an
// interface that Pile implements.
//
// Pile.Pour/merges-in-any-order: Pour takes no input that another replica
// of Pile can be. The remedy is to give Pour an input of an interface that
// Pile implements.

// The unclaimed classifications of Pile:
//
// The shape:writer of Pile.Absorb has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Pile.Pour has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Pile.Add has the stance silent. A lifecycle call
// returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Pile writes.
type (
	PileHandle = conform.Handle[Pile]
	PileCheck  = conform.Check[Pile]
	PileOption = conform.Option[Pile]
)

// PileChecks contains the name of each callable of Pile, with the
// identities of its checks.
var PileChecks = struct {
	Absorb pileAbsorb
	Pour   pilePour
	Add    pileAdd
	Count  pileCount
}{
	Absorb: pileAbsorb{"Pile.Absorb"},
	Pour:   pilePour{"Pile.Pour"},
	Add:    pileAdd{"Pile.Add"},
	Count:  pileCount{"Pile.Count"},
}

// pileAbsorb is the name of Pile.Absorb, with the identities of its
// checks.
type pileAbsorb struct {
	conform.Callable[Pile]
}

// Smoke returns the identity of the check that Absorb returns without a
// panic for derived inputs.
func (c pileAbsorb) Smoke() conform.ID[Pile] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Absorb returns
// without a panic for derived inputs.
func (c pileAbsorb) SmokeGenerated() conform.ID[Pile] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Absorb returns an error
// for a cancelled context.
func (c pileAbsorb) Cancel() conform.ID[Pile] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Absorb returns an
// error for a cancelled context.
func (c pileAbsorb) CancelGenerated() conform.ID[Pile] {
	return c.Check("cancel#generated")
}

// pilePour is the name of Pile.Pour, with the identities of its checks.
type pilePour struct {
	conform.Callable[Pile]
}

// pileAdd is the name of Pile.Add, with the identities of its checks.
type pileAdd struct {
	conform.Callable[Pile]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c pileAdd) Smoke() conform.ID[Pile] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c pileAdd) Cancel() conform.ID[Pile] {
	return c.Check("cancel")
}

// pileCount is the name of Pile.Count, with the identities of its checks.
type pileCount struct {
	conform.Callable[Pile]
}

// Smoke returns the identity of the check that Count returns without a
// panic for derived inputs.
func (c pileCount) Smoke() conform.ID[Pile] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Count returns an error for
// a cancelled context.
func (c pileCount) Cancel() conform.ID[Pile] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Count returns the
// zero value beside an error.
func (c pileCount) ZeroOnError() conform.ID[Pile] {
	return c.Check("zero-on-error")
}

// PileSuite returns the checks of Pile, and the calls and the budgets of
// its callables.
func PileSuite() conform.Suite[Pile] {
	return conform.Suite[Pile]{
		Declaration: "Pile",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Pile]{
			{
				ID:          PileChecks.Absorb.Smoke(),
				Statement:   "Absorb returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PileHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, int(42))
					})
				},
				Defect: pileAbsorbPanics,
			},
			{
				ID:          PileChecks.Absorb.SmokeGenerated(),
				Statement:   "Absorb returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PileHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, drawnN)
					})
				},
				Defect: pileAbsorbPanics,
			},
			{
				ID:          PileChecks.Absorb.Cancel(),
				Statement:   "Absorb returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, int(42))
					})
				},
				Defect: pileAbsorbIgnoresCancel,
			},
			{
				ID:          PileChecks.Absorb.CancelGenerated(),
				Statement:   "Absorb returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PileHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, drawnN)
					})
				},
				Defect: pileAbsorbIgnoresCancel,
			},
			{
				ID:          PileChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PileHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: pileAddPanics,
			},
			{
				ID:          PileChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: pileAddIgnoresCancel,
			},
			{
				ID:          PileChecks.Count.Smoke(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PileHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: pileCountPanics,
			},
			{
				ID:          PileChecks.Count.Cancel(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: pileCountIgnoresCancel,
			},
			{
				ID:          PileChecks.Count.ZeroOnError(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx)
					}, int(0))
				},
				Defect: pileCountValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Pile.Pour",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter into of Pour",
				Remedy:   "author a sample of the parameter into with the directive sample",
			},
			{
				Claim:    spec.MergesInAnyOrder,
				Callable: "Pile.Absorb",
				Gate:     spec.TakesAReplica,
				Why:      "Absorb takes no input that another replica of Pile can be",
				Remedy:   "give Absorb an input of an interface that Pile implements",
			},
			{
				Claim:    spec.MergesInAnyOrder,
				Callable: "Pile.Pour",
				Gate:     spec.TakesAReplica,
				Why:      "Pour takes no input that another replica of Pile can be",
				Remedy:   "give Pour an input of an interface that Pile implements",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Pile.Absorb",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Pile.Pour",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Pile.Add",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Pile]{
			{
				Calls:    pileAbsorbCalls,
				Over:     newPileOver,
				Callable: PileChecks.Absorb.Callable,
			},
			{
				Calls:    pilePourCalls,
				Over:     newPileOver,
				Callable: PileChecks.Pour.Callable,
			},
			{
				Calls:    pileAddCalls,
				Over:     newPileOver,
				Callable: PileChecks.Add.Callable,
			},
			{
				Calls:    pileCountCalls,
				Over:     newPileOver,
				Callable: PileChecks.Count.Callable,
			},
		},
	}
}

// RunPileSuite applies the checks of Pile, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunPileSuite[T Pile](t *testing.T, newSubject func() T, opts ...PileOption) {
	t.Helper()
	RunPileSuiteFromHandle(t, func(*PileHandle) T { return newSubject() }, opts...)
}

// RunPileSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPileSuiteFromHandle[T Pile](
	t *testing.T,
	newSubject func(h *PileHandle) T,
	opts ...PileOption,
) {
	t.Helper()
	conform.Run(t, PileSuite(), pileSubjects(newSubject), opts...)
}

// FuzzPileSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPileSuite[T Pile](
	f *testing.F,
	callable conform.Callable[Pile],
	newSubject func(h *PileHandle) T,
	opts ...PileOption,
) {
	f.Helper()
	conform.Fuzz(f, PileSuite(), callable, pileSubjects(newSubject), opts...)
}

// RunPileBudgets applies the allocation ceilings of Pile and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunPileBudgets[T Pile](t *testing.T, newSubject func(h *PileHandle) T, opts ...PileOption) {
	t.Helper()
	conform.RunBudgets(t, PileSuite(), pileSubjects(newSubject), opts...)
}

// BenchPileSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPileSuite[T Pile](
	b conform.B,
	callable conform.Callable[Pile],
	newSubject func(h *PileHandle) T,
	opts ...PileOption,
) {
	b.Helper()
	conform.Bench(b, PileSuite(), callable, pileSubjects(newSubject), opts...)
}

// pileSubjects returns the factory of the subjects of Pile that calls
// newSubject with the handle of each check.
func pileSubjects[T Pile](newSubject func(h *PileHandle) T) func(h *PileHandle) Pile {
	return func(h *PileHandle) Pile { return newSubject(h) }
}

// pileAbsorbCalls generates the calls of Pile.Absorb with the inputs that
// each case draws.
var pileAbsorbCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pile) error {
	n := c.Draw(prop.Of[int](), "n")
	return func(ctx context.Context, subject Pile) error {
		return subject.Absorb(ctx, n)
	}
})

// pilePourCalls generates the calls of Pile.Pour with the inputs that each
// case draws.
var pilePourCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pile) error {
	intoDrainResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"into-drain",
	)
	var intoDrainAt atomic.Int64
	into := &pileSinkInput{
		onDrain: func(context.Context) (r0 error) {
			i := int(intoDrainAt.Add(1) - 1)
			if i >= len(intoDrainResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return intoDrainResults[i]
		},
	}
	return func(ctx context.Context, subject Pile) error {
		return subject.Pour(ctx, into)
	}
})

// pileAddCalls generates the calls of Pile.Add with the inputs that each
// case draws.
var pileAddCalls = prop.Composite(func(*prop.Case) func(context.Context, Pile) error {
	return func(ctx context.Context, subject Pile) error {
		return subject.Add(ctx)
	}
})

// pileCountCalls generates the calls of Pile.Count with the inputs that
// each case draws.
var pileCountCalls = prop.Composite(func(*prop.Case) func(context.Context, Pile) error {
	return func(ctx context.Context, subject Pile) error {
		_, err := subject.Count(ctx)
		return err
	}
})

// pileSinkInputEmbedded is Sink under the name of the field that
// pileSinkInput embeds.
type pileSinkInputEmbedded = Sink

// pileSinkInput implements Sink. Its methods return the results that the
// generators of the suite of Pile draw for them, in turn.
type pileSinkInput struct {
	pileSinkInputEmbedded
	onDrain func(a0 context.Context) error
}

// Drain calls the function of Drain.
func (in pileSinkInput) Drain(a0 context.Context) error {
	return in.onDrain(a0)
}

// pileOverEmbedded is Pile under the name of the field that pileOver
// embeds.
type pileOverEmbedded = Pile

// pileOver is a subject of Pile whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type pileOver struct {
	pileOverEmbedded
	spend func()
}

// Absorb calls spend and returns its zero results.
func (o pileOver) Absorb(context.Context, int) (r0 error) {
	o.spend()
	return r0
}

// Pour calls spend and returns its zero results.
func (o pileOver) Pour(context.Context, Sink) (r0 error) {
	o.spend()
	return r0
}

// Add calls spend and returns its zero results.
func (o pileOver) Add(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o pileOver) Count(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newPileOver returns the planted defect of a ceiling of Pile, whose
// callables call spend.
func newPileOver(_ *PileHandle, spend func()) Pile {
	return pileOver{spend: spend}
}

// pileAbsorbPanicsDefect is the planted defect of a Pile whose Absorb
// panics.
type pileAbsorbPanicsDefect struct {
	pileOverEmbedded
}

// Absorb is the override of the planted defect.
func (*pileAbsorbPanicsDefect) Absorb(p0 context.Context, p1 int) (r0 error) {
	panic("the planted defect of a Pile whose Absorb panics")
}

// pileAbsorbPanics returns the planted defect of a Pile whose Absorb
// panics over a fresh subject.
func pileAbsorbPanics(h *PileHandle) Pile {
	return &pileAbsorbPanicsDefect{pileOverEmbedded: h.Fresh()}
}

// pileAbsorbIgnoresCancelDefect is the planted defect of a Pile whose
// Absorb ignores the context.
type pileAbsorbIgnoresCancelDefect struct {
	pileOverEmbedded
}

// Absorb is the override of the planted defect.
func (d *pileAbsorbIgnoresCancelDefect) Absorb(p0 context.Context, p1 int) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.pileOverEmbedded.Absorb(p0, p1)
}

// pileAbsorbIgnoresCancel returns the planted defect of a Pile whose
// Absorb ignores the context over a fresh subject.
func pileAbsorbIgnoresCancel(h *PileHandle) Pile {
	return &pileAbsorbIgnoresCancelDefect{pileOverEmbedded: h.Fresh()}
}

// pileAddPanicsDefect is the planted defect of a Pile whose Add panics.
type pileAddPanicsDefect struct {
	pileOverEmbedded
}

// Add is the override of the planted defect.
func (*pileAddPanicsDefect) Add(p0 context.Context) (r0 error) {
	panic("the planted defect of a Pile whose Add panics")
}

// pileAddPanics returns the planted defect of a Pile whose Add panics over
// a fresh subject.
func pileAddPanics(h *PileHandle) Pile {
	return &pileAddPanicsDefect{pileOverEmbedded: h.Fresh()}
}

// pileAddIgnoresCancelDefect is the planted defect of a Pile whose Add
// ignores the context.
type pileAddIgnoresCancelDefect struct {
	pileOverEmbedded
}

// Add is the override of the planted defect.
func (d *pileAddIgnoresCancelDefect) Add(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.pileOverEmbedded.Add(p0)
}

// pileAddIgnoresCancel returns the planted defect of a Pile whose Add
// ignores the context over a fresh subject.
func pileAddIgnoresCancel(h *PileHandle) Pile {
	return &pileAddIgnoresCancelDefect{pileOverEmbedded: h.Fresh()}
}

// pileCountPanicsDefect is the planted defect of a Pile whose Count
// panics.
type pileCountPanicsDefect struct {
	pileOverEmbedded
}

// Count is the override of the planted defect.
func (*pileCountPanicsDefect) Count(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Pile whose Count panics")
}

// pileCountPanics returns the planted defect of a Pile whose Count panics
// over a fresh subject.
func pileCountPanics(h *PileHandle) Pile {
	return &pileCountPanicsDefect{pileOverEmbedded: h.Fresh()}
}

// pileCountIgnoresCancelDefect is the planted defect of a Pile whose Count
// ignores the context.
type pileCountIgnoresCancelDefect struct {
	pileOverEmbedded
}

// Count is the override of the planted defect.
func (d *pileCountIgnoresCancelDefect) Count(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.pileOverEmbedded.Count(p0)
}

// pileCountIgnoresCancel returns the planted defect of a Pile whose Count
// ignores the context over a fresh subject.
func pileCountIgnoresCancel(h *PileHandle) Pile {
	return &pileCountIgnoresCancelDefect{pileOverEmbedded: h.Fresh()}
}

// pileCountValueBesideDefect is the planted defect of a Pile whose Count
// returns a value beside an error.
type pileCountValueBesideDefect struct {
	pileOverEmbedded
}

// Count is the override of the planted defect.
func (d *pileCountValueBesideDefect) Count(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.pileOverEmbedded.Count(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// pileCountValueBeside returns the planted defect of a Pile whose Count
// returns a value beside an error over a fresh subject.
func pileCountValueBeside(h *PileHandle) Pile {
	return &pileCountValueBesideDefect{pileOverEmbedded: h.Fresh()}
}

// The derived checks of Run:
//
// Run.Advance/smoke: Advance returns without a panic for derived inputs.
//
// Run.Advance/smoke#generated: Advance returns without a panic for derived
// inputs.
//
// Run.Advance/cancel: Advance returns an error for a cancelled context.
//
// Run.Advance/cancel#generated: Advance returns an error for a cancelled
// context.
//
// Run.Advance/follows-its-graph: each call of Advance leaves the state
// that Status reads at the state before it or at one of its successors,
// over any sequence of calls.
//
// Run.Retry/smoke: Retry returns without a panic for derived inputs.
//
// Run.Retry/smoke#generated: Retry returns without a panic for derived
// inputs.
//
// Run.Retry/cancel: Retry returns an error for a cancelled context.
//
// Run.Retry/cancel#generated: Retry returns an error for a cancelled
// context.
//
// Run.Retry/follows-its-graph: each call of Retry leaves the state that
// Status reads at the state before it or at one of its successors, over
// any sequence of calls.
//
// Run.Restart/smoke: Restart returns without a panic for derived inputs.
//
// Run.Restart/smoke#generated: Restart returns without a panic for derived
// inputs.
//
// Run.Restart/cancel: Restart returns an error for a cancelled context.
//
// Run.Restart/cancel#generated: Restart returns an error for a cancelled
// context.
//
// Run.Toggle/smoke: Toggle returns without a panic for derived inputs.
//
// Run.Toggle/smoke#generated: Toggle returns without a panic for derived
// inputs.
//
// Run.Toggle/cancel: Toggle returns an error for a cancelled context.
//
// Run.Toggle/cancel#generated: Toggle returns an error for a cancelled
// context.
//
// Run.Toggle/follows-its-graph: each call of Toggle leaves the state that
// Status reads at the state before it or at one of its successors, over
// any sequence of calls.
//
// Run.Status/smoke: Status returns without a panic for derived inputs.
//
// Run.Status/smoke#generated: Status returns without a panic for derived
// inputs.
//
// Run.Status/cancel: Status returns an error for a cancelled context.
//
// Run.Status/cancel#generated: Status returns an error for a cancelled
// context.

// The derived checks of Run that witness their laws at the shallow end:
//
// Run.Advance/smoke
//
// Run.Retry/smoke
//
// Run.Restart/smoke
//
// Run.Toggle/smoke
//
// Run.Status/smoke

// The refused claims of Run:
//
// Run.Restart/follows-its-graph: the variable that transitions names for
// Restart is no graph of the states that observe reads. The remedy is to
// declare transitions as a map from each state to the list of its
// successors.

// The unclaimed classifications of Run:
//
// The shape:writer of Run.Advance has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Run.Retry has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:writer of Run.Restart has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Run.Toggle has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Run writes.
type (
	RunHandle = conform.Handle[Run]
	RunCheck  = conform.Check[Run]
	RunOption = conform.Option[Run]
)

// RunChecks contains the name of each callable of Run, with the identities
// of its checks.
var RunChecks = struct {
	Advance runAdvance
	Retry   runRetry
	Restart runRestart
	Toggle  runToggle
	Status  runStatus
}{
	Advance: runAdvance{"Run.Advance"},
	Retry:   runRetry{"Run.Retry"},
	Restart: runRestart{"Run.Restart"},
	Toggle:  runToggle{"Run.Toggle"},
	Status:  runStatus{"Run.Status"},
}

// runAdvance is the name of Run.Advance, with the identities of its
// checks.
type runAdvance struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Advance returns without a
// panic for derived inputs.
func (c runAdvance) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Advance returns
// without a panic for derived inputs.
func (c runAdvance) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Advance returns an error
// for a cancelled context.
func (c runAdvance) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Advance returns
// an error for a cancelled context.
func (c runAdvance) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// FollowsItsGraph returns the identity of the check that each call of
// Advance leaves the state that Status reads at the state before it or at
// one of its successors, over any sequence of calls.
func (c runAdvance) FollowsItsGraph() conform.ID[Run] {
	return c.Check("follows-its-graph")
}

// runRetry is the name of Run.Retry, with the identities of its checks.
type runRetry struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Retry returns without a
// panic for derived inputs.
func (c runRetry) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Retry returns
// without a panic for derived inputs.
func (c runRetry) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Retry returns an error for
// a cancelled context.
func (c runRetry) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Retry returns an
// error for a cancelled context.
func (c runRetry) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// FollowsItsGraph returns the identity of the check that each call of
// Retry leaves the state that Status reads at the state before it or at
// one of its successors, over any sequence of calls.
func (c runRetry) FollowsItsGraph() conform.ID[Run] {
	return c.Check("follows-its-graph")
}

// runRestart is the name of Run.Restart, with the identities of its
// checks.
type runRestart struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Restart returns without a
// panic for derived inputs.
func (c runRestart) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Restart returns
// without a panic for derived inputs.
func (c runRestart) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Restart returns an error
// for a cancelled context.
func (c runRestart) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Restart returns
// an error for a cancelled context.
func (c runRestart) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// runToggle is the name of Run.Toggle, with the identities of its checks.
type runToggle struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Toggle returns without a
// panic for derived inputs.
func (c runToggle) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Toggle returns
// without a panic for derived inputs.
func (c runToggle) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Toggle returns an error
// for a cancelled context.
func (c runToggle) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Toggle returns an
// error for a cancelled context.
func (c runToggle) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// FollowsItsGraph returns the identity of the check that each call of
// Toggle leaves the state that Status reads at the state before it or at
// one of its successors, over any sequence of calls.
func (c runToggle) FollowsItsGraph() conform.ID[Run] {
	return c.Check("follows-its-graph")
}

// runStatus is the name of Run.Status, with the identities of its checks.
type runStatus struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Status returns without a
// panic for derived inputs.
func (c runStatus) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Status returns
// without a panic for derived inputs.
func (c runStatus) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Status returns an error
// for a cancelled context.
func (c runStatus) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Status returns an
// error for a cancelled context.
func (c runStatus) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// RunSuite returns the checks of Run, and the calls and the budgets of its
// callables.
func RunSuite() conform.Suite[Run] {
	return conform.Suite[Run]{
		Declaration: "Run",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Run]{
			{
				ID:          RunChecks.Advance.Smoke(),
				Statement:   "Advance returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, string("key-a"))
					})
				},
				Defect: runAdvancePanics,
			},
			{
				ID:          RunChecks.Advance.SmokeGenerated(),
				Statement:   "Advance returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, drawnKey)
					})
				},
				Defect: runAdvancePanics,
			},
			{
				ID:          RunChecks.Advance.Cancel(),
				Statement:   "Advance returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, string("key-a"))
					})
				},
				Defect: runAdvanceIgnoresCancel,
			},
			{
				ID:          RunChecks.Advance.CancelGenerated(),
				Statement:   "Advance returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, drawnKey)
					})
				},
				Defect: runAdvanceIgnoresCancel,
			},
			{
				ID:          RunChecks.Advance.FollowsItsGraph(),
				Statement:   "each call of Advance leaves the state that Status reads at the state before it or at one of its successors, over any sequence of calls",
				Claim:       spec.FollowsItsGraph,
				Probe:       spec.StatesMoveAlongTheGraph,
				Reason:      spec.ForbiddenTransition,
				Law:         spec.StatesFollowTheGraph,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *RunHandle, c *prop.Case) {
					probe.StatesMoveAlongTheGraph(
						h,
						c,
						runActions,
						"Advance",
						func(ctx context.Context, args []any) conform.Output {
							values, err := runStatusCall(
								ctx,
								h.Subject(),
								[]any{args[0], bool(true)},
							)
							return conform.OutputOf(err, values...)
						},
						Graph,
					)
				},
				Defect: runAdvanceAdmitsAnyTransitionAfterThree,
			},
			{
				ID:          RunChecks.Retry.Smoke(),
				Statement:   "Retry returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, string("key-a"))
					})
				},
				Defect: runRetryPanics,
			},
			{
				ID:          RunChecks.Retry.SmokeGenerated(),
				Statement:   "Retry returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, drawnKey)
					})
				},
				Defect: runRetryPanics,
			},
			{
				ID:          RunChecks.Retry.Cancel(),
				Statement:   "Retry returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, string("key-a"))
					})
				},
				Defect: runRetryIgnoresCancel,
			},
			{
				ID:          RunChecks.Retry.CancelGenerated(),
				Statement:   "Retry returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, drawnKey)
					})
				},
				Defect: runRetryIgnoresCancel,
			},
			{
				ID:          RunChecks.Retry.FollowsItsGraph(),
				Statement:   "each call of Retry leaves the state that Status reads at the state before it or at one of its successors, over any sequence of calls",
				Claim:       spec.FollowsItsGraph,
				Probe:       spec.StatesMoveAlongTheGraph,
				Reason:      spec.ForbiddenTransition,
				Law:         spec.StatesFollowTheGraph,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *RunHandle, c *prop.Case) {
					probe.StatesMoveAlongTheGraph(
						h,
						c,
						runActions,
						"Retry",
						func(ctx context.Context, args []any) conform.Output {
							values, err := runStatusCall(
								ctx,
								h.Subject(),
								[]any{args[0], bool(true)},
							)
							return conform.OutputOf(err, values...)
						},
						Retries,
					)
				},
				Defect: runRetryAdmitsAnyTransitionAfterThree,
			},
			{
				ID:          RunChecks.Restart.Smoke(),
				Statement:   "Restart returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, string("key-a"))
					})
				},
				Defect: runRestartPanics,
			},
			{
				ID:          RunChecks.Restart.SmokeGenerated(),
				Statement:   "Restart returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, drawnKey)
					})
				},
				Defect: runRestartPanics,
			},
			{
				ID:          RunChecks.Restart.Cancel(),
				Statement:   "Restart returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, string("key-a"))
					})
				},
				Defect: runRestartIgnoresCancel,
			},
			{
				ID:          RunChecks.Restart.CancelGenerated(),
				Statement:   "Restart returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, drawnKey)
					})
				},
				Defect: runRestartIgnoresCancel,
			},
			{
				ID:          RunChecks.Toggle.Smoke(),
				Statement:   "Toggle returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Toggle(ctx, string("key-a"), func() bool { return false })
					})
				},
				Defect: runTogglePanics,
			},
			{
				ID:          RunChecks.Toggle.SmokeGenerated(),
				Statement:   "Toggle returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
					drawnVerbose := func() func() bool {
						var verboseAt atomic.Int64
						return func() (r0 bool) {
							i := int(verboseAt.Add(1) - 1)
							if i >= len(verboseResults) {
								return r0
							}
							return verboseResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Toggle(ctx, drawnKey, drawnVerbose())
					})
				},
				Defect: runTogglePanics,
			},
			{
				ID:          RunChecks.Toggle.Cancel(),
				Statement:   "Toggle returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Toggle(ctx, string("key-a"), func() bool { return false })
					})
				},
				Defect: runToggleIgnoresCancel,
			},
			{
				ID:          RunChecks.Toggle.CancelGenerated(),
				Statement:   "Toggle returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
					drawnVerbose := func() func() bool {
						var verboseAt atomic.Int64
						return func() (r0 bool) {
							i := int(verboseAt.Add(1) - 1)
							if i >= len(verboseResults) {
								return r0
							}
							return verboseResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Toggle(ctx, drawnKey, drawnVerbose())
					})
				},
				Defect: runToggleIgnoresCancel,
			},
			{
				ID:          RunChecks.Toggle.FollowsItsGraph(),
				Statement:   "each call of Toggle leaves the state that Status reads at the state before it or at one of its successors, over any sequence of calls",
				Claim:       spec.FollowsItsGraph,
				Probe:       spec.StatesMoveAlongTheGraph,
				Reason:      spec.ForbiddenTransition,
				Law:         spec.StatesFollowTheGraph,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *RunHandle, c *prop.Case) {
					probe.StatesMoveAlongTheGraph(
						h,
						c,
						runActions,
						"Toggle",
						func(ctx context.Context, args []any) conform.Output {
							values, err := runStatusCall(
								ctx,
								h.Subject(),
								[]any{args[0], bool(true)},
							)
							return conform.OutputOf(err, values...)
						},
						Graph,
					)
				},
				Defect: runToggleAdmitsAnyTransitionAfterThree,
			},
			{
				ID:          RunChecks.Status.Smoke(),
				Statement:   "Status returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, string("key-a"), bool(true))
						return err
					})
				},
				Defect: runStatusPanics,
			},
			{
				ID:          RunChecks.Status.SmokeGenerated(),
				Statement:   "Status returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnVerbose := c.Draw(prop.Of[bool](), "verbose")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, drawnKey, drawnVerbose)
						return err
					})
				},
				Defect: runStatusPanics,
			},
			{
				ID:          RunChecks.Status.Cancel(),
				Statement:   "Status returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, string("key-a"), bool(true))
						return err
					})
				},
				Defect: runStatusIgnoresCancel,
			},
			{
				ID:          RunChecks.Status.CancelGenerated(),
				Statement:   "Status returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnVerbose := c.Draw(prop.Of[bool](), "verbose")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, drawnKey, drawnVerbose)
						return err
					})
				},
				Defect: runStatusIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.FollowsItsGraph,
				Callable: "Run.Restart",
				Gate:     spec.GraphOfStates,
				Why:      "the variable that transitions names for Restart is no graph of the states that observe reads",
				Remedy:   "declare transitions as a map from each state to the list of its successors",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Run.Advance",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Run.Retry",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Run.Restart",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Run.Toggle",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Run]{
			{
				Calls:    runAdvanceCalls,
				Over:     newRunOver,
				Callable: RunChecks.Advance.Callable,
			},
			{
				Calls:    runRetryCalls,
				Over:     newRunOver,
				Callable: RunChecks.Retry.Callable,
			},
			{
				Calls:    runRestartCalls,
				Over:     newRunOver,
				Callable: RunChecks.Restart.Callable,
			},
			{
				Calls:    runToggleCalls,
				Over:     newRunOver,
				Callable: RunChecks.Toggle.Callable,
			},
			{
				Calls:    runStatusCalls,
				Over:     newRunOver,
				Callable: RunChecks.Status.Callable,
			},
		},
	}
}

// RunRunSuite applies the checks of Run, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunRunSuite[T Run](t *testing.T, newSubject func() T, opts ...RunOption) {
	t.Helper()
	RunRunSuiteFromHandle(t, func(*RunHandle) T { return newSubject() }, opts...)
}

// RunRunSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRunSuiteFromHandle[T Run](
	t *testing.T,
	newSubject func(h *RunHandle) T,
	opts ...RunOption,
) {
	t.Helper()
	conform.Run(t, RunSuite(), runSubjects(newSubject), opts...)
}

// FuzzRunSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRunSuite[T Run](
	f *testing.F,
	callable conform.Callable[Run],
	newSubject func(h *RunHandle) T,
	opts ...RunOption,
) {
	f.Helper()
	conform.Fuzz(f, RunSuite(), callable, runSubjects(newSubject), opts...)
}

// RunRunBudgets applies the allocation ceilings of Run and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunRunBudgets[T Run](t *testing.T, newSubject func(h *RunHandle) T, opts ...RunOption) {
	t.Helper()
	conform.RunBudgets(t, RunSuite(), runSubjects(newSubject), opts...)
}

// BenchRunSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchRunSuite[T Run](
	b conform.B,
	callable conform.Callable[Run],
	newSubject func(h *RunHandle) T,
	opts ...RunOption,
) {
	b.Helper()
	conform.Bench(b, RunSuite(), callable, runSubjects(newSubject), opts...)
}

// runSubjects returns the factory of the subjects of Run that calls
// newSubject with the handle of each check.
func runSubjects[T Run](newSubject func(h *RunHandle) T) func(h *RunHandle) Run {
	return func(h *RunHandle) Run { return newSubject(h) }
}

// runAdvanceCalls generates the calls of Run.Advance with the inputs that
// each case draws.
var runAdvanceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Run) error {
		return subject.Advance(ctx, key)
	}
})

// runRetryCalls generates the calls of Run.Retry with the inputs that each
// case draws.
var runRetryCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Run) error {
		return subject.Retry(ctx, key)
	}
})

// runRestartCalls generates the calls of Run.Restart with the inputs that
// each case draws.
var runRestartCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Run) error {
		return subject.Restart(ctx, key)
	}
})

// runToggleCalls generates the calls of Run.Toggle with the inputs that
// each case draws.
var runToggleCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
	var verboseAt atomic.Int64
	verbose := func() (r0 bool) {
		i := int(verboseAt.Add(1) - 1)
		if i >= len(verboseResults) {
			return r0
		}
		return verboseResults[i]
	}
	return func(ctx context.Context, subject Run) error {
		return subject.Toggle(ctx, key, verbose)
	}
})

// runStatusCalls generates the calls of Run.Status with the inputs that
// each case draws.
var runStatusCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	verbose := c.Draw(prop.Of[bool](), "verbose")
	return func(ctx context.Context, subject Run) error {
		_, err := subject.Status(ctx, key, verbose)
		return err
	}
})

// runOverEmbedded is Run under the name of the field that runOver embeds.
type runOverEmbedded = Run

// runOver is a subject of Run whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type runOver struct {
	runOverEmbedded
	spend func()
}

// Advance calls spend and returns its zero results.
func (o runOver) Advance(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Retry calls spend and returns its zero results.
func (o runOver) Retry(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Restart calls spend and returns its zero results.
func (o runOver) Restart(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Toggle calls spend and returns its zero results.
func (o runOver) Toggle(context.Context, string, func() bool) (r0 error) {
	o.spend()
	return r0
}

// Status calls spend and returns its zero results.
func (o runOver) Status(context.Context, string, bool) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newRunOver returns the planted defect of a ceiling of Run, whose
// callables call spend.
func newRunOver(_ *RunHandle, spend func()) Run {
	return runOver{spend: spend}
}

// runAdvancePanicsDefect is the planted defect of a Run whose Advance
// panics.
type runAdvancePanicsDefect struct {
	runOverEmbedded
}

// Advance is the override of the planted defect.
func (*runAdvancePanicsDefect) Advance(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Run whose Advance panics")
}

// runAdvancePanics returns the planted defect of a Run whose Advance
// panics over a fresh subject.
func runAdvancePanics(h *RunHandle) Run {
	return &runAdvancePanicsDefect{runOverEmbedded: h.Fresh()}
}

// runAdvanceIgnoresCancelDefect is the planted defect of a Run whose
// Advance ignores the context.
type runAdvanceIgnoresCancelDefect struct {
	runOverEmbedded
}

// Advance is the override of the planted defect.
func (d *runAdvanceIgnoresCancelDefect) Advance(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Advance(p0, p1)
}

// runAdvanceIgnoresCancel returns the planted defect of a Run whose
// Advance ignores the context over a fresh subject.
func runAdvanceIgnoresCancel(h *RunHandle) Run {
	return &runAdvanceIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runAdvanceAdmitsAnyTransitionAfterThreeDefect is the planted defect of a
// Run whose Advance makes its transition twice from its fourth call on.
type runAdvanceAdmitsAnyTransitionAfterThreeDefect struct {
	runOverEmbedded
	calls int
}

// Advance is the override of the planted defect.
func (d *runAdvanceAdmitsAnyTransitionAfterThreeDefect) Advance(
	p0 context.Context,
	p1 string,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		r0 = d.runOverEmbedded.Advance(p0, p1)
	}
	return d.runOverEmbedded.Advance(p0, p1)
}

// runAdvanceAdmitsAnyTransitionAfterThree returns the planted defect of a
// Run whose Advance makes its transition twice from its fourth call on
// over a fresh subject.
func runAdvanceAdmitsAnyTransitionAfterThree(h *RunHandle) Run {
	return &runAdvanceAdmitsAnyTransitionAfterThreeDefect{runOverEmbedded: h.Fresh()}
}

// runRetryPanicsDefect is the planted defect of a Run whose Retry panics.
type runRetryPanicsDefect struct {
	runOverEmbedded
}

// Retry is the override of the planted defect.
func (*runRetryPanicsDefect) Retry(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Run whose Retry panics")
}

// runRetryPanics returns the planted defect of a Run whose Retry panics
// over a fresh subject.
func runRetryPanics(h *RunHandle) Run {
	return &runRetryPanicsDefect{runOverEmbedded: h.Fresh()}
}

// runRetryIgnoresCancelDefect is the planted defect of a Run whose Retry
// ignores the context.
type runRetryIgnoresCancelDefect struct {
	runOverEmbedded
}

// Retry is the override of the planted defect.
func (d *runRetryIgnoresCancelDefect) Retry(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Retry(p0, p1)
}

// runRetryIgnoresCancel returns the planted defect of a Run whose Retry
// ignores the context over a fresh subject.
func runRetryIgnoresCancel(h *RunHandle) Run {
	return &runRetryIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runRetryAdmitsAnyTransitionAfterThreeDefect is the planted defect of a
// Run whose Retry makes its transition twice from its fourth call on.
type runRetryAdmitsAnyTransitionAfterThreeDefect struct {
	runOverEmbedded
	calls int
}

// Retry is the override of the planted defect.
func (d *runRetryAdmitsAnyTransitionAfterThreeDefect) Retry(
	p0 context.Context,
	p1 string,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		r0 = d.runOverEmbedded.Retry(p0, p1)
	}
	return d.runOverEmbedded.Retry(p0, p1)
}

// runRetryAdmitsAnyTransitionAfterThree returns the planted defect of a
// Run whose Retry makes its transition twice from its fourth call on over
// a fresh subject.
func runRetryAdmitsAnyTransitionAfterThree(h *RunHandle) Run {
	return &runRetryAdmitsAnyTransitionAfterThreeDefect{runOverEmbedded: h.Fresh()}
}

// runRestartPanicsDefect is the planted defect of a Run whose Restart
// panics.
type runRestartPanicsDefect struct {
	runOverEmbedded
}

// Restart is the override of the planted defect.
func (*runRestartPanicsDefect) Restart(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Run whose Restart panics")
}

// runRestartPanics returns the planted defect of a Run whose Restart
// panics over a fresh subject.
func runRestartPanics(h *RunHandle) Run {
	return &runRestartPanicsDefect{runOverEmbedded: h.Fresh()}
}

// runRestartIgnoresCancelDefect is the planted defect of a Run whose
// Restart ignores the context.
type runRestartIgnoresCancelDefect struct {
	runOverEmbedded
}

// Restart is the override of the planted defect.
func (d *runRestartIgnoresCancelDefect) Restart(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Restart(p0, p1)
}

// runRestartIgnoresCancel returns the planted defect of a Run whose
// Restart ignores the context over a fresh subject.
func runRestartIgnoresCancel(h *RunHandle) Run {
	return &runRestartIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runTogglePanicsDefect is the planted defect of a Run whose Toggle
// panics.
type runTogglePanicsDefect struct {
	runOverEmbedded
}

// Toggle is the override of the planted defect.
func (*runTogglePanicsDefect) Toggle(p0 context.Context, p1 string, p2 func() bool) (r0 error) {
	panic("the planted defect of a Run whose Toggle panics")
}

// runTogglePanics returns the planted defect of a Run whose Toggle panics
// over a fresh subject.
func runTogglePanics(h *RunHandle) Run {
	return &runTogglePanicsDefect{runOverEmbedded: h.Fresh()}
}

// runToggleIgnoresCancelDefect is the planted defect of a Run whose Toggle
// ignores the context.
type runToggleIgnoresCancelDefect struct {
	runOverEmbedded
}

// Toggle is the override of the planted defect.
func (d *runToggleIgnoresCancelDefect) Toggle(
	p0 context.Context,
	p1 string,
	p2 func() bool,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Toggle(p0, p1, p2)
}

// runToggleIgnoresCancel returns the planted defect of a Run whose Toggle
// ignores the context over a fresh subject.
func runToggleIgnoresCancel(h *RunHandle) Run {
	return &runToggleIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runToggleAdmitsAnyTransitionAfterThreeDefect is the planted defect of a
// Run whose Toggle makes its transition twice from its fourth call on.
type runToggleAdmitsAnyTransitionAfterThreeDefect struct {
	runOverEmbedded
	calls int
}

// Toggle is the override of the planted defect.
func (d *runToggleAdmitsAnyTransitionAfterThreeDefect) Toggle(
	p0 context.Context,
	p1 string,
	p2 func() bool,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		r0 = d.runOverEmbedded.Toggle(p0, p1, p2)
	}
	return d.runOverEmbedded.Toggle(p0, p1, p2)
}

// runToggleAdmitsAnyTransitionAfterThree returns the planted defect of a
// Run whose Toggle makes its transition twice from its fourth call on over
// a fresh subject.
func runToggleAdmitsAnyTransitionAfterThree(h *RunHandle) Run {
	return &runToggleAdmitsAnyTransitionAfterThreeDefect{runOverEmbedded: h.Fresh()}
}

// runStatusPanicsDefect is the planted defect of a Run whose Status
// panics.
type runStatusPanicsDefect struct {
	runOverEmbedded
}

// Status is the override of the planted defect.
func (*runStatusPanicsDefect) Status(p0 context.Context, p1 string, p2 bool) (r0 string, r1 error) {
	panic("the planted defect of a Run whose Status panics")
}

// runStatusPanics returns the planted defect of a Run whose Status panics
// over a fresh subject.
func runStatusPanics(h *RunHandle) Run {
	return &runStatusPanicsDefect{runOverEmbedded: h.Fresh()}
}

// runStatusIgnoresCancelDefect is the planted defect of a Run whose Status
// ignores the context.
type runStatusIgnoresCancelDefect struct {
	runOverEmbedded
}

// Status is the override of the planted defect.
func (d *runStatusIgnoresCancelDefect) Status(
	p0 context.Context,
	p1 string,
	p2 bool,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.runOverEmbedded.Status(p0, p1, p2)
}

// runStatusIgnoresCancel returns the planted defect of a Run whose Status
// ignores the context over a fresh subject.
func runStatusIgnoresCancel(h *RunHandle) Run {
	return &runStatusIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runActions returns the actions of the generated sequences of Run over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func runActions(c *prop.Case, subject func() Run) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Advance",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "advance.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Advance", args)
				values, err := runAdvanceCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Retry",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "retry.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Retry", args)
				values, err := runRetryCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Restart",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "restart.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Restart", args)
				values, err := runRestartCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Toggle",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "toggle.key")
				verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
				var verboseAt atomic.Int64
				verbose := func() (r0 bool) {
					i := int(verboseAt.Add(1) - 1)
					if i >= len(verboseResults) {
						return r0
					}
					return verboseResults[i]
				}
				return []any{key, verbose}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Toggle", args, args[0])
				values, err := runToggleCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Status",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "status.key")
				verbose := c.Draw(prop.Of[bool](), "status.verbose")
				return []any{key, verbose}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Status", args, args[0])
				values, err := runStatusCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// runAdvanceCall calls Advance on subject with the arguments args, and
// returns the values of the call beside its error.
func runAdvanceCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Advance(ctx, a0)
}

// runRetryCall calls Retry on subject with the arguments args, and returns
// the values of the call beside its error.
func runRetryCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Retry(ctx, a0)
}

// runRestartCall calls Restart on subject with the arguments args, and
// returns the values of the call beside its error.
func runRestartCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Restart(ctx, a0)
}

// runToggleCall calls Toggle on subject with the arguments args, and
// returns the values of the call beside its error.
func runToggleCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(func() bool)
	return nil, subject.Toggle(ctx, a0, a1)
}

// runStatusCall calls Status on subject with the arguments args, and
// returns the values of the call beside its error.
func runStatusCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(bool)
	r0, err := subject.Status(ctx, a0, a1)
	return []any{r0}, err
}

// The derived checks of Log:
//
// Log.Append/smoke: Append returns without a panic for derived inputs.
//
// Log.Append/smoke#generated: Append returns without a panic for derived
// inputs.
//
// Log.Append/cancel: Append returns an error for a cancelled context.
//
// Log.Append/cancel#generated: Append returns an error for a cancelled
// context.
//
// Log.Append/zero-on-error: Append returns the zero value beside an error.
//
// Log.Append/zero-on-error#generated: Append returns the zero value beside
// an error.
//
// Log.Append/appender: Append returns a position above the position of the
// append before it.
//
// Log.Append/appender#generated: Append returns a position above the
// position of the append before it.
//
// Log.Append/appends-keep-ascending: Append returns a position above the
// one before it, over any sequence of calls.

// The derived checks of Log that witness their laws at the shallow end:
//
// Log.Append/smoke
//
// Log.Append/appender
//
// Log.Append/appender#generated

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Log writes.
type (
	LogHandle = conform.Handle[Log]
	LogCheck  = conform.Check[Log]
	LogOption = conform.Option[Log]
)

// LogChecks contains the name of each callable of Log, with the identities
// of its checks.
var LogChecks = struct {
	Append logAppend
}{
	Append: logAppend{"Log.Append"},
}

// logAppend is the name of Log.Append, with the identities of its checks.
type logAppend struct {
	conform.Callable[Log]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c logAppend) Smoke() conform.ID[Log] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Append returns
// without a panic for derived inputs.
func (c logAppend) SmokeGenerated() conform.ID[Log] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c logAppend) Cancel() conform.ID[Log] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Append returns an
// error for a cancelled context.
func (c logAppend) CancelGenerated() conform.ID[Log] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Append returns the
// zero value beside an error.
func (c logAppend) ZeroOnError() conform.ID[Log] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Append
// returns the zero value beside an error.
func (c logAppend) ZeroOnErrorGenerated() conform.ID[Log] {
	return c.Check("zero-on-error#generated")
}

// Appender returns the identity of the check that Append returns a
// position above the position of the append before it.
func (c logAppend) Appender() conform.ID[Log] {
	return c.Check("appender")
}

// AppenderGenerated returns the identity of the check that Append returns
// a position above the position of the append before it.
func (c logAppend) AppenderGenerated() conform.ID[Log] {
	return c.Check("appender#generated")
}

// AppendsKeepAscending returns the identity of the check that Append
// returns a position above the one before it, over any sequence of calls.
func (c logAppend) AppendsKeepAscending() conform.ID[Log] {
	return c.Check("appends-keep-ascending")
}

// LogSuite returns the checks of Log, and the calls and the budgets of its
// callables.
func LogSuite() conform.Suite[Log] {
	return conform.Suite[Log]{
		Declaration: "Log",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Log]{
			{
				ID:          LogChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LogHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, string("entry-a"))
						return err
					})
				},
				Defect: logAppendPanics,
			},
			{
				ID:          LogChecks.Append.SmokeGenerated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, drawnEntry)
						return err
					})
				},
				Defect: logAppendPanics,
			},
			{
				ID:          LogChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, string("entry-a"))
						return err
					})
				},
				Defect: logAppendIgnoresCancel,
			},
			{
				ID:          LogChecks.Append.CancelGenerated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, drawnEntry)
						return err
					})
				},
				Defect: logAppendIgnoresCancel,
			},
			{
				ID:          LogChecks.Append.ZeroOnError(),
				Statement:   "Append returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, string("entry-a"))
					}, int64(0))
				},
				Defect: logAppendValueBeside,
			},
			{
				ID:          LogChecks.Append.ZeroOnErrorGenerated(),
				Statement:   "Append returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.ZeroBesideError(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, drawnEntry)
					}, int64(0))
				},
				Defect: logAppendValueBeside,
			},
			{
				ID:          LogChecks.Append.Appender(),
				Statement:   "Append returns a position above the position of the append before it",
				Claim:       spec.Appender,
				Probe:       spec.RepeatAscends,
				Reason:      spec.NotAscending,
				Law:         spec.AppendReturnsAboveTheLast,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LogHandle) {
					probe.RepeatAscends(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, string("entry-a"))
					})
				},
				Defect: logAppendDescends,
			},
			{
				ID:          LogChecks.Append.AppenderGenerated(),
				Statement:   "Append returns a position above the position of the append before it",
				Claim:       spec.Appender,
				Probe:       spec.RepeatAscends,
				Reason:      spec.NotAscending,
				Law:         spec.AppendReturnsAboveTheLast,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.RepeatAscends(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, drawnEntry)
					})
				},
				Defect: logAppendDescends,
			},
			{
				ID:          LogChecks.Append.AppendsKeepAscending(),
				Statement:   "Append returns a position above the one before it, over any sequence of calls",
				Claim:       spec.AppendsKeepAscending,
				Probe:       spec.PositionsRise,
				Reason:      spec.NotAscending,
				Law:         spec.AppendReturnsAboveTheLast,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Log]{LogChecks.Append.Appender()},
				Property: func(h *LogHandle, c *prop.Case) {
					probe.PositionsRise[Log, int64](h, c, logActions, "Append")
				},
				Defect: logAppendDescendsAfterThree,
			},
		},
		Calls: []conform.Call[Log]{
			{
				Calls:    logAppendCalls,
				Over:     newLogOver,
				Callable: LogChecks.Append.Callable,
			},
		},
	}
}

// RunLogSuite applies the checks of Log, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunLogSuite[T Log](t *testing.T, newSubject func() T, opts ...LogOption) {
	t.Helper()
	RunLogSuiteFromHandle(t, func(*LogHandle) T { return newSubject() }, opts...)
}

// RunLogSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLogSuiteFromHandle[T Log](
	t *testing.T,
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	t.Helper()
	conform.Run(t, LogSuite(), logSubjects(newSubject), opts...)
}

// FuzzLogSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLogSuite[T Log](
	f *testing.F,
	callable conform.Callable[Log],
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	f.Helper()
	conform.Fuzz(f, LogSuite(), callable, logSubjects(newSubject), opts...)
}

// RunLogBudgets applies the allocation ceilings of Log and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunLogBudgets[T Log](t *testing.T, newSubject func(h *LogHandle) T, opts ...LogOption) {
	t.Helper()
	conform.RunBudgets(t, LogSuite(), logSubjects(newSubject), opts...)
}

// BenchLogSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLogSuite[T Log](
	b conform.B,
	callable conform.Callable[Log],
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	b.Helper()
	conform.Bench(b, LogSuite(), callable, logSubjects(newSubject), opts...)
}

// logSubjects returns the factory of the subjects of Log that calls
// newSubject with the handle of each check.
func logSubjects[T Log](newSubject func(h *LogHandle) T) func(h *LogHandle) Log {
	return func(h *LogHandle) Log { return newSubject(h) }
}

// logAppendCalls generates the calls of Log.Append with the inputs that
// each case draws.
var logAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Log) error {
	entry := c.Draw(prop.Of[string](), "entry")
	return func(ctx context.Context, subject Log) error {
		_, err := subject.Append(ctx, entry)
		return err
	}
})

// logOverEmbedded is Log under the name of the field that logOver embeds.
type logOverEmbedded = Log

// logOver is a subject of Log whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type logOver struct {
	logOverEmbedded
	spend func()
}

// Append calls spend and returns its zero results.
func (o logOver) Append(context.Context, string) (r0 int64, r1 error) {
	o.spend()
	return r0, r1
}

// newLogOver returns the planted defect of a ceiling of Log, whose
// callables call spend.
func newLogOver(_ *LogHandle, spend func()) Log {
	return logOver{spend: spend}
}

// logAppendPanicsDefect is the planted defect of a Log whose Append
// panics.
type logAppendPanicsDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (*logAppendPanicsDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	panic("the planted defect of a Log whose Append panics")
}

// logAppendPanics returns the planted defect of a Log whose Append panics
// over a fresh subject.
func logAppendPanics(h *LogHandle) Log {
	return &logAppendPanicsDefect{logOverEmbedded: h.Fresh()}
}

// logAppendIgnoresCancelDefect is the planted defect of a Log whose Append
// ignores the context.
type logAppendIgnoresCancelDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (d *logAppendIgnoresCancelDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.logOverEmbedded.Append(p0, p1)
}

// logAppendIgnoresCancel returns the planted defect of a Log whose Append
// ignores the context over a fresh subject.
func logAppendIgnoresCancel(h *LogHandle) Log {
	return &logAppendIgnoresCancelDefect{logOverEmbedded: h.Fresh()}
}

// logAppendValueBesideDefect is the planted defect of a Log whose Append
// returns a value beside an error.
type logAppendValueBesideDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (d *logAppendValueBesideDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	r0, r1 = d.logOverEmbedded.Append(p0, p1)
	if r1 != nil {
		r0 = int64(7)
	}
	return r0, r1
}

// logAppendValueBeside returns the planted defect of a Log whose Append
// returns a value beside an error over a fresh subject.
func logAppendValueBeside(h *LogHandle) Log {
	return &logAppendValueBesideDefect{logOverEmbedded: h.Fresh()}
}

// logAppendDescendsDefect is the planted defect of a Log whose Append
// returns a position that is not above the last.
type logAppendDescendsDefect struct {
	logOverEmbedded
	calls int
}

// Append is the override of the planted defect.
func (d *logAppendDescendsDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls == 1 {
		r0 = int64(7)
	}
	return r0, r1
}

// logAppendDescends returns the planted defect of a Log whose Append
// returns a position that is not above the last over a fresh subject.
func logAppendDescends(h *LogHandle) Log {
	return &logAppendDescendsDefect{logOverEmbedded: h.Fresh()}
}

// logAppendDescendsAfterThreeDefect is the planted defect of a Log whose
// Append returns the position 0 from its fourth call on.
type logAppendDescendsAfterThreeDefect struct {
	logOverEmbedded
	calls int
}

// Append is the override of the planted defect.
func (d *logAppendDescendsAfterThreeDefect) Append(
	p0 context.Context,
	p1 string,
) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls > 3 {
		return r0, r1
	}
	return d.logOverEmbedded.Append(p0, p1)
}

// logAppendDescendsAfterThree returns the planted defect of a Log whose
// Append returns the position 0 from its fourth call on over a fresh
// subject.
func logAppendDescendsAfterThree(h *LogHandle) Log {
	return &logAppendDescendsAfterThreeDefect{logOverEmbedded: h.Fresh()}
}

// logActions returns the actions of the generated sequences of Log over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func logActions(c *prop.Case, subject func() Log) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Append",
			Input: func(c *prop.Case, _ struct{}) any {
				entry := c.Draw(prop.SampledFrom(keysString...), "append.entry")
				return []any{entry}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Append", args, args[0])
				values, err := logAppendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// logAppendCall calls Append on subject with the arguments args, and
// returns the values of the call beside its error.
func logAppendCall(ctx context.Context, subject Log, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Append(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Feed:
//
// Feed.Put/smoke: Put returns without a panic for derived inputs.
//
// Feed.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Feed.Put/cancel: Put returns an error for a cancelled context.
//
// Feed.Put/cancel#generated: Put returns an error for a cancelled context.
//
// Feed.All/smoke: All returns without a panic for derived inputs.
//
// Feed.All/no-duplicates: no element appears twice in one All.
//
// Feed.All/no-duplicates#generated: no element appears twice in one All.
//
// Feed.All/streams-what-was-stored: All yields each value that Put stored.
//
// Feed.All/streams-what-was-stored#generated: All yields each value that
// Put stored.
//
// Feed.All/yields-each-once: All never yields an element twice, over any
// sequence of calls.
//
// Feed.Each/smoke: Each returns without a panic for derived inputs.
//
// Feed.Each/no-duplicates: no element appears twice in one Each.
//
// Feed.Each/no-duplicates#generated: no element appears twice in one Each.
//
// Feed.Each/streams-what-was-stored: Each yields each value that Put
// stored.
//
// Feed.Each/streams-what-was-stored#generated: Each yields each value that
// Put stored.
//
// Feed.Each/yields-each-once: Each never yields an element twice, over any
// sequence of calls.
//
// Feed.Pairs/smoke: Pairs returns without a panic for derived inputs.
//
// Feed.Pairs/no-duplicates: no element appears twice in one Pairs.
//
// Feed.Pairs/no-duplicates#generated: no element appears twice in one
// Pairs.
//
// Feed.Pairs/streams-what-was-stored: Pairs yields each value that Put
// stored.
//
// Feed.Pairs/streams-what-was-stored#generated: Pairs yields each value
// that Put stored.
//
// Feed.Pairs/yields-each-once: Pairs never yields an element twice, over
// any sequence of calls.

// The derived checks of Feed that witness their laws at the shallow end:
//
// Feed.Put/smoke
//
// Feed.All/smoke
//
// Feed.All/no-duplicates
//
// Feed.All/no-duplicates#generated
//
// Feed.All/streams-what-was-stored
//
// Feed.All/streams-what-was-stored#generated
//
// Feed.Each/smoke
//
// Feed.Each/no-duplicates
//
// Feed.Each/no-duplicates#generated
//
// Feed.Each/streams-what-was-stored
//
// Feed.Each/streams-what-was-stored#generated
//
// Feed.Pairs/smoke
//
// Feed.Pairs/no-duplicates
//
// Feed.Pairs/no-duplicates#generated
//
// Feed.Pairs/streams-what-was-stored
//
// Feed.Pairs/streams-what-was-stored#generated

// The refused claims of Feed:
//
// Feed.All/cancel: All cannot fail. The remedy is to give All an error
// result.
//
// Feed.Each/cancel: Each cannot fail. The remedy is to give Each an error
// result.
//
// Feed.Pairs/cancel: Pairs cannot fail. The remedy is to give Pairs an
// error result.

// The unclaimed classifications of Feed:
//
// The shape:writer of Feed.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Feed writes.
type (
	FeedHandle = conform.Handle[Feed]
	FeedCheck  = conform.Check[Feed]
	FeedOption = conform.Option[Feed]
)

// FeedChecks contains the name of each callable of Feed, with the
// identities of its checks.
var FeedChecks = struct {
	Put   feedPut
	All   feedAll
	Each  feedEach
	Pairs feedPairs
}{
	Put:   feedPut{"Feed.Put"},
	All:   feedAll{"Feed.All"},
	Each:  feedEach{"Feed.Each"},
	Pairs: feedPairs{"Feed.Pairs"},
}

// feedPut is the name of Feed.Put, with the identities of its checks.
type feedPut struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c feedPut) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c feedPut) SmokeGenerated() conform.ID[Feed] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c feedPut) Cancel() conform.ID[Feed] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c feedPut) CancelGenerated() conform.ID[Feed] {
	return c.Check("cancel#generated")
}

// feedAll is the name of Feed.All, with the identities of its checks.
type feedAll struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that All returns without a panic
// for derived inputs.
func (c feedAll) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// NoDuplicates returns the identity of the check that no element appears
// twice in one All.
func (c feedAll) NoDuplicates() conform.ID[Feed] {
	return c.Check("no-duplicates")
}

// NoDuplicatesGenerated returns the identity of the check that no element
// appears twice in one All.
func (c feedAll) NoDuplicatesGenerated() conform.ID[Feed] {
	return c.Check("no-duplicates#generated")
}

// StreamsWhatWasStored returns the identity of the check that All yields
// each value that Put stored.
func (c feedAll) StreamsWhatWasStored() conform.ID[Feed] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that All
// yields each value that Put stored.
func (c feedAll) StreamsWhatWasStoredGenerated() conform.ID[Feed] {
	return c.Check("streams-what-was-stored#generated")
}

// YieldsEachOnce returns the identity of the check that All never yields
// an element twice, over any sequence of calls.
func (c feedAll) YieldsEachOnce() conform.ID[Feed] {
	return c.Check("yields-each-once")
}

// feedEach is the name of Feed.Each, with the identities of its checks.
type feedEach struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that Each returns without a
// panic for derived inputs.
func (c feedEach) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// NoDuplicates returns the identity of the check that no element appears
// twice in one Each.
func (c feedEach) NoDuplicates() conform.ID[Feed] {
	return c.Check("no-duplicates")
}

// NoDuplicatesGenerated returns the identity of the check that no element
// appears twice in one Each.
func (c feedEach) NoDuplicatesGenerated() conform.ID[Feed] {
	return c.Check("no-duplicates#generated")
}

// StreamsWhatWasStored returns the identity of the check that Each yields
// each value that Put stored.
func (c feedEach) StreamsWhatWasStored() conform.ID[Feed] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Each yields each value that Put stored.
func (c feedEach) StreamsWhatWasStoredGenerated() conform.ID[Feed] {
	return c.Check("streams-what-was-stored#generated")
}

// YieldsEachOnce returns the identity of the check that Each never yields
// an element twice, over any sequence of calls.
func (c feedEach) YieldsEachOnce() conform.ID[Feed] {
	return c.Check("yields-each-once")
}

// feedPairs is the name of Feed.Pairs, with the identities of its checks.
type feedPairs struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that Pairs returns without a
// panic for derived inputs.
func (c feedPairs) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// NoDuplicates returns the identity of the check that no element appears
// twice in one Pairs.
func (c feedPairs) NoDuplicates() conform.ID[Feed] {
	return c.Check("no-duplicates")
}

// NoDuplicatesGenerated returns the identity of the check that no element
// appears twice in one Pairs.
func (c feedPairs) NoDuplicatesGenerated() conform.ID[Feed] {
	return c.Check("no-duplicates#generated")
}

// StreamsWhatWasStored returns the identity of the check that Pairs yields
// each value that Put stored.
func (c feedPairs) StreamsWhatWasStored() conform.ID[Feed] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Pairs yields each value that Put stored.
func (c feedPairs) StreamsWhatWasStoredGenerated() conform.ID[Feed] {
	return c.Check("streams-what-was-stored#generated")
}

// YieldsEachOnce returns the identity of the check that Pairs never yields
// an element twice, over any sequence of calls.
func (c feedPairs) YieldsEachOnce() conform.ID[Feed] {
	return c.Check("yields-each-once")
}

// FeedSuite returns the checks of Feed, and the calls and the budgets of
// its callables.
func FeedSuite() conform.Suite[Feed] {
	return conform.Suite[Feed]{
		Declaration: "Feed",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Feed]{
			{
				ID:          FeedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("item-a"))
					})
				},
				Defect: feedPutPanics,
			},
			{
				ID:          FeedChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnItem)
					})
				},
				Defect: feedPutPanics,
			},
			{
				ID:          FeedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *FeedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("item-a"))
					})
				},
				Defect: feedPutIgnoresCancel,
			},
			{
				ID:          FeedChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnItem)
					})
				},
				Defect: feedPutIgnoresCancel,
			},
			{
				ID:          FeedChecks.All.Smoke(),
				Statement:   "All returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().All(ctx)
						return nil
					})
				},
				Defect: feedAllPanics,
			},
			{
				ID:          FeedChecks.All.NoDuplicates(),
				Statement:   "no element appears twice in one All",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsTwice,
			},
			{
				ID:          FeedChecks.All.NoDuplicatesGenerated(),
				Statement:   "no element appears twice in one All",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsTwice,
			},
			{
				ID:          FeedChecks.All.StreamsWhatWasStored(),
				Statement:   "All yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsNothing,
			},
			{
				ID:          FeedChecks.All.StreamsWhatWasStoredGenerated(),
				Statement:   "All yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsNothing,
			},
			{
				ID:          FeedChecks.All.YieldsEachOnce(),
				Statement:   "All never yields an element twice, over any sequence of calls",
				Claim:       spec.YieldsEachOnce,
				Probe:       spec.DrainsWithoutRepeats,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Feed]{FeedChecks.All.NoDuplicates()},
				Property: func(h *FeedHandle, c *prop.Case) {
					probe.DrainsWithoutRepeats[Feed, string](h, c, feedActions, "All")
				},
				Defect: feedAllRepeatsAfterThree,
			},
			{
				ID:          FeedChecks.Each.Smoke(),
				Statement:   "Each returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Each(ctx)
						return nil
					})
				},
				Defect: feedEachPanics,
			},
			{
				ID:          FeedChecks.Each.NoDuplicates(),
				Statement:   "no element appears twice in one Each",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsTwice,
			},
			{
				ID:          FeedChecks.Each.NoDuplicatesGenerated(),
				Statement:   "no element appears twice in one Each",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsTwice,
			},
			{
				ID:          FeedChecks.Each.StreamsWhatWasStored(),
				Statement:   "Each yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsNothing,
			},
			{
				ID:          FeedChecks.Each.StreamsWhatWasStoredGenerated(),
				Statement:   "Each yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsNothing,
			},
			{
				ID:          FeedChecks.Each.YieldsEachOnce(),
				Statement:   "Each never yields an element twice, over any sequence of calls",
				Claim:       spec.YieldsEachOnce,
				Probe:       spec.DrainsWithoutRepeats,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Feed]{FeedChecks.Each.NoDuplicates()},
				Property: func(h *FeedHandle, c *prop.Case) {
					probe.DrainsWithoutRepeats[Feed, string](h, c, feedActions, "Each")
				},
				Defect: feedEachRepeatsAfterThree,
			},
			{
				ID:          FeedChecks.Pairs.Smoke(),
				Statement:   "Pairs returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Pairs(ctx)
						return nil
					})
				},
				Defect: feedPairsPanics,
			},
			{
				ID:          FeedChecks.Pairs.NoDuplicates(),
				Statement:   "no element appears twice in one Pairs",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.NoElementTwice(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := string("item-a")
								if calls > 1 {
									w = string("item-b")
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsTwice,
			},
			{
				ID:          FeedChecks.Pairs.NoDuplicatesGenerated(),
				Statement:   "no element appears twice in one Pairs",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.NoElementTwice(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := drawnItem[0]
								if calls > 1 {
									w = drawnItem[1]
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsTwice,
			},
			{
				ID:          FeedChecks.Pairs.StreamsWhatWasStored(),
				Statement:   "Pairs yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := string("item-a")
								if calls > 1 {
									w = string("item-b")
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsNothing,
			},
			{
				ID:          FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
				Statement:   "Pairs yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := drawnItem[0]
								if calls > 1 {
									w = drawnItem[1]
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsNothing,
			},
			{
				ID:          FeedChecks.Pairs.YieldsEachOnce(),
				Statement:   "Pairs never yields an element twice, over any sequence of calls",
				Claim:       spec.YieldsEachOnce,
				Probe:       spec.DrainsWithoutRepeats,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Feed]{FeedChecks.Pairs.NoDuplicates()},
				Property: func(h *FeedHandle, c *prop.Case) {
					probe.DrainsWithoutRepeats[Feed, probe.Pair[string, int]](
						h,
						c,
						feedActions,
						"Pairs",
					)
				},
				Defect: feedPairsRepeatsAfterThree,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Feed.All",
				Gate:     spec.Failure,
				Why:      "All cannot fail",
				Remedy:   "give All an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Feed.Each",
				Gate:     spec.Failure,
				Why:      "Each cannot fail",
				Remedy:   "give Each an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Feed.Pairs",
				Gate:     spec.Failure,
				Why:      "Pairs cannot fail",
				Remedy:   "give Pairs an error result",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Feed.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Feed]{
			{
				Calls:    feedPutCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.Put.Callable,
			},
			{
				Calls:    feedAllCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.All.Callable,
			},
			{
				Calls:    feedEachCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.Each.Callable,
			},
			{
				Calls:    feedPairsCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.Pairs.Callable,
			},
		},
	}
}

// RunFeedSuite applies the checks of Feed, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunFeedSuite[T Feed](t *testing.T, newSubject func() T, opts ...FeedOption) {
	t.Helper()
	RunFeedSuiteFromHandle(t, func(*FeedHandle) T { return newSubject() }, opts...)
}

// RunFeedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunFeedSuiteFromHandle[T Feed](
	t *testing.T,
	newSubject func(h *FeedHandle) T,
	opts ...FeedOption,
) {
	t.Helper()
	conform.Run(t, FeedSuite(), feedSubjects(newSubject), opts...)
}

// FuzzFeedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzFeedSuite[T Feed](
	f *testing.F,
	callable conform.Callable[Feed],
	newSubject func(h *FeedHandle) T,
	opts ...FeedOption,
) {
	f.Helper()
	conform.Fuzz(f, FeedSuite(), callable, feedSubjects(newSubject), opts...)
}

// RunFeedBudgets applies the allocation ceilings of Feed and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunFeedBudgets[T Feed](t *testing.T, newSubject func(h *FeedHandle) T, opts ...FeedOption) {
	t.Helper()
	conform.RunBudgets(t, FeedSuite(), feedSubjects(newSubject), opts...)
}

// BenchFeedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchFeedSuite[T Feed](
	b conform.B,
	callable conform.Callable[Feed],
	newSubject func(h *FeedHandle) T,
	opts ...FeedOption,
) {
	b.Helper()
	conform.Bench(b, FeedSuite(), callable, feedSubjects(newSubject), opts...)
}

// feedSubjects returns the factory of the subjects of Feed that calls
// newSubject with the handle of each check.
func feedSubjects[T Feed](newSubject func(h *FeedHandle) T) func(h *FeedHandle) Feed {
	return func(h *FeedHandle) Feed { return newSubject(h) }
}

// feedPutCalls generates the calls of Feed.Put with the inputs that each
// case draws.
var feedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Feed) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Feed) error {
		return subject.Put(ctx, item)
	}
})

// feedAllCalls generates the calls of Feed.All with the inputs that each
// case draws.
var feedAllCalls = prop.Composite(func(*prop.Case) func(context.Context, Feed) error {
	return func(ctx context.Context, subject Feed) error {
		_ = subject.All(ctx)
		return nil
	}
})

// feedEachCalls generates the calls of Feed.Each with the inputs that each
// case draws.
var feedEachCalls = prop.Composite(func(*prop.Case) func(context.Context, Feed) error {
	return func(ctx context.Context, subject Feed) error {
		_ = subject.Each(ctx)
		return nil
	}
})

// feedPairsCalls generates the calls of Feed.Pairs with the inputs that
// each case draws.
var feedPairsCalls = prop.Composite(func(*prop.Case) func(context.Context, Feed) error {
	return func(ctx context.Context, subject Feed) error {
		_ = subject.Pairs(ctx)
		return nil
	}
})

// feedOverEmbedded is Feed under the name of the field that feedOver
// embeds.
type feedOverEmbedded = Feed

// feedOver is a subject of Feed whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type feedOver struct {
	feedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o feedOver) Put(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// All calls spend and returns its zero results.
func (o feedOver) All(context.Context) (r0 iter.Seq2[string, error]) {
	o.spend()
	return r0
}

// Each calls spend and returns its zero results.
func (o feedOver) Each(context.Context) (r0 iter.Seq[string]) {
	o.spend()
	return r0
}

// Pairs calls spend and returns its zero results.
func (o feedOver) Pairs(context.Context) (r0 iter.Seq2[string, int]) {
	o.spend()
	return r0
}

// newFeedOver returns the planted defect of a ceiling of Feed, whose
// callables call spend.
func newFeedOver(_ *FeedHandle, spend func()) Feed {
	return feedOver{spend: spend}
}

// feedPutPanicsDefect is the planted defect of a Feed whose Put panics.
type feedPutPanicsDefect struct {
	feedOverEmbedded
}

// Put is the override of the planted defect.
func (*feedPutPanicsDefect) Put(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Feed whose Put panics")
}

// feedPutPanics returns the planted defect of a Feed whose Put panics over
// a fresh subject.
func feedPutPanics(h *FeedHandle) Feed {
	return &feedPutPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedPutIgnoresCancelDefect is the planted defect of a Feed whose Put
// ignores the context.
type feedPutIgnoresCancelDefect struct {
	feedOverEmbedded
}

// Put is the override of the planted defect.
func (d *feedPutIgnoresCancelDefect) Put(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.feedOverEmbedded.Put(p0, p1)
}

// feedPutIgnoresCancel returns the planted defect of a Feed whose Put
// ignores the context over a fresh subject.
func feedPutIgnoresCancel(h *FeedHandle) Feed {
	return &feedPutIgnoresCancelDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllPanicsDefect is the planted defect of a Feed whose All panics.
type feedAllPanicsDefect struct {
	feedOverEmbedded
}

// All is the override of the planted defect.
func (*feedAllPanicsDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	panic("the planted defect of a Feed whose All panics")
}

// feedAllPanics returns the planted defect of a Feed whose All panics over
// a fresh subject.
func feedAllPanics(h *FeedHandle) Feed {
	return &feedAllPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllDrainsTwiceDefect is the planted defect of a Feed whose All
// returns one element twice.
type feedAllDrainsTwiceDefect struct {
	feedOverEmbedded
}

// All is the override of the planted defect.
func (d *feedAllDrainsTwiceDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = d.feedOverEmbedded.All(p0)
	if len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedAllDrainsTwice returns the planted defect of a Feed whose All
// returns one element twice over a fresh subject.
func feedAllDrainsTwice(h *FeedHandle) Feed {
	return &feedAllDrainsTwiceDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllDrainsNothingDefect is the planted defect of a Feed whose All
// loses what Put stored.
type feedAllDrainsNothingDefect struct {
	feedOverEmbedded
}

// All is the override of the planted defect.
func (*feedAllDrainsNothingDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = func(func(string, error) bool) {}
	return r0
}

// feedAllDrainsNothing returns the planted defect of a Feed whose All
// loses what Put stored over a fresh subject.
func feedAllDrainsNothing(h *FeedHandle) Feed {
	return &feedAllDrainsNothingDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllRepeatsAfterThreeDefect is the planted defect of a Feed whose All
// yields its first element twice from its fourth call on.
type feedAllRepeatsAfterThreeDefect struct {
	feedOverEmbedded
	calls int
}

// All is the override of the planted defect.
func (d *feedAllRepeatsAfterThreeDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = d.feedOverEmbedded.All(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedAllRepeatsAfterThree returns the planted defect of a Feed whose All
// yields its first element twice from its fourth call on over a fresh
// subject.
func feedAllRepeatsAfterThree(h *FeedHandle) Feed {
	return &feedAllRepeatsAfterThreeDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachPanicsDefect is the planted defect of a Feed whose Each panics.
type feedEachPanicsDefect struct {
	feedOverEmbedded
}

// Each is the override of the planted defect.
func (*feedEachPanicsDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	panic("the planted defect of a Feed whose Each panics")
}

// feedEachPanics returns the planted defect of a Feed whose Each panics
// over a fresh subject.
func feedEachPanics(h *FeedHandle) Feed {
	return &feedEachPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachDrainsTwiceDefect is the planted defect of a Feed whose Each
// returns one element twice.
type feedEachDrainsTwiceDefect struct {
	feedOverEmbedded
}

// Each is the override of the planted defect.
func (d *feedEachDrainsTwiceDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = d.feedOverEmbedded.Each(p0)
	if len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedEachDrainsTwice returns the planted defect of a Feed whose Each
// returns one element twice over a fresh subject.
func feedEachDrainsTwice(h *FeedHandle) Feed {
	return &feedEachDrainsTwiceDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachDrainsNothingDefect is the planted defect of a Feed whose Each
// loses what Put stored.
type feedEachDrainsNothingDefect struct {
	feedOverEmbedded
}

// Each is the override of the planted defect.
func (*feedEachDrainsNothingDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = func(func(string) bool) {}
	return r0
}

// feedEachDrainsNothing returns the planted defect of a Feed whose Each
// loses what Put stored over a fresh subject.
func feedEachDrainsNothing(h *FeedHandle) Feed {
	return &feedEachDrainsNothingDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachRepeatsAfterThreeDefect is the planted defect of a Feed whose
// Each yields its first element twice from its fourth call on.
type feedEachRepeatsAfterThreeDefect struct {
	feedOverEmbedded
	calls int
}

// Each is the override of the planted defect.
func (d *feedEachRepeatsAfterThreeDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = d.feedOverEmbedded.Each(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedEachRepeatsAfterThree returns the planted defect of a Feed whose
// Each yields its first element twice from its fourth call on over a fresh
// subject.
func feedEachRepeatsAfterThree(h *FeedHandle) Feed {
	return &feedEachRepeatsAfterThreeDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsPanicsDefect is the planted defect of a Feed whose Pairs
// panics.
type feedPairsPanicsDefect struct {
	feedOverEmbedded
}

// Pairs is the override of the planted defect.
func (*feedPairsPanicsDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	panic("the planted defect of a Feed whose Pairs panics")
}

// feedPairsPanics returns the planted defect of a Feed whose Pairs panics
// over a fresh subject.
func feedPairsPanics(h *FeedHandle) Feed {
	return &feedPairsPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsDrainsTwiceDefect is the planted defect of a Feed whose Pairs
// returns one element twice.
type feedPairsDrainsTwiceDefect struct {
	feedOverEmbedded
}

// Pairs is the override of the planted defect.
func (d *feedPairsDrainsTwiceDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = d.feedOverEmbedded.Pairs(p0)
	if len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedPairsDrainsTwice returns the planted defect of a Feed whose Pairs
// returns one element twice over a fresh subject.
func feedPairsDrainsTwice(h *FeedHandle) Feed {
	return &feedPairsDrainsTwiceDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsDrainsNothingDefect is the planted defect of a Feed whose Pairs
// loses what Put stored.
type feedPairsDrainsNothingDefect struct {
	feedOverEmbedded
}

// Pairs is the override of the planted defect.
func (*feedPairsDrainsNothingDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = func(func(string, int) bool) {}
	return r0
}

// feedPairsDrainsNothing returns the planted defect of a Feed whose Pairs
// loses what Put stored over a fresh subject.
func feedPairsDrainsNothing(h *FeedHandle) Feed {
	return &feedPairsDrainsNothingDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsRepeatsAfterThreeDefect is the planted defect of a Feed whose
// Pairs yields its first element twice from its fourth call on.
type feedPairsRepeatsAfterThreeDefect struct {
	feedOverEmbedded
	calls int
}

// Pairs is the override of the planted defect.
func (d *feedPairsRepeatsAfterThreeDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = d.feedOverEmbedded.Pairs(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedPairsRepeatsAfterThree returns the planted defect of a Feed whose
// Pairs yields its first element twice from its fourth call on over a
// fresh subject.
func feedPairsRepeatsAfterThree(h *FeedHandle) Feed {
	return &feedPairsRepeatsAfterThreeDefect{feedOverEmbedded: h.Fresh()}
}

// feedActions returns the actions of the generated sequences of Feed over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func feedActions(c *prop.Case, subject func() Feed) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.Of[string](), "put.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Put", args)
				values, err := feedPutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "All",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "All", args)
				values, err := feedAllCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Each",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Each", args)
				values, err := feedEachCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Pairs",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Pairs", args)
				values, err := feedPairsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// feedPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func feedPutCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Put(ctx, a0)
}

// feedAllCall calls All on subject with the arguments args, and returns
// the values of the call beside its error.
func feedAllCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	r0 := subject.All(ctx)
	var drained []string
	for e, err := range r0 {
		if err != nil {
			return []any{drained}, err
		}
		drained = append(drained, e)
	}
	return []any{drained}, nil
}

// feedEachCall calls Each on subject with the arguments args, and returns
// the values of the call beside its error.
func feedEachCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	r0 := subject.Each(ctx)
	return []any{slices.Collect(r0)}, nil
}

// feedPairsCall calls Pairs on subject with the arguments args, and
// returns the values of the call beside its error.
func feedPairsCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	r0 := subject.Pairs(ctx)
	var drained []probe.Pair[string, int]
	for k, v := range r0 {
		drained = append(drained, probe.Pair[string, int]{First: k, Second: v})
	}
	return []any{drained}, nil
}

// The derived checks of Register:
//
// Register.Set/smoke: Set returns without a panic for derived inputs.
//
// Register.Set/smoke#generated: Set returns without a panic for derived
// inputs.
//
// Register.Set/cancel: Set returns an error for a cancelled context.
//
// Register.Set/cancel#generated: Set returns an error for a cancelled
// context.
//
// Register.Get/smoke: Get returns without a panic for derived inputs.
//
// Register.Get/cancel: Get returns an error for a cancelled context.
//
// Register.Get/zero-on-error: Get returns the zero value beside an error.
//
// Register.Get/agrees-with-the-reference: every sequence of calls returns
// on Register what it returns on the reference.

// The derived checks of Register that witness their laws at the shallow
// end:
//
// Register.Set/smoke
//
// Register.Get/smoke

// The unclaimed classifications of Register:
//
// The shape:writer of Register.Set has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Register writes.
type (
	RegisterHandle = conform.Handle[Register]
	RegisterCheck  = conform.Check[Register]
	RegisterOption = conform.Option[Register]
)

// RegisterChecks contains the name of each callable of Register, with the
// identities of its checks.
var RegisterChecks = struct {
	Set registerSet
	Get registerGet
}{
	Set: registerSet{"Register.Set"},
	Get: registerGet{"Register.Get"},
}

// registerSet is the name of Register.Set, with the identities of its
// checks.
type registerSet struct {
	conform.Callable[Register]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c registerSet) Smoke() conform.ID[Register] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Set returns
// without a panic for derived inputs.
func (c registerSet) SmokeGenerated() conform.ID[Register] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Set returns an error for a
// cancelled context.
func (c registerSet) Cancel() conform.ID[Register] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Set returns an
// error for a cancelled context.
func (c registerSet) CancelGenerated() conform.ID[Register] {
	return c.Check("cancel#generated")
}

// registerGet is the name of Register.Get, with the identities of its
// checks.
type registerGet struct {
	conform.Callable[Register]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c registerGet) Smoke() conform.ID[Register] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c registerGet) Cancel() conform.ID[Register] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c registerGet) ZeroOnError() conform.ID[Register] {
	return c.Check("zero-on-error")
}

// AgreesWithTheReference returns the identity of the check that every
// sequence of calls returns on Register what it returns on the reference.
func (c registerGet) AgreesWithTheReference() conform.ID[Register] {
	return "Register/agrees-with-the-reference"
}

// RegisterSuite returns the checks of Register, and the calls and the
// budgets of its callables.
func RegisterSuite() conform.Suite[Register] {
	return conform.Suite[Register]{
		Declaration: "Register",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Register]{
			{
				ID:          RegisterChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegisterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, string("v-a"))
					})
				},
				Defect: registerSetPanics,
			},
			{
				ID:          RegisterChecks.Set.SmokeGenerated(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegisterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[string](), "v")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnV)
					})
				},
				Defect: registerSetPanics,
			},
			{
				ID:          RegisterChecks.Set.Cancel(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegisterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, string("v-a"))
					})
				},
				Defect: registerSetIgnoresCancel,
			},
			{
				ID:          RegisterChecks.Set.CancelGenerated(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegisterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[string](), "v")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnV)
					})
				},
				Defect: registerSetIgnoresCancel,
			},
			{
				ID:          RegisterChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegisterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx)
						return err
					})
				},
				Defect: registerGetPanics,
			},
			{
				ID:          RegisterChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegisterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx)
						return err
					})
				},
				Defect: registerGetIgnoresCancel,
			},
			{
				ID:          RegisterChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegisterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx)
					}, string(""))
				},
				Defect: registerGetValueBeside,
			},
			{
				ID:        RegisterChecks.Get.AgreesWithTheReference(),
				Statement: "every sequence of calls returns on Register what it returns on the reference",
				Claim:     spec.AgreesWithTheReference,
				Probe:     spec.SequencesAgreeWithTheReference,
				Reason:    spec.DivergesFromTheReference,
				Law:       spec.EverySequenceAgreesWithTheReference,
				Machineries: []spec.MachineryID{
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *RegisterHandle, c *prop.Case) {
					probe.SequencesAgreeWithTheReference(h, c, registerActions, registerReference)
				},
				Defect: registerGetDivergesAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Register.Set",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Register]{
			{
				Calls:    registerSetCalls,
				Over:     newRegisterOver,
				Callable: RegisterChecks.Set.Callable,
			},
			{
				Calls:    registerGetCalls,
				Over:     newRegisterOver,
				Callable: RegisterChecks.Get.Callable,
			},
		},
	}
}

// RunRegisterSuite applies the checks of Register, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunRegisterSuite[T Register](t *testing.T, newSubject func() T, opts ...RegisterOption) {
	t.Helper()
	RunRegisterSuiteFromHandle(t, func(*RegisterHandle) T { return newSubject() }, opts...)
}

// RunRegisterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRegisterSuiteFromHandle[T Register](
	t *testing.T,
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	t.Helper()
	conform.Run(t, RegisterSuite(), registerSubjects(newSubject), opts...)
}

// FuzzRegisterSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRegisterSuite[T Register](
	f *testing.F,
	callable conform.Callable[Register],
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	f.Helper()
	conform.Fuzz(f, RegisterSuite(), callable, registerSubjects(newSubject), opts...)
}

// RunRegisterBudgets applies the allocation ceilings of Register and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRegisterBudgets[T Register](
	t *testing.T,
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	t.Helper()
	conform.RunBudgets(t, RegisterSuite(), registerSubjects(newSubject), opts...)
}

// BenchRegisterSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchRegisterSuite[T Register](
	b conform.B,
	callable conform.Callable[Register],
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	b.Helper()
	conform.Bench(b, RegisterSuite(), callable, registerSubjects(newSubject), opts...)
}

// registerSubjects returns the factory of the subjects of Register that
// calls newSubject with the handle of each check.
func registerSubjects[T Register](
	newSubject func(h *RegisterHandle) T,
) func(h *RegisterHandle) Register {
	return func(h *RegisterHandle) Register { return newSubject(h) }
}

// registerSetCalls generates the calls of Register.Set with the inputs
// that each case draws.
var registerSetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Register) error {
	v := c.Draw(prop.Of[string](), "v")
	return func(ctx context.Context, subject Register) error {
		return subject.Set(ctx, v)
	}
})

// registerGetCalls generates the calls of Register.Get with the inputs
// that each case draws.
var registerGetCalls = prop.Composite(func(*prop.Case) func(context.Context, Register) error {
	return func(ctx context.Context, subject Register) error {
		_, err := subject.Get(ctx)
		return err
	}
})

// registerOverEmbedded is Register under the name of the field that
// registerOver embeds.
type registerOverEmbedded = Register

// registerOver is a subject of Register whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type registerOver struct {
	registerOverEmbedded
	spend func()
}

// Set calls spend and returns its zero results.
func (o registerOver) Set(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o registerOver) Get(context.Context) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newRegisterOver returns the planted defect of a ceiling of Register,
// whose callables call spend.
func newRegisterOver(_ *RegisterHandle, spend func()) Register {
	return registerOver{spend: spend}
}

// registerSetPanicsDefect is the planted defect of a Register whose Set
// panics.
type registerSetPanicsDefect struct {
	registerOverEmbedded
}

// Set is the override of the planted defect.
func (*registerSetPanicsDefect) Set(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Register whose Set panics")
}

// registerSetPanics returns the planted defect of a Register whose Set
// panics over a fresh subject.
func registerSetPanics(h *RegisterHandle) Register {
	return &registerSetPanicsDefect{registerOverEmbedded: h.Fresh()}
}

// registerSetIgnoresCancelDefect is the planted defect of a Register whose
// Set ignores the context.
type registerSetIgnoresCancelDefect struct {
	registerOverEmbedded
}

// Set is the override of the planted defect.
func (d *registerSetIgnoresCancelDefect) Set(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.registerOverEmbedded.Set(p0, p1)
}

// registerSetIgnoresCancel returns the planted defect of a Register whose
// Set ignores the context over a fresh subject.
func registerSetIgnoresCancel(h *RegisterHandle) Register {
	return &registerSetIgnoresCancelDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetPanicsDefect is the planted defect of a Register whose Get
// panics.
type registerGetPanicsDefect struct {
	registerOverEmbedded
}

// Get is the override of the planted defect.
func (*registerGetPanicsDefect) Get(p0 context.Context) (r0 string, r1 error) {
	panic("the planted defect of a Register whose Get panics")
}

// registerGetPanics returns the planted defect of a Register whose Get
// panics over a fresh subject.
func registerGetPanics(h *RegisterHandle) Register {
	return &registerGetPanicsDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetIgnoresCancelDefect is the planted defect of a Register whose
// Get ignores the context.
type registerGetIgnoresCancelDefect struct {
	registerOverEmbedded
}

// Get is the override of the planted defect.
func (d *registerGetIgnoresCancelDefect) Get(p0 context.Context) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registerOverEmbedded.Get(p0)
}

// registerGetIgnoresCancel returns the planted defect of a Register whose
// Get ignores the context over a fresh subject.
func registerGetIgnoresCancel(h *RegisterHandle) Register {
	return &registerGetIgnoresCancelDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetValueBesideDefect is the planted defect of a Register whose
// Get returns a value beside an error.
type registerGetValueBesideDefect struct {
	registerOverEmbedded
}

// Get is the override of the planted defect.
func (d *registerGetValueBesideDefect) Get(p0 context.Context) (r0 string, r1 error) {
	r0, r1 = d.registerOverEmbedded.Get(p0)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// registerGetValueBeside returns the planted defect of a Register whose
// Get returns a value beside an error over a fresh subject.
func registerGetValueBeside(h *RegisterHandle) Register {
	return &registerGetValueBesideDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetDivergesAfterThreeDefect is the planted defect of a Register
// whose first callable with a value result returns another value from its
// fourth call on.
type registerGetDivergesAfterThreeDefect struct {
	registerOverEmbedded
	calls int
}

// Get is the override of the planted defect.
func (d *registerGetDivergesAfterThreeDefect) Get(p0 context.Context) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.registerOverEmbedded.Get(p0)
	d.calls++
	if d.calls > 3 {
		r0 = string("sample-b")
	}
	return r0, r1
}

// registerGetDivergesAfterThree returns the planted defect of a Register
// whose first callable with a value result returns another value from its
// fourth call on over a fresh subject.
func registerGetDivergesAfterThree(h *RegisterHandle) Register {
	return &registerGetDivergesAfterThreeDefect{registerOverEmbedded: h.Fresh()}
}

// registerActions returns the actions of the generated sequences of
// Register over the subject that subject returns. The actions of a case
// take their keys from the pools that it draws.
func registerActions(c *prop.Case, subject func() Register) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Set",
			Input: func(c *prop.Case, _ struct{}) any {
				v := c.Draw(prop.Of[string](), "set.v")
				return []any{v}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Set", args)
				values, err := registerSetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args)
				values, err := registerGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// registerSetCall calls Set on subject with the arguments args, and
// returns the values of the call beside its error.
func registerSetCall(ctx context.Context, subject Register, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Set(ctx, a0)
}

// registerGetCall calls Get on subject with the arguments args, and
// returns the values of the call beside its error.
func registerGetCall(ctx context.Context, subject Register, args []any) ([]any, error) {
	r0, err := subject.Get(ctx)
	return Unstamped("Get", []any{r0}), err
}

// registerReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func registerReference(r Register) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Set":
			values, err = registerSetCall(ctx, r, args)
		case "Get":
			values, err = registerGetCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Double:
//
// Double/smoke: Double returns without a panic for derived inputs.
//
// Double/smoke#generated: Double returns without a panic for derived
// inputs.
//
// Double/cancel: Double returns an error for a cancelled context.
//
// Double/cancel#generated: Double returns an error for a cancelled
// context.
//
// Double/zero-on-error: Double returns the zero value beside an error.
//
// Double/zero-on-error#generated: Double returns the zero value beside an
// error.
//
// Double/agrees-with-the-reference: every sequence of calls returns on
// Double what it returns on the reference.

// The derived checks of Double that witness their laws at the shallow end:
//
// Double/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Double writes.
type (
	DoubleHandle = conform.Handle[DoubleFunc]
	DoubleCheck  = conform.Check[DoubleFunc]
	DoubleOption = conform.Option[DoubleFunc]
)

// DoubleFunc is the type of the function Double. Each subject of the suite
// of Double has it.
type DoubleFunc func(ctx context.Context, n int) (int, error)

// DoubleChecks contains the name of each callable of Double, with the
// identities of its checks.
var DoubleChecks = doubleCallable{"Double"}

// doubleCallable is the name of Double, with the identities of its checks.
type doubleCallable struct {
	conform.Callable[DoubleFunc]
}

// Smoke returns the identity of the check that Double returns without a
// panic for derived inputs.
func (c doubleCallable) Smoke() conform.ID[DoubleFunc] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Double returns
// without a panic for derived inputs.
func (c doubleCallable) SmokeGenerated() conform.ID[DoubleFunc] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Double returns an error
// for a cancelled context.
func (c doubleCallable) Cancel() conform.ID[DoubleFunc] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Double returns an
// error for a cancelled context.
func (c doubleCallable) CancelGenerated() conform.ID[DoubleFunc] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Double returns the
// zero value beside an error.
func (c doubleCallable) ZeroOnError() conform.ID[DoubleFunc] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Double
// returns the zero value beside an error.
func (c doubleCallable) ZeroOnErrorGenerated() conform.ID[DoubleFunc] {
	return c.Check("zero-on-error#generated")
}

// AgreesWithTheReference returns the identity of the check that every
// sequence of calls returns on Double what it returns on the reference.
func (c doubleCallable) AgreesWithTheReference() conform.ID[DoubleFunc] {
	return "Double/agrees-with-the-reference"
}

// DoubleSuite returns the checks of Double, and the calls and the budgets
// of its callables.
func DoubleSuite() conform.Suite[DoubleFunc] {
	return conform.Suite[DoubleFunc]{
		Declaration: "Double",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[DoubleFunc]{
			{
				ID:          DoubleChecks.Smoke(),
				Statement:   "Double returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DoubleHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, int(42))
						return err
					})
				},
				Defect: doublePanics,
			},
			{
				ID:          DoubleChecks.SmokeGenerated(),
				Statement:   "Double returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubleHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnN)
						return err
					})
				},
				Defect: doublePanics,
			},
			{
				ID:          DoubleChecks.Cancel(),
				Statement:   "Double returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DoubleHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, int(42))
						return err
					})
				},
				Defect: doubleIgnoresCancel,
			},
			{
				ID:          DoubleChecks.CancelGenerated(),
				Statement:   "Double returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubleHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnN)
						return err
					})
				},
				Defect: doubleIgnoresCancel,
			},
			{
				ID:          DoubleChecks.ZeroOnError(),
				Statement:   "Double returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DoubleHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject()(ctx, int(42))
					}, int(0))
				},
				Defect: doubleValueBeside,
			},
			{
				ID:          DoubleChecks.ZeroOnErrorGenerated(),
				Statement:   "Double returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubleHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject()(ctx, drawnN)
					}, int(0))
				},
				Defect: doubleValueBeside,
			},
			{
				ID:        DoubleChecks.AgreesWithTheReference(),
				Statement: "every sequence of calls returns on Double what it returns on the reference",
				Claim:     spec.AgreesWithTheReference,
				Probe:     spec.SequencesAgreeWithTheReference,
				Reason:    spec.DivergesFromTheReference,
				Law:       spec.EverySequenceAgreesWithTheReference,
				Machineries: []spec.MachineryID{
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *DoubleHandle, c *prop.Case) {
					probe.SequencesAgreeWithTheReference(h, c, doubleActions, doubleReference)
				},
				Defect: doubleDivergesAfterThree,
			},
		},
		Calls: []conform.Call[DoubleFunc]{
			{
				Calls:    doubleCalls,
				Over:     newDoubleOver,
				Callable: DoubleChecks.Callable,
			},
		},
	}
}

// RunDoubleSuite applies the checks of Double, and the checks that opts
// add, to subject.
func RunDoubleSuite(t *testing.T, subject DoubleFunc, opts ...DoubleOption) {
	t.Helper()
	RunDoubleSuiteFromHandle(t, doubleSubjects(subject), opts...)
}

// RunDoubleSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDoubleSuiteFromHandle(
	t *testing.T,
	newSubject func(h *DoubleHandle) DoubleFunc,
	opts ...DoubleOption,
) {
	t.Helper()
	conform.Run(t, DoubleSuite(), newSubject, opts...)
}

// FuzzDoubleSuite registers on f the fuzz target of the property checks of
// Double on subject.
func FuzzDoubleSuite(f *testing.F, subject DoubleFunc, opts ...DoubleOption) {
	f.Helper()
	conform.Fuzz(f, DoubleSuite(), DoubleChecks.Callable, doubleSubjects(subject), opts...)
}

// RunDoubleBudgets applies the allocation ceiling of Double and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunDoubleBudgets(t *testing.T, subject DoubleFunc, opts ...DoubleOption) {
	t.Helper()
	conform.RunBudgets(t, DoubleSuite(), doubleSubjects(subject), opts...)
}

// BenchDoubleSuite measures subject under the ceilings of the budget of
// Double.
func BenchDoubleSuite(b conform.B, subject DoubleFunc, opts ...DoubleOption) {
	b.Helper()
	conform.Bench(b, DoubleSuite(), DoubleChecks.Callable, doubleSubjects(subject), opts...)
}

// doubleSubjects returns the factory of the subjects of Double that
// returns subject for each handle.
func doubleSubjects(subject DoubleFunc) func(*DoubleHandle) DoubleFunc {
	return func(*DoubleHandle) DoubleFunc { return subject }
}

// doubleCalls generates the calls of Double with the inputs that each case
// draws.
var doubleCalls = prop.Composite(func(c *prop.Case) func(context.Context, DoubleFunc) error {
	n := c.Draw(prop.Of[int](), "n")
	return func(ctx context.Context, subject DoubleFunc) error {
		_, err := subject(ctx, n)
		return err
	}
})

// newDoubleOver returns the planted defect of a ceiling of Double, which
// calls spend and returns its zero results.
func newDoubleOver(_ *DoubleHandle, spend func()) DoubleFunc {
	return func(context.Context, int) (r0 int, r1 error) {
		spend()
		return r0, r1
	}
}

// doublePanics returns the planted defect of a Double whose Double panics
// over the fresh function.
func doublePanics(h *DoubleHandle) DoubleFunc {
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		panic("the planted defect of a Double whose Double panics")
	}
}

// doubleIgnoresCancel returns the planted defect of a Double whose Double
// ignores the context over the fresh function.
func doubleIgnoresCancel(h *DoubleHandle) DoubleFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		if p0.Err() != nil {
			return r0, r1
		}
		return fresh(p0, p1)
	}
}

// doubleValueBeside returns the planted defect of a Double whose Double
// returns a value beside an error over the fresh function.
func doubleValueBeside(h *DoubleHandle) DoubleFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		r0, r1 = fresh(p0, p1)
		if r1 != nil {
			r0 = int(7)
		}
		return r0, r1
	}
}

// doubleDivergesAfterThree returns the planted defect of a Double whose
// first callable with a value result returns another value from its fourth
// call on over the fresh function.
func doubleDivergesAfterThree(h *DoubleHandle) DoubleFunc {
	fresh := h.Fresh()
	calls := 0
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		if err := p0.Err(); err != nil {
			r1 = err
			return r0, r1
		}
		r0, r1 = fresh(p0, p1)
		calls++
		if calls > 3 {
			r0 = int(7)
		}
		return r0, r1
	}
}

// doubleActions returns the actions of the generated sequences of Double
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func doubleActions(c *prop.Case, subject func() DoubleFunc) []stateful.Action[struct{}] {
	keysInt := c.Draw(
		prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Double",
			Input: func(c *prop.Case, _ struct{}) any {
				n := c.Draw(prop.SampledFrom(keysInt...), "double.n")
				return []any{n}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Double", args, args[0])
				values, err := doubleCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// doubleCall calls Double on subject with the arguments args, and returns
// the values of the call beside its error.
func doubleCall(ctx context.Context, subject DoubleFunc, args []any) ([]any, error) {
	a0, _ := args[0].(int)
	r0, err := subject(ctx, a0)
	return []any{r0}, err
}

// doubleReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func doubleReference(r DoubleFunc) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Double":
			values, err = doubleCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Cached:
//
// Cached.Get/smoke: Get returns without a panic for derived inputs.
//
// Cached.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Cached.Get/cancel: Get returns an error for a cancelled context.
//
// Cached.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Cached.Get/zero-on-error: Get returns the zero value beside an error.
//
// Cached.Get/zero-on-error#generated: Get returns the zero value beside an
// error.
//
// Cached.Get/cache-agrees-over-sequences: Get returns what Load returns
// for its arguments, over any sequence of calls.
//
// Cached.Load/smoke: Load returns without a panic for derived inputs.
//
// Cached.Load/smoke#generated: Load returns without a panic for derived
// inputs.
//
// Cached.Load/cancel: Load returns an error for a cancelled context.
//
// Cached.Load/cancel#generated: Load returns an error for a cancelled
// context.
//
// Cached.Load/zero-on-error: Load returns the zero value beside an error.
//
// Cached.Load/zero-on-error#generated: Load returns the zero value beside
// an error.
//
// Cached.Peek/smoke: Peek returns without a panic for derived inputs.
//
// Cached.Peek/smoke#generated: Peek returns without a panic for derived
// inputs.
//
// Cached.Peek/cancel: Peek returns an error for a cancelled context.
//
// Cached.Peek/cancel#generated: Peek returns an error for a cancelled
// context.
//
// Cached.Peek/zero-on-error: Peek returns the zero value beside an error.
//
// Cached.Peek/zero-on-error#generated: Peek returns the zero value beside
// an error.
//
// Cached.Peek/cache-agrees-over-sequences: Peek returns what Fetch returns
// for its arguments, over any sequence of calls.
//
// Cached.Fetch/smoke: Fetch returns without a panic for derived inputs.
//
// Cached.Fetch/smoke#generated: Fetch returns without a panic for derived
// inputs.
//
// Cached.Fetch/cancel: Fetch returns an error for a cancelled context.
//
// Cached.Fetch/cancel#generated: Fetch returns an error for a cancelled
// context.
//
// Cached.Fetch/zero-on-error: Fetch returns the zero value beside an
// error.
//
// Cached.Fetch/zero-on-error#generated: Fetch returns the zero value
// beside an error.

// The derived checks of Cached that witness their laws at the shallow end:
//
// Cached.Get/smoke
//
// Cached.Load/smoke
//
// Cached.Peek/smoke
//
// Cached.Fetch/smoke

// The refused claims of Cached:
//
// Cached.Get/cache: no callable of the declaration is the partner put of
// Get. The remedy is to name the partner of Get with the param of the
// classification.
//
// Cached.Peek/cache: no callable of the declaration is the partner put of
// Peek. The remedy is to name the partner of Peek with the param of the
// classification.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Cached writes.
type (
	CachedHandle = conform.Handle[Cached]
	CachedCheck  = conform.Check[Cached]
	CachedOption = conform.Option[Cached]
)

// CachedChecks contains the name of each callable of Cached, with the
// identities of its checks.
var CachedChecks = struct {
	Get   cachedGet
	Load  cachedLoad
	Peek  cachedPeek
	Fetch cachedFetch
}{
	Get:   cachedGet{"Cached.Get"},
	Load:  cachedLoad{"Cached.Load"},
	Peek:  cachedPeek{"Cached.Peek"},
	Fetch: cachedFetch{"Cached.Fetch"},
}

// cachedGet is the name of Cached.Get, with the identities of its checks.
type cachedGet struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c cachedGet) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c cachedGet) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c cachedGet) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c cachedGet) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c cachedGet) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c cachedGet) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// CacheAgreesOverSequences returns the identity of the check that Get
// returns what Load returns for its arguments, over any sequence of calls.
func (c cachedGet) CacheAgreesOverSequences() conform.ID[Cached] {
	return c.Check("cache-agrees-over-sequences")
}

// cachedLoad is the name of Cached.Load, with the identities of its
// checks.
type cachedLoad struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Load returns without a
// panic for derived inputs.
func (c cachedLoad) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Load returns
// without a panic for derived inputs.
func (c cachedLoad) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Load returns an error for
// a cancelled context.
func (c cachedLoad) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Load returns an
// error for a cancelled context.
func (c cachedLoad) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Load returns the zero
// value beside an error.
func (c cachedLoad) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Load returns
// the zero value beside an error.
func (c cachedLoad) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// cachedPeek is the name of Cached.Peek, with the identities of its
// checks.
type cachedPeek struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Peek returns without a
// panic for derived inputs.
func (c cachedPeek) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Peek returns
// without a panic for derived inputs.
func (c cachedPeek) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Peek returns an error for
// a cancelled context.
func (c cachedPeek) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Peek returns an
// error for a cancelled context.
func (c cachedPeek) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Peek returns the zero
// value beside an error.
func (c cachedPeek) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Peek returns
// the zero value beside an error.
func (c cachedPeek) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// CacheAgreesOverSequences returns the identity of the check that Peek
// returns what Fetch returns for its arguments, over any sequence of
// calls.
func (c cachedPeek) CacheAgreesOverSequences() conform.ID[Cached] {
	return c.Check("cache-agrees-over-sequences")
}

// cachedFetch is the name of Cached.Fetch, with the identities of its
// checks.
type cachedFetch struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Fetch returns without a
// panic for derived inputs.
func (c cachedFetch) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Fetch returns
// without a panic for derived inputs.
func (c cachedFetch) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Fetch returns an error for
// a cancelled context.
func (c cachedFetch) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Fetch returns an
// error for a cancelled context.
func (c cachedFetch) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Fetch returns the
// zero value beside an error.
func (c cachedFetch) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Fetch
// returns the zero value beside an error.
func (c cachedFetch) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// CachedSuite returns the checks of Cached, and the calls and the budgets
// of its callables.
func CachedSuite() conform.Suite[Cached] {
	return conform.Suite[Cached]{
		Declaration: "Cached",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Cached]{
			{
				ID:          CachedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedGetPanics,
			},
			{
				ID:          CachedChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedGetPanics,
			},
			{
				ID:          CachedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedGetIgnoresCancel,
			},
			{
				ID:          CachedChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedGetIgnoresCancel,
			},
			{
				ID:          CachedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: cachedGetValueBeside,
			},
			{
				ID:          CachedChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: cachedGetValueBeside,
			},
			{
				ID:          CachedChecks.Get.CacheAgreesOverSequences(),
				Statement:   "Get returns what Load returns for its arguments, over any sequence of calls",
				Claim:       spec.CacheAgreesOverSequences,
				Probe:       spec.ReadsAgreeWithTheBacking,
				Reason:      spec.CacheDisagrees,
				Law:         spec.CacheAgreesWithItsBacking,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CachedHandle, c *prop.Case) {
					probe.ReadsAgreeWithTheBacking(
						h,
						c,
						cachedActions,
						"Get",
						func(ctx context.Context, args []any) conform.Output {
							values, err := cachedLoadCall(ctx, h.Subject(), []any{string("id-a")})
							return conform.OutputOf(err, values...)
						},
					)
				},
				Defect: cachedGetServesAStaleValueAfterThree,
			},
			{
				ID:          CachedChecks.Load.Smoke(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, string("id-a"))
						return err
					})
				},
				Defect: cachedLoadPanics,
			},
			{
				ID:          CachedChecks.Load.SmokeGenerated(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnID := c.Draw(prop.Of[string](), "id")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, drawnID)
						return err
					})
				},
				Defect: cachedLoadPanics,
			},
			{
				ID:          CachedChecks.Load.Cancel(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, string("id-a"))
						return err
					})
				},
				Defect: cachedLoadIgnoresCancel,
			},
			{
				ID:          CachedChecks.Load.CancelGenerated(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnID := c.Draw(prop.Of[string](), "id")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, drawnID)
						return err
					})
				},
				Defect: cachedLoadIgnoresCancel,
			},
			{
				ID:          CachedChecks.Load.ZeroOnError(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Load(ctx, string("id-a"))
					}, Item{})
				},
				Defect: cachedLoadValueBeside,
			},
			{
				ID:          CachedChecks.Load.ZeroOnErrorGenerated(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnID := c.Draw(prop.Of[string](), "id")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Load(ctx, drawnID)
					}, Item{})
				},
				Defect: cachedLoadValueBeside,
			},
			{
				ID:          CachedChecks.Peek.Smoke(),
				Statement:   "Peek returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedPeekPanics,
			},
			{
				ID:          CachedChecks.Peek.SmokeGenerated(),
				Statement:   "Peek returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedPeekPanics,
			},
			{
				ID:          CachedChecks.Peek.Cancel(),
				Statement:   "Peek returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedPeekIgnoresCancel,
			},
			{
				ID:          CachedChecks.Peek.CancelGenerated(),
				Statement:   "Peek returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedPeekIgnoresCancel,
			},
			{
				ID:          CachedChecks.Peek.ZeroOnError(),
				Statement:   "Peek returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Peek(ctx, string("key-a"))
					}, Item{})
				},
				Defect: cachedPeekValueBeside,
			},
			{
				ID:          CachedChecks.Peek.ZeroOnErrorGenerated(),
				Statement:   "Peek returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Peek(ctx, drawnKey)
					}, Item{})
				},
				Defect: cachedPeekValueBeside,
			},
			{
				ID:          CachedChecks.Peek.CacheAgreesOverSequences(),
				Statement:   "Peek returns what Fetch returns for its arguments, over any sequence of calls",
				Claim:       spec.CacheAgreesOverSequences,
				Probe:       spec.ReadsAgreeWithTheBacking,
				Reason:      spec.CacheDisagrees,
				Law:         spec.CacheAgreesWithItsBacking,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CachedHandle, c *prop.Case) {
					probe.ReadsAgreeWithTheBacking(
						h,
						c,
						cachedActions,
						"Peek",
						func(ctx context.Context, args []any) conform.Output {
							values, err := cachedFetchCall(ctx, h.Subject(), []any{int(42)})
							return conform.OutputOf(err, values...)
						},
					)
				},
				Defect: cachedPeekServesAStaleValueAfterThree,
			},
			{
				ID:          CachedChecks.Fetch.Smoke(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, int(42))
						return err
					})
				},
				Defect: cachedFetchPanics,
			},
			{
				ID:          CachedChecks.Fetch.SmokeGenerated(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[int](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedFetchPanics,
			},
			{
				ID:          CachedChecks.Fetch.Cancel(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, int(42))
						return err
					})
				},
				Defect: cachedFetchIgnoresCancel,
			},
			{
				ID:          CachedChecks.Fetch.CancelGenerated(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[int](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedFetchIgnoresCancel,
			},
			{
				ID:          CachedChecks.Fetch.ZeroOnError(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Fetch(ctx, int(42))
					}, Item{})
				},
				Defect: cachedFetchValueBeside,
			},
			{
				ID:          CachedChecks.Fetch.ZeroOnErrorGenerated(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[int](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Fetch(ctx, drawnKey)
					}, Item{})
				},
				Defect: cachedFetchValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cache,
				Callable: "Cached.Get",
				Gate:     spec.PartnerFound,
				Why:      "no callable of the declaration is the partner put of Get",
				Remedy:   "name the partner of Get with the param of the classification",
			},
			{
				Claim:    spec.Cache,
				Callable: "Cached.Peek",
				Gate:     spec.PartnerFound,
				Why:      "no callable of the declaration is the partner put of Peek",
				Remedy:   "name the partner of Peek with the param of the classification",
			},
		},
		Calls: []conform.Call[Cached]{
			{
				Calls:    cachedGetCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Get.Callable,
			},
			{
				Calls:    cachedLoadCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Load.Callable,
			},
			{
				Calls:    cachedPeekCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Peek.Callable,
			},
			{
				Calls:    cachedFetchCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Fetch.Callable,
			},
		},
	}
}

// RunCachedSuite applies the checks of Cached, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCachedSuite[T Cached](t *testing.T, newSubject func() T, opts ...CachedOption) {
	t.Helper()
	RunCachedSuiteFromHandle(t, func(*CachedHandle) T { return newSubject() }, opts...)
}

// RunCachedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCachedSuiteFromHandle[T Cached](
	t *testing.T,
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	t.Helper()
	conform.Run(t, CachedSuite(), cachedSubjects(newSubject), opts...)
}

// FuzzCachedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCachedSuite[T Cached](
	f *testing.F,
	callable conform.Callable[Cached],
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	f.Helper()
	conform.Fuzz(f, CachedSuite(), callable, cachedSubjects(newSubject), opts...)
}

// RunCachedBudgets applies the allocation ceilings of Cached and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCachedBudgets[T Cached](
	t *testing.T,
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	t.Helper()
	conform.RunBudgets(t, CachedSuite(), cachedSubjects(newSubject), opts...)
}

// BenchCachedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchCachedSuite[T Cached](
	b conform.B,
	callable conform.Callable[Cached],
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	b.Helper()
	conform.Bench(b, CachedSuite(), callable, cachedSubjects(newSubject), opts...)
}

// cachedSubjects returns the factory of the subjects of Cached that calls
// newSubject with the handle of each check.
func cachedSubjects[T Cached](newSubject func(h *CachedHandle) T) func(h *CachedHandle) Cached {
	return func(h *CachedHandle) Cached { return newSubject(h) }
}

// cachedGetCalls generates the calls of Cached.Get with the inputs that
// each case draws.
var cachedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// cachedLoadCalls generates the calls of Cached.Load with the inputs that
// each case draws.
var cachedLoadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	id := c.Draw(prop.Of[string](), "id")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Load(ctx, id)
		return err
	}
})

// cachedPeekCalls generates the calls of Cached.Peek with the inputs that
// each case draws.
var cachedPeekCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Peek(ctx, key)
		return err
	}
})

// cachedFetchCalls generates the calls of Cached.Fetch with the inputs
// that each case draws.
var cachedFetchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[int](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Fetch(ctx, key)
		return err
	}
})

// cachedOverEmbedded is Cached under the name of the field that cachedOver
// embeds.
type cachedOverEmbedded = Cached

// cachedOver is a subject of Cached whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type cachedOver struct {
	cachedOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o cachedOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Load calls spend and returns its zero results.
func (o cachedOver) Load(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Peek calls spend and returns its zero results.
func (o cachedOver) Peek(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Fetch calls spend and returns its zero results.
func (o cachedOver) Fetch(context.Context, int) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newCachedOver returns the planted defect of a ceiling of Cached, whose
// callables call spend.
func newCachedOver(_ *CachedHandle, spend func()) Cached {
	return cachedOver{spend: spend}
}

// cachedGetPanicsDefect is the planted defect of a Cached whose Get
// panics.
type cachedGetPanicsDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (*cachedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Get panics")
}

// cachedGetPanics returns the planted defect of a Cached whose Get panics
// over a fresh subject.
func cachedGetPanics(h *CachedHandle) Cached {
	return &cachedGetPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetIgnoresCancelDefect is the planted defect of a Cached whose Get
// ignores the context.
type cachedGetIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (d *cachedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Get(p0, p1)
}

// cachedGetIgnoresCancel returns the planted defect of a Cached whose Get
// ignores the context over a fresh subject.
func cachedGetIgnoresCancel(h *CachedHandle) Cached {
	return &cachedGetIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetValueBesideDefect is the planted defect of a Cached whose Get
// returns a value beside an error.
type cachedGetValueBesideDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (d *cachedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedGetValueBeside returns the planted defect of a Cached whose Get
// returns a value beside an error over a fresh subject.
func cachedGetValueBeside(h *CachedHandle) Cached {
	return &cachedGetValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetServesAStaleValueAfterThreeDefect is the planted defect of a
// Cached whose Get returns its first value again from its fourth call on.
type cachedGetServesAStaleValueAfterThreeDefect struct {
	cachedOverEmbedded
	calls int
	first *Item
}

// Get is the override of the planted defect.
func (d *cachedGetServesAStaleValueAfterThreeDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.cachedOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// cachedGetServesAStaleValueAfterThree returns the planted defect of a
// Cached whose Get returns its first value again from its fourth call on
// over a fresh subject.
func cachedGetServesAStaleValueAfterThree(h *CachedHandle) Cached {
	return &cachedGetServesAStaleValueAfterThreeDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedLoadPanicsDefect is the planted defect of a Cached whose Load
// panics.
type cachedLoadPanicsDefect struct {
	cachedOverEmbedded
}

// Load is the override of the planted defect.
func (*cachedLoadPanicsDefect) Load(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Load panics")
}

// cachedLoadPanics returns the planted defect of a Cached whose Load
// panics over a fresh subject.
func cachedLoadPanics(h *CachedHandle) Cached {
	return &cachedLoadPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedLoadIgnoresCancelDefect is the planted defect of a Cached whose
// Load ignores the context.
type cachedLoadIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Load is the override of the planted defect.
func (d *cachedLoadIgnoresCancelDefect) Load(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Load(p0, p1)
}

// cachedLoadIgnoresCancel returns the planted defect of a Cached whose
// Load ignores the context over a fresh subject.
func cachedLoadIgnoresCancel(h *CachedHandle) Cached {
	return &cachedLoadIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedLoadValueBesideDefect is the planted defect of a Cached whose Load
// returns a value beside an error.
type cachedLoadValueBesideDefect struct {
	cachedOverEmbedded
}

// Load is the override of the planted defect.
func (d *cachedLoadValueBesideDefect) Load(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Load(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedLoadValueBeside returns the planted defect of a Cached whose Load
// returns a value beside an error over a fresh subject.
func cachedLoadValueBeside(h *CachedHandle) Cached {
	return &cachedLoadValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekPanicsDefect is the planted defect of a Cached whose Peek
// panics.
type cachedPeekPanicsDefect struct {
	cachedOverEmbedded
}

// Peek is the override of the planted defect.
func (*cachedPeekPanicsDefect) Peek(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Peek panics")
}

// cachedPeekPanics returns the planted defect of a Cached whose Peek
// panics over a fresh subject.
func cachedPeekPanics(h *CachedHandle) Cached {
	return &cachedPeekPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekIgnoresCancelDefect is the planted defect of a Cached whose
// Peek ignores the context.
type cachedPeekIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Peek is the override of the planted defect.
func (d *cachedPeekIgnoresCancelDefect) Peek(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Peek(p0, p1)
}

// cachedPeekIgnoresCancel returns the planted defect of a Cached whose
// Peek ignores the context over a fresh subject.
func cachedPeekIgnoresCancel(h *CachedHandle) Cached {
	return &cachedPeekIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekValueBesideDefect is the planted defect of a Cached whose Peek
// returns a value beside an error.
type cachedPeekValueBesideDefect struct {
	cachedOverEmbedded
}

// Peek is the override of the planted defect.
func (d *cachedPeekValueBesideDefect) Peek(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Peek(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedPeekValueBeside returns the planted defect of a Cached whose Peek
// returns a value beside an error over a fresh subject.
func cachedPeekValueBeside(h *CachedHandle) Cached {
	return &cachedPeekValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekServesAStaleValueAfterThreeDefect is the planted defect of a
// Cached whose Peek returns its first value again from its fourth call on.
type cachedPeekServesAStaleValueAfterThreeDefect struct {
	cachedOverEmbedded
	calls int
	first *Item
}

// Peek is the override of the planted defect.
func (d *cachedPeekServesAStaleValueAfterThreeDefect) Peek(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.cachedOverEmbedded.Peek(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// cachedPeekServesAStaleValueAfterThree returns the planted defect of a
// Cached whose Peek returns its first value again from its fourth call on
// over a fresh subject.
func cachedPeekServesAStaleValueAfterThree(h *CachedHandle) Cached {
	return &cachedPeekServesAStaleValueAfterThreeDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchPanicsDefect is the planted defect of a Cached whose Fetch
// panics.
type cachedFetchPanicsDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (*cachedFetchPanicsDefect) Fetch(p0 context.Context, p1 int) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Fetch panics")
}

// cachedFetchPanics returns the planted defect of a Cached whose Fetch
// panics over a fresh subject.
func cachedFetchPanics(h *CachedHandle) Cached {
	return &cachedFetchPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchIgnoresCancelDefect is the planted defect of a Cached whose
// Fetch ignores the context.
type cachedFetchIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *cachedFetchIgnoresCancelDefect) Fetch(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Fetch(p0, p1)
}

// cachedFetchIgnoresCancel returns the planted defect of a Cached whose
// Fetch ignores the context over a fresh subject.
func cachedFetchIgnoresCancel(h *CachedHandle) Cached {
	return &cachedFetchIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchValueBesideDefect is the planted defect of a Cached whose
// Fetch returns a value beside an error.
type cachedFetchValueBesideDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *cachedFetchValueBesideDefect) Fetch(p0 context.Context, p1 int) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Fetch(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedFetchValueBeside returns the planted defect of a Cached whose
// Fetch returns a value beside an error over a fresh subject.
func cachedFetchValueBeside(h *CachedHandle) Cached {
	return &cachedFetchValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedActions returns the actions of the generated sequences of Cached
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func cachedActions(c *prop.Case, subject func() Cached) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	keysInt := c.Draw(
		prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args, args[0])
				values, err := cachedGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Load",
			Input: func(c *prop.Case, _ struct{}) any {
				id := c.Draw(prop.SampledFrom(keysString...), "load.id")
				return []any{id}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Load", args, args[0])
				values, err := cachedLoadCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Peek",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "peek.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Peek", args, args[0])
				values, err := cachedPeekCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Fetch",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysInt...), "fetch.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Fetch", args, args[0])
				values, err := cachedFetchCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// cachedGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func cachedGetCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// cachedLoadCall calls Load on subject with the arguments args, and
// returns the values of the call beside its error.
func cachedLoadCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Load(ctx, a0)
	return []any{r0}, err
}

// cachedPeekCall calls Peek on subject with the arguments args, and
// returns the values of the call beside its error.
func cachedPeekCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Peek(ctx, a0)
	return []any{r0}, err
}

// cachedFetchCall calls Fetch on subject with the arguments args, and
// returns the values of the call beside its error.
func cachedFetchCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(int)
	r0, err := subject.Fetch(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Index:
//
// Index.Put/smoke: Put returns without a panic for derived inputs.
//
// Index.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Index.Put/cancel: Put returns an error for a cancelled context.
//
// Index.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Index.Get/smoke: Get returns without a panic for derived inputs.
//
// Index.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Index.Get/cancel: Get returns an error for a cancelled context.
//
// Index.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Index.Get/zero-on-error: Get returns the zero value beside an error.
//
// Index.Get/zero-on-error#generated: Get returns the zero value beside an
// error.
//
// Index.Get/read-after-write: Get returns what Put stored before it.
//
// Index.Get/read-after-write#generated: Get returns what Put stored before
// it.
//
// Index.Get/reads-each-last-write: Get returns what Put wrote right before
// it, over any sequence of calls.
//
// Index.Find/smoke: Find returns without a panic for derived inputs.
//
// Index.Find/smoke#generated: Find returns without a panic for derived
// inputs.
//
// Index.Find/cancel: Find returns an error for a cancelled context.
//
// Index.Find/cancel#generated: Find returns an error for a cancelled
// context.
//
// Index.Find/read-after-write: Find returns what Put stored before it.
//
// Index.Find/read-after-write#generated: Find returns what Put stored
// before it.
//
// Index.Find/reads-each-last-write: Find returns what Put wrote right
// before it, over any sequence of calls.
//
// Index.ByRev/smoke: ByRev returns without a panic for derived inputs.
//
// Index.ByRev/smoke#generated: ByRev returns without a panic for derived
// inputs.
//
// Index.ByRev/cancel: ByRev returns an error for a cancelled context.
//
// Index.ByRev/cancel#generated: ByRev returns an error for a cancelled
// context.
//
// Index.ByRev/zero-on-error: ByRev returns the zero value beside an error.
//
// Index.ByRev/zero-on-error#generated: ByRev returns the zero value beside
// an error.
//
// Index.ByRev/read-after-write: ByRev returns what Put stored before it.
//
// Index.ByRev/read-after-write#generated: ByRev returns what Put stored
// before it.
//
// Index.ByRev/reads-each-last-write: ByRev returns what Put wrote right
// before it, over any sequence of calls.
//
// Index.Lookup/smoke: Lookup returns without a panic for derived inputs.
//
// Index.Lookup/smoke#generated: Lookup returns without a panic for derived
// inputs.
//
// Index.Lookup/cancel: Lookup returns an error for a cancelled context.
//
// Index.Lookup/cancel#generated: Lookup returns an error for a cancelled
// context.
//
// Index.Lookup/read-after-write: Lookup returns what Put stored before it.
//
// Index.Lookup/read-after-write#generated: Lookup returns what Put stored
// before it.
//
// Index.Lookup/reads-each-last-write: Lookup returns what Put wrote right
// before it, over any sequence of calls.

// The derived checks of Index that witness their laws at the shallow end:
//
// Index.Put/smoke
//
// Index.Get/smoke
//
// Index.Get/read-after-write
//
// Index.Get/read-after-write#generated
//
// Index.Find/smoke
//
// Index.Find/read-after-write
//
// Index.Find/read-after-write#generated
//
// Index.ByRev/smoke
//
// Index.ByRev/read-after-write
//
// Index.ByRev/read-after-write#generated
//
// Index.Lookup/smoke
//
// Index.Lookup/read-after-write
//
// Index.Lookup/read-after-write#generated

// The unclaimed classifications of Index:
//
// The shape:writer of Index.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Index writes.
type (
	IndexHandle = conform.Handle[Index]
	IndexCheck  = conform.Check[Index]
	IndexOption = conform.Option[Index]
)

// IndexChecks contains the name of each callable of Index, with the
// identities of its checks.
var IndexChecks = struct {
	Put    indexPut
	Get    indexGet
	Find   indexFind
	ByRev  indexByRev
	Lookup indexLookup
}{
	Put:    indexPut{"Index.Put"},
	Get:    indexGet{"Index.Get"},
	Find:   indexFind{"Index.Find"},
	ByRev:  indexByRev{"Index.ByRev"},
	Lookup: indexLookup{"Index.Lookup"},
}

// indexPut is the name of Index.Put, with the identities of its checks.
type indexPut struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c indexPut) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c indexPut) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c indexPut) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c indexPut) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// indexGet is the name of Index.Get, with the identities of its checks.
type indexGet struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c indexGet) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c indexGet) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c indexGet) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c indexGet) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c indexGet) ZeroOnError() conform.ID[Index] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c indexGet) ZeroOnErrorGenerated() conform.ID[Index] {
	return c.Check("zero-on-error#generated")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c indexGet) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that Get
// returns what Put stored before it.
func (c indexGet) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Get returns
// what Put wrote right before it, over any sequence of calls.
func (c indexGet) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// indexFind is the name of Index.Find, with the identities of its checks.
type indexFind struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Find returns without a
// panic for derived inputs.
func (c indexFind) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Find returns
// without a panic for derived inputs.
func (c indexFind) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Find returns an error for
// a cancelled context.
func (c indexFind) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Find returns an
// error for a cancelled context.
func (c indexFind) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ReadAfterWrite returns the identity of the check that Find returns what
// Put stored before it.
func (c indexFind) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that Find
// returns what Put stored before it.
func (c indexFind) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Find returns
// what Put wrote right before it, over any sequence of calls.
func (c indexFind) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// indexByRev is the name of Index.ByRev, with the identities of its
// checks.
type indexByRev struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that ByRev returns without a
// panic for derived inputs.
func (c indexByRev) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that ByRev returns
// without a panic for derived inputs.
func (c indexByRev) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that ByRev returns an error for
// a cancelled context.
func (c indexByRev) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that ByRev returns an
// error for a cancelled context.
func (c indexByRev) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that ByRev returns the
// zero value beside an error.
func (c indexByRev) ZeroOnError() conform.ID[Index] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that ByRev
// returns the zero value beside an error.
func (c indexByRev) ZeroOnErrorGenerated() conform.ID[Index] {
	return c.Check("zero-on-error#generated")
}

// ReadAfterWrite returns the identity of the check that ByRev returns what
// Put stored before it.
func (c indexByRev) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that ByRev
// returns what Put stored before it.
func (c indexByRev) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that ByRev returns
// what Put wrote right before it, over any sequence of calls.
func (c indexByRev) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// indexLookup is the name of Index.Lookup, with the identities of its
// checks.
type indexLookup struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Lookup returns without a
// panic for derived inputs.
func (c indexLookup) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Lookup returns
// without a panic for derived inputs.
func (c indexLookup) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Lookup returns an error
// for a cancelled context.
func (c indexLookup) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Lookup returns an
// error for a cancelled context.
func (c indexLookup) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ReadAfterWrite returns the identity of the check that Lookup returns
// what Put stored before it.
func (c indexLookup) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that Lookup
// returns what Put stored before it.
func (c indexLookup) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Lookup returns
// what Put wrote right before it, over any sequence of calls.
func (c indexLookup) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// IndexSuite returns the checks of Index, and the calls and the budgets of
// its callables.
func IndexSuite() conform.Suite[Index] {
	return conform.Suite[Index]{
		Declaration: "Index",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Index]{
			{
				ID:          IndexChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: indexPutPanics,
			},
			{
				ID:          IndexChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: indexPutPanics,
			},
			{
				ID:          IndexChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: indexPutIgnoresCancel,
			},
			{
				ID:          IndexChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: indexPutIgnoresCancel,
			},
			{
				ID:          IndexChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: indexGetPanics,
			},
			{
				ID:          IndexChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: indexGetPanics,
			},
			{
				ID:          IndexChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: indexGetIgnoresCancel,
			},
			{
				ID:          IndexChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: indexGetIgnoresCancel,
			},
			{
				ID:          IndexChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: indexGetValueBeside,
			},
			{
				ID:          IndexChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: indexGetValueBeside,
			},
			{
				ID:          IndexChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, (Item{Key: "Key-a"}).Key)
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexGetReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Get.ReadAfterWriteGenerated(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnIt.Key)
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexGetReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Get.ReadsEachLastWrite(),
				Statement:   "Get returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.Get.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexGetCall(ctx, h.Subject(), []any{w.Key})
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexGetReadsAStaleValueAfterThree,
			},
			{
				ID:          IndexChecks.Find.Smoke(),
				Statement:   "Find returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Find(ctx, string("key-a"), func(Item) bool { return false })
						return err
					})
				},
				Defect: indexFindPanics,
			},
			{
				ID:          IndexChecks.Find.SmokeGenerated(),
				Statement:   "Find returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
					drawnKeep := func() func(Item) bool {
						var keepAt atomic.Int64
						return func(Item) (r0 bool) {
							i := int(keepAt.Add(1) - 1)
							if i >= len(keepResults) {
								return r0
							}
							return keepResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Find(ctx, drawnKey, drawnKeep())
						return err
					})
				},
				Defect: indexFindPanics,
			},
			{
				ID:          IndexChecks.Find.Cancel(),
				Statement:   "Find returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Find(ctx, string("key-a"), func(Item) bool { return false })
						return err
					})
				},
				Defect: indexFindIgnoresCancel,
			},
			{
				ID:          IndexChecks.Find.CancelGenerated(),
				Statement:   "Find returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
					drawnKeep := func() func(Item) bool {
						var keepAt atomic.Int64
						return func(Item) (r0 bool) {
							i := int(keepAt.Add(1) - 1)
							if i >= len(keepResults) {
								return r0
							}
							return keepResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Find(ctx, drawnKey, drawnKeep())
						return err
					})
				},
				Defect: indexFindIgnoresCancel,
			},
			{
				ID:          IndexChecks.Find.ReadAfterWrite(),
				Statement:   "Find returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().
							Find(ctx, (Item{Key: "Key-a"}).Key, func(Item) bool { return false })
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexFindReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Find.ReadAfterWriteGenerated(),
				Statement:   "Find returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
					drawnKeep := func() func(Item) bool {
						var keepAt atomic.Int64
						return func(Item) (r0 bool) {
							i := int(keepAt.Add(1) - 1)
							if i >= len(keepResults) {
								return r0
							}
							return keepResults[i]
						}
					}
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Find(ctx, drawnIt.Key, drawnKeep())
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexFindReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Find.ReadsEachLastWrite(),
				Statement:   "Find returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.Find.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexFindCall(
								ctx,
								h.Subject(),
								[]any{w.Key, func(Item) bool { return false }},
							)
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexFindReadsAStaleValueAfterThree,
			},
			{
				ID:          IndexChecks.ByRev.Smoke(),
				Statement:   "ByRev returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, int(42))
						return err
					})
				},
				Defect: indexByRevPanics,
			},
			{
				ID:          IndexChecks.ByRev.SmokeGenerated(),
				Statement:   "ByRev returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, drawnRev)
						return err
					})
				},
				Defect: indexByRevPanics,
			},
			{
				ID:          IndexChecks.ByRev.Cancel(),
				Statement:   "ByRev returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, int(42))
						return err
					})
				},
				Defect: indexByRevIgnoresCancel,
			},
			{
				ID:          IndexChecks.ByRev.CancelGenerated(),
				Statement:   "ByRev returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, drawnRev)
						return err
					})
				},
				Defect: indexByRevIgnoresCancel,
			},
			{
				ID:          IndexChecks.ByRev.ZeroOnError(),
				Statement:   "ByRev returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, int(42))
					}, Item{})
				},
				Defect: indexByRevValueBeside,
			},
			{
				ID:          IndexChecks.ByRev.ZeroOnErrorGenerated(),
				Statement:   "ByRev returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, drawnRev)
					}, Item{})
				},
				Defect: indexByRevValueBeside,
			},
			{
				ID:          IndexChecks.ByRev.ReadAfterWrite(),
				Statement:   "ByRev returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, (Item{Key: "Key-a"}).Rev)
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexByRevReadsAnotherValue,
			},
			{
				ID:          IndexChecks.ByRev.ReadAfterWriteGenerated(),
				Statement:   "ByRev returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, drawnIt.Rev)
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexByRevReadsAnotherValue,
			},
			{
				ID:          IndexChecks.ByRev.ReadsEachLastWrite(),
				Statement:   "ByRev returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.ByRev.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexByRevCall(ctx, h.Subject(), []any{w.Rev})
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexByRevReadsAStaleValueAfterThree,
			},
			{
				ID:          IndexChecks.Lookup.Smoke(),
				Statement:   "Lookup returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, string("key-a"), int(42))
						return err
					})
				},
				Defect: indexLookupPanics,
			},
			{
				ID:          IndexChecks.Lookup.SmokeGenerated(),
				Statement:   "Lookup returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, drawnKey, drawnRev)
						return err
					})
				},
				Defect: indexLookupPanics,
			},
			{
				ID:          IndexChecks.Lookup.Cancel(),
				Statement:   "Lookup returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, string("key-a"), int(42))
						return err
					})
				},
				Defect: indexLookupIgnoresCancel,
			},
			{
				ID:          IndexChecks.Lookup.CancelGenerated(),
				Statement:   "Lookup returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, drawnKey, drawnRev)
						return err
					})
				},
				Defect: indexLookupIgnoresCancel,
			},
			{
				ID:          IndexChecks.Lookup.ReadAfterWrite(),
				Statement:   "Lookup returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().
							Lookup(ctx, (Item{Key: "Key-a"}).Key, (Item{Key: "Key-a"}).Rev)
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexLookupReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Lookup.ReadAfterWriteGenerated(),
				Statement:   "Lookup returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Lookup(ctx, drawnIt.Key, drawnIt.Rev)
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexLookupReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Lookup.ReadsEachLastWrite(),
				Statement:   "Lookup returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.Lookup.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexLookupCall(ctx, h.Subject(), []any{w.Key, w.Rev})
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexLookupReadsAStaleValueAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Index.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Index]{
			{
				Calls:    indexPutCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Put.Callable,
			},
			{
				Calls:    indexGetCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Get.Callable,
			},
			{
				Calls:    indexFindCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Find.Callable,
			},
			{
				Calls:    indexByRevCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.ByRev.Callable,
			},
			{
				Calls:    indexLookupCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Lookup.Callable,
			},
		},
	}
}

// RunIndexSuite applies the checks of Index, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunIndexSuite[T Index](t *testing.T, newSubject func() T, opts ...IndexOption) {
	t.Helper()
	RunIndexSuiteFromHandle(t, func(*IndexHandle) T { return newSubject() }, opts...)
}

// RunIndexSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunIndexSuiteFromHandle[T Index](
	t *testing.T,
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	t.Helper()
	conform.Run(t, IndexSuite(), indexSubjects(newSubject), opts...)
}

// FuzzIndexSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzIndexSuite[T Index](
	f *testing.F,
	callable conform.Callable[Index],
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	f.Helper()
	conform.Fuzz(f, IndexSuite(), callable, indexSubjects(newSubject), opts...)
}

// RunIndexBudgets applies the allocation ceilings of Index and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunIndexBudgets[T Index](
	t *testing.T,
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	t.Helper()
	conform.RunBudgets(t, IndexSuite(), indexSubjects(newSubject), opts...)
}

// BenchIndexSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchIndexSuite[T Index](
	b conform.B,
	callable conform.Callable[Index],
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	b.Helper()
	conform.Bench(b, IndexSuite(), callable, indexSubjects(newSubject), opts...)
}

// indexSubjects returns the factory of the subjects of Index that calls
// newSubject with the handle of each check.
func indexSubjects[T Index](newSubject func(h *IndexHandle) T) func(h *IndexHandle) Index {
	return func(h *IndexHandle) Index { return newSubject(h) }
}

// indexPutCalls generates the calls of Index.Put with the inputs that each
// case draws.
var indexPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Index) error {
		return subject.Put(ctx, it)
	}
})

// indexGetCalls generates the calls of Index.Get with the inputs that each
// case draws.
var indexGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Index) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// indexFindCalls generates the calls of Index.Find with the inputs that
// each case draws.
var indexFindCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	key := c.Draw(prop.Of[string](), "key")
	keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
	var keepAt atomic.Int64
	keep := func(Item) (r0 bool) {
		i := int(keepAt.Add(1) - 1)
		if i >= len(keepResults) {
			return r0
		}
		return keepResults[i]
	}
	return func(ctx context.Context, subject Index) error {
		_, err := subject.Find(ctx, key, keep)
		return err
	}
})

// indexByRevCalls generates the calls of Index.ByRev with the inputs that
// each case draws.
var indexByRevCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	rev := c.Draw(prop.Of[int](), "rev")
	return func(ctx context.Context, subject Index) error {
		_, err := subject.ByRev(ctx, rev)
		return err
	}
})

// indexLookupCalls generates the calls of Index.Lookup with the inputs
// that each case draws.
var indexLookupCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	key := c.Draw(prop.Of[string](), "key")
	rev := c.Draw(prop.Of[int](), "rev")
	return func(ctx context.Context, subject Index) error {
		_, err := subject.Lookup(ctx, key, rev)
		return err
	}
})

// indexOverEmbedded is Index under the name of the field that indexOver
// embeds.
type indexOverEmbedded = Index

// indexOver is a subject of Index whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type indexOver struct {
	indexOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o indexOver) Put(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o indexOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Find calls spend and returns its zero results.
func (o indexOver) Find(context.Context, string, func(Item) bool) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// ByRev calls spend and returns its zero results.
func (o indexOver) ByRev(context.Context, int) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Lookup calls spend and returns its zero results.
func (o indexOver) Lookup(context.Context, string, int) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newIndexOver returns the planted defect of a ceiling of Index, whose
// callables call spend.
func newIndexOver(_ *IndexHandle, spend func()) Index {
	return indexOver{spend: spend}
}

// indexPutPanicsDefect is the planted defect of a Index whose Put panics.
type indexPutPanicsDefect struct {
	indexOverEmbedded
}

// Put is the override of the planted defect.
func (*indexPutPanicsDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Index whose Put panics")
}

// indexPutPanics returns the planted defect of a Index whose Put panics
// over a fresh subject.
func indexPutPanics(h *IndexHandle) Index {
	return &indexPutPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexPutIgnoresCancelDefect is the planted defect of a Index whose Put
// ignores the context.
type indexPutIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Put is the override of the planted defect.
func (d *indexPutIgnoresCancelDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.indexOverEmbedded.Put(p0, p1)
}

// indexPutIgnoresCancel returns the planted defect of a Index whose Put
// ignores the context over a fresh subject.
func indexPutIgnoresCancel(h *IndexHandle) Index {
	return &indexPutIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetPanicsDefect is the planted defect of a Index whose Get panics.
type indexGetPanicsDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (*indexGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose Get panics")
}

// indexGetPanics returns the planted defect of a Index whose Get panics
// over a fresh subject.
func indexGetPanics(h *IndexHandle) Index {
	return &indexGetPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetIgnoresCancelDefect is the planted defect of a Index whose Get
// ignores the context.
type indexGetIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (d *indexGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.Get(p0, p1)
}

// indexGetIgnoresCancel returns the planted defect of a Index whose Get
// ignores the context over a fresh subject.
func indexGetIgnoresCancel(h *IndexHandle) Index {
	return &indexGetIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetValueBesideDefect is the planted defect of a Index whose Get
// returns a value beside an error.
type indexGetValueBesideDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (d *indexGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.indexOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexGetValueBeside returns the planted defect of a Index whose Get
// returns a value beside an error over a fresh subject.
func indexGetValueBeside(h *IndexHandle) Index {
	return &indexGetValueBesideDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetReadsAnotherValueDefect is the planted defect of a Index whose
// Get returns a value that Put did not write.
type indexGetReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (d *indexGetReadsAnotherValueDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexGetReadsAnotherValue returns the planted defect of a Index whose
// Get returns a value that Put did not write over a fresh subject.
func indexGetReadsAnotherValue(h *IndexHandle) Index {
	return &indexGetReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose Get returns its first value again from its fourth call on.
type indexGetReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// Get is the override of the planted defect.
func (d *indexGetReadsAStaleValueAfterThreeDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexGetReadsAStaleValueAfterThree returns the planted defect of a Index
// whose Get returns its first value again from its fourth call on over a
// fresh subject.
func indexGetReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexGetReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindPanicsDefect is the planted defect of a Index whose Find
// panics.
type indexFindPanicsDefect struct {
	indexOverEmbedded
}

// Find is the override of the planted defect.
func (*indexFindPanicsDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose Find panics")
}

// indexFindPanics returns the planted defect of a Index whose Find panics
// over a fresh subject.
func indexFindPanics(h *IndexHandle) Index {
	return &indexFindPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindIgnoresCancelDefect is the planted defect of a Index whose Find
// ignores the context.
type indexFindIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Find is the override of the planted defect.
func (d *indexFindIgnoresCancelDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.Find(p0, p1, p2)
}

// indexFindIgnoresCancel returns the planted defect of a Index whose Find
// ignores the context over a fresh subject.
func indexFindIgnoresCancel(h *IndexHandle) Index {
	return &indexFindIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindReadsAnotherValueDefect is the planted defect of a Index whose
// Find returns a value that Put did not write.
type indexFindReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// Find is the override of the planted defect.
func (d *indexFindReadsAnotherValueDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.Find(p0, p1, p2)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexFindReadsAnotherValue returns the planted defect of a Index whose
// Find returns a value that Put did not write over a fresh subject.
func indexFindReadsAnotherValue(h *IndexHandle) Index {
	return &indexFindReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose Find returns its first value again from its fourth call on.
type indexFindReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// Find is the override of the planted defect.
func (d *indexFindReadsAStaleValueAfterThreeDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.Find(p0, p1, p2)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexFindReadsAStaleValueAfterThree returns the planted defect of a
// Index whose Find returns its first value again from its fourth call on
// over a fresh subject.
func indexFindReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexFindReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevPanicsDefect is the planted defect of a Index whose ByRev
// panics.
type indexByRevPanicsDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (*indexByRevPanicsDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose ByRev panics")
}

// indexByRevPanics returns the planted defect of a Index whose ByRev
// panics over a fresh subject.
func indexByRevPanics(h *IndexHandle) Index {
	return &indexByRevPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevIgnoresCancelDefect is the planted defect of a Index whose
// ByRev ignores the context.
type indexByRevIgnoresCancelDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (d *indexByRevIgnoresCancelDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.ByRev(p0, p1)
}

// indexByRevIgnoresCancel returns the planted defect of a Index whose
// ByRev ignores the context over a fresh subject.
func indexByRevIgnoresCancel(h *IndexHandle) Index {
	return &indexByRevIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevValueBesideDefect is the planted defect of a Index whose ByRev
// returns a value beside an error.
type indexByRevValueBesideDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (d *indexByRevValueBesideDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	r0, r1 = d.indexOverEmbedded.ByRev(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexByRevValueBeside returns the planted defect of a Index whose ByRev
// returns a value beside an error over a fresh subject.
func indexByRevValueBeside(h *IndexHandle) Index {
	return &indexByRevValueBesideDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevReadsAnotherValueDefect is the planted defect of a Index whose
// ByRev returns a value that Put did not write.
type indexByRevReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (d *indexByRevReadsAnotherValueDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.ByRev(p0, p1)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexByRevReadsAnotherValue returns the planted defect of a Index whose
// ByRev returns a value that Put did not write over a fresh subject.
func indexByRevReadsAnotherValue(h *IndexHandle) Index {
	return &indexByRevReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose ByRev returns its first value again from its fourth call on.
type indexByRevReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// ByRev is the override of the planted defect.
func (d *indexByRevReadsAStaleValueAfterThreeDefect) ByRev(
	p0 context.Context,
	p1 int,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.ByRev(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexByRevReadsAStaleValueAfterThree returns the planted defect of a
// Index whose ByRev returns its first value again from its fourth call on
// over a fresh subject.
func indexByRevReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexByRevReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupPanicsDefect is the planted defect of a Index whose Lookup
// panics.
type indexLookupPanicsDefect struct {
	indexOverEmbedded
}

// Lookup is the override of the planted defect.
func (*indexLookupPanicsDefect) Lookup(p0 context.Context, p1 string, p2 int) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose Lookup panics")
}

// indexLookupPanics returns the planted defect of a Index whose Lookup
// panics over a fresh subject.
func indexLookupPanics(h *IndexHandle) Index {
	return &indexLookupPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupIgnoresCancelDefect is the planted defect of a Index whose
// Lookup ignores the context.
type indexLookupIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Lookup is the override of the planted defect.
func (d *indexLookupIgnoresCancelDefect) Lookup(
	p0 context.Context,
	p1 string,
	p2 int,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.Lookup(p0, p1, p2)
}

// indexLookupIgnoresCancel returns the planted defect of a Index whose
// Lookup ignores the context over a fresh subject.
func indexLookupIgnoresCancel(h *IndexHandle) Index {
	return &indexLookupIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupReadsAnotherValueDefect is the planted defect of a Index
// whose Lookup returns a value that Put did not write.
type indexLookupReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// Lookup is the override of the planted defect.
func (d *indexLookupReadsAnotherValueDefect) Lookup(
	p0 context.Context,
	p1 string,
	p2 int,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.Lookup(p0, p1, p2)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexLookupReadsAnotherValue returns the planted defect of a Index whose
// Lookup returns a value that Put did not write over a fresh subject.
func indexLookupReadsAnotherValue(h *IndexHandle) Index {
	return &indexLookupReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose Lookup returns its first value again from its fourth call
// on.
type indexLookupReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// Lookup is the override of the planted defect.
func (d *indexLookupReadsAStaleValueAfterThreeDefect) Lookup(
	p0 context.Context,
	p1 string,
	p2 int,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.Lookup(p0, p1, p2)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexLookupReadsAStaleValueAfterThree returns the planted defect of a
// Index whose Lookup returns its first value again from its fourth call on
// over a fresh subject.
func indexLookupReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexLookupReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexActions returns the actions of the generated sequences of Index
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func indexActions(c *prop.Case, subject func() Index) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	keysInt := c.Draw(
		prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "put.it")
				it.Key = c.Draw(prop.SampledFrom(keysString...), "put.it.key")
				it.Rev = c.Draw(prop.SampledFrom(keysInt...), "put.it.rev")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				key0, _ := args[0].(Item)
				key1, _ := args[0].(Item)
				call := c.History().Invoke(client, "Put", args, key0.Key, key1.Rev)
				values, err := indexPutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args, args[0])
				values, err := indexGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Find",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "find.key")
				keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
				var keepAt atomic.Int64
				keep := func(Item) (r0 bool) {
					i := int(keepAt.Add(1) - 1)
					if i >= len(keepResults) {
						return r0
					}
					return keepResults[i]
				}
				return []any{key, keep}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Find", args)
				values, err := indexFindCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "ByRev",
			Input: func(c *prop.Case, _ struct{}) any {
				rev := c.Draw(prop.SampledFrom(keysInt...), "by-rev.rev")
				return []any{rev}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "ByRev", args, args[0])
				values, err := indexByRevCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Lookup",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "lookup.key")
				rev := c.Draw(prop.Of[int](), "lookup.rev")
				return []any{key, rev}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Lookup", args)
				values, err := indexLookupCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// indexPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func indexPutCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Put(ctx, a0)
}

// indexGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func indexGetCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// indexFindCall calls Find on subject with the arguments args, and returns
// the values of the call beside its error.
func indexFindCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(func(Item) bool)
	r0, err := subject.Find(ctx, a0, a1)
	return []any{r0}, err
}

// indexByRevCall calls ByRev on subject with the arguments args, and
// returns the values of the call beside its error.
func indexByRevCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(int)
	r0, err := subject.ByRev(ctx, a0)
	return []any{r0}, err
}

// indexLookupCall calls Lookup on subject with the arguments args, and
// returns the values of the call beside its error.
func indexLookupCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(int)
	r0, err := subject.Lookup(ctx, a0, a1)
	return []any{r0}, err
}

// The derived checks of Stamped:
//
// Stamped.Store/smoke: Store returns without a panic for derived inputs.
//
// Stamped.Store/smoke#generated: Store returns without a panic for derived
// inputs.
//
// Stamped.Store/cancel: Store returns an error for a cancelled context.
//
// Stamped.Store/cancel#generated: Store returns an error for a cancelled
// context.
//
// Stamped.Get/smoke: Get returns without a panic for derived inputs.
//
// Stamped.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Stamped.Get/cancel: Get returns an error for a cancelled context.
//
// Stamped.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Stamped.Get/zero-on-error: Get returns the zero value beside an error.
//
// Stamped.Get/zero-on-error#generated: Get returns the zero value beside
// an error.
//
// Stamped.Get/monotonic-reads: Get never returns an older version than one
// that it returned before.
//
// Stamped.Get/monotonic-reads#generated: Get never returns an older
// version than one that it returned before.
//
// Stamped.Get/reads-never-go-back: Get never returns a version below one
// that it returned before, over any sequence of calls.

// The derived checks of Stamped that witness their laws at the shallow
// end:
//
// Stamped.Store/smoke
//
// Stamped.Get/smoke
//
// Stamped.Get/monotonic-reads
//
// Stamped.Get/monotonic-reads#generated

// The unclaimed classifications of Stamped:
//
// The shape:writer of Stamped.Store has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Stamped writes.
type (
	StampedHandle = conform.Handle[Stamped]
	StampedCheck  = conform.Check[Stamped]
	StampedOption = conform.Option[Stamped]
)

// StampedChecks contains the name of each callable of Stamped, with the
// identities of its checks.
var StampedChecks = struct {
	Store stampedStore
	Get   stampedGet
}{
	Store: stampedStore{"Stamped.Store"},
	Get:   stampedGet{"Stamped.Get"},
}

// stampedStore is the name of Stamped.Store, with the identities of its
// checks.
type stampedStore struct {
	conform.Callable[Stamped]
}

// Smoke returns the identity of the check that Store returns without a
// panic for derived inputs.
func (c stampedStore) Smoke() conform.ID[Stamped] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Store returns
// without a panic for derived inputs.
func (c stampedStore) SmokeGenerated() conform.ID[Stamped] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Store returns an error for
// a cancelled context.
func (c stampedStore) Cancel() conform.ID[Stamped] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Store returns an
// error for a cancelled context.
func (c stampedStore) CancelGenerated() conform.ID[Stamped] {
	return c.Check("cancel#generated")
}

// stampedGet is the name of Stamped.Get, with the identities of its
// checks.
type stampedGet struct {
	conform.Callable[Stamped]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c stampedGet) Smoke() conform.ID[Stamped] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c stampedGet) SmokeGenerated() conform.ID[Stamped] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c stampedGet) Cancel() conform.ID[Stamped] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c stampedGet) CancelGenerated() conform.ID[Stamped] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c stampedGet) ZeroOnError() conform.ID[Stamped] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c stampedGet) ZeroOnErrorGenerated() conform.ID[Stamped] {
	return c.Check("zero-on-error#generated")
}

// MonotonicReads returns the identity of the check that Get never returns
// an older version than one that it returned before.
func (c stampedGet) MonotonicReads() conform.ID[Stamped] {
	return c.Check("monotonic-reads")
}

// MonotonicReadsGenerated returns the identity of the check that Get never
// returns an older version than one that it returned before.
func (c stampedGet) MonotonicReadsGenerated() conform.ID[Stamped] {
	return c.Check("monotonic-reads#generated")
}

// ReadsNeverGoBack returns the identity of the check that Get never
// returns a version below one that it returned before, over any sequence
// of calls.
func (c stampedGet) ReadsNeverGoBack() conform.ID[Stamped] {
	return c.Check("reads-never-go-back")
}

// StampedSuite returns the checks of Stamped, and the calls and the
// budgets of its callables.
func StampedSuite() conform.Suite[Stamped] {
	return conform.Suite[Stamped]{
		Declaration: "Stamped",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Stamped]{
			{
				ID:          StampedChecks.Store.Smoke(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: stampedStorePanics,
			},
			{
				ID:          StampedChecks.Store.SmokeGenerated(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnIt)
					})
				},
				Defect: stampedStorePanics,
			},
			{
				ID:          StampedChecks.Store.Cancel(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: stampedStoreIgnoresCancel,
			},
			{
				ID:          StampedChecks.Store.CancelGenerated(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnIt)
					})
				},
				Defect: stampedStoreIgnoresCancel,
			},
			{
				ID:          StampedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: stampedGetPanics,
			},
			{
				ID:          StampedChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: stampedGetPanics,
			},
			{
				ID:          StampedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: stampedGetIgnoresCancel,
			},
			{
				ID:          StampedChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: stampedGetIgnoresCancel,
			},
			{
				ID:          StampedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: stampedGetValueBeside,
			},
			{
				ID:          StampedChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: stampedGetValueBeside,
			},
			{
				ID:          StampedChecks.Get.MonotonicReads(),
				Statement:   "Get never returns an older version than one that it returned before",
				Claim:       spec.MonotonicReads,
				Probe:       spec.ReadsNeverRegress,
				Reason:      spec.Regressed,
				Law:         spec.ReadsNeverGoBackwards,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.ReadsNeverRegress(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, (Item{Key: "Key-a"}).Key)
					}, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Item{Key: "Key-a"})
					}, func(v Item) int { return v.Rev })
				},
				Defect: stampedGetVersionDescends,
			},
			{
				ID:          StampedChecks.Get.MonotonicReadsGenerated(),
				Statement:   "Get never returns an older version than one that it returned before",
				Claim:       spec.MonotonicReads,
				Probe:       spec.ReadsNeverRegress,
				Reason:      spec.Regressed,
				Law:         spec.ReadsNeverGoBackwards,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.ReadsNeverRegress(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnIt.Key)
					}, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnIt)
					}, func(v Item) int { return v.Rev })
				},
				Defect: stampedGetVersionDescends,
			},
			{
				ID:          StampedChecks.Get.ReadsNeverGoBack(),
				Statement:   "Get never returns a version below one that it returned before, over any sequence of calls",
				Claim:       spec.ReadsNeverGoBack,
				Probe:       spec.ReadVersionsNeverFall,
				Reason:      spec.Regressed,
				Law:         spec.ReadsNeverGoBackwards,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Stamped]{StampedChecks.Get.MonotonicReads()},
				Property: func(h *StampedHandle, c *prop.Case) {
					probe.ReadVersionsNeverFall(
						h,
						c,
						stampedActions,
						"Get",
						func(v Item) int { return v.Rev },
					)
				},
				Defect: stampedGetReadsThePastAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Stamped.Store",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Stamped]{
			{
				Calls:    stampedStoreCalls,
				Over:     newStampedOver,
				Callable: StampedChecks.Store.Callable,
			},
			{
				Calls:    stampedGetCalls,
				Over:     newStampedOver,
				Callable: StampedChecks.Get.Callable,
			},
		},
	}
}

// RunStampedSuite applies the checks of Stamped, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunStampedSuite[T Stamped](t *testing.T, newSubject func() T, opts ...StampedOption) {
	t.Helper()
	RunStampedSuiteFromHandle(t, func(*StampedHandle) T { return newSubject() }, opts...)
}

// RunStampedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStampedSuiteFromHandle[T Stamped](
	t *testing.T,
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	t.Helper()
	conform.Run(t, StampedSuite(), stampedSubjects(newSubject), opts...)
}

// FuzzStampedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzStampedSuite[T Stamped](
	f *testing.F,
	callable conform.Callable[Stamped],
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	f.Helper()
	conform.Fuzz(f, StampedSuite(), callable, stampedSubjects(newSubject), opts...)
}

// RunStampedBudgets applies the allocation ceilings of Stamped and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunStampedBudgets[T Stamped](
	t *testing.T,
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	t.Helper()
	conform.RunBudgets(t, StampedSuite(), stampedSubjects(newSubject), opts...)
}

// BenchStampedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchStampedSuite[T Stamped](
	b conform.B,
	callable conform.Callable[Stamped],
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	b.Helper()
	conform.Bench(b, StampedSuite(), callable, stampedSubjects(newSubject), opts...)
}

// stampedSubjects returns the factory of the subjects of Stamped that
// calls newSubject with the handle of each check.
func stampedSubjects[T Stamped](
	newSubject func(h *StampedHandle) T,
) func(h *StampedHandle) Stamped {
	return func(h *StampedHandle) Stamped { return newSubject(h) }
}

// stampedStoreCalls generates the calls of Stamped.Store with the inputs
// that each case draws.
var stampedStoreCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stamped) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Stamped) error {
		return subject.Store(ctx, it)
	}
})

// stampedGetCalls generates the calls of Stamped.Get with the inputs that
// each case draws.
var stampedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stamped) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Stamped) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// stampedOverEmbedded is Stamped under the name of the field that
// stampedOver embeds.
type stampedOverEmbedded = Stamped

// stampedOver is a subject of Stamped whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type stampedOver struct {
	stampedOverEmbedded
	spend func()
}

// Store calls spend and returns its zero results.
func (o stampedOver) Store(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o stampedOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newStampedOver returns the planted defect of a ceiling of Stamped, whose
// callables call spend.
func newStampedOver(_ *StampedHandle, spend func()) Stamped {
	return stampedOver{spend: spend}
}

// stampedStorePanicsDefect is the planted defect of a Stamped whose Store
// panics.
type stampedStorePanicsDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (*stampedStorePanicsDefect) Store(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Stamped whose Store panics")
}

// stampedStorePanics returns the planted defect of a Stamped whose Store
// panics over a fresh subject.
func stampedStorePanics(h *StampedHandle) Stamped {
	return &stampedStorePanicsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedStoreIgnoresCancelDefect is the planted defect of a Stamped whose
// Store ignores the context.
type stampedStoreIgnoresCancelDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (d *stampedStoreIgnoresCancelDefect) Store(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.stampedOverEmbedded.Store(p0, p1)
}

// stampedStoreIgnoresCancel returns the planted defect of a Stamped whose
// Store ignores the context over a fresh subject.
func stampedStoreIgnoresCancel(h *StampedHandle) Stamped {
	return &stampedStoreIgnoresCancelDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetPanicsDefect is the planted defect of a Stamped whose Get
// panics.
type stampedGetPanicsDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (*stampedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Stamped whose Get panics")
}

// stampedGetPanics returns the planted defect of a Stamped whose Get
// panics over a fresh subject.
func stampedGetPanics(h *StampedHandle) Stamped {
	return &stampedGetPanicsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetIgnoresCancelDefect is the planted defect of a Stamped whose
// Get ignores the context.
type stampedGetIgnoresCancelDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (d *stampedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.stampedOverEmbedded.Get(p0, p1)
}

// stampedGetIgnoresCancel returns the planted defect of a Stamped whose
// Get ignores the context over a fresh subject.
func stampedGetIgnoresCancel(h *StampedHandle) Stamped {
	return &stampedGetIgnoresCancelDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetValueBesideDefect is the planted defect of a Stamped whose Get
// returns a value beside an error.
type stampedGetValueBesideDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (d *stampedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.stampedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// stampedGetValueBeside returns the planted defect of a Stamped whose Get
// returns a value beside an error over a fresh subject.
func stampedGetValueBeside(h *StampedHandle) Stamped {
	return &stampedGetValueBesideDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetVersionDescendsDefect is the planted defect of a Stamped whose
// Get returns an older version on the second read.
type stampedGetVersionDescendsDefect struct {
	stampedOverEmbedded
	calls int
}

// Get is the override of the planted defect.
func (d *stampedGetVersionDescendsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0 = Item{Key: "Key-b"}
	if d.calls == 1 {
		r0.Rev = 1
	} else {
		r0.Rev = 0
	}
	return r0, r1
}

// stampedGetVersionDescends returns the planted defect of a Stamped whose
// Get returns an older version on the second read over a fresh subject.
func stampedGetVersionDescends(h *StampedHandle) Stamped {
	return &stampedGetVersionDescendsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetReadsThePastAfterThreeDefect is the planted defect of a
// Stamped whose Get returns its first value again from its fourth call on.
type stampedGetReadsThePastAfterThreeDefect struct {
	stampedOverEmbedded
	calls int
	first *Item
}

// Get is the override of the planted defect.
func (d *stampedGetReadsThePastAfterThreeDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.stampedOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// stampedGetReadsThePastAfterThree returns the planted defect of a Stamped
// whose Get returns its first value again from its fourth call on over a
// fresh subject.
func stampedGetReadsThePastAfterThree(h *StampedHandle) Stamped {
	return &stampedGetReadsThePastAfterThreeDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedActions returns the actions of the generated sequences of Stamped
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func stampedActions(c *prop.Case, subject func() Stamped) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Store",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "store.it")
				it.Key = c.Draw(prop.SampledFrom(keysString...), "store.it.key")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				key0, _ := args[0].(Item)
				call := c.History().Invoke(client, "Store", args, key0.Key)
				values, err := stampedStoreCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args, args[0])
				values, err := stampedGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// stampedStoreCall calls Store on subject with the arguments args, and
// returns the values of the call beside its error.
func stampedStoreCall(ctx context.Context, subject Stamped, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Store(ctx, a0)
}

// stampedGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func stampedGetCall(ctx context.Context, subject Stamped, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Versioned:
//
// Versioned.Put/smoke: Put returns without a panic for derived inputs.
//
// Versioned.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Versioned.Put/cancel: Put returns an error for a cancelled context.
//
// Versioned.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Versioned.Get/smoke: Get returns without a panic for derived inputs.
//
// Versioned.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Versioned.Get/cancel: Get returns an error for a cancelled context.
//
// Versioned.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Versioned.Get/zero-on-error: Get returns the zero value beside an error.
//
// Versioned.Get/zero-on-error#generated: Get returns the zero value beside
// an error.

// The derived checks of Versioned that witness their laws at the shallow
// end:
//
// Versioned.Put/smoke
//
// Versioned.Get/smoke

// The refused claims of Versioned:
//
// Versioned.Get/writes-follow-reads: the partner write of Get returns no
// value of the type that Get returns. The remedy is to return the written
// value from the partner write of Get.
//
// Versioned.Get/writes-follow-their-reads: the partner write of Get
// returns no value of the type that Get returns. The remedy is to return
// the written value from the partner write of Get.

// The unclaimed classifications of Versioned:
//
// The shape:writer of Versioned.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Versioned writes.
type (
	VersionedHandle = conform.Handle[Versioned]
	VersionedCheck  = conform.Check[Versioned]
	VersionedOption = conform.Option[Versioned]
)

// VersionedChecks contains the name of each callable of Versioned, with
// the identities of its checks.
var VersionedChecks = struct {
	Put versionedPut
	Get versionedGet
}{
	Put: versionedPut{"Versioned.Put"},
	Get: versionedGet{"Versioned.Get"},
}

// versionedPut is the name of Versioned.Put, with the identities of its
// checks.
type versionedPut struct {
	conform.Callable[Versioned]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c versionedPut) Smoke() conform.ID[Versioned] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c versionedPut) SmokeGenerated() conform.ID[Versioned] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c versionedPut) Cancel() conform.ID[Versioned] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c versionedPut) CancelGenerated() conform.ID[Versioned] {
	return c.Check("cancel#generated")
}

// versionedGet is the name of Versioned.Get, with the identities of its
// checks.
type versionedGet struct {
	conform.Callable[Versioned]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c versionedGet) Smoke() conform.ID[Versioned] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c versionedGet) SmokeGenerated() conform.ID[Versioned] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c versionedGet) Cancel() conform.ID[Versioned] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c versionedGet) CancelGenerated() conform.ID[Versioned] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c versionedGet) ZeroOnError() conform.ID[Versioned] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c versionedGet) ZeroOnErrorGenerated() conform.ID[Versioned] {
	return c.Check("zero-on-error#generated")
}

// VersionedSuite returns the checks of Versioned, and the calls and the
// budgets of its callables.
func VersionedSuite() conform.Suite[Versioned] {
	return conform.Suite[Versioned]{
		Declaration: "Versioned",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Versioned]{
			{
				ID:          VersionedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VersionedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: versionedPutPanics,
			},
			{
				ID:          VersionedChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: versionedPutPanics,
			},
			{
				ID:          VersionedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VersionedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: versionedPutIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: versionedPutIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VersionedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: versionedGetPanics,
			},
			{
				ID:          VersionedChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: versionedGetPanics,
			},
			{
				ID:          VersionedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VersionedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: versionedGetIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: versionedGetIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VersionedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: versionedGetValueBeside,
			},
			{
				ID:          VersionedChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: versionedGetValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.WritesFollowReads,
				Callable: "Versioned.Get",
				Gate:     spec.WriteReturnsTheValue,
				Why:      "the partner write of Get returns no value of the type that Get returns",
				Remedy:   "return the written value from the partner write of Get",
			},
			{
				Claim:    spec.WritesFollowTheirReads,
				Callable: "Versioned.Get",
				Gate:     spec.WriteReturnsTheValue,
				Why:      "the partner write of Get returns no value of the type that Get returns",
				Remedy:   "return the written value from the partner write of Get",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Versioned.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Versioned]{
			{
				Calls:    versionedPutCalls,
				Over:     newVersionedOver,
				Callable: VersionedChecks.Put.Callable,
			},
			{
				Calls:    versionedGetCalls,
				Over:     newVersionedOver,
				Callable: VersionedChecks.Get.Callable,
			},
		},
	}
}

// RunVersionedSuite applies the checks of Versioned, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunVersionedSuite[T Versioned](t *testing.T, newSubject func() T, opts ...VersionedOption) {
	t.Helper()
	RunVersionedSuiteFromHandle(t, func(*VersionedHandle) T { return newSubject() }, opts...)
}

// RunVersionedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunVersionedSuiteFromHandle[T Versioned](
	t *testing.T,
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	t.Helper()
	conform.Run(t, VersionedSuite(), versionedSubjects(newSubject), opts...)
}

// FuzzVersionedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzVersionedSuite[T Versioned](
	f *testing.F,
	callable conform.Callable[Versioned],
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	f.Helper()
	conform.Fuzz(f, VersionedSuite(), callable, versionedSubjects(newSubject), opts...)
}

// RunVersionedBudgets applies the allocation ceilings of Versioned and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunVersionedBudgets[T Versioned](
	t *testing.T,
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	t.Helper()
	conform.RunBudgets(t, VersionedSuite(), versionedSubjects(newSubject), opts...)
}

// BenchVersionedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchVersionedSuite[T Versioned](
	b conform.B,
	callable conform.Callable[Versioned],
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	b.Helper()
	conform.Bench(b, VersionedSuite(), callable, versionedSubjects(newSubject), opts...)
}

// versionedSubjects returns the factory of the subjects of Versioned that
// calls newSubject with the handle of each check.
func versionedSubjects[T Versioned](
	newSubject func(h *VersionedHandle) T,
) func(h *VersionedHandle) Versioned {
	return func(h *VersionedHandle) Versioned { return newSubject(h) }
}

// versionedPutCalls generates the calls of Versioned.Put with the inputs
// that each case draws.
var versionedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Versioned) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Versioned) error {
		return subject.Put(ctx, it)
	}
})

// versionedGetCalls generates the calls of Versioned.Get with the inputs
// that each case draws.
var versionedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Versioned) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Versioned) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// versionedOverEmbedded is Versioned under the name of the field that
// versionedOver embeds.
type versionedOverEmbedded = Versioned

// versionedOver is a subject of Versioned whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type versionedOver struct {
	versionedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o versionedOver) Put(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o versionedOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newVersionedOver returns the planted defect of a ceiling of Versioned,
// whose callables call spend.
func newVersionedOver(_ *VersionedHandle, spend func()) Versioned {
	return versionedOver{spend: spend}
}

// versionedPutPanicsDefect is the planted defect of a Versioned whose Put
// panics.
type versionedPutPanicsDefect struct {
	versionedOverEmbedded
}

// Put is the override of the planted defect.
func (*versionedPutPanicsDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Versioned whose Put panics")
}

// versionedPutPanics returns the planted defect of a Versioned whose Put
// panics over a fresh subject.
func versionedPutPanics(h *VersionedHandle) Versioned {
	return &versionedPutPanicsDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedPutIgnoresCancelDefect is the planted defect of a Versioned
// whose Put ignores the context.
type versionedPutIgnoresCancelDefect struct {
	versionedOverEmbedded
}

// Put is the override of the planted defect.
func (d *versionedPutIgnoresCancelDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.versionedOverEmbedded.Put(p0, p1)
}

// versionedPutIgnoresCancel returns the planted defect of a Versioned
// whose Put ignores the context over a fresh subject.
func versionedPutIgnoresCancel(h *VersionedHandle) Versioned {
	return &versionedPutIgnoresCancelDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedGetPanicsDefect is the planted defect of a Versioned whose Get
// panics.
type versionedGetPanicsDefect struct {
	versionedOverEmbedded
}

// Get is the override of the planted defect.
func (*versionedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Versioned whose Get panics")
}

// versionedGetPanics returns the planted defect of a Versioned whose Get
// panics over a fresh subject.
func versionedGetPanics(h *VersionedHandle) Versioned {
	return &versionedGetPanicsDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedGetIgnoresCancelDefect is the planted defect of a Versioned
// whose Get ignores the context.
type versionedGetIgnoresCancelDefect struct {
	versionedOverEmbedded
}

// Get is the override of the planted defect.
func (d *versionedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.versionedOverEmbedded.Get(p0, p1)
}

// versionedGetIgnoresCancel returns the planted defect of a Versioned
// whose Get ignores the context over a fresh subject.
func versionedGetIgnoresCancel(h *VersionedHandle) Versioned {
	return &versionedGetIgnoresCancelDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedGetValueBesideDefect is the planted defect of a Versioned whose
// Get returns a value beside an error.
type versionedGetValueBesideDefect struct {
	versionedOverEmbedded
}

// Get is the override of the planted defect.
func (d *versionedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.versionedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// versionedGetValueBeside returns the planted defect of a Versioned whose
// Get returns a value beside an error over a fresh subject.
func versionedGetValueBeside(h *VersionedHandle) Versioned {
	return &versionedGetValueBesideDefect{versionedOverEmbedded: h.Fresh()}
}

// The derived checks of Accounts:
//
// Accounts.Transfer/smoke: Transfer returns without a panic for derived
// inputs.
//
// Accounts.Transfer/smoke#generated: Transfer returns without a panic for
// derived inputs.
//
// Accounts.Transfer/cancel: Transfer returns an error for a cancelled
// context.
//
// Accounts.Transfer/cancel#generated: Transfer returns an error for a
// cancelled context.
//
// Accounts.Transfer/conserves-the-quantity: each call of Transfer keeps
// the total that Total reads, over any sequence of calls.
//
// Accounts.Total/smoke: Total returns without a panic for derived inputs.
//
// Accounts.Total/cancel: Total returns an error for a cancelled context.
//
// Accounts.Total/zero-on-error: Total returns the zero value beside an
// error.
//
// Accounts.Rename/smoke: Rename returns without a panic for derived
// inputs.
//
// Accounts.Rename/smoke#generated: Rename returns without a panic for
// derived inputs.
//
// Accounts.Rename/cancel: Rename returns an error for a cancelled context.
//
// Accounts.Rename/cancel#generated: Rename returns an error for a
// cancelled context.
//
// Accounts.Names/smoke: Names returns without a panic for derived inputs.
//
// Accounts.Names/cancel: Names returns an error for a cancelled context.
//
// Accounts.Names/zero-on-error: Names returns the zero value beside an
// error.
//
// Accounts.Fold/smoke: Fold returns without a panic for derived inputs.
//
// Accounts.Fold/smoke#generated: Fold returns without a panic for derived
// inputs.
//
// Accounts.Fold/cancel: Fold returns an error for a cancelled context.
//
// Accounts.Fold/cancel#generated: Fold returns an error for a cancelled
// context.
//
// Accounts.Audit/smoke: Audit returns without a panic for derived inputs.
//
// Accounts.Audit/cancel: Audit returns an error for a cancelled context.

// The derived checks of Accounts that witness their laws at the shallow
// end:
//
// Accounts.Transfer/smoke
//
// Accounts.Total/smoke
//
// Accounts.Rename/smoke
//
// Accounts.Names/smoke
//
// Accounts.Fold/smoke
//
// Accounts.Audit/smoke

// The refused claims of Accounts:
//
// Accounts.Rename/conserves-the-quantity: the reader that observe names
// for Rename returns no number. The remedy is to name with observe a
// callable that returns the total that Rename moves as a number.
//
// Accounts.Fold/conserves-the-quantity: the reader that observe names for
// Fold returns no number. The remedy is to name with observe a callable
// that returns the total that Fold moves as a number.

// The unclaimed classifications of Accounts:
//
// The shape:writer of Accounts.Transfer has the stance silent. A write is
// an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Accounts.Rename has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:writer of Accounts.Fold has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Accounts.Audit has the stance silent. A lifecycle
// call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Accounts writes.
type (
	AccountsHandle = conform.Handle[Accounts]
	AccountsCheck  = conform.Check[Accounts]
	AccountsOption = conform.Option[Accounts]
)

// AccountsChecks contains the name of each callable of Accounts, with the
// identities of its checks.
var AccountsChecks = struct {
	Transfer accountsTransfer
	Total    accountsTotal
	Rename   accountsRename
	Names    accountsNames
	Fold     accountsFold
	Audit    accountsAudit
}{
	Transfer: accountsTransfer{"Accounts.Transfer"},
	Total:    accountsTotal{"Accounts.Total"},
	Rename:   accountsRename{"Accounts.Rename"},
	Names:    accountsNames{"Accounts.Names"},
	Fold:     accountsFold{"Accounts.Fold"},
	Audit:    accountsAudit{"Accounts.Audit"},
}

// accountsTransfer is the name of Accounts.Transfer, with the identities
// of its checks.
type accountsTransfer struct {
	conform.Callable[Accounts]
}

// Smoke returns the identity of the check that Transfer returns without a
// panic for derived inputs.
func (c accountsTransfer) Smoke() conform.ID[Accounts] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Transfer returns
// without a panic for derived inputs.
func (c accountsTransfer) SmokeGenerated() conform.ID[Accounts] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Transfer returns an error
// for a cancelled context.
func (c accountsTransfer) Cancel() conform.ID[Accounts] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Transfer returns
// an error for a cancelled context.
func (c accountsTransfer) CancelGenerated() conform.ID[Accounts] {
	return c.Check("cancel#generated")
}

// ConservesTheQuantity returns the identity of the check that each call of
// Transfer keeps the total that Total reads, over any sequence of calls.
func (c accountsTransfer) ConservesTheQuantity() conform.ID[Accounts] {
	return c.Check("conserves-the-quantity")
}

// accountsTotal is the name of Accounts.Total, with the identities of its
// checks.
type accountsTotal struct {
	conform.Callable[Accounts]
}

// Smoke returns the identity of the check that Total returns without a
// panic for derived inputs.
func (c accountsTotal) Smoke() conform.ID[Accounts] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Total returns an error for
// a cancelled context.
func (c accountsTotal) Cancel() conform.ID[Accounts] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Total returns the
// zero value beside an error.
func (c accountsTotal) ZeroOnError() conform.ID[Accounts] {
	return c.Check("zero-on-error")
}

// accountsRename is the name of Accounts.Rename, with the identities of
// its checks.
type accountsRename struct {
	conform.Callable[Accounts]
}

// Smoke returns the identity of the check that Rename returns without a
// panic for derived inputs.
func (c accountsRename) Smoke() conform.ID[Accounts] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Rename returns
// without a panic for derived inputs.
func (c accountsRename) SmokeGenerated() conform.ID[Accounts] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Rename returns an error
// for a cancelled context.
func (c accountsRename) Cancel() conform.ID[Accounts] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Rename returns an
// error for a cancelled context.
func (c accountsRename) CancelGenerated() conform.ID[Accounts] {
	return c.Check("cancel#generated")
}

// accountsNames is the name of Accounts.Names, with the identities of its
// checks.
type accountsNames struct {
	conform.Callable[Accounts]
}

// Smoke returns the identity of the check that Names returns without a
// panic for derived inputs.
func (c accountsNames) Smoke() conform.ID[Accounts] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Names returns an error for
// a cancelled context.
func (c accountsNames) Cancel() conform.ID[Accounts] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Names returns the
// zero value beside an error.
func (c accountsNames) ZeroOnError() conform.ID[Accounts] {
	return c.Check("zero-on-error")
}

// accountsFold is the name of Accounts.Fold, with the identities of its
// checks.
type accountsFold struct {
	conform.Callable[Accounts]
}

// Smoke returns the identity of the check that Fold returns without a
// panic for derived inputs.
func (c accountsFold) Smoke() conform.ID[Accounts] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Fold returns
// without a panic for derived inputs.
func (c accountsFold) SmokeGenerated() conform.ID[Accounts] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Fold returns an error for
// a cancelled context.
func (c accountsFold) Cancel() conform.ID[Accounts] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Fold returns an
// error for a cancelled context.
func (c accountsFold) CancelGenerated() conform.ID[Accounts] {
	return c.Check("cancel#generated")
}

// accountsAudit is the name of Accounts.Audit, with the identities of its
// checks.
type accountsAudit struct {
	conform.Callable[Accounts]
}

// Smoke returns the identity of the check that Audit returns without a
// panic for derived inputs.
func (c accountsAudit) Smoke() conform.ID[Accounts] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Audit returns an error for
// a cancelled context.
func (c accountsAudit) Cancel() conform.ID[Accounts] {
	return c.Check("cancel")
}

// AccountsSuite returns the checks of Accounts, and the calls and the
// budgets of its callables.
func AccountsSuite() conform.Suite[Accounts] {
	return conform.Suite[Accounts]{
		Declaration: "Accounts",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Accounts]{
			{
				ID:          AccountsChecks.Transfer.Smoke(),
				Statement:   "Transfer returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Transfer(ctx, Move{From: "From-a"})
					})
				},
				Defect: accountsTransferPanics,
			},
			{
				ID:          AccountsChecks.Transfer.SmokeGenerated(),
				Statement:   "Transfer returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountsHandle, c *prop.Case) {
					drawnM := c.Draw(prop.Of[Move](), "m")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Transfer(ctx, drawnM)
					})
				},
				Defect: accountsTransferPanics,
			},
			{
				ID:          AccountsChecks.Transfer.Cancel(),
				Statement:   "Transfer returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Transfer(ctx, Move{From: "From-a"})
					})
				},
				Defect: accountsTransferIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Transfer.CancelGenerated(),
				Statement:   "Transfer returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountsHandle, c *prop.Case) {
					drawnM := c.Draw(prop.Of[Move](), "m")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Transfer(ctx, drawnM)
					})
				},
				Defect: accountsTransferIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Transfer.ConservesTheQuantity(),
				Statement:   "each call of Transfer keeps the total that Total reads, over any sequence of calls",
				Claim:       spec.ConservesTheQuantity,
				Probe:       spec.MovesKeepTheTotal,
				Reason:      spec.NotConserved,
				Law:         spec.MovesConserveTheTotal,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *AccountsHandle, c *prop.Case) {
					probe.MovesKeepTheTotal(
						h,
						c,
						accountsActions,
						"Transfer",
						func(ctx context.Context) (int64, error) {
							return h.Subject().Total(ctx)
						},
					)
				},
				Defect: accountsTotalLeaksOneUnit,
			},
			{
				ID:          AccountsChecks.Total.Smoke(),
				Statement:   "Total returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Total(ctx)
						return err
					})
				},
				Defect: accountsTotalPanics,
			},
			{
				ID:          AccountsChecks.Total.Cancel(),
				Statement:   "Total returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Total(ctx)
						return err
					})
				},
				Defect: accountsTotalIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Total.ZeroOnError(),
				Statement:   "Total returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Total(ctx)
					}, int64(0))
				},
				Defect: accountsTotalValueBeside,
			},
			{
				ID:          AccountsChecks.Rename.Smoke(),
				Statement:   "Rename returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Rename(ctx, string("from-a"), string("to-a"))
					})
				},
				Defect: accountsRenamePanics,
			},
			{
				ID:          AccountsChecks.Rename.SmokeGenerated(),
				Statement:   "Rename returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountsHandle, c *prop.Case) {
					drawnFrom := c.Draw(prop.Of[string](), "from")
					drawnTo := c.Draw(prop.Of[string](), "to")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Rename(ctx, drawnFrom, drawnTo)
					})
				},
				Defect: accountsRenamePanics,
			},
			{
				ID:          AccountsChecks.Rename.Cancel(),
				Statement:   "Rename returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Rename(ctx, string("from-a"), string("to-a"))
					})
				},
				Defect: accountsRenameIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Rename.CancelGenerated(),
				Statement:   "Rename returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountsHandle, c *prop.Case) {
					drawnFrom := c.Draw(prop.Of[string](), "from")
					drawnTo := c.Draw(prop.Of[string](), "to")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Rename(ctx, drawnFrom, drawnTo)
					})
				},
				Defect: accountsRenameIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Names.Smoke(),
				Statement:   "Names returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Names(ctx)
						return err
					})
				},
				Defect: accountsNamesPanics,
			},
			{
				ID:          AccountsChecks.Names.Cancel(),
				Statement:   "Names returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Names(ctx)
						return err
					})
				},
				Defect: accountsNamesIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Names.ZeroOnError(),
				Statement:   "Names returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Names(ctx)
					}, []string(nil))
				},
				Defect: accountsNamesValueBeside,
			},
			{
				ID:          AccountsChecks.Fold.Smoke(),
				Statement:   "Fold returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Fold(ctx, string("to-a"))
					})
				},
				Defect: accountsFoldPanics,
			},
			{
				ID:          AccountsChecks.Fold.SmokeGenerated(),
				Statement:   "Fold returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountsHandle, c *prop.Case) {
					drawnTo := c.Draw(prop.Of[string](), "to")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Fold(ctx, drawnTo)
					})
				},
				Defect: accountsFoldPanics,
			},
			{
				ID:          AccountsChecks.Fold.Cancel(),
				Statement:   "Fold returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Fold(ctx, string("to-a"))
					})
				},
				Defect: accountsFoldIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Fold.CancelGenerated(),
				Statement:   "Fold returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountsHandle, c *prop.Case) {
					drawnTo := c.Draw(prop.Of[string](), "to")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Fold(ctx, drawnTo)
					})
				},
				Defect: accountsFoldIgnoresCancel,
			},
			{
				ID:          AccountsChecks.Audit.Smoke(),
				Statement:   "Audit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Audit(ctx)
					})
				},
				Defect: accountsAuditPanics,
			},
			{
				ID:          AccountsChecks.Audit.Cancel(),
				Statement:   "Audit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Audit(ctx)
					})
				},
				Defect: accountsAuditIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ConservesTheQuantity,
				Callable: "Accounts.Rename",
				Gate:     spec.NumericTotal,
				Why:      "the reader that observe names for Rename returns no number",
				Remedy:   "name with observe a callable that returns the total that Rename moves as a number",
			},
			{
				Claim:    spec.ConservesTheQuantity,
				Callable: "Accounts.Fold",
				Gate:     spec.NumericTotal,
				Why:      "the reader that observe names for Fold returns no number",
				Remedy:   "name with observe a callable that returns the total that Fold moves as a number",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Accounts.Transfer",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Accounts.Rename",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Accounts.Fold",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Accounts.Audit",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Accounts]{
			{
				Calls:    accountsTransferCalls,
				Over:     newAccountsOver,
				Callable: AccountsChecks.Transfer.Callable,
			},
			{
				Calls:    accountsTotalCalls,
				Over:     newAccountsOver,
				Callable: AccountsChecks.Total.Callable,
			},
			{
				Calls:    accountsRenameCalls,
				Over:     newAccountsOver,
				Callable: AccountsChecks.Rename.Callable,
			},
			{
				Calls:    accountsNamesCalls,
				Over:     newAccountsOver,
				Callable: AccountsChecks.Names.Callable,
			},
			{
				Calls:    accountsFoldCalls,
				Over:     newAccountsOver,
				Callable: AccountsChecks.Fold.Callable,
			},
			{
				Calls:    accountsAuditCalls,
				Over:     newAccountsOver,
				Callable: AccountsChecks.Audit.Callable,
			},
		},
	}
}

// RunAccountsSuite applies the checks of Accounts, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunAccountsSuite[T Accounts](t *testing.T, newSubject func() T, opts ...AccountsOption) {
	t.Helper()
	RunAccountsSuiteFromHandle(t, func(*AccountsHandle) T { return newSubject() }, opts...)
}

// RunAccountsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunAccountsSuiteFromHandle[T Accounts](
	t *testing.T,
	newSubject func(h *AccountsHandle) T,
	opts ...AccountsOption,
) {
	t.Helper()
	conform.Run(t, AccountsSuite(), accountsSubjects(newSubject), opts...)
}

// FuzzAccountsSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzAccountsSuite[T Accounts](
	f *testing.F,
	callable conform.Callable[Accounts],
	newSubject func(h *AccountsHandle) T,
	opts ...AccountsOption,
) {
	f.Helper()
	conform.Fuzz(f, AccountsSuite(), callable, accountsSubjects(newSubject), opts...)
}

// RunAccountsBudgets applies the allocation ceilings of Accounts and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunAccountsBudgets[T Accounts](
	t *testing.T,
	newSubject func(h *AccountsHandle) T,
	opts ...AccountsOption,
) {
	t.Helper()
	conform.RunBudgets(t, AccountsSuite(), accountsSubjects(newSubject), opts...)
}

// BenchAccountsSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchAccountsSuite[T Accounts](
	b conform.B,
	callable conform.Callable[Accounts],
	newSubject func(h *AccountsHandle) T,
	opts ...AccountsOption,
) {
	b.Helper()
	conform.Bench(b, AccountsSuite(), callable, accountsSubjects(newSubject), opts...)
}

// accountsSubjects returns the factory of the subjects of Accounts that
// calls newSubject with the handle of each check.
func accountsSubjects[T Accounts](
	newSubject func(h *AccountsHandle) T,
) func(h *AccountsHandle) Accounts {
	return func(h *AccountsHandle) Accounts { return newSubject(h) }
}

// accountsTransferCalls generates the calls of Accounts.Transfer with the
// inputs that each case draws.
var accountsTransferCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, Accounts) error {
		m := c.Draw(prop.Of[Move](), "m")
		return func(ctx context.Context, subject Accounts) error {
			return subject.Transfer(ctx, m)
		}
	},
)

// accountsTotalCalls generates the calls of Accounts.Total with the inputs
// that each case draws.
var accountsTotalCalls = prop.Composite(func(*prop.Case) func(context.Context, Accounts) error {
	return func(ctx context.Context, subject Accounts) error {
		_, err := subject.Total(ctx)
		return err
	}
})

// accountsRenameCalls generates the calls of Accounts.Rename with the
// inputs that each case draws.
var accountsRenameCalls = prop.Composite(func(c *prop.Case) func(context.Context, Accounts) error {
	from := c.Draw(prop.Of[string](), "from")
	to := c.Draw(prop.Of[string](), "to")
	return func(ctx context.Context, subject Accounts) error {
		return subject.Rename(ctx, from, to)
	}
})

// accountsNamesCalls generates the calls of Accounts.Names with the inputs
// that each case draws.
var accountsNamesCalls = prop.Composite(func(*prop.Case) func(context.Context, Accounts) error {
	return func(ctx context.Context, subject Accounts) error {
		_, err := subject.Names(ctx)
		return err
	}
})

// accountsFoldCalls generates the calls of Accounts.Fold with the inputs
// that each case draws.
var accountsFoldCalls = prop.Composite(func(c *prop.Case) func(context.Context, Accounts) error {
	to := c.Draw(prop.Of[string](), "to")
	return func(ctx context.Context, subject Accounts) error {
		return subject.Fold(ctx, to)
	}
})

// accountsAuditCalls generates the calls of Accounts.Audit with the inputs
// that each case draws.
var accountsAuditCalls = prop.Composite(func(*prop.Case) func(context.Context, Accounts) error {
	return func(ctx context.Context, subject Accounts) error {
		return subject.Audit(ctx)
	}
})

// accountsOverEmbedded is Accounts under the name of the field that
// accountsOver embeds.
type accountsOverEmbedded = Accounts

// accountsOver is a subject of Accounts whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type accountsOver struct {
	accountsOverEmbedded
	spend func()
}

// Transfer calls spend and returns its zero results.
func (o accountsOver) Transfer(context.Context, Move) (r0 error) {
	o.spend()
	return r0
}

// Total calls spend and returns its zero results.
func (o accountsOver) Total(context.Context) (r0 int64, r1 error) {
	o.spend()
	return r0, r1
}

// Rename calls spend and returns its zero results.
func (o accountsOver) Rename(context.Context, string, string) (r0 error) {
	o.spend()
	return r0
}

// Names calls spend and returns its zero results.
func (o accountsOver) Names(context.Context) (r0 []string, r1 error) {
	o.spend()
	return r0, r1
}

// Fold calls spend and returns its zero results.
func (o accountsOver) Fold(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Audit calls spend and returns its zero results.
func (o accountsOver) Audit(context.Context) (r0 error) {
	o.spend()
	return r0
}

// newAccountsOver returns the planted defect of a ceiling of Accounts,
// whose callables call spend.
func newAccountsOver(_ *AccountsHandle, spend func()) Accounts {
	return accountsOver{spend: spend}
}

// accountsTransferPanicsDefect is the planted defect of a Accounts whose
// Transfer panics.
type accountsTransferPanicsDefect struct {
	accountsOverEmbedded
}

// Transfer is the override of the planted defect.
func (*accountsTransferPanicsDefect) Transfer(p0 context.Context, p1 Move) (r0 error) {
	panic("the planted defect of a Accounts whose Transfer panics")
}

// accountsTransferPanics returns the planted defect of a Accounts whose
// Transfer panics over a fresh subject.
func accountsTransferPanics(h *AccountsHandle) Accounts {
	return &accountsTransferPanicsDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsTransferIgnoresCancelDefect is the planted defect of a Accounts
// whose Transfer ignores the context.
type accountsTransferIgnoresCancelDefect struct {
	accountsOverEmbedded
}

// Transfer is the override of the planted defect.
func (d *accountsTransferIgnoresCancelDefect) Transfer(p0 context.Context, p1 Move) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.accountsOverEmbedded.Transfer(p0, p1)
}

// accountsTransferIgnoresCancel returns the planted defect of a Accounts
// whose Transfer ignores the context over a fresh subject.
func accountsTransferIgnoresCancel(h *AccountsHandle) Accounts {
	return &accountsTransferIgnoresCancelDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsTotalLeaksOneUnitDefect is the planted defect of a Accounts
// whose Total reads one unit more on each call.
type accountsTotalLeaksOneUnitDefect struct {
	accountsOverEmbedded
	calls int
}

// Total is the override of the planted defect.
func (d *accountsTotalLeaksOneUnitDefect) Total(p0 context.Context) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.accountsOverEmbedded.Total(p0)
	d.calls++
	r0 += int64(d.calls)
	return r0, r1
}

// accountsTotalLeaksOneUnit returns the planted defect of a Accounts whose
// Total reads one unit more on each call over a fresh subject.
func accountsTotalLeaksOneUnit(h *AccountsHandle) Accounts {
	return &accountsTotalLeaksOneUnitDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsTotalPanicsDefect is the planted defect of a Accounts whose
// Total panics.
type accountsTotalPanicsDefect struct {
	accountsOverEmbedded
}

// Total is the override of the planted defect.
func (*accountsTotalPanicsDefect) Total(p0 context.Context) (r0 int64, r1 error) {
	panic("the planted defect of a Accounts whose Total panics")
}

// accountsTotalPanics returns the planted defect of a Accounts whose Total
// panics over a fresh subject.
func accountsTotalPanics(h *AccountsHandle) Accounts {
	return &accountsTotalPanicsDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsTotalIgnoresCancelDefect is the planted defect of a Accounts
// whose Total ignores the context.
type accountsTotalIgnoresCancelDefect struct {
	accountsOverEmbedded
}

// Total is the override of the planted defect.
func (d *accountsTotalIgnoresCancelDefect) Total(p0 context.Context) (r0 int64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.accountsOverEmbedded.Total(p0)
}

// accountsTotalIgnoresCancel returns the planted defect of a Accounts
// whose Total ignores the context over a fresh subject.
func accountsTotalIgnoresCancel(h *AccountsHandle) Accounts {
	return &accountsTotalIgnoresCancelDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsTotalValueBesideDefect is the planted defect of a Accounts whose
// Total returns a value beside an error.
type accountsTotalValueBesideDefect struct {
	accountsOverEmbedded
}

// Total is the override of the planted defect.
func (d *accountsTotalValueBesideDefect) Total(p0 context.Context) (r0 int64, r1 error) {
	r0, r1 = d.accountsOverEmbedded.Total(p0)
	if r1 != nil {
		r0 = int64(7)
	}
	return r0, r1
}

// accountsTotalValueBeside returns the planted defect of a Accounts whose
// Total returns a value beside an error over a fresh subject.
func accountsTotalValueBeside(h *AccountsHandle) Accounts {
	return &accountsTotalValueBesideDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsRenamePanicsDefect is the planted defect of a Accounts whose
// Rename panics.
type accountsRenamePanicsDefect struct {
	accountsOverEmbedded
}

// Rename is the override of the planted defect.
func (*accountsRenamePanicsDefect) Rename(p0 context.Context, p1, p2 string) (r0 error) {
	panic("the planted defect of a Accounts whose Rename panics")
}

// accountsRenamePanics returns the planted defect of a Accounts whose
// Rename panics over a fresh subject.
func accountsRenamePanics(h *AccountsHandle) Accounts {
	return &accountsRenamePanicsDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsRenameIgnoresCancelDefect is the planted defect of a Accounts
// whose Rename ignores the context.
type accountsRenameIgnoresCancelDefect struct {
	accountsOverEmbedded
}

// Rename is the override of the planted defect.
func (d *accountsRenameIgnoresCancelDefect) Rename(p0 context.Context, p1, p2 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.accountsOverEmbedded.Rename(p0, p1, p2)
}

// accountsRenameIgnoresCancel returns the planted defect of a Accounts
// whose Rename ignores the context over a fresh subject.
func accountsRenameIgnoresCancel(h *AccountsHandle) Accounts {
	return &accountsRenameIgnoresCancelDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsNamesPanicsDefect is the planted defect of a Accounts whose
// Names panics.
type accountsNamesPanicsDefect struct {
	accountsOverEmbedded
}

// Names is the override of the planted defect.
func (*accountsNamesPanicsDefect) Names(p0 context.Context) (r0 []string, r1 error) {
	panic("the planted defect of a Accounts whose Names panics")
}

// accountsNamesPanics returns the planted defect of a Accounts whose Names
// panics over a fresh subject.
func accountsNamesPanics(h *AccountsHandle) Accounts {
	return &accountsNamesPanicsDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsNamesIgnoresCancelDefect is the planted defect of a Accounts
// whose Names ignores the context.
type accountsNamesIgnoresCancelDefect struct {
	accountsOverEmbedded
}

// Names is the override of the planted defect.
func (d *accountsNamesIgnoresCancelDefect) Names(p0 context.Context) (r0 []string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.accountsOverEmbedded.Names(p0)
}

// accountsNamesIgnoresCancel returns the planted defect of a Accounts
// whose Names ignores the context over a fresh subject.
func accountsNamesIgnoresCancel(h *AccountsHandle) Accounts {
	return &accountsNamesIgnoresCancelDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsNamesValueBesideDefect is the planted defect of a Accounts whose
// Names returns a value beside an error.
type accountsNamesValueBesideDefect struct {
	accountsOverEmbedded
}

// Names is the override of the planted defect.
func (d *accountsNamesValueBesideDefect) Names(p0 context.Context) (r0 []string, r1 error) {
	r0, r1 = d.accountsOverEmbedded.Names(p0)
	if r1 != nil {
		r0 = []string{"sample-b"}
	}
	return r0, r1
}

// accountsNamesValueBeside returns the planted defect of a Accounts whose
// Names returns a value beside an error over a fresh subject.
func accountsNamesValueBeside(h *AccountsHandle) Accounts {
	return &accountsNamesValueBesideDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsFoldPanicsDefect is the planted defect of a Accounts whose Fold
// panics.
type accountsFoldPanicsDefect struct {
	accountsOverEmbedded
}

// Fold is the override of the planted defect.
func (*accountsFoldPanicsDefect) Fold(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Accounts whose Fold panics")
}

// accountsFoldPanics returns the planted defect of a Accounts whose Fold
// panics over a fresh subject.
func accountsFoldPanics(h *AccountsHandle) Accounts {
	return &accountsFoldPanicsDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsFoldIgnoresCancelDefect is the planted defect of a Accounts
// whose Fold ignores the context.
type accountsFoldIgnoresCancelDefect struct {
	accountsOverEmbedded
}

// Fold is the override of the planted defect.
func (d *accountsFoldIgnoresCancelDefect) Fold(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.accountsOverEmbedded.Fold(p0, p1)
}

// accountsFoldIgnoresCancel returns the planted defect of a Accounts whose
// Fold ignores the context over a fresh subject.
func accountsFoldIgnoresCancel(h *AccountsHandle) Accounts {
	return &accountsFoldIgnoresCancelDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsAuditPanicsDefect is the planted defect of a Accounts whose
// Audit panics.
type accountsAuditPanicsDefect struct {
	accountsOverEmbedded
}

// Audit is the override of the planted defect.
func (*accountsAuditPanicsDefect) Audit(p0 context.Context) (r0 error) {
	panic("the planted defect of a Accounts whose Audit panics")
}

// accountsAuditPanics returns the planted defect of a Accounts whose Audit
// panics over a fresh subject.
func accountsAuditPanics(h *AccountsHandle) Accounts {
	return &accountsAuditPanicsDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsAuditIgnoresCancelDefect is the planted defect of a Accounts
// whose Audit ignores the context.
type accountsAuditIgnoresCancelDefect struct {
	accountsOverEmbedded
}

// Audit is the override of the planted defect.
func (d *accountsAuditIgnoresCancelDefect) Audit(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.accountsOverEmbedded.Audit(p0)
}

// accountsAuditIgnoresCancel returns the planted defect of a Accounts
// whose Audit ignores the context over a fresh subject.
func accountsAuditIgnoresCancel(h *AccountsHandle) Accounts {
	return &accountsAuditIgnoresCancelDefect{accountsOverEmbedded: h.Fresh()}
}

// accountsActions returns the actions of the generated sequences of
// Accounts over the subject that subject returns. The actions of a case
// take their keys from the pools that it draws.
func accountsActions(c *prop.Case, subject func() Accounts) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Transfer",
			Input: func(c *prop.Case, _ struct{}) any {
				m := c.Draw(prop.Of[Move](), "transfer.m")
				return []any{m}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Transfer", args)
				values, err := accountsTransferCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Total",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Total", args)
				values, err := accountsTotalCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Rename",
			Input: func(c *prop.Case, _ struct{}) any {
				from := c.Draw(prop.SampledFrom(keysString...), "rename.from")
				to := c.Draw(prop.Of[string](), "rename.to")
				return []any{from, to}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Rename", args, args[0])
				values, err := accountsRenameCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Names",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Names", args)
				values, err := accountsNamesCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Fold",
			Input: func(c *prop.Case, _ struct{}) any {
				to := c.Draw(prop.Of[string](), "fold.to")
				return []any{to}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Fold", args)
				values, err := accountsFoldCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Audit",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Audit", args)
				values, err := accountsAuditCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// accountsTransferCall calls Transfer on subject with the arguments args,
// and returns the values of the call beside its error.
func accountsTransferCall(ctx context.Context, subject Accounts, args []any) ([]any, error) {
	a0, _ := args[0].(Move)
	return nil, subject.Transfer(ctx, a0)
}

// accountsTotalCall calls Total on subject with the arguments args, and
// returns the values of the call beside its error.
func accountsTotalCall(ctx context.Context, subject Accounts, args []any) ([]any, error) {
	r0, err := subject.Total(ctx)
	return []any{r0}, err
}

// accountsRenameCall calls Rename on subject with the arguments args, and
// returns the values of the call beside its error.
func accountsRenameCall(ctx context.Context, subject Accounts, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	return nil, subject.Rename(ctx, a0, a1)
}

// accountsNamesCall calls Names on subject with the arguments args, and
// returns the values of the call beside its error.
func accountsNamesCall(ctx context.Context, subject Accounts, args []any) ([]any, error) {
	r0, err := subject.Names(ctx)
	return []any{r0}, err
}

// accountsFoldCall calls Fold on subject with the arguments args, and
// returns the values of the call beside its error.
func accountsFoldCall(ctx context.Context, subject Accounts, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Fold(ctx, a0)
}

// accountsAuditCall calls Audit on subject with the arguments args, and
// returns the values of the call beside its error.
func accountsAuditCall(ctx context.Context, subject Accounts, args []any) ([]any, error) {
	return nil, subject.Audit(ctx)
}

// The derived checks of Journal:
//
// Journal.Record/smoke: Record returns without a panic for derived inputs.
//
// Journal.Record/smoke#generated: Record returns without a panic for
// derived inputs.
//
// Journal.Record/cancel: Record returns an error for a cancelled context.
//
// Journal.Record/cancel#generated: Record returns an error for a cancelled
// context.
//
// Journal.Record/failure-is-reported: Record returns an error when a call
// of its dependency fails.
//
// Journal.Record/failure-is-reported#generated: Record returns an error
// when a call of its dependency fails.
//
// Journal.Record/no-partial-write: Entry finds no effect of a Record that
// failed with its dependency.
//
// Journal.Record/no-partial-write#generated: Entry finds no effect of a
// Record that failed with its dependency.
//
// Journal.Record/no-trace-of-a-failed-call: Entry finds no effect of a
// Record that failed with its dependencies, over any sequence of calls.
//
// Journal.Entry/smoke: Entry returns without a panic for derived inputs.
//
// Journal.Entry/smoke#generated: Entry returns without a panic for derived
// inputs.
//
// Journal.Entry/cancel: Entry returns an error for a cancelled context.
//
// Journal.Entry/cancel#generated: Entry returns an error for a cancelled
// context.
//
// Journal.Entry/failure-is-reported: Entry returns an error when a call of
// its dependency fails.
//
// Journal.Entry/failure-is-reported#generated: Entry returns an error when
// a call of its dependency fails.
//
// Journal.Entry/zero-on-error: Entry returns the zero value beside an
// error.
//
// Journal.Entry/zero-on-error#generated: Entry returns the zero value
// beside an error.

// The derived checks of Journal that witness their laws at the shallow
// end:
//
// Journal.Record/smoke
//
// Journal.Record/no-partial-write
//
// Journal.Record/no-partial-write#generated
//
// Journal.Entry/smoke

// The unclaimed classifications of Journal:
//
// The shape:writer of Journal.Record has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Journal writes.
type (
	JournalHandle = conform.Handle[Journal]
	JournalCheck  = conform.Check[Journal]
	JournalOption = conform.Option[Journal]
)

// JournalChecks contains the name of each callable of Journal, with the
// identities of its checks.
var JournalChecks = struct {
	Record journalRecord
	Entry  journalEntry
}{
	Record: journalRecord{"Journal.Record"},
	Entry:  journalEntry{"Journal.Entry"},
}

// journalRecord is the name of Journal.Record, with the identities of its
// checks.
type journalRecord struct {
	conform.Callable[Journal]
}

// Smoke returns the identity of the check that Record returns without a
// panic for derived inputs.
func (c journalRecord) Smoke() conform.ID[Journal] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Record returns
// without a panic for derived inputs.
func (c journalRecord) SmokeGenerated() conform.ID[Journal] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Record returns an error
// for a cancelled context.
func (c journalRecord) Cancel() conform.ID[Journal] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Record returns an
// error for a cancelled context.
func (c journalRecord) CancelGenerated() conform.ID[Journal] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Record returns
// an error when a call of its dependency fails.
func (c journalRecord) FailureIsReported() conform.ID[Journal] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Record
// returns an error when a call of its dependency fails.
func (c journalRecord) FailureIsReportedGenerated() conform.ID[Journal] {
	return c.Check("failure-is-reported#generated")
}

// NoPartialWrite returns the identity of the check that Entry finds no
// effect of a Record that failed with its dependency.
func (c journalRecord) NoPartialWrite() conform.ID[Journal] {
	return c.Check("no-partial-write")
}

// NoPartialWriteGenerated returns the identity of the check that Entry
// finds no effect of a Record that failed with its dependency.
func (c journalRecord) NoPartialWriteGenerated() conform.ID[Journal] {
	return c.Check("no-partial-write#generated")
}

// NoTraceOfAFailedCall returns the identity of the check that Entry finds
// no effect of a Record that failed with its dependencies, over any
// sequence of calls.
func (c journalRecord) NoTraceOfAFailedCall() conform.ID[Journal] {
	return c.Check("no-trace-of-a-failed-call")
}

// journalEntry is the name of Journal.Entry, with the identities of its
// checks.
type journalEntry struct {
	conform.Callable[Journal]
}

// Smoke returns the identity of the check that Entry returns without a
// panic for derived inputs.
func (c journalEntry) Smoke() conform.ID[Journal] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Entry returns
// without a panic for derived inputs.
func (c journalEntry) SmokeGenerated() conform.ID[Journal] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Entry returns an error for
// a cancelled context.
func (c journalEntry) Cancel() conform.ID[Journal] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Entry returns an
// error for a cancelled context.
func (c journalEntry) CancelGenerated() conform.ID[Journal] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Entry returns
// an error when a call of its dependency fails.
func (c journalEntry) FailureIsReported() conform.ID[Journal] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Entry
// returns an error when a call of its dependency fails.
func (c journalEntry) FailureIsReportedGenerated() conform.ID[Journal] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Entry returns the
// zero value beside an error.
func (c journalEntry) ZeroOnError() conform.ID[Journal] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Entry
// returns the zero value beside an error.
func (c journalEntry) ZeroOnErrorGenerated() conform.ID[Journal] {
	return c.Check("zero-on-error#generated")
}

// JournalSuite returns the checks of Journal, and the calls and the
// budgets of its callables.
func JournalSuite() conform.Suite[Journal] {
	return conform.Suite[Journal]{
		Declaration: "Journal",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Journal]{
			{
				ID:          JournalChecks.Record.Smoke(),
				Statement:   "Record returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *JournalHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: journalRecordPanics,
			},
			{
				ID:          JournalChecks.Record.SmokeGenerated(),
				Statement:   "Record returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, drawnIt)
					})
				},
				Defect: journalRecordPanics,
			},
			{
				ID:          JournalChecks.Record.Cancel(),
				Statement:   "Record returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *JournalHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: journalRecordIgnoresCancel,
			},
			{
				ID:          JournalChecks.Record.CancelGenerated(),
				Statement:   "Record returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, drawnIt)
					})
				},
				Defect: journalRecordIgnoresCancel,
			},
			{
				ID:        JournalChecks.Record.FailureIsReported(),
				Statement: "Record returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *JournalHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: journalRecordSwallowsFailure,
			},
			{
				ID:        JournalChecks.Record.FailureIsReportedGenerated(),
				Statement: "Record returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, drawnIt)
					})
				},
				Defect: journalRecordSwallowsFailure,
			},
			{
				ID:        JournalChecks.Record.NoPartialWrite(),
				Statement: "Entry finds no effect of a Record that failed with its dependency",
				Claim:     spec.NoPartialWrite,
				Probe:     spec.FailsWithoutEffect,
				Reason:    spec.PartialWrite,
				Law:       spec.FailedCallLeavesNoEffect,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Shallow: true,
				Apply: func(h *JournalHandle) {
					probe.FailsWithoutEffect(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, Item{Key: "Key-a"})
					}, func(ctx context.Context) (Item, error) {
						return h.Subject().Entry(ctx, (Item{Key: "Key-a"}).Key)
					})
				},
				Defect: journalEntryHalfApplied,
			},
			{
				ID:        JournalChecks.Record.NoPartialWriteGenerated(),
				Statement: "Entry finds no effect of a Record that failed with its dependency",
				Claim:     spec.NoPartialWrite,
				Probe:     spec.FailsWithoutEffect,
				Reason:    spec.PartialWrite,
				Law:       spec.FailedCallLeavesNoEffect,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Shallow: true,
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.FailsWithoutEffect(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, drawnIt)
					}, func(ctx context.Context) (Item, error) {
						return h.Subject().Entry(ctx, drawnIt.Key)
					})
				},
				Defect: journalEntryHalfApplied,
			},
			{
				ID:        JournalChecks.Record.NoTraceOfAFailedCall(),
				Statement: "Entry finds no effect of a Record that failed with its dependencies, over any sequence of calls",
				Claim:     spec.NoTraceOfAFailedCall,
				Probe:     spec.FailedCallsLeaveNoTrace,
				Reason:    spec.PartialWrite,
				Law:       spec.FailedCallLeavesNoEffect,
				Machineries: []spec.MachineryID{
					spec.GeneratedSequences,
					spec.CallThatFailsOnDemand,
				},
				Entails: []conform.ID[Journal]{JournalChecks.Record.NoPartialWrite()},
				Property: func(h *JournalHandle, c *prop.Case) {
					probe.FailedCallsLeaveNoTrace(
						h,
						c,
						journalActions,
						"Record",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := journalEntryCall(ctx, h.Subject(), []any{w.Key})
							return conform.OutputOf(err, values...)
						},
					)
				},
				Defect: journalEntryHalfAppliedAfterThree,
			},
			{
				ID:          JournalChecks.Entry.Smoke(),
				Statement:   "Entry returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *JournalHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Entry(ctx, string("key-a"))
						return err
					})
				},
				Defect: journalEntryPanics,
			},
			{
				ID:          JournalChecks.Entry.SmokeGenerated(),
				Statement:   "Entry returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Entry(ctx, drawnKey)
						return err
					})
				},
				Defect: journalEntryPanics,
			},
			{
				ID:          JournalChecks.Entry.Cancel(),
				Statement:   "Entry returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *JournalHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Entry(ctx, string("key-a"))
						return err
					})
				},
				Defect: journalEntryIgnoresCancel,
			},
			{
				ID:          JournalChecks.Entry.CancelGenerated(),
				Statement:   "Entry returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Entry(ctx, drawnKey)
						return err
					})
				},
				Defect: journalEntryIgnoresCancel,
			},
			{
				ID:        JournalChecks.Entry.FailureIsReported(),
				Statement: "Entry returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *JournalHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Entry(ctx, string("key-a"))
						return err
					})
				},
				Defect: journalEntrySwallowsFailure,
			},
			{
				ID:        JournalChecks.Entry.FailureIsReportedGenerated(),
				Statement: "Entry returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Entry(ctx, drawnKey)
						return err
					})
				},
				Defect: journalEntrySwallowsFailure,
			},
			{
				ID:          JournalChecks.Entry.ZeroOnError(),
				Statement:   "Entry returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *JournalHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Entry(ctx, string("key-a"))
					}, Item{})
				},
				Defect: journalEntryValueBeside,
			},
			{
				ID:          JournalChecks.Entry.ZeroOnErrorGenerated(),
				Statement:   "Entry returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Entry(ctx, drawnKey)
					}, Item{})
				},
				Defect: journalEntryValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Journal.Record",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Journal]{
			{
				Calls:    journalRecordCalls,
				Over:     newJournalOver,
				Callable: JournalChecks.Record.Callable,
			},
			{
				Calls:    journalEntryCalls,
				Over:     newJournalOver,
				Callable: JournalChecks.Entry.Callable,
			},
		},
	}
}

// RunJournalSuite applies the checks of Journal, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunJournalSuite[T Journal](t *testing.T, newSubject func() T, opts ...JournalOption) {
	t.Helper()
	RunJournalSuiteFromHandle(t, func(*JournalHandle) T { return newSubject() }, opts...)
}

// RunJournalSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunJournalSuiteFromHandle[T Journal](
	t *testing.T,
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	t.Helper()
	conform.Run(t, JournalSuite(), journalSubjects(newSubject), opts...)
}

// FuzzJournalSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzJournalSuite[T Journal](
	f *testing.F,
	callable conform.Callable[Journal],
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	f.Helper()
	conform.Fuzz(f, JournalSuite(), callable, journalSubjects(newSubject), opts...)
}

// RunJournalBudgets applies the allocation ceilings of Journal and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunJournalBudgets[T Journal](
	t *testing.T,
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	t.Helper()
	conform.RunBudgets(t, JournalSuite(), journalSubjects(newSubject), opts...)
}

// BenchJournalSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchJournalSuite[T Journal](
	b conform.B,
	callable conform.Callable[Journal],
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	b.Helper()
	conform.Bench(b, JournalSuite(), callable, journalSubjects(newSubject), opts...)
}

// journalSubjects returns the factory of the subjects of Journal that
// calls newSubject with the handle of each check.
func journalSubjects[T Journal](
	newSubject func(h *JournalHandle) T,
) func(h *JournalHandle) Journal {
	return func(h *JournalHandle) Journal { return newSubject(h) }
}

// journalRecordCalls generates the calls of Journal.Record with the inputs
// that each case draws.
var journalRecordCalls = prop.Composite(func(c *prop.Case) func(context.Context, Journal) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Journal) error {
		return subject.Record(ctx, it)
	}
})

// journalEntryCalls generates the calls of Journal.Entry with the inputs
// that each case draws.
var journalEntryCalls = prop.Composite(func(c *prop.Case) func(context.Context, Journal) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Journal) error {
		_, err := subject.Entry(ctx, key)
		return err
	}
})

// journalOverEmbedded is Journal under the name of the field that
// journalOver embeds.
type journalOverEmbedded = Journal

// journalOver is a subject of Journal whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type journalOver struct {
	journalOverEmbedded
	spend func()
}

// Record calls spend and returns its zero results.
func (o journalOver) Record(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Entry calls spend and returns its zero results.
func (o journalOver) Entry(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newJournalOver returns the planted defect of a ceiling of Journal, whose
// callables call spend.
func newJournalOver(_ *JournalHandle, spend func()) Journal {
	return journalOver{spend: spend}
}

// journalRecordPanicsDefect is the planted defect of a Journal whose
// Record panics.
type journalRecordPanicsDefect struct {
	journalOverEmbedded
}

// Record is the override of the planted defect.
func (*journalRecordPanicsDefect) Record(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Journal whose Record panics")
}

// journalRecordPanics returns the planted defect of a Journal whose Record
// panics over a fresh subject.
func journalRecordPanics(h *JournalHandle) Journal {
	return &journalRecordPanicsDefect{journalOverEmbedded: h.Fresh()}
}

// journalRecordIgnoresCancelDefect is the planted defect of a Journal
// whose Record ignores the context.
type journalRecordIgnoresCancelDefect struct {
	journalOverEmbedded
}

// Record is the override of the planted defect.
func (d *journalRecordIgnoresCancelDefect) Record(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.journalOverEmbedded.Record(p0, p1)
}

// journalRecordIgnoresCancel returns the planted defect of a Journal whose
// Record ignores the context over a fresh subject.
func journalRecordIgnoresCancel(h *JournalHandle) Journal {
	return &journalRecordIgnoresCancelDefect{journalOverEmbedded: h.Fresh()}
}

// journalRecordSwallowsFailureDefect is the planted defect of a Journal
// whose Record returns no error when its dependency fails.
type journalRecordSwallowsFailureDefect struct {
	journalOverEmbedded
	h *JournalHandle
}

// Record is the override of the planted defect.
func (d *journalRecordSwallowsFailureDefect) Record(p0 context.Context, p1 Item) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.journalOverEmbedded.Record(p0, p1)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// journalRecordSwallowsFailure returns the planted defect of a Journal
// whose Record returns no error when its dependency fails over a fresh
// subject.
func journalRecordSwallowsFailure(h *JournalHandle) Journal {
	return &journalRecordSwallowsFailureDefect{journalOverEmbedded: h.Fresh(), h: h}
}

// journalEntryHalfAppliedDefect is the planted defect of a Journal whose
// Entry finds a part of a Record that failed.
type journalEntryHalfAppliedDefect struct {
	journalOverEmbedded
	h *JournalHandle
}

// Entry is the override of the planted defect.
func (d *journalEntryHalfAppliedDefect) Entry(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	failed := d.h.Injector().Failures() > 0
	r0, r1 = d.journalOverEmbedded.Entry(p0, p1)
	if failed {
		r0 = Item{Key: "Key-b"}
		r1 = nil
	}
	return r0, r1
}

// journalEntryHalfApplied returns the planted defect of a Journal whose
// Entry finds a part of a Record that failed over a fresh subject.
func journalEntryHalfApplied(h *JournalHandle) Journal {
	return &journalEntryHalfAppliedDefect{journalOverEmbedded: h.Fresh(), h: h}
}

// journalEntryHalfAppliedAfterThreeDefect is the planted defect of a
// Journal whose Entry finds a part of a failed Record from its fourth read
// on.
type journalEntryHalfAppliedAfterThreeDefect struct {
	journalOverEmbedded
	calls int
	h     *JournalHandle
}

// Entry is the override of the planted defect.
func (d *journalEntryHalfAppliedAfterThreeDefect) Entry(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	failed := d.h.Injector().Failures() > 0
	r0, r1 = d.journalOverEmbedded.Entry(p0, p1)
	d.calls++
	if failed && d.calls > 3 {
		r0 = Item{Key: "Key-b"}
		r1 = nil
	}
	return r0, r1
}

// journalEntryHalfAppliedAfterThree returns the planted defect of a
// Journal whose Entry finds a part of a failed Record from its fourth read
// on over a fresh subject.
func journalEntryHalfAppliedAfterThree(h *JournalHandle) Journal {
	return &journalEntryHalfAppliedAfterThreeDefect{journalOverEmbedded: h.Fresh(), h: h}
}

// journalEntryPanicsDefect is the planted defect of a Journal whose Entry
// panics.
type journalEntryPanicsDefect struct {
	journalOverEmbedded
}

// Entry is the override of the planted defect.
func (*journalEntryPanicsDefect) Entry(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Journal whose Entry panics")
}

// journalEntryPanics returns the planted defect of a Journal whose Entry
// panics over a fresh subject.
func journalEntryPanics(h *JournalHandle) Journal {
	return &journalEntryPanicsDefect{journalOverEmbedded: h.Fresh()}
}

// journalEntryIgnoresCancelDefect is the planted defect of a Journal whose
// Entry ignores the context.
type journalEntryIgnoresCancelDefect struct {
	journalOverEmbedded
}

// Entry is the override of the planted defect.
func (d *journalEntryIgnoresCancelDefect) Entry(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.journalOverEmbedded.Entry(p0, p1)
}

// journalEntryIgnoresCancel returns the planted defect of a Journal whose
// Entry ignores the context over a fresh subject.
func journalEntryIgnoresCancel(h *JournalHandle) Journal {
	return &journalEntryIgnoresCancelDefect{journalOverEmbedded: h.Fresh()}
}

// journalEntrySwallowsFailureDefect is the planted defect of a Journal
// whose Entry returns no error when its dependency fails.
type journalEntrySwallowsFailureDefect struct {
	journalOverEmbedded
	h *JournalHandle
}

// Entry is the override of the planted defect.
func (d *journalEntrySwallowsFailureDefect) Entry(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.journalOverEmbedded.Entry(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(Item)
		r1 = nil
	}
	return r0, r1
}

// journalEntrySwallowsFailure returns the planted defect of a Journal
// whose Entry returns no error when its dependency fails over a fresh
// subject.
func journalEntrySwallowsFailure(h *JournalHandle) Journal {
	return &journalEntrySwallowsFailureDefect{journalOverEmbedded: h.Fresh(), h: h}
}

// journalEntryValueBesideDefect is the planted defect of a Journal whose
// Entry returns a value beside an error.
type journalEntryValueBesideDefect struct {
	journalOverEmbedded
}

// Entry is the override of the planted defect.
func (d *journalEntryValueBesideDefect) Entry(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.journalOverEmbedded.Entry(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// journalEntryValueBeside returns the planted defect of a Journal whose
// Entry returns a value beside an error over a fresh subject.
func journalEntryValueBeside(h *JournalHandle) Journal {
	return &journalEntryValueBesideDefect{journalOverEmbedded: h.Fresh()}
}

// journalActions returns the actions of the generated sequences of Journal
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func journalActions(c *prop.Case, subject func() Journal) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Record",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "record.it")
				it.Key = c.Draw(prop.SampledFrom(keysString...), "record.it.key")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				key0, _ := args[0].(Item)
				call := c.History().Invoke(client, "Record", args, key0.Key)
				values, err := journalRecordCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Entry",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "entry.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Entry", args, args[0])
				values, err := journalEntryCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// journalRecordCall calls Record on subject with the arguments args, and
// returns the values of the call beside its error.
func journalRecordCall(ctx context.Context, subject Journal, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Record(ctx, a0)
}

// journalEntryCall calls Entry on subject with the arguments args, and
// returns the values of the call beside its error.
func journalEntryCall(ctx context.Context, subject Journal, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Entry(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Stock:
//
// Stock.Items/smoke: Items returns without a panic for derived inputs.
//
// Stock.Items/streams-what-was-stored: Items yields each value that Add
// stored.
//
// Stock.Items/streams-what-was-stored#generated: Items yields each value
// that Add stored.
//
// Stock.Items/iteration-sees-mutations: an iteration of Items yields each
// value that Add writes and no value that Remove removes while it runs.
//
// Stock.Each/smoke: Each returns without a panic for derived inputs.
//
// Stock.Each/streams-what-was-stored: Each yields each value that Add
// stored.
//
// Stock.Each/streams-what-was-stored#generated: Each yields each value
// that Add stored.
//
// Stock.Each/iteration-sees-mutations: an iteration of Each yields each
// value that Add writes and no value that Remove removes while it runs.
//
// Stock.Counts/smoke: Counts returns without a panic for derived inputs.
//
// Stock.Counts/streams-what-was-stored: Counts yields each value that Add
// stored.
//
// Stock.Counts/streams-what-was-stored#generated: Counts yields each value
// that Add stored.
//
// Stock.Counts/iteration-sees-mutations: an iteration of Counts yields
// each value that Add writes and no value that Remove removes while it
// runs.
//
// Stock.Sizes/smoke: Sizes returns without a panic for derived inputs.
//
// Stock.Sizes/streams-what-was-stored: Sizes yields each value that Add
// stored.
//
// Stock.Sizes/streams-what-was-stored#generated: Sizes yields each value
// that Add stored.
//
// Stock.Snapshot/smoke: Snapshot returns without a panic for derived
// inputs.
//
// Stock.Snapshot/cancel: Snapshot returns an error for a cancelled
// context.
//
// Stock.Snapshot/zero-on-error: Snapshot returns the zero value beside an
// error.
//
// Stock.Add/smoke: Add returns without a panic for derived inputs.
//
// Stock.Add/smoke#generated: Add returns without a panic for derived
// inputs.
//
// Stock.Add/cancel: Add returns an error for a cancelled context.
//
// Stock.Add/cancel#generated: Add returns an error for a cancelled
// context.
//
// Stock.Remove/smoke: Remove returns without a panic for derived inputs.
//
// Stock.Remove/smoke#generated: Remove returns without a panic for derived
// inputs.
//
// Stock.Remove/cancel: Remove returns an error for a cancelled context.
//
// Stock.Remove/cancel#generated: Remove returns an error for a cancelled
// context.

// The derived checks of Stock that witness their laws at the shallow end:
//
// Stock.Items/smoke
//
// Stock.Items/streams-what-was-stored
//
// Stock.Items/streams-what-was-stored#generated
//
// Stock.Each/smoke
//
// Stock.Each/streams-what-was-stored
//
// Stock.Each/streams-what-was-stored#generated
//
// Stock.Counts/smoke
//
// Stock.Counts/streams-what-was-stored
//
// Stock.Counts/streams-what-was-stored#generated
//
// Stock.Sizes/smoke
//
// Stock.Sizes/streams-what-was-stored
//
// Stock.Sizes/streams-what-was-stored#generated
//
// Stock.Snapshot/smoke
//
// Stock.Add/smoke
//
// Stock.Remove/smoke

// The refused claims of Stock:
//
// Stock.Items/cancel: Items cannot fail. The remedy is to give Items an
// error result.
//
// Stock.Each/cancel: Each cannot fail. The remedy is to give Each an error
// result.
//
// Stock.Counts/cancel: Counts cannot fail. The remedy is to give Counts an
// error result.
//
// Stock.Sizes/cancel: Sizes cannot fail. The remedy is to give Sizes an
// error result.
//
// Stock.Sizes/iteration-sees-mutations: Sizes returns no stream of a value
// that Add and Remove take. The remedy is to name with mutate and delete
// the callables that write and remove an element of the stream of Sizes.
//
// Stock.Snapshot/iteration-sees-mutations: Snapshot returns no stream of a
// value that Add and Remove take. The remedy is to name with mutate and
// delete the callables that write and remove an element of the stream of
// Snapshot.

// The unclaimed classifications of Stock:
//
// The shape:writer of Stock.Add has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:deleter of Stock.Remove has the stance silent. A delete is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Stock writes.
type (
	StockHandle = conform.Handle[Stock]
	StockCheck  = conform.Check[Stock]
	StockOption = conform.Option[Stock]
)

// StockChecks contains the name of each callable of Stock, with the
// identities of its checks.
var StockChecks = struct {
	Items    stockItems
	Each     stockEach
	Counts   stockCounts
	Sizes    stockSizes
	Snapshot stockSnapshot
	Add      stockAdd
	Remove   stockRemove
}{
	Items:    stockItems{"Stock.Items"},
	Each:     stockEach{"Stock.Each"},
	Counts:   stockCounts{"Stock.Counts"},
	Sizes:    stockSizes{"Stock.Sizes"},
	Snapshot: stockSnapshot{"Stock.Snapshot"},
	Add:      stockAdd{"Stock.Add"},
	Remove:   stockRemove{"Stock.Remove"},
}

// stockItems is the name of Stock.Items, with the identities of its
// checks.
type stockItems struct {
	conform.Callable[Stock]
}

// Smoke returns the identity of the check that Items returns without a
// panic for derived inputs.
func (c stockItems) Smoke() conform.ID[Stock] {
	return c.Check("smoke")
}

// StreamsWhatWasStored returns the identity of the check that Items yields
// each value that Add stored.
func (c stockItems) StreamsWhatWasStored() conform.ID[Stock] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Items yields each value that Add stored.
func (c stockItems) StreamsWhatWasStoredGenerated() conform.ID[Stock] {
	return c.Check("streams-what-was-stored#generated")
}

// IterationSeesMutations returns the identity of the check that an
// iteration of Items yields each value that Add writes and no value that
// Remove removes while it runs.
func (c stockItems) IterationSeesMutations() conform.ID[Stock] {
	return c.Check("iteration-sees-mutations")
}

// stockEach is the name of Stock.Each, with the identities of its checks.
type stockEach struct {
	conform.Callable[Stock]
}

// Smoke returns the identity of the check that Each returns without a
// panic for derived inputs.
func (c stockEach) Smoke() conform.ID[Stock] {
	return c.Check("smoke")
}

// StreamsWhatWasStored returns the identity of the check that Each yields
// each value that Add stored.
func (c stockEach) StreamsWhatWasStored() conform.ID[Stock] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Each yields each value that Add stored.
func (c stockEach) StreamsWhatWasStoredGenerated() conform.ID[Stock] {
	return c.Check("streams-what-was-stored#generated")
}

// IterationSeesMutations returns the identity of the check that an
// iteration of Each yields each value that Add writes and no value that
// Remove removes while it runs.
func (c stockEach) IterationSeesMutations() conform.ID[Stock] {
	return c.Check("iteration-sees-mutations")
}

// stockCounts is the name of Stock.Counts, with the identities of its
// checks.
type stockCounts struct {
	conform.Callable[Stock]
}

// Smoke returns the identity of the check that Counts returns without a
// panic for derived inputs.
func (c stockCounts) Smoke() conform.ID[Stock] {
	return c.Check("smoke")
}

// StreamsWhatWasStored returns the identity of the check that Counts
// yields each value that Add stored.
func (c stockCounts) StreamsWhatWasStored() conform.ID[Stock] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Counts yields each value that Add stored.
func (c stockCounts) StreamsWhatWasStoredGenerated() conform.ID[Stock] {
	return c.Check("streams-what-was-stored#generated")
}

// IterationSeesMutations returns the identity of the check that an
// iteration of Counts yields each value that Add writes and no value that
// Remove removes while it runs.
func (c stockCounts) IterationSeesMutations() conform.ID[Stock] {
	return c.Check("iteration-sees-mutations")
}

// stockSizes is the name of Stock.Sizes, with the identities of its
// checks.
type stockSizes struct {
	conform.Callable[Stock]
}

// Smoke returns the identity of the check that Sizes returns without a
// panic for derived inputs.
func (c stockSizes) Smoke() conform.ID[Stock] {
	return c.Check("smoke")
}

// StreamsWhatWasStored returns the identity of the check that Sizes yields
// each value that Add stored.
func (c stockSizes) StreamsWhatWasStored() conform.ID[Stock] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Sizes yields each value that Add stored.
func (c stockSizes) StreamsWhatWasStoredGenerated() conform.ID[Stock] {
	return c.Check("streams-what-was-stored#generated")
}

// stockSnapshot is the name of Stock.Snapshot, with the identities of its
// checks.
type stockSnapshot struct {
	conform.Callable[Stock]
}

// Smoke returns the identity of the check that Snapshot returns without a
// panic for derived inputs.
func (c stockSnapshot) Smoke() conform.ID[Stock] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Snapshot returns an error
// for a cancelled context.
func (c stockSnapshot) Cancel() conform.ID[Stock] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Snapshot returns the
// zero value beside an error.
func (c stockSnapshot) ZeroOnError() conform.ID[Stock] {
	return c.Check("zero-on-error")
}

// stockAdd is the name of Stock.Add, with the identities of its checks.
type stockAdd struct {
	conform.Callable[Stock]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c stockAdd) Smoke() conform.ID[Stock] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Add returns
// without a panic for derived inputs.
func (c stockAdd) SmokeGenerated() conform.ID[Stock] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c stockAdd) Cancel() conform.ID[Stock] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Add returns an
// error for a cancelled context.
func (c stockAdd) CancelGenerated() conform.ID[Stock] {
	return c.Check("cancel#generated")
}

// stockRemove is the name of Stock.Remove, with the identities of its
// checks.
type stockRemove struct {
	conform.Callable[Stock]
}

// Smoke returns the identity of the check that Remove returns without a
// panic for derived inputs.
func (c stockRemove) Smoke() conform.ID[Stock] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Remove returns
// without a panic for derived inputs.
func (c stockRemove) SmokeGenerated() conform.ID[Stock] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Remove returns an error
// for a cancelled context.
func (c stockRemove) Cancel() conform.ID[Stock] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Remove returns an
// error for a cancelled context.
func (c stockRemove) CancelGenerated() conform.ID[Stock] {
	return c.Check("cancel#generated")
}

// StockSuite returns the checks of Stock, and the calls and the budgets of
// its callables.
func StockSuite() conform.Suite[Stock] {
	return conform.Suite[Stock]{
		Declaration: "Stock",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Stock]{
			{
				ID:          StockChecks.Items.Smoke(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Items(ctx)
						return nil
					})
				},
				Defect: stockItemsPanics,
			},
			{
				ID:          StockChecks.Items.StreamsWhatWasStored(),
				Statement:   "Items yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Items(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: stockItemsDrainsNothing,
			},
			{
				ID:          StockChecks.Items.StreamsWhatWasStoredGenerated(),
				Statement:   "Items yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Items(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: stockItemsDrainsNothing,
			},
			{
				ID:          StockChecks.Items.IterationSeesMutations(),
				Statement:   "an iteration of Items yields each value that Add writes and no value that Remove removes while it runs",
				Claim:       spec.IterationSeesMutations,
				Probe:       spec.IterationsSeeTheMutations,
				Reason:      spec.StaleIteration,
				Law:         spec.IterationReflectsItsSource,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *StockHandle, c *prop.Case) {
					probe.IterationsSeeTheMutations(
						h,
						c,
						stockActions,
						func(ctx context.Context) (func() (string, bool), func(), error) {
							r0 := h.Subject().Items(ctx)
							next, stop := iter.Pull2(r0)
							return func() (string, bool) {
								v, err, ok := next()
								return v, ok && err == nil
							}, stop, nil
						},
						"Add",
						"Remove",
						func(args []any) string {
							v, _ := args[0].(string)
							return v
						},
						func(args []any) string {
							v, _ := args[0].(string)
							return v
						},
					)
				},
				Defect: stockItemsIteratesACopy,
			},
			{
				ID:          StockChecks.Each.Smoke(),
				Statement:   "Each returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Each(ctx)
						return nil
					})
				},
				Defect: stockEachPanics,
			},
			{
				ID:          StockChecks.Each.StreamsWhatWasStored(),
				Statement:   "Each yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: stockEachDrainsNothing,
			},
			{
				ID:          StockChecks.Each.StreamsWhatWasStoredGenerated(),
				Statement:   "Each yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: stockEachDrainsNothing,
			},
			{
				ID:          StockChecks.Each.IterationSeesMutations(),
				Statement:   "an iteration of Each yields each value that Add writes and no value that Remove removes while it runs",
				Claim:       spec.IterationSeesMutations,
				Probe:       spec.IterationsSeeTheMutations,
				Reason:      spec.StaleIteration,
				Law:         spec.IterationReflectsItsSource,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *StockHandle, c *prop.Case) {
					probe.IterationsSeeTheMutations(
						h,
						c,
						stockActions,
						func(ctx context.Context) (func() (string, bool), func(), error) {
							r0 := h.Subject().Each(ctx)
							next, stop := iter.Pull(r0)
							return next, stop, nil
						},
						"Add",
						"Remove",
						func(args []any) string {
							v, _ := args[0].(string)
							return v
						},
						func(args []any) string {
							v, _ := args[0].(string)
							return v
						},
					)
				},
				Defect: stockEachIteratesACopy,
			},
			{
				ID:          StockChecks.Counts.Smoke(),
				Statement:   "Counts returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Counts(ctx)
						return nil
					})
				},
				Defect: stockCountsPanics,
			},
			{
				ID:          StockChecks.Counts.StreamsWhatWasStored(),
				Statement:   "Counts yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Counts(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := string("item-a")
								if calls > 1 {
									w = string("item-b")
								}
								return h.Subject().Add(ctx, w)
							}
						}(),
					)
				},
				Defect: stockCountsDrainsNothing,
			},
			{
				ID:          StockChecks.Counts.StreamsWhatWasStoredGenerated(),
				Statement:   "Counts yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Counts(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := drawnItem[0]
								if calls > 1 {
									w = drawnItem[1]
								}
								return h.Subject().Add(ctx, w)
							}
						}(),
					)
				},
				Defect: stockCountsDrainsNothing,
			},
			{
				ID:          StockChecks.Counts.IterationSeesMutations(),
				Statement:   "an iteration of Counts yields each value that Add writes and no value that Remove removes while it runs",
				Claim:       spec.IterationSeesMutations,
				Probe:       spec.IterationsSeeTheMutations,
				Reason:      spec.StaleIteration,
				Law:         spec.IterationReflectsItsSource,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *StockHandle, c *prop.Case) {
					probe.IterationsSeeTheMutations(
						h,
						c,
						stockActions,
						func(ctx context.Context) (func() (string, bool), func(), error) {
							r0 := h.Subject().Counts(ctx)
							next, stop := iter.Pull2(r0)
							return func() (string, bool) {
								v, _, ok := next()
								return v, ok
							}, stop, nil
						},
						"Add",
						"Remove",
						func(args []any) string {
							v, _ := args[0].(string)
							return v
						},
						func(args []any) string {
							v, _ := args[0].(string)
							return v
						},
					)
				},
				Defect: stockCountsIteratesACopy,
			},
			{
				ID:          StockChecks.Sizes.Smoke(),
				Statement:   "Sizes returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Sizes(ctx)
						return nil
					})
				},
				Defect: stockSizesPanics,
			},
			{
				ID:          StockChecks.Sizes.StreamsWhatWasStored(),
				Statement:   "Sizes yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]int, error) {
						r0 := h.Subject().Sizes(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: stockSizesDrainsNothing,
			},
			{
				ID:          StockChecks.Sizes.StreamsWhatWasStoredGenerated(),
				Statement:   "Sizes yields each value that Add stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]int, error) {
						r0 := h.Subject().Sizes(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: stockSizesDrainsNothing,
			},
			{
				ID:          StockChecks.Snapshot.Smoke(),
				Statement:   "Snapshot returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Snapshot(ctx)
						return err
					})
				},
				Defect: stockSnapshotPanics,
			},
			{
				ID:          StockChecks.Snapshot.Cancel(),
				Statement:   "Snapshot returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StockHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Snapshot(ctx)
						return err
					})
				},
				Defect: stockSnapshotIgnoresCancel,
			},
			{
				ID:          StockChecks.Snapshot.ZeroOnError(),
				Statement:   "Snapshot returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StockHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Snapshot(ctx)
					}, []string(nil))
				},
				Defect: stockSnapshotValueBeside,
			},
			{
				ID:          StockChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("item-a"))
					})
				},
				Defect: stockAddPanics,
			},
			{
				ID:          StockChecks.Add.SmokeGenerated(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnItem)
					})
				},
				Defect: stockAddPanics,
			},
			{
				ID:          StockChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StockHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("item-a"))
					})
				},
				Defect: stockAddIgnoresCancel,
			},
			{
				ID:          StockChecks.Add.CancelGenerated(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnItem)
					})
				},
				Defect: stockAddIgnoresCancel,
			},
			{
				ID:          StockChecks.Remove.Smoke(),
				Statement:   "Remove returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Remove(ctx, string("item-a"))
					})
				},
				Defect: stockRemovePanics,
			},
			{
				ID:          StockChecks.Remove.SmokeGenerated(),
				Statement:   "Remove returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Remove(ctx, drawnItem)
					})
				},
				Defect: stockRemovePanics,
			},
			{
				ID:          StockChecks.Remove.Cancel(),
				Statement:   "Remove returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StockHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Remove(ctx, string("item-a"))
					})
				},
				Defect: stockRemoveIgnoresCancel,
			},
			{
				ID:          StockChecks.Remove.CancelGenerated(),
				Statement:   "Remove returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StockHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Remove(ctx, drawnItem)
					})
				},
				Defect: stockRemoveIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Stock.Items",
				Gate:     spec.Failure,
				Why:      "Items cannot fail",
				Remedy:   "give Items an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Stock.Each",
				Gate:     spec.Failure,
				Why:      "Each cannot fail",
				Remedy:   "give Each an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Stock.Counts",
				Gate:     spec.Failure,
				Why:      "Counts cannot fail",
				Remedy:   "give Counts an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Stock.Sizes",
				Gate:     spec.Failure,
				Why:      "Sizes cannot fail",
				Remedy:   "give Sizes an error result",
			},
			{
				Claim:    spec.IterationSeesMutations,
				Callable: "Stock.Sizes",
				Gate:     spec.StreamMutations,
				Why:      "Sizes returns no stream of a value that Add and Remove take",
				Remedy:   "name with mutate and delete the callables that write and remove an element of the stream of Sizes",
			},
			{
				Claim:    spec.IterationSeesMutations,
				Callable: "Stock.Snapshot",
				Gate:     spec.StreamMutations,
				Why:      "Snapshot returns no stream of a value that Add and Remove take",
				Remedy:   "name with mutate and delete the callables that write and remove an element of the stream of Snapshot",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Stock.Add",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:deleter",
				Callable:       "Stock.Remove",
				Stance:         spec.Silent,
				Note:           "A delete is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Stock]{
			{
				Calls:    stockItemsCalls,
				Over:     newStockOver,
				Callable: StockChecks.Items.Callable,
			},
			{
				Calls:    stockEachCalls,
				Over:     newStockOver,
				Callable: StockChecks.Each.Callable,
			},
			{
				Calls:    stockCountsCalls,
				Over:     newStockOver,
				Callable: StockChecks.Counts.Callable,
			},
			{
				Calls:    stockSizesCalls,
				Over:     newStockOver,
				Callable: StockChecks.Sizes.Callable,
			},
			{
				Calls:    stockSnapshotCalls,
				Over:     newStockOver,
				Callable: StockChecks.Snapshot.Callable,
			},
			{
				Calls:    stockAddCalls,
				Over:     newStockOver,
				Callable: StockChecks.Add.Callable,
			},
			{
				Calls:    stockRemoveCalls,
				Over:     newStockOver,
				Callable: StockChecks.Remove.Callable,
			},
		},
	}
}

// RunStockSuite applies the checks of Stock, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunStockSuite[T Stock](t *testing.T, newSubject func() T, opts ...StockOption) {
	t.Helper()
	RunStockSuiteFromHandle(t, func(*StockHandle) T { return newSubject() }, opts...)
}

// RunStockSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStockSuiteFromHandle[T Stock](
	t *testing.T,
	newSubject func(h *StockHandle) T,
	opts ...StockOption,
) {
	t.Helper()
	conform.Run(t, StockSuite(), stockSubjects(newSubject), opts...)
}

// FuzzStockSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzStockSuite[T Stock](
	f *testing.F,
	callable conform.Callable[Stock],
	newSubject func(h *StockHandle) T,
	opts ...StockOption,
) {
	f.Helper()
	conform.Fuzz(f, StockSuite(), callable, stockSubjects(newSubject), opts...)
}

// RunStockBudgets applies the allocation ceilings of Stock and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunStockBudgets[T Stock](
	t *testing.T,
	newSubject func(h *StockHandle) T,
	opts ...StockOption,
) {
	t.Helper()
	conform.RunBudgets(t, StockSuite(), stockSubjects(newSubject), opts...)
}

// BenchStockSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchStockSuite[T Stock](
	b conform.B,
	callable conform.Callable[Stock],
	newSubject func(h *StockHandle) T,
	opts ...StockOption,
) {
	b.Helper()
	conform.Bench(b, StockSuite(), callable, stockSubjects(newSubject), opts...)
}

// stockSubjects returns the factory of the subjects of Stock that calls
// newSubject with the handle of each check.
func stockSubjects[T Stock](newSubject func(h *StockHandle) T) func(h *StockHandle) Stock {
	return func(h *StockHandle) Stock { return newSubject(h) }
}

// stockItemsCalls generates the calls of Stock.Items with the inputs that
// each case draws.
var stockItemsCalls = prop.Composite(func(*prop.Case) func(context.Context, Stock) error {
	return func(ctx context.Context, subject Stock) error {
		_ = subject.Items(ctx)
		return nil
	}
})

// stockEachCalls generates the calls of Stock.Each with the inputs that
// each case draws.
var stockEachCalls = prop.Composite(func(*prop.Case) func(context.Context, Stock) error {
	return func(ctx context.Context, subject Stock) error {
		_ = subject.Each(ctx)
		return nil
	}
})

// stockCountsCalls generates the calls of Stock.Counts with the inputs
// that each case draws.
var stockCountsCalls = prop.Composite(func(*prop.Case) func(context.Context, Stock) error {
	return func(ctx context.Context, subject Stock) error {
		_ = subject.Counts(ctx)
		return nil
	}
})

// stockSizesCalls generates the calls of Stock.Sizes with the inputs that
// each case draws.
var stockSizesCalls = prop.Composite(func(*prop.Case) func(context.Context, Stock) error {
	return func(ctx context.Context, subject Stock) error {
		_ = subject.Sizes(ctx)
		return nil
	}
})

// stockSnapshotCalls generates the calls of Stock.Snapshot with the inputs
// that each case draws.
var stockSnapshotCalls = prop.Composite(func(*prop.Case) func(context.Context, Stock) error {
	return func(ctx context.Context, subject Stock) error {
		_, err := subject.Snapshot(ctx)
		return err
	}
})

// stockAddCalls generates the calls of Stock.Add with the inputs that each
// case draws.
var stockAddCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stock) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Stock) error {
		return subject.Add(ctx, item)
	}
})

// stockRemoveCalls generates the calls of Stock.Remove with the inputs
// that each case draws.
var stockRemoveCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stock) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Stock) error {
		return subject.Remove(ctx, item)
	}
})

// stockOverEmbedded is Stock under the name of the field that stockOver
// embeds.
type stockOverEmbedded = Stock

// stockOver is a subject of Stock whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type stockOver struct {
	stockOverEmbedded
	spend func()
}

// Items calls spend and returns its zero results.
func (o stockOver) Items(context.Context) (r0 iter.Seq2[string, error]) {
	o.spend()
	return r0
}

// Each calls spend and returns its zero results.
func (o stockOver) Each(context.Context) (r0 iter.Seq[string]) {
	o.spend()
	return r0
}

// Counts calls spend and returns its zero results.
func (o stockOver) Counts(context.Context) (r0 iter.Seq2[string, int]) {
	o.spend()
	return r0
}

// Sizes calls spend and returns its zero results.
func (o stockOver) Sizes(context.Context) (r0 iter.Seq[int]) {
	o.spend()
	return r0
}

// Snapshot calls spend and returns its zero results.
func (o stockOver) Snapshot(context.Context) (r0 []string, r1 error) {
	o.spend()
	return r0, r1
}

// Add calls spend and returns its zero results.
func (o stockOver) Add(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Remove calls spend and returns its zero results.
func (o stockOver) Remove(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// newStockOver returns the planted defect of a ceiling of Stock, whose
// callables call spend.
func newStockOver(_ *StockHandle, spend func()) Stock {
	return stockOver{spend: spend}
}

// stockItemsPanicsDefect is the planted defect of a Stock whose Items
// panics.
type stockItemsPanicsDefect struct {
	stockOverEmbedded
}

// Items is the override of the planted defect.
func (*stockItemsPanicsDefect) Items(p0 context.Context) (r0 iter.Seq2[string, error]) {
	panic("the planted defect of a Stock whose Items panics")
}

// stockItemsPanics returns the planted defect of a Stock whose Items
// panics over a fresh subject.
func stockItemsPanics(h *StockHandle) Stock {
	return &stockItemsPanicsDefect{stockOverEmbedded: h.Fresh()}
}

// stockItemsDrainsNothingDefect is the planted defect of a Stock whose
// Items loses what Add stored.
type stockItemsDrainsNothingDefect struct {
	stockOverEmbedded
}

// Items is the override of the planted defect.
func (*stockItemsDrainsNothingDefect) Items(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = func(func(string, error) bool) {}
	return r0
}

// stockItemsDrainsNothing returns the planted defect of a Stock whose
// Items loses what Add stored over a fresh subject.
func stockItemsDrainsNothing(h *StockHandle) Stock {
	return &stockItemsDrainsNothingDefect{stockOverEmbedded: h.Fresh()}
}

// stockItemsIteratesACopyDefect is the planted defect of a Stock whose
// Items iterates a copy of its source.
type stockItemsIteratesACopyDefect struct {
	stockOverEmbedded
}

// Items is the override of the planted defect.
func (d *stockItemsIteratesACopyDefect) Items(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = d.stockOverEmbedded.Items(p0)
	var firsts []string
	var seconds []error
	for first, second := range r0 {
		firsts, seconds = append(firsts, first), append(seconds, second)
	}
	r0 = func(yield func(string, error) bool) {
		for i := range firsts {
			if !yield(firsts[i], seconds[i]) {
				return
			}
		}
	}
	return r0
}

// stockItemsIteratesACopy returns the planted defect of a Stock whose
// Items iterates a copy of its source over a fresh subject.
func stockItemsIteratesACopy(h *StockHandle) Stock {
	return &stockItemsIteratesACopyDefect{stockOverEmbedded: h.Fresh()}
}

// stockEachPanicsDefect is the planted defect of a Stock whose Each
// panics.
type stockEachPanicsDefect struct {
	stockOverEmbedded
}

// Each is the override of the planted defect.
func (*stockEachPanicsDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	panic("the planted defect of a Stock whose Each panics")
}

// stockEachPanics returns the planted defect of a Stock whose Each panics
// over a fresh subject.
func stockEachPanics(h *StockHandle) Stock {
	return &stockEachPanicsDefect{stockOverEmbedded: h.Fresh()}
}

// stockEachDrainsNothingDefect is the planted defect of a Stock whose Each
// loses what Add stored.
type stockEachDrainsNothingDefect struct {
	stockOverEmbedded
}

// Each is the override of the planted defect.
func (*stockEachDrainsNothingDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = func(func(string) bool) {}
	return r0
}

// stockEachDrainsNothing returns the planted defect of a Stock whose Each
// loses what Add stored over a fresh subject.
func stockEachDrainsNothing(h *StockHandle) Stock {
	return &stockEachDrainsNothingDefect{stockOverEmbedded: h.Fresh()}
}

// stockEachIteratesACopyDefect is the planted defect of a Stock whose Each
// iterates a copy of its source.
type stockEachIteratesACopyDefect struct {
	stockOverEmbedded
}

// Each is the override of the planted defect.
func (d *stockEachIteratesACopyDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = d.stockOverEmbedded.Each(p0)
	r0 = slices.Values(slices.Collect(r0))
	return r0
}

// stockEachIteratesACopy returns the planted defect of a Stock whose Each
// iterates a copy of its source over a fresh subject.
func stockEachIteratesACopy(h *StockHandle) Stock {
	return &stockEachIteratesACopyDefect{stockOverEmbedded: h.Fresh()}
}

// stockCountsPanicsDefect is the planted defect of a Stock whose Counts
// panics.
type stockCountsPanicsDefect struct {
	stockOverEmbedded
}

// Counts is the override of the planted defect.
func (*stockCountsPanicsDefect) Counts(p0 context.Context) (r0 iter.Seq2[string, int]) {
	panic("the planted defect of a Stock whose Counts panics")
}

// stockCountsPanics returns the planted defect of a Stock whose Counts
// panics over a fresh subject.
func stockCountsPanics(h *StockHandle) Stock {
	return &stockCountsPanicsDefect{stockOverEmbedded: h.Fresh()}
}

// stockCountsDrainsNothingDefect is the planted defect of a Stock whose
// Counts loses what Add stored.
type stockCountsDrainsNothingDefect struct {
	stockOverEmbedded
}

// Counts is the override of the planted defect.
func (*stockCountsDrainsNothingDefect) Counts(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = func(func(string, int) bool) {}
	return r0
}

// stockCountsDrainsNothing returns the planted defect of a Stock whose
// Counts loses what Add stored over a fresh subject.
func stockCountsDrainsNothing(h *StockHandle) Stock {
	return &stockCountsDrainsNothingDefect{stockOverEmbedded: h.Fresh()}
}

// stockCountsIteratesACopyDefect is the planted defect of a Stock whose
// Counts iterates a copy of its source.
type stockCountsIteratesACopyDefect struct {
	stockOverEmbedded
}

// Counts is the override of the planted defect.
func (d *stockCountsIteratesACopyDefect) Counts(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = d.stockOverEmbedded.Counts(p0)
	var firsts []string
	var seconds []int
	for first, second := range r0 {
		firsts, seconds = append(firsts, first), append(seconds, second)
	}
	r0 = func(yield func(string, int) bool) {
		for i := range firsts {
			if !yield(firsts[i], seconds[i]) {
				return
			}
		}
	}
	return r0
}

// stockCountsIteratesACopy returns the planted defect of a Stock whose
// Counts iterates a copy of its source over a fresh subject.
func stockCountsIteratesACopy(h *StockHandle) Stock {
	return &stockCountsIteratesACopyDefect{stockOverEmbedded: h.Fresh()}
}

// stockSizesPanicsDefect is the planted defect of a Stock whose Sizes
// panics.
type stockSizesPanicsDefect struct {
	stockOverEmbedded
}

// Sizes is the override of the planted defect.
func (*stockSizesPanicsDefect) Sizes(p0 context.Context) (r0 iter.Seq[int]) {
	panic("the planted defect of a Stock whose Sizes panics")
}

// stockSizesPanics returns the planted defect of a Stock whose Sizes
// panics over a fresh subject.
func stockSizesPanics(h *StockHandle) Stock {
	return &stockSizesPanicsDefect{stockOverEmbedded: h.Fresh()}
}

// stockSizesDrainsNothingDefect is the planted defect of a Stock whose
// Sizes loses what Add stored.
type stockSizesDrainsNothingDefect struct {
	stockOverEmbedded
}

// Sizes is the override of the planted defect.
func (*stockSizesDrainsNothingDefect) Sizes(p0 context.Context) (r0 iter.Seq[int]) {
	r0 = func(func(int) bool) {}
	return r0
}

// stockSizesDrainsNothing returns the planted defect of a Stock whose
// Sizes loses what Add stored over a fresh subject.
func stockSizesDrainsNothing(h *StockHandle) Stock {
	return &stockSizesDrainsNothingDefect{stockOverEmbedded: h.Fresh()}
}

// stockSnapshotPanicsDefect is the planted defect of a Stock whose
// Snapshot panics.
type stockSnapshotPanicsDefect struct {
	stockOverEmbedded
}

// Snapshot is the override of the planted defect.
func (*stockSnapshotPanicsDefect) Snapshot(p0 context.Context) (r0 []string, r1 error) {
	panic("the planted defect of a Stock whose Snapshot panics")
}

// stockSnapshotPanics returns the planted defect of a Stock whose Snapshot
// panics over a fresh subject.
func stockSnapshotPanics(h *StockHandle) Stock {
	return &stockSnapshotPanicsDefect{stockOverEmbedded: h.Fresh()}
}

// stockSnapshotIgnoresCancelDefect is the planted defect of a Stock whose
// Snapshot ignores the context.
type stockSnapshotIgnoresCancelDefect struct {
	stockOverEmbedded
}

// Snapshot is the override of the planted defect.
func (d *stockSnapshotIgnoresCancelDefect) Snapshot(p0 context.Context) (r0 []string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.stockOverEmbedded.Snapshot(p0)
}

// stockSnapshotIgnoresCancel returns the planted defect of a Stock whose
// Snapshot ignores the context over a fresh subject.
func stockSnapshotIgnoresCancel(h *StockHandle) Stock {
	return &stockSnapshotIgnoresCancelDefect{stockOverEmbedded: h.Fresh()}
}

// stockSnapshotValueBesideDefect is the planted defect of a Stock whose
// Snapshot returns a value beside an error.
type stockSnapshotValueBesideDefect struct {
	stockOverEmbedded
}

// Snapshot is the override of the planted defect.
func (d *stockSnapshotValueBesideDefect) Snapshot(p0 context.Context) (r0 []string, r1 error) {
	r0, r1 = d.stockOverEmbedded.Snapshot(p0)
	if r1 != nil {
		r0 = []string{"sample-b"}
	}
	return r0, r1
}

// stockSnapshotValueBeside returns the planted defect of a Stock whose
// Snapshot returns a value beside an error over a fresh subject.
func stockSnapshotValueBeside(h *StockHandle) Stock {
	return &stockSnapshotValueBesideDefect{stockOverEmbedded: h.Fresh()}
}

// stockAddPanicsDefect is the planted defect of a Stock whose Add panics.
type stockAddPanicsDefect struct {
	stockOverEmbedded
}

// Add is the override of the planted defect.
func (*stockAddPanicsDefect) Add(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Stock whose Add panics")
}

// stockAddPanics returns the planted defect of a Stock whose Add panics
// over a fresh subject.
func stockAddPanics(h *StockHandle) Stock {
	return &stockAddPanicsDefect{stockOverEmbedded: h.Fresh()}
}

// stockAddIgnoresCancelDefect is the planted defect of a Stock whose Add
// ignores the context.
type stockAddIgnoresCancelDefect struct {
	stockOverEmbedded
}

// Add is the override of the planted defect.
func (d *stockAddIgnoresCancelDefect) Add(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.stockOverEmbedded.Add(p0, p1)
}

// stockAddIgnoresCancel returns the planted defect of a Stock whose Add
// ignores the context over a fresh subject.
func stockAddIgnoresCancel(h *StockHandle) Stock {
	return &stockAddIgnoresCancelDefect{stockOverEmbedded: h.Fresh()}
}

// stockRemovePanicsDefect is the planted defect of a Stock whose Remove
// panics.
type stockRemovePanicsDefect struct {
	stockOverEmbedded
}

// Remove is the override of the planted defect.
func (*stockRemovePanicsDefect) Remove(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Stock whose Remove panics")
}

// stockRemovePanics returns the planted defect of a Stock whose Remove
// panics over a fresh subject.
func stockRemovePanics(h *StockHandle) Stock {
	return &stockRemovePanicsDefect{stockOverEmbedded: h.Fresh()}
}

// stockRemoveIgnoresCancelDefect is the planted defect of a Stock whose
// Remove ignores the context.
type stockRemoveIgnoresCancelDefect struct {
	stockOverEmbedded
}

// Remove is the override of the planted defect.
func (d *stockRemoveIgnoresCancelDefect) Remove(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.stockOverEmbedded.Remove(p0, p1)
}

// stockRemoveIgnoresCancel returns the planted defect of a Stock whose
// Remove ignores the context over a fresh subject.
func stockRemoveIgnoresCancel(h *StockHandle) Stock {
	return &stockRemoveIgnoresCancelDefect{stockOverEmbedded: h.Fresh()}
}

// stockActions returns the actions of the generated sequences of Stock
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func stockActions(c *prop.Case, subject func() Stock) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Items",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Items", args)
				values, err := stockItemsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Each",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Each", args)
				values, err := stockEachCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Counts",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Counts", args)
				values, err := stockCountsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Sizes",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Sizes", args)
				values, err := stockSizesCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Snapshot",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Snapshot", args)
				values, err := stockSnapshotCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Add",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.Of[string](), "add.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Add", args)
				values, err := stockAddCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Remove",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.SampledFrom(keysString...), "remove.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Remove", args, args[0])
				values, err := stockRemoveCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// stockItemsCall calls Items on subject with the arguments args, and
// returns the values of the call beside its error.
func stockItemsCall(ctx context.Context, subject Stock, args []any) ([]any, error) {
	r0 := subject.Items(ctx)
	var drained []string
	for e, err := range r0 {
		if err != nil {
			return []any{drained}, err
		}
		drained = append(drained, e)
	}
	return []any{drained}, nil
}

// stockEachCall calls Each on subject with the arguments args, and returns
// the values of the call beside its error.
func stockEachCall(ctx context.Context, subject Stock, args []any) ([]any, error) {
	r0 := subject.Each(ctx)
	return []any{slices.Collect(r0)}, nil
}

// stockCountsCall calls Counts on subject with the arguments args, and
// returns the values of the call beside its error.
func stockCountsCall(ctx context.Context, subject Stock, args []any) ([]any, error) {
	r0 := subject.Counts(ctx)
	var drained []probe.Pair[string, int]
	for k, v := range r0 {
		drained = append(drained, probe.Pair[string, int]{First: k, Second: v})
	}
	return []any{drained}, nil
}

// stockSizesCall calls Sizes on subject with the arguments args, and
// returns the values of the call beside its error.
func stockSizesCall(ctx context.Context, subject Stock, args []any) ([]any, error) {
	r0 := subject.Sizes(ctx)
	return []any{slices.Collect(r0)}, nil
}

// stockSnapshotCall calls Snapshot on subject with the arguments args, and
// returns the values of the call beside its error.
func stockSnapshotCall(ctx context.Context, subject Stock, args []any) ([]any, error) {
	r0, err := subject.Snapshot(ctx)
	return []any{r0}, err
}

// stockAddCall calls Add on subject with the arguments args, and returns
// the values of the call beside its error.
func stockAddCall(ctx context.Context, subject Stock, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Add(ctx, a0)
}

// stockRemoveCall calls Remove on subject with the arguments args, and
// returns the values of the call beside its error.
func stockRemoveCall(ctx context.Context, subject Stock, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Remove(ctx, a0)
}

// The derived checks of Outbox:
//
// Outbox.Append/smoke: Append returns without a panic for derived inputs.
//
// Outbox.Append/smoke#generated: Append returns without a panic for
// derived inputs.
//
// Outbox.Append/cancel: Append returns an error for a cancelled context.
//
// Outbox.Append/cancel#generated: Append returns an error for a cancelled
// context.
//
// Outbox.Append/delivers-each-append: each message that Append records
// arrives on the channel of Subscribe.
//
// Outbox.Subscribe/smoke: Subscribe returns without a panic for derived
// inputs.
//
// Outbox.Subscribe/cancel: Subscribe returns an error for a cancelled
// context.
//
// Outbox.Note/smoke: Note returns without a panic for derived inputs.
//
// Outbox.Note/smoke#generated: Note returns without a panic for derived
// inputs.
//
// Outbox.Note/cancel: Note returns an error for a cancelled context.
//
// Outbox.Note/cancel#generated: Note returns an error for a cancelled
// context.
//
// Outbox.Notes/smoke: Notes returns without a panic for derived inputs.
//
// Outbox.Notes/cancel: Notes returns an error for a cancelled context.
//
// Outbox.Log/smoke: Log returns without a panic for derived inputs.
//
// Outbox.Log/smoke#generated: Log returns without a panic for derived
// inputs.
//
// Outbox.Log/cancel: Log returns an error for a cancelled context.
//
// Outbox.Log/cancel#generated: Log returns an error for a cancelled
// context.
//
// Outbox.Logged/smoke: Logged returns without a panic for derived inputs.
//
// Outbox.Logged/cancel: Logged returns an error for a cancelled context.
//
// Outbox.Logged/zero-on-error: Logged returns the zero value beside an
// error.
//
// Outbox.Post/smoke: Post returns without a panic for derived inputs.
//
// Outbox.Post/smoke#generated: Post returns without a panic for derived
// inputs.
//
// Outbox.Post/cancel: Post returns an error for a cancelled context.
//
// Outbox.Post/cancel#generated: Post returns an error for a cancelled
// context.
//
// Outbox.Queue/smoke: Queue returns without a panic for derived inputs.
//
// Outbox.Queue/cancel: Queue returns an error for a cancelled context.
//
// Outbox.Push/smoke: Push returns without a panic for derived inputs.
//
// Outbox.Push/smoke#generated: Push returns without a panic for derived
// inputs.
//
// Outbox.Push/cancel: Push returns an error for a cancelled context.
//
// Outbox.Push/cancel#generated: Push returns an error for a cancelled
// context.
//
// Outbox.Flush/smoke: Flush returns without a panic for derived inputs.
//
// Outbox.Flush/cancel: Flush returns an error for a cancelled context.
//
// Outbox.Pour/smoke: Pour returns without a panic for derived inputs.
//
// Outbox.Pour/smoke#generated: Pour returns without a panic for derived
// inputs.
//
// Outbox.Pour/cancel: Pour returns an error for a cancelled context.
//
// Outbox.Pour/cancel#generated: Pour returns an error for a cancelled
// context.
//
// Outbox.Drain/smoke: Drain returns without a panic for derived inputs.

// The derived checks of Outbox that witness their laws at the shallow end:
//
// Outbox.Append/smoke
//
// Outbox.Subscribe/smoke
//
// Outbox.Note/smoke
//
// Outbox.Notes/smoke
//
// Outbox.Log/smoke
//
// Outbox.Logged/smoke
//
// Outbox.Post/smoke
//
// Outbox.Queue/smoke
//
// Outbox.Push/smoke
//
// Outbox.Flush/smoke
//
// Outbox.Pour/smoke
//
// Outbox.Drain/smoke

// The refused claims of Outbox:
//
// Outbox.Subscribe/zero-on-error: eidos derives no sample of a result of
// Subscribe. The remedy is to author a sample of the result type with the
// directive sample.
//
// Outbox.Notes/zero-on-error: eidos derives no sample of a result of
// Notes. The remedy is to author a sample of the result type with the
// directive sample.
//
// Outbox.Queue/zero-on-error: eidos derives no sample of a result of
// Queue. The remedy is to author a sample of the result type with the
// directive sample.
//
// Outbox.Drain/cancel: Drain cannot fail. The remedy is to give Drain an
// error result.
//
// Outbox.Note/delivers-each-append: Notes returns no channel of a value
// that Note passes. The remedy is to subscribe through a callable that
// returns a channel of the values that Note passes.
//
// Outbox.Log/delivers-each-append: Logged returns no channel of a value
// that Log passes. The remedy is to subscribe through a callable that
// returns a channel of the values that Log passes.
//
// Outbox.Post/delivers-each-append: Queue returns no channel of a value
// that Post passes. The remedy is to subscribe through a callable that
// returns a channel of the values that Post passes.
//
// Outbox.Push/delivers-each-append: Flush returns no channel of a value
// that Push passes. The remedy is to subscribe through a callable that
// returns a channel of the values that Push passes.
//
// Outbox.Pour/delivers-each-append: Drain returns no channel of a value
// that Pour passes. The remedy is to subscribe through a callable that
// returns a channel of the values that Pour passes.

// The unclaimed classifications of Outbox:
//
// The shape:writer of Outbox.Append has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Outbox.Note has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Outbox.Log has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Outbox.Post has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Outbox.Push has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Outbox.Flush has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The shape:writer of Outbox.Pour has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Outbox writes.
type (
	OutboxHandle = conform.Handle[Outbox]
	OutboxCheck  = conform.Check[Outbox]
	OutboxOption = conform.Option[Outbox]
)

// OutboxChecks contains the name of each callable of Outbox, with the
// identities of its checks.
var OutboxChecks = struct {
	Append    outboxAppend
	Subscribe outboxSubscribe
	Note      outboxNote
	Notes     outboxNotes
	Log       outboxLog
	Logged    outboxLogged
	Post      outboxPost
	Queue     outboxQueue
	Push      outboxPush
	Flush     outboxFlush
	Pour      outboxPour
	Drain     outboxDrain
}{
	Append:    outboxAppend{"Outbox.Append"},
	Subscribe: outboxSubscribe{"Outbox.Subscribe"},
	Note:      outboxNote{"Outbox.Note"},
	Notes:     outboxNotes{"Outbox.Notes"},
	Log:       outboxLog{"Outbox.Log"},
	Logged:    outboxLogged{"Outbox.Logged"},
	Post:      outboxPost{"Outbox.Post"},
	Queue:     outboxQueue{"Outbox.Queue"},
	Push:      outboxPush{"Outbox.Push"},
	Flush:     outboxFlush{"Outbox.Flush"},
	Pour:      outboxPour{"Outbox.Pour"},
	Drain:     outboxDrain{"Outbox.Drain"},
}

// outboxAppend is the name of Outbox.Append, with the identities of its
// checks.
type outboxAppend struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c outboxAppend) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Append returns
// without a panic for derived inputs.
func (c outboxAppend) SmokeGenerated() conform.ID[Outbox] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c outboxAppend) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Append returns an
// error for a cancelled context.
func (c outboxAppend) CancelGenerated() conform.ID[Outbox] {
	return c.Check("cancel#generated")
}

// DeliversEachAppend returns the identity of the check that each message
// that Append records arrives on the channel of Subscribe.
func (c outboxAppend) DeliversEachAppend() conform.ID[Outbox] {
	return c.Check("delivers-each-append")
}

// outboxSubscribe is the name of Outbox.Subscribe, with the identities of
// its checks.
type outboxSubscribe struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Subscribe returns without a
// panic for derived inputs.
func (c outboxSubscribe) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Subscribe returns an error
// for a cancelled context.
func (c outboxSubscribe) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// outboxNote is the name of Outbox.Note, with the identities of its
// checks.
type outboxNote struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Note returns without a
// panic for derived inputs.
func (c outboxNote) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Note returns
// without a panic for derived inputs.
func (c outboxNote) SmokeGenerated() conform.ID[Outbox] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Note returns an error for
// a cancelled context.
func (c outboxNote) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Note returns an
// error for a cancelled context.
func (c outboxNote) CancelGenerated() conform.ID[Outbox] {
	return c.Check("cancel#generated")
}

// outboxNotes is the name of Outbox.Notes, with the identities of its
// checks.
type outboxNotes struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Notes returns without a
// panic for derived inputs.
func (c outboxNotes) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Notes returns an error for
// a cancelled context.
func (c outboxNotes) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// outboxLog is the name of Outbox.Log, with the identities of its checks.
type outboxLog struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Log returns without a panic
// for derived inputs.
func (c outboxLog) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Log returns
// without a panic for derived inputs.
func (c outboxLog) SmokeGenerated() conform.ID[Outbox] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Log returns an error for a
// cancelled context.
func (c outboxLog) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Log returns an
// error for a cancelled context.
func (c outboxLog) CancelGenerated() conform.ID[Outbox] {
	return c.Check("cancel#generated")
}

// outboxLogged is the name of Outbox.Logged, with the identities of its
// checks.
type outboxLogged struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Logged returns without a
// panic for derived inputs.
func (c outboxLogged) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Logged returns an error
// for a cancelled context.
func (c outboxLogged) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Logged returns the
// zero value beside an error.
func (c outboxLogged) ZeroOnError() conform.ID[Outbox] {
	return c.Check("zero-on-error")
}

// outboxPost is the name of Outbox.Post, with the identities of its
// checks.
type outboxPost struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Post returns without a
// panic for derived inputs.
func (c outboxPost) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Post returns
// without a panic for derived inputs.
func (c outboxPost) SmokeGenerated() conform.ID[Outbox] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Post returns an error for
// a cancelled context.
func (c outboxPost) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Post returns an
// error for a cancelled context.
func (c outboxPost) CancelGenerated() conform.ID[Outbox] {
	return c.Check("cancel#generated")
}

// outboxQueue is the name of Outbox.Queue, with the identities of its
// checks.
type outboxQueue struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Queue returns without a
// panic for derived inputs.
func (c outboxQueue) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Queue returns an error for
// a cancelled context.
func (c outboxQueue) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// outboxPush is the name of Outbox.Push, with the identities of its
// checks.
type outboxPush struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Push returns without a
// panic for derived inputs.
func (c outboxPush) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Push returns
// without a panic for derived inputs.
func (c outboxPush) SmokeGenerated() conform.ID[Outbox] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Push returns an error for
// a cancelled context.
func (c outboxPush) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Push returns an
// error for a cancelled context.
func (c outboxPush) CancelGenerated() conform.ID[Outbox] {
	return c.Check("cancel#generated")
}

// outboxFlush is the name of Outbox.Flush, with the identities of its
// checks.
type outboxFlush struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Flush returns without a
// panic for derived inputs.
func (c outboxFlush) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Flush returns an error for
// a cancelled context.
func (c outboxFlush) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// outboxPour is the name of Outbox.Pour, with the identities of its
// checks.
type outboxPour struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Pour returns without a
// panic for derived inputs.
func (c outboxPour) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Pour returns
// without a panic for derived inputs.
func (c outboxPour) SmokeGenerated() conform.ID[Outbox] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Pour returns an error for
// a cancelled context.
func (c outboxPour) Cancel() conform.ID[Outbox] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Pour returns an
// error for a cancelled context.
func (c outboxPour) CancelGenerated() conform.ID[Outbox] {
	return c.Check("cancel#generated")
}

// outboxDrain is the name of Outbox.Drain, with the identities of its
// checks.
type outboxDrain struct {
	conform.Callable[Outbox]
}

// Smoke returns the identity of the check that Drain returns without a
// panic for derived inputs.
func (c outboxDrain) Smoke() conform.ID[Outbox] {
	return c.Check("smoke")
}

// OutboxSuite returns the checks of Outbox, and the calls and the budgets
// of its callables.
func OutboxSuite() conform.Suite[Outbox] {
	return conform.Suite[Outbox]{
		Declaration: "Outbox",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Outbox]{
			{
				ID:          OutboxChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxAppendPanics,
			},
			{
				ID:          OutboxChecks.Append.SmokeGenerated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnE)
					})
				},
				Defect: outboxAppendPanics,
			},
			{
				ID:          OutboxChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxAppendIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Append.CancelGenerated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnE)
					})
				},
				Defect: outboxAppendIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Append.DeliversEachAppend(),
				Statement:   "each message that Append records arrives on the channel of Subscribe",
				Claim:       spec.DeliversEachAppend,
				Probe:       spec.AppendsArrive,
				Reason:      spec.LostMessage,
				Law:         spec.EachAppendIsDelivered,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *OutboxHandle, c *prop.Case) {
					probe.AppendsArrive(
						h,
						c,
						outboxActions,
						"Append",
						func(ctx context.Context) (<-chan Event, error) {
							r0, err := h.Subject().Subscribe(ctx)
							return r0, err
						},
						func(args []any) Event {
							v, _ := args[0].(Event)
							return v
						},
					)
				},
				Defect: outboxAppendDropsAppendsAfterThree,
				Around: synctest.Test,
			},
			{
				ID:          OutboxChecks.Subscribe.Smoke(),
				Statement:   "Subscribe returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Subscribe(ctx)
						return err
					})
				},
				Defect: outboxSubscribePanics,
			},
			{
				ID:          OutboxChecks.Subscribe.Cancel(),
				Statement:   "Subscribe returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Subscribe(ctx)
						return err
					})
				},
				Defect: outboxSubscribeIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Note.Smoke(),
				Statement:   "Note returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Note(ctx, string("text-a"))
					})
				},
				Defect: outboxNotePanics,
			},
			{
				ID:          OutboxChecks.Note.SmokeGenerated(),
				Statement:   "Note returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnText := c.Draw(prop.Of[string](), "text")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Note(ctx, drawnText)
					})
				},
				Defect: outboxNotePanics,
			},
			{
				ID:          OutboxChecks.Note.Cancel(),
				Statement:   "Note returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Note(ctx, string("text-a"))
					})
				},
				Defect: outboxNoteIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Note.CancelGenerated(),
				Statement:   "Note returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnText := c.Draw(prop.Of[string](), "text")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Note(ctx, drawnText)
					})
				},
				Defect: outboxNoteIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Notes.Smoke(),
				Statement:   "Notes returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Notes(ctx)
						return err
					})
				},
				Defect: outboxNotesPanics,
			},
			{
				ID:          OutboxChecks.Notes.Cancel(),
				Statement:   "Notes returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Notes(ctx)
						return err
					})
				},
				Defect: outboxNotesIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Log.Smoke(),
				Statement:   "Log returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Log(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxLogPanics,
			},
			{
				ID:          OutboxChecks.Log.SmokeGenerated(),
				Statement:   "Log returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Log(ctx, drawnE)
					})
				},
				Defect: outboxLogPanics,
			},
			{
				ID:          OutboxChecks.Log.Cancel(),
				Statement:   "Log returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Log(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxLogIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Log.CancelGenerated(),
				Statement:   "Log returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Log(ctx, drawnE)
					})
				},
				Defect: outboxLogIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Logged.Smoke(),
				Statement:   "Logged returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Logged(ctx)
						return err
					})
				},
				Defect: outboxLoggedPanics,
			},
			{
				ID:          OutboxChecks.Logged.Cancel(),
				Statement:   "Logged returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Logged(ctx)
						return err
					})
				},
				Defect: outboxLoggedIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Logged.ZeroOnError(),
				Statement:   "Logged returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]Event, error) {
						return h.Subject().Logged(ctx)
					}, []Event(nil))
				},
				Defect: outboxLoggedValueBeside,
			},
			{
				ID:          OutboxChecks.Post.Smoke(),
				Statement:   "Post returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxPostPanics,
			},
			{
				ID:          OutboxChecks.Post.SmokeGenerated(),
				Statement:   "Post returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, drawnE)
					})
				},
				Defect: outboxPostPanics,
			},
			{
				ID:          OutboxChecks.Post.Cancel(),
				Statement:   "Post returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxPostIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Post.CancelGenerated(),
				Statement:   "Post returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, drawnE)
					})
				},
				Defect: outboxPostIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Queue.Smoke(),
				Statement:   "Queue returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Queue(ctx)
						return err
					})
				},
				Defect: outboxQueuePanics,
			},
			{
				ID:          OutboxChecks.Queue.Cancel(),
				Statement:   "Queue returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Queue(ctx)
						return err
					})
				},
				Defect: outboxQueueIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Push.Smoke(),
				Statement:   "Push returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxPushPanics,
			},
			{
				ID:          OutboxChecks.Push.SmokeGenerated(),
				Statement:   "Push returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, drawnE)
					})
				},
				Defect: outboxPushPanics,
			},
			{
				ID:          OutboxChecks.Push.Cancel(),
				Statement:   "Push returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxPushIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Push.CancelGenerated(),
				Statement:   "Push returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, drawnE)
					})
				},
				Defect: outboxPushIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Flush.Smoke(),
				Statement:   "Flush returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Flush(ctx)
					})
				},
				Defect: outboxFlushPanics,
			},
			{
				ID:          OutboxChecks.Flush.Cancel(),
				Statement:   "Flush returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Flush(ctx)
					})
				},
				Defect: outboxFlushIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Pour.Smoke(),
				Statement:   "Pour returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pour(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxPourPanics,
			},
			{
				ID:          OutboxChecks.Pour.SmokeGenerated(),
				Statement:   "Pour returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pour(ctx, drawnE)
					})
				},
				Defect: outboxPourPanics,
			},
			{
				ID:          OutboxChecks.Pour.Cancel(),
				Statement:   "Pour returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OutboxHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pour(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: outboxPourIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Pour.CancelGenerated(),
				Statement:   "Pour returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *OutboxHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pour(ctx, drawnE)
					})
				},
				Defect: outboxPourIgnoresCancel,
			},
			{
				ID:          OutboxChecks.Drain.Smoke(),
				Statement:   "Drain returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OutboxHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Drain(ctx)
						return nil
					})
				},
				Defect: outboxDrainPanics,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ZeroOnError,
				Callable: "Outbox.Subscribe",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Subscribe",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Outbox.Notes",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Notes",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Outbox.Queue",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Queue",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Outbox.Drain",
				Gate:     spec.Failure,
				Why:      "Drain cannot fail",
				Remedy:   "give Drain an error result",
			},
			{
				Claim:    spec.DeliversEachAppend,
				Callable: "Outbox.Note",
				Gate:     spec.DeliversOnAChannel,
				Why:      "Notes returns no channel of a value that Note passes",
				Remedy:   "subscribe through a callable that returns a channel of the values that Note passes",
			},
			{
				Claim:    spec.DeliversEachAppend,
				Callable: "Outbox.Log",
				Gate:     spec.DeliversOnAChannel,
				Why:      "Logged returns no channel of a value that Log passes",
				Remedy:   "subscribe through a callable that returns a channel of the values that Log passes",
			},
			{
				Claim:    spec.DeliversEachAppend,
				Callable: "Outbox.Post",
				Gate:     spec.DeliversOnAChannel,
				Why:      "Queue returns no channel of a value that Post passes",
				Remedy:   "subscribe through a callable that returns a channel of the values that Post passes",
			},
			{
				Claim:    spec.DeliversEachAppend,
				Callable: "Outbox.Push",
				Gate:     spec.DeliversOnAChannel,
				Why:      "Flush returns no channel of a value that Push passes",
				Remedy:   "subscribe through a callable that returns a channel of the values that Push passes",
			},
			{
				Claim:    spec.DeliversEachAppend,
				Callable: "Outbox.Pour",
				Gate:     spec.DeliversOnAChannel,
				Why:      "Drain returns no channel of a value that Pour passes",
				Remedy:   "subscribe through a callable that returns a channel of the values that Pour passes",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Outbox.Append",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Outbox.Note",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Outbox.Log",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Outbox.Post",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Outbox.Push",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Outbox.Flush",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Outbox.Pour",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Outbox]{
			{
				Calls:    outboxAppendCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Append.Callable,
			},
			{
				Calls:    outboxSubscribeCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Subscribe.Callable,
			},
			{
				Calls:    outboxNoteCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Note.Callable,
			},
			{
				Calls:    outboxNotesCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Notes.Callable,
			},
			{
				Calls:    outboxLogCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Log.Callable,
			},
			{
				Calls:    outboxLoggedCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Logged.Callable,
			},
			{
				Calls:    outboxPostCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Post.Callable,
			},
			{
				Calls:    outboxQueueCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Queue.Callable,
			},
			{
				Calls:    outboxPushCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Push.Callable,
			},
			{
				Calls:    outboxFlushCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Flush.Callable,
			},
			{
				Calls:    outboxPourCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Pour.Callable,
			},
			{
				Calls:    outboxDrainCalls,
				Over:     newOutboxOver,
				Callable: OutboxChecks.Drain.Callable,
			},
		},
	}
}

// RunOutboxSuite applies the checks of Outbox, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunOutboxSuite[T Outbox](t *testing.T, newSubject func() T, opts ...OutboxOption) {
	t.Helper()
	RunOutboxSuiteFromHandle(t, func(*OutboxHandle) T { return newSubject() }, opts...)
}

// RunOutboxSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunOutboxSuiteFromHandle[T Outbox](
	t *testing.T,
	newSubject func(h *OutboxHandle) T,
	opts ...OutboxOption,
) {
	t.Helper()
	conform.Run(t, OutboxSuite(), outboxSubjects(newSubject), opts...)
}

// FuzzOutboxSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzOutboxSuite[T Outbox](
	f *testing.F,
	callable conform.Callable[Outbox],
	newSubject func(h *OutboxHandle) T,
	opts ...OutboxOption,
) {
	f.Helper()
	conform.Fuzz(f, OutboxSuite(), callable, outboxSubjects(newSubject), opts...)
}

// RunOutboxBudgets applies the allocation ceilings of Outbox and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunOutboxBudgets[T Outbox](
	t *testing.T,
	newSubject func(h *OutboxHandle) T,
	opts ...OutboxOption,
) {
	t.Helper()
	conform.RunBudgets(t, OutboxSuite(), outboxSubjects(newSubject), opts...)
}

// BenchOutboxSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchOutboxSuite[T Outbox](
	b conform.B,
	callable conform.Callable[Outbox],
	newSubject func(h *OutboxHandle) T,
	opts ...OutboxOption,
) {
	b.Helper()
	conform.Bench(b, OutboxSuite(), callable, outboxSubjects(newSubject), opts...)
}

// outboxSubjects returns the factory of the subjects of Outbox that calls
// newSubject with the handle of each check.
func outboxSubjects[T Outbox](newSubject func(h *OutboxHandle) T) func(h *OutboxHandle) Outbox {
	return func(h *OutboxHandle) Outbox { return newSubject(h) }
}

// outboxAppendCalls generates the calls of Outbox.Append with the inputs
// that each case draws.
var outboxAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Outbox) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Outbox) error {
		return subject.Append(ctx, e)
	}
})

// outboxSubscribeCalls generates the calls of Outbox.Subscribe with the
// inputs that each case draws.
var outboxSubscribeCalls = prop.Composite(func(*prop.Case) func(context.Context, Outbox) error {
	return func(ctx context.Context, subject Outbox) error {
		_, err := subject.Subscribe(ctx)
		return err
	}
})

// outboxNoteCalls generates the calls of Outbox.Note with the inputs that
// each case draws.
var outboxNoteCalls = prop.Composite(func(c *prop.Case) func(context.Context, Outbox) error {
	text := c.Draw(prop.Of[string](), "text")
	return func(ctx context.Context, subject Outbox) error {
		return subject.Note(ctx, text)
	}
})

// outboxNotesCalls generates the calls of Outbox.Notes with the inputs
// that each case draws.
var outboxNotesCalls = prop.Composite(func(*prop.Case) func(context.Context, Outbox) error {
	return func(ctx context.Context, subject Outbox) error {
		_, err := subject.Notes(ctx)
		return err
	}
})

// outboxLogCalls generates the calls of Outbox.Log with the inputs that
// each case draws.
var outboxLogCalls = prop.Composite(func(c *prop.Case) func(context.Context, Outbox) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Outbox) error {
		return subject.Log(ctx, e)
	}
})

// outboxLoggedCalls generates the calls of Outbox.Logged with the inputs
// that each case draws.
var outboxLoggedCalls = prop.Composite(func(*prop.Case) func(context.Context, Outbox) error {
	return func(ctx context.Context, subject Outbox) error {
		_, err := subject.Logged(ctx)
		return err
	}
})

// outboxPostCalls generates the calls of Outbox.Post with the inputs that
// each case draws.
var outboxPostCalls = prop.Composite(func(c *prop.Case) func(context.Context, Outbox) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Outbox) error {
		return subject.Post(ctx, e)
	}
})

// outboxQueueCalls generates the calls of Outbox.Queue with the inputs
// that each case draws.
var outboxQueueCalls = prop.Composite(func(*prop.Case) func(context.Context, Outbox) error {
	return func(ctx context.Context, subject Outbox) error {
		_, err := subject.Queue(ctx)
		return err
	}
})

// outboxPushCalls generates the calls of Outbox.Push with the inputs that
// each case draws.
var outboxPushCalls = prop.Composite(func(c *prop.Case) func(context.Context, Outbox) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Outbox) error {
		return subject.Push(ctx, e)
	}
})

// outboxFlushCalls generates the calls of Outbox.Flush with the inputs
// that each case draws.
var outboxFlushCalls = prop.Composite(func(*prop.Case) func(context.Context, Outbox) error {
	return func(ctx context.Context, subject Outbox) error {
		return subject.Flush(ctx)
	}
})

// outboxPourCalls generates the calls of Outbox.Pour with the inputs that
// each case draws.
var outboxPourCalls = prop.Composite(func(c *prop.Case) func(context.Context, Outbox) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Outbox) error {
		return subject.Pour(ctx, e)
	}
})

// outboxDrainCalls generates the calls of Outbox.Drain with the inputs
// that each case draws.
var outboxDrainCalls = prop.Composite(func(*prop.Case) func(context.Context, Outbox) error {
	return func(ctx context.Context, subject Outbox) error {
		subject.Drain(ctx)
		return nil
	}
})

// outboxOverEmbedded is Outbox under the name of the field that outboxOver
// embeds.
type outboxOverEmbedded = Outbox

// outboxOver is a subject of Outbox whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type outboxOver struct {
	outboxOverEmbedded
	spend func()
}

// Append calls spend and returns its zero results.
func (o outboxOver) Append(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Subscribe calls spend and returns its zero results.
func (o outboxOver) Subscribe(context.Context) (r0 <-chan Event, r1 error) {
	o.spend()
	return r0, r1
}

// Note calls spend and returns its zero results.
func (o outboxOver) Note(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Notes calls spend and returns its zero results.
func (o outboxOver) Notes(context.Context) (r0 <-chan Event, r1 error) {
	o.spend()
	return r0, r1
}

// Log calls spend and returns its zero results.
func (o outboxOver) Log(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Logged calls spend and returns its zero results.
func (o outboxOver) Logged(context.Context) (r0 []Event, r1 error) {
	o.spend()
	return r0, r1
}

// Post calls spend and returns its zero results.
func (o outboxOver) Post(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Queue calls spend and returns its zero results.
func (o outboxOver) Queue(context.Context) (r0 chan<- Event, r1 error) {
	o.spend()
	return r0, r1
}

// Push calls spend and returns its zero results.
func (o outboxOver) Push(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Flush calls spend and returns its zero results.
func (o outboxOver) Flush(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Pour calls spend and returns its zero results.
func (o outboxOver) Pour(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Drain calls spend and returns its zero results.
func (o outboxOver) Drain(context.Context) {
	o.spend()
}

// newOutboxOver returns the planted defect of a ceiling of Outbox, whose
// callables call spend.
func newOutboxOver(_ *OutboxHandle, spend func()) Outbox {
	return outboxOver{spend: spend}
}

// outboxAppendPanicsDefect is the planted defect of a Outbox whose Append
// panics.
type outboxAppendPanicsDefect struct {
	outboxOverEmbedded
}

// Append is the override of the planted defect.
func (*outboxAppendPanicsDefect) Append(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Outbox whose Append panics")
}

// outboxAppendPanics returns the planted defect of a Outbox whose Append
// panics over a fresh subject.
func outboxAppendPanics(h *OutboxHandle) Outbox {
	return &outboxAppendPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxAppendIgnoresCancelDefect is the planted defect of a Outbox whose
// Append ignores the context.
type outboxAppendIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Append is the override of the planted defect.
func (d *outboxAppendIgnoresCancelDefect) Append(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.outboxOverEmbedded.Append(p0, p1)
}

// outboxAppendIgnoresCancel returns the planted defect of a Outbox whose
// Append ignores the context over a fresh subject.
func outboxAppendIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxAppendIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxAppendDropsAppendsAfterThreeDefect is the planted defect of a
// Outbox whose Append records nothing from its fourth call on.
type outboxAppendDropsAppendsAfterThreeDefect struct {
	outboxOverEmbedded
	calls int
}

// Append is the override of the planted defect.
func (d *outboxAppendDropsAppendsAfterThreeDefect) Append(p0 context.Context, p1 Event) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		return r0
	}
	return d.outboxOverEmbedded.Append(p0, p1)
}

// outboxAppendDropsAppendsAfterThree returns the planted defect of a
// Outbox whose Append records nothing from its fourth call on over a fresh
// subject.
func outboxAppendDropsAppendsAfterThree(h *OutboxHandle) Outbox {
	return &outboxAppendDropsAppendsAfterThreeDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxSubscribePanicsDefect is the planted defect of a Outbox whose
// Subscribe panics.
type outboxSubscribePanicsDefect struct {
	outboxOverEmbedded
}

// Subscribe is the override of the planted defect.
func (*outboxSubscribePanicsDefect) Subscribe(p0 context.Context) (r0 <-chan Event, r1 error) {
	panic("the planted defect of a Outbox whose Subscribe panics")
}

// outboxSubscribePanics returns the planted defect of a Outbox whose
// Subscribe panics over a fresh subject.
func outboxSubscribePanics(h *OutboxHandle) Outbox {
	return &outboxSubscribePanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxSubscribeIgnoresCancelDefect is the planted defect of a Outbox
// whose Subscribe ignores the context.
type outboxSubscribeIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Subscribe is the override of the planted defect.
func (d *outboxSubscribeIgnoresCancelDefect) Subscribe(
	p0 context.Context,
) (r0 <-chan Event, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.outboxOverEmbedded.Subscribe(p0)
}

// outboxSubscribeIgnoresCancel returns the planted defect of a Outbox
// whose Subscribe ignores the context over a fresh subject.
func outboxSubscribeIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxSubscribeIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxNotePanicsDefect is the planted defect of a Outbox whose Note
// panics.
type outboxNotePanicsDefect struct {
	outboxOverEmbedded
}

// Note is the override of the planted defect.
func (*outboxNotePanicsDefect) Note(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Outbox whose Note panics")
}

// outboxNotePanics returns the planted defect of a Outbox whose Note
// panics over a fresh subject.
func outboxNotePanics(h *OutboxHandle) Outbox {
	return &outboxNotePanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxNoteIgnoresCancelDefect is the planted defect of a Outbox whose
// Note ignores the context.
type outboxNoteIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Note is the override of the planted defect.
func (d *outboxNoteIgnoresCancelDefect) Note(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.outboxOverEmbedded.Note(p0, p1)
}

// outboxNoteIgnoresCancel returns the planted defect of a Outbox whose
// Note ignores the context over a fresh subject.
func outboxNoteIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxNoteIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxNotesPanicsDefect is the planted defect of a Outbox whose Notes
// panics.
type outboxNotesPanicsDefect struct {
	outboxOverEmbedded
}

// Notes is the override of the planted defect.
func (*outboxNotesPanicsDefect) Notes(p0 context.Context) (r0 <-chan Event, r1 error) {
	panic("the planted defect of a Outbox whose Notes panics")
}

// outboxNotesPanics returns the planted defect of a Outbox whose Notes
// panics over a fresh subject.
func outboxNotesPanics(h *OutboxHandle) Outbox {
	return &outboxNotesPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxNotesIgnoresCancelDefect is the planted defect of a Outbox whose
// Notes ignores the context.
type outboxNotesIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Notes is the override of the planted defect.
func (d *outboxNotesIgnoresCancelDefect) Notes(p0 context.Context) (r0 <-chan Event, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.outboxOverEmbedded.Notes(p0)
}

// outboxNotesIgnoresCancel returns the planted defect of a Outbox whose
// Notes ignores the context over a fresh subject.
func outboxNotesIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxNotesIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxLogPanicsDefect is the planted defect of a Outbox whose Log
// panics.
type outboxLogPanicsDefect struct {
	outboxOverEmbedded
}

// Log is the override of the planted defect.
func (*outboxLogPanicsDefect) Log(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Outbox whose Log panics")
}

// outboxLogPanics returns the planted defect of a Outbox whose Log panics
// over a fresh subject.
func outboxLogPanics(h *OutboxHandle) Outbox {
	return &outboxLogPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxLogIgnoresCancelDefect is the planted defect of a Outbox whose Log
// ignores the context.
type outboxLogIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Log is the override of the planted defect.
func (d *outboxLogIgnoresCancelDefect) Log(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.outboxOverEmbedded.Log(p0, p1)
}

// outboxLogIgnoresCancel returns the planted defect of a Outbox whose Log
// ignores the context over a fresh subject.
func outboxLogIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxLogIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxLoggedPanicsDefect is the planted defect of a Outbox whose Logged
// panics.
type outboxLoggedPanicsDefect struct {
	outboxOverEmbedded
}

// Logged is the override of the planted defect.
func (*outboxLoggedPanicsDefect) Logged(p0 context.Context) (r0 []Event, r1 error) {
	panic("the planted defect of a Outbox whose Logged panics")
}

// outboxLoggedPanics returns the planted defect of a Outbox whose Logged
// panics over a fresh subject.
func outboxLoggedPanics(h *OutboxHandle) Outbox {
	return &outboxLoggedPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxLoggedIgnoresCancelDefect is the planted defect of a Outbox whose
// Logged ignores the context.
type outboxLoggedIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Logged is the override of the planted defect.
func (d *outboxLoggedIgnoresCancelDefect) Logged(p0 context.Context) (r0 []Event, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.outboxOverEmbedded.Logged(p0)
}

// outboxLoggedIgnoresCancel returns the planted defect of a Outbox whose
// Logged ignores the context over a fresh subject.
func outboxLoggedIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxLoggedIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxLoggedValueBesideDefect is the planted defect of a Outbox whose
// Logged returns a value beside an error.
type outboxLoggedValueBesideDefect struct {
	outboxOverEmbedded
}

// Logged is the override of the planted defect.
func (d *outboxLoggedValueBesideDefect) Logged(p0 context.Context) (r0 []Event, r1 error) {
	r0, r1 = d.outboxOverEmbedded.Logged(p0)
	if r1 != nil {
		r0 = []Event{{Key: "Key-b"}}
	}
	return r0, r1
}

// outboxLoggedValueBeside returns the planted defect of a Outbox whose
// Logged returns a value beside an error over a fresh subject.
func outboxLoggedValueBeside(h *OutboxHandle) Outbox {
	return &outboxLoggedValueBesideDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxPostPanicsDefect is the planted defect of a Outbox whose Post
// panics.
type outboxPostPanicsDefect struct {
	outboxOverEmbedded
}

// Post is the override of the planted defect.
func (*outboxPostPanicsDefect) Post(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Outbox whose Post panics")
}

// outboxPostPanics returns the planted defect of a Outbox whose Post
// panics over a fresh subject.
func outboxPostPanics(h *OutboxHandle) Outbox {
	return &outboxPostPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxPostIgnoresCancelDefect is the planted defect of a Outbox whose
// Post ignores the context.
type outboxPostIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Post is the override of the planted defect.
func (d *outboxPostIgnoresCancelDefect) Post(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.outboxOverEmbedded.Post(p0, p1)
}

// outboxPostIgnoresCancel returns the planted defect of a Outbox whose
// Post ignores the context over a fresh subject.
func outboxPostIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxPostIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxQueuePanicsDefect is the planted defect of a Outbox whose Queue
// panics.
type outboxQueuePanicsDefect struct {
	outboxOverEmbedded
}

// Queue is the override of the planted defect.
func (*outboxQueuePanicsDefect) Queue(p0 context.Context) (r0 chan<- Event, r1 error) {
	panic("the planted defect of a Outbox whose Queue panics")
}

// outboxQueuePanics returns the planted defect of a Outbox whose Queue
// panics over a fresh subject.
func outboxQueuePanics(h *OutboxHandle) Outbox {
	return &outboxQueuePanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxQueueIgnoresCancelDefect is the planted defect of a Outbox whose
// Queue ignores the context.
type outboxQueueIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Queue is the override of the planted defect.
func (d *outboxQueueIgnoresCancelDefect) Queue(p0 context.Context) (r0 chan<- Event, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.outboxOverEmbedded.Queue(p0)
}

// outboxQueueIgnoresCancel returns the planted defect of a Outbox whose
// Queue ignores the context over a fresh subject.
func outboxQueueIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxQueueIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxPushPanicsDefect is the planted defect of a Outbox whose Push
// panics.
type outboxPushPanicsDefect struct {
	outboxOverEmbedded
}

// Push is the override of the planted defect.
func (*outboxPushPanicsDefect) Push(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Outbox whose Push panics")
}

// outboxPushPanics returns the planted defect of a Outbox whose Push
// panics over a fresh subject.
func outboxPushPanics(h *OutboxHandle) Outbox {
	return &outboxPushPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxPushIgnoresCancelDefect is the planted defect of a Outbox whose
// Push ignores the context.
type outboxPushIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Push is the override of the planted defect.
func (d *outboxPushIgnoresCancelDefect) Push(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.outboxOverEmbedded.Push(p0, p1)
}

// outboxPushIgnoresCancel returns the planted defect of a Outbox whose
// Push ignores the context over a fresh subject.
func outboxPushIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxPushIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxFlushPanicsDefect is the planted defect of a Outbox whose Flush
// panics.
type outboxFlushPanicsDefect struct {
	outboxOverEmbedded
}

// Flush is the override of the planted defect.
func (*outboxFlushPanicsDefect) Flush(p0 context.Context) (r0 error) {
	panic("the planted defect of a Outbox whose Flush panics")
}

// outboxFlushPanics returns the planted defect of a Outbox whose Flush
// panics over a fresh subject.
func outboxFlushPanics(h *OutboxHandle) Outbox {
	return &outboxFlushPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxFlushIgnoresCancelDefect is the planted defect of a Outbox whose
// Flush ignores the context.
type outboxFlushIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Flush is the override of the planted defect.
func (d *outboxFlushIgnoresCancelDefect) Flush(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.outboxOverEmbedded.Flush(p0)
}

// outboxFlushIgnoresCancel returns the planted defect of a Outbox whose
// Flush ignores the context over a fresh subject.
func outboxFlushIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxFlushIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxPourPanicsDefect is the planted defect of a Outbox whose Pour
// panics.
type outboxPourPanicsDefect struct {
	outboxOverEmbedded
}

// Pour is the override of the planted defect.
func (*outboxPourPanicsDefect) Pour(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Outbox whose Pour panics")
}

// outboxPourPanics returns the planted defect of a Outbox whose Pour
// panics over a fresh subject.
func outboxPourPanics(h *OutboxHandle) Outbox {
	return &outboxPourPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxPourIgnoresCancelDefect is the planted defect of a Outbox whose
// Pour ignores the context.
type outboxPourIgnoresCancelDefect struct {
	outboxOverEmbedded
}

// Pour is the override of the planted defect.
func (d *outboxPourIgnoresCancelDefect) Pour(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.outboxOverEmbedded.Pour(p0, p1)
}

// outboxPourIgnoresCancel returns the planted defect of a Outbox whose
// Pour ignores the context over a fresh subject.
func outboxPourIgnoresCancel(h *OutboxHandle) Outbox {
	return &outboxPourIgnoresCancelDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxDrainPanicsDefect is the planted defect of a Outbox whose Drain
// panics.
type outboxDrainPanicsDefect struct {
	outboxOverEmbedded
}

// Drain is the override of the planted defect.
func (*outboxDrainPanicsDefect) Drain(p0 context.Context) {
	panic("the planted defect of a Outbox whose Drain panics")
}

// outboxDrainPanics returns the planted defect of a Outbox whose Drain
// panics over a fresh subject.
func outboxDrainPanics(h *OutboxHandle) Outbox {
	return &outboxDrainPanicsDefect{outboxOverEmbedded: h.Fresh()}
}

// outboxActions returns the actions of the generated sequences of Outbox
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func outboxActions(c *prop.Case, subject func() Outbox) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Append",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "append.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Append", args)
				values, err := outboxAppendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Subscribe",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Subscribe", args)
				values, err := outboxSubscribeCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Note",
			Input: func(c *prop.Case, _ struct{}) any {
				text := c.Draw(prop.Of[string](), "note.text")
				return []any{text}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Note", args)
				values, err := outboxNoteCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Notes",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Notes", args)
				values, err := outboxNotesCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Log",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "log.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Log", args)
				values, err := outboxLogCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Logged",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Logged", args)
				values, err := outboxLoggedCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Post",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "post.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Post", args)
				values, err := outboxPostCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Queue",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Queue", args)
				values, err := outboxQueueCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Push",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "push.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Push", args)
				values, err := outboxPushCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Flush",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Flush", args)
				values, err := outboxFlushCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Pour",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "pour.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Pour", args)
				values, err := outboxPourCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Drain",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Drain", args)
				values, err := outboxDrainCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// outboxAppendCall calls Append on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxAppendCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Append(ctx, a0)
}

// outboxSubscribeCall calls Subscribe on subject with the arguments args,
// and returns the values of the call beside its error.
func outboxSubscribeCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	r0, err := subject.Subscribe(ctx)
	return []any{r0}, err
}

// outboxNoteCall calls Note on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxNoteCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Note(ctx, a0)
}

// outboxNotesCall calls Notes on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxNotesCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	r0, err := subject.Notes(ctx)
	return []any{r0}, err
}

// outboxLogCall calls Log on subject with the arguments args, and returns
// the values of the call beside its error.
func outboxLogCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Log(ctx, a0)
}

// outboxLoggedCall calls Logged on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxLoggedCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	r0, err := subject.Logged(ctx)
	return []any{r0}, err
}

// outboxPostCall calls Post on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxPostCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Post(ctx, a0)
}

// outboxQueueCall calls Queue on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxQueueCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	r0, err := subject.Queue(ctx)
	return []any{r0}, err
}

// outboxPushCall calls Push on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxPushCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Push(ctx, a0)
}

// outboxFlushCall calls Flush on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxFlushCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	return nil, subject.Flush(ctx)
}

// outboxPourCall calls Pour on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxPourCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Pour(ctx, a0)
}

// outboxDrainCall calls Drain on subject with the arguments args, and
// returns the values of the call beside its error.
func outboxDrainCall(ctx context.Context, subject Outbox, args []any) ([]any, error) {
	subject.Drain(ctx)
	return nil, nil
}

// The derived checks of Broker:
//
// Broker.Send/smoke: Send returns without a panic for derived inputs.
//
// Broker.Send/smoke#generated: Send returns without a panic for derived
// inputs.
//
// Broker.Send/cancel: Send returns an error for a cancelled context.
//
// Broker.Send/cancel#generated: Send returns an error for a cancelled
// context.
//
// Broker.Send/delivers-by-its-mode: each event that Send publishes arrives
// on the channel of Listen as its delivery mode states.
//
// Broker.Resend/smoke: Resend returns without a panic for derived inputs.
//
// Broker.Resend/smoke#generated: Resend returns without a panic for
// derived inputs.
//
// Broker.Resend/cancel: Resend returns an error for a cancelled context.
//
// Broker.Resend/cancel#generated: Resend returns an error for a cancelled
// context.
//
// Broker.Resend/redelivers-on-demand: each event that Resend redelivers
// arrives on the channel of Listen once more.
//
// Broker.Listen/smoke: Listen returns without a panic for derived inputs.
//
// Broker.Listen/cancel: Listen returns an error for a cancelled context.
//
// Broker.Fire/smoke: Fire returns without a panic for derived inputs.
//
// Broker.Fire/smoke#generated: Fire returns without a panic for derived
// inputs.
//
// Broker.Fire/cancel: Fire returns an error for a cancelled context.
//
// Broker.Fire/cancel#generated: Fire returns an error for a cancelled
// context.
//
// Broker.Fire/delivers-by-its-mode: each event that Fire publishes arrives
// on the channel of Tap as its delivery mode states.
//
// Broker.Refire/smoke: Refire returns without a panic for derived inputs.
//
// Broker.Refire/smoke#generated: Refire returns without a panic for
// derived inputs.
//
// Broker.Refire/cancel: Refire returns an error for a cancelled context.
//
// Broker.Refire/cancel#generated: Refire returns an error for a cancelled
// context.
//
// Broker.Tap/smoke: Tap returns without a panic for derived inputs.

// The derived checks of Broker that witness their laws at the shallow end:
//
// Broker.Send/smoke
//
// Broker.Resend/smoke
//
// Broker.Listen/smoke
//
// Broker.Fire/smoke
//
// Broker.Refire/smoke
//
// Broker.Tap/smoke

// The refused claims of Broker:
//
// Broker.Listen/zero-on-error: eidos derives no sample of a result of
// Listen. The remedy is to author a sample of the result type with the
// directive sample.
//
// Broker.Tap/cancel: Tap cannot fail. The remedy is to give Tap an error
// result.
//
// Broker.Tap/zero-on-error: Tap cannot fail. The remedy is to give Tap an
// error result.
//
// Broker.Refire/redelivers-on-demand: Fire states no delivery mode
// at-least-once, so a redelivery of Refire need not deliver an event
// again. The remedy is to state the delivery mode of Fire, whose check
// judges the redeliveries of Refire under any mode.

// The unclaimed classifications of Broker:
//
// The shape:writer of Broker.Send has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Broker.Resend has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Broker.Fire has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Broker.Refire has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Broker writes.
type (
	BrokerHandle = conform.Handle[Broker]
	BrokerCheck  = conform.Check[Broker]
	BrokerOption = conform.Option[Broker]
)

// BrokerChecks contains the name of each callable of Broker, with the
// identities of its checks.
var BrokerChecks = struct {
	Send   brokerSend
	Resend brokerResend
	Listen brokerListen
	Fire   brokerFire
	Refire brokerRefire
	Tap    brokerTap
}{
	Send:   brokerSend{"Broker.Send"},
	Resend: brokerResend{"Broker.Resend"},
	Listen: brokerListen{"Broker.Listen"},
	Fire:   brokerFire{"Broker.Fire"},
	Refire: brokerRefire{"Broker.Refire"},
	Tap:    brokerTap{"Broker.Tap"},
}

// brokerSend is the name of Broker.Send, with the identities of its
// checks.
type brokerSend struct {
	conform.Callable[Broker]
}

// Smoke returns the identity of the check that Send returns without a
// panic for derived inputs.
func (c brokerSend) Smoke() conform.ID[Broker] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Send returns
// without a panic for derived inputs.
func (c brokerSend) SmokeGenerated() conform.ID[Broker] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Send returns an error for
// a cancelled context.
func (c brokerSend) Cancel() conform.ID[Broker] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Send returns an
// error for a cancelled context.
func (c brokerSend) CancelGenerated() conform.ID[Broker] {
	return c.Check("cancel#generated")
}

// DeliversByItsMode returns the identity of the check that each event that
// Send publishes arrives on the channel of Listen as its delivery mode
// states.
func (c brokerSend) DeliversByItsMode() conform.ID[Broker] {
	return c.Check("delivers-by-its-mode")
}

// brokerResend is the name of Broker.Resend, with the identities of its
// checks.
type brokerResend struct {
	conform.Callable[Broker]
}

// Smoke returns the identity of the check that Resend returns without a
// panic for derived inputs.
func (c brokerResend) Smoke() conform.ID[Broker] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Resend returns
// without a panic for derived inputs.
func (c brokerResend) SmokeGenerated() conform.ID[Broker] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Resend returns an error
// for a cancelled context.
func (c brokerResend) Cancel() conform.ID[Broker] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Resend returns an
// error for a cancelled context.
func (c brokerResend) CancelGenerated() conform.ID[Broker] {
	return c.Check("cancel#generated")
}

// RedeliversOnDemand returns the identity of the check that each event
// that Resend redelivers arrives on the channel of Listen once more.
func (c brokerResend) RedeliversOnDemand() conform.ID[Broker] {
	return c.Check("redelivers-on-demand")
}

// brokerListen is the name of Broker.Listen, with the identities of its
// checks.
type brokerListen struct {
	conform.Callable[Broker]
}

// Smoke returns the identity of the check that Listen returns without a
// panic for derived inputs.
func (c brokerListen) Smoke() conform.ID[Broker] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Listen returns an error
// for a cancelled context.
func (c brokerListen) Cancel() conform.ID[Broker] {
	return c.Check("cancel")
}

// brokerFire is the name of Broker.Fire, with the identities of its
// checks.
type brokerFire struct {
	conform.Callable[Broker]
}

// Smoke returns the identity of the check that Fire returns without a
// panic for derived inputs.
func (c brokerFire) Smoke() conform.ID[Broker] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Fire returns
// without a panic for derived inputs.
func (c brokerFire) SmokeGenerated() conform.ID[Broker] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Fire returns an error for
// a cancelled context.
func (c brokerFire) Cancel() conform.ID[Broker] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Fire returns an
// error for a cancelled context.
func (c brokerFire) CancelGenerated() conform.ID[Broker] {
	return c.Check("cancel#generated")
}

// DeliversByItsMode returns the identity of the check that each event that
// Fire publishes arrives on the channel of Tap as its delivery mode
// states.
func (c brokerFire) DeliversByItsMode() conform.ID[Broker] {
	return c.Check("delivers-by-its-mode")
}

// brokerRefire is the name of Broker.Refire, with the identities of its
// checks.
type brokerRefire struct {
	conform.Callable[Broker]
}

// Smoke returns the identity of the check that Refire returns without a
// panic for derived inputs.
func (c brokerRefire) Smoke() conform.ID[Broker] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Refire returns
// without a panic for derived inputs.
func (c brokerRefire) SmokeGenerated() conform.ID[Broker] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Refire returns an error
// for a cancelled context.
func (c brokerRefire) Cancel() conform.ID[Broker] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Refire returns an
// error for a cancelled context.
func (c brokerRefire) CancelGenerated() conform.ID[Broker] {
	return c.Check("cancel#generated")
}

// brokerTap is the name of Broker.Tap, with the identities of its checks.
type brokerTap struct {
	conform.Callable[Broker]
}

// Smoke returns the identity of the check that Tap returns without a panic
// for derived inputs.
func (c brokerTap) Smoke() conform.ID[Broker] {
	return c.Check("smoke")
}

// BrokerSuite returns the checks of Broker, and the calls and the budgets
// of its callables.
func BrokerSuite() conform.Suite[Broker] {
	return conform.Suite[Broker]{
		Declaration: "Broker",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Broker]{
			{
				ID:          BrokerChecks.Send.Smoke(),
				Statement:   "Send returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BrokerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Send(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerSendPanics,
			},
			{
				ID:          BrokerChecks.Send.SmokeGenerated(),
				Statement:   "Send returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Send(ctx, drawnE)
					})
				},
				Defect: brokerSendPanics,
			},
			{
				ID:          BrokerChecks.Send.Cancel(),
				Statement:   "Send returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BrokerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Send(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerSendIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Send.CancelGenerated(),
				Statement:   "Send returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Send(ctx, drawnE)
					})
				},
				Defect: brokerSendIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Send.DeliversByItsMode(),
				Statement:   "each event that Send publishes arrives on the channel of Listen as its delivery mode states",
				Claim:       spec.DeliversByItsMode,
				Probe:       spec.DeliveriesMatchTheMode,
				Reason:      spec.WrongDeliveryCount,
				Law:         spec.DeliveriesFollowTheMode,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *BrokerHandle, c *prop.Case) {
					probe.DeliveriesMatchTheMode(
						h,
						c,
						brokerActions,
						"Send",
						func(ctx context.Context) (<-chan Event, error) {
							r0, err := h.Subject().Listen(ctx)
							return r0, err
						},
						func(args []any) Event {
							v, _ := args[0].(Event)
							return v
						},
						"at-least-once",
					)
				},
				Defect: brokerSendBreaksTheModeAfterThree,
				Around: synctest.Test,
			},
			{
				ID:          BrokerChecks.Resend.Smoke(),
				Statement:   "Resend returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BrokerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Resend(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerResendPanics,
			},
			{
				ID:          BrokerChecks.Resend.SmokeGenerated(),
				Statement:   "Resend returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Resend(ctx, drawnE)
					})
				},
				Defect: brokerResendPanics,
			},
			{
				ID:          BrokerChecks.Resend.Cancel(),
				Statement:   "Resend returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BrokerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Resend(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerResendIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Resend.CancelGenerated(),
				Statement:   "Resend returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Resend(ctx, drawnE)
					})
				},
				Defect: brokerResendIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Resend.RedeliversOnDemand(),
				Statement:   "each event that Resend redelivers arrives on the channel of Listen once more",
				Claim:       spec.RedeliversOnDemand,
				Probe:       spec.RedeliveriesArrive,
				Reason:      spec.NotRedelivered,
				Law:         spec.RedeliveryDeliversAgain,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *BrokerHandle, c *prop.Case) {
					probe.RedeliveriesArrive(
						h,
						c,
						brokerActions,
						"Send",
						"Resend",
						func(ctx context.Context) (<-chan Event, error) {
							r0, err := h.Subject().Listen(ctx)
							return r0, err
						},
						func(args []any) Event {
							v, _ := args[0].(Event)
							return v
						},
						func(args []any) Event {
							v, _ := args[0].(Event)
							return v
						},
					)
				},
				Defect: brokerResendSkipsARedeliveryAfterThree,
				Around: synctest.Test,
			},
			{
				ID:          BrokerChecks.Listen.Smoke(),
				Statement:   "Listen returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BrokerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Listen(ctx)
						return err
					})
				},
				Defect: brokerListenPanics,
			},
			{
				ID:          BrokerChecks.Listen.Cancel(),
				Statement:   "Listen returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BrokerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Listen(ctx)
						return err
					})
				},
				Defect: brokerListenIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Fire.Smoke(),
				Statement:   "Fire returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BrokerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Fire(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerFirePanics,
			},
			{
				ID:          BrokerChecks.Fire.SmokeGenerated(),
				Statement:   "Fire returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Fire(ctx, drawnE)
					})
				},
				Defect: brokerFirePanics,
			},
			{
				ID:          BrokerChecks.Fire.Cancel(),
				Statement:   "Fire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BrokerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Fire(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerFireIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Fire.CancelGenerated(),
				Statement:   "Fire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Fire(ctx, drawnE)
					})
				},
				Defect: brokerFireIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Fire.DeliversByItsMode(),
				Statement:   "each event that Fire publishes arrives on the channel of Tap as its delivery mode states",
				Claim:       spec.DeliversByItsMode,
				Probe:       spec.DeliveriesMatchTheMode,
				Reason:      spec.WrongDeliveryCount,
				Law:         spec.DeliveriesFollowTheMode,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *BrokerHandle, c *prop.Case) {
					probe.DeliveriesMatchTheMode(
						h,
						c,
						brokerActions,
						"Fire",
						func(ctx context.Context) (<-chan Event, error) {
							r0 := h.Subject().Tap(ctx)
							return r0, nil
						},
						func(args []any) Event {
							v, _ := args[0].(Event)
							return v
						},
						"at-most-once",
					)
				},
				Defect: brokerFireBreaksTheModeAfterThree,
				Around: synctest.Test,
			},
			{
				ID:          BrokerChecks.Refire.Smoke(),
				Statement:   "Refire returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BrokerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Refire(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerRefirePanics,
			},
			{
				ID:          BrokerChecks.Refire.SmokeGenerated(),
				Statement:   "Refire returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Refire(ctx, drawnE)
					})
				},
				Defect: brokerRefirePanics,
			},
			{
				ID:          BrokerChecks.Refire.Cancel(),
				Statement:   "Refire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BrokerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Refire(ctx, Event{Key: "Key-a"})
					})
				},
				Defect: brokerRefireIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Refire.CancelGenerated(),
				Statement:   "Refire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BrokerHandle, c *prop.Case) {
					drawnE := c.Draw(prop.Of[Event](), "e")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Refire(ctx, drawnE)
					})
				},
				Defect: brokerRefireIgnoresCancel,
			},
			{
				ID:          BrokerChecks.Tap.Smoke(),
				Statement:   "Tap returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BrokerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Tap(ctx)
						return nil
					})
				},
				Defect: brokerTapPanics,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ZeroOnError,
				Callable: "Broker.Listen",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Listen",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Broker.Tap",
				Gate:     spec.Failure,
				Why:      "Tap cannot fail",
				Remedy:   "give Tap an error result",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Broker.Tap",
				Gate:     spec.Failure,
				Why:      "Tap cannot fail",
				Remedy:   "give Tap an error result",
			},
			{
				Claim:    spec.RedeliversOnDemand,
				Callable: "Broker.Refire",
				Gate:     spec.AtLeastOnceMode,
				Why:      "Fire states no delivery mode at-least-once, so a redelivery of Refire need not deliver an event again",
				Remedy:   "state the delivery mode of Fire, whose check judges the redeliveries of Refire under any mode",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Broker.Send",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Broker.Resend",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Broker.Fire",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Broker.Refire",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Broker]{
			{
				Calls:    brokerSendCalls,
				Over:     newBrokerOver,
				Callable: BrokerChecks.Send.Callable,
			},
			{
				Calls:    brokerResendCalls,
				Over:     newBrokerOver,
				Callable: BrokerChecks.Resend.Callable,
			},
			{
				Calls:    brokerListenCalls,
				Over:     newBrokerOver,
				Callable: BrokerChecks.Listen.Callable,
			},
			{
				Calls:    brokerFireCalls,
				Over:     newBrokerOver,
				Callable: BrokerChecks.Fire.Callable,
			},
			{
				Calls:    brokerRefireCalls,
				Over:     newBrokerOver,
				Callable: BrokerChecks.Refire.Callable,
			},
			{
				Calls:    brokerTapCalls,
				Over:     newBrokerOver,
				Callable: BrokerChecks.Tap.Callable,
			},
		},
	}
}

// RunBrokerSuite applies the checks of Broker, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunBrokerSuite[T Broker](t *testing.T, newSubject func() T, opts ...BrokerOption) {
	t.Helper()
	RunBrokerSuiteFromHandle(t, func(*BrokerHandle) T { return newSubject() }, opts...)
}

// RunBrokerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunBrokerSuiteFromHandle[T Broker](
	t *testing.T,
	newSubject func(h *BrokerHandle) T,
	opts ...BrokerOption,
) {
	t.Helper()
	conform.Run(t, BrokerSuite(), brokerSubjects(newSubject), opts...)
}

// FuzzBrokerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzBrokerSuite[T Broker](
	f *testing.F,
	callable conform.Callable[Broker],
	newSubject func(h *BrokerHandle) T,
	opts ...BrokerOption,
) {
	f.Helper()
	conform.Fuzz(f, BrokerSuite(), callable, brokerSubjects(newSubject), opts...)
}

// RunBrokerBudgets applies the allocation ceilings of Broker and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunBrokerBudgets[T Broker](
	t *testing.T,
	newSubject func(h *BrokerHandle) T,
	opts ...BrokerOption,
) {
	t.Helper()
	conform.RunBudgets(t, BrokerSuite(), brokerSubjects(newSubject), opts...)
}

// BenchBrokerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchBrokerSuite[T Broker](
	b conform.B,
	callable conform.Callable[Broker],
	newSubject func(h *BrokerHandle) T,
	opts ...BrokerOption,
) {
	b.Helper()
	conform.Bench(b, BrokerSuite(), callable, brokerSubjects(newSubject), opts...)
}

// brokerSubjects returns the factory of the subjects of Broker that calls
// newSubject with the handle of each check.
func brokerSubjects[T Broker](newSubject func(h *BrokerHandle) T) func(h *BrokerHandle) Broker {
	return func(h *BrokerHandle) Broker { return newSubject(h) }
}

// brokerSendCalls generates the calls of Broker.Send with the inputs that
// each case draws.
var brokerSendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Broker) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Broker) error {
		return subject.Send(ctx, e)
	}
})

// brokerResendCalls generates the calls of Broker.Resend with the inputs
// that each case draws.
var brokerResendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Broker) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Broker) error {
		return subject.Resend(ctx, e)
	}
})

// brokerListenCalls generates the calls of Broker.Listen with the inputs
// that each case draws.
var brokerListenCalls = prop.Composite(func(*prop.Case) func(context.Context, Broker) error {
	return func(ctx context.Context, subject Broker) error {
		_, err := subject.Listen(ctx)
		return err
	}
})

// brokerFireCalls generates the calls of Broker.Fire with the inputs that
// each case draws.
var brokerFireCalls = prop.Composite(func(c *prop.Case) func(context.Context, Broker) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Broker) error {
		return subject.Fire(ctx, e)
	}
})

// brokerRefireCalls generates the calls of Broker.Refire with the inputs
// that each case draws.
var brokerRefireCalls = prop.Composite(func(c *prop.Case) func(context.Context, Broker) error {
	e := c.Draw(prop.Of[Event](), "e")
	return func(ctx context.Context, subject Broker) error {
		return subject.Refire(ctx, e)
	}
})

// brokerTapCalls generates the calls of Broker.Tap with the inputs that
// each case draws.
var brokerTapCalls = prop.Composite(func(*prop.Case) func(context.Context, Broker) error {
	return func(ctx context.Context, subject Broker) error {
		_ = subject.Tap(ctx)
		return nil
	}
})

// brokerOverEmbedded is Broker under the name of the field that brokerOver
// embeds.
type brokerOverEmbedded = Broker

// brokerOver is a subject of Broker whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type brokerOver struct {
	brokerOverEmbedded
	spend func()
}

// Send calls spend and returns its zero results.
func (o brokerOver) Send(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Resend calls spend and returns its zero results.
func (o brokerOver) Resend(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Listen calls spend and returns its zero results.
func (o brokerOver) Listen(context.Context) (r0 chan Event, r1 error) {
	o.spend()
	return r0, r1
}

// Fire calls spend and returns its zero results.
func (o brokerOver) Fire(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Refire calls spend and returns its zero results.
func (o brokerOver) Refire(context.Context, Event) (r0 error) {
	o.spend()
	return r0
}

// Tap calls spend and returns its zero results.
func (o brokerOver) Tap(context.Context) (r0 <-chan Event) {
	o.spend()
	return r0
}

// newBrokerOver returns the planted defect of a ceiling of Broker, whose
// callables call spend.
func newBrokerOver(_ *BrokerHandle, spend func()) Broker {
	return brokerOver{spend: spend}
}

// brokerSendPanicsDefect is the planted defect of a Broker whose Send
// panics.
type brokerSendPanicsDefect struct {
	brokerOverEmbedded
}

// Send is the override of the planted defect.
func (*brokerSendPanicsDefect) Send(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Broker whose Send panics")
}

// brokerSendPanics returns the planted defect of a Broker whose Send
// panics over a fresh subject.
func brokerSendPanics(h *BrokerHandle) Broker {
	return &brokerSendPanicsDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerSendIgnoresCancelDefect is the planted defect of a Broker whose
// Send ignores the context.
type brokerSendIgnoresCancelDefect struct {
	brokerOverEmbedded
}

// Send is the override of the planted defect.
func (d *brokerSendIgnoresCancelDefect) Send(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.brokerOverEmbedded.Send(p0, p1)
}

// brokerSendIgnoresCancel returns the planted defect of a Broker whose
// Send ignores the context over a fresh subject.
func brokerSendIgnoresCancel(h *BrokerHandle) Broker {
	return &brokerSendIgnoresCancelDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerSendBreaksTheModeAfterThreeDefect is the planted defect of a
// Broker whose Send breaks its delivery mode from its fourth call on.
type brokerSendBreaksTheModeAfterThreeDefect struct {
	brokerOverEmbedded
	calls int
}

// Send is the override of the planted defect.
func (d *brokerSendBreaksTheModeAfterThreeDefect) Send(p0 context.Context, p1 Event) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		return r0
	}
	return d.brokerOverEmbedded.Send(p0, p1)
}

// brokerSendBreaksTheModeAfterThree returns the planted defect of a Broker
// whose Send breaks its delivery mode from its fourth call on over a fresh
// subject.
func brokerSendBreaksTheModeAfterThree(h *BrokerHandle) Broker {
	return &brokerSendBreaksTheModeAfterThreeDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerResendPanicsDefect is the planted defect of a Broker whose Resend
// panics.
type brokerResendPanicsDefect struct {
	brokerOverEmbedded
}

// Resend is the override of the planted defect.
func (*brokerResendPanicsDefect) Resend(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Broker whose Resend panics")
}

// brokerResendPanics returns the planted defect of a Broker whose Resend
// panics over a fresh subject.
func brokerResendPanics(h *BrokerHandle) Broker {
	return &brokerResendPanicsDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerResendIgnoresCancelDefect is the planted defect of a Broker whose
// Resend ignores the context.
type brokerResendIgnoresCancelDefect struct {
	brokerOverEmbedded
}

// Resend is the override of the planted defect.
func (d *brokerResendIgnoresCancelDefect) Resend(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.brokerOverEmbedded.Resend(p0, p1)
}

// brokerResendIgnoresCancel returns the planted defect of a Broker whose
// Resend ignores the context over a fresh subject.
func brokerResendIgnoresCancel(h *BrokerHandle) Broker {
	return &brokerResendIgnoresCancelDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerResendSkipsARedeliveryAfterThreeDefect is the planted defect of a
// Broker whose Resend does nothing from its fourth call on.
type brokerResendSkipsARedeliveryAfterThreeDefect struct {
	brokerOverEmbedded
	calls int
}

// Resend is the override of the planted defect.
func (d *brokerResendSkipsARedeliveryAfterThreeDefect) Resend(
	p0 context.Context,
	p1 Event,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		return r0
	}
	return d.brokerOverEmbedded.Resend(p0, p1)
}

// brokerResendSkipsARedeliveryAfterThree returns the planted defect of a
// Broker whose Resend does nothing from its fourth call on over a fresh
// subject.
func brokerResendSkipsARedeliveryAfterThree(h *BrokerHandle) Broker {
	return &brokerResendSkipsARedeliveryAfterThreeDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerListenPanicsDefect is the planted defect of a Broker whose Listen
// panics.
type brokerListenPanicsDefect struct {
	brokerOverEmbedded
}

// Listen is the override of the planted defect.
func (*brokerListenPanicsDefect) Listen(p0 context.Context) (r0 chan Event, r1 error) {
	panic("the planted defect of a Broker whose Listen panics")
}

// brokerListenPanics returns the planted defect of a Broker whose Listen
// panics over a fresh subject.
func brokerListenPanics(h *BrokerHandle) Broker {
	return &brokerListenPanicsDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerListenIgnoresCancelDefect is the planted defect of a Broker whose
// Listen ignores the context.
type brokerListenIgnoresCancelDefect struct {
	brokerOverEmbedded
}

// Listen is the override of the planted defect.
func (d *brokerListenIgnoresCancelDefect) Listen(p0 context.Context) (r0 chan Event, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.brokerOverEmbedded.Listen(p0)
}

// brokerListenIgnoresCancel returns the planted defect of a Broker whose
// Listen ignores the context over a fresh subject.
func brokerListenIgnoresCancel(h *BrokerHandle) Broker {
	return &brokerListenIgnoresCancelDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerFirePanicsDefect is the planted defect of a Broker whose Fire
// panics.
type brokerFirePanicsDefect struct {
	brokerOverEmbedded
}

// Fire is the override of the planted defect.
func (*brokerFirePanicsDefect) Fire(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Broker whose Fire panics")
}

// brokerFirePanics returns the planted defect of a Broker whose Fire
// panics over a fresh subject.
func brokerFirePanics(h *BrokerHandle) Broker {
	return &brokerFirePanicsDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerFireIgnoresCancelDefect is the planted defect of a Broker whose
// Fire ignores the context.
type brokerFireIgnoresCancelDefect struct {
	brokerOverEmbedded
}

// Fire is the override of the planted defect.
func (d *brokerFireIgnoresCancelDefect) Fire(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.brokerOverEmbedded.Fire(p0, p1)
}

// brokerFireIgnoresCancel returns the planted defect of a Broker whose
// Fire ignores the context over a fresh subject.
func brokerFireIgnoresCancel(h *BrokerHandle) Broker {
	return &brokerFireIgnoresCancelDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerFireBreaksTheModeAfterThreeDefect is the planted defect of a
// Broker whose Fire breaks its delivery mode from its fourth call on.
type brokerFireBreaksTheModeAfterThreeDefect struct {
	brokerOverEmbedded
	calls int
}

// Fire is the override of the planted defect.
func (d *brokerFireBreaksTheModeAfterThreeDefect) Fire(p0 context.Context, p1 Event) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		r0 = d.brokerOverEmbedded.Fire(p0, p1)
	}
	return d.brokerOverEmbedded.Fire(p0, p1)
}

// brokerFireBreaksTheModeAfterThree returns the planted defect of a Broker
// whose Fire breaks its delivery mode from its fourth call on over a fresh
// subject.
func brokerFireBreaksTheModeAfterThree(h *BrokerHandle) Broker {
	return &brokerFireBreaksTheModeAfterThreeDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerRefirePanicsDefect is the planted defect of a Broker whose Refire
// panics.
type brokerRefirePanicsDefect struct {
	brokerOverEmbedded
}

// Refire is the override of the planted defect.
func (*brokerRefirePanicsDefect) Refire(p0 context.Context, p1 Event) (r0 error) {
	panic("the planted defect of a Broker whose Refire panics")
}

// brokerRefirePanics returns the planted defect of a Broker whose Refire
// panics over a fresh subject.
func brokerRefirePanics(h *BrokerHandle) Broker {
	return &brokerRefirePanicsDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerRefireIgnoresCancelDefect is the planted defect of a Broker whose
// Refire ignores the context.
type brokerRefireIgnoresCancelDefect struct {
	brokerOverEmbedded
}

// Refire is the override of the planted defect.
func (d *brokerRefireIgnoresCancelDefect) Refire(p0 context.Context, p1 Event) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.brokerOverEmbedded.Refire(p0, p1)
}

// brokerRefireIgnoresCancel returns the planted defect of a Broker whose
// Refire ignores the context over a fresh subject.
func brokerRefireIgnoresCancel(h *BrokerHandle) Broker {
	return &brokerRefireIgnoresCancelDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerTapPanicsDefect is the planted defect of a Broker whose Tap
// panics.
type brokerTapPanicsDefect struct {
	brokerOverEmbedded
}

// Tap is the override of the planted defect.
func (*brokerTapPanicsDefect) Tap(p0 context.Context) (r0 <-chan Event) {
	panic("the planted defect of a Broker whose Tap panics")
}

// brokerTapPanics returns the planted defect of a Broker whose Tap panics
// over a fresh subject.
func brokerTapPanics(h *BrokerHandle) Broker {
	return &brokerTapPanicsDefect{brokerOverEmbedded: h.Fresh()}
}

// brokerActions returns the actions of the generated sequences of Broker
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func brokerActions(c *prop.Case, subject func() Broker) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Send",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "send.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Send", args)
				values, err := brokerSendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Resend",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "resend.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Resend", args)
				values, err := brokerResendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Listen",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Listen", args)
				values, err := brokerListenCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Fire",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "fire.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Fire", args)
				values, err := brokerFireCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Refire",
			Input: func(c *prop.Case, _ struct{}) any {
				e := c.Draw(prop.Of[Event](), "refire.e")
				return []any{e}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Refire", args)
				values, err := brokerRefireCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Tap",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Tap", args)
				values, err := brokerTapCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// brokerSendCall calls Send on subject with the arguments args, and
// returns the values of the call beside its error.
func brokerSendCall(ctx context.Context, subject Broker, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Send(ctx, a0)
}

// brokerResendCall calls Resend on subject with the arguments args, and
// returns the values of the call beside its error.
func brokerResendCall(ctx context.Context, subject Broker, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Resend(ctx, a0)
}

// brokerListenCall calls Listen on subject with the arguments args, and
// returns the values of the call beside its error.
func brokerListenCall(ctx context.Context, subject Broker, args []any) ([]any, error) {
	r0, err := subject.Listen(ctx)
	return []any{r0}, err
}

// brokerFireCall calls Fire on subject with the arguments args, and
// returns the values of the call beside its error.
func brokerFireCall(ctx context.Context, subject Broker, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Fire(ctx, a0)
}

// brokerRefireCall calls Refire on subject with the arguments args, and
// returns the values of the call beside its error.
func brokerRefireCall(ctx context.Context, subject Broker, args []any) ([]any, error) {
	a0, _ := args[0].(Event)
	return nil, subject.Refire(ctx, a0)
}

// brokerTapCall calls Tap on subject with the arguments args, and returns
// the values of the call beside its error.
func brokerTapCall(ctx context.Context, subject Broker, args []any) ([]any, error) {
	r0 := subject.Tap(ctx)
	return []any{r0}, nil
}

// The derived checks of Watches:
//
// Watches.Watch/smoke: Watch returns without a panic for derived inputs.
//
// Watches.Watch/smoke#generated: Watch returns without a panic for derived
// inputs.
//
// Watches.Watch/cancel: Watch returns an error for a cancelled context.
//
// Watches.Watch/cancel#generated: Watch returns an error for a cancelled
// context.
//
// Watches.Watch/watch-sees-each-trigger: a subscription of Watch returns
// each update that Set causes for its key, in order, and no update after
// its stop.
//
// Watches.Set/smoke: Set returns without a panic for derived inputs.
//
// Watches.Set/smoke#generated: Set returns without a panic for derived
// inputs.
//
// Watches.Set/cancel: Set returns an error for a cancelled context.
//
// Watches.Set/cancel#generated: Set returns an error for a cancelled
// context.
//
// Watches.Ticks/smoke: Ticks returns without a panic for derived inputs.
//
// Watches.Ticks/smoke#generated: Ticks returns without a panic for derived
// inputs.
//
// Watches.Ticks/cancel: Ticks returns an error for a cancelled context.
//
// Watches.Ticks/cancel#generated: Ticks returns an error for a cancelled
// context.
//
// Watches.Ticks/watch-sees-each-trigger: a subscription of Ticks returns
// each update that Tick causes for its key, in order, and no update after
// its stop.
//
// Watches.Tick/smoke: Tick returns without a panic for derived inputs.
//
// Watches.Tick/smoke#generated: Tick returns without a panic for derived
// inputs.
//
// Watches.Tick/cancel: Tick returns an error for a cancelled context.
//
// Watches.Tick/cancel#generated: Tick returns an error for a cancelled
// context.
//
// Watches.Peeks/smoke: Peeks returns without a panic for derived inputs.
//
// Watches.Peeks/smoke#generated: Peeks returns without a panic for derived
// inputs.
//
// Watches.Peeks/cancel: Peeks returns an error for a cancelled context.
//
// Watches.Peeks/cancel#generated: Peeks returns an error for a cancelled
// context.
//
// Watches.Peeked/smoke: Peeked returns without a panic for derived inputs.
//
// Watches.Peeked/smoke#generated: Peeked returns without a panic for
// derived inputs.
//
// Watches.Peeked/cancel: Peeked returns an error for a cancelled context.
//
// Watches.Peeked/cancel#generated: Peeked returns an error for a cancelled
// context.
//
// Watches.Finds/smoke: Finds returns without a panic for derived inputs.
//
// Watches.Finds/smoke#generated: Finds returns without a panic for derived
// inputs.
//
// Watches.Finds/cancel: Finds returns an error for a cancelled context.
//
// Watches.Finds/cancel#generated: Finds returns an error for a cancelled
// context.
//
// Watches.Found/smoke: Found returns without a panic for derived inputs.
//
// Watches.Found/smoke#generated: Found returns without a panic for derived
// inputs.
//
// Watches.Found/cancel: Found returns an error for a cancelled context.
//
// Watches.Found/cancel#generated: Found returns an error for a cancelled
// context.
//
// Watches.Reads/smoke: Reads returns without a panic for derived inputs.
//
// Watches.Reads/smoke#generated: Reads returns without a panic for derived
// inputs.
//
// Watches.Reads/cancel: Reads returns an error for a cancelled context.
//
// Watches.Reads/cancel#generated: Reads returns an error for a cancelled
// context.
//
// Watches.Written/smoke: Written returns without a panic for derived
// inputs.
//
// Watches.Written/smoke#generated: Written returns without a panic for
// derived inputs.
//
// Watches.Written/cancel: Written returns an error for a cancelled
// context.
//
// Watches.Written/cancel#generated: Written returns an error for a
// cancelled context.
//
// Watches.Cancels/smoke: Cancels returns without a panic for derived
// inputs.
//
// Watches.Cancels/smoke#generated: Cancels returns without a panic for
// derived inputs.
//
// Watches.Cancels/cancel: Cancels returns an error for a cancelled
// context.
//
// Watches.Cancels/cancel#generated: Cancels returns an error for a
// cancelled context.
//
// Watches.Cancelled/smoke: Cancelled returns without a panic for derived
// inputs.
//
// Watches.Cancelled/smoke#generated: Cancelled returns without a panic for
// derived inputs.
//
// Watches.Cancelled/cancel: Cancelled returns an error for a cancelled
// context.
//
// Watches.Cancelled/cancel#generated: Cancelled returns an error for a
// cancelled context.
//
// Watches.Closes/smoke: Closes returns without a panic for derived inputs.
//
// Watches.Closes/smoke#generated: Closes returns without a panic for
// derived inputs.
//
// Watches.Closes/cancel: Closes returns an error for a cancelled context.
//
// Watches.Closes/cancel#generated: Closes returns an error for a cancelled
// context.
//
// Watches.Closed/smoke: Closed returns without a panic for derived inputs.
//
// Watches.Closed/smoke#generated: Closed returns without a panic for
// derived inputs.
//
// Watches.Closed/cancel: Closed returns an error for a cancelled context.
//
// Watches.Closed/cancel#generated: Closed returns an error for a cancelled
// context.
//
// Watches.Fields/smoke: Fields returns without a panic for derived inputs.
//
// Watches.Fields/smoke#generated: Fields returns without a panic for
// derived inputs.
//
// Watches.Fields/cancel: Fields returns an error for a cancelled context.
//
// Watches.Fields/cancel#generated: Fields returns an error for a cancelled
// context.
//
// Watches.Filled/smoke: Filled returns without a panic for derived inputs.
//
// Watches.Filled/smoke#generated: Filled returns without a panic for
// derived inputs.
//
// Watches.Filled/cancel: Filled returns an error for a cancelled context.
//
// Watches.Filled/cancel#generated: Filled returns an error for a cancelled
// context.
//
// Watches.Follows/smoke: Follows returns without a panic for derived
// inputs.
//
// Watches.Follows/smoke#generated: Follows returns without a panic for
// derived inputs.
//
// Watches.Follows/cancel: Follows returns an error for a cancelled
// context.
//
// Watches.Follows/cancel#generated: Follows returns an error for a
// cancelled context.
//
// Watches.Touch/smoke: Touch returns without a panic for derived inputs.
//
// Watches.Touch/smoke#generated: Touch returns without a panic for derived
// inputs.
//
// Watches.Touch/cancel: Touch returns an error for a cancelled context.
//
// Watches.Touch/cancel#generated: Touch returns an error for a cancelled
// context.

// The derived checks of Watches that witness their laws at the shallow
// end:
//
// Watches.Watch/smoke
//
// Watches.Set/smoke
//
// Watches.Ticks/smoke
//
// Watches.Tick/smoke
//
// Watches.Peeks/smoke
//
// Watches.Peeked/smoke
//
// Watches.Finds/smoke
//
// Watches.Found/smoke
//
// Watches.Reads/smoke
//
// Watches.Written/smoke
//
// Watches.Cancels/smoke
//
// Watches.Cancelled/smoke
//
// Watches.Closes/smoke
//
// Watches.Closed/smoke
//
// Watches.Fields/smoke
//
// Watches.Filled/smoke
//
// Watches.Follows/smoke
//
// Watches.Touch/smoke

// The refused claims of Watches:
//
// Watches.Ticks/zero-on-error: eidos derives no sample of a result of
// Ticks. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Peeks/zero-on-error: eidos derives no sample of a result of
// Peeks. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Finds/zero-on-error: eidos derives no sample of a result of
// Finds. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Reads/zero-on-error: eidos derives no sample of a result of
// Reads. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Cancels/zero-on-error: eidos derives no sample of a result of
// Cancels. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Closes/zero-on-error: eidos derives no sample of a result of
// Closes. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Fields/zero-on-error: eidos derives no sample of a result of
// Fields. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Follows/zero-on-error: eidos derives no sample of a result of
// Follows. The remedy is to author a sample of the result type with the
// directive sample.
//
// Watches.Peeks/watch-sees-each-trigger: the methods that next and stop
// name on the subscription of Peeks, or the value that Peeked passes, do
// not fit a read of updates. The remedy is to name with next a method that
// takes a wait and returns the update and a flag, and with stop a method
// without parameters and results.
//
// Watches.Finds/watch-sees-each-trigger: the methods that next and stop
// name on the subscription of Finds, or the value that Found passes, do
// not fit a read of updates. The remedy is to name with next a method that
// takes a wait and returns the update and a flag, and with stop a method
// without parameters and results.
//
// Watches.Reads/watch-sees-each-trigger: the methods that next and stop
// name on the subscription of Reads, or the value that Written passes, do
// not fit a read of updates. The remedy is to name with next a method that
// takes a wait and returns the update and a flag, and with stop a method
// without parameters and results.
//
// Watches.Cancels/watch-sees-each-trigger: the methods that next and stop
// name on the subscription of Cancels, or the value that Cancelled passes,
// do not fit a read of updates. The remedy is to name with next a method
// that takes a wait and returns the update and a flag, and with stop a
// method without parameters and results.
//
// Watches.Closes/watch-sees-each-trigger: the methods that next and stop
// name on the subscription of Closes, or the value that Closed passes, do
// not fit a read of updates. The remedy is to name with next a method that
// takes a wait and returns the update and a flag, and with stop a method
// without parameters and results.
//
// Watches.Fields/watch-sees-each-trigger: the methods that next and stop
// name on the subscription of Fields, or the value that Filled passes, do
// not fit a read of updates. The remedy is to name with next a method that
// takes a wait and returns the update and a flag, and with stop a method
// without parameters and results.
//
// Watches.Follows/watch-sees-each-trigger: the methods that next and stop
// name on the subscription of Follows, or the value that Touch passes, do
// not fit a read of updates. The remedy is to name with next a method that
// takes a wait and returns the update and a flag, and with stop a method
// without parameters and results.

// The unclaimed classifications of Watches:
//
// The shape:multi-arg-writer of Watches.Set has the stance silent. A write
// is an input of the other checks, and its result is an error alone.
//
// The shape:multi-arg-writer of Watches.Tick has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Watches.Peeked has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Watches.Found has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Watches.Written has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Watches.Cancelled has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Watches.Closed has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Watches.Filled has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:writer of Watches.Touch has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Watches writes.
type (
	WatchesHandle = conform.Handle[Watches]
	WatchesCheck  = conform.Check[Watches]
	WatchesOption = conform.Option[Watches]
)

// WatchesChecks contains the name of each callable of Watches, with the
// identities of its checks.
var WatchesChecks = struct {
	Watch     watchesWatch
	Set       watchesSet
	Ticks     watchesTicks
	Tick      watchesTick
	Peeks     watchesPeeks
	Peeked    watchesPeeked
	Finds     watchesFinds
	Found     watchesFound
	Reads     watchesReads
	Written   watchesWritten
	Cancels   watchesCancels
	Cancelled watchesCancelled
	Closes    watchesCloses
	Closed    watchesClosed
	Fields    watchesFields
	Filled    watchesFilled
	Follows   watchesFollows
	Touch     watchesTouch
}{
	Watch:     watchesWatch{"Watches.Watch"},
	Set:       watchesSet{"Watches.Set"},
	Ticks:     watchesTicks{"Watches.Ticks"},
	Tick:      watchesTick{"Watches.Tick"},
	Peeks:     watchesPeeks{"Watches.Peeks"},
	Peeked:    watchesPeeked{"Watches.Peeked"},
	Finds:     watchesFinds{"Watches.Finds"},
	Found:     watchesFound{"Watches.Found"},
	Reads:     watchesReads{"Watches.Reads"},
	Written:   watchesWritten{"Watches.Written"},
	Cancels:   watchesCancels{"Watches.Cancels"},
	Cancelled: watchesCancelled{"Watches.Cancelled"},
	Closes:    watchesCloses{"Watches.Closes"},
	Closed:    watchesClosed{"Watches.Closed"},
	Fields:    watchesFields{"Watches.Fields"},
	Filled:    watchesFilled{"Watches.Filled"},
	Follows:   watchesFollows{"Watches.Follows"},
	Touch:     watchesTouch{"Watches.Touch"},
}

// watchesWatch is the name of Watches.Watch, with the identities of its
// checks.
type watchesWatch struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Watch returns without a
// panic for derived inputs.
func (c watchesWatch) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Watch returns
// without a panic for derived inputs.
func (c watchesWatch) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Watch returns an error for
// a cancelled context.
func (c watchesWatch) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Watch returns an
// error for a cancelled context.
func (c watchesWatch) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// WatchSeesEachTrigger returns the identity of the check that a
// subscription of Watch returns each update that Set causes for its key,
// in order, and no update after its stop.
func (c watchesWatch) WatchSeesEachTrigger() conform.ID[Watches] {
	return c.Check("watch-sees-each-trigger")
}

// watchesSet is the name of Watches.Set, with the identities of its
// checks.
type watchesSet struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c watchesSet) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Set returns
// without a panic for derived inputs.
func (c watchesSet) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Set returns an error for a
// cancelled context.
func (c watchesSet) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Set returns an
// error for a cancelled context.
func (c watchesSet) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesTicks is the name of Watches.Ticks, with the identities of its
// checks.
type watchesTicks struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Ticks returns without a
// panic for derived inputs.
func (c watchesTicks) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Ticks returns
// without a panic for derived inputs.
func (c watchesTicks) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Ticks returns an error for
// a cancelled context.
func (c watchesTicks) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Ticks returns an
// error for a cancelled context.
func (c watchesTicks) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// WatchSeesEachTrigger returns the identity of the check that a
// subscription of Ticks returns each update that Tick causes for its key,
// in order, and no update after its stop.
func (c watchesTicks) WatchSeesEachTrigger() conform.ID[Watches] {
	return c.Check("watch-sees-each-trigger")
}

// watchesTick is the name of Watches.Tick, with the identities of its
// checks.
type watchesTick struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Tick returns without a
// panic for derived inputs.
func (c watchesTick) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Tick returns
// without a panic for derived inputs.
func (c watchesTick) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Tick returns an error for
// a cancelled context.
func (c watchesTick) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Tick returns an
// error for a cancelled context.
func (c watchesTick) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesPeeks is the name of Watches.Peeks, with the identities of its
// checks.
type watchesPeeks struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Peeks returns without a
// panic for derived inputs.
func (c watchesPeeks) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Peeks returns
// without a panic for derived inputs.
func (c watchesPeeks) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Peeks returns an error for
// a cancelled context.
func (c watchesPeeks) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Peeks returns an
// error for a cancelled context.
func (c watchesPeeks) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesPeeked is the name of Watches.Peeked, with the identities of its
// checks.
type watchesPeeked struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Peeked returns without a
// panic for derived inputs.
func (c watchesPeeked) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Peeked returns
// without a panic for derived inputs.
func (c watchesPeeked) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Peeked returns an error
// for a cancelled context.
func (c watchesPeeked) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Peeked returns an
// error for a cancelled context.
func (c watchesPeeked) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesFinds is the name of Watches.Finds, with the identities of its
// checks.
type watchesFinds struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Finds returns without a
// panic for derived inputs.
func (c watchesFinds) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Finds returns
// without a panic for derived inputs.
func (c watchesFinds) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Finds returns an error for
// a cancelled context.
func (c watchesFinds) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Finds returns an
// error for a cancelled context.
func (c watchesFinds) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesFound is the name of Watches.Found, with the identities of its
// checks.
type watchesFound struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Found returns without a
// panic for derived inputs.
func (c watchesFound) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Found returns
// without a panic for derived inputs.
func (c watchesFound) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Found returns an error for
// a cancelled context.
func (c watchesFound) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Found returns an
// error for a cancelled context.
func (c watchesFound) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesReads is the name of Watches.Reads, with the identities of its
// checks.
type watchesReads struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Reads returns without a
// panic for derived inputs.
func (c watchesReads) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Reads returns
// without a panic for derived inputs.
func (c watchesReads) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Reads returns an error for
// a cancelled context.
func (c watchesReads) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Reads returns an
// error for a cancelled context.
func (c watchesReads) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesWritten is the name of Watches.Written, with the identities of
// its checks.
type watchesWritten struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Written returns without a
// panic for derived inputs.
func (c watchesWritten) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Written returns
// without a panic for derived inputs.
func (c watchesWritten) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Written returns an error
// for a cancelled context.
func (c watchesWritten) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Written returns
// an error for a cancelled context.
func (c watchesWritten) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesCancels is the name of Watches.Cancels, with the identities of
// its checks.
type watchesCancels struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Cancels returns without a
// panic for derived inputs.
func (c watchesCancels) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Cancels returns
// without a panic for derived inputs.
func (c watchesCancels) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Cancels returns an error
// for a cancelled context.
func (c watchesCancels) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Cancels returns
// an error for a cancelled context.
func (c watchesCancels) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesCancelled is the name of Watches.Cancelled, with the identities
// of its checks.
type watchesCancelled struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Cancelled returns without a
// panic for derived inputs.
func (c watchesCancelled) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Cancelled returns
// without a panic for derived inputs.
func (c watchesCancelled) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Cancelled returns an error
// for a cancelled context.
func (c watchesCancelled) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Cancelled returns
// an error for a cancelled context.
func (c watchesCancelled) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesCloses is the name of Watches.Closes, with the identities of its
// checks.
type watchesCloses struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Closes returns without a
// panic for derived inputs.
func (c watchesCloses) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Closes returns
// without a panic for derived inputs.
func (c watchesCloses) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Closes returns an error
// for a cancelled context.
func (c watchesCloses) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Closes returns an
// error for a cancelled context.
func (c watchesCloses) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesClosed is the name of Watches.Closed, with the identities of its
// checks.
type watchesClosed struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Closed returns without a
// panic for derived inputs.
func (c watchesClosed) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Closed returns
// without a panic for derived inputs.
func (c watchesClosed) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Closed returns an error
// for a cancelled context.
func (c watchesClosed) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Closed returns an
// error for a cancelled context.
func (c watchesClosed) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesFields is the name of Watches.Fields, with the identities of its
// checks.
type watchesFields struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Fields returns without a
// panic for derived inputs.
func (c watchesFields) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Fields returns
// without a panic for derived inputs.
func (c watchesFields) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Fields returns an error
// for a cancelled context.
func (c watchesFields) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Fields returns an
// error for a cancelled context.
func (c watchesFields) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesFilled is the name of Watches.Filled, with the identities of its
// checks.
type watchesFilled struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Filled returns without a
// panic for derived inputs.
func (c watchesFilled) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Filled returns
// without a panic for derived inputs.
func (c watchesFilled) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Filled returns an error
// for a cancelled context.
func (c watchesFilled) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Filled returns an
// error for a cancelled context.
func (c watchesFilled) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesFollows is the name of Watches.Follows, with the identities of
// its checks.
type watchesFollows struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Follows returns without a
// panic for derived inputs.
func (c watchesFollows) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Follows returns
// without a panic for derived inputs.
func (c watchesFollows) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Follows returns an error
// for a cancelled context.
func (c watchesFollows) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Follows returns
// an error for a cancelled context.
func (c watchesFollows) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// watchesTouch is the name of Watches.Touch, with the identities of its
// checks.
type watchesTouch struct {
	conform.Callable[Watches]
}

// Smoke returns the identity of the check that Touch returns without a
// panic for derived inputs.
func (c watchesTouch) Smoke() conform.ID[Watches] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Touch returns
// without a panic for derived inputs.
func (c watchesTouch) SmokeGenerated() conform.ID[Watches] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Touch returns an error for
// a cancelled context.
func (c watchesTouch) Cancel() conform.ID[Watches] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Touch returns an
// error for a cancelled context.
func (c watchesTouch) CancelGenerated() conform.ID[Watches] {
	return c.Check("cancel#generated")
}

// WatchesSuite returns the checks of Watches, and the calls and the
// budgets of its callables.
func WatchesSuite() conform.Suite[Watches] {
	return conform.Suite[Watches]{
		Declaration: "Watches",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Watches]{
			{
				ID:          WatchesChecks.Watch.Smoke(),
				Statement:   "Watch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Watch(ctx, string("scope-a"), string("key-a"))
						return err
					})
				},
				Defect: watchesWatchPanics,
			},
			{
				ID:          WatchesChecks.Watch.SmokeGenerated(),
				Statement:   "Watch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnScope := c.Draw(prop.Of[string](), "scope")
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Watch(ctx, drawnScope, drawnKey)
						return err
					})
				},
				Defect: watchesWatchPanics,
			},
			{
				ID:          WatchesChecks.Watch.Cancel(),
				Statement:   "Watch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Watch(ctx, string("scope-a"), string("key-a"))
						return err
					})
				},
				Defect: watchesWatchIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Watch.CancelGenerated(),
				Statement:   "Watch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnScope := c.Draw(prop.Of[string](), "scope")
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Watch(ctx, drawnScope, drawnKey)
						return err
					})
				},
				Defect: watchesWatchIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Watch.WatchSeesEachTrigger(),
				Statement:   "a subscription of Watch returns each update that Set causes for its key, in order, and no update after its stop",
				Claim:       spec.WatchSeesEachTrigger,
				Probe:       spec.WatchesSeeTheTriggers,
				Reason:      spec.MissedUpdate,
				Law:         spec.WatchReturnsEachTrigger,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *WatchesHandle, c *prop.Case) {
					probe.WatchesSeeTheTriggers(
						h,
						c,
						watchesActions,
						"Watch",
						"Set",
						Subscription.Next,
						Subscription.Stop,
						func(ctx context.Context, args []any) conform.Output {
							values, err := watchesSetCall(
								ctx,
								h.Subject(),
								[]any{args[0], args[1], args[2]},
							)
							return conform.OutputOf(err, values...)
						},
						func(args []any) any {
							return []any{args[1], args[0]}
						},
						func(args []any) any {
							return []any{args[0], args[1]}
						},
						func(args []any) Update {
							v, _ := args[2].(Update)
							return v
						},
					)
				},
				Defect: watchesSetMissesUpdatesAfterThree,
				Around: synctest.Test,
			},
			{
				ID:          WatchesChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Set(ctx, string("key-a"), string("scope-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesSetPanics,
			},
			{
				ID:          WatchesChecks.Set.SmokeGenerated(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnScope := c.Draw(prop.Of[string](), "scope")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnKey, drawnScope, drawnU)
					})
				},
				Defect: watchesSetPanics,
			},
			{
				ID:          WatchesChecks.Set.Cancel(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Set(ctx, string("key-a"), string("scope-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesSetIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Set.CancelGenerated(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnScope := c.Draw(prop.Of[string](), "scope")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnKey, drawnScope, drawnU)
					})
				},
				Defect: watchesSetIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Ticks.Smoke(),
				Statement:   "Ticks returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Ticks(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesTicksPanics,
			},
			{
				ID:          WatchesChecks.Ticks.SmokeGenerated(),
				Statement:   "Ticks returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Ticks(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesTicksPanics,
			},
			{
				ID:          WatchesChecks.Ticks.Cancel(),
				Statement:   "Ticks returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Ticks(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesTicksIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Ticks.CancelGenerated(),
				Statement:   "Ticks returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Ticks(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesTicksIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Ticks.WatchSeesEachTrigger(),
				Statement:   "a subscription of Ticks returns each update that Tick causes for its key, in order, and no update after its stop",
				Claim:       spec.WatchSeesEachTrigger,
				Probe:       spec.WatchesSeeTheTriggers,
				Reason:      spec.MissedUpdate,
				Law:         spec.WatchReturnsEachTrigger,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *WatchesHandle, c *prop.Case) {
					probe.WatchesSeeTheTriggers(
						h,
						c,
						watchesActions,
						"Ticks",
						"Tick",
						(*Ticker).Next,
						(*Ticker).Stop,
						func(ctx context.Context, args []any) conform.Output {
							values, err := watchesTickCall(
								ctx,
								h.Subject(),
								[]any{args[0], args[1], args[2]},
							)
							return conform.OutputOf(err, values...)
						},
						func(args []any) any {
							return []any{args[0]}
						},
						func(args []any) any {
							return []any{args[0]}
						},
						func(args []any) Update {
							v, _ := args[2].(Update)
							return v
						},
					)
				},
				Defect: watchesTickMissesUpdatesAfterThree,
				Around: synctest.Test,
			},
			{
				ID:          WatchesChecks.Tick.Smoke(),
				Statement:   "Tick returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Tick(ctx, string("key-a"), string("note-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesTickPanics,
			},
			{
				ID:          WatchesChecks.Tick.SmokeGenerated(),
				Statement:   "Tick returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnNote := c.Draw(prop.Of[string](), "note")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Tick(ctx, drawnKey, drawnNote, drawnU)
					})
				},
				Defect: watchesTickPanics,
			},
			{
				ID:          WatchesChecks.Tick.Cancel(),
				Statement:   "Tick returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Tick(ctx, string("key-a"), string("note-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesTickIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Tick.CancelGenerated(),
				Statement:   "Tick returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnNote := c.Draw(prop.Of[string](), "note")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Tick(ctx, drawnKey, drawnNote, drawnU)
					})
				},
				Defect: watchesTickIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Peeks.Smoke(),
				Statement:   "Peeks returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peeks(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesPeeksPanics,
			},
			{
				ID:          WatchesChecks.Peeks.SmokeGenerated(),
				Statement:   "Peeks returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peeks(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesPeeksPanics,
			},
			{
				ID:          WatchesChecks.Peeks.Cancel(),
				Statement:   "Peeks returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peeks(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesPeeksIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Peeks.CancelGenerated(),
				Statement:   "Peeks returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peeks(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesPeeksIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Peeked.Smoke(),
				Statement:   "Peeked returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Peeked(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesPeekedPanics,
			},
			{
				ID:          WatchesChecks.Peeked.SmokeGenerated(),
				Statement:   "Peeked returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Peeked(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesPeekedPanics,
			},
			{
				ID:          WatchesChecks.Peeked.Cancel(),
				Statement:   "Peeked returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Peeked(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesPeekedIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Peeked.CancelGenerated(),
				Statement:   "Peeked returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Peeked(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesPeekedIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Finds.Smoke(),
				Statement:   "Finds returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Finds(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesFindsPanics,
			},
			{
				ID:          WatchesChecks.Finds.SmokeGenerated(),
				Statement:   "Finds returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Finds(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesFindsPanics,
			},
			{
				ID:          WatchesChecks.Finds.Cancel(),
				Statement:   "Finds returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Finds(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesFindsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Finds.CancelGenerated(),
				Statement:   "Finds returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Finds(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesFindsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Found.Smoke(),
				Statement:   "Found returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Found(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesFoundPanics,
			},
			{
				ID:          WatchesChecks.Found.SmokeGenerated(),
				Statement:   "Found returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Found(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesFoundPanics,
			},
			{
				ID:          WatchesChecks.Found.Cancel(),
				Statement:   "Found returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Found(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesFoundIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Found.CancelGenerated(),
				Statement:   "Found returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Found(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesFoundIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Reads.Smoke(),
				Statement:   "Reads returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Reads(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesReadsPanics,
			},
			{
				ID:          WatchesChecks.Reads.SmokeGenerated(),
				Statement:   "Reads returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Reads(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesReadsPanics,
			},
			{
				ID:          WatchesChecks.Reads.Cancel(),
				Statement:   "Reads returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Reads(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesReadsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Reads.CancelGenerated(),
				Statement:   "Reads returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Reads(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesReadsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Written.Smoke(),
				Statement:   "Written returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Written(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesWrittenPanics,
			},
			{
				ID:          WatchesChecks.Written.SmokeGenerated(),
				Statement:   "Written returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Written(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesWrittenPanics,
			},
			{
				ID:          WatchesChecks.Written.Cancel(),
				Statement:   "Written returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Written(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesWrittenIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Written.CancelGenerated(),
				Statement:   "Written returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Written(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesWrittenIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Cancels.Smoke(),
				Statement:   "Cancels returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Cancels(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesCancelsPanics,
			},
			{
				ID:          WatchesChecks.Cancels.SmokeGenerated(),
				Statement:   "Cancels returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Cancels(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesCancelsPanics,
			},
			{
				ID:          WatchesChecks.Cancels.Cancel(),
				Statement:   "Cancels returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Cancels(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesCancelsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Cancels.CancelGenerated(),
				Statement:   "Cancels returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Cancels(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesCancelsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Cancelled.Smoke(),
				Statement:   "Cancelled returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Cancelled(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesCancelledPanics,
			},
			{
				ID:          WatchesChecks.Cancelled.SmokeGenerated(),
				Statement:   "Cancelled returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Cancelled(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesCancelledPanics,
			},
			{
				ID:          WatchesChecks.Cancelled.Cancel(),
				Statement:   "Cancelled returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Cancelled(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesCancelledIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Cancelled.CancelGenerated(),
				Statement:   "Cancelled returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Cancelled(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesCancelledIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Closes.Smoke(),
				Statement:   "Closes returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Closes(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesClosesPanics,
			},
			{
				ID:          WatchesChecks.Closes.SmokeGenerated(),
				Statement:   "Closes returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Closes(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesClosesPanics,
			},
			{
				ID:          WatchesChecks.Closes.Cancel(),
				Statement:   "Closes returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Closes(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesClosesIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Closes.CancelGenerated(),
				Statement:   "Closes returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Closes(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesClosesIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Closed.Smoke(),
				Statement:   "Closed returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Closed(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesClosedPanics,
			},
			{
				ID:          WatchesChecks.Closed.SmokeGenerated(),
				Statement:   "Closed returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Closed(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesClosedPanics,
			},
			{
				ID:          WatchesChecks.Closed.Cancel(),
				Statement:   "Closed returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Closed(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesClosedIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Closed.CancelGenerated(),
				Statement:   "Closed returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Closed(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesClosedIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Fields.Smoke(),
				Statement:   "Fields returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fields(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesFieldsPanics,
			},
			{
				ID:          WatchesChecks.Fields.SmokeGenerated(),
				Statement:   "Fields returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fields(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesFieldsPanics,
			},
			{
				ID:          WatchesChecks.Fields.Cancel(),
				Statement:   "Fields returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fields(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesFieldsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Fields.CancelGenerated(),
				Statement:   "Fields returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fields(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesFieldsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Filled.Smoke(),
				Statement:   "Filled returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Filled(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesFilledPanics,
			},
			{
				ID:          WatchesChecks.Filled.SmokeGenerated(),
				Statement:   "Filled returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Filled(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesFilledPanics,
			},
			{
				ID:          WatchesChecks.Filled.Cancel(),
				Statement:   "Filled returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Filled(ctx, string("key-a"), Update{Key: "Key-a"})
					})
				},
				Defect: watchesFilledIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Filled.CancelGenerated(),
				Statement:   "Filled returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnU := c.Draw(prop.Of[Update](), "u")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Filled(ctx, drawnKey, drawnU)
					})
				},
				Defect: watchesFilledIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Follows.Smoke(),
				Statement:   "Follows returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Follows(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesFollowsPanics,
			},
			{
				ID:          WatchesChecks.Follows.SmokeGenerated(),
				Statement:   "Follows returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Follows(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesFollowsPanics,
			},
			{
				ID:          WatchesChecks.Follows.Cancel(),
				Statement:   "Follows returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Follows(ctx, string("key-a"))
						return err
					})
				},
				Defect: watchesFollowsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Follows.CancelGenerated(),
				Statement:   "Follows returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Follows(ctx, drawnKey)
						return err
					})
				},
				Defect: watchesFollowsIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Touch.Smoke(),
				Statement:   "Touch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WatchesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Touch(ctx, string("key-a"))
					})
				},
				Defect: watchesTouchPanics,
			},
			{
				ID:          WatchesChecks.Touch.SmokeGenerated(),
				Statement:   "Touch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Touch(ctx, drawnKey)
					})
				},
				Defect: watchesTouchPanics,
			},
			{
				ID:          WatchesChecks.Touch.Cancel(),
				Statement:   "Touch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WatchesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Touch(ctx, string("key-a"))
					})
				},
				Defect: watchesTouchIgnoresCancel,
			},
			{
				ID:          WatchesChecks.Touch.CancelGenerated(),
				Statement:   "Touch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WatchesHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Touch(ctx, drawnKey)
					})
				},
				Defect: watchesTouchIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Ticks",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Ticks",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Peeks",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Peeks",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Finds",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Finds",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Reads",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Reads",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Cancels",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Cancels",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Closes",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Closes",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Fields",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Fields",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Watches.Follows",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Follows",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.WatchSeesEachTrigger,
				Callable: "Watches.Peeks",
				Gate:     spec.SubscriptionMethods,
				Why:      "the methods that next and stop name on the subscription of Peeks, or the value that Peeked passes, do not fit a read of updates",
				Remedy:   "name with next a method that takes a wait and returns the update and a flag, and with stop a method without parameters and results",
			},
			{
				Claim:    spec.WatchSeesEachTrigger,
				Callable: "Watches.Finds",
				Gate:     spec.SubscriptionMethods,
				Why:      "the methods that next and stop name on the subscription of Finds, or the value that Found passes, do not fit a read of updates",
				Remedy:   "name with next a method that takes a wait and returns the update and a flag, and with stop a method without parameters and results",
			},
			{
				Claim:    spec.WatchSeesEachTrigger,
				Callable: "Watches.Reads",
				Gate:     spec.SubscriptionMethods,
				Why:      "the methods that next and stop name on the subscription of Reads, or the value that Written passes, do not fit a read of updates",
				Remedy:   "name with next a method that takes a wait and returns the update and a flag, and with stop a method without parameters and results",
			},
			{
				Claim:    spec.WatchSeesEachTrigger,
				Callable: "Watches.Cancels",
				Gate:     spec.SubscriptionMethods,
				Why:      "the methods that next and stop name on the subscription of Cancels, or the value that Cancelled passes, do not fit a read of updates",
				Remedy:   "name with next a method that takes a wait and returns the update and a flag, and with stop a method without parameters and results",
			},
			{
				Claim:    spec.WatchSeesEachTrigger,
				Callable: "Watches.Closes",
				Gate:     spec.SubscriptionMethods,
				Why:      "the methods that next and stop name on the subscription of Closes, or the value that Closed passes, do not fit a read of updates",
				Remedy:   "name with next a method that takes a wait and returns the update and a flag, and with stop a method without parameters and results",
			},
			{
				Claim:    spec.WatchSeesEachTrigger,
				Callable: "Watches.Fields",
				Gate:     spec.SubscriptionMethods,
				Why:      "the methods that next and stop name on the subscription of Fields, or the value that Filled passes, do not fit a read of updates",
				Remedy:   "name with next a method that takes a wait and returns the update and a flag, and with stop a method without parameters and results",
			},
			{
				Claim:    spec.WatchSeesEachTrigger,
				Callable: "Watches.Follows",
				Gate:     spec.SubscriptionMethods,
				Why:      "the methods that next and stop name on the subscription of Follows, or the value that Touch passes, do not fit a read of updates",
				Remedy:   "name with next a method that takes a wait and returns the update and a flag, and with stop a method without parameters and results",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:multi-arg-writer",
				Callable:       "Watches.Set",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:multi-arg-writer",
				Callable:       "Watches.Tick",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Watches.Peeked",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Watches.Found",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Watches.Written",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Watches.Cancelled",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Watches.Closed",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Watches.Filled",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Watches.Touch",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Watches]{
			{
				Calls:    watchesWatchCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Watch.Callable,
			},
			{
				Calls:    watchesSetCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Set.Callable,
			},
			{
				Calls:    watchesTicksCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Ticks.Callable,
			},
			{
				Calls:    watchesTickCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Tick.Callable,
			},
			{
				Calls:    watchesPeeksCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Peeks.Callable,
			},
			{
				Calls:    watchesPeekedCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Peeked.Callable,
			},
			{
				Calls:    watchesFindsCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Finds.Callable,
			},
			{
				Calls:    watchesFoundCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Found.Callable,
			},
			{
				Calls:    watchesReadsCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Reads.Callable,
			},
			{
				Calls:    watchesWrittenCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Written.Callable,
			},
			{
				Calls:    watchesCancelsCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Cancels.Callable,
			},
			{
				Calls:    watchesCancelledCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Cancelled.Callable,
			},
			{
				Calls:    watchesClosesCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Closes.Callable,
			},
			{
				Calls:    watchesClosedCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Closed.Callable,
			},
			{
				Calls:    watchesFieldsCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Fields.Callable,
			},
			{
				Calls:    watchesFilledCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Filled.Callable,
			},
			{
				Calls:    watchesFollowsCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Follows.Callable,
			},
			{
				Calls:    watchesTouchCalls,
				Over:     newWatchesOver,
				Callable: WatchesChecks.Touch.Callable,
			},
		},
	}
}

// RunWatchesSuite applies the checks of Watches, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunWatchesSuite[T Watches](t *testing.T, newSubject func() T, opts ...WatchesOption) {
	t.Helper()
	RunWatchesSuiteFromHandle(t, func(*WatchesHandle) T { return newSubject() }, opts...)
}

// RunWatchesSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunWatchesSuiteFromHandle[T Watches](
	t *testing.T,
	newSubject func(h *WatchesHandle) T,
	opts ...WatchesOption,
) {
	t.Helper()
	conform.Run(t, WatchesSuite(), watchesSubjects(newSubject), opts...)
}

// FuzzWatchesSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzWatchesSuite[T Watches](
	f *testing.F,
	callable conform.Callable[Watches],
	newSubject func(h *WatchesHandle) T,
	opts ...WatchesOption,
) {
	f.Helper()
	conform.Fuzz(f, WatchesSuite(), callable, watchesSubjects(newSubject), opts...)
}

// RunWatchesBudgets applies the allocation ceilings of Watches and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunWatchesBudgets[T Watches](
	t *testing.T,
	newSubject func(h *WatchesHandle) T,
	opts ...WatchesOption,
) {
	t.Helper()
	conform.RunBudgets(t, WatchesSuite(), watchesSubjects(newSubject), opts...)
}

// BenchWatchesSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchWatchesSuite[T Watches](
	b conform.B,
	callable conform.Callable[Watches],
	newSubject func(h *WatchesHandle) T,
	opts ...WatchesOption,
) {
	b.Helper()
	conform.Bench(b, WatchesSuite(), callable, watchesSubjects(newSubject), opts...)
}

// watchesSubjects returns the factory of the subjects of Watches that
// calls newSubject with the handle of each check.
func watchesSubjects[T Watches](
	newSubject func(h *WatchesHandle) T,
) func(h *WatchesHandle) Watches {
	return func(h *WatchesHandle) Watches { return newSubject(h) }
}

// watchesWatchCalls generates the calls of Watches.Watch with the inputs
// that each case draws.
var watchesWatchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	scope := c.Draw(prop.Of[string](), "scope")
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Watch(ctx, scope, key)
		return err
	}
})

// watchesSetCalls generates the calls of Watches.Set with the inputs that
// each case draws.
var watchesSetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	scope := c.Draw(prop.Of[string](), "scope")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Set(ctx, key, scope, u)
	}
})

// watchesTicksCalls generates the calls of Watches.Ticks with the inputs
// that each case draws.
var watchesTicksCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Ticks(ctx, key)
		return err
	}
})

// watchesTickCalls generates the calls of Watches.Tick with the inputs
// that each case draws.
var watchesTickCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	note := c.Draw(prop.Of[string](), "note")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Tick(ctx, key, note, u)
	}
})

// watchesPeeksCalls generates the calls of Watches.Peeks with the inputs
// that each case draws.
var watchesPeeksCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Peeks(ctx, key)
		return err
	}
})

// watchesPeekedCalls generates the calls of Watches.Peeked with the inputs
// that each case draws.
var watchesPeekedCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Peeked(ctx, key, u)
	}
})

// watchesFindsCalls generates the calls of Watches.Finds with the inputs
// that each case draws.
var watchesFindsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Finds(ctx, key)
		return err
	}
})

// watchesFoundCalls generates the calls of Watches.Found with the inputs
// that each case draws.
var watchesFoundCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Found(ctx, key, u)
	}
})

// watchesReadsCalls generates the calls of Watches.Reads with the inputs
// that each case draws.
var watchesReadsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Reads(ctx, key)
		return err
	}
})

// watchesWrittenCalls generates the calls of Watches.Written with the
// inputs that each case draws.
var watchesWrittenCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Written(ctx, key, u)
	}
})

// watchesCancelsCalls generates the calls of Watches.Cancels with the
// inputs that each case draws.
var watchesCancelsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Cancels(ctx, key)
		return err
	}
})

// watchesCancelledCalls generates the calls of Watches.Cancelled with the
// inputs that each case draws.
var watchesCancelledCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Cancelled(ctx, key, u)
	}
})

// watchesClosesCalls generates the calls of Watches.Closes with the inputs
// that each case draws.
var watchesClosesCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Closes(ctx, key)
		return err
	}
})

// watchesClosedCalls generates the calls of Watches.Closed with the inputs
// that each case draws.
var watchesClosedCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Closed(ctx, key, u)
	}
})

// watchesFieldsCalls generates the calls of Watches.Fields with the inputs
// that each case draws.
var watchesFieldsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Fields(ctx, key)
		return err
	}
})

// watchesFilledCalls generates the calls of Watches.Filled with the inputs
// that each case draws.
var watchesFilledCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	u := c.Draw(prop.Of[Update](), "u")
	return func(ctx context.Context, subject Watches) error {
		return subject.Filled(ctx, key, u)
	}
})

// watchesFollowsCalls generates the calls of Watches.Follows with the
// inputs that each case draws.
var watchesFollowsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		_, err := subject.Follows(ctx, key)
		return err
	}
})

// watchesTouchCalls generates the calls of Watches.Touch with the inputs
// that each case draws.
var watchesTouchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Watches) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Watches) error {
		return subject.Touch(ctx, key)
	}
})

// watchesOverEmbedded is Watches under the name of the field that
// watchesOver embeds.
type watchesOverEmbedded = Watches

// watchesOver is a subject of Watches whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type watchesOver struct {
	watchesOverEmbedded
	spend func()
}

// Watch calls spend and returns its zero results.
func (o watchesOver) Watch(context.Context, string, string) (r0 Subscription, r1 error) {
	o.spend()
	return r0, r1
}

// Set calls spend and returns its zero results.
func (o watchesOver) Set(context.Context, string, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Ticks calls spend and returns its zero results.
func (o watchesOver) Ticks(context.Context, string) (r0 *Ticker, r1 error) {
	o.spend()
	return r0, r1
}

// Tick calls spend and returns its zero results.
func (o watchesOver) Tick(context.Context, string, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Peeks calls spend and returns its zero results.
func (o watchesOver) Peeks(context.Context, string) (r0 Subscription, r1 error) {
	o.spend()
	return r0, r1
}

// Peeked calls spend and returns its zero results.
func (o watchesOver) Peeked(context.Context, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Finds calls spend and returns its zero results.
func (o watchesOver) Finds(context.Context, string) (r0 Subscription, r1 error) {
	o.spend()
	return r0, r1
}

// Found calls spend and returns its zero results.
func (o watchesOver) Found(context.Context, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Reads calls spend and returns its zero results.
func (o watchesOver) Reads(context.Context, string) (r0 Subscription, r1 error) {
	o.spend()
	return r0, r1
}

// Written calls spend and returns its zero results.
func (o watchesOver) Written(context.Context, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Cancels calls spend and returns its zero results.
func (o watchesOver) Cancels(context.Context, string) (r0 Subscription, r1 error) {
	o.spend()
	return r0, r1
}

// Cancelled calls spend and returns its zero results.
func (o watchesOver) Cancelled(context.Context, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Closes calls spend and returns its zero results.
func (o watchesOver) Closes(context.Context, string) (r0 Subscription, r1 error) {
	o.spend()
	return r0, r1
}

// Closed calls spend and returns its zero results.
func (o watchesOver) Closed(context.Context, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Fields calls spend and returns its zero results.
func (o watchesOver) Fields(context.Context, string) (r0 *Ticker, r1 error) {
	o.spend()
	return r0, r1
}

// Filled calls spend and returns its zero results.
func (o watchesOver) Filled(context.Context, string, Update) (r0 error) {
	o.spend()
	return r0
}

// Follows calls spend and returns its zero results.
func (o watchesOver) Follows(context.Context, string) (r0 Subscription, r1 error) {
	o.spend()
	return r0, r1
}

// Touch calls spend and returns its zero results.
func (o watchesOver) Touch(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// newWatchesOver returns the planted defect of a ceiling of Watches, whose
// callables call spend.
func newWatchesOver(_ *WatchesHandle, spend func()) Watches {
	return watchesOver{spend: spend}
}

// watchesWatchPanicsDefect is the planted defect of a Watches whose Watch
// panics.
type watchesWatchPanicsDefect struct {
	watchesOverEmbedded
}

// Watch is the override of the planted defect.
func (*watchesWatchPanicsDefect) Watch(
	p0 context.Context,
	p1, p2 string,
) (r0 Subscription, r1 error) {
	panic("the planted defect of a Watches whose Watch panics")
}

// watchesWatchPanics returns the planted defect of a Watches whose Watch
// panics over a fresh subject.
func watchesWatchPanics(h *WatchesHandle) Watches {
	return &watchesWatchPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesWatchIgnoresCancelDefect is the planted defect of a Watches whose
// Watch ignores the context.
type watchesWatchIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Watch is the override of the planted defect.
func (d *watchesWatchIgnoresCancelDefect) Watch(
	p0 context.Context,
	p1, p2 string,
) (r0 Subscription, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Watch(p0, p1, p2)
}

// watchesWatchIgnoresCancel returns the planted defect of a Watches whose
// Watch ignores the context over a fresh subject.
func watchesWatchIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesWatchIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesSetMissesUpdatesAfterThreeDefect is the planted defect of a
// Watches whose Set does nothing from its fourth call on.
type watchesSetMissesUpdatesAfterThreeDefect struct {
	watchesOverEmbedded
	calls int
}

// Set is the override of the planted defect.
func (d *watchesSetMissesUpdatesAfterThreeDefect) Set(
	p0 context.Context,
	p1, p2 string,
	p3 Update,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		return r0
	}
	return d.watchesOverEmbedded.Set(p0, p1, p2, p3)
}

// watchesSetMissesUpdatesAfterThree returns the planted defect of a
// Watches whose Set does nothing from its fourth call on over a fresh
// subject.
func watchesSetMissesUpdatesAfterThree(h *WatchesHandle) Watches {
	return &watchesSetMissesUpdatesAfterThreeDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesSetPanicsDefect is the planted defect of a Watches whose Set
// panics.
type watchesSetPanicsDefect struct {
	watchesOverEmbedded
}

// Set is the override of the planted defect.
func (*watchesSetPanicsDefect) Set(p0 context.Context, p1, p2 string, p3 Update) (r0 error) {
	panic("the planted defect of a Watches whose Set panics")
}

// watchesSetPanics returns the planted defect of a Watches whose Set
// panics over a fresh subject.
func watchesSetPanics(h *WatchesHandle) Watches {
	return &watchesSetPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesSetIgnoresCancelDefect is the planted defect of a Watches whose
// Set ignores the context.
type watchesSetIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Set is the override of the planted defect.
func (d *watchesSetIgnoresCancelDefect) Set(
	p0 context.Context,
	p1, p2 string,
	p3 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Set(p0, p1, p2, p3)
}

// watchesSetIgnoresCancel returns the planted defect of a Watches whose
// Set ignores the context over a fresh subject.
func watchesSetIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesSetIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesTicksPanicsDefect is the planted defect of a Watches whose Ticks
// panics.
type watchesTicksPanicsDefect struct {
	watchesOverEmbedded
}

// Ticks is the override of the planted defect.
func (*watchesTicksPanicsDefect) Ticks(p0 context.Context, p1 string) (r0 *Ticker, r1 error) {
	panic("the planted defect of a Watches whose Ticks panics")
}

// watchesTicksPanics returns the planted defect of a Watches whose Ticks
// panics over a fresh subject.
func watchesTicksPanics(h *WatchesHandle) Watches {
	return &watchesTicksPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesTicksIgnoresCancelDefect is the planted defect of a Watches whose
// Ticks ignores the context.
type watchesTicksIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Ticks is the override of the planted defect.
func (d *watchesTicksIgnoresCancelDefect) Ticks(
	p0 context.Context,
	p1 string,
) (r0 *Ticker, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Ticks(p0, p1)
}

// watchesTicksIgnoresCancel returns the planted defect of a Watches whose
// Ticks ignores the context over a fresh subject.
func watchesTicksIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesTicksIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesTickMissesUpdatesAfterThreeDefect is the planted defect of a
// Watches whose Tick does nothing from its fourth call on.
type watchesTickMissesUpdatesAfterThreeDefect struct {
	watchesOverEmbedded
	calls int
}

// Tick is the override of the planted defect.
func (d *watchesTickMissesUpdatesAfterThreeDefect) Tick(
	p0 context.Context,
	p1, p2 string,
	p3 Update,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		return r0
	}
	return d.watchesOverEmbedded.Tick(p0, p1, p2, p3)
}

// watchesTickMissesUpdatesAfterThree returns the planted defect of a
// Watches whose Tick does nothing from its fourth call on over a fresh
// subject.
func watchesTickMissesUpdatesAfterThree(h *WatchesHandle) Watches {
	return &watchesTickMissesUpdatesAfterThreeDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesTickPanicsDefect is the planted defect of a Watches whose Tick
// panics.
type watchesTickPanicsDefect struct {
	watchesOverEmbedded
}

// Tick is the override of the planted defect.
func (*watchesTickPanicsDefect) Tick(p0 context.Context, p1, p2 string, p3 Update) (r0 error) {
	panic("the planted defect of a Watches whose Tick panics")
}

// watchesTickPanics returns the planted defect of a Watches whose Tick
// panics over a fresh subject.
func watchesTickPanics(h *WatchesHandle) Watches {
	return &watchesTickPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesTickIgnoresCancelDefect is the planted defect of a Watches whose
// Tick ignores the context.
type watchesTickIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Tick is the override of the planted defect.
func (d *watchesTickIgnoresCancelDefect) Tick(
	p0 context.Context,
	p1, p2 string,
	p3 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Tick(p0, p1, p2, p3)
}

// watchesTickIgnoresCancel returns the planted defect of a Watches whose
// Tick ignores the context over a fresh subject.
func watchesTickIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesTickIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesPeeksPanicsDefect is the planted defect of a Watches whose Peeks
// panics.
type watchesPeeksPanicsDefect struct {
	watchesOverEmbedded
}

// Peeks is the override of the planted defect.
func (*watchesPeeksPanicsDefect) Peeks(p0 context.Context, p1 string) (r0 Subscription, r1 error) {
	panic("the planted defect of a Watches whose Peeks panics")
}

// watchesPeeksPanics returns the planted defect of a Watches whose Peeks
// panics over a fresh subject.
func watchesPeeksPanics(h *WatchesHandle) Watches {
	return &watchesPeeksPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesPeeksIgnoresCancelDefect is the planted defect of a Watches whose
// Peeks ignores the context.
type watchesPeeksIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Peeks is the override of the planted defect.
func (d *watchesPeeksIgnoresCancelDefect) Peeks(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Peeks(p0, p1)
}

// watchesPeeksIgnoresCancel returns the planted defect of a Watches whose
// Peeks ignores the context over a fresh subject.
func watchesPeeksIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesPeeksIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesPeekedPanicsDefect is the planted defect of a Watches whose
// Peeked panics.
type watchesPeekedPanicsDefect struct {
	watchesOverEmbedded
}

// Peeked is the override of the planted defect.
func (*watchesPeekedPanicsDefect) Peeked(p0 context.Context, p1 string, p2 Update) (r0 error) {
	panic("the planted defect of a Watches whose Peeked panics")
}

// watchesPeekedPanics returns the planted defect of a Watches whose Peeked
// panics over a fresh subject.
func watchesPeekedPanics(h *WatchesHandle) Watches {
	return &watchesPeekedPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesPeekedIgnoresCancelDefect is the planted defect of a Watches
// whose Peeked ignores the context.
type watchesPeekedIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Peeked is the override of the planted defect.
func (d *watchesPeekedIgnoresCancelDefect) Peeked(
	p0 context.Context,
	p1 string,
	p2 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Peeked(p0, p1, p2)
}

// watchesPeekedIgnoresCancel returns the planted defect of a Watches whose
// Peeked ignores the context over a fresh subject.
func watchesPeekedIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesPeekedIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFindsPanicsDefect is the planted defect of a Watches whose Finds
// panics.
type watchesFindsPanicsDefect struct {
	watchesOverEmbedded
}

// Finds is the override of the planted defect.
func (*watchesFindsPanicsDefect) Finds(p0 context.Context, p1 string) (r0 Subscription, r1 error) {
	panic("the planted defect of a Watches whose Finds panics")
}

// watchesFindsPanics returns the planted defect of a Watches whose Finds
// panics over a fresh subject.
func watchesFindsPanics(h *WatchesHandle) Watches {
	return &watchesFindsPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFindsIgnoresCancelDefect is the planted defect of a Watches whose
// Finds ignores the context.
type watchesFindsIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Finds is the override of the planted defect.
func (d *watchesFindsIgnoresCancelDefect) Finds(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Finds(p0, p1)
}

// watchesFindsIgnoresCancel returns the planted defect of a Watches whose
// Finds ignores the context over a fresh subject.
func watchesFindsIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesFindsIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFoundPanicsDefect is the planted defect of a Watches whose Found
// panics.
type watchesFoundPanicsDefect struct {
	watchesOverEmbedded
}

// Found is the override of the planted defect.
func (*watchesFoundPanicsDefect) Found(p0 context.Context, p1 string, p2 Update) (r0 error) {
	panic("the planted defect of a Watches whose Found panics")
}

// watchesFoundPanics returns the planted defect of a Watches whose Found
// panics over a fresh subject.
func watchesFoundPanics(h *WatchesHandle) Watches {
	return &watchesFoundPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFoundIgnoresCancelDefect is the planted defect of a Watches whose
// Found ignores the context.
type watchesFoundIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Found is the override of the planted defect.
func (d *watchesFoundIgnoresCancelDefect) Found(
	p0 context.Context,
	p1 string,
	p2 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Found(p0, p1, p2)
}

// watchesFoundIgnoresCancel returns the planted defect of a Watches whose
// Found ignores the context over a fresh subject.
func watchesFoundIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesFoundIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesReadsPanicsDefect is the planted defect of a Watches whose Reads
// panics.
type watchesReadsPanicsDefect struct {
	watchesOverEmbedded
}

// Reads is the override of the planted defect.
func (*watchesReadsPanicsDefect) Reads(p0 context.Context, p1 string) (r0 Subscription, r1 error) {
	panic("the planted defect of a Watches whose Reads panics")
}

// watchesReadsPanics returns the planted defect of a Watches whose Reads
// panics over a fresh subject.
func watchesReadsPanics(h *WatchesHandle) Watches {
	return &watchesReadsPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesReadsIgnoresCancelDefect is the planted defect of a Watches whose
// Reads ignores the context.
type watchesReadsIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Reads is the override of the planted defect.
func (d *watchesReadsIgnoresCancelDefect) Reads(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Reads(p0, p1)
}

// watchesReadsIgnoresCancel returns the planted defect of a Watches whose
// Reads ignores the context over a fresh subject.
func watchesReadsIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesReadsIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesWrittenPanicsDefect is the planted defect of a Watches whose
// Written panics.
type watchesWrittenPanicsDefect struct {
	watchesOverEmbedded
}

// Written is the override of the planted defect.
func (*watchesWrittenPanicsDefect) Written(p0 context.Context, p1 string, p2 Update) (r0 error) {
	panic("the planted defect of a Watches whose Written panics")
}

// watchesWrittenPanics returns the planted defect of a Watches whose
// Written panics over a fresh subject.
func watchesWrittenPanics(h *WatchesHandle) Watches {
	return &watchesWrittenPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesWrittenIgnoresCancelDefect is the planted defect of a Watches
// whose Written ignores the context.
type watchesWrittenIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Written is the override of the planted defect.
func (d *watchesWrittenIgnoresCancelDefect) Written(
	p0 context.Context,
	p1 string,
	p2 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Written(p0, p1, p2)
}

// watchesWrittenIgnoresCancel returns the planted defect of a Watches
// whose Written ignores the context over a fresh subject.
func watchesWrittenIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesWrittenIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesCancelsPanicsDefect is the planted defect of a Watches whose
// Cancels panics.
type watchesCancelsPanicsDefect struct {
	watchesOverEmbedded
}

// Cancels is the override of the planted defect.
func (*watchesCancelsPanicsDefect) Cancels(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	panic("the planted defect of a Watches whose Cancels panics")
}

// watchesCancelsPanics returns the planted defect of a Watches whose
// Cancels panics over a fresh subject.
func watchesCancelsPanics(h *WatchesHandle) Watches {
	return &watchesCancelsPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesCancelsIgnoresCancelDefect is the planted defect of a Watches
// whose Cancels ignores the context.
type watchesCancelsIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Cancels is the override of the planted defect.
func (d *watchesCancelsIgnoresCancelDefect) Cancels(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Cancels(p0, p1)
}

// watchesCancelsIgnoresCancel returns the planted defect of a Watches
// whose Cancels ignores the context over a fresh subject.
func watchesCancelsIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesCancelsIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesCancelledPanicsDefect is the planted defect of a Watches whose
// Cancelled panics.
type watchesCancelledPanicsDefect struct {
	watchesOverEmbedded
}

// Cancelled is the override of the planted defect.
func (*watchesCancelledPanicsDefect) Cancelled(
	p0 context.Context,
	p1 string,
	p2 Update,
) (r0 error) {
	panic("the planted defect of a Watches whose Cancelled panics")
}

// watchesCancelledPanics returns the planted defect of a Watches whose
// Cancelled panics over a fresh subject.
func watchesCancelledPanics(h *WatchesHandle) Watches {
	return &watchesCancelledPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesCancelledIgnoresCancelDefect is the planted defect of a Watches
// whose Cancelled ignores the context.
type watchesCancelledIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Cancelled is the override of the planted defect.
func (d *watchesCancelledIgnoresCancelDefect) Cancelled(
	p0 context.Context,
	p1 string,
	p2 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Cancelled(p0, p1, p2)
}

// watchesCancelledIgnoresCancel returns the planted defect of a Watches
// whose Cancelled ignores the context over a fresh subject.
func watchesCancelledIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesCancelledIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesClosesPanicsDefect is the planted defect of a Watches whose
// Closes panics.
type watchesClosesPanicsDefect struct {
	watchesOverEmbedded
}

// Closes is the override of the planted defect.
func (*watchesClosesPanicsDefect) Closes(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	panic("the planted defect of a Watches whose Closes panics")
}

// watchesClosesPanics returns the planted defect of a Watches whose Closes
// panics over a fresh subject.
func watchesClosesPanics(h *WatchesHandle) Watches {
	return &watchesClosesPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesClosesIgnoresCancelDefect is the planted defect of a Watches
// whose Closes ignores the context.
type watchesClosesIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Closes is the override of the planted defect.
func (d *watchesClosesIgnoresCancelDefect) Closes(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Closes(p0, p1)
}

// watchesClosesIgnoresCancel returns the planted defect of a Watches whose
// Closes ignores the context over a fresh subject.
func watchesClosesIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesClosesIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesClosedPanicsDefect is the planted defect of a Watches whose
// Closed panics.
type watchesClosedPanicsDefect struct {
	watchesOverEmbedded
}

// Closed is the override of the planted defect.
func (*watchesClosedPanicsDefect) Closed(p0 context.Context, p1 string, p2 Update) (r0 error) {
	panic("the planted defect of a Watches whose Closed panics")
}

// watchesClosedPanics returns the planted defect of a Watches whose Closed
// panics over a fresh subject.
func watchesClosedPanics(h *WatchesHandle) Watches {
	return &watchesClosedPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesClosedIgnoresCancelDefect is the planted defect of a Watches
// whose Closed ignores the context.
type watchesClosedIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Closed is the override of the planted defect.
func (d *watchesClosedIgnoresCancelDefect) Closed(
	p0 context.Context,
	p1 string,
	p2 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Closed(p0, p1, p2)
}

// watchesClosedIgnoresCancel returns the planted defect of a Watches whose
// Closed ignores the context over a fresh subject.
func watchesClosedIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesClosedIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFieldsPanicsDefect is the planted defect of a Watches whose
// Fields panics.
type watchesFieldsPanicsDefect struct {
	watchesOverEmbedded
}

// Fields is the override of the planted defect.
func (*watchesFieldsPanicsDefect) Fields(p0 context.Context, p1 string) (r0 *Ticker, r1 error) {
	panic("the planted defect of a Watches whose Fields panics")
}

// watchesFieldsPanics returns the planted defect of a Watches whose Fields
// panics over a fresh subject.
func watchesFieldsPanics(h *WatchesHandle) Watches {
	return &watchesFieldsPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFieldsIgnoresCancelDefect is the planted defect of a Watches
// whose Fields ignores the context.
type watchesFieldsIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Fields is the override of the planted defect.
func (d *watchesFieldsIgnoresCancelDefect) Fields(
	p0 context.Context,
	p1 string,
) (r0 *Ticker, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Fields(p0, p1)
}

// watchesFieldsIgnoresCancel returns the planted defect of a Watches whose
// Fields ignores the context over a fresh subject.
func watchesFieldsIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesFieldsIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFilledPanicsDefect is the planted defect of a Watches whose
// Filled panics.
type watchesFilledPanicsDefect struct {
	watchesOverEmbedded
}

// Filled is the override of the planted defect.
func (*watchesFilledPanicsDefect) Filled(p0 context.Context, p1 string, p2 Update) (r0 error) {
	panic("the planted defect of a Watches whose Filled panics")
}

// watchesFilledPanics returns the planted defect of a Watches whose Filled
// panics over a fresh subject.
func watchesFilledPanics(h *WatchesHandle) Watches {
	return &watchesFilledPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFilledIgnoresCancelDefect is the planted defect of a Watches
// whose Filled ignores the context.
type watchesFilledIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Filled is the override of the planted defect.
func (d *watchesFilledIgnoresCancelDefect) Filled(
	p0 context.Context,
	p1 string,
	p2 Update,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Filled(p0, p1, p2)
}

// watchesFilledIgnoresCancel returns the planted defect of a Watches whose
// Filled ignores the context over a fresh subject.
func watchesFilledIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesFilledIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFollowsPanicsDefect is the planted defect of a Watches whose
// Follows panics.
type watchesFollowsPanicsDefect struct {
	watchesOverEmbedded
}

// Follows is the override of the planted defect.
func (*watchesFollowsPanicsDefect) Follows(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	panic("the planted defect of a Watches whose Follows panics")
}

// watchesFollowsPanics returns the planted defect of a Watches whose
// Follows panics over a fresh subject.
func watchesFollowsPanics(h *WatchesHandle) Watches {
	return &watchesFollowsPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesFollowsIgnoresCancelDefect is the planted defect of a Watches
// whose Follows ignores the context.
type watchesFollowsIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Follows is the override of the planted defect.
func (d *watchesFollowsIgnoresCancelDefect) Follows(
	p0 context.Context,
	p1 string,
) (r0 Subscription, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.watchesOverEmbedded.Follows(p0, p1)
}

// watchesFollowsIgnoresCancel returns the planted defect of a Watches
// whose Follows ignores the context over a fresh subject.
func watchesFollowsIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesFollowsIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesTouchPanicsDefect is the planted defect of a Watches whose Touch
// panics.
type watchesTouchPanicsDefect struct {
	watchesOverEmbedded
}

// Touch is the override of the planted defect.
func (*watchesTouchPanicsDefect) Touch(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Watches whose Touch panics")
}

// watchesTouchPanics returns the planted defect of a Watches whose Touch
// panics over a fresh subject.
func watchesTouchPanics(h *WatchesHandle) Watches {
	return &watchesTouchPanicsDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesTouchIgnoresCancelDefect is the planted defect of a Watches whose
// Touch ignores the context.
type watchesTouchIgnoresCancelDefect struct {
	watchesOverEmbedded
}

// Touch is the override of the planted defect.
func (d *watchesTouchIgnoresCancelDefect) Touch(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.watchesOverEmbedded.Touch(p0, p1)
}

// watchesTouchIgnoresCancel returns the planted defect of a Watches whose
// Touch ignores the context over a fresh subject.
func watchesTouchIgnoresCancel(h *WatchesHandle) Watches {
	return &watchesTouchIgnoresCancelDefect{watchesOverEmbedded: h.Fresh()}
}

// watchesActions returns the actions of the generated sequences of Watches
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func watchesActions(c *prop.Case, subject func() Watches) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Watch",
			Input: func(c *prop.Case, _ struct{}) any {
				scope := c.Draw(prop.Of[string](), "watch.scope")
				key := c.Draw(prop.Of[string](), "watch.key")
				return []any{scope, key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Watch", args)
				values, err := watchesWatchCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Set",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "set.key")
				scope := c.Draw(prop.Of[string](), "set.scope")
				u := c.Draw(prop.Of[Update](), "set.u")
				return []any{key, scope, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Set", args)
				values, err := watchesSetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Ticks",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "ticks.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Ticks", args, args[0])
				values, err := watchesTicksCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Tick",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "tick.key")
				note := c.Draw(prop.Of[string](), "tick.note")
				u := c.Draw(prop.Of[Update](), "tick.u")
				return []any{key, note, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Tick", args, args[0])
				values, err := watchesTickCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Peeks",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "peeks.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Peeks", args, args[0])
				values, err := watchesPeeksCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Peeked",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "peeked.key")
				u := c.Draw(prop.Of[Update](), "peeked.u")
				return []any{key, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Peeked", args, args[0])
				values, err := watchesPeekedCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Finds",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "finds.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Finds", args, args[0])
				values, err := watchesFindsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Found",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "found.key")
				u := c.Draw(prop.Of[Update](), "found.u")
				return []any{key, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Found", args, args[0])
				values, err := watchesFoundCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Reads",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "reads.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Reads", args, args[0])
				values, err := watchesReadsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Written",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "written.key")
				u := c.Draw(prop.Of[Update](), "written.u")
				return []any{key, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Written", args, args[0])
				values, err := watchesWrittenCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Cancels",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "cancels.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Cancels", args, args[0])
				values, err := watchesCancelsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Cancelled",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "cancelled.key")
				u := c.Draw(prop.Of[Update](), "cancelled.u")
				return []any{key, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Cancelled", args, args[0])
				values, err := watchesCancelledCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Closes",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "closes.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Closes", args, args[0])
				values, err := watchesClosesCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Closed",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "closed.key")
				u := c.Draw(prop.Of[Update](), "closed.u")
				return []any{key, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Closed", args, args[0])
				values, err := watchesClosedCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Fields",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "fields.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Fields", args, args[0])
				values, err := watchesFieldsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Filled",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "filled.key")
				u := c.Draw(prop.Of[Update](), "filled.u")
				return []any{key, u}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Filled", args, args[0])
				values, err := watchesFilledCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Follows",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "follows.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Follows", args, args[0])
				values, err := watchesFollowsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Touch",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "touch.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Touch", args, args[0])
				values, err := watchesTouchCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// watchesWatchCall calls Watch on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesWatchCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	r0, err := subject.Watch(ctx, a0, a1)
	return []any{r0}, err
}

// watchesSetCall calls Set on subject with the arguments args, and returns
// the values of the call beside its error.
func watchesSetCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	a2, _ := args[2].(Update)
	return nil, subject.Set(ctx, a0, a1, a2)
}

// watchesTicksCall calls Ticks on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesTicksCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Ticks(ctx, a0)
	return []any{r0}, err
}

// watchesTickCall calls Tick on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesTickCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	a2, _ := args[2].(Update)
	return nil, subject.Tick(ctx, a0, a1, a2)
}

// watchesPeeksCall calls Peeks on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesPeeksCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Peeks(ctx, a0)
	return []any{r0}, err
}

// watchesPeekedCall calls Peeked on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesPeekedCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(Update)
	return nil, subject.Peeked(ctx, a0, a1)
}

// watchesFindsCall calls Finds on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesFindsCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Finds(ctx, a0)
	return []any{r0}, err
}

// watchesFoundCall calls Found on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesFoundCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(Update)
	return nil, subject.Found(ctx, a0, a1)
}

// watchesReadsCall calls Reads on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesReadsCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Reads(ctx, a0)
	return []any{r0}, err
}

// watchesWrittenCall calls Written on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesWrittenCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(Update)
	return nil, subject.Written(ctx, a0, a1)
}

// watchesCancelsCall calls Cancels on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesCancelsCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Cancels(ctx, a0)
	return []any{r0}, err
}

// watchesCancelledCall calls Cancelled on subject with the arguments args,
// and returns the values of the call beside its error.
func watchesCancelledCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(Update)
	return nil, subject.Cancelled(ctx, a0, a1)
}

// watchesClosesCall calls Closes on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesClosesCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Closes(ctx, a0)
	return []any{r0}, err
}

// watchesClosedCall calls Closed on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesClosedCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(Update)
	return nil, subject.Closed(ctx, a0, a1)
}

// watchesFieldsCall calls Fields on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesFieldsCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Fields(ctx, a0)
	return []any{r0}, err
}

// watchesFilledCall calls Filled on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesFilledCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(Update)
	return nil, subject.Filled(ctx, a0, a1)
}

// watchesFollowsCall calls Follows on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesFollowsCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Follows(ctx, a0)
	return []any{r0}, err
}

// watchesTouchCall calls Touch on subject with the arguments args, and
// returns the values of the call beside its error.
func watchesTouchCall(ctx context.Context, subject Watches, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Touch(ctx, a0)
}
