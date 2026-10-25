// Package deep declares the level of a tool, which only the packages of the
// tree lib import.
package deep

// Level is the level of a tool, with a declared default.
//
//dokimi:default 3
type Level int

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
