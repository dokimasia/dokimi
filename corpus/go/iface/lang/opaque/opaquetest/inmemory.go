// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package opaquetest contains the suite and the double of Store. InMemory is
// the subject of the suite.
package opaquetest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/opaque"
)

// InMemory implements [opaque.Store] in memory. Each put adds the length of
// the set of its request to a total, so it reads the field of the type without
// an exported field. The corpus applies the suite of Store to it as a subject.
//
//dokimi:subject opaque.Store new=NewInMemory
type InMemory struct {
	total int
	mu    sync.Mutex
}

// NewInMemory returns an InMemory without puts.
func NewInMemory() *InMemory { return &InMemory{} }

// Put implements [opaque.Store]. It returns the error of a cancelled context.
func (s *InMemory) Put(ctx context.Context, req *opaque.PutRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total += req.Set.Len()
	return nil
}

// The compiler checks that InMemory implements [opaque.Store].
var _ opaque.Store = (*InMemory)(nil)
