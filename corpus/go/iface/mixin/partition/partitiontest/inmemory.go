// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package partitiontest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package partitiontest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/partition"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("partitiontest: not found")

// InMemory implements [partition.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject partition.Mixed new=NewInMemory
type InMemory struct {
	parts map[string]map[string]string
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [partition.Mixed].
var _ partition.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{parts: map[string]map[string]string{}}
}

// Put implements [partition.Mixed].
func (s *InMemory) Put(ctx context.Context, part, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.parts[part] == nil {
		s.parts[part] = map[string]string{}
	}
	s.parts[part][key] = value
	return nil
}

// Read implements [partition.Mixed].
func (s *InMemory) Read(ctx context.Context, part, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.parts[part][key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("partitiontest: nil context")
	}
	return ctx.Err()
}
