// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package conservativetest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package conservativetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/conservative"
)

// InMemory implements [conservative.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject conservative.Mixed new=NewInMemory
type InMemory struct {
	buckets map[string]int
	reserve int
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [conservative.Mixed].
var _ conservative.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{buckets: map[string]int{}} }

// Apply implements [conservative.Mixed].
func (s *InMemory) Apply(ctx context.Context, d conservative.Delta) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buckets[d.Key] += d.Amount
	s.reserve -= d.Amount
	return nil
}

// Total implements [conservative.Mixed].
func (s *InMemory) Total(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	total := s.reserve
	for _, amount := range s.buckets {
		total += amount
	}
	return total, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("conservativetest: nil context")
	}
	return ctx.Err()
}
