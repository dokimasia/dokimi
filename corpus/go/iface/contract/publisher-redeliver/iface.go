// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisherredeliver is the contract-axis fixture that arms the
// publisher's redeliver role: the at-least-once bound proven through an
// actual redelivery rather than on the single publish.
//
// Republish re-offers a message the broker already delivered, so each
// subscriber sees the duplicate the mode permits — which is exactly what
// tells at-least-once apart from exactly-once. Its unarmed siblings
// exercise the role's optional omission, and their headers say so.
package publisherredeliver

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=publisherredelivertest/
type Contract interface {
	// Publish is the publisher contract's publish role.
	//dokimi:contract publisher role=publish mode=at-least-once
	Publish(ctx context.Context, v Value) error

	// Republish is the publisher contract's redeliver role: it re-offers a
	// message to every current subscriber, the duplicate this mode permits.
	//dokimi:contract publisher role=redeliver
	Republish(ctx context.Context, v Value) error

	// Subscribe is the publisher contract's subscribe role.
	//dokimi:contract publisher role=subscribe
	Subscribe(ctx context.Context) (<-chan Value, error)
}
