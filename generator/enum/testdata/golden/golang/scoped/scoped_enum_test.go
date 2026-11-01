package scoped_test

import (
	testing2 "testing"

	"example.com/shop/scoped"
	assert2 "go.dokimi.dev/assert"
)

// TestStatus checks the functions and the methods of the enumeration
// Status.
func TestStatus(t *testing2.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing2.T) {
		t.Parallel()
		seen := map[scoped.Status]bool{}
		for _, v := range []scoped.Status{
			scoped.Draft,
			scoped.Published,
		} {
			assert2.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("StatusValues", func(t *testing2.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing2.T) {
				t.Parallel()
				want := []scoped.Status{
					scoped.Draft,
					scoped.Published,
				}
				assert2.Equal(
					t,
					scoped.StatusValues(),
					want,
					"StatusValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing2.T) {
			t.Parallel()
			first, second := scoped.StatusValues(), scoped.StatusValues()
			assert2.NotEqual(
				t,
				&first[0],
				&second[0],
				"the two calls return two slices",
				assert2.ByIdentity(),
			)
		})
	})

	t.Run("String", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing2.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []scoped.Status{
				scoped.Draft,
				scoped.Published,
			} {
				assert2.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing2.T) {
				t.Parallel()
				out := scoped.Status(3)
				for _, v := range []scoped.Status{
					scoped.Draft,
					scoped.Published,
				} {
					assert2.NotEqual(
						t,
						out.String(),
						v.String(),
						"an undeclared value has a text of its own",
					)
				}
			},
		)
	})

	t.Run("IsValid", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Status
		}{
			{name: "accepts Draft", give: scoped.Draft},
			{name: "accepts Published", give: scoped.Published},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				assert2.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing2.T) {
			t.Parallel()
			assert2.False(t, scoped.Status(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing2.T) {
			t.Parallel()
			var zero scoped.Status
			assert2.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseStatus", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Status
		}{
			{name: "returns Draft for its text", give: scoped.Draft},
			{name: "returns Published for its text", give: scoped.Published},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				got, err := scoped.ParseStatus(tt.give.String())
				assert2.NoError(t, err, "ParseStatus accepts the text of the variant")
				assert2.Equal(t, got, tt.give, "ParseStatus returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing2.T) {
			t.Parallel()
			_, err := scoped.ParseStatus("??????????")
			assert2.ErrorIs(
				t,
				err,
				scoped.ErrUnknownStatus,
				"ParseStatus refuses the text with ErrUnknownStatus",
			)
			assert2.Contains(t, err.Error(), "??????????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Status
		}{
			{name: "encodes Draft as its text", give: scoped.Draft},
			{name: "encodes Published as its text", give: scoped.Published},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert2.NoError(t, err, "MarshalText encodes the variant")
				assert2.Equal(
					t,
					string(text),
					tt.give.String(),
					"the encoding of the variant is its text",
				)
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing2.T) {
			t.Parallel()
			out := scoped.Status(3)
			_, err := out.MarshalText()
			assert2.ErrorIs(
				t,
				err,
				scoped.ErrUnknownStatus,
				"MarshalText refuses the value with ErrUnknownStatus",
			)
			assert2.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Status
		}{
			{name: "decodes Draft from its encoding", give: scoped.Draft},
			{name: "decodes Published from its encoding", give: scoped.Published},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert2.NoError(t, err, "MarshalText encodes the variant")
				var got scoped.Status
				assert2.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert2.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing2.T) {
			t.Parallel()
			var got scoped.Status
			assert2.HasError(
				t,
				got.UnmarshalText([]byte("??????????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestLevel checks the functions and the methods of the enumeration Level.
func TestLevel(t *testing2.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing2.T) {
		t.Parallel()
		seen := map[scoped.Level]bool{}
		for _, v := range []scoped.Level{
			scoped.Dim,
			scoped.Bright,
		} {
			assert2.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("LevelValues", func(t *testing2.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing2.T) {
				t.Parallel()
				want := []scoped.Level{
					scoped.Dim,
					scoped.Bright,
				}
				assert2.Equal(
					t,
					scoped.LevelValues(),
					want,
					"LevelValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing2.T) {
			t.Parallel()
			first, second := scoped.LevelValues(), scoped.LevelValues()
			assert2.NotEqual(
				t,
				&first[0],
				&second[0],
				"the two calls return two slices",
				assert2.ByIdentity(),
			)
		})
	})

	t.Run("String", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns a text of its own for each variant", func(t *testing2.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []scoped.Level{
				scoped.Dim,
				scoped.Bright,
			} {
				assert2.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing2.T) {
				t.Parallel()
				out := scoped.Level(2.5)
				for _, v := range []scoped.Level{
					scoped.Dim,
					scoped.Bright,
				} {
					assert2.NotEqual(
						t,
						out.String(),
						v.String(),
						"an undeclared value has a text of its own",
					)
				}
			},
		)
	})

	t.Run("IsValid", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Level
		}{
			{name: "accepts Dim", give: scoped.Dim},
			{name: "accepts Bright", give: scoped.Bright},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				assert2.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing2.T) {
			t.Parallel()
			assert2.False(t, scoped.Level(2.5).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing2.T) {
			t.Parallel()
			var zero scoped.Level
			assert2.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseLevel", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Level
		}{
			{name: "returns Dim for its text", give: scoped.Dim},
			{name: "returns Bright for its text", give: scoped.Bright},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				got, err := scoped.ParseLevel(tt.give.String())
				assert2.NoError(t, err, "ParseLevel accepts the text of the variant")
				assert2.Equal(t, got, tt.give, "ParseLevel returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing2.T) {
			t.Parallel()
			_, err := scoped.ParseLevel("???????")
			assert2.ErrorIs(
				t,
				err,
				scoped.ErrUnknownLevel,
				"ParseLevel refuses the text with ErrUnknownLevel",
			)
			assert2.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Level
		}{
			{name: "encodes Dim as its text", give: scoped.Dim},
			{name: "encodes Bright as its text", give: scoped.Bright},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert2.NoError(t, err, "MarshalText encodes the variant")
				assert2.Equal(
					t,
					string(text),
					tt.give.String(),
					"the encoding of the variant is its text",
				)
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing2.T) {
			t.Parallel()
			out := scoped.Level(2.5)
			_, err := out.MarshalText()
			assert2.ErrorIs(
				t,
				err,
				scoped.ErrUnknownLevel,
				"MarshalText refuses the value with ErrUnknownLevel",
			)
			assert2.Contains(
				t,
				err.Error(),
				out.String(),
				"the error contains the text of the value",
			)
		})
	})

	t.Run("UnmarshalText", func(t *testing2.T) {
		t.Parallel()

		tests := []struct {
			name string
			give scoped.Level
		}{
			{name: "decodes Dim from its encoding", give: scoped.Dim},
			{name: "decodes Bright from its encoding", give: scoped.Bright},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing2.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert2.NoError(t, err, "MarshalText encodes the variant")
				var got scoped.Level
				assert2.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert2.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing2.T) {
			t.Parallel()
			var got scoped.Level
			assert2.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
