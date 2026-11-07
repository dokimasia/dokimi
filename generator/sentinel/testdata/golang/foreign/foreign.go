// Package foreign reads the policies of the store. Its error value ErrMissing
// wraps the error value ErrNotFound of the package store, which is no parent
// in the family of foreign, also where foreign declares an ErrNotFound of its
// own.
//
//dokimi:sentinel prefix=off
package foreign

import (
	"errors"
	"fmt"

	"example.com/shop/store"
)

// The error values of the policies.
var (
	// ErrMissing is returned for a policy that the store does not have.
	ErrMissing = fmt.Errorf("%w: no policy", store.ErrNotFound)

	// ErrStale is returned for a policy that the store replaced.
	ErrStale = errors.New("foreign: stale policy")

	// ErrNotFound is returned for a procedure without a policy.
	ErrNotFound = errors.New("foreign: not found")
)
