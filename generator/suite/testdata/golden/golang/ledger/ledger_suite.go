package ledger

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

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

// CreditsInPlace returns the identity of the check that Credit allocates
// nothing once the account has a balance.
func (c ledgerCredit) CreditsInPlace() conform.ID[Ledger] {
	return c.Check("credits-in-place")
}

// ledgerBalance is the name of Ledger.Balance, with the identities of its
// checks.
type ledgerBalance struct {
	conform.Callable[Ledger]
}

// KeepsTheSum returns the identity of the check that Balance returns the
// sum of the credits of the account.
func (c ledgerBalance) KeepsTheSum() conform.ID[Ledger] {
	return c.Check("keeps-the-sum")
}

// LedgerSuite returns the checks of Ledger, and the calls and the budgets
// of its callables.
func LedgerSuite() conform.Suite[Ledger] {
	return conform.Suite[Ledger]{
		Declaration: "Ledger",
		Checks: []conform.Check[Ledger]{
			{
				ID:        LedgerChecks.Credit.CreditsInPlace(),
				Statement: "Credit allocates nothing once the account has a balance",
				Apply:     CreditsInPlace,
				Defect:    copiesBalances,
				Counted:   true,
			},
			{
				ID:        LedgerChecks.Balance.KeepsTheSum(),
				Statement: "Balance returns the sum of the credits of the account",
				Apply:     KeepsTheSum,
				Defect:    dropsCredits,
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
