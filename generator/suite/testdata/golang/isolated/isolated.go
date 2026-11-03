// Package isolated keeps the lists and the registers of transactional stores,
// whose transactions are serializable or have snapshot isolation. It also
// settles transactions without a handle of the package.
package isolated

import (
	"context"
	"database/sql/driver"
	"errors"
)

// ErrClosed is the error of a settlement of a closed transaction.
var ErrClosed = errors.New("isolated: closed")

// Tx is the handle of a transaction of Lists.
type Tx struct{ ID int64 }

// Lists keeps lists of numbers, which transactions append to. Each callable
// of a transaction takes its handle.
//
//dokimi:suite
type Lists interface {
	// Begin opens a transaction under the label.
	//
	//dokimi:contract tx role=begin closed=ErrClosed id=lists
	Begin(ctx context.Context, label string) (Tx, error)

	// Commit applies the transaction.
	//
	//dokimi:contract tx role=commit id=lists
	Commit(ctx context.Context, tx Tx) error

	// Rollback discards the transaction.
	//
	//dokimi:contract tx role=rollback id=lists
	Rollback(ctx context.Context, tx Tx) error

	// Append appends v to the list of key in the transaction.
	//
	//dokimi:mixin serializable read=Get
	Append(ctx context.Context, tx Tx, key string, v int) error

	// Get returns the list of key in the transaction.
	Get(ctx context.Context, tx Tx, key string) ([]int, error)
}

// Registers keeps one text under each key, which transactions write through
// the methods of their handle.
//
//dokimi:suite
type Registers interface {
	// Begin opens a transaction and returns its handle.
	//
	//dokimi:contract tx role=begin closed=ErrClosed id=registers
	Begin(ctx context.Context) (Txn, error)

	// Latest returns the handle of the latest transaction.
	Latest(ctx context.Context) (Txn, error)
}

// Txn is the handle of a transaction of Registers.
type Txn interface {
	// Commit applies the transaction.
	//
	//dokimi:contract tx role=commit id=registers
	Commit(ctx context.Context) error

	// Rollback discards the transaction.
	//
	//dokimi:contract tx role=rollback id=registers
	Rollback(ctx context.Context) error

	// Set writes v under key in the transaction.
	//
	//dokimi:mixin snapshot-isolation read=Load
	Set(ctx context.Context, key, v string) error

	// Load returns the text under key in the transaction.
	Load(ctx context.Context, key string) (string, error)
}

// Latch settles transactions without a handle.
//
//dokimi:suite
type Latch interface {
	// Begin opens a transaction, and waits until it can.
	//
	//dokimi:contract tx role=begin closed=ErrClosed id=latch
	Begin(ctx context.Context)

	// Commit applies the transaction.
	//
	//dokimi:contract tx role=commit id=latch
	Commit(ctx context.Context) error

	// Rollback discards the transaction.
	//
	//dokimi:contract tx role=rollback id=latch
	Rollback(ctx context.Context) error
}

// Driver settles transactions through the handle of the package driver.
//
//dokimi:suite
type Driver interface {
	// Begin opens a transaction and returns its handle.
	//
	//dokimi:contract tx role=begin closed=ErrClosed id=driver
	Begin(ctx context.Context) (driver.Tx, error)

	// Commit applies the transaction.
	//
	//dokimi:contract tx role=commit id=driver
	Commit(ctx context.Context, tx driver.Tx) error

	// Rollback discards the transaction.
	//
	//dokimi:contract tx role=rollback id=driver
	Rollback(ctx context.Context, tx driver.Tx) error
}
