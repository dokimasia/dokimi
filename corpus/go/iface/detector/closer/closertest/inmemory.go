// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package closertest contains the double and the suite of Closer. InMemory is
// the control of the suite.
package closertest

import (
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/closer"
)

// InMemory implements [closer.Closer] in memory. The corpus applies the
// suite of Closer to it as a control.
//
//dokimi:subject closer.Closer new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	closed bool
}

// The compiler checks that InMemory implements [closer.Closer].
var _ closer.Closer = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Close implements [closer.Closer].
func (s *InMemory) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}
