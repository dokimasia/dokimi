// Package bad holds subjects that the generator reports.
package bad

import "example.com/shop/store"

// Lister lists items, and no directive suite marks it.
type Lister interface {
	// List returns every item.
	List() []store.Item
}

// Plain is a Lister of a declaration without a suite.
//
//dokimi:subject Lister new=NewPlain
type Plain struct{}

// NewPlain returns a Plain.
func NewPlain() *Plain { return &Plain{} }

// List returns no item.
func (*Plain) List() []store.Item { return nil }

// Twice is a subject whose factory takes two parameters.
//
//dokimi:subject store.Store new=NewTwice
type Twice struct{}

// NewTwice returns a Twice.
func NewTwice(_, _ string) *Twice { return &Twice{} }

// Paired is a subject whose factory returns two results.
//
//dokimi:subject store.Store new=NewPaired
type Paired struct{}

// NewPaired returns a Paired.
func NewPaired() (*Paired, error) { return &Paired{}, nil }

// Nothing is a subject whose factory returns no result.
//
//dokimi:subject store.Store new=NewNothing
type Nothing struct{}

// NewNothing returns nothing.
func NewNothing() {}
