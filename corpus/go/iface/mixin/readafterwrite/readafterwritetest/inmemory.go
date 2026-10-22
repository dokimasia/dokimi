// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package readafterwritetest contains the double and the suite of Mixed.
// InMemory is the control of the suite.
package readafterwritetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/readafterwrite"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("readafterwritetest: not found")

// InMemory implements [readafterwrite.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a control.
//
//dokimi:subject readafterwrite.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items map[string]string
}

// The compiler checks that InMemory implements [readafterwrite.Mixed].
var _ readafterwrite.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("readafterwritetest: nil context")
	}
	return ctx.Err()
}

// Write implements [readafterwrite.Mixed].
func (s *InMemory) Write(ctx context.Context, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
	return nil
}

// Read implements [readafterwrite.Mixed].
func (s *InMemory) Read(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}
