// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package embeddedforeigntest contains the double and the suite of Stream.
// InMemory is the control of the suite.
package embeddedforeigntest

import (
	"context"
	"errors"
	"io"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/embeddedforeign"
)

// ErrNotFound is the error of Read for a key without a value.
var ErrNotFound = errors.New("embeddedforeigntest: not found")

// InMemory implements [embeddedforeign.Stream] in memory. The corpus applies
// the suite of Stream to it as a control.
//
//dokimi:subject embeddedforeign.Stream new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	items  map[string]string
	closed bool
}

// The compiler checks that InMemory implements [embeddedforeign.Stream] and the
// interface [io.Closer] that Stream embeds.
var (
	_ embeddedforeign.Stream = (*InMemory)(nil)
	_ io.Closer              = (*InMemory)(nil)
)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Read implements [embeddedforeign.Stream].
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

// Close implements [io.Closer]. A second call returns nil too.
func (s *InMemory) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

// Put stores value under key. A test seeds the stream through it, because the
// interface declares no writer.
func (s *InMemory) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("embeddedforeigntest: nil context")
	}
	return ctx.Err()
}
