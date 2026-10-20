// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package lang_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/plugin"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/internal/lang"
)

// templateName is the template that the tree of each part of the specs
// contains.
const templateName = "fixture.tmpl"

// fixturePart is the part of a language of the specs. Its target has the
// name of its language, and its tree contains the template.
type fixturePart struct {
	lang symbol.Lang
}

// Lang returns the language of the part.
func (p fixturePart) Lang() symbol.Lang { return p.lang }

// Target returns the target of the part.
func (p fixturePart) Target() plugin.Target { return plugin.Target(p.lang) }

// Options returns the template tree of the part.
func (fixturePart) Options() []sdk.TargetOption {
	return []sdk.TargetOption{sdk.Templates(fstest.MapFS{templateName: {Data: []byte("{{range .Records}}{{end}}")}})}
}

// match is a match of a subject of one language.
type match struct {
	lang symbol.Lang
}

// Lang returns the language of the subject.
func (m match) Lang() symbol.Lang { return m.lang }

// The parts of the specs, one for Go and one for TypeScript.
var (
	golangPart     = fixturePart{lang: "golang"}
	typescriptPart = fixturePart{lang: "typescript"}
	parts          = []fixturePart{golangPart, typescriptPart}
)

// errHandler is the error of a handler of the specs.
var errHandler = errors.New("lang: the error of the handler")

// Served runs a handler with the part of the language of its subject, and
// Register declares the target of each part on a plugin.
func TestPart(t *testing.T) {
	t.Parallel()

	t.Run("Served", func(t *testing.T) {
		t.Parallel()

		t.Run("runs the handler with the part of the language of the match", func(t *testing.T) {
			t.Parallel()

			var got []fixturePart
			h := lang.Served(parts, func(_ match, p fixturePart, _ *sdk.Emitter) error {
				got = append(got, p)
				return nil
			})
			assert.NoError(t, h(match{lang: "typescript"}, nil), "the handler runs")
			expect.Equal(t, got, []fixturePart{typescriptPart}, "the handler runs once with the part of TypeScript")
		})

		t.Run("does nothing for a match of a language without a part", func(t *testing.T) {
			t.Parallel()

			ran := false
			h := lang.Served(parts, func(match, fixturePart, *sdk.Emitter) error {
				ran = true
				return nil
			})
			assert.NoError(t, h(match{lang: "scripted"}, nil), "the handler of another language returns no error")
			expect.False(t, ran, "the handler does not run for a language without a part")
		})

		t.Run("returns the error of the handler", func(t *testing.T) {
			t.Parallel()

			h := lang.Served(parts, func(match, fixturePart, *sdk.Emitter) error { return errHandler })
			expect.ErrorIs(t, h(match{lang: "golang"}, nil), errHandler, "the error of the handler returns")
		})
	})

	t.Run("Register", func(t *testing.T) {
		t.Parallel()

		t.Run("declares the target of each part with its options", func(t *testing.T) {
			t.Parallel()

			b := sdk.NewPlugin("fixture").
				Version("0.1.0").
				Output(plugin.Output{Per: plugin.PerSource, Word: "fixture", Template: templateName})
			lang.Register(b, parts)
			built := b.Handle(sdk.OnStruct(func(*sdk.StructMatch, *sdk.Emitter) error { return nil })).Generator()
			templates, provides := built.(plugin.TemplateProvider)
			assert.True(t, provides, "the plugin declares templates")
			assert.Equal(t, templates.TemplateTargets(), []plugin.Target{"golang", "typescript"},
				"the plugin serves the target of each part")
			for _, target := range templates.TemplateTargets() {
				tree, declared := templates.Templates(target)
				assert.True(t, declared, "the target "+string(target)+" has a template tree")
				_, err := fs.Stat(tree, templateName)
				expect.NoError(t, err, "the tree of "+string(target)+" is the tree of its part")
			}
		})
	})
}
