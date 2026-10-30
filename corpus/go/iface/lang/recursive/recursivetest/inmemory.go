// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package recursivetest

import "go.dokimi.dev/dokimi/corpus/go/iface/lang/recursive"

// InMemory implements [recursive.Counter]. It walks each value to its end, so
// each value of a recursive type that the suite draws has to end. The corpus
// applies the suite of Counter to it as a subject.
//
//dokimi:subject recursive.Counter new=NewInMemory
type InMemory struct{}

// NewInMemory returns an InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Nodes implements [recursive.Counter].
func (c *InMemory) Nodes(root recursive.Node) int {
	if root == nil {
		return 0
	}
	n := 1
	for _, child := range root.Children() {
		n += c.Nodes(child)
	}
	return n
}

// Steps implements [recursive.Counter].
func (*InMemory) Steps(start recursive.Step) int {
	n := 0
	for s := start; s != nil; s = s("next") {
		n++
	}
	return n
}

// The compiler checks that InMemory implements [recursive.Counter].
var _ recursive.Counter = (*InMemory)(nil)
