// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package seededreadertest contains the double and the suite of Catalog.
// InMemory is the subject of the suite.
package seededreadertest

import (
	"context"
	"errors"
	"maps"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/seededreader"
)

// InMemory implements [seededreader.Catalog] in memory. The corpus applies the
// suite of Catalog to it as a subject.
//
//dokimi:subject seededreader.Catalog new=NewInMemory
type InMemory struct {
	docs map[seededreader.Key]seededreader.Body
}

// The compiler checks that InMemory implements [seededreader.Catalog].
var _ seededreader.Catalog = (*InMemory)(nil)

// NewInMemory returns a catalog of one document, whose key and body are the
// declared defaults of Key and Body.
func NewInMemory() *InMemory {
	return NewInMemoryOf(map[seededreader.Key]seededreader.Body{"test-doc": "test-contents"})
}

// NewInMemoryOf returns a catalog of a copy of docs.
func NewInMemoryOf(docs map[seededreader.Key]seededreader.Body) *InMemory {
	return &InMemory{docs: maps.Clone(docs)}
}

// Lookup implements [seededreader.Catalog].
func (s *InMemory) Lookup(ctx context.Context, key seededreader.Key) (seededreader.Body, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	return s.docs[key], nil
}

// Len implements [seededreader.Catalog].
func (s *InMemory) Len(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	return len(s.docs), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("seededreadertest: nil context")
	}
	return ctx.Err()
}
