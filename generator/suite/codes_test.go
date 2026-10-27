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
			WantMeaning: "a function under the directive check has no statement, it, its defect, a door of its " +
				"defect, its options or its environment have another signature than a check, or it has the identity " +
				"of a check that the test of its suite or a budget adds",
		},
		{
			Name: "UnknownExcusal", Give: suite.UnknownExcusal, WantCode: "DOKIMISUITE-0006",
			WantMeaning: "the directive excuse names a check that the suite of the declaration of its subject " +
				"does not have",
		},
		{
			Name: "BadBudget", Give: suite.BadBudget, WantCode: "DOKIMISUITE-0007",
			WantMeaning: "the directive budget states a ceiling that no check can apply, a warm-up that no benchmark " +
				"can apply, or a budget that no suite applies",
		},
		{
			Name: "RepeatedSubject", Give: suite.RepeatedSubject, WantCode: "DOKIMISUITE-0008",
			WantMeaning: "a struct has two directives subject of one declaration",
		},
		{
			Name: "Unimportable", Give: suite.Unimportable, WantCode: "DOKIMISUITE-0009",
			WantMeaning: "a callable of a declaration under the directive suite refers to a declaration that the " +
				"file of its suite cannot import",
		},
		{
			Name: "Ungenerated", Give: suite.Ungenerated, WantCode: "DOKIMISUITE-0010",
			WantMeaning: "a type in the inputs of a callable of a declaration under the directive suite has no " +
				"generator",
		},
		{
			Name: "BadOptions", Give: suite.BadOptions, WantCode: "DOKIMISUITE-0011",
			WantMeaning: "the function of the param options of the directive subject is not a function of the " +
				"options of an application",
		},
	}
	for _, c := range codes {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, c)
		})
	}
}
