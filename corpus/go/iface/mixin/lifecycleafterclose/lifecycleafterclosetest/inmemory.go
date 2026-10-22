// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package lifecycleafterclosetest contains the double and the suite of Mixed.
// InMemory is the control of the suite.
package lifecycleafterclosetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/lifecycleafterclose"
)

// InMemory implements [lifecycleafterclose.Mixed] in memory. The corpus
// applies the suite of Mixed to it as a control.
//
//dokimi:subject lifecycleafterclose.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	closed bool
	works  int
}

// The compiler checks that InMemory implements [lifecycleafterclose.Mixed].
var _ lifecycleafterclose.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Close implements [lifecycleafterclose.Mixed].
func (s *InMemory) Close(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

// Work implements [lifecycleafterclose.Mixed].
func (s *InMemory) Work(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return lifecycleafterclose.ErrClosed
	}
	s.works++
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("lifecycleafterclosetest: nil context")
	}
	return ctx.Err()
}
