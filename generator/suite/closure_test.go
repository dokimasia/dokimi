// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"

	"go.dokimi.dev/dokimi/generator/suite"
)

// names returns the names of the variables lets.
func names(lets []suite.Let) []string {
	out := make([]string, 0, len(lets))
	for _, l := range lets {
		out = append(out, l.Name)
	}
	return out
}

func TestClosure(t *testing.T) {
	t.Parallel()

	families := generate(t, goTree, familiesPackage).suites
	closure := func(t *testing.T, decl, field, segment string, arg int) *suite.Closure {
		t.Helper()
		s := families[decl]
		d := derivedNamed(t, &s, field, segment)
		assert.True(t, arg < len(d.Args), "the probe of the check takes the argument")
		c := d.Args[arg].Closure
		assert.NotNil(t, c, "the argument is a closure")
		return c
	}

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("takes the two inputs of a call that combines two values", func(t *testing.T) {
			t.Parallel()

			tallies := generate(t, goTree, drawnPackage).suites["Tallies"]
			c := derivedNamed(t, &tallies, "Merge", "commutative").Args[0].Closure
			assert.NotNil(t, c, "the call is a closure")
			assert.NotNil(t, c.Param, "the closure takes the two values")
			expect.Equal(t, c.Param.Name, "x", "the first value is x")
			expect.Equal(t, c.Param.Other, "y", "the second value is y")
			expect.Equal(t, spelled(c.Param.Type), "Tally", "the values have the type of the inputs")
			expect.Equal(t, c.Args, []suite.Expr{{Name: "x"}, {Name: "y"}}, "the call passes the two values")
			expect.Equal(t, c.Results, []suite.Expr{{Name: "r0"}, {Word: "nil"}},
				"the closure returns the value beside a nil error")
		})

		t.Run("passes the key of each of the two writes to a batch read", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			d := derivedNamed(t, &shelf, "GetMany", "reads-each-key")
			read, write := d.Args[0].Closure, d.Args[1].Closure
			assert.NotNil(t, read, "the batch read is a closure")
			assert.NotNil(t, write, "the write is a closure")
			assert.Length(t, read.Args, 3, "the batch read passes the context and two keys")
			assert.Length(t, write.Lets, 2, "the write declares its value and sets its key")
			assert.Length(t, write.Later, 1, "the second write changes its key")
			expect.Equal(t, write.Lets[1].Select, "Key", "the first write sets the key")
			expect.Equal(t, write.Later[0].Select, "Key", "the second write changes the key")
			expect.Equal(t, read.Args[1:], []suite.Expr{write.Lets[1].Value, write.Later[0].Value},
				"the batch read passes the key of each write")
		})

		t.Run("passes the value of each of the two writes to a batch read of the values", func(t *testing.T) {
			t.Parallel()

			tags := generate(t, goTree, classifiedPackage).suites["Tags"]
			d := derivedNamed(t, &tags, "Has", "reads-each-key")
			read, write := d.Args[0].Closure, d.Args[1].Closure
			assert.NotNil(t, read, "the batch read is a closure")
			assert.NotNil(t, write, "the write is a closure")
			assert.Equal(t, names(write.Lets), []string{"w"}, "the write declares the value alone")
			assert.Equal(t, names(write.Later), []string{"w"}, "the second write changes the value")
			expect.Equal(t, read.Args[1:], []suite.Expr{write.Lets[0].Value, write.Later[0].Value},
				"the batch read passes the value of each write")
		})

		t.Run("returns the error of a call that returns a stream beside it", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "Scan", "over-match").Args[0].Closure
			assert.NotNil(t, c, "the drain is a closure")
			expect.Equal(t, c.Assign, []string{"r0", "err"}, "the closure assigns the stream and the error")
			expect.Equal(t, c.Results, []suite.Expr{{Name: "r0"}, {Name: "err"}},
				"the closure drains the stream after the error of the call")
		})

		t.Run("drains a stream with an error into the list of its elements", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "All", "streams-what-was-stored").Args[0].Closure
			assert.NotNil(t, c, "the drain is a closure")
			expect.Equal(t, c.Drain, "until-error", "the closure drains the elements until the first error")
			assert.Length(t, c.Returns, 1, "the closure returns a list of one type")
			expect.Equal(t, spelled(c.Returns[0]), "Item", "the list contains the items")
			expect.False(t, c.Pair, "an element is no pair")
		})

		t.Run("drains a stream of elements into a list", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "Items", "streams-what-was-stored").Args[0].Closure
			assert.NotNil(t, c, "the drain is a closure")
			expect.Equal(t, c.Drain, "elements", "the closure drains each element")
			assert.Length(t, c.Returns, 1, "the closure returns a list of one type")
			expect.Equal(t, spelled(c.Returns[0]), "Item", "the list contains the items")
		})

		t.Run("drains a stream of pairs into a list of pairs", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "Keyed", "streams-what-was-stored").Args[0].Closure
			assert.NotNil(t, c, "the drain is a closure")
			expect.Equal(t, c.Drain, "pairs", "the closure drains each pair")
			assert.Length(t, c.Returns, 2, "the closure returns a list of pairs of two types")
			expect.Equal(t, spelled(c.Returns[0]), "string", "the first of a pair is the key")
			expect.Equal(t, spelled(c.Returns[1]), "Item", "the second of a pair is the item")
			expect.True(t, c.Pair, "an element is a pair")
		})

		t.Run("takes the position of an indexed call", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "At", "respects-its-bounds").Args[0].Closure
			assert.NotNil(t, c, "the call is a closure")
			assert.NotNil(t, c.Param, "the closure takes the position")
			expect.Equal(t, c.Param.Name, "position", "the parameter is the position")
			expect.Equal(t, spelled(c.Param.Type), "int", "the position is an int")
			expect.Equal(t, c.Args, []suite.Expr{{Word: "ctx"}, {Name: "position"}},
				"the call passes the position to the integer input")
		})

		t.Run("passes the first sample to an input of an indexed call that is no integer", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "Slot", "respects-its-bounds").Args[0].Closure
			assert.NotNil(t, c, "the call is a closure")
			assert.Length(t, c.Args, 3, "the call passes the context, the name and the position")
			expect.NotNil(t, c.Args[1].Value, "the name takes its sample")
			expect.Equal(t, c.Args[2].Name, "position", "the call passes the position to the integer input")
		})

		t.Run("converts the position to an unsigned input", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "Slot", "respects-its-bounds").Args[0].Closure
			assert.NotNil(t, c, "the call is a closure")
			assert.Length(t, c.Args, 3, "the call passes the context, the name and the position")
			assert.NotNil(t, c.Args[2].Convert, "the call converts the position")
			expect.Equal(t, spelled(c.Args[2].Convert), "uint", "the position of a slot is a uint")
		})

		t.Run("converts the position to an integer input of another type", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "AtOffset", "respects-its-bounds").Args[0].Closure
			assert.NotNil(t, c, "the call is a closure")
			assert.Length(t, c.Args, 2, "the call passes the context and the offset")
			assert.NotNil(t, c.Args[1].Convert, "the call converts the position")
			expect.Equal(t, spelled(c.Args[1].Convert), "int64", "the offset is an int64")
			expect.Equal(t, c.Args[1].Name, "position", "the call converts the position")
		})

		t.Run("calls the sibling of a wrapper on a new subject", func(t *testing.T) {
			t.Parallel()

			shelf := generate(t, goTree, classifiedPackage).suites["Shelf"]
			c := derivedNamed(t, &shelf, "TracedGet", "agrees-with-its-sibling").Args[1].Closure
			assert.NotNil(t, c, "the sibling is a closure")
			expect.Equal(t, c.Receiver, &suite.Expr{Word: "h.New()"}, "the sibling runs on a new subject")
			expect.Equal(t, c.Method, "Get", "the sibling is Get")
		})

		t.Run("returns the call of a callable that returns one value beside an error", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Store", "Get", "zero-on-error", 0)
			expect.Empty(t, c.Assign, "the closure assigns nothing")
			expect.Empty(t, c.Results, "the closure returns the results of the call")
			expect.Length(t, c.Returns, 1, "the closure returns the value")
			expect.True(t, c.Error, "the closure returns the error")
		})

		t.Run("returns the two values of a callable as a pair", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Store", "GetWithMeta", "zeroes-on-error", 0)
			expect.True(t, c.Pair, "the closure returns a pair")
			expect.Equal(t, c.Assign, []string{"r0", "r1", "err"}, "the closure assigns the values and the error")
			assert.Length(t, c.Results, 2, "the closure returns the pair and the error")
			expect.NotNil(t, c.Results[0].Pair, "the first result is the pair of the values")
		})

		t.Run("returns the call of a lookup as it is", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Store", "Load", "zero-on-miss", 0)
			expect.True(t, c.Flag, "the closure returns the flag")
			expect.Empty(t, c.Assign, "the closure returns the results of the call")
		})

		t.Run("returns the first value of a callable that returns a flag beside it", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Store", "Load", "cacheable", 0)
			expect.Equal(t, c.Assign, []string{"r0", "_"}, "the closure leaves the flag out")
			expect.Equal(t, c.Results, []suite.Expr{{Name: "r0"}, {Word: "nil"}},
				"the closure returns the value beside a nil error")
		})

		t.Run("returns the error of a callable that returns values beside it", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Store", "Get", "cancel", 0)
			expect.Equal(t, c.Assign, []string{"_", "err"}, "the closure leaves the value out")
			expect.Equal(t, c.Results, []suite.Expr{{Name: "err"}}, "the closure returns the error")
		})

		t.Run("returns the call of a callable that returns an error alone", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Store", "Put", "smoke", 0)
			expect.Empty(t, c.Assign, "the closure assigns nothing")
			expect.False(t, c.Discard, "the closure returns the call")
		})

		t.Run("returns nil after a callable without an error", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Plain", "Set", "smoke", 0)
			expect.True(t, c.Discard, "the closure makes the call a statement")
			expect.Equal(t, c.Results, []suite.Expr{{Word: "nil"}}, "the closure returns nil")
		})

		t.Run("passes the written value to the input after the key that the read shares", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Keyed", "Get", "read-after-write", 1)
			expect.Equal(t, names(c.Lets), []string{"w"}, "the write declares the written value")
			assert.Length(t, c.Args, 3, "the write takes the context, the key and the value")
			expect.NotNil(t, c.Args[1].Value, "the key takes its sample")
			expect.Equal(t, c.Args[2], suite.Expr{Name: "w"}, "the value takes the written value")
			expect.Equal(t, c.Results, []suite.Expr{{Name: "w"}, {Name: "err"}},
				"the write returns the written value beside its error")
		})

		t.Run("returns the written value beside nil after a write without an error", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Plain", "Get", "read-after-write", 1)
			expect.True(t, c.Discard, "the write is a statement")
			expect.Equal(t, c.Results, []suite.Expr{{Name: "w"}, {Word: "nil"}},
				"the write returns the written value beside nil")
		})

		t.Run("shares the sample of the parameter of the write that the read takes by its name", func(t *testing.T) {
			t.Parallel()

			read := closure(t, "Tagged", "Get", "read-after-write", 0)
			write := closure(t, "Tagged", "Get", "read-after-write", 1)
			assert.Length(t, read.Args, 2, "the read takes the context and the tag")
			assert.Length(t, write.Args, 3, "the write takes the context, the value and the tag")
			expect.Equal(t, read.Args[1], write.Args[2], "the read takes the tag that the write passes")
		})

		t.Run("changes the field of the second write that the read does not select", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Updater", "Put", "updater", 0)
			assert.Length(t, c.Later, 1, "the closure changes the written value on its later calls")
			expect.Equal(t, c.Later[0].Select, "Body", "the second write changes the field Body")
		})

		t.Run("keeps the key of a second write that the read shares", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "KeyedUpdater", "Put", "updater", 0)
			expect.Equal(t, names(c.Later), []string{"w"}, "the second write changes the value alone")
			expect.NotNil(t, c.Args[1].Value, "the key keeps its sample")
		})

		t.Run("writes the second value into another partition", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Scoped", "Set", "partition", 0)
			expect.Equal(
				t,
				names(c.Lets),
				[]string{"w", "s0", "s1"},
				"the closure declares the value and the selectors",
			)
			expect.Equal(t, names(c.Later), []string{"w", "s0", "s1"},
				"the later calls change the value, the scope and the key")
			expect.Equal(t, c.Args[1:], []suite.Expr{{Name: "s0"}, {Name: "s1"}, {Name: "w"}},
				"the call takes the selectors and the value")
		})

		t.Run("writes the second sample on the later calls of a write of a drain", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Bag", "Items", "no-duplicates", 1)
			assert.Length(t, c.Later, 1, "the closure changes the written value on its later calls")
			expect.Equal(t, c.Later[0].Select, "", "the later calls write the whole second sample")
		})

		t.Run("changes the field of the write between two reads of a snapshot", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Snapshots", "Get", "point-in-time", 1)
			assert.Length(t, c.Lets, 2, "the write declares its value and changes one field of it")
			expect.NotEqual(t, c.Lets[1].Select, "", "the second variable changes a field")
			assert.Length(t, c.Later, 1, "the later writes change the field again")
			expect.Equal(t, c.Later[0].Select, c.Lets[1].Select, "the later writes change the same field")
			expect.NotEqual(t, c.Later[0].Value, c.Lets[1].Value, "the later writes take another value of the field")
		})

		t.Run("changes the judged field on the later calls of a guarded write", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Guarded", "Put", "if-match", 0)
			assert.Length(t, c.Lets, 2, "the first write declares the value and sets the judged field")
			expect.Equal(t, c.Lets[1].Select, "Body", "the first write sets the field that the predicate judges")
			assert.Length(t, c.Later, 1, "the later writes change the value")
			expect.NotEqual(t, c.Later[0].Value, c.Lets[1].Value, "the later writes take another value of the field")
			expect.Equal(t, c.Later[0].Select, "Body", "the later writes change the field that the predicate judges")
			expect.NotNil(t, c.Args[1].Value, "the metadata takes its sample")
			expect.Equal(t, c.Args[2], suite.Expr{Name: "w"}, "the write takes the value that has the judged field")
			expect.Empty(t, c.Results, "the closure returns the call of the write")
		})

		t.Run("passes the value of the second guarded write to the predicate", func(t *testing.T) {
			t.Parallel()

			write := closure(t, "Guarded", "Put", "if-match", 0)
			match := closure(t, "Guarded", "Put", "if-match", 1)
			expect.Equal(t, match.Lets, []suite.Let{write.Lets[0], write.Later[0]},
				"the predicate judges the value of the second write")
			expect.Equal(t, match.Args[1], write.Args[1], "the predicate takes the metadata of the write")
			expect.Equal(t, match.Args[2], suite.Expr{Name: "w"}, "the predicate takes the value")
			expect.True(t, match.Error, "the predicate returns its verdict beside its error")
		})

		t.Run("returns the value that a write returns of the type that it writes", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Stamped", "Get", "read-after-write", 1)
			expect.Equal(t, names(c.Lets), []string{"w"}, "the write declares the written value")
			expect.Length(t, c.Returns, 1, "the write returns one value beside its error")
			expect.Empty(t, c.Results, "the write returns the value that the call returns")
		})

		t.Run("passes its own sample to a function parameter of a read", func(t *testing.T) {
			t.Parallel()

			read := closure(t, "Decoded", "Get", "read-after-write", 0)
			write := closure(t, "Decoded", "Get", "read-after-write", 1)
			assert.Length(t, read.Args, 3, "the read takes the context, the key and the function")
			expect.Equal(t, read.Args[1], write.Args[1], "the read takes the key that the write passes")
			expect.NotNil(t, read.Args[2].Value, "the function takes a sample")
			expect.NotEqual(t, read.Args[2], read.Args[1], "the function takes the sample of its own type")
		})

		t.Run("takes the body of a scope as its parameter", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Scope", "Run", "transaction", 0)
			assert.NotNil(t, c.Param, "the closure takes a parameter")
			expect.True(t, c.Param.Body, "the parameter is the body of the scope")
			expect.Equal(t, c.Args[1], suite.Expr{Name: "body"}, "the scope runs the body")
		})

		t.Run("takes the value of a validator as its parameter", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Screened", "Store", "validates", 0)
			assert.NotNil(t, c.Param, "the closure takes a parameter")
			expect.Equal(t, c.Param.Name, "v", "the parameter is the value")
			expect.Equal(t, c.Args[1], suite.Expr{Name: "v"}, "the call takes the parameter")
		})

		t.Run("returns the value of the inverse of a codec", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Codec", "Encode", "codec", 1)
			expect.Equal(t, c.Method, "Decode", "the inverse calls Decode")
			expect.Equal(t, c.Param.Name, "v", "the inverse takes the encoding")
			expect.Empty(t, c.Assign, "the inverse returns the results of the call")
		})

		t.Run("returns the input of a round trip beside its result", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Codec", "Encode", "codec", 0)
			expect.Equal(t, names(c.Lets), []string{"in"}, "the round trip declares its input")
			expect.True(t, c.Pair, "the round trip returns a pair")
			assert.NotNil(t, c.Results[0].Pair, "the first result is the pair")
			expect.Equal(t, c.Results[0].Pair.Values[0], suite.Expr{Name: "in"}, "the pair opens with the input")
		})

		t.Run("adapts the hook of the probe to the function type of the registrar", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Events", "Fire", "hooks", 1)
			assert.NotNil(t, c.Param, "the registrar takes the hook")
			expect.True(t, c.Param.Hook, "the parameter is the hook")
			expect.NotNil(t, c.Args[0].Hook, "the registrar takes a function of its own type")
		})

		t.Run("calls the method of the handle that a param names", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Opener", "Open", "cursor-open", 1)
			expect.Equal(t, c.Receiver, &suite.Expr{Name: "handle"}, "the closure calls the handle")
			expect.Equal(t, c.Method, "Close", "the closure calls the method that close-method names")
		})

		t.Run("passes the handle to the partner that closes it", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Ledger", "Begin", "tx", 1)
			expect.Equal(t, c.Method, "Commit", "the closure calls the partner commit")
			expect.Equal(t, c.Args[1], suite.Expr{Name: "handle"}, "the partner takes the handle")
		})

		t.Run("returns the delay that a registration passes", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Timers", "At", "scheduled", 0)
			expect.Equal(t, names(c.Lets), []string{"d"}, "the registration declares its delay")
			expect.Equal(t, c.Results[0], suite.Expr{Name: "d"}, "the registration returns its delay")
		})

		t.Run("returns the stamp of each entry of a replay", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Log", "Append", "chain-causal", 1)
			expect.Equal(t, c.Stamps, "Seq", "the replay returns the field that the param version names")
		})

		t.Run("passes the value of a param to the first input of the call", func(t *testing.T) {
			t.Parallel()

			c := closure(t, "Renderer", "Render", "xss-safe", 0)
			assert.NotNil(t, c.Args[1].Value, "the first input takes a value")
			expect.Equal(t, c.Args[1].Value.Kind, emit.ValueName, "the first input takes the declaration Unsafe")
		})
	})
}
