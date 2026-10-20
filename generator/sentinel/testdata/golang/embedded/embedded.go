// Package embedded declares error types whose hooks an embed promotes.
//
//dokimi:sentinel
package embedded

import (
	"errors"
	"fmt"
)

// ErrNotFound is the error value that the embedded hooks compare with.
var ErrNotFound = errors.New("embedded: not found")

// BaseError has every hook, and other error types embed it.
type BaseError struct {
	// Op is the operation that failed.
	Op string
	// Cause is the cause of the failure.
	Cause error
}

// Error returns the message of the error.
func (e BaseError) Error() string { return fmt.Sprintf("embedded: %s: %v", e.Op, e.Cause) }

// Is reports whether the cause matches target.
func (e BaseError) Is(target error) bool { return errors.Is(e.Cause, target) }

// Unwrap returns the cause.
func (e BaseError) Unwrap() error { return e.Cause }

// NotFoundError declares nothing, and its embed promotes every hook.
type NotFoundError struct {
	BaseError
}

// LostError embeds its base by pointer. The base of its zero value is nil,
// so a check cannot assign the cause through the embed.
type LostError struct {
	*lostBase
}

// lostBase is the base of LostError, and its hooks accept a nil base.
type lostBase struct {
	// Cause is the cause of the failure.
	Cause error
}

// Error returns the message of the error.
func (e *lostBase) Error() string { return "embedded: lost" }

// Is reports whether the cause matches target.
func (e *lostBase) Is(target error) bool { return e != nil && errors.Is(e.Cause, target) }

// Unwrap returns the cause.
func (e *lostBase) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
