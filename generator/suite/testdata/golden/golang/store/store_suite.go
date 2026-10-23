package store

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

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

// KeepsTheItem returns the identity of the check that Get returns the item
// that Put stored.
func (c storeGet) KeepsTheItem() conform.ID[Store] {
	return c.Check("keeps-the-item")
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

// StoreSuite returns the checks of Store, and the calls and the budgets of
// its callables.
func StoreSuite() conform.Suite[Store] {
	return conform.Suite[Store]{
		Declaration: "Store",
		Checks: []conform.Check[Store]{
			{
				ID:        StoreChecks.Get.KeepsTheItem(),
				Statement: "Get returns the item that Put stored",
				Apply:     KeepsTheItem,
				Defect:    forgetful,
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
	conform.Run(t, StoreSuite(), func(h *StoreHandle) Store { return newSubject(h) }, opts...)
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
	conform.Fuzz(
		f,
		StoreSuite(),
		callable,
		func(h *StoreHandle) Store { return newSubject(h) },
		opts...)
}

// RunStoreBudgets applies the allocation ceilings of Store, and the
// ceilings that opts add, to a fresh subject of newSubject for each case.
func RunStoreBudgets[T Store](
	t *testing.T,
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	t.Helper()
	conform.RunBudgets(
		t,
		StoreSuite(),
		func(h *StoreHandle) Store { return newSubject(h) },
		opts...)
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
	conform.Bench(
		b,
		StoreSuite(),
		callable,
		func(h *StoreHandle) Store { return newSubject(h) },
		opts...)
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
