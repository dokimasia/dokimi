package scoped_test

import (
	testing2 "testing"

	"example.com/shop/scoped"
	assert2 "go.dokimi.dev/assert"
	conform2 "go.dokimi.dev/conform"
)

// TestHyphenateSuite checks the suite of Hyphenate and the two forms of
// its application.
func TestHyphenateSuite(t *testing2.T) {
	t.Parallel()

	t.Run("HyphenateSuite", func(t *testing2.T) {
		t.Parallel()
		conform2.Falsify(t, scoped.HyphenateSuite())
	})

	t.Run("Hyphenate", func(t *testing2.T) {
		t.Parallel()
		scoped.RunHyphenateSuite(t, scoped.Hyphenate)
	})

	t.Run("RunHyphenateSuite", func(t *testing2.T) {
		t.Parallel()
		scoped.RunHyphenateSuite(
			t,
			scoped.Hyphenate,
			conform2.With(scoped.HyphenateCheck{
				ID:        "Hyphenate/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *scoped.HyphenateHandle) {
					assert2.Equal(
						h,
						h.Subject(),
						scoped.HyphenateFunc(scoped.Hyphenate),
						h.Statement(),
					)
				},
				Defect: func(*scoped.HyphenateHandle) scoped.HyphenateFunc { return nil },
			}),
			conform2.ExcuseSuite[scoped.HyphenateFunc](
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunHyphenateSuiteFromHandle", func(t *testing2.T) {
		t.Parallel()
		scoped.RunHyphenateSuiteFromHandle(
			t,
			func(*scoped.HyphenateHandle) scoped.HyphenateFunc {
				return scoped.Hyphenate
			},
			conform2.With(scoped.HyphenateCheck{
				ID:        "Hyphenate/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *scoped.HyphenateHandle) {
					assert2.Equal(
						h,
						h.Subject(),
						scoped.HyphenateFunc(scoped.Hyphenate),
						h.Statement(),
					)
				},
				Defect: func(*scoped.HyphenateHandle) scoped.HyphenateFunc { return nil },
			}),
			conform2.ExcuseSuite[scoped.HyphenateFunc](
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("HyphenateChecks", func(t *testing2.T) {
		t.Parallel()
		assert2.Equal(
			t,
			string(scoped.HyphenateChecks.LatencyBudget()),
			"Hyphenate/latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert2.Equal(
			t,
			string(scoped.HyphenateChecks.MeanLatencyBudget()),
			"Hyphenate/mean-latency-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// BenchmarkHyphenate measures Hyphenate.
func BenchmarkHyphenate(b *testing2.B) {
	scoped.BenchHyphenateSuite(b, scoped.Hyphenate)
}

// FuzzHyphenate fuzzes the property checks of Hyphenate.
func FuzzHyphenate(f *testing2.F) {
	scoped.FuzzHyphenateSuite(f, scoped.Hyphenate)
}
