// Package specced keeps the items of the shop in stores that only its spec
// builds.
package specced

import (
	"context"

	"example.com/shop/store"
)

// Store keeps no item. Its spec builds it through a factory of the external
// test package, which the directive names alone.
//
//dokimi:subject store.Store new=newStore
type Store struct{}

// Get returns no item.
func (*Store) Get(context.Context, string) (store.Item, error) {
	return store.Item{}, nil
}

// Put drops the item.
func (*Store) Put(context.Context, store.Item) error {
	return nil
}

// Vault keeps no item. Its spec builds it through a factory of the external
// test package, which the directive names by its import path.
//
//dokimi:subject store.Store new=example.com/shop/specced_test.newVault
type Vault struct{}

// Get returns no item.
func (*Vault) Get(context.Context, string) (store.Item, error) {
	return store.Item{}, nil
}

// Put drops the item.
func (*Vault) Put(context.Context, store.Item) error {
	return nil
}
