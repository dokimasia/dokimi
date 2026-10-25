// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	gofrontend "go.dokimi.dev/eidos/lang/go/frontend"
	gorules "go.dokimi.dev/eidos/lang/go/rules"
	tsfrontend "go.dokimi.dev/eidos/lang/typescript/frontend"
	tsrules "go.dokimi.dev/eidos/lang/typescript/rules"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/frontendtest"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/output"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/rulestest"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// Brand is the brand of the trees of the specs, so a carrier sets a
// directive as //dokimi:builder.
const Brand output.Brand = "dokimi"

// Spec describes a fixture: the environment of its stores, the tree, the
// packages of the fixture, the target of the plan, and the plugins that the
// fixture registers and runs.
type Spec struct {
	// Getenv reads the environment through which the frontends locate their
	// stores. A nil Getenv reads the environment of the process, as the
	// command line of eidos does.
	Getenv func(key string) string
	// Dir is the directory of the tree, slash-separated and relative to the
	// directory of the spec.
	Dir string
	// Packages are the paths of the packages of the tree that the fixture
	// contains: the import path of a Go package, such as
	// example.com/shop/catalog, the path of a TypeScript module relative to
	// the tree, such as catalog, or the package line of a scripted file, such
	// as shop.
	Packages []string
	// Target is the target of the plan of a generator's phase call. It is
	// zero for an annotator.
	Target plugin.Target
	// Types is the spoke of the backend of the plan of a generator's phase
	// call, such as the Go backend. A phase call translates types through it,
	// and previews the package of a routed file where the spoke spells file
	// names and derives packages. It is nil for an annotator.
	Types plugin.TypeSpeller
	// Plugins are the plugins whose keys the fixture registers and whose
	// directives it validates, the plugin under test among them.
	Plugins []plugin.Plugin
	// Annotators stamp the fixture in order before the plugin under test
	// runs, as the annotators of earlier turns stamp a run.
	Annotators []plugin.Annotator
}

// Fixture is a plugintest fixture over the packages of a tree, with the
// declarations of those packages.
//
// # Concurrency
//
// A Fixture belongs to one test, as a plugintest fixture does.
type Fixture struct {
	*plugintest.Fixture
	// decls are the declarations of the packages of the fixture, in the
	// order of the walk of each package.
	decls []node.Declaration
}

// New returns a fixture over the packages of the tree that s describes. It
// loads the tree through [plugintest.Loaded] under [Brand], with the stores
// that the frontends locate through [Spec.Getenv], as a run of the
// composition loads it. So a global type of TypeScript resolves to the
// library of TypeScript, and an import of the standard library of Go resolves
// to its package. The fixture has the frontends and the rules of Go,
// TypeScript and the scripted language of eidos. No plugin of dokimi serves
// the scripted language, so a spec loads a scripted tree to show that a
// plugin leaves the declarations of another language alone. New registers the
// plugins and the annotators of s. It seals the graph to list the
// declarations of the packages, as the first phase call seals it. The
// annotators of s then stamp the fixture in order. New stops the test where a
// frontend locates no store, where [plugintest.Loaded] stops it and where an
// annotator returns an error. It fails the test for each Error that an
// annotator reports.
//
// Each call loads the tree again, so a plugintest setup that calls New
// returns a fresh fixture on every call.
func New(tb assert.TB, s *Spec) *Fixture {
	tb.Helper()

	getenv := s.Getenv
	if getenv == nil {
		getenv = os.Getenv
	}
	frontends := []plugin.Frontend{gofrontend.New(nil), tsfrontend.New(), frontendtest.NewScripted()}
	f := &Fixture{Fixture: plugintest.Loaded(tb, plugintest.Source{
		Tree:      os.DirFS(filepath.FromSlash(s.Dir)),
		Stores:    stores(tb, frontends, getenv),
		Brand:     Brand,
		Frontends: frontends,
		Rules:     []rules.SourceRules{gorules.New(), tsrules.New(), rulestest.Scripted()},
		Plugins:   slices.Concat(s.Plugins, annotators(s.Annotators)),
		Packages:  s.Packages,
	})}
	f.Target, f.Types = s.Target, s.Types
	f.Graph.Freeze()
	for pkg := range f.Graph.Packages() {
		f.decls = slices.AppendSeq(f.decls, node.Declarations(pkg))
	}
	for _, a := range s.Annotators {
		res := f.Annotate(tb, a)
		assert.NoError(tb, res.Err, "the annotator "+string(a.Name())+" annotates")
		for d := range res.Sink.All() {
			expect.NotEqual(tb, d.Severity, diag.SeverityError, fmt.Sprintf(
				"the annotator %s reports no Error, and reports %s at %s: %s", a.Name(), d.Code, d.Pos, d.Msg))
		}
	}
	return f
}

// Find returns the declaration of a kind in the packages of the fixture,
// with an owner and a name. owner is the name of the declaration that the
// declaration is a member of, and empty for a declaration at file level.
// Find stops the test where no declaration or more than one matches.
func (f *Fixture) Find(tb assert.TB, kind symbol.Kind, owner, name string) node.Declaration {
	tb.Helper()

	var found []node.Declaration
	for _, d := range f.decls {
		if id := d.Identity(); id.Kind == kind && id.Owner == owner && id.Name == name {
			found = append(found, d)
		}
	}
	assert.Length(tb, found, 1, fmt.Sprintf("the fixture declares one %s %q in %q", kind, name, owner))
	return found[0]
}

// stores returns the stores that the frontends locate through getenv, by
// name, as the command line of eidos locates them. It stops the test where a
// frontend locates no store.
func stores(tb assert.TB, frontends []plugin.Frontend, getenv func(key string) string) map[string]fs.FS {
	tb.Helper()

	out := map[string]fs.FS{}
	for _, front := range frontends {
		if locator, locates := front.(plugin.StoreLocator); locates {
			located, err := locator.Stores(getenv)
			assert.NoError(tb, err, "the frontend "+string(front.Name())+" locates its stores")
			maps.Copy(out, located)
		}
	}
	return out
}

// annotators returns the annotators as plugins, in order.
func annotators(as []plugin.Annotator) []plugin.Plugin {
	out := make([]plugin.Plugin, 0, len(as))
	for _, a := range as {
		out = append(out, a)
	}
	return out
}
