// Package sealed declares interfaces under the directive stub with an
// unexported method that their doubles cannot implement.
package sealed

import "example.com/shop/counter"

// Moved is an interface under the directive stub whose double is in the
// package sealedtest, which cannot implement seal.
//
//dokimi:stub
//dokimi:out path=sealedtest/
type Moved interface {
	// Name returns the name.
	Name() string
	// seal limits the implementations of Moved to the package.
	seal()
}

// Named is an interface under the directive stub that embeds Tally of the
// package counter. Only a type of counter implements the method seal of
// Tally.
//
//dokimi:stub
type Named interface {
	counter.Tally
}
