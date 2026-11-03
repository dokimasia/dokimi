// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sample is the mixin-axis fixture for the sample mixin, which
// declares the builder of the inputs of each derived call of a callable.
//
// Process accepts only an input that its builder returns, so a check that
// drew its input would fail every call. The builder is a function of the
// package of the declaration, which the file of the suite imports.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package sample

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=sampletest/
type Mixed interface {
	// Process returns the processed form of an input that NewInput returns.
	// It fails for any other input.
	//
	//dokimi:mixin sample builder=NewInput
	Process(ctx context.Context, input string) (string, error)
}

// NewInput returns the input of each derived call of Process.
func NewInput() (input string) { return "sample:a" }
