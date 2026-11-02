// Package silent has a machine without a reply, whose simulation has a call
// of a client.
package silent

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

// Asker is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped
//dokimi:out path=silentsim/
type Asker interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}

// Mute is a machine whose package has no simulation of its own.
//
//dokimi:sim start=Started stop=Stopped
type Mute interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
