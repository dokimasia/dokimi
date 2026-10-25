// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// The calls of a suite override the generator of each enumeration, of each
// struct with a starting value, of each function type and of each interface
// that their inputs take at some depth.
func TestOverride(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			givePkgs []string
			giveDecl string
			want     [][]string
		}{
			{
				name: "records the members of each enumeration that the calls take", givePkgs: []string{meterPackage},
				giveDecl: "Meter", want: [][]string{{"Grams", "Litres"}},
			},
			{
				name:     "records the members of each enumeration at each depth of the inputs",
				givePkgs: []string{callsPackage}, giveDecl: "Board",
				want: [][]string{{"Round", "Square"}, {"Light", "Dark"}, {"Red", "Blue"}, {"Warm", "Cool"}},
			},
			{
				name:     "records the members of an enumeration inside a result of a generated function",
				givePkgs: []string{inputsPackage}, giveDecl: "Pipeline", want: [][]string{{"Raw", "Cooked"}},
			},
			{
				name: "records no enumeration for calls that take none", givePkgs: []string{storePackage},
				giveDecl: "Store", want: [][]string{},
			},
			{
				name:     "does not record an enumeration that the file of the suite cannot import",
				givePkgs: []string{wrapPackage, toolPackage, deepPackage}, giveDecl: "Tuned",
				want: [][]string{{"Percent", "Steps"}},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, tt.givePkgs...).suites[tt.giveDecl]
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

		dynamic := []struct {
			name     string
			giveDecl string
			want     bool
		}{
			{
				name:     "records the generator of any for a variadic parameter of any type",
				giveDecl: "Printer",
				want:     true,
			},
			{
				name:     "records the generator of any for a parameter of the empty interface",
				giveDecl: "Recorder",
				want:     true,
			},
			{name: "records the generator of any for a field of any type of an input", giveDecl: "Journal", want: true},
			{name: "records no generator of any for calls that take no value of any type", giveDecl: "Plain"},
		}
		for _, tt := range dynamic {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, loosePackage).suites[tt.giveDecl]
				expect.Equal(t, got.Dynamic, tt.want, "the suite draws a value of any type where its calls take one")
			})
		}

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

		t.Run("does not record the starting value of a type that the file of the suite cannot import",
			func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, wrapPackage, toolPackage, deepPackage).suites["Tuned"]
				expect.Equal(t, startsOf(got.Starts), []start{{
					name: "TunedSettingValues", spelling: "Setting", value: "composite Value", label: "setting",
					fields: []string{"Value setting.value"},
				}}, "the spec and the level of a tool are types of the internal package of the tree lib")
			})

		labeler := []fn{{
			field: "OnLabel", signature: "(string) (r0 Receipt, r1 error)", index: "labelerLabelAt",
			results: []result{
				{name: "labelerLabelResults0", label: "labeler-label-0", spelling: "Receipt"},
				{
					name: "labelerLabelResults1", label: "labeler-label-1", spelling: "error", isError: true,
					errors: []string{"ErrTorn"},
				},
			},
		}}
		generators := []struct {
			name     string
			givePkgs []string
			giveDecl string
			want     []generator
		}{
			{
				name:     "records the generator of each function type and each interface at each depth of the inputs",
				givePkgs: []string{nestedPackage}, giveDecl: "Courier",
				want: []generator{
					{spelling: "Labeler", input: "CourierLabelerInput", funcs: labeler},
					{spelling: "func() error", funcs: []fn{{
						signature: "() (r0 error)", index: "funcAt",
						results: []result{{name: "funcResults", label: "func", spelling: "error", isError: true}},
					}}},
					{spelling: "Ack", funcs: []fn{{
						signature: "(string) (r0 error)", index: "ackAt",
						results: []result{{name: "ackResults", label: "ack", spelling: "error", isError: true}},
					}}},
					{spelling: "Scale[int]", input: "CourierScaleInput", funcs: []fn{{
						field: "OnWeigh", signature: "() (r0 int)", index: "scaleWeighAt",
						results: []result{{name: "scaleWeighResults", label: "scale-weigh", spelling: "int"}},
					}}},
					{spelling: "func()", funcs: []fn{{signature: "() ()"}}},
				},
			},
			{
				name:     "records the generator of an interface inside a result of a generated function",
				givePkgs: []string{nestedPackage}, giveDecl: "Hire",
				want: []generator{
					{spelling: "Labeler", input: "HireLabelerInput", funcs: labeler},
					{spelling: "func()", funcs: []fn{{signature: "() ()"}}},
				},
			},
			{
				name:     "records the generator of each type at depth that the file of the suite can import",
				givePkgs: []string{wrapPackage, toolPackage, deepPackage}, giveDecl: "Tuned",
				want: []generator{{spelling: "Note", input: "TunedNoteInput", funcs: []fn{{
					field: "OnWrite", signature: "(string) (r0 error)", index: "noteWriteAt",
					results: []result{
						{name: "noteWriteResults", label: "note-write", spelling: "error", isError: true},
					},
				}}}},
			},
			{
				name:     "does not record a generator for calls without a function type or an interface",
				givePkgs: []string{storePackage}, giveDecl: "Store", want: []generator{},
			},
		}
		for _, tt := range generators {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, tt.givePkgs...).suites[tt.giveDecl]
				gens := make([]generator, 0, len(got.Generated))
				for _, g := range got.Generated {
					gens = append(gens, generator{spelling: spelled(g.Type), input: g.Input, funcs: fnsOf(g.Funcs)})
				}
				expect.Equal(t, gens, tt.want, "the suite generates the values of each such type once")
			})
		}

		t.Run("leaves the import of a type at depth to the layout where the plan does not preview a package",
			func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				s := specOf(goTree, p, wrapPackage, toolPackage, deepPackage)
				s.Types = nil
				got := recorded(t, gentest.New(t, s), p).suites["Tuned"]
				spellings := make([]string, 0, len(got.Generated))
				for _, g := range got.Generated {
					spellings = append(spellings, spelled(g.Type))
				}
				expect.Equal(t, spellings, []string{
					"deep.Probe", "Note", "func() deep.Level", "deep.Hook", "Dial", "Gauge", "Calibrator",
				}, "the suite generates the values of each such type, and the layout checks its references")
			})
	})
}

// generator is the generator of the values of a function type or of an
// interface, with the spelling of the type and its generated functions.
type generator struct {
	spelling, input string
	funcs           []fn
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
