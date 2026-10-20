package catalog

import (
	"maps"
	"testing"
	"time"

	"go.dokimi.dev/assert"
)

// TestAddressBuilderConstruction pins what the constructors produce.
func TestAddressBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Address
		seed.Street = "Street-a"
		seed.City = "City-a"
		seed.Grid = [2]int{42}
		seed.Timeout = time.Duration(1)
		got := NewAddressBuilderFrom(seed).Build()
		assert.Equal(t, got.Street, seed.Street, "From keeps the member it was given")
		assert.Equal(t, got.City, seed.City, "From keeps the member it was given")
		assert.Equal(t, got.Grid, seed.Grid, "From keeps the member it was given")
		assert.Equal(t, got.Timeout, seed.Timeout, "From keeps the member it was given")
	})

	t.Run("Street carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Dam 1"
		got := NewAddressBuilder().Build()
		assert.Equal(t, got.Street, want, "the constructor applies the declared default")
	})

	t.Run("City carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Amsterdam"
		got := NewAddressBuilder().Build()
		assert.Equal(t, got.City, want, "the constructor applies the declared default")
	})

	t.Run("Timeout carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want time.Duration
		want = time.Second
		got := NewAddressBuilder().Build()
		assert.Equal(t, got.Timeout, want, "the constructor applies the declared default")
	})
}

// TestAddressBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestAddressBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithStreet writes Street", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Street-a"
		got := NewAddressBuilder().WithStreet(want).Build()
		assert.Equal(t, got.Street, want, "the setter writes what it was given")
	})

	t.Run("WithStreet replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Street-a"
		want = "Street-b"
		got := NewAddressBuilder().WithStreet(first).WithStreet(want).Build()
		assert.Equal(t, got.Street, want, "the later call wins")
	})

	t.Run("WithCity writes City", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "City-a"
		got := NewAddressBuilder().WithCity(want).Build()
		assert.Equal(t, got.City, want, "the setter writes what it was given")
	})

	t.Run("WithCity replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "City-a"
		want = "City-b"
		got := NewAddressBuilder().WithCity(first).WithCity(want).Build()
		assert.Equal(t, got.City, want, "the later call wins")
	})

	t.Run("WithFloor writes Floor", func(t *testing.T) {
		t.Parallel()

		var want int
		want = 42
		got := NewAddressBuilder().WithFloor(want).Build()
		assert.NotNil(t, got.Floor, "the setter addresses the value it was given")
	})

	t.Run("WithLines writes Lines", func(t *testing.T) {
		t.Parallel()

		var want []string
		want = []string{"Lines-a"}
		got := NewAddressBuilder().WithLines(want...).Build()
		assert.Length(t, got.Lines, 1, "the setter replaces the sequence")
	})

	t.Run("WithLabels writes Labels", func(t *testing.T) {
		t.Parallel()

		var want map[string]string
		want = map[string]string{"Labels-a": "Labels-a"}
		got := NewAddressBuilder().WithLabels(want).Build()
		assert.Length(t, got.Labels, 1, "the setter replaces the collection")
	})

	t.Run("WithTags writes Tags", func(t *testing.T) {
		t.Parallel()

		var want map[string]struct{}
		want = map[string]struct{}{"Tags-a": struct{}{}}
		got := NewAddressBuilder().WithTags(want).Build()
		assert.Length(t, got.Tags, 1, "the setter replaces the collection")
	})

	t.Run("WithRaw writes Raw", func(t *testing.T) {
		t.Parallel()

		var want []byte
		want = []byte{1}
		got := NewAddressBuilder().WithRaw(want).Build()
		assert.Length(t, got.Raw, 1, "the setter replaces the sequence")
	})

	t.Run("WithGrid writes Grid", func(t *testing.T) {
		t.Parallel()

		var want [2]int
		want = [2]int{42}
		got := NewAddressBuilder().WithGrid(want).Build()
		assert.Equal(t, got.Grid, want, "the setter writes what it was given")
	})

	t.Run("WithGrid replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want [2]int
		first = [2]int{42}
		want = [2]int{7}
		got := NewAddressBuilder().WithGrid(first).WithGrid(want).Build()
		assert.Equal(t, got.Grid, want, "the later call wins")
	})

	t.Run("WithTimeout writes Timeout", func(t *testing.T) {
		t.Parallel()

		var want time.Duration
		want = time.Duration(1)
		got := NewAddressBuilder().WithTimeout(want).Build()
		assert.Equal(t, got.Timeout, want, "the setter writes what it was given")
	})

	t.Run("WithTimeout replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want time.Duration
		first = time.Duration(1)
		want = time.Duration(2)
		got := NewAddressBuilder().WithTimeout(first).WithTimeout(want).Build()
		assert.Equal(t, got.Timeout, want, "the later call wins")
	})

	t.Run("AppendLines keeps what Lines holds", func(t *testing.T) {
		t.Parallel()

		got := NewAddressBuilder().WithLines(*new(string)).AppendLines(*new(string)).Build()
		assert.Length(t, got.Lines, 2, "the appender adds to the sequence")
	})

	t.Run("WithLabelsEntry adds one key to Labels", func(t *testing.T) {
		t.Parallel()

		got := NewAddressBuilder().WithLabels(nil).WithLabelsEntry(*new(string), *new(string)).Build()
		assert.Length(t, got.Labels, 1, "the setter adds the key to an empty collection")
	})

	t.Run("WithLabelsEntries adds each key to Labels", func(t *testing.T) {
		t.Parallel()

		entries := map[string]string{*new(string): *new(string)}
		got := NewAddressBuilder().WithLabels(nil).WithLabelsEntries(entries).Build()
		assert.Length(t, got.Labels, 1, "the setter merges the entries into an empty collection")
	})

	t.Run("WithTagsEntry adds one key to Tags", func(t *testing.T) {
		t.Parallel()

		got := NewAddressBuilder().WithTags(nil).WithTagsEntry(*new(string)).Build()
		assert.Length(t, got.Tags, 1, "the setter adds the key to an empty collection")
	})

	t.Run("WithTagsEntries adds each key to Tags", func(t *testing.T) {
		t.Parallel()

		got := NewAddressBuilder().WithTags(nil).WithTagsEntries(*new(string), *new(string)).Build()
		assert.Length(t, got.Tags, 1, "the setter adds two equal keys as one entry")
	})

	t.Run("WithRawString writes the bytes of the text to Raw", func(t *testing.T) {
		t.Parallel()

		got := NewAddressBuilder().WithRawString("a").Build()
		assert.Equal(t, got.Raw, []byte("a"), "the setter writes the bytes of the text")
	})
}

// TestAddressBuilderClone pins that a clone can be configured without the
// original seeing it.
func TestAddressBuilderClone(t *testing.T) {
	t.Parallel()

	t.Run("Lines is not shared with a clone", func(t *testing.T) {
		t.Parallel()

		var want, other []string
		want = []string{"Lines-a"}
		other = []string{"Lines-b"}
		original := NewAddressBuilder().WithLines(want...)
		clone := original.Clone()
		clone.Mutate(func(v *Address) { copy(v.Lines, other) })
		assert.NotEqual(t, original.Build().Lines, other, "writing through a clone must not reach the original")
	})

	t.Run("Labels is not shared with a clone", func(t *testing.T) {
		t.Parallel()

		var want, other map[string]string
		want = map[string]string{"Labels-a": "Labels-a"}
		other = map[string]string{"Labels-b": "Labels-a"}
		original := NewAddressBuilder().WithLabels(want)
		clone := original.Clone()
		clone.Mutate(func(v *Address) { maps.Copy(v.Labels, other) })
		assert.Length(t, original.Build().Labels, 1, "writing through a clone must not reach the original")
	})

	t.Run("Tags is not shared with a clone", func(t *testing.T) {
		t.Parallel()

		var want, other map[string]struct{}
		want = map[string]struct{}{"Tags-a": struct{}{}}
		other = map[string]struct{}{"Tags-b": struct{}{}}
		original := NewAddressBuilder().WithTags(want)
		clone := original.Clone()
		clone.Mutate(func(v *Address) { maps.Copy(v.Tags, other) })
		assert.Length(t, original.Build().Tags, 1, "writing through a clone must not reach the original")
	})

	t.Run("Raw is not shared with a clone", func(t *testing.T) {
		t.Parallel()

		var want, other []byte
		want = []byte{1}
		other = []byte{2}
		original := NewAddressBuilder().WithRaw(want)
		clone := original.Clone()
		clone.Mutate(func(v *Address) { copy(v.Raw, other) })
		assert.NotEqual(t, original.Build().Raw, other, "writing through a clone must not reach the original")
	})
}

// TestAddressBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestAddressBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Address
		record := func(v *Address) { got = append(got, v) }
		NewAddressBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewAddressBuilder()
		original.WithLines(*new(string))
		original.WithLabels(map[string]string{*new(string): *new(string)})
		original.WithTags(map[string]struct{}{*new(string): {}})
		original.WithRaw([]byte{0})
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewAddressBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestProfileBuilderConstruction pins what the constructors produce.
func TestProfileBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Profile
		seed.Active = true
		got := NewProfileBuilderFrom(seed).Build()
		assert.Equal(t, got.Active, seed.Active, "From keeps the member it was given")
	})
}

// TestProfileBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestProfileBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithActive writes Active", func(t *testing.T) {
		t.Parallel()

		var want bool
		want = true
		got := NewProfileBuilder().WithActive(want).Build()
		assert.Equal(t, got.Active, want, "the setter writes what it was given")
	})

	t.Run("WithActive replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want bool
		first = true
		want = false
		got := NewProfileBuilder().WithActive(first).WithActive(want).Build()
		assert.Equal(t, got.Active, want, "the later call wins")
	})

	t.Run("WithRole returns its builder", func(t *testing.T) {
		t.Parallel()

		b := NewProfileBuilder()
		assert.Equal(t, b.WithRole(b.Build().Role), b, "the setter returns its builder for the next call",
			assert.ByIdentity())
	})

	t.Run("WithAudit returns its builder", func(t *testing.T) {
		t.Parallel()

		b := NewProfileBuilder()
		assert.Equal(t, b.WithAudit(b.Build().Audit), b, "the setter returns its builder for the next call",
			assert.ByIdentity())
	})
}

// TestProfileBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestProfileBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Profile
		record := func(v *Profile) { got = append(got, v) }
		NewProfileBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewProfileBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewProfileBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestEnvelopeBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestEnvelopeBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithNotifier returns its builder", func(t *testing.T) {
		t.Parallel()

		b := NewEnvelopeBuilder()
		assert.Equal(t, b.WithNotifier(b.Build().Notifier), b, "the setter returns its builder for the next call",
			assert.ByIdentity())
	})

	t.Run("WithCodes returns its builder", func(t *testing.T) {
		t.Parallel()

		b := NewEnvelopeBuilder()
		assert.Equal(t, b.WithCodes(b.Build().Codes), b, "the setter returns its builder for the next call",
			assert.ByIdentity())
	})
}

// TestEnvelopeBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestEnvelopeBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Envelope
		record := func(v *Envelope) { got = append(got, v) }
		NewEnvelopeBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewEnvelopeBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewEnvelopeBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestHiddenBuilderConstruction pins what the constructors produce.
func TestHiddenBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Hidden
		seed.Name = "Name-a"
		got := NewHiddenBuilderFrom(seed).Build()
		assert.Equal(t, got.Name, seed.Name, "From keeps the member it was given")
	})
}

// TestHiddenBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestHiddenBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithName writes Name", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Name-a"
		got := NewHiddenBuilder().WithName(want).Build()
		assert.Equal(t, got.Name, want, "the setter writes what it was given")
	})

	t.Run("WithName replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Name-a"
		want = "Name-b"
		got := NewHiddenBuilder().WithName(first).WithName(want).Build()
		assert.Equal(t, got.Name, want, "the later call wins")
	})
}

// TestHiddenBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestHiddenBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Hidden
		record := func(v *Hidden) { got = append(got, v) }
		NewHiddenBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewHiddenBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewHiddenBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestOrderBuilderConstruction pins what the constructors produce.
func TestOrderBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Order
		seed.Codes = Codes([]string{"Codes-a"})
		seed.Marker = struct{}{}
		got := NewOrderBuilderFrom(seed).Build()
		assert.Equal(t, got.Codes, seed.Codes, "From keeps the member it was given")
		assert.Equal(t, got.Marker, seed.Marker, "From keeps the member it was given")
	})
}

// TestOrderBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestOrderBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithRemarks writes Remarks", func(t *testing.T) {
		t.Parallel()

		var want Remarks
		want = []string{"Remarks-a"}
		got := NewOrderBuilder().WithRemarks(want...).Build()
		assert.Length(t, got.Remarks, 1, "the setter replaces the sequence")
	})

	t.Run("WithCodes writes Codes", func(t *testing.T) {
		t.Parallel()

		var want Codes
		want = Codes([]string{"Codes-a"})
		got := NewOrderBuilder().WithCodes(want).Build()
		assert.Equal(t, got.Codes, want, "the setter writes what it was given")
	})

	t.Run("WithCodes replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want Codes
		first = Codes([]string{"Codes-a"})
		want = Codes([]string{"Codes-b"})
		got := NewOrderBuilder().WithCodes(first).WithCodes(want).Build()
		assert.Equal(t, got.Codes, want, "the later call wins")
	})

	t.Run("WithFlags writes Flags", func(t *testing.T) {
		t.Parallel()

		var want map[string]struct{}
		want = map[string]struct{}{"Flags-a": struct{}{}}
		got := NewOrderBuilder().WithFlags(want).Build()
		assert.Length(t, got.Flags, 1, "the setter replaces the collection")
	})

	t.Run("WithMarker writes Marker", func(t *testing.T) {
		t.Parallel()

		var want struct{}
		want = struct{}{}
		got := NewOrderBuilder().WithMarker(want).Build()
		assert.Equal(t, got.Marker, want, "the setter writes what it was given")
	})

	t.Run("WithMarks writes Marks", func(t *testing.T) {
		t.Parallel()

		var want []struct{}
		want = []struct{}{struct{}{}}
		got := NewOrderBuilder().WithMarks(want...).Build()
		assert.Length(t, got.Marks, 1, "the setter replaces the sequence")
	})

	t.Run("AppendRemarks keeps what Remarks holds", func(t *testing.T) {
		t.Parallel()

		got := NewOrderBuilder().WithRemarks(*new(string)).AppendRemarks(*new(string)).Build()
		assert.Length(t, got.Remarks, 2, "the appender adds to the sequence")
	})

	t.Run("WithFlagsEntry adds one key to Flags", func(t *testing.T) {
		t.Parallel()

		got := NewOrderBuilder().WithFlags(nil).WithFlagsEntry(*new(string)).Build()
		assert.Length(t, got.Flags, 1, "the setter adds the key to an empty collection")
	})

	t.Run("WithFlagsEntries adds each key to Flags", func(t *testing.T) {
		t.Parallel()

		got := NewOrderBuilder().WithFlags(nil).WithFlagsEntries(*new(string), *new(string)).Build()
		assert.Length(t, got.Flags, 1, "the setter adds two equal keys as one entry")
	})

	t.Run("AppendMarks keeps what Marks holds", func(t *testing.T) {
		t.Parallel()

		got := NewOrderBuilder().WithMarks(*new(struct{})).AppendMarks(*new(struct{})).Build()
		assert.Length(t, got.Marks, 2, "the appender adds to the sequence")
	})
}

// TestOrderBuilderClone pins that a clone can be configured without the
// original seeing it.
func TestOrderBuilderClone(t *testing.T) {
	t.Parallel()

	t.Run("Remarks is not shared with a clone", func(t *testing.T) {
		t.Parallel()

		var want, other Remarks
		want = []string{"Remarks-a"}
		other = []string{"Remarks-b"}
		original := NewOrderBuilder().WithRemarks(want...)
		clone := original.Clone()
		clone.Mutate(func(v *Order) { copy(v.Remarks, other) })
		assert.NotEqual(t, original.Build().Remarks, other, "writing through a clone must not reach the original")
	})

	t.Run("Flags is not shared with a clone", func(t *testing.T) {
		t.Parallel()

		var want, other map[string]struct{}
		want = map[string]struct{}{"Flags-a": struct{}{}}
		other = map[string]struct{}{"Flags-b": struct{}{}}
		original := NewOrderBuilder().WithFlags(want)
		clone := original.Clone()
		clone.Mutate(func(v *Order) { maps.Copy(v.Flags, other) })
		assert.Length(t, original.Build().Flags, 1, "writing through a clone must not reach the original")
	})
}

// TestOrderBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestOrderBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Order
		record := func(v *Order) { got = append(got, v) }
		NewOrderBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewOrderBuilder()
		original.WithRemarks(*new(string))
		original.WithFlags(map[string]struct{}{*new(string): {}})
		original.WithMarks(*new(struct{}))
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewOrderBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestLooseBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestLooseBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithItems writes Items", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithItems(*new(Notifier)).Build()
		assert.Length(t, got.Items, 1, "the setter replaces the sequence")
	})

	t.Run("AppendItems keeps what Items holds", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithItems(*new(Notifier)).AppendItems(*new(Notifier)).Build()
		assert.Length(t, got.Items, 2, "the appender adds to the sequence")
	})

	t.Run("WithIndex writes Index", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithIndex(map[string]Notifier{*new(string): *new(Notifier)}).Build()
		assert.Length(t, got.Index, 1, "the setter replaces the collection")
	})

	t.Run("WithIndexEntry adds one key to Index", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithIndex(nil).WithIndexEntry(*new(string), *new(Notifier)).Build()
		assert.Length(t, got.Index, 1, "the setter adds the key to an empty collection")
	})

	t.Run("WithIndexEntries adds each key to Index", func(t *testing.T) {
		t.Parallel()

		entries := map[string]Notifier{*new(string): *new(Notifier)}
		got := NewLooseBuilder().WithIndex(nil).WithIndexEntries(entries).Build()
		assert.Length(t, got.Index, 1, "the setter merges the entries into an empty collection")
	})

	t.Run("WithKeys writes Keys", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithKeys(map[Notifier]struct{}{*new(Notifier): {}}).Build()
		assert.Length(t, got.Keys, 1, "the setter replaces the collection")
	})

	t.Run("WithKeysEntry adds one key to Keys", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithKeys(nil).WithKeysEntry(*new(Notifier)).Build()
		assert.Length(t, got.Keys, 1, "the setter adds the key to an empty collection")
	})

	t.Run("WithKeysEntries adds each key to Keys", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithKeys(nil).WithKeysEntries(*new(Notifier), *new(Notifier)).Build()
		assert.Length(t, got.Keys, 1, "the setter adds two equal keys as one entry")
	})

	t.Run("WithExtra writes Extra", func(t *testing.T) {
		t.Parallel()

		got := NewLooseBuilder().WithExtra(*new(Notifier)).Build()
		assert.NotNil(t, got.Extra, "the setter addresses the value it was given")
	})

	t.Run("WithHook returns its builder", func(t *testing.T) {
		t.Parallel()

		b := NewLooseBuilder()
		assert.Equal(t, b.WithHook(b.Build().Hook), b, "the setter returns its builder for the next call",
			assert.ByIdentity())
	})
}

// TestLooseBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestLooseBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Loose
		record := func(v *Loose) { got = append(got, v) }
		NewLooseBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewLooseBuilder()
		original.WithItems(*new(Notifier))
		original.WithIndex(map[string]Notifier{*new(string): *new(Notifier)})
		original.WithKeys(map[Notifier]struct{}{*new(Notifier): {}})
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewLooseBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}
