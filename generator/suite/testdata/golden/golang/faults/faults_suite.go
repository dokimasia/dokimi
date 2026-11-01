package faults

import (
	"context"
	"sync/atomic"
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
// Cache.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Cache.Get/cancel: Get returns an error for a cancelled context.
//
// Cache.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Cache.Get/failure-is-reported: Get returns an error when a call of its
// dependency fails.
//
// Cache.Get/failure-is-reported#generated: Get returns an error when a
// call of its dependency fails.
//
// Cache.Get/zero-on-error: Get returns the zero value beside an error.
//
// Cache.Get/zero-on-error#generated: Get returns the zero value beside an
// error.
//
// Cache.Len/smoke: Len returns without a panic for derived inputs.
//
// Cache.Len/deterministic: Len returns the same value for the same inputs.

// The derived checks of Cache that witness their laws at the shallow end:
//
// Cache.Get/smoke
//
// Cache.Len/smoke
//
// Cache.Len/deterministic

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Cache writes.
type (
	CacheHandle = conform.Handle[Cache]
	CacheCheck  = conform.Check[Cache]
	CacheOption = conform.Option[Cache]
)

// CacheChecks contains the name of each callable of Cache, with the
// identities of its checks.
var CacheChecks = struct {
	Get cacheGet
	Len cacheLen
}{
	Get: cacheGet{"Cache.Get"},
	Len: cacheLen{"Cache.Len"},
}

// cacheGet is the name of Cache.Get, with the identities of its checks.
type cacheGet struct {
	conform.Callable[Cache]
}

// ServesTheStoredValue returns the identity of the check that Get returns
// the value that the backing store keeps, after a failed call of the
// store.
func (c cacheGet) ServesTheStoredValue() conform.ID[Cache] {
	return c.Check("serves-the-stored-value")
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c cacheGet) Smoke() conform.ID[Cache] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c cacheGet) SmokeGenerated() conform.ID[Cache] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c cacheGet) Cancel() conform.ID[Cache] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c cacheGet) CancelGenerated() conform.ID[Cache] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Get returns an
// error when a call of its dependency fails.
func (c cacheGet) FailureIsReported() conform.ID[Cache] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Get
// returns an error when a call of its dependency fails.
func (c cacheGet) FailureIsReportedGenerated() conform.ID[Cache] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c cacheGet) ZeroOnError() conform.ID[Cache] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c cacheGet) ZeroOnErrorGenerated() conform.ID[Cache] {
	return c.Check("zero-on-error#generated")
}

// cacheLen is the name of Cache.Len, with the identities of its checks.
type cacheLen struct {
	conform.Callable[Cache]
}

// Smoke returns the identity of the check that Len returns without a panic
// for derived inputs.
func (c cacheLen) Smoke() conform.ID[Cache] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Len returns the
// same value for the same inputs.
func (c cacheLen) Deterministic() conform.ID[Cache] {
	return c.Check("deterministic")
}

// CacheSuite returns the checks of Cache, and the calls and the budgets of
// its callables.
func CacheSuite() conform.Suite[Cache] {
	return conform.Suite[Cache]{
		Declaration: "Cache",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Cache]{
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
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cacheGetPanics,
			},
			{
				ID:          CacheChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CacheHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: cacheGetPanics,
			},
			{
				ID:          CacheChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CacheHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cacheGetIgnoresCancel,
			},
			{
				ID:          CacheChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CacheHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: cacheGetIgnoresCancel,
			},
			{
				ID:        CacheChecks.Get.FailureIsReported(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *CacheHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cacheGetSwallowsFailure,
			},
			{
				ID:        CacheChecks.Get.FailureIsReportedGenerated(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *CacheHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: cacheGetSwallowsFailure,
			},
			{
				ID:          CacheChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CacheHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: cacheGetValueBeside,
			},
			{
				ID:          CacheChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CacheHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Value{})
				},
				Defect: cacheGetValueBeside,
			},
			{
				ID:          CacheChecks.Len.Smoke(),
				Statement:   "Len returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CacheHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Len()
						return nil
					})
				},
				Defect: cacheLenPanics,
			},
			{
				ID:          CacheChecks.Len.Deterministic(),
				Statement:   "Len returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Cache]{CacheChecks.Len.Smoke()},
				Apply: func(h *CacheHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Len()
						return r0, nil
					})
				},
				Defect: cacheLenDrifts,
			},
			{
				ID:          CacheChecks.Get.ServesTheStoredValue(),
				Statement:   "Get returns the value that the backing store keeps, after a failed call of the store",
				Apply:       ServesTheStoredValue,
				Defect:      servesStaleValues,
				Doors:       conform.Doors[Cache]{Faults: faulty},
				Machineries: []spec.MachineryID{spec.CallThatFailsOnDemand},
			},
		},
		Calls: []conform.Call[Cache]{
			{
				Calls:    cacheGetCalls,
				Over:     newCacheOver,
				Callable: CacheChecks.Get.Callable,
			},
			{
				Calls:    cacheLenCalls,
				Over:     newCacheOver,
				Callable: CacheChecks.Len.Callable,
			},
		},
	}
}

// RunCacheSuite applies the checks of Cache, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunCacheSuite[T Cache](t *testing.T, newSubject func() T, opts ...CacheOption) {
	t.Helper()
	RunCacheSuiteFromHandle(t, func(*CacheHandle) T { return newSubject() }, opts...)
}

// RunCacheSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCacheSuiteFromHandle[T Cache](
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
func FuzzCacheSuite[T Cache](
	f *testing.F,
	callable conform.Callable[Cache],
	newSubject func(h *CacheHandle) T,
	opts ...CacheOption,
) {
	f.Helper()
	conform.Fuzz(f, CacheSuite(), callable, cacheSubjects(newSubject), opts...)
}

// RunCacheBudgets applies the allocation ceilings of Cache and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunCacheBudgets[T Cache](
	t *testing.T,
	newSubject func(h *CacheHandle) T,
	opts ...CacheOption,
) {
	t.Helper()
	conform.RunBudgets(t, CacheSuite(), cacheSubjects(newSubject), opts...)
}

// BenchCacheSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchCacheSuite[T Cache](
	b conform.B,
	callable conform.Callable[Cache],
	newSubject func(h *CacheHandle) T,
	opts ...CacheOption,
) {
	b.Helper()
	conform.Bench(b, CacheSuite(), callable, cacheSubjects(newSubject), opts...)
}

// cacheSubjects returns the factory of the subjects of Cache that calls
// newSubject with the handle of each check.
func cacheSubjects[T Cache](newSubject func(h *CacheHandle) T) func(h *CacheHandle) Cache {
	return func(h *CacheHandle) Cache { return newSubject(h) }
}

// cacheGetCalls generates the calls of Cache.Get with the inputs that each
// case draws.
var cacheGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cache) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Cache) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// cacheLenCalls generates the calls of Cache.Len with the inputs that each
// case draws.
var cacheLenCalls = prop.Composite(func(*prop.Case) func(context.Context, Cache) error {
	return func(_ context.Context, subject Cache) error {
		_ = subject.Len()
		return nil
	}
})

// cacheOverEmbedded is Cache under the name of the field that cacheOver
// embeds.
type cacheOverEmbedded = Cache

// cacheOver is a subject of Cache whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type cacheOver struct {
	cacheOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o cacheOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Len calls spend and returns its zero results.
func (o cacheOver) Len() (r0 int) {
	o.spend()
	return r0
}

// newCacheOver returns the planted defect of a ceiling of Cache, whose
// callables call spend.
func newCacheOver(_ *CacheHandle, spend func()) Cache {
	return cacheOver{spend: spend}
}

// cacheGetPanicsDefect is the planted defect of a Cache whose Get panics.
type cacheGetPanicsDefect struct {
	cacheOverEmbedded
}

// Get is the override of the planted defect.
func (*cacheGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Cache whose Get panics")
}

// cacheGetPanics returns the planted defect of a Cache whose Get panics
// over a fresh subject.
func cacheGetPanics(h *CacheHandle) Cache {
	return &cacheGetPanicsDefect{cacheOverEmbedded: h.Fresh()}
}

// cacheGetIgnoresCancelDefect is the planted defect of a Cache whose Get
// ignores the context.
type cacheGetIgnoresCancelDefect struct {
	cacheOverEmbedded
}

// Get is the override of the planted defect.
func (d *cacheGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cacheOverEmbedded.Get(p0, p1)
}

// cacheGetIgnoresCancel returns the planted defect of a Cache whose Get
// ignores the context over a fresh subject.
func cacheGetIgnoresCancel(h *CacheHandle) Cache {
	return &cacheGetIgnoresCancelDefect{cacheOverEmbedded: h.Fresh()}
}

// cacheGetSwallowsFailureDefect is the planted defect of a Cache whose Get
// returns no error when its dependency fails.
type cacheGetSwallowsFailureDefect struct {
	cacheOverEmbedded
	h *CacheHandle
}

// Get is the override of the planted defect.
func (d *cacheGetSwallowsFailureDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.cacheOverEmbedded.Get(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(Value)
		r1 = nil
	}
	return r0, r1
}

// cacheGetSwallowsFailure returns the planted defect of a Cache whose Get
// returns no error when its dependency fails over a fresh subject.
func cacheGetSwallowsFailure(h *CacheHandle) Cache {
	return &cacheGetSwallowsFailureDefect{cacheOverEmbedded: h.Fresh(), h: h}
}

// cacheGetValueBesideDefect is the planted defect of a Cache whose Get
// returns a value beside an error.
type cacheGetValueBesideDefect struct {
	cacheOverEmbedded
}

// Get is the override of the planted defect.
func (d *cacheGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.cacheOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// cacheGetValueBeside returns the planted defect of a Cache whose Get
// returns a value beside an error over a fresh subject.
func cacheGetValueBeside(h *CacheHandle) Cache {
	return &cacheGetValueBesideDefect{cacheOverEmbedded: h.Fresh()}
}

// cacheLenPanicsDefect is the planted defect of a Cache whose Len panics.
type cacheLenPanicsDefect struct {
	cacheOverEmbedded
}

// Len is the override of the planted defect.
func (*cacheLenPanicsDefect) Len() (r0 int) {
	panic("the planted defect of a Cache whose Len panics")
}

// cacheLenPanics returns the planted defect of a Cache whose Len panics
// over a fresh subject.
func cacheLenPanics(h *CacheHandle) Cache {
	return &cacheLenPanicsDefect{cacheOverEmbedded: h.Fresh()}
}

// cacheLenDriftsDefect is the planted defect of a Cache whose Len returns
// another value on the repeat.
type cacheLenDriftsDefect struct {
	cacheOverEmbedded
	calls int
}

// Len is the override of the planted defect.
func (d *cacheLenDriftsDefect) Len() (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// cacheLenDrifts returns the planted defect of a Cache whose Len returns
// another value on the repeat over a fresh subject.
func cacheLenDrifts(h *CacheHandle) Cache {
	return &cacheLenDriftsDefect{cacheOverEmbedded: h.Fresh()}
}

// The derived checks of Writer:
//
// Writer.Put/smoke: Put returns without a panic for derived inputs.
//
// Writer.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Writer.Put/cancel: Put returns an error for a cancelled context.
//
// Writer.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Writer.Put/failure-is-reported: Put returns an error when a call of its
// dependency fails.
//
// Writer.Put/failure-is-reported#generated: Put returns an error when a
// call of its dependency fails.
//
// Writer.Put/no-partial-write: Get finds no effect of a Put that failed
// with its dependency.
//
// Writer.Put/no-partial-write#generated: Get finds no effect of a Put that
// failed with its dependency.
//
// Writer.Get/smoke: Get returns without a panic for derived inputs.
//
// Writer.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Writer.Get/cancel: Get returns an error for a cancelled context.
//
// Writer.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Writer.Get/failure-is-reported: Get returns an error when a call of its
// dependency fails.
//
// Writer.Get/failure-is-reported#generated: Get returns an error when a
// call of its dependency fails.
//
// Writer.Get/zero-on-error: Get returns the zero value beside an error.
//
// Writer.Get/zero-on-error#generated: Get returns the zero value beside an
// error.

// The derived checks of Writer that witness their laws at the shallow end:
//
// Writer.Put/smoke
//
// Writer.Get/smoke

// The unclaimed classifications of Writer:
//
// The shape:writer of Writer.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Writer writes.
type (
	WriterHandle = conform.Handle[Writer]
	WriterCheck  = conform.Check[Writer]
	WriterOption = conform.Option[Writer]
)

// WriterChecks contains the name of each callable of Writer, with the
// identities of its checks.
var WriterChecks = struct {
	Put writerPut
	Get writerGet
}{
	Put: writerPut{"Writer.Put"},
	Get: writerGet{"Writer.Get"},
}

// writerPut is the name of Writer.Put, with the identities of its checks.
type writerPut struct {
	conform.Callable[Writer]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c writerPut) Smoke() conform.ID[Writer] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c writerPut) SmokeGenerated() conform.ID[Writer] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c writerPut) Cancel() conform.ID[Writer] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c writerPut) CancelGenerated() conform.ID[Writer] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Put returns an
// error when a call of its dependency fails.
func (c writerPut) FailureIsReported() conform.ID[Writer] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Put
// returns an error when a call of its dependency fails.
func (c writerPut) FailureIsReportedGenerated() conform.ID[Writer] {
	return c.Check("failure-is-reported#generated")
}

// NoPartialWrite returns the identity of the check that Get finds no
// effect of a Put that failed with its dependency.
func (c writerPut) NoPartialWrite() conform.ID[Writer] {
	return c.Check("no-partial-write")
}

// NoPartialWriteGenerated returns the identity of the check that Get finds
// no effect of a Put that failed with its dependency.
func (c writerPut) NoPartialWriteGenerated() conform.ID[Writer] {
	return c.Check("no-partial-write#generated")
}

// writerGet is the name of Writer.Get, with the identities of its checks.
type writerGet struct {
	conform.Callable[Writer]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c writerGet) Smoke() conform.ID[Writer] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c writerGet) SmokeGenerated() conform.ID[Writer] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c writerGet) Cancel() conform.ID[Writer] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c writerGet) CancelGenerated() conform.ID[Writer] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Get returns an
// error when a call of its dependency fails.
func (c writerGet) FailureIsReported() conform.ID[Writer] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Get
// returns an error when a call of its dependency fails.
func (c writerGet) FailureIsReportedGenerated() conform.ID[Writer] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c writerGet) ZeroOnError() conform.ID[Writer] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c writerGet) ZeroOnErrorGenerated() conform.ID[Writer] {
	return c.Check("zero-on-error#generated")
}

// WriterSuite returns the checks of Writer, and the calls and the budgets
// of its callables.
func WriterSuite() conform.Suite[Writer] {
	return conform.Suite[Writer]{
		Declaration: "Writer",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Writer]{
			{
				ID:          WriterChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WriterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: writerPutPanics,
			},
			{
				ID:          WriterChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnV)
					})
				},
				Defect: writerPutPanics,
			},
			{
				ID:          WriterChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WriterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: writerPutIgnoresCancel,
			},
			{
				ID:          WriterChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnV)
					})
				},
				Defect: writerPutIgnoresCancel,
			},
			{
				ID:        WriterChecks.Put.FailureIsReported(),
				Statement: "Put returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *WriterHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: writerPutSwallowsFailure,
			},
			{
				ID:        WriterChecks.Put.FailureIsReportedGenerated(),
				Statement: "Put returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnV)
					})
				},
				Defect: writerPutSwallowsFailure,
			},
			{
				ID:        WriterChecks.Put.NoPartialWrite(),
				Statement: "Get finds no effect of a Put that failed with its dependency",
				Claim:     spec.NoPartialWrite,
				Probe:     spec.FailsWithoutEffect,
				Reason:    spec.PartialWrite,
				Law:       spec.FailedCallLeavesNoEffect,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *WriterHandle) {
					probe.FailsWithoutEffect(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					}, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					})
				},
				Defect: writerGetHalfApplied,
			},
			{
				ID:        WriterChecks.Put.NoPartialWriteGenerated(),
				Statement: "Get finds no effect of a Put that failed with its dependency",
				Claim:     spec.NoPartialWrite,
				Probe:     spec.FailsWithoutEffect,
				Reason:    spec.PartialWrite,
				Law:       spec.FailedCallLeavesNoEffect,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.FailsWithoutEffect(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnV)
					}, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, drawnV.Key)
					})
				},
				Defect: writerGetHalfApplied,
			},
			{
				ID:          WriterChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WriterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: writerGetPanics,
			},
			{
				ID:          WriterChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: writerGetPanics,
			},
			{
				ID:          WriterChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WriterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: writerGetIgnoresCancel,
			},
			{
				ID:          WriterChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: writerGetIgnoresCancel,
			},
			{
				ID:        WriterChecks.Get.FailureIsReported(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *WriterHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: writerGetSwallowsFailure,
			},
			{
				ID:        WriterChecks.Get.FailureIsReportedGenerated(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: writerGetSwallowsFailure,
			},
			{
				ID:          WriterChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WriterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: writerGetValueBeside,
			},
			{
				ID:          WriterChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *WriterHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Value{})
				},
				Defect: writerGetValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Writer.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Writer]{
			{
				Calls:    writerPutCalls,
				Over:     newWriterOver,
				Callable: WriterChecks.Put.Callable,
			},
			{
				Calls:    writerGetCalls,
				Over:     newWriterOver,
				Callable: WriterChecks.Get.Callable,
			},
		},
	}
}

// RunWriterSuite applies the checks of Writer, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunWriterSuite[T Writer](t *testing.T, newSubject func() T, opts ...WriterOption) {
	t.Helper()
	RunWriterSuiteFromHandle(t, func(*WriterHandle) T { return newSubject() }, opts...)
}

// RunWriterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunWriterSuiteFromHandle[T Writer](
	t *testing.T,
	newSubject func(h *WriterHandle) T,
	opts ...WriterOption,
) {
	t.Helper()
	conform.Run(t, WriterSuite(), writerSubjects(newSubject), opts...)
}

// FuzzWriterSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzWriterSuite[T Writer](
	f *testing.F,
	callable conform.Callable[Writer],
	newSubject func(h *WriterHandle) T,
	opts ...WriterOption,
) {
	f.Helper()
	conform.Fuzz(f, WriterSuite(), callable, writerSubjects(newSubject), opts...)
}

// RunWriterBudgets applies the allocation ceilings of Writer and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunWriterBudgets[T Writer](
	t *testing.T,
	newSubject func(h *WriterHandle) T,
	opts ...WriterOption,
) {
	t.Helper()
	conform.RunBudgets(t, WriterSuite(), writerSubjects(newSubject), opts...)
}

// BenchWriterSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchWriterSuite[T Writer](
	b conform.B,
	callable conform.Callable[Writer],
	newSubject func(h *WriterHandle) T,
	opts ...WriterOption,
) {
	b.Helper()
	conform.Bench(b, WriterSuite(), callable, writerSubjects(newSubject), opts...)
}

// writerSubjects returns the factory of the subjects of Writer that calls
// newSubject with the handle of each check.
func writerSubjects[T Writer](newSubject func(h *WriterHandle) T) func(h *WriterHandle) Writer {
	return func(h *WriterHandle) Writer { return newSubject(h) }
}

// writerPutCalls generates the calls of Writer.Put with the inputs that
// each case draws.
var writerPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Writer) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Writer) error {
		return subject.Put(ctx, v)
	}
})

// writerGetCalls generates the calls of Writer.Get with the inputs that
// each case draws.
var writerGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Writer) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Writer) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// writerOverEmbedded is Writer under the name of the field that writerOver
// embeds.
type writerOverEmbedded = Writer

// writerOver is a subject of Writer whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type writerOver struct {
	writerOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o writerOver) Put(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o writerOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newWriterOver returns the planted defect of a ceiling of Writer, whose
// callables call spend.
func newWriterOver(_ *WriterHandle, spend func()) Writer {
	return writerOver{spend: spend}
}

// writerPutPanicsDefect is the planted defect of a Writer whose Put
// panics.
type writerPutPanicsDefect struct {
	writerOverEmbedded
}

// Put is the override of the planted defect.
func (*writerPutPanicsDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Writer whose Put panics")
}

// writerPutPanics returns the planted defect of a Writer whose Put panics
// over a fresh subject.
func writerPutPanics(h *WriterHandle) Writer {
	return &writerPutPanicsDefect{writerOverEmbedded: h.Fresh()}
}

// writerPutIgnoresCancelDefect is the planted defect of a Writer whose Put
// ignores the context.
type writerPutIgnoresCancelDefect struct {
	writerOverEmbedded
}

// Put is the override of the planted defect.
func (d *writerPutIgnoresCancelDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.writerOverEmbedded.Put(p0, p1)
}

// writerPutIgnoresCancel returns the planted defect of a Writer whose Put
// ignores the context over a fresh subject.
func writerPutIgnoresCancel(h *WriterHandle) Writer {
	return &writerPutIgnoresCancelDefect{writerOverEmbedded: h.Fresh()}
}

// writerPutSwallowsFailureDefect is the planted defect of a Writer whose
// Put returns no error when its dependency fails.
type writerPutSwallowsFailureDefect struct {
	writerOverEmbedded
	h *WriterHandle
}

// Put is the override of the planted defect.
func (d *writerPutSwallowsFailureDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.writerOverEmbedded.Put(p0, p1)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// writerPutSwallowsFailure returns the planted defect of a Writer whose
// Put returns no error when its dependency fails over a fresh subject.
func writerPutSwallowsFailure(h *WriterHandle) Writer {
	return &writerPutSwallowsFailureDefect{writerOverEmbedded: h.Fresh(), h: h}
}

// writerGetHalfAppliedDefect is the planted defect of a Writer whose Get
// finds a part of a Put that failed.
type writerGetHalfAppliedDefect struct {
	writerOverEmbedded
	h *WriterHandle
}

// Get is the override of the planted defect.
func (d *writerGetHalfAppliedDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	failed := d.h.Injector().Failures() > 0
	r0, r1 = d.writerOverEmbedded.Get(p0, p1)
	if failed {
		r0 = Value{Key: "Key-b"}
		r1 = nil
	}
	return r0, r1
}

// writerGetHalfApplied returns the planted defect of a Writer whose Get
// finds a part of a Put that failed over a fresh subject.
func writerGetHalfApplied(h *WriterHandle) Writer {
	return &writerGetHalfAppliedDefect{writerOverEmbedded: h.Fresh(), h: h}
}

// writerGetPanicsDefect is the planted defect of a Writer whose Get
// panics.
type writerGetPanicsDefect struct {
	writerOverEmbedded
}

// Get is the override of the planted defect.
func (*writerGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Writer whose Get panics")
}

// writerGetPanics returns the planted defect of a Writer whose Get panics
// over a fresh subject.
func writerGetPanics(h *WriterHandle) Writer {
	return &writerGetPanicsDefect{writerOverEmbedded: h.Fresh()}
}

// writerGetIgnoresCancelDefect is the planted defect of a Writer whose Get
// ignores the context.
type writerGetIgnoresCancelDefect struct {
	writerOverEmbedded
}

// Get is the override of the planted defect.
func (d *writerGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.writerOverEmbedded.Get(p0, p1)
}

// writerGetIgnoresCancel returns the planted defect of a Writer whose Get
// ignores the context over a fresh subject.
func writerGetIgnoresCancel(h *WriterHandle) Writer {
	return &writerGetIgnoresCancelDefect{writerOverEmbedded: h.Fresh()}
}

// writerGetSwallowsFailureDefect is the planted defect of a Writer whose
// Get returns no error when its dependency fails.
type writerGetSwallowsFailureDefect struct {
	writerOverEmbedded
	h *WriterHandle
}

// Get is the override of the planted defect.
func (d *writerGetSwallowsFailureDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.writerOverEmbedded.Get(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(Value)
		r1 = nil
	}
	return r0, r1
}

// writerGetSwallowsFailure returns the planted defect of a Writer whose
// Get returns no error when its dependency fails over a fresh subject.
func writerGetSwallowsFailure(h *WriterHandle) Writer {
	return &writerGetSwallowsFailureDefect{writerOverEmbedded: h.Fresh(), h: h}
}

// writerGetValueBesideDefect is the planted defect of a Writer whose Get
// returns a value beside an error.
type writerGetValueBesideDefect struct {
	writerOverEmbedded
}

// Get is the override of the planted defect.
func (d *writerGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.writerOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// writerGetValueBeside returns the planted defect of a Writer whose Get
// returns a value beside an error over a fresh subject.
func writerGetValueBeside(h *WriterHandle) Writer {
	return &writerGetValueBesideDefect{writerOverEmbedded: h.Fresh()}
}

// The derived checks of Syncer:
//
// Syncer.Sync/smoke: Sync returns without a panic for derived inputs.
//
// Syncer.Sync/cancel: Sync returns an error for a cancelled context.
//
// Syncer.Sync/failure-is-reported: Sync returns an error when a call of
// its dependency fails.
//
// Syncer.Sync/retry-succeeds: Sync succeeds within the declared attempts.
//
// Syncer.Sync/retry-succeeds-under-faults: Sync succeeds within the
// declared attempts after a call of its dependency fails.

// The derived checks of Syncer that witness their laws at the shallow end:
//
// Syncer.Sync/smoke
//
// Syncer.Sync/retry-succeeds

// The unclaimed classifications of Syncer:
//
// The shape:lifecycle of Syncer.Sync has the stance silent. A lifecycle
// call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Syncer writes.
type (
	SyncerHandle = conform.Handle[Syncer]
	SyncerCheck  = conform.Check[Syncer]
	SyncerOption = conform.Option[Syncer]
)

// SyncerChecks contains the name of each callable of Syncer, with the
// identities of its checks.
var SyncerChecks = struct {
	Sync syncerSync
}{
	Sync: syncerSync{"Syncer.Sync"},
}

// syncerSync is the name of Syncer.Sync, with the identities of its
// checks.
type syncerSync struct {
	conform.Callable[Syncer]
}

// Smoke returns the identity of the check that Sync returns without a
// panic for derived inputs.
func (c syncerSync) Smoke() conform.ID[Syncer] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Sync returns an error for
// a cancelled context.
func (c syncerSync) Cancel() conform.ID[Syncer] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Sync returns an
// error when a call of its dependency fails.
func (c syncerSync) FailureIsReported() conform.ID[Syncer] {
	return c.Check("failure-is-reported")
}

// RetrySucceeds returns the identity of the check that Sync succeeds
// within the declared attempts.
func (c syncerSync) RetrySucceeds() conform.ID[Syncer] {
	return c.Check("retry-succeeds")
}

// RetrySucceedsUnderFaults returns the identity of the check that Sync
// succeeds within the declared attempts after a call of its dependency
// fails.
func (c syncerSync) RetrySucceedsUnderFaults() conform.ID[Syncer] {
	return c.Check("retry-succeeds-under-faults")
}

// SyncerSuite returns the checks of Syncer, and the calls and the budgets
// of its callables.
func SyncerSuite() conform.Suite[Syncer] {
	return conform.Suite[Syncer]{
		Declaration: "Syncer",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Syncer]{
			{
				ID:          SyncerChecks.Sync.Smoke(),
				Statement:   "Sync returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SyncerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					})
				},
				Defect: syncerSyncPanics,
			},
			{
				ID:          SyncerChecks.Sync.Cancel(),
				Statement:   "Sync returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SyncerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					})
				},
				Defect: syncerSyncIgnoresCancel,
			},
			{
				ID:        SyncerChecks.Sync.FailureIsReported(),
				Statement: "Sync returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *SyncerHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					})
				},
				Defect: syncerSyncSwallowsFailure,
			},
			{
				ID:          SyncerChecks.Sync.RetrySucceeds(),
				Statement:   "Sync succeeds within the declared attempts",
				Claim:       spec.RetrySucceeds,
				Probe:       spec.SucceedsWithin,
				Reason:      spec.NeverConverged,
				Law:         spec.CallConvergesWithinTheDeclaredAttempts,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SyncerHandle) {
					probe.SucceedsWithin(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					}, 3)
				},
				Defect: syncerSyncFailsAlways,
			},
			{
				ID:        SyncerChecks.Sync.RetrySucceedsUnderFaults(),
				Statement: "Sync succeeds within the declared attempts after a call of its dependency fails",
				Claim:     spec.RetrySucceedsUnderFaults,
				Probe:     spec.SucceedsAfterAFailure,
				Reason:    spec.NeverConverged,
				Law:       spec.CallConvergesWithinTheDeclaredAttempts,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Entails: []conform.ID[Syncer]{SyncerChecks.Sync.RetrySucceeds()},
				Apply: func(h *SyncerHandle) {
					probe.SucceedsAfterAFailure(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					}, 3)
				},
				Defect: syncerSyncGivesUpEarly,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Syncer.Sync",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Syncer]{
			{
				Calls:    syncerSyncCalls,
				Over:     newSyncerOver,
				Callable: SyncerChecks.Sync.Callable,
			},
		},
	}
}

// RunSyncerSuite applies the checks of Syncer, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunSyncerSuite[T Syncer](t *testing.T, newSubject func() T, opts ...SyncerOption) {
	t.Helper()
	RunSyncerSuiteFromHandle(t, func(*SyncerHandle) T { return newSubject() }, opts...)
}

// RunSyncerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunSyncerSuiteFromHandle[T Syncer](
	t *testing.T,
	newSubject func(h *SyncerHandle) T,
	opts ...SyncerOption,
) {
	t.Helper()
	conform.Run(t, SyncerSuite(), syncerSubjects(newSubject), opts...)
}

// FuzzSyncerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzSyncerSuite[T Syncer](
	f *testing.F,
	callable conform.Callable[Syncer],
	newSubject func(h *SyncerHandle) T,
	opts ...SyncerOption,
) {
	f.Helper()
	conform.Fuzz(f, SyncerSuite(), callable, syncerSubjects(newSubject), opts...)
}

// RunSyncerBudgets applies the allocation ceilings of Syncer and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunSyncerBudgets[T Syncer](
	t *testing.T,
	newSubject func(h *SyncerHandle) T,
	opts ...SyncerOption,
) {
	t.Helper()
	conform.RunBudgets(t, SyncerSuite(), syncerSubjects(newSubject), opts...)
}

// BenchSyncerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchSyncerSuite[T Syncer](
	b conform.B,
	callable conform.Callable[Syncer],
	newSubject func(h *SyncerHandle) T,
	opts ...SyncerOption,
) {
	b.Helper()
	conform.Bench(b, SyncerSuite(), callable, syncerSubjects(newSubject), opts...)
}

// syncerSubjects returns the factory of the subjects of Syncer that calls
// newSubject with the handle of each check.
func syncerSubjects[T Syncer](newSubject func(h *SyncerHandle) T) func(h *SyncerHandle) Syncer {
	return func(h *SyncerHandle) Syncer { return newSubject(h) }
}

// syncerSyncCalls generates the calls of Syncer.Sync with the inputs that
// each case draws.
var syncerSyncCalls = prop.Composite(func(*prop.Case) func(context.Context, Syncer) error {
	return func(ctx context.Context, subject Syncer) error {
		return subject.Sync(ctx)
	}
})

// syncerOverEmbedded is Syncer under the name of the field that syncerOver
// embeds.
type syncerOverEmbedded = Syncer

// syncerOver is a subject of Syncer whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type syncerOver struct {
	syncerOverEmbedded
	spend func()
}

// Sync calls spend and returns its zero results.
func (o syncerOver) Sync(context.Context) (r0 error) {
	o.spend()
	return r0
}

// newSyncerOver returns the planted defect of a ceiling of Syncer, whose
// callables call spend.
func newSyncerOver(_ *SyncerHandle, spend func()) Syncer {
	return syncerOver{spend: spend}
}

// syncerSyncPanicsDefect is the planted defect of a Syncer whose Sync
// panics.
type syncerSyncPanicsDefect struct {
	syncerOverEmbedded
}

// Sync is the override of the planted defect.
func (*syncerSyncPanicsDefect) Sync(p0 context.Context) (r0 error) {
	panic("the planted defect of a Syncer whose Sync panics")
}

// syncerSyncPanics returns the planted defect of a Syncer whose Sync
// panics over a fresh subject.
func syncerSyncPanics(h *SyncerHandle) Syncer {
	return &syncerSyncPanicsDefect{syncerOverEmbedded: h.Fresh()}
}

// syncerSyncIgnoresCancelDefect is the planted defect of a Syncer whose
// Sync ignores the context.
type syncerSyncIgnoresCancelDefect struct {
	syncerOverEmbedded
}

// Sync is the override of the planted defect.
func (d *syncerSyncIgnoresCancelDefect) Sync(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.syncerOverEmbedded.Sync(p0)
}

// syncerSyncIgnoresCancel returns the planted defect of a Syncer whose
// Sync ignores the context over a fresh subject.
func syncerSyncIgnoresCancel(h *SyncerHandle) Syncer {
	return &syncerSyncIgnoresCancelDefect{syncerOverEmbedded: h.Fresh()}
}

// syncerSyncSwallowsFailureDefect is the planted defect of a Syncer whose
// Sync returns no error when its dependency fails.
type syncerSyncSwallowsFailureDefect struct {
	syncerOverEmbedded
	h *SyncerHandle
}

// Sync is the override of the planted defect.
func (d *syncerSyncSwallowsFailureDefect) Sync(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.syncerOverEmbedded.Sync(p0)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// syncerSyncSwallowsFailure returns the planted defect of a Syncer whose
// Sync returns no error when its dependency fails over a fresh subject.
func syncerSyncSwallowsFailure(h *SyncerHandle) Syncer {
	return &syncerSyncSwallowsFailureDefect{syncerOverEmbedded: h.Fresh(), h: h}
}

// syncerSyncFailsAlwaysDefect is the planted defect of a Syncer whose Sync
// never succeeds.
type syncerSyncFailsAlwaysDefect struct {
	syncerOverEmbedded
}

// Sync is the override of the planted defect.
func (d *syncerSyncFailsAlwaysDefect) Sync(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	r0 = d.syncerOverEmbedded.Sync(p0)
	if r0 == nil {
		r0 = conform.ErrGenerated
	}
	return r0
}

// syncerSyncFailsAlways returns the planted defect of a Syncer whose Sync
// never succeeds over a fresh subject.
func syncerSyncFailsAlways(h *SyncerHandle) Syncer {
	return &syncerSyncFailsAlwaysDefect{syncerOverEmbedded: h.Fresh()}
}

// syncerSyncGivesUpEarlyDefect is the planted defect of a Syncer whose
// Sync keeps failing once its dependency failed.
type syncerSyncGivesUpEarlyDefect struct {
	syncerOverEmbedded
	h *SyncerHandle
}

// Sync is the override of the planted defect.
func (d *syncerSyncGivesUpEarlyDefect) Sync(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	r0 = d.syncerOverEmbedded.Sync(p0)
	if d.h.Injector().Failures() > 0 {
		r0 = conform.ErrGenerated
	}
	return r0
}

// syncerSyncGivesUpEarly returns the planted defect of a Syncer whose Sync
// keeps failing once its dependency failed over a fresh subject.
func syncerSyncGivesUpEarly(h *SyncerHandle) Syncer {
	return &syncerSyncGivesUpEarlyDefect{syncerOverEmbedded: h.Fresh(), h: h}
}

// The derived checks of Account:
//
// Account.Run/smoke: Run returns without a panic for derived inputs.
//
// Account.Run/smoke#generated: Run returns without a panic for derived
// inputs.
//
// Account.Run/cancel: Run returns an error for a cancelled context.
//
// Account.Run/cancel#generated: Run returns an error for a cancelled
// context.
//
// Account.Run/failure-is-reported: Run returns an error when a call of its
// dependency fails.
//
// Account.Run/failure-is-reported#generated: Run returns an error when a
// call of its dependency fails.
//
// Account.Run/transaction: Get finds nothing that a failed Run wrote
// through Put.
//
// Account.Run/transaction#generated: Get finds nothing that a failed Run
// wrote through Put.
//
// Account.Run/transaction-under-faults: Get finds nothing that Run wrote
// through Put before its commit failed.
//
// Account.Run/transaction-under-faults#generated: Get finds nothing that
// Run wrote through Put before its commit failed.
//
// Account.Put/smoke: Put returns without a panic for derived inputs.
//
// Account.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Account.Put/cancel: Put returns an error for a cancelled context.
//
// Account.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Account.Put/failure-is-reported: Put returns an error when a call of its
// dependency fails.
//
// Account.Put/failure-is-reported#generated: Put returns an error when a
// call of its dependency fails.
//
// Account.Get/smoke: Get returns without a panic for derived inputs.
//
// Account.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Account.Get/cancel: Get returns an error for a cancelled context.
//
// Account.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Account.Get/failure-is-reported: Get returns an error when a call of its
// dependency fails.
//
// Account.Get/failure-is-reported#generated: Get returns an error when a
// call of its dependency fails.
//
// Account.Get/zero-on-error: Get returns the zero value beside an error.
//
// Account.Get/zero-on-error#generated: Get returns the zero value beside
// an error.

// The derived checks of Account that witness their laws at the shallow
// end:
//
// Account.Run/smoke
//
// Account.Put/smoke
//
// Account.Get/smoke

// The unclaimed classifications of Account:
//
// The shape:writer of Account.Run has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Account.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Account writes.
type (
	AccountHandle = conform.Handle[Account]
	AccountCheck  = conform.Check[Account]
	AccountOption = conform.Option[Account]
)

// AccountChecks contains the name of each callable of Account, with the
// identities of its checks.
var AccountChecks = struct {
	Run accountRun
	Put accountPut
	Get accountGet
}{
	Run: accountRun{"Account.Run"},
	Put: accountPut{"Account.Put"},
	Get: accountGet{"Account.Get"},
}

// accountRun is the name of Account.Run, with the identities of its
// checks.
type accountRun struct {
	conform.Callable[Account]
}

// Smoke returns the identity of the check that Run returns without a panic
// for derived inputs.
func (c accountRun) Smoke() conform.ID[Account] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Run returns
// without a panic for derived inputs.
func (c accountRun) SmokeGenerated() conform.ID[Account] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Run returns an error for a
// cancelled context.
func (c accountRun) Cancel() conform.ID[Account] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Run returns an
// error for a cancelled context.
func (c accountRun) CancelGenerated() conform.ID[Account] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Run returns an
// error when a call of its dependency fails.
func (c accountRun) FailureIsReported() conform.ID[Account] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Run
// returns an error when a call of its dependency fails.
func (c accountRun) FailureIsReportedGenerated() conform.ID[Account] {
	return c.Check("failure-is-reported#generated")
}

// Transaction returns the identity of the check that Get finds nothing
// that a failed Run wrote through Put.
func (c accountRun) Transaction() conform.ID[Account] {
	return c.Check("transaction")
}

// TransactionGenerated returns the identity of the check that Get finds
// nothing that a failed Run wrote through Put.
func (c accountRun) TransactionGenerated() conform.ID[Account] {
	return c.Check("transaction#generated")
}

// TransactionUnderFaults returns the identity of the check that Get finds
// nothing that Run wrote through Put before its commit failed.
func (c accountRun) TransactionUnderFaults() conform.ID[Account] {
	return c.Check("transaction-under-faults")
}

// TransactionUnderFaultsGenerated returns the identity of the check that
// Get finds nothing that Run wrote through Put before its commit failed.
func (c accountRun) TransactionUnderFaultsGenerated() conform.ID[Account] {
	return c.Check("transaction-under-faults#generated")
}

// accountPut is the name of Account.Put, with the identities of its
// checks.
type accountPut struct {
	conform.Callable[Account]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c accountPut) Smoke() conform.ID[Account] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c accountPut) SmokeGenerated() conform.ID[Account] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c accountPut) Cancel() conform.ID[Account] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c accountPut) CancelGenerated() conform.ID[Account] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Put returns an
// error when a call of its dependency fails.
func (c accountPut) FailureIsReported() conform.ID[Account] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Put
// returns an error when a call of its dependency fails.
func (c accountPut) FailureIsReportedGenerated() conform.ID[Account] {
	return c.Check("failure-is-reported#generated")
}

// accountGet is the name of Account.Get, with the identities of its
// checks.
type accountGet struct {
	conform.Callable[Account]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c accountGet) Smoke() conform.ID[Account] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c accountGet) SmokeGenerated() conform.ID[Account] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c accountGet) Cancel() conform.ID[Account] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c accountGet) CancelGenerated() conform.ID[Account] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Get returns an
// error when a call of its dependency fails.
func (c accountGet) FailureIsReported() conform.ID[Account] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Get
// returns an error when a call of its dependency fails.
func (c accountGet) FailureIsReportedGenerated() conform.ID[Account] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c accountGet) ZeroOnError() conform.ID[Account] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c accountGet) ZeroOnErrorGenerated() conform.ID[Account] {
	return c.Check("zero-on-error#generated")
}

// AccountSuite returns the checks of Account, and the calls and the
// budgets of its callables.
func AccountSuite() conform.Suite[Account] {
	return conform.Suite[Account]{
		Declaration: "Account",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Account]{
			{
				ID:          AccountChecks.Run.Smoke(),
				Statement:   "Run returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, func(ctx context.Context) error { return nil })
					})
				},
				Defect: accountRunPanics,
			},
			{
				ID:          AccountChecks.Run.SmokeGenerated(),
				Statement:   "Run returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					bodyResults := c.Draw(
						prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
						"body",
					)
					drawnBody := func() func(ctx context.Context) error {
						var bodyAt atomic.Int64
						return func(context.Context) (r0 error) {
							i := int(bodyAt.Add(1) - 1)
							if i >= len(bodyResults) {
								r0 = conform.ErrGenerated
								return r0
							}
							return bodyResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, drawnBody())
					})
				},
				Defect: accountRunPanics,
			},
			{
				ID:          AccountChecks.Run.Cancel(),
				Statement:   "Run returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, func(ctx context.Context) error { return nil })
					})
				},
				Defect: accountRunIgnoresCancel,
			},
			{
				ID:          AccountChecks.Run.CancelGenerated(),
				Statement:   "Run returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					bodyResults := c.Draw(
						prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
						"body",
					)
					drawnBody := func() func(ctx context.Context) error {
						var bodyAt atomic.Int64
						return func(context.Context) (r0 error) {
							i := int(bodyAt.Add(1) - 1)
							if i >= len(bodyResults) {
								r0 = conform.ErrGenerated
								return r0
							}
							return bodyResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, drawnBody())
					})
				},
				Defect: accountRunIgnoresCancel,
			},
			{
				ID:        AccountChecks.Run.FailureIsReported(),
				Statement: "Run returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *AccountHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, func(ctx context.Context) error { return nil })
					})
				},
				Defect: accountRunSwallowsFailure,
			},
			{
				ID:        AccountChecks.Run.FailureIsReportedGenerated(),
				Statement: "Run returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *AccountHandle, c *prop.Case) {
					bodyResults := c.Draw(
						prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
						"body",
					)
					drawnBody := func() func(ctx context.Context) error {
						var bodyAt atomic.Int64
						return func(context.Context) (r0 error) {
							i := int(bodyAt.Add(1) - 1)
							if i >= len(bodyResults) {
								r0 = conform.ErrGenerated
								return r0
							}
							return bodyResults[i]
						}
					}
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, drawnBody())
					})
				},
				Defect: accountRunSwallowsFailure,
			},
			{
				ID:          AccountChecks.Run.Transaction(),
				Statement:   "Get finds nothing that a failed Run wrote through Put",
				Claim:       spec.Transaction,
				Probe:       spec.RollsBackOnFailure,
				Reason:      spec.NotRolledBack,
				Law:         spec.FailedScopeLeavesNothing,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountHandle) {
					probe.RollsBackOnFailure(h, func(ctx context.Context, body probe.Act) error {
						return h.Subject().Run(ctx, body)
					}, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					}, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					}, ErrMissing)
				},
				Defect: accountGetKeepsFailedWrites,
			},
			{
				ID:          AccountChecks.Run.TransactionGenerated(),
				Statement:   "Get finds nothing that a failed Run wrote through Put",
				Claim:       spec.Transaction,
				Probe:       spec.RollsBackOnFailure,
				Reason:      spec.NotRolledBack,
				Law:         spec.FailedScopeLeavesNothing,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.RollsBackOnFailure(h, func(ctx context.Context, body probe.Act) error {
						return h.Subject().Run(ctx, body)
					}, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnKey, drawnV)
					}, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					}, ErrMissing)
				},
				Defect: accountGetKeepsFailedWrites,
			},
			{
				ID:        AccountChecks.Run.TransactionUnderFaults(),
				Statement: "Get finds nothing that Run wrote through Put before its commit failed",
				Claim:     spec.TransactionUnderFaults,
				Probe:     spec.RollsBackUnderFaults,
				Reason:    spec.NotRolledBack,
				Law:       spec.FailedCommitLeavesNothing,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *AccountHandle) {
					probe.RollsBackUnderFaults(h, func(ctx context.Context, body probe.Act) error {
						return h.Subject().Run(ctx, body)
					}, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					}, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					}, ErrMissing)
				},
				Defect: accountGetCommitsOnFailure,
			},
			{
				ID:        AccountChecks.Run.TransactionUnderFaultsGenerated(),
				Statement: "Get finds nothing that Run wrote through Put before its commit failed",
				Claim:     spec.TransactionUnderFaults,
				Probe:     spec.RollsBackUnderFaults,
				Reason:    spec.NotRolledBack,
				Law:       spec.FailedCommitLeavesNothing,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.RollsBackUnderFaults(h, func(ctx context.Context, body probe.Act) error {
						return h.Subject().Run(ctx, body)
					}, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnKey, drawnV)
					}, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					}, ErrMissing)
				},
				Defect: accountGetCommitsOnFailure,
			},
			{
				ID:          AccountChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: accountPutPanics,
			},
			{
				ID:          AccountChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnKey, drawnV)
					})
				},
				Defect: accountPutPanics,
			},
			{
				ID:          AccountChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: accountPutIgnoresCancel,
			},
			{
				ID:          AccountChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnKey, drawnV)
					})
				},
				Defect: accountPutIgnoresCancel,
			},
			{
				ID:        AccountChecks.Put.FailureIsReported(),
				Statement: "Put returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *AccountHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: accountPutSwallowsFailure,
			},
			{
				ID:        AccountChecks.Put.FailureIsReportedGenerated(),
				Statement: "Put returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnKey, drawnV)
					})
				},
				Defect: accountPutSwallowsFailure,
			},
			{
				ID:          AccountChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AccountHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: accountGetPanics,
			},
			{
				ID:          AccountChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: accountGetPanics,
			},
			{
				ID:          AccountChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: accountGetIgnoresCancel,
			},
			{
				ID:          AccountChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: accountGetIgnoresCancel,
			},
			{
				ID:        AccountChecks.Get.FailureIsReported(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *AccountHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: accountGetSwallowsFailure,
			},
			{
				ID:        AccountChecks.Get.FailureIsReportedGenerated(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: accountGetSwallowsFailure,
			},
			{
				ID:          AccountChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AccountHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: accountGetValueBeside,
			},
			{
				ID:          AccountChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *AccountHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Value{})
				},
				Defect: accountGetValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Account.Run",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Account.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Account]{
			{
				Calls:    accountRunCalls,
				Over:     newAccountOver,
				Callable: AccountChecks.Run.Callable,
			},
			{
				Calls:    accountPutCalls,
				Over:     newAccountOver,
				Callable: AccountChecks.Put.Callable,
			},
			{
				Calls:    accountGetCalls,
				Over:     newAccountOver,
				Callable: AccountChecks.Get.Callable,
			},
		},
	}
}

// RunAccountSuite applies the checks of Account, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunAccountSuite[T Account](t *testing.T, newSubject func() T, opts ...AccountOption) {
	t.Helper()
	RunAccountSuiteFromHandle(t, func(*AccountHandle) T { return newSubject() }, opts...)
}

// RunAccountSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunAccountSuiteFromHandle[T Account](
	t *testing.T,
	newSubject func(h *AccountHandle) T,
	opts ...AccountOption,
) {
	t.Helper()
	conform.Run(t, AccountSuite(), accountSubjects(newSubject), opts...)
}

// FuzzAccountSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzAccountSuite[T Account](
	f *testing.F,
	callable conform.Callable[Account],
	newSubject func(h *AccountHandle) T,
	opts ...AccountOption,
) {
	f.Helper()
	conform.Fuzz(f, AccountSuite(), callable, accountSubjects(newSubject), opts...)
}

// RunAccountBudgets applies the allocation ceilings of Account and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunAccountBudgets[T Account](
	t *testing.T,
	newSubject func(h *AccountHandle) T,
	opts ...AccountOption,
) {
	t.Helper()
	conform.RunBudgets(t, AccountSuite(), accountSubjects(newSubject), opts...)
}

// BenchAccountSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchAccountSuite[T Account](
	b conform.B,
	callable conform.Callable[Account],
	newSubject func(h *AccountHandle) T,
	opts ...AccountOption,
) {
	b.Helper()
	conform.Bench(b, AccountSuite(), callable, accountSubjects(newSubject), opts...)
}

// accountSubjects returns the factory of the subjects of Account that
// calls newSubject with the handle of each check.
func accountSubjects[T Account](
	newSubject func(h *AccountHandle) T,
) func(h *AccountHandle) Account {
	return func(h *AccountHandle) Account { return newSubject(h) }
}

// accountRunCalls generates the calls of Account.Run with the inputs that
// each case draws.
var accountRunCalls = prop.Composite(func(c *prop.Case) func(context.Context, Account) error {
	bodyResults := c.Draw(prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)), "body")
	var bodyAt atomic.Int64
	body := func(context.Context) (r0 error) {
		i := int(bodyAt.Add(1) - 1)
		if i >= len(bodyResults) {
			r0 = conform.ErrGenerated
			return r0
		}
		return bodyResults[i]
	}
	return func(ctx context.Context, subject Account) error {
		return subject.Run(ctx, body)
	}
})

// accountPutCalls generates the calls of Account.Put with the inputs that
// each case draws.
var accountPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Account) error {
	key := c.Draw(prop.Of[string](), "key")
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Account) error {
		return subject.Put(ctx, key, v)
	}
})

// accountGetCalls generates the calls of Account.Get with the inputs that
// each case draws.
var accountGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Account) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Account) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// accountOverEmbedded is Account under the name of the field that
// accountOver embeds.
type accountOverEmbedded = Account

// accountOver is a subject of Account whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type accountOver struct {
	accountOverEmbedded
	spend func()
}

// Run calls spend and returns its zero results.
func (o accountOver) Run(context.Context, func(ctx context.Context) error) (r0 error) {
	o.spend()
	return r0
}

// Put calls spend and returns its zero results.
func (o accountOver) Put(context.Context, string, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o accountOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newAccountOver returns the planted defect of a ceiling of Account, whose
// callables call spend.
func newAccountOver(_ *AccountHandle, spend func()) Account {
	return accountOver{spend: spend}
}

// accountRunPanicsDefect is the planted defect of a Account whose Run
// panics.
type accountRunPanicsDefect struct {
	accountOverEmbedded
}

// Run is the override of the planted defect.
func (*accountRunPanicsDefect) Run(
	p0 context.Context,
	p1 func(ctx context.Context) error,
) (r0 error) {
	panic("the planted defect of a Account whose Run panics")
}

// accountRunPanics returns the planted defect of a Account whose Run
// panics over a fresh subject.
func accountRunPanics(h *AccountHandle) Account {
	return &accountRunPanicsDefect{accountOverEmbedded: h.Fresh()}
}

// accountRunIgnoresCancelDefect is the planted defect of a Account whose
// Run ignores the context.
type accountRunIgnoresCancelDefect struct {
	accountOverEmbedded
}

// Run is the override of the planted defect.
func (d *accountRunIgnoresCancelDefect) Run(
	p0 context.Context,
	p1 func(ctx context.Context) error,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.accountOverEmbedded.Run(p0, p1)
}

// accountRunIgnoresCancel returns the planted defect of a Account whose
// Run ignores the context over a fresh subject.
func accountRunIgnoresCancel(h *AccountHandle) Account {
	return &accountRunIgnoresCancelDefect{accountOverEmbedded: h.Fresh()}
}

// accountRunSwallowsFailureDefect is the planted defect of a Account whose
// Run returns no error when its dependency fails.
type accountRunSwallowsFailureDefect struct {
	accountOverEmbedded
	h *AccountHandle
}

// Run is the override of the planted defect.
func (d *accountRunSwallowsFailureDefect) Run(
	p0 context.Context,
	p1 func(ctx context.Context) error,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.accountOverEmbedded.Run(p0, p1)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// accountRunSwallowsFailure returns the planted defect of a Account whose
// Run returns no error when its dependency fails over a fresh subject.
func accountRunSwallowsFailure(h *AccountHandle) Account {
	return &accountRunSwallowsFailureDefect{accountOverEmbedded: h.Fresh(), h: h}
}

// accountGetKeepsFailedWritesDefect is the planted defect of a Account
// whose Run keeps what a failed body wrote.
type accountGetKeepsFailedWritesDefect struct {
	accountOverEmbedded
	h *AccountHandle
}

// Get is the override of the planted defect.
func (d *accountGetKeepsFailedWritesDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.h.Injector().Failures() == 0 {
		return r0, r1
	}
	return d.accountOverEmbedded.Get(p0, p1)
}

// accountGetKeepsFailedWrites returns the planted defect of a Account
// whose Run keeps what a failed body wrote over a fresh subject.
func accountGetKeepsFailedWrites(h *AccountHandle) Account {
	return &accountGetKeepsFailedWritesDefect{accountOverEmbedded: h.Fresh(), h: h}
}

// accountGetCommitsOnFailureDefect is the planted defect of a Account
// whose Get finds what Run wrote before its commit failed.
type accountGetCommitsOnFailureDefect struct {
	accountOverEmbedded
	h *AccountHandle
}

// Get is the override of the planted defect.
func (d *accountGetCommitsOnFailureDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	failed := d.h.Injector().Failures() > 0
	r0, r1 = d.accountOverEmbedded.Get(p0, p1)
	if failed {
		r0 = Value{Key: "Key-b"}
		r1 = nil
	}
	return r0, r1
}

// accountGetCommitsOnFailure returns the planted defect of a Account whose
// Get finds what Run wrote before its commit failed over a fresh subject.
func accountGetCommitsOnFailure(h *AccountHandle) Account {
	return &accountGetCommitsOnFailureDefect{accountOverEmbedded: h.Fresh(), h: h}
}

// accountPutPanicsDefect is the planted defect of a Account whose Put
// panics.
type accountPutPanicsDefect struct {
	accountOverEmbedded
}

// Put is the override of the planted defect.
func (*accountPutPanicsDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	panic("the planted defect of a Account whose Put panics")
}

// accountPutPanics returns the planted defect of a Account whose Put
// panics over a fresh subject.
func accountPutPanics(h *AccountHandle) Account {
	return &accountPutPanicsDefect{accountOverEmbedded: h.Fresh()}
}

// accountPutIgnoresCancelDefect is the planted defect of a Account whose
// Put ignores the context.
type accountPutIgnoresCancelDefect struct {
	accountOverEmbedded
}

// Put is the override of the planted defect.
func (d *accountPutIgnoresCancelDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.accountOverEmbedded.Put(p0, p1, p2)
}

// accountPutIgnoresCancel returns the planted defect of a Account whose
// Put ignores the context over a fresh subject.
func accountPutIgnoresCancel(h *AccountHandle) Account {
	return &accountPutIgnoresCancelDefect{accountOverEmbedded: h.Fresh()}
}

// accountPutSwallowsFailureDefect is the planted defect of a Account whose
// Put returns no error when its dependency fails.
type accountPutSwallowsFailureDefect struct {
	accountOverEmbedded
	h *AccountHandle
}

// Put is the override of the planted defect.
func (d *accountPutSwallowsFailureDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.accountOverEmbedded.Put(p0, p1, p2)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// accountPutSwallowsFailure returns the planted defect of a Account whose
// Put returns no error when its dependency fails over a fresh subject.
func accountPutSwallowsFailure(h *AccountHandle) Account {
	return &accountPutSwallowsFailureDefect{accountOverEmbedded: h.Fresh(), h: h}
}

// accountGetPanicsDefect is the planted defect of a Account whose Get
// panics.
type accountGetPanicsDefect struct {
	accountOverEmbedded
}

// Get is the override of the planted defect.
func (*accountGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Account whose Get panics")
}

// accountGetPanics returns the planted defect of a Account whose Get
// panics over a fresh subject.
func accountGetPanics(h *AccountHandle) Account {
	return &accountGetPanicsDefect{accountOverEmbedded: h.Fresh()}
}

// accountGetIgnoresCancelDefect is the planted defect of a Account whose
// Get ignores the context.
type accountGetIgnoresCancelDefect struct {
	accountOverEmbedded
}

// Get is the override of the planted defect.
func (d *accountGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.accountOverEmbedded.Get(p0, p1)
}

// accountGetIgnoresCancel returns the planted defect of a Account whose
// Get ignores the context over a fresh subject.
func accountGetIgnoresCancel(h *AccountHandle) Account {
	return &accountGetIgnoresCancelDefect{accountOverEmbedded: h.Fresh()}
}

// accountGetSwallowsFailureDefect is the planted defect of a Account whose
// Get returns no error when its dependency fails.
type accountGetSwallowsFailureDefect struct {
	accountOverEmbedded
	h *AccountHandle
}

// Get is the override of the planted defect.
func (d *accountGetSwallowsFailureDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.accountOverEmbedded.Get(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(Value)
		r1 = nil
	}
	return r0, r1
}

// accountGetSwallowsFailure returns the planted defect of a Account whose
// Get returns no error when its dependency fails over a fresh subject.
func accountGetSwallowsFailure(h *AccountHandle) Account {
	return &accountGetSwallowsFailureDefect{accountOverEmbedded: h.Fresh(), h: h}
}

// accountGetValueBesideDefect is the planted defect of a Account whose Get
// returns a value beside an error.
type accountGetValueBesideDefect struct {
	accountOverEmbedded
}

// Get is the override of the planted defect.
func (d *accountGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.accountOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// accountGetValueBeside returns the planted defect of a Account whose Get
// returns a value beside an error over a fresh subject.
func accountGetValueBeside(h *AccountHandle) Account {
	return &accountGetValueBesideDefect{accountOverEmbedded: h.Fresh()}
}

// The derived checks of Breaker:
//
// Breaker.Do/smoke: Do returns without a panic for derived inputs.
//
// Breaker.Do/cancel: Do returns an error for a cancelled context.
//
// Breaker.Do/failure-is-reported: Do returns an error when a call of its
// dependency fails.
//
// Breaker.Do/fails-fast: Do refuses a call without its operation after the
// operation failed repeatedly.

// The derived checks of Breaker that witness their laws at the shallow
// end:
//
// Breaker.Do/smoke

// The unclaimed classifications of Breaker:
//
// The shape:lifecycle of Breaker.Do has the stance silent. A lifecycle
// call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Breaker writes.
type (
	BreakerHandle = conform.Handle[Breaker]
	BreakerCheck  = conform.Check[Breaker]
	BreakerOption = conform.Option[Breaker]
)

// BreakerChecks contains the name of each callable of Breaker, with the
// identities of its checks.
var BreakerChecks = struct {
	Do breakerDo
}{
	Do: breakerDo{"Breaker.Do"},
}

// breakerDo is the name of Breaker.Do, with the identities of its checks.
type breakerDo struct {
	conform.Callable[Breaker]
}

// Smoke returns the identity of the check that Do returns without a panic
// for derived inputs.
func (c breakerDo) Smoke() conform.ID[Breaker] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Do returns an error for a
// cancelled context.
func (c breakerDo) Cancel() conform.ID[Breaker] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Do returns an
// error when a call of its dependency fails.
func (c breakerDo) FailureIsReported() conform.ID[Breaker] {
	return c.Check("failure-is-reported")
}

// FailsFast returns the identity of the check that Do refuses a call
// without its operation after the operation failed repeatedly.
func (c breakerDo) FailsFast() conform.ID[Breaker] {
	return c.Check("fails-fast")
}

// BreakerSuite returns the checks of Breaker, and the calls and the
// budgets of its callables.
func BreakerSuite() conform.Suite[Breaker] {
	return conform.Suite[Breaker]{
		Declaration: "Breaker",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Breaker]{
			{
				ID:          BreakerChecks.Do.Smoke(),
				Statement:   "Do returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BreakerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Do(ctx)
					})
				},
				Defect: breakerDoPanics,
			},
			{
				ID:          BreakerChecks.Do.Cancel(),
				Statement:   "Do returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BreakerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Do(ctx)
					})
				},
				Defect: breakerDoIgnoresCancel,
			},
			{
				ID:        BreakerChecks.Do.FailureIsReported(),
				Statement: "Do returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *BreakerHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Do(ctx)
					})
				},
				Defect: breakerDoSwallowsFailure,
			},
			{
				ID:        BreakerChecks.Do.FailsFast(),
				Statement: "Do refuses a call without its operation after the operation failed repeatedly",
				Claim:     spec.FailsFast,
				Probe:     spec.RefusesWithoutTheDependency,
				Reason:    spec.NeverOpened,
				Law:       spec.OpenBreakerSkipsTheOperation,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *BreakerHandle) {
					probe.RefusesWithoutTheDependency(h, func(ctx context.Context) error {
						return h.Subject().Do(ctx)
					})
				},
				Defect: breakerDoNeverOpens,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Breaker.Do",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Breaker]{
			{
				Calls:    breakerDoCalls,
				Over:     newBreakerOver,
				Callable: BreakerChecks.Do.Callable,
			},
		},
	}
}

// RunBreakerSuite applies the checks of Breaker, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunBreakerSuite[T Breaker](t *testing.T, newSubject func() T, opts ...BreakerOption) {
	t.Helper()
	RunBreakerSuiteFromHandle(t, func(*BreakerHandle) T { return newSubject() }, opts...)
}

// RunBreakerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunBreakerSuiteFromHandle[T Breaker](
	t *testing.T,
	newSubject func(h *BreakerHandle) T,
	opts ...BreakerOption,
) {
	t.Helper()
	conform.Run(t, BreakerSuite(), breakerSubjects(newSubject), opts...)
}

// FuzzBreakerSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzBreakerSuite[T Breaker](
	f *testing.F,
	callable conform.Callable[Breaker],
	newSubject func(h *BreakerHandle) T,
	opts ...BreakerOption,
) {
	f.Helper()
	conform.Fuzz(f, BreakerSuite(), callable, breakerSubjects(newSubject), opts...)
}

// RunBreakerBudgets applies the allocation ceilings of Breaker and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunBreakerBudgets[T Breaker](
	t *testing.T,
	newSubject func(h *BreakerHandle) T,
	opts ...BreakerOption,
) {
	t.Helper()
	conform.RunBudgets(t, BreakerSuite(), breakerSubjects(newSubject), opts...)
}

// BenchBreakerSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchBreakerSuite[T Breaker](
	b conform.B,
	callable conform.Callable[Breaker],
	newSubject func(h *BreakerHandle) T,
	opts ...BreakerOption,
) {
	b.Helper()
	conform.Bench(b, BreakerSuite(), callable, breakerSubjects(newSubject), opts...)
}

// breakerSubjects returns the factory of the subjects of Breaker that
// calls newSubject with the handle of each check.
func breakerSubjects[T Breaker](
	newSubject func(h *BreakerHandle) T,
) func(h *BreakerHandle) Breaker {
	return func(h *BreakerHandle) Breaker { return newSubject(h) }
}

// breakerDoCalls generates the calls of Breaker.Do with the inputs that
// each case draws.
var breakerDoCalls = prop.Composite(func(*prop.Case) func(context.Context, Breaker) error {
	return func(ctx context.Context, subject Breaker) error {
		return subject.Do(ctx)
	}
})

// breakerOverEmbedded is Breaker under the name of the field that
// breakerOver embeds.
type breakerOverEmbedded = Breaker

// breakerOver is a subject of Breaker whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type breakerOver struct {
	breakerOverEmbedded
	spend func()
}

// Do calls spend and returns its zero results.
func (o breakerOver) Do(context.Context) (r0 error) {
	o.spend()
	return r0
}

// newBreakerOver returns the planted defect of a ceiling of Breaker, whose
// callables call spend.
func newBreakerOver(_ *BreakerHandle, spend func()) Breaker {
	return breakerOver{spend: spend}
}

// breakerDoPanicsDefect is the planted defect of a Breaker whose Do
// panics.
type breakerDoPanicsDefect struct {
	breakerOverEmbedded
}

// Do is the override of the planted defect.
func (*breakerDoPanicsDefect) Do(p0 context.Context) (r0 error) {
	panic("the planted defect of a Breaker whose Do panics")
}

// breakerDoPanics returns the planted defect of a Breaker whose Do panics
// over a fresh subject.
func breakerDoPanics(h *BreakerHandle) Breaker {
	return &breakerDoPanicsDefect{breakerOverEmbedded: h.Fresh()}
}

// breakerDoIgnoresCancelDefect is the planted defect of a Breaker whose Do
// ignores the context.
type breakerDoIgnoresCancelDefect struct {
	breakerOverEmbedded
}

// Do is the override of the planted defect.
func (d *breakerDoIgnoresCancelDefect) Do(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.breakerOverEmbedded.Do(p0)
}

// breakerDoIgnoresCancel returns the planted defect of a Breaker whose Do
// ignores the context over a fresh subject.
func breakerDoIgnoresCancel(h *BreakerHandle) Breaker {
	return &breakerDoIgnoresCancelDefect{breakerOverEmbedded: h.Fresh()}
}

// breakerDoSwallowsFailureDefect is the planted defect of a Breaker whose
// Do returns no error when its dependency fails.
type breakerDoSwallowsFailureDefect struct {
	breakerOverEmbedded
	h *BreakerHandle
}

// Do is the override of the planted defect.
func (d *breakerDoSwallowsFailureDefect) Do(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.breakerOverEmbedded.Do(p0)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// breakerDoSwallowsFailure returns the planted defect of a Breaker whose
// Do returns no error when its dependency fails over a fresh subject.
func breakerDoSwallowsFailure(h *BreakerHandle) Breaker {
	return &breakerDoSwallowsFailureDefect{breakerOverEmbedded: h.Fresh(), h: h}
}

// breakerDoNeverOpensDefect is the planted defect of a Breaker whose Do
// keeps calling an operation that keeps failing.
type breakerDoNeverOpensDefect struct {
	breakerOverEmbedded
	h *BreakerHandle
}

// Do is the override of the planted defect.
func (d *breakerDoNeverOpensDefect) Do(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return d.h.Fresh().Do(p0)
}

// breakerDoNeverOpens returns the planted defect of a Breaker whose Do
// keeps calling an operation that keeps failing over a fresh subject.
func breakerDoNeverOpens(h *BreakerHandle) Breaker {
	return &breakerDoNeverOpensDefect{breakerOverEmbedded: h.Fresh(), h: h}
}

// The derived checks of Balancer:
//
// Balancer.Route/smoke: Route returns without a panic for derived inputs.
//
// Balancer.Route/smoke#generated: Route returns without a panic for
// derived inputs.
//
// Balancer.Route/cancel: Route returns an error for a cancelled context.
//
// Balancer.Route/cancel#generated: Route returns an error for a cancelled
// context.
//
// Balancer.Route/failure-is-reported: Route returns an error when a call
// of its dependency fails.
//
// Balancer.Route/failure-is-reported#generated: Route returns an error
// when a call of its dependency fails.
//
// Balancer.Route/keeps-its-instance: Route sends each call with one key to
// one instance.
//
// Balancer.Route/keeps-its-instance#generated: Route sends each call with
// one key to one instance.

// The derived checks of Balancer that witness their laws at the shallow
// end:
//
// Balancer.Route/smoke

// The unclaimed classifications of Balancer:
//
// The shape:writer of Balancer.Route has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Balancer writes.
type (
	BalancerHandle = conform.Handle[Balancer]
	BalancerCheck  = conform.Check[Balancer]
	BalancerOption = conform.Option[Balancer]
)

// BalancerChecks contains the name of each callable of Balancer, with the
// identities of its checks.
var BalancerChecks = struct {
	Route balancerRoute
}{
	Route: balancerRoute{"Balancer.Route"},
}

// balancerRoute is the name of Balancer.Route, with the identities of its
// checks.
type balancerRoute struct {
	conform.Callable[Balancer]
}

// Smoke returns the identity of the check that Route returns without a
// panic for derived inputs.
func (c balancerRoute) Smoke() conform.ID[Balancer] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Route returns
// without a panic for derived inputs.
func (c balancerRoute) SmokeGenerated() conform.ID[Balancer] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Route returns an error for
// a cancelled context.
func (c balancerRoute) Cancel() conform.ID[Balancer] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Route returns an
// error for a cancelled context.
func (c balancerRoute) CancelGenerated() conform.ID[Balancer] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Route returns
// an error when a call of its dependency fails.
func (c balancerRoute) FailureIsReported() conform.ID[Balancer] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Route
// returns an error when a call of its dependency fails.
func (c balancerRoute) FailureIsReportedGenerated() conform.ID[Balancer] {
	return c.Check("failure-is-reported#generated")
}

// KeepsItsInstance returns the identity of the check that Route sends each
// call with one key to one instance.
func (c balancerRoute) KeepsItsInstance() conform.ID[Balancer] {
	return c.Check("keeps-its-instance")
}

// KeepsItsInstanceGenerated returns the identity of the check that Route
// sends each call with one key to one instance.
func (c balancerRoute) KeepsItsInstanceGenerated() conform.ID[Balancer] {
	return c.Check("keeps-its-instance#generated")
}

// BalancerSuite returns the checks of Balancer, and the calls and the
// budgets of its callables.
func BalancerSuite() conform.Suite[Balancer] {
	return conform.Suite[Balancer]{
		Declaration: "Balancer",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Balancer]{
			{
				ID:          BalancerChecks.Route.Smoke(),
				Statement:   "Route returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BalancerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, string("key-a"))
					})
				},
				Defect: balancerRoutePanics,
			},
			{
				ID:          BalancerChecks.Route.SmokeGenerated(),
				Statement:   "Route returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BalancerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, drawnKey)
					})
				},
				Defect: balancerRoutePanics,
			},
			{
				ID:          BalancerChecks.Route.Cancel(),
				Statement:   "Route returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BalancerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, string("key-a"))
					})
				},
				Defect: balancerRouteIgnoresCancel,
			},
			{
				ID:          BalancerChecks.Route.CancelGenerated(),
				Statement:   "Route returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BalancerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, drawnKey)
					})
				},
				Defect: balancerRouteIgnoresCancel,
			},
			{
				ID:        BalancerChecks.Route.FailureIsReported(),
				Statement: "Route returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *BalancerHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, string("key-a"))
					})
				},
				Defect: balancerRouteSwallowsFailure,
			},
			{
				ID:        BalancerChecks.Route.FailureIsReportedGenerated(),
				Statement: "Route returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *BalancerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, drawnKey)
					})
				},
				Defect: balancerRouteSwallowsFailure,
			},
			{
				ID:        BalancerChecks.Route.KeepsItsInstance(),
				Statement: "Route sends each call with one key to one instance",
				Claim:     spec.KeepsItsInstance,
				Probe:     spec.UsesOneInstance,
				Reason:    spec.MovedInstance,
				Law:       spec.OneKeyOneInstance,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *BalancerHandle) {
					probe.UsesOneInstance(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, string("key-a"))
					})
				},
				Defect: balancerRouteBalancesEachCall,
			},
			{
				ID:        BalancerChecks.Route.KeepsItsInstanceGenerated(),
				Statement: "Route sends each call with one key to one instance",
				Claim:     spec.KeepsItsInstance,
				Probe:     spec.UsesOneInstance,
				Reason:    spec.MovedInstance,
				Law:       spec.OneKeyOneInstance,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *BalancerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UsesOneInstance(h, func(ctx context.Context) error {
						return h.Subject().Route(ctx, drawnKey)
					})
				},
				Defect: balancerRouteBalancesEachCall,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Balancer.Route",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Balancer]{
			{
				Calls:    balancerRouteCalls,
				Over:     newBalancerOver,
				Callable: BalancerChecks.Route.Callable,
			},
		},
	}
}

// RunBalancerSuite applies the checks of Balancer, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunBalancerSuite[T Balancer](t *testing.T, newSubject func() T, opts ...BalancerOption) {
	t.Helper()
	RunBalancerSuiteFromHandle(t, func(*BalancerHandle) T { return newSubject() }, opts...)
}

// RunBalancerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunBalancerSuiteFromHandle[T Balancer](
	t *testing.T,
	newSubject func(h *BalancerHandle) T,
	opts ...BalancerOption,
) {
	t.Helper()
	conform.Run(t, BalancerSuite(), balancerSubjects(newSubject), opts...)
}

// FuzzBalancerSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzBalancerSuite[T Balancer](
	f *testing.F,
	callable conform.Callable[Balancer],
	newSubject func(h *BalancerHandle) T,
	opts ...BalancerOption,
) {
	f.Helper()
	conform.Fuzz(f, BalancerSuite(), callable, balancerSubjects(newSubject), opts...)
}

// RunBalancerBudgets applies the allocation ceilings of Balancer and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunBalancerBudgets[T Balancer](
	t *testing.T,
	newSubject func(h *BalancerHandle) T,
	opts ...BalancerOption,
) {
	t.Helper()
	conform.RunBudgets(t, BalancerSuite(), balancerSubjects(newSubject), opts...)
}

// BenchBalancerSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchBalancerSuite[T Balancer](
	b conform.B,
	callable conform.Callable[Balancer],
	newSubject func(h *BalancerHandle) T,
	opts ...BalancerOption,
) {
	b.Helper()
	conform.Bench(b, BalancerSuite(), callable, balancerSubjects(newSubject), opts...)
}

// balancerSubjects returns the factory of the subjects of Balancer that
// calls newSubject with the handle of each check.
func balancerSubjects[T Balancer](
	newSubject func(h *BalancerHandle) T,
) func(h *BalancerHandle) Balancer {
	return func(h *BalancerHandle) Balancer { return newSubject(h) }
}

// balancerRouteCalls generates the calls of Balancer.Route with the inputs
// that each case draws.
var balancerRouteCalls = prop.Composite(func(c *prop.Case) func(context.Context, Balancer) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Balancer) error {
		return subject.Route(ctx, key)
	}
})

// balancerOverEmbedded is Balancer under the name of the field that
// balancerOver embeds.
type balancerOverEmbedded = Balancer

// balancerOver is a subject of Balancer whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type balancerOver struct {
	balancerOverEmbedded
	spend func()
}

// Route calls spend and returns its zero results.
func (o balancerOver) Route(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// newBalancerOver returns the planted defect of a ceiling of Balancer,
// whose callables call spend.
func newBalancerOver(_ *BalancerHandle, spend func()) Balancer {
	return balancerOver{spend: spend}
}

// balancerRoutePanicsDefect is the planted defect of a Balancer whose
// Route panics.
type balancerRoutePanicsDefect struct {
	balancerOverEmbedded
}

// Route is the override of the planted defect.
func (*balancerRoutePanicsDefect) Route(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Balancer whose Route panics")
}

// balancerRoutePanics returns the planted defect of a Balancer whose Route
// panics over a fresh subject.
func balancerRoutePanics(h *BalancerHandle) Balancer {
	return &balancerRoutePanicsDefect{balancerOverEmbedded: h.Fresh()}
}

// balancerRouteIgnoresCancelDefect is the planted defect of a Balancer
// whose Route ignores the context.
type balancerRouteIgnoresCancelDefect struct {
	balancerOverEmbedded
}

// Route is the override of the planted defect.
func (d *balancerRouteIgnoresCancelDefect) Route(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.balancerOverEmbedded.Route(p0, p1)
}

// balancerRouteIgnoresCancel returns the planted defect of a Balancer
// whose Route ignores the context over a fresh subject.
func balancerRouteIgnoresCancel(h *BalancerHandle) Balancer {
	return &balancerRouteIgnoresCancelDefect{balancerOverEmbedded: h.Fresh()}
}

// balancerRouteSwallowsFailureDefect is the planted defect of a Balancer
// whose Route returns no error when its dependency fails.
type balancerRouteSwallowsFailureDefect struct {
	balancerOverEmbedded
	h *BalancerHandle
}

// Route is the override of the planted defect.
func (d *balancerRouteSwallowsFailureDefect) Route(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.balancerOverEmbedded.Route(p0, p1)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// balancerRouteSwallowsFailure returns the planted defect of a Balancer
// whose Route returns no error when its dependency fails over a fresh
// subject.
func balancerRouteSwallowsFailure(h *BalancerHandle) Balancer {
	return &balancerRouteSwallowsFailureDefect{balancerOverEmbedded: h.Fresh(), h: h}
}

// balancerRouteBalancesEachCallDefect is the planted defect of a Balancer
// whose Route sends each call to another instance.
type balancerRouteBalancesEachCallDefect struct {
	balancerOverEmbedded
	h *BalancerHandle
}

// Route is the override of the planted defect.
func (d *balancerRouteBalancesEachCallDefect) Route(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return d.h.Fresh().Route(p0, p1)
}

// balancerRouteBalancesEachCall returns the planted defect of a Balancer
// whose Route sends each call to another instance over a fresh subject.
func balancerRouteBalancesEachCall(h *BalancerHandle) Balancer {
	return &balancerRouteBalancesEachCallDefect{balancerOverEmbedded: h.Fresh(), h: h}
}

// The derived checks of Copier:
//
// Copier.Copy/smoke: Copy returns without a panic for derived inputs.
//
// Copier.Copy/smoke#generated: Copy returns without a panic for derived
// inputs.
//
// Copier.Copy/cancel: Copy returns an error for a cancelled context.
//
// Copier.Copy/cancel#generated: Copy returns an error for a cancelled
// context.
//
// Copier.Copy/failure-is-reported: Copy returns an error when a call of
// its dependency fails.
//
// Copier.Copy/failure-is-reported#generated: Copy returns an error when a
// call of its dependency fails.
//
// Copier.Copy/releases-on-every-path: Copy releases each resource that it
// acquired when a call of its dependency fails.
//
// Copier.Copy/releases-on-every-path#generated: Copy releases each
// resource that it acquired when a call of its dependency fails.
//
// Copier.Open/smoke: Open returns without a panic for derived inputs.
//
// Copier.Open/cancel: Open returns an error for a cancelled context.
//
// Copier.Open/failure-is-reported: Open returns an error when a call of
// its dependency fails.
//
// Copier.Close/smoke: Close returns without a panic for derived inputs.
//
// Copier.Close/cancel: Close returns an error for a cancelled context.
//
// Copier.Close/failure-is-reported: Close returns an error when a call of
// its dependency fails.

// The derived checks of Copier that witness their laws at the shallow end:
//
// Copier.Copy/smoke
//
// Copier.Open/smoke
//
// Copier.Close/smoke

// The unclaimed classifications of Copier:
//
// The shape:writer of Copier.Copy has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Copier.Open has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The shape:lifecycle of Copier.Close has the stance silent. A lifecycle
// call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Copier writes.
type (
	CopierHandle = conform.Handle[Copier]
	CopierCheck  = conform.Check[Copier]
	CopierOption = conform.Option[Copier]
)

// CopierChecks contains the name of each callable of Copier, with the
// identities of its checks.
var CopierChecks = struct {
	Copy  copierCopy
	Open  copierOpen
	Close copierClose
}{
	Copy:  copierCopy{"Copier.Copy"},
	Open:  copierOpen{"Copier.Open"},
	Close: copierClose{"Copier.Close"},
}

// copierCopy is the name of Copier.Copy, with the identities of its
// checks.
type copierCopy struct {
	conform.Callable[Copier]
}

// Smoke returns the identity of the check that Copy returns without a
// panic for derived inputs.
func (c copierCopy) Smoke() conform.ID[Copier] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Copy returns
// without a panic for derived inputs.
func (c copierCopy) SmokeGenerated() conform.ID[Copier] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Copy returns an error for
// a cancelled context.
func (c copierCopy) Cancel() conform.ID[Copier] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Copy returns an
// error for a cancelled context.
func (c copierCopy) CancelGenerated() conform.ID[Copier] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Copy returns an
// error when a call of its dependency fails.
func (c copierCopy) FailureIsReported() conform.ID[Copier] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Copy
// returns an error when a call of its dependency fails.
func (c copierCopy) FailureIsReportedGenerated() conform.ID[Copier] {
	return c.Check("failure-is-reported#generated")
}

// ReleasesOnEveryPath returns the identity of the check that Copy releases
// each resource that it acquired when a call of its dependency fails.
func (c copierCopy) ReleasesOnEveryPath() conform.ID[Copier] {
	return c.Check("releases-on-every-path")
}

// ReleasesOnEveryPathGenerated returns the identity of the check that Copy
// releases each resource that it acquired when a call of its dependency
// fails.
func (c copierCopy) ReleasesOnEveryPathGenerated() conform.ID[Copier] {
	return c.Check("releases-on-every-path#generated")
}

// copierOpen is the name of Copier.Open, with the identities of its
// checks.
type copierOpen struct {
	conform.Callable[Copier]
}

// Smoke returns the identity of the check that Open returns without a
// panic for derived inputs.
func (c copierOpen) Smoke() conform.ID[Copier] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Open returns an error for
// a cancelled context.
func (c copierOpen) Cancel() conform.ID[Copier] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Open returns an
// error when a call of its dependency fails.
func (c copierOpen) FailureIsReported() conform.ID[Copier] {
	return c.Check("failure-is-reported")
}

// copierClose is the name of Copier.Close, with the identities of its
// checks.
type copierClose struct {
	conform.Callable[Copier]
}

// Smoke returns the identity of the check that Close returns without a
// panic for derived inputs.
func (c copierClose) Smoke() conform.ID[Copier] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Close returns an error for
// a cancelled context.
func (c copierClose) Cancel() conform.ID[Copier] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Close returns
// an error when a call of its dependency fails.
func (c copierClose) FailureIsReported() conform.ID[Copier] {
	return c.Check("failure-is-reported")
}

// CopierSuite returns the checks of Copier, and the calls and the budgets
// of its callables.
func CopierSuite() conform.Suite[Copier] {
	return conform.Suite[Copier]{
		Declaration: "Copier",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Copier]{
			{
				ID:          CopierChecks.Copy.Smoke(),
				Statement:   "Copy returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CopierHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, string("key-a"))
					})
				},
				Defect: copierCopyPanics,
			},
			{
				ID:          CopierChecks.Copy.SmokeGenerated(),
				Statement:   "Copy returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CopierHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, drawnKey)
					})
				},
				Defect: copierCopyPanics,
			},
			{
				ID:          CopierChecks.Copy.Cancel(),
				Statement:   "Copy returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CopierHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, string("key-a"))
					})
				},
				Defect: copierCopyIgnoresCancel,
			},
			{
				ID:          CopierChecks.Copy.CancelGenerated(),
				Statement:   "Copy returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CopierHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, drawnKey)
					})
				},
				Defect: copierCopyIgnoresCancel,
			},
			{
				ID:        CopierChecks.Copy.FailureIsReported(),
				Statement: "Copy returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *CopierHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, string("key-a"))
					})
				},
				Defect: copierCopySwallowsFailure,
			},
			{
				ID:        CopierChecks.Copy.FailureIsReportedGenerated(),
				Statement: "Copy returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *CopierHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, drawnKey)
					})
				},
				Defect: copierCopySwallowsFailure,
			},
			{
				ID:        CopierChecks.Copy.ReleasesOnEveryPath(),
				Statement: "Copy releases each resource that it acquired when a call of its dependency fails",
				Claim:     spec.ReleasesOnEveryPath,
				Probe:     spec.ReleasesWhatItAcquires,
				Reason:    spec.LeakedResource,
				Law:       spec.FailurePathReleases,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *CopierHandle) {
					probe.ReleasesWhatItAcquires(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, string("key-a"))
					}, "Open", "Close")
				},
				Defect: copierCopyOpensWithoutClosing,
			},
			{
				ID:        CopierChecks.Copy.ReleasesOnEveryPathGenerated(),
				Statement: "Copy releases each resource that it acquired when a call of its dependency fails",
				Claim:     spec.ReleasesOnEveryPath,
				Probe:     spec.ReleasesWhatItAcquires,
				Reason:    spec.LeakedResource,
				Law:       spec.FailurePathReleases,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *CopierHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReleasesWhatItAcquires(h, func(ctx context.Context) error {
						return h.Subject().Copy(ctx, drawnKey)
					}, "Open", "Close")
				},
				Defect: copierCopyOpensWithoutClosing,
			},
			{
				ID:          CopierChecks.Open.Smoke(),
				Statement:   "Open returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CopierHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Open(ctx)
					})
				},
				Defect: copierOpenPanics,
			},
			{
				ID:          CopierChecks.Open.Cancel(),
				Statement:   "Open returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CopierHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Open(ctx)
					})
				},
				Defect: copierOpenIgnoresCancel,
			},
			{
				ID:        CopierChecks.Open.FailureIsReported(),
				Statement: "Open returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *CopierHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Open(ctx)
					})
				},
				Defect: copierOpenSwallowsFailure,
			},
			{
				ID:          CopierChecks.Close.Smoke(),
				Statement:   "Close returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CopierHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Close(ctx)
					})
				},
				Defect: copierClosePanics,
			},
			{
				ID:          CopierChecks.Close.Cancel(),
				Statement:   "Close returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CopierHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Close(ctx)
					})
				},
				Defect: copierCloseIgnoresCancel,
			},
			{
				ID:        CopierChecks.Close.FailureIsReported(),
				Statement: "Close returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *CopierHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Close(ctx)
					})
				},
				Defect: copierCloseSwallowsFailure,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Copier.Copy",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Copier.Open",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Copier.Close",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Copier]{
			{
				Calls:    copierCopyCalls,
				Over:     newCopierOver,
				Callable: CopierChecks.Copy.Callable,
			},
			{
				Calls:    copierOpenCalls,
				Over:     newCopierOver,
				Callable: CopierChecks.Open.Callable,
			},
			{
				Calls:    copierCloseCalls,
				Over:     newCopierOver,
				Callable: CopierChecks.Close.Callable,
			},
		},
	}
}

// RunCopierSuite applies the checks of Copier, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCopierSuite[T Copier](t *testing.T, newSubject func() T, opts ...CopierOption) {
	t.Helper()
	RunCopierSuiteFromHandle(t, func(*CopierHandle) T { return newSubject() }, opts...)
}

// RunCopierSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCopierSuiteFromHandle[T Copier](
	t *testing.T,
	newSubject func(h *CopierHandle) T,
	opts ...CopierOption,
) {
	t.Helper()
	conform.Run(t, CopierSuite(), copierSubjects(newSubject), opts...)
}

// FuzzCopierSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCopierSuite[T Copier](
	f *testing.F,
	callable conform.Callable[Copier],
	newSubject func(h *CopierHandle) T,
	opts ...CopierOption,
) {
	f.Helper()
	conform.Fuzz(f, CopierSuite(), callable, copierSubjects(newSubject), opts...)
}

// RunCopierBudgets applies the allocation ceilings of Copier and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCopierBudgets[T Copier](
	t *testing.T,
	newSubject func(h *CopierHandle) T,
	opts ...CopierOption,
) {
	t.Helper()
	conform.RunBudgets(t, CopierSuite(), copierSubjects(newSubject), opts...)
}

// BenchCopierSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchCopierSuite[T Copier](
	b conform.B,
	callable conform.Callable[Copier],
	newSubject func(h *CopierHandle) T,
	opts ...CopierOption,
) {
	b.Helper()
	conform.Bench(b, CopierSuite(), callable, copierSubjects(newSubject), opts...)
}

// copierSubjects returns the factory of the subjects of Copier that calls
// newSubject with the handle of each check.
func copierSubjects[T Copier](newSubject func(h *CopierHandle) T) func(h *CopierHandle) Copier {
	return func(h *CopierHandle) Copier { return newSubject(h) }
}

// copierCopyCalls generates the calls of Copier.Copy with the inputs that
// each case draws.
var copierCopyCalls = prop.Composite(func(c *prop.Case) func(context.Context, Copier) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Copier) error {
		return subject.Copy(ctx, key)
	}
})

// copierOpenCalls generates the calls of Copier.Open with the inputs that
// each case draws.
var copierOpenCalls = prop.Composite(func(*prop.Case) func(context.Context, Copier) error {
	return func(ctx context.Context, subject Copier) error {
		return subject.Open(ctx)
	}
})

// copierCloseCalls generates the calls of Copier.Close with the inputs
// that each case draws.
var copierCloseCalls = prop.Composite(func(*prop.Case) func(context.Context, Copier) error {
	return func(ctx context.Context, subject Copier) error {
		return subject.Close(ctx)
	}
})

// copierOverEmbedded is Copier under the name of the field that copierOver
// embeds.
type copierOverEmbedded = Copier

// copierOver is a subject of Copier whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type copierOver struct {
	copierOverEmbedded
	spend func()
}

// Copy calls spend and returns its zero results.
func (o copierOver) Copy(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Open calls spend and returns its zero results.
func (o copierOver) Open(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Close calls spend and returns its zero results.
func (o copierOver) Close(context.Context) (r0 error) {
	o.spend()
	return r0
}

// newCopierOver returns the planted defect of a ceiling of Copier, whose
// callables call spend.
func newCopierOver(_ *CopierHandle, spend func()) Copier {
	return copierOver{spend: spend}
}

// copierCopyPanicsDefect is the planted defect of a Copier whose Copy
// panics.
type copierCopyPanicsDefect struct {
	copierOverEmbedded
}

// Copy is the override of the planted defect.
func (*copierCopyPanicsDefect) Copy(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Copier whose Copy panics")
}

// copierCopyPanics returns the planted defect of a Copier whose Copy
// panics over a fresh subject.
func copierCopyPanics(h *CopierHandle) Copier {
	return &copierCopyPanicsDefect{copierOverEmbedded: h.Fresh()}
}

// copierCopyIgnoresCancelDefect is the planted defect of a Copier whose
// Copy ignores the context.
type copierCopyIgnoresCancelDefect struct {
	copierOverEmbedded
}

// Copy is the override of the planted defect.
func (d *copierCopyIgnoresCancelDefect) Copy(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.copierOverEmbedded.Copy(p0, p1)
}

// copierCopyIgnoresCancel returns the planted defect of a Copier whose
// Copy ignores the context over a fresh subject.
func copierCopyIgnoresCancel(h *CopierHandle) Copier {
	return &copierCopyIgnoresCancelDefect{copierOverEmbedded: h.Fresh()}
}

// copierCopySwallowsFailureDefect is the planted defect of a Copier whose
// Copy returns no error when its dependency fails.
type copierCopySwallowsFailureDefect struct {
	copierOverEmbedded
	h *CopierHandle
}

// Copy is the override of the planted defect.
func (d *copierCopySwallowsFailureDefect) Copy(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.copierOverEmbedded.Copy(p0, p1)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// copierCopySwallowsFailure returns the planted defect of a Copier whose
// Copy returns no error when its dependency fails over a fresh subject.
func copierCopySwallowsFailure(h *CopierHandle) Copier {
	return &copierCopySwallowsFailureDefect{copierOverEmbedded: h.Fresh(), h: h}
}

// copierCopyOpensWithoutClosingDefect is the planted defect of a Copier
// whose Copy keeps a resource that it acquired before a failure.
type copierCopyOpensWithoutClosingDefect struct {
	copierOverEmbedded
	leak *conform.Dependency
}

// Copy is the override of the planted defect.
func (d *copierCopyOpensWithoutClosingDefect) Copy(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	r0 = d.copierOverEmbedded.Copy(p0, p1)
	if r0 != nil {
		_ = d.leak.Call("Open")
	}
	return r0
}

// copierCopyOpensWithoutClosing returns the planted defect of a Copier
// whose Copy keeps a resource that it acquired before a failure over a
// fresh subject.
func copierCopyOpensWithoutClosing(h *CopierHandle) Copier {
	return &copierCopyOpensWithoutClosingDefect{
		copierOverEmbedded: h.Fresh(),
		leak:               h.Injector().Dependency(),
	}
}

// copierOpenPanicsDefect is the planted defect of a Copier whose Open
// panics.
type copierOpenPanicsDefect struct {
	copierOverEmbedded
}

// Open is the override of the planted defect.
func (*copierOpenPanicsDefect) Open(p0 context.Context) (r0 error) {
	panic("the planted defect of a Copier whose Open panics")
}

// copierOpenPanics returns the planted defect of a Copier whose Open
// panics over a fresh subject.
func copierOpenPanics(h *CopierHandle) Copier {
	return &copierOpenPanicsDefect{copierOverEmbedded: h.Fresh()}
}

// copierOpenIgnoresCancelDefect is the planted defect of a Copier whose
// Open ignores the context.
type copierOpenIgnoresCancelDefect struct {
	copierOverEmbedded
}

// Open is the override of the planted defect.
func (d *copierOpenIgnoresCancelDefect) Open(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.copierOverEmbedded.Open(p0)
}

// copierOpenIgnoresCancel returns the planted defect of a Copier whose
// Open ignores the context over a fresh subject.
func copierOpenIgnoresCancel(h *CopierHandle) Copier {
	return &copierOpenIgnoresCancelDefect{copierOverEmbedded: h.Fresh()}
}

// copierOpenSwallowsFailureDefect is the planted defect of a Copier whose
// Open returns no error when its dependency fails.
type copierOpenSwallowsFailureDefect struct {
	copierOverEmbedded
	h *CopierHandle
}

// Open is the override of the planted defect.
func (d *copierOpenSwallowsFailureDefect) Open(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.copierOverEmbedded.Open(p0)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// copierOpenSwallowsFailure returns the planted defect of a Copier whose
// Open returns no error when its dependency fails over a fresh subject.
func copierOpenSwallowsFailure(h *CopierHandle) Copier {
	return &copierOpenSwallowsFailureDefect{copierOverEmbedded: h.Fresh(), h: h}
}

// copierClosePanicsDefect is the planted defect of a Copier whose Close
// panics.
type copierClosePanicsDefect struct {
	copierOverEmbedded
}

// Close is the override of the planted defect.
func (*copierClosePanicsDefect) Close(p0 context.Context) (r0 error) {
	panic("the planted defect of a Copier whose Close panics")
}

// copierClosePanics returns the planted defect of a Copier whose Close
// panics over a fresh subject.
func copierClosePanics(h *CopierHandle) Copier {
	return &copierClosePanicsDefect{copierOverEmbedded: h.Fresh()}
}

// copierCloseIgnoresCancelDefect is the planted defect of a Copier whose
// Close ignores the context.
type copierCloseIgnoresCancelDefect struct {
	copierOverEmbedded
}

// Close is the override of the planted defect.
func (d *copierCloseIgnoresCancelDefect) Close(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.copierOverEmbedded.Close(p0)
}

// copierCloseIgnoresCancel returns the planted defect of a Copier whose
// Close ignores the context over a fresh subject.
func copierCloseIgnoresCancel(h *CopierHandle) Copier {
	return &copierCloseIgnoresCancelDefect{copierOverEmbedded: h.Fresh()}
}

// copierCloseSwallowsFailureDefect is the planted defect of a Copier whose
// Close returns no error when its dependency fails.
type copierCloseSwallowsFailureDefect struct {
	copierOverEmbedded
	h *CopierHandle
}

// Close is the override of the planted defect.
func (d *copierCloseSwallowsFailureDefect) Close(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.copierOverEmbedded.Close(p0)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// copierCloseSwallowsFailure returns the planted defect of a Copier whose
// Close returns no error when its dependency fails over a fresh subject.
func copierCloseSwallowsFailure(h *CopierHandle) Copier {
	return &copierCloseSwallowsFailureDefect{copierOverEmbedded: h.Fresh(), h: h}
}

// The derived checks of Undeclared:
//
// Undeclared.Put/smoke: Put returns without a panic for derived inputs.
//
// Undeclared.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Undeclared.Put/cancel: Put returns an error for a cancelled context.
//
// Undeclared.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Undeclared.Get/smoke: Get returns without a panic for derived inputs.
//
// Undeclared.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Undeclared.Get/cancel: Get returns an error for a cancelled context.
//
// Undeclared.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Undeclared.Get/zero-on-error: Get returns the zero value beside an
// error.
//
// Undeclared.Get/zero-on-error#generated: Get returns the zero value
// beside an error.

// The derived checks of Undeclared that witness their laws at the shallow
// end:
//
// Undeclared.Put/smoke
//
// Undeclared.Get/smoke

// The unclaimed classifications of Undeclared:
//
// The shape:writer of Undeclared.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The mixin:atomic of Undeclared.Put has the stance needs. The claim
// no-partial-write needs the param machineries=[faults] of the directive
// suite.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Undeclared writes.
type (
	UndeclaredHandle = conform.Handle[Undeclared]
	UndeclaredCheck  = conform.Check[Undeclared]
	UndeclaredOption = conform.Option[Undeclared]
)

// UndeclaredChecks contains the name of each callable of Undeclared, with
// the identities of its checks.
var UndeclaredChecks = struct {
	Put undeclaredPut
	Get undeclaredGet
}{
	Put: undeclaredPut{"Undeclared.Put"},
	Get: undeclaredGet{"Undeclared.Get"},
}

// undeclaredPut is the name of Undeclared.Put, with the identities of its
// checks.
type undeclaredPut struct {
	conform.Callable[Undeclared]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c undeclaredPut) Smoke() conform.ID[Undeclared] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c undeclaredPut) SmokeGenerated() conform.ID[Undeclared] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c undeclaredPut) Cancel() conform.ID[Undeclared] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c undeclaredPut) CancelGenerated() conform.ID[Undeclared] {
	return c.Check("cancel#generated")
}

// undeclaredGet is the name of Undeclared.Get, with the identities of its
// checks.
type undeclaredGet struct {
	conform.Callable[Undeclared]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c undeclaredGet) Smoke() conform.ID[Undeclared] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c undeclaredGet) SmokeGenerated() conform.ID[Undeclared] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c undeclaredGet) Cancel() conform.ID[Undeclared] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c undeclaredGet) CancelGenerated() conform.ID[Undeclared] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c undeclaredGet) ZeroOnError() conform.ID[Undeclared] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c undeclaredGet) ZeroOnErrorGenerated() conform.ID[Undeclared] {
	return c.Check("zero-on-error#generated")
}

// UndeclaredSuite returns the checks of Undeclared, and the calls and the
// budgets of its callables.
func UndeclaredSuite() conform.Suite[Undeclared] {
	return conform.Suite[Undeclared]{
		Declaration: "Undeclared",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Undeclared]{
			{
				ID:          UndeclaredChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *UndeclaredHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: undeclaredPutPanics,
			},
			{
				ID:          UndeclaredChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *UndeclaredHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnV)
					})
				},
				Defect: undeclaredPutPanics,
			},
			{
				ID:          UndeclaredChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *UndeclaredHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: undeclaredPutIgnoresCancel,
			},
			{
				ID:          UndeclaredChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *UndeclaredHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[Value](), "v")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnV)
					})
				},
				Defect: undeclaredPutIgnoresCancel,
			},
			{
				ID:          UndeclaredChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *UndeclaredHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: undeclaredGetPanics,
			},
			{
				ID:          UndeclaredChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *UndeclaredHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: undeclaredGetPanics,
			},
			{
				ID:          UndeclaredChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *UndeclaredHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: undeclaredGetIgnoresCancel,
			},
			{
				ID:          UndeclaredChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *UndeclaredHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: undeclaredGetIgnoresCancel,
			},
			{
				ID:          UndeclaredChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *UndeclaredHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: undeclaredGetValueBeside,
			},
			{
				ID:          UndeclaredChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *UndeclaredHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Value{})
				},
				Defect: undeclaredGetValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Undeclared.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "mixin:atomic",
				Callable:       "Undeclared.Put",
				Stance:         spec.Needs,
				Note:           "The claim no-partial-write needs the param machineries=[faults] of the directive suite.",
			},
		},
		Calls: []conform.Call[Undeclared]{
			{
				Calls:    undeclaredPutCalls,
				Over:     newUndeclaredOver,
				Callable: UndeclaredChecks.Put.Callable,
			},
			{
				Calls:    undeclaredGetCalls,
				Over:     newUndeclaredOver,
				Callable: UndeclaredChecks.Get.Callable,
			},
		},
	}
}

// RunUndeclaredSuite applies the checks of Undeclared, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunUndeclaredSuite[T Undeclared](t *testing.T, newSubject func() T, opts ...UndeclaredOption) {
	t.Helper()
	RunUndeclaredSuiteFromHandle(t, func(*UndeclaredHandle) T { return newSubject() }, opts...)
}

// RunUndeclaredSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunUndeclaredSuiteFromHandle[T Undeclared](
	t *testing.T,
	newSubject func(h *UndeclaredHandle) T,
	opts ...UndeclaredOption,
) {
	t.Helper()
	conform.Run(t, UndeclaredSuite(), undeclaredSubjects(newSubject), opts...)
}

// FuzzUndeclaredSuite registers on f the fuzz target of the property
// checks of callable, on a subject that newSubject builds from the handle
// of each case.
func FuzzUndeclaredSuite[T Undeclared](
	f *testing.F,
	callable conform.Callable[Undeclared],
	newSubject func(h *UndeclaredHandle) T,
	opts ...UndeclaredOption,
) {
	f.Helper()
	conform.Fuzz(f, UndeclaredSuite(), callable, undeclaredSubjects(newSubject), opts...)
}

// RunUndeclaredBudgets applies the allocation ceilings of Undeclared and
// its checks that count allocations, and the ceilings and the checks that
// opts add, to a fresh subject of newSubject for each case.
func RunUndeclaredBudgets[T Undeclared](
	t *testing.T,
	newSubject func(h *UndeclaredHandle) T,
	opts ...UndeclaredOption,
) {
	t.Helper()
	conform.RunBudgets(t, UndeclaredSuite(), undeclaredSubjects(newSubject), opts...)
}

// BenchUndeclaredSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchUndeclaredSuite[T Undeclared](
	b conform.B,
	callable conform.Callable[Undeclared],
	newSubject func(h *UndeclaredHandle) T,
	opts ...UndeclaredOption,
) {
	b.Helper()
	conform.Bench(b, UndeclaredSuite(), callable, undeclaredSubjects(newSubject), opts...)
}

// undeclaredSubjects returns the factory of the subjects of Undeclared
// that calls newSubject with the handle of each check.
func undeclaredSubjects[T Undeclared](
	newSubject func(h *UndeclaredHandle) T,
) func(h *UndeclaredHandle) Undeclared {
	return func(h *UndeclaredHandle) Undeclared { return newSubject(h) }
}

// undeclaredPutCalls generates the calls of Undeclared.Put with the inputs
// that each case draws.
var undeclaredPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Undeclared) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Undeclared) error {
		return subject.Put(ctx, v)
	}
})

// undeclaredGetCalls generates the calls of Undeclared.Get with the inputs
// that each case draws.
var undeclaredGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Undeclared) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Undeclared) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// undeclaredOverEmbedded is Undeclared under the name of the field that
// undeclaredOver embeds.
type undeclaredOverEmbedded = Undeclared

// undeclaredOver is a subject of Undeclared whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type undeclaredOver struct {
	undeclaredOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o undeclaredOver) Put(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o undeclaredOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newUndeclaredOver returns the planted defect of a ceiling of Undeclared,
// whose callables call spend.
func newUndeclaredOver(_ *UndeclaredHandle, spend func()) Undeclared {
	return undeclaredOver{spend: spend}
}

// undeclaredPutPanicsDefect is the planted defect of a Undeclared whose
// Put panics.
type undeclaredPutPanicsDefect struct {
	undeclaredOverEmbedded
}

// Put is the override of the planted defect.
func (*undeclaredPutPanicsDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Undeclared whose Put panics")
}

// undeclaredPutPanics returns the planted defect of a Undeclared whose Put
// panics over a fresh subject.
func undeclaredPutPanics(h *UndeclaredHandle) Undeclared {
	return &undeclaredPutPanicsDefect{undeclaredOverEmbedded: h.Fresh()}
}

// undeclaredPutIgnoresCancelDefect is the planted defect of a Undeclared
// whose Put ignores the context.
type undeclaredPutIgnoresCancelDefect struct {
	undeclaredOverEmbedded
}

// Put is the override of the planted defect.
func (d *undeclaredPutIgnoresCancelDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.undeclaredOverEmbedded.Put(p0, p1)
}

// undeclaredPutIgnoresCancel returns the planted defect of a Undeclared
// whose Put ignores the context over a fresh subject.
func undeclaredPutIgnoresCancel(h *UndeclaredHandle) Undeclared {
	return &undeclaredPutIgnoresCancelDefect{undeclaredOverEmbedded: h.Fresh()}
}

// undeclaredGetPanicsDefect is the planted defect of a Undeclared whose
// Get panics.
type undeclaredGetPanicsDefect struct {
	undeclaredOverEmbedded
}

// Get is the override of the planted defect.
func (*undeclaredGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Undeclared whose Get panics")
}

// undeclaredGetPanics returns the planted defect of a Undeclared whose Get
// panics over a fresh subject.
func undeclaredGetPanics(h *UndeclaredHandle) Undeclared {
	return &undeclaredGetPanicsDefect{undeclaredOverEmbedded: h.Fresh()}
}

// undeclaredGetIgnoresCancelDefect is the planted defect of a Undeclared
// whose Get ignores the context.
type undeclaredGetIgnoresCancelDefect struct {
	undeclaredOverEmbedded
}

// Get is the override of the planted defect.
func (d *undeclaredGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.undeclaredOverEmbedded.Get(p0, p1)
}

// undeclaredGetIgnoresCancel returns the planted defect of a Undeclared
// whose Get ignores the context over a fresh subject.
func undeclaredGetIgnoresCancel(h *UndeclaredHandle) Undeclared {
	return &undeclaredGetIgnoresCancelDefect{undeclaredOverEmbedded: h.Fresh()}
}

// undeclaredGetValueBesideDefect is the planted defect of a Undeclared
// whose Get returns a value beside an error.
type undeclaredGetValueBesideDefect struct {
	undeclaredOverEmbedded
}

// Get is the override of the planted defect.
func (d *undeclaredGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.undeclaredOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// undeclaredGetValueBeside returns the planted defect of a Undeclared
// whose Get returns a value beside an error over a fresh subject.
func undeclaredGetValueBeside(h *UndeclaredHandle) Undeclared {
	return &undeclaredGetValueBesideDefect{undeclaredOverEmbedded: h.Fresh()}
}

// The derived checks of Gravity:
//
// Gravity.Get/smoke: Get returns without a panic for derived inputs.
//
// Gravity.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Gravity.Get/cancel: Get returns an error for a cancelled context.
//
// Gravity.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Gravity.Get/failure-is-reported: Get returns an error when a call of its
// dependency fails.
//
// Gravity.Get/failure-is-reported#generated: Get returns an error when a
// call of its dependency fails.
//
// Gravity.Get/zero-on-error: Get returns the zero value beside an error.
//
// Gravity.Get/zero-on-error#generated: Get returns the zero value beside
// an error.

// The derived checks of Gravity that witness their laws at the shallow
// end:
//
// Gravity.Get/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Gravity writes.
type (
	GravityHandle = conform.Handle[Gravity]
	GravityCheck  = conform.Check[Gravity]
	GravityOption = conform.Option[Gravity]
)

// GravityChecks contains the name of each callable of Gravity, with the
// identities of its checks.
var GravityChecks = struct {
	Get gravityGet
}{
	Get: gravityGet{"Gravity.Get"},
}

// gravityGet is the name of Gravity.Get, with the identities of its
// checks.
type gravityGet struct {
	conform.Callable[Gravity]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c gravityGet) Smoke() conform.ID[Gravity] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c gravityGet) SmokeGenerated() conform.ID[Gravity] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c gravityGet) Cancel() conform.ID[Gravity] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c gravityGet) CancelGenerated() conform.ID[Gravity] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Get returns an
// error when a call of its dependency fails.
func (c gravityGet) FailureIsReported() conform.ID[Gravity] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Get
// returns an error when a call of its dependency fails.
func (c gravityGet) FailureIsReportedGenerated() conform.ID[Gravity] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c gravityGet) ZeroOnError() conform.ID[Gravity] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c gravityGet) ZeroOnErrorGenerated() conform.ID[Gravity] {
	return c.Check("zero-on-error#generated")
}

// GravitySuite returns the checks of Gravity, and the calls and the
// budgets of its callables.
func GravitySuite() conform.Suite[Gravity] {
	return conform.Suite[Gravity]{
		Declaration: "Gravity",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Gravity]{
			{
				ID:          GravityChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *GravityHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: gravityGetPanics,
			},
			{
				ID:          GravityChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *GravityHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: gravityGetPanics,
			},
			{
				ID:          GravityChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *GravityHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: gravityGetIgnoresCancel,
			},
			{
				ID:          GravityChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *GravityHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: gravityGetIgnoresCancel,
			},
			{
				ID:        GravityChecks.Get.FailureIsReported(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *GravityHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: gravityGetSwallowsFailure,
			},
			{
				ID:        GravityChecks.Get.FailureIsReportedGenerated(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *GravityHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: gravityGetSwallowsFailure,
			},
			{
				ID:          GravityChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *GravityHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: gravityGetValueBeside,
			},
			{
				ID:          GravityChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *GravityHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Value{})
				},
				Defect: gravityGetValueBeside,
			},
		},
		Calls: []conform.Call[Gravity]{
			{
				Calls:    gravityGetCalls,
				Over:     newGravityOver,
				Callable: GravityChecks.Get.Callable,
			},
		},
	}
}

// RunGravitySuite applies the checks of Gravity, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunGravitySuite[T Gravity](t *testing.T, newSubject func() T, opts ...GravityOption) {
	t.Helper()
	RunGravitySuiteFromHandle(t, func(*GravityHandle) T { return newSubject() }, opts...)
}

// RunGravitySuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunGravitySuiteFromHandle[T Gravity](
	t *testing.T,
	newSubject func(h *GravityHandle) T,
	opts ...GravityOption,
) {
	t.Helper()
	conform.Run(t, GravitySuite(), gravitySubjects(newSubject), opts...)
}

// FuzzGravitySuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzGravitySuite[T Gravity](
	f *testing.F,
	callable conform.Callable[Gravity],
	newSubject func(h *GravityHandle) T,
	opts ...GravityOption,
) {
	f.Helper()
	conform.Fuzz(f, GravitySuite(), callable, gravitySubjects(newSubject), opts...)
}

// RunGravityBudgets applies the allocation ceilings of Gravity and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunGravityBudgets[T Gravity](
	t *testing.T,
	newSubject func(h *GravityHandle) T,
	opts ...GravityOption,
) {
	t.Helper()
	conform.RunBudgets(t, GravitySuite(), gravitySubjects(newSubject), opts...)
}

// BenchGravitySuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchGravitySuite[T Gravity](
	b conform.B,
	callable conform.Callable[Gravity],
	newSubject func(h *GravityHandle) T,
	opts ...GravityOption,
) {
	b.Helper()
	conform.Bench(b, GravitySuite(), callable, gravitySubjects(newSubject), opts...)
}

// gravitySubjects returns the factory of the subjects of Gravity that
// calls newSubject with the handle of each check.
func gravitySubjects[T Gravity](
	newSubject func(h *GravityHandle) T,
) func(h *GravityHandle) Gravity {
	return func(h *GravityHandle) Gravity { return newSubject(h) }
}

// gravityGetCalls generates the calls of Gravity.Get with the inputs that
// each case draws.
var gravityGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Gravity) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Gravity) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// gravityOverEmbedded is Gravity under the name of the field that
// gravityOver embeds.
type gravityOverEmbedded = Gravity

// gravityOver is a subject of Gravity whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type gravityOver struct {
	gravityOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o gravityOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newGravityOver returns the planted defect of a ceiling of Gravity, whose
// callables call spend.
func newGravityOver(_ *GravityHandle, spend func()) Gravity {
	return gravityOver{spend: spend}
}

// gravityGetPanicsDefect is the planted defect of a Gravity whose Get
// panics.
type gravityGetPanicsDefect struct {
	gravityOverEmbedded
}

// Get is the override of the planted defect.
func (*gravityGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Gravity whose Get panics")
}

// gravityGetPanics returns the planted defect of a Gravity whose Get
// panics over a fresh subject.
func gravityGetPanics(h *GravityHandle) Gravity {
	return &gravityGetPanicsDefect{gravityOverEmbedded: h.Fresh()}
}

// gravityGetIgnoresCancelDefect is the planted defect of a Gravity whose
// Get ignores the context.
type gravityGetIgnoresCancelDefect struct {
	gravityOverEmbedded
}

// Get is the override of the planted defect.
func (d *gravityGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.gravityOverEmbedded.Get(p0, p1)
}

// gravityGetIgnoresCancel returns the planted defect of a Gravity whose
// Get ignores the context over a fresh subject.
func gravityGetIgnoresCancel(h *GravityHandle) Gravity {
	return &gravityGetIgnoresCancelDefect{gravityOverEmbedded: h.Fresh()}
}

// gravityGetSwallowsFailureDefect is the planted defect of a Gravity whose
// Get returns no error when its dependency fails.
type gravityGetSwallowsFailureDefect struct {
	gravityOverEmbedded
	h *GravityHandle
}

// Get is the override of the planted defect.
func (d *gravityGetSwallowsFailureDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.gravityOverEmbedded.Get(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(Value)
		r1 = nil
	}
	return r0, r1
}

// gravityGetSwallowsFailure returns the planted defect of a Gravity whose
// Get returns no error when its dependency fails over a fresh subject.
func gravityGetSwallowsFailure(h *GravityHandle) Gravity {
	return &gravityGetSwallowsFailureDefect{gravityOverEmbedded: h.Fresh(), h: h}
}

// gravityGetValueBesideDefect is the planted defect of a Gravity whose Get
// returns a value beside an error.
type gravityGetValueBesideDefect struct {
	gravityOverEmbedded
}

// Get is the override of the planted defect.
func (d *gravityGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.gravityOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// gravityGetValueBeside returns the planted defect of a Gravity whose Get
// returns a value beside an error over a fresh subject.
func gravityGetValueBeside(h *GravityHandle) Gravity {
	return &gravityGetValueBesideDefect{gravityOverEmbedded: h.Fresh()}
}

// The derived checks of Mirror:
//
// Mirror.Get/smoke: Get returns without a panic for derived inputs.
//
// Mirror.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Mirror.Get/cancel: Get returns an error for a cancelled context.
//
// Mirror.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Mirror.Get/failure-is-reported: Get returns an error when a call of its
// dependency fails.
//
// Mirror.Get/failure-is-reported#generated: Get returns an error when a
// call of its dependency fails.
//
// Mirror.Get/zero-on-error: Get returns the zero value beside an error.
//
// Mirror.Get/zero-on-error#generated: Get returns the zero value beside an
// error.

// The derived checks of Mirror that witness their laws at the shallow end:
//
// Mirror.Get/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Mirror writes.
type (
	MirrorHandle = conform.Handle[MirrorMethods]
	MirrorCheck  = conform.Check[MirrorMethods]
	MirrorOption = conform.Option[MirrorMethods]
)

// MirrorMethods is the interface of the exported methods of Mirror. Each
// subject of the suite of Mirror implements it.
type MirrorMethods interface {
	Get(ctx context.Context, key string) (Value, error)
}

// MirrorChecks contains the name of each callable of Mirror, with the
// identities of its checks.
var MirrorChecks = struct {
	Get mirrorGet
}{
	Get: mirrorGet{"Mirror.Get"},
}

// mirrorGet is the name of Mirror.Get, with the identities of its checks.
type mirrorGet struct {
	conform.Callable[MirrorMethods]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c mirrorGet) Smoke() conform.ID[MirrorMethods] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c mirrorGet) SmokeGenerated() conform.ID[MirrorMethods] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c mirrorGet) Cancel() conform.ID[MirrorMethods] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c mirrorGet) CancelGenerated() conform.ID[MirrorMethods] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Get returns an
// error when a call of its dependency fails.
func (c mirrorGet) FailureIsReported() conform.ID[MirrorMethods] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Get
// returns an error when a call of its dependency fails.
func (c mirrorGet) FailureIsReportedGenerated() conform.ID[MirrorMethods] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c mirrorGet) ZeroOnError() conform.ID[MirrorMethods] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c mirrorGet) ZeroOnErrorGenerated() conform.ID[MirrorMethods] {
	return c.Check("zero-on-error#generated")
}

// MirrorSuite returns the checks of Mirror, and the calls and the budgets
// of its callables.
func MirrorSuite() conform.Suite[MirrorMethods] {
	return conform.Suite[MirrorMethods]{
		Declaration: "Mirror",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[MirrorMethods]{
			{
				ID:          MirrorChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MirrorHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: mirrorGetPanics,
			},
			{
				ID:          MirrorChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *MirrorHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: mirrorGetPanics,
			},
			{
				ID:          MirrorChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MirrorHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: mirrorGetIgnoresCancel,
			},
			{
				ID:          MirrorChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *MirrorHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: mirrorGetIgnoresCancel,
			},
			{
				ID:        MirrorChecks.Get.FailureIsReported(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *MirrorHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: mirrorGetSwallowsFailure,
			},
			{
				ID:        MirrorChecks.Get.FailureIsReportedGenerated(),
				Statement: "Get returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *MirrorHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: mirrorGetSwallowsFailure,
			},
			{
				ID:          MirrorChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MirrorHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: mirrorGetValueBeside,
			},
			{
				ID:          MirrorChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *MirrorHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Value{})
				},
				Defect: mirrorGetValueBeside,
			},
		},
		Calls: []conform.Call[MirrorMethods]{
			{
				Calls:    mirrorGetCalls,
				Over:     newMirrorOver,
				Callable: MirrorChecks.Get.Callable,
			},
		},
	}
}

// RunMirrorSuite applies the checks of Mirror, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunMirrorSuite[T MirrorMethods](t *testing.T, newSubject func() T, opts ...MirrorOption) {
	t.Helper()
	RunMirrorSuiteFromHandle(t, func(*MirrorHandle) T { return newSubject() }, opts...)
}

// RunMirrorSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunMirrorSuiteFromHandle[T MirrorMethods](
	t *testing.T,
	newSubject func(h *MirrorHandle) T,
	opts ...MirrorOption,
) {
	t.Helper()
	conform.Run(t, MirrorSuite(), mirrorSubjects(newSubject), opts...)
}

// FuzzMirrorSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzMirrorSuite[T MirrorMethods](
	f *testing.F,
	callable conform.Callable[MirrorMethods],
	newSubject func(h *MirrorHandle) T,
	opts ...MirrorOption,
) {
	f.Helper()
	conform.Fuzz(f, MirrorSuite(), callable, mirrorSubjects(newSubject), opts...)
}

// RunMirrorBudgets applies the allocation ceilings of Mirror and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunMirrorBudgets[T MirrorMethods](
	t *testing.T,
	newSubject func(h *MirrorHandle) T,
	opts ...MirrorOption,
) {
	t.Helper()
	conform.RunBudgets(t, MirrorSuite(), mirrorSubjects(newSubject), opts...)
}

// BenchMirrorSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchMirrorSuite[T MirrorMethods](
	b conform.B,
	callable conform.Callable[MirrorMethods],
	newSubject func(h *MirrorHandle) T,
	opts ...MirrorOption,
) {
	b.Helper()
	conform.Bench(b, MirrorSuite(), callable, mirrorSubjects(newSubject), opts...)
}

// mirrorSubjects returns the factory of the subjects of Mirror that calls
// newSubject with the handle of each check.
func mirrorSubjects[T MirrorMethods](
	newSubject func(h *MirrorHandle) T,
) func(h *MirrorHandle) MirrorMethods {
	return func(h *MirrorHandle) MirrorMethods { return newSubject(h) }
}

// mirrorGetCalls generates the calls of Mirror.Get with the inputs that
// each case draws.
var mirrorGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, MirrorMethods) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject MirrorMethods) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// mirrorOverEmbedded is MirrorMethods under the name of the field that
// mirrorOver embeds.
type mirrorOverEmbedded = MirrorMethods

// mirrorOver is a subject of Mirror whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type mirrorOver struct {
	mirrorOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o mirrorOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newMirrorOver returns the planted defect of a ceiling of Mirror, whose
// callables call spend.
func newMirrorOver(_ *MirrorHandle, spend func()) MirrorMethods {
	return mirrorOver{spend: spend}
}

// mirrorGetPanicsDefect is the planted defect of a Mirror whose Get
// panics.
type mirrorGetPanicsDefect struct {
	mirrorOverEmbedded
}

// Get is the override of the planted defect.
func (*mirrorGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Mirror whose Get panics")
}

// mirrorGetPanics returns the planted defect of a Mirror whose Get panics
// over a fresh subject.
func mirrorGetPanics(h *MirrorHandle) MirrorMethods {
	return &mirrorGetPanicsDefect{mirrorOverEmbedded: h.Fresh()}
}

// mirrorGetIgnoresCancelDefect is the planted defect of a Mirror whose Get
// ignores the context.
type mirrorGetIgnoresCancelDefect struct {
	mirrorOverEmbedded
}

// Get is the override of the planted defect.
func (d *mirrorGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.mirrorOverEmbedded.Get(p0, p1)
}

// mirrorGetIgnoresCancel returns the planted defect of a Mirror whose Get
// ignores the context over a fresh subject.
func mirrorGetIgnoresCancel(h *MirrorHandle) MirrorMethods {
	return &mirrorGetIgnoresCancelDefect{mirrorOverEmbedded: h.Fresh()}
}

// mirrorGetSwallowsFailureDefect is the planted defect of a Mirror whose
// Get returns no error when its dependency fails.
type mirrorGetSwallowsFailureDefect struct {
	mirrorOverEmbedded
	h *MirrorHandle
}

// Get is the override of the planted defect.
func (d *mirrorGetSwallowsFailureDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.mirrorOverEmbedded.Get(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(Value)
		r1 = nil
	}
	return r0, r1
}

// mirrorGetSwallowsFailure returns the planted defect of a Mirror whose
// Get returns no error when its dependency fails over a fresh subject.
func mirrorGetSwallowsFailure(h *MirrorHandle) MirrorMethods {
	return &mirrorGetSwallowsFailureDefect{mirrorOverEmbedded: h.Fresh(), h: h}
}

// mirrorGetValueBesideDefect is the planted defect of a Mirror whose Get
// returns a value beside an error.
type mirrorGetValueBesideDefect struct {
	mirrorOverEmbedded
}

// Get is the override of the planted defect.
func (d *mirrorGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.mirrorOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// mirrorGetValueBeside returns the planted defect of a Mirror whose Get
// returns a value beside an error over a fresh subject.
func mirrorGetValueBeside(h *MirrorHandle) MirrorMethods {
	return &mirrorGetValueBesideDefect{mirrorOverEmbedded: h.Fresh()}
}

// The derived checks of Fetch:
//
// Fetch/smoke: Fetch returns without a panic for derived inputs.
//
// Fetch/smoke#generated: Fetch returns without a panic for derived inputs.
//
// Fetch/cancel: Fetch returns an error for a cancelled context.
//
// Fetch/cancel#generated: Fetch returns an error for a cancelled context.
//
// Fetch/failure-is-reported: Fetch returns an error when a call of its
// dependency fails.
//
// Fetch/failure-is-reported#generated: Fetch returns an error when a call
// of its dependency fails.
//
// Fetch/zero-on-error: Fetch returns the zero value beside an error.
//
// Fetch/zero-on-error#generated: Fetch returns the zero value beside an
// error.
//
// Fetch/retry-succeeds: Fetch succeeds within the declared attempts.
//
// Fetch/retry-succeeds#generated: Fetch succeeds within the declared
// attempts.
//
// Fetch/retry-succeeds-under-faults: Fetch succeeds within the declared
// attempts after a call of its dependency fails.
//
// Fetch/retry-succeeds-under-faults#generated: Fetch succeeds within the
// declared attempts after a call of its dependency fails.

// The derived checks of Fetch that witness their laws at the shallow end:
//
// Fetch/smoke
//
// Fetch/retry-succeeds
//
// Fetch/retry-succeeds#generated

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Fetch writes.
type (
	FetchHandle = conform.Handle[FetchFunc]
	FetchCheck  = conform.Check[FetchFunc]
	FetchOption = conform.Option[FetchFunc]
)

// FetchFunc is the type of the function Fetch. Each subject of the suite
// of Fetch has it.
type FetchFunc func(ctx context.Context, key string) (string, error)

// FetchChecks contains the name of each callable of Fetch, with the
// identities of its checks.
var FetchChecks = fetchCallable{"Fetch"}

// fetchCallable is the name of Fetch, with the identities of its checks.
type fetchCallable struct {
	conform.Callable[FetchFunc]
}

// Smoke returns the identity of the check that Fetch returns without a
// panic for derived inputs.
func (c fetchCallable) Smoke() conform.ID[FetchFunc] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Fetch returns
// without a panic for derived inputs.
func (c fetchCallable) SmokeGenerated() conform.ID[FetchFunc] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Fetch returns an error for
// a cancelled context.
func (c fetchCallable) Cancel() conform.ID[FetchFunc] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Fetch returns an
// error for a cancelled context.
func (c fetchCallable) CancelGenerated() conform.ID[FetchFunc] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Fetch returns
// an error when a call of its dependency fails.
func (c fetchCallable) FailureIsReported() conform.ID[FetchFunc] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Fetch
// returns an error when a call of its dependency fails.
func (c fetchCallable) FailureIsReportedGenerated() conform.ID[FetchFunc] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Fetch returns the
// zero value beside an error.
func (c fetchCallable) ZeroOnError() conform.ID[FetchFunc] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Fetch
// returns the zero value beside an error.
func (c fetchCallable) ZeroOnErrorGenerated() conform.ID[FetchFunc] {
	return c.Check("zero-on-error#generated")
}

// RetrySucceeds returns the identity of the check that Fetch succeeds
// within the declared attempts.
func (c fetchCallable) RetrySucceeds() conform.ID[FetchFunc] {
	return c.Check("retry-succeeds")
}

// RetrySucceedsGenerated returns the identity of the check that Fetch
// succeeds within the declared attempts.
func (c fetchCallable) RetrySucceedsGenerated() conform.ID[FetchFunc] {
	return c.Check("retry-succeeds#generated")
}

// RetrySucceedsUnderFaults returns the identity of the check that Fetch
// succeeds within the declared attempts after a call of its dependency
// fails.
func (c fetchCallable) RetrySucceedsUnderFaults() conform.ID[FetchFunc] {
	return c.Check("retry-succeeds-under-faults")
}

// RetrySucceedsUnderFaultsGenerated returns the identity of the check that
// Fetch succeeds within the declared attempts after a call of its
// dependency fails.
func (c fetchCallable) RetrySucceedsUnderFaultsGenerated() conform.ID[FetchFunc] {
	return c.Check("retry-succeeds-under-faults#generated")
}

// FetchSuite returns the checks of Fetch, and the calls and the budgets of
// its callables.
func FetchSuite() conform.Suite[FetchFunc] {
	return conform.Suite[FetchFunc]{
		Declaration: "Fetch",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[FetchFunc]{
			{
				ID:          FetchChecks.Smoke(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FetchHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					})
				},
				Defect: fetchPanics,
			},
			{
				ID:          FetchChecks.SmokeGenerated(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FetchHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					})
				},
				Defect: fetchPanics,
			},
			{
				ID:          FetchChecks.Cancel(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *FetchHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					})
				},
				Defect: fetchIgnoresCancel,
			},
			{
				ID:          FetchChecks.CancelGenerated(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FetchHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					})
				},
				Defect: fetchIgnoresCancel,
			},
			{
				ID:        FetchChecks.FailureIsReported(),
				Statement: "Fetch returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *FetchHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					})
				},
				Defect: fetchSwallowsFailure,
			},
			{
				ID:        FetchChecks.FailureIsReportedGenerated(),
				Statement: "Fetch returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *FetchHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					})
				},
				Defect: fetchSwallowsFailure,
			},
			{
				ID:          FetchChecks.ZeroOnError(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *FetchHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject()(ctx, string("key-a"))
					}, string(""))
				},
				Defect: fetchValueBeside,
			},
			{
				ID:          FetchChecks.ZeroOnErrorGenerated(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FetchHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject()(ctx, drawnKey)
					}, string(""))
				},
				Defect: fetchValueBeside,
			},
			{
				ID:          FetchChecks.RetrySucceeds(),
				Statement:   "Fetch succeeds within the declared attempts",
				Claim:       spec.RetrySucceeds,
				Probe:       spec.SucceedsWithin,
				Reason:      spec.NeverConverged,
				Law:         spec.CallConvergesWithinTheDeclaredAttempts,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FetchHandle) {
					probe.SucceedsWithin(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					}, 3)
				},
				Defect: fetchFailsAlways,
			},
			{
				ID:          FetchChecks.RetrySucceedsGenerated(),
				Statement:   "Fetch succeeds within the declared attempts",
				Claim:       spec.RetrySucceeds,
				Probe:       spec.SucceedsWithin,
				Reason:      spec.NeverConverged,
				Law:         spec.CallConvergesWithinTheDeclaredAttempts,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FetchHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.SucceedsWithin(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					}, 3)
				},
				Defect: fetchFailsAlways,
			},
			{
				ID:        FetchChecks.RetrySucceedsUnderFaults(),
				Statement: "Fetch succeeds within the declared attempts after a call of its dependency fails",
				Claim:     spec.RetrySucceedsUnderFaults,
				Probe:     spec.SucceedsAfterAFailure,
				Reason:    spec.NeverConverged,
				Law:       spec.CallConvergesWithinTheDeclaredAttempts,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Entails: []conform.ID[FetchFunc]{FetchChecks.RetrySucceeds()},
				Apply: func(h *FetchHandle) {
					probe.SucceedsAfterAFailure(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					}, 3)
				},
				Defect: fetchGivesUpEarly,
			},
			{
				ID:        FetchChecks.RetrySucceedsUnderFaultsGenerated(),
				Statement: "Fetch succeeds within the declared attempts after a call of its dependency fails",
				Claim:     spec.RetrySucceedsUnderFaults,
				Probe:     spec.SucceedsAfterAFailure,
				Reason:    spec.NeverConverged,
				Law:       spec.CallConvergesWithinTheDeclaredAttempts,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Entails: []conform.ID[FetchFunc]{FetchChecks.RetrySucceedsGenerated()},
				Property: func(h *FetchHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.SucceedsAfterAFailure(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					}, 3)
				},
				Defect: fetchGivesUpEarly,
			},
		},
		Calls: []conform.Call[FetchFunc]{
			{
				Calls:    fetchCalls,
				Over:     newFetchOver,
				Callable: FetchChecks.Callable,
			},
		},
	}
}

// RunFetchSuite applies the checks of Fetch, and the checks that opts add,
// to subject.
func RunFetchSuite(t *testing.T, subject FetchFunc, opts ...FetchOption) {
	t.Helper()
	RunFetchSuiteFromHandle(t, fetchSubjects(subject), opts...)
}

// RunFetchSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunFetchSuiteFromHandle(
	t *testing.T,
	newSubject func(h *FetchHandle) FetchFunc,
	opts ...FetchOption,
) {
	t.Helper()
	conform.Run(t, FetchSuite(), newSubject, opts...)
}

// FuzzFetchSuite registers on f the fuzz target of the property checks of
// Fetch on subject.
func FuzzFetchSuite(f *testing.F, subject FetchFunc, opts ...FetchOption) {
	f.Helper()
	conform.Fuzz(f, FetchSuite(), FetchChecks.Callable, fetchSubjects(subject), opts...)
}

// RunFetchBudgets applies the allocation ceiling of Fetch and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunFetchBudgets(t *testing.T, subject FetchFunc, opts ...FetchOption) {
	t.Helper()
	conform.RunBudgets(t, FetchSuite(), fetchSubjects(subject), opts...)
}

// BenchFetchSuite measures subject under the ceilings of the budget of
// Fetch.
func BenchFetchSuite(b conform.B, subject FetchFunc, opts ...FetchOption) {
	b.Helper()
	conform.Bench(b, FetchSuite(), FetchChecks.Callable, fetchSubjects(subject), opts...)
}

// fetchSubjects returns the factory of the subjects of Fetch that returns
// subject for each handle.
func fetchSubjects(subject FetchFunc) func(*FetchHandle) FetchFunc {
	return func(*FetchHandle) FetchFunc { return subject }
}

// fetchCalls generates the calls of Fetch with the inputs that each case
// draws.
var fetchCalls = prop.Composite(func(c *prop.Case) func(context.Context, FetchFunc) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject FetchFunc) error {
		_, err := subject(ctx, key)
		return err
	}
})

// newFetchOver returns the planted defect of a ceiling of Fetch, which
// calls spend and returns its zero results.
func newFetchOver(_ *FetchHandle, spend func()) FetchFunc {
	return func(context.Context, string) (r0 string, r1 error) {
		spend()
		return r0, r1
	}
}

// fetchPanics returns the planted defect of a Fetch whose Fetch panics
// over the fresh function.
func fetchPanics(h *FetchHandle) FetchFunc {
	return func(p0 context.Context, p1 string) (r0 string, r1 error) {
		panic("the planted defect of a Fetch whose Fetch panics")
	}
}

// fetchIgnoresCancel returns the planted defect of a Fetch whose Fetch
// ignores the context over the fresh function.
func fetchIgnoresCancel(h *FetchHandle) FetchFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 string, r1 error) {
		if p0.Err() != nil {
			return r0, r1
		}
		return fresh(p0, p1)
	}
}

// fetchSwallowsFailure returns the planted defect of a Fetch whose Fetch
// returns no error when its dependency fails over the fresh function.
func fetchSwallowsFailure(h *FetchHandle) FetchFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 string, r1 error) {
		if err := p0.Err(); err != nil {
			r1 = err
			return r0, r1
		}
		before := h.Injector().Failures()
		r0, r1 = fresh(p0, p1)
		if r1 != nil && h.Injector().Failures() > before {
			r0 = *new(string)
			r1 = nil
		}
		return r0, r1
	}
}

// fetchValueBeside returns the planted defect of a Fetch whose Fetch
// returns a value beside an error over the fresh function.
func fetchValueBeside(h *FetchHandle) FetchFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 string, r1 error) {
		r0, r1 = fresh(p0, p1)
		if r1 != nil {
			r0 = string("sample-b")
		}
		return r0, r1
	}
}

// fetchFailsAlways returns the planted defect of a Fetch whose Fetch never
// succeeds over the fresh function.
func fetchFailsAlways(h *FetchHandle) FetchFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 string, r1 error) {
		if err := p0.Err(); err != nil {
			r1 = err
			return r0, r1
		}
		r0, r1 = fresh(p0, p1)
		if r1 == nil {
			r0 = *new(string)
			r1 = conform.ErrGenerated
		}
		return r0, r1
	}
}

// fetchGivesUpEarly returns the planted defect of a Fetch whose Fetch
// keeps failing once its dependency failed over the fresh function.
func fetchGivesUpEarly(h *FetchHandle) FetchFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 string, r1 error) {
		if err := p0.Err(); err != nil {
			r1 = err
			return r0, r1
		}
		r0, r1 = fresh(p0, p1)
		if h.Injector().Failures() > 0 {
			r0 = *new(string)
			r1 = conform.ErrGenerated
		}
		return r0, r1
	}
}

// The derived checks of Pick:
//
// Pick/smoke: Pick returns without a panic for derived inputs.
//
// Pick/smoke#generated: Pick returns without a panic for derived inputs.
//
// Pick/cancel: Pick returns an error for a cancelled context.
//
// Pick/cancel#generated: Pick returns an error for a cancelled context.
//
// Pick/failure-is-reported: Pick returns an error when a call of its
// dependency fails.
//
// Pick/failure-is-reported#generated: Pick returns an error when a call of
// its dependency fails.
//
// Pick/keeps-its-instance: Pick sends each call with one key to one
// instance.
//
// Pick/keeps-its-instance#generated: Pick sends each call with one key to
// one instance.

// The derived checks of Pick that witness their laws at the shallow end:
//
// Pick/smoke

// The unclaimed classifications of Pick:
//
// The shape:writer of Pick has the stance silent. A write is an input of
// the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Pick writes.
type (
	PickHandle = conform.Handle[PickFunc]
	PickCheck  = conform.Check[PickFunc]
	PickOption = conform.Option[PickFunc]
)

// PickFunc is the type of the function Pick. Each subject of the suite of
// Pick has it.
type PickFunc func(ctx context.Context, key string) error

// PickChecks contains the name of each callable of Pick, with the
// identities of its checks.
var PickChecks = pickCallable{"Pick"}

// pickCallable is the name of Pick, with the identities of its checks.
type pickCallable struct {
	conform.Callable[PickFunc]
}

// Smoke returns the identity of the check that Pick returns without a
// panic for derived inputs.
func (c pickCallable) Smoke() conform.ID[PickFunc] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Pick returns
// without a panic for derived inputs.
func (c pickCallable) SmokeGenerated() conform.ID[PickFunc] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Pick returns an error for
// a cancelled context.
func (c pickCallable) Cancel() conform.ID[PickFunc] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Pick returns an
// error for a cancelled context.
func (c pickCallable) CancelGenerated() conform.ID[PickFunc] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Pick returns an
// error when a call of its dependency fails.
func (c pickCallable) FailureIsReported() conform.ID[PickFunc] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Pick
// returns an error when a call of its dependency fails.
func (c pickCallable) FailureIsReportedGenerated() conform.ID[PickFunc] {
	return c.Check("failure-is-reported#generated")
}

// KeepsItsInstance returns the identity of the check that Pick sends each
// call with one key to one instance.
func (c pickCallable) KeepsItsInstance() conform.ID[PickFunc] {
	return c.Check("keeps-its-instance")
}

// KeepsItsInstanceGenerated returns the identity of the check that Pick
// sends each call with one key to one instance.
func (c pickCallable) KeepsItsInstanceGenerated() conform.ID[PickFunc] {
	return c.Check("keeps-its-instance#generated")
}

// PickSuite returns the checks of Pick, and the calls and the budgets of
// its callables.
func PickSuite() conform.Suite[PickFunc] {
	return conform.Suite[PickFunc]{
		Declaration: "Pick",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[PickFunc]{
			{
				ID:          PickChecks.Smoke(),
				Statement:   "Pick returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PickHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject()(ctx, string("key-a"))
					})
				},
				Defect: pickPanics,
			},
			{
				ID:          PickChecks.SmokeGenerated(),
				Statement:   "Pick returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PickHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject()(ctx, drawnKey)
					})
				},
				Defect: pickPanics,
			},
			{
				ID:          PickChecks.Cancel(),
				Statement:   "Pick returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PickHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject()(ctx, string("key-a"))
					})
				},
				Defect: pickIgnoresCancel,
			},
			{
				ID:          PickChecks.CancelGenerated(),
				Statement:   "Pick returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PickHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject()(ctx, drawnKey)
					})
				},
				Defect: pickIgnoresCancel,
			},
			{
				ID:        PickChecks.FailureIsReported(),
				Statement: "Pick returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *PickHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject()(ctx, string("key-a"))
					})
				},
				Defect: pickSwallowsFailure,
			},
			{
				ID:        PickChecks.FailureIsReportedGenerated(),
				Statement: "Pick returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *PickHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject()(ctx, drawnKey)
					})
				},
				Defect: pickSwallowsFailure,
			},
			{
				ID:        PickChecks.KeepsItsInstance(),
				Statement: "Pick sends each call with one key to one instance",
				Claim:     spec.KeepsItsInstance,
				Probe:     spec.UsesOneInstance,
				Reason:    spec.MovedInstance,
				Law:       spec.OneKeyOneInstance,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *PickHandle) {
					probe.UsesOneInstance(h, func(ctx context.Context) error {
						return h.Subject()(ctx, string("key-a"))
					})
				},
				Defect: pickBalancesEachCall,
			},
			{
				ID:        PickChecks.KeepsItsInstanceGenerated(),
				Statement: "Pick sends each call with one key to one instance",
				Claim:     spec.KeepsItsInstance,
				Probe:     spec.UsesOneInstance,
				Reason:    spec.MovedInstance,
				Law:       spec.OneKeyOneInstance,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *PickHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UsesOneInstance(h, func(ctx context.Context) error {
						return h.Subject()(ctx, drawnKey)
					})
				},
				Defect: pickBalancesEachCall,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Pick",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[PickFunc]{
			{
				Calls:    pickCalls,
				Over:     newPickOver,
				Callable: PickChecks.Callable,
			},
		},
	}
}

// RunPickSuite applies the checks of Pick, and the checks that opts add,
// to subject.
func RunPickSuite(t *testing.T, subject PickFunc, opts ...PickOption) {
	t.Helper()
	RunPickSuiteFromHandle(t, pickSubjects(subject), opts...)
}

// RunPickSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPickSuiteFromHandle(
	t *testing.T,
	newSubject func(h *PickHandle) PickFunc,
	opts ...PickOption,
) {
	t.Helper()
	conform.Run(t, PickSuite(), newSubject, opts...)
}

// FuzzPickSuite registers on f the fuzz target of the property checks of
// Pick on subject.
func FuzzPickSuite(f *testing.F, subject PickFunc, opts ...PickOption) {
	f.Helper()
	conform.Fuzz(f, PickSuite(), PickChecks.Callable, pickSubjects(subject), opts...)
}

// RunPickBudgets applies the allocation ceiling of Pick and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunPickBudgets(t *testing.T, subject PickFunc, opts ...PickOption) {
	t.Helper()
	conform.RunBudgets(t, PickSuite(), pickSubjects(subject), opts...)
}

// BenchPickSuite measures subject under the ceilings of the budget of
// Pick.
func BenchPickSuite(b conform.B, subject PickFunc, opts ...PickOption) {
	b.Helper()
	conform.Bench(b, PickSuite(), PickChecks.Callable, pickSubjects(subject), opts...)
}

// pickSubjects returns the factory of the subjects of Pick that returns
// subject for each handle.
func pickSubjects(subject PickFunc) func(*PickHandle) PickFunc {
	return func(*PickHandle) PickFunc { return subject }
}

// pickCalls generates the calls of Pick with the inputs that each case
// draws.
var pickCalls = prop.Composite(func(c *prop.Case) func(context.Context, PickFunc) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject PickFunc) error {
		return subject(ctx, key)
	}
})

// newPickOver returns the planted defect of a ceiling of Pick, which calls
// spend and returns its zero results.
func newPickOver(_ *PickHandle, spend func()) PickFunc {
	return func(context.Context, string) (r0 error) {
		spend()
		return r0
	}
}

// pickPanics returns the planted defect of a Pick whose Pick panics over
// the fresh function.
func pickPanics(h *PickHandle) PickFunc {
	return func(p0 context.Context, p1 string) (r0 error) {
		panic("the planted defect of a Pick whose Pick panics")
	}
}

// pickIgnoresCancel returns the planted defect of a Pick whose Pick
// ignores the context over the fresh function.
func pickIgnoresCancel(h *PickHandle) PickFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 error) {
		if p0.Err() != nil {
			return r0
		}
		return fresh(p0, p1)
	}
}

// pickSwallowsFailure returns the planted defect of a Pick whose Pick
// returns no error when its dependency fails over the fresh function.
func pickSwallowsFailure(h *PickHandle) PickFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 error) {
		if err := p0.Err(); err != nil {
			r0 = err
			return r0
		}
		before := h.Injector().Failures()
		r0 = fresh(p0, p1)
		if r0 != nil && h.Injector().Failures() > before {
			r0 = nil
		}
		return r0
	}
}

// pickBalancesEachCall returns the planted defect of a Pick whose Pick
// sends each call to another instance over the fresh function.
func pickBalancesEachCall(h *PickHandle) PickFunc {
	return func(p0 context.Context, p1 string) (r0 error) {
		if err := p0.Err(); err != nil {
			r0 = err
			return r0
		}
		return h.Fresh()(p0, p1)
	}
}

// The derived checks of Load:
//
// Load/smoke: Load returns without a panic for derived inputs.
//
// Load/smoke#generated: Load returns without a panic for derived inputs.
//
// Load/cancel: Load returns an error for a cancelled context.
//
// Load/cancel#generated: Load returns an error for a cancelled context.
//
// Load/failure-is-reported: Load returns an error when a call of its
// dependency fails.
//
// Load/failure-is-reported#generated: Load returns an error when a call of
// its dependency fails.
//
// Load/zero-on-error: Load returns the zero value beside an error.
//
// Load/zero-on-error#generated: Load returns the zero value beside an
// error.
//
// Load/releases-on-every-path: Load releases each resource that it
// acquired when a call of its dependency fails.
//
// Load/releases-on-every-path#generated: Load releases each resource that
// it acquired when a call of its dependency fails.

// The derived checks of Load that witness their laws at the shallow end:
//
// Load/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Load writes.
type (
	LoadHandle = conform.Handle[LoadFunc]
	LoadCheck  = conform.Check[LoadFunc]
	LoadOption = conform.Option[LoadFunc]
)

// LoadFunc is the type of the function Load. Each subject of the suite of
// Load has it.
type LoadFunc func(ctx context.Context, key string) (Value, error)

// LoadChecks contains the name of each callable of Load, with the
// identities of its checks.
var LoadChecks = loadCallable{"Load"}

// loadCallable is the name of Load, with the identities of its checks.
type loadCallable struct {
	conform.Callable[LoadFunc]
}

// Smoke returns the identity of the check that Load returns without a
// panic for derived inputs.
func (c loadCallable) Smoke() conform.ID[LoadFunc] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Load returns
// without a panic for derived inputs.
func (c loadCallable) SmokeGenerated() conform.ID[LoadFunc] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Load returns an error for
// a cancelled context.
func (c loadCallable) Cancel() conform.ID[LoadFunc] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Load returns an
// error for a cancelled context.
func (c loadCallable) CancelGenerated() conform.ID[LoadFunc] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Load returns an
// error when a call of its dependency fails.
func (c loadCallable) FailureIsReported() conform.ID[LoadFunc] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Load
// returns an error when a call of its dependency fails.
func (c loadCallable) FailureIsReportedGenerated() conform.ID[LoadFunc] {
	return c.Check("failure-is-reported#generated")
}

// ZeroOnError returns the identity of the check that Load returns the zero
// value beside an error.
func (c loadCallable) ZeroOnError() conform.ID[LoadFunc] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Load returns
// the zero value beside an error.
func (c loadCallable) ZeroOnErrorGenerated() conform.ID[LoadFunc] {
	return c.Check("zero-on-error#generated")
}

// ReleasesOnEveryPath returns the identity of the check that Load releases
// each resource that it acquired when a call of its dependency fails.
func (c loadCallable) ReleasesOnEveryPath() conform.ID[LoadFunc] {
	return c.Check("releases-on-every-path")
}

// ReleasesOnEveryPathGenerated returns the identity of the check that Load
// releases each resource that it acquired when a call of its dependency
// fails.
func (c loadCallable) ReleasesOnEveryPathGenerated() conform.ID[LoadFunc] {
	return c.Check("releases-on-every-path#generated")
}

// LoadSuite returns the checks of Load, and the calls and the budgets of
// its callables.
func LoadSuite() conform.Suite[LoadFunc] {
	return conform.Suite[LoadFunc]{
		Declaration: "Load",
		Package:     "example.com/shop/faults",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[LoadFunc]{
			{
				ID:          LoadChecks.Smoke(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LoadHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					})
				},
				Defect: loadPanics,
			},
			{
				ID:          LoadChecks.SmokeGenerated(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LoadHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					})
				},
				Defect: loadPanics,
			},
			{
				ID:          LoadChecks.Cancel(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LoadHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					})
				},
				Defect: loadIgnoresCancel,
			},
			{
				ID:          LoadChecks.CancelGenerated(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LoadHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					})
				},
				Defect: loadIgnoresCancel,
			},
			{
				ID:        LoadChecks.FailureIsReported(),
				Statement: "Load returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *LoadHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					})
				},
				Defect: loadSwallowsFailure,
			},
			{
				ID:        LoadChecks.FailureIsReportedGenerated(),
				Statement: "Load returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *LoadHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					})
				},
				Defect: loadSwallowsFailure,
			},
			{
				ID:          LoadChecks.ZeroOnError(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LoadHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject()(ctx, string("key-a"))
					}, Value{})
				},
				Defect: loadValueBeside,
			},
			{
				ID:          LoadChecks.ZeroOnErrorGenerated(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LoadHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject()(ctx, drawnKey)
					}, Value{})
				},
				Defect: loadValueBeside,
			},
			{
				ID:        LoadChecks.ReleasesOnEveryPath(),
				Statement: "Load releases each resource that it acquired when a call of its dependency fails",
				Claim:     spec.ReleasesOnEveryPath,
				Probe:     spec.ReleasesWhatItAcquires,
				Reason:    spec.LeakedResource,
				Law:       spec.FailurePathReleases,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *LoadHandle) {
					probe.ReleasesWhatItAcquires(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, string("key-a"))
						return err
					}, "Connect", "Disconnect")
				},
				Defect: loadOpensWithoutClosing,
			},
			{
				ID:        LoadChecks.ReleasesOnEveryPathGenerated(),
				Statement: "Load releases each resource that it acquired when a call of its dependency fails",
				Claim:     spec.ReleasesOnEveryPath,
				Probe:     spec.ReleasesWhatItAcquires,
				Reason:    spec.LeakedResource,
				Law:       spec.FailurePathReleases,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *LoadHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReleasesWhatItAcquires(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnKey)
						return err
					}, "Connect", "Disconnect")
				},
				Defect: loadOpensWithoutClosing,
			},
		},
		Calls: []conform.Call[LoadFunc]{
			{
				Calls:    loadCalls,
				Over:     newLoadOver,
				Callable: LoadChecks.Callable,
			},
		},
	}
}

// RunLoadSuite applies the checks of Load, and the checks that opts add,
// to subject.
func RunLoadSuite(t *testing.T, subject LoadFunc, opts ...LoadOption) {
	t.Helper()
	RunLoadSuiteFromHandle(t, loadSubjects(subject), opts...)
}

// RunLoadSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLoadSuiteFromHandle(
	t *testing.T,
	newSubject func(h *LoadHandle) LoadFunc,
	opts ...LoadOption,
) {
	t.Helper()
	conform.Run(t, LoadSuite(), newSubject, opts...)
}

// FuzzLoadSuite registers on f the fuzz target of the property checks of
// Load on subject.
func FuzzLoadSuite(f *testing.F, subject LoadFunc, opts ...LoadOption) {
	f.Helper()
	conform.Fuzz(f, LoadSuite(), LoadChecks.Callable, loadSubjects(subject), opts...)
}

// RunLoadBudgets applies the allocation ceiling of Load and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunLoadBudgets(t *testing.T, subject LoadFunc, opts ...LoadOption) {
	t.Helper()
	conform.RunBudgets(t, LoadSuite(), loadSubjects(subject), opts...)
}

// BenchLoadSuite measures subject under the ceilings of the budget of
// Load.
func BenchLoadSuite(b conform.B, subject LoadFunc, opts ...LoadOption) {
	b.Helper()
	conform.Bench(b, LoadSuite(), LoadChecks.Callable, loadSubjects(subject), opts...)
}

// loadSubjects returns the factory of the subjects of Load that returns
// subject for each handle.
func loadSubjects(subject LoadFunc) func(*LoadHandle) LoadFunc {
	return func(*LoadHandle) LoadFunc { return subject }
}

// loadCalls generates the calls of Load with the inputs that each case
// draws.
var loadCalls = prop.Composite(func(c *prop.Case) func(context.Context, LoadFunc) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject LoadFunc) error {
		_, err := subject(ctx, key)
		return err
	}
})

// newLoadOver returns the planted defect of a ceiling of Load, which calls
// spend and returns its zero results.
func newLoadOver(_ *LoadHandle, spend func()) LoadFunc {
	return func(context.Context, string) (r0 Value, r1 error) {
		spend()
		return r0, r1
	}
}

// loadPanics returns the planted defect of a Load whose Load panics over
// the fresh function.
func loadPanics(h *LoadHandle) LoadFunc {
	return func(p0 context.Context, p1 string) (r0 Value, r1 error) {
		panic("the planted defect of a Load whose Load panics")
	}
}

// loadIgnoresCancel returns the planted defect of a Load whose Load
// ignores the context over the fresh function.
func loadIgnoresCancel(h *LoadHandle) LoadFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 Value, r1 error) {
		if p0.Err() != nil {
			return r0, r1
		}
		return fresh(p0, p1)
	}
}

// loadSwallowsFailure returns the planted defect of a Load whose Load
// returns no error when its dependency fails over the fresh function.
func loadSwallowsFailure(h *LoadHandle) LoadFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 Value, r1 error) {
		if err := p0.Err(); err != nil {
			r1 = err
			return r0, r1
		}
		before := h.Injector().Failures()
		r0, r1 = fresh(p0, p1)
		if r1 != nil && h.Injector().Failures() > before {
			r0 = *new(Value)
			r1 = nil
		}
		return r0, r1
	}
}

// loadValueBeside returns the planted defect of a Load whose Load returns
// a value beside an error over the fresh function.
func loadValueBeside(h *LoadHandle) LoadFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 string) (r0 Value, r1 error) {
		r0, r1 = fresh(p0, p1)
		if r1 != nil {
			r0 = Value{Key: "Key-b"}
		}
		return r0, r1
	}
}

// loadOpensWithoutClosing returns the planted defect of a Load whose Load
// keeps a resource that it acquired before a failure over the fresh
// function.
func loadOpensWithoutClosing(h *LoadHandle) LoadFunc {
	fresh := h.Fresh()
	leak := h.Injector().Dependency()
	return func(p0 context.Context, p1 string) (r0 Value, r1 error) {
		if err := p0.Err(); err != nil {
			r1 = err
			return r0, r1
		}
		r0, r1 = fresh(p0, p1)
		if r1 != nil {
			_ = leak.Call("Connect")
		}
		return r0, r1
	}
}
