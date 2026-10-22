// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package paginationtest contains the double and the suite of Contract.
// InMemory is the control of the suite.
package paginationtest

import (
	"context"
	"errors"
	"maps"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/pagination"
)

// PageSize is the number of values on a page. It is one, so two stored values
// take two pages.
const PageSize = 1

// InMemory implements [pagination.Contract] in memory. The corpus applies
// the suite of Contract to it as a control.
//
//dokimi:subject pagination.Contract new=NewInMemory
type InMemory struct {
	values map[string]pagination.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [pagination.Contract].
var _ pagination.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]pagination.Value{}} }

// Put implements [pagination.Contract].
func (s *InMemory) Put(ctx context.Context, v pagination.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	return nil
}

// Fetch implements [pagination.Contract].
func (s *InMemory) Fetch(ctx context.Context, cur pagination.Cursor) (pagination.Page, error) {
	if err := contextErr(ctx); err != nil {
		return pagination.Page{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := slices.Sorted(maps.Keys(s.values))
	from, _ := slices.BinarySearch(keys, string(cur))
	if from < len(keys) && keys[from] == string(cur) {
		from++
	}

	var page pagination.Page
	for _, k := range keys[from:] {
		if len(page.Items) == PageSize {
			return page, nil
		}
		page.Items = append(page.Items, s.values[k])
		page.Next = pagination.Cursor(k)
	}
	return page, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("paginationtest: nil context")
	}
	return ctx.Err()
}
