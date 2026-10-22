// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package replica is a Log of the fixture causal-chain that follows another
// log and does not admit an entry of its own.
package replica

import (
	"context"
	"errors"

	causalchain "go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/memory"
)

// ErrReadOnly reports an append to a replica.
var ErrReadOnly = errors.New("replica: a replica admits no entry")

// Log is a read-only follower of a source log. It replays the entries of its
// source and refuses each append, so it does not pass the checks that Replay
// returns the body that Append admitted.
//
//dokimi:subject causalchain.Log new=New
//dokimi:excuse Log.Replay/keeps-the-body reason="a replica refuses each append, so no body reaches its replay"
//dokimi:excuse Log.Replay/keeps-every-body reason="a replica refuses each append, so no body reaches its replay"
type Log struct {
	source causalchain.Log
}

// New returns a replica of an empty log in memory.
func New() *Log {
	return Follow(memory.New())
}

// Follow returns a replica of the log source.
func Follow(source causalchain.Log) *Log {
	return &Log{source: source}
}

// Append refuses the entry with ErrReadOnly, because the source admits each
// entry.
func (*Log) Append(context.Context, causalchain.Entry) error {
	return ErrReadOnly
}

// Replay returns the entries of the source in the order of their positions.
func (l *Log) Replay(ctx context.Context) ([]causalchain.Entry, error) {
	//nolint:wrapcheck // a replica reports the errors of its source as the source reports them
	return l.source.Replay(ctx)
}
