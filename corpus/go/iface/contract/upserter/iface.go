// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package upserter is the contract-axis fixture for the upserter contract:
// an upsert paired with the read that observes it.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package upserter

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=upsertertest/
type Contract interface {
	// Put is the upserter contract's writer role.
	//dokimi:contract upserter role=writer
	Put(ctx context.Context, v Value) error

	// Get is the upserter contract's reader role.
	//dokimi:contract upserter role=reader
	Get(ctx context.Context, key string) (Value, error)
}
