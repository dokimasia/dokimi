// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest_test

import (
	"errors"
	"strconv"
	"testing"

	"go.dokimi.dev/assert"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/frontendtest"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/gentest"
)

// The package shipping of the Go tree, the TypeScript tree and its module,
// the scripted tree and its package, a tree that does not exist, and a
// package that the Go tree does not declare.
const (
	shippingPackage = "example.com/shop/shipping"
	tsTree          = "testdata/typescript"
	catalogModule   = "catalog"
	scriptedTree    = "testdata/scripted"
	shopPackage     = "shop"
	missingTree     = "testdata/missing"
	missingPackage  = "example.com/shop/missing"
)

// The declarations that the cases find beside Address and its field Street:
// the fields of the class of the TypeScript tree, the struct of the scripted
// tree, and a name that no tree declares.
const (
	tsStreetName = "street"
	tsLabelsName = "labels"
	widgetName   = "Widget"
	missingName  = "Missing"
)

// The plugins of the cases: an annotator that stamps a fact on each struct,
// one that reports an Error at each struct, and one that returns an error.
const (
	stamperID  plugin.ID = "stamper"
	reporterID plugin.ID = "reporter"
	breakerID  plugin.ID = "breaker"
)

// stampedKey is the key of the fact that the stamper stamps on each struct.
const stampedKey meta.KeyName = "stamper.stamped"

// The openings of the contracts of the records that a stopped fixture
// writes.
const (
	storesContract   = "the frontend " + string(golang.Name) + " locates its stores"
	loadContract     = "the source tree loads"
	packageContract  = "the load declares the package "
	reportedContract = "the annotator reporter reports no Error, and reports "
	errorContract    = "the annotator breaker annotates"
)

// errBroken is the error that the breaker returns from each match.
var errBroken = errors.New("gentest_test: the annotator breaks")

// reported is the code of the Error that the reporter reports at each struct.
var reported = diag.MustRegister("GENTESTTEST", diag.CodeSpec{
	Number:  1,
	Meaning: "a struct of a fixture that the reporter of the tests of gentest reports",
})

// A fixture loads the packages of a tree in each of its languages under the
// brand of dokimi, stamps them through the annotators of its spec, and finds
// their declarations.
func TestFixture(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("loads the packages of a Go tree", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: goTree, Packages: []string{catalogPackage}})
			got := f.Find(t, symbol.KindStruct, "", addressName).Identity()
			assert.Equal(t, got.Package, catalogPackage, "the fixture loads Address from the package catalog")
		})

		t.Run("loads the modules of a TypeScript tree", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: tsTree, Packages: []string{catalogModule}})
			got := f.Find(t, symbol.KindField, addressName, tsStreetName).Identity()
			assert.Equal(t, got.Package, catalogModule, "the fixture loads the class Address from the module catalog")
		})

		t.Run("resolves a global type of a TypeScript tree to the library of TypeScript", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: tsTree, Packages: []string{catalogModule}})
			field, isField := f.Find(t, symbol.KindField, addressName, tsLabelsName).(*node.Field)
			assert.True(t, isField, "labels is a field")
			assert.NotNil(t, field.Type, "labels declares its type")
			assert.True(t, typescript.Library(field.Type.Target.Package),
				"the fixture loads the store of the library, which declares Map")
		})

		t.Run("loads the packages of a tree of the scripted language of eidos", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: scriptedTree, Packages: []string{shopPackage}})
			got := f.Find(t, symbol.KindStruct, "", widgetName).Identity()
			assert.Equal(t, got.Lang, frontendtest.ScriptedLang, "the fixture loads Widget in the scripted language")
		})

		t.Run("sets the target of the plan of the spec", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: goTree, Packages: []string{catalogPackage}, Target: golang.Target})
			assert.Equal(t, f.Target, golang.Target, "a phase call of a generator runs under the target of the spec")
		})

		t.Run("sets the spoke of the plan of the spec", func(t *testing.T) {
			t.Parallel()

			spoke := gobackend.New()
			f := gentest.New(t, &gentest.Spec{
				Dir: goTree, Packages: []string{catalogPackage}, Target: golang.Target, Types: spoke,
			})
			assert.Equal(t, f.Types, plugin.TypeSpeller(spoke),
				"a phase call of a generator translates types and previews routes through the spoke of the spec",
				assert.ByIdentity())
		})

		t.Run("stamps the fixture through the annotators of the spec", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{
				Dir: goTree, Packages: []string{catalogPackage}, Annotators: []plugin.Annotator{stamper()},
			})
			key, registered := meta.Lookup[bool](f.Keys, stampedKey)
			assert.True(t, registered, "the fixture registers the key of the stamper")
			got, _ := meta.Get(f.Facts, f.Find(t, symbol.KindStruct, "", addressName).Identity(), key)
			assert.True(t, got, "the stamper stamps Address before New returns")
		})

		t.Run("stops at a tree that does not load", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "a tree that does not exist", func(tb assert.TB) {
				gentest.New(tb, &gentest.Spec{Dir: missingTree, Packages: []string{catalogPackage}})
			})
			assert.Length(t, records, 1, "the fixture stops at the load")
			assert.Equal(t, records[0].Contract, loadContract, "the fixture stops at the load")
		})

		t.Run("stops at a package that the tree does not declare", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "a spec of a package that the tree lacks", func(tb assert.TB) {
				gentest.New(tb, &gentest.Spec{Dir: goTree, Packages: []string{missingPackage}})
			})
			assert.Length(t, records, 1, "the fixture stops at the packages")
			assert.Equal(t, records[0].Contract, packageContract+missingPackage, "the fixture stops at the package")
		})

		t.Run("stops at a frontend that locates no store", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "an environment without the roots of Go", func(tb assert.TB) {
				gentest.New(tb, &gentest.Spec{
					Dir: goTree, Packages: []string{catalogPackage}, Getenv: func(string) string { return "" },
				})
			})
			assert.Length(t, records, 1, "the fixture stops at the stores")
			assert.Equal(t, records[0].Contract, storesContract, "the fixture stops at the frontend of Go")
		})

		t.Run("fails at an Error that an annotator reports", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "an annotator that reports an Error", func(tb assert.TB) {
				gentest.New(tb, &gentest.Spec{
					Dir: goTree, Packages: []string{catalogPackage}, Annotators: []plugin.Annotator{reporter()},
				})
			})
			assert.Length(t, records, 2, "the fixture fails once for each Error, at Address and at Warehouse")
			assert.HasPrefix(t, records[0].Contract, reportedContract, "the fixture fails at the annotator")
		})

		t.Run("stops at an annotator that returns an error", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "an annotator that returns an error", func(tb assert.TB) {
				gentest.New(tb, &gentest.Spec{
					Dir: goTree, Packages: []string{catalogPackage}, Annotators: []plugin.Annotator{breaker()},
				})
			})
			assert.Length(t, records, 1, "the fixture stops at the error")
			assert.Equal(t, records[0].Contract, errorContract, "the fixture stops at the annotator")
		})
	})

	t.Run("Find", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the declaration of the kind that the owner declares under the name", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: goTree, Packages: []string{catalogPackage}})
			got := f.Find(t, symbol.KindField, addressName, streetName).Identity()
			assert.Equal(t, got.Name, streetName, "the declaration is the field Street of Address")
		})

		t.Run("returns the declaration of the kind among declarations with the same name", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: goTree, Packages: []string{catalogPackage, shippingPackage}})
			got := f.Find(t, symbol.KindStruct, "", addressName).Identity()
			assert.Equal(t, got.Package, catalogPackage, "Find returns the struct Address and not the function Address")
		})

		t.Run("stops at a name that the fixture does not declare", func(t *testing.T) {
			t.Parallel()

			f := gentest.New(t, &gentest.Spec{Dir: goTree, Packages: []string{catalogPackage}})
			records := assert.Rejects(t, "a name without a declaration", func(tb assert.TB) {
				f.Find(tb, symbol.KindStruct, "", missingName)
			})
			assert.Length(t, records, 1, "Find stops at the count of the declarations")
			assert.Equal(t, records[0].Contract,
				"the fixture declares one "+symbol.KindStruct.String()+" "+strconv.Quote(missingName)+` in ""`,
				"the contract of the record contains the kind, the name and the owner")
		})
	})
}

// stamper returns an annotator that stamps the fact of [stampedKey] on each
// struct.
func stamper() plugin.Annotator {
	var key meta.Key[bool]
	return sdk.NewPlugin(stamperID).
		Version(version).
		Keys(func(r *meta.Registry) error {
			if err := r.ClaimNamespace(stampedKey.Namespace()); err != nil {
				return err
			}
			var err error
			key, err = meta.Register[bool](r, meta.KeySpec{
				Name: stampedKey, Kinds: []symbol.Kind{symbol.KindStruct}, Doc: "marks each struct of a fixture",
			})
			return err
		}).
		Handle(sdk.OnStruct(func(_ *sdk.StructMatch, st *sdk.Stamper) error {
			sdk.Stamp(st, key, true)
			return nil
		})).
		Annotator()
}

// reporter returns an annotator that reports an Error at each struct.
func reporter() plugin.Annotator {
	return sdk.NewPlugin(reporterID).
		Version(version).
		Handle(sdk.OnStruct(func(m *sdk.StructMatch, _ *sdk.Stamper) error {
			m.Errorf(reported, "%s is reported", m.Struct.Name)
			return nil
		})).
		Annotator()
}

// breaker returns an annotator that returns an error from each match.
func breaker() plugin.Annotator {
	return sdk.NewPlugin(breakerID).
		Version(version).
		Handle(sdk.OnStruct(func(*sdk.StructMatch, *sdk.Stamper) error { return errBroken })).
		Annotator()
}
