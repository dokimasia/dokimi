// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package puretest contains the double and the suite of Pure. InMemory is the
// control of the suite.
package puretest

import (
	"go.dokimi.dev/dokimi/corpus/go/iface/detector/pure"
)

// InMemory implements [pure.Pure] in memory. The corpus applies the suite of
// Pure to it as a control.
//
//dokimi:subject pure.Pure new=NewInMemory
type InMemory struct{ label string }

// The compiler checks that InMemory implements [pure.Pure].
var _ pure.Pure = (*InMemory)(nil)

// NewInMemory returns an InMemory with the label control.
func NewInMemory() *InMemory { return NewInMemoryLabelled("control") }

// NewInMemoryLabelled returns an InMemory whose description names label.
func NewInMemoryLabelled(label string) *InMemory { return &InMemory{label: label} }

// Describe implements [pure.Pure].
func (s *InMemory) Describe() string { return "in-memory: " + s.label }
