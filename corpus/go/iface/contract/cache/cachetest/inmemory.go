// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package cachetest contains the double and the suite of Contract. InMemory is
// the control of the suite.
package cachetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/cache"
)

// ErrNotFound is the error of Lookup and Fetch for a key that neither the
// cache nor the backing store has.
var ErrNotFound = errors.New("cachetest: no value under that key")

// InMemory implements [cache.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject cache.Contract new=NewInMemory
type InMemory struct {
	mu      sync.Mutex
	backing map[string]cache.Value
	cached  map[string]cache.Value
	fetches int
}

// The compiler checks that InMemory implements [cache.Contract].
var _ cache.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		backing: map[string]cache.Value{},
		cached:  map[string]cache.Value{},
	}
}

// Store implements [cache.Contract].
func (s *InMemory) Store(ctx context.Context, key string, v cache.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backing[key] = v
	return nil
}

// Lookup implements [cache.Contract].
func (s *InMemory) Lookup(ctx context.Context, key string) (cache.Value, error) {
	if err := contextErr(ctx); err != nil {
		return cache.Value{}, err
	}
	s.mu.Lock()
	if v, hit := s.cached[key]; hit {
		s.mu.Unlock()
		return v, nil
	}
	s.mu.Unlock()

	v, err := s.Fetch(ctx, key)
	if err != nil {
		return cache.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cached[key] = v
	return v, nil
}

// Forget deletes key from the backing store and leaves the cache unchanged.
// A Lookup of key after it returns the cached value, which a pass-through
// store cannot return.
func (s *InMemory) Forget(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.backing, key)
}

// Fetch implements [cache.Contract].
func (s *InMemory) Fetch(ctx context.Context, key string) (cache.Value, error) {
	if err := contextErr(ctx); err != nil {
		return cache.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fetches++
	v, present := s.backing[key]
	if !present {
		return cache.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("cachetest: nil context")
	}
	return ctx.Err()
}
