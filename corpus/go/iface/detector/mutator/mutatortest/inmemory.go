// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package mutatortest contains the double and the suite of Mutator. InMemory is
// the control of the suite.
package mutatortest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/mutator"
)

// InMemory implements [mutator.Mutator] in memory. The corpus applies the
// suite of Mutator to it as a control.
//
//dokimi:subject mutator.Mutator new=NewInMemory
type InMemory struct {
	touched map[string]int
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [mutator.Mutator].
var _ mutator.Mutator = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{touched: map[string]int{}} }

// Touch implements [mutator.Mutator].
func (s *InMemory) Touch(ctx context.Context, key string) {
	if ctx == nil || ctx.Err() != nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.touched[key]++
}
