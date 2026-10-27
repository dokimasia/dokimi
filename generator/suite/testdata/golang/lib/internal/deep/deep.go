// Package deep declares the level of a tool, which only the packages of the
// tree lib import.
package deep

import "errors"

// ErrWorn reports a worn tool.
var ErrWorn = errors.New("deep: the tool is worn")

// Level is the level of a tool, with a declared default.
//
//dokimi:default 3
type Level int

// High is the highest level of a tool. It is untyped, so Level is a defined
// type without members.
const High = 9

// Wide is the widest pitch of a tool.
const Wide = 12

// Mode is the mode of a tool.
type Mode int

// The modes of a tool.
const (
	Manual Mode = iota + 1
	Auto
)

// Spec is the spec of a tool, with a declared default depth.
type Spec struct {
	Depth int `default:"1"`
}

// Probe probes a tool.
type Probe interface {
	// Depth returns the depth of the probe.
	Depth() int
}

// Hook runs at the level of a tool.
type Hook func() error

// Shape is the shape of a tool. Its method area is not exported, so the
// package deep alone implements it.
type Shape interface {
	// Area returns the area of the shape.
	Area() int
	area()
}

// Square is the shape of a square tool.
type Square struct {
	Side int
}

// Area returns the area of the square.
func (s Square) Area() int { return s.Side * s.Side }

func (Square) area() {}
