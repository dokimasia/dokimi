package cached

import (
	"testing"
)

// TestCacheConformsToCache applies the checks of the suite of Cache to
// Cache.
func TestCacheConformsToCache(t *testing.T) {
	t.Parallel()

	RunCacheSuite(t, New)
}

// BenchmarkCacheConformsToCache measures each callable of the suite of
// Cache on Cache.
func BenchmarkCacheConformsToCache(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		BenchCacheSuite(b, CacheChecks.Get.Callable, func(*CacheHandle) *Cache { return New() })
	})
	b.Run("Set", func(b *testing.B) {
		BenchCacheSuite(b, CacheChecks.Set.Callable, func(*CacheHandle) *Cache { return New() })
	})
}

// FuzzCacheConformsToCacheGet fuzzes the property checks of Cache.Get on
// Cache.
func FuzzCacheConformsToCacheGet(f *testing.F) {
	FuzzCacheSuite(f, CacheChecks.Get.Callable, func(*CacheHandle) *Cache { return New() })
}

// FuzzCacheConformsToCacheSet fuzzes the property checks of Cache.Set on
// Cache.
func FuzzCacheConformsToCacheSet(f *testing.F) {
	FuzzCacheSuite(f, CacheChecks.Set.Callable, func(*CacheHandle) *Cache { return New() })
}
