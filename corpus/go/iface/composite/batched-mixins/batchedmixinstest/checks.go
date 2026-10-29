// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package batchedmixinstest contains the suite and the double of Batched, the
// Batched in memory, and the checks that a consumer adds to the suite.
package batchedmixinstest

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"go.dokimi.dev/assert"

	batchedmixins "go.dokimi.dev/dokimi/corpus/go/iface/composite/batched-mixins"
)

// ListsAtMostTheLimit checks that List returns at most 50 keys after 51
// writes.
//
//dokimi:check Batched.List defect=unbounded
func ListsAtMostTheLimit(h *BatchedHandle) {
	s := h.Subject()
	for i := range listLimit + 1 {
		assert.NoError(h, s.Put(h.Context(), fmt.Sprintf("key-%02d", i), "value"), "Put stores the key")
	}
	got, err := s.List(h.Context())
	assert.NoError(h, err, "List reads the keys")
	assert.InRange(h, len(got), 0, listLimit, h.Statement())
}

// unbounded returns a Batched whose List returns every stored key. Each call
// returns the error of a cancelled context.
func unbounded(h *BatchedHandle) batchedmixins.Batched {
	values := map[string]string{}
	return NewBatchedStub(h,
		WithBatchedPut(func(ctx context.Context, key, value string) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			values[key] = value
			return nil
		}),
		WithBatchedRead(func(ctx context.Context, key string) (string, error) {
			if err := ctx.Err(); err != nil {
				return "", err
			}
			return values[key], nil
		}),
		WithBatchedList(func(ctx context.Context) ([]string, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			return slices.Sorted(maps.Keys(values)), nil
		}))
}
