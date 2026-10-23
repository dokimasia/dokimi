// Package machine runs a state machine of the shop, whose states the package
// does not export.
package machine

// state is a state of the machine.
//
//dokimi:enum text=comment
type state int

// The states of the machine.
const (
	idle    state = iota // idle
	running              // running
)

// Current returns the text of the first state of the machine.
func Current() string { return idle.String() }

// Next returns the text of the state after the first one.
func Next() string { return running.String() }
