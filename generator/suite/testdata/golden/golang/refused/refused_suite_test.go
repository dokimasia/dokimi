package refused_test

import (
	"testing"

	"example.com/shop/refused"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// hybridClockSubjectEmbedded is refused.HybridClock under the name of the
// field that hybridClockSubject embeds.
type hybridClockSubjectEmbedded = refused.HybridClock

// hybridClockSubject is a value of HybridClock that only
// TestHybridClockSuite builds.
type hybridClockSubject struct{ hybridClockSubjectEmbedded }

// TestHybridClockSuite checks the suite of HybridClock and the two forms
// of its application.
func TestHybridClockSuite(t *testing.T) {
	t.Parallel()

	t.Run("HybridClockSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, refused.HybridClockSuite())
	})

	t.Run("RunHybridClockSuite", func(t *testing.T) {
		t.Parallel()
		var built *hybridClockSubject
		refused.RunHybridClockSuite(
			t,
			func() *hybridClockSubject {
				built = &hybridClockSubject{}
				return built
			},
			conform.With(refused.HybridClockCheck{
				ID:        "HybridClock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *refused.HybridClockHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						refused.HybridClock(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*refused.HybridClockHandle) refused.HybridClock { return &hybridClockSubject{} },
			}),
			conform.Excuse(
				refused.HybridClockChecks.Update.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.Update.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.MaxOffset.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.MaxOffset.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.MaxOffset.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunHybridClockSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *hybridClockSubject
		refused.RunHybridClockSuiteFromHandle(
			t,
			func(*refused.HybridClockHandle) *hybridClockSubject {
				built = &hybridClockSubject{}
				return built
			},
			conform.With(refused.HybridClockCheck{
				ID:        "HybridClock/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *refused.HybridClockHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						refused.HybridClock(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*refused.HybridClockHandle) refused.HybridClock { return &hybridClockSubject{} },
			}),
			conform.Excuse(
				refused.HybridClockChecks.Update.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.Update.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.MaxOffset.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.MaxOffset.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.HybridClockChecks.MaxOffset.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// hybridClockIdle returns the planted defect of a ceiling of HybridClock
// with a spend of nothing, a subject that answers each call at once.
func hybridClockIdle(h *refused.HybridClockHandle) refused.HybridClock {
	return refused.HybridClockSuite().Calls[0].Over(h, func() {})
}

// TestHybridClockSuiteBudgets applies the allocation ceilings of the suite
// of HybridClock to hybridClockIdle with every check of the suite excused,
// and a ceiling of 0 to each callable without a budget whose calls build
// no input of their own.
func TestHybridClockSuiteBudgets(t *testing.T) {
	refused.RunHybridClockBudgets(
		t,
		hybridClockIdle,
		conform.Ceilings(
			refused.HybridClockChecks.Update.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			refused.HybridClockChecks.MaxOffset.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			refused.HybridClockChecks.Update.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.Update.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHybridClockSuiteUpdate registers the fuzz target of
// HybridClock.Update on hybridClockIdle, with every check of the suite
// excused.
func FuzzHybridClockSuiteUpdate(f *testing.F) {
	refused.FuzzHybridClockSuite(
		f,
		refused.HybridClockChecks.Update.Callable,
		hybridClockIdle,
		conform.Excuse(
			refused.HybridClockChecks.Update.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.Update.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHybridClockSuiteMaxOffset registers the fuzz target of
// HybridClock.MaxOffset on hybridClockIdle, with every check of the suite
// excused.
func FuzzHybridClockSuiteMaxOffset(f *testing.F) {
	refused.FuzzHybridClockSuite(
		f,
		refused.HybridClockChecks.MaxOffset.Callable,
		hybridClockIdle,
		conform.Excuse(
			refused.HybridClockChecks.Update.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.Update.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.HybridClockChecks.MaxOffset.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkHybridClockSuite measures each callable of the suite of
// HybridClock on hybridClockIdle.
func BenchmarkHybridClockSuite(b *testing.B) {
	b.Run("Update", func(b *testing.B) {
		refused.BenchHybridClockSuite(b, refused.HybridClockChecks.Update.Callable, hybridClockIdle)
	})
	b.Run("MaxOffset", func(b *testing.B) {
		refused.BenchHybridClockSuite(
			b,
			refused.HybridClockChecks.MaxOffset.Callable,
			hybridClockIdle,
		)
	})
}

// doubledSubjectEmbedded is refused.Doubled under the name of the field
// that doubledSubject embeds.
type doubledSubjectEmbedded = refused.Doubled

// doubledSubject is a value of Doubled that only TestDoubledSuite builds.
type doubledSubject struct{ doubledSubjectEmbedded }

// TestDoubledSuite checks the suite of Doubled and the two forms of its
// application.
func TestDoubledSuite(t *testing.T) {
	t.Parallel()

	t.Run("DoubledSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, refused.DoubledSuite())
	})

	t.Run("RunDoubledSuite", func(t *testing.T) {
		t.Parallel()
		var built *doubledSubject
		refused.RunDoubledSuite(
			t,
			func() *doubledSubject {
				built = &doubledSubject{}
				return built
			},
			conform.With(refused.DoubledCheck{
				ID:        "Doubled/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *refused.DoubledHandle) {
					got := h.Subject()
					assert.Equal(h, got, refused.Doubled(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*refused.DoubledHandle) refused.Doubled { return &doubledSubject{} },
			}),
			conform.Excuse(
				refused.DoubledChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.DoubledChecks.Read.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.DoubledChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.DoubledChecks.Read.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunDoubledSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *doubledSubject
		refused.RunDoubledSuiteFromHandle(
			t,
			func(*refused.DoubledHandle) *doubledSubject {
				built = &doubledSubject{}
				return built
			},
			conform.With(refused.DoubledCheck{
				ID:        "Doubled/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *refused.DoubledHandle) {
					got := h.Subject()
					assert.Equal(h, got, refused.Doubled(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*refused.DoubledHandle) refused.Doubled { return &doubledSubject{} },
			}),
			conform.Excuse(
				refused.DoubledChecks.Read.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.DoubledChecks.Read.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.DoubledChecks.Read.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.DoubledChecks.Read.CancelGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// doubledIdle returns the planted defect of a ceiling of Doubled with a
// spend of nothing, a subject that answers each call at once.
func doubledIdle(h *refused.DoubledHandle) refused.Doubled {
	return refused.DoubledSuite().Calls[0].Over(h, func() {})
}

// TestDoubledSuiteBudgets applies the allocation ceilings of the suite of
// Doubled to doubledIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestDoubledSuiteBudgets(t *testing.T) {
	refused.RunDoubledBudgets(
		t,
		doubledIdle,
		conform.Ceilings(
			refused.DoubledChecks.Read.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			refused.DoubledChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.DoubledChecks.Read.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.DoubledChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.DoubledChecks.Read.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDoubledSuiteRead registers the fuzz target of Doubled.Read on
// doubledIdle, with every check of the suite excused.
func FuzzDoubledSuiteRead(f *testing.F) {
	refused.FuzzDoubledSuite(
		f,
		refused.DoubledChecks.Read.Callable,
		doubledIdle,
		conform.Excuse(
			refused.DoubledChecks.Read.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.DoubledChecks.Read.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.DoubledChecks.Read.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			refused.DoubledChecks.Read.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkDoubledSuite measures each callable of the suite of Doubled on
// doubledIdle.
func BenchmarkDoubledSuite(b *testing.B) {
	b.Run("Read", func(b *testing.B) {
		refused.BenchDoubledSuite(b, refused.DoubledChecks.Read.Callable, doubledIdle)
	})
}

// TestTitleSuite checks the suite of Title and the two forms of its
// application.
func TestTitleSuite(t *testing.T) {
	t.Parallel()

	t.Run("TitleSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, refused.TitleSuite())
	})

	t.Run("Title", func(t *testing.T) {
		t.Parallel()
		refused.RunTitleSuite(t, refused.Title)
	})

	t.Run("RunTitleSuite", func(t *testing.T) {
		t.Parallel()
		refused.RunTitleSuite(
			t,
			refused.Title,
			conform.With(refused.TitleCheck{
				ID:        "Title/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *refused.TitleHandle) {
					assert.Equal(h, h.Subject(), refused.TitleFunc(refused.Title), h.Statement())
				},
				Defect: func(*refused.TitleHandle) refused.TitleFunc { return nil },
			}),
			conform.Excuse(
				refused.TitleChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.TitleChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunTitleSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		refused.RunTitleSuiteFromHandle(
			t,
			func(*refused.TitleHandle) refused.TitleFunc {
				return refused.Title
			},
			conform.With(refused.TitleCheck{
				ID:        "Title/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *refused.TitleHandle) {
					assert.Equal(h, h.Subject(), refused.TitleFunc(refused.Title), h.Statement())
				},
				Defect: func(*refused.TitleHandle) refused.TitleFunc { return nil },
			}),
			conform.Excuse(
				refused.TitleChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				refused.TitleChecks.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// BenchmarkTitle measures Title.
func BenchmarkTitle(b *testing.B) {
	refused.BenchTitleSuite(b, refused.Title)
}

// FuzzTitle fuzzes the property checks of Title.
func FuzzTitle(f *testing.F) {
	refused.FuzzTitleSuite(f, refused.Title)
}
