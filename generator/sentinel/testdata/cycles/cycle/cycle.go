// Package cycle declares two error values that wrap each other, which Go
// refuses to initialize. Its tree has no module file, so no build of a module
// compiles it.
//
//dokimi:sentinel
package cycle

import "fmt"

// The error values that wrap each other, and an error value that wraps one
// of them.
var (
	// ErrFirst wraps ErrSecond.
	ErrFirst = fmt.Errorf("%w: first", ErrSecond)

	// ErrSecond wraps ErrFirst.
	ErrSecond = fmt.Errorf("%w: second", ErrFirst)

	// ErrStart wraps ErrFirst.
	ErrStart = fmt.Errorf("%w: start", ErrFirst)
)
