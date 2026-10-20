// Package mixed declares an interface whose methods have the mixins
// deprecated and order-after.
package mixed

import (
	"context"
	"errors"
)

// ErrNotReady is what Commit reports before Prepare.
var ErrNotReady = errors.New("mixed: not ready")

// Session commits after it prepares, and keeps an old method.
//
//dokimi:stub
type Session interface {
	// Prepare runs before Commit.
	Prepare(ctx context.Context) error
	// Commit is valid after Prepare alone.
	//
	//dokimi:mixin order-after fn=Prepare unready=ErrNotReady
	Commit(ctx context.Context) error
	// Old is scheduled for removal.
	//
	//dokimi:mixin deprecated
	Old(key string) (string, error)
}
