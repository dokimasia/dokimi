// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deprecatedtest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package deprecatedtest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/deprecated"
)

// ErrNotFound is the error of Old and New for a key without a value.
var ErrNotFound = errors.New("deprecatedtest: not found")

// InMemory implements [deprecated.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject deprecated.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items map[string]string
}

// The compiler checks that InMemory implements [deprecated.Mixed].
var _ deprecated.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Old implements [deprecated.Mixed].
func (s *InMemory) Old(ctx context.Context, key string) (string, error) {
	return s.New(ctx, key)
}

// New implements [deprecated.Mixed].
func (s *InMemory) New(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// Put stores value under key. A test seeds the store through it, because the
// interface declares no writer.
func (s *InMemory) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("deprecatedtest: nil context")
	}
	return ctx.Err()
}
