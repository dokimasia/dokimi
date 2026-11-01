// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package referencetest

import (
	"context"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/reference"
)

// SliceLog implements [reference.Journal] over a slice of entries. The corpus
// applies the suite of Journal to it with a map of positions as its
// reference.
//
//dokimi:subject reference.Journal new=NewSliceLog options=SliceLogOptions
type SliceLog struct {
	entries []string
	mu      sync.Mutex
}

// The compiler checks that SliceLog implements [reference.Journal].
var _ reference.Journal = (*SliceLog)(nil)

// NewSliceLog returns an empty SliceLog.
func NewSliceLog() *SliceLog { return &SliceLog{} }

// SliceLogOptions returns the reference of a SliceLog: a journal over a map
// of positions.
func SliceLogOptions() []JournalOption {
	return []JournalOption{conform.Reference(func(*JournalHandle) reference.Journal {
		return &mapLog{entries: map[int]string{}}
	})}
}

// Append implements [reference.Journal].
func (l *SliceLog) Append(ctx context.Context, entry string) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, entry)
	return len(l.entries) - 1, nil
}

// At implements [reference.Journal].
func (l *SliceLog) At(ctx context.Context, position int) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if position < 0 || position >= len(l.entries) {
		return "", reference.ErrRange
	}
	return l.entries[position], nil
}

// Len implements [reference.Journal].
func (l *SliceLog) Len(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.entries), nil
}

// mapLog is the reference of a SliceLog: the entries in a map by their
// positions.
type mapLog struct {
	entries map[int]string
}

// Append implements [reference.Journal].
func (l *mapLog) Append(_ context.Context, entry string) (int, error) {
	position := len(l.entries)
	l.entries[position] = entry
	return position, nil
}

// At implements [reference.Journal].
func (l *mapLog) At(_ context.Context, position int) (string, error) {
	entry, stored := l.entries[position]
	if !stored {
		return "", reference.ErrRange
	}
	return entry, nil
}

// Len implements [reference.Journal].
func (l *mapLog) Len(context.Context) (int, error) { return len(l.entries), nil }
