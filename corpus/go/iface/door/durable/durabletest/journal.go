// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package durabletest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/durable"
)

// ErrCut is the error of a write to a medium after its cut.
var ErrCut = errors.New("durabletest: the medium is cut")

// Medium is a medium in memory: the records of the values by key, which each
// journal over the medium reads and writes. A cut medium refuses each write
// after a number of writes, as a crash stops a disk.
type Medium struct {
	records map[string]string
	left    int
	mu      sync.Mutex
	cut     bool
}

// NewMedium returns an empty medium without a cut.
func NewMedium() *Medium { return &Medium{records: map[string]string{}} }

// Cut makes the medium refuse each write after n more writes.
func (m *Medium) Cut(n int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.left, m.cut = n, true
}

// Heal ends the cut of the medium.
func (m *Medium) Heal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cut = false
}

// write writes the record of the key, and returns ErrCut where the cut of the
// medium refuses the write.
func (m *Medium) write(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cut {
		if m.left == 0 {
			return ErrCut
		}
		m.left--
	}
	m.records[key] = value
	return nil
}

// read returns the record of the key, and whether the medium has one.
func (m *Medium) read(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, found := m.records[key]
	return value, found
}

// Journal is the Store over a medium. Its Put writes each value in one write
// of the medium, so a cut of the medium leaves the value before the Put or the
// whole new value. The corpus applies the suite of Store to it as a control,
// with its doors.
//
//dokimi:subject durable.Store new=NewJournal options=JournalOptions
type Journal struct {
	medium *Medium
}

// NewJournal returns a journal over an empty medium.
func NewJournal() *Journal { return &Journal{medium: NewMedium()} }

// Put writes the value to the record of the key.
func (j *Journal) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return j.medium.write(key, value)
}

// Get reads the record of the key.
func (j *Journal) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	value, found := j.medium.read(key)
	if !found {
		return "", durable.ErrNotFound
	}
	return value, nil
}

// JournalOptions returns the restart and the crash of a journal. The restart
// opens a new journal over the medium of the journal. The crash cuts the
// medium after 0 or 1 writes of the call.
func JournalOptions() []StoreOption {
	return []StoreOption{conform.Restartable(reopen), conform.Crashable(crash)}
}

// reopen returns a new journal over the medium of the journal s, as a restart
// of its process opens its disk again.
func reopen(_ *StoreHandle, s durable.Store) durable.Store {
	return &Journal{medium: mediumOf(s)}
}

// crash cuts the medium of the journal s after 0 or 1 writes, which the random
// source of the check draws, runs call, and returns a new journal over the
// medium.
func crash(h *StoreHandle, s durable.Store, call func()) durable.Store {
	medium := mediumOf(s)
	medium.Cut(int(h.Rand().Uint64() % 2))
	call()
	medium.Heal()
	return &Journal{medium: medium}
}

// mediumOf returns the medium of the journal s, and a new medium for a store
// that is no journal.
func mediumOf(s durable.Store) *Medium {
	if j, isJournal := s.(*Journal); isJournal {
		return j.medium
	}
	return NewMedium()
}

// The compiler checks that Journal implements [durable.Store].
var _ durable.Store = (*Journal)(nil)
