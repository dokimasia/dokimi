package durable_test

import (
	"testing"

	"example.com/shop/durable"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// keeperSubjectEmbedded is durable.Keeper under the name of the field that
// keeperSubject embeds.
type keeperSubjectEmbedded = durable.Keeper

// keeperSubject is a value of Keeper that only TestKeeperSuite builds.
type keeperSubject struct{ keeperSubjectEmbedded }

// TestKeeperSuite checks the suite of Keeper and the two forms of its
// application.
func TestKeeperSuite(t *testing.T) {
	t.Parallel()

	t.Run("KeeperSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, durable.KeeperSuite())
	})

	t.Run("RunKeeperSuite", func(t *testing.T) {
		t.Parallel()
		var built *keeperSubject
		durable.RunKeeperSuite(
			t,
			func() *keeperSubject {
				built = &keeperSubject{}
				return built
			},
			conform.With(durable.KeeperCheck{
				ID:        "Keeper/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *durable.KeeperHandle) {
					got := h.Subject()
					assert.Equal(h, got, durable.Keeper(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*durable.KeeperHandle) durable.Keeper { return &keeperSubject{} },
			}),
			conform.Excuse(
				durable.KeeperChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Put.LeavesTheWholeValue(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.KeepsTheLastPut(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunKeeperSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *keeperSubject
		durable.RunKeeperSuiteFromHandle(
			t,
			func(*durable.KeeperHandle) *keeperSubject {
				built = &keeperSubject{}
				return built
			},
			conform.With(durable.KeeperCheck{
				ID:        "Keeper/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *durable.KeeperHandle) {
					got := h.Subject()
					assert.Equal(h, got, durable.Keeper(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*durable.KeeperHandle) durable.Keeper { return &keeperSubject{} },
			}),
			conform.Excuse(
				durable.KeeperChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Put.LeavesTheWholeValue(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				durable.KeeperChecks.Get.KeepsTheLastPut(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("KeeperChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(durable.KeeperChecks.Put.LeavesTheWholeValue()),
			"Keeper.Put/leaves-the-whole-value",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(durable.KeeperChecks.Get.KeepsEachPutAcrossARestart()),
			"Keeper.Get/keeps-each-put-across-a-restart",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(durable.KeeperChecks.Get.KeepsTheLastPut()),
			"Keeper.Get/keeps-the-last-put",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// keeperIdle returns the planted defect of a ceiling of Keeper with a
// spend of nothing, a subject that answers each call at once.
func keeperIdle(h *durable.KeeperHandle) durable.Keeper {
	return durable.KeeperSuite().Calls[0].Over(h, func() {})
}

// TestKeeperSuiteBudgets applies the allocation ceilings of the suite of
// Keeper to keeperIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestKeeperSuiteBudgets(t *testing.T) {
	durable.RunKeeperBudgets(
		t,
		keeperIdle,
		conform.Ceilings(durable.KeeperChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(durable.KeeperChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			durable.KeeperChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Put.LeavesTheWholeValue(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsTheLastPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeeperSuitePut registers the fuzz target of Keeper.Put on
// keeperIdle, with every check of the suite excused.
func FuzzKeeperSuitePut(f *testing.F) {
	durable.FuzzKeeperSuite(
		f,
		durable.KeeperChecks.Put.Callable,
		keeperIdle,
		conform.Excuse(
			durable.KeeperChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Put.LeavesTheWholeValue(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsTheLastPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzKeeperSuiteGet registers the fuzz target of Keeper.Get on
// keeperIdle, with every check of the suite excused.
func FuzzKeeperSuiteGet(f *testing.F) {
	durable.FuzzKeeperSuite(
		f,
		durable.KeeperChecks.Get.Callable,
		keeperIdle,
		conform.Excuse(
			durable.KeeperChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Put.LeavesTheWholeValue(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsTheLastPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkKeeperSuite measures each callable of the suite of Keeper on
// keeperIdle.
func BenchmarkKeeperSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		durable.BenchKeeperSuite(b, durable.KeeperChecks.Put.Callable, keeperIdle)
	})
	b.Run("Get", func(b *testing.B) {
		durable.BenchKeeperSuite(b, durable.KeeperChecks.Get.Callable, keeperIdle)
	})
}
