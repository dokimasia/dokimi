// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package idempotentclosetest contains the double and the suite of Closer.
// InMemory is the subject of the suite.
package idempotentclosetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/idempotentclose"
)

// InMemory implements [idempotentclose.Closer] in memory. The corpus applies
// the suite of Closer to it as a subject.
//
//dokimi:subject idempotentclose.Closer new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	closed bool
}

// The compiler checks that InMemory implements [idempotentclose.Closer].
var _ idempotentclose.Closer = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Close implements [idempotentclose.Closer].
func (s *InMemory) Close(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

// Stats implements [idempotentclose.Closer].
func (s *InMemory) Stats(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0, nil
	}
	return 1, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("idempotentclosetest: nil context")
	}
	return ctx.Err()
}
