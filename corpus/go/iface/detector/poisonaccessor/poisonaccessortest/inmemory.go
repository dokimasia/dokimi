// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package poisonaccessortest contains the double and the suite of
// PoisonAccessor. InMemory is the control of the suite.
package poisonaccessortest

import (
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/poisonaccessor"
)

// InMemory implements [poisonaccessor.PoisonAccessor] in memory. The corpus
// applies the suite of PoisonAccessor to it as a control.
//
//dokimi:subject poisonaccessor.PoisonAccessor new=NewInMemory
type InMemory struct {
	mu        sync.Mutex
	poisoned  bool
	callCount int
}

// The compiler checks that InMemory implements [poisonaccessor.PoisonAccessor].
var _ poisonaccessor.PoisonAccessor = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Poisoned returns an InMemory whose Err returns
// [poisonaccessor.ErrPoisoned].
func Poisoned() *InMemory { return &InMemory{poisoned: true} }

// Err implements [poisonaccessor.PoisonAccessor].
func (s *InMemory) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callCount++
	if s.poisoned {
		return poisonaccessor.ErrPoisoned
	}
	return nil
}
