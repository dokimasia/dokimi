// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package txwithretrytest contains the double and the suite of TxWithRetry.
// InMemory is the control of the suite.
package txwithretrytest

import (
	"context"
	"errors"
	"sync"

	txwithretry "go.dokimi.dev/dokimi/corpus/go/iface/composite/tx-with-retry"
)

// The states of a transaction. A transaction starts Idle, Begin opens it, and
// it ends Committed or RolledBack.
const (
	Idle       = "idle"
	Open       = "open"
	Committed  = "committed"
	RolledBack = "rolled back"
)

// InMemory implements [txwithretry.TxWithRetry] in memory. The corpus applies
// the suite of TxWithRetry to it as a control.
//
//dokimi:subject txwithretry.TxWithRetry new=NewInMemory
type InMemory struct {
	state     string
	transient int
	mu        sync.Mutex
}

// The compiler checks that InMemory implements [txwithretry.TxWithRetry].
var _ txwithretry.TxWithRetry = (*InMemory)(nil)

// NewInMemory returns a transaction that has not begun, and whose commits
// succeed.
func NewInMemory() *InMemory { return NewInMemoryFailing(0) }

// NewInMemoryFailing returns a transaction that has not begun, and whose first
// transientCommits commits fail with [txwithretry.ErrTransient].
func NewInMemoryFailing(transientCommits int) *InMemory {
	return &InMemory{state: Idle, transient: transientCommits}
}

// Begin implements [txwithretry.TxWithRetry]. It refuses to open a second
// transaction.
func (s *InMemory) Begin(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state == Open {
		return txwithretry.ErrClosed
	}
	s.state = Open
	return nil
}

// Commit implements [txwithretry.TxWithRetry]. A transient failure leaves the
// transaction open, so a retry can commit it.
func (s *InMemory) Commit(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state != Open {
		return txwithretry.ErrClosed
	}
	if s.transient > 0 {
		s.transient--
		return txwithretry.ErrTransient
	}
	s.state = Committed
	return nil
}

// Rollback implements [txwithretry.TxWithRetry]. It refuses a transaction that
// is not open.
func (s *InMemory) Rollback(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state != Open {
		return txwithretry.ErrClosed
	}
	s.state = RolledBack
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("txwithretrytest: nil context")
	}
	return ctx.Err()
}
