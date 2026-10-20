// Package remote declares an interface that embeds an interface of another
// package.
package remote

import "example.com/shop/other"

// Remote embeds other.Thing. A load without the package other leaves a gap
// in its members.
//
//dokimi:stub
type Remote interface {
	other.Thing
}
