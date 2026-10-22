// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package validatestest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package validatestest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/validates"
)

// ErrNotFound is the error of Read for a key that no call of Store wrote.
var ErrNotFound = errors.New("validatestest: not found")

// ErrInvalid is the error of Validate and Store for a payload without a key.
var ErrInvalid = errors.New("validatestest: invalid payload")

// InMemory implements [validates.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject validates.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items map[string]validates.Payload
}

// The compiler checks that InMemory implements [validates.Mixed].
var _ validates.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]validates.Payload{}} }

// Store implements [validates.Mixed].
func (s *InMemory) Store(ctx context.Context, v validates.Payload) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if err := s.Validate(v); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[v.Key] = v
	return nil
}

// Validate implements [validates.Mixed].
func (*InMemory) Validate(v validates.Payload) error {
	if v.Key == "" {
		return ErrInvalid
	}
	return nil
}

// Read implements [validates.Mixed].
func (s *InMemory) Read(ctx context.Context, key string) (validates.Payload, error) {
	if err := contextErr(ctx); err != nil {
		return validates.Payload{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[key]
	if !ok {
		return validates.Payload{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("validatestest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the payloads of prior. A restart after a
// crash keeps every payload that Store wrote. Do not use prior after the call.
// The two stores share the payloads but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{items: prior.items}
}
