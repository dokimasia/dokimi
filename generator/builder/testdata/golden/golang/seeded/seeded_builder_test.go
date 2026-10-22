package seeded

import (
	"runtime"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
)

// TestAccountBuilder checks the constructors and the methods of
// AccountBuilder.
func TestAccountBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewAccountBuilder", func(t *testing.T) {
		t.Parallel()

		t.Run("applies the declared default of Owner", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "owner"
			got := NewAccountBuilder().Build()
			assert.Equal(
				t,
				got.Owner,
				want,
				"NewAccountBuilder applies the declared default of Owner",
			)
		})
	})

	t.Run("NewAccountBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed Account
			seed.Owner = "Owner-a"
			seed.Plan = "Plan-a"
			got := NewAccountBuilderFrom(seed).Build()
			expect.Equal(t, got.Owner, seed.Owner, "NewAccountBuilderFrom keeps Owner")
			expect.Equal(t, got.Plan, seed.Plan, "NewAccountBuilderFrom keeps Plan")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewAccountBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewAccountBuilderFrom(want).Build(),
				want,
				"NewAccountBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithOwner", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Owner", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Owner-a"
			got := NewAccountBuilder().WithOwner(want).Build()
			assert.Equal(t, got.Owner, want, "WithOwner writes its argument to Owner")
		})

		t.Run("replaces the value of Owner", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Owner-a"
			want = "Owner-b"
			got := NewAccountBuilder().WithOwner(first).WithOwner(want).Build()
			assert.Equal(t, got.Owner, want, "the second call of WithOwner replaces the first")
		})
	})

	t.Run("WithPlan", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Plan", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Plan-a"
			got := NewAccountBuilder().WithPlan(want).Build()
			assert.Equal(t, got.Plan, want, "WithPlan writes its argument to Plan")
		})

		t.Run("replaces the value of Plan", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Plan-a"
			want = "Plan-b"
			got := NewAccountBuilder().WithPlan(first).WithPlan(want).Build()
			assert.Equal(t, got.Plan, want, "the second call of WithPlan replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Account
			record := func(v *Account) { got = append(got, v) }
			NewAccountBuilder().Mutate(record).Mutate(record)
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

			original := NewAccountBuilder()
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

			b := NewAccountBuilder()
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

			var seed Account
			seed.Owner = "Owner-a"
			seed.Plan = "Plan-a"
			b := NewAccountBuilder()
			b.Mutate(func(v *Account) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestPolicyBuilder checks the constructors and the methods of
// PolicyBuilder.
func TestPolicyBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewPolicyBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed Policy
			seed.Attempts = 42
			got := NewPolicyBuilderFrom(seed).Build()
			assert.Equal(t, got.Attempts, seed.Attempts, "NewPolicyBuilderFrom keeps Attempts")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewPolicyBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewPolicyBuilderFrom(want).Build(),
				want,
				"NewPolicyBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithAttempts", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Attempts", func(t *testing.T) {
			t.Parallel()

			var want int
			want = 42
			got := NewPolicyBuilder().WithAttempts(want).Build()
			assert.Equal(t, got.Attempts, want, "WithAttempts writes its argument to Attempts")
		})

		t.Run("replaces the value of Attempts", func(t *testing.T) {
			t.Parallel()

			var first, want int
			first = 42
			want = 7
			got := NewPolicyBuilder().WithAttempts(first).WithAttempts(want).Build()
			assert.Equal(
				t,
				got.Attempts,
				want,
				"the second call of WithAttempts replaces the first",
			)
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Policy
			record := func(v *Policy) { got = append(got, v) }
			NewPolicyBuilder().Mutate(record).Mutate(record)
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

			original := NewPolicyBuilder()
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

			b := NewPolicyBuilder()
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

			var seed Policy
			seed.Attempts = 42
			b := NewPolicyBuilder()
			b.Mutate(func(v *Policy) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestTierBuilder checks the constructors and the methods of TierBuilder.
func TestTierBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewTierBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed Tier
			seed.Level = 42
			got := NewTierBuilderFrom(seed).Build()
			assert.Equal(t, got.Level, seed.Level, "NewTierBuilderFrom keeps Level")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewTierBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewTierBuilderFrom(want).Build(),
				want,
				"NewTierBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithLevel", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Level", func(t *testing.T) {
			t.Parallel()

			var want int
			want = 42
			got := NewTierBuilder().WithLevel(want).Build()
			assert.Equal(t, got.Level, want, "WithLevel writes its argument to Level")
		})

		t.Run("replaces the value of Level", func(t *testing.T) {
			t.Parallel()

			var first, want int
			first = 42
			want = 7
			got := NewTierBuilder().WithLevel(first).WithLevel(want).Build()
			assert.Equal(t, got.Level, want, "the second call of WithLevel replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Tier
			record := func(v *Tier) { got = append(got, v) }
			NewTierBuilder().Mutate(record).Mutate(record)
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

			original := NewTierBuilder()
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

			b := NewTierBuilder()
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

			var seed Tier
			seed.Level = 42
			b := NewTierBuilder()
			b.Mutate(func(v *Tier) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// TestDeploymentBuilder checks the constructors and the methods of
// DeploymentBuilder.
func TestDeploymentBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewDeploymentBuilder", func(t *testing.T) {
		t.Parallel()

		t.Run("applies the declared default of Mode", func(t *testing.T) {
			t.Parallel()

			var want Code
			want = "true"
			got := NewDeploymentBuilder().Build()
			assert.Equal(
				t,
				got.Mode,
				want,
				"NewDeploymentBuilder applies the declared default of Mode",
			)
		})

		t.Run("applies the declared default of Region", func(t *testing.T) {
			t.Parallel()

			var want string
			want = Region
			got := NewDeploymentBuilder().Build()
			assert.Equal(
				t,
				got.Region,
				want,
				"NewDeploymentBuilder applies the declared default of Region",
			)
		})

		t.Run("applies the declared default of Zone", func(t *testing.T) {
			t.Parallel()

			var want string
			want = zone()
			got := NewDeploymentBuilder().Build()
			assert.Equal(
				t,
				got.Zone,
				want,
				"NewDeploymentBuilder applies the declared default of Zone",
			)
		})

		t.Run("applies the declared default of System", func(t *testing.T) {
			t.Parallel()

			var want string
			want = runtime.GOOS
			got := NewDeploymentBuilder().Build()
			assert.Equal(
				t,
				got.System,
				want,
				"NewDeploymentBuilder applies the declared default of System",
			)
		})

		t.Run("applies the declared default of Name", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "web"
			got := NewDeploymentBuilder().Build()
			assert.Equal(
				t,
				got.Name,
				want,
				"NewDeploymentBuilder applies the declared default of Name",
			)
		})

		t.Run("applies the declared default of Label", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "a b"
			got := NewDeploymentBuilder().Build()
			assert.Equal(
				t,
				got.Label,
				want,
				"NewDeploymentBuilder applies the declared default of Label",
			)
		})

		t.Run("applies the declared default of Flag", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "true"
			got := NewDeploymentBuilder().Build()
			assert.Equal(
				t,
				got.Flag,
				want,
				"NewDeploymentBuilder applies the declared default of Flag",
			)
		})
	})

	t.Run("NewDeploymentBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
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
			expect.Equal(t, got.Mode, seed.Mode, "NewDeploymentBuilderFrom keeps Mode")
			expect.Equal(t, got.Region, seed.Region, "NewDeploymentBuilderFrom keeps Region")
			expect.Equal(t, got.Zone, seed.Zone, "NewDeploymentBuilderFrom keeps Zone")
			expect.Equal(t, got.System, seed.System, "NewDeploymentBuilderFrom keeps System")
			expect.Equal(t, got.Name, seed.Name, "NewDeploymentBuilderFrom keeps Name")
			expect.Equal(t, got.Label, seed.Label, "NewDeploymentBuilderFrom keeps Label")
			expect.Equal(t, got.Flag, seed.Flag, "NewDeploymentBuilderFrom keeps Flag")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := NewDeploymentBuilder()
			want := b.Build()
			assert.Equal(
				t,
				NewDeploymentBuilderFrom(want).Build(),
				want,
				"NewDeploymentBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithMode", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Mode", func(t *testing.T) {
			t.Parallel()

			var want Code
			want = Code("Mode-a")
			got := NewDeploymentBuilder().WithMode(want).Build()
			assert.Equal(t, got.Mode, want, "WithMode writes its argument to Mode")
		})

		t.Run("replaces the value of Mode", func(t *testing.T) {
			t.Parallel()

			var first, want Code
			first = Code("Mode-a")
			want = Code("Mode-b")
			got := NewDeploymentBuilder().WithMode(first).WithMode(want).Build()
			assert.Equal(t, got.Mode, want, "the second call of WithMode replaces the first")
		})
	})

	t.Run("WithRegion", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Region", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Region-a"
			got := NewDeploymentBuilder().WithRegion(want).Build()
			assert.Equal(t, got.Region, want, "WithRegion writes its argument to Region")
		})

		t.Run("replaces the value of Region", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Region-a"
			want = "Region-b"
			got := NewDeploymentBuilder().WithRegion(first).WithRegion(want).Build()
			assert.Equal(t, got.Region, want, "the second call of WithRegion replaces the first")
		})
	})

	t.Run("WithZone", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Zone", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Zone-a"
			got := NewDeploymentBuilder().WithZone(want).Build()
			assert.Equal(t, got.Zone, want, "WithZone writes its argument to Zone")
		})

		t.Run("replaces the value of Zone", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Zone-a"
			want = "Zone-b"
			got := NewDeploymentBuilder().WithZone(first).WithZone(want).Build()
			assert.Equal(t, got.Zone, want, "the second call of WithZone replaces the first")
		})
	})

	t.Run("WithSystem", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to System", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "System-a"
			got := NewDeploymentBuilder().WithSystem(want).Build()
			assert.Equal(t, got.System, want, "WithSystem writes its argument to System")
		})

		t.Run("replaces the value of System", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "System-a"
			want = "System-b"
			got := NewDeploymentBuilder().WithSystem(first).WithSystem(want).Build()
			assert.Equal(t, got.System, want, "the second call of WithSystem replaces the first")
		})
	})

	t.Run("WithName", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Name", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Name-a"
			got := NewDeploymentBuilder().WithName(want).Build()
			assert.Equal(t, got.Name, want, "WithName writes its argument to Name")
		})

		t.Run("replaces the value of Name", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Name-a"
			want = "Name-b"
			got := NewDeploymentBuilder().WithName(first).WithName(want).Build()
			assert.Equal(t, got.Name, want, "the second call of WithName replaces the first")
		})
	})

	t.Run("WithLabel", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Label", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Label-a"
			got := NewDeploymentBuilder().WithLabel(want).Build()
			assert.Equal(t, got.Label, want, "WithLabel writes its argument to Label")
		})

		t.Run("replaces the value of Label", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Label-a"
			want = "Label-b"
			got := NewDeploymentBuilder().WithLabel(first).WithLabel(want).Build()
			assert.Equal(t, got.Label, want, "the second call of WithLabel replaces the first")
		})
	})

	t.Run("WithFlag", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Flag", func(t *testing.T) {
			t.Parallel()

			var want string
			want = "Flag-a"
			got := NewDeploymentBuilder().WithFlag(want).Build()
			assert.Equal(t, got.Flag, want, "WithFlag writes its argument to Flag")
		})

		t.Run("replaces the value of Flag", func(t *testing.T) {
			t.Parallel()

			var first, want string
			first = "Flag-a"
			want = "Flag-b"
			got := NewDeploymentBuilder().WithFlag(first).WithFlag(want).Build()
			assert.Equal(t, got.Flag, want, "the second call of WithFlag replaces the first")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*Deployment
			record := func(v *Deployment) { got = append(got, v) }
			NewDeploymentBuilder().Mutate(record).Mutate(record)
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

			original := NewDeploymentBuilder()
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

			b := NewDeploymentBuilder()
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

			var seed Deployment
			seed.Mode = Code("Mode-a")
			seed.Region = "Region-a"
			seed.Zone = "Zone-a"
			seed.System = "System-a"
			seed.Name = "Name-a"
			seed.Label = "Label-a"
			seed.Flag = "Flag-a"
			b := NewDeploymentBuilder()
			b.Mutate(func(v *Deployment) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}
