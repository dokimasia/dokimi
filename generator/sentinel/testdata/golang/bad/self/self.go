// Package self lists itself as a neighbour.
//
//dokimi:sentinel no-overlap-with=[example.com/shop/bad/self]
package self

import "errors"

// ErrSelf matches itself.
var ErrSelf = errors.New("self: matched")
