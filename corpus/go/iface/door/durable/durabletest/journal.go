// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package durabletest

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"go.dokimi.dev/conform"
	"go.dokimi.dev/medium"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/durable"
)

// Journal is the Store over a medium in memory. Its Put writes the value to a
// temporary file, syncs it, renames it to the file of the key and syncs the
// directory, so a crash leaves the value before the Put or the whole new
// value. The corpus applies the suite of Store to it as a subject, with its
// doors.
//
//dokimi:subject durable.Store new=NewJournal options=JournalOptions
type Journal struct {
	medium *medium.Memory
}

// NewJournal returns a journal over an empty medium.
func NewJournal() *Journal { return &Journal{medium: medium.NewMemory()} }

// Put writes the value to the file of the key.
func (j *Journal) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	name := hex.EncodeToString([]byte(key))
	f, err := j.medium.OpenFile("t"+name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return failed("create", err)
	}
	defer func() { _ = f.Close() }()
	if _, err := io.WriteString(f, value); err != nil {
		return failed("write", err)
	}
	if err := f.Sync(); err != nil {
		return failed("sync", err)
	}
	if err := j.medium.Rename("t"+name, "k"+name); err != nil {
		return failed("rename", err)
	}
	return failed("sync the directory", j.medium.SyncDir("."))
}

// Get reads the file of the key.
func (j *Journal) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	f, err := j.medium.OpenFile("k"+hex.EncodeToString([]byte(key)), os.O_RDONLY, 0)
	if errors.Is(err, fs.ErrNotExist) {
		return "", durable.ErrNotFound
	}
	if err != nil {
		return "", failed("open", err)
	}
	defer func() { _ = f.Close() }()
	value, err := io.ReadAll(f)
	if err != nil {
		return "", failed("read", err)
	}
	return string(value), nil
}

// JournalOptions returns the restart and the crash of a journal. The restart
// opens a new journal over the medium of the journal. The crash falls in
// place of one of the next three persistence points of the medium.
func JournalOptions() []StoreOption {
	return []StoreOption{conform.Restartable(reopen), conform.Crashable(crash)}
}

// reopen returns a new journal over the medium of the journal s, as a restart
// of its process opens its disk again.
func reopen(_ *StoreHandle, s durable.Store) durable.Store {
	return &Journal{medium: mediumOf(s)}
}

// crash arms a crash of the medium of the journal s in place of one of its
// next three persistence points, which the random source of the check draws,
// runs call, crashes the medium and returns a new journal over it. A Put has
// two persistence points, so a crash at the third keeps the whole value of a
// Put that call makes.
func crash(h *StoreHandle, s durable.Store, call func()) durable.Store {
	m := mediumOf(s)
	m.CrashAt(m.Points() + int(h.Rand().Uint64()%3))
	call()
	m.Crash()
	return &Journal{medium: m}
}

// mediumOf returns the medium of the journal s, and a new medium for a store
// that is no journal.
func mediumOf(s durable.Store) *medium.Memory {
	if j, isJournal := s.(*Journal); isJournal {
		return j.medium
	}
	return medium.NewMemory()
}

// failed returns the error err of the operation op on the medium with the
// prefix of the package, and nil for a nil err.
func failed(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("durabletest: %s: %w", op, err)
}

// The compiler checks that Journal implements [durable.Store].
var _ durable.Store = (*Journal)(nil)
