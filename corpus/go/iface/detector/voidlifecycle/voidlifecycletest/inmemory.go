// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package voidlifecycletest contains the double and the suite of VoidLifecycle.
// InMemory is the control of the suite.
package voidlifecycletest

import (
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/voidlifecycle"
)

// InMemory implements [voidlifecycle.VoidLifecycle] in memory. The corpus
// applies the suite of VoidLifecycle to it as a control.
//
//dokimi:subject voidlifecycle.VoidLifecycle new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	stops int
}

// The compiler checks that InMemory implements [voidlifecycle.VoidLifecycle].
var _ voidlifecycle.VoidLifecycle = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Stop implements [voidlifecycle.VoidLifecycle].
func (s *InMemory) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stops++
}
