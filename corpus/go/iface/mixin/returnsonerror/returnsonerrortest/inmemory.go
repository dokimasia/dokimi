// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package returnsonerrortest contains the double and the suite of Mixed.
// InMemory is the subject of the suite.
package returnsonerrortest

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/returnsonerror"
)

// InMemory implements [returnsonerror.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject returnsonerror.Mixed new=NewInMemory
type InMemory struct {
	records map[string][]byte
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [returnsonerror.Mixed].
var _ returnsonerror.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{records: map[string][]byte{}} }

// Put implements [returnsonerror.Mixed].
func (s *InMemory) Put(ctx context.Context, key string, record []byte) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[key] = slices.Clone(record)
	return nil
}

// Read implements [returnsonerror.Mixed].
func (s *InMemory) Read(ctx context.Context, req *returnsonerror.Request) (returnsonerror.Result, error) {
	if req == nil {
		return returnsonerror.Result{}, errors.New("returnsonerrortest: nil request")
	}
	record, err := s.fill(ctx, *req)
	return returnsonerror.Result{Record: record, Size: len(record) - len(req.Dst)}, err
}

// Fill implements [returnsonerror.Mixed].
func (s *InMemory) Fill(ctx context.Context, req returnsonerror.Request) ([]byte, error) {
	return s.fill(ctx, req)
}

// fill appends the record of req.Key to req.Dst and returns the buffer. It
// returns req.Dst with every error.
func (s *InMemory) fill(ctx context.Context, req returnsonerror.Request) ([]byte, error) {
	if err := contextErr(ctx); err != nil {
		return req.Dst, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	record, found := s.records[req.Key]
	if !found {
		return req.Dst, fmt.Errorf("returnsonerrortest: key %q: %w", req.Key, returnsonerror.ErrNotFound)
	}
	return append(req.Dst, record...), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("returnsonerrortest: nil context")
	}
	return ctx.Err()
}
