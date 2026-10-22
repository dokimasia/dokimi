// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sideeffecttest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package sideeffecttest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/sideeffect"
)

// InMemory implements [sideeffect.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject sideeffect.Mixed new=NewInMemory
type InMemory struct {
	mu      sync.Mutex
	touches map[string]int
}

// The compiler checks that InMemory implements [sideeffect.Mixed].
var _ sideeffect.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{touches: map[string]int{}} }

// Touch implements [sideeffect.Mixed].
func (s *InMemory) Touch(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.touches[key]++
	return nil
}

// Observed implements [sideeffect.Mixed].
func (s *InMemory) Observed(ctx context.Context, key string) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.touches[key], nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("sideeffecttest: nil context")
	}
	return ctx.Err()
}
