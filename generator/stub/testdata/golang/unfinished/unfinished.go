// Package unfinished declares an interface that does not compile, because
// its method returns iter.Seq without a type argument.
package unfinished

import "iter"

// Feed returns a stream without the type of its elements.
//
//dokimi:stub
type Feed interface {
	// Names returns iter.Seq without its type argument.
	Names() iter.Seq
}
