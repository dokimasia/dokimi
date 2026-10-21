package memory

import (
	"testing"
)

// TestStoreConformsToStore applies the checks of the suite of Store to
// Store.
func TestStoreConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuite(t, New)
}

// TestClockedConformsToStore applies the checks of the suite of Store to
// Clocked.
func TestClockedConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuiteFromHandle(t, NewClocked)
}
