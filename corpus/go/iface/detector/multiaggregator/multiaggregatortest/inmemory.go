// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package multiaggregatortest contains the double and the suite of
// MultiAggregator. InMemory is the subject of the suite.
package multiaggregatortest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/multiaggregator"
)

// InMemory implements [multiaggregator.MultiAggregator] in memory. The
// corpus applies the suite of MultiAggregator to it as a subject.
//
//dokimi:subject multiaggregator.MultiAggregator new=NewInMemory
type InMemory struct {
	items []int
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [multiaggregator.MultiAggregator].
var _ multiaggregator.MultiAggregator = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Stats implements [multiaggregator.MultiAggregator].
func (s *InMemory) Stats(ctx context.Context) (int, int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var sum int
	for _, n := range s.items {
		sum += n
	}
	return len(s.items), sum, nil
}

// Add adds n to the items that Stats counts and sums. A test seeds the
// aggregator through it.
func (s *InMemory) Add(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, n)
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("multiaggregatortest: nil context")
	}
	return ctx.Err()
}
