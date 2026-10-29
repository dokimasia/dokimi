package durable

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Keeper writes.
type (
	KeeperHandle = conform.Handle[Keeper]
	KeeperCheck  = conform.Check[Keeper]
	KeeperOption = conform.Option[Keeper]
)

// KeeperChecks contains the name of each callable of Keeper, with the
// identities of its checks.
var KeeperChecks = struct {
	Put keeperPut
	Get keeperGet
}{
	Put: keeperPut{"Keeper.Put"},
	Get: keeperGet{"Keeper.Get"},
}

// keeperPut is the name of Keeper.Put, with the identities of its checks.
type keeperPut struct {
	conform.Callable[Keeper]
}

// LeavesTheWholeValue returns the identity of the check that a crash
// during Put leaves the whole value or no value.
func (c keeperPut) LeavesTheWholeValue() conform.ID[Keeper] {
	return c.Check("leaves-the-whole-value")
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c keeperPut) Smoke() conform.ID[Keeper] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c keeperPut) Cancel() conform.ID[Keeper] {
	return c.Check("cancel")
}

// keeperGet is the name of Keeper.Get, with the identities of its checks.
type keeperGet struct {
	conform.Callable[Keeper]
}

// KeepsEachPutAcrossARestart returns the identity of the check that Get
// returns each value that Put stored before a restart.
func (c keeperGet) KeepsEachPutAcrossARestart() conform.ID[Keeper] {
	return c.Check("keeps-each-put-across-a-restart")
}

// KeepsTheLastPut returns the identity of the check that Get returns the
// value of the last Put of a key.
func (c keeperGet) KeepsTheLastPut() conform.ID[Keeper] {
	return c.Check("keeps-the-last-put")
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c keeperGet) Smoke() conform.ID[Keeper] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c keeperGet) Cancel() conform.ID[Keeper] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c keeperGet) ZeroOnError() conform.ID[Keeper] {
	return c.Check("zero-on-error")
}

// KeeperSuite returns the checks of Keeper, and the calls and the budgets
// of its callables.
func KeeperSuite() conform.Suite[Keeper] {
	return conform.Suite[Keeper]{
		Declaration: "Keeper",
		Checks: []conform.Check[Keeper]{
			{
				ID:          KeeperChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeeperHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: keeperPutPanics,
			},
			{
				ID:          KeeperChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeeperHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: keeperPutSucceedsAlways,
			},
			{
				ID:          KeeperChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *KeeperHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: keeperGetPanics,
			},
			{
				ID:          KeeperChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeeperHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: keeperGetSucceedsAlways,
			},
			{
				ID:          KeeperChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *KeeperHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, string(""))
				},
				Defect: keeperGetValueBeside,
			},
			{
				ID:        KeeperChecks.Put.LeavesTheWholeValue(),
				Statement: "a crash during Put leaves the whole value or no value",
				Apply:     LeavesTheWholeValue,
				Defect:    tearsWrites,
				Doors:     conform.Doors[Keeper]{Reference: reference, Second: second, Crash: tear},
			},
			{
				ID:        KeeperChecks.Get.KeepsEachPutAcrossARestart(),
				Statement: "Get returns each value that Put stored before a restart",
				Apply:     KeepsEachPutAcrossARestart,
				Defect:    forgetsOnRestart,
				Doors:     conform.Doors[Keeper]{Restart: forget},
			},
			{
				ID:        KeeperChecks.Get.KeepsTheLastPut(),
				Statement: "Get returns the value of the last Put of a key",
				Property:  KeepsTheLastPut,
				Defect:    keepsTheFirstPut,
				Options:   putShares,
			},
		},
		Calls: []conform.Call[Keeper]{
			{
				Calls:    keeperPutCalls,
				Over:     newKeeperOver,
				Callable: KeeperChecks.Put.Callable,
			},
			{
				Calls:    keeperGetCalls,
				Over:     newKeeperOver,
				Callable: KeeperChecks.Get.Callable,
			},
		},
	}
}

// RunKeeperSuite applies the checks of Keeper, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunKeeperSuite[T Keeper](t *testing.T, newSubject func() T, opts ...KeeperOption) {
	t.Helper()
	RunKeeperSuiteFromHandle(t, func(*KeeperHandle) T { return newSubject() }, opts...)
}

// RunKeeperSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunKeeperSuiteFromHandle[T Keeper](
	t *testing.T,
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	t.Helper()
	conform.Run(t, KeeperSuite(), keeperSubjects(newSubject), opts...)
}

// FuzzKeeperSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzKeeperSuite[T Keeper](
	f *testing.F,
	callable conform.Callable[Keeper],
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	f.Helper()
	conform.Fuzz(f, KeeperSuite(), callable, keeperSubjects(newSubject), opts...)
}

// RunKeeperBudgets applies the allocation ceilings of Keeper and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunKeeperBudgets[T Keeper](
	t *testing.T,
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	t.Helper()
	conform.RunBudgets(t, KeeperSuite(), keeperSubjects(newSubject), opts...)
}

// BenchKeeperSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchKeeperSuite[T Keeper](
	b conform.B,
	callable conform.Callable[Keeper],
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	b.Helper()
	conform.Bench(b, KeeperSuite(), callable, keeperSubjects(newSubject), opts...)
}

// keeperSubjects returns the factory of the subjects of Keeper that calls
// newSubject with the handle of each check.
func keeperSubjects[T Keeper](newSubject func(h *KeeperHandle) T) func(h *KeeperHandle) Keeper {
	return func(h *KeeperHandle) Keeper { return newSubject(h) }
}

// keeperPutCalls generates the calls of Keeper.Put with the inputs that
// each case draws.
var keeperPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Keeper) error {
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Keeper) error {
		return subject.Put(ctx, key, value)
	}
})

// keeperGetCalls generates the calls of Keeper.Get with the inputs that
// each case draws.
var keeperGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Keeper) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Keeper) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// keeperOverEmbedded is Keeper under the name of the field that keeperOver
// embeds.
type keeperOverEmbedded = Keeper

// keeperOver is a subject of Keeper whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type keeperOver struct {
	keeperOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o keeperOver) Put(context.Context, string, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o keeperOver) Get(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newKeeperOver returns the planted defect of a ceiling of Keeper, whose
// callables call spend.
func newKeeperOver(_ *KeeperHandle, spend func()) Keeper {
	return keeperOver{spend: spend}
}

// keeperPutPanicsDefect is the planted defect of a Keeper whose Put
// panics.
type keeperPutPanicsDefect struct {
	keeperOverEmbedded
}

// Put is the override of the planted defect.
func (*keeperPutPanicsDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	panic("the planted defect of a Keeper whose Put panics")
}

// keeperPutPanics returns the planted defect of a Keeper whose Put panics
// over a fresh subject.
func keeperPutPanics(h *KeeperHandle) Keeper {
	return &keeperPutPanicsDefect{keeperOverEmbedded: h.Fresh()}
}

// keeperPutSucceedsAlwaysDefect is the planted defect of a Keeper whose
// Put ignores the context.
type keeperPutSucceedsAlwaysDefect struct {
	keeperOverEmbedded
}

// Put is the override of the planted defect.
func (*keeperPutSucceedsAlwaysDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	return r0
}

// keeperPutSucceedsAlways returns the planted defect of a Keeper whose Put
// ignores the context over a fresh subject.
func keeperPutSucceedsAlways(h *KeeperHandle) Keeper {
	return &keeperPutSucceedsAlwaysDefect{keeperOverEmbedded: h.Fresh()}
}

// keeperGetPanicsDefect is the planted defect of a Keeper whose Get
// panics.
type keeperGetPanicsDefect struct {
	keeperOverEmbedded
}

// Get is the override of the planted defect.
func (*keeperGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Keeper whose Get panics")
}

// keeperGetPanics returns the planted defect of a Keeper whose Get panics
// over a fresh subject.
func keeperGetPanics(h *KeeperHandle) Keeper {
	return &keeperGetPanicsDefect{keeperOverEmbedded: h.Fresh()}
}

// keeperGetSucceedsAlwaysDefect is the planted defect of a Keeper whose
// Get ignores the context.
type keeperGetSucceedsAlwaysDefect struct {
	keeperOverEmbedded
}

// Get is the override of the planted defect.
func (*keeperGetSucceedsAlwaysDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	return r0, r1
}

// keeperGetSucceedsAlways returns the planted defect of a Keeper whose Get
// ignores the context over a fresh subject.
func keeperGetSucceedsAlways(h *KeeperHandle) Keeper {
	return &keeperGetSucceedsAlwaysDefect{keeperOverEmbedded: h.Fresh()}
}

// keeperGetValueBesideDefect is the planted defect of a Keeper whose Get
// returns a value beside an error.
type keeperGetValueBesideDefect struct {
	keeperOverEmbedded
}

// Get is the override of the planted defect.
func (*keeperGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	r0 = string("sample-b")
	r1 = conform.ErrGenerated
	return r0, r1
}

// keeperGetValueBeside returns the planted defect of a Keeper whose Get
// returns a value beside an error over a fresh subject.
func keeperGetValueBeside(h *KeeperHandle) Keeper {
	return &keeperGetValueBesideDefect{keeperOverEmbedded: h.Fresh()}
}
