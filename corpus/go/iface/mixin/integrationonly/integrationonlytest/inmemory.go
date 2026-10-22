// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package integrationonlytest contains the double and the suite of Mixed.
// InMemory is the control of the suite.
package integrationonlytest

import (
	"context"
	"errors"
	"strings"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/integrationonly"
)

// ErrBadDSN is the error of Connect for a data source name without a scheme.
var ErrBadDSN = errors.New("integrationonlytest: malformed dsn")

// InMemory implements [integrationonly.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a control.
//
//dokimi:subject integrationonly.Mixed new=NewInMemory
type InMemory struct {
	mu        sync.Mutex
	connected bool
}

// The compiler checks that InMemory implements [integrationonly.Mixed].
var _ integrationonly.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Connect implements [integrationonly.Mixed].
func (s *InMemory) Connect(ctx context.Context, dsn string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if !strings.Contains(dsn, "://") {
		return ErrBadDSN
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connected = true
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("integrationonlytest: nil context")
	}
	return ctx.Err()
}
