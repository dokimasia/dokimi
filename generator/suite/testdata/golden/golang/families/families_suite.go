package families

import (
	"context"
	"iter"
	"slices"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Store:
//
// Store.Put/smoke: Put returns without a panic for derived inputs.
//
// Store.Put/cancel: Put returns an error for a cancelled context.
//
// Store.Get/smoke: Get returns without a panic for derived inputs.
//
// Store.Get/cancel: Get returns an error for a cancelled context.
//
// Store.Get/zero-on-error: Get returns the zero value beside an error.
//
// Store.Get/read-after-write: Get returns what Put stored before it.
//
// Store.Get/monotonic-reads: Get never returns an older version than one
// that it returned before.
//
// Store.Load/smoke: Load returns without a panic for derived inputs.
//
// Store.Load/zero-on-miss: Load returns the zero value beside a false
// presence flag.
//
// Store.Load/cacheable: Load returns the same value on a repeat.
//
// Store.Split/smoke: Split returns without a panic for derived inputs.
//
// Store.Find/smoke: Find returns without a panic for derived inputs.
//
// Store.Find/cancel: Find returns an error for a cancelled context.
//
// Store.Level/smoke: Level returns without a panic for derived inputs.
//
// Store.Level/cancel: Level returns an error for a cancelled context.
//
// Store.Level/zero-on-error: Level returns the zero value beside an error.
//
// Store.Raise/smoke: Raise returns without a panic for derived inputs.
//
// Store.Raise/cancel: Raise returns an error for a cancelled context.
//
// Store.Inspect/smoke: Inspect returns without a panic for derived inputs.
//
// Store.GetWithMeta/smoke: GetWithMeta returns without a panic for derived
// inputs.
//
// Store.GetWithMeta/cancel: GetWithMeta returns an error for a cancelled
// context.
//
// Store.GetWithMeta/zeroes-on-error: GetWithMeta returns the zero values
// beside an error.

// The derived checks of Store that witness their laws at the shallow end:
//
// Store.Put/smoke
//
// Store.Get/smoke
//
// Store.Get/read-after-write
//
// Store.Get/monotonic-reads
//
// Store.Load/smoke
//
// Store.Load/cacheable
//
// Store.Split/smoke
//
// Store.Find/smoke
//
// Store.Level/smoke
//
// Store.Raise/smoke
//
// Store.Inspect/smoke
//
// Store.GetWithMeta/smoke

// The refused claims of Store:
//
// Store.Load/cancel: Load cannot fail. The remedy is to give Load an error
// result.
//
// Store.Split/pure: Split does not return exactly one value. The remedy is
// to return one value from Split.
//
// Store.Inspect/zeroes-on-miss: Inspect returns no presence flag. The
// remedy is to return a presence flag beside the value of Inspect.

// The unclaimed classifications of Store:
//
// The shape:writer of Store.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:batch-reader of Store.Find has the stance gap. No claim states
// what a batch read returns for a key that it does not find.
//
// The shape:writer of Store.Raise has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Store writes.
type (
	StoreHandle = conform.Handle[Store]
	StoreCheck  = conform.Check[Store]
	StoreOption = conform.Option[Store]
)

// StoreChecks contains the name of each callable of Store, with the
// identities of its checks.
var StoreChecks = struct {
	Put         storePut
	Get         storeGet
	Load        storeLoad
	Split       storeSplit
	Find        storeFind
	Level       storeLevel
	Raise       storeRaise
	Inspect     storeInspect
	GetWithMeta storeGetWithMeta
}{
	Put:         storePut{"Store.Put"},
	Get:         storeGet{"Store.Get"},
	Load:        storeLoad{"Store.Load"},
	Split:       storeSplit{"Store.Split"},
	Find:        storeFind{"Store.Find"},
	Level:       storeLevel{"Store.Level"},
	Raise:       storeRaise{"Store.Raise"},
	Inspect:     storeInspect{"Store.Inspect"},
	GetWithMeta: storeGetWithMeta{"Store.GetWithMeta"},
}

// storePut is the name of Store.Put, with the identities of its checks.
type storePut struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c storePut) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c storePut) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// storeGet is the name of Store.Get, with the identities of its checks.
type storeGet struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c storeGet) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c storeGet) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c storeGet) ZeroOnError() conform.ID[Store] {
	return c.Check("zero-on-error")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c storeGet) ReadAfterWrite() conform.ID[Store] {
	return c.Check("read-after-write")
}

// MonotonicReads returns the identity of the check that Get never returns
// an older version than one that it returned before.
func (c storeGet) MonotonicReads() conform.ID[Store] {
	return c.Check("monotonic-reads")
}

// storeLoad is the name of Store.Load, with the identities of its checks.
type storeLoad struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Load returns without a
// panic for derived inputs.
func (c storeLoad) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// ZeroOnMiss returns the identity of the check that Load returns the zero
// value beside a false presence flag.
func (c storeLoad) ZeroOnMiss() conform.ID[Store] {
	return c.Check("zero-on-miss")
}

// Cacheable returns the identity of the check that Load returns the same
// value on a repeat.
func (c storeLoad) Cacheable() conform.ID[Store] {
	return c.Check("cacheable")
}

// storeSplit is the name of Store.Split, with the identities of its
// checks.
type storeSplit struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Split returns without a
// panic for derived inputs.
func (c storeSplit) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// storeFind is the name of Store.Find, with the identities of its checks.
type storeFind struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Find returns without a
// panic for derived inputs.
func (c storeFind) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Find returns an error for
// a cancelled context.
func (c storeFind) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// storeLevel is the name of Store.Level, with the identities of its
// checks.
type storeLevel struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Level returns without a
// panic for derived inputs.
func (c storeLevel) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Level returns an error for
// a cancelled context.
func (c storeLevel) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Level returns the
// zero value beside an error.
func (c storeLevel) ZeroOnError() conform.ID[Store] {
	return c.Check("zero-on-error")
}

// storeRaise is the name of Store.Raise, with the identities of its
// checks.
type storeRaise struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Raise returns without a
// panic for derived inputs.
func (c storeRaise) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Raise returns an error for
// a cancelled context.
func (c storeRaise) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// storeInspect is the name of Store.Inspect, with the identities of its
// checks.
type storeInspect struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that Inspect returns without a
// panic for derived inputs.
func (c storeInspect) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// storeGetWithMeta is the name of Store.GetWithMeta, with the identities
// of its checks.
type storeGetWithMeta struct {
	conform.Callable[Store]
}

// Smoke returns the identity of the check that GetWithMeta returns without
// a panic for derived inputs.
func (c storeGetWithMeta) Smoke() conform.ID[Store] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that GetWithMeta returns an
// error for a cancelled context.
func (c storeGetWithMeta) Cancel() conform.ID[Store] {
	return c.Check("cancel")
}

// ZeroesOnError returns the identity of the check that GetWithMeta returns
// the zero values beside an error.
func (c storeGetWithMeta) ZeroesOnError() conform.ID[Store] {
	return c.Check("zeroes-on-error")
}

// StoreSuite returns the checks of Store, and the calls and the budgets of
// its callables.
func StoreSuite() conform.Suite[Store] {
	return conform.Suite[Store]{
		Declaration: "Store",
		Checks: []conform.Check[Store]{
			{
				ID:          StoreChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: storePutPanics,
			},
			{
				ID:          StoreChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: storePutIgnoresCancel,
			},
			{
				ID:          StoreChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: storeGetPanics,
			},
			{
				ID:          StoreChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: storeGetIgnoresCancel,
			},
			{
				ID:          StoreChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: storeGetValueBeside,
			},
			{
				ID:          StoreChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) (Value, error) {
						w := Value{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: storeGetReadsAnotherValue,
			},
			{
				ID:          StoreChecks.Get.MonotonicReads(),
				Statement:   "Get never returns an older version than one that it returned before",
				Claim:       spec.MonotonicReads,
				Probe:       spec.ReadsNeverRegress,
				Reason:      spec.Regressed,
				Law:         spec.ReadsNeverGoBackwards,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.ReadsNeverRegress(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					}, func(v Value) int64 { return v.Rev })
				},
				Defect: storeGetVersionDescends,
			},
			{
				ID:          StoreChecks.Load.Smoke(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Load(ctx, string("key-a"))
						return nil
					})
				},
				Defect: storeLoadPanics,
			},
			{
				ID:          StoreChecks.Load.ZeroOnMiss(),
				Statement:   "Load returns the zero value beside a false presence flag",
				Claim:       spec.ZeroOnMiss,
				Probe:       spec.ZeroBesideFalse,
				Reason:      spec.ValueBesideFalse,
				Law:         spec.NoValueBesideAFalseFlag,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.ZeroBesideFalse(h, func(ctx context.Context) (Value, bool) {
						return h.Subject().Load(ctx, string("key-a"))
					}, Value{})
				},
				Defect: storeLoadValueBesideAFalseFlag,
			},
			{
				ID:          StoreChecks.Load.Cacheable(),
				Statement:   "Load returns the same value on a repeat",
				Claim:       spec.Cacheable,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (Value, error) {
						r0, _ := h.Subject().Load(ctx, string("key-a"))
						return r0, nil
					})
				},
				Defect: storeLoadDrifts,
			},
			{
				ID:          StoreChecks.Split.Smoke(),
				Statement:   "Split returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Split(string("in-a"))
						return nil
					})
				},
				Defect: storeSplitPanics,
			},
			{
				ID:          StoreChecks.Find.Smoke(),
				Statement:   "Find returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Find(ctx, string("keys-a"))
						return err
					})
				},
				Defect: storeFindPanics,
			},
			{
				ID:          StoreChecks.Find.Cancel(),
				Statement:   "Find returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Find(ctx, string("keys-a"))
						return err
					})
				},
				Defect: storeFindIgnoresCancel,
			},
			{
				ID:          StoreChecks.Level.Smoke(),
				Statement:   "Level returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Level(ctx, string("key-a"))
						return err
					})
				},
				Defect: storeLevelPanics,
			},
			{
				ID:          StoreChecks.Level.Cancel(),
				Statement:   "Level returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Level(ctx, string("key-a"))
						return err
					})
				},
				Defect: storeLevelIgnoresCancel,
			},
			{
				ID:          StoreChecks.Level.ZeroOnError(),
				Statement:   "Level returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Level, error) {
						return h.Subject().Level(ctx, string("key-a"))
					}, Level(0))
				},
				Defect: storeLevelValueBeside,
			},
			{
				ID:          StoreChecks.Raise.Smoke(),
				Statement:   "Raise returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Raise(ctx, Level(Base))
					})
				},
				Defect: storeRaisePanics,
			},
			{
				ID:          StoreChecks.Raise.Cancel(),
				Statement:   "Raise returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Raise(ctx, Level(Base))
					})
				},
				Defect: storeRaiseIgnoresCancel,
			},
			{
				ID:          StoreChecks.Inspect.Smoke(),
				Statement:   "Inspect returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Inspect(string("key-a"))
						return nil
					})
				},
				Defect: storeInspectPanics,
			},
			{
				ID:          StoreChecks.GetWithMeta.Smoke(),
				Statement:   "GetWithMeta returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StoreHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, _, err := h.Subject().GetWithMeta(ctx, string("key-a"))
						return err
					})
				},
				Defect: storeGetWithMetaPanics,
			},
			{
				ID:          StoreChecks.GetWithMeta.Cancel(),
				Statement:   "GetWithMeta returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, _, err := h.Subject().GetWithMeta(ctx, string("key-a"))
						return err
					})
				},
				Defect: storeGetWithMetaIgnoresCancel,
			},
			{
				ID:          StoreChecks.GetWithMeta.ZeroesOnError(),
				Statement:   "GetWithMeta returns the zero values beside an error",
				Claim:       spec.ZeroesOnError,
				Probe:       spec.ZeroPairBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StoreHandle) {
					probe.ZeroPairBesideError(
						h,
						func(ctx context.Context) (probe.Pair[Value, Meta], error) {
							r0, r1, err := h.Subject().GetWithMeta(ctx, string("key-a"))
							return probe.Pair[Value, Meta]{First: r0, Second: r1}, err
						},
						probe.Pair[Value, Meta]{First: Value{}, Second: Meta{}},
					)
				},
				Defect: storeGetWithMetaValuesBeside,
			},
		},
		Calls: []conform.Call[Store]{
			{
				Calls:    storePutCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Put.Callable,
			},
			{
				Calls:    storeGetCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Get.Callable,
			},
			{
				Calls:    storeLoadCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Load.Callable,
			},
			{
				Calls:    storeSplitCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Split.Callable,
			},
			{
				Calls:    storeFindCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Find.Callable,
			},
			{
				Calls:    storeLevelCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Level.Callable,
			},
			{
				Calls:    storeRaiseCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Raise.Callable,
			},
			{
				Calls:    storeInspectCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.Inspect.Callable,
			},
			{
				Calls:    storeGetWithMetaCalls,
				Over:     newStoreOver,
				Callable: StoreChecks.GetWithMeta.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(storeLevelValues),
		},
	}
}

// RunStoreSuite applies the checks of Store, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunStoreSuite[T Store](t *testing.T, newSubject func() T, opts ...StoreOption) {
	t.Helper()
	RunStoreSuiteFromHandle(t, func(*StoreHandle) T { return newSubject() }, opts...)
}

// RunStoreSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStoreSuiteFromHandle[T Store](
	t *testing.T,
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	t.Helper()
	conform.Run(t, StoreSuite(), storeSubjects(newSubject), opts...)
}

// FuzzStoreSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzStoreSuite[T Store](
	f *testing.F,
	callable conform.Callable[Store],
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	f.Helper()
	conform.Fuzz(f, StoreSuite(), callable, storeSubjects(newSubject), opts...)
}

// RunStoreBudgets applies the allocation ceilings of Store and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunStoreBudgets[T Store](
	t *testing.T,
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	t.Helper()
	conform.RunBudgets(t, StoreSuite(), storeSubjects(newSubject), opts...)
}

// BenchStoreSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchStoreSuite[T Store](
	b conform.B,
	callable conform.Callable[Store],
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
) {
	b.Helper()
	conform.Bench(b, StoreSuite(), callable, storeSubjects(newSubject), opts...)
}

// storeSubjects returns the factory of the subjects of Store that calls
// newSubject with the handle of each check.
func storeSubjects[T Store](newSubject func(h *StoreHandle) T) func(h *StoreHandle) Store {
	return func(h *StoreHandle) Store { return newSubject(h) }
}

// storePutCalls generates the calls of Store.Put with the inputs that each
// case draws.
var storePutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Store) error {
		return subject.Put(ctx, v)
	}
})

// storeGetCalls generates the calls of Store.Get with the inputs that each
// case draws.
var storeGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Store) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// storeLoadCalls generates the calls of Store.Load with the inputs that
// each case draws.
var storeLoadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Store) error {
		_, _ = subject.Load(ctx, key)
		return nil
	}
})

// storeSplitCalls generates the calls of Store.Split with the inputs that
// each case draws.
var storeSplitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	in := c.Draw(prop.Of[string](), "in")
	return func(_ context.Context, subject Store) error {
		_, _ = subject.Split(in)
		return nil
	}
})

// storeFindCalls generates the calls of Store.Find with the inputs that
// each case draws.
var storeFindCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	keys := c.Draw(prop.List(prop.Of[string]()), "keys")
	return func(ctx context.Context, subject Store) error {
		_, err := subject.Find(ctx, keys...)
		return err
	}
})

// storeLevelCalls generates the calls of Store.Level with the inputs that
// each case draws.
var storeLevelCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Store) error {
		_, err := subject.Level(ctx, key)
		return err
	}
})

// storeRaiseCalls generates the calls of Store.Raise with the inputs that
// each case draws.
var storeRaiseCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	l := c.Draw(prop.Of[Level](), "l")
	return func(ctx context.Context, subject Store) error {
		return subject.Raise(ctx, l)
	}
})

// storeInspectCalls generates the calls of Store.Inspect with the inputs
// that each case draws.
var storeInspectCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(_ context.Context, subject Store) error {
		_, _, _ = subject.Inspect(key)
		return nil
	}
})

// storeGetWithMetaCalls generates the calls of Store.GetWithMeta with the
// inputs that each case draws.
var storeGetWithMetaCalls = prop.Composite(func(c *prop.Case) func(context.Context, Store) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Store) error {
		_, _, err := subject.GetWithMeta(ctx, key)
		return err
	}
})

// storeLevelValues generates the values of Level that keep or change its
// starting value.
var storeLevelValues = prop.Composite(func(c *prop.Case) Level {
	changed := c.Draw(prop.Of[Level](), "level")
	if c.Draw(prop.Boolean(), "level.changed") {
		return changed
	}
	return Base
})

// storeOverEmbedded is Store under the name of the field that storeOver
// embeds.
type storeOverEmbedded = Store

// storeOver is a subject of Store whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type storeOver struct {
	storeOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o storeOver) Put(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o storeOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Load calls spend and returns its zero results.
func (o storeOver) Load(context.Context, string) (r0 Value, r1 bool) {
	o.spend()
	return r0, r1
}

// Split calls spend and returns its zero results.
func (o storeOver) Split(string) (r0, r1 string) {
	o.spend()
	return r0, r1
}

// Find calls spend and returns its zero results.
func (o storeOver) Find(context.Context, ...string) (r0 []Value, r1 error) {
	o.spend()
	return r0, r1
}

// Level calls spend and returns its zero results.
func (o storeOver) Level(context.Context, string) (r0 Level, r1 error) {
	o.spend()
	return r0, r1
}

// Raise calls spend and returns its zero results.
func (o storeOver) Raise(context.Context, Level) (r0 error) {
	o.spend()
	return r0
}

// Inspect calls spend and returns its zero results.
func (o storeOver) Inspect(string) (r0 Value, r1 Meta, r2 bool) {
	o.spend()
	return r0, r1, r2
}

// GetWithMeta calls spend and returns its zero results.
func (o storeOver) GetWithMeta(context.Context, string) (r0 Value, r1 Meta, r2 error) {
	o.spend()
	return r0, r1, r2
}

// newStoreOver returns the planted defect of a ceiling of Store, whose
// callables call spend.
func newStoreOver(_ *StoreHandle, spend func()) Store {
	return storeOver{spend: spend}
}

// storePutPanicsDefect is the planted defect of a Store whose Put panics.
type storePutPanicsDefect struct {
	storeOverEmbedded
}

// Put is the override of the planted defect.
func (*storePutPanicsDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Store whose Put panics")
}

// storePutPanics returns the planted defect of a Store whose Put panics
// over a fresh subject.
func storePutPanics(h *StoreHandle) Store {
	return &storePutPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storePutIgnoresCancelDefect is the planted defect of a Store whose Put
// ignores the context.
type storePutIgnoresCancelDefect struct {
	storeOverEmbedded
}

// Put is the override of the planted defect.
func (d *storePutIgnoresCancelDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.storeOverEmbedded.Put(p0, p1)
}

// storePutIgnoresCancel returns the planted defect of a Store whose Put
// ignores the context over a fresh subject.
func storePutIgnoresCancel(h *StoreHandle) Store {
	return &storePutIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetPanicsDefect is the planted defect of a Store whose Get panics.
type storeGetPanicsDefect struct {
	storeOverEmbedded
}

// Get is the override of the planted defect.
func (*storeGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Store whose Get panics")
}

// storeGetPanics returns the planted defect of a Store whose Get panics
// over a fresh subject.
func storeGetPanics(h *StoreHandle) Store {
	return &storeGetPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetIgnoresCancelDefect is the planted defect of a Store whose Get
// ignores the context.
type storeGetIgnoresCancelDefect struct {
	storeOverEmbedded
}

// Get is the override of the planted defect.
func (d *storeGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.storeOverEmbedded.Get(p0, p1)
}

// storeGetIgnoresCancel returns the planted defect of a Store whose Get
// ignores the context over a fresh subject.
func storeGetIgnoresCancel(h *StoreHandle) Store {
	return &storeGetIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetValueBesideDefect is the planted defect of a Store whose Get
// returns a value beside an error.
type storeGetValueBesideDefect struct {
	storeOverEmbedded
}

// Get is the override of the planted defect.
func (d *storeGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.storeOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// storeGetValueBeside returns the planted defect of a Store whose Get
// returns a value beside an error over a fresh subject.
func storeGetValueBeside(h *StoreHandle) Store {
	return &storeGetValueBesideDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetReadsAnotherValueDefect is the planted defect of a Store whose
// Get returns a value that Put did not write.
type storeGetReadsAnotherValueDefect struct {
	storeOverEmbedded
}

// Get is the override of the planted defect.
func (d *storeGetReadsAnotherValueDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.storeOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// storeGetReadsAnotherValue returns the planted defect of a Store whose
// Get returns a value that Put did not write over a fresh subject.
func storeGetReadsAnotherValue(h *StoreHandle) Store {
	return &storeGetReadsAnotherValueDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetVersionDescendsDefect is the planted defect of a Store whose Get
// returns an older version on the second read.
type storeGetVersionDescendsDefect struct {
	storeOverEmbedded
	calls int
}

// Get is the override of the planted defect.
func (d *storeGetVersionDescendsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0 = Value{Key: "Key-b"}
	if d.calls == 1 {
		r0.Rev = 1
	} else {
		r0.Rev = 0
	}
	return r0, r1
}

// storeGetVersionDescends returns the planted defect of a Store whose Get
// returns an older version on the second read over a fresh subject.
func storeGetVersionDescends(h *StoreHandle) Store {
	return &storeGetVersionDescendsDefect{storeOverEmbedded: h.Fresh()}
}

// storeLoadPanicsDefect is the planted defect of a Store whose Load
// panics.
type storeLoadPanicsDefect struct {
	storeOverEmbedded
}

// Load is the override of the planted defect.
func (*storeLoadPanicsDefect) Load(p0 context.Context, p1 string) (r0 Value, r1 bool) {
	panic("the planted defect of a Store whose Load panics")
}

// storeLoadPanics returns the planted defect of a Store whose Load panics
// over a fresh subject.
func storeLoadPanics(h *StoreHandle) Store {
	return &storeLoadPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeLoadValueBesideAFalseFlagDefect is the planted defect of a Store
// whose Load returns a value beside a false presence flag.
type storeLoadValueBesideAFalseFlagDefect struct {
	storeOverEmbedded
}

// Load is the override of the planted defect.
func (d *storeLoadValueBesideAFalseFlagDefect) Load(
	p0 context.Context,
	p1 string,
) (r0 Value, r1 bool) {
	r0, r1 = d.storeOverEmbedded.Load(p0, p1)
	if !r1 {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// storeLoadValueBesideAFalseFlag returns the planted defect of a Store
// whose Load returns a value beside a false presence flag over a fresh
// subject.
func storeLoadValueBesideAFalseFlag(h *StoreHandle) Store {
	return &storeLoadValueBesideAFalseFlagDefect{storeOverEmbedded: h.Fresh()}
}

// storeLoadDriftsDefect is the planted defect of a Store whose Load
// returns another value on the repeat.
type storeLoadDriftsDefect struct {
	storeOverEmbedded
	calls int
}

// Load is the override of the planted defect.
func (d *storeLoadDriftsDefect) Load(p0 context.Context, p1 string) (r0 Value, r1 bool) {
	d.calls++
	if d.calls > 1 {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// storeLoadDrifts returns the planted defect of a Store whose Load returns
// another value on the repeat over a fresh subject.
func storeLoadDrifts(h *StoreHandle) Store {
	return &storeLoadDriftsDefect{storeOverEmbedded: h.Fresh()}
}

// storeSplitPanicsDefect is the planted defect of a Store whose Split
// panics.
type storeSplitPanicsDefect struct {
	storeOverEmbedded
}

// Split is the override of the planted defect.
func (*storeSplitPanicsDefect) Split(p0 string) (r0, r1 string) {
	panic("the planted defect of a Store whose Split panics")
}

// storeSplitPanics returns the planted defect of a Store whose Split
// panics over a fresh subject.
func storeSplitPanics(h *StoreHandle) Store {
	return &storeSplitPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeFindPanicsDefect is the planted defect of a Store whose Find
// panics.
type storeFindPanicsDefect struct {
	storeOverEmbedded
}

// Find is the override of the planted defect.
func (*storeFindPanicsDefect) Find(p0 context.Context, p1 ...string) (r0 []Value, r1 error) {
	panic("the planted defect of a Store whose Find panics")
}

// storeFindPanics returns the planted defect of a Store whose Find panics
// over a fresh subject.
func storeFindPanics(h *StoreHandle) Store {
	return &storeFindPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeFindIgnoresCancelDefect is the planted defect of a Store whose Find
// ignores the context.
type storeFindIgnoresCancelDefect struct {
	storeOverEmbedded
}

// Find is the override of the planted defect.
func (d *storeFindIgnoresCancelDefect) Find(
	p0 context.Context,
	p1 ...string,
) (r0 []Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.storeOverEmbedded.Find(p0, p1...)
}

// storeFindIgnoresCancel returns the planted defect of a Store whose Find
// ignores the context over a fresh subject.
func storeFindIgnoresCancel(h *StoreHandle) Store {
	return &storeFindIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}

// storeLevelPanicsDefect is the planted defect of a Store whose Level
// panics.
type storeLevelPanicsDefect struct {
	storeOverEmbedded
}

// Level is the override of the planted defect.
func (*storeLevelPanicsDefect) Level(p0 context.Context, p1 string) (r0 Level, r1 error) {
	panic("the planted defect of a Store whose Level panics")
}

// storeLevelPanics returns the planted defect of a Store whose Level
// panics over a fresh subject.
func storeLevelPanics(h *StoreHandle) Store {
	return &storeLevelPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeLevelIgnoresCancelDefect is the planted defect of a Store whose
// Level ignores the context.
type storeLevelIgnoresCancelDefect struct {
	storeOverEmbedded
}

// Level is the override of the planted defect.
func (d *storeLevelIgnoresCancelDefect) Level(p0 context.Context, p1 string) (r0 Level, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.storeOverEmbedded.Level(p0, p1)
}

// storeLevelIgnoresCancel returns the planted defect of a Store whose
// Level ignores the context over a fresh subject.
func storeLevelIgnoresCancel(h *StoreHandle) Store {
	return &storeLevelIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}

// storeLevelValueBesideDefect is the planted defect of a Store whose Level
// returns a value beside an error.
type storeLevelValueBesideDefect struct {
	storeOverEmbedded
}

// Level is the override of the planted defect.
func (d *storeLevelValueBesideDefect) Level(p0 context.Context, p1 string) (r0 Level, r1 error) {
	r0, r1 = d.storeOverEmbedded.Level(p0, p1)
	if r1 != nil {
		r0 = Level(7)
	}
	return r0, r1
}

// storeLevelValueBeside returns the planted defect of a Store whose Level
// returns a value beside an error over a fresh subject.
func storeLevelValueBeside(h *StoreHandle) Store {
	return &storeLevelValueBesideDefect{storeOverEmbedded: h.Fresh()}
}

// storeRaisePanicsDefect is the planted defect of a Store whose Raise
// panics.
type storeRaisePanicsDefect struct {
	storeOverEmbedded
}

// Raise is the override of the planted defect.
func (*storeRaisePanicsDefect) Raise(p0 context.Context, p1 Level) (r0 error) {
	panic("the planted defect of a Store whose Raise panics")
}

// storeRaisePanics returns the planted defect of a Store whose Raise
// panics over a fresh subject.
func storeRaisePanics(h *StoreHandle) Store {
	return &storeRaisePanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeRaiseIgnoresCancelDefect is the planted defect of a Store whose
// Raise ignores the context.
type storeRaiseIgnoresCancelDefect struct {
	storeOverEmbedded
}

// Raise is the override of the planted defect.
func (d *storeRaiseIgnoresCancelDefect) Raise(p0 context.Context, p1 Level) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.storeOverEmbedded.Raise(p0, p1)
}

// storeRaiseIgnoresCancel returns the planted defect of a Store whose
// Raise ignores the context over a fresh subject.
func storeRaiseIgnoresCancel(h *StoreHandle) Store {
	return &storeRaiseIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}

// storeInspectPanicsDefect is the planted defect of a Store whose Inspect
// panics.
type storeInspectPanicsDefect struct {
	storeOverEmbedded
}

// Inspect is the override of the planted defect.
func (*storeInspectPanicsDefect) Inspect(p0 string) (r0 Value, r1 Meta, r2 bool) {
	panic("the planted defect of a Store whose Inspect panics")
}

// storeInspectPanics returns the planted defect of a Store whose Inspect
// panics over a fresh subject.
func storeInspectPanics(h *StoreHandle) Store {
	return &storeInspectPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetWithMetaPanicsDefect is the planted defect of a Store whose
// GetWithMeta panics.
type storeGetWithMetaPanicsDefect struct {
	storeOverEmbedded
}

// GetWithMeta is the override of the planted defect.
func (*storeGetWithMetaPanicsDefect) GetWithMeta(
	p0 context.Context,
	p1 string,
) (r0 Value, r1 Meta, r2 error) {
	panic("the planted defect of a Store whose GetWithMeta panics")
}

// storeGetWithMetaPanics returns the planted defect of a Store whose
// GetWithMeta panics over a fresh subject.
func storeGetWithMetaPanics(h *StoreHandle) Store {
	return &storeGetWithMetaPanicsDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetWithMetaIgnoresCancelDefect is the planted defect of a Store
// whose GetWithMeta ignores the context.
type storeGetWithMetaIgnoresCancelDefect struct {
	storeOverEmbedded
}

// GetWithMeta is the override of the planted defect.
func (d *storeGetWithMetaIgnoresCancelDefect) GetWithMeta(
	p0 context.Context,
	p1 string,
) (r0 Value, r1 Meta, r2 error) {
	if p0.Err() != nil {
		return r0, r1, r2
	}
	return d.storeOverEmbedded.GetWithMeta(p0, p1)
}

// storeGetWithMetaIgnoresCancel returns the planted defect of a Store
// whose GetWithMeta ignores the context over a fresh subject.
func storeGetWithMetaIgnoresCancel(h *StoreHandle) Store {
	return &storeGetWithMetaIgnoresCancelDefect{storeOverEmbedded: h.Fresh()}
}

// storeGetWithMetaValuesBesideDefect is the planted defect of a Store
// whose GetWithMeta returns values beside an error.
type storeGetWithMetaValuesBesideDefect struct {
	storeOverEmbedded
}

// GetWithMeta is the override of the planted defect.
func (d *storeGetWithMetaValuesBesideDefect) GetWithMeta(
	p0 context.Context,
	p1 string,
) (r0 Value, r1 Meta, r2 error) {
	r0, r1, r2 = d.storeOverEmbedded.GetWithMeta(p0, p1)
	if r2 != nil {
		r0 = Value{Key: "Key-b"}
		r1 = Meta{Revision: 7}
	}
	return r0, r1, r2
}

// storeGetWithMetaValuesBeside returns the planted defect of a Store whose
// GetWithMeta returns values beside an error over a fresh subject.
func storeGetWithMetaValuesBeside(h *StoreHandle) Store {
	return &storeGetWithMetaValuesBesideDefect{storeOverEmbedded: h.Fresh()}
}

// The derived checks of Keyed:
//
// Keyed.Put/smoke: Put returns without a panic for derived inputs.
//
// Keyed.Put/cancel: Put returns an error for a cancelled context.
//
// Keyed.Get/smoke: Get returns without a panic for derived inputs.
//
// Keyed.Get/cancel: Get returns an error for a cancelled context.
//
// Keyed.Get/zero-on-error: Get returns the zero value beside an error.
//
// Keyed.Get/read-after-write: Get returns what Put stored before it.
//
// Keyed.Delete/smoke: Delete returns without a panic for derived inputs.
//
// Keyed.Delete/cancel: Delete returns an error for a cancelled context.

// The derived checks of Keyed that witness their laws at the shallow end:
//
// Keyed.Put/smoke
//
// Keyed.Get/smoke
//
// Keyed.Get/read-after-write
//
// Keyed.Delete/smoke

// The refused claims of Keyed:
//
// Keyed.Delete/delete-removes: no callable of the declaration is the
// partner put of Delete. The remedy is to name the partner of Delete with
// the param of the classification.

// The unclaimed classifications of Keyed:
//
// The shape:composite-writer of Keyed.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.
//
// The shape:deleter of Keyed.Delete has the stance silent. A delete is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Keyed writes.
type (
	KeyedHandle = conform.Handle[Keyed]
	KeyedCheck  = conform.Check[Keyed]
	KeyedOption = conform.Option[Keyed]
)

// KeyedChecks contains the name of each callable of Keyed, with the
// identities of its checks.
var KeyedChecks = struct {
	Put    keyedPut
	Get    keyedGet
	Delete keyedDelete
}{
	Put:    keyedPut{"Keyed.Put"},
	Get:    keyedGet{"Keyed.Get"},
	Delete: keyedDelete{"Keyed.Delete"},
}

// keyedPut is the name of Keyed.Put, with the identities of its checks.
type keyedPut struct {
	conform.Callable[Keyed]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c keyedPut) Smoke() conform.ID[Keyed] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c keyedPut) Cancel() conform.ID[Keyed] {
	return c.Check("cancel")
}

// keyedGet is the name of Keyed.Get, with the identities of its checks.
type keyedGet struct {
	conform.Callable[Keyed]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c keyedGet) Smoke() conform.ID[Keyed] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c keyedGet) Cancel() conform.ID[Keyed] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c keyedGet) ZeroOnError() conform.ID[Keyed] {
	return c.Check("zero-on-error")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c keyedGet) ReadAfterWrite() conform.ID[Keyed] {
	return c.Check("read-after-write")
}

// keyedDelete is the name of Keyed.Delete, with the identities of its
// checks.
type keyedDelete struct {
	conform.Callable[Keyed]
}

// Smoke returns the identity of the check that Delete returns without a
// panic for derived inputs.
func (c keyedDelete) Smoke() conform.ID[Keyed] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Delete returns an error
// for a cancelled context.
func (c keyedDelete) Cancel() conform.ID[Keyed] {
	return c.Check("cancel")
}

// KeyedSuite returns the checks of Keyed, and the calls and the budgets of
// its callables.
func KeyedSuite() conform.Suite[Keyed] {
	return conform.Suite[Keyed]{
		Declaration: "Keyed",
		Checks: []conform.Check[Keyed]{
			{
				ID:          KeyedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeyedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: keyedPutPanics,
			},
			{
				ID:          KeyedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: keyedPutIgnoresCancel,
			},
			{
				ID:          KeyedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeyedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: keyedGetPanics,
			},
			{
				ID:          KeyedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: keyedGetIgnoresCancel,
			},
			{
				ID:          KeyedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: keyedGetValueBeside,
			},
			{
				ID:          KeyedChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeyedHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, func(ctx context.Context) (Value, error) {
						w := Value{Key: "Key-a"}
						err := h.Subject().Put(ctx, string("key-a"), w)
						return w, err
					})
				},
				Defect: keyedGetReadsAnotherValue,
			},
			{
				ID:          KeyedChecks.Delete.Smoke(),
				Statement:   "Delete returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeyedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Delete(ctx, string("key-a"))
					})
				},
				Defect: keyedDeletePanics,
			},
			{
				ID:          KeyedChecks.Delete.Cancel(),
				Statement:   "Delete returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Delete(ctx, string("key-a"))
					})
				},
				Defect: keyedDeleteIgnoresCancel,
			},
		},
		Calls: []conform.Call[Keyed]{
			{
				Calls:    keyedPutCalls,
				Over:     newKeyedOver,
				Callable: KeyedChecks.Put.Callable,
			},
			{
				Calls:    keyedGetCalls,
				Over:     newKeyedOver,
				Callable: KeyedChecks.Get.Callable,
			},
			{
				Calls:    keyedDeleteCalls,
				Over:     newKeyedOver,
				Callable: KeyedChecks.Delete.Callable,
			},
		},
	}
}

// RunKeyedSuite applies the checks of Keyed, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunKeyedSuite[T Keyed](t *testing.T, newSubject func() T, opts ...KeyedOption) {
	t.Helper()
	RunKeyedSuiteFromHandle(t, func(*KeyedHandle) T { return newSubject() }, opts...)
}

// RunKeyedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunKeyedSuiteFromHandle[T Keyed](
	t *testing.T,
	newSubject func(h *KeyedHandle) T,
	opts ...KeyedOption,
) {
	t.Helper()
	conform.Run(t, KeyedSuite(), keyedSubjects(newSubject), opts...)
}

// FuzzKeyedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzKeyedSuite[T Keyed](
	f *testing.F,
	callable conform.Callable[Keyed],
	newSubject func(h *KeyedHandle) T,
	opts ...KeyedOption,
) {
	f.Helper()
	conform.Fuzz(f, KeyedSuite(), callable, keyedSubjects(newSubject), opts...)
}

// RunKeyedBudgets applies the allocation ceilings of Keyed and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunKeyedBudgets[T Keyed](
	t *testing.T,
	newSubject func(h *KeyedHandle) T,
	opts ...KeyedOption,
) {
	t.Helper()
	conform.RunBudgets(t, KeyedSuite(), keyedSubjects(newSubject), opts...)
}

// BenchKeyedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchKeyedSuite[T Keyed](
	b conform.B,
	callable conform.Callable[Keyed],
	newSubject func(h *KeyedHandle) T,
	opts ...KeyedOption,
) {
	b.Helper()
	conform.Bench(b, KeyedSuite(), callable, keyedSubjects(newSubject), opts...)
}

// keyedSubjects returns the factory of the subjects of Keyed that calls
// newSubject with the handle of each check.
func keyedSubjects[T Keyed](newSubject func(h *KeyedHandle) T) func(h *KeyedHandle) Keyed {
	return func(h *KeyedHandle) Keyed { return newSubject(h) }
}

// keyedPutCalls generates the calls of Keyed.Put with the inputs that each
// case draws.
var keyedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Keyed) error {
	key := c.Draw(prop.Of[string](), "key")
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Keyed) error {
		return subject.Put(ctx, key, v)
	}
})

// keyedGetCalls generates the calls of Keyed.Get with the inputs that each
// case draws.
var keyedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Keyed) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Keyed) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// keyedDeleteCalls generates the calls of Keyed.Delete with the inputs
// that each case draws.
var keyedDeleteCalls = prop.Composite(func(c *prop.Case) func(context.Context, Keyed) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Keyed) error {
		return subject.Delete(ctx, key)
	}
})

// keyedOverEmbedded is Keyed under the name of the field that keyedOver
// embeds.
type keyedOverEmbedded = Keyed

// keyedOver is a subject of Keyed whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type keyedOver struct {
	keyedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o keyedOver) Put(context.Context, string, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o keyedOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Delete calls spend and returns its zero results.
func (o keyedOver) Delete(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// newKeyedOver returns the planted defect of a ceiling of Keyed, whose
// callables call spend.
func newKeyedOver(_ *KeyedHandle, spend func()) Keyed {
	return keyedOver{spend: spend}
}

// keyedPutPanicsDefect is the planted defect of a Keyed whose Put panics.
type keyedPutPanicsDefect struct {
	keyedOverEmbedded
}

// Put is the override of the planted defect.
func (*keyedPutPanicsDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	panic("the planted defect of a Keyed whose Put panics")
}

// keyedPutPanics returns the planted defect of a Keyed whose Put panics
// over a fresh subject.
func keyedPutPanics(h *KeyedHandle) Keyed {
	return &keyedPutPanicsDefect{keyedOverEmbedded: h.Fresh()}
}

// keyedPutIgnoresCancelDefect is the planted defect of a Keyed whose Put
// ignores the context.
type keyedPutIgnoresCancelDefect struct {
	keyedOverEmbedded
}

// Put is the override of the planted defect.
func (d *keyedPutIgnoresCancelDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.keyedOverEmbedded.Put(p0, p1, p2)
}

// keyedPutIgnoresCancel returns the planted defect of a Keyed whose Put
// ignores the context over a fresh subject.
func keyedPutIgnoresCancel(h *KeyedHandle) Keyed {
	return &keyedPutIgnoresCancelDefect{keyedOverEmbedded: h.Fresh()}
}

// keyedGetPanicsDefect is the planted defect of a Keyed whose Get panics.
type keyedGetPanicsDefect struct {
	keyedOverEmbedded
}

// Get is the override of the planted defect.
func (*keyedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Keyed whose Get panics")
}

// keyedGetPanics returns the planted defect of a Keyed whose Get panics
// over a fresh subject.
func keyedGetPanics(h *KeyedHandle) Keyed {
	return &keyedGetPanicsDefect{keyedOverEmbedded: h.Fresh()}
}

// keyedGetIgnoresCancelDefect is the planted defect of a Keyed whose Get
// ignores the context.
type keyedGetIgnoresCancelDefect struct {
	keyedOverEmbedded
}

// Get is the override of the planted defect.
func (d *keyedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.keyedOverEmbedded.Get(p0, p1)
}

// keyedGetIgnoresCancel returns the planted defect of a Keyed whose Get
// ignores the context over a fresh subject.
func keyedGetIgnoresCancel(h *KeyedHandle) Keyed {
	return &keyedGetIgnoresCancelDefect{keyedOverEmbedded: h.Fresh()}
}

// keyedGetValueBesideDefect is the planted defect of a Keyed whose Get
// returns a value beside an error.
type keyedGetValueBesideDefect struct {
	keyedOverEmbedded
}

// Get is the override of the planted defect.
func (d *keyedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.keyedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// keyedGetValueBeside returns the planted defect of a Keyed whose Get
// returns a value beside an error over a fresh subject.
func keyedGetValueBeside(h *KeyedHandle) Keyed {
	return &keyedGetValueBesideDefect{keyedOverEmbedded: h.Fresh()}
}

// keyedGetReadsAnotherValueDefect is the planted defect of a Keyed whose
// Get returns a value that Put did not write.
type keyedGetReadsAnotherValueDefect struct {
	keyedOverEmbedded
}

// Get is the override of the planted defect.
func (d *keyedGetReadsAnotherValueDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.keyedOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// keyedGetReadsAnotherValue returns the planted defect of a Keyed whose
// Get returns a value that Put did not write over a fresh subject.
func keyedGetReadsAnotherValue(h *KeyedHandle) Keyed {
	return &keyedGetReadsAnotherValueDefect{keyedOverEmbedded: h.Fresh()}
}

// keyedDeletePanicsDefect is the planted defect of a Keyed whose Delete
// panics.
type keyedDeletePanicsDefect struct {
	keyedOverEmbedded
}

// Delete is the override of the planted defect.
func (*keyedDeletePanicsDefect) Delete(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Keyed whose Delete panics")
}

// keyedDeletePanics returns the planted defect of a Keyed whose Delete
// panics over a fresh subject.
func keyedDeletePanics(h *KeyedHandle) Keyed {
	return &keyedDeletePanicsDefect{keyedOverEmbedded: h.Fresh()}
}

// keyedDeleteIgnoresCancelDefect is the planted defect of a Keyed whose
// Delete ignores the context.
type keyedDeleteIgnoresCancelDefect struct {
	keyedOverEmbedded
}

// Delete is the override of the planted defect.
func (d *keyedDeleteIgnoresCancelDefect) Delete(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.keyedOverEmbedded.Delete(p0, p1)
}

// keyedDeleteIgnoresCancel returns the planted defect of a Keyed whose
// Delete ignores the context over a fresh subject.
func keyedDeleteIgnoresCancel(h *KeyedHandle) Keyed {
	return &keyedDeleteIgnoresCancelDefect{keyedOverEmbedded: h.Fresh()}
}

// The derived checks of Stamped:
//
// Stamped.Store/smoke: Store returns without a panic for derived inputs.
//
// Stamped.Store/cancel: Store returns an error for a cancelled context.
//
// Stamped.Store/zero-on-error: Store returns the zero value beside an
// error.
//
// Stamped.Get/smoke: Get returns without a panic for derived inputs.
//
// Stamped.Get/cancel: Get returns an error for a cancelled context.
//
// Stamped.Get/zero-on-error: Get returns the zero value beside an error.
//
// Stamped.Get/read-after-write: Get returns what Store stored before it.
//
// Stamped.Get/writes-follow-reads: a Store after Get has a version above
// the version that Get returned.

// The derived checks of Stamped that witness their laws at the shallow
// end:
//
// Stamped.Store/smoke
//
// Stamped.Get/smoke
//
// Stamped.Get/read-after-write
//
// Stamped.Get/writes-follow-reads

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Stamped writes.
type (
	StampedHandle = conform.Handle[Stamped]
	StampedCheck  = conform.Check[Stamped]
	StampedOption = conform.Option[Stamped]
)

// StampedChecks contains the name of each callable of Stamped, with the
// identities of its checks.
var StampedChecks = struct {
	Store stampedStore
	Get   stampedGet
}{
	Store: stampedStore{"Stamped.Store"},
	Get:   stampedGet{"Stamped.Get"},
}

// stampedStore is the name of Stamped.Store, with the identities of its
// checks.
type stampedStore struct {
	conform.Callable[Stamped]
}

// Smoke returns the identity of the check that Store returns without a
// panic for derived inputs.
func (c stampedStore) Smoke() conform.ID[Stamped] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Store returns an error for
// a cancelled context.
func (c stampedStore) Cancel() conform.ID[Stamped] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Store returns the
// zero value beside an error.
func (c stampedStore) ZeroOnError() conform.ID[Stamped] {
	return c.Check("zero-on-error")
}

// stampedGet is the name of Stamped.Get, with the identities of its
// checks.
type stampedGet struct {
	conform.Callable[Stamped]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c stampedGet) Smoke() conform.ID[Stamped] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c stampedGet) Cancel() conform.ID[Stamped] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c stampedGet) ZeroOnError() conform.ID[Stamped] {
	return c.Check("zero-on-error")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Store stored before it.
func (c stampedGet) ReadAfterWrite() conform.ID[Stamped] {
	return c.Check("read-after-write")
}

// WritesFollowReads returns the identity of the check that a Store after
// Get has a version above the version that Get returned.
func (c stampedGet) WritesFollowReads() conform.ID[Stamped] {
	return c.Check("writes-follow-reads")
}

// StampedSuite returns the checks of Stamped, and the calls and the
// budgets of its callables.
func StampedSuite() conform.Suite[Stamped] {
	return conform.Suite[Stamped]{
		Declaration: "Stamped",
		Checks: []conform.Check[Stamped]{
			{
				ID:          StampedChecks.Store.Smoke(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Store(ctx, Value{Key: "Key-a"})
						return err
					})
				},
				Defect: stampedStorePanics,
			},
			{
				ID:          StampedChecks.Store.Cancel(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Store(ctx, Value{Key: "Key-a"})
						return err
					})
				},
				Defect: stampedStoreIgnoresCancel,
			},
			{
				ID:          StampedChecks.Store.ZeroOnError(),
				Statement:   "Store returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					}, Value{})
				},
				Defect: stampedStoreValueBeside,
			},
			{
				ID:          StampedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: stampedGetPanics,
			},
			{
				ID:          StampedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: stampedGetIgnoresCancel,
			},
			{
				ID:          StampedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: stampedGetValueBeside,
			},
			{
				ID:          StampedChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Store stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) (Value, error) {
						w := Value{Key: "Key-a"}
						return h.Subject().Store(ctx, w)
					})
				},
				Defect: stampedGetReadsAnotherValue,
			},
			{
				ID:          StampedChecks.Get.WritesFollowReads(),
				Statement:   "a Store after Get has a version above the version that Get returned",
				Claim:       spec.WritesFollowReads,
				Probe:       spec.WritesAdvance,
				Reason:      spec.UnorderedWrite,
				Law:         spec.WriteFollowsTheReadBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.WritesAdvance(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) (Value, error) {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					}, func(v Value) int64 { return v.Rev })
				},
				Defect: stampedStoreStampsThePast,
			},
		},
		Calls: []conform.Call[Stamped]{
			{
				Calls:    stampedStoreCalls,
				Over:     newStampedOver,
				Callable: StampedChecks.Store.Callable,
			},
			{
				Calls:    stampedGetCalls,
				Over:     newStampedOver,
				Callable: StampedChecks.Get.Callable,
			},
		},
	}
}

// RunStampedSuite applies the checks of Stamped, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunStampedSuite[T Stamped](t *testing.T, newSubject func() T, opts ...StampedOption) {
	t.Helper()
	RunStampedSuiteFromHandle(t, func(*StampedHandle) T { return newSubject() }, opts...)
}

// RunStampedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStampedSuiteFromHandle[T Stamped](
	t *testing.T,
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	t.Helper()
	conform.Run(t, StampedSuite(), stampedSubjects(newSubject), opts...)
}

// FuzzStampedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzStampedSuite[T Stamped](
	f *testing.F,
	callable conform.Callable[Stamped],
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	f.Helper()
	conform.Fuzz(f, StampedSuite(), callable, stampedSubjects(newSubject), opts...)
}

// RunStampedBudgets applies the allocation ceilings of Stamped and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunStampedBudgets[T Stamped](
	t *testing.T,
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	t.Helper()
	conform.RunBudgets(t, StampedSuite(), stampedSubjects(newSubject), opts...)
}

// BenchStampedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchStampedSuite[T Stamped](
	b conform.B,
	callable conform.Callable[Stamped],
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	b.Helper()
	conform.Bench(b, StampedSuite(), callable, stampedSubjects(newSubject), opts...)
}

// stampedSubjects returns the factory of the subjects of Stamped that
// calls newSubject with the handle of each check.
func stampedSubjects[T Stamped](
	newSubject func(h *StampedHandle) T,
) func(h *StampedHandle) Stamped {
	return func(h *StampedHandle) Stamped { return newSubject(h) }
}

// stampedStoreCalls generates the calls of Stamped.Store with the inputs
// that each case draws.
var stampedStoreCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stamped) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Stamped) error {
		_, err := subject.Store(ctx, v)
		return err
	}
})

// stampedGetCalls generates the calls of Stamped.Get with the inputs that
// each case draws.
var stampedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stamped) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Stamped) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// stampedOverEmbedded is Stamped under the name of the field that
// stampedOver embeds.
type stampedOverEmbedded = Stamped

// stampedOver is a subject of Stamped whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type stampedOver struct {
	stampedOverEmbedded
	spend func()
}

// Store calls spend and returns its zero results.
func (o stampedOver) Store(context.Context, Value) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Get calls spend and returns its zero results.
func (o stampedOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newStampedOver returns the planted defect of a ceiling of Stamped, whose
// callables call spend.
func newStampedOver(_ *StampedHandle, spend func()) Stamped {
	return stampedOver{spend: spend}
}

// stampedStorePanicsDefect is the planted defect of a Stamped whose Store
// panics.
type stampedStorePanicsDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (*stampedStorePanicsDefect) Store(p0 context.Context, p1 Value) (r0 Value, r1 error) {
	panic("the planted defect of a Stamped whose Store panics")
}

// stampedStorePanics returns the planted defect of a Stamped whose Store
// panics over a fresh subject.
func stampedStorePanics(h *StampedHandle) Stamped {
	return &stampedStorePanicsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedStoreIgnoresCancelDefect is the planted defect of a Stamped whose
// Store ignores the context.
type stampedStoreIgnoresCancelDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (d *stampedStoreIgnoresCancelDefect) Store(p0 context.Context, p1 Value) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.stampedOverEmbedded.Store(p0, p1)
}

// stampedStoreIgnoresCancel returns the planted defect of a Stamped whose
// Store ignores the context over a fresh subject.
func stampedStoreIgnoresCancel(h *StampedHandle) Stamped {
	return &stampedStoreIgnoresCancelDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedStoreValueBesideDefect is the planted defect of a Stamped whose
// Store returns a value beside an error.
type stampedStoreValueBesideDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (d *stampedStoreValueBesideDefect) Store(p0 context.Context, p1 Value) (r0 Value, r1 error) {
	r0, r1 = d.stampedOverEmbedded.Store(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// stampedStoreValueBeside returns the planted defect of a Stamped whose
// Store returns a value beside an error over a fresh subject.
func stampedStoreValueBeside(h *StampedHandle) Stamped {
	return &stampedStoreValueBesideDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetPanicsDefect is the planted defect of a Stamped whose Get
// panics.
type stampedGetPanicsDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (*stampedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Stamped whose Get panics")
}

// stampedGetPanics returns the planted defect of a Stamped whose Get
// panics over a fresh subject.
func stampedGetPanics(h *StampedHandle) Stamped {
	return &stampedGetPanicsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetIgnoresCancelDefect is the planted defect of a Stamped whose
// Get ignores the context.
type stampedGetIgnoresCancelDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (d *stampedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.stampedOverEmbedded.Get(p0, p1)
}

// stampedGetIgnoresCancel returns the planted defect of a Stamped whose
// Get ignores the context over a fresh subject.
func stampedGetIgnoresCancel(h *StampedHandle) Stamped {
	return &stampedGetIgnoresCancelDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetValueBesideDefect is the planted defect of a Stamped whose Get
// returns a value beside an error.
type stampedGetValueBesideDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (d *stampedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.stampedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// stampedGetValueBeside returns the planted defect of a Stamped whose Get
// returns a value beside an error over a fresh subject.
func stampedGetValueBeside(h *StampedHandle) Stamped {
	return &stampedGetValueBesideDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetReadsAnotherValueDefect is the planted defect of a Stamped
// whose Get returns a value that Store did not write.
type stampedGetReadsAnotherValueDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (d *stampedGetReadsAnotherValueDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.stampedOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// stampedGetReadsAnotherValue returns the planted defect of a Stamped
// whose Get returns a value that Store did not write over a fresh subject.
func stampedGetReadsAnotherValue(h *StampedHandle) Stamped {
	return &stampedGetReadsAnotherValueDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedStoreStampsThePastDefect is the planted defect of a Stamped whose
// Store is not ordered after the read before it.
type stampedStoreStampsThePastDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (d *stampedStoreStampsThePastDefect) Store(p0 context.Context, p1 Value) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.stampedOverEmbedded.Store(p0, p1)
	r0.Rev = 0
	return r0, r1
}

// stampedStoreStampsThePast returns the planted defect of a Stamped whose
// Store is not ordered after the read before it over a fresh subject.
func stampedStoreStampsThePast(h *StampedHandle) Stamped {
	return &stampedStoreStampsThePastDefect{stampedOverEmbedded: h.Fresh()}
}

// The derived checks of Decoded:
//
// Decoded.Put/smoke: Put returns without a panic for derived inputs.
//
// Decoded.Put/cancel: Put returns an error for a cancelled context.
//
// Decoded.Get/smoke: Get returns without a panic for derived inputs.
//
// Decoded.Get/cancel: Get returns an error for a cancelled context.
//
// Decoded.Get/read-after-write: Get returns what Put stored before it.

// The derived checks of Decoded that witness their laws at the shallow
// end:
//
// Decoded.Put/smoke
//
// Decoded.Get/smoke
//
// Decoded.Get/read-after-write

// The unclaimed classifications of Decoded:
//
// The shape:composite-writer of Decoded.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Decoded writes.
type (
	DecodedHandle = conform.Handle[Decoded]
	DecodedCheck  = conform.Check[Decoded]
	DecodedOption = conform.Option[Decoded]
)

// DecodedChecks contains the name of each callable of Decoded, with the
// identities of its checks.
var DecodedChecks = struct {
	Put decodedPut
	Get decodedGet
}{
	Put: decodedPut{"Decoded.Put"},
	Get: decodedGet{"Decoded.Get"},
}

// decodedPut is the name of Decoded.Put, with the identities of its
// checks.
type decodedPut struct {
	conform.Callable[Decoded]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c decodedPut) Smoke() conform.ID[Decoded] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c decodedPut) Cancel() conform.ID[Decoded] {
	return c.Check("cancel")
}

// decodedGet is the name of Decoded.Get, with the identities of its
// checks.
type decodedGet struct {
	conform.Callable[Decoded]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c decodedGet) Smoke() conform.ID[Decoded] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c decodedGet) Cancel() conform.ID[Decoded] {
	return c.Check("cancel")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c decodedGet) ReadAfterWrite() conform.ID[Decoded] {
	return c.Check("read-after-write")
}

// DecodedSuite returns the checks of Decoded, and the calls and the
// budgets of its callables.
func DecodedSuite() conform.Suite[Decoded] {
	return conform.Suite[Decoded]{
		Declaration: "Decoded",
		Checks: []conform.Check[Decoded]{
			{
				ID:          DecodedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DecodedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: decodedPutPanics,
			},
			{
				ID:          DecodedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DecodedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: decodedPutIgnoresCancel,
			},
			{
				ID:          DecodedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DecodedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Get(ctx, string("key-a"), func(Value) Value { return Value{} })
						return err
					})
				},
				Defect: decodedGetPanics,
			},
			{
				ID:          DecodedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DecodedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Get(ctx, string("key-a"), func(Value) Value { return Value{} })
						return err
					})
				},
				Defect: decodedGetIgnoresCancel,
			},
			{
				ID:          DecodedChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DecodedHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Value, error) {
						return h.Subject().
							Get(ctx, string("key-a"), func(Value) Value { return Value{} })
					}, func(ctx context.Context) (Value, error) {
						w := Value{Key: "Key-a"}
						err := h.Subject().Put(ctx, string("key-a"), w)
						return w, err
					})
				},
				Defect: decodedGetReadsAnotherValue,
			},
		},
		Calls: []conform.Call[Decoded]{
			{
				Calls:    decodedPutCalls,
				Over:     newDecodedOver,
				Callable: DecodedChecks.Put.Callable,
			},
			{
				Calls:    decodedGetCalls,
				Over:     newDecodedOver,
				Callable: DecodedChecks.Get.Callable,
			},
		},
	}
}

// RunDecodedSuite applies the checks of Decoded, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunDecodedSuite[T Decoded](t *testing.T, newSubject func() T, opts ...DecodedOption) {
	t.Helper()
	RunDecodedSuiteFromHandle(t, func(*DecodedHandle) T { return newSubject() }, opts...)
}

// RunDecodedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDecodedSuiteFromHandle[T Decoded](
	t *testing.T,
	newSubject func(h *DecodedHandle) T,
	opts ...DecodedOption,
) {
	t.Helper()
	conform.Run(t, DecodedSuite(), decodedSubjects(newSubject), opts...)
}

// FuzzDecodedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzDecodedSuite[T Decoded](
	f *testing.F,
	callable conform.Callable[Decoded],
	newSubject func(h *DecodedHandle) T,
	opts ...DecodedOption,
) {
	f.Helper()
	conform.Fuzz(f, DecodedSuite(), callable, decodedSubjects(newSubject), opts...)
}

// RunDecodedBudgets applies the allocation ceilings of Decoded and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunDecodedBudgets[T Decoded](
	t *testing.T,
	newSubject func(h *DecodedHandle) T,
	opts ...DecodedOption,
) {
	t.Helper()
	conform.RunBudgets(t, DecodedSuite(), decodedSubjects(newSubject), opts...)
}

// BenchDecodedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchDecodedSuite[T Decoded](
	b conform.B,
	callable conform.Callable[Decoded],
	newSubject func(h *DecodedHandle) T,
	opts ...DecodedOption,
) {
	b.Helper()
	conform.Bench(b, DecodedSuite(), callable, decodedSubjects(newSubject), opts...)
}

// decodedSubjects returns the factory of the subjects of Decoded that
// calls newSubject with the handle of each check.
func decodedSubjects[T Decoded](
	newSubject func(h *DecodedHandle) T,
) func(h *DecodedHandle) Decoded {
	return func(h *DecodedHandle) Decoded { return newSubject(h) }
}

// decodedPutCalls generates the calls of Decoded.Put with the inputs that
// each case draws.
var decodedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Decoded) error {
	key := c.Draw(prop.Of[string](), "key")
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Decoded) error {
		return subject.Put(ctx, key, v)
	}
})

// decodedGetCalls generates the calls of Decoded.Get with the inputs that
// each case draws.
var decodedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Decoded) error {
	key := c.Draw(prop.Of[string](), "key")
	decodeResults := c.Draw(prop.List(prop.Of[Value]()), "decode")
	var decodeAt atomic.Int64
	decode := func(Value) (r0 Value) {
		i := int(decodeAt.Add(1) - 1)
		if i >= len(decodeResults) {
			return r0
		}
		return decodeResults[i]
	}
	return func(ctx context.Context, subject Decoded) error {
		_, err := subject.Get(ctx, key, decode)
		return err
	}
})

// decodedOverEmbedded is Decoded under the name of the field that
// decodedOver embeds.
type decodedOverEmbedded = Decoded

// decodedOver is a subject of Decoded whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type decodedOver struct {
	decodedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o decodedOver) Put(context.Context, string, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o decodedOver) Get(context.Context, string, func(Value) Value) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newDecodedOver returns the planted defect of a ceiling of Decoded, whose
// callables call spend.
func newDecodedOver(_ *DecodedHandle, spend func()) Decoded {
	return decodedOver{spend: spend}
}

// decodedPutPanicsDefect is the planted defect of a Decoded whose Put
// panics.
type decodedPutPanicsDefect struct {
	decodedOverEmbedded
}

// Put is the override of the planted defect.
func (*decodedPutPanicsDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	panic("the planted defect of a Decoded whose Put panics")
}

// decodedPutPanics returns the planted defect of a Decoded whose Put
// panics over a fresh subject.
func decodedPutPanics(h *DecodedHandle) Decoded {
	return &decodedPutPanicsDefect{decodedOverEmbedded: h.Fresh()}
}

// decodedPutIgnoresCancelDefect is the planted defect of a Decoded whose
// Put ignores the context.
type decodedPutIgnoresCancelDefect struct {
	decodedOverEmbedded
}

// Put is the override of the planted defect.
func (d *decodedPutIgnoresCancelDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.decodedOverEmbedded.Put(p0, p1, p2)
}

// decodedPutIgnoresCancel returns the planted defect of a Decoded whose
// Put ignores the context over a fresh subject.
func decodedPutIgnoresCancel(h *DecodedHandle) Decoded {
	return &decodedPutIgnoresCancelDefect{decodedOverEmbedded: h.Fresh()}
}

// decodedGetPanicsDefect is the planted defect of a Decoded whose Get
// panics.
type decodedGetPanicsDefect struct {
	decodedOverEmbedded
}

// Get is the override of the planted defect.
func (*decodedGetPanicsDefect) Get(
	p0 context.Context,
	p1 string,
	p2 func(Value) Value,
) (r0 Value, r1 error) {
	panic("the planted defect of a Decoded whose Get panics")
}

// decodedGetPanics returns the planted defect of a Decoded whose Get
// panics over a fresh subject.
func decodedGetPanics(h *DecodedHandle) Decoded {
	return &decodedGetPanicsDefect{decodedOverEmbedded: h.Fresh()}
}

// decodedGetIgnoresCancelDefect is the planted defect of a Decoded whose
// Get ignores the context.
type decodedGetIgnoresCancelDefect struct {
	decodedOverEmbedded
}

// Get is the override of the planted defect.
func (d *decodedGetIgnoresCancelDefect) Get(
	p0 context.Context,
	p1 string,
	p2 func(Value) Value,
) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.decodedOverEmbedded.Get(p0, p1, p2)
}

// decodedGetIgnoresCancel returns the planted defect of a Decoded whose
// Get ignores the context over a fresh subject.
func decodedGetIgnoresCancel(h *DecodedHandle) Decoded {
	return &decodedGetIgnoresCancelDefect{decodedOverEmbedded: h.Fresh()}
}

// decodedGetReadsAnotherValueDefect is the planted defect of a Decoded
// whose Get returns a value that Put did not write.
type decodedGetReadsAnotherValueDefect struct {
	decodedOverEmbedded
}

// Get is the override of the planted defect.
func (d *decodedGetReadsAnotherValueDefect) Get(
	p0 context.Context,
	p1 string,
	p2 func(Value) Value,
) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.decodedOverEmbedded.Get(p0, p1, p2)
	if r1 == nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// decodedGetReadsAnotherValue returns the planted defect of a Decoded
// whose Get returns a value that Put did not write over a fresh subject.
func decodedGetReadsAnotherValue(h *DecodedHandle) Decoded {
	return &decodedGetReadsAnotherValueDefect{decodedOverEmbedded: h.Fresh()}
}

// The derived checks of Expiring:
//
// Expiring.Put/smoke: Put returns without a panic for derived inputs.
//
// Expiring.Put/cancel: Put returns an error for a cancelled context.
//
// Expiring.Read/smoke: Read returns without a panic for derived inputs.
//
// Expiring.Read/cancel: Read returns an error for a cancelled context.
//
// Expiring.Read/zero-on-error: Read returns the zero value beside an
// error.
//
// Expiring.Read/expires: Read fails with the declared sentinel once the
// clock passes the declared lifetime.

// The derived checks of Expiring that witness their laws at the shallow
// end:
//
// Expiring.Put/smoke
//
// Expiring.Read/smoke

// The unclaimed classifications of Expiring:
//
// The shape:writer of Expiring.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Expiring writes.
type (
	ExpiringHandle = conform.Handle[Expiring]
	ExpiringCheck  = conform.Check[Expiring]
	ExpiringOption = conform.Option[Expiring]
)

// ExpiringChecks contains the name of each callable of Expiring, with the
// identities of its checks.
var ExpiringChecks = struct {
	Put  expiringPut
	Read expiringRead
}{
	Put:  expiringPut{"Expiring.Put"},
	Read: expiringRead{"Expiring.Read"},
}

// expiringPut is the name of Expiring.Put, with the identities of its
// checks.
type expiringPut struct {
	conform.Callable[Expiring]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c expiringPut) Smoke() conform.ID[Expiring] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c expiringPut) Cancel() conform.ID[Expiring] {
	return c.Check("cancel")
}

// expiringRead is the name of Expiring.Read, with the identities of its
// checks.
type expiringRead struct {
	conform.Callable[Expiring]
}

// Smoke returns the identity of the check that Read returns without a
// panic for derived inputs.
func (c expiringRead) Smoke() conform.ID[Expiring] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Read returns an error for
// a cancelled context.
func (c expiringRead) Cancel() conform.ID[Expiring] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Read returns the zero
// value beside an error.
func (c expiringRead) ZeroOnError() conform.ID[Expiring] {
	return c.Check("zero-on-error")
}

// Expires returns the identity of the check that Read fails with the
// declared sentinel once the clock passes the declared lifetime.
func (c expiringRead) Expires() conform.ID[Expiring] {
	return c.Check("expires")
}

// ExpiringSuite returns the checks of Expiring, and the calls and the
// budgets of its callables.
func ExpiringSuite() conform.Suite[Expiring] {
	return conform.Suite[Expiring]{
		Declaration: "Expiring",
		Checks: []conform.Check[Expiring]{
			{
				ID:          ExpiringChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ExpiringHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: expiringPutPanics,
			},
			{
				ID:          ExpiringChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ExpiringHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: expiringPutIgnoresCancel,
			},
			{
				ID:          ExpiringChecks.Read.Smoke(),
				Statement:   "Read returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ExpiringHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, string("key-a"))
						return err
					})
				},
				Defect: expiringReadPanics,
			},
			{
				ID:          ExpiringChecks.Read.Cancel(),
				Statement:   "Read returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ExpiringHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, string("key-a"))
						return err
					})
				},
				Defect: expiringReadIgnoresCancel,
			},
			{
				ID:          ExpiringChecks.Read.ZeroOnError(),
				Statement:   "Read returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ExpiringHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Read(ctx, string("key-a"))
					}, Value{})
				},
				Defect: expiringReadValueBeside,
			},
			{
				ID:          ExpiringChecks.Read.Expires(),
				Statement:   "Read fails with the declared sentinel once the clock passes the declared lifetime",
				Claim:       spec.Expires,
				Probe:       spec.ExpiresAfter,
				Reason:      spec.NotExpired,
				Law:         spec.EntryLapsesOnTheClock,
				Machineries: []spec.MachineryID{spec.ClockTheRunAdvances},
				Apply: func(h *ExpiringHandle) {
					probe.ExpiresAfter(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Read(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					}, time.Minute, ErrMissing)
				},
				Defect: expiringReadAnswersASample,
				Around: synctest.Test,
			},
		},
		Calls: []conform.Call[Expiring]{
			{
				Calls:    expiringPutCalls,
				Over:     newExpiringOver,
				Callable: ExpiringChecks.Put.Callable,
			},
			{
				Calls:    expiringReadCalls,
				Over:     newExpiringOver,
				Callable: ExpiringChecks.Read.Callable,
			},
		},
	}
}

// RunExpiringSuite applies the checks of Expiring, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunExpiringSuite[T Expiring](t *testing.T, newSubject func() T, opts ...ExpiringOption) {
	t.Helper()
	RunExpiringSuiteFromHandle(t, func(*ExpiringHandle) T { return newSubject() }, opts...)
}

// RunExpiringSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunExpiringSuiteFromHandle[T Expiring](
	t *testing.T,
	newSubject func(h *ExpiringHandle) T,
	opts ...ExpiringOption,
) {
	t.Helper()
	conform.Run(t, ExpiringSuite(), expiringSubjects(newSubject), opts...)
}

// FuzzExpiringSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzExpiringSuite[T Expiring](
	f *testing.F,
	callable conform.Callable[Expiring],
	newSubject func(h *ExpiringHandle) T,
	opts ...ExpiringOption,
) {
	f.Helper()
	conform.Fuzz(f, ExpiringSuite(), callable, expiringSubjects(newSubject), opts...)
}

// RunExpiringBudgets applies the allocation ceilings of Expiring and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunExpiringBudgets[T Expiring](
	t *testing.T,
	newSubject func(h *ExpiringHandle) T,
	opts ...ExpiringOption,
) {
	t.Helper()
	conform.RunBudgets(t, ExpiringSuite(), expiringSubjects(newSubject), opts...)
}

// BenchExpiringSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchExpiringSuite[T Expiring](
	b conform.B,
	callable conform.Callable[Expiring],
	newSubject func(h *ExpiringHandle) T,
	opts ...ExpiringOption,
) {
	b.Helper()
	conform.Bench(b, ExpiringSuite(), callable, expiringSubjects(newSubject), opts...)
}

// expiringSubjects returns the factory of the subjects of Expiring that
// calls newSubject with the handle of each check.
func expiringSubjects[T Expiring](
	newSubject func(h *ExpiringHandle) T,
) func(h *ExpiringHandle) Expiring {
	return func(h *ExpiringHandle) Expiring { return newSubject(h) }
}

// expiringPutCalls generates the calls of Expiring.Put with the inputs
// that each case draws.
var expiringPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Expiring) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Expiring) error {
		return subject.Put(ctx, v)
	}
})

// expiringReadCalls generates the calls of Expiring.Read with the inputs
// that each case draws.
var expiringReadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Expiring) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Expiring) error {
		_, err := subject.Read(ctx, key)
		return err
	}
})

// expiringOverEmbedded is Expiring under the name of the field that
// expiringOver embeds.
type expiringOverEmbedded = Expiring

// expiringOver is a subject of Expiring whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type expiringOver struct {
	expiringOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o expiringOver) Put(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Read calls spend and returns its zero results.
func (o expiringOver) Read(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newExpiringOver returns the planted defect of a ceiling of Expiring,
// whose callables call spend.
func newExpiringOver(_ *ExpiringHandle, spend func()) Expiring {
	return expiringOver{spend: spend}
}

// expiringPutPanicsDefect is the planted defect of a Expiring whose Put
// panics.
type expiringPutPanicsDefect struct {
	expiringOverEmbedded
}

// Put is the override of the planted defect.
func (*expiringPutPanicsDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Expiring whose Put panics")
}

// expiringPutPanics returns the planted defect of a Expiring whose Put
// panics over a fresh subject.
func expiringPutPanics(h *ExpiringHandle) Expiring {
	return &expiringPutPanicsDefect{expiringOverEmbedded: h.Fresh()}
}

// expiringPutIgnoresCancelDefect is the planted defect of a Expiring whose
// Put ignores the context.
type expiringPutIgnoresCancelDefect struct {
	expiringOverEmbedded
}

// Put is the override of the planted defect.
func (d *expiringPutIgnoresCancelDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.expiringOverEmbedded.Put(p0, p1)
}

// expiringPutIgnoresCancel returns the planted defect of a Expiring whose
// Put ignores the context over a fresh subject.
func expiringPutIgnoresCancel(h *ExpiringHandle) Expiring {
	return &expiringPutIgnoresCancelDefect{expiringOverEmbedded: h.Fresh()}
}

// expiringReadPanicsDefect is the planted defect of a Expiring whose Read
// panics.
type expiringReadPanicsDefect struct {
	expiringOverEmbedded
}

// Read is the override of the planted defect.
func (*expiringReadPanicsDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Expiring whose Read panics")
}

// expiringReadPanics returns the planted defect of a Expiring whose Read
// panics over a fresh subject.
func expiringReadPanics(h *ExpiringHandle) Expiring {
	return &expiringReadPanicsDefect{expiringOverEmbedded: h.Fresh()}
}

// expiringReadIgnoresCancelDefect is the planted defect of a Expiring
// whose Read ignores the context.
type expiringReadIgnoresCancelDefect struct {
	expiringOverEmbedded
}

// Read is the override of the planted defect.
func (d *expiringReadIgnoresCancelDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.expiringOverEmbedded.Read(p0, p1)
}

// expiringReadIgnoresCancel returns the planted defect of a Expiring whose
// Read ignores the context over a fresh subject.
func expiringReadIgnoresCancel(h *ExpiringHandle) Expiring {
	return &expiringReadIgnoresCancelDefect{expiringOverEmbedded: h.Fresh()}
}

// expiringReadValueBesideDefect is the planted defect of a Expiring whose
// Read returns a value beside an error.
type expiringReadValueBesideDefect struct {
	expiringOverEmbedded
}

// Read is the override of the planted defect.
func (d *expiringReadValueBesideDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.expiringOverEmbedded.Read(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// expiringReadValueBeside returns the planted defect of a Expiring whose
// Read returns a value beside an error over a fresh subject.
func expiringReadValueBeside(h *ExpiringHandle) Expiring {
	return &expiringReadValueBesideDefect{expiringOverEmbedded: h.Fresh()}
}

// expiringReadAnswersASampleDefect is the planted defect of a Expiring
// whose Read returns an entry after its lifetime.
type expiringReadAnswersASampleDefect struct {
	expiringOverEmbedded
}

// Read is the override of the planted defect.
func (*expiringReadAnswersASampleDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = Value{Key: "Key-b"}
	return r0, r1
}

// expiringReadAnswersASample returns the planted defect of a Expiring
// whose Read returns an entry after its lifetime over a fresh subject.
func expiringReadAnswersASample(h *ExpiringHandle) Expiring {
	return &expiringReadAnswersASampleDefect{expiringOverEmbedded: h.Fresh()}
}

// The derived checks of Cached:
//
// Cached.Get/smoke: Get returns without a panic for derived inputs.
//
// Cached.Get/cancel: Get returns an error for a cancelled context.
//
// Cached.Get/zero-on-error: Get returns the zero value beside an error.
//
// Cached.Get/cache: Get returns what Fetch returns.
//
// Cached.Fetch/smoke: Fetch returns without a panic for derived inputs.
//
// Cached.Fetch/cancel: Fetch returns an error for a cancelled context.
//
// Cached.Fetch/zero-on-error: Fetch returns the zero value beside an
// error.
//
// Cached.Store/smoke: Store returns without a panic for derived inputs.
//
// Cached.Store/cancel: Store returns an error for a cancelled context.

// The derived checks of Cached that witness their laws at the shallow end:
//
// Cached.Get/smoke
//
// Cached.Get/cache
//
// Cached.Fetch/smoke
//
// Cached.Store/smoke

// The unclaimed classifications of Cached:
//
// The shape:writer of Cached.Store has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Cached writes.
type (
	CachedHandle = conform.Handle[Cached]
	CachedCheck  = conform.Check[Cached]
	CachedOption = conform.Option[Cached]
)

// CachedChecks contains the name of each callable of Cached, with the
// identities of its checks.
var CachedChecks = struct {
	Get   cachedGet
	Fetch cachedFetch
	Store cachedStore
}{
	Get:   cachedGet{"Cached.Get"},
	Fetch: cachedFetch{"Cached.Fetch"},
	Store: cachedStore{"Cached.Store"},
}

// cachedGet is the name of Cached.Get, with the identities of its checks.
type cachedGet struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c cachedGet) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c cachedGet) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c cachedGet) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// Cache returns the identity of the check that Get returns what Fetch
// returns.
func (c cachedGet) Cache() conform.ID[Cached] {
	return c.Check("cache")
}

// cachedFetch is the name of Cached.Fetch, with the identities of its
// checks.
type cachedFetch struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Fetch returns without a
// panic for derived inputs.
func (c cachedFetch) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Fetch returns an error for
// a cancelled context.
func (c cachedFetch) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Fetch returns the
// zero value beside an error.
func (c cachedFetch) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// cachedStore is the name of Cached.Store, with the identities of its
// checks.
type cachedStore struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Store returns without a
// panic for derived inputs.
func (c cachedStore) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Store returns an error for
// a cancelled context.
func (c cachedStore) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CachedSuite returns the checks of Cached, and the calls and the budgets
// of its callables.
func CachedSuite() conform.Suite[Cached] {
	return conform.Suite[Cached]{
		Declaration: "Cached",
		Checks: []conform.Check[Cached]{
			{
				ID:          CachedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedGetPanics,
			},
			{
				ID:          CachedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedGetIgnoresCancel,
			},
			{
				ID:          CachedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: cachedGetValueBeside,
			},
			{
				ID:          CachedChecks.Get.Cache(),
				Statement:   "Get returns what Fetch returns",
				Claim:       spec.Cache,
				Probe:       spec.AgreesWithTheBacking,
				Reason:      spec.CacheDisagrees,
				Law:         spec.CacheAgreesWithItsBacking,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.AgreesWithTheBacking(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) (Value, error) {
						return h.Subject().Fetch(ctx, string("key-a"))
					}, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: cachedGetAnswersASample,
			},
			{
				ID:          CachedChecks.Fetch.Smoke(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedFetchPanics,
			},
			{
				ID:          CachedChecks.Fetch.Cancel(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedFetchIgnoresCancel,
			},
			{
				ID:          CachedChecks.Fetch.ZeroOnError(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Fetch(ctx, string("key-a"))
					}, Value{})
				},
				Defect: cachedFetchValueBeside,
			},
			{
				ID:          CachedChecks.Store.Smoke(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: cachedStorePanics,
			},
			{
				ID:          CachedChecks.Store.Cancel(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: cachedStoreIgnoresCancel,
			},
		},
		Calls: []conform.Call[Cached]{
			{
				Calls:    cachedGetCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Get.Callable,
			},
			{
				Calls:    cachedFetchCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Fetch.Callable,
			},
			{
				Calls:    cachedStoreCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Store.Callable,
			},
		},
	}
}

// RunCachedSuite applies the checks of Cached, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCachedSuite[T Cached](t *testing.T, newSubject func() T, opts ...CachedOption) {
	t.Helper()
	RunCachedSuiteFromHandle(t, func(*CachedHandle) T { return newSubject() }, opts...)
}

// RunCachedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCachedSuiteFromHandle[T Cached](
	t *testing.T,
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	t.Helper()
	conform.Run(t, CachedSuite(), cachedSubjects(newSubject), opts...)
}

// FuzzCachedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCachedSuite[T Cached](
	f *testing.F,
	callable conform.Callable[Cached],
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	f.Helper()
	conform.Fuzz(f, CachedSuite(), callable, cachedSubjects(newSubject), opts...)
}

// RunCachedBudgets applies the allocation ceilings of Cached and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCachedBudgets[T Cached](
	t *testing.T,
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	t.Helper()
	conform.RunBudgets(t, CachedSuite(), cachedSubjects(newSubject), opts...)
}

// BenchCachedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchCachedSuite[T Cached](
	b conform.B,
	callable conform.Callable[Cached],
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	b.Helper()
	conform.Bench(b, CachedSuite(), callable, cachedSubjects(newSubject), opts...)
}

// cachedSubjects returns the factory of the subjects of Cached that calls
// newSubject with the handle of each check.
func cachedSubjects[T Cached](newSubject func(h *CachedHandle) T) func(h *CachedHandle) Cached {
	return func(h *CachedHandle) Cached { return newSubject(h) }
}

// cachedGetCalls generates the calls of Cached.Get with the inputs that
// each case draws.
var cachedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// cachedFetchCalls generates the calls of Cached.Fetch with the inputs
// that each case draws.
var cachedFetchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Fetch(ctx, key)
		return err
	}
})

// cachedStoreCalls generates the calls of Cached.Store with the inputs
// that each case draws.
var cachedStoreCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Cached) error {
		return subject.Store(ctx, v)
	}
})

// cachedOverEmbedded is Cached under the name of the field that cachedOver
// embeds.
type cachedOverEmbedded = Cached

// cachedOver is a subject of Cached whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type cachedOver struct {
	cachedOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o cachedOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Fetch calls spend and returns its zero results.
func (o cachedOver) Fetch(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Store calls spend and returns its zero results.
func (o cachedOver) Store(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// newCachedOver returns the planted defect of a ceiling of Cached, whose
// callables call spend.
func newCachedOver(_ *CachedHandle, spend func()) Cached {
	return cachedOver{spend: spend}
}

// cachedGetPanicsDefect is the planted defect of a Cached whose Get
// panics.
type cachedGetPanicsDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (*cachedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Cached whose Get panics")
}

// cachedGetPanics returns the planted defect of a Cached whose Get panics
// over a fresh subject.
func cachedGetPanics(h *CachedHandle) Cached {
	return &cachedGetPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetIgnoresCancelDefect is the planted defect of a Cached whose Get
// ignores the context.
type cachedGetIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (d *cachedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Get(p0, p1)
}

// cachedGetIgnoresCancel returns the planted defect of a Cached whose Get
// ignores the context over a fresh subject.
func cachedGetIgnoresCancel(h *CachedHandle) Cached {
	return &cachedGetIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetValueBesideDefect is the planted defect of a Cached whose Get
// returns a value beside an error.
type cachedGetValueBesideDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (d *cachedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// cachedGetValueBeside returns the planted defect of a Cached whose Get
// returns a value beside an error over a fresh subject.
func cachedGetValueBeside(h *CachedHandle) Cached {
	return &cachedGetValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetAnswersASampleDefect is the planted defect of a Cached whose
// Get returns a value that Fetch never had.
type cachedGetAnswersASampleDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (*cachedGetAnswersASampleDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = Value{Key: "Key-b"}
	return r0, r1
}

// cachedGetAnswersASample returns the planted defect of a Cached whose Get
// returns a value that Fetch never had over a fresh subject.
func cachedGetAnswersASample(h *CachedHandle) Cached {
	return &cachedGetAnswersASampleDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchPanicsDefect is the planted defect of a Cached whose Fetch
// panics.
type cachedFetchPanicsDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (*cachedFetchPanicsDefect) Fetch(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Cached whose Fetch panics")
}

// cachedFetchPanics returns the planted defect of a Cached whose Fetch
// panics over a fresh subject.
func cachedFetchPanics(h *CachedHandle) Cached {
	return &cachedFetchPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchIgnoresCancelDefect is the planted defect of a Cached whose
// Fetch ignores the context.
type cachedFetchIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *cachedFetchIgnoresCancelDefect) Fetch(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Fetch(p0, p1)
}

// cachedFetchIgnoresCancel returns the planted defect of a Cached whose
// Fetch ignores the context over a fresh subject.
func cachedFetchIgnoresCancel(h *CachedHandle) Cached {
	return &cachedFetchIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchValueBesideDefect is the planted defect of a Cached whose
// Fetch returns a value beside an error.
type cachedFetchValueBesideDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *cachedFetchValueBesideDefect) Fetch(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Fetch(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// cachedFetchValueBeside returns the planted defect of a Cached whose
// Fetch returns a value beside an error over a fresh subject.
func cachedFetchValueBeside(h *CachedHandle) Cached {
	return &cachedFetchValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedStorePanicsDefect is the planted defect of a Cached whose Store
// panics.
type cachedStorePanicsDefect struct {
	cachedOverEmbedded
}

// Store is the override of the planted defect.
func (*cachedStorePanicsDefect) Store(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Cached whose Store panics")
}

// cachedStorePanics returns the planted defect of a Cached whose Store
// panics over a fresh subject.
func cachedStorePanics(h *CachedHandle) Cached {
	return &cachedStorePanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedStoreIgnoresCancelDefect is the planted defect of a Cached whose
// Store ignores the context.
type cachedStoreIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Store is the override of the planted defect.
func (d *cachedStoreIgnoresCancelDefect) Store(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.cachedOverEmbedded.Store(p0, p1)
}

// cachedStoreIgnoresCancel returns the planted defect of a Cached whose
// Store ignores the context over a fresh subject.
func cachedStoreIgnoresCancel(h *CachedHandle) Cached {
	return &cachedStoreIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// The derived checks of Tagged:
//
// Tagged.Put/smoke: Put returns without a panic for derived inputs.
//
// Tagged.Put/cancel: Put returns an error for a cancelled context.
//
// Tagged.Get/smoke: Get returns without a panic for derived inputs.
//
// Tagged.Get/cancel: Get returns an error for a cancelled context.
//
// Tagged.Get/zero-on-error: Get returns the zero value beside an error.
//
// Tagged.Get/read-after-write: Get returns what Put stored before it.

// The derived checks of Tagged that witness their laws at the shallow end:
//
// Tagged.Put/smoke
//
// Tagged.Get/smoke
//
// Tagged.Get/read-after-write

// The unclaimed classifications of Tagged:
//
// The shape:composite-writer of Tagged.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Tagged writes.
type (
	TaggedHandle = conform.Handle[Tagged]
	TaggedCheck  = conform.Check[Tagged]
	TaggedOption = conform.Option[Tagged]
)

// TaggedChecks contains the name of each callable of Tagged, with the
// identities of its checks.
var TaggedChecks = struct {
	Put taggedPut
	Get taggedGet
}{
	Put: taggedPut{"Tagged.Put"},
	Get: taggedGet{"Tagged.Get"},
}

// taggedPut is the name of Tagged.Put, with the identities of its checks.
type taggedPut struct {
	conform.Callable[Tagged]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c taggedPut) Smoke() conform.ID[Tagged] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c taggedPut) Cancel() conform.ID[Tagged] {
	return c.Check("cancel")
}

// taggedGet is the name of Tagged.Get, with the identities of its checks.
type taggedGet struct {
	conform.Callable[Tagged]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c taggedGet) Smoke() conform.ID[Tagged] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c taggedGet) Cancel() conform.ID[Tagged] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c taggedGet) ZeroOnError() conform.ID[Tagged] {
	return c.Check("zero-on-error")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c taggedGet) ReadAfterWrite() conform.ID[Tagged] {
	return c.Check("read-after-write")
}

// TaggedSuite returns the checks of Tagged, and the calls and the budgets
// of its callables.
func TaggedSuite() conform.Suite[Tagged] {
	return conform.Suite[Tagged]{
		Declaration: "Tagged",
		Checks: []conform.Check[Tagged]{
			{
				ID:          TaggedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TaggedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"}, string("tag-a"))
					})
				},
				Defect: taggedPutPanics,
			},
			{
				ID:          TaggedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TaggedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"}, string("tag-a"))
					})
				},
				Defect: taggedPutIgnoresCancel,
			},
			{
				ID:          TaggedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TaggedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("tag-a"))
						return err
					})
				},
				Defect: taggedGetPanics,
			},
			{
				ID:          TaggedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TaggedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("tag-a"))
						return err
					})
				},
				Defect: taggedGetIgnoresCancel,
			},
			{
				ID:          TaggedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TaggedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("tag-a"))
					}, Value{})
				},
				Defect: taggedGetValueBeside,
			},
			{
				ID:          TaggedChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TaggedHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("tag-a"))
					}, func(ctx context.Context) (Value, error) {
						w := Value{Key: "Key-a"}
						err := h.Subject().Put(ctx, w, string("tag-a"))
						return w, err
					})
				},
				Defect: taggedGetReadsAnotherValue,
			},
		},
		Calls: []conform.Call[Tagged]{
			{
				Calls:    taggedPutCalls,
				Over:     newTaggedOver,
				Callable: TaggedChecks.Put.Callable,
			},
			{
				Calls:    taggedGetCalls,
				Over:     newTaggedOver,
				Callable: TaggedChecks.Get.Callable,
			},
		},
	}
}

// RunTaggedSuite applies the checks of Tagged, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunTaggedSuite[T Tagged](t *testing.T, newSubject func() T, opts ...TaggedOption) {
	t.Helper()
	RunTaggedSuiteFromHandle(t, func(*TaggedHandle) T { return newSubject() }, opts...)
}

// RunTaggedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTaggedSuiteFromHandle[T Tagged](
	t *testing.T,
	newSubject func(h *TaggedHandle) T,
	opts ...TaggedOption,
) {
	t.Helper()
	conform.Run(t, TaggedSuite(), taggedSubjects(newSubject), opts...)
}

// FuzzTaggedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTaggedSuite[T Tagged](
	f *testing.F,
	callable conform.Callable[Tagged],
	newSubject func(h *TaggedHandle) T,
	opts ...TaggedOption,
) {
	f.Helper()
	conform.Fuzz(f, TaggedSuite(), callable, taggedSubjects(newSubject), opts...)
}

// RunTaggedBudgets applies the allocation ceilings of Tagged and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunTaggedBudgets[T Tagged](
	t *testing.T,
	newSubject func(h *TaggedHandle) T,
	opts ...TaggedOption,
) {
	t.Helper()
	conform.RunBudgets(t, TaggedSuite(), taggedSubjects(newSubject), opts...)
}

// BenchTaggedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchTaggedSuite[T Tagged](
	b conform.B,
	callable conform.Callable[Tagged],
	newSubject func(h *TaggedHandle) T,
	opts ...TaggedOption,
) {
	b.Helper()
	conform.Bench(b, TaggedSuite(), callable, taggedSubjects(newSubject), opts...)
}

// taggedSubjects returns the factory of the subjects of Tagged that calls
// newSubject with the handle of each check.
func taggedSubjects[T Tagged](newSubject func(h *TaggedHandle) T) func(h *TaggedHandle) Tagged {
	return func(h *TaggedHandle) Tagged { return newSubject(h) }
}

// taggedPutCalls generates the calls of Tagged.Put with the inputs that
// each case draws.
var taggedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tagged) error {
	v := c.Draw(prop.Of[Value](), "v")
	tag := c.Draw(prop.Of[string](), "tag")
	return func(ctx context.Context, subject Tagged) error {
		return subject.Put(ctx, v, tag)
	}
})

// taggedGetCalls generates the calls of Tagged.Get with the inputs that
// each case draws.
var taggedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tagged) error {
	tag := c.Draw(prop.Of[string](), "tag")
	return func(ctx context.Context, subject Tagged) error {
		_, err := subject.Get(ctx, tag)
		return err
	}
})

// taggedOverEmbedded is Tagged under the name of the field that taggedOver
// embeds.
type taggedOverEmbedded = Tagged

// taggedOver is a subject of Tagged whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type taggedOver struct {
	taggedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o taggedOver) Put(context.Context, Value, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o taggedOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newTaggedOver returns the planted defect of a ceiling of Tagged, whose
// callables call spend.
func newTaggedOver(_ *TaggedHandle, spend func()) Tagged {
	return taggedOver{spend: spend}
}

// taggedPutPanicsDefect is the planted defect of a Tagged whose Put
// panics.
type taggedPutPanicsDefect struct {
	taggedOverEmbedded
}

// Put is the override of the planted defect.
func (*taggedPutPanicsDefect) Put(p0 context.Context, p1 Value, p2 string) (r0 error) {
	panic("the planted defect of a Tagged whose Put panics")
}

// taggedPutPanics returns the planted defect of a Tagged whose Put panics
// over a fresh subject.
func taggedPutPanics(h *TaggedHandle) Tagged {
	return &taggedPutPanicsDefect{taggedOverEmbedded: h.Fresh()}
}

// taggedPutIgnoresCancelDefect is the planted defect of a Tagged whose Put
// ignores the context.
type taggedPutIgnoresCancelDefect struct {
	taggedOverEmbedded
}

// Put is the override of the planted defect.
func (d *taggedPutIgnoresCancelDefect) Put(p0 context.Context, p1 Value, p2 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.taggedOverEmbedded.Put(p0, p1, p2)
}

// taggedPutIgnoresCancel returns the planted defect of a Tagged whose Put
// ignores the context over a fresh subject.
func taggedPutIgnoresCancel(h *TaggedHandle) Tagged {
	return &taggedPutIgnoresCancelDefect{taggedOverEmbedded: h.Fresh()}
}

// taggedGetPanicsDefect is the planted defect of a Tagged whose Get
// panics.
type taggedGetPanicsDefect struct {
	taggedOverEmbedded
}

// Get is the override of the planted defect.
func (*taggedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Tagged whose Get panics")
}

// taggedGetPanics returns the planted defect of a Tagged whose Get panics
// over a fresh subject.
func taggedGetPanics(h *TaggedHandle) Tagged {
	return &taggedGetPanicsDefect{taggedOverEmbedded: h.Fresh()}
}

// taggedGetIgnoresCancelDefect is the planted defect of a Tagged whose Get
// ignores the context.
type taggedGetIgnoresCancelDefect struct {
	taggedOverEmbedded
}

// Get is the override of the planted defect.
func (d *taggedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.taggedOverEmbedded.Get(p0, p1)
}

// taggedGetIgnoresCancel returns the planted defect of a Tagged whose Get
// ignores the context over a fresh subject.
func taggedGetIgnoresCancel(h *TaggedHandle) Tagged {
	return &taggedGetIgnoresCancelDefect{taggedOverEmbedded: h.Fresh()}
}

// taggedGetValueBesideDefect is the planted defect of a Tagged whose Get
// returns a value beside an error.
type taggedGetValueBesideDefect struct {
	taggedOverEmbedded
}

// Get is the override of the planted defect.
func (d *taggedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.taggedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// taggedGetValueBeside returns the planted defect of a Tagged whose Get
// returns a value beside an error over a fresh subject.
func taggedGetValueBeside(h *TaggedHandle) Tagged {
	return &taggedGetValueBesideDefect{taggedOverEmbedded: h.Fresh()}
}

// taggedGetReadsAnotherValueDefect is the planted defect of a Tagged whose
// Get returns a value that Put did not write.
type taggedGetReadsAnotherValueDefect struct {
	taggedOverEmbedded
}

// Get is the override of the planted defect.
func (d *taggedGetReadsAnotherValueDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.taggedOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// taggedGetReadsAnotherValue returns the planted defect of a Tagged whose
// Get returns a value that Put did not write over a fresh subject.
func taggedGetReadsAnotherValue(h *TaggedHandle) Tagged {
	return &taggedGetReadsAnotherValueDefect{taggedOverEmbedded: h.Fresh()}
}

// The derived checks of Snapshots:
//
// Snapshots.Store/smoke: Store returns without a panic for derived inputs.
//
// Snapshots.Store/cancel: Store returns an error for a cancelled context.
//
// Snapshots.Get/smoke: Get returns without a panic for derived inputs.
//
// Snapshots.Get/cancel: Get returns an error for a cancelled context.
//
// Snapshots.Get/zero-on-error: Get returns the zero value beside an error.
//
// Snapshots.Get/point-in-time: Get returns the same value twice across a
// Store between the reads.

// The derived checks of Snapshots that witness their laws at the shallow
// end:
//
// Snapshots.Store/smoke
//
// Snapshots.Get/smoke
//
// Snapshots.Get/point-in-time

// The unclaimed classifications of Snapshots:
//
// The shape:writer of Snapshots.Store has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Snapshots writes.
type (
	SnapshotsHandle = conform.Handle[Snapshots]
	SnapshotsCheck  = conform.Check[Snapshots]
	SnapshotsOption = conform.Option[Snapshots]
)

// SnapshotsChecks contains the name of each callable of Snapshots, with
// the identities of its checks.
var SnapshotsChecks = struct {
	Store snapshotsStore
	Get   snapshotsGet
}{
	Store: snapshotsStore{"Snapshots.Store"},
	Get:   snapshotsGet{"Snapshots.Get"},
}

// snapshotsStore is the name of Snapshots.Store, with the identities of
// its checks.
type snapshotsStore struct {
	conform.Callable[Snapshots]
}

// Smoke returns the identity of the check that Store returns without a
// panic for derived inputs.
func (c snapshotsStore) Smoke() conform.ID[Snapshots] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Store returns an error for
// a cancelled context.
func (c snapshotsStore) Cancel() conform.ID[Snapshots] {
	return c.Check("cancel")
}

// snapshotsGet is the name of Snapshots.Get, with the identities of its
// checks.
type snapshotsGet struct {
	conform.Callable[Snapshots]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c snapshotsGet) Smoke() conform.ID[Snapshots] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c snapshotsGet) Cancel() conform.ID[Snapshots] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c snapshotsGet) ZeroOnError() conform.ID[Snapshots] {
	return c.Check("zero-on-error")
}

// PointInTime returns the identity of the check that Get returns the same
// value twice across a Store between the reads.
func (c snapshotsGet) PointInTime() conform.ID[Snapshots] {
	return c.Check("point-in-time")
}

// SnapshotsSuite returns the checks of Snapshots, and the calls and the
// budgets of its callables.
func SnapshotsSuite() conform.Suite[Snapshots] {
	return conform.Suite[Snapshots]{
		Declaration: "Snapshots",
		Checks: []conform.Check[Snapshots]{
			{
				ID:          SnapshotsChecks.Store.Smoke(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SnapshotsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: snapshotsStorePanics,
			},
			{
				ID:          SnapshotsChecks.Store.Cancel(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SnapshotsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: snapshotsStoreIgnoresCancel,
			},
			{
				ID:          SnapshotsChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SnapshotsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: snapshotsGetPanics,
			},
			{
				ID:          SnapshotsChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SnapshotsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: snapshotsGetIgnoresCancel,
			},
			{
				ID:          SnapshotsChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SnapshotsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: snapshotsGetValueBeside,
			},
			{
				ID:          SnapshotsChecks.Get.PointInTime(),
				Statement:   "Get returns the same value twice across a Store between the reads",
				Claim:       spec.PointInTime,
				Probe:       spec.PinnedAcrossWrite,
				Reason:      spec.UnpinnedRead,
				Law:         spec.ReadKeepsItsSnapshot,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SnapshotsHandle) {
					probe.PinnedAcrossWrite(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) error {
						w := Value{Key: "Key-a"}
						w.Body = string("Body-b")
						return h.Subject().Store(ctx, w)
					})
				},
				Defect: snapshotsGetDrifts,
			},
		},
		Calls: []conform.Call[Snapshots]{
			{
				Calls:    snapshotsStoreCalls,
				Over:     newSnapshotsOver,
				Callable: SnapshotsChecks.Store.Callable,
			},
			{
				Calls:    snapshotsGetCalls,
				Over:     newSnapshotsOver,
				Callable: SnapshotsChecks.Get.Callable,
			},
		},
	}
}

// RunSnapshotsSuite applies the checks of Snapshots, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunSnapshotsSuite[T Snapshots](t *testing.T, newSubject func() T, opts ...SnapshotsOption) {
	t.Helper()
	RunSnapshotsSuiteFromHandle(t, func(*SnapshotsHandle) T { return newSubject() }, opts...)
}

// RunSnapshotsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunSnapshotsSuiteFromHandle[T Snapshots](
	t *testing.T,
	newSubject func(h *SnapshotsHandle) T,
	opts ...SnapshotsOption,
) {
	t.Helper()
	conform.Run(t, SnapshotsSuite(), snapshotsSubjects(newSubject), opts...)
}

// FuzzSnapshotsSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzSnapshotsSuite[T Snapshots](
	f *testing.F,
	callable conform.Callable[Snapshots],
	newSubject func(h *SnapshotsHandle) T,
	opts ...SnapshotsOption,
) {
	f.Helper()
	conform.Fuzz(f, SnapshotsSuite(), callable, snapshotsSubjects(newSubject), opts...)
}

// RunSnapshotsBudgets applies the allocation ceilings of Snapshots and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunSnapshotsBudgets[T Snapshots](
	t *testing.T,
	newSubject func(h *SnapshotsHandle) T,
	opts ...SnapshotsOption,
) {
	t.Helper()
	conform.RunBudgets(t, SnapshotsSuite(), snapshotsSubjects(newSubject), opts...)
}

// BenchSnapshotsSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchSnapshotsSuite[T Snapshots](
	b conform.B,
	callable conform.Callable[Snapshots],
	newSubject func(h *SnapshotsHandle) T,
	opts ...SnapshotsOption,
) {
	b.Helper()
	conform.Bench(b, SnapshotsSuite(), callable, snapshotsSubjects(newSubject), opts...)
}

// snapshotsSubjects returns the factory of the subjects of Snapshots that
// calls newSubject with the handle of each check.
func snapshotsSubjects[T Snapshots](
	newSubject func(h *SnapshotsHandle) T,
) func(h *SnapshotsHandle) Snapshots {
	return func(h *SnapshotsHandle) Snapshots { return newSubject(h) }
}

// snapshotsStoreCalls generates the calls of Snapshots.Store with the
// inputs that each case draws.
var snapshotsStoreCalls = prop.Composite(func(c *prop.Case) func(context.Context, Snapshots) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Snapshots) error {
		return subject.Store(ctx, v)
	}
})

// snapshotsGetCalls generates the calls of Snapshots.Get with the inputs
// that each case draws.
var snapshotsGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Snapshots) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Snapshots) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// snapshotsOverEmbedded is Snapshots under the name of the field that
// snapshotsOver embeds.
type snapshotsOverEmbedded = Snapshots

// snapshotsOver is a subject of Snapshots whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type snapshotsOver struct {
	snapshotsOverEmbedded
	spend func()
}

// Store calls spend and returns its zero results.
func (o snapshotsOver) Store(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o snapshotsOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newSnapshotsOver returns the planted defect of a ceiling of Snapshots,
// whose callables call spend.
func newSnapshotsOver(_ *SnapshotsHandle, spend func()) Snapshots {
	return snapshotsOver{spend: spend}
}

// snapshotsStorePanicsDefect is the planted defect of a Snapshots whose
// Store panics.
type snapshotsStorePanicsDefect struct {
	snapshotsOverEmbedded
}

// Store is the override of the planted defect.
func (*snapshotsStorePanicsDefect) Store(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Snapshots whose Store panics")
}

// snapshotsStorePanics returns the planted defect of a Snapshots whose
// Store panics over a fresh subject.
func snapshotsStorePanics(h *SnapshotsHandle) Snapshots {
	return &snapshotsStorePanicsDefect{snapshotsOverEmbedded: h.Fresh()}
}

// snapshotsStoreIgnoresCancelDefect is the planted defect of a Snapshots
// whose Store ignores the context.
type snapshotsStoreIgnoresCancelDefect struct {
	snapshotsOverEmbedded
}

// Store is the override of the planted defect.
func (d *snapshotsStoreIgnoresCancelDefect) Store(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.snapshotsOverEmbedded.Store(p0, p1)
}

// snapshotsStoreIgnoresCancel returns the planted defect of a Snapshots
// whose Store ignores the context over a fresh subject.
func snapshotsStoreIgnoresCancel(h *SnapshotsHandle) Snapshots {
	return &snapshotsStoreIgnoresCancelDefect{snapshotsOverEmbedded: h.Fresh()}
}

// snapshotsGetPanicsDefect is the planted defect of a Snapshots whose Get
// panics.
type snapshotsGetPanicsDefect struct {
	snapshotsOverEmbedded
}

// Get is the override of the planted defect.
func (*snapshotsGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Snapshots whose Get panics")
}

// snapshotsGetPanics returns the planted defect of a Snapshots whose Get
// panics over a fresh subject.
func snapshotsGetPanics(h *SnapshotsHandle) Snapshots {
	return &snapshotsGetPanicsDefect{snapshotsOverEmbedded: h.Fresh()}
}

// snapshotsGetIgnoresCancelDefect is the planted defect of a Snapshots
// whose Get ignores the context.
type snapshotsGetIgnoresCancelDefect struct {
	snapshotsOverEmbedded
}

// Get is the override of the planted defect.
func (d *snapshotsGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.snapshotsOverEmbedded.Get(p0, p1)
}

// snapshotsGetIgnoresCancel returns the planted defect of a Snapshots
// whose Get ignores the context over a fresh subject.
func snapshotsGetIgnoresCancel(h *SnapshotsHandle) Snapshots {
	return &snapshotsGetIgnoresCancelDefect{snapshotsOverEmbedded: h.Fresh()}
}

// snapshotsGetValueBesideDefect is the planted defect of a Snapshots whose
// Get returns a value beside an error.
type snapshotsGetValueBesideDefect struct {
	snapshotsOverEmbedded
}

// Get is the override of the planted defect.
func (d *snapshotsGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.snapshotsOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// snapshotsGetValueBeside returns the planted defect of a Snapshots whose
// Get returns a value beside an error over a fresh subject.
func snapshotsGetValueBeside(h *SnapshotsHandle) Snapshots {
	return &snapshotsGetValueBesideDefect{snapshotsOverEmbedded: h.Fresh()}
}

// snapshotsGetDriftsDefect is the planted defect of a Snapshots whose Get
// returns another value on the second read.
type snapshotsGetDriftsDefect struct {
	snapshotsOverEmbedded
	calls int
}

// Get is the override of the planted defect.
func (d *snapshotsGetDriftsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls > 1 {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// snapshotsGetDrifts returns the planted defect of a Snapshots whose Get
// returns another value on the second read over a fresh subject.
func snapshotsGetDrifts(h *SnapshotsHandle) Snapshots {
	return &snapshotsGetDriftsDefect{snapshotsOverEmbedded: h.Fresh()}
}

// The derived checks of Bag:
//
// Bag.Add/smoke: Add returns without a panic for derived inputs.
//
// Bag.Add/cancel: Add returns an error for a cancelled context.
//
// Bag.Items/smoke: Items returns without a panic for derived inputs.
//
// Bag.Items/cancel: Items returns an error for a cancelled context.
//
// Bag.Items/zero-on-error: Items returns the zero value beside an error.
//
// Bag.Items/no-duplicates: no element appears twice in one Items.
//
// Bag.Items/stable-order: two calls of Items return the elements in the
// same order.

// The derived checks of Bag that witness their laws at the shallow end:
//
// Bag.Add/smoke
//
// Bag.Items/smoke
//
// Bag.Items/no-duplicates
//
// Bag.Items/stable-order

// The unclaimed classifications of Bag:
//
// The shape:writer of Bag.Add has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Bag writes.
type (
	BagHandle = conform.Handle[Bag]
	BagCheck  = conform.Check[Bag]
	BagOption = conform.Option[Bag]
)

// BagChecks contains the name of each callable of Bag, with the identities
// of its checks.
var BagChecks = struct {
	Add   bagAdd
	Items bagItems
}{
	Add:   bagAdd{"Bag.Add"},
	Items: bagItems{"Bag.Items"},
}

// bagAdd is the name of Bag.Add, with the identities of its checks.
type bagAdd struct {
	conform.Callable[Bag]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c bagAdd) Smoke() conform.ID[Bag] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c bagAdd) Cancel() conform.ID[Bag] {
	return c.Check("cancel")
}

// bagItems is the name of Bag.Items, with the identities of its checks.
type bagItems struct {
	conform.Callable[Bag]
}

// Smoke returns the identity of the check that Items returns without a
// panic for derived inputs.
func (c bagItems) Smoke() conform.ID[Bag] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Items returns an error for
// a cancelled context.
func (c bagItems) Cancel() conform.ID[Bag] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Items returns the
// zero value beside an error.
func (c bagItems) ZeroOnError() conform.ID[Bag] {
	return c.Check("zero-on-error")
}

// NoDuplicates returns the identity of the check that no element appears
// twice in one Items.
func (c bagItems) NoDuplicates() conform.ID[Bag] {
	return c.Check("no-duplicates")
}

// StableOrder returns the identity of the check that two calls of Items
// return the elements in the same order.
func (c bagItems) StableOrder() conform.ID[Bag] {
	return c.Check("stable-order")
}

// BagSuite returns the checks of Bag, and the calls and the budgets of its
// callables.
func BagSuite() conform.Suite[Bag] {
	return conform.Suite[Bag]{
		Declaration: "Bag",
		Checks: []conform.Check[Bag]{
			{
				ID:          BagChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BagHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: bagAddPanics,
			},
			{
				ID:          BagChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BagHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: bagAddIgnoresCancel,
			},
			{
				ID:          BagChecks.Items.Smoke(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BagHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx)
						return err
					})
				},
				Defect: bagItemsPanics,
			},
			{
				ID:          BagChecks.Items.Cancel(),
				Statement:   "Items returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BagHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx)
						return err
					})
				},
				Defect: bagItemsIgnoresCancel,
			},
			{
				ID:          BagChecks.Items.ZeroOnError(),
				Statement:   "Items returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BagHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]Value, error) {
						return h.Subject().Items(ctx)
					}, []Value(nil))
				},
				Defect: bagItemsValueBeside,
			},
			{
				ID:          BagChecks.Items.NoDuplicates(),
				Statement:   "no element appears twice in one Items",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BagHandle) {
					probe.NoElementTwice(h, func(ctx context.Context) ([]Value, error) {
						return h.Subject().Items(ctx)
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Value{Key: "Key-a"}
							if calls > 1 {
								w = Value{Key: "Key-b"}
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: bagItemsDrainsTwice,
			},
			{
				ID:          BagChecks.Items.StableOrder(),
				Statement:   "two calls of Items return the elements in the same order",
				Claim:       spec.StableOrder,
				Probe:       spec.DrainsInOneOrder,
				Reason:      spec.ReorderedDrain,
				Law:         spec.TwoDrainsAgreeOnOrder,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BagHandle) {
					probe.DrainsInOneOrder(h, func(ctx context.Context) ([]Value, error) {
						return h.Subject().Items(ctx)
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Value{Key: "Key-a"}
							if calls > 1 {
								w = Value{Key: "Key-b"}
							}
							return h.Subject().Add(ctx, w)
						}
					}())
				},
				Defect: bagItemsReordersDrain,
			},
		},
		Calls: []conform.Call[Bag]{
			{
				Calls:    bagAddCalls,
				Over:     newBagOver,
				Callable: BagChecks.Add.Callable,
			},
			{
				Calls:    bagItemsCalls,
				Over:     newBagOver,
				Callable: BagChecks.Items.Callable,
			},
		},
	}
}

// RunBagSuite applies the checks of Bag, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunBagSuite[T Bag](t *testing.T, newSubject func() T, opts ...BagOption) {
	t.Helper()
	RunBagSuiteFromHandle(t, func(*BagHandle) T { return newSubject() }, opts...)
}

// RunBagSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunBagSuiteFromHandle[T Bag](
	t *testing.T,
	newSubject func(h *BagHandle) T,
	opts ...BagOption,
) {
	t.Helper()
	conform.Run(t, BagSuite(), bagSubjects(newSubject), opts...)
}

// FuzzBagSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzBagSuite[T Bag](
	f *testing.F,
	callable conform.Callable[Bag],
	newSubject func(h *BagHandle) T,
	opts ...BagOption,
) {
	f.Helper()
	conform.Fuzz(f, BagSuite(), callable, bagSubjects(newSubject), opts...)
}

// RunBagBudgets applies the allocation ceilings of Bag and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunBagBudgets[T Bag](t *testing.T, newSubject func(h *BagHandle) T, opts ...BagOption) {
	t.Helper()
	conform.RunBudgets(t, BagSuite(), bagSubjects(newSubject), opts...)
}

// BenchBagSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchBagSuite[T Bag](
	b conform.B,
	callable conform.Callable[Bag],
	newSubject func(h *BagHandle) T,
	opts ...BagOption,
) {
	b.Helper()
	conform.Bench(b, BagSuite(), callable, bagSubjects(newSubject), opts...)
}

// bagSubjects returns the factory of the subjects of Bag that calls
// newSubject with the handle of each check.
func bagSubjects[T Bag](newSubject func(h *BagHandle) T) func(h *BagHandle) Bag {
	return func(h *BagHandle) Bag { return newSubject(h) }
}

// bagAddCalls generates the calls of Bag.Add with the inputs that each
// case draws.
var bagAddCalls = prop.Composite(func(c *prop.Case) func(context.Context, Bag) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Bag) error {
		return subject.Add(ctx, v)
	}
})

// bagItemsCalls generates the calls of Bag.Items with the inputs that each
// case draws.
var bagItemsCalls = prop.Composite(func(*prop.Case) func(context.Context, Bag) error {
	return func(ctx context.Context, subject Bag) error {
		_, err := subject.Items(ctx)
		return err
	}
})

// bagOverEmbedded is Bag under the name of the field that bagOver embeds.
type bagOverEmbedded = Bag

// bagOver is a subject of Bag whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type bagOver struct {
	bagOverEmbedded
	spend func()
}

// Add calls spend and returns its zero results.
func (o bagOver) Add(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Items calls spend and returns its zero results.
func (o bagOver) Items(context.Context) (r0 []Value, r1 error) {
	o.spend()
	return r0, r1
}

// newBagOver returns the planted defect of a ceiling of Bag, whose
// callables call spend.
func newBagOver(_ *BagHandle, spend func()) Bag {
	return bagOver{spend: spend}
}

// bagAddPanicsDefect is the planted defect of a Bag whose Add panics.
type bagAddPanicsDefect struct {
	bagOverEmbedded
}

// Add is the override of the planted defect.
func (*bagAddPanicsDefect) Add(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Bag whose Add panics")
}

// bagAddPanics returns the planted defect of a Bag whose Add panics over a
// fresh subject.
func bagAddPanics(h *BagHandle) Bag {
	return &bagAddPanicsDefect{bagOverEmbedded: h.Fresh()}
}

// bagAddIgnoresCancelDefect is the planted defect of a Bag whose Add
// ignores the context.
type bagAddIgnoresCancelDefect struct {
	bagOverEmbedded
}

// Add is the override of the planted defect.
func (d *bagAddIgnoresCancelDefect) Add(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.bagOverEmbedded.Add(p0, p1)
}

// bagAddIgnoresCancel returns the planted defect of a Bag whose Add
// ignores the context over a fresh subject.
func bagAddIgnoresCancel(h *BagHandle) Bag {
	return &bagAddIgnoresCancelDefect{bagOverEmbedded: h.Fresh()}
}

// bagItemsPanicsDefect is the planted defect of a Bag whose Items panics.
type bagItemsPanicsDefect struct {
	bagOverEmbedded
}

// Items is the override of the planted defect.
func (*bagItemsPanicsDefect) Items(p0 context.Context) (r0 []Value, r1 error) {
	panic("the planted defect of a Bag whose Items panics")
}

// bagItemsPanics returns the planted defect of a Bag whose Items panics
// over a fresh subject.
func bagItemsPanics(h *BagHandle) Bag {
	return &bagItemsPanicsDefect{bagOverEmbedded: h.Fresh()}
}

// bagItemsIgnoresCancelDefect is the planted defect of a Bag whose Items
// ignores the context.
type bagItemsIgnoresCancelDefect struct {
	bagOverEmbedded
}

// Items is the override of the planted defect.
func (d *bagItemsIgnoresCancelDefect) Items(p0 context.Context) (r0 []Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.bagOverEmbedded.Items(p0)
}

// bagItemsIgnoresCancel returns the planted defect of a Bag whose Items
// ignores the context over a fresh subject.
func bagItemsIgnoresCancel(h *BagHandle) Bag {
	return &bagItemsIgnoresCancelDefect{bagOverEmbedded: h.Fresh()}
}

// bagItemsValueBesideDefect is the planted defect of a Bag whose Items
// returns a value beside an error.
type bagItemsValueBesideDefect struct {
	bagOverEmbedded
}

// Items is the override of the planted defect.
func (d *bagItemsValueBesideDefect) Items(p0 context.Context) (r0 []Value, r1 error) {
	r0, r1 = d.bagOverEmbedded.Items(p0)
	if r1 != nil {
		r0 = []Value{{Key: "Key-b"}}
	}
	return r0, r1
}

// bagItemsValueBeside returns the planted defect of a Bag whose Items
// returns a value beside an error over a fresh subject.
func bagItemsValueBeside(h *BagHandle) Bag {
	return &bagItemsValueBesideDefect{bagOverEmbedded: h.Fresh()}
}

// bagItemsDrainsTwiceDefect is the planted defect of a Bag whose Items
// returns one element twice.
type bagItemsDrainsTwiceDefect struct {
	bagOverEmbedded
}

// Items is the override of the planted defect.
func (d *bagItemsDrainsTwiceDefect) Items(p0 context.Context) (r0 []Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.bagOverEmbedded.Items(p0)
	if len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0, r1
}

// bagItemsDrainsTwice returns the planted defect of a Bag whose Items
// returns one element twice over a fresh subject.
func bagItemsDrainsTwice(h *BagHandle) Bag {
	return &bagItemsDrainsTwiceDefect{bagOverEmbedded: h.Fresh()}
}

// bagItemsReordersDrainDefect is the planted defect of a Bag whose Items
// reorders the elements between calls.
type bagItemsReordersDrainDefect struct {
	bagOverEmbedded
	calls int
}

// Items is the override of the planted defect.
func (d *bagItemsReordersDrainDefect) Items(p0 context.Context) (r0 []Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.bagOverEmbedded.Items(p0)
	d.calls++
	if d.calls > 1 {
		slices.Reverse(r0)
	}
	return r0, r1
}

// bagItemsReordersDrain returns the planted defect of a Bag whose Items
// reorders the elements between calls over a fresh subject.
func bagItemsReordersDrain(h *BagHandle) Bag {
	return &bagItemsReordersDrainDefect{bagOverEmbedded: h.Fresh()}
}

// The derived checks of Updater:
//
// Updater.Put/smoke: Put returns without a panic for derived inputs.
//
// Updater.Put/cancel: Put returns an error for a cancelled context.
//
// Updater.Put/updater: a second Put for one key replaces the first.
//
// Updater.Get/smoke: Get returns without a panic for derived inputs.
//
// Updater.Get/cancel: Get returns an error for a cancelled context.
//
// Updater.Get/zero-on-error: Get returns the zero value beside an error.

// The derived checks of Updater that witness their laws at the shallow
// end:
//
// Updater.Put/smoke
//
// Updater.Get/smoke

// The unclaimed classifications of Updater:
//
// The shape:writer of Updater.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Updater writes.
type (
	UpdaterHandle = conform.Handle[Updater]
	UpdaterCheck  = conform.Check[Updater]
	UpdaterOption = conform.Option[Updater]
)

// UpdaterChecks contains the name of each callable of Updater, with the
// identities of its checks.
var UpdaterChecks = struct {
	Put updaterPut
	Get updaterGet
}{
	Put: updaterPut{"Updater.Put"},
	Get: updaterGet{"Updater.Get"},
}

// updaterPut is the name of Updater.Put, with the identities of its
// checks.
type updaterPut struct {
	conform.Callable[Updater]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c updaterPut) Smoke() conform.ID[Updater] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c updaterPut) Cancel() conform.ID[Updater] {
	return c.Check("cancel")
}

// Updater returns the identity of the check that a second Put for one key
// replaces the first.
func (c updaterPut) Updater() conform.ID[Updater] {
	return c.Check("updater")
}

// updaterGet is the name of Updater.Get, with the identities of its
// checks.
type updaterGet struct {
	conform.Callable[Updater]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c updaterGet) Smoke() conform.ID[Updater] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c updaterGet) Cancel() conform.ID[Updater] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c updaterGet) ZeroOnError() conform.ID[Updater] {
	return c.Check("zero-on-error")
}

// UpdaterSuite returns the checks of Updater, and the calls and the
// budgets of its callables.
func UpdaterSuite() conform.Suite[Updater] {
	return conform.Suite[Updater]{
		Declaration: "Updater",
		Checks: []conform.Check[Updater]{
			{
				ID:          UpdaterChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *UpdaterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: updaterPutPanics,
			},
			{
				ID:          UpdaterChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *UpdaterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: updaterPutIgnoresCancel,
			},
			{
				ID:          UpdaterChecks.Put.Updater(),
				Statement:   "a second Put for one key replaces the first",
				Claim:       spec.Updater,
				Probe:       spec.SecondWriteReplaces,
				Reason:      spec.UnreplacedWrite,
				Law:         spec.SecondWriteReplacesTheFirst,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *UpdaterHandle) {
					probe.SecondWriteReplaces(h, func() func(ctx context.Context) (Value, error) {
						calls := 0
						return func(ctx context.Context) (Value, error) {
							calls++
							w := Value{Key: "Key-a"}
							w.Body = string("Body-a")
							if calls > 1 {
								w.Body = string("Body-b")
							}
							err := h.Subject().Put(ctx, w)
							return w, err
						}
					}(), func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					})
				},
				Defect: updaterPutStoresNothing,
			},
			{
				ID:          UpdaterChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *UpdaterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: updaterGetPanics,
			},
			{
				ID:          UpdaterChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *UpdaterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: updaterGetIgnoresCancel,
			},
			{
				ID:          UpdaterChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *UpdaterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: updaterGetValueBeside,
			},
		},
		Calls: []conform.Call[Updater]{
			{
				Calls:    updaterPutCalls,
				Over:     newUpdaterOver,
				Callable: UpdaterChecks.Put.Callable,
			},
			{
				Calls:    updaterGetCalls,
				Over:     newUpdaterOver,
				Callable: UpdaterChecks.Get.Callable,
			},
		},
	}
}

// RunUpdaterSuite applies the checks of Updater, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunUpdaterSuite[T Updater](t *testing.T, newSubject func() T, opts ...UpdaterOption) {
	t.Helper()
	RunUpdaterSuiteFromHandle(t, func(*UpdaterHandle) T { return newSubject() }, opts...)
}

// RunUpdaterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunUpdaterSuiteFromHandle[T Updater](
	t *testing.T,
	newSubject func(h *UpdaterHandle) T,
	opts ...UpdaterOption,
) {
	t.Helper()
	conform.Run(t, UpdaterSuite(), updaterSubjects(newSubject), opts...)
}

// FuzzUpdaterSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzUpdaterSuite[T Updater](
	f *testing.F,
	callable conform.Callable[Updater],
	newSubject func(h *UpdaterHandle) T,
	opts ...UpdaterOption,
) {
	f.Helper()
	conform.Fuzz(f, UpdaterSuite(), callable, updaterSubjects(newSubject), opts...)
}

// RunUpdaterBudgets applies the allocation ceilings of Updater and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunUpdaterBudgets[T Updater](
	t *testing.T,
	newSubject func(h *UpdaterHandle) T,
	opts ...UpdaterOption,
) {
	t.Helper()
	conform.RunBudgets(t, UpdaterSuite(), updaterSubjects(newSubject), opts...)
}

// BenchUpdaterSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchUpdaterSuite[T Updater](
	b conform.B,
	callable conform.Callable[Updater],
	newSubject func(h *UpdaterHandle) T,
	opts ...UpdaterOption,
) {
	b.Helper()
	conform.Bench(b, UpdaterSuite(), callable, updaterSubjects(newSubject), opts...)
}

// updaterSubjects returns the factory of the subjects of Updater that
// calls newSubject with the handle of each check.
func updaterSubjects[T Updater](
	newSubject func(h *UpdaterHandle) T,
) func(h *UpdaterHandle) Updater {
	return func(h *UpdaterHandle) Updater { return newSubject(h) }
}

// updaterPutCalls generates the calls of Updater.Put with the inputs that
// each case draws.
var updaterPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Updater) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Updater) error {
		return subject.Put(ctx, v)
	}
})

// updaterGetCalls generates the calls of Updater.Get with the inputs that
// each case draws.
var updaterGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Updater) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Updater) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// updaterOverEmbedded is Updater under the name of the field that
// updaterOver embeds.
type updaterOverEmbedded = Updater

// updaterOver is a subject of Updater whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type updaterOver struct {
	updaterOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o updaterOver) Put(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o updaterOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newUpdaterOver returns the planted defect of a ceiling of Updater, whose
// callables call spend.
func newUpdaterOver(_ *UpdaterHandle, spend func()) Updater {
	return updaterOver{spend: spend}
}

// updaterPutPanicsDefect is the planted defect of a Updater whose Put
// panics.
type updaterPutPanicsDefect struct {
	updaterOverEmbedded
}

// Put is the override of the planted defect.
func (*updaterPutPanicsDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Updater whose Put panics")
}

// updaterPutPanics returns the planted defect of a Updater whose Put
// panics over a fresh subject.
func updaterPutPanics(h *UpdaterHandle) Updater {
	return &updaterPutPanicsDefect{updaterOverEmbedded: h.Fresh()}
}

// updaterPutIgnoresCancelDefect is the planted defect of a Updater whose
// Put ignores the context.
type updaterPutIgnoresCancelDefect struct {
	updaterOverEmbedded
}

// Put is the override of the planted defect.
func (d *updaterPutIgnoresCancelDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.updaterOverEmbedded.Put(p0, p1)
}

// updaterPutIgnoresCancel returns the planted defect of a Updater whose
// Put ignores the context over a fresh subject.
func updaterPutIgnoresCancel(h *UpdaterHandle) Updater {
	return &updaterPutIgnoresCancelDefect{updaterOverEmbedded: h.Fresh()}
}

// updaterPutStoresNothingDefect is the planted defect of a Updater whose
// Put keeps the first write.
type updaterPutStoresNothingDefect struct {
	updaterOverEmbedded
}

// Put is the override of the planted defect.
func (*updaterPutStoresNothingDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// updaterPutStoresNothing returns the planted defect of a Updater whose
// Put keeps the first write over a fresh subject.
func updaterPutStoresNothing(h *UpdaterHandle) Updater {
	return &updaterPutStoresNothingDefect{updaterOverEmbedded: h.Fresh()}
}

// updaterGetPanicsDefect is the planted defect of a Updater whose Get
// panics.
type updaterGetPanicsDefect struct {
	updaterOverEmbedded
}

// Get is the override of the planted defect.
func (*updaterGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Updater whose Get panics")
}

// updaterGetPanics returns the planted defect of a Updater whose Get
// panics over a fresh subject.
func updaterGetPanics(h *UpdaterHandle) Updater {
	return &updaterGetPanicsDefect{updaterOverEmbedded: h.Fresh()}
}

// updaterGetIgnoresCancelDefect is the planted defect of a Updater whose
// Get ignores the context.
type updaterGetIgnoresCancelDefect struct {
	updaterOverEmbedded
}

// Get is the override of the planted defect.
func (d *updaterGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.updaterOverEmbedded.Get(p0, p1)
}

// updaterGetIgnoresCancel returns the planted defect of a Updater whose
// Get ignores the context over a fresh subject.
func updaterGetIgnoresCancel(h *UpdaterHandle) Updater {
	return &updaterGetIgnoresCancelDefect{updaterOverEmbedded: h.Fresh()}
}

// updaterGetValueBesideDefect is the planted defect of a Updater whose Get
// returns a value beside an error.
type updaterGetValueBesideDefect struct {
	updaterOverEmbedded
}

// Get is the override of the planted defect.
func (d *updaterGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.updaterOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// updaterGetValueBeside returns the planted defect of a Updater whose Get
// returns a value beside an error over a fresh subject.
func updaterGetValueBeside(h *UpdaterHandle) Updater {
	return &updaterGetValueBesideDefect{updaterOverEmbedded: h.Fresh()}
}

// The derived checks of KeyedUpdater:
//
// KeyedUpdater.Put/smoke: Put returns without a panic for derived inputs.
//
// KeyedUpdater.Put/cancel: Put returns an error for a cancelled context.
//
// KeyedUpdater.Put/updater: a second Put for one key replaces the first.
//
// KeyedUpdater.Get/smoke: Get returns without a panic for derived inputs.
//
// KeyedUpdater.Get/cancel: Get returns an error for a cancelled context.
//
// KeyedUpdater.Get/zero-on-error: Get returns the zero value beside an
// error.

// The derived checks of KeyedUpdater that witness their laws at the
// shallow end:
//
// KeyedUpdater.Put/smoke
//
// KeyedUpdater.Get/smoke

// The unclaimed classifications of KeyedUpdater:
//
// The shape:composite-writer of KeyedUpdater.Put has the stance silent. A
// write is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of KeyedUpdater writes.
type (
	KeyedUpdaterHandle = conform.Handle[KeyedUpdater]
	KeyedUpdaterCheck  = conform.Check[KeyedUpdater]
	KeyedUpdaterOption = conform.Option[KeyedUpdater]
)

// KeyedUpdaterChecks contains the name of each callable of KeyedUpdater,
// with the identities of its checks.
var KeyedUpdaterChecks = struct {
	Put keyedUpdaterPut
	Get keyedUpdaterGet
}{
	Put: keyedUpdaterPut{"KeyedUpdater.Put"},
	Get: keyedUpdaterGet{"KeyedUpdater.Get"},
}

// keyedUpdaterPut is the name of KeyedUpdater.Put, with the identities of
// its checks.
type keyedUpdaterPut struct {
	conform.Callable[KeyedUpdater]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c keyedUpdaterPut) Smoke() conform.ID[KeyedUpdater] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c keyedUpdaterPut) Cancel() conform.ID[KeyedUpdater] {
	return c.Check("cancel")
}

// Updater returns the identity of the check that a second Put for one key
// replaces the first.
func (c keyedUpdaterPut) Updater() conform.ID[KeyedUpdater] {
	return c.Check("updater")
}

// keyedUpdaterGet is the name of KeyedUpdater.Get, with the identities of
// its checks.
type keyedUpdaterGet struct {
	conform.Callable[KeyedUpdater]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c keyedUpdaterGet) Smoke() conform.ID[KeyedUpdater] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c keyedUpdaterGet) Cancel() conform.ID[KeyedUpdater] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c keyedUpdaterGet) ZeroOnError() conform.ID[KeyedUpdater] {
	return c.Check("zero-on-error")
}

// KeyedUpdaterSuite returns the checks of KeyedUpdater, and the calls and
// the budgets of its callables.
func KeyedUpdaterSuite() conform.Suite[KeyedUpdater] {
	return conform.Suite[KeyedUpdater]{
		Declaration: "KeyedUpdater",
		Checks: []conform.Check[KeyedUpdater]{
			{
				ID:          KeyedUpdaterChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeyedUpdaterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: keyedUpdaterPutPanics,
			},
			{
				ID:          KeyedUpdaterChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedUpdaterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: keyedUpdaterPutIgnoresCancel,
			},
			{
				ID:          KeyedUpdaterChecks.Put.Updater(),
				Statement:   "a second Put for one key replaces the first",
				Claim:       spec.Updater,
				Probe:       spec.SecondWriteReplaces,
				Reason:      spec.UnreplacedWrite,
				Law:         spec.SecondWriteReplacesTheFirst,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedUpdaterHandle) {
					probe.SecondWriteReplaces(h, func() func(ctx context.Context) (Value, error) {
						calls := 0
						return func(ctx context.Context) (Value, error) {
							calls++
							w := Value{Key: "Key-a"}
							w.Key = string("Key-a")
							if calls > 1 {
								w.Key = string("Key-b")
							}
							err := h.Subject().Put(ctx, string("key-a"), w)
							return w, err
						}
					}(), func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					})
				},
				Defect: keyedUpdaterPutStoresNothing,
			},
			{
				ID:          KeyedUpdaterChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeyedUpdaterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: keyedUpdaterGetPanics,
			},
			{
				ID:          KeyedUpdaterChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedUpdaterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: keyedUpdaterGetIgnoresCancel,
			},
			{
				ID:          KeyedUpdaterChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeyedUpdaterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: keyedUpdaterGetValueBeside,
			},
		},
		Calls: []conform.Call[KeyedUpdater]{
			{
				Calls:    keyedUpdaterPutCalls,
				Over:     newKeyedUpdaterOver,
				Callable: KeyedUpdaterChecks.Put.Callable,
			},
			{
				Calls:    keyedUpdaterGetCalls,
				Over:     newKeyedUpdaterOver,
				Callable: KeyedUpdaterChecks.Get.Callable,
			},
		},
	}
}

// RunKeyedUpdaterSuite applies the checks of KeyedUpdater, and the checks
// that opts add, to a fresh subject of newSubject for each check.
func RunKeyedUpdaterSuite[T KeyedUpdater](
	t *testing.T,
	newSubject func() T,
	opts ...KeyedUpdaterOption,
) {
	t.Helper()
	RunKeyedUpdaterSuiteFromHandle(t, func(*KeyedUpdaterHandle) T { return newSubject() }, opts...)
}

// RunKeyedUpdaterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunKeyedUpdaterSuiteFromHandle[T KeyedUpdater](
	t *testing.T,
	newSubject func(h *KeyedUpdaterHandle) T,
	opts ...KeyedUpdaterOption,
) {
	t.Helper()
	conform.Run(t, KeyedUpdaterSuite(), keyedUpdaterSubjects(newSubject), opts...)
}

// FuzzKeyedUpdaterSuite registers on f the fuzz target of the property
// checks of callable, on a subject that newSubject builds from the handle
// of each case.
func FuzzKeyedUpdaterSuite[T KeyedUpdater](
	f *testing.F,
	callable conform.Callable[KeyedUpdater],
	newSubject func(h *KeyedUpdaterHandle) T,
	opts ...KeyedUpdaterOption,
) {
	f.Helper()
	conform.Fuzz(f, KeyedUpdaterSuite(), callable, keyedUpdaterSubjects(newSubject), opts...)
}

// RunKeyedUpdaterBudgets applies the allocation ceilings of KeyedUpdater
// and its checks that count allocations, and the ceilings and the checks
// that opts add, to a fresh subject of newSubject for each case.
func RunKeyedUpdaterBudgets[T KeyedUpdater](
	t *testing.T,
	newSubject func(h *KeyedUpdaterHandle) T,
	opts ...KeyedUpdaterOption,
) {
	t.Helper()
	conform.RunBudgets(t, KeyedUpdaterSuite(), keyedUpdaterSubjects(newSubject), opts...)
}

// BenchKeyedUpdaterSuite measures callable on a subject of newSubject,
// under the ceilings of its budget.
func BenchKeyedUpdaterSuite[T KeyedUpdater](
	b conform.B,
	callable conform.Callable[KeyedUpdater],
	newSubject func(h *KeyedUpdaterHandle) T,
	opts ...KeyedUpdaterOption,
) {
	b.Helper()
	conform.Bench(b, KeyedUpdaterSuite(), callable, keyedUpdaterSubjects(newSubject), opts...)
}

// keyedUpdaterSubjects returns the factory of the subjects of KeyedUpdater
// that calls newSubject with the handle of each check.
func keyedUpdaterSubjects[T KeyedUpdater](
	newSubject func(h *KeyedUpdaterHandle) T,
) func(h *KeyedUpdaterHandle) KeyedUpdater {
	return func(h *KeyedUpdaterHandle) KeyedUpdater { return newSubject(h) }
}

// keyedUpdaterPutCalls generates the calls of KeyedUpdater.Put with the
// inputs that each case draws.
var keyedUpdaterPutCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, KeyedUpdater) error {
		key := c.Draw(prop.Of[string](), "key")
		v := c.Draw(prop.Of[Value](), "v")
		return func(ctx context.Context, subject KeyedUpdater) error {
			return subject.Put(ctx, key, v)
		}
	},
)

// keyedUpdaterGetCalls generates the calls of KeyedUpdater.Get with the
// inputs that each case draws.
var keyedUpdaterGetCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, KeyedUpdater) error {
		key := c.Draw(prop.Of[string](), "key")
		return func(ctx context.Context, subject KeyedUpdater) error {
			_, err := subject.Get(ctx, key)
			return err
		}
	},
)

// keyedUpdaterOverEmbedded is KeyedUpdater under the name of the field
// that keyedUpdaterOver embeds.
type keyedUpdaterOverEmbedded = KeyedUpdater

// keyedUpdaterOver is a subject of KeyedUpdater whose callables call spend
// and return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type keyedUpdaterOver struct {
	keyedUpdaterOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o keyedUpdaterOver) Put(context.Context, string, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o keyedUpdaterOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newKeyedUpdaterOver returns the planted defect of a ceiling of
// KeyedUpdater, whose callables call spend.
func newKeyedUpdaterOver(_ *KeyedUpdaterHandle, spend func()) KeyedUpdater {
	return keyedUpdaterOver{spend: spend}
}

// keyedUpdaterPutPanicsDefect is the planted defect of a KeyedUpdater
// whose Put panics.
type keyedUpdaterPutPanicsDefect struct {
	keyedUpdaterOverEmbedded
}

// Put is the override of the planted defect.
func (*keyedUpdaterPutPanicsDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	panic("the planted defect of a KeyedUpdater whose Put panics")
}

// keyedUpdaterPutPanics returns the planted defect of a KeyedUpdater whose
// Put panics over a fresh subject.
func keyedUpdaterPutPanics(h *KeyedUpdaterHandle) KeyedUpdater {
	return &keyedUpdaterPutPanicsDefect{keyedUpdaterOverEmbedded: h.Fresh()}
}

// keyedUpdaterPutIgnoresCancelDefect is the planted defect of a
// KeyedUpdater whose Put ignores the context.
type keyedUpdaterPutIgnoresCancelDefect struct {
	keyedUpdaterOverEmbedded
}

// Put is the override of the planted defect.
func (d *keyedUpdaterPutIgnoresCancelDefect) Put(
	p0 context.Context,
	p1 string,
	p2 Value,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.keyedUpdaterOverEmbedded.Put(p0, p1, p2)
}

// keyedUpdaterPutIgnoresCancel returns the planted defect of a
// KeyedUpdater whose Put ignores the context over a fresh subject.
func keyedUpdaterPutIgnoresCancel(h *KeyedUpdaterHandle) KeyedUpdater {
	return &keyedUpdaterPutIgnoresCancelDefect{keyedUpdaterOverEmbedded: h.Fresh()}
}

// keyedUpdaterPutStoresNothingDefect is the planted defect of a
// KeyedUpdater whose Put keeps the first write.
type keyedUpdaterPutStoresNothingDefect struct {
	keyedUpdaterOverEmbedded
}

// Put is the override of the planted defect.
func (*keyedUpdaterPutStoresNothingDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// keyedUpdaterPutStoresNothing returns the planted defect of a
// KeyedUpdater whose Put keeps the first write over a fresh subject.
func keyedUpdaterPutStoresNothing(h *KeyedUpdaterHandle) KeyedUpdater {
	return &keyedUpdaterPutStoresNothingDefect{keyedUpdaterOverEmbedded: h.Fresh()}
}

// keyedUpdaterGetPanicsDefect is the planted defect of a KeyedUpdater
// whose Get panics.
type keyedUpdaterGetPanicsDefect struct {
	keyedUpdaterOverEmbedded
}

// Get is the override of the planted defect.
func (*keyedUpdaterGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a KeyedUpdater whose Get panics")
}

// keyedUpdaterGetPanics returns the planted defect of a KeyedUpdater whose
// Get panics over a fresh subject.
func keyedUpdaterGetPanics(h *KeyedUpdaterHandle) KeyedUpdater {
	return &keyedUpdaterGetPanicsDefect{keyedUpdaterOverEmbedded: h.Fresh()}
}

// keyedUpdaterGetIgnoresCancelDefect is the planted defect of a
// KeyedUpdater whose Get ignores the context.
type keyedUpdaterGetIgnoresCancelDefect struct {
	keyedUpdaterOverEmbedded
}

// Get is the override of the planted defect.
func (d *keyedUpdaterGetIgnoresCancelDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.keyedUpdaterOverEmbedded.Get(p0, p1)
}

// keyedUpdaterGetIgnoresCancel returns the planted defect of a
// KeyedUpdater whose Get ignores the context over a fresh subject.
func keyedUpdaterGetIgnoresCancel(h *KeyedUpdaterHandle) KeyedUpdater {
	return &keyedUpdaterGetIgnoresCancelDefect{keyedUpdaterOverEmbedded: h.Fresh()}
}

// keyedUpdaterGetValueBesideDefect is the planted defect of a KeyedUpdater
// whose Get returns a value beside an error.
type keyedUpdaterGetValueBesideDefect struct {
	keyedUpdaterOverEmbedded
}

// Get is the override of the planted defect.
func (d *keyedUpdaterGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.keyedUpdaterOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// keyedUpdaterGetValueBeside returns the planted defect of a KeyedUpdater
// whose Get returns a value beside an error over a fresh subject.
func keyedUpdaterGetValueBeside(h *KeyedUpdaterHandle) KeyedUpdater {
	return &keyedUpdaterGetValueBesideDefect{keyedUpdaterOverEmbedded: h.Fresh()}
}

// The derived checks of Codec:
//
// Codec.Encode/smoke: Encode returns without a panic for derived inputs.
//
// Codec.Encode/cancel: Encode returns an error for a cancelled context.
//
// Codec.Encode/zero-on-error: Encode returns the zero value beside an
// error.
//
// Codec.Encode/codec: Decode returns what went into Encode.
//
// Codec.Decode/smoke: Decode returns without a panic for derived inputs.
//
// Codec.Decode/cancel: Decode returns an error for a cancelled context.
//
// Codec.Decode/zero-on-error: Decode returns the zero value beside an
// error.

// The derived checks of Codec that witness their laws at the shallow end:
//
// Codec.Encode/smoke
//
// Codec.Encode/codec
//
// Codec.Decode/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Codec writes.
type (
	CodecHandle = conform.Handle[Codec]
	CodecCheck  = conform.Check[Codec]
	CodecOption = conform.Option[Codec]
)

// CodecChecks contains the name of each callable of Codec, with the
// identities of its checks.
var CodecChecks = struct {
	Encode codecEncode
	Decode codecDecode
}{
	Encode: codecEncode{"Codec.Encode"},
	Decode: codecDecode{"Codec.Decode"},
}

// codecEncode is the name of Codec.Encode, with the identities of its
// checks.
type codecEncode struct {
	conform.Callable[Codec]
}

// Smoke returns the identity of the check that Encode returns without a
// panic for derived inputs.
func (c codecEncode) Smoke() conform.ID[Codec] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Encode returns an error
// for a cancelled context.
func (c codecEncode) Cancel() conform.ID[Codec] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Encode returns the
// zero value beside an error.
func (c codecEncode) ZeroOnError() conform.ID[Codec] {
	return c.Check("zero-on-error")
}

// Codec returns the identity of the check that Decode returns what went
// into Encode.
func (c codecEncode) Codec() conform.ID[Codec] {
	return c.Check("codec")
}

// codecDecode is the name of Codec.Decode, with the identities of its
// checks.
type codecDecode struct {
	conform.Callable[Codec]
}

// Smoke returns the identity of the check that Decode returns without a
// panic for derived inputs.
func (c codecDecode) Smoke() conform.ID[Codec] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Decode returns an error
// for a cancelled context.
func (c codecDecode) Cancel() conform.ID[Codec] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Decode returns the
// zero value beside an error.
func (c codecDecode) ZeroOnError() conform.ID[Codec] {
	return c.Check("zero-on-error")
}

// CodecSuite returns the checks of Codec, and the calls and the budgets of
// its callables.
func CodecSuite() conform.Suite[Codec] {
	return conform.Suite[Codec]{
		Declaration: "Codec",
		Checks: []conform.Check[Codec]{
			{
				ID:          CodecChecks.Encode.Smoke(),
				Statement:   "Encode returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CodecHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Encode(ctx, string("in-a"))
						return err
					})
				},
				Defect: codecEncodePanics,
			},
			{
				ID:          CodecChecks.Encode.Cancel(),
				Statement:   "Encode returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CodecHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Encode(ctx, string("in-a"))
						return err
					})
				},
				Defect: codecEncodeIgnoresCancel,
			},
			{
				ID:          CodecChecks.Encode.ZeroOnError(),
				Statement:   "Encode returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CodecHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Encode(ctx, string("in-a"))
					}, string(""))
				},
				Defect: codecEncodeValueBeside,
			},
			{
				ID:          CodecChecks.Encode.Codec(),
				Statement:   "Decode returns what went into Encode",
				Claim:       spec.Codec,
				Probe:       spec.RoundTrips,
				Reason:      spec.BrokenRoundTrip,
				Law:         spec.EncodingSurvivesItsRoundTrip,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CodecHandle) {
					probe.RoundTrips(
						h,
						func(ctx context.Context) (probe.Pair[string, string], error) {
							in := string("in-a")
							r0, err := h.Subject().Encode(ctx, in)
							return probe.Pair[string, string]{First: in, Second: r0}, err
						},
						func(ctx context.Context, v string) (string, error) {
							return h.Subject().Decode(ctx, v)
						},
					)
				},
				Defect: codecEncodeDrifts,
			},
			{
				ID:          CodecChecks.Decode.Smoke(),
				Statement:   "Decode returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CodecHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Decode(ctx, string("in-a"))
						return err
					})
				},
				Defect: codecDecodePanics,
			},
			{
				ID:          CodecChecks.Decode.Cancel(),
				Statement:   "Decode returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CodecHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Decode(ctx, string("in-a"))
						return err
					})
				},
				Defect: codecDecodeIgnoresCancel,
			},
			{
				ID:          CodecChecks.Decode.ZeroOnError(),
				Statement:   "Decode returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CodecHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Decode(ctx, string("in-a"))
					}, string(""))
				},
				Defect: codecDecodeValueBeside,
			},
		},
		Calls: []conform.Call[Codec]{
			{
				Calls:    codecEncodeCalls,
				Over:     newCodecOver,
				Callable: CodecChecks.Encode.Callable,
			},
			{
				Calls:    codecDecodeCalls,
				Over:     newCodecOver,
				Callable: CodecChecks.Decode.Callable,
			},
		},
	}
}

// RunCodecSuite applies the checks of Codec, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunCodecSuite[T Codec](t *testing.T, newSubject func() T, opts ...CodecOption) {
	t.Helper()
	RunCodecSuiteFromHandle(t, func(*CodecHandle) T { return newSubject() }, opts...)
}

// RunCodecSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCodecSuiteFromHandle[T Codec](
	t *testing.T,
	newSubject func(h *CodecHandle) T,
	opts ...CodecOption,
) {
	t.Helper()
	conform.Run(t, CodecSuite(), codecSubjects(newSubject), opts...)
}

// FuzzCodecSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCodecSuite[T Codec](
	f *testing.F,
	callable conform.Callable[Codec],
	newSubject func(h *CodecHandle) T,
	opts ...CodecOption,
) {
	f.Helper()
	conform.Fuzz(f, CodecSuite(), callable, codecSubjects(newSubject), opts...)
}

// RunCodecBudgets applies the allocation ceilings of Codec and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunCodecBudgets[T Codec](
	t *testing.T,
	newSubject func(h *CodecHandle) T,
	opts ...CodecOption,
) {
	t.Helper()
	conform.RunBudgets(t, CodecSuite(), codecSubjects(newSubject), opts...)
}

// BenchCodecSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchCodecSuite[T Codec](
	b conform.B,
	callable conform.Callable[Codec],
	newSubject func(h *CodecHandle) T,
	opts ...CodecOption,
) {
	b.Helper()
	conform.Bench(b, CodecSuite(), callable, codecSubjects(newSubject), opts...)
}

// codecSubjects returns the factory of the subjects of Codec that calls
// newSubject with the handle of each check.
func codecSubjects[T Codec](newSubject func(h *CodecHandle) T) func(h *CodecHandle) Codec {
	return func(h *CodecHandle) Codec { return newSubject(h) }
}

// codecEncodeCalls generates the calls of Codec.Encode with the inputs
// that each case draws.
var codecEncodeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Codec) error {
	in := c.Draw(prop.Of[string](), "in")
	return func(ctx context.Context, subject Codec) error {
		_, err := subject.Encode(ctx, in)
		return err
	}
})

// codecDecodeCalls generates the calls of Codec.Decode with the inputs
// that each case draws.
var codecDecodeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Codec) error {
	in := c.Draw(prop.Of[string](), "in")
	return func(ctx context.Context, subject Codec) error {
		_, err := subject.Decode(ctx, in)
		return err
	}
})

// codecOverEmbedded is Codec under the name of the field that codecOver
// embeds.
type codecOverEmbedded = Codec

// codecOver is a subject of Codec whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type codecOver struct {
	codecOverEmbedded
	spend func()
}

// Encode calls spend and returns its zero results.
func (o codecOver) Encode(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// Decode calls spend and returns its zero results.
func (o codecOver) Decode(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newCodecOver returns the planted defect of a ceiling of Codec, whose
// callables call spend.
func newCodecOver(_ *CodecHandle, spend func()) Codec {
	return codecOver{spend: spend}
}

// codecEncodePanicsDefect is the planted defect of a Codec whose Encode
// panics.
type codecEncodePanicsDefect struct {
	codecOverEmbedded
}

// Encode is the override of the planted defect.
func (*codecEncodePanicsDefect) Encode(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Codec whose Encode panics")
}

// codecEncodePanics returns the planted defect of a Codec whose Encode
// panics over a fresh subject.
func codecEncodePanics(h *CodecHandle) Codec {
	return &codecEncodePanicsDefect{codecOverEmbedded: h.Fresh()}
}

// codecEncodeIgnoresCancelDefect is the planted defect of a Codec whose
// Encode ignores the context.
type codecEncodeIgnoresCancelDefect struct {
	codecOverEmbedded
}

// Encode is the override of the planted defect.
func (d *codecEncodeIgnoresCancelDefect) Encode(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.codecOverEmbedded.Encode(p0, p1)
}

// codecEncodeIgnoresCancel returns the planted defect of a Codec whose
// Encode ignores the context over a fresh subject.
func codecEncodeIgnoresCancel(h *CodecHandle) Codec {
	return &codecEncodeIgnoresCancelDefect{codecOverEmbedded: h.Fresh()}
}

// codecEncodeValueBesideDefect is the planted defect of a Codec whose
// Encode returns a value beside an error.
type codecEncodeValueBesideDefect struct {
	codecOverEmbedded
}

// Encode is the override of the planted defect.
func (d *codecEncodeValueBesideDefect) Encode(p0 context.Context, p1 string) (r0 string, r1 error) {
	r0, r1 = d.codecOverEmbedded.Encode(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// codecEncodeValueBeside returns the planted defect of a Codec whose
// Encode returns a value beside an error over a fresh subject.
func codecEncodeValueBeside(h *CodecHandle) Codec {
	return &codecEncodeValueBesideDefect{codecOverEmbedded: h.Fresh()}
}

// codecEncodeDriftsDefect is the planted defect of a Codec whose Encode
// returns another encoding on every call.
type codecEncodeDriftsDefect struct {
	codecOverEmbedded
	calls int
}

// Encode is the override of the planted defect.
func (d *codecEncodeDriftsDefect) Encode(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls > 1 {
		r0 = string("sample-b")
	}
	return r0, r1
}

// codecEncodeDrifts returns the planted defect of a Codec whose Encode
// returns another encoding on every call over a fresh subject.
func codecEncodeDrifts(h *CodecHandle) Codec {
	return &codecEncodeDriftsDefect{codecOverEmbedded: h.Fresh()}
}

// codecDecodePanicsDefect is the planted defect of a Codec whose Decode
// panics.
type codecDecodePanicsDefect struct {
	codecOverEmbedded
}

// Decode is the override of the planted defect.
func (*codecDecodePanicsDefect) Decode(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Codec whose Decode panics")
}

// codecDecodePanics returns the planted defect of a Codec whose Decode
// panics over a fresh subject.
func codecDecodePanics(h *CodecHandle) Codec {
	return &codecDecodePanicsDefect{codecOverEmbedded: h.Fresh()}
}

// codecDecodeIgnoresCancelDefect is the planted defect of a Codec whose
// Decode ignores the context.
type codecDecodeIgnoresCancelDefect struct {
	codecOverEmbedded
}

// Decode is the override of the planted defect.
func (d *codecDecodeIgnoresCancelDefect) Decode(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.codecOverEmbedded.Decode(p0, p1)
}

// codecDecodeIgnoresCancel returns the planted defect of a Codec whose
// Decode ignores the context over a fresh subject.
func codecDecodeIgnoresCancel(h *CodecHandle) Codec {
	return &codecDecodeIgnoresCancelDefect{codecOverEmbedded: h.Fresh()}
}

// codecDecodeValueBesideDefect is the planted defect of a Codec whose
// Decode returns a value beside an error.
type codecDecodeValueBesideDefect struct {
	codecOverEmbedded
}

// Decode is the override of the planted defect.
func (d *codecDecodeValueBesideDefect) Decode(p0 context.Context, p1 string) (r0 string, r1 error) {
	r0, r1 = d.codecOverEmbedded.Decode(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// codecDecodeValueBeside returns the planted defect of a Codec whose
// Decode returns a value beside an error over a fresh subject.
func codecDecodeValueBeside(h *CodecHandle) Codec {
	return &codecDecodeValueBesideDefect{codecOverEmbedded: h.Fresh()}
}

// The derived checks of Pager:
//
// Pager.Fetch/smoke: Fetch returns without a panic for derived inputs.
//
// Pager.Fetch/cancel: Fetch returns an error for a cancelled context.
//
// Pager.Fetch/zero-on-error: Fetch returns the zero value beside an error.
//
// Pager.Fetch/pagination: Fetch resumes from the cursor of its own page.

// The derived checks of Pager that witness their laws at the shallow end:
//
// Pager.Fetch/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Pager writes.
type (
	PagerHandle = conform.Handle[Pager]
	PagerCheck  = conform.Check[Pager]
	PagerOption = conform.Option[Pager]
)

// PagerChecks contains the name of each callable of Pager, with the
// identities of its checks.
var PagerChecks = struct {
	Fetch pagerFetch
}{
	Fetch: pagerFetch{"Pager.Fetch"},
}

// pagerFetch is the name of Pager.Fetch, with the identities of its
// checks.
type pagerFetch struct {
	conform.Callable[Pager]
}

// Smoke returns the identity of the check that Fetch returns without a
// panic for derived inputs.
func (c pagerFetch) Smoke() conform.ID[Pager] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Fetch returns an error for
// a cancelled context.
func (c pagerFetch) Cancel() conform.ID[Pager] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Fetch returns the
// zero value beside an error.
func (c pagerFetch) ZeroOnError() conform.ID[Pager] {
	return c.Check("zero-on-error")
}

// Pagination returns the identity of the check that Fetch resumes from the
// cursor of its own page.
func (c pagerFetch) Pagination() conform.ID[Pager] {
	return c.Check("pagination")
}

// PagerSuite returns the checks of Pager, and the calls and the budgets of
// its callables.
func PagerSuite() conform.Suite[Pager] {
	return conform.Suite[Pager]{
		Declaration: "Pager",
		Checks: []conform.Check[Pager]{
			{
				ID:          PagerChecks.Fetch.Smoke(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PagerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, int(42))
						return err
					})
				},
				Defect: pagerFetchPanics,
			},
			{
				ID:          PagerChecks.Fetch.Cancel(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PagerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, int(42))
						return err
					})
				},
				Defect: pagerFetchIgnoresCancel,
			},
			{
				ID:          PagerChecks.Fetch.ZeroOnError(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PagerHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (*Page, error) {
						return h.Subject().Fetch(ctx, int(42))
					}, (*Page)(nil))
				},
				Defect: pagerFetchValueBeside,
			},
			{
				ID:          PagerChecks.Fetch.Pagination(),
				Statement:   "Fetch resumes from the cursor of its own page",
				Claim:       spec.Pagination,
				Probe:       spec.ContinuesFromTheCursor,
				Reason:      spec.RepeatedPage,
				Law:         spec.WalkResumesWhereItLeftOff,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PagerHandle) {
					probe.ContinuesFromTheCursor(
						h,
						func(ctx context.Context, v int) (*Page, error) {
							return h.Subject().Fetch(ctx, v)
						},
						func(v *Page) int { return v.Next },
					)
				},
				Defect: pagerFetchServesOnePage,
			},
		},
		Calls: []conform.Call[Pager]{
			{
				Calls:    pagerFetchCalls,
				Over:     newPagerOver,
				Callable: PagerChecks.Fetch.Callable,
			},
		},
	}
}

// RunPagerSuite applies the checks of Pager, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunPagerSuite[T Pager](t *testing.T, newSubject func() T, opts ...PagerOption) {
	t.Helper()
	RunPagerSuiteFromHandle(t, func(*PagerHandle) T { return newSubject() }, opts...)
}

// RunPagerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPagerSuiteFromHandle[T Pager](
	t *testing.T,
	newSubject func(h *PagerHandle) T,
	opts ...PagerOption,
) {
	t.Helper()
	conform.Run(t, PagerSuite(), pagerSubjects(newSubject), opts...)
}

// FuzzPagerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPagerSuite[T Pager](
	f *testing.F,
	callable conform.Callable[Pager],
	newSubject func(h *PagerHandle) T,
	opts ...PagerOption,
) {
	f.Helper()
	conform.Fuzz(f, PagerSuite(), callable, pagerSubjects(newSubject), opts...)
}

// RunPagerBudgets applies the allocation ceilings of Pager and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunPagerBudgets[T Pager](
	t *testing.T,
	newSubject func(h *PagerHandle) T,
	opts ...PagerOption,
) {
	t.Helper()
	conform.RunBudgets(t, PagerSuite(), pagerSubjects(newSubject), opts...)
}

// BenchPagerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPagerSuite[T Pager](
	b conform.B,
	callable conform.Callable[Pager],
	newSubject func(h *PagerHandle) T,
	opts ...PagerOption,
) {
	b.Helper()
	conform.Bench(b, PagerSuite(), callable, pagerSubjects(newSubject), opts...)
}

// pagerSubjects returns the factory of the subjects of Pager that calls
// newSubject with the handle of each check.
func pagerSubjects[T Pager](newSubject func(h *PagerHandle) T) func(h *PagerHandle) Pager {
	return func(h *PagerHandle) Pager { return newSubject(h) }
}

// pagerFetchCalls generates the calls of Pager.Fetch with the inputs that
// each case draws.
var pagerFetchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pager) error {
	cursor := c.Draw(prop.Of[int](), "cursor")
	return func(ctx context.Context, subject Pager) error {
		_, err := subject.Fetch(ctx, cursor)
		return err
	}
})

// pagerOverEmbedded is Pager under the name of the field that pagerOver
// embeds.
type pagerOverEmbedded = Pager

// pagerOver is a subject of Pager whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type pagerOver struct {
	pagerOverEmbedded
	spend func()
}

// Fetch calls spend and returns its zero results.
func (o pagerOver) Fetch(context.Context, int) (r0 *Page, r1 error) {
	o.spend()
	return r0, r1
}

// newPagerOver returns the planted defect of a ceiling of Pager, whose
// callables call spend.
func newPagerOver(_ *PagerHandle, spend func()) Pager {
	return pagerOver{spend: spend}
}

// pagerFetchPanicsDefect is the planted defect of a Pager whose Fetch
// panics.
type pagerFetchPanicsDefect struct {
	pagerOverEmbedded
}

// Fetch is the override of the planted defect.
func (*pagerFetchPanicsDefect) Fetch(p0 context.Context, p1 int) (r0 *Page, r1 error) {
	panic("the planted defect of a Pager whose Fetch panics")
}

// pagerFetchPanics returns the planted defect of a Pager whose Fetch
// panics over a fresh subject.
func pagerFetchPanics(h *PagerHandle) Pager {
	return &pagerFetchPanicsDefect{pagerOverEmbedded: h.Fresh()}
}

// pagerFetchIgnoresCancelDefect is the planted defect of a Pager whose
// Fetch ignores the context.
type pagerFetchIgnoresCancelDefect struct {
	pagerOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *pagerFetchIgnoresCancelDefect) Fetch(p0 context.Context, p1 int) (r0 *Page, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.pagerOverEmbedded.Fetch(p0, p1)
}

// pagerFetchIgnoresCancel returns the planted defect of a Pager whose
// Fetch ignores the context over a fresh subject.
func pagerFetchIgnoresCancel(h *PagerHandle) Pager {
	return &pagerFetchIgnoresCancelDefect{pagerOverEmbedded: h.Fresh()}
}

// pagerFetchValueBesideDefect is the planted defect of a Pager whose Fetch
// returns a value beside an error.
type pagerFetchValueBesideDefect struct {
	pagerOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *pagerFetchValueBesideDefect) Fetch(p0 context.Context, p1 int) (r0 *Page, r1 error) {
	r0, r1 = d.pagerOverEmbedded.Fetch(p0, p1)
	if r1 != nil {
		r0 = &Page{Next: 7}
	}
	return r0, r1
}

// pagerFetchValueBeside returns the planted defect of a Pager whose Fetch
// returns a value beside an error over a fresh subject.
func pagerFetchValueBeside(h *PagerHandle) Pager {
	return &pagerFetchValueBesideDefect{pagerOverEmbedded: h.Fresh()}
}

// pagerFetchServesOnePageDefect is the planted defect of a Pager whose
// Fetch returns the first page to every cursor.
type pagerFetchServesOnePageDefect struct {
	pagerOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *pagerFetchServesOnePageDefect) Fetch(p0 context.Context, p1 int) (r0 *Page, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	p1 = *new(int)
	return d.pagerOverEmbedded.Fetch(p0, p1)
}

// pagerFetchServesOnePage returns the planted defect of a Pager whose
// Fetch returns the first page to every cursor over a fresh subject.
func pagerFetchServesOnePage(h *PagerHandle) Pager {
	return &pagerFetchServesOnePageDefect{pagerOverEmbedded: h.Fresh()}
}

// The derived checks of Log:
//
// Log.Append/smoke: Append returns without a panic for derived inputs.
//
// Log.Append/cancel: Append returns an error for a cancelled context.
//
// Log.Append/chain-causal: Replay returns the entries of Append in the
// order of their stamps.
//
// Log.Append/chain: Append grows Replay by the entry that it admits.
//
// Log.Replay/smoke: Replay returns without a panic for derived inputs.
//
// Log.Replay/cancel: Replay returns an error for a cancelled context.
//
// Log.Replay/zero-on-error: Replay returns the zero value beside an error.

// The derived checks of Log that witness their laws at the shallow end:
//
// Log.Append/smoke
//
// Log.Append/chain-causal
//
// Log.Replay/smoke

// The refused claims of Log:
//
// Log.Append/causal: Append does not return exactly one value. The remedy
// is to return one value from Append.

// The unclaimed classifications of Log:
//
// The shape:writer of Log.Append has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Log writes.
type (
	LogHandle = conform.Handle[Log]
	LogCheck  = conform.Check[Log]
	LogOption = conform.Option[Log]
)

// LogChecks contains the name of each callable of Log, with the identities
// of its checks.
var LogChecks = struct {
	Append logAppend
	Replay logReplay
}{
	Append: logAppend{"Log.Append"},
	Replay: logReplay{"Log.Replay"},
}

// logAppend is the name of Log.Append, with the identities of its checks.
type logAppend struct {
	conform.Callable[Log]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c logAppend) Smoke() conform.ID[Log] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c logAppend) Cancel() conform.ID[Log] {
	return c.Check("cancel")
}

// ChainCausal returns the identity of the check that Replay returns the
// entries of Append in the order of their stamps.
func (c logAppend) ChainCausal() conform.ID[Log] {
	return c.Check("chain-causal")
}

// Chain returns the identity of the check that Append grows Replay by the
// entry that it admits.
func (c logAppend) Chain() conform.ID[Log] {
	return c.Check("chain")
}

// logReplay is the name of Log.Replay, with the identities of its checks.
type logReplay struct {
	conform.Callable[Log]
}

// Smoke returns the identity of the check that Replay returns without a
// panic for derived inputs.
func (c logReplay) Smoke() conform.ID[Log] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Replay returns an error
// for a cancelled context.
func (c logReplay) Cancel() conform.ID[Log] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Replay returns the
// zero value beside an error.
func (c logReplay) ZeroOnError() conform.ID[Log] {
	return c.Check("zero-on-error")
}

// LogSuite returns the checks of Log, and the calls and the budgets of its
// callables.
func LogSuite() conform.Suite[Log] {
	return conform.Suite[Log]{
		Declaration: "Log",
		Checks: []conform.Check[Log]{
			{
				ID:          LogChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LogHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Entry{ID: "ID-a"})
					})
				},
				Defect: logAppendPanics,
			},
			{
				ID:          LogChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Entry{ID: "ID-a"})
					})
				},
				Defect: logAppendIgnoresCancel,
			},
			{
				ID:          LogChecks.Append.ChainCausal(),
				Statement:   "Replay returns the entries of Append in the order of their stamps",
				Claim:       spec.ChainCausal,
				Probe:       spec.ReplaysInStampOrder,
				Reason:      spec.UnorderedHistory,
				Law:         spec.HistoryReplaysInItsStampedOrder,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LogHandle) {
					probe.ReplaysInStampOrder(h, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Entry{ID: "ID-a"}
							if calls > 1 {
								w = Entry{ID: "ID-b"}
							}
							return h.Subject().Append(ctx, w)
						}
					}(), func(ctx context.Context) ([]int64, error) {
						r0, err := h.Subject().Replay(ctx)
						stamps := make([]int64, 0, len(r0))
						for _, e := range r0 {
							stamps = append(stamps, e.Seq)
						}
						return stamps, err
					})
				},
				Defect: logReplayReplaysReversed,
			},
			{
				ID:          LogChecks.Append.Chain(),
				Statement:   "Append grows Replay by the entry that it admits",
				Claim:       spec.Chain,
				Probe:       spec.ReplayGrows,
				Reason:      spec.ReplayNotGrown,
				Law:         spec.AdmittedEntryGrowsTheReplay,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.ReplayGrows(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Entry{ID: "ID-a"})
					}, func(ctx context.Context) ([]Entry, error) {
						return h.Subject().Replay(ctx)
					})
				},
				Defect: logAppendSucceedsAlways,
			},
			{
				ID:          LogChecks.Replay.Smoke(),
				Statement:   "Replay returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LogHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Replay(ctx)
						return err
					})
				},
				Defect: logReplayPanics,
			},
			{
				ID:          LogChecks.Replay.Cancel(),
				Statement:   "Replay returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Replay(ctx)
						return err
					})
				},
				Defect: logReplayIgnoresCancel,
			},
			{
				ID:          LogChecks.Replay.ZeroOnError(),
				Statement:   "Replay returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]Entry, error) {
						return h.Subject().Replay(ctx)
					}, []Entry(nil))
				},
				Defect: logReplayValueBeside,
			},
		},
		Calls: []conform.Call[Log]{
			{
				Calls:    logAppendCalls,
				Over:     newLogOver,
				Callable: LogChecks.Append.Callable,
			},
			{
				Calls:    logReplayCalls,
				Over:     newLogOver,
				Callable: LogChecks.Replay.Callable,
			},
		},
	}
}

// RunLogSuite applies the checks of Log, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunLogSuite[T Log](t *testing.T, newSubject func() T, opts ...LogOption) {
	t.Helper()
	RunLogSuiteFromHandle(t, func(*LogHandle) T { return newSubject() }, opts...)
}

// RunLogSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLogSuiteFromHandle[T Log](
	t *testing.T,
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	t.Helper()
	conform.Run(t, LogSuite(), logSubjects(newSubject), opts...)
}

// FuzzLogSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLogSuite[T Log](
	f *testing.F,
	callable conform.Callable[Log],
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	f.Helper()
	conform.Fuzz(f, LogSuite(), callable, logSubjects(newSubject), opts...)
}

// RunLogBudgets applies the allocation ceilings of Log and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunLogBudgets[T Log](t *testing.T, newSubject func(h *LogHandle) T, opts ...LogOption) {
	t.Helper()
	conform.RunBudgets(t, LogSuite(), logSubjects(newSubject), opts...)
}

// BenchLogSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLogSuite[T Log](
	b conform.B,
	callable conform.Callable[Log],
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	b.Helper()
	conform.Bench(b, LogSuite(), callable, logSubjects(newSubject), opts...)
}

// logSubjects returns the factory of the subjects of Log that calls
// newSubject with the handle of each check.
func logSubjects[T Log](newSubject func(h *LogHandle) T) func(h *LogHandle) Log {
	return func(h *LogHandle) Log { return newSubject(h) }
}

// logAppendCalls generates the calls of Log.Append with the inputs that
// each case draws.
var logAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Log) error {
	e := c.Draw(prop.Of[Entry](), "e")
	return func(ctx context.Context, subject Log) error {
		return subject.Append(ctx, e)
	}
})

// logReplayCalls generates the calls of Log.Replay with the inputs that
// each case draws.
var logReplayCalls = prop.Composite(func(*prop.Case) func(context.Context, Log) error {
	return func(ctx context.Context, subject Log) error {
		_, err := subject.Replay(ctx)
		return err
	}
})

// logOverEmbedded is Log under the name of the field that logOver embeds.
type logOverEmbedded = Log

// logOver is a subject of Log whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type logOver struct {
	logOverEmbedded
	spend func()
}

// Append calls spend and returns its zero results.
func (o logOver) Append(context.Context, Entry) (r0 error) {
	o.spend()
	return r0
}

// Replay calls spend and returns its zero results.
func (o logOver) Replay(context.Context) (r0 []Entry, r1 error) {
	o.spend()
	return r0, r1
}

// newLogOver returns the planted defect of a ceiling of Log, whose
// callables call spend.
func newLogOver(_ *LogHandle, spend func()) Log {
	return logOver{spend: spend}
}

// logAppendPanicsDefect is the planted defect of a Log whose Append
// panics.
type logAppendPanicsDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (*logAppendPanicsDefect) Append(p0 context.Context, p1 Entry) (r0 error) {
	panic("the planted defect of a Log whose Append panics")
}

// logAppendPanics returns the planted defect of a Log whose Append panics
// over a fresh subject.
func logAppendPanics(h *LogHandle) Log {
	return &logAppendPanicsDefect{logOverEmbedded: h.Fresh()}
}

// logAppendIgnoresCancelDefect is the planted defect of a Log whose Append
// ignores the context.
type logAppendIgnoresCancelDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (d *logAppendIgnoresCancelDefect) Append(p0 context.Context, p1 Entry) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.logOverEmbedded.Append(p0, p1)
}

// logAppendIgnoresCancel returns the planted defect of a Log whose Append
// ignores the context over a fresh subject.
func logAppendIgnoresCancel(h *LogHandle) Log {
	return &logAppendIgnoresCancelDefect{logOverEmbedded: h.Fresh()}
}

// logReplayReplaysReversedDefect is the planted defect of a Log whose
// Replay returns the history of Append reversed.
type logReplayReplaysReversedDefect struct {
	logOverEmbedded
}

// Replay is the override of the planted defect.
func (d *logReplayReplaysReversedDefect) Replay(p0 context.Context) (r0 []Entry, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.logOverEmbedded.Replay(p0)
	slices.Reverse(r0)
	return r0, r1
}

// logReplayReplaysReversed returns the planted defect of a Log whose
// Replay returns the history of Append reversed over a fresh subject.
func logReplayReplaysReversed(h *LogHandle) Log {
	return &logReplayReplaysReversedDefect{logOverEmbedded: h.Fresh()}
}

// logAppendSucceedsAlwaysDefect is the planted defect of a Log whose
// Append returns success and appends nothing.
type logAppendSucceedsAlwaysDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (*logAppendSucceedsAlwaysDefect) Append(p0 context.Context, p1 Entry) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// logAppendSucceedsAlways returns the planted defect of a Log whose Append
// returns success and appends nothing over a fresh subject.
func logAppendSucceedsAlways(h *LogHandle) Log {
	return &logAppendSucceedsAlwaysDefect{logOverEmbedded: h.Fresh()}
}

// logReplayPanicsDefect is the planted defect of a Log whose Replay
// panics.
type logReplayPanicsDefect struct {
	logOverEmbedded
}

// Replay is the override of the planted defect.
func (*logReplayPanicsDefect) Replay(p0 context.Context) (r0 []Entry, r1 error) {
	panic("the planted defect of a Log whose Replay panics")
}

// logReplayPanics returns the planted defect of a Log whose Replay panics
// over a fresh subject.
func logReplayPanics(h *LogHandle) Log {
	return &logReplayPanicsDefect{logOverEmbedded: h.Fresh()}
}

// logReplayIgnoresCancelDefect is the planted defect of a Log whose Replay
// ignores the context.
type logReplayIgnoresCancelDefect struct {
	logOverEmbedded
}

// Replay is the override of the planted defect.
func (d *logReplayIgnoresCancelDefect) Replay(p0 context.Context) (r0 []Entry, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.logOverEmbedded.Replay(p0)
}

// logReplayIgnoresCancel returns the planted defect of a Log whose Replay
// ignores the context over a fresh subject.
func logReplayIgnoresCancel(h *LogHandle) Log {
	return &logReplayIgnoresCancelDefect{logOverEmbedded: h.Fresh()}
}

// logReplayValueBesideDefect is the planted defect of a Log whose Replay
// returns a value beside an error.
type logReplayValueBesideDefect struct {
	logOverEmbedded
}

// Replay is the override of the planted defect.
func (d *logReplayValueBesideDefect) Replay(p0 context.Context) (r0 []Entry, r1 error) {
	r0, r1 = d.logOverEmbedded.Replay(p0)
	if r1 != nil {
		r0 = []Entry{{ID: "ID-b"}}
	}
	return r0, r1
}

// logReplayValueBeside returns the planted defect of a Log whose Replay
// returns a value beside an error over a fresh subject.
func logReplayValueBeside(h *LogHandle) Log {
	return &logReplayValueBesideDefect{logOverEmbedded: h.Fresh()}
}

// The derived checks of Opener:
//
// Opener.Open/smoke: Open returns without a panic for derived inputs.
//
// Opener.Open/cancel: Open returns an error for a cancelled context.
//
// Opener.Open/cursor-open: the cursor of Open fails with the declared
// sentinel once it is closed.

// The derived checks of Opener that witness their laws at the shallow end:
//
// Opener.Open/smoke

// The refused claims of Opener:
//
// Opener.Open/zero-on-error: eidos derives no sample of a result of Open.
// The remedy is to author a sample of the result type with the directive
// sample.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Opener writes.
type (
	OpenerHandle = conform.Handle[Opener]
	OpenerCheck  = conform.Check[Opener]
	OpenerOption = conform.Option[Opener]
)

// OpenerChecks contains the name of each callable of Opener, with the
// identities of its checks.
var OpenerChecks = struct {
	Open openerOpen
}{
	Open: openerOpen{"Opener.Open"},
}

// openerOpen is the name of Opener.Open, with the identities of its
// checks.
type openerOpen struct {
	conform.Callable[Opener]
}

// Smoke returns the identity of the check that Open returns without a
// panic for derived inputs.
func (c openerOpen) Smoke() conform.ID[Opener] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Open returns an error for
// a cancelled context.
func (c openerOpen) Cancel() conform.ID[Opener] {
	return c.Check("cancel")
}

// CursorOpen returns the identity of the check that the cursor of Open
// fails with the declared sentinel once it is closed.
func (c openerOpen) CursorOpen() conform.ID[Opener] {
	return c.Check("cursor-open")
}

// OpenerSuite returns the checks of Opener, and the calls and the budgets
// of its callables.
func OpenerSuite() conform.Suite[Opener] {
	return conform.Suite[Opener]{
		Declaration: "Opener",
		Checks: []conform.Check[Opener]{
			{
				ID:          OpenerChecks.Open.Smoke(),
				Statement:   "Open returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *OpenerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Open(ctx)
						return err
					})
				},
				Defect: openerOpenPanics,
			},
			{
				ID:          OpenerChecks.Open.Cancel(),
				Statement:   "Open returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OpenerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Open(ctx)
						return err
					})
				},
				Defect: openerOpenIgnoresCancel,
			},
			{
				ID:          OpenerChecks.Open.CursorOpen(),
				Statement:   "the cursor of Open fails with the declared sentinel once it is closed",
				Claim:       spec.CursorOpen,
				Probe:       spec.SealsAfterClose,
				Reason:      spec.AliveAfterClose,
				Law:         spec.WorkAfterCloseIsRefused,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *OpenerHandle) {
					probe.SealsAfterClose(h, func(ctx context.Context) (Cursor, error) {
						return h.Subject().Open(ctx)
					}, func(ctx context.Context, handle Cursor) error {
						return handle.Close(ctx)
					}, func(ctx context.Context, handle Cursor) error {
						_, _, err := handle.Next(ctx)
						return err
					}, ErrClosed)
				},
				Defect: openerOpenOpensALiveHandle,
			},
		},
		Calls: []conform.Call[Opener]{
			{
				Calls:    openerOpenCalls,
				Over:     newOpenerOver,
				Callable: OpenerChecks.Open.Callable,
			},
		},
	}
}

// RunOpenerSuite applies the checks of Opener, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunOpenerSuite[T Opener](t *testing.T, newSubject func() T, opts ...OpenerOption) {
	t.Helper()
	RunOpenerSuiteFromHandle(t, func(*OpenerHandle) T { return newSubject() }, opts...)
}

// RunOpenerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunOpenerSuiteFromHandle[T Opener](
	t *testing.T,
	newSubject func(h *OpenerHandle) T,
	opts ...OpenerOption,
) {
	t.Helper()
	conform.Run(t, OpenerSuite(), openerSubjects(newSubject), opts...)
}

// FuzzOpenerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzOpenerSuite[T Opener](
	f *testing.F,
	callable conform.Callable[Opener],
	newSubject func(h *OpenerHandle) T,
	opts ...OpenerOption,
) {
	f.Helper()
	conform.Fuzz(f, OpenerSuite(), callable, openerSubjects(newSubject), opts...)
}

// RunOpenerBudgets applies the allocation ceilings of Opener and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunOpenerBudgets[T Opener](
	t *testing.T,
	newSubject func(h *OpenerHandle) T,
	opts ...OpenerOption,
) {
	t.Helper()
	conform.RunBudgets(t, OpenerSuite(), openerSubjects(newSubject), opts...)
}

// BenchOpenerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchOpenerSuite[T Opener](
	b conform.B,
	callable conform.Callable[Opener],
	newSubject func(h *OpenerHandle) T,
	opts ...OpenerOption,
) {
	b.Helper()
	conform.Bench(b, OpenerSuite(), callable, openerSubjects(newSubject), opts...)
}

// openerSubjects returns the factory of the subjects of Opener that calls
// newSubject with the handle of each check.
func openerSubjects[T Opener](newSubject func(h *OpenerHandle) T) func(h *OpenerHandle) Opener {
	return func(h *OpenerHandle) Opener { return newSubject(h) }
}

// openerOpenCalls generates the calls of Opener.Open with the inputs that
// each case draws.
var openerOpenCalls = prop.Composite(func(*prop.Case) func(context.Context, Opener) error {
	return func(ctx context.Context, subject Opener) error {
		_, err := subject.Open(ctx)
		return err
	}
})

// openerOverEmbedded is Opener under the name of the field that openerOver
// embeds.
type openerOverEmbedded = Opener

// openerOver is a subject of Opener whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type openerOver struct {
	openerOverEmbedded
	spend func()
}

// Open calls spend and returns its zero results.
func (o openerOver) Open(context.Context) (r0 Cursor, r1 error) {
	o.spend()
	return r0, r1
}

// newOpenerOver returns the planted defect of a ceiling of Opener, whose
// callables call spend.
func newOpenerOver(_ *OpenerHandle, spend func()) Opener {
	return openerOver{spend: spend}
}

// openerOpenPanicsDefect is the planted defect of a Opener whose Open
// panics.
type openerOpenPanicsDefect struct {
	openerOverEmbedded
}

// Open is the override of the planted defect.
func (*openerOpenPanicsDefect) Open(p0 context.Context) (r0 Cursor, r1 error) {
	panic("the planted defect of a Opener whose Open panics")
}

// openerOpenPanics returns the planted defect of a Opener whose Open
// panics over a fresh subject.
func openerOpenPanics(h *OpenerHandle) Opener {
	return &openerOpenPanicsDefect{openerOverEmbedded: h.Fresh()}
}

// openerOpenIgnoresCancelDefect is the planted defect of a Opener whose
// Open ignores the context.
type openerOpenIgnoresCancelDefect struct {
	openerOverEmbedded
}

// Open is the override of the planted defect.
func (d *openerOpenIgnoresCancelDefect) Open(p0 context.Context) (r0 Cursor, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.openerOverEmbedded.Open(p0)
}

// openerOpenIgnoresCancel returns the planted defect of a Opener whose
// Open ignores the context over a fresh subject.
func openerOpenIgnoresCancel(h *OpenerHandle) Opener {
	return &openerOpenIgnoresCancelDefect{openerOverEmbedded: h.Fresh()}
}

// openerOpenOpensALiveHandleDefect is the planted defect of a Opener whose
// Open returns a cursor that drains after its close.
type openerOpenOpensALiveHandleDefect struct {
	openerOverEmbedded
}

// Open is the override of the planted defect.
func (d *openerOpenOpensALiveHandleDefect) Open(p0 context.Context) (r0 Cursor, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.openerOverEmbedded.Open(p0)
	r0 = openerOpenOpensALiveHandleHandle{r0}
	return r0, r1
}

// openerOpenOpensALiveHandleHandle is a handle whose Next returns its zero
// results and no error, also after its close.
type openerOpenOpensALiveHandleHandle struct {
	Cursor
}

// Next returns its zero results and no error.
func (openerOpenOpensALiveHandleHandle) Next(p0 context.Context) (r0 Value, r1 bool, r2 error) {
	return r0, r1, r2
}

// openerOpenOpensALiveHandle returns the planted defect of a Opener whose
// Open returns a cursor that drains after its close over a fresh subject.
func openerOpenOpensALiveHandle(h *OpenerHandle) Opener {
	return &openerOpenOpensALiveHandleDefect{openerOverEmbedded: h.Fresh()}
}

// The derived checks of Ledger:
//
// Ledger.Begin/smoke: Begin returns without a panic for derived inputs.
//
// Ledger.Begin/cancel: Begin returns an error for a cancelled context.
//
// Ledger.Begin/zero-on-error: Begin returns the zero value beside an
// error.
//
// Ledger.Begin/tx: Rollback of a transaction that Begin opened and Commit
// settled fails with the declared sentinel.
//
// Ledger.Commit/smoke: Commit returns without a panic for derived inputs.
//
// Ledger.Commit/cancel: Commit returns an error for a cancelled context.
//
// Ledger.Rollback/smoke: Rollback returns without a panic for derived
// inputs.
//
// Ledger.Rollback/cancel: Rollback returns an error for a cancelled
// context.

// The derived checks of Ledger that witness their laws at the shallow end:
//
// Ledger.Begin/smoke
//
// Ledger.Commit/smoke
//
// Ledger.Rollback/smoke

// The unclaimed classifications of Ledger:
//
// The shape:writer of Ledger.Commit has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Ledger.Rollback has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Ledger writes.
type (
	LedgerHandle = conform.Handle[Ledger]
	LedgerCheck  = conform.Check[Ledger]
	LedgerOption = conform.Option[Ledger]
)

// LedgerChecks contains the name of each callable of Ledger, with the
// identities of its checks.
var LedgerChecks = struct {
	Begin    ledgerBegin
	Commit   ledgerCommit
	Rollback ledgerRollback
}{
	Begin:    ledgerBegin{"Ledger.Begin"},
	Commit:   ledgerCommit{"Ledger.Commit"},
	Rollback: ledgerRollback{"Ledger.Rollback"},
}

// ledgerBegin is the name of Ledger.Begin, with the identities of its
// checks.
type ledgerBegin struct {
	conform.Callable[Ledger]
}

// Smoke returns the identity of the check that Begin returns without a
// panic for derived inputs.
func (c ledgerBegin) Smoke() conform.ID[Ledger] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Begin returns an error for
// a cancelled context.
func (c ledgerBegin) Cancel() conform.ID[Ledger] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Begin returns the
// zero value beside an error.
func (c ledgerBegin) ZeroOnError() conform.ID[Ledger] {
	return c.Check("zero-on-error")
}

// Tx returns the identity of the check that Rollback of a transaction that
// Begin opened and Commit settled fails with the declared sentinel.
func (c ledgerBegin) Tx() conform.ID[Ledger] {
	return c.Check("tx")
}

// ledgerCommit is the name of Ledger.Commit, with the identities of its
// checks.
type ledgerCommit struct {
	conform.Callable[Ledger]
}

// Smoke returns the identity of the check that Commit returns without a
// panic for derived inputs.
func (c ledgerCommit) Smoke() conform.ID[Ledger] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Commit returns an error
// for a cancelled context.
func (c ledgerCommit) Cancel() conform.ID[Ledger] {
	return c.Check("cancel")
}

// ledgerRollback is the name of Ledger.Rollback, with the identities of
// its checks.
type ledgerRollback struct {
	conform.Callable[Ledger]
}

// Smoke returns the identity of the check that Rollback returns without a
// panic for derived inputs.
func (c ledgerRollback) Smoke() conform.ID[Ledger] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Rollback returns an error
// for a cancelled context.
func (c ledgerRollback) Cancel() conform.ID[Ledger] {
	return c.Check("cancel")
}

// LedgerSuite returns the checks of Ledger, and the calls and the budgets
// of its callables.
func LedgerSuite() conform.Suite[Ledger] {
	return conform.Suite[Ledger]{
		Declaration: "Ledger",
		Checks: []conform.Check[Ledger]{
			{
				ID:          LedgerChecks.Begin.Smoke(),
				Statement:   "Begin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LedgerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx)
						return err
					})
				},
				Defect: ledgerBeginPanics,
			},
			{
				ID:          LedgerChecks.Begin.Cancel(),
				Statement:   "Begin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LedgerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx)
						return err
					})
				},
				Defect: ledgerBeginIgnoresCancel,
			},
			{
				ID:          LedgerChecks.Begin.ZeroOnError(),
				Statement:   "Begin returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LedgerHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Tx, error) {
						return h.Subject().Begin(ctx)
					}, Tx{})
				},
				Defect: ledgerBeginValueBeside,
			},
			{
				ID:          LedgerChecks.Begin.Tx(),
				Statement:   "Rollback of a transaction that Begin opened and Commit settled fails with the declared sentinel",
				Claim:       spec.Tx,
				Probe:       spec.SealsAfterClose,
				Reason:      spec.AliveAfterClose,
				Law:         spec.SettledTransactionIsNotSettledAgain,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LedgerHandle) {
					probe.SealsAfterClose(h, func(ctx context.Context) (Tx, error) {
						return h.Subject().Begin(ctx)
					}, func(ctx context.Context, handle Tx) error {
						return h.Subject().Commit(ctx, handle)
					}, func(ctx context.Context, handle Tx) error {
						return h.Subject().Rollback(ctx, handle)
					}, ErrSettled)
				},
				Defect: ledgerRollbackSettlesTwice,
			},
			{
				ID:          LedgerChecks.Commit.Smoke(),
				Statement:   "Commit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LedgerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx, Tx{ID: 42})
					})
				},
				Defect: ledgerCommitPanics,
			},
			{
				ID:          LedgerChecks.Commit.Cancel(),
				Statement:   "Commit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LedgerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx, Tx{ID: 42})
					})
				},
				Defect: ledgerCommitIgnoresCancel,
			},
			{
				ID:          LedgerChecks.Rollback.Smoke(),
				Statement:   "Rollback returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LedgerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx, Tx{ID: 42})
					})
				},
				Defect: ledgerRollbackPanics,
			},
			{
				ID:          LedgerChecks.Rollback.Cancel(),
				Statement:   "Rollback returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LedgerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx, Tx{ID: 42})
					})
				},
				Defect: ledgerRollbackIgnoresCancel,
			},
		},
		Calls: []conform.Call[Ledger]{
			{
				Calls:    ledgerBeginCalls,
				Over:     newLedgerOver,
				Callable: LedgerChecks.Begin.Callable,
			},
			{
				Calls:    ledgerCommitCalls,
				Over:     newLedgerOver,
				Callable: LedgerChecks.Commit.Callable,
			},
			{
				Calls:    ledgerRollbackCalls,
				Over:     newLedgerOver,
				Callable: LedgerChecks.Rollback.Callable,
			},
		},
	}
}

// RunLedgerSuite applies the checks of Ledger, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLedgerSuite[T Ledger](t *testing.T, newSubject func() T, opts ...LedgerOption) {
	t.Helper()
	RunLedgerSuiteFromHandle(t, func(*LedgerHandle) T { return newSubject() }, opts...)
}

// RunLedgerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLedgerSuiteFromHandle[T Ledger](
	t *testing.T,
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	t.Helper()
	conform.Run(t, LedgerSuite(), ledgerSubjects(newSubject), opts...)
}

// FuzzLedgerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLedgerSuite[T Ledger](
	f *testing.F,
	callable conform.Callable[Ledger],
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	f.Helper()
	conform.Fuzz(f, LedgerSuite(), callable, ledgerSubjects(newSubject), opts...)
}

// RunLedgerBudgets applies the allocation ceilings of Ledger and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLedgerBudgets[T Ledger](
	t *testing.T,
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	t.Helper()
	conform.RunBudgets(t, LedgerSuite(), ledgerSubjects(newSubject), opts...)
}

// BenchLedgerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLedgerSuite[T Ledger](
	b conform.B,
	callable conform.Callable[Ledger],
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	b.Helper()
	conform.Bench(b, LedgerSuite(), callable, ledgerSubjects(newSubject), opts...)
}

// ledgerSubjects returns the factory of the subjects of Ledger that calls
// newSubject with the handle of each check.
func ledgerSubjects[T Ledger](newSubject func(h *LedgerHandle) T) func(h *LedgerHandle) Ledger {
	return func(h *LedgerHandle) Ledger { return newSubject(h) }
}

// ledgerBeginCalls generates the calls of Ledger.Begin with the inputs
// that each case draws.
var ledgerBeginCalls = prop.Composite(func(*prop.Case) func(context.Context, Ledger) error {
	return func(ctx context.Context, subject Ledger) error {
		_, err := subject.Begin(ctx)
		return err
	}
})

// ledgerCommitCalls generates the calls of Ledger.Commit with the inputs
// that each case draws.
var ledgerCommitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Ledger) error {
	tx := c.Draw(prop.Of[Tx](), "tx")
	return func(ctx context.Context, subject Ledger) error {
		return subject.Commit(ctx, tx)
	}
})

// ledgerRollbackCalls generates the calls of Ledger.Rollback with the
// inputs that each case draws.
var ledgerRollbackCalls = prop.Composite(func(c *prop.Case) func(context.Context, Ledger) error {
	tx := c.Draw(prop.Of[Tx](), "tx")
	return func(ctx context.Context, subject Ledger) error {
		return subject.Rollback(ctx, tx)
	}
})

// ledgerOverEmbedded is Ledger under the name of the field that ledgerOver
// embeds.
type ledgerOverEmbedded = Ledger

// ledgerOver is a subject of Ledger whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type ledgerOver struct {
	ledgerOverEmbedded
	spend func()
}

// Begin calls spend and returns its zero results.
func (o ledgerOver) Begin(context.Context) (r0 Tx, r1 error) {
	o.spend()
	return r0, r1
}

// Commit calls spend and returns its zero results.
func (o ledgerOver) Commit(context.Context, Tx) (r0 error) {
	o.spend()
	return r0
}

// Rollback calls spend and returns its zero results.
func (o ledgerOver) Rollback(context.Context, Tx) (r0 error) {
	o.spend()
	return r0
}

// newLedgerOver returns the planted defect of a ceiling of Ledger, whose
// callables call spend.
func newLedgerOver(_ *LedgerHandle, spend func()) Ledger {
	return ledgerOver{spend: spend}
}

// ledgerBeginPanicsDefect is the planted defect of a Ledger whose Begin
// panics.
type ledgerBeginPanicsDefect struct {
	ledgerOverEmbedded
}

// Begin is the override of the planted defect.
func (*ledgerBeginPanicsDefect) Begin(p0 context.Context) (r0 Tx, r1 error) {
	panic("the planted defect of a Ledger whose Begin panics")
}

// ledgerBeginPanics returns the planted defect of a Ledger whose Begin
// panics over a fresh subject.
func ledgerBeginPanics(h *LedgerHandle) Ledger {
	return &ledgerBeginPanicsDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerBeginIgnoresCancelDefect is the planted defect of a Ledger whose
// Begin ignores the context.
type ledgerBeginIgnoresCancelDefect struct {
	ledgerOverEmbedded
}

// Begin is the override of the planted defect.
func (d *ledgerBeginIgnoresCancelDefect) Begin(p0 context.Context) (r0 Tx, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.ledgerOverEmbedded.Begin(p0)
}

// ledgerBeginIgnoresCancel returns the planted defect of a Ledger whose
// Begin ignores the context over a fresh subject.
func ledgerBeginIgnoresCancel(h *LedgerHandle) Ledger {
	return &ledgerBeginIgnoresCancelDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerBeginValueBesideDefect is the planted defect of a Ledger whose
// Begin returns a value beside an error.
type ledgerBeginValueBesideDefect struct {
	ledgerOverEmbedded
}

// Begin is the override of the planted defect.
func (d *ledgerBeginValueBesideDefect) Begin(p0 context.Context) (r0 Tx, r1 error) {
	r0, r1 = d.ledgerOverEmbedded.Begin(p0)
	if r1 != nil {
		r0 = Tx{ID: 7}
	}
	return r0, r1
}

// ledgerBeginValueBeside returns the planted defect of a Ledger whose
// Begin returns a value beside an error over a fresh subject.
func ledgerBeginValueBeside(h *LedgerHandle) Ledger {
	return &ledgerBeginValueBesideDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerRollbackSettlesTwiceDefect is the planted defect of a Ledger whose
// Rollback settles a transaction that Commit settled.
type ledgerRollbackSettlesTwiceDefect struct {
	ledgerOverEmbedded
}

// Rollback is the override of the planted defect.
func (*ledgerRollbackSettlesTwiceDefect) Rollback(p0 context.Context, p1 Tx) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// ledgerRollbackSettlesTwice returns the planted defect of a Ledger whose
// Rollback settles a transaction that Commit settled over a fresh subject.
func ledgerRollbackSettlesTwice(h *LedgerHandle) Ledger {
	return &ledgerRollbackSettlesTwiceDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerCommitPanicsDefect is the planted defect of a Ledger whose Commit
// panics.
type ledgerCommitPanicsDefect struct {
	ledgerOverEmbedded
}

// Commit is the override of the planted defect.
func (*ledgerCommitPanicsDefect) Commit(p0 context.Context, p1 Tx) (r0 error) {
	panic("the planted defect of a Ledger whose Commit panics")
}

// ledgerCommitPanics returns the planted defect of a Ledger whose Commit
// panics over a fresh subject.
func ledgerCommitPanics(h *LedgerHandle) Ledger {
	return &ledgerCommitPanicsDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerCommitIgnoresCancelDefect is the planted defect of a Ledger whose
// Commit ignores the context.
type ledgerCommitIgnoresCancelDefect struct {
	ledgerOverEmbedded
}

// Commit is the override of the planted defect.
func (d *ledgerCommitIgnoresCancelDefect) Commit(p0 context.Context, p1 Tx) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.ledgerOverEmbedded.Commit(p0, p1)
}

// ledgerCommitIgnoresCancel returns the planted defect of a Ledger whose
// Commit ignores the context over a fresh subject.
func ledgerCommitIgnoresCancel(h *LedgerHandle) Ledger {
	return &ledgerCommitIgnoresCancelDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerRollbackPanicsDefect is the planted defect of a Ledger whose
// Rollback panics.
type ledgerRollbackPanicsDefect struct {
	ledgerOverEmbedded
}

// Rollback is the override of the planted defect.
func (*ledgerRollbackPanicsDefect) Rollback(p0 context.Context, p1 Tx) (r0 error) {
	panic("the planted defect of a Ledger whose Rollback panics")
}

// ledgerRollbackPanics returns the planted defect of a Ledger whose
// Rollback panics over a fresh subject.
func ledgerRollbackPanics(h *LedgerHandle) Ledger {
	return &ledgerRollbackPanicsDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerRollbackIgnoresCancelDefect is the planted defect of a Ledger
// whose Rollback ignores the context.
type ledgerRollbackIgnoresCancelDefect struct {
	ledgerOverEmbedded
}

// Rollback is the override of the planted defect.
func (d *ledgerRollbackIgnoresCancelDefect) Rollback(p0 context.Context, p1 Tx) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.ledgerOverEmbedded.Rollback(p0, p1)
}

// ledgerRollbackIgnoresCancel returns the planted defect of a Ledger whose
// Rollback ignores the context over a fresh subject.
func ledgerRollbackIgnoresCancel(h *LedgerHandle) Ledger {
	return &ledgerRollbackIgnoresCancelDefect{ledgerOverEmbedded: h.Fresh()}
}

// The derived checks of Guarded:
//
// Guarded.Put/smoke: Put returns without a panic for derived inputs.
//
// Guarded.Put/cancel: Put returns an error for a cancelled context.
//
// Guarded.Put/if-match: Put refuses a value that Match rejects.
//
// Guarded.Match/smoke: Match returns without a panic for derived inputs.
//
// Guarded.Match/cancel: Match returns an error for a cancelled context.

// The derived checks of Guarded that witness their laws at the shallow
// end:
//
// Guarded.Put/smoke
//
// Guarded.Match/smoke

// The unclaimed classifications of Guarded:
//
// The shape:composite-writer of Guarded.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Guarded writes.
type (
	GuardedHandle = conform.Handle[Guarded]
	GuardedCheck  = conform.Check[Guarded]
	GuardedOption = conform.Option[Guarded]
)

// GuardedChecks contains the name of each callable of Guarded, with the
// identities of its checks.
var GuardedChecks = struct {
	Put   guardedPut
	Match guardedMatch
}{
	Put:   guardedPut{"Guarded.Put"},
	Match: guardedMatch{"Guarded.Match"},
}

// guardedPut is the name of Guarded.Put, with the identities of its
// checks.
type guardedPut struct {
	conform.Callable[Guarded]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c guardedPut) Smoke() conform.ID[Guarded] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c guardedPut) Cancel() conform.ID[Guarded] {
	return c.Check("cancel")
}

// IfMatch returns the identity of the check that Put refuses a value that
// Match rejects.
func (c guardedPut) IfMatch() conform.ID[Guarded] {
	return c.Check("if-match")
}

// guardedMatch is the name of Guarded.Match, with the identities of its
// checks.
type guardedMatch struct {
	conform.Callable[Guarded]
}

// Smoke returns the identity of the check that Match returns without a
// panic for derived inputs.
func (c guardedMatch) Smoke() conform.ID[Guarded] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Match returns an error for
// a cancelled context.
func (c guardedMatch) Cancel() conform.ID[Guarded] {
	return c.Check("cancel")
}

// GuardedSuite returns the checks of Guarded, and the calls and the
// budgets of its callables.
func GuardedSuite() conform.Suite[Guarded] {
	return conform.Suite[Guarded]{
		Declaration: "Guarded",
		Checks: []conform.Check[Guarded]{
			{
				ID:          GuardedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *GuardedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Meta{Revision: 42}, &Value{Key: "Key-a"})
					})
				},
				Defect: guardedPutPanics,
			},
			{
				ID:          GuardedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *GuardedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Meta{Revision: 42}, &Value{Key: "Key-a"})
					})
				},
				Defect: guardedPutIgnoresCancel,
			},
			{
				ID:          GuardedChecks.Put.IfMatch(),
				Statement:   "Put refuses a value that Match rejects",
				Claim:       spec.IfMatch,
				Probe:       spec.RefusesWhatThePredicateRejects,
				Reason:      spec.UnmatchedWrite,
				Law:         spec.WriteIsRefusedWhereThePredicateRejectsIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *GuardedHandle) {
					probe.RefusesWhatThePredicateRejects(h, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := &Value{Key: "Key-a"}
							w.Body = string("Body-a")
							if calls > 1 {
								w.Body = string("Body-b")
							}
							return h.Subject().Put(ctx, Meta{Revision: 42}, w)
						}
					}(), func(ctx context.Context) (bool, error) {
						w := &Value{Key: "Key-a"}
						w.Body = string("Body-b")
						return h.Subject().Match(ctx, Meta{Revision: 42}, w)
					})
				},
				Defect: guardedPutIgnoresThePredicate,
			},
			{
				ID:          GuardedChecks.Match.Smoke(),
				Statement:   "Match returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *GuardedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Match(ctx, Meta{Revision: 42}, &Value{Key: "Key-a"})
						return err
					})
				},
				Defect: guardedMatchPanics,
			},
			{
				ID:          GuardedChecks.Match.Cancel(),
				Statement:   "Match returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *GuardedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Match(ctx, Meta{Revision: 42}, &Value{Key: "Key-a"})
						return err
					})
				},
				Defect: guardedMatchIgnoresCancel,
			},
		},
		Calls: []conform.Call[Guarded]{
			{
				Calls:    guardedPutCalls,
				Over:     newGuardedOver,
				Callable: GuardedChecks.Put.Callable,
			},
			{
				Calls:    guardedMatchCalls,
				Over:     newGuardedOver,
				Callable: GuardedChecks.Match.Callable,
			},
		},
	}
}

// RunGuardedSuite applies the checks of Guarded, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunGuardedSuite[T Guarded](t *testing.T, newSubject func() T, opts ...GuardedOption) {
	t.Helper()
	RunGuardedSuiteFromHandle(t, func(*GuardedHandle) T { return newSubject() }, opts...)
}

// RunGuardedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunGuardedSuiteFromHandle[T Guarded](
	t *testing.T,
	newSubject func(h *GuardedHandle) T,
	opts ...GuardedOption,
) {
	t.Helper()
	conform.Run(t, GuardedSuite(), guardedSubjects(newSubject), opts...)
}

// FuzzGuardedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzGuardedSuite[T Guarded](
	f *testing.F,
	callable conform.Callable[Guarded],
	newSubject func(h *GuardedHandle) T,
	opts ...GuardedOption,
) {
	f.Helper()
	conform.Fuzz(f, GuardedSuite(), callable, guardedSubjects(newSubject), opts...)
}

// RunGuardedBudgets applies the allocation ceilings of Guarded and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunGuardedBudgets[T Guarded](
	t *testing.T,
	newSubject func(h *GuardedHandle) T,
	opts ...GuardedOption,
) {
	t.Helper()
	conform.RunBudgets(t, GuardedSuite(), guardedSubjects(newSubject), opts...)
}

// BenchGuardedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchGuardedSuite[T Guarded](
	b conform.B,
	callable conform.Callable[Guarded],
	newSubject func(h *GuardedHandle) T,
	opts ...GuardedOption,
) {
	b.Helper()
	conform.Bench(b, GuardedSuite(), callable, guardedSubjects(newSubject), opts...)
}

// guardedSubjects returns the factory of the subjects of Guarded that
// calls newSubject with the handle of each check.
func guardedSubjects[T Guarded](
	newSubject func(h *GuardedHandle) T,
) func(h *GuardedHandle) Guarded {
	return func(h *GuardedHandle) Guarded { return newSubject(h) }
}

// guardedPutCalls generates the calls of Guarded.Put with the inputs that
// each case draws.
var guardedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Guarded) error {
	meta := c.Draw(prop.Of[Meta](), "meta")
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Guarded) error {
		return subject.Put(ctx, meta, &v)
	}
})

// guardedMatchCalls generates the calls of Guarded.Match with the inputs
// that each case draws.
var guardedMatchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Guarded) error {
	meta := c.Draw(prop.Of[Meta](), "meta")
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Guarded) error {
		_, err := subject.Match(ctx, meta, &v)
		return err
	}
})

// guardedOverEmbedded is Guarded under the name of the field that
// guardedOver embeds.
type guardedOverEmbedded = Guarded

// guardedOver is a subject of Guarded whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type guardedOver struct {
	guardedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o guardedOver) Put(context.Context, Meta, *Value) (r0 error) {
	o.spend()
	return r0
}

// Match calls spend and returns its zero results.
func (o guardedOver) Match(context.Context, Meta, *Value) (r0 bool, r1 error) {
	o.spend()
	return r0, r1
}

// newGuardedOver returns the planted defect of a ceiling of Guarded, whose
// callables call spend.
func newGuardedOver(_ *GuardedHandle, spend func()) Guarded {
	return guardedOver{spend: spend}
}

// guardedPutPanicsDefect is the planted defect of a Guarded whose Put
// panics.
type guardedPutPanicsDefect struct {
	guardedOverEmbedded
}

// Put is the override of the planted defect.
func (*guardedPutPanicsDefect) Put(p0 context.Context, p1 Meta, p2 *Value) (r0 error) {
	panic("the planted defect of a Guarded whose Put panics")
}

// guardedPutPanics returns the planted defect of a Guarded whose Put
// panics over a fresh subject.
func guardedPutPanics(h *GuardedHandle) Guarded {
	return &guardedPutPanicsDefect{guardedOverEmbedded: h.Fresh()}
}

// guardedPutIgnoresCancelDefect is the planted defect of a Guarded whose
// Put ignores the context.
type guardedPutIgnoresCancelDefect struct {
	guardedOverEmbedded
}

// Put is the override of the planted defect.
func (d *guardedPutIgnoresCancelDefect) Put(p0 context.Context, p1 Meta, p2 *Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.guardedOverEmbedded.Put(p0, p1, p2)
}

// guardedPutIgnoresCancel returns the planted defect of a Guarded whose
// Put ignores the context over a fresh subject.
func guardedPutIgnoresCancel(h *GuardedHandle) Guarded {
	return &guardedPutIgnoresCancelDefect{guardedOverEmbedded: h.Fresh()}
}

// guardedPutIgnoresThePredicateDefect is the planted defect of a Guarded
// whose Put accepts a write that Match rejects.
type guardedPutIgnoresThePredicateDefect struct {
	guardedOverEmbedded
	calls int
}

// Put is the override of the planted defect.
func (d *guardedPutIgnoresThePredicateDefect) Put(
	p0 context.Context,
	p1 Meta,
	p2 *Value,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 1 {
		return r0
	}
	return d.guardedOverEmbedded.Put(p0, p1, p2)
}

// guardedPutIgnoresThePredicate returns the planted defect of a Guarded
// whose Put accepts a write that Match rejects over a fresh subject.
func guardedPutIgnoresThePredicate(h *GuardedHandle) Guarded {
	return &guardedPutIgnoresThePredicateDefect{guardedOverEmbedded: h.Fresh()}
}

// guardedMatchPanicsDefect is the planted defect of a Guarded whose Match
// panics.
type guardedMatchPanicsDefect struct {
	guardedOverEmbedded
}

// Match is the override of the planted defect.
func (*guardedMatchPanicsDefect) Match(p0 context.Context, p1 Meta, p2 *Value) (r0 bool, r1 error) {
	panic("the planted defect of a Guarded whose Match panics")
}

// guardedMatchPanics returns the planted defect of a Guarded whose Match
// panics over a fresh subject.
func guardedMatchPanics(h *GuardedHandle) Guarded {
	return &guardedMatchPanicsDefect{guardedOverEmbedded: h.Fresh()}
}

// guardedMatchIgnoresCancelDefect is the planted defect of a Guarded whose
// Match ignores the context.
type guardedMatchIgnoresCancelDefect struct {
	guardedOverEmbedded
}

// Match is the override of the planted defect.
func (d *guardedMatchIgnoresCancelDefect) Match(
	p0 context.Context,
	p1 Meta,
	p2 *Value,
) (r0 bool, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.guardedOverEmbedded.Match(p0, p1, p2)
}

// guardedMatchIgnoresCancel returns the planted defect of a Guarded whose
// Match ignores the context over a fresh subject.
func guardedMatchIgnoresCancel(h *GuardedHandle) Guarded {
	return &guardedMatchIgnoresCancelDefect{guardedOverEmbedded: h.Fresh()}
}

// The derived checks of Scope:
//
// Scope.Run/smoke: Run returns without a panic for derived inputs.
//
// Scope.Run/cancel: Run returns an error for a cancelled context.
//
// Scope.Run/transaction: Get finds nothing that a failed Run wrote through
// Put.
//
// Scope.Put/smoke: Put returns without a panic for derived inputs.
//
// Scope.Put/cancel: Put returns an error for a cancelled context.
//
// Scope.Get/smoke: Get returns without a panic for derived inputs.
//
// Scope.Get/cancel: Get returns an error for a cancelled context.
//
// Scope.Get/zero-on-error: Get returns the zero value beside an error.

// The derived checks of Scope that witness their laws at the shallow end:
//
// Scope.Run/smoke
//
// Scope.Put/smoke
//
// Scope.Get/smoke

// The unclaimed classifications of Scope:
//
// The shape:writer of Scope.Run has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Scope.Put has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Scope writes.
type (
	ScopeHandle = conform.Handle[Scope]
	ScopeCheck  = conform.Check[Scope]
	ScopeOption = conform.Option[Scope]
)

// ScopeChecks contains the name of each callable of Scope, with the
// identities of its checks.
var ScopeChecks = struct {
	Run scopeRun
	Put scopePut
	Get scopeGet
}{
	Run: scopeRun{"Scope.Run"},
	Put: scopePut{"Scope.Put"},
	Get: scopeGet{"Scope.Get"},
}

// scopeRun is the name of Scope.Run, with the identities of its checks.
type scopeRun struct {
	conform.Callable[Scope]
}

// Smoke returns the identity of the check that Run returns without a panic
// for derived inputs.
func (c scopeRun) Smoke() conform.ID[Scope] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Run returns an error for a
// cancelled context.
func (c scopeRun) Cancel() conform.ID[Scope] {
	return c.Check("cancel")
}

// Transaction returns the identity of the check that Get finds nothing
// that a failed Run wrote through Put.
func (c scopeRun) Transaction() conform.ID[Scope] {
	return c.Check("transaction")
}

// scopePut is the name of Scope.Put, with the identities of its checks.
type scopePut struct {
	conform.Callable[Scope]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c scopePut) Smoke() conform.ID[Scope] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c scopePut) Cancel() conform.ID[Scope] {
	return c.Check("cancel")
}

// scopeGet is the name of Scope.Get, with the identities of its checks.
type scopeGet struct {
	conform.Callable[Scope]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c scopeGet) Smoke() conform.ID[Scope] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c scopeGet) Cancel() conform.ID[Scope] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c scopeGet) ZeroOnError() conform.ID[Scope] {
	return c.Check("zero-on-error")
}

// ScopeSuite returns the checks of Scope, and the calls and the budgets of
// its callables.
func ScopeSuite() conform.Suite[Scope] {
	return conform.Suite[Scope]{
		Declaration: "Scope",
		Checks: []conform.Check[Scope]{
			{
				ID:          ScopeChecks.Run.Smoke(),
				Statement:   "Run returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScopeHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, func(ctx context.Context) error { return nil })
					})
				},
				Defect: scopeRunPanics,
			},
			{
				ID:          ScopeChecks.Run.Cancel(),
				Statement:   "Run returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopeHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, func(ctx context.Context) error { return nil })
					})
				},
				Defect: scopeRunIgnoresCancel,
			},
			{
				ID:          ScopeChecks.Run.Transaction(),
				Statement:   "Get finds nothing that a failed Run wrote through Put",
				Claim:       spec.Transaction,
				Probe:       spec.RollsBackOnFailure,
				Reason:      spec.NotRolledBack,
				Law:         spec.FailedScopeLeavesNothing,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopeHandle) {
					probe.RollsBackOnFailure(h, func(ctx context.Context, body probe.Act) error {
						return h.Subject().Run(ctx, body)
					}, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					}, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					}, ErrMissing)
				},
				Defect: scopeGetKeepsFailedWrites,
			},
			{
				ID:          ScopeChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScopeHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: scopePutPanics,
			},
			{
				ID:          ScopeChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopeHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), Value{Key: "Key-a"})
					})
				},
				Defect: scopePutIgnoresCancel,
			},
			{
				ID:          ScopeChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScopeHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: scopeGetPanics,
			},
			{
				ID:          ScopeChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopeHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: scopeGetIgnoresCancel,
			},
			{
				ID:          ScopeChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopeHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: scopeGetValueBeside,
			},
		},
		Calls: []conform.Call[Scope]{
			{
				Calls:    scopeRunCalls,
				Over:     newScopeOver,
				Callable: ScopeChecks.Run.Callable,
			},
			{
				Calls:    scopePutCalls,
				Over:     newScopeOver,
				Callable: ScopeChecks.Put.Callable,
			},
			{
				Calls:    scopeGetCalls,
				Over:     newScopeOver,
				Callable: ScopeChecks.Get.Callable,
			},
		},
	}
}

// RunScopeSuite applies the checks of Scope, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunScopeSuite[T Scope](t *testing.T, newSubject func() T, opts ...ScopeOption) {
	t.Helper()
	RunScopeSuiteFromHandle(t, func(*ScopeHandle) T { return newSubject() }, opts...)
}

// RunScopeSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunScopeSuiteFromHandle[T Scope](
	t *testing.T,
	newSubject func(h *ScopeHandle) T,
	opts ...ScopeOption,
) {
	t.Helper()
	conform.Run(t, ScopeSuite(), scopeSubjects(newSubject), opts...)
}

// FuzzScopeSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzScopeSuite[T Scope](
	f *testing.F,
	callable conform.Callable[Scope],
	newSubject func(h *ScopeHandle) T,
	opts ...ScopeOption,
) {
	f.Helper()
	conform.Fuzz(f, ScopeSuite(), callable, scopeSubjects(newSubject), opts...)
}

// RunScopeBudgets applies the allocation ceilings of Scope and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunScopeBudgets[T Scope](
	t *testing.T,
	newSubject func(h *ScopeHandle) T,
	opts ...ScopeOption,
) {
	t.Helper()
	conform.RunBudgets(t, ScopeSuite(), scopeSubjects(newSubject), opts...)
}

// BenchScopeSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchScopeSuite[T Scope](
	b conform.B,
	callable conform.Callable[Scope],
	newSubject func(h *ScopeHandle) T,
	opts ...ScopeOption,
) {
	b.Helper()
	conform.Bench(b, ScopeSuite(), callable, scopeSubjects(newSubject), opts...)
}

// scopeSubjects returns the factory of the subjects of Scope that calls
// newSubject with the handle of each check.
func scopeSubjects[T Scope](newSubject func(h *ScopeHandle) T) func(h *ScopeHandle) Scope {
	return func(h *ScopeHandle) Scope { return newSubject(h) }
}

// scopeRunCalls generates the calls of Scope.Run with the inputs that each
// case draws.
var scopeRunCalls = prop.Composite(func(c *prop.Case) func(context.Context, Scope) error {
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
	return func(ctx context.Context, subject Scope) error {
		return subject.Run(ctx, body)
	}
})

// scopePutCalls generates the calls of Scope.Put with the inputs that each
// case draws.
var scopePutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Scope) error {
	key := c.Draw(prop.Of[string](), "key")
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Scope) error {
		return subject.Put(ctx, key, v)
	}
})

// scopeGetCalls generates the calls of Scope.Get with the inputs that each
// case draws.
var scopeGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Scope) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Scope) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// scopeOverEmbedded is Scope under the name of the field that scopeOver
// embeds.
type scopeOverEmbedded = Scope

// scopeOver is a subject of Scope whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type scopeOver struct {
	scopeOverEmbedded
	spend func()
}

// Run calls spend and returns its zero results.
func (o scopeOver) Run(context.Context, func(ctx context.Context) error) (r0 error) {
	o.spend()
	return r0
}

// Put calls spend and returns its zero results.
func (o scopeOver) Put(context.Context, string, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o scopeOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newScopeOver returns the planted defect of a ceiling of Scope, whose
// callables call spend.
func newScopeOver(_ *ScopeHandle, spend func()) Scope {
	return scopeOver{spend: spend}
}

// scopeRunPanicsDefect is the planted defect of a Scope whose Run panics.
type scopeRunPanicsDefect struct {
	scopeOverEmbedded
}

// Run is the override of the planted defect.
func (*scopeRunPanicsDefect) Run(
	p0 context.Context,
	p1 func(ctx context.Context) error,
) (r0 error) {
	panic("the planted defect of a Scope whose Run panics")
}

// scopeRunPanics returns the planted defect of a Scope whose Run panics
// over a fresh subject.
func scopeRunPanics(h *ScopeHandle) Scope {
	return &scopeRunPanicsDefect{scopeOverEmbedded: h.Fresh()}
}

// scopeRunIgnoresCancelDefect is the planted defect of a Scope whose Run
// ignores the context.
type scopeRunIgnoresCancelDefect struct {
	scopeOverEmbedded
}

// Run is the override of the planted defect.
func (d *scopeRunIgnoresCancelDefect) Run(
	p0 context.Context,
	p1 func(ctx context.Context) error,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.scopeOverEmbedded.Run(p0, p1)
}

// scopeRunIgnoresCancel returns the planted defect of a Scope whose Run
// ignores the context over a fresh subject.
func scopeRunIgnoresCancel(h *ScopeHandle) Scope {
	return &scopeRunIgnoresCancelDefect{scopeOverEmbedded: h.Fresh()}
}

// scopeGetKeepsFailedWritesDefect is the planted defect of a Scope whose
// Run keeps what a failed body wrote.
type scopeGetKeepsFailedWritesDefect struct {
	scopeOverEmbedded
}

// Get is the override of the planted defect.
func (*scopeGetKeepsFailedWritesDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	return r0, r1
}

// scopeGetKeepsFailedWrites returns the planted defect of a Scope whose
// Run keeps what a failed body wrote over a fresh subject.
func scopeGetKeepsFailedWrites(h *ScopeHandle) Scope {
	return &scopeGetKeepsFailedWritesDefect{scopeOverEmbedded: h.Fresh()}
}

// scopePutPanicsDefect is the planted defect of a Scope whose Put panics.
type scopePutPanicsDefect struct {
	scopeOverEmbedded
}

// Put is the override of the planted defect.
func (*scopePutPanicsDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	panic("the planted defect of a Scope whose Put panics")
}

// scopePutPanics returns the planted defect of a Scope whose Put panics
// over a fresh subject.
func scopePutPanics(h *ScopeHandle) Scope {
	return &scopePutPanicsDefect{scopeOverEmbedded: h.Fresh()}
}

// scopePutIgnoresCancelDefect is the planted defect of a Scope whose Put
// ignores the context.
type scopePutIgnoresCancelDefect struct {
	scopeOverEmbedded
}

// Put is the override of the planted defect.
func (d *scopePutIgnoresCancelDefect) Put(p0 context.Context, p1 string, p2 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.scopeOverEmbedded.Put(p0, p1, p2)
}

// scopePutIgnoresCancel returns the planted defect of a Scope whose Put
// ignores the context over a fresh subject.
func scopePutIgnoresCancel(h *ScopeHandle) Scope {
	return &scopePutIgnoresCancelDefect{scopeOverEmbedded: h.Fresh()}
}

// scopeGetPanicsDefect is the planted defect of a Scope whose Get panics.
type scopeGetPanicsDefect struct {
	scopeOverEmbedded
}

// Get is the override of the planted defect.
func (*scopeGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Scope whose Get panics")
}

// scopeGetPanics returns the planted defect of a Scope whose Get panics
// over a fresh subject.
func scopeGetPanics(h *ScopeHandle) Scope {
	return &scopeGetPanicsDefect{scopeOverEmbedded: h.Fresh()}
}

// scopeGetIgnoresCancelDefect is the planted defect of a Scope whose Get
// ignores the context.
type scopeGetIgnoresCancelDefect struct {
	scopeOverEmbedded
}

// Get is the override of the planted defect.
func (d *scopeGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.scopeOverEmbedded.Get(p0, p1)
}

// scopeGetIgnoresCancel returns the planted defect of a Scope whose Get
// ignores the context over a fresh subject.
func scopeGetIgnoresCancel(h *ScopeHandle) Scope {
	return &scopeGetIgnoresCancelDefect{scopeOverEmbedded: h.Fresh()}
}

// scopeGetValueBesideDefect is the planted defect of a Scope whose Get
// returns a value beside an error.
type scopeGetValueBesideDefect struct {
	scopeOverEmbedded
}

// Get is the override of the planted defect.
func (d *scopeGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.scopeOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// scopeGetValueBeside returns the planted defect of a Scope whose Get
// returns a value beside an error over a fresh subject.
func scopeGetValueBeside(h *ScopeHandle) Scope {
	return &scopeGetValueBesideDefect{scopeOverEmbedded: h.Fresh()}
}

// The derived checks of Events:
//
// Events.Fire/smoke: Fire returns without a panic for derived inputs.
//
// Events.Fire/cancel: Fire returns an error for a cancelled context.
//
// Events.Fire/hooks: Fire runs a hook that OnEvent installed.
//
// Events.OnEvent/smoke: OnEvent returns without a panic for derived
// inputs.

// The derived checks of Events that witness their laws at the shallow end:
//
// Events.Fire/smoke
//
// Events.OnEvent/smoke

// The unclaimed classifications of Events:
//
// The shape:lifecycle of Events.Fire has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The shape:mutator of Events.OnEvent has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Events writes.
type (
	EventsHandle = conform.Handle[Events]
	EventsCheck  = conform.Check[Events]
	EventsOption = conform.Option[Events]
)

// EventsChecks contains the name of each callable of Events, with the
// identities of its checks.
var EventsChecks = struct {
	Fire    eventsFire
	OnEvent eventsOnEvent
}{
	Fire:    eventsFire{"Events.Fire"},
	OnEvent: eventsOnEvent{"Events.OnEvent"},
}

// eventsFire is the name of Events.Fire, with the identities of its
// checks.
type eventsFire struct {
	conform.Callable[Events]
}

// Smoke returns the identity of the check that Fire returns without a
// panic for derived inputs.
func (c eventsFire) Smoke() conform.ID[Events] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Fire returns an error for
// a cancelled context.
func (c eventsFire) Cancel() conform.ID[Events] {
	return c.Check("cancel")
}

// Hooks returns the identity of the check that Fire runs a hook that
// OnEvent installed.
func (c eventsFire) Hooks() conform.ID[Events] {
	return c.Check("hooks")
}

// eventsOnEvent is the name of Events.OnEvent, with the identities of its
// checks.
type eventsOnEvent struct {
	conform.Callable[Events]
}

// Smoke returns the identity of the check that OnEvent returns without a
// panic for derived inputs.
func (c eventsOnEvent) Smoke() conform.ID[Events] {
	return c.Check("smoke")
}

// EventsSuite returns the checks of Events, and the calls and the budgets
// of its callables.
func EventsSuite() conform.Suite[Events] {
	return conform.Suite[Events]{
		Declaration: "Events",
		Checks: []conform.Check[Events]{
			{
				ID:          EventsChecks.Fire.Smoke(),
				Statement:   "Fire returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *EventsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Fire(ctx)
					})
				},
				Defect: eventsFirePanics,
			},
			{
				ID:          EventsChecks.Fire.Cancel(),
				Statement:   "Fire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *EventsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Fire(ctx)
					})
				},
				Defect: eventsFireIgnoresCancel,
			},
			{
				ID:          EventsChecks.Fire.Hooks(),
				Statement:   "Fire runs a hook that OnEvent installed",
				Claim:       spec.Hooks,
				Probe:       spec.RunsTheHook,
				Reason:      spec.HookNotRun,
				Law:         spec.RegisteredHookRuns,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *EventsHandle) {
					probe.RunsTheHook(h, func(ctx context.Context) error {
						return h.Subject().Fire(ctx)
					}, func(ctx context.Context, hook func()) error {
						h.Subject().OnEvent(func(_ string) { hook() })
						return nil
					})
				},
				Defect: eventsOnEventSwallowsHooks,
			},
			{
				ID:          EventsChecks.OnEvent.Smoke(),
				Statement:   "OnEvent returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *EventsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().OnEvent(func(event string) {})
						return nil
					})
				},
				Defect: eventsOnEventPanics,
			},
		},
		Calls: []conform.Call[Events]{
			{
				Calls:    eventsFireCalls,
				Over:     newEventsOver,
				Callable: EventsChecks.Fire.Callable,
			},
			{
				Calls:    eventsOnEventCalls,
				Over:     newEventsOver,
				Callable: EventsChecks.OnEvent.Callable,
			},
		},
	}
}

// RunEventsSuite applies the checks of Events, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunEventsSuite[T Events](t *testing.T, newSubject func() T, opts ...EventsOption) {
	t.Helper()
	RunEventsSuiteFromHandle(t, func(*EventsHandle) T { return newSubject() }, opts...)
}

// RunEventsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunEventsSuiteFromHandle[T Events](
	t *testing.T,
	newSubject func(h *EventsHandle) T,
	opts ...EventsOption,
) {
	t.Helper()
	conform.Run(t, EventsSuite(), eventsSubjects(newSubject), opts...)
}

// FuzzEventsSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzEventsSuite[T Events](
	f *testing.F,
	callable conform.Callable[Events],
	newSubject func(h *EventsHandle) T,
	opts ...EventsOption,
) {
	f.Helper()
	conform.Fuzz(f, EventsSuite(), callable, eventsSubjects(newSubject), opts...)
}

// RunEventsBudgets applies the allocation ceilings of Events and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunEventsBudgets[T Events](
	t *testing.T,
	newSubject func(h *EventsHandle) T,
	opts ...EventsOption,
) {
	t.Helper()
	conform.RunBudgets(t, EventsSuite(), eventsSubjects(newSubject), opts...)
}

// BenchEventsSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchEventsSuite[T Events](
	b conform.B,
	callable conform.Callable[Events],
	newSubject func(h *EventsHandle) T,
	opts ...EventsOption,
) {
	b.Helper()
	conform.Bench(b, EventsSuite(), callable, eventsSubjects(newSubject), opts...)
}

// eventsSubjects returns the factory of the subjects of Events that calls
// newSubject with the handle of each check.
func eventsSubjects[T Events](newSubject func(h *EventsHandle) T) func(h *EventsHandle) Events {
	return func(h *EventsHandle) Events { return newSubject(h) }
}

// eventsFireCalls generates the calls of Events.Fire with the inputs that
// each case draws.
var eventsFireCalls = prop.Composite(func(*prop.Case) func(context.Context, Events) error {
	return func(ctx context.Context, subject Events) error {
		return subject.Fire(ctx)
	}
})

// eventsOnEventCalls generates the calls of Events.OnEvent with the inputs
// that each case draws.
var eventsOnEventCalls = prop.Composite(func(*prop.Case) func(context.Context, Events) error {
	fn := func(string) {}
	return func(_ context.Context, subject Events) error {
		subject.OnEvent(fn)
		return nil
	}
})

// eventsOverEmbedded is Events under the name of the field that eventsOver
// embeds.
type eventsOverEmbedded = Events

// eventsOver is a subject of Events whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type eventsOver struct {
	eventsOverEmbedded
	spend func()
}

// Fire calls spend and returns its zero results.
func (o eventsOver) Fire(context.Context) (r0 error) {
	o.spend()
	return r0
}

// OnEvent calls spend and returns its zero results.
func (o eventsOver) OnEvent(func(event string)) {
	o.spend()
}

// newEventsOver returns the planted defect of a ceiling of Events, whose
// callables call spend.
func newEventsOver(_ *EventsHandle, spend func()) Events {
	return eventsOver{spend: spend}
}

// eventsFirePanicsDefect is the planted defect of a Events whose Fire
// panics.
type eventsFirePanicsDefect struct {
	eventsOverEmbedded
}

// Fire is the override of the planted defect.
func (*eventsFirePanicsDefect) Fire(p0 context.Context) (r0 error) {
	panic("the planted defect of a Events whose Fire panics")
}

// eventsFirePanics returns the planted defect of a Events whose Fire
// panics over a fresh subject.
func eventsFirePanics(h *EventsHandle) Events {
	return &eventsFirePanicsDefect{eventsOverEmbedded: h.Fresh()}
}

// eventsFireIgnoresCancelDefect is the planted defect of a Events whose
// Fire ignores the context.
type eventsFireIgnoresCancelDefect struct {
	eventsOverEmbedded
}

// Fire is the override of the planted defect.
func (d *eventsFireIgnoresCancelDefect) Fire(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.eventsOverEmbedded.Fire(p0)
}

// eventsFireIgnoresCancel returns the planted defect of a Events whose
// Fire ignores the context over a fresh subject.
func eventsFireIgnoresCancel(h *EventsHandle) Events {
	return &eventsFireIgnoresCancelDefect{eventsOverEmbedded: h.Fresh()}
}

// eventsOnEventSwallowsHooksDefect is the planted defect of a Events whose
// OnEvent never installs the hook.
type eventsOnEventSwallowsHooksDefect struct {
	eventsOverEmbedded
}

// OnEvent is the override of the planted defect.
func (*eventsOnEventSwallowsHooksDefect) OnEvent(p0 func(event string)) {
	return
}

// eventsOnEventSwallowsHooks returns the planted defect of a Events whose
// OnEvent never installs the hook over a fresh subject.
func eventsOnEventSwallowsHooks(h *EventsHandle) Events {
	return &eventsOnEventSwallowsHooksDefect{eventsOverEmbedded: h.Fresh()}
}

// eventsOnEventPanicsDefect is the planted defect of a Events whose
// OnEvent panics.
type eventsOnEventPanicsDefect struct {
	eventsOverEmbedded
}

// OnEvent is the override of the planted defect.
func (*eventsOnEventPanicsDefect) OnEvent(p0 func(event string)) {
	panic("the planted defect of a Events whose OnEvent panics")
}

// eventsOnEventPanics returns the planted defect of a Events whose OnEvent
// panics over a fresh subject.
func eventsOnEventPanics(h *EventsHandle) Events {
	return &eventsOnEventPanicsDefect{eventsOverEmbedded: h.Fresh()}
}

// The derived checks of Timers:
//
// Timers.At/smoke: At returns without a panic for derived inputs.
//
// Timers.At/cancel: At returns an error for a cancelled context.
//
// Timers.At/scheduled: a task that At registers runs once the clock passes
// its instant, and not before.
//
// Timers.Fired/smoke: Fired returns without a panic for derived inputs.
//
// Timers.Fired/cancel: Fired returns an error for a cancelled context.
//
// Timers.Fired/zero-on-error: Fired returns the zero value beside an
// error.
//
// Timers.Slow/smoke: Slow returns without a panic for derived inputs.
//
// Timers.Slow/cancel: Slow returns an error for a cancelled context.
//
// Timers.Slow/within: Slow returns before the clock passes the declared
// budget.

// The derived checks of Timers that witness their laws at the shallow end:
//
// Timers.At/smoke
//
// Timers.Fired/smoke
//
// Timers.Slow/smoke

// The unclaimed classifications of Timers:
//
// The shape:writer of Timers.At has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:writer of Timers.Slow has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Timers writes.
type (
	TimersHandle = conform.Handle[Timers]
	TimersCheck  = conform.Check[Timers]
	TimersOption = conform.Option[Timers]
)

// TimersChecks contains the name of each callable of Timers, with the
// identities of its checks.
var TimersChecks = struct {
	At    timersAt
	Fired timersFired
	Slow  timersSlow
}{
	At:    timersAt{"Timers.At"},
	Fired: timersFired{"Timers.Fired"},
	Slow:  timersSlow{"Timers.Slow"},
}

// timersAt is the name of Timers.At, with the identities of its checks.
type timersAt struct {
	conform.Callable[Timers]
}

// Smoke returns the identity of the check that At returns without a panic
// for derived inputs.
func (c timersAt) Smoke() conform.ID[Timers] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that At returns an error for a
// cancelled context.
func (c timersAt) Cancel() conform.ID[Timers] {
	return c.Check("cancel")
}

// Scheduled returns the identity of the check that a task that At
// registers runs once the clock passes its instant, and not before.
func (c timersAt) Scheduled() conform.ID[Timers] {
	return c.Check("scheduled")
}

// timersFired is the name of Timers.Fired, with the identities of its
// checks.
type timersFired struct {
	conform.Callable[Timers]
}

// Smoke returns the identity of the check that Fired returns without a
// panic for derived inputs.
func (c timersFired) Smoke() conform.ID[Timers] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Fired returns an error for
// a cancelled context.
func (c timersFired) Cancel() conform.ID[Timers] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Fired returns the
// zero value beside an error.
func (c timersFired) ZeroOnError() conform.ID[Timers] {
	return c.Check("zero-on-error")
}

// timersSlow is the name of Timers.Slow, with the identities of its
// checks.
type timersSlow struct {
	conform.Callable[Timers]
}

// Smoke returns the identity of the check that Slow returns without a
// panic for derived inputs.
func (c timersSlow) Smoke() conform.ID[Timers] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Slow returns an error for
// a cancelled context.
func (c timersSlow) Cancel() conform.ID[Timers] {
	return c.Check("cancel")
}

// Within returns the identity of the check that Slow returns before the
// clock passes the declared budget.
func (c timersSlow) Within() conform.ID[Timers] {
	return c.Check("within")
}

// TimersSuite returns the checks of Timers, and the calls and the budgets
// of its callables.
func TimersSuite() conform.Suite[Timers] {
	return conform.Suite[Timers]{
		Declaration: "Timers",
		Checks: []conform.Check[Timers]{
			{
				ID:          TimersChecks.At.Smoke(),
				Statement:   "At returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TimersHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().At(ctx, time.Duration(1))
					})
				},
				Defect: timersAtPanics,
			},
			{
				ID:          TimersChecks.At.Cancel(),
				Statement:   "At returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TimersHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().At(ctx, time.Duration(1))
					})
				},
				Defect: timersAtIgnoresCancel,
			},
			{
				ID:          TimersChecks.At.Scheduled(),
				Statement:   "a task that At registers runs once the clock passes its instant, and not before",
				Claim:       spec.Scheduled,
				Probe:       spec.RunsOnSchedule,
				Reason:      spec.OffSchedule,
				Law:         spec.TaskRunsOnceItsInstantPasses,
				Machineries: []spec.MachineryID{spec.ClockTheRunAdvances},
				Apply: func(h *TimersHandle) {
					probe.RunsOnSchedule(h, func(ctx context.Context) (time.Duration, error) {
						d := time.Duration(1)
						err := h.Subject().At(ctx, d)
						return d, err
					}, func(ctx context.Context) (int, error) {
						return h.Subject().Fired(ctx, time.Duration(1))
					})
				},
				Defect: timersFiredCountDrifts,
				Around: synctest.Test,
			},
			{
				ID:          TimersChecks.Fired.Smoke(),
				Statement:   "Fired returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TimersHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fired(ctx, time.Duration(1))
						return err
					})
				},
				Defect: timersFiredPanics,
			},
			{
				ID:          TimersChecks.Fired.Cancel(),
				Statement:   "Fired returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TimersHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fired(ctx, time.Duration(1))
						return err
					})
				},
				Defect: timersFiredIgnoresCancel,
			},
			{
				ID:          TimersChecks.Fired.ZeroOnError(),
				Statement:   "Fired returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TimersHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Fired(ctx, time.Duration(1))
					}, int(0))
				},
				Defect: timersFiredValueBeside,
			},
			{
				ID:          TimersChecks.Slow.Smoke(),
				Statement:   "Slow returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TimersHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Slow(ctx, string("key-a"))
					})
				},
				Defect: timersSlowPanics,
			},
			{
				ID:          TimersChecks.Slow.Cancel(),
				Statement:   "Slow returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TimersHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Slow(ctx, string("key-a"))
					})
				},
				Defect: timersSlowIgnoresCancel,
			},
			{
				ID:          TimersChecks.Slow.Within(),
				Statement:   "Slow returns before the clock passes the declared budget",
				Claim:       spec.Within,
				Probe:       spec.ReturnsWithin,
				Reason:      spec.OverranBudget,
				Law:         spec.CallReturnsWithinItsBudget,
				Machineries: []spec.MachineryID{spec.ClockTheRunAdvances},
				Apply: func(h *TimersHandle) {
					probe.ReturnsWithin(h, func(ctx context.Context) error {
						return h.Subject().Slow(ctx, string("key-a"))
					}, 5*time.Second)
				},
				Defect: timersSlowSlow,
				Around: synctest.Test,
			},
		},
		Calls: []conform.Call[Timers]{
			{
				Calls:    timersAtCalls,
				Over:     newTimersOver,
				Callable: TimersChecks.At.Callable,
			},
			{
				Calls:    timersFiredCalls,
				Over:     newTimersOver,
				Callable: TimersChecks.Fired.Callable,
			},
			{
				Calls:    timersSlowCalls,
				Over:     newTimersOver,
				Callable: TimersChecks.Slow.Callable,
			},
		},
	}
}

// RunTimersSuite applies the checks of Timers, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunTimersSuite[T Timers](t *testing.T, newSubject func() T, opts ...TimersOption) {
	t.Helper()
	RunTimersSuiteFromHandle(t, func(*TimersHandle) T { return newSubject() }, opts...)
}

// RunTimersSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTimersSuiteFromHandle[T Timers](
	t *testing.T,
	newSubject func(h *TimersHandle) T,
	opts ...TimersOption,
) {
	t.Helper()
	conform.Run(t, TimersSuite(), timersSubjects(newSubject), opts...)
}

// FuzzTimersSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTimersSuite[T Timers](
	f *testing.F,
	callable conform.Callable[Timers],
	newSubject func(h *TimersHandle) T,
	opts ...TimersOption,
) {
	f.Helper()
	conform.Fuzz(f, TimersSuite(), callable, timersSubjects(newSubject), opts...)
}

// RunTimersBudgets applies the allocation ceilings of Timers and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunTimersBudgets[T Timers](
	t *testing.T,
	newSubject func(h *TimersHandle) T,
	opts ...TimersOption,
) {
	t.Helper()
	conform.RunBudgets(t, TimersSuite(), timersSubjects(newSubject), opts...)
}

// BenchTimersSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchTimersSuite[T Timers](
	b conform.B,
	callable conform.Callable[Timers],
	newSubject func(h *TimersHandle) T,
	opts ...TimersOption,
) {
	b.Helper()
	conform.Bench(b, TimersSuite(), callable, timersSubjects(newSubject), opts...)
}

// timersSubjects returns the factory of the subjects of Timers that calls
// newSubject with the handle of each check.
func timersSubjects[T Timers](newSubject func(h *TimersHandle) T) func(h *TimersHandle) Timers {
	return func(h *TimersHandle) Timers { return newSubject(h) }
}

// timersAtCalls generates the calls of Timers.At with the inputs that each
// case draws.
var timersAtCalls = prop.Composite(func(c *prop.Case) func(context.Context, Timers) error {
	after := c.Draw(prop.Of[time.Duration](), "after")
	return func(ctx context.Context, subject Timers) error {
		return subject.At(ctx, after)
	}
})

// timersFiredCalls generates the calls of Timers.Fired with the inputs
// that each case draws.
var timersFiredCalls = prop.Composite(func(c *prop.Case) func(context.Context, Timers) error {
	after := c.Draw(prop.Of[time.Duration](), "after")
	return func(ctx context.Context, subject Timers) error {
		_, err := subject.Fired(ctx, after)
		return err
	}
})

// timersSlowCalls generates the calls of Timers.Slow with the inputs that
// each case draws.
var timersSlowCalls = prop.Composite(func(c *prop.Case) func(context.Context, Timers) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Timers) error {
		return subject.Slow(ctx, key)
	}
})

// timersOverEmbedded is Timers under the name of the field that timersOver
// embeds.
type timersOverEmbedded = Timers

// timersOver is a subject of Timers whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type timersOver struct {
	timersOverEmbedded
	spend func()
}

// At calls spend and returns its zero results.
func (o timersOver) At(context.Context, time.Duration) (r0 error) {
	o.spend()
	return r0
}

// Fired calls spend and returns its zero results.
func (o timersOver) Fired(context.Context, time.Duration) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// Slow calls spend and returns its zero results.
func (o timersOver) Slow(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// newTimersOver returns the planted defect of a ceiling of Timers, whose
// callables call spend.
func newTimersOver(_ *TimersHandle, spend func()) Timers {
	return timersOver{spend: spend}
}

// timersAtPanicsDefect is the planted defect of a Timers whose At panics.
type timersAtPanicsDefect struct {
	timersOverEmbedded
}

// At is the override of the planted defect.
func (*timersAtPanicsDefect) At(p0 context.Context, p1 time.Duration) (r0 error) {
	panic("the planted defect of a Timers whose At panics")
}

// timersAtPanics returns the planted defect of a Timers whose At panics
// over a fresh subject.
func timersAtPanics(h *TimersHandle) Timers {
	return &timersAtPanicsDefect{timersOverEmbedded: h.Fresh()}
}

// timersAtIgnoresCancelDefect is the planted defect of a Timers whose At
// ignores the context.
type timersAtIgnoresCancelDefect struct {
	timersOverEmbedded
}

// At is the override of the planted defect.
func (d *timersAtIgnoresCancelDefect) At(p0 context.Context, p1 time.Duration) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.timersOverEmbedded.At(p0, p1)
}

// timersAtIgnoresCancel returns the planted defect of a Timers whose At
// ignores the context over a fresh subject.
func timersAtIgnoresCancel(h *TimersHandle) Timers {
	return &timersAtIgnoresCancelDefect{timersOverEmbedded: h.Fresh()}
}

// timersFiredCountDriftsDefect is the planted defect of a Timers that runs
// each task of At the moment it registers.
type timersFiredCountDriftsDefect struct {
	timersOverEmbedded
	calls int
}

// Fired is the override of the planted defect.
func (d *timersFiredCountDriftsDefect) Fired(
	p0 context.Context,
	p1 time.Duration,
) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0, r1
}

// timersFiredCountDrifts returns the planted defect of a Timers that runs
// each task of At the moment it registers over a fresh subject.
func timersFiredCountDrifts(h *TimersHandle) Timers {
	return &timersFiredCountDriftsDefect{timersOverEmbedded: h.Fresh()}
}

// timersFiredPanicsDefect is the planted defect of a Timers whose Fired
// panics.
type timersFiredPanicsDefect struct {
	timersOverEmbedded
}

// Fired is the override of the planted defect.
func (*timersFiredPanicsDefect) Fired(p0 context.Context, p1 time.Duration) (r0 int, r1 error) {
	panic("the planted defect of a Timers whose Fired panics")
}

// timersFiredPanics returns the planted defect of a Timers whose Fired
// panics over a fresh subject.
func timersFiredPanics(h *TimersHandle) Timers {
	return &timersFiredPanicsDefect{timersOverEmbedded: h.Fresh()}
}

// timersFiredIgnoresCancelDefect is the planted defect of a Timers whose
// Fired ignores the context.
type timersFiredIgnoresCancelDefect struct {
	timersOverEmbedded
}

// Fired is the override of the planted defect.
func (d *timersFiredIgnoresCancelDefect) Fired(
	p0 context.Context,
	p1 time.Duration,
) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.timersOverEmbedded.Fired(p0, p1)
}

// timersFiredIgnoresCancel returns the planted defect of a Timers whose
// Fired ignores the context over a fresh subject.
func timersFiredIgnoresCancel(h *TimersHandle) Timers {
	return &timersFiredIgnoresCancelDefect{timersOverEmbedded: h.Fresh()}
}

// timersFiredValueBesideDefect is the planted defect of a Timers whose
// Fired returns a value beside an error.
type timersFiredValueBesideDefect struct {
	timersOverEmbedded
}

// Fired is the override of the planted defect.
func (d *timersFiredValueBesideDefect) Fired(
	p0 context.Context,
	p1 time.Duration,
) (r0 int, r1 error) {
	r0, r1 = d.timersOverEmbedded.Fired(p0, p1)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// timersFiredValueBeside returns the planted defect of a Timers whose
// Fired returns a value beside an error over a fresh subject.
func timersFiredValueBeside(h *TimersHandle) Timers {
	return &timersFiredValueBesideDefect{timersOverEmbedded: h.Fresh()}
}

// timersSlowPanicsDefect is the planted defect of a Timers whose Slow
// panics.
type timersSlowPanicsDefect struct {
	timersOverEmbedded
}

// Slow is the override of the planted defect.
func (*timersSlowPanicsDefect) Slow(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Timers whose Slow panics")
}

// timersSlowPanics returns the planted defect of a Timers whose Slow
// panics over a fresh subject.
func timersSlowPanics(h *TimersHandle) Timers {
	return &timersSlowPanicsDefect{timersOverEmbedded: h.Fresh()}
}

// timersSlowIgnoresCancelDefect is the planted defect of a Timers whose
// Slow ignores the context.
type timersSlowIgnoresCancelDefect struct {
	timersOverEmbedded
}

// Slow is the override of the planted defect.
func (d *timersSlowIgnoresCancelDefect) Slow(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.timersOverEmbedded.Slow(p0, p1)
}

// timersSlowIgnoresCancel returns the planted defect of a Timers whose
// Slow ignores the context over a fresh subject.
func timersSlowIgnoresCancel(h *TimersHandle) Timers {
	return &timersSlowIgnoresCancelDefect{timersOverEmbedded: h.Fresh()}
}

// timersSlowSlowDefect is the planted defect of a Timers whose Slow waits
// longer than its budget.
type timersSlowSlowDefect struct {
	timersOverEmbedded
	clock assert.Clock
}

// Slow is the override of the planted defect.
func (d *timersSlowSlowDefect) Slow(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.clock.Sleep(10 * time.Second)
	return d.timersOverEmbedded.Slow(p0, p1)
}

// timersSlowSlow returns the planted defect of a Timers whose Slow waits
// longer than its budget over a fresh subject.
func timersSlowSlow(h *TimersHandle) Timers {
	return &timersSlowSlowDefect{timersOverEmbedded: h.Fresh(), clock: h.Clock()}
}

// The derived checks of Election:
//
// Election.Campaign/smoke: Campaign returns without a panic for derived
// inputs.
//
// Election.Campaign/cancel: Campaign returns an error for a cancelled
// context.
//
// Election.Campaign/leader-election: IsLeader reports no leadership that
// Campaign took once Resign has run.
//
// Election.Resign/smoke: Resign returns without a panic for derived
// inputs.
//
// Election.Resign/cancel: Resign returns an error for a cancelled context.
//
// Election.IsLeader/smoke: IsLeader returns without a panic for derived
// inputs.

// The derived checks of Election that witness their laws at the shallow
// end:
//
// Election.Campaign/smoke
//
// Election.Resign/smoke
//
// Election.IsLeader/smoke

// The refused claims of Election:
//
// Election.IsLeader/cancel: IsLeader cannot fail. The remedy is to give
// IsLeader an error result.
//
// Election.IsLeader/zero-on-error: IsLeader cannot fail. The remedy is to
// give IsLeader an error result.

// The unclaimed classifications of Election:
//
// The shape:lifecycle of Election.Campaign has the stance silent. A
// lifecycle call returns nothing that a check could compare.
//
// The shape:lifecycle of Election.Resign has the stance silent. A
// lifecycle call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Election writes.
type (
	ElectionHandle = conform.Handle[Election]
	ElectionCheck  = conform.Check[Election]
	ElectionOption = conform.Option[Election]
)

// ElectionChecks contains the name of each callable of Election, with the
// identities of its checks.
var ElectionChecks = struct {
	Campaign electionCampaign
	Resign   electionResign
	IsLeader electionIsLeader
}{
	Campaign: electionCampaign{"Election.Campaign"},
	Resign:   electionResign{"Election.Resign"},
	IsLeader: electionIsLeader{"Election.IsLeader"},
}

// electionCampaign is the name of Election.Campaign, with the identities
// of its checks.
type electionCampaign struct {
	conform.Callable[Election]
}

// Smoke returns the identity of the check that Campaign returns without a
// panic for derived inputs.
func (c electionCampaign) Smoke() conform.ID[Election] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Campaign returns an error
// for a cancelled context.
func (c electionCampaign) Cancel() conform.ID[Election] {
	return c.Check("cancel")
}

// LeaderElection returns the identity of the check that IsLeader reports
// no leadership that Campaign took once Resign has run.
func (c electionCampaign) LeaderElection() conform.ID[Election] {
	return c.Check("leader-election")
}

// electionResign is the name of Election.Resign, with the identities of
// its checks.
type electionResign struct {
	conform.Callable[Election]
}

// Smoke returns the identity of the check that Resign returns without a
// panic for derived inputs.
func (c electionResign) Smoke() conform.ID[Election] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Resign returns an error
// for a cancelled context.
func (c electionResign) Cancel() conform.ID[Election] {
	return c.Check("cancel")
}

// electionIsLeader is the name of Election.IsLeader, with the identities
// of its checks.
type electionIsLeader struct {
	conform.Callable[Election]
}

// Smoke returns the identity of the check that IsLeader returns without a
// panic for derived inputs.
func (c electionIsLeader) Smoke() conform.ID[Election] {
	return c.Check("smoke")
}

// ElectionSuite returns the checks of Election, and the calls and the
// budgets of its callables.
func ElectionSuite() conform.Suite[Election] {
	return conform.Suite[Election]{
		Declaration: "Election",
		Checks: []conform.Check[Election]{
			{
				ID:          ElectionChecks.Campaign.Smoke(),
				Statement:   "Campaign returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ElectionHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Campaign(ctx)
					})
				},
				Defect: electionCampaignPanics,
			},
			{
				ID:          ElectionChecks.Campaign.Cancel(),
				Statement:   "Campaign returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ElectionHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Campaign(ctx)
					})
				},
				Defect: electionCampaignIgnoresCancel,
			},
			{
				ID:          ElectionChecks.Campaign.LeaderElection(),
				Statement:   "IsLeader reports no leadership that Campaign took once Resign has run",
				Claim:       spec.LeaderElection,
				Probe:       spec.ResignsLeadership,
				Reason:      spec.StillLeading,
				Law:         spec.ResignedLeadershipIsReleased,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ElectionHandle) {
					probe.ResignsLeadership(h, func(ctx context.Context) error {
						return h.Subject().Campaign(ctx)
					}, func(ctx context.Context) error {
						return h.Subject().Resign(ctx)
					}, func(ctx context.Context) (bool, error) {
						r0 := h.Subject().IsLeader(ctx)
						return r0, nil
					})
				},
				Defect: electionIsLeaderLeadsOn,
			},
			{
				ID:          ElectionChecks.Resign.Smoke(),
				Statement:   "Resign returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ElectionHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Resign(ctx)
					})
				},
				Defect: electionResignPanics,
			},
			{
				ID:          ElectionChecks.Resign.Cancel(),
				Statement:   "Resign returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ElectionHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Resign(ctx)
					})
				},
				Defect: electionResignIgnoresCancel,
			},
			{
				ID:          ElectionChecks.IsLeader.Smoke(),
				Statement:   "IsLeader returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ElectionHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().IsLeader(ctx)
						return nil
					})
				},
				Defect: electionIsLeaderPanics,
			},
		},
		Calls: []conform.Call[Election]{
			{
				Calls:    electionCampaignCalls,
				Over:     newElectionOver,
				Callable: ElectionChecks.Campaign.Callable,
			},
			{
				Calls:    electionResignCalls,
				Over:     newElectionOver,
				Callable: ElectionChecks.Resign.Callable,
			},
			{
				Calls:    electionIsLeaderCalls,
				Over:     newElectionOver,
				Callable: ElectionChecks.IsLeader.Callable,
			},
		},
	}
}

// RunElectionSuite applies the checks of Election, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunElectionSuite[T Election](t *testing.T, newSubject func() T, opts ...ElectionOption) {
	t.Helper()
	RunElectionSuiteFromHandle(t, func(*ElectionHandle) T { return newSubject() }, opts...)
}

// RunElectionSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunElectionSuiteFromHandle[T Election](
	t *testing.T,
	newSubject func(h *ElectionHandle) T,
	opts ...ElectionOption,
) {
	t.Helper()
	conform.Run(t, ElectionSuite(), electionSubjects(newSubject), opts...)
}

// FuzzElectionSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzElectionSuite[T Election](
	f *testing.F,
	callable conform.Callable[Election],
	newSubject func(h *ElectionHandle) T,
	opts ...ElectionOption,
) {
	f.Helper()
	conform.Fuzz(f, ElectionSuite(), callable, electionSubjects(newSubject), opts...)
}

// RunElectionBudgets applies the allocation ceilings of Election and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunElectionBudgets[T Election](
	t *testing.T,
	newSubject func(h *ElectionHandle) T,
	opts ...ElectionOption,
) {
	t.Helper()
	conform.RunBudgets(t, ElectionSuite(), electionSubjects(newSubject), opts...)
}

// BenchElectionSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchElectionSuite[T Election](
	b conform.B,
	callable conform.Callable[Election],
	newSubject func(h *ElectionHandle) T,
	opts ...ElectionOption,
) {
	b.Helper()
	conform.Bench(b, ElectionSuite(), callable, electionSubjects(newSubject), opts...)
}

// electionSubjects returns the factory of the subjects of Election that
// calls newSubject with the handle of each check.
func electionSubjects[T Election](
	newSubject func(h *ElectionHandle) T,
) func(h *ElectionHandle) Election {
	return func(h *ElectionHandle) Election { return newSubject(h) }
}

// electionCampaignCalls generates the calls of Election.Campaign with the
// inputs that each case draws.
var electionCampaignCalls = prop.Composite(func(*prop.Case) func(context.Context, Election) error {
	return func(ctx context.Context, subject Election) error {
		return subject.Campaign(ctx)
	}
})

// electionResignCalls generates the calls of Election.Resign with the
// inputs that each case draws.
var electionResignCalls = prop.Composite(func(*prop.Case) func(context.Context, Election) error {
	return func(ctx context.Context, subject Election) error {
		return subject.Resign(ctx)
	}
})

// electionIsLeaderCalls generates the calls of Election.IsLeader with the
// inputs that each case draws.
var electionIsLeaderCalls = prop.Composite(func(*prop.Case) func(context.Context, Election) error {
	return func(ctx context.Context, subject Election) error {
		_ = subject.IsLeader(ctx)
		return nil
	}
})

// electionOverEmbedded is Election under the name of the field that
// electionOver embeds.
type electionOverEmbedded = Election

// electionOver is a subject of Election whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type electionOver struct {
	electionOverEmbedded
	spend func()
}

// Campaign calls spend and returns its zero results.
func (o electionOver) Campaign(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Resign calls spend and returns its zero results.
func (o electionOver) Resign(context.Context) (r0 error) {
	o.spend()
	return r0
}

// IsLeader calls spend and returns its zero results.
func (o electionOver) IsLeader(context.Context) (r0 bool) {
	o.spend()
	return r0
}

// newElectionOver returns the planted defect of a ceiling of Election,
// whose callables call spend.
func newElectionOver(_ *ElectionHandle, spend func()) Election {
	return electionOver{spend: spend}
}

// electionCampaignPanicsDefect is the planted defect of a Election whose
// Campaign panics.
type electionCampaignPanicsDefect struct {
	electionOverEmbedded
}

// Campaign is the override of the planted defect.
func (*electionCampaignPanicsDefect) Campaign(p0 context.Context) (r0 error) {
	panic("the planted defect of a Election whose Campaign panics")
}

// electionCampaignPanics returns the planted defect of a Election whose
// Campaign panics over a fresh subject.
func electionCampaignPanics(h *ElectionHandle) Election {
	return &electionCampaignPanicsDefect{electionOverEmbedded: h.Fresh()}
}

// electionCampaignIgnoresCancelDefect is the planted defect of a Election
// whose Campaign ignores the context.
type electionCampaignIgnoresCancelDefect struct {
	electionOverEmbedded
}

// Campaign is the override of the planted defect.
func (d *electionCampaignIgnoresCancelDefect) Campaign(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.electionOverEmbedded.Campaign(p0)
}

// electionCampaignIgnoresCancel returns the planted defect of a Election
// whose Campaign ignores the context over a fresh subject.
func electionCampaignIgnoresCancel(h *ElectionHandle) Election {
	return &electionCampaignIgnoresCancelDefect{electionOverEmbedded: h.Fresh()}
}

// electionIsLeaderLeadsOnDefect is the planted defect of a Election whose
// IsLeader reports leadership after Resign.
type electionIsLeaderLeadsOnDefect struct {
	electionOverEmbedded
}

// IsLeader is the override of the planted defect.
func (*electionIsLeaderLeadsOnDefect) IsLeader(p0 context.Context) (r0 bool) {
	r0 = true
	return r0
}

// electionIsLeaderLeadsOn returns the planted defect of a Election whose
// IsLeader reports leadership after Resign over a fresh subject.
func electionIsLeaderLeadsOn(h *ElectionHandle) Election {
	return &electionIsLeaderLeadsOnDefect{electionOverEmbedded: h.Fresh()}
}

// electionResignPanicsDefect is the planted defect of a Election whose
// Resign panics.
type electionResignPanicsDefect struct {
	electionOverEmbedded
}

// Resign is the override of the planted defect.
func (*electionResignPanicsDefect) Resign(p0 context.Context) (r0 error) {
	panic("the planted defect of a Election whose Resign panics")
}

// electionResignPanics returns the planted defect of a Election whose
// Resign panics over a fresh subject.
func electionResignPanics(h *ElectionHandle) Election {
	return &electionResignPanicsDefect{electionOverEmbedded: h.Fresh()}
}

// electionResignIgnoresCancelDefect is the planted defect of a Election
// whose Resign ignores the context.
type electionResignIgnoresCancelDefect struct {
	electionOverEmbedded
}

// Resign is the override of the planted defect.
func (d *electionResignIgnoresCancelDefect) Resign(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.electionOverEmbedded.Resign(p0)
}

// electionResignIgnoresCancel returns the planted defect of a Election
// whose Resign ignores the context over a fresh subject.
func electionResignIgnoresCancel(h *ElectionHandle) Election {
	return &electionResignIgnoresCancelDefect{electionOverEmbedded: h.Fresh()}
}

// electionIsLeaderPanicsDefect is the planted defect of a Election whose
// IsLeader panics.
type electionIsLeaderPanicsDefect struct {
	electionOverEmbedded
}

// IsLeader is the override of the planted defect.
func (*electionIsLeaderPanicsDefect) IsLeader(p0 context.Context) (r0 bool) {
	panic("the planted defect of a Election whose IsLeader panics")
}

// electionIsLeaderPanics returns the planted defect of a Election whose
// IsLeader panics over a fresh subject.
func electionIsLeaderPanics(h *ElectionHandle) Election {
	return &electionIsLeaderPanicsDefect{electionOverEmbedded: h.Fresh()}
}

// The derived checks of Pool:
//
// Pool.Get/smoke: Get returns without a panic for derived inputs.
//
// Pool.Get/cancel: Get returns an error for a cancelled context.
//
// Pool.Get/zero-on-error: Get returns the zero value beside an error.
//
// Pool.Get/pool: Put takes back what Get lent.
//
// Pool.Put/smoke: Put returns without a panic for derived inputs.
//
// Pool.Put/cancel: Put returns an error for a cancelled context.

// The derived checks of Pool that witness their laws at the shallow end:
//
// Pool.Get/smoke
//
// Pool.Put/smoke

// The unclaimed classifications of Pool:
//
// The shape:writer of Pool.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Pool writes.
type (
	PoolHandle = conform.Handle[Pool]
	PoolCheck  = conform.Check[Pool]
	PoolOption = conform.Option[Pool]
)

// PoolChecks contains the name of each callable of Pool, with the
// identities of its checks.
var PoolChecks = struct {
	Get poolGet
	Put poolPut
}{
	Get: poolGet{"Pool.Get"},
	Put: poolPut{"Pool.Put"},
}

// poolGet is the name of Pool.Get, with the identities of its checks.
type poolGet struct {
	conform.Callable[Pool]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c poolGet) Smoke() conform.ID[Pool] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c poolGet) Cancel() conform.ID[Pool] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c poolGet) ZeroOnError() conform.ID[Pool] {
	return c.Check("zero-on-error")
}

// Pool returns the identity of the check that Put takes back what Get
// lent.
func (c poolGet) Pool() conform.ID[Pool] {
	return c.Check("pool")
}

// poolPut is the name of Pool.Put, with the identities of its checks.
type poolPut struct {
	conform.Callable[Pool]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c poolPut) Smoke() conform.ID[Pool] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c poolPut) Cancel() conform.ID[Pool] {
	return c.Check("cancel")
}

// PoolSuite returns the checks of Pool, and the calls and the budgets of
// its callables.
func PoolSuite() conform.Suite[Pool] {
	return conform.Suite[Pool]{
		Declaration: "Pool",
		Checks: []conform.Check[Pool]{
			{
				ID:          PoolChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PoolHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx)
						return err
					})
				},
				Defect: poolGetPanics,
			},
			{
				ID:          PoolChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PoolHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx)
						return err
					})
				},
				Defect: poolGetIgnoresCancel,
			},
			{
				ID:          PoolChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PoolHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx)
					}, Value{})
				},
				Defect: poolGetValueBeside,
			},
			{
				ID:          PoolChecks.Get.Pool(),
				Statement:   "Put takes back what Get lent",
				Claim:       spec.Pool,
				Probe:       spec.TakesBackWhatItLent,
				Reason:      spec.UnreturnedLoan,
				Law:         spec.BorrowedValueIsTakenBack,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PoolHandle) {
					probe.TakesBackWhatItLent(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx)
					}, func(ctx context.Context, v Value) error {
						return h.Subject().Put(ctx, v)
					})
				},
				Defect: poolPutPutRefuses,
			},
			{
				ID:          PoolChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PoolHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: poolPutPanics,
			},
			{
				ID:          PoolChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PoolHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: poolPutIgnoresCancel,
			},
		},
		Calls: []conform.Call[Pool]{
			{
				Calls:    poolGetCalls,
				Over:     newPoolOver,
				Callable: PoolChecks.Get.Callable,
			},
			{
				Calls:    poolPutCalls,
				Over:     newPoolOver,
				Callable: PoolChecks.Put.Callable,
			},
		},
	}
}

// RunPoolSuite applies the checks of Pool, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunPoolSuite[T Pool](t *testing.T, newSubject func() T, opts ...PoolOption) {
	t.Helper()
	RunPoolSuiteFromHandle(t, func(*PoolHandle) T { return newSubject() }, opts...)
}

// RunPoolSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPoolSuiteFromHandle[T Pool](
	t *testing.T,
	newSubject func(h *PoolHandle) T,
	opts ...PoolOption,
) {
	t.Helper()
	conform.Run(t, PoolSuite(), poolSubjects(newSubject), opts...)
}

// FuzzPoolSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPoolSuite[T Pool](
	f *testing.F,
	callable conform.Callable[Pool],
	newSubject func(h *PoolHandle) T,
	opts ...PoolOption,
) {
	f.Helper()
	conform.Fuzz(f, PoolSuite(), callable, poolSubjects(newSubject), opts...)
}

// RunPoolBudgets applies the allocation ceilings of Pool and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunPoolBudgets[T Pool](t *testing.T, newSubject func(h *PoolHandle) T, opts ...PoolOption) {
	t.Helper()
	conform.RunBudgets(t, PoolSuite(), poolSubjects(newSubject), opts...)
}

// BenchPoolSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPoolSuite[T Pool](
	b conform.B,
	callable conform.Callable[Pool],
	newSubject func(h *PoolHandle) T,
	opts ...PoolOption,
) {
	b.Helper()
	conform.Bench(b, PoolSuite(), callable, poolSubjects(newSubject), opts...)
}

// poolSubjects returns the factory of the subjects of Pool that calls
// newSubject with the handle of each check.
func poolSubjects[T Pool](newSubject func(h *PoolHandle) T) func(h *PoolHandle) Pool {
	return func(h *PoolHandle) Pool { return newSubject(h) }
}

// poolGetCalls generates the calls of Pool.Get with the inputs that each
// case draws.
var poolGetCalls = prop.Composite(func(*prop.Case) func(context.Context, Pool) error {
	return func(ctx context.Context, subject Pool) error {
		_, err := subject.Get(ctx)
		return err
	}
})

// poolPutCalls generates the calls of Pool.Put with the inputs that each
// case draws.
var poolPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pool) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Pool) error {
		return subject.Put(ctx, v)
	}
})

// poolOverEmbedded is Pool under the name of the field that poolOver
// embeds.
type poolOverEmbedded = Pool

// poolOver is a subject of Pool whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type poolOver struct {
	poolOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o poolOver) Get(context.Context) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Put calls spend and returns its zero results.
func (o poolOver) Put(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// newPoolOver returns the planted defect of a ceiling of Pool, whose
// callables call spend.
func newPoolOver(_ *PoolHandle, spend func()) Pool {
	return poolOver{spend: spend}
}

// poolGetPanicsDefect is the planted defect of a Pool whose Get panics.
type poolGetPanicsDefect struct {
	poolOverEmbedded
}

// Get is the override of the planted defect.
func (*poolGetPanicsDefect) Get(p0 context.Context) (r0 Value, r1 error) {
	panic("the planted defect of a Pool whose Get panics")
}

// poolGetPanics returns the planted defect of a Pool whose Get panics over
// a fresh subject.
func poolGetPanics(h *PoolHandle) Pool {
	return &poolGetPanicsDefect{poolOverEmbedded: h.Fresh()}
}

// poolGetIgnoresCancelDefect is the planted defect of a Pool whose Get
// ignores the context.
type poolGetIgnoresCancelDefect struct {
	poolOverEmbedded
}

// Get is the override of the planted defect.
func (d *poolGetIgnoresCancelDefect) Get(p0 context.Context) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.poolOverEmbedded.Get(p0)
}

// poolGetIgnoresCancel returns the planted defect of a Pool whose Get
// ignores the context over a fresh subject.
func poolGetIgnoresCancel(h *PoolHandle) Pool {
	return &poolGetIgnoresCancelDefect{poolOverEmbedded: h.Fresh()}
}

// poolGetValueBesideDefect is the planted defect of a Pool whose Get
// returns a value beside an error.
type poolGetValueBesideDefect struct {
	poolOverEmbedded
}

// Get is the override of the planted defect.
func (d *poolGetValueBesideDefect) Get(p0 context.Context) (r0 Value, r1 error) {
	r0, r1 = d.poolOverEmbedded.Get(p0)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// poolGetValueBeside returns the planted defect of a Pool whose Get
// returns a value beside an error over a fresh subject.
func poolGetValueBeside(h *PoolHandle) Pool {
	return &poolGetValueBesideDefect{poolOverEmbedded: h.Fresh()}
}

// poolPutPutRefusesDefect is the planted defect of a Pool whose Put
// refuses what Get lent.
type poolPutPutRefusesDefect struct {
	poolOverEmbedded
}

// Put is the override of the planted defect.
func (*poolPutPutRefusesDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	r0 = conform.ErrGenerated
	return r0
}

// poolPutPutRefuses returns the planted defect of a Pool whose Put refuses
// what Get lent over a fresh subject.
func poolPutPutRefuses(h *PoolHandle) Pool {
	return &poolPutPutRefusesDefect{poolOverEmbedded: h.Fresh()}
}

// poolPutPanicsDefect is the planted defect of a Pool whose Put panics.
type poolPutPanicsDefect struct {
	poolOverEmbedded
}

// Put is the override of the planted defect.
func (*poolPutPanicsDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Pool whose Put panics")
}

// poolPutPanics returns the planted defect of a Pool whose Put panics over
// a fresh subject.
func poolPutPanics(h *PoolHandle) Pool {
	return &poolPutPanicsDefect{poolOverEmbedded: h.Fresh()}
}

// poolPutIgnoresCancelDefect is the planted defect of a Pool whose Put
// ignores the context.
type poolPutIgnoresCancelDefect struct {
	poolOverEmbedded
}

// Put is the override of the planted defect.
func (d *poolPutIgnoresCancelDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.poolOverEmbedded.Put(p0, p1)
}

// poolPutIgnoresCancel returns the planted defect of a Pool whose Put
// ignores the context over a fresh subject.
func poolPutIgnoresCancel(h *PoolHandle) Pool {
	return &poolPutIgnoresCancelDefect{poolOverEmbedded: h.Fresh()}
}

// The derived checks of Batch:
//
// Batch.Write/smoke: Write returns without a panic for derived inputs.
//
// Batch.Write/cancel: Write returns an error for a cancelled context.
//
// Batch.Write/batch-writer: Read finds no part of a value that Write
// refuses.
//
// Batch.Read/smoke: Read returns without a panic for derived inputs.
//
// Batch.Read/cancel: Read returns an error for a cancelled context.
//
// Batch.Read/zero-on-error: Read returns the zero value beside an error.

// The derived checks of Batch that witness their laws at the shallow end:
//
// Batch.Write/smoke
//
// Batch.Read/smoke

// The unclaimed classifications of Batch:
//
// The shape:writer of Batch.Write has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Batch writes.
type (
	BatchHandle = conform.Handle[Batch]
	BatchCheck  = conform.Check[Batch]
	BatchOption = conform.Option[Batch]
)

// BatchChecks contains the name of each callable of Batch, with the
// identities of its checks.
var BatchChecks = struct {
	Write batchWrite
	Read  batchRead
}{
	Write: batchWrite{"Batch.Write"},
	Read:  batchRead{"Batch.Read"},
}

// batchWrite is the name of Batch.Write, with the identities of its
// checks.
type batchWrite struct {
	conform.Callable[Batch]
}

// Smoke returns the identity of the check that Write returns without a
// panic for derived inputs.
func (c batchWrite) Smoke() conform.ID[Batch] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Write returns an error for
// a cancelled context.
func (c batchWrite) Cancel() conform.ID[Batch] {
	return c.Check("cancel")
}

// BatchWriter returns the identity of the check that Read finds no part of
// a value that Write refuses.
func (c batchWrite) BatchWriter() conform.ID[Batch] {
	return c.Check("batch-writer")
}

// batchRead is the name of Batch.Read, with the identities of its checks.
type batchRead struct {
	conform.Callable[Batch]
}

// Smoke returns the identity of the check that Read returns without a
// panic for derived inputs.
func (c batchRead) Smoke() conform.ID[Batch] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Read returns an error for
// a cancelled context.
func (c batchRead) Cancel() conform.ID[Batch] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Read returns the zero
// value beside an error.
func (c batchRead) ZeroOnError() conform.ID[Batch] {
	return c.Check("zero-on-error")
}

// BatchSuite returns the checks of Batch, and the calls and the budgets of
// its callables.
func BatchSuite() conform.Suite[Batch] {
	return conform.Suite[Batch]{
		Declaration: "Batch",
		Checks: []conform.Check[Batch]{
			{
				ID:          BatchChecks.Write.Smoke(),
				Statement:   "Write returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BatchHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Write(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: batchWritePanics,
			},
			{
				ID:          BatchChecks.Write.Cancel(),
				Statement:   "Write returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BatchHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Write(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: batchWriteIgnoresCancel,
			},
			{
				ID:          BatchChecks.Write.BatchWriter(),
				Statement:   "Read finds no part of a value that Write refuses",
				Claim:       spec.BatchWriter,
				Probe:       spec.RefusedLeavesNothing,
				Reason:      spec.PartialWrite,
				Law:         spec.RefusedWriteLeavesNothing,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BatchHandle) {
					probe.RefusedLeavesNothing(h, func(ctx context.Context, v Value) error {
						return h.Subject().Write(ctx, v)
					}, func(ctx context.Context) (Value, error) {
						return h.Subject().Read(ctx, (Value{Key: "Key-a"}).Key)
					}, Unkeyed, Value{})
				},
				Defect: batchReadFilesRefusedWrite,
			},
			{
				ID:          BatchChecks.Read.Smoke(),
				Statement:   "Read returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BatchHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, string("key-a"))
						return err
					})
				},
				Defect: batchReadPanics,
			},
			{
				ID:          BatchChecks.Read.Cancel(),
				Statement:   "Read returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BatchHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, string("key-a"))
						return err
					})
				},
				Defect: batchReadIgnoresCancel,
			},
			{
				ID:          BatchChecks.Read.ZeroOnError(),
				Statement:   "Read returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BatchHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Read(ctx, string("key-a"))
					}, Value{})
				},
				Defect: batchReadValueBeside,
			},
		},
		Calls: []conform.Call[Batch]{
			{
				Calls:    batchWriteCalls,
				Over:     newBatchOver,
				Callable: BatchChecks.Write.Callable,
			},
			{
				Calls:    batchReadCalls,
				Over:     newBatchOver,
				Callable: BatchChecks.Read.Callable,
			},
		},
	}
}

// RunBatchSuite applies the checks of Batch, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunBatchSuite[T Batch](t *testing.T, newSubject func() T, opts ...BatchOption) {
	t.Helper()
	RunBatchSuiteFromHandle(t, func(*BatchHandle) T { return newSubject() }, opts...)
}

// RunBatchSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunBatchSuiteFromHandle[T Batch](
	t *testing.T,
	newSubject func(h *BatchHandle) T,
	opts ...BatchOption,
) {
	t.Helper()
	conform.Run(t, BatchSuite(), batchSubjects(newSubject), opts...)
}

// FuzzBatchSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzBatchSuite[T Batch](
	f *testing.F,
	callable conform.Callable[Batch],
	newSubject func(h *BatchHandle) T,
	opts ...BatchOption,
) {
	f.Helper()
	conform.Fuzz(f, BatchSuite(), callable, batchSubjects(newSubject), opts...)
}

// RunBatchBudgets applies the allocation ceilings of Batch and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunBatchBudgets[T Batch](
	t *testing.T,
	newSubject func(h *BatchHandle) T,
	opts ...BatchOption,
) {
	t.Helper()
	conform.RunBudgets(t, BatchSuite(), batchSubjects(newSubject), opts...)
}

// BenchBatchSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchBatchSuite[T Batch](
	b conform.B,
	callable conform.Callable[Batch],
	newSubject func(h *BatchHandle) T,
	opts ...BatchOption,
) {
	b.Helper()
	conform.Bench(b, BatchSuite(), callable, batchSubjects(newSubject), opts...)
}

// batchSubjects returns the factory of the subjects of Batch that calls
// newSubject with the handle of each check.
func batchSubjects[T Batch](newSubject func(h *BatchHandle) T) func(h *BatchHandle) Batch {
	return func(h *BatchHandle) Batch { return newSubject(h) }
}

// batchWriteCalls generates the calls of Batch.Write with the inputs that
// each case draws.
var batchWriteCalls = prop.Composite(func(c *prop.Case) func(context.Context, Batch) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Batch) error {
		return subject.Write(ctx, v)
	}
})

// batchReadCalls generates the calls of Batch.Read with the inputs that
// each case draws.
var batchReadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Batch) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Batch) error {
		_, err := subject.Read(ctx, key)
		return err
	}
})

// batchOverEmbedded is Batch under the name of the field that batchOver
// embeds.
type batchOverEmbedded = Batch

// batchOver is a subject of Batch whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type batchOver struct {
	batchOverEmbedded
	spend func()
}

// Write calls spend and returns its zero results.
func (o batchOver) Write(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Read calls spend and returns its zero results.
func (o batchOver) Read(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newBatchOver returns the planted defect of a ceiling of Batch, whose
// callables call spend.
func newBatchOver(_ *BatchHandle, spend func()) Batch {
	return batchOver{spend: spend}
}

// batchWritePanicsDefect is the planted defect of a Batch whose Write
// panics.
type batchWritePanicsDefect struct {
	batchOverEmbedded
}

// Write is the override of the planted defect.
func (*batchWritePanicsDefect) Write(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Batch whose Write panics")
}

// batchWritePanics returns the planted defect of a Batch whose Write
// panics over a fresh subject.
func batchWritePanics(h *BatchHandle) Batch {
	return &batchWritePanicsDefect{batchOverEmbedded: h.Fresh()}
}

// batchWriteIgnoresCancelDefect is the planted defect of a Batch whose
// Write ignores the context.
type batchWriteIgnoresCancelDefect struct {
	batchOverEmbedded
}

// Write is the override of the planted defect.
func (d *batchWriteIgnoresCancelDefect) Write(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.batchOverEmbedded.Write(p0, p1)
}

// batchWriteIgnoresCancel returns the planted defect of a Batch whose
// Write ignores the context over a fresh subject.
func batchWriteIgnoresCancel(h *BatchHandle) Batch {
	return &batchWriteIgnoresCancelDefect{batchOverEmbedded: h.Fresh()}
}

// batchReadFilesRefusedWriteDefect is the planted defect of a Batch whose
// Write files a part of a write that it refused.
type batchReadFilesRefusedWriteDefect struct {
	batchOverEmbedded
}

// Read is the override of the planted defect.
func (*batchReadFilesRefusedWriteDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = Value{Key: "Key-b"}
	return r0, r1
}

// batchReadFilesRefusedWrite returns the planted defect of a Batch whose
// Write files a part of a write that it refused over a fresh subject.
func batchReadFilesRefusedWrite(h *BatchHandle) Batch {
	return &batchReadFilesRefusedWriteDefect{batchOverEmbedded: h.Fresh()}
}

// batchReadPanicsDefect is the planted defect of a Batch whose Read
// panics.
type batchReadPanicsDefect struct {
	batchOverEmbedded
}

// Read is the override of the planted defect.
func (*batchReadPanicsDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Batch whose Read panics")
}

// batchReadPanics returns the planted defect of a Batch whose Read panics
// over a fresh subject.
func batchReadPanics(h *BatchHandle) Batch {
	return &batchReadPanicsDefect{batchOverEmbedded: h.Fresh()}
}

// batchReadIgnoresCancelDefect is the planted defect of a Batch whose Read
// ignores the context.
type batchReadIgnoresCancelDefect struct {
	batchOverEmbedded
}

// Read is the override of the planted defect.
func (d *batchReadIgnoresCancelDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.batchOverEmbedded.Read(p0, p1)
}

// batchReadIgnoresCancel returns the planted defect of a Batch whose Read
// ignores the context over a fresh subject.
func batchReadIgnoresCancel(h *BatchHandle) Batch {
	return &batchReadIgnoresCancelDefect{batchOverEmbedded: h.Fresh()}
}

// batchReadValueBesideDefect is the planted defect of a Batch whose Read
// returns a value beside an error.
type batchReadValueBesideDefect struct {
	batchOverEmbedded
}

// Read is the override of the planted defect.
func (d *batchReadValueBesideDefect) Read(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.batchOverEmbedded.Read(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// batchReadValueBeside returns the planted defect of a Batch whose Read
// returns a value beside an error over a fresh subject.
func batchReadValueBeside(h *BatchHandle) Batch {
	return &batchReadValueBesideDefect{batchOverEmbedded: h.Fresh()}
}

// The derived checks of Screened:
//
// Screened.Store/smoke: Store returns without a panic for derived inputs.
//
// Screened.Store/cancel: Store returns an error for a cancelled context.
//
// Screened.Store/validates: Store refuses the value that Validate refuses.
//
// Screened.Validate/smoke: Validate returns without a panic for derived
// inputs.

// The derived checks of Screened that witness their laws at the shallow
// end:
//
// Screened.Store/smoke
//
// Screened.Store/validates
//
// Screened.Validate/smoke

// The unclaimed classifications of Screened:
//
// The shape:writer of Screened.Store has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Screened.Validate has the stance silent. A write is
// an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Screened writes.
type (
	ScreenedHandle = conform.Handle[Screened]
	ScreenedCheck  = conform.Check[Screened]
	ScreenedOption = conform.Option[Screened]
)

// ScreenedChecks contains the name of each callable of Screened, with the
// identities of its checks.
var ScreenedChecks = struct {
	Store    screenedStore
	Validate screenedValidate
}{
	Store:    screenedStore{"Screened.Store"},
	Validate: screenedValidate{"Screened.Validate"},
}

// screenedStore is the name of Screened.Store, with the identities of its
// checks.
type screenedStore struct {
	conform.Callable[Screened]
}

// Smoke returns the identity of the check that Store returns without a
// panic for derived inputs.
func (c screenedStore) Smoke() conform.ID[Screened] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Store returns an error for
// a cancelled context.
func (c screenedStore) Cancel() conform.ID[Screened] {
	return c.Check("cancel")
}

// Validates returns the identity of the check that Store refuses the value
// that Validate refuses.
func (c screenedStore) Validates() conform.ID[Screened] {
	return c.Check("validates")
}

// screenedValidate is the name of Screened.Validate, with the identities
// of its checks.
type screenedValidate struct {
	conform.Callable[Screened]
}

// Smoke returns the identity of the check that Validate returns without a
// panic for derived inputs.
func (c screenedValidate) Smoke() conform.ID[Screened] {
	return c.Check("smoke")
}

// ScreenedSuite returns the checks of Screened, and the calls and the
// budgets of its callables.
func ScreenedSuite() conform.Suite[Screened] {
	return conform.Suite[Screened]{
		Declaration: "Screened",
		Checks: []conform.Check[Screened]{
			{
				ID:          ScreenedChecks.Store.Smoke(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScreenedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: screenedStorePanics,
			},
			{
				ID:          ScreenedChecks.Store.Cancel(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScreenedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: screenedStoreIgnoresCancel,
			},
			{
				ID:          ScreenedChecks.Store.Validates(),
				Statement:   "Store refuses the value that Validate refuses",
				Claim:       spec.Validates,
				Probe:       spec.RefusesTheInvalid,
				Reason:      spec.UnscreenedWrite,
				Law:         spec.DeclaredInvalidValueIsRefused,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScreenedHandle) {
					probe.RefusesTheInvalid(h, func(ctx context.Context, v Value) error {
						return h.Subject().Store(ctx, v)
					}, func(ctx context.Context, v Value) error {
						return h.Subject().Validate(v)
					}, Invalid)
				},
				Defect: screenedStoreAcceptsTheInvalid,
			},
			{
				ID:          ScreenedChecks.Validate.Smoke(),
				Statement:   "Validate returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScreenedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Validate(Value{Key: "Key-a"})
					})
				},
				Defect: screenedValidatePanics,
			},
		},
		Calls: []conform.Call[Screened]{
			{
				Calls:    screenedStoreCalls,
				Over:     newScreenedOver,
				Callable: ScreenedChecks.Store.Callable,
			},
			{
				Calls:    screenedValidateCalls,
				Over:     newScreenedOver,
				Callable: ScreenedChecks.Validate.Callable,
			},
		},
	}
}

// RunScreenedSuite applies the checks of Screened, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunScreenedSuite[T Screened](t *testing.T, newSubject func() T, opts ...ScreenedOption) {
	t.Helper()
	RunScreenedSuiteFromHandle(t, func(*ScreenedHandle) T { return newSubject() }, opts...)
}

// RunScreenedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunScreenedSuiteFromHandle[T Screened](
	t *testing.T,
	newSubject func(h *ScreenedHandle) T,
	opts ...ScreenedOption,
) {
	t.Helper()
	conform.Run(t, ScreenedSuite(), screenedSubjects(newSubject), opts...)
}

// FuzzScreenedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzScreenedSuite[T Screened](
	f *testing.F,
	callable conform.Callable[Screened],
	newSubject func(h *ScreenedHandle) T,
	opts ...ScreenedOption,
) {
	f.Helper()
	conform.Fuzz(f, ScreenedSuite(), callable, screenedSubjects(newSubject), opts...)
}

// RunScreenedBudgets applies the allocation ceilings of Screened and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunScreenedBudgets[T Screened](
	t *testing.T,
	newSubject func(h *ScreenedHandle) T,
	opts ...ScreenedOption,
) {
	t.Helper()
	conform.RunBudgets(t, ScreenedSuite(), screenedSubjects(newSubject), opts...)
}

// BenchScreenedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchScreenedSuite[T Screened](
	b conform.B,
	callable conform.Callable[Screened],
	newSubject func(h *ScreenedHandle) T,
	opts ...ScreenedOption,
) {
	b.Helper()
	conform.Bench(b, ScreenedSuite(), callable, screenedSubjects(newSubject), opts...)
}

// screenedSubjects returns the factory of the subjects of Screened that
// calls newSubject with the handle of each check.
func screenedSubjects[T Screened](
	newSubject func(h *ScreenedHandle) T,
) func(h *ScreenedHandle) Screened {
	return func(h *ScreenedHandle) Screened { return newSubject(h) }
}

// screenedStoreCalls generates the calls of Screened.Store with the inputs
// that each case draws.
var screenedStoreCalls = prop.Composite(func(c *prop.Case) func(context.Context, Screened) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Screened) error {
		return subject.Store(ctx, v)
	}
})

// screenedValidateCalls generates the calls of Screened.Validate with the
// inputs that each case draws.
var screenedValidateCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, Screened) error {
		v := c.Draw(prop.Of[Value](), "v")
		return func(_ context.Context, subject Screened) error {
			return subject.Validate(v)
		}
	},
)

// screenedOverEmbedded is Screened under the name of the field that
// screenedOver embeds.
type screenedOverEmbedded = Screened

// screenedOver is a subject of Screened whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type screenedOver struct {
	screenedOverEmbedded
	spend func()
}

// Store calls spend and returns its zero results.
func (o screenedOver) Store(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Validate calls spend and returns its zero results.
func (o screenedOver) Validate(Value) (r0 error) {
	o.spend()
	return r0
}

// newScreenedOver returns the planted defect of a ceiling of Screened,
// whose callables call spend.
func newScreenedOver(_ *ScreenedHandle, spend func()) Screened {
	return screenedOver{spend: spend}
}

// screenedStorePanicsDefect is the planted defect of a Screened whose
// Store panics.
type screenedStorePanicsDefect struct {
	screenedOverEmbedded
}

// Store is the override of the planted defect.
func (*screenedStorePanicsDefect) Store(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Screened whose Store panics")
}

// screenedStorePanics returns the planted defect of a Screened whose Store
// panics over a fresh subject.
func screenedStorePanics(h *ScreenedHandle) Screened {
	return &screenedStorePanicsDefect{screenedOverEmbedded: h.Fresh()}
}

// screenedStoreIgnoresCancelDefect is the planted defect of a Screened
// whose Store ignores the context.
type screenedStoreIgnoresCancelDefect struct {
	screenedOverEmbedded
}

// Store is the override of the planted defect.
func (d *screenedStoreIgnoresCancelDefect) Store(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.screenedOverEmbedded.Store(p0, p1)
}

// screenedStoreIgnoresCancel returns the planted defect of a Screened
// whose Store ignores the context over a fresh subject.
func screenedStoreIgnoresCancel(h *ScreenedHandle) Screened {
	return &screenedStoreIgnoresCancelDefect{screenedOverEmbedded: h.Fresh()}
}

// screenedStoreAcceptsTheInvalidDefect is the planted defect of a Screened
// whose Store accepts a value that Validate refuses.
type screenedStoreAcceptsTheInvalidDefect struct {
	screenedOverEmbedded
}

// Store is the override of the planted defect.
func (*screenedStoreAcceptsTheInvalidDefect) Store(p0 context.Context, p1 Value) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// screenedStoreAcceptsTheInvalid returns the planted defect of a Screened
// whose Store accepts a value that Validate refuses over a fresh subject.
func screenedStoreAcceptsTheInvalid(h *ScreenedHandle) Screened {
	return &screenedStoreAcceptsTheInvalidDefect{screenedOverEmbedded: h.Fresh()}
}

// screenedValidatePanicsDefect is the planted defect of a Screened whose
// Validate panics.
type screenedValidatePanicsDefect struct {
	screenedOverEmbedded
}

// Validate is the override of the planted defect.
func (*screenedValidatePanicsDefect) Validate(p0 Value) (r0 error) {
	panic("the planted defect of a Screened whose Validate panics")
}

// screenedValidatePanics returns the planted defect of a Screened whose
// Validate panics over a fresh subject.
func screenedValidatePanics(h *ScreenedHandle) Screened {
	return &screenedValidatePanicsDefect{screenedOverEmbedded: h.Fresh()}
}

// The derived checks of Limiter:
//
// Limiter.Run/smoke: Run returns without a panic for derived inputs.
//
// Limiter.Run/cancel: Run returns an error for a cancelled context.
//
// Limiter.Run/rate-limit: Run admits its declared burst and refuses the
// call after it.

// The derived checks of Limiter that witness their laws at the shallow
// end:
//
// Limiter.Run/smoke

// The unclaimed classifications of Limiter:
//
// The shape:writer of Limiter.Run has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Limiter writes.
type (
	LimiterHandle = conform.Handle[Limiter]
	LimiterCheck  = conform.Check[Limiter]
	LimiterOption = conform.Option[Limiter]
)

// LimiterChecks contains the name of each callable of Limiter, with the
// identities of its checks.
var LimiterChecks = struct {
	Run limiterRun
}{
	Run: limiterRun{"Limiter.Run"},
}

// limiterRun is the name of Limiter.Run, with the identities of its
// checks.
type limiterRun struct {
	conform.Callable[Limiter]
}

// Smoke returns the identity of the check that Run returns without a panic
// for derived inputs.
func (c limiterRun) Smoke() conform.ID[Limiter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Run returns an error for a
// cancelled context.
func (c limiterRun) Cancel() conform.ID[Limiter] {
	return c.Check("cancel")
}

// RateLimit returns the identity of the check that Run admits its declared
// burst and refuses the call after it.
func (c limiterRun) RateLimit() conform.ID[Limiter] {
	return c.Check("rate-limit")
}

// LimiterSuite returns the checks of Limiter, and the calls and the
// budgets of its callables.
func LimiterSuite() conform.Suite[Limiter] {
	return conform.Suite[Limiter]{
		Declaration: "Limiter",
		Checks: []conform.Check[Limiter]{
			{
				ID:          LimiterChecks.Run.Smoke(),
				Statement:   "Run returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LimiterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, string("key-a"))
					})
				},
				Defect: limiterRunPanics,
			},
			{
				ID:          LimiterChecks.Run.Cancel(),
				Statement:   "Run returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LimiterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, string("key-a"))
					})
				},
				Defect: limiterRunIgnoresCancel,
			},
			{
				ID:          LimiterChecks.Run.RateLimit(),
				Statement:   "Run admits its declared burst and refuses the call after it",
				Claim:       spec.RateLimit,
				Probe:       spec.RefusesPastTheBurst,
				Reason:      spec.PastTheBurst,
				Law:         spec.BurstIsSpentThenRefused,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LimiterHandle) {
					probe.RefusesPastTheBurst(h, func(ctx context.Context) error {
						return h.Subject().Run(ctx, string("key-a"))
					}, 10, ErrLimited)
				},
				Defect: limiterRunSucceedsAlways,
			},
		},
		Calls: []conform.Call[Limiter]{
			{
				Calls:    limiterRunCalls,
				Over:     newLimiterOver,
				Callable: LimiterChecks.Run.Callable,
			},
		},
	}
}

// RunLimiterSuite applies the checks of Limiter, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLimiterSuite[T Limiter](t *testing.T, newSubject func() T, opts ...LimiterOption) {
	t.Helper()
	RunLimiterSuiteFromHandle(t, func(*LimiterHandle) T { return newSubject() }, opts...)
}

// RunLimiterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLimiterSuiteFromHandle[T Limiter](
	t *testing.T,
	newSubject func(h *LimiterHandle) T,
	opts ...LimiterOption,
) {
	t.Helper()
	conform.Run(t, LimiterSuite(), limiterSubjects(newSubject), opts...)
}

// FuzzLimiterSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLimiterSuite[T Limiter](
	f *testing.F,
	callable conform.Callable[Limiter],
	newSubject func(h *LimiterHandle) T,
	opts ...LimiterOption,
) {
	f.Helper()
	conform.Fuzz(f, LimiterSuite(), callable, limiterSubjects(newSubject), opts...)
}

// RunLimiterBudgets applies the allocation ceilings of Limiter and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLimiterBudgets[T Limiter](
	t *testing.T,
	newSubject func(h *LimiterHandle) T,
	opts ...LimiterOption,
) {
	t.Helper()
	conform.RunBudgets(t, LimiterSuite(), limiterSubjects(newSubject), opts...)
}

// BenchLimiterSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchLimiterSuite[T Limiter](
	b conform.B,
	callable conform.Callable[Limiter],
	newSubject func(h *LimiterHandle) T,
	opts ...LimiterOption,
) {
	b.Helper()
	conform.Bench(b, LimiterSuite(), callable, limiterSubjects(newSubject), opts...)
}

// limiterSubjects returns the factory of the subjects of Limiter that
// calls newSubject with the handle of each check.
func limiterSubjects[T Limiter](
	newSubject func(h *LimiterHandle) T,
) func(h *LimiterHandle) Limiter {
	return func(h *LimiterHandle) Limiter { return newSubject(h) }
}

// limiterRunCalls generates the calls of Limiter.Run with the inputs that
// each case draws.
var limiterRunCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limiter) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Limiter) error {
		return subject.Run(ctx, key)
	}
})

// limiterOverEmbedded is Limiter under the name of the field that
// limiterOver embeds.
type limiterOverEmbedded = Limiter

// limiterOver is a subject of Limiter whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type limiterOver struct {
	limiterOverEmbedded
	spend func()
}

// Run calls spend and returns its zero results.
func (o limiterOver) Run(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// newLimiterOver returns the planted defect of a ceiling of Limiter, whose
// callables call spend.
func newLimiterOver(_ *LimiterHandle, spend func()) Limiter {
	return limiterOver{spend: spend}
}

// limiterRunPanicsDefect is the planted defect of a Limiter whose Run
// panics.
type limiterRunPanicsDefect struct {
	limiterOverEmbedded
}

// Run is the override of the planted defect.
func (*limiterRunPanicsDefect) Run(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Limiter whose Run panics")
}

// limiterRunPanics returns the planted defect of a Limiter whose Run
// panics over a fresh subject.
func limiterRunPanics(h *LimiterHandle) Limiter {
	return &limiterRunPanicsDefect{limiterOverEmbedded: h.Fresh()}
}

// limiterRunIgnoresCancelDefect is the planted defect of a Limiter whose
// Run ignores the context.
type limiterRunIgnoresCancelDefect struct {
	limiterOverEmbedded
}

// Run is the override of the planted defect.
func (d *limiterRunIgnoresCancelDefect) Run(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.limiterOverEmbedded.Run(p0, p1)
}

// limiterRunIgnoresCancel returns the planted defect of a Limiter whose
// Run ignores the context over a fresh subject.
func limiterRunIgnoresCancel(h *LimiterHandle) Limiter {
	return &limiterRunIgnoresCancelDefect{limiterOverEmbedded: h.Fresh()}
}

// limiterRunSucceedsAlwaysDefect is the planted defect of a Limiter whose
// Run admits calls past the declared burst.
type limiterRunSucceedsAlwaysDefect struct {
	limiterOverEmbedded
}

// Run is the override of the planted defect.
func (*limiterRunSucceedsAlwaysDefect) Run(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// limiterRunSucceedsAlways returns the planted defect of a Limiter whose
// Run admits calls past the declared burst over a fresh subject.
func limiterRunSucceedsAlways(h *LimiterHandle) Limiter {
	return &limiterRunSucceedsAlwaysDefect{limiterOverEmbedded: h.Fresh()}
}

// The derived checks of Renderer:
//
// Renderer.Render/smoke: Render returns without a panic for derived
// inputs.
//
// Renderer.Render/cancel: Render returns an error for a cancelled context.
//
// Renderer.Render/zero-on-error: Render returns the zero value beside an
// error.
//
// Renderer.Render/xss-safe: Render does not return the value that its
// directive declares unsafe.
//
// Renderer.Classify/smoke: Classify returns without a panic for derived
// inputs.
//
// Renderer.Classify/cancel: Classify returns an error for a cancelled
// context.
//
// Renderer.Classify/zero-on-error: Classify returns the zero value beside
// an error.
//
// Renderer.Classify/total: Classify succeeds for the input that its
// directive declares at the edge of the domain.

// The derived checks of Renderer that witness their laws at the shallow
// end:
//
// Renderer.Render/smoke
//
// Renderer.Render/xss-safe
//
// Renderer.Classify/smoke
//
// Renderer.Classify/total

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Renderer writes.
type (
	RendererHandle = conform.Handle[Renderer]
	RendererCheck  = conform.Check[Renderer]
	RendererOption = conform.Option[Renderer]
)

// RendererChecks contains the name of each callable of Renderer, with the
// identities of its checks.
var RendererChecks = struct {
	Render   rendererRender
	Classify rendererClassify
}{
	Render:   rendererRender{"Renderer.Render"},
	Classify: rendererClassify{"Renderer.Classify"},
}

// rendererRender is the name of Renderer.Render, with the identities of
// its checks.
type rendererRender struct {
	conform.Callable[Renderer]
}

// Smoke returns the identity of the check that Render returns without a
// panic for derived inputs.
func (c rendererRender) Smoke() conform.ID[Renderer] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Render returns an error
// for a cancelled context.
func (c rendererRender) Cancel() conform.ID[Renderer] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Render returns the
// zero value beside an error.
func (c rendererRender) ZeroOnError() conform.ID[Renderer] {
	return c.Check("zero-on-error")
}

// XSSSafe returns the identity of the check that Render does not return
// the value that its directive declares unsafe.
func (c rendererRender) XSSSafe() conform.ID[Renderer] {
	return c.Check("xss-safe")
}

// rendererClassify is the name of Renderer.Classify, with the identities
// of its checks.
type rendererClassify struct {
	conform.Callable[Renderer]
}

// Smoke returns the identity of the check that Classify returns without a
// panic for derived inputs.
func (c rendererClassify) Smoke() conform.ID[Renderer] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Classify returns an error
// for a cancelled context.
func (c rendererClassify) Cancel() conform.ID[Renderer] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Classify returns the
// zero value beside an error.
func (c rendererClassify) ZeroOnError() conform.ID[Renderer] {
	return c.Check("zero-on-error")
}

// Total returns the identity of the check that Classify succeeds for the
// input that its directive declares at the edge of the domain.
func (c rendererClassify) Total() conform.ID[Renderer] {
	return c.Check("total")
}

// RendererSuite returns the checks of Renderer, and the calls and the
// budgets of its callables.
func RendererSuite() conform.Suite[Renderer] {
	return conform.Suite[Renderer]{
		Declaration: "Renderer",
		Checks: []conform.Check[Renderer]{
			{
				ID:          RendererChecks.Render.Smoke(),
				Statement:   "Render returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RendererHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Render(ctx, string("in-a"))
						return err
					})
				},
				Defect: rendererRenderPanics,
			},
			{
				ID:          RendererChecks.Render.Cancel(),
				Statement:   "Render returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RendererHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Render(ctx, string("in-a"))
						return err
					})
				},
				Defect: rendererRenderIgnoresCancel,
			},
			{
				ID:          RendererChecks.Render.ZeroOnError(),
				Statement:   "Render returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RendererHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Render(ctx, string("in-a"))
					}, string(""))
				},
				Defect: rendererRenderValueBeside,
			},
			{
				ID:          RendererChecks.Render.XSSSafe(),
				Statement:   "Render does not return the value that its directive declares unsafe",
				Claim:       spec.XSSSafe,
				Probe:       spec.EscapesTheUnsafe,
				Reason:      spec.UnescapedOutput,
				Law:         spec.UnsafeValueDoesNotSurvive,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RendererHandle) {
					probe.EscapesTheUnsafe(h, func(ctx context.Context) (string, error) {
						return h.Subject().Render(ctx, Unsafe)
					}, Unsafe)
				},
				Defect: rendererRenderEchoes,
			},
			{
				ID:          RendererChecks.Classify.Smoke(),
				Statement:   "Classify returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RendererHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Classify(ctx, string("in-a"))
						return err
					})
				},
				Defect: rendererClassifyPanics,
			},
			{
				ID:          RendererChecks.Classify.Cancel(),
				Statement:   "Classify returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RendererHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Classify(ctx, string("in-a"))
						return err
					})
				},
				Defect: rendererClassifyIgnoresCancel,
			},
			{
				ID:          RendererChecks.Classify.ZeroOnError(),
				Statement:   "Classify returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RendererHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Classify(ctx, string("in-a"))
					}, string(""))
				},
				Defect: rendererClassifyValueBeside,
			},
			{
				ID:          RendererChecks.Classify.Total(),
				Statement:   "Classify succeeds for the input that its directive declares at the edge of the domain",
				Claim:       spec.Total,
				Probe:       spec.AcceptsTheEdge,
				Reason:      spec.RefusedTheEdge,
				Law:         spec.CallAcceptsItsDeclaredDomain,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RendererHandle) {
					probe.AcceptsTheEdge(h, func(ctx context.Context, v string) error {
						_, err := h.Subject().Classify(ctx, v)
						return err
					}, Edge)
				},
				Defect: rendererClassifyFailsAlways,
			},
		},
		Calls: []conform.Call[Renderer]{
			{
				Calls:    rendererRenderCalls,
				Over:     newRendererOver,
				Callable: RendererChecks.Render.Callable,
			},
			{
				Calls:    rendererClassifyCalls,
				Over:     newRendererOver,
				Callable: RendererChecks.Classify.Callable,
			},
		},
	}
}

// RunRendererSuite applies the checks of Renderer, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunRendererSuite[T Renderer](t *testing.T, newSubject func() T, opts ...RendererOption) {
	t.Helper()
	RunRendererSuiteFromHandle(t, func(*RendererHandle) T { return newSubject() }, opts...)
}

// RunRendererSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRendererSuiteFromHandle[T Renderer](
	t *testing.T,
	newSubject func(h *RendererHandle) T,
	opts ...RendererOption,
) {
	t.Helper()
	conform.Run(t, RendererSuite(), rendererSubjects(newSubject), opts...)
}

// FuzzRendererSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRendererSuite[T Renderer](
	f *testing.F,
	callable conform.Callable[Renderer],
	newSubject func(h *RendererHandle) T,
	opts ...RendererOption,
) {
	f.Helper()
	conform.Fuzz(f, RendererSuite(), callable, rendererSubjects(newSubject), opts...)
}

// RunRendererBudgets applies the allocation ceilings of Renderer and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRendererBudgets[T Renderer](
	t *testing.T,
	newSubject func(h *RendererHandle) T,
	opts ...RendererOption,
) {
	t.Helper()
	conform.RunBudgets(t, RendererSuite(), rendererSubjects(newSubject), opts...)
}

// BenchRendererSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchRendererSuite[T Renderer](
	b conform.B,
	callable conform.Callable[Renderer],
	newSubject func(h *RendererHandle) T,
	opts ...RendererOption,
) {
	b.Helper()
	conform.Bench(b, RendererSuite(), callable, rendererSubjects(newSubject), opts...)
}

// rendererSubjects returns the factory of the subjects of Renderer that
// calls newSubject with the handle of each check.
func rendererSubjects[T Renderer](
	newSubject func(h *RendererHandle) T,
) func(h *RendererHandle) Renderer {
	return func(h *RendererHandle) Renderer { return newSubject(h) }
}

// rendererRenderCalls generates the calls of Renderer.Render with the
// inputs that each case draws.
var rendererRenderCalls = prop.Composite(func(c *prop.Case) func(context.Context, Renderer) error {
	in := c.Draw(prop.Of[string](), "in")
	return func(ctx context.Context, subject Renderer) error {
		_, err := subject.Render(ctx, in)
		return err
	}
})

// rendererClassifyCalls generates the calls of Renderer.Classify with the
// inputs that each case draws.
var rendererClassifyCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, Renderer) error {
		in := c.Draw(prop.Of[string](), "in")
		return func(ctx context.Context, subject Renderer) error {
			_, err := subject.Classify(ctx, in)
			return err
		}
	},
)

// rendererOverEmbedded is Renderer under the name of the field that
// rendererOver embeds.
type rendererOverEmbedded = Renderer

// rendererOver is a subject of Renderer whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type rendererOver struct {
	rendererOverEmbedded
	spend func()
}

// Render calls spend and returns its zero results.
func (o rendererOver) Render(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// Classify calls spend and returns its zero results.
func (o rendererOver) Classify(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newRendererOver returns the planted defect of a ceiling of Renderer,
// whose callables call spend.
func newRendererOver(_ *RendererHandle, spend func()) Renderer {
	return rendererOver{spend: spend}
}

// rendererRenderPanicsDefect is the planted defect of a Renderer whose
// Render panics.
type rendererRenderPanicsDefect struct {
	rendererOverEmbedded
}

// Render is the override of the planted defect.
func (*rendererRenderPanicsDefect) Render(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Renderer whose Render panics")
}

// rendererRenderPanics returns the planted defect of a Renderer whose
// Render panics over a fresh subject.
func rendererRenderPanics(h *RendererHandle) Renderer {
	return &rendererRenderPanicsDefect{rendererOverEmbedded: h.Fresh()}
}

// rendererRenderIgnoresCancelDefect is the planted defect of a Renderer
// whose Render ignores the context.
type rendererRenderIgnoresCancelDefect struct {
	rendererOverEmbedded
}

// Render is the override of the planted defect.
func (d *rendererRenderIgnoresCancelDefect) Render(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.rendererOverEmbedded.Render(p0, p1)
}

// rendererRenderIgnoresCancel returns the planted defect of a Renderer
// whose Render ignores the context over a fresh subject.
func rendererRenderIgnoresCancel(h *RendererHandle) Renderer {
	return &rendererRenderIgnoresCancelDefect{rendererOverEmbedded: h.Fresh()}
}

// rendererRenderValueBesideDefect is the planted defect of a Renderer
// whose Render returns a value beside an error.
type rendererRenderValueBesideDefect struct {
	rendererOverEmbedded
}

// Render is the override of the planted defect.
func (d *rendererRenderValueBesideDefect) Render(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	r0, r1 = d.rendererOverEmbedded.Render(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// rendererRenderValueBeside returns the planted defect of a Renderer whose
// Render returns a value beside an error over a fresh subject.
func rendererRenderValueBeside(h *RendererHandle) Renderer {
	return &rendererRenderValueBesideDefect{rendererOverEmbedded: h.Fresh()}
}

// rendererRenderEchoesDefect is the planted defect of a Renderer whose
// Render returns its input unchanged.
type rendererRenderEchoesDefect struct {
	rendererOverEmbedded
}

// Render is the override of the planted defect.
func (*rendererRenderEchoesDefect) Render(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = p1
	return r0, r1
}

// rendererRenderEchoes returns the planted defect of a Renderer whose
// Render returns its input unchanged over a fresh subject.
func rendererRenderEchoes(h *RendererHandle) Renderer {
	return &rendererRenderEchoesDefect{rendererOverEmbedded: h.Fresh()}
}

// rendererClassifyPanicsDefect is the planted defect of a Renderer whose
// Classify panics.
type rendererClassifyPanicsDefect struct {
	rendererOverEmbedded
}

// Classify is the override of the planted defect.
func (*rendererClassifyPanicsDefect) Classify(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Renderer whose Classify panics")
}

// rendererClassifyPanics returns the planted defect of a Renderer whose
// Classify panics over a fresh subject.
func rendererClassifyPanics(h *RendererHandle) Renderer {
	return &rendererClassifyPanicsDefect{rendererOverEmbedded: h.Fresh()}
}

// rendererClassifyIgnoresCancelDefect is the planted defect of a Renderer
// whose Classify ignores the context.
type rendererClassifyIgnoresCancelDefect struct {
	rendererOverEmbedded
}

// Classify is the override of the planted defect.
func (d *rendererClassifyIgnoresCancelDefect) Classify(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.rendererOverEmbedded.Classify(p0, p1)
}

// rendererClassifyIgnoresCancel returns the planted defect of a Renderer
// whose Classify ignores the context over a fresh subject.
func rendererClassifyIgnoresCancel(h *RendererHandle) Renderer {
	return &rendererClassifyIgnoresCancelDefect{rendererOverEmbedded: h.Fresh()}
}

// rendererClassifyValueBesideDefect is the planted defect of a Renderer
// whose Classify returns a value beside an error.
type rendererClassifyValueBesideDefect struct {
	rendererOverEmbedded
}

// Classify is the override of the planted defect.
func (d *rendererClassifyValueBesideDefect) Classify(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	r0, r1 = d.rendererOverEmbedded.Classify(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// rendererClassifyValueBeside returns the planted defect of a Renderer
// whose Classify returns a value beside an error over a fresh subject.
func rendererClassifyValueBeside(h *RendererHandle) Renderer {
	return &rendererClassifyValueBesideDefect{rendererOverEmbedded: h.Fresh()}
}

// rendererClassifyFailsAlwaysDefect is the planted defect of a Renderer
// whose Classify fails at the edge of its domain.
type rendererClassifyFailsAlwaysDefect struct {
	rendererOverEmbedded
}

// Classify is the override of the planted defect.
func (d *rendererClassifyFailsAlwaysDefect) Classify(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.rendererOverEmbedded.Classify(p0, p1)
	if r1 == nil {
		r0 = *new(string)
		r1 = conform.ErrGenerated
	}
	return r0, r1
}

// rendererClassifyFailsAlways returns the planted defect of a Renderer
// whose Classify fails at the edge of its domain over a fresh subject.
func rendererClassifyFailsAlways(h *RendererHandle) Renderer {
	return &rendererClassifyFailsAlwaysDefect{rendererOverEmbedded: h.Fresh()}
}

// The derived checks of Lifecycle:
//
// Lifecycle.Stop/smoke: Stop returns without a panic for derived inputs.
//
// Lifecycle.Stop/repeat-survives: Stop returns without a panic when it is
// called twice.
//
// Lifecycle.Put/smoke: Put returns without a panic for derived inputs.
//
// Lifecycle.Put/cancel: Put returns an error for a cancelled context.
//
// Lifecycle.Put/idempotent: a second identical Put succeeds after the
// first.
//
// Lifecycle.Sync/smoke: Sync returns without a panic for derived inputs.
//
// Lifecycle.Sync/cancel: Sync returns an error for a cancelled context.
//
// Lifecycle.Sync/retry-succeeds: Sync succeeds within the declared
// attempts.

// The derived checks of Lifecycle that witness their laws at the shallow
// end:
//
// Lifecycle.Stop/smoke
//
// Lifecycle.Put/smoke
//
// Lifecycle.Sync/smoke
//
// Lifecycle.Sync/retry-succeeds

// The unclaimed classifications of Lifecycle:
//
// The shape:composite-writer of Lifecycle.Put has the stance silent. A
// write is an input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Lifecycle.Sync has the stance silent. A lifecycle
// call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Lifecycle writes.
type (
	LifecycleHandle = conform.Handle[Lifecycle]
	LifecycleCheck  = conform.Check[Lifecycle]
	LifecycleOption = conform.Option[Lifecycle]
)

// LifecycleChecks contains the name of each callable of Lifecycle, with
// the identities of its checks.
var LifecycleChecks = struct {
	Stop lifecycleStop
	Put  lifecyclePut
	Sync lifecycleSync
}{
	Stop: lifecycleStop{"Lifecycle.Stop"},
	Put:  lifecyclePut{"Lifecycle.Put"},
	Sync: lifecycleSync{"Lifecycle.Sync"},
}

// lifecycleStop is the name of Lifecycle.Stop, with the identities of its
// checks.
type lifecycleStop struct {
	conform.Callable[Lifecycle]
}

// Smoke returns the identity of the check that Stop returns without a
// panic for derived inputs.
func (c lifecycleStop) Smoke() conform.ID[Lifecycle] {
	return c.Check("smoke")
}

// RepeatSurvives returns the identity of the check that Stop returns
// without a panic when it is called twice.
func (c lifecycleStop) RepeatSurvives() conform.ID[Lifecycle] {
	return c.Check("repeat-survives")
}

// lifecyclePut is the name of Lifecycle.Put, with the identities of its
// checks.
type lifecyclePut struct {
	conform.Callable[Lifecycle]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c lifecyclePut) Smoke() conform.ID[Lifecycle] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c lifecyclePut) Cancel() conform.ID[Lifecycle] {
	return c.Check("cancel")
}

// Idempotent returns the identity of the check that a second identical Put
// succeeds after the first.
func (c lifecyclePut) Idempotent() conform.ID[Lifecycle] {
	return c.Check("idempotent")
}

// lifecycleSync is the name of Lifecycle.Sync, with the identities of its
// checks.
type lifecycleSync struct {
	conform.Callable[Lifecycle]
}

// Smoke returns the identity of the check that Sync returns without a
// panic for derived inputs.
func (c lifecycleSync) Smoke() conform.ID[Lifecycle] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Sync returns an error for
// a cancelled context.
func (c lifecycleSync) Cancel() conform.ID[Lifecycle] {
	return c.Check("cancel")
}

// RetrySucceeds returns the identity of the check that Sync succeeds
// within the declared attempts.
func (c lifecycleSync) RetrySucceeds() conform.ID[Lifecycle] {
	return c.Check("retry-succeeds")
}

// LifecycleSuite returns the checks of Lifecycle, and the calls and the
// budgets of its callables.
func LifecycleSuite() conform.Suite[Lifecycle] {
	return conform.Suite[Lifecycle]{
		Declaration: "Lifecycle",
		Checks: []conform.Check[Lifecycle]{
			{
				ID:          LifecycleChecks.Stop.Smoke(),
				Statement:   "Stop returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LifecycleHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Stop()
						return nil
					})
				},
				Defect: lifecycleStopPanics,
			},
			{
				ID:          LifecycleChecks.Stop.RepeatSurvives(),
				Statement:   "Stop returns without a panic when it is called twice",
				Claim:       spec.RepeatSurvives,
				Probe:       spec.SurvivesRepeat,
				Reason:      spec.PanickedRepeat,
				Law:         spec.RepeatedCallSurvives,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LifecycleHandle) {
					probe.SurvivesRepeat(h, func(ctx context.Context) error {
						h.Subject().Stop()
						return nil
					})
				},
				Defect: lifecycleStopPanicsOnRepeat,
			},
			{
				ID:          LifecycleChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LifecycleHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: lifecyclePutPanics,
			},
			{
				ID:          LifecycleChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LifecycleHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: lifecyclePutIgnoresCancel,
			},
			{
				ID:          LifecycleChecks.Put.Idempotent(),
				Statement:   "a second identical Put succeeds after the first",
				Claim:       spec.Idempotent,
				Probe:       spec.RepeatIsQuiet,
				Reason:      spec.RefusedRepeat,
				Law:         spec.RepeatIsAccepted,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LifecycleHandle) {
					probe.RepeatIsQuiet(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: lifecyclePutRefusesRepeat,
			},
			{
				ID:          LifecycleChecks.Sync.Smoke(),
				Statement:   "Sync returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LifecycleHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					})
				},
				Defect: lifecycleSyncPanics,
			},
			{
				ID:          LifecycleChecks.Sync.Cancel(),
				Statement:   "Sync returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LifecycleHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					})
				},
				Defect: lifecycleSyncIgnoresCancel,
			},
			{
				ID:          LifecycleChecks.Sync.RetrySucceeds(),
				Statement:   "Sync succeeds within the declared attempts",
				Claim:       spec.RetrySucceeds,
				Probe:       spec.SucceedsWithin,
				Reason:      spec.NeverConverged,
				Law:         spec.CallConvergesWithinTheDeclaredAttempts,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LifecycleHandle) {
					probe.SucceedsWithin(h, func(ctx context.Context) error {
						return h.Subject().Sync(ctx)
					}, 3)
				},
				Defect: lifecycleSyncFailsAlways,
			},
		},
		Calls: []conform.Call[Lifecycle]{
			{
				Calls:    lifecycleStopCalls,
				Over:     newLifecycleOver,
				Callable: LifecycleChecks.Stop.Callable,
			},
			{
				Calls:    lifecyclePutCalls,
				Over:     newLifecycleOver,
				Callable: LifecycleChecks.Put.Callable,
			},
			{
				Calls:    lifecycleSyncCalls,
				Over:     newLifecycleOver,
				Callable: LifecycleChecks.Sync.Callable,
			},
		},
	}
}

// RunLifecycleSuite applies the checks of Lifecycle, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunLifecycleSuite[T Lifecycle](t *testing.T, newSubject func() T, opts ...LifecycleOption) {
	t.Helper()
	RunLifecycleSuiteFromHandle(t, func(*LifecycleHandle) T { return newSubject() }, opts...)
}

// RunLifecycleSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLifecycleSuiteFromHandle[T Lifecycle](
	t *testing.T,
	newSubject func(h *LifecycleHandle) T,
	opts ...LifecycleOption,
) {
	t.Helper()
	conform.Run(t, LifecycleSuite(), lifecycleSubjects(newSubject), opts...)
}

// FuzzLifecycleSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLifecycleSuite[T Lifecycle](
	f *testing.F,
	callable conform.Callable[Lifecycle],
	newSubject func(h *LifecycleHandle) T,
	opts ...LifecycleOption,
) {
	f.Helper()
	conform.Fuzz(f, LifecycleSuite(), callable, lifecycleSubjects(newSubject), opts...)
}

// RunLifecycleBudgets applies the allocation ceilings of Lifecycle and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLifecycleBudgets[T Lifecycle](
	t *testing.T,
	newSubject func(h *LifecycleHandle) T,
	opts ...LifecycleOption,
) {
	t.Helper()
	conform.RunBudgets(t, LifecycleSuite(), lifecycleSubjects(newSubject), opts...)
}

// BenchLifecycleSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchLifecycleSuite[T Lifecycle](
	b conform.B,
	callable conform.Callable[Lifecycle],
	newSubject func(h *LifecycleHandle) T,
	opts ...LifecycleOption,
) {
	b.Helper()
	conform.Bench(b, LifecycleSuite(), callable, lifecycleSubjects(newSubject), opts...)
}

// lifecycleSubjects returns the factory of the subjects of Lifecycle that
// calls newSubject with the handle of each check.
func lifecycleSubjects[T Lifecycle](
	newSubject func(h *LifecycleHandle) T,
) func(h *LifecycleHandle) Lifecycle {
	return func(h *LifecycleHandle) Lifecycle { return newSubject(h) }
}

// lifecycleStopCalls generates the calls of Lifecycle.Stop with the inputs
// that each case draws.
var lifecycleStopCalls = prop.Composite(func(*prop.Case) func(context.Context, Lifecycle) error {
	return func(_ context.Context, subject Lifecycle) error {
		subject.Stop()
		return nil
	}
})

// lifecyclePutCalls generates the calls of Lifecycle.Put with the inputs
// that each case draws.
var lifecyclePutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lifecycle) error {
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Lifecycle) error {
		return subject.Put(ctx, key, value)
	}
})

// lifecycleSyncCalls generates the calls of Lifecycle.Sync with the inputs
// that each case draws.
var lifecycleSyncCalls = prop.Composite(func(*prop.Case) func(context.Context, Lifecycle) error {
	return func(ctx context.Context, subject Lifecycle) error {
		return subject.Sync(ctx)
	}
})

// lifecycleOverEmbedded is Lifecycle under the name of the field that
// lifecycleOver embeds.
type lifecycleOverEmbedded = Lifecycle

// lifecycleOver is a subject of Lifecycle whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type lifecycleOver struct {
	lifecycleOverEmbedded
	spend func()
}

// Stop calls spend and returns its zero results.
func (o lifecycleOver) Stop() {
	o.spend()
}

// Put calls spend and returns its zero results.
func (o lifecycleOver) Put(context.Context, string, string) (r0 error) {
	o.spend()
	return r0
}

// Sync calls spend and returns its zero results.
func (o lifecycleOver) Sync(context.Context) (r0 error) {
	o.spend()
	return r0
}

// newLifecycleOver returns the planted defect of a ceiling of Lifecycle,
// whose callables call spend.
func newLifecycleOver(_ *LifecycleHandle, spend func()) Lifecycle {
	return lifecycleOver{spend: spend}
}

// lifecycleStopPanicsDefect is the planted defect of a Lifecycle whose
// Stop panics.
type lifecycleStopPanicsDefect struct {
	lifecycleOverEmbedded
}

// Stop is the override of the planted defect.
func (*lifecycleStopPanicsDefect) Stop() {
	panic("the planted defect of a Lifecycle whose Stop panics")
}

// lifecycleStopPanics returns the planted defect of a Lifecycle whose Stop
// panics over a fresh subject.
func lifecycleStopPanics(h *LifecycleHandle) Lifecycle {
	return &lifecycleStopPanicsDefect{lifecycleOverEmbedded: h.Fresh()}
}

// lifecycleStopPanicsOnRepeatDefect is the planted defect of a Lifecycle
// whose Stop panics on the repeat.
type lifecycleStopPanicsOnRepeatDefect struct {
	lifecycleOverEmbedded
	calls int
}

// Stop is the override of the planted defect.
func (d *lifecycleStopPanicsOnRepeatDefect) Stop() {
	d.calls++
	if d.calls > 1 {
		panic("the planted defect of a Lifecycle whose Stop panics on the repeat")
	}
	d.lifecycleOverEmbedded.Stop()
}

// lifecycleStopPanicsOnRepeat returns the planted defect of a Lifecycle
// whose Stop panics on the repeat over a fresh subject.
func lifecycleStopPanicsOnRepeat(h *LifecycleHandle) Lifecycle {
	return &lifecycleStopPanicsOnRepeatDefect{lifecycleOverEmbedded: h.Fresh()}
}

// lifecyclePutPanicsDefect is the planted defect of a Lifecycle whose Put
// panics.
type lifecyclePutPanicsDefect struct {
	lifecycleOverEmbedded
}

// Put is the override of the planted defect.
func (*lifecyclePutPanicsDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	panic("the planted defect of a Lifecycle whose Put panics")
}

// lifecyclePutPanics returns the planted defect of a Lifecycle whose Put
// panics over a fresh subject.
func lifecyclePutPanics(h *LifecycleHandle) Lifecycle {
	return &lifecyclePutPanicsDefect{lifecycleOverEmbedded: h.Fresh()}
}

// lifecyclePutIgnoresCancelDefect is the planted defect of a Lifecycle
// whose Put ignores the context.
type lifecyclePutIgnoresCancelDefect struct {
	lifecycleOverEmbedded
}

// Put is the override of the planted defect.
func (d *lifecyclePutIgnoresCancelDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.lifecycleOverEmbedded.Put(p0, p1, p2)
}

// lifecyclePutIgnoresCancel returns the planted defect of a Lifecycle
// whose Put ignores the context over a fresh subject.
func lifecyclePutIgnoresCancel(h *LifecycleHandle) Lifecycle {
	return &lifecyclePutIgnoresCancelDefect{lifecycleOverEmbedded: h.Fresh()}
}

// lifecyclePutRefusesRepeatDefect is the planted defect of a Lifecycle
// whose Put fails on the second call.
type lifecyclePutRefusesRepeatDefect struct {
	lifecycleOverEmbedded
	calls int
}

// Put is the override of the planted defect.
func (d *lifecyclePutRefusesRepeatDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 1 {
		r0 = conform.ErrGenerated
		return r0
	}
	return d.lifecycleOverEmbedded.Put(p0, p1, p2)
}

// lifecyclePutRefusesRepeat returns the planted defect of a Lifecycle
// whose Put fails on the second call over a fresh subject.
func lifecyclePutRefusesRepeat(h *LifecycleHandle) Lifecycle {
	return &lifecyclePutRefusesRepeatDefect{lifecycleOverEmbedded: h.Fresh()}
}

// lifecycleSyncPanicsDefect is the planted defect of a Lifecycle whose
// Sync panics.
type lifecycleSyncPanicsDefect struct {
	lifecycleOverEmbedded
}

// Sync is the override of the planted defect.
func (*lifecycleSyncPanicsDefect) Sync(p0 context.Context) (r0 error) {
	panic("the planted defect of a Lifecycle whose Sync panics")
}

// lifecycleSyncPanics returns the planted defect of a Lifecycle whose Sync
// panics over a fresh subject.
func lifecycleSyncPanics(h *LifecycleHandle) Lifecycle {
	return &lifecycleSyncPanicsDefect{lifecycleOverEmbedded: h.Fresh()}
}

// lifecycleSyncIgnoresCancelDefect is the planted defect of a Lifecycle
// whose Sync ignores the context.
type lifecycleSyncIgnoresCancelDefect struct {
	lifecycleOverEmbedded
}

// Sync is the override of the planted defect.
func (d *lifecycleSyncIgnoresCancelDefect) Sync(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.lifecycleOverEmbedded.Sync(p0)
}

// lifecycleSyncIgnoresCancel returns the planted defect of a Lifecycle
// whose Sync ignores the context over a fresh subject.
func lifecycleSyncIgnoresCancel(h *LifecycleHandle) Lifecycle {
	return &lifecycleSyncIgnoresCancelDefect{lifecycleOverEmbedded: h.Fresh()}
}

// lifecycleSyncFailsAlwaysDefect is the planted defect of a Lifecycle
// whose Sync never succeeds.
type lifecycleSyncFailsAlwaysDefect struct {
	lifecycleOverEmbedded
}

// Sync is the override of the planted defect.
func (d *lifecycleSyncFailsAlwaysDefect) Sync(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	r0 = d.lifecycleOverEmbedded.Sync(p0)
	if r0 == nil {
		r0 = conform.ErrGenerated
	}
	return r0
}

// lifecycleSyncFailsAlways returns the planted defect of a Lifecycle whose
// Sync never succeeds over a fresh subject.
func lifecycleSyncFailsAlways(h *LifecycleHandle) Lifecycle {
	return &lifecycleSyncFailsAlwaysDefect{lifecycleOverEmbedded: h.Fresh()}
}

// The derived checks of Shards:
//
// Shards.Put/smoke: Put returns without a panic for derived inputs.
//
// Shards.Put/cancel: Put returns an error for a cancelled context.
//
// Shards.Put/partition: Get does not return in one partition what Put
// wrote in another.
//
// Shards.Get/smoke: Get returns without a panic for derived inputs.
//
// Shards.Get/cancel: Get returns an error for a cancelled context.
//
// Shards.Get/zero-on-error: Get returns the zero value beside an error.

// The derived checks of Shards that witness their laws at the shallow end:
//
// Shards.Put/smoke
//
// Shards.Get/smoke

// The unclaimed classifications of Shards:
//
// The shape:writer of Shards.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Shards writes.
type (
	ShardsHandle = conform.Handle[Shards]
	ShardsCheck  = conform.Check[Shards]
	ShardsOption = conform.Option[Shards]
)

// ShardsChecks contains the name of each callable of Shards, with the
// identities of its checks.
var ShardsChecks = struct {
	Put shardsPut
	Get shardsGet
}{
	Put: shardsPut{"Shards.Put"},
	Get: shardsGet{"Shards.Get"},
}

// shardsPut is the name of Shards.Put, with the identities of its checks.
type shardsPut struct {
	conform.Callable[Shards]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c shardsPut) Smoke() conform.ID[Shards] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c shardsPut) Cancel() conform.ID[Shards] {
	return c.Check("cancel")
}

// Partition returns the identity of the check that Get does not return in
// one partition what Put wrote in another.
func (c shardsPut) Partition() conform.ID[Shards] {
	return c.Check("partition")
}

// shardsGet is the name of Shards.Get, with the identities of its checks.
type shardsGet struct {
	conform.Callable[Shards]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c shardsGet) Smoke() conform.ID[Shards] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c shardsGet) Cancel() conform.ID[Shards] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c shardsGet) ZeroOnError() conform.ID[Shards] {
	return c.Check("zero-on-error")
}

// ShardsSuite returns the checks of Shards, and the calls and the budgets
// of its callables.
func ShardsSuite() conform.Suite[Shards] {
	return conform.Suite[Shards]{
		Declaration: "Shards",
		Checks: []conform.Check[Shards]{
			{
				ID:          ShardsChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShardsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: shardsPutPanics,
			},
			{
				ID:          ShardsChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShardsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Value{Key: "Key-a"})
					})
				},
				Defect: shardsPutIgnoresCancel,
			},
			{
				ID:          ShardsChecks.Put.Partition(),
				Statement:   "Get does not return in one partition what Put wrote in another",
				Claim:       spec.Partition,
				Probe:       spec.KeepsPartitionsApart,
				Reason:      spec.CrossedPartitions,
				Law:         spec.OnePartitionDoesNotServeAnother,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShardsHandle) {
					probe.KeepsPartitionsApart(h, func() func(ctx context.Context) (Value, error) {
						calls := 0
						return func(ctx context.Context) (Value, error) {
							calls++
							w := Value{Key: "Key-a"}
							if calls > 1 {
								w = Value{Key: "Key-b"}
							}
							err := h.Subject().Put(ctx, w)
							return w, err
						}
					}(), func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					})
				},
				Defect: shardsPutSucceedsAlways,
			},
			{
				ID:          ShardsChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShardsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: shardsGetPanics,
			},
			{
				ID:          ShardsChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShardsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: shardsGetIgnoresCancel,
			},
			{
				ID:          ShardsChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShardsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: shardsGetValueBeside,
			},
		},
		Calls: []conform.Call[Shards]{
			{
				Calls:    shardsPutCalls,
				Over:     newShardsOver,
				Callable: ShardsChecks.Put.Callable,
			},
			{
				Calls:    shardsGetCalls,
				Over:     newShardsOver,
				Callable: ShardsChecks.Get.Callable,
			},
		},
	}
}

// RunShardsSuite applies the checks of Shards, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunShardsSuite[T Shards](t *testing.T, newSubject func() T, opts ...ShardsOption) {
	t.Helper()
	RunShardsSuiteFromHandle(t, func(*ShardsHandle) T { return newSubject() }, opts...)
}

// RunShardsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunShardsSuiteFromHandle[T Shards](
	t *testing.T,
	newSubject func(h *ShardsHandle) T,
	opts ...ShardsOption,
) {
	t.Helper()
	conform.Run(t, ShardsSuite(), shardsSubjects(newSubject), opts...)
}

// FuzzShardsSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzShardsSuite[T Shards](
	f *testing.F,
	callable conform.Callable[Shards],
	newSubject func(h *ShardsHandle) T,
	opts ...ShardsOption,
) {
	f.Helper()
	conform.Fuzz(f, ShardsSuite(), callable, shardsSubjects(newSubject), opts...)
}

// RunShardsBudgets applies the allocation ceilings of Shards and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunShardsBudgets[T Shards](
	t *testing.T,
	newSubject func(h *ShardsHandle) T,
	opts ...ShardsOption,
) {
	t.Helper()
	conform.RunBudgets(t, ShardsSuite(), shardsSubjects(newSubject), opts...)
}

// BenchShardsSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchShardsSuite[T Shards](
	b conform.B,
	callable conform.Callable[Shards],
	newSubject func(h *ShardsHandle) T,
	opts ...ShardsOption,
) {
	b.Helper()
	conform.Bench(b, ShardsSuite(), callable, shardsSubjects(newSubject), opts...)
}

// shardsSubjects returns the factory of the subjects of Shards that calls
// newSubject with the handle of each check.
func shardsSubjects[T Shards](newSubject func(h *ShardsHandle) T) func(h *ShardsHandle) Shards {
	return func(h *ShardsHandle) Shards { return newSubject(h) }
}

// shardsPutCalls generates the calls of Shards.Put with the inputs that
// each case draws.
var shardsPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shards) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(ctx context.Context, subject Shards) error {
		return subject.Put(ctx, v)
	}
})

// shardsGetCalls generates the calls of Shards.Get with the inputs that
// each case draws.
var shardsGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shards) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Shards) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// shardsOverEmbedded is Shards under the name of the field that shardsOver
// embeds.
type shardsOverEmbedded = Shards

// shardsOver is a subject of Shards whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type shardsOver struct {
	shardsOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o shardsOver) Put(context.Context, Value) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o shardsOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// newShardsOver returns the planted defect of a ceiling of Shards, whose
// callables call spend.
func newShardsOver(_ *ShardsHandle, spend func()) Shards {
	return shardsOver{spend: spend}
}

// shardsPutPanicsDefect is the planted defect of a Shards whose Put
// panics.
type shardsPutPanicsDefect struct {
	shardsOverEmbedded
}

// Put is the override of the planted defect.
func (*shardsPutPanicsDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	panic("the planted defect of a Shards whose Put panics")
}

// shardsPutPanics returns the planted defect of a Shards whose Put panics
// over a fresh subject.
func shardsPutPanics(h *ShardsHandle) Shards {
	return &shardsPutPanicsDefect{shardsOverEmbedded: h.Fresh()}
}

// shardsPutIgnoresCancelDefect is the planted defect of a Shards whose Put
// ignores the context.
type shardsPutIgnoresCancelDefect struct {
	shardsOverEmbedded
}

// Put is the override of the planted defect.
func (d *shardsPutIgnoresCancelDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.shardsOverEmbedded.Put(p0, p1)
}

// shardsPutIgnoresCancel returns the planted defect of a Shards whose Put
// ignores the context over a fresh subject.
func shardsPutIgnoresCancel(h *ShardsHandle) Shards {
	return &shardsPutIgnoresCancelDefect{shardsOverEmbedded: h.Fresh()}
}

// shardsPutSucceedsAlwaysDefect is the planted defect of a Shards whose
// Put keeps no partition.
type shardsPutSucceedsAlwaysDefect struct {
	shardsOverEmbedded
}

// Put is the override of the planted defect.
func (*shardsPutSucceedsAlwaysDefect) Put(p0 context.Context, p1 Value) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// shardsPutSucceedsAlways returns the planted defect of a Shards whose Put
// keeps no partition over a fresh subject.
func shardsPutSucceedsAlways(h *ShardsHandle) Shards {
	return &shardsPutSucceedsAlwaysDefect{shardsOverEmbedded: h.Fresh()}
}

// shardsGetPanicsDefect is the planted defect of a Shards whose Get
// panics.
type shardsGetPanicsDefect struct {
	shardsOverEmbedded
}

// Get is the override of the planted defect.
func (*shardsGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Shards whose Get panics")
}

// shardsGetPanics returns the planted defect of a Shards whose Get panics
// over a fresh subject.
func shardsGetPanics(h *ShardsHandle) Shards {
	return &shardsGetPanicsDefect{shardsOverEmbedded: h.Fresh()}
}

// shardsGetIgnoresCancelDefect is the planted defect of a Shards whose Get
// ignores the context.
type shardsGetIgnoresCancelDefect struct {
	shardsOverEmbedded
}

// Get is the override of the planted defect.
func (d *shardsGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shardsOverEmbedded.Get(p0, p1)
}

// shardsGetIgnoresCancel returns the planted defect of a Shards whose Get
// ignores the context over a fresh subject.
func shardsGetIgnoresCancel(h *ShardsHandle) Shards {
	return &shardsGetIgnoresCancelDefect{shardsOverEmbedded: h.Fresh()}
}

// shardsGetValueBesideDefect is the planted defect of a Shards whose Get
// returns a value beside an error.
type shardsGetValueBesideDefect struct {
	shardsOverEmbedded
}

// Get is the override of the planted defect.
func (d *shardsGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.shardsOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// shardsGetValueBeside returns the planted defect of a Shards whose Get
// returns a value beside an error over a fresh subject.
func shardsGetValueBeside(h *ShardsHandle) Shards {
	return &shardsGetValueBesideDefect{shardsOverEmbedded: h.Fresh()}
}

// The derived checks of Scoped:
//
// Scoped.Set/smoke: Set returns without a panic for derived inputs.
//
// Scoped.Set/cancel: Set returns an error for a cancelled context.
//
// Scoped.Set/partition: Get does not return in one partition what Set
// wrote in another.
//
// Scoped.Get/smoke: Get returns without a panic for derived inputs.
//
// Scoped.Get/cancel: Get returns an error for a cancelled context.

// The derived checks of Scoped that witness their laws at the shallow end:
//
// Scoped.Set/smoke
//
// Scoped.Get/smoke

// The unclaimed classifications of Scoped:
//
// The shape:multi-arg-writer of Scoped.Set has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Scoped writes.
type (
	ScopedHandle = conform.Handle[Scoped]
	ScopedCheck  = conform.Check[Scoped]
	ScopedOption = conform.Option[Scoped]
)

// ScopedChecks contains the name of each callable of Scoped, with the
// identities of its checks.
var ScopedChecks = struct {
	Set scopedSet
	Get scopedGet
}{
	Set: scopedSet{"Scoped.Set"},
	Get: scopedGet{"Scoped.Get"},
}

// scopedSet is the name of Scoped.Set, with the identities of its checks.
type scopedSet struct {
	conform.Callable[Scoped]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c scopedSet) Smoke() conform.ID[Scoped] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Set returns an error for a
// cancelled context.
func (c scopedSet) Cancel() conform.ID[Scoped] {
	return c.Check("cancel")
}

// Partition returns the identity of the check that Get does not return in
// one partition what Set wrote in another.
func (c scopedSet) Partition() conform.ID[Scoped] {
	return c.Check("partition")
}

// scopedGet is the name of Scoped.Get, with the identities of its checks.
type scopedGet struct {
	conform.Callable[Scoped]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c scopedGet) Smoke() conform.ID[Scoped] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c scopedGet) Cancel() conform.ID[Scoped] {
	return c.Check("cancel")
}

// ScopedSuite returns the checks of Scoped, and the calls and the budgets
// of its callables.
func ScopedSuite() conform.Suite[Scoped] {
	return conform.Suite[Scoped]{
		Declaration: "Scoped",
		Checks: []conform.Check[Scoped]{
			{
				ID:          ScopedChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScopedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Set(ctx, string("scope-a"), string("key-a"), string("value-a"))
					})
				},
				Defect: scopedSetPanics,
			},
			{
				ID:          ScopedChecks.Set.Cancel(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Set(ctx, string("scope-a"), string("key-a"), string("value-a"))
					})
				},
				Defect: scopedSetIgnoresCancel,
			},
			{
				ID:          ScopedChecks.Set.Partition(),
				Statement:   "Get does not return in one partition what Set wrote in another",
				Claim:       spec.Partition,
				Probe:       spec.KeepsPartitionsApart,
				Reason:      spec.CrossedPartitions,
				Law:         spec.OnePartitionDoesNotServeAnother,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopedHandle) {
					probe.KeepsPartitionsApart(h, func() func(ctx context.Context) (string, error) {
						calls := 0
						return func(ctx context.Context) (string, error) {
							calls++
							w := string("value-a")
							s0 := string("scope-a")
							s1 := string("key-a")
							if calls > 1 {
								w = string("value-b")
								s0 = string("scope-b")
								s1 = string("key-b")
							}
							err := h.Subject().Set(ctx, s0, s1, w)
							return w, err
						}
					}(), func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, string("scope-a"), string("key-a"))
					})
				},
				Defect: scopedSetSucceedsAlways,
			},
			{
				ID:          ScopedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ScopedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("scope-a"), string("key-a"))
						return err
					})
				},
				Defect: scopedGetPanics,
			},
			{
				ID:          ScopedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ScopedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("scope-a"), string("key-a"))
						return err
					})
				},
				Defect: scopedGetIgnoresCancel,
			},
		},
		Calls: []conform.Call[Scoped]{
			{
				Calls:    scopedSetCalls,
				Over:     newScopedOver,
				Callable: ScopedChecks.Set.Callable,
			},
			{
				Calls:    scopedGetCalls,
				Over:     newScopedOver,
				Callable: ScopedChecks.Get.Callable,
			},
		},
	}
}

// RunScopedSuite applies the checks of Scoped, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunScopedSuite[T Scoped](t *testing.T, newSubject func() T, opts ...ScopedOption) {
	t.Helper()
	RunScopedSuiteFromHandle(t, func(*ScopedHandle) T { return newSubject() }, opts...)
}

// RunScopedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunScopedSuiteFromHandle[T Scoped](
	t *testing.T,
	newSubject func(h *ScopedHandle) T,
	opts ...ScopedOption,
) {
	t.Helper()
	conform.Run(t, ScopedSuite(), scopedSubjects(newSubject), opts...)
}

// FuzzScopedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzScopedSuite[T Scoped](
	f *testing.F,
	callable conform.Callable[Scoped],
	newSubject func(h *ScopedHandle) T,
	opts ...ScopedOption,
) {
	f.Helper()
	conform.Fuzz(f, ScopedSuite(), callable, scopedSubjects(newSubject), opts...)
}

// RunScopedBudgets applies the allocation ceilings of Scoped and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunScopedBudgets[T Scoped](
	t *testing.T,
	newSubject func(h *ScopedHandle) T,
	opts ...ScopedOption,
) {
	t.Helper()
	conform.RunBudgets(t, ScopedSuite(), scopedSubjects(newSubject), opts...)
}

// BenchScopedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchScopedSuite[T Scoped](
	b conform.B,
	callable conform.Callable[Scoped],
	newSubject func(h *ScopedHandle) T,
	opts ...ScopedOption,
) {
	b.Helper()
	conform.Bench(b, ScopedSuite(), callable, scopedSubjects(newSubject), opts...)
}

// scopedSubjects returns the factory of the subjects of Scoped that calls
// newSubject with the handle of each check.
func scopedSubjects[T Scoped](newSubject func(h *ScopedHandle) T) func(h *ScopedHandle) Scoped {
	return func(h *ScopedHandle) Scoped { return newSubject(h) }
}

// scopedSetCalls generates the calls of Scoped.Set with the inputs that
// each case draws.
var scopedSetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Scoped) error {
	scope := c.Draw(prop.Of[string](), "scope")
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Scoped) error {
		return subject.Set(ctx, scope, key, value)
	}
})

// scopedGetCalls generates the calls of Scoped.Get with the inputs that
// each case draws.
var scopedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Scoped) error {
	scope := c.Draw(prop.Of[string](), "scope")
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Scoped) error {
		_, err := subject.Get(ctx, scope, key)
		return err
	}
})

// scopedOverEmbedded is Scoped under the name of the field that scopedOver
// embeds.
type scopedOverEmbedded = Scoped

// scopedOver is a subject of Scoped whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type scopedOver struct {
	scopedOverEmbedded
	spend func()
}

// Set calls spend and returns its zero results.
func (o scopedOver) Set(context.Context, string, string, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o scopedOver) Get(context.Context, string, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newScopedOver returns the planted defect of a ceiling of Scoped, whose
// callables call spend.
func newScopedOver(_ *ScopedHandle, spend func()) Scoped {
	return scopedOver{spend: spend}
}

// scopedSetPanicsDefect is the planted defect of a Scoped whose Set
// panics.
type scopedSetPanicsDefect struct {
	scopedOverEmbedded
}

// Set is the override of the planted defect.
func (*scopedSetPanicsDefect) Set(p0 context.Context, p1, p2, p3 string) (r0 error) {
	panic("the planted defect of a Scoped whose Set panics")
}

// scopedSetPanics returns the planted defect of a Scoped whose Set panics
// over a fresh subject.
func scopedSetPanics(h *ScopedHandle) Scoped {
	return &scopedSetPanicsDefect{scopedOverEmbedded: h.Fresh()}
}

// scopedSetIgnoresCancelDefect is the planted defect of a Scoped whose Set
// ignores the context.
type scopedSetIgnoresCancelDefect struct {
	scopedOverEmbedded
}

// Set is the override of the planted defect.
func (d *scopedSetIgnoresCancelDefect) Set(p0 context.Context, p1, p2, p3 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.scopedOverEmbedded.Set(p0, p1, p2, p3)
}

// scopedSetIgnoresCancel returns the planted defect of a Scoped whose Set
// ignores the context over a fresh subject.
func scopedSetIgnoresCancel(h *ScopedHandle) Scoped {
	return &scopedSetIgnoresCancelDefect{scopedOverEmbedded: h.Fresh()}
}

// scopedSetSucceedsAlwaysDefect is the planted defect of a Scoped whose
// Set keeps no partition.
type scopedSetSucceedsAlwaysDefect struct {
	scopedOverEmbedded
}

// Set is the override of the planted defect.
func (*scopedSetSucceedsAlwaysDefect) Set(p0 context.Context, p1, p2, p3 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// scopedSetSucceedsAlways returns the planted defect of a Scoped whose Set
// keeps no partition over a fresh subject.
func scopedSetSucceedsAlways(h *ScopedHandle) Scoped {
	return &scopedSetSucceedsAlwaysDefect{scopedOverEmbedded: h.Fresh()}
}

// scopedGetPanicsDefect is the planted defect of a Scoped whose Get
// panics.
type scopedGetPanicsDefect struct {
	scopedOverEmbedded
}

// Get is the override of the planted defect.
func (*scopedGetPanicsDefect) Get(p0 context.Context, p1, p2 string) (r0 string, r1 error) {
	panic("the planted defect of a Scoped whose Get panics")
}

// scopedGetPanics returns the planted defect of a Scoped whose Get panics
// over a fresh subject.
func scopedGetPanics(h *ScopedHandle) Scoped {
	return &scopedGetPanicsDefect{scopedOverEmbedded: h.Fresh()}
}

// scopedGetIgnoresCancelDefect is the planted defect of a Scoped whose Get
// ignores the context.
type scopedGetIgnoresCancelDefect struct {
	scopedOverEmbedded
}

// Get is the override of the planted defect.
func (d *scopedGetIgnoresCancelDefect) Get(
	p0 context.Context,
	p1, p2 string,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.scopedOverEmbedded.Get(p0, p1, p2)
}

// scopedGetIgnoresCancel returns the planted defect of a Scoped whose Get
// ignores the context over a fresh subject.
func scopedGetIgnoresCancel(h *ScopedHandle) Scoped {
	return &scopedGetIgnoresCancelDefect{scopedOverEmbedded: h.Fresh()}
}

// The derived checks of Plain:
//
// Plain.Set/smoke: Set returns without a panic for derived inputs.
//
// Plain.Get/smoke: Get returns without a panic for derived inputs.
//
// Plain.Get/cancel: Get returns an error for a cancelled context.
//
// Plain.Get/zero-on-error: Get returns the zero value beside an error.
//
// Plain.Get/read-after-write: Get returns what Set stored before it.
//
// Plain.Watch/smoke: Watch returns without a panic for derived inputs.
//
// Plain.Watch/cancel: Watch returns an error for a cancelled context.
//
// Plain.All/smoke: All returns without a panic for derived inputs.

// The derived checks of Plain that witness their laws at the shallow end:
//
// Plain.Set/smoke
//
// Plain.Get/smoke
//
// Plain.Get/read-after-write
//
// Plain.Watch/smoke
//
// Plain.All/smoke

// The refused claims of Plain:
//
// Plain.Watch/zero-on-error: eidos derives no sample of a result of Watch.
// The remedy is to author a sample of the result type with the directive
// sample.

// The unclaimed classifications of Plain:
//
// The shape:mutator of Plain.Set has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:stream-reader of Plain.All has the stance gap. No claim states
// what a stream returns after its end.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Plain writes.
type (
	PlainHandle = conform.Handle[Plain]
	PlainCheck  = conform.Check[Plain]
	PlainOption = conform.Option[Plain]
)

// PlainChecks contains the name of each callable of Plain, with the
// identities of its checks.
var PlainChecks = struct {
	Set   plainSet
	Get   plainGet
	Watch plainWatch
	All   plainAll
}{
	Set:   plainSet{"Plain.Set"},
	Get:   plainGet{"Plain.Get"},
	Watch: plainWatch{"Plain.Watch"},
	All:   plainAll{"Plain.All"},
}

// plainSet is the name of Plain.Set, with the identities of its checks.
type plainSet struct {
	conform.Callable[Plain]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c plainSet) Smoke() conform.ID[Plain] {
	return c.Check("smoke")
}

// plainGet is the name of Plain.Get, with the identities of its checks.
type plainGet struct {
	conform.Callable[Plain]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c plainGet) Smoke() conform.ID[Plain] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c plainGet) Cancel() conform.ID[Plain] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c plainGet) ZeroOnError() conform.ID[Plain] {
	return c.Check("zero-on-error")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Set stored before it.
func (c plainGet) ReadAfterWrite() conform.ID[Plain] {
	return c.Check("read-after-write")
}

// plainWatch is the name of Plain.Watch, with the identities of its
// checks.
type plainWatch struct {
	conform.Callable[Plain]
}

// Smoke returns the identity of the check that Watch returns without a
// panic for derived inputs.
func (c plainWatch) Smoke() conform.ID[Plain] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Watch returns an error for
// a cancelled context.
func (c plainWatch) Cancel() conform.ID[Plain] {
	return c.Check("cancel")
}

// plainAll is the name of Plain.All, with the identities of its checks.
type plainAll struct {
	conform.Callable[Plain]
}

// Smoke returns the identity of the check that All returns without a panic
// for derived inputs.
func (c plainAll) Smoke() conform.ID[Plain] {
	return c.Check("smoke")
}

// PlainSuite returns the checks of Plain, and the calls and the budgets of
// its callables.
func PlainSuite() conform.Suite[Plain] {
	return conform.Suite[Plain]{
		Declaration: "Plain",
		Checks: []conform.Check[Plain]{
			{
				ID:          PlainChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlainHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Set(Value{Key: "Key-a"})
						return nil
					})
				},
				Defect: plainSetPanics,
			},
			{
				ID:          PlainChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlainHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: plainGetPanics,
			},
			{
				ID:          PlainChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PlainHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: plainGetIgnoresCancel,
			},
			{
				ID:          PlainChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PlainHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Value{})
				},
				Defect: plainGetValueBeside,
			},
			{
				ID:          PlainChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Set stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlainHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Value, error) {
						return h.Subject().Get(ctx, (Value{Key: "Key-a"}).Key)
					}, func(ctx context.Context) (Value, error) {
						w := Value{Key: "Key-a"}
						h.Subject().Set(w)
						return w, nil
					})
				},
				Defect: plainGetReadsAnotherValue,
			},
			{
				ID:          PlainChecks.Watch.Smoke(),
				Statement:   "Watch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlainHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Watch(ctx)
						return err
					})
				},
				Defect: plainWatchPanics,
			},
			{
				ID:          PlainChecks.Watch.Cancel(),
				Statement:   "Watch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PlainHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Watch(ctx)
						return err
					})
				},
				Defect: plainWatchIgnoresCancel,
			},
			{
				ID:          PlainChecks.All.Smoke(),
				Statement:   "All returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlainHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().All()
						return nil
					})
				},
				Defect: plainAllPanics,
			},
		},
		Calls: []conform.Call[Plain]{
			{
				Calls:    plainSetCalls,
				Over:     newPlainOver,
				Callable: PlainChecks.Set.Callable,
			},
			{
				Calls:    plainGetCalls,
				Over:     newPlainOver,
				Callable: PlainChecks.Get.Callable,
			},
			{
				Calls:    plainWatchCalls,
				Over:     newPlainOver,
				Callable: PlainChecks.Watch.Callable,
			},
			{
				Calls:    plainAllCalls,
				Over:     newPlainOver,
				Callable: PlainChecks.All.Callable,
			},
		},
	}
}

// RunPlainSuite applies the checks of Plain, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunPlainSuite[T Plain](t *testing.T, newSubject func() T, opts ...PlainOption) {
	t.Helper()
	RunPlainSuiteFromHandle(t, func(*PlainHandle) T { return newSubject() }, opts...)
}

// RunPlainSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPlainSuiteFromHandle[T Plain](
	t *testing.T,
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	t.Helper()
	conform.Run(t, PlainSuite(), plainSubjects(newSubject), opts...)
}

// FuzzPlainSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPlainSuite[T Plain](
	f *testing.F,
	callable conform.Callable[Plain],
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	f.Helper()
	conform.Fuzz(f, PlainSuite(), callable, plainSubjects(newSubject), opts...)
}

// RunPlainBudgets applies the allocation ceilings of Plain and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunPlainBudgets[T Plain](
	t *testing.T,
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	t.Helper()
	conform.RunBudgets(t, PlainSuite(), plainSubjects(newSubject), opts...)
}

// BenchPlainSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPlainSuite[T Plain](
	b conform.B,
	callable conform.Callable[Plain],
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	b.Helper()
	conform.Bench(b, PlainSuite(), callable, plainSubjects(newSubject), opts...)
}

// plainSubjects returns the factory of the subjects of Plain that calls
// newSubject with the handle of each check.
func plainSubjects[T Plain](newSubject func(h *PlainHandle) T) func(h *PlainHandle) Plain {
	return func(h *PlainHandle) Plain { return newSubject(h) }
}

// plainSetCalls generates the calls of Plain.Set with the inputs that each
// case draws.
var plainSetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Plain) error {
	v := c.Draw(prop.Of[Value](), "v")
	return func(_ context.Context, subject Plain) error {
		subject.Set(v)
		return nil
	}
})

// plainGetCalls generates the calls of Plain.Get with the inputs that each
// case draws.
var plainGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Plain) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Plain) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// plainWatchCalls generates the calls of Plain.Watch with the inputs that
// each case draws.
var plainWatchCalls = prop.Composite(func(*prop.Case) func(context.Context, Plain) error {
	return func(ctx context.Context, subject Plain) error {
		_, err := subject.Watch(ctx)
		return err
	}
})

// plainAllCalls generates the calls of Plain.All with the inputs that each
// case draws.
var plainAllCalls = prop.Composite(func(*prop.Case) func(context.Context, Plain) error {
	return func(_ context.Context, subject Plain) error {
		_ = subject.All()
		return nil
	}
})

// plainOverEmbedded is Plain under the name of the field that plainOver
// embeds.
type plainOverEmbedded = Plain

// plainOver is a subject of Plain whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type plainOver struct {
	plainOverEmbedded
	spend func()
}

// Set calls spend and returns its zero results.
func (o plainOver) Set(Value) {
	o.spend()
}

// Get calls spend and returns its zero results.
func (o plainOver) Get(context.Context, string) (r0 Value, r1 error) {
	o.spend()
	return r0, r1
}

// Watch calls spend and returns its zero results.
func (o plainOver) Watch(context.Context) (r0 <-chan Value, r1 error) {
	o.spend()
	return r0, r1
}

// All calls spend and returns its zero results.
func (o plainOver) All() (r0 iter.Seq[Value]) {
	o.spend()
	return r0
}

// newPlainOver returns the planted defect of a ceiling of Plain, whose
// callables call spend.
func newPlainOver(_ *PlainHandle, spend func()) Plain {
	return plainOver{spend: spend}
}

// plainSetPanicsDefect is the planted defect of a Plain whose Set panics.
type plainSetPanicsDefect struct {
	plainOverEmbedded
}

// Set is the override of the planted defect.
func (*plainSetPanicsDefect) Set(p0 Value) {
	panic("the planted defect of a Plain whose Set panics")
}

// plainSetPanics returns the planted defect of a Plain whose Set panics
// over a fresh subject.
func plainSetPanics(h *PlainHandle) Plain {
	return &plainSetPanicsDefect{plainOverEmbedded: h.Fresh()}
}

// plainGetPanicsDefect is the planted defect of a Plain whose Get panics.
type plainGetPanicsDefect struct {
	plainOverEmbedded
}

// Get is the override of the planted defect.
func (*plainGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	panic("the planted defect of a Plain whose Get panics")
}

// plainGetPanics returns the planted defect of a Plain whose Get panics
// over a fresh subject.
func plainGetPanics(h *PlainHandle) Plain {
	return &plainGetPanicsDefect{plainOverEmbedded: h.Fresh()}
}

// plainGetIgnoresCancelDefect is the planted defect of a Plain whose Get
// ignores the context.
type plainGetIgnoresCancelDefect struct {
	plainOverEmbedded
}

// Get is the override of the planted defect.
func (d *plainGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.plainOverEmbedded.Get(p0, p1)
}

// plainGetIgnoresCancel returns the planted defect of a Plain whose Get
// ignores the context over a fresh subject.
func plainGetIgnoresCancel(h *PlainHandle) Plain {
	return &plainGetIgnoresCancelDefect{plainOverEmbedded: h.Fresh()}
}

// plainGetValueBesideDefect is the planted defect of a Plain whose Get
// returns a value beside an error.
type plainGetValueBesideDefect struct {
	plainOverEmbedded
}

// Get is the override of the planted defect.
func (d *plainGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	r0, r1 = d.plainOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// plainGetValueBeside returns the planted defect of a Plain whose Get
// returns a value beside an error over a fresh subject.
func plainGetValueBeside(h *PlainHandle) Plain {
	return &plainGetValueBesideDefect{plainOverEmbedded: h.Fresh()}
}

// plainGetReadsAnotherValueDefect is the planted defect of a Plain whose
// Get returns a value that Set did not write.
type plainGetReadsAnotherValueDefect struct {
	plainOverEmbedded
}

// Get is the override of the planted defect.
func (d *plainGetReadsAnotherValueDefect) Get(p0 context.Context, p1 string) (r0 Value, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.plainOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Value{Key: "Key-b"}
	}
	return r0, r1
}

// plainGetReadsAnotherValue returns the planted defect of a Plain whose
// Get returns a value that Set did not write over a fresh subject.
func plainGetReadsAnotherValue(h *PlainHandle) Plain {
	return &plainGetReadsAnotherValueDefect{plainOverEmbedded: h.Fresh()}
}

// plainWatchPanicsDefect is the planted defect of a Plain whose Watch
// panics.
type plainWatchPanicsDefect struct {
	plainOverEmbedded
}

// Watch is the override of the planted defect.
func (*plainWatchPanicsDefect) Watch(p0 context.Context) (r0 <-chan Value, r1 error) {
	panic("the planted defect of a Plain whose Watch panics")
}

// plainWatchPanics returns the planted defect of a Plain whose Watch
// panics over a fresh subject.
func plainWatchPanics(h *PlainHandle) Plain {
	return &plainWatchPanicsDefect{plainOverEmbedded: h.Fresh()}
}

// plainWatchIgnoresCancelDefect is the planted defect of a Plain whose
// Watch ignores the context.
type plainWatchIgnoresCancelDefect struct {
	plainOverEmbedded
}

// Watch is the override of the planted defect.
func (d *plainWatchIgnoresCancelDefect) Watch(p0 context.Context) (r0 <-chan Value, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.plainOverEmbedded.Watch(p0)
}

// plainWatchIgnoresCancel returns the planted defect of a Plain whose
// Watch ignores the context over a fresh subject.
func plainWatchIgnoresCancel(h *PlainHandle) Plain {
	return &plainWatchIgnoresCancelDefect{plainOverEmbedded: h.Fresh()}
}

// plainAllPanicsDefect is the planted defect of a Plain whose All panics.
type plainAllPanicsDefect struct {
	plainOverEmbedded
}

// All is the override of the planted defect.
func (*plainAllPanicsDefect) All() (r0 iter.Seq[Value]) {
	panic("the planted defect of a Plain whose All panics")
}

// plainAllPanics returns the planted defect of a Plain whose All panics
// over a fresh subject.
func plainAllPanics(h *PlainHandle) Plain {
	return &plainAllPanicsDefect{plainOverEmbedded: h.Fresh()}
}

// The derived checks of Sessions:
//
// Sessions.Open/smoke: Open returns without a panic for derived inputs.
//
// Sessions.Open/cancel: Open returns an error for a cancelled context.
//
// Sessions.Open/zero-on-error: Open returns the zero value beside an
// error.

// The derived checks of Sessions that witness their laws at the shallow
// end:
//
// Sessions.Open/smoke

// The refused claims of Sessions:
//
// Sessions.Open/cursor-open: the handle that Open returns is not an
// interface, so no planted handle can replace it. The remedy is to return
// the handle of Open as an interface.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Sessions writes.
type (
	SessionsHandle = conform.Handle[Sessions]
	SessionsCheck  = conform.Check[Sessions]
	SessionsOption = conform.Option[Sessions]
)

// SessionsChecks contains the name of each callable of Sessions, with the
// identities of its checks.
var SessionsChecks = struct {
	Open sessionsOpen
}{
	Open: sessionsOpen{"Sessions.Open"},
}

// sessionsOpen is the name of Sessions.Open, with the identities of its
// checks.
type sessionsOpen struct {
	conform.Callable[Sessions]
}

// Smoke returns the identity of the check that Open returns without a
// panic for derived inputs.
func (c sessionsOpen) Smoke() conform.ID[Sessions] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Open returns an error for
// a cancelled context.
func (c sessionsOpen) Cancel() conform.ID[Sessions] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Open returns the zero
// value beside an error.
func (c sessionsOpen) ZeroOnError() conform.ID[Sessions] {
	return c.Check("zero-on-error")
}

// SessionsSuite returns the checks of Sessions, and the calls and the
// budgets of its callables.
func SessionsSuite() conform.Suite[Sessions] {
	return conform.Suite[Sessions]{
		Declaration: "Sessions",
		Checks: []conform.Check[Sessions]{
			{
				ID:          SessionsChecks.Open.Smoke(),
				Statement:   "Open returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SessionsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Open(ctx)
						return err
					})
				},
				Defect: sessionsOpenPanics,
			},
			{
				ID:          SessionsChecks.Open.Cancel(),
				Statement:   "Open returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SessionsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Open(ctx)
						return err
					})
				},
				Defect: sessionsOpenIgnoresCancel,
			},
			{
				ID:          SessionsChecks.Open.ZeroOnError(),
				Statement:   "Open returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SessionsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (*Session, error) {
						return h.Subject().Open(ctx)
					}, (*Session)(nil))
				},
				Defect: sessionsOpenValueBeside,
			},
		},
		Calls: []conform.Call[Sessions]{
			{
				Calls:    sessionsOpenCalls,
				Over:     newSessionsOver,
				Callable: SessionsChecks.Open.Callable,
			},
		},
	}
}

// RunSessionsSuite applies the checks of Sessions, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunSessionsSuite[T Sessions](t *testing.T, newSubject func() T, opts ...SessionsOption) {
	t.Helper()
	RunSessionsSuiteFromHandle(t, func(*SessionsHandle) T { return newSubject() }, opts...)
}

// RunSessionsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunSessionsSuiteFromHandle[T Sessions](
	t *testing.T,
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	t.Helper()
	conform.Run(t, SessionsSuite(), sessionsSubjects(newSubject), opts...)
}

// FuzzSessionsSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzSessionsSuite[T Sessions](
	f *testing.F,
	callable conform.Callable[Sessions],
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	f.Helper()
	conform.Fuzz(f, SessionsSuite(), callable, sessionsSubjects(newSubject), opts...)
}

// RunSessionsBudgets applies the allocation ceilings of Sessions and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunSessionsBudgets[T Sessions](
	t *testing.T,
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	t.Helper()
	conform.RunBudgets(t, SessionsSuite(), sessionsSubjects(newSubject), opts...)
}

// BenchSessionsSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchSessionsSuite[T Sessions](
	b conform.B,
	callable conform.Callable[Sessions],
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	b.Helper()
	conform.Bench(b, SessionsSuite(), callable, sessionsSubjects(newSubject), opts...)
}

// sessionsSubjects returns the factory of the subjects of Sessions that
// calls newSubject with the handle of each check.
func sessionsSubjects[T Sessions](
	newSubject func(h *SessionsHandle) T,
) func(h *SessionsHandle) Sessions {
	return func(h *SessionsHandle) Sessions { return newSubject(h) }
}

// sessionsOpenCalls generates the calls of Sessions.Open with the inputs
// that each case draws.
var sessionsOpenCalls = prop.Composite(func(*prop.Case) func(context.Context, Sessions) error {
	return func(ctx context.Context, subject Sessions) error {
		_, err := subject.Open(ctx)
		return err
	}
})

// sessionsOverEmbedded is Sessions under the name of the field that
// sessionsOver embeds.
type sessionsOverEmbedded = Sessions

// sessionsOver is a subject of Sessions whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type sessionsOver struct {
	sessionsOverEmbedded
	spend func()
}

// Open calls spend and returns its zero results.
func (o sessionsOver) Open(context.Context) (r0 *Session, r1 error) {
	o.spend()
	return r0, r1
}

// newSessionsOver returns the planted defect of a ceiling of Sessions,
// whose callables call spend.
func newSessionsOver(_ *SessionsHandle, spend func()) Sessions {
	return sessionsOver{spend: spend}
}

// sessionsOpenPanicsDefect is the planted defect of a Sessions whose Open
// panics.
type sessionsOpenPanicsDefect struct {
	sessionsOverEmbedded
}

// Open is the override of the planted defect.
func (*sessionsOpenPanicsDefect) Open(p0 context.Context) (r0 *Session, r1 error) {
	panic("the planted defect of a Sessions whose Open panics")
}

// sessionsOpenPanics returns the planted defect of a Sessions whose Open
// panics over a fresh subject.
func sessionsOpenPanics(h *SessionsHandle) Sessions {
	return &sessionsOpenPanicsDefect{sessionsOverEmbedded: h.Fresh()}
}

// sessionsOpenIgnoresCancelDefect is the planted defect of a Sessions
// whose Open ignores the context.
type sessionsOpenIgnoresCancelDefect struct {
	sessionsOverEmbedded
}

// Open is the override of the planted defect.
func (d *sessionsOpenIgnoresCancelDefect) Open(p0 context.Context) (r0 *Session, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.sessionsOverEmbedded.Open(p0)
}

// sessionsOpenIgnoresCancel returns the planted defect of a Sessions whose
// Open ignores the context over a fresh subject.
func sessionsOpenIgnoresCancel(h *SessionsHandle) Sessions {
	return &sessionsOpenIgnoresCancelDefect{sessionsOverEmbedded: h.Fresh()}
}

// sessionsOpenValueBesideDefect is the planted defect of a Sessions whose
// Open returns a value beside an error.
type sessionsOpenValueBesideDefect struct {
	sessionsOverEmbedded
}

// Open is the override of the planted defect.
func (d *sessionsOpenValueBesideDefect) Open(p0 context.Context) (r0 *Session, r1 error) {
	r0, r1 = d.sessionsOverEmbedded.Open(p0)
	if r1 != nil {
		r0 = &Session{ID: 7}
	}
	return r0, r1
}

// sessionsOpenValueBeside returns the planted defect of a Sessions whose
// Open returns a value beside an error over a fresh subject.
func sessionsOpenValueBeside(h *SessionsHandle) Sessions {
	return &sessionsOpenValueBesideDefect{sessionsOverEmbedded: h.Fresh()}
}
