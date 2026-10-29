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
