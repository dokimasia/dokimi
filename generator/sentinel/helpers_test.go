// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel_test

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/sentinel"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees of the specs, the packages and modules that two specs or more
// read, and the template of the family.
const (
	goTree           = "testdata/golang"
	tsTree           = "testdata/typescript"
	shopModule       = "example.com/shop"
	storePackage     = "example.com/shop/store"
	neighbourPackage = "example.com/shop/neighbour"
	embeddedPackage  = "example.com/shop/embedded"
	noprefixPackage  = "example.com/shop/noprefix"
	splitPackage     = "example.com/shop/split"
	typedPackage     = "example.com/shop/typed"
	storeModule      = "store"
	testTemplate     = "sentinel.test.tmpl"
)

// generated is what sentinel records over one fixture: the checks of each
// package by its name, and the findings of the run.
type generated struct {
	checks   map[string]sentinel.Checks
	findings []diag.Diag
}

// fixture returns a fixture over the packages pkgs of the tree dir, loaded
// under the context of tb, with the target of the plan and the plugin
// sentinel p.
func fixture(tb assert.TB, target plugin.Target, dir string, p plugin.Plugin, pkgs ...string) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{Dir: dir, Packages: pkgs, Target: target, Plugins: []plugin.Plugin{p}})
}

// generate runs a new sentinel over the packages pkgs of the tree dir with
// the target of the plan, and returns what it records.
func generate(t *testing.T, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	return run(t, sentinel.New(), target, dir, pkgs...)
}

// run runs the sentinel p over the packages pkgs of the tree dir with the
// target of the plan, and returns what it records. It checks that the family
// claims its files through its template.
func run(t *testing.T, p plugin.Generator, target plugin.Target, dir string, pkgs ...string) generated {
	t.Helper()

	res := fixture(t, target, dir, p, pkgs...).Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	got := generated{checks: map[string]sentinel.Checks{}, findings: slices.Collect(res.Sink.All())}
	for name, r := range gentest.Records(t, res, testTemplate, sentinel.TestTag,
		func(c sentinel.Checks) string { return c.Subject }) {
		got.checks[name] = r.Data
	}
	return got
}

// statements returns the statements of a method of a record as text: the
// names that each statement assigns, an equals sign and the value, as in
// got = new NotFoundError("key-a").
func statements(stmts []emit.Stmt) string {
	out := make([]string, 0, len(stmts))
	for i := range stmts {
		out = append(out, strings.Join(stmts[i].Names, ", ")+" = "+spelled(stmts[i].Value.Val))
	}
	return strings.Join(out, "; ")
}

// spelled returns a value of a statement in the syntax of its language: a
// construction, an address, a list, a composite with named fields, or the
// text of a literal.
func spelled(v *emit.Value) string {
	parts := make([]string, 0, len(v.Args)+len(v.Fields))
	for i := range v.Args {
		parts = append(parts, spelled(&v.Args[i]))
	}
	for i := range v.Fields {
		if name := v.Fields[i].Name; name != "" {
			parts = append(parts, name+": "+spelled(&v.Fields[i].Value))
			continue
		}
		parts = append(parts, spelled(&v.Fields[i].Value))
	}
	switch {
	case v.Kind == emit.ValueConstruction:
		return "new " + v.Type.Spelling + "(" + strings.Join(parts, ", ") + ")"
	case v.Kind == emit.ValueAddress:
		return "&" + spelled(v.Inner)
	case v.Kind == emit.ValueComposite && v.Type.Form == symbol.FormList:
		return "[" + strings.Join(parts, ", ") + "]"
	case v.Kind == emit.ValueComposite:
		return v.Type.Spelling + "{" + strings.Join(parts, ", ") + "}"
	case v.Literal == emit.LiteralString:
		return strconv.Quote(v.Text)
	default:
		return v.Text
	}
}
