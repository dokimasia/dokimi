// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package recursive is the fixture of inputs whose generated results contain
// their own type: a tree of nodes and a state function. The suite generates
// such a type through prop.Recursive, so each value ends at its base values.
package recursive

// Node is a node of a tree. Its children are nodes.
type Node interface {
	// Name returns the name of the node.
	Name() string
	// Children returns the children of the node.
	Children() []Node
}

// Step is a step of a state machine. It returns the next step, and nil at the
// end.
type Step func(input string) Step

// Counter counts the parts of recursive values.
//
//dokimi:suite
//dokimi:out path=recursivetest/
type Counter interface {
	// Nodes counts the nodes of the tree under root.
	Nodes(root Node) int
	// Steps counts the steps of a run of the machine from start.
	Steps(start Step) int
}
