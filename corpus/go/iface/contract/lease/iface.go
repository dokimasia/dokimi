// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package lease is the contract-axis fixture for the lease contract:
// an acquire balanced by exactly one release.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package lease

import (
	"context"
	"errors"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// ErrHeld is what Acquire reports when the lease is standing — the contract's
// own sentinel, declared beside the shape so the directive can name it.
var ErrHeld = errors.New("lease: the lease is already held")

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=leasetest/
type Contract interface {
	// Acquire is the lease contract's acquire role.
	//dokimi:contract lease role=acquire timeout=100ms taken=ErrHeld
	Acquire(ctx context.Context, key string) error

	// Release is the lease contract's release role.
	//dokimi:contract lease role=release
	Release(ctx context.Context, key string) error
}
