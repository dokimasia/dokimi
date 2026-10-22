// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package built is the integration fixture for the default results of a
// double, which the builders of the result types construct.
//
// Ledger has a method for each row of the table of default results: a value
// of a struct with a builder, a pointer to one, a value beside a flag, an
// instance of a generic struct with a builder, a struct without a builder and
// a struct whose builder is in another package. The test of the double of
// Ledger checks the default results of each method.
//
// The builders and the double are in one directory, so the double calls each
// builder in its own package.
//
//dokimi:out path=builttest/
package built

import (
	"context"

	"go.dokimi.dev/dokimi/corpus/go/integration/built/billing"
)

// Account is a struct with a builder, which AccountDefaults seeds.
//
//dokimi:builder
type Account struct {
	ID    string
	Email string
}

// AccountDefaults returns the account that each builder of Account starts
// from.
func AccountDefaults() Account {
	return Account{ID: "acct-1", Email: "owner@example.com"}
}

// Page is a generic struct with a builder, which builds each instance of
// Page.
//
//dokimi:builder
type Page[T any] struct {
	Next  string
	Items []T
}

// Note is a struct without a builder.
type Note struct {
	Text string
}

// Ledger returns a value of each type of the table of default results.
//
//dokimi:stub
type Ledger interface {
	// Get returns an account.
	Get(ctx context.Context, id string) (Account, error)

	// Find returns a pointer to an account.
	Find(ctx context.Context, id string) (*Account, error)

	// Lookup returns an account and whether the ledger has it. The flag
	// keeps its zero value, so an unconfigured call reports a miss.
	Lookup(id string) (Account, bool)

	// List returns a page of accounts.
	List(ctx context.Context) (Page[Account], error)

	// Note returns a struct without a builder, which keeps its zero value.
	Note(id string) Note

	// Bill returns a struct whose builder is in another package. The double
	// refers to a builder of its own package alone, so the invoice keeps its
	// zero value.
	Bill(id string) billing.Invoice
}
