// Package scoped declares at package scope the names of the packages that the
// surfaces of Status and Level import.
package scoped

// The variables that take the names of the packages of the surfaces.
var (
	errors  = "the error of a text that matches no variant"
	fmt     = "the text of a value outside the variants"
	strconv = "the text of a number"
)

// Status is the state of an order.
//
//dokimi:enum
type Status int

// The states of an order.
const (
	Draft Status = iota + 1
	Published
)

// Level is a level of a lamp, a float.
//
//dokimi:enum
type Level float64

// The levels of a lamp.
const (
	Dim    Level = 0.5
	Bright Level = 1.5
)
