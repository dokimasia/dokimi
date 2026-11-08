// Package plain is a machine without a cancel event, whose simulation picks
// its processes uniformly and checks the isolation of its transactions.
package plain

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

// Beta is an effect without a bound of its call.
//
//dokimi:perform func=beta
type Beta struct{}

func (Beta) effect() {}

// Batched is an interface that implements Effect. The anatomy leaves it out,
// because each effect of a machine is a struct.
type Batched interface {
	Effect
	// Batch returns the size of the batch of the effect.
	Batch() int
}

// Plain is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped tolerate=[pause:1] require=[pause:1] strategy=uniform isolation=serializable
//dokimi:out path=plainsim/
type Plain interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
