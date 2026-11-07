// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"

	"go.dokimi.dev/dokimi/generator/suite"
)

func TestWorld(t *testing.T) {
	t.Parallel()

	vaults := generate(t, goTree, vaultPackage).suites
	worldOf := func(t *testing.T, declaration string) *suite.World {
		t.Helper()
		w := vaults[declaration].World
		assert.NotNil(t, w, "the suite of "+declaration+" has a world")
		return w
	}

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("names the world and the function of its operations after the declaration", func(t *testing.T) {
			t.Parallel()

			w := worldOf(t, "Vault")
			expect.Equal(t, []string{w.Name, w.Operations}, []string{"VaultWorld", "VaultOperations"},
				"the variable and the function of the world open with the declaration")
		})

		t.Run("reads the declared error of a corrupted record from the mixin durable", func(t *testing.T) {
			t.Parallel()

			w := worldOf(t, "Vault")
			assert.NotNil(t, w.Corrupt, "the world has the error that the param corrupt refers to")
			assert.NotNil(t, w.Corrupt.Value, "the declared error is a value")
			expect.Equal(
				t,
				[]any{w.Corrupt.Value.Kind, w.Corrupt.Value.Callee.Name},
				[]any{emit.ValueName, "ErrCorrupt"},
				"the declared error is the name ErrCorrupt",
			)
		})

		t.Run("leaves out the declared error where the mixin durable names none", func(t *testing.T) {
			t.Parallel()

			expect.Nil(t, worldOf(t, "Journal").Corrupt, "the mixin durable of Append has no param corrupt")
		})

		clients := []struct {
			name string
			give string
			want int
		}{
			{name: "takes three clients where each operation has the mixin concurrent", give: "Registry", want: 3},
			{name: "takes one client where an operation lacks the mixin concurrent", give: "Journal", want: 0},
			{name: "takes one client where no operation has the mixin concurrent", give: "Vault", want: 0},
		}
		for _, tt := range clients {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				expect.Equal(
					t,
					worldOf(t, tt.give).Clients,
					tt.want,
					"the world has the clients that its operations admit",
				)
			})
		}

		t.Run("leaves the world out of a suite without a check of a world probe", func(t *testing.T) {
			t.Parallel()

			short := generate(t, goTree, sealedPackage).suites["Short"]
			expect.Nil(t, short.World, "no callable of Short has the mixin durable")
		})
	})
}
