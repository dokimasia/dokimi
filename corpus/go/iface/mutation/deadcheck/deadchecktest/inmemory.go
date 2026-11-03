// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deadchecktest contains the suite of Store. InMemory is the subject
// of the suite.
package deadchecktest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mutation/deadcheck"
)

// errNotFound is the error of Get for a key without a value.
var errNotFound = errors.New("deadchecktest: not found")

// InMemory implements [deadcheck.Store] in memory. The corpus applies the
// suite of Store to it as a subject.
//
//dokimi:subject deadcheck.Store new=NewInMemory
type InMemory struct {
	items  map[string]string
	mu     sync.Mutex
	closed bool
}

// The compiler checks that InMemory implements [deadcheck.Store].
var _ deadcheck.Store = (*InMemory)(nil)

// NewInMemory returns an empty InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Get implements [deadcheck.Store].
func (s *InMemory) Get(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[key]
	if !ok {
		return "", errNotFound
	}
	return v, nil
}

// Put implements [deadcheck.Store].
func (s *InMemory) Put(ctx context.Context, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
	return nil
}

// Close implements [deadcheck.Store]. A nil InMemory closes without an error.
// The methods do not read the mark of Close, so a check does not observe a
// change of Close.
func (s *InMemory) Close() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

// Len returns the number of the values of s. No check of the suite calls it.
func (s *InMemory) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("deadchecktest: nil context")
	}
	return ctx.Err()
}
