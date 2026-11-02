// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package concurrenttest contains the double and the suite of Mixed. InMemory
// is the subject of the suite, and a map without a lock is its reference.
package concurrenttest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/concurrent"
)

// InMemory implements [concurrent.Mixed] in memory under a lock. The corpus
// applies the suite of Mixed to it as a subject, with a map without a lock as
// its reference.
//
//dokimi:subject concurrent.Mixed new=NewInMemory options=InMemoryOptions
type InMemory struct {
	values map[string]string
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [concurrent.Mixed].
var _ concurrent.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]string{}} }

// InMemoryOptions returns the reference of an InMemory: a map without a lock,
// which the checker of a history calls one operation at a time.
func InMemoryOptions() []MixedOption {
	return []MixedOption{conform.Reference(func(*MixedHandle) concurrent.Mixed {
		return sequential{}
	})}
}

// Put implements [concurrent.Mixed].
func (s *InMemory) Put(ctx context.Context, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = value
	return nil
}

// Get implements [concurrent.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values[key], nil
}

// sequential implements [concurrent.Mixed] over a map without a lock, so it
// is the reference of a caller at a time.
type sequential map[string]string

// Put implements [concurrent.Mixed].
func (s sequential) Put(ctx context.Context, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s[key] = value
	return nil
}

// Get implements [concurrent.Mixed].
func (s sequential) Get(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	return s[key], nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("concurrenttest: nil context")
	}
	return ctx.Err()
}
