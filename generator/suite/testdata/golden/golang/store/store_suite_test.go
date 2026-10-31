package store_test

import (
	"testing"

	"example.com/shop/store"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// storeSubjectEmbedded is store.Store under the name of the field that
// storeSubject embeds.
type storeSubjectEmbedded = store.Store

// storeSubject is a value of Store that only TestStoreSuite builds.
type storeSubject struct{ storeSubjectEmbedded }

// TestStoreSuite checks the suite of Store and the two forms of its
// application.
func TestStoreSuite(t *testing.T) {
	t.Parallel()

	t.Run("StoreSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, store.StoreSuite())
	})

	t.Run("RunStoreSuite", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		store.RunStoreSuite(
			t,
			func() *storeSubject {
				built = &storeSubject{}
				return built
			},
			conform.With(store.StoreCheck{
				ID:        "Store/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *store.StoreHandle) {
					got := h.Subject()
					assert.Equal(h, got, store.Store(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*store.StoreHandle) store.Store { return &storeSubject{} },
			}),
			conform.Excuse(
				store.StoreChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Put.RefusesACancelledPut(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStoreSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		store.RunStoreSuiteFromHandle(
			t,
			func(*store.StoreHandle) *storeSubject {
				built = &storeSubject{}
				return built
			},
			conform.With(store.StoreCheck{
				ID:        "Store/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *store.StoreHandle) {
					got := h.Subject()
					assert.Equal(h, got, store.Store(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*store.StoreHandle) store.Store { return &storeSubject{} },
			}),
			conform.Excuse(
				store.StoreChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				store.StoreChecks.Put.RefusesACancelledPut(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("StoreChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(store.StoreChecks.Put.RefusesACancelledPut()),
			"Store.Put/refuses-a-cancelled-put",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// storeIdle returns the planted defect of a ceiling of Store with a spend
// of nothing, a subject that answers each call at once.
func storeIdle(h *store.StoreHandle) store.Store {
	return store.StoreSuite().Calls[0].Over(h, func() {})
}

// TestStoreSuiteBudgets applies the allocation ceilings of the suite of
// Store to storeIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestStoreSuiteBudgets(t *testing.T) {
	store.RunStoreBudgets(
		t,
		storeIdle,
		conform.Ceilings(store.StoreChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(store.StoreChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			store.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.RefusesACancelledPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteGet registers the fuzz target of Store.Get on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuiteGet(f *testing.F) {
	store.FuzzStoreSuite(
		f,
		store.StoreChecks.Get.Callable,
		storeIdle,
		conform.Excuse(
			store.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.RefusesACancelledPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuitePut registers the fuzz target of Store.Put on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuitePut(f *testing.F) {
	store.FuzzStoreSuite(
		f,
		store.StoreChecks.Put.Callable,
		storeIdle,
		conform.Excuse(
			store.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			store.StoreChecks.Put.RefusesACancelledPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStoreSuite measures each callable of the suite of Store on
// storeIdle.
func BenchmarkStoreSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		store.BenchStoreSuite(b, store.StoreChecks.Get.Callable, storeIdle)
	})
	b.Run("Put", func(b *testing.B) {
		store.BenchStoreSuite(b, store.StoreChecks.Put.Callable, storeIdle)
	})
}
