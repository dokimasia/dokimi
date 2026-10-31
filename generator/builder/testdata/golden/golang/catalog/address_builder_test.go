package catalog_test

import (
	"maps"
	"slices"
	"testing"
	"time"

	"example.com/shop/catalog"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
)

// TestAddressBuilder checks the constructors and the methods of
// AddressBuilder.
func TestAddressBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewAddressBuilder", func(t *testing.T) {
		t.Parallel()

		t.Run("applies the declared default of Street", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Dam 1"
			got := catalog.NewAddressBuilder().Build()
			assert.Equal(
				t,
				got.Street,
				want,
				"NewAddressBuilder applies the declared default of Street",
			)
		})

		t.Run("applies the declared default of City", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Amsterdam"
			got := catalog.NewAddressBuilder().Build()
			assert.Equal(
				t,
				got.City,
				want,
				"NewAddressBuilder applies the declared default of City",
			)
		})

		t.Run("applies the declared default of Timeout", func(t *testing.T) {
			t.Parallel()

			var want time.Duration
			want = time.Second
			got := catalog.NewAddressBuilder().Build()
			assert.Equal(
				t,
				got.Timeout,
				want,
				"NewAddressBuilder applies the declared default of Timeout",
			)
		})
	})

	t.Run("NewAddressBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed catalog.Address
			seed.Street = "Street-a"
			seed.City = "City-a"
			seed.Grid = [2]int{42}
			seed.Timeout = time.Duration(1)
			got := catalog.NewAddressBuilderFrom(seed).Build()
			expect.Equal(t, got.Street, seed.Street, "NewAddressBuilderFrom keeps Street")
			expect.Equal(t, got.City, seed.City, "NewAddressBuilderFrom keeps City")
			expect.Equal(t, got.Grid, seed.Grid, "NewAddressBuilderFrom keeps Grid")
			expect.Equal(t, got.Timeout, seed.Timeout, "NewAddressBuilderFrom keeps Timeout")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewAddressBuilder()
			b.WithLines("")
			b.WithLabels(map[string]string{"": ""})
			b.WithTags(map[string]struct{}{"": {}})
			b.WithRaw([]byte{0})
			want := b.Build()
			assert.Equal(
				t,
				catalog.NewAddressBuilderFrom(want).Build(),
				want,
				"NewAddressBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithStreet", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Street", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Street-a"
			got := catalog.NewAddressBuilder().WithStreet(want).Build()
			assert.Equal(t, got.Street, want, "WithStreet writes its argument to Street")
		})

		t.Run("replaces the value of Street", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Street-a"
			want = "Street-b"
			got := catalog.NewAddressBuilder().WithStreet(first).WithStreet(want).Build()
			assert.Equal(t, got.Street, want, "the second call of WithStreet replaces the first")
		})
	})

	t.Run("WithCity", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to City", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "City-a"
			got := catalog.NewAddressBuilder().WithCity(want).Build()
			assert.Equal(t, got.City, want, "WithCity writes its argument to City")
		})

		t.Run("replaces the value of City", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "City-a"
			want = "City-b"
			got := catalog.NewAddressBuilder().WithCity(first).WithCity(want).Build()
			assert.Equal(t, got.City, want, "the second call of WithCity replaces the first")
		})
	})

	t.Run("WithFloor", func(t *testing.T) {
		t.Parallel()

		t.Run("points Floor at its argument", func(t *testing.T) {
			t.Parallel()

			var want int
			want = 42
			got := catalog.NewAddressBuilder().WithFloor(want).Build()
			assert.NotNil(t, got.Floor, "WithFloor points Floor at a value")
			assert.Equal(t, *got.Floor, want, "WithFloor points Floor at its argument")
		})

		t.Run("replaces the value of Floor", func(t *testing.T) {
			t.Parallel()

			var first, want int
			first = 42
			want = 7
			got := catalog.NewAddressBuilder().WithFloor(first).WithFloor(want).Build()
			assert.NotNil(t, got.Floor, "WithFloor points Floor at a value")
			assert.Equal(t, *got.Floor, want, "the second call of WithFloor replaces the first")
		})
	})

	t.Run("WithLines", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Lines", func(t *testing.T) {
			t.Parallel()

			var want []string
			want = []string{"Lines-a"}
			got := catalog.NewAddressBuilder().WithLines(want...).Build()
			assert.Equal(t, got.Lines, want, "WithLines writes its argument to Lines")
		})

		t.Run("replaces the value of Lines", func(t *testing.T) {
			t.Parallel()

			var first, want []string
			first = []string{"Lines-a"}
			want = []string{"Lines-b"}
			got := catalog.NewAddressBuilder().WithLines(first...).WithLines(want...).Build()
			assert.Equal(t, got.Lines, want, "the second call of WithLines replaces the first")
		})
	})

	t.Run("WithLabels", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Labels", func(t *testing.T) {
			t.Parallel()

			var want map[string]string
			want = map[string]string{"Labels-a": "Labels-a"}
			got := catalog.NewAddressBuilder().WithLabels(want).Build()
			assert.Equal(t, got.Labels, want, "WithLabels writes its argument to Labels")
		})

		t.Run("replaces the value of Labels", func(t *testing.T) {
			t.Parallel()

			var first, want map[string]string
			first = map[string]string{"Labels-a": "Labels-a"}
			want = map[string]string{"Labels-b": "Labels-a"}
			got := catalog.NewAddressBuilder().WithLabels(first).WithLabels(want).Build()
			assert.Equal(t, got.Labels, want, "the second call of WithLabels replaces the first")
		})
	})

	t.Run("WithTags", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Tags", func(t *testing.T) {
			t.Parallel()

			var want map[string]struct{}
			want = map[string]struct{}{"Tags-a": {}}
			got := catalog.NewAddressBuilder().WithTags(want).Build()
			assert.Equal(t, got.Tags, want, "WithTags writes its argument to Tags")
		})

		t.Run("replaces the value of Tags", func(t *testing.T) {
			t.Parallel()

			var first, want map[string]struct{}
			first = map[string]struct{}{"Tags-a": {}}
			want = map[string]struct{}{"Tags-b": {}}
			got := catalog.NewAddressBuilder().WithTags(first).WithTags(want).Build()
			assert.Equal(t, got.Tags, want, "the second call of WithTags replaces the first")
		})
	})

	t.Run("WithRaw", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Raw", func(t *testing.T) {
			t.Parallel()

			var want []byte
			want = []byte{1}
			got := catalog.NewAddressBuilder().WithRaw(want).Build()
			assert.Equal(t, got.Raw, want, "WithRaw writes its argument to Raw")
		})

		t.Run("replaces the value of Raw", func(t *testing.T) {
			t.Parallel()

			var first, want []byte
			first = []byte{1}
			want = []byte{2}
			got := catalog.NewAddressBuilder().WithRaw(first).WithRaw(want).Build()
			assert.Equal(t, got.Raw, want, "the second call of WithRaw replaces the first")
		})
	})

	t.Run("WithGrid", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Grid", func(t *testing.T) {
			t.Parallel()

			var want [2]int
			want = [2]int{42}
			got := catalog.NewAddressBuilder().WithGrid(want).Build()
			assert.Equal(t, got.Grid, want, "WithGrid writes its argument to Grid")
		})

		t.Run("replaces the value of Grid", func(t *testing.T) {
			t.Parallel()

			var first, want [2]int
			first = [2]int{42}
			want = [2]int{7}
			got := catalog.NewAddressBuilder().WithGrid(first).WithGrid(want).Build()
			assert.Equal(t, got.Grid, want, "the second call of WithGrid replaces the first")
		})
	})

	t.Run("WithTimeout", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Timeout", func(t *testing.T) {
			t.Parallel()

			var want time.Duration
			want = time.Duration(1)
			got := catalog.NewAddressBuilder().WithTimeout(want).Build()
			assert.Equal(t, got.Timeout, want, "WithTimeout writes its argument to Timeout")
		})

		t.Run("replaces the value of Timeout", func(t *testing.T) {
			t.Parallel()

			var first, want time.Duration
			first = time.Duration(1)
			want = time.Duration(2)
			got := catalog.NewAddressBuilder().WithTimeout(first).WithTimeout(want).Build()
			assert.Equal(t, got.Timeout, want, "the second call of WithTimeout replaces the first")
		})
	})

	t.Run("AppendLines", func(t *testing.T) {
		t.Parallel()

		t.Run("appends to the elements of Lines", func(t *testing.T) {
			t.Parallel()

			var first, want []string
			first = []string{"Lines-a"}
			want = []string{"Lines-b"}
			got := catalog.NewAddressBuilder().WithLines(first...).AppendLines(want...).Build()
			assert.Equal(
				t,
				got.Lines,
				slices.Concat(first, want),
				"AppendLines appends its elements after those of WithLines",
			)
		})
	})

	t.Run("WithLabelsEntry", func(t *testing.T) {
		t.Parallel()

		t.Run("adds one key to Labels", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewAddressBuilder().WithLabels(nil).WithLabelsEntry("", "").Build()
			assert.Equal(
				t,
				got.Labels,
				map[string]string{"": ""},
				"WithLabelsEntry adds the key to a nil map",
			)
		})
	})

	t.Run("WithLabelsEntries", func(t *testing.T) {
		t.Parallel()

		t.Run("adds each key to Labels", func(t *testing.T) {
			t.Parallel()

			var first, want map[string]string
			first = map[string]string{"Labels-a": "Labels-a"}
			want = map[string]string{"Labels-b": "Labels-a"}
			merged := maps.Clone(first)
			maps.Copy(merged, want)
			got := catalog.NewAddressBuilder().
				WithLabels(nil).
				WithLabelsEntries(first).
				WithLabelsEntries(want).
				Build()
			assert.Equal(
				t,
				got.Labels,
				merged,
				"WithLabelsEntries keeps the entries of its earlier call",
			)
		})
	})

	t.Run("WithTagsEntry", func(t *testing.T) {
		t.Parallel()

		t.Run("adds one key to Tags", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewAddressBuilder().WithTags(nil).WithTagsEntry("").Build()
			assert.Equal(
				t,
				got.Tags,
				map[string]struct{}{"": {}},
				"WithTagsEntry adds the key to a nil set",
			)
		})
	})

	t.Run("WithTagsEntries", func(t *testing.T) {
		t.Parallel()

		t.Run("adds each key to Tags", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewAddressBuilder().WithTags(nil).WithTagsEntries("", "").Build()
			assert.Equal(
				t,
				got.Tags,
				map[string]struct{}{"": {}},
				"WithTagsEntries adds two equal keys as one key",
			)
		})
	})

	t.Run("WithRawString", func(t *testing.T) {
		t.Parallel()

		t.Run("writes the bytes of the text to Raw", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewAddressBuilder().WithRawString("a").Build()
			assert.Equal(t, got.Raw, []byte("a"), "WithRawString writes the bytes of the text")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*catalog.Address
			record := func(v *catalog.Address) { got = append(got, v) }
			catalog.NewAddressBuilder().Mutate(record).Mutate(record)
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

			original := catalog.NewAddressBuilder()
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

			b := catalog.NewAddressBuilder()
			b.WithLines("")
			b.WithLabels(map[string]string{"": ""})
			b.WithTags(map[string]struct{}{"": {}})
			b.WithRaw([]byte{0})
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})

		t.Run("copies Lines into the clone", func(t *testing.T) {
			t.Parallel()

			var want, other []string
			want = []string{"Lines-a"}
			other = []string{"Lines-b"}
			original := catalog.NewAddressBuilder().WithLines(want...)
			before := slices.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *catalog.Address) { copy(v.Lines, other) })
			assert.Equal(
				t,
				original.Build().Lines,
				before,
				"a change of Lines in the clone leaves the original unchanged",
			)
		})

		t.Run("copies Labels into the clone", func(t *testing.T) {
			t.Parallel()

			var want, other map[string]string
			want = map[string]string{"Labels-a": "Labels-a"}
			other = map[string]string{"Labels-b": "Labels-a"}
			original := catalog.NewAddressBuilder().WithLabels(want)
			before := maps.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *catalog.Address) { maps.Copy(v.Labels, other) })
			assert.Equal(
				t,
				original.Build().Labels,
				before,
				"a change of Labels in the clone leaves the original unchanged",
			)
		})

		t.Run("copies Tags into the clone", func(t *testing.T) {
			t.Parallel()

			var want, other map[string]struct{}
			want = map[string]struct{}{"Tags-a": {}}
			other = map[string]struct{}{"Tags-b": {}}
			original := catalog.NewAddressBuilder().WithTags(want)
			before := maps.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *catalog.Address) { maps.Copy(v.Tags, other) })
			assert.Equal(
				t,
				original.Build().Tags,
				before,
				"a change of Tags in the clone leaves the original unchanged",
			)
		})

		t.Run("copies Raw into the clone", func(t *testing.T) {
			t.Parallel()

			var want, other []byte
			want = []byte{1}
			other = []byte{2}
			original := catalog.NewAddressBuilder().WithRaw(want)
			before := slices.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *catalog.Address) { copy(v.Raw, other) })
			assert.Equal(
				t,
				original.Build().Raw,
				before,
				"a change of Raw in the clone leaves the original unchanged",
			)
		})
	})

	t.Run("Build", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value under construction", func(t *testing.T) {
			t.Parallel()

			var seed catalog.Address
			seed.Street = "Street-a"
			seed.City = "City-a"
			seed.Grid = [2]int{42}
			seed.Timeout = time.Duration(1)
			b := catalog.NewAddressBuilder()
			b.Mutate(func(v *catalog.Address) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestProfileBuilder checks the constructors and the methods of
// ProfileBuilder.
func TestProfileBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewProfileBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed catalog.Profile
			seed.Active = true
			got := catalog.NewProfileBuilderFrom(seed).Build()
			assert.Equal(t, got.Active, seed.Active, "NewProfileBuilderFrom keeps Active")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewProfileBuilder()
			want := b.Build()
			assert.Equal(
				t,
				catalog.NewProfileBuilderFrom(want).Build(),
				want,
				"NewProfileBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithActive", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Active", func(t *testing.T) {
			t.Parallel()

			var want bool
			want = true
			got := catalog.NewProfileBuilder().WithActive(want).Build()
			assert.Equal(t, got.Active, want, "WithActive writes its argument to Active")
		})

		t.Run("replaces the value of Active", func(t *testing.T) {
			t.Parallel()

			var first, want bool
			first = true
			want = false
			got := catalog.NewProfileBuilder().WithActive(first).WithActive(want).Build()
			assert.Equal(t, got.Active, want, "the second call of WithActive replaces the first")
		})
	})

	t.Run("WithRole", func(t *testing.T) {
		t.Parallel()

		t.Run("returns its builder", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewProfileBuilder()
			assert.Equal(
				t,
				b.WithRole(b.Build().Role),
				b,
				"WithRole returns its builder for the next call",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("WithAudit", func(t *testing.T) {
		t.Parallel()

		t.Run("returns its builder", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewProfileBuilder()
			assert.Equal(
				t,
				b.WithAudit(b.Build().Audit),
				b,
				"WithAudit returns its builder for the next call",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*catalog.Profile
			record := func(v *catalog.Profile) { got = append(got, v) }
			catalog.NewProfileBuilder().Mutate(record).Mutate(record)
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

			original := catalog.NewProfileBuilder()
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

			b := catalog.NewProfileBuilder()
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

			var seed catalog.Profile
			seed.Active = true
			b := catalog.NewProfileBuilder()
			b.Mutate(func(v *catalog.Profile) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestEnvelopeBuilder checks the constructors and the methods of
// EnvelopeBuilder.
func TestEnvelopeBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewEnvelopeBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewEnvelopeBuilder()
			want := b.Build()
			assert.Equal(
				t,
				catalog.NewEnvelopeBuilderFrom(want).Build(),
				want,
				"NewEnvelopeBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithNotifier", func(t *testing.T) {
		t.Parallel()

		t.Run("returns its builder", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewEnvelopeBuilder()
			assert.Equal(
				t,
				b.WithNotifier(b.Build().Notifier),
				b,
				"WithNotifier returns its builder for the next call",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("WithCodes", func(t *testing.T) {
		t.Parallel()

		t.Run("returns its builder", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewEnvelopeBuilder()
			assert.Equal(
				t,
				b.WithCodes(b.Build().Codes),
				b,
				"WithCodes returns its builder for the next call",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*catalog.Envelope
			record := func(v *catalog.Envelope) { got = append(got, v) }
			catalog.NewEnvelopeBuilder().Mutate(record).Mutate(record)
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

			original := catalog.NewEnvelopeBuilder()
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

			b := catalog.NewEnvelopeBuilder()
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})
	})
}

// TestHiddenBuilder checks the constructors and the methods of
// HiddenBuilder.
func TestHiddenBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewHiddenBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed catalog.Hidden
			seed.Name = "Name-a"
			got := catalog.NewHiddenBuilderFrom(seed).Build()
			assert.Equal(t, got.Name, seed.Name, "NewHiddenBuilderFrom keeps Name")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewHiddenBuilder()
			want := b.Build()
			assert.Equal(
				t,
				catalog.NewHiddenBuilderFrom(want).Build(),
				want,
				"NewHiddenBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithName", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Name", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Name-a"
			got := catalog.NewHiddenBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := catalog.NewHiddenBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*catalog.Hidden
			record := func(v *catalog.Hidden) { got = append(got, v) }
			catalog.NewHiddenBuilder().Mutate(record).Mutate(record)
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

			original := catalog.NewHiddenBuilder()
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

			b := catalog.NewHiddenBuilder()
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

			var seed catalog.Hidden
			seed.Name = "Name-a"
			b := catalog.NewHiddenBuilder()
			b.Mutate(func(v *catalog.Hidden) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestOrderBuilder checks the constructors and the methods of
// OrderBuilder.
func TestOrderBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewOrderBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed catalog.Order
			seed.Codes = catalog.Codes{"Codes-a"}
			seed.Marker = struct{}{}
			got := catalog.NewOrderBuilderFrom(seed).Build()
			expect.Equal(t, got.Codes, seed.Codes, "NewOrderBuilderFrom keeps Codes")
			expect.Equal(t, got.Marker, seed.Marker, "NewOrderBuilderFrom keeps Marker")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewOrderBuilder()
			b.WithRemarks("")
			b.WithFlags(map[string]struct{}{"": {}})
			b.WithMarks(struct{}{})
			want := b.Build()
			assert.Equal(
				t,
				catalog.NewOrderBuilderFrom(want).Build(),
				want,
				"NewOrderBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithRemarks", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Remarks", func(t *testing.T) {
			t.Parallel()

			var want catalog.Remarks
			want = []string{"Remarks-a"}
			got := catalog.NewOrderBuilder().WithRemarks(want...).Build()
			assert.Equal(t, got.Remarks, want, "WithRemarks writes its argument to Remarks")
		})

		t.Run("replaces the value of Remarks", func(t *testing.T) {
			t.Parallel()

			var first, want catalog.Remarks
			first = []string{"Remarks-a"}
			want = []string{"Remarks-b"}
			got := catalog.NewOrderBuilder().WithRemarks(first...).WithRemarks(want...).Build()
			assert.Equal(t, got.Remarks, want, "the second call of WithRemarks replaces the first")
		})
	})

	t.Run("WithCodes", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Codes", func(t *testing.T) {
			t.Parallel()

			var want catalog.Codes
			want = catalog.Codes{"Codes-a"}
			got := catalog.NewOrderBuilder().WithCodes(want).Build()
			assert.Equal(t, got.Codes, want, "WithCodes writes its argument to Codes")
		})

		t.Run("replaces the value of Codes", func(t *testing.T) {
			t.Parallel()

			var first, want catalog.Codes
			first = catalog.Codes{"Codes-a"}
			want = catalog.Codes{"Codes-b"}
			got := catalog.NewOrderBuilder().WithCodes(first).WithCodes(want).Build()
			assert.Equal(t, got.Codes, want, "the second call of WithCodes replaces the first")
		})
	})

	t.Run("WithFlags", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Flags", func(t *testing.T) {
			t.Parallel()

			var want map[string]struct{}
			want = map[string]struct{}{"Flags-a": {}}
			got := catalog.NewOrderBuilder().WithFlags(want).Build()
			assert.Equal(t, got.Flags, want, "WithFlags writes its argument to Flags")
		})

		t.Run("replaces the value of Flags", func(t *testing.T) {
			t.Parallel()

			var first, want map[string]struct{}
			first = map[string]struct{}{"Flags-a": {}}
			want = map[string]struct{}{"Flags-b": {}}
			got := catalog.NewOrderBuilder().WithFlags(first).WithFlags(want).Build()
			assert.Equal(t, got.Flags, want, "the second call of WithFlags replaces the first")
		})
	})

	t.Run("WithMarker", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Marker", func(t *testing.T) {
			t.Parallel()

			var want struct{}
			want = struct{}{}
			got := catalog.NewOrderBuilder().WithMarker(want).Build()
			assert.Equal(t, got.Marker, want, "WithMarker writes its argument to Marker")
		})
	})

	t.Run("WithMarks", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Marks", func(t *testing.T) {
			t.Parallel()

			var want []struct{}
			want = []struct{}{{}}
			got := catalog.NewOrderBuilder().WithMarks(want...).Build()
			assert.Equal(t, got.Marks, want, "WithMarks writes its argument to Marks")
		})
	})

	t.Run("AppendRemarks", func(t *testing.T) {
		t.Parallel()

		t.Run("appends to the elements of Remarks", func(t *testing.T) {
			t.Parallel()

			var first, want catalog.Remarks
			first = []string{"Remarks-a"}
			want = []string{"Remarks-b"}
			got := catalog.NewOrderBuilder().WithRemarks(first...).AppendRemarks(want...).Build()
			assert.Equal(
				t,
				got.Remarks,
				slices.Concat(first, want),
				"AppendRemarks appends its elements after those of WithRemarks",
			)
		})
	})

	t.Run("WithFlagsEntry", func(t *testing.T) {
		t.Parallel()

		t.Run("adds one key to Flags", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewOrderBuilder().WithFlags(nil).WithFlagsEntry("").Build()
			assert.Equal(
				t,
				got.Flags,
				map[string]struct{}{"": {}},
				"WithFlagsEntry adds the key to a nil set",
			)
		})
	})

	t.Run("WithFlagsEntries", func(t *testing.T) {
		t.Parallel()

		t.Run("adds each key to Flags", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewOrderBuilder().WithFlags(nil).WithFlagsEntries("", "").Build()
			assert.Equal(
				t,
				got.Flags,
				map[string]struct{}{"": {}},
				"WithFlagsEntries adds two equal keys as one key",
			)
		})
	})

	t.Run("AppendMarks", func(t *testing.T) {
		t.Parallel()

		t.Run("appends to the elements of Marks", func(t *testing.T) {
			t.Parallel()

			var want []struct{}
			want = []struct{}{{}}
			got := catalog.NewOrderBuilder().WithMarks(want...).AppendMarks(want...).Build()
			assert.Equal(
				t,
				got.Marks,
				slices.Concat(want, want),
				"AppendMarks appends its elements after those of WithMarks",
			)
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*catalog.Order
			record := func(v *catalog.Order) { got = append(got, v) }
			catalog.NewOrderBuilder().Mutate(record).Mutate(record)
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

			original := catalog.NewOrderBuilder()
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

			b := catalog.NewOrderBuilder()
			b.WithRemarks("")
			b.WithFlags(map[string]struct{}{"": {}})
			b.WithMarks(struct{}{})
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})

		t.Run("copies Remarks into the clone", func(t *testing.T) {
			t.Parallel()

			var want, other catalog.Remarks
			want = []string{"Remarks-a"}
			other = []string{"Remarks-b"}
			original := catalog.NewOrderBuilder().WithRemarks(want...)
			before := slices.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *catalog.Order) { copy(v.Remarks, other) })
			assert.Equal(
				t,
				original.Build().Remarks,
				before,
				"a change of Remarks in the clone leaves the original unchanged",
			)
		})

		t.Run("copies Flags into the clone", func(t *testing.T) {
			t.Parallel()

			var want, other map[string]struct{}
			want = map[string]struct{}{"Flags-a": {}}
			other = map[string]struct{}{"Flags-b": {}}
			original := catalog.NewOrderBuilder().WithFlags(want)
			before := maps.Clone(want)
			clone := original.Clone()
			clone.Mutate(func(v *catalog.Order) { maps.Copy(v.Flags, other) })
			assert.Equal(
				t,
				original.Build().Flags,
				before,
				"a change of Flags in the clone leaves the original unchanged",
			)
		})
	})

	t.Run("Build", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the value under construction", func(t *testing.T) {
			t.Parallel()

			var seed catalog.Order
			seed.Codes = catalog.Codes{"Codes-a"}
			seed.Marker = struct{}{}
			b := catalog.NewOrderBuilder()
			b.Mutate(func(v *catalog.Order) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestLooseBuilder checks the constructors and the methods of
// LooseBuilder.
func TestLooseBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewLooseBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewLooseBuilder()
			b.WithItems(nil)
			b.WithIndex(map[string]catalog.Notifier{"": nil})
			b.WithKeys(map[catalog.Notifier]struct{}{nil: {}})
			want := b.Build()
			assert.Equal(
				t,
				catalog.NewLooseBuilderFrom(want).Build(),
				want,
				"NewLooseBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithItems", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Items", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewLooseBuilder().WithItems(nil).Build()
			assert.Equal(
				t,
				got.Items,
				[]catalog.Notifier{nil},
				"WithItems writes its argument to Items",
			)
		})
	})

	t.Run("AppendItems", func(t *testing.T) {
		t.Parallel()

		t.Run("appends to the elements of Items", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewLooseBuilder().WithItems(nil).AppendItems(nil).Build()
			assert.Equal(t, got.Items, []catalog.Notifier{nil, nil},
				"AppendItems appends its elements after those of WithItems")
		})
	})

	t.Run("WithIndex", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Index", func(t *testing.T) {
			t.Parallel()

			want := map[string]catalog.Notifier{"": nil}
			got := catalog.NewLooseBuilder().WithIndex(want).Build()
			assert.Equal(t, got.Index, want, "WithIndex writes its argument to Index")
		})
	})

	t.Run("WithIndexEntry", func(t *testing.T) {
		t.Parallel()

		t.Run("adds one key to Index", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewLooseBuilder().WithIndex(nil).WithIndexEntry("", nil).Build()
			assert.Equal(
				t,
				got.Index,
				map[string]catalog.Notifier{"": nil},
				"WithIndexEntry adds the key to a nil map",
			)
		})
	})

	t.Run("WithIndexEntries", func(t *testing.T) {
		t.Parallel()

		t.Run("adds each key to Index", func(t *testing.T) {
			t.Parallel()

			entries := map[string]catalog.Notifier{"": nil}
			got := catalog.NewLooseBuilder().WithIndex(nil).WithIndexEntries(entries).Build()
			assert.Equal(
				t,
				got.Index,
				entries,
				"WithIndexEntries copies its entries into a nil map",
			)
		})
	})

	t.Run("WithKeys", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Keys", func(t *testing.T) {
			t.Parallel()

			want := map[catalog.Notifier]struct{}{nil: {}}
			got := catalog.NewLooseBuilder().WithKeys(want).Build()
			assert.Equal(t, got.Keys, want, "WithKeys writes its argument to Keys")
		})
	})

	t.Run("WithKeysEntry", func(t *testing.T) {
		t.Parallel()

		t.Run("adds one key to Keys", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewLooseBuilder().WithKeys(nil).WithKeysEntry(nil).Build()
			assert.Equal(
				t,
				got.Keys,
				map[catalog.Notifier]struct{}{nil: {}},
				"WithKeysEntry adds the key to a nil set",
			)
		})
	})

	t.Run("WithKeysEntries", func(t *testing.T) {
		t.Parallel()

		t.Run("adds each key to Keys", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewLooseBuilder().WithKeys(nil).WithKeysEntries(nil, nil).Build()
			assert.Equal(
				t,
				got.Keys,
				map[catalog.Notifier]struct{}{nil: {}},
				"WithKeysEntries adds two equal keys as one key",
			)
		})
	})

	t.Run("WithExtra", func(t *testing.T) {
		t.Parallel()

		t.Run("points Extra at its argument", func(t *testing.T) {
			t.Parallel()

			got := catalog.NewLooseBuilder().WithExtra(nil).Build()
			assert.NotNil(t, got.Extra, "WithExtra points Extra at a value")
			assert.Equal(t, *got.Extra, nil, "WithExtra points Extra at its argument")
		})
	})

	t.Run("WithHook", func(t *testing.T) {
		t.Parallel()

		t.Run("returns its builder", func(t *testing.T) {
			t.Parallel()

			b := catalog.NewLooseBuilder()
			assert.Equal(
				t,
				b.WithHook(b.Build().Hook),
				b,
				"WithHook returns its builder for the next call",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*catalog.Loose
			record := func(v *catalog.Loose) { got = append(got, v) }
			catalog.NewLooseBuilder().Mutate(record).Mutate(record)
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

			original := catalog.NewLooseBuilder()
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

			b := catalog.NewLooseBuilder()
			b.WithItems(nil)
			b.WithIndex(map[string]catalog.Notifier{"": nil})
			b.WithKeys(map[catalog.Notifier]struct{}{nil: {}})
			assert.Equal(
				t,
				b.Clone().Build(),
				b.Build(),
				"the clone builds the value of its original",
			)
		})
	})
}
