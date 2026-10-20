package commented

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestGradeVariants pins the declared set itself, which every other check
// here takes as given.
func TestGradeVariants(t *testing.T) {
	t.Parallel()

	all := GradeValues()

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Grade]bool{}
		for _, v := range all {
			assert.False(t, seen[v], "each variant must hold a distinct value")
			seen[v] = true
		}
	})

	t.Run("every variant renders the value it declares", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, Poor.String(), "poor", "the declared value is the textual form")
		assert.Equal(t, Fair.String(), "fair enough", "the declared value is the textual form")
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
		assert.False(t, Grade(3).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is Poor", func(t *testing.T) {
		t.Parallel()
		var zero Grade
		assert.True(t, zero.IsValid(), "the zero value is a declared variant")
		assert.Equal(t, zero, Poor, "the zero value is Poor")
	})
}

// TestGradeText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestGradeText(t *testing.T) {
	t.Parallel()

	all := GradeValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseGrade(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseGrade("????????????")
		assert.ErrorIs(t, err, ErrUnknownGrade, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Grade(3)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestGradeEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestGradeEncoding(t *testing.T) {
	t.Parallel()

	all := GradeValues()

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
			var got Grade
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Grade
		assert.HasError(t, got.UnmarshalText([]byte("????????????")), "unknown text must not decode")
	})
}

// TestMoodVariants pins the declared set itself, which every other check
// here takes as given.
func TestMoodVariants(t *testing.T) {
	t.Parallel()

	all := MoodValues()

	t.Run("declares exactly 2 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 2, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Mood]bool{}
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
		assert.False(t, Mood(2).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is Calm", func(t *testing.T) {
		t.Parallel()
		var zero Mood
		assert.True(t, zero.IsValid(), "the zero value is a declared variant")
		assert.Equal(t, zero, Calm, "the zero value is Calm")
	})
}

// TestMoodText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestMoodText(t *testing.T) {
	t.Parallel()

	all := MoodValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseMood(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseMood("??????")
		assert.ErrorIs(t, err, ErrUnknownMood, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Mood(2)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestMoodEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestMoodEncoding(t *testing.T) {
	t.Parallel()

	all := MoodValues()

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
			var got Mood
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Mood
		assert.HasError(t, got.UnmarshalText([]byte("??????")), "unknown text must not decode")
	})
}
