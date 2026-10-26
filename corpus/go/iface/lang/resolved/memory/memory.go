// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package memory is the Applier of the fixture resolved in memory.
package memory

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/resolved"
)

// Applier is an applier in memory. It keeps the limit of the resolution of
// each request that it applied, by the key of the request.
type Applier struct {
	limits map[string]int
	mu     sync.Mutex
}

// New returns an applier without requests.
func New() *Applier { return &Applier{limits: map[string]int{}} }

// Apply implements [resolved.Applier].
func (a *Applier) Apply(_ context.Context, req resolved.ApplyRequest) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.limits[req.Key] = req.Resolution.Limit()
	return nil
}

// Limit implements [resolved.Applier].
func (a *Applier) Limit(_ context.Context, key string) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.limits[key], nil
}

// The compiler checks that Applier implements [resolved.Applier].
var _ resolved.Applier = (*Applier)(nil)
