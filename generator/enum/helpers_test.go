// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees of the specs, the packages and modules that two specs or more
// read, and the templates of the two families.
const (
	goTree             = "testdata/golang"
	tsTree             = "testdata/typescript"
	statusPackage      = "example.com/shop/status"
	regionPackage      = "example.com/shop/region"
	quietPackage       = "example.com/shop/quiet"
	commentedPackage   = "example.com/shop/commented"
	handwrittenPackage = "example.com/shop/handwritten"
	statusModule       = "status"
	enumTemplate       = "enum.tmpl"
	testTemplate       = "enum.test.tmpl"
)

// The enumerations that two specs or more read.
const (
	statusName  = "Status"
	regionName  = "Region"
	levelName   = "Level"
	weekdayName = "Weekday"
)

// generated is what enum records over one fixture: the model of the surface,
// the names of the declarations that its record declares beside the struct of
// the enumeration, the checks and the record of the checks of each
// enumeration by its name, and the findings of the run.
type generated struct {
	models   map[string]enum.Model
	beside   map[string][]string
	checks   map[string]enum.Checks
	tests    map[string]*emit.Struct
	findings []diag.Diag
}

// fixture returns a fixture over the packages pkgs of the tree dir, loaded
// under the context of tb, with the target of the plan and the plugin enum
// p.
func fixture(tb assert.TB, target plugin.Target, dir string, p plugin.Plugin, pkgs ...string) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{Dir: dir, Packages: pkgs, Target: target, Plugins: []plugin.Plugin{p}})
}

// generate runs a new enum over the packages pkgs of the tree dir with the
// target of the plan, and returns what it records.
func generate(t *testing.T, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	return run(t, enum.New(), target, dir, pkgs...)
}

// run runs the enum p over the packages pkgs of the tree dir with the target
// of the plan, and returns what it records. It checks that each family
// claims its files through its template.
func run(t *testing.T, p plugin.Generator, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	res := fixture(t, target, dir, p, pkgs...).Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	got := generated{
		models:   map[string]enum.Model{},
		checks:   map[string]enum.Checks{},
		tests:    map[string]*emit.Struct{},
		findings: slices.Collect(res.Sink.All()),
	}
	subject := func(m enum.Model) string { return m.Subject }
	for name, r := range gentest.Records(t, res, enumTemplate, "", subject) {
		got.models[name] = r.Data
	}
	got.beside = gentest.Beside(res, subject)
	for name, r := range gentest.Records(t, res, testTemplate, enum.TestTag,
		func(c enum.Checks) string { return c.Subject }) {
		got.checks[name] = r.Data
		got.tests[name] = r.Decl
	}
	return got
}

// texts returns the names and the texts of the variants.
func texts(vs []enum.Variant) []string {
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		out = append(out, v.Name+"="+v.Text)
	}
	return out
}
