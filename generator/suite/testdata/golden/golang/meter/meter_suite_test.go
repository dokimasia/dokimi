package meter_test

import (
	"testing"

	"example.com/shop/meter"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// meterSubjectEmbedded is meter.Meter under the name of the field that
// meterSubject embeds.
type meterSubjectEmbedded = meter.Meter

// meterSubject is a value of Meter that only TestMeterSuite builds.
type meterSubject struct{ meterSubjectEmbedded }

// TestMeterSuite checks the suite of Meter and the two forms of its
// application.
func TestMeterSuite(t *testing.T) {
	t.Parallel()

	t.Run("MeterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, meter.MeterSuite())
	})

	t.Run("RunMeterSuite", func(t *testing.T) {
		t.Parallel()
		var built *meterSubject
		meter.RunMeterSuite(
			t,
			func() *meterSubject {
				built = &meterSubject{}
				return built
			},
			conform.With(meter.MeterCheck{
				ID:        "Meter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *meter.MeterHandle) {
					got := h.Subject()
					assert.Equal(h, got, meter.Meter(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*meter.MeterHandle) meter.Meter { return &meterSubject{} },
			}),
			conform.Excuse(
				meter.MeterChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Record.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Record.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Merge.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Merge.NilSafe(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Total.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Total.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunMeterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *meterSubject
		meter.RunMeterSuiteFromHandle(
			t,
			func(*meter.MeterHandle) *meterSubject {
				built = &meterSubject{}
				return built
			},
			conform.With(meter.MeterCheck{
				ID:        "Meter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *meter.MeterHandle) {
					got := h.Subject()
					assert.Equal(h, got, meter.Meter(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*meter.MeterHandle) meter.Meter { return &meterSubject{} },
			}),
			conform.Excuse(
				meter.MeterChecks.Count.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Count.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Count.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Record.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Record.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Merge.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Merge.NilSafe(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Total.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.MeterChecks.Total.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("MeterChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(meter.MeterChecks.Count.AllocationBudget()),
			"Meter.Count/allocation-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(meter.MeterChecks.Count.ByteBudget()),
			"Meter.Count/byte-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(meter.MeterChecks.Count.LatencyBudget()),
			"Meter.Count/latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(meter.MeterChecks.Count.MeanLatencyBudget()),
			"Meter.Count/mean-latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(meter.MeterChecks.Record.LatencyBudget()),
			"Meter.Record/latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// TestMeterSuiteAllocs falsifies the allocation ceilings and the checks
// that count allocations of the suite of Meter.
func TestMeterSuiteAllocs(t *testing.T) {
	t.Run("MeterSuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, meter.MeterSuite())
	})
}

// meterIdle returns the planted defect of a ceiling of Meter with a spend
// of nothing, a subject that answers each call at once.
func meterIdle(h *meter.MeterHandle) meter.Meter {
	return meter.MeterSuite().Calls[0].Over(h, func() {})
}

// TestMeterSuiteBudgets applies the allocation ceilings of the suite of
// Meter to meterIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestMeterSuiteBudgets(t *testing.T) {
	meter.RunMeterBudgets(
		t,
		meterIdle,
		conform.Ceilings(meter.MeterChecks.Merge.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(meter.MeterChecks.Total.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			meter.MeterChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.NilSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMeterSuiteCount registers the fuzz target of Meter.Count on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteCount(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Count.Callable,
		meterIdle,
		conform.Excuse(
			meter.MeterChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.NilSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMeterSuiteRecord registers the fuzz target of Meter.Record on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteRecord(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Record.Callable,
		meterIdle,
		conform.Excuse(
			meter.MeterChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.NilSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMeterSuiteMerge registers the fuzz target of Meter.Merge on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteMerge(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Merge.Callable,
		meterIdle,
		conform.Excuse(
			meter.MeterChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.NilSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMeterSuiteTotal registers the fuzz target of Meter.Total on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteTotal(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Total.Callable,
		meterIdle,
		conform.Excuse(
			meter.MeterChecks.Count.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Count.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Record.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Merge.NilSafe(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			meter.MeterChecks.Total.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkMeterSuite measures each callable of the suite of Meter on
// meterIdle.
func BenchmarkMeterSuite(b *testing.B) {
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuite(b, meter.MeterChecks.Count.Callable, meterIdle)
	})
	b.Run("Record", func(b *testing.B) {
		meter.BenchMeterSuite(b, meter.MeterChecks.Record.Callable, meterIdle)
	})
	b.Run("Merge", func(b *testing.B) {
		meter.BenchMeterSuite(b, meter.MeterChecks.Merge.Callable, meterIdle)
	})
	b.Run("Total", func(b *testing.B) {
		meter.BenchMeterSuite(b, meter.MeterChecks.Total.Callable, meterIdle)
	})
}

// BenchmarkMeterSuiteParallel measures each callable of the suite of Meter
// that admits concurrent calls on meterIdle, on GOMAXPROCS goroutines at
// once.
func BenchmarkMeterSuiteParallel(b *testing.B) {
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuiteParallel(b, meter.MeterChecks.Count.Callable, meterIdle)
	})
}

// TestWeighSuite checks the suite of Weigh and the two forms of its
// application.
func TestWeighSuite(t *testing.T) {
	t.Parallel()

	t.Run("WeighSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, meter.WeighSuite())
	})

	t.Run("Weigh", func(t *testing.T) {
		t.Parallel()
		meter.RunWeighSuite(t, meter.Weigh)
	})

	t.Run("RunWeighSuite", func(t *testing.T) {
		t.Parallel()
		meter.RunWeighSuite(
			t,
			meter.Weigh,
			conform.With(meter.WeighCheck{
				ID:        "Weigh/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *meter.WeighHandle) {
					assert.Equal(h, h.Subject(), meter.WeighFunc(meter.Weigh), h.Statement())
				},
				Defect: func(*meter.WeighHandle) meter.WeighFunc { return nil },
			}),
			conform.Excuse(
				meter.WeighChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.WeighChecks.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunWeighSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		meter.RunWeighSuiteFromHandle(
			t,
			func(*meter.WeighHandle) meter.WeighFunc {
				return meter.Weigh
			},
			conform.With(meter.WeighCheck{
				ID:        "Weigh/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *meter.WeighHandle) {
					assert.Equal(h, h.Subject(), meter.WeighFunc(meter.Weigh), h.Statement())
				},
				Defect: func(*meter.WeighHandle) meter.WeighFunc { return nil },
			}),
			conform.Excuse(
				meter.WeighChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				meter.WeighChecks.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("WeighChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(meter.WeighChecks.AllocationBudget()),
			"Weigh/allocation-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(meter.WeighChecks.LatencyBudget()),
			"Weigh/latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// TestWeighSuiteAllocs falsifies the allocation ceilings and the checks
// that count allocations of the suite of Weigh.
func TestWeighSuiteAllocs(t *testing.T) {
	t.Run("WeighSuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, meter.WeighSuite())
	})
}

// TestWeighAllocs applies the allocation ceiling and the checks that count
// allocations of Weigh to Weigh.
func TestWeighAllocs(t *testing.T) {
	meter.RunWeighBudgets(t, meter.Weigh)
}

// BenchmarkWeigh measures Weigh.
func BenchmarkWeigh(b *testing.B) {
	meter.BenchWeighSuite(b, meter.Weigh)
}

// BenchmarkWeighParallel measures Weigh on GOMAXPROCS goroutines at once.
func BenchmarkWeighParallel(b *testing.B) {
	meter.BenchWeighSuiteParallel(b, meter.Weigh)
}

// FuzzWeigh fuzzes the property checks of Weigh.
func FuzzWeigh(f *testing.F) {
	meter.FuzzWeighSuite(f, meter.Weigh)
}
