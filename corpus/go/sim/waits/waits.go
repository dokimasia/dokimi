// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package waits is a leader that campaigns for the term of a scope and
// watches the terms of the scope through a store. A call of a client makes
// the leader step down: it resigns its campaign and ends its watch. Another
// call makes it campaign and watch again. The leader is a sans-IO machine.
package waits

import (
	"context"
	"errors"
)

// The errors of the leaders and of their store.
var (
	// ErrUnavailable is the error of a campaign at a store that is down.
	ErrUnavailable = errors.New("waits: the store is unavailable")
	// ErrStopped is the error of a call that the stop of its leader ends.
	ErrStopped = errors.New("waits: the leader stopped")
)

// Store keeps the candidacies of the sessions of the leaders for the terms of
// the scopes. A session is one incarnation of a leader.
type Store interface {
	// Campaign enters a candidacy of the session for the term of the scope.
	//
	//dokimi:fault errors=[ErrUnavailable]
	Campaign(ctx context.Context, scope, session string) (Candidacy, error)
	// Term returns the term of the scope, and 0 for a scope without a term.
	Term(scope string) uint64
}

// Candidacy is the candidacy of a campaign in the store.
type Candidacy interface {
	// Term returns the term that the store granted to the candidacy, and 0
	// before the grant.
	Term() uint64
	// Resign ends the candidacy, and the term that the store granted to it.
	Resign()
}

// Token links an effect to the events of its result.
//
//dokimi:token
type Token uint64

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of a leader.
type Event interface{ event() }

// Effect is an effect of a leader.
type Effect interface{ effect() }

// Started is the start of a leader, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of a leader.
type Stopped struct{}

// StepDown is an inbound call of a client that makes the leader resign its
// campaign and end its watch. Its reply has the term of the leader.
type StepDown struct{ Reply ReplyToken }

// Rejoin is an inbound call of a client that makes the leader campaign and
// watch again.
type Rejoin struct{ Reply ReplyToken }

// Lead is an inbound call of a client that waits until the leader has a term,
// and returns the term.
type Lead struct{ Reply ReplyToken }

// Campaigned is the grant of the term Term to the campaign Token, or the
// failure of the campaign. A Campaigned with an error is the last event of
// its campaign.
type Campaigned struct {
	Err   error
	Token Token
	Term  uint64
}

// TermEnded is the end of the term of the campaign Token, and the last event
// of the campaign.
type TermEnded struct{ Token Token }

// Watched is the term of the scope of the watch Token, after a change.
type Watched struct {
	Token Token
	Term  uint64
}

func (Started) event()    {}
func (Stopped) event()    {}
func (StepDown) event()   {}
func (Rejoin) event()     {}
func (Lead) event()       {}
func (Campaigned) event() {}
func (TermEnded) event()  {}
func (Watched) event()    {}

// Reply is the reply of a leader to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Campaign campaigns for the term of Scope. Its wait delivers the grant of a
// term, and the end of the term.
//
//dokimi:perform func=campaign dependency=Store result=Campaigned end=TermEnded
type Campaign struct {
	Scope string
	Token Token
}

// Resign ends the campaign Token.
//
//dokimi:perform ends=[Campaign]
type Resign struct{ Token Token }

// Watch watches the terms of Scope.
//
//dokimi:perform func=watch
type Watch struct {
	Scope string
	Token Token
}

// Unwatch ends the watch Token.
//
//dokimi:perform ends=[Watch]
type Unwatch struct{ Token Token }

func (Reply) effect()    {}
func (Campaign) effect() {}
func (Resign) effect()   {}
func (Watch) effect()    {}
func (Unwatch) effect()  {}

// Machine is the machine of a leader.
//
//dokimi:sim start=Started stop=Stopped tolerate=[outage:1, lose-after:1]
//dokimi:out path=waitssim/
type Machine interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the machine stopped itself, and nil
	// while it runs.
	Err() error
}

// Leader campaigns for the term of its scope and watches the terms of the
// scope once it starts. A step down resigns its campaign and ends its watch,
// and a rejoin campaigns and watches again. The leader does not campaign
// again after the failure or the end of a campaign. A lead waits until the
// leader has a term, and on its stop the leader replies to each lead that
// waits.
type Leader struct {
	scope    string
	leads    []ReplyToken
	tokens   Token
	campaign Token
	watch    Token
	term     uint64
}

// NewLeader returns a leader of the scope.
func NewLeader(scope string) *Leader { return &Leader{scope: scope} }

// Step applies ev and appends its effects to out.
func (l *Leader) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return l.join(append(out, Reply{To: ev.Reply}))
	case Rejoin:
		return l.join(append(out, Reply{To: ev.Reply}))
	case Lead:
		if l.term != 0 {
			return append(out, Reply{Value: l.term, To: ev.Reply})
		}
		l.leads = append(l.leads, ev.Reply)
		return out
	case StepDown:
		out = append(out, Reply{Value: l.term, To: ev.Reply})
		if l.campaign != 0 {
			out = append(out, Resign{Token: l.campaign})
		}
		if l.watch != 0 {
			out = append(out, Unwatch{Token: l.watch})
			l.watch = 0
		}
		return out
	case Campaigned:
		switch {
		case ev.Token != l.campaign:
		case ev.Err != nil:
			l.campaign = 0
		default:
			l.term = ev.Term
			for _, reply := range l.leads {
				out = append(out, Reply{Value: l.term, To: reply})
			}
			l.leads = nil
		}
		return out
	case TermEnded:
		if ev.Token == l.campaign {
			l.campaign, l.term = 0, 0
		}
		return out
	case Stopped:
		for _, reply := range l.leads {
			out = append(out, Reply{Err: ErrStopped, To: reply})
		}
		l.leads = nil
		return out
	default:
		return out
	}
}

// Err returns nil, because a leader never stops itself.
func (*Leader) Err() error { return nil }

// join campaigns where the leader has no campaign, and watches where it has
// no watch.
func (l *Leader) join(out []Effect) []Effect {
	if l.campaign == 0 {
		l.tokens++
		l.campaign = l.tokens
		out = append(out, Campaign{Scope: l.scope, Token: l.campaign})
	}
	if l.watch == 0 {
		l.tokens++
		l.watch = l.tokens
		out = append(out, Watch{Scope: l.scope, Token: l.watch})
	}
	return out
}
