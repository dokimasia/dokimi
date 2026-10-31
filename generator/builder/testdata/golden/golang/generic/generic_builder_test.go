package generic_test

import (
	"slices"
	"testing"

	"example.com/shop/generic"
	"go.dokimi.dev/assert"
)

// TestPageBuilder checks the constructors and the methods of PageBuilder.
func TestPageBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewPageBuilder", func(t *testing.T) {
		t.Parallel()

		t.Run("applies the declared default of Next", func(t *testing.T) {
			t.Parallel()

			var want int
			want = 1
			got := generic.NewPageBuilder[int]().Build()
			assert.Equal(t, got.Next, want, "NewPageBuilder applies the declared default of Next")
		})
	})

	t.Run("NewPageBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed generic.Page[int]
			seed.Next = 42
			got := generic.NewPageBuilderFrom[int](seed).Build()
			assert.Equal(t, got.Next, seed.Next, "NewPageBuilderFrom keeps Next")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := generic.NewPageBuilder[int]()
			b.WithItems(0)
			want := b.Build()
			assert.Equal(
				t,
				generic.NewPageBuilderFrom[int](want).Build(),
				want,
				"NewPageBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithItems", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Items", func(t *testing.T) {
			t.Parallel()

			var want []int
			want = []int{42}
			got := generic.NewPageBuilder[int]().WithItems(want...).Build()
			assert.Equal(t, got.Items, want, "WithItems writes its argument to Items")
		})

		t.Run("replaces the value of Items", func(t *testing.T) {
			t.Parallel()

			var first, want []int
			first = []int{42}
			want = []int{7}
			got := generic.NewPageBuilder[int]().WithItems(first...).WithItems(want...).Build()
			assert.Equal(t, got.Items, want, "the second call of WithItems replaces the first")
		})
	})

	t.Run("WithNext", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Next", func(t *testing.T) {
			t.Parallel()

			var want int
			want = 42
			got := generic.NewPageBuilder[int]().WithNext(want).Build()
			assert.Equal(t, got.Next, want, "WithNext writes its argument to Next")
		})

		t.Run("replaces the value of Next", func(t *testing.T) {
			t.Parallel()

			var first, want int
			first = 42
			want = 7
			got := generic.NewPageBuilder[int]().WithNext(first).WithNext(want).Build()
			assert.Equal(t, got.Next, want, "the second call of WithNext replaces the first")
		})
	})

	t.Run("AppendItems", func(t *testing.T) {
		t.Parallel()

		t.Run("appends to the elements of Items", func(t *testing.T) {
			t.Parallel()

			var first, want []int
			first = []int{42}
			want = []int{7}
			got := generic.NewPageBuilder[int]().WithItems(first...).AppendItems(want...).Build()
			assert.Equal(
				t,
				got.Items,
				slices.Concat(first, want),
				"AppendItems appends its elements after those of WithItems",
			)
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*generic.Page[int]
			record := func(v *generic.Page[int]) { got = append(got, v) }
			generic.NewPageBuilder[int]().Mutate(record).Mutate(record)
			assert.Length(t, got, 2, "Mutate calls each function once")
			assert.Equal(
				t,
				got[0],
				got[1],
				"the two functions get the same value",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("Clone", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a builder of its own", func(t *testing.T) {
			t.Parallel()

			original := generic.NewPageBuilder[int]()
			assert.NotEqual(
				t,
				original.Clone(),
				original,
				"Clone returns a builder of its own",
				assert.ByIdentity(),
			)
		})

		t.Run("builds the value of its original", func(t *testing.T) {
			t.Parallel()

			b := generic.NewPageBuilder[int]()
			b.WithItems(0)
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})

		t.Run("copies Items into the clone", func(t *testing.T) {
			t.Parallel()

			var want, other []int
			want = []int{42}
			other = []int{7}
			original := generic.NewPageBuilder[int]().WithItems(want...)
			before := slices.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *generic.Page[int]) { copy(v.Items, other) })
			assert.Equal(
				t,
				original.Build().Items,
				before,
				"a change of Items in the clone leaves the original unchanged",
			)
		})
	})

	t.Run("Build", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value under construction", func(t *testing.T) {
			t.Parallel()

			var seed generic.Page[int]
			seed.Next = 42
			b := generic.NewPageBuilder[int]()
			b.Mutate(func(v *generic.Page[int]) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}
