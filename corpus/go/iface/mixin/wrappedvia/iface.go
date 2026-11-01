// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package wrappedvia is the mixin-axis fixture for the wrapped-via mixin. The
// mixin declares that the annotated callable delegates to the sibling that its
// param fn names, and adds a cross-cutting concern around the call, such as
// logging, metrics or retries.
//
// The check of the mixin calls the callable on the subject and the sibling on
// a new subject with one input, and compares their results. A wrapper that
// stops delegating returns what its sibling does not.
//
// The fixture has no negated form. eidos refuses the negation of a mixin
// directive, because a mixin is opt-in and deleting the directive is the
// suppression.
package wrappedvia

import (
	"context"
)

// Value is the payload the fixture reads.
type Value struct{ Key, Body string }

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=wrappedviatest/
type Mixed interface {
	// Get returns the value of the key.
	Get(ctx context.Context, key string) (Value, error)

	// TracedGet returns what Get returns, and annotates the error of Get. The
	// param fn names a sibling, so the delegate is a method of this
	// interface.
	//dokimi:mixin wrapped-via fn=Get
	TracedGet(ctx context.Context, key string) (Value, error)
}
