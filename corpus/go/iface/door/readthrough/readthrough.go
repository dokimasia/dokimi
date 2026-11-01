// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package readthrough is the fixture of a cache in front of a store whose
// calls fail on demand. The cache reads through to its store on a miss, and
// writes through to its store before it keeps a value. The door Faults of the
// subject builds the cache over a double of the store, whose methods ask the
// injector of the handle.
package readthrough

import (
	"context"
	"errors"
)

// ErrMissing is the error of a read of a key without a value.
var ErrMissing = errors.New("readthrough: missing")

// Value is a value that the cache keeps under its key.
type Value struct{ Key, Body string }

// Store is the store behind a cache.
//
//dokimi:stub
//dokimi:out path=readthroughtest/
type Store interface {
	// Get returns the value under the key, and ErrMissing for a key without a
	// value.
	Get(ctx context.Context, key string) (Value, error)

	// Put stores the value under its key.
	Put(ctx context.Context, v Value) error
}

// Cache keeps the values of a store in memory.
//
//dokimi:suite machineries=[faults]
//dokimi:out path=readthroughtest/
type Cache interface {
	// Get returns the value under the key from memory, and from the store on a
	// miss.
	Get(ctx context.Context, key string) (Value, error)

	// Put stores the value in the store, and keeps it in memory once the store
	// has it. It retries a failed write of the store, up to three attempts.
	//
	//dokimi:mixin atomic read=Get
	//dokimi:mixin retry-succeeds attempts=3
	Put(ctx context.Context, v Value) error
}
