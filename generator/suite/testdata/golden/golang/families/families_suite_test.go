package families

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// storeSubjectEmbedded is Store under the name of the field that
// storeSubject embeds.
type storeSubjectEmbedded = Store

// storeSubject is a value of Store that only TestStoreSuite builds.
type storeSubject struct{ storeSubjectEmbedded }

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
				StoreChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.MonotonicReads(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Load.ZeroOnMiss(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Load.Cacheable(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Split.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Level.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Level.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Level.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Raise.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Raise.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Inspect.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.GetWithMeta.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.GetWithMeta.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.GetWithMeta.ZeroesOnError(),
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
				StoreChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Get.MonotonicReads(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Load.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Load.ZeroOnMiss(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Load.Cacheable(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Split.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Level.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Level.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Level.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Raise.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Raise.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.Inspect.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.GetWithMeta.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.GetWithMeta.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StoreChecks.GetWithMeta.ZeroesOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// storeIdle returns the planted defect of a ceiling of Store with a spend
// of nothing, a subject that answers each call at once.
func storeIdle(h *StoreHandle) Store {
	return StoreSuite().Calls[0].Over(h, func() {})
}

// TestStoreSuiteBudgets applies the allocation ceilings of the suite of
// Store to storeIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestStoreSuiteBudgets(t *testing.T) {
	RunStoreBudgets(
		t,
		storeIdle,
		conform.Ceilings(StoreChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.Load.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.Split.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.Find.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.Level.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.Raise.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.Inspect.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StoreChecks.GetWithMeta.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuitePut registers the fuzz target of Store.Put on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuitePut(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Put.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteGet registers the fuzz target of Store.Get on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuiteGet(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Get.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteLoad registers the fuzz target of Store.Load on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuiteLoad(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Load.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteSplit registers the fuzz target of Store.Split on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteSplit(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Split.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteFind registers the fuzz target of Store.Find on storeIdle,
// with every check of the suite excused.
func FuzzStoreSuiteFind(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Find.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteLevel registers the fuzz target of Store.Level on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteLevel(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Level.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteRaise registers the fuzz target of Store.Raise on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteRaise(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Raise.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteInspect registers the fuzz target of Store.Inspect on
// storeIdle, with every check of the suite excused.
func FuzzStoreSuiteInspect(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.Inspect.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStoreSuiteGetWithMeta registers the fuzz target of Store.GetWithMeta
// on storeIdle, with every check of the suite excused.
func FuzzStoreSuiteGetWithMeta(f *testing.F) {
	FuzzStoreSuite(
		f,
		StoreChecks.GetWithMeta.Callable,
		storeIdle,
		conform.Excuse(
			StoreChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Get.MonotonicReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Load.Cacheable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Split.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Level.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Raise.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.Inspect.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StoreChecks.GetWithMeta.ZeroesOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStoreSuite measures each callable of the suite of Store on
// storeIdle.
func BenchmarkStoreSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Put.Callable, storeIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Get.Callable, storeIdle)
	})
	b.Run("Load", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Load.Callable, storeIdle)
	})
	b.Run("Split", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Split.Callable, storeIdle)
	})
	b.Run("Find", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Find.Callable, storeIdle)
	})
	b.Run("Level", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Level.Callable, storeIdle)
	})
	b.Run("Raise", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Raise.Callable, storeIdle)
	})
	b.Run("Inspect", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.Inspect.Callable, storeIdle)
	})
	b.Run("GetWithMeta", func(b *testing.B) {
		BenchStoreSuite(b, StoreChecks.GetWithMeta.Callable, storeIdle)
	})
}

// keyedSubjectEmbedded is Keyed under the name of the field that
// keyedSubject embeds.
type keyedSubjectEmbedded = Keyed

// keyedSubject is a value of Keyed that only TestKeyedSuite builds.
type keyedSubject struct{ keyedSubjectEmbedded }

// TestKeyedSuite checks the suite of Keyed and the two forms of its
// application.
func TestKeyedSuite(t *testing.T) {
	t.Parallel()

	t.Run("KeyedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, KeyedSuite())
	})

	t.Run("RunKeyedSuite", func(t *testing.T) {
		t.Parallel()
		var built *keyedSubject
		RunKeyedSuite(
			t,
			func() *keyedSubject {
				built = &keyedSubject{}
				return built
			},
			conform.With(KeyedCheck{
				ID:        "Keyed/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *KeyedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Keyed(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*KeyedHandle) Keyed { return &keyedSubject{} },
			}),
			conform.Excuse(
				KeyedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Delete.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Delete.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunKeyedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *keyedSubject
		RunKeyedSuiteFromHandle(
			t,
			func(*KeyedHandle) *keyedSubject {
				built = &keyedSubject{}
				return built
			},
			conform.With(KeyedCheck{
				ID:        "Keyed/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *KeyedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Keyed(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*KeyedHandle) Keyed { return &keyedSubject{} },
			}),
			conform.Excuse(
				KeyedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Delete.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedChecks.Delete.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// keyedIdle returns the planted defect of a ceiling of Keyed with a spend
// of nothing, a subject that answers each call at once.
func keyedIdle(h *KeyedHandle) Keyed {
	return KeyedSuite().Calls[0].Over(h, func() {})
}

// TestKeyedSuiteBudgets applies the allocation ceilings of the suite of
// Keyed to keyedIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestKeyedSuiteBudgets(t *testing.T) {
	RunKeyedBudgets(
		t,
		keyedIdle,
		conform.Ceilings(KeyedChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(KeyedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(KeyedChecks.Delete.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedSuitePut registers the fuzz target of Keyed.Put on keyedIdle,
// with every check of the suite excused.
func FuzzKeyedSuitePut(f *testing.F) {
	FuzzKeyedSuite(
		f,
		KeyedChecks.Put.Callable,
		keyedIdle,
		conform.Excuse(
			KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedSuiteGet registers the fuzz target of Keyed.Get on keyedIdle,
// with every check of the suite excused.
func FuzzKeyedSuiteGet(f *testing.F) {
	FuzzKeyedSuite(
		f,
		KeyedChecks.Get.Callable,
		keyedIdle,
		conform.Excuse(
			KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedSuiteDelete registers the fuzz target of Keyed.Delete on
// keyedIdle, with every check of the suite excused.
func FuzzKeyedSuiteDelete(f *testing.F) {
	FuzzKeyedSuite(
		f,
		KeyedChecks.Delete.Callable,
		keyedIdle,
		conform.Excuse(
			KeyedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedChecks.Delete.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkKeyedSuite measures each callable of the suite of Keyed on
// keyedIdle.
func BenchmarkKeyedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchKeyedSuite(b, KeyedChecks.Put.Callable, keyedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchKeyedSuite(b, KeyedChecks.Get.Callable, keyedIdle)
	})
	b.Run("Delete", func(b *testing.B) {
		BenchKeyedSuite(b, KeyedChecks.Delete.Callable, keyedIdle)
	})
}

// stampedSubjectEmbedded is Stamped under the name of the field that
// stampedSubject embeds.
type stampedSubjectEmbedded = Stamped

// stampedSubject is a value of Stamped that only TestStampedSuite builds.
type stampedSubject struct{ stampedSubjectEmbedded }

// TestStampedSuite checks the suite of Stamped and the two forms of its
// application.
func TestStampedSuite(t *testing.T) {
	t.Parallel()

	t.Run("StampedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, StampedSuite())
	})

	t.Run("RunStampedSuite", func(t *testing.T) {
		t.Parallel()
		var built *stampedSubject
		RunStampedSuite(
			t,
			func() *stampedSubject {
				built = &stampedSubject{}
				return built
			},
			conform.With(StampedCheck{
				ID:        "Stamped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *StampedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Stamped(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*StampedHandle) Stamped { return &stampedSubject{} },
			}),
			conform.Excuse(
				StampedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Store.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.WritesFollowReads(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunStampedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *stampedSubject
		RunStampedSuiteFromHandle(
			t,
			func(*StampedHandle) *stampedSubject {
				built = &stampedSubject{}
				return built
			},
			conform.With(StampedCheck{
				ID:        "Stamped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *StampedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Stamped(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*StampedHandle) Stamped { return &stampedSubject{} },
			}),
			conform.Excuse(
				StampedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Store.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				StampedChecks.Get.WritesFollowReads(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// stampedIdle returns the planted defect of a ceiling of Stamped with a
// spend of nothing, a subject that answers each call at once.
func stampedIdle(h *StampedHandle) Stamped {
	return StampedSuite().Calls[0].Over(h, func() {})
}

// TestStampedSuiteBudgets applies the allocation ceilings of the suite of
// Stamped to stampedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestStampedSuiteBudgets(t *testing.T) {
	RunStampedBudgets(
		t,
		stampedIdle,
		conform.Ceilings(StampedChecks.Store.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(StampedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Store.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.WritesFollowReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStampedSuiteStore registers the fuzz target of Stamped.Store on
// stampedIdle, with every check of the suite excused.
func FuzzStampedSuiteStore(f *testing.F) {
	FuzzStampedSuite(
		f,
		StampedChecks.Store.Callable,
		stampedIdle,
		conform.Excuse(
			StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Store.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.WritesFollowReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzStampedSuiteGet registers the fuzz target of Stamped.Get on
// stampedIdle, with every check of the suite excused.
func FuzzStampedSuiteGet(f *testing.F) {
	FuzzStampedSuite(
		f,
		StampedChecks.Get.Callable,
		stampedIdle,
		conform.Excuse(
			StampedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Store.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			StampedChecks.Get.WritesFollowReads(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkStampedSuite measures each callable of the suite of Stamped on
// stampedIdle.
func BenchmarkStampedSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		BenchStampedSuite(b, StampedChecks.Store.Callable, stampedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchStampedSuite(b, StampedChecks.Get.Callable, stampedIdle)
	})
}

// decodedSubjectEmbedded is Decoded under the name of the field that
// decodedSubject embeds.
type decodedSubjectEmbedded = Decoded

// decodedSubject is a value of Decoded that only TestDecodedSuite builds.
type decodedSubject struct{ decodedSubjectEmbedded }

// TestDecodedSuite checks the suite of Decoded and the two forms of its
// application.
func TestDecodedSuite(t *testing.T) {
	t.Parallel()

	t.Run("DecodedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, DecodedSuite())
	})

	t.Run("RunDecodedSuite", func(t *testing.T) {
		t.Parallel()
		var built *decodedSubject
		RunDecodedSuite(
			t,
			func() *decodedSubject {
				built = &decodedSubject{}
				return built
			},
			conform.With(DecodedCheck{
				ID:        "Decoded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *DecodedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Decoded(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*DecodedHandle) Decoded { return &decodedSubject{} },
			}),
			conform.Excuse(
				DecodedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunDecodedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *decodedSubject
		RunDecodedSuiteFromHandle(
			t,
			func(*DecodedHandle) *decodedSubject {
				built = &decodedSubject{}
				return built
			},
			conform.With(DecodedCheck{
				ID:        "Decoded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *DecodedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Decoded(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*DecodedHandle) Decoded { return &decodedSubject{} },
			}),
			conform.Excuse(
				DecodedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				DecodedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// decodedIdle returns the planted defect of a ceiling of Decoded with a
// spend of nothing, a subject that answers each call at once.
func decodedIdle(h *DecodedHandle) Decoded {
	return DecodedSuite().Calls[0].Over(h, func() {})
}

// TestDecodedSuiteBudgets applies the allocation ceilings of the suite of
// Decoded to decodedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestDecodedSuiteBudgets(t *testing.T) {
	RunDecodedBudgets(
		t,
		decodedIdle,
		conform.Ceilings(DecodedChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(DecodedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			DecodedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDecodedSuitePut registers the fuzz target of Decoded.Put on
// decodedIdle, with every check of the suite excused.
func FuzzDecodedSuitePut(f *testing.F) {
	FuzzDecodedSuite(
		f,
		DecodedChecks.Put.Callable,
		decodedIdle,
		conform.Excuse(
			DecodedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDecodedSuiteGet registers the fuzz target of Decoded.Get on
// decodedIdle, with every check of the suite excused.
func FuzzDecodedSuiteGet(f *testing.F) {
	FuzzDecodedSuite(
		f,
		DecodedChecks.Get.Callable,
		decodedIdle,
		conform.Excuse(
			DecodedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			DecodedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkDecodedSuite measures each callable of the suite of Decoded on
// decodedIdle.
func BenchmarkDecodedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchDecodedSuite(b, DecodedChecks.Put.Callable, decodedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchDecodedSuite(b, DecodedChecks.Get.Callable, decodedIdle)
	})
}

// expiringSubjectEmbedded is Expiring under the name of the field that
// expiringSubject embeds.
type expiringSubjectEmbedded = Expiring

// expiringSubject is a value of Expiring that only TestExpiringSuite
// builds.
type expiringSubject struct{ expiringSubjectEmbedded }

// TestExpiringSuite checks the suite of Expiring and the two forms of its
// application.
func TestExpiringSuite(t *testing.T) {
	t.Parallel()

	t.Run("ExpiringSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ExpiringSuite())
	})

	t.Run("RunExpiringSuite", func(t *testing.T) {
		t.Parallel()
		var built *expiringSubject
		RunExpiringSuite(
			t,
			func() *expiringSubject {
				built = &expiringSubject{}
				return built
			},
			conform.With(ExpiringCheck{
				ID:        "Expiring/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ExpiringHandle) {
					got := h.Subject()
					assert.Equal(h, got, Expiring(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ExpiringHandle) Expiring { return &expiringSubject{} },
			}),
			conform.Excuse(
				ExpiringChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.Expires(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunExpiringSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *expiringSubject
		RunExpiringSuiteFromHandle(
			t,
			func(*ExpiringHandle) *expiringSubject {
				built = &expiringSubject{}
				return built
			},
			conform.With(ExpiringCheck{
				ID:        "Expiring/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ExpiringHandle) {
					got := h.Subject()
					assert.Equal(h, got, Expiring(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ExpiringHandle) Expiring { return &expiringSubject{} },
			}),
			conform.Excuse(
				ExpiringChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ExpiringChecks.Read.Expires(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// expiringIdle returns the planted defect of a ceiling of Expiring with a
// spend of nothing, a subject that answers each call at once.
func expiringIdle(h *ExpiringHandle) Expiring {
	return ExpiringSuite().Calls[0].Over(h, func() {})
}

// TestExpiringSuiteBudgets applies the allocation ceilings of the suite of
// Expiring to expiringIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestExpiringSuiteBudgets(t *testing.T) {
	RunExpiringBudgets(
		t,
		expiringIdle,
		conform.Ceilings(ExpiringChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ExpiringChecks.Read.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			ExpiringChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Expires(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzExpiringSuitePut registers the fuzz target of Expiring.Put on
// expiringIdle, with every check of the suite excused.
func FuzzExpiringSuitePut(f *testing.F) {
	FuzzExpiringSuite(
		f,
		ExpiringChecks.Put.Callable,
		expiringIdle,
		conform.Excuse(
			ExpiringChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Expires(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzExpiringSuiteRead registers the fuzz target of Expiring.Read on
// expiringIdle, with every check of the suite excused.
func FuzzExpiringSuiteRead(f *testing.F) {
	FuzzExpiringSuite(
		f,
		ExpiringChecks.Read.Callable,
		expiringIdle,
		conform.Excuse(
			ExpiringChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ExpiringChecks.Read.Expires(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkExpiringSuite measures each callable of the suite of Expiring
// on expiringIdle.
func BenchmarkExpiringSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchExpiringSuite(b, ExpiringChecks.Put.Callable, expiringIdle)
	})
	b.Run("Read", func(b *testing.B) {
		BenchExpiringSuite(b, ExpiringChecks.Read.Callable, expiringIdle)
	})
}

// cachedSubjectEmbedded is Cached under the name of the field that
// cachedSubject embeds.
type cachedSubjectEmbedded = Cached

// cachedSubject is a value of Cached that only TestCachedSuite builds.
type cachedSubject struct{ cachedSubjectEmbedded }

// TestCachedSuite checks the suite of Cached and the two forms of its
// application.
func TestCachedSuite(t *testing.T) {
	t.Parallel()

	t.Run("CachedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, CachedSuite())
	})

	t.Run("RunCachedSuite", func(t *testing.T) {
		t.Parallel()
		var built *cachedSubject
		RunCachedSuite(
			t,
			func() *cachedSubject {
				built = &cachedSubject{}
				return built
			},
			conform.With(CachedCheck{
				ID:        "Cached/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CachedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Cached(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CachedHandle) Cached { return &cachedSubject{} },
			}),
			conform.Excuse(
				CachedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Get.Cache(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCachedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *cachedSubject
		RunCachedSuiteFromHandle(
			t,
			func(*CachedHandle) *cachedSubject {
				built = &cachedSubject{}
				return built
			},
			conform.With(CachedCheck{
				ID:        "Cached/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CachedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Cached(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CachedHandle) Cached { return &cachedSubject{} },
			}),
			conform.Excuse(
				CachedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Get.Cache(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CachedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// cachedIdle returns the planted defect of a ceiling of Cached with a
// spend of nothing, a subject that answers each call at once.
func cachedIdle(h *CachedHandle) Cached {
	return CachedSuite().Calls[0].Over(h, func() {})
}

// TestCachedSuiteBudgets applies the allocation ceilings of the suite of
// Cached to cachedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCachedSuiteBudgets(t *testing.T) {
	RunCachedBudgets(
		t,
		cachedIdle,
		conform.Ceilings(CachedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(CachedChecks.Fetch.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(CachedChecks.Store.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteGet registers the fuzz target of Cached.Get on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteGet(f *testing.F) {
	FuzzCachedSuite(
		f,
		CachedChecks.Get.Callable,
		cachedIdle,
		conform.Excuse(
			CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteFetch registers the fuzz target of Cached.Fetch on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteFetch(f *testing.F) {
	FuzzCachedSuite(
		f,
		CachedChecks.Fetch.Callable,
		cachedIdle,
		conform.Excuse(
			CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCachedSuiteStore registers the fuzz target of Cached.Store on
// cachedIdle, with every check of the suite excused.
func FuzzCachedSuiteStore(f *testing.F) {
	FuzzCachedSuite(
		f,
		CachedChecks.Store.Callable,
		cachedIdle,
		conform.Excuse(
			CachedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Get.Cache(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CachedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCachedSuite measures each callable of the suite of Cached on
// cachedIdle.
func BenchmarkCachedSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		BenchCachedSuite(b, CachedChecks.Get.Callable, cachedIdle)
	})
	b.Run("Fetch", func(b *testing.B) {
		BenchCachedSuite(b, CachedChecks.Fetch.Callable, cachedIdle)
	})
	b.Run("Store", func(b *testing.B) {
		BenchCachedSuite(b, CachedChecks.Store.Callable, cachedIdle)
	})
}

// taggedSubjectEmbedded is Tagged under the name of the field that
// taggedSubject embeds.
type taggedSubjectEmbedded = Tagged

// taggedSubject is a value of Tagged that only TestTaggedSuite builds.
type taggedSubject struct{ taggedSubjectEmbedded }

// TestTaggedSuite checks the suite of Tagged and the two forms of its
// application.
func TestTaggedSuite(t *testing.T) {
	t.Parallel()

	t.Run("TaggedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, TaggedSuite())
	})

	t.Run("RunTaggedSuite", func(t *testing.T) {
		t.Parallel()
		var built *taggedSubject
		RunTaggedSuite(
			t,
			func() *taggedSubject {
				built = &taggedSubject{}
				return built
			},
			conform.With(TaggedCheck{
				ID:        "Tagged/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TaggedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Tagged(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TaggedHandle) Tagged { return &taggedSubject{} },
			}),
			conform.Excuse(
				TaggedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTaggedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *taggedSubject
		RunTaggedSuiteFromHandle(
			t,
			func(*TaggedHandle) *taggedSubject {
				built = &taggedSubject{}
				return built
			},
			conform.With(TaggedCheck{
				ID:        "Tagged/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TaggedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Tagged(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TaggedHandle) Tagged { return &taggedSubject{} },
			}),
			conform.Excuse(
				TaggedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TaggedChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// taggedIdle returns the planted defect of a ceiling of Tagged with a
// spend of nothing, a subject that answers each call at once.
func taggedIdle(h *TaggedHandle) Tagged {
	return TaggedSuite().Calls[0].Over(h, func() {})
}

// TestTaggedSuiteBudgets applies the allocation ceilings of the suite of
// Tagged to taggedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestTaggedSuiteBudgets(t *testing.T) {
	RunTaggedBudgets(
		t,
		taggedIdle,
		conform.Ceilings(TaggedChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TaggedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			TaggedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTaggedSuitePut registers the fuzz target of Tagged.Put on
// taggedIdle, with every check of the suite excused.
func FuzzTaggedSuitePut(f *testing.F) {
	FuzzTaggedSuite(
		f,
		TaggedChecks.Put.Callable,
		taggedIdle,
		conform.Excuse(
			TaggedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTaggedSuiteGet registers the fuzz target of Tagged.Get on
// taggedIdle, with every check of the suite excused.
func FuzzTaggedSuiteGet(f *testing.F) {
	FuzzTaggedSuite(
		f,
		TaggedChecks.Get.Callable,
		taggedIdle,
		conform.Excuse(
			TaggedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TaggedChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkTaggedSuite measures each callable of the suite of Tagged on
// taggedIdle.
func BenchmarkTaggedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchTaggedSuite(b, TaggedChecks.Put.Callable, taggedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchTaggedSuite(b, TaggedChecks.Get.Callable, taggedIdle)
	})
}

// snapshotsSubjectEmbedded is Snapshots under the name of the field that
// snapshotsSubject embeds.
type snapshotsSubjectEmbedded = Snapshots

// snapshotsSubject is a value of Snapshots that only TestSnapshotsSuite
// builds.
type snapshotsSubject struct{ snapshotsSubjectEmbedded }

// TestSnapshotsSuite checks the suite of Snapshots and the two forms of
// its application.
func TestSnapshotsSuite(t *testing.T) {
	t.Parallel()

	t.Run("SnapshotsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, SnapshotsSuite())
	})

	t.Run("RunSnapshotsSuite", func(t *testing.T) {
		t.Parallel()
		var built *snapshotsSubject
		RunSnapshotsSuite(
			t,
			func() *snapshotsSubject {
				built = &snapshotsSubject{}
				return built
			},
			conform.With(SnapshotsCheck{
				ID:        "Snapshots/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *SnapshotsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Snapshots(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*SnapshotsHandle) Snapshots { return &snapshotsSubject{} },
			}),
			conform.Excuse(
				SnapshotsChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.PointInTime(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSnapshotsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *snapshotsSubject
		RunSnapshotsSuiteFromHandle(
			t,
			func(*SnapshotsHandle) *snapshotsSubject {
				built = &snapshotsSubject{}
				return built
			},
			conform.With(SnapshotsCheck{
				ID:        "Snapshots/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *SnapshotsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Snapshots(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*SnapshotsHandle) Snapshots { return &snapshotsSubject{} },
			}),
			conform.Excuse(
				SnapshotsChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SnapshotsChecks.Get.PointInTime(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// snapshotsIdle returns the planted defect of a ceiling of Snapshots with
// a spend of nothing, a subject that answers each call at once.
func snapshotsIdle(h *SnapshotsHandle) Snapshots {
	return SnapshotsSuite().Calls[0].Over(h, func() {})
}

// TestSnapshotsSuiteBudgets applies the allocation ceilings of the suite
// of Snapshots to snapshotsIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestSnapshotsSuiteBudgets(t *testing.T) {
	RunSnapshotsBudgets(
		t,
		snapshotsIdle,
		conform.Ceilings(SnapshotsChecks.Store.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(SnapshotsChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			SnapshotsChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.PointInTime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSnapshotsSuiteStore registers the fuzz target of Snapshots.Store on
// snapshotsIdle, with every check of the suite excused.
func FuzzSnapshotsSuiteStore(f *testing.F) {
	FuzzSnapshotsSuite(
		f,
		SnapshotsChecks.Store.Callable,
		snapshotsIdle,
		conform.Excuse(
			SnapshotsChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.PointInTime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSnapshotsSuiteGet registers the fuzz target of Snapshots.Get on
// snapshotsIdle, with every check of the suite excused.
func FuzzSnapshotsSuiteGet(f *testing.F) {
	FuzzSnapshotsSuite(
		f,
		SnapshotsChecks.Get.Callable,
		snapshotsIdle,
		conform.Excuse(
			SnapshotsChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SnapshotsChecks.Get.PointInTime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSnapshotsSuite measures each callable of the suite of Snapshots
// on snapshotsIdle.
func BenchmarkSnapshotsSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		BenchSnapshotsSuite(b, SnapshotsChecks.Store.Callable, snapshotsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchSnapshotsSuite(b, SnapshotsChecks.Get.Callable, snapshotsIdle)
	})
}

// bagSubjectEmbedded is Bag under the name of the field that bagSubject
// embeds.
type bagSubjectEmbedded = Bag

// bagSubject is a value of Bag that only TestBagSuite builds.
type bagSubject struct{ bagSubjectEmbedded }

// TestBagSuite checks the suite of Bag and the two forms of its
// application.
func TestBagSuite(t *testing.T) {
	t.Parallel()

	t.Run("BagSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, BagSuite())
	})

	t.Run("RunBagSuite", func(t *testing.T) {
		t.Parallel()
		var built *bagSubject
		RunBagSuite(
			t,
			func() *bagSubject {
				built = &bagSubject{}
				return built
			},
			conform.With(BagCheck{
				ID:        "Bag/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *BagHandle) {
					got := h.Subject()
					assert.Equal(h, got, Bag(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*BagHandle) Bag { return &bagSubject{} },
			}),
			conform.Excuse(
				BagChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.StableOrder(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBagSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *bagSubject
		RunBagSuiteFromHandle(
			t,
			func(*BagHandle) *bagSubject {
				built = &bagSubject{}
				return built
			},
			conform.With(BagCheck{
				ID:        "Bag/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *BagHandle) {
					got := h.Subject()
					assert.Equal(h, got, Bag(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*BagHandle) Bag { return &bagSubject{} },
			}),
			conform.Excuse(
				BagChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.NoDuplicates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BagChecks.Items.StableOrder(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// bagIdle returns the planted defect of a ceiling of Bag with a spend of
// nothing, a subject that answers each call at once.
func bagIdle(h *BagHandle) Bag {
	return BagSuite().Calls[0].Over(h, func() {})
}

// TestBagSuiteBudgets applies the allocation ceilings of the suite of Bag
// to bagIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestBagSuiteBudgets(t *testing.T) {
	RunBagBudgets(
		t,
		bagIdle,
		conform.Ceilings(BagChecks.Add.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BagChecks.Items.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			BagChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.StableOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBagSuiteAdd registers the fuzz target of Bag.Add on bagIdle, with
// every check of the suite excused.
func FuzzBagSuiteAdd(f *testing.F) {
	FuzzBagSuite(
		f,
		BagChecks.Add.Callable,
		bagIdle,
		conform.Excuse(
			BagChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.StableOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBagSuiteItems registers the fuzz target of Bag.Items on bagIdle,
// with every check of the suite excused.
func FuzzBagSuiteItems(f *testing.F) {
	FuzzBagSuite(
		f,
		BagChecks.Items.Callable,
		bagIdle,
		conform.Excuse(
			BagChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.NoDuplicates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BagChecks.Items.StableOrder(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBagSuite measures each callable of the suite of Bag on bagIdle.
func BenchmarkBagSuite(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		BenchBagSuite(b, BagChecks.Add.Callable, bagIdle)
	})
	b.Run("Items", func(b *testing.B) {
		BenchBagSuite(b, BagChecks.Items.Callable, bagIdle)
	})
}

// updaterSubjectEmbedded is Updater under the name of the field that
// updaterSubject embeds.
type updaterSubjectEmbedded = Updater

// updaterSubject is a value of Updater that only TestUpdaterSuite builds.
type updaterSubject struct{ updaterSubjectEmbedded }

// TestUpdaterSuite checks the suite of Updater and the two forms of its
// application.
func TestUpdaterSuite(t *testing.T) {
	t.Parallel()

	t.Run("UpdaterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, UpdaterSuite())
	})

	t.Run("RunUpdaterSuite", func(t *testing.T) {
		t.Parallel()
		var built *updaterSubject
		RunUpdaterSuite(
			t,
			func() *updaterSubject {
				built = &updaterSubject{}
				return built
			},
			conform.With(UpdaterCheck{
				ID:        "Updater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *UpdaterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Updater(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*UpdaterHandle) Updater { return &updaterSubject{} },
			}),
			conform.Excuse(
				UpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunUpdaterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *updaterSubject
		RunUpdaterSuiteFromHandle(
			t,
			func(*UpdaterHandle) *updaterSubject {
				built = &updaterSubject{}
				return built
			},
			conform.With(UpdaterCheck{
				ID:        "Updater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *UpdaterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Updater(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*UpdaterHandle) Updater { return &updaterSubject{} },
			}),
			conform.Excuse(
				UpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				UpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// updaterIdle returns the planted defect of a ceiling of Updater with a
// spend of nothing, a subject that answers each call at once.
func updaterIdle(h *UpdaterHandle) Updater {
	return UpdaterSuite().Calls[0].Over(h, func() {})
}

// TestUpdaterSuiteBudgets applies the allocation ceilings of the suite of
// Updater to updaterIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestUpdaterSuiteBudgets(t *testing.T) {
	RunUpdaterBudgets(
		t,
		updaterIdle,
		conform.Ceilings(UpdaterChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(UpdaterChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			UpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzUpdaterSuitePut registers the fuzz target of Updater.Put on
// updaterIdle, with every check of the suite excused.
func FuzzUpdaterSuitePut(f *testing.F) {
	FuzzUpdaterSuite(
		f,
		UpdaterChecks.Put.Callable,
		updaterIdle,
		conform.Excuse(
			UpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzUpdaterSuiteGet registers the fuzz target of Updater.Get on
// updaterIdle, with every check of the suite excused.
func FuzzUpdaterSuiteGet(f *testing.F) {
	FuzzUpdaterSuite(
		f,
		UpdaterChecks.Get.Callable,
		updaterIdle,
		conform.Excuse(
			UpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			UpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkUpdaterSuite measures each callable of the suite of Updater on
// updaterIdle.
func BenchmarkUpdaterSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchUpdaterSuite(b, UpdaterChecks.Put.Callable, updaterIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchUpdaterSuite(b, UpdaterChecks.Get.Callable, updaterIdle)
	})
}

// keyedUpdaterSubjectEmbedded is KeyedUpdater under the name of the field
// that keyedUpdaterSubject embeds.
type keyedUpdaterSubjectEmbedded = KeyedUpdater

// keyedUpdaterSubject is a value of KeyedUpdater that only
// TestKeyedUpdaterSuite builds.
type keyedUpdaterSubject struct{ keyedUpdaterSubjectEmbedded }

// TestKeyedUpdaterSuite checks the suite of KeyedUpdater and the two forms
// of its application.
func TestKeyedUpdaterSuite(t *testing.T) {
	t.Parallel()

	t.Run("KeyedUpdaterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, KeyedUpdaterSuite())
	})

	t.Run("RunKeyedUpdaterSuite", func(t *testing.T) {
		t.Parallel()
		var built *keyedUpdaterSubject
		RunKeyedUpdaterSuite(
			t,
			func() *keyedUpdaterSubject {
				built = &keyedUpdaterSubject{}
				return built
			},
			conform.With(KeyedUpdaterCheck{
				ID:        "KeyedUpdater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *KeyedUpdaterHandle) {
					got := h.Subject()
					assert.Equal(h, got, KeyedUpdater(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*KeyedUpdaterHandle) KeyedUpdater { return &keyedUpdaterSubject{} },
			}),
			conform.Excuse(
				KeyedUpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunKeyedUpdaterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *keyedUpdaterSubject
		RunKeyedUpdaterSuiteFromHandle(
			t,
			func(*KeyedUpdaterHandle) *keyedUpdaterSubject {
				built = &keyedUpdaterSubject{}
				return built
			},
			conform.With(KeyedUpdaterCheck{
				ID:        "KeyedUpdater/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *KeyedUpdaterHandle) {
					got := h.Subject()
					assert.Equal(h, got, KeyedUpdater(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*KeyedUpdaterHandle) KeyedUpdater { return &keyedUpdaterSubject{} },
			}),
			conform.Excuse(
				KeyedUpdaterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Put.Updater(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				KeyedUpdaterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// keyedUpdaterIdle returns the planted defect of a ceiling of KeyedUpdater
// with a spend of nothing, a subject that answers each call at once.
func keyedUpdaterIdle(h *KeyedUpdaterHandle) KeyedUpdater {
	return KeyedUpdaterSuite().Calls[0].Over(h, func() {})
}

// TestKeyedUpdaterSuiteBudgets applies the allocation ceilings of the
// suite of KeyedUpdater to keyedUpdaterIdle with every check of the suite
// excused, and a ceiling of 0 to each callable without a budget whose
// calls build no input of their own.
func TestKeyedUpdaterSuiteBudgets(t *testing.T) {
	RunKeyedUpdaterBudgets(
		t,
		keyedUpdaterIdle,
		conform.Ceilings(KeyedUpdaterChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(KeyedUpdaterChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			KeyedUpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedUpdaterSuitePut registers the fuzz target of KeyedUpdater.Put
// on keyedUpdaterIdle, with every check of the suite excused.
func FuzzKeyedUpdaterSuitePut(f *testing.F) {
	FuzzKeyedUpdaterSuite(
		f,
		KeyedUpdaterChecks.Put.Callable,
		keyedUpdaterIdle,
		conform.Excuse(
			KeyedUpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeyedUpdaterSuiteGet registers the fuzz target of KeyedUpdater.Get
// on keyedUpdaterIdle, with every check of the suite excused.
func FuzzKeyedUpdaterSuiteGet(f *testing.F) {
	FuzzKeyedUpdaterSuite(
		f,
		KeyedUpdaterChecks.Get.Callable,
		keyedUpdaterIdle,
		conform.Excuse(
			KeyedUpdaterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Put.Updater(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			KeyedUpdaterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkKeyedUpdaterSuite measures each callable of the suite of
// KeyedUpdater on keyedUpdaterIdle.
func BenchmarkKeyedUpdaterSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchKeyedUpdaterSuite(b, KeyedUpdaterChecks.Put.Callable, keyedUpdaterIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchKeyedUpdaterSuite(b, KeyedUpdaterChecks.Get.Callable, keyedUpdaterIdle)
	})
}

// codecSubjectEmbedded is Codec under the name of the field that
// codecSubject embeds.
type codecSubjectEmbedded = Codec

// codecSubject is a value of Codec that only TestCodecSuite builds.
type codecSubject struct{ codecSubjectEmbedded }

// TestCodecSuite checks the suite of Codec and the two forms of its
// application.
func TestCodecSuite(t *testing.T) {
	t.Parallel()

	t.Run("CodecSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, CodecSuite())
	})

	t.Run("RunCodecSuite", func(t *testing.T) {
		t.Parallel()
		var built *codecSubject
		RunCodecSuite(
			t,
			func() *codecSubject {
				built = &codecSubject{}
				return built
			},
			conform.With(CodecCheck{
				ID:        "Codec/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CodecHandle) {
					got := h.Subject()
					assert.Equal(h, got, Codec(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CodecHandle) Codec { return &codecSubject{} },
			}),
			conform.Excuse(
				CodecChecks.Encode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Encode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Encode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Encode.Codec(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Decode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Decode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Decode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCodecSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *codecSubject
		RunCodecSuiteFromHandle(
			t,
			func(*CodecHandle) *codecSubject {
				built = &codecSubject{}
				return built
			},
			conform.With(CodecCheck{
				ID:        "Codec/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CodecHandle) {
					got := h.Subject()
					assert.Equal(h, got, Codec(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CodecHandle) Codec { return &codecSubject{} },
			}),
			conform.Excuse(
				CodecChecks.Encode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Encode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Encode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Encode.Codec(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Decode.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Decode.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				CodecChecks.Decode.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// codecIdle returns the planted defect of a ceiling of Codec with a spend
// of nothing, a subject that answers each call at once.
func codecIdle(h *CodecHandle) Codec {
	return CodecSuite().Calls[0].Over(h, func() {})
}

// TestCodecSuiteBudgets applies the allocation ceilings of the suite of
// Codec to codecIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestCodecSuiteBudgets(t *testing.T) {
	RunCodecBudgets(
		t,
		codecIdle,
		conform.Ceilings(CodecChecks.Encode.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(CodecChecks.Decode.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			CodecChecks.Encode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.Codec(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCodecSuiteEncode registers the fuzz target of Codec.Encode on
// codecIdle, with every check of the suite excused.
func FuzzCodecSuiteEncode(f *testing.F) {
	FuzzCodecSuite(
		f,
		CodecChecks.Encode.Callable,
		codecIdle,
		conform.Excuse(
			CodecChecks.Encode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.Codec(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCodecSuiteDecode registers the fuzz target of Codec.Decode on
// codecIdle, with every check of the suite excused.
func FuzzCodecSuiteDecode(f *testing.F) {
	FuzzCodecSuite(
		f,
		CodecChecks.Decode.Callable,
		codecIdle,
		conform.Excuse(
			CodecChecks.Encode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Encode.Codec(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			CodecChecks.Decode.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCodecSuite measures each callable of the suite of Codec on
// codecIdle.
func BenchmarkCodecSuite(b *testing.B) {
	b.Run("Encode", func(b *testing.B) {
		BenchCodecSuite(b, CodecChecks.Encode.Callable, codecIdle)
	})
	b.Run("Decode", func(b *testing.B) {
		BenchCodecSuite(b, CodecChecks.Decode.Callable, codecIdle)
	})
}

// pagerSubjectEmbedded is Pager under the name of the field that
// pagerSubject embeds.
type pagerSubjectEmbedded = Pager

// pagerSubject is a value of Pager that only TestPagerSuite builds.
type pagerSubject struct{ pagerSubjectEmbedded }

// TestPagerSuite checks the suite of Pager and the two forms of its
// application.
func TestPagerSuite(t *testing.T) {
	t.Parallel()

	t.Run("PagerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PagerSuite())
	})

	t.Run("RunPagerSuite", func(t *testing.T) {
		t.Parallel()
		var built *pagerSubject
		RunPagerSuite(
			t,
			func() *pagerSubject {
				built = &pagerSubject{}
				return built
			},
			conform.With(PagerCheck{
				ID:        "Pager/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PagerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Pager(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PagerHandle) Pager { return &pagerSubject{} },
			}),
			conform.Excuse(
				PagerChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PagerChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PagerChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PagerChecks.Fetch.Pagination(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPagerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *pagerSubject
		RunPagerSuiteFromHandle(
			t,
			func(*PagerHandle) *pagerSubject {
				built = &pagerSubject{}
				return built
			},
			conform.With(PagerCheck{
				ID:        "Pager/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PagerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Pager(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PagerHandle) Pager { return &pagerSubject{} },
			}),
			conform.Excuse(
				PagerChecks.Fetch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PagerChecks.Fetch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PagerChecks.Fetch.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PagerChecks.Fetch.Pagination(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// pagerIdle returns the planted defect of a ceiling of Pager with a spend
// of nothing, a subject that answers each call at once.
func pagerIdle(h *PagerHandle) Pager {
	return PagerSuite().Calls[0].Over(h, func() {})
}

// TestPagerSuiteBudgets applies the allocation ceilings of the suite of
// Pager to pagerIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPagerSuiteBudgets(t *testing.T) {
	RunPagerBudgets(
		t,
		pagerIdle,
		conform.Ceilings(PagerChecks.Fetch.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			PagerChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PagerChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PagerChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PagerChecks.Fetch.Pagination(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPagerSuiteFetch registers the fuzz target of Pager.Fetch on
// pagerIdle, with every check of the suite excused.
func FuzzPagerSuiteFetch(f *testing.F) {
	FuzzPagerSuite(
		f,
		PagerChecks.Fetch.Callable,
		pagerIdle,
		conform.Excuse(
			PagerChecks.Fetch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PagerChecks.Fetch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PagerChecks.Fetch.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PagerChecks.Fetch.Pagination(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPagerSuite measures each callable of the suite of Pager on
// pagerIdle.
func BenchmarkPagerSuite(b *testing.B) {
	b.Run("Fetch", func(b *testing.B) {
		BenchPagerSuite(b, PagerChecks.Fetch.Callable, pagerIdle)
	})
}

// logSubjectEmbedded is Log under the name of the field that logSubject
// embeds.
type logSubjectEmbedded = Log

// logSubject is a value of Log that only TestLogSuite builds.
type logSubject struct{ logSubjectEmbedded }

// TestLogSuite checks the suite of Log and the two forms of its
// application.
func TestLogSuite(t *testing.T) {
	t.Parallel()

	t.Run("LogSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, LogSuite())
	})

	t.Run("RunLogSuite", func(t *testing.T) {
		t.Parallel()
		var built *logSubject
		RunLogSuite(
			t,
			func() *logSubject {
				built = &logSubject{}
				return built
			},
			conform.With(LogCheck{
				ID:        "Log/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LogHandle) {
					got := h.Subject()
					assert.Equal(h, got, Log(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LogHandle) Log { return &logSubject{} },
			}),
			conform.Excuse(
				LogChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Append.ChainCausal(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Append.Chain(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Replay.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Replay.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Replay.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLogSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *logSubject
		RunLogSuiteFromHandle(
			t,
			func(*LogHandle) *logSubject {
				built = &logSubject{}
				return built
			},
			conform.With(LogCheck{
				ID:        "Log/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LogHandle) {
					got := h.Subject()
					assert.Equal(h, got, Log(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LogHandle) Log { return &logSubject{} },
			}),
			conform.Excuse(
				LogChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Append.ChainCausal(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Append.Chain(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Replay.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Replay.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LogChecks.Replay.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// logIdle returns the planted defect of a ceiling of Log with a spend of
// nothing, a subject that answers each call at once.
func logIdle(h *LogHandle) Log {
	return LogSuite().Calls[0].Over(h, func() {})
}

// TestLogSuiteBudgets applies the allocation ceilings of the suite of Log
// to logIdle with every check of the suite excused, and a ceiling of 0 to
// each callable without a budget whose calls build no input of their own.
func TestLogSuiteBudgets(t *testing.T) {
	RunLogBudgets(
		t,
		logIdle,
		conform.Ceilings(LogChecks.Append.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LogChecks.Replay.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.ChainCausal(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLogSuiteAppend registers the fuzz target of Log.Append on logIdle,
// with every check of the suite excused.
func FuzzLogSuiteAppend(f *testing.F) {
	FuzzLogSuite(
		f,
		LogChecks.Append.Callable,
		logIdle,
		conform.Excuse(
			LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.ChainCausal(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLogSuiteReplay registers the fuzz target of Log.Replay on logIdle,
// with every check of the suite excused.
func FuzzLogSuiteReplay(f *testing.F) {
	FuzzLogSuite(
		f,
		LogChecks.Replay.Callable,
		logIdle,
		conform.Excuse(
			LogChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.ChainCausal(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LogChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLogSuite measures each callable of the suite of Log on logIdle.
func BenchmarkLogSuite(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		BenchLogSuite(b, LogChecks.Append.Callable, logIdle)
	})
	b.Run("Replay", func(b *testing.B) {
		BenchLogSuite(b, LogChecks.Replay.Callable, logIdle)
	})
}

// openerSubjectEmbedded is Opener under the name of the field that
// openerSubject embeds.
type openerSubjectEmbedded = Opener

// openerSubject is a value of Opener that only TestOpenerSuite builds.
type openerSubject struct{ openerSubjectEmbedded }

// TestOpenerSuite checks the suite of Opener and the two forms of its
// application.
func TestOpenerSuite(t *testing.T) {
	t.Parallel()

	t.Run("OpenerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, OpenerSuite())
	})

	t.Run("RunOpenerSuite", func(t *testing.T) {
		t.Parallel()
		var built *openerSubject
		RunOpenerSuite(
			t,
			func() *openerSubject {
				built = &openerSubject{}
				return built
			},
			conform.With(OpenerCheck{
				ID:        "Opener/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *OpenerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Opener(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*OpenerHandle) Opener { return &openerSubject{} },
			}),
			conform.Excuse(
				OpenerChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				OpenerChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				OpenerChecks.Open.CursorOpen(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunOpenerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *openerSubject
		RunOpenerSuiteFromHandle(
			t,
			func(*OpenerHandle) *openerSubject {
				built = &openerSubject{}
				return built
			},
			conform.With(OpenerCheck{
				ID:        "Opener/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *OpenerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Opener(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*OpenerHandle) Opener { return &openerSubject{} },
			}),
			conform.Excuse(
				OpenerChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				OpenerChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				OpenerChecks.Open.CursorOpen(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// openerIdle returns the planted defect of a ceiling of Opener with a
// spend of nothing, a subject that answers each call at once.
func openerIdle(h *OpenerHandle) Opener {
	return OpenerSuite().Calls[0].Over(h, func() {})
}

// TestOpenerSuiteBudgets applies the allocation ceilings of the suite of
// Opener to openerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestOpenerSuiteBudgets(t *testing.T) {
	RunOpenerBudgets(
		t,
		openerIdle,
		conform.Ceilings(OpenerChecks.Open.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			OpenerChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			OpenerChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			OpenerChecks.Open.CursorOpen(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzOpenerSuiteOpen registers the fuzz target of Opener.Open on
// openerIdle, with every check of the suite excused.
func FuzzOpenerSuiteOpen(f *testing.F) {
	FuzzOpenerSuite(
		f,
		OpenerChecks.Open.Callable,
		openerIdle,
		conform.Excuse(
			OpenerChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			OpenerChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			OpenerChecks.Open.CursorOpen(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkOpenerSuite measures each callable of the suite of Opener on
// openerIdle.
func BenchmarkOpenerSuite(b *testing.B) {
	b.Run("Open", func(b *testing.B) {
		BenchOpenerSuite(b, OpenerChecks.Open.Callable, openerIdle)
	})
}

// ledgerSubjectEmbedded is Ledger under the name of the field that
// ledgerSubject embeds.
type ledgerSubjectEmbedded = Ledger

// ledgerSubject is a value of Ledger that only TestLedgerSuite builds.
type ledgerSubject struct{ ledgerSubjectEmbedded }

// TestLedgerSuite checks the suite of Ledger and the two forms of its
// application.
func TestLedgerSuite(t *testing.T) {
	t.Parallel()

	t.Run("LedgerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, LedgerSuite())
	})

	t.Run("RunLedgerSuite", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		RunLedgerSuite(
			t,
			func() *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LedgerHandle) Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				LedgerChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Begin.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Begin.Tx(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLedgerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		RunLedgerSuiteFromHandle(
			t,
			func(*LedgerHandle) *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LedgerHandle) Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				LedgerChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Begin.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Begin.Tx(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// ledgerIdle returns the planted defect of a ceiling of Ledger with a
// spend of nothing, a subject that answers each call at once.
func ledgerIdle(h *LedgerHandle) Ledger {
	return LedgerSuite().Calls[0].Over(h, func() {})
}

// TestLedgerSuiteBudgets applies the allocation ceilings of the suite of
// Ledger to ledgerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLedgerSuiteBudgets(t *testing.T) {
	RunLedgerBudgets(
		t,
		ledgerIdle,
		conform.Ceilings(LedgerChecks.Begin.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LedgerChecks.Commit.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LedgerChecks.Rollback.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteBegin registers the fuzz target of Ledger.Begin on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteBegin(f *testing.F) {
	FuzzLedgerSuite(
		f,
		LedgerChecks.Begin.Callable,
		ledgerIdle,
		conform.Excuse(
			LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteCommit registers the fuzz target of Ledger.Commit on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteCommit(f *testing.F) {
	FuzzLedgerSuite(
		f,
		LedgerChecks.Commit.Callable,
		ledgerIdle,
		conform.Excuse(
			LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteRollback registers the fuzz target of Ledger.Rollback on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteRollback(f *testing.F) {
	FuzzLedgerSuite(
		f,
		LedgerChecks.Rollback.Callable,
		ledgerIdle,
		conform.Excuse(
			LedgerChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Begin.Tx(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLedgerSuite measures each callable of the suite of Ledger on
// ledgerIdle.
func BenchmarkLedgerSuite(b *testing.B) {
	b.Run("Begin", func(b *testing.B) {
		BenchLedgerSuite(b, LedgerChecks.Begin.Callable, ledgerIdle)
	})
	b.Run("Commit", func(b *testing.B) {
		BenchLedgerSuite(b, LedgerChecks.Commit.Callable, ledgerIdle)
	})
	b.Run("Rollback", func(b *testing.B) {
		BenchLedgerSuite(b, LedgerChecks.Rollback.Callable, ledgerIdle)
	})
}

// guardedSubjectEmbedded is Guarded under the name of the field that
// guardedSubject embeds.
type guardedSubjectEmbedded = Guarded

// guardedSubject is a value of Guarded that only TestGuardedSuite builds.
type guardedSubject struct{ guardedSubjectEmbedded }

// TestGuardedSuite checks the suite of Guarded and the two forms of its
// application.
func TestGuardedSuite(t *testing.T) {
	t.Parallel()

	t.Run("GuardedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, GuardedSuite())
	})

	t.Run("RunGuardedSuite", func(t *testing.T) {
		t.Parallel()
		var built *guardedSubject
		RunGuardedSuite(
			t,
			func() *guardedSubject {
				built = &guardedSubject{}
				return built
			},
			conform.With(GuardedCheck{
				ID:        "Guarded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *GuardedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Guarded(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*GuardedHandle) Guarded { return &guardedSubject{} },
			}),
			conform.Excuse(
				GuardedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Put.IfMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Match.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Match.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunGuardedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *guardedSubject
		RunGuardedSuiteFromHandle(
			t,
			func(*GuardedHandle) *guardedSubject {
				built = &guardedSubject{}
				return built
			},
			conform.With(GuardedCheck{
				ID:        "Guarded/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *GuardedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Guarded(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*GuardedHandle) Guarded { return &guardedSubject{} },
			}),
			conform.Excuse(
				GuardedChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Put.IfMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Match.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				GuardedChecks.Match.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// guardedIdle returns the planted defect of a ceiling of Guarded with a
// spend of nothing, a subject that answers each call at once.
func guardedIdle(h *GuardedHandle) Guarded {
	return GuardedSuite().Calls[0].Over(h, func() {})
}

// TestGuardedSuiteBudgets applies the allocation ceilings of the suite of
// Guarded to guardedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestGuardedSuiteBudgets(t *testing.T) {
	RunGuardedBudgets(
		t,
		guardedIdle,
		conform.Ceilings(GuardedChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(GuardedChecks.Match.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			GuardedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Put.IfMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Match.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Match.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzGuardedSuitePut registers the fuzz target of Guarded.Put on
// guardedIdle, with every check of the suite excused.
func FuzzGuardedSuitePut(f *testing.F) {
	FuzzGuardedSuite(
		f,
		GuardedChecks.Put.Callable,
		guardedIdle,
		conform.Excuse(
			GuardedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Put.IfMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Match.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Match.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzGuardedSuiteMatch registers the fuzz target of Guarded.Match on
// guardedIdle, with every check of the suite excused.
func FuzzGuardedSuiteMatch(f *testing.F) {
	FuzzGuardedSuite(
		f,
		GuardedChecks.Match.Callable,
		guardedIdle,
		conform.Excuse(
			GuardedChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Put.IfMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Match.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			GuardedChecks.Match.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkGuardedSuite measures each callable of the suite of Guarded on
// guardedIdle.
func BenchmarkGuardedSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchGuardedSuite(b, GuardedChecks.Put.Callable, guardedIdle)
	})
	b.Run("Match", func(b *testing.B) {
		BenchGuardedSuite(b, GuardedChecks.Match.Callable, guardedIdle)
	})
}

// scopeSubjectEmbedded is Scope under the name of the field that
// scopeSubject embeds.
type scopeSubjectEmbedded = Scope

// scopeSubject is a value of Scope that only TestScopeSuite builds.
type scopeSubject struct{ scopeSubjectEmbedded }

// TestScopeSuite checks the suite of Scope and the two forms of its
// application.
func TestScopeSuite(t *testing.T) {
	t.Parallel()

	t.Run("ScopeSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ScopeSuite())
	})

	t.Run("RunScopeSuite", func(t *testing.T) {
		t.Parallel()
		var built *scopeSubject
		RunScopeSuite(
			t,
			func() *scopeSubject {
				built = &scopeSubject{}
				return built
			},
			conform.With(ScopeCheck{
				ID:        "Scope/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ScopeHandle) {
					got := h.Subject()
					assert.Equal(h, got, Scope(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ScopeHandle) Scope { return &scopeSubject{} },
			}),
			conform.Excuse(
				ScopeChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Run.Transaction(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunScopeSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *scopeSubject
		RunScopeSuiteFromHandle(
			t,
			func(*ScopeHandle) *scopeSubject {
				built = &scopeSubject{}
				return built
			},
			conform.With(ScopeCheck{
				ID:        "Scope/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ScopeHandle) {
					got := h.Subject()
					assert.Equal(h, got, Scope(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ScopeHandle) Scope { return &scopeSubject{} },
			}),
			conform.Excuse(
				ScopeChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Run.Transaction(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopeChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// scopeIdle returns the planted defect of a ceiling of Scope with a spend
// of nothing, a subject that answers each call at once.
func scopeIdle(h *ScopeHandle) Scope {
	return ScopeSuite().Calls[0].Over(h, func() {})
}

// TestScopeSuiteBudgets applies the allocation ceilings of the suite of
// Scope to scopeIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestScopeSuiteBudgets(t *testing.T) {
	RunScopeBudgets(
		t,
		scopeIdle,
		conform.Ceilings(ScopeChecks.Run.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ScopeChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ScopeChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopeSuiteRun registers the fuzz target of Scope.Run on scopeIdle,
// with every check of the suite excused.
func FuzzScopeSuiteRun(f *testing.F) {
	FuzzScopeSuite(
		f,
		ScopeChecks.Run.Callable,
		scopeIdle,
		conform.Excuse(
			ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopeSuitePut registers the fuzz target of Scope.Put on scopeIdle,
// with every check of the suite excused.
func FuzzScopeSuitePut(f *testing.F) {
	FuzzScopeSuite(
		f,
		ScopeChecks.Put.Callable,
		scopeIdle,
		conform.Excuse(
			ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopeSuiteGet registers the fuzz target of Scope.Get on scopeIdle,
// with every check of the suite excused.
func FuzzScopeSuiteGet(f *testing.F) {
	FuzzScopeSuite(
		f,
		ScopeChecks.Get.Callable,
		scopeIdle,
		conform.Excuse(
			ScopeChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopeChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkScopeSuite measures each callable of the suite of Scope on
// scopeIdle.
func BenchmarkScopeSuite(b *testing.B) {
	b.Run("Run", func(b *testing.B) {
		BenchScopeSuite(b, ScopeChecks.Run.Callable, scopeIdle)
	})
	b.Run("Put", func(b *testing.B) {
		BenchScopeSuite(b, ScopeChecks.Put.Callable, scopeIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchScopeSuite(b, ScopeChecks.Get.Callable, scopeIdle)
	})
}

// eventsSubjectEmbedded is Events under the name of the field that
// eventsSubject embeds.
type eventsSubjectEmbedded = Events

// eventsSubject is a value of Events that only TestEventsSuite builds.
type eventsSubject struct{ eventsSubjectEmbedded }

// TestEventsSuite checks the suite of Events and the two forms of its
// application.
func TestEventsSuite(t *testing.T) {
	t.Parallel()

	t.Run("EventsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, EventsSuite())
	})

	t.Run("RunEventsSuite", func(t *testing.T) {
		t.Parallel()
		var built *eventsSubject
		RunEventsSuite(
			t,
			func() *eventsSubject {
				built = &eventsSubject{}
				return built
			},
			conform.With(EventsCheck{
				ID:        "Events/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *EventsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Events(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*EventsHandle) Events { return &eventsSubject{} },
			}),
			conform.Excuse(
				EventsChecks.Fire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				EventsChecks.Fire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				EventsChecks.Fire.Hooks(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				EventsChecks.OnEvent.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunEventsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *eventsSubject
		RunEventsSuiteFromHandle(
			t,
			func(*EventsHandle) *eventsSubject {
				built = &eventsSubject{}
				return built
			},
			conform.With(EventsCheck{
				ID:        "Events/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *EventsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Events(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*EventsHandle) Events { return &eventsSubject{} },
			}),
			conform.Excuse(
				EventsChecks.Fire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				EventsChecks.Fire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				EventsChecks.Fire.Hooks(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				EventsChecks.OnEvent.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// eventsIdle returns the planted defect of a ceiling of Events with a
// spend of nothing, a subject that answers each call at once.
func eventsIdle(h *EventsHandle) Events {
	return EventsSuite().Calls[0].Over(h, func() {})
}

// TestEventsSuiteBudgets applies the allocation ceilings of the suite of
// Events to eventsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestEventsSuiteBudgets(t *testing.T) {
	RunEventsBudgets(
		t,
		eventsIdle,
		conform.Ceilings(EventsChecks.Fire.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(EventsChecks.OnEvent.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			EventsChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.Fire.Hooks(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.OnEvent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzEventsSuiteFire registers the fuzz target of Events.Fire on
// eventsIdle, with every check of the suite excused.
func FuzzEventsSuiteFire(f *testing.F) {
	FuzzEventsSuite(
		f,
		EventsChecks.Fire.Callable,
		eventsIdle,
		conform.Excuse(
			EventsChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.Fire.Hooks(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.OnEvent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzEventsSuiteOnEvent registers the fuzz target of Events.OnEvent on
// eventsIdle, with every check of the suite excused.
func FuzzEventsSuiteOnEvent(f *testing.F) {
	FuzzEventsSuite(
		f,
		EventsChecks.OnEvent.Callable,
		eventsIdle,
		conform.Excuse(
			EventsChecks.Fire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.Fire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.Fire.Hooks(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			EventsChecks.OnEvent.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkEventsSuite measures each callable of the suite of Events on
// eventsIdle.
func BenchmarkEventsSuite(b *testing.B) {
	b.Run("Fire", func(b *testing.B) {
		BenchEventsSuite(b, EventsChecks.Fire.Callable, eventsIdle)
	})
	b.Run("OnEvent", func(b *testing.B) {
		BenchEventsSuite(b, EventsChecks.OnEvent.Callable, eventsIdle)
	})
}

// timersSubjectEmbedded is Timers under the name of the field that
// timersSubject embeds.
type timersSubjectEmbedded = Timers

// timersSubject is a value of Timers that only TestTimersSuite builds.
type timersSubject struct{ timersSubjectEmbedded }

// TestTimersSuite checks the suite of Timers and the two forms of its
// application.
func TestTimersSuite(t *testing.T) {
	t.Parallel()

	t.Run("TimersSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, TimersSuite())
	})

	t.Run("RunTimersSuite", func(t *testing.T) {
		t.Parallel()
		var built *timersSubject
		RunTimersSuite(
			t,
			func() *timersSubject {
				built = &timersSubject{}
				return built
			},
			conform.With(TimersCheck{
				ID:        "Timers/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TimersHandle) {
					got := h.Subject()
					assert.Equal(h, got, Timers(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TimersHandle) Timers { return &timersSubject{} },
			}),
			conform.Excuse(
				TimersChecks.At.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.At.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.At.Scheduled(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Fired.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Fired.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Fired.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Slow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Slow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Slow.Within(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTimersSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *timersSubject
		RunTimersSuiteFromHandle(
			t,
			func(*TimersHandle) *timersSubject {
				built = &timersSubject{}
				return built
			},
			conform.With(TimersCheck{
				ID:        "Timers/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *TimersHandle) {
					got := h.Subject()
					assert.Equal(h, got, Timers(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*TimersHandle) Timers { return &timersSubject{} },
			}),
			conform.Excuse(
				TimersChecks.At.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.At.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.At.Scheduled(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Fired.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Fired.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Fired.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Slow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Slow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				TimersChecks.Slow.Within(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// timersIdle returns the planted defect of a ceiling of Timers with a
// spend of nothing, a subject that answers each call at once.
func timersIdle(h *TimersHandle) Timers {
	return TimersSuite().Calls[0].Over(h, func() {})
}

// TestTimersSuiteBudgets applies the allocation ceilings of the suite of
// Timers to timersIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestTimersSuiteBudgets(t *testing.T) {
	RunTimersBudgets(
		t,
		timersIdle,
		conform.Ceilings(TimersChecks.At.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TimersChecks.Fired.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(TimersChecks.Slow.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTimersSuiteAt registers the fuzz target of Timers.At on timersIdle,
// with every check of the suite excused.
func FuzzTimersSuiteAt(f *testing.F) {
	FuzzTimersSuite(
		f,
		TimersChecks.At.Callable,
		timersIdle,
		conform.Excuse(
			TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTimersSuiteFired registers the fuzz target of Timers.Fired on
// timersIdle, with every check of the suite excused.
func FuzzTimersSuiteFired(f *testing.F) {
	FuzzTimersSuite(
		f,
		TimersChecks.Fired.Callable,
		timersIdle,
		conform.Excuse(
			TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTimersSuiteSlow registers the fuzz target of Timers.Slow on
// timersIdle, with every check of the suite excused.
func FuzzTimersSuiteSlow(f *testing.F) {
	FuzzTimersSuite(
		f,
		TimersChecks.Slow.Callable,
		timersIdle,
		conform.Excuse(
			TimersChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.At.Scheduled(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Fired.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			TimersChecks.Slow.Within(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkTimersSuite measures each callable of the suite of Timers on
// timersIdle.
func BenchmarkTimersSuite(b *testing.B) {
	b.Run("At", func(b *testing.B) {
		BenchTimersSuite(b, TimersChecks.At.Callable, timersIdle)
	})
	b.Run("Fired", func(b *testing.B) {
		BenchTimersSuite(b, TimersChecks.Fired.Callable, timersIdle)
	})
	b.Run("Slow", func(b *testing.B) {
		BenchTimersSuite(b, TimersChecks.Slow.Callable, timersIdle)
	})
}

// electionSubjectEmbedded is Election under the name of the field that
// electionSubject embeds.
type electionSubjectEmbedded = Election

// electionSubject is a value of Election that only TestElectionSuite
// builds.
type electionSubject struct{ electionSubjectEmbedded }

// TestElectionSuite checks the suite of Election and the two forms of its
// application.
func TestElectionSuite(t *testing.T) {
	t.Parallel()

	t.Run("ElectionSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ElectionSuite())
	})

	t.Run("RunElectionSuite", func(t *testing.T) {
		t.Parallel()
		var built *electionSubject
		RunElectionSuite(
			t,
			func() *electionSubject {
				built = &electionSubject{}
				return built
			},
			conform.With(ElectionCheck{
				ID:        "Election/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ElectionHandle) {
					got := h.Subject()
					assert.Equal(h, got, Election(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ElectionHandle) Election { return &electionSubject{} },
			}),
			conform.Excuse(
				ElectionChecks.Campaign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Campaign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Campaign.LeaderElection(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Resign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Resign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.IsLeader.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunElectionSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *electionSubject
		RunElectionSuiteFromHandle(
			t,
			func(*ElectionHandle) *electionSubject {
				built = &electionSubject{}
				return built
			},
			conform.With(ElectionCheck{
				ID:        "Election/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ElectionHandle) {
					got := h.Subject()
					assert.Equal(h, got, Election(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ElectionHandle) Election { return &electionSubject{} },
			}),
			conform.Excuse(
				ElectionChecks.Campaign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Campaign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Campaign.LeaderElection(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Resign.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.Resign.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ElectionChecks.IsLeader.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// electionIdle returns the planted defect of a ceiling of Election with a
// spend of nothing, a subject that answers each call at once.
func electionIdle(h *ElectionHandle) Election {
	return ElectionSuite().Calls[0].Over(h, func() {})
}

// TestElectionSuiteBudgets applies the allocation ceilings of the suite of
// Election to electionIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestElectionSuiteBudgets(t *testing.T) {
	RunElectionBudgets(
		t,
		electionIdle,
		conform.Ceilings(ElectionChecks.Campaign.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ElectionChecks.Resign.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ElectionChecks.IsLeader.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzElectionSuiteCampaign registers the fuzz target of Election.Campaign
// on electionIdle, with every check of the suite excused.
func FuzzElectionSuiteCampaign(f *testing.F) {
	FuzzElectionSuite(
		f,
		ElectionChecks.Campaign.Callable,
		electionIdle,
		conform.Excuse(
			ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzElectionSuiteResign registers the fuzz target of Election.Resign on
// electionIdle, with every check of the suite excused.
func FuzzElectionSuiteResign(f *testing.F) {
	FuzzElectionSuite(
		f,
		ElectionChecks.Resign.Callable,
		electionIdle,
		conform.Excuse(
			ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzElectionSuiteIsLeader registers the fuzz target of Election.IsLeader
// on electionIdle, with every check of the suite excused.
func FuzzElectionSuiteIsLeader(f *testing.F) {
	FuzzElectionSuite(
		f,
		ElectionChecks.IsLeader.Callable,
		electionIdle,
		conform.Excuse(
			ElectionChecks.Campaign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Campaign.LeaderElection(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.Resign.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ElectionChecks.IsLeader.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkElectionSuite measures each callable of the suite of Election
// on electionIdle.
func BenchmarkElectionSuite(b *testing.B) {
	b.Run("Campaign", func(b *testing.B) {
		BenchElectionSuite(b, ElectionChecks.Campaign.Callable, electionIdle)
	})
	b.Run("Resign", func(b *testing.B) {
		BenchElectionSuite(b, ElectionChecks.Resign.Callable, electionIdle)
	})
	b.Run("IsLeader", func(b *testing.B) {
		BenchElectionSuite(b, ElectionChecks.IsLeader.Callable, electionIdle)
	})
}

// poolSubjectEmbedded is Pool under the name of the field that poolSubject
// embeds.
type poolSubjectEmbedded = Pool

// poolSubject is a value of Pool that only TestPoolSuite builds.
type poolSubject struct{ poolSubjectEmbedded }

// TestPoolSuite checks the suite of Pool and the two forms of its
// application.
func TestPoolSuite(t *testing.T) {
	t.Parallel()

	t.Run("PoolSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PoolSuite())
	})

	t.Run("RunPoolSuite", func(t *testing.T) {
		t.Parallel()
		var built *poolSubject
		RunPoolSuite(
			t,
			func() *poolSubject {
				built = &poolSubject{}
				return built
			},
			conform.With(PoolCheck{
				ID:        "Pool/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PoolHandle) {
					got := h.Subject()
					assert.Equal(h, got, Pool(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PoolHandle) Pool { return &poolSubject{} },
			}),
			conform.Excuse(
				PoolChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Get.Pool(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPoolSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *poolSubject
		RunPoolSuiteFromHandle(
			t,
			func(*PoolHandle) *poolSubject {
				built = &poolSubject{}
				return built
			},
			conform.With(PoolCheck{
				ID:        "Pool/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PoolHandle) {
					got := h.Subject()
					assert.Equal(h, got, Pool(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PoolHandle) Pool { return &poolSubject{} },
			}),
			conform.Excuse(
				PoolChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Get.Pool(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PoolChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// poolIdle returns the planted defect of a ceiling of Pool with a spend of
// nothing, a subject that answers each call at once.
func poolIdle(h *PoolHandle) Pool {
	return PoolSuite().Calls[0].Over(h, func() {})
}

// TestPoolSuiteBudgets applies the allocation ceilings of the suite of
// Pool to poolIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestPoolSuiteBudgets(t *testing.T) {
	RunPoolBudgets(
		t,
		poolIdle,
		conform.Ceilings(PoolChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PoolChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			PoolChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPoolSuiteGet registers the fuzz target of Pool.Get on poolIdle, with
// every check of the suite excused.
func FuzzPoolSuiteGet(f *testing.F) {
	FuzzPoolSuite(
		f,
		PoolChecks.Get.Callable,
		poolIdle,
		conform.Excuse(
			PoolChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPoolSuitePut registers the fuzz target of Pool.Put on poolIdle, with
// every check of the suite excused.
func FuzzPoolSuitePut(f *testing.F) {
	FuzzPoolSuite(
		f,
		PoolChecks.Put.Callable,
		poolIdle,
		conform.Excuse(
			PoolChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Get.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PoolChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPoolSuite measures each callable of the suite of Pool on
// poolIdle.
func BenchmarkPoolSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		BenchPoolSuite(b, PoolChecks.Get.Callable, poolIdle)
	})
	b.Run("Put", func(b *testing.B) {
		BenchPoolSuite(b, PoolChecks.Put.Callable, poolIdle)
	})
}

// batchSubjectEmbedded is Batch under the name of the field that
// batchSubject embeds.
type batchSubjectEmbedded = Batch

// batchSubject is a value of Batch that only TestBatchSuite builds.
type batchSubject struct{ batchSubjectEmbedded }

// TestBatchSuite checks the suite of Batch and the two forms of its
// application.
func TestBatchSuite(t *testing.T) {
	t.Parallel()

	t.Run("BatchSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, BatchSuite())
	})

	t.Run("RunBatchSuite", func(t *testing.T) {
		t.Parallel()
		var built *batchSubject
		RunBatchSuite(
			t,
			func() *batchSubject {
				built = &batchSubject{}
				return built
			},
			conform.With(BatchCheck{
				ID:        "Batch/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *BatchHandle) {
					got := h.Subject()
					assert.Equal(h, got, Batch(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*BatchHandle) Batch { return &batchSubject{} },
			}),
			conform.Excuse(
				BatchChecks.Write.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Write.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Write.BatchWriter(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBatchSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *batchSubject
		RunBatchSuiteFromHandle(
			t,
			func(*BatchHandle) *batchSubject {
				built = &batchSubject{}
				return built
			},
			conform.With(BatchCheck{
				ID:        "Batch/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *BatchHandle) {
					got := h.Subject()
					assert.Equal(h, got, Batch(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*BatchHandle) Batch { return &batchSubject{} },
			}),
			conform.Excuse(
				BatchChecks.Write.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Write.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Write.BatchWriter(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BatchChecks.Read.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// batchIdle returns the planted defect of a ceiling of Batch with a spend
// of nothing, a subject that answers each call at once.
func batchIdle(h *BatchHandle) Batch {
	return BatchSuite().Calls[0].Over(h, func() {})
}

// TestBatchSuiteBudgets applies the allocation ceilings of the suite of
// Batch to batchIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestBatchSuiteBudgets(t *testing.T) {
	RunBatchBudgets(
		t,
		batchIdle,
		conform.Ceilings(BatchChecks.Write.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BatchChecks.Read.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			BatchChecks.Write.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Write.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Write.BatchWriter(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBatchSuiteWrite registers the fuzz target of Batch.Write on
// batchIdle, with every check of the suite excused.
func FuzzBatchSuiteWrite(f *testing.F) {
	FuzzBatchSuite(
		f,
		BatchChecks.Write.Callable,
		batchIdle,
		conform.Excuse(
			BatchChecks.Write.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Write.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Write.BatchWriter(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBatchSuiteRead registers the fuzz target of Batch.Read on batchIdle,
// with every check of the suite excused.
func FuzzBatchSuiteRead(f *testing.F) {
	FuzzBatchSuite(
		f,
		BatchChecks.Read.Callable,
		batchIdle,
		conform.Excuse(
			BatchChecks.Write.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Write.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Write.BatchWriter(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BatchChecks.Read.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBatchSuite measures each callable of the suite of Batch on
// batchIdle.
func BenchmarkBatchSuite(b *testing.B) {
	b.Run("Write", func(b *testing.B) {
		BenchBatchSuite(b, BatchChecks.Write.Callable, batchIdle)
	})
	b.Run("Read", func(b *testing.B) {
		BenchBatchSuite(b, BatchChecks.Read.Callable, batchIdle)
	})
}

// screenedSubjectEmbedded is Screened under the name of the field that
// screenedSubject embeds.
type screenedSubjectEmbedded = Screened

// screenedSubject is a value of Screened that only TestScreenedSuite
// builds.
type screenedSubject struct{ screenedSubjectEmbedded }

// TestScreenedSuite checks the suite of Screened and the two forms of its
// application.
func TestScreenedSuite(t *testing.T) {
	t.Parallel()

	t.Run("ScreenedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ScreenedSuite())
	})

	t.Run("RunScreenedSuite", func(t *testing.T) {
		t.Parallel()
		var built *screenedSubject
		RunScreenedSuite(
			t,
			func() *screenedSubject {
				built = &screenedSubject{}
				return built
			},
			conform.With(ScreenedCheck{
				ID:        "Screened/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ScreenedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Screened(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ScreenedHandle) Screened { return &screenedSubject{} },
			}),
			conform.Excuse(
				ScreenedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScreenedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScreenedChecks.Store.Validates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScreenedChecks.Validate.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunScreenedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *screenedSubject
		RunScreenedSuiteFromHandle(
			t,
			func(*ScreenedHandle) *screenedSubject {
				built = &screenedSubject{}
				return built
			},
			conform.With(ScreenedCheck{
				ID:        "Screened/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ScreenedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Screened(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ScreenedHandle) Screened { return &screenedSubject{} },
			}),
			conform.Excuse(
				ScreenedChecks.Store.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScreenedChecks.Store.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScreenedChecks.Store.Validates(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScreenedChecks.Validate.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// screenedIdle returns the planted defect of a ceiling of Screened with a
// spend of nothing, a subject that answers each call at once.
func screenedIdle(h *ScreenedHandle) Screened {
	return ScreenedSuite().Calls[0].Over(h, func() {})
}

// TestScreenedSuiteBudgets applies the allocation ceilings of the suite of
// Screened to screenedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestScreenedSuiteBudgets(t *testing.T) {
	RunScreenedBudgets(
		t,
		screenedIdle,
		conform.Ceilings(ScreenedChecks.Store.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ScreenedChecks.Validate.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			ScreenedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Store.Validates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Validate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScreenedSuiteStore registers the fuzz target of Screened.Store on
// screenedIdle, with every check of the suite excused.
func FuzzScreenedSuiteStore(f *testing.F) {
	FuzzScreenedSuite(
		f,
		ScreenedChecks.Store.Callable,
		screenedIdle,
		conform.Excuse(
			ScreenedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Store.Validates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Validate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScreenedSuiteValidate registers the fuzz target of Screened.Validate
// on screenedIdle, with every check of the suite excused.
func FuzzScreenedSuiteValidate(f *testing.F) {
	FuzzScreenedSuite(
		f,
		ScreenedChecks.Validate.Callable,
		screenedIdle,
		conform.Excuse(
			ScreenedChecks.Store.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Store.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Store.Validates(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScreenedChecks.Validate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkScreenedSuite measures each callable of the suite of Screened
// on screenedIdle.
func BenchmarkScreenedSuite(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		BenchScreenedSuite(b, ScreenedChecks.Store.Callable, screenedIdle)
	})
	b.Run("Validate", func(b *testing.B) {
		BenchScreenedSuite(b, ScreenedChecks.Validate.Callable, screenedIdle)
	})
}

// limiterSubjectEmbedded is Limiter under the name of the field that
// limiterSubject embeds.
type limiterSubjectEmbedded = Limiter

// limiterSubject is a value of Limiter that only TestLimiterSuite builds.
type limiterSubject struct{ limiterSubjectEmbedded }

// TestLimiterSuite checks the suite of Limiter and the two forms of its
// application.
func TestLimiterSuite(t *testing.T) {
	t.Parallel()

	t.Run("LimiterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, LimiterSuite())
	})

	t.Run("RunLimiterSuite", func(t *testing.T) {
		t.Parallel()
		var built *limiterSubject
		RunLimiterSuite(
			t,
			func() *limiterSubject {
				built = &limiterSubject{}
				return built
			},
			conform.With(LimiterCheck{
				ID:        "Limiter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LimiterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Limiter(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LimiterHandle) Limiter { return &limiterSubject{} },
			}),
			conform.Excuse(
				LimiterChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimiterChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimiterChecks.Run.RateLimit(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLimiterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *limiterSubject
		RunLimiterSuiteFromHandle(
			t,
			func(*LimiterHandle) *limiterSubject {
				built = &limiterSubject{}
				return built
			},
			conform.With(LimiterCheck{
				ID:        "Limiter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LimiterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Limiter(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LimiterHandle) Limiter { return &limiterSubject{} },
			}),
			conform.Excuse(
				LimiterChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimiterChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LimiterChecks.Run.RateLimit(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// limiterIdle returns the planted defect of a ceiling of Limiter with a
// spend of nothing, a subject that answers each call at once.
func limiterIdle(h *LimiterHandle) Limiter {
	return LimiterSuite().Calls[0].Over(h, func() {})
}

// TestLimiterSuiteBudgets applies the allocation ceilings of the suite of
// Limiter to limiterIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLimiterSuiteBudgets(t *testing.T) {
	RunLimiterBudgets(
		t,
		limiterIdle,
		conform.Ceilings(LimiterChecks.Run.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			LimiterChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimiterChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimiterChecks.Run.RateLimit(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLimiterSuiteRun registers the fuzz target of Limiter.Run on
// limiterIdle, with every check of the suite excused.
func FuzzLimiterSuiteRun(f *testing.F) {
	FuzzLimiterSuite(
		f,
		LimiterChecks.Run.Callable,
		limiterIdle,
		conform.Excuse(
			LimiterChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimiterChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LimiterChecks.Run.RateLimit(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLimiterSuite measures each callable of the suite of Limiter on
// limiterIdle.
func BenchmarkLimiterSuite(b *testing.B) {
	b.Run("Run", func(b *testing.B) {
		BenchLimiterSuite(b, LimiterChecks.Run.Callable, limiterIdle)
	})
}

// rendererSubjectEmbedded is Renderer under the name of the field that
// rendererSubject embeds.
type rendererSubjectEmbedded = Renderer

// rendererSubject is a value of Renderer that only TestRendererSuite
// builds.
type rendererSubject struct{ rendererSubjectEmbedded }

// TestRendererSuite checks the suite of Renderer and the two forms of its
// application.
func TestRendererSuite(t *testing.T) {
	t.Parallel()

	t.Run("RendererSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, RendererSuite())
	})

	t.Run("RunRendererSuite", func(t *testing.T) {
		t.Parallel()
		var built *rendererSubject
		RunRendererSuite(
			t,
			func() *rendererSubject {
				built = &rendererSubject{}
				return built
			},
			conform.With(RendererCheck{
				ID:        "Renderer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *RendererHandle) {
					got := h.Subject()
					assert.Equal(h, got, Renderer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*RendererHandle) Renderer { return &rendererSubject{} },
			}),
			conform.Excuse(
				RendererChecks.Render.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Render.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Render.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Render.XSSSafe(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.Total(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunRendererSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *rendererSubject
		RunRendererSuiteFromHandle(
			t,
			func(*RendererHandle) *rendererSubject {
				built = &rendererSubject{}
				return built
			},
			conform.With(RendererCheck{
				ID:        "Renderer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *RendererHandle) {
					got := h.Subject()
					assert.Equal(h, got, Renderer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*RendererHandle) Renderer { return &rendererSubject{} },
			}),
			conform.Excuse(
				RendererChecks.Render.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Render.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Render.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Render.XSSSafe(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				RendererChecks.Classify.Total(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// rendererIdle returns the planted defect of a ceiling of Renderer with a
// spend of nothing, a subject that answers each call at once.
func rendererIdle(h *RendererHandle) Renderer {
	return RendererSuite().Calls[0].Over(h, func() {})
}

// TestRendererSuiteBudgets applies the allocation ceilings of the suite of
// Renderer to rendererIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRendererSuiteBudgets(t *testing.T) {
	RunRendererBudgets(
		t,
		rendererIdle,
		conform.Ceilings(RendererChecks.Render.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(RendererChecks.Classify.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			RendererChecks.Render.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.XSSSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Total(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRendererSuiteRender registers the fuzz target of Renderer.Render on
// rendererIdle, with every check of the suite excused.
func FuzzRendererSuiteRender(f *testing.F) {
	FuzzRendererSuite(
		f,
		RendererChecks.Render.Callable,
		rendererIdle,
		conform.Excuse(
			RendererChecks.Render.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.XSSSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Total(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRendererSuiteClassify registers the fuzz target of Renderer.Classify
// on rendererIdle, with every check of the suite excused.
func FuzzRendererSuiteClassify(f *testing.F) {
	FuzzRendererSuite(
		f,
		RendererChecks.Classify.Callable,
		rendererIdle,
		conform.Excuse(
			RendererChecks.Render.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Render.XSSSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			RendererChecks.Classify.Total(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkRendererSuite measures each callable of the suite of Renderer
// on rendererIdle.
func BenchmarkRendererSuite(b *testing.B) {
	b.Run("Render", func(b *testing.B) {
		BenchRendererSuite(b, RendererChecks.Render.Callable, rendererIdle)
	})
	b.Run("Classify", func(b *testing.B) {
		BenchRendererSuite(b, RendererChecks.Classify.Callable, rendererIdle)
	})
}

// lifecycleSubjectEmbedded is Lifecycle under the name of the field that
// lifecycleSubject embeds.
type lifecycleSubjectEmbedded = Lifecycle

// lifecycleSubject is a value of Lifecycle that only TestLifecycleSuite
// builds.
type lifecycleSubject struct{ lifecycleSubjectEmbedded }

// TestLifecycleSuite checks the suite of Lifecycle and the two forms of
// its application.
func TestLifecycleSuite(t *testing.T) {
	t.Parallel()

	t.Run("LifecycleSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, LifecycleSuite())
	})

	t.Run("RunLifecycleSuite", func(t *testing.T) {
		t.Parallel()
		var built *lifecycleSubject
		RunLifecycleSuite(
			t,
			func() *lifecycleSubject {
				built = &lifecycleSubject{}
				return built
			},
			conform.With(LifecycleCheck{
				ID:        "Lifecycle/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LifecycleHandle) {
					got := h.Subject()
					assert.Equal(h, got, Lifecycle(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LifecycleHandle) Lifecycle { return &lifecycleSubject{} },
			}),
			conform.Excuse(
				LifecycleChecks.Stop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Stop.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Put.Idempotent(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Sync.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Sync.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Sync.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLifecycleSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *lifecycleSubject
		RunLifecycleSuiteFromHandle(
			t,
			func(*LifecycleHandle) *lifecycleSubject {
				built = &lifecycleSubject{}
				return built
			},
			conform.With(LifecycleCheck{
				ID:        "Lifecycle/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LifecycleHandle) {
					got := h.Subject()
					assert.Equal(h, got, Lifecycle(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LifecycleHandle) Lifecycle { return &lifecycleSubject{} },
			}),
			conform.Excuse(
				LifecycleChecks.Stop.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Stop.RepeatSurvives(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Put.Idempotent(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Sync.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Sync.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LifecycleChecks.Sync.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// lifecycleIdle returns the planted defect of a ceiling of Lifecycle with
// a spend of nothing, a subject that answers each call at once.
func lifecycleIdle(h *LifecycleHandle) Lifecycle {
	return LifecycleSuite().Calls[0].Over(h, func() {})
}

// TestLifecycleSuiteBudgets applies the allocation ceilings of the suite
// of Lifecycle to lifecycleIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLifecycleSuiteBudgets(t *testing.T) {
	RunLifecycleBudgets(
		t,
		lifecycleIdle,
		conform.Ceilings(LifecycleChecks.Stop.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LifecycleChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LifecycleChecks.Sync.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLifecycleSuiteStop registers the fuzz target of Lifecycle.Stop on
// lifecycleIdle, with every check of the suite excused.
func FuzzLifecycleSuiteStop(f *testing.F) {
	FuzzLifecycleSuite(
		f,
		LifecycleChecks.Stop.Callable,
		lifecycleIdle,
		conform.Excuse(
			LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLifecycleSuitePut registers the fuzz target of Lifecycle.Put on
// lifecycleIdle, with every check of the suite excused.
func FuzzLifecycleSuitePut(f *testing.F) {
	FuzzLifecycleSuite(
		f,
		LifecycleChecks.Put.Callable,
		lifecycleIdle,
		conform.Excuse(
			LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLifecycleSuiteSync registers the fuzz target of Lifecycle.Sync on
// lifecycleIdle, with every check of the suite excused.
func FuzzLifecycleSuiteSync(f *testing.F) {
	FuzzLifecycleSuite(
		f,
		LifecycleChecks.Sync.Callable,
		lifecycleIdle,
		conform.Excuse(
			LifecycleChecks.Stop.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Stop.RepeatSurvives(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Put.Idempotent(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LifecycleChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLifecycleSuite measures each callable of the suite of Lifecycle
// on lifecycleIdle.
func BenchmarkLifecycleSuite(b *testing.B) {
	b.Run("Stop", func(b *testing.B) {
		BenchLifecycleSuite(b, LifecycleChecks.Stop.Callable, lifecycleIdle)
	})
	b.Run("Put", func(b *testing.B) {
		BenchLifecycleSuite(b, LifecycleChecks.Put.Callable, lifecycleIdle)
	})
	b.Run("Sync", func(b *testing.B) {
		BenchLifecycleSuite(b, LifecycleChecks.Sync.Callable, lifecycleIdle)
	})
}

// shardsSubjectEmbedded is Shards under the name of the field that
// shardsSubject embeds.
type shardsSubjectEmbedded = Shards

// shardsSubject is a value of Shards that only TestShardsSuite builds.
type shardsSubject struct{ shardsSubjectEmbedded }

// TestShardsSuite checks the suite of Shards and the two forms of its
// application.
func TestShardsSuite(t *testing.T) {
	t.Parallel()

	t.Run("ShardsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ShardsSuite())
	})

	t.Run("RunShardsSuite", func(t *testing.T) {
		t.Parallel()
		var built *shardsSubject
		RunShardsSuite(
			t,
			func() *shardsSubject {
				built = &shardsSubject{}
				return built
			},
			conform.With(ShardsCheck{
				ID:        "Shards/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ShardsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Shards(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ShardsHandle) Shards { return &shardsSubject{} },
			}),
			conform.Excuse(
				ShardsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Put.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunShardsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *shardsSubject
		RunShardsSuiteFromHandle(
			t,
			func(*ShardsHandle) *shardsSubject {
				built = &shardsSubject{}
				return built
			},
			conform.With(ShardsCheck{
				ID:        "Shards/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ShardsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Shards(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ShardsHandle) Shards { return &shardsSubject{} },
			}),
			conform.Excuse(
				ShardsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Put.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ShardsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// shardsIdle returns the planted defect of a ceiling of Shards with a
// spend of nothing, a subject that answers each call at once.
func shardsIdle(h *ShardsHandle) Shards {
	return ShardsSuite().Calls[0].Over(h, func() {})
}

// TestShardsSuiteBudgets applies the allocation ceilings of the suite of
// Shards to shardsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestShardsSuiteBudgets(t *testing.T) {
	RunShardsBudgets(
		t,
		shardsIdle,
		conform.Ceilings(ShardsChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ShardsChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			ShardsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Put.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShardsSuitePut registers the fuzz target of Shards.Put on
// shardsIdle, with every check of the suite excused.
func FuzzShardsSuitePut(f *testing.F) {
	FuzzShardsSuite(
		f,
		ShardsChecks.Put.Callable,
		shardsIdle,
		conform.Excuse(
			ShardsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Put.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShardsSuiteGet registers the fuzz target of Shards.Get on
// shardsIdle, with every check of the suite excused.
func FuzzShardsSuiteGet(f *testing.F) {
	FuzzShardsSuite(
		f,
		ShardsChecks.Get.Callable,
		shardsIdle,
		conform.Excuse(
			ShardsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Put.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ShardsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkShardsSuite measures each callable of the suite of Shards on
// shardsIdle.
func BenchmarkShardsSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchShardsSuite(b, ShardsChecks.Put.Callable, shardsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchShardsSuite(b, ShardsChecks.Get.Callable, shardsIdle)
	})
}

// scopedSubjectEmbedded is Scoped under the name of the field that
// scopedSubject embeds.
type scopedSubjectEmbedded = Scoped

// scopedSubject is a value of Scoped that only TestScopedSuite builds.
type scopedSubject struct{ scopedSubjectEmbedded }

// TestScopedSuite checks the suite of Scoped and the two forms of its
// application.
func TestScopedSuite(t *testing.T) {
	t.Parallel()

	t.Run("ScopedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ScopedSuite())
	})

	t.Run("RunScopedSuite", func(t *testing.T) {
		t.Parallel()
		var built *scopedSubject
		RunScopedSuite(
			t,
			func() *scopedSubject {
				built = &scopedSubject{}
				return built
			},
			conform.With(ScopedCheck{
				ID:        "Scoped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ScopedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Scoped(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ScopedHandle) Scoped { return &scopedSubject{} },
			}),
			conform.Excuse(
				ScopedChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Set.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunScopedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *scopedSubject
		RunScopedSuiteFromHandle(
			t,
			func(*ScopedHandle) *scopedSubject {
				built = &scopedSubject{}
				return built
			},
			conform.With(ScopedCheck{
				ID:        "Scoped/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ScopedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Scoped(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ScopedHandle) Scoped { return &scopedSubject{} },
			}),
			conform.Excuse(
				ScopedChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Set.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Set.Partition(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ScopedChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// scopedIdle returns the planted defect of a ceiling of Scoped with a
// spend of nothing, a subject that answers each call at once.
func scopedIdle(h *ScopedHandle) Scoped {
	return ScopedSuite().Calls[0].Over(h, func() {})
}

// TestScopedSuiteBudgets applies the allocation ceilings of the suite of
// Scoped to scopedIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestScopedSuiteBudgets(t *testing.T) {
	RunScopedBudgets(
		t,
		scopedIdle,
		conform.Ceilings(ScopedChecks.Set.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ScopedChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			ScopedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopedSuiteSet registers the fuzz target of Scoped.Set on
// scopedIdle, with every check of the suite excused.
func FuzzScopedSuiteSet(f *testing.F) {
	FuzzScopedSuite(
		f,
		ScopedChecks.Set.Callable,
		scopedIdle,
		conform.Excuse(
			ScopedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzScopedSuiteGet registers the fuzz target of Scoped.Get on
// scopedIdle, with every check of the suite excused.
func FuzzScopedSuiteGet(f *testing.F) {
	FuzzScopedSuite(
		f,
		ScopedChecks.Get.Callable,
		scopedIdle,
		conform.Excuse(
			ScopedChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Set.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Set.Partition(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ScopedChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkScopedSuite measures each callable of the suite of Scoped on
// scopedIdle.
func BenchmarkScopedSuite(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		BenchScopedSuite(b, ScopedChecks.Set.Callable, scopedIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchScopedSuite(b, ScopedChecks.Get.Callable, scopedIdle)
	})
}

// plainSubjectEmbedded is Plain under the name of the field that
// plainSubject embeds.
type plainSubjectEmbedded = Plain

// plainSubject is a value of Plain that only TestPlainSuite builds.
type plainSubject struct{ plainSubjectEmbedded }

// TestPlainSuite checks the suite of Plain and the two forms of its
// application.
func TestPlainSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlainSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PlainSuite())
	})

	t.Run("RunPlainSuite", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		RunPlainSuite(
			t,
			func() *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlainHandle) Plain { return &plainSubject{} },
			}),
			conform.Excuse(
				PlainChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Watch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Watch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPlainSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		RunPlainSuiteFromHandle(
			t,
			func(*PlainHandle) *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlainHandle) Plain { return &plainSubject{} },
			}),
			conform.Excuse(
				PlainChecks.Set.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Watch.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.Watch.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				PlainChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// plainIdle returns the planted defect of a ceiling of Plain with a spend
// of nothing, a subject that answers each call at once.
func plainIdle(h *PlainHandle) Plain {
	return PlainSuite().Calls[0].Over(h, func() {})
}

// TestPlainSuiteBudgets applies the allocation ceilings of the suite of
// Plain to plainIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPlainSuiteBudgets(t *testing.T) {
	RunPlainBudgets(
		t,
		plainIdle,
		conform.Ceilings(PlainChecks.Set.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PlainChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PlainChecks.Watch.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(PlainChecks.All.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteSet registers the fuzz target of Plain.Set on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteSet(f *testing.F) {
	FuzzPlainSuite(
		f,
		PlainChecks.Set.Callable,
		plainIdle,
		conform.Excuse(
			PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteGet registers the fuzz target of Plain.Get on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteGet(f *testing.F) {
	FuzzPlainSuite(
		f,
		PlainChecks.Get.Callable,
		plainIdle,
		conform.Excuse(
			PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteWatch registers the fuzz target of Plain.Watch on
// plainIdle, with every check of the suite excused.
func FuzzPlainSuiteWatch(f *testing.F) {
	FuzzPlainSuite(
		f,
		PlainChecks.Watch.Callable,
		plainIdle,
		conform.Excuse(
			PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuiteAll registers the fuzz target of Plain.All on plainIdle,
// with every check of the suite excused.
func FuzzPlainSuiteAll(f *testing.F) {
	FuzzPlainSuite(
		f,
		PlainChecks.All.Callable,
		plainIdle,
		conform.Excuse(
			PlainChecks.Set.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.Watch.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			PlainChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPlainSuite measures each callable of the suite of Plain on
// plainIdle.
func BenchmarkPlainSuite(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		BenchPlainSuite(b, PlainChecks.Set.Callable, plainIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchPlainSuite(b, PlainChecks.Get.Callable, plainIdle)
	})
	b.Run("Watch", func(b *testing.B) {
		BenchPlainSuite(b, PlainChecks.Watch.Callable, plainIdle)
	})
	b.Run("All", func(b *testing.B) {
		BenchPlainSuite(b, PlainChecks.All.Callable, plainIdle)
	})
}

// sessionsSubjectEmbedded is Sessions under the name of the field that
// sessionsSubject embeds.
type sessionsSubjectEmbedded = Sessions

// sessionsSubject is a value of Sessions that only TestSessionsSuite
// builds.
type sessionsSubject struct{ sessionsSubjectEmbedded }

// TestSessionsSuite checks the suite of Sessions and the two forms of its
// application.
func TestSessionsSuite(t *testing.T) {
	t.Parallel()

	t.Run("SessionsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, SessionsSuite())
	})

	t.Run("RunSessionsSuite", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		RunSessionsSuite(
			t,
			func() *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *SessionsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Sessions(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*SessionsHandle) Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				SessionsChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SessionsChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SessionsChecks.Open.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSessionsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		RunSessionsSuiteFromHandle(
			t,
			func(*SessionsHandle) *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *SessionsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Sessions(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*SessionsHandle) Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				SessionsChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SessionsChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SessionsChecks.Open.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// sessionsIdle returns the planted defect of a ceiling of Sessions with a
// spend of nothing, a subject that answers each call at once.
func sessionsIdle(h *SessionsHandle) Sessions {
	return SessionsSuite().Calls[0].Over(h, func() {})
}

// TestSessionsSuiteBudgets applies the allocation ceilings of the suite of
// Sessions to sessionsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestSessionsSuiteBudgets(t *testing.T) {
	RunSessionsBudgets(
		t,
		sessionsIdle,
		conform.Ceilings(SessionsChecks.Open.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			SessionsChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SessionsChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SessionsChecks.Open.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSessionsSuiteOpen registers the fuzz target of Sessions.Open on
// sessionsIdle, with every check of the suite excused.
func FuzzSessionsSuiteOpen(f *testing.F) {
	FuzzSessionsSuite(
		f,
		SessionsChecks.Open.Callable,
		sessionsIdle,
		conform.Excuse(
			SessionsChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SessionsChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			SessionsChecks.Open.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSessionsSuite measures each callable of the suite of Sessions
// on sessionsIdle.
func BenchmarkSessionsSuite(b *testing.B) {
	b.Run("Open", func(b *testing.B) {
		BenchSessionsSuite(b, SessionsChecks.Open.Callable, sessionsIdle)
	})
}
