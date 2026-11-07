// Package vault keeps values on a medium, apart from the process of the vault
// that writes them.
package vault

import (
	"context"
	"errors"
)

// ErrCorrupt is the error of a read of a value whose record a flipped bit
// corrupted.
var ErrCorrupt = errors.New("vault: the record is corrupt")

// Vault keeps values by key on a medium.
//
//dokimi:suite machineries=[medium]
type Vault interface {
	// Put stores the value under the key.
	//
	//dokimi:mixin durable corrupt=ErrCorrupt
	Put(ctx context.Context, key, value string) error

	// Get returns the value under the key.
	Get(ctx context.Context, key string) (string, error)
}

// Journal appends values by key on a medium, and returns the position of each
// append.
//
//dokimi:suite machineries=[medium]
type Journal interface {
	// Append stores the value under the key, and returns its position. Callers
	// may append at once, and read one at a time.
	//
	//dokimi:mixin durable
	//dokimi:mixin concurrent
	Append(ctx context.Context, key, value string) (int64, error)

	// Get returns the value under the key.
	Get(ctx context.Context, key string) (string, error)
}

// Ledger keeps values by key on a medium, and reports a key without a value
// with a presence flag.
//
//dokimi:suite machineries=[medium]
type Ledger interface {
	// Record stores the value under the key.
	//
	//dokimi:mixin durable
	Record(key, value string) error

	// Lookup returns the value under the key, and whether the ledger has one.
	Lookup(key string) (string, bool)
}

// Registry keeps values by key on a medium for concurrent callers, and agrees
// with its reference through crashes and faults.
//
//dokimi:suite machineries=[medium, reference, faults]
type Registry interface {
	// Set stores the value under the key.
	//
	//dokimi:mixin durable corrupt=ErrCorrupt
	//dokimi:mixin concurrent
	Set(ctx context.Context, key, value string) error

	// Value returns the value under the key.
	//
	//dokimi:mixin concurrent
	Value(ctx context.Context, key string) (string, error)

	// Await returns the value under the key once a call of Set stores one.
	//
	//dokimi:mixin waits
	Await(ctx context.Context, key string) (string, error)

	// Add stores the value under a new key, and returns the key.
	//
	//dokimi:mixin creates
	//dokimi:mixin concurrent
	Add(ctx context.Context, value string) (string, error)

	// Insert stores the value under a new key, and returns the entry with the
	// key.
	//
	//dokimi:mixin creates key=Key
	//dokimi:mixin concurrent
	Insert(ctx context.Context, value string) (*Entry, error)

	// Stamp returns a new number, which no callable takes.
	//
	//dokimi:mixin creates
	//dokimi:mixin concurrent
	Stamp(ctx context.Context) (int64, error)

	// Touch records a use of the registry, and returns nothing.
	//
	//dokimi:mixin creates
	//dokimi:mixin concurrent
	Touch(ctx context.Context)
}

// Entry is a value of Registry with the key that Insert gave it.
type Entry struct {
	Key   string
	Value string
}
