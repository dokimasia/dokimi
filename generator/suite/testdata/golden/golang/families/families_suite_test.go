package families_test

import (
	"testing"

	"example.com/shop/families"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// storeSubjectEmbedded is families.Store under the name of the field that
// storeSubject embeds.
type storeSubjectEmbedded = families.Store

// storeSubject is a value of Store that only TestStoreSuite builds.
type storeSubject struct{ storeSubjectEmbedded }

// TestStoreSuite checks the suite of Store and the two forms of its
// application.
func TestStoreSuite(t *testing.T) {
	t.Parallel()

	t.Run("StoreSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.StoreSuite())
	})

	t.Run("RunStoreSuite", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		families.RunStoreSuite(
			t,
			func() *storeSubject {
				built = &storeSubject{}
				return built
			},
			conform.With(families.StoreCheck{
				ID:        "Store/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.StoreHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Store(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.StoreHandle) families.Store { return &storeSubject{} },
			}),
			conform.Excuse(
				families.StoreChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.MonotonicReads(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Load.ZeroOnMiss(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Load.Cacheable(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Split.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Level.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Level.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Level.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Raise.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Raise.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Inspect.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.GetWithMeta.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.GetWithMeta.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.GetWithMeta.ZeroesOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStoreSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *storeSubject
		families.RunStoreSuiteFromHandle(
			t,
			func(*families.StoreHandle) *storeSubject {
				built = &storeSubject{}
				return built
			},
			conform.With(families.StoreCheck{
				ID:        "Store/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.StoreHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Store(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.StoreHandle) families.Store { return &storeSubject{} },
			}),
			conform.Excuse(
				families.StoreChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Get.MonotonicReads(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Load.ZeroOnMiss(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Load.Cacheable(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Split.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Level.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Level.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Level.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Raise.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Raise.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.Inspect.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.GetWithMeta.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.GetWithMeta.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StoreChecks.GetWithMeta.ZeroesOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// storeIdle returns the planted defect of a ceiling of Store with a spend
// of nothing, a subject that answers each call at once.
func storeIdle(h *families.StoreHandle) families.Store {
	return families.StoreSuite().Calls[0].Over(h, func() {})
}

// TestStoreSuiteBudgets applies the allocation ceilings of the suite of
// Store to storeIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestStoreSuiteBudgets(t *testing.T) {
	families.RunStoreBudgets(
		t,
		storeIdle,
		conform.Ceilings(families.StoreChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(families.StoreChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			families.StoreChecks.Load.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.StoreChecks.Split.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.StoreChecks.Find.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.StoreChecks.Level.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.StoreChecks.Raise.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.StoreChecks.Inspect.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.StoreChecks.GetWithMeta.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuitePut registers the fuzz target of Store.Put on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuitePut(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Put.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteGet registers the fuzz target of Store.Get on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuiteGet(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Get.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteLoad registers the fuzz target of Store.Load on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuiteLoad(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Load.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteSplit registers the fuzz target of Store.Split on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteSplit(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Split.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteFind registers the fuzz target of Store.Find on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuiteFind(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Find.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteLevel registers the fuzz target of Store.Level on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteLevel(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Level.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteRaise registers the fuzz target of Store.Raise on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteRaise(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Raise.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteInspect registers the fuzz target of Store.Inspect on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteInspect(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.Inspect.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteGetWithMeta registers the fuzz target of Store.GetWithMeta
// on storeIdle, with every check of the suite excused.
func FuzzStoreSuiteGetWithMeta(f *testing.F) {
	families.FuzzStoreSuite(
		f,
		families.StoreChecks.GetWithMeta.Callable,
		storeIdle,
		conform.Excuse(
			families.StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStoreSuite measures each callable of the suite of Store on
// storeIdle.
func BenchmarkStoreSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Put.Callable, storeIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Get.Callable, storeIdle)
	})
	b.Run("Load", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Load.Callable, storeIdle)
	})
	b.Run("Split", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Split.Callable, storeIdle)
	})
	b.Run("Find", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Find.Callable, storeIdle)
	})
	b.Run("Level", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Level.Callable, storeIdle)
	})
	b.Run("Raise", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Raise.Callable, storeIdle)
	})
	b.Run("Inspect", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.Inspect.Callable, storeIdle)
	})
	b.Run("GetWithMeta", func(b *testing.B) {
		families.BenchStoreSuite(b, families.StoreChecks.GetWithMeta.Callable, storeIdle)
	})
}

// keyedSubjectEmbedded is families.Keyed under the name of the field that
// keyedSubject embeds.
type keyedSubjectEmbedded = families.Keyed

// keyedSubject is a value of Keyed that only TestKeyedSuite builds.
type keyedSubject struct{ keyedSubjectEmbedded }

// TestKeyedSuite checks the suite of Keyed and the two forms of its
// application.
func TestKeyedSuite(t *testing.T) {
	t.Parallel()

	t.Run("KeyedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.KeyedSuite())
	})

	t.Run("RunKeyedSuite", func(t *testing.T) {
		t.Parallel()
		var built *keyedSubject
		families.RunKeyedSuite(
			t,
			func() *keyedSubject {
				built = &keyedSubject{}
				return built
			},
			conform.With(families.KeyedCheck{
				ID:        "Keyed/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.KeyedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Keyed(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.KeyedHandle) families.Keyed { return &keyedSubject{} },
			}),
			conform.Excuse(
				families.KeyedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Delete.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Delete.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunKeyedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *keyedSubject
		families.RunKeyedSuiteFromHandle(
			t,
			func(*families.KeyedHandle) *keyedSubject {
				built = &keyedSubject{}
				return built
			},
			conform.With(families.KeyedCheck{
				ID:        "Keyed/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.KeyedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Keyed(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.KeyedHandle) families.Keyed { return &keyedSubject{} },
			}),
			conform.Excuse(
				families.KeyedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Delete.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedChecks.Delete.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// keyedIdle returns the planted defect of a ceiling of Keyed with a spend
// of nothing, a subject that answers each call at once.
func keyedIdle(h *families.KeyedHandle) families.Keyed {
	return families.KeyedSuite().Calls[0].Over(h, func() {})
}

// TestKeyedSuiteBudgets applies the allocation ceilings of the suite of
// Keyed to keyedIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestKeyedSuiteBudgets(t *testing.T) {
	families.RunKeyedBudgets(
		t,
		keyedIdle,
		conform.Ceilings(families.KeyedChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(families.KeyedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			families.KeyedChecks.Delete.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedSuitePut registers the fuzz target of Keyed.Put on keyedIdle,
// with every check of the suite excused.
func FuzzKeyedSuitePut(f *testing.F) {
	families.FuzzKeyedSuite(
		f,
		families.KeyedChecks.Put.Callable,
		keyedIdle,
		conform.Excuse(
			families.KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedSuiteGet registers the fuzz target of Keyed.Get on keyedIdle,
// with every check of the suite excused.
func FuzzKeyedSuiteGet(f *testing.F) {
	families.FuzzKeyedSuite(
		f,
		families.KeyedChecks.Get.Callable,
		keyedIdle,
		conform.Excuse(
			families.KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedSuiteDelete registers the fuzz target of Keyed.Delete on
// keyedIdle, with every check of the suite excused.
func FuzzKeyedSuiteDelete(f *testing.F) {
	families.FuzzKeyedSuite(
		f,
		families.KeyedChecks.Delete.Callable,
		keyedIdle,
		conform.Excuse(
			families.KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkKeyedSuite measures each callable of the suite of Keyed on
// keyedIdle.
func BenchmarkKeyedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchKeyedSuite(b, families.KeyedChecks.Put.Callable, keyedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchKeyedSuite(b, families.KeyedChecks.Get.Callable, keyedIdle)
	})
	b.Run("Delete", func(b *testing.B) {
		families.BenchKeyedSuite(b, families.KeyedChecks.Delete.Callable, keyedIdle)
	})
}

// stampedSubjectEmbedded is families.Stamped under the name of the field
// that stampedSubject embeds.
type stampedSubjectEmbedded = families.Stamped

// stampedSubject is a value of Stamped that only TestStampedSuite builds.
type stampedSubject struct{ stampedSubjectEmbedded }

// TestStampedSuite checks the suite of Stamped and the two forms of its
// application.
func TestStampedSuite(t *testing.T) {
	t.Parallel()

	t.Run("StampedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.StampedSuite())
	})

	t.Run("RunStampedSuite", func(t *testing.T) {
		t.Parallel()
		var built *stampedSubject
		families.RunStampedSuite(
			t,
			func() *stampedSubject {
				built = &stampedSubject{}
				return built
			},
			conform.With(families.StampedCheck{
				ID:        "Stamped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.StampedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Stamped(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.StampedHandle) families.Stamped { return &stampedSubject{} },
			}),
			conform.Excuse(
				families.StampedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Store.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.WritesFollowReads(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStampedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *stampedSubject
		families.RunStampedSuiteFromHandle(
			t,
			func(*families.StampedHandle) *stampedSubject {
				built = &stampedSubject{}
				return built
			},
			conform.With(families.StampedCheck{
				ID:        "Stamped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.StampedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Stamped(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.StampedHandle) families.Stamped { return &stampedSubject{} },
			}),
			conform.Excuse(
				families.StampedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Store.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.StampedChecks.Get.WritesFollowReads(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// stampedIdle returns the planted defect of a ceiling of Stamped with a
// spend of nothing, a subject that answers each call at once.
func stampedIdle(h *families.StampedHandle) families.Stamped {
	return families.StampedSuite().Calls[0].Over(h, func() {})
}

// TestStampedSuiteBudgets applies the allocation ceilings of the suite of
// Stamped to stampedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestStampedSuiteBudgets(t *testing.T) {
	families.RunStampedBudgets(
		t,
		stampedIdle,
		conform.Ceilings(
			families.StampedChecks.Store.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.StampedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Store.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.WritesFollowReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStampedSuiteStore registers the fuzz target of Stamped.Store on
// stampedIdle, with every check of the suite excused.
func FuzzStampedSuiteStore(f *testing.F) {
	families.FuzzStampedSuite(
		f,
		families.StampedChecks.Store.Callable,
		stampedIdle,
		conform.Excuse(
			families.StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Store.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.WritesFollowReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStampedSuiteGet registers the fuzz target of Stamped.Get on
// stampedIdle, with every check of the suite excused.
func FuzzStampedSuiteGet(f *testing.F) {
	families.FuzzStampedSuite(
		f,
		families.StampedChecks.Get.Callable,
		stampedIdle,
		conform.Excuse(
			families.StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Store.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.StampedChecks.Get.WritesFollowReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStampedSuite measures each callable of the suite of Stamped on
// stampedIdle.
func BenchmarkStampedSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		families.BenchStampedSuite(b, families.StampedChecks.Store.Callable, stampedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchStampedSuite(b, families.StampedChecks.Get.Callable, stampedIdle)
	})
}

// decodedSubjectEmbedded is families.Decoded under the name of the field
// that decodedSubject embeds.
type decodedSubjectEmbedded = families.Decoded

// decodedSubject is a value of Decoded that only TestDecodedSuite builds.
type decodedSubject struct{ decodedSubjectEmbedded }

// TestDecodedSuite checks the suite of Decoded and the two forms of its
// application.
func TestDecodedSuite(t *testing.T) {
	t.Parallel()

	t.Run("DecodedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.DecodedSuite())
	})

	t.Run("RunDecodedSuite", func(t *testing.T) {
		t.Parallel()
		var built *decodedSubject
		families.RunDecodedSuite(
			t,
			func() *decodedSubject {
				built = &decodedSubject{}
				return built
			},
			conform.With(families.DecodedCheck{
				ID:        "Decoded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.DecodedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Decoded(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.DecodedHandle) families.Decoded { return &decodedSubject{} },
			}),
			conform.Excuse(
				families.DecodedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunDecodedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *decodedSubject
		families.RunDecodedSuiteFromHandle(
			t,
			func(*families.DecodedHandle) *decodedSubject {
				built = &decodedSubject{}
				return built
			},
			conform.With(families.DecodedCheck{
				ID:        "Decoded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.DecodedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Decoded(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.DecodedHandle) families.Decoded { return &decodedSubject{} },
			}),
			conform.Excuse(
				families.DecodedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.DecodedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// decodedIdle returns the planted defect of a ceiling of Decoded with a
// spend of nothing, a subject that answers each call at once.
func decodedIdle(h *families.DecodedHandle) families.Decoded {
	return families.DecodedSuite().Calls[0].Over(h, func() {})
}

// TestDecodedSuiteBudgets applies the allocation ceilings of the suite of
// Decoded to decodedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestDecodedSuiteBudgets(t *testing.T) {
	families.RunDecodedBudgets(
		t,
		decodedIdle,
		conform.Ceilings(
			families.DecodedChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.DecodedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.DecodedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDecodedSuitePut registers the fuzz target of Decoded.Put on
// decodedIdle, with every check of the suite excused.
func FuzzDecodedSuitePut(f *testing.F) {
	families.FuzzDecodedSuite(
		f,
		families.DecodedChecks.Put.Callable,
		decodedIdle,
		conform.Excuse(
			families.DecodedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDecodedSuiteGet registers the fuzz target of Decoded.Get on
// decodedIdle, with every check of the suite excused.
func FuzzDecodedSuiteGet(f *testing.F) {
	families.FuzzDecodedSuite(
		f,
		families.DecodedChecks.Get.Callable,
		decodedIdle,
		conform.Excuse(
			families.DecodedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.DecodedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkDecodedSuite measures each callable of the suite of Decoded on
// decodedIdle.
func BenchmarkDecodedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchDecodedSuite(b, families.DecodedChecks.Put.Callable, decodedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchDecodedSuite(b, families.DecodedChecks.Get.Callable, decodedIdle)
	})
}

// expiringSubjectEmbedded is families.Expiring under the name of the field
// that expiringSubject embeds.
type expiringSubjectEmbedded = families.Expiring

// expiringSubject is a value of Expiring that only TestExpiringSuite
// builds.
type expiringSubject struct{ expiringSubjectEmbedded }

// TestExpiringSuite checks the suite of Expiring and the two forms of its
// application.
func TestExpiringSuite(t *testing.T) {
	t.Parallel()

	t.Run("ExpiringSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.ExpiringSuite())
	})

	t.Run("RunExpiringSuite", func(t *testing.T) {
		t.Parallel()
		var built *expiringSubject
		families.RunExpiringSuite(
			t,
			func() *expiringSubject {
				built = &expiringSubject{}
				return built
			},
			conform.With(families.ExpiringCheck{
				ID:        "Expiring/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ExpiringHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Expiring(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.ExpiringHandle) families.Expiring { return &expiringSubject{} },
			}),
			conform.Excuse(
				families.ExpiringChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.Expires(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunExpiringSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *expiringSubject
		families.RunExpiringSuiteFromHandle(
			t,
			func(*families.ExpiringHandle) *expiringSubject {
				built = &expiringSubject{}
				return built
			},
			conform.With(families.ExpiringCheck{
				ID:        "Expiring/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ExpiringHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Expiring(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.ExpiringHandle) families.Expiring { return &expiringSubject{} },
			}),
			conform.Excuse(
				families.ExpiringChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ExpiringChecks.Read.Expires(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// expiringIdle returns the planted defect of a ceiling of Expiring with a
// spend of nothing, a subject that answers each call at once.
func expiringIdle(h *families.ExpiringHandle) families.Expiring {
	return families.ExpiringSuite().Calls[0].Over(h, func() {})
}

// TestExpiringSuiteBudgets applies the allocation ceilings of the suite of
// Expiring to expiringIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestExpiringSuiteBudgets(t *testing.T) {
	families.RunExpiringBudgets(
		t,
		expiringIdle,
		conform.Ceilings(
			families.ExpiringChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.ExpiringChecks.Read.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.ExpiringChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Expires(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzExpiringSuitePut registers the fuzz target of Expiring.Put on
// expiringIdle, with every check of the suite excused.
func FuzzExpiringSuitePut(f *testing.F) {
	families.FuzzExpiringSuite(
		f,
		families.ExpiringChecks.Put.Callable,
		expiringIdle,
		conform.Excuse(
			families.ExpiringChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Expires(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzExpiringSuiteRead registers the fuzz target of Expiring.Read on
// expiringIdle, with every check of the suite excused.
func FuzzExpiringSuiteRead(f *testing.F) {
	families.FuzzExpiringSuite(
		f,
		families.ExpiringChecks.Read.Callable,
		expiringIdle,
		conform.Excuse(
			families.ExpiringChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ExpiringChecks.Read.Expires(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkExpiringSuite measures each callable of the suite of Expiring
// on expiringIdle.
func BenchmarkExpiringSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchExpiringSuite(b, families.ExpiringChecks.Put.Callable, expiringIdle)
	})
	b.Run("Read", func(b *testing.B) {
		families.BenchExpiringSuite(b, families.ExpiringChecks.Read.Callable, expiringIdle)
	})
}

// cachedSubjectEmbedded is families.Cached under the name of the field
// that cachedSubject embeds.
type cachedSubjectEmbedded = families.Cached

// cachedSubject is a value of Cached that only TestCachedSuite builds.
type cachedSubject struct{ cachedSubjectEmbedded }

// TestCachedSuite checks the suite of Cached and the two forms of its
// application.
func TestCachedSuite(t *testing.T) {
	t.Parallel()

	t.Run("CachedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.CachedSuite())
	})

	t.Run("RunCachedSuite", func(t *testing.T) {
		t.Parallel()
		var built *cachedSubject
		families.RunCachedSuite(
			t,
			func() *cachedSubject {
				built = &cachedSubject{}
				return built
			},
			conform.With(families.CachedCheck{
				ID:        "Cached/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.CachedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Cached(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.CachedHandle) families.Cached { return &cachedSubject{} },
			}),
			conform.Excuse(
				families.CachedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Get.Cache(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCachedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *cachedSubject
		families.RunCachedSuiteFromHandle(
			t,
			func(*families.CachedHandle) *cachedSubject {
				built = &cachedSubject{}
				return built
			},
			conform.With(families.CachedCheck{
				ID:        "Cached/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.CachedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Cached(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.CachedHandle) families.Cached { return &cachedSubject{} },
			}),
			conform.Excuse(
				families.CachedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Get.Cache(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CachedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// cachedIdle returns the planted defect of a ceiling of Cached with a
// spend of nothing, a subject that answers each call at once.
func cachedIdle(h *families.CachedHandle) families.Cached {
	return families.CachedSuite().Calls[0].Over(h, func() {})
}

// TestCachedSuiteBudgets applies the allocation ceilings of the suite of
// Cached to cachedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCachedSuiteBudgets(t *testing.T) {
	families.RunCachedBudgets(
		t,
		cachedIdle,
		conform.Ceilings(
			families.CachedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.CachedChecks.Fetch.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.CachedChecks.Store.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteGet registers the fuzz target of Cached.Get on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteGet(f *testing.F) {
	families.FuzzCachedSuite(
		f,
		families.CachedChecks.Get.Callable,
		cachedIdle,
		conform.Excuse(
			families.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteFetch registers the fuzz target of Cached.Fetch on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteFetch(f *testing.F) {
	families.FuzzCachedSuite(
		f,
		families.CachedChecks.Fetch.Callable,
		cachedIdle,
		conform.Excuse(
			families.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteStore registers the fuzz target of Cached.Store on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteStore(f *testing.F) {
	families.FuzzCachedSuite(
		f,
		families.CachedChecks.Store.Callable,
		cachedIdle,
		conform.Excuse(
			families.CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCachedSuite measures each callable of the suite of Cached on
// cachedIdle.
func BenchmarkCachedSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		families.BenchCachedSuite(b, families.CachedChecks.Get.Callable, cachedIdle)
	})
	b.Run("Fetch", func(b *testing.B) {
		families.BenchCachedSuite(b, families.CachedChecks.Fetch.Callable, cachedIdle)
	})
	b.Run("Store", func(b *testing.B) {
		families.BenchCachedSuite(b, families.CachedChecks.Store.Callable, cachedIdle)
	})
}

// taggedSubjectEmbedded is families.Tagged under the name of the field
// that taggedSubject embeds.
type taggedSubjectEmbedded = families.Tagged

// taggedSubject is a value of Tagged that only TestTaggedSuite builds.
type taggedSubject struct{ taggedSubjectEmbedded }

// TestTaggedSuite checks the suite of Tagged and the two forms of its
// application.
func TestTaggedSuite(t *testing.T) {
	t.Parallel()

	t.Run("TaggedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.TaggedSuite())
	})

	t.Run("RunTaggedSuite", func(t *testing.T) {
		t.Parallel()
		var built *taggedSubject
		families.RunTaggedSuite(
			t,
			func() *taggedSubject {
				built = &taggedSubject{}
				return built
			},
			conform.With(families.TaggedCheck{
				ID:        "Tagged/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.TaggedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Tagged(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.TaggedHandle) families.Tagged { return &taggedSubject{} },
			}),
			conform.Excuse(
				families.TaggedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTaggedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *taggedSubject
		families.RunTaggedSuiteFromHandle(
			t,
			func(*families.TaggedHandle) *taggedSubject {
				built = &taggedSubject{}
				return built
			},
			conform.With(families.TaggedCheck{
				ID:        "Tagged/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.TaggedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Tagged(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.TaggedHandle) families.Tagged { return &taggedSubject{} },
			}),
			conform.Excuse(
				families.TaggedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TaggedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// taggedIdle returns the planted defect of a ceiling of Tagged with a
// spend of nothing, a subject that answers each call at once.
func taggedIdle(h *families.TaggedHandle) families.Tagged {
	return families.TaggedSuite().Calls[0].Over(h, func() {})
}

// TestTaggedSuiteBudgets applies the allocation ceilings of the suite of
// Tagged to taggedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestTaggedSuiteBudgets(t *testing.T) {
	families.RunTaggedBudgets(
		t,
		taggedIdle,
		conform.Ceilings(
			families.TaggedChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.TaggedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.TaggedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTaggedSuitePut registers the fuzz target of Tagged.Put on
// taggedIdle, with every check of the suite excused.
func FuzzTaggedSuitePut(f *testing.F) {
	families.FuzzTaggedSuite(
		f,
		families.TaggedChecks.Put.Callable,
		taggedIdle,
		conform.Excuse(
			families.TaggedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTaggedSuiteGet registers the fuzz target of Tagged.Get on
// taggedIdle, with every check of the suite excused.
func FuzzTaggedSuiteGet(f *testing.F) {
	families.FuzzTaggedSuite(
		f,
		families.TaggedChecks.Get.Callable,
		taggedIdle,
		conform.Excuse(
			families.TaggedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TaggedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkTaggedSuite measures each callable of the suite of Tagged on
// taggedIdle.
func BenchmarkTaggedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchTaggedSuite(b, families.TaggedChecks.Put.Callable, taggedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchTaggedSuite(b, families.TaggedChecks.Get.Callable, taggedIdle)
	})
}

// snapshotsSubjectEmbedded is families.Snapshots under the name of the
// field that snapshotsSubject embeds.
type snapshotsSubjectEmbedded = families.Snapshots

// snapshotsSubject is a value of Snapshots that only TestSnapshotsSuite
// builds.
type snapshotsSubject struct{ snapshotsSubjectEmbedded }

// TestSnapshotsSuite checks the suite of Snapshots and the two forms of
// its application.
func TestSnapshotsSuite(t *testing.T) {
	t.Parallel()

	t.Run("SnapshotsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.SnapshotsSuite())
	})

	t.Run("RunSnapshotsSuite", func(t *testing.T) {
		t.Parallel()
		var built *snapshotsSubject
		families.RunSnapshotsSuite(
			t,
			func() *snapshotsSubject {
				built = &snapshotsSubject{}
				return built
			},
			conform.With(families.SnapshotsCheck{
				ID:        "Snapshots/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.SnapshotsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Snapshots(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.SnapshotsHandle) families.Snapshots { return &snapshotsSubject{} },
			}),
			conform.Excuse(
				families.SnapshotsChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.PointInTime(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSnapshotsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *snapshotsSubject
		families.RunSnapshotsSuiteFromHandle(
			t,
			func(*families.SnapshotsHandle) *snapshotsSubject {
				built = &snapshotsSubject{}
				return built
			},
			conform.With(families.SnapshotsCheck{
				ID:        "Snapshots/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.SnapshotsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Snapshots(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.SnapshotsHandle) families.Snapshots { return &snapshotsSubject{} },
			}),
			conform.Excuse(
				families.SnapshotsChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SnapshotsChecks.Get.PointInTime(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// snapshotsIdle returns the planted defect of a ceiling of Snapshots with
// a spend of nothing, a subject that answers each call at once.
func snapshotsIdle(h *families.SnapshotsHandle) families.Snapshots {
	return families.SnapshotsSuite().Calls[0].Over(h, func() {})
}

// TestSnapshotsSuiteBudgets applies the allocation ceilings of the suite
// of Snapshots to snapshotsIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestSnapshotsSuiteBudgets(t *testing.T) {
	families.RunSnapshotsBudgets(
		t,
		snapshotsIdle,
		conform.Ceilings(
			families.SnapshotsChecks.Store.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.SnapshotsChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.SnapshotsChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.PointInTime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSnapshotsSuiteStore registers the fuzz target of Snapshots.Store on
// snapshotsIdle, with every check of the suite excused.
func FuzzSnapshotsSuiteStore(f *testing.F) {
	families.FuzzSnapshotsSuite(
		f,
		families.SnapshotsChecks.Store.Callable,
		snapshotsIdle,
		conform.Excuse(
			families.SnapshotsChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.PointInTime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSnapshotsSuiteGet registers the fuzz target of Snapshots.Get on
// snapshotsIdle, with every check of the suite excused.
func FuzzSnapshotsSuiteGet(f *testing.F) {
	families.FuzzSnapshotsSuite(
		f,
		families.SnapshotsChecks.Get.Callable,
		snapshotsIdle,
		conform.Excuse(
			families.SnapshotsChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SnapshotsChecks.Get.PointInTime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSnapshotsSuite measures each callable of the suite of Snapshots
// on snapshotsIdle.
func BenchmarkSnapshotsSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		families.BenchSnapshotsSuite(b, families.SnapshotsChecks.Store.Callable, snapshotsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchSnapshotsSuite(b, families.SnapshotsChecks.Get.Callable, snapshotsIdle)
	})
}

// bagSubjectEmbedded is families.Bag under the name of the field that
// bagSubject embeds.
type bagSubjectEmbedded = families.Bag

// bagSubject is a value of Bag that only TestBagSuite builds.
type bagSubject struct{ bagSubjectEmbedded }

// TestBagSuite checks the suite of Bag and the two forms of its
// application.
func TestBagSuite(t *testing.T) {
	t.Parallel()

	t.Run("BagSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.BagSuite())
	})

	t.Run("RunBagSuite", func(t *testing.T) {
		t.Parallel()
		var built *bagSubject
		families.RunBagSuite(
			t,
			func() *bagSubject {
				built = &bagSubject{}
				return built
			},
			conform.With(families.BagCheck{
				ID:        "Bag/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.BagHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Bag(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.BagHandle) families.Bag { return &bagSubject{} },
			}),
			conform.Excuse(
				families.BagChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.StableOrder(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBagSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *bagSubject
		families.RunBagSuiteFromHandle(
			t,
			func(*families.BagHandle) *bagSubject {
				built = &bagSubject{}
				return built
			},
			conform.With(families.BagCheck{
				ID:        "Bag/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.BagHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Bag(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.BagHandle) families.Bag { return &bagSubject{} },
			}),
			conform.Excuse(
				families.BagChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BagChecks.Items.StableOrder(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// bagIdle returns the planted defect of a ceiling of Bag with a spend of
// nothing, a subject that answers each call at once.
func bagIdle(h *families.BagHandle) families.Bag {
	return families.BagSuite().Calls[0].Over(h, func() {})
}

// TestBagSuiteBudgets applies the allocation ceilings of the suite of Bag
// to bagIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestBagSuiteBudgets(t *testing.T) {
	families.RunBagBudgets(
		t,
		bagIdle,
		conform.Ceilings(families.BagChecks.Add.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(families.BagChecks.Items.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			families.BagChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.StableOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBagSuiteAdd registers the fuzz target of Bag.Add on bagIdle, with
// every check of the suite excused.
func FuzzBagSuiteAdd(f *testing.F) {
	families.FuzzBagSuite(
		f,
		families.BagChecks.Add.Callable,
		bagIdle,
		conform.Excuse(
			families.BagChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.StableOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBagSuiteItems registers the fuzz target of Bag.Items on bagIdle,
// with every check of the suite excused.
func FuzzBagSuiteItems(f *testing.F) {
	families.FuzzBagSuite(
		f,
		families.BagChecks.Items.Callable,
		bagIdle,
		conform.Excuse(
			families.BagChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BagChecks.Items.StableOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBagSuite measures each callable of the suite of Bag on bagIdle.
func BenchmarkBagSuite(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		families.BenchBagSuite(b, families.BagChecks.Add.Callable, bagIdle)
	})
	b.Run("Items", func(b *testing.B) {
		families.BenchBagSuite(b, families.BagChecks.Items.Callable, bagIdle)
	})
}

// updaterSubjectEmbedded is families.Updater under the name of the field
// that updaterSubject embeds.
type updaterSubjectEmbedded = families.Updater

// updaterSubject is a value of Updater that only TestUpdaterSuite builds.
type updaterSubject struct{ updaterSubjectEmbedded }

// TestUpdaterSuite checks the suite of Updater and the two forms of its
// application.
func TestUpdaterSuite(t *testing.T) {
	t.Parallel()

	t.Run("UpdaterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.UpdaterSuite())
	})

	t.Run("RunUpdaterSuite", func(t *testing.T) {
		t.Parallel()
		var built *updaterSubject
		families.RunUpdaterSuite(
			t,
			func() *updaterSubject {
				built = &updaterSubject{}
				return built
			},
			conform.With(families.UpdaterCheck{
				ID:        "Updater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.UpdaterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Updater(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.UpdaterHandle) families.Updater { return &updaterSubject{} },
			}),
			conform.Excuse(
				families.UpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunUpdaterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *updaterSubject
		families.RunUpdaterSuiteFromHandle(
			t,
			func(*families.UpdaterHandle) *updaterSubject {
				built = &updaterSubject{}
				return built
			},
			conform.With(families.UpdaterCheck{
				ID:        "Updater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.UpdaterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Updater(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.UpdaterHandle) families.Updater { return &updaterSubject{} },
			}),
			conform.Excuse(
				families.UpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.UpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// updaterIdle returns the planted defect of a ceiling of Updater with a
// spend of nothing, a subject that answers each call at once.
func updaterIdle(h *families.UpdaterHandle) families.Updater {
	return families.UpdaterSuite().Calls[0].Over(h, func() {})
}

// TestUpdaterSuiteBudgets applies the allocation ceilings of the suite of
// Updater to updaterIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestUpdaterSuiteBudgets(t *testing.T) {
	families.RunUpdaterBudgets(
		t,
		updaterIdle,
		conform.Ceilings(
			families.UpdaterChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.UpdaterChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.UpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzUpdaterSuitePut registers the fuzz target of Updater.Put on
// updaterIdle, with every check of the suite excused.
func FuzzUpdaterSuitePut(f *testing.F) {
	families.FuzzUpdaterSuite(
		f,
		families.UpdaterChecks.Put.Callable,
		updaterIdle,
		conform.Excuse(
			families.UpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzUpdaterSuiteGet registers the fuzz target of Updater.Get on
// updaterIdle, with every check of the suite excused.
func FuzzUpdaterSuiteGet(f *testing.F) {
	families.FuzzUpdaterSuite(
		f,
		families.UpdaterChecks.Get.Callable,
		updaterIdle,
		conform.Excuse(
			families.UpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.UpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkUpdaterSuite measures each callable of the suite of Updater on
// updaterIdle.
func BenchmarkUpdaterSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchUpdaterSuite(b, families.UpdaterChecks.Put.Callable, updaterIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchUpdaterSuite(b, families.UpdaterChecks.Get.Callable, updaterIdle)
	})
}

// keyedUpdaterSubjectEmbedded is families.KeyedUpdater under the name of
// the field that keyedUpdaterSubject embeds.
type keyedUpdaterSubjectEmbedded = families.KeyedUpdater

// keyedUpdaterSubject is a value of KeyedUpdater that only
// TestKeyedUpdaterSuite builds.
type keyedUpdaterSubject struct{ keyedUpdaterSubjectEmbedded }

// TestKeyedUpdaterSuite checks the suite of KeyedUpdater and the two forms
// of its application.
func TestKeyedUpdaterSuite(t *testing.T) {
	t.Parallel()

	t.Run("KeyedUpdaterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.KeyedUpdaterSuite())
	})

	t.Run("RunKeyedUpdaterSuite", func(t *testing.T) {
		t.Parallel()
		var built *keyedUpdaterSubject
		families.RunKeyedUpdaterSuite(
			t,
			func() *keyedUpdaterSubject {
				built = &keyedUpdaterSubject{}
				return built
			},
			conform.With(families.KeyedUpdaterCheck{
				ID:        "KeyedUpdater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.KeyedUpdaterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.KeyedUpdater(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.KeyedUpdaterHandle) families.KeyedUpdater { return &keyedUpdaterSubject{} },
			}),
			conform.Excuse(
				families.KeyedUpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunKeyedUpdaterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *keyedUpdaterSubject
		families.RunKeyedUpdaterSuiteFromHandle(
			t,
			func(*families.KeyedUpdaterHandle) *keyedUpdaterSubject {
				built = &keyedUpdaterSubject{}
				return built
			},
			conform.With(families.KeyedUpdaterCheck{
				ID:        "KeyedUpdater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.KeyedUpdaterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.KeyedUpdater(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.KeyedUpdaterHandle) families.KeyedUpdater { return &keyedUpdaterSubject{} },
			}),
			conform.Excuse(
				families.KeyedUpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.KeyedUpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// keyedUpdaterIdle returns the planted defect of a ceiling of KeyedUpdater
// with a spend of nothing, a subject that answers each call at once.
func keyedUpdaterIdle(h *families.KeyedUpdaterHandle) families.KeyedUpdater {
	return families.KeyedUpdaterSuite().Calls[0].Over(h, func() {})
}

// TestKeyedUpdaterSuiteBudgets applies the allocation ceilings of the
// suite of KeyedUpdater to keyedUpdaterIdle with every check of the suite
// excused, and a ceiling of 0 to each callable without a budget whose
// calls build no input of their own.
func TestKeyedUpdaterSuiteBudgets(t *testing.T) {
	families.RunKeyedUpdaterBudgets(
		t,
		keyedUpdaterIdle,
		conform.Ceilings(
			families.KeyedUpdaterChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.KeyedUpdaterChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedUpdaterSuitePut registers the fuzz target of KeyedUpdater.Put
// on keyedUpdaterIdle, with every check of the suite excused.
func FuzzKeyedUpdaterSuitePut(f *testing.F) {
	families.FuzzKeyedUpdaterSuite(
		f,
		families.KeyedUpdaterChecks.Put.Callable,
		keyedUpdaterIdle,
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedUpdaterSuiteGet registers the fuzz target of KeyedUpdater.Get
// on keyedUpdaterIdle, with every check of the suite excused.
func FuzzKeyedUpdaterSuiteGet(f *testing.F) {
	families.FuzzKeyedUpdaterSuite(
		f,
		families.KeyedUpdaterChecks.Get.Callable,
		keyedUpdaterIdle,
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.KeyedUpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkKeyedUpdaterSuite measures each callable of the suite of
// KeyedUpdater on keyedUpdaterIdle.
func BenchmarkKeyedUpdaterSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchKeyedUpdaterSuite(
			b,
			families.KeyedUpdaterChecks.Put.Callable,
			keyedUpdaterIdle,
		)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchKeyedUpdaterSuite(
			b,
			families.KeyedUpdaterChecks.Get.Callable,
			keyedUpdaterIdle,
		)
	})
}

// codecSubjectEmbedded is families.Codec under the name of the field that
// codecSubject embeds.
type codecSubjectEmbedded = families.Codec

// codecSubject is a value of Codec that only TestCodecSuite builds.
type codecSubject struct{ codecSubjectEmbedded }

// TestCodecSuite checks the suite of Codec and the two forms of its
// application.
func TestCodecSuite(t *testing.T) {
	t.Parallel()

	t.Run("CodecSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.CodecSuite())
	})

	t.Run("RunCodecSuite", func(t *testing.T) {
		t.Parallel()
		var built *codecSubject
		families.RunCodecSuite(
			t,
			func() *codecSubject {
				built = &codecSubject{}
				return built
			},
			conform.With(families.CodecCheck{
				ID:        "Codec/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.CodecHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Codec(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.CodecHandle) families.Codec { return &codecSubject{} },
			}),
			conform.Excuse(
				families.CodecChecks.Encode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Encode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Encode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Encode.Codec(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Decode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Decode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Decode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCodecSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *codecSubject
		families.RunCodecSuiteFromHandle(
			t,
			func(*families.CodecHandle) *codecSubject {
				built = &codecSubject{}
				return built
			},
			conform.With(families.CodecCheck{
				ID:        "Codec/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.CodecHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Codec(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.CodecHandle) families.Codec { return &codecSubject{} },
			}),
			conform.Excuse(
				families.CodecChecks.Encode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Encode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Encode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Encode.Codec(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Decode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Decode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.CodecChecks.Decode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// codecIdle returns the planted defect of a ceiling of Codec with a spend
// of nothing, a subject that answers each call at once.
func codecIdle(h *families.CodecHandle) families.Codec {
	return families.CodecSuite().Calls[0].Over(h, func() {})
}

// TestCodecSuiteBudgets applies the allocation ceilings of the suite of
// Codec to codecIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestCodecSuiteBudgets(t *testing.T) {
	families.RunCodecBudgets(
		t,
		codecIdle,
		conform.Ceilings(
			families.CodecChecks.Encode.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.CodecChecks.Decode.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.CodecChecks.Encode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.Codec(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCodecSuiteEncode registers the fuzz target of Codec.Encode on
// codecIdle, with every check of the suite excused.
func FuzzCodecSuiteEncode(f *testing.F) {
	families.FuzzCodecSuite(
		f,
		families.CodecChecks.Encode.Callable,
		codecIdle,
		conform.Excuse(
			families.CodecChecks.Encode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.Codec(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCodecSuiteDecode registers the fuzz target of Codec.Decode on
// codecIdle, with every check of the suite excused.
func FuzzCodecSuiteDecode(f *testing.F) {
	families.FuzzCodecSuite(
		f,
		families.CodecChecks.Decode.Callable,
		codecIdle,
		conform.Excuse(
			families.CodecChecks.Encode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Encode.Codec(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.CodecChecks.Decode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCodecSuite measures each callable of the suite of Codec on
// codecIdle.
func BenchmarkCodecSuite(b *testing.B) {
	b.Run("Encode", func(b *testing.B) {
		families.BenchCodecSuite(b, families.CodecChecks.Encode.Callable, codecIdle)
	})
	b.Run("Decode", func(b *testing.B) {
		families.BenchCodecSuite(b, families.CodecChecks.Decode.Callable, codecIdle)
	})
}

// pagerSubjectEmbedded is families.Pager under the name of the field that
// pagerSubject embeds.
type pagerSubjectEmbedded = families.Pager

// pagerSubject is a value of Pager that only TestPagerSuite builds.
type pagerSubject struct{ pagerSubjectEmbedded }

// TestPagerSuite checks the suite of Pager and the two forms of its
// application.
func TestPagerSuite(t *testing.T) {
	t.Parallel()

	t.Run("PagerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.PagerSuite())
	})

	t.Run("RunPagerSuite", func(t *testing.T) {
		t.Parallel()
		var built *pagerSubject
		families.RunPagerSuite(
			t,
			func() *pagerSubject {
				built = &pagerSubject{}
				return built
			},
			conform.With(families.PagerCheck{
				ID:        "Pager/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.PagerHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Pager(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.PagerHandle) families.Pager { return &pagerSubject{} },
			}),
			conform.Excuse(
				families.PagerChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PagerChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PagerChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PagerChecks.Fetch.Pagination(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPagerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *pagerSubject
		families.RunPagerSuiteFromHandle(
			t,
			func(*families.PagerHandle) *pagerSubject {
				built = &pagerSubject{}
				return built
			},
			conform.With(families.PagerCheck{
				ID:        "Pager/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.PagerHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Pager(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.PagerHandle) families.Pager { return &pagerSubject{} },
			}),
			conform.Excuse(
				families.PagerChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PagerChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PagerChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PagerChecks.Fetch.Pagination(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// pagerIdle returns the planted defect of a ceiling of Pager with a spend
// of nothing, a subject that answers each call at once.
func pagerIdle(h *families.PagerHandle) families.Pager {
	return families.PagerSuite().Calls[0].Over(h, func() {})
}

// TestPagerSuiteBudgets applies the allocation ceilings of the suite of
// Pager to pagerIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPagerSuiteBudgets(t *testing.T) {
	families.RunPagerBudgets(
		t,
		pagerIdle,
		conform.Ceilings(
			families.PagerChecks.Fetch.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.PagerChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PagerChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PagerChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PagerChecks.Fetch.Pagination(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPagerSuiteFetch registers the fuzz target of Pager.Fetch on
// pagerIdle, with every check of the suite excused.
func FuzzPagerSuiteFetch(f *testing.F) {
	families.FuzzPagerSuite(
		f,
		families.PagerChecks.Fetch.Callable,
		pagerIdle,
		conform.Excuse(
			families.PagerChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PagerChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PagerChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PagerChecks.Fetch.Pagination(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPagerSuite measures each callable of the suite of Pager on
// pagerIdle.
func BenchmarkPagerSuite(b *testing.B) {
	b.Run("Fetch", func(b *testing.B) {
		families.BenchPagerSuite(b, families.PagerChecks.Fetch.Callable, pagerIdle)
	})
}

// logSubjectEmbedded is families.Log under the name of the field that
// logSubject embeds.
type logSubjectEmbedded = families.Log

// logSubject is a value of Log that only TestLogSuite builds.
type logSubject struct{ logSubjectEmbedded }

// TestLogSuite checks the suite of Log and the two forms of its
// application.
func TestLogSuite(t *testing.T) {
	t.Parallel()

	t.Run("LogSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.LogSuite())
	})

	t.Run("RunLogSuite", func(t *testing.T) {
		t.Parallel()
		var built *logSubject
		families.RunLogSuite(
			t,
			func() *logSubject {
				built = &logSubject{}
				return built
			},
			conform.With(families.LogCheck{
				ID:        "Log/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LogHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Log(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.LogHandle) families.Log { return &logSubject{} },
			}),
			conform.Excuse(
				families.LogChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Append.ChainCausal(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Append.Chain(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Replay.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Replay.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Replay.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLogSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *logSubject
		families.RunLogSuiteFromHandle(
			t,
			func(*families.LogHandle) *logSubject {
				built = &logSubject{}
				return built
			},
			conform.With(families.LogCheck{
				ID:        "Log/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LogHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Log(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.LogHandle) families.Log { return &logSubject{} },
			}),
			conform.Excuse(
				families.LogChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Append.ChainCausal(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Append.Chain(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Replay.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Replay.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LogChecks.Replay.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// logIdle returns the planted defect of a ceiling of Log with a spend of
// nothing, a subject that answers each call at once.
func logIdle(h *families.LogHandle) families.Log {
	return families.LogSuite().Calls[0].Over(h, func() {})
}

// TestLogSuiteBudgets applies the allocation ceilings of the suite of Log
// to logIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestLogSuiteBudgets(t *testing.T) {
	families.RunLogBudgets(
		t,
		logIdle,
		conform.Ceilings(
			families.LogChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.LogChecks.Replay.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.ChainCausal(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLogSuiteAppend registers the fuzz target of Log.Append on logIdle,
// with every check of the suite excused.
func FuzzLogSuiteAppend(f *testing.F) {
	families.FuzzLogSuite(
		f,
		families.LogChecks.Append.Callable,
		logIdle,
		conform.Excuse(
			families.LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.ChainCausal(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLogSuiteReplay registers the fuzz target of Log.Replay on logIdle,
// with every check of the suite excused.
func FuzzLogSuiteReplay(f *testing.F) {
	families.FuzzLogSuite(
		f,
		families.LogChecks.Replay.Callable,
		logIdle,
		conform.Excuse(
			families.LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.ChainCausal(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LogChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLogSuite measures each callable of the suite of Log on logIdle.
func BenchmarkLogSuite(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		families.BenchLogSuite(b, families.LogChecks.Append.Callable, logIdle)
	})
	b.Run("Replay", func(b *testing.B) {
		families.BenchLogSuite(b, families.LogChecks.Replay.Callable, logIdle)
	})
}

// openerSubjectEmbedded is families.Opener under the name of the field
// that openerSubject embeds.
type openerSubjectEmbedded = families.Opener

// openerSubject is a value of Opener that only TestOpenerSuite builds.
type openerSubject struct{ openerSubjectEmbedded }

// TestOpenerSuite checks the suite of Opener and the two forms of its
// application.
func TestOpenerSuite(t *testing.T) {
	t.Parallel()

	t.Run("OpenerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.OpenerSuite())
	})

	t.Run("RunOpenerSuite", func(t *testing.T) {
		t.Parallel()
		var built *openerSubject
		families.RunOpenerSuite(
			t,
			func() *openerSubject {
				built = &openerSubject{}
				return built
			},
			conform.With(families.OpenerCheck{
				ID:        "Opener/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.OpenerHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Opener(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.OpenerHandle) families.Opener { return &openerSubject{} },
			}),
			conform.Excuse(
				families.OpenerChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.OpenerChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.OpenerChecks.Open.CursorOpen(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunOpenerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *openerSubject
		families.RunOpenerSuiteFromHandle(
			t,
			func(*families.OpenerHandle) *openerSubject {
				built = &openerSubject{}
				return built
			},
			conform.With(families.OpenerCheck{
				ID:        "Opener/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.OpenerHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Opener(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.OpenerHandle) families.Opener { return &openerSubject{} },
			}),
			conform.Excuse(
				families.OpenerChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.OpenerChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.OpenerChecks.Open.CursorOpen(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// openerIdle returns the planted defect of a ceiling of Opener with a
// spend of nothing, a subject that answers each call at once.
func openerIdle(h *families.OpenerHandle) families.Opener {
	return families.OpenerSuite().Calls[0].Over(h, func() {})
}

// TestOpenerSuiteBudgets applies the allocation ceilings of the suite of
// Opener to openerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestOpenerSuiteBudgets(t *testing.T) {
	families.RunOpenerBudgets(
		t,
		openerIdle,
		conform.Ceilings(
			families.OpenerChecks.Open.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.OpenerChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.OpenerChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.OpenerChecks.Open.CursorOpen(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOpenerSuiteOpen registers the fuzz target of Opener.Open on
// openerIdle, with every check of the suite excused.
func FuzzOpenerSuiteOpen(f *testing.F) {
	families.FuzzOpenerSuite(
		f,
		families.OpenerChecks.Open.Callable,
		openerIdle,
		conform.Excuse(
			families.OpenerChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.OpenerChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.OpenerChecks.Open.CursorOpen(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkOpenerSuite measures each callable of the suite of Opener on
// openerIdle.
func BenchmarkOpenerSuite(b *testing.B) {
	b.Run("Open", func(b *testing.B) {
		families.BenchOpenerSuite(b, families.OpenerChecks.Open.Callable, openerIdle)
	})
}

// ledgerSubjectEmbedded is families.Ledger under the name of the field
// that ledgerSubject embeds.
type ledgerSubjectEmbedded = families.Ledger

// ledgerSubject is a value of Ledger that only TestLedgerSuite builds.
type ledgerSubject struct{ ledgerSubjectEmbedded }

// TestLedgerSuite checks the suite of Ledger and the two forms of its
// application.
func TestLedgerSuite(t *testing.T) {
	t.Parallel()

	t.Run("LedgerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.LedgerSuite())
	})

	t.Run("RunLedgerSuite", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		families.RunLedgerSuite(
			t,
			func() *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(families.LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.LedgerHandle) families.Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				families.LedgerChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Begin.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Begin.Tx(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLedgerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		families.RunLedgerSuiteFromHandle(
			t,
			func(*families.LedgerHandle) *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(families.LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.LedgerHandle) families.Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				families.LedgerChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Begin.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Begin.Tx(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LedgerChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// ledgerIdle returns the planted defect of a ceiling of Ledger with a
// spend of nothing, a subject that answers each call at once.
func ledgerIdle(h *families.LedgerHandle) families.Ledger {
	return families.LedgerSuite().Calls[0].Over(h, func() {})
}

// TestLedgerSuiteBudgets applies the allocation ceilings of the suite of
// Ledger to ledgerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLedgerSuiteBudgets(t *testing.T) {
	families.RunLedgerBudgets(
		t,
		ledgerIdle,
		conform.Ceilings(
			families.LedgerChecks.Begin.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.LedgerChecks.Commit.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.LedgerChecks.Rollback.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteBegin registers the fuzz target of Ledger.Begin on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteBegin(f *testing.F) {
	families.FuzzLedgerSuite(
		f,
		families.LedgerChecks.Begin.Callable,
		ledgerIdle,
		conform.Excuse(
			families.LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteCommit registers the fuzz target of Ledger.Commit on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteCommit(f *testing.F) {
	families.FuzzLedgerSuite(
		f,
		families.LedgerChecks.Commit.Callable,
		ledgerIdle,
		conform.Excuse(
			families.LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteRollback registers the fuzz target of Ledger.Rollback on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteRollback(f *testing.F) {
	families.FuzzLedgerSuite(
		f,
		families.LedgerChecks.Rollback.Callable,
		ledgerIdle,
		conform.Excuse(
			families.LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLedgerSuite measures each callable of the suite of Ledger on
// ledgerIdle.
func BenchmarkLedgerSuite(b *testing.B) {
	b.Run("Begin", func(b *testing.B) {
		families.BenchLedgerSuite(b, families.LedgerChecks.Begin.Callable, ledgerIdle)
	})
	b.Run("Commit", func(b *testing.B) {
		families.BenchLedgerSuite(b, families.LedgerChecks.Commit.Callable, ledgerIdle)
	})
	b.Run("Rollback", func(b *testing.B) {
		families.BenchLedgerSuite(b, families.LedgerChecks.Rollback.Callable, ledgerIdle)
	})
}

// guardedSubjectEmbedded is families.Guarded under the name of the field
// that guardedSubject embeds.
type guardedSubjectEmbedded = families.Guarded

// guardedSubject is a value of Guarded that only TestGuardedSuite builds.
type guardedSubject struct{ guardedSubjectEmbedded }

// TestGuardedSuite checks the suite of Guarded and the two forms of its
// application.
func TestGuardedSuite(t *testing.T) {
	t.Parallel()

	t.Run("GuardedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.GuardedSuite())
	})

	t.Run("RunGuardedSuite", func(t *testing.T) {
		t.Parallel()
		var built *guardedSubject
		families.RunGuardedSuite(
			t,
			func() *guardedSubject {
				built = &guardedSubject{}
				return built
			},
			conform.With(families.GuardedCheck{
				ID:        "Guarded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.GuardedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Guarded(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.GuardedHandle) families.Guarded { return &guardedSubject{} },
			}),
			conform.Excuse(
				families.GuardedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Put.IfMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Match.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Match.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunGuardedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *guardedSubject
		families.RunGuardedSuiteFromHandle(
			t,
			func(*families.GuardedHandle) *guardedSubject {
				built = &guardedSubject{}
				return built
			},
			conform.With(families.GuardedCheck{
				ID:        "Guarded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.GuardedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Guarded(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.GuardedHandle) families.Guarded { return &guardedSubject{} },
			}),
			conform.Excuse(
				families.GuardedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Put.IfMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Match.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.GuardedChecks.Match.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// guardedIdle returns the planted defect of a ceiling of Guarded with a
// spend of nothing, a subject that answers each call at once.
func guardedIdle(h *families.GuardedHandle) families.Guarded {
	return families.GuardedSuite().Calls[0].Over(h, func() {})
}

// TestGuardedSuiteBudgets applies the allocation ceilings of the suite of
// Guarded to guardedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestGuardedSuiteBudgets(t *testing.T) {
	families.RunGuardedBudgets(
		t,
		guardedIdle,
		conform.Ceilings(
			families.GuardedChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.GuardedChecks.Match.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.GuardedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Put.IfMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Match.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Match.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzGuardedSuitePut registers the fuzz target of Guarded.Put on
// guardedIdle, with every check of the suite excused.
func FuzzGuardedSuitePut(f *testing.F) {
	families.FuzzGuardedSuite(
		f,
		families.GuardedChecks.Put.Callable,
		guardedIdle,
		conform.Excuse(
			families.GuardedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Put.IfMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Match.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Match.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzGuardedSuiteMatch registers the fuzz target of Guarded.Match on
// guardedIdle, with every check of the suite excused.
func FuzzGuardedSuiteMatch(f *testing.F) {
	families.FuzzGuardedSuite(
		f,
		families.GuardedChecks.Match.Callable,
		guardedIdle,
		conform.Excuse(
			families.GuardedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Put.IfMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Match.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.GuardedChecks.Match.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkGuardedSuite measures each callable of the suite of Guarded on
// guardedIdle.
func BenchmarkGuardedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchGuardedSuite(b, families.GuardedChecks.Put.Callable, guardedIdle)
	})
	b.Run("Match", func(b *testing.B) {
		families.BenchGuardedSuite(b, families.GuardedChecks.Match.Callable, guardedIdle)
	})
}

// scopeSubjectEmbedded is families.Scope under the name of the field that
// scopeSubject embeds.
type scopeSubjectEmbedded = families.Scope

// scopeSubject is a value of Scope that only TestScopeSuite builds.
type scopeSubject struct{ scopeSubjectEmbedded }

// TestScopeSuite checks the suite of Scope and the two forms of its
// application.
func TestScopeSuite(t *testing.T) {
	t.Parallel()

	t.Run("ScopeSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.ScopeSuite())
	})

	t.Run("RunScopeSuite", func(t *testing.T) {
		t.Parallel()
		var built *scopeSubject
		families.RunScopeSuite(
			t,
			func() *scopeSubject {
				built = &scopeSubject{}
				return built
			},
			conform.With(families.ScopeCheck{
				ID:        "Scope/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ScopeHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Scope(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.ScopeHandle) families.Scope { return &scopeSubject{} },
			}),
			conform.Excuse(
				families.ScopeChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Run.Transaction(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunScopeSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *scopeSubject
		families.RunScopeSuiteFromHandle(
			t,
			func(*families.ScopeHandle) *scopeSubject {
				built = &scopeSubject{}
				return built
			},
			conform.With(families.ScopeCheck{
				ID:        "Scope/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ScopeHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Scope(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.ScopeHandle) families.Scope { return &scopeSubject{} },
			}),
			conform.Excuse(
				families.ScopeChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Run.Transaction(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopeChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// scopeIdle returns the planted defect of a ceiling of Scope with a spend
// of nothing, a subject that answers each call at once.
func scopeIdle(h *families.ScopeHandle) families.Scope {
	return families.ScopeSuite().Calls[0].Over(h, func() {})
}

// TestScopeSuiteBudgets applies the allocation ceilings of the suite of
// Scope to scopeIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestScopeSuiteBudgets(t *testing.T) {
	families.RunScopeBudgets(
		t,
		scopeIdle,
		conform.Ceilings(families.ScopeChecks.Run.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(families.ScopeChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(families.ScopeChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			families.ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopeSuiteRun registers the fuzz target of Scope.Run on scopeIdle,
// with every check of the suite excused.
func FuzzScopeSuiteRun(f *testing.F) {
	families.FuzzScopeSuite(
		f,
		families.ScopeChecks.Run.Callable,
		scopeIdle,
		conform.Excuse(
			families.ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopeSuitePut registers the fuzz target of Scope.Put on scopeIdle,
// with every check of the suite excused.
func FuzzScopeSuitePut(f *testing.F) {
	families.FuzzScopeSuite(
		f,
		families.ScopeChecks.Put.Callable,
		scopeIdle,
		conform.Excuse(
			families.ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopeSuiteGet registers the fuzz target of Scope.Get on scopeIdle,
// with every check of the suite excused.
func FuzzScopeSuiteGet(f *testing.F) {
	families.FuzzScopeSuite(
		f,
		families.ScopeChecks.Get.Callable,
		scopeIdle,
		conform.Excuse(
			families.ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkScopeSuite measures each callable of the suite of Scope on
// scopeIdle.
func BenchmarkScopeSuite(b *testing.B) {
	b.Run("Run", func(b *testing.B) {
		families.BenchScopeSuite(b, families.ScopeChecks.Run.Callable, scopeIdle)
	})
	b.Run("Put", func(b *testing.B) {
		families.BenchScopeSuite(b, families.ScopeChecks.Put.Callable, scopeIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchScopeSuite(b, families.ScopeChecks.Get.Callable, scopeIdle)
	})
}

// eventsSubjectEmbedded is families.Events under the name of the field
// that eventsSubject embeds.
type eventsSubjectEmbedded = families.Events

// eventsSubject is a value of Events that only TestEventsSuite builds.
type eventsSubject struct{ eventsSubjectEmbedded }

// TestEventsSuite checks the suite of Events and the two forms of its
// application.
func TestEventsSuite(t *testing.T) {
	t.Parallel()

	t.Run("EventsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.EventsSuite())
	})

	t.Run("RunEventsSuite", func(t *testing.T) {
		t.Parallel()
		var built *eventsSubject
		families.RunEventsSuite(
			t,
			func() *eventsSubject {
				built = &eventsSubject{}
				return built
			},
			conform.With(families.EventsCheck{
				ID:        "Events/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.EventsHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Events(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.EventsHandle) families.Events { return &eventsSubject{} },
			}),
			conform.Excuse(
				families.EventsChecks.Fire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.EventsChecks.Fire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.EventsChecks.Fire.Hooks(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.EventsChecks.OnEvent.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunEventsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *eventsSubject
		families.RunEventsSuiteFromHandle(
			t,
			func(*families.EventsHandle) *eventsSubject {
				built = &eventsSubject{}
				return built
			},
			conform.With(families.EventsCheck{
				ID:        "Events/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.EventsHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Events(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.EventsHandle) families.Events { return &eventsSubject{} },
			}),
			conform.Excuse(
				families.EventsChecks.Fire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.EventsChecks.Fire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.EventsChecks.Fire.Hooks(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.EventsChecks.OnEvent.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// eventsIdle returns the planted defect of a ceiling of Events with a
// spend of nothing, a subject that answers each call at once.
func eventsIdle(h *families.EventsHandle) families.Events {
	return families.EventsSuite().Calls[0].Over(h, func() {})
}

// TestEventsSuiteBudgets applies the allocation ceilings of the suite of
// Events to eventsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestEventsSuiteBudgets(t *testing.T) {
	families.RunEventsBudgets(
		t,
		eventsIdle,
		conform.Ceilings(
			families.EventsChecks.Fire.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.EventsChecks.OnEvent.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.EventsChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.Fire.Hooks(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.OnEvent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzEventsSuiteFire registers the fuzz target of Events.Fire on
// eventsIdle, with every check of the suite excused.
func FuzzEventsSuiteFire(f *testing.F) {
	families.FuzzEventsSuite(
		f,
		families.EventsChecks.Fire.Callable,
		eventsIdle,
		conform.Excuse(
			families.EventsChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.Fire.Hooks(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.OnEvent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzEventsSuiteOnEvent registers the fuzz target of Events.OnEvent on
// eventsIdle, with every check of the suite excused.
func FuzzEventsSuiteOnEvent(f *testing.F) {
	families.FuzzEventsSuite(
		f,
		families.EventsChecks.OnEvent.Callable,
		eventsIdle,
		conform.Excuse(
			families.EventsChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.Fire.Hooks(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.EventsChecks.OnEvent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkEventsSuite measures each callable of the suite of Events on
// eventsIdle.
func BenchmarkEventsSuite(b *testing.B) {
	b.Run("Fire", func(b *testing.B) {
		families.BenchEventsSuite(b, families.EventsChecks.Fire.Callable, eventsIdle)
	})
	b.Run("OnEvent", func(b *testing.B) {
		families.BenchEventsSuite(b, families.EventsChecks.OnEvent.Callable, eventsIdle)
	})
}

// timersSubjectEmbedded is families.Timers under the name of the field
// that timersSubject embeds.
type timersSubjectEmbedded = families.Timers

// timersSubject is a value of Timers that only TestTimersSuite builds.
type timersSubject struct{ timersSubjectEmbedded }

// TestTimersSuite checks the suite of Timers and the two forms of its
// application.
func TestTimersSuite(t *testing.T) {
	t.Parallel()

	t.Run("TimersSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.TimersSuite())
	})

	t.Run("RunTimersSuite", func(t *testing.T) {
		t.Parallel()
		var built *timersSubject
		families.RunTimersSuite(
			t,
			func() *timersSubject {
				built = &timersSubject{}
				return built
			},
			conform.With(families.TimersCheck{
				ID:        "Timers/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.TimersHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Timers(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.TimersHandle) families.Timers { return &timersSubject{} },
			}),
			conform.Excuse(
				families.TimersChecks.At.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.At.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.At.Scheduled(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Fired.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Fired.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Fired.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Slow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Slow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Slow.Within(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTimersSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *timersSubject
		families.RunTimersSuiteFromHandle(
			t,
			func(*families.TimersHandle) *timersSubject {
				built = &timersSubject{}
				return built
			},
			conform.With(families.TimersCheck{
				ID:        "Timers/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.TimersHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Timers(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.TimersHandle) families.Timers { return &timersSubject{} },
			}),
			conform.Excuse(
				families.TimersChecks.At.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.At.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.At.Scheduled(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Fired.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Fired.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Fired.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Slow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Slow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.TimersChecks.Slow.Within(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// timersIdle returns the planted defect of a ceiling of Timers with a
// spend of nothing, a subject that answers each call at once.
func timersIdle(h *families.TimersHandle) families.Timers {
	return families.TimersSuite().Calls[0].Over(h, func() {})
}

// TestTimersSuiteBudgets applies the allocation ceilings of the suite of
// Timers to timersIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestTimersSuiteBudgets(t *testing.T) {
	families.RunTimersBudgets(
		t,
		timersIdle,
		conform.Ceilings(families.TimersChecks.At.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			families.TimersChecks.Fired.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.TimersChecks.Slow.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTimersSuiteAt registers the fuzz target of Timers.At on timersIdle,
// with every check of the suite excused.
func FuzzTimersSuiteAt(f *testing.F) {
	families.FuzzTimersSuite(
		f,
		families.TimersChecks.At.Callable,
		timersIdle,
		conform.Excuse(
			families.TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTimersSuiteFired registers the fuzz target of Timers.Fired on
// timersIdle, with every check of the suite excused.
func FuzzTimersSuiteFired(f *testing.F) {
	families.FuzzTimersSuite(
		f,
		families.TimersChecks.Fired.Callable,
		timersIdle,
		conform.Excuse(
			families.TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTimersSuiteSlow registers the fuzz target of Timers.Slow on
// timersIdle, with every check of the suite excused.
func FuzzTimersSuiteSlow(f *testing.F) {
	families.FuzzTimersSuite(
		f,
		families.TimersChecks.Slow.Callable,
		timersIdle,
		conform.Excuse(
			families.TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkTimersSuite measures each callable of the suite of Timers on
// timersIdle.
func BenchmarkTimersSuite(b *testing.B) {
	b.Run("At", func(b *testing.B) {
		families.BenchTimersSuite(b, families.TimersChecks.At.Callable, timersIdle)
	})
	b.Run("Fired", func(b *testing.B) {
		families.BenchTimersSuite(b, families.TimersChecks.Fired.Callable, timersIdle)
	})
	b.Run("Slow", func(b *testing.B) {
		families.BenchTimersSuite(b, families.TimersChecks.Slow.Callable, timersIdle)
	})
}

// electionSubjectEmbedded is families.Election under the name of the field
// that electionSubject embeds.
type electionSubjectEmbedded = families.Election

// electionSubject is a value of Election that only TestElectionSuite
// builds.
type electionSubject struct{ electionSubjectEmbedded }

// TestElectionSuite checks the suite of Election and the two forms of its
// application.
func TestElectionSuite(t *testing.T) {
	t.Parallel()

	t.Run("ElectionSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.ElectionSuite())
	})

	t.Run("RunElectionSuite", func(t *testing.T) {
		t.Parallel()
		var built *electionSubject
		families.RunElectionSuite(
			t,
			func() *electionSubject {
				built = &electionSubject{}
				return built
			},
			conform.With(families.ElectionCheck{
				ID:        "Election/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ElectionHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Election(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.ElectionHandle) families.Election { return &electionSubject{} },
			}),
			conform.Excuse(
				families.ElectionChecks.Campaign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Campaign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Campaign.LeaderElection(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Resign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Resign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.IsLeader.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunElectionSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *electionSubject
		families.RunElectionSuiteFromHandle(
			t,
			func(*families.ElectionHandle) *electionSubject {
				built = &electionSubject{}
				return built
			},
			conform.With(families.ElectionCheck{
				ID:        "Election/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ElectionHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Election(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.ElectionHandle) families.Election { return &electionSubject{} },
			}),
			conform.Excuse(
				families.ElectionChecks.Campaign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Campaign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Campaign.LeaderElection(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Resign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.Resign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ElectionChecks.IsLeader.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// electionIdle returns the planted defect of a ceiling of Election with a
// spend of nothing, a subject that answers each call at once.
func electionIdle(h *families.ElectionHandle) families.Election {
	return families.ElectionSuite().Calls[0].Over(h, func() {})
}

// TestElectionSuiteBudgets applies the allocation ceilings of the suite of
// Election to electionIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestElectionSuiteBudgets(t *testing.T) {
	families.RunElectionBudgets(
		t,
		electionIdle,
		conform.Ceilings(
			families.ElectionChecks.Campaign.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.ElectionChecks.Resign.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.ElectionChecks.IsLeader.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzElectionSuiteCampaign registers the fuzz target of Election.Campaign
// on electionIdle, with every check of the suite excused.
func FuzzElectionSuiteCampaign(f *testing.F) {
	families.FuzzElectionSuite(
		f,
		families.ElectionChecks.Campaign.Callable,
		electionIdle,
		conform.Excuse(
			families.ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzElectionSuiteResign registers the fuzz target of Election.Resign on
// electionIdle, with every check of the suite excused.
func FuzzElectionSuiteResign(f *testing.F) {
	families.FuzzElectionSuite(
		f,
		families.ElectionChecks.Resign.Callable,
		electionIdle,
		conform.Excuse(
			families.ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzElectionSuiteIsLeader registers the fuzz target of Election.IsLeader
// on electionIdle, with every check of the suite excused.
func FuzzElectionSuiteIsLeader(f *testing.F) {
	families.FuzzElectionSuite(
		f,
		families.ElectionChecks.IsLeader.Callable,
		electionIdle,
		conform.Excuse(
			families.ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkElectionSuite measures each callable of the suite of Election
// on electionIdle.
func BenchmarkElectionSuite(b *testing.B) {
	b.Run("Campaign", func(b *testing.B) {
		families.BenchElectionSuite(b, families.ElectionChecks.Campaign.Callable, electionIdle)
	})
	b.Run("Resign", func(b *testing.B) {
		families.BenchElectionSuite(b, families.ElectionChecks.Resign.Callable, electionIdle)
	})
	b.Run("IsLeader", func(b *testing.B) {
		families.BenchElectionSuite(b, families.ElectionChecks.IsLeader.Callable, electionIdle)
	})
}

// poolSubjectEmbedded is families.Pool under the name of the field that
// poolSubject embeds.
type poolSubjectEmbedded = families.Pool

// poolSubject is a value of Pool that only TestPoolSuite builds.
type poolSubject struct{ poolSubjectEmbedded }

// TestPoolSuite checks the suite of Pool and the two forms of its
// application.
func TestPoolSuite(t *testing.T) {
	t.Parallel()

	t.Run("PoolSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.PoolSuite())
	})

	t.Run("RunPoolSuite", func(t *testing.T) {
		t.Parallel()
		var built *poolSubject
		families.RunPoolSuite(
			t,
			func() *poolSubject {
				built = &poolSubject{}
				return built
			},
			conform.With(families.PoolCheck{
				ID:        "Pool/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.PoolHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Pool(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.PoolHandle) families.Pool { return &poolSubject{} },
			}),
			conform.Excuse(
				families.PoolChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Get.Pool(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPoolSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *poolSubject
		families.RunPoolSuiteFromHandle(
			t,
			func(*families.PoolHandle) *poolSubject {
				built = &poolSubject{}
				return built
			},
			conform.With(families.PoolCheck{
				ID:        "Pool/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.PoolHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Pool(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.PoolHandle) families.Pool { return &poolSubject{} },
			}),
			conform.Excuse(
				families.PoolChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Get.Pool(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PoolChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// poolIdle returns the planted defect of a ceiling of Pool with a spend of
// nothing, a subject that answers each call at once.
func poolIdle(h *families.PoolHandle) families.Pool {
	return families.PoolSuite().Calls[0].Over(h, func() {})
}

// TestPoolSuiteBudgets applies the allocation ceilings of the suite of
// Pool to poolIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestPoolSuiteBudgets(t *testing.T) {
	families.RunPoolBudgets(
		t,
		poolIdle,
		conform.Ceilings(families.PoolChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(families.PoolChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			families.PoolChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPoolSuiteGet registers the fuzz target of Pool.Get on poolIdle, with
// every check of the suite excused.
func FuzzPoolSuiteGet(f *testing.F) {
	families.FuzzPoolSuite(
		f,
		families.PoolChecks.Get.Callable,
		poolIdle,
		conform.Excuse(
			families.PoolChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPoolSuitePut registers the fuzz target of Pool.Put on poolIdle, with
// every check of the suite excused.
func FuzzPoolSuitePut(f *testing.F) {
	families.FuzzPoolSuite(
		f,
		families.PoolChecks.Put.Callable,
		poolIdle,
		conform.Excuse(
			families.PoolChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Get.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PoolChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPoolSuite measures each callable of the suite of Pool on
// poolIdle.
func BenchmarkPoolSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		families.BenchPoolSuite(b, families.PoolChecks.Get.Callable, poolIdle)
	})
	b.Run("Put", func(b *testing.B) {
		families.BenchPoolSuite(b, families.PoolChecks.Put.Callable, poolIdle)
	})
}

// batchSubjectEmbedded is families.Batch under the name of the field that
// batchSubject embeds.
type batchSubjectEmbedded = families.Batch

// batchSubject is a value of Batch that only TestBatchSuite builds.
type batchSubject struct{ batchSubjectEmbedded }

// TestBatchSuite checks the suite of Batch and the two forms of its
// application.
func TestBatchSuite(t *testing.T) {
	t.Parallel()

	t.Run("BatchSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.BatchSuite())
	})

	t.Run("RunBatchSuite", func(t *testing.T) {
		t.Parallel()
		var built *batchSubject
		families.RunBatchSuite(
			t,
			func() *batchSubject {
				built = &batchSubject{}
				return built
			},
			conform.With(families.BatchCheck{
				ID:        "Batch/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.BatchHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Batch(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.BatchHandle) families.Batch { return &batchSubject{} },
			}),
			conform.Excuse(
				families.BatchChecks.Write.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Write.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Write.BatchWriter(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBatchSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *batchSubject
		families.RunBatchSuiteFromHandle(
			t,
			func(*families.BatchHandle) *batchSubject {
				built = &batchSubject{}
				return built
			},
			conform.With(families.BatchCheck{
				ID:        "Batch/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.BatchHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Batch(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.BatchHandle) families.Batch { return &batchSubject{} },
			}),
			conform.Excuse(
				families.BatchChecks.Write.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Write.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Write.BatchWriter(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.BatchChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// batchIdle returns the planted defect of a ceiling of Batch with a spend
// of nothing, a subject that answers each call at once.
func batchIdle(h *families.BatchHandle) families.Batch {
	return families.BatchSuite().Calls[0].Over(h, func() {})
}

// TestBatchSuiteBudgets applies the allocation ceilings of the suite of
// Batch to batchIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestBatchSuiteBudgets(t *testing.T) {
	families.RunBatchBudgets(
		t,
		batchIdle,
		conform.Ceilings(
			families.BatchChecks.Write.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.BatchChecks.Read.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.BatchChecks.Write.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Write.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Write.BatchWriter(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBatchSuiteWrite registers the fuzz target of Batch.Write on
// batchIdle, with every check of the suite excused.
func FuzzBatchSuiteWrite(f *testing.F) {
	families.FuzzBatchSuite(
		f,
		families.BatchChecks.Write.Callable,
		batchIdle,
		conform.Excuse(
			families.BatchChecks.Write.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Write.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Write.BatchWriter(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBatchSuiteRead registers the fuzz target of Batch.Read on batchIdle,
// with every check of the suite excused.
func FuzzBatchSuiteRead(f *testing.F) {
	families.FuzzBatchSuite(
		f,
		families.BatchChecks.Read.Callable,
		batchIdle,
		conform.Excuse(
			families.BatchChecks.Write.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Write.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Write.BatchWriter(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.BatchChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBatchSuite measures each callable of the suite of Batch on
// batchIdle.
func BenchmarkBatchSuite(b *testing.B) {
	b.Run("Write", func(b *testing.B) {
		families.BenchBatchSuite(b, families.BatchChecks.Write.Callable, batchIdle)
	})
	b.Run("Read", func(b *testing.B) {
		families.BenchBatchSuite(b, families.BatchChecks.Read.Callable, batchIdle)
	})
}

// screenedSubjectEmbedded is families.Screened under the name of the field
// that screenedSubject embeds.
type screenedSubjectEmbedded = families.Screened

// screenedSubject is a value of Screened that only TestScreenedSuite
// builds.
type screenedSubject struct{ screenedSubjectEmbedded }

// TestScreenedSuite checks the suite of Screened and the two forms of its
// application.
func TestScreenedSuite(t *testing.T) {
	t.Parallel()

	t.Run("ScreenedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.ScreenedSuite())
	})

	t.Run("RunScreenedSuite", func(t *testing.T) {
		t.Parallel()
		var built *screenedSubject
		families.RunScreenedSuite(
			t,
			func() *screenedSubject {
				built = &screenedSubject{}
				return built
			},
			conform.With(families.ScreenedCheck{
				ID:        "Screened/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ScreenedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Screened(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.ScreenedHandle) families.Screened { return &screenedSubject{} },
			}),
			conform.Excuse(
				families.ScreenedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScreenedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScreenedChecks.Store.Validates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScreenedChecks.Validate.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunScreenedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *screenedSubject
		families.RunScreenedSuiteFromHandle(
			t,
			func(*families.ScreenedHandle) *screenedSubject {
				built = &screenedSubject{}
				return built
			},
			conform.With(families.ScreenedCheck{
				ID:        "Screened/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ScreenedHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Screened(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.ScreenedHandle) families.Screened { return &screenedSubject{} },
			}),
			conform.Excuse(
				families.ScreenedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScreenedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScreenedChecks.Store.Validates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScreenedChecks.Validate.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// screenedIdle returns the planted defect of a ceiling of Screened with a
// spend of nothing, a subject that answers each call at once.
func screenedIdle(h *families.ScreenedHandle) families.Screened {
	return families.ScreenedSuite().Calls[0].Over(h, func() {})
}

// TestScreenedSuiteBudgets applies the allocation ceilings of the suite of
// Screened to screenedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestScreenedSuiteBudgets(t *testing.T) {
	families.RunScreenedBudgets(
		t,
		screenedIdle,
		conform.Ceilings(
			families.ScreenedChecks.Store.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.ScreenedChecks.Validate.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.ScreenedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Store.Validates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Validate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScreenedSuiteStore registers the fuzz target of Screened.Store on
// screenedIdle, with every check of the suite excused.
func FuzzScreenedSuiteStore(f *testing.F) {
	families.FuzzScreenedSuite(
		f,
		families.ScreenedChecks.Store.Callable,
		screenedIdle,
		conform.Excuse(
			families.ScreenedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Store.Validates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Validate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScreenedSuiteValidate registers the fuzz target of Screened.Validate
// on screenedIdle, with every check of the suite excused.
func FuzzScreenedSuiteValidate(f *testing.F) {
	families.FuzzScreenedSuite(
		f,
		families.ScreenedChecks.Validate.Callable,
		screenedIdle,
		conform.Excuse(
			families.ScreenedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Store.Validates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScreenedChecks.Validate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkScreenedSuite measures each callable of the suite of Screened
// on screenedIdle.
func BenchmarkScreenedSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		families.BenchScreenedSuite(b, families.ScreenedChecks.Store.Callable, screenedIdle)
	})
	b.Run("Validate", func(b *testing.B) {
		families.BenchScreenedSuite(b, families.ScreenedChecks.Validate.Callable, screenedIdle)
	})
}

// limiterSubjectEmbedded is families.Limiter under the name of the field
// that limiterSubject embeds.
type limiterSubjectEmbedded = families.Limiter

// limiterSubject is a value of Limiter that only TestLimiterSuite builds.
type limiterSubject struct{ limiterSubjectEmbedded }

// TestLimiterSuite checks the suite of Limiter and the two forms of its
// application.
func TestLimiterSuite(t *testing.T) {
	t.Parallel()

	t.Run("LimiterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.LimiterSuite())
	})

	t.Run("RunLimiterSuite", func(t *testing.T) {
		t.Parallel()
		var built *limiterSubject
		families.RunLimiterSuite(
			t,
			func() *limiterSubject {
				built = &limiterSubject{}
				return built
			},
			conform.With(families.LimiterCheck{
				ID:        "Limiter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LimiterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Limiter(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.LimiterHandle) families.Limiter { return &limiterSubject{} },
			}),
			conform.Excuse(
				families.LimiterChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LimiterChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LimiterChecks.Run.RateLimit(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLimiterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *limiterSubject
		families.RunLimiterSuiteFromHandle(
			t,
			func(*families.LimiterHandle) *limiterSubject {
				built = &limiterSubject{}
				return built
			},
			conform.With(families.LimiterCheck{
				ID:        "Limiter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LimiterHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Limiter(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.LimiterHandle) families.Limiter { return &limiterSubject{} },
			}),
			conform.Excuse(
				families.LimiterChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LimiterChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LimiterChecks.Run.RateLimit(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// limiterIdle returns the planted defect of a ceiling of Limiter with a
// spend of nothing, a subject that answers each call at once.
func limiterIdle(h *families.LimiterHandle) families.Limiter {
	return families.LimiterSuite().Calls[0].Over(h, func() {})
}

// TestLimiterSuiteBudgets applies the allocation ceilings of the suite of
// Limiter to limiterIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLimiterSuiteBudgets(t *testing.T) {
	families.RunLimiterBudgets(
		t,
		limiterIdle,
		conform.Ceilings(
			families.LimiterChecks.Run.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.LimiterChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LimiterChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LimiterChecks.Run.RateLimit(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimiterSuiteRun registers the fuzz target of Limiter.Run on
// limiterIdle, with every check of the suite excused.
func FuzzLimiterSuiteRun(f *testing.F) {
	families.FuzzLimiterSuite(
		f,
		families.LimiterChecks.Run.Callable,
		limiterIdle,
		conform.Excuse(
			families.LimiterChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LimiterChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LimiterChecks.Run.RateLimit(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLimiterSuite measures each callable of the suite of Limiter on
// limiterIdle.
func BenchmarkLimiterSuite(b *testing.B) {
	b.Run("Run", func(b *testing.B) {
		families.BenchLimiterSuite(b, families.LimiterChecks.Run.Callable, limiterIdle)
	})
}

// rendererSubjectEmbedded is families.Renderer under the name of the field
// that rendererSubject embeds.
type rendererSubjectEmbedded = families.Renderer

// rendererSubject is a value of Renderer that only TestRendererSuite
// builds.
type rendererSubject struct{ rendererSubjectEmbedded }

// TestRendererSuite checks the suite of Renderer and the two forms of its
// application.
func TestRendererSuite(t *testing.T) {
	t.Parallel()

	t.Run("RendererSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.RendererSuite())
	})

	t.Run("RunRendererSuite", func(t *testing.T) {
		t.Parallel()
		var built *rendererSubject
		families.RunRendererSuite(
			t,
			func() *rendererSubject {
				built = &rendererSubject{}
				return built
			},
			conform.With(families.RendererCheck{
				ID:        "Renderer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.RendererHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Renderer(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.RendererHandle) families.Renderer { return &rendererSubject{} },
			}),
			conform.Excuse(
				families.RendererChecks.Render.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Render.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Render.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Render.XSSSafe(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.Total(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunRendererSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *rendererSubject
		families.RunRendererSuiteFromHandle(
			t,
			func(*families.RendererHandle) *rendererSubject {
				built = &rendererSubject{}
				return built
			},
			conform.With(families.RendererCheck{
				ID:        "Renderer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.RendererHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Renderer(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.RendererHandle) families.Renderer { return &rendererSubject{} },
			}),
			conform.Excuse(
				families.RendererChecks.Render.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Render.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Render.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Render.XSSSafe(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.RendererChecks.Classify.Total(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// rendererIdle returns the planted defect of a ceiling of Renderer with a
// spend of nothing, a subject that answers each call at once.
func rendererIdle(h *families.RendererHandle) families.Renderer {
	return families.RendererSuite().Calls[0].Over(h, func() {})
}

// TestRendererSuiteBudgets applies the allocation ceilings of the suite of
// Renderer to rendererIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRendererSuiteBudgets(t *testing.T) {
	families.RunRendererBudgets(
		t,
		rendererIdle,
		conform.Ceilings(
			families.RendererChecks.Render.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.RendererChecks.Classify.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.RendererChecks.Render.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.XSSSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Total(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRendererSuiteRender registers the fuzz target of Renderer.Render on
// rendererIdle, with every check of the suite excused.
func FuzzRendererSuiteRender(f *testing.F) {
	families.FuzzRendererSuite(
		f,
		families.RendererChecks.Render.Callable,
		rendererIdle,
		conform.Excuse(
			families.RendererChecks.Render.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.XSSSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Total(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRendererSuiteClassify registers the fuzz target of Renderer.Classify
// on rendererIdle, with every check of the suite excused.
func FuzzRendererSuiteClassify(f *testing.F) {
	families.FuzzRendererSuite(
		f,
		families.RendererChecks.Classify.Callable,
		rendererIdle,
		conform.Excuse(
			families.RendererChecks.Render.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Render.XSSSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.RendererChecks.Classify.Total(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkRendererSuite measures each callable of the suite of Renderer
// on rendererIdle.
func BenchmarkRendererSuite(b *testing.B) {
	b.Run("Render", func(b *testing.B) {
		families.BenchRendererSuite(b, families.RendererChecks.Render.Callable, rendererIdle)
	})
	b.Run("Classify", func(b *testing.B) {
		families.BenchRendererSuite(b, families.RendererChecks.Classify.Callable, rendererIdle)
	})
}

// lifecycleSubjectEmbedded is families.Lifecycle under the name of the
// field that lifecycleSubject embeds.
type lifecycleSubjectEmbedded = families.Lifecycle

// lifecycleSubject is a value of Lifecycle that only TestLifecycleSuite
// builds.
type lifecycleSubject struct{ lifecycleSubjectEmbedded }

// TestLifecycleSuite checks the suite of Lifecycle and the two forms of
// its application.
func TestLifecycleSuite(t *testing.T) {
	t.Parallel()

	t.Run("LifecycleSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.LifecycleSuite())
	})

	t.Run("RunLifecycleSuite", func(t *testing.T) {
		t.Parallel()
		var built *lifecycleSubject
		families.RunLifecycleSuite(
			t,
			func() *lifecycleSubject {
				built = &lifecycleSubject{}
				return built
			},
			conform.With(families.LifecycleCheck{
				ID:        "Lifecycle/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LifecycleHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Lifecycle(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.LifecycleHandle) families.Lifecycle { return &lifecycleSubject{} },
			}),
			conform.Excuse(
				families.LifecycleChecks.Stop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Stop.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Put.Idempotent(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Sync.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Sync.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Sync.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLifecycleSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *lifecycleSubject
		families.RunLifecycleSuiteFromHandle(
			t,
			func(*families.LifecycleHandle) *lifecycleSubject {
				built = &lifecycleSubject{}
				return built
			},
			conform.With(families.LifecycleCheck{
				ID:        "Lifecycle/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.LifecycleHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Lifecycle(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.LifecycleHandle) families.Lifecycle { return &lifecycleSubject{} },
			}),
			conform.Excuse(
				families.LifecycleChecks.Stop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Stop.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Put.Idempotent(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Sync.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Sync.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.LifecycleChecks.Sync.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// lifecycleIdle returns the planted defect of a ceiling of Lifecycle with
// a spend of nothing, a subject that answers each call at once.
func lifecycleIdle(h *families.LifecycleHandle) families.Lifecycle {
	return families.LifecycleSuite().Calls[0].Over(h, func() {})
}

// TestLifecycleSuiteBudgets applies the allocation ceilings of the suite
// of Lifecycle to lifecycleIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLifecycleSuiteBudgets(t *testing.T) {
	families.RunLifecycleBudgets(
		t,
		lifecycleIdle,
		conform.Ceilings(
			families.LifecycleChecks.Stop.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.LifecycleChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.LifecycleChecks.Sync.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLifecycleSuiteStop registers the fuzz target of Lifecycle.Stop on
// lifecycleIdle, with every check of the suite excused.
func FuzzLifecycleSuiteStop(f *testing.F) {
	families.FuzzLifecycleSuite(
		f,
		families.LifecycleChecks.Stop.Callable,
		lifecycleIdle,
		conform.Excuse(
			families.LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLifecycleSuitePut registers the fuzz target of Lifecycle.Put on
// lifecycleIdle, with every check of the suite excused.
func FuzzLifecycleSuitePut(f *testing.F) {
	families.FuzzLifecycleSuite(
		f,
		families.LifecycleChecks.Put.Callable,
		lifecycleIdle,
		conform.Excuse(
			families.LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLifecycleSuiteSync registers the fuzz target of Lifecycle.Sync on
// lifecycleIdle, with every check of the suite excused.
func FuzzLifecycleSuiteSync(f *testing.F) {
	families.FuzzLifecycleSuite(
		f,
		families.LifecycleChecks.Sync.Callable,
		lifecycleIdle,
		conform.Excuse(
			families.LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLifecycleSuite measures each callable of the suite of Lifecycle
// on lifecycleIdle.
func BenchmarkLifecycleSuite(b *testing.B) {
	b.Run("Stop", func(b *testing.B) {
		families.BenchLifecycleSuite(b, families.LifecycleChecks.Stop.Callable, lifecycleIdle)
	})
	b.Run("Put", func(b *testing.B) {
		families.BenchLifecycleSuite(b, families.LifecycleChecks.Put.Callable, lifecycleIdle)
	})
	b.Run("Sync", func(b *testing.B) {
		families.BenchLifecycleSuite(b, families.LifecycleChecks.Sync.Callable, lifecycleIdle)
	})
}

// shardsSubjectEmbedded is families.Shards under the name of the field
// that shardsSubject embeds.
type shardsSubjectEmbedded = families.Shards

// shardsSubject is a value of Shards that only TestShardsSuite builds.
type shardsSubject struct{ shardsSubjectEmbedded }

// TestShardsSuite checks the suite of Shards and the two forms of its
// application.
func TestShardsSuite(t *testing.T) {
	t.Parallel()

	t.Run("ShardsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.ShardsSuite())
	})

	t.Run("RunShardsSuite", func(t *testing.T) {
		t.Parallel()
		var built *shardsSubject
		families.RunShardsSuite(
			t,
			func() *shardsSubject {
				built = &shardsSubject{}
				return built
			},
			conform.With(families.ShardsCheck{
				ID:        "Shards/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ShardsHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Shards(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.ShardsHandle) families.Shards { return &shardsSubject{} },
			}),
			conform.Excuse(
				families.ShardsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Put.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunShardsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *shardsSubject
		families.RunShardsSuiteFromHandle(
			t,
			func(*families.ShardsHandle) *shardsSubject {
				built = &shardsSubject{}
				return built
			},
			conform.With(families.ShardsCheck{
				ID:        "Shards/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ShardsHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Shards(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.ShardsHandle) families.Shards { return &shardsSubject{} },
			}),
			conform.Excuse(
				families.ShardsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Put.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ShardsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// shardsIdle returns the planted defect of a ceiling of Shards with a
// spend of nothing, a subject that answers each call at once.
func shardsIdle(h *families.ShardsHandle) families.Shards {
	return families.ShardsSuite().Calls[0].Over(h, func() {})
}

// TestShardsSuiteBudgets applies the allocation ceilings of the suite of
// Shards to shardsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestShardsSuiteBudgets(t *testing.T) {
	families.RunShardsBudgets(
		t,
		shardsIdle,
		conform.Ceilings(
			families.ShardsChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.ShardsChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.ShardsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Put.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShardsSuitePut registers the fuzz target of Shards.Put on
// shardsIdle, with every check of the suite excused.
func FuzzShardsSuitePut(f *testing.F) {
	families.FuzzShardsSuite(
		f,
		families.ShardsChecks.Put.Callable,
		shardsIdle,
		conform.Excuse(
			families.ShardsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Put.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShardsSuiteGet registers the fuzz target of Shards.Get on
// shardsIdle, with every check of the suite excused.
func FuzzShardsSuiteGet(f *testing.F) {
	families.FuzzShardsSuite(
		f,
		families.ShardsChecks.Get.Callable,
		shardsIdle,
		conform.Excuse(
			families.ShardsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Put.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ShardsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkShardsSuite measures each callable of the suite of Shards on
// shardsIdle.
func BenchmarkShardsSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		families.BenchShardsSuite(b, families.ShardsChecks.Put.Callable, shardsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchShardsSuite(b, families.ShardsChecks.Get.Callable, shardsIdle)
	})
}

// scopedSubjectEmbedded is families.Scoped under the name of the field
// that scopedSubject embeds.
type scopedSubjectEmbedded = families.Scoped

// scopedSubject is a value of Scoped that only TestScopedSuite builds.
type scopedSubject struct{ scopedSubjectEmbedded }

// TestScopedSuite checks the suite of Scoped and the two forms of its
// application.
func TestScopedSuite(t *testing.T) {
	t.Parallel()

	t.Run("ScopedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.ScopedSuite())
	})

	t.Run("RunScopedSuite", func(t *testing.T) {
		t.Parallel()
		var built *scopedSubject
		families.RunScopedSuite(
			t,
			func() *scopedSubject {
				built = &scopedSubject{}
				return built
			},
			conform.With(families.ScopedCheck{
				ID:        "Scoped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ScopedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Scoped(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.ScopedHandle) families.Scoped { return &scopedSubject{} },
			}),
			conform.Excuse(
				families.ScopedChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Set.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunScopedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *scopedSubject
		families.RunScopedSuiteFromHandle(
			t,
			func(*families.ScopedHandle) *scopedSubject {
				built = &scopedSubject{}
				return built
			},
			conform.With(families.ScopedCheck{
				ID:        "Scoped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.ScopedHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Scoped(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.ScopedHandle) families.Scoped { return &scopedSubject{} },
			}),
			conform.Excuse(
				families.ScopedChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Set.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.ScopedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// scopedIdle returns the planted defect of a ceiling of Scoped with a
// spend of nothing, a subject that answers each call at once.
func scopedIdle(h *families.ScopedHandle) families.Scoped {
	return families.ScopedSuite().Calls[0].Over(h, func() {})
}

// TestScopedSuiteBudgets applies the allocation ceilings of the suite of
// Scoped to scopedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestScopedSuiteBudgets(t *testing.T) {
	families.RunScopedBudgets(
		t,
		scopedIdle,
		conform.Ceilings(
			families.ScopedChecks.Set.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			families.ScopedChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.ScopedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopedSuiteSet registers the fuzz target of Scoped.Set on
// scopedIdle, with every check of the suite excused.
func FuzzScopedSuiteSet(f *testing.F) {
	families.FuzzScopedSuite(
		f,
		families.ScopedChecks.Set.Callable,
		scopedIdle,
		conform.Excuse(
			families.ScopedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopedSuiteGet registers the fuzz target of Scoped.Get on
// scopedIdle, with every check of the suite excused.
func FuzzScopedSuiteGet(f *testing.F) {
	families.FuzzScopedSuite(
		f,
		families.ScopedChecks.Get.Callable,
		scopedIdle,
		conform.Excuse(
			families.ScopedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.ScopedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkScopedSuite measures each callable of the suite of Scoped on
// scopedIdle.
func BenchmarkScopedSuite(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		families.BenchScopedSuite(b, families.ScopedChecks.Set.Callable, scopedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchScopedSuite(b, families.ScopedChecks.Get.Callable, scopedIdle)
	})
}

// plainSubjectEmbedded is families.Plain under the name of the field that
// plainSubject embeds.
type plainSubjectEmbedded = families.Plain

// plainSubject is a value of Plain that only TestPlainSuite builds.
type plainSubject struct{ plainSubjectEmbedded }

// TestPlainSuite checks the suite of Plain and the two forms of its
// application.
func TestPlainSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlainSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.PlainSuite())
	})

	t.Run("RunPlainSuite", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		families.RunPlainSuite(
			t,
			func() *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(families.PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.PlainHandle) families.Plain { return &plainSubject{} },
			}),
			conform.Excuse(
				families.PlainChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Watch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Watch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPlainSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		families.RunPlainSuiteFromHandle(
			t,
			func(*families.PlainHandle) *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(families.PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, families.Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*families.PlainHandle) families.Plain { return &plainSubject{} },
			}),
			conform.Excuse(
				families.PlainChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Watch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.Watch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.PlainChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// plainIdle returns the planted defect of a ceiling of Plain with a spend
// of nothing, a subject that answers each call at once.
func plainIdle(h *families.PlainHandle) families.Plain {
	return families.PlainSuite().Calls[0].Over(h, func() {})
}

// TestPlainSuiteBudgets applies the allocation ceilings of the suite of
// Plain to plainIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPlainSuiteBudgets(t *testing.T) {
	families.RunPlainBudgets(
		t,
		plainIdle,
		conform.Ceilings(families.PlainChecks.Set.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(families.PlainChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			families.PlainChecks.Watch.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(families.PlainChecks.All.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			families.PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteSet registers the fuzz target of Plain.Set on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteSet(f *testing.F) {
	families.FuzzPlainSuite(
		f,
		families.PlainChecks.Set.Callable,
		plainIdle,
		conform.Excuse(
			families.PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteGet registers the fuzz target of Plain.Get on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteGet(f *testing.F) {
	families.FuzzPlainSuite(
		f,
		families.PlainChecks.Get.Callable,
		plainIdle,
		conform.Excuse(
			families.PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteWatch registers the fuzz target of Plain.Watch on
// plainIdle, with every check of the suite excused.
func FuzzPlainSuiteWatch(f *testing.F) {
	families.FuzzPlainSuite(
		f,
		families.PlainChecks.Watch.Callable,
		plainIdle,
		conform.Excuse(
			families.PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteAll registers the fuzz target of Plain.All on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteAll(f *testing.F) {
	families.FuzzPlainSuite(
		f,
		families.PlainChecks.All.Callable,
		plainIdle,
		conform.Excuse(
			families.PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPlainSuite measures each callable of the suite of Plain on
// plainIdle.
func BenchmarkPlainSuite(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		families.BenchPlainSuite(b, families.PlainChecks.Set.Callable, plainIdle)
	})
	b.Run("Get", func(b *testing.B) {
		families.BenchPlainSuite(b, families.PlainChecks.Get.Callable, plainIdle)
	})
	b.Run("Watch", func(b *testing.B) {
		families.BenchPlainSuite(b, families.PlainChecks.Watch.Callable, plainIdle)
	})
	b.Run("All", func(b *testing.B) {
		families.BenchPlainSuite(b, families.PlainChecks.All.Callable, plainIdle)
	})
}

// sessionsSubjectEmbedded is families.Sessions under the name of the field
// that sessionsSubject embeds.
type sessionsSubjectEmbedded = families.Sessions

// sessionsSubject is a value of Sessions that only TestSessionsSuite
// builds.
type sessionsSubject struct{ sessionsSubjectEmbedded }

// TestSessionsSuite checks the suite of Sessions and the two forms of its
// application.
func TestSessionsSuite(t *testing.T) {
	t.Parallel()

	t.Run("SessionsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, families.SessionsSuite())
	})

	t.Run("RunSessionsSuite", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		families.RunSessionsSuite(
			t,
			func() *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(families.SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.SessionsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Sessions(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.SessionsHandle) families.Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				families.SessionsChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SessionsChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SessionsChecks.Open.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSessionsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		families.RunSessionsSuiteFromHandle(
			t,
			func(*families.SessionsHandle) *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(families.SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *families.SessionsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						families.Sessions(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*families.SessionsHandle) families.Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				families.SessionsChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SessionsChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				families.SessionsChecks.Open.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// sessionsIdle returns the planted defect of a ceiling of Sessions with a
// spend of nothing, a subject that answers each call at once.
func sessionsIdle(h *families.SessionsHandle) families.Sessions {
	return families.SessionsSuite().Calls[0].Over(h, func() {})
}

// TestSessionsSuiteBudgets applies the allocation ceilings of the suite of
// Sessions to sessionsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestSessionsSuiteBudgets(t *testing.T) {
	families.RunSessionsBudgets(
		t,
		sessionsIdle,
		conform.Ceilings(
			families.SessionsChecks.Open.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			families.SessionsChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SessionsChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SessionsChecks.Open.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSessionsSuiteOpen registers the fuzz target of Sessions.Open on
// sessionsIdle, with every check of the suite excused.
func FuzzSessionsSuiteOpen(f *testing.F) {
	families.FuzzSessionsSuite(
		f,
		families.SessionsChecks.Open.Callable,
		sessionsIdle,
		conform.Excuse(
			families.SessionsChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SessionsChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			families.SessionsChecks.Open.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSessionsSuite measures each callable of the suite of Sessions
// on sessionsIdle.
func BenchmarkSessionsSuite(b *testing.B) {
	b.Run("Open", func(b *testing.B) {
		families.BenchSessionsSuite(b, families.SessionsChecks.Open.Callable, sessionsIdle)
	})
}
