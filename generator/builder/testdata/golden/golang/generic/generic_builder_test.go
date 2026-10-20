package generic

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestPageBuilderConstruction pins what the constructors produce.
func TestPageBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Page[int]
		seed.Next = 42
		got := NewPageBuilderFrom[int](seed).Build()
		assert.Equal(t, got.Next, seed.Next, "From keeps the member it was given")
	})

	t.Run("Next carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want int
		want = 1
		got := NewPageBuilder[int]().Build()
		assert.Equal(t, got.Next, want, "the constructor applies the declared default")
	})
}

// TestPageBuilderSetters pins that each setter writes the member it names,
// and writes nothing else.
func TestPageBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithItems writes Items", func(t *testing.T) {
		t.Parallel()

		var want []int
		want = []int{42}
		got := NewPageBuilder[int]().WithItems(want...).Build()
		assert.Length(t, got.Items, 1, "the setter replaces the sequence")
	})

	t.Run("WithNext writes Next", func(t *testing.T) {
		t.Parallel()

		var want int
		want = 42
		got := NewPageBuilder[int]().WithNext(want).Build()
		assert.Equal(t, got.Next, want, "the setter writes what it was given")
	})

	t.Run("WithNext replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want int
		first = 42
		want = 7
		got := NewPageBuilder[int]().WithNext(first).WithNext(want).Build()
		assert.Equal(t, got.Next, want, "the later call wins")
	})

	t.Run("AppendItems keeps what Items holds", func(t *testing.T) {
		t.Parallel()

		got := NewPageBuilder[int]().WithItems(*new(int)).AppendItems(*new(int)).Build()
		assert.Length(t, got.Items, 2, "the appender adds to the sequence")
	})
}

// TestPageBuilderClone pins that a clone can be configured without the
// original seeing it.
func TestPageBuilderClone(t *testing.T) {
	t.Parallel()

	t.Run("Items is not shared with a clone", func(t *testing.T) {
		t.Parallel()

		var want, other []int
		want = []int{42}
		other = []int{7}
		original := NewPageBuilder[int]().WithItems(want...)
		clone := original.Clone()
		clone.Mutate(func(v *Page[int]) { copy(v.Items, other) })
		assert.NotEqual(t, original.Build().Items, other, "writing through a clone must not reach the original")
	})
}

// TestPageBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestPageBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Page[int]
		record := func(v *Page[int]) { got = append(got, v) }
		NewPageBuilder[int]().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewPageBuilder[int]()
		original.WithItems(*new(int))
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewPageBuilderFrom[int](want).Build(), want, "From keeps the value it was given")
	})
}
