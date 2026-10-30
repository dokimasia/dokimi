package generic

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Cache:
//
// Cache.Get/smoke: Get returns without a panic for derived inputs.
//
// Cache.Get/zero-on-miss: Get returns the zero value beside a false
// presence flag.
//
// Cache.Set/smoke: Set returns without a panic for derived inputs.

// The derived checks of Cache that witness their laws at the shallow end:
//
// Cache.Get/smoke
//
// Cache.Set/smoke

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

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c cacheGet) Smoke() conform.ID[Cache[int, int]] {
	return c.Check("smoke")
}

// ZeroOnMiss returns the identity of the check that Get returns the zero
// value beside a false presence flag.
func (c cacheGet) ZeroOnMiss() conform.ID[Cache[int, int]] {
	return c.Check("zero-on-miss")
}

// cacheSet is the name of Cache.Set, with the identities of its checks.
type cacheSet struct {
	conform.Callable[Cache[int, int]]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c cacheSet) Smoke() conform.ID[Cache[int, int]] {
	return c.Check("smoke")
}

// CacheSuite returns the checks of Cache, and the calls and the budgets of
// its callables.
func CacheSuite() conform.Suite[Cache[int, int]] {
	return conform.Suite[Cache[int, int]]{
		Declaration: "Cache",
		Package:     "example.com/shop/generic",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Cache[int, int]]{
			{
				ID:          CacheChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CacheHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Get(int(42))
						return nil
					})
				},
				Defect: cacheGetPanics,
			},
			{
				ID:          CacheChecks.Get.ZeroOnMiss(),
				Statement:   "Get returns the zero value beside a false presence flag",
				Claim:       spec.ZeroOnMiss,
				Probe:       spec.ZeroBesideFalse,
				Reason:      spec.ValueBesideFalse,
				Law:         spec.NoValueBesideAFalseFlag,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CacheHandle) {
					probe.ZeroBesideFalse(h, func(ctx context.Context) (int, bool) {
						return h.Subject().Get(int(42))
					}, int(0))
				},
				Defect: cacheGetValueBesideAFalseFlag,
			},
			{
				ID:          CacheChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CacheHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Set(int(42), int(42))
						return nil
					})
				},
				Defect: cacheSetPanics,
			},
		},
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

// RunCacheBudgets applies the allocation ceilings of Cache and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
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

// cacheGetPanicsDefect is the planted defect of a Cache whose Get panics.
type cacheGetPanicsDefect struct {
	cacheOverEmbedded
}

// Get is the override of the planted defect.
func (*cacheGetPanicsDefect) Get(p0 int) (r0 int, r1 bool) {
	panic("the planted defect of a Cache whose Get panics")
}

// cacheGetPanics returns the planted defect of a Cache whose Get panics
// over a fresh subject.
func cacheGetPanics(h *CacheHandle) Cache[int, int] {
	return &cacheGetPanicsDefect{cacheOverEmbedded: h.Fresh()}
}

// cacheGetValueBesideAFalseFlagDefect is the planted defect of a Cache
// whose Get returns a value beside a false presence flag.
type cacheGetValueBesideAFalseFlagDefect struct {
	cacheOverEmbedded
}

// Get is the override of the planted defect.
func (d *cacheGetValueBesideAFalseFlagDefect) Get(p0 int) (r0 int, r1 bool) {
	r0, r1 = d.cacheOverEmbedded.Get(p0)
	if !r1 {
		r0 = int(7)
	}
	return r0, r1
}

// cacheGetValueBesideAFalseFlag returns the planted defect of a Cache
// whose Get returns a value beside a false presence flag over a fresh
// subject.
func cacheGetValueBesideAFalseFlag(h *CacheHandle) Cache[int, int] {
	return &cacheGetValueBesideAFalseFlagDefect{cacheOverEmbedded: h.Fresh()}
}

// cacheSetPanicsDefect is the planted defect of a Cache whose Set panics.
type cacheSetPanicsDefect struct {
	cacheOverEmbedded
}

// Set is the override of the planted defect.
func (*cacheSetPanicsDefect) Set(p0, p1 int) {
	panic("the planted defect of a Cache whose Set panics")
}

// cacheSetPanics returns the planted defect of a Cache whose Set panics
// over a fresh subject.
func cacheSetPanics(h *CacheHandle) Cache[int, int] {
	return &cacheSetPanicsDefect{cacheOverEmbedded: h.Fresh()}
}
