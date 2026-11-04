// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package clocktest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package clocktest

import (
	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/clock"
)

// InMemory implements [clock.Mixed] on a clock of assert. The corpus applies
// the suite of Mixed to it as a subject.
//
//dokimi:subject clock.Mixed new=NewInMemoryFromHandle
type InMemory struct {
	clk assert.Clock
}

// The compiler checks that InMemory implements [clock.Mixed].
var _ clock.Mixed = (*InMemory)(nil)

// NewInMemoryFromHandle returns an InMemory on the clock of the handle h, so
// a check that advances the clock moves the time of the subject.
func NewInMemoryFromHandle(h *MixedHandle) *InMemory { return &InMemory{clk: h.Clock()} }

// Now implements [clock.Mixed].
func (c *InMemory) Now() int64 { return c.clk.Now().UnixNano() }
