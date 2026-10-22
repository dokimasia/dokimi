// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package cursortest contains the double and the suite of Contract. InMemory is
// the control of the suite.
package cursortest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/cursor"
)

// InMemory implements [cursor.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject cursor.Contract new=NewInMemory
type InMemory struct {
	values []cursor.Value
	at     int
	mu     sync.Mutex
	closed bool
}

// The compiler checks that InMemory implements [cursor.Contract].
var _ cursor.Contract = (*InMemory)(nil)

// NewInMemory returns a cursor over two values.
func NewInMemory() *InMemory {
	return NewInMemoryOf(
		cursor.Value{Key: "a", Body: "first"},
		cursor.Value{Key: "b", Body: "second"},
	)
}

// NewInMemoryOf returns a cursor over the values.
func NewInMemoryOf(values ...cursor.Value) *InMemory { return &InMemory{values: values} }

// Next implements [cursor.Contract].
func (s *InMemory) Next(ctx context.Context) (cursor.Value, bool, error) {
	if err := contextErr(ctx); err != nil {
		return cursor.Value{}, false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return cursor.Value{}, false, cursor.ErrClosed
	}
	if s.at >= len(s.values) {
		return cursor.Value{}, false, nil
	}
	v := s.values[s.at]
	s.at++
	return v, true, nil
}

// Close implements [cursor.Contract].
func (s *InMemory) Close(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("cursortest: nil context")
	}
	return ctx.Err()
}
