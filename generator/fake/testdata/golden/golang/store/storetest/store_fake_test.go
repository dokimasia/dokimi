package storetest_test

import (
	"testing"

	"example.com/shop/store/storetest"
)

// TestStoreFake applies the checks of the suite of Store to StoreFake.
func TestStoreFake(t *testing.T) {
	t.Parallel()

	storetest.RunStoreSuite(t, storetest.NewStoreFake)
}
