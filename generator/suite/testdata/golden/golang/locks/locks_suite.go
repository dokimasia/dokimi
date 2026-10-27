package locks

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

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

// lockRelease is the name of Lock.Release, with the identities of its
// checks.
type lockRelease struct {
	conform.Callable[Lock]
}

// lockHeld is the name of Lock.Held, with the identities of its checks.
type lockHeld struct {
	conform.Callable[Lock]
}

// LockSuite returns the checks of Lock, and the calls and the budgets of
// its callables.
func LockSuite() conform.Suite[Lock] {
	return conform.Suite[Lock]{
		Declaration: "Lock",
		Checks: []conform.Check[Lock]{
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

// GateSuite returns the checks of Gate, and the calls and the budgets of
// its callables.
func GateSuite() conform.Suite[Gate] {
	return conform.Suite[Gate]{
		Declaration: "Gate",
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

// AwaitSuite returns the checks of Await, and the calls and the budgets of
// its callables.
func AwaitSuite() conform.Suite[AwaitFunc] {
	return conform.Suite[AwaitFunc]{
		Declaration: "Await",
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
