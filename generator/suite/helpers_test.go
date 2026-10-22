// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

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
	sessionPackage  = "example.com/shop/session"
	suiteTemplate   = "suite.tmpl"
	conformTemplate = "conform.test.tmpl"
)

// generated is what suite records over one fixture: the suite of each
// declaration by its name, the record of each test of the family of the tests
// by the name of the test, the origins that test their package from inside,
// and the findings of the run.
type generated struct {
	suites   map[string]suite.Suite
	tests    map[string]gentest.Record[suite.Conformance]
	beside   map[string][]string
	inside   map[symbol.Identity]bool
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
	return recorded(t, fixture(t, dir, p, pkgs...), p)
}

// ref returns the reference of a record to the function named name of the
// fixture f.
func ref(t *testing.T, f *gentest.Fixture, name string) *emit.TypeRef {
	t.Helper()

	return &emit.TypeRef{Spelling: name, Target: f.Find(t, symbol.KindFunction, "", name).Identity()}
}

// phases runs the annotate half and then the generate half of the suite p over
// the fixture f, and returns the result of each.
func phases(t *testing.T, f *gentest.Fixture, p plugin.Plugin) (stamps, records plugintest.Result) {
	t.Helper()

	stamps = f.Annotate(t, p)
	assert.NoError(t, stamps.Err, "the annotate half runs whole")
	records = f.Generate(t, p)
	assert.NoError(t, records.Err, "the generate half runs whole")
	return stamps, records
}

// recorded runs both halves of the suite p over the fixture f, and returns
// what they record and report.
func recorded(t *testing.T, f *gentest.Fixture, p plugin.Plugin) generated {
	t.Helper()

	stamps, res := phases(t, f, p)
	got := generated{
		suites:   map[string]suite.Suite{},
		inside:   map[symbol.Identity]bool{},
		findings: slices.Concat(slices.Collect(stamps.Sink.All()), slices.Collect(res.Sink.All())),
	}
	for u := range res.Emit.Units() {
		for _, origin := range u.Inside {
			got.inside[origin] = true
		}
	}
	declaration := func(s suite.Suite) string { return s.Declaration }
	for name, r := range gentest.Records(t, res, suiteTemplate, "", declaration) {
		got.suites[name] = r.Data
	}
	got.beside = gentest.Beside(res, declaration)
	got.tests = gentest.Records(t, res, conformTemplate, suite.TestTag,
		func(c suite.Conformance) string { return c.Test })
	return got
}
