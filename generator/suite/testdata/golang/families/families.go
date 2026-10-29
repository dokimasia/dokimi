// Package families declares an interface for each family of the probes of the
// derived checks, under the classifications of the shape catalog.
package families

import (
	"context"
	"errors"
	"iter"
	"time"
)

// The sentinels that the classifications of the package name.
var (
	ErrMissing = errors.New("families: missing")
	ErrClosed  = errors.New("families: closed")
	ErrLimited = errors.New("families: limited")
	ErrSettled = errors.New("families: settled")
)

// Unsafe is the value that Render must not return.
var Unsafe = `<script>alert(1)</script>`

// Edge is the edge of the domain of Classify.
var Edge = ""

// Value is a value of the stores of the package.
type Value struct {
	Key  string
	Body string
	Rev  int64
}

// Invalid is a value that Validate refuses, and Unkeyed a value that Write
// refuses.
var (
	Invalid = Value{Body: "invalid"}
	Unkeyed = Value{Body: "unkeyed"}
)

// Level is a level, whose declared default is the untyped constant Base.
//
//dokimi:default Base
type Level int

// Base is the base level.
const Base = 2

// Meta is the revision of a value.
type Meta struct{ Revision int }

// Entry is one entry of a log.
type Entry struct {
	ID  string
	Seq int64
}

// Page is one page of a paged read.
type Page struct {
	Next  int
	Items []Value
}

// Cursor drains the values that Open opens.
type Cursor interface {
	Close(ctx context.Context) error
	Next(ctx context.Context) (Value, bool, error)
}

// Tx is a transaction.
type Tx struct{ ID int64 }

// Store reads and writes values.
//
//dokimi:suite
type Store interface {
	// Put stores v.
	Put(ctx context.Context, v Value) error

	// Get returns the value under key.
	//
	//dokimi:mixin read-after-write write=Put
	//dokimi:mixin monotonic-reads version=Rev
	Get(ctx context.Context, key string) (Value, error)

	// Load returns the value under key and whether it exists, the same on a
	// repeat.
	//
	//dokimi:mixin cacheable
	Load(ctx context.Context, key string) (Value, bool)

	// Split splits in into two halves.
	//
	//dokimi:mixin pure
	Split(in string) (string, string)

	// Find returns the values under keys.
	Find(ctx context.Context, keys ...string) ([]Value, error)

	// Level returns the level of key.
	Level(ctx context.Context, key string) (Level, error)

	// Raise raises the level of the store to l.
	Raise(ctx context.Context, l Level) error

	// Inspect returns the value under key with its revision, and whether it
	// exists.
	Inspect(key string) (Value, Meta, bool)

	// GetWithMeta returns the value under key with its revision.
	GetWithMeta(ctx context.Context, key string) (Value, Meta, error)
}

// Keyed stores values under keys.
//
//dokimi:suite
type Keyed interface {
	// Put stores v under key.
	Put(ctx context.Context, key string, v Value) error

	// Get returns the value under key.
	//
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, key string) (Value, error)

	// Delete removes the value under key.
	//
	//dokimi:mixin delete-removes read=Get sentinel=ErrMissing
	Delete(ctx context.Context, key string) error
}

// Expiring keeps each value for a minute.
//
//dokimi:suite
type Expiring interface {
	// Put stores v.
	Put(ctx context.Context, v Value) error

	// Read returns the value under key, and ErrMissing once its lifetime has
	// passed.
	//
	//dokimi:mixin ttl duration=1m put=Put read=Read not-found=ErrMissing
	Read(ctx context.Context, key string) (Value, error)
}

// Cached reads a value through a cache of its backing.
//
//dokimi:suite
type Cached interface {
	// Get returns the cached value under key.
	//
	//dokimi:contract cache role=cache
	Get(ctx context.Context, key string) (Value, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing
	Fetch(ctx context.Context, key string) (Value, error)

	// Store stores v in the backing.
	Store(ctx context.Context, v Value) error
}

// Tagged stores values under tags, which follow the value.
//
//dokimi:suite
type Tagged interface {
	// Put stores v under tag.
	Put(ctx context.Context, v Value, tag string) error

	// Get returns the value under tag.
	//
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, tag string) (Value, error)
}

// Snapshots reads a value across a write.
//
//dokimi:suite
type Snapshots interface {
	// Store stores v.
	Store(ctx context.Context, v Value) error

	// Get returns the value under key as the snapshot of the read sees it.
	//
	//dokimi:mixin point-in-time write=Store
	Get(ctx context.Context, key string) (Value, error)
}

// Bag drains its values.
//
//dokimi:suite
type Bag interface {
	// Add adds v.
	Add(ctx context.Context, v Value) error

	// Items returns the values.
	//
	//dokimi:mixin no-duplicates
	//dokimi:mixin stable-order
	Items(ctx context.Context) ([]Value, error)
}

// Updater replaces a value.
//
//dokimi:suite
type Updater interface {
	// Put replaces the value under v.Key.
	//
	//dokimi:contract updater role=writer id=valued
	Put(ctx context.Context, v Value) error

	// Get returns the value under key.
	//
	//dokimi:contract updater role=reader id=valued
	Get(ctx context.Context, key string) (Value, error)
}

// KeyedUpdater replaces the value under a key.
//
//dokimi:suite
type KeyedUpdater interface {
	// Put replaces the value under key.
	//
	//dokimi:contract updater role=writer id=keyed
	Put(ctx context.Context, key string, v Value) error

	// Get returns the value under key.
	//
	//dokimi:contract updater role=reader id=keyed
	Get(ctx context.Context, key string) (Value, error)
}

// Codec encodes a text.
//
//dokimi:suite
type Codec interface {
	// Encode encodes in.
	//
	//dokimi:contract codec role=forward fidelity=exact
	Encode(ctx context.Context, in string) (string, error)

	// Decode decodes in.
	//
	//dokimi:contract codec role=inverse
	Decode(ctx context.Context, in string) (string, error)
}

// Pager reads pages.
//
//dokimi:suite
type Pager interface {
	// Fetch returns the page at cursor.
	//
	//dokimi:contract pagination role=reader cursor=Next
	Fetch(ctx context.Context, cursor int) (*Page, error)
}

// Log appends entries.
//
//dokimi:suite
type Log interface {
	// Append appends e.
	//
	//dokimi:contract chain role=append
	//dokimi:mixin causal version=Seq
	Append(ctx context.Context, e Entry) error

	// Replay returns the entries.
	//
	//dokimi:contract chain role=replay
	Replay(ctx context.Context) ([]Entry, error)
}

// Opener opens cursors.
//
//dokimi:suite
type Opener interface {
	// Open opens a cursor.
	//
	//dokimi:contract cursor role=open id=cursors next-method=Next close-method=Close sentinel=ErrClosed
	Open(ctx context.Context) (Cursor, error)
}

// Ledger settles transactions.
//
//dokimi:suite
type Ledger interface {
	// Begin opens a transaction.
	//
	//dokimi:contract tx role=begin closed=ErrSettled
	Begin(ctx context.Context) (Tx, error)

	// Commit settles tx.
	//
	//dokimi:contract tx role=commit
	Commit(ctx context.Context, tx Tx) error

	// Rollback undoes tx.
	//
	//dokimi:contract tx role=rollback
	Rollback(ctx context.Context, tx Tx) error
}

// Guarded writes a value only where the stored value of its key matches it.
// It takes each value by its pointer, after the metadata of the write.
//
//dokimi:suite
type Guarded interface {
	// Put writes v with the metadata meta where Match accepts it.
	//
	//dokimi:contract if-match role=writer field=Body
	Put(ctx context.Context, meta Meta, v *Value) error

	// Match reports whether the stored value of the key of v has the body of
	// v.
	//
	//dokimi:contract if-match role=match
	Match(ctx context.Context, meta Meta, v *Value) (bool, error)
}

// Scope runs a body in a transaction.
//
//dokimi:suite
type Scope interface {
	// Run runs body in a transaction.
	//
	//dokimi:contract transaction role=fn not-found=ErrMissing write=Put read=Get
	Run(ctx context.Context, body func(ctx context.Context) error) error

	// Put stores v under key.
	Put(ctx context.Context, key string, v Value) error

	// Get returns the value under key.
	Get(ctx context.Context, key string) (Value, error)
}

// Events fires events.
//
//dokimi:suite
type Events interface {
	// Fire fires an event.
	//
	//dokimi:mixin hooks register=OnEvent
	Fire(ctx context.Context) error

	// OnEvent registers fn.
	OnEvent(fn func(event string))
}

// Timers runs work on the clock.
//
//dokimi:suite
type Timers interface {
	// At registers a task that runs once after has passed.
	//
	//dokimi:mixin scheduled schedule=At runs=Fired
	At(ctx context.Context, after time.Duration) error

	// Fired returns the number of tasks that ran with the delay after.
	Fired(ctx context.Context, after time.Duration) (int, error)

	// Slow returns within its budget.
	//
	//dokimi:mixin timeout duration=5s
	Slow(ctx context.Context, key string) error
}

// Election elects a leader.
//
//dokimi:suite
type Election interface {
	// Campaign takes the leadership.
	//
	//dokimi:contract leader-election role=campaign
	Campaign(ctx context.Context) error

	// Resign gives up the leadership.
	//
	//dokimi:contract leader-election role=resign
	Resign(ctx context.Context) error

	// IsLeader reports whether the subject leads.
	//
	//dokimi:contract leader-election role=is-leader
	IsLeader(ctx context.Context) bool
}

// Pool lends values.
//
//dokimi:suite
type Pool interface {
	// Get lends a value.
	//
	//dokimi:contract pool role=get
	Get(ctx context.Context) (Value, error)

	// Put takes v back.
	//
	//dokimi:contract pool role=put
	Put(ctx context.Context, v Value) error
}

// Batch writes values and refuses a value without a key.
//
//dokimi:suite
type Batch interface {
	// Write writes v.
	//
	//dokimi:contract batch-writer role=writer mode=atomic refused=Unkeyed
	Write(ctx context.Context, v Value) error

	// Read returns the value under key.
	//
	//dokimi:contract batch-writer role=reader
	Read(ctx context.Context, key string) (Value, error)
}

// Screened stores the values that Validate accepts.
//
//dokimi:suite
type Screened interface {
	// Store stores v.
	//
	//dokimi:mixin validates fn=Validate invalid=Invalid
	Store(ctx context.Context, v Value) error

	// Validate refuses an invalid v.
	Validate(v Value) error
}

// Limiter admits a burst.
//
//dokimi:suite
type Limiter interface {
	// Run admits one call.
	//
	//dokimi:contract rate-limit role=fn rate=100 burst=10 limited=ErrLimited
	Run(ctx context.Context, key string) error
}

// Renderer renders texts.
//
//dokimi:suite
type Renderer interface {
	// Render renders in.
	//
	//dokimi:mixin xss-safe unsafe=Unsafe
	Render(ctx context.Context, in string) (string, error)

	// Classify classifies in.
	//
	//dokimi:mixin total domain=strings edge=Edge
	Classify(ctx context.Context, in string) (string, error)
}

// Lifecycle stops.
//
//dokimi:suite
type Lifecycle interface {
	// Stop stops the subject.
	Stop()

	// Put stores the value under key, and accepts a second identical call.
	//
	//dokimi:mixin idempotent
	Put(ctx context.Context, key, value string) error

	// Sync succeeds within three attempts.
	//
	//dokimi:mixin retry-succeeds attempts=3
	Sync(ctx context.Context) error
}

// Shards keeps one partition from serving another.
//
//dokimi:suite
type Shards interface {
	// Put stores v in the partition of v.Key.
	//
	//dokimi:mixin partition read=Get
	Put(ctx context.Context, v Value) error

	// Get returns the value under key.
	Get(ctx context.Context, key string) (Value, error)
}

// Scoped keeps the values of one scope from serving another.
//
//dokimi:suite
type Scoped interface {
	// Set stores value under key in scope.
	//
	//dokimi:mixin partition read=Get
	Set(ctx context.Context, scope, key, value string) error

	// Get returns the value under key in scope.
	Get(ctx context.Context, scope, key string) (string, error)
}

// Plain stores values without an error.
//
//dokimi:suite
type Plain interface {
	// Set stores v.
	Set(v Value)

	// Get returns the value under key.
	//
	//dokimi:mixin read-after-write write=Set
	Get(ctx context.Context, key string) (Value, error)

	// Watch returns the values as they arrive.
	Watch(ctx context.Context) (<-chan Value, error)

	// All returns the values.
	All() iter.Seq[Value]
}

// Session is a handle that Begin opens.
type Session struct{ ID int64 }

// Next returns the next value of s.
func (s *Session) Next(ctx context.Context) (Value, bool, error) { return Value{}, false, ctx.Err() }

// Close closes s.
func (s *Session) Close(ctx context.Context) error { return ctx.Err() }

// Sessions opens sessions.
//
//dokimi:suite
type Sessions interface {
	// Open opens a session.
	//
	//dokimi:contract cursor role=open id=sessions next-method=Next close-method=Close sentinel=ErrClosed
	Open(ctx context.Context) (*Session, error)
}
