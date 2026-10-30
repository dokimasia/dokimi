// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package monotonicwritestest contains the double and the suite of Mixed.
// InMemory is the subject of the suite.
package monotonicwritestest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/monotonicwrites"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("monotonicwritestest: not found")

// InMemory implements [monotonicwrites.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a subject.
//
//dokimi:subject monotonicwrites.Mixed new=NewInMemory
type InMemory struct {
	values map[string]monotonicwrites.Value
	rev    int64
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [monotonicwrites.Mixed].
var _ monotonicwrites.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]monotonicwrites.Value{}}
}

// Store implements [monotonicwrites.Mixed].
func (s *InMemory) Store(ctx context.Context, v monotonicwrites.Value) (monotonicwrites.Value, error) {
	if err := contextErr(ctx); err != nil {
		return monotonicwrites.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rev++
	v.Rev = s.rev
	s.values[v.Key] = v
	return v, nil
}

// Get implements [monotonicwrites.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (monotonicwrites.Value, error) {
	if err := contextErr(ctx); err != nil {
		return monotonicwrites.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return monotonicwrites.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("monotonicwritestest: nil context")
	}
	return ctx.Err()
}
