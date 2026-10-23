package meter

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// meterSubjectEmbedded is Meter under the name of the field that
// meterSubject embeds.
type meterSubjectEmbedded = Meter

// meterSubject is a value of Meter that only TestMeterSuite builds.
type meterSubject struct{ meterSubjectEmbedded }

// TestMeterSuite checks the suite of Meter and the two forms of its
// application.
func TestMeterSuite(t *testing.T) {
	t.Parallel()

	t.Run("MeterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, MeterSuite())
	})

	t.Run("RunMeterSuite", func(t *testing.T) {
		t.Parallel()
		var built *meterSubject
		RunMeterSuite(
			t,
			func() *meterSubject {
				built = &meterSubject{}
				return built
			},
			conform.With(MeterCheck{
				ID:        "Meter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *MeterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Meter(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*MeterHandle) Meter { return &meterSubject{} },
			}),
		)
	})

	t.Run("RunMeterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *meterSubject
		RunMeterSuiteFromHandle(
			t,
			func(*MeterHandle) *meterSubject {
				built = &meterSubject{}
				return built
			},
			conform.With(MeterCheck{
				ID:        "Meter/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *MeterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Meter(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*MeterHandle) Meter { return &meterSubject{} },
			}),
		)
	})

	t.Run("MeterChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(MeterChecks.Count.AllocationBudget()),
			"Meter.Count/allocation-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(MeterChecks.Count.ByteBudget()),
			"Meter.Count/byte-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(MeterChecks.Count.LatencyBudget()),
			"Meter.Count/latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(MeterChecks.Count.MeanLatencyBudget()),
			"Meter.Count/mean-latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(MeterChecks.Record.LatencyBudget()),
			"Meter.Record/latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// TestMeterSuiteAllocs falsifies the allocation ceilings of the suite of
// Meter.
func TestMeterSuiteAllocs(t *testing.T) {
	t.Run("MeterSuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, MeterSuite())
	})
}

// meterIdle returns the planted defect of a ceiling of Meter with a spend
// of nothing, a subject that answers each call at once.
func meterIdle(h *MeterHandle) Meter {
	return MeterSuite().Calls[0].Over(h, func() {})
}

// TestMeterSuiteBudgets applies the allocation ceilings of the suite of
// Meter to meterIdle, and a ceiling of 0 to each callable without a budget
// whose calls build no input of their own.
func TestMeterSuiteBudgets(t *testing.T) {
	RunMeterBudgets(
		t,
		meterIdle,
		conform.Ceilings(MeterChecks.Merge.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(MeterChecks.Total.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzMeterSuiteCount registers the fuzz target of Meter.Count on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteCount(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Count.Callable,
		meterIdle,
	)
}

// FuzzMeterSuiteRecord registers the fuzz target of Meter.Record on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteRecord(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Record.Callable,
		meterIdle,
	)
}

// FuzzMeterSuiteMerge registers the fuzz target of Meter.Merge on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteMerge(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Merge.Callable,
		meterIdle,
	)
}

// FuzzMeterSuiteTotal registers the fuzz target of Meter.Total on
// meterIdle, with every check of the suite excused.
func FuzzMeterSuiteTotal(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Total.Callable,
		meterIdle,
	)
}

// BenchmarkMeterSuite measures each callable of the suite of Meter on
// meterIdle.
func BenchmarkMeterSuite(b *testing.B) {
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuite(b, MeterChecks.Count.Callable, meterIdle)
	})
	b.Run("Record", func(b *testing.B) {
		BenchMeterSuite(b, MeterChecks.Record.Callable, meterIdle)
	})
	b.Run("Merge", func(b *testing.B) {
		BenchMeterSuite(b, MeterChecks.Merge.Callable, meterIdle)
	})
	b.Run("Total", func(b *testing.B) {
		BenchMeterSuite(b, MeterChecks.Total.Callable, meterIdle)
	})
}

// BenchmarkMeterSuiteParallel measures each callable of the suite of Meter
// that admits concurrent calls on meterIdle, on GOMAXPROCS goroutines at
// once.
func BenchmarkMeterSuiteParallel(b *testing.B) {
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuiteParallel(b, MeterChecks.Count.Callable, meterIdle)
	})
}

// TestWeighSuite checks the suite of Weigh and the two forms of its
// application.
func TestWeighSuite(t *testing.T) {
	t.Parallel()

	t.Run("WeighSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, WeighSuite())
	})

	t.Run("Weigh", func(t *testing.T) {
		t.Parallel()
		RunWeighSuite(t, Weigh)
	})

	t.Run("RunWeighSuite", func(t *testing.T) {
		t.Parallel()
		RunWeighSuite(
			t,
			Weigh,
			conform.With(WeighCheck{
				ID:        "Weigh/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *WeighHandle) {
					assert.Equal(h, h.Subject(), WeighFunc(Weigh), h.Statement())
				},
				Defect: func(*WeighHandle) WeighFunc { return nil },
			}),
		)
	})

	t.Run("RunWeighSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		RunWeighSuiteFromHandle(
			t,
			func(*WeighHandle) WeighFunc {
				return Weigh
			},
			conform.With(WeighCheck{
				ID:        "Weigh/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *WeighHandle) {
					assert.Equal(h, h.Subject(), WeighFunc(Weigh), h.Statement())
				},
				Defect: func(*WeighHandle) WeighFunc { return nil },
			}),
		)
	})

	t.Run("WeighChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(WeighChecks.AllocationBudget()),
			"Weigh/allocation-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(WeighChecks.LatencyBudget()),
			"Weigh/latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// TestWeighSuiteAllocs falsifies the allocation ceilings of the suite of
// Weigh.
func TestWeighSuiteAllocs(t *testing.T) {
	t.Run("WeighSuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, WeighSuite())
	})
}

// TestWeighAllocs applies the allocation ceiling of Weigh to Weigh.
func TestWeighAllocs(t *testing.T) {
	RunWeighBudgets(t, Weigh)
}

// BenchmarkWeigh measures Weigh.
func BenchmarkWeigh(b *testing.B) {
	BenchWeighSuite(b, Weigh)
}

// BenchmarkWeighParallel measures Weigh on GOMAXPROCS goroutines at once.
func BenchmarkWeighParallel(b *testing.B) {
	BenchWeighSuiteParallel(b, Weigh)
}

// FuzzWeigh fuzzes the property checks of Weigh.
func FuzzWeigh(f *testing.F) {
	FuzzWeighSuite(f, Weigh)
}
