// Package inputs moves the items of the shop through a pipeline.
package inputs

import (
	"context"
	"errors"
	"io"
)

// ErrDrained reports a feed without items.
var ErrDrained = errors.New("inputs: the feed is drained")

// ErrMissing reports an item that the feed does not have.
var ErrMissing = errors.New("inputs: the item is missing")

// Kind is the kind of an item.
type Kind int

// The kinds of an item.
const (
	Raw Kind = iota + 1
	Cooked
)

// Base is the part of an item that names its owner.
type Base struct {
	Owner string
}

// Item is one item of a pipeline.
//
//dokimi:builder
type Item struct {
	Base
	ID    string
	Count int `default:"1"`
	Kind  Kind
}

// audit is the part of a line that names who wrote it.
type audit struct {
	By string
}

// Line is one line of an order, with a declared default quantity and a
// declared default rank that no rule reads.
type Line struct {
	audit
	SKU  string
	Qty  int    `default:"1"`
	Rank int    `default:"high"`
	Memo string `prop:"-"`
}

// Note is the note of an item.
//
//dokimi:sample value=DefaultNote()
type Note struct {
	Text string
}

// DefaultNote returns the note of a new item.
func DefaultNote() Note { return Note{Text: "new"} }

// Handler handles one item.
type Handler func(item Item) error

// Option configures an item.
type Option func(item *Item)

// Tag is the tag of an item, with a declared default.
//
//dokimi:default "new"
type Tag string

// Code is the code of an item, with an authored sample.
//
//dokimi:sample value=DefaultCode()
type Code int

// DefaultCode returns the code of a new item.
func DefaultCode() Code { return 1 }

// Level is the level of an item, whose declared default no rule reads.
//
//dokimi:default unknown
type Level int

// Text is a transparent alias of a string.
type Text = string

// Bag is a list of values of one type.
type Bag[T any] []T

// Loop1 and Loop2 are defined types over each other, which no compiler
// accepts.
type (
	Loop1 Loop2
	Loop2 Loop1
)

// Feed feeds the items of a pipeline.
type Feed interface {
	// Next returns the next item, and false where the feed has none.
	//
	//dokimi:fault errors=[ErrDrained]
	Next(ctx context.Context) (Item, bool, error)
	// Get returns the item with the identity id.
	//
	//dokimi:mixin not-found sentinel=ErrMissing
	Get(ctx context.Context, id string) (Item, error)
	// Seek moves the feed to the first of the keys.
	Seek(keys ...string)
	// Close closes the feed.
	Close() error
}

// Source reads values of one type.
type Source[T any] interface {
	// Read returns the next value.
	Read() T
}

// Pipeline moves the items of a feed.
//
//dokimi:suite
type Pipeline interface {
	// Pump moves the items of the feed, and returns how many it moved.
	Pump(ctx context.Context, from Feed) (int, error)
	// Merge moves the items of two feeds.
	Merge(left, right Feed) error
	// Each calls fn for each item until fn returns false.
	Each(fn func(Item) bool)
	// Handle passes each item to h.
	Handle(h Handler) error
	// Do runs op once.
	Do(op func() (int, error))
	// Notify calls fn for each event.
	Notify(fn func(event string))
	// Apply calls fn with the counts of the items.
	Apply(fn func(counts ...int))
	// Fold folds the items of the channel in.
	Fold(ctx context.Context, items <-chan Item) error
	// Put stores the item with its line and its note.
	Put(item Item, line Line, note Note) error
	// Count counts the values of two sources.
	Count(ints Source[int], texts Source[string]) int
	// Spin takes a type that no compiler accepts.
	Spin(l Loop1)
	// Copy copies the items of the reader.
	Copy(src io.Reader) error
	// Configure configures the item with the name.
	Configure(name string, opts ...Option)
	// Subscribe moves the items of each feed.
	Subscribe(feeds ...Feed)
	// Label labels the items.
	Label(tag Tag, code Code, level Level, text Text, bag Bag[int])
}
