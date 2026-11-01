// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package wrappedviatest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package wrappedviatest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/wrappedvia"
)

// ErrNotFound is the error of Get for a key without a value. Get wraps it with
// the key.
var ErrNotFound = errors.New("wrappedviatest: not found")

// InMemory implements [wrappedvia.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject wrappedvia.Mixed new=NewInMemory
type InMemory struct {
	values map[string]wrappedvia.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [wrappedvia.Mixed].
var _ wrappedvia.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]wrappedvia.Value{}}
}

// Get implements [wrappedvia.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (wrappedvia.Value, error) {
	if err := contextErr(ctx); err != nil {
		return wrappedvia.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return wrappedvia.Value{}, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	return v, nil
}

// TracedGet implements [wrappedvia.Mixed]. It delegates to Get, and wraps the
// error of Get with the name of the call.
func (s *InMemory) TracedGet(ctx context.Context, key string) (wrappedvia.Value, error) {
	v, err := s.Get(ctx, key)
	if err != nil {
		return v, fmt.Errorf("traced get: %w", err)
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("wrappedviatest: nil context")
	}
	return ctx.Err()
}
