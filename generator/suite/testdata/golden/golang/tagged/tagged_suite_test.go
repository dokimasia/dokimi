package tagged

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// hasherSubjectEmbedded is Hasher under the name of the field that
// hasherSubject embeds.
type hasherSubjectEmbedded = Hasher

// hasherSubject is a value of Hasher that only TestHasherSuite builds.
type hasherSubject struct{ hasherSubjectEmbedded }

// TestHasherSuite checks the suite of Hasher and the two forms of its
// application.
func TestHasherSuite(t *testing.T) {
	t.Parallel()

	t.Run("HasherSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, HasherSuite())
	})

	t.Run("RunHasherSuite", func(t *testing.T) {
		t.Parallel()
		var built *hasherSubject
		RunHasherSuite(
			t,
			func() *hasherSubject {
				built = &hasherSubject{}
				return built
			},
			conform.With(HasherCheck{
				ID:        "Hasher/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *HasherHandle) {
					got := h.Subject()
					assert.Equal(h, got, Hasher(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*HasherHandle) Hasher { return &hasherSubject{} },
			}),
		)
	})

	t.Run("RunHasherSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *hasherSubject
		RunHasherSuiteFromHandle(
			t,
			func(*HasherHandle) *hasherSubject {
				built = &hasherSubject{}
				return built
			},
			conform.With(HasherCheck{
				ID:        "Hasher/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *HasherHandle) {
					got := h.Subject()
					assert.Equal(h, got, Hasher(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*HasherHandle) Hasher { return &hasherSubject{} },
			}),
		)
	})
}

// hasherIdle returns the planted defect of a ceiling of Hasher with a
// spend of nothing, a subject that answers each call at once.
func hasherIdle(h *HasherHandle) Hasher {
	return HasherSuite().Calls[0].Over(h, func() {})
}

// TestHasherSuiteBudgets applies the allocation ceilings of the suite of
// Hasher to hasherIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestHasherSuiteBudgets(t *testing.T) {
	RunHasherBudgets(
		t,
		hasherIdle,
		conform.Ceilings(HasherChecks.HashTagged.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			HasherChecks.CombineTagged.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(HasherChecks.HashRoles.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(HasherChecks.Truncate.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(HasherChecks.Seal.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzHasherSuiteHashTagged registers the fuzz target of Hasher.HashTagged
// on hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteHashTagged(f *testing.F) {
	FuzzHasherSuite(
		f,
		HasherChecks.HashTagged.Callable,
		hasherIdle,
	)
}

// FuzzHasherSuiteCombineTagged registers the fuzz target of
// Hasher.CombineTagged on hasherIdle, with every check of the suite
// excused.
func FuzzHasherSuiteCombineTagged(f *testing.F) {
	FuzzHasherSuite(
		f,
		HasherChecks.CombineTagged.Callable,
		hasherIdle,
	)
}

// FuzzHasherSuiteHashRoles registers the fuzz target of Hasher.HashRoles
// on hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteHashRoles(f *testing.F) {
	FuzzHasherSuite(
		f,
		HasherChecks.HashRoles.Callable,
		hasherIdle,
	)
}

// FuzzHasherSuiteTruncate registers the fuzz target of Hasher.Truncate on
// hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteTruncate(f *testing.F) {
	FuzzHasherSuite(
		f,
		HasherChecks.Truncate.Callable,
		hasherIdle,
	)
}

// FuzzHasherSuiteSeal registers the fuzz target of Hasher.Seal on
// hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteSeal(f *testing.F) {
	FuzzHasherSuite(
		f,
		HasherChecks.Seal.Callable,
		hasherIdle,
	)
}

// BenchmarkHasherSuite measures each callable of the suite of Hasher on
// hasherIdle.
func BenchmarkHasherSuite(b *testing.B) {
	b.Run("HashTagged", func(b *testing.B) {
		BenchHasherSuite(b, HasherChecks.HashTagged.Callable, hasherIdle)
	})
	b.Run("CombineTagged", func(b *testing.B) {
		BenchHasherSuite(b, HasherChecks.CombineTagged.Callable, hasherIdle)
	})
	b.Run("HashRoles", func(b *testing.B) {
		BenchHasherSuite(b, HasherChecks.HashRoles.Callable, hasherIdle)
	})
	b.Run("Truncate", func(b *testing.B) {
		BenchHasherSuite(b, HasherChecks.Truncate.Callable, hasherIdle)
	})
	b.Run("Seal", func(b *testing.B) {
		BenchHasherSuite(b, HasherChecks.Seal.Callable, hasherIdle)
	})
}

// TestDigestSuite checks the suite of Digest and the two forms of its
// application.
func TestDigestSuite(t *testing.T) {
	t.Parallel()

	t.Run("DigestSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, DigestSuite())
	})

	t.Run("Digest", func(t *testing.T) {
		t.Parallel()
		RunDigestSuite(t, Digest)
	})

	t.Run("RunDigestSuite", func(t *testing.T) {
		t.Parallel()
		RunDigestSuite(
			t,
			Digest,
			conform.With(DigestCheck{
				ID:        "Digest/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *DigestHandle) {
					assert.Equal(h, h.Subject(), DigestFunc(Digest), h.Statement())
				},
				Defect: func(*DigestHandle) DigestFunc { return nil },
			}),
		)
	})

	t.Run("RunDigestSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		RunDigestSuiteFromHandle(
			t,
			func(*DigestHandle) DigestFunc {
				return Digest
			},
			conform.With(DigestCheck{
				ID:        "Digest/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *DigestHandle) {
					assert.Equal(h, h.Subject(), DigestFunc(Digest), h.Statement())
				},
				Defect: func(*DigestHandle) DigestFunc { return nil },
			}),
		)
	})
}

// BenchmarkDigest measures Digest.
func BenchmarkDigest(b *testing.B) {
	BenchDigestSuite(b, Digest)
}

// FuzzDigest fuzzes the property checks of Digest.
func FuzzDigest(f *testing.F) {
	FuzzDigestSuite(f, Digest)
}
