// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package codeclossytest contains the double and the suite of Contract.
// InMemory is the subject of the suite.
package codeclossytest

import (
	"context"
	"errors"
	"strings"

	codeclossy "go.dokimi.dev/dokimi/corpus/go/iface/contract/codec-lossy"
)

// InMemory implements [codeclossy.Contract] in memory. The corpus applies
// the suite of Contract to it as a subject.
//
//dokimi:subject codeclossy.Contract new=NewInMemory
type InMemory struct{}

// The compiler checks that InMemory implements [codeclossy.Contract].
var _ codeclossy.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Encode implements [codeclossy.Contract].
func (*InMemory) Encode(ctx context.Context, in string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	return strings.ToLower(in), nil
}

// Decode implements [codeclossy.Contract].
func (*InMemory) Decode(ctx context.Context, in string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	return in, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("codeclossytest: nil context")
	}
	return ctx.Err()
}
