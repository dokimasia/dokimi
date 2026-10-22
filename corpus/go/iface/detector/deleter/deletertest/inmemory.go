// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deletertest contains the double and the suite of Deleter. InMemory is
// the control of the suite.
package deletertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/deleter"
)

// InMemory implements [deleter.Deleter] in memory. The corpus applies the
// suite of Deleter to it as a control.
//
//dokimi:subject deleter.Deleter new=NewInMemory
type InMemory struct {
	items map[string]string
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [deleter.Deleter].
var _ deleter.Deleter = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Delete implements [deleter.Deleter].
func (s *InMemory) Delete(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
	return nil
}

// Set stores value under key. A test seeds the store through it.
func (s *InMemory) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("deletertest: nil context")
	}
	return ctx.Err()
}
