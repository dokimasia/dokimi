// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package dynamictest contains the suite of Diary. InMemory is the control of
// the suite.
package dynamictest

import (
	"math/big"
	"sync"
	"time"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/dynamic"
)

// InMemory implements [dynamic.Diary] in memory. The corpus applies the suite
// of Diary to it as a control.
//
//dokimi:subject dynamic.Diary new=NewInMemory
type InMemory struct {
	zone    *time.Location
	scale   *big.Int
	records []string
	mu      sync.Mutex
}

// NewInMemory returns an empty InMemory in the zone UTC at the scale 1.
func NewInMemory() *InMemory { return &InMemory{zone: time.UTC, scale: big.NewInt(1)} }

// Print implements [dynamic.Diary].
func (d *InMemory) Print(msg string, _ ...any) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.records = append(d.records, msg)
}

// Field implements [dynamic.Diary].
func (d *InMemory) Field(name string, _ any) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.records = append(d.records, name)
}

// In implements [dynamic.Diary].
func (d *InMemory) In(loc *time.Location) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.zone = loc
}

// Scale implements [dynamic.Diary].
func (d *InMemory) Scale(n *big.Int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.scale = n
}

// Count implements [dynamic.Diary].
func (d *InMemory) Count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.records)
}

// The compiler checks that InMemory implements [dynamic.Diary].
var _ dynamic.Diary = (*InMemory)(nil)
