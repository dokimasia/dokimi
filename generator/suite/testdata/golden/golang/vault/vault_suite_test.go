package vault_test

import (
	"testing"

	"example.com/shop/vault"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// vaultSubjectEmbedded is vault.Vault under the name of the field that
// vaultSubject embeds.
type vaultSubjectEmbedded = vault.Vault

// vaultSubject is a value of Vault that only TestVaultSuite builds.
type vaultSubject struct{ vaultSubjectEmbedded }

// TestVaultSuite checks the suite of Vault and the two forms of its
// application.
func TestVaultSuite(t *testing.T) {
	t.Parallel()

	t.Run("VaultSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, vault.VaultSuite())
	})

	t.Run("RunVaultSuite", func(t *testing.T) {
		t.Parallel()
		var built *vaultSubject
		vault.RunVaultSuite(
			t,
			func() *vaultSubject {
				built = &vaultSubject{}
				return built
			},
			conform.With(vault.VaultCheck{
				ID:        "Vault/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.VaultHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Vault(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.VaultHandle) vault.Vault { return &vaultSubject{} },
			}),
			conform.ExcuseSuite[vault.Vault]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunVaultSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *vaultSubject
		vault.RunVaultSuiteFromHandle(
			t,
			func(*vault.VaultHandle) *vaultSubject {
				built = &vaultSubject{}
				return built
			},
			conform.With(vault.VaultCheck{
				ID:        "Vault/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.VaultHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Vault(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.VaultHandle) vault.Vault { return &vaultSubject{} },
			}),
			conform.ExcuseSuite[vault.Vault]("the test applies the check of its option alone"),
		)
	})
}

// vaultSuite is the suite of Vault that the tests of the planted defect of
// a ceiling apply.
var vaultSuite = vault.VaultSuite()

// vaultIdle returns the planted defect of a ceiling of Vault with a spend
// of nothing. Each call of the defect returns its zero results at once.
func vaultIdle(h *vault.VaultHandle) vault.Vault {
	return vaultSuite.Calls[0].Over(h, func() {})
}

// TestVaultSuiteBudgets applies the allocation ceilings of the suite of
// Vault to vaultIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestVaultSuiteBudgets(t *testing.T) {
	vault.RunVaultBudgets(
		t,
		vaultIdle,
		conform.Ceilings(vault.VaultChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(vault.VaultChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.ExcuseSuite[vault.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuitePut registers the fuzz target of Vault.Put on vaultIdle,
// with every check of the suite excused.
func FuzzVaultSuitePut(f *testing.F) {
	vault.FuzzVaultSuite(
		f,
		vault.VaultChecks.Put.Callable,
		vaultIdle,
		conform.ExcuseSuite[vault.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzVaultSuiteGet registers the fuzz target of Vault.Get on vaultIdle,
// with every check of the suite excused.
func FuzzVaultSuiteGet(f *testing.F) {
	vault.FuzzVaultSuite(
		f,
		vault.VaultChecks.Get.Callable,
		vaultIdle,
		conform.ExcuseSuite[vault.Vault](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkVaultSuite measures each callable of the suite of Vault on
// vaultIdle.
func BenchmarkVaultSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		vault.BenchVaultSuite(b, vault.VaultChecks.Put.Callable, vaultIdle)
	})
	b.Run("Get", func(b *testing.B) {
		vault.BenchVaultSuite(b, vault.VaultChecks.Get.Callable, vaultIdle)
	})
}

// journalSubjectEmbedded is vault.Journal under the name of the field that
// journalSubject embeds.
type journalSubjectEmbedded = vault.Journal

// journalSubject is a value of Journal that only TestJournalSuite builds.
type journalSubject struct{ journalSubjectEmbedded }

// TestJournalSuite checks the suite of Journal and the two forms of its
// application.
func TestJournalSuite(t *testing.T) {
	t.Parallel()

	t.Run("JournalSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, vault.JournalSuite())
	})

	t.Run("RunJournalSuite", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		vault.RunJournalSuite(
			t,
			func() *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(vault.JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.JournalHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Journal(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.JournalHandle) vault.Journal { return &journalSubject{} },
			}),
			conform.ExcuseSuite[vault.Journal]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunJournalSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		vault.RunJournalSuiteFromHandle(
			t,
			func(*vault.JournalHandle) *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(vault.JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.JournalHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Journal(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.JournalHandle) vault.Journal { return &journalSubject{} },
			}),
			conform.ExcuseSuite[vault.Journal]("the test applies the check of its option alone"),
		)
	})
}

// journalSuite is the suite of Journal that the tests of the planted
// defect of a ceiling apply.
var journalSuite = vault.JournalSuite()

// journalIdle returns the planted defect of a ceiling of Journal with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func journalIdle(h *vault.JournalHandle) vault.Journal {
	return journalSuite.Calls[0].Over(h, func() {})
}

// TestJournalSuiteBudgets applies the allocation ceilings of the suite of
// Journal to journalIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestJournalSuiteBudgets(t *testing.T) {
	vault.RunJournalBudgets(
		t,
		journalIdle,
		conform.Ceilings(
			vault.JournalChecks.Append.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(vault.JournalChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.ExcuseSuite[vault.Journal](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzJournalSuiteAppend registers the fuzz target of Journal.Append on
// journalIdle, with every check of the suite excused.
func FuzzJournalSuiteAppend(f *testing.F) {
	vault.FuzzJournalSuite(
		f,
		vault.JournalChecks.Append.Callable,
		journalIdle,
		conform.ExcuseSuite[vault.Journal](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzJournalSuiteGet registers the fuzz target of Journal.Get on
// journalIdle, with every check of the suite excused.
func FuzzJournalSuiteGet(f *testing.F) {
	vault.FuzzJournalSuite(
		f,
		vault.JournalChecks.Get.Callable,
		journalIdle,
		conform.ExcuseSuite[vault.Journal](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkJournalSuite measures each callable of the suite of Journal on
// journalIdle.
func BenchmarkJournalSuite(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		vault.BenchJournalSuite(b, vault.JournalChecks.Append.Callable, journalIdle)
	})
	b.Run("Get", func(b *testing.B) {
		vault.BenchJournalSuite(b, vault.JournalChecks.Get.Callable, journalIdle)
	})
}

// BenchmarkJournalSuiteParallel measures each callable of the suite of
// Journal that admits concurrent calls on journalIdle, on GOMAXPROCS
// goroutines at once.
func BenchmarkJournalSuiteParallel(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		vault.BenchJournalSuiteParallel(b, vault.JournalChecks.Append.Callable, journalIdle)
	})
}

// ledgerSubjectEmbedded is vault.Ledger under the name of the field that
// ledgerSubject embeds.
type ledgerSubjectEmbedded = vault.Ledger

// ledgerSubject is a value of Ledger that only TestLedgerSuite builds.
type ledgerSubject struct{ ledgerSubjectEmbedded }

// TestLedgerSuite checks the suite of Ledger and the two forms of its
// application.
func TestLedgerSuite(t *testing.T) {
	t.Parallel()

	t.Run("LedgerSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, vault.LedgerSuite())
	})

	t.Run("RunLedgerSuite", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		vault.RunLedgerSuite(
			t,
			func() *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(vault.LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.LedgerHandle) vault.Ledger { return &ledgerSubject{} },
			}),
			conform.ExcuseSuite[vault.Ledger]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunLedgerSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *ledgerSubject
		vault.RunLedgerSuiteFromHandle(
			t,
			func(*vault.LedgerHandle) *ledgerSubject {
				built = &ledgerSubject{}
				return built
			},
			conform.With(vault.LedgerCheck{
				ID:        "Ledger/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.LedgerHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Ledger(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.LedgerHandle) vault.Ledger { return &ledgerSubject{} },
			}),
			conform.ExcuseSuite[vault.Ledger]("the test applies the check of its option alone"),
		)
	})
}

// ledgerSuite is the suite of Ledger that the tests of the planted defect
// of a ceiling apply.
var ledgerSuite = vault.LedgerSuite()

// ledgerIdle returns the planted defect of a ceiling of Ledger with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func ledgerIdle(h *vault.LedgerHandle) vault.Ledger {
	return ledgerSuite.Calls[0].Over(h, func() {})
}

// TestLedgerSuiteBudgets applies the allocation ceilings of the suite of
// Ledger to ledgerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestLedgerSuiteBudgets(t *testing.T) {
	vault.RunLedgerBudgets(
		t,
		ledgerIdle,
		conform.Ceilings(
			vault.LedgerChecks.Record.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			vault.LedgerChecks.Lookup.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[vault.Ledger](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzLedgerSuiteRecord registers the fuzz target of Ledger.Record on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteRecord(f *testing.F) {
	vault.FuzzLedgerSuite(
		f,
		vault.LedgerChecks.Record.Callable,
		ledgerIdle,
		conform.ExcuseSuite[vault.Ledger](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzLedgerSuiteLookup registers the fuzz target of Ledger.Lookup on
// ledgerIdle, with every check of the suite excused.
func FuzzLedgerSuiteLookup(f *testing.F) {
	vault.FuzzLedgerSuite(
		f,
		vault.LedgerChecks.Lookup.Callable,
		ledgerIdle,
		conform.ExcuseSuite[vault.Ledger](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkLedgerSuite measures each callable of the suite of Ledger on
// ledgerIdle.
func BenchmarkLedgerSuite(b *testing.B) {
	b.Run("Record", func(b *testing.B) {
		vault.BenchLedgerSuite(b, vault.LedgerChecks.Record.Callable, ledgerIdle)
	})
	b.Run("Lookup", func(b *testing.B) {
		vault.BenchLedgerSuite(b, vault.LedgerChecks.Lookup.Callable, ledgerIdle)
	})
}

// registrySubjectEmbedded is vault.Registry under the name of the field
// that registrySubject embeds.
type registrySubjectEmbedded = vault.Registry

// registrySubject is a value of Registry that only TestRegistrySuite
// builds.
type registrySubject struct{ registrySubjectEmbedded }

// TestRegistrySuite checks the suite of Registry and the two forms of its
// application.
func TestRegistrySuite(t *testing.T) {
	t.Parallel()

	t.Run("RegistrySuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, vault.RegistrySuite())
	})

	t.Run("RunRegistrySuite", func(t *testing.T) {
		t.Parallel()
		var built *registrySubject
		vault.RunRegistrySuite(
			t,
			func() *registrySubject {
				built = &registrySubject{}
				return built
			},
			conform.With(vault.RegistryCheck{
				ID:        "Registry/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.RegistryHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Registry(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.RegistryHandle) vault.Registry { return &registrySubject{} },
			}),
			conform.ExcuseSuite[vault.Registry]("the test applies the check of its option alone"),
		)
	})

	t.Run("RunRegistrySuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *registrySubject
		vault.RunRegistrySuiteFromHandle(
			t,
			func(*vault.RegistryHandle) *registrySubject {
				built = &registrySubject{}
				return built
			},
			conform.With(vault.RegistryCheck{
				ID:        "Registry/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *vault.RegistryHandle) {
					got := h.Subject()
					assert.Equal(h, got, vault.Registry(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*vault.RegistryHandle) vault.Registry { return &registrySubject{} },
			}),
			conform.ExcuseSuite[vault.Registry]("the test applies the check of its option alone"),
		)
	})
}

// registrySuite is the suite of Registry that the tests of the planted
// defect of a ceiling apply.
var registrySuite = vault.RegistrySuite()

// registryIdle returns the planted defect of a ceiling of Registry with a
// spend of nothing. Each call of the defect returns its zero results at
// once.
func registryIdle(h *vault.RegistryHandle) vault.Registry {
	return registrySuite.Calls[0].Over(h, func() {})
}

// TestRegistrySuiteBudgets applies the allocation ceilings of the suite of
// Registry to registryIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRegistrySuiteBudgets(t *testing.T) {
	vault.RunRegistryBudgets(
		t,
		registryIdle,
		conform.Ceilings(vault.RegistryChecks.Set.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			vault.RegistryChecks.Value.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(vault.RegistryChecks.Add.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			vault.RegistryChecks.Insert.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			vault.RegistryChecks.Stamp.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			vault.RegistryChecks.Touch.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRegistrySuiteSet registers the fuzz target of Registry.Set on
// registryIdle, with every check of the suite excused.
func FuzzRegistrySuiteSet(f *testing.F) {
	vault.FuzzRegistrySuite(
		f,
		vault.RegistryChecks.Set.Callable,
		registryIdle,
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRegistrySuiteValue registers the fuzz target of Registry.Value on
// registryIdle, with every check of the suite excused.
func FuzzRegistrySuiteValue(f *testing.F) {
	vault.FuzzRegistrySuite(
		f,
		vault.RegistryChecks.Value.Callable,
		registryIdle,
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRegistrySuiteAwait registers the fuzz target of Registry.Await on
// registryIdle, with every check of the suite excused.
func FuzzRegistrySuiteAwait(f *testing.F) {
	vault.FuzzRegistrySuite(
		f,
		vault.RegistryChecks.Await.Callable,
		registryIdle,
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRegistrySuiteAdd registers the fuzz target of Registry.Add on
// registryIdle, with every check of the suite excused.
func FuzzRegistrySuiteAdd(f *testing.F) {
	vault.FuzzRegistrySuite(
		f,
		vault.RegistryChecks.Add.Callable,
		registryIdle,
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRegistrySuiteInsert registers the fuzz target of Registry.Insert on
// registryIdle, with every check of the suite excused.
func FuzzRegistrySuiteInsert(f *testing.F) {
	vault.FuzzRegistrySuite(
		f,
		vault.RegistryChecks.Insert.Callable,
		registryIdle,
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRegistrySuiteStamp registers the fuzz target of Registry.Stamp on
// registryIdle, with every check of the suite excused.
func FuzzRegistrySuiteStamp(f *testing.F) {
	vault.FuzzRegistrySuite(
		f,
		vault.RegistryChecks.Stamp.Callable,
		registryIdle,
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// FuzzRegistrySuiteTouch registers the fuzz target of Registry.Touch on
// registryIdle, with every check of the suite excused.
func FuzzRegistrySuiteTouch(f *testing.F) {
	vault.FuzzRegistrySuite(
		f,
		vault.RegistryChecks.Touch.Callable,
		registryIdle,
		conform.ExcuseSuite[vault.Registry](
			"the planted defect of a ceiling returns the zero results of each call",
		),
	)
}

// BenchmarkRegistrySuite measures each callable of the suite of Registry
// on registryIdle.
func BenchmarkRegistrySuite(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		vault.BenchRegistrySuite(b, vault.RegistryChecks.Set.Callable, registryIdle)
	})
	b.Run("Value", func(b *testing.B) {
		vault.BenchRegistrySuite(b, vault.RegistryChecks.Value.Callable, registryIdle)
	})
	b.Run("Add", func(b *testing.B) {
		vault.BenchRegistrySuite(b, vault.RegistryChecks.Add.Callable, registryIdle)
	})
	b.Run("Insert", func(b *testing.B) {
		vault.BenchRegistrySuite(b, vault.RegistryChecks.Insert.Callable, registryIdle)
	})
	b.Run("Stamp", func(b *testing.B) {
		vault.BenchRegistrySuite(b, vault.RegistryChecks.Stamp.Callable, registryIdle)
	})
	b.Run("Touch", func(b *testing.B) {
		vault.BenchRegistrySuite(b, vault.RegistryChecks.Touch.Callable, registryIdle)
	})
}

// BenchmarkRegistrySuiteParallel measures each callable of the suite of
// Registry that admits concurrent calls on registryIdle, on GOMAXPROCS
// goroutines at once.
func BenchmarkRegistrySuiteParallel(b *testing.B) {
	b.Run("Set", func(b *testing.B) {
		vault.BenchRegistrySuiteParallel(b, vault.RegistryChecks.Set.Callable, registryIdle)
	})
	b.Run("Value", func(b *testing.B) {
		vault.BenchRegistrySuiteParallel(b, vault.RegistryChecks.Value.Callable, registryIdle)
	})
	b.Run("Add", func(b *testing.B) {
		vault.BenchRegistrySuiteParallel(b, vault.RegistryChecks.Add.Callable, registryIdle)
	})
	b.Run("Insert", func(b *testing.B) {
		vault.BenchRegistrySuiteParallel(b, vault.RegistryChecks.Insert.Callable, registryIdle)
	})
	b.Run("Stamp", func(b *testing.B) {
		vault.BenchRegistrySuiteParallel(b, vault.RegistryChecks.Stamp.Callable, registryIdle)
	})
	b.Run("Touch", func(b *testing.B) {
		vault.BenchRegistrySuiteParallel(b, vault.RegistryChecks.Touch.Callable, registryIdle)
	})
}
