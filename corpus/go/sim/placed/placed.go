// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package placed is a system of cells that the agent of each host starts and
// stops on its host. Each cell writes its values to a store that every cell
// shares. The agents and the cells are sans-IO machines.
package placed

import (
	"context"
	"errors"
	"slices"
	"time"
)

// The errors of the calls of the clients.
var (
	// ErrPlaced is the error of an add of a cell that the agent placed.
	ErrPlaced = errors.New("placed: the cell is placed")
	// ErrUnplaced is the error of a drop of a cell that the agent has not
	// placed.
	ErrUnplaced = errors.New("placed: the cell is not placed")
	// ErrStopped is the error of a call that the stop of its machine ends.
	ErrStopped = errors.New("placed: the machine stopped")
	// ErrUnavailable is the error of a write to a store that is down.
	ErrUnavailable = errors.New("placed: the store is unavailable")
)

// Store keeps the writes of the cells. Every cell of the system writes to
// it.
type Store interface {
	// Write writes value for cell at the time at.
	//
	//dokimi:fault errors=[ErrUnavailable]
	Write(ctx context.Context, cell, value string, at time.Time) error
}

// Token links an effect to the event of its result.
//
//dokimi:token
type Token uint64

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of an agent or of a cell.
type Event interface{ event() }

// Effect is an effect of an agent or of a cell.
type Effect interface{ effect() }

// Started is the start of an agent, an inbound call. It is also the first
// event of a cell, the inbound call of the agent that starts the cell.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of an agent or of a cell.
type Stopped struct{}

// Add is an inbound call of a client that makes the agent start Cell.
type Add struct {
	Cell  string
	Reply ReplyToken
}

// Drop is an inbound call of a client that makes the agent stop Cell.
type Drop struct {
	Cell  string
	Reply ReplyToken
}

// Set is an inbound call of a client that makes a cell write Value.
type Set struct {
	Value string
	Reply ReplyToken
}

// Await is an inbound call of a client that waits for the next write of a
// cell, and returns its value.
type Await struct{ Reply ReplyToken }

// Placed is the result of the start Token of a cell.
type Placed struct {
	Err   error
	Token Token
}

// Unplaced is the result of the stop Token of a cell.
type Unplaced struct{ Token Token }

// Written is the result of the write Token, and the event of its failure.
type Written struct {
	Err   error
	Token Token
}

func (Started) event()  {}
func (Stopped) event()  {}
func (Add) event()      {}
func (Drop) event()     {}
func (Set) event()      {}
func (Await) event()    {}
func (Placed) event()   {}
func (Unplaced) event() {}
func (Written) event()  {}

// Reply is the reply of an agent or of a cell to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Place starts Cell on the host of the agent. The writes of the cell have
// the timeout Timeout.
//
//dokimi:perform func=place
type Place struct {
	Cell    string
	Token   Token
	Timeout time.Duration
}

// Unplace stops Cell.
//
//dokimi:perform func=unplace
type Unplace struct {
	Cell  string
	Token Token
}

// Write writes Value to the store, within Timeout.
//
//dokimi:perform func=write dependency=Store result=Written timeout=Timeout
type Write struct {
	Value   string
	Token   Token
	Timeout time.Duration
}

func (Reply) effect()   {}
func (Place) effect()   {}
func (Unplace) effect() {}
func (Write) effect()   {}

// Node is the machine of an agent and of a cell.
//
// +dokimi:sim start=Started stop=Stopped require=[crash:0.05] \
// tolerate=[crash:1, restart:1, pause:1, outage:1, slow:1, lose-after:1, clock-jump:1]
//
//dokimi:out path=placedsim/
type Node interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the machine stopped itself, and nil
	// while it runs.
	Err() error
}

// call is a call that waits for the result of an effect: the token of the
// effect, the reply token of the call, the cell of an add or a drop, and the
// value of a set.
type call struct {
	cell  string
	value string
	token Token
	reply ReplyToken
}

// Agent starts and stops the cells of its host. Each cell that it starts
// writes with the timeout of the agent. On its stop it replies to each call
// that waits.
type Agent struct {
	placed  []string
	calls   []call
	tokens  Token
	timeout time.Duration
}

// NewAgent returns an agent whose cells write with the timeout timeout.
func NewAgent(timeout time.Duration) *Agent { return &Agent{timeout: timeout} }

// Step applies ev and appends its effects to out.
func (a *Agent) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Add:
		if slices.Contains(a.placed, ev.Cell) || a.waits(ev.Cell) {
			return append(out, Reply{Err: ErrPlaced, To: ev.Reply})
		}
		a.tokens++
		a.calls = append(a.calls, call{cell: ev.Cell, token: a.tokens, reply: ev.Reply})
		return append(out, Place{Cell: ev.Cell, Token: a.tokens, Timeout: a.timeout})
	case Drop:
		if !slices.Contains(a.placed, ev.Cell) || a.waits(ev.Cell) {
			return append(out, Reply{Err: ErrUnplaced, To: ev.Reply})
		}
		a.tokens++
		a.calls = append(a.calls, call{cell: ev.Cell, token: a.tokens, reply: ev.Reply})
		return append(out, Unplace{Cell: ev.Cell, Token: a.tokens})
	case Placed:
		c, waits := a.take(ev.Token)
		if !waits {
			return out
		}
		if ev.Err == nil {
			a.placed = append(a.placed, c.cell)
		}
		return append(out, Reply{Err: ev.Err, To: c.reply})
	case Unplaced:
		c, waits := a.take(ev.Token)
		if !waits {
			return out
		}
		a.placed = slices.DeleteFunc(a.placed, func(cell string) bool { return cell == c.cell })
		return append(out, Reply{To: c.reply})
	case Stopped:
		for _, c := range a.calls {
			out = append(out, Reply{Err: ErrStopped, To: c.reply})
		}
		a.calls = nil
		return out
	default:
		return out
	}
}

// Err returns nil, because an agent never stops itself.
func (*Agent) Err() error { return nil }

// waits reports whether a call of the agent waits on the cell.
func (a *Agent) waits(cell string) bool {
	return slices.ContainsFunc(a.calls, func(c call) bool { return c.cell == cell })
}

// take removes the call that waits for the effect token, and reports whether
// one waits.
func (a *Agent) take(token Token) (call, bool) {
	at := slices.IndexFunc(a.calls, func(c call) bool { return c.token == token })
	if at < 0 {
		return call{}, false
	}
	c := a.calls[at]
	a.calls = slices.Delete(a.calls, at, at+1)
	return c, true
}

// Cell writes each value of its calls to the store before it replies. Its
// writes have its timeout. Each await returns the value of the next write
// that succeeds. On its stop it replies to each call that waits.
type Cell struct {
	calls   []call
	awaits  []ReplyToken
	tokens  Token
	timeout time.Duration
}

// NewCell returns a cell whose writes have the timeout timeout.
func NewCell(timeout time.Duration) *Cell { return &Cell{timeout: timeout} }

// Step applies ev and appends its effects to out.
func (c *Cell) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Set:
		c.tokens++
		c.calls = append(c.calls, call{value: ev.Value, token: c.tokens, reply: ev.Reply})
		return append(out, Write{Value: ev.Value, Token: c.tokens, Timeout: c.timeout})
	case Await:
		c.awaits = append(c.awaits, ev.Reply)
		return out
	case Written:
		at := slices.IndexFunc(c.calls, func(w call) bool { return w.token == ev.Token })
		if at < 0 {
			return out
		}
		w := c.calls[at]
		c.calls = slices.Delete(c.calls, at, at+1)
		out = append(out, Reply{Err: ev.Err, To: w.reply})
		if ev.Err != nil {
			return out
		}
		for _, reply := range c.awaits {
			out = append(out, Reply{Value: w.value, To: reply})
		}
		c.awaits = nil
		return out
	case Stopped:
		for _, w := range c.calls {
			out = append(out, Reply{Err: ErrStopped, To: w.reply})
		}
		for _, reply := range c.awaits {
			out = append(out, Reply{Err: ErrStopped, To: reply})
		}
		c.calls, c.awaits = nil, nil
		return out
	default:
		return out
	}
}

// Err returns nil, because a cell never stops itself.
func (*Cell) Err() error { return nil }
