package slug_test

import (
	"testing"

	"example.com/shop/slug"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// TestSlugSuite checks the suite of Slug and the two forms of its
// application.
func TestSlugSuite(t *testing.T) {
	t.Parallel()

	t.Run("SlugSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, slug.SlugSuite())
	})

	t.Run("Slug", func(t *testing.T) {
		t.Parallel()
		slug.RunSlugSuite(t, slug.Slug)
	})

	t.Run("RunSlugSuite", func(t *testing.T) {
		t.Parallel()
		slug.RunSlugSuite(
			t,
			slug.Slug,
			conform.With(slug.SlugCheck{
				ID:        "Slug/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *slug.SlugHandle) {
					assert.Equal(h, h.Subject(), slug.SlugFunc(slug.Slug), h.Statement())
				},
				Defect: func(*slug.SlugHandle) slug.SlugFunc { return nil },
			}),
			conform.Excuse(
				slug.SlugChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				slug.SlugChecks.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				slug.SlugChecks.KeepsTheOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				slug.SlugChecks.KeepsTheWords(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSlugSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		slug.RunSlugSuiteFromHandle(
			t,
			func(*slug.SlugHandle) slug.SlugFunc {
				return slug.Slug
			},
			conform.With(slug.SlugCheck{
				ID:        "Slug/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *slug.SlugHandle) {
					assert.Equal(h, h.Subject(), slug.SlugFunc(slug.Slug), h.Statement())
				},
				Defect: func(*slug.SlugHandle) slug.SlugFunc { return nil },
			}),
			conform.Excuse(
				slug.SlugChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				slug.SlugChecks.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				slug.SlugChecks.KeepsTheOrder(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				slug.SlugChecks.KeepsTheWords(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("SlugChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(slug.SlugChecks.KeepsTheOrder()),
			"Slug/keeps-the-order",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(slug.SlugChecks.KeepsTheWords()),
			"Slug/keeps-the-words",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// BenchmarkSlug measures Slug.
func BenchmarkSlug(b *testing.B) {
	slug.BenchSlugSuite(b, slug.Slug)
}

// FuzzSlug fuzzes the property checks of Slug.
func FuzzSlug(f *testing.F) {
	slug.FuzzSlugSuite(f, slug.Slug)
}
