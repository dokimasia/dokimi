// Package placed is an agent that writes the value of each call to a store
// that every agent shares.
package placed

import (
	"context"
	"errors"
	"time"
)

// ErrUnavailable is the error of a call to a store that is down.
var ErrUnavailable = errors.New("placed: the store is unavailable")

// ErrStale is the error of a read that returns a value older than the last
// write.
var ErrStale = errors.New("placed: the value is stale")

// ErrStopped is the error of a call that the stop of its agent ends.
var ErrStopped = errors.New("placed: stopped")

// Store keeps the value of each cell.
type Store interface {
	// Write writes value for cell.
	//
	//dokimi:fault errors=[ErrUnavailable]
	Write(ctx context.Context, cell, value string) error

	// Read returns the value of cell.
	//
	//dokimi:fault errors=[ErrStale, ErrUnavailable]
	Read(ctx context.Context, cell string) (string, error)
}

// Token links an effect to the event of its result.
//
//dokimi:token
type Token uint64

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of an agent.
type Event interface{ event() }

// Effect is an effect of an agent.
type Effect interface{ effect() }

// Started is the start of an agent, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of an agent.
type Stopped struct{}

// Assign is an inbound call that writes Value to the cell of the agent.
type Assign struct {
	Value string
	Reply ReplyToken
}

// Await is an inbound call that waits for the next write of the agent, and
// returns its value.
type Await struct{ Reply ReplyToken }

// Written is the result of the write Token, and the event of its failure.
type Written struct {
	Err   error
	Token Token
}

func (Started) event() {}
func (Stopped) event() {}
func (Assign) event()  {}
func (Await) event()   {}
func (Written) event() {}

// Reply is the reply of an agent to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Write writes Value to the cell of the agent through the store, within
// Timeout.
//
//dokimi:perform func=write dependency=Store result=Written timeout=Timeout
type Write struct {
	Value   string
	Token   Token
	Timeout time.Duration
}

func (Reply) effect() {}
func (Write) effect() {}

// Machine is the machine of an agent.
//
//dokimi:sim start=Started stop=Stopped tolerate=[crash:1, outage:1, lose-after:1]
//dokimi:out path=placedsim/
type Machine interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the agent stopped itself.
	Err() error
}

// Agent writes the value of each call to the store before it replies. Each
// of its writes has the timeout of the agent. Each await returns the value of
// the next write that succeeds. On its stop it replies to each waiting call.
type Agent struct {
	waiting map[Token]assigned
	awaits  []ReplyToken
	tokens  Token
	timeout time.Duration
}

// assigned is an assign that waits for its write: its value and the reply
// token of its call.
type assigned struct {
	value string
	reply ReplyToken
}

// New returns an agent whose writes have the timeout timeout.
func New(timeout time.Duration) *Agent {
	return &Agent{waiting: map[Token]assigned{}, timeout: timeout}
}

// Step applies ev and appends its effects to out.
func (a *Agent) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Assign:
		a.tokens++
		a.waiting[a.tokens] = assigned{value: ev.Value, reply: ev.Reply}
		return append(out, Write{Value: ev.Value, Token: a.tokens, Timeout: a.timeout})
	case Await:
		a.awaits = append(a.awaits, ev.Reply)
		return out
	case Written:
		call, waits := a.waiting[ev.Token]
		if !waits {
			return out
		}
		delete(a.waiting, ev.Token)
		out = append(out, Reply{Err: ev.Err, To: call.reply})
		if ev.Err != nil {
			return out
		}
		for _, reply := range a.awaits {
			out = append(out, Reply{Value: call.value, To: reply})
		}
		a.awaits = nil
		return out
	case Stopped:
		for token := Token(1); token <= a.tokens; token++ {
			if call, waits := a.waiting[token]; waits {
				out = append(out, Reply{Err: ErrStopped, To: call.reply})
			}
		}
		for _, reply := range a.awaits {
			out = append(out, Reply{Err: ErrStopped, To: reply})
		}
		a.awaits = nil
		return out
	default:
		return out
	}
}

// Err returns nil, because an agent never stops itself.
func (*Agent) Err() error { return nil }
