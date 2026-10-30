// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package perkey contains a control of the suite of Mixed: a counter that
// keeps a total for each key and sums the totals on each read.
package perkey

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/commutative"
)

// Counter implements [commutative.Mixed] with a total for each key of a delta.
// The corpus applies the suite of Mixed to it as a control of the counters.
//
//dokimi:subject commutative.Mixed new=New
type Counter struct {
	totals map[string]int
	mu     sync.Mutex
}

// The compiler checks that Counter implements [commutative.Mixed].
var _ commutative.Mixed = (*Counter)(nil)

// New returns a Counter with a total of 0.
func New() *Counter { return &Counter{totals: map[string]int{}} }

// Apply implements [commutative.Mixed].
func (c *Counter) Apply(ctx context.Context, d commutative.Delta) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totals[d.Key] += d.Amount
	return nil
}

// Total implements [commutative.Mixed].
func (c *Counter) Total(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	sum := 0
	for _, total := range c.totals {
		sum += total
	}
	return sum, nil
}
