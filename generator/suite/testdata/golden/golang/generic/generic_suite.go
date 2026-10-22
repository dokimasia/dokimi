package generic

import (
	"testing"

	"go.dokimi.dev/conform"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Cache writes.
type (
	CacheHandle = conform.Handle[Cache[int, int]]
	CacheCheck  = conform.Check[Cache[int, int]]
	CacheOption = conform.Option[Cache[int, int]]
)

// CacheChecks contains the name of each callable of Cache, with the
// identities of its checks.
var CacheChecks = struct {
	Get cacheGet
	Set cacheSet
}{
	Get: cacheGet{"Cache.Get"},
	Set: cacheSet{"Cache.Set"},
}

// cacheGet is the name of Cache.Get, with the identities of its checks.
type cacheGet struct {
	conform.Callable[Cache[int, int]]
}

// cacheSet is the name of Cache.Set, with the identities of its checks.
type cacheSet struct {
	conform.Callable[Cache[int, int]]
}

// CacheSuite returns the checks of Cache.
func CacheSuite() conform.Suite[Cache[int, int]] {
	return conform.Suite[Cache[int, int]]{Declaration: "Cache"}
}

// RunCacheSuite applies the checks of Cache, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunCacheSuite[T Cache[int, int]](t *testing.T, newSubject func() T, opts ...CacheOption) {
	t.Helper()
	RunCacheSuiteFromHandle(t, func(*CacheHandle) T { return newSubject() }, opts...)
}

// RunCacheSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCacheSuiteFromHandle[T Cache[int, int]](
	t *testing.T,
	newSubject func(h *CacheHandle) T,
	opts ...CacheOption,
) {
	t.Helper()
	conform.Run(
		t,
		CacheSuite(),
		func(h *CacheHandle) Cache[int, int] { return newSubject(h) },
		opts...)
}
