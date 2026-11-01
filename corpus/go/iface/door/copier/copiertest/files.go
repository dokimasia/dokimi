// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package copiertest contains the suite of Copier and the subject of the
// suite.
package copiertest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/copier"
)

// Files is the Copier over a disk, which keeps the names of its open files.
// The corpus applies the suite of Copier to it as a subject, over a disk whose
// calls fail on demand.
//
//dokimi:subject copier.Copier new=NewFiles options=FilesOptions
type Files struct {
	disk *conform.Dependency
	open map[string]int
	mu   sync.Mutex
}

// The compiler checks that Files implements [copier.Copier].
var _ copier.Copier = (*Files)(nil)

// NewFiles returns a copier over a disk whose calls never fail.
func NewFiles() *Files { return &Files{open: map[string]int{}} }

// Copy implements [copier.Copier].
func (f *Files) Copy(ctx context.Context, name, _ string) error {
	if err := f.Open(ctx, name); err != nil {
		return err
	}
	written := f.ask("Write")
	return errors.Join(written, f.Close(ctx, name))
}

// Open implements [copier.Copier].
func (f *Files) Open(ctx context.Context, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := f.ask("Open"); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.open[name]++
	return nil
}

// Close implements [copier.Copier]. A file that is not open has nothing to
// close.
func (f *Files) Close(ctx context.Context, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := f.ask("Close"); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.open[name] > 0 {
		f.open[name]--
	}
	return nil
}

// ask returns the error of the call of the callable of the disk, and nil for a
// copier without a disk.
func (f *Files) ask(callable string) error {
	if f.disk == nil {
		return nil
	}
	if err := f.disk.Call(callable); err != nil {
		return fmt.Errorf("copiertest: %s: %w", callable, err)
	}
	return nil
}

// FilesOptions returns the door of a copier over a disk whose calls fail on
// demand.
func FilesOptions() []CopierOption {
	return []CopierOption{conform.Faultable(faulty)}
}

// faulty returns a copier over a disk that asks a dependency of the injector.
func faulty(_ *CopierHandle, inj *conform.Injector) copier.Copier {
	return &Files{disk: inj.Dependency(), open: map[string]int{}}
}
