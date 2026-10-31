package memory_test

import (
	"testing"
	"time"

	"example.com/shop/memory"
	"example.com/shop/store"
	"go.dokimi.dev/conform"
)

// TestStoreConformsToStore applies the checks of the suite of Store to
// Store.
func TestStoreConformsToStore(t *testing.T) {
	t.Parallel()

	store.RunStoreSuite(t, memory.New)
}

// BenchmarkStoreConformsToStore measures each callable of the suite of
// Store on Store.
func BenchmarkStoreConformsToStore(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		store.BenchStoreSuite(
			b,
			store.StoreChecks.Get.Callable,
			func(*store.StoreHandle) *memory.Store { return memory.New() },
		)
	})
	b.Run("Put", func(b *testing.B) {
		store.BenchStoreSuite(
			b,
			store.StoreChecks.Put.Callable,
			func(*store.StoreHandle) *memory.Store { return memory.New() },
		)
	})
}

// FuzzStoreConformsToStoreGet fuzzes the property checks of Store.Get on
// Store.
func FuzzStoreConformsToStoreGet(f *testing.F) {
	store.FuzzStoreSuite(
		f,
		store.StoreChecks.Get.Callable,
		func(*store.StoreHandle) *memory.Store { return memory.New() },
	)
}

// FuzzStoreConformsToStorePut fuzzes the property checks of Store.Put on
// Store.
func FuzzStoreConformsToStorePut(f *testing.F) {
	store.FuzzStoreSuite(
		f,
		store.StoreChecks.Put.Callable,
		func(*store.StoreHandle) *memory.Store { return memory.New() },
	)
}

// TestClockedConformsToStore applies the checks of the suite of Store to
// Clocked.
func TestClockedConformsToStore(t *testing.T) {
	t.Parallel()

	store.RunStoreSuiteFromHandle(
		t,
		memory.NewClocked,
		conform.Cases[store.Store](20),
	)
}

// TestClockedConformsToStoreAllocs applies the allocation ceilings and the
// checks that count allocations of the suite of Store to Clocked.
func TestClockedConformsToStoreAllocs(t *testing.T) {
	store.RunStoreBudgets(
		t,
		memory.NewClocked,
		conform.Cases[store.Store](20),
		conform.Ceilings(
			store.StoreChecks.Get.Callable,
			conform.Budget{
				Allocs:  new(uint64(2)),
				Latency: time.Millisecond,
			},
		),
	)
}

// BenchmarkClockedConformsToStore measures each callable of the suite of
// Store on Clocked.
func BenchmarkClockedConformsToStore(b *testing.B) {
	opts := []store.StoreOption{
		conform.Ceilings(
			store.StoreChecks.Get.Callable,
			conform.Budget{
				Allocs:  new(uint64(2)),
				Latency: time.Millisecond,
			},
		),
	}
	b.Run("Get", func(b *testing.B) {
		store.BenchStoreSuite(b, store.StoreChecks.Get.Callable, memory.NewClocked, opts...)
	})
	b.Run("Put", func(b *testing.B) {
		store.BenchStoreSuite(b, store.StoreChecks.Put.Callable, memory.NewClocked, opts...)
	})
}

// FuzzClockedConformsToStoreGet fuzzes the property checks of Store.Get on
// Clocked.
func FuzzClockedConformsToStoreGet(f *testing.F) {
	store.FuzzStoreSuite(f, store.StoreChecks.Get.Callable, memory.NewClocked)
}

// FuzzClockedConformsToStorePut fuzzes the property checks of Store.Put on
// Clocked.
func FuzzClockedConformsToStorePut(f *testing.F) {
	store.FuzzStoreSuite(f, store.StoreChecks.Put.Callable, memory.NewClocked)
}
