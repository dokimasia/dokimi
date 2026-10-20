// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package multireturn is the language-axis fixture for more than two return values, past the point where a generator can assume
// a value-and-error pair.
//
// This axis varies the Go type system rather than the classification, so
// these break generators independently of any directive.
package multireturn

import (
	"context"
)

// Wide returns four values. Anything that assumes the last return is an error
// and the first is the value produces wrong field names here.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=multireturntest/
type Wide interface {
	Quad(ctx context.Context, id string) (string, int, bool, error)
	Triple(ctx context.Context, id string) (string, int, bool)
	NoError(ctx context.Context, id string) (string, int)
}
