package ledger_test

import (
	"testing"

	"example.com/shop/ledger"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// ledgerSubjectEmbedded is ledger.Ledger under the name of the field that
// ledgerSubject embeds.
type ledgerSubjectEmbedded = ledger.Ledger

// ledgerSubject is a value of Ledger that only TestLedgerSuite builds.
type ledgerSubject struct{ ledgerSubjectEmbedded }

// TestLedgerSuite checks the suite of Ledger and the two forms of its
// application.
func TestLedgerSuite(t *testing.T) {
	t.Parallel()

	t.Run("LedgerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ledger.LedgerSuite())
	})

	t.Run("RunLedgerSuite", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		ledger.RunLedgerSuite(
			t,
			func() *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(ledger.LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ledger.LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, ledger.Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ledger.LedgerHandle) ledger.Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				ledger.LedgerChecks.Credit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ledger.LedgerChecks.Balance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ledger.LedgerChecks.Balance.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLedgerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		ledger.RunLedgerSuiteFromHandle(
			t,
			func(*ledger.LedgerHandle) *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(ledger.LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ledger.LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, ledger.Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ledger.LedgerHandle) ledger.Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				ledger.LedgerChecks.Credit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ledger.LedgerChecks.Balance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				ledger.LedgerChecks.Balance.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// ledgerIdle returns the planted defect of a ceiling of Ledger with a
// spend of nothing, a subject that answers each call at once.
func ledgerIdle(h *ledger.LedgerHandle) ledger.Ledger {
	return ledger.LedgerSuite().Calls[0].Over(h, func() {})
}

// TestLedgerSuiteBudgets applies the allocation ceilings of the suite of
// Ledger to ledgerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLedgerSuiteBudgets(t *testing.T) {
	ledger.RunLedgerBudgets(
		t,
		ledgerIdle,
		conform.Ceilings(
			ledger.LedgerChecks.Credit.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			ledger.LedgerChecks.Balance.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			ledger.LedgerChecks.Credit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ledger.LedgerChecks.Balance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ledger.LedgerChecks.Balance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteCredit registers the fuzz target of Ledger.Credit on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteCredit(f *testing.F) {
	ledger.FuzzLedgerSuite(
		f,
		ledger.LedgerChecks.Credit.Callable,
		ledgerIdle,
		conform.Excuse(
			ledger.LedgerChecks.Credit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ledger.LedgerChecks.Balance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ledger.LedgerChecks.Balance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteBalance registers the fuzz target of Ledger.Balance on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteBalance(f *testing.F) {
	ledger.FuzzLedgerSuite(
		f,
		ledger.LedgerChecks.Balance.Callable,
		ledgerIdle,
		conform.Excuse(
			ledger.LedgerChecks.Credit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ledger.LedgerChecks.Balance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			ledger.LedgerChecks.Balance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLedgerSuite measures each callable of the suite of Ledger on
// ledgerIdle.
func BenchmarkLedgerSuite(b *testing.B) {
	b.Run("Credit", func(b *testing.B) {
		ledger.BenchLedgerSuite(b, ledger.LedgerChecks.Credit.Callable, ledgerIdle)
	})
	b.Run("Balance", func(b *testing.B) {
		ledger.BenchLedgerSuite(b, ledger.LedgerChecks.Balance.Callable, ledgerIdle)
	})
}

// TestBookConformsToLedger applies the checks of the suite of Ledger to
// Book.
func TestBookConformsToLedger(t *testing.T) {
	t.Parallel()

	ledger.RunLedgerSuite(t, ledger.NewBook)
}

// BenchmarkBookConformsToLedger measures each callable of the suite of
// Ledger on Book.
func BenchmarkBookConformsToLedger(b *testing.B) {
	b.Run("Credit", func(b *testing.B) {
		ledger.BenchLedgerSuite(
			b,
			ledger.LedgerChecks.Credit.Callable,
			func(*ledger.LedgerHandle) *ledger.Book { return ledger.NewBook() },
		)
	})
	b.Run("Balance", func(b *testing.B) {
		ledger.BenchLedgerSuite(
			b,
			ledger.LedgerChecks.Balance.Callable,
			func(*ledger.LedgerHandle) *ledger.Book { return ledger.NewBook() },
		)
	})
}

// FuzzBookConformsToLedgerCredit fuzzes the property checks of
// Ledger.Credit on Book.
func FuzzBookConformsToLedgerCredit(f *testing.F) {
	ledger.FuzzLedgerSuite(
		f,
		ledger.LedgerChecks.Credit.Callable,
		func(*ledger.LedgerHandle) *ledger.Book { return ledger.NewBook() },
	)
}

// FuzzBookConformsToLedgerBalance fuzzes the property checks of
// Ledger.Balance on Book.
func FuzzBookConformsToLedgerBalance(f *testing.F) {
	ledger.FuzzLedgerSuite(
		f,
		ledger.LedgerChecks.Balance.Callable,
		func(*ledger.LedgerHandle) *ledger.Book { return ledger.NewBook() },
	)
}
