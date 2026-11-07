// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package hostedtest

import (
	"context"

	"go.dokimi.dev/conform"
	"go.dokimi.dev/host"
	"go.dokimi.dev/medium"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/hosted"
)

// admissionName is the name of the dependency that admits the calls of a store
// on a host.
const admissionName = "admission"

// Flushed is the Store whose flusher runs on a goroutine of its host. A Put
// adds its value to the pending batch, and waits on a channel for the flush
// of the batch. The goroutine flushes the batch every interval, so a Put
// waits for a timer of the host. The corpus applies the suite of Store to it
// as a subject, in native mode under the option Blocking. No goroutine of
// the store waits while it has locked its mutex.
//
//dokimi:subject hosted.Store new=NewFlushed options=FlushedOptions
type Flushed struct {
	*store
	env *host.Env
}

// The compiler checks that Flushed implements [hosted.Store].
var _ hosted.Store = (*Flushed)(nil)

// NewFlushed returns a Flushed without a host over a new medium in memory,
// which writes each Put at once and admits every call.
func NewFlushed() *Flushed { return OpenFlushed(medium.NewMemory()) }

// OpenFlushed returns a Flushed without a host over fsys, which reads the
// values that an earlier store over fsys wrote.
func OpenFlushed(fsys medium.FS) *Flushed { return &Flushed{store: openStore(fsys, nil, nil)} }

// HostFlushed returns a Flushed over the disk of env, whose flusher runs on a
// goroutine of env, and which the dependency admission of env admits the
// calls of.
func HostFlushed(env *host.Env) *Flushed {
	s := &Flushed{
		store: openStore(env.Disk(), env.Dependency(admissionName), func() { env.Clock().Sleep(interval) }),
		env:   env,
	}
	env.Go(func() { s.flushEach(env) })
	return s
}

// FlushedOptions returns the options of the suite of Store on Flushed: the
// doors Open, Faults and Host, the reference, and the option Blocking.
func FlushedOptions() []StoreOption {
	return []StoreOption{
		conform.Durable(func(_ *StoreHandle, m medium.FS) hosted.Store { return OpenFlushed(m) }),
		conform.Faultable(func(_ *StoreHandle, inj *conform.Injector) hosted.Store {
			return &Flushed{store: openStore(medium.NewMemory(), inj.Dependency(), nil)}
		}),
		conform.Hosted(func(_ *StoreHandle, env *host.Env) hosted.Store { return HostFlushed(env) }),
		conform.Reference(func(*StoreHandle) hosted.Store { return newReference() }),
		conform.Blocking[hosted.Store](),
	}
}

// Put implements [hosted.Store].
func (s *Flushed) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.admit("Put"); err != nil {
		return err
	}
	if s.env == nil {
		return s.putNow(key, value)
	}
	b := s.enqueue(key, value)
	select {
	case <-b.flushed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
