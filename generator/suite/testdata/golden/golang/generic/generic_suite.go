package generic

import (
	"testing"

	"go.dokimi.dev/dokimi-go/conform"
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

// RunCache applies the checks of Cache, and the checks that opts add, to a
// fresh subject of newSubject for each check.
func RunCache[T Cache[int, int]](t *testing.T, newSubject func() T, opts ...CacheOption) {
	t.Helper()
	conform.Verify(t, CacheSuite(), func(*CacheHandle) Cache[int, int] { return newSubject() }, opts...)
}

// RunCacheOn applies the same checks to a subject that newSubject builds
// on the handle of each check.
func RunCacheOn[T Cache[int, int]](t *testing.T, newSubject func(h *CacheHandle) T, opts ...CacheOption) {
	t.Helper()
	conform.Verify(t, CacheSuite(), func(h *CacheHandle) Cache[int, int] { return newSubject(h) }, opts...)
}

// ControlCache applies the same checks to a fresh control of newControl
// for each check. A control is a correct implementation that differs from
// the subject.
func ControlCache[T Cache[int, int]](t *testing.T, newControl func() T, opts ...CacheOption) {
	t.Helper()
	conform.Control(t, CacheSuite(), func(*CacheHandle) Cache[int, int] { return newControl() }, opts...)
}

// ControlCacheOn applies the same checks to a control that newControl
// builds on the handle of each check.
func ControlCacheOn[T Cache[int, int]](t *testing.T, newControl func(h *CacheHandle) T, opts ...CacheOption) {
	t.Helper()
	conform.Control(t, CacheSuite(), func(h *CacheHandle) Cache[int, int] { return newControl(h) }, opts...)
}
