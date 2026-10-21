// Package memory keeps the items of the shop in memory.
package memory

import (
	"context"
	"errors"

	"example.com/shop/store"
)

// errNotFound reports an identity without an item.
var errNotFound = errors.New("memory: not found")

// Store is the store of items in memory.
//
//dokimi:subject store.Store new=New
type Store struct {
	items map[string]store.Item
}

// New returns an empty store.
func New() *Store {
	return &Store{items: map[string]store.Item{}}
}

// Get returns the item with the identity id.
func (s *Store) Get(_ context.Context, id string) (store.Item, error) {
	item, found := s.items[id]
	if !found {
		return store.Item{}, errNotFound
	}
	return item, nil
}

// Put stores the item.
func (s *Store) Put(_ context.Context, item store.Item) error {
	s.items[item.ID] = item
	return nil
}

// Clocked is the store of items in memory that a factory builds on the
// handle of a check.
//
//dokimi:subject store.Store new=NewClocked
type Clocked struct {
	*Store
}

// NewClocked returns an empty store for the handle of a check.
func NewClocked(handle any) *Clocked {
	_ = handle
	return &Clocked{Store: New()}
}
