package leadertest_test

import (
	"testing"

	"example.com/shop/leader"
	"example.com/shop/leader/leadertest"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// primarySubjectEmbedded is leader.Primary under the name of the field
// that primarySubject embeds.
type primarySubjectEmbedded = leader.Primary

// primarySubject is a value of Primary that only TestPrimarySuite builds.
type primarySubject struct{ primarySubjectEmbedded }

// TestPrimarySuite checks the suite of Primary and the two forms of its
// application.
func TestPrimarySuite(t *testing.T) {
	t.Parallel()

	t.Run("PrimarySuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, leadertest.PrimarySuite(), falsifyOptions()...)
	})

	t.Run("RunPrimarySuite", func(t *testing.T) {
		t.Parallel()
		var built *primarySubject
		leadertest.RunPrimarySuite(
			t,
			func() *primarySubject {
				built = &primarySubject{}
				return built
			},
			conform.With(leadertest.PrimaryCheck{
				ID:        "Primary/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *leadertest.PrimaryHandle) {
					got := h.Subject()
					assert.Equal(h, got, leader.Primary(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*leadertest.PrimaryHandle) leader.Primary { return &primarySubject{} },
			}),
			conform.ExcuseSuite[leader.Primary]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunPrimarySuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *primarySubject
		leadertest.RunPrimarySuiteFromHandle(
			t,
			func(*leadertest.PrimaryHandle) *primarySubject {
				built = &primarySubject{}
				return built
			},
			conform.With(leadertest.PrimaryCheck{
				ID:        "Primary/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *leadertest.PrimaryHandle) {
					got := h.Subject()
					assert.Equal(h, got, leader.Primary(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*leadertest.PrimaryHandle) leader.Primary { return &primarySubject{} },
			}),
			conform.ExcuseSuite[leader.Primary]("the test applies the check of its option alone"),
		)
	})

	t.Run("PrimaryChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(leadertest.PrimaryChecks.Commit.AllocationBudget()),
			"Primary.Commit/allocation-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// TestPrimarySuiteAllocs falsifies the allocation ceilings and the checks
// that count allocations of the suite of Primary.
func TestPrimarySuiteAllocs(t *testing.T) {
	t.Run("PrimarySuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, leadertest.PrimarySuite(), falsifyOptions()...)
	})
}

// primarySuite is the suite of Primary that the tests of the planted
// defect of a ceiling apply.
var primarySuite = leadertest.PrimarySuite()

// primaryIdle returns the planted defect of a ceiling of Primary with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func primaryIdle(h *leadertest.PrimaryHandle) leader.Primary {
	return primarySuite.Calls[0].Over(h, func() {})
}

// TestPrimarySuiteBudgets applies the allocation ceilings of the suite of
// Primary to primaryIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestPrimarySuiteBudgets(t *testing.T) {
	leadertest.RunPrimaryBudgets(
		t,
		primaryIdle,
		conform.ExcuseSuite[leader.Primary](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzPrimarySuiteCommit registers the fuzz target of Primary.Commit on
// primaryIdle, with every check of the suite excused.
func FuzzPrimarySuiteCommit(f *testing.F) {
	leadertest.FuzzPrimarySuite(
		f,
		leadertest.PrimaryChecks.Commit.Callable,
		primaryIdle,
		conform.ExcuseSuite[leader.Primary](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkPrimarySuite measures each callable of the suite of Primary on
// primaryIdle.
func BenchmarkPrimarySuite(b *testing.B) {
	b.Run("Commit", func(b *testing.B) {
		leadertest.BenchPrimarySuite(b, leadertest.PrimaryChecks.Commit.Callable, primaryIdle)
	})
}
