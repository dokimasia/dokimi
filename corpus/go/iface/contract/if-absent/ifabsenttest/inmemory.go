// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ifabsenttest contains the double and the suite of Contract. InMemory
// is the subject of the suite.
package ifabsenttest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	ifabsent "go.dokimi.dev/dokimi/corpus/go/iface/contract/if-absent"
)

// ErrPresent is the error of Put for a key that the store has. It wraps
// [ifabsent.ErrExists].
var ErrPresent = fmt.Errorf("ifabsenttest: key already present: %w", ifabsent.ErrExists)

// InMemory implements [ifabsent.Contract] in memory. The corpus applies the
// suite of Contract to it as a subject.
//
//dokimi:subject ifabsent.Contract new=NewInMemory
type InMemory struct {
	values map[string]ifabsent.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [ifabsent.Contract].
var _ ifabsent.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]ifabsent.Value{}} }

// Put implements [ifabsent.Contract].
func (s *InMemory) Put(ctx context.Context, v ifabsent.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, present := s.values[v.Key]; present {
		return ErrPresent
	}
	s.values[v.Key] = v
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("ifabsenttest: nil context")
	}
	return ctx.Err()
}
