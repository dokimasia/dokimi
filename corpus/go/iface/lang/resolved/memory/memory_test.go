// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package memory_test

import (
	"context"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/resolved"
	"go.dokimi.dev/dokimi/corpus/go/iface/lang/resolved/memory"
	"go.dokimi.dev/dokimi/corpus/go/iface/lang/resolved/resolvedtest"
)

// The applier in memory passes the suite of Applier with the resolutions of a
// registry, and a check of its own that reads them. The application falsifies
// that check on its planted defect first, with the same resolutions.
func TestApplier(t *testing.T) {
	t.Parallel()

	registry := resolved.NewRegistry(3)
	resolvedtest.RunApplierSuite(t, memory.New,
		conform.Using[resolved.Applier](prop.Just(registry.Resolve())), conform.With(keepsTheLimit()))
}

// keepsTheLimit returns the check that Limit returns the limit of the
// resolution of the request that Apply applied. Its planted defect is the
// applier in memory whose Limit returns 0. The defect passes the check under
// the zero resolution of the suite, so the check needs the resolutions of the
// application.
func keepsTheLimit() resolvedtest.ApplierCheck {
	return resolvedtest.ApplierCheck{
		ID:        resolvedtest.ApplierChecks.Limit.Check("keeps-the-limit"),
		Statement: "Limit returns the limit of the resolution of the request that Apply applied",
		Property: func(h *resolvedtest.ApplierHandle, c *prop.Case) {
			req := c.Draw(prop.Of[resolved.ApplyRequest](), "req")
			s := h.Subject()
			assert.NoError(h, s.Apply(h.Context(), req), "Apply applies the request")
			got, err := s.Limit(h.Context(), req.Key)
			assert.NoError(h, err, "Limit reads the limit of the request")
			assert.Equal(h, got, req.Resolution.Limit(), h.Statement())
		},
		Defect: func(h *resolvedtest.ApplierHandle) resolved.Applier {
			return resolvedtest.NewApplierStub(h, resolvedtest.ApplierStubDelegateTo(memory.New()),
				resolvedtest.WithApplierLimit(func(ctx context.Context, _ string) (int, error) { return 0, ctx.Err() }))
		},
	}
}
