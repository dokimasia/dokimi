// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package multireturntest contains the double and the suite of Wide. InMemory
// is the subject of the suite.
package multireturntest

import (
	"context"
	"errors"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/multireturn"
)

// ErrNotFound is the error of Quad for an identity without a value.
var ErrNotFound = errors.New("multireturntest: not found")

// InMemory implements [multireturn.Wide] in memory. The corpus applies the
// suite of Wide to it as a subject.
//
//dokimi:subject multireturn.Wide new=NewInMemory
type InMemory struct{ items map[string]string }

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Put stores value under id. A test seeds the store through it, because the
// interface declares no writer.
func (s *InMemory) Put(id, value string) { s.items[id] = value }

// Quad implements [multireturn.Wide].
func (s *InMemory) Quad(ctx context.Context, id string) (string, int, bool, error) {
	if err := contextErr(ctx); err != nil {
		return "", 0, false, err
	}
	v, ok := s.items[id]
	if !ok {
		return "", 0, false, ErrNotFound
	}
	return v, len(v), true, nil
}

// Triple implements [multireturn.Wide].
func (s *InMemory) Triple(_ context.Context, id string) (string, int, bool) {
	v, ok := s.items[id]
	return v, len(v), ok
}

// NoError implements [multireturn.Wide].
func (s *InMemory) NoError(_ context.Context, id string) (string, int) {
	v := s.items[id]
	return v, len(v)
}

// PartialZero is a defective Wide. Its Quad returns ErrNotFound beside nonzero
// values in its second and third results, which a check of the first result
// alone accepts.
type PartialZero struct{ *InMemory }

// Quad returns the zero value of its first result alone beside ErrNotFound.
func (PartialZero) Quad(context.Context, string) (string, int, bool, error) {
	return "", 7, true, ErrNotFound
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("multireturntest: nil context")
	}
	return ctx.Err()
}

// The compiler checks that InMemory and PartialZero implement
// [multireturn.Wide].
var (
	_ multireturn.Wide = (*InMemory)(nil)
	_ multireturn.Wide = PartialZero{}
)
