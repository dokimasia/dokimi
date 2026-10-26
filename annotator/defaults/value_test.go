// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package defaults_test

import (
	"sync"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/meta"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/annotator/defaults"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// valuesPackage is the package with a default on a field of each kind of
// type.
const valuesPackage = "example.com/shop/values"

// lifted is the value of the text of a default: its kind, the kind of its
// literal and its text, and whether the rules read a value.
type lifted struct {
	kind    emit.ValueKind
	literal emit.LiteralKind
	text    string
	read    bool
}

// The rules of a language read the text of a declared default as a value of
// the type of its field, and the text of a type of text as its content.
func TestValue(t *testing.T) {
	t.Parallel()

	t.Run("Value", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give string
			want lifted
		}{
			{
				name: "returns the literal that the rules read", give: "Number",
				want: lifted{kind: emit.ValueLiteral, literal: emit.LiteralInt, text: "7", read: true},
			},
			{
				name: "returns the content of a text that the rules do not read", give: "Text",
				want: lifted{kind: emit.ValueLiteral, literal: emit.LiteralString, text: "Dam 1", read: true},
			},
			{
				name: "returns the constant that the rules read from its name", give: "Named",
				want: lifted{kind: emit.ValueName, read: true},
			},
			{
				name: "returns the content of a text of a declared type of text", give: "Code",
				want: lifted{kind: emit.ValueLiteral, literal: emit.LiteralString, text: "abc", read: true},
			},
			{name: "reports false for a text that is no number", give: "Bad"},
			{name: "reports false for a text of a declared type of a number", give: "Count"},
			{name: "reports false for a text of a struct", give: "Spot"},
			{name: "reports false for a text of a cycle of declared types", give: "Loop"},
		}
		values := valuesOf(t)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, found := values[tt.give]
				assert.True(t, found, "the field "+tt.give+" declares a default")
				expect.Equal(t, got, tt.want, "the rules read the text as a value of the type of the field")
			})
		}
	})
}

// valuesOf returns the value of the default of each field of Values, by the
// name of the field. A probe reads each value through the rules of the field
// after defaults stamps the text of each default.
func valuesOf(t *testing.T) map[string]lifted {
	t.Helper()

	var mu sync.Mutex
	got := map[string]lifted{}
	probe := sdk.NewPlugin("probe").Version("0.1.0").Handle(sdk.OnField(func(m *sdk.FieldMatch, _ *sdk.Stamper) error {
		f := m.Field
		text, declared := sdk.FactOf(m, f.ID, meta.Named[string](defaults.KeyValue))
		if !declared {
			return nil
		}
		b := m.Rules()
		scope := b.View().ScopeOf(f.ID)
		value, read := defaults.Value(&b, &scope, f.Type, text)
		mu.Lock()
		defer mu.Unlock()
		got[f.Name] = lifted{kind: value.Kind, literal: value.Literal, text: value.Text, read: read}
		return nil
	})).Annotator()
	f := gentest.New(t, &gentest.Spec{
		Dir: goTree, Packages: []string{valuesPackage}, Plugins: []plugin.Plugin{probe},
		Annotators: []plugin.Annotator{defaults.New()},
	})
	assert.NoError(t, f.Annotate(t, probe).Err, "the probe runs whole")
	return got
}
