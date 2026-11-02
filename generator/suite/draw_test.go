// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"maps"
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/suite"
)

// drawnOf returns the draw that the expression e refers to. It stops the test
// where e refers to no draw.
func drawnOf(t *testing.T, e *suite.Expr) *suite.Drawn {
	t.Helper()

	assert.NotNil(t, e.Drawn, "the expression refers to a draw")
	return e.Drawn
}

// drawNames returns the name of each draw of draws.
func drawNames(draws []*suite.Draw) []string {
	out := make([]string, 0, len(draws))
	for _, d := range draws {
		out = append(out, d.Name)
	}
	return out
}

func TestDraw(t *testing.T) {
	t.Parallel()

	families := generate(t, goTree, familiesPackage).suites
	maps.Copy(families, generate(t, goTree, drawnPackage).suites)
	stores := generate(t, goTree, derivedPackage).suites
	derived := func(t *testing.T, decl, field, segment string) *suite.Derived {
		t.Helper()
		s := families[decl]
		return derivedNamed(t, &s, field, segment)
	}

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("draws the input that a write and its read share once", func(t *testing.T) {
			t.Parallel()

			d := derived(t, "Keyed", "Get", "read-after-write#generated")
			expect.Equal(t, drawNames(d.Draws), []string{"drawnKey", "drawnV"}, "the check draws the key and the value")
			read, write := d.Args[0].Closure, d.Args[1].Closure
			expect.Equal(t, drawnOf(t, &read.Args[1]).Draw.Name, "drawnKey", "the read takes the drawn key")
			expect.Equal(t, drawnOf(t, &write.Args[1]).Draw.Name, "drawnKey", "the write takes the same key")
			expect.Equal(t, drawnOf(t, &write.Lets[0].Value).Draw.Name, "drawnV", "the write writes the drawn value")
		})

		t.Run("selects the field of the drawn value that the read takes", func(t *testing.T) {
			t.Parallel()

			s := stores["Store"]
			d := derivedNamed(t, &s, "Get", "read-after-write#generated")
			read := d.Args[0].Closure
			assert.NotNil(t, read.Args[1].Of, "the read selects a field of a value")
			expect.Equal(t, drawnOf(t, read.Args[1].Of).Draw.Name, "drawnS", "the read selects a field of the session")
			expect.Equal(t, read.Args[1].Select, "ID", "the read takes the ID of the drawn session")
		})

		t.Run("keeps the sample that a parameter states", func(t *testing.T) {
			t.Parallel()

			c := derived(t, "Tallies", "Limit", "smoke#generated").Args[0].Closure
			expect.Nil(t, c.Args[1].Drawn, "the count takes its sample")
			expect.NotNil(t, c.Args[1].Value, "the count takes the value of its sample")
			expect.Equal(t, drawnOf(t, &c.Args[2]).Draw.Name, "drawnSku", "the SKU takes a drawn value")
		})

		t.Run("draws a generated function for a parameter of a function type", func(t *testing.T) {
			t.Parallel()

			d := derived(t, "Tallies", "Each", "smoke#generated")
			assert.Length(t, d.Draws, 1, "the check draws the function")
			expect.Equal(t, d.Draws[0].Name, "drawnFn", "the draw has the name of the parameter")
			expect.Length(t, d.Draws[0].Funcs, 1, "the draw is a generated function")
		})

		t.Run("takes the address of a drawn value for a pointer", func(t *testing.T) {
			t.Parallel()

			d := derived(t, "Tallies", "Save", "smoke#generated")
			drawn := drawnOf(t, &d.Args[0].Closure.Args[1])
			expect.True(t, drawn.Address, "the call takes the address of the drawn item")
			expect.Equal(t, spelled(drawn.Draw.Type), "Item", "the check draws a value of the element type")
		})

		t.Run("draws two distinct values for the second write to another partition", func(t *testing.T) {
			t.Parallel()

			d := derived(t, "Scoped", "Set", "partition#generated")
			expect.Equal(t, drawNames(d.Draws), []string{"drawnValue", "drawnScope", "drawnKey"},
				"the check draws the value, the scope and the key of the write")
			for _, draw := range d.Draws {
				expect.True(t, draw.Pair, "the check draws two values of "+draw.Name)
			}
			write := d.Args[0].Closure
			for _, l := range write.Later {
				expect.True(t, drawnOf(t, &l.Value).Second, "the second write takes the second value of "+l.Name)
			}
		})

		t.Run("draws two distinct values of the field that a second write changes", func(t *testing.T) {
			t.Parallel()

			d := derived(t, "Updater", "Put", "updater#generated")
			write := d.Args[0].Closure
			assert.Length(t, write.Later, 1, "the second write changes one field")
			changed := drawnOf(t, &write.Later[0].Value)
			expect.Equal(t, changed.Draw.Name, "drawnVBody", "the check draws the field Body of the value")
			expect.Equal(t, changed.Draw.Label, "v.body", "the label joins the input and the field")
			expect.True(t, changed.Second, "the second write takes the second value of the field")
		})

		t.Run("keeps the samples of the field of a written value with a sample of its own", func(t *testing.T) {
			t.Parallel()

			write := derived(t, "Notes", "Put", "updater#generated").Args[0].Closure
			expect.Equal(t, drawnOf(t, &write.Args[1]).Draw.Name, "drawnShelf", "the write takes the drawn shelf")
			assert.Length(t, write.Later, 1, "the second write changes one field")
			for _, l := range append(write.Lets, write.Later...) {
				expect.Nil(t, l.Value.Drawn, "the note keeps the sample of "+l.Name)
				expect.NotNil(t, l.Value.Value, "the note takes the value of the sample of "+l.Name)
			}
		})

		t.Run("keeps the sample of a variadic parameter of functions", func(t *testing.T) {
			t.Parallel()

			c := derived(t, "Tallies", "Visit", "smoke#generated").Args[0].Closure
			expect.Equal(t, drawnOf(t, &c.Args[1]).Draw.Name, "drawnSku", "the SKU takes a drawn value")
			expect.Nil(t, c.Args[2].Drawn, "the hooks take their sample")
			expect.NotNil(t, c.Args[2].Value, "the hooks take the value of their sample")
		})

		t.Run("keeps the sample of a generated function for its second value", func(t *testing.T) {
			t.Parallel()

			write := derived(t, "Hooked", "Set", "partition#generated").Args[0].Closure
			fn := func(lets []suite.Let) *suite.Expr {
				at := slices.IndexFunc(lets, func(l suite.Let) bool { return l.Name == "s2" })
				assert.True(t, at >= 0, "the write declares the hook of its partition")
				return &lets[at].Value
			}
			expect.Length(t, drawnOf(t, fn(write.Lets)).Draw.Funcs, 1, "the first write takes a generated hook")
			expect.Nil(t, fn(write.Later).Drawn, "the second write takes the sample of the hook")
		})

		t.Run("keeps the sample of an input with a part that a call consumes", func(t *testing.T) {
			t.Parallel()

			c := derived(t, "Tallies", "Schedule", "smoke#generated").Args[0].Closure
			expect.Equal(t, drawnOf(t, &c.Args[1]).Draw.Name, "drawnSku", "the SKU takes a drawn value")
			expect.Nil(t, c.Args[2].Drawn, "the plan with a hook takes its sample")
			expect.NotNil(t, c.Args[2].Value, "the plan takes the value of its sample")
		})

		tracker := generate(t, goTree, nestedPackage).suites["Tracker"]
		samples := []struct {
			name string
			give string
		}{
			{name: "keeps the sample of a generated function whose results are generated functions", give: "Walk"},
			{
				name: "keeps the sample of a generated implementation whose methods return generated functions",
				give: "Tick",
			},
		}
		for _, tt := range samples {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				derivedNamed(t, &tracker, tt.give, "deterministic")
				drawn := slices.ContainsFunc(tracker.Derived, func(d suite.Derived) bool {
					return d.Field == tt.give && d.Segment == "deterministic#generated"
				})
				expect.False(t, drawn, "the input of "+tt.give+" keeps its sample, so the check draws no input")
			})
		}

		t.Run("draws an input of a type that contains itself", func(t *testing.T) {
			t.Parallel()

			c := derived(t, "Tallies", "Link", "smoke#generated").Args[0].Closure
			expect.Equal(t, drawnOf(t, &c.Args[1]).Draw.Name, "drawnN", "the node takes a drawn value")
		})

		t.Run("records no generated binding of a check whose draws take no choice of the case", func(t *testing.T) {
			t.Parallel()

			tallies := families["Tallies"]
			derivedNamed(t, &tallies, "Subscribe", "smoke")
			drawn := slices.ContainsFunc(tallies.Derived, func(d suite.Derived) bool {
				return d.Field == "Subscribe" && d.Segment == "smoke#generated"
			})
			expect.False(t, drawn, "a generated function without results takes no choice of the case")
		})

		t.Run("draws two values for two parameters of one name and two types", func(t *testing.T) {
			t.Parallel()

			d := derived(t, "Labels", "Get", "read-after-write#generated")
			expect.Equal(t, drawNames(d.Draws), []string{"drawnName", "drawnName0", "drawnV"},
				"the check numbers the name of the second draw of name")
			expect.NotEqual(t, spelled(d.Draws[0].Type), spelled(d.Draws[1].Type),
				"the two draws of name have two types")
		})
	})
}
