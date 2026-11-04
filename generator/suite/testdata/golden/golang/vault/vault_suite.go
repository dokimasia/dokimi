package vault

import (
	"context"
	"errors"
	"io/fs"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
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
// over a medium whose syncs persist nothing over a new medium in memory.
func vaultPutAcknowledgesBeforeSync(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultPutAcknowledgesBeforeSyncOpen opens a subject of the application
// over the medium m that probe.Unsynced wraps, as the planted defect of a
// Vault over a medium whose syncs persist nothing.
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
// Put writes the zero value under the key before its first value over a
// new medium in memory.
func vaultPutTrustsAPartialRecord(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultPutTrustsAPartialRecordOpen opens the planted defect of a Vault
// whose Put writes the zero value under the key before its first value
// over a subject that the application opens over the medium m.
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
// returns the zero value of a corrupted record over a new medium in
// memory.
func vaultGetIgnoresChecksum(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultGetIgnoresChecksumOpen opens the planted defect of a Vault whose
// Get returns the zero value of a corrupted record over a subject that the
// application opens over the medium m.
func vaultGetIgnoresChecksumOpen(h *VaultHandle, m medium.FS) Vault {
	return &vaultGetIgnoresChecksumDefect{vaultOverEmbedded: h.FreshOpen(m)}
}

// vaultPutHidesAFullMedium returns the planted defect of a Vault over a
// medium that reports a full write as written over a new medium in memory.
func vaultPutHidesAFullMedium(h *VaultHandle) Vault {
	return h.Open(medium.NewMemory())
}

// vaultPutHidesAFullMediumOpen opens a subject of the application over the
// medium m that probe.Unfilled wraps, as the planted defect of a Vault
// over a medium that reports a full write as written.
func vaultPutHidesAFullMediumOpen(h *VaultHandle, m medium.FS) Vault {
	return h.FreshOpen(probe.Unfilled(m))
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
		Calls: []conform.Call[Journal]{
			{
				Calls:    journalAppendCalls,
				Over:     newJournalOver,
				Callable: JournalChecks.Append.Callable,
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
// Journal over a medium whose syncs persist nothing over a new medium in
// memory.
func journalAppendAcknowledgesBeforeSync(h *JournalHandle) Journal {
	return h.Open(medium.NewMemory())
}

// journalAppendAcknowledgesBeforeSyncOpen opens a subject of the
// application over the medium m that probe.Unsynced wraps, as the planted
// defect of a Journal over a medium whose syncs persist nothing.
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
// first value over a new medium in memory.
func journalAppendTrustsAPartialRecord(h *JournalHandle) Journal {
	return h.Open(medium.NewMemory())
}

// journalAppendTrustsAPartialRecordOpen opens the planted defect of a
// Journal whose Append writes the zero value under the key before its
// first value over a subject that the application opens over the medium m.
func journalAppendTrustsAPartialRecordOpen(h *JournalHandle, m medium.FS) Journal {
	return &journalAppendTrustsAPartialRecordDefect{journalOverEmbedded: h.FreshOpen(m)}
}

// journalAppendHidesAFullMedium returns the planted defect of a Journal
// over a medium that reports a full write as written over a new medium in
// memory.
func journalAppendHidesAFullMedium(h *JournalHandle) Journal {
	return h.Open(medium.NewMemory())
}

// journalAppendHidesAFullMediumOpen opens a subject of the application
// over the medium m that probe.Unfilled wraps, as the planted defect of a
// Journal over a medium that reports a full write as written.
func journalAppendHidesAFullMediumOpen(h *JournalHandle, m medium.FS) Journal {
	return h.FreshOpen(probe.Unfilled(m))
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
// Ledger over a medium whose syncs persist nothing over a new medium in
// memory.
func ledgerRecordAcknowledgesBeforeSync(h *LedgerHandle) Ledger {
	return h.Open(medium.NewMemory())
}

// ledgerRecordAcknowledgesBeforeSyncOpen opens a subject of the
// application over the medium m that probe.Unsynced wraps, as the planted
// defect of a Ledger over a medium whose syncs persist nothing.
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
// whose Record writes the zero value under the key before its first value
// over a new medium in memory.
func ledgerRecordTrustsAPartialRecord(h *LedgerHandle) Ledger {
	return h.Open(medium.NewMemory())
}

// ledgerRecordTrustsAPartialRecordOpen opens the planted defect of a
// Ledger whose Record writes the zero value under the key before its first
// value over a subject that the application opens over the medium m.
func ledgerRecordTrustsAPartialRecordOpen(h *LedgerHandle, m medium.FS) Ledger {
	return &ledgerRecordTrustsAPartialRecordDefect{ledgerOverEmbedded: h.FreshOpen(m)}
}

// ledgerRecordHidesAFullMedium returns the planted defect of a Ledger over
// a medium that reports a full write as written over a new medium in
// memory.
func ledgerRecordHidesAFullMedium(h *LedgerHandle) Ledger {
	return h.Open(medium.NewMemory())
}

// ledgerRecordHidesAFullMediumOpen opens a subject of the application over
// the medium m that probe.Unfilled wraps, as the planted defect of a
// Ledger over a medium that reports a full write as written.
func ledgerRecordHidesAFullMediumOpen(h *LedgerHandle, m medium.FS) Ledger {
	return h.FreshOpen(probe.Unfilled(m))
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
