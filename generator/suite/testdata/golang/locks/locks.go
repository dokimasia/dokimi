// Package locks grants the locks of a resource.
package locks

import (
	"context"
	"errors"
)

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

// Board keeps the last notice for concurrent callers, and a reference
// implementation checks it.
//
//dokimi:suite machineries=[reference]
type Board interface {
	// Post stores notice.
	//dokimi:mixin concurrent
	Post(ctx context.Context, notice string) error
	// Read returns the last notice.
	//dokimi:mixin concurrent
	Read(ctx context.Context) (string, error)
}

// Shelf counts the items under each key for concurrent readers.
//
//dokimi:suite
type Shelf interface {
	// Add adds one item under key.
	Add(ctx context.Context, key string) error
	// Count returns the count of the items under key.
	//dokimi:mixin concurrent-readers
	Count(ctx context.Context, key string) (int, error)
}

// Loader loads the value of a key.
type Loader func(ctx context.Context, key string) (string, error)

// Flight shares one load of a key among the calls of the key at once.
//
//dokimi:suite machineries=[faults]
type Flight interface {
	// Do returns the value that load loads for key.
	//dokimi:contract single-flight role=fn
	Do(ctx context.Context, key string, load Loader) (string, error)
}

// ErrLocked is the error of a lock of a key that another caller holds.
var ErrLocked = errors.New("locks: the key is locked")

// Locker leases each key to one caller at a time.
//
//dokimi:suite
type Locker interface {
	// Lock takes the lease of key.
	//dokimi:contract lease role=acquire taken=ErrLocked
	Lock(ctx context.Context, key string) error
	// Unlock returns the lease of key.
	//dokimi:contract lease role=release
	Unlock(ctx context.Context, key string) error
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

// Worker runs the loop of a worker.
//
//dokimi:suite
type Worker interface {
	// Run runs the loop until ctx ends, and returns nil when ctx ends.
	//dokimi:mixin runs-until-cancelled
	Run(ctx context.Context) error
	// Stop ends the loop.
	Stop()
}

// Serve serves until ctx ends, and returns nil when ctx ends.
//
//dokimi:suite
//dokimi:mixin runs-until-cancelled
func Serve(ctx context.Context) error {
	<-ctx.Done()
	return nil
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

// Loop is a worker in memory.
//
//dokimi:subject Worker new=NewLoop
type Loop struct {
	stopped chan struct{}
}

// NewLoop returns a loop that runs.
func NewLoop() *Loop {
	return &Loop{stopped: make(chan struct{}, 1)}
}

// Run waits until ctx ends or Stop ends the loop, and returns nil.
func (l *Loop) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
	case <-l.stopped:
	}
	return nil
}

// Stop ends the loop.
func (l *Loop) Stop() {
	select {
	case l.stopped <- struct{}{}:
	default:
	}
}
