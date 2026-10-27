// Package badbudget states the budgets that suite refuses.
package badbudget

import "example.com/shop/meter"

// Gauge states a budget that suite refuses on each method.
//
//dokimi:suite
type Gauge interface {
	// Empty states no ceiling.
	//
	//dokimi:budget
	Empty()
	// Negative states a count of allocations below 0.
	//
	//dokimi:budget allocs=-1
	Negative()
	// Shrinking states a count of bytes below 0.
	//
	//dokimi:budget bytes=-2
	Shrinking()
	// Fast states a p99 that is no duration.
	//
	//dokimi:budget p99=fast
	Fast()
	// Instant states a p99 of no time.
	//
	//dokimi:budget p99=0s
	Instant()
	// Steady states a mean that is no duration.
	//
	//dokimi:budget p99=1ms mean=steady
	Steady()
	// Average states a mean without a p99.
	//
	//dokimi:budget mean=1us
	Average()
	// Named names a callable on a method.
	//
	//dokimi:budget Gauge.Named allocs=1
	Named()
	// Twice states two budgets.
	//
	//dokimi:budget allocs=1
	//dokimi:budget allocs=2
	Twice()
	// hidden states a budget of a method that is not exported.
	//
	//dokimi:budget allocs=1
	hidden()
	// Cold states a warm-up below 0.
	//
	//dokimi:budget warmup=-1
	Cold()
	// Warm states a warm-up without a ceiling.
	//
	//dokimi:budget warmup=10
	Warm()
	// Block waits for another call and states a budget.
	//dokimi:mixin waits
	//dokimi:budget allocs=1
	Block()
}

// Patient states a budget of a callable that waits for another call.
//
//dokimi:subject Gauge new=NewGauge
//dokimi:budget Gauge.Block allocs=1
type Patient struct{ Gauge }

// Doze waits for another call and states a budget.
//
//dokimi:suite
//dokimi:mixin waits
//dokimi:budget allocs=1
func Doze() {}

// Rewarmed states a budget of a callable whose declaration states a warm-up.
//
//dokimi:subject Gauge new=NewGauge
//dokimi:budget Gauge.Warm p99=1ms
type Rewarmed struct{ Gauge }

// NewGauge returns a gauge of the specs.
func NewGauge() Gauge { return nil }

// Unsuited has no directive suite.
type Unsuited interface {
	// Get states a budget of a declaration without a suite.
	//
	//dokimi:budget allocs=1
	Get()
}

// Loose states a budget of a function without a suite.
//
//dokimi:budget allocs=1
func Loose() {}

// Unnamed states a budget on a subject without a callable.
//
//dokimi:subject meter.Meter new=New
//dokimi:budget allocs=1
type Unnamed struct{ meter.Meter }

// Lost states a budget of a callable that the declaration lacks.
//
//dokimi:subject meter.Meter new=New
//dokimi:budget Meter.Missing allocs=1
type Lost struct{ meter.Meter }

// Doubled states a budget of a callable whose declaration states one.
//
//dokimi:subject meter.Meter new=New
//dokimi:budget Meter.Count allocs=2
type Doubled struct{ meter.Meter }

// Repeated states two budgets of one callable.
//
//dokimi:subject meter.Meter new=New
//dokimi:budget Meter.Total p99=1ms
//dokimi:budget Meter.Total p99=2ms
type Repeated struct{ meter.Meter }

// Invalid states a ceiling that is no duration.
//
//dokimi:subject meter.Meter new=New
//dokimi:budget Meter.Total p99=fast
type Invalid struct{ meter.Meter }

// Unexcused excuses a check of Count that the suite of Meter does not have.
//
//dokimi:subject meter.Meter new=New
//dokimi:excuse Meter.Count/keeps-the-count reason="the meter loses each sample"
type Unexcused struct{ meter.Meter }

// Unbound states a budget on a struct that is no subject.
//
//dokimi:budget Meter.Total p99=1ms
type Unbound struct{ meter.Meter }

// New returns a subject of the specs.
func New() meter.Meter { return nil }
