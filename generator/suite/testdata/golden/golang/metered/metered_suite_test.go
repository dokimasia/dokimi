package metered_test

import (
	"testing"
	"time"

	"example.com/shop/meter"
	"example.com/shop/metered"
	"go.dokimi.dev/conform"
)

// TestMeterConformsToMeter applies the checks of the suite of Meter to
// Meter.
func TestMeterConformsToMeter(t *testing.T) {
	t.Parallel()

	meter.RunMeterSuite(t, metered.New)
}

// TestMeterConformsToMeterAllocs applies the allocation ceilings and the
// checks that count allocations of the suite of Meter to Meter.
func TestMeterConformsToMeterAllocs(t *testing.T) {
	meter.RunMeterBudgets(
		t,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
		conform.Ceilings(
			meter.MeterChecks.Merge.Callable,
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
	opts := []meter.MeterOption{
		conform.Ceilings(
			meter.MeterChecks.Merge.Callable,
			conform.Budget{
				Allocs:  new(uint64(0)),
				Latency: time.Millisecond,
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Count.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Record", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Record.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Merge", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Merge.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Total", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Total.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
}

// BenchmarkMeterConformsToMeterParallel measures each callable of the
// suite of Meter that admits concurrent calls on Meter, on GOMAXPROCS
// goroutines at once.
func BenchmarkMeterConformsToMeterParallel(b *testing.B) {
	opts := []meter.MeterOption{
		conform.Ceilings(
			meter.MeterChecks.Merge.Callable,
			conform.Budget{
				Allocs:  new(uint64(0)),
				Latency: time.Millisecond,
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuiteParallel(
			b,
			meter.MeterChecks.Count.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
}

// FuzzMeterConformsToMeterCount fuzzes the property checks of Meter.Count
// on Meter.
func FuzzMeterConformsToMeterCount(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Count.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}

// FuzzMeterConformsToMeterRecord fuzzes the property checks of
// Meter.Record on Meter.
func FuzzMeterConformsToMeterRecord(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Record.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}

// FuzzMeterConformsToMeterMerge fuzzes the property checks of Meter.Merge
// on Meter.
func FuzzMeterConformsToMeterMerge(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Merge.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}

// FuzzMeterConformsToMeterTotal fuzzes the property checks of Meter.Total
// on Meter.
func FuzzMeterConformsToMeterTotal(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Total.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}

// TestPlainConformsToMeter applies the checks of the suite of Meter to
// Plain.
func TestPlainConformsToMeter(t *testing.T) {
	t.Parallel()

	meter.RunMeterSuite(
		t,
		metered.New,
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// TestPlainConformsToMeterAllocs applies the allocation ceilings and the
// checks that count allocations of the suite of Meter to Plain.
func TestPlainConformsToMeterAllocs(t *testing.T) {
	meter.RunMeterBudgets(
		t,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// BenchmarkPlainConformsToMeter measures each callable of the suite of
// Meter on Plain.
func BenchmarkPlainConformsToMeter(b *testing.B) {
	opts := []meter.MeterOption{
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	}
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Count.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Record", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Record.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Merge", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Merge.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Total", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Total.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
}

// BenchmarkPlainConformsToMeterParallel measures each callable of the
// suite of Meter that admits concurrent calls on Plain, on GOMAXPROCS
// goroutines at once.
func BenchmarkPlainConformsToMeterParallel(b *testing.B) {
	opts := []meter.MeterOption{
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	}
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuiteParallel(
			b,
			meter.MeterChecks.Count.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
}

// FuzzPlainConformsToMeterCount fuzzes the property checks of Meter.Count
// on Plain.
func FuzzPlainConformsToMeterCount(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Count.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// FuzzPlainConformsToMeterRecord fuzzes the property checks of
// Meter.Record on Plain.
func FuzzPlainConformsToMeterRecord(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Record.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// FuzzPlainConformsToMeterMerge fuzzes the property checks of Meter.Merge
// on Plain.
func FuzzPlainConformsToMeterMerge(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Merge.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// FuzzPlainConformsToMeterTotal fuzzes the property checks of Meter.Total
// on Plain.
func FuzzPlainConformsToMeterTotal(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Total.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
		conform.Excuse(
			meter.MeterChecks.Count.LatencyBudget(),
			"the meter scans every sample",
		),
	)
}

// TestWarmedConformsToMeter applies the checks of the suite of Meter to
// Warmed.
func TestWarmedConformsToMeter(t *testing.T) {
	t.Parallel()

	meter.RunMeterSuite(t, metered.New)
}

// TestWarmedConformsToMeterAllocs applies the allocation ceilings and the
// checks that count allocations of the suite of Meter to Warmed.
func TestWarmedConformsToMeterAllocs(t *testing.T) {
	meter.RunMeterBudgets(
		t,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
		conform.Ceilings(
			meter.MeterChecks.Merge.Callable,
			conform.Budget{
				Warmup: new(300),
			},
		),
	)
}

// BenchmarkWarmedConformsToMeter measures each callable of the suite of
// Meter on Warmed.
func BenchmarkWarmedConformsToMeter(b *testing.B) {
	opts := []meter.MeterOption{
		conform.Ceilings(
			meter.MeterChecks.Merge.Callable,
			conform.Budget{
				Warmup: new(300),
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Count.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Record", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Record.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Merge", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Merge.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
	b.Run("Total", func(b *testing.B) {
		meter.BenchMeterSuite(
			b,
			meter.MeterChecks.Total.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
}

// BenchmarkWarmedConformsToMeterParallel measures each callable of the
// suite of Meter that admits concurrent calls on Warmed, on GOMAXPROCS
// goroutines at once.
func BenchmarkWarmedConformsToMeterParallel(b *testing.B) {
	opts := []meter.MeterOption{
		conform.Ceilings(
			meter.MeterChecks.Merge.Callable,
			conform.Budget{
				Warmup: new(300),
			},
		),
	}
	b.Run("Count", func(b *testing.B) {
		meter.BenchMeterSuiteParallel(
			b,
			meter.MeterChecks.Count.Callable,
			func(*meter.MeterHandle) *metered.Meter { return metered.New() },
			opts...)
	})
}

// FuzzWarmedConformsToMeterCount fuzzes the property checks of Meter.Count
// on Warmed.
func FuzzWarmedConformsToMeterCount(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Count.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}

// FuzzWarmedConformsToMeterRecord fuzzes the property checks of
// Meter.Record on Warmed.
func FuzzWarmedConformsToMeterRecord(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Record.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}

// FuzzWarmedConformsToMeterMerge fuzzes the property checks of Meter.Merge
// on Warmed.
func FuzzWarmedConformsToMeterMerge(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Merge.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}

// FuzzWarmedConformsToMeterTotal fuzzes the property checks of Meter.Total
// on Warmed.
func FuzzWarmedConformsToMeterTotal(f *testing.F) {
	meter.FuzzMeterSuite(
		f,
		meter.MeterChecks.Total.Callable,
		func(*meter.MeterHandle) *metered.Meter { return metered.New() },
	)
}
