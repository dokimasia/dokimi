// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package watchertest contains the double and the suite of Contract. InMemory
// is the control of the suite.
package watchertest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/watcher"
)

// watcherBuffer is the number of changes that a subscription buffers.
const watcherBuffer = 8

// ErrUnwatchable is the error of Watch for the empty key.
var ErrUnwatchable = errors.New("watchertest: a watch needs a key")

// ErrFull is the error of Trigger for a subscription whose buffer is full.
var ErrFull = errors.New("watchertest: a watcher is too far behind to take the change")

// subscription is the subscription of Watch to the changes of one key.
type subscription struct {
	ch   chan watcher.Value
	stop func()
	once sync.Once
}

// Next returns the next change within timeout, and false where none arrives.
func (s *subscription) Next(timeout time.Duration) (watcher.Value, bool) {
	select {
	case v := <-s.ch:
		return v, true
	case <-time.After(timeout):
		return watcher.Value{}, false
	}
}

// Stop ends the subscription. A second call does nothing.
func (s *subscription) Stop() { s.once.Do(s.stop) }

// InMemory implements [watcher.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject watcher.Contract new=NewInMemory
type InMemory struct {
	mu       sync.Mutex
	watchers map[string][]*subscription
}

// The compiler checks that InMemory implements [watcher.Contract].
var _ watcher.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{watchers: map[string][]*subscription{}}
}

// Watch implements [watcher.Contract].
func (s *InMemory) Watch(ctx context.Context, key string) (watcher.Subscription, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	if key == "" {
		return nil, ErrUnwatchable
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	sub := &subscription{ch: make(chan watcher.Value, watcherBuffer)}
	sub.stop = func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		subs := s.watchers[key]
		for i, candidate := range subs {
			if candidate == sub {
				s.watchers[key] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
	s.watchers[key] = append(s.watchers[key], sub)
	return sub, nil
}

// Trigger implements [watcher.Contract].
func (s *InMemory) Trigger(ctx context.Context, key string, v watcher.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sub := range s.watchers[key] {
		select {
		case sub.ch <- v:
		default:
			return ErrFull
		}
	}
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("watchertest: nil context")
	}
	return ctx.Err()
}
