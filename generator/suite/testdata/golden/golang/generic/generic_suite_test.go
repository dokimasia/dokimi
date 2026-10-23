package generic

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// cacheSubjectEmbedded is Cache[int, int] under the name of the field that
// cacheSubject embeds.
type cacheSubjectEmbedded = Cache[int, int]

// cacheSubject is a value of Cache that only TestCacheSuite builds.
type cacheSubject struct{ cacheSubjectEmbedded }

// TestCacheSuite checks the suite of Cache and the two forms of its
// application.
func TestCacheSuite(t *testing.T) {
	t.Parallel()

	t.Run("CacheSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, CacheSuite())
	})

	t.Run("RunCacheSuite", func(t *testing.T) {
		t.Parallel()
		var built *cacheSubject
		RunCacheSuite(
			t,
			func() *cacheSubject {
				built = &cacheSubject{}
				return built
			},
			conform.With(CacheCheck{
				ID:        "Cache/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CacheHandle) {
					got := h.Subject()
					assert.Equal(h, got, Cache[int, int](built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CacheHandle) Cache[int, int] { return &cacheSubject{} },
			}),
		)
	})

	t.Run("RunCacheSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *cacheSubject
		RunCacheSuiteFromHandle(
			t,
			func(*CacheHandle) *cacheSubject {
				built = &cacheSubject{}
				return built
			},
			conform.With(CacheCheck{
				ID:        "Cache/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CacheHandle) {
					got := h.Subject()
					assert.Equal(h, got, Cache[int, int](built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CacheHandle) Cache[int, int] { return &cacheSubject{} },
			}),
		)
	})
}

// cacheIdle returns the planted defect of a ceiling of Cache with a spend
// of nothing, a subject that answers each call at once.
func cacheIdle(h *CacheHandle) Cache[int, int] {
	return CacheSuite().Calls[0].Over(h, func() {})
}

// TestCacheSuiteBudgets applies the allocation ceilings of the suite of
// Cache to cacheIdle, and a ceiling of 0 to each callable without a budget
// whose calls build no input of their own.
func TestCacheSuiteBudgets(t *testing.T) {
	RunCacheBudgets(
		t,
		cacheIdle,
		conform.Ceilings(CacheChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(CacheChecks.Set.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzCacheSuiteGet registers the fuzz target of Cache.Get on cacheIdle,
// with every check of the suite excused.
func FuzzCacheSuiteGet(f *testing.F) {
	FuzzCacheSuite(
		f,
		CacheChecks.Get.Callable,
		cacheIdle,
	)
}

// FuzzCacheSuiteSet registers the fuzz target of Cache.Set on cacheIdle,
// with every check of the suite excused.
func FuzzCacheSuiteSet(f *testing.F) {
	FuzzCacheSuite(
		f,
		CacheChecks.Set.Callable,
		cacheIdle,
	)
}

// BenchmarkCacheSuite measures each callable of the suite of Cache on
// cacheIdle.
func BenchmarkCacheSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		BenchCacheSuite(b, CacheChecks.Get.Callable, cacheIdle)
	})
	b.Run("Set", func(b *testing.B) {
		BenchCacheSuite(b, CacheChecks.Set.Callable, cacheIdle)
	})
}
