package cachetest

import (
	"context"
	"errors"
	"sync"

	"example.com/shop/cache"
)

// ErrNoValue is the error of a load of CacheFake for a key without a
// value.
var ErrNoValue = errors.New("cachetest: no value under the key")

// CacheFake implements Cache in memory, and is safe for concurrent use. It
// keeps the values that its methods store in a map under their keys. The
// test of CacheFake applies the suite of Cache to it.
type CacheFake struct {
	values map[string]cache.Value
	mu     sync.Mutex
}

// The compiler checks that CacheFake implements Cache.
var _ cache.Cache = (*CacheFake)(nil)

// NewCacheFake returns an empty CacheFake.
func NewCacheFake() *CacheFake {
	return &CacheFake{values: map[string]cache.Value{}}
}

// Store stores v under key.
func (f *CacheFake) Store(ctx context.Context, key string, v cache.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[key] = v
	return nil
}

// Lookup returns the value under key, and ErrNoValue where the fake has
// none.
func (f *CacheFake) Lookup(ctx context.Context, key string) (cache.Value, error) {
	var r0 cache.Value
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r0, present := f.values[key]
	if !present {
		return r0, ErrNoValue
	}
	return r0, nil
}

// Fetch returns the value under key, and ErrNoValue where the fake has
// none.
func (f *CacheFake) Fetch(ctx context.Context, key string) (cache.Value, error) {
	var r0 cache.Value
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r0, present := f.values[key]
	if !present {
		return r0, ErrNoValue
	}
	return r0, nil
}
