// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisheratmostonce is the contract-axis fixture for the publisher contract's
// at-most-once bound: each subscriber receives a published message zero or one times — loss permitted, duplicates forbidden.
//
// The mode rides the directive, and the bound law it selects counts each
// subscriber's copies of a published message against it. The redeliver role
// stays undeclared here, so the bound is exercised on the single publish and
// the header says the arm went unarmed — `publisher-redeliver` and the
// exactly-once sibling are the fixtures that declare the role.
package publisheratmostonce

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=publisheratmostoncetest/
type Contract interface {
	// Publish is the publisher contract's publish role.
	//dokimi:contract publisher role=publish mode=at-most-once
	Publish(ctx context.Context, v Value) error

	// Subscribe is the publisher contract's subscribe role.
	//dokimi:contract publisher role=subscribe
	Subscribe(ctx context.Context) (<-chan Value, error)
}
