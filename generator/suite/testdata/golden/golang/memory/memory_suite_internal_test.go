package memory

import (
	"testing"
)

// TestHiddenConformsToStore applies the checks of the suite of Store to
// Hidden.
func TestHiddenConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuite(t, newHidden)
}

// FuzzHiddenGet fuzzes the property checks of Store.Get on Hidden.
func FuzzHiddenGet(f *testing.F) {
	FuzzStoreSuite(f, StoreChecks.Get.Callable, func(*StoreHandle) *Hidden { return newHidden() })
}

// FuzzHiddenPut fuzzes the property checks of Store.Put on Hidden.
func FuzzHiddenPut(f *testing.F) {
	FuzzStoreSuite(f, StoreChecks.Put.Callable, func(*StoreHandle) *Hidden { return newHidden() })
}
