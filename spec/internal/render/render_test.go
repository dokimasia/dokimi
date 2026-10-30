// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package render_test

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/files"

	"go.dokimi.dev/dokimi/spec/internal/render"
)

// minimal is a spec with one entry of each kind, in the form that people
// edit.
func minimal() files.Tree {
	return files.Tree{
		"VERSION": files.Text("1.0.0\n"),
		"laws.yaml": files.Text(`laws:
  no-value-beside-an-error:
    statement: A call that fails returns no value.
    depth: nothing-beyond-a-caller
    distinct-from: []
    since: 1.0.0
`),
		"claims.yaml": files.Text(`claims:
  zero-on-error:
    generator: suite
    licences:
      - shape: reader
    unless: default-on-error
    gates: [failure, returns-beside-error]
    law: no-value-beside-an-error
    machineries: [nothing-beyond-a-caller]
    statement: "{callable} returns the zero value beside an error"
    probe: zero-beside-error
    partners:
      write: {mode: named, param: write}
    inputs:
      call: {from: callable}
      zero: {from: zero}
    defect:
      form: value-beside
      clause: "a {declaration} whose {callable} returns a value beside an error"
      reason: value-beside-error
    since: 1.0.0
`),
		"gates.yaml": files.Text(`gates:
  failure:
    passes: The error model of the callable is not none.
    why: "{callable} cannot fail"
    remedy: "declare an error result"
    since: 1.0.0
  returns-beside-error:
    passes: The error model of the callable is last-return.
    why: "{callable} returns no value when it fails"
    remedy: "none, because the language returns no value beside an error"
    since: 1.0.0
`),
		"probes.yaml": files.Text(`probes:
  zero-beside-error:
    summary: A call that fails returns the zero value of its result.
    inputs:
      - {name: call, kind: call}
      - {name: zero, kind: value}
    reasons: [value-beside-error]
    vacuous: No call fails.
    since: 1.0.0
`),
		"reasons.yaml": files.Text(`reasons:
  value-beside-error:
    meaning: The call failed and returned a value <x> other than the zero value.
    since: 1.0.0
`),
		"defects.yaml": files.Text(`forms:
  value-beside:
    changes: callable
    does: A failed call returns the second sample of the value instead of the zero value.
    min-steps: 1
    since: 1.0.0
`),
		"machineries.yaml": files.Text(`machineries:
  nothing-beyond-a-caller:
    supplies: A fixed sequence of calls on one subject, with derived inputs.
    by: runtime
    since: 1.0.0
`),
		"stances.yaml": files.Text(`stances:
  silent:
    meaning: The classification implies nothing that a check could falsify.
    closes: Nothing.
    since: 1.0.0
unclaimed:
  mixin:deprecated:
    stance: silent
    note: The mixin documents the callable.
`),
	}
}

func TestRender(t *testing.T) {
	t.Parallel()

	t.Run("Files", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the files of the spec in the order of the render", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, render.Files(), []string{
				"laws", "claims", "gates", "probes", "reasons", "defects", "machineries", "stances",
			}, "the spec has eight files")
		})
	})

	t.Run("Render", func(t *testing.T) {
		t.Parallel()

		t.Run("writes each section as an array of entries in the order of the file", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			dir := files.Workspace(t, tree)

			assert.NoError(t, render.Render(dir), "the spec renders")

			var got struct {
				Gates []struct {
					ID string `json:"id"`
				} `json:"gates"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "gates.json"))), &got),
				"gates.json decodes",
			)
			assert.Equal(
				t,
				[]string{got.Gates[0].ID, got.Gates[1].ID},
				[]string{"failure", "returns-beside-error"},
				"the gates keep the order of gates.yaml",
			)
		})

		t.Run("writes the inputs of a claim as a list with the name of each input", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the spec renders")

			var got struct {
				Claims []struct {
					Partners []map[string]string `json:"partners"`
					Inputs   []map[string]string `json:"inputs"`
				} `json:"claims"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "claims.json"))), &got),
				"claims.json decodes",
			)
			assert.Equal(t, got.Claims[0].Inputs, []map[string]string{
				{"name": "call", "from": "callable"},
				{"name": "zero", "from": "zero"},
			}, "the inputs keep the order of claims.yaml")
			assert.Equal(
				t,
				got.Claims[0].Partners,
				[]map[string]string{{"name": "write", "mode": "named", "param": "write"}},
				"the partners are a list with the name of each",
			)
		})

		t.Run("writes a mixin under unless as an object with the mixin", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the spec renders")

			var got struct {
				Claims []struct {
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
				got.Claims[0].Unless,
				map[string]string{"mixin": "default-on-error"},
				"the mixin is under the key mixin",
			)
		})

		t.Run("returns an error for a directory that does not exist", func(t *testing.T) {
			t.Parallel()

			err := render.Render(filepath.Join(t.TempDir(), "missing"))

			assert.ErrorIs(t, err, fs.ErrNotExist, "the missing directory fails the render")
		})

		t.Run("returns an error for names.gen.go that it cannot write", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["names.gen.go/keep"] = files.Text("a directory in place of the file\n")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.HasError(t, err, "the directory in place of names.gen.go fails the render")
			assert.Contains(t, err.Error(), "names.gen.go:", "the error names the file")
		})

		t.Run("keeps the inputs of a probe that are a list already", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the spec renders")

			var got struct {
				Probes []struct {
					Inputs []map[string]string `json:"inputs"`
				} `json:"probes"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "probes.json"))), &got),
				"probes.json decodes",
			)
			assert.Equal(t, got.Probes[0].Inputs, []map[string]string{
				{"name": "call", "kind": "call"},
				{"name": "zero", "kind": "value"},
			}, "the inputs keep the list of probes.yaml")
		})

		t.Run("writes indented JSON without escaping HTML", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the spec renders")

			got := files.Read(t, filepath.Join(dir, "reasons.json"))
			assert.HasPrefix(
				t,
				got,
				"{\n  \"reasons\": [\n    {\n      \"id\": \"value-beside-error\",",
				"the file is indented by two spaces",
			)
			assert.Contains(t, got, "a value <x> other", "the file keeps the angle brackets of the meaning")
		})

		t.Run(
			"writes a schema of draft 2020-12 that refuses an identifier with upper-case letters",
			func(t *testing.T) {
				t.Parallel()
				dir := files.Workspace(t, minimal())

				assert.NoError(t, render.Render(dir), "the spec renders")

				var schema jsonschema.Schema
				assert.NoError(
					t,
					json.Unmarshal(
						[]byte(files.Read(t, filepath.Join(dir, "schemas", "reasons.schema.json"))),
						&schema,
					),
					"the schema decodes",
				)
				assert.Equal(
					t,
					schema.Schema,
					"https://json-schema.org/draft/2020-12/schema",
					"the schema states its draft",
				)
				resolved, err := schema.Resolve(nil)
				assert.NoError(t, err, "the schema resolves")
				doc := map[string]any{
					"reasons": []any{map[string]any{"id": "Value_Beside", "meaning": "m", "since": "1.0.0"}},
				}
				assert.HasError(t, resolved.Validate(doc), "the schema refuses the identifier")
			},
		)

		t.Run("writes a schema that accepts its rendered file", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the spec renders")

			for _, name := range render.Files() {
				var schema jsonschema.Schema
				assert.NoError(
					t,
					json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "schemas", name+".schema.json"))), &schema),
					"the schema decodes",
				)
				resolved, err := schema.Resolve(nil)
				assert.NoError(t, err, "the schema resolves")
				var doc any
				assert.NoError(
					t,
					json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, name+".json"))), &doc),
					"the rendered file decodes",
				)
				assert.NoError(t, resolved.Validate(doc), "the schema accepts the rendered file")
			}
		})

		t.Run("writes the schema of the suite record", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the spec renders")

			var schema jsonschema.Schema
			assert.NoError(t, json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "schemas", "suite.schema.json"))),
				&schema), "the schema decodes")
			assert.Equal(t, schema.Schema, "https://json-schema.org/draft/2020-12/schema",
				"the schema has the draft of the other schemas")
			resolved, err := schema.Resolve(nil)
			assert.NoError(t, err, "the schema resolves")
			record := func(seed, result string) map[string]any {
				return map[string]any{
					"spec": "1.0.0",
					"suite": map[string]any{
						"language":    "go",
						"package":     "example.com/store",
						"declaration": "Store",
					},
					"generator":   map[string]any{"name": "dokimi", "version": "v0.1.0"},
					"runtime":     map[string]any{"name": "go.dokimi.dev/conform", "version": "v0.1.0"},
					"seed":        seed,
					"application": "subject",
					"checks": []any{map[string]any{
						"identity": "Store.Read/zero-on-error", "key": "store.read/zero-on-error", "result": result,
					}},
					"refusals":  []any{},
					"unclaimed": []any{},
				}
			}
			assert.NoError(t, resolved.Validate(record("0x5f3a1c9e27d4b860", "pass")), "the schema accepts a record")
			later := record("0x5f3a1c9e27d4b860", "pass")
			later["trace"] = "a field of a later minor version"
			later["checks"].([]any)[0].(map[string]any)["steps"] = 3
			assert.NoError(t, resolved.Validate(later), "the schema accepts the fields of a later minor version")
			assert.HasError(t, resolved.Validate(record("0x5f3a", "pass")), "the schema refuses a short seed")
			assert.HasError(t, resolved.Validate(record("0x5f3a1c9e27d4b860", "green")),
				"the schema refuses a result outside the spec")
		})

		t.Run("writes a constant for each entry in the order of the files", func(t *testing.T) {
			t.Parallel()
			dir := files.Workspace(t, minimal())

			assert.NoError(t, render.Render(dir), "the spec renders")

			got := files.Read(t, filepath.Join(dir, "names.gen.go"))
			assert.Matches(t, got, `ZeroOnError +ClaimID += "zero-on-error"`, "names.gen.go declares the claim")
			assert.Matches(
				t,
				got,
				`ReturnsBesideError +GateID += "returns-beside-error"`,
				"names.gen.go declares the gate",
			)
			assert.ContainsInOrder(t, got, []string{
				"NoValueBesideAnError", "ZeroOnError", "Failure", "ReturnsBesideError", "ZeroBesideError",
				"ValueBesideError", "ValueBeside", "NothingBeyondACaller", "Silent",
			}, "names.gen.go declares the entries in the order of the files")
		})

		t.Run("writes the spelling of each identifier in each language by its section", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["stances.yaml"] = files.Text(`stances:
  ttl-gap:
    meaning: The classification has a law and no claim.
    closes: A claim.
    since: 1.0.0
unclaimed: {}
`)
			dir := files.Workspace(t, tree)

			assert.NoError(t, render.Render(dir), "the spec renders")

			var got struct {
				Names     map[string]map[string]map[string]string `json:"names"`
				Languages []string                                `json:"languages"`
			}
			assert.NoError(
				t,
				json.Unmarshal([]byte(files.Read(t, filepath.Join(dir, "naming.json"))), &got),
				"naming.json decodes",
			)
			assert.Equal(t, got.Languages, []string{"go", "typescript"}, "the table names Go and TypeScript")
			assert.Equal(t, got.Names["claims"]["zero-on-error"], map[string]string{
				"go": "ZeroOnError", "typescript": "zeroOnError",
			}, "the claim has its spelling in each language")
			assert.Equal(t, got.Names["stances"]["ttl-gap"], map[string]string{
				"go": "TTLGap", "typescript": "ttlGap",
			}, "Go spells the initialism in capitals and TypeScript does not")
			assert.Length(t, got.Names, 8, "each section with constants has its names")
		})

		t.Run("returns an error for naming.json that it cannot write", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["naming.json/keep"] = files.Text("a directory in place of the file\n")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.HasError(t, err, "the directory in place of naming.json fails the render")
			assert.Contains(t, err.Error(), "naming.json", "the error names the file")
		})

		t.Run("returns an error for a field that the entry does not have", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["reasons.yaml"] = files.Text("reasons:\n  value-beside-error:\n    meanin: typo\n    since: 1.0.0\n")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.HasError(t, err, "the misspelled field fails the render")
			assert.Contains(
				t,
				err.Error(),
				`reasons.yaml: json: unknown field "meanin"`,
				"the error names the file and the field",
			)
		})

		tests := []struct {
			name     string
			giveFile string
			giveText string
			wantErr  string
		}{
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
				name:     "returns an error for a nested value that does not decode",
				giveFile: "laws.yaml",
				giveText: "laws:\n  a:\n    distinct-from: [{law: !!int x}]\n",
				wantErr:  "laws.yaml: section laws: a: distinct-from: line 3:",
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
			{
				name:     "returns an error for an identifier that starts with an article",
				giveFile: "laws.yaml",
				giveText: "laws:\n  a-read-returns:\n    statement: s\n",
				wantErr:  "laws.yaml: section laws: a-read-returns: the identifier starts with the article a",
			},
			{
				name: "returns an error for an identifier that is no Go name", giveFile: "laws.yaml",
				giveText: "laws:\n  9a:\n    statement: s\n", wantErr: "names.gen.go:",
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
			delete(tree, "stances.yaml")
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

		t.Run("returns an error for the schema of a record that it cannot write", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["schemas/suite.schema.json/keep"] = files.Text("a directory in place of the schema\n")
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.HasError(t, err, "the directory in place of the schema fails the render")
			assert.Contains(t, err.Error(), "the schema of the suite record:", "the error names the record")
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

		t.Run("returns an error for two entries with one identifier", func(t *testing.T) {
			t.Parallel()
			tree := minimal()
			tree["reasons.yaml"] = files.Text(
				"reasons:\n  failure:\n    meaning: A reason with the name of a gate.\n    since: 1.0.0\n",
			)
			dir := files.Workspace(t, tree)

			err := render.Render(dir)

			assert.HasError(t, err, "the shared identifier fails the render")
			assert.Contains(
				t,
				err.Error(),
				"the identifier failure names a gate and a reason",
				"the error names the identifier and both kinds",
			)
		})
	})
}
