// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package roledtypestest contains the double and the suite of Store. InMemory
// is the control of the suite.
package roledtypestest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/roledtypes"
)

// ErrNotFound is the error of Get for a key without a payload.
var ErrNotFound = errors.New("roledtypestest: not found")

// InMemory implements [roledtypes.Store] in memory. The corpus applies the
// suite of Store to it as a control.
//
//dokimi:subject roledtypes.Store new=NewInMemory
type InMemory struct {
	payloads map[roledtypes.Key]roledtypes.Payload
	mu       sync.Mutex
}

// The compiler checks that InMemory implements [roledtypes.Store].
var _ roledtypes.Store = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{payloads: map[roledtypes.Key]roledtypes.Payload{}}
}

// Put implements [roledtypes.Store].
func (s *InMemory) Put(ctx context.Context, key roledtypes.Key, payload roledtypes.Payload) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payloads[key] = payload
	return nil
}

// Get implements [roledtypes.Store].
func (s *InMemory) Get(ctx context.Context, key roledtypes.Key) (roledtypes.Payload, error) {
	if err := contextErr(ctx); err != nil {
		return roledtypes.Payload{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.payloads[key]
	if !ok {
		return roledtypes.Payload{}, ErrNotFound
	}
	return p, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("roledtypestest: nil context")
	}
	return ctx.Err()
}
