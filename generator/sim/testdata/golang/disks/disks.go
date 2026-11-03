// Package disks is a machine whose simulation tolerates the faults of the
// media of its processes.
package disks

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

// Flush is an effect that syncs the medium of a process.
//
//dokimi:perform func=flush
type Flush struct{}

func (Flush) effect() {}

// Disks is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped tolerate=[torn-write:1, corruption:1, full-disk:1] require=[full-disk:0.5]
//dokimi:out path=diskssim/
type Disks interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
