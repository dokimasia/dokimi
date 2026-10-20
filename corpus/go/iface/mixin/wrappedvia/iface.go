// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package wrappedvia is the mixin-axis fixture for the wrappedvia mixin, which
// declares that the annotated callable delegates to the named sibling,
// wrapping it with cross-cutting concerns — logging, metrics, retries.
//
// Delegation is a fact about the implementation's insides: both methods sit
// on one subject, and nothing observable at the interface floor says whether
// one ran through the other. The stamp documents the relationship and the
// resolver validates the named sibling; what a delegating call preserves is
// the model tier's question, asked with instrumentation the floor does not
// have.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package wrappedvia

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=wrappedviatest/
type Mixed interface {
	// Open delegates to Cause, wrapping it. The fn parameter resolves to
	// a sibling, so the delegate has to be reachable from this interface.
	//dokimi:mixin wrapped-via fn=Cause
	Open(ctx context.Context, name string) error

	// Cause is the delegate fn names.
	Cause(ctx context.Context) error
}
