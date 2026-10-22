// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package concurrenttest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package concurrenttest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/concurrent"
)

// InMemory implements [concurrent.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject concurrent.Mixed new=NewInMemory
type InMemory struct {
	counts map[string]int
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [concurrent.Mixed].
var _ concurrent.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{counts: map[string]int{}} }

// Bump implements [concurrent.Mixed].
func (s *InMemory) Bump(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[key]++
	return nil
}

// Count implements [concurrent.Mixed].
func (s *InMemory) Count(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var total int
	for _, n := range s.counts {
		total += n
	}
	return total, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("concurrenttest: nil context")
	}
	return ctx.Err()
}
