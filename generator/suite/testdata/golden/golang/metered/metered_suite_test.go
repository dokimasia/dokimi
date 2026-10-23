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

// TestMeterConformsToMeterAllocs applies the allocation ceilings of the
// suite of Meter to Meter.
func TestMeterConformsToMeterAllocs(t *testing.T) {
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

// BenchmarkMeterConformsToMeter measures each callable of the suite of
// Meter on Meter.
func BenchmarkMeterConformsToMeter(b *testing.B) {
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

// BenchmarkMeterConformsToMeterParallel measures each callable of the
// suite of Meter that admits concurrent calls on Meter, on GOMAXPROCS
// goroutines at once.
func BenchmarkMeterConformsToMeterParallel(b *testing.B) {
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

// FuzzMeterConformsToMeterCount fuzzes the property checks of Meter.Count
// on Meter.
func FuzzMeterConformsToMeterCount(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Count.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzMeterConformsToMeterRecord fuzzes the property checks of
// Meter.Record on Meter.
func FuzzMeterConformsToMeterRecord(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Record.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzMeterConformsToMeterMerge fuzzes the property checks of Meter.Merge
// on Meter.
func FuzzMeterConformsToMeterMerge(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Merge.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzMeterConformsToMeterTotal fuzzes the property checks of Meter.Total
// on Meter.
func FuzzMeterConformsToMeterTotal(f *testing.F) {
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

// TestPlainConformsToMeterAllocs applies the allocation ceilings of the
// suite of Meter to Plain.
func TestPlainConformsToMeterAllocs(t *testing.T) {
	RunMeterBudgets(
		t,
		func(*MeterHandle) *Meter { return New() },
		conform.Excuse(
			MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// BenchmarkPlainConformsToMeter measures each callable of the suite of
// Meter on Plain.
func BenchmarkPlainConformsToMeter(b *testing.B) {
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

// BenchmarkPlainConformsToMeterParallel measures each callable of the
// suite of Meter that admits concurrent calls on Plain, on GOMAXPROCS
// goroutines at once.
func BenchmarkPlainConformsToMeterParallel(b *testing.B) {
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

// FuzzPlainConformsToMeterCount fuzzes the property checks of Meter.Count
// on Plain.
func FuzzPlainConformsToMeterCount(f *testing.F) {
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

// FuzzPlainConformsToMeterRecord fuzzes the property checks of
// Meter.Record on Plain.
func FuzzPlainConformsToMeterRecord(f *testing.F) {
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

// FuzzPlainConformsToMeterMerge fuzzes the property checks of Meter.Merge
// on Plain.
func FuzzPlainConformsToMeterMerge(f *testing.F) {
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

// FuzzPlainConformsToMeterTotal fuzzes the property checks of Meter.Total
// on Plain.
func FuzzPlainConformsToMeterTotal(f *testing.F) {
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

// TestWarmedConformsToMeterAllocs applies the allocation ceilings of the
// suite of Meter to Warmed.
func TestWarmedConformsToMeterAllocs(t *testing.T) {
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

// BenchmarkWarmedConformsToMeter measures each callable of the suite of
// Meter on Warmed.
func BenchmarkWarmedConformsToMeter(b *testing.B) {
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

// BenchmarkWarmedConformsToMeterParallel measures each callable of the
// suite of Meter that admits concurrent calls on Warmed, on GOMAXPROCS
// goroutines at once.
func BenchmarkWarmedConformsToMeterParallel(b *testing.B) {
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

// FuzzWarmedConformsToMeterCount fuzzes the property checks of Meter.Count
// on Warmed.
func FuzzWarmedConformsToMeterCount(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Count.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzWarmedConformsToMeterRecord fuzzes the property checks of
// Meter.Record on Warmed.
func FuzzWarmedConformsToMeterRecord(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Record.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzWarmedConformsToMeterMerge fuzzes the property checks of Meter.Merge
// on Warmed.
func FuzzWarmedConformsToMeterMerge(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Merge.Callable, func(*MeterHandle) *Meter { return New() })
}

// FuzzWarmedConformsToMeterTotal fuzzes the property checks of Meter.Total
// on Warmed.
func FuzzWarmedConformsToMeterTotal(f *testing.F) {
	FuzzMeterSuite(f, MeterChecks.Total.Callable, func(*MeterHandle) *Meter { return New() })
}
