// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package durable is the fixture of the doors of a handle. Its store keeps
// values on a medium apart from the store, so a new store over the medium
// reads the values of an earlier one. The suite restarts a store and crashes
// the medium of a store during a Put. The journal of durabletest is the
// control of both doors.
package durable

import (
	"context"
	"errors"
)

// ErrNotFound is the error of Get for a key without a value.
var ErrNotFound = errors.New("durable: not found")

// Store keeps values by key on a medium.
//
//dokimi:suite
//dokimi:out path=durabletest/
type Store interface {
	// Put stores the value under the key.
	Put(ctx context.Context, key, value string) error
	// Get returns the value under the key, and ErrNotFound for a key without
	// a value.
	Get(ctx context.Context, key string) (string, error)
}
