// Package locks grants the locks of a resource.
package locks

import "context"

// Lock grants a resource to one holder at a time.
//
//dokimi:suite
type Lock interface {
	// Acquire waits until the caller holds the lock, or until ctx ends.
	//dokimi:mixin waits
	//dokimi:mixin concurrent
	Acquire(ctx context.Context) error
	// Release releases the lock that the caller holds.
	//dokimi:mixin concurrent
	Release()
	// Held reports whether a caller holds the lock.
	Held() bool
}

// Gate lets one caller pass at a time.
//
//dokimi:suite
type Gate interface {
	// Pass waits until the gate opens, or until ctx ends.
	//dokimi:mixin waits
	Pass(ctx context.Context) error
}

// Await waits until ctx ends.
//
//dokimi:suite
//dokimi:mixin waits
func Await(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// Mutex is a lock in memory.
//
//dokimi:subject Lock new=NewMutex
type Mutex struct {
	held chan struct{}
}

// NewMutex returns a mutex that no caller holds.
func NewMutex() *Mutex {
	return &Mutex{held: make(chan struct{}, 1)}
}

// Acquire waits until the caller holds the mutex, or until ctx ends.
func (m *Mutex) Acquire(ctx context.Context) error {
	select {
	case m.held <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release releases the mutex.
func (m *Mutex) Release() {
	select {
	case <-m.held:
	default:
	}
}

// Held reports whether a caller holds the mutex.
func (m *Mutex) Held() bool {
	return len(m.held) == 1
}

// Turnstile is a gate in memory.
//
//dokimi:subject Gate new=NewTurnstile
type Turnstile struct {
	open chan struct{}
}

// NewTurnstile returns an open turnstile.
func NewTurnstile() *Turnstile {
	t := &Turnstile{open: make(chan struct{}, 1)}
	t.open <- struct{}{}
	return t
}

// Pass waits until the turnstile opens, or until ctx ends.
func (t *Turnstile) Pass(ctx context.Context) error {
	select {
	case <-t.open:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
