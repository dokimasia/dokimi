package memory

import (
	"testing"

	"go.dokimi.dev/conform"
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

// TestForgetfulConformsToStore applies the checks of the suite of Store to
// Forgetful.
func TestForgetfulConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuite(
		t,
		NewForgetful,
		conform.Excuse(
			StoreChecks.Get.KeepsTheItem(),
			"the store drops each item",
		),
		conform.Excuse(
			StoreChecks.Put.RefusesACancelledPut(),
			"the store drops each item before it reads the context",
		),
	)
}
