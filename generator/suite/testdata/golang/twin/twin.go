// Package twin declares a second Store under the directive suite.
package twin

import "context"

// Store keeps the names of items by their identity.
//
//dokimi:suite
type Store interface {
	// Get returns the name of the item with the identity id.
	Get(ctx context.Context, id string) (string, error)
}
