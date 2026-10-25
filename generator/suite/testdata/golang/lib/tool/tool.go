// Package tool declares a tool whose exported methods return a type of an
// internal package of the tree lib.
package tool

import "example.com/shop/lib/internal/deep"

// Tool is a tool of the shop.
type Tool struct {
	level deep.Level
}

// Name returns the name of the tool.
func (t *Tool) Name() string { return "tool" }

// Level returns the level of the tool.
func (t *Tool) Level() deep.Level { return t.level }

// Gauge is the interface of a tool with a level.
type Gauge interface {
	// Level returns the level of the tool.
	Level() deep.Level
}

// Knob sets the level of a tool.
type Knob struct{}

// Set sets the level of the tool to l.
func (k *Knob) Set(l deep.Level) {}

// Unit is the unit of a setting.
type Unit int

// The units of a setting.
const (
	Percent Unit = iota + 1
	Steps
)

// Setting is a setting of a tool, with a declared default value.
type Setting struct {
	Value int `default:"5"`
}

// Note notes the setting of a tool.
type Note interface {
	// Write writes the text of the note.
	Write(text string) error
}

// Dial dials a tool to the level l.
type Dial func(l deep.Level) error

// Calibrator calibrates a tool.
type Calibrator interface {
	// Calibrate calibrates the tool to the level l.
	Calibrate(l deep.Level) error
}

// Request is the request of a tool. Each type of the internal package of the
// tree lib, and each type whose value refers to one, is on one field alone.
type Request struct {
	Mode       deep.Mode
	Unit       Unit
	Spec       deep.Spec
	Setting    Setting
	Level      deep.Level
	Probe      deep.Probe
	Note       Note
	Read       func() deep.Level
	Hook       deep.Hook
	Dial       Dial
	Gauge      Gauge
	Calibrator Calibrator
}
