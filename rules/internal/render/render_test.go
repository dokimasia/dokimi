// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package render_test

import (
	"encoding/json"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"

	"go.dokimi.dev/dokimi/rules/internal/render"
)

// minimal is a set of rules with one entry of each kind, in the form that
// people edit.
func minimal() files.Tree {
	return files.Tree{
		"claims.yaml": files.Text(`claims:
  zero-on-error:
    generator: suite
    licences:
      - shape: reader
    unless: default-on-error
    gates: [failure, returns-beside-error]
    partners:
      write: {mode: named, param: write}
    inputs:
      call: {from: callable}
      zero: {from: zero}
    clause: "a {declaration} whose {callable} returns a value <x> beside an error"
`),
		"laws.yaml": files.Text(`laws:
  no-value-beside-an-error:
    depth: nothing-beyond-a-caller
    distinct-from: []
`),
		"defects.yaml": files.Text(`forms:
  value-beside:
    changes: callable
    min-steps: 1
`),
		"gates.yaml": files.Text(`gates:
  failure:
    passes: The error model of the callable is not none.
    why: "{callable} cannot fail"
    remedy: "declare an error result"
  returns-beside-error:
    passes: The error model of the callable is last-return.
    why: "{callable} returns no value when it fails"
    remedy: "none, because the language returns no value beside an error"
`),
		"unclaimed.yaml": files.Text(`unclaimed:
  mixin:deprecated:
    stance: silent
    note: The mixin documents the callable.
`),
		"references.yaml": files.Text(`references:
  persister:
    contract: persister
    state: map
    roles:
      writer: store
      reader: load
    since: 1.0.0
`),
	}
}

// schemaOf returns the resolved schema of the rendered file name in dir.
func schemaOf(t *testing.T, dir, name string) *jsonschema.Resolved {
	t.Helper()
	var schema jsonschema.Schema
	assert.NoError(
		t,
		json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "schemas", name+".schema.json"))), &schema),
		"the schema of "+name+" decodes",
	)
	assert.Equal(t, schema.Schema, "https://json-schema.org/draft/2020-12/schema",
		"the schema of "+name+" states its draft")
	resolved, err := schema.Resolve(nil)
	assert.NoError(t, err, "the schema of "+name+" resolves")
	return resolved
}

func TestRender(t *testing.T) {
	t.Parallel()

	t.Run("Files", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the files of the rules in the order of the render", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, render.Files(), []string{"claims", "laws", "defects", "gates", "unclaimed", "references"},
				"the rules have six files")
		})
	})

	t.Run("Render", func(t *testing.T) {
		t.Parallel()

		t.Run("writes each section as an object from the identifier of each entry to the entry", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the rules render")

			var got struct {
				Gates map[string]map[string]string `json:"gates"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "gates.json"))), &got),
				"gates.json decodes",
			)
			assert.Equal(t, slices.Sorted(maps.Keys(got.Gates)), []string{"failure", "returns-beside-error"},
				"the gates are the keys of the section")
			assert.Equal(t, got.Gates["failure"]["why"], "{callable} cannot fail", "the gate has its reason")
		})

		t.Run("writes the inputs and the partners of a claim as lists with the name of each", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the rules render")

			var got struct {
				Claims map[string]struct {
					Partners []map[string]string `json:"partners"`
					Inputs   []map[string]string `json:"inputs"`
				} `json:"claims"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "claims.json"))), &got),
				"claims.json decodes",
			)
			assert.Equal(t, got.Claims["zero-on-error"].Inputs, []map[string]string{
				{"name": "call", "from": "callable"},
				{"name": "zero", "from": "zero"},
			}, "the inputs keep the order of claims.yaml")
			assert.Equal(
				t,
				got.Claims["zero-on-error"].Partners,
				[]map[string]string{{"name": "write", "mode": "named", "param": "write"}},
				"the partners are a list with the name of each",
			)
		})

		t.Run("keeps the inputs of a claim that are a list already", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["claims.yaml"] = files.Text(`claims:
  smoke:
    generator: suite
    inputs: [{name: call, from: callable}]
    clause: "a {declaration} whose {callable} panics"
`)
			dir := files.Workspace(t, tree)

			assert.NoError(t, render.Render(dir), "the rules render")

			var got struct {
				Claims map[string]struct {
					Inputs []map[string]string `json:"inputs"`
				} `json:"claims"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "claims.json"))), &got),
				"claims.json decodes",
			)
			assert.Equal(t, got.Claims["smoke"].Inputs, []map[string]string{{"name": "call", "from": "callable"}},
				"the inputs keep the list of claims.yaml")
		})

		t.Run("writes a schema that refuses an identifier of the spec with upper-case letters", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the rules render")

			law := func(depth string) map[string]any {
				return map[string]any{"laws": map[string]any{"no-value-beside-an-error": map[string]any{
					"depth": depth, "distinct-from": []any{},
				}}}
			}
			resolved := schemaOf(t, dir, "laws")
			assert.NoError(t, resolved.Validate(law("nothing-beyond-a-caller")), "the schema accepts an identifier")
			assert.HasError(t, resolved.Validate(law("Nothing_Beyond")), "the schema refuses the identifier")
		})

		t.Run("writes a mixin under unless as an object with the mixin", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the rules render")

			var got struct {
				Claims map[string]struct {
					Unless map[string]string `json:"unless"`
				} `json:"claims"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "claims.json"))), &got),
				"claims.json decodes",
			)
			assert.Equal(
				t,
				got.Claims["zero-on-error"].Unless,
				map[string]string{"mixin": "default-on-error"},
				"the mixin is under the key mixin",
			)
		})

		t.Run("writes indented JSON without escaping HTML", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the rules render")

			got := files.Read(t, filepath.Join(dir, "claims.json"))
			assert.HasPrefix(t, got, "{\n  \"claims\": {\n    \"zero-on-error\": {\n      \"generator\": \"suite\",",
				"the file is indented by two spaces")
			assert.Contains(t, got, "a value <x> beside", "the file keeps the angle brackets of the clause")
		})

		t.Run("writes a schema that accepts its rendered file", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the rules render")

			for _, name := range render.Files() {
				var doc any
				assert.NoError(
					t,
					json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, name+".json"))), &doc),
					"the rendered file decodes",
				)
				assert.NoError(t, schemaOf(t, dir, name).Validate(doc), "the schema accepts "+name+".json")
			}
		})

		t.Run("writes a schema of the claims that refuses a field of the runtime contract", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the rules render")

			claim := func(field string) map[string]any {
				c := map[string]any{"generator": "suite", "inputs": []any{}, "clause": "c"}
				if field != "" {
					c[field] = "s"
				}
				return map[string]any{"claims": map[string]any{"smoke": c}}
			}
			resolved := schemaOf(t, dir, "claims")
			assert.NoError(t, resolved.Validate(claim("")), "the schema accepts a rule")
			assert.HasError(t, resolved.Validate(claim("statement")), "the schema refuses the statement of the spec")
		})

		t.Run(
			"writes a schema of the reference forms with the states and the effects of the rules",
			func(t *testing.T) {
				t.Parallel()
				dir := files.Workspace(t, minimal())

				assert.NoError(t, render.Render(dir), "the rules render")

				form := func(state, effect string) map[string]any {
					return map[string]any{"references": map[string]any{"persister": map[string]any{
						"contract": "persister",
						"state":    state,
						"roles":    map[string]any{"writer": effect},
						"since":    "1.0.0",
					}}}
				}
				resolved := schemaOf(t, dir, "references")
				assert.NoError(t, resolved.Validate(form("list", "append")),
					"the schema accepts a state and an effect of the rules")
				assert.HasError(
					t,
					resolved.Validate(form("tree", "store")),
					"the schema refuses a state outside the rules",
				)
				assert.HasError(t, resolved.Validate(form("map", "delete")),
					"the schema refuses an effect outside the rules")
			},
		)

		t.Run("returns an error for a directory that does not exist", func(t *testing.T) {
			t.Parallel()

			err := render.Render(filepath.Join(t.TempDir(), "missing"))

			assert.ErrorIs(t, err, fs.ErrNotExist, "the missing directory fails the render")
		})

		tests := []struct {
			name     string
			giveFile string
			giveText string
			wantErr  string
		}{
			{
				name: "returns an error for a field that the entry does not have", giveFile: "gates.yaml",
				giveText: "gates:\n  failure:\n    whyy: typo\n", wantErr: `gates.yaml: json: unknown field "whyy"`,
			},
			{
				name: "returns an error for a field of the runtime contract", giveFile: "gates.yaml",
				giveText: "gates:\n  failure:\n    since: 1.0.0\n", wantErr: `gates.yaml: json: unknown field "since"`,
			},
			{
				name: "returns an error for a section that the file does not have", giveFile: "gates.yaml",
				giveText: "gates: {}\nstances: {}\n", wantErr: `gates.yaml: json: unknown field "stances"`,
			},
			{
				name: "returns an error for a file that does not parse", giveFile: "laws.yaml",
				giveText: "laws: [\n", wantErr: "laws.yaml: yaml:",
			},
			{
				name: "returns an error for an empty file", giveFile: "laws.yaml",
				giveText: "", wantErr: "laws.yaml: the file is no mapping of sections",
			},
			{
				name: "returns an error for a file that is no mapping of sections", giveFile: "laws.yaml",
				giveText: "- a\n", wantErr: "laws.yaml: the file is no mapping of sections",
			},
			{
				name:     "returns an error for a section that is no mapping",
				giveFile: "laws.yaml",
				giveText: "laws: [a]\n",
				wantErr:  "laws.yaml: section laws: the section laws is no mapping of identifiers to entries",
			},
			{
				name: "returns an error for an entry that is no mapping", giveFile: "laws.yaml",
				giveText: "laws:\n  a: b\n", wantErr: "laws.yaml: section laws: a: the entry a is no mapping",
			},
			{
				name:     "returns an error for an input that is no mapping",
				giveFile: "claims.yaml",
				giveText: "claims:\n  a:\n    inputs:\n      call: callable\n",
				wantErr:  "claims.yaml: section claims: a: inputs: the entry call is no mapping",
			},
			{
				name:     "returns an error for an element of a list that does not decode",
				giveFile: "laws.yaml",
				giveText: "laws:\n  a:\n    distinct-from: [{law: !!int x}]\n",
				wantErr:  "laws.yaml: section laws: a: distinct-from: line 3:",
			},
			{
				name:     "returns an error for a value of a mapping that does not decode",
				giveFile: "references.yaml",
				giveText: "references:\n  a:\n    roles: {writer: !!int x}\n",
				wantErr:  "references.yaml: section references: a: roles: line 3:",
			},
			{
				name:     "returns an error for a number that JSON cannot state",
				giveFile: "laws.yaml",
				giveText: "laws:\n  a:\n    depth: .nan\n",
				wantErr:  "laws.yaml: section laws: a: depth: line 3: json: unsupported value: NaN",
			},
			{
				name: "returns an error for an unless that is neither a mixin nor a param", giveFile: "claims.yaml",
				giveText: "claims:\n  a:\n    unless: [x]\n", wantErr: "claims.yaml: json: cannot unmarshal array",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				tree := minimal()
				tree[tt.giveFile] = files.Text(tt.giveText)
				dir := files.Workspace(t, tree)

				err := render.Render(dir)

				assert.HasError(t, err, "the render fails")
				assert.Contains(t, err.Error(), tt.wantErr, "the error names the file and the fault")
			})
		}

		t.Run("returns an error for a file that is missing", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			delete(tree, "unclaimed.yaml")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.ErrorIs(t, err, fs.ErrNotExist, "the missing file fails the render")
		})

		t.Run("returns an error for a rendered file that it cannot write", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["laws.json/keep"] = files.Text("a directory in place of the file\n")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			pathErr := assert.ErrorAs[*fs.PathError](t, err, "the directory in place of laws.json fails the render")
			assert.Equal(t, pathErr.Path, "laws.json", "the error comes from the creation of laws.json")
		})

		t.Run("returns an error for a schema that it cannot write", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["schemas/laws.schema.json/keep"] = files.Text("a directory in place of the schema\n")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.HasError(t, err, "the directory in place of the schema fails the render")
			assert.Contains(t, err.Error(), "the schema of laws:", "the error names the schema")
		})

		t.Run("returns an error for a directory that it cannot write", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["schemas"] = files.Text("a file in place of the directory\n")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.HasError(t, err, "the file in place of the directory fails the render")
			assert.Contains(t, err.Error(), "mkdir", "the error comes from the creation of the directory")
		})
	})
}
