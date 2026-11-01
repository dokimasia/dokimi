// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
)

// sequencedPackage is the package of the Go tree that declares the
// classifications whose claims generated sequences check.
const sequencedPackage = "example.com/shop/sequenced"

func TestSequence(t *testing.T) {
	t.Parallel()

	generated := generate(t, goTree, sequencedPackage)
	wrapped := generate(t, goTree, wrapPackage, toolPackage, deepPackage).suites
	suiteOf := func(name string) *suite.Suite {
		s := generated.suites[name]
		return &s
	}
	defectOf := func(t *testing.T, s *suite.Suite, d *suite.Derived) *suite.Defect {
		t.Helper()
		at := slices.IndexFunc(s.Defects, func(o *suite.Defect) bool { return o.Name == d.Defect })
		assert.True(t, at >= 0, "the suite plants the defect "+d.Defect)
		return s.Defects[at]
	}
	stepOf := func(t *testing.T, s *suite.Suite, name string) *suite.Step {
		t.Helper()
		assert.NotNil(t, s.Sequences, "the suite of "+s.Declaration+" has generated sequences")
		at := slices.IndexFunc(s.Sequences.Steps, func(st suite.Step) bool { return st.Name == name })
		assert.True(t, at >= 0, "the sequences have a step of "+name)
		return &s.Sequences.Steps[at]
	}

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the count of an integer observer and the bounds to the probe of a bound", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Counter"), "Inc", "keeps-within-its-bounds")
			assert.Length(t, d.Args, 4, "the probe takes the actions, the observer, the floor and the ceiling")
			assert.NotNil(t, d.Args[1].Countable, "the observer is a call that returns the use")
			expect.Equal(t, d.Args[1].Countable.Count, "int", "Value returns the use as an integer")
			expect.Equal(t, d.Args[2].Expr.Word, "0", "the floor is the default of the param min")
			expect.Equal(t, d.Args[3].Expr.Word, "3", "the ceiling is the param limit")
			defect := defectOf(t, suiteOf("Counter"), d)
			expect.Equal(t, defect.Grow, "int", "the planted defect raises the integer of Value")
			expect.Equal(t, defect.Limit, "3", "the planted defect raises the use past the limit")
		})

		t.Run("passes the count of an unsigned observer to the probe of a bound", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Counter"), "Grow", "keeps-within-its-bounds")
			assert.NotNil(t, d.Args[1].Countable, "the observer is a call that returns the use")
			expect.Equal(t, d.Args[1].Countable.Count, "int", "Size returns the use as an unsigned integer")
		})

		t.Run("passes the length of a list observer to the probe of a bound", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Notes"), "Push", "keeps-within-its-bounds")
			assert.NotNil(t, d.Args[1].Countable, "the observer is a call that returns the use")
			expect.Equal(t, d.Args[1].Countable.Count, "len", "the use is the length of the list of Items")
			expect.Length(t, d.Args[1].Countable.Args, 1, "the call of Items passes the sample of its prefix")
			expect.Equal(t, defectOf(t, suiteOf("Notes"), d).Grow, "len", "the planted defect adds to the list")
		})

		t.Run("plants a defect that grows the text of a text observer", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Notes"), "Append", "keeps-within-its-bounds")
			expect.Equal(t, d.Args[1].Countable.Count, "len", "the use is the length of the text of Text")
			expect.Equal(t, defectOf(t, suiteOf("Notes"), d).Grow, "text", "the planted defect adds to the text")
		})

		t.Run("passes no argument for a variadic parameter of functions", func(t *testing.T) {
			t.Parallel()

			inc := stepOf(t, suiteOf("Counter"), "Inc")
			assert.Length(t, inc.Params, 2, "the step of Inc has the context and the options")
			expect.Equal(t, inc.Params[1].Arg, -1, "the call of Inc passes no option")
		})

		applies := []struct {
			name      string
			give      string
			giveField string
			giveCheck string
			want      []string
		}{
			{
				name: "applies the observer of a workflow with the sample of an input that the call lacks",
				give: "Run", giveField: "Advance", giveCheck: "follows-its-graph", want: []string{"args[0]", ""},
			},
			{
				name: "applies the observer with its own sample of an input of another type in the call",
				give: "Run", giveField: "Toggle", giveCheck: "follows-its-graph", want: []string{"args[0]", ""},
			},
			{
				name: "applies the backing with its own sample of an input of another name in the cache",
				give: "Cached", giveField: "Get", giveCheck: "cache-agrees-over-sequences", want: []string{""},
			},
			{
				name: "applies the backing with its own sample of an input of another type in the cache",
				give: "Cached", giveField: "Peek", giveCheck: "cache-agrees-over-sequences", want: []string{""},
			},
		}
		for _, tt := range applies {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				d := derivedNamed(t, suiteOf(tt.give), tt.giveField, tt.giveCheck)
				apply := d.Args[2].Apply
				assert.NotNil(t, apply, "the probe applies a callable to the arguments of "+tt.giveField)
				names := make([]string, 0, len(apply.Args))
				for _, a := range apply.Args {
					names = append(names, a.Name)
				}
				expect.Equal(t, names, tt.want, "the callable takes an argument of the call, or a sample for no name")
			})
		}

		typeArgs := []struct {
			name      string
			give      string
			giveField string
			giveCheck string
			want      []string
		}{
			{
				name: "passes the type of the value of a monotonic read to its probe", give: "Counter",
				giveField: "Value", giveCheck: "never-goes-back", want: []string{"int"},
			},
			{
				name: "passes the type of the position of an append to its probe", give: "Log",
				giveField: "Append", giveCheck: "appends-keep-ascending", want: []string{"int64"},
			},
			{
				name: "passes the type of the values of a stream of values and errors to its probe", give: "Feed",
				giveField: "All", giveCheck: "yields-each-once", want: []string{"string"},
			},
			{
				name: "passes the type of the elements of a stream to its probe", give: "Feed",
				giveField: "Each", giveCheck: "yields-each-once", want: []string{"string"},
			},
			{
				name: "passes the two types of the elements of a stream of pairs to its probe", give: "Feed",
				giveField: "Pairs", giveCheck: "yields-each-once", want: []string{"string", "int"},
			},
		}
		for _, tt := range typeArgs {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				d := derivedNamed(t, suiteOf(tt.give), tt.giveField, tt.giveCheck)
				assert.Length(t, d.TypeArgs, 1, "the probe takes one type argument after the subject type")
				spellings := make([]string, 0, len(d.TypeArgs[0]))
				for _, ref := range d.TypeArgs[0] {
					spellings = append(spellings, ref.Spelling)
				}
				expect.Equal(t, spellings, tt.want, "the type argument is the type of the values")
			})
		}

		t.Run("passes the reference to the probe of a merge of an interface of the declaration", func(t *testing.T) {
			t.Parallel()

			set := suiteOf("Set")
			d := derivedNamed(t, set, "Merge", "merges-in-any-order")
			assert.NotNil(t, set.Sequences, "the suite has generated sequences")
			expect.Equal(t, set.Sequences.Reference, "SetReference", "the merges apply their operations by name")
			expect.Equal(t, d.Args[1].Ref, set.Sequences.Reference, "the probe applies the merge by its name")
		})

		t.Run("binds a merge of another subject of the declaration", func(t *testing.T) {
			t.Parallel()

			derivedNamed(t, suiteOf("Set"), "Join", "merges-in-any-order")
		})

		t.Run("binds a merge of a pointer to the struct", func(t *testing.T) {
			t.Parallel()

			derivedNamed(t, suiteOf("Tally"), "Merge", "merges-in-any-order")
		})

		t.Run("binds the check of a declaration to its reference", func(t *testing.T) {
			t.Parallel()

			register := suiteOf("Register")
			d := derivedNamed(t, register, "Get", "agrees-with-the-reference")
			expect.Equal(t, d.Whole, "Register/agrees-with-the-reference", "the identity names the declaration")
			assert.NotNil(t, register.Sequences, "the suite has generated sequences")
			expect.Equal(t, d.Args[1].Ref, register.Sequences.Reference, "the probe takes the reference")
		})

		t.Run("names the call of a function after the function", func(t *testing.T) {
			t.Parallel()

			double := suiteOf("Double")
			derivedNamed(t, double, "", "agrees-with-the-reference")
			expect.Equal(t, stepOf(t, double, "Double").Call, "DoubleCall", "the call of Double is DoubleCall")
		})

		t.Run("drops the reference of a declaration whose checks apply no operation by its name", func(t *testing.T) {
			t.Parallel()

			counter := suiteOf("Counter")
			assert.NotNil(t, counter.Sequences, "the suite has generated sequences")
			expect.Equal(t, counter.Sequences.Reference, "", "no check of Counter applies an operation by its name")
		})

		t.Run("applies the read with the own sample of its function input", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Index"), "Find", "reads-each-last-write")
			readback := d.Args[2].Apply
			assert.NotNil(t, readback, "the readback applies Find to the arguments of Put")
			assert.Length(t, readback.Args, 2, "Find takes the key and keep")
			expect.Equal(t, readback.Args[0].Select, "Key", "Find takes the field Key of the written item")
			expect.NotNil(t, readback.Args[1].Value, "Find takes the sample of keep")
		})

		t.Run("applies a read of two fields of the written value with one written value", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Index"), "Lookup", "reads-each-last-write")
			readback := d.Args[2].Apply
			assert.NotNil(t, readback, "the readback applies Lookup to the arguments of Put")
			expect.Length(t, readback.Lets, 1, "Lookup reads the written item once")
			expect.Equal(t, readback.Args[1].Select, "Rev", "Lookup takes the field Rev of the written item")
		})

		t.Run("gives the written value the key of each pool that a read selects it by", func(t *testing.T) {
			t.Parallel()

			put := stepOf(t, suiteOf("Index"), "Put")
			fields := make([]string, 0, len(put.Params[1].Fields))
			for _, f := range put.Params[1].Fields {
				fields = append(fields, f.Name)
			}
			expect.Equal(t, fields, []string{"Key", "Rev"}, "the reads select the item by its key and its revision")
		})

		t.Run("gives the written value the key of a pair of a check of a single caller", func(t *testing.T) {
			t.Parallel()

			store := stepOf(t, suiteOf("Stamped"), "Store")
			assert.Length(t, store.Params[1].Fields, 1, "the write of the pair of reads-never-regress takes a key")
			expect.Equal(t, store.Params[1].Fields[0].Name, "Key", "Get selects the item by its key")
		})

		t.Run("passes the values of each call through the function of the param normalize", func(t *testing.T) {
			t.Parallel()

			register := suiteOf("Register")
			assert.NotNil(t, register.Sequences, "the suite has generated sequences")
			assert.NotNil(t, register.Sequences.Normalize, "the directive suite of Register names Unstamped")
			expect.Equal(t, register.Sequences.Normalize.Spelling, "Unstamped", "the function is Unstamped")
		})

		t.Run("reports each function of the param normalize of another type", func(t *testing.T) {
			t.Parallel()

			var reported []string
			for _, d := range generated.findings {
				if d.Code == suite.BadNormalize {
					reported = append(reported, d.Msg)
				}
			}
			const tail = ", and the function of the param has the type func(string, []any) []any, so the suite " +
				"records the values of each call as they are"
			expect.Permutation(t, reported, []string{
				"the function Counted of the param normalize of Counter has the type func([]any) []any" + tail,
				"the function Noted of the param normalize of Notes has the type func(string, []any)" + tail,
				"the function Logged of the param normalize of Log has the type func(string, []any) ([]any, error)" + tail,
			}, "each function of another type is a finding at its directive")
			expect.Nil(t, suiteOf("Counter").Sequences.Normalize, "the suite of Counter passes the values as they are")
		})

		t.Run("writes no sequences for a declaration whose checks of sequences the suite refuses", func(t *testing.T) {
			t.Parallel()

			expect.Nil(t, suiteOf("Pile").Sequences, "the suite refuses the merges of Pile")
		})

		t.Run("binds the graph of a workflow whose variable has a stated map type", func(t *testing.T) {
			t.Parallel()

			derivedNamed(t, suiteOf("Run"), "Retry", "follows-its-graph")
			derivedNamed(t, suiteOf("Run"), "Advance", "follows-its-graph")
		})

		refusals := []struct {
			name     string
			give     *suite.Suite
			giveFrom string
			want     string
		}{
			{
				name: "refuses a bound whose observer returns no use", give: suiteOf("Notes"), giveFrom: "Pin",
				want: "keeps-within-its-bounds",
			},
			{
				name: "refuses a bound whose observer returns no value", give: suiteOf("Notes"), giveFrom: "Pop",
				want: "keeps-within-its-bounds",
			},
			{
				name: "refuses a bound whose observer returns a fraction", give: suiteOf("Counter"),
				giveFrom: "Throttle", want: "keeps-within-its-bounds",
			},
			{
				name: "refuses a merge whose input no other subject can be", give: suiteOf("Pile"), giveFrom: "Absorb",
				want: "merges-in-any-order",
			},
			{
				name: "refuses a merge of an interface that the declaration does not implement", give: suiteOf("Pile"),
				giveFrom: "Pour", want: "merges-in-any-order",
			},
			{
				name: "refuses a merge of two inputs", give: suiteOf("Set"), giveFrom: "Weigh",
				want: "merges-in-any-order",
			},
			{
				name: "refuses a workflow whose transitions are no graph", give: suiteOf("Run"), giveFrom: "Restart",
				want: "follows-its-graph",
			},
			{
				name: "refuses the reads of sequences of a write that returns no value", give: suiteOf("Versioned"),
				giveFrom: "Get", want: "writes-follow-their-reads",
			},
			{
				name: "refuses a check whose apply passes a sample that the suite cannot import",
				give: new(wrapped["Played"]), giveFrom: "Play", want: "follows-its-graph",
			},
			{
				name: "refuses a check whose observer takes a sample that the suite cannot import",
				give: new(wrapped["Gauged"]), giveFrom: "Fill", want: "keeps-within-its-bounds",
			},
		}
		for _, tt := range refusals {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				refused := slices.ContainsFunc(tt.give.Derivation.Refusals, func(r derive.Refusal) bool {
					return r.Callable == tt.giveFrom && string(r.Claim) == tt.want
				})
				expect.True(t, refused, "the suite refuses the claim "+tt.want+" of "+tt.giveFrom)
				derived := slices.ContainsFunc(tt.give.Derived, func(d suite.Derived) bool {
					return d.Field == tt.giveFrom && strings.HasPrefix(d.Segment, tt.want)
				})
				expect.False(t, derived, "the suite derives no check of the refused claim")
			})
		}
	})
}
