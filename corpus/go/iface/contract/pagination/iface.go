// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package pagination is the contract-axis fixture for the pagination contract:
// a reader whose protocol cursors through pages.
//
// Every role the contract declares is present, because a contract is a
// multi-callable protocol and a fixture missing a partner exercises the
// validator's failure path rather than the contract itself.
package pagination

import (
	"context"
)

// Value is the payload the contract's roles carry.
type Value struct{ Key, Body string }

// Cursor is where a walk resumes; its zero value is both where a walk
// starts and where it ends.
type Cursor string

// Page is what one call answers: the entries it carries and the cursor a
// caller passes back to go on.
//
// A struct rather than a spread of results, because the cursor is the
// member the directive names and a member has to live on something.
type Page struct {
	Items []Value
	Next  Cursor
}

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=paginationtest/
type Contract interface {
	// Fetch is the pagination contract's reader role, and hosts the
	// directive that names its partners. It answers one page and where the
	// walk resumes — the shape both walk claims need, because a keyed read
	// has no cursor to resume from.
	//dokimi:contract pagination role=reader cursor=Next
	Fetch(ctx context.Context, cur Cursor) (Page, error)

	// Put stores an entry the walk will visit. Without a writer the
	// sequences seed nothing and every walk holds over empty pages.
	Put(ctx context.Context, v Value) error
}
