// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package overmatch is the mixin-axis fixture for the overmatch mixin, which
// declares that a drain yields at least the elements required of it,
// and may yield more.
//
// The interface carries an append beside the drain: a claim about what a
// drain yields needs something to have gone in, or the generated subtest
// passes over an empty collection.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package overmatch

import (
	"context"
)

// Value is the element the collection holds.
type Value struct{ Key, Body string }

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=overmatchtest/
type Mixed interface {
	// Add puts an element in for the drain to yield.
	Add(ctx context.Context, v Value) error

	// Items drains the collection, yielding a superset of what was required.
	//dokimi:mixin over-match
	Items(ctx context.Context) ([]Value, error)
}
