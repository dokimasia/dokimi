// Package refused keeps the clock and the tally of a node, whose callables
// refuse the claims that their declarations do not make.
package refused

import "context"

// HybridClock is the hybrid clock of a node.
//
//dokimi:suite
type HybridClock interface {
	// Update merges the instant of a peer into the clock. It returns an error
	// for an instant beyond the max offset of the clock.
	//
	//dokimi:refuse zero-on-error reason="Update fails only for an instant beyond the max offset of the clock"
	//dokimi:refuse cancel reason="Update returns before it reads its context"
	Update(ctx context.Context, observed int64) (int64, error)

	// MaxOffset returns the max offset of the clock.
	MaxOffset(ctx context.Context) (int64, error)
}

// Tally records amounts, and refuses a claim that its write does not have.
//
//dokimi:suite
type Tally interface {
	// Record records the amount.
	//
	//dokimi:refuse deterministic reason="Record writes the amount"
	Record(ctx context.Context, amount int64) error
}

// Doubled refuses one claim twice.
//
//dokimi:suite
type Doubled interface {
	// Read returns the value of the key.
	//
	//dokimi:refuse zero-on-error reason="Read never fails"
	//dokimi:refuse zero-on-error reason="Read never fails at all"
	Read(ctx context.Context, key string) (string, error)
}

// Unsuited has no directive suite.
type Unsuited interface {
	// Get returns the count.
	//
	//dokimi:refuse smoke reason="Get has no suite"
	Get() int
}

// Title returns the title of a name.
//
//dokimi:suite
//dokimi:refuse deterministic reason="Title salts each name"
func Title(name string) string { return name }

// Untitled returns a name without a directive suite.
//
//dokimi:refuse deterministic reason="Untitled has no suite"
func Untitled(name string) string { return name }

// Watch reads a clock, and refuses a claim that its read does not have.
//
//dokimi:suite
type Watch interface {
	// Now returns the current instant of the watch.
	//
	//dokimi:mixin time-aware
	//dokimi:refuse zero-on-error reason="Now never fails"
	Now() int64
}

// Spinner spins, and refuses a claim that its spin does not have.
//
//dokimi:suite
type Spinner interface {
	// Spin turns the spinner by the turns.
	//
	//dokimi:refuse cancel reason="Spin takes no context"
	Spin(turns int)
}

// Opaque is a value without an exported field.
type Opaque struct{ v int }

// Marker marks opaque values, and refuses a claim of a callable without
// claims.
//
//dokimi:suite
type Marker interface {
	// Mark marks the value.
	//
	//dokimi:refuse deterministic reason="Mark returns nothing"
	Mark(o Opaque)
}

// Pair reads and writes numbers, and refuses a claim that its write does not
// have.
//
//dokimi:suite
type Pair interface {
	// Read returns the number under the key.
	//
	//dokimi:refuse zero-on-error reason="Read returns the last number on an error"
	Read(ctx context.Context, key string) (int64, error)

	// Write stores the number under the key.
	//
	//dokimi:refuse cancel reason="Write returns before it reads its context"
	//dokimi:refuse zero-on-error reason="Write returns no number"
	Write(ctx context.Context, key string, n int64) error

	// Len returns the count of the keys.
	Len() int
}

// Shout returns the text in capitals, and refuses a claim that it does not
// have.
//
//dokimi:suite
//dokimi:refuse cancel reason="Shout takes no context"
func Shout(text string) string { return text }
