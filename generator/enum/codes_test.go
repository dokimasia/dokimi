// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum_test

import (
	"testing"

	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "Empty", Give: enum.Empty, WantCode: "DOKIMIENUM-0001",
			WantMeaning: "an enumeration under the directive enum has no variant",
		},
		{
			Name: "Duplicate", Give: enum.Duplicate, WantCode: "DOKIMIENUM-0002",
			WantMeaning: "two variants of an enumeration under the directive enum have one value or one text",
		},
		{
			Name: "Conflict", Give: enum.Conflict, WantCode: "DOKIMIENUM-0003",
			WantMeaning: "a declaration of the package of an enumeration under the directive enum has the name " +
				"of a member of its surface and another form",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
