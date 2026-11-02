package faults_test

import (
	"testing"

	"example.com/shop/faults"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// cacheSubjectEmbedded is faults.Cache under the name of the field that
// cacheSubject embeds.
type cacheSubjectEmbedded = faults.Cache

// cacheSubject is a value of Cache that only TestCacheSuite builds.
type cacheSubject struct{ cacheSubjectEmbedded }

// TestCacheSuite checks the suite of Cache and the two forms of its
// application.
func TestCacheSuite(t *testing.T) {
	t.Parallel()

	t.Run("CacheSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.CacheSuite())
	})

	t.Run("RunCacheSuite", func(t *testing.T) {
		t.Parallel()
		var built *cacheSubject
		faults.RunCacheSuite(
			t,
			func() *cacheSubject {
				built = &cacheSubject{}
				return built
			},
			conform.With(faults.CacheCheck{
				ID:        "Cache/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.CacheHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Cache(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.CacheHandle) faults.Cache { return &cacheSubject{} },
			}),
			conform.Excuse(
				faults.CacheChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Len.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Len.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.ServesTheStoredValue(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCacheSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *cacheSubject
		faults.RunCacheSuiteFromHandle(
			t,
			func(*faults.CacheHandle) *cacheSubject {
				built = &cacheSubject{}
				return built
			},
			conform.With(faults.CacheCheck{
				ID:        "Cache/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.CacheHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Cache(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.CacheHandle) faults.Cache { return &cacheSubject{} },
			}),
			conform.Excuse(
				faults.CacheChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Len.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Len.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CacheChecks.Get.ServesTheStoredValue(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("CacheChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(faults.CacheChecks.Get.ServesTheStoredValue()),
			"Cache.Get/serves-the-stored-value",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// cacheIdle returns the planted defect of a ceiling of Cache with a spend
// of nothing, a subject that answers each call at once.
func cacheIdle(h *faults.CacheHandle) faults.Cache {
	return faults.CacheSuite().Calls[0].Over(h, func() {})
}

// TestCacheSuiteBudgets applies the allocation ceilings of the suite of
// Cache to cacheIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestCacheSuiteBudgets(t *testing.T) {
	faults.RunCacheBudgets(
		t,
		cacheIdle,
		conform.Ceilings(faults.CacheChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(faults.CacheChecks.Len.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			faults.CacheChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Len.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ServesTheStoredValue(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCacheSuiteGet registers the fuzz target of Cache.Get on cacheIdle,
// with every check of the suite excused.
func FuzzCacheSuiteGet(f *testing.F) {
	faults.FuzzCacheSuite(
		f,
		faults.CacheChecks.Get.Callable,
		cacheIdle,
		conform.Excuse(
			faults.CacheChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Len.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ServesTheStoredValue(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCacheSuiteLen registers the fuzz target of Cache.Len on cacheIdle,
// with every check of the suite excused.
func FuzzCacheSuiteLen(f *testing.F) {
	faults.FuzzCacheSuite(
		f,
		faults.CacheChecks.Len.Callable,
		cacheIdle,
		conform.Excuse(
			faults.CacheChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Len.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Len.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CacheChecks.Get.ServesTheStoredValue(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCacheSuite measures each callable of the suite of Cache on
// cacheIdle.
func BenchmarkCacheSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		faults.BenchCacheSuite(b, faults.CacheChecks.Get.Callable, cacheIdle)
	})
	b.Run("Len", func(b *testing.B) {
		faults.BenchCacheSuite(b, faults.CacheChecks.Len.Callable, cacheIdle)
	})
}

// writerSubjectEmbedded is faults.Writer under the name of the field that
// writerSubject embeds.
type writerSubjectEmbedded = faults.Writer

// writerSubject is a value of Writer that only TestWriterSuite builds.
type writerSubject struct{ writerSubjectEmbedded }

// TestWriterSuite checks the suite of Writer and the two forms of its
// application.
func TestWriterSuite(t *testing.T) {
	t.Parallel()

	t.Run("WriterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.WriterSuite())
	})

	t.Run("RunWriterSuite", func(t *testing.T) {
		t.Parallel()
		var built *writerSubject
		faults.RunWriterSuite(
			t,
			func() *writerSubject {
				built = &writerSubject{}
				return built
			},
			conform.With(faults.WriterCheck{
				ID:        "Writer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.WriterHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Writer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.WriterHandle) faults.Writer { return &writerSubject{} },
			}),
			conform.Excuse(
				faults.WriterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.NoPartialWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.NoPartialWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.NoTraceOfAFailedCall(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunWriterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *writerSubject
		faults.RunWriterSuiteFromHandle(
			t,
			func(*faults.WriterHandle) *writerSubject {
				built = &writerSubject{}
				return built
			},
			conform.With(faults.WriterCheck{
				ID:        "Writer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.WriterHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Writer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.WriterHandle) faults.Writer { return &writerSubject{} },
			}),
			conform.Excuse(
				faults.WriterChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.NoPartialWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.NoPartialWriteGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Put.NoTraceOfAFailedCall(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.WriterChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// writerIdle returns the planted defect of a ceiling of Writer with a
// spend of nothing, a subject that answers each call at once.
func writerIdle(h *faults.WriterHandle) faults.Writer {
	return faults.WriterSuite().Calls[0].Over(h, func() {})
}

// TestWriterSuiteBudgets applies the allocation ceilings of the suite of
// Writer to writerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestWriterSuiteBudgets(t *testing.T) {
	faults.RunWriterBudgets(
		t,
		writerIdle,
		conform.Ceilings(faults.WriterChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(faults.WriterChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			faults.WriterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoPartialWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoPartialWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoTraceOfAFailedCall(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWriterSuitePut registers the fuzz target of Writer.Put on
// writerIdle, with every check of the suite excused.
func FuzzWriterSuitePut(f *testing.F) {
	faults.FuzzWriterSuite(
		f,
		faults.WriterChecks.Put.Callable,
		writerIdle,
		conform.Excuse(
			faults.WriterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoPartialWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoPartialWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoTraceOfAFailedCall(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzWriterSuiteGet registers the fuzz target of Writer.Get on
// writerIdle, with every check of the suite excused.
func FuzzWriterSuiteGet(f *testing.F) {
	faults.FuzzWriterSuite(
		f,
		faults.WriterChecks.Get.Callable,
		writerIdle,
		conform.Excuse(
			faults.WriterChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoPartialWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoPartialWriteGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Put.NoTraceOfAFailedCall(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.WriterChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkWriterSuite measures each callable of the suite of Writer on
// writerIdle.
func BenchmarkWriterSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		faults.BenchWriterSuite(b, faults.WriterChecks.Put.Callable, writerIdle)
	})
	b.Run("Get", func(b *testing.B) {
		faults.BenchWriterSuite(b, faults.WriterChecks.Get.Callable, writerIdle)
	})
}

// syncerSubjectEmbedded is faults.Syncer under the name of the field that
// syncerSubject embeds.
type syncerSubjectEmbedded = faults.Syncer

// syncerSubject is a value of Syncer that only TestSyncerSuite builds.
type syncerSubject struct{ syncerSubjectEmbedded }

// TestSyncerSuite checks the suite of Syncer and the two forms of its
// application.
func TestSyncerSuite(t *testing.T) {
	t.Parallel()

	t.Run("SyncerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.SyncerSuite())
	})

	t.Run("RunSyncerSuite", func(t *testing.T) {
		t.Parallel()
		var built *syncerSubject
		faults.RunSyncerSuite(
			t,
			func() *syncerSubject {
				built = &syncerSubject{}
				return built
			},
			conform.With(faults.SyncerCheck{
				ID:        "Syncer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.SyncerHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Syncer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.SyncerHandle) faults.Syncer { return &syncerSubject{} },
			}),
			conform.Excuse(
				faults.SyncerChecks.Sync.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.RetrySucceedsUnderFaults(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSyncerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *syncerSubject
		faults.RunSyncerSuiteFromHandle(
			t,
			func(*faults.SyncerHandle) *syncerSubject {
				built = &syncerSubject{}
				return built
			},
			conform.With(faults.SyncerCheck{
				ID:        "Syncer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.SyncerHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Syncer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.SyncerHandle) faults.Syncer { return &syncerSubject{} },
			}),
			conform.Excuse(
				faults.SyncerChecks.Sync.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.SyncerChecks.Sync.RetrySucceedsUnderFaults(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// syncerIdle returns the planted defect of a ceiling of Syncer with a
// spend of nothing, a subject that answers each call at once.
func syncerIdle(h *faults.SyncerHandle) faults.Syncer {
	return faults.SyncerSuite().Calls[0].Over(h, func() {})
}

// TestSyncerSuiteBudgets applies the allocation ceilings of the suite of
// Syncer to syncerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestSyncerSuiteBudgets(t *testing.T) {
	faults.RunSyncerBudgets(
		t,
		syncerIdle,
		conform.Ceilings(faults.SyncerChecks.Sync.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			faults.SyncerChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.RetrySucceedsUnderFaults(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSyncerSuiteSync registers the fuzz target of Syncer.Sync on
// syncerIdle, with every check of the suite excused.
func FuzzSyncerSuiteSync(f *testing.F) {
	faults.FuzzSyncerSuite(
		f,
		faults.SyncerChecks.Sync.Callable,
		syncerIdle,
		conform.Excuse(
			faults.SyncerChecks.Sync.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.RetrySucceeds(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.SyncerChecks.Sync.RetrySucceedsUnderFaults(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSyncerSuite measures each callable of the suite of Syncer on
// syncerIdle.
func BenchmarkSyncerSuite(b *testing.B) {
	b.Run("Sync", func(b *testing.B) {
		faults.BenchSyncerSuite(b, faults.SyncerChecks.Sync.Callable, syncerIdle)
	})
}

// accountSubjectEmbedded is faults.Account under the name of the field
// that accountSubject embeds.
type accountSubjectEmbedded = faults.Account

// accountSubject is a value of Account that only TestAccountSuite builds.
type accountSubject struct{ accountSubjectEmbedded }

// TestAccountSuite checks the suite of Account and the two forms of its
// application.
func TestAccountSuite(t *testing.T) {
	t.Parallel()

	t.Run("AccountSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.AccountSuite())
	})

	t.Run("RunAccountSuite", func(t *testing.T) {
		t.Parallel()
		var built *accountSubject
		faults.RunAccountSuite(
			t,
			func() *accountSubject {
				built = &accountSubject{}
				return built
			},
			conform.With(faults.AccountCheck{
				ID:        "Account/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.AccountHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Account(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.AccountHandle) faults.Account { return &accountSubject{} },
			}),
			conform.Excuse(
				faults.AccountChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.Transaction(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.TransactionGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.TransactionUnderFaults(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.TransactionUnderFaultsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunAccountSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *accountSubject
		faults.RunAccountSuiteFromHandle(
			t,
			func(*faults.AccountHandle) *accountSubject {
				built = &accountSubject{}
				return built
			},
			conform.With(faults.AccountCheck{
				ID:        "Account/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.AccountHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Account(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.AccountHandle) faults.Account { return &accountSubject{} },
			}),
			conform.Excuse(
				faults.AccountChecks.Run.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.Transaction(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.TransactionGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.TransactionUnderFaults(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Run.TransactionUnderFaultsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Put.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.AccountChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// accountIdle returns the planted defect of a ceiling of Account with a
// spend of nothing, a subject that answers each call at once.
func accountIdle(h *faults.AccountHandle) faults.Account {
	return faults.AccountSuite().Calls[0].Over(h, func() {})
}

// TestAccountSuiteBudgets applies the allocation ceilings of the suite of
// Account to accountIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestAccountSuiteBudgets(t *testing.T) {
	faults.RunAccountBudgets(
		t,
		accountIdle,
		conform.Ceilings(faults.AccountChecks.Run.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(faults.AccountChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(faults.AccountChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			faults.AccountChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaults(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaultsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountSuiteRun registers the fuzz target of Account.Run on
// accountIdle, with every check of the suite excused.
func FuzzAccountSuiteRun(f *testing.F) {
	faults.FuzzAccountSuite(
		f,
		faults.AccountChecks.Run.Callable,
		accountIdle,
		conform.Excuse(
			faults.AccountChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaults(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaultsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountSuitePut registers the fuzz target of Account.Put on
// accountIdle, with every check of the suite excused.
func FuzzAccountSuitePut(f *testing.F) {
	faults.FuzzAccountSuite(
		f,
		faults.AccountChecks.Put.Callable,
		accountIdle,
		conform.Excuse(
			faults.AccountChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaults(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaultsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAccountSuiteGet registers the fuzz target of Account.Get on
// accountIdle, with every check of the suite excused.
func FuzzAccountSuiteGet(f *testing.F) {
	faults.FuzzAccountSuite(
		f,
		faults.AccountChecks.Get.Callable,
		accountIdle,
		conform.Excuse(
			faults.AccountChecks.Run.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.Transaction(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaults(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Run.TransactionUnderFaultsGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Put.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.AccountChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkAccountSuite measures each callable of the suite of Account on
// accountIdle.
func BenchmarkAccountSuite(b *testing.B) {
	b.Run("Run", func(b *testing.B) {
		faults.BenchAccountSuite(b, faults.AccountChecks.Run.Callable, accountIdle)
	})
	b.Run("Put", func(b *testing.B) {
		faults.BenchAccountSuite(b, faults.AccountChecks.Put.Callable, accountIdle)
	})
	b.Run("Get", func(b *testing.B) {
		faults.BenchAccountSuite(b, faults.AccountChecks.Get.Callable, accountIdle)
	})
}

// breakerSubjectEmbedded is faults.Breaker under the name of the field
// that breakerSubject embeds.
type breakerSubjectEmbedded = faults.Breaker

// breakerSubject is a value of Breaker that only TestBreakerSuite builds.
type breakerSubject struct{ breakerSubjectEmbedded }

// TestBreakerSuite checks the suite of Breaker and the two forms of its
// application.
func TestBreakerSuite(t *testing.T) {
	t.Parallel()

	t.Run("BreakerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.BreakerSuite())
	})

	t.Run("RunBreakerSuite", func(t *testing.T) {
		t.Parallel()
		var built *breakerSubject
		faults.RunBreakerSuite(
			t,
			func() *breakerSubject {
				built = &breakerSubject{}
				return built
			},
			conform.With(faults.BreakerCheck{
				ID:        "Breaker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.BreakerHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Breaker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.BreakerHandle) faults.Breaker { return &breakerSubject{} },
			}),
			conform.Excuse(
				faults.BreakerChecks.Do.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BreakerChecks.Do.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BreakerChecks.Do.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BreakerChecks.Do.FailsFast(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBreakerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *breakerSubject
		faults.RunBreakerSuiteFromHandle(
			t,
			func(*faults.BreakerHandle) *breakerSubject {
				built = &breakerSubject{}
				return built
			},
			conform.With(faults.BreakerCheck{
				ID:        "Breaker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.BreakerHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Breaker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.BreakerHandle) faults.Breaker { return &breakerSubject{} },
			}),
			conform.Excuse(
				faults.BreakerChecks.Do.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BreakerChecks.Do.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BreakerChecks.Do.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BreakerChecks.Do.FailsFast(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// breakerIdle returns the planted defect of a ceiling of Breaker with a
// spend of nothing, a subject that answers each call at once.
func breakerIdle(h *faults.BreakerHandle) faults.Breaker {
	return faults.BreakerSuite().Calls[0].Over(h, func() {})
}

// TestBreakerSuiteBudgets applies the allocation ceilings of the suite of
// Breaker to breakerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestBreakerSuiteBudgets(t *testing.T) {
	faults.RunBreakerBudgets(
		t,
		breakerIdle,
		conform.Ceilings(faults.BreakerChecks.Do.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			faults.BreakerChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BreakerChecks.Do.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BreakerChecks.Do.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BreakerChecks.Do.FailsFast(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBreakerSuiteDo registers the fuzz target of Breaker.Do on
// breakerIdle, with every check of the suite excused.
func FuzzBreakerSuiteDo(f *testing.F) {
	faults.FuzzBreakerSuite(
		f,
		faults.BreakerChecks.Do.Callable,
		breakerIdle,
		conform.Excuse(
			faults.BreakerChecks.Do.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BreakerChecks.Do.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BreakerChecks.Do.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BreakerChecks.Do.FailsFast(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBreakerSuite measures each callable of the suite of Breaker on
// breakerIdle.
func BenchmarkBreakerSuite(b *testing.B) {
	b.Run("Do", func(b *testing.B) {
		faults.BenchBreakerSuite(b, faults.BreakerChecks.Do.Callable, breakerIdle)
	})
}

// balancerSubjectEmbedded is faults.Balancer under the name of the field
// that balancerSubject embeds.
type balancerSubjectEmbedded = faults.Balancer

// balancerSubject is a value of Balancer that only TestBalancerSuite
// builds.
type balancerSubject struct{ balancerSubjectEmbedded }

// TestBalancerSuite checks the suite of Balancer and the two forms of its
// application.
func TestBalancerSuite(t *testing.T) {
	t.Parallel()

	t.Run("BalancerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.BalancerSuite())
	})

	t.Run("RunBalancerSuite", func(t *testing.T) {
		t.Parallel()
		var built *balancerSubject
		faults.RunBalancerSuite(
			t,
			func() *balancerSubject {
				built = &balancerSubject{}
				return built
			},
			conform.With(faults.BalancerCheck{
				ID:        "Balancer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.BalancerHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Balancer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.BalancerHandle) faults.Balancer { return &balancerSubject{} },
			}),
			conform.Excuse(
				faults.BalancerChecks.Route.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.KeepsItsInstance(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.KeepsItsInstanceGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBalancerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *balancerSubject
		faults.RunBalancerSuiteFromHandle(
			t,
			func(*faults.BalancerHandle) *balancerSubject {
				built = &balancerSubject{}
				return built
			},
			conform.With(faults.BalancerCheck{
				ID:        "Balancer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.BalancerHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Balancer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.BalancerHandle) faults.Balancer { return &balancerSubject{} },
			}),
			conform.Excuse(
				faults.BalancerChecks.Route.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.KeepsItsInstance(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.BalancerChecks.Route.KeepsItsInstanceGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// balancerIdle returns the planted defect of a ceiling of Balancer with a
// spend of nothing, a subject that answers each call at once.
func balancerIdle(h *faults.BalancerHandle) faults.Balancer {
	return faults.BalancerSuite().Calls[0].Over(h, func() {})
}

// TestBalancerSuiteBudgets applies the allocation ceilings of the suite of
// Balancer to balancerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestBalancerSuiteBudgets(t *testing.T) {
	faults.RunBalancerBudgets(
		t,
		balancerIdle,
		conform.Ceilings(
			faults.BalancerChecks.Route.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			faults.BalancerChecks.Route.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.KeepsItsInstance(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.KeepsItsInstanceGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBalancerSuiteRoute registers the fuzz target of Balancer.Route on
// balancerIdle, with every check of the suite excused.
func FuzzBalancerSuiteRoute(f *testing.F) {
	faults.FuzzBalancerSuite(
		f,
		faults.BalancerChecks.Route.Callable,
		balancerIdle,
		conform.Excuse(
			faults.BalancerChecks.Route.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.KeepsItsInstance(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.BalancerChecks.Route.KeepsItsInstanceGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBalancerSuite measures each callable of the suite of Balancer
// on balancerIdle.
func BenchmarkBalancerSuite(b *testing.B) {
	b.Run("Route", func(b *testing.B) {
		faults.BenchBalancerSuite(b, faults.BalancerChecks.Route.Callable, balancerIdle)
	})
}

// copierSubjectEmbedded is faults.Copier under the name of the field that
// copierSubject embeds.
type copierSubjectEmbedded = faults.Copier

// copierSubject is a value of Copier that only TestCopierSuite builds.
type copierSubject struct{ copierSubjectEmbedded }

// TestCopierSuite checks the suite of Copier and the two forms of its
// application.
func TestCopierSuite(t *testing.T) {
	t.Parallel()

	t.Run("CopierSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.CopierSuite())
	})

	t.Run("RunCopierSuite", func(t *testing.T) {
		t.Parallel()
		var built *copierSubject
		faults.RunCopierSuite(
			t,
			func() *copierSubject {
				built = &copierSubject{}
				return built
			},
			conform.With(faults.CopierCheck{
				ID:        "Copier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.CopierHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Copier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.CopierHandle) faults.Copier { return &copierSubject{} },
			}),
			conform.Excuse(
				faults.CopierChecks.Copy.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.ReleasesOnEveryPath(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.ReleasesOnEveryPathGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Open.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Close.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Close.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Close.FailureIsReported(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCopierSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *copierSubject
		faults.RunCopierSuiteFromHandle(
			t,
			func(*faults.CopierHandle) *copierSubject {
				built = &copierSubject{}
				return built
			},
			conform.With(faults.CopierCheck{
				ID:        "Copier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.CopierHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Copier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.CopierHandle) faults.Copier { return &copierSubject{} },
			}),
			conform.Excuse(
				faults.CopierChecks.Copy.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.ReleasesOnEveryPath(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Copy.ReleasesOnEveryPathGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Open.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Open.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Close.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Close.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.CopierChecks.Close.FailureIsReported(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// copierIdle returns the planted defect of a ceiling of Copier with a
// spend of nothing, a subject that answers each call at once.
func copierIdle(h *faults.CopierHandle) faults.Copier {
	return faults.CopierSuite().Calls[0].Over(h, func() {})
}

// TestCopierSuiteBudgets applies the allocation ceilings of the suite of
// Copier to copierIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCopierSuiteBudgets(t *testing.T) {
	faults.RunCopierBudgets(
		t,
		copierIdle,
		conform.Ceilings(faults.CopierChecks.Copy.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(faults.CopierChecks.Open.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			faults.CopierChecks.Close.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			faults.CopierChecks.Copy.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPath(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPathGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCopierSuiteCopy registers the fuzz target of Copier.Copy on
// copierIdle, with every check of the suite excused.
func FuzzCopierSuiteCopy(f *testing.F) {
	faults.FuzzCopierSuite(
		f,
		faults.CopierChecks.Copy.Callable,
		copierIdle,
		conform.Excuse(
			faults.CopierChecks.Copy.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPath(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPathGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCopierSuiteOpen registers the fuzz target of Copier.Open on
// copierIdle, with every check of the suite excused.
func FuzzCopierSuiteOpen(f *testing.F) {
	faults.FuzzCopierSuite(
		f,
		faults.CopierChecks.Open.Callable,
		copierIdle,
		conform.Excuse(
			faults.CopierChecks.Copy.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPath(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPathGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCopierSuiteClose registers the fuzz target of Copier.Close on
// copierIdle, with every check of the suite excused.
func FuzzCopierSuiteClose(f *testing.F) {
	faults.FuzzCopierSuite(
		f,
		faults.CopierChecks.Close.Callable,
		copierIdle,
		conform.Excuse(
			faults.CopierChecks.Copy.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPath(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Copy.ReleasesOnEveryPathGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Open.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.CopierChecks.Close.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCopierSuite measures each callable of the suite of Copier on
// copierIdle.
func BenchmarkCopierSuite(b *testing.B) {
	b.Run("Copy", func(b *testing.B) {
		faults.BenchCopierSuite(b, faults.CopierChecks.Copy.Callable, copierIdle)
	})
	b.Run("Open", func(b *testing.B) {
		faults.BenchCopierSuite(b, faults.CopierChecks.Open.Callable, copierIdle)
	})
	b.Run("Close", func(b *testing.B) {
		faults.BenchCopierSuite(b, faults.CopierChecks.Close.Callable, copierIdle)
	})
}

// undeclaredSubjectEmbedded is faults.Undeclared under the name of the
// field that undeclaredSubject embeds.
type undeclaredSubjectEmbedded = faults.Undeclared

// undeclaredSubject is a value of Undeclared that only TestUndeclaredSuite
// builds.
type undeclaredSubject struct{ undeclaredSubjectEmbedded }

// TestUndeclaredSuite checks the suite of Undeclared and the two forms of
// its application.
func TestUndeclaredSuite(t *testing.T) {
	t.Parallel()

	t.Run("UndeclaredSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.UndeclaredSuite())
	})

	t.Run("RunUndeclaredSuite", func(t *testing.T) {
		t.Parallel()
		var built *undeclaredSubject
		faults.RunUndeclaredSuite(
			t,
			func() *undeclaredSubject {
				built = &undeclaredSubject{}
				return built
			},
			conform.With(faults.UndeclaredCheck{
				ID:        "Undeclared/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.UndeclaredHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						faults.Undeclared(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*faults.UndeclaredHandle) faults.Undeclared { return &undeclaredSubject{} },
			}),
			conform.Excuse(
				faults.UndeclaredChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunUndeclaredSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *undeclaredSubject
		faults.RunUndeclaredSuiteFromHandle(
			t,
			func(*faults.UndeclaredHandle) *undeclaredSubject {
				built = &undeclaredSubject{}
				return built
			},
			conform.With(faults.UndeclaredCheck{
				ID:        "Undeclared/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.UndeclaredHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						faults.Undeclared(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*faults.UndeclaredHandle) faults.Undeclared { return &undeclaredSubject{} },
			}),
			conform.Excuse(
				faults.UndeclaredChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Put.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Put.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.UndeclaredChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// undeclaredIdle returns the planted defect of a ceiling of Undeclared
// with a spend of nothing, a subject that answers each call at once.
func undeclaredIdle(h *faults.UndeclaredHandle) faults.Undeclared {
	return faults.UndeclaredSuite().Calls[0].Over(h, func() {})
}

// TestUndeclaredSuiteBudgets applies the allocation ceilings of the suite
// of Undeclared to undeclaredIdle with every check of the suite excused,
// and a ceiling of 0 to each callable without a budget whose calls build
// no input of their own.
func TestUndeclaredSuiteBudgets(t *testing.T) {
	faults.RunUndeclaredBudgets(
		t,
		undeclaredIdle,
		conform.Ceilings(
			faults.UndeclaredChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			faults.UndeclaredChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzUndeclaredSuitePut registers the fuzz target of Undeclared.Put on
// undeclaredIdle, with every check of the suite excused.
func FuzzUndeclaredSuitePut(f *testing.F) {
	faults.FuzzUndeclaredSuite(
		f,
		faults.UndeclaredChecks.Put.Callable,
		undeclaredIdle,
		conform.Excuse(
			faults.UndeclaredChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzUndeclaredSuiteGet registers the fuzz target of Undeclared.Get on
// undeclaredIdle, with every check of the suite excused.
func FuzzUndeclaredSuiteGet(f *testing.F) {
	faults.FuzzUndeclaredSuite(
		f,
		faults.UndeclaredChecks.Get.Callable,
		undeclaredIdle,
		conform.Excuse(
			faults.UndeclaredChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Put.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.UndeclaredChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkUndeclaredSuite measures each callable of the suite of
// Undeclared on undeclaredIdle.
func BenchmarkUndeclaredSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		faults.BenchUndeclaredSuite(b, faults.UndeclaredChecks.Put.Callable, undeclaredIdle)
	})
	b.Run("Get", func(b *testing.B) {
		faults.BenchUndeclaredSuite(b, faults.UndeclaredChecks.Get.Callable, undeclaredIdle)
	})
}

// gravitySubjectEmbedded is faults.Gravity under the name of the field
// that gravitySubject embeds.
type gravitySubjectEmbedded = faults.Gravity

// gravitySubject is a value of Gravity that only TestGravitySuite builds.
type gravitySubject struct{ gravitySubjectEmbedded }

// TestGravitySuite checks the suite of Gravity and the two forms of its
// application.
func TestGravitySuite(t *testing.T) {
	t.Parallel()

	t.Run("GravitySuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.GravitySuite())
	})

	t.Run("RunGravitySuite", func(t *testing.T) {
		t.Parallel()
		var built *gravitySubject
		faults.RunGravitySuite(
			t,
			func() *gravitySubject {
				built = &gravitySubject{}
				return built
			},
			conform.With(faults.GravityCheck{
				ID:        "Gravity/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.GravityHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Gravity(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.GravityHandle) faults.Gravity { return &gravitySubject{} },
			}),
			conform.Excuse(
				faults.GravityChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunGravitySuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *gravitySubject
		faults.RunGravitySuiteFromHandle(
			t,
			func(*faults.GravityHandle) *gravitySubject {
				built = &gravitySubject{}
				return built
			},
			conform.With(faults.GravityCheck{
				ID:        "Gravity/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.GravityHandle) {
					got := h.Subject()
					assert.Equal(h, got, faults.Gravity(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*faults.GravityHandle) faults.Gravity { return &gravitySubject{} },
			}),
			conform.Excuse(
				faults.GravityChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.GravityChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// gravityIdle returns the planted defect of a ceiling of Gravity with a
// spend of nothing, a subject that answers each call at once.
func gravityIdle(h *faults.GravityHandle) faults.Gravity {
	return faults.GravitySuite().Calls[0].Over(h, func() {})
}

// TestGravitySuiteBudgets applies the allocation ceilings of the suite of
// Gravity to gravityIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestGravitySuiteBudgets(t *testing.T) {
	faults.RunGravityBudgets(
		t,
		gravityIdle,
		conform.Ceilings(faults.GravityChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			faults.GravityChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzGravitySuiteGet registers the fuzz target of Gravity.Get on
// gravityIdle, with every check of the suite excused.
func FuzzGravitySuiteGet(f *testing.F) {
	faults.FuzzGravitySuite(
		f,
		faults.GravityChecks.Get.Callable,
		gravityIdle,
		conform.Excuse(
			faults.GravityChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.GravityChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkGravitySuite measures each callable of the suite of Gravity on
// gravityIdle.
func BenchmarkGravitySuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		faults.BenchGravitySuite(b, faults.GravityChecks.Get.Callable, gravityIdle)
	})
}

// mirrorSubjectEmbedded is faults.MirrorMethods under the name of the
// field that mirrorSubject embeds.
type mirrorSubjectEmbedded = faults.MirrorMethods

// mirrorSubject is a value of Mirror that only TestMirrorSuite builds.
type mirrorSubject struct{ mirrorSubjectEmbedded }

// TestMirrorSuite checks the suite of Mirror and the two forms of its
// application.
func TestMirrorSuite(t *testing.T) {
	t.Parallel()

	t.Run("MirrorSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.MirrorSuite())
	})

	t.Run("RunMirrorSuite", func(t *testing.T) {
		t.Parallel()
		var built *mirrorSubject
		faults.RunMirrorSuite(
			t,
			func() *mirrorSubject {
				built = &mirrorSubject{}
				return built
			},
			conform.With(faults.MirrorCheck{
				ID:        "Mirror/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.MirrorHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						faults.MirrorMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*faults.MirrorHandle) faults.MirrorMethods { return &mirrorSubject{} },
			}),
			conform.Excuse(
				faults.MirrorChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunMirrorSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *mirrorSubject
		faults.RunMirrorSuiteFromHandle(
			t,
			func(*faults.MirrorHandle) *mirrorSubject {
				built = &mirrorSubject{}
				return built
			},
			conform.With(faults.MirrorCheck{
				ID:        "Mirror/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *faults.MirrorHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						faults.MirrorMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*faults.MirrorHandle) faults.MirrorMethods { return &mirrorSubject{} },
			}),
			conform.Excuse(
				faults.MirrorChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.MirrorChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// mirrorIdle returns the planted defect of a ceiling of Mirror with a
// spend of nothing, a subject that answers each call at once.
func mirrorIdle(h *faults.MirrorHandle) faults.MirrorMethods {
	return faults.MirrorSuite().Calls[0].Over(h, func() {})
}

// TestMirrorSuiteBudgets applies the allocation ceilings of the suite of
// Mirror to mirrorIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestMirrorSuiteBudgets(t *testing.T) {
	faults.RunMirrorBudgets(
		t,
		mirrorIdle,
		conform.Ceilings(faults.MirrorChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			faults.MirrorChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMirrorSuiteGet registers the fuzz target of Mirror.Get on
// mirrorIdle, with every check of the suite excused.
func FuzzMirrorSuiteGet(f *testing.F) {
	faults.FuzzMirrorSuite(
		f,
		faults.MirrorChecks.Get.Callable,
		mirrorIdle,
		conform.Excuse(
			faults.MirrorChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.FailureIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.FailureIsReportedGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			faults.MirrorChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkMirrorSuite measures each callable of the suite of Mirror on
// mirrorIdle.
func BenchmarkMirrorSuite(b *testing.B) {
	b.Run("Get", func(b *testing.B) {
		faults.BenchMirrorSuite(b, faults.MirrorChecks.Get.Callable, mirrorIdle)
	})
}

// TestFetchSuite checks the suite of Fetch and the two forms of its
// application.
func TestFetchSuite(t *testing.T) {
	t.Parallel()

	t.Run("FetchSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.FetchSuite())
	})

	t.Run("Fetch", func(t *testing.T) {
		t.Parallel()
		faults.RunFetchSuite(t, faults.Fetch)
	})

	t.Run("RunFetchSuite", func(t *testing.T) {
		t.Parallel()
		faults.RunFetchSuite(
			t,
			faults.Fetch,
			conform.With(faults.FetchCheck{
				ID:        "Fetch/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *faults.FetchHandle) {
					assert.Equal(h, h.Subject(), faults.FetchFunc(faults.Fetch), h.Statement())
				},
				Defect: func(*faults.FetchHandle) faults.FetchFunc { return nil },
			}),
			conform.Excuse(
				faults.FetchChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceedsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceedsUnderFaults(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceedsUnderFaultsGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunFetchSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		faults.RunFetchSuiteFromHandle(
			t,
			func(*faults.FetchHandle) faults.FetchFunc {
				return faults.Fetch
			},
			conform.With(faults.FetchCheck{
				ID:        "Fetch/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *faults.FetchHandle) {
					assert.Equal(h, h.Subject(), faults.FetchFunc(faults.Fetch), h.Statement())
				},
				Defect: func(*faults.FetchHandle) faults.FetchFunc { return nil },
			}),
			conform.Excuse(
				faults.FetchChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceeds(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceedsGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceedsUnderFaults(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.FetchChecks.RetrySucceedsUnderFaultsGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// BenchmarkFetch measures Fetch.
func BenchmarkFetch(b *testing.B) {
	faults.BenchFetchSuite(b, faults.Fetch)
}

// FuzzFetch fuzzes the property checks of Fetch.
func FuzzFetch(f *testing.F) {
	faults.FuzzFetchSuite(f, faults.Fetch)
}

// TestPickSuite checks the suite of Pick and the two forms of its
// application.
func TestPickSuite(t *testing.T) {
	t.Parallel()

	t.Run("PickSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.PickSuite())
	})

	t.Run("Pick", func(t *testing.T) {
		t.Parallel()
		faults.RunPickSuite(t, faults.Pick)
	})

	t.Run("RunPickSuite", func(t *testing.T) {
		t.Parallel()
		faults.RunPickSuite(
			t,
			faults.Pick,
			conform.With(faults.PickCheck{
				ID:        "Pick/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *faults.PickHandle) {
					assert.Equal(h, h.Subject(), faults.PickFunc(faults.Pick), h.Statement())
				},
				Defect: func(*faults.PickHandle) faults.PickFunc { return nil },
			}),
			conform.Excuse(
				faults.PickChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.KeepsItsInstance(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.KeepsItsInstanceGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPickSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		faults.RunPickSuiteFromHandle(
			t,
			func(*faults.PickHandle) faults.PickFunc {
				return faults.Pick
			},
			conform.With(faults.PickCheck{
				ID:        "Pick/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *faults.PickHandle) {
					assert.Equal(h, h.Subject(), faults.PickFunc(faults.Pick), h.Statement())
				},
				Defect: func(*faults.PickHandle) faults.PickFunc { return nil },
			}),
			conform.Excuse(
				faults.PickChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.KeepsItsInstance(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.PickChecks.KeepsItsInstanceGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// BenchmarkPick measures Pick.
func BenchmarkPick(b *testing.B) {
	faults.BenchPickSuite(b, faults.Pick)
}

// FuzzPick fuzzes the property checks of Pick.
func FuzzPick(f *testing.F) {
	faults.FuzzPickSuite(f, faults.Pick)
}

// TestLoadSuite checks the suite of Load and the two forms of its
// application.
func TestLoadSuite(t *testing.T) {
	t.Parallel()

	t.Run("LoadSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, faults.LoadSuite())
	})

	t.Run("Load", func(t *testing.T) {
		t.Parallel()
		faults.RunLoadSuite(t, faults.Load)
	})

	t.Run("RunLoadSuite", func(t *testing.T) {
		t.Parallel()
		faults.RunLoadSuite(
			t,
			faults.Load,
			conform.With(faults.LoadCheck{
				ID:        "Load/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *faults.LoadHandle) {
					assert.Equal(h, h.Subject(), faults.LoadFunc(faults.Load), h.Statement())
				},
				Defect: func(*faults.LoadHandle) faults.LoadFunc { return nil },
			}),
			conform.Excuse(
				faults.LoadChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ReleasesOnEveryPath(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ReleasesOnEveryPathGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLoadSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		faults.RunLoadSuiteFromHandle(
			t,
			func(*faults.LoadHandle) faults.LoadFunc {
				return faults.Load
			},
			conform.With(faults.LoadCheck{
				ID:        "Load/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *faults.LoadHandle) {
					assert.Equal(h, h.Subject(), faults.LoadFunc(faults.Load), h.Statement())
				},
				Defect: func(*faults.LoadHandle) faults.LoadFunc { return nil },
			}),
			conform.Excuse(
				faults.LoadChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.FailureIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.FailureIsReportedGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ReleasesOnEveryPath(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				faults.LoadChecks.ReleasesOnEveryPathGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// BenchmarkLoad measures Load.
func BenchmarkLoad(b *testing.B) {
	faults.BenchLoadSuite(b, faults.Load)
}

// FuzzLoad fuzzes the property checks of Load.
func FuzzLoad(f *testing.F) {
	faults.FuzzLoadSuite(f, faults.Load)
}
