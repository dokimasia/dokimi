package locks

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
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
// The mixin:concurrent of Lock.Acquire has the stance needs. The claim
// linearizable needs the param machineries=[reference] of the directive
// suite.
//
// The mixin:waits of Lock.Acquire has the stance silent. The mixin
// documents that the callable can wait for another call.

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
				Note:           "The claim linearizable needs the param machineries=[reference] of the directive suite.",
			},
			{
				Classification: "mixin:waits",
				Callable:       "Lock.Acquire",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable can wait for another call.",
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

// The derived checks of Board:
//
// Board.Post/smoke: Post returns without a panic for derived inputs.
//
// Board.Post/smoke#generated: Post returns without a panic for derived
// inputs.
//
// Board.Post/cancel: Post returns an error for a cancelled context.
//
// Board.Post/cancel#generated: Post returns an error for a cancelled
// context.
//
// Board.Read/smoke: Read returns without a panic for derived inputs.
//
// Board.Read/cancel: Read returns an error for a cancelled context.
//
// Board.Read/zero-on-error: Read returns the zero value beside an error.
//
// Board.Read/agrees-with-the-reference: every sequence of calls returns on
// Board what it returns on the reference.
//
// Board.Read/linearizable: every concurrent history of Board has an order
// in which the reference returns the same results.

// The derived checks of Board that witness their laws at the shallow end:
//
// Board.Post/smoke
//
// Board.Read/smoke

// The unclaimed classifications of Board:
//
// The shape:writer of Board.Post has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Board writes.
type (
	BoardHandle = conform.Handle[Board]
	BoardCheck  = conform.Check[Board]
	BoardOption = conform.Option[Board]
)

// BoardChecks contains the name of each callable of Board, with the
// identities of its checks.
var BoardChecks = struct {
	Post boardPost
	Read boardRead
}{
	Post: boardPost{"Board.Post"},
	Read: boardRead{"Board.Read"},
}

// boardPost is the name of Board.Post, with the identities of its checks.
type boardPost struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Post returns without a
// panic for derived inputs.
func (c boardPost) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Post returns
// without a panic for derived inputs.
func (c boardPost) SmokeGenerated() conform.ID[Board] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Post returns an error for
// a cancelled context.
func (c boardPost) Cancel() conform.ID[Board] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Post returns an
// error for a cancelled context.
func (c boardPost) CancelGenerated() conform.ID[Board] {
	return c.Check("cancel#generated")
}

// boardRead is the name of Board.Read, with the identities of its checks.
type boardRead struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Read returns without a
// panic for derived inputs.
func (c boardRead) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Read returns an error for
// a cancelled context.
func (c boardRead) Cancel() conform.ID[Board] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Read returns the zero
// value beside an error.
func (c boardRead) ZeroOnError() conform.ID[Board] {
	return c.Check("zero-on-error")
}

// AgreesWithTheReference returns the identity of the check that every
// sequence of calls returns on Board what it returns on the reference.
func (c boardRead) AgreesWithTheReference() conform.ID[Board] {
	return "Board/agrees-with-the-reference"
}

// Linearizable returns the identity of the check that every concurrent
// history of Board has an order in which the reference returns the same
// results.
func (c boardRead) Linearizable() conform.ID[Board] {
	return "Board/linearizable"
}

// BoardSuite returns the checks of Board, and the calls and the budgets of
// its callables.
func BoardSuite() conform.Suite[Board] {
	return conform.Suite[Board]{
		Declaration: "Board",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Board]{
			{
				ID:          BoardChecks.Post.Smoke(),
				Statement:   "Post returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, string("notice-a"))
					})
				},
				Defect: boardPostPanics,
			},
			{
				ID:          BoardChecks.Post.SmokeGenerated(),
				Statement:   "Post returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BoardHandle, c *prop.Case) {
					drawnNotice := c.Draw(prop.Of[string](), "notice")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, drawnNotice)
					})
				},
				Defect: boardPostPanics,
			},
			{
				ID:          BoardChecks.Post.Cancel(),
				Statement:   "Post returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BoardHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, string("notice-a"))
					})
				},
				Defect: boardPostIgnoresCancel,
			},
			{
				ID:          BoardChecks.Post.CancelGenerated(),
				Statement:   "Post returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *BoardHandle, c *prop.Case) {
					drawnNotice := c.Draw(prop.Of[string](), "notice")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Post(ctx, drawnNotice)
					})
				},
				Defect: boardPostIgnoresCancel,
			},
			{
				ID:          BoardChecks.Read.Smoke(),
				Statement:   "Read returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx)
						return err
					})
				},
				Defect: boardReadPanics,
			},
			{
				ID:          BoardChecks.Read.Cancel(),
				Statement:   "Read returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BoardHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx)
						return err
					})
				},
				Defect: boardReadIgnoresCancel,
			},
			{
				ID:          BoardChecks.Read.ZeroOnError(),
				Statement:   "Read returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BoardHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Read(ctx)
					}, string(""))
				},
				Defect: boardReadValueBeside,
			},
			{
				ID:        BoardChecks.Read.AgreesWithTheReference(),
				Statement: "every sequence of calls returns on Board what it returns on the reference",
				Claim:     spec.AgreesWithTheReference,
				Probe:     spec.SequencesAgreeWithTheReference,
				Reason:    spec.DivergesFromTheReference,
				Law:       spec.EverySequenceAgreesWithTheReference,
				Machineries: []spec.MachineryID{
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *BoardHandle, c *prop.Case) {
					probe.SequencesAgreeWithTheReference(h, c, boardActions, boardReference)
				},
				Defect: boardReadDivergesAfterThree,
			},
			{
				ID:        BoardChecks.Read.Linearizable(),
				Statement: "every concurrent history of Board has an order in which the reference returns the same results",
				Claim:     spec.Linearizable,
				Probe:     spec.LinearizableHistory,
				Reason:    spec.NotLinearizable,
				Law:       spec.ConcurrentHistoryIsLinearizable,
				Machineries: []spec.MachineryID{
					spec.ConcurrentCallers,
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *BoardHandle, c *prop.Case) {
					probe.LinearizableHistory(
						h,
						c,
						boardActions,
						boardReference,
						2,
						history.Budget(200000),
					)
				},
				Defect: boardRacyReadWrite,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Board.Post",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Board]{
			{
				Calls:      boardPostCalls,
				Over:       newBoardOver,
				Callable:   BoardChecks.Post.Callable,
				Concurrent: true,
			},
			{
				Calls:      boardReadCalls,
				Over:       newBoardOver,
				Callable:   BoardChecks.Read.Callable,
				Concurrent: true,
			},
		},
	}
}

// RunBoardSuite applies the checks of Board, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunBoardSuite[T Board](t *testing.T, newSubject func() T, opts ...BoardOption) {
	t.Helper()
	RunBoardSuiteFromHandle(t, func(*BoardHandle) T { return newSubject() }, opts...)
}

// RunBoardSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunBoardSuiteFromHandle[T Board](
	t *testing.T,
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	t.Helper()
	conform.Run(t, BoardSuite(), boardSubjects(newSubject), opts...)
}

// FuzzBoardSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzBoardSuite[T Board](
	f *testing.F,
	callable conform.Callable[Board],
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	f.Helper()
	conform.Fuzz(f, BoardSuite(), callable, boardSubjects(newSubject), opts...)
}

// RunBoardBudgets applies the allocation ceilings of Board and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunBoardBudgets[T Board](
	t *testing.T,
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	t.Helper()
	conform.RunBudgets(t, BoardSuite(), boardSubjects(newSubject), opts...)
}

// BenchBoardSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchBoardSuite[T Board](
	b conform.B,
	callable conform.Callable[Board],
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	b.Helper()
	conform.Bench(b, BoardSuite(), callable, boardSubjects(newSubject), opts...)
}

// BenchBoardSuiteParallel measures callable on a subject of newSubject on
// GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchBoardSuiteParallel[T Board](
	b *testing.B,
	callable conform.Callable[Board],
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	b.Helper()
	conform.BenchParallel(b, BoardSuite(), callable, boardSubjects(newSubject), opts...)
}

// boardSubjects returns the factory of the subjects of Board that calls
// newSubject with the handle of each check.
func boardSubjects[T Board](newSubject func(h *BoardHandle) T) func(h *BoardHandle) Board {
	return func(h *BoardHandle) Board { return newSubject(h) }
}

// boardPostCalls generates the calls of Board.Post with the inputs that
// each case draws.
var boardPostCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	notice := c.Draw(prop.Of[string](), "notice")
	return func(ctx context.Context, subject Board) error {
		return subject.Post(ctx, notice)
	}
})

// boardReadCalls generates the calls of Board.Read with the inputs that
// each case draws.
var boardReadCalls = prop.Composite(func(*prop.Case) func(context.Context, Board) error {
	return func(ctx context.Context, subject Board) error {
		_, err := subject.Read(ctx)
		return err
	}
})

// boardOverEmbedded is Board under the name of the field that boardOver
// embeds.
type boardOverEmbedded = Board

// boardOver is a subject of Board whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type boardOver struct {
	boardOverEmbedded
	spend func()
}

// Post calls spend and returns its zero results.
func (o boardOver) Post(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Read calls spend and returns its zero results.
func (o boardOver) Read(context.Context) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newBoardOver returns the planted defect of a ceiling of Board, whose
// callables call spend.
func newBoardOver(_ *BoardHandle, spend func()) Board {
	return boardOver{spend: spend}
}

// boardPostPanicsDefect is the planted defect of a Board whose Post
// panics.
type boardPostPanicsDefect struct {
	boardOverEmbedded
}

// Post is the override of the planted defect.
func (*boardPostPanicsDefect) Post(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Board whose Post panics")
}

// boardPostPanics returns the planted defect of a Board whose Post panics
// over a fresh subject.
func boardPostPanics(h *BoardHandle) Board {
	return &boardPostPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardPostIgnoresCancelDefect is the planted defect of a Board whose Post
// ignores the context.
type boardPostIgnoresCancelDefect struct {
	boardOverEmbedded
}

// Post is the override of the planted defect.
func (d *boardPostIgnoresCancelDefect) Post(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.boardOverEmbedded.Post(p0, p1)
}

// boardPostIgnoresCancel returns the planted defect of a Board whose Post
// ignores the context over a fresh subject.
func boardPostIgnoresCancel(h *BoardHandle) Board {
	return &boardPostIgnoresCancelDefect{boardOverEmbedded: h.Fresh()}
}

// boardReadPanicsDefect is the planted defect of a Board whose Read
// panics.
type boardReadPanicsDefect struct {
	boardOverEmbedded
}

// Read is the override of the planted defect.
func (*boardReadPanicsDefect) Read(p0 context.Context) (r0 string, r1 error) {
	panic("the planted defect of a Board whose Read panics")
}

// boardReadPanics returns the planted defect of a Board whose Read panics
// over a fresh subject.
func boardReadPanics(h *BoardHandle) Board {
	return &boardReadPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardReadIgnoresCancelDefect is the planted defect of a Board whose Read
// ignores the context.
type boardReadIgnoresCancelDefect struct {
	boardOverEmbedded
}

// Read is the override of the planted defect.
func (d *boardReadIgnoresCancelDefect) Read(p0 context.Context) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.boardOverEmbedded.Read(p0)
}

// boardReadIgnoresCancel returns the planted defect of a Board whose Read
// ignores the context over a fresh subject.
func boardReadIgnoresCancel(h *BoardHandle) Board {
	return &boardReadIgnoresCancelDefect{boardOverEmbedded: h.Fresh()}
}

// boardReadValueBesideDefect is the planted defect of a Board whose Read
// returns a value beside an error.
type boardReadValueBesideDefect struct {
	boardOverEmbedded
}

// Read is the override of the planted defect.
func (d *boardReadValueBesideDefect) Read(p0 context.Context) (r0 string, r1 error) {
	r0, r1 = d.boardOverEmbedded.Read(p0)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// boardReadValueBeside returns the planted defect of a Board whose Read
// returns a value beside an error over a fresh subject.
func boardReadValueBeside(h *BoardHandle) Board {
	return &boardReadValueBesideDefect{boardOverEmbedded: h.Fresh()}
}

// boardReadDivergesAfterThreeDefect is the planted defect of a Board whose
// first callable with a value result returns another value from its fourth
// call on.
type boardReadDivergesAfterThreeDefect struct {
	boardOverEmbedded
	calls int
}

// Read is the override of the planted defect.
func (d *boardReadDivergesAfterThreeDefect) Read(p0 context.Context) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.boardOverEmbedded.Read(p0)
	d.calls++
	if d.calls > 3 {
		r0 = string("sample-b")
	}
	return r0, r1
}

// boardReadDivergesAfterThree returns the planted defect of a Board whose
// first callable with a value result returns another value from its fourth
// call on over a fresh subject.
func boardReadDivergesAfterThree(h *BoardHandle) Board {
	return &boardReadDivergesAfterThreeDefect{boardOverEmbedded: h.Fresh()}
}

// boardRacyReadWriteDefect is the planted defect of a Board whose first
// callable with a value result returns another value while another call is
// in progress.
type boardRacyReadWriteDefect struct {
	boardOverEmbedded
	h      *BoardHandle
	marked atomic.Int32
}

// Post is the override of the planted defect.
func (d *boardRacyReadWriteDefect) Post(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	return d.boardOverEmbedded.Post(p0, p1)
}

// Read is the override of the planted defect.
func (d *boardRacyReadWriteDefect) Read(p0 context.Context) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.marked.Load() > 0 {
		r0 = string("sample-b")
		return r0, r1
	}
	return d.boardOverEmbedded.Read(p0)
}

// boardRacyReadWrite returns the planted defect of a Board whose first
// callable with a value result returns another value while another call is
// in progress over a fresh subject.
func boardRacyReadWrite(h *BoardHandle) Board {
	return &boardRacyReadWriteDefect{boardOverEmbedded: h.Fresh(), h: h}
}

// boardActions returns the actions of the generated sequences of Board
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func boardActions(c *prop.Case, subject func() Board) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Post",
			Input: func(c *prop.Case, _ struct{}) any {
				notice := c.Draw(prop.Of[string](), "post.notice")
				return []any{notice}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Post", args)
				values, err := boardPostCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Read",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Read", args)
				values, err := boardReadCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// boardPostCall calls Post on subject with the arguments args, and returns
// the values of the call beside its error.
func boardPostCall(ctx context.Context, subject Board, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Post(ctx, a0)
}

// boardReadCall calls Read on subject with the arguments args, and returns
// the values of the call beside its error.
func boardReadCall(ctx context.Context, subject Board, args []any) ([]any, error) {
	r0, err := subject.Read(ctx)
	return []any{r0}, err
}

// boardReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func boardReference(r Board) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Post":
			values, err = boardPostCall(ctx, r, args)
		case "Read":
			values, err = boardReadCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Shelf:
//
// Shelf.Add/smoke: Add returns without a panic for derived inputs.
//
// Shelf.Add/smoke#generated: Add returns without a panic for derived
// inputs.
//
// Shelf.Add/cancel: Add returns an error for a cancelled context.
//
// Shelf.Add/cancel#generated: Add returns an error for a cancelled
// context.
//
// Shelf.Count/smoke: Count returns without a panic for derived inputs.
//
// Shelf.Count/smoke#generated: Count returns without a panic for derived
// inputs.
//
// Shelf.Count/cancel: Count returns an error for a cancelled context.
//
// Shelf.Count/cancel#generated: Count returns an error for a cancelled
// context.
//
// Shelf.Count/zero-on-error: Count returns the zero value beside an error.
//
// Shelf.Count/zero-on-error#generated: Count returns the zero value beside
// an error.
//
// Shelf.Count/reads-agree-under-concurrency: Count returns beside other
// reads what it returns with the same arguments when it runs alone.

// The derived checks of Shelf that witness their laws at the shallow end:
//
// Shelf.Add/smoke
//
// Shelf.Count/smoke

// The unclaimed classifications of Shelf:
//
// The shape:writer of Shelf.Add has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Shelf writes.
type (
	ShelfHandle = conform.Handle[Shelf]
	ShelfCheck  = conform.Check[Shelf]
	ShelfOption = conform.Option[Shelf]
)

// ShelfChecks contains the name of each callable of Shelf, with the
// identities of its checks.
var ShelfChecks = struct {
	Add   shelfAdd
	Count shelfCount
}{
	Add:   shelfAdd{"Shelf.Add"},
	Count: shelfCount{"Shelf.Count"},
}

// shelfAdd is the name of Shelf.Add, with the identities of its checks.
type shelfAdd struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c shelfAdd) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Add returns
// without a panic for derived inputs.
func (c shelfAdd) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c shelfAdd) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Add returns an
// error for a cancelled context.
func (c shelfAdd) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// shelfCount is the name of Shelf.Count, with the identities of its
// checks.
type shelfCount struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Count returns without a
// panic for derived inputs.
func (c shelfCount) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Count returns
// without a panic for derived inputs.
func (c shelfCount) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Count returns an error for
// a cancelled context.
func (c shelfCount) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Count returns an
// error for a cancelled context.
func (c shelfCount) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Count returns the
// zero value beside an error.
func (c shelfCount) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Count
// returns the zero value beside an error.
func (c shelfCount) ZeroOnErrorGenerated() conform.ID[Shelf] {
	return c.Check("zero-on-error#generated")
}

// ReadsAgreeUnderConcurrency returns the identity of the check that Count
// returns beside other reads what it returns with the same arguments when
// it runs alone.
func (c shelfCount) ReadsAgreeUnderConcurrency() conform.ID[Shelf] {
	return c.Check("reads-agree-under-concurrency")
}

// ShelfSuite returns the checks of Shelf, and the calls and the budgets of
// its callables.
func ShelfSuite() conform.Suite[Shelf] {
	return conform.Suite[Shelf]{
		Declaration: "Shelf",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Shelf]{
			{
				ID:          ShelfChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("key-a"))
					})
				},
				Defect: shelfAddPanics,
			},
			{
				ID:          ShelfChecks.Add.SmokeGenerated(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnKey)
					})
				},
				Defect: shelfAddPanics,
			},
			{
				ID:          ShelfChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("key-a"))
					})
				},
				Defect: shelfAddIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Add.CancelGenerated(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnKey)
					})
				},
				Defect: shelfAddIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Count.Smoke(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx, string("key-a"))
						return err
					})
				},
				Defect: shelfCountPanics,
			},
			{
				ID:          ShelfChecks.Count.SmokeGenerated(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx, drawnKey)
						return err
					})
				},
				Defect: shelfCountPanics,
			},
			{
				ID:          ShelfChecks.Count.Cancel(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx, string("key-a"))
						return err
					})
				},
				Defect: shelfCountIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Count.CancelGenerated(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx, drawnKey)
						return err
					})
				},
				Defect: shelfCountIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Count.ZeroOnError(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx, string("key-a"))
					}, int(0))
				},
				Defect: shelfCountValueBeside,
			},
			{
				ID:          ShelfChecks.Count.ZeroOnErrorGenerated(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx, drawnKey)
					}, int(0))
				},
				Defect: shelfCountValueBeside,
			},
			{
				ID:          ShelfChecks.Count.ReadsAgreeUnderConcurrency(),
				Statement:   "Count returns beside other reads what it returns with the same arguments when it runs alone",
				Claim:       spec.ReadsAgreeUnderConcurrency,
				Probe:       spec.ConcurrentReadsAgree,
				Reason:      spec.InconsistentConcurrentRead,
				Law:         spec.ReadsBesideReadsReturnWhatOneReadReturns,
				Machineries: []spec.MachineryID{spec.ConcurrentCallers, spec.GeneratedSequences},
				Property: func(h *ShelfHandle, c *prop.Case) {
					probe.ConcurrentReadsAgree(
						h,
						c,
						shelfActions,
						"Count",
						func(ctx context.Context, args []any) conform.Output {
							values, err := shelfCountCall(ctx, h.Subject(), []any{args[0]})
							return conform.OutputOf(err, values...)
						},
						2,
					)
				},
				Defect: shelfCountRacyRead,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Shelf.Add",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Shelf]{
			{
				Calls:    shelfAddCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Add.Callable,
			},
			{
				Calls:      shelfCountCalls,
				Over:       newShelfOver,
				Callable:   ShelfChecks.Count.Callable,
				Concurrent: true,
			},
		},
	}
}

// RunShelfSuite applies the checks of Shelf, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunShelfSuite[T Shelf](t *testing.T, newSubject func() T, opts ...ShelfOption) {
	t.Helper()
	RunShelfSuiteFromHandle(t, func(*ShelfHandle) T { return newSubject() }, opts...)
}

// RunShelfSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunShelfSuiteFromHandle[T Shelf](
	t *testing.T,
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	t.Helper()
	conform.Run(t, ShelfSuite(), shelfSubjects(newSubject), opts...)
}

// FuzzShelfSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzShelfSuite[T Shelf](
	f *testing.F,
	callable conform.Callable[Shelf],
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	f.Helper()
	conform.Fuzz(f, ShelfSuite(), callable, shelfSubjects(newSubject), opts...)
}

// RunShelfBudgets applies the allocation ceilings of Shelf and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunShelfBudgets[T Shelf](
	t *testing.T,
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	t.Helper()
	conform.RunBudgets(t, ShelfSuite(), shelfSubjects(newSubject), opts...)
}

// BenchShelfSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchShelfSuite[T Shelf](
	b conform.B,
	callable conform.Callable[Shelf],
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	b.Helper()
	conform.Bench(b, ShelfSuite(), callable, shelfSubjects(newSubject), opts...)
}

// BenchShelfSuiteParallel measures callable on a subject of newSubject on
// GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchShelfSuiteParallel[T Shelf](
	b *testing.B,
	callable conform.Callable[Shelf],
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	b.Helper()
	conform.BenchParallel(b, ShelfSuite(), callable, shelfSubjects(newSubject), opts...)
}

// shelfSubjects returns the factory of the subjects of Shelf that calls
// newSubject with the handle of each check.
func shelfSubjects[T Shelf](newSubject func(h *ShelfHandle) T) func(h *ShelfHandle) Shelf {
	return func(h *ShelfHandle) Shelf { return newSubject(h) }
}

// shelfAddCalls generates the calls of Shelf.Add with the inputs that each
// case draws.
var shelfAddCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Shelf) error {
		return subject.Add(ctx, key)
	}
})

// shelfCountCalls generates the calls of Shelf.Count with the inputs that
// each case draws.
var shelfCountCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.Count(ctx, key)
		return err
	}
})

// shelfOverEmbedded is Shelf under the name of the field that shelfOver
// embeds.
type shelfOverEmbedded = Shelf

// shelfOver is a subject of Shelf whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type shelfOver struct {
	shelfOverEmbedded
	spend func()
}

// Add calls spend and returns its zero results.
func (o shelfOver) Add(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o shelfOver) Count(context.Context, string) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newShelfOver returns the planted defect of a ceiling of Shelf, whose
// callables call spend.
func newShelfOver(_ *ShelfHandle, spend func()) Shelf {
	return shelfOver{spend: spend}
}

// shelfAddPanicsDefect is the planted defect of a Shelf whose Add panics.
type shelfAddPanicsDefect struct {
	shelfOverEmbedded
}

// Add is the override of the planted defect.
func (*shelfAddPanicsDefect) Add(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Shelf whose Add panics")
}

// shelfAddPanics returns the planted defect of a Shelf whose Add panics
// over a fresh subject.
func shelfAddPanics(h *ShelfHandle) Shelf {
	return &shelfAddPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAddIgnoresCancelDefect is the planted defect of a Shelf whose Add
// ignores the context.
type shelfAddIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Add is the override of the planted defect.
func (d *shelfAddIgnoresCancelDefect) Add(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.shelfOverEmbedded.Add(p0, p1)
}

// shelfAddIgnoresCancel returns the planted defect of a Shelf whose Add
// ignores the context over a fresh subject.
func shelfAddIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfAddIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfCountPanicsDefect is the planted defect of a Shelf whose Count
// panics.
type shelfCountPanicsDefect struct {
	shelfOverEmbedded
}

// Count is the override of the planted defect.
func (*shelfCountPanicsDefect) Count(p0 context.Context, p1 string) (r0 int, r1 error) {
	panic("the planted defect of a Shelf whose Count panics")
}

// shelfCountPanics returns the planted defect of a Shelf whose Count
// panics over a fresh subject.
func shelfCountPanics(h *ShelfHandle) Shelf {
	return &shelfCountPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfCountIgnoresCancelDefect is the planted defect of a Shelf whose
// Count ignores the context.
type shelfCountIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Count is the override of the planted defect.
func (d *shelfCountIgnoresCancelDefect) Count(p0 context.Context, p1 string) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.Count(p0, p1)
}

// shelfCountIgnoresCancel returns the planted defect of a Shelf whose
// Count ignores the context over a fresh subject.
func shelfCountIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfCountIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfCountValueBesideDefect is the planted defect of a Shelf whose Count
// returns a value beside an error.
type shelfCountValueBesideDefect struct {
	shelfOverEmbedded
}

// Count is the override of the planted defect.
func (d *shelfCountValueBesideDefect) Count(p0 context.Context, p1 string) (r0 int, r1 error) {
	r0, r1 = d.shelfOverEmbedded.Count(p0, p1)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// shelfCountValueBeside returns the planted defect of a Shelf whose Count
// returns a value beside an error over a fresh subject.
func shelfCountValueBeside(h *ShelfHandle) Shelf {
	return &shelfCountValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfCountRacyReadDefect is the planted defect of a Shelf whose Count
// returns the results of another read that overlaps it.
type shelfCountRacyReadDefect struct {
	shelfOverEmbedded
	h    *ShelfHandle
	last atomic.Value
}

// Count is the override of the planted defect.
func (d *shelfCountRacyReadDefect) Count(p0 context.Context, p1 string) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.shelfOverEmbedded.Count(p0, p1)
	d.last.Store([]any{r0, r1})
	d.h.Yield()
	kept, _ := d.last.Load().([]any)
	r0, _ = kept[0].(int)
	r1, _ = kept[1].(error)
	return r0, r1
}

// shelfCountRacyRead returns the planted defect of a Shelf whose Count
// returns the results of another read that overlaps it over a fresh
// subject.
func shelfCountRacyRead(h *ShelfHandle) Shelf {
	return &shelfCountRacyReadDefect{shelfOverEmbedded: h.Fresh(), h: h}
}

// shelfActions returns the actions of the generated sequences of Shelf
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func shelfActions(c *prop.Case, subject func() Shelf) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Add",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "add.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Add", args)
				values, err := shelfAddCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Count",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "count.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Count", args, args[0])
				values, err := shelfCountCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// shelfAddCall calls Add on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfAddCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Add(ctx, a0)
}

// shelfCountCall calls Count on subject with the arguments args, and
// returns the values of the call beside its error.
func shelfCountCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Count(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Flight:
//
// Flight.Do/smoke: Do returns without a panic for derived inputs.
//
// Flight.Do/smoke#generated: Do returns without a panic for derived
// inputs.
//
// Flight.Do/cancel: Do returns an error for a cancelled context.
//
// Flight.Do/cancel#generated: Do returns an error for a cancelled context.
//
// Flight.Do/failure-is-reported: Do returns an error when a call of its
// dependency fails.
//
// Flight.Do/failure-is-reported#generated: Do returns an error when a call
// of its dependency fails.
//
// Flight.Do/computes-once: calls of Do with one key at once share one
// computation of their result.
//
// Flight.Do/computes-once#generated: calls of Do with one key at once
// share one computation of their result.

// The derived checks of Flight that witness their laws at the shallow end:
//
// Flight.Do/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Flight writes.
type (
	FlightHandle = conform.Handle[Flight]
	FlightCheck  = conform.Check[Flight]
	FlightOption = conform.Option[Flight]
)

// FlightChecks contains the name of each callable of Flight, with the
// identities of its checks.
var FlightChecks = struct {
	Do flightDo
}{
	Do: flightDo{"Flight.Do"},
}

// flightDo is the name of Flight.Do, with the identities of its checks.
type flightDo struct {
	conform.Callable[Flight]
}

// Smoke returns the identity of the check that Do returns without a panic
// for derived inputs.
func (c flightDo) Smoke() conform.ID[Flight] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Do returns without
// a panic for derived inputs.
func (c flightDo) SmokeGenerated() conform.ID[Flight] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Do returns an error for a
// cancelled context.
func (c flightDo) Cancel() conform.ID[Flight] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Do returns an
// error for a cancelled context.
func (c flightDo) CancelGenerated() conform.ID[Flight] {
	return c.Check("cancel#generated")
}

// FailureIsReported returns the identity of the check that Do returns an
// error when a call of its dependency fails.
func (c flightDo) FailureIsReported() conform.ID[Flight] {
	return c.Check("failure-is-reported")
}

// FailureIsReportedGenerated returns the identity of the check that Do
// returns an error when a call of its dependency fails.
func (c flightDo) FailureIsReportedGenerated() conform.ID[Flight] {
	return c.Check("failure-is-reported#generated")
}

// ComputesOnce returns the identity of the check that calls of Do with one
// key at once share one computation of their result.
func (c flightDo) ComputesOnce() conform.ID[Flight] {
	return c.Check("computes-once")
}

// ComputesOnceGenerated returns the identity of the check that calls of Do
// with one key at once share one computation of their result.
func (c flightDo) ComputesOnceGenerated() conform.ID[Flight] {
	return c.Check("computes-once#generated")
}

// FlightSuite returns the checks of Flight, and the calls and the budgets
// of its callables.
func FlightSuite() conform.Suite[Flight] {
	return conform.Suite[Flight]{
		Declaration: "Flight",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Flight]{
			{
				ID:          FlightChecks.Do.Smoke(),
				Statement:   "Do returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FlightHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Do(ctx, string("key-a"), Loader(func(ctx context.Context, key string) (string, error) { return "", nil }))
						return err
					})
				},
				Defect: flightDoPanics,
			},
			{
				ID:          FlightChecks.Do.SmokeGenerated(),
				Statement:   "Do returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FlightHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					loadResults0 := c.Draw(prop.List(prop.Of[string]()), "load-0")
					loadResults1 := c.Draw(
						prop.List(
							prop.SampledFrom[error](nil, conform.ErrGenerated),
							prop.MinSize(len(loadResults0)),
							prop.MaxSize(len(loadResults0)),
						),
						"load-1",
					)
					drawnLoad := func() Loader {
						var loadAt atomic.Int64
						return func(context.Context, string) (r0 string, r1 error) {
							i := int(loadAt.Add(1) - 1)
							if i >= len(loadResults0) {
								r1 = conform.ErrGenerated
								return r0, r1
							}
							return loadResults0[i],
								loadResults1[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Do(ctx, drawnKey, drawnLoad())
						return err
					})
				},
				Defect: flightDoPanics,
			},
			{
				ID:          FlightChecks.Do.Cancel(),
				Statement:   "Do returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *FlightHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Do(ctx, string("key-a"), Loader(func(ctx context.Context, key string) (string, error) { return "", nil }))
						return err
					})
				},
				Defect: flightDoIgnoresCancel,
			},
			{
				ID:          FlightChecks.Do.CancelGenerated(),
				Statement:   "Do returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FlightHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					loadResults0 := c.Draw(prop.List(prop.Of[string]()), "load-0")
					loadResults1 := c.Draw(
						prop.List(
							prop.SampledFrom[error](nil, conform.ErrGenerated),
							prop.MinSize(len(loadResults0)),
							prop.MaxSize(len(loadResults0)),
						),
						"load-1",
					)
					drawnLoad := func() Loader {
						var loadAt atomic.Int64
						return func(context.Context, string) (r0 string, r1 error) {
							i := int(loadAt.Add(1) - 1)
							if i >= len(loadResults0) {
								r1 = conform.ErrGenerated
								return r0, r1
							}
							return loadResults0[i],
								loadResults1[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Do(ctx, drawnKey, drawnLoad())
						return err
					})
				},
				Defect: flightDoIgnoresCancel,
			},
			{
				ID:        FlightChecks.Do.FailureIsReported(),
				Statement: "Do returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *FlightHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Do(ctx, string("key-a"), Loader(func(ctx context.Context, key string) (string, error) { return "", nil }))
						return err
					})
				},
				Defect: flightDoSwallowsFailure,
			},
			{
				ID:        FlightChecks.Do.FailureIsReportedGenerated(),
				Statement: "Do returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Probe:     spec.ReportsAFailedDependency,
				Reason:    spec.SwallowedFailure,
				Law:       spec.FailedDependencyIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *FlightHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					loadResults0 := c.Draw(prop.List(prop.Of[string]()), "load-0")
					loadResults1 := c.Draw(
						prop.List(
							prop.SampledFrom[error](nil, conform.ErrGenerated),
							prop.MinSize(len(loadResults0)),
							prop.MaxSize(len(loadResults0)),
						),
						"load-1",
					)
					drawnLoad := func() Loader {
						var loadAt atomic.Int64
						return func(context.Context, string) (r0 string, r1 error) {
							i := int(loadAt.Add(1) - 1)
							if i >= len(loadResults0) {
								r1 = conform.ErrGenerated
								return r0, r1
							}
							return loadResults0[i],
								loadResults1[i]
						}
					}
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Do(ctx, drawnKey, drawnLoad())
						return err
					})
				},
				Defect: flightDoSwallowsFailure,
			},
			{
				ID:          FlightChecks.Do.ComputesOnce(),
				Statement:   "calls of Do with one key at once share one computation of their result",
				Claim:       spec.ComputesOnce,
				Probe:       spec.SharesOneComputation,
				Reason:      spec.RepeatedComputation,
				Law:         spec.CallsAtOnceShareOneComputation,
				Machineries: []spec.MachineryID{spec.ConcurrentCallers, spec.CallThatFailsOnDemand},
				Apply: func(h *FlightHandle) {
					probe.SharesOneComputation(h, func(ctx context.Context) (string, error) {
						return h.Subject().
							Do(ctx, string("key-a"), func(context.Context, string) (r0 string, r1 error) {
								_ = h.Injector().Dependency().Call("Do")
								return r0, r1
							})
					}, 2)
				},
				Defect: flightDoComputesPerCall,
				Around: synctest.Test,
			},
			{
				ID:        FlightChecks.Do.ComputesOnceGenerated(),
				Statement: "calls of Do with one key at once share one computation of their result",
				Claim:     spec.ComputesOnce,
				Probe:     spec.SharesOneComputation,
				Reason:    spec.RepeatedComputation,
				Law:       spec.CallsAtOnceShareOneComputation,
				Machineries: []spec.MachineryID{
					spec.ConcurrentCallers,
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
				},
				Property: func(h *FlightHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.SharesOneComputation(h, func(ctx context.Context) (string, error) {
						return h.Subject().
							Do(ctx, drawnKey, func(context.Context, string) (r0 string, r1 error) {
								_ = h.Injector().Dependency().Call("Do")
								return r0, r1
							})
					}, 2)
				},
				Defect: flightDoComputesPerCall,
				Around: synctest.Test,
			},
		},
		Calls: []conform.Call[Flight]{
			{
				Calls:    flightDoCalls,
				Over:     newFlightOver,
				Callable: FlightChecks.Do.Callable,
			},
		},
	}
}

// RunFlightSuite applies the checks of Flight, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunFlightSuite[T Flight](t *testing.T, newSubject func() T, opts ...FlightOption) {
	t.Helper()
	RunFlightSuiteFromHandle(t, func(*FlightHandle) T { return newSubject() }, opts...)
}

// RunFlightSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunFlightSuiteFromHandle[T Flight](
	t *testing.T,
	newSubject func(h *FlightHandle) T,
	opts ...FlightOption,
) {
	t.Helper()
	conform.Run(t, FlightSuite(), flightSubjects(newSubject), opts...)
}

// FuzzFlightSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzFlightSuite[T Flight](
	f *testing.F,
	callable conform.Callable[Flight],
	newSubject func(h *FlightHandle) T,
	opts ...FlightOption,
) {
	f.Helper()
	conform.Fuzz(f, FlightSuite(), callable, flightSubjects(newSubject), opts...)
}

// RunFlightBudgets applies the allocation ceilings of Flight and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunFlightBudgets[T Flight](
	t *testing.T,
	newSubject func(h *FlightHandle) T,
	opts ...FlightOption,
) {
	t.Helper()
	conform.RunBudgets(t, FlightSuite(), flightSubjects(newSubject), opts...)
}

// BenchFlightSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchFlightSuite[T Flight](
	b conform.B,
	callable conform.Callable[Flight],
	newSubject func(h *FlightHandle) T,
	opts ...FlightOption,
) {
	b.Helper()
	conform.Bench(b, FlightSuite(), callable, flightSubjects(newSubject), opts...)
}

// flightSubjects returns the factory of the subjects of Flight that calls
// newSubject with the handle of each check.
func flightSubjects[T Flight](newSubject func(h *FlightHandle) T) func(h *FlightHandle) Flight {
	return func(h *FlightHandle) Flight { return newSubject(h) }
}

// flightDoCalls generates the calls of Flight.Do with the inputs that each
// case draws.
var flightDoCalls = prop.Composite(func(c *prop.Case) func(context.Context, Flight) error {
	key := c.Draw(prop.Of[string](), "key")
	loadResults0 := c.Draw(prop.List(prop.Of[string]()), "load-0")
	loadResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(loadResults0)),
			prop.MaxSize(len(loadResults0)),
		),
		"load-1",
	)
	var loadAt atomic.Int64
	load := func(context.Context, string) (r0 string, r1 error) {
		i := int(loadAt.Add(1) - 1)
		if i >= len(loadResults0) {
			r1 = conform.ErrGenerated
			return r0, r1
		}
		return loadResults0[i],
			loadResults1[i]
	}
	return func(ctx context.Context, subject Flight) error {
		_, err := subject.Do(ctx, key, load)
		return err
	}
})

// flightOverEmbedded is Flight under the name of the field that flightOver
// embeds.
type flightOverEmbedded = Flight

// flightOver is a subject of Flight whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type flightOver struct {
	flightOverEmbedded
	spend func()
}

// Do calls spend and returns its zero results.
func (o flightOver) Do(context.Context, string, Loader) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newFlightOver returns the planted defect of a ceiling of Flight, whose
// callables call spend.
func newFlightOver(_ *FlightHandle, spend func()) Flight {
	return flightOver{spend: spend}
}

// flightDoPanicsDefect is the planted defect of a Flight whose Do panics.
type flightDoPanicsDefect struct {
	flightOverEmbedded
}

// Do is the override of the planted defect.
func (*flightDoPanicsDefect) Do(p0 context.Context, p1 string, p2 Loader) (r0 string, r1 error) {
	panic("the planted defect of a Flight whose Do panics")
}

// flightDoPanics returns the planted defect of a Flight whose Do panics
// over a fresh subject.
func flightDoPanics(h *FlightHandle) Flight {
	return &flightDoPanicsDefect{flightOverEmbedded: h.Fresh()}
}

// flightDoIgnoresCancelDefect is the planted defect of a Flight whose Do
// ignores the context.
type flightDoIgnoresCancelDefect struct {
	flightOverEmbedded
}

// Do is the override of the planted defect.
func (d *flightDoIgnoresCancelDefect) Do(
	p0 context.Context,
	p1 string,
	p2 Loader,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.flightOverEmbedded.Do(p0, p1, p2)
}

// flightDoIgnoresCancel returns the planted defect of a Flight whose Do
// ignores the context over a fresh subject.
func flightDoIgnoresCancel(h *FlightHandle) Flight {
	return &flightDoIgnoresCancelDefect{flightOverEmbedded: h.Fresh()}
}

// flightDoSwallowsFailureDefect is the planted defect of a Flight whose Do
// returns no error when its dependency fails.
type flightDoSwallowsFailureDefect struct {
	flightOverEmbedded
	h *FlightHandle
}

// Do is the override of the planted defect.
func (d *flightDoSwallowsFailureDefect) Do(
	p0 context.Context,
	p1 string,
	p2 Loader,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.flightOverEmbedded.Do(p0, p1, p2)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(string)
		r1 = nil
	}
	return r0, r1
}

// flightDoSwallowsFailure returns the planted defect of a Flight whose Do
// returns no error when its dependency fails over a fresh subject.
func flightDoSwallowsFailure(h *FlightHandle) Flight {
	return &flightDoSwallowsFailureDefect{flightOverEmbedded: h.Fresh(), h: h}
}

// flightDoComputesPerCallDefect is the planted defect of a Flight whose Do
// computes its result for each call.
type flightDoComputesPerCallDefect struct {
	flightOverEmbedded
	h *FlightHandle
}

// Do is the override of the planted defect.
func (d *flightDoComputesPerCallDefect) Do(
	p0 context.Context,
	p1 string,
	p2 Loader,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	return d.h.Fresh().Do(p0, p1, p2)
}

// flightDoComputesPerCall returns the planted defect of a Flight whose Do
// computes its result for each call over a fresh subject.
func flightDoComputesPerCall(h *FlightHandle) Flight {
	return &flightDoComputesPerCallDefect{flightOverEmbedded: h.Fresh(), h: h}
}

// The derived checks of Locker:
//
// Locker.Lock/smoke: Lock returns without a panic for derived inputs.
//
// Locker.Lock/smoke#generated: Lock returns without a panic for derived
// inputs.
//
// Locker.Lock/cancel: Lock returns an error for a cancelled context.
//
// Locker.Lock/cancel#generated: Lock returns an error for a cancelled
// context.
//
// Locker.Lock/lease: a second Lock for one key fails with the declared
// sentinel.
//
// Locker.Lock/lease#generated: a second Lock for one key fails with the
// declared sentinel.
//
// Locker.Lock/release-frees-the-lease: Lock grants a lease again after
// Unlock released it.
//
// Locker.Lock/release-frees-the-lease#generated: Lock grants a lease again
// after Unlock released it.
//
// Locker.Lock/one-holder-at-a-time: calls of Lock with one key at once
// grant the lease once.
//
// Locker.Lock/one-holder-at-a-time#generated: calls of Lock with one key
// at once grant the lease once.
//
// Locker.Unlock/smoke: Unlock returns without a panic for derived inputs.
//
// Locker.Unlock/smoke#generated: Unlock returns without a panic for
// derived inputs.
//
// Locker.Unlock/cancel: Unlock returns an error for a cancelled context.
//
// Locker.Unlock/cancel#generated: Unlock returns an error for a cancelled
// context.

// The derived checks of Locker that witness their laws at the shallow end:
//
// Locker.Lock/smoke
//
// Locker.Lock/lease
//
// Locker.Lock/lease#generated
//
// Locker.Unlock/smoke

// The unclaimed classifications of Locker:
//
// The shape:writer of Locker.Lock has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Locker.Unlock has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Locker writes.
type (
	LockerHandle = conform.Handle[Locker]
	LockerCheck  = conform.Check[Locker]
	LockerOption = conform.Option[Locker]
)

// LockerChecks contains the name of each callable of Locker, with the
// identities of its checks.
var LockerChecks = struct {
	Lock   lockerLock
	Unlock lockerUnlock
}{
	Lock:   lockerLock{"Locker.Lock"},
	Unlock: lockerUnlock{"Locker.Unlock"},
}

// lockerLock is the name of Locker.Lock, with the identities of its
// checks.
type lockerLock struct {
	conform.Callable[Locker]
}

// Smoke returns the identity of the check that Lock returns without a
// panic for derived inputs.
func (c lockerLock) Smoke() conform.ID[Locker] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Lock returns
// without a panic for derived inputs.
func (c lockerLock) SmokeGenerated() conform.ID[Locker] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Lock returns an error for
// a cancelled context.
func (c lockerLock) Cancel() conform.ID[Locker] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Lock returns an
// error for a cancelled context.
func (c lockerLock) CancelGenerated() conform.ID[Locker] {
	return c.Check("cancel#generated")
}

// Lease returns the identity of the check that a second Lock for one key
// fails with the declared sentinel.
func (c lockerLock) Lease() conform.ID[Locker] {
	return c.Check("lease")
}

// LeaseGenerated returns the identity of the check that a second Lock for
// one key fails with the declared sentinel.
func (c lockerLock) LeaseGenerated() conform.ID[Locker] {
	return c.Check("lease#generated")
}

// ReleaseFreesTheLease returns the identity of the check that Lock grants
// a lease again after Unlock released it.
func (c lockerLock) ReleaseFreesTheLease() conform.ID[Locker] {
	return c.Check("release-frees-the-lease")
}

// ReleaseFreesTheLeaseGenerated returns the identity of the check that
// Lock grants a lease again after Unlock released it.
func (c lockerLock) ReleaseFreesTheLeaseGenerated() conform.ID[Locker] {
	return c.Check("release-frees-the-lease#generated")
}

// OneHolderAtATime returns the identity of the check that calls of Lock
// with one key at once grant the lease once.
func (c lockerLock) OneHolderAtATime() conform.ID[Locker] {
	return c.Check("one-holder-at-a-time")
}

// OneHolderAtATimeGenerated returns the identity of the check that calls
// of Lock with one key at once grant the lease once.
func (c lockerLock) OneHolderAtATimeGenerated() conform.ID[Locker] {
	return c.Check("one-holder-at-a-time#generated")
}

// lockerUnlock is the name of Locker.Unlock, with the identities of its
// checks.
type lockerUnlock struct {
	conform.Callable[Locker]
}

// Smoke returns the identity of the check that Unlock returns without a
// panic for derived inputs.
func (c lockerUnlock) Smoke() conform.ID[Locker] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Unlock returns
// without a panic for derived inputs.
func (c lockerUnlock) SmokeGenerated() conform.ID[Locker] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Unlock returns an error
// for a cancelled context.
func (c lockerUnlock) Cancel() conform.ID[Locker] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Unlock returns an
// error for a cancelled context.
func (c lockerUnlock) CancelGenerated() conform.ID[Locker] {
	return c.Check("cancel#generated")
}

// LockerSuite returns the checks of Locker, and the calls and the budgets
// of its callables.
func LockerSuite() conform.Suite[Locker] {
	return conform.Suite[Locker]{
		Declaration: "Locker",
		Package:     "example.com/shop/locks",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Locker]{
			{
				ID:          LockerChecks.Lock.Smoke(),
				Statement:   "Lock returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, string("key-a"))
					})
				},
				Defect: lockerLockPanics,
			},
			{
				ID:          LockerChecks.Lock.SmokeGenerated(),
				Statement:   "Lock returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, drawnKey)
					})
				},
				Defect: lockerLockPanics,
			},
			{
				ID:          LockerChecks.Lock.Cancel(),
				Statement:   "Lock returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, string("key-a"))
					})
				},
				Defect: lockerLockIgnoresCancel,
			},
			{
				ID:          LockerChecks.Lock.CancelGenerated(),
				Statement:   "Lock returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, drawnKey)
					})
				},
				Defect: lockerLockIgnoresCancel,
			},
			{
				ID:          LockerChecks.Lock.Lease(),
				Statement:   "a second Lock for one key fails with the declared sentinel",
				Claim:       spec.Lease,
				Probe:       spec.RefusesTheSecondWrite,
				Reason:      spec.UnguardedWrite,
				Law:         spec.GrantedLeaseIsNotGrantedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockerHandle) {
					probe.RefusesTheSecondWrite(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, string("key-a"))
					}, ErrLocked)
				},
				Defect: lockerLockSucceedsAlways,
			},
			{
				ID:          LockerChecks.Lock.LeaseGenerated(),
				Statement:   "a second Lock for one key fails with the declared sentinel",
				Claim:       spec.Lease,
				Probe:       spec.RefusesTheSecondWrite,
				Reason:      spec.UnguardedWrite,
				Law:         spec.GrantedLeaseIsNotGrantedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.RefusesTheSecondWrite(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, drawnKey)
					}, ErrLocked)
				},
				Defect: lockerLockSucceedsAlways,
			},
			{
				ID:          LockerChecks.Lock.ReleaseFreesTheLease(),
				Statement:   "Lock grants a lease again after Unlock released it",
				Claim:       spec.ReleaseFreesTheLease,
				Probe:       spec.GrantsAfterRelease,
				Reason:      spec.KeptAfterRelease,
				Law:         spec.ReleasedLeaseIsGrantedAgain,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockerHandle) {
					probe.GrantsAfterRelease(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, string("key-a"))
					}, func(ctx context.Context) error {
						return h.Subject().Unlock(ctx, string("key-a"))
					})
				},
				Defect: lockerUnlockReleaseKeepsTheLease,
			},
			{
				ID:          LockerChecks.Lock.ReleaseFreesTheLeaseGenerated(),
				Statement:   "Lock grants a lease again after Unlock released it",
				Claim:       spec.ReleaseFreesTheLease,
				Probe:       spec.GrantsAfterRelease,
				Reason:      spec.KeptAfterRelease,
				Law:         spec.ReleasedLeaseIsGrantedAgain,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.GrantsAfterRelease(h, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, drawnKey)
					}, func(ctx context.Context) error {
						return h.Subject().Unlock(ctx, drawnKey)
					})
				},
				Defect: lockerUnlockReleaseKeepsTheLease,
			},
			{
				ID:          LockerChecks.Lock.OneHolderAtATime(),
				Statement:   "calls of Lock with one key at once grant the lease once",
				Claim:       spec.OneHolderAtATime,
				Probe:       spec.ExclusiveUnderConcurrency,
				Reason:      spec.UnguardedWrite,
				Law:         spec.GrantedLeaseIsNotGrantedTwice,
				Machineries: []spec.MachineryID{spec.ConcurrentCallers},
				Entails:     []conform.ID[Locker]{LockerChecks.Lock.Lease()},
				Property: func(h *LockerHandle, c *prop.Case) {
					probe.ExclusiveUnderConcurrency(h, c, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, string("key-a"))
					}, ErrLocked, 2)
				},
				Defect: lockerLockCheckThenGrant,
			},
			{
				ID:          LockerChecks.Lock.OneHolderAtATimeGenerated(),
				Statement:   "calls of Lock with one key at once grant the lease once",
				Claim:       spec.OneHolderAtATime,
				Probe:       spec.ExclusiveUnderConcurrency,
				Reason:      spec.UnguardedWrite,
				Law:         spec.GrantedLeaseIsNotGrantedTwice,
				Machineries: []spec.MachineryID{spec.ConcurrentCallers, spec.GeneratedInputs},
				Entails:     []conform.ID[Locker]{LockerChecks.Lock.LeaseGenerated()},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ExclusiveUnderConcurrency(h, c, func(ctx context.Context) error {
						return h.Subject().Lock(ctx, drawnKey)
					}, ErrLocked, 2)
				},
				Defect: lockerLockCheckThenGrant,
			},
			{
				ID:          LockerChecks.Unlock.Smoke(),
				Statement:   "Unlock returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Unlock(ctx, string("key-a"))
					})
				},
				Defect: lockerUnlockPanics,
			},
			{
				ID:          LockerChecks.Unlock.SmokeGenerated(),
				Statement:   "Unlock returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Unlock(ctx, drawnKey)
					})
				},
				Defect: lockerUnlockPanics,
			},
			{
				ID:          LockerChecks.Unlock.Cancel(),
				Statement:   "Unlock returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Unlock(ctx, string("key-a"))
					})
				},
				Defect: lockerUnlockIgnoresCancel,
			},
			{
				ID:          LockerChecks.Unlock.CancelGenerated(),
				Statement:   "Unlock returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Unlock(ctx, drawnKey)
					})
				},
				Defect: lockerUnlockIgnoresCancel,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Locker.Lock",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Locker.Unlock",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Locker]{
			{
				Calls:    lockerLockCalls,
				Over:     newLockerOver,
				Callable: LockerChecks.Lock.Callable,
			},
			{
				Calls:    lockerUnlockCalls,
				Over:     newLockerOver,
				Callable: LockerChecks.Unlock.Callable,
			},
		},
	}
}

// RunLockerSuite applies the checks of Locker, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLockerSuite[T Locker](t *testing.T, newSubject func() T, opts ...LockerOption) {
	t.Helper()
	RunLockerSuiteFromHandle(t, func(*LockerHandle) T { return newSubject() }, opts...)
}

// RunLockerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLockerSuiteFromHandle[T Locker](
	t *testing.T,
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	t.Helper()
	conform.Run(t, LockerSuite(), lockerSubjects(newSubject), opts...)
}

// FuzzLockerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLockerSuite[T Locker](
	f *testing.F,
	callable conform.Callable[Locker],
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	f.Helper()
	conform.Fuzz(f, LockerSuite(), callable, lockerSubjects(newSubject), opts...)
}

// RunLockerBudgets applies the allocation ceilings of Locker and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLockerBudgets[T Locker](
	t *testing.T,
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	t.Helper()
	conform.RunBudgets(t, LockerSuite(), lockerSubjects(newSubject), opts...)
}

// BenchLockerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLockerSuite[T Locker](
	b conform.B,
	callable conform.Callable[Locker],
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	b.Helper()
	conform.Bench(b, LockerSuite(), callable, lockerSubjects(newSubject), opts...)
}

// lockerSubjects returns the factory of the subjects of Locker that calls
// newSubject with the handle of each check.
func lockerSubjects[T Locker](newSubject func(h *LockerHandle) T) func(h *LockerHandle) Locker {
	return func(h *LockerHandle) Locker { return newSubject(h) }
}

// lockerLockCalls generates the calls of Locker.Lock with the inputs that
// each case draws.
var lockerLockCalls = prop.Composite(func(c *prop.Case) func(context.Context, Locker) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Locker) error {
		return subject.Lock(ctx, key)
	}
})

// lockerUnlockCalls generates the calls of Locker.Unlock with the inputs
// that each case draws.
var lockerUnlockCalls = prop.Composite(func(c *prop.Case) func(context.Context, Locker) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Locker) error {
		return subject.Unlock(ctx, key)
	}
})

// lockerOverEmbedded is Locker under the name of the field that lockerOver
// embeds.
type lockerOverEmbedded = Locker

// lockerOver is a subject of Locker whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type lockerOver struct {
	lockerOverEmbedded
	spend func()
}

// Lock calls spend and returns its zero results.
func (o lockerOver) Lock(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Unlock calls spend and returns its zero results.
func (o lockerOver) Unlock(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// newLockerOver returns the planted defect of a ceiling of Locker, whose
// callables call spend.
func newLockerOver(_ *LockerHandle, spend func()) Locker {
	return lockerOver{spend: spend}
}

// lockerLockPanicsDefect is the planted defect of a Locker whose Lock
// panics.
type lockerLockPanicsDefect struct {
	lockerOverEmbedded
}

// Lock is the override of the planted defect.
func (*lockerLockPanicsDefect) Lock(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Locker whose Lock panics")
}

// lockerLockPanics returns the planted defect of a Locker whose Lock
// panics over a fresh subject.
func lockerLockPanics(h *LockerHandle) Locker {
	return &lockerLockPanicsDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerLockIgnoresCancelDefect is the planted defect of a Locker whose
// Lock ignores the context.
type lockerLockIgnoresCancelDefect struct {
	lockerOverEmbedded
}

// Lock is the override of the planted defect.
func (d *lockerLockIgnoresCancelDefect) Lock(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.lockerOverEmbedded.Lock(p0, p1)
}

// lockerLockIgnoresCancel returns the planted defect of a Locker whose
// Lock ignores the context over a fresh subject.
func lockerLockIgnoresCancel(h *LockerHandle) Locker {
	return &lockerLockIgnoresCancelDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerLockSucceedsAlwaysDefect is the planted defect of a Locker whose
// Lock grants a lease that it holds.
type lockerLockSucceedsAlwaysDefect struct {
	lockerOverEmbedded
}

// Lock is the override of the planted defect.
func (*lockerLockSucceedsAlwaysDefect) Lock(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// lockerLockSucceedsAlways returns the planted defect of a Locker whose
// Lock grants a lease that it holds over a fresh subject.
func lockerLockSucceedsAlways(h *LockerHandle) Locker {
	return &lockerLockSucceedsAlwaysDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerUnlockReleaseKeepsTheLeaseDefect is the planted defect of a Locker
// whose Unlock keeps the lease.
type lockerUnlockReleaseKeepsTheLeaseDefect struct {
	lockerOverEmbedded
}

// Unlock is the override of the planted defect.
func (*lockerUnlockReleaseKeepsTheLeaseDefect) Unlock(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// lockerUnlockReleaseKeepsTheLease returns the planted defect of a Locker
// whose Unlock keeps the lease over a fresh subject.
func lockerUnlockReleaseKeepsTheLease(h *LockerHandle) Locker {
	return &lockerUnlockReleaseKeepsTheLeaseDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerLockCheckThenGrantDefect is the planted defect of a Locker whose
// Lock grants a lease to two calls at once.
type lockerLockCheckThenGrantDefect struct {
	lockerOverEmbedded
	h      *LockerHandle
	marked atomic.Int32
}

// Lock is the override of the planted defect.
func (d *lockerLockCheckThenGrantDefect) Lock(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	if d.marked.Load() > 1 {
		return r0
	}
	return d.lockerOverEmbedded.Lock(p0, p1)
}

// lockerLockCheckThenGrant returns the planted defect of a Locker whose
// Lock grants a lease to two calls at once over a fresh subject.
func lockerLockCheckThenGrant(h *LockerHandle) Locker {
	return &lockerLockCheckThenGrantDefect{lockerOverEmbedded: h.Fresh(), h: h}
}

// lockerUnlockPanicsDefect is the planted defect of a Locker whose Unlock
// panics.
type lockerUnlockPanicsDefect struct {
	lockerOverEmbedded
}

// Unlock is the override of the planted defect.
func (*lockerUnlockPanicsDefect) Unlock(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Locker whose Unlock panics")
}

// lockerUnlockPanics returns the planted defect of a Locker whose Unlock
// panics over a fresh subject.
func lockerUnlockPanics(h *LockerHandle) Locker {
	return &lockerUnlockPanicsDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerUnlockIgnoresCancelDefect is the planted defect of a Locker whose
// Unlock ignores the context.
type lockerUnlockIgnoresCancelDefect struct {
	lockerOverEmbedded
}

// Unlock is the override of the planted defect.
func (d *lockerUnlockIgnoresCancelDefect) Unlock(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.lockerOverEmbedded.Unlock(p0, p1)
}

// lockerUnlockIgnoresCancel returns the planted defect of a Locker whose
// Unlock ignores the context over a fresh subject.
func lockerUnlockIgnoresCancel(h *LockerHandle) Locker {
	return &lockerUnlockIgnoresCancelDefect{lockerOverEmbedded: h.Fresh()}
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
