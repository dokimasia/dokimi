// Package loose declares the logs of the shop, whose calls take values of any
// type and pointers that the property engine reads as values.
package loose

import (
	"math/big"
	"time"
)

// Printer prints messages with values of any type.
//
//dokimi:suite
type Printer interface {
	// Print prints the message with its values.
	Print(msg string, values ...any)
}

// Recorder records the values of fields.
//
//dokimi:suite
type Recorder interface {
	// Field records the value of the field.
	Field(name string, value interface{})
}

// Entry is an entry of a journal, whose value has any type.
type Entry struct {
	// Value is the value of the entry.
	Value any
}

// Journal keeps entries whose values have any type.
//
//dokimi:suite
type Journal interface {
	// Put keeps the entry.
	Put(e Entry)
}

// Plain prints messages alone.
//
//dokimi:suite
type Plain interface {
	// Print prints the message.
	Print(msg string)
}

// Zoned is a clock with the pointers that the property engine reads as values.
//
//dokimi:suite
type Zoned interface {
	// In moves the clock into the zone loc.
	In(loc *time.Location)
	// Scale scales the clock by n.
	Scale(n *big.Int)
	// Ratio sets the ratio of the clock to r.
	Ratio(r *big.Rat)
}
