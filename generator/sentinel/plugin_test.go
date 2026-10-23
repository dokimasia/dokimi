// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/sentinel"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The packages that sentinel reports.
const (
	emptyPackage    = "example.com/shop/bad/empty"
	selfPackage     = "example.com/shop/bad/self"
	distantPackage  = "example.com/shop/bad/distant"
	strangerPackage = "example.com/shop/bad/stranger"
	aliasPackage    = "example.com/shop/alias"
)

// sentinel records the checks of the error values and of the error types of
// each package under its directive, and passes the plugin suite.
func TestSentinel(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		suites := []struct {
			name   string
			target plugin.Target
			dir    string
			pkgs   []string
		}{
			{
				name: "passes the plugin suite over two Go packages", target: golang.Target,
				dir: goTree, pkgs: []string{storePackage, neighbourPackage},
			},
			{
				name: "passes the plugin suite over the TypeScript module store", target: typescript.Target,
				dir: tsTree, pkgs: []string{storeModule},
			},
		}
		for _, tt := range suites {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
					p := sentinel.New()
					return p, fixture(tb, tt.target, tt.dir, p, tt.pkgs...).Fixture
				})
			})
		}

		t.Run("checks each package from outside", func(t *testing.T) {
			t.Parallel()

			gentest.External(t, sentinel.New(), sentinel.TestTag)
		})

		sentinels := []struct {
			name string
			pkg  string
			give string
			want string
		}{
			{
				name: "records each variable whose name the rules accept for an error value", pkg: storePackage,
				give: "store", want: "ErrNotFound ErrConflict ErrClosed",
			},
			{
				name: "records the error values of each file of the package", pkg: splitPackage, give: "split",
				want: "ErrRead ErrWrite",
			},
			{name: "records no error value of a package of error types", pkg: typedPackage, give: "typed"},
		}
		for _, tt := range sentinels {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, recorded := generate(t, golang.Target, goTree, tt.pkg).checks[tt.give]
				assert.True(t, recorded, "sentinel records the checks of "+tt.give)
				expect.Equal(t, names(got.Sentinels), tt.want, "the error values are in the order of the source")
			})
		}

		kinds := []struct {
			name string
			pkg  string
			give string
			want []string
		}{
			{
				name: "records what a check reads of each public error type",
				pkg:  storePackage, give: "store",
				want: []string{
					"NotFoundError pointer compares members=Key:text,Attempt",
					"WrappedError pointer unwraps cause=Cause members=Op:text",
				},
			},
			{
				name: "records what an embed by value promotes", pkg: embeddedPackage,
				give: "embedded",
				want: []string{
					"BaseError compares unwraps cause=Cause members=Op:text",
					"LostError pointer compares unwraps",
					"NotFoundError compares unwraps cause=Cause",
				},
			},
			{
				name: "records an error type with a value receiver", pkg: typedPackage, give: "typed",
				want: []string{"LimitError members=Name:text"},
			},
		}
		for _, tt := range kinds {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, golang.Target, goTree, tt.pkg).checks[tt.give]
				expect.Equal(t, describe(got.Kinds), tt.want, "each error type has the facts of its contract")
			})
		}

		prefixes := []struct {
			name   string
			pkg    string
			give   string
			option string
			want   string
		}{
			{
				name: "opens each message with the prefix of the directive", pkg: storePackage, give: "store",
				want: "kv: ",
			},
			{
				name: "opens each message with the name of the package by default", pkg: neighbourPackage,
				give: "neighbour", want: "neighbour: ",
			},
			{name: "leaves out the opening where the prefix is off", pkg: noprefixPackage, give: "noprefix"},
			{
				name: "opens each message with the option prefix where the directive states none",
				pkg:  neighbourPackage, give: "neighbour", option: "acme", want: "acme: ",
			},
			{
				name: "opens each message with the prefix of the directive over the option", pkg: storePackage,
				give: "store", option: "acme", want: "kv: ",
			},
		}
		for _, tt := range prefixes {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := sentinel.New()
				p.(plugin.OptionsProvider).Options().(*sentinel.Options).Prefix = tt.option
				got := run(t, p, golang.Target, goTree, tt.pkg).checks[tt.give]
				expect.Equal(t, got.Prefix, tt.want, "the prefix opens each message")
			})
		}

		t.Run("records the error values of each package of the param no-overlap-with", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, storePackage, neighbourPackage).checks["store"]
			assert.Length(t, got.Neighbours, 1, "the store lists one neighbour")
			expect.Equal(t, got.Neighbours[0].Name, "neighbour", "the neighbour has the name of its package")
			expect.Equal(t, names(got.Neighbours[0].Sentinels), "ErrNotFound ErrConflict",
				"the neighbour has its error values")
		})

		t.Run("records a re-export of an error value of another package", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, storePackage, relayPackage).checks["relay"]
			reexports := map[string]bool{}
			for _, s := range got.Sentinels {
				reexports[s.Name] = s.Reexport
			}
			expect.Equal(t, reexports, map[string]bool{"ErrNotFound": true, "ErrDropped": false},
				"ErrNotFound is the error value of store, and ErrDropped is a value of relay")
		})

		t.Run("records no re-export of an error value of its own package", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, aliasPackage).checks["alias"]
			reexports := map[string]bool{}
			for _, s := range got.Sentinels {
				reexports[s.Name] = s.Reexport
			}
			expect.Equal(t, reexports, map[string]bool{"ErrBase": false, "ErrSame": false},
				"ErrSame is ErrBase of the same package under a second name")
		})

		t.Run("records each error type under the param remote as remote", func(t *testing.T) {
			t.Parallel()

			res := generate(t, golang.Target, goTree, storePackage, relayPackage)
			remote := map[string]bool{}
			for _, k := range res.checks["relay"].Kinds {
				remote[k.Name] = k.Remote
			}
			expect.Equal(t, remote, map[string]bool{"CallError": true, "TimeoutError": false},
				"CallError carries the message of the remote store")
			unknown := slices.ContainsFunc(res.findings, func(d diag.Diag) bool {
				return d.Code == sentinel.UnknownRemote
			})
			expect.False(t, unknown, "CallError is an error type of relay, so sentinel reports no UnknownRemote")
		})

		t.Run("keeps a package of the param that the run did not load without error values", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, distantPackage).checks["distant"]
			const missing = "example.com/shop/missing"
			expect.Equal(t, got.Neighbours, []sentinel.Neighbour{{Path: missing, Name: missing}},
				"the checks show the package under its path")
		})

		origins := []struct {
			name   string
			target plugin.Target
			dir    string
			pkg    string
			give   string
			kind   symbol.Kind
			want   string
		}{
			{
				name: "takes the origin of the first error value", target: golang.Target, dir: goTree,
				pkg: storePackage, give: "store", kind: symbol.KindVariable, want: "ErrNotFound",
			},
			{
				name:   "takes the origin of the first error type of a package without error values",
				target: golang.Target, dir: goTree, pkg: typedPackage, give: "typed", kind: symbol.KindStruct,
				want: "LimitError",
			},
			{
				name:   "takes the origin of the first error class of a module in the order of its source",
				target: typescript.Target, dir: tsTree, pkg: storeModule, give: storeModule, kind: symbol.KindStruct,
				want: "NotFoundError",
			},
		}
		for _, tt := range origins {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := sentinel.New()
				f := fixture(t, tt.target, tt.dir, p, tt.pkg)
				res := f.Generate(t, p)
				assert.NoError(t, res.Err, "the generator runs whole")
				record := gentest.Records(t, res, testTemplate, sentinel.TestTag,
					func(c sentinel.Checks) string { return c.Subject })[tt.give]
				assert.NotNil(t, record.Decl, "the record of the package is a struct")
				expect.Equal(t, record.Decl.Origin, f.Find(t, tt.kind, "", tt.want).Identity(),
					"the record has the origin of the first error in the order of the source")
			})
		}

		values := []struct {
			name string
			pkg  string
			give string
			want map[string]string
		}{
			{
				name: "declares got as the address of a composite of the sample of each member", pkg: storePackage,
				give: "store",
				want: map[string]string{
					"NotFoundError": `got = &NotFoundError{Key: "Key-a", Attempt: 42}`,
					"WrappedError":  `got = &WrappedError{Op: "Op-a"}`,
				},
			},
			{
				name: "declares got as a composite of an error type with a value receiver", pkg: typedPackage,
				give: "typed", want: map[string]string{"LimitError": `got = LimitError{Name: "Name-a"}`},
			},
		}
		for _, tt := range values {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := sentinel.New()
				res := fixture(t, golang.Target, goTree, p, tt.pkg).Generate(t, p)
				assert.NoError(t, res.Err, "the generator runs whole")
				record := gentest.Records(t, res, testTemplate, sentinel.TestTag,
					func(c sentinel.Checks) string { return c.Subject })[tt.give]
				assert.NotNil(t, record.Decl, "the record of the package is a struct")
				got := map[string]string{}
				for _, k := range record.Data.Kinds {
					got[k.Name] = statements(record.Decl.Methods.Items()[k.Method].Body.Stmts)
				}
				expect.Equal(t, got, tt.want, "the method of each error type declares got")
			})
		}

		t.Run("records no construction of an error type of a language without constructors", func(t *testing.T) {
			t.Parallel()

			for _, k := range generate(t, golang.Target, goTree, storePackage).checks["store"].Kinds {
				expect.Equal(t, k.CauseMethod, -1, k.Name+" takes its cause by an assignment")
				expect.Equal(t, k.PeersMethod, -1, k.Name+" has no list of the other error types")
			}
		})

		t.Run("records no method of an error type without members", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, embeddedPackage).checks["embedded"]
			expect.Equal(t, got.Kinds[1].Method, -1, "LostError has no members")
		})

		findings := []struct {
			name     string
			pkg      string
			wantCode diag.Code
			wantMsg  string
			want     bool
		}{
			{
				name: "reports NoErrors at a package without errors", pkg: emptyPackage,
				wantCode: sentinel.NoErrors,
				wantMsg:  "the package empty has no error value and no error type, so it has no checks",
			},
			{
				name: "reports OwnNeighbour at a package that lists itself", pkg: selfPackage,
				wantCode: sentinel.OwnNeighbour, want: true,
				wantMsg: "the package self lists itself under no-overlap-with, and each error value matches itself",
			},
			{
				name: "reports Unloaded at a package that lists a package that the run did not load",
				pkg:  distantPackage, wantCode: sentinel.Unloaded, want: true,
				wantMsg: "the package distant lists example.com/shop/missing under no-overlap-with, and the run did " +
					"not load it, so the check of its error values compares with nothing",
			},
			{
				name:     "reports UnknownRemote at a remote type that is no error type",
				pkg:      strangerPackage,
				wantCode: sentinel.UnknownRemote,
				want:     true,
				wantMsg: "Point under remote is no error type of the package stranger, so its checks keep the message " +
					"of its zero value and its prefix",
			},
		}
		for _, tt := range findings {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := sentinel.New()
				f := fixture(t, golang.Target, goTree, p, tt.pkg)
				res := f.Generate(t, p)
				assert.NoError(t, res.Err, "the generator runs whole")
				var at node.Declaration
				for pkg := range f.Graph.Packages() {
					if pkg.ID.Package == tt.pkg {
						at = pkg
					}
				}
				assert.NotNil(t, at, "the fixture contains the package "+tt.pkg)
				got := gentest.Finding(t, res, at)
				expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of its case")
				expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the package and the reason")
				records := gentest.Records(t, res, testTemplate, sentinel.TestTag,
					func(c sentinel.Checks) string { return c.Subject })
				expect.Equal(t, len(records) == 1, tt.want, "sentinel records the checks where the finding allows")
			})
		}
	})
}

// names returns the names of the error values, joined by spaces.
func names(sentinels []sentinel.Sentinel) string {
	out := make([]string, 0, len(sentinels))
	for _, s := range sentinels {
		out = append(out, s.Name)
	}
	return strings.Join(out, " ")
}

// describe returns the facts of each error type: its name, its pointer, its
// hooks, its cause and its members, where it has them.
func describe(kinds []sentinel.Kind) []string {
	out := make([]string, 0, len(kinds))
	for _, k := range kinds {
		parts := []string{k.Name}
		if k.Pointer {
			parts = append(parts, "pointer")
		}
		if k.Compares {
			parts = append(parts, "compares")
		}
		if k.Unwraps {
			parts = append(parts, "unwraps")
		}
		if k.Cause != "" {
			parts = append(parts, "cause="+k.Cause)
		}
		members := make([]string, 0, len(k.Members))
		for _, m := range k.Members {
			if m.Text {
				members = append(members, m.Name+":text")
			} else {
				members = append(members, m.Name)
			}
		}
		if len(members) > 0 {
			parts = append(parts, "members="+strings.Join(members, ","))
		}
		out = append(out, strings.Join(parts, " "))
	}
	return out
}
