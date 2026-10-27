// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert"
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
			{
				name:     "records only the members that the file of the suite can refer to",
				givePkgs: []string{wrapPackage, toolPackage, deepPackage}, giveDecl: "Limited",
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
					name: "RouterItemValues0", spelling: "inputs.Item", value: "composite Count", label: "item",
					fields: []string{"ID item.id", "Count item.count", "Kind item.kind", "Base item.base"},
				},
			}, "two structs of one name take two generators, and a generic type or an alias takes none")
		})

		t.Run("takes the declared defaults of a struct of another package in place of its builder", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, pathsPackage, inputsPackage).suites["Router"]
			seeds := map[string]string{}
			for _, s := range startsOf(got.Starts) {
				seeds[s.name] = s.seed + s.value
			}
			expect.Equal(t, seeds["RouterItemValues0"], "composite Count",
				"the builder of inputs.Item is in the package inputs, and the suite of Router is in the package paths")
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

		starts := []struct {
			name string
			give string
			want []start
		}{
			{
				name: "leaves a field at its zero value where the file of the suite cannot refer to its default",
				give: "LimitedLimitsValues",
				want: []start{{
					name: "LimitedLimitsValues", spelling: "tool.Limits", value: "composite Count", label: "limits",
					fields: []string{"Level limits.level", "Count limits.count"},
				}},
			},
			{
				name: "takes the declared defaults where the file of the suite cannot refer to the authored sample",
				give: "LimitedGripValues",
				want: []start{{
					name: "LimitedGripValues", spelling: "tool.Grip", value: "composite Force", label: "grip",
					fields: []string{"Force grip.force"},
				}},
			},
			{
				name: "starts a struct without a field that the property engine reads from its zero value",
				give: "LimitedClampValues",
				want: []start{{name: "LimitedClampValues", spelling: "tool.Clamp", value: "composite", label: "clamp"}},
			},
			{
				name: "starts a struct that embeds a type that its package does not export from its zero value",
				give: "LimitedFrameValues",
				want: []start{{name: "LimitedFrameValues", spelling: "tool.Frame", value: "composite", label: "frame"}},
			},
			{
				name: "does not record the starting value of a struct whose embedded type the property engine reads",
				give: "LimitedMountValues", want: []start{},
			},
			{
				name: "takes the declared default of a defined type where the file of the suite cannot refer to its sample",
				give: "LimitedSpanValues",
				want: []start{{
					name: "LimitedSpanValues", spelling: "tool.Span", value: "literal 8", label: "span", whole: true,
				}},
			},
			{
				name: "does not record a defined type whose declared default the file of the suite cannot refer to",
				give: "LimitedPitchValues", want: []start{},
			},
		}
		for _, tt := range starts {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, wrapPackage, toolPackage, deepPackage).suites["Limited"]
				named := []start{}
				for _, s := range startsOf(got.Starts) {
					if s.name == tt.give {
						named = append(named, s)
					}
				}
				expect.Equal(t, named, tt.want,
					"the suite takes the first source of the value that its file can refer to")
			})
		}

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

		t.Run("bounds the base values of each type whose generated results contain the type", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, nestedPackage).suites["Tracker"]
			leaves := map[string]int{}
			for _, g := range got.Generated {
				leaves[spelled(g.Type)] = g.Leaves
			}
			expect.Equal(t, leaves, map[string]int{
				"Node": 10, "Branch": 10, "Step": 10, "Ping": 10, "Pong": 10, "func()": 0, "Pair": 0, "Leaf": 0,
				"Rally": 0, "Linker": 0,
			}, "a result of the type, a field of a result, the result of a function type and two interfaces that "+
				"return each other make a type recursive, and a receipt, a pair of leaves, a rally or a linker "+
				"does not")
		})

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
					"interface{ Depth() deep.Level }",
				}, "the suite generates the values of each such type, and the layout checks its references")
			})

		t.Run("records a stream of each channel that a subject receives from", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Router"]
			streams := make([][2]string, 0, len(got.Streams))
			for _, s := range got.Streams {
				streams = append(streams, [2]string{spelled(s.Type), spelled(s.Elem)})
			}
			expect.Equal(
				t,
				streams,
				[][2]string{{"<-chan string", "string"}, {"chan string", "string"}, {"Ticks", "int"}},
				"a receive channel, a channel of both directions and a defined channel type each take one stream",
			)
		})

		t.Run("records the sink of each send channel", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Router"]
			sinks := make([]string, 0, len(got.Sinks))
			for _, s := range got.Sinks {
				sinks = append(sinks, spelled(s))
			}
			expect.Equal(t, sinks, []string{"chan<- string"},
				"a field, a parameter and a result of an input share the sink of one send channel type")
		})

		t.Run("records a stream of the channel of a timer that a clock returns", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Placer"]
			assert.Length(t, got.Streams, 1, "the suite of Placer has one stream")
			expect.Equal(t, [2]string{spelled(got.Streams[0].Type), spelled(got.Streams[0].Elem)},
				[2]string{"<-chan time.Time", "time.Time"}, "the method C of a timer returns a stream of times")
		})

		t.Run("records an implementation of an interface with an unexported method", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Drawer"]
			assert.Length(t, got.Variants, 1, "the suite of Drawer has the variants of Shape")
			expect.Equal(t, spelled(got.Variants[0].Type), "Shape", "the variants generate a shape")
			members := make([]variant, 0, len(got.Variants[0].Members))
			for _, m := range got.Variants[0].Members {
				members = append(members, variant{spelling: spelled(m.Type), label: m.Label, address: m.Address})
			}
			expect.Equal(t, members, []variant{
				{spelling: "Circle", label: "circle"},
				{spelling: "Corner", label: "corner"},
				{spelling: "Square", label: "square", address: true},
			}, "a shape is a circle, a corner or the address of a square, and not a blob, a dot or a polygon")
		})

		t.Run("records a reader of drawn bytes for each interface of io and a buffered reader", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Std"]
			readers := make([][3]string, 0, len(got.Readers))
			for _, r := range got.Readers {
				readers = append(readers, [3]string{spelled(r.Type), r.Label, r.Wrap.Callee.Name})
			}
			expect.Equal(t, readers, [][3]string{
				{"io.Writer", "writer", ""},
				{"io.Reader", "reader", ""},
				{"io.ReadCloser", "read-closer", ""},
				{"*bufio.Reader", "reader", "NewReader"},
			}, "an interface of io takes the bytes, and a buffered reader reads them")
		})

		t.Run("records the value of each type that takes one fixed value", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Std"]
			values := make([][2]string, 0, len(got.Constants))
			for _, c := range got.Constants {
				value := c.Value.Callee.Package + "." + c.Value.Callee.Name
				if c.Value.Kind == emit.ValueLiteral {
					value = c.Value.Literal.String()
				}
				values = append(values, [2]string{spelled(c.Type), value})
			}
			expect.Equal(t, values, [][2]string{
				{"*slog.Logger", "log/slog.New"},
				{"*log.Logger", "log.New"},
				{"*bufio.Writer", "bufio.NewWriter"},
				{"context.Context", "context.Background"},
				{"unsafe.Pointer", "nil"},
			}, "a logger and a writer take a constructor, a context takes the background, and a pointer takes nil")
		})

		errors := []struct {
			name     string
			giveDecl string
			want     bool
		}{
			{
				name:     "records the generator of error for an error outside the results of a function",
				giveDecl: "Std",
				want:     true,
			},
			{name: "records no generator of error for calls that take no error", giveDecl: "Placer"},
		}
		for _, tt := range errors {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, kindsPackage).suites[tt.giveDecl]
				expect.Equal(t, got.Errors, tt.want, "the suite draws an error where its calls take one")
			})
		}

		t.Run("records the zero value of each struct without a readable field", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Std"]
			expect.Equal(t, startsOf(got.Starts), []start{
				{name: "StdMutexValues", spelling: "sync.Mutex", value: "composite", label: "mutex"},
				{name: "StdBoxValues", spelling: "Box[int]", value: "composite", label: "box"},
				{name: "StdBuilderValues", spelling: "strings.Builder", value: "composite", label: "builder"},
				{name: "StdStructValues", spelling: "struct{ n int }", value: "composite", label: "struct"},
			}, "a struct of the standard library, a generic struct and an inline struct start from their zero value")
		})

		t.Run("records a stream of a channel in a field of a struct of the standard library", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Std"]
			streams := make([]string, 0, len(got.Streams))
			for _, s := range got.Streams {
				streams = append(streams, spelled(s.Type))
			}
			expect.Equal(t, streams, []string{"<-chan Time"}, "the field C of a timer is a receive channel")
		})

		reports := []struct {
			name     string
			givePkgs []string
			giveDecl string
			want     []string
		}{
			{
				name:     "reports Ungenerated for an interface whose package exports no implementation of it",
				givePkgs: []string{kindsPackage}, giveDecl: "Drawer",
				want: []string{
					"Drawer.Pick takes Sealed in a result of Picker.Pick in its parameter p. The interface Sealed has an " +
						"unexported method, and the package example.com/shop/kinds exports no implementation that the " +
						"file of the suite can import. Each run of Drawer.Pick fails before its first case until the test " +
						"process registers a generator of Sealed",
					"Drawer.Test takes testing.TB in its parameter tb. The interface testing.TB of the standard library " +
						"has an unexported method, so a test cannot build an implementation of it. Each run of Drawer.Test " +
						"fails before its first case until the test process registers a generator of testing.TB",
				},
			},
			{
				name:     "reports Ungenerated for a struct with two fields of one JSON name",
				givePkgs: []string{kindsPackage},
				giveDecl: "Std",
				want: []string{
					"Std.Pair takes Twins in an element of []Twins in its parameter t. The fields A and B of Twins are " +
						"both named \"x\". Each run of Std.Pair fails before its first case until the test process " +
						"registers a generator of Twins",
					"Std.Meet takes Clash in its parameter c. The fields Name and Circle of Clash are both named " +
						"\"Circle\". Each run of Std.Meet fails before its first case until the test process registers a " +
						"generator of Clash",
				},
			},
			{
				name:     "reports Ungenerated for each type whose generated code the file of the suite cannot import",
				givePkgs: []string{wrapPackage, toolPackage, deepPackage}, giveDecl: "Tuned",
				want: []string{
					"Tuned.Tune takes deep.Probe in the field Probe of tool.Request in its parameter r. The type " +
						"deep.Probe refers to a declaration that the file of the suite cannot import. Each run of " +
						"Tuned.Tune fails before its first case until the test process registers a generator of deep.Probe",
					"Tuned.Tune takes func() deep.Level in the field Read of tool.Request in its parameter r. The type " +
						"func() deep.Level refers to a declaration that the file of the suite cannot import. Each run of " +
						"Tuned.Tune fails before its first case until the test process registers a generator of " +
						"func() deep.Level",
					"Tuned.Tune takes deep.Hook in the field Hook of tool.Request in its parameter r. The type " +
						"deep.Hook refers to a declaration that the file of the suite cannot import. Each run of " +
						"Tuned.Tune fails before its first case until the test process registers a generator of deep.Hook",
					"Tuned.Tune takes Dial in the field Dial of tool.Request in its parameter r. The signature of Dial " +
						"refers to a declaration that the file of the suite cannot import. Each run of Tuned.Tune fails " +
						"before its first case until the test process registers a generator of Dial",
					"Tuned.Tune takes Gauge in the field Gauge of tool.Request in its parameter r. A method of Gauge " +
						"refers to a declaration that the file of the suite cannot import. Each run of Tuned.Tune fails " +
						"before its first case until the test process registers a generator of Gauge",
					"Tuned.Tune takes Calibrator in the field Calibrator of tool.Request in its parameter r. A method " +
						"of Calibrator refers to a declaration that the file of the suite cannot import. Each run of " +
						"Tuned.Tune fails before its first case until the test process registers a generator of Calibrator",
					"Tuned.Tune takes <-chan deep.Level in the field Levels of tool.Request in its parameter r. The " +
						"channel <-chan deep.Level refers to a declaration that the file of the suite cannot import. Each " +
						"run of Tuned.Tune fails before its first case until the test process registers a generator of " +
						"<-chan deep.Level",
					"Tuned.Tune takes interface{ Depth() deep.Level } in the field Depth of tool.Request in its " +
						"parameter r. A method of interface{ Depth() deep.Level } refers to a declaration that the file " +
						"of the suite cannot import. Each run of Tuned.Tune fails before its first case until the test " +
						"process registers a generator of interface{ Depth() deep.Level }",
					"Tuned.Tune takes Shaped in the field Shaped of tool.Request in its parameter r. The interface " +
						"Shaped has an unexported method, and the package example.com/shop/lib/internal/deep exports no " +
						"implementation that the file of the suite can import. Each run of Tuned.Tune fails before its " +
						"first case until the test process registers a generator of Shaped",
				},
			},
			{
				name:     "reports no Ungenerated for calls whose types each have a generator",
				givePkgs: []string{kindsPackage},
				giveDecl: "Router",
			},
		}
		for _, tt := range reports {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := ungenerated(generate(t, goTree, tt.givePkgs...), tt.giveDecl)
				expect.Equal(t, got, tt.want, "the suite reports each type without a generator once, at its first path",
					assert.EquateEmpty())
			})
		}

		t.Run("leaves the references of a value to the layout where the plan does not preview a package",
			func(t *testing.T) {
				t.Parallel()

				p := suite.New()
				s := specOf(goTree, p, wrapPackage, toolPackage, deepPackage)
				s.Types = nil
				got := recorded(t, gentest.New(t, s), p).suites["Limited"]
				values := map[string]string{}
				for _, start := range startsOf(got.Starts) {
					values[start.name] = start.value
				}
				expect.Equal(t, values["LimitedLimitsValues"], "composite Level Count",
					"the suite keeps the default of the internal package, and the layout checks its references")
			})
	})
}

// generator is the generator of the values of a function type or of an
// interface, with the spelling of the type and its generated functions.
type generator struct {
	spelling, input string
	funcs           []fn
}

// variant is one implementation of an interface with an unexported method,
// with the spelling of its type.
type variant struct {
	spelling, label string
	address         bool
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
