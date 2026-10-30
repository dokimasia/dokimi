// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ratelimittest contains the double and the suite of Contract. InMemory
// is the subject of the suite.
package ratelimittest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	ratelimit "go.dokimi.dev/dokimi/corpus/go/iface/contract/rate-limit"
)

// ErrLimited is the error of Run once the bucket of the limiter is empty.
var ErrLimited = ratelimit.ErrLimited

// The burst and the period of the limiter of NewInMemory.
const (
	burst  = 10
	period = time.Second / 100
)

// InMemory implements [ratelimit.Contract] in memory. The corpus applies the
// suite of Contract to it as a subject. It is a token bucket that gains one
// call per period up to the burst.
//
//dokimi:subject ratelimit.Contract new=NewInMemory
type InMemory struct {
	filled time.Time
	clock  assert.Clock
	burst  int
	period time.Duration
	tokens int
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [ratelimit.Contract].
var _ ratelimit.Contract = (*InMemory)(nil)

// NewInMemory returns a limiter on a controlled clock at the Unix epoch, with a
// burst of 10 calls and one more call every 10 ms.
func NewInMemory() *InMemory {
	return NewInMemoryOn(assert.NewControlled(time.Unix(0, 0).UTC()), burst, period)
}

// NewInMemoryOn returns a limiter on clk with a full bucket of burst calls.
func NewInMemoryOn(clk assert.Clock, burst int, period time.Duration) *InMemory {
	return &InMemory{
		clock:  clk,
		burst:  burst,
		period: period,
		tokens: burst,
		filled: clk.Now(),
	}
}

// Run implements [ratelimit.Contract].
func (s *InMemory) Run(ctx context.Context, _ string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refill()
	if s.tokens == 0 {
		return ErrLimited
	}
	s.tokens--
	return nil
}

// refill adds one call for each whole period since the last fill, up to the
// burst.
func (s *InMemory) refill() {
	elapsed := s.clock.Now().Sub(s.filled)
	if elapsed < s.period {
		return
	}
	periods := int(elapsed / s.period)
	s.filled = s.filled.Add(time.Duration(periods) * s.period)
	s.tokens = min(s.tokens+periods, s.burst)
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("ratelimittest: nil context")
	}
	return ctx.Err()
}
