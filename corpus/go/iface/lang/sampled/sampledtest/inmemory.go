// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sampledtest contains the double and the suite of Labeller. InMemory
// is the control of the suite.
package sampledtest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/sampled"
)

// InMemory implements [sampled.Labeller] in memory. The corpus applies the
// suite of Labeller to it as a control.
//
//dokimi:subject sampled.Labeller new=NewInMemory
type InMemory struct {
	mu      sync.Mutex
	known   map[sampled.Tag]bool
	applied map[sampled.Tag]sampled.Level
}

// The compiler checks that InMemory implements [sampled.Labeller].
var _ sampled.Labeller = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		known:   map[sampled.Tag]bool{sampled.NewTestTag(): true},
		applied: map[sampled.Tag]sampled.Level{},
	}
}

// Apply implements [sampled.Labeller].
func (s *InMemory) Apply(ctx context.Context, tag sampled.Tag, level sampled.Level) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.known[tag] {
		return sampled.ErrUnknownTag
	}
	s.applied[tag] = level
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("sampledtest: nil context")
	}
	return ctx.Err()
}
