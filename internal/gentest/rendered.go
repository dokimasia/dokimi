// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest

import (
	"io/fs"
	"slices"
	"testing/fstest"
	"text/template"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/sdk/backend"
	"go.dokimi.dev/eidos/sdk/backendtest"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
)

// Rendered renders the units of the phase call res through the backend b, as
// the render of a plan renders them after the plugin p generated them over
// the fixture f. The plan runs the generators before ahead of p, as the plan
// of the composition runs builder ahead of stub, so the store of res contains
// their units too. The layout of the plan routes each file beside its source,
// or where an out directive moves it, and the file of each family of one file
// per plan to the root of the tree. The render reads the template tree, the
// helpers and the overrides that each generator of the plan declares for the
// target of b. Rendered returns the files of p by their routed paths, so a
// spec matches them with golden.MatchTree. It stops the test where a
// generator of the plan declares no template tree for the target of b, and
// where the layout or the render reports a finding, because such a finding
// is a defect of a template or of what a generator records.
func Rendered(
	tb assert.TB, f *Fixture, res plugintest.Result, p plugin.Plugin, b backend.Backend, before ...plugin.Plugin,
) fstest.MapFS {
	tb.Helper()

	target := b.Target()
	plan := append(slices.Clone(before), p)
	fixture := &backendtest.Fixture{
		Emit:      res.Emit,
		Trees:     map[plugin.ID]fs.FS{},
		Funcs:     map[plugin.ID]template.FuncMap{},
		Overrides: map[plugin.ID][]string{},
		Layout:    f.Layout(tb, plan...),
	}
	for _, q := range plan {
		templates, provides := q.(plugin.TemplateProvider)
		assert.True(tb, provides, "the plugin declares templates")
		tree, declared := templates.Templates(target)
		assert.True(tb, declared, "the plugin declares a template tree for the target "+string(target))
		fixture.Schedule = append(fixture.Schedule, q.Name())
		fixture.Trees[q.Name()] = tree
		fixture.Funcs[q.Name()] = templates.TemplateFuncs(target)
		fixture.Overrides[q.Name()] = templates.Overrides(target)
	}
	files, findings := backendtest.Rendered(tb, func(assert.TB) (backend.Backend, *backendtest.Fixture) {
		return b, fixture
	})
	assert.Empty(tb, findings, "the render reports no finding")
	out := fstest.MapFS{}
	for i := range files {
		if slices.Contains(files[i].Plugins, p.Name()) {
			out[files[i].Path] = &fstest.MapFile{Data: files[i].Body}
		}
	}
	return out
}
