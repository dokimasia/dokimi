// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package totaltest contains the double and the suite of Mixed. InMemory is the
// control of the suite.
package totaltest

import (
	"context"
	"errors"
	"strings"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/total"
)

// InMemory implements [total.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a control.
//
//dokimi:subject total.Mixed new=NewInMemory
type InMemory struct{}

// The compiler checks that InMemory implements [total.Mixed].
var _ total.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Classify implements [total.Mixed].
func (*InMemory) Classify(ctx context.Context, in string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}

	if strings.TrimSpace(in) == "" {
		return "empty", nil
	}
	return "populated", nil
}

// Normalize implements [total.Mixed].
func (*InMemory) Normalize(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("totaltest: nil context")
	}
	return ctx.Err()
}
