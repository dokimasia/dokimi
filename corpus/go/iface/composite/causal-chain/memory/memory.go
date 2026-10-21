// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package memory is the Log of the fixture causal-chain in memory.
package memory

import (
	"context"
	"errors"
	"slices"
	"sync"

	causalchain "go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain"
)

// ErrUnadmitted reports an entry that depends on an entry that the log does
// not have.
var ErrUnadmitted = errors.New("memory: the entry depends on an entry that the log does not have")

// Log is the causal log in memory.
//
//dokimi:subject causalchain.Log new=New
type Log struct {
	entries []causalchain.Entry
	mu      sync.Mutex
}

// New returns an empty log.
func New() *Log {
	return &Log{}
}

// Append admits the entry once the log has each entry that it depends on, and
// stamps the entry with its position. It returns ErrUnadmitted for an entry
// that depends on an entry that the log does not have.
func (l *Log) Append(ctx context.Context, e causalchain.Entry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, id := range e.DependsOn {
		if !slices.ContainsFunc(l.entries, func(have causalchain.Entry) bool { return have.ID == id }) {
			return ErrUnadmitted
		}
	}
	e.Seq = int64(len(l.entries))
	l.entries = append(l.entries, e)
	return nil
}

// Replay returns the entries in the order of their positions.
func (l *Log) Replay(ctx context.Context) ([]causalchain.Entry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return slices.Clone(l.entries), nil
}
