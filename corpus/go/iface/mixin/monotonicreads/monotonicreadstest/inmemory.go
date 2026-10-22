// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package monotonicreadstest contains the double and the suite of Mixed.
// InMemory is the control of the suite.
package monotonicreadstest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/monotonicreads"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("monotonicreadstest: not found")

// InMemory implements [monotonicreads.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a control.
//
//dokimi:subject monotonicreads.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	rev    int64
	values map[string]monotonicreads.Value
}

// The compiler checks that InMemory implements [monotonicreads.Mixed].
var _ monotonicreads.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]monotonicreads.Value{}}
}

// Store implements [monotonicreads.Mixed].
func (s *InMemory) Store(ctx context.Context, v monotonicreads.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rev++
	v.Rev = s.rev
	s.values[v.Key] = v
	return nil
}

// Get implements [monotonicreads.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (monotonicreads.Value, error) {
	if err := contextErr(ctx); err != nil {
		return monotonicreads.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return monotonicreads.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("monotonicreadstest: nil context")
	}
	return ctx.Err()
}
