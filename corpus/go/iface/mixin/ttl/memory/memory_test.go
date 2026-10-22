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

	ttltest.RunMixedSuiteFromHandle(t, newStore)
}

// The store in memory passes the allocation ceiling of Read, which the
// declaration of Mixed states.
func TestStoreAllocs(t *testing.T) {
	ttltest.RunMixedBudgets(t, newStore)
}

// BenchmarkStore measures each callable of Mixed on the store in memory,
// under the ceilings of the budget of Read.
func BenchmarkStore(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		ttltest.BenchMixedSuite(b, ttltest.MixedChecks.Put.Callable, newStore)
	})
	b.Run("Read", func(b *testing.B) {
		ttltest.BenchMixedSuite(b, ttltest.MixedChecks.Read.Callable, newStore)
	})
}

// newStore returns an empty store on the clock of the handle h.
func newStore(h *ttltest.MixedHandle) *memory.Store {
	return memory.New(h.Clock())
}
