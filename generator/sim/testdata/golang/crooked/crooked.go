// Package crooked has a machine whose simulation has functions with the
// signatures of other directives.
package crooked

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

// Reply is the reply of a machine to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

func (Reply) effect() {}

// Crooked is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped
//dokimi:out path=crookedsim/
type Crooked interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
