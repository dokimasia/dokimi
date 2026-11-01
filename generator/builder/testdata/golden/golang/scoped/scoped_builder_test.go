package scoped_test

import (
	maps2 "maps"
	slices2 "slices"
	testing2 "testing"

	"example.com/shop/scoped"
	assert2 "go.dokimi.dev/assert"
)

// TestBasketBuilder checks the constructors and the methods of
// BasketBuilder.
func TestBasketBuilder(t *testing2.T) {
	t.Parallel()

	t.Run("NewBasketBuilderFrom", func(t *testing2.T) {
		t.Parallel()

		t.Run("keeps the value of its argument", func(t *testing2.T) {
			t.Parallel()

			b := scoped.NewBasketBuilder()
			b.WithItems("")
			b.WithCounts(map[string]int{"": 0})
			want := b.Build()
			assert2.Equal(
				t,
				scoped.NewBasketBuilderFrom(want).Build(),
				want,
				"NewBasketBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithItems", func(t *testing2.T) {
		t.Parallel()

		t.Run("writes its argument to Items", func(t *testing2.T) {
			t.Parallel()

			var want []string
			want = []string{"Items-a"}
			got := scoped.NewBasketBuilder().WithItems(want...).Build()
			assert2.Equal(t, got.Items, want, "WithItems writes its argument to Items")
		})

		t.Run("replaces the value of Items", func(t *testing2.T) {
			t.Parallel()

			var first, want []string
			first = []string{"Items-a"}
			want = []string{"Items-b"}
			got := scoped.NewBasketBuilder().WithItems(first...).WithItems(want...).Build()
			assert2.Equal(t, got.Items, want, "the second call of WithItems replaces the first")
		})
	})

	t.Run("WithCounts", func(t *testing2.T) {
		t.Parallel()

		t.Run("writes its argument to Counts", func(t *testing2.T) {
			t.Parallel()

			var want map[string]int
			want = map[string]int{"Counts-a": 42}
			got := scoped.NewBasketBuilder().WithCounts(want).Build()
			assert2.Equal(t, got.Counts, want, "WithCounts writes its argument to Counts")
		})

		t.Run("replaces the value of Counts", func(t *testing2.T) {
			t.Parallel()

			var first, want map[string]int
			first = map[string]int{"Counts-a": 42}
			want = map[string]int{"Counts-b": 42}
			got := scoped.NewBasketBuilder().WithCounts(first).WithCounts(want).Build()
			assert2.Equal(t, got.Counts, want, "the second call of WithCounts replaces the first")
		})
	})

	t.Run("AppendItems", func(t *testing2.T) {
		t.Parallel()

		t.Run("appends to the elements of Items", func(t *testing2.T) {
			t.Parallel()

			var first, want []string
			first = []string{"Items-a"}
			want = []string{"Items-b"}
			got := scoped.NewBasketBuilder().WithItems(first...).AppendItems(want...).Build()
			assert2.Equal(
				t,
				got.Items,
				slices2.Concat(first, want),
				"AppendItems appends its elements after those of WithItems",
			)
		})
	})

	t.Run("WithCountsEntry", func(t *testing2.T) {
		t.Parallel()

		t.Run("adds one key to Counts", func(t *testing2.T) {
			t.Parallel()

			got := scoped.NewBasketBuilder().WithCounts(nil).WithCountsEntry("", 0).Build()
			assert2.Equal(
				t,
				got.Counts,
				map[string]int{"": 0},
				"WithCountsEntry adds the key to a nil map",
			)
		})
	})

	t.Run("WithCountsEntries", func(t *testing2.T) {
		t.Parallel()

		t.Run("adds each key to Counts", func(t *testing2.T) {
			t.Parallel()

			var first, want map[string]int
			first = map[string]int{"Counts-a": 42}
			want = map[string]int{"Counts-b": 42}
			merged := maps2.Clone(first)
			maps2.Copy(merged, want)
			got := scoped.NewBasketBuilder().
				WithCounts(nil).
				WithCountsEntries(first).
				WithCountsEntries(want).
				Build()
			assert2.Equal(
				t,
				got.Counts,
				merged,
				"WithCountsEntries keeps the entries of its earlier call",
			)
		})
	})

	t.Run("Mutate", func(t *testing2.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing2.T) {
			t.Parallel()

			var got []*scoped.Basket
			record := func(v *scoped.Basket) { got = append(got, v) }
			scoped.NewBasketBuilder().Mutate(record).Mutate(record)
			assert2.Length(t, got, 2, "Mutate calls each function once")
			assert2.Equal(
				t,
				got[0],
				got[1],
				"the two functions get the same value",
				assert2.ByIdentity(),
			)
		})
	})

	t.Run("Clone", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns a builder of its own", func(t *testing2.T) {
			t.Parallel()

			original := scoped.NewBasketBuilder()
			assert2.NotEqual(
				t,
				original.Clone(),
				original,
				"Clone returns a builder of its own",
				assert2.ByIdentity(),
			)
		})

		t.Run("builds the value of its original", func(t *testing2.T) {
			t.Parallel()

			b := scoped.NewBasketBuilder()
			b.WithItems("")
			b.WithCounts(map[string]int{"": 0})
			assert2.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})

		t.Run("copies Items into the clone", func(t *testing2.T) {
			t.Parallel()

			var want, other []string
			want = []string{"Items-a"}
			other = []string{"Items-b"}
			original := scoped.NewBasketBuilder().WithItems(want...)
			before := slices2.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *scoped.Basket) { copy(v.Items, other) })
			assert2.Equal(
				t,
				original.Build().Items,
				before,
				"a change of Items in the clone leaves the original unchanged",
			)
		})

		t.Run("copies Counts into the clone", func(t *testing2.T) {
			t.Parallel()

			var want, other map[string]int
			want = map[string]int{"Counts-a": 42}
			other = map[string]int{"Counts-b": 42}
			original := scoped.NewBasketBuilder().WithCounts(want)
			before := maps2.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *scoped.Basket) { maps2.Copy(v.Counts, other) })
			assert2.Equal(
				t,
				original.Build().Counts,
				before,
				"a change of Counts in the clone leaves the original unchanged",
			)
		})
	})
}
