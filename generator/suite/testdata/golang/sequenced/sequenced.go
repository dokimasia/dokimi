// Package sequenced declares the classifications whose claims generated
// sequences check: the bounds of a counter and of a list of notes, the merges
// of a set, a tally and a pile, the graphs of a run, the positions of a log,
// the streams of a feed, the references of a register and of a function, the
// reads of an index, the reads of versioned items, and the functions of the
// param normalize.
package sequenced

import (
	"context"
	"iter"
)

// Item is a stored item.
type Item struct {
	Key string
	Rev int
}

// Option configures one call of Counter.Inc.
type Option func(*int)

// Graph is the graph of the states of a run. Its initializer states its type.
var Graph = map[string][]string{"": {"running"}, "running": {"done"}}

// Retries is the graph of the retries of a run, with its type stated.
var Retries map[string][]string = map[string][]string{"running": {""}}

// Order is the list of the states of a run, which is no graph.
var Order []string = []string{"", "running", "done"}

// Counter counts calls.
//
//dokimi:suite normalize=Counted
type Counter interface {
	// Inc adds one to the count, which is at most 3. Each option configures
	// the call.
	//
	//dokimi:mixin bounded limit=3 observe=Value
	Inc(ctx context.Context, opts ...Option) error
	// Value returns the count.
	//
	//dokimi:mixin monotonic
	Value(ctx context.Context) (int, error)
	// Grow adds a byte to the buffer, which has at most 64 bytes.
	//
	//dokimi:mixin bounded limit=64 observe=Size
	Grow(ctx context.Context) error
	// Size returns the count of the bytes of the buffer.
	Size(ctx context.Context) (uint64, error)
	// Throttle adds to the load, which is at most 1.
	//
	//dokimi:mixin bounded limit=1 observe=Load
	Throttle(ctx context.Context) error
	// Load returns the load.
	Load(ctx context.Context) (float64, error)
}

// Counted returns the values of a call without the name of its callable,
// which no function of the param normalize does.
func Counted(values []any) []any { return values }

// Notes keeps a list of items, a text and a pinned item.
//
//dokimi:suite normalize=Noted
type Notes interface {
	// Push adds an item to the list, which has at most 4 items.
	//
	//dokimi:mixin bounded limit=4 observe=Items
	Push(ctx context.Context, item string) error
	// Pop removes the last item of the list.
	//
	//dokimi:mixin bounded limit=1 observe=Push
	Pop(ctx context.Context) error
	// Items returns the items of the list that start with prefix.
	Items(ctx context.Context, prefix string) ([]string, error)
	// Append appends to the text, which has at most 8 bytes.
	//
	//dokimi:mixin bounded limit=8 observe=Text
	Append(ctx context.Context, text string) error
	// Text returns the text.
	Text(ctx context.Context) (string, error)
	// Pin pins it.
	//
	//dokimi:mixin bounded limit=1 observe=Head
	Pin(ctx context.Context, it Item) error
	// Head returns the pinned item.
	Head(ctx context.Context) (Item, error)
}

// Noted returns no values, which no function of the param normalize does.
func Noted(operation string, values []any) {}

// Replica is a peer that a merge reads.
type Replica interface {
	// Items returns the items of the peer.
	Items(ctx context.Context) ([]string, error)
}

// Sink is a peer that a merge drains into.
type Sink interface {
	// Drain drains the sink.
	Drain(ctx context.Context) error
}

// Set is a set of items that merges the items of other replicas.
//
//dokimi:suite
type Set interface {
	// Merge adds the items of peer.
	//
	//dokimi:mixin crdt-merge write=Add read=Items
	Merge(ctx context.Context, peer Replica) error
	// Join adds the items of other.
	//
	//dokimi:mixin crdt-merge write=Add read=Items
	Join(ctx context.Context, other Set) error
	// Weigh adds the items of peer weight times.
	//
	//dokimi:mixin crdt-merge write=Add read=Items
	Weigh(ctx context.Context, peer Replica, weight int) error
	// Add adds an item.
	Add(ctx context.Context, item string) error
	// Items returns the items.
	Items(ctx context.Context) ([]string, error)
}

// Tally is a count that merges the count of another tally.
//
//dokimi:suite
type Tally struct {
	n int
}

// Merge adds the count of other.
//
//dokimi:mixin crdt-merge write=Add read=Count
func (t *Tally) Merge(ctx context.Context, other *Tally) error {
	t.n += other.n
	return ctx.Err()
}

// Add adds one to the count.
func (t *Tally) Add(ctx context.Context) error {
	t.n++
	return ctx.Err()
}

// Count returns the count.
func (t *Tally) Count(ctx context.Context) (int, error) { return t.n, ctx.Err() }

// Pile is a count that merges a number.
//
//dokimi:suite
type Pile interface {
	// Absorb adds n to the count.
	//
	//dokimi:mixin crdt-merge write=Add read=Count
	Absorb(ctx context.Context, n int) error
	// Pour drains the count into into.
	//
	//dokimi:mixin crdt-merge write=Add read=Count
	Pour(ctx context.Context, into Sink) error
	// Add adds one to the count.
	Add(ctx context.Context) error
	// Count returns the count.
	Count(ctx context.Context) (int, error)
}

// Run moves through the states of a run.
//
//dokimi:suite
type Run interface {
	// Advance moves the run of key to its next state.
	//
	//dokimi:contract workflow role=fn id=forward transitions=Graph observe=Status
	Advance(ctx context.Context, key string) error
	// Retry moves the run of key back to the state before.
	//
	//dokimi:contract workflow role=fn id=back transitions=Retries observe=Status
	Retry(ctx context.Context, key string) error
	// Restart moves the run of key to its first state.
	//
	//dokimi:contract workflow role=fn id=again transitions=Order observe=Status
	Restart(ctx context.Context, key string) error
	// Toggle moves the run of key to its next state, and reports the move
	// where verbose returns true.
	//
	//dokimi:contract workflow role=fn id=toggle transitions=Graph observe=Status
	Toggle(ctx context.Context, key string, verbose func() bool) error
	// Status returns the state of the run of key, and the history of the run
	// where verbose is true.
	Status(ctx context.Context, key string, verbose bool) (string, error)
}

// Log appends entries at rising positions.
//
//dokimi:suite normalize=Logged
type Log interface {
	// Append appends an entry and returns its position.
	//
	//dokimi:contract appender role=fn
	Append(ctx context.Context, entry string) (int64, error)
}

// Logged returns the values of a call beside an error, which no function of
// the param normalize does.
func Logged(operation string, values []any) ([]any, error) { return values, nil }

// Feed yields its items.
//
//dokimi:suite
type Feed interface {
	// Put adds an item.
	Put(ctx context.Context, item string) error
	// All yields each item, and the error that ends the stream.
	//
	//dokimi:mixin no-duplicates
	All(ctx context.Context) iter.Seq2[string, error]
	// Each yields each item.
	//
	//dokimi:mixin no-duplicates
	Each(ctx context.Context) iter.Seq[string]
	// Pairs yields each item with its count.
	//
	//dokimi:mixin no-duplicates
	Pairs(ctx context.Context) iter.Seq2[string, int]
}

// Register keeps one value, and a reference implementation checks it.
//
//dokimi:suite machineries=[reference] normalize=Unstamped
type Register interface {
	// Set stores v.
	Set(ctx context.Context, v string) error
	// Get returns the stored value.
	Get(ctx context.Context) (string, error)
}

// Unstamped returns the values of a call of a Register as they are.
func Unstamped(operation string, values []any) []any { return values }

// Double returns twice n, and a reference implementation checks it.
//
//dokimi:suite machineries=[reference]
func Double(ctx context.Context, n int) (int, error) { return 2 * n, ctx.Err() }

// Cached reads items through two caches. The backing of each cache takes
// another input than the cache.
//
//dokimi:suite
type Cached interface {
	// Get returns the item of the key from the first cache.
	//
	//dokimi:contract cache role=cache id=first
	Get(ctx context.Context, key string) (Item, error)
	// Load returns the item of the identity from the backing of the first
	// cache.
	//
	//dokimi:contract cache role=backing id=first
	Load(ctx context.Context, id string) (Item, error)
	// Peek returns the item of the key from the second cache.
	//
	//dokimi:contract cache role=cache id=second
	Peek(ctx context.Context, key string) (Item, error)
	// Fetch returns the item at the position key from the backing of the
	// second cache.
	//
	//dokimi:contract cache role=backing id=second
	Fetch(ctx context.Context, key int) (Item, error)
}

// Index finds items.
//
//dokimi:suite
type Index interface {
	// Put stores it.
	Put(ctx context.Context, it Item) error
	// Get returns the stored item of the key.
	//
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, key string) (Item, error)
	// Find returns the stored item of the key that keep accepts.
	//
	//dokimi:mixin read-after-write write=Put
	Find(ctx context.Context, key string, keep func(Item) bool) (Item, error)
	// ByRev returns the stored item of the revision.
	//
	//dokimi:mixin read-after-write write=Put
	ByRev(ctx context.Context, rev int) (Item, error)
	// Lookup returns the stored item of the key and the revision.
	//
	//dokimi:mixin read-after-write write=Put
	Lookup(ctx context.Context, key string, rev int) (Item, error)
}

// Stamped keeps items that the store stamps with a revision. The claim of the
// reads of sequences of Get has no partner, and the claim of its single
// caller pairs it with Store.
//
//dokimi:suite
type Stamped interface {
	// Store stores it.
	Store(ctx context.Context, it Item) error
	// Get returns the item of the key.
	//
	//dokimi:mixin monotonic-reads version=Rev
	Get(ctx context.Context, key string) (Item, error)
}

// Versioned keeps versioned items.
//
//dokimi:suite
type Versioned interface {
	// Put stores it.
	Put(ctx context.Context, it Item) error
	// Get returns the item of the key.
	//
	//dokimi:mixin writes-follow-reads version=Rev
	Get(ctx context.Context, key string) (Item, error)
}
