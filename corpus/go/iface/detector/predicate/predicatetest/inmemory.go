// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package predicatetest contains the double and the suite of Predicate.
// InMemory is the control of the suite.
package predicatetest

import (
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/predicate"
)

// InMemory implements [predicate.Predicate] in memory. The corpus applies
// the suite of Predicate to it as a control.
//
//dokimi:subject predicate.Predicate new=NewInMemory
type InMemory struct {
	items []string
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [predicate.Predicate].
var _ predicate.Predicate = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// IsEmpty implements [predicate.Predicate].
func (s *InMemory) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items) == 0
}
