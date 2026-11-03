package cachetest_test

import (
	"testing"

	"example.com/shop/cache/cachetest"
)

// TestCacheFake applies the checks of the suite of Cache to CacheFake.
func TestCacheFake(t *testing.T) {
	t.Parallel()

	cachetest.RunCacheSuite(t, cachetest.NewCacheFake)
}
