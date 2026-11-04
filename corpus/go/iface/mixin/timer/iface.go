// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package timer is the mixin-axis fixture of a command with the signature of
// a query, beside the fixture sideeffect.
//
// Stop has the shape predicate, whose claim predicate expects a repeat to
// return the same truth value. Stop stops the timer and reports whether the
// timer was active, so a repeat returns false. The mixin side-effect refuses
// the claim with the gate stable-subject.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package timer

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=timertest/
type Mixed interface {
	// Stop stops the timer, and reports whether the timer was active.
	//
	//dokimi:mixin side-effect
	Stop() bool
}
