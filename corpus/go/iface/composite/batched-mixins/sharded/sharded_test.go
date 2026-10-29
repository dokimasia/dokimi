// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sharded_test

import (
	"context"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"

	batchedmixins "go.dokimi.dev/dokimi/corpus/go/iface/composite/batched-mixins"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/batched-mixins/batchedmixinstest"
	"go.dokimi.dev/dokimi/corpus/go/iface/composite/batched-mixins/sharded"
)

// The sharded Batched passes the suite of Batched and a check against a
// reference. The Batched in memory of the suite is the reference, because it
// keeps every key in one map.
func TestBatched(t *testing.T) {
	t.Parallel()

	batchedmixinstest.RunBatchedSuite(t, func() *sharded.Batched { return sharded.New(4) },
		conform.Reference(func(*batchedmixinstest.BatchedHandle) batchedmixins.Batched {
			return batchedmixinstest.NewInMemory()
		}),
		conform.With(listsAsTheReference()))
}

// listsAsTheReference returns the check that List returns the keys that the
// reference lists after the same writes. Its planted defect is the sharded
// Batched whose List returns the first key alone.
func listsAsTheReference() batchedmixinstest.BatchedCheck {
	return batchedmixinstest.BatchedCheck{
		ID:        batchedmixinstest.BatchedChecks.List.Check("lists-as-the-reference"),
		Statement: "List returns the keys that the reference lists after the same writes",
		Apply: func(h *batchedmixinstest.BatchedHandle) {
			s, reference := h.Subject(), h.Reference()
			for _, key := range []string{"cherry", "apple", "date", "banana"} {
				assert.NoError(h, s.Put(h.Context(), key, "fruit"), "Put stores the key in the subject")
				assert.NoError(h, reference.Put(h.Context(), key, "fruit"), "Put stores the key in the reference")
			}
			want, err := reference.List(h.Context())
			assert.NoError(h, err, "List reads the keys of the reference")
			got, err := s.List(h.Context())
			assert.NoError(h, err, "List reads the keys of the subject")
			assert.Equal(h, got, want, h.Statement())
		},
		Defect: func(h *batchedmixinstest.BatchedHandle) batchedmixins.Batched {
			inner := sharded.New(4)
			return batchedmixinstest.NewBatchedStub(h, batchedmixinstest.BatchedStubDelegateTo(inner),
				batchedmixinstest.WithBatchedList(func(ctx context.Context) ([]string, error) {
					keys, err := inner.List(ctx)
					return keys[:min(1, len(keys))], err
				}))
		},
	}
}
