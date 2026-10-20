// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package timeaware is the mixin-axis fixture for the timeaware mixin, which
// marks a callable whose behaviour depends on a clock — without saying which
// quantity, because the quantities belong to the claims layered on top
// (`ttl`, `timeout`, `windowed`).
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package timeaware

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=timeawaretest/
type Mixed interface {
	// Touch records that the key was seen, at whatever the clock reads.
	Touch(ctx context.Context, key string) error

	// AgeOf reports how long ago, which is the dependency the mixin marks:
	// the answer moves when the clock does and not otherwise.
	//dokimi:mixin time-aware
	AgeOf(ctx context.Context, key string) (int64, error)
}
