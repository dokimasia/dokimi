package classified_test

import (
	"testing"

	"example.com/shop/classified"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// shelfSubjectEmbedded is classified.Shelf under the name of the field
// that shelfSubject embeds.
type shelfSubjectEmbedded = classified.Shelf

// shelfSubject is a value of Shelf that only TestShelfSuite builds.
type shelfSubject struct{ shelfSubjectEmbedded }

// TestShelfSuite checks the suite of Shelf and the two forms of its
// application.
func TestShelfSuite(t *testing.T) {
	t.Parallel()

	t.Run("ShelfSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, classified.ShelfSuite())
	})

	t.Run("RunShelfSuite", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		classified.RunShelfSuite(
			t,
			func() *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(classified.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.ShelfHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Shelf(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.ShelfHandle) classified.Shelf { return &shelfSubject{} },
			}),
			conform.Excuse(
				classified.ShelfChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.ReadsEachKey(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.All.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Items.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.OverMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.OverMatchGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Keyed.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.OverMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.OverMatchGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Len.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Len.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Len.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.RespectsItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.RespectsItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.RespectsItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Err.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Err.FreshReportsNoPoison(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunShelfSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		classified.RunShelfSuiteFromHandle(
			t,
			func(*classified.ShelfHandle) *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(classified.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.ShelfHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Shelf(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.ShelfHandle) classified.Shelf { return &shelfSubject{} },
			}),
			conform.Excuse(
				classified.ShelfChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.ReadsEachKey(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.All.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.All.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Items.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Items.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.OverMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.OverMatchGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Keyed.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.OverMatch(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.OverMatchGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Len.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Len.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Len.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.At.RespectsItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.AtOffset.RespectsItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.RespectsItsBounds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Err.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.ShelfChecks.Err.FreshReportsNoPoison(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// shelfIdle returns the planted defect of a ceiling of Shelf with a spend
// of nothing, a subject that answers each call at once.
func shelfIdle(h *classified.ShelfHandle) classified.Shelf {
	return classified.ShelfSuite().Calls[0].Over(h, func() {})
}

// TestShelfSuiteBudgets applies the allocation ceilings of the suite of
// Shelf to shelfIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestShelfSuiteBudgets(t *testing.T) {
	classified.RunShelfBudgets(
		t,
		shelfIdle,
		conform.Ceilings(
			classified.ShelfChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.GetMany.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.All.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Items.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Scan.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Keyed.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Snapshot.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Len.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.At.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.AtOffset.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Slot.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.TracedGet.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Err.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuitePut registers the fuzz target of Shelf.Put on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuitePut(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Put.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteGetMany registers the fuzz target of Shelf.GetMany on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteGetMany(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.GetMany.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteAll registers the fuzz target of Shelf.All on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteAll(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.All.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteItems registers the fuzz target of Shelf.Items on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteItems(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Items.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteScan registers the fuzz target of Shelf.Scan on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteScan(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Scan.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteKeyed registers the fuzz target of Shelf.Keyed on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteKeyed(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Keyed.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteSnapshot registers the fuzz target of Shelf.Snapshot on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteSnapshot(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Snapshot.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteLen registers the fuzz target of Shelf.Len on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteLen(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Len.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteAt registers the fuzz target of Shelf.At on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteAt(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.At.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteAtOffset registers the fuzz target of Shelf.AtOffset on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteAtOffset(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.AtOffset.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteSlot registers the fuzz target of Shelf.Slot on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteSlot(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Slot.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteGet registers the fuzz target of Shelf.Get on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteGet(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Get.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteTracedGet registers the fuzz target of Shelf.TracedGet on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteTracedGet(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.TracedGet.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteErr registers the fuzz target of Shelf.Err on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteErr(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Err.Callable,
		shelfIdle,
		conform.Excuse(
			classified.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.GetMany.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.All.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Scan.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStored(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatch(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.OverMatchGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Len.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.At.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.AtOffset.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBounds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Slot.RespectsItsBoundsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSibling(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.ShelfChecks.Err.FreshReportsNoPoison(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkShelfSuite measures each callable of the suite of Shelf on
// shelfIdle.
func BenchmarkShelfSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Put.Callable, shelfIdle)
	})
	b.Run("GetMany", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.GetMany.Callable, shelfIdle)
	})
	b.Run("All", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.All.Callable, shelfIdle)
	})
	b.Run("Items", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Items.Callable, shelfIdle)
	})
	b.Run("Scan", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Scan.Callable, shelfIdle)
	})
	b.Run("Keyed", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Keyed.Callable, shelfIdle)
	})
	b.Run("Snapshot", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Snapshot.Callable, shelfIdle)
	})
	b.Run("Len", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Len.Callable, shelfIdle)
	})
	b.Run("At", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.At.Callable, shelfIdle)
	})
	b.Run("AtOffset", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.AtOffset.Callable, shelfIdle)
	})
	b.Run("Slot", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Slot.Callable, shelfIdle)
	})
	b.Run("Get", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Get.Callable, shelfIdle)
	})
	b.Run("TracedGet", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.TracedGet.Callable, shelfIdle)
	})
	b.Run("Err", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Err.Callable, shelfIdle)
	})
}

// tagsSubjectEmbedded is classified.Tags under the name of the field that
// tagsSubject embeds.
type tagsSubjectEmbedded = classified.Tags

// tagsSubject is a value of Tags that only TestTagsSuite builds.
type tagsSubject struct{ tagsSubjectEmbedded }

// TestTagsSuite checks the suite of Tags and the two forms of its
// application.
func TestTagsSuite(t *testing.T) {
	t.Parallel()

	t.Run("TagsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, classified.TagsSuite())
	})

	t.Run("RunTagsSuite", func(t *testing.T) {
		t.Parallel()
		var built *tagsSubject
		classified.RunTagsSuite(
			t,
			func() *tagsSubject {
				built = &tagsSubject{}
				return built
			},
			conform.With(classified.TagsCheck{
				ID:        "Tags/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.TagsHandle) {
					got := h.Subject()
					assert.Equal(h, got, classified.Tags(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*classified.TagsHandle) classified.Tags { return &tagsSubject{} },
			}),
			conform.Excuse(
				classified.TagsChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.ReadsEachKey(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.ReadsEachKeyGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTagsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *tagsSubject
		classified.RunTagsSuiteFromHandle(
			t,
			func(*classified.TagsHandle) *tagsSubject {
				built = &tagsSubject{}
				return built
			},
			conform.With(classified.TagsCheck{
				ID:        "Tags/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.TagsHandle) {
					got := h.Subject()
					assert.Equal(h, got, classified.Tags(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*classified.TagsHandle) classified.Tags { return &tagsSubject{} },
			}),
			conform.Excuse(
				classified.TagsChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Add.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Add.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Add.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.ReadsEachKey(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.TagsChecks.Has.ReadsEachKeyGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// tagsIdle returns the planted defect of a ceiling of Tags with a spend of
// nothing, a subject that answers each call at once.
func tagsIdle(h *classified.TagsHandle) classified.Tags {
	return classified.TagsSuite().Calls[0].Over(h, func() {})
}

// TestTagsSuiteBudgets applies the allocation ceilings of the suite of
// Tags to tagsIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestTagsSuiteBudgets(t *testing.T) {
	classified.RunTagsBudgets(
		t,
		tagsIdle,
		conform.Ceilings(
			classified.TagsChecks.Add.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.TagsChecks.Has.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			classified.TagsChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTagsSuiteAdd registers the fuzz target of Tags.Add on tagsIdle, with
// every check of the suite excused.
func FuzzTagsSuiteAdd(f *testing.F) {
	classified.FuzzTagsSuite(
		f,
		classified.TagsChecks.Add.Callable,
		tagsIdle,
		conform.Excuse(
			classified.TagsChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzTagsSuiteHas registers the fuzz target of Tags.Has on tagsIdle, with
// every check of the suite excused.
func FuzzTagsSuiteHas(f *testing.F) {
	classified.FuzzTagsSuite(
		f,
		classified.TagsChecks.Has.Callable,
		tagsIdle,
		conform.Excuse(
			classified.TagsChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Add.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.ReadsEachKey(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.TagsChecks.Has.ReadsEachKeyGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkTagsSuite measures each callable of the suite of Tags on
// tagsIdle.
func BenchmarkTagsSuite(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		classified.BenchTagsSuite(b, classified.TagsChecks.Add.Callable, tagsIdle)
	})
	b.Run("Has", func(b *testing.B) {
		classified.BenchTagsSuite(b, classified.TagsChecks.Has.Callable, tagsIdle)
	})
}

// vaultSubjectEmbedded is classified.Vault under the name of the field
// that vaultSubject embeds.
type vaultSubjectEmbedded = classified.Vault

// vaultSubject is a value of Vault that only TestVaultSuite builds.
type vaultSubject struct{ vaultSubjectEmbedded }

// TestVaultSuite checks the suite of Vault and the two forms of its
// application.
func TestVaultSuite(t *testing.T) {
	t.Parallel()

	t.Run("VaultSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, classified.VaultSuite())
	})

	t.Run("RunVaultSuite", func(t *testing.T) {
		t.Parallel()
		var built *vaultSubject
		classified.RunVaultSuite(
			t,
			func() *vaultSubject {
				built = &vaultSubject{}
				return built
			},
			conform.With(classified.VaultCheck{
				ID:        "Vault/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.VaultHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Vault(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.VaultHandle) classified.Vault { return &vaultSubject{} },
			}),
			conform.Excuse(
				classified.VaultChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.Chain(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.ChainGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Replay.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Replay.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Replay.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.VerifiesWhatItAppended(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.Lease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.LeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.OneHolderAtATime(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.Pool(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.CountsEachLoan(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunVaultSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *vaultSubject
		classified.RunVaultSuiteFromHandle(
			t,
			func(*classified.VaultHandle) *vaultSubject {
				built = &vaultSubject{}
				return built
			},
			conform.With(classified.VaultCheck{
				ID:        "Vault/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.VaultHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Vault(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.VaultHandle) classified.Vault { return &vaultSubject{} },
			}),
			conform.Excuse(
				classified.VaultChecks.Append.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.Chain(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Append.ChainGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Replay.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Replay.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Replay.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.VerifiesWhatItAppended(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.Lease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.LeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.OneHolderAtATime(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Release.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Borrow.Pool(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Return.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				classified.VaultChecks.Stats.CountsEachLoan(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// vaultIdle returns the planted defect of a ceiling of Vault with a spend
// of nothing, a subject that answers each call at once.
func vaultIdle(h *classified.VaultHandle) classified.Vault {
	return classified.VaultSuite().Calls[0].Over(h, func() {})
}

// TestVaultSuiteBudgets applies the allocation ceilings of the suite of
// Vault to vaultIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestVaultSuiteBudgets(t *testing.T) {
	classified.RunVaultBudgets(
		t,
		vaultIdle,
		conform.Ceilings(
			classified.VaultChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Replay.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Verify.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Acquire.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Release.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Borrow.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Return.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Stats.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteAppend registers the fuzz target of Vault.Append on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteAppend(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Append.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteReplay registers the fuzz target of Vault.Replay on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteReplay(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Replay.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteVerify registers the fuzz target of Vault.Verify on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteVerify(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Verify.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteAcquire registers the fuzz target of Vault.Acquire on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteAcquire(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Acquire.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteRelease registers the fuzz target of Vault.Release on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteRelease(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Release.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteBorrow registers the fuzz target of Vault.Borrow on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteBorrow(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Borrow.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteReturn registers the fuzz target of Vault.Return on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteReturn(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Return.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzVaultSuiteStats registers the fuzz target of Vault.Stats on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteStats(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Stats.Callable,
		vaultIdle,
		conform.Excuse(
			classified.VaultChecks.Append.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.Chain(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Append.ChainGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Replay.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppended(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.Lease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.LeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLease(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATime(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Acquire.OneHolderAtATimeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Release.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Borrow.Pool(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Return.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			classified.VaultChecks.Stats.CountsEachLoan(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkVaultSuite measures each callable of the suite of Vault on
// vaultIdle.
func BenchmarkVaultSuite(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Append.Callable, vaultIdle)
	})
	b.Run("Replay", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Replay.Callable, vaultIdle)
	})
	b.Run("Verify", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Verify.Callable, vaultIdle)
	})
	b.Run("Acquire", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Acquire.Callable, vaultIdle)
	})
	b.Run("Release", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Release.Callable, vaultIdle)
	})
	b.Run("Borrow", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Borrow.Callable, vaultIdle)
	})
	b.Run("Return", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Return.Callable, vaultIdle)
	})
	b.Run("Stats", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Stats.Callable, vaultIdle)
	})
}
