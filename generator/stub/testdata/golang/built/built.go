// Package built declares a store whose methods return values of types with a
// builder, and of types without one.
package built

import (
	"context"

	"example.com/shop/billing"
)

// Account is a value with a builder.
//
//dokimi:builder
type Account struct {
	ID    string
	Email string
}

// Page is a generic value with a builder, which builds each instance of Page.
//
//dokimi:builder
type Page[T any] struct {
	Items []T
	Next  string
}

// Note is a value without a builder.
type Note struct {
	Text string
}

// Store returns values of types with a builder, and of types without one.
//
//dokimi:stub
type Store interface {
	// Get returns an account.
	Get(ctx context.Context, id string) (Account, error)
	// Find returns a pointer to an account.
	Find(ctx context.Context, id string) (*Account, error)
	// Lookup returns an account and whether the store has it.
	Lookup(id string) (Account, bool)
	// List returns a page of accounts.
	List(ctx context.Context) (Page[Account], error)
	// Note returns a value of a type without a builder.
	Note(id string) Note
	// Bill returns a value of a type whose builder is in another package.
	Bill(id string) billing.Invoice
}
