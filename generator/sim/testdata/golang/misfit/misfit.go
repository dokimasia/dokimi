// Package misfit declares effects whose directives perform do not fit.
package misfit

import (
	"context"

	"go.dokimi.dev/sim"
)

// Store stores values under keys.
type Store interface {
	// Put stores value under key.
	Put(ctx context.Context, key, value string) error

	// Get returns the value under key.
	Get(ctx context.Context, key string) (string, error)
}

// Token links an effect to the event of its result.
//
//dokimi:token
type Token uint64

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of a node.
type Event interface{ event() }

// Effect is an effect of a node.
type Effect interface{ effect() }

// Started is the start of a node.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of a node.
type Stopped struct{}

// Stored is the result of a put.
type Stored struct {
	Err   error
	Token Token
}

// Untokened is a result without a field of the token.
type Untokened struct{ Err error }

// Errorless is a result without a field of the error.
type Errorless struct{ Token Token }

// Valueless is the result of a get without a field of the value.
type Valueless struct {
	Err   error
	Token Token
}

func (Started) event()   {}
func (Stopped) event()   {}
func (Stored) event()    {}
func (Untokened) event() {}
func (Errorless) event() {}
func (Valueless) event() {}

// Twice names both a function and a call.
//
//dokimi:perform func=twice call=Store.Put result=Stored
type Twice struct {
	Key   string
	Value string
	Token Token
}

// Neither names neither a function nor a call.
//
//dokimi:perform timeout=Timeout
type Neither struct{ Token Token }

// Resultless names a call without a result.
//
//dokimi:perform call=Store.Put
type Resultless struct {
	Key   string
	Value string
	Token Token
}

// Keyless has no field for the key of the put.
//
//dokimi:perform call=Store.Put result=Stored
type Keyless struct {
	Value string
	Token Token
}

// Lost puts into an event without a field of the token.
//
//dokimi:perform call=Store.Put result=Untokened
type Lost struct {
	Key   string
	Value string
	Token Token
}

// Silent puts into an event without a field of the error.
//
//dokimi:perform call=Store.Put result=Errorless
type Silent struct {
	Key   string
	Value string
	Token Token
}

// Blind gets into an event without a field of the value.
//
//dokimi:perform call=Store.Get result=Valueless
type Blind struct {
	Key   string
	Token Token
}

// Unhosted names a function of the package in place of a method.
//
//dokimi:perform call=helper result=Stored
type Unhosted struct{ Token Token }

// Tokenless has no field of a token for the event of its result.
//
//dokimi:perform call=Store.Put result=Stored
type Tokenless struct {
	Key   string
	Value string
}

// Routed names a dependency beside a call.
//
//dokimi:perform call=Store.Put dependency=Store result=Stored
type Routed struct {
	Key   string
	Value string
	Token Token
}

// Failless names a dependency without the event of its failure.
//
//dokimi:perform func=perform dependency=Store
type Failless struct{ Token Token }

// Orphan names the event of a failure without a dependency.
//
//dokimi:perform func=perform result=Stored
type Orphan struct{ Token Token }

// Strayed names an event as its dependency.
//
//dokimi:perform func=stray dependency=Stored result=Stored
type Strayed struct{ Token Token }

// Unmarked has no field of a token for the event of its failure.
//
//dokimi:perform func=unmark dependency=Store result=Stored
type Unmarked struct{ Key string }

// helper is no method of a dependency.
func helper(context.Context) error { return nil }

// stray performs Strayed.
func stray(*sim.World[Event, Effect], *sim.Process[Event, Effect], Strayed) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

// unmark performs Unmarked.
func unmark(*sim.World[Event, Effect], *sim.Process[Event, Effect], Unmarked) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

func (Twice) effect()      {}
func (Neither) effect()    {}
func (Resultless) effect() {}
func (Keyless) effect()    {}
func (Lost) effect()       {}
func (Silent) effect()     {}
func (Blind) effect()      {}
func (Unhosted) effect()   {}
func (Tokenless) effect()  {}
func (Routed) effect()     {}
func (Failless) effect()   {}
func (Orphan) effect()     {}
func (Strayed) effect()    {}
func (Unmarked) effect()   {}

// Machine is the machine of a node.
//
//dokimi:sim start=Started stop=Stopped
type Machine interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the node stopped itself.
	Err() error
}
