// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub_test

import (
	"testing"

	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "Constraint", Give: stub.Constraint, WantCode: "DOKIMISTUB-0001",
			WantMeaning: "a double cannot implement an interface that only constrains type parameters",
		},
		{
			Name: "IncompleteMembers", Give: stub.IncompleteMembers, WantCode: "DOKIMISTUB-0002",
			WantMeaning: "a double cannot implement a declaration whose members are incomplete",
		},
		{
			Name: "Unimplementable", Give: stub.Unimplementable, WantCode: "DOKIMISTUB-0003",
			WantMeaning: "a double cannot implement a member of its declaration, or cannot refer to the declaration",
		},
		{
			Name: "Collision", Give: stub.Collision, WantCode: "DOKIMISTUB-0004",
			WantMeaning: "a double would have two members of one name",
		},
		{
			Name: "Empty", Give: stub.Empty, WantCode: "DOKIMISTUB-0005",
			WantMeaning: "a double of a declaration without a member implements nothing",
		},
		{
			Name: "Incomparable", Give: stub.Incomparable, WantCode: "DOKIMISTUB-0006",
			WantMeaning: "a partition fault needs a parameter whose arguments compare with a key",
		},
		{
			Name: "BadField", Give: stub.BadField, WantCode: "DOKIMISTUB-0007",
			WantMeaning: "the param field of a double needs a name template that uses {{.Method}}",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
