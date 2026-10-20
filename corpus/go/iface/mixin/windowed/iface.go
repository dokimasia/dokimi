// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package windowed is the mixin-axis fixture for the windowed mixin, which
// declares that a count covers a bounded interval rather than all of history.
//
// The directive names the pair it governs and the interval itself, because a
// window nobody stated is a number the generator would have to invent.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package windowed

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=windowedtest/
type Mixed interface {
	// Record adds one occurrence at the clock's current reading.
	//dokimi:mixin windowed incr=Record count=CountIn window=1m
	Record(ctx context.Context, key string) error

	// CountIn reports occurrences inside the declared window, which is what
	// makes the count bounded rather than cumulative.
	CountIn(ctx context.Context, key string) (int, error)
}
