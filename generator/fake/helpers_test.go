// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fake_test

import (
	"testing"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	goannotate "go.dokimi.dev/eidos/lang/go/annotate"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	"go.dokimi.dev/eidos/plugin/shape/catalog"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/authored"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/generator/fake"
	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The tree of the specs, the packages that two specs or more read, and the
// templates of the two families.
const (
	goTree         = "testdata/golang"
	storePackage   = "example.com/shop/store"
	cachePackage   = "example.com/shop/cache"
	logPackage     = "example.com/shop/log"
	genericPackage = "example.com/shop/generic"
	clashPackage   = "example.com/shop/clash"
	quietPackage   = "example.com/shop/quiet"
	fakeTemplate   = "fake.tmpl"
	testTemplate   = "fake.test.tmpl"
)

// run is one run of suite and fake over a fixture: the fixture, the plugins
// and the result of the phase call of fake.
type run struct {
	f     *gentest.Fixture
	suite sdk.Dual
	fake  plugin.Generator
	res   plugintest.Result
}

// generate runs both halves of a new suite and then a new fake over the
// packages pkgs of the tree of the specs, as a composition runs the suite of a
// declaration before its fake. The fixture has the Go backend as its spoke,
// so a phase call previews the package of each routed file.
func generate(t *testing.T, pkgs ...string) run {
	t.Helper()

	return phases(t, gobackend.New(), pkgs...)
}

// phases runs both halves of a new suite and then a new fake over the
// packages pkgs of the tree of the specs with the spoke types, which is nil
// for a plan that previews no package. The annotators of the composition stamp
// the fixture: the authored samples and witnesses, the facts of Go, the
// declared defaults, the names of the builders, the faults of the methods and
// the classifications of the shape catalog.
func phases(t *testing.T, types plugin.TypeSpeller, pkgs ...string) run {
	t.Helper()

	s, p := suite.New(), fake.New()
	f := gentest.New(t, &gentest.Spec{
		Dir: goTree, Packages: pkgs, Target: golang.Target, Types: types, Plugins: []plugin.Plugin{s, p},
		Annotators: append([]plugin.Annotator{
			authored.Sample(), authored.Witness(), goannotate.New(), defaults.New(), builder.New(), fault.New(),
		}, catalog.Annotators()...),
	})
	assert.NoError(t, f.Annotate(t, s).Err, "the annotate half of suite runs whole")
	assert.NoError(t, f.Generate(t, s).Err, "the generate half of suite runs whole")
	res := f.Generate(t, p)
	assert.NoError(t, res.Err, "fake runs whole")
	return run{f: f, suite: s, fake: p, res: res}
}

// models returns the model of each fake of the run by the name of its
// declaration. It checks that the family claims its files through its
// template.
func (r run) models(t *testing.T) map[string]gentest.Record[fake.Model] {
	t.Helper()

	return gentest.Records(t, r.res, fakeTemplate, "", func(m fake.Model) string { return m.Subject })
}

// tests returns the record of the test of each fake of the run by the name of
// its declaration. It checks that the family of the tests claims its files
// through its template.
func (r run) tests(t *testing.T) map[string]gentest.Record[fake.Test] {
	t.Helper()

	return gentest.Records(t, r.res, testTemplate, fake.TestTag, func(c fake.Test) string { return c.Subject })
}
