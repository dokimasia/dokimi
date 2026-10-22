// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package multiargwritertest contains the double and the suite of
// MultiArgWriter. InMemory is the control of the suite.
package multiargwritertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/multiargwriter"
)

// InMemory implements [multiargwriter.MultiArgWriter] in memory. The corpus
// applies the suite of MultiArgWriter to it as a control.
//
//dokimi:subject multiargwriter.MultiArgWriter new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	values map[string]string
}

// The compiler checks that InMemory implements [multiargwriter.MultiArgWriter].
var _ multiargwriter.MultiArgWriter = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]string{}} }

// Set implements [multiargwriter.MultiArgWriter].
func (s *InMemory) Set(ctx context.Context, key, body, mime string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = mime + ":" + body
	return nil
}

// Get implements [multiargwriter.MultiArgWriter].
func (s *InMemory) Get(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return "", multiargwriter.ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("multiargwritertest: nil context")
	}
	return ctx.Err()
}
