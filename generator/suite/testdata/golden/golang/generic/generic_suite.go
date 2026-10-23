package generic

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
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

// CacheSuite returns the checks of Cache, and the calls and the budgets of
// its callables.
func CacheSuite() conform.Suite[Cache[int, int]] {
	return conform.Suite[Cache[int, int]]{
		Declaration: "Cache",
		Calls: []conform.Call[Cache[int, int]]{
			{
				Calls:    cacheGetCalls,
				Over:     newCacheOver,
				Callable: CacheChecks.Get.Callable,
			},
			{
				Calls:    cacheSetCalls,
				Over:     newCacheOver,
				Callable: CacheChecks.Set.Callable,
			},
		},
	}
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
	conform.Run(t, CacheSuite(), cacheSubjects(newSubject), opts...)
}

// FuzzCacheSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCacheSuite[T Cache[int, int]](
	f *testing.F,
	callable conform.Callable[Cache[int, int]],
	newSubject func(h *CacheHandle) T,
	opts ...CacheOption,
) {
	f.Helper()
	conform.Fuzz(f, CacheSuite(), callable, cacheSubjects(newSubject), opts...)
}

// RunCacheBudgets applies the allocation ceilings of Cache, and the
// ceilings that opts add, to a fresh subject of newSubject for each case.
func RunCacheBudgets[T Cache[int, int]](
	t *testing.T,
	newSubject func(h *CacheHandle) T,
	opts ...CacheOption,
) {
	t.Helper()
	conform.RunBudgets(t, CacheSuite(), cacheSubjects(newSubject), opts...)
}

// BenchCacheSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchCacheSuite[T Cache[int, int]](
	b conform.B,
	callable conform.Callable[Cache[int, int]],
	newSubject func(h *CacheHandle) T,
	opts ...CacheOption,
) {
	b.Helper()
	conform.Bench(b, CacheSuite(), callable, cacheSubjects(newSubject), opts...)
}

// cacheSubjects returns the factory of the subjects of Cache that calls
// newSubject with the handle of each check.
func cacheSubjects[T Cache[int, int]](
	newSubject func(h *CacheHandle) T,
) func(h *CacheHandle) Cache[int, int] {
	return func(h *CacheHandle) Cache[int, int] { return newSubject(h) }
}

// cacheGetCalls generates the calls of Cache.Get with the inputs that each
// case draws.
var cacheGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cache[int, int]) error {
	key := c.Draw(prop.Of[int](), "key")
	return func(_ context.Context, subject Cache[int, int]) error {
		_, _ = subject.Get(key)
		return nil
	}
})

// cacheSetCalls generates the calls of Cache.Set with the inputs that each
// case draws.
var cacheSetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cache[int, int]) error {
	key := c.Draw(prop.Of[int](), "key")
	value := c.Draw(prop.Of[int](), "value")
	return func(_ context.Context, subject Cache[int, int]) error {
		subject.Set(key, value)
		return nil
	}
})

// cacheOverEmbedded is Cache[int, int] under the name of the field that
// cacheOver embeds.
type cacheOverEmbedded = Cache[int, int]

// cacheOver is a subject of Cache whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type cacheOver struct {
	cacheOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o cacheOver) Get(int) (r0 int, r1 bool) {
	o.spend()
	return r0, r1
}

// Set calls spend and returns its zero results.
func (o cacheOver) Set(int, int) {
	o.spend()
}

// newCacheOver returns the planted defect of a ceiling of Cache, whose
// callables call spend.
func newCacheOver(_ *CacheHandle, spend func()) Cache[int, int] {
	return cacheOver{spend: spend}
}
