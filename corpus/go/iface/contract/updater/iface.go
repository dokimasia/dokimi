// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package updater is the contract-axis fixture for the updater contract:
// an update paired with the read that observes it.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package updater

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=updatertest/
type Contract interface {
	// Put is the updater contract's writer role.
	//dokimi:contract updater role=writer
	Put(ctx context.Context, v Value) error

	// Get is the updater contract's reader role.
	//dokimi:contract updater role=reader
	Get(ctx context.Context, key string) (Value, error)
}
