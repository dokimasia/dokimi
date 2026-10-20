package split

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestSignalVariants pins the declared set itself, which every other check
// here takes as given.
func TestSignalVariants(t *testing.T) {
	t.Parallel()

	all := SignalValues()

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Signal]bool{}
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
		assert.False(t, Signal(4).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Signal
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestSignalText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestSignalText(t *testing.T) {
	t.Parallel()

	all := SignalValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseSignal(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseSignal("??????")
		assert.ErrorIs(t, err, ErrUnknownSignal, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Signal(4)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestSignalEncoding pins the boundary where a value most often crosses
// out of the program's control.
func TestSignalEncoding(t *testing.T) {
	t.Parallel()

	all := SignalValues()

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
			var got Signal
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Signal
		assert.HasError(t, got.UnmarshalText([]byte("??????")), "unknown text must not decode")
	})
}
