// Package store is the store that the nodes of the remote register read
// from.
package store

import "context"

// Value is a value of the store.
type Value string

// Reader reads the values of the store.
type Reader interface {
	// Read returns the value under key.
	Read(ctx context.Context, key string) (Value, error)
}
