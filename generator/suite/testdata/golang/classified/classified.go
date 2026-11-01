// Package classified declares the classifications whose claims one caller
// checks: a batch read, a stream, an indexed read, a wrapper and an accessor
// of a failure on a shelf, a batch read of the values that the writes store
// on a set of tags, and the verify of a chain, the release of a lease and the
// stats of a pool on a vault.
package classified

import (
	"context"
	"errors"
	"iter"
)

// ErrHeld reports an acquire of a held lease.
var ErrHeld = errors.New("classified: the lease is held")

// Item is a stored item.
type Item struct {
	Key  string
	Body string
}

// Usage is the accounting of the pool of a vault.
type Usage struct {
	Lent int
	Free int
}

// Stream is a stream of items.
type Stream = iter.Seq[Item]

// Shelf keeps items.
//
//dokimi:suite
type Shelf interface {
	// Put stores it.
	Put(ctx context.Context, it Item) error
	// GetMany returns the items of the keys.
	GetMany(ctx context.Context, keys ...string) ([]Item, error)
	// All returns each stored item, and the error that ends the stream.
	All(ctx context.Context) iter.Seq2[Item, error]
	// Items returns each stored item.
	Items(ctx context.Context) Stream
	// Scan returns each stored item, or the error that prevents the scan.
	//
	//dokimi:mixin over-match
	Scan(ctx context.Context) (iter.Seq[Item], error)
	// Keyed returns each stored item by its key.
	Keyed(ctx context.Context) iter.Seq2[string, Item]
	// Snapshot returns the list of the stored items.
	//
	//dokimi:mixin over-match
	Snapshot(ctx context.Context) ([]Item, error)
	// Len returns the number of items.
	Len(ctx context.Context) (int, error)
	// At returns the item at the position i.
	//
	//dokimi:mixin indexed by=Len
	At(ctx context.Context, i int) (Item, error)
	// AtOffset returns the item at the offset i.
	//
	//dokimi:mixin indexed by=Len
	AtOffset(ctx context.Context, i int64) (Item, error)
	// Slot returns the item at the position i of the section name.
	//
	//dokimi:mixin indexed by=Len
	Slot(ctx context.Context, name string, i uint) (Item, error)
	// Get returns the item of the key.
	Get(ctx context.Context, key string) (Item, error)
	// TracedGet returns what Get returns and traces the call.
	//
	//dokimi:mixin wrapped-via fn=Get
	TracedGet(ctx context.Context, key string) (Item, error)
	// Err reports the failure of the shelf.
	Err() error
}

// Tags keeps a set of tags.
//
//dokimi:suite
type Tags interface {
	// Add stores the tag.
	Add(ctx context.Context, tag string) error
	// Has returns the stored tags of tags.
	Has(ctx context.Context, tags ...string) ([]string, error)
}

// Vault is a chain of items, a lease of keys and a pool of items.
//
//dokimi:suite
type Vault interface {
	// Append appends it to the chain.
	//
	//dokimi:contract chain role=append
	Append(ctx context.Context, it Item) error
	// Replay returns the chain.
	//
	//dokimi:contract chain role=replay
	Replay(ctx context.Context) ([]Item, error)
	// Verify reports a change of the chain.
	//
	//dokimi:contract chain role=verify
	Verify(ctx context.Context) error
	// Acquire takes the lease of the key.
	//
	//dokimi:contract lease role=acquire taken=ErrHeld
	Acquire(ctx context.Context, key string) error
	// Release returns the lease of the key.
	//
	//dokimi:contract lease role=release
	Release(ctx context.Context, key string) error
	// Borrow lends an item of the pool.
	//
	//dokimi:contract pool role=get
	Borrow(ctx context.Context) (Item, error)
	// Return takes back an item.
	//
	//dokimi:contract pool role=put
	Return(ctx context.Context, it Item) error
	// Stats returns the accounting of the pool.
	//
	//dokimi:contract pool role=stats
	Stats(ctx context.Context) (Usage, error)
}
