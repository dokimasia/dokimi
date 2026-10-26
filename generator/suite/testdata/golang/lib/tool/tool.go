// Package tool declares a tool whose exported methods return a type of an
// internal package of the tree lib.
package tool

import (
	"errors"

	"example.com/shop/lib/internal/deep"
)

// ErrJammed reports a jammed tool.
var ErrJammed = errors.New("tool: the tool is jammed")

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

// draft is the unit of a setting before its first save, which the package
// does not export.
const draft Unit = 0

// Phase is the phase of a tool, whose members the package does not export.
type Phase int

// The phases of a tool.
const (
	warmup Phase = iota + 1
	steady
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

// Limits are the limits of a tool. The default of Level is a constant of the
// internal package of the tree lib.
type Limits struct {
	Level deep.Level `default:"deep.High"`
	Count int        `default:"2"`
}

// Grip is the grip of a tool. Its authored sample calls a function that the
// package does not export, and its field slack is not exported.
//
//dokimi:sample value=defaultGrip()
type Grip struct {
	Force int `default:"4"`
	slack int `default:"1"`
}

// defaultGrip returns the grip of a new tool.
func defaultGrip() Grip { return Grip{Force: 3} }

// Clamp is the clamp of a tool, which has no exported field.
type Clamp struct {
	width int
}

// brace is a part of a frame, which the package does not export.
type brace struct {
	size int
}

// Frame is the frame of a tool, which embeds a type that the package does not
// export.
type Frame struct {
	brace
}

// Mount is the mount of a tool, which embeds the clamp beside a field that
// the package does not export.
type Mount struct {
	Clamp
	pitch int
}

// Span is the span of a tool. Its authored sample calls a function that the
// package does not export.
//
//dokimi:sample value=defaultSpan()
//dokimi:default 8
type Span int

// defaultSpan returns the span of a new tool.
func defaultSpan() Span { return 6 }

// Pitch is the pitch of a tool, whose declared default is a constant of the
// internal package of the tree lib.
//
//dokimi:default deep.Wide
type Pitch int

// Sensor senses a tool. The errors of its methods are error values of the
// internal package of the tree lib and of this package.
type Sensor interface {
	// Sense senses the tool.
	//
	//dokimi:fault errors=[deep.ErrWorn]
	Sense() error
	// Probe probes the tool.
	//
	//dokimi:fault errors=[deep.ErrWorn, ErrJammed]
	Probe() error
	// Read reads the setting under the key.
	//
	//dokimi:mixin not-found sentinel=deep.ErrWorn
	Read(key string) (int, error)
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
