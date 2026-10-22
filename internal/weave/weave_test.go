// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package weave_test

import (
	"testing"

	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/weave"
)

// The fact of a double is a key of suite, so the gate of stub binds where suite
// claims its namespace.
func TestWeave(t *testing.T) {
	t.Parallel()

	t.Run("KeyDouble", func(t *testing.T) {
		t.Parallel()

		t.Run("is in the namespace of the keys of suite", func(t *testing.T) {
			t.Parallel()

			expect.Equal(t, weave.KeyDouble.Namespace(), suite.KeyDeclaration.Namespace(),
				"the key of the double has the namespace that suite claims")
		})
	})
}
