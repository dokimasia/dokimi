// Package handles declares interfaces whose methods return interfaces of
// the package, which the double of each producer doubles too.
package handles

import (
	"context"

	"example.com/shop/other"
)

// Archive opens logs through the interface Opener that it embeds. It sorts
// before every other producer of the package.
//
//dokimi:stub
type Archive interface {
	Opener
	// Size returns the number of logs, which is no handle.
	Size() int
}

// Opener opens logs.
type Opener interface {
	// OpenLog returns a log, a handle without a directive.
	OpenLog(name string) (Log, error)
}

// Log is a handle that Archive returns through Opener.
type Log interface {
	// Append appends a line to the log.
	Append(line string) error
}

// Producer opens cursors and subscriptions.
//
//dokimi:stub
type Producer interface {
	// Open returns a cursor, a handle without a directive.
	Open(ctx context.Context) (Cursor, error)
	// Resume returns a cursor too, so two methods return the handle.
	Resume(token string) Cursor
	// Watch returns a subscription, which has a directive of its own.
	Watch() Subscription
	// Nothing returns an interface without methods.
	Nothing() Empty
	// Elsewhere returns an interface of another package.
	Elsewhere() other.Thing
}

// Reader also returns a cursor, and the double of Producer already doubles
// it.
//
//dokimi:stub
type Reader interface {
	// Read returns a cursor.
	Read() Cursor
}

// Cursor is a handle that Open returns.
type Cursor interface {
	// Next returns the next value and whether it exists.
	Next(ctx context.Context) (string, bool, error)
	// Close releases the cursor.
	Close(ctx context.Context) error
}

// Subscription has a directive, so it gets a double of its own.
//
//dokimi:stub
type Subscription interface {
	// Cancel ends the subscription.
	Cancel()
}

// Empty declares no method.
type Empty interface{}
