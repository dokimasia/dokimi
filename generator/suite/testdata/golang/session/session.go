// Package session keeps the sessions of signed-in users.
package session

import "context"

// Session is the session of one signed-in user.
type Session struct {
	ID   string
	User string
}

// Sessions keeps sessions by their ID.
//
//dokimi:suite
type Sessions interface {
	// Put stores the session s.
	Put(ctx context.Context, s Session) error
	// Get returns the session with the ID id.
	Get(ctx context.Context, id string) (Session, error)
}
