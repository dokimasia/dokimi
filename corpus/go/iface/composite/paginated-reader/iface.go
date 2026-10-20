// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package paginatedreader stacks the reader detector with the pagination
// contract on the same method.
//
// It earns a composite fixture because the contract changes what the detector
// generates rather than adding to it. A bare reader's suite asserts one call
// and one result. Once the method is pagination's reader role the same
// signature has to be driven as a walk — call, take the cursor, call again —
// and the assertion becomes that the second page is not the first, which does
// not exist for an unpaginated reader.
//
// A generator that reads the axes independently emits a single-call reader
// subtest and a pagination subtest that never runs, which passes.
package paginatedreader

import (
	"context"
)

// Value is the payload the pages carry.
type Value struct{ Key, Body string }

// Page is what one call answers. Next is an integer offset rather than an
// opaque token, so the derivation is held to the member's type rather than
// to a shape one fixture happens to have.
type Page struct {
	Items []Value
	Next  int
}

// PaginatedReader is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=paginatedreadertest/
type PaginatedReader interface {
	// Fetch is a reader by signature and pagination's reader role by
	// directive. The cursor parameter is what the contract keys on.
	//dokimi:contract pagination role=reader cursor=Next
	Fetch(ctx context.Context, cursor int) (Page, error)

	// Add stores an entry the walk will visit, so a page has something to
	// follow it. An operation nobody can call through the interface is one
	// no check can reach.
	Add(ctx context.Context, v Value) error
}
