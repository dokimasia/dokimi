// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package generic reads values of a type parameter through a cache, the
// fixture of a generic fake.
package generic

import "context"

// Cache reads values of V under keys of K through a cache, and fills its
// backing through Store. The type parameter of the keys is the second.
//
//dokimi:suite
//dokimi:fake
//dokimi:witness V=int K=string
//dokimi:out path=generictest/
type Cache[V any, K comparable] interface {
	// Store fills the backing with v under key.
	Store(ctx context.Context, key K, v V) error

	// Lookup returns the cached value under key.
	//
	//dokimi:contract cache role=cache
	Lookup(ctx context.Context, key K) (V, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing
	Fetch(ctx context.Context, key K) (V, error)
}
