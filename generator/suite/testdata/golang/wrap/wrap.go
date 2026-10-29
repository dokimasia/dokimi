// Package wrap wraps the tools of the tree lib, whose internal package it
// cannot import.
package wrap

import (
	"context"

	"example.com/shop/lib/tool"
)

// Reach stores and reads the pitch of a tool, whose declared default is a
// constant of the internal package of the tree lib.
//
//dokimi:suite
type Reach interface {
	// Store stores p.
	Store(ctx context.Context, p tool.Pitch) error

	// Current returns the pitch that Store stored.
	//
	//dokimi:mixin read-after-write write=Store
	Current(ctx context.Context) (tool.Pitch, error)

	// Gear returns the gear of the tool.
	Gear(ctx context.Context) (tool.Gear, error)
}

// Wrapped is a struct under the directive suite that promotes Level, whose
// result is a type of the internal package of the tree lib.
//
//dokimi:suite
type Wrapped struct {
	*tool.Tool
}

// Measured is an interface under the directive suite that embeds Gauge, whose
// method Level returns a type of the internal package of the tree lib.
//
//dokimi:suite
type Measured interface {
	tool.Gauge
}

// Adjusted is a struct under the directive suite that promotes Set, whose
// parameter is a type of the internal package of the tree lib.
//
//dokimi:suite
type Adjusted struct {
	*tool.Knob
}

// Tuned is an interface under the directive suite whose callables take types
// of the tree lib. Their fields and their methods refer to the internal
// package of the tree.
//
//dokimi:suite
type Tuned interface {
	// Tune tunes the tool of the request.
	Tune(r tool.Request) error
	// Watch watches the gauge.
	Watch(g tool.Gauge)
	// Turn turns the tool with the dial.
	Turn(d tool.Dial)
}

// Limited is an interface under the directive suite whose callables take
// values that refer to the internal package of the tree lib, or to
// declarations that the package tool does not export.
//
//dokimi:suite
type Limited interface {
	// Limit limits the tool.
	Limit(l tool.Limits)
	// Grip grips the tool.
	Grip(g tool.Grip)
	// Clamp clamps the tool.
	Clamp(c tool.Clamp)
	// Frame frames the tool.
	Frame(f tool.Frame)
	// Mount mounts the tool.
	Mount(m tool.Mount)
	// Span spans the tool.
	Span(s tool.Span)
	// Pitch pitches the tool.
	Pitch(p tool.Pitch)
	// Measure measures the tool in the unit.
	Measure(u tool.Unit)
	// Phase moves the tool to the phase.
	Phase(p tool.Phase)
	// Sense senses the tool with the sensor.
	Sense(s tool.Sensor) error
}

// Plain is a struct under the directive suite in the same file, whose
// callables refer to no internal package.
//
//dokimi:suite
type Plain struct {
	name string
}

// Name returns the name.
func (p *Plain) Name() string { return p.name }
