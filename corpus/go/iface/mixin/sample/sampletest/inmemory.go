// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sampletest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package sampletest

import (
	"context"
	"errors"
	"strings"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/sample"
)

// ErrUnprocessable is the error of Process for an input that NewInput did not
// return.
var ErrUnprocessable = errors.New("sampletest: input is not sampled")

// inputPrefix is the prefix of each input that NewInput returns.
const inputPrefix = "sample:"

// InMemory implements [sample.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a subject.
//
//dokimi:subject sample.Mixed new=NewInMemory
type InMemory struct {
	mu   sync.Mutex
	next int
}

// The compiler checks that InMemory implements [sample.Mixed].
var _ sample.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// NewInput implements [sample.Mixed].
func (s *InMemory) NewInput(ctx context.Context) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	return inputPrefix + string(rune('a'+(s.next-1)%26)), nil
}

// Process implements [sample.Mixed].
func (*InMemory) Process(ctx context.Context, input string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	if !strings.HasPrefix(input, inputPrefix) {
		return "", ErrUnprocessable
	}
	return strings.ToUpper(strings.TrimPrefix(input, inputPrefix)), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("sampletest: nil context")
	}
	return ctx.Err()
}
