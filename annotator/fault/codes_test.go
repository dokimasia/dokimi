// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fault_test

import (
	"testing"

	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "Infallible", Give: fault.Infallible, WantCode: "DOKIMIFAULT-0001",
			WantMeaning: "a fault of a double needs a method that can fail",
		},
		{
			Name: "BadRetry", Give: fault.BadRetry, WantCode: "DOKIMIFAULT-0002",
			WantMeaning: "a retry schedule lets an attempt from the second on succeed",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
