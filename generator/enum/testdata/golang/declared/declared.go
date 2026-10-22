// Package declared declares members of the surfaces of its enumerations.
package declared

import "errors"

// ErrUnknownLevel reports a text of no level. The parse function of the
// surface of Level wraps it.
var ErrUnknownLevel = errors.New("declared: unknown level")

// Level is the level of a log record, whose package declares the error of
// its parse function.
//
//dokimi:enum
type Level int

// The levels of a log record.
const (
	Debug Level = iota + 1
	Info
	Warn
)

// Color is a color of a light, whose package declares the list of its
// variants.
//
//dokimi:enum
type Color int

// The colors of a light.
const (
	Red Color = iota + 1
	Green
)

// ColorValues returns the colors in the order of their declaration.
func ColorValues() []Color { return []Color{Red, Green} }

// Lane is a lane of a road, whose package declares its parse function and
// no render.
//
//dokimi:enum
type Lane int

// The lanes of a road.
const (
	Left Lane = iota + 1
	Right
)

// ParseLane returns the lane whose name is s.
func ParseLane(s string) (Lane, error) {
	switch s {
	case "Left":
		return Left, nil
	case "Right":
		return Right, nil
	}
	return 0, errors.New("declared: unknown lane")
}

// Coil is a coil of a spring. A method of another type of its package has the
// name of its parse function.
//
//dokimi:enum
type Coil int

// The coils of a spring.
const (
	Snug Coil = iota + 1
	Slack
)

// Label is the label of a coil.
type Label string

// ParseCoil returns the coil of the label, as a method of a label.
func (l Label) ParseCoil() (Coil, error) {
	if l == "slack" {
		return Slack, nil
	}
	return Snug, nil
}
