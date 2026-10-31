package cached_test

import (
	"testing"

	"example.com/shop/cached"
	"example.com/shop/generic"
	"go.dokimi.dev/conform"
)

// TestCacheConformsToCache applies the checks of the suite of Cache to
// Cache.
func TestCacheConformsToCache(t *testing.T) {
	t.Parallel()

	generic.RunCacheSuite(
		t,
		cached.New,
		conform.Excuse(
			generic.CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	)
}

// BenchmarkCacheConformsToCache measures each callable of the suite of
// Cache on Cache.
func BenchmarkCacheConformsToCache(b *testing.B) {
	opts := []generic.CacheOption{
		conform.Excuse(
			generic.CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	}
	b.Run("Get", func(b *testing.B) {
		generic.BenchCacheSuite(
			b,
			generic.CacheChecks.Get.Callable,
			func(*generic.CacheHandle) *cached.Cache { return cached.New() },
			opts...)
	})
	b.Run("Set", func(b *testing.B) {
		generic.BenchCacheSuite(
			b,
			generic.CacheChecks.Set.Callable,
			func(*generic.CacheHandle) *cached.Cache { return cached.New() },
			opts...)
	})
}

// FuzzCacheConformsToCacheGet fuzzes the property checks of Cache.Get on
// Cache.
func FuzzCacheConformsToCacheGet(f *testing.F) {
	generic.FuzzCacheSuite(
		f,
		generic.CacheChecks.Get.Callable,
		func(*generic.CacheHandle) *cached.Cache { return cached.New() },
		conform.Excuse(
			generic.CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	)
}

// FuzzCacheConformsToCacheSet fuzzes the property checks of Cache.Set on
// Cache.
func FuzzCacheConformsToCacheSet(f *testing.F) {
	generic.FuzzCacheSuite(
		f,
		generic.CacheChecks.Set.Callable,
		func(*generic.CacheHandle) *cached.Cache { return cached.New() },
		conform.Excuse(
			generic.CacheChecks.Get.ZeroOnMiss(),
			"the cache returns the zero value of int for each miss already",
		),
	)
}
