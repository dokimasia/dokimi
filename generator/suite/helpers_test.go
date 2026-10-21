// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees of the specs, the packages that two specs or more read, and the
// templates of the two families.
const (
	goTree          = "testdata/golang"
	storePackage    = "example.com/shop/store"
	genericPackage  = "example.com/shop/generic"
	catalogPackage  = "example.com/shop/catalog"
	slugPackage     = "example.com/shop/slug"
	memoryPackage   = "example.com/shop/memory"
	cachedPackage   = "example.com/shop/cached"
	suiteTemplate   = "suite.tmpl"
	conformTemplate = "conform.test.tmpl"
)

// generated is what suite records over one fixture: the suite of each
// declaration, the record of the conformance test of each subject, by their
// names, and the findings of the run.
type generated struct {
	suites   map[string]suite.Suite
	tests    map[string]gentest.Record[suite.Conformance]
	beside   map[string][]string
	findings []diag.Diag
}

// fixture returns a fixture over the packages pkgs of the tree dir with the
// plugin suite p and the target golang.
func fixture(tb assert.TB, dir string, p plugin.Plugin, pkgs ...string) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{Dir: dir, Packages: pkgs, Target: golang.Target, Plugins: []plugin.Plugin{p}})
}

// generate runs a new suite over the packages pkgs of the tree dir, and
// returns what it records. It checks that each family claims its files
// through its template.
func generate(t *testing.T, dir string, pkgs ...string) generated {
	t.Helper()

	p := suite.New()
	return recorded(t, fixture(t, dir, p, pkgs...).Generate(t, p))
}

// recorded returns what the phase call res of suite records.
func recorded(t *testing.T, res plugintest.Result) generated {
	t.Helper()

	assert.NoError(t, res.Err, "the generator runs whole")
	got := generated{suites: map[string]suite.Suite{}, findings: slices.Collect(res.Sink.All())}
	declaration := func(s suite.Suite) string { return s.Declaration }
	for name, r := range gentest.Records(t, res, suiteTemplate, "", declaration) {
		got.suites[name] = r.Data
	}
	got.beside = gentest.Beside(res, declaration)
	got.tests = gentest.Records(t, res, conformTemplate, suite.TestTag,
		func(c suite.Conformance) string { return c.Subject })
	return got
}
