// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sagatest contains the double and the suite of Contract. InMemory is
// the control of the suite.
package sagatest

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/saga"
)

// ErrNotApplied is the error of a compensation of a step that the saga did not
// apply.
var ErrNotApplied = errors.New("sagatest: no applied step to compensate")

// ErrAlreadyApplied is the error of a step for a key that the saga applied.
var ErrAlreadyApplied = errors.New("sagatest: a step for that key is already applied")

// InMemory implements [saga.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject saga.Contract new=NewInMemory
type InMemory struct {
	// applied and shipped are the values of the two steps. Each step has a list
	// of its own, so a compensation of one step leaves the other step unchanged.
	applied []saga.Value
	shipped []saga.Value
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [saga.Contract].
var _ saga.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Step implements [saga.Contract].
func (s *InMemory) Step(ctx context.Context, v saga.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if slices.ContainsFunc(s.applied, func(a saga.Value) bool { return a.Key == v.Key }) {
		return ErrAlreadyApplied
	}
	s.applied = append(s.applied, v)
	return nil
}

// Compensate implements [saga.Contract].
func (s *InMemory) Compensate(ctx context.Context, v saga.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	at := slices.Index(s.applied, v)
	if at < 0 {
		return ErrNotApplied
	}
	s.applied = slices.Delete(s.applied, at, at+1)
	return nil
}

// Ship implements [saga.Contract].
func (s *InMemory) Ship(ctx context.Context, v saga.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if slices.ContainsFunc(s.shipped, func(a saga.Value) bool { return a.Key == v.Key }) {
		return ErrAlreadyApplied
	}
	s.shipped = append(s.shipped, v)
	return nil
}

// UndoShip implements [saga.Contract].
func (s *InMemory) UndoShip(ctx context.Context, v saga.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	at := slices.Index(s.shipped, v)
	if at < 0 {
		return ErrNotApplied
	}
	s.shipped = slices.Delete(s.shipped, at, at+1)
	return nil
}

// State implements [saga.Contract].
func (s *InMemory) State(ctx context.Context) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	parts := make([]string, 0, len(s.applied)+len(s.shipped))
	for _, v := range s.applied {
		parts = append(parts, "a"+v.Key+"\x00"+v.Body)
	}
	for _, v := range s.shipped {
		parts = append(parts, "s"+v.Key+"\x00"+v.Body)
	}
	return strings.Join(parts, "\x1f"), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("sagatest: nil context")
	}
	return ctx.Err()
}
