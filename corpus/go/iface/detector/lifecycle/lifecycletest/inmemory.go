// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package lifecycletest contains the double and the suite of Lifecycle.
// InMemory is the control of the suite.
package lifecycletest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/lifecycle"
)

// InMemory implements [lifecycle.Lifecycle] in memory. The corpus applies
// the suite of Lifecycle to it as a control.
//
//dokimi:subject lifecycle.Lifecycle new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	closed bool
}

// The compiler checks that InMemory implements [lifecycle.Lifecycle].
var _ lifecycle.Lifecycle = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Close implements [lifecycle.Lifecycle].
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
		return errors.New("lifecycletest: nil context")
	}
	return ctx.Err()
}
