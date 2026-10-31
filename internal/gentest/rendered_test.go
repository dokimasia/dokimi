// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest_test

import (
	"io/fs"
	"strings"
	"testing"
	"text/template"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	golang "go.dokimi.dev/eidos/lang/go"
	gobackend "go.dokimi.dev/eidos/lang/go/backend"
	tsbackend "go.dokimi.dev/eidos/lang/typescript/backend"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/internal/gentest"
)

// The template of the primary family of the recorder that writes a comment
// for each record, and a template that writes no Go.
const (
	commentBody = "{{range .Records}}// {{.Data.Name}} is recorded.\n{{end}}"
	brokenBody  = "{{range .Records}}func {{end}}"
)

// The paths at which the render routes the two files of the recorder, in the
// directory of the source file of Address.
const (
	recordPath = "catalog/address_recorder.go"
	checkPath  = "catalog/address_recorder_test.go"
)

// tracerID is the name of a recorder that the plan runs before the recorder.
// It sorts after recorder by name, so the order of the plan differs from the
// order of the names.
const tracerID plugin.ID = "tracer"

// The helper of the recorder and a template that calls it, and a template
// that calls the helper wrapdocs of the shared vocabulary of eidos.
const (
	shoutHelper = "shout"
	shoutBody   = "{{range .Records}}// {{shout .Data.Name}} is recorded.\n{{end}}"
	wrapHelper  = "wrapdocs"
	wrapBody    = "{{range .Records}}// {{wrapdocs .Data.Name}}\n{{end}}"
)

// Rendered renders what a generator recorded through a backend, with the
// templates of the generator, and returns the files by their routed paths.
func TestRendered(t *testing.T) {
	t.Parallel()

	t.Run("Rendered", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each file that the render writes by its routed path", func(t *testing.T) {
			t.Parallel()

			p, f, res := generate(t, recording{records: []string{addressName}, body: commentBody})
			got := gentest.Rendered(t, f, res, p, gobackend.New())
			paths, err := fs.Glob(got, "*/*")
			assert.NoError(t, err, "the files of the render list")
			assert.Equal(t, paths, []string{recordPath, checkPath},
				"the render writes the file of each family of the recorder beside the source")
			expect.Contains(t, string(got[recordPath].Data), "// Address is recorded.",
				"the file contains what the template of the recorder writes")
		})

		t.Run("writes the file of a data family of one file per plan at the root", func(t *testing.T) {
			t.Parallel()

			p, f, res := generate(t, recording{listed: []string{addressName, warehouseName}})
			got := gentest.Rendered(t, f, res, p, gobackend.New())

			assert.Equal(t, string(got[listFile].Data), "Address\nWarehouse\n",
				"the data file of the plan has the record of each struct")
		})

		t.Run("returns the files of the plugin alone where the plan runs a generator before it", func(t *testing.T) {
			t.Parallel()

			tracer := recorder(recording{id: tracerID, records: []string{warehouseName}, body: commentBody})
			p := recorder(recording{records: []string{addressName}, body: commentBody})
			f := gentest.New(t, &gentest.Spec{
				Dir: goTree, Packages: []string{catalogPackage}, Target: golang.Target,
				Plugins: []plugin.Plugin{tracer, p},
			})
			assert.NoError(t, f.Generate(t, tracer).Err, "the tracer runs whole")
			res := f.Generate(t, p)
			assert.NoError(t, res.Err, "the recorder runs whole")

			got := gentest.Rendered(t, f, res, p, gobackend.New(), tracer)

			paths, err := fs.Glob(got, "*/*")
			assert.NoError(t, err, "the files of the render list")
			assert.Equal(t, paths, []string{recordPath, checkPath},
				"the render returns the files of the recorder alone")
		})

		t.Run("renders the helpers of the plugin", func(t *testing.T) {
			t.Parallel()

			p, f, res := generate(t, recording{
				records: []string{addressName}, body: shoutBody,
				helpers: template.FuncMap{shoutHelper: strings.ToUpper},
			})
			got := gentest.Rendered(t, f, res, p, gobackend.New())

			expect.Contains(t, string(got[recordPath].Data), "// ADDRESS is recorded.",
				"the template of the recorder calls the helper of the recorder")
		})

		t.Run("renders the override of the generator that the plan runs last", func(t *testing.T) {
			t.Parallel()

			tracer := recorder(recording{
				id: tracerID, overrides: template.FuncMap{wrapHelper: func(string) string { return string(tracerID) }},
			})
			p := recorder(recording{
				records: []string{addressName}, body: wrapBody,
				overrides: template.FuncMap{wrapHelper: func(string) string { return string(recorderID) }},
			})
			f := gentest.New(t, &gentest.Spec{
				Dir: goTree, Packages: []string{catalogPackage}, Target: golang.Target,
				Plugins: []plugin.Plugin{tracer, p},
			})
			res := f.Generate(t, p)
			assert.NoError(t, res.Err, "the recorder runs whole")

			got := gentest.Rendered(t, f, res, p, gobackend.New(), tracer)

			expect.Contains(t, string(got[recordPath].Data), "// recorder\n",
				"the override of the recorder replaces the override of the tracer, which the plan runs before it")
		})

		t.Run("fails where the plugin declares no templates", func(t *testing.T) {
			t.Parallel()

			_, f, res := generate(t, recording{records: []string{addressName}, body: commentBody})
			records := assert.Rejects(t, "a plugin without templates", func(tb assert.TB) {
				gentest.Rendered(tb, f, res, bare{}, gobackend.New())
			})
			assert.Length(t, records, 1, "the helper stops at the templates")
			expect.Equal(t, records[0].Contract, "the plugin declares templates", "the helper stops at the plugin")
		})

		t.Run("fails where the plugin declares no template tree for the target of the backend", func(t *testing.T) {
			t.Parallel()

			p, f, res := generate(t, recording{records: []string{addressName}, body: commentBody})
			records := assert.Rejects(t, "a backend of a target that the recorder does not serve", func(tb assert.TB) {
				gentest.Rendered(tb, f, res, p, tsbackend.New())
			})
			assert.Length(t, records, 1, "the helper stops at the template tree")
			expect.Equal(t, records[0].Contract, "the plugin declares a template tree for the target typescript",
				"the helper stops at the target of the backend")
		})

		t.Run("fails where the render reports a finding", func(t *testing.T) {
			t.Parallel()

			p, f, res := generate(t, recording{records: []string{addressName}, body: brokenBody})
			records := assert.Rejects(t, "a template that writes no Go", func(tb assert.TB) {
				gentest.Rendered(tb, f, res, p, gobackend.New())
			})
			assert.Length(t, records, 1, "the helper stops at the findings of the render")
			expect.Equal(t, records[0].Contract, "the render reports no finding", "the helper stops at the finding")
		})
	})
}

// bare is a plugin that declares nothing but its name.
type bare struct{}

// Name returns the name of bare.
func (bare) Name() plugin.ID { return "bare" }
