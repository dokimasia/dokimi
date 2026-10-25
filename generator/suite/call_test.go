// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/suite"
)

// The call generator of each callable of a suite draws the inputs of a call
// from the shapes of their types, and the planted defect of a ceiling takes
// the signature of the callable.
func TestCall(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("records a call of each callable in the order of the member walk", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, meterPackage).suites["Meter"]
			names := make([][3]string, 0, len(got.Calls))
			for _, c := range got.Calls {
				names = append(names, [3]string{c.Var, c.Field, c.Type})
			}
			expect.Equal(t, names, [][3]string{
				{"MeterCountCalls", "Count", "MeterCount"},
				{"MeterRecordCalls", "Record", "MeterRecord"},
				{"MeterMergeCalls", "Merge", "MeterMerge"},
				{"MeterTotalCalls", "Total", "MeterTotal"},
			}, "each call joins Calls to the type of its callable")
		})

		t.Run("records the call of a function without a field", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, meterPackage).suites["Weigh"]
			assert.Length(t, got.Calls, 1, "a function is the one callable of its suite")
			expect.Equal(t, [3]string{got.Calls[0].Var, got.Calls[0].Field, got.Calls[0].Type},
				[3]string{"WeighCalls", "", "WeighCallable"}, "the call of a function calls the subject itself")
		})

		tests := []struct {
			name        string
			givePkg     string
			giveDecl    string
			give        string
			wantParams  []param
			wantAssign  string
			wantReturns string
		}{
			{
				name: "draws the value of each input after the context", givePkg: meterPackage, giveDecl: "Meter",
				give:       "Count",
				wantParams: []param{{context: true}, {name: "unit", label: "unit", spelling: "Unit"}},
				wantAssign: "_, err :=", wantReturns: "err",
			},
			{
				name: "draws a list of the element type for a variadic parameter", givePkg: meterPackage,
				giveDecl: "Meter", give: "Record",
				wantParams: []param{
					{context: true}, {name: "samples", label: "samples", spelling: "Sample", variadic: true},
				},
			},
			{
				name: "draws a pointer that can be nil under the mixin nil-safe", givePkg: meterPackage,
				giveDecl: "Meter", give: "Merge",
				wantParams:  []param{{name: "s", label: "s", spelling: "*Sample"}},
				wantReturns: "nil",
			},
			{
				name: "draws no input for a callable without parameters", givePkg: meterPackage, giveDecl: "Meter",
				give: "Total", wantAssign: "_ =", wantReturns: "nil",
			},
			{
				name: "draws the value of a pointer without the mixin nil-safe", givePkg: callsPackage,
				giveDecl: "Board", give: "Place",
				wantParams: []param{{name: "p", label: "p", spelling: "Point", address: true}},
			},
			{
				name: "renames an input with a name that the generated call takes", givePkg: callsPackage,
				giveDecl: "Board", give: "Clash",
				wantParams: []param{
					{context: true},
					{name: "c0", label: "c", spelling: "string"},
					{name: "subject0", label: "subject", spelling: "string"},
					{name: "err0", label: "err", spelling: "string"},
					{name: "prop0", label: "prop", spelling: "int"},
				},
			},
			{
				name: "names each input without a name by its position", givePkg: callsPackage, giveDecl: "Board",
				give: "Blank",
				wantParams: []param{
					{name: "a0", label: "a0", spelling: "string"}, {name: "a1", label: "a1", spelling: "int"},
				},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name of a predeclared type", givePkg: callsPackage,
				giveDecl: "Board", give: "Label",
				wantParams:  []param{{name: "string0", label: "string", spelling: "string"}},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name of the package of an element type", givePkg: callsPackage,
				giveDecl: "Board", give: "Wait",
				wantParams: []param{
					{name: "time0", label: "time", spelling: "int"},
					{name: "delays", label: "delays", spelling: "[]time.Duration"},
				},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name of the package of a type argument", givePkg: callsPackage,
				giveDecl: "Board", give: "Hold",
				wantParams: []param{
					{name: "time0", label: "time", spelling: "int"},
					{name: "pair", label: "pair", spelling: "Pair[time.Month]"},
				},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name of the package of the declaration", givePkg: callsPackage,
				giveDecl: "Board", give: "Repeat",
				wantParams:  []param{{name: "calls0", label: "calls", spelling: "int"}},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name under which the file imports a package", givePkg: callsPackage,
				giveDecl: "Board", give: "Stock",
				wantParams: []param{
					{name: "shop0", label: "shop", spelling: "int"},
					{name: "item", label: "item", spelling: "shop.Item"},
				},
				wantReturns: "nil",
			},
			{
				name: "names an input with the blank name by its position", givePkg: callsPackage, giveDecl: "Board",
				give: "Discard", wantParams: []param{{name: "a0", label: "a0", spelling: "string"}},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name that an input before it takes", givePkg: callsPackage,
				giveDecl: "Board", give: "Twice",
				wantParams: []param{
					{name: "c0", label: "c", spelling: "string"}, {name: "c00", label: "c0", spelling: "string"},
				},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name of the list of a channel before it", givePkg: callsPackage,
				giveDecl: "Board", give: "Drain",
				wantParams: []param{
					{name: "intsValues", label: "ints", spelling: "int"},
					{name: "intsValues0", label: "ints-values", spelling: "int"},
				},
				wantReturns: "nil",
			},
			{
				name: "renames an input with the name of a part of a generated function", givePkg: callsPackage,
				giveDecl: "Board", give: "Gather",
				wantParams: []param{
					{name: "fn", label: "fn", spelling: "func() int"},
					{name: "fnResults0", label: "fn-results", spelling: "int"},
					{name: "fnAt0", label: "fn-at", spelling: "int"},
				},
				wantReturns: "nil",
			},
			{
				name: "draws the inputs of a generic interface at its witnesses", givePkg: genericPackage,
				giveDecl: "Cache", give: "Get",
				wantParams: []param{{name: "key", label: "key", spelling: "int"}}, wantAssign: "_, _ =",
				wantReturns: "nil",
			},
			{
				name: "draws the inputs of a method of a struct", givePkg: catalogPackage, giveDecl: "Catalog",
				give: "Add",
				wantParams: []param{
					{name: "id", label: "id", spelling: "string"}, {name: "name", label: "name", spelling: "string"},
				},
				wantReturns: "nil",
			},
			{
				name: "draws a pointer to a zone as the pointer that the property engine reads", givePkg: loosePackage,
				giveDecl: "Zoned", give: "In",
				wantParams:  []param{{name: "loc", label: "loc", spelling: "*time.Location"}},
				wantReturns: "nil",
			},
			{
				name:    "draws a pointer to a big integer as the pointer that the property engine reads",
				givePkg: loosePackage, giveDecl: "Zoned", give: "Scale",
				wantParams:  []param{{name: "n", label: "n", spelling: "*big.Int"}},
				wantReturns: "nil",
			},
			{
				name:    "draws a pointer to a big rational as the pointer that the property engine reads",
				givePkg: loosePackage, giveDecl: "Zoned", give: "Ratio",
				wantParams:  []param{{name: "r", label: "r", spelling: "*big.Rat"}},
				wantReturns: "nil",
			},
			{
				name: "draws the inputs of a function", givePkg: meterPackage, giveDecl: "Weigh", give: "",
				wantParams: []param{
					{name: "s", label: "s", spelling: "Sample"},
					{name: "density", label: "density", spelling: "float64"},
				},
				wantAssign: "_ =", wantReturns: "nil",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := called(t, generate(t, goTree, tt.givePkg).suites[tt.giveDecl], tt.give)
				params := make([]param, 0, len(got.Params))
				for _, p := range got.Params {
					params = append(params, paramOf(p))
				}
				expect.Equal(t, params, tt.wantParams, "the call draws each input of the callable",
					assert.EquateEmpty())
				expect.Equal(t, got.Assign, tt.wantAssign, "the call assigns the results of the callable")
				expect.Equal(t, got.Returns, tt.wantReturns, "the call returns the error of the callable or nil")
			})
		}

		t.Run("renames an input with the name that the package of a type declares", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, callsPackage, storePackage).suites["Board"], "Keep")
			params := make([]param, 0, len(got.Params))
			for _, p := range got.Params {
				params = append(params, paramOf(p))
			}
			expect.Equal(t, params, []param{
				{name: "store0", label: "store", spelling: "int"}, {name: "item", label: "item", spelling: "shop.Item"},
			}, "the generated file imports the package of the item under the name store")
		})

		budgets := []struct {
			name       string
			giveDecl   string
			give       string
			wantBudget suite.Budget
		}{
			{
				name: "records the budget of a method", giveDecl: "Meter", give: "Count",
				wantBudget: suite.Budget{
					Allocs: new(int64(1)), Bytes: new(int64(96)), Latency: 2 * time.Microsecond,
					Mean: time.Microsecond,
				},
			},
			{
				name: "records the budget of a function", giveDecl: "Weigh",
				wantBudget: suite.Budget{Allocs: new(int64(0)), Latency: 100 * time.Nanosecond},
			},
			{
				name: "records the warm-up of a method", giveDecl: "Meter", give: "Record",
				wantBudget: suite.Budget{Warmup: new(int64(500)), Latency: time.Millisecond},
			},
			{name: "records no budget of a callable without one", giveDecl: "Meter", give: "Total"},
		}
		for _, tt := range budgets {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := called(t, generate(t, goTree, meterPackage).suites[tt.giveDecl], tt.give)
				expect.Equal(t, got.Budget, tt.wantBudget,
					"the call has the budget that the facts of its callable state")
			})
		}

		concurrency := []struct {
			name     string
			giveDecl string
			give     string
			want     bool
		}{
			{name: "reports the mixin concurrent of a method", giveDecl: "Meter", give: "Count", want: true},
			{name: "reports the mixin concurrent of a function", giveDecl: "Weigh", want: true},
			{name: "reports no concurrency of a method without the mixin", giveDecl: "Meter", give: "Record"},
		}
		for _, tt := range concurrency {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := called(t, generate(t, goTree, meterPackage).suites[tt.giveDecl], tt.give)
				expect.Equal(t, got.Concurrent, tt.want, "a parallel benchmark measures a concurrent callable")
			})
		}

		t.Run("reports the mixin concurrent-readers of a method", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, pathsPackage, inputsPackage).suites["Router"], "Scan")
			expect.True(t, got.Concurrent, "a parallel benchmark measures a callable with concurrent readers")
		})

		defects := []struct {
			name     string
			givePkg  string
			giveDecl string
			give     string
			want     string
		}{
			{
				name: "writes the defect of a method with its parameters without names", givePkg: meterPackage,
				giveDecl: "Meter", give: "Count", want: "Count(context.Context, Unit) (r0 int, r1 error)",
			},
			{
				name: "writes the defect of a variadic method with its list", givePkg: meterPackage, giveDecl: "Meter",
				give: "Record", want: "Record(context.Context, ...Sample) (r0 error)",
			},
			{
				name: "writes the defect of a method without results", givePkg: callsPackage, giveDecl: "Board",
				give: "Blank", want: "Blank(string, int) ()",
			},
			{
				name: "writes the defect of a generic interface at its witnesses", givePkg: genericPackage,
				giveDecl: "Cache", give: "Get", want: "Get(int) (r0 int, r1 bool)",
			},
			{
				name: "writes the defect of a function", givePkg: meterPackage, giveDecl: "Weigh",
				want: "Weigh(Sample, float64) (r0 float64)",
			},
		}
		for _, tt := range defects {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := called(t, generate(t, goTree, tt.givePkg).suites[tt.giveDecl], tt.give)
				assert.NotNil(t, got.Over, "the call has the signature of its planted defect")
				expect.Equal(t, signature(got.Over), tt.want,
					"the defect takes the types of the parameters and names each result")
			})
		}

		funcs := []struct {
			name     string
			give     string
			wantName string
			want     []fn
		}{
			{
				name: "draws a generated function for a parameter of a function type", give: "Each", wantName: "fn",
				want: []fn{{
					signature: "(Item) (r0 bool)", index: "fnAt",
					results: []result{{name: "fnResults", label: "fn", spelling: "bool"}},
				}},
			},
			{
				name: "draws a generated function for a parameter of a defined function type", give: "Handle",
				wantName: "h",
				want: []fn{{
					signature: "(Item) (r0 error)", index: "hAt",
					results: []result{{name: "hResults", label: "h", spelling: "error", isError: true}},
				}},
			},
			{
				name: "draws a list of each result of a generated function", give: "Do", wantName: "op",
				want: []fn{{
					signature: "() (r0 int, r1 error)", index: "opAt",
					results: []result{
						{name: "opResults0", label: "op-0", spelling: "int"},
						{name: "opResults1", label: "op-1", spelling: "error", isError: true},
					},
				}},
			},
			{
				name: "draws no list for a generated function without results", give: "Notify", wantName: "fn",
				want: []fn{{signature: "(string) ()"}},
			},
			{
				name: "draws a generated function of a variadic function type", give: "Apply", wantName: "fn",
				want: []fn{{signature: "(...int) ()"}},
			},
			{
				name: "draws a generated function for each method of an interface", give: "Pump", wantName: "from",
				want: []fn{
					{
						field: "OnNext", signature: "(context.Context) (r0 Item, r1 bool, r2 error)",
						index: "fromNextAt",
						results: []result{
							{name: "fromNextResults0", label: "from-next-0", spelling: "Item"},
							{name: "fromNextResults1", label: "from-next-1", spelling: "bool"},
							{
								name: "fromNextResults2", label: "from-next-2", spelling: "error", isError: true,
								errors: []string{"ErrDrained"},
							},
						},
					},
					{
						field: "OnGet", signature: "(context.Context, string) (r0 Item, r1 error)", index: "fromGetAt",
						results: []result{
							{name: "fromGetResults0", label: "from-get-0", spelling: "Item"},
							{
								name: "fromGetResults1", label: "from-get-1", spelling: "error", isError: true,
								errors: []string{"ErrMissing"},
							},
						},
					},
					{field: "OnSeek", signature: "(...string) ()"},
					{
						field: "OnClose", signature: "() (r0 error)", index: "fromCloseAt",
						results: []result{
							{name: "fromCloseResults", label: "from-close", spelling: "error", isError: true},
						},
					},
				},
			},
		}
		for _, tt := range funcs {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := called(t, generate(t, goTree, inputsPackage).suites["Pipeline"], tt.give)
				p := got.Params[len(got.Params)-1]
				expect.Equal(t, p.Name, tt.wantName, "the generated function takes the name of its parameter")
				expect.Equal(t, fnsOf(p.Funcs), tt.want, "the call draws the results of each generated function")
			})
		}

		t.Run("draws an implementation of an interface outside the standard library", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, inputsPackage).suites["Pipeline"], "Pump")
			assert.Length(t, got.Params, 2, "Pump takes a context and a feed")
			expect.Equal(t, got.Params[1].Input, "PipelineFeedInput", "the feed takes the implementation of Feed")
		})

		omitted := []struct {
			name string
			give string
		}{
			{name: "passes no values for a variadic parameter of functions", give: "Configure"},
			{name: "passes no values for a variadic parameter of an interface", give: "Subscribe"},
		}
		for _, tt := range omitted {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := called(t, generate(t, goTree, inputsPackage).suites["Pipeline"], tt.give)
				last := got.Params[len(got.Params)-1]
				expect.True(t, last.Omitted, "no generator draws a value of the elements")
				expect.Empty(t, last.Funcs, "the call draws no function for the parameter")
			})
		}

		t.Run("draws a list for a channel that the callable receives from", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, inputsPackage).suites["Pipeline"], "Fold")
			assert.Length(t, got.Params, 2, "Fold takes a context and a channel")
			p := got.Params[1]
			expect.Equal(t, paramOf(p), param{name: "itemsValues", label: "items", spelling: "Item"},
				"the call draws a list of the elements of the channel")
			expect.True(t, p.Receive, "each call passes a new channel of the list")
		})

		t.Run("draws the channel itself for a channel that the callable sends to", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, pathsPackage, inputsPackage).suites["Router"], "Ship")
			assert.Length(t, got.Params, 1, "Ship takes a channel")
			expect.Equal(t, paramOf(got.Params[0]), param{name: "out", label: "out", spelling: "chan<- string"},
				"the call draws a value of the channel")
			expect.False(t, got.Params[0].Receive, "a channel that the callable sends to takes no list")
		})

		t.Run("passes the elements of a variadic function type as its list", func(t *testing.T) {
			t.Parallel()

			got := called(t, generate(t, goTree, inputsPackage).suites["Pipeline"], "Apply")
			assert.Length(t, got.Params, 1, "Apply takes a function")
			assert.Length(t, got.Params[0].Funcs, 1, "the function type has one generated function")
			params := got.Params[0].Funcs[0].Signature.Params.Items()
			assert.Length(t, params, 1, "the function type takes one list")
			expect.Equal(t, params[0].Variadic, symbol.VariadicPositional, "the list is variadic")
			expect.Equal(t, spelled(params[0].Type), "int", "the list takes the type of its elements")
		})
	})
}

// param is a parameter of a call with the spelling of its type.
type param struct {
	name, label, spelling      string
	context, variadic, address bool
}

// paramOf returns the parameter p with the spelling of its type and of the
// arguments of the type.
func paramOf(p suite.Param) param {
	out := param{name: p.Name, label: p.Label, context: p.Context, variadic: p.Variadic, address: p.Address}
	if p.Type != nil {
		out.spelling = spelled(p.Type)
	}
	return out
}
