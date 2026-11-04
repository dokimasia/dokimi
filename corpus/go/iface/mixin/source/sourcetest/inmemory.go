// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sourcetest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package sourcetest

import (
	"sync/atomic"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/source"
)

// golden is the increment of the state of SplitMix64.
const golden = 0x9e3779b97f4a7c15

// InMemory implements [source.Mixed] with SplitMix64, whose state is one
// integer. The corpus applies the suite of Mixed to it as a subject.
//
//dokimi:subject source.Mixed new=NewInMemory
type InMemory struct {
	state atomic.Uint64
}

// The compiler checks that InMemory implements [source.Mixed].
var _ source.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory at the state 0.
func NewInMemory() *InMemory { return &InMemory{} }

// Uint64 implements [source.Mixed].
func (s *InMemory) Uint64() uint64 {
	z := s.state.Add(golden)
	z = (z ^ z>>30) * 0xbf58476d1ce4e5b9
	z = (z ^ z>>27) * 0x94d049bb133111eb
	return z ^ z>>31
}
