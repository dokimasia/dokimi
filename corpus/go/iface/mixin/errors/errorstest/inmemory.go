// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package errorstest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package errorstest

import (
	"context"
	stderrors "errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/errors"
)

// gone is the key for which Get returns [errors.ErrGone].
const gone = "gone"

// InMemory implements [errors.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a control.
//
//dokimi:subject errors.Mixed new=NewInMemory
type InMemory struct {
	items map[string]string
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [errors.Mixed].
var _ errors.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Get implements [errors.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	if key == gone {
		return "", errors.ErrGone
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[key]
	if !ok {
		return "", errors.ErrNotFound
	}
	return v, nil
}

// Put stores value under key. A test seeds the store through it, because the
// interface declares no writer.
func (s *InMemory) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}

// GoneKey returns the key for which Get returns [errors.ErrGone].
func GoneKey() string { return gone }

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return stderrors.New("errorstest: nil context")
	}
	return ctx.Err()
}
