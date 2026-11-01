// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package copier is the fixture of a callable that acquires a resource of a
// disk whose calls fail on demand. Copy opens a file of the disk, writes to
// it and closes it, also after a failed write.
package copier

import "context"

// Copier copies values to the files of a disk.
//
//dokimi:suite machineries=[faults]
//dokimi:out path=copiertest/
type Copier interface {
	// Copy opens the file with the name, writes the value to it and closes
	// it. It closes the file after a failed write as well.
	//
	//dokimi:mixin leak-free open=Open close=Close
	Copy(ctx context.Context, name, value string) error

	// Open opens the file with the name.
	Open(ctx context.Context, name string) error

	// Close closes the file with the name.
	Close(ctx context.Context, name string) error
}
