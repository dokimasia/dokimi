// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package pointintimetest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package pointintimetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/pointintime"
)

// ErrNotFound is the error of Get for a key that no call of Store wrote.
var ErrNotFound = errors.New("pointintimetest: not found")

// InMemory implements [pointintime.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject pointintime.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	source map[string]pointintime.Value
	view   map[string]pointintime.Value
}

// The compiler checks that InMemory implements [pointintime.Mixed].
var _ pointintime.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		source: map[string]pointintime.Value{},
		view:   map[string]pointintime.Value{},
	}
}

// Store implements [pointintime.Mixed].
func (s *InMemory) Store(ctx context.Context, v pointintime.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.source[v.Key] = v
	return nil
}

// Get implements [pointintime.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (pointintime.Value, error) {
	if err := contextErr(ctx); err != nil {
		return pointintime.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if pinnedValue, pinned := s.view[key]; pinned {
		return pinnedValue, nil
	}
	v, ok := s.source[key]
	if !ok {
		return pointintime.Value{}, ErrNotFound
	}
	s.view[key] = v
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("pointintimetest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the source of prior and an empty view. A
// restart after a crash keeps the source and loses the view. Do not use prior
// after the call. The two stores share the source but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{source: prior.source, view: map[string]pointintime.Value{}}
}
