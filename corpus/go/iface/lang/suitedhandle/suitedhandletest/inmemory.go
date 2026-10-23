// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package suitedhandletest contains the doubles and the suites of File and FS.
// InMemoryFile and InMemoryFS are the controls of the suites.
package suitedhandletest

import (
	"go.dokimi.dev/dokimi/corpus/go/iface/lang/suitedhandle"
)

// InMemoryFile implements [suitedhandle.File] in memory. The corpus applies
// the suite of File to it as a control.
//
//dokimi:subject suitedhandle.File new=NewInMemoryFile
type InMemoryFile struct {
	name string
}

// NewInMemoryFile returns a new InMemoryFile without a name.
func NewInMemoryFile() *InMemoryFile { return &InMemoryFile{} }

// Name implements [suitedhandle.File].
func (f *InMemoryFile) Name() string { return f.name }

// Close implements [suitedhandle.File].
func (*InMemoryFile) Close() error { return nil }

// InMemoryFS implements [suitedhandle.FS] in memory. The corpus applies the
// suite of FS to it as a control.
//
//dokimi:subject suitedhandle.FS new=NewInMemoryFS
type InMemoryFS struct {
	files map[string]*InMemoryFile
}

// NewInMemoryFS returns a new InMemoryFS without files.
func NewInMemoryFS() *InMemoryFS { return &InMemoryFS{files: map[string]*InMemoryFile{}} }

// Open implements [suitedhandle.FS].
func (s *InMemoryFS) Open(name string) (suitedhandle.File, error) {
	f, found := s.files[name]
	if !found {
		return nil, suitedhandle.ErrNotFound
	}
	return f, nil
}

// The compiler checks that InMemoryFile implements [suitedhandle.File], and
// that InMemoryFS implements [suitedhandle.FS].
var (
	_ suitedhandle.File = (*InMemoryFile)(nil)
	_ suitedhandle.FS   = (*InMemoryFS)(nil)
)
