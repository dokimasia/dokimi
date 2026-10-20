// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisherexactlyonce is the contract-axis fixture for the publisher contract's
// exactly-once bound: each subscriber receives a published message exactly once — neither loss nor duplication.
//
// The mode rides the directive, and the bound law it selects counts each
// subscriber's copies of a published message against it. The redeliver role
// is armed here through Replay, whose dedupe is the exactly-once claim
// itself: the law re-offers the published message and the duplicate must
// never reach a subscriber — which is what tells this mode apart from the
// at-least-once fixture that redelivers without asking.
package publisherexactlyonce

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=publisherexactlyoncetest/
type Contract interface {
	// Publish is the publisher contract's publish role.
	//dokimi:contract publisher role=publish mode=exactly-once
	Publish(ctx context.Context, v Value) error

	// Replay is the publisher contract's redeliver role: it re-offers a
	// message the broker may already have delivered, and the exactly-once
	// claim is that the duplicate never reaches a subscriber.
	//dokimi:contract publisher role=redeliver
	Replay(ctx context.Context, v Value) error

	// Subscribe is the publisher contract's subscribe role.
	//dokimi:contract publisher role=subscribe
	Subscribe(ctx context.Context) (<-chan Value, error)
}
