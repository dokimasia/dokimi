// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package puretest contains the double and the suite of Mixed. InMemory is the
// subject of the suite.
package puretest

import (
	"strings"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/pure"
)

// InMemory implements [pure.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a subject.
//
//dokimi:subject pure.Mixed new=NewInMemory
type InMemory struct{}

// The compiler checks that InMemory implements [pure.Mixed].
var _ pure.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Derive implements [pure.Mixed].
func (*InMemory) Derive(input string) string { return strings.ToUpper(input) }
