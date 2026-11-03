// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package cache reads values through a cache of a backing store, the fixture
// of a fake of the contract cache with a writer.
package cache

import "context"

// Value is a value of the cache.
type Value struct {
	Body string
}

// Cache reads values through a cache, and fills its backing through Store.
//
//dokimi:suite
//dokimi:fake absent=ErrNoValue
//dokimi:out path=cachetest/
type Cache interface {
	// Store fills the backing with v under key.
	Store(ctx context.Context, key string, v Value) error

	// Lookup returns the cached value under key.
	//
	//dokimi:contract cache role=cache
	Lookup(ctx context.Context, key string) (Value, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing
	Fetch(ctx context.Context, key string) (Value, error)
}
