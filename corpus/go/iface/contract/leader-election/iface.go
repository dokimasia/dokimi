// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package leaderelection is the contract-axis fixture for the leader-election contract:
// a leadership protocol.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
//
// The model tier rides the twin floor. ref.LeaseTracker is the near miss:
// it models exclusion, but it is keyed by what acquire takes and Campaign
// takes nothing, so the store has no key to hold the leadership under.
package leaderelection

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=leaderelectiontest/
type Contract interface {
	// Campaign is the leader-election contract's campaign role.
	//dokimi:contract leader-election role=campaign
	Campaign(ctx context.Context) error

	// Resign is the leader-election contract's resign role.
	//dokimi:contract leader-election role=resign
	Resign(ctx context.Context) error

	// IsLeader is the leader-election contract's is-leader role.
	//dokimi:contract leader-election role=is-leader
	IsLeader(ctx context.Context) bool
}
