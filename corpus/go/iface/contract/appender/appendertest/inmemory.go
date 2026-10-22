// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package appendertest contains the double and the suite of Contract. InMemory
// is the control of the suite.
package appendertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/appender"
)

// InMemory implements [appender.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject appender.Contract new=NewInMemory
type InMemory struct {
	mu      sync.Mutex
	entries []appender.Value
}

// The compiler checks that InMemory implements [appender.Contract].
var _ appender.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Run implements [appender.Contract].
func (s *InMemory) Run(ctx context.Context, v appender.Value) (int64, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, v)
	return int64(len(s.entries) - 1), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("appendertest: nil context")
	}
	return ctx.Err()
}
