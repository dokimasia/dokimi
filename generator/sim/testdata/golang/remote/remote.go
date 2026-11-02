// Package remote is a register whose nodes store each write through a
// writer.
package remote

import (
	"context"
	"errors"
	"time"

	"example.com/shop/remote/store"
)

// ErrFull is the error of a put to a writer without room.
var ErrFull = errors.New("remote: full")

// ErrStopped is the error of a call that the stop of its node ends.
var ErrStopped = errors.New("remote: stopped")

// Writer stores values under keys.
type Writer interface {
	// Put stores value under key.
	//
	//dokimi:fault errors=[ErrFull]
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

// Started is the start of a node, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of a node.
type Stopped struct{}

// Write is an inbound call that writes Value under Key.
type Write struct {
	Key   string
	Value string
	Reply ReplyToken
}

// Stored is the result of the put Token.
type Stored struct {
	Err   error
	Token Token
}

// Fetched is the result of the get Token.
type Fetched struct {
	Value string
	Err   error
	Token Token
}

// Watch is an inbound call that waits for the stop of the node.
type Watch struct{ Reply ReplyToken }

// Loaded is the result of the load Token.
type Loaded struct {
	Value store.Value
	Err   error
	Token Token
}

func (Started) event() {}
func (Stopped) event() {}
func (Write) event()   {}
func (Stored) event()  {}
func (Fetched) event() {}
func (Watch) event()   {}
func (Loaded) event()  {}

// Reply is the reply of a node to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Put stores Value under Key through the writer of the node, within Timeout.
//
//dokimi:perform call=Writer.Put result=Stored timeout=Timeout
type Put struct {
	Key     string
	Value   string
	Token   Token
	Timeout time.Duration
}

// Get reads the value under Key through the writer of the node.
//
//dokimi:perform call=Writer.Get result=Fetched
type Get struct {
	Key   string
	Token Token
}

// Load reads the value under Key through the reader of the store.
//
//dokimi:perform call=store.Reader.Read result=Loaded
type Load struct {
	Key   string
	Token Token
}

func (Reply) effect() {}
func (Put) effect()   {}
func (Get) effect()   {}
func (Load) effect()  {}

// Machine is the machine of a node.
//
//dokimi:sim start=Started stop=Stopped
//dokimi:out path=remotesim/
type Machine interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the node stopped itself.
	Err() error
}

// Node is a node of the register. It puts each write and reads it back
// before it replies with the value that it read. On its stop it replies to
// each waiting write and to each watch.
type Node struct {
	waiting  map[Token]ReplyToken
	watching []ReplyToken
	tokens   Token
}

// New returns a node.
func New() *Node { return &Node{waiting: map[Token]ReplyToken{}} }

// Step applies ev and appends its effects to out.
func (n *Node) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Write:
		n.tokens += 2
		n.waiting[n.tokens] = ev.Reply
		return append(out, Put{Key: ev.Key, Value: ev.Value, Token: n.tokens - 1, Timeout: time.Second},
			Get{Key: ev.Key, Token: n.tokens})
	case Fetched:
		reply, waits := n.waiting[ev.Token]
		if !waits {
			return out
		}
		delete(n.waiting, ev.Token)
		return append(out, Reply{Value: ev.Value, Err: ev.Err, To: reply})
	case Watch:
		n.watching = append(n.watching, ev.Reply)
		return out
	case Stopped:
		for token := Token(1); token <= n.tokens; token++ {
			if reply, waits := n.waiting[token]; waits {
				out = append(out, Reply{Err: ErrStopped, To: reply})
			}
		}
		for _, reply := range n.watching {
			out = append(out, Reply{Err: ErrStopped, To: reply})
		}
		n.watching = nil
		return out
	default:
		return out
	}
}

// Err returns nil, because a node never stops itself.
func (*Node) Err() error { return nil }
