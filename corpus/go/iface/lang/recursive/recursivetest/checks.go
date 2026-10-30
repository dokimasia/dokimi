// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package recursivetest contains the suite of Counter, the checks that a
// consumer adds to the suite, and the subject of the suite.
package recursivetest

import (
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/prop"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/recursive"
)

// CountsTheRoot checks that Nodes counts the root of a tree that the case
// draws.
//
//dokimi:check Counter.Nodes defect=countsNothing
func CountsTheRoot(h *CounterHandle, c *prop.Case) {
	root := c.Draw(prop.Of[recursive.Node](), "root")
	assert.True(h, h.Subject().Nodes(root) >= 1, h.Statement())
}

// TakesTheFirstStep checks that Steps counts the first step of a machine that
// the case draws.
//
//dokimi:check Counter.Steps defect=countsNothing
func TakesTheFirstStep(h *CounterHandle, c *prop.Case) {
	start := c.Draw(prop.Of[recursive.Step](), "start")
	assert.True(h, h.Subject().Steps(start) >= 1, h.Statement())
}

// nothing is a Counter that counts no node and no step.
type nothing struct{}

// Nodes counts no node.
func (nothing) Nodes(recursive.Node) int { return 0 }

// Steps counts no step.
func (nothing) Steps(recursive.Step) int { return 0 }

// countsNothing returns the planted defect of both checks, a Counter that
// counts nothing.
func countsNothing(*CounterHandle) recursive.Counter { return nothing{} }
