package classified_test

import (
	"testing"

	"example.com/shop/classified"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// shelfSubjectEmbedded is classified.Shelf under the name of the field
// that shelfSubject embeds.
type shelfSubjectEmbedded = classified.Shelf

// shelfSubject is a value of Shelf that only TestShelfSuite builds.
type shelfSubject struct{ shelfSubjectEmbedded }

// TestShelfSuite checks the suite of Shelf and the two forms of its
// application.
func TestShelfSuite(t *testing.T) {
	t.Parallel()

	t.Run("ShelfSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, classified.ShelfSuite())
	})

	t.Run("RunShelfSuite", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		classified.RunShelfSuite(
			t,
			func() *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(classified.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.ShelfHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Shelf(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.ShelfHandle) classified.Shelf { return &shelfSubject{} },
			}),
			conform.ExcuseSuite[classified.Shelf]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunShelfSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *shelfSubject
		classified.RunShelfSuiteFromHandle(
			t,
			func(*classified.ShelfHandle) *shelfSubject {
				built = &shelfSubject{}
				return built
			},
			conform.With(classified.ShelfCheck{
				ID:        "Shelf/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.ShelfHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Shelf(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.ShelfHandle) classified.Shelf { return &shelfSubject{} },
			}),
			conform.ExcuseSuite[classified.Shelf]("the test applies the check of its option alone"),
		)
	})
}

// shelfSuite is the suite of Shelf that the tests of the planted defect of
// a ceiling apply.
var shelfSuite = classified.ShelfSuite()

// shelfIdle returns the planted defect of a ceiling of Shelf with a spend
// of nothing. Each call of the defect returns its zero results at once.
func shelfIdle(h *classified.ShelfHandle) classified.Shelf {
	return shelfSuite.Calls[0].Over(h, func() {})
}

// TestShelfSuiteBudgets applies the allocation ceilings of the suite of
// Shelf to shelfIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestShelfSuiteBudgets(t *testing.T) {
	classified.RunShelfBudgets(
		t,
		shelfIdle,
		conform.Ceilings(
			classified.ShelfChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.GetMany.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.All.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Items.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Scan.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Keyed.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Snapshot.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Len.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.At.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.AtOffset.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Slot.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.TracedGet.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.ShelfChecks.Err.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuitePut registers the fuzz target of Shelf.Put on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuitePut(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Put.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteGetMany registers the fuzz target of Shelf.GetMany on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteGetMany(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.GetMany.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteAll registers the fuzz target of Shelf.All on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteAll(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.All.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteItems registers the fuzz target of Shelf.Items on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteItems(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Items.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteScan registers the fuzz target of Shelf.Scan on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteScan(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Scan.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteKeyed registers the fuzz target of Shelf.Keyed on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteKeyed(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Keyed.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteSnapshot registers the fuzz target of Shelf.Snapshot on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteSnapshot(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Snapshot.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteLen registers the fuzz target of Shelf.Len on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteLen(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Len.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteAt registers the fuzz target of Shelf.At on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteAt(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.At.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteAtOffset registers the fuzz target of Shelf.AtOffset on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteAtOffset(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.AtOffset.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteSlot registers the fuzz target of Shelf.Slot on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteSlot(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Slot.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteGet registers the fuzz target of Shelf.Get on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteGet(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Get.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteTracedGet registers the fuzz target of Shelf.TracedGet on
// shelfIdle, with every check of the suite excused.
func FuzzShelfSuiteTracedGet(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.TracedGet.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzShelfSuiteErr registers the fuzz target of Shelf.Err on shelfIdle,
// with every check of the suite excused.
func FuzzShelfSuiteErr(f *testing.F) {
	classified.FuzzShelfSuite(
		f,
		classified.ShelfChecks.Err.Callable,
		shelfIdle,
		conform.ExcuseSuite[classified.Shelf](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkShelfSuite measures each callable of the suite of Shelf on
// shelfIdle.
func BenchmarkShelfSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Put.Callable, shelfIdle)
	})
	b.Run("GetMany", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.GetMany.Callable, shelfIdle)
	})
	b.Run("All", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.All.Callable, shelfIdle)
	})
	b.Run("Items", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Items.Callable, shelfIdle)
	})
	b.Run("Scan", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Scan.Callable, shelfIdle)
	})
	b.Run("Keyed", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Keyed.Callable, shelfIdle)
	})
	b.Run("Snapshot", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Snapshot.Callable, shelfIdle)
	})
	b.Run("Len", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Len.Callable, shelfIdle)
	})
	b.Run("At", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.At.Callable, shelfIdle)
	})
	b.Run("AtOffset", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.AtOffset.Callable, shelfIdle)
	})
	b.Run("Slot", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Slot.Callable, shelfIdle)
	})
	b.Run("Get", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Get.Callable, shelfIdle)
	})
	b.Run("TracedGet", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.TracedGet.Callable, shelfIdle)
	})
	b.Run("Err", func(b *testing.B) {
		classified.BenchShelfSuite(b, classified.ShelfChecks.Err.Callable, shelfIdle)
	})
}

// tagsSubjectEmbedded is classified.Tags under the name of the field that
// tagsSubject embeds.
type tagsSubjectEmbedded = classified.Tags

// tagsSubject is a value of Tags that only TestTagsSuite builds.
type tagsSubject struct{ tagsSubjectEmbedded }

// TestTagsSuite checks the suite of Tags and the two forms of its
// application.
func TestTagsSuite(t *testing.T) {
	t.Parallel()

	t.Run("TagsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, classified.TagsSuite())
	})

	t.Run("RunTagsSuite", func(t *testing.T) {
		t.Parallel()
		var built *tagsSubject
		classified.RunTagsSuite(
			t,
			func() *tagsSubject {
				built = &tagsSubject{}
				return built
			},
			conform.With(classified.TagsCheck{
				ID:        "Tags/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.TagsHandle) {
					got := h.Subject()
					assert.Equal(h, got, classified.Tags(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*classified.TagsHandle) classified.Tags { return &tagsSubject{} },
			}),
			conform.ExcuseSuite[classified.Tags]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunTagsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *tagsSubject
		classified.RunTagsSuiteFromHandle(
			t,
			func(*classified.TagsHandle) *tagsSubject {
				built = &tagsSubject{}
				return built
			},
			conform.With(classified.TagsCheck{
				ID:        "Tags/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.TagsHandle) {
					got := h.Subject()
					assert.Equal(h, got, classified.Tags(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*classified.TagsHandle) classified.Tags { return &tagsSubject{} },
			}),
			conform.ExcuseSuite[classified.Tags]("the test applies the check of its option alone"),
		)
	})
}

// tagsSuite is the suite of Tags that the tests of the planted defect of a
// ceiling apply.
var tagsSuite = classified.TagsSuite()

// tagsIdle returns the planted defect of a ceiling of Tags with a spend of
// nothing. Each call of the defect returns its zero results at once.
func tagsIdle(h *classified.TagsHandle) classified.Tags {
	return tagsSuite.Calls[0].Over(h, func() {})
}

// TestTagsSuiteBudgets applies the allocation ceilings of the suite of
// Tags to tagsIdle with every check of the suite excused, and a ceiling of
// 0 to each callable without a budget whose calls build no input of their
// own.
func TestTagsSuiteBudgets(t *testing.T) {
	classified.RunTagsBudgets(
		t,
		tagsIdle,
		conform.Ceilings(
			classified.TagsChecks.Add.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.TagsChecks.Has.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[classified.Tags](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTagsSuiteAdd registers the fuzz target of Tags.Add on tagsIdle, with
// every check of the suite excused.
func FuzzTagsSuiteAdd(f *testing.F) {
	classified.FuzzTagsSuite(
		f,
		classified.TagsChecks.Add.Callable,
		tagsIdle,
		conform.ExcuseSuite[classified.Tags](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzTagsSuiteHas registers the fuzz target of Tags.Has on tagsIdle, with
// every check of the suite excused.
func FuzzTagsSuiteHas(f *testing.F) {
	classified.FuzzTagsSuite(
		f,
		classified.TagsChecks.Has.Callable,
		tagsIdle,
		conform.ExcuseSuite[classified.Tags](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkTagsSuite measures each callable of the suite of Tags on
// tagsIdle.
func BenchmarkTagsSuite(b *testing.B) {
	b.Run("Add", func(b *testing.B) {
		classified.BenchTagsSuite(b, classified.TagsChecks.Add.Callable, tagsIdle)
	})
	b.Run("Has", func(b *testing.B) {
		classified.BenchTagsSuite(b, classified.TagsChecks.Has.Callable, tagsIdle)
	})
}

// vaultSubjectEmbedded is classified.Vault under the name of the field
// that vaultSubject embeds.
type vaultSubjectEmbedded = classified.Vault

// vaultSubject is a value of Vault that only TestVaultSuite builds.
type vaultSubject struct{ vaultSubjectEmbedded }

// TestVaultSuite checks the suite of Vault and the two forms of its
// application.
func TestVaultSuite(t *testing.T) {
	t.Parallel()

	t.Run("VaultSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, classified.VaultSuite())
	})

	t.Run("RunVaultSuite", func(t *testing.T) {
		t.Parallel()
		var built *vaultSubject
		classified.RunVaultSuite(
			t,
			func() *vaultSubject {
				built = &vaultSubject{}
				return built
			},
			conform.With(classified.VaultCheck{
				ID:        "Vault/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.VaultHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Vault(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.VaultHandle) classified.Vault { return &vaultSubject{} },
			}),
			conform.ExcuseSuite[classified.Vault]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunVaultSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *vaultSubject
		classified.RunVaultSuiteFromHandle(
			t,
			func(*classified.VaultHandle) *vaultSubject {
				built = &vaultSubject{}
				return built
			},
			conform.With(classified.VaultCheck{
				ID:        "Vault/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *classified.VaultHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						classified.Vault(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*classified.VaultHandle) classified.Vault { return &vaultSubject{} },
			}),
			conform.ExcuseSuite[classified.Vault]("the test applies the check of its option alone"),
		)
	})
}

// vaultSuite is the suite of Vault that the tests of the planted defect of
// a ceiling apply.
var vaultSuite = classified.VaultSuite()

// vaultIdle returns the planted defect of a ceiling of Vault with a spend
// of nothing. Each call of the defect returns its zero results at once.
func vaultIdle(h *classified.VaultHandle) classified.Vault {
	return vaultSuite.Calls[0].Over(h, func() {})
}

// TestVaultSuiteBudgets applies the allocation ceilings of the suite of
// Vault to vaultIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestVaultSuiteBudgets(t *testing.T) {
	classified.RunVaultBudgets(
		t,
		vaultIdle,
		conform.Ceilings(
			classified.VaultChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Replay.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Verify.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Acquire.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Release.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Borrow.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Return.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			classified.VaultChecks.Stats.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteAppend registers the fuzz target of Vault.Append on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteAppend(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Append.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteReplay registers the fuzz target of Vault.Replay on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteReplay(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Replay.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteVerify registers the fuzz target of Vault.Verify on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteVerify(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Verify.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteAcquire registers the fuzz target of Vault.Acquire on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteAcquire(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Acquire.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteRelease registers the fuzz target of Vault.Release on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteRelease(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Release.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteBorrow registers the fuzz target of Vault.Borrow on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteBorrow(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Borrow.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteReturn registers the fuzz target of Vault.Return on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteReturn(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Return.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteStats registers the fuzz target of Vault.Stats on
// vaultIdle, with every check of the suite excused.
func FuzzVaultSuiteStats(f *testing.F) {
	classified.FuzzVaultSuite(
		f,
		classified.VaultChecks.Stats.Callable,
		vaultIdle,
		conform.ExcuseSuite[classified.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkVaultSuite measures each callable of the suite of Vault on
// vaultIdle.
func BenchmarkVaultSuite(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Append.Callable, vaultIdle)
	})
	b.Run("Replay", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Replay.Callable, vaultIdle)
	})
	b.Run("Verify", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Verify.Callable, vaultIdle)
	})
	b.Run("Acquire", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Acquire.Callable, vaultIdle)
	})
	b.Run("Release", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Release.Callable, vaultIdle)
	})
	b.Run("Borrow", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Borrow.Callable, vaultIdle)
	})
	b.Run("Return", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Return.Callable, vaultIdle)
	})
	b.Run("Stats", func(b *testing.B) {
		classified.BenchVaultSuite(b, classified.VaultChecks.Stats.Callable, vaultIdle)
	})
}
