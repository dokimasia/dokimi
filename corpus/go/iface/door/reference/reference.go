// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package reference is the fixture of the door Reference. Each declaration
// states the machinery reference, so its suite checks that every sequence of
// calls returns on a subject what it returns on a reference implementation.
// The subjects of referencetest pass their references through the door.
package reference

import (
	"context"
	"errors"
)

// The errors of the declarations.
var (
	// ErrNotFound is the error of Store.Get for a key without a value.
	ErrNotFound = errors.New("reference: not found")

	// ErrRange is the error of Journal.At for a position past the last entry.
	ErrRange = errors.New("reference: position out of range")

	// ErrEmpty is the error of Fifo.Pop for an empty queue.
	ErrEmpty = errors.New("reference: empty queue")
)

// Store keeps values by key. Each store assigns the revisions of its own
// writes, so the check of a sequence compares its calls without the revisions.
//
//dokimi:suite machineries=[reference] normalize=Revisionless
//dokimi:out path=referencetest/
type Store interface {
	// Get returns the value under the key, and ErrNotFound for a key without
	// a value.
	Get(ctx context.Context, key string) (string, error)

	// Put stores the value under the key and returns the revision of the
	// write.
	Put(ctx context.Context, key, value string) (int64, error)

	// Delete removes the value under the key.
	Delete(ctx context.Context, key string) error
}

// Revisionless returns the values of a call of the callable operation of a
// Store without the revision that Put returns.
func Revisionless(operation string, values []any) []any {
	if operation == "Put" {
		return nil
	}
	return values
}

// Journal keeps entries in the order of their appends.
//
//dokimi:suite machineries=[reference]
//dokimi:out path=referencetest/
type Journal interface {
	// Append appends the entry and returns its position, which counts from 0.
	Append(ctx context.Context, entry string) (int, error)

	// At returns the entry at the position, and ErrRange for a position
	// outside the entries.
	At(ctx context.Context, position int) (string, error)

	// Len returns the count of the entries.
	Len(ctx context.Context) (int, error)
}

// Fifo keeps items in the order of their pushes. The oldest item leaves
// first.
//
//dokimi:suite machineries=[reference]
//dokimi:out path=referencetest/
type Fifo interface {
	// Push adds the item.
	Push(ctx context.Context, item string) error

	// Pop removes and returns the oldest item, and ErrEmpty for an empty
	// queue.
	Pop(ctx context.Context) (string, error)

	// Len returns the count of the items.
	Len(ctx context.Context) (int, error)
}
