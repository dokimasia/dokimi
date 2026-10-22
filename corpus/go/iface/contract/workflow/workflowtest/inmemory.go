// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package workflowtest contains the double and the suite of Contract. InMemory
// is the control of the suite.
package workflowtest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/workflow"
)

// ErrTerminal is the error of Run for a key in a state without a transition.
var ErrTerminal = errors.New("workflowtest: no transition out of the last state")

// InMemory implements [workflow.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject workflow.Contract new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	state map[string]string
}

// The compiler checks that InMemory implements [workflow.Contract].
var _ workflow.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{state: map[string]string{}}
}

// Run implements [workflow.Contract].
func (s *InMemory) Run(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.state[key]
	out := workflow.Transitions[current]
	if len(out) == 0 {
		return fmt.Errorf("%w: %q is terminal", ErrTerminal, current)
	}
	s.state[key] = out[0]
	return nil
}

// State implements [workflow.Contract].
func (s *InMemory) State(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state[key], nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("workflowtest: nil context")
	}
	return ctx.Err()
}
