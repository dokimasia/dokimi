// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package monotonictest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package monotonictest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/monotonic"
)

// InMemory implements [monotonic.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject monotonic.Mixed new=NewInMemory
type InMemory struct {
	mu      sync.Mutex
	version int64
}

// The compiler checks that InMemory implements [monotonic.Mixed].
var _ monotonic.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Version implements [monotonic.Mixed].
func (s *InMemory) Version(ctx context.Context) (int64, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.version, nil
}

// Advance implements [monotonic.Mixed].
func (s *InMemory) Advance(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.version++
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("monotonictest: nil context")
	}
	return ctx.Err()
}
