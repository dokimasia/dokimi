package leadertest

import (
	"context"
	"testing"

	"example.com/shop/leader"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Primary writes.
type (
	PrimaryHandle = conform.Handle[leader.Primary]
	PrimaryCheck  = conform.Check[leader.Primary]
	PrimaryOption = conform.Option[leader.Primary]
)

// PrimaryChecks contains the name of each callable of Primary, with the
// identities of its checks.
var PrimaryChecks = struct {
	Commit primaryCommit
}{
	Commit: primaryCommit{"Primary.Commit"},
}

// primaryCommit is the name of Primary.Commit, with the identities of its
// checks.
type primaryCommit struct {
	conform.Callable[leader.Primary]
}

// Smoke returns the identity of the check that Commit returns without a
// panic for derived inputs.
func (c primaryCommit) Smoke() conform.ID[leader.Primary] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Commit returns the
// same value for the same inputs.
func (c primaryCommit) Deterministic() conform.ID[leader.Primary] {
	return c.Check("deterministic")
}

// AllocationBudget returns the identity of the check
// Primary.Commit/allocation-budget.
func (c primaryCommit) AllocationBudget() conform.ID[leader.Primary] {
	return c.Check("allocation-budget")
}

// PrimarySuite returns the checks of Primary, and the calls and the
// budgets of its callables.
func PrimarySuite() conform.Suite[leader.Primary] {
	return conform.Suite[leader.Primary]{
		Declaration: "Primary",
		Package:     "example.com/shop/leader",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[leader.Primary]{
			{
				ID:          PrimaryChecks.Commit.Smoke(),
				Statement:   "Commit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PrimaryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Commit()
						return nil
					})
				},
				Defect: primaryCommitPanics,
			},
			{
				ID:          PrimaryChecks.Commit.Deterministic(),
				Statement:   "Commit returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[leader.Primary]{PrimaryChecks.Commit.Smoke()},
				Apply: func(h *PrimaryHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Commit()
						return r0, nil
					})
				},
				Defect: primaryCommitDrifts,
			},
		},
		Calls: []conform.Call[leader.Primary]{
			{
				Calls:    primaryCommitCalls,
				Over:     newPrimaryOver,
				Callable: PrimaryChecks.Commit.Callable,
				Budget: conform.Budget{
					Allocs: new(uint64(0)),
				},
			},
		},
	}
}

// RunPrimarySuite applies the checks of Primary, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunPrimarySuite[T leader.Primary](t *testing.T, newSubject func() T, opts ...PrimaryOption) {
	t.Helper()
	RunPrimarySuiteFromHandle(t, func(*PrimaryHandle) T { return newSubject() }, opts...)
}

// RunPrimarySuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPrimarySuiteFromHandle[T leader.Primary](
	t *testing.T,
	newSubject func(h *PrimaryHandle) T,
	opts ...PrimaryOption,
) {
	t.Helper()
	conform.Run(t, PrimarySuite(), primarySubjects(newSubject), opts...)
}

// FuzzPrimarySuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPrimarySuite[T leader.Primary](
	f *testing.F,
	callable conform.Callable[leader.Primary],
	newSubject func(h *PrimaryHandle) T,
	opts ...PrimaryOption,
) {
	f.Helper()
	conform.Fuzz(f, PrimarySuite(), callable, primarySubjects(newSubject), opts...)
}

// RunPrimaryBudgets applies the allocation ceilings of Primary and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunPrimaryBudgets[T leader.Primary](
	t *testing.T,
	newSubject func(h *PrimaryHandle) T,
	opts ...PrimaryOption,
) {
	t.Helper()
	conform.RunBudgets(t, PrimarySuite(), primarySubjects(newSubject), opts...)
}

// BenchPrimarySuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchPrimarySuite[T leader.Primary](
	b conform.B,
	callable conform.Callable[leader.Primary],
	newSubject func(h *PrimaryHandle) T,
	opts ...PrimaryOption,
) {
	b.Helper()
	conform.Bench(b, PrimarySuite(), callable, primarySubjects(newSubject), opts...)
}

// primarySubjects returns the factory of the subjects of Primary that
// calls newSubject with the handle of each check.
func primarySubjects[T leader.Primary](
	newSubject func(h *PrimaryHandle) T,
) func(h *PrimaryHandle) leader.Primary {
	return func(h *PrimaryHandle) leader.Primary { return newSubject(h) }
}

// primaryCommitCalls generates the calls of Primary.Commit with the inputs
// that each case draws.
var primaryCommitCalls = prop.Composite(
	func(*prop.Case) func(context.Context, leader.Primary) error {
		return func(_ context.Context, subject leader.Primary) error {
			_ = subject.Commit()
			return nil
		}
	},
)

// primaryOverEmbedded is leader.Primary under the name of the field that
// primaryOver embeds.
type primaryOverEmbedded = leader.Primary

// primaryOver is a subject of Primary whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type primaryOver struct {
	primaryOverEmbedded
	spend func()
}

// Commit calls spend and returns its zero results.
func (o primaryOver) Commit() (r0 int) {
	o.spend()
	return r0
}

// newPrimaryOver returns the planted defect of a ceiling of Primary, whose
// callables call spend.
func newPrimaryOver(_ *PrimaryHandle, spend func()) leader.Primary {
	return primaryOver{spend: spend}
}

// primaryCommitPanicsDefect is the planted defect of a Primary whose
// Commit panics.
type primaryCommitPanicsDefect struct {
	primaryOverEmbedded
}

// Commit is the override of the planted defect.
func (*primaryCommitPanicsDefect) Commit() (r0 int) {
	panic("the planted defect of a Primary whose Commit panics")
}

// primaryCommitPanics returns the planted defect of a Primary whose Commit
// panics over a fresh subject.
func primaryCommitPanics(h *PrimaryHandle) leader.Primary {
	return &primaryCommitPanicsDefect{primaryOverEmbedded: h.Fresh()}
}

// primaryCommitDriftsDefect is the planted defect of a Primary whose
// Commit returns another value on the repeat.
type primaryCommitDriftsDefect struct {
	primaryOverEmbedded
	calls int
}

// Commit is the override of the planted defect.
func (d *primaryCommitDriftsDefect) Commit() (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// primaryCommitDrifts returns the planted defect of a Primary whose Commit
// returns another value on the repeat over a fresh subject.
func primaryCommitDrifts(h *PrimaryHandle) leader.Primary {
	return &primaryCommitDriftsDefect{primaryOverEmbedded: h.Fresh()}
}
