package region_test

import (
	"testing"

	"example.com/shop/region"
	"go.dokimi.dev/assert"
)

// TestRegion checks the functions and the methods of the enumeration
// Region.
func TestRegion(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[region.Region]bool{}
		for _, v := range []region.Region{
			region.US,
			region.EU,
			region.AP,
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
				want := []region.Region{
					region.US,
					region.EU,
					region.AP,
				}
				assert.Equal(
					t,
					region.RegionValues(),
					want,
					"RegionValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := region.RegionValues(), region.RegionValues()
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
			give region.Region
			want string
		}{
			{name: "returns the declared text of US", give: region.US, want: "us-east"},
			{name: "returns the declared text of EU", give: region.EU, want: "eu-west"},
			{name: "returns the declared text of AP", give: region.AP, want: "ap-south"},
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
			for _, v := range []region.Region{
				region.US,
				region.EU,
				region.AP,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := region.Region("")
				for _, v := range []region.Region{
					region.US,
					region.EU,
					region.AP,
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
			give region.Region
		}{
			{name: "accepts US", give: region.US},
			{name: "accepts EU", give: region.EU},
			{name: "accepts AP", give: region.AP},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, region.Region("").IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero region.Region
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseRegion", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give region.Region
		}{
			{name: "returns US for its text", give: region.US},
			{name: "returns EU for its text", give: region.EU},
			{name: "returns AP for its text", give: region.AP},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := region.ParseRegion(tt.give.String())
				assert.NoError(t, err, "ParseRegion accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseRegion returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := region.ParseRegion("?????????")
			assert.ErrorIs(
				t,
				err,
				region.ErrUnknownRegion,
				"ParseRegion refuses the text with ErrUnknownRegion",
			)
			assert.Contains(t, err.Error(), "?????????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give region.Region
		}{
			{name: "encodes US as its text", give: region.US},
			{name: "encodes EU as its text", give: region.EU},
			{name: "encodes AP as its text", give: region.AP},
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
			out := region.Region("")
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				region.ErrUnknownRegion,
				"MarshalText refuses the value with ErrUnknownRegion",
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
			give region.Region
		}{
			{name: "decodes US from its encoding", give: region.US},
			{name: "decodes EU from its encoding", give: region.EU},
			{name: "decodes AP from its encoding", give: region.AP},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got region.Region
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got region.Region
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
		seen := map[region.Zone]bool{}
		for _, v := range []region.Zone{
			region.North,
			region.South,
			region.Unset,
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
				want := []region.Zone{
					region.North,
					region.South,
					region.Unset,
				}
				assert.Equal(
					t,
					region.ZoneValues(),
					want,
					"ZoneValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := region.ZoneValues(), region.ZoneValues()
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
			give region.Zone
			want string
		}{
			{name: "returns the declared text of North", give: region.North, want: "north"},
			{name: "returns the declared text of South", give: region.South, want: "south"},
			{name: "returns the declared text of Unset", give: region.Unset, want: ""},
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
			for _, v := range []region.Zone{
				region.North,
				region.South,
				region.Unset,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := region.Zone("??????")
				for _, v := range []region.Zone{
					region.North,
					region.South,
					region.Unset,
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
			give region.Zone
		}{
			{name: "accepts North", give: region.North},
			{name: "accepts South", give: region.South},
			{name: "accepts Unset", give: region.Unset},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, region.Zone("??????").IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("accepts the zero value Unset", func(t *testing.T) {
			t.Parallel()
			var zero region.Zone
			assert.Equal(t, zero, region.Unset, "the zero value is Unset")
			assert.True(t, zero.IsValid(), "IsValid accepts the zero value")
		})
	})

	t.Run("ParseZone", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give region.Zone
		}{
			{name: "returns North for its text", give: region.North},
			{name: "returns South for its text", give: region.South},
			{name: "returns Unset for its text", give: region.Unset},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := region.ParseZone(tt.give.String())
				assert.NoError(t, err, "ParseZone accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseZone returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := region.ParseZone("??????")
			assert.ErrorIs(
				t,
				err,
				region.ErrUnknownZone,
				"ParseZone refuses the text with ErrUnknownZone",
			)
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give region.Zone
		}{
			{name: "encodes North as its text", give: region.North},
			{name: "encodes South as its text", give: region.South},
			{name: "encodes Unset as its text", give: region.Unset},
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
			out := region.Zone("??????")
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				region.ErrUnknownZone,
				"MarshalText refuses the value with ErrUnknownZone",
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
			give region.Zone
		}{
			{name: "decodes North from its encoding", give: region.North},
			{name: "decodes South from its encoding", give: region.South},
			{name: "decodes Unset from its encoding", give: region.Unset},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got region.Zone
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got region.Zone
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
		seen := map[region.Lane]bool{}
		for _, v := range []region.Lane{
			region.Left,
			region.Right,
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
				want := []region.Lane{
					region.Left,
					region.Right,
				}
				assert.Equal(
					t,
					region.LaneValues(),
					want,
					"LaneValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := region.LaneValues(), region.LaneValues()
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
			give region.Lane
			want string
		}{
			{name: "returns the declared text of Left", give: region.Left, want: "left"},
			{name: "returns the declared text of Right", give: region.Right, want: "right"},
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
			for _, v := range []region.Lane{
				region.Left,
				region.Right,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})

		t.Run(
			"renders a value outside the declared set with a text of its own",
			func(t *testing.T) {
				t.Parallel()
				out := region.Lane("")
				for _, v := range []region.Lane{
					region.Left,
					region.Right,
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
			give region.Lane
		}{
			{name: "accepts Left", give: region.Left},
			{name: "accepts Right", give: region.Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses a value outside the declared set", func(t *testing.T) {
			t.Parallel()
			assert.False(t, region.Lane("").IsValid(), "IsValid refuses an undeclared value")
		})

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero region.Lane
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseLane", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give region.Lane
		}{
			{name: "returns Left for its text", give: region.Left},
			{name: "returns Right for its text", give: region.Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := region.ParseLane(tt.give.String())
				assert.NoError(t, err, "ParseLane accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseLane returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := region.ParseLane("??????")
			assert.ErrorIs(
				t,
				err,
				region.ErrUnknownLane,
				"ParseLane refuses the text with ErrUnknownLane",
			)
			assert.Contains(t, err.Error(), "??????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give region.Lane
		}{
			{name: "encodes Left as its text", give: region.Left},
			{name: "encodes Right as its text", give: region.Right},
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
			out := region.Lane("")
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				region.ErrUnknownLane,
				"MarshalText refuses the value with ErrUnknownLane",
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
			give region.Lane
		}{
			{name: "decodes Left from its encoding", give: region.Left},
			{name: "decodes Right from its encoding", give: region.Right},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got region.Lane
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got region.Lane
			assert.HasError(
				t,
				got.UnmarshalText([]byte("??????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
