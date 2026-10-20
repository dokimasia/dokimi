// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sampled is the language-axis fixture for AUTHORED VALUES: a
// //dokimi:sample names the function a generated check writes a type
// with, where the derived value would break the run — a tag the
// subject validates against a vocabulary its structure does not state.
//
// The directive is load-bearing here, not decorative. Apply refuses a
// tag outside the vocabulary, a refused first call leaves the
// idempotent claim vacuous, and the census fails a run whose claim
// never engaged — so without the authored sample the generated suite
// cannot be green against a validating subject.
//
// Level carries both value directives at once, so the generated
// fixture also pins their precedence: the sample outranks the default,
// because the fixture is what checks write and that is the sample
// directive's half of the composition — defaults seed a constructor.
package sampled

import (
	"context"
	"errors"
)

// ErrUnknownTag is what Apply reports for a tag outside the subject's
// vocabulary — the refusal that makes a derived "test-apply" value
// unable to engage the idempotent claim.
var ErrUnknownTag = errors.New("sampled: unknown tag")

// Tag is a validated vocabulary word. A derived sample satisfies the
// compiler and fails the vocabulary, which is the case the sample
// directive exists for.
//
//dokimi:sample value=NewTestTag()
type Tag string

// NewTestTag is the function the sample directive names: a value the
// vocabulary admits.
func NewTestTag() Tag { return "release" }

// Level carries both directives, with different values, so the
// generated fixture proves which one wins.
//
//dokimi:sample value=NewTestLevel()
//dokimi:default 3
type Level int

// NewTestLevel is the authored value the fixture must draw over the
// declared default.
func NewTestLevel() Level { return 7 }

// Labeller is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=sampledtest/
type Labeller interface {
	// Apply is idempotent for a known tag and refuses an unknown one
	// with [ErrUnknownTag]. The refusal is the point: a refused first
	// call leaves the idempotent claim vacuous, so only an authored
	// value the vocabulary admits lets the claim engage.
	//dokimi:mixin idempotent
	Apply(ctx context.Context, tag Tag, level Level) error
}
