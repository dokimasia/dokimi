// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package associative is the mixin-axis fixture for the associative mixin,
// which declares that a callable of two inputs of one type returns one value
// for three inputs in either grouping.
//
// Merge combines two deltas into a delta, so its result has the type of its
// inputs. Merge is commutative as well, so the planted defect of each of its
// two claims keeps the law of the other. The defect of commutative returns
// its first input, which keeps the grouping. The defect of associative keeps
// the value of its first call, which keeps the order.
//
// Apply keeps the mixin on a callable of one input, a batch of deltas. The law
// of the mixin needs two inputs of one type, so the gate
// two-inputs-of-one-type refuses its claim.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package associative

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
//dokimi:out path=associativetest/
type Mixed interface {
	// Apply folds a batch of deltas in.
	//dokimi:mixin associative
	Apply(ctx context.Context, ds []Delta) error

	// Total observes the fold.
	Total(ctx context.Context) (int, error)

	// Merge returns the delta with the smaller key of a and b and the sum of
	// their amounts.
	//
	//dokimi:mixin associative
	//dokimi:mixin commutative
	Merge(a, b Delta) Delta
}
