// Package noprefix declares error values whose messages do not open with the
// name of the package.
//
//dokimi:sentinel prefix=off
package noprefix

import "errors"

// The error values of the retries.
var (
	// ErrRetry asks for another attempt.
	ErrRetry = errors.New("retry the operation")
	// ErrGiveUp reports that no attempt remains.
	ErrGiveUp = errors.New("no attempts remain")
)
