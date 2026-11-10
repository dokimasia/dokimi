// Package defectless is a machine whose simulation has an invariant and a
// convergence without a planted defect.
package defectless

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

// Undefended is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped
//dokimi:out path=defectlesssim/
type Undefended interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
