// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ledger is a ledger of lists. A transaction appends to lists and
// reads them, and its node applies it at once and replies once it has
// persisted it. The node is a sans-IO machine.
package ledger

import (
	"errors"
	"slices"
)

// ErrStopped is the error of a call that the stop of the node ends.
var ErrStopped = errors.New("ledger: the node stopped")

// Token links an effect to the event of its result.
//
//dokimi:token
type Token uint64

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of the node.
type Event interface{ event() }

// Effect is an effect of the node.
type Effect interface{ effect() }

// Started is the start of the node, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of the node.
type Stopped struct{}

// Cancelled ends the wait of the caller of the inbound call Reply.
type Cancelled struct{ Reply ReplyToken }

// Txn is an inbound call of a client that applies Ops in one transaction.
// Each operation is a list-append micro-operation of assert-go:
// ["append", key, value] or ["read", key, nil].
type Txn struct {
	Ops   []any
	Reply ReplyToken
}

// Persisted is the result of the persist Token.
type Persisted struct{ Token Token }

func (Started) event()   {}
func (Stopped) event()   {}
func (Cancelled) event() {}
func (Txn) event()       {}
func (Persisted) event() {}

// Reply is the reply of the node to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Persist persists the transaction Token.
//
//dokimi:perform func=persist
type Persist struct{ Token Token }

func (Reply) effect()   {}
func (Persist) effect() {}

// Ledger is the machine of the node of the ledger.
//
// +dokimi:sim start=Started stop=Stopped cancel=Cancelled isolation=serializable \
// tolerate=[crash:1, restart:1, pause:1, clock-jump:1]
//
//dokimi:out path=ledgersim/
type Ledger interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the node stopped itself, and nil while
	// it runs.
	Err() error
}

// Disk is the durable store of the lists. A crash of the node keeps it.
type Disk map[string][]any

// Node is the node of the ledger. It applies each transaction at once, and
// replies to it once it has persisted it. A transaction without operations
// waits for the stop of the node.
type Node struct {
	disk    Disk
	pending []pending
	waiting []ReplyToken
	tokens  Token
}

// pending is a transaction that waits for its persist: the token of the
// persist, the reply token of its call and its completed operations.
type pending struct {
	ops   []any
	token Token
	reply ReplyToken
}

// NewNode returns a node whose lists are on disk.
func NewNode(disk Disk) *Node { return &Node{disk: disk} }

// Step applies ev and appends its effects to out.
func (n *Node) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Txn:
		if len(ev.Ops) == 0 {
			n.waiting = append(n.waiting, ev.Reply)
			return out
		}
		n.tokens++
		n.pending = append(n.pending, pending{ops: Apply(n.disk, ev.Ops), token: n.tokens, reply: ev.Reply})
		return append(out, Persist{Token: n.tokens})
	case Persisted:
		at := slices.IndexFunc(n.pending, func(p pending) bool { return p.token == ev.Token })
		if at < 0 {
			return out
		}
		p := n.pending[at]
		n.pending = slices.Delete(n.pending, at, at+1)
		return append(out, Reply{Value: p.ops, To: p.reply})
	case Cancelled:
		n.pending = slices.DeleteFunc(n.pending, func(p pending) bool { return p.reply == ev.Reply })
		n.waiting = slices.DeleteFunc(n.waiting, func(r ReplyToken) bool { return r == ev.Reply })
		return out
	case Stopped:
		for _, p := range n.pending {
			out = append(out, Reply{Value: p.ops, To: p.reply})
		}
		for _, r := range n.waiting {
			out = append(out, Reply{Err: ErrStopped, To: r})
		}
		n.pending, n.waiting = nil, nil
		return out
	default:
		return out
	}
}

// Err returns nil, because the node never stops itself.
func (*Node) Err() error { return nil }

// Apply applies the micro-operations ops to the lists of disk in order, and
// returns them completed: each read with the list that it read.
func Apply(disk Disk, ops []any) []any {
	out := make([]any, 0, len(ops))
	for _, op := range ops {
		micro, _ := op.([]any)
		kind, _ := micro[0].(string)
		key, _ := micro[1].(string)
		if kind == "append" {
			disk[key] = append(disk[key], micro[2])
			out = append(out, micro)
			continue
		}
		out = append(out, []any{kind, key, append([]any{}, disk[key]...)})
	}
	return out
}
