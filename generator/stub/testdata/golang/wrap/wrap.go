// Package wrap wraps the tools of the tree lib, whose internal package it
// cannot import.
package wrap

import "example.com/shop/lib/tool"

// Wrapped is a struct under the directive stub that promotes Level, whose
// result is a type of the internal package of the tree lib.
//
//dokimi:stub
type Wrapped struct {
	*tool.Tool
}

// Measured is an interface under the directive stub that embeds Gauge, whose
// method Level returns a type of the internal package of the tree lib.
//
//dokimi:stub
type Measured interface {
	tool.Gauge
}

// Adjusted is a struct under the directive stub that promotes Set, whose
// parameter is a type of the internal package of the tree lib.
//
//dokimi:stub
type Adjusted struct {
	*tool.Knob
}
