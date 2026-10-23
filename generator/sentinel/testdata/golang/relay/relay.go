// Package relay relays the calls of the shop to a remote store. It
// re-exports an error value of store, and carries the errors of the remote
// store in CallError.
//
//dokimi:sentinel remote=[CallError]
package relay

import (
	"errors"

	"example.com/shop/store"
)

// The error values of the relay.
var (
	// ErrNotFound re-exports the error value of store, so a caller matches it
	// without an import of store.
	ErrNotFound = store.ErrNotFound
	// ErrDropped reports a call that the relay dropped.
	ErrDropped = errors.New("relay: dropped")
)

// CallError carries the error of the remote store in place of the result of
// a call, with the message of the remote store unchanged.
type CallError struct {
	// Message is the message of the remote error.
	Message string
	// Code is the class of the remote error.
	Code int
}

// Error returns the message of the remote error.
func (e *CallError) Error() string { return e.Message }

// TimeoutError reports a call that the remote store did not answer in time.
type TimeoutError struct {
	// Op is the operation that timed out.
	Op string
}

// Error returns the message of the error.
func (e *TimeoutError) Error() string { return "relay: " + e.Op + " timed out" }
