package declared_test

import (
	"testing"

	"example.com/shop/declared"
	"go.dokimi.dev/assert"
)

// TestLevel checks the functions and the methods of the enumeration Level.
func TestLevel(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[declared.Level]bool{}
		for _, v := range []declared.Level{
			declared.Debug,
			declared.Info,
			declared.Warn,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("LevelValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []declared.Level{
					declared.Debug,
					declared.Info,
					declared.Warn,
				}
				assert.Equal(
					t,
					declared.LevelValues(),
					want,
					"LevelValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := declared.LevelValues(), declared.LevelValues()
			assert.NotEqual(
				t,
				&first[0],
				&second[0],
				"the two calls return two slices",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []declared.Level{
				declared.Debug,
				declared.Info,
				declared.Warn,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := declared.Level(4)
				for _, v := range []declared.Level{
					declared.Debug,
					declared.Info,
					declared.Warn,
				} {
					assert.NotEqual(
						t,
						out.String(),
						v.String(),
						"an undeclared value has a text of its own",
					)
				}
			},
		)
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Level
		}{
			{name: "accepts Debug", give: declared.Debug},
			{name: "accepts Info", give: declared.Info},
			{name: "accepts Warn", give: declared.Warn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, declared.Level(4).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero declared.Level
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseLevel", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Level
		}{
			{name: "returns Debug for its text", give: declared.Debug},
			{name: "returns Info for its text", give: declared.Info},
			{name: "returns Warn for its text", give: declared.Warn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := declared.ParseLevel(tt.give.String())
				assert.NoError(t, err, "ParseLevel accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseLevel returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := declared.ParseLevel("??????")
			assert.ErrorIs(
				t,
				err,
				declared.ErrUnknownLevel,
				"ParseLevel refuses the text with ErrUnknownLevel",
			)
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Level
		}{
			{name: "encodes Debug as its text", give: declared.Debug},
			{name: "encodes Info as its text", give: declared.Info},
			{name: "encodes Warn as its text", give: declared.Warn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				assert.Equal(
					t,
					string(text),
					tt.give.String(),
					"the encoding of the variant is its text",
				)
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			out := declared.Level(4)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				declared.ErrUnknownLevel,
				"MarshalText refuses the value with ErrUnknownLevel",
			)
			assert.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Level
		}{
			{name: "decodes Debug from its encoding", give: declared.Debug},
			{name: "decodes Info from its encoding", give: declared.Info},
			{name: "decodes Warn from its encoding", give: declared.Warn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got declared.Level
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got declared.Level
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestColor checks the functions and the methods of the enumeration Color.
func TestColor(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[declared.Color]bool{}
		for _, v := range []declared.Color{
			declared.Red,
			declared.Green,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("ColorValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []declared.Color{
					declared.Red,
					declared.Green,
				}
				assert.Equal(
					t,
					declared.ColorValues(),
					want,
					"ColorValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := declared.ColorValues(), declared.ColorValues()
			assert.NotEqual(
				t,
				&first[0],
				&second[0],
				"the two calls return two slices",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []declared.Color{
				declared.Red,
				declared.Green,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := declared.Color(3)
				for _, v := range []declared.Color{
					declared.Red,
					declared.Green,
				} {
					assert.NotEqual(
						t,
						out.String(),
						v.String(),
						"an undeclared value has a text of its own",
					)
				}
			},
		)
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Color
		}{
			{name: "accepts Red", give: declared.Red},
			{name: "accepts Green", give: declared.Green},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, declared.Color(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero declared.Color
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseColor", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Color
		}{
			{name: "returns Red for its text", give: declared.Red},
			{name: "returns Green for its text", give: declared.Green},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := declared.ParseColor(tt.give.String())
				assert.NoError(t, err, "ParseColor accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseColor returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := declared.ParseColor("??????")
			assert.ErrorIs(
				t,
				err,
				declared.ErrUnknownColor,
				"ParseColor refuses the text with ErrUnknownColor",
			)
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Color
		}{
			{name: "encodes Red as its text", give: declared.Red},
			{name: "encodes Green as its text", give: declared.Green},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				assert.Equal(
					t,
					string(text),
					tt.give.String(),
					"the encoding of the variant is its text",
				)
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			out := declared.Color(3)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				declared.ErrUnknownColor,
				"MarshalText refuses the value with ErrUnknownColor",
			)
			assert.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Color
		}{
			{name: "decodes Red from its encoding", give: declared.Red},
			{name: "decodes Green from its encoding", give: declared.Green},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got declared.Color
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got declared.Color
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestLane checks the functions and the methods of the enumeration Lane.
func TestLane(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[declared.Lane]bool{}
		for _, v := range []declared.Lane{
			declared.Left,
			declared.Right,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("LaneValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []declared.Lane{
					declared.Left,
					declared.Right,
				}
				assert.Equal(
					t,
					declared.LaneValues(),
					want,
					"LaneValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := declared.LaneValues(), declared.LaneValues()
			assert.NotEqual(
				t,
				&first[0],
				&second[0],
				"the two calls return two slices",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []declared.Lane{
				declared.Left,
				declared.Right,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := declared.Lane(3)
				for _, v := range []declared.Lane{
					declared.Left,
					declared.Right,
				} {
					assert.NotEqual(
						t,
						out.String(),
						v.String(),
						"an undeclared value has a text of its own",
					)
				}
			},
		)
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Lane
		}{
			{name: "accepts Left", give: declared.Left},
			{name: "accepts Right", give: declared.Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, declared.Lane(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero declared.Lane
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseLane", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Lane
		}{
			{name: "returns Left for its text", give: declared.Left},
			{name: "returns Right for its text", give: declared.Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := declared.ParseLane(tt.give.String())
				assert.NoError(t, err, "ParseLane accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseLane returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := declared.ParseLane("??????")
			assert.HasError(t, err, "ParseLane refuses the text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Lane
		}{
			{name: "encodes Left as its text", give: declared.Left},
			{name: "encodes Right as its text", give: declared.Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				assert.Equal(
					t,
					string(text),
					tt.give.String(),
					"the encoding of the variant is its text",
				)
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			out := declared.Lane(3)
			_, err := out.MarshalText()
			assert.HasError(t, err, "MarshalText refuses the value")
			assert.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Lane
		}{
			{name: "decodes Left from its encoding", give: declared.Left},
			{name: "decodes Right from its encoding", give: declared.Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got declared.Lane
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got declared.Lane
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestCoil checks the functions and the methods of the enumeration Coil.
func TestCoil(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[declared.Coil]bool{}
		for _, v := range []declared.Coil{
			declared.Snug,
			declared.Slack,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("CoilValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []declared.Coil{
					declared.Snug,
					declared.Slack,
				}
				assert.Equal(
					t,
					declared.CoilValues(),
					want,
					"CoilValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := declared.CoilValues(), declared.CoilValues()
			assert.NotEqual(
				t,
				&first[0],
				&second[0],
				"the two calls return two slices",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []declared.Coil{
				declared.Snug,
				declared.Slack,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := declared.Coil(3)
				for _, v := range []declared.Coil{
					declared.Snug,
					declared.Slack,
				} {
					assert.NotEqual(
						t,
						out.String(),
						v.String(),
						"an undeclared value has a text of its own",
					)
				}
			},
		)
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Coil
		}{
			{name: "accepts Snug", give: declared.Snug},
			{name: "accepts Slack", give: declared.Slack},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, declared.Coil(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero declared.Coil
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseCoil", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Coil
		}{
			{name: "returns Snug for its text", give: declared.Snug},
			{name: "returns Slack for its text", give: declared.Slack},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := declared.ParseCoil(tt.give.String())
				assert.NoError(t, err, "ParseCoil accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseCoil returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := declared.ParseCoil("??????")
			assert.ErrorIs(
				t,
				err,
				declared.ErrUnknownCoil,
				"ParseCoil refuses the text with ErrUnknownCoil",
			)
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Coil
		}{
			{name: "encodes Snug as its text", give: declared.Snug},
			{name: "encodes Slack as its text", give: declared.Slack},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				assert.Equal(
					t,
					string(text),
					tt.give.String(),
					"the encoding of the variant is its text",
				)
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			out := declared.Coil(3)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				declared.ErrUnknownCoil,
				"MarshalText refuses the value with ErrUnknownCoil",
			)
			assert.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give declared.Coil
		}{
			{name: "decodes Snug from its encoding", give: declared.Snug},
			{name: "decodes Slack from its encoding", give: declared.Slack},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got declared.Coil
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got declared.Coil
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
