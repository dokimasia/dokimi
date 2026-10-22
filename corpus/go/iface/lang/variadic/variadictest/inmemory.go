// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package variadictest contains the double and the suite of Finder. InMemory is
// the control of the suite.
package variadictest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/variadic"
)

// ErrNoKeys is the error of Find for a call without keys.
var ErrNoKeys = errors.New("variadictest: no keys requested")

// InMemory implements [variadic.Finder] in memory. The corpus applies the
// suite of Finder to it as a control.
//
//dokimi:subject variadic.Finder new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items map[string]string
}

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Put stores value under key. A test seeds the finder through it, because the
// interface declares no writer.
func (s *InMemory) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}

// Find implements [variadic.Finder].
func (s *InMemory) Find(ctx context.Context, keys ...string) ([]string, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, ErrNoKeys
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if v, ok := s.items[k]; ok {
			out = append(out, v)
		}
	}
	return out, nil
}

// FindWithLimit implements [variadic.Finder].
func (s *InMemory) FindWithLimit(ctx context.Context, limit int, keys ...string) ([]string, error) {
	found, err := s.Find(ctx, keys...)
	if err != nil {
		return nil, err
	}
	if limit >= 0 && len(found) > limit {
		return found[:limit], nil
	}
	return found, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("variadictest: nil context")
	}
	return ctx.Err()
}

// The compiler checks that InMemory implements [variadic.Finder].
var _ variadic.Finder = (*InMemory)(nil)
