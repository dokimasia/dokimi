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
