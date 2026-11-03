// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package deadchecktest_test

import (
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mutation/deadcheck/deadchecktest"
)

// The spec of Len is the only test that covers Len, and it states a bound
// alone. The mutants of Len survive it, so a mutation record lists them as
// survivors that no check covered. The spec builds its store without
// NewInMemory, so a mutant of the constructor does not end its test binary.
func TestInMemory(t *testing.T) {
	t.Parallel()

	t.Run("Len", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a number that is not negative", func(t *testing.T) {
			t.Parallel()

			assert.True(t, (&deadchecktest.InMemory{}).Len() >= 0, "the length of a store is not negative")
		})
	})
}
