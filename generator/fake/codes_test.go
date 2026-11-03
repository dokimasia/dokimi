// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fake_test

import (
	"testing"

	"go.dokimi.dev/dokimi/generator/fake"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "NoSuite", Give: fake.NoSuite, WantCode: "DOKIMIFAKE-0001",
			WantMeaning: "an interface under the directive fake has no directive suite",
		},
		{
			Name: "BadAbsent", Give: fake.BadAbsent, WantCode: "DOKIMIFAKE-0002",
			WantMeaning: "the param absent of the directive fake is no name template of a sentinel",
		},
		{
			Name: "Refused", Give: fake.Refused, WantCode: "DOKIMIFAKE-0003",
			WantMeaning: "an interface under the directive fake has a callable that no reference form covers",
		},
		{
			Name: "Unimplementable", Give: fake.Unimplementable, WantCode: "DOKIMIFAKE-0004",
			WantMeaning: "a fake cannot implement an interface under the directive fake, or its test cannot refer " +
				"to the interface",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
