package seeded

import (
	"runtime"
	"testing"

	"go.dokimi.dev/assert"
)

// TestAccountBuilderConstruction pins what the constructors produce.
func TestAccountBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Account
		seed.Owner = "Owner-a"
		seed.Plan = "Plan-a"
		got := NewAccountBuilderFrom(seed).Build()
		assert.Equal(t, got.Owner, seed.Owner, "From keeps the member it was given")
		assert.Equal(t, got.Plan, seed.Plan, "From keeps the member it was given")
	})

	t.Run("Owner carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "owner"
		got := NewAccountBuilder().Build()
		assert.Equal(t, got.Owner, want, "the constructor applies the declared default")
	})
}

// TestAccountBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestAccountBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithOwner writes Owner", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Owner-a"
		got := NewAccountBuilder().WithOwner(want).Build()
		assert.Equal(t, got.Owner, want, "the setter writes what it was given")
	})

	t.Run("WithOwner replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Owner-a"
		want = "Owner-b"
		got := NewAccountBuilder().WithOwner(first).WithOwner(want).Build()
		assert.Equal(t, got.Owner, want, "the later call wins")
	})

	t.Run("WithPlan writes Plan", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Plan-a"
		got := NewAccountBuilder().WithPlan(want).Build()
		assert.Equal(t, got.Plan, want, "the setter writes what it was given")
	})

	t.Run("WithPlan replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Plan-a"
		want = "Plan-b"
		got := NewAccountBuilder().WithPlan(first).WithPlan(want).Build()
		assert.Equal(t, got.Plan, want, "the later call wins")
	})
}

// TestAccountBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestAccountBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Account
		record := func(v *Account) { got = append(got, v) }
		NewAccountBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewAccountBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewAccountBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestPolicyBuilderConstruction pins what the constructors produce.
func TestPolicyBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Policy
		seed.Attempts = 42
		got := NewPolicyBuilderFrom(seed).Build()
		assert.Equal(t, got.Attempts, seed.Attempts, "From keeps the member it was given")
	})
}

// TestPolicyBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestPolicyBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithAttempts writes Attempts", func(t *testing.T) {
		t.Parallel()

		var want int
		want = 42
		got := NewPolicyBuilder().WithAttempts(want).Build()
		assert.Equal(t, got.Attempts, want, "the setter writes what it was given")
	})

	t.Run("WithAttempts replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want int
		first = 42
		want = 7
		got := NewPolicyBuilder().WithAttempts(first).WithAttempts(want).Build()
		assert.Equal(t, got.Attempts, want, "the later call wins")
	})
}

// TestPolicyBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestPolicyBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Policy
		record := func(v *Policy) { got = append(got, v) }
		NewPolicyBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewPolicyBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewPolicyBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestTierBuilderConstruction pins what the constructors produce.
func TestTierBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Tier
		seed.Level = 42
		got := NewTierBuilderFrom(seed).Build()
		assert.Equal(t, got.Level, seed.Level, "From keeps the member it was given")
	})
}

// TestTierBuilderSetters pins that each setter writes the member it names,
// and writes nothing else.
func TestTierBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithLevel writes Level", func(t *testing.T) {
		t.Parallel()

		var want int
		want = 42
		got := NewTierBuilder().WithLevel(want).Build()
		assert.Equal(t, got.Level, want, "the setter writes what it was given")
	})

	t.Run("WithLevel replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want int
		first = 42
		want = 7
		got := NewTierBuilder().WithLevel(first).WithLevel(want).Build()
		assert.Equal(t, got.Level, want, "the later call wins")
	})
}

// TestTierBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestTierBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Tier
		record := func(v *Tier) { got = append(got, v) }
		NewTierBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewTierBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewTierBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}

// TestDeploymentBuilderConstruction pins what the constructors produce.
func TestDeploymentBuilderConstruction(t *testing.T) {
	t.Parallel()

	t.Run("From round-trips what it was seeded with", func(t *testing.T) {
		t.Parallel()

		var seed Deployment
		seed.Mode = Code("Mode-a")
		seed.Region = "Region-a"
		seed.Zone = "Zone-a"
		seed.System = "System-a"
		seed.Name = "Name-a"
		seed.Label = "Label-a"
		seed.Flag = "Flag-a"
		got := NewDeploymentBuilderFrom(seed).Build()
		assert.Equal(t, got.Mode, seed.Mode, "From keeps the member it was given")
		assert.Equal(t, got.Region, seed.Region, "From keeps the member it was given")
		assert.Equal(t, got.Zone, seed.Zone, "From keeps the member it was given")
		assert.Equal(t, got.System, seed.System, "From keeps the member it was given")
		assert.Equal(t, got.Name, seed.Name, "From keeps the member it was given")
		assert.Equal(t, got.Label, seed.Label, "From keeps the member it was given")
		assert.Equal(t, got.Flag, seed.Flag, "From keeps the member it was given")
	})

	t.Run("Mode carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want Code
		want = "true"
		got := NewDeploymentBuilder().Build()
		assert.Equal(t, got.Mode, want, "the constructor applies the declared default")
	})

	t.Run("Region carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = Region
		got := NewDeploymentBuilder().Build()
		assert.Equal(t, got.Region, want, "the constructor applies the declared default")
	})

	t.Run("Zone carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = zone()
		got := NewDeploymentBuilder().Build()
		assert.Equal(t, got.Zone, want, "the constructor applies the declared default")
	})

	t.Run("System carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = runtime.GOOS
		got := NewDeploymentBuilder().Build()
		assert.Equal(t, got.System, want, "the constructor applies the declared default")
	})

	t.Run("Name carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "web"
		got := NewDeploymentBuilder().Build()
		assert.Equal(t, got.Name, want, "the constructor applies the declared default")
	})

	t.Run("Label carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "a b"
		got := NewDeploymentBuilder().Build()
		assert.Equal(t, got.Label, want, "the constructor applies the declared default")
	})

	t.Run("Flag carries its declared default", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "true"
		got := NewDeploymentBuilder().Build()
		assert.Equal(t, got.Flag, want, "the constructor applies the declared default")
	})
}

// TestDeploymentBuilderSetters pins that each setter writes the member it
// names, and writes nothing else.
func TestDeploymentBuilderSetters(t *testing.T) {
	t.Parallel()

	t.Run("WithMode writes Mode", func(t *testing.T) {
		t.Parallel()

		var want Code
		want = Code("Mode-a")
		got := NewDeploymentBuilder().WithMode(want).Build()
		assert.Equal(t, got.Mode, want, "the setter writes what it was given")
	})

	t.Run("WithMode replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want Code
		first = Code("Mode-a")
		want = Code("Mode-b")
		got := NewDeploymentBuilder().WithMode(first).WithMode(want).Build()
		assert.Equal(t, got.Mode, want, "the later call wins")
	})

	t.Run("WithRegion writes Region", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Region-a"
		got := NewDeploymentBuilder().WithRegion(want).Build()
		assert.Equal(t, got.Region, want, "the setter writes what it was given")
	})

	t.Run("WithRegion replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Region-a"
		want = "Region-b"
		got := NewDeploymentBuilder().WithRegion(first).WithRegion(want).Build()
		assert.Equal(t, got.Region, want, "the later call wins")
	})

	t.Run("WithZone writes Zone", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Zone-a"
		got := NewDeploymentBuilder().WithZone(want).Build()
		assert.Equal(t, got.Zone, want, "the setter writes what it was given")
	})

	t.Run("WithZone replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Zone-a"
		want = "Zone-b"
		got := NewDeploymentBuilder().WithZone(first).WithZone(want).Build()
		assert.Equal(t, got.Zone, want, "the later call wins")
	})

	t.Run("WithSystem writes System", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "System-a"
		got := NewDeploymentBuilder().WithSystem(want).Build()
		assert.Equal(t, got.System, want, "the setter writes what it was given")
	})

	t.Run("WithSystem replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "System-a"
		want = "System-b"
		got := NewDeploymentBuilder().WithSystem(first).WithSystem(want).Build()
		assert.Equal(t, got.System, want, "the later call wins")
	})

	t.Run("WithName writes Name", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Name-a"
		got := NewDeploymentBuilder().WithName(want).Build()
		assert.Equal(t, got.Name, want, "the setter writes what it was given")
	})

	t.Run("WithName replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Name-a"
		want = "Name-b"
		got := NewDeploymentBuilder().WithName(first).WithName(want).Build()
		assert.Equal(t, got.Name, want, "the later call wins")
	})

	t.Run("WithLabel writes Label", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Label-a"
		got := NewDeploymentBuilder().WithLabel(want).Build()
		assert.Equal(t, got.Label, want, "the setter writes what it was given")
	})

	t.Run("WithLabel replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Label-a"
		want = "Label-b"
		got := NewDeploymentBuilder().WithLabel(first).WithLabel(want).Build()
		assert.Equal(t, got.Label, want, "the later call wins")
	})

	t.Run("WithFlag writes Flag", func(t *testing.T) {
		t.Parallel()

		var want string
		want = "Flag-a"
		got := NewDeploymentBuilder().WithFlag(want).Build()
		assert.Equal(t, got.Flag, want, "the setter writes what it was given")
	})

	t.Run("WithFlag replaces what was already there", func(t *testing.T) {
		t.Parallel()

		var first, want string
		first = "Flag-a"
		want = "Flag-b"
		got := NewDeploymentBuilder().WithFlag(first).WithFlag(want).Build()
		assert.Equal(t, got.Flag, want, "the later call wins")
	})
}

// TestDeploymentBuilderMethods pins that Mutate edits the value under
// construction, and that Clone and From keep that value.
func TestDeploymentBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("Mutate passes the value under construction", func(t *testing.T) {
		t.Parallel()

		var got []*Deployment
		record := func(v *Deployment) { got = append(got, v) }
		NewDeploymentBuilder().Mutate(record).Mutate(record)
		assert.Length(t, got, 2, "Mutate calls each function once")
		assert.Equal(t, got[0], got[1], "the two functions get the same value", assert.ByIdentity())
	})

	t.Run("Clone and From keep the value under construction", func(t *testing.T) {
		t.Parallel()

		original := NewDeploymentBuilder()
		want := original.Build()
		clone := original.Clone()
		assert.NotEqual(t, clone, original, "Clone returns a builder of its own", assert.ByIdentity())
		assert.Equal(t, clone.Build(), want, "a clone builds the value of its original")
		assert.Equal(t, NewDeploymentBuilderFrom(want).Build(), want, "From keeps the value it was given")
	})
}
