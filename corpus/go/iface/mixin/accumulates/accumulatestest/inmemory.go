// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package accumulatestest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package accumulatestest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/accumulates"
)

// InMemory implements [accumulates.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject accumulates.Mixed new=NewInMemory
type InMemory struct {
	totals map[string]int
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [accumulates.Mixed].
var _ accumulates.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{totals: map[string]int{}} }

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("accumulatestest: nil context")
	}
	return ctx.Err()
}

// Add implements [accumulates.Mixed].
func (s *InMemory) Add(ctx context.Context, key string, amount int) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totals[key] += amount
	return nil
}

// Total implements [accumulates.Mixed].
func (s *InMemory) Total(ctx context.Context, key string) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totals[key], nil
}
