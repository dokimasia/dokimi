package refused

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of HybridClock:
//
// HybridClock.Update/smoke: Update returns without a panic for derived
// inputs.
//
// HybridClock.Update/smoke#generated: Update returns without a panic for
// derived inputs.
//
// HybridClock.MaxOffset/smoke: MaxOffset returns without a panic for
// derived inputs.
//
// HybridClock.MaxOffset/cancel: MaxOffset returns an error for a cancelled
// context.
//
// HybridClock.MaxOffset/zero-on-error: MaxOffset returns the zero value
// beside an error.

// The derived checks of HybridClock that witness their laws at the shallow
// end:
//
// HybridClock.Update/smoke
//
// HybridClock.MaxOffset/smoke

// The refused claims of HybridClock:
//
// HybridClock.Update/cancel: Update returns before it reads its context.
// The remedy is to remove the directive refuse of the claim where Update
// makes it.
//
// HybridClock.Update/zero-on-error: Update fails only for an instant
// beyond the max offset of the clock. The remedy is to remove the
// directive refuse of the claim where Update makes it.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of HybridClock writes.
type (
	HybridClockHandle = conform.Handle[HybridClock]
	HybridClockCheck  = conform.Check[HybridClock]
	HybridClockOption = conform.Option[HybridClock]
)

// HybridClockChecks contains the name of each callable of HybridClock,
// with the identities of its checks.
var HybridClockChecks = struct {
	Update    hybridClockUpdate
	MaxOffset hybridClockMaxOffset
}{
	Update:    hybridClockUpdate{"HybridClock.Update"},
	MaxOffset: hybridClockMaxOffset{"HybridClock.MaxOffset"},
}

// hybridClockUpdate is the name of HybridClock.Update, with the identities
// of its checks.
type hybridClockUpdate struct {
	conform.Callable[HybridClock]
}

// Smoke returns the identity of the check that Update returns without a
// panic for derived inputs.
func (c hybridClockUpdate) Smoke() conform.ID[HybridClock] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Update returns
// without a panic for derived inputs.
func (c hybridClockUpdate) SmokeGenerated() conform.ID[HybridClock] {
	return c.Check("smoke#generated")
}

// hybridClockMaxOffset is the name of HybridClock.MaxOffset, with the
// identities of its checks.
type hybridClockMaxOffset struct {
	conform.Callable[HybridClock]
}

// Smoke returns the identity of the check that MaxOffset returns without a
// panic for derived inputs.
func (c hybridClockMaxOffset) Smoke() conform.ID[HybridClock] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that MaxOffset returns an error
// for a cancelled context.
func (c hybridClockMaxOffset) Cancel() conform.ID[HybridClock] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that MaxOffset returns the
// zero value beside an error.
func (c hybridClockMaxOffset) ZeroOnError() conform.ID[HybridClock] {
	return c.Check("zero-on-error")
}

// HybridClockSuite returns the checks of HybridClock, and the calls and
// the budgets of its callables.
func HybridClockSuite() conform.Suite[HybridClock] {
	return conform.Suite[HybridClock]{
		Declaration: "HybridClock",
		Package:     "example.com/shop/refused",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[HybridClock]{
			{
				ID:          HybridClockChecks.Update.Smoke(),
				Statement:   "Update returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HybridClockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Update(ctx, int64(42))
						return err
					})
				},
				Defect: hybridClockUpdatePanics,
			},
			{
				ID:          HybridClockChecks.Update.SmokeGenerated(),
				Statement:   "Update returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *HybridClockHandle, c *prop.Case) {
					drawnObserved := c.Draw(prop.Of[int64](), "observed")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Update(ctx, drawnObserved)
						return err
					})
				},
				Defect: hybridClockUpdatePanics,
			},
			{
				ID:          HybridClockChecks.MaxOffset.Smoke(),
				Statement:   "MaxOffset returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HybridClockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().MaxOffset(ctx)
						return err
					})
				},
				Defect: hybridClockMaxOffsetPanics,
			},
			{
				ID:          HybridClockChecks.MaxOffset.Cancel(),
				Statement:   "MaxOffset returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *HybridClockHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().MaxOffset(ctx)
						return err
					})
				},
				Defect: hybridClockMaxOffsetIgnoresCancel,
			},
			{
				ID:          HybridClockChecks.MaxOffset.ZeroOnError(),
				Statement:   "MaxOffset returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *HybridClockHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int64, error) {
						return h.Subject().MaxOffset(ctx)
					}, int64(0))
				},
				Defect: hybridClockMaxOffsetValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "HybridClock.Update",
				Gate:     spec.Unrefused,
				Why:      "Update returns before it reads its context",
				Remedy:   "remove the directive refuse of the claim where Update makes it",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "HybridClock.Update",
				Gate:     spec.Unrefused,
				Why:      "Update fails only for an instant beyond the max offset of the clock",
				Remedy:   "remove the directive refuse of the claim where Update makes it",
			},
		},
		Calls: []conform.Call[HybridClock]{
			{
				Calls:    hybridClockUpdateCalls,
				Over:     newHybridClockOver,
				Callable: HybridClockChecks.Update.Callable,
			},
			{
				Calls:    hybridClockMaxOffsetCalls,
				Over:     newHybridClockOver,
				Callable: HybridClockChecks.MaxOffset.Callable,
			},
		},
	}
}

// RunHybridClockSuite applies the checks of HybridClock, and the checks
// that opts add, to a fresh subject of newSubject for each check.
func RunHybridClockSuite[T HybridClock](
	t *testing.T,
	newSubject func() T,
	opts ...HybridClockOption,
) {
	t.Helper()
	RunHybridClockSuiteFromHandle(t, func(*HybridClockHandle) T { return newSubject() }, opts...)
}

// RunHybridClockSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunHybridClockSuiteFromHandle[T HybridClock](
	t *testing.T,
	newSubject func(h *HybridClockHandle) T,
	opts ...HybridClockOption,
) {
	t.Helper()
	conform.Run(t, HybridClockSuite(), hybridClockSubjects(newSubject), opts...)
}

// FuzzHybridClockSuite registers on f the fuzz target of the property
// checks of callable, on a subject that newSubject builds from the handle
// of each case.
func FuzzHybridClockSuite[T HybridClock](
	f *testing.F,
	callable conform.Callable[HybridClock],
	newSubject func(h *HybridClockHandle) T,
	opts ...HybridClockOption,
) {
	f.Helper()
	conform.Fuzz(f, HybridClockSuite(), callable, hybridClockSubjects(newSubject), opts...)
}

// RunHybridClockBudgets applies the allocation ceilings of HybridClock and
// its checks that count allocations, and the ceilings and the checks that
// opts add, to a fresh subject of newSubject for each case.
func RunHybridClockBudgets[T HybridClock](
	t *testing.T,
	newSubject func(h *HybridClockHandle) T,
	opts ...HybridClockOption,
) {
	t.Helper()
	conform.RunBudgets(t, HybridClockSuite(), hybridClockSubjects(newSubject), opts...)
}

// BenchHybridClockSuite measures callable on a subject of newSubject,
// under the ceilings of its budget.
func BenchHybridClockSuite[T HybridClock](
	b conform.B,
	callable conform.Callable[HybridClock],
	newSubject func(h *HybridClockHandle) T,
	opts ...HybridClockOption,
) {
	b.Helper()
	conform.Bench(b, HybridClockSuite(), callable, hybridClockSubjects(newSubject), opts...)
}

// hybridClockSubjects returns the factory of the subjects of HybridClock
// that calls newSubject with the handle of each check.
func hybridClockSubjects[T HybridClock](
	newSubject func(h *HybridClockHandle) T,
) func(h *HybridClockHandle) HybridClock {
	return func(h *HybridClockHandle) HybridClock { return newSubject(h) }
}

// hybridClockUpdateCalls generates the calls of HybridClock.Update with
// the inputs that each case draws.
var hybridClockUpdateCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, HybridClock) error {
		observed := c.Draw(prop.Of[int64](), "observed")
		return func(ctx context.Context, subject HybridClock) error {
			_, err := subject.Update(ctx, observed)
			return err
		}
	},
)

// hybridClockMaxOffsetCalls generates the calls of HybridClock.MaxOffset
// with the inputs that each case draws.
var hybridClockMaxOffsetCalls = prop.Composite(
	func(*prop.Case) func(context.Context, HybridClock) error {
		return func(ctx context.Context, subject HybridClock) error {
			_, err := subject.MaxOffset(ctx)
			return err
		}
	},
)

// hybridClockOverEmbedded is HybridClock under the name of the field that
// hybridClockOver embeds.
type hybridClockOverEmbedded = HybridClock

// hybridClockOver is a subject of HybridClock whose callables call spend
// and return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type hybridClockOver struct {
	hybridClockOverEmbedded
	spend func()
}

// Update calls spend and returns its zero results.
func (o hybridClockOver) Update(context.Context, int64) (r0 int64, r1 error) {
	o.spend()
	return r0, r1
}

// MaxOffset calls spend and returns its zero results.
func (o hybridClockOver) MaxOffset(context.Context) (r0 int64, r1 error) {
	o.spend()
	return r0, r1
}

// newHybridClockOver returns the planted defect of a ceiling of
// HybridClock, whose callables call spend.
func newHybridClockOver(_ *HybridClockHandle, spend func()) HybridClock {
	return hybridClockOver{spend: spend}
}

// hybridClockUpdatePanicsDefect is the planted defect of a HybridClock
// whose Update panics.
type hybridClockUpdatePanicsDefect struct {
	hybridClockOverEmbedded
}

// Update is the override of the planted defect.
func (*hybridClockUpdatePanicsDefect) Update(p0 context.Context, p1 int64) (r0 int64, r1 error) {
	panic("the planted defect of a HybridClock whose Update panics")
}

// hybridClockUpdatePanics returns the planted defect of a HybridClock
// whose Update panics over a fresh subject.
func hybridClockUpdatePanics(h *HybridClockHandle) HybridClock {
	return &hybridClockUpdatePanicsDefect{hybridClockOverEmbedded: h.Fresh()}
}

// hybridClockMaxOffsetPanicsDefect is the planted defect of a HybridClock
// whose MaxOffset panics.
type hybridClockMaxOffsetPanicsDefect struct {
	hybridClockOverEmbedded
}

// MaxOffset is the override of the planted defect.
func (*hybridClockMaxOffsetPanicsDefect) MaxOffset(p0 context.Context) (r0 int64, r1 error) {
	panic("the planted defect of a HybridClock whose MaxOffset panics")
}

// hybridClockMaxOffsetPanics returns the planted defect of a HybridClock
// whose MaxOffset panics over a fresh subject.
func hybridClockMaxOffsetPanics(h *HybridClockHandle) HybridClock {
	return &hybridClockMaxOffsetPanicsDefect{hybridClockOverEmbedded: h.Fresh()}
}

// hybridClockMaxOffsetIgnoresCancelDefect is the planted defect of a
// HybridClock whose MaxOffset ignores the context.
type hybridClockMaxOffsetIgnoresCancelDefect struct {
	hybridClockOverEmbedded
}

// MaxOffset is the override of the planted defect.
func (d *hybridClockMaxOffsetIgnoresCancelDefect) MaxOffset(
	p0 context.Context,
) (r0 int64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.hybridClockOverEmbedded.MaxOffset(p0)
}

// hybridClockMaxOffsetIgnoresCancel returns the planted defect of a
// HybridClock whose MaxOffset ignores the context over a fresh subject.
func hybridClockMaxOffsetIgnoresCancel(h *HybridClockHandle) HybridClock {
	return &hybridClockMaxOffsetIgnoresCancelDefect{hybridClockOverEmbedded: h.Fresh()}
}

// hybridClockMaxOffsetValueBesideDefect is the planted defect of a
// HybridClock whose MaxOffset returns a value beside an error.
type hybridClockMaxOffsetValueBesideDefect struct {
	hybridClockOverEmbedded
}

// MaxOffset is the override of the planted defect.
func (d *hybridClockMaxOffsetValueBesideDefect) MaxOffset(p0 context.Context) (r0 int64, r1 error) {
	r0, r1 = d.hybridClockOverEmbedded.MaxOffset(p0)
	if r1 != nil {
		r0 = int64(7)
	}
	return r0, r1
}

// hybridClockMaxOffsetValueBeside returns the planted defect of a
// HybridClock whose MaxOffset returns a value beside an error over a fresh
// subject.
func hybridClockMaxOffsetValueBeside(h *HybridClockHandle) HybridClock {
	return &hybridClockMaxOffsetValueBesideDefect{hybridClockOverEmbedded: h.Fresh()}
}

// The derived checks of Doubled:
//
// Doubled.Read/smoke: Read returns without a panic for derived inputs.
//
// Doubled.Read/smoke#generated: Read returns without a panic for derived
// inputs.
//
// Doubled.Read/cancel: Read returns an error for a cancelled context.
//
// Doubled.Read/cancel#generated: Read returns an error for a cancelled
// context.

// The derived checks of Doubled that witness their laws at the shallow
// end:
//
// Doubled.Read/smoke

// The refused claims of Doubled:
//
// Doubled.Read/zero-on-error: Read never fails. The remedy is to remove
// the directive refuse of the claim where Read makes it.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Doubled writes.
type (
	DoubledHandle = conform.Handle[Doubled]
	DoubledCheck  = conform.Check[Doubled]
	DoubledOption = conform.Option[Doubled]
)

// DoubledChecks contains the name of each callable of Doubled, with the
// identities of its checks.
var DoubledChecks = struct {
	Read doubledRead
}{
	Read: doubledRead{"Doubled.Read"},
}

// doubledRead is the name of Doubled.Read, with the identities of its
// checks.
type doubledRead struct {
	conform.Callable[Doubled]
}

// Smoke returns the identity of the check that Read returns without a
// panic for derived inputs.
func (c doubledRead) Smoke() conform.ID[Doubled] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Read returns
// without a panic for derived inputs.
func (c doubledRead) SmokeGenerated() conform.ID[Doubled] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Read returns an error for
// a cancelled context.
func (c doubledRead) Cancel() conform.ID[Doubled] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Read returns an
// error for a cancelled context.
func (c doubledRead) CancelGenerated() conform.ID[Doubled] {
	return c.Check("cancel#generated")
}

// DoubledSuite returns the checks of Doubled, and the calls and the
// budgets of its callables.
func DoubledSuite() conform.Suite[Doubled] {
	return conform.Suite[Doubled]{
		Declaration: "Doubled",
		Package:     "example.com/shop/refused",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Doubled]{
			{
				ID:          DoubledChecks.Read.Smoke(),
				Statement:   "Read returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DoubledHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, string("key-a"))
						return err
					})
				},
				Defect: doubledReadPanics,
			},
			{
				ID:          DoubledChecks.Read.SmokeGenerated(),
				Statement:   "Read returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubledHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, drawnKey)
						return err
					})
				},
				Defect: doubledReadPanics,
			},
			{
				ID:          DoubledChecks.Read.Cancel(),
				Statement:   "Read returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DoubledHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, string("key-a"))
						return err
					})
				},
				Defect: doubledReadIgnoresCancel,
			},
			{
				ID:          DoubledChecks.Read.CancelGenerated(),
				Statement:   "Read returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubledHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, drawnKey)
						return err
					})
				},
				Defect: doubledReadIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ZeroOnError,
				Callable: "Doubled.Read",
				Gate:     spec.Unrefused,
				Why:      "Read never fails",
				Remedy:   "remove the directive refuse of the claim where Read makes it",
			},
		},
		Calls: []conform.Call[Doubled]{
			{
				Calls:    doubledReadCalls,
				Over:     newDoubledOver,
				Callable: DoubledChecks.Read.Callable,
			},
		},
	}
}

// RunDoubledSuite applies the checks of Doubled, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunDoubledSuite[T Doubled](t *testing.T, newSubject func() T, opts ...DoubledOption) {
	t.Helper()
	RunDoubledSuiteFromHandle(t, func(*DoubledHandle) T { return newSubject() }, opts...)
}

// RunDoubledSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDoubledSuiteFromHandle[T Doubled](
	t *testing.T,
	newSubject func(h *DoubledHandle) T,
	opts ...DoubledOption,
) {
	t.Helper()
	conform.Run(t, DoubledSuite(), doubledSubjects(newSubject), opts...)
}

// FuzzDoubledSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzDoubledSuite[T Doubled](
	f *testing.F,
	callable conform.Callable[Doubled],
	newSubject func(h *DoubledHandle) T,
	opts ...DoubledOption,
) {
	f.Helper()
	conform.Fuzz(f, DoubledSuite(), callable, doubledSubjects(newSubject), opts...)
}

// RunDoubledBudgets applies the allocation ceilings of Doubled and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunDoubledBudgets[T Doubled](
	t *testing.T,
	newSubject func(h *DoubledHandle) T,
	opts ...DoubledOption,
) {
	t.Helper()
	conform.RunBudgets(t, DoubledSuite(), doubledSubjects(newSubject), opts...)
}

// BenchDoubledSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchDoubledSuite[T Doubled](
	b conform.B,
	callable conform.Callable[Doubled],
	newSubject func(h *DoubledHandle) T,
	opts ...DoubledOption,
) {
	b.Helper()
	conform.Bench(b, DoubledSuite(), callable, doubledSubjects(newSubject), opts...)
}

// doubledSubjects returns the factory of the subjects of Doubled that
// calls newSubject with the handle of each check.
func doubledSubjects[T Doubled](
	newSubject func(h *DoubledHandle) T,
) func(h *DoubledHandle) Doubled {
	return func(h *DoubledHandle) Doubled { return newSubject(h) }
}

// doubledReadCalls generates the calls of Doubled.Read with the inputs
// that each case draws.
var doubledReadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Doubled) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Doubled) error {
		_, err := subject.Read(ctx, key)
		return err
	}
})

// doubledOverEmbedded is Doubled under the name of the field that
// doubledOver embeds.
type doubledOverEmbedded = Doubled

// doubledOver is a subject of Doubled whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type doubledOver struct {
	doubledOverEmbedded
	spend func()
}

// Read calls spend and returns its zero results.
func (o doubledOver) Read(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newDoubledOver returns the planted defect of a ceiling of Doubled, whose
// callables call spend.
func newDoubledOver(_ *DoubledHandle, spend func()) Doubled {
	return doubledOver{spend: spend}
}

// doubledReadPanicsDefect is the planted defect of a Doubled whose Read
// panics.
type doubledReadPanicsDefect struct {
	doubledOverEmbedded
}

// Read is the override of the planted defect.
func (*doubledReadPanicsDefect) Read(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Doubled whose Read panics")
}

// doubledReadPanics returns the planted defect of a Doubled whose Read
// panics over a fresh subject.
func doubledReadPanics(h *DoubledHandle) Doubled {
	return &doubledReadPanicsDefect{doubledOverEmbedded: h.Fresh()}
}

// doubledReadIgnoresCancelDefect is the planted defect of a Doubled whose
// Read ignores the context.
type doubledReadIgnoresCancelDefect struct {
	doubledOverEmbedded
}

// Read is the override of the planted defect.
func (d *doubledReadIgnoresCancelDefect) Read(p0 context.Context, p1 string) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.doubledOverEmbedded.Read(p0, p1)
}

// doubledReadIgnoresCancel returns the planted defect of a Doubled whose
// Read ignores the context over a fresh subject.
func doubledReadIgnoresCancel(h *DoubledHandle) Doubled {
	return &doubledReadIgnoresCancelDefect{doubledOverEmbedded: h.Fresh()}
}

// The derived checks of Title:
//
// Title/smoke: Title returns without a panic for derived inputs.
//
// Title/smoke#generated: Title returns without a panic for derived inputs.

// The derived checks of Title that witness their laws at the shallow end:
//
// Title/smoke

// The refused claims of Title:
//
// Title/deterministic: Title salts each name. The remedy is to remove the
// directive refuse of the claim where Title makes it.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Title writes.
type (
	TitleHandle = conform.Handle[TitleFunc]
	TitleCheck  = conform.Check[TitleFunc]
	TitleOption = conform.Option[TitleFunc]
)

// TitleFunc is the type of the function Title. Each subject of the suite
// of Title has it.
type TitleFunc func(name string) string

// TitleChecks contains the name of each callable of Title, with the
// identities of its checks.
var TitleChecks = titleCallable{"Title"}

// titleCallable is the name of Title, with the identities of its checks.
type titleCallable struct {
	conform.Callable[TitleFunc]
}

// Smoke returns the identity of the check that Title returns without a
// panic for derived inputs.
func (c titleCallable) Smoke() conform.ID[TitleFunc] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Title returns
// without a panic for derived inputs.
func (c titleCallable) SmokeGenerated() conform.ID[TitleFunc] {
	return c.Check("smoke#generated")
}

// TitleSuite returns the checks of Title, and the calls and the budgets of
// its callables.
func TitleSuite() conform.Suite[TitleFunc] {
	return conform.Suite[TitleFunc]{
		Declaration: "Title",
		Package:     "example.com/shop/refused",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[TitleFunc]{
			{
				ID:          TitleChecks.Smoke(),
				Statement:   "Title returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TitleHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject()(string("name-a"))
						return nil
					})
				},
				Defect: titlePanics,
			},
			{
				ID:          TitleChecks.SmokeGenerated(),
				Statement:   "Title returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TitleHandle, c *prop.Case) {
					drawnName := c.Draw(prop.Of[string](), "name")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject()(drawnName)
						return nil
					})
				},
				Defect: titlePanics,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Deterministic,
				Callable: "Title",
				Gate:     spec.Unrefused,
				Why:      "Title salts each name",
				Remedy:   "remove the directive refuse of the claim where Title makes it",
			},
		},
		Calls: []conform.Call[TitleFunc]{
			{
				Calls:    titleCalls,
				Over:     newTitleOver,
				Callable: TitleChecks.Callable,
			},
		},
	}
}

// RunTitleSuite applies the checks of Title, and the checks that opts add,
// to subject.
func RunTitleSuite(t *testing.T, subject TitleFunc, opts ...TitleOption) {
	t.Helper()
	RunTitleSuiteFromHandle(t, titleSubjects(subject), opts...)
}

// RunTitleSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTitleSuiteFromHandle(
	t *testing.T,
	newSubject func(h *TitleHandle) TitleFunc,
	opts ...TitleOption,
) {
	t.Helper()
	conform.Run(t, TitleSuite(), newSubject, opts...)
}

// FuzzTitleSuite registers on f the fuzz target of the property checks of
// Title on subject.
func FuzzTitleSuite(f *testing.F, subject TitleFunc, opts ...TitleOption) {
	f.Helper()
	conform.Fuzz(f, TitleSuite(), TitleChecks.Callable, titleSubjects(subject), opts...)
}

// RunTitleBudgets applies the allocation ceiling of Title and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunTitleBudgets(t *testing.T, subject TitleFunc, opts ...TitleOption) {
	t.Helper()
	conform.RunBudgets(t, TitleSuite(), titleSubjects(subject), opts...)
}

// BenchTitleSuite measures subject under the ceilings of the budget of
// Title.
func BenchTitleSuite(b conform.B, subject TitleFunc, opts ...TitleOption) {
	b.Helper()
	conform.Bench(b, TitleSuite(), TitleChecks.Callable, titleSubjects(subject), opts...)
}

// titleSubjects returns the factory of the subjects of Title that returns
// subject for each handle.
func titleSubjects(subject TitleFunc) func(*TitleHandle) TitleFunc {
	return func(*TitleHandle) TitleFunc { return subject }
}

// titleCalls generates the calls of Title with the inputs that each case
// draws.
var titleCalls = prop.Composite(func(c *prop.Case) func(context.Context, TitleFunc) error {
	name := c.Draw(prop.Of[string](), "name")
	return func(_ context.Context, subject TitleFunc) error {
		_ = subject(name)
		return nil
	}
})

// newTitleOver returns the planted defect of a ceiling of Title, which
// calls spend and returns its zero results.
func newTitleOver(_ *TitleHandle, spend func()) TitleFunc {
	return func(string) (r0 string) {
		spend()
		return r0
	}
}

// titlePanics returns the planted defect of a Title whose Title panics
// over the fresh function.
func titlePanics(h *TitleHandle) TitleFunc {
	return func(p0 string) (r0 string) {
		panic("the planted defect of a Title whose Title panics")
	}
}
