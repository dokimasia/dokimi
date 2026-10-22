// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package replica_test

import (
	"testing"

	"go.dokimi.dev/assert"

	causalchain "go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/memory"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/causal-chain/replica"
)

// A replica replays the entries of its source and does not admit an entry of
// its own.
func TestLog(t *testing.T) {
	t.Parallel()

	t.Run("Follow", func(t *testing.T) {
		t.Parallel()

		t.Run("replays the entries that the source admitted", func(t *testing.T) {
			t.Parallel()

			source := memory.New()
			assert.NoError(t, source.Append(t.Context(), causalchain.RootEntry()), "the source admits the root entry")
			got, err := replica.Follow(source).Replay(t.Context())
			assert.NoError(t, err, "the replica replays its source")
			assert.Length(t, got, 1, "the replica returns the entry of its source")
		})
	})

	t.Run("Append", func(t *testing.T) {
		t.Parallel()

		t.Run("refuses the entry with ErrReadOnly", func(t *testing.T) {
			t.Parallel()

			err := replica.New().Append(t.Context(), causalchain.RootEntry())
			assert.ErrorIs(t, err, replica.ErrReadOnly, "a replica admits no entry")
		})
	})
}
