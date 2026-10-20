// Package client declares an interface under the directive that returns a
// handle of the package handles.
package client

import "example.com/shop/handles"

// Client returns a cursor of the package handles, whose double a producer of
// that package writes.
//
//dokimi:stub
type Client interface {
	// Cursor returns a cursor of the package handles.
	Cursor() handles.Cursor
}
