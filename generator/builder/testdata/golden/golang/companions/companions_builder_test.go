package companions

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestTakenBuilder checks the constructors and the methods of
// TakenBuilder.
func TestTakenBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewTakenBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed Taken
			seed.Name = "Name-a"
			got := NewTakenBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewTakenBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewTakenBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewTakenBuilderFrom(want).Build(),
				want,
				"NewTakenBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithName", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Name", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Name-a"
			got := NewTakenBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := NewTakenBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Taken
			record := func(v *Taken) { got = append(got, v) }
			NewTakenBuilder().Mutate(record).Mutate(record)
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

			original := NewTakenBuilder()
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

			b := NewTakenBuilder()
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})
	})

	t.Run("Build", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value under construction", func(t *testing.T) {
			t.Parallel()

			var seed Taken
			seed.Name = "Name-a"
			b := NewTakenBuilder()
			b.Mutate(func(v *Taken) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestPairedBuilder checks the constructors and the methods of
// PairedBuilder.
func TestPairedBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewPairedBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed Paired
			seed.Name = "Name-a"
			got := NewPairedBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewPairedBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewPairedBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewPairedBuilderFrom(want).Build(),
				want,
				"NewPairedBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithName", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Name", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Name-a"
			got := NewPairedBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := NewPairedBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Paired
			record := func(v *Paired) { got = append(got, v) }
			NewPairedBuilder().Mutate(record).Mutate(record)
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

			original := NewPairedBuilder()
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

			b := NewPairedBuilder()
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})
	})

	t.Run("Build", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value under construction", func(t *testing.T) {
			t.Parallel()

			var seed Paired
			seed.Name = "Name-a"
			b := NewPairedBuilder()
			b.Mutate(func(v *Paired) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestBareBuilder checks the constructors and the methods of BareBuilder.
func TestBareBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewBareBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed Bare
			seed.Name = "Name-a"
			got := NewBareBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewBareBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewBareBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewBareBuilderFrom(want).Build(),
				want,
				"NewBareBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithName", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Name", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Name-a"
			got := NewBareBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := NewBareBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Bare
			record := func(v *Bare) { got = append(got, v) }
			NewBareBuilder().Mutate(record).Mutate(record)
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

			original := NewBareBuilder()
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

			b := NewBareBuilder()
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})
	})

	t.Run("Build", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value under construction", func(t *testing.T) {
			t.Parallel()

			var seed Bare
			seed.Name = "Name-a"
			b := NewBareBuilder()
			b.Mutate(func(v *Bare) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestOtherBuilder checks the constructors and the methods of
// OtherBuilder.
func TestOtherBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewOtherBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed Other
			seed.Name = "Name-a"
			got := NewOtherBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewOtherBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewOtherBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewOtherBuilderFrom(want).Build(),
				want,
				"NewOtherBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithName", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Name", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Name-a"
			got := NewOtherBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := NewOtherBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Other
			record := func(v *Other) { got = append(got, v) }
			NewOtherBuilder().Mutate(record).Mutate(record)
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

			original := NewOtherBuilder()
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

			b := NewOtherBuilder()
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})
	})

	t.Run("Build", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value under construction", func(t *testing.T) {
			t.Parallel()

			var seed Other
			seed.Name = "Name-a"
			b := NewOtherBuilder()
			b.Mutate(func(v *Other) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}
