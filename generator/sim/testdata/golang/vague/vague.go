// Package vague has the generic replies whose instantiations the anatomy
// cannot list.
package vague

// ReplyToken links an inbound call to the reply to it.
//
//dokimi:token reply
type ReplyToken uint64

// Event is an event of the machine.
type Event interface{ event() }

// Effect is an effect of the machine.
type Effect interface{ effect() }

// Started is the start of the machine, an inbound call.
type Started struct{ Reply ReplyToken }

// Stopped is the stop of the machine.
type Stopped struct{}

func (Started) event() {}
func (Stopped) event() {}

// Sized is a constraint with a method and no union.
type Sized interface {
	// Size returns the size of the value.
	Size() int
}

// Loose is a reply generic over every type.
//
//dokimi:reply
type Loose[T any] struct {
	Value T
	Err   error
	To    ReplyToken
}

// Rough is a reply generic over a union with an approximation.
//
//dokimi:reply
type Rough[T ~int | string] struct {
	Value T
	Err   error
	To    ReplyToken
}

// Twin is a reply generic over two type parameters.
//
//dokimi:reply
type Twin[A struct{} | int, B struct{} | int] struct {
	Value A
	Other B
	Err   error
	To    ReplyToken
}

// Bare is a reply generic over a named constraint without a union.
//
//dokimi:reply
type Bare[T Sized] struct {
	Value T
	Err   error
	To    ReplyToken
}

func (Loose[T]) effect()   {}
func (Rough[T]) effect()   {}
func (Twin[A, B]) effect() {}
func (Bare[T]) effect()    {}

// Vague is the machine of the package.
//
//dokimi:sim start=Started stop=Stopped
//dokimi:out path=vaguesim/
type Vague interface {
	// Step applies ev and appends its effects to out.
	Step(ev Event, out []Effect) []Effect
	// Err returns the error of the machine.
	Err() error
}
