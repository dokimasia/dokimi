// Package bad has the machine interfaces that sim refuses, each for one
// defect of its methods, of its lifecycle or of its params.
package bad

import "go.dokimi.dev/sim"

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Name is a token that is no unsigned integer type.
//
//dokimi:token
type Name string

// Event is an event of a machine.
type Event interface{ event() }

// Effect is an effect of a machine.
type Effect interface{ effect() }

// Started is the start of a machine, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of a machine.
type Stopped struct{}

// Twice is an event with two fields of a reply token.
type Twice struct{ First, Second ReplyToken }

func (Started) event() {}
func (Stopped) event() {}
func (Twice) event()   {}

// Ack is a struct with a reply token that is no event.
type Ack struct{ Reply ReplyToken }

// system spawns no process. It is the system of each machine interface of
// the package, so the simulation of each interface has a system.
//
//dokimi:system
func system(*sim.World[Event, Effect]) {}

// NoStep has no step method.
//
//dokimi:sim start=Started stop=Stopped
type NoStep interface {
	// Err returns the error of the machine.
	Err() error
}

// BadStep has a step method without the slice of the effects.
//
//dokimi:sim start=Started stop=Stopped
type BadStep interface {
	// Step applies ev.
	Step(ev Event) []Effect
	// Err returns the error of the machine.
	Err() error
}

// ListlessStep takes one effect in place of the slice of the effects.
//
//dokimi:sim start=Started stop=Stopped
type ListlessStep interface {
	// Step applies ev and returns its effects.
	Step(ev Event, out Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// ConcreteStep steps on a struct in place of an event interface.
//
//dokimi:sim start=Started stop=Stopped
type ConcreteStep interface {
	// Step applies ev and appends its effects to out.
	Step(ev Started, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// TwoStepReturns has a step method that also returns an error.
//
//dokimi:sim start=Started stop=Stopped
type TwoStepReturns interface {
	// Step applies ev and returns its effects and an error.
	Step(ev Event, out []Effect) ([]Effect, error)
	// Err returns the error of the machine.
	Err() error
}

// PointerReturn has a step method that returns one effect.
//
//dokimi:sim start=Started stop=Stopped
type PointerReturn interface {
	// Step applies ev and returns its last effect.
	Step(ev Event, out []Effect) *Effect
	// Err returns the error of the machine.
	Err() error
}

// MismatchedStep has a step method that returns events in place of effects.
//
//dokimi:sim start=Started stop=Stopped
type MismatchedStep interface {
	// Step applies ev and returns the events that follow from it.
	Step(ev Event, out []Effect) []Event
	// Err returns the error of the machine.
	Err() error
}

// ConcreteEffects appends structs in place of effects of an interface.
//
//dokimi:sim start=Started stop=Stopped
type ConcreteEffects interface {
	// Step applies ev and appends its acks to out.
	Step(ev Event, out []Ack) []Ack
	// Err returns the error of the machine.
	Err() error
}

// StartOutside starts with a struct that is no event.
//
//dokimi:sim start=Ack stop=Stopped
type StartOutside interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// NoErr has no error method.
//
//dokimi:sim start=Started stop=Stopped
type NoErr interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
}

// BadErr has an error method that returns a text.
//
//dokimi:sim start=Started stop=Stopped
type BadErr interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine as a text.
	Err() string
}

// ParamErr has an error method that takes a parameter.
//
//dokimi:sim start=Started stop=Stopped
type ParamErr interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine for the code.
	Err(code int) error
}

// StartWithoutToken starts with an event without a reply token.
//
//dokimi:sim start=Stopped stop=Stopped
type StartWithoutToken interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// StartTwice starts with an event with two reply tokens.
//
//dokimi:sim start=Twice stop=Stopped
type StartTwice interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// StopOutside stops with a struct that is no event.
//
//dokimi:sim start=Started stop=Ack
type StopOutside interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// CancelWithoutToken cancels with an event without a reply token.
//
//dokimi:sim start=Started stop=Stopped cancel=Stopped
type CancelWithoutToken interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// CountlessFault tolerates a fault without a count.
//
//dokimi:sim start=Started stop=Stopped tolerate=[crash]
type CountlessFault interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// UnknownFault tolerates a fault that the world does not have.
//
//dokimi:sim start=Started stop=Stopped tolerate=[meteor:1]
type UnknownFault interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// HugeCount tolerates more crashes at once than an int counts.
//
//dokimi:sim start=Started stop=Stopped tolerate=[crash:99999999999999999999]
type HugeCount interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// ZeroCount tolerates no crash at once.
//
//dokimi:sim start=Started stop=Stopped tolerate=[crash:0]
type ZeroCount interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// WholeShare requires a fault in more than every case.
//
//dokimi:sim start=Started stop=Stopped require=[crash:2]
type WholeShare interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// ZeroShare requires a fault in no case.
//
//dokimi:sim start=Started stop=Stopped require=[crash:0]
type ZeroShare interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// UnknownShare requires a fault that the world does not have.
//
//dokimi:sim start=Started stop=Stopped require=[meteor:0.5]
type UnknownShare interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// ShallowStrategy has a strategy of PCT without depth.
//
//dokimi:sim start=Started stop=Stopped strategy=pct:0
type ShallowStrategy interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// HugeDepth has a strategy of PCT deeper than an int counts.
//
//dokimi:sim start=Started stop=Stopped strategy=pct:99999999999999999999
type HugeDepth interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// UnknownStrategy has a strategy that the world does not have.
//
//dokimi:sim start=Started stop=Stopped strategy=random:3
type UnknownStrategy interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
