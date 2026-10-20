// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package boundedcache is the mixin-axis fixture for a bound on a KEYED
// store, beside `bounded`, which bounds a collection.
//
// The two shapes take the same stamp and want different oracles. A
// bounded collection answers a clamped list, and every store model
// answers everything it was fed, so the first clamped read reads as
// disagreement — that fixture rides the twin floor and says so. A
// bounded cache answers one key at a time, and a key it evicted is a
// MISS. Absence is legal there, which is what makes an unbounded
// reference usable: the comparison is one-sided, a subject hit must
// agree and a subject miss never disagrees, and eviction stops looking
// like a wrong answer.
//
// The bound is small on purpose. Drawn sequences have to write more
// distinct keys than the capacity or the eviction path never runs, and a
// law that never reaches its own subject reads as coverage.
package boundedcache

import (
	"context"
)

// Key identifies one entry.
//
// A directive, not a tag: a named type has no field to carry one, which is
// the only case the directive form remains for once a field can reach for
// a struct tag instead.
//
//dokimi:default "test-key"
type Key string

// Value is what an entry holds.
type Value struct {
	Body string `default:"test-value"`
}

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=boundedcachetest/
type Mixed interface {
	// Put stores value under key. At capacity the policy picks a victim
	// to make room; overwriting a held key evicts nothing. WHICH entry
	// dies is deliberately unstated — a victim policy is the consumer's
	// claim, not this interface's.
	Put(ctx context.Context, key Key, value Value) error

	// Get reports the value and whether it was there. A miss is an
	// answer rather than an error, which is the half of the declaration
	// that makes the bound checkable: with an error channel a caller
	// cannot tell an evicted key from a broken read.
	//dokimi:mixin cacheable
	Get(ctx context.Context, key Key) (Value, bool)

	// Len counts held entries, and never answers past the bound.
	//dokimi:mixin bounded limit=2
	Len(ctx context.Context) (int, error)
}
