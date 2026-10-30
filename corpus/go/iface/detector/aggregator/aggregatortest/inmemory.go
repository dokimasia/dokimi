// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package aggregatortest contains the double and the suite of Aggregator.
// InMemory is the subject of the suite.
package aggregatortest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/aggregator"
)

// InMemory implements [aggregator.Aggregator] in memory. The corpus applies
// the suite of Aggregator to it as a subject.
//
//dokimi:subject aggregator.Aggregator new=NewInMemory
type InMemory struct {
	items []string
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [aggregator.Aggregator].
var _ aggregator.Aggregator = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Count implements [aggregator.Aggregator].
func (s *InMemory) Count(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items), nil
}

// Add adds item to the items that Count counts. A test seeds the aggregator
// through it.
func (s *InMemory) Add(item string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("aggregatortest: nil context")
	}
	return ctx.Err()
}
