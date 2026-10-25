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

// Plain is a struct under the directive suite in the same file, whose
// callables refer to no internal package.
//
//dokimi:suite
type Plain struct {
	name string
}

// Name returns the name.
func (p *Plain) Name() string { return p.name }
