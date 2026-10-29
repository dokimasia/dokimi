package ledger

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// ledgerSubjectEmbedded is Ledger under the name of the field that
// ledgerSubject embeds.
type ledgerSubjectEmbedded = Ledger

// ledgerSubject is a value of Ledger that only TestLedgerSuite builds.
type ledgerSubject struct{ ledgerSubjectEmbedded }

// TestLedgerSuite checks the suite of Ledger and the two forms of its
// application.
func TestLedgerSuite(t *testing.T) {
	t.Parallel()

	t.Run("LedgerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, LedgerSuite())
	})

	t.Run("RunLedgerSuite", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		RunLedgerSuite(
			t,
			func() *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LedgerHandle) Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				LedgerChecks.Credit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Balance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Balance.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLedgerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		RunLedgerSuiteFromHandle(
			t,
			func(*LedgerHandle) *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*LedgerHandle) Ledger { return &ledgerSubject{} },
			}),
			conform.Excuse(
				LedgerChecks.Credit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Balance.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				LedgerChecks.Balance.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// ledgerIdle returns the planted defect of a ceiling of Ledger with a
// spend of nothing, a subject that answers each call at once.
func ledgerIdle(h *LedgerHandle) Ledger {
	return LedgerSuite().Calls[0].Over(h, func() {})
}

// TestLedgerSuiteBudgets applies the allocation ceilings of the suite of
// Ledger to ledgerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLedgerSuiteBudgets(t *testing.T) {
	RunLedgerBudgets(
		t,
		ledgerIdle,
		conform.Ceilings(LedgerChecks.Credit.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(LedgerChecks.Balance.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			LedgerChecks.Credit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Balance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Balance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteCredit registers the fuzz target of Ledger.Credit on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteCredit(f *testing.F) {
	FuzzLedgerSuite(
		f,
		LedgerChecks.Credit.Callable,
		ledgerIdle,
		conform.Excuse(
			LedgerChecks.Credit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Balance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Balance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLedgerSuiteBalance registers the fuzz target of Ledger.Balance on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteBalance(f *testing.F) {
	FuzzLedgerSuite(
		f,
		LedgerChecks.Balance.Callable,
		ledgerIdle,
		conform.Excuse(
			LedgerChecks.Credit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Balance.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			LedgerChecks.Balance.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLedgerSuite measures each callable of the suite of Ledger on
// ledgerIdle.
func BenchmarkLedgerSuite(b *testing.B) {
	b.Run("Credit", func(b *testing.B) {
		BenchLedgerSuite(b, LedgerChecks.Credit.Callable, ledgerIdle)
	})
	b.Run("Balance", func(b *testing.B) {
		BenchLedgerSuite(b, LedgerChecks.Balance.Callable, ledgerIdle)
	})
}

// TestBookConformsToLedger applies the checks of the suite of Ledger to
// Book.
func TestBookConformsToLedger(t *testing.T) {
	t.Parallel()

	RunLedgerSuite(t, NewBook)
}

// BenchmarkBookConformsToLedger measures each callable of the suite of
// Ledger on Book.
func BenchmarkBookConformsToLedger(b *testing.B) {
	b.Run("Credit", func(b *testing.B) {
		BenchLedgerSuite(
			b,
			LedgerChecks.Credit.Callable,
			func(*LedgerHandle) *Book { return NewBook() },
		)
	})
	b.Run("Balance", func(b *testing.B) {
		BenchLedgerSuite(
			b,
			LedgerChecks.Balance.Callable,
			func(*LedgerHandle) *Book { return NewBook() },
		)
	})
}

// FuzzBookConformsToLedgerCredit fuzzes the property checks of
// Ledger.Credit on Book.
func FuzzBookConformsToLedgerCredit(f *testing.F) {
	FuzzLedgerSuite(f, LedgerChecks.Credit.Callable, func(*LedgerHandle) *Book { return NewBook() })
}

// FuzzBookConformsToLedgerBalance fuzzes the property checks of
// Ledger.Balance on Book.
func FuzzBookConformsToLedgerBalance(f *testing.F) {
	FuzzLedgerSuite(
		f,
		LedgerChecks.Balance.Callable,
		func(*LedgerHandle) *Book { return NewBook() },
	)
}
