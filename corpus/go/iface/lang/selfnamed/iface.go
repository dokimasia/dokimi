// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package selfnamed is the language-axis fixture for interfaces with a method
// of their own name. A type that embeds such an interface has a field of that
// name, so a generated method of the same name does not compile.
package selfnamed

// Source returns the items of a feed, and false after the last item. Its one
// method has the name of the interface.
type Source interface {
	Source() (string, bool)
}

// Filter selects the items of a feed. Its first method has the name of the
// interface, and its second takes a [Source].
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=selfnamedtest/
type Filter interface {
	// Filter reports whether the item passes.
	Filter(item string) bool
	// Drain returns the number of items of src that pass.
	Drain(src Source) int
}
