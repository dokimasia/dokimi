// Package faulty has the machine interfaces whose events and effects sim
// refuses: an effect without a performer, a reply without a field of its
// directive, a performer that the package of the simulation lacks, and calls
// without a reply.
package faulty

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of a machine.
type Event interface{ event() }

// Effect is an effect of a machine.
type Effect interface{ effect() }

// Started is the start of a machine, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of a machine.
type Stopped struct{}

func (Started) event() {}
func (Stopped) event() {}

// Sender is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped
type Sender interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// Send is an effect without the directive perform.
type Send struct{ To string }

func (Send) effect() {}
