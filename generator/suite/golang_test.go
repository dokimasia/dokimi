// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The part of Go declares the target golang with the template of each family,
// and its templates write the suites and the conformance tests in the
// conventions of Go.
func TestGolang(t *testing.T) {
	t.Parallel()

	t.Run("declares the template of each family for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, suite.New(), golang.Target, suiteTemplate, conformTemplate)
	})

	t.Run("writes the golden tree golang", func(t *testing.T) {
		t.Parallel()

		p := suite.New()
		f := fixture(t, goTree, p, storePackage, genericPackage, catalogPackage, slugPackage, memoryPackage,
			cachedPackage, sessionPackage)
		_, res := phases(t, f, p)
		golden.MatchTree(t, "golang", gentest.Rendered(t, res, p, gobackend.New()), golden.ShouldUpdate())
	})
}
