// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ttltest contains the suite and the double of Mixed, and the checks
// that a consumer adds to the suite.
package ttltest

import (
	"context"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl"
)

// lifetime is the duration of the mixin ttl of Read.
const lifetime = time.Minute

// ExpiresAfterTheLifetime checks that Read reports ErrExpired once the
// lifetime of the value has passed.
//
//dokimi:check ttl.Mixed.Read defect=keepsForever
func ExpiresAfterTheLifetime(h *MixedHandle) {
	s := h.Subject()
	v := ttl.Value{Key: "session", Body: "token"}
	assert.NoError(h, s.Put(h.Context(), v), "Put stores the value")
	got, err := s.Read(h.Context(), v.Key)
	assert.NoError(h, err, "Read returns the value within its lifetime")
	assert.Equal(h, got, v, "Read returns the value that Put stored")
	h.Advance(lifetime)
	_, err = s.Read(h.Context(), v.Key)
	assert.ErrorIs(h, err, ttl.ErrExpired, h.Statement())
}

// keepsForever returns a Mixed on the clock of the check whose Read returns
// each value after its lifetime.
func keepsForever(h *MixedHandle) ttl.Mixed {
	values := map[string]ttl.Value{}
	return NewMixedStub(h, MixedStubWithClock(h.Clock()),
		WithMixedPut(func(_ context.Context, v ttl.Value) error {
			values[v.Key] = v
			return nil
		}),
		WithMixedRead(func(_ context.Context, key string) (ttl.Value, error) {
			v, found := values[key]
			if !found {
				return ttl.Value{}, ttl.ErrExpired
			}
			return v, nil
		}))
}
