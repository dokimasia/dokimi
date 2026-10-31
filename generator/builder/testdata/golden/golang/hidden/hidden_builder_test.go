package hidden_test

import (
	"testing"

	"example.com/shop/hidden"
	"go.dokimi.dev/assert"
)

// TestSiteBuilder checks the constructors and the methods of SiteBuilder.
func TestSiteBuilder(t *testing.T) {
	t.Parallel()

	t.Run("NewSiteBuilderFrom", func(t *testing.T) {
		t.Parallel()

		t.Run("keeps each member of its value", func(t *testing.T) {
			t.Parallel()

			var seed hidden.Site
			seed.Zone = hidden.NewZone()
			got := hidden.NewSiteBuilderFrom(seed).Build()
			assert.Equal(t, got.Zone, seed.Zone, "NewSiteBuilderFrom keeps Zone")
		})

		t.Run("keeps the value of its argument", func(t *testing.T) {
			t.Parallel()

			b := hidden.NewSiteBuilder()
			b.WithHosts(hidden.Host{})
			want := b.Build()
			assert.Equal(
				t,
				hidden.NewSiteBuilderFrom(want).Build(),
				want,
				"NewSiteBuilderFrom keeps the value of its argument",
			)
		})
	})

	t.Run("WithZone", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Zone", func(t *testing.T) {
			t.Parallel()

			var want hidden.Zone
			want = hidden.NewZone()
			got := hidden.NewSiteBuilder().WithZone(want).Build()
			assert.Equal(t, got.Zone, want, "WithZone writes its argument to Zone")
		})
	})

	t.Run("WithHost", func(t *testing.T) {
		t.Parallel()

		t.Run("returns its builder", func(t *testing.T) {
			t.Parallel()

			b := hidden.NewSiteBuilder()
			assert.Equal(
				t,
				b.WithHost(b.Build().Host),
				b,
				"WithHost returns its builder for the next call",
				assert.ByIdentity(),
			)
		})
	})

	t.Run("WithHosts", func(t *testing.T) {
		t.Parallel()

		t.Run("writes its argument to Hosts", func(t *testing.T) {
			t.Parallel()

			got := hidden.NewSiteBuilder().WithHosts(hidden.Host{}).Build()
			assert.Equal(t, got.Hosts, []hidden.Host{{}}, "WithHosts writes its argument to Hosts")
		})
	})

	t.Run("AppendHosts", func(t *testing.T) {
		t.Parallel()

		t.Run("appends to the elements of Hosts", func(t *testing.T) {
			t.Parallel()

			got := hidden.NewSiteBuilder().
				WithHosts(hidden.Host{}).
				AppendHosts(hidden.Host{}).
				Build()
			assert.Equal(t, got.Hosts, []hidden.Host{{}, {}},
				"AppendHosts appends its elements after those of WithHosts")
		})
	})

	t.Run("Mutate", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the value under construction", func(t *testing.T) {
			t.Parallel()

			var got []*hidden.Site
			record := func(v *hidden.Site) { got = append(got, v) }
			hidden.NewSiteBuilder().Mutate(record).Mutate(record)
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

			original := hidden.NewSiteBuilder()
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

			b := hidden.NewSiteBuilder()
			b.WithHosts(hidden.Host{})
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

			var seed hidden.Site
			seed.Zone = hidden.NewZone()
			b := hidden.NewSiteBuilder()
			b.Mutate(func(v *hidden.Site) { *v = seed })
			assert.Equal(t, b.Build(), seed, "Build returns the value that Mutate sets")
		})
	})
}

// [Pool] is generic, and its type parameters have no witnesses that a test
// outside its package can refer to, so the file has no checks of
// PoolBuilder. A test in the package of a consumer instantiates it.
