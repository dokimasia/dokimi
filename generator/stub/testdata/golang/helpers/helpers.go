// Package helpers declares an interface whose stub of a method would have
// two helpers of one name.
package helpers

import "errors"

// The error values of Fetch, whose base names are both Gone.
var (
	// ErrGone reports an item that the store deleted.
	ErrGone = errors.New("helpers: gone")
	// Gone reports the same under a name that the rules of Go do not read as
	// the name of an error value.
	Gone = errors.New("helpers: gone again")
)

// Fetcher has a method whose directive fault lists two error values with one
// base name, so both helpers are FaultGone.
//
//dokimi:stub
type Fetcher interface {
	// Fetch can fail.
	//
	//dokimi:fault errors=[ErrGone, Gone]
	Fetch() error
}
