package memory

import (
	"testing"
	"time"

	"example.com/shop/store"
	"go.dokimi.dev/conform"
)

// TestStoreConformsToStore applies the checks of the suite of Store to
// Store.
func TestStoreConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuite(t, New)
}

// BenchmarkStore measures each callable of the suite of Store on Store.
func BenchmarkStore(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Get.Callable, func(*StoreHandle) *Store { return New() })
	})
	b.Run("Put", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Put.Callable, func(*StoreHandle) *Store { return New() })
	})
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

	RunStoreSuiteFromHandle(
		t,
		NewClocked,
		conform.Cases[store.Store](20),
	)
}

// TestClockedAllocs applies the allocation ceilings of the suite of Store
// to Clocked.
func TestClockedAllocs(t *testing.T) {
	RunStoreBudgets(
		t,
		NewClocked,
		conform.Cases[store.Store](20),
		conform.Ceilings(
			StoreChecks.Get.Callable,
			conform.Budget{
				Allocs:  new(uint64(2)),
				Latency: time.Millisecond,
			},
		),
	)
}

// BenchmarkClocked measures each callable of the suite of Store on
// Clocked.
func BenchmarkClocked(b *testing.B) {
	opts := []StoreOption{
		conform.Ceilings(
			StoreChecks.Get.Callable,
			conform.Budget{
				Allocs:  new(uint64(2)),
				Latency: time.Millisecond,
			},
		),
	}
	b.Run("Get", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Get.Callable, NewClocked, opts...)
	})
	b.Run("Put", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Put.Callable, NewClocked, opts...)
	})
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

// BenchmarkForgetful measures each callable of the suite of Store on
// Forgetful.
func BenchmarkForgetful(b *testing.B) {
	opts := []StoreOption{
		conform.Excuse(
			StoreChecks.Get.KeepsTheItem(),
			"the store drops each item",
		),
		conform.Excuse(
			StoreChecks.Put.RefusesACancelledPut(),
			"the store drops each item before it reads the context",
		),
	}
	b.Run("Get", func(b *testing.B) {
		BenchStoreSuite(
			b,
			StoreChecks.Get.Callable,
			func(*StoreHandle) *Forgetful { return NewForgetful() },
			opts...)
	})
	b.Run("Put", func(b *testing.B) {
		BenchStoreSuite(
			b,
			StoreChecks.Put.Callable,
			func(*StoreHandle) *Forgetful { return NewForgetful() },
			opts...)
	})
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
