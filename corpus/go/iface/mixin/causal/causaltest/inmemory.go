// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package causaltest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package causaltest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/causal"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("causaltest: not found")

// InMemory implements [causal.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a control.
//
//dokimi:subject causal.Mixed new=NewInMemory
type InMemory struct {
	values map[string]causal.Value
	rev    int64
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [causal.Mixed].
var _ causal.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]causal.Value{}}
}

// Store implements [causal.Mixed].
func (s *InMemory) Store(ctx context.Context, v causal.Value) (causal.Value, error) {
	if err := contextErr(ctx); err != nil {
		return causal.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rev++
	v.Rev = s.rev
	s.values[v.Key] = v
	return v, nil
}

// Get implements [causal.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (causal.Value, error) {
	if err := contextErr(ctx); err != nil {
		return causal.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return causal.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("causaltest: nil context")
	}
	return ctx.Err()
}
