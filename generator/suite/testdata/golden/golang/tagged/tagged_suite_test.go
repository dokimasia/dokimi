package tagged_test

import (
	"testing"

	"example.com/shop/tagged"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// hasherSubjectEmbedded is tagged.Hasher under the name of the field that
// hasherSubject embeds.
type hasherSubjectEmbedded = tagged.Hasher

// hasherSubject is a value of Hasher that only TestHasherSuite builds.
type hasherSubject struct{ hasherSubjectEmbedded }

// TestHasherSuite checks the suite of Hasher and the two forms of its
// application.
func TestHasherSuite(t *testing.T) {
	t.Parallel()

	t.Run("HasherSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, tagged.HasherSuite())
	})

	t.Run("RunHasherSuite", func(t *testing.T) {
		t.Parallel()
		var built *hasherSubject
		tagged.RunHasherSuite(
			t,
			func() *hasherSubject {
				built = &hasherSubject{}
				return built
			},
			conform.With(tagged.HasherCheck{
				ID:        "Hasher/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *tagged.HasherHandle) {
					got := h.Subject()
					assert.Equal(h, got, tagged.Hasher(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*tagged.HasherHandle) tagged.Hasher { return &hasherSubject{} },
			}),
			conform.Excuse(
				tagged.HasherChecks.HashTagged.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.HashTagged.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.CombineTagged.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.CombineTagged.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.HashRoles.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.HashRoles.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Truncate.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Truncate.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Seal.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Seal.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunHasherSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *hasherSubject
		tagged.RunHasherSuiteFromHandle(
			t,
			func(*tagged.HasherHandle) *hasherSubject {
				built = &hasherSubject{}
				return built
			},
			conform.With(tagged.HasherCheck{
				ID:        "Hasher/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *tagged.HasherHandle) {
					got := h.Subject()
					assert.Equal(h, got, tagged.Hasher(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*tagged.HasherHandle) tagged.Hasher { return &hasherSubject{} },
			}),
			conform.Excuse(
				tagged.HasherChecks.HashTagged.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.HashTagged.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.CombineTagged.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.CombineTagged.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.HashRoles.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.HashRoles.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Truncate.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Truncate.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Seal.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.HasherChecks.Seal.Cancel(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// hasherIdle returns the planted defect of a ceiling of Hasher with a
// spend of nothing, a subject that answers each call at once.
func hasherIdle(h *tagged.HasherHandle) tagged.Hasher {
	return tagged.HasherSuite().Calls[0].Over(h, func() {})
}

// TestHasherSuiteBudgets applies the allocation ceilings of the suite of
// Hasher to hasherIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestHasherSuiteBudgets(t *testing.T) {
	tagged.RunHasherBudgets(
		t,
		hasherIdle,
		conform.Ceilings(
			tagged.HasherChecks.HashTagged.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			tagged.HasherChecks.CombineTagged.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			tagged.HasherChecks.HashRoles.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			tagged.HasherChecks.Truncate.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(tagged.HasherChecks.Seal.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHasherSuiteHashTagged registers the fuzz target of Hasher.HashTagged
// on hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteHashTagged(f *testing.F) {
	tagged.FuzzHasherSuite(
		f,
		tagged.HasherChecks.HashTagged.Callable,
		hasherIdle,
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHasherSuiteCombineTagged registers the fuzz target of
// Hasher.CombineTagged on hasherIdle, with every check of the suite
// excused.
func FuzzHasherSuiteCombineTagged(f *testing.F) {
	tagged.FuzzHasherSuite(
		f,
		tagged.HasherChecks.CombineTagged.Callable,
		hasherIdle,
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHasherSuiteHashRoles registers the fuzz target of Hasher.HashRoles
// on hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteHashRoles(f *testing.F) {
	tagged.FuzzHasherSuite(
		f,
		tagged.HasherChecks.HashRoles.Callable,
		hasherIdle,
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHasherSuiteTruncate registers the fuzz target of Hasher.Truncate on
// hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteTruncate(f *testing.F) {
	tagged.FuzzHasherSuite(
		f,
		tagged.HasherChecks.Truncate.Callable,
		hasherIdle,
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzHasherSuiteSeal registers the fuzz target of Hasher.Seal on
// hasherIdle, with every check of the suite excused.
func FuzzHasherSuiteSeal(f *testing.F) {
	tagged.FuzzHasherSuite(
		f,
		tagged.HasherChecks.Seal.Callable,
		hasherIdle,
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.CombineTagged.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.HashRoles.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Truncate.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			tagged.HasherChecks.Seal.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkHasherSuite measures each callable of the suite of Hasher on
// hasherIdle.
func BenchmarkHasherSuite(b *testing.B) {
	b.Run("HashTagged", func(b *testing.B) {
		tagged.BenchHasherSuite(b, tagged.HasherChecks.HashTagged.Callable, hasherIdle)
	})
	b.Run("CombineTagged", func(b *testing.B) {
		tagged.BenchHasherSuite(b, tagged.HasherChecks.CombineTagged.Callable, hasherIdle)
	})
	b.Run("HashRoles", func(b *testing.B) {
		tagged.BenchHasherSuite(b, tagged.HasherChecks.HashRoles.Callable, hasherIdle)
	})
	b.Run("Truncate", func(b *testing.B) {
		tagged.BenchHasherSuite(b, tagged.HasherChecks.Truncate.Callable, hasherIdle)
	})
	b.Run("Seal", func(b *testing.B) {
		tagged.BenchHasherSuite(b, tagged.HasherChecks.Seal.Callable, hasherIdle)
	})
}

// TestDigestSuite checks the suite of Digest and the two forms of its
// application.
func TestDigestSuite(t *testing.T) {
	t.Parallel()

	t.Run("DigestSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, tagged.DigestSuite())
	})

	t.Run("Digest", func(t *testing.T) {
		t.Parallel()
		tagged.RunDigestSuite(t, tagged.Digest)
	})

	t.Run("RunDigestSuite", func(t *testing.T) {
		t.Parallel()
		tagged.RunDigestSuite(
			t,
			tagged.Digest,
			conform.With(tagged.DigestCheck{
				ID:        "Digest/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *tagged.DigestHandle) {
					assert.Equal(h, h.Subject(), tagged.DigestFunc(tagged.Digest), h.Statement())
				},
				Defect: func(*tagged.DigestHandle) tagged.DigestFunc { return nil },
			}),
			conform.Excuse(
				tagged.DigestChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.DigestChecks.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunDigestSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		tagged.RunDigestSuiteFromHandle(
			t,
			func(*tagged.DigestHandle) tagged.DigestFunc {
				return tagged.Digest
			},
			conform.With(tagged.DigestCheck{
				ID:        "Digest/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *tagged.DigestHandle) {
					assert.Equal(h, h.Subject(), tagged.DigestFunc(tagged.Digest), h.Statement())
				},
				Defect: func(*tagged.DigestHandle) tagged.DigestFunc { return nil },
			}),
			conform.Excuse(
				tagged.DigestChecks.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				tagged.DigestChecks.Deterministic(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// BenchmarkDigest measures Digest.
func BenchmarkDigest(b *testing.B) {
	tagged.BenchDigestSuite(b, tagged.Digest)
}

// FuzzDigest fuzzes the property checks of Digest.
func FuzzDigest(f *testing.F) {
	tagged.FuzzDigestSuite(f, tagged.Digest)
}
