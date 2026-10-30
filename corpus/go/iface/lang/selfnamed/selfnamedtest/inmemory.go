// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package selfnamedtest contains the double and the suite of Filter.
// InMemory is the subject of the suite.
package selfnamedtest

import (
	"go.dokimi.dev/dokimi/corpus/go/iface/lang/selfnamed"
)

// InMemory implements [selfnamed.Filter] in memory. An item passes when it is
// not empty. The corpus applies the suite of Filter to it as a subject.
//
//dokimi:subject selfnamed.Filter new=NewInMemory
type InMemory struct{}

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Filter implements [selfnamed.Filter].
func (*InMemory) Filter(item string) bool { return item != "" }

// Drain implements [selfnamed.Filter].
func (f *InMemory) Drain(src selfnamed.Source) int {
	n := 0
	for item, ok := src.Source(); ok; item, ok = src.Source() {
		if f.Filter(item) {
			n++
		}
	}
	return n
}

// The compiler checks that InMemory implements [selfnamed.Filter].
var _ selfnamed.Filter = (*InMemory)(nil)
