package handwritten

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestWeekdayVariants pins the declared set itself, which every other
// check here takes as given.
func TestWeekdayVariants(t *testing.T) {
	t.Parallel()

	all := []Weekday{
		Monday,
		Tuesday,
		Wednesday,
	}

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Weekday]bool{}
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

	t.Run("a value outside the set is not valid", func(t *testing.T) {
		t.Parallel()
		assert.False(t, Weekday(4).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Weekday
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestWeekdayEncoding pins the boundary where a value most often crosses
// out of the program's control.
func TestWeekdayEncoding(t *testing.T) {
	t.Parallel()

	all := []Weekday{
		Monday,
		Tuesday,
		Wednesday,
	}

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
			var got Weekday
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Weekday
		assert.HasError(t, got.UnmarshalText([]byte("??????????")), "unknown text must not decode")
	})
}
