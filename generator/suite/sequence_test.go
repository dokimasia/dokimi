// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/spec"

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

		vaults := generate(t, goTree, vaultPackage).suites

		t.Run("keeps the steps and the reference of a world without the actions of sequences", func(t *testing.T) {
			t.Parallel()

			v := vaults["Vault"]
			assert.NotNil(t, v.Sequences, "the world of Vault calls the steps of its sequences")
			expect.Equal(t, []string{v.Sequences.Actions, v.Sequences.Reference}, []string{"", "VaultReference"},
				"no check of Vault runs the actions, and the world applies its history to the reference")
			expect.Length(t, v.Sequences.Steps, 2, "the world calls Put and Get")
		})

		t.Run("keeps the actions of a world beside a check of generated sequences", func(t *testing.T) {
			t.Parallel()

			r := vaults["Registry"]
			assert.NotNil(t, r.Sequences, "the suite of Registry has generated sequences")
			expect.Equal(t, r.Sequences.Actions, "RegistryActions", "agrees-with-the-reference runs the actions")
		})

		t.Run("adds the first value of a call of the mixin creates to the pool of its type", func(t *testing.T) {
			t.Parallel()

			r := vaults["Registry"]
			p := stepOf(t, &r, "Add").Produced
			assert.NotNil(t, p, "a call of Add produces a key")
			expect.Equal(t, []any{p.Pool, p.Type.Spelling, p.Field, p.Pointer},
				[]any{"keysString", "string", "", false},
				"the pool of the strings takes the key that Add returns")
		})

		t.Run("adds the field of the first value that the param key of the mixin creates names", func(t *testing.T) {
			t.Parallel()

			r := vaults["Registry"]
			p := stepOf(t, &r, "Insert").Produced
			assert.NotNil(t, p, "a call of Insert produces a key")
			expect.Equal(t, []any{p.Pool, p.Type.Spelling, p.Field, p.Pointer},
				[]any{"keysString", "*Entry", "Key", true},
				"the pool of the strings takes the field Key of the entry that Insert returns")
		})

		produced := []struct {
			name string
			give string
		}{
			{name: "produces no key for a callable without the mixin creates", give: "Set"},
			{name: "produces no key of a type that no pool has", give: "Stamp"},
			{name: "produces no key for a callable of the mixin creates without a result", give: "Touch"},
		}
		for _, tt := range produced {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				r := vaults["Registry"]
				expect.Nil(t, stepOf(t, &r, tt.give).Produced, "the world adds no key of "+tt.give+" to a pool")
			})
		}

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

		t.Run("gives no key of a pool to a write whose builder returns its inputs", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, sealedPackage).suites["Shelf"]
			put := stepOf(t, &shelf, "Put")
			assert.Length(t, put.Params, 3, "Put takes the context, the key and the number")
			expect.Equal(t, put.Params[1].Pool, "", "Put takes its key from its builder")
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

		t.Run("passes the total that observe reads to the probe of a conserved quantity", func(t *testing.T) {
			t.Parallel()

			accounts := suiteOf("Accounts")
			d := derivedNamed(t, accounts, "Transfer", "conserves-the-quantity")
			assert.Length(t, d.Args, 3, "the probe takes the actions, the operation and the total")
			expect.Equal(t, d.Args[1].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Transfer"))},
				"the operation is the name of Transfer")
			total := d.Args[2].Closure
			assert.NotNil(t, total, "the total is a closure")
			expect.Equal(t, total.Method, "Total", "the total calls the reader that observe names")
			assert.Length(t, total.Returns, 1, "the total is one value")
			expect.Equal(t, spelled(total.Returns[0]), "int64", "the total has the type of the value of Total")
			defect := defectOf(t, accounts, d)
			expect.Equal(t, defect.Body, "leak", "the planted defect adds a unit to the total on each call")
			expect.Equal(t, defect.Method.Name, "Total", "the planted defect changes the reader of the total")
			expect.True(t, defect.Counts, "the planted defect counts the calls of Total")
		})

		t.Run("applies the read of the mixin atomic to the arguments of a call that fails", func(t *testing.T) {
			t.Parallel()

			journal := suiteOf("Journal")
			d := derivedNamed(t, journal, "Record", "no-trace-of-a-failed-call")
			assert.Length(t, d.Args, 3, "the probe takes the actions, the operation and the readback")
			readback := d.Args[2].Apply
			assert.NotNil(t, readback, "the readback applies Entry to the arguments of Record")
			assert.Length(t, readback.Args, 1, "Entry takes the key")
			expect.Equal(t, readback.Args[0].Select, "Key", "Entry takes the field Key of the recorded item")
			defect := defectOf(t, journal, d)
			expect.Equal(t, defect.Body, "half-applied-after", "the planted defect reads a sample after a failure")
			expect.Equal(t, defect.Method.Name, "Entry", "the planted defect changes the read")
			expect.Equal(t, defect.After, 3, "the planted defect breaks from the fourth read on")
			expect.True(t, defect.Counts, "the planted defect counts the reads")
			expect.True(t, defect.Keeps, "the planted defect reads the injector of its handle")
		})

		pulls := []struct {
			name      string
			give      string
			wantDrain string
			wantElems []string
		}{
			{
				name: "pulls the values of a stream of values and errors", give: "Items", wantDrain: "until-error",
				wantElems: []string{"string", "error"},
			},
			{
				name: "pulls the elements of a stream", give: "Each", wantDrain: "elements",
				wantElems: []string{"string"},
			},
			{
				name: "pulls the first value of each pair of a stream of pairs", give: "Counts", wantDrain: "pairs",
				wantElems: []string{"string", "int"},
			},
		}
		for _, tt := range pulls {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				stock := suiteOf("Stock")
				d := derivedNamed(t, stock, tt.give, "iteration-sees-mutations")
				assert.Length(t, d.Args, 6,
					"the probe takes the actions, the stream, two operations and two projections")
				stream := d.Args[1].Closure
				assert.NotNil(t, stream, "the stream is a closure")
				expect.True(t, stream.Pull, "the closure starts a pull iteration of the stream")
				expect.Equal(t, stream.Drain, tt.wantDrain, "the iteration pulls the stream by its drain")
				expect.Equal(t, spelled(stream.Returns[0]), "string", "the iteration yields the items")
				expect.Equal(t, d.Args[2].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Add"))},
					"the mutation is Add")
				expect.Equal(t, d.Args[3].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Remove"))},
					"the delete is Remove")
				for _, a := range d.Args[4:] {
					assert.NotNil(t, a.Written, "each projection returns an argument of a call")
					expect.Equal(t, a.Written.Index, 0, "Add and Remove take the item first")
				}
				defect := defectOf(t, stock, d)
				expect.Equal(t, defect.Body, "copy", "the planted defect iterates a copy of the stream")
				expect.Equal(t, defect.Stream, "r0", "the stream is the first result")
				elems := make([]string, 0, len(defect.Elems))
				for _, e := range defect.Elems {
					elems = append(elems, spelled(e))
				}
				expect.Equal(t, elems, tt.wantElems, "the copy keeps each value of each element")
			})
		}

		deliveries := []struct {
			name          string
			give          string
			giveField     string
			giveCheck     string
			wantSubscribe string
			wantBody      string
			wantWritten   []int
			wantArgs      int
			wantAt        int
		}{
			{
				name: "passes the channel of the subscription of an outbox to its probe", give: "Outbox",
				giveField: "Append", giveCheck: "delivers-each-append", wantSubscribe: "Subscribe",
				wantBody: "zero-after", wantWritten: []int{3}, wantArgs: 4, wantAt: 2,
			},
			{
				name: "passes the channel of a publisher at least once to its probe", give: "Broker",
				giveField: "Send", giveCheck: "delivers-by-its-mode", wantSubscribe: "Listen",
				wantBody: "zero-after", wantWritten: []int{3}, wantArgs: 5, wantAt: 2,
			},
			{
				name: "plants a second delivery of each event of a publisher at most once", give: "Broker",
				giveField: "Fire", giveCheck: "delivers-by-its-mode", wantSubscribe: "Tap",
				wantBody: "twice-after", wantWritten: []int{3}, wantArgs: 5, wantAt: 2,
			},
			{
				name: "passes the channel of the subscription of a redelivery to its probe", give: "Broker",
				giveField: "Resend", giveCheck: "redelivers-on-demand", wantSubscribe: "Listen",
				wantBody: "zero-after", wantWritten: []int{4, 5}, wantArgs: 6, wantAt: 3,
			},
		}
		for _, tt := range deliveries {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				s := suiteOf(tt.give)
				d := derivedNamed(t, s, tt.giveField, tt.giveCheck)
				assert.Length(t, d.Args, tt.wantArgs, "the probe takes the actions, the operations, the channel "+
					"and the projections")
				subscribe := d.Args[tt.wantAt].Closure
				assert.NotNil(t, subscribe, "the subscription is a closure")
				expect.True(t, subscribe.Channel, "the closure returns the channel of the subscription")
				expect.Equal(t, subscribe.Method, tt.wantSubscribe, "the closure calls the role subscribe")
				expect.Equal(t, spelled(subscribe.Returns[0]), "Event", "the channel delivers events")
				for _, i := range tt.wantWritten {
					assert.NotNil(t, d.Args[i].Written, "the projection returns an argument of a call")
					expect.Equal(t, d.Args[i].Written.Index, 0, "each call passes the event first")
				}
				expect.True(t, d.Bubble, "the probe waits for the deliveries inside a bubble")
				expect.Equal(t, defectOf(t, s, d).Body, tt.wantBody,
					"the planted defect breaks from the fourth call on")
			})
		}

		t.Run("passes the delivery mode of a publisher to its probe", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Broker"), "Fire", "delivers-by-its-mode")
			expect.Equal(t, d.Args[4].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "at-most-once"))},
				"the mode is the text of the param mode")
		})

		t.Run("passes the methods of a pointer to a subscription to the probe of a watch", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, suiteOf("Watches"), "Ticks", "watch-sees-each-trigger")
			assert.NotNil(t, d.Args[3].Method, "next is a method of the subscription")
			expect.Equal(t, spelled(d.Args[3].Method.Type), "*Ticker",
				"the method belongs to the pointer that Ticks returns")
			assert.NotNil(t, d.Args[7].Key, "the key of a trigger is a projection")
			expect.Equal(t, d.Args[7].Key.Args, []int{0},
				"the key of Tick leaves out the note that Ticks does not take")
		})

		t.Run("passes the methods of the subscription and the keys of a watch to its probe", func(t *testing.T) {
			t.Parallel()

			watches := suiteOf("Watches")
			d := derivedNamed(t, watches, "Watch", "watch-sees-each-trigger")
			assert.Length(t, d.Args, 9, "the probe takes the actions, two operations, two methods, the retrigger, "+
				"two keys and the update")
			expect.Equal(t, d.Args[2].Expr, &suite.Expr{Value: new(emit.Literal(emit.LiteralString, "Set"))},
				"the trigger is Set")
			methods := make([]string, 0, 2)
			for _, a := range d.Args[3:5] {
				assert.NotNil(t, a.Method, "next and stop are methods of the subscription")
				expect.Equal(t, spelled(a.Method.Type), "Subscription", "the subscription is the value of Watch")
				methods = append(methods, a.Method.Name)
			}
			expect.Equal(t, methods, []string{"Next", "Stop"}, "the params next and stop name the methods")
			assert.NotNil(t, d.Args[5].Apply, "the retrigger applies Set")
			expect.Length(t, d.Args[5].Apply.Args, 3, "Set takes the arguments of its own call")
			assert.NotNil(t, d.Args[6].Key, "the key of a watch is a projection")
			expect.Equal(t, d.Args[6].Key.Args, []int{1, 0}, "the key of Watch lists its key before its scope")
			assert.NotNil(t, d.Args[7].Key, "the key of a trigger is a projection")
			expect.Equal(t, d.Args[7].Key.Args, []int{0, 1}, "the key of Set lists its key before its scope")
			assert.NotNil(t, d.Args[8].Written, "the update is an argument of a call of Set")
			expect.Equal(t, d.Args[8].Written.Index, 2, "Set takes the update third")
			expect.True(t, d.Bubble, "the probe reads the updates inside a bubble")
			defect := defectOf(t, watches, d)
			expect.Equal(t, defect.Method.Name, "Set", "the planted defect changes the trigger")
			expect.Equal(t, defect.Body, "zero-after", "the planted defect drops each trigger from the fourth on")
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

		gated := []struct {
			name      string
			give      string
			giveFrom  string
			giveClaim string
			wantGate  spec.GateID
		}{
			{
				name: "refuses a total that is no number", give: "Accounts", giveFrom: "Rename",
				giveClaim: "conserves-the-quantity", wantGate: spec.NumericTotal,
			},
			{
				name: "refuses a total without a value", give: "Accounts", giveFrom: "Fold",
				giveClaim: "conserves-the-quantity", wantGate: spec.NumericTotal,
			},
			{
				name: "refuses an iteration whose elements the mutation does not take", give: "Stock",
				giveFrom: "Sizes", giveClaim: "iteration-sees-mutations", wantGate: spec.StreamMutations,
			},
			{
				name: "refuses an iteration of a callable without a stream", give: "Stock", giveFrom: "Snapshot",
				giveClaim: "iteration-sees-mutations", wantGate: spec.StreamMutations,
			},
			{
				name: "refuses an append of no value of the channel", give: "Outbox", giveFrom: "Note",
				giveClaim: "delivers-each-append", wantGate: spec.DeliversOnAChannel,
			},
			{
				name: "refuses a subscription that returns no channel", give: "Outbox", giveFrom: "Log",
				giveClaim: "delivers-each-append", wantGate: spec.DeliversOnAChannel,
			},
			{
				name: "refuses a subscription to a channel that a caller sends to", give: "Outbox", giveFrom: "Post",
				giveClaim: "delivers-each-append", wantGate: spec.DeliversOnAChannel,
			},
			{
				name: "refuses a subscription without a value", give: "Outbox", giveFrom: "Push",
				giveClaim: "delivers-each-append", wantGate: spec.DeliversOnAChannel,
			},
			{
				name: "refuses a subscription without results", give: "Outbox", giveFrom: "Pour",
				giveClaim: "delivers-each-append", wantGate: spec.DeliversOnAChannel,
			},
			{
				name: "refuses a redelivery of a publisher at most once", give: "Broker", giveFrom: "Refire",
				giveClaim: "redelivers-on-demand", wantGate: spec.AtLeastOnceMode,
			},
			{
				name: "refuses a subscription whose next takes no wait", give: "Watches", giveFrom: "Peeks",
				giveClaim: "watch-sees-each-trigger", wantGate: spec.SubscriptionMethods,
			},
			{
				name: "refuses a subscription whose next takes a key", give: "Watches", giveFrom: "Finds",
				giveClaim: "watch-sees-each-trigger", wantGate: spec.SubscriptionMethods,
			},
			{
				name: "refuses a subscription whose next returns no flag", give: "Watches", giveFrom: "Reads",
				giveClaim: "watch-sees-each-trigger", wantGate: spec.SubscriptionMethods,
			},
			{
				name: "refuses a subscription whose stop takes a parameter", give: "Watches", giveFrom: "Cancels",
				giveClaim: "watch-sees-each-trigger", wantGate: spec.SubscriptionMethods,
			},
			{
				name: "refuses a subscription whose stop returns an error", give: "Watches", giveFrom: "Closes",
				giveClaim: "watch-sees-each-trigger", wantGate: spec.SubscriptionMethods,
			},
			{
				name: "refuses a subscription whose next is a field", give: "Watches", giveFrom: "Fields",
				giveClaim: "watch-sees-each-trigger", wantGate: spec.SubscriptionMethods,
			},
			{
				name: "refuses a watch whose trigger passes no update", give: "Watches", giveFrom: "Follows",
				giveClaim: "watch-sees-each-trigger", wantGate: spec.SubscriptionMethods,
			},
		}
		for _, tt := range gated {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				s := suiteOf(tt.give)
				refused := slices.ContainsFunc(s.Derivation.Refusals, func(r derive.Refusal) bool {
					return r.Callable == tt.giveFrom && string(r.Claim) == tt.giveClaim && r.Gate == tt.wantGate
				})
				expect.True(t, refused, "the suite refuses the claim "+tt.giveClaim+" of "+tt.giveFrom+" by the gate "+
					string(tt.wantGate))
				derived := slices.ContainsFunc(s.Derived, func(d suite.Derived) bool {
					return d.Field == tt.giveFrom && strings.HasPrefix(d.Segment, tt.giveClaim)
				})
				expect.False(t, derived, "the suite derives no check of the refused claim")
			})
		}
	})
}
