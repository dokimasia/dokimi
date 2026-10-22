// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	goannotate "go.dokimi.dev/eidos/lang/go/annotate"
	"go.dokimi.dev/eidos/plugin/shape/catalog"
	"go.dokimi.dev/eidos/sdk/authored"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/generator/stub"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees of the specs, the packages and modules that two specs or more
// read, the templates of the two families, and the declarations under the
// directive that two specs or more read.
const (
	goTree          = "testdata/golang"
	tsTree          = "testdata/typescript"
	catalogPackage  = "example.com/shop/catalog"
	parityPackage   = "example.com/shop/parity"
	genericPackage  = "example.com/shop/generic"
	edgePackage     = "example.com/shop/edge"
	handlesPackage  = "example.com/shop/handles"
	mixedPackage    = "example.com/shop/mixed"
	sequencePackage = "example.com/shop/sequence"
	builtPackage    = "example.com/shop/built"
	billingPackage  = "example.com/shop/billing"
	storeModule     = "store"
	clockModule     = "clock"
	handlesModule   = "handles"
	streamModule    = "stream"
	mixedModule     = "mixed"
	pointModule     = "point"
	builtModule     = "built"
	billingModule   = "billing"
	stubTemplate    = "stub.tmpl"
	testTemplate    = "stub.test.tmpl"
	storeName       = "Store"
	clockName       = "Clock"
	repoName        = "Repo"
	edgeName        = "Edge"
)

// generated is what stub records over one fixture: the model, the record of
// each double, the names of the declarations that the record declares after
// the double and the record of its tests, by the name of the subject, and the
// findings of the run.
type generated struct {
	models   map[string]stub.Model
	doubles  map[string]*emit.Struct
	beside   map[string][]string
	tests    map[string]*emit.Struct
	findings []diag.Diag
}

// fixture returns a fixture over the packages pkgs of the tree dir, loaded
// under the context of tb, with the target of the plan and the plugin stub
// p. The annotators of the composition stamp the fixture.
func fixture(tb assert.TB, target plugin.Target, dir string, p plugin.Plugin, pkgs ...string) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{
		Dir:        dir,
		Packages:   pkgs,
		Target:     target,
		Plugins:    []plugin.Plugin{p},
		Annotators: annotators(),
	})
}

// annotators returns the annotators of the composition: the authored samples
// and witnesses, the facts of Go, the declared defaults, the names of the
// builders, the faults of the methods and the classifications of the shape
// catalog.
func annotators() []plugin.Annotator {
	return append([]plugin.Annotator{
		authored.Sample(), authored.Witness(), goannotate.New(), defaults.New(), builder.New(), fault.New(),
	}, catalog.Annotators()...)
}

// generate runs a new stub over the packages pkgs of the tree dir with the
// target of the plan, and returns what it records.
func generate(t *testing.T, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	return run(t, stub.New(), target, dir, pkgs...)
}

// run runs the plugin p over the packages pkgs of the tree dir with the
// target of the plan, and returns what it records. It checks that each family
// claims its files through its template.
func run(t *testing.T, p plugin.Generator, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	res := fixture(t, target, dir, p, pkgs...).Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	got := generated{
		models:   map[string]stub.Model{},
		doubles:  map[string]*emit.Struct{},
		tests:    map[string]*emit.Struct{},
		findings: slices.Collect(res.Sink.All()),
	}
	subject := func(m stub.Model) string { return m.Subject }
	for name, r := range gentest.Records(t, res, stubTemplate, "", subject) {
		got.models[name] = r.Data
		got.doubles[name] = r.Decl
	}
	got.beside = gentest.Beside(res, subject)
	for name, r := range gentest.Records(t, res, testTemplate, stub.TestTag, func(c stub.Checks) string { return c.Subject }) {
		got.tests[name] = r.Decl
	}
	return got
}

// method returns the method named name of the model m. It stops the test
// where m has no such method.
func method(t *testing.T, m stub.Model, name string) stub.Method {
	t.Helper()

	at := slices.IndexFunc(m.Methods, func(c stub.Method) bool { return c.Name == name })
	assert.NotEqual(t, at, -1, "the model of "+m.Subject+" has the method "+name)
	return m.Methods[at]
}

// declared returns the method named name of the record r. It stops the test
// where r has no such method.
func declared(t *testing.T, r *emit.Struct, name string) *emit.Method {
	t.Helper()

	assert.NotNil(t, r, "stub records the record")
	at := slices.IndexFunc(r.Methods.Items(), func(m *emit.Method) bool { return m.Name == name })
	assert.NotEqual(t, at, -1, "the record "+r.Name+" has the method "+name)
	return r.Methods.Items()[at]
}

// paramFields returns the field of each parameter, joined by spaces.
func paramFields(ps []stub.Param) string {
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Field)
	}
	return strings.Join(out, " ")
}

// resultFields returns the field of each result, joined by spaces.
func resultFields(rs []stub.Result) string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.Field)
	}
	return strings.Join(out, " ")
}
