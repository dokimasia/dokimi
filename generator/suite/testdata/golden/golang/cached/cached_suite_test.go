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

// FuzzCacheGet fuzzes the property checks of Cache.Get on Cache.
func FuzzCacheGet(f *testing.F) {
	FuzzCacheSuite(f, CacheChecks.Get.Callable, func(*CacheHandle) *Cache { return New() })
}

// FuzzCacheSet fuzzes the property checks of Cache.Set on Cache.
func FuzzCacheSet(f *testing.F) {
	FuzzCacheSuite(f, CacheChecks.Set.Callable, func(*CacheHandle) *Cache { return New() })
}
