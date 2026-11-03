package generictest_test

import (
	"testing"

	"example.com/shop/generic/generictest"
)

// TestCacheFake applies the checks of the suite of Cache to CacheFake.
func TestCacheFake(t *testing.T) {
	t.Parallel()

	generictest.RunCacheSuite(t, generictest.NewCacheFake[int, string])
}
