package vault

import (
	"context"
	"errors"
	"io/fs"
	"sync/atomic"
	"testing"

	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/host"
	"go.dokimi.dev/medium"
	"go.dokimi.dev/spec"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Vault writes.
type (
	VaultHandle = conform.Handle[Vault]
	VaultCheck  = conform.Check[Vault]
	VaultOption = conform.Option[Vault]
)

// VaultChecks contains the name of each callable of Vault, with the
// identities of its checks.
var VaultChecks = struct {
	Put vaultPut
	Get vaultGet
}{
	Put: vaultPut{"Vault.Put"},
	Get: vaultGet{"Vault.Get"},
}

// vaultPut is the name of Vault.Put, with the identities of its checks.
type vaultPut struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c vaultPut) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c vaultPut) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// SyncedWritesSurvive returns the identity of the check that Get returns
// each value that Put wrote with success before a crash at any persistence
// point.
func (c vaultPut) SyncedWritesSurvive() conform.ID[Vault] {
	return c.Check("synced-writes-survive")
}

// NoTornRecord returns the identity of the check that Get returns each
// value of Put whole or not at all after a crash during a torn write.
func (c vaultPut) NoTornRecord() conform.ID[Vault] {
	return c.Check("no-torn-record")
}

// CorruptionIsReported returns the identity of the check that Get fails
// with the declared error of a corrupted record, and returns no value of
// it.
func (c vaultPut) CorruptionIsReported() conform.ID[Vault] {
	return c.Check("corruption-is-reported")
}

// FullMediumIsReported returns the identity of the check that Put fails on
// a full medium, and Get returns each value as before the write.
func (c vaultPut) FullMediumIsReported() conform.ID[Vault] {
	return c.Check("full-medium-is-reported")
}

// CompletesAfterRecovery returns the identity of the check that each call
// of Vault completes within the bound of the world once every fault ends
// and the subject opens again.
func (c vaultPut) CompletesAfterRecovery() conform.ID[Vault] {
	return "Vault/completes-after-recovery"
}

// vaultGet is the name of Vault.Get, with the identities of its checks.
type vaultGet struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c vaultGet) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c vaultGet) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c vaultGet) ZeroOnError() conform.ID[Vault] {
	return c.Check("zero-on-error")
}

// VaultSuite returns the checks of Vault, and the calls and the budgets of
// its callables.
func VaultSuite() conform.Suite[Vault] {
	return conform.Suite[Vault]{
		Declaration: "Vault",
		Package:     "example.com/shop/vault",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Vault]{
			{
				ID:          VaultChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: vaultPutPanics,
			},
			{
				ID:          VaultChecks.Put.Smoke().Generated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnKey, drawnValue)
					})
				},
				Defect: vaultPutPanics,
			},
			{
				ID:          VaultChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: vaultPutIgnoresCancel,
			},
			{
				ID:          VaultChecks.Put.Cancel().Generated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnKey, drawnValue)
					})
				},
				Defect: vaultPutIgnoresCancel,
			},
			{
				ID:          VaultChecks.Put.SyncedWritesSurvive(),
				Statement:   "Get returns each value that Put wrote with success before a crash at any persistence point",
				Claim:       spec.SyncedWritesSurvive,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *VaultHandle) {
					probe.SurvivesEveryCrashPoint(h, probe.Durable[Vault, string]{
						Write: func(ctx context.Context, s Vault, i int) error {
							if i == 0 {
								return s.Put(ctx, string("key-a"), string("value-a"))
							}
							return s.Put(ctx, string("key-b"), string("value-b"))
						},
						Read: func(ctx context.Context, s Vault, i int) (string, error) {
							if i == 0 {
								return s.Get(ctx, string("key-a"))
							}
							return s.Get(ctx, string("key-b"))
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: vaultPutAcknowledgesBeforeSync,
				Doors:  conform.Doors[Vault]{Open: vaultPutAcknowledgesBeforeSyncOpen},
			},
			{
				ID:          VaultChecks.Put.NoTornRecord(),
				Statement:   "Get returns each value of Put whole or not at all after a crash during a torn write",
				Claim:       spec.NoTornRecord,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *VaultHandle) {
					probe.WholeAtEveryCrashPoint(h, probe.Durable[Vault, string]{
						Write: func(ctx context.Context, s Vault, i int) error {
							if i == 0 {
								return s.Put(ctx, string("key-a"), string("value-a"))
							}
							return s.Put(ctx, string("key-b"), string("value-b"))
						},
						Read: func(ctx context.Context, s Vault, i int) (string, error) {
							if i == 0 {
								return s.Get(ctx, string("key-a"))
							}
							return s.Get(ctx, string("key-b"))
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: vaultPutTrustsAPartialRecord,
				Doors:  conform.Doors[Vault]{Open: vaultPutTrustsAPartialRecordOpen},
			},
			{
				ID:          VaultChecks.Put.CorruptionIsReported(),
				Statement:   "Get fails with the declared error of a corrupted record, and returns no value of it",
				Claim:       spec.CorruptionIsReported,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *VaultHandle) {
					probe.ReportsEachFlippedBit(h, probe.Durable[Vault, string]{
						Write: func(ctx context.Context, s Vault, i int) error {
							if i == 0 {
								return s.Put(ctx, string("key-a"), string("value-a"))
							}
							return s.Put(ctx, string("key-b"), string("value-b"))
						},
						Read: func(ctx context.Context, s Vault, i int) (string, error) {
							if i == 0 {
								return s.Get(ctx, string("key-a"))
							}
							return s.Get(ctx, string("key-b"))
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					}, ErrCorrupt)
				},
				Defect: vaultGetIgnoresChecksum,
				Doors:  conform.Doors[Vault]{Open: vaultGetIgnoresChecksumOpen},
			},
			{
				ID:          VaultChecks.Put.FullMediumIsReported(),
				Statement:   "Put fails on a full medium, and Get returns each value as before the write",
				Claim:       spec.FullMediumIsReported,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *VaultHandle) {
					probe.RefusesOnAFullMedium(h, probe.Durable[Vault, string]{
						Write: func(ctx context.Context, s Vault, i int) error {
							if i == 0 {
								return s.Put(ctx, string("key-a"), string("value-a"))
							}
							return s.Put(ctx, string("key-b"), string("value-b"))
						},
						Read: func(ctx context.Context, s Vault, i int) (string, error) {
							if i == 0 {
								return s.Get(ctx, string("key-a"))
							}
							return s.Get(ctx, string("key-b"))
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: vaultPutHidesAFullMedium,
				Doors:  conform.Doors[Vault]{Open: vaultPutHidesAFullMediumOpen},
			},
			{
				ID:        VaultChecks.Put.CompletesAfterRecovery(),
				Statement: "each call of Vault completes within the bound of the world once every fault ends and the subject opens again",
				Claim:     spec.CompletesAfterRecovery,
				Machineries: []spec.MachineryID{
					spec.Quiescence,
					spec.GeneratedSequences,
					spec.KillableProcessOrBreakableMedium,
				},
				Property: func(h *VaultHandle, c *prop.Case) {
					probe.WorldRecovery(h, c, vaultWorld)
				},
				Defect:  vaultBlocksAfterRecovery,
				Doors:   conform.Doors[Vault]{Host: vaultBlocksAfterRecoveryHost},
				Around:  conform.Bubble,
				Options: probe.WorldOptions[Vault],
			},
			{
				ID:          VaultChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: vaultGetPanics,
			},
			{
				ID:          VaultChecks.Get.Smoke().Generated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: vaultGetPanics,
			},
			{
				ID:          VaultChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: vaultGetIgnoresCancel,
			},
			{
				ID:          VaultChecks.Get.Cancel().Generated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: vaultGetIgnoresCancel,
			},
			{
				ID:          VaultChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, string(""))
				},
				Defect: vaultGetValueBeside,
			},
			{
				ID:          VaultChecks.Get.ZeroOnError().Generated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, string(""))
				},
				Defect: vaultGetValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Vault.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Vault]{
			{
				Calls:    vaultPutCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Put.Callable,
			},
			{
				Calls:    vaultGetCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Get.Callable,
			},
		},
	}
}

// RunVaultSuite applies the checks of Vault, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunVaultSuite[T Vault](t *testing.T, newSubject func() T, opts ...VaultOption) {
	t.Helper()
	RunVaultSuiteFromHandle(t, func(*VaultHandle) T { return newSubject() }, opts...)
}

// RunVaultSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunVaultSuiteFromHandle[T Vault](
	t *testing.T,
	newSubject func(h *VaultHandle) T,
	opts ...VaultOption,
) {
	t.Helper()
	conform.Run(t, VaultSuite(), vaultSubjects(newSubject), opts...)
}

// FuzzVaultSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzVaultSuite[T Vault](
	f *testing.F,
	callable conform.Callable[Vault],
	newSubject func(h *VaultHandle) T,
	opts ...VaultOption,
) {
	f.Helper()
	conform.Fuzz(f, VaultSuite(), callable, vaultSubjects(newSubject), opts...)
}

// RunVaultBudgets applies the allocation ceilings of Vault and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunVaultBudgets[T Vault](
	t *testing.T,
	newSubject func(h *VaultHandle) T,
	opts ...VaultOption,
) {
	t.Helper()
	conform.RunBudgets(t, VaultSuite(), vaultSubjects(newSubject), opts...)
}

// BenchVaultSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchVaultSuite[T Vault](
	b conform.B,
	callable conform.Callable[Vault],
	newSubject func(h *VaultHandle) T,
	opts ...VaultOption,
) {
	b.Helper()
	conform.Bench(b, VaultSuite(), callable, vaultSubjects(newSubject), opts...)
}

// vaultSubjects returns the factory of the subjects of Vault that calls
// newSubject with the handle of each check.
func vaultSubjects[T Vault](newSubject func(h *VaultHandle) T) func(h *VaultHandle) Vault {
	return func(h *VaultHandle) Vault { return newSubject(h) }
}

// vaultPutCalls generates the calls of Vault.Put with the inputs that each
// case draws.
var vaultPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Vault) error {
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Vault) error {
		return subject.Put(ctx, key, value)
	}
})

// vaultGetCalls generates the calls of Vault.Get with the inputs that each
// case draws.
var vaultGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Vault) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Vault) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// vaultOverEmbedded is Vault under the name of the field that vaultOver
// embeds.
type vaultOverEmbedded = Vault

// vaultOver is a subject of Vault whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type vaultOver struct {
	vaultOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o vaultOver) Put(context.Context, string, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o vaultOver) Get(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newVaultOver returns the planted defect of a ceiling of Vault, whose
// callables call spend.
func newVaultOver(_ *VaultHandle, spend func()) Vault {
	return vaultOver{spend: spend}
}

// vaultPutPanicsDefect is the planted defect of a Vault whose Put panics.
type vaultPutPanicsDefect struct {
	vaultOverEmbedded
}

// Put is the override of the planted defect.
func (*vaultPutPanicsDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	panic("the planted defect of a Vault whose Put panics")
}

// vaultPutPanics returns the planted defect of a Vault whose Put panics
// over a fresh subject.
func vaultPutPanics(h *VaultHandle) Vault {
	return &vaultPutPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultPutIgnoresCancelDefect is the planted defect of a Vault whose Put
// ignores the context.
type vaultPutIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Put is the override of the planted defect.
func (d *vaultPutIgnoresCancelDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.vaultOverEmbedded.Put(p0, p1, p2)
}

// vaultPutIgnoresCancel returns the planted defect of a Vault whose Put
// ignores the context over a fresh subject.
func vaultPutIgnoresCancel(h *VaultHandle) Vault {
	return &vaultPutIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultPutAcknowledgesBeforeSync returns the planted defect of a Vault
// over a medium whose syncs persist nothing. The door of its check opens
// the defect over a new medium in memory.
func vaultPutAcknowledgesBeforeSync(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultPutAcknowledgesBeforeSyncOpen opens the planted defect of a Vault
// over a medium whose syncs persist nothing. The application opens the
// subject over the medium m that probe.Unsynced wraps.
func vaultPutAcknowledgesBeforeSyncOpen(h *VaultHandle, m medium.FS) Vault {
	return h.FreshOpen(probe.Unsynced(m))
}

// vaultPutTrustsAPartialRecordDefect is the planted defect of a Vault
// whose Put writes the zero value under the key before its first value.
type vaultPutTrustsAPartialRecordDefect struct {
	vaultOverEmbedded
	wrote bool
}

// Put is the override of the planted defect.
func (d *vaultPutTrustsAPartialRecordDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	if !d.wrote {
		d.wrote = true
		r0 = d.vaultOverEmbedded.Put(p0, p1, *new(string))
		if r0 != nil {
			return r0
		}
	}
	return d.vaultOverEmbedded.Put(p0, p1, p2)
}

// vaultPutTrustsAPartialRecord returns the planted defect of a Vault whose
// Put writes the zero value under the key before its first value. The door
// of its check opens the defect over a new medium in memory.
func vaultPutTrustsAPartialRecord(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultPutTrustsAPartialRecordOpen opens the planted defect of a Vault
// whose Put writes the zero value under the key before its first value.
// The defect wraps a subject that the application opens over the medium m.
func vaultPutTrustsAPartialRecordOpen(h *VaultHandle, m medium.FS) Vault {
	return &vaultPutTrustsAPartialRecordDefect{vaultOverEmbedded: h.FreshOpen(m)}
}

// vaultGetIgnoresChecksumDefect is the planted defect of a Vault whose Get
// returns the zero value of a corrupted record.
type vaultGetIgnoresChecksumDefect struct {
	vaultOverEmbedded
}

// Get is the override of the planted defect.
func (d *vaultGetIgnoresChecksumDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.vaultOverEmbedded.Get(p0, p1)
	if errors.Is(r1, ErrCorrupt) {
		r0 = *new(string)
		r1 = nil
	}
	return r0, r1
}

// vaultGetIgnoresChecksum returns the planted defect of a Vault whose Get
// returns the zero value of a corrupted record. The door of its check
// opens the defect over a new medium in memory.
func vaultGetIgnoresChecksum(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultGetIgnoresChecksumOpen opens the planted defect of a Vault whose
// Get returns the zero value of a corrupted record. The defect wraps a
// subject that the application opens over the medium m.
func vaultGetIgnoresChecksumOpen(h *VaultHandle, m medium.FS) Vault {
	return &vaultGetIgnoresChecksumDefect{vaultOverEmbedded: h.FreshOpen(m)}
}

// vaultPutHidesAFullMedium returns the planted defect of a Vault over a
// medium that reports a full write as written. The door of its check opens
// the defect over a new medium in memory.
func vaultPutHidesAFullMedium(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultPutHidesAFullMediumOpen opens the planted defect of a Vault over a
// medium that reports a full write as written. The application opens the
// subject over the medium m that probe.Unfilled wraps.
func vaultPutHidesAFullMediumOpen(h *VaultHandle, m medium.FS) Vault {
	return h.FreshOpen(probe.Unfilled(m))
}

// vaultBlocksAfterRecoveryDefect is the planted defect of a Vault whose
// each call after the first restart blocks until its case ends.
type vaultBlocksAfterRecoveryDefect struct {
	vaultOverEmbedded
	lifetime context.Context
}

// Put is the override of the planted defect.
func (d *vaultBlocksAfterRecoveryDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	if d.lifetime != nil {
		r0 = host.Park(p0)
		return r0
	}
	return d.vaultOverEmbedded.Put(p0, p1, p2)
}

// Get is the override of the planted defect.
func (d *vaultBlocksAfterRecoveryDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.vaultOverEmbedded.Get(p0, p1)
}

// vaultBlocksAfterRecovery returns the subject of the application in place
// of the planted defect of a Vault whose each call after the first restart
// blocks until its case ends. The door Host of its check opens the defect
// in a world.
func vaultBlocksAfterRecovery(h *VaultHandle) Vault {
	return h.Fresh()
}

// vaultBlocksAfterRecoveryHost opens the planted defect of a Vault whose
// each call after the first restart blocks until its case ends. The defect
// wraps the subject that the application opens over the environment env of
// an incarnation. From the second incarnation of the host of env on, each
// callable of the defect blocks until its context ends.
func vaultBlocksAfterRecoveryHost(h *VaultHandle, env *host.Env) Vault {
	d := &vaultBlocksAfterRecoveryDefect{vaultOverEmbedded: h.FreshHost(env)}
	if env.Incarnation() > 1 {
		d.lifetime = env.Lifetime()
	}
	return d
}

// vaultGetPanicsDefect is the planted defect of a Vault whose Get panics.
type vaultGetPanicsDefect struct {
	vaultOverEmbedded
}

// Get is the override of the planted defect.
func (*vaultGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Vault whose Get panics")
}

// vaultGetPanics returns the planted defect of a Vault whose Get panics
// over a fresh subject.
func vaultGetPanics(h *VaultHandle) Vault {
	return &vaultGetPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultGetIgnoresCancelDefect is the planted defect of a Vault whose Get
// ignores the context.
type vaultGetIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Get is the override of the planted defect.
func (d *vaultGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.vaultOverEmbedded.Get(p0, p1)
}

// vaultGetIgnoresCancel returns the planted defect of a Vault whose Get
// ignores the context over a fresh subject.
func vaultGetIgnoresCancel(h *VaultHandle) Vault {
	return &vaultGetIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultGetValueBesideDefect is the planted defect of a Vault whose Get
// returns a value beside an error.
type vaultGetValueBesideDefect struct {
	vaultOverEmbedded
}

// Get is the override of the planted defect.
func (d *vaultGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	r0, r1 = d.vaultOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// vaultGetValueBeside returns the planted defect of a Vault whose Get
// returns a value beside an error over a fresh subject.
func vaultGetValueBeside(h *VaultHandle) Vault {
	return &vaultGetValueBesideDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func vaultPutCall(ctx context.Context, subject Vault, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	return nil, subject.Put(ctx, a0, a1)
}

// vaultGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func vaultGetCall(ctx context.Context, subject Vault, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// vaultReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func vaultReference(r Vault) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Put":
			values, err = vaultPutCall(ctx, r, args)
		case "Get":
			values, err = vaultGetCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// vaultWorld is the world of the subject of Vault, which the checks of its
// world probes run.
var vaultWorld = probe.World[Vault]{
	Operations: vaultOperations,
	Reference:  vaultReference,
	Corrupt:    ErrCorrupt,
}

// vaultOperations returns the operations of the world of Vault in the case
// c, over the pools of the keys that c chooses.
func vaultOperations(c *prop.Case) []probe.Operation[Vault] {
	keysString := probe.NewPool(
		c,
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []probe.Operation[Vault]{
		{
			Name: "Put",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "put.key")
				value := c.Draw(prop.Of[string](), "put.value")
				return []any{key, value}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Vault, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := vaultPutCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "get.key")
				return []any{key}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Vault, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := vaultGetCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
	}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Journal writes.
type (
	JournalHandle = conform.Handle[Journal]
	JournalCheck  = conform.Check[Journal]
	JournalOption = conform.Option[Journal]
)

// JournalChecks contains the name of each callable of Journal, with the
// identities of its checks.
var JournalChecks = struct {
	Append journalAppend
	Get    journalGet
}{
	Append: journalAppend{"Journal.Append"},
	Get:    journalGet{"Journal.Get"},
}

// journalAppend is the name of Journal.Append, with the identities of its
// checks.
type journalAppend struct {
	conform.Callable[Journal]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c journalAppend) Smoke() conform.ID[Journal] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c journalAppend) Cancel() conform.ID[Journal] {
	return c.Check("cancel")
}

// SyncedWritesSurvive returns the identity of the check that Get returns
// each value that Append wrote with success before a crash at any
// persistence point.
func (c journalAppend) SyncedWritesSurvive() conform.ID[Journal] {
	return c.Check("synced-writes-survive")
}

// NoTornRecord returns the identity of the check that Get returns each
// value of Append whole or not at all after a crash during a torn write.
func (c journalAppend) NoTornRecord() conform.ID[Journal] {
	return c.Check("no-torn-record")
}

// FullMediumIsReported returns the identity of the check that Append fails
// on a full medium, and Get returns each value as before the write.
func (c journalAppend) FullMediumIsReported() conform.ID[Journal] {
	return c.Check("full-medium-is-reported")
}

// CompletesAfterRecovery returns the identity of the check that each call
// of Journal completes within the bound of the world once every fault ends
// and the subject opens again.
func (c journalAppend) CompletesAfterRecovery() conform.ID[Journal] {
	return "Journal/completes-after-recovery"
}

// journalGet is the name of Journal.Get, with the identities of its
// checks.
type journalGet struct {
	conform.Callable[Journal]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c journalGet) Smoke() conform.ID[Journal] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c journalGet) Cancel() conform.ID[Journal] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c journalGet) ZeroOnError() conform.ID[Journal] {
	return c.Check("zero-on-error")
}

// JournalSuite returns the checks of Journal, and the calls and the
// budgets of its callables.
func JournalSuite() conform.Suite[Journal] {
	return conform.Suite[Journal]{
		Declaration: "Journal",
		Package:     "example.com/shop/vault",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Journal]{
			{
				ID:          JournalChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *JournalHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, string("key-a"), string("value-a"))
						return err
					})
				},
				Defect: journalAppendPanics,
			},
			{
				ID:          JournalChecks.Append.Smoke().Generated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, drawnKey, drawnValue)
						return err
					})
				},
				Defect: journalAppendPanics,
			},
			{
				ID:          JournalChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *JournalHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, string("key-a"), string("value-a"))
						return err
					})
				},
				Defect: journalAppendIgnoresCancel,
			},
			{
				ID:          JournalChecks.Append.Cancel().Generated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, drawnKey, drawnValue)
						return err
					})
				},
				Defect: journalAppendIgnoresCancel,
			},
			{
				ID:          JournalChecks.Append.SyncedWritesSurvive(),
				Statement:   "Get returns each value that Append wrote with success before a crash at any persistence point",
				Claim:       spec.SyncedWritesSurvive,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *JournalHandle) {
					probe.SurvivesEveryCrashPoint(h, probe.Durable[Journal, string]{
						Write: func(ctx context.Context, s Journal, i int) error {
							if i == 0 {
								_, err := s.Append(ctx, string("key-a"), string("value-a"))
								return err
							}
							_, err := s.Append(ctx, string("key-b"), string("value-b"))
							return err
						},
						Read: func(ctx context.Context, s Journal, i int) (string, error) {
							if i == 0 {
								return s.Get(ctx, string("key-a"))
							}
							return s.Get(ctx, string("key-b"))
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: journalAppendAcknowledgesBeforeSync,
				Doors:  conform.Doors[Journal]{Open: journalAppendAcknowledgesBeforeSyncOpen},
			},
			{
				ID:          JournalChecks.Append.NoTornRecord(),
				Statement:   "Get returns each value of Append whole or not at all after a crash during a torn write",
				Claim:       spec.NoTornRecord,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *JournalHandle) {
					probe.WholeAtEveryCrashPoint(h, probe.Durable[Journal, string]{
						Write: func(ctx context.Context, s Journal, i int) error {
							if i == 0 {
								_, err := s.Append(ctx, string("key-a"), string("value-a"))
								return err
							}
							_, err := s.Append(ctx, string("key-b"), string("value-b"))
							return err
						},
						Read: func(ctx context.Context, s Journal, i int) (string, error) {
							if i == 0 {
								return s.Get(ctx, string("key-a"))
							}
							return s.Get(ctx, string("key-b"))
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: journalAppendTrustsAPartialRecord,
				Doors:  conform.Doors[Journal]{Open: journalAppendTrustsAPartialRecordOpen},
			},
			{
				ID:          JournalChecks.Append.FullMediumIsReported(),
				Statement:   "Append fails on a full medium, and Get returns each value as before the write",
				Claim:       spec.FullMediumIsReported,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *JournalHandle) {
					probe.RefusesOnAFullMedium(h, probe.Durable[Journal, string]{
						Write: func(ctx context.Context, s Journal, i int) error {
							if i == 0 {
								_, err := s.Append(ctx, string("key-a"), string("value-a"))
								return err
							}
							_, err := s.Append(ctx, string("key-b"), string("value-b"))
							return err
						},
						Read: func(ctx context.Context, s Journal, i int) (string, error) {
							if i == 0 {
								return s.Get(ctx, string("key-a"))
							}
							return s.Get(ctx, string("key-b"))
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: journalAppendHidesAFullMedium,
				Doors:  conform.Doors[Journal]{Open: journalAppendHidesAFullMediumOpen},
			},
			{
				ID:        JournalChecks.Append.CompletesAfterRecovery(),
				Statement: "each call of Journal completes within the bound of the world once every fault ends and the subject opens again",
				Claim:     spec.CompletesAfterRecovery,
				Machineries: []spec.MachineryID{
					spec.Quiescence,
					spec.GeneratedSequences,
					spec.KillableProcessOrBreakableMedium,
				},
				Property: func(h *JournalHandle, c *prop.Case) {
					probe.WorldRecovery(h, c, journalWorld)
				},
				Defect:  journalBlocksAfterRecovery,
				Doors:   conform.Doors[Journal]{Host: journalBlocksAfterRecoveryHost},
				Around:  conform.Bubble,
				Options: probe.WorldOptions[Journal],
			},
			{
				ID:          JournalChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *JournalHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: journalGetPanics,
			},
			{
				ID:          JournalChecks.Get.Smoke().Generated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: journalGetPanics,
			},
			{
				ID:          JournalChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *JournalHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: journalGetIgnoresCancel,
			},
			{
				ID:          JournalChecks.Get.Cancel().Generated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: journalGetIgnoresCancel,
			},
			{
				ID:          JournalChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *JournalHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, string(""))
				},
				Defect: journalGetValueBeside,
			},
			{
				ID:          JournalChecks.Get.ZeroOnError().Generated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *JournalHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, string(""))
				},
				Defect: journalGetValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "mixin:concurrent",
				Callable:       "Journal.Append",
				Stance:         spec.Needs,
				Note:           "The claim linearizable needs the param machineries=[reference] of the directive suite.",
			},
		},
		Calls: []conform.Call[Journal]{
			{
				Calls:      journalAppendCalls,
				Over:       newJournalOver,
				Callable:   JournalChecks.Append.Callable,
				Concurrent: true,
			},
			{
				Calls:    journalGetCalls,
				Over:     newJournalOver,
				Callable: JournalChecks.Get.Callable,
			},
		},
	}
}

// RunJournalSuite applies the checks of Journal, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunJournalSuite[T Journal](t *testing.T, newSubject func() T, opts ...JournalOption) {
	t.Helper()
	RunJournalSuiteFromHandle(t, func(*JournalHandle) T { return newSubject() }, opts...)
}

// RunJournalSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunJournalSuiteFromHandle[T Journal](
	t *testing.T,
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	t.Helper()
	conform.Run(t, JournalSuite(), journalSubjects(newSubject), opts...)
}

// FuzzJournalSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzJournalSuite[T Journal](
	f *testing.F,
	callable conform.Callable[Journal],
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	f.Helper()
	conform.Fuzz(f, JournalSuite(), callable, journalSubjects(newSubject), opts...)
}

// RunJournalBudgets applies the allocation ceilings of Journal and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunJournalBudgets[T Journal](
	t *testing.T,
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	t.Helper()
	conform.RunBudgets(t, JournalSuite(), journalSubjects(newSubject), opts...)
}

// BenchJournalSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchJournalSuite[T Journal](
	b conform.B,
	callable conform.Callable[Journal],
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	b.Helper()
	conform.Bench(b, JournalSuite(), callable, journalSubjects(newSubject), opts...)
}

// BenchJournalSuiteParallel measures callable on a subject of newSubject
// on GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchJournalSuiteParallel[T Journal](
	b *testing.B,
	callable conform.Callable[Journal],
	newSubject func(h *JournalHandle) T,
	opts ...JournalOption,
) {
	b.Helper()
	conform.BenchParallel(b, JournalSuite(), callable, journalSubjects(newSubject), opts...)
}

// journalSubjects returns the factory of the subjects of Journal that
// calls newSubject with the handle of each check.
func journalSubjects[T Journal](
	newSubject func(h *JournalHandle) T,
) func(h *JournalHandle) Journal {
	return func(h *JournalHandle) Journal { return newSubject(h) }
}

// journalAppendCalls generates the calls of Journal.Append with the inputs
// that each case draws.
var journalAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Journal) error {
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Journal) error {
		_, err := subject.Append(ctx, key, value)
		return err
	}
})

// journalGetCalls generates the calls of Journal.Get with the inputs that
// each case draws.
var journalGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Journal) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Journal) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// journalOverEmbedded is Journal under the name of the field that
// journalOver embeds.
type journalOverEmbedded = Journal

// journalOver is a subject of Journal whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type journalOver struct {
	journalOverEmbedded
	spend func()
}

// Append calls spend and returns its zero results.
func (o journalOver) Append(context.Context, string, string) (r0 int64, r1 error) {
	o.spend()
	return r0, r1
}

// Get calls spend and returns its zero results.
func (o journalOver) Get(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newJournalOver returns the planted defect of a ceiling of Journal, whose
// callables call spend.
func newJournalOver(_ *JournalHandle, spend func()) Journal {
	return journalOver{spend: spend}
}

// journalAppendPanicsDefect is the planted defect of a Journal whose
// Append panics.
type journalAppendPanicsDefect struct {
	journalOverEmbedded
}

// Append is the override of the planted defect.
func (*journalAppendPanicsDefect) Append(p0 context.Context, p1, p2 string) (r0 int64, r1 error) {
	panic("the planted defect of a Journal whose Append panics")
}

// journalAppendPanics returns the planted defect of a Journal whose Append
// panics over a fresh subject.
func journalAppendPanics(h *JournalHandle) Journal {
	return &journalAppendPanicsDefect{journalOverEmbedded: h.Fresh()}
}

// journalAppendIgnoresCancelDefect is the planted defect of a Journal
// whose Append ignores the context.
type journalAppendIgnoresCancelDefect struct {
	journalOverEmbedded
}

// Append is the override of the planted defect.
func (d *journalAppendIgnoresCancelDefect) Append(
	p0 context.Context,
	p1, p2 string,
) (r0 int64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.journalOverEmbedded.Append(p0, p1, p2)
}

// journalAppendIgnoresCancel returns the planted defect of a Journal whose
// Append ignores the context over a fresh subject.
func journalAppendIgnoresCancel(h *JournalHandle) Journal {
	return &journalAppendIgnoresCancelDefect{journalOverEmbedded: h.Fresh()}
}

// journalAppendAcknowledgesBeforeSync returns the planted defect of a
// Journal over a medium whose syncs persist nothing. The door of its check
// opens the defect over a new medium in memory.
func journalAppendAcknowledgesBeforeSync(h *JournalHandle) Journal {
	return h.Open(medium.NewMemory())
}

// journalAppendAcknowledgesBeforeSyncOpen opens the planted defect of a
// Journal over a medium whose syncs persist nothing. The application opens
// the subject over the medium m that probe.Unsynced wraps.
func journalAppendAcknowledgesBeforeSyncOpen(h *JournalHandle, m medium.FS) Journal {
	return h.FreshOpen(probe.Unsynced(m))
}

// journalAppendTrustsAPartialRecordDefect is the planted defect of a
// Journal whose Append writes the zero value under the key before its
// first value.
type journalAppendTrustsAPartialRecordDefect struct {
	journalOverEmbedded
	wrote bool
}

// Append is the override of the planted defect.
func (d *journalAppendTrustsAPartialRecordDefect) Append(
	p0 context.Context,
	p1, p2 string,
) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if !d.wrote {
		d.wrote = true
		r0, r1 = d.journalOverEmbedded.Append(p0, p1, *new(string))
		if r1 != nil {
			return r0, r1
		}
	}
	return d.journalOverEmbedded.Append(p0, p1, p2)
}

// journalAppendTrustsAPartialRecord returns the planted defect of a
// Journal whose Append writes the zero value under the key before its
// first value. The door of its check opens the defect over a new medium in
// memory.
func journalAppendTrustsAPartialRecord(h *JournalHandle) Journal {
	return h.Open(medium.NewMemory())
}

// journalAppendTrustsAPartialRecordOpen opens the planted defect of a
// Journal whose Append writes the zero value under the key before its
// first value. The defect wraps a subject that the application opens over
// the medium m.
func journalAppendTrustsAPartialRecordOpen(h *JournalHandle, m medium.FS) Journal {
	return &journalAppendTrustsAPartialRecordDefect{journalOverEmbedded: h.FreshOpen(m)}
}

// journalAppendHidesAFullMedium returns the planted defect of a Journal
// over a medium that reports a full write as written. The door of its
// check opens the defect over a new medium in memory.
func journalAppendHidesAFullMedium(h *JournalHandle) Journal {
	return h.Open(medium.NewMemory())
}

// journalAppendHidesAFullMediumOpen opens the planted defect of a Journal
// over a medium that reports a full write as written. The application
// opens the subject over the medium m that probe.Unfilled wraps.
func journalAppendHidesAFullMediumOpen(h *JournalHandle, m medium.FS) Journal {
	return h.FreshOpen(probe.Unfilled(m))
}

// journalBlocksAfterRecoveryDefect is the planted defect of a Journal
// whose each call after the first restart blocks until its case ends.
type journalBlocksAfterRecoveryDefect struct {
	journalOverEmbedded
	lifetime context.Context
}

// Append is the override of the planted defect.
func (d *journalBlocksAfterRecoveryDefect) Append(
	p0 context.Context,
	p1, p2 string,
) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.journalOverEmbedded.Append(p0, p1, p2)
}

// Get is the override of the planted defect.
func (d *journalBlocksAfterRecoveryDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.journalOverEmbedded.Get(p0, p1)
}

// journalBlocksAfterRecovery returns the subject of the application in
// place of the planted defect of a Journal whose each call after the first
// restart blocks until its case ends. The door Host of its check opens the
// defect in a world.
func journalBlocksAfterRecovery(h *JournalHandle) Journal {
	return h.Fresh()
}

// journalBlocksAfterRecoveryHost opens the planted defect of a Journal
// whose each call after the first restart blocks until its case ends. The
// defect wraps the subject that the application opens over the environment
// env of an incarnation. From the second incarnation of the host of env
// on, each callable of the defect blocks until its context ends.
func journalBlocksAfterRecoveryHost(h *JournalHandle, env *host.Env) Journal {
	d := &journalBlocksAfterRecoveryDefect{journalOverEmbedded: h.FreshHost(env)}
	if env.Incarnation() > 1 {
		d.lifetime = env.Lifetime()
	}
	return d
}

// journalGetPanicsDefect is the planted defect of a Journal whose Get
// panics.
type journalGetPanicsDefect struct {
	journalOverEmbedded
}

// Get is the override of the planted defect.
func (*journalGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Journal whose Get panics")
}

// journalGetPanics returns the planted defect of a Journal whose Get
// panics over a fresh subject.
func journalGetPanics(h *JournalHandle) Journal {
	return &journalGetPanicsDefect{journalOverEmbedded: h.Fresh()}
}

// journalGetIgnoresCancelDefect is the planted defect of a Journal whose
// Get ignores the context.
type journalGetIgnoresCancelDefect struct {
	journalOverEmbedded
}

// Get is the override of the planted defect.
func (d *journalGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.journalOverEmbedded.Get(p0, p1)
}

// journalGetIgnoresCancel returns the planted defect of a Journal whose
// Get ignores the context over a fresh subject.
func journalGetIgnoresCancel(h *JournalHandle) Journal {
	return &journalGetIgnoresCancelDefect{journalOverEmbedded: h.Fresh()}
}

// journalGetValueBesideDefect is the planted defect of a Journal whose Get
// returns a value beside an error.
type journalGetValueBesideDefect struct {
	journalOverEmbedded
}

// Get is the override of the planted defect.
func (d *journalGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	r0, r1 = d.journalOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// journalGetValueBeside returns the planted defect of a Journal whose Get
// returns a value beside an error over a fresh subject.
func journalGetValueBeside(h *JournalHandle) Journal {
	return &journalGetValueBesideDefect{journalOverEmbedded: h.Fresh()}
}

// journalAppendCall calls Append on subject with the arguments args, and
// returns the values of the call beside its error.
func journalAppendCall(ctx context.Context, subject Journal, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	r0, err := subject.Append(ctx, a0, a1)
	return []any{r0}, err
}

// journalGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func journalGetCall(ctx context.Context, subject Journal, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// journalReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func journalReference(r Journal) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Append":
			values, err = journalAppendCall(ctx, r, args)
		case "Get":
			values, err = journalGetCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// journalWorld is the world of the subject of Journal, which the checks of
// its world probes run.
var journalWorld = probe.World[Journal]{
	Operations: journalOperations,
	Reference:  journalReference,
}

// journalOperations returns the operations of the world of Journal in the
// case c, over the pools of the keys that c chooses.
func journalOperations(c *prop.Case) []probe.Operation[Journal] {
	keysString := probe.NewPool(
		c,
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []probe.Operation[Journal]{
		{
			Name: "Append",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "append.key")
				value := c.Draw(prop.Of[string](), "append.value")
				return []any{key, value}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Journal, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := journalAppendCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "get.key")
				return []any{key}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Journal, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := journalGetCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
	}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Ledger writes.
type (
	LedgerHandle = conform.Handle[Ledger]
	LedgerCheck  = conform.Check[Ledger]
	LedgerOption = conform.Option[Ledger]
)

// LedgerChecks contains the name of each callable of Ledger, with the
// identities of its checks.
var LedgerChecks = struct {
	Record ledgerRecord
	Lookup ledgerLookup
}{
	Record: ledgerRecord{"Ledger.Record"},
	Lookup: ledgerLookup{"Ledger.Lookup"},
}

// ledgerRecord is the name of Ledger.Record, with the identities of its
// checks.
type ledgerRecord struct {
	conform.Callable[Ledger]
}

// Smoke returns the identity of the check that Record returns without a
// panic for derived inputs.
func (c ledgerRecord) Smoke() conform.ID[Ledger] {
	return c.Check("smoke")
}

// SyncedWritesSurvive returns the identity of the check that Lookup
// returns each value that Record wrote with success before a crash at any
// persistence point.
func (c ledgerRecord) SyncedWritesSurvive() conform.ID[Ledger] {
	return c.Check("synced-writes-survive")
}

// NoTornRecord returns the identity of the check that Lookup returns each
// value of Record whole or not at all after a crash during a torn write.
func (c ledgerRecord) NoTornRecord() conform.ID[Ledger] {
	return c.Check("no-torn-record")
}

// FullMediumIsReported returns the identity of the check that Record fails
// on a full medium, and Lookup returns each value as before the write.
func (c ledgerRecord) FullMediumIsReported() conform.ID[Ledger] {
	return c.Check("full-medium-is-reported")
}

// CompletesAfterRecovery returns the identity of the check that each call
// of Ledger completes within the bound of the world once every fault ends
// and the subject opens again.
func (c ledgerRecord) CompletesAfterRecovery() conform.ID[Ledger] {
	return "Ledger/completes-after-recovery"
}

// ledgerLookup is the name of Ledger.Lookup, with the identities of its
// checks.
type ledgerLookup struct {
	conform.Callable[Ledger]
}

// Smoke returns the identity of the check that Lookup returns without a
// panic for derived inputs.
func (c ledgerLookup) Smoke() conform.ID[Ledger] {
	return c.Check("smoke")
}

// ZeroOnMiss returns the identity of the check that Lookup returns the
// zero value beside a false presence flag.
func (c ledgerLookup) ZeroOnMiss() conform.ID[Ledger] {
	return c.Check("zero-on-miss")
}

// LedgerSuite returns the checks of Ledger, and the calls and the budgets
// of its callables.
func LedgerSuite() conform.Suite[Ledger] {
	return conform.Suite[Ledger]{
		Declaration: "Ledger",
		Package:     "example.com/shop/vault",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Ledger]{
			{
				ID:          LedgerChecks.Record.Smoke(),
				Statement:   "Record returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LedgerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Record(string("key-a"), string("value-a"))
					})
				},
				Defect: ledgerRecordPanics,
			},
			{
				ID:          LedgerChecks.Record.Smoke().Generated(),
				Statement:   "Record returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LedgerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Record(drawnKey, drawnValue)
					})
				},
				Defect: ledgerRecordPanics,
			},
			{
				ID:          LedgerChecks.Record.SyncedWritesSurvive(),
				Statement:   "Lookup returns each value that Record wrote with success before a crash at any persistence point",
				Claim:       spec.SyncedWritesSurvive,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *LedgerHandle) {
					probe.SurvivesEveryCrashPoint(h, probe.Durable[Ledger, string]{
						Write: func(ctx context.Context, s Ledger, i int) error {
							if i == 0 {
								return s.Record(string("key-a"), string("value-a"))
							}
							return s.Record(string("key-b"), string("value-b"))
						},
						Read: func(ctx context.Context, s Ledger, i int) (string, error) {
							if i == 0 {
								r0, ok := s.Lookup(string("key-a"))
								if !ok {
									return r0, fs.ErrNotExist
								}
								return r0, nil
							}
							r0, ok := s.Lookup(string("key-b"))
							if !ok {
								return r0, fs.ErrNotExist
							}
							return r0, nil
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: ledgerRecordAcknowledgesBeforeSync,
				Doors:  conform.Doors[Ledger]{Open: ledgerRecordAcknowledgesBeforeSyncOpen},
			},
			{
				ID:          LedgerChecks.Record.NoTornRecord(),
				Statement:   "Lookup returns each value of Record whole or not at all after a crash during a torn write",
				Claim:       spec.NoTornRecord,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *LedgerHandle) {
					probe.WholeAtEveryCrashPoint(h, probe.Durable[Ledger, string]{
						Write: func(ctx context.Context, s Ledger, i int) error {
							if i == 0 {
								return s.Record(string("key-a"), string("value-a"))
							}
							return s.Record(string("key-b"), string("value-b"))
						},
						Read: func(ctx context.Context, s Ledger, i int) (string, error) {
							if i == 0 {
								r0, ok := s.Lookup(string("key-a"))
								if !ok {
									return r0, fs.ErrNotExist
								}
								return r0, nil
							}
							r0, ok := s.Lookup(string("key-b"))
							if !ok {
								return r0, fs.ErrNotExist
							}
							return r0, nil
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: ledgerRecordTrustsAPartialRecord,
				Doors:  conform.Doors[Ledger]{Open: ledgerRecordTrustsAPartialRecordOpen},
			},
			{
				ID:          LedgerChecks.Record.FullMediumIsReported(),
				Statement:   "Record fails on a full medium, and Lookup returns each value as before the write",
				Claim:       spec.FullMediumIsReported,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *LedgerHandle) {
					probe.RefusesOnAFullMedium(h, probe.Durable[Ledger, string]{
						Write: func(ctx context.Context, s Ledger, i int) error {
							if i == 0 {
								return s.Record(string("key-a"), string("value-a"))
							}
							return s.Record(string("key-b"), string("value-b"))
						},
						Read: func(ctx context.Context, s Ledger, i int) (string, error) {
							if i == 0 {
								r0, ok := s.Lookup(string("key-a"))
								if !ok {
									return r0, fs.ErrNotExist
								}
								return r0, nil
							}
							r0, ok := s.Lookup(string("key-b"))
							if !ok {
								return r0, fs.ErrNotExist
							}
							return r0, nil
						},
						Want: func(i int) string {
							if i == 0 {
								return string("value-a")
							}
							return string("value-b")
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{string("key-a")}
							}
							return []any{string("key-b")}
						},
					})
				},
				Defect: ledgerRecordHidesAFullMedium,
				Doors:  conform.Doors[Ledger]{Open: ledgerRecordHidesAFullMediumOpen},
			},
			{
				ID:        LedgerChecks.Record.CompletesAfterRecovery(),
				Statement: "each call of Ledger completes within the bound of the world once every fault ends and the subject opens again",
				Claim:     spec.CompletesAfterRecovery,
				Machineries: []spec.MachineryID{
					spec.Quiescence,
					spec.GeneratedSequences,
					spec.KillableProcessOrBreakableMedium,
				},
				Property: func(h *LedgerHandle, c *prop.Case) {
					probe.WorldRecovery(h, c, ledgerWorld)
				},
				Defect:  ledgerBlocksAfterRecovery,
				Doors:   conform.Doors[Ledger]{Host: ledgerBlocksAfterRecoveryHost},
				Around:  conform.Bubble,
				Options: probe.WorldOptions[Ledger],
			},
			{
				ID:          LedgerChecks.Lookup.Smoke(),
				Statement:   "Lookup returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LedgerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Lookup(string("key-a"))
						return nil
					})
				},
				Defect: ledgerLookupPanics,
			},
			{
				ID:          LedgerChecks.Lookup.Smoke().Generated(),
				Statement:   "Lookup returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LedgerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Lookup(drawnKey)
						return nil
					})
				},
				Defect: ledgerLookupPanics,
			},
			{
				ID:          LedgerChecks.Lookup.ZeroOnMiss(),
				Statement:   "Lookup returns the zero value beside a false presence flag",
				Claim:       spec.ZeroOnMiss,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LedgerHandle) {
					probe.ZeroBesideFalse(h, func(ctx context.Context) (string, bool) {
						return h.Subject().Lookup(string("key-a"))
					}, string(""))
				},
				Defect: ledgerLookupValueBesideAFalseFlag,
			},
			{
				ID:          LedgerChecks.Lookup.ZeroOnMiss().Generated(),
				Statement:   "Lookup returns the zero value beside a false presence flag",
				Claim:       spec.ZeroOnMiss,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LedgerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideFalse(h, func(ctx context.Context) (string, bool) {
						return h.Subject().Lookup(drawnKey)
					}, string(""))
				},
				Defect: ledgerLookupValueBesideAFalseFlag,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Ledger.Record",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Ledger]{
			{
				Calls:    ledgerRecordCalls,
				Over:     newLedgerOver,
				Callable: LedgerChecks.Record.Callable,
			},
			{
				Calls:    ledgerLookupCalls,
				Over:     newLedgerOver,
				Callable: LedgerChecks.Lookup.Callable,
			},
		},
	}
}

// RunLedgerSuite applies the checks of Ledger, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLedgerSuite[T Ledger](t *testing.T, newSubject func() T, opts ...LedgerOption) {
	t.Helper()
	RunLedgerSuiteFromHandle(t, func(*LedgerHandle) T { return newSubject() }, opts...)
}

// RunLedgerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLedgerSuiteFromHandle[T Ledger](
	t *testing.T,
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	t.Helper()
	conform.Run(t, LedgerSuite(), ledgerSubjects(newSubject), opts...)
}

// FuzzLedgerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLedgerSuite[T Ledger](
	f *testing.F,
	callable conform.Callable[Ledger],
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	f.Helper()
	conform.Fuzz(f, LedgerSuite(), callable, ledgerSubjects(newSubject), opts...)
}

// RunLedgerBudgets applies the allocation ceilings of Ledger and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLedgerBudgets[T Ledger](
	t *testing.T,
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	t.Helper()
	conform.RunBudgets(t, LedgerSuite(), ledgerSubjects(newSubject), opts...)
}

// BenchLedgerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLedgerSuite[T Ledger](
	b conform.B,
	callable conform.Callable[Ledger],
	newSubject func(h *LedgerHandle) T,
	opts ...LedgerOption,
) {
	b.Helper()
	conform.Bench(b, LedgerSuite(), callable, ledgerSubjects(newSubject), opts...)
}

// ledgerSubjects returns the factory of the subjects of Ledger that calls
// newSubject with the handle of each check.
func ledgerSubjects[T Ledger](newSubject func(h *LedgerHandle) T) func(h *LedgerHandle) Ledger {
	return func(h *LedgerHandle) Ledger { return newSubject(h) }
}

// ledgerRecordCalls generates the calls of Ledger.Record with the inputs
// that each case draws.
var ledgerRecordCalls = prop.Composite(func(c *prop.Case) func(context.Context, Ledger) error {
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(_ context.Context, subject Ledger) error {
		return subject.Record(key, value)
	}
})

// ledgerLookupCalls generates the calls of Ledger.Lookup with the inputs
// that each case draws.
var ledgerLookupCalls = prop.Composite(func(c *prop.Case) func(context.Context, Ledger) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(_ context.Context, subject Ledger) error {
		_, _ = subject.Lookup(key)
		return nil
	}
})

// ledgerOverEmbedded is Ledger under the name of the field that ledgerOver
// embeds.
type ledgerOverEmbedded = Ledger

// ledgerOver is a subject of Ledger whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type ledgerOver struct {
	ledgerOverEmbedded
	spend func()
}

// Record calls spend and returns its zero results.
func (o ledgerOver) Record(string, string) (r0 error) {
	o.spend()
	return r0
}

// Lookup calls spend and returns its zero results.
func (o ledgerOver) Lookup(string) (r0 string, r1 bool) {
	o.spend()
	return r0, r1
}

// newLedgerOver returns the planted defect of a ceiling of Ledger, whose
// callables call spend.
func newLedgerOver(_ *LedgerHandle, spend func()) Ledger {
	return ledgerOver{spend: spend}
}

// ledgerRecordPanicsDefect is the planted defect of a Ledger whose Record
// panics.
type ledgerRecordPanicsDefect struct {
	ledgerOverEmbedded
}

// Record is the override of the planted defect.
func (*ledgerRecordPanicsDefect) Record(p0, p1 string) (r0 error) {
	panic("the planted defect of a Ledger whose Record panics")
}

// ledgerRecordPanics returns the planted defect of a Ledger whose Record
// panics over a fresh subject.
func ledgerRecordPanics(h *LedgerHandle) Ledger {
	return &ledgerRecordPanicsDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerRecordAcknowledgesBeforeSync returns the planted defect of a
// Ledger over a medium whose syncs persist nothing. The door of its check
// opens the defect over a new medium in memory.
func ledgerRecordAcknowledgesBeforeSync(h *LedgerHandle) Ledger {
	return h.Open(medium.NewMemory())
}

// ledgerRecordAcknowledgesBeforeSyncOpen opens the planted defect of a
// Ledger over a medium whose syncs persist nothing. The application opens
// the subject over the medium m that probe.Unsynced wraps.
func ledgerRecordAcknowledgesBeforeSyncOpen(h *LedgerHandle, m medium.FS) Ledger {
	return h.FreshOpen(probe.Unsynced(m))
}

// ledgerRecordTrustsAPartialRecordDefect is the planted defect of a Ledger
// whose Record writes the zero value under the key before its first value.
type ledgerRecordTrustsAPartialRecordDefect struct {
	ledgerOverEmbedded
	wrote bool
}

// Record is the override of the planted defect.
func (d *ledgerRecordTrustsAPartialRecordDefect) Record(p0, p1 string) (r0 error) {
	if !d.wrote {
		d.wrote = true
		r0 = d.ledgerOverEmbedded.Record(p0, *new(string))
		if r0 != nil {
			return r0
		}
	}
	return d.ledgerOverEmbedded.Record(p0, p1)
}

// ledgerRecordTrustsAPartialRecord returns the planted defect of a Ledger
// whose Record writes the zero value under the key before its first value.
// The door of its check opens the defect over a new medium in memory.
func ledgerRecordTrustsAPartialRecord(h *LedgerHandle) Ledger {
	return h.Open(medium.NewMemory())
}

// ledgerRecordTrustsAPartialRecordOpen opens the planted defect of a
// Ledger whose Record writes the zero value under the key before its first
// value. The defect wraps a subject that the application opens over the
// medium m.
func ledgerRecordTrustsAPartialRecordOpen(h *LedgerHandle, m medium.FS) Ledger {
	return &ledgerRecordTrustsAPartialRecordDefect{ledgerOverEmbedded: h.FreshOpen(m)}
}

// ledgerRecordHidesAFullMedium returns the planted defect of a Ledger over
// a medium that reports a full write as written. The door of its check
// opens the defect over a new medium in memory.
func ledgerRecordHidesAFullMedium(h *LedgerHandle) Ledger {
	return h.Open(medium.NewMemory())
}

// ledgerRecordHidesAFullMediumOpen opens the planted defect of a Ledger
// over a medium that reports a full write as written. The application
// opens the subject over the medium m that probe.Unfilled wraps.
func ledgerRecordHidesAFullMediumOpen(h *LedgerHandle, m medium.FS) Ledger {
	return h.FreshOpen(probe.Unfilled(m))
}

// ledgerBlocksAfterRecoveryDefect is the planted defect of a Ledger whose
// each call after the first restart blocks until its case ends.
type ledgerBlocksAfterRecoveryDefect struct {
	ledgerOverEmbedded
	lifetime context.Context
}

// Record is the override of the planted defect.
func (d *ledgerBlocksAfterRecoveryDefect) Record(p0, p1 string) (r0 error) {
	if d.lifetime != nil {
		r0 = host.Park(d.lifetime)
		return r0
	}
	return d.ledgerOverEmbedded.Record(p0, p1)
}

// Lookup is the override of the planted defect.
func (d *ledgerBlocksAfterRecoveryDefect) Lookup(p0 string) (r0 string, r1 bool) {
	if d.lifetime != nil {
		_ = host.Park(d.lifetime)
		return r0, r1
	}
	return d.ledgerOverEmbedded.Lookup(p0)
}

// ledgerBlocksAfterRecovery returns the subject of the application in
// place of the planted defect of a Ledger whose each call after the first
// restart blocks until its case ends. The door Host of its check opens the
// defect in a world.
func ledgerBlocksAfterRecovery(h *LedgerHandle) Ledger {
	return h.Fresh()
}

// ledgerBlocksAfterRecoveryHost opens the planted defect of a Ledger whose
// each call after the first restart blocks until its case ends. The defect
// wraps the subject that the application opens over the environment env of
// an incarnation. From the second incarnation of the host of env on, each
// callable of the defect blocks until its context ends.
func ledgerBlocksAfterRecoveryHost(h *LedgerHandle, env *host.Env) Ledger {
	d := &ledgerBlocksAfterRecoveryDefect{ledgerOverEmbedded: h.FreshHost(env)}
	if env.Incarnation() > 1 {
		d.lifetime = env.Lifetime()
	}
	return d
}

// ledgerLookupPanicsDefect is the planted defect of a Ledger whose Lookup
// panics.
type ledgerLookupPanicsDefect struct {
	ledgerOverEmbedded
}

// Lookup is the override of the planted defect.
func (*ledgerLookupPanicsDefect) Lookup(p0 string) (r0 string, r1 bool) {
	panic("the planted defect of a Ledger whose Lookup panics")
}

// ledgerLookupPanics returns the planted defect of a Ledger whose Lookup
// panics over a fresh subject.
func ledgerLookupPanics(h *LedgerHandle) Ledger {
	return &ledgerLookupPanicsDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerLookupValueBesideAFalseFlagDefect is the planted defect of a
// Ledger whose Lookup returns a value beside a false presence flag.
type ledgerLookupValueBesideAFalseFlagDefect struct {
	ledgerOverEmbedded
}

// Lookup is the override of the planted defect.
func (d *ledgerLookupValueBesideAFalseFlagDefect) Lookup(p0 string) (r0 string, r1 bool) {
	r0, r1 = d.ledgerOverEmbedded.Lookup(p0)
	if !r1 {
		r0 = string("sample-b")
	}
	return r0, r1
}

// ledgerLookupValueBesideAFalseFlag returns the planted defect of a Ledger
// whose Lookup returns a value beside a false presence flag over a fresh
// subject.
func ledgerLookupValueBesideAFalseFlag(h *LedgerHandle) Ledger {
	return &ledgerLookupValueBesideAFalseFlagDefect{ledgerOverEmbedded: h.Fresh()}
}

// ledgerRecordCall calls Record on subject with the arguments args, and
// returns the values of the call beside its error.
func ledgerRecordCall(ctx context.Context, subject Ledger, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	return nil, subject.Record(a0, a1)
}

// ledgerLookupCall calls Lookup on subject with the arguments args, and
// returns the values of the call beside its error.
func ledgerLookupCall(ctx context.Context, subject Ledger, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, ok := subject.Lookup(a0)
	return []any{r0, ok}, nil
}

// ledgerReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func ledgerReference(r Ledger) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Record":
			values, err = ledgerRecordCall(ctx, r, args)
		case "Lookup":
			values, err = ledgerLookupCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// ledgerWorld is the world of the subject of Ledger, which the checks of
// its world probes run.
var ledgerWorld = probe.World[Ledger]{
	Operations: ledgerOperations,
	Reference:  ledgerReference,
}

// ledgerOperations returns the operations of the world of Ledger in the
// case c, over the pools of the keys that c chooses.
func ledgerOperations(c *prop.Case) []probe.Operation[Ledger] {
	keysString := probe.NewPool(
		c,
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []probe.Operation[Ledger]{
		{
			Name: "Record",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "record.key")
				value := c.Draw(prop.Of[string](), "record.value")
				return []any{key, value}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Ledger, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := ledgerRecordCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
		{
			Name: "Lookup",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "lookup.key")
				return []any{key}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Ledger, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := ledgerLookupCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
	}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Registry writes.
type (
	RegistryHandle = conform.Handle[Registry]
	RegistryCheck  = conform.Check[Registry]
	RegistryOption = conform.Option[Registry]
)

// RegistryChecks contains the name of each callable of Registry, with the
// identities of its checks.
var RegistryChecks = struct {
	Set    registrySet
	Value  registryValue
	Await  registryAwait
	Add    registryAdd
	Insert registryInsert
	Stamp  registryStamp
	Touch  registryTouch
}{
	Set:    registrySet{"Registry.Set"},
	Value:  registryValue{"Registry.Value"},
	Await:  registryAwait{"Registry.Await"},
	Add:    registryAdd{"Registry.Add"},
	Insert: registryInsert{"Registry.Insert"},
	Stamp:  registryStamp{"Registry.Stamp"},
	Touch:  registryTouch{"Registry.Touch"},
}

// registrySet is the name of Registry.Set, with the identities of its
// checks.
type registrySet struct {
	conform.Callable[Registry]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c registrySet) Smoke() conform.ID[Registry] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Set returns an error for a
// cancelled context.
func (c registrySet) Cancel() conform.ID[Registry] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Set returns an
// error when a call of its dependency fails.
func (c registrySet) FailureIsReported() conform.ID[Registry] {
	return c.Check("failure-is-reported")
}

// CompletesAfterRecovery returns the identity of the check that each call
// of Registry completes within the bound of the world once every fault
// ends and the subject opens again.
func (c registrySet) CompletesAfterRecovery() conform.ID[Registry] {
	return "Registry/completes-after-recovery"
}

// registryValue is the name of Registry.Value, with the identities of its
// checks.
type registryValue struct {
	conform.Callable[Registry]
}

// Smoke returns the identity of the check that Value returns without a
// panic for derived inputs.
func (c registryValue) Smoke() conform.ID[Registry] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Value returns an error for
// a cancelled context.
func (c registryValue) Cancel() conform.ID[Registry] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Value returns
// an error when a call of its dependency fails.
func (c registryValue) FailureIsReported() conform.ID[Registry] {
	return c.Check("failure-is-reported")
}

// ZeroOnError returns the identity of the check that Value returns the
// zero value beside an error.
func (c registryValue) ZeroOnError() conform.ID[Registry] {
	return c.Check("zero-on-error")
}

// AgreesWithTheReference returns the identity of the check that every
// sequence of calls returns on Registry what it returns on the reference.
func (c registryValue) AgreesWithTheReference() conform.ID[Registry] {
	return "Registry/agrees-with-the-reference"
}

// Linearizable returns the identity of the check that every concurrent
// history of Registry has an order in which the reference returns the same
// results.
func (c registryValue) Linearizable() conform.ID[Registry] {
	return "Registry/linearizable"
}

// LinearizableUnderFaults returns the identity of the check that the
// history of the calls of Registry through crashes and faults has a
// linearization that the reference accepts.
func (c registryValue) LinearizableUnderFaults() conform.ID[Registry] {
	return "Registry/linearizable-under-faults"
}

// registryAwait is the name of Registry.Await, with the identities of its
// checks.
type registryAwait struct {
	conform.Callable[Registry]
}

// Cancel returns the identity of the check that Await returns an error for
// a cancelled context.
func (c registryAwait) Cancel() conform.ID[Registry] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Await returns
// an error when a call of its dependency fails.
func (c registryAwait) FailureIsReported() conform.ID[Registry] {
	return c.Check("failure-is-reported")
}

// registryAdd is the name of Registry.Add, with the identities of its
// checks.
type registryAdd struct {
	conform.Callable[Registry]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c registryAdd) Smoke() conform.ID[Registry] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c registryAdd) Cancel() conform.ID[Registry] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Add returns an
// error when a call of its dependency fails.
func (c registryAdd) FailureIsReported() conform.ID[Registry] {
	return c.Check("failure-is-reported")
}

// ZeroOnError returns the identity of the check that Add returns the zero
// value beside an error.
func (c registryAdd) ZeroOnError() conform.ID[Registry] {
	return c.Check("zero-on-error")
}

// registryInsert is the name of Registry.Insert, with the identities of
// its checks.
type registryInsert struct {
	conform.Callable[Registry]
}

// Smoke returns the identity of the check that Insert returns without a
// panic for derived inputs.
func (c registryInsert) Smoke() conform.ID[Registry] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Insert returns an error
// for a cancelled context.
func (c registryInsert) Cancel() conform.ID[Registry] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Insert returns
// an error when a call of its dependency fails.
func (c registryInsert) FailureIsReported() conform.ID[Registry] {
	return c.Check("failure-is-reported")
}

// ZeroOnError returns the identity of the check that Insert returns the
// zero value beside an error.
func (c registryInsert) ZeroOnError() conform.ID[Registry] {
	return c.Check("zero-on-error")
}

// registryStamp is the name of Registry.Stamp, with the identities of its
// checks.
type registryStamp struct {
	conform.Callable[Registry]
}

// Smoke returns the identity of the check that Stamp returns without a
// panic for derived inputs.
func (c registryStamp) Smoke() conform.ID[Registry] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Stamp returns an error for
// a cancelled context.
func (c registryStamp) Cancel() conform.ID[Registry] {
	return c.Check("cancel")
}

// FailureIsReported returns the identity of the check that Stamp returns
// an error when a call of its dependency fails.
func (c registryStamp) FailureIsReported() conform.ID[Registry] {
	return c.Check("failure-is-reported")
}

// ZeroOnError returns the identity of the check that Stamp returns the
// zero value beside an error.
func (c registryStamp) ZeroOnError() conform.ID[Registry] {
	return c.Check("zero-on-error")
}

// registryTouch is the name of Registry.Touch, with the identities of its
// checks.
type registryTouch struct {
	conform.Callable[Registry]
}

// Smoke returns the identity of the check that Touch returns without a
// panic for derived inputs.
func (c registryTouch) Smoke() conform.ID[Registry] {
	return c.Check("smoke")
}

// RegistrySuite returns the checks of Registry, and the calls and the
// budgets of its callables.
func RegistrySuite() conform.Suite[Registry] {
	return conform.Suite[Registry]{
		Declaration: "Registry",
		Package:     "example.com/shop/vault",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Registry]{
			{
				ID:          RegistryChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: registrySetPanics,
			},
			{
				ID:          RegistryChecks.Set.Smoke().Generated(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnKey, drawnValue)
					})
				},
				Defect: registrySetPanics,
			},
			{
				ID:          RegistryChecks.Set.Cancel(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: registrySetIgnoresCancel,
			},
			{
				ID:          RegistryChecks.Set.Cancel().Generated(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnKey, drawnValue)
					})
				},
				Defect: registrySetIgnoresCancel,
			},
			{
				ID:        RegistryChecks.Set.FailureIsReported(),
				Statement: "Set returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *RegistryHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, string("key-a"), string("value-a"))
					})
				},
				Defect: registrySetSwallowsFailure,
			},
			{
				ID:        RegistryChecks.Set.FailureIsReported().Generated(),
				Statement: "Set returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnKey, drawnValue)
					})
				},
				Defect: registrySetSwallowsFailure,
			},
			{
				ID:        RegistryChecks.Set.CompletesAfterRecovery(),
				Statement: "each call of Registry completes within the bound of the world once every fault ends and the subject opens again",
				Claim:     spec.CompletesAfterRecovery,
				Machineries: []spec.MachineryID{
					spec.Quiescence,
					spec.GeneratedSequences,
					spec.KillableProcessOrBreakableMedium,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					probe.WorldRecovery(h, c, registryWorld)
				},
				Defect:  registryBlocksAfterRecovery,
				Doors:   conform.Doors[Registry]{Host: registryBlocksAfterRecoveryHost},
				Around:  conform.Bubble,
				Options: probe.WorldOptions[Registry],
			},
			{
				ID:          RegistryChecks.Value.Smoke(),
				Statement:   "Value returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx, string("key-a"))
						return err
					})
				},
				Defect: registryValuePanics,
			},
			{
				ID:          RegistryChecks.Value.Smoke().Generated(),
				Statement:   "Value returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx, drawnKey)
						return err
					})
				},
				Defect: registryValuePanics,
			},
			{
				ID:          RegistryChecks.Value.Cancel(),
				Statement:   "Value returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx, string("key-a"))
						return err
					})
				},
				Defect: registryValueIgnoresCancel,
			},
			{
				ID:          RegistryChecks.Value.Cancel().Generated(),
				Statement:   "Value returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx, drawnKey)
						return err
					})
				},
				Defect: registryValueIgnoresCancel,
			},
			{
				ID:        RegistryChecks.Value.FailureIsReported(),
				Statement: "Value returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *RegistryHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx, string("key-a"))
						return err
					})
				},
				Defect: registryValueSwallowsFailure,
			},
			{
				ID:        RegistryChecks.Value.FailureIsReported().Generated(),
				Statement: "Value returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx, drawnKey)
						return err
					})
				},
				Defect: registryValueSwallowsFailure,
			},
			{
				ID:          RegistryChecks.Value.ZeroOnError(),
				Statement:   "Value returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Value(ctx, string("key-a"))
					}, string(""))
				},
				Defect: registryValueValueBeside,
			},
			{
				ID:          RegistryChecks.Value.ZeroOnError().Generated(),
				Statement:   "Value returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Value(ctx, drawnKey)
					}, string(""))
				},
				Defect: registryValueValueBeside,
			},
			{
				ID:        RegistryChecks.Value.AgreesWithTheReference(),
				Statement: "every sequence of calls returns on Registry what it returns on the reference",
				Claim:     spec.AgreesWithTheReference,
				Machineries: []spec.MachineryID{
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					probe.SequencesAgreeWithTheReference(h, c, registryActions, registryReference)
				},
				Defect: registryValueDivergesAfterThree,
			},
			{
				ID:        RegistryChecks.Value.Linearizable(),
				Statement: "every concurrent history of Registry has an order in which the reference returns the same results",
				Claim:     spec.Linearizable,
				Machineries: []spec.MachineryID{
					spec.ConcurrentCallers,
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					probe.LinearizableHistory(
						h,
						c,
						registryActions,
						registryReference,
						2,
						history.Budget(200000),
					)
				},
				Defect: registryRacyReadWrite,
			},
			{
				ID:        RegistryChecks.Value.LinearizableUnderFaults(),
				Statement: "the history of the calls of Registry through crashes and faults has a linearization that the reference accepts",
				Claim:     spec.LinearizableUnderFaults,
				Machineries: []spec.MachineryID{
					spec.Quiescence,
					spec.ConcurrentCallers,
					spec.GeneratedSequences,
					spec.KillableProcessOrBreakableMedium,
					spec.CallThatFailsOnDemand,
					spec.ReferenceImplementation,
					spec.ClockTheRunAdvances,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					probe.FaultyWorldHistory(h, c, registryWorld, history.Budget(200000))
				},
				Defect:  registryAmnesiacRecovery,
				Doors:   conform.Doors[Registry]{Host: registryAmnesiacRecoveryHost},
				Around:  conform.Bubble,
				Options: probe.WorldOptions[Registry],
			},
			{
				ID:          RegistryChecks.Await.Cancel(),
				Statement:   "Await returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Await(ctx, string("key-a"))
						return err
					})
				},
				Defect: registryAwaitIgnoresCancel,
			},
			{
				ID:          RegistryChecks.Await.Cancel().Generated(),
				Statement:   "Await returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Await(ctx, drawnKey)
						return err
					})
				},
				Defect: registryAwaitIgnoresCancel,
			},
			{
				ID:        RegistryChecks.Await.FailureIsReported(),
				Statement: "Await returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *RegistryHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Await(ctx, string("key-a"))
						return err
					})
				},
				Defect: registryAwaitSwallowsFailure,
			},
			{
				ID:        RegistryChecks.Await.FailureIsReported().Generated(),
				Statement: "Await returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Await(ctx, drawnKey)
						return err
					})
				},
				Defect: registryAwaitSwallowsFailure,
			},
			{
				ID:          RegistryChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Add(ctx, string("value-a"))
						return err
					})
				},
				Defect: registryAddPanics,
			},
			{
				ID:          RegistryChecks.Add.Smoke().Generated(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Add(ctx, drawnValue)
						return err
					})
				},
				Defect: registryAddPanics,
			},
			{
				ID:          RegistryChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Add(ctx, string("value-a"))
						return err
					})
				},
				Defect: registryAddIgnoresCancel,
			},
			{
				ID:          RegistryChecks.Add.Cancel().Generated(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Add(ctx, drawnValue)
						return err
					})
				},
				Defect: registryAddIgnoresCancel,
			},
			{
				ID:        RegistryChecks.Add.FailureIsReported(),
				Statement: "Add returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *RegistryHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Add(ctx, string("value-a"))
						return err
					})
				},
				Defect: registryAddSwallowsFailure,
			},
			{
				ID:        RegistryChecks.Add.FailureIsReported().Generated(),
				Statement: "Add returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Add(ctx, drawnValue)
						return err
					})
				},
				Defect: registryAddSwallowsFailure,
			},
			{
				ID:          RegistryChecks.Add.ZeroOnError(),
				Statement:   "Add returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Add(ctx, string("value-a"))
					}, string(""))
				},
				Defect: registryAddValueBeside,
			},
			{
				ID:          RegistryChecks.Add.ZeroOnError().Generated(),
				Statement:   "Add returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Add(ctx, drawnValue)
					}, string(""))
				},
				Defect: registryAddValueBeside,
			},
			{
				ID:          RegistryChecks.Insert.Smoke(),
				Statement:   "Insert returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Insert(ctx, string("value-a"))
						return err
					})
				},
				Defect: registryInsertPanics,
			},
			{
				ID:          RegistryChecks.Insert.Smoke().Generated(),
				Statement:   "Insert returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Insert(ctx, drawnValue)
						return err
					})
				},
				Defect: registryInsertPanics,
			},
			{
				ID:          RegistryChecks.Insert.Cancel(),
				Statement:   "Insert returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Insert(ctx, string("value-a"))
						return err
					})
				},
				Defect: registryInsertIgnoresCancel,
			},
			{
				ID:          RegistryChecks.Insert.Cancel().Generated(),
				Statement:   "Insert returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Insert(ctx, drawnValue)
						return err
					})
				},
				Defect: registryInsertIgnoresCancel,
			},
			{
				ID:        RegistryChecks.Insert.FailureIsReported(),
				Statement: "Insert returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *RegistryHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Insert(ctx, string("value-a"))
						return err
					})
				},
				Defect: registryInsertSwallowsFailure,
			},
			{
				ID:        RegistryChecks.Insert.FailureIsReported().Generated(),
				Statement: "Insert returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.GeneratedInputs,
					spec.NothingBeyondACaller,
				},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Insert(ctx, drawnValue)
						return err
					})
				},
				Defect: registryInsertSwallowsFailure,
			},
			{
				ID:          RegistryChecks.Insert.ZeroOnError(),
				Statement:   "Insert returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (*Entry, error) {
						return h.Subject().Insert(ctx, string("value-a"))
					}, (*Entry)(nil))
				},
				Defect: registryInsertValueBeside,
			},
			{
				ID:          RegistryChecks.Insert.ZeroOnError().Generated(),
				Statement:   "Insert returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegistryHandle, c *prop.Case) {
					drawnValue := c.Draw(prop.Of[string](), "value")
					probe.ZeroBesideError(h, func(ctx context.Context) (*Entry, error) {
						return h.Subject().Insert(ctx, drawnValue)
					}, (*Entry)(nil))
				},
				Defect: registryInsertValueBeside,
			},
			{
				ID:          RegistryChecks.Stamp.Smoke(),
				Statement:   "Stamp returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Stamp(ctx)
						return err
					})
				},
				Defect: registryStampPanics,
			},
			{
				ID:          RegistryChecks.Stamp.Cancel(),
				Statement:   "Stamp returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Stamp(ctx)
						return err
					})
				},
				Defect: registryStampIgnoresCancel,
			},
			{
				ID:        RegistryChecks.Stamp.FailureIsReported(),
				Statement: "Stamp returns an error when a call of its dependency fails",
				Claim:     spec.FailureIsReported,
				Machineries: []spec.MachineryID{
					spec.CallThatFailsOnDemand,
					spec.NothingBeyondACaller,
				},
				Apply: func(h *RegistryHandle) {
					probe.ReportsAFailedDependency(h, func(ctx context.Context) error {
						_, err := h.Subject().Stamp(ctx)
						return err
					})
				},
				Defect: registryStampSwallowsFailure,
			},
			{
				ID:          RegistryChecks.Stamp.ZeroOnError(),
				Statement:   "Stamp returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistryHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Stamp(ctx)
					}, int64(0))
				},
				Defect: registryStampValueBeside,
			},
			{
				ID:          RegistryChecks.Touch.Smoke(),
				Statement:   "Touch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistryHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Touch(ctx)
						return nil
					})
				},
				Defect: registryTouchPanics,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.SyncedWritesSurvive,
				Callable: "Registry.Set",
				Gate:     spec.OnePartner,
				Why:      "two callables of the declaration can be the partner read of Set",
				Remedy:   "name the partner of Set with the param of the classification",
			},
			{
				Claim:    spec.NoTornRecord,
				Callable: "Registry.Set",
				Gate:     spec.OnePartner,
				Why:      "two callables of the declaration can be the partner read of Set",
				Remedy:   "name the partner of Set with the param of the classification",
			},
			{
				Claim:    spec.CorruptionIsReported,
				Callable: "Registry.Set",
				Gate:     spec.OnePartner,
				Why:      "two callables of the declaration can be the partner read of Set",
				Remedy:   "name the partner of Set with the param of the classification",
			},
			{
				Claim:    spec.FullMediumIsReported,
				Callable: "Registry.Set",
				Gate:     spec.OnePartner,
				Why:      "two callables of the declaration can be the partner read of Set",
				Remedy:   "name the partner of Set with the param of the classification",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Registry.Await",
				Gate:     spec.ReturnsUnaided,
				Why:      "Await waits for the call of another callable on a fresh subject",
				Remedy:   "remove the mixin waits where Await returns on a fresh subject",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Registry.Touch",
				Gate:     spec.Failure,
				Why:      "Touch cannot fail",
				Remedy:   "give Touch an error result",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Registry.Set",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "mixin:waits",
				Callable:       "Registry.Await",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable can wait for another call.",
			},
			{
				Classification: "mixin:creates",
				Callable:       "Registry.Add",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable creates a record and returns the key of the record.",
			},
			{
				Classification: "mixin:creates",
				Callable:       "Registry.Insert",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable creates a record and returns the key of the record.",
			},
			{
				Classification: "mixin:creates",
				Callable:       "Registry.Stamp",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable creates a record and returns the key of the record.",
			},
			{
				Classification: "mixin:creates",
				Callable:       "Registry.Touch",
				Stance:         spec.Silent,
				Note:           "The mixin documents that the callable creates a record and returns the key of the record.",
			},
		},
		Calls: []conform.Call[Registry]{
			{
				Calls:      registrySetCalls,
				Over:       newRegistryOver,
				Callable:   RegistryChecks.Set.Callable,
				Concurrent: true,
			},
			{
				Calls:      registryValueCalls,
				Over:       newRegistryOver,
				Callable:   RegistryChecks.Value.Callable,
				Concurrent: true,
			},
			{
				Calls:    registryAwaitCalls,
				Over:     newRegistryOver,
				Callable: RegistryChecks.Await.Callable,
				Waits:    true,
			},
			{
				Calls:      registryAddCalls,
				Over:       newRegistryOver,
				Callable:   RegistryChecks.Add.Callable,
				Concurrent: true,
			},
			{
				Calls:      registryInsertCalls,
				Over:       newRegistryOver,
				Callable:   RegistryChecks.Insert.Callable,
				Concurrent: true,
			},
			{
				Calls:      registryStampCalls,
				Over:       newRegistryOver,
				Callable:   RegistryChecks.Stamp.Callable,
				Concurrent: true,
			},
			{
				Calls:      registryTouchCalls,
				Over:       newRegistryOver,
				Callable:   RegistryChecks.Touch.Callable,
				Concurrent: true,
			},
		},
	}
}

// RunRegistrySuite applies the checks of Registry, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunRegistrySuite[T Registry](t *testing.T, newSubject func() T, opts ...RegistryOption) {
	t.Helper()
	RunRegistrySuiteFromHandle(t, func(*RegistryHandle) T { return newSubject() }, opts...)
}

// RunRegistrySuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRegistrySuiteFromHandle[T Registry](
	t *testing.T,
	newSubject func(h *RegistryHandle) T,
	opts ...RegistryOption,
) {
	t.Helper()
	conform.Run(t, RegistrySuite(), registrySubjects(newSubject), opts...)
}

// FuzzRegistrySuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRegistrySuite[T Registry](
	f *testing.F,
	callable conform.Callable[Registry],
	newSubject func(h *RegistryHandle) T,
	opts ...RegistryOption,
) {
	f.Helper()
	conform.Fuzz(f, RegistrySuite(), callable, registrySubjects(newSubject), opts...)
}

// RunRegistryBudgets applies the allocation ceilings of Registry and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRegistryBudgets[T Registry](
	t *testing.T,
	newSubject func(h *RegistryHandle) T,
	opts ...RegistryOption,
) {
	t.Helper()
	conform.RunBudgets(t, RegistrySuite(), registrySubjects(newSubject), opts...)
}

// BenchRegistrySuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchRegistrySuite[T Registry](
	b conform.B,
	callable conform.Callable[Registry],
	newSubject func(h *RegistryHandle) T,
	opts ...RegistryOption,
) {
	b.Helper()
	conform.Bench(b, RegistrySuite(), callable, registrySubjects(newSubject), opts...)
}

// BenchRegistrySuiteParallel measures callable on a subject of newSubject
// on GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchRegistrySuiteParallel[T Registry](
	b *testing.B,
	callable conform.Callable[Registry],
	newSubject func(h *RegistryHandle) T,
	opts ...RegistryOption,
) {
	b.Helper()
	conform.BenchParallel(b, RegistrySuite(), callable, registrySubjects(newSubject), opts...)
}

// registrySubjects returns the factory of the subjects of Registry that
// calls newSubject with the handle of each check.
func registrySubjects[T Registry](
	newSubject func(h *RegistryHandle) T,
) func(h *RegistryHandle) Registry {
	return func(h *RegistryHandle) Registry { return newSubject(h) }
}

// registrySetCalls generates the calls of Registry.Set with the inputs
// that each case draws.
var registrySetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Registry) error {
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Registry) error {
		return subject.Set(ctx, key, value)
	}
})

// registryValueCalls generates the calls of Registry.Value with the inputs
// that each case draws.
var registryValueCalls = prop.Composite(func(c *prop.Case) func(context.Context, Registry) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Registry) error {
		_, err := subject.Value(ctx, key)
		return err
	}
})

// registryAwaitCalls generates the calls of Registry.Await with the inputs
// that each case draws.
var registryAwaitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Registry) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Registry) error {
		_, err := subject.Await(ctx, key)
		return err
	}
})

// registryAddCalls generates the calls of Registry.Add with the inputs
// that each case draws.
var registryAddCalls = prop.Composite(func(c *prop.Case) func(context.Context, Registry) error {
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Registry) error {
		_, err := subject.Add(ctx, value)
		return err
	}
})

// registryInsertCalls generates the calls of Registry.Insert with the
// inputs that each case draws.
var registryInsertCalls = prop.Composite(func(c *prop.Case) func(context.Context, Registry) error {
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Registry) error {
		_, err := subject.Insert(ctx, value)
		return err
	}
})

// registryStampCalls generates the calls of Registry.Stamp with the inputs
// that each case draws.
var registryStampCalls = prop.Composite(func(*prop.Case) func(context.Context, Registry) error {
	return func(ctx context.Context, subject Registry) error {
		_, err := subject.Stamp(ctx)
		return err
	}
})

// registryTouchCalls generates the calls of Registry.Touch with the inputs
// that each case draws.
var registryTouchCalls = prop.Composite(func(*prop.Case) func(context.Context, Registry) error {
	return func(ctx context.Context, subject Registry) error {
		subject.Touch(ctx)
		return nil
	}
})

// registryOverEmbedded is Registry under the name of the field that
// registryOver embeds.
type registryOverEmbedded = Registry

// registryOver is a subject of Registry whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type registryOver struct {
	registryOverEmbedded
	spend func()
}

// Set calls spend and returns its zero results.
func (o registryOver) Set(context.Context, string, string) (r0 error) {
	o.spend()
	return r0
}

// Value calls spend and returns its zero results.
func (o registryOver) Value(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// Await calls spend and returns its zero results.
func (o registryOver) Await(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// Add calls spend and returns its zero results.
func (o registryOver) Add(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// Insert calls spend and returns its zero results.
func (o registryOver) Insert(context.Context, string) (r0 *Entry, r1 error) {
	o.spend()
	return r0, r1
}

// Stamp calls spend and returns its zero results.
func (o registryOver) Stamp(context.Context) (r0 int64, r1 error) {
	o.spend()
	return r0, r1
}

// Touch calls spend and returns its zero results.
func (o registryOver) Touch(context.Context) {
	o.spend()
}

// newRegistryOver returns the planted defect of a ceiling of Registry,
// whose callables call spend.
func newRegistryOver(_ *RegistryHandle, spend func()) Registry {
	return registryOver{spend: spend}
}

// registrySetPanicsDefect is the planted defect of a Registry whose Set
// panics.
type registrySetPanicsDefect struct {
	registryOverEmbedded
}

// Set is the override of the planted defect.
func (*registrySetPanicsDefect) Set(p0 context.Context, p1, p2 string) (r0 error) {
	panic("the planted defect of a Registry whose Set panics")
}

// registrySetPanics returns the planted defect of a Registry whose Set
// panics over a fresh subject.
func registrySetPanics(h *RegistryHandle) Registry {
	return &registrySetPanicsDefect{registryOverEmbedded: h.Fresh()}
}

// registrySetIgnoresCancelDefect is the planted defect of a Registry whose
// Set ignores the context.
type registrySetIgnoresCancelDefect struct {
	registryOverEmbedded
}

// Set is the override of the planted defect.
func (d *registrySetIgnoresCancelDefect) Set(p0 context.Context, p1, p2 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.registryOverEmbedded.Set(p0, p1, p2)
}

// registrySetIgnoresCancel returns the planted defect of a Registry whose
// Set ignores the context over a fresh subject.
func registrySetIgnoresCancel(h *RegistryHandle) Registry {
	return &registrySetIgnoresCancelDefect{registryOverEmbedded: h.Fresh()}
}

// registrySetSwallowsFailureDefect is the planted defect of a Registry
// whose Set returns no error when its dependency fails.
type registrySetSwallowsFailureDefect struct {
	registryOverEmbedded
	h *RegistryHandle
}

// Set is the override of the planted defect.
func (d *registrySetSwallowsFailureDefect) Set(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	before := d.h.Injector().Failures()
	r0 = d.registryOverEmbedded.Set(p0, p1, p2)
	if r0 != nil && d.h.Injector().Failures() > before {
		r0 = nil
	}
	return r0
}

// registrySetSwallowsFailure returns the planted defect of a Registry
// whose Set returns no error when its dependency fails over a fresh
// subject.
func registrySetSwallowsFailure(h *RegistryHandle) Registry {
	return &registrySetSwallowsFailureDefect{registryOverEmbedded: h.Fresh(), h: h}
}

// registryBlocksAfterRecoveryDefect is the planted defect of a Registry
// whose each call after the first restart blocks until its case ends.
type registryBlocksAfterRecoveryDefect struct {
	registryOverEmbedded
	lifetime context.Context
}

// Set is the override of the planted defect.
func (d *registryBlocksAfterRecoveryDefect) Set(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	if d.lifetime != nil {
		r0 = host.Park(p0)
		return r0
	}
	return d.registryOverEmbedded.Set(p0, p1, p2)
}

// Value is the override of the planted defect.
func (d *registryBlocksAfterRecoveryDefect) Value(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.registryOverEmbedded.Value(p0, p1)
}

// Await is the override of the planted defect.
func (d *registryBlocksAfterRecoveryDefect) Await(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.registryOverEmbedded.Await(p0, p1)
}

// Add is the override of the planted defect.
func (d *registryBlocksAfterRecoveryDefect) Add(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.registryOverEmbedded.Add(p0, p1)
}

// Insert is the override of the planted defect.
func (d *registryBlocksAfterRecoveryDefect) Insert(
	p0 context.Context,
	p1 string,
) (r0 *Entry, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.registryOverEmbedded.Insert(p0, p1)
}

// Stamp is the override of the planted defect.
func (d *registryBlocksAfterRecoveryDefect) Stamp(p0 context.Context) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.registryOverEmbedded.Stamp(p0)
}

// Touch is the override of the planted defect.
func (d *registryBlocksAfterRecoveryDefect) Touch(p0 context.Context) {
	if d.lifetime != nil {
		_ = host.Park(p0)
		return
	}
	d.registryOverEmbedded.Touch(p0)
}

// registryBlocksAfterRecovery returns the subject of the application in
// place of the planted defect of a Registry whose each call after the
// first restart blocks until its case ends. The door Host of its check
// opens the defect in a world.
func registryBlocksAfterRecovery(h *RegistryHandle) Registry {
	return h.Fresh()
}

// registryBlocksAfterRecoveryHost opens the planted defect of a Registry
// whose each call after the first restart blocks until its case ends. The
// defect wraps the subject that the application opens over the environment
// env of an incarnation. From the second incarnation of the host of env
// on, each callable of the defect blocks until its context ends.
func registryBlocksAfterRecoveryHost(h *RegistryHandle, env *host.Env) Registry {
	d := &registryBlocksAfterRecoveryDefect{registryOverEmbedded: h.FreshHost(env)}
	if env.Incarnation() > 1 {
		d.lifetime = env.Lifetime()
	}
	return d
}

// registryValuePanicsDefect is the planted defect of a Registry whose
// Value panics.
type registryValuePanicsDefect struct {
	registryOverEmbedded
}

// Value is the override of the planted defect.
func (*registryValuePanicsDefect) Value(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Registry whose Value panics")
}

// registryValuePanics returns the planted defect of a Registry whose Value
// panics over a fresh subject.
func registryValuePanics(h *RegistryHandle) Registry {
	return &registryValuePanicsDefect{registryOverEmbedded: h.Fresh()}
}

// registryValueIgnoresCancelDefect is the planted defect of a Registry
// whose Value ignores the context.
type registryValueIgnoresCancelDefect struct {
	registryOverEmbedded
}

// Value is the override of the planted defect.
func (d *registryValueIgnoresCancelDefect) Value(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registryOverEmbedded.Value(p0, p1)
}

// registryValueIgnoresCancel returns the planted defect of a Registry
// whose Value ignores the context over a fresh subject.
func registryValueIgnoresCancel(h *RegistryHandle) Registry {
	return &registryValueIgnoresCancelDefect{registryOverEmbedded: h.Fresh()}
}

// registryValueSwallowsFailureDefect is the planted defect of a Registry
// whose Value returns no error when its dependency fails.
type registryValueSwallowsFailureDefect struct {
	registryOverEmbedded
	h *RegistryHandle
}

// Value is the override of the planted defect.
func (d *registryValueSwallowsFailureDefect) Value(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.registryOverEmbedded.Value(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(string)
		r1 = nil
	}
	return r0, r1
}

// registryValueSwallowsFailure returns the planted defect of a Registry
// whose Value returns no error when its dependency fails over a fresh
// subject.
func registryValueSwallowsFailure(h *RegistryHandle) Registry {
	return &registryValueSwallowsFailureDefect{registryOverEmbedded: h.Fresh(), h: h}
}

// registryValueValueBesideDefect is the planted defect of a Registry whose
// Value returns a value beside an error.
type registryValueValueBesideDefect struct {
	registryOverEmbedded
}

// Value is the override of the planted defect.
func (d *registryValueValueBesideDefect) Value(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	r0, r1 = d.registryOverEmbedded.Value(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// registryValueValueBeside returns the planted defect of a Registry whose
// Value returns a value beside an error over a fresh subject.
func registryValueValueBeside(h *RegistryHandle) Registry {
	return &registryValueValueBesideDefect{registryOverEmbedded: h.Fresh()}
}

// registryValueDivergesAfterThreeDefect is the planted defect of a
// Registry whose first callable with a value result returns another value
// from its fourth call on.
type registryValueDivergesAfterThreeDefect struct {
	registryOverEmbedded
	calls int
}

// Value is the override of the planted defect.
func (d *registryValueDivergesAfterThreeDefect) Value(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.registryOverEmbedded.Value(p0, p1)
	d.calls++
	if d.calls > 3 {
		r0 = string("sample-b")
	}
	return r0, r1
}

// registryValueDivergesAfterThree returns the planted defect of a Registry
// whose first callable with a value result returns another value from its
// fourth call on over a fresh subject.
func registryValueDivergesAfterThree(h *RegistryHandle) Registry {
	return &registryValueDivergesAfterThreeDefect{registryOverEmbedded: h.Fresh()}
}

// registryRacyReadWriteDefect is the planted defect of a Registry whose
// first callable with a value result returns another value while another
// call is in progress.
type registryRacyReadWriteDefect struct {
	registryOverEmbedded
	h      *RegistryHandle
	marked atomic.Int32
}

// Set is the override of the planted defect.
func (d *registryRacyReadWriteDefect) Set(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	return d.registryOverEmbedded.Set(p0, p1, p2)
}

// Value is the override of the planted defect.
func (d *registryRacyReadWriteDefect) Value(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.marked.Load() > 0 {
		r0 = string("sample-b")
		return r0, r1
	}
	return d.registryOverEmbedded.Value(p0, p1)
}

// Await is the override of the planted defect.
func (d *registryRacyReadWriteDefect) Await(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	return d.registryOverEmbedded.Await(p0, p1)
}

// Add is the override of the planted defect.
func (d *registryRacyReadWriteDefect) Add(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	return d.registryOverEmbedded.Add(p0, p1)
}

// Insert is the override of the planted defect.
func (d *registryRacyReadWriteDefect) Insert(p0 context.Context, p1 string) (r0 *Entry, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	return d.registryOverEmbedded.Insert(p0, p1)
}

// Stamp is the override of the planted defect.
func (d *registryRacyReadWriteDefect) Stamp(p0 context.Context) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	return d.registryOverEmbedded.Stamp(p0)
}

// Touch is the override of the planted defect.
func (d *registryRacyReadWriteDefect) Touch(p0 context.Context) {
	d.marked.Add(1)
	defer d.marked.Add(-1)
	d.h.Yield()
	d.registryOverEmbedded.Touch(p0)
}

// registryRacyReadWrite returns the planted defect of a Registry whose
// first callable with a value result returns another value while another
// call is in progress over a fresh subject.
func registryRacyReadWrite(h *RegistryHandle) Registry {
	return &registryRacyReadWriteDefect{registryOverEmbedded: h.Fresh(), h: h}
}

// registryAmnesiacRecovery returns the subject of the application in place
// of the planted defect of a Registry whose each restart opens the subject
// over a new medium. The door Host of its check opens the defect in a
// world.
func registryAmnesiacRecovery(h *RegistryHandle) Registry {
	return h.Fresh()
}

// registryAmnesiacRecoveryHost opens the planted defect of a Registry
// whose each restart opens the subject over a new medium. The function
// probe.Amnesiac opens the defect over the environment env of an
// incarnation.
func registryAmnesiacRecoveryHost(h *RegistryHandle, env *host.Env) Registry {
	return probe.Amnesiac(h, env)
}

// registryAwaitIgnoresCancelDefect is the planted defect of a Registry
// whose Await ignores the context.
type registryAwaitIgnoresCancelDefect struct {
	registryOverEmbedded
}

// Await is the override of the planted defect.
func (d *registryAwaitIgnoresCancelDefect) Await(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registryOverEmbedded.Await(p0, p1)
}

// registryAwaitIgnoresCancel returns the planted defect of a Registry
// whose Await ignores the context over a fresh subject.
func registryAwaitIgnoresCancel(h *RegistryHandle) Registry {
	return &registryAwaitIgnoresCancelDefect{registryOverEmbedded: h.Fresh()}
}

// registryAwaitSwallowsFailureDefect is the planted defect of a Registry
// whose Await returns no error when its dependency fails.
type registryAwaitSwallowsFailureDefect struct {
	registryOverEmbedded
	h *RegistryHandle
}

// Await is the override of the planted defect.
func (d *registryAwaitSwallowsFailureDefect) Await(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.registryOverEmbedded.Await(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(string)
		r1 = nil
	}
	return r0, r1
}

// registryAwaitSwallowsFailure returns the planted defect of a Registry
// whose Await returns no error when its dependency fails over a fresh
// subject.
func registryAwaitSwallowsFailure(h *RegistryHandle) Registry {
	return &registryAwaitSwallowsFailureDefect{registryOverEmbedded: h.Fresh(), h: h}
}

// registryAddPanicsDefect is the planted defect of a Registry whose Add
// panics.
type registryAddPanicsDefect struct {
	registryOverEmbedded
}

// Add is the override of the planted defect.
func (*registryAddPanicsDefect) Add(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Registry whose Add panics")
}

// registryAddPanics returns the planted defect of a Registry whose Add
// panics over a fresh subject.
func registryAddPanics(h *RegistryHandle) Registry {
	return &registryAddPanicsDefect{registryOverEmbedded: h.Fresh()}
}

// registryAddIgnoresCancelDefect is the planted defect of a Registry whose
// Add ignores the context.
type registryAddIgnoresCancelDefect struct {
	registryOverEmbedded
}

// Add is the override of the planted defect.
func (d *registryAddIgnoresCancelDefect) Add(p0 context.Context, p1 string) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registryOverEmbedded.Add(p0, p1)
}

// registryAddIgnoresCancel returns the planted defect of a Registry whose
// Add ignores the context over a fresh subject.
func registryAddIgnoresCancel(h *RegistryHandle) Registry {
	return &registryAddIgnoresCancelDefect{registryOverEmbedded: h.Fresh()}
}

// registryAddSwallowsFailureDefect is the planted defect of a Registry
// whose Add returns no error when its dependency fails.
type registryAddSwallowsFailureDefect struct {
	registryOverEmbedded
	h *RegistryHandle
}

// Add is the override of the planted defect.
func (d *registryAddSwallowsFailureDefect) Add(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.registryOverEmbedded.Add(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(string)
		r1 = nil
	}
	return r0, r1
}

// registryAddSwallowsFailure returns the planted defect of a Registry
// whose Add returns no error when its dependency fails over a fresh
// subject.
func registryAddSwallowsFailure(h *RegistryHandle) Registry {
	return &registryAddSwallowsFailureDefect{registryOverEmbedded: h.Fresh(), h: h}
}

// registryAddValueBesideDefect is the planted defect of a Registry whose
// Add returns a value beside an error.
type registryAddValueBesideDefect struct {
	registryOverEmbedded
}

// Add is the override of the planted defect.
func (d *registryAddValueBesideDefect) Add(p0 context.Context, p1 string) (r0 string, r1 error) {
	r0, r1 = d.registryOverEmbedded.Add(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// registryAddValueBeside returns the planted defect of a Registry whose
// Add returns a value beside an error over a fresh subject.
func registryAddValueBeside(h *RegistryHandle) Registry {
	return &registryAddValueBesideDefect{registryOverEmbedded: h.Fresh()}
}

// registryInsertPanicsDefect is the planted defect of a Registry whose
// Insert panics.
type registryInsertPanicsDefect struct {
	registryOverEmbedded
}

// Insert is the override of the planted defect.
func (*registryInsertPanicsDefect) Insert(p0 context.Context, p1 string) (r0 *Entry, r1 error) {
	panic("the planted defect of a Registry whose Insert panics")
}

// registryInsertPanics returns the planted defect of a Registry whose
// Insert panics over a fresh subject.
func registryInsertPanics(h *RegistryHandle) Registry {
	return &registryInsertPanicsDefect{registryOverEmbedded: h.Fresh()}
}

// registryInsertIgnoresCancelDefect is the planted defect of a Registry
// whose Insert ignores the context.
type registryInsertIgnoresCancelDefect struct {
	registryOverEmbedded
}

// Insert is the override of the planted defect.
func (d *registryInsertIgnoresCancelDefect) Insert(
	p0 context.Context,
	p1 string,
) (r0 *Entry, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registryOverEmbedded.Insert(p0, p1)
}

// registryInsertIgnoresCancel returns the planted defect of a Registry
// whose Insert ignores the context over a fresh subject.
func registryInsertIgnoresCancel(h *RegistryHandle) Registry {
	return &registryInsertIgnoresCancelDefect{registryOverEmbedded: h.Fresh()}
}

// registryInsertSwallowsFailureDefect is the planted defect of a Registry
// whose Insert returns no error when its dependency fails.
type registryInsertSwallowsFailureDefect struct {
	registryOverEmbedded
	h *RegistryHandle
}

// Insert is the override of the planted defect.
func (d *registryInsertSwallowsFailureDefect) Insert(
	p0 context.Context,
	p1 string,
) (r0 *Entry, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.registryOverEmbedded.Insert(p0, p1)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(*Entry)
		r1 = nil
	}
	return r0, r1
}

// registryInsertSwallowsFailure returns the planted defect of a Registry
// whose Insert returns no error when its dependency fails over a fresh
// subject.
func registryInsertSwallowsFailure(h *RegistryHandle) Registry {
	return &registryInsertSwallowsFailureDefect{registryOverEmbedded: h.Fresh(), h: h}
}

// registryInsertValueBesideDefect is the planted defect of a Registry
// whose Insert returns a value beside an error.
type registryInsertValueBesideDefect struct {
	registryOverEmbedded
}

// Insert is the override of the planted defect.
func (d *registryInsertValueBesideDefect) Insert(
	p0 context.Context,
	p1 string,
) (r0 *Entry, r1 error) {
	r0, r1 = d.registryOverEmbedded.Insert(p0, p1)
	if r1 != nil {
		r0 = &Entry{Key: "Key-b"}
	}
	return r0, r1
}

// registryInsertValueBeside returns the planted defect of a Registry whose
// Insert returns a value beside an error over a fresh subject.
func registryInsertValueBeside(h *RegistryHandle) Registry {
	return &registryInsertValueBesideDefect{registryOverEmbedded: h.Fresh()}
}

// registryStampPanicsDefect is the planted defect of a Registry whose
// Stamp panics.
type registryStampPanicsDefect struct {
	registryOverEmbedded
}

// Stamp is the override of the planted defect.
func (*registryStampPanicsDefect) Stamp(p0 context.Context) (r0 int64, r1 error) {
	panic("the planted defect of a Registry whose Stamp panics")
}

// registryStampPanics returns the planted defect of a Registry whose Stamp
// panics over a fresh subject.
func registryStampPanics(h *RegistryHandle) Registry {
	return &registryStampPanicsDefect{registryOverEmbedded: h.Fresh()}
}

// registryStampIgnoresCancelDefect is the planted defect of a Registry
// whose Stamp ignores the context.
type registryStampIgnoresCancelDefect struct {
	registryOverEmbedded
}

// Stamp is the override of the planted defect.
func (d *registryStampIgnoresCancelDefect) Stamp(p0 context.Context) (r0 int64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registryOverEmbedded.Stamp(p0)
}

// registryStampIgnoresCancel returns the planted defect of a Registry
// whose Stamp ignores the context over a fresh subject.
func registryStampIgnoresCancel(h *RegistryHandle) Registry {
	return &registryStampIgnoresCancelDefect{registryOverEmbedded: h.Fresh()}
}

// registryStampSwallowsFailureDefect is the planted defect of a Registry
// whose Stamp returns no error when its dependency fails.
type registryStampSwallowsFailureDefect struct {
	registryOverEmbedded
	h *RegistryHandle
}

// Stamp is the override of the planted defect.
func (d *registryStampSwallowsFailureDefect) Stamp(p0 context.Context) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	before := d.h.Injector().Failures()
	r0, r1 = d.registryOverEmbedded.Stamp(p0)
	if r1 != nil && d.h.Injector().Failures() > before {
		r0 = *new(int64)
		r1 = nil
	}
	return r0, r1
}

// registryStampSwallowsFailure returns the planted defect of a Registry
// whose Stamp returns no error when its dependency fails over a fresh
// subject.
func registryStampSwallowsFailure(h *RegistryHandle) Registry {
	return &registryStampSwallowsFailureDefect{registryOverEmbedded: h.Fresh(), h: h}
}

// registryStampValueBesideDefect is the planted defect of a Registry whose
// Stamp returns a value beside an error.
type registryStampValueBesideDefect struct {
	registryOverEmbedded
}

// Stamp is the override of the planted defect.
func (d *registryStampValueBesideDefect) Stamp(p0 context.Context) (r0 int64, r1 error) {
	r0, r1 = d.registryOverEmbedded.Stamp(p0)
	if r1 != nil {
		r0 = int64(7)
	}
	return r0, r1
}

// registryStampValueBeside returns the planted defect of a Registry whose
// Stamp returns a value beside an error over a fresh subject.
func registryStampValueBeside(h *RegistryHandle) Registry {
	return &registryStampValueBesideDefect{registryOverEmbedded: h.Fresh()}
}

// registryTouchPanicsDefect is the planted defect of a Registry whose
// Touch panics.
type registryTouchPanicsDefect struct {
	registryOverEmbedded
}

// Touch is the override of the planted defect.
func (*registryTouchPanicsDefect) Touch(p0 context.Context) {
	panic("the planted defect of a Registry whose Touch panics")
}

// registryTouchPanics returns the planted defect of a Registry whose Touch
// panics over a fresh subject.
func registryTouchPanics(h *RegistryHandle) Registry {
	return &registryTouchPanicsDefect{registryOverEmbedded: h.Fresh()}
}

// registryActions returns the actions of the generated sequences of
// Registry over the subject that subject returns. The actions of a case
// take their keys from the pools that it draws.
func registryActions(c *prop.Case, subject func() Registry) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Set",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "set.key")
				value := c.Draw(prop.Of[string](), "set.value")
				return []any{key, value}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Set", args, args[0])
				values, err := registrySetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Value",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "value.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Value", args, args[0])
				values, err := registryValueCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Add",
			Input: func(c *prop.Case, _ struct{}) any {
				value := c.Draw(prop.SampledFrom(keysString...), "add.value")
				return []any{value}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Add", args, args[0])
				values, err := registryAddCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Insert",
			Input: func(c *prop.Case, _ struct{}) any {
				value := c.Draw(prop.SampledFrom(keysString...), "insert.value")
				return []any{value}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Insert", args, args[0])
				values, err := registryInsertCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Stamp",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Stamp", args)
				values, err := registryStampCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Touch",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Touch", args)
				values, err := registryTouchCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// registrySetCall calls Set on subject with the arguments args, and
// returns the values of the call beside its error.
func registrySetCall(ctx context.Context, subject Registry, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	return nil, subject.Set(ctx, a0, a1)
}

// registryValueCall calls Value on subject with the arguments args, and
// returns the values of the call beside its error.
func registryValueCall(ctx context.Context, subject Registry, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Value(ctx, a0)
	return []any{r0}, err
}

// registryAwaitCall calls Await on subject with the arguments args, and
// returns the values of the call beside its error.
func registryAwaitCall(ctx context.Context, subject Registry, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Await(ctx, a0)
	return []any{r0}, err
}

// registryAddCall calls Add on subject with the arguments args, and
// returns the values of the call beside its error.
func registryAddCall(ctx context.Context, subject Registry, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Add(ctx, a0)
	return []any{r0}, err
}

// registryInsertCall calls Insert on subject with the arguments args, and
// returns the values of the call beside its error.
func registryInsertCall(ctx context.Context, subject Registry, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Insert(ctx, a0)
	return []any{r0}, err
}

// registryStampCall calls Stamp on subject with the arguments args, and
// returns the values of the call beside its error.
func registryStampCall(ctx context.Context, subject Registry, args []any) ([]any, error) {
	r0, err := subject.Stamp(ctx)
	return []any{r0}, err
}

// registryTouchCall calls Touch on subject with the arguments args, and
// returns the values of the call beside its error.
func registryTouchCall(ctx context.Context, subject Registry, args []any) ([]any, error) {
	subject.Touch(ctx)
	return nil, nil
}

// registryReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func registryReference(r Registry) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Set":
			values, err = registrySetCall(ctx, r, args)
		case "Value":
			values, err = registryValueCall(ctx, r, args)
		case "Await":
			values, err = registryAwaitCall(ctx, r, args)
		case "Add":
			values, err = registryAddCall(ctx, r, args)
		case "Insert":
			values, err = registryInsertCall(ctx, r, args)
		case "Stamp":
			values, err = registryStampCall(ctx, r, args)
		case "Touch":
			values, err = registryTouchCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// registryWorld is the world of the subject of Registry, which the checks
// of its world probes run.
var registryWorld = probe.World[Registry]{
	Operations: registryOperations,
	Reference:  registryReference,
	Corrupt:    ErrCorrupt,
	Clients:    3,
}

// registryOperations returns the operations of the world of Registry in
// the case c, over the pools of the keys that c chooses.
func registryOperations(c *prop.Case) []probe.Operation[Registry] {
	keysString := probe.NewPool(
		c,
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []probe.Operation[Registry]{
		{
			Name: "Set",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "set.key")
				value := c.Draw(prop.Of[string](), "set.value")
				return []any{key, value}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Registry, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := registrySetCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
		{
			Name: "Value",
			Input: func(c *prop.Case) any {
				key := keysString.Pick(c, "value.key")
				return []any{key}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Registry, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := registryValueCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
		{
			Name: "Add",
			Input: func(c *prop.Case) any {
				value := keysString.Pick(c, "add.value")
				return []any{value}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Registry, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := registryAddCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
			Produced: func(out conform.Output) {
				if v, ok := out.Values[0].(string); ok {
					keysString.Add(v)
				}
			},
		},
		{
			Name: "Insert",
			Input: func(c *prop.Case) any {
				value := keysString.Pick(c, "insert.value")
				return []any{value}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				args, _ := input.([]any)
				return []any{args[0]}
			},
			Call: func(ctx context.Context, s Registry, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := registryInsertCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
			Produced: func(out conform.Output) {
				if v, ok := out.Values[0].(*Entry); ok && v != nil {
					keysString.Add(v.Key)
				}
			},
		},
		{
			Name: "Stamp",
			Input: func(c *prop.Case) any {
				return []any{}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				return nil
			},
			Call: func(ctx context.Context, s Registry, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := registryStampCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
		{
			Name: "Touch",
			Input: func(c *prop.Case) any {
				return []any{}
			},
			Args: func(input any) []any {
				args, _ := input.([]any)
				return args
			},
			Keys: func(input any) []any {
				return nil
			},
			Call: func(ctx context.Context, s Registry, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := registryTouchCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
	}
}
