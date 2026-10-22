// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum_test

import (
	"encoding/json"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/enum"
)

// The decoder of a config file takes a name template or a method template,
// and refuses any other text with the rule that the text breaks.
func TestNameTemplate(t *testing.T) {
	t.Parallel()

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		accepted := []struct {
			name string
			give string
		}{
			{name: "accepts a name before the name of the enumeration", give: "Parse{{.Enum}}"},
			{name: "accepts a name after the name of the enumeration", give: "{{.Enum}}Values"},
			{name: "accepts a space, a hyphen or an underscore between words", give: "parse {{.Enum}}-from_text"},
			{name: "accepts spaces inside the braces of the action", give: "Parse{{ .Enum }}"},
			{name: "accepts a digit after the first letter", give: "V2{{.Enum}}"},
		}
		for _, tt := range accepted {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				quoted, err := json.Marshal(tt.give)
				assert.NoError(t, err, "the text encodes as JSON")
				var got enum.NameTemplate
				assert.NoError(t, json.Unmarshal(quoted, &got), "the decoder of a config file takes the template")
				expect.Equal(t, got, enum.NameTemplate(tt.give), "the template keeps its text")
			})
		}

		refused := []struct {
			name    string
			give    string
			wantErr string
		}{
			{
				name: "refuses a template without the name of the enumeration", give: "Parse",
				wantErr: `enum: the name template "Parse" does not use {{.Enum}}`,
			},
			{name: "refuses the empty template", give: "", wantErr: "enum: the name template is empty"},
			{
				name: "refuses a field other than Enum", give: "Parse{{.Name}}",
				wantErr: `enum: the name template "Parse{{.Name}}" has an action other than {{.Enum}}`,
			},
			{
				name: "refuses a field of the name of the enumeration", give: "Parse{{.Enum.Name}}",
				wantErr: `enum: the name template "Parse{{.Enum.Name}}" has an action other than {{.Enum}}`,
			},
			{
				name: "refuses a constant in place of the field", give: `Parse{{"Status"}}`,
				wantErr: `enum: the name template "Parse{{\"Status\"}}" has an action other than {{.Enum}}`,
			},
			{
				name: "refuses an argument after the field", give: `Parse{{.Enum "x"}}`,
				wantErr: `enum: the name template "Parse{{.Enum \"x\"}}" has an action other than {{.Enum}}`,
			},
			{
				name: "refuses a pipeline of the field", give: "Parse{{.Enum | len}}",
				wantErr: `enum: the name template "Parse{{.Enum | len}}" has an action other than {{.Enum}}`,
			},
			{
				name: "refuses the declaration of a variable", give: "Parse{{$name := .Enum}}{{.Enum}}",
				wantErr: `enum: the name template "Parse{{$name := .Enum}}{{.Enum}}" has an action other than ` +
					"{{.Enum}}",
			},
			{
				name: "refuses a condition", give: "{{if .Enum}}Parse{{end}}{{.Enum}}",
				wantErr: `enum: the name template "{{if .Enum}}Parse{{end}}{{.Enum}}" has an action other than ` +
					"{{.Enum}}",
			},
			{
				name: "refuses a mark in the text", give: "Parse!{{.Enum}}",
				wantErr: `enum: the name template "Parse!{{.Enum}}" contains the character '!'`,
			},
			{
				name: "refuses a template that starts with a digit", give: "2{{.Enum}}",
				wantErr: `enum: the name template "2{{.Enum}}" starts with '2' instead of a letter`,
			},
			{
				name: "refuses a template that starts with a separator", give: "-{{.Enum}}",
				wantErr: `enum: the name template "-{{.Enum}}" starts with '-' instead of a letter`,
			},
		}
		for _, tt := range refused {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := enum.DefaultParseFunction
				err := got.UnmarshalText([]byte(tt.give))
				assert.HasError(t, err, "the text is refused")
				expect.Equal(t, err.Error(), tt.wantErr, "the error contains the template and the rule")
				expect.Equal(t, got, enum.DefaultParseFunction, "a refused text leaves the option unchanged")
			})
		}

		t.Run("refuses a template that does not parse", func(t *testing.T) {
			t.Parallel()

			var got enum.NameTemplate
			err := got.UnmarshalText([]byte("Parse{{.Enum"))
			assert.HasError(t, err, "the text is refused")
			expect.Contains(t, err.Error(), `enum: the name template "Parse{{.Enum" does not parse: template: name:1:`,
				"the error contains the template and the error of the parser")
		})
	})

	t.Run("MethodTemplate", func(t *testing.T) {
		t.Parallel()

		t.Run("UnmarshalText", func(t *testing.T) {
			t.Parallel()

			accepted := []struct {
				name string
				give string
			}{
				{name: "accepts a name without the name of the enumeration", give: "Valid"},
				{name: "accepts words that a space separates", give: "is valid"},
				{name: "accepts the name of the enumeration", give: "IsA{{.Enum}}"},
			}
			for _, tt := range accepted {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					quoted, err := json.Marshal(tt.give)
					assert.NoError(t, err, "the text encodes as JSON")
					var got enum.MethodTemplate
					assert.NoError(t, json.Unmarshal(quoted, &got), "the decoder of a config file takes the template")
					expect.Equal(t, got, enum.MethodTemplate(tt.give), "the template keeps its text")
				})
			}

			t.Run("refuses a text that is no name template", func(t *testing.T) {
				t.Parallel()

				got := enum.DefaultValidMethod
				err := got.UnmarshalText([]byte("Is{{.Name}}"))
				assert.HasError(t, err, "the text is refused")
				expect.Equal(t, err.Error(), `enum: the name template "Is{{.Name}}" has an action other than {{.Enum}}`,
					"the error contains the template and the rule")
				expect.Equal(t, got, enum.DefaultValidMethod, "a refused text leaves the option unchanged")
			})
		})
	})
}
