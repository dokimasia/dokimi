package storetest

import (
	"context"
	"errors"
	"sync"

	"example.com/shop/store"
)

// ErrStoreFakeNotFound is the error of a load of StoreFake for a key
// without a value.
var ErrStoreFakeNotFound = errors.New("storetest: no value under the key")

// StoreFake implements Store in memory, and is safe for concurrent use. It
// keeps the values that its methods store in a map under their keys. The
// test of StoreFake applies the suite of Store to it.
type StoreFake struct {
	values map[string]store.Value
	mu     sync.Mutex
}

// The compiler checks that StoreFake implements Store.
var _ store.Store = (*StoreFake)(nil)

// NewStoreFake returns an empty StoreFake.
func NewStoreFake() *StoreFake {
	return &StoreFake{values: map[string]store.Value{}}
}

// Put stores v under v.Key.
func (f *StoreFake) Put(ctx context.Context, v store.Value) (store.Value, error) {
	var r0 store.Value
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[v.Key] = v
	return v, nil
}

// Get returns the value under key, and ErrStoreFakeNotFound where the fake
// has none.
func (f *StoreFake) Get(ctx context.Context, key string) (store.Value, error) {
	var r0 store.Value
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r0, present := f.values[key]
	if !present {
		return r0, ErrStoreFakeNotFound
	}
	return r0, nil
}
