// Package generic is a machine whose replies are generic over the values of
// its calls.
package generic

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of the machine.
type Event interface{ event() }

// Effect is an effect of the machine.
type Effect interface{ effect() }

// Started is the start of the machine, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of the machine.
type Stopped struct{}

func (Started) event() {}
func (Stopped) event() {}

// Ack is the value of the reply to a report.
type Ack struct{ Count int }

// Result is the constraint of the value of a reply.
type Result interface {
	struct{} | Ack
}

// Reply is the reply of the machine to an inbound call, generic over a
// named constraint.
//
//dokimi:reply
type Reply[T Result] struct {
	Value T
	Err   error
	To    ReplyToken
}

// Note is a reply generic over a union that its type parameter states.
//
//dokimi:reply
type Note[T struct{} | int] struct {
	Value T
	Err   error
	To    ReplyToken
}

func (Reply[T]) effect() {}
func (Note[T]) effect()  {}

// Replier is a machine that replies to its start.
//
//dokimi:sim start=Started stop=Stopped
//dokimi:out path=genericsim/
type Replier interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the machine stopped itself.
	Err() error
}

// Node replies to its start.
type Node struct{}

// Step applies ev and appends its effects to out.
func (Node) Step(ev Event, out []Effect) []Effect {
	if s, started := ev.(Started); started {
		return append(out, Reply[struct{}]{To: s.Reply})
	}
	return out
}

// Err returns nil, because the machine never stops itself.
func (Node) Err() error { return nil }
