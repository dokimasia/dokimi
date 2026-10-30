package ledger

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Ledger:
//
// Ledger.Credit/smoke: Credit returns without a panic for derived inputs.
//
// Ledger.Balance/smoke: Balance returns without a panic for derived
// inputs.
//
// Ledger.Balance/deterministic: Balance returns the same value for the
// same inputs.

// The derived checks of Ledger that witness their laws at the shallow end:
//
// Ledger.Credit/smoke
//
// Ledger.Balance/smoke
//
// Ledger.Balance/deterministic

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
	Credit  ledgerCredit
	Balance ledgerBalance
}{
	Credit:  ledgerCredit{"Ledger.Credit"},
	Balance: ledgerBalance{"Ledger.Balance"},
}

// ledgerCredit is the name of Ledger.Credit, with the identities of its
// checks.
type ledgerCredit struct {
	conform.Callable[Ledger]
}

// Smoke returns the identity of the check that Credit returns without a
// panic for derived inputs.
func (c ledgerCredit) Smoke() conform.ID[Ledger] {
	return c.Check("smoke")
}

// ledgerBalance is the name of Ledger.Balance, with the identities of its
// checks.
type ledgerBalance struct {
	conform.Callable[Ledger]
}

// Smoke returns the identity of the check that Balance returns without a
// panic for derived inputs.
func (c ledgerBalance) Smoke() conform.ID[Ledger] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Balance returns the
// same value for the same inputs.
func (c ledgerBalance) Deterministic() conform.ID[Ledger] {
	return c.Check("deterministic")
}

// LedgerSuite returns the checks of Ledger, and the calls and the budgets
// of its callables.
func LedgerSuite() conform.Suite[Ledger] {
	return conform.Suite[Ledger]{
		Declaration: "Ledger",
		Package:     "example.com/shop/ledger",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Ledger]{
			{
				ID:          LedgerChecks.Credit.Smoke(),
				Statement:   "Credit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LedgerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Credit(string("account-a"), int(42))
						return nil
					})
				},
				Defect: ledgerCreditPanics,
			},
			{
				ID:          LedgerChecks.Balance.Smoke(),
				Statement:   "Balance returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LedgerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Balance(string("account-a"))
						return nil
					})
				},
				Defect: ledgerBalancePanics,
			},
			{
				ID:          LedgerChecks.Balance.Deterministic(),
				Statement:   "Balance returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Ledger]{LedgerChecks.Balance.Smoke()},
				Apply: func(h *LedgerHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Balance(string("account-a"))
						return r0, nil
					})
				},
				Defect: ledgerBalanceDrifts,
			},
		},
		Calls: []conform.Call[Ledger]{
			{
				Calls:    ledgerCreditCalls,
				Over:     newLedgerOver,
				Callable: LedgerChecks.Credit.Callable,
			},
			{
				Calls:    ledgerBalanceCalls,
				Over:     newLedgerOver,
				Callable: LedgerChecks.Balance.Callable,
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

// ledgerCreditCalls generates the calls of Ledger.Credit with the inputs
// that each case draws.
var ledgerCreditCalls = prop.Composite(func(c *prop.Case) func(context.Context, Ledger) error {
	account := c.Draw(prop.Of[string](), "account")
	amount := c.Draw(prop.Of[int](), "amount")
	return func(_ context.Context, subject Ledger) error {
		subject.Credit(account, amount)
		return nil
	}
})

// ledgerBalanceCalls generates the calls of Ledger.Balance with the inputs
// that each case draws.
var ledgerBalanceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Ledger) error {
	account := c.Draw(prop.Of[string](), "account")
	return func(_ context.Context, subject Ledger) error {
		_ = subject.Balance(account)
		return nil
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

// Credit calls spend and returns its zero results.
func (o ledgerOver) Credit(string, int) {
	o.spend()
}

// Balance calls spend and returns its zero results.
func (o ledgerOver) Balance(string) (r0 int) {
	o.spend()
	return r0
}

// newLedgerOver returns the planted defect of a ceiling of Ledger, whose
// callables call spend.
func newLedgerOver(_ *LedgerHandle, spend func()) Ledger {
	return ledgerOver{spend: spend}
}

// ledgerCreditPanicsDefect is the planted defect of a Ledger whose Credit
// panics.
type ledgerCreditPanicsDefect struct {
	ledgerOverEmbedded
}

// Credit is the override of the planted defect.
func (*ledgerCreditPanicsDefect) Credit(p0 string, p1 int) {
	panic("the planted defect of a Ledger whose Credit panics")
}

// ledgerCreditPanics returns the planted defect of a Ledger whose Credit
// panics over a fresh subject.
func ledgerCreditPanics(h *LedgerHandle) Ledger {
	return &ledgerCreditPanicsDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerBalancePanicsDefect is the planted defect of a Ledger whose
// Balance panics.
type ledgerBalancePanicsDefect struct {
	ledgerOverEmbedded
}

// Balance is the override of the planted defect.
func (*ledgerBalancePanicsDefect) Balance(p0 string) (r0 int) {
	panic("the planted defect of a Ledger whose Balance panics")
}

// ledgerBalancePanics returns the planted defect of a Ledger whose Balance
// panics over a fresh subject.
func ledgerBalancePanics(h *LedgerHandle) Ledger {
	return &ledgerBalancePanicsDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerBalanceDriftsDefect is the planted defect of a Ledger whose
// Balance returns another value on the repeat.
type ledgerBalanceDriftsDefect struct {
	ledgerOverEmbedded
	calls int
}

// Balance is the override of the planted defect.
func (d *ledgerBalanceDriftsDefect) Balance(p0 string) (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// ledgerBalanceDrifts returns the planted defect of a Ledger whose Balance
// returns another value on the repeat over a fresh subject.
func ledgerBalanceDrifts(h *LedgerHandle) Ledger {
	return &ledgerBalanceDriftsDefect{ledgerOverEmbedded: h.Fresh()}
}
