// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package atomictest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package atomictest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/atomic"
)

// ErrNotFound is the error of Read for a key without an entry.
var ErrNotFound = errors.New("atomictest: not found")

// ErrHalfEntry is the error of Write for an entry with one of its two halves.
// Write then stores no part of the entry.
var ErrHalfEntry = errors.New("atomictest: refuses an entry missing one half")

// InMemory implements [atomic.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a control.
//
//dokimi:subject atomic.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items map[string]atomic.Entry
}

// The compiler checks that InMemory implements [atomic.Mixed].
var _ atomic.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]atomic.Entry{}} }

// Write implements [atomic.Mixed].
func (s *InMemory) Write(ctx context.Context, e atomic.Entry) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if (e.Left == "") != (e.Right == "") {
		return ErrHalfEntry
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[e.Key] = e
	return nil
}

// Read implements [atomic.Mixed].
func (s *InMemory) Read(ctx context.Context, key string) (atomic.Entry, error) {
	if err := contextErr(ctx); err != nil {
		return atomic.Entry{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.items[key]
	if !ok {
		return atomic.Entry{}, ErrNotFound
	}
	return e, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("atomictest: nil context")
	}
	return ctx.Err()
}
