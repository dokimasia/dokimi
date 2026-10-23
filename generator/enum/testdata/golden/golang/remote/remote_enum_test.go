package remote

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestSeason checks the functions and the methods of the enumeration
// Season.
func TestSeason(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[Season]bool{}
		for _, v := range []Season{
			Spring,
			Autumn,
		} {
			assert.False(t, seen[v], "each variant has a value of its own")
			seen[v] = true
		}
	})

	t.Run("SeasonValues", func(t *testing.T) {
		t.Parallel()

		t.Run(
			"returns the declared variants in the order of their declaration",
			func(t *testing.T) {
				t.Parallel()
				want := []Season{
					Spring,
					Autumn,
				}
				assert.Equal(t, SeasonValues(), want, "SeasonValues returns the declared variants")
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := SeasonValues(), SeasonValues()
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
			for _, v := range []Season{
				Spring,
				Autumn,
			} {
				assert.False(t, seen[v.String()], "each variant has a text of its own")
				seen[v.String()] = true
			}
		})
	})

	t.Run("IsValid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Season
		}{
			{name: "accepts Spring", give: Spring},
			{name: "accepts Autumn", give: Autumn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero Season
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseSeason", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Season
		}{
			{name: "returns Spring for its text", give: Spring},
			{name: "returns Autumn for its text", give: Autumn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := ParseSeason(tt.give.String())
				assert.NoError(t, err, "ParseSeason accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseSeason returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := ParseSeason("???????")
			assert.ErrorIs(
				t,
				err,
				ErrUnknownSeason,
				"ParseSeason refuses the text with ErrUnknownSeason",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give Season
		}{
			{name: "encodes Spring as its text", give: Spring},
			{name: "encodes Autumn as its text", give: Autumn},
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
			var out Season
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				ErrUnknownSeason,
				"MarshalText refuses the value with ErrUnknownSeason",
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
			give Season
		}{
			{name: "decodes Spring from its encoding", give: Spring},
			{name: "decodes Autumn from its encoding", give: Autumn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got Season
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got Season
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
