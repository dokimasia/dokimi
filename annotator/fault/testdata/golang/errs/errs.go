// Package errs declares the error values of the remote calls of the shop.
package errs

import "errors"

// ErrTimeout reports a remote call that did not answer in time.
var ErrTimeout = errors.New("errs: timeout")
