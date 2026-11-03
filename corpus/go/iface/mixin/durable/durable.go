// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package durable is the mixin-axis fixture for the durable mixin. A write
// that returned success survives a crash at any persistence point, recovery
// returns each value whole or not at all, a read reports a corrupted record,
// and a write on a full medium fails without effect.
package durable

import (
	"context"
	"errors"
)

// ErrNotFound is the error of Get for a key without a value.
var ErrNotFound = errors.New("durable: not found")

// ErrCorrupt is the error of Get for a value whose record a flipped bit
// corrupted.
var ErrCorrupt = errors.New("durable: the record is corrupt")

// Vault keeps values by key on a medium.
//
//dokimi:suite machineries=[medium]
//dokimi:out path=durabletest/
type Vault interface {
	// Put stores the value under the key, and returns once a sync covers
	// it.
	//
	//dokimi:mixin durable corrupt=ErrCorrupt
	Put(ctx context.Context, key, value string) error

	// Get returns the value under the key, ErrNotFound for a key without a
	// value, and ErrCorrupt for a value whose record a flipped bit
	// corrupted.
	Get(ctx context.Context, key string) (string, error)
}
