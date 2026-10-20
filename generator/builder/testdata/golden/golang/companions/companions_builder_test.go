package companions

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestTakenBuilderConstruction pins what the constructors produce.
func TestTakenBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Taken
		seed.Name = "Name-a"
		got := NewTakenBuilderFrom(seed).Build()
		assert.Equal(t, got.Name, seed.Name, "From keeps the member it was given")
	})
}

// TestTakenBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestTakenBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithName writes Name", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Name-a"
		got := NewTakenBuilder().WithName(want).Build()
		assert.Equal(t, got.Name, want, "the setter writes what it was given")
	})

	t.Run("WithName replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Name-a"
		want = "Name-b"
		got := NewTakenBuilder().WithName(first).WithName(want).Build()
		assert.Equal(t, got.Name, want, "the later call wins")
	})
}

// TestTakenBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestTakenBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Taken
		record := func(v *Taken) { got = append(got, v) }
		NewTakenBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewTakenBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewTakenBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestPairedBuilderConstruction pins what the constructors produce.
func TestPairedBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Paired
		seed.Name = "Name-a"
		got := NewPairedBuilderFrom(seed).Build()
		assert.Equal(t, got.Name, seed.Name, "From keeps the member it was given")
	})
}

// TestPairedBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestPairedBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithName writes Name", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Name-a"
		got := NewPairedBuilder().WithName(want).Build()
		assert.Equal(t, got.Name, want, "the setter writes what it was given")
	})

	t.Run("WithName replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Name-a"
		want = "Name-b"
		got := NewPairedBuilder().WithName(first).WithName(want).Build()
		assert.Equal(t, got.Name, want, "the later call wins")
	})
}

// TestPairedBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestPairedBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Paired
		record := func(v *Paired) { got = append(got, v) }
		NewPairedBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewPairedBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewPairedBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestBareBuilderConstruction pins what the constructors produce.
func TestBareBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Bare
		seed.Name = "Name-a"
		got := NewBareBuilderFrom(seed).Build()
		assert.Equal(t, got.Name, seed.Name, "From keeps the member it was given")
	})
}

// TestBareBuilderSetters pins that each setter writes the member it names,
// and writes nothing else.
func TestBareBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithName writes Name", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Name-a"
		got := NewBareBuilder().WithName(want).Build()
		assert.Equal(t, got.Name, want, "the setter writes what it was given")
	})

	t.Run("WithName replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Name-a"
		want = "Name-b"
		got := NewBareBuilder().WithName(first).WithName(want).Build()
		assert.Equal(t, got.Name, want, "the later call wins")
	})
}

// TestBareBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestBareBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Bare
		record := func(v *Bare) { got = append(got, v) }
		NewBareBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewBareBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewBareBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestOtherBuilderConstruction pins what the constructors produce.
func TestOtherBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Other
		seed.Name = "Name-a"
		got := NewOtherBuilderFrom(seed).Build()
		assert.Equal(t, got.Name, seed.Name, "From keeps the member it was given")
	})
}

// TestOtherBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestOtherBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithName writes Name", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Name-a"
		got := NewOtherBuilder().WithName(want).Build()
		assert.Equal(t, got.Name, want, "the setter writes what it was given")
	})

	t.Run("WithName replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Name-a"
		want = "Name-b"
		got := NewOtherBuilder().WithName(first).WithName(want).Build()
		assert.Equal(t, got.Name, want, "the later call wins")
	})
}

// TestOtherBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestOtherBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Other
		record := func(v *Other) { got = append(got, v) }
		NewOtherBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewOtherBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewOtherBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}
