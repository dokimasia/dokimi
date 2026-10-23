package remote

import "example.com/shop/errs"

// Timeout reports whether err is the timeout of a remote call.
func Timeout(err error) bool {
	return err == errs.ErrTimeout
}
