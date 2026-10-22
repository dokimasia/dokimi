// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package nocontexttest contains the double and the suite of Calculator.
// InMemory is the control of the suite.
package nocontexttest

import (
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/nocontext"
)

// ErrDivideByZero is the error of a division by zero.
var ErrDivideByZero = errors.New("nocontexttest: divide by zero")

// InMemory implements [nocontext.Calculator] in memory. The corpus applies
// the suite of Calculator to it as a control.
//
//dokimi:subject nocontext.Calculator new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	calls int
}

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Add implements [nocontext.Calculator].
func (s *InMemory) Add(a, b int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	return a + b
}

// Divide implements [nocontext.Calculator].
func (s *InMemory) Divide(a, b int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

// Reset implements [nocontext.Calculator].
func (s *InMemory) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = 0
}

// The compiler checks that InMemory implements [nocontext.Calculator].
var _ nocontext.Calculator = (*InMemory)(nil)
