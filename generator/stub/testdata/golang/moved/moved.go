// Package moved declares the interfaces whose doubles the directive out moves
// into the package movedtest. A double there cannot refer to an unexported
// interface of this package.
//
//dokimi:out path=movedtest/
package moved

import "context"

// cache keeps values by key. Only this package can refer to it.
//
//dokimi:stub
type cache interface {
	// Get returns the value under the key.
	Get(ctx context.Context, key string) (string, error)
}

// Store keeps values by key.
//
//dokimi:stub
type Store interface {
	// Get returns the value under the key.
	Get(ctx context.Context, key string) (string, error)
}

// Cached returns the value of the key a in c, so that cache has a use.
func Cached(ctx context.Context, c cache) (string, error) { return c.Get(ctx, "a") }
