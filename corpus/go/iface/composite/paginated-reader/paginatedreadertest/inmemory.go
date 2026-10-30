// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package paginatedreadertest contains the double and the suite of
// PaginatedReader. InMemory is the subject of the suite.
package paginatedreadertest

import (
	"context"
	"errors"
	"slices"
	"sync"

	paginatedreader "go.dokimi.dev/dokimi/corpus/go/iface/composite/paginated-reader"
)

// Start is the cursor at which a walk begins, and End is the cursor of the
// last page. Both are zero, so a walk ends on an empty page.
const (
	Start = 0
	End   = 0
)

// pageSize is the number of values on a page.
const pageSize = 2

// ErrUnknownCursor is the error of Fetch for a cursor that the reader did not
// issue.
var ErrUnknownCursor = errors.New("paginatedreadertest: cursor was not issued by this reader")

// InMemory implements [paginatedreader.PaginatedReader] in memory. The corpus
// applies the suite of PaginatedReader to it as a subject. A cursor is a token
// that the reader issued, so Fetch refuses an integer that no page returned.
//
//dokimi:subject paginatedreader.PaginatedReader new=NewInMemory
type InMemory struct {
	issued map[int]bool
	values []paginatedreader.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements
// [paginatedreader.PaginatedReader].
var _ paginatedreader.PaginatedReader = (*InMemory)(nil)

// NewInMemory returns a reader over three values, so a walk has two pages.
func NewInMemory() *InMemory {
	return NewInMemoryOf(
		paginatedreader.Value{Key: "a", Body: "first"},
		paginatedreader.Value{Key: "b", Body: "second"},
		paginatedreader.Value{Key: "c", Body: "third"},
	)
}

// NewInMemoryOf returns a reader over the values.
func NewInMemoryOf(values ...paginatedreader.Value) *InMemory {
	return &InMemory{values: values, issued: map[int]bool{Start: true}}
}

// Fetch implements [paginatedreader.PaginatedReader]. It issues the cursor of
// the next page before it returns the page.
func (s *InMemory) Fetch(ctx context.Context, cursor int) (paginatedreader.Page, error) {
	if err := contextErr(ctx); err != nil {
		return paginatedreader.Page{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.issued[cursor] {
		return paginatedreader.Page{}, ErrUnknownCursor
	}
	if cursor >= len(s.values) {
		return paginatedreader.Page{}, nil
	}
	to := min(cursor+pageSize, len(s.values))
	page := paginatedreader.Page{Items: slices.Clone(s.values[cursor:to])}
	if to < len(s.values) {
		s.issued[to] = true
		page.Next = to
	}
	return page, nil
}

// Add implements [paginatedreader.PaginatedReader]. It appends the value, so
// each issued cursor keeps its position.
func (s *InMemory) Add(ctx context.Context, v paginatedreader.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = append(s.values, v)
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("paginatedreadertest: nil context")
	}
	return ctx.Err()
}
