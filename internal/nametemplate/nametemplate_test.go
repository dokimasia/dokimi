// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package nametemplate_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/internal/nametemplate"
)

// Segments returns the texts around the actions of a name template, and
// refuses each template that breaks a rule with the rule in its message.
func TestNameTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Segments", func(t *testing.T) {
		t.Parallel()

		accepted := []struct {
			name      string
			give      string
			giveField string
			want      []string
		}{
			{
				name: "returns the text before the field", give: "On{{.Method}}", giveField: "Method",
				want: []string{"On", ""},
			},
			{
				name: "returns the text after the field", give: "{{.Method}}Stub", giveField: "Method",
				want: []string{"", "Stub"},
			},
			{
				name: "returns each text between two uses of the field", give: "a{{.Enum}}b{{.Enum}}",
				giveField: "Enum", want: []string{"a", "b", ""},
			},
			{
				name: "returns one text for a template without an action", give: "is valid", giveField: "Enum",
				want: []string{"is valid"},
			},
			{
				name: "accepts each separator between words", give: "parse {{.Enum}}-from_text", giveField: "Enum",
				want: []string{"parse ", "-from_text"},
			},
			{
				name: "accepts spaces inside the braces of the action", give: "Parse{{ .Enum }}", giveField: "Enum",
				want: []string{"Parse", ""},
			},
			{
				name: "accepts a digit after the first letter", give: "V2{{.Enum}}", giveField: "Enum",
				want: []string{"V2", ""},
			},
		}
		for _, tt := range accepted {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, err := nametemplate.Segments(tt.give, tt.giveField)
				assert.NoError(t, err, "the template follows each rule")
				expect.Equal(t, got, tt.want, "the texts around the actions")
			})
		}

		refused := []struct {
			name    string
			give    string
			wantErr string
		}{
			{
				name: "refuses a field other than the field of the generator", give: "On{{.Name}}",
				wantErr: `the name template "On{{.Name}}" has an action other than {{.Method}}`,
			},
			{
				name: "refuses a field of the field of the generator", give: "On{{.Method.Name}}",
				wantErr: `the name template "On{{.Method.Name}}" has an action other than {{.Method}}`,
			},
			{
				name: "refuses a constant in place of the field", give: `On{{"Get"}}`,
				wantErr: `the name template "On{{\"Get\"}}" has an action other than {{.Method}}`,
			},
			{
				name: "refuses an argument after the field", give: `On{{.Method "x"}}`,
				wantErr: `the name template "On{{.Method \"x\"}}" has an action other than {{.Method}}`,
			},
			{
				name: "refuses a pipeline of the field", give: "On{{.Method | len}}",
				wantErr: `the name template "On{{.Method | len}}" has an action other than {{.Method}}`,
			},
			{
				name: "refuses the declaration of a variable", give: "On{{$name := .Method}}{{.Method}}",
				wantErr: `the name template "On{{$name := .Method}}{{.Method}}" has an action other than {{.Method}}`,
			},
			{
				name: "refuses a condition", give: "{{if .Method}}On{{end}}{{.Method}}",
				wantErr: `the name template "{{if .Method}}On{{end}}{{.Method}}" has an action other than {{.Method}}`,
			},
			{
				name: "refuses a mark in the text", give: "On!{{.Method}}",
				wantErr: `the name template "On!{{.Method}}" contains the character '!'`,
			},
			{
				name: "refuses a mark at the start of the text after the field", give: "On{{.Method}}!",
				wantErr: `the name template "On{{.Method}}!" contains the character '!'`,
			},
			{
				name: "refuses a template that starts with a digit", give: "2{{.Method}}",
				wantErr: `the name template "2{{.Method}}" starts with '2' instead of a letter`,
			},
			{
				name: "refuses a template that starts with a separator", give: "-{{.Method}}",
				wantErr: `the name template "-{{.Method}}" starts with '-' instead of a letter`,
			},
		}
		for _, tt := range refused {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, err := nametemplate.Segments(tt.give, "Method")
				assert.HasError(t, err, "the template breaks a rule")
				expect.Equal(t, err.Error(), tt.wantErr, "the error contains the template and the rule")
				expect.Nil(t, got, "a refused template has no texts")
			})
		}

		t.Run("refuses the empty template", func(t *testing.T) {
			t.Parallel()

			_, err := nametemplate.Segments("", "Method")
			assert.HasError(t, err, "an empty template has no name")
			expect.Equal(t, err.Error(), `the name template "" is empty`,
				"the error contains the template and the rule")
		})

		t.Run("refuses a template that does not parse", func(t *testing.T) {
			t.Parallel()

			_, err := nametemplate.Segments("On{{.Method", "Method")
			assert.HasError(t, err, "the text is refused")
			expect.Contains(t, err.Error(), `the name template "On{{.Method" does not parse: template: name:1:`,
				"the error contains the template and the error of the parser")
		})
	})
}
