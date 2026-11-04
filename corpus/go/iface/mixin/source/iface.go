// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package source is the mixin-axis fixture of a call with a side effect and
// without a context, beside the fixture sideeffect.
//
// Uint64 has the shape computation, whose claim deterministic expects a
// repeat to return the same value. Each call returns the next number of the
// source, so the mixin side-effect refuses the claim with the gate
// stable-subject.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package source

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=sourcetest/
type Mixed interface {
	// Uint64 returns the next number of the source.
	//
	//dokimi:mixin side-effect
	Uint64() uint64
}
