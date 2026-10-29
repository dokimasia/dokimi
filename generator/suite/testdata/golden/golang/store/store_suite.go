package store

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Store:
//
// Store.Get/smoke: Get returns without a panic for derived inputs.
//
// Store.Get/cancel: Get returns an error for a cancelled context.
//
// Store.Get/zero-on-error: Get returns the zero value beside an error.
//
// Store.Put/smoke: Put returns without a panic for derived inputs.
//
// Store.Put/cancel: Put returns an error for a cancelled context.

// The derived checks of Store that witness their laws at the shallow end:
//
// Store.Get/smoke
//
// Store.Put/smoke

// The unclaimed classifications of Store:
//
// The shape:writer of Store.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Store writes.
type (
	StoreHandle = conform.Handle[Store]
	StoreCheck  = conform.Check[Store]
	StoreOption = conform.Option[Store]
)

// StoreChecks contains the name of each callable of Store, with the
// identities of its checks.
var StoreChecks = struct {
	Get storeGet
	Put storePut
}{
	Get: storeGet{"Store.Get"},
	Put: storePut{"Store.Put"},
}

// storeGet is the name of Store.Get, with the identities of its checks.
type storeGet struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c storeGet) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c storeGet) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c storeGet) ZeroOnError() conform.ID[Store] {
	return c.Check("zero-on-error")
}

// storePut is the name of Store.Put, with the identities of its checks.
type storePut struct {
	conform.Callable[Store]
}

// RefusesACancelledPut returns the identity of the check that Put reports
// a cancelled context.
func (c storePut) RefusesACancelledPut() conform.ID[Store] {
	return c.Check("refuses-a-cancelled-put")
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c storePut) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c storePut) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// StoreSuite returns the checks of Store, and the calls and the budgets of
// its callables.
func StoreSuite() conform.Suite[Store] {
	return conform.Suite[Store]{
		Declaration: "Store",
		Checks: []conform.Check[Store]{
			{
				ID:          StoreChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("id-a"))
						return err
					})
				},
				Defect: storeGetPanics,
			},
			{
				ID:          StoreChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("id-a"))
						return err
					})
				},
				Defect: storeGetIgnoresCancel,
			},
			{
				ID:          StoreChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("id-a"))
					}, Item{})
				},
				Defect: storeGetValueBeside,
			},
			{
				ID:          StoreChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{ID: "ID-a"})
					})
				},
				Defect: storePutPanics,
			},
			{
				ID:          StoreChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{ID: "ID-a"})
					})
				},
				Defect: storePutIgnoresCancel,
			},
			{
				ID:        StoreChecks.Put.RefusesACancelledPut(),
				Statement: "Put reports a cancelled context",
				Apply:     RefusesACancelledPut,
				Defect:    careless,
			},
		},
		Calls: []conform.Call[Store]{
			{
				Calls:    storeGetCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Get.Callable,
			},
			{
				Calls:    storePutCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Put.Callable,
			},
		},
	}
}

// RunStoreSuite applies the checks of Store, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunStoreSuite[T Store](t *testing.T, newSubject func() T, opts ...StoreOption) {
	t.Helper()
	RunStoreSuiteFromHandle(t, func(*StoreHandle) T { return newSubject() }, opts...)
}

// RunStoreSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStoreSuiteFromHandle[T Store](
	t *testing.T,
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	t.Helper()
	conform.Run(t, StoreSuite(), storeSubjects(newSubject), opts...)
}

// FuzzStoreSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzStoreSuite[T Store](
	f *testing.F,
	callable conform.Callable[Store],
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	f.Helper()
	conform.Fuzz(f, StoreSuite(), callable, storeSubjects(newSubject), opts...)
}

// RunStoreBudgets applies the allocation ceilings of Store and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunStoreBudgets[T Store](
	t *testing.T,
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	t.Helper()
	conform.RunBudgets(t, StoreSuite(), storeSubjects(newSubject), opts...)
}

// BenchStoreSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchStoreSuite[T Store](
	b conform.B,
	callable conform.Callable[Store],
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	b.Helper()
	conform.Bench(b, StoreSuite(), callable, storeSubjects(newSubject), opts...)
}

// storeSubjects returns the factory of the subjects of Store that calls
// newSubject with the handle of each check.
func storeSubjects[T Store](newSubject func(h *StoreHandle) T) func(h *StoreHandle) Store {
	return func(h *StoreHandle) Store { return newSubject(h) }
}

// storeGetCalls generates the calls of Store.Get with the inputs that each
// case draws.
var storeGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	id := c.Draw(prop.Of[string](), "id")
	return func(ctx context.Context, subject Store) error {
		_, err := subject.Get(ctx, id)
		return err
	}
})

// storePutCalls generates the calls of Store.Put with the inputs that each
// case draws.
var storePutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	item := c.Draw(prop.Of[Item](), "item")
	return func(ctx context.Context, subject Store) error {
		return subject.Put(ctx, item)
	}
})

// storeOverEmbedded is Store under the name of the field that storeOver
// embeds.
type storeOverEmbedded = Store

// storeOver is a subject of Store whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type storeOver struct {
	storeOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o storeOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Put calls spend and returns its zero results.
func (o storeOver) Put(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// newStoreOver returns the planted defect of a ceiling of Store, whose
// callables call spend.
func newStoreOver(_ *StoreHandle, spend func()) Store {
	return storeOver{spend: spend}
}

// storeGetPanicsDefect is the planted defect of a Store whose Get panics.
type storeGetPanicsDefect struct {
	storeOverEmbedded
}

// Get is the override of the planted defect.
func (*storeGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Store whose Get panics")
}

// storeGetPanics returns the planted defect of a Store whose Get panics
// over a fresh subject.
func storeGetPanics(h *StoreHandle) Store {
	return &storeGetPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetIgnoresCancelDefect is the planted defect of a Store whose Get
// ignores the context.
type storeGetIgnoresCancelDefect struct {
	storeOverEmbedded
}

// Get is the override of the planted defect.
func (d *storeGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.storeOverEmbedded.Get(p0, p1)
}

// storeGetIgnoresCancel returns the planted defect of a Store whose Get
// ignores the context over a fresh subject.
func storeGetIgnoresCancel(h *StoreHandle) Store {
	return &storeGetIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetValueBesideDefect is the planted defect of a Store whose Get
// returns a value beside an error.
type storeGetValueBesideDefect struct {
	storeOverEmbedded
}

// Get is the override of the planted defect.
func (d *storeGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.storeOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{ID: "ID-b"}
	}
	return r0, r1
}

// storeGetValueBeside returns the planted defect of a Store whose Get
// returns a value beside an error over a fresh subject.
func storeGetValueBeside(h *StoreHandle) Store {
	return &storeGetValueBesideDefect{storeOverEmbedded: h.Fresh()}
}

// storePutPanicsDefect is the planted defect of a Store whose Put panics.
type storePutPanicsDefect struct {
	storeOverEmbedded
}

// Put is the override of the planted defect.
func (*storePutPanicsDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Store whose Put panics")
}

// storePutPanics returns the planted defect of a Store whose Put panics
// over a fresh subject.
func storePutPanics(h *StoreHandle) Store {
	return &storePutPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storePutIgnoresCancelDefect is the planted defect of a Store whose Put
// ignores the context.
type storePutIgnoresCancelDefect struct {
	storeOverEmbedded
}

// Put is the override of the planted defect.
func (d *storePutIgnoresCancelDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.storeOverEmbedded.Put(p0, p1)
}

// storePutIgnoresCancel returns the planted defect of a Store whose Put
// ignores the context over a fresh subject.
func storePutIgnoresCancel(h *StoreHandle) Store {
	return &storePutIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}
