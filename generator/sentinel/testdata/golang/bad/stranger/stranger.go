// Package stranger lists a remote type that is no error type of the package.
//
//dokimi:sentinel remote=[Point]
package stranger

import "errors"

// ErrLost reports a lost point.
var ErrLost = errors.New("stranger: lost")

// Point is a point on a plane, which is no error.
type Point struct {
	// X is the distance from the vertical axis.
	X int
}
