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

// Hidden is a subject whose factory its package does not export, so the
// external test package cannot call it.
//
//dokimi:subject store.Store new=newHidden
type Hidden struct{}

// newHidden returns a new Hidden.
func newHidden() *Hidden { return &Hidden{} }

// Priceless is a subject that excuses a check of Get that the suite of Store
// does not have.
//
//dokimi:subject store.Store new=NewPriceless
//dokimi:excuse Store.Get/keeps-the-price reason="the store has no prices"
type Priceless struct{}

// NewPriceless returns a Priceless.
func NewPriceless() *Priceless { return &Priceless{} }

// Misplaced is a subject that excuses a check of Get under the callable Put.
//
//dokimi:subject store.Store new=NewMisplaced
//dokimi:excuse Store.Put/keeps-the-item reason="the store drops each item"
type Misplaced struct{}

// NewMisplaced returns a Misplaced.
func NewMisplaced() *Misplaced { return &Misplaced{} }

// Fetchless is a subject that excuses a check of a callable that Store does
// not have.
//
//dokimi:subject store.Store new=NewFetchless
//dokimi:excuse Store.Fetch/finds-the-item reason="the store cannot fetch"
type Fetchless struct{}

// NewFetchless returns a Fetchless.
func NewFetchless() *Fetchless { return &Fetchless{} }

// Configured is a subject whose function of the options takes a parameter.
//
//dokimi:subject store.Store new=NewConfigured options=Configure
type Configured struct{}

// NewConfigured returns a Configured.
func NewConfigured() *Configured { return &Configured{} }

// Configure returns the options of the level.
func Configure(level int) []any { return []any{level} }

// Doubled is a subject whose function of the options returns two results.
//
//dokimi:subject store.Store new=NewDoubled options=DoubledOptions
type Doubled struct{}

// NewDoubled returns a Doubled.
func NewDoubled() *Doubled { return &Doubled{} }

// DoubledOptions returns the options and an error.
func DoubledOptions() ([]any, error) { return nil, nil }

// Single is a subject whose function of the options returns one option.
//
//dokimi:subject store.Store new=NewSingle options=SingleOption
type Single struct{}

// NewSingle returns a Single.
func NewSingle() *Single { return &Single{} }

// SingleOption returns one option.
func SingleOption() any { return nil }

// Secret is a subject whose function of the options its package does not
// export, so the external test package cannot call it.
//
//dokimi:subject store.Store new=NewSecret options=secretOptions
type Secret struct{}

// NewSecret returns a Secret.
func NewSecret() *Secret { return &Secret{} }

// secretOptions returns no option.
func secretOptions() []any { return nil }
