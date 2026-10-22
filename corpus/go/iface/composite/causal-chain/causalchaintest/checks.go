// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package causalchaintest contains the suite and the double of Log, and the
// checks that a consumer adds to the suite.
package causalchaintest

import (
	"context"
	"slices"

	"go.dokimi.dev/assert"

	causalchain "go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain"
)

// KeepsTheBody checks that Replay returns the body that Append admitted.
//
//dokimi:check Log.Replay defect=losesBodies
func KeepsTheBody(h *LogHandle) {
	s := h.Subject()
	assert.NoError(h, s.Append(h.Context(), causalchain.RootEntry()), "Append admits the root entry")
	got, err := s.Replay(h.Context())
	assert.NoError(h, err, "Replay reads the log")
	assert.Length(h, got, 1, "Replay returns the root entry")
	assert.Equal(h, got[0].Body, causalchain.RootEntry().Body, h.Statement())
}

// losesBodies returns a Log that admits each entry without its body.
func losesBodies(h *LogHandle) causalchain.Log {
	var entries []causalchain.Entry
	return NewLogStub(h,
		WithLogAppend(func(_ context.Context, e causalchain.Entry) error {
			e.Body = ""
			entries = append(entries, e)
			return nil
		}),
		WithLogReplay(func(context.Context) ([]causalchain.Entry, error) {
			return slices.Clone(entries), nil
		}))
}
