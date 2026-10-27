package durable

import (
	"testing"

	"go.dokimi.dev/conform"
)

// TestJournalConformsToKeeper applies the checks of the suite of Keeper to
// Journal.
func TestJournalConformsToKeeper(t *testing.T) {
	t.Parallel()

	RunKeeperSuite(t, NewJournal, JournalOptions()...)
}

// BenchmarkJournalConformsToKeeper measures each callable of the suite of
// Keeper on Journal.
func BenchmarkJournalConformsToKeeper(b *testing.B) {
	opts := JournalOptions()
	b.Run("Put", func(b *testing.B) {
		BenchKeeperSuite(
			b,
			KeeperChecks.Put.Callable,
			func(*KeeperHandle) *Journal { return NewJournal() },
			opts...)
	})
	b.Run("Get", func(b *testing.B) {
		BenchKeeperSuite(
			b,
			KeeperChecks.Get.Callable,
			func(*KeeperHandle) *Journal { return NewJournal() },
			opts...)
	})
}

// FuzzJournalConformsToKeeperPut fuzzes the property checks of Keeper.Put
// on Journal.
func FuzzJournalConformsToKeeperPut(f *testing.F) {
	FuzzKeeperSuite(
		f,
		KeeperChecks.Put.Callable,
		func(*KeeperHandle) *Journal { return NewJournal() },
		JournalOptions()...)
}

// FuzzJournalConformsToKeeperGet fuzzes the property checks of Keeper.Get
// on Journal.
func FuzzJournalConformsToKeeperGet(f *testing.F) {
	FuzzKeeperSuite(
		f,
		KeeperChecks.Get.Callable,
		func(*KeeperHandle) *Journal { return NewJournal() },
		JournalOptions()...)
}

// TestTunedConformsToKeeper applies the checks of the suite of Keeper to
// Tuned.
func TestTunedConformsToKeeper(t *testing.T) {
	t.Parallel()

	opts := append([]KeeperOption{
		conform.Cases[Keeper](20),
		conform.Excuse(
			KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, JournalOptions()...)
	RunKeeperSuite(t, NewTuned, opts...)
}

// TestTunedConformsToKeeperAllocs applies the allocation ceilings and the
// checks that count allocations of the suite of Keeper to Tuned.
func TestTunedConformsToKeeperAllocs(t *testing.T) {
	opts := append([]KeeperOption{
		conform.Cases[Keeper](20),
		conform.Ceilings(
			KeeperChecks.Get.Callable,
			conform.Budget{
				Allocs: new(uint64(2)),
			},
		),
		conform.Excuse(
			KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, JournalOptions()...)
	RunKeeperBudgets(t, func(*KeeperHandle) *Tuned { return NewTuned() }, opts...)
}

// BenchmarkTunedConformsToKeeper measures each callable of the suite of
// Keeper on Tuned.
func BenchmarkTunedConformsToKeeper(b *testing.B) {
	opts := append([]KeeperOption{
		conform.Ceilings(
			KeeperChecks.Get.Callable,
			conform.Budget{
				Allocs: new(uint64(2)),
			},
		),
		conform.Excuse(
			KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, JournalOptions()...)
	b.Run("Put", func(b *testing.B) {
		BenchKeeperSuite(
			b,
			KeeperChecks.Put.Callable,
			func(*KeeperHandle) *Tuned { return NewTuned() },
			opts...)
	})
	b.Run("Get", func(b *testing.B) {
		BenchKeeperSuite(
			b,
			KeeperChecks.Get.Callable,
			func(*KeeperHandle) *Tuned { return NewTuned() },
			opts...)
	})
}

// FuzzTunedConformsToKeeperPut fuzzes the property checks of Keeper.Put on
// Tuned.
func FuzzTunedConformsToKeeperPut(f *testing.F) {
	opts := append([]KeeperOption{
		conform.Excuse(
			KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, JournalOptions()...)
	FuzzKeeperSuite(
		f,
		KeeperChecks.Put.Callable,
		func(*KeeperHandle) *Tuned { return NewTuned() },
		opts...)
}

// FuzzTunedConformsToKeeperGet fuzzes the property checks of Keeper.Get on
// Tuned.
func FuzzTunedConformsToKeeperGet(f *testing.F) {
	opts := append([]KeeperOption{
		conform.Excuse(
			KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, JournalOptions()...)
	FuzzKeeperSuite(
		f,
		KeeperChecks.Get.Callable,
		func(*KeeperHandle) *Tuned { return NewTuned() },
		opts...)
}
