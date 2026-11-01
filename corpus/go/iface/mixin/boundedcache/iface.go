// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package boundedcache is the fixture of a bound on a keyed store, beside the
// fixture bounded, which bounds a collection.
//
// A bounded cache returns one key at a time, and a key that it evicted is a
// miss. A miss is a legal result, so a check compares a hit with an unbounded
// reference and never reports a miss as a wrong value.
//
// The bound is small. A sequence writes more distinct keys than the capacity,
// so the eviction runs within a few calls.
package boundedcache

import (
	"context"
)

// Key identifies one entry.
//
// A named type has no field for a struct tag, so it takes its default from
// the directive.
//
//dokimi:default "test-key"
type Key string

// Value is the value of an entry.
type Value struct {
	Body string `default:"test-value"`
}

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=boundedcachetest/
type Mixed interface {
	// Put stores value under key. At its capacity the cache evicts an entry
	// to make room, and a write of a stored key evicts nothing. The interface
	// does not state which entry the cache evicts, because the eviction
	// policy is a claim of the consumer. The use of Put is the count that
	// Len returns, and the limit is two entries.
	//dokimi:mixin bounded limit=2 observe=Len
	Put(ctx context.Context, key Key, value Value) error

	// Get returns the value of key and whether the cache has it. A miss is
	// a result and not an error, so a caller can tell an evicted key from a
	// failed read.
	//dokimi:mixin cacheable
	Get(ctx context.Context, key Key) (Value, bool)

	// Len returns the count of the stored entries.
	Len(ctx context.Context) (int, error)
}
