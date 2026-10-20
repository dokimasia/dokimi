// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package saga is the contract-axis fixture for the saga contract:
// two steps, each paired with the compensation that undoes it.
//
// Two pairs because compensation is an ORDER. A saga that fails partway
// unwinds the steps before the failure, last first, and leaves the failed
// one alone — none of which a single-step fixture can state.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package saga

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=sagatest/
type Contract interface {
	// Step is the saga contract's step role.
	//dokimi:contract saga role=step id=step
	Step(ctx context.Context, v Value) error

	// Compensate is the saga contract's compensate role.
	//dokimi:contract saga role=compensate id=step
	Compensate(ctx context.Context, v Value) error

	// Ship is a SECOND step, with its own compensator.
	//
	// Two pairs rather than one, and the ordering is the reason. A saga
	// compensates in reverse from the step that failed, so a single-step
	// fixture leaves that whole claim unstated: there is nothing before the
	// failure to undo. With two, a failure in the second has to unwind the
	// first and only the first.
	//
	// The id ship separates this pair from the pair of Step.
	//dokimi:contract saga role=step id=ship
	Ship(ctx context.Context, v Value) error

	// UndoShip undoes what Ship applied.
	//dokimi:contract saga role=compensate id=ship
	UndoShip(ctx context.Context, v Value) error

	// State reports the applied-step fingerprint, in order — the
	// observation a compensation claim compares before and after a failed
	// run. A string rather than a listing because the comparison is
	// equality, and one composed value cannot half-match.
	State(ctx context.Context) (string, error)
}
