// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package associativetest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package associativetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/associative"
)

// InMemory implements [associative.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject associative.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	total int
}

// The compiler checks that InMemory implements [associative.Mixed].
var _ associative.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Apply implements [associative.Mixed].
func (s *InMemory) Apply(ctx context.Context, ds []associative.Delta) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	batch := 0
	for _, d := range ds {
		batch += d.Amount
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total += batch
	return nil
}

// Total implements [associative.Mixed].
func (s *InMemory) Total(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.total, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("associativetest: nil context")
	}
	return ctx.Err()
}
