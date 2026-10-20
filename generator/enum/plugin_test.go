// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The packages and modules that this spec alone reads.
const (
	priorityPackage = "example.com/shop/priority"
	splitPackage    = "example.com/shop/split"
	badPackage      = "example.com/shop/bad"
	badModule       = "bad"
)

// The enumerations that this spec alone reads.
const (
	priorityName = "Priority"
	zoneName     = "Zone"
	signalName   = "Signal"
	modeName     = "Mode"
	flagName     = "Flag"
	headerName   = "Header"
	gradeName    = "Grade"
)

// enum records the surface of each enumeration under its directive and the
// checks of the surface, in Go and in TypeScript, and passes the plugin
// suite.
func TestEnum(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the plugin suite over the Go package status", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := enum.New()
				return p, fixture(tb, golang.Target, goTree, p, statusPackage).Fixture
			})
		})

		t.Run("passes the plugin suite over the TypeScript module status", func(t *testing.T) {
			t.Parallel()

			plugintest.RunPluginSuite(t, func(tb assert.TB) (plugin.Plugin, *plugintest.Fixture) {
				p := enum.New()
				return p, fixture(tb, typescript.Target, tsTree, p, statusModule).Fixture
			})
		})

		t.Run("tests each enumeration from outside its package", func(t *testing.T) {
			t.Parallel()

			gentest.External(t, enum.New(), enum.TestTag)
		})

		variants := []struct {
			name        string
			target      plugin.Target
			tree        string
			pkg         string
			give        string
			want        []string
			wantTextual bool
		}{
			{
				name: "records the name of each variant of a number as its text", target: golang.Target,
				tree: goTree, pkg: statusPackage, give: statusName,
				want: []string{`Draft="Draft"`, `Published="Published"`, `Archived="Archived"`},
			},
			{
				name: "records the value of each variant of text as its text", target: golang.Target, tree: goTree,
				pkg: regionPackage, give: regionName, wantTextual: true,
				want: []string{`US="us-east"`, `EU="eu-west"`, `AP="ap-south"`},
			},
			{
				name: "records the variants that another file of the package declares", target: golang.Target,
				tree: goTree, pkg: splitPackage, give: signalName,
				want: []string{`Green="Green"`, `Amber="Amber"`, `Red="Red"`},
			},
			{
				name:   "records the name of each member of a TypeScript enum of numbers as its text",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: statusName,
				want: []string{`Draft="Draft"`, `Published="Published"`, `Archived="Archived"`},
			},
			{
				name:   "records the value of each member of a TypeScript enum of text as its text",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: regionName, wantTextual: true,
				want: []string{`US="us-east"`, `EU="eu-west"`, `AP="ap-south"`},
			},
			{
				name:   "records the trailing comment of each variant as its text under text=comment",
				target: golang.Target, tree: goTree, pkg: commentedPackage, give: gradeName,
				want: []string{`Poor="poor"`, `Fair="fair enough"`, `Good="Good"`},
			},
			{
				name:   "records the name of each variant with a trailing comment as its text without text=comment",
				target: golang.Target, tree: goTree, pkg: commentedPackage, give: "Mood",
				want: []string{`Calm="Calm"`, `Angry="Angry"`},
			},
			{
				name:   "records the trailing comment of each TypeScript member as its text under text=comment",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: gradeName,
				want: []string{`Poor="poor"`, `Fair="fair enough"`, `Good="Good"`},
			},
		}
		for _, tt := range variants {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, recorded := generate(t, tt.target, tt.tree, tt.pkg).models[tt.give]
				assert.True(t, recorded, "enum records the surface of "+tt.give)
				expect.Equal(t, texts(got.Variants), tt.want, "each variant has the text of its form")
				expect.Equal(t, got.Textual, tt.wantTextual, "the text of a variant of text is its value")
			})
		}

		t.Run("records the name of the declaration of each variant as a value", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, splitPackage).models[signalName]
			callees := make([]string, 0, len(got.Variants))
			for _, v := range got.Variants {
				expect.Equal(t, v.Value.Kind, emit.ValueName, "the value of "+v.Name+" is a name")
				callees = append(callees, v.Value.Callee.Package+"."+v.Value.Callee.Name)
			}
			assert.Equal(t, callees, []string{splitPackage + ".Green", splitPackage + ".Amber", splitPackage + ".Red"},
				"each value names the constant of its variant in the package of the enumeration")
		})

		stated := []struct {
			name   string
			target plugin.Target
			tree   string
			pkg    string
			give   string
			want   []string
		}{
			{
				name:   "states the text of each variant with a trailing comment under text=comment",
				target: golang.Target, tree: goTree, pkg: commentedPackage, give: gradeName,
				want: []string{"Poor", "Fair"},
			},
			{
				name: "states the text of each variant of text", target: golang.Target, tree: goTree,
				pkg: regionPackage, give: regionName, want: []string{"US", "EU", "AP"},
			},
			{
				name: "states the text of no variant whose text is its name", target: golang.Target, tree: goTree,
				pkg: statusPackage, give: statusName,
			},
			{
				name:   "states the text of each TypeScript member with a trailing comment under text=comment",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: gradeName,
				want: []string{"Poor", "Fair"},
			},
		}
		for _, tt := range stated {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, recorded := generate(t, tt.target, tt.tree, tt.pkg).checks[tt.give]
				assert.True(t, recorded, "enum records the checks of "+tt.give)
				var names []string
				for _, v := range got.Variants {
					if v.Stated {
						names = append(names, v.Name)
					}
				}
				expect.Equal(t, names, tt.want, "the checks pin the text of each variant whose declaration states it")
				expect.Equal(t, got.Stated, tt.want != nil, "the checks pin a text where a variant states one")
			})
		}

		surfaces := []struct {
			name   string
			target plugin.Target
			tree   string
			pkg    string
			give   string
			want   enum.Members
			calls  enum.Members
		}{
			{
				name: "records each member of the surface of an enumeration without one", target: golang.Target,
				tree: goTree, pkg: statusPackage, give: statusName,
				want:  enum.Members{Render: true, Parse: true, Values: true, Valid: true, Encode: true, Decode: true},
				calls: enum.Members{Render: true, Parse: true, Values: true, Valid: true, Encode: true, Decode: true},
			},
			{
				name:   "leaves out the render, the parse function and the list where the type renders itself",
				target: golang.Target, tree: goTree, pkg: handwrittenPackage, give: weekdayName,
				want:  enum.Members{Valid: true, Encode: true, Decode: true},
				calls: enum.Members{Render: true, Valid: true, Encode: true, Decode: true},
			},
			{
				name:   "leaves out the surface under methods=off and checks the members of the type",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: modeName,
				calls: enum.Members{Render: true, Valid: true, Encode: true},
			},
			{
				name:   "leaves out the surface under methods=off and checks no member without one",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Plain",
			},
			{
				name:   "keeps the decode where the type declares the encode, and checks no round trip",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Tone",
				want:  enum.Members{Render: true, Parse: true, Values: true, Valid: true, Decode: true},
				calls: enum.Members{Render: true, Parse: true, Values: true, Valid: true, Encode: true},
			},
			{
				name:   "leaves out the encode where neither the surface nor the package declares a parse function",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Phase",
				want:  enum.Members{Valid: true},
				calls: enum.Members{Render: true, Valid: true},
			},
			{
				name:   "leaves out the decode that the type declares, and keeps the encode",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Speed",
				want:  enum.Members{Render: true, Parse: true, Values: true, Valid: true, Encode: true},
				calls: enum.Members{Render: true, Parse: true, Values: true, Valid: true, Encode: true},
			},
			{
				name:   "records the functions of the surface of a TypeScript enum",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: statusName,
				want:  enum.Members{Render: true, Parse: true, Values: true, Valid: true},
				calls: enum.Members{Render: true, Parse: true, Values: true, Valid: true},
			},
			{
				name:   "leaves out the functions that go with a render that the module declares",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: modeName,
				want:  enum.Members{Valid: true},
				calls: enum.Members{Render: true, Valid: true},
			},
		}
		for _, tt := range surfaces {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, tt.target, tt.tree, tt.pkg)
				expect.Equal(t, got.models[tt.give].Generate, tt.want, "the surface declares what is missing")
				assert.Contains(t, got.checks, tt.give, "enum records the checks of "+tt.give)
				expect.Equal(t, got.checks[tt.give].Calls, tt.calls,
					"the checks call each member that the surface or the type declares")
				expect.Equal(t, got.checks[tt.give].Generate, tt.want, "the checks know what the surface declares")
			})
		}

		t.Run("records no surface where the directive leaves out every member", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, quietPackage)
			expect.NotContains(t, got.models, modeName, "Mode has no surface")
			expect.Contains(t, got.models, levelName, "Level has a surface")
		})

		zeros := []struct {
			name   string
			pkg    string
			give   string
			want   string
			absent bool
		}{
			{
				name: "records the variant whose value is the zero value", pkg: priorityPackage, give: priorityName,
				want: "Low",
			},
			{
				name: "records the zero variant where it is the last variant", pkg: regionPackage, give: zoneName,
				want: "Unset",
			},
			{
				name: "records no zero variant where no variant has the zero value", pkg: statusPackage,
				give: statusName, absent: true,
			},
		}
		for _, tt := range zeros {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, golang.Target, goTree, tt.pkg).checks[tt.give]
				if tt.absent {
					expect.Nil(t, got.Zero, tt.give+" has no zero variant")
					return
				}
				assert.NotNil(t, got.Zero, tt.give+" has a zero variant")
				expect.Equal(t, got.Zero.Name, tt.want, "the zero variant has the zero value")
			})
		}

		probes := []struct {
			name           string
			target         plugin.Target
			tree           string
			pkg            string
			give           string
			wantUnknown    string
			wantOutOfRange string
		}{
			{
				name: "records a text longer than every text, and the value after the last", target: golang.Target,
				tree: goTree, pkg: statusPackage, give: statusName, wantUnknown: `"??????????"`,
				wantOutOfRange: "4",
			},
			{
				name: "records the empty text as the value outside an enumeration of text", target: golang.Target,
				tree: goTree, pkg: regionPackage, give: regionName, wantUnknown: `"?????????"`,
				wantOutOfRange: `""`,
			},
			{
				name:   "records a text of marks as the value outside an enumeration with the empty text",
				target: golang.Target, tree: goTree, pkg: regionPackage, give: zoneName, wantUnknown: `"??????"`,
				wantOutOfRange: `"??????"`,
			},
			{
				name: "records the value after the last member of a TypeScript enum", target: typescript.Target,
				tree: tsTree, pkg: statusModule, give: statusName, wantUnknown: `"??????????"`,
				wantOutOfRange: "4",
			},
			{
				name:   "records no value outside an enumeration that takes every value of its type",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Switch", wantUnknown: `"?????"`,
			},
			{
				name:   "records a text longer than every text that a trailing comment states",
				target: golang.Target, tree: goTree, pkg: commentedPackage, give: gradeName,
				wantUnknown: `"????????????"`, wantOutOfRange: "3",
			},
		}
		for _, tt := range probes {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, tt.target, tt.tree, tt.pkg).checks[tt.give]
				expect.Equal(t, got.Unknown, tt.wantUnknown, "no variant has the unknown text")
				expect.Equal(t, got.OutOfRange, tt.wantOutOfRange, "no variant has the value outside the set")
			})
		}

		t.Run("records the package that opens the message of the error of the parse function", func(t *testing.T) {
			t.Parallel()

			got := generate(t, golang.Target, goTree, statusPackage).models[statusName]
			expect.Equal(t, got.Package, "status", "the message of ErrUnknownStatus opens with status")
		})

		t.Run("names the parse function after the option parse-word", func(t *testing.T) {
			t.Parallel()

			p := enum.New()
			p.(plugin.OptionsProvider).Options().(*enum.Options).ParseWord = "decode"
			got := run(t, p, golang.Target, goTree, statusPackage)
			expect.Equal(t, got.models[statusName].Names.Parse, "DecodeStatus", "the parse function is DecodeStatus")
		})

		t.Run("names the error of the parse function after the option sentinel-word", func(t *testing.T) {
			t.Parallel()

			p := enum.New()
			p.(plugin.OptionsProvider).Options().(*enum.Options).SentinelWord = "unrecognised"
			got := run(t, p, golang.Target, goTree, statusPackage)
			expect.Equal(t, got.models[statusName].Names.Sentinel, "ErrUnrecognisedStatus",
				"the error is ErrUnrecognisedStatus")
		})

		words := []struct {
			name         string
			give         func(*enum.Options)
			wantParse    string
			wantSentinel string
		}{
			{
				name: "names the parse function after the word parse where the option is empty",
				give: func(o *enum.Options) { o.ParseWord = "" }, wantParse: "ParseStatus",
				wantSentinel: "ErrUnknownStatus",
			},
			{
				name: "names the error after the word unknown where the option is empty",
				give: func(o *enum.Options) { o.SentinelWord = "" }, wantParse: "ParseStatus",
				wantSentinel: "ErrUnknownStatus",
			},
		}
		for _, tt := range words {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := enum.New()
				options := p.(plugin.OptionsProvider).Options().(*enum.Options)
				expect.Equal(t, *options, enum.Options{ParseWord: enum.ParseWord, SentinelWord: enum.SentinelWord},
					"the options have the default words")
				tt.give(options)
				got := run(t, p, golang.Target, goTree, statusPackage).models[statusName].Names
				expect.Equal(t, got.Parse, tt.wantParse, "the parse function has the default word")
				expect.Equal(t, got.Sentinel, tt.wantSentinel, "the error has the default word")
			})
		}

		t.Run("records no checks under the option no-tests", func(t *testing.T) {
			t.Parallel()

			p := enum.New()
			p.(plugin.OptionsProvider).Options().(*enum.Options).NoTests = true
			got := run(t, p, golang.Target, goTree, statusPackage)
			expect.Contains(t, got.models, statusName, "enum records the surface of Status")
			expect.Empty(t, got.checks, "enum records no checks")
		})

		findings := []struct {
			name     string
			target   plugin.Target
			tree     string
			pkg      string
			give     string
			wantCode diag.Code
			wantMsg  string
		}{
			{
				name: "reports Duplicate at an enumeration whose two variants have one value", target: golang.Target,
				tree: goTree, pkg: badPackage, give: "Shade", wantCode: enum.Duplicate,
				wantMsg: "a variant of Shade after Dark has the value of Dark, so it has no surface",
			},
			{
				name: "reports Duplicate at an enumeration whose two variants have one text", target: golang.Target,
				tree: goTree, pkg: badPackage, give: "Code", wantCode: enum.Duplicate,
				wantMsg: "the variants CodeA and A of Code have one text, so it has no surface",
			},
			{
				name:   "reports Duplicate at an enumeration whose two trailing comments state one text",
				target: golang.Target, tree: goTree, pkg: badPackage, give: "Size", wantCode: enum.Duplicate,
				wantMsg: "the variants Small and Large of Size have one text, so it has no surface",
			},
			{
				name:   "reports Duplicate at a TypeScript enum whose two members have one value",
				target: typescript.Target, tree: tsTree, pkg: badModule, give: "Shade", wantCode: enum.Duplicate,
				wantMsg: "a variant of Shade after Dark has the value of Dark, so it has no surface",
			},
			{
				name: "reports Empty at a TypeScript enum without a member", target: typescript.Target, tree: tsTree,
				pkg: badModule, give: "Nothing", wantCode: enum.Empty,
				wantMsg: "Nothing has no variant, so it has no surface",
			},
		}
		for _, tt := range findings {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := enum.New()
				f := fixture(t, tt.target, tt.tree, p, tt.pkg)
				res := f.Generate(t, p)
				assert.NoError(t, res.Err, "the generator runs whole")
				at := f.Find(t, symbol.KindEnum, "", tt.give)
				got := gentest.Finding(t, res, at)
				expect.Equal(t, got.Code, tt.wantCode, "the finding has the code of its case")
				expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the enumeration and the reason")
				gentest.Unrecorded(t, res, at)
			})
		}
	})
}
