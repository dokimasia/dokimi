// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/suite"
)

// The suite declares an implementation of each interface that its calls
// take, whose methods return the results that a call generator draws for
// them.
func TestInput(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("declares one implementation of each interface at each type argument", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, inputsPackage).suites["Pipeline"]
			names := make([][2]string, 0, len(got.Inputs))
			for _, in := range got.Inputs {
				names = append(names, [2]string{in.Name, spelled(in.Type)})
			}
			expect.Equal(t, names, [][2]string{
				{"PipelineFeedInput", "Feed"},
				{"PipelineSourceInput", "Source[int]"},
				{"PipelineSourceInput0", "Source[string]"},
			}, "two parameters of one interface share its implementation")
		})

		tests := []struct {
			name string
			give int
			want [][2]string
		}{
			{
				name: "names each parameter of a method of an implementation by its position", give: 0,
				want: [][2]string{
					{"OnNext", "Next(a0 context.Context) (Item, bool, error)"},
					{"OnGet", "Get(a0 context.Context, a1 string) (Item, error)"},
					{"OnSeek", "Seek(a0 ...string) ()"},
					{"OnClose", "Close() (error)"},
				},
			},
			{
				name: "declares the methods of a generic interface at the type arguments of its parameter", give: 2,
				want: [][2]string{{"OnRead", "Read() (string)"}},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, inputsPackage).suites["Pipeline"]
				assert.True(t, tt.give < len(got.Inputs), "the suite declares the implementation of its case")
				methods := make([][2]string, 0, len(got.Inputs[tt.give].Methods))
				for _, m := range got.Inputs[tt.give].Methods {
					methods = append(methods, [2]string{m.Field, signature(m.Method)})
				}
				expect.Equal(t, methods, tt.want, "each method passes its parameters to the function of its field")
			})
		}

		t.Run("draws the generator of an interface of the standard library", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, inputsPackage).suites["Pipeline"], "Copy")
			assert.Length(t, got.Params, 1, "Copy takes a reader")
			expect.Equal(t, got.Params[0].Input, "", "a reader of io takes the reader of drawn bytes")
			expect.Equal(t, spelled(got.Params[0].Type), "io.Reader", "the call draws a value of the reader")
		})

		t.Run("declares an implementation of an interface of the standard library", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, kindsPackage).suites["Std"], "From")
			assert.Length(t, got.Params, 1, "From takes a reader from")
			expect.Equal(t, got.Params[0].Input, "StdReaderFromInput",
				"an interface of io without the methods of a reader of bytes takes a generated implementation")
		})

		t.Run("declares an implementation of each interface without a declaration of its own", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Odd"]
			names := make([][2]string, 0, len(got.Inputs))
			for _, in := range got.Inputs {
				names = append(names, [2]string{in.Name, spelled(in.Type)})
			}
			expect.Equal(t, names, [][2]string{
				{"OddStringerInput", "fmt.Stringer"},
				{"OddAnythingInput", "any"},
				{"OddInterfaceInput", "interface {\n\t\tNamer\n\t\tSize() int\n\t\tReset()\n\t}"},
			}, "a defined type implements the interface that it refers to, and an inline interface its own body")
		})

		t.Run("declares the methods of the types that an inline interface embeds", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, kindsPackage).suites["Odd"]
			assert.Length(t, got.Inputs, 3, "the suite of Odd declares three implementations")
			methods := make([][2]string, 0, len(got.Inputs[2].Methods))
			for _, m := range got.Inputs[2].Methods {
				methods = append(methods, [2]string{m.Field, signature(m.Method)})
			}
			expect.Equal(t, methods, [][2]string{
				{"OnSize", "Size() (int)"}, {"OnReset", "Reset() ()"}, {"OnName", "Name() (string)"},
			}, "the implementation has the methods of the body and the method of the embedded Namer")
		})

		t.Run("takes the end value of each channel after its drawn results", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, kindsPackage).suites["Router"], "Watch")
			assert.Length(t, got.Params, 1, "Watch takes a source")
			expect.Equal(t, fnsOf(got.Params[0].Funcs), []fn{
				{
					field: "OnUpdates", signature: "() (r0 <-chan string)", index: "sUpdatesAt",
					results: []result{{
						name: "sUpdatesResults", label: "s-updates", spelling: "<-chan string", end: "sUpdatesEnd",
					}},
				},
				{
					field: "OnResults", signature: "() (r0 chan<- string)", index: "sResultsAt",
					results: []result{{
						name: "sResultsResults", label: "s-results", spelling: "chan<- string", end: "sResultsEnd",
						sink: true,
					}},
				},
				{
					field:     "OnPipe",
					signature: "() (r0 chan string)",
					index:     "sPipeAt",
					results: []result{
						{name: "sPipeResults", label: "s-pipe", spelling: "chan string", end: "sPipeEnd"},
					},
				},
			}, "a receive channel and a channel of both directions end closed, and a send channel ends at the sink")
		})

		inputs := []struct {
			name     string
			giveDecl string
			givePkgs []string
			want     [][2]string
		}{
			{
				name: "declares an implementation of each of two interfaces", giveDecl: "Router",
				givePkgs: []string{pathsPackage, inputsPackage},
				want: [][2]string{
					{"RouterDepotInput", "Depot"}, {"RouterLedgerInput", "Ledger"}, {"RouterRowsInput", "Rows"},
				},
			},
			{
				name: "declares the implementation of an interface that a function takes", giveDecl: "Route",
				givePkgs: []string{pathsPackage, inputsPackage}, want: [][2]string{{"RouteDepotInput", "Depot"}},
			},
			{
				name: "declares the implementation of each interface that the calls take at depth", giveDecl: "Courier",
				givePkgs: []string{nestedPackage},
				want:     [][2]string{{"CourierLabelerInput", "Labeler"}, {"CourierScaleInput", "Scale[int]"}},
			},
			{
				name:     "does not declare an implementation whose methods the file of the suite cannot import",
				giveDecl: "Tuned",
				givePkgs: []string{wrapPackage, toolPackage, deepPackage},
				want:     [][2]string{{"TunedNoteInput", "Note"}},
			},
		}
		for _, tt := range inputs {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := generate(t, goTree, tt.givePkgs...).suites[tt.giveDecl]
				names := make([][2]string, 0, len(got.Inputs))
				for _, in := range got.Inputs {
					names = append(names, [2]string{in.Name, spelled(in.Type)})
				}
				expect.Equal(t, names, tt.want, "the suite declares one implementation of each interface")
			})
		}

		t.Run("takes the declared errors of each method of an interface", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, pathsPackage, inputsPackage).suites["Router"], "Serve")
			assert.Length(t, got.Params, 2, "Serve takes a depot and a ledger")
			errs := make([][]string, 0, len(got.Params[0].Funcs))
			for _, f := range fnsOf(got.Params[0].Funcs) {
				errs = append(errs, f.results[len(f.results)-1].errors)
			}
			expect.Equal(t, errs, [][]string{{"ErrLost"}, {"ErrGone"}},
				"the errors that the directive fault lists replace the sentinels of the mixins, and a mixin "+
					"adds its sentinel alone")
		})

		t.Run("leaves out each declared error that the file of the suite cannot refer to", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, wrapPackage, toolPackage, deepPackage).suites["Limited"], "Sense")
			assert.Length(t, got.Params, 1, "Sense takes a sensor")
			errs := make([][]string, 0, len(got.Params[0].Funcs))
			for _, f := range fnsOf(got.Params[0].Funcs) {
				errs = append(errs, f.results[len(f.results)-1].errors)
			}
			expect.Equal(t, errs, [][]string{nil, {"ErrJammed"}, nil},
				"an error value of the internal package of the tree lib is no declared error of Sense, Probe or Read")
		})

		t.Run("takes no declared error from the sentinel of a contract", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, pathsPackage, inputsPackage).suites["Router"], "Read")
			assert.Length(t, got.Params, 1, "Read takes the rows")
			errs := make([][]string, 0, len(got.Params[0].Funcs))
			for _, f := range fnsOf(got.Params[0].Funcs) {
				errs = append(errs, f.results[len(f.results)-1].errors)
			}
			expect.Equal(t, errs, [][]string{nil, nil}, "the errors of a method come from its faults and its mixins")
		})

		t.Run("declares no implementation for calls without an interface", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, storePackage).suites["Store"]
			expect.Equal(t, got.Inputs, []suite.Input(nil), "the calls of Store take no interface")
		})
	})
}
