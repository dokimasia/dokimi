// Package store declares error values and error types with each hook. Its
// messages open with kv rather than with the name of the package.
//
//dokimi:sentinel prefix=kv no-overlap-with=[example.com/shop/neighbour]
package store

import (
	"errors"
	"fmt"
)

// The error values of the store.
var (
	// ErrNotFound reports a key that was never written.
	ErrNotFound = errors.New("kv: not found")
	// ErrConflict reports a write that lost a race.
	ErrConflict = errors.New("kv: conflict")
	// ErrClosed reports a use after the close.
	ErrClosed = errors.New("kv: closed")
)

// limit is the most keys that the store holds, and no error value.
var limit = 64

// NotFoundError reports a missing key, and compares itself with ErrNotFound.
type NotFoundError struct {
	// Key is the missing key.
	Key string
	// Attempt is the number of the attempt.
	Attempt int
}

// Error returns the message of the error.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("kv: key %q not found after %d attempts", e.Key, e.Attempt)
}

// Is reports the error as a kind of ErrNotFound.
func (*NotFoundError) Is(target error) bool { return target == ErrNotFound }

// WrappedError wraps the cause of a failed operation.
type WrappedError struct {
	// Op is the operation that failed.
	Op string
	// Cause is the cause of the failure.
	Cause error
}

// Error returns the message of the error.
func (e *WrappedError) Error() string { return fmt.Sprintf("kv: %s: %v", e.Op, e.Cause) }

// Unwrap returns the cause.
func (e *WrappedError) Unwrap() error { return e.Cause }

// note is an unexported error type, which has no checks.
type note struct{}

// Error returns the message of the note.
func (note) Error() string { return "kv: note" }
