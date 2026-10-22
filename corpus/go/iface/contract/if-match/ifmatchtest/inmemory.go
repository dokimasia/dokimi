// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ifmatchtest contains the double and the suite of Contract. InMemory
// is the control of the suite.
package ifmatchtest

import (
	"context"
	"errors"
	"sync"

	ifmatch "go.dokimi.dev/dokimi/corpus/go/iface/contract/if-match"
)

// ErrAbsent is the error of Match for a key without a value.
var ErrAbsent = errors.New("ifmatchtest: no value under that key")

// ErrMismatch is the error of Put for a value whose body differs from the
// stored body of its key.
var ErrMismatch = errors.New("ifmatchtest: stored value does not match")

// InMemory implements [ifmatch.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject ifmatch.Contract new=NewInMemory
type InMemory struct {
	values map[string]ifmatch.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [ifmatch.Contract].
var _ ifmatch.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]ifmatch.Value{}} }

// Put implements [ifmatch.Contract].
func (s *InMemory) Put(ctx context.Context, v ifmatch.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if current, present := s.values[v.Key]; present && !matches(current, v) {
		return ErrMismatch
	}
	s.values[v.Key] = v
	return nil
}

// Match implements [ifmatch.Contract].
func (s *InMemory) Match(ctx context.Context, v ifmatch.Value) (bool, error) {
	if err := contextErr(ctx); err != nil {
		return false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, present := s.values[v.Key]
	if !present {
		return false, ErrAbsent
	}
	return matches(current, v), nil
}

// matches reports whether current has the body of want. Put and Match share
// it, so both compare a pair the same way.
func matches(current, want ifmatch.Value) bool { return current.Body == want.Body }

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("ifmatchtest: nil context")
	}
	return ctx.Err()
}
