// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package kindstest contains the suite of Every, the checks that a consumer
// adds to the suite, and the subject of the suite.
package kindstest

import (
	"fmt"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/prop"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/kinds"
)

// MeasuresTheShape checks that Measure returns the area of a shape that the
// case draws.
//
//dokimi:check Every.Measure defect=measuresNothing
func MeasuresTheShape(h *EveryHandle, c *prop.Case) {
	s := c.Draw(prop.Of[kinds.Shape](), "shape")
	assert.Equal(h, h.Subject().Measure(s), s.Area(), h.Statement())
}

// CountsTheNames checks that Index returns the number of the keys of a map
// of names that the case draws.
//
//dokimi:check Every.Index defect=measuresNothing
func CountsTheNames(h *EveryHandle, c *prop.Case) {
	names := c.Draw(prop.Of[map[fmt.Stringer]int](), "names")
	assert.Equal(h, h.Subject().Index(names), len(names), h.Statement())
}

// DrainsTheFeed checks that Drain sends each buffered value of the channels
// of a feed that the case draws.
//
//dokimi:check Every.Drain defect=measuresNothing
func DrainsTheFeed(h *EveryHandle, c *prop.Case) {
	f := c.Draw(prop.Of[kinds.Feed](), "feed")
	want := len(f.In) + len(f.Both) + len(f.Ticks)
	assert.Equal(h, h.Subject().Drain(f), want, h.Statement())
}

// nothing is an InMemory whose Measure, Index and Drain return their zero
// results without a read of their input.
type nothing struct {
	*InMemory
}

// Measure returns 0.
func (nothing) Measure(kinds.Shape) float64 { return 0 }

// Index returns 0.
func (nothing) Index(map[fmt.Stringer]int) int { return 0 }

// Drain returns 0.
func (nothing) Drain(kinds.Feed) int { return 0 }

// measuresNothing returns the planted defect of each check, an InMemory whose
// Measure, Index and Drain return their zero results.
func measuresNothing(*EveryHandle) kinds.Every { return nothing{InMemory: NewInMemory()} }
