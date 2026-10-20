// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisher is the contract-axis fixture for the publisher contract:
// a publish paired with the subscribe that delivers it.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package publisher

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=publishertest/
type Contract interface {
	// Publish is the publisher contract's publish role.
	//dokimi:contract publisher role=publish
	Publish(ctx context.Context, v Value) error

	// Subscribe is the publisher contract's subscribe role.
	//dokimi:contract publisher role=subscribe
	Subscribe(ctx context.Context) (<-chan Value, error)
}
