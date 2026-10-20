// Package parity declares the store of the golden files of the stub
// generator of testkit, so its double compares with theirs.
package parity

import (
	"context"
	"errors"
)

// The error values that the double of Get arms.
var (
	// ErrNotFound reports an ID without an item.
	ErrNotFound = errors.New("parity: not found")
	// ErrGone reports an item that the store deleted.
	ErrGone = errors.New("parity: gone")
)

// Store has a method with named results, a method with unnamed results, and
// a method without results.
//
//dokimi:stub
type Store interface {
	// Get returns the item of an ID, selects the calls of a partition fault
	// by the ID, and has a helper of each of two error values.
	//
	//dokimi:fault errors=[ErrNotFound, ErrGone] partition=id
	Get(ctx context.Context, id string) (item string, err error)
	// List returns every item, and recovers on the third attempt.
	//
	//dokimi:fault retry=3
	List() ([]string, error)
	// Close releases the store.
	Close()
}
