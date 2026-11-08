// Package waits is a leader that campaigns for a term and watches a feed,
// and ends both waits.
package waits

import (
	"context"
	"errors"
)

// Token links an effect to the events of its result.
//
//dokimi:token
type Token uint64

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// ErrUnavailable is the error of a store that is down.
var ErrUnavailable = errors.New("waits: the store is unavailable")

// Store is the store of the campaigns.
type Store interface {
	// Campaign campaigns for a term of the scope.
	//
	//dokimi:fault errors=[ErrUnavailable]
	Campaign(ctx context.Context, scope string) error
}

// Event is an event of the leader.
type Event interface{ event() }

// Effect is an effect of the leader.
type Effect interface{ effect() }

// Started is the start of the leader, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of the leader.
type Stopped struct{}

// Campaigned is the grant of a term to a campaign, or the failure of the
// campaign.
type Campaigned struct {
	Err   error
	Token Token
	Term  uint64
}

// TermEnded is the end of the term of a campaign.
type TermEnded struct{ Token Token }

// Watched is an item of a watch.
type Watched struct {
	Err   error
	Token Token
}

func (Started) event()    {}
func (Stopped) event()    {}
func (Campaigned) event() {}
func (TermEnded) event()  {}
func (Watched) event()    {}

// Reply is the reply to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Campaign campaigns for a term of the scope.
//
//dokimi:perform func=campaign dependency=Store result=Campaigned end=TermEnded
type Campaign struct {
	Scope string
	Token Token
}

// Resign ends the campaign Token.
//
//dokimi:perform ends=[Campaign] func=resign
type Resign struct{ Token Token }

// Watch watches the feed of the leader.
//
//dokimi:perform func=watch
type Watch struct{ Token Token }

// Unwatch ends the watch Token.
//
//dokimi:perform ends=[Watch]
type Unwatch struct{ Token Token }

func (Reply) effect()    {}
func (Campaign) effect() {}
func (Resign) effect()   {}
func (Watch) effect()    {}
func (Unwatch) effect()  {}

// Leader is the machine of the leader.
//
//dokimi:sim start=Started stop=Stopped tolerate=[outage:1, lose-after:1] mean-steps=1 max-steps=200
//dokimi:out path=waitssim/
type Leader interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the leader stopped itself.
	Err() error
}

// leader campaigns once it starts, and replies to its start on its stop. It
// resigns and watches the feed once its campaign has a term, and unwatches
// once its watch delivers an item.
type leader struct {
	tokens   Token
	campaign Token
	watch    Token
	start    ReplyToken
}

// New returns a leader.
func New() Leader { return &leader{} }

// Step applies ev.
func (l *leader) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		l.tokens++
		l.start, l.campaign = ev.Reply, l.tokens
		return append(out, Campaign{Scope: "s", Token: l.campaign})
	case Stopped:
		return append(out, Reply{To: l.start})
	case Campaigned:
		if ev.Token == l.campaign && ev.Err == nil {
			l.tokens++
			l.watch = l.tokens
			return append(out, Resign{Token: l.campaign}, Watch{Token: l.watch})
		}
	case Watched:
		if ev.Token == l.watch && ev.Err == nil {
			return append(out, Unwatch{Token: l.watch})
		}
	}
	return out
}

// Err returns nil, because the leader does not stop itself.
func (*leader) Err() error { return nil }
