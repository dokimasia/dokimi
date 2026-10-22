// Package metered measures the items of the shop in memory.
package metered

import (
	"context"

	"example.com/shop/meter"
)

// Meter measures the items of the shop in memory, under a budget of its own
// for Merge.
//
//dokimi:subject meter.Meter new=New
//dokimi:budget Meter.Merge allocs=0 p99=1ms
type Meter struct {
	samples []meter.Sample
}

// New returns an empty meter.
func New() *Meter {
	return &Meter{}
}

// Count returns the number of samples of the unit.
func (m *Meter) Count(_ context.Context, unit meter.Unit) (int, error) {
	n := 0
	for _, s := range m.samples {
		if s.Unit == unit {
			n++
		}
	}
	return n, nil
}

// Record stores the samples.
func (m *Meter) Record(_ context.Context, samples ...meter.Sample) error {
	m.samples = append(m.samples, samples...)
	return nil
}

// Merge folds the sample in, and takes nil for none.
func (m *Meter) Merge(s *meter.Sample) {
	if s != nil {
		m.samples = append(m.samples, *s)
	}
}

// Total returns the sum of the values of the samples.
func (m *Meter) Total() int {
	total := 0
	for _, s := range m.samples {
		total += s.Value
	}
	return total
}

// Plain measures the items of the shop in memory, under the budget of the
// declaration alone, and excuses the p99 ceiling of Count.
//
//dokimi:subject meter.Meter new=New
//dokimi:excuse Meter.Count/latency-budget reason="the meter scans every sample"
type Plain struct {
	*Meter
}

// Warmed measures the items of the shop in memory, with a warm-up of its own
// for the benchmarks of Merge.
//
//dokimi:subject meter.Meter new=New
//dokimi:budget Meter.Merge warmup=300
type Warmed struct {
	*Meter
}
