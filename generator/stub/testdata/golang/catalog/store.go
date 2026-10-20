// Package catalog declares the declarations under the directive stub that
// the specs read.
package catalog

import (
	"context"
	"time"
)

// Item is what a store keeps.
type Item struct {
	// ID identifies the item.
	ID string
}

// Store has a method of each form of a callable, and the method that it
// embeds from Closer.
//
//dokimi:stub
type Store interface {
	// Get can fail, and takes a context.
	Get(ctx context.Context, id string) (Item, error)
	// Put has parameters without names, and an error as its one result.
	Put(context.Context, Item) error
	// Len cannot fail.
	Len() int
	// Lookup has named results and a presence flag.
	Lookup(key string) (item Item, ok bool)
	// Log is variadic and has no result. Its variadic parameter has elements
	// with a sample.
	Log(format string, args ...string)
	// Split has two results without names beside its error.
	Split(s string) (string, string, error)
	Closer
}

// Closer releases the resources of a store.
type Closer interface {
	// Close can fail.
	Close() error
}

// Clock is a struct whose double implements its methods.
//
//dokimi:stub
type Clock struct {
	// zone is a field, which a double of the struct leaves out.
	zone string
}

// Now returns the instant of the clock.
func (c *Clock) Now() time.Time { return time.Time{} }

// Zone returns the zone of the clock.
func (c Clock) Zone() string { return c.zone }

// Every calls fn at each interval d until a call of stop. Its parameter fn
// and its result stop have no sample.
func (c *Clock) Every(d time.Duration, fn func()) (stop func()) { return func() {} }
