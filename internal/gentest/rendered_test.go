// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest_test

import (
	"io/fs"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
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

// Rendered renders what a generator recorded through a backend, with the
// templates of the generator, and returns the files by their routed paths.
func TestRendered(t *testing.T) {
	t.Parallel()

	t.Run("Rendered", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each file that the render writes by its routed path", func(t *testing.T) {
			t.Parallel()

			p, _, res := generate(t, recording{records: []string{addressName}, body: commentBody})
			got := gentest.Rendered(t, res, p, gobackend.New())
			paths, err := fs.Glob(got, "*/*")
			assert.NoError(t, err, "the files of the render list")
			assert.Equal(t, paths, []string{recordPath, checkPath},
				"the render writes the file of each family of the recorder beside the source")
			expect.Contains(t, string(got[recordPath].Data), "// Address is recorded.",
				"the file contains what the template of the recorder writes")
		})

		t.Run("fails where the plugin declares no templates", func(t *testing.T) {
			t.Parallel()

			_, _, res := generate(t, recording{records: []string{addressName}, body: commentBody})
			records := assert.Rejects(t, "a plugin without templates", func(tb assert.TB) {
				gentest.Rendered(tb, res, bare{}, gobackend.New())
			})
			assert.Length(t, records, 1, "the helper stops at the templates")
			expect.Equal(t, records[0].Contract, "the plugin declares templates", "the helper stops at the plugin")
		})

		t.Run("fails where the plugin declares no template tree for the target of the backend", func(t *testing.T) {
			t.Parallel()

			p, _, res := generate(t, recording{records: []string{addressName}, body: commentBody})
			records := assert.Rejects(t, "a backend of a target that the recorder does not serve", func(tb assert.TB) {
				gentest.Rendered(tb, res, p, tsbackend.New())
			})
			assert.Length(t, records, 1, "the helper stops at the template tree")
			expect.Equal(t, records[0].Contract, "the plugin declares a template tree for the target typescript",
				"the helper stops at the target of the backend")
		})

		t.Run("fails where the render reports a finding", func(t *testing.T) {
			t.Parallel()

			p, _, res := generate(t, recording{records: []string{addressName}, body: brokenBody})
			records := assert.Rejects(t, "a template that writes no Go", func(tb assert.TB) {
				gentest.Rendered(tb, res, p, gobackend.New())
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
