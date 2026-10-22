// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package readertest contains the double and the suite of Reader. InMemory is
// the control of the suite.
package readertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/reader"
)

// InMemory implements [reader.Reader] in memory. The corpus applies the
// suite of Reader to it as a control.
//
//dokimi:subject reader.Reader new=NewInMemory
type InMemory struct {
	values map[string]reader.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [reader.Reader].
var _ reader.Reader = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]reader.Value{}} }

// Get implements [reader.Reader].
func (s *InMemory) Get(ctx context.Context, key string) (reader.Value, error) {
	if err := contextErr(ctx); err != nil {
		return reader.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return reader.Value{}, reader.ErrNotFound
	}
	return v, nil
}

// Put stores v under its key. A test seeds the reader through it.
func (s *InMemory) Put(v reader.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("readertest: nil context")
	}
	return ctx.Err()
}
