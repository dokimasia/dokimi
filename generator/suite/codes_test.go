// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package registers its codes under its own prefix when it initializes.
func TestCodes(t *testing.T) {
	t.Parallel()

	codes := []gentest.Code{
		{
			Name: "Uninstantiable", Give: suite.Uninstantiable, WantCode: "DOKIMISUITE-0001",
			WantMeaning: "a declaration under the directive suite has type parameters that its suite " +
				"cannot instantiate",
		},
		{
			Name: "NoSuite", Give: suite.NoSuite, WantCode: "DOKIMISUITE-0002",
			WantMeaning: "the declaration of the directive subject has no directive suite",
		},
		{
			Name: "BadFactory", Give: suite.BadFactory, WantCode: "DOKIMISUITE-0003",
			WantMeaning: "the function of the param new of the directive subject is not a factory of subjects",
		},
		{
			Name: "UnknownCallable", Give: suite.UnknownCallable, WantCode: "DOKIMISUITE-0004",
			WantMeaning: "the directive check has a callable that no declaration under the directive suite " +
				"has, or that two of them have",
		},
		{
			Name: "BadCheck", Give: suite.BadCheck, WantCode: "DOKIMISUITE-0005",
			WantMeaning: "a function under the directive check has no statement, or it or its defect has " +
				"another signature than a check",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
