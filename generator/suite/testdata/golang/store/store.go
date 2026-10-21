// Package store keeps the items of the shop.
package store

import "context"

// Item is one item of the shop.
type Item struct {
	ID   string
	Name string
}

// Store keeps items by their identity. Its directive has the prefix of the
// plugin.
//
//dokimi:suite:suite
type Store interface {
	// Get returns the item with the identity id.
	Get(ctx context.Context, id string) (Item, error)
	// Put stores the item.
	Put(ctx context.Context, item Item) error
}
