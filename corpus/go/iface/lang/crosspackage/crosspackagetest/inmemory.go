// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package crosspackagetest contains the suite of Fitter. InMemory is the
// subject of the suite.
package crosspackagetest

import (
	"context"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/crosspackage"
	"go.dokimi.dev/dokimi/corpus/go/iface/lang/crosspackage/lib/part"
)

// InMemory implements [crosspackage.Fitter] in memory. It counts the parts
// that it fits within their limits and the senses of each sensor, and returns
// the error of each sense. The corpus applies the suite of Fitter to it as a
// subject.
//
//dokimi:subject crosspackage.Fitter new=NewInMemory
type InMemory struct {
	fitted int
	sensed int
	mu     sync.Mutex
}

// NewInMemory returns an InMemory without fits.
func NewInMemory() *InMemory { return &InMemory{} }

// Fit implements [crosspackage.Fitter]. It returns the error of a cancelled
// context.
func (f *InMemory) Fit(ctx context.Context, p part.Part) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fitted += p.Count
	return nil
}

// Limit implements [crosspackage.Fitter]. It returns the error of a cancelled
// context.
func (f *InMemory) Limit(ctx context.Context, l part.Limits) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fitted = min(f.fitted, l.Count)
	return nil
}

// Sense implements [crosspackage.Fitter].
func (f *InMemory) Sense(_ context.Context, s part.Sensor) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sensed++
	if err := s.Sense(); err != nil {
		return fmt.Errorf("crosspackagetest: the sensor fails: %w", err)
	}
	return nil
}

// The compiler checks that InMemory implements [crosspackage.Fitter].
var _ crosspackage.Fitter = (*InMemory)(nil)
