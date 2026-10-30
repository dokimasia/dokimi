// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package retrysucceedstest contains the double and the suite of Mixed.
// InMemory is the subject of the suite.
package retrysucceedstest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/retrysucceeds"
)

// ErrTransient is the error of each of the first failuresBeforeSuccess calls of
// a key.
var ErrTransient = errors.New("retrysucceedstest: transient failure")

// failuresBeforeSuccess is the number of calls of a key that fail before a
// call of the key succeeds.
const failuresBeforeSuccess = 2

// InMemory implements [retrysucceeds.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a subject.
//
//dokimi:subject retrysucceeds.Mixed new=NewInMemory
type InMemory struct {
	attempts map[string]int
	mu       sync.Mutex
}

// The compiler checks that InMemory implements [retrysucceeds.Mixed].
var _ retrysucceeds.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{attempts: map[string]int{}} }

// Call implements [retrysucceeds.Mixed].
func (s *InMemory) Call(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attempts[key]++
	if s.attempts[key] <= failuresBeforeSuccess {
		return ErrTransient
	}
	return nil
}

// Attempts implements [retrysucceeds.Mixed].
func (s *InMemory) Attempts(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var total int
	for _, n := range s.attempts {
		total += n
	}
	return total, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("retrysucceedstest: nil context")
	}
	return ctx.Err()
}
