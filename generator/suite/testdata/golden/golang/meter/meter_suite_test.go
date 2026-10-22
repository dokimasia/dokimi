package meter

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// meterSubject is a value of Meter that only TestMeterSuite builds.
type meterSubject struct{ Meter }

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
}

// TestMeterSuiteAllocs falsifies the allocation ceilings of the suite of
// Meter.
func TestMeterSuiteAllocs(t *testing.T) {
	t.Run("MeterSuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, MeterSuite())
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
