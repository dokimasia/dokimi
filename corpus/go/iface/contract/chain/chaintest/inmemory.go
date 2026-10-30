// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package chaintest contains the double and the suite of Contract. InMemory is
// the subject of the suite.
package chaintest

import (
	"context"
	"errors"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/chain"
)

// ErrBroken is the error of Verify for a chain whose digest does not match its
// entries.
var ErrBroken = errors.New("chaintest: chain digest does not match")

// BreakKey is the key of an entry that Append records without a link into the
// digest, so Verify reports the chain as broken.
const BreakKey = "chaintest-break"

// InMemory implements [chain.Contract] in memory. The corpus applies the
// suite of Contract to it as a subject.
//
//dokimi:subject chain.Contract new=NewInMemory
type InMemory struct {
	digest  string
	entries []chain.Entry
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [chain.Contract].
var _ chain.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Append implements [chain.Contract].
func (s *InMemory) Append(ctx context.Context, e chain.Entry) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	if e.Key != BreakKey {
		s.digest = link(s.digest, e)
	}
	return nil
}

// Replay implements [chain.Contract].
func (s *InMemory) Replay(ctx context.Context) ([]chain.Entry, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.entries), nil
}

// Verify implements [chain.Contract].
func (s *InMemory) Verify(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var want string
	for _, e := range s.entries {
		want = link(want, e)
	}
	if want != s.digest {
		return ErrBroken
	}
	return nil
}

// link appends the entry e to the digest prev. Two different chains have
// different digests.
func link(prev string, e chain.Entry) string {
	return prev + "|" + e.Key + "=" + e.Body
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("chaintest: nil context")
	}
	return ctx.Err()
}
