// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package cache is the contract-axis fixture for the cache contract:
// a cache in front of a backing store.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
//
// The model tier rides the twin floor, and what a store would need is the
// loader Fetch calls on a miss: the cache's state is whatever that
// returned, and no stamp names it. The same bar ref.NewCoalescer fails.
package cache

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=cachetest/
type Contract interface {
	// Store fills the backing store, and is on the interface for the reason
	// tx's Put is: a contract stating that a cached read agrees with its
	// backing states something about a store with something IN it, and a
	// state no conformance check can establish through the interface is one
	// no check can reach. Held off it, the agreement bound against two
	// misses — which agree by both finding nothing.
	Store(ctx context.Context, key string, v Value) error

	// Lookup is the cache contract's cache role.
	//dokimi:contract cache role=cache
	Lookup(ctx context.Context, key string) (Value, error)

	// Fetch is the cache contract's backing role.
	//dokimi:contract cache role=backing
	Fetch(ctx context.Context, key string) (Value, error)
}
