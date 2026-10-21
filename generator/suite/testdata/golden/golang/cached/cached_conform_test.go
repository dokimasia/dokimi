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
