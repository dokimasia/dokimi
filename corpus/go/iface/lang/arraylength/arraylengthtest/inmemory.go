// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package arraylengthtest contains the suite of Stager. InMemory is the
// subject of the suite.
package arraylengthtest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/arraylength"
	"go.dokimi.dev/dokimi/corpus/go/iface/lang/arraylength/sizes"
)

// InMemory implements [arraylength.Stager] in memory. It counts the bytes
// that it stages. The corpus applies the suite of Stager to it as a subject.
//
//dokimi:subject arraylength.Stager new=NewInMemory
type InMemory struct {
	staged int
	mu     sync.Mutex
}

// NewInMemory returns an InMemory without staged bytes.
func NewInMemory() *InMemory { return &InMemory{} }

// Stage implements [arraylength.Stager]. It returns the error of a cancelled
// context.
func (s *InMemory) Stage(ctx context.Context, req arraylength.Request) error {
	return s.add(ctx, len(req.Salt)+len(req.Pad))
}

// Seal implements [arraylength.Stager]. It returns the error of a cancelled
// context.
func (s *InMemory) Seal(ctx context.Context, salt [sizes.SaltSize]byte) error {
	return s.add(ctx, len(salt))
}

// Pad implements [arraylength.Stager]. It returns the error of a cancelled
// context.
func (s *InMemory) Pad(ctx context.Context, pad [arraylength.PadSize]byte) error {
	return s.add(ctx, len(pad))
}

// Header implements [arraylength.Stager]. It returns the error of a cancelled
// context.
func (s *InMemory) Header(ctx context.Context, header [sizes.SaltSize + arraylength.PadSize]byte) error {
	return s.add(ctx, len(header))
}

// Key implements [arraylength.Stager]. It returns the error of a cancelled
// context.
func (s *InMemory) Key(ctx context.Context, req arraylength.KeyRequest) error {
	return s.add(ctx, len(req.Key))
}

// add counts n staged bytes, and returns the error of a cancelled context
// without a count.
func (s *InMemory) add(ctx context.Context, n int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.staged += n
	return nil
}

// The compiler checks that InMemory implements [arraylength.Stager].
var _ arraylength.Stager = (*InMemory)(nil)
