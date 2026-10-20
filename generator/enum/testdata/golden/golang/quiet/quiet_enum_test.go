package quiet

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestModeVariants pins the declared set itself, which every other check
// here takes as given.
func TestModeVariants(t *testing.T) {
	t.Parallel()

	all := []Mode{
		Off,
		On,
	}

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Mode]bool{}
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
		assert.False(t, Mode(2).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is Off", func(t *testing.T) {
		t.Parallel()
		var zero Mode
		assert.True(t, zero.IsValid(), "the zero value is a declared variant")
		assert.Equal(t, zero, Off, "the zero value is Off")
	})
}

// TestModeEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestModeEncoding(t *testing.T) {
	t.Parallel()

	all := []Mode{
		Off,
		On,
	}

	t.Run("every variant encodes to its textual form", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			text, err := v.MarshalText()
			assert.NoError(t, err, "a declared variant must encode")
			assert.Equal(t, string(text), v.String(), "the encoding must match the textual form")
		}
	})
}

// TestPlainVariants pins the declared set itself, which every other check
// here takes as given.
func TestPlainVariants(t *testing.T) {
	t.Parallel()

	all := []Plain{
		Loose,
		Tight,
	}

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Plain]bool{}
		for _, v := range all {
			assert.False(t, seen[v], "each variant must hold a distinct value")
			seen[v] = true
		}
	})
}

// TestToneVariants pins the declared set itself, which every other check
// here takes as given.
func TestToneVariants(t *testing.T) {
	t.Parallel()

	all := ToneValues()

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Tone]bool{}
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
		assert.False(t, Tone(3).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Tone
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestToneText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestToneText(t *testing.T) {
	t.Parallel()

	all := ToneValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseTone(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseTone("?????")
		assert.ErrorIs(t, err, ErrUnknownTone, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Tone(3)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestToneEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestToneEncoding(t *testing.T) {
	t.Parallel()

	all := ToneValues()

	t.Run("every variant encodes to its textual form", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			text, err := v.MarshalText()
			assert.NoError(t, err, "a declared variant must encode")
			assert.Equal(t, string(text), v.String(), "the encoding must match the textual form")
		}
	})
}

// TestSwitchVariants pins the declared set itself, which every other check
// here takes as given.
func TestSwitchVariants(t *testing.T) {
	t.Parallel()

	all := SwitchValues()

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Switch]bool{}
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

	t.Run("the zero value is Down", func(t *testing.T) {
		t.Parallel()
		var zero Switch
		assert.True(t, zero.IsValid(), "the zero value is a declared variant")
		assert.Equal(t, zero, Down, "the zero value is Down")
	})
}

// TestSwitchText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestSwitchText(t *testing.T) {
	t.Parallel()

	all := SwitchValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseSwitch(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseSwitch("?????")
		assert.ErrorIs(t, err, ErrUnknownSwitch, "unknown text must be refused with the declared refusal")
	})
}

// TestSwitchEncoding pins the boundary where a value most often crosses
// out of the program's control.
func TestSwitchEncoding(t *testing.T) {
	t.Parallel()

	all := SwitchValues()

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
			var got Switch
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Switch
		assert.HasError(t, got.UnmarshalText([]byte("?????")), "unknown text must not decode")
	})
}

// TestPhaseVariants pins the declared set itself, which every other check
// here takes as given.
func TestPhaseVariants(t *testing.T) {
	t.Parallel()

	all := []Phase{
		Warm,
		Cool,
	}

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Phase]bool{}
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
		assert.False(t, Phase(3).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Phase
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestSpeedVariants pins the declared set itself, which every other check
// here takes as given.
func TestSpeedVariants(t *testing.T) {
	t.Parallel()

	all := SpeedValues()

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Speed]bool{}
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
		assert.False(t, Speed(3).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Speed
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestSpeedText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestSpeedText(t *testing.T) {
	t.Parallel()

	all := SpeedValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseSpeed(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseSpeed("?????")
		assert.ErrorIs(t, err, ErrUnknownSpeed, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Speed(3)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestSpeedEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestSpeedEncoding(t *testing.T) {
	t.Parallel()

	all := SpeedValues()

	t.Run("every variant encodes to its textual form", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			text, err := v.MarshalText()
			assert.NoError(t, err, "a declared variant must encode")
			assert.Equal(t, string(text), v.String(), "the encoding must match the textual form")
		}
	})
}

// TestLevelVariants pins the declared set itself, which every other check
// here takes as given.
func TestLevelVariants(t *testing.T) {
	t.Parallel()

	all := LevelValues()

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Level]bool{}
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
		assert.False(t, Level(2.5).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Level
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestLevelText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestLevelText(t *testing.T) {
	t.Parallel()

	all := LevelValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseLevel(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseLevel("???????")
		assert.ErrorIs(t, err, ErrUnknownLevel, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Level(2.5)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestLevelEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestLevelEncoding(t *testing.T) {
	t.Parallel()

	all := LevelValues()

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
			var got Level
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Level
		assert.HasError(t, got.UnmarshalText([]byte("???????")), "unknown text must not decode")
	})
}

// TestFlagVariants pins the declared set itself, which every other check
// here takes as given.
func TestFlagVariants(t *testing.T) {
	t.Parallel()

	all := FlagValues()

	t.Run("declares exactly 1 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 1, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Flag]bool{}
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
		assert.False(t, Flag(false).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Flag
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestFlagText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestFlagText(t *testing.T) {
	t.Parallel()

	all := FlagValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseFlag(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseFlag("????")
		assert.ErrorIs(t, err, ErrUnknownFlag, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Flag(false)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestFlagEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestFlagEncoding(t *testing.T) {
	t.Parallel()

	all := FlagValues()

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
			var got Flag
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Flag
		assert.HasError(t, got.UnmarshalText([]byte("????")), "unknown text must not decode")
	})
}
