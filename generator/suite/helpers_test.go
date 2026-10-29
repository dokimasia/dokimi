// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	goannotate "go.dokimi.dev/eidos/lang/go/annotate"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	"go.dokimi.dev/eidos/plugin/shape/catalog"
	"go.dokimi.dev/eidos/sdk/authored"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/generator/builder"
	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The trees of the specs, the packages that two specs or more read, and the
// templates of the two families.
const (
	goTree             = "testdata/golang"
	storePackage       = "example.com/shop/store"
	genericPackage     = "example.com/shop/generic"
	catalogPackage     = "example.com/shop/catalog"
	slugPackage        = "example.com/shop/slug"
	memoryPackage      = "example.com/shop/memory"
	cachedPackage      = "example.com/shop/cached"
	sessionPackage     = "example.com/shop/session"
	meterPackage       = "example.com/shop/meter"
	meteredPackage     = "example.com/shop/metered"
	callsPackage       = "example.com/shop/calls"
	inputsPackage      = "example.com/shop/inputs"
	pathsPackage       = "example.com/shop/paths"
	speccedPackage     = "example.com/shop/specced"
	speccedTestPackage = "example.com/shop/specced_test"
	badBudgetPackage   = "example.com/shop/badbudget"
	dualPackage        = "example.com/shop/dual"
	ledgerPackage      = "example.com/shop/ledger"
	loosePackage       = "example.com/shop/loose"
	nestedPackage      = "example.com/shop/nested"
	wrapPackage        = "example.com/shop/wrap"
	kindsPackage       = "example.com/shop/kinds"
	durablePackage     = "example.com/shop/durable"
	taggedPackage      = "example.com/shop/tagged"
	locksPackage       = "example.com/shop/locks"
	toolPackage        = "example.com/shop/lib/tool"
	deepPackage        = "example.com/shop/lib/internal/deep"
	derivedPackage     = "example.com/shop/derived"
	suiteTemplate      = "suite.tmpl"
	conformTemplate    = "conform.test.tmpl"
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

// fixture returns a fixture over the packages pkgs of the tree dir, as
// [specOf] describes it.
func fixture(tb assert.TB, dir string, p plugin.Plugin, pkgs ...string) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, specOf(dir, p, pkgs...))
}

// specOf returns the spec of a fixture over the packages pkgs of the tree dir
// with the plugin suite p, the target golang and the Go backend as its spoke,
// so a phase call previews the package of each routed file. The annotators of
// the composition stamp the fixture: the authored samples and witnesses, the
// facts of Go, the declared defaults, the names of the builders, the faults of
// the methods and the classifications of the shape catalog.
func specOf(dir string, p plugin.Plugin, pkgs ...string) *gentest.Spec {
	annotators := append([]plugin.Annotator{
		authored.Sample(), authored.Witness(), goannotate.New(), defaults.New(), builder.New(), fault.New(),
	}, catalog.Annotators()...)
	return &gentest.Spec{
		Dir:        dir,
		Packages:   pkgs,
		Target:     golang.Target,
		Types:      gobackend.New(),
		Plugins:    []plugin.Plugin{p},
		Annotators: annotators,
	}
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

// called returns the call of the callable whose field is field in the suite
// s, and the call of a function for an empty field. It stops the test where
// the suite has no such call.
func called(t *testing.T, s suite.Suite, field string) suite.Call {
	t.Helper()

	at := slices.IndexFunc(s.Calls, func(c suite.Call) bool { return c.Field == field })
	assert.True(t, at >= 0, "the suite of "+s.Declaration+" has a call of "+field)
	return s.Calls[at]
}

// signature returns the signature of the method m: its name, then the name
// and the spelling of the type of each parameter, and the name and the
// spelling of the type of each result. A parameter or a result without a
// name has the spelling alone.
func signature(m *emit.Method) string {
	params := make([]string, 0, m.Params.Len())
	for _, p := range m.Params.Items() {
		spelling := spelled(p.Type)
		if p.Variadic != symbol.VariadicNone {
			spelling = "..." + spelling
		}
		params = append(params, strings.TrimSpace(p.Name+" "+spelling))
	}
	results := make([]string, 0, m.Returns.Len())
	for _, r := range m.Returns.Items() {
		results = append(results, strings.TrimSpace(r.Name+" "+spelled(r.Type)))
	}
	return m.Name + "(" + strings.Join(params, ", ") + ") (" + strings.Join(results, ", ") + ")"
}

// spelled returns the spelling of the type ref and of its type arguments, as
// Pair[time.Month].
func spelled(ref *emit.TypeRef) string {
	if len(ref.Args) == 0 {
		return ref.Spelling
	}
	args := make([]string, 0, len(ref.Args))
	for _, a := range ref.Args {
		args = append(args, spelled(a))
	}
	return ref.Spelling + "[" + strings.Join(args, ", ") + "]"
}

// fn is a generated function with the signature of its literal and the
// spellings of its results.
type fn struct {
	field, signature, index string
	results                 []result
}

// result is the list of one result of a generated function, with the
// spelling of its type, the names of its declared errors and the variable
// of the end value of a channel.
type result struct {
	name, label, spelling, end string
	errors                     []string
	isError, sink              bool
}

// fnsOf returns each function of funcs with the signature of its literal and
// the spellings of its results.
func fnsOf(funcs []suite.Func) []fn {
	out := make([]fn, 0, len(funcs))
	for _, f := range funcs {
		got := fn{field: f.Field, signature: signature(f.Signature), index: f.Index}
		for _, r := range f.Results {
			res := result{
				name: r.Name, label: r.Label, spelling: spelled(r.Type), end: r.End, isError: r.Error, sink: r.Sink,
			}
			for _, e := range r.Errors {
				res.errors = append(res.errors, e.Callee.Name)
			}
			got.results = append(got.results, res)
		}
		out = append(out, got)
	}
	return out
}

// ungenerated returns the message of each finding Ungenerated of got about a
// callable of the declaration named decl, in the order of the findings. The
// message of the finding opens with the callable.
func ungenerated(got generated, decl string) []string {
	var out []string
	for _, d := range got.findings {
		if d.Code == suite.Ungenerated && (strings.HasPrefix(d.Msg, decl+".") || strings.HasPrefix(d.Msg, decl+" ")) {
			out = append(out, d.Msg)
		}
	}
	return out
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
