// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder_test

import (
	"testing"

	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "BadDefault", Give: builder.BadDefault, WantCode: "DOKIMIBUILDER-0001",
			WantMeaning: "the declared default of a field is no value of the type of the field",
		},
		{
			Name: "NotStruct", Give: builder.NotStruct, WantCode: "DOKIMIBUILDER-0002",
			WantMeaning: "an alias under the directive builder refers to no struct of the workspace, " +
				"or the alias or its struct is generic",
		},
		{
			Name: "Unconstructed", Give: builder.Unconstructed, WantCode: "DOKIMIBUILDER-0003",
			WantMeaning: "the constructor of a class under the directive builder takes a parameter " +
				"that a caller must pass",
		},
		{
			Name: "Collision", Give: builder.Collision, WantCode: "DOKIMIBUILDER-0004",
			WantMeaning: "two methods of the builder of a declaration would have one name",
		},
		{
			Name: "BadTag", Give: builder.BadTag, WantCode: "DOKIMIBUILDER-0005",
			WantMeaning: "the tag builder of a field has another value than -",
		},
		{
			Name: "Empty", Give: builder.Empty, WantCode: "DOKIMIBUILDER-0006",
			WantMeaning: "a declaration under the directive builder has no field that a builder can set",
		},
		{
			Name: "Unimportable", Give: builder.Unimportable, WantCode: "DOKIMIBUILDER-0007",
			WantMeaning: "the builder of a declaration refers to a declaration that the file of the builder " +
				"cannot import",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
