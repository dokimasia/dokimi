package durable_test

import (
	"testing"

	"example.com/shop/durable"
	"go.dokimi.dev/conform"
)

// TestJournalConformsToKeeper applies the checks of the suite of Keeper to
// Journal.
func TestJournalConformsToKeeper(t *testing.T) {
	t.Parallel()

	durable.RunKeeperSuite(t, durable.NewJournal, durable.JournalOptions()...)
}

// BenchmarkJournalConformsToKeeper measures each callable of the suite of
// Keeper on Journal.
func BenchmarkJournalConformsToKeeper(b *testing.B) {
	opts := durable.JournalOptions()
	b.Run("Put", func(b *testing.B) {
		durable.BenchKeeperSuite(
			b,
			durable.KeeperChecks.Put.Callable,
			func(*durable.KeeperHandle) *durable.Journal { return durable.NewJournal() },
			opts...)
	})
	b.Run("Get", func(b *testing.B) {
		durable.BenchKeeperSuite(
			b,
			durable.KeeperChecks.Get.Callable,
			func(*durable.KeeperHandle) *durable.Journal { return durable.NewJournal() },
			opts...)
	})
}

// FuzzJournalConformsToKeeperPut fuzzes the property checks of Keeper.Put
// on Journal.
func FuzzJournalConformsToKeeperPut(f *testing.F) {
	durable.FuzzKeeperSuite(
		f,
		durable.KeeperChecks.Put.Callable,
		func(*durable.KeeperHandle) *durable.Journal { return durable.NewJournal() },
		durable.JournalOptions()...)
}

// FuzzJournalConformsToKeeperGet fuzzes the property checks of Keeper.Get
// on Journal.
func FuzzJournalConformsToKeeperGet(f *testing.F) {
	durable.FuzzKeeperSuite(
		f,
		durable.KeeperChecks.Get.Callable,
		func(*durable.KeeperHandle) *durable.Journal { return durable.NewJournal() },
		durable.JournalOptions()...)
}

// TestTunedConformsToKeeper applies the checks of the suite of Keeper to
// Tuned.
func TestTunedConformsToKeeper(t *testing.T) {
	t.Parallel()

	opts := append([]durable.KeeperOption{
		conform.Cases[durable.Keeper](20),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, durable.JournalOptions()...)
	durable.RunKeeperSuite(t, durable.NewTuned, opts...)
}

// TestTunedConformsToKeeperAllocs applies the allocation ceilings and the
// checks that count allocations of the suite of Keeper to Tuned.
func TestTunedConformsToKeeperAllocs(t *testing.T) {
	opts := append([]durable.KeeperOption{
		conform.Cases[durable.Keeper](20),
		conform.Ceilings(
			durable.KeeperChecks.Get.Callable,
			conform.Budget{
				Allocs: new(uint64(2)),
			},
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, durable.JournalOptions()...)
	durable.RunKeeperBudgets(
		t,
		func(*durable.KeeperHandle) *durable.Tuned { return durable.NewTuned() },
		opts...)
}

// BenchmarkTunedConformsToKeeper measures each callable of the suite of
// Keeper on Tuned.
func BenchmarkTunedConformsToKeeper(b *testing.B) {
	opts := append([]durable.KeeperOption{
		conform.Ceilings(
			durable.KeeperChecks.Get.Callable,
			conform.Budget{
				Allocs: new(uint64(2)),
			},
		),
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, durable.JournalOptions()...)
	b.Run("Put", func(b *testing.B) {
		durable.BenchKeeperSuite(
			b,
			durable.KeeperChecks.Put.Callable,
			func(*durable.KeeperHandle) *durable.Tuned { return durable.NewTuned() },
			opts...)
	})
	b.Run("Get", func(b *testing.B) {
		durable.BenchKeeperSuite(
			b,
			durable.KeeperChecks.Get.Callable,
			func(*durable.KeeperHandle) *durable.Tuned { return durable.NewTuned() },
			opts...)
	})
}

// FuzzTunedConformsToKeeperPut fuzzes the property checks of Keeper.Put on
// Tuned.
func FuzzTunedConformsToKeeperPut(f *testing.F) {
	opts := append([]durable.KeeperOption{
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, durable.JournalOptions()...)
	durable.FuzzKeeperSuite(
		f,
		durable.KeeperChecks.Put.Callable,
		func(*durable.KeeperHandle) *durable.Tuned { return durable.NewTuned() },
		opts...)
}

// FuzzTunedConformsToKeeperGet fuzzes the property checks of Keeper.Get on
// Tuned.
func FuzzTunedConformsToKeeperGet(f *testing.F) {
	opts := append([]durable.KeeperOption{
		conform.Excuse(
			durable.KeeperChecks.Get.KeepsEachPutAcrossARestart(),
			"the tuned journal keeps no medium",
		),
	}, durable.JournalOptions()...)
	durable.FuzzKeeperSuite(
		f,
		durable.KeeperChecks.Get.Callable,
		func(*durable.KeeperHandle) *durable.Tuned { return durable.NewTuned() },
		opts...)
}
