// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package durabletest contains the subject and the suite of Vault. Files is
// the subject of the suite.
package durabletest

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"os"

	"go.dokimi.dev/conform"
	"go.dokimi.dev/medium"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/durable"
)

// crcSize is the size of the CRC-32 after the value of a record.
const crcSize = 4

// Files implements [durable.Vault] over a medium.FS, with one file for each
// key: the value, then its CRC-32. A Put writes a temporary file, syncs it,
// renames it to the file of the key and syncs the directory, so a crash
// leaves the old value or the new one. The corpus applies the suite of Vault
// to it as a subject, with the door Open of FilesOptions.
//
//dokimi:subject durable.Vault new=NewFiles options=FilesOptions
type Files struct {
	fs medium.FS
}

// The compiler checks that Files implements [durable.Vault].
var _ durable.Vault = (*Files)(nil)

// NewFiles returns a Files over a new medium in memory.
func NewFiles() *Files { return OpenFiles(medium.NewMemory()) }

// OpenFiles returns a Files over m, which reads the values that an earlier
// Files over m wrote.
func OpenFiles(m medium.FS) *Files { return &Files{fs: m} }

// FilesOptions returns the options of the suite of Vault on Files: the door
// Open, which opens a Files over the medium of each subject of a check.
func FilesOptions() []VaultOption {
	return []VaultOption{
		conform.Durable(func(_ *VaultHandle, m medium.FS) durable.Vault { return OpenFiles(m) }),
	}
}

// Put implements [durable.Vault].
func (s *Files) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	name, temporary := fileOf(key), "t"+hex.EncodeToString([]byte(key))
	if err := s.write(temporary, value); err != nil {
		return err
	}
	if err := s.fs.Rename(temporary, name); err != nil {
		return failed("rename", err)
	}
	return failed("sync the directory", s.fs.SyncDir("."))
}

// Get implements [durable.Vault].
func (s *Files) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	f, err := s.fs.OpenFile(fileOf(key), os.O_RDONLY, 0)
	if errors.Is(err, fs.ErrNotExist) {
		return "", durable.ErrNotFound
	}
	if err != nil {
		return "", failed("open", err)
	}
	defer func() { _ = f.Close() }()
	record, err := io.ReadAll(f)
	if err != nil {
		return "", failed("read", err)
	}
	if len(record) < crcSize {
		return "", durable.ErrCorrupt
	}
	value, sum := record[:len(record)-crcSize], record[len(record)-crcSize:]
	if crc32.ChecksumIEEE(value) != binary.BigEndian.Uint32(sum) {
		return "", durable.ErrCorrupt
	}
	return string(value), nil
}

// write writes the record of value to the file name and syncs the file.
func (s *Files) write(name, value string) error {
	f, err := s.fs.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return failed("create", err)
	}
	defer func() { _ = f.Close() }()
	record := binary.BigEndian.AppendUint32([]byte(value), crc32.ChecksumIEEE([]byte(value)))
	if _, err := f.Write(record); err != nil {
		return failed("write", err)
	}
	return failed("sync", f.Sync())
}

// fileOf returns the name of the file of key: the letter k and the
// hexadecimal of the key, so every key has a valid name.
func fileOf(key string) string { return "k" + hex.EncodeToString([]byte(key)) }

// failed returns the error err of the operation op on the medium with the
// prefix of the package, and nil for a nil err.
func failed(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("durabletest: %s: %w", op, err)
}
