// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel_test

import (
	"testing"

	"go.dokimi.dev/dokimi/generator/sentinel"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "NoErrors", Give: sentinel.NoErrors, WantCode: "DOKIMISENTINEL-0001",
			WantMeaning: "a package under the directive sentinel has no error value and no error type",
		},
		{
			Name: "OwnNeighbour", Give: sentinel.OwnNeighbour, WantCode: "DOKIMISENTINEL-0002",
			WantMeaning: "a package under the directive sentinel lists itself under the param no-overlap-with",
		},
		{
			Name: "Unloaded", Give: sentinel.Unloaded, WantCode: "DOKIMISENTINEL-0003",
			WantMeaning: "a package under the param no-overlap-with of the directive sentinel is not in the run",
		},
		{
			Name: "Unconstructed", Give: sentinel.Unconstructed, WantCode: "DOKIMISENTINEL-0004",
			WantMeaning: "an error type under the directive sentinel has no public constructor that a check can call",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
