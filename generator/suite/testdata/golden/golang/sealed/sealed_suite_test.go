package sealed_test

import (
	"testing"

	"example.com/shop/sealed"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// aeadSubjectEmbedded is sealed.AEAD under the name of the field that
// aeadSubject embeds.
type aeadSubjectEmbedded = sealed.AEAD

// aeadSubject is a value of AEAD that only TestAEADSuite builds.
type aeadSubject struct{ aeadSubjectEmbedded }

// TestAEADSuite checks the suite of AEAD and the two forms of its
// application.
func TestAEADSuite(t *testing.T) {
	t.Parallel()

	t.Run("AEADSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sealed.AEADSuite())
	})

	t.Run("RunAEADSuite", func(t *testing.T) {
		t.Parallel()
		var built *aeadSubject
		sealed.RunAEADSuite(
			t,
			func() *aeadSubject {
				built = &aeadSubject{}
				return built
			},
			conform.With(sealed.AEADCheck{
				ID:        "AEAD/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.AEADHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.AEAD(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.AEADHandle) sealed.AEAD { return &aeadSubject{} },
			}),
			conform.Excuse(
				sealed.AEADChecks.Seal.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Seal.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.NonceSize.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.NonceSize.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Overhead.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Overhead.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunAEADSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *aeadSubject
		sealed.RunAEADSuiteFromHandle(
			t,
			func(*sealed.AEADHandle) *aeadSubject {
				built = &aeadSubject{}
				return built
			},
			conform.With(sealed.AEADCheck{
				ID:        "AEAD/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.AEADHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.AEAD(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.AEADHandle) sealed.AEAD { return &aeadSubject{} },
			}),
			conform.Excuse(
				sealed.AEADChecks.Seal.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Seal.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Open.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.NonceSize.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.NonceSize.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Overhead.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.AEADChecks.Overhead.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// aeadIdle returns the planted defect of a ceiling of AEAD with a spend of
// nothing, a subject that answers each call at once.
func aeadIdle(h *sealed.AEADHandle) sealed.AEAD {
	return sealed.AEADSuite().Calls[0].Over(h, func() {})
}

// TestAEADSuiteBudgets applies the allocation ceilings of the suite of
// AEAD to aeadIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestAEADSuiteBudgets(t *testing.T) {
	sealed.RunAEADBudgets(
		t,
		aeadIdle,
		conform.Ceilings(
			sealed.AEADChecks.NonceSize.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			sealed.AEADChecks.Overhead.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			sealed.AEADChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Seal.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAEADSuiteSeal registers the fuzz target of AEAD.Seal on aeadIdle,
// with every check of the suite excused.
func FuzzAEADSuiteSeal(f *testing.F) {
	sealed.FuzzAEADSuite(
		f,
		sealed.AEADChecks.Seal.Callable,
		aeadIdle,
		conform.Excuse(
			sealed.AEADChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Seal.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAEADSuiteOpen registers the fuzz target of AEAD.Open on aeadIdle,
// with every check of the suite excused.
func FuzzAEADSuiteOpen(f *testing.F) {
	sealed.FuzzAEADSuite(
		f,
		sealed.AEADChecks.Open.Callable,
		aeadIdle,
		conform.Excuse(
			sealed.AEADChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Seal.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAEADSuiteNonceSize registers the fuzz target of AEAD.NonceSize on
// aeadIdle, with every check of the suite excused.
func FuzzAEADSuiteNonceSize(f *testing.F) {
	sealed.FuzzAEADSuite(
		f,
		sealed.AEADChecks.NonceSize.Callable,
		aeadIdle,
		conform.Excuse(
			sealed.AEADChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Seal.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzAEADSuiteOverhead registers the fuzz target of AEAD.Overhead on
// aeadIdle, with every check of the suite excused.
func FuzzAEADSuiteOverhead(f *testing.F) {
	sealed.FuzzAEADSuite(
		f,
		sealed.AEADChecks.Overhead.Callable,
		aeadIdle,
		conform.Excuse(
			sealed.AEADChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Seal.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Open.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.NonceSize.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.AEADChecks.Overhead.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkAEADSuite measures each callable of the suite of AEAD on
// aeadIdle.
func BenchmarkAEADSuite(b *testing.B) {
	b.Run("Seal", func(b *testing.B) {
		sealed.BenchAEADSuite(b, sealed.AEADChecks.Seal.Callable, aeadIdle)
	})
	b.Run("Open", func(b *testing.B) {
		sealed.BenchAEADSuite(b, sealed.AEADChecks.Open.Callable, aeadIdle)
	})
	b.Run("NonceSize", func(b *testing.B) {
		sealed.BenchAEADSuite(b, sealed.AEADChecks.NonceSize.Callable, aeadIdle)
	})
	b.Run("Overhead", func(b *testing.B) {
		sealed.BenchAEADSuite(b, sealed.AEADChecks.Overhead.Callable, aeadIdle)
	})
}

// metronomeSubjectEmbedded is sealed.Metronome under the name of the field
// that metronomeSubject embeds.
type metronomeSubjectEmbedded = sealed.Metronome

// metronomeSubject is a value of Metronome that only TestMetronomeSuite
// builds.
type metronomeSubject struct{ metronomeSubjectEmbedded }

// TestMetronomeSuite checks the suite of Metronome and the two forms of
// its application.
func TestMetronomeSuite(t *testing.T) {
	t.Parallel()

	t.Run("MetronomeSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sealed.MetronomeSuite())
	})

	t.Run("RunMetronomeSuite", func(t *testing.T) {
		t.Parallel()
		var built *metronomeSubject
		sealed.RunMetronomeSuite(
			t,
			func() *metronomeSubject {
				built = &metronomeSubject{}
				return built
			},
			conform.With(sealed.MetronomeCheck{
				ID:        "Metronome/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.MetronomeHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sealed.Metronome(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sealed.MetronomeHandle) sealed.Metronome { return &metronomeSubject{} },
			}),
			conform.Excuse(
				sealed.MetronomeChecks.Tick.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Tick.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Tick.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Pace.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Pace.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Pace.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunMetronomeSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *metronomeSubject
		sealed.RunMetronomeSuiteFromHandle(
			t,
			func(*sealed.MetronomeHandle) *metronomeSubject {
				built = &metronomeSubject{}
				return built
			},
			conform.With(sealed.MetronomeCheck{
				ID:        "Metronome/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.MetronomeHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						sealed.Metronome(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*sealed.MetronomeHandle) sealed.Metronome { return &metronomeSubject{} },
			}),
			conform.Excuse(
				sealed.MetronomeChecks.Tick.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Tick.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Tick.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Pace.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Pace.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MetronomeChecks.Pace.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// metronomeIdle returns the planted defect of a ceiling of Metronome with
// a spend of nothing, a subject that answers each call at once.
func metronomeIdle(h *sealed.MetronomeHandle) sealed.Metronome {
	return sealed.MetronomeSuite().Calls[0].Over(h, func() {})
}

// TestMetronomeSuiteBudgets applies the allocation ceilings of the suite
// of Metronome to metronomeIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestMetronomeSuiteBudgets(t *testing.T) {
	sealed.RunMetronomeBudgets(
		t,
		metronomeIdle,
		conform.Excuse(
			sealed.MetronomeChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Tick.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMetronomeSuiteTick registers the fuzz target of Metronome.Tick on
// metronomeIdle, with every check of the suite excused.
func FuzzMetronomeSuiteTick(f *testing.F) {
	sealed.FuzzMetronomeSuite(
		f,
		sealed.MetronomeChecks.Tick.Callable,
		metronomeIdle,
		conform.Excuse(
			sealed.MetronomeChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Tick.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMetronomeSuitePace registers the fuzz target of Metronome.Pace on
// metronomeIdle, with every check of the suite excused.
func FuzzMetronomeSuitePace(f *testing.F) {
	sealed.FuzzMetronomeSuite(
		f,
		sealed.MetronomeChecks.Pace.Callable,
		metronomeIdle,
		conform.Excuse(
			sealed.MetronomeChecks.Tick.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Tick.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Tick.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MetronomeChecks.Pace.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkMetronomeSuite measures each callable of the suite of Metronome
// on metronomeIdle.
func BenchmarkMetronomeSuite(b *testing.B) {
	b.Run("Tick", func(b *testing.B) {
		sealed.BenchMetronomeSuite(b, sealed.MetronomeChecks.Tick.Callable, metronomeIdle)
	})
	b.Run("Pace", func(b *testing.B) {
		sealed.BenchMetronomeSuite(b, sealed.MetronomeChecks.Pace.Callable, metronomeIdle)
	})
}

// lockerSubjectEmbedded is sealed.Locker under the name of the field that
// lockerSubject embeds.
type lockerSubjectEmbedded = sealed.Locker

// lockerSubject is a value of Locker that only TestLockerSuite builds.
type lockerSubject struct{ lockerSubjectEmbedded }

// TestLockerSuite checks the suite of Locker and the two forms of its
// application.
func TestLockerSuite(t *testing.T) {
	t.Parallel()

	t.Run("LockerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sealed.LockerSuite())
	})

	t.Run("RunLockerSuite", func(t *testing.T) {
		t.Parallel()
		var built *lockerSubject
		sealed.RunLockerSuite(
			t,
			func() *lockerSubject {
				built = &lockerSubject{}
				return built
			},
			conform.With(sealed.LockerCheck{
				ID:        "Locker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.LockerHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Locker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.LockerHandle) sealed.Locker { return &lockerSubject{} },
			}),
			conform.Excuse(
				sealed.LockerChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.SyncedWritesSurvive(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.NoTornRecord(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.FullMediumIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunLockerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *lockerSubject
		sealed.RunLockerSuiteFromHandle(
			t,
			func(*sealed.LockerHandle) *lockerSubject {
				built = &lockerSubject{}
				return built
			},
			conform.With(sealed.LockerCheck{
				ID:        "Locker/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.LockerHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Locker(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.LockerHandle) sealed.Locker { return &lockerSubject{} },
			}),
			conform.Excuse(
				sealed.LockerChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.SyncedWritesSurvive(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.NoTornRecord(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Put.FullMediumIsReported(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.LockerChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// lockerIdle returns the planted defect of a ceiling of Locker with a
// spend of nothing, a subject that answers each call at once.
func lockerIdle(h *sealed.LockerHandle) sealed.Locker {
	return sealed.LockerSuite().Calls[0].Over(h, func() {})
}

// TestLockerSuiteBudgets applies the allocation ceilings of the suite of
// Locker to lockerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLockerSuiteBudgets(t *testing.T) {
	sealed.RunLockerBudgets(
		t,
		lockerIdle,
		conform.Ceilings(sealed.LockerChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			sealed.LockerChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.SyncedWritesSurvive(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.NoTornRecord(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.FullMediumIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockerSuitePut registers the fuzz target of Locker.Put on
// lockerIdle, with every check of the suite excused.
func FuzzLockerSuitePut(f *testing.F) {
	sealed.FuzzLockerSuite(
		f,
		sealed.LockerChecks.Put.Callable,
		lockerIdle,
		conform.Excuse(
			sealed.LockerChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.SyncedWritesSurvive(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.NoTornRecord(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.FullMediumIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzLockerSuiteGet registers the fuzz target of Locker.Get on
// lockerIdle, with every check of the suite excused.
func FuzzLockerSuiteGet(f *testing.F) {
	sealed.FuzzLockerSuite(
		f,
		sealed.LockerChecks.Get.Callable,
		lockerIdle,
		conform.Excuse(
			sealed.LockerChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.SyncedWritesSurvive(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.NoTornRecord(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Put.FullMediumIsReported(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.LockerChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkLockerSuite measures each callable of the suite of Locker on
// lockerIdle.
func BenchmarkLockerSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		sealed.BenchLockerSuite(b, sealed.LockerChecks.Put.Callable, lockerIdle)
	})
	b.Run("Get", func(b *testing.B) {
		sealed.BenchLockerSuite(b, sealed.LockerChecks.Get.Callable, lockerIdle)
	})
}

// shelfSubjectEmbedded is sealed.Shelf under the name of the field that
// shelfSubject embeds.
type shelfSubjectEmbedded = sealed.Shelf

// shelfSubject is a value of Shelf that only TestShelfSuite builds.
type shelfSubject struct{ shelfSubjectEmbedded }

// TestShelfSuite checks the suite of Shelf and the two forms of its
// application.
func TestShelfSuite(t *testing.T) {
	t.Parallel()

	t.Run("ShelfSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sealed.ShelfSuite())
	})

	t.Run("RunShelfSuite", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		sealed.RunShelfSuite(
			t,
			func() *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(sealed.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.ShelfHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Shelf(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.ShelfHandle) sealed.Shelf { return &shelfSubject{} },
			}),
			conform.Excuse(
				sealed.ShelfChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunShelfSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		sealed.RunShelfSuiteFromHandle(
			t,
			func(*sealed.ShelfHandle) *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(sealed.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.ShelfHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Shelf(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.ShelfHandle) sealed.Shelf { return &shelfSubject{} },
			}),
			conform.Excuse(
				sealed.ShelfChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.SmokeGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.CancelGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ZeroOnErrorGenerated(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.ShelfChecks.Get.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// shelfIdle returns the planted defect of a ceiling of Shelf with a spend
// of nothing, a subject that answers each call at once.
func shelfIdle(h *sealed.ShelfHandle) sealed.Shelf {
	return sealed.ShelfSuite().Calls[0].Over(h, func() {})
}

// TestShelfSuiteBudgets applies the allocation ceilings of the suite of
// Shelf to shelfIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestShelfSuiteBudgets(t *testing.T) {
	sealed.RunShelfBudgets(
		t,
		shelfIdle,
		conform.Ceilings(sealed.ShelfChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			sealed.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuitePut registers the fuzz target of Shelf.Put on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuitePut(f *testing.F) {
	sealed.FuzzShelfSuite(
		f,
		sealed.ShelfChecks.Put.Callable,
		shelfIdle,
		conform.Excuse(
			sealed.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzShelfSuiteGet registers the fuzz target of Shelf.Get on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteGet(f *testing.F) {
	sealed.FuzzShelfSuite(
		f,
		sealed.ShelfChecks.Get.Callable,
		shelfIdle,
		conform.Excuse(
			sealed.ShelfChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.SmokeGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.CancelGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ZeroOnErrorGenerated(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.ShelfChecks.Get.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkShelfSuite measures each callable of the suite of Shelf on
// shelfIdle.
func BenchmarkShelfSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		sealed.BenchShelfSuite(b, sealed.ShelfChecks.Put.Callable, shelfIdle)
	})
	b.Run("Get", func(b *testing.B) {
		sealed.BenchShelfSuite(b, sealed.ShelfChecks.Get.Callable, shelfIdle)
	})
}

// mixerSubjectEmbedded is sealed.Mixer under the name of the field that
// mixerSubject embeds.
type mixerSubjectEmbedded = sealed.Mixer

// mixerSubject is a value of Mixer that only TestMixerSuite builds.
type mixerSubject struct{ mixerSubjectEmbedded }

// TestMixerSuite checks the suite of Mixer and the two forms of its
// application.
func TestMixerSuite(t *testing.T) {
	t.Parallel()

	t.Run("MixerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sealed.MixerSuite())
	})

	t.Run("RunMixerSuite", func(t *testing.T) {
		t.Parallel()
		var built *mixerSubject
		sealed.RunMixerSuite(
			t,
			func() *mixerSubject {
				built = &mixerSubject{}
				return built
			},
			conform.With(sealed.MixerCheck{
				ID:        "Mixer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.MixerHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Mixer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.MixerHandle) sealed.Mixer { return &mixerSubject{} },
			}),
			conform.Excuse(
				sealed.MixerChecks.Mix.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MixerChecks.Mix.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunMixerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *mixerSubject
		sealed.RunMixerSuiteFromHandle(
			t,
			func(*sealed.MixerHandle) *mixerSubject {
				built = &mixerSubject{}
				return built
			},
			conform.With(sealed.MixerCheck{
				ID:        "Mixer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.MixerHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Mixer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.MixerHandle) sealed.Mixer { return &mixerSubject{} },
			}),
			conform.Excuse(
				sealed.MixerChecks.Mix.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.MixerChecks.Mix.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// mixerIdle returns the planted defect of a ceiling of Mixer with a spend
// of nothing, a subject that answers each call at once.
func mixerIdle(h *sealed.MixerHandle) sealed.Mixer {
	return sealed.MixerSuite().Calls[0].Over(h, func() {})
}

// TestMixerSuiteBudgets applies the allocation ceilings of the suite of
// Mixer to mixerIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestMixerSuiteBudgets(t *testing.T) {
	sealed.RunMixerBudgets(
		t,
		mixerIdle,
		conform.Excuse(
			sealed.MixerChecks.Mix.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MixerChecks.Mix.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzMixerSuiteMix registers the fuzz target of Mixer.Mix on mixerIdle,
// with every check of the suite excused.
func FuzzMixerSuiteMix(f *testing.F) {
	sealed.FuzzMixerSuite(
		f,
		sealed.MixerChecks.Mix.Callable,
		mixerIdle,
		conform.Excuse(
			sealed.MixerChecks.Mix.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.MixerChecks.Mix.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkMixerSuite measures each callable of the suite of Mixer on
// mixerIdle.
func BenchmarkMixerSuite(b *testing.B) {
	b.Run("Mix", func(b *testing.B) {
		sealed.BenchMixerSuite(b, sealed.MixerChecks.Mix.Callable, mixerIdle)
	})
}

// depotSubjectEmbedded is sealed.Depot under the name of the field that
// depotSubject embeds.
type depotSubjectEmbedded = sealed.Depot

// depotSubject is a value of Depot that only TestDepotSuite builds.
type depotSubject struct{ depotSubjectEmbedded }

// TestDepotSuite checks the suite of Depot and the two forms of its
// application.
func TestDepotSuite(t *testing.T) {
	t.Parallel()

	t.Run("DepotSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, sealed.DepotSuite())
	})

	t.Run("RunDepotSuite", func(t *testing.T) {
		t.Parallel()
		var built *depotSubject
		sealed.RunDepotSuite(
			t,
			func() *depotSubject {
				built = &depotSubject{}
				return built
			},
			conform.With(sealed.DepotCheck{
				ID:        "Depot/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.DepotHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Depot(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.DepotHandle) sealed.Depot { return &depotSubject{} },
			}),
			conform.Excuse(
				sealed.DepotChecks.Stock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Stock.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunDepotSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *depotSubject
		sealed.RunDepotSuiteFromHandle(
			t,
			func(*sealed.DepotHandle) *depotSubject {
				built = &depotSubject{}
				return built
			},
			conform.With(sealed.DepotCheck{
				ID:        "Depot/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *sealed.DepotHandle) {
					got := h.Subject()
					assert.Equal(h, got, sealed.Depot(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*sealed.DepotHandle) sealed.Depot { return &depotSubject{} },
			}),
			conform.Excuse(
				sealed.DepotChecks.Stock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Stock.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.ReadAfterWrite(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				sealed.DepotChecks.Find.ReadsEachLastWrite(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// depotIdle returns the planted defect of a ceiling of Depot with a spend
// of nothing, a subject that answers each call at once.
func depotIdle(h *sealed.DepotHandle) sealed.Depot {
	return sealed.DepotSuite().Calls[0].Over(h, func() {})
}

// TestDepotSuiteBudgets applies the allocation ceilings of the suite of
// Depot to depotIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestDepotSuiteBudgets(t *testing.T) {
	sealed.RunDepotBudgets(
		t,
		depotIdle,
		conform.Excuse(
			sealed.DepotChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Stock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDepotSuiteStock registers the fuzz target of Depot.Stock on
// depotIdle, with every check of the suite excused.
func FuzzDepotSuiteStock(f *testing.F) {
	sealed.FuzzDepotSuite(
		f,
		sealed.DepotChecks.Stock.Callable,
		depotIdle,
		conform.Excuse(
			sealed.DepotChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Stock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzDepotSuiteFind registers the fuzz target of Depot.Find on depotIdle,
// with every check of the suite excused.
func FuzzDepotSuiteFind(f *testing.F) {
	sealed.FuzzDepotSuite(
		f,
		sealed.DepotChecks.Find.Callable,
		depotIdle,
		conform.Excuse(
			sealed.DepotChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Stock.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.ReadAfterWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			sealed.DepotChecks.Find.ReadsEachLastWrite(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkDepotSuite measures each callable of the suite of Depot on
// depotIdle.
func BenchmarkDepotSuite(b *testing.B) {
	b.Run("Stock", func(b *testing.B) {
		sealed.BenchDepotSuite(b, sealed.DepotChecks.Stock.Callable, depotIdle)
	})
	b.Run("Find", func(b *testing.B) {
		sealed.BenchDepotSuite(b, sealed.DepotChecks.Find.Callable, depotIdle)
	})
}
