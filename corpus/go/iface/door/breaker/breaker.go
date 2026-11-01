// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package breaker is the fixture of a circuit breaker over an operation whose
// calls fail on demand. The breaker opens after three failures of the
// operation in a row, and refuses each later call without the operation.
package breaker

import (
	"context"
	"errors"
)

// ErrOpen is the error of a call that an open breaker refuses.
var ErrOpen = errors.New("breaker: open")

// Breaker guards an operation.
//
//dokimi:suite machineries=[faults]
//dokimi:out path=breakertest/
type Breaker interface {
	// Do calls the operation, and returns ErrOpen without the operation once
	// the breaker is open.
	//
	//dokimi:contract circuit-breaker role=fn
	Do(ctx context.Context) error
}
