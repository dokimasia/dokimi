// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package codectest contains the double and the suite of Contract. InMemory is
// the control of the suite.
package codectest

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/codec"
)

// InMemory implements [codec.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject codec.Contract new=NewInMemory
type InMemory struct{}

// The compiler checks that InMemory implements [codec.Contract].
var _ codec.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Encode implements [codec.Contract].
func (*InMemory) Encode(ctx context.Context, in string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(in)), nil
}

// Decode implements [codec.Contract].
func (*InMemory) Decode(ctx context.Context, in string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", fmt.Errorf("codectest: decode: %w", err)
	}
	return string(raw), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("codectest: nil context")
	}
	return ctx.Err()
}
