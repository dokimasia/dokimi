// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package hostedtest contains the suite of Store, its reference, and the
// subjects Flushed and Scheduled, which share the state and the disk format
// of a store.
package hostedtest

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"go.dokimi.dev/host"
	"go.dokimi.dev/medium"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/hosted"
)

// interval is the time between two flushes of a store on a host, on the
// monotonic clock of the host.
const interval = 10 * time.Millisecond

// store is the state that Flushed and Scheduled share: the disk of the store,
// its admission, the values that a flush made durable, and the batch that
// waits for the next flush. A store without a host writes each Put to its
// disk at once. A store on a host leaves its writes to the flusher of the
// host. Each call first asks the admission whether it may proceed, as a
// service behind admission control does, and returns the error of a failed
// admission.
//
// The disk has one file for each key, whose name is the letter k and the
// hexadecimal of the key, and whose content is the value. A flush writes the
// value of each key of its batch to a temporary file, syncs it, and renames it
// to the file of the key. It then syncs the directory, so a crash keeps each
// value whole. A Get reads the durable values alone, so a crash never takes
// back a value that a Get returned.
type store struct {
	fs medium.FS
	// admission is the dependency that admits each call, and nil for a
	// store that admits every call.
	admission *host.Dependency
	durable   map[string]string
	pending   *batch
	mu        sync.Mutex
	// direct serializes the writes of a store without a host.
	direct sync.Mutex
}

// batch is the values that wait for one flush, and whether the flush ended.
// flushed is closed once done is set.
type batch struct {
	values  map[string]string
	flushed chan struct{}
	done    bool
}

// newBatch returns an empty batch.
func newBatch() *batch {
	return &batch{values: map[string]string{}, flushed: make(chan struct{})}
}

// openStore returns a store over fsys with the values of its files, which
// admission admits the calls of. It syncs the directory of fsys first, so a
// rename that a crash of the process left pending is durable before a Get
// returns its value. A store on a host passes the sleep of its clock as
// retry, and opens again after each failed call of the disk. A store without
// a host passes nil, and opens once.
func openStore(fsys medium.FS, admission *host.Dependency, retry func()) *store {
	for {
		durable, err := load(fsys)
		if err == nil || retry == nil {
			return &store{fs: fsys, admission: admission, durable: durable, pending: newBatch()}
		}
		retry()
	}
}

// load syncs the directory of fsys and reads the value of each key.
func load(fsys medium.FS) (map[string]string, error) {
	if err := fsys.SyncDir("."); err != nil {
		return map[string]string{}, failed("sync the directory", err)
	}
	entries, err := fsys.ReadDir(".")
	if err != nil {
		return map[string]string{}, failed("list", err)
	}
	durable := map[string]string{}
	for _, e := range entries {
		hexKey, isKey := strings.CutPrefix(e.Name(), "k")
		key, err := hex.DecodeString(hexKey)
		if !isKey || err != nil {
			continue
		}
		value, err := read(fsys, e.Name())
		if err != nil {
			return map[string]string{}, err
		}
		durable[string(key)] = value
	}
	return durable, nil
}

// read returns the content of the file name of fsys.
func read(fsys medium.FS, name string) (string, error) {
	f, err := fsys.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return "", failed("open", err)
	}
	b, err := io.ReadAll(f)
	if err = errors.Join(err, f.Close()); err != nil {
		return "", failed("read", err)
	}
	return string(b), nil
}

// Get returns the durable value under the key.
func (s *store) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if err := s.admit("Get"); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	value, found := s.durable[key]
	if !found {
		return "", hosted.ErrNotFound
	}
	return value, nil
}

// admit asks the admission of the store whether the call op may proceed, and
// returns the error of a failed admission. A store without an admission admits
// every call.
func (s *store) admit(op string) error {
	if s.admission == nil {
		return nil
	}
	return failed("admit "+op, s.admission.Call(op))
}

// putNow writes the value under the key to the disk at once, as a store
// without a host does. The caller has admitted the Put.
func (s *store) putNow(key, value string) error {
	s.direct.Lock()
	defer s.direct.Unlock()
	values := map[string]string{key: value}
	if err := flush(s.fs, values); err != nil {
		return err
	}
	s.apply(values)
	return nil
}

// enqueue adds the value under the key to the pending batch, and returns the
// batch.
func (s *store) enqueue(key, value string) *batch {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending.values[key] = value
	return s.pending
}

// flushEach flushes the pending batch every interval of the clock of env,
// until the lifetime of env ends. It writes a batch again after each failed
// call of the disk, an interval later, so each Put of the batch returns once
// the batch is durable.
func (s *store) flushEach(env *host.Env) {
	life := env.Lifetime()
	for {
		env.Clock().Sleep(interval)
		if life.Err() != nil {
			return
		}
		s.mu.Lock()
		b := s.pending
		s.pending = newBatch()
		s.mu.Unlock()
		for len(b.values) > 0 && flush(s.fs, b.values) != nil {
			env.Clock().Sleep(interval)
			if life.Err() != nil {
				return
			}
		}
		s.apply(b.values)
		s.mu.Lock()
		b.done = true
		s.mu.Unlock()
		close(b.flushed)
	}
}

// ended reports whether the flush of b has ended.
func (s *store) ended(b *batch) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return b.done
}

// apply makes values durable in the state of the store.
func (s *store) apply(values map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.durable, values)
}

// flush writes values to fsys in the order of their keys, and syncs the
// directory.
func flush(fsys medium.FS, values map[string]string) error {
	for _, key := range slices.Sorted(maps.Keys(values)) {
		name := hex.EncodeToString([]byte(key))
		f, err := fsys.OpenFile("t"+name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
		if err != nil {
			return failed("create", err)
		}
		_, err = io.WriteString(f, values[key])
		if err == nil {
			err = f.Sync()
		}
		if err = errors.Join(err, f.Close()); err != nil {
			return failed("write", err)
		}
		if err := fsys.Rename("t"+name, "k"+name); err != nil {
			return failed("rename", err)
		}
	}
	return failed("sync the directory", fsys.SyncDir("."))
}

// failed returns the error err of the operation op on the disk with the
// prefix of the package, and nil for a nil err.
func failed(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("hostedtest: %s: %w", op, err)
}
