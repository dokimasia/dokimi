// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sampletest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package sampletest

import (
	"context"
	"errors"
	"strings"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/sample"
)

// ErrUnprocessable is the error of Process for an input that NewInput does
// not return.
var ErrUnprocessable = errors.New("sampletest: input is not sampled")

// inputPrefix is the prefix of each input that sample.NewInput returns.
const inputPrefix = "sample:"

// InMemory implements [sample.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a subject.
//
//dokimi:subject sample.Mixed new=NewInMemory
type InMemory struct{}

// The compiler checks that InMemory implements [sample.Mixed].
var _ sample.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

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
