package isolated_test

import (
	"testing"

	"example.com/shop/isolated"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// listsSubjectEmbedded is isolated.Lists under the name of the field that
// listsSubject embeds.
type listsSubjectEmbedded = isolated.Lists

// listsSubject is a value of Lists that only TestListsSuite builds.
type listsSubject struct{ listsSubjectEmbedded }

// TestListsSuite checks the suite of Lists and the two forms of its
// application.
func TestListsSuite(t *testing.T) {
	t.Parallel()

	t.Run("ListsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, isolated.ListsSuite())
	})

	t.Run("RunListsSuite", func(t *testing.T) {
		t.Parallel()
		var built *listsSubject
		isolated.RunListsSuite(
			t,
			func() *listsSubject {
				built = &listsSubject{}
				return built
			},
			conform.With(isolated.ListsCheck{
				ID:        "Lists/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.ListsHandle) {
					got := h.Subject()
					assert.Equal(h, got, isolated.Lists(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*isolated.ListsHandle) isolated.Lists { return &listsSubject{} },
			}),
			conform.Excuse(
				isolated.ListsChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.Serializable(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunListsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *listsSubject
		isolated.RunListsSuiteFromHandle(
			t,
			func(*isolated.ListsHandle) *listsSubject {
				built = &listsSubject{}
				return built
			},
			conform.With(isolated.ListsCheck{
				ID:        "Lists/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.ListsHandle) {
					got := h.Subject()
					assert.Equal(h, got, isolated.Lists(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*isolated.ListsHandle) isolated.Lists { return &listsSubject{} },
			}),
			conform.Excuse(
				isolated.ListsChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Commit.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Rollback.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Append.Serializable(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.ListsChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// listsIdle returns the planted defect of a ceiling of Lists with a spend
// of nothing, a subject that answers each call at once.
func listsIdle(h *isolated.ListsHandle) isolated.Lists {
	return isolated.ListsSuite().Calls[0].Over(h, func() {})
}

// TestListsSuiteBudgets applies the allocation ceilings of the suite of
// Lists to listsIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestListsSuiteBudgets(t *testing.T) {
	isolated.RunListsBudgets(
		t,
		listsIdle,
		conform.Ceilings(
			isolated.ListsChecks.Begin.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			isolated.ListsChecks.Commit.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			isolated.ListsChecks.Rollback.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			isolated.ListsChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(isolated.ListsChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			isolated.ListsChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Serializable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzListsSuiteBegin registers the fuzz target of Lists.Begin on
// listsIdle, with every check of the suite excused.
func FuzzListsSuiteBegin(f *testing.F) {
	isolated.FuzzListsSuite(
		f,
		isolated.ListsChecks.Begin.Callable,
		listsIdle,
		conform.Excuse(
			isolated.ListsChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Serializable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzListsSuiteCommit registers the fuzz target of Lists.Commit on
// listsIdle, with every check of the suite excused.
func FuzzListsSuiteCommit(f *testing.F) {
	isolated.FuzzListsSuite(
		f,
		isolated.ListsChecks.Commit.Callable,
		listsIdle,
		conform.Excuse(
			isolated.ListsChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Serializable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzListsSuiteRollback registers the fuzz target of Lists.Rollback on
// listsIdle, with every check of the suite excused.
func FuzzListsSuiteRollback(f *testing.F) {
	isolated.FuzzListsSuite(
		f,
		isolated.ListsChecks.Rollback.Callable,
		listsIdle,
		conform.Excuse(
			isolated.ListsChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Serializable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzListsSuiteAppend registers the fuzz target of Lists.Append on
// listsIdle, with every check of the suite excused.
func FuzzListsSuiteAppend(f *testing.F) {
	isolated.FuzzListsSuite(
		f,
		isolated.ListsChecks.Append.Callable,
		listsIdle,
		conform.Excuse(
			isolated.ListsChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Serializable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzListsSuiteGet registers the fuzz target of Lists.Get on listsIdle,
// with every check of the suite excused.
func FuzzListsSuiteGet(f *testing.F) {
	isolated.FuzzListsSuite(
		f,
		isolated.ListsChecks.Get.Callable,
		listsIdle,
		conform.Excuse(
			isolated.ListsChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Begin.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Commit.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Rollback.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Append.Serializable(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.ListsChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkListsSuite measures each callable of the suite of Lists on
// listsIdle.
func BenchmarkListsSuite(b *testing.B) {
	b.Run("Begin", func(b *testing.B) {
		isolated.BenchListsSuite(b, isolated.ListsChecks.Begin.Callable, listsIdle)
	})
	b.Run("Commit", func(b *testing.B) {
		isolated.BenchListsSuite(b, isolated.ListsChecks.Commit.Callable, listsIdle)
	})
	b.Run("Rollback", func(b *testing.B) {
		isolated.BenchListsSuite(b, isolated.ListsChecks.Rollback.Callable, listsIdle)
	})
	b.Run("Append", func(b *testing.B) {
		isolated.BenchListsSuite(b, isolated.ListsChecks.Append.Callable, listsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		isolated.BenchListsSuite(b, isolated.ListsChecks.Get.Callable, listsIdle)
	})
}

// registersSubjectEmbedded is isolated.Registers under the name of the
// field that registersSubject embeds.
type registersSubjectEmbedded = isolated.Registers

// registersSubject is a value of Registers that only TestRegistersSuite
// builds.
type registersSubject struct{ registersSubjectEmbedded }

// TestRegistersSuite checks the suite of Registers and the two forms of
// its application.
func TestRegistersSuite(t *testing.T) {
	t.Parallel()

	t.Run("RegistersSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, isolated.RegistersSuite())
	})

	t.Run("RunRegistersSuite", func(t *testing.T) {
		t.Parallel()
		var built *registersSubject
		isolated.RunRegistersSuite(
			t,
			func() *registersSubject {
				built = &registersSubject{}
				return built
			},
			conform.With(isolated.RegistersCheck{
				ID:        "Registers/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.RegistersHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						isolated.Registers(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*isolated.RegistersHandle) isolated.Registers { return &registersSubject{} },
			}),
			conform.Excuse(
				isolated.RegistersChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Begin.SnapshotIsolation(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Latest.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Latest.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunRegistersSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *registersSubject
		isolated.RunRegistersSuiteFromHandle(
			t,
			func(*isolated.RegistersHandle) *registersSubject {
				built = &registersSubject{}
				return built
			},
			conform.With(isolated.RegistersCheck{
				ID:        "Registers/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.RegistersHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						isolated.Registers(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*isolated.RegistersHandle) isolated.Registers { return &registersSubject{} },
			}),
			conform.Excuse(
				isolated.RegistersChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Begin.SnapshotIsolation(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Latest.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.RegistersChecks.Latest.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// registersIdle returns the planted defect of a ceiling of Registers with
// a spend of nothing, a subject that answers each call at once.
func registersIdle(h *isolated.RegistersHandle) isolated.Registers {
	return isolated.RegistersSuite().Calls[0].Over(h, func() {})
}

// TestRegistersSuiteBudgets applies the allocation ceilings of the suite
// of Registers to registersIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRegistersSuiteBudgets(t *testing.T) {
	isolated.RunRegistersBudgets(
		t,
		registersIdle,
		conform.Ceilings(
			isolated.RegistersChecks.Begin.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			isolated.RegistersChecks.Latest.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			isolated.RegistersChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Begin.SnapshotIsolation(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Latest.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Latest.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRegistersSuiteBegin registers the fuzz target of Registers.Begin on
// registersIdle, with every check of the suite excused.
func FuzzRegistersSuiteBegin(f *testing.F) {
	isolated.FuzzRegistersSuite(
		f,
		isolated.RegistersChecks.Begin.Callable,
		registersIdle,
		conform.Excuse(
			isolated.RegistersChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Begin.SnapshotIsolation(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Latest.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Latest.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRegistersSuiteLatest registers the fuzz target of Registers.Latest
// on registersIdle, with every check of the suite excused.
func FuzzRegistersSuiteLatest(f *testing.F) {
	isolated.FuzzRegistersSuite(
		f,
		isolated.RegistersChecks.Latest.Callable,
		registersIdle,
		conform.Excuse(
			isolated.RegistersChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Begin.SnapshotIsolation(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Latest.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.RegistersChecks.Latest.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkRegistersSuite measures each callable of the suite of Registers
// on registersIdle.
func BenchmarkRegistersSuite(b *testing.B) {
	b.Run("Begin", func(b *testing.B) {
		isolated.BenchRegistersSuite(b, isolated.RegistersChecks.Begin.Callable, registersIdle)
	})
	b.Run("Latest", func(b *testing.B) {
		isolated.BenchRegistersSuite(b, isolated.RegistersChecks.Latest.Callable, registersIdle)
	})
}

// latchSubjectEmbedded is isolated.Latch under the name of the field that
// latchSubject embeds.
type latchSubjectEmbedded = isolated.Latch

// latchSubject is a value of Latch that only TestLatchSuite builds.
type latchSubject struct{ latchSubjectEmbedded }

// TestLatchSuite checks the suite of Latch and the two forms of its
// application.
func TestLatchSuite(t *testing.T) {
	t.Parallel()

	t.Run("LatchSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, isolated.LatchSuite())
	})

	t.Run("RunLatchSuite", func(t *testing.T) {
		t.Parallel()
		var built *latchSubject
		isolated.RunLatchSuite(
			t,
			func() *latchSubject {
				built = &latchSubject{}
				return built
			},
			conform.With(isolated.LatchCheck{
				ID:        "Latch/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.LatchHandle) {
					got := h.Subject()
					assert.Equal(h, got, isolated.Latch(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*isolated.LatchHandle) isolated.Latch { return &latchSubject{} },
			}),
			conform.Excuse(
				isolated.LatchChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLatchSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *latchSubject
		isolated.RunLatchSuiteFromHandle(
			t,
			func(*isolated.LatchHandle) *latchSubject {
				built = &latchSubject{}
				return built
			},
			conform.With(isolated.LatchCheck{
				ID:        "Latch/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.LatchHandle) {
					got := h.Subject()
					assert.Equal(h, got, isolated.Latch(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*isolated.LatchHandle) isolated.Latch { return &latchSubject{} },
			}),
			conform.Excuse(
				isolated.LatchChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Commit.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Commit.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Rollback.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.LatchChecks.Rollback.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// latchIdle returns the planted defect of a ceiling of Latch with a spend
// of nothing, a subject that answers each call at once.
func latchIdle(h *isolated.LatchHandle) isolated.Latch {
	return isolated.LatchSuite().Calls[0].Over(h, func() {})
}

// TestLatchSuiteBudgets applies the allocation ceilings of the suite of
// Latch to latchIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestLatchSuiteBudgets(t *testing.T) {
	isolated.RunLatchBudgets(
		t,
		latchIdle,
		conform.Ceilings(
			isolated.LatchChecks.Begin.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			isolated.LatchChecks.Commit.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			isolated.LatchChecks.Rollback.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			isolated.LatchChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLatchSuiteBegin registers the fuzz target of Latch.Begin on
// latchIdle, with every check of the suite excused.
func FuzzLatchSuiteBegin(f *testing.F) {
	isolated.FuzzLatchSuite(
		f,
		isolated.LatchChecks.Begin.Callable,
		latchIdle,
		conform.Excuse(
			isolated.LatchChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLatchSuiteCommit registers the fuzz target of Latch.Commit on
// latchIdle, with every check of the suite excused.
func FuzzLatchSuiteCommit(f *testing.F) {
	isolated.FuzzLatchSuite(
		f,
		isolated.LatchChecks.Commit.Callable,
		latchIdle,
		conform.Excuse(
			isolated.LatchChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLatchSuiteRollback registers the fuzz target of Latch.Rollback on
// latchIdle, with every check of the suite excused.
func FuzzLatchSuiteRollback(f *testing.F) {
	isolated.FuzzLatchSuite(
		f,
		isolated.LatchChecks.Rollback.Callable,
		latchIdle,
		conform.Excuse(
			isolated.LatchChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Commit.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.LatchChecks.Rollback.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLatchSuite measures each callable of the suite of Latch on
// latchIdle.
func BenchmarkLatchSuite(b *testing.B) {
	b.Run("Begin", func(b *testing.B) {
		isolated.BenchLatchSuite(b, isolated.LatchChecks.Begin.Callable, latchIdle)
	})
	b.Run("Commit", func(b *testing.B) {
		isolated.BenchLatchSuite(b, isolated.LatchChecks.Commit.Callable, latchIdle)
	})
	b.Run("Rollback", func(b *testing.B) {
		isolated.BenchLatchSuite(b, isolated.LatchChecks.Rollback.Callable, latchIdle)
	})
}

// driverSubjectEmbedded is isolated.Driver under the name of the field
// that driverSubject embeds.
type driverSubjectEmbedded = isolated.Driver

// driverSubject is a value of Driver that only TestDriverSuite builds.
type driverSubject struct{ driverSubjectEmbedded }

// TestDriverSuite checks the suite of Driver and the two forms of its
// application.
func TestDriverSuite(t *testing.T) {
	t.Parallel()

	t.Run("DriverSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, isolated.DriverSuite())
	})

	t.Run("RunDriverSuite", func(t *testing.T) {
		t.Parallel()
		var built *driverSubject
		isolated.RunDriverSuite(
			t,
			func() *driverSubject {
				built = &driverSubject{}
				return built
			},
			conform.With(isolated.DriverCheck{
				ID:        "Driver/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.DriverHandle) {
					got := h.Subject()
					assert.Equal(h, got, isolated.Driver(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*isolated.DriverHandle) isolated.Driver { return &driverSubject{} },
			}),
			conform.Excuse(
				isolated.DriverChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.DriverChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunDriverSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *driverSubject
		isolated.RunDriverSuiteFromHandle(
			t,
			func(*isolated.DriverHandle) *driverSubject {
				built = &driverSubject{}
				return built
			},
			conform.With(isolated.DriverCheck{
				ID:        "Driver/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *isolated.DriverHandle) {
					got := h.Subject()
					assert.Equal(h, got, isolated.Driver(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*isolated.DriverHandle) isolated.Driver { return &driverSubject{} },
			}),
			conform.Excuse(
				isolated.DriverChecks.Begin.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				isolated.DriverChecks.Begin.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// driverIdle returns the planted defect of a ceiling of Driver with a
// spend of nothing, a subject that answers each call at once.
func driverIdle(h *isolated.DriverHandle) isolated.Driver {
	return isolated.DriverSuite().Calls[0].Over(h, func() {})
}

// TestDriverSuiteBudgets applies the allocation ceilings of the suite of
// Driver to driverIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestDriverSuiteBudgets(t *testing.T) {
	isolated.RunDriverBudgets(
		t,
		driverIdle,
		conform.Ceilings(
			isolated.DriverChecks.Begin.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			isolated.DriverChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.DriverChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDriverSuiteBegin registers the fuzz target of Driver.Begin on
// driverIdle, with every check of the suite excused.
func FuzzDriverSuiteBegin(f *testing.F) {
	isolated.FuzzDriverSuite(
		f,
		isolated.DriverChecks.Begin.Callable,
		driverIdle,
		conform.Excuse(
			isolated.DriverChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.DriverChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDriverSuiteCommit registers the fuzz target of Driver.Commit on
// driverIdle, with every check of the suite excused.
func FuzzDriverSuiteCommit(f *testing.F) {
	isolated.FuzzDriverSuite(
		f,
		isolated.DriverChecks.Commit.Callable,
		driverIdle,
		conform.Excuse(
			isolated.DriverChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.DriverChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDriverSuiteRollback registers the fuzz target of Driver.Rollback on
// driverIdle, with every check of the suite excused.
func FuzzDriverSuiteRollback(f *testing.F) {
	isolated.FuzzDriverSuite(
		f,
		isolated.DriverChecks.Rollback.Callable,
		driverIdle,
		conform.Excuse(
			isolated.DriverChecks.Begin.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			isolated.DriverChecks.Begin.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkDriverSuite measures each callable of the suite of Driver on
// driverIdle.
func BenchmarkDriverSuite(b *testing.B) {
	b.Run("Begin", func(b *testing.B) {
		isolated.BenchDriverSuite(b, isolated.DriverChecks.Begin.Callable, driverIdle)
	})
	b.Run("Commit", func(b *testing.B) {
		isolated.BenchDriverSuite(b, isolated.DriverChecks.Commit.Callable, driverIdle)
	})
	b.Run("Rollback", func(b *testing.B) {
		isolated.BenchDriverSuite(b, isolated.DriverChecks.Rollback.Callable, driverIdle)
	})
}
