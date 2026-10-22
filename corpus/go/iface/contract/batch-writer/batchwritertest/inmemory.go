// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package batchwritertest contains the double and the suite of Contract.
// InMemory is the control of the suite.
package batchwritertest

import (
	"context"
	"errors"
	"sync"

	batchwriter "go.dokimi.dev/dokimi/corpus/go/iface/contract/batch-writer"
)

// ErrUnkeyed is the error of Put for a value without a key.
var ErrUnkeyed = errors.New("batchwritertest: a value needs a key")

// ErrNotFound is the error of Get for a key without a value.
var ErrNotFound = errors.New("batchwritertest: no value under that key")

// InMemory implements [batchwriter.Contract] in memory. The corpus applies
// the suite of Contract to it as a control.
//
//dokimi:subject batchwriter.Contract new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	values map[string]batchwriter.Value
	order  []string
}

// The compiler checks that InMemory implements [batchwriter.Contract].
var _ batchwriter.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]batchwriter.Value{}}
}

// Put implements [batchwriter.Contract].
func (s *InMemory) Put(ctx context.Context, v batchwriter.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if v.Key == "" {
		return ErrUnkeyed
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, present := s.values[v.Key]; !present {
		s.order = append(s.order, v.Key)
	}
	s.values[v.Key] = v
	return nil
}

// Get implements [batchwriter.Contract].
func (s *InMemory) Get(ctx context.Context, key string) (batchwriter.Value, error) {
	if err := contextErr(ctx); err != nil {
		return batchwriter.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return batchwriter.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("batchwritertest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the values of prior. A restart after a crash
// keeps every value that Put wrote. Do not use prior after the call. The two
// stores share the values but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{values: prior.values, order: prior.order}
}
