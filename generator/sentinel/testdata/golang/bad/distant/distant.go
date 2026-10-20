// Package distant lists a neighbour that the run does not load.
//
//dokimi:sentinel no-overlap-with=[example.com/shop/missing]
package distant

import "errors"

// ErrDistant is the error value of the package.
var ErrDistant = errors.New("distant: far away")
