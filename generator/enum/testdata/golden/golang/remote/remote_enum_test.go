package remote_test

import (
	"testing"

	"example.com/shop/remote"
	"go.dokimi.dev/assert"
)

// TestSeason checks the functions and the methods of the enumeration
// Season.
func TestSeason(t *testing.T) {
	t.Parallel()

	t.Run("gives each variant a value of its own", func(t *testing.T) {
		t.Parallel()
		seen := map[remote.Season]bool{}
		for _, v := range []remote.Season{
			remote.Spring,
			remote.Autumn,
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
				want := []remote.Season{
					remote.Spring,
					remote.Autumn,
				}
				assert.Equal(
					t,
					remote.SeasonValues(),
					want,
					"SeasonValues returns the declared variants",
				)
			},
		)

		t.Run("returns a new slice for each call", func(t *testing.T) {
			t.Parallel()
			first, second := remote.SeasonValues(), remote.SeasonValues()
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
			for _, v := range []remote.Season{
				remote.Spring,
				remote.Autumn,
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
			give remote.Season
		}{
			{name: "accepts Spring", give: remote.Spring},
			{name: "accepts Autumn", give: remote.Autumn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				assert.True(t, tt.give.IsValid(), "IsValid accepts each declared variant")
			})
		}

		t.Run("refuses the zero value", func(t *testing.T) {
			t.Parallel()
			var zero remote.Season
			assert.False(t, zero.IsValid(), "IsValid refuses the zero value")
		})
	})

	t.Run("ParseSeason", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give remote.Season
		}{
			{name: "returns Spring for its text", give: remote.Spring},
			{name: "returns Autumn for its text", give: remote.Autumn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := remote.ParseSeason(tt.give.String())
				assert.NoError(t, err, "ParseSeason accepts the text of the variant")
				assert.Equal(t, got, tt.give, "ParseSeason returns the variant of the text")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			_, err := remote.ParseSeason("???????")
			assert.ErrorIs(
				t,
				err,
				remote.ErrUnknownSeason,
				"ParseSeason refuses the text with ErrUnknownSeason",
			)
			assert.Contains(t, err.Error(), "???????", "the error contains the refused text")
		})
	})

	t.Run("MarshalText", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give remote.Season
		}{
			{name: "encodes Spring as its text", give: remote.Spring},
			{name: "encodes Autumn as its text", give: remote.Autumn},
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
			var out remote.Season
			_, err := out.MarshalText()
			assert.ErrorIs(
				t,
				err,
				remote.ErrUnknownSeason,
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
			give remote.Season
		}{
			{name: "decodes Spring from its encoding", give: remote.Spring},
			{name: "decodes Autumn from its encoding", give: remote.Autumn},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				text, err := tt.give.MarshalText()
				assert.NoError(t, err, "MarshalText encodes the variant")
				var got remote.Season
				assert.NoError(t, got.UnmarshalText(text), "UnmarshalText decodes the encoding")
				assert.Equal(t, got, tt.give, "UnmarshalText returns the encoded variant")
			})
		}

		t.Run("refuses a text that no variant has", func(t *testing.T) {
			t.Parallel()
			var got remote.Season
			assert.HasError(
				t,
				got.UnmarshalText([]byte("???????")),
				"UnmarshalText refuses the text",
			)
		})
	})
}
