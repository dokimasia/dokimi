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
	// Append stores the value under the key, and returns its position.
	//
	//dokimi:mixin durable
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
