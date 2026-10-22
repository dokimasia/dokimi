// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package cacheabletest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package cacheabletest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/cacheable"
)

// InMemory implements [cacheable.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject cacheable.Mixed new=NewInMemory
type InMemory struct {
	cache  map[string]string
	misses int
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [cacheable.Mixed].
var _ cacheable.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{cache: map[string]string{}}
}

// Get implements [cacheable.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.cache[key]; ok {
		return v, nil
	}
	s.misses++
	v := "computed:" + key
	s.cache[key] = v
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("cacheabletest: nil context")
	}
	return ctx.Err()
}
