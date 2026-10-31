// Package hidden declares a double whose results have samples that only this
// package can refer to.
package hidden

import "context"

// Host is a host whose samples are the results of unexported functions.
//
//dokimi:sample value=newHost() alternate=otherHost()
type Host struct {
	// Name is the name of the host.
	Name string
}

// newHost returns the sample of Host.
func newHost() Host { return Host{Name: "eu"} }

// otherHost returns the alternate sample of Host.
func otherHost() Host { return Host{Name: "us"} }

// Hosts finds hosts.
//
//dokimi:stub
type Hosts interface {
	// Find returns the host of the name.
	Find(ctx context.Context, name string) (Host, error)
}

// store keeps values by key. Only this package can refer to it.
//
//dokimi:stub
type store interface {
	// Get returns the value under the key.
	Get(ctx context.Context, key string) (string, error)
	// Open opens a session of the store.
	Open(ctx context.Context) (session, error)
}

// session is a session of a store. Only this package can refer to it.
type session interface {
	// Close ends the session.
	Close() error
}

// Value returns the value of the key a in s, so that store has a use.
func Value(ctx context.Context, s store) (string, error) { return s.Get(ctx, "a") }

// node is a node that only this package can refer to.
type node struct {
	// ID is the identity of the node.
	ID int
}

// Finder finds values. The witness of T is a type that only this package can
// refer to.
//
//dokimi:stub
//dokimi:witness T=node
type Finder[T any] interface {
	// Find returns the value of the name.
	Find(ctx context.Context, name string) (T, error)
}
