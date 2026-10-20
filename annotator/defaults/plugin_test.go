// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package defaults_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The Go tree and its packages, and the TypeScript tree and its module.
const (
	goTree         = "testdata/golang"
	catalogPackage = "example.com/shop/catalog"
	taggedPackage  = "example.com/shop/tagged"
	tsTree         = "testdata/typescript"
	catalogModule  = "catalog"
)

// The struct Address of the Go tree and its fields, and the class Address of
// the TypeScript tree and its fields. The Go tree declares the defined types
// Key and Level, and the TypeScript tree the type alias Level.
const (
	addressName  = "Address"
	streetName   = "Street"
	cityName     = "City"
	countryName  = "Country"
	floorName    = "Floor"
	linesName    = "Lines"
	tsStreetName = "street"
	tsCityName   = "city"
	keyName      = "Key"
	levelName    = "Level"
)

// The defaults that the trees declare on the fields of Address, on the type
// Key of the Go tree and on the type Level of the TypeScript tree.
const (
	streetDefault  = "Dam 1"
	cityDefault    = "Amsterdam"
	countryDefault = "NL"
	floorDefault   = "0"
	keyDefault     = "test-key"
	levelDefault   = "3"
)

// The struct Config of the package tagged and its fields, the second key of
// the tag of Region, and the default under each key.
const (
	configName      = "Config"
	regionName      = "Region"
	zoneName        = "Zone"
	fallbackTag     = "fallback"
	fallbackDefault = "west"
	regionDefault   = "east"
)

// otherRegistrant is the registrant that claims the namespace of defaults
// before defaults registers its keys.
const otherRegistrant = "other"

// stamp is what defaults stamps on one field: the text of the value, the
// source, and whether it stamps the value at all.
type stamp struct {
	value   string
	source  string
	stamped bool
}

// defaults stamps the declared defaults of the fields of Go and
// TypeScript, the directive before the tag, and passes the plugin suite.
func TestDefaults(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the plugin suite over the Go package catalog", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := defaults.New()
				return p, fixture(tb, goTree, catalogPackage, p).Fixture
			})
		})

		t.Run("passes the plugin suite over the TypeScript module catalog", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := defaults.New()
				return p, fixture(tb, tsTree, catalogModule, p).Fixture
			})
		})

		tests := []struct {
			name string
			give string
			want stamp
		}{
			{
				name: "stamps the value of the directive on a Go field under the directive", give: cityName,
				want: stamp{value: cityDefault, source: defaults.SourceDirective, stamped: true},
			},
			{
				name: "stamps the value of the directive on a Go field with the tag too", give: streetName,
				want: stamp{value: streetDefault, source: defaults.SourceDirective, stamped: true},
			},
			{
				name: "stamps the value of the tag on a Go field with the tag", give: countryName,
				want: stamp{value: countryDefault, source: defaults.SourceTag, stamped: true},
			},
			{
				name: "stamps the text of a number in the tag", give: floorName,
				want: stamp{value: floorDefault, source: defaults.SourceTag, stamped: true},
			},
			{name: "stamps nothing on a Go field without a default", give: linesName},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := defaults.New()
				f := fixture(t, goTree, catalogPackage, p)
				assert.Equal(t, stamped(t, f, p, symbol.KindField, addressName, tt.give), tt.want,
					"defaults stamps the field "+tt.give+" as the case states")
			})
		}

		t.Run("stamps the value of the directive on a TypeScript field", func(t *testing.T) {
			t.Parallel()

			p := defaults.New()
			f := fixture(t, tsTree, catalogModule, p)
			assert.Equal(t, stamped(t, f, p, symbol.KindField, addressName, tsStreetName),
				stamp{value: streetDefault, source: defaults.SourceDirective, stamped: true},
				"defaults stamps the directive on the field street")
		})

		t.Run("stamps nothing on a TypeScript field without the directive", func(t *testing.T) {
			t.Parallel()

			p := defaults.New()
			f := fixture(t, tsTree, catalogModule, p)
			assert.Equal(t, stamped(t, f, p, symbol.KindField, addressName, tsCityName), stamp{},
				"the rules of TypeScript state no tags, so city has no default")
		})

		{
			tests := []struct {
				name     string
				giveTree string
				givePkg  string
				give     string
				want     stamp
			}{
				{
					name: "stamps the value of the directive on a Go defined type", giveTree: goTree,
					givePkg: catalogPackage, give: keyName,
					want: stamp{value: keyDefault, source: defaults.SourceDirective, stamped: true},
				},
				{
					name: "stamps nothing on a Go defined type without the directive", giveTree: goTree,
					givePkg: catalogPackage, give: levelName,
				},
				{
					name: "stamps the value of the directive on a TypeScript type alias", giveTree: tsTree,
					givePkg: catalogModule, give: levelName,
					want: stamp{value: levelDefault, source: defaults.SourceDirective, stamped: true},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := defaults.New()
					f := fixture(t, tt.giveTree, tt.givePkg, p)
					assert.Equal(t, stamped(t, f, p, symbol.KindAlias, "", tt.give), tt.want,
						"defaults stamps the type "+tt.give+" as the case states")
				})
			}
		}

		t.Run("reads the tag default without a config", func(t *testing.T) {
			t.Parallel()

			p := defaults.New()
			f := fixture(t, goTree, taggedPackage, p)
			assert.Equal(t, stamped(t, f, p, symbol.KindField, configName, regionName),
				stamp{value: regionDefault, source: defaults.SourceTag, stamped: true},
				"defaults reads the tag of the key default")
		})

		t.Run("stamps nothing on a field whose tag has an empty default", func(t *testing.T) {
			t.Parallel()

			p := defaults.New()
			f := fixture(t, goTree, taggedPackage, p)
			assert.Equal(t, stamped(t, f, p, symbol.KindField, configName, zoneName), stamp{},
				"an empty tag value declares no default")
		})

		t.Run("reads the tag of the key that the option tag-key sets", func(t *testing.T) {
			t.Parallel()

			p := defaults.New()
			p.(plugin.OptionsProvider).Options().(*defaults.Options).TagKey = fallbackTag
			f := fixture(t, goTree, taggedPackage, p)
			assert.Equal(t, stamped(t, f, p, symbol.KindField, configName, regionName),
				stamp{value: fallbackDefault, source: defaults.SourceTag, stamped: true},
				"defaults reads the tag of the key that the option sets")
		})

		t.Run("reads no tag where the option no-tags is set", func(t *testing.T) {
			t.Parallel()

			p := defaults.New()
			p.(plugin.OptionsProvider).Options().(*defaults.Options).NoTags = true
			f := fixture(t, goTree, taggedPackage, p)
			assert.Equal(t, stamped(t, f, p, symbol.KindField, configName, regionName), stamp{},
				"the directive is the only form of a default under no-tags")
		})

		t.Run("returns from Keys the error of a namespace that another registrant claimed", func(t *testing.T) {
			t.Parallel()

			r := meta.NewRegistry()
			assert.NoError(t, r.For(otherRegistrant).ClaimNamespace(defaults.KeyValue.Namespace()),
				"another registrant claims the namespace of defaults")
			keys, provides := defaults.New().(plugin.KeyProvider)
			assert.True(t, provides, "defaults registers keys")
			assert.HasError(t, keys.Keys(r.For(string(defaults.ID))), "defaults cannot claim its namespace")
		})

		t.Run("returns from Keys the error of a key registered with another type", func(t *testing.T) {
			t.Parallel()

			r := meta.NewRegistry().For(string(defaults.ID))
			assert.NoError(t, r.ClaimNamespace(defaults.KeyValue.Namespace()), "defaults claims its namespace")
			_, err := meta.Register[int64](r, meta.KeySpec{Name: defaults.KeyValue, Doc: "a value of another type"})
			assert.NoError(t, err, "a key of another type takes the name of the value")
			keys, provides := defaults.New().(plugin.KeyProvider)
			assert.True(t, provides, "defaults registers keys")
			assert.HasError(t, keys.Keys(r), "defaults cannot register its key of the value again")
		})
	})
}

// fixture returns a fixture over the package pkg of the tree dir, loaded
// under the context of tb, with the plugin defaults p.
func fixture(tb assert.TB, dir, pkg string, p plugin.Plugin) *gentest.Fixture {
	tb.Helper()

	return gentest.New(tb, &gentest.Spec{Dir: dir, Packages: []string{pkg}, Plugins: []plugin.Plugin{p}})
}

// stamped runs defaults p over the fixture f, and returns what p stamps on
// the declaration named name of the given kind. The declaration is a member
// of the declaration named owner, or at file level for an empty owner.
func stamped(t *testing.T, f *gentest.Fixture, p plugin.Plugin, kind symbol.Kind, owner, name string) stamp {
	t.Helper()

	assert.NoError(t, f.Annotate(t, p).Err, "the annotator runs whole")
	valueKey, registered := meta.Lookup[string](f.Keys, defaults.KeyValue)
	assert.True(t, registered, "defaults registers the key of the value")
	sourceKey, registered := meta.Lookup[string](f.Keys, defaults.KeySource)
	assert.True(t, registered, "defaults registers the key of the source")
	id := f.Find(t, kind, owner, name).Identity()
	value, held := meta.Get(f.Facts, id, valueKey)
	source, _ := meta.Get(f.Facts, id, sourceKey)
	return stamp{value: value, source: source, stamped: held}
}
