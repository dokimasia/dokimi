// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest

import (
	"io/fs"
	"testing/fstest"
	"text/template"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/eidos/sdk/backend"
	"go.dokimi.dev/eidos/sdk/backendtest"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/plugintest"
)

// Rendered renders the units of the phase call res through the backend b, as
// the render of a plan renders them after the plugin p generated them. The
// render reads the template tree, the helpers and the overrides that p
// declares for the target of b. Rendered returns the files by their routed
// paths, so a spec matches them with golden.MatchTree. It stops the test
// where p declares no template tree for the target of b, and where the
// render reports a finding, because such a finding is a defect of a template
// or of what p records.
func Rendered(tb assert.TB, res plugintest.Result, p plugin.Plugin, b backend.Backend) fstest.MapFS {
	tb.Helper()

	target := b.Target()
	templates, provides := p.(plugin.TemplateProvider)
	assert.True(tb, provides, "the plugin declares templates")
	tree, declared := templates.Templates(target)
	assert.True(tb, declared, "the plugin declares a template tree for the target "+string(target))
	files, findings := backendtest.Rendered(tb, func(assert.TB) (backend.Backend, *backendtest.Fixture) {
		return b, &backendtest.Fixture{
			Emit:      res.Emit,
			Schedule:  []plugin.ID{p.Name()},
			Trees:     map[plugin.ID]fs.FS{p.Name(): tree},
			Funcs:     map[plugin.ID]template.FuncMap{p.Name(): templates.TemplateFuncs(target)},
			Overrides: map[plugin.ID][]string{p.Name(): templates.Overrides(target)},
		}
	})
	assert.Empty(tb, findings, "the render reports no finding")
	out := fstest.MapFS{}
	for i := range files {
		out[files[i].Path] = &fstest.MapFile{Data: files[i].Body}
	}
	return out
}
