// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package commutative is the mixin-axis fixture for the commutative mixin, which
// declares that two writes applied in either order reach the same state.
//
// The interface carries a total beside the step: an algebraic claim about a
// fold needs somewhere to read the fold's result, or the generated subtest
// passes over a state nothing can see.
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
	// Apply folds one delta in. Swapping two applications changes nothing.
	//dokimi:mixin commutative
	Apply(ctx context.Context, d Delta) error

	// Total observes the fold, which is what makes the claim about Apply
	// checkable at all.
	Total(ctx context.Context) (int, error)
}
