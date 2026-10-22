// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package nilsafetest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package nilsafetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/nilsafe"
)

// ErrNilPayload is the error of Store for a nil payload.
var ErrNilPayload = errors.New("nilsafetest: nil payload")

// InMemory implements [nilsafe.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject nilsafe.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	values map[string]nilsafe.Payload
}

// The compiler checks that InMemory implements [nilsafe.Mixed].
var _ nilsafe.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]nilsafe.Payload{}}
}

// Store implements [nilsafe.Mixed].
func (s *InMemory) Store(ctx context.Context, v *nilsafe.Payload) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if v == nil {
		return ErrNilPayload
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = *v
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("nilsafetest: nil context")
	}
	return ctx.Err()
}
