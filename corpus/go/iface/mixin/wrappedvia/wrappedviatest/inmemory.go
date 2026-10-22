// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package wrappedviatest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package wrappedviatest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/wrappedvia"
)

// ErrUnderlying is the cause that the error of Open wraps. Cause returns it
// after that error.
var ErrUnderlying = errors.New("wrappedviatest: underlying failure")

// failing is the name that Open refuses.
const failing = "closed"

// InMemory implements [wrappedvia.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject wrappedvia.Mixed new=NewInMemory
type InMemory struct {
	last error
	mu   sync.Mutex
}

// The compiler checks that InMemory implements [wrappedvia.Mixed].
var _ wrappedvia.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Open implements [wrappedvia.Mixed].
func (s *InMemory) Open(ctx context.Context, name string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if name != failing {
		return nil
	}
	err := fmt.Errorf("opening %q: %w", name, ErrUnderlying)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = ErrUnderlying
	return err
}

// Cause implements [wrappedvia.Mixed].
func (s *InMemory) Cause(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last
}

// FailingName returns the name that Open refuses.
func FailingName() string { return failing }

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("wrappedviatest: nil context")
	}
	return ctx.Err()
}
