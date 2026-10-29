package tagged

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Hasher:
//
// Hasher.HashTagged/smoke: HashTagged returns without a panic for derived
// inputs.
//
// Hasher.HashTagged/deterministic: HashTagged returns the same value for
// the same inputs.
//
// Hasher.CombineTagged/smoke: CombineTagged returns without a panic for
// derived inputs.
//
// Hasher.CombineTagged/deterministic: CombineTagged returns the same value
// for the same inputs.
//
// Hasher.HashRoles/smoke: HashRoles returns without a panic for derived
// inputs.
//
// Hasher.HashRoles/deterministic: HashRoles returns the same value for the
// same inputs.
//
// Hasher.Truncate/smoke: Truncate returns without a panic for derived
// inputs.
//
// Hasher.Truncate/deterministic: Truncate returns the same value for the
// same inputs.
//
// Hasher.Seal/smoke: Seal returns without a panic for derived inputs.
//
// Hasher.Seal/cancel: Seal returns an error for a cancelled context.

// The derived checks of Hasher that witness their laws at the shallow end:
//
// Hasher.HashTagged/smoke
//
// Hasher.HashTagged/deterministic
//
// Hasher.CombineTagged/smoke
//
// Hasher.CombineTagged/deterministic
//
// Hasher.HashRoles/smoke
//
// Hasher.HashRoles/deterministic
//
// Hasher.Truncate/smoke
//
// Hasher.Truncate/deterministic
//
// Hasher.Seal/smoke

// The unclaimed classifications of Hasher:
//
// The shape:writer of Hasher.Seal has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Hasher writes.
type (
	HasherHandle = conform.Handle[Hasher]
	HasherCheck  = conform.Check[Hasher]
	HasherOption = conform.Option[Hasher]
)

// HasherChecks contains the name of each callable of Hasher, with the
// identities of its checks.
var HasherChecks = struct {
	HashTagged    hasherHashTagged
	CombineTagged hasherCombineTagged
	HashRoles     hasherHashRoles
	Truncate      hasherTruncate
	Seal          hasherSeal
}{
	HashTagged:    hasherHashTagged{"Hasher.HashTagged"},
	CombineTagged: hasherCombineTagged{"Hasher.CombineTagged"},
	HashRoles:     hasherHashRoles{"Hasher.HashRoles"},
	Truncate:      hasherTruncate{"Hasher.Truncate"},
	Seal:          hasherSeal{"Hasher.Seal"},
}

// hasherHashTagged is the name of Hasher.HashTagged, with the identities
// of its checks.
type hasherHashTagged struct {
	conform.Callable[Hasher]
}

// Smoke returns the identity of the check that HashTagged returns without
// a panic for derived inputs.
func (c hasherHashTagged) Smoke() conform.ID[Hasher] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that HashTagged returns
// the same value for the same inputs.
func (c hasherHashTagged) Deterministic() conform.ID[Hasher] {
	return c.Check("deterministic")
}

// hasherCombineTagged is the name of Hasher.CombineTagged, with the
// identities of its checks.
type hasherCombineTagged struct {
	conform.Callable[Hasher]
}

// Smoke returns the identity of the check that CombineTagged returns
// without a panic for derived inputs.
func (c hasherCombineTagged) Smoke() conform.ID[Hasher] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that CombineTagged
// returns the same value for the same inputs.
func (c hasherCombineTagged) Deterministic() conform.ID[Hasher] {
	return c.Check("deterministic")
}

// hasherHashRoles is the name of Hasher.HashRoles, with the identities of
// its checks.
type hasherHashRoles struct {
	conform.Callable[Hasher]
}

// Smoke returns the identity of the check that HashRoles returns without a
// panic for derived inputs.
func (c hasherHashRoles) Smoke() conform.ID[Hasher] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that HashRoles returns
// the same value for the same inputs.
func (c hasherHashRoles) Deterministic() conform.ID[Hasher] {
	return c.Check("deterministic")
}

// hasherTruncate is the name of Hasher.Truncate, with the identities of
// its checks.
type hasherTruncate struct {
	conform.Callable[Hasher]
}

// Smoke returns the identity of the check that Truncate returns without a
// panic for derived inputs.
func (c hasherTruncate) Smoke() conform.ID[Hasher] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Truncate returns
// the same value for the same inputs.
func (c hasherTruncate) Deterministic() conform.ID[Hasher] {
	return c.Check("deterministic")
}

// hasherSeal is the name of Hasher.Seal, with the identities of its
// checks.
type hasherSeal struct {
	conform.Callable[Hasher]
}

// Smoke returns the identity of the check that Seal returns without a
// panic for derived inputs.
func (c hasherSeal) Smoke() conform.ID[Hasher] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Seal returns an error for
// a cancelled context.
func (c hasherSeal) Cancel() conform.ID[Hasher] {
	return c.Check("cancel")
}

// HasherSuite returns the checks of Hasher, and the calls and the budgets
// of its callables.
func HasherSuite() conform.Suite[Hasher] {
	return conform.Suite[Hasher]{
		Declaration: "Hasher",
		Checks: []conform.Check[Hasher]{
			{
				ID:          HasherChecks.HashTagged.Smoke(),
				Statement:   "HashTagged returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HasherHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().HashTagged(Role(1), []byte{1})
						return nil
					})
				},
				Defect: hasherHashTaggedPanics,
			},
			{
				ID:          HasherChecks.HashTagged.Deterministic(),
				Statement:   "HashTagged returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Hasher]{HasherChecks.HashTagged.Smoke()},
				Apply: func(h *HasherHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) ([32]byte, error) {
						r0 := h.Subject().HashTagged(Role(1), []byte{1})
						return r0, nil
					})
				},
				Defect: hasherHashTaggedDrifts,
			},
			{
				ID:          HasherChecks.CombineTagged.Smoke(),
				Statement:   "CombineTagged returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HasherHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().CombineTagged(Role(129), [32]byte{1}, [32]byte{1})
						return nil
					})
				},
				Defect: hasherCombineTaggedPanics,
			},
			{
				ID:          HasherChecks.CombineTagged.Deterministic(),
				Statement:   "CombineTagged returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Hasher]{HasherChecks.CombineTagged.Smoke()},
				Apply: func(h *HasherHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) ([32]byte, error) {
						r0 := h.Subject().CombineTagged(Role(129), [32]byte{1}, [32]byte{1})
						return r0, nil
					})
				},
				Defect: hasherCombineTaggedDrifts,
			},
			{
				ID:          HasherChecks.HashRoles.Smoke(),
				Statement:   "HashRoles returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HasherHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().HashRoles(Role(2))
						return nil
					})
				},
				Defect: hasherHashRolesPanics,
			},
			{
				ID:          HasherChecks.HashRoles.Deterministic(),
				Statement:   "HashRoles returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Hasher]{HasherChecks.HashRoles.Smoke()},
				Apply: func(h *HasherHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) ([32]byte, error) {
						r0 := h.Subject().HashRoles(Role(2))
						return r0, nil
					})
				},
				Defect: hasherHashRolesDrifts,
			},
			{
				ID:          HasherChecks.Truncate.Smoke(),
				Statement:   "Truncate returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HasherHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Truncate([32]byte{1}, Size(16))
						return nil
					})
				},
				Defect: hasherTruncatePanics,
			},
			{
				ID:          HasherChecks.Truncate.Deterministic(),
				Statement:   "Truncate returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Hasher]{HasherChecks.Truncate.Smoke()},
				Apply: func(h *HasherHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) ([]byte, error) {
						r0 := h.Subject().Truncate([32]byte{1}, Size(16))
						return r0, nil
					})
				},
				Defect: hasherTruncateDrifts,
			},
			{
				ID:          HasherChecks.Seal.Smoke(),
				Statement:   "Seal returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HasherHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Seal(ctx, [32]byte{1})
					})
				},
				Defect: hasherSealPanics,
			},
			{
				ID:          HasherChecks.Seal.Cancel(),
				Statement:   "Seal returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *HasherHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Seal(ctx, [32]byte{1})
					})
				},
				Defect: hasherSealIgnoresCancel,
			},
		},
		Calls: []conform.Call[Hasher]{
			{
				Calls:    hasherHashTaggedCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.HashTagged.Callable,
			},
			{
				Calls:    hasherCombineTaggedCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.CombineTagged.Callable,
			},
			{
				Calls:    hasherHashRolesCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.HashRoles.Callable,
			},
			{
				Calls:    hasherTruncateCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.Truncate.Callable,
			},
			{
				Calls:    hasherSealCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.Seal.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(hasherSizeValues),
		},
	}
}

// RunHasherSuite applies the checks of Hasher, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunHasherSuite[T Hasher](t *testing.T, newSubject func() T, opts ...HasherOption) {
	t.Helper()
	RunHasherSuiteFromHandle(t, func(*HasherHandle) T { return newSubject() }, opts...)
}

// RunHasherSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunHasherSuiteFromHandle[T Hasher](
	t *testing.T,
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	t.Helper()
	conform.Run(t, HasherSuite(), hasherSubjects(newSubject), opts...)
}

// FuzzHasherSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzHasherSuite[T Hasher](
	f *testing.F,
	callable conform.Callable[Hasher],
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	f.Helper()
	conform.Fuzz(f, HasherSuite(), callable, hasherSubjects(newSubject), opts...)
}

// RunHasherBudgets applies the allocation ceilings of Hasher and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunHasherBudgets[T Hasher](
	t *testing.T,
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	t.Helper()
	conform.RunBudgets(t, HasherSuite(), hasherSubjects(newSubject), opts...)
}

// BenchHasherSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchHasherSuite[T Hasher](
	b conform.B,
	callable conform.Callable[Hasher],
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	b.Helper()
	conform.Bench(b, HasherSuite(), callable, hasherSubjects(newSubject), opts...)
}

// hasherSubjects returns the factory of the subjects of Hasher that calls
// newSubject with the handle of each check.
func hasherSubjects[T Hasher](newSubject func(h *HasherHandle) T) func(h *HasherHandle) Hasher {
	return func(h *HasherHandle) Hasher { return newSubject(h) }
}

// hasherHashTaggedCalls generates the calls of Hasher.HashTagged with the
// inputs that each case draws.
var hasherHashTaggedCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hasher) error {
	r := Role(1)
	data := c.Draw(prop.Of[[]byte](), "data")
	return func(_ context.Context, subject Hasher) error {
		_ = subject.HashTagged(r, data)
		return nil
	}
})

// hasherCombineTaggedCalls generates the calls of Hasher.CombineTagged
// with the inputs that each case draws.
var hasherCombineTaggedCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, Hasher) error {
		r := Role(129)
		left := c.Draw(prop.Of[[32]byte](), "left")
		right := c.Draw(prop.Of[[32]byte](), "right")
		return func(_ context.Context, subject Hasher) error {
			_ = subject.CombineTagged(r, left, right)
			return nil
		}
	},
)

// hasherHashRolesCalls generates the calls of Hasher.HashRoles with the
// inputs that each case draws.
var hasherHashRolesCalls = prop.Composite(func(*prop.Case) func(context.Context, Hasher) error {
	roles := []Role{Role(2)}
	return func(_ context.Context, subject Hasher) error {
		_ = subject.HashRoles(roles...)
		return nil
	}
})

// hasherTruncateCalls generates the calls of Hasher.Truncate with the
// inputs that each case draws.
var hasherTruncateCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hasher) error {
	d := c.Draw(prop.Of[[32]byte](), "d")
	n := c.Draw(prop.Of[Size](), "n")
	return func(_ context.Context, subject Hasher) error {
		_ = subject.Truncate(d, n)
		return nil
	}
})

// hasherSealCalls generates the calls of Hasher.Seal with the inputs that
// each case draws.
var hasherSealCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hasher) error {
	d := c.Draw(prop.Of[[32]byte](), "d")
	return func(ctx context.Context, subject Hasher) error {
		return subject.Seal(ctx, d)
	}
})

// hasherSizeValues generates the values of Size that keep or change its
// starting value.
var hasherSizeValues = prop.Composite(func(c *prop.Case) Size {
	changed := c.Draw(prop.Of[Size](), "size")
	if c.Draw(prop.Boolean(), "size.changed") {
		return changed
	}
	return 16
})

// hasherOverEmbedded is Hasher under the name of the field that hasherOver
// embeds.
type hasherOverEmbedded = Hasher

// hasherOver is a subject of Hasher whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type hasherOver struct {
	hasherOverEmbedded
	spend func()
}

// HashTagged calls spend and returns its zero results.
func (o hasherOver) HashTagged(Role, []byte) (r0 [32]byte) {
	o.spend()
	return r0
}

// CombineTagged calls spend and returns its zero results.
func (o hasherOver) CombineTagged(Role, [32]byte, [32]byte) (r0 [32]byte) {
	o.spend()
	return r0
}

// HashRoles calls spend and returns its zero results.
func (o hasherOver) HashRoles(...Role) (r0 [32]byte) {
	o.spend()
	return r0
}

// Truncate calls spend and returns its zero results.
func (o hasherOver) Truncate([32]byte, Size) (r0 []byte) {
	o.spend()
	return r0
}

// Seal calls spend and returns its zero results.
func (o hasherOver) Seal(context.Context, [32]byte) (r0 error) {
	o.spend()
	return r0
}

// newHasherOver returns the planted defect of a ceiling of Hasher, whose
// callables call spend.
func newHasherOver(_ *HasherHandle, spend func()) Hasher {
	return hasherOver{spend: spend}
}

// hasherHashTaggedPanicsDefect is the planted defect of a Hasher whose
// HashTagged panics.
type hasherHashTaggedPanicsDefect struct {
	hasherOverEmbedded
}

// HashTagged is the override of the planted defect.
func (*hasherHashTaggedPanicsDefect) HashTagged(p0 Role, p1 []byte) (r0 [32]byte) {
	panic("the planted defect of a Hasher whose HashTagged panics")
}

// hasherHashTaggedPanics returns the planted defect of a Hasher whose
// HashTagged panics over a fresh subject.
func hasherHashTaggedPanics(h *HasherHandle) Hasher {
	return &hasherHashTaggedPanicsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherHashTaggedDriftsDefect is the planted defect of a Hasher whose
// HashTagged returns another value on the repeat.
type hasherHashTaggedDriftsDefect struct {
	hasherOverEmbedded
	calls int
}

// HashTagged is the override of the planted defect.
func (d *hasherHashTaggedDriftsDefect) HashTagged(p0 Role, p1 []byte) (r0 [32]byte) {
	d.calls++
	if d.calls > 1 {
		r0 = [32]byte{2}
	}
	return r0
}

// hasherHashTaggedDrifts returns the planted defect of a Hasher whose
// HashTagged returns another value on the repeat over a fresh subject.
func hasherHashTaggedDrifts(h *HasherHandle) Hasher {
	return &hasherHashTaggedDriftsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherCombineTaggedPanicsDefect is the planted defect of a Hasher whose
// CombineTagged panics.
type hasherCombineTaggedPanicsDefect struct {
	hasherOverEmbedded
}

// CombineTagged is the override of the planted defect.
func (*hasherCombineTaggedPanicsDefect) CombineTagged(p0 Role, p1, p2 [32]byte) (r0 [32]byte) {
	panic("the planted defect of a Hasher whose CombineTagged panics")
}

// hasherCombineTaggedPanics returns the planted defect of a Hasher whose
// CombineTagged panics over a fresh subject.
func hasherCombineTaggedPanics(h *HasherHandle) Hasher {
	return &hasherCombineTaggedPanicsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherCombineTaggedDriftsDefect is the planted defect of a Hasher whose
// CombineTagged returns another value on the repeat.
type hasherCombineTaggedDriftsDefect struct {
	hasherOverEmbedded
	calls int
}

// CombineTagged is the override of the planted defect.
func (d *hasherCombineTaggedDriftsDefect) CombineTagged(p0 Role, p1, p2 [32]byte) (r0 [32]byte) {
	d.calls++
	if d.calls > 1 {
		r0 = [32]byte{2}
	}
	return r0
}

// hasherCombineTaggedDrifts returns the planted defect of a Hasher whose
// CombineTagged returns another value on the repeat over a fresh subject.
func hasherCombineTaggedDrifts(h *HasherHandle) Hasher {
	return &hasherCombineTaggedDriftsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherHashRolesPanicsDefect is the planted defect of a Hasher whose
// HashRoles panics.
type hasherHashRolesPanicsDefect struct {
	hasherOverEmbedded
}

// HashRoles is the override of the planted defect.
func (*hasherHashRolesPanicsDefect) HashRoles(p0 ...Role) (r0 [32]byte) {
	panic("the planted defect of a Hasher whose HashRoles panics")
}

// hasherHashRolesPanics returns the planted defect of a Hasher whose
// HashRoles panics over a fresh subject.
func hasherHashRolesPanics(h *HasherHandle) Hasher {
	return &hasherHashRolesPanicsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherHashRolesDriftsDefect is the planted defect of a Hasher whose
// HashRoles returns another value on the repeat.
type hasherHashRolesDriftsDefect struct {
	hasherOverEmbedded
	calls int
}

// HashRoles is the override of the planted defect.
func (d *hasherHashRolesDriftsDefect) HashRoles(p0 ...Role) (r0 [32]byte) {
	d.calls++
	if d.calls > 1 {
		r0 = [32]byte{2}
	}
	return r0
}

// hasherHashRolesDrifts returns the planted defect of a Hasher whose
// HashRoles returns another value on the repeat over a fresh subject.
func hasherHashRolesDrifts(h *HasherHandle) Hasher {
	return &hasherHashRolesDriftsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherTruncatePanicsDefect is the planted defect of a Hasher whose
// Truncate panics.
type hasherTruncatePanicsDefect struct {
	hasherOverEmbedded
}

// Truncate is the override of the planted defect.
func (*hasherTruncatePanicsDefect) Truncate(p0 [32]byte, p1 Size) (r0 []byte) {
	panic("the planted defect of a Hasher whose Truncate panics")
}

// hasherTruncatePanics returns the planted defect of a Hasher whose
// Truncate panics over a fresh subject.
func hasherTruncatePanics(h *HasherHandle) Hasher {
	return &hasherTruncatePanicsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherTruncateDriftsDefect is the planted defect of a Hasher whose
// Truncate returns another value on the repeat.
type hasherTruncateDriftsDefect struct {
	hasherOverEmbedded
	calls int
}

// Truncate is the override of the planted defect.
func (d *hasherTruncateDriftsDefect) Truncate(p0 [32]byte, p1 Size) (r0 []byte) {
	d.calls++
	if d.calls > 1 {
		r0 = []byte{2}
	}
	return r0
}

// hasherTruncateDrifts returns the planted defect of a Hasher whose
// Truncate returns another value on the repeat over a fresh subject.
func hasherTruncateDrifts(h *HasherHandle) Hasher {
	return &hasherTruncateDriftsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherSealPanicsDefect is the planted defect of a Hasher whose Seal
// panics.
type hasherSealPanicsDefect struct {
	hasherOverEmbedded
}

// Seal is the override of the planted defect.
func (*hasherSealPanicsDefect) Seal(p0 context.Context, p1 [32]byte) (r0 error) {
	panic("the planted defect of a Hasher whose Seal panics")
}

// hasherSealPanics returns the planted defect of a Hasher whose Seal
// panics over a fresh subject.
func hasherSealPanics(h *HasherHandle) Hasher {
	return &hasherSealPanicsDefect{hasherOverEmbedded: h.Fresh()}
}

// hasherSealIgnoresCancelDefect is the planted defect of a Hasher whose
// Seal ignores the context.
type hasherSealIgnoresCancelDefect struct {
	hasherOverEmbedded
}

// Seal is the override of the planted defect.
func (d *hasherSealIgnoresCancelDefect) Seal(p0 context.Context, p1 [32]byte) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.hasherOverEmbedded.Seal(p0, p1)
}

// hasherSealIgnoresCancel returns the planted defect of a Hasher whose
// Seal ignores the context over a fresh subject.
func hasherSealIgnoresCancel(h *HasherHandle) Hasher {
	return &hasherSealIgnoresCancelDefect{hasherOverEmbedded: h.Fresh()}
}

// The derived checks of Digest:
//
// Digest/smoke: Digest returns without a panic for derived inputs.
//
// Digest/deterministic: Digest returns the same value for the same inputs.

// The derived checks of Digest that witness their laws at the shallow end:
//
// Digest/smoke
//
// Digest/deterministic

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Digest writes.
type (
	DigestHandle = conform.Handle[DigestFunc]
	DigestCheck  = conform.Check[DigestFunc]
	DigestOption = conform.Option[DigestFunc]
)

// DigestFunc is the type of the function Digest. Each subject of the suite
// of Digest has it.
type DigestFunc func(r Role, data []byte) [32]byte

// DigestChecks contains the name of each callable of Digest, with the
// identities of its checks.
var DigestChecks = digestCallable{"Digest"}

// digestCallable is the name of Digest, with the identities of its checks.
type digestCallable struct {
	conform.Callable[DigestFunc]
}

// Smoke returns the identity of the check that Digest returns without a
// panic for derived inputs.
func (c digestCallable) Smoke() conform.ID[DigestFunc] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Digest returns the
// same value for the same inputs.
func (c digestCallable) Deterministic() conform.ID[DigestFunc] {
	return c.Check("deterministic")
}

// DigestSuite returns the checks of Digest, and the calls and the budgets
// of its callables.
func DigestSuite() conform.Suite[DigestFunc] {
	return conform.Suite[DigestFunc]{
		Declaration: "Digest",
		Checks: []conform.Check[DigestFunc]{
			{
				ID:          DigestChecks.Smoke(),
				Statement:   "Digest returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DigestHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject()(Role(3), []byte{1})
						return nil
					})
				},
				Defect: digestPanics,
			},
			{
				ID:          DigestChecks.Deterministic(),
				Statement:   "Digest returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[DigestFunc]{DigestChecks.Smoke()},
				Apply: func(h *DigestHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) ([32]byte, error) {
						r0 := h.Subject()(Role(3), []byte{1})
						return r0, nil
					})
				},
				Defect: digestDrifts,
			},
		},
		Calls: []conform.Call[DigestFunc]{
			{
				Calls:    digestCalls,
				Over:     newDigestOver,
				Callable: DigestChecks.Callable,
			},
		},
	}
}

// RunDigestSuite applies the checks of Digest, and the checks that opts
// add, to subject.
func RunDigestSuite(t *testing.T, subject DigestFunc, opts ...DigestOption) {
	t.Helper()
	RunDigestSuiteFromHandle(t, digestSubjects(subject), opts...)
}

// RunDigestSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDigestSuiteFromHandle(
	t *testing.T,
	newSubject func(h *DigestHandle) DigestFunc,
	opts ...DigestOption,
) {
	t.Helper()
	conform.Run(t, DigestSuite(), newSubject, opts...)
}

// FuzzDigestSuite registers on f the fuzz target of the property checks of
// Digest on subject.
func FuzzDigestSuite(f *testing.F, subject DigestFunc, opts ...DigestOption) {
	f.Helper()
	conform.Fuzz(f, DigestSuite(), DigestChecks.Callable, digestSubjects(subject), opts...)
}

// RunDigestBudgets applies the allocation ceiling of Digest and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunDigestBudgets(t *testing.T, subject DigestFunc, opts ...DigestOption) {
	t.Helper()
	conform.RunBudgets(t, DigestSuite(), digestSubjects(subject), opts...)
}

// BenchDigestSuite measures subject under the ceilings of the budget of
// Digest.
func BenchDigestSuite(b conform.B, subject DigestFunc, opts ...DigestOption) {
	b.Helper()
	conform.Bench(b, DigestSuite(), DigestChecks.Callable, digestSubjects(subject), opts...)
}

// digestSubjects returns the factory of the subjects of Digest that
// returns subject for each handle.
func digestSubjects(subject DigestFunc) func(*DigestHandle) DigestFunc {
	return func(*DigestHandle) DigestFunc { return subject }
}

// digestCalls generates the calls of Digest with the inputs that each case
// draws.
var digestCalls = prop.Composite(func(c *prop.Case) func(context.Context, DigestFunc) error {
	r := Role(3)
	data := c.Draw(prop.Of[[]byte](), "data")
	return func(_ context.Context, subject DigestFunc) error {
		_ = subject(r, data)
		return nil
	}
})

// newDigestOver returns the planted defect of a ceiling of Digest, which
// calls spend and returns its zero results.
func newDigestOver(_ *DigestHandle, spend func()) DigestFunc {
	return func(Role, []byte) (r0 [32]byte) {
		spend()
		return r0
	}
}

// digestPanics returns the planted defect of a Digest whose Digest panics
// over the fresh function.
func digestPanics(h *DigestHandle) DigestFunc {
	return func(p0 Role, p1 []byte) (r0 [32]byte) {
		panic("the planted defect of a Digest whose Digest panics")
	}
}

// digestDrifts returns the planted defect of a Digest whose Digest returns
// another value on the repeat over the fresh function.
func digestDrifts(h *DigestHandle) DigestFunc {
	calls := 0
	return func(p0 Role, p1 []byte) (r0 [32]byte) {
		calls++
		if calls > 1 {
			r0 = [32]byte{2}
		}
		return r0
	}
}
