package memory

import (
	"testing"

	"go.dokimi.dev/conform"
)

// TestStoreConformsToStore applies the checks of the suite of Store to
// Store.
func TestStoreConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuite(t, New)
}

// FuzzStoreGet fuzzes the property checks of Store.Get on Store.
func FuzzStoreGet(f *testing.F) {
	FuzzStoreSuite(f, StoreChecks.Get.Callable, func(*StoreHandle) *Store { return New() })
}

// FuzzStorePut fuzzes the property checks of Store.Put on Store.
func FuzzStorePut(f *testing.F) {
	FuzzStoreSuite(f, StoreChecks.Put.Callable, func(*StoreHandle) *Store { return New() })
}

// TestClockedConformsToStore applies the checks of the suite of Store to
// Clocked.
func TestClockedConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuiteFromHandle(t, NewClocked)
}

// FuzzClockedGet fuzzes the property checks of Store.Get on Clocked.
func FuzzClockedGet(f *testing.F) {
	FuzzStoreSuite(f, StoreChecks.Get.Callable, NewClocked)
}

// FuzzClockedPut fuzzes the property checks of Store.Put on Clocked.
func FuzzClockedPut(f *testing.F) {
	FuzzStoreSuite(f, StoreChecks.Put.Callable, NewClocked)
}

// TestForgetfulConformsToStore applies the checks of the suite of Store to
// Forgetful.
func TestForgetfulConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuite(
		t,
		NewForgetful,
		conform.Excuse(
			StoreChecks.Get.KeepsTheItem(),
			"the store drops each item",
		),
		conform.Excuse(
			StoreChecks.Put.RefusesACancelledPut(),
			"the store drops each item before it reads the context",
		),
	)
}

// FuzzForgetfulGet fuzzes the property checks of Store.Get on Forgetful.
func FuzzForgetfulGet(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Get.Callable,
		func(*StoreHandle) *Forgetful { return NewForgetful() },
		conform.Excuse(
			StoreChecks.Get.KeepsTheItem(),
			"the store drops each item",
		),
		conform.Excuse(
			StoreChecks.Put.RefusesACancelledPut(),
			"the store drops each item before it reads the context",
		),
	)
}

// FuzzForgetfulPut fuzzes the property checks of Store.Put on Forgetful.
func FuzzForgetfulPut(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Put.Callable,
		func(*StoreHandle) *Forgetful { return NewForgetful() },
		conform.Excuse(
			StoreChecks.Get.KeepsTheItem(),
			"the store drops each item",
		),
		conform.Excuse(
			StoreChecks.Put.RefusesACancelledPut(),
			"the store drops each item before it reads the context",
		),
	)
}
