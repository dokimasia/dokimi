package buffered_test

import (
	"testing"

	"example.com/shop/buffered"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// recordsSubjectEmbedded is buffered.Records under the name of the field
// that recordsSubject embeds.
type recordsSubjectEmbedded = buffered.Records

// recordsSubject is a value of Records that only TestRecordsSuite builds.
type recordsSubject struct{ recordsSubjectEmbedded }

// TestRecordsSuite checks the suite of Records and the two forms of its
// application.
func TestRecordsSuite(t *testing.T) {
	t.Parallel()

	t.Run("RecordsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, buffered.RecordsSuite())
	})

	t.Run("RunRecordsSuite", func(t *testing.T) {
		t.Parallel()
		var built *recordsSubject
		buffered.RunRecordsSuite(
			t,
			func() *recordsSubject {
				built = &recordsSubject{}
				return built
			},
			conform.With(buffered.RecordsCheck{
				ID:        "Records/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *buffered.RecordsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						buffered.Records(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*buffered.RecordsHandle) buffered.Records { return &recordsSubject{} },
			}),
			conform.ExcuseSuite[buffered.Records]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunRecordsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *recordsSubject
		buffered.RunRecordsSuiteFromHandle(
			t,
			func(*buffered.RecordsHandle) *recordsSubject {
				built = &recordsSubject{}
				return built
			},
			conform.With(buffered.RecordsCheck{
				ID:        "Records/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *buffered.RecordsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						buffered.Records(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*buffered.RecordsHandle) buffered.Records { return &recordsSubject{} },
			}),
			conform.ExcuseSuite[buffered.Records]("the test applies the check of its option alone"),
		)
	})
}

// recordsSuite is the suite of Records that the tests of the planted
// defect of a ceiling apply.
var recordsSuite = buffered.RecordsSuite()

// recordsIdle returns the planted defect of a ceiling of Records with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func recordsIdle(h *buffered.RecordsHandle) buffered.Records {
	return recordsSuite.Calls[0].Over(h, func() {})
}

// TestRecordsSuiteBudgets applies the allocation ceilings of the suite of
// Records to recordsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRecordsSuiteBudgets(t *testing.T) {
	buffered.RunRecordsBudgets(
		t,
		recordsIdle,
		conform.Ceilings(
			buffered.RecordsChecks.Read.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			buffered.RecordsChecks.Fill.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			buffered.RecordsChecks.Peek.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[buffered.Records](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRecordsSuiteRead registers the fuzz target of Records.Read on
// recordsIdle, with every check of the suite excused.
func FuzzRecordsSuiteRead(f *testing.F) {
	buffered.FuzzRecordsSuite(
		f,
		buffered.RecordsChecks.Read.Callable,
		recordsIdle,
		conform.ExcuseSuite[buffered.Records](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRecordsSuiteFill registers the fuzz target of Records.Fill on
// recordsIdle, with every check of the suite excused.
func FuzzRecordsSuiteFill(f *testing.F) {
	buffered.FuzzRecordsSuite(
		f,
		buffered.RecordsChecks.Fill.Callable,
		recordsIdle,
		conform.ExcuseSuite[buffered.Records](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRecordsSuitePeek registers the fuzz target of Records.Peek on
// recordsIdle, with every check of the suite excused.
func FuzzRecordsSuitePeek(f *testing.F) {
	buffered.FuzzRecordsSuite(
		f,
		buffered.RecordsChecks.Peek.Callable,
		recordsIdle,
		conform.ExcuseSuite[buffered.Records](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkRecordsSuite measures each callable of the suite of Records on
// recordsIdle.
func BenchmarkRecordsSuite(b *testing.B) {
	b.Run("Read", func(b *testing.B) {
		buffered.BenchRecordsSuite(b, buffered.RecordsChecks.Read.Callable, recordsIdle)
	})
	b.Run("Fill", func(b *testing.B) {
		buffered.BenchRecordsSuite(b, buffered.RecordsChecks.Fill.Callable, recordsIdle)
	})
	b.Run("Peek", func(b *testing.B) {
		buffered.BenchRecordsSuite(b, buffered.RecordsChecks.Peek.Callable, recordsIdle)
	})
}
