package store

import (
	"testing"

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

// StoreSuite returns the checks of Store.
func StoreSuite() conform.Suite[Store] {
	return conform.Suite[Store]{Declaration: "Store", Checks: []conform.Check[Store]{
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
	}}
}

// RunStoreSuite applies the checks of Store, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunStoreSuite[T Store](t *testing.T, newSubject func() T, opts ...StoreOption) {
	t.Helper()
	RunStoreSuiteFromHandle(t, func(*StoreHandle) T { return newSubject() }, opts...)
}

// RunStoreSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStoreSuiteFromHandle[T Store](t *testing.T, newSubject func(h *StoreHandle) T, opts ...StoreOption) {
	t.Helper()
	conform.Run(t, StoreSuite(), func(h *StoreHandle) Store { return newSubject(h) }, opts...)
}
