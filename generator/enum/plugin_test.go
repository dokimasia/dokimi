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
	customPackage   = "example.com/shop/custom"
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

		t.Run("checks each enumeration from outside its package", func(t *testing.T) {
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

		{
			tests := []struct {
				name    string
				givePkg string
				give    string
				want    bool
			}{
				{
					name:    "records that the render converts a value of text whose text is its value",
					givePkg: regionPackage, give: regionName, want: true,
				},
				{
					name:    "records that the render converts no value of text whose trailing comment states its text",
					givePkg: regionPackage, give: "Lane",
				},
				{
					name:    "records that the render converts no value of a number",
					givePkg: statusPackage, give: statusName,
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					got, recorded := generate(t, golang.Target, goTree, tt.givePkg).models[tt.give]
					assert.True(t, recorded, "enum records the surface of "+tt.give)
					expect.Equal(t, got.Converts, tt.want, "the render converts a value where each text is the value")
				})
			}
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
				want: enum.Members{
					Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
				},
				calls: enum.Members{
					Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
				},
			},
			{
				name:   "leaves out each member that the package declares",
				target: golang.Target, tree: goTree, pkg: handwrittenPackage, give: weekdayName,
				want: enum.Members{Values: true, Valid: true, Encode: true, Decode: true},
				calls: enum.Members{
					Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
				},
			},
			{
				name:   "wraps the error that the package declares in the parse function",
				target: golang.Target, tree: goTree, pkg: declaredPackage, give: "Level",
				want: enum.Members{Render: true, Parse: true, Values: true, Valid: true, Encode: true, Decode: true},
				calls: enum.Members{
					Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
				},
			},
			{
				name:   "leaves out the list of the variants that the package declares",
				target: golang.Target, tree: goTree, pkg: declaredPackage, give: "Color",
				want: enum.Members{Render: true, Parse: true, Error: true, Valid: true, Encode: true, Decode: true},
				calls: enum.Members{
					Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
				},
			},
			{
				name:   "keeps the render beside a parse function that the package declares",
				target: golang.Target, tree: goTree, pkg: declaredPackage, give: "Lane",
				want:  enum.Members{Render: true, Values: true, Valid: true, Encode: true, Decode: true},
				calls: enum.Members{Render: true, Parse: true, Values: true, Valid: true, Encode: true, Decode: true},
			},
			{
				name:   "writes the parse function beside a method of its name on another type",
				target: golang.Target, tree: goTree, pkg: declaredPackage, give: "Coil",
				want: enum.Members{
					Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
				},
				calls: enum.Members{
					Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
				},
			},
			{
				name:   "checks no parse function without a render",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Gauge",
			},
			{
				name:   "checks the members of the type without a surface under methods=off",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: modeName,
				calls: enum.Members{Render: true, Valid: true, Encode: true},
			},
			{
				name:   "checks no member of a type without members under methods=off",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Plain",
			},
			{
				name:   "keeps the decode of a type that declares the encode",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Tone",
				want:  enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true, Decode: true},
				calls: enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true},
			},
			{
				name:   "leaves out the encode where neither the surface nor the package declares a parse function",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Phase",
				want:  enum.Members{Values: true, Valid: true},
				calls: enum.Members{Render: true, Values: true, Valid: true},
			},
			{
				name:   "keeps the encode of a type that declares the decode",
				target: golang.Target, tree: goTree, pkg: quietPackage, give: "Speed",
				want:  enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true},
				calls: enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true},
			},
			{
				name:   "records the functions of the surface of a TypeScript enum",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: statusName,
				want:  enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true},
				calls: enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true},
			},
			{
				name:   "leaves out the parse function beside a render that the module declares",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: modeName,
				want:  enum.Members{Values: true, Valid: true},
				calls: enum.Members{Render: true, Values: true, Valid: true},
			},
			{
				name:   "throws the error class that the module declares from the parse function",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: "Level",
				want:  enum.Members{Render: true, Parse: true, Values: true, Valid: true},
				calls: enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true},
			},
			{
				name:   "leaves out the type guard that the module declares",
				target: typescript.Target, tree: tsTree, pkg: statusModule, give: "Phase",
				want:  enum.Members{Render: true, Parse: true, Error: true, Values: true},
				calls: enum.Members{Render: true, Parse: true, Error: true, Values: true, Valid: true},
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
				name: "records the probes outside an enumeration of numbers", target: golang.Target,
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

		t.Run("starts from the default template of each name", func(t *testing.T) {
			t.Parallel()

			got := enum.New().(plugin.OptionsProvider).Options().(*enum.Options)
			assert.Equal(t, *got, enum.Options{
				ParseFunction: enum.DefaultParseFunction, ValuesFunction: enum.DefaultValuesFunction,
				ParseError: enum.DefaultParseError, ValidMethod: enum.DefaultValidMethod,
			}, "the options have the default templates")
		})

		templates := []struct {
			name string
			give func(*enum.Options)
			want enum.Names
		}{
			{
				name: "names the parse function after the option parse-function",
				give: func(o *enum.Options) { o.ParseFunction = "Decode{{.Enum}}" },
				want: enum.Names{
					Render: "String", Parse: "DecodeStatus", Values: "StatusValues", Valid: "IsValid",
					Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnknownStatus",
				},
			},
			{
				name: "names the function that lists the variants after the option values-function",
				give: func(o *enum.Options) { o.ValuesFunction = "Every{{.Enum}}" },
				want: enum.Names{
					Render: "String", Parse: "ParseStatus", Values: "EveryStatus", Valid: "IsValid",
					Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnknownStatus",
				},
			},
			{
				name: "names the error of the parse function after the option parse-error",
				give: func(o *enum.Options) { o.ParseError = "Unrecognised{{.Enum}}" },
				want: enum.Names{
					Render: "String", Parse: "ParseStatus", Values: "StatusValues", Valid: "IsValid",
					Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnrecognisedStatus",
				},
			},
			{
				name: "names the validity method after the option valid-method",
				give: func(o *enum.Options) { o.ValidMethod = "Valid" },
				want: enum.Names{
					Render: "String", Parse: "ParseStatus", Values: "StatusValues", Valid: "Valid",
					Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnknownStatus",
				},
			},
			{
				name: "writes the words of a template in PascalCase",
				give: func(o *enum.Options) { o.ParseFunction = "parse {{.Enum}} text"; o.ValidMethod = "is_known" },
				want: enum.Names{
					Render: "String", Parse: "ParseStatusText", Values: "StatusValues", Valid: "IsKnown",
					Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnknownStatus",
				},
			},
			{
				name: "names each member after its default template where each option is empty",
				give: func(o *enum.Options) { *o = enum.Options{} },
				want: enum.Names{
					Render: "String", Parse: "ParseStatus", Values: "StatusValues", Valid: "IsValid",
					Encode: "MarshalText", Decode: "UnmarshalText", Sentinel: "ErrUnknownStatus",
				},
			},
		}
		for _, tt := range templates {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := enum.New()
				tt.give(p.(plugin.OptionsProvider).Options().(*enum.Options))
				got := run(t, p, golang.Target, goTree, statusPackage).models[statusName].Names
				assert.Equal(t, got, tt.want, "the surface of Status has the names that the options render")
			})
		}

		faults := []struct {
			name    string
			give    func(*enum.Options)
			wantErr string
		}{
			{
				name:    "returns the error of an option parse-function that does not render",
				give:    func(o *enum.Options) { o.ParseFunction = "Parse{{.Name}}" },
				wantErr: `enum: the name template "Parse{{.Name}}" has an action other than {{.Enum}}`,
			},
			{
				name:    "returns the error of an option values-function that does not render",
				give:    func(o *enum.Options) { o.ValuesFunction = "{{.Enum}}!" },
				wantErr: `enum: the name template "{{.Enum}}!" contains the character '!'`,
			},
			{
				name:    "returns the error of an option parse-error that does not render",
				give:    func(o *enum.Options) { o.ParseError = "{{" },
				wantErr: `enum: the name template "{{" does not parse`,
			},
			{
				name:    "returns the error of an option valid-method that does not render",
				give:    func(o *enum.Options) { o.ValidMethod = "2valid" },
				wantErr: `enum: the name template "2valid" starts with '2' instead of a letter`,
			},
		}
		for _, tt := range faults {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := enum.New()
				tt.give(p.(plugin.OptionsProvider).Options().(*enum.Options))
				res := fixture(t, golang.Target, goTree, p, statusPackage).Generate(t, p)
				assert.HasError(t, res.Err, "the generator stops at the template")
				expect.Contains(t, res.Err.Error(), tt.wantErr, "the error contains the template and the rule")
			})
		}

		{
			tests := []struct {
				name       string
				giveMethod enum.MethodTemplate
				want       enum.Members
				wantValid  string
			}{
				{
					name:       "writes a second validity method beside a method of another name",
					giveMethod: enum.DefaultValidMethod,
					want: enum.Members{
						Render: true, Parse: true, Error: true, Values: true, Valid: true, Encode: true, Decode: true,
					},
					wantValid: "IsValid",
				},
				{
					name:       "takes the validity method that the type declares under the option valid-method",
					giveMethod: "Valid",
					want: enum.Members{
						Render: true, Parse: true, Error: true, Values: true, Encode: true, Decode: true,
					},
					wantValid: "Valid",
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					p := enum.New()
					p.(plugin.OptionsProvider).Options().(*enum.Options).ValidMethod = tt.giveMethod
					got := run(t, p, golang.Target, goTree, customPackage)
					expect.Equal(t, got.models["Gear"].Generate, tt.want, "the surface declares what is missing")
					expect.True(t, got.checks["Gear"].Calls.Valid, "the checks call the validity method")
					expect.Equal(t, got.checks["Gear"].Names.Valid, tt.wantValid,
						"the checks call the method that the option names")
				})
			}
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

		conflicts := []struct {
			name    string
			target  plugin.Target
			tree    string
			pkg     string
			give    string
			atKind  symbol.Kind
			atOwner string
			at      string
			wantMsg string
		}{
			{
				name: "reports Conflict at a parse function of bytes", target: golang.Target, tree: goTree,
				pkg: badPackage, give: "Tint", atKind: symbol.KindFunction, at: "ParseTint",
				wantMsg: "ParseTint has the name of the parse function of Tint and another form, so Tint has no " +
					"surface. The parse function is a function of a text that returns a variant and an error",
			},
			{
				name:   "reports Conflict at a constant under the name of the error of the parse function",
				target: golang.Target, tree: goTree, pkg: badPackage, give: "Hue", atKind: symbol.KindConstant,
				at: "ErrUnknownHue",
				wantMsg: "ErrUnknownHue has the name of the error of the parse function of Hue and another form, " +
					"so Hue has no surface. The error of the parse function is a variable",
			},
			{
				name: "reports Conflict at a render with a parameter", target: golang.Target, tree: goTree,
				pkg: badPackage, give: "Beam", atKind: symbol.KindMethod, atOwner: "Beam", at: "String",
				wantMsg: "String has the name of the render of Beam and another form, so Beam has no surface. " +
					"The render is a method of the enumeration without a parameter that returns a text",
			},
			{
				name: "reports Conflict at a parse function that returns another type", target: golang.Target,
				tree: goTree, pkg: badPackage, give: "Knot", atKind: symbol.KindFunction, at: "ParseKnot",
				wantMsg: "ParseKnot has the name of the parse function of Knot and another form, so Knot has no " +
					"surface. The parse function is a function of a text that returns a variant and an error",
			},
			{
				name:   "reports Conflict at a type under the name of the error of the parse function",
				target: golang.Target, tree: goTree, pkg: badPackage, give: "Knot", atKind: symbol.KindStruct,
				at: "ErrUnknownKnot",
				wantMsg: "ErrUnknownKnot has the name of the error of the parse function of Knot and another form, " +
					"so Knot has no surface. The error of the parse function is a variable",
			},
			{
				name:   "reports Conflict at a list of texts under the name of the list of the variants",
				target: golang.Target, tree: goTree, pkg: badPackage, give: "Knot", atKind: symbol.KindFunction,
				at: "KnotValues",
				wantMsg: "KnotValues has the name of the list of the variants of Knot and another form, so Knot has " +
					"no surface. The list of the variants is a function without a parameter that returns a list of " +
					"the enumeration",
			},
			{
				name: "reports Conflict at an encode that returns a text", target: golang.Target, tree: goTree,
				pkg: badPackage, give: "Knot", atKind: symbol.KindMethod, atOwner: "Knot", at: "MarshalText",
				wantMsg: "MarshalText has the name of the encode of Knot and another form, so Knot has no surface. " +
					"The encode is a method of the enumeration without a parameter that returns bytes and an error",
			},
			{
				name: "reports Conflict at a parse function without an error", target: golang.Target, tree: goTree,
				pkg: badPackage, give: "Loop", atKind: symbol.KindFunction, at: "ParseLoop",
				wantMsg: "ParseLoop has the name of the parse function of Loop and another form, so Loop has no " +
					"surface. The parse function is a function of a text that returns a variant and an error",
			},
			{
				name:   "reports Conflict at a function under the name of the error of the parse function",
				target: golang.Target, tree: goTree, pkg: badPackage, give: "Loop", atKind: symbol.KindFunction,
				at: "ErrUnknownLoop",
				wantMsg: "ErrUnknownLoop has the name of the error of the parse function of Loop and another form, " +
					"so Loop has no surface. The error of the parse function is a variable",
			},
			{
				name:   "reports Conflict at an array under the name of the list of the variants",
				target: golang.Target, tree: goTree, pkg: badPackage, give: "Loop", atKind: symbol.KindFunction,
				at: "LoopValues",
				wantMsg: "LoopValues has the name of the list of the variants of Loop and another form, so Loop has " +
					"no surface. The list of the variants is a function without a parameter that returns a list of " +
					"the enumeration",
			},
			{
				name: "reports Conflict at a parse function that takes a context", target: golang.Target,
				tree: goTree, pkg: badPackage, give: "Mesh", atKind: symbol.KindFunction, at: "ParseMesh",
				wantMsg: "ParseMesh has the name of the parse function of Mesh and another form, so Mesh has no " +
					"surface. The parse function is a function of a text that returns a variant and an error",
			},
			{
				name:   "reports Conflict at a variable under the name of the list of the variants",
				target: golang.Target, tree: goTree, pkg: badPackage, give: "Mesh", atKind: symbol.KindVariable,
				at: "MeshValues",
				wantMsg: "MeshValues has the name of the list of the variants of Mesh and another form, so Mesh has " +
					"no surface. The list of the variants is a function without a parameter that returns a list of " +
					"the enumeration",
			},
			{
				name: "reports Conflict at a validity method that returns a text", target: golang.Target,
				tree: goTree, pkg: badPackage, give: "Glow", atKind: symbol.KindMethod, atOwner: "Glow", at: "IsValid",
				wantMsg: "IsValid has the name of the validity method of Glow and another form, so Glow has no " +
					"surface. The validity method is a method of the enumeration without a parameter that returns " +
					"a truth value",
			},
			{
				name: "reports Conflict at a TypeScript parse function of a number", target: typescript.Target,
				tree: tsTree, pkg: badModule, give: "Tint", atKind: symbol.KindFunction, at: "parseTint",
				wantMsg: "parseTint has the name of the parse function of Tint and another form, so Tint has no " +
					"surface. The parse function is a function of a text that returns a variant",
			},
			{
				name:   "reports Conflict at a TypeScript function under the name of the error of the parse function",
				target: typescript.Target, tree: tsTree, pkg: badModule, give: "Hue", atKind: symbol.KindFunction,
				at: "UnknownHueError",
				wantMsg: "UnknownHueError has the name of the error of the parse function of Hue and another form, " +
					"so Hue has no surface. The error of the parse function is a class",
			},
		}
		for _, tt := range conflicts {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p := enum.New()
				f := fixture(t, tt.target, tt.tree, p, tt.pkg)
				res := f.Generate(t, p)
				assert.NoError(t, res.Err, "the generator runs whole")
				got := gentest.Finding(t, res, f.Find(t, tt.atKind, tt.atOwner, tt.at))
				expect.Equal(t, got.Code, enum.Conflict, "the finding has the code Conflict")
				expect.Equal(t, got.Msg, tt.wantMsg, "the message contains the declaration, the member and its form")
				gentest.Unrecorded(t, res, f.Find(t, symbol.KindEnum, "", tt.give))
			})
		}
	})
}
