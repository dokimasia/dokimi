// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deadcheck is the mutation-axis fixture for the findings of a
// mutation record. The check of Close kills no mutant of the code that it
// covers, and no check covers the method Len of the subject. The spec of the
// command report reads a record of dokimi-mutate-go over the subject and
// finds the dead check, the subsumed checks and the survivors of Len.
package deadcheck

import (
	"context"
)

// Store keeps values by key until it closes.
//
//dokimi:suite
//dokimi:out path=deadchecktest/
type Store interface {
	// Get returns the value that Put stored under key.
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, key string) (string, error)

	// Put stores value under key.
	Put(ctx context.Context, key, value string) error

	// Close marks the store closed.
	Close() error
}
