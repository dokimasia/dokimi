// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package rootedtest contains the suite of Tally and the check that a
// consumer adds to the suite. The kit does not import the composition root,
// so the planted defect of the check wraps the reference of the handle.
package rootedtest

import (
	"context"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted"
)

// incs is the number of calls of Inc that the check makes on the subject and
// on the reference.
const incs = 3

// CountsEachInc checks that Count returns the count of the reference after
// the same calls of Inc.
//
//dokimi:check rooted.Tally.Count defect=skipsAnInc
func CountsEachInc(h *TallyHandle) {
	s, ref := h.Subject(), h.Reference()
	for range incs {
		assert.NoError(h, s.Inc(h.Context()), "Inc counts on the subject")
		assert.NoError(h, ref.Inc(h.Context()), "Inc counts on the reference")
	}
	got, err := s.Count(h.Context())
	assert.NoError(h, err, "Count reads the subject")
	want, err := ref.Count(h.Context())
	assert.NoError(h, err, "Count reads the reference")
	assert.Equal(h, got, want, h.Statement())
}

// skipsAnInc returns the planted defect of CountsEachInc: a tally over the
// reference of the handle that does not count the second call of Inc.
func skipsAnInc(h *TallyHandle) rooted.Tally {
	return &skipping{Tally: h.Reference()}
}

// skipping is a tally that does not count the second call of Inc.
type skipping struct {
	rooted.Tally
	calls int
}

// Inc adds one to the count, except on its second call.
func (s *skipping) Inc(ctx context.Context) error {
	s.calls++
	if s.calls == 2 {
		return nil
	}
	//nolint:wrapcheck // the defect reports the errors of the reference as the reference reports them
	return s.Tally.Inc(ctx)
}
