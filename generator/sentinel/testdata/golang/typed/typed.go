// Package typed declares error types and no error value.
//
//dokimi:sentinel
package typed

import "fmt"

// LimitError reports an exceeded limit through a value receiver.
type LimitError struct {
	// Name is the name of the limit.
	Name string
	// count is the number of the attempts, which a check does not write.
	count int
}

// Error returns the message of the error.
func (e LimitError) Error() string {
	return fmt.Sprintf("typed: limit %s exceeded after %d attempts", e.Name, e.count)
}

// RangeError reports a value outside a range of any type, and has no checks.
type RangeError[T any] struct {
	// Value is the value outside the range.
	Value T
}

// Error returns the message of the error.
func (e RangeError[T]) Error() string { return fmt.Sprintf("typed: %v out of range", e.Value) }
