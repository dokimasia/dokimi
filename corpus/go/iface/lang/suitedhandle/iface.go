// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package suitedhandle is the language-axis fixture for a declaration under
// the directive suite alone that a method of a stubbed interface returns. The
// declaration has one double, which the file of the producer contains.
package suitedhandle

import "errors"

// ErrNotFound reports a name that the file system does not have.
var ErrNotFound = errors.New("suitedhandle: file not found")

// File is an open file. The directive suite asks for its double, and the
// method Open of [FS] returns it.
//
//dokimi:suite
//dokimi:out path=suitedhandletest/
type File interface {
	// Name returns the name of the file.
	Name() string
	// Close closes the file.
	Close() error
}

// FS opens the files of a file system.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=suitedhandletest/
type FS interface {
	// Open returns the file with the name.
	Open(name string) (File, error)
}
