// Package durable keeps values on a medium apart from the process of its
// keeper.
package durable

import "context"

// Keeper keeps values by key.
//
//dokimi:suite
type Keeper interface {
	// Put stores the value under the key.
	Put(ctx context.Context, key, value string) error
	// Get returns the value under the key.
	Get(ctx context.Context, key string) (string, error)
}
