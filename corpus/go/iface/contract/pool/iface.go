// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package pool is the contract-axis fixture for the pool contract:
// a pool whose every Get is returned by a Put.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package pool

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Stats is the accounting the pool laws balance — on the interface per
// the database/sql.DBStats precedent, and returned by the stats role of
// the contract so the generator knows where the numbers live.
type Stats struct{ Gets, Puts, Outstanding int }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=pooltest/
type Contract interface {
	// Get is the pool contract's get role.
	//dokimi:contract pool role=get
	Get(ctx context.Context) (Value, error)

	// Put is the pool contract's put role.
	//dokimi:contract pool role=put
	Put(ctx context.Context, v Value) error

	// Stats is the pool contract's stats role: the accounting the
	// balanced and leak-free laws verify against the cycle count.
	//dokimi:contract pool role=stats
	Stats(ctx context.Context) (Stats, error)
}
