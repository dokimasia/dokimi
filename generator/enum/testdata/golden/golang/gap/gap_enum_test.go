package gap_test

import (
	"testing"

	"example.com/shop/gap"
	"go.dokimi.dev/assert"
)

// TestGap checks the functions and the methods of the enumeration Gap.
func TestGap(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[gap.Gap]bool{}
		for _, v := range []gap.Gap{
			gap.Two,
			gap.Three,
			gap.Five,
			gap.Six,
			gap.Eleven,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("GapValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []gap.Gap{
					gap.Two,
					gap.Three,
					gap.Five,
					gap.Six,
					gap.Eleven,
				}
				assert.Equal(t, gap.GapValues(), want, "GapValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := gap.GapValues(), gap.GapValues()
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
			for _, v := range []gap.Gap{
				gap.Two,
				gap.Three,
				gap.Five,
				gap.Six,
				gap.Eleven,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := gap.Gap(12)
				for _, v := range []gap.Gap{
					gap.Two,
					gap.Three,
					gap.Five,
					gap.Six,
					gap.Eleven,
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
			give gap.Gap
		}{
			{name: "accepts Two", give: gap.Two},
			{name: "accepts Three", give: gap.Three},
			{name: "accepts Five", give: gap.Five},
			{name: "accepts Six", give: gap.Six},
			{name: "accepts Eleven", give: gap.Eleven},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, gap.Gap(12).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero gap.Gap
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseGap", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Gap
		}{
			{name: "returns Two for its text", give: gap.Two},
			{name: "returns Three for its text", give: gap.Three},
			{name: "returns Five for its text", give: gap.Five},
			{name: "returns Six for its text", give: gap.Six},
			{name: "returns Eleven for its text", give: gap.Eleven},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := gap.ParseGap(tt.give.String())
				assert.NoError(t, err, "ParseGap accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseGap returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := gap.ParseGap("???????")
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownGap,
				"ParseGap refuses the text with ErrUnknownGap",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Gap
		}{
			{name: "encodes Two as its text", give: gap.Two},
			{name: "encodes Three as its text", give: gap.Three},
			{name: "encodes Five as its text", give: gap.Five},
			{name: "encodes Six as its text", give: gap.Six},
			{name: "encodes Eleven as its text", give: gap.Eleven},
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
			out := gap.Gap(12)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownGap,
				"MarshalText refuses the value with ErrUnknownGap",
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
			give gap.Gap
		}{
			{name: "decodes Two from its encoding", give: gap.Two},
			{name: "decodes Three from its encoding", give: gap.Three},
			{name: "decodes Five from its encoding", give: gap.Five},
			{name: "decodes Six from its encoding", give: gap.Six},
			{name: "decodes Eleven from its encoding", give: gap.Eleven},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got gap.Gap
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got gap.Gap
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestUnum checks the functions and the methods of the enumeration Unum.
func TestUnum(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[gap.Unum]bool{}
		for _, v := range []gap.Unum{
			gap.Naught,
			gap.One,
			gap.Far,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("UnumValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []gap.Unum{
					gap.Naught,
					gap.One,
					gap.Far,
				}
				assert.Equal(t, gap.UnumValues(), want, "UnumValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := gap.UnumValues(), gap.UnumValues()
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
			for _, v := range []gap.Unum{
				gap.Naught,
				gap.One,
				gap.Far,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := gap.Unum(254)
				for _, v := range []gap.Unum{
					gap.Naught,
					gap.One,
					gap.Far,
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
			give gap.Unum
		}{
			{name: "accepts Naught", give: gap.Naught},
			{name: "accepts One", give: gap.One},
			{name: "accepts Far", give: gap.Far},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, gap.Unum(254).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Naught", func(t *testing.T) {
			t.Parallel()
			var zero gap.Unum
			assert.Equal(t, zero, gap.Naught, "the zero value is Naught")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseUnum", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Unum
		}{
			{name: "returns Naught for its text", give: gap.Naught},
			{name: "returns One for its text", give: gap.One},
			{name: "returns Far for its text", give: gap.Far},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := gap.ParseUnum(tt.give.String())
				assert.NoError(t, err, "ParseUnum accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseUnum returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := gap.ParseUnum("???????")
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownUnum,
				"ParseUnum refuses the text with ErrUnknownUnum",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Unum
		}{
			{name: "encodes Naught as its text", give: gap.Naught},
			{name: "encodes One as its text", give: gap.One},
			{name: "encodes Far as its text", give: gap.Far},
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
			out := gap.Unum(254)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownUnum,
				"MarshalText refuses the value with ErrUnknownUnum",
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
			give gap.Unum
		}{
			{name: "decodes Naught from its encoding", give: gap.Naught},
			{name: "decodes One from its encoding", give: gap.One},
			{name: "decodes Far from its encoding", give: gap.Far},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got gap.Unum
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got gap.Unum
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestPrime checks the functions and the methods of the enumeration Prime.
func TestPrime(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[gap.Prime]bool{}
		for _, v := range []gap.Prime{
			gap.P2,
			gap.P3,
			gap.P5,
			gap.P7,
			gap.P11,
			gap.P13,
			gap.P17,
			gap.P19,
			gap.P23,
			gap.P29,
			gap.P31,
			gap.P37,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("PrimeValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []gap.Prime{
					gap.P2,
					gap.P3,
					gap.P5,
					gap.P7,
					gap.P11,
					gap.P13,
					gap.P17,
					gap.P19,
					gap.P23,
					gap.P29,
					gap.P31,
					gap.P37,
				}
				assert.Equal(
					t,
					gap.PrimeValues(),
					want,
					"PrimeValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := gap.PrimeValues(), gap.PrimeValues()
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
			for _, v := range []gap.Prime{
				gap.P2,
				gap.P3,
				gap.P5,
				gap.P7,
				gap.P11,
				gap.P13,
				gap.P17,
				gap.P19,
				gap.P23,
				gap.P29,
				gap.P31,
				gap.P37,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := gap.Prime(38)
				for _, v := range []gap.Prime{
					gap.P2,
					gap.P3,
					gap.P5,
					gap.P7,
					gap.P11,
					gap.P13,
					gap.P17,
					gap.P19,
					gap.P23,
					gap.P29,
					gap.P31,
					gap.P37,
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
			give gap.Prime
		}{
			{name: "accepts P2", give: gap.P2},
			{name: "accepts P3", give: gap.P3},
			{name: "accepts P5", give: gap.P5},
			{name: "accepts P7", give: gap.P7},
			{name: "accepts P11", give: gap.P11},
			{name: "accepts P13", give: gap.P13},
			{name: "accepts P17", give: gap.P17},
			{name: "accepts P19", give: gap.P19},
			{name: "accepts P23", give: gap.P23},
			{name: "accepts P29", give: gap.P29},
			{name: "accepts P31", give: gap.P31},
			{name: "accepts P37", give: gap.P37},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, gap.Prime(38).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero gap.Prime
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParsePrime", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Prime
		}{
			{name: "returns P2 for its text", give: gap.P2},
			{name: "returns P3 for its text", give: gap.P3},
			{name: "returns P5 for its text", give: gap.P5},
			{name: "returns P7 for its text", give: gap.P7},
			{name: "returns P11 for its text", give: gap.P11},
			{name: "returns P13 for its text", give: gap.P13},
			{name: "returns P17 for its text", give: gap.P17},
			{name: "returns P19 for its text", give: gap.P19},
			{name: "returns P23 for its text", give: gap.P23},
			{name: "returns P29 for its text", give: gap.P29},
			{name: "returns P31 for its text", give: gap.P31},
			{name: "returns P37 for its text", give: gap.P37},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := gap.ParsePrime(tt.give.String())
				assert.NoError(t, err, "ParsePrime accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParsePrime returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := gap.ParsePrime("????")
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownPrime,
				"ParsePrime refuses the text with ErrUnknownPrime",
			)
			assert.Contains(t, err.Error(), "????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Prime
		}{
			{name: "encodes P2 as its text", give: gap.P2},
			{name: "encodes P3 as its text", give: gap.P3},
			{name: "encodes P5 as its text", give: gap.P5},
			{name: "encodes P7 as its text", give: gap.P7},
			{name: "encodes P11 as its text", give: gap.P11},
			{name: "encodes P13 as its text", give: gap.P13},
			{name: "encodes P17 as its text", give: gap.P17},
			{name: "encodes P19 as its text", give: gap.P19},
			{name: "encodes P23 as its text", give: gap.P23},
			{name: "encodes P29 as its text", give: gap.P29},
			{name: "encodes P31 as its text", give: gap.P31},
			{name: "encodes P37 as its text", give: gap.P37},
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
			out := gap.Prime(38)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownPrime,
				"MarshalText refuses the value with ErrUnknownPrime",
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
			give gap.Prime
		}{
			{name: "decodes P2 from its encoding", give: gap.P2},
			{name: "decodes P3 from its encoding", give: gap.P3},
			{name: "decodes P5 from its encoding", give: gap.P5},
			{name: "decodes P7 from its encoding", give: gap.P7},
			{name: "decodes P11 from its encoding", give: gap.P11},
			{name: "decodes P13 from its encoding", give: gap.P13},
			{name: "decodes P17 from its encoding", give: gap.P17},
			{name: "decodes P19 from its encoding", give: gap.P19},
			{name: "decodes P23 from its encoding", give: gap.P23},
			{name: "decodes P29 from its encoding", give: gap.P29},
			{name: "decodes P31 from its encoding", give: gap.P31},
			{name: "decodes P37 from its encoding", give: gap.P37},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got gap.Prime
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got gap.Prime
			assert.HasError(t, got.UnmarshalText([]byte("????")), "UnmarshalText refuses the text")
		})
	})
}

// TestRank checks the functions and the methods of the enumeration Rank.
func TestRank(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[gap.Rank]bool{}
		for _, v := range []gap.Rank{
			gap.Gold,
			gap.Silver,
			gap.Bronze,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("RankValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []gap.Rank{
					gap.Gold,
					gap.Silver,
					gap.Bronze,
				}
				assert.Equal(t, gap.RankValues(), want, "RankValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := gap.RankValues(), gap.RankValues()
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
			for _, v := range []gap.Rank{
				gap.Gold,
				gap.Silver,
				gap.Bronze,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := gap.Rank(4)
				for _, v := range []gap.Rank{
					gap.Gold,
					gap.Silver,
					gap.Bronze,
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
			give gap.Rank
		}{
			{name: "accepts Gold", give: gap.Gold},
			{name: "accepts Silver", give: gap.Silver},
			{name: "accepts Bronze", give: gap.Bronze},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, gap.Rank(4).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero gap.Rank
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseRank", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Rank
		}{
			{name: "returns Gold for its text", give: gap.Gold},
			{name: "returns Silver for its text", give: gap.Silver},
			{name: "returns Bronze for its text", give: gap.Bronze},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := gap.ParseRank(tt.give.String())
				assert.NoError(t, err, "ParseRank accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseRank returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := gap.ParseRank("???????")
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownRank,
				"ParseRank refuses the text with ErrUnknownRank",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Rank
		}{
			{name: "encodes Gold as its text", give: gap.Gold},
			{name: "encodes Silver as its text", give: gap.Silver},
			{name: "encodes Bronze as its text", give: gap.Bronze},
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
			out := gap.Rank(4)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownRank,
				"MarshalText refuses the value with ErrUnknownRank",
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
			give gap.Rank
		}{
			{name: "decodes Gold from its encoding", give: gap.Gold},
			{name: "decodes Silver from its encoding", give: gap.Silver},
			{name: "decodes Bronze from its encoding", give: gap.Bronze},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got gap.Rank
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got gap.Rank
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestOdd checks the functions and the methods of the enumeration Odd.
func TestOdd(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[gap.Odd]bool{}
		for _, v := range []gap.Odd{
			gap.O1,
			gap.O3,
			gap.O5,
			gap.O7,
			gap.O9,
			gap.O11,
			gap.O13,
			gap.O15,
			gap.O17,
			gap.O19,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("OddValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []gap.Odd{
					gap.O1,
					gap.O3,
					gap.O5,
					gap.O7,
					gap.O9,
					gap.O11,
					gap.O13,
					gap.O15,
					gap.O17,
					gap.O19,
				}
				assert.Equal(t, gap.OddValues(), want, "OddValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := gap.OddValues(), gap.OddValues()
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
			for _, v := range []gap.Odd{
				gap.O1,
				gap.O3,
				gap.O5,
				gap.O7,
				gap.O9,
				gap.O11,
				gap.O13,
				gap.O15,
				gap.O17,
				gap.O19,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := gap.Odd(20)
				for _, v := range []gap.Odd{
					gap.O1,
					gap.O3,
					gap.O5,
					gap.O7,
					gap.O9,
					gap.O11,
					gap.O13,
					gap.O15,
					gap.O17,
					gap.O19,
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
			give gap.Odd
		}{
			{name: "accepts O1", give: gap.O1},
			{name: "accepts O3", give: gap.O3},
			{name: "accepts O5", give: gap.O5},
			{name: "accepts O7", give: gap.O7},
			{name: "accepts O9", give: gap.O9},
			{name: "accepts O11", give: gap.O11},
			{name: "accepts O13", give: gap.O13},
			{name: "accepts O15", give: gap.O15},
			{name: "accepts O17", give: gap.O17},
			{name: "accepts O19", give: gap.O19},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, gap.Odd(20).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero gap.Odd
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseOdd", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Odd
		}{
			{name: "returns O1 for its text", give: gap.O1},
			{name: "returns O3 for its text", give: gap.O3},
			{name: "returns O5 for its text", give: gap.O5},
			{name: "returns O7 for its text", give: gap.O7},
			{name: "returns O9 for its text", give: gap.O9},
			{name: "returns O11 for its text", give: gap.O11},
			{name: "returns O13 for its text", give: gap.O13},
			{name: "returns O15 for its text", give: gap.O15},
			{name: "returns O17 for its text", give: gap.O17},
			{name: "returns O19 for its text", give: gap.O19},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := gap.ParseOdd(tt.give.String())
				assert.NoError(t, err, "ParseOdd accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseOdd returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := gap.ParseOdd("????")
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownOdd,
				"ParseOdd refuses the text with ErrUnknownOdd",
			)
			assert.Contains(t, err.Error(), "????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Odd
		}{
			{name: "encodes O1 as its text", give: gap.O1},
			{name: "encodes O3 as its text", give: gap.O3},
			{name: "encodes O5 as its text", give: gap.O5},
			{name: "encodes O7 as its text", give: gap.O7},
			{name: "encodes O9 as its text", give: gap.O9},
			{name: "encodes O11 as its text", give: gap.O11},
			{name: "encodes O13 as its text", give: gap.O13},
			{name: "encodes O15 as its text", give: gap.O15},
			{name: "encodes O17 as its text", give: gap.O17},
			{name: "encodes O19 as its text", give: gap.O19},
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
			out := gap.Odd(20)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownOdd,
				"MarshalText refuses the value with ErrUnknownOdd",
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
			give gap.Odd
		}{
			{name: "decodes O1 from its encoding", give: gap.O1},
			{name: "decodes O3 from its encoding", give: gap.O3},
			{name: "decodes O5 from its encoding", give: gap.O5},
			{name: "decodes O7 from its encoding", give: gap.O7},
			{name: "decodes O9 from its encoding", give: gap.O9},
			{name: "decodes O11 from its encoding", give: gap.O11},
			{name: "decodes O13 from its encoding", give: gap.O13},
			{name: "decodes O15 from its encoding", give: gap.O15},
			{name: "decodes O17 from its encoding", give: gap.O17},
			{name: "decodes O19 from its encoding", give: gap.O19},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got gap.Odd
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got gap.Odd
			assert.HasError(t, got.UnmarshalText([]byte("????")), "UnmarshalText refuses the text")
		})
	})
}

// TestNum checks the functions and the methods of the enumeration Num.
func TestNum(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[gap.Num]bool{}
		for _, v := range []gap.Num{
			gap.Minus,
			gap.Zero,
			gap.Plus,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("NumValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []gap.Num{
					gap.Minus,
					gap.Zero,
					gap.Plus,
				}
				assert.Equal(t, gap.NumValues(), want, "NumValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := gap.NumValues(), gap.NumValues()
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
			for _, v := range []gap.Num{
				gap.Minus,
				gap.Zero,
				gap.Plus,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := gap.Num(2)
				for _, v := range []gap.Num{
					gap.Minus,
					gap.Zero,
					gap.Plus,
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
			give gap.Num
		}{
			{name: "accepts Minus", give: gap.Minus},
			{name: "accepts Zero", give: gap.Zero},
			{name: "accepts Plus", give: gap.Plus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, gap.Num(2).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Zero", func(t *testing.T) {
			t.Parallel()
			var zero gap.Num
			assert.Equal(t, zero, gap.Zero, "the zero value is Zero")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseNum", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Num
		}{
			{name: "returns Minus for its text", give: gap.Minus},
			{name: "returns Zero for its text", give: gap.Zero},
			{name: "returns Plus for its text", give: gap.Plus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := gap.ParseNum(tt.give.String())
				assert.NoError(t, err, "ParseNum accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseNum returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := gap.ParseNum("??????")
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownNum,
				"ParseNum refuses the text with ErrUnknownNum",
			)
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give gap.Num
		}{
			{name: "encodes Minus as its text", give: gap.Minus},
			{name: "encodes Zero as its text", give: gap.Zero},
			{name: "encodes Plus as its text", give: gap.Plus},
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
			out := gap.Num(2)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				gap.ErrUnknownNum,
				"MarshalText refuses the value with ErrUnknownNum",
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
			give gap.Num
		}{
			{name: "decodes Minus from its encoding", give: gap.Minus},
			{name: "decodes Zero from its encoding", give: gap.Zero},
			{name: "decodes Plus from its encoding", give: gap.Plus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got gap.Num
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got gap.Num
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
