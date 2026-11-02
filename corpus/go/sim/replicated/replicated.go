// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package replicated is a register of keys whose primary copies each write
// to a backup before it applies and acknowledges the write. Its nodes are
// sans-IO machines.
package replicated

import (
	"errors"
	"slices"
	"time"
)

// ErrStopped is the error of a call that the stop of its node ends.
var ErrStopped = errors.New("replicated: the node stopped")

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

// Started is the start of a node, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of a node.
type Stopped struct{}

// Cancelled ends the wait of the caller of the inbound call Reply.
type Cancelled struct{ Reply ReplyToken }

// Put is an inbound call of a client that writes Value under Key.
type Put struct {
	Key   string
	Value string
	Reply ReplyToken
}

// Get is an inbound call of a client that reads Key.
type Get struct {
	Key   string
	Reply ReplyToken
}

// Watch is an inbound call of a client that waits for the next write of Key.
type Watch struct {
	Key   string
	Reply ReplyToken
}

// Copy is the inbound call of the primary that copies Value under Key to
// the backup.
type Copy struct {
	Key   string
	Value string
	Reply ReplyToken
}

// Copied is the result of the copy Token.
type Copied struct {
	Err   error
	Token Token
}

func (Started) event()   {}
func (Stopped) event()   {}
func (Cancelled) event() {}
func (Put) event()       {}
func (Get) event()       {}
func (Watch) event()     {}
func (Copy) event()      {}
func (Copied) event()    {}

// Reply is the reply of a node to an inbound call. The backup can hold the
// copy of a write that the stop of the primary ends, so the outcome of that
// write is unknown.
//
//dokimi:reply unknown=[ErrStopped]
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Send copies Value under Key to the backup, within Timeout.
//
//dokimi:perform func=send timeout=Timeout
type Send struct {
	Key     string
	Value   string
	Token   Token
	Timeout time.Duration
}

func (Reply) effect() {}
func (Send) effect()  {}

// Node is the machine of a node of the register.
//
// +dokimi:sim start=Started stop=Stopped cancel=Cancelled require=[crash:0.05] \
// tolerate=[crash:1, restart:1, partition:1, pause:1, drop:1, delay:1, duplicate:1, clock-jump:1]
//
//dokimi:out path=replicatedsim/
type Node interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the node stopped itself, and nil while
	// it runs.
	Err() error
}

// Disk is the durable store of a node. A crash of the node keeps it.
type Disk map[string]string

// Primary is the primary of the register. It copies each write to the
// backup, and applies the write to its disk and acknowledges it once the
// backup has it. It serves the reads and the watches from its disk.
type Primary struct {
	disk    Disk
	writes  []write
	watches []watch
	tokens  Token
}

// write is a write that waits for its copy: its key and value, the token of
// its copy and the reply token of its call.
type write struct {
	key   string
	value string
	token Token
	reply ReplyToken
}

// watch is a watch that waits for the next write of its key.
type watch struct {
	key   string
	reply ReplyToken
}

// NewPrimary returns a primary that applies its writes to disk.
func NewPrimary(disk Disk) *Primary { return &Primary{disk: disk} }

// Step applies ev and appends its effects to out.
func (p *Primary) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Put:
		p.tokens++
		p.writes = append(p.writes, write{key: ev.Key, value: ev.Value, token: p.tokens, reply: ev.Reply})
		return append(out, Send{Key: ev.Key, Value: ev.Value, Token: p.tokens, Timeout: time.Second})
	case Copied:
		at := slices.IndexFunc(p.writes, func(w write) bool { return w.token == ev.Token })
		if at < 0 {
			return out
		}
		w := p.writes[at]
		p.writes = slices.Delete(p.writes, at, at+1)
		if ev.Err != nil {
			return append(out, Reply{Err: ev.Err, To: w.reply})
		}
		p.disk[w.key] = w.value
		out = append(out, Reply{To: w.reply})
		for _, watched := range p.watches {
			if watched.key == w.key {
				out = append(out, Reply{Value: w.value, To: watched.reply})
			}
		}
		p.watches = slices.DeleteFunc(p.watches, func(watched watch) bool { return watched.key == w.key })
		return out
	case Get:
		return append(out, Reply{Value: p.disk[ev.Key], To: ev.Reply})
	case Watch:
		p.watches = append(p.watches, watch{key: ev.Key, reply: ev.Reply})
		return out
	case Cancelled:
		p.writes = slices.DeleteFunc(p.writes, func(w write) bool { return w.reply == ev.Reply })
		p.watches = slices.DeleteFunc(p.watches, func(w watch) bool { return w.reply == ev.Reply })
		return out
	case Stopped:
		for _, w := range p.writes {
			out = append(out, Reply{Err: ErrStopped, To: w.reply})
		}
		for _, w := range p.watches {
			out = append(out, Reply{Err: ErrStopped, To: w.reply})
		}
		p.writes, p.watches = nil, nil
		return out
	default:
		return out
	}
}

// Err returns nil, because a primary never stops itself.
func (*Primary) Err() error { return nil }

// Backup is the backup of the register. It applies each copy to its disk and
// replies to it, and serves the reads from its disk.
type Backup struct {
	disk Disk
}

// NewBackup returns a backup that applies its copies to disk.
func NewBackup(disk Disk) *Backup { return &Backup{disk: disk} }

// Step applies ev and appends its effects to out.
func (b *Backup) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Copy:
		b.disk[ev.Key] = ev.Value
		return append(out, Reply{To: ev.Reply})
	case Get:
		return append(out, Reply{Value: b.disk[ev.Key], To: ev.Reply})
	default:
		return out
	}
}

// Err returns nil, because a backup never stops itself.
func (*Backup) Err() error { return nil }
