package slug

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// TestSlugSuite checks the suite of Slug and the two forms of its
// application.
func TestSlugSuite(t *testing.T) {
	t.Parallel()

	t.Run("SlugSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, SlugSuite())
	})

	t.Run("Slug", func(t *testing.T) {
		t.Parallel()
		RunSlugSuite(t, Slug)
	})

	t.Run("RunSlugSuite", func(t *testing.T) {
		t.Parallel()
		RunSlugSuite(
			t,
			Slug,
			conform.With(SlugCheck{
				ID:        "Slug/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *SlugHandle) {
					assert.Equal(h, h.Subject(), SlugFunc(Slug), h.Statement())
				},
				Defect: func(*SlugHandle) SlugFunc { return nil },
			}),
			conform.Excuse(
				SlugChecks.KeepsTheOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SlugChecks.KeepsTheWords(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSlugSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		RunSlugSuiteFromHandle(
			t,
			func(*SlugHandle) SlugFunc {
				return Slug
			},
			conform.With(SlugCheck{
				ID:        "Slug/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *SlugHandle) {
					assert.Equal(h, h.Subject(), SlugFunc(Slug), h.Statement())
				},
				Defect: func(*SlugHandle) SlugFunc { return nil },
			}),
			conform.Excuse(
				SlugChecks.KeepsTheOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				SlugChecks.KeepsTheWords(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("SlugChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(SlugChecks.KeepsTheOrder()),
			"Slug/keeps-the-order",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(SlugChecks.KeepsTheWords()),
			"Slug/keeps-the-words",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// BenchmarkSlug measures Slug.
func BenchmarkSlug(b *testing.B) {
	BenchSlugSuite(b, Slug)
}

// FuzzSlug fuzzes the property checks of Slug.
func FuzzSlug(f *testing.F) {
	FuzzSlugSuite(f, Slug)
}
