package region

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestRegion checks the functions and the methods of the enumeration
// Region.
func TestRegion(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Region]bool{}
		for _, v := range []Region{
			US,
			EU,
			AP,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("RegionValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []Region{
					US,
					EU,
					AP,
				}
				assert.Equal(t, RegionValues(), want, "RegionValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := RegionValues(), RegionValues()
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

		tests := []struct {
			name string
			give Region
			want string
		}{
			{name: "returns the declared text of US", give: US, want: "us-east"},
			{name: "returns the declared text of EU", give: EU, want: "eu-west"},
			{name: "returns the declared text of AP", give: AP, want: "ap-south"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(
					t,
					tt.give.String(),
					tt.want,
					"String returns the declared text of the variant",
				)
			})
		}

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []Region{
				US,
				EU,
				AP,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Region("")
				for _, v := range []Region{
					US,
					EU,
					AP,
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
			give Region
		}{
			{name: "accepts US", give: US},
			{name: "accepts EU", give: EU},
			{name: "accepts AP", give: AP},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Region("").IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Region
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseRegion", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Region
		}{
			{name: "returns US for its text", give: US},
			{name: "returns EU for its text", give: EU},
			{name: "returns AP for its text", give: AP},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseRegion(tt.give.String())
				assert.NoError(t, err, "ParseRegion accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseRegion returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseRegion("?????????")
			assert.ErrorIs(
				t,
				err,
				ErrUnknownRegion,
				"ParseRegion refuses the text with ErrUnknownRegion",
			)
			assert.Contains(t, err.Error(), "?????????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Region
		}{
			{name: "encodes US as its text", give: US},
			{name: "encodes EU as its text", give: EU},
			{name: "encodes AP as its text", give: AP},
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
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Region
		}{
			{name: "decodes US from its encoding", give: US},
			{name: "decodes EU from its encoding", give: EU},
			{name: "decodes AP from its encoding", give: AP},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Region
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Region
			assert.HasError(
				t,
				got.UnmarshalText([]byte("?????????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}

// TestZone checks the functions and the methods of the enumeration Zone.
func TestZone(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Zone]bool{}
		for _, v := range []Zone{
			North,
			South,
			Unset,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("ZoneValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []Zone{
					North,
					South,
					Unset,
				}
				assert.Equal(t, ZoneValues(), want, "ZoneValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := ZoneValues(), ZoneValues()
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

		tests := []struct {
			name string
			give Zone
			want string
		}{
			{name: "returns the declared text of North", give: North, want: "north"},
			{name: "returns the declared text of South", give: South, want: "south"},
			{name: "returns the declared text of Unset", give: Unset, want: ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(
					t,
					tt.give.String(),
					tt.want,
					"String returns the declared text of the variant",
				)
			})
		}

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []Zone{
				North,
				South,
				Unset,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Zone("??????")
				for _, v := range []Zone{
					North,
					South,
					Unset,
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
			give Zone
		}{
			{name: "accepts North", give: North},
			{name: "accepts South", give: South},
			{name: "accepts Unset", give: Unset},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Zone("??????").IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Unset", func(t *testing.T) {
			t.Parallel()
			var zero Zone
			assert.Equal(t, zero, Unset, "the zero value is Unset")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseZone", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Zone
		}{
			{name: "returns North for its text", give: North},
			{name: "returns South for its text", give: South},
			{name: "returns Unset for its text", give: Unset},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseZone(tt.give.String())
				assert.NoError(t, err, "ParseZone accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseZone returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseZone("??????")
			assert.ErrorIs(t, err, ErrUnknownZone, "ParseZone refuses the text with ErrUnknownZone")
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Zone
		}{
			{name: "encodes North as its text", give: North},
			{name: "encodes South as its text", give: South},
			{name: "encodes Unset as its text", give: Unset},
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
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Zone
		}{
			{name: "decodes North from its encoding", give: North},
			{name: "decodes South from its encoding", give: South},
			{name: "decodes Unset from its encoding", give: Unset},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Zone
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Zone
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
		seen := map[Lane]bool{}
		for _, v := range []Lane{
			Left,
			Right,
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
				want := []Lane{
					Left,
					Right,
				}
				assert.Equal(t, LaneValues(), want, "LaneValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := LaneValues(), LaneValues()
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

		tests := []struct {
			name string
			give Lane
			want string
		}{
			{name: "returns the declared text of Left", give: Left, want: "left"},
			{name: "returns the declared text of Right", give: Right, want: "right"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(
					t,
					tt.give.String(),
					tt.want,
					"String returns the declared text of the variant",
				)
			})
		}

		t.Run("returns a text of its own for each variant", func(t *testing.T) {
			t.Parallel()
			seen := map[string]bool{}
			for _, v := range []Lane{
				Left,
				Right,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := Lane("")
				for _, v := range []Lane{
					Left,
					Right,
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
			give Lane
		}{
			{name: "accepts Left", give: Left},
			{name: "accepts Right", give: Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, Lane("").IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Lane
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseLane", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Lane
		}{
			{name: "returns Left for its text", give: Left},
			{name: "returns Right for its text", give: Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseLane(tt.give.String())
				assert.NoError(t, err, "ParseLane accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseLane returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseLane("??????")
			assert.ErrorIs(t, err, ErrUnknownLane, "ParseLane refuses the text with ErrUnknownLane")
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Lane
		}{
			{name: "encodes Left as its text", give: Left},
			{name: "encodes Right as its text", give: Right},
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
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Lane
		}{
			{name: "decodes Left from its encoding", give: Left},
			{name: "decodes Right from its encoding", give: Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Lane
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Lane
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
