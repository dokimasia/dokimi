package catalog

import (
	"errors"
	"fmt"
)

// The error values of the catalog.
var (
	// ErrNoAccount reports an account that does not exist.
	ErrNoAccount = errors.New("catalog: no account")
	// ErrFrozen reports a change of a frozen account.
	ErrFrozen = errors.New("catalog: account frozen")
)

// LimitError reports an account over its limit.
type LimitError struct {
	// Owner is the owner of the account.
	Owner string
}

// Error returns the message of the error.
func (e *LimitError) Error() string { return fmt.Sprintf("catalog: account of %q over its limit", e.Owner) }
