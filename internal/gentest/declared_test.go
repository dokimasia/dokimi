// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package gentest_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	golang "go.dokimi.dev/eidos/lang/go"
	"go.dokimi.dev/eidos/lang/typescript"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/internal/gentest"
)

// The spelling and the meaning of the code recorded, a code that no package
// registers, and a template that the tree of the recorder lacks.
const (
	recordedCode    = "GENTESTTEST-0002"
	recordedMeaning = "a declaration of a fixture that the recorder of the tests of gentest reports"
	missingTemplate = "missing.tmpl"
)

// unregistered is a code that no package registers.
var unregistered = diag.Code{Prefix: "GENTESTTEST", Number: 99}

// nameOnly is a plugin with a name and without templates.
type nameOnly struct{}

// Name returns the name of the plugin.
func (nameOnly) Name() plugin.ID { return "name-only" }

// Templates checks the template tree that a plugin declares for a target,
// External checks that a family of a plugin tests its package from outside,
// and Codes checks the spelling and the meaning of the codes of a plugin.
func TestDeclared(t *testing.T) {
	t.Parallel()

	t.Run("Templates", func(t *testing.T) {
		t.Parallel()

		t.Run("passes for a tree of the target with each template", func(t *testing.T) {
			t.Parallel()

			gentest.Templates(t, recorder(recording{}), golang.Target, recordFamily, checkFamily)
		})

		t.Run("fails for a plugin without templates", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "a plugin without templates", func(tb assert.TB) {
				gentest.Templates(tb, nameOnly{}, golang.Target, recordFamily)
			})
			assert.Length(t, records, 1, "the helper stops at the plugin")
			expect.Equal(t, records[0].Contract, "the plugin declares templates", "the helper stops at the plugin")
		})

		t.Run("fails for a target without a tree", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "a target without a tree", func(tb assert.TB) {
				gentest.Templates(tb, recorder(recording{}), typescript.Target, recordFamily)
			})
			assert.Length(t, records, 2, "the helper fails at the target and stops at the tree")
			expect.Equal(
				t,
				records[0].Contract,
				"the plugin serves the target typescript",
				"the helper fails at the target",
			)
			expect.Equal(
				t,
				records[1].Contract,
				"the target typescript has a template tree",
				"the helper stops at the tree",
			)
		})

		t.Run("fails for a tree without a template", func(t *testing.T) {
			t.Parallel()

			records := assert.Rejects(t, "a tree without a template", func(tb assert.TB) {
				gentest.Templates(tb, recorder(recording{}), golang.Target, recordFamily, missingTemplate)
			})
			assert.Length(t, records, 1, "the helper fails at the missing template alone")
			expect.Equal(t, records[0].Contract, "the tree of golang contains "+missingTemplate,
				"the helper fails at the template")
		})
	})

	t.Run("External", func(t *testing.T) {
		t.Parallel()

		t.Run("passes for a family whose files check their package from outside", func(t *testing.T) {
			t.Parallel()

			gentest.External(t, recorder(recording{external: true}), sdk.TagTest)
		})

		tests := []struct {
			name         string
			give         plugin.Plugin
			giveTag      string
			wantContract string
		}{
			{
				name: "fails for a plugin without families", give: nameOnly{}, giveTag: sdk.TagTest,
				wantContract: "the plugin declares its families",
			},
			{
				name: "fails for a tag of no family", give: recorder(recording{external: true}), giveTag: "bench",
				wantContract: "the plugin declares the family of the tag bench",
			},
			{
				name: "fails for a family whose files are in the package that they check",
				give: recorder(recording{}), giveTag: sdk.TagTest,
				wantContract: "the files of the family test their package from outside",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				records := assert.Rejects(t, "a family that is not external", func(tb assert.TB) {
					gentest.External(tb, tt.give, tt.giveTag)
				})
				assert.Length(t, records, 1, "the helper fails once")
				expect.Equal(t, records[0].Contract, tt.wantContract, "the helper fails at the difference")
			})
		}
	})

	t.Run("Codes", func(t *testing.T) {
		t.Parallel()

		t.Run("passes for a complete code", func(t *testing.T) {
			t.Parallel()

			gentest.Codes(t, gentest.Code{
				Name: "recorded", Give: recorded, WantCode: recordedCode, WantMeaning: recordedMeaning,
			})
		})

		tests := []struct {
			name         string
			give         gentest.Code
			wantContract string
		}{
			{
				name: "fails for a code with another spelling",
				give: gentest.Code{
					Name: "recorded", Give: recorded, WantCode: "GENTESTTEST-0003", WantMeaning: recordedMeaning,
				},
				wantContract: "recorded has the prefix of its package and its number",
			},
			{
				name: "fails for a code that no package registers",
				give: gentest.Code{
					Name: "unregistered", Give: unregistered, WantCode: "GENTESTTEST-0099",
				},
				wantContract: "the package registers unregistered",
			},
			{
				name: "fails for a code with another meaning",
				give: gentest.Code{
					Name: "recorded", Give: recorded, WantCode: recordedCode, WantMeaning: "another meaning",
				},
				wantContract: "the meaning of recorded anchors the code in the index",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				records := assert.Rejects(t, "a code that differs from its case", func(tb assert.TB) {
					gentest.Codes(tb, tt.give)
				})
				assert.Length(t, records, 1, "the helper fails once")
				expect.Equal(t, records[0].Contract, tt.wantContract, "the helper fails at the difference")
			})
		}
	})
}
