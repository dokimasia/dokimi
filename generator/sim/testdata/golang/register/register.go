// Package register is a register of one value whose nodes are sans-IO
// machines.
package register

import (
	"errors"
	"slices"
	"time"
)

// ErrStopped is the error of a write that the stop of its node ends.
var ErrStopped = errors.New("register: stopped")

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

// Write is an inbound call that writes Value.
type Write struct {
	Value string
	Reply ReplyToken
}

// Committed is the result of the commit Token.
type Committed struct {
	Err   error
	Token Token
}

// Watch is an inbound call that waits for the stop of the node.
type Watch struct{ Reply ReplyToken }

func (Started) event()   {}
func (Stopped) event()   {}
func (Cancelled) event() {}
func (Write) event()     {}
func (Committed) event() {}
func (Watch) event()     {}

// Reply is the reply of a node to an inbound call. A write that the stop of
// its node ends can have committed, so its outcome is unknown.
//
//dokimi:reply unknown=[ErrStopped]
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Commit stores Value on the disk of the node, within Timeout.
//
//dokimi:perform func=commit timeout=Timeout
type Commit struct {
	Value   string
	Token   Token
	Timeout time.Duration
}

func (Reply) effect()  {}
func (Commit) effect() {}

// Machine is the machine of a node.
//
//dokimi:sim start=Started stop=Stopped cancel=Cancelled tolerate=[crash:1, partition:1] require=[crash:0.1] strategy=pct:2
//dokimi:out path=registersim/
type Machine interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the node stopped itself.
	Err() error
}

// Node is a node of the register. It commits each write before it replies
// with the written value. On its stop it replies to each waiting write and
// to each watch.
type Node struct {
	waiting  map[Token]pending
	watching []ReplyToken
	tokens   Token
}

// pending is a write that waits for its commit: its value and its reply
// token.
type pending struct {
	value string
	reply ReplyToken
}

// New returns a node.
func New() *Node { return &Node{waiting: map[Token]pending{}} }

// Step applies ev and appends its effects to out.
func (n *Node) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Write:
		n.tokens++
		n.waiting[n.tokens] = pending{value: ev.Value, reply: ev.Reply}
		return append(out, Commit{Value: ev.Value, Token: n.tokens, Timeout: time.Second})
	case Committed:
		p, waits := n.waiting[ev.Token]
		if !waits {
			return out
		}
		delete(n.waiting, ev.Token)
		return append(out, Reply{Value: p.value, Err: ev.Err, To: p.reply})
	case Watch:
		n.watching = append(n.watching, ev.Reply)
		return out
	case Cancelled:
		for token := Token(1); token <= n.tokens; token++ {
			if p, waits := n.waiting[token]; waits && p.reply == ev.Reply {
				delete(n.waiting, token)
			}
		}
		n.watching = slices.DeleteFunc(n.watching, func(r ReplyToken) bool { return r == ev.Reply })
		return out
	case Stopped:
		for token := Token(1); token <= n.tokens; token++ {
			if p, waits := n.waiting[token]; waits {
				out = append(out, Reply{Err: ErrStopped, To: p.reply})
			}
		}
		for _, r := range n.watching {
			out = append(out, Reply{Err: ErrStopped, To: r})
		}
		n.watching = nil
		return out
	default:
		return out
	}
}

// Err returns nil, because a node never stops itself.
func (*Node) Err() error { return nil }
