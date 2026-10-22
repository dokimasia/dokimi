package metered

import (
	"testing"
	"time"

	"go.dokimi.dev/conform"
)

// TestMeterConformsToMeter applies the checks of the suite of Meter to
// Meter.
func TestMeterConformsToMeter(t *testing.T) {
	t.Parallel()

	RunMeterSuite(t, New)
}

// TestMeterAllocs applies the allocation ceilings of the suite of Meter to
// Meter.
func TestMeterAllocs(t *testing.T) {
	RunMeterBudgets(
		t,
		func(*MeterHandle) *Meter { return New() },
		conform.Ceilings(
			MeterChecks.Merge.Callable,
			conform.Budget{
				Allocs:  new(uint64(0)),
				Latency: time.Millisecond,
			},
		),
	)
}

// BenchmarkMeter measures each callable of the suite of Meter on Meter.
func BenchmarkMeter(b *testing.B) {
	opts := []MeterOption{
		conform.Ceilings(
			MeterChecks.Merge.Callable,
			conform.Budget{
				Allocs:  new(uint64(0)),
				Latency: time.Millisecond,
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Count.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Record", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Record.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Merge", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Merge.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Total", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Total.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
}

// BenchmarkMeterParallel measures each callable of the suite of Meter that
// admits concurrent calls on Meter, on GOMAXPROCS goroutines at once.
func BenchmarkMeterParallel(b *testing.B) {
	opts := []MeterOption{
		conform.Ceilings(
			MeterChecks.Merge.Callable,
			conform.Budget{
				Allocs:  new(uint64(0)),
				Latency: time.Millisecond,
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuiteParallel(
			b,
			MeterChecks.Count.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
}

// FuzzMeterCount fuzzes the property checks of Meter.Count on Meter.
func FuzzMeterCount(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Count.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzMeterRecord fuzzes the property checks of Meter.Record on Meter.
func FuzzMeterRecord(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Record.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzMeterMerge fuzzes the property checks of Meter.Merge on Meter.
func FuzzMeterMerge(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Merge.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzMeterTotal fuzzes the property checks of Meter.Total on Meter.
func FuzzMeterTotal(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Total.Callable, func(*MeterHandle) *Meter { return New() })
}

// TestPlainConformsToMeter applies the checks of the suite of Meter to
// Plain.
func TestPlainConformsToMeter(t *testing.T) {
	t.Parallel()

	RunMeterSuite(
		t,
		New,
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// TestPlainAllocs applies the allocation ceilings of the suite of Meter to
// Plain.
func TestPlainAllocs(t *testing.T) {
	RunMeterBudgets(
		t,
		func(*MeterHandle) *Meter { return New() },
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// BenchmarkPlain measures each callable of the suite of Meter on Plain.
func BenchmarkPlain(b *testing.B) {
	opts := []MeterOption{
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	}
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Count.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Record", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Record.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Merge", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Merge.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Total", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Total.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
}

// BenchmarkPlainParallel measures each callable of the suite of Meter that
// admits concurrent calls on Plain, on GOMAXPROCS goroutines at once.
func BenchmarkPlainParallel(b *testing.B) {
	opts := []MeterOption{
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	}
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuiteParallel(
			b,
			MeterChecks.Count.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
}

// FuzzPlainCount fuzzes the property checks of Meter.Count on Plain.
func FuzzPlainCount(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Count.Callable,
		func(*MeterHandle) *Meter { return New() },
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// FuzzPlainRecord fuzzes the property checks of Meter.Record on Plain.
func FuzzPlainRecord(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Record.Callable,
		func(*MeterHandle) *Meter { return New() },
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// FuzzPlainMerge fuzzes the property checks of Meter.Merge on Plain.
func FuzzPlainMerge(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Merge.Callable,
		func(*MeterHandle) *Meter { return New() },
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// FuzzPlainTotal fuzzes the property checks of Meter.Total on Plain.
func FuzzPlainTotal(f *testing.F) {
	FuzzMeterSuite(
		f,
		MeterChecks.Total.Callable,
		func(*MeterHandle) *Meter { return New() },
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// TestWarmedConformsToMeter applies the checks of the suite of Meter to
// Warmed.
func TestWarmedConformsToMeter(t *testing.T) {
	t.Parallel()

	RunMeterSuite(t, New)
}

// TestWarmedAllocs applies the allocation ceilings of the suite of Meter
// to Warmed.
func TestWarmedAllocs(t *testing.T) {
	RunMeterBudgets(
		t,
		func(*MeterHandle) *Meter { return New() },
		conform.Ceilings(
			MeterChecks.Merge.Callable,
			conform.Budget{
				Warmup: new(300),
			},
		),
	)
}

// BenchmarkWarmed measures each callable of the suite of Meter on Warmed.
func BenchmarkWarmed(b *testing.B) {
	opts := []MeterOption{
		conform.Ceilings(
			MeterChecks.Merge.Callable,
			conform.Budget{
				Warmup: new(300),
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Count.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Record", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Record.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Merge", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Merge.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
	b.Run("Total", func(b *testing.B) {
		BenchMeterSuite(
			b,
			MeterChecks.Total.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
}

// BenchmarkWarmedParallel measures each callable of the suite of Meter
// that admits concurrent calls on Warmed, on GOMAXPROCS goroutines at
// once.
func BenchmarkWarmedParallel(b *testing.B) {
	opts := []MeterOption{
		conform.Ceilings(
			MeterChecks.Merge.Callable,
			conform.Budget{
				Warmup: new(300),
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		BenchMeterSuiteParallel(
			b,
			MeterChecks.Count.Callable,
			func(*MeterHandle) *Meter { return New() },
			opts...)
	})
}

// FuzzWarmedCount fuzzes the property checks of Meter.Count on Warmed.
func FuzzWarmedCount(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Count.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzWarmedRecord fuzzes the property checks of Meter.Record on Warmed.
func FuzzWarmedRecord(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Record.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzWarmedMerge fuzzes the property checks of Meter.Merge on Warmed.
func FuzzWarmedMerge(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Merge.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzWarmedTotal fuzzes the property checks of Meter.Total on Warmed.
func FuzzWarmedTotal(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Total.Callable, func(*MeterHandle) *Meter { return New() })
}
