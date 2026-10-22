// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package memory_test

import (
	"testing"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl/memory"
	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl/ttltest"
)

// The store in memory passes the suite of Mixed on the clock of each check,
// because the check of the lifetime advances that clock.
func TestStore(t *testing.T) {
	t.Parallel()

	ttltest.RunMixedSuiteFromHandle(t, func(h *ttltest.MixedHandle) *memory.Store {
		return memory.New(h.Clock())
	})
}
