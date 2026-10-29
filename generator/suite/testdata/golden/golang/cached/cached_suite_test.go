package cached

import (
	"testing"

	"go.dokimi.dev/conform"
)

// TestCacheConformsToCache applies the checks of the suite of Cache to
// Cache.
func TestCacheConformsToCache(t *testing.T) {
	t.Parallel()

	RunCacheSuite(
		t,
		New,
		conform.Excuse(
			CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	)
}

// BenchmarkCacheConformsToCache measures each callable of the suite of
// Cache on Cache.
func BenchmarkCacheConformsToCache(b *testing.B) {
	opts := []CacheOption{
		conform.Excuse(
			CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	}
	b.Run("Get", func(b *testing.B) {
		BenchCacheSuite(
			b,
			CacheChecks.Get.Callable,
			func(*CacheHandle) *Cache { return New() },
			opts...)
	})
	b.Run("Set", func(b *testing.B) {
		BenchCacheSuite(
			b,
			CacheChecks.Set.Callable,
			func(*CacheHandle) *Cache { return New() },
			opts...)
	})
}

// FuzzCacheConformsToCacheGet fuzzes the property checks of Cache.Get on
// Cache.
func FuzzCacheConformsToCacheGet(f *testing.F) {
	FuzzCacheSuite(
		f,
		CacheChecks.Get.Callable,
		func(*CacheHandle) *Cache { return New() },
		conform.Excuse(
			CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	)
}

// FuzzCacheConformsToCacheSet fuzzes the property checks of Cache.Set on
// Cache.
func FuzzCacheConformsToCacheSet(f *testing.F) {
	FuzzCacheSuite(
		f,
		CacheChecks.Set.Callable,
		func(*CacheHandle) *Cache { return New() },
		conform.Excuse(
			CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	)
}
