// Package store declares a store whose methods configure the faults of
// their doubles.
package store

import (
	"context"
	"errors"
)

// The error values that the double of Get arms.
var (
	// ErrNotFound reports an ID without an item.
	ErrNotFound = errors.New("store: not found")
	// ErrGone reports an item that the store deleted.
	ErrGone = errors.New("store: gone")
)

// Store keeps items.
type Store interface {
	// Get returns the item of an ID, and recovers on the third attempt. A
	// partition fault selects its calls by the ID, and its double arms two
	// error values.
	//
	//dokimi:fault errors=[ErrNotFound, ErrGone] retry=3 partition=id
	Get(ctx context.Context, id string) (string, error)
	// List returns every item, and recovers on the second attempt.
	//
	//dokimi:fault retry=2
	List() ([]string, error)
	// Close releases the store, and configures no fault.
	Close() error
}
