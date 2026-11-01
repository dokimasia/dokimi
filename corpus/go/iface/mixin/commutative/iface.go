// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package commutative is the mixin-axis fixture for the commutative mixin,
// which declares that a callable of two inputs of one type returns one value
// for the two inputs in either order.
//
// Sum adds the amounts of two deltas, so its result has another type than its
// inputs. Its planted defect therefore cannot return its first input, and it
// calls Sum with a zero second input instead.
//
// Apply keeps the mixin on a callable of one input. The law of the mixin
// needs two inputs of one type, so the gate two-inputs-of-one-type refuses
// its claim.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package commutative

import (
	"context"
)

// Delta is one contribution to the fold.
type Delta struct {
	Key    string
	Amount int
}

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=commutativetest/
type Mixed interface {
	// Apply folds one delta in.
	//dokimi:mixin commutative
	Apply(ctx context.Context, d Delta) error

	// Total observes the fold.
	Total(ctx context.Context) (int, error)

	// Sum returns the sum of the amounts of a and b.
	//
	//dokimi:mixin commutative
	Sum(a, b Delta) int
}
