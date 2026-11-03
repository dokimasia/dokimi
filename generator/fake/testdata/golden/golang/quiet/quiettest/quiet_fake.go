package quiettest

import (
	"context"
	"errors"
	"sync"

	"example.com/shop/quiet"
)

// ErrQuietFakeNotFound is the error of a load of QuietFake for a key
// without a value.
var ErrQuietFakeNotFound = errors.New("quiettest: no value under the key")

// QuietFake implements Quiet in memory, and is safe for concurrent use. It
// keeps the values that its methods store in a map under their keys. The
// test of QuietFake applies the suite of Quiet to it.
type QuietFake struct {
	values map[string]quiet.Value
	mu     sync.Mutex
}

// The compiler checks that QuietFake implements Quiet.
var _ quiet.Quiet = (*QuietFake)(nil)

// NewQuietFake returns an empty QuietFake.
func NewQuietFake() *QuietFake {
	return &QuietFake{values: map[string]quiet.Value{}}
}

// Put stores v under v.Key.
func (f *QuietFake) Put(ctx context.Context, v quiet.Value) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[v.Key] = v
}

// Get returns the value under key, and ErrQuietFakeNotFound where the fake
// has none.
func (f *QuietFake) Get(ctx context.Context, key string) (quiet.Value, error) {
	var r0 quiet.Value
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	r0, present := f.values[key]
	if !present {
		return r0, ErrQuietFakeNotFound
	}
	return r0, nil
}
