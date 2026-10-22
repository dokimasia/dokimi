// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package lookuptest contains the double and the suite of Lookup. InMemory is
// the control of the suite.
package lookuptest

import (
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/lookup"
)

// InMemory implements [lookup.Lookup] in memory. The corpus applies the
// suite of Lookup to it as a control.
//
//dokimi:subject lookup.Lookup new=NewInMemory
type InMemory struct {
	values map[string]lookup.Value
	metas  map[string]lookup.Meta
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [lookup.Lookup].
var _ lookup.Lookup = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		values: map[string]lookup.Value{},
		metas:  map[string]lookup.Meta{},
	}
}

// Inspect implements [lookup.Lookup].
func (s *InMemory) Inspect(key string) (lookup.Value, lookup.Meta, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return lookup.Value{}, lookup.Meta{}, false
	}
	return v, s.metas[key], true
}

// Put stores v and its metadata m under the key of v. A test seeds the lookup
// through it.
func (s *InMemory) Put(v lookup.Value, m lookup.Meta) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	s.metas[v.Key] = m
}
