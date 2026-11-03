package clashtest

import (
	"context"
	"errors"
	"sync"

	"example.com/shop/clash"
)

// ErrShadowFakeNotFound is the error of a load of ShadowFake for a key
// without a value.
var ErrShadowFakeNotFound = errors.New("clashtest: no value under the key")

// ShadowFake implements Shadow in memory, and is safe for concurrent use.
// It keeps the values that its methods store in a map under their keys.
// The test of ShadowFake applies the suite of Shadow to it.
type ShadowFake struct {
	values map[string]clash.Value
	mu     sync.Mutex
}

// The compiler checks that ShadowFake implements Shadow.
var _ clash.Shadow = (*ShadowFake)(nil)

// NewShadowFake returns an empty ShadowFake.
func NewShadowFake() *ShadowFake {
	return &ShadowFake{values: map[string]clash.Value{}}
}

// Put stores r0 under r0.Present.
func (f0 *ShadowFake) Put(f context.Context, r0 clash.Value) (clash.Value, error) {
	var r00 clash.Value
	if err := f.Err(); err != nil {
		return r00, err
	}
	f0.mu.Lock()
	defer f0.mu.Unlock()
	f0.values[r0.Present] = r0
	return r0, nil
}

// Get returns the value under present, and ErrShadowFakeNotFound where the
// fake has none.
func (f0 *ShadowFake) Get(f context.Context, present string) (clash.Value, error) {
	var r0 clash.Value
	if err := f.Err(); err != nil {
		return r0, err
	}
	f0.mu.Lock()
	defer f0.mu.Unlock()
	r0, present0 := f0.values[present]
	if !present0 {
		return r0, ErrShadowFakeNotFound
	}
	return r0, nil
}
