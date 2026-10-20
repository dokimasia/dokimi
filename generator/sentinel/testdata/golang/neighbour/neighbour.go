// Package neighbour declares the error values that the error values of the
// store do not match.
//
//dokimi:sentinel
package neighbour

import "errors"

// The error values of the neighbour.
var (
	// ErrNotFound reports a missing key.
	ErrNotFound = errors.New("neighbour: not found")
	// ErrConflict reports a losing write.
	ErrConflict = errors.New("neighbour: conflict")
)
