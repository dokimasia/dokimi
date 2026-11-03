// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fake_test

import (
	"maps"
	"strings"
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"
	"go.dokimi.dev/eidos/sdk/toolchain"

	"go.dokimi.dev/dokimi/generator/fake"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// shopModule is the module of the Go tree.
const shopModule = "example.com/shop"

// The part of Go declares the target golang with the templates of the two
// families. Its templates write the fakes and their tests in the conventions
// of Go.
func TestGolang(t *testing.T) {
	t.Parallel()

	t.Run("declares the templates of the families for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, fake.New(), golang.Target, fakeTemplate, testTemplate)
	})

	t.Run("writes the tests of the fakes into the external test package", func(t *testing.T) {
		t.Parallel()

		gentest.External(t, fake.New(), fake.TestTag)
	})

	t.Run("writes the fakes of the golden tree golang", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "golang", rendered(t), golden.ShouldUpdate())
	})

	t.Run("writes fakes that pass go vet", func(t *testing.T) {
		t.Parallel()

		if !toolchain.Require(t, gotesting.New()) {
			return
		}
		files := rendered(t)
		// The tests of the fakes call the suites, whose module the tree does
		// not require, so go vet checks the fakes alone.
		maps.DeleteFunc(files, func(path string, _ *fstest.MapFile) bool { return strings.HasSuffix(path, "_test.go") })
		gotesting.AssertVets(t.Context(), t, gentest.Project(files, goTree, shopModule))
	})
}

// rendered runs both halves of suite and then fake over the packages of the
// Go tree whose fakes render, and returns the files that the Go backend
// renders from the records of fake.
func rendered(t *testing.T) fstest.MapFS {
	t.Helper()

	r := generate(t, storePackage, cachePackage, logPackage, genericPackage, clashPackage, quietPackage)
	return gentest.Rendered(t, r.f, r.res, r.fake, gobackend.New(), r.suite)
}
