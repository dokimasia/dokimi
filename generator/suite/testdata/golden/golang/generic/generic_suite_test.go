package generic_test

import (
	"testing"

	"example.com/shop/generic"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// cacheSubjectEmbedded is generic.Cache[int, int] under the name of the
// field that cacheSubject embeds.
type cacheSubjectEmbedded = generic.Cache[int, int]

// cacheSubject is a value of Cache that only TestCacheSuite builds.
type cacheSubject struct{ cacheSubjectEmbedded }

// TestCacheSuite checks the suite of Cache and the two forms of its
// application.
func TestCacheSuite(t *testing.T) {
	t.Parallel()

	t.Run("CacheSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, generic.CacheSuite())
	})

	t.Run("RunCacheSuite", func(t *testing.T) {
		t.Parallel()
		var built *cacheSubject
		generic.RunCacheSuite(
			t,
			func() *cacheSubject {
				built = &cacheSubject{}
				return built
			},
			conform.With(generic.CacheCheck{
				ID:        "Cache/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *generic.CacheHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						generic.Cache[int, int](built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*generic.CacheHandle) generic.Cache[int, int] { return &cacheSubject{} },
			}),
			conform.ExcuseSuite[generic.Cache[int, int]](
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCacheSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *cacheSubject
		generic.RunCacheSuiteFromHandle(
			t,
			func(*generic.CacheHandle) *cacheSubject {
				built = &cacheSubject{}
				return built
			},
			conform.With(generic.CacheCheck{
				ID:        "Cache/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *generic.CacheHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						generic.Cache[int, int](built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*generic.CacheHandle) generic.Cache[int, int] { return &cacheSubject{} },
			}),
			conform.ExcuseSuite[generic.Cache[int, int]](
				"the test applies the check of its option alone",
			),
		)
	})
}

// cacheSuite is the suite of Cache that the tests of the planted defect of
// a ceiling apply.
var cacheSuite = generic.CacheSuite()

// cacheIdle returns the planted defect of a ceiling of Cache with a spend
// of nothing. Each call of the defect returns its zero results at once.
func cacheIdle(h *generic.CacheHandle) generic.Cache[int, int] {
	return cacheSuite.Calls[0].Over(h, func() {})
}

// TestCacheSuiteBudgets applies the allocation ceilings of the suite of
// Cache to cacheIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestCacheSuiteBudgets(t *testing.T) {
	generic.RunCacheBudgets(
		t,
		cacheIdle,
		conform.Ceilings(generic.CacheChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(generic.CacheChecks.Set.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.ExcuseSuite[generic.Cache[int, int]](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzCacheSuiteGet registers the fuzz target of Cache.Get on cacheIdle,
// with every check of the suite excused.
func FuzzCacheSuiteGet(f *testing.F) {
	generic.FuzzCacheSuite(
		f,
		generic.CacheChecks.Get.Callable,
		cacheIdle,
		conform.ExcuseSuite[generic.Cache[int, int]](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzCacheSuiteSet registers the fuzz target of Cache.Set on cacheIdle,
// with every check of the suite excused.
func FuzzCacheSuiteSet(f *testing.F) {
	generic.FuzzCacheSuite(
		f,
		generic.CacheChecks.Set.Callable,
		cacheIdle,
		conform.ExcuseSuite[generic.Cache[int, int]](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkCacheSuite measures each callable of the suite of Cache on
// cacheIdle.
func BenchmarkCacheSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		generic.BenchCacheSuite(b, generic.CacheChecks.Get.Callable, cacheIdle)
	})
	b.Run("Set", func(b *testing.B) {
		generic.BenchCacheSuite(b, generic.CacheChecks.Set.Callable, cacheIdle)
	})
}
