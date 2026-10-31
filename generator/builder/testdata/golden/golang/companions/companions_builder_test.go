package companions_test

import (
	"testing"

	"example.com/shop/companions"
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

			var seed companions.Taken
			seed.Name = "Name-a"
			got := companions.NewTakenBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewTakenBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := companions.NewTakenBuilder()
			want := b.Build()
			assert.Equal(
				t,
				companions.NewTakenBuilderFrom(want).Build(),
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
			got := companions.NewTakenBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := companions.NewTakenBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*companions.Taken
			record := func(v *companions.Taken) { got = append(got, v) }
			companions.NewTakenBuilder().Mutate(record).Mutate(record)
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

			original := companions.NewTakenBuilder()
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

			b := companions.NewTakenBuilder()
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

			var seed companions.Taken
			seed.Name = "Name-a"
			b := companions.NewTakenBuilder()
			b.Mutate(func(v *companions.Taken) { *v = seed })
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

			var seed companions.Paired
			seed.Name = "Name-a"
			got := companions.NewPairedBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewPairedBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := companions.NewPairedBuilder()
			want := b.Build()
			assert.Equal(
				t,
				companions.NewPairedBuilderFrom(want).Build(),
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
			got := companions.NewPairedBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := companions.NewPairedBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*companions.Paired
			record := func(v *companions.Paired) { got = append(got, v) }
			companions.NewPairedBuilder().Mutate(record).Mutate(record)
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

			original := companions.NewPairedBuilder()
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

			b := companions.NewPairedBuilder()
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

			var seed companions.Paired
			seed.Name = "Name-a"
			b := companions.NewPairedBuilder()
			b.Mutate(func(v *companions.Paired) { *v = seed })
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

			var seed companions.Bare
			seed.Name = "Name-a"
			got := companions.NewBareBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewBareBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := companions.NewBareBuilder()
			want := b.Build()
			assert.Equal(
				t,
				companions.NewBareBuilderFrom(want).Build(),
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
			got := companions.NewBareBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := companions.NewBareBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*companions.Bare
			record := func(v *companions.Bare) { got = append(got, v) }
			companions.NewBareBuilder().Mutate(record).Mutate(record)
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

			original := companions.NewBareBuilder()
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

			b := companions.NewBareBuilder()
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

			var seed companions.Bare
			seed.Name = "Name-a"
			b := companions.NewBareBuilder()
			b.Mutate(func(v *companions.Bare) { *v = seed })
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

			var seed companions.Other
			seed.Name = "Name-a"
			got := companions.NewOtherBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewOtherBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := companions.NewOtherBuilder()
			want := b.Build()
			assert.Equal(
				t,
				companions.NewOtherBuilderFrom(want).Build(),
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
			got := companions.NewOtherBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := companions.NewOtherBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*companions.Other
			record := func(v *companions.Other) { got = append(got, v) }
			companions.NewOtherBuilder().Mutate(record).Mutate(record)
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

			original := companions.NewOtherBuilder()
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

			b := companions.NewOtherBuilder()
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

			var seed companions.Other
			seed.Name = "Name-a"
			b := companions.NewOtherBuilder()
			b.Mutate(func(v *companions.Other) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}
