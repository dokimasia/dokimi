// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package tamperevidenttest contains the double and the suite of Mixed.
// InMemory is the control of the suite.
package tamperevidenttest

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/tamperevident"
)

// ErrTampered is the error of Verify for a body that changed after Store.
var ErrTampered = errors.New("tamperevidenttest: stored value does not match its digest")

// InMemory implements [tamperevident.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a control.
//
//dokimi:subject tamperevident.Mixed new=NewInMemory
type InMemory struct {
	body   string
	digest string
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [tamperevident.Mixed].
var _ tamperevident.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Store implements [tamperevident.Mixed].
func (s *InMemory) Store(ctx context.Context, body string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.body, s.digest = body, digestOf(body)
	return nil
}

// Corrupt implements [tamperevident.Mixed].
func (s *InMemory) Corrupt(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.body += "!"
	return nil
}

// Verify implements [tamperevident.Mixed].
func (s *InMemory) Verify(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if digestOf(s.body) != s.digest {
		return ErrTampered
	}
	return nil
}

// digestOf returns a digest of body. Two different bodies have different
// digests.
func digestOf(body string) string {
	return body + "#" + strconv.Itoa(len(body))
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("tamperevidenttest: nil context")
	}
	return ctx.Err()
}
