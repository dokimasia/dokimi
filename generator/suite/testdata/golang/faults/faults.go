// Package faults declares a subject for each claim of the machinery
// call-that-fails-on-demand, a declaration without the machinery and a
// declaration with a machinery that the spec does not have.
package faults

import (
	"context"
	"errors"
)

// ErrMissing is the error of a read of a key without a value.
var ErrMissing = errors.New("faults: missing")

// Value is a value of the stores of the package.
type Value struct {
	Key  string
	Body string
}

// Cache reads values through a cache in front of its backing store.
//
//dokimi:suite machineries=[faults]
type Cache interface {
	// Get returns the value under key.
	Get(ctx context.Context, key string) (Value, error)

	// Len returns the number of cached values.
	Len() int
}

// Writer applies each write whole or not at all.
//
//dokimi:suite machineries=[faults]
type Writer interface {
	// Put stores v.
	//
	//dokimi:mixin atomic read=Get
	Put(ctx context.Context, v Value) error

	// Get returns the value under key.
	Get(ctx context.Context, key string) (Value, error)
}

// Syncer syncs with its remote, and retries a failed call of the remote.
//
//dokimi:suite machineries=[faults]
type Syncer interface {
	// Sync succeeds within three attempts.
	//
	//dokimi:mixin retry-succeeds attempts=3
	Sync(ctx context.Context) error
}

// Account runs a body in a transaction over its store.
//
//dokimi:suite machineries=[faults]
type Account interface {
	// Run runs body in a transaction.
	//
	//dokimi:contract transaction role=fn not-found=ErrMissing write=Put read=Get
	Run(ctx context.Context, body func(ctx context.Context) error) error

	// Put stores v under key.
	Put(ctx context.Context, key string, v Value) error

	// Get returns the value under key.
	Get(ctx context.Context, key string) (Value, error)
}

// Breaker guards an operation, and refuses a call once the operation failed
// repeatedly.
//
//dokimi:suite machineries=[faults]
type Breaker interface {
	// Do calls the operation.
	//
	//dokimi:contract circuit-breaker role=fn
	Do(ctx context.Context) error
}

// Balancer sends each call to one of its instances.
//
//dokimi:suite machineries=[faults]
type Balancer interface {
	// Route sends a call with key to an instance.
	//
	//dokimi:mixin sticky key=key
	Route(ctx context.Context, key string) error
}

// Copier copies a value through a resource that it opens and closes.
//
//dokimi:suite machineries=[faults]
type Copier interface {
	// Copy copies the value under key.
	//
	//dokimi:mixin leak-free open=Open close=Close
	Copy(ctx context.Context, key string) error

	// Open opens the resource of a copy.
	Open(ctx context.Context) error

	// Close closes the resource of a copy.
	Close(ctx context.Context) error
}

// Undeclared applies each write whole or not at all, and states no machinery.
//
//dokimi:suite
type Undeclared interface {
	// Put stores v.
	//
	//dokimi:mixin atomic read=Get
	Put(ctx context.Context, v Value) error

	// Get returns the value under key.
	Get(ctx context.Context, key string) (Value, error)
}

// Gravity states a machinery that the spec does not have beside faults.
//
//dokimi:suite machineries=[faults, gravity]
type Gravity interface {
	// Get returns the value under key.
	Get(ctx context.Context, key string) (Value, error)
}

// Mirror keeps a copy of the values of its remote store.
//
//dokimi:suite machineries=[faults]
type Mirror struct{}

// Get returns the value under key from the remote store.
func (*Mirror) Get(ctx context.Context, key string) (Value, error) {
	return Value{Key: key}, ctx.Err()
}

// Fetch returns the body under key from the remote store, within three
// attempts.
//
//dokimi:suite machineries=[faults]
//dokimi:mixin retry-succeeds attempts=3
func Fetch(ctx context.Context, key string) (string, error) {
	return key, ctx.Err()
}

// Pick sends a call with key to one of the instances of the remote store.
//
//dokimi:suite machineries=[faults]
//dokimi:mixin sticky key=key
func Pick(ctx context.Context, key string) error {
	return ctx.Err()
}

// Load loads the value under key through a connection that it opens and
// closes.
//
//dokimi:suite machineries=[faults]
//dokimi:mixin leak-free open=Connect close=Disconnect
func Load(ctx context.Context, key string) (Value, error) {
	return Value{Key: key}, ctx.Err()
}

// Connect opens the connection of a load.
func Connect(ctx context.Context) error {
	return ctx.Err()
}

// Disconnect closes the connection of a load.
func Disconnect(ctx context.Context) error {
	return ctx.Err()
}
