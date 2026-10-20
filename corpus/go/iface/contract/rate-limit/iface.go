// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ratelimit is the contract-axis fixture for the rate-limit contract:
// a call bounded by a rate and a burst.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package ratelimit

import (
	"context"
	"errors"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// ErrLimited reports a call refused because the burst is spent.
//
// On the interface's package rather than the subject's, because limited=
// names it and the directive resolves through this package's var scope. A
// sentinel a consumer cannot reach is one their implementation cannot
// report, and the claim is about identity: any error would let a subject
// that fails for its own reasons pass as a limiter.
var ErrLimited = errors.New("ratelimit: burst spent")

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=ratelimittest/
type Contract interface {
	// Run is the rate-limit contract's fn role, and hosts the directive
	// that names its partners.
	//dokimi:contract rate-limit role=fn rate=100 burst=10 limited=ErrLimited
	Run(ctx context.Context, key string) error
}
