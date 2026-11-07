// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package hosted is the fixture of the door Hosted. Its store keeps values by
// key on the disk of a host, and a flusher of the store writes the pending
// values to the disk on a timer. A Put returns once a flush covers its value,
// and each call first asks a dependency of the host for admission. The suite
// runs the store in a world that crashes it and strikes the faults of its
// disk and of its dependency. hostedtest has two subjects: Flushed runs its
// flusher on a goroutine in native mode under the option Blocking, and
// Scheduled runs it as a task of the host in Tasks mode.
package hosted

import (
	"context"
	"errors"
)

// ErrNotFound is the error of Get for a key without a value.
var ErrNotFound = errors.New("hosted: not found")

// Store keeps values by key on the disk of its host.
//
//dokimi:suite machineries=[medium, reference, faults]
//dokimi:out path=hostedtest/
type Store interface {
	// Put stores the value under the key, and returns once a flush makes it
	// durable.
	//
	//dokimi:mixin durable
	//dokimi:mixin concurrent
	Put(ctx context.Context, key, value string) error

	// Get returns the durable value under the key, and ErrNotFound for a key
	// without a value.
	//
	//dokimi:mixin concurrent
	Get(ctx context.Context, key string) (string, error)
}
