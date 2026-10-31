// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	goannotate "go.dokimi.dev/eidos/lang/go/annotate"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	"go.dokimi.dev/eidos/sdk/authored"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees of the specs, the packages and modules that two specs or more
// read, and the templates of the two families.
const (
	goTree            = "testdata/golang"
	tsTree            = "testdata/typescript"
	catalogPackage    = "example.com/shop/catalog"
	seededPackage     = "example.com/shop/seeded"
	genericPackage    = "example.com/shop/generic"
	companionsPackage = "example.com/shop/companions"
	hiddenPackage     = "example.com/shop/hidden"
	catalogModule     = "catalog"
	aliasedModule     = "aliased"
	builderTemplate   = "builder.tmpl"
	testTemplate      = "builder.test.tmpl"
)

// The declarations under the directive that two specs or more read.
const (
	addressName = "Address"
	orderName   = "Order"
	accountName = "Account"
)

// generated is what builder records over one fixture: the model, the start
// body, the names of the functions that the record declares beside the
// builder, the checks and the record of the checks of each builder by the
// name of its subject, and the findings of both halves.
type generated struct {
	models    map[string]builder.Model
	starts    map[string]*emit.Method
	functions map[string][]string
	checks    map[string]builder.Checks
	tests     map[string]*emit.Struct
	findings  []diag.Diag
}

// fixture returns a fixture over the packages pkgs of the tree dir, loaded
// under the context of tb, with the target and the spoke of the plan and the
// plugin builder p. The annotators of the composition stamp the fixture: the
// authored samples and witnesses, the facts of Go and the defaults.
func fixture(tb assert.TB, target plugin.Target, dir string, p plugin.Plugin, pkgs ...string) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{
		Dir:        dir,
		Packages:   pkgs,
		Target:     target,
		Types:      spokes[target],
		Plugins:    []plugin.Plugin{p},
		Annotators: []plugin.Annotator{authored.Sample(), authored.Witness(), goannotate.New(), defaults.New()},
	})
}

// spokes are the spokes of the fixtures by their target. The Go backend
// previews the package of each routed file, as the plan of the composition
// previews it.
var spokes = map[plugin.Target]plugin.TypeSpeller{golang.Target: gobackend.New()}

// generate runs a new builder over the packages pkgs of the tree dir with
// the target of the plan, and returns what it records.
func generate(t *testing.T, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	return run(t, builder.New(), target, dir, pkgs...)
}

// phases runs the annotate half and then the generate half of the builder p
// over the fixture f, as a run of the composition runs them. It returns the
// result of each half, and stops the test where a half returns an error.
func phases(t *testing.T, f *gentest.Fixture, p plugin.Plugin) (stamps, records plugintest.Result) {
	t.Helper()

	stamps = f.Annotate(t, p)
	assert.NoError(t, stamps.Err, "the annotate half runs whole")
	records = f.Generate(t, p)
	assert.NoError(t, records.Err, "the generate half runs whole")
	return stamps, records
}

// run runs the builder p over the packages pkgs of the tree dir with the
// target of the plan, and returns what it records and the findings of both
// halves. It checks that each family claims its files through its template.
func run(t *testing.T, p plugin.Generator, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	stamps, res := phases(t, fixture(t, target, dir, p, pkgs...), p)
	got := generated{
		models:   map[string]builder.Model{},
		starts:   map[string]*emit.Method{},
		checks:   map[string]builder.Checks{},
		tests:    map[string]*emit.Struct{},
		findings: slices.Concat(slices.Collect(stamps.Sink.All()), slices.Collect(res.Sink.All())),
	}
	subject := func(m builder.Model) string { return m.Subject }
	for name, r := range gentest.Records(t, res, builderTemplate, "", subject) {
		got.models[name] = r.Data
		got.starts[name] = r.Decl.Types.Items()[0].(*emit.Method)
	}
	got.functions = gentest.Beside(res, subject)
	checks := gentest.Records(t, res, testTemplate, builder.TestTag, func(c builder.Checks) string { return c.Subject })
	for name, r := range checks {
		got.checks[name] = r.Data
		got.tests[name] = r.Decl
	}
	return got
}

// member returns the member named name of the model m. It stops the test
// where m has no such member.
func member(t *testing.T, m builder.Model, name string) builder.Member {
	t.Helper()

	at := slices.IndexFunc(m.Members, func(mem builder.Member) bool { return mem.Field == name })
	assert.NotEqual(t, at, -1, "the model of "+m.Subject+" has the member "+name)
	return m.Members[at]
}

// fields returns the fields of the members, joined by spaces.
func fields(ms []builder.Member) string {
	out := make([]string, 0, len(ms))
	for _, m := range ms {
		out = append(out, m.Field)
	}
	return strings.Join(out, " ")
}

// seedOf returns the name of the function whose result the first statement
// of the start body declares as the value. It stops the test where the
// start body declares no value.
func seedOf(t *testing.T, start *emit.Method) string {
	t.Helper()

	stmts := start.Body.Stmts
	assert.NotEmpty(t, stmts, "the start body declares the value")
	assert.True(t, stmts[0].Declare, "the first statement declares the value")
	return stmts[0].Value.Fn.Ref.Name
}
