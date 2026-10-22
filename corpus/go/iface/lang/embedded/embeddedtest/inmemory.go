// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package embeddedtest contains the double and the suite of Composed. InMemory
// is the control of the suite.
package embeddedtest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/embedded"
)

// ErrNotFound is the error of Get for an identity without a value.
var ErrNotFound = errors.New("embeddedtest: not found")

// InMemory implements [embedded.Composed] in memory. The corpus applies the
// suite of Composed to it as a control.
//
//dokimi:subject embedded.Composed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	items  map[string]string
	closed bool
}

// The compiler checks that InMemory implements each of the three interfaces.
var (
	_ embedded.Base     = (*InMemory)(nil)
	_ embedded.Closer   = (*InMemory)(nil)
	_ embedded.Composed = (*InMemory)(nil)
)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Get implements [embedded.Composed].
func (s *InMemory) Get(ctx context.Context, id string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[id]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// Put stores value under id. A test seeds the store through it, because the
// interface declares no writer.
func (s *InMemory) Put(id, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = value
}

// Ping implements [embedded.Base].
func (*InMemory) Ping(ctx context.Context) error { return contextErr(ctx) }

// Close implements [embedded.Closer]. A second call returns nil too.
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
		return errors.New("embeddedtest: nil context")
	}
	return ctx.Err()
}
