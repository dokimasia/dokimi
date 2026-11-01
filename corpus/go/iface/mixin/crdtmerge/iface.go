// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package crdtmerge is the fixture of the mixin crdt-merge, which declares
// that replicas reach the same state whatever the order of their merges.
//
// The check makes two replicas diverge through the write that the param write
// names, merges them in both orders, and compares the states that the read of
// the param read returns. The interface therefore declares the merge, the
// write and the read.
//
// The fixture has no negated form. A mixin is opt-in, so eidos refuses the
// negation of its directive.
package crdtmerge

import (
	"context"
)

// Replica is the peer that Merge reads. A subject of Mixed implements it, so
// a check can pass one replica to the merge of another.
//
//dokimi:stub
//dokimi:out path=crdtmergetest/
type Replica interface {
	Items(ctx context.Context) ([]string, error)
}

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=crdtmergetest/
type Mixed interface {
	// Merge adds the items of peer.
	//dokimi:mixin crdt-merge write=Add read=Items
	Merge(ctx context.Context, peer Replica) error

	// Add adds an item, which makes the replicas diverge.
	Add(ctx context.Context, item string) error

	// Items returns the items of the replica.
	Items(ctx context.Context) ([]string, error)
}
