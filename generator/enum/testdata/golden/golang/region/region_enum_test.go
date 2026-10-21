package region

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestRegionVariants pins the declared set itself, which every other check
// here takes as given.
func TestRegionVariants(t *testing.T) {
	t.Parallel()

	all := RegionValues()

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Region]bool{}
		for _, v := range all {
			assert.False(t, seen[v], "each variant must hold a distinct value")
			seen[v] = true
		}
	})

	t.Run("every variant renders the value it declares", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, US.String(), "us-east", "the declared value is the textual form")
		assert.Equal(t, EU.String(), "eu-west", "the declared value is the textual form")
		assert.Equal(t, AP.String(), "ap-south", "the declared value is the textual form")
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
		assert.False(t, Region("").IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Region
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestRegionText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestRegionText(t *testing.T) {
	t.Parallel()

	all := RegionValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseRegion(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseRegion("?????????")
		assert.ErrorIs(t, err, ErrUnknownRegion, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Region("")
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestRegionEncoding pins the boundary where a value most often crosses
// out of the program's control.
func TestRegionEncoding(t *testing.T) {
	t.Parallel()

	all := RegionValues()

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
			var got Region
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Region
		assert.HasError(t, got.UnmarshalText([]byte("?????????")), "unknown text must not decode")
	})
}

// TestZoneVariants pins the declared set itself, which every other check
// here takes as given.
func TestZoneVariants(t *testing.T) {
	t.Parallel()

	all := ZoneValues()

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Zone]bool{}
		for _, v := range all {
			assert.False(t, seen[v], "each variant must hold a distinct value")
			seen[v] = true
		}
	})

	t.Run("every variant renders the value it declares", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, North.String(), "north", "the declared value is the textual form")
		assert.Equal(t, South.String(), "south", "the declared value is the textual form")
		assert.Equal(t, Unset.String(), "", "the declared value is the textual form")
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
		assert.False(t, Zone("??????").IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is Unset", func(t *testing.T) {
		t.Parallel()
		var zero Zone
		assert.True(t, zero.IsValid(), "the zero value is a declared variant")
		assert.Equal(t, zero, Unset, "the zero value is Unset")
	})
}

// TestZoneText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestZoneText(t *testing.T) {
	t.Parallel()

	all := ZoneValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseZone(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseZone("??????")
		assert.ErrorIs(t, err, ErrUnknownZone, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Zone("??????")
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestZoneEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestZoneEncoding(t *testing.T) {
	t.Parallel()

	all := ZoneValues()

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
			var got Zone
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Zone
		assert.HasError(t, got.UnmarshalText([]byte("??????")), "unknown text must not decode")
	})
}

// TestLaneVariants pins the declared set itself, which every other check
// here takes as given.
func TestLaneVariants(t *testing.T) {
	t.Parallel()

	all := LaneValues()

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Lane]bool{}
		for _, v := range all {
			assert.False(t, seen[v], "each variant must hold a distinct value")
			seen[v] = true
		}
	})

	t.Run("every variant renders the value it declares", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Left.String(), "left", "the declared value is the textual form")
		assert.Equal(t, Right.String(), "right", "the declared value is the textual form")
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
		assert.False(t, Lane("").IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Lane
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestLaneText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestLaneText(t *testing.T) {
	t.Parallel()

	all := LaneValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseLane(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseLane("??????")
		assert.ErrorIs(t, err, ErrUnknownLane, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Lane("")
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestLaneEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestLaneEncoding(t *testing.T) {
	t.Parallel()

	all := LaneValues()

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
			var got Lane
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Lane
		assert.HasError(t, got.UnmarshalText([]byte("??????")), "unknown text must not decode")
	})
}
