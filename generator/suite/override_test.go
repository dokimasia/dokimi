// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"

	"go.dokimi.dev/dokimi/generator/suite"
)

// The calls of a suite override the generator of each enumeration and of each
// struct with a starting value that their inputs take at some depth.
func TestOverride(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			givePkg  string
			giveDecl string
			want     [][]string
		}{
			{
				name: "records the members of each enumeration that the calls take", givePkg: meterPackage,
				giveDecl: "Meter", want: [][]string{{"Grams", "Litres"}},
			},
			{
				name: "records the members of each enumeration at each depth of the inputs", givePkg: callsPackage,
				giveDecl: "Board",
				want:     [][]string{{"Round", "Square"}, {"Light", "Dark"}, {"Red", "Blue"}, {"Warm", "Cool"}},
			},
			{
				name:    "records the members of an enumeration inside a result of a generated function",
				givePkg: inputsPackage, giveDecl: "Pipeline", want: [][]string{{"Raw", "Cooked"}},
			},
			{
				name: "records no enumeration for calls that take none", givePkg: storePackage, giveDecl: "Store",
				want: [][]string{},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, tt.givePkg).suites[tt.giveDecl]
				members := make([][]string, 0, len(got.Using))
				for _, e := range got.Using {
					names := make([]string, 0, len(e.Members))
					for _, v := range e.Members {
						expect.Equal(t, v.Kind, emit.ValueName, "a member is the name of its constant")
						names = append(names, v.Callee.Name)
					}
					members = append(members, names)
				}
				expect.Equal(t, members, tt.want, "the suite overrides each enumeration with its members")
			})
		}

		t.Run("records the members of each enumeration on each kind of path", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, pathsPackage, inputsPackage).suites["Router"]
			members := make([][]string, 0, len(got.Using))
			for _, e := range got.Using {
				names := make([]string, 0, len(e.Members))
				for _, v := range e.Members {
					names = append(names, v.Callee.Name)
				}
				members = append(members, names)
			}
			expect.Equal(t, members, [][]string{
				{"Light", "Heavy"}, {"Slow", "Fast"}, {"North", "South"}, {"Small", "Large"}, {"Raw", "Cooked"},
			}, "a variadic list, a received channel, a pointer, a generated result and a field each lead to an "+
				"enumeration, and a field that is not exported leads to none")
		})

		t.Run("records the starting value of each struct on a path of the inputs", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, pathsPackage, inputsPackage).suites["Router"]
			expect.Equal(t, startsOf(got.Starts), []start{
				{
					name: "RouterItemValues", spelling: "Item", value: "composite Count", label: "item",
					fields: []string{"Count item.count", "Name item.name"},
				},
				{
					name: "RouterItemValues0", spelling: "inputs.Item", seed: "NewItemBuilder", label: "item",
					fields: []string{"ID item.id", "Count item.count", "Kind item.kind", "Base item.base"},
				},
			}, "two structs of one name take two generators, and a generic type or an alias takes none")
		})

		t.Run("records the starting value of each struct with one", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, inputsPackage).suites["Pipeline"]
			expect.Equal(t, startsOf(got.Starts), []start{
				{
					name: "PipelineItemValues", spelling: "Item", seed: "NewItemBuilder", label: "item",
					fields: []string{"ID item.id", "Count item.count", "Kind item.kind", "Base item.base"},
				},
				{
					name: "PipelineLineValues", spelling: "Line", value: "composite Qty", label: "line",
					fields: []string{"SKU line.sku", "Qty line.qty", "Rank line.rank"},
				},
				{
					name: "PipelineNoteValues", spelling: "Note", value: "call DefaultNote", label: "note",
					fields: []string{"Text note.text"},
				},
				{name: "PipelineTagValues", spelling: "Tag", value: "literal new", label: "tag", whole: true},
				{name: "PipelineCodeValues", spelling: "Code", value: "call DefaultCode", label: "code", whole: true},
			}, "the value of a builder, the declared defaults and the authored sample each start a type")
		})

		t.Run("records no starting value of a struct without one", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, callsPackage).suites["Board"]
			expect.Empty(t, got.Starts, "no struct that Board takes has a sample, a builder or a declared default")
		})
	})
}

// start is the generator of the values of a type with a starting value, with
// the spellings of its type, of its seed and of its starting value.
type start struct {
	name, spelling, seed, value, label string
	fields                             []string
	whole                              bool
}

// startsOf returns each generator of starts with the spellings of its type,
// of its seed and of its starting value. A field joins its name and its
// label.
func startsOf(starts []suite.Start) []start {
	out := make([]start, 0, len(starts))
	for _, s := range starts {
		got := start{name: s.Name, spelling: spelled(s.Type), label: s.Label, whole: s.Whole}
		if s.Seed != nil {
			got.seed = s.Seed.Spelling
		}
		switch s.Value.Kind {
		case emit.ValueComposite:
			got.value = "composite"
			for _, f := range s.Value.Fields {
				got.value += " " + f.Name
			}
		case emit.ValueCall:
			got.value = "call " + s.Value.Callee.Name
		case emit.ValueLiteral:
			got.value = "literal " + s.Value.Text
		default:
			// The start of a struct with a builder has a seed and no value.
		}
		for _, f := range s.Fields {
			got.fields = append(got.fields, f.Name+" "+f.Label)
		}
		out = append(out, got)
	}
	return out
}
