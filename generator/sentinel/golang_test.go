// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel_test

import (
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	gotesting "go.dokimi.dev/eidos/lang/go/testing"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/toolchain"

	"go.dokimi.dev/dokimi/generator/sentinel"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The part of Go declares the target golang with the template of the family
// and its helper, and opens each message with the name of its package.
func TestGolang(t *testing.T) {
	t.Parallel()

	t.Run("declares the template of the family for the target golang", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, sentinel.New(), golang.Target, testTemplate)
	})

	t.Run("writes the checks of the golden tree golang", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "golang", rendered(t), golden.ShouldUpdate())
	})

	t.Run("writes Go that vets and whose tests pass", func(t *testing.T) {
		t.Parallel()

		a := gotesting.New()
		if !toolchain.Require(t, a) {
			return
		}
		project := gentest.Project(rendered(t), goTree, shopModule)
		gotesting.AssertVets(t.Context(), t, project)
		toolchain.AssertTestsPass(t.Context(), t, a, project)
	})

	t.Run("writes the name of a package in Pascal case", func(t *testing.T) {
		t.Parallel()

		templates, provides := sentinel.New().(plugin.TemplateProvider)
		assert.True(t, provides, "sentinel declares templates")
		pascal, declared := templates.TemplateFuncs(golang.Target)["pascal"].(func(string) string)
		assert.True(t, declared, "the target golang declares the helper pascal")
		expect.Equal(t, pascal("store"), "Store", "the test of the package store is TestStoreSentinels")
	})
}

// rendered runs sentinel over the packages of the Go tree whose checks build,
// and returns the files that the Go backend renders from them.
func rendered(t *testing.T) fstest.MapFS {
	t.Helper()

	p := sentinel.New()
	res := fixture(t, golang.Target, goTree, p, storePackage, neighbourPackage, embeddedPackage, noprefixPackage,
		splitPackage, typedPackage).Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	return gentest.Rendered(t, res, p, gobackend.New())
}
