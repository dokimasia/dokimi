// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sim_test

import (
	"testing"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"

	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/generator/sim"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The tree of the specs, the packages that two specs or more read, and the
// templates of the two families.
const (
	goTree             = "testdata/golang"
	registerPackage    = "example.com/shop/register"
	registerSimPackage = "example.com/shop/register/registersim"
	plainPackage       = "example.com/shop/plain"
	disksPackage       = "example.com/shop/disks"
	disksSimPackage    = "example.com/shop/disks/diskssim"
	plainSimPackage    = "example.com/shop/plain/plainsim"
	genericPackage     = "example.com/shop/generic"
	genericSimPackage  = "example.com/shop/generic/genericsim"
	remotePackage      = "example.com/shop/remote"
	remoteSimPackage   = "example.com/shop/remote/remotesim"
	storePackage       = "example.com/shop/remote/store"
	cellsPackage       = "example.com/shop/cells"
	cellsSimPackage    = "example.com/shop/cells/cellssim"
	tasksSimPackage    = "example.com/shop/cells/taskssim"
	placedPackage      = "example.com/shop/placed"
	placedSimPackage   = "example.com/shop/placed/placedsim"
	waitsPackage       = "example.com/shop/waits"
	waitsSimPackage    = "example.com/shop/waits/waitssim"
	hostedPackage      = "example.com/shop/hosted"
	hostedSimPackage   = "example.com/shop/hosted/hostedsim"
	simTemplate        = "sim.tmpl"
	testTemplate       = "sim.test.tmpl"
)

// run is one run of both halves of sim over a fixture: the plugin, the
// fixture, and the result of each half.
type run struct {
	p      plugin.Plugin
	f      *gentest.Fixture
	stamps plugintest.Result
	res    plugintest.Result
}

// fixture loads the packages pkgs of the tree of the specs for the plugin p,
// with the target golang and the Go backend as its spoke, so a phase call
// previews the package of each routed file.
func fixture(tb assert.TB, p plugin.Plugin, pkgs ...string) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{
		Dir: goTree, Packages: pkgs, Target: golang.Target, Types: gobackend.New(), Plugins: []plugin.Plugin{p},
	})
}

// generate runs the annotate half and then the generate half of a new sim
// over the packages pkgs of the tree of the specs.
func generate(t *testing.T, pkgs ...string) run {
	t.Helper()

	p := sim.New()
	f := fixture(t, p, pkgs...)
	stamps := f.Annotate(t, p)
	assert.NoError(t, stamps.Err, "the annotate half runs whole")
	res := f.Generate(t, p)
	assert.NoError(t, res.Err, "the generate half runs whole")
	return run{p: p, f: f, stamps: stamps, res: res}
}

// generateWithFaults runs the annotator fault, then the annotate half and the
// generate half of a new sim, over the packages pkgs of the tree of the
// specs, as a composition runs fault before sim.
func generateWithFaults(t *testing.T, pkgs ...string) run {
	t.Helper()

	faults, p := fault.New(), sim.New()
	f := gentest.New(t, &gentest.Spec{
		Dir: goTree, Packages: pkgs, Target: golang.Target, Types: gobackend.New(),
		Plugins: []plugin.Plugin{faults, p},
	})
	assert.NoError(t, f.Annotate(t, faults).Err, "the annotator fault runs whole")
	stamps := f.Annotate(t, p)
	assert.NoError(t, stamps.Err, "the annotate half runs whole")
	res := f.Generate(t, p)
	assert.NoError(t, res.Err, "the generate half runs whole")
	return run{p: p, f: f, stamps: stamps, res: res}
}

// models returns the model of each simulation of the run by the name of its
// machine interface. It checks that the family claims its files through its
// template.
func (r run) models(t *testing.T) map[string]sim.Model {
	t.Helper()

	out := map[string]sim.Model{}
	for name, rec := range gentest.Records(t, r.res, simTemplate, "", func(m sim.Model) string { return m.Machine }) {
		out[name] = rec.Data
	}
	return out
}

// checks returns the record of the tests of each simulation of the run by
// the name of its machine interface. It checks that the family of the tests
// claims its files through its template.
func (r run) checks(t *testing.T) map[string]gentest.Record[sim.Checks] {
	t.Helper()

	return gentest.Records(t, r.res, testTemplate, sim.TestTag, func(c sim.Checks) string { return c.Machine })
}
