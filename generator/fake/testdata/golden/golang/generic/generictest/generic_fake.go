package generictest

import (
	"context"
	"errors"
	"sync"

	"example.com/shop/generic"
)

// ErrCacheFakeNotFound is the error of a load of CacheFake for a key
// without a value.
var ErrCacheFakeNotFound = errors.New("generictest: no value under the key")

// CacheFake implements Cache in memory, and is safe for concurrent use. It
// keeps the values that its methods store in a map under their keys. The
// test of CacheFake applies the suite of Cache to it.
type CacheFake[V any, K comparable] struct {
	values map[K]V
	mu     sync.Mutex
}

// The compiler checks that CacheFake implements Cache.
var _ generic.Cache[int, string] = (*CacheFake[int, string])(nil)

// NewCacheFake returns an empty CacheFake.
func NewCacheFake[V any, K comparable]() *CacheFake[V, K] {
	return &CacheFake[V, K]{values: map[K]V{}}
}

// Store stores v under key.
func (f *CacheFake[V, K]) Store(ctx context.Context, key K, v V) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[key] = v
	return nil
}

// Lookup returns the value under key, and ErrCacheFakeNotFound where the
// fake has none.
func (f *CacheFake[V, K]) Lookup(ctx context.Context, key K) (V, error) {
	var r0 V
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r0, present := f.values[key]
	if !present {
		return r0, ErrCacheFakeNotFound
	}
	return r0, nil
}

// Fetch returns the value under key, and ErrCacheFakeNotFound where the
// fake has none.
func (f *CacheFake[V, K]) Fetch(ctx context.Context, key K) (V, error) {
	var r0 V
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r0, present := f.values[key]
	if !present {
		return r0, ErrCacheFakeNotFound
	}
	return r0, nil
}
