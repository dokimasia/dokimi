// Package hosted is a member that keeps the value of each key in the storage
// of its host, and reads the address of a peer from a registry that every
// member shares.
package hosted

import (
	"context"
	"errors"
)

// ErrUnavailable is the error of a call to a storage that is down.
var ErrUnavailable = errors.New("hosted: the storage is unavailable")

// ErrStopped is the error of a call that the stop of its member ends.
var ErrStopped = errors.New("hosted: stopped")

// Storage keeps the values of one member.
type Storage interface {
	// Write writes value under key.
	//
	//dokimi:fault errors=[ErrUnavailable]
	Write(ctx context.Context, key, value string) error

	// Read returns the value under key.
	//
	//dokimi:fault errors=[ErrUnavailable]
	Read(ctx context.Context, key string) (string, error)
}

// Registry keeps the address of each member.
type Registry interface {
	// Lookup returns the address of the member named name.
	Lookup(ctx context.Context, name string) (string, error)
}

// Token links an effect to the event of its result.
//
//dokimi:token
type Token uint64

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of a member.
type Event interface{ event() }

// Effect is an effect of a member.
type Effect interface{ effect() }

// Started is the start of a member, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of a member.
type Stopped struct{}

// Put is an inbound call that writes Value under Key.
type Put struct {
	Key   string
	Value string
	Reply ReplyToken
}

// Get is an inbound call that returns the value under Key.
type Get struct {
	Key   string
	Reply ReplyToken
}

// Find is an inbound call that returns the address of the member Name.
type Find struct {
	Name  string
	Reply ReplyToken
}

// Await is an inbound call that waits for the next put of the member, and
// returns its value.
type Await struct{ Reply ReplyToken }

// Persisted is the result of the persist Token, and the event of its
// failure.
type Persisted struct {
	Err   error
	Token Token
}

// Loaded is the result of the load Token: the value under its key, or the
// error of the read.
type Loaded struct {
	Err   error
	Value string
	Token Token
}

// Found is the result of the lookup Token: the address of the member, or the
// error of the lookup.
type Found struct {
	Err     error
	Address string
	Token   Token
}

func (Started) event()   {}
func (Stopped) event()   {}
func (Put) event()       {}
func (Get) event()       {}
func (Find) event()      {}
func (Await) event()     {}
func (Persisted) event() {}
func (Loaded) event()    {}
func (Found) event()     {}

// Reply is the reply of a member to an inbound call.
//
//dokimi:reply
type Reply struct {
	Value any
	Err   error
	To    ReplyToken
}

// Persist writes Value under Key through the storage of the host of the
// member.
//
//dokimi:perform func=persist dependency=Storage result=Persisted scope=host
type Persist struct {
	Key   string
	Value string
	Token Token
}

// Load reads the value under Key through the storage of the host of the
// member.
//
//dokimi:perform call=Storage.Read result=Loaded scope=host
type Load struct {
	Key   string
	Token Token
}

// Lookup reads the address of the member Name from the registry.
//
//dokimi:perform call=Registry.Lookup result=Found scope=world
type Lookup struct {
	Name  string
	Token Token
}

func (Reply) effect()   {}
func (Persist) effect() {}
func (Load) effect()    {}
func (Lookup) effect()  {}

// Machine is the machine of a member.
//
//dokimi:sim start=Started stop=Stopped tolerate=[crash:1, outage:1]
//dokimi:out path=hostedsim/
type Machine interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error with which the member stopped itself.
	Err() error
}

// Member replies to each inbound call with the result of its effect. Each
// await returns the value of the next put that succeeds. On its stop it
// replies to each call that waits.
type Member struct {
	waiting map[Token]waiting
	awaits  []ReplyToken
	tokens  Token
}

// waiting is a call that waits for the result of its effect: the reply token
// of the call, and the value of a put.
type waiting struct {
	value string
	reply ReplyToken
}

// New returns a member without a call.
func New() *Member {
	return &Member{waiting: map[Token]waiting{}}
}

// Step applies ev and appends its effects to out.
func (m *Member) Step(ev Event, out []Effect) []Effect {
	switch ev := ev.(type) {
	case Started:
		return append(out, Reply{To: ev.Reply})
	case Put:
		return append(out, Persist{Key: ev.Key, Value: ev.Value, Token: m.call(ev.Reply, ev.Value)})
	case Get:
		return append(out, Load{Key: ev.Key, Token: m.call(ev.Reply, "")})
	case Find:
		return append(out, Lookup{Name: ev.Name, Token: m.call(ev.Reply, "")})
	case Await:
		m.awaits = append(m.awaits, ev.Reply)
		return out
	case Persisted:
		call, waits := m.waiting[ev.Token]
		out = m.reply(out, ev.Token, nil, ev.Err)
		if !waits || ev.Err != nil {
			return out
		}
		for _, reply := range m.awaits {
			out = append(out, Reply{Value: call.value, To: reply})
		}
		m.awaits = nil
		return out
	case Loaded:
		return m.reply(out, ev.Token, ev.Value, ev.Err)
	case Found:
		return m.reply(out, ev.Token, ev.Address, ev.Err)
	case Stopped:
		for token := Token(1); token <= m.tokens; token++ {
			out = m.reply(out, token, nil, ErrStopped)
		}
		for _, reply := range m.awaits {
			out = append(out, Reply{Err: ErrStopped, To: reply})
		}
		m.awaits = nil
		return out
	default:
		return out
	}
}

// call returns a new token of an effect for the inbound call of the reply
// token reply, with the value of a put.
func (m *Member) call(reply ReplyToken, value string) Token {
	m.tokens++
	m.waiting[m.tokens] = waiting{value: value, reply: reply}
	return m.tokens
}

// reply appends the reply to the call that waits for the result of the
// token to out, and returns out alone where no call waits for it.
func (m *Member) reply(out []Effect, token Token, value any, err error) []Effect {
	call, waits := m.waiting[token]
	if !waits {
		return out
	}
	delete(m.waiting, token)
	return append(out, Reply{Value: value, Err: err, To: call.reply})
}

// Err returns nil, because a member never stops itself.
func (*Member) Err() error { return nil }
