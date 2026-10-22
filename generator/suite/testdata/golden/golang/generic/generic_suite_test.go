package generic

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// cacheSubject is a value of Cache that only TestCacheSuite builds.
type cacheSubject struct{ Cache[int, int] }

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
