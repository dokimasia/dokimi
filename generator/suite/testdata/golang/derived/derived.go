// Package derived keeps the sessions of signed-in users, under the
// classifications of the shape catalog.
package derived

import (
	"context"
	"errors"
)

// ErrNotFound is the error of a read of a session that the store does not
// have.
var ErrNotFound = errors.New("derived: not found")

// Session is the session of one signed-in user, with a field that code
// outside the package cannot set.
type Session struct {
	ID    string
	User  string
	trace string
}

// Count is a count of visits, with a declared default.
//
//dokimi:default 3
type Count int

// Name is the name of a visitor, without a declared default.
type Name string

// Visit is one visit, with a declared default of its field Pages.
type Visit struct {
	Pages int `default:"5"`
	Note  string
}

// Store keeps the sessions of signed-in users.
//
//dokimi:suite
type Store interface {
	// Put stores s under s.ID.
	Put(ctx context.Context, s Session) error

	// Get returns the session with the ID id.
	//
	//dokimi:mixin not-found sentinel=ErrNotFound
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, id string) (Session, error)

	// Find returns the session with the ID id, and nil for none.
	Find(ctx context.Context, id string) (*Session, error)

	// Count counts the visits of the visitor.
	Count(
		//dokimi:sample value=7
		stated Count,
		//dokimi:sample value=seven
		refused Count,
		plain Count,
		name Name,
		visit Visit,
	) int
}
