// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package singleflighttest

import (
	"context"
	"fmt"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/singleflight"
)

// Ledgered implements [singleflight.Contract] over a ledger of its flights.
// The first call of a key records a flight in the ledger before it computes,
// and Flights reads the count of the ledger. The corpus applies the suite of
// Contract to it as a subject, with a ledger whose calls fail on demand.
//
//dokimi:subject singleflight.Contract new=NewLedgered options=LedgeredOptions
type Ledgered struct {
	ledger *ledger
	calls  map[string]*call
	mu     sync.Mutex
}

// The compiler checks that Ledgered implements [singleflight.Contract].
var _ singleflight.Contract = (*Ledgered)(nil)

// NewLedgered returns a Ledgered over a ledger in memory.
func NewLedgered() *Ledgered {
	return &Ledgered{ledger: &ledger{}, calls: map[string]*call{}}
}

// LedgeredOptions returns the door of a Ledgered over a ledger whose calls
// fail on demand.
func LedgeredOptions() []ContractOption {
	return []ContractOption{conform.Faultable(faultyLedgered)}
}

// faultyLedgered returns a Ledgered over a ledger that calls a dependency of
// inj before each record and each count.
func faultyLedgered(_ *ContractHandle, inj *conform.Injector) singleflight.Contract {
	return &Ledgered{ledger: &ledger{dependency: inj.Dependency()}, calls: map[string]*call{}}
}

// Run implements [singleflight.Contract]. A flight whose record fails ends
// with the error of the ledger, which each call of the flight returns. The
// next call of the key starts a new flight.
func (s *Ledgered) Run(ctx context.Context, key string, compute func() string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}

	s.mu.Lock()
	if c, seen := s.calls[key]; seen {
		s.mu.Unlock()
		<-c.done
		return c.val, c.err
	}
	c := &call{done: make(chan struct{})}
	s.calls[key] = c
	s.mu.Unlock()

	c.err = s.ledger.record()
	if c.err != nil {
		s.mu.Lock()
		delete(s.calls, key)
		s.mu.Unlock()
	} else if compute != nil {
		c.val = compute()
	}
	close(c.done)
	return c.val, c.err
}

// Flights implements [singleflight.Contract].
func (s *Ledgered) Flights(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	return s.ledger.count()
}

// ledger counts the flights of a Ledgered. A ledger with a dependency calls
// it before each record and each count, so its calls fail on demand.
type ledger struct {
	dependency *conform.Dependency
	flights    int
	mu         sync.Mutex
}

// record records one flight.
func (l *ledger) record() error {
	if err := l.ask("Record"); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flights++
	return nil
}

// count returns the number of the recorded flights.
func (l *ledger) count() (int, error) {
	if err := l.ask("Count"); err != nil {
		return 0, err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flights, nil
}

// ask calls the callable of the dependency of the ledger. A ledger without a
// dependency asks nothing and returns nil.
func (l *ledger) ask(callable string) error {
	if l.dependency == nil {
		return nil
	}
	if err := l.dependency.Call(callable); err != nil {
		return fmt.Errorf("singleflighttest: the ledger refuses %s: %w", callable, err)
	}
	return nil
}
