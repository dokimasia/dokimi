// Package scoped declares error values and error types. Its external test
// package declares at package scope the names of the packages that their
// checks import.
//
//dokimi:sentinel
package scoped

import (
	"errors"
	"fmt"
)

// The error values of the package.
var (
	// ErrMissing reports a key that was never written.
	ErrMissing = errors.New("scoped: missing")
	// ErrStale reports a read of an old version.
	ErrStale = errors.New("scoped: stale")
)

// MissingError reports a missing key, and compares itself with ErrMissing.
type MissingError struct {
	// Key is the missing key.
	Key string
}

// Error returns the message of the error.
func (e *MissingError) Error() string { return fmt.Sprintf("scoped: key %q missing", e.Key) }

// Is reports the error as a kind of ErrMissing.
func (*MissingError) Is(target error) bool { return target == ErrMissing }

// WrappedError wraps the cause of a failed operation.
type WrappedError struct {
	// Cause is the cause of the failure.
	Cause error
}

// Error returns the message of the error.
func (e *WrappedError) Error() string { return fmt.Sprintf("scoped: %v", e.Cause) }

// Unwrap returns the cause.
func (e *WrappedError) Unwrap() error { return e.Cause }
