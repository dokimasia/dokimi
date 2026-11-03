// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/diag"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/symbol"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/gentest"
)

// messagesOf returns the message of each finding of got with the code, in
// the order of the messages.
func messagesOf(got generated, code diag.Code) []string {
	var out []string
	for _, d := range got.findings {
		if d.Code == code {
			out = append(out, d.Msg)
		}
	}
	slices.Sort(out)
	return out
}

func TestSample(t *testing.T) {
	t.Parallel()

	p := suite.New()
	f := fixture(t, goTree, p, sealedPackage)
	got := recorded(t, f, p)
	aead := got.suites["AEAD"]

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("takes each input of a derived call from the builder of the subject", func(t *testing.T) {
			t.Parallel()

			call := derivedNamed(t, &aead, "Seal", "smoke").Args[0].Closure
			assert.NotNil(t, call, "the call is a closure")
			assert.Length(t, call.Builds, 1, "the call declares the values of one builder")
			build := call.Builds[0]
			expect.Equal(t, []any{build.Builder.Func.Spelling, build.Builder.Subject}, []any{"SealInputs", true},
				"the builder of Seal takes the subject")
			expect.Equal(t, build.Names, []string{"builtSeal0", "builtSeal1", "builtSeal2", "builtSeal3"},
				"the call declares a value of each input")
			names := make([]string, 0, len(call.Args))
			for _, a := range call.Args {
				assert.NotNil(t, a.Built, "each argument is a value of the builder")
				names = append(names, a.Built.Name)
			}
			expect.Equal(t, names, build.Names, "the call passes the values of the builder in order")
		})

		t.Run("lists the machinery of the inputs of a builder that takes the subject", func(t *testing.T) {
			t.Parallel()

			metronome := got.suites["Metronome"]
			expect.Equal(t, derivedNamed(t, &aead, "Seal", "smoke").Machineries,
				[]string{"InputsBuiltFromTheSubject", "NothingBeyondACaller"},
				"the builder of Seal takes the subject, which it can call")
			expect.Equal(t, derivedNamed(t, &metronome, "Tick", "smoke").Machineries, []string{"NothingBeyondACaller"},
				"the builder of Tick takes no subject")
		})

		t.Run("records no generated binding of a callable with a builder", func(t *testing.T) {
			t.Parallel()

			expect.False(t, slices.ContainsFunc(aead.Derived, func(d suite.Derived) bool {
				return d.Field == "Seal" && d.Generated
			}), "a check of Seal draws no input")
		})

		t.Run("calls the builder in the call generator of a callable", func(t *testing.T) {
			t.Parallel()

			c := called(t, aead, "Open")
			assert.NotNil(t, c.Builder, "the call generator calls the builder")
			expect.Equal(t, []any{c.Builder.Func.Spelling, c.Builder.Subject}, []any{"OpenInputs", true},
				"the builder of Open takes the subject")
			names := make([]string, 0, len(c.Params))
			for _, p := range c.Params {
				names = append(names, p.Name)
			}
			expect.Equal(t, names, []string{"dst", "nonce", "ciphertext", "additionalData"},
				"the builder returns a value of each parameter")
		})

		t.Run("calls a builder without the subject", func(t *testing.T) {
			t.Parallel()

			metronome := got.suites["Metronome"]
			c := called(t, metronome, "Tick")
			assert.NotNil(t, c.Builder, "the call generator calls the builder")
			expect.False(t, c.Builder.Subject, "the builder of Tick takes no subject")
			call := derivedNamed(t, &metronome, "Tick", "smoke").Args[0].Closure
			assert.NotNil(t, call, "the call is a closure")
			expect.Equal(t, call.Builds[0].Names, []string{"builtTick0"}, "the call declares the value of n")
		})

		t.Run("names an unnamed parameter of a call generator with a builder by its position", func(t *testing.T) {
			t.Parallel()

			c := called(t, got.suites["Metronome"], "Pace")
			assert.Length(t, c.Params, 2, "the call generator takes the context and the int")
			expect.Equal(t, c.Params[1].Name, "a1", "the int takes the name of its position")
		})

		t.Run("names each value of a builder with a name that no other value or package takes", func(t *testing.T) {
			t.Parallel()

			c := called(t, got.suites["Mixer"], "Mix")
			names := make([]string, 0, len(c.Params))
			for _, p := range c.Params {
				names = append(names, p.Name)
			}
			expect.Equal(t, names, []string{"", "a1", "time0", "subject0", "subject00", "track"},
				"each value takes its own name, or the first name that the call, a package or a value before it "+
					"does not take")
			expect.False(t, c.Params[5].Address, "the call passes the pointer that the builder returns")
		})

		t.Run("declares the values of two builders and of a written field in one closure", func(t *testing.T) {
			t.Parallel()

			depot := got.suites["Depot"]
			d := derivedNamed(t, &depot, "Find", "read-after-write")
			reads := func(a suite.Arg) bool { return a.Closure != nil && a.Closure.Method == "Find" }
			at := slices.IndexFunc(d.Args, reads)
			assert.True(t, at >= 0, "the check reads through a closure")
			read := d.Args[at].Closure
			builds := make([]string, 0, len(read.Builds))
			for _, b := range read.Builds {
				builds = append(builds, b.Builder.Func.Spelling+" "+strings.Join(b.Names, ","))
			}
			expect.Equal(t, builds, []string{"Stocked builtStock0", "Found _,builtFind1"},
				"the read calls the builder of the write and its own builder once each")
			sku := read.Args[1].Of
			assert.NotNil(t, sku, "the SKU is a field of the written item")
			assert.NotNil(t, sku.Built, "the written item is a value of the builder of Stock")
			expect.Equal(t, sku.Built.Name, "builtStock0", "the read takes the SKU of the item that it declares")
		})

		t.Run("calls the builder in place for a value outside a closure", func(t *testing.T) {
			t.Parallel()

			locker := got.suites["Locker"]
			durable := derivedNamed(t, &locker, "Put", "synced-writes-survive").Args[0].Durable
			assert.NotNil(t, durable, "the argument is the calls of a durable writer")
			want := durable.Wants[0].Built
			assert.NotNil(t, want, "the value of a sample is a value of the builder")
			expect.Equal(t, []string{want.Name, want.Values()}, []string{"", "_, v"},
				"the value calls the builder in place and keeps its second value")
		})

		t.Run("reports BadBuilder at a builder of another signature", func(t *testing.T) {
			t.Parallel()

			expect.Equal(t, messagesOf(got, suite.BadBuilder), []string{
				"the builder Key of Short.Put returns (string), and a builder returns one value of the type of each " +
					"input of Put in order, (string, int), so Short has no suite",
				"the builder StallInputs of Stall.Add takes (Stall), and a builder takes the subject of the suite or " +
					"nothing, so Stall has no suite",
				"the builder Strings of Wrong.Put returns (string, string), and a builder returns one value of the " +
					"type of each input of Put in order, (string, int), so Wrong has no suite",
				"the builder Takes of Stamp takes (int), and a builder takes the subject of the suite or nothing, " +
					"so Stamp has no suite",
				"the builder Takes of Taker.Get takes (int), and a builder takes the subject of the suite or nothing, " +
					"so Taker has no suite",
			}, "suite reports each builder of another signature")
			for _, name := range []string{"Short", "Wrong", "Taker", "Stamp", "Stall"} {
				_, found := got.suites[name]
				expect.False(t, found, "suite records no suite of "+name)
			}
			for _, at := range []node.Declaration{
				f.Find(t, symbol.KindMethod, "Short", "Put"), f.Find(t, symbol.KindFunction, "", "Stamp"),
			} {
				expect.True(t, slices.ContainsFunc(got.findings, func(d diag.Diag) bool {
					return d.Code == suite.BadBuilder && d.Pos == at.Position()
				}), "suite reports the builder of "+at.Identity().Name+" at its callable")
			}
		})

		t.Run("reports Unimportable at a builder that the file of the suite cannot import", func(t *testing.T) {
			t.Parallel()

			messages := messagesOf(got, suite.Unimportable)
			assert.Length(t, messages, 1, "suite reports the builder of Concealed")
			expect.True(t, strings.HasPrefix(messages[0], "the builder hiddenInputs of Concealed.Count is a function "+
				"that the file of the suite in the package example.com/shop/sealed/sealedtest cannot import, so "+
				"Concealed has no suite: "), "the message names the builder and the package of the suite")
			_, found := got.suites["Concealed"]
			expect.False(t, found, "suite records no suite of Concealed")
		})

		t.Run("leaves the import of a builder to the layout where the plan previews no package", func(t *testing.T) {
			t.Parallel()

			p := suite.New()
			s := specOf(goTree, p, sealedPackage)
			s.Types = nil
			_, found := recorded(t, gentest.New(t, s), p).suites["Concealed"]
			expect.True(t, found, "suite records the suite of Concealed, and the layout refuses its builder")
		})
	})
}
