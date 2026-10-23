// Package alias declares an error value under a second name of its own
// package, which is no re-export.
//
//dokimi:sentinel
package alias

import "errors"

// The error values of the package.
var (
	// ErrBase reports the base failure.
	ErrBase = errors.New("alias: base")
	// ErrSame is ErrBase under a second name.
	ErrSame = ErrBase
)
