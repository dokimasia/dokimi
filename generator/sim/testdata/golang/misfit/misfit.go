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

// Ticket links an effect to the event of its result, as Token does.
//
//dokimi:token
type Ticket uint64

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

// Ended is the end of a wait.
type Ended struct{ Token Token }

func (Started) event()   {}
func (Stopped) event()   {}
func (Stored) event()    {}
func (Untokened) event() {}
func (Errorless) event() {}
func (Valueless) event() {}
func (Ended) event()     {}

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

// Watch waits for the items of a feed.
//
//dokimi:perform func=watch
type Watch struct{ Token Token }

// EndsCall ends a wait beside a call.
//
//dokimi:perform ends=[Watch] call=Store.Put result=Stored
type EndsCall struct{ Token Token }

// EndsDependency ends a wait beside a dependency.
//
//dokimi:perform ends=[Watch] func=watch dependency=Store
type EndsDependency struct{ Token Token }

// EndsResult ends a wait beside a result.
//
//dokimi:perform ends=[Watch] result=Stored
type EndsResult struct{ Token Token }

// EndsTimeout ends a wait beside a timeout.
//
//dokimi:perform ends=[Watch] timeout=Timeout
type EndsTimeout struct{ Token Token }

// EndsEnd ends a wait beside an end event.
//
//dokimi:perform ends=[Watch] end=Ended
type EndsEnd struct{ Token Token }

// Endless names an end event without a dependency.
//
//dokimi:perform func=watch end=Ended
type Endless struct{ Token Token }

// Untied ends a wait without a field of a token.
//
//dokimi:perform ends=[Watch]
type Untied struct{}

// Misended ends the wait of an event.
//
//dokimi:perform ends=[Stored]
type Misended struct{ Token Token }

// Unaligned ends the wait of an effect without a token.
//
//dokimi:perform ends=[Tokenless]
type Unaligned struct{ Token Token }

// Unended names an end event without a field of the token.
//
//dokimi:perform func=unend dependency=Store result=Stored end=Untokened
type Unended struct{ Token Token }

// Ticketed waits for the items of a feed under a ticket.
//
//dokimi:perform func=ticket
type Ticketed struct{ Ticket Ticket }

// Mistyped ends the wait of an effect whose token has another type.
//
//dokimi:perform ends=[Ticketed]
type Mistyped struct{ Token Token }

// Pointed names a function that takes a pointer to another effect.
//
//dokimi:perform func=point
type Pointed struct{ Token Token }

// Misvalued names a function that takes another effect.
//
//dokimi:perform func=misvalue
type Misvalued struct{ Token Token }

// Unscoped names a scope for a function without a dependency.
//
//dokimi:perform func=watch scope=host
type Unscoped struct{ Token Token }

// EndsScope ends a wait beside a scope.
//
//dokimi:perform ends=[Watch] scope=host
type EndsScope struct{ Token Token }

// Hosted puts through the store of its host. It is the first effect of the
// machine whose performer binds.
//
//dokimi:perform call=Store.Put result=Stored scope=host
type Hosted struct {
	Key   string
	Value string
	Token Token
}

// Shared puts through the store that every host shares.
//
//dokimi:perform call=Store.Put result=Stored
type Shared struct {
	Key   string
	Value string
	Token Token
}

// Zoned puts through the store of its host, as Hosted does. Its name sorts
// after Watch, an effect without a dependency, so the check of its scope
// meets an effect without one.
//
//dokimi:perform call=Store.Put result=Stored scope=host
type Zoned struct {
	Key   string
	Value string
	Token Token
}

// helper is no method of a dependency.
func helper(context.Context) error { return nil }

// watch performs Watch.
func watch(*sim.World[Event, Effect], *sim.Process[Event, Effect], Watch) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

// unend performs Unended.
func unend(*sim.World[Event, Effect], *sim.Process[Event, Effect], Unended) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

// point takes a pointer to Watch in place of Pointed.
func point(*sim.World[Event, Effect], *sim.Process[Event, Effect], *Watch) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

// misvalue takes Watch in place of Misvalued.
func misvalue(*sim.World[Event, Effect], *sim.Process[Event, Effect], Watch) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

// ticket performs Ticketed.
func ticket(*sim.World[Event, Effect], *sim.Process[Event, Effect], Ticketed) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

// stray performs Strayed.
func stray(*sim.World[Event, Effect], *sim.Process[Event, Effect], Strayed) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

// unmark performs Unmarked.
func unmark(*sim.World[Event, Effect], *sim.Process[Event, Effect], Unmarked) sim.Outcome[Event, Effect] {
	return sim.Outcome[Event, Effect]{}
}

func (Twice) effect()          {}
func (Neither) effect()        {}
func (Resultless) effect()     {}
func (Keyless) effect()        {}
func (Lost) effect()           {}
func (Silent) effect()         {}
func (Blind) effect()          {}
func (Unhosted) effect()       {}
func (Tokenless) effect()      {}
func (Routed) effect()         {}
func (Failless) effect()       {}
func (Orphan) effect()         {}
func (Strayed) effect()        {}
func (Unmarked) effect()       {}
func (Watch) effect()          {}
func (EndsCall) effect()       {}
func (EndsDependency) effect() {}
func (EndsResult) effect()     {}
func (EndsTimeout) effect()    {}
func (EndsEnd) effect()        {}
func (Endless) effect()        {}
func (Untied) effect()         {}
func (Misended) effect()       {}
func (Unaligned) effect()      {}
func (Unended) effect()        {}
func (Ticketed) effect()       {}
func (Mistyped) effect()       {}
func (Pointed) effect()        {}
func (Misvalued) effect()      {}
func (Unscoped) effect()       {}
func (EndsScope) effect()      {}
func (Hosted) effect()         {}
func (Shared) effect()         {}
func (Zoned) effect()          {}

// Machine is the machine of a node.
//
//dokimi:sim start=Started stop=Stopped
type Machine interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the node stopped itself.
	Err() error
}
