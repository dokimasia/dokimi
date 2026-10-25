// Package wrap wraps the tools of the tree lib, whose internal package it
// cannot import.
package wrap

import "example.com/shop/lib/tool"

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

// Plain is a struct under the directive suite in the same file, whose
// callables refer to no internal package.
//
//dokimi:suite
type Plain struct {
	name string
}

// Name returns the name.
func (p *Plain) Name() string { return p.name }
