// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ifmatch is the contract-axis fixture for the if-match contract:
// a write conditional on a predicate.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
//
// The predicate has two spellings and this fixture uses the callable one:
// Match plays the role match of the contract. `pred=` is the other form and
// carries an expression — `pred=Version==Expected` — which the resolver
// deliberately leaves verbatim, because there is no callable in it.
//
// A conformance check has to call the predicate, so only the callable form is
// reachable from this tier. The expression form is a declaration the model tier
// can act on and this one cannot.
package ifmatch

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=ifmatchtest/
type Contract interface {
	// Put is the if-match contract's writer role.
	//dokimi:contract if-match role=writer field=Body
	Put(ctx context.Context, v Value) error

	// Match is the predicate the write is conditional on. It answers about the
	// value Put takes, because a predicate over anything else is one the write
	// cannot be conditional on.
	//dokimi:contract if-match role=match
	Match(ctx context.Context, v Value) (bool, error)
}
