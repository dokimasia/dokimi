package store

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// storeSubject is a value of Store that only TestStoreSuite builds.
type storeSubject struct{ Store }

// TestStoreSuite checks the suite of Store and the two forms of its
// application.
func TestStoreSuite(t *testing.T) {
	t.Parallel()

	t.Run("StoreSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, StoreSuite())
	})

	t.Run("RunStoreSuite", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		RunStoreSuite(
			t,
			func() *storeSubject {
				built = &storeSubject{}
				return built
			},
			conform.With(StoreCheck{
				ID:        "Store/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *StoreHandle) {
					got := h.Subject()
					assert.Equal(h, got, Store(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*StoreHandle) Store { return &storeSubject{} },
			}),
			conform.Excuse(
				StoreChecks.Get.KeepsTheItem(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Put.RefusesACancelledPut(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStoreSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		RunStoreSuiteFromHandle(
			t,
			func(*StoreHandle) *storeSubject {
				built = &storeSubject{}
				return built
			},
			conform.With(StoreCheck{
				ID:        "Store/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *StoreHandle) {
					got := h.Subject()
					assert.Equal(h, got, Store(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*StoreHandle) Store { return &storeSubject{} },
			}),
			conform.Excuse(
				StoreChecks.Get.KeepsTheItem(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Put.RefusesACancelledPut(),
				"the test applies the check of its option alone",
			),
		)
	})
}
