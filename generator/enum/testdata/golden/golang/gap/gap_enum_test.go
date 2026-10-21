package gap

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestGapVariants pins the declared set itself, which every other check
// here takes as given.
func TestGapVariants(t *testing.T) {
	t.Parallel()

	all := GapValues()

	t.Run("declares exactly 5 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 5, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Gap]bool{}
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
		assert.False(t, Gap(12).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Gap
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestGapText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestGapText(t *testing.T) {
	t.Parallel()

	all := GapValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseGap(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseGap("???????")
		assert.ErrorIs(t, err, ErrUnknownGap, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Gap(12)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestGapEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestGapEncoding(t *testing.T) {
	t.Parallel()

	all := GapValues()

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
			var got Gap
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Gap
		assert.HasError(t, got.UnmarshalText([]byte("???????")), "unknown text must not decode")
	})
}

// TestUnumVariants pins the declared set itself, which every other check
// here takes as given.
func TestUnumVariants(t *testing.T) {
	t.Parallel()

	all := UnumValues()

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Unum]bool{}
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
		assert.False(t, Unum(254).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is Naught", func(t *testing.T) {
		t.Parallel()
		var zero Unum
		assert.True(t, zero.IsValid(), "the zero value is a declared variant")
		assert.Equal(t, zero, Naught, "the zero value is Naught")
	})
}

// TestUnumText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestUnumText(t *testing.T) {
	t.Parallel()

	all := UnumValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseUnum(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseUnum("???????")
		assert.ErrorIs(t, err, ErrUnknownUnum, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Unum(254)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestUnumEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestUnumEncoding(t *testing.T) {
	t.Parallel()

	all := UnumValues()

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
			var got Unum
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Unum
		assert.HasError(t, got.UnmarshalText([]byte("???????")), "unknown text must not decode")
	})
}

// TestPrimeVariants pins the declared set itself, which every other check
// here takes as given.
func TestPrimeVariants(t *testing.T) {
	t.Parallel()

	all := PrimeValues()

	t.Run("declares exactly 12 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 12, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Prime]bool{}
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
		assert.False(t, Prime(38).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Prime
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestPrimeText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestPrimeText(t *testing.T) {
	t.Parallel()

	all := PrimeValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParsePrime(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParsePrime("????")
		assert.ErrorIs(t, err, ErrUnknownPrime, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Prime(38)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestPrimeEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestPrimeEncoding(t *testing.T) {
	t.Parallel()

	all := PrimeValues()

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
			var got Prime
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Prime
		assert.HasError(t, got.UnmarshalText([]byte("????")), "unknown text must not decode")
	})
}

// TestRankVariants pins the declared set itself, which every other check
// here takes as given.
func TestRankVariants(t *testing.T) {
	t.Parallel()

	all := RankValues()

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Rank]bool{}
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
		assert.False(t, Rank(4).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Rank
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestRankText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestRankText(t *testing.T) {
	t.Parallel()

	all := RankValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseRank(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseRank("???????")
		assert.ErrorIs(t, err, ErrUnknownRank, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Rank(4)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestRankEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestRankEncoding(t *testing.T) {
	t.Parallel()

	all := RankValues()

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
			var got Rank
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Rank
		assert.HasError(t, got.UnmarshalText([]byte("???????")), "unknown text must not decode")
	})
}

// TestOddVariants pins the declared set itself, which every other check
// here takes as given.
func TestOddVariants(t *testing.T) {
	t.Parallel()

	all := OddValues()

	t.Run("declares exactly 10 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 10, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Odd]bool{}
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
		assert.False(t, Odd(20).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is not a declared variant", func(t *testing.T) {
		t.Parallel()
		var zero Odd
		assert.False(t, zero.IsValid(), "the zero value must not be a variant")
	})
}

// TestOddText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestOddText(t *testing.T) {
	t.Parallel()

	all := OddValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseOdd(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseOdd("????")
		assert.ErrorIs(t, err, ErrUnknownOdd, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Odd(20)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestOddEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestOddEncoding(t *testing.T) {
	t.Parallel()

	all := OddValues()

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
			var got Odd
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Odd
		assert.HasError(t, got.UnmarshalText([]byte("????")), "unknown text must not decode")
	})
}

// TestNumVariants pins the declared set itself, which every other check
// here takes as given.
func TestNumVariants(t *testing.T) {
	t.Parallel()

	all := NumValues()

	t.Run("declares exactly 3 variants", func(t *testing.T) {
		t.Parallel()
		assert.Length(t, all, 3, "the declared set must not change unnoticed")
	})

	t.Run("no two variants share a value", func(t *testing.T) {
		t.Parallel()
		seen := map[Num]bool{}
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
		assert.False(t, Num(2).IsValid(), "an undeclared value must not be valid")
	})

	t.Run("every declared variant is valid", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			assert.True(t, v.IsValid(), "a declared variant must be valid")
		}
	})

	t.Run("the zero value is Zero", func(t *testing.T) {
		t.Parallel()
		var zero Num
		assert.True(t, zero.IsValid(), "the zero value is a declared variant")
		assert.Equal(t, zero, Zero, "the zero value is Zero")
	})
}

// TestNumText pins the round trip a value makes through anything that
// carries it as text — a log line, a JSON field, a database column.
func TestNumText(t *testing.T) {
	t.Parallel()

	all := NumValues()

	t.Run("every variant survives a text round trip", func(t *testing.T) {
		t.Parallel()
		for _, v := range all {
			got, err := ParseNum(v.String())
			assert.NoError(t, err, "a declared variant must parse back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused", func(t *testing.T) {
		t.Parallel()
		_, err := ParseNum("??????")
		assert.ErrorIs(t, err, ErrUnknownNum, "unknown text must be refused with the declared refusal")
	})

	t.Run("a value outside the set does not render as one inside it", func(t *testing.T) {
		t.Parallel()
		out := Num(2)
		for _, v := range all {
			assert.NotEqual(t, out.String(), v.String(), "an undeclared value must render distinctly")
		}
	})
}

// TestNumEncoding pins the boundary where a value most often crosses out
// of the program's control.
func TestNumEncoding(t *testing.T) {
	t.Parallel()

	all := NumValues()

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
			var got Num
			assert.NoError(t, got.UnmarshalText(text), "and decode back")
			assert.Equal(t, got, v, "the round trip must be lossless")
		}
	})

	t.Run("text naming no variant is refused on decode", func(t *testing.T) {
		t.Parallel()
		var got Num
		assert.HasError(t, got.UnmarshalText([]byte("??????")), "unknown text must not decode")
	})
}
