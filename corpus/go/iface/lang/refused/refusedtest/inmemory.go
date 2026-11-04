// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package refusedtest contains the double and the suite of Clock. InMemory is
// the subject of the suite.
package refusedtest

import (
	"fmt"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/refused"
)

// maxOffset is the max offset of the wall time of a peer beyond the wall time
// of the clock.
const maxOffset = 500 * time.Millisecond

// InMemory implements [refused.Clock] with the update rule of a hybrid
// logical clock, on a clock of assert. The corpus applies the suite of Clock
// to it as a subject.
//
//dokimi:subject refused.Clock new=NewInMemoryFromHandle
type InMemory struct {
	clk  assert.Clock
	last refused.Instant
	mu   sync.Mutex
}

// The compiler checks that InMemory implements [refused.Clock].
var _ refused.Clock = (*InMemory)(nil)

// NewInMemoryFromHandle returns an InMemory on the clock of the handle h.
func NewInMemoryFromHandle(h *ClockHandle) *InMemory { return &InMemory{clk: h.Clock()} }

// Update implements [refused.Clock].
func (c *InMemory) Update(observed refused.Instant) (refused.Instant, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clk.Now().UnixNano()
	if observed.Wall > now+maxOffset.Nanoseconds() {
		return refused.Instant{}, fmt.Errorf("refusedtest: wall time %d: %w", observed.Wall, refused.ErrOffset)
	}
	next := refused.Instant{Wall: max(now, c.last.Wall, observed.Wall)}
	switch next.Wall {
	case c.last.Wall:
		next.Logical = c.last.Logical + 1
		if next.Wall == observed.Wall {
			next.Logical = max(c.last.Logical, observed.Logical) + 1
		}
	case observed.Wall:
		next.Logical = observed.Logical + 1
	}
	c.last = next
	return next, nil
}
