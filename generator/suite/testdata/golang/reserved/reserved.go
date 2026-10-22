// Package reserved binds a check of a method with the segment of a check of
// the test of a suite.
package reserved

import "context"

// Shelf keeps the names of items by their identity.
//
//dokimi:suite
type Shelf interface {
	// Get returns the name of the item with the identity id.
	Get(ctx context.Context, id string) (string, error)
}

// BuildsTheSubject checks that Get returns the name of the item.
//
//dokimi:check Shelf.Get defect=empty
func BuildsTheSubject(h any) {}

// empty returns a shelf without items.
func empty(any) Shelf { return nil }
