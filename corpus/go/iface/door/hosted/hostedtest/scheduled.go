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

// Scheduled is the Store whose flusher runs as a task of its host. A Put adds
// its value to the pending batch, and waits for the flush of the batch through
// Env.Wait, so it blocks through the host alone. The corpus applies the suite
// of Store to it as a subject, in Tasks mode, which replays each case whole.
//
//dokimi:subject hosted.Store new=NewScheduled options=ScheduledOptions
type Scheduled struct {
	*store
	env *host.Env
}

// The compiler checks that Scheduled implements [hosted.Store].
var _ hosted.Store = (*Scheduled)(nil)

// NewScheduled returns a Scheduled without a host over a new medium in
// memory, which writes each Put at once and admits every call.
func NewScheduled() *Scheduled { return OpenScheduled(medium.NewMemory()) }

// OpenScheduled returns a Scheduled without a host over fsys, which reads the
// values that an earlier store over fsys wrote.
func OpenScheduled(fsys medium.FS) *Scheduled { return &Scheduled{store: openStore(fsys, nil, nil)} }

// HostScheduled returns a Scheduled over the disk of env, whose flusher runs
// as a task of env, and which the dependency admission of env admits the
// calls of.
func HostScheduled(env *host.Env) *Scheduled {
	s := &Scheduled{
		store: openStore(env.Disk(), env.Dependency(admissionName), func() { env.Clock().Sleep(interval) }),
		env:   env,
	}
	env.Go(func() { s.flushEach(env) })
	return s
}

// ScheduledOptions returns the options of the suite of Store on Scheduled:
// the doors Open, Faults and Host, the reference, and the option Tasks.
func ScheduledOptions() []StoreOption {
	return []StoreOption{
		conform.Durable(func(_ *StoreHandle, m medium.FS) hosted.Store { return OpenScheduled(m) }),
		conform.Faultable(func(_ *StoreHandle, inj *conform.Injector) hosted.Store {
			return &Scheduled{store: openStore(medium.NewMemory(), inj.Dependency(), nil)}
		}),
		conform.Hosted(func(_ *StoreHandle, env *host.Env) hosted.Store { return HostScheduled(env) }),
		conform.Reference(func(*StoreHandle) hosted.Store { return newReference() }),
		conform.Tasks[hosted.Store](),
	}
}

// Put implements [hosted.Store].
func (s *Scheduled) Put(ctx context.Context, key, value string) error {
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
	s.env.Wait(func() bool { return s.ended(b) || ctx.Err() != nil })
	if !s.ended(b) {
		return ctx.Err()
	}
	return nil
}
