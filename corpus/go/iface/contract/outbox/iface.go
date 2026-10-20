// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package outbox is the contract-axis fixture for the outbox contract:
// a transactional outbox and its reader.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package outbox

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=outboxtest/
type Contract interface {
	// Append is the outbox contract's append role.
	//dokimi:contract outbox role=append
	Append(ctx context.Context, v Value) error

	// Subscribe is the outbox contract's subscribe role.
	//dokimi:contract outbox role=subscribe
	Subscribe(ctx context.Context) (<-chan Value, error)
}
