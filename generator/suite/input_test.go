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
			expect.Equal(t, got.Params[0].Input, "", "a reader of the standard library takes its registration")
			expect.Equal(t, spelled(got.Params[0].Type), "io.Reader", "the call draws a value of the reader")
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
