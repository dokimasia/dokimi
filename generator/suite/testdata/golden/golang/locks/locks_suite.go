package locks

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Lock:
//
// Lock.Acquire/cancel: Acquire returns an error for a cancelled context.
//
// Lock.Release/smoke: Release returns without a panic for derived inputs.
//
// Lock.Release/repeat-survives: Release returns without a panic when it is
// called twice.
//
// Lock.Held/smoke: Held returns without a panic for derived inputs.
//
// Lock.Held/predicate: Held returns the same answer twice for an unchanged
// subject.

// The derived checks of Lock that witness their laws at the shallow end:
//
// Lock.Release/smoke
//
// Lock.Held/smoke
//
// Lock.Held/predicate

// The unclaimed classifications of Lock:
//
// The shape:lifecycle of Lock.Acquire has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The mixin:concurrent of Lock.Acquire has the stance needs. The claims of
// the mixin need concurrent callers.
//
// The mixin:waits of Lock.Acquire has the stance silent. The mixin
// documents that the callable can wait for another call.
//
// The mixin:concurrent of Lock.Release has the stance needs. The claims of
// the mixin need concurrent callers.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Lock writes.
type (
	LockHandle = conform.Handle[Lock]
	LockCheck  = conform.Check[Lock]
	LockOption = conform.Option[Lock]
)

// LockChecks contains the name of each callable of Lock, with the
// identities of its checks.
var LockChecks = struct {
	Acquire lockAcquire
	Release lockRelease
	Held    lockHeld
}{
	Acquire: lockAcquire{"Lock.Acquire"},
	Release: lockRelease{"Lock.Release"},
	Held:    lockHeld{"Lock.Held"},
}

// lockAcquire is the name of Lock.Acquire, with the identities of its
// checks.
type lockAcquire struct {
	conform.Callable[Lock]
}

// ReleaseEndsTheWait returns the identity of the check that Release ends
// the wait of Acquire.
func (c lockAcquire) ReleaseEndsTheWait() conform.ID[Lock] {
	return c.Check("release-ends-the-wait")
}

// Cancel returns the identity of the check that Acquire returns an error
// for a cancelled context.
func (c lockAcquire) Cancel() conform.ID[Lock] {
	return c.Check("cancel")
}

// lockRelease is the name of Lock.Release, with the identities of its
// checks.
type lockRelease struct {
	conform.Callable[Lock]
}

// Smoke returns the identity of the check that Release returns without a
// panic for derived inputs.
func (c lockRelease) Smoke() conform.ID[Lock] {
	return c.Check("smoke")
}

// RepeatSurvives returns the identity of the check that Release returns
// without a panic when it is called twice.
func (c lockRelease) RepeatSurvives() conform.ID[Lock] {
	return c.Check("repeat-survives")
}

// lockHeld is the name of Lock.Held, with the identities of its checks.
type lockHeld struct {
	conform.Callable[Lock]
}

// Smoke returns the identity of the check that Held returns without a
// panic for derived inputs.
func (c lockHeld) Smoke() conform.ID[Lock] {
	return c.Check("smoke")
}

// Predicate returns the identity of the check that Held returns the same
// answer twice for an unchanged subject.
func (c lockHeld) Predicate() conform.ID[Lock] {
	return c.Check("predicate")
}

// LockSuite returns the checks of Lock, and the calls and the budgets of
// its callables.
func LockSuite() conform.Suite[Lock] {
	return conform.Suite[Lock]{
		Declaration: "Lock",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Lock]{
			{
				ID:          LockChecks.Acquire.Cancel(),
				Statement:   "Acquire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx)
					})
				},
				Defect: lockAcquireIgnoresCancel,
			},
			{
				ID:          LockChecks.Release.Smoke(),
				Statement:   "Release returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Release()
						return nil
					})
				},
				Defect: lockReleasePanics,
			},
			{
				ID:          LockChecks.Release.RepeatSurvives(),
				Statement:   "Release returns without a panic when it is called twice",
				Claim:       spec.RepeatSurvives,
				Probe:       spec.SurvivesRepeat,
				Reason:      spec.PanickedRepeat,
				Law:         spec.RepeatedCallSurvives,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockHandle) {
					probe.SurvivesRepeat(h, func(ctx context.Context) error {
						h.Subject().Release()
						return nil
					})
				},
				Defect: lockReleasePanicsOnRepeat,
			},
			{
				ID:          LockChecks.Held.Smoke(),
				Statement:   "Held returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Held()
						return nil
					})
				},
				Defect: lockHeldPanics,
			},
			{
				ID:          LockChecks.Held.Predicate(),
				Statement:   "Held returns the same answer twice for an unchanged subject",
				Claim:       spec.Predicate,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (bool, error) {
						r0 := h.Subject().Held()
						return r0, nil
					})
				},
				Defect: lockHeldDrifts,
			},
			{
				ID:        LockChecks.Acquire.ReleaseEndsTheWait(),
				Statement: "Release ends the wait of Acquire",
				Apply:     ReleaseEndsTheWait,
				Defect:    neverReleases,
				Around:    bubble,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Lock.Acquire",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "mixin:concurrent",
				Callable:       "Lock.Acquire",
				Stance:         spec.Needs,
				Note:           "The claims of the mixin need concurrent callers.",
			},
			{
				Classification: "mixin:waits",
				Callable:       "Lock.Acquire",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable can wait for another call.",
			},
			{
				Classification: "mixin:concurrent",
				Callable:       "Lock.Release",
				Stance:         spec.Needs,
				Note:           "The claims of the mixin need concurrent callers.",
			},
		},
		Calls: []conform.Call[Lock]{
			{
				Calls:      lockAcquireCalls,
				Over:       newLockOver,
				Callable:   LockChecks.Acquire.Callable,
				Concurrent: true,
				Waits:      true,
			},
			{
				Calls:      lockReleaseCalls,
				Over:       newLockOver,
				Callable:   LockChecks.Release.Callable,
				Concurrent: true,
			},
			{
				Calls:    lockHeldCalls,
				Over:     newLockOver,
				Callable: LockChecks.Held.Callable,
			},
		},
	}
}

// RunLockSuite applies the checks of Lock, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunLockSuite[T Lock](t *testing.T, newSubject func() T, opts ...LockOption) {
	t.Helper()
	RunLockSuiteFromHandle(t, func(*LockHandle) T { return newSubject() }, opts...)
}

// RunLockSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLockSuiteFromHandle[T Lock](
	t *testing.T,
	newSubject func(h *LockHandle) T,
	opts ...LockOption,
) {
	t.Helper()
	conform.Run(t, LockSuite(), lockSubjects(newSubject), opts...)
}

// FuzzLockSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLockSuite[T Lock](
	f *testing.F,
	callable conform.Callable[Lock],
	newSubject func(h *LockHandle) T,
	opts ...LockOption,
) {
	f.Helper()
	conform.Fuzz(f, LockSuite(), callable, lockSubjects(newSubject), opts...)
}

// RunLockBudgets applies the allocation ceilings of Lock and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunLockBudgets[T Lock](t *testing.T, newSubject func(h *LockHandle) T, opts ...LockOption) {
	t.Helper()
	conform.RunBudgets(t, LockSuite(), lockSubjects(newSubject), opts...)
}

// BenchLockSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLockSuite[T Lock](
	b conform.B,
	callable conform.Callable[Lock],
	newSubject func(h *LockHandle) T,
	opts ...LockOption,
) {
	b.Helper()
	conform.Bench(b, LockSuite(), callable, lockSubjects(newSubject), opts...)
}

// BenchLockSuiteParallel measures callable on a subject of newSubject on
// GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchLockSuiteParallel[T Lock](
	b *testing.B,
	callable conform.Callable[Lock],
	newSubject func(h *LockHandle) T,
	opts ...LockOption,
) {
	b.Helper()
	conform.BenchParallel(b, LockSuite(), callable, lockSubjects(newSubject), opts...)
}

// lockSubjects returns the factory of the subjects of Lock that calls
// newSubject with the handle of each check.
func lockSubjects[T Lock](newSubject func(h *LockHandle) T) func(h *LockHandle) Lock {
	return func(h *LockHandle) Lock { return newSubject(h) }
}

// lockAcquireCalls generates the calls of Lock.Acquire with the inputs
// that each case draws.
var lockAcquireCalls = prop.Composite(func(*prop.Case) func(context.Context, Lock) error {
	return func(ctx context.Context, subject Lock) error {
		return subject.Acquire(ctx)
	}
})

// lockReleaseCalls generates the calls of Lock.Release with the inputs
// that each case draws.
var lockReleaseCalls = prop.Composite(func(*prop.Case) func(context.Context, Lock) error {
	return func(_ context.Context, subject Lock) error {
		subject.Release()
		return nil
	}
})

// lockHeldCalls generates the calls of Lock.Held with the inputs that each
// case draws.
var lockHeldCalls = prop.Composite(func(*prop.Case) func(context.Context, Lock) error {
	return func(_ context.Context, subject Lock) error {
		_ = subject.Held()
		return nil
	}
})

// lockOverEmbedded is Lock under the name of the field that lockOver
// embeds.
type lockOverEmbedded = Lock

// lockOver is a subject of Lock whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type lockOver struct {
	lockOverEmbedded
	spend func()
}

// Acquire calls spend and returns its zero results.
func (o lockOver) Acquire(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Release calls spend and returns its zero results.
func (o lockOver) Release() {
	o.spend()
}

// Held calls spend and returns its zero results.
func (o lockOver) Held() (r0 bool) {
	o.spend()
	return r0
}

// newLockOver returns the planted defect of a ceiling of Lock, whose
// callables call spend.
func newLockOver(_ *LockHandle, spend func()) Lock {
	return lockOver{spend: spend}
}

// lockAcquireIgnoresCancelDefect is the planted defect of a Lock whose
// Acquire ignores the context.
type lockAcquireIgnoresCancelDefect struct {
	lockOverEmbedded
}

// Acquire is the override of the planted defect.
func (d *lockAcquireIgnoresCancelDefect) Acquire(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.lockOverEmbedded.Acquire(p0)
}

// lockAcquireIgnoresCancel returns the planted defect of a Lock whose
// Acquire ignores the context over a fresh subject.
func lockAcquireIgnoresCancel(h *LockHandle) Lock {
	return &lockAcquireIgnoresCancelDefect{lockOverEmbedded: h.Fresh()}
}

// lockReleasePanicsDefect is the planted defect of a Lock whose Release
// panics.
type lockReleasePanicsDefect struct {
	lockOverEmbedded
}

// Release is the override of the planted defect.
func (*lockReleasePanicsDefect) Release() {
	panic("the planted defect of a Lock whose Release panics")
}

// lockReleasePanics returns the planted defect of a Lock whose Release
// panics over a fresh subject.
func lockReleasePanics(h *LockHandle) Lock {
	return &lockReleasePanicsDefect{lockOverEmbedded: h.Fresh()}
}

// lockReleasePanicsOnRepeatDefect is the planted defect of a Lock whose
// Release panics on the repeat.
type lockReleasePanicsOnRepeatDefect struct {
	lockOverEmbedded
	calls int
}

// Release is the override of the planted defect.
func (d *lockReleasePanicsOnRepeatDefect) Release() {
	d.calls++
	if d.calls > 1 {
		panic("the planted defect of a Lock whose Release panics on the repeat")
	}
	d.lockOverEmbedded.Release()
}

// lockReleasePanicsOnRepeat returns the planted defect of a Lock whose
// Release panics on the repeat over a fresh subject.
func lockReleasePanicsOnRepeat(h *LockHandle) Lock {
	return &lockReleasePanicsOnRepeatDefect{lockOverEmbedded: h.Fresh()}
}

// lockHeldPanicsDefect is the planted defect of a Lock whose Held panics.
type lockHeldPanicsDefect struct {
	lockOverEmbedded
}

// Held is the override of the planted defect.
func (*lockHeldPanicsDefect) Held() (r0 bool) {
	panic("the planted defect of a Lock whose Held panics")
}

// lockHeldPanics returns the planted defect of a Lock whose Held panics
// over a fresh subject.
func lockHeldPanics(h *LockHandle) Lock {
	return &lockHeldPanicsDefect{lockOverEmbedded: h.Fresh()}
}

// lockHeldDriftsDefect is the planted defect of a Lock whose Held returns
// another answer on the repeat.
type lockHeldDriftsDefect struct {
	lockOverEmbedded
	calls int
}

// Held is the override of the planted defect.
func (d *lockHeldDriftsDefect) Held() (r0 bool) {
	d.calls++
	if d.calls > 1 {
		r0 = bool(true)
	}
	return r0
}

// lockHeldDrifts returns the planted defect of a Lock whose Held returns
// another answer on the repeat over a fresh subject.
func lockHeldDrifts(h *LockHandle) Lock {
	return &lockHeldDriftsDefect{lockOverEmbedded: h.Fresh()}
}

// The derived checks of Gate:
//
// Gate.Pass/cancel: Pass returns an error for a cancelled context.

// The unclaimed classifications of Gate:
//
// The shape:lifecycle of Gate.Pass has the stance silent. A lifecycle call
// returns nothing that a check could compare.
//
// The mixin:waits of Gate.Pass has the stance silent. The mixin documents
// that the callable can wait for another call.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Gate writes.
type (
	GateHandle = conform.Handle[Gate]
	GateCheck  = conform.Check[Gate]
	GateOption = conform.Option[Gate]
)

// GateChecks contains the name of each callable of Gate, with the
// identities of its checks.
var GateChecks = struct {
	Pass gatePass
}{
	Pass: gatePass{"Gate.Pass"},
}

// gatePass is the name of Gate.Pass, with the identities of its checks.
type gatePass struct {
	conform.Callable[Gate]
}

// Cancel returns the identity of the check that Pass returns an error for
// a cancelled context.
func (c gatePass) Cancel() conform.ID[Gate] {
	return c.Check("cancel")
}

// GateSuite returns the checks of Gate, and the calls and the budgets of
// its callables.
func GateSuite() conform.Suite[Gate] {
	return conform.Suite[Gate]{
		Declaration: "Gate",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Gate]{
			{
				ID:          GateChecks.Pass.Cancel(),
				Statement:   "Pass returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *GateHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pass(ctx)
					})
				},
				Defect: gatePassIgnoresCancel,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Gate.Pass",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "mixin:waits",
				Callable:       "Gate.Pass",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable can wait for another call.",
			},
		},
		Calls: []conform.Call[Gate]{
			{
				Calls:    gatePassCalls,
				Over:     newGateOver,
				Callable: GateChecks.Pass.Callable,
				Waits:    true,
			},
		},
	}
}

// RunGateSuite applies the checks of Gate, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunGateSuite[T Gate](t *testing.T, newSubject func() T, opts ...GateOption) {
	t.Helper()
	RunGateSuiteFromHandle(t, func(*GateHandle) T { return newSubject() }, opts...)
}

// RunGateSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunGateSuiteFromHandle[T Gate](
	t *testing.T,
	newSubject func(h *GateHandle) T,
	opts ...GateOption,
) {
	t.Helper()
	conform.Run(t, GateSuite(), gateSubjects(newSubject), opts...)
}

// FuzzGateSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzGateSuite[T Gate](
	f *testing.F,
	callable conform.Callable[Gate],
	newSubject func(h *GateHandle) T,
	opts ...GateOption,
) {
	f.Helper()
	conform.Fuzz(f, GateSuite(), callable, gateSubjects(newSubject), opts...)
}

// RunGateBudgets applies the allocation ceilings of Gate and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunGateBudgets[T Gate](t *testing.T, newSubject func(h *GateHandle) T, opts ...GateOption) {
	t.Helper()
	conform.RunBudgets(t, GateSuite(), gateSubjects(newSubject), opts...)
}

// BenchGateSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchGateSuite[T Gate](
	b conform.B,
	callable conform.Callable[Gate],
	newSubject func(h *GateHandle) T,
	opts ...GateOption,
) {
	b.Helper()
	conform.Bench(b, GateSuite(), callable, gateSubjects(newSubject), opts...)
}

// gateSubjects returns the factory of the subjects of Gate that calls
// newSubject with the handle of each check.
func gateSubjects[T Gate](newSubject func(h *GateHandle) T) func(h *GateHandle) Gate {
	return func(h *GateHandle) Gate { return newSubject(h) }
}

// gatePassCalls generates the calls of Gate.Pass with the inputs that each
// case draws.
var gatePassCalls = prop.Composite(func(*prop.Case) func(context.Context, Gate) error {
	return func(ctx context.Context, subject Gate) error {
		return subject.Pass(ctx)
	}
})

// gateOverEmbedded is Gate under the name of the field that gateOver
// embeds.
type gateOverEmbedded = Gate

// gateOver is a subject of Gate whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type gateOver struct {
	gateOverEmbedded
	spend func()
}

// Pass calls spend and returns its zero results.
func (o gateOver) Pass(context.Context) (r0 error) {
	o.spend()
	return r0
}

// newGateOver returns the planted defect of a ceiling of Gate, whose
// callables call spend.
func newGateOver(_ *GateHandle, spend func()) Gate {
	return gateOver{spend: spend}
}

// gatePassIgnoresCancelDefect is the planted defect of a Gate whose Pass
// ignores the context.
type gatePassIgnoresCancelDefect struct {
	gateOverEmbedded
}

// Pass is the override of the planted defect.
func (d *gatePassIgnoresCancelDefect) Pass(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.gateOverEmbedded.Pass(p0)
}

// gatePassIgnoresCancel returns the planted defect of a Gate whose Pass
// ignores the context over a fresh subject.
func gatePassIgnoresCancel(h *GateHandle) Gate {
	return &gatePassIgnoresCancelDefect{gateOverEmbedded: h.Fresh()}
}

// The derived checks of Await:
//
// Await/cancel: Await returns an error for a cancelled context.

// The unclaimed classifications of Await:
//
// The shape:lifecycle of Await has the stance silent. A lifecycle call
// returns nothing that a check could compare.
//
// The mixin:waits of Await has the stance silent. The mixin documents that
// the callable can wait for another call.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Await writes.
type (
	AwaitHandle = conform.Handle[AwaitFunc]
	AwaitCheck  = conform.Check[AwaitFunc]
	AwaitOption = conform.Option[AwaitFunc]
)

// AwaitFunc is the type of the function Await. Each subject of the suite
// of Await has it.
type AwaitFunc func(ctx context.Context) error

// AwaitChecks contains the name of each callable of Await, with the
// identities of its checks.
var AwaitChecks = awaitCallable{"Await"}

// awaitCallable is the name of Await, with the identities of its checks.
type awaitCallable struct {
	conform.Callable[AwaitFunc]
}

// Cancel returns the identity of the check that Await returns an error for
// a cancelled context.
func (c awaitCallable) Cancel() conform.ID[AwaitFunc] {
	return c.Check("cancel")
}

// AwaitSuite returns the checks of Await, and the calls and the budgets of
// its callables.
func AwaitSuite() conform.Suite[AwaitFunc] {
	return conform.Suite[AwaitFunc]{
		Declaration: "Await",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[AwaitFunc]{
			{
				ID:          AwaitChecks.Cancel(),
				Statement:   "Await returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *AwaitHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject()(ctx)
					})
				},
				Defect: awaitIgnoresCancel,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Await",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "mixin:waits",
				Callable:       "Await",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable can wait for another call.",
			},
		},
		Calls: []conform.Call[AwaitFunc]{
			{
				Calls:    awaitCalls,
				Over:     newAwaitOver,
				Callable: AwaitChecks.Callable,
				Waits:    true,
			},
		},
	}
}

// RunAwaitSuite applies the checks of Await, and the checks that opts add,
// to subject.
func RunAwaitSuite(t *testing.T, subject AwaitFunc, opts ...AwaitOption) {
	t.Helper()
	RunAwaitSuiteFromHandle(t, awaitSubjects(subject), opts...)
}

// RunAwaitSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunAwaitSuiteFromHandle(
	t *testing.T,
	newSubject func(h *AwaitHandle) AwaitFunc,
	opts ...AwaitOption,
) {
	t.Helper()
	conform.Run(t, AwaitSuite(), newSubject, opts...)
}

// FuzzAwaitSuite registers on f the fuzz target of the property checks of
// Await on subject.
func FuzzAwaitSuite(f *testing.F, subject AwaitFunc, opts ...AwaitOption) {
	f.Helper()
	conform.Fuzz(f, AwaitSuite(), AwaitChecks.Callable, awaitSubjects(subject), opts...)
}

// RunAwaitBudgets applies the allocation ceiling of Await and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunAwaitBudgets(t *testing.T, subject AwaitFunc, opts ...AwaitOption) {
	t.Helper()
	conform.RunBudgets(t, AwaitSuite(), awaitSubjects(subject), opts...)
}

// BenchAwaitSuite measures subject under the ceilings of the budget of
// Await.
func BenchAwaitSuite(b conform.B, subject AwaitFunc, opts ...AwaitOption) {
	b.Helper()
	conform.Bench(b, AwaitSuite(), AwaitChecks.Callable, awaitSubjects(subject), opts...)
}

// awaitSubjects returns the factory of the subjects of Await that returns
// subject for each handle.
func awaitSubjects(subject AwaitFunc) func(*AwaitHandle) AwaitFunc {
	return func(*AwaitHandle) AwaitFunc { return subject }
}

// awaitCalls generates the calls of Await with the inputs that each case
// draws.
var awaitCalls = prop.Composite(func(*prop.Case) func(context.Context, AwaitFunc) error {
	return func(ctx context.Context, subject AwaitFunc) error {
		return subject(ctx)
	}
})

// newAwaitOver returns the planted defect of a ceiling of Await, which
// calls spend and returns its zero results.
func newAwaitOver(_ *AwaitHandle, spend func()) AwaitFunc {
	return func(context.Context) (r0 error) {
		spend()
		return r0
	}
}

// awaitIgnoresCancel returns the planted defect of a Await whose Await
// ignores the context over the fresh function.
func awaitIgnoresCancel(h *AwaitHandle) AwaitFunc {
	fresh := h.Fresh()
	return func(p0 context.Context) (r0 error) {
		if p0.Err() != nil {
			return r0
		}
		return fresh(p0)
	}
}

// The derived checks of Worker:
//
// Worker.Stop/smoke: Stop returns without a panic for derived inputs.
//
// Worker.Stop/repeat-survives: Stop returns without a panic when it is
// called twice.

// The derived checks of Worker that witness their laws at the shallow end:
//
// Worker.Stop/smoke

// The refused claims of Worker:
//
// Worker.Run/cancel: Run runs until its context ends. The remedy is to
// check Run with a check of the consumer that ends its context.

// The unclaimed classifications of Worker:
//
// The shape:lifecycle of Worker.Run has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The mixin:runs-until-cancelled of Worker.Run has the stance silent. The
// mixin documents that the callable runs until its context ends.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Worker writes.
type (
	WorkerHandle = conform.Handle[Worker]
	WorkerCheck  = conform.Check[Worker]
	WorkerOption = conform.Option[Worker]
)

// WorkerChecks contains the name of each callable of Worker, with the
// identities of its checks.
var WorkerChecks = struct {
	Run  workerRun
	Stop workerStop
}{
	Run:  workerRun{"Worker.Run"},
	Stop: workerStop{"Worker.Stop"},
}

// workerRun is the name of Worker.Run, with the identities of its checks.
type workerRun struct {
	conform.Callable[Worker]
}

// workerStop is the name of Worker.Stop, with the identities of its
// checks.
type workerStop struct {
	conform.Callable[Worker]
}

// Smoke returns the identity of the check that Stop returns without a
// panic for derived inputs.
func (c workerStop) Smoke() conform.ID[Worker] {
	return c.Check("smoke")
}

// RepeatSurvives returns the identity of the check that Stop returns
// without a panic when it is called twice.
func (c workerStop) RepeatSurvives() conform.ID[Worker] {
	return c.Check("repeat-survives")
}

// WorkerSuite returns the checks of Worker, and the calls and the budgets
// of its callables.
func WorkerSuite() conform.Suite[Worker] {
	return conform.Suite[Worker]{
		Declaration: "Worker",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Worker]{
			{
				ID:          WorkerChecks.Stop.Smoke(),
				Statement:   "Stop returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WorkerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Stop()
						return nil
					})
				},
				Defect: workerStopPanics,
			},
			{
				ID:          WorkerChecks.Stop.RepeatSurvives(),
				Statement:   "Stop returns without a panic when it is called twice",
				Claim:       spec.RepeatSurvives,
				Probe:       spec.SurvivesRepeat,
				Reason:      spec.PanickedRepeat,
				Law:         spec.RepeatedCallSurvives,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *WorkerHandle) {
					probe.SurvivesRepeat(h, func(ctx context.Context) error {
						h.Subject().Stop()
						return nil
					})
				},
				Defect: workerStopPanicsOnRepeat,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Worker.Run",
				Gate:     spec.ReturnsBeforeCancel,
				Why:      "Run runs until its context ends",
				Remedy:   "check Run with a check of the consumer that ends its context",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Worker.Run",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "mixin:runs-until-cancelled",
				Callable:       "Worker.Run",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable runs until its context ends.",
			},
		},
		Calls: []conform.Call[Worker]{
			{
				Calls:              workerRunCalls,
				Over:               newWorkerOver,
				Callable:           WorkerChecks.Run.Callable,
				RunsUntilCancelled: true,
			},
			{
				Calls:    workerStopCalls,
				Over:     newWorkerOver,
				Callable: WorkerChecks.Stop.Callable,
			},
		},
	}
}

// RunWorkerSuite applies the checks of Worker, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunWorkerSuite[T Worker](t *testing.T, newSubject func() T, opts ...WorkerOption) {
	t.Helper()
	RunWorkerSuiteFromHandle(t, func(*WorkerHandle) T { return newSubject() }, opts...)
}

// RunWorkerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunWorkerSuiteFromHandle[T Worker](
	t *testing.T,
	newSubject func(h *WorkerHandle) T,
	opts ...WorkerOption,
) {
	t.Helper()
	conform.Run(t, WorkerSuite(), workerSubjects(newSubject), opts...)
}

// FuzzWorkerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzWorkerSuite[T Worker](
	f *testing.F,
	callable conform.Callable[Worker],
	newSubject func(h *WorkerHandle) T,
	opts ...WorkerOption,
) {
	f.Helper()
	conform.Fuzz(f, WorkerSuite(), callable, workerSubjects(newSubject), opts...)
}

// RunWorkerBudgets applies the allocation ceilings of Worker and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunWorkerBudgets[T Worker](
	t *testing.T,
	newSubject func(h *WorkerHandle) T,
	opts ...WorkerOption,
) {
	t.Helper()
	conform.RunBudgets(t, WorkerSuite(), workerSubjects(newSubject), opts...)
}

// BenchWorkerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchWorkerSuite[T Worker](
	b conform.B,
	callable conform.Callable[Worker],
	newSubject func(h *WorkerHandle) T,
	opts ...WorkerOption,
) {
	b.Helper()
	conform.Bench(b, WorkerSuite(), callable, workerSubjects(newSubject), opts...)
}

// workerSubjects returns the factory of the subjects of Worker that calls
// newSubject with the handle of each check.
func workerSubjects[T Worker](newSubject func(h *WorkerHandle) T) func(h *WorkerHandle) Worker {
	return func(h *WorkerHandle) Worker { return newSubject(h) }
}

// workerRunCalls generates the calls of Worker.Run with the inputs that
// each case draws.
var workerRunCalls = prop.Composite(func(*prop.Case) func(context.Context, Worker) error {
	return func(ctx context.Context, subject Worker) error {
		return subject.Run(ctx)
	}
})

// workerStopCalls generates the calls of Worker.Stop with the inputs that
// each case draws.
var workerStopCalls = prop.Composite(func(*prop.Case) func(context.Context, Worker) error {
	return func(_ context.Context, subject Worker) error {
		subject.Stop()
		return nil
	}
})

// workerOverEmbedded is Worker under the name of the field that workerOver
// embeds.
type workerOverEmbedded = Worker

// workerOver is a subject of Worker whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type workerOver struct {
	workerOverEmbedded
	spend func()
}

// Run calls spend and returns its zero results.
func (o workerOver) Run(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Stop calls spend and returns its zero results.
func (o workerOver) Stop() {
	o.spend()
}

// newWorkerOver returns the planted defect of a ceiling of Worker, whose
// callables call spend.
func newWorkerOver(_ *WorkerHandle, spend func()) Worker {
	return workerOver{spend: spend}
}

// workerStopPanicsDefect is the planted defect of a Worker whose Stop
// panics.
type workerStopPanicsDefect struct {
	workerOverEmbedded
}

// Stop is the override of the planted defect.
func (*workerStopPanicsDefect) Stop() {
	panic("the planted defect of a Worker whose Stop panics")
}

// workerStopPanics returns the planted defect of a Worker whose Stop
// panics over a fresh subject.
func workerStopPanics(h *WorkerHandle) Worker {
	return &workerStopPanicsDefect{workerOverEmbedded: h.Fresh()}
}

// workerStopPanicsOnRepeatDefect is the planted defect of a Worker whose
// Stop panics on the repeat.
type workerStopPanicsOnRepeatDefect struct {
	workerOverEmbedded
	calls int
}

// Stop is the override of the planted defect.
func (d *workerStopPanicsOnRepeatDefect) Stop() {
	d.calls++
	if d.calls > 1 {
		panic("the planted defect of a Worker whose Stop panics on the repeat")
	}
	d.workerOverEmbedded.Stop()
}

// workerStopPanicsOnRepeat returns the planted defect of a Worker whose
// Stop panics on the repeat over a fresh subject.
func workerStopPanicsOnRepeat(h *WorkerHandle) Worker {
	return &workerStopPanicsOnRepeatDefect{workerOverEmbedded: h.Fresh()}
}

// The refused claims of Serve:
//
// Serve/cancel: Serve runs until its context ends. The remedy is to check
// Serve with a check of the consumer that ends its context.

// The unclaimed classifications of Serve:
//
// The shape:lifecycle of Serve has the stance silent. A lifecycle call
// returns nothing that a check could compare.
//
// The mixin:runs-until-cancelled of Serve has the stance silent. The mixin
// documents that the callable runs until its context ends.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Serve writes.
type (
	ServeHandle = conform.Handle[ServeFunc]
	ServeCheck  = conform.Check[ServeFunc]
	ServeOption = conform.Option[ServeFunc]
)

// ServeFunc is the type of the function Serve. Each subject of the suite
// of Serve has it.
type ServeFunc func(ctx context.Context) error

// ServeChecks contains the name of each callable of Serve, with the
// identities of its checks.
var ServeChecks = serveCallable{"Serve"}

// serveCallable is the name of Serve, with the identities of its checks.
type serveCallable struct {
	conform.Callable[ServeFunc]
}

// ServeSuite returns the checks of Serve, and the calls and the budgets of
// its callables.
func ServeSuite() conform.Suite[ServeFunc] {
	return conform.Suite[ServeFunc]{
		Declaration: "Serve",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Serve",
				Gate:     spec.ReturnsBeforeCancel,
				Why:      "Serve runs until its context ends",
				Remedy:   "check Serve with a check of the consumer that ends its context",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Serve",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "mixin:runs-until-cancelled",
				Callable:       "Serve",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable runs until its context ends.",
			},
		},
		Calls: []conform.Call[ServeFunc]{
			{
				Calls:              serveCalls,
				Over:               newServeOver,
				Callable:           ServeChecks.Callable,
				RunsUntilCancelled: true,
			},
		},
	}
}

// RunServeSuite applies the checks of Serve, and the checks that opts add,
// to subject.
func RunServeSuite(t *testing.T, subject ServeFunc, opts ...ServeOption) {
	t.Helper()
	RunServeSuiteFromHandle(t, serveSubjects(subject), opts...)
}

// RunServeSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunServeSuiteFromHandle(
	t *testing.T,
	newSubject func(h *ServeHandle) ServeFunc,
	opts ...ServeOption,
) {
	t.Helper()
	conform.Run(t, ServeSuite(), newSubject, opts...)
}

// FuzzServeSuite registers on f the fuzz target of the property checks of
// Serve on subject.
func FuzzServeSuite(f *testing.F, subject ServeFunc, opts ...ServeOption) {
	f.Helper()
	conform.Fuzz(f, ServeSuite(), ServeChecks.Callable, serveSubjects(subject), opts...)
}

// RunServeBudgets applies the allocation ceiling of Serve and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunServeBudgets(t *testing.T, subject ServeFunc, opts ...ServeOption) {
	t.Helper()
	conform.RunBudgets(t, ServeSuite(), serveSubjects(subject), opts...)
}

// BenchServeSuite measures subject under the ceilings of the budget of
// Serve.
func BenchServeSuite(b conform.B, subject ServeFunc, opts ...ServeOption) {
	b.Helper()
	conform.Bench(b, ServeSuite(), ServeChecks.Callable, serveSubjects(subject), opts...)
}

// serveSubjects returns the factory of the subjects of Serve that returns
// subject for each handle.
func serveSubjects(subject ServeFunc) func(*ServeHandle) ServeFunc {
	return func(*ServeHandle) ServeFunc { return subject }
}

// serveCalls generates the calls of Serve with the inputs that each case
// draws.
var serveCalls = prop.Composite(func(*prop.Case) func(context.Context, ServeFunc) error {
	return func(ctx context.Context, subject ServeFunc) error {
		return subject(ctx)
	}
})

// newServeOver returns the planted defect of a ceiling of Serve, which
// calls spend and returns its zero results.
func newServeOver(_ *ServeHandle, spend func()) ServeFunc {
	return func(context.Context) (r0 error) {
		spend()
		return r0
	}
}
