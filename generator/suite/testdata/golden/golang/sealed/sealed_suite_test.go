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
			conform.ExcuseSuite[sealed.AEAD]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sealed.AEAD]("the test applies the check of its option alone"),
		)
	})
}

// aeadSuite is the suite of AEAD that the tests of the planted defect of a
// ceiling apply.
var aeadSuite = sealed.AEADSuite()

// aeadIdle returns the planted defect of a ceiling of AEAD with a spend of
// nothing. Each call of the defect returns its zero results at once.
func aeadIdle(h *sealed.AEADHandle) sealed.AEAD {
	return aeadSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sealed.AEAD](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.AEAD](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.AEAD](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.AEAD](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.AEAD](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sealed.Metronome]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sealed.Metronome]("the test applies the check of its option alone"),
		)
	})
}

// metronomeSuite is the suite of Metronome that the tests of the planted
// defect of a ceiling apply.
var metronomeSuite = sealed.MetronomeSuite()

// metronomeIdle returns the planted defect of a ceiling of Metronome with
// a spend of nothing. Each call of the defect returns its zero results at
// once.
func metronomeIdle(h *sealed.MetronomeHandle) sealed.Metronome {
	return metronomeSuite.Calls[0].Over(h, func() {})
}

// TestMetronomeSuiteBudgets applies the allocation ceilings of the suite
// of Metronome to metronomeIdle with every check of the suite excused, and
// a ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestMetronomeSuiteBudgets(t *testing.T) {
	sealed.RunMetronomeBudgets(
		t,
		metronomeIdle,
		conform.ExcuseSuite[sealed.Metronome](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Metronome](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Metronome](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sealed.Locker]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sealed.Locker]("the test applies the check of its option alone"),
		)
	})
}

// lockerSuite is the suite of Locker that the tests of the planted defect
// of a ceiling apply.
var lockerSuite = sealed.LockerSuite()

// lockerIdle returns the planted defect of a ceiling of Locker with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func lockerIdle(h *sealed.LockerHandle) sealed.Locker {
	return lockerSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sealed.Locker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Locker](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Locker](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sealed.Shelf]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sealed.Shelf]("the test applies the check of its option alone"),
		)
	})
}

// shelfSuite is the suite of Shelf that the tests of the planted defect of
// a ceiling apply.
var shelfSuite = sealed.ShelfSuite()

// shelfIdle returns the planted defect of a ceiling of Shelf with a spend
// of nothing. Each call of the defect returns its zero results at once.
func shelfIdle(h *sealed.ShelfHandle) sealed.Shelf {
	return shelfSuite.Calls[0].Over(h, func() {})
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
		conform.ExcuseSuite[sealed.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sealed.Mixer]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sealed.Mixer]("the test applies the check of its option alone"),
		)
	})
}

// mixerSuite is the suite of Mixer that the tests of the planted defect of
// a ceiling apply.
var mixerSuite = sealed.MixerSuite()

// mixerIdle returns the planted defect of a ceiling of Mixer with a spend
// of nothing. Each call of the defect returns its zero results at once.
func mixerIdle(h *sealed.MixerHandle) sealed.Mixer {
	return mixerSuite.Calls[0].Over(h, func() {})
}

// TestMixerSuiteBudgets applies the allocation ceilings of the suite of
// Mixer to mixerIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestMixerSuiteBudgets(t *testing.T) {
	sealed.RunMixerBudgets(
		t,
		mixerIdle,
		conform.ExcuseSuite[sealed.Mixer](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Mixer](
			"the planted defect of a ceiling returns the zero results of each call",
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
			conform.ExcuseSuite[sealed.Depot]("the test applies the check of its option alone"),
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
			conform.ExcuseSuite[sealed.Depot]("the test applies the check of its option alone"),
		)
	})
}

// depotSuite is the suite of Depot that the tests of the planted defect of
// a ceiling apply.
var depotSuite = sealed.DepotSuite()

// depotIdle returns the planted defect of a ceiling of Depot with a spend
// of nothing. Each call of the defect returns its zero results at once.
func depotIdle(h *sealed.DepotHandle) sealed.Depot {
	return depotSuite.Calls[0].Over(h, func() {})
}

// TestDepotSuiteBudgets applies the allocation ceilings of the suite of
// Depot to depotIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestDepotSuiteBudgets(t *testing.T) {
	sealed.RunDepotBudgets(
		t,
		depotIdle,
		conform.ExcuseSuite[sealed.Depot](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Depot](
			"the planted defect of a ceiling returns the zero results of each call",
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
		conform.ExcuseSuite[sealed.Depot](
			"the planted defect of a ceiling returns the zero results of each call",
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
