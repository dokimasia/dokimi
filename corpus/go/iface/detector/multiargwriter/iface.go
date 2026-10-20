// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package multiargwriter is the detector-axis fixture for the multiargwriter shape:
// three or more non-context parameters — past the composite pair, at the
// arity boundary the detector documents and draws.
//
// No directive appears here: the classification comes from the signature
// alone, so a detector that misfires shows up as wrong generated output
// rather than as a directive that was misread.
package multiargwriter

import (
	"context"
	"errors"
)

// ErrNotFound is what Get reports for a key Set never filed.
var ErrNotFound = errors.New("multiargwriter: not found")

// Value is the stored form.
type Value struct{ Key, Body string }

// MultiArgWriter is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=multiargwritertest/
type MultiArgWriter interface {
	// Set carries three non-context parameters — past the composite pair,
	// into the boundary the detector draws at three or more.
	Set(ctx context.Context, key, body, mime string) error

	// Get reads back what Set filed, and is here so the model tier has
	// something to compare. A lone writer answers only an error, so both
	// sides return nil for every call a correct subject makes and the
	// differential has nothing to disagree about — which left the shape's
	// own action constructor emitted by nothing in the corpus.
	Get(ctx context.Context, key string) (string, error)
}
