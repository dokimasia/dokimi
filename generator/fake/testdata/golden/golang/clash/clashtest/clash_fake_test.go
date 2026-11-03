package clashtest_test

import (
	"testing"

	"example.com/shop/clash/clashtest"
)

// TestShadowFake applies the checks of the suite of Shadow to ShadowFake.
func TestShadowFake(t *testing.T) {
	t.Parallel()

	clashtest.RunShadowSuite(t, clashtest.NewShadowFake)
}
