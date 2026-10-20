// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deleter is the detector-axis fixture for the deleter shape:
// a removal named for it — the key of the entity to remove in, only an
// error out.
//
// Structurally this is the writer shape, and the name is the
// discrimination: nothing in the signature separates a Delete from a
// Put, and the two license opposite expectations — a law deriving
// write-then-read-back against a delete asserts the reverse of correct
// behaviour. The name gate is what keeps that from happening.
//
// No directive appears here: the classification comes from the
// signature and the name alone, so a detector that misfires shows up
// as wrong generated output rather than as a directive that was
// misread.
package deleter

import (
	"context"
)

// Deleter is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=deletertest/
type Deleter interface {
	Delete(ctx context.Context, key string) error
}
