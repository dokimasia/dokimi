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

		t.Run("declares no implementation for calls without an interface", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, storePackage).suites["Store"]
			expect.Equal(t, got.Inputs, []suite.Input(nil), "the calls of Store take no interface")
		})
	})
}
