// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package causal is the mixin-axis fixture for the causal mixin, which
// declares that reads respect the happens-before order of the writes
// that produced them.
//
// The interface carries a write and a read because the claim relates them:
// a read alone has nothing to be consistent with, and the generated subtest
// would pass by having nothing to check.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package causal

import (
	"context"
)

// Value is the payload the store holds. Rev is the store-assigned ordering
// stamp the causal claim is judged by — the member the directive's
// version= names.
type Value struct {
	Key, Body string
	Rev       int64
}

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=causaltest/
type Mixed interface {
	// Store is the write the read is judged against. It answers the stored
	// value, because the store assigns Rev and a caller who never sees the
	// stamp has no order to be consistent with.
	Store(ctx context.Context, v Value) (Value, error)

	// Get reads a key. The mixin says a read never observes an effect without
	// having observed its cause, and version= names the stamp that orders
	// the causes.
	//dokimi:mixin causal version=Rev
	Get(ctx context.Context, key string) (Value, error)
}
