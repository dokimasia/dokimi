// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package defaultonerror is the mixin-axis fixture for the defaultonerror mixin, which
// declares that a failed read hands back the type's declared default rather
// than an undefined value.
//
// The result is a named type carrying //dokimi:default, because that is
// where the vocabulary looks for "the type's default": a claim about the
// error path's expected answer needs the declaration to state that answer,
// or there is nothing to hold the failed read against. The mixin REPLACES
// the zero-on-error claim the reader shape would otherwise license — the
// two contradict on one method, and the authored statement wins.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package defaultonerror

import (
	"context"
)

// Body is what the store answers for a key. The directive declares the
// resting value a failed read owes — the mixin's whole claim — and
// [DefaultBody] is the same value as a Go constant, so the subject and
// the seats never respell it.
//
//dokimi:default "fallback-body"
type Body string

// DefaultBody is the declared default as a value. The directive above
// carries the same spelling for the derivation; this is the one Go
// code reaches for. A var rather than a const, because a typed const
// block coalesces the named type into an enum node — and an enum
// cannot carry the default directive.
var DefaultBody = Body("fallback-body")

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=defaultonerrortest/
type Mixed interface {
	// Store is the write the read is judged against.
	Store(ctx context.Context, key string, body Body) error

	// Get reads a key, answering a failed read with the declared default
	// beside the error rather than with whatever the lookup left behind.
	//dokimi:mixin default-on-error
	Get(ctx context.Context, key string) (Body, error)
}
