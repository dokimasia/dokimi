// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package causalchain is the composite fixture pairing the chain contract
// with the causal mixin: an append-only log whose admission policy is
// causality — an entry lands only after everything it depends on.
//
// The pairing exists because the replay-causality rule needs both stamps on
// one method: the chain says the log replays, the causal claim says the
// replay respects the dependency order, and neither alone states the rule.
package causalchain

import (
	"context"
)

// RootEntry is the fixture's declared Entry.
//
// Declared rather than derived, because DependsOn is an admission
// PRECONDITION and not a payload. A derived entry fills every member it can
// and so names a dependency that never landed, which this log refuses whole:
// the append fails, the replay finds nothing, and the chain check skips
// having asserted anything. An entry depending on nothing is the one that
// always lands, which is what a claim about the log's length needs.
func RootEntry() Entry {
	return Entry{ID: "root", Body: "held"}
}

// Entry is the log's element: an identifier, the entries it depends on, the
// payload they order, and the position the log stamped it at.
//
// Seq is the store's, not the caller's. causal's version= names the member
// carrying the ordering stamp, and a stamp the caller chose would make the
// replay-order claim a question about the check's own arithmetic rather than
// about the log.
//
//dokimi:sample value=RootEntry()
type Entry struct {
	ID        string
	DependsOn []string
	Body      string
	Seq       int64
}

// Log is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=causalchaintest/
type Log interface {
	// Append admits an entry only after everything it depends on, and
	// stamps it with the position it took.
	//
	// Both stamps together are the claim: chain says the log replays,
	// causal says which member orders it, and the replay coming back in
	// that order is what neither states alone.
	//dokimi:contract chain role=append
	//dokimi:mixin causal version=Seq
	Append(ctx context.Context, e Entry) error

	// Replay is the chain contract's replay role.
	//dokimi:contract chain role=replay
	Replay(ctx context.Context) ([]Entry, error)
}
