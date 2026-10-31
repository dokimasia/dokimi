// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum_test

import (
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	"go.dokimi.dev/eidos/lang/typescript"
	tsbackend "go.dokimi.dev/eidos/lang/typescript/backend"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/generator/enum"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The part of TypeScript declares the target typescript with the template
// of each family and its helper, and names the functions of the surface of an
// enum. It refers to each member through its enum, and quotes the text of each
// member. The checks import each function of the surface by its module.
func TestTypeScript(t *testing.T) {
	t.Parallel()

	t.Run("declares the template of each family for the target typescript", func(t *testing.T) {
		t.Parallel()

		gentest.Templates(t, enum.New(), typescript.Target, enumTemplate, testTemplate)
	})

	t.Run("writes the golden tree typescript", func(t *testing.T) {
		t.Parallel()

		golden.MatchTree(t, "typescript", renderedTypeScript(t), golden.ShouldUpdate())
	})

	t.Run("names each declaration of the surface", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, statusModule).models[statusName].Names
		assert.Equal(t, got, enum.Names{
			Render: "statusToString", Parse: "parseStatus", Values: "statusValues", Valid: "isStatus",
			Sentinel: "UnknownStatusError",
		}, "the surface of Status has the names of TypeScript, and no pair of text")
	})

	templates := []struct {
		name string
		give func(*enum.Options)
		want enum.Names
	}{
		{
			name: "names the parse function in camelCase after the option parse-function",
			give: func(o *enum.Options) { o.ParseFunction = "Decode{{.Enum}}" },
			want: enum.Names{
				Render: "statusToString", Parse: "decodeStatus", Values: "statusValues", Valid: "isStatus",
				Sentinel: "UnknownStatusError",
			},
		},
		{
			name: "names the function that lists the variants in camelCase after the option values-function",
			give: func(o *enum.Options) { o.ValuesFunction = "every {{.Enum}}" },
			want: enum.Names{
				Render: "statusToString", Parse: "parseStatus", Values: "everyStatus", Valid: "isStatus",
				Sentinel: "UnknownStatusError",
			},
		},
		{
			name: "names the error class with the suffix Error after the option parse-error",
			give: func(o *enum.Options) { o.ParseError = "unrecognised {{.Enum}}" },
			want: enum.Names{
				Render: "statusToString", Parse: "parseStatus", Values: "statusValues", Valid: "isStatus",
				Sentinel: "UnrecognisedStatusError",
			},
		},
		{
			name: "names the type guard after the enum under the option valid-method",
			give: func(o *enum.Options) { o.ValidMethod = "Valid" },
			want: enum.Names{
				Render: "statusToString", Parse: "parseStatus", Values: "statusValues", Valid: "isStatus",
				Sentinel: "UnknownStatusError",
			},
		},
	}
	for _, tt := range templates {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := enum.New()
			tt.give(p.(plugin.OptionsProvider).Options().(*enum.Options))
			got := run(t, p, typescript.Target, tsTree, statusModule).models[statusName].Names
			assert.Equal(t, got, tt.want, "the surface of Status has the names that the options render")
		})
	}

	t.Run("refers to a member whose name is no identifier through brackets", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, statusModule).models[headerName]
		members := make([]string, 0, len(got.Variants))
		for _, v := range got.Variants {
			members = append(members, v.Member)
		}
		assert.Equal(t, members, []string{`["content-type"]`, ".Accept", ".Html", `["1st"]`},
			"the reference to a member whose name is no identifier is in brackets")
		want := []string{`content-type="ct"`, `Accept="a"`, `Html="<html>"`, `1st="first"`}
		expect.Equal(t, texts(got.Variants), want, "each text is a string of TypeScript, without escapes of HTML")
	})

	t.Run("prints a value outside the members without a conversion", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, statusModule).models[statusName]
		expect.Nil(t, got.Fallback, "the render prints a value outside the members through String")
	})

	t.Run("declares nothing beside an enum", func(t *testing.T) {
		t.Parallel()

		got := generate(t, typescript.Target, tsTree, statusModule)
		expect.Empty(t, got.beside[statusName], "the checks import each function of the surface by its module")
		expect.Empty(t, gentest.Referents(got.tests[statusName]), "the checks of Status refer to no field")
	})

	modules := []struct {
		name   string
		helper string
		want   string
	}{
		{name: "names the module of the surface beside a spec", helper: "surfaceModule", want: "./status.enum"},
		{name: "names the module of the enum beside a spec", helper: "sourceModule", want: "./status"},
	}
	for _, tt := range modules {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			templates, provides := enum.New().(plugin.TemplateProvider)
			assert.True(t, provides, "enum declares templates")
			module, declared := templates.TemplateFuncs(typescript.Target)[tt.helper].(func(string) string)
			assert.True(t, declared, "the target typescript declares the helper "+tt.helper)
			expect.Equal(t, module("status.enum.test.ts"), tt.want, "the test file imports the module")
		})
	}
}

// renderedTypeScript runs enum over the module status of the TypeScript
// tree, and returns the files that the TypeScript backend renders from it.
func renderedTypeScript(t *testing.T) fstest.MapFS {
	t.Helper()

	p := enum.New()
	f := fixture(t, typescript.Target, tsTree, p, statusModule)
	res := f.Generate(t, p)
	assert.NoError(t, res.Err, "the generator runs whole")
	return gentest.Rendered(t, f, res, p, tsbackend.New())
}
