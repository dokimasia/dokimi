package gap

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestGap checks the functions and the methods of the enumeration Gap.
func TestGap(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Gap]bool{}
		for _, v := range []Gap{
			Two,
			Three,
			Five,
			Six,
			Eleven,
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
				want := []Gap{
					Two,
					Three,
					Five,
					Six,
					Eleven,
				}
				assert.Equal(t, GapValues(), want, "GapValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := GapValues(), GapValues()
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
			for _, v := range []Gap{
				Two,
				Three,
				Five,
				Six,
				Eleven,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Gap(12)
				for _, v := range []Gap{
					Two,
					Three,
					Five,
					Six,
					Eleven,
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
			give Gap
		}{
			{name: "accepts Two", give: Two},
			{name: "accepts Three", give: Three},
			{name: "accepts Five", give: Five},
			{name: "accepts Six", give: Six},
			{name: "accepts Eleven", give: Eleven},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Gap(12).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Gap
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseGap", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Gap
		}{
			{name: "returns Two for its text", give: Two},
			{name: "returns Three for its text", give: Three},
			{name: "returns Five for its text", give: Five},
			{name: "returns Six for its text", give: Six},
			{name: "returns Eleven for its text", give: Eleven},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseGap(tt.give.String())
				assert.NoError(t, err, "ParseGap accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseGap returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseGap("???????")
			assert.ErrorIs(t, err, ErrUnknownGap, "ParseGap refuses the text with ErrUnknownGap")
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Gap
		}{
			{name: "encodes Two as its text", give: Two},
			{name: "encodes Three as its text", give: Three},
			{name: "encodes Five as its text", give: Five},
			{name: "encodes Six as its text", give: Six},
			{name: "encodes Eleven as its text", give: Eleven},
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
			out := Gap(12)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownGap,
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
			give Gap
		}{
			{name: "decodes Two from its encoding", give: Two},
			{name: "decodes Three from its encoding", give: Three},
			{name: "decodes Five from its encoding", give: Five},
			{name: "decodes Six from its encoding", give: Six},
			{name: "decodes Eleven from its encoding", give: Eleven},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Gap
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Gap
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
		seen := map[Unum]bool{}
		for _, v := range []Unum{
			Naught,
			One,
			Far,
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
				want := []Unum{
					Naught,
					One,
					Far,
				}
				assert.Equal(t, UnumValues(), want, "UnumValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := UnumValues(), UnumValues()
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
			for _, v := range []Unum{
				Naught,
				One,
				Far,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Unum(254)
				for _, v := range []Unum{
					Naught,
					One,
					Far,
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
			give Unum
		}{
			{name: "accepts Naught", give: Naught},
			{name: "accepts One", give: One},
			{name: "accepts Far", give: Far},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Unum(254).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Naught", func(t *testing.T) {
			t.Parallel()
			var zero Unum
			assert.Equal(t, zero, Naught, "the zero value is Naught")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseUnum", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Unum
		}{
			{name: "returns Naught for its text", give: Naught},
			{name: "returns One for its text", give: One},
			{name: "returns Far for its text", give: Far},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseUnum(tt.give.String())
				assert.NoError(t, err, "ParseUnum accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseUnum returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseUnum("???????")
			assert.ErrorIs(t, err, ErrUnknownUnum, "ParseUnum refuses the text with ErrUnknownUnum")
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Unum
		}{
			{name: "encodes Naught as its text", give: Naught},
			{name: "encodes One as its text", give: One},
			{name: "encodes Far as its text", give: Far},
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
			out := Unum(254)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownUnum,
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
			give Unum
		}{
			{name: "decodes Naught from its encoding", give: Naught},
			{name: "decodes One from its encoding", give: One},
			{name: "decodes Far from its encoding", give: Far},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Unum
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Unum
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
		seen := map[Prime]bool{}
		for _, v := range []Prime{
			P2,
			P3,
			P5,
			P7,
			P11,
			P13,
			P17,
			P19,
			P23,
			P29,
			P31,
			P37,
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
				want := []Prime{
					P2,
					P3,
					P5,
					P7,
					P11,
					P13,
					P17,
					P19,
					P23,
					P29,
					P31,
					P37,
				}
				assert.Equal(t, PrimeValues(), want, "PrimeValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := PrimeValues(), PrimeValues()
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
			for _, v := range []Prime{
				P2,
				P3,
				P5,
				P7,
				P11,
				P13,
				P17,
				P19,
				P23,
				P29,
				P31,
				P37,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Prime(38)
				for _, v := range []Prime{
					P2,
					P3,
					P5,
					P7,
					P11,
					P13,
					P17,
					P19,
					P23,
					P29,
					P31,
					P37,
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
			give Prime
		}{
			{name: "accepts P2", give: P2},
			{name: "accepts P3", give: P3},
			{name: "accepts P5", give: P5},
			{name: "accepts P7", give: P7},
			{name: "accepts P11", give: P11},
			{name: "accepts P13", give: P13},
			{name: "accepts P17", give: P17},
			{name: "accepts P19", give: P19},
			{name: "accepts P23", give: P23},
			{name: "accepts P29", give: P29},
			{name: "accepts P31", give: P31},
			{name: "accepts P37", give: P37},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Prime(38).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Prime
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParsePrime", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Prime
		}{
			{name: "returns P2 for its text", give: P2},
			{name: "returns P3 for its text", give: P3},
			{name: "returns P5 for its text", give: P5},
			{name: "returns P7 for its text", give: P7},
			{name: "returns P11 for its text", give: P11},
			{name: "returns P13 for its text", give: P13},
			{name: "returns P17 for its text", give: P17},
			{name: "returns P19 for its text", give: P19},
			{name: "returns P23 for its text", give: P23},
			{name: "returns P29 for its text", give: P29},
			{name: "returns P31 for its text", give: P31},
			{name: "returns P37 for its text", give: P37},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParsePrime(tt.give.String())
				assert.NoError(t, err, "ParsePrime accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParsePrime returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParsePrime("????")
			assert.ErrorIs(
				t,
				err,
				ErrUnknownPrime,
				"ParsePrime refuses the text with ErrUnknownPrime",
			)
			assert.Contains(t, err.Error(), "????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Prime
		}{
			{name: "encodes P2 as its text", give: P2},
			{name: "encodes P3 as its text", give: P3},
			{name: "encodes P5 as its text", give: P5},
			{name: "encodes P7 as its text", give: P7},
			{name: "encodes P11 as its text", give: P11},
			{name: "encodes P13 as its text", give: P13},
			{name: "encodes P17 as its text", give: P17},
			{name: "encodes P19 as its text", give: P19},
			{name: "encodes P23 as its text", give: P23},
			{name: "encodes P29 as its text", give: P29},
			{name: "encodes P31 as its text", give: P31},
			{name: "encodes P37 as its text", give: P37},
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
			out := Prime(38)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownPrime,
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
			give Prime
		}{
			{name: "decodes P2 from its encoding", give: P2},
			{name: "decodes P3 from its encoding", give: P3},
			{name: "decodes P5 from its encoding", give: P5},
			{name: "decodes P7 from its encoding", give: P7},
			{name: "decodes P11 from its encoding", give: P11},
			{name: "decodes P13 from its encoding", give: P13},
			{name: "decodes P17 from its encoding", give: P17},
			{name: "decodes P19 from its encoding", give: P19},
			{name: "decodes P23 from its encoding", give: P23},
			{name: "decodes P29 from its encoding", give: P29},
			{name: "decodes P31 from its encoding", give: P31},
			{name: "decodes P37 from its encoding", give: P37},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Prime
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Prime
			assert.HasError(t, got.UnmarshalText([]byte("????")), "UnmarshalText refuses the text")
		})
	})
}

// TestRank checks the functions and the methods of the enumeration Rank.
func TestRank(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Rank]bool{}
		for _, v := range []Rank{
			Gold,
			Silver,
			Bronze,
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
				want := []Rank{
					Gold,
					Silver,
					Bronze,
				}
				assert.Equal(t, RankValues(), want, "RankValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := RankValues(), RankValues()
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
			for _, v := range []Rank{
				Gold,
				Silver,
				Bronze,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Rank(4)
				for _, v := range []Rank{
					Gold,
					Silver,
					Bronze,
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
			give Rank
		}{
			{name: "accepts Gold", give: Gold},
			{name: "accepts Silver", give: Silver},
			{name: "accepts Bronze", give: Bronze},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Rank(4).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Rank
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseRank", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Rank
		}{
			{name: "returns Gold for its text", give: Gold},
			{name: "returns Silver for its text", give: Silver},
			{name: "returns Bronze for its text", give: Bronze},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseRank(tt.give.String())
				assert.NoError(t, err, "ParseRank accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseRank returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseRank("???????")
			assert.ErrorIs(t, err, ErrUnknownRank, "ParseRank refuses the text with ErrUnknownRank")
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Rank
		}{
			{name: "encodes Gold as its text", give: Gold},
			{name: "encodes Silver as its text", give: Silver},
			{name: "encodes Bronze as its text", give: Bronze},
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
			out := Rank(4)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownRank,
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
			give Rank
		}{
			{name: "decodes Gold from its encoding", give: Gold},
			{name: "decodes Silver from its encoding", give: Silver},
			{name: "decodes Bronze from its encoding", give: Bronze},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Rank
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Rank
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
		seen := map[Odd]bool{}
		for _, v := range []Odd{
			O1,
			O3,
			O5,
			O7,
			O9,
			O11,
			O13,
			O15,
			O17,
			O19,
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
				want := []Odd{
					O1,
					O3,
					O5,
					O7,
					O9,
					O11,
					O13,
					O15,
					O17,
					O19,
				}
				assert.Equal(t, OddValues(), want, "OddValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := OddValues(), OddValues()
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
			for _, v := range []Odd{
				O1,
				O3,
				O5,
				O7,
				O9,
				O11,
				O13,
				O15,
				O17,
				O19,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Odd(20)
				for _, v := range []Odd{
					O1,
					O3,
					O5,
					O7,
					O9,
					O11,
					O13,
					O15,
					O17,
					O19,
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
			give Odd
		}{
			{name: "accepts O1", give: O1},
			{name: "accepts O3", give: O3},
			{name: "accepts O5", give: O5},
			{name: "accepts O7", give: O7},
			{name: "accepts O9", give: O9},
			{name: "accepts O11", give: O11},
			{name: "accepts O13", give: O13},
			{name: "accepts O15", give: O15},
			{name: "accepts O17", give: O17},
			{name: "accepts O19", give: O19},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Odd(20).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Odd
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseOdd", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Odd
		}{
			{name: "returns O1 for its text", give: O1},
			{name: "returns O3 for its text", give: O3},
			{name: "returns O5 for its text", give: O5},
			{name: "returns O7 for its text", give: O7},
			{name: "returns O9 for its text", give: O9},
			{name: "returns O11 for its text", give: O11},
			{name: "returns O13 for its text", give: O13},
			{name: "returns O15 for its text", give: O15},
			{name: "returns O17 for its text", give: O17},
			{name: "returns O19 for its text", give: O19},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseOdd(tt.give.String())
				assert.NoError(t, err, "ParseOdd accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseOdd returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseOdd("????")
			assert.ErrorIs(t, err, ErrUnknownOdd, "ParseOdd refuses the text with ErrUnknownOdd")
			assert.Contains(t, err.Error(), "????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Odd
		}{
			{name: "encodes O1 as its text", give: O1},
			{name: "encodes O3 as its text", give: O3},
			{name: "encodes O5 as its text", give: O5},
			{name: "encodes O7 as its text", give: O7},
			{name: "encodes O9 as its text", give: O9},
			{name: "encodes O11 as its text", give: O11},
			{name: "encodes O13 as its text", give: O13},
			{name: "encodes O15 as its text", give: O15},
			{name: "encodes O17 as its text", give: O17},
			{name: "encodes O19 as its text", give: O19},
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
			out := Odd(20)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownOdd,
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
			give Odd
		}{
			{name: "decodes O1 from its encoding", give: O1},
			{name: "decodes O3 from its encoding", give: O3},
			{name: "decodes O5 from its encoding", give: O5},
			{name: "decodes O7 from its encoding", give: O7},
			{name: "decodes O9 from its encoding", give: O9},
			{name: "decodes O11 from its encoding", give: O11},
			{name: "decodes O13 from its encoding", give: O13},
			{name: "decodes O15 from its encoding", give: O15},
			{name: "decodes O17 from its encoding", give: O17},
			{name: "decodes O19 from its encoding", give: O19},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Odd
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Odd
			assert.HasError(t, got.UnmarshalText([]byte("????")), "UnmarshalText refuses the text")
		})
	})
}

// TestNum checks the functions and the methods of the enumeration Num.
func TestNum(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Num]bool{}
		for _, v := range []Num{
			Minus,
			Zero,
			Plus,
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
				want := []Num{
					Minus,
					Zero,
					Plus,
				}
				assert.Equal(t, NumValues(), want, "NumValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := NumValues(), NumValues()
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
			for _, v := range []Num{
				Minus,
				Zero,
				Plus,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Num(2)
				for _, v := range []Num{
					Minus,
					Zero,
					Plus,
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
			give Num
		}{
			{name: "accepts Minus", give: Minus},
			{name: "accepts Zero", give: Zero},
			{name: "accepts Plus", give: Plus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Num(2).IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Zero", func(t *testing.T) {
			t.Parallel()
			var zero Num
			assert.Equal(t, zero, Zero, "the zero value is Zero")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseNum", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Num
		}{
			{name: "returns Minus for its text", give: Minus},
			{name: "returns Zero for its text", give: Zero},
			{name: "returns Plus for its text", give: Plus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseNum(tt.give.String())
				assert.NoError(t, err, "ParseNum accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseNum returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseNum("??????")
			assert.ErrorIs(t, err, ErrUnknownNum, "ParseNum refuses the text with ErrUnknownNum")
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Num
		}{
			{name: "encodes Minus as its text", give: Minus},
			{name: "encodes Zero as its text", give: Zero},
			{name: "encodes Plus as its text", give: Plus},
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
			out := Num(2)
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownNum,
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
			give Num
		}{
			{name: "decodes Minus from its encoding", give: Minus},
			{name: "decodes Zero from its encoding", give: Zero},
			{name: "decodes Plus from its encoding", give: Plus},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Num
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Num
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
