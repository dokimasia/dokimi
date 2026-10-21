package remote

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestSeasonVariants pins the declared set itself, which every other check
// here takes as given.
func TestSeasonVariants(t *testing.T) {
	t.Parallel()

	all := SeasonValues()

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Season]bool{}
		for _, v := range all {
			assert.False(t, seen[v], "each variant must hold a distinct value")
			seen[v] = true
		}
	})

	t.Run("no two variants share a textual form", func(t *testing.T) {
		t.Parallel()
		seen := map[string]bool{}
		for _, v := range all {
			assert.False(t, seen[v.String()], "each variant must render distinctly")
			seen[v.String()] = true
		}
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Season
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestSeasonText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestSeasonText(t *testing.T) {
	t.Parallel()

	all := SeasonValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseSeason(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseSeason("???????")
		assert.ErrorIs(t, err, ErrUnknownSeason, "unknown text must be refused with the declared refusal")
	})
}

// TestSeasonEncoding pins the boundary where a value most often crosses
// out of the program's control.
func TestSeasonEncoding(t *testing.T) {
	t.Parallel()

	all := SeasonValues()

	t.Run("every variant encodes to its textual form", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			text, err := v.MarshalText()
			assert.NoError(t, err, "a declared variant must encode")
			assert.Equal(t, string(text), v.String(), "the encoding must match the textual form")
		}
	})

	t.Run("every variant survives an encoding round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			text, err := v.MarshalText()
			assert.NoError(t, err, "a declared variant must encode")
			var got Season
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Season
		assert.HasError(t, got.UnmarshalText([]byte("???????")), "unknown text must not decode")
	})
}
