// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package injectionsafetest contains the double and the suite of Mixed.
// InMemory is the control of the suite.
package injectionsafetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/injectionsafe"
)

// ErrMissing is the error of a read of a key without a value.
var ErrMissing = errors.New("injectionsafetest: nothing stored under the key")

// InMemory implements [injectionsafe.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a control.
//
//dokimi:subject injectionsafe.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	values map[string]string
}

// The compiler checks that InMemory implements [injectionsafe.Mixed].
var _ injectionsafe.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]string{}} }

// Store implements [injectionsafe.Mixed].
func (s *InMemory) Store(ctx context.Context, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = value
	return nil
}

// Load implements [injectionsafe.Mixed].
func (s *InMemory) Load(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.values[key]
	if !ok {
		return "", ErrMissing
	}
	return value, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("injectionsafetest: nil context")
	}
	return ctx.Err()
}
