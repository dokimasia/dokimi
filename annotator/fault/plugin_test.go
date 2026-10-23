// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fault_test

import (
	"fmt"
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/fault"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The Go tree and its packages, and the TypeScript tree and its module.
const (
	goTree        = "testdata/golang"
	storePackage  = "example.com/shop/store"
	badPackage    = "example.com/shop/bad"
	remotePackage = "example.com/shop/remote"
	tsTree        = "testdata/typescript"
	storeModule   = "store"
)

// The interfaces of the trees.
const (
	storeName  = "Store"
	ledgerName = "Ledger"
	remoteName = "Remote"
)

// otherRegistrant is the registrant that claims the namespace of fault
// before fault registers its keys.
const otherRegistrant = "other"

// stamp is what fault stamps on one method: the names of the error values of
// errors, the attempt of retry, the name of the parameter of partition, and
// whether it stamps the retry and the partition.
type stamp struct {
	errors      []string
	partition   string
	retry       int64
	retried     bool
	partitioned bool
}

// fault stamps the error values, the retry and the partition of the directive
// on a Go or a TypeScript method, reports each fault of an instance at the
// directive, and passes the plugin suite.
func TestFault(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the plugin suite over the Go package store", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := fault.New()
				return p, fixture(tb, goTree, storePackage, p).Fixture
			})
		})

		t.Run("passes the plugin suite over the TypeScript module store", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := fault.New()
				return p, fixture(tb, tsTree, storeModule, p).Fixture
			})
		})

		tests := []struct {
			name      string
			giveTree  string
			givePkg   string
			giveOwner string
			give      string
			want      stamp
		}{
			{
				name:     "stamps each param of the directive on a Go method",
				giveTree: goTree, givePkg: storePackage, giveOwner: storeName, give: "Get",
				want: stamp{
					errors: []string{"ErrNotFound", "ErrGone"}, retry: 3, retried: true,
					partition: "id", partitioned: true,
				},
			},
			{
				name:     "stamps an error of a package that another file of the package imports",
				giveTree: goTree, givePkg: remotePackage, giveOwner: remoteName, give: "Call",
				want: stamp{errors: []string{"ErrTimeout"}},
			},
			{
				name:     "stamps the attempt of retry on a Go method without a partition",
				giveTree: goTree, givePkg: storePackage, giveOwner: storeName, give: "List",
				want: stamp{retry: 2, retried: true},
			},
			{
				name:     "stamps nothing on a method without the directive",
				giveTree: goTree, givePkg: storePackage, giveOwner: storeName, give: "Close",
			},
			{
				name:     "stamps each param of the directive on a TypeScript method",
				giveTree: tsTree, givePkg: storeModule, giveOwner: storeName, give: "get",
				want: stamp{errors: []string{"notFound"}, retry: 3, retried: true, partition: "id", partitioned: true},
			},
			{
				name:     "stamps the attempt of retry on a TypeScript method that returns a value",
				giveTree: tsTree, givePkg: storeModule, giveOwner: storeName, give: "size",
				want: stamp{retry: 2, retried: true},
			},
			{
				name:     "stamps the parameter of partition of an instance whose retry it reports",
				giveTree: goTree, givePkg: badPackage, giveOwner: ledgerName, give: "Put",
				want: stamp{partition: "key", partitioned: true},
			},
			{
				name:     "stamps nothing on a Go method that cannot fail",
				giveTree: goTree, givePkg: badPackage, giveOwner: ledgerName, give: "Len",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := fault.New()
				f := fixture(t, tt.giveTree, tt.givePkg, p)
				assert.Equal(t, stamped(t, f, p, tt.giveOwner, tt.give), tt.want,
					"fault stamps the method "+tt.give+" as the case states")
			})
		}

		t.Run("reports each fault of an instance at the line of the directive", func(t *testing.T) {
			t.Parallel()

			p := fault.New()
			res := fixture(t, goTree, badPackage, p).Annotate(t, p)
			assert.NoError(t, res.Err, "the annotator runs whole")
			var got []string
			for d := range res.Sink.All() {
				got = append(got, fmt.Sprintf("%d %s: %s", d.Pos.Line, d.Code, d.Msg))
			}
			slices.Sort(got)
			assert.Equal(t, got, []string{
				"13 DOKIMIFAULT-0001: Len cannot fail, so a fault of its double has no error to return",
				"9 DOKIMIFAULT-0002: retry=1 is below 2, the first attempt that a retry schedule can let succeed",
			}, "fault reports the method that cannot fail and the retry below 2")
		})

		t.Run("returns from Keys the error of a namespace that another registrant claimed", func(t *testing.T) {
			t.Parallel()

			r := meta.NewRegistry()
			assert.NoError(t, r.For(otherRegistrant).ClaimNamespace(fault.KeyRetry.Namespace()),
				"another registrant claims the namespace of fault")
			keys, provides := fault.New().(plugin.KeyProvider)
			assert.True(t, provides, "fault registers keys")
			assert.HasError(t, keys.Keys(r.For(string(fault.ID))), "fault cannot claim its namespace")
		})

		t.Run("returns from Keys the error of a key registered with another type", func(t *testing.T) {
			t.Parallel()

			r := meta.NewRegistry().For(string(fault.ID))
			assert.NoError(t, r.ClaimNamespace(fault.KeyRetry.Namespace()), "fault claims its namespace")
			_, err := meta.Register[string](r, meta.KeySpec{Name: fault.KeyRetry, Doc: "a retry of another type"})
			assert.NoError(t, err, "a key of another type takes the name of the retry")
			keys, provides := fault.New().(plugin.KeyProvider)
			assert.True(t, provides, "fault registers keys")
			assert.HasError(t, keys.Keys(r), "fault cannot register its key of retry again")
		})
	})
}

// fixture returns a fixture over the package pkg of the tree dir, loaded
// under the context of tb, with the plugin fault p.
func fixture(tb assert.TB, dir, pkg string, p plugin.Plugin) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{Dir: dir, Packages: []string{pkg}, Plugins: []plugin.Plugin{p}})
}

// stamped runs fault p over the fixture f, and returns what p stamps on the
// method named name of the interface named owner. It stops the test where
// the stamped partition is no parameter of the method.
func stamped(t *testing.T, f *gentest.Fixture, p plugin.Plugin, owner, name string) stamp {
	t.Helper()

	assert.NoError(t, f.Annotate(t, p).Err, "the annotator runs whole")
	errorsKey, registered := meta.Lookup[[]symbol.Identity](f.Keys, fault.KeyErrors)
	assert.True(t, registered, "fault registers the key of errors")
	retryKey, registered := meta.Lookup[int64](f.Keys, fault.KeyRetry)
	assert.True(t, registered, "fault registers the key of retry")
	partitionKey, registered := meta.Lookup[symbol.Identity](f.Keys, fault.KeyPartition)
	assert.True(t, registered, "fault registers the key of partition")
	method, isMethod := f.Find(t, symbol.KindMethod, owner, name).(*node.Method)
	assert.True(t, isMethod, "the fixture declares the method "+name)
	values, _ := meta.Get(f.Facts, method.ID, errorsKey)
	retry, retried := meta.Get(f.Facts, method.ID, retryKey)
	partition, partitioned := meta.Get(f.Facts, method.ID, partitionKey)
	got := stamp{retry: retry, retried: retried, partitioned: partitioned}
	for _, v := range values {
		got.errors = append(got.errors, v.Name)
	}
	if partitioned {
		at := slices.IndexFunc(method.Params, func(p *node.Param) bool { return p.ID == partition })
		assert.NotEqual(t, at, -1, "the partition is a parameter of "+name)
		got.partition = method.Params[at].Name
	}
	return got
}
