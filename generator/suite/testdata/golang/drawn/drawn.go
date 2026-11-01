// Package drawn combines the tallies of baskets, under derived checks whose
// generated bindings draw the inputs of their calls.
package drawn

import "context"

// Tally counts the items and the value of a basket.
type Tally struct {
	Items int
	Value int
}

// Point is a point on a grid.
type Point struct {
	X int
	Y int
}

// Item is an item of a basket.
type Item struct {
	SKU   string
	Count int
}

// Tallies combines and records the tallies of baskets.
//
//dokimi:suite
type Tallies interface {
	// Merge returns the tally of two baskets together.
	//
	//dokimi:mixin commutative
	//dokimi:mixin associative
	Merge(a, b Tally) Tally

	// Distance returns the distance between two points.
	//
	//dokimi:mixin commutative
	Distance(a, b Point) int

	// Each calls fn with the SKU of each item.
	Each(ctx context.Context, fn func(sku string) error) error

	// Save stores the item that item points to.
	Save(ctx context.Context, item *Item) error

	// Limit limits the count of the items with the SKU sku to n.
	Limit(
		ctx context.Context,
		//dokimi:sample value=7
		n int,
		sku string,
	) error

	// Count returns the number of items.
	Count(ctx context.Context) (int, error)

	// Visit visits the item with the SKU sku, and calls each hook with it.
	Visit(ctx context.Context, sku string, hooks ...func(sku string) error) error

	// Schedule schedules the plan p for the items with the SKU sku.
	Schedule(ctx context.Context, sku string, p Plan) error

	// Link links the node n into the chain of the items.
	Link(ctx context.Context, n Node) error

	// Subscribe calls fn with the SKU of each item that changes.
	Subscribe(ctx context.Context, fn func(sku string)) error
}

// Plan is a plan of the items, with a hook that runs when the plan is done.
type Plan struct {
	Name   string
	OnDone func() error
}

// Node is a node of a chain of items.
type Node struct {
	SKU  string
	Next *Node
}

// Note is the note of an item.
type Note struct {
	SKU  string
	Text string
}

// FirstNote returns the note that a check writes first.
func FirstNote() Note { return Note{SKU: "sku-1", Text: "first"} }

// Notes keeps a note for each item.
//
//dokimi:suite
type Notes interface {
	// Put replaces the note under n.SKU in the shelf shelf.
	//
	//dokimi:contract updater role=writer id=noted
	Put(
		ctx context.Context,
		shelf string,
		//dokimi:sample value=FirstNote()
		n Note,
	) error

	// Get returns the note under sku.
	//
	//dokimi:contract updater role=reader id=noted
	Get(ctx context.Context, sku string) (Note, error)
}

// Hooked keeps a value in each scope, and calls a hook on each call.
//
//dokimi:suite
type Hooked interface {
	// Set stores value in scope and calls fn.
	//
	//dokimi:mixin partition read=Get
	Set(ctx context.Context, scope string, value int, fn func() int) error

	// Get returns the value in scope and calls fn.
	Get(ctx context.Context, scope string, fn func() int) (int, error)
}

// Found reports whether a lookup found its key.
type Found struct {
	Found bool
}

// Present reports whether an item is present.
type Present bool

// Enabled reports whether an item is enabled.
//
//dokimi:sample value=false alternate=true
type Enabled bool

// Switch is the state of the switch of an item.
type Switch struct {
	On Flag
}

// Flag is a flag of an item.
type Flag bool

// Amount is an amount of an item.
//
//dokimi:sample value=5 alternate=0
type Amount int

// Ratio is the ratio of two amounts.
//
//dokimi:sample value=0.5 alternate=0
type Ratio float64

// Code is the code of an item.
//
//dokimi:sample value=code alternate=""
type Code string

// Lookups finds the keys of the items.
//
//dokimi:suite
type Lookups interface {
	// Has reports whether the key key has an item.
	Has(ctx context.Context, key string) (Found, error)

	// IsPresent reports whether the key key has an item.
	IsPresent(ctx context.Context, key string) (Present, error)

	// IsEnabled reports whether the item with the key key is enabled.
	IsEnabled(ctx context.Context, key string) (Enabled, error)

	// State returns the state of the switch of the item with the key key.
	State(ctx context.Context, key string) (Switch, error)

	// Toggles returns the toggles of the item with the key key.
	Toggles(ctx context.Context, key string) ([]bool, error)

	// Amount returns the amount of the item with the key key.
	Amount(ctx context.Context, key string) (Amount, error)

	// Ratio returns the ratio of the item with the key key.
	Ratio(ctx context.Context, key string) (Ratio, error)

	// Code returns the code of the item with the key key.
	Code(ctx context.Context, key string) (Code, error)
}

// Label is the label of an item.
type Label struct {
	Name string
	Text string
}

// Labels keeps a label under each name.
//
//dokimi:suite
type Labels interface {
	// Put stores the label v under name.
	Put(ctx context.Context, name string, v Label) error

	// Get returns the label under the name that name returns.
	//
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, name func() string) (Label, error)
}
