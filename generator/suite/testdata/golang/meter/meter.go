// Package meter measures the items of the shop.
package meter

import "context"

// Unit is the unit of a measure.
type Unit int

// The units of a measure.
const (
	Grams Unit = iota + 1
	Litres
)

// Sample is one measure of an item.
type Sample struct {
	ID    string
	Value int
	Unit  Unit
}

// Meter measures the items of the shop.
//
//dokimi:suite
type Meter interface {
	// Count returns the number of samples of the unit.
	//
	//dokimi:budget allocs=1 bytes=96 p99=2us mean=1us
	//dokimi:mixin concurrent
	Count(ctx context.Context, unit Unit) (int, error)
	// Record stores the samples.
	//
	//dokimi:budget p99=1ms warmup=500
	Record(ctx context.Context, samples ...Sample) error
	// Merge folds the sample in, and takes nil for none.
	//
	//dokimi:mixin nil-safe
	Merge(s *Sample)
	// Total returns the sum of the values of the samples.
	Total() int
}

// Weigh returns the weight of the sample in grams.
//
//dokimi:suite
//dokimi:budget allocs=0 p99=100ns
//dokimi:mixin concurrent
func Weigh(s Sample, density float64) float64 {
	return float64(s.Value) * density
}
