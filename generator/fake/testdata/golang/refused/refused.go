// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package refused contains declarations under the directive fake that no
// reference form covers or that no fake can implement, the fixtures of the
// refusals of fake. The roles of each declaration are in a contract instance
// of their own.
package refused

import (
	"context"
	"time"
)

// Value is a value of the stores.
type Value struct {
	Key  string
	Body string
}

// Expiring stores values for a minute, a reader with a mixin beside its role.
//
//dokimi:suite
//dokimi:fake
type Expiring interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer id=expiring
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under key, and ErrGone once its minute passed.
	//
	//dokimi:contract persister role=reader id=expiring
	//dokimi:mixin ttl duration=1m put=Put read=Get not-found=ErrGone
	Get(ctx context.Context, key string) (Value, error)
}

// ErrGone is the error of Get for a value whose minute passed.
var ErrGone = context.DeadlineExceeded

// Counted stores values and counts them, a callable without a role.
//
//dokimi:suite
//dokimi:fake
type Counted interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer id=counted
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under key.
	//
	//dokimi:contract persister role=reader id=counted
	Get(ctx context.Context, key string) (Value, error)

	// Since returns the time since the first value.
	Since(ctx context.Context) (time.Duration, error)
}

// Unkeyed reads by an id that no field of the value has, a load whose key
// selects nothing of the store.
//
//dokimi:suite
//dokimi:fake
type Unkeyed interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer id=unkeyed
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under id.
	//
	//dokimi:contract persister role=reader id=unkeyed
	Get(ctx context.Context, id int) (Value, error)
}

// Unsuited has no directive suite.
//
//dokimi:fake
type Unsuited interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer id=unsuited
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under key.
	//
	//dokimi:contract persister role=reader id=unsuited
	Get(ctx context.Context, key string) (Value, error)
}

// Misnamed has a param absent that reads a field of no name template.
//
//dokimi:suite
//dokimi:fake absent=Err{{.Method}}
type Misnamed interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer id=misnamed
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under key.
	//
	//dokimi:contract persister role=reader id=misnamed
	Get(ctx context.Context, key string) (Value, error)
}

// hidden is not exported, so the test of its fake cannot refer to it.
//
//dokimi:suite
//dokimi:fake
type hidden interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer id=hidden
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under key.
	//
	//dokimi:contract persister role=reader id=hidden
	Get(ctx context.Context, key string) (Value, error)
}

// Sealed has a method that its package does not export.
//
//dokimi:suite
//dokimi:fake
type Sealed interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer id=sealed
	Put(ctx context.Context, v Value) (Value, error)

	// get returns the value under key.
	get(ctx context.Context, key string) (Value, error)
}

// Nothing has no method.
//
//dokimi:suite
//dokimi:fake
type Nothing interface{}

// Bytes keys its values by a slice, which no Go map can key.
//
//dokimi:suite
//dokimi:fake
type Bytes interface {
	// Store fills the backing with v under key.
	Store(ctx context.Context, key []byte, v Value) error

	// Lookup returns the cached value under key.
	//
	//dokimi:contract cache role=cache id=bytes
	Lookup(ctx context.Context, key []byte) (Value, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing id=bytes
	Fetch(ctx context.Context, key []byte) (Value, error)
}

// Mapped keys its values by a map, which no Go map can key.
//
//dokimi:suite
//dokimi:fake
type Mapped interface {
	// Store fills the backing with v under key.
	Store(ctx context.Context, key map[string]string, v Value) error

	// Lookup returns the cached value under key.
	//
	//dokimi:contract cache role=cache id=mapped
	Lookup(ctx context.Context, key map[string]string) (Value, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing id=mapped
	Fetch(ctx context.Context, key map[string]string) (Value, error)
}

// Record is a value whose key is unexported, so a fake in another package
// cannot read it.
type Record struct {
	key  string
	Body string
}

// Private reads by a key that only an unexported field of the record has.
//
//dokimi:suite
//dokimi:fake
type Private interface {
	// Put stores r.
	//
	//dokimi:contract persister role=writer id=private
	Put(ctx context.Context, r Record) (Record, error)

	// Get returns the record under key.
	//
	//dokimi:contract persister role=reader id=private
	Get(ctx context.Context, key string) (Record, error)
}

// Loose keys its values by a type parameter whose constraint is not
// comparable.
//
//dokimi:suite
//dokimi:fake
//dokimi:witness K=string
type Loose[K any] interface {
	// Store fills the backing with v under key.
	Store(ctx context.Context, key K, v Value) error

	// Lookup returns the cached value under key.
	//
	//dokimi:contract cache role=cache id=loose
	Lookup(ctx context.Context, key K) (Value, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing id=loose
	Fetch(ctx context.Context, key K) (Value, error)
}

// Keyable has type-set elements, so it only constrains type parameters.
type Keyable interface {
	~string | ~float64
}

// Tagged has a type parameter without a witness.
//
//dokimi:suite
//dokimi:fake
type Tagged[K Keyable] interface {
	// Store fills the backing with v under key.
	Store(ctx context.Context, key K, v Value) error

	// Lookup returns the cached value under key.
	//
	//dokimi:contract cache role=cache id=tagged
	Lookup(ctx context.Context, key K) (Value, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing id=tagged
	Fetch(ctx context.Context, key K) (Value, error)
}
