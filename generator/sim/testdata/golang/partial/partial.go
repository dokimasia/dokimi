// Package partial has a reply without the field of its value, and effects
// whose performers the package of the simulation lacks.
package partial

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

// Reply is a reply whose directive names a field that it lacks.
//
//dokimi:reply value=Result
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Flush is an effect whose performer is not in the package of the
// simulation.
//
//dokimi:perform func=flush
type Flush struct{}

// Wrong is an effect whose performer takes a world in place of a process.
//
//dokimi:perform func=wrong
type Wrong struct{}

// Pair is an effect whose performer takes the effect Wrong.
//
//dokimi:perform func=pair
type Pair struct{}

// Odd is an effect whose performer returns an outcome of the package of the
// simulation.
//
//dokimi:perform func=odd
type Odd struct{}

// Skew is an effect whose performer returns an outcome whose events are the
// effects.
//
//dokimi:perform func=skew
type Skew struct{}

func (Reply) effect() {}
func (Flush) effect() {}
func (Wrong) effect() {}
func (Pair) effect()  {}
func (Odd) effect()   {}
func (Skew) effect()  {}

// Flusher is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped
//dokimi:out path=partialsim/
type Flusher interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
