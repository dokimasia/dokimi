// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package batchreader is the detector-axis fixture for the batch-reader shape,
// which reads many keys in one call. The variadic list of keys marks the
// shape. The writer Put stores the values that the batch read returns.
//
// No directive appears here. The classification comes from the signature
// alone, so a detector that misfires shows up as wrong generated output
// rather than as a directive that was misread.
package batchreader

import (
	"context"
)

// Value is the payload the fixture reads.
type Value struct{ Key, Body string }

// BatchReader is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=batchreadertest/
type BatchReader interface {
	GetAll(ctx context.Context, keys ...string) ([]Value, error)
	Put(ctx context.Context, v Value) error
}
