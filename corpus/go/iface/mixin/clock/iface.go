// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package clock is the mixin-axis fixture of a time-aware read without a
// context, beside the fixture timeaware.
//
// Now has the shape computation, whose claim deterministic expects a repeat
// to return the same value. Now reads the clock, so the mixin time-aware
// refuses the claim with the gate stable-subject.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package clock

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=clocktest/
type Mixed interface {
	// Now returns the current instant of the clock, in nanoseconds since the
	// Unix epoch.
	//
	//dokimi:mixin time-aware
	Now() int64
}
