package shaped_test

import (
	"testing"

	"example.com/shop/shaped"
	"go.dokimi.dev/assert"
)

// TestRatio checks the functions and the methods of the enumeration Ratio.
func TestRatio(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[shaped.Ratio]bool{}
		for _, v := range []shaped.Ratio{
			shaped.Single,
			shaped.Double,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("RatioValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []shaped.Ratio{
					shaped.Single,
					shaped.Double,
				}
				assert.Equal(
					t,
					shaped.RatioValues(),
					want,
					"RatioValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := shaped.RatioValues(), shaped.RatioValues()
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
			for _, v := range []shaped.Ratio{
				shaped.Single,
				shaped.Double,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := shaped.Ratio(3)
				for _, v := range []shaped.Ratio{
					shaped.Single,
					shaped.Double,
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
			give shaped.Ratio
		}{
			{name: "accepts Single", give: shaped.Single},
			{name: "accepts Double", give: shaped.Double},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, shaped.Ratio(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero shaped.Ratio
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseRatio", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give shaped.Ratio
		}{
			{name: "returns Single for its text", give: shaped.Single},
			{name: "returns Double for its text", give: shaped.Double},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := shaped.ParseRatio(tt.give.String())
				assert.NoError(t, err, "ParseRatio accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseRatio returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := shaped.ParseRatio("???????")
			assert.ErrorIs(
				t,
				err,
				shaped.ErrUnknownRatio,
				"ParseRatio refuses the text with ErrUnknownRatio",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give shaped.Ratio
		}{
			{name: "encodes Single as its text", give: shaped.Single},
			{name: "encodes Double as its text", give: shaped.Double},
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
			out := shaped.Ratio(3)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				shaped.ErrUnknownRatio,
				"MarshalText refuses the value with ErrUnknownRatio",
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
			give shaped.Ratio
		}{
			{name: "decodes Single from its encoding", give: shaped.Single},
			{name: "decodes Double from its encoding", give: shaped.Double},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got shaped.Ratio
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got shaped.Ratio
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestTier checks the functions and the methods of the enumeration Tier.
func TestTier(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[shaped.Tier]bool{}
		for _, v := range []shaped.Tier{
			shaped.Bronze,
			shaped.Silver,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("TierValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []shaped.Tier{
					shaped.Bronze,
					shaped.Silver,
				}
				assert.Equal(
					t,
					shaped.TierValues(),
					want,
					"TierValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := shaped.TierValues(), shaped.TierValues()
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
			for _, v := range []shaped.Tier{
				shaped.Bronze,
				shaped.Silver,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := shaped.Tier(3)
				for _, v := range []shaped.Tier{
					shaped.Bronze,
					shaped.Silver,
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
			give shaped.Tier
		}{
			{name: "accepts Bronze", give: shaped.Bronze},
			{name: "accepts Silver", give: shaped.Silver},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, shaped.Tier(3).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero shaped.Tier
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseTier", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give shaped.Tier
		}{
			{name: "returns Bronze for its text", give: shaped.Bronze},
			{name: "returns Silver for its text", give: shaped.Silver},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := shaped.ParseTier(tt.give.String())
				assert.NoError(t, err, "ParseTier accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseTier returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := shaped.ParseTier("???????")
			assert.ErrorIs(
				t,
				err,
				shaped.ErrUnknownTier,
				"ParseTier refuses the text with ErrUnknownTier",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give shaped.Tier
		}{
			{name: "encodes Bronze as its text", give: shaped.Bronze},
			{name: "encodes Silver as its text", give: shaped.Silver},
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
			out := shaped.Tier(3)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				shaped.ErrUnknownTier,
				"MarshalText refuses the value with ErrUnknownTier",
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
			give shaped.Tier
		}{
			{name: "decodes Bronze from its encoding", give: shaped.Bronze},
			{name: "decodes Silver from its encoding", give: shaped.Silver},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got shaped.Tier
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got shaped.Tier
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
