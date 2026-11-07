package sealed

import (
	"context"
	"testing"
	"time"

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

// The names that a test of AEAD writes.
type (
	AEADHandle = conform.Handle[AEAD]
	AEADCheck  = conform.Check[AEAD]
	AEADOption = conform.Option[AEAD]
)

// AEADChecks contains the name of each callable of AEAD, with the
// identities of its checks.
var AEADChecks = struct {
	Seal      aeadSeal
	Open      aeadOpen
	NonceSize aeadNonceSize
	Overhead  aeadOverhead
}{
	Seal:      aeadSeal{"AEAD.Seal"},
	Open:      aeadOpen{"AEAD.Open"},
	NonceSize: aeadNonceSize{"AEAD.NonceSize"},
	Overhead:  aeadOverhead{"AEAD.Overhead"},
}

// aeadSeal is the name of AEAD.Seal, with the identities of its checks.
type aeadSeal struct {
	conform.Callable[AEAD]
}

// Smoke returns the identity of the check that Seal returns without a
// panic for derived inputs.
func (c aeadSeal) Smoke() conform.ID[AEAD] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Seal returns the
// same value for the same inputs.
func (c aeadSeal) Deterministic() conform.ID[AEAD] {
	return c.Check("deterministic")
}

// aeadOpen is the name of AEAD.Open, with the identities of its checks.
type aeadOpen struct {
	conform.Callable[AEAD]
}

// Smoke returns the identity of the check that Open returns without a
// panic for derived inputs.
func (c aeadOpen) Smoke() conform.ID[AEAD] {
	return c.Check("smoke")
}

// aeadNonceSize is the name of AEAD.NonceSize, with the identities of its
// checks.
type aeadNonceSize struct {
	conform.Callable[AEAD]
}

// Smoke returns the identity of the check that NonceSize returns without a
// panic for derived inputs.
func (c aeadNonceSize) Smoke() conform.ID[AEAD] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that NonceSize returns
// the same value for the same inputs.
func (c aeadNonceSize) Deterministic() conform.ID[AEAD] {
	return c.Check("deterministic")
}

// aeadOverhead is the name of AEAD.Overhead, with the identities of its
// checks.
type aeadOverhead struct {
	conform.Callable[AEAD]
}

// Smoke returns the identity of the check that Overhead returns without a
// panic for derived inputs.
func (c aeadOverhead) Smoke() conform.ID[AEAD] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Overhead returns
// the same value for the same inputs.
func (c aeadOverhead) Deterministic() conform.ID[AEAD] {
	return c.Check("deterministic")
}

// AEADSuite returns the checks of AEAD, and the calls and the budgets of
// its callables.
func AEADSuite() conform.Suite[AEAD] {
	return conform.Suite[AEAD]{
		Declaration: "AEAD",
		Package:     "example.com/shop/sealed",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[AEAD]{
			{
				ID:        AEADChecks.Seal.Smoke(),
				Statement: "Seal returns without a panic for derived inputs",
				Claim:     spec.Smoke,
				Machineries: []spec.MachineryID{
					spec.InputsBuiltFromTheSubject,
					spec.NothingBeyondACaller,
				},
				Shallow: true,
				Apply: func(h *AEADHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtSeal0, builtSeal1, builtSeal2, builtSeal3 := SealInputs(h.Subject())
						h.Subject().Seal(builtSeal0, builtSeal1, builtSeal2, builtSeal3)
						return nil
					})
				},
				Defect: aeadSealPanics,
			},
			{
				ID:        AEADChecks.Seal.Deterministic(),
				Statement: "Seal returns the same value for the same inputs",
				Claim:     spec.Deterministic,
				Machineries: []spec.MachineryID{
					spec.InputsBuiltFromTheSubject,
					spec.NothingBeyondACaller,
				},
				Shallow: true,
				Entails: []conform.ID[AEAD]{AEADChecks.Seal.Smoke()},
				Apply: func(h *AEADHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) ([]byte, error) {
						builtSeal0, builtSeal1, builtSeal2, builtSeal3 := SealInputs(h.Subject())
						r0 := h.Subject().Seal(builtSeal0, builtSeal1, builtSeal2, builtSeal3)
						return r0, nil
					})
				},
				Defect: aeadSealDrifts,
			},
			{
				ID:        AEADChecks.Open.Smoke(),
				Statement: "Open returns without a panic for derived inputs",
				Claim:     spec.Smoke,
				Machineries: []spec.MachineryID{
					spec.InputsBuiltFromTheSubject,
					spec.NothingBeyondACaller,
				},
				Shallow: true,
				Apply: func(h *AEADHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtOpen0, builtOpen1, builtOpen2, builtOpen3 := OpenInputs(h.Subject())
						_, err := h.Subject().Open(builtOpen0, builtOpen1, builtOpen2, builtOpen3)
						return err
					})
				},
				Defect: aeadOpenPanics,
			},
			{
				ID:          AEADChecks.NonceSize.Smoke(),
				Statement:   "NonceSize returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AEADHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().NonceSize()
						return nil
					})
				},
				Defect: aeadNonceSizePanics,
			},
			{
				ID:          AEADChecks.NonceSize.Deterministic(),
				Statement:   "NonceSize returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[AEAD]{AEADChecks.NonceSize.Smoke()},
				Apply: func(h *AEADHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().NonceSize()
						return r0, nil
					})
				},
				Defect: aeadNonceSizeDrifts,
			},
			{
				ID:          AEADChecks.Overhead.Smoke(),
				Statement:   "Overhead returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *AEADHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Overhead()
						return nil
					})
				},
				Defect: aeadOverheadPanics,
			},
			{
				ID:          AEADChecks.Overhead.Deterministic(),
				Statement:   "Overhead returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[AEAD]{AEADChecks.Overhead.Smoke()},
				Apply: func(h *AEADHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Overhead()
						return r0, nil
					})
				},
				Defect: aeadOverheadDrifts,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "mixin:sample",
				Callable:       "AEAD.Seal",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
			{
				Classification: "mixin:sample",
				Callable:       "AEAD.Open",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
		},
		Calls: []conform.Call[AEAD]{
			{
				Calls:    aeadSealCalls,
				Over:     newAEADOver,
				Callable: AEADChecks.Seal.Callable,
			},
			{
				Calls:    aeadOpenCalls,
				Over:     newAEADOver,
				Callable: AEADChecks.Open.Callable,
			},
			{
				Calls:    aeadNonceSizeCalls,
				Over:     newAEADOver,
				Callable: AEADChecks.NonceSize.Callable,
			},
			{
				Calls:    aeadOverheadCalls,
				Over:     newAEADOver,
				Callable: AEADChecks.Overhead.Callable,
			},
		},
	}
}

// RunAEADSuite applies the checks of AEAD, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunAEADSuite[T AEAD](t *testing.T, newSubject func() T, opts ...AEADOption) {
	t.Helper()
	RunAEADSuiteFromHandle(t, func(*AEADHandle) T { return newSubject() }, opts...)
}

// RunAEADSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunAEADSuiteFromHandle[T AEAD](
	t *testing.T,
	newSubject func(h *AEADHandle) T,
	opts ...AEADOption,
) {
	t.Helper()
	conform.Run(t, AEADSuite(), aeadSubjects(newSubject), opts...)
}

// FuzzAEADSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzAEADSuite[T AEAD](
	f *testing.F,
	callable conform.Callable[AEAD],
	newSubject func(h *AEADHandle) T,
	opts ...AEADOption,
) {
	f.Helper()
	conform.Fuzz(f, AEADSuite(), callable, aeadSubjects(newSubject), opts...)
}

// RunAEADBudgets applies the allocation ceilings of AEAD and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunAEADBudgets[T AEAD](t *testing.T, newSubject func(h *AEADHandle) T, opts ...AEADOption) {
	t.Helper()
	conform.RunBudgets(t, AEADSuite(), aeadSubjects(newSubject), opts...)
}

// BenchAEADSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchAEADSuite[T AEAD](
	b conform.B,
	callable conform.Callable[AEAD],
	newSubject func(h *AEADHandle) T,
	opts ...AEADOption,
) {
	b.Helper()
	conform.Bench(b, AEADSuite(), callable, aeadSubjects(newSubject), opts...)
}

// aeadSubjects returns the factory of the subjects of AEAD that calls
// newSubject with the handle of each check.
func aeadSubjects[T AEAD](newSubject func(h *AEADHandle) T) func(h *AEADHandle) AEAD {
	return func(h *AEADHandle) AEAD { return newSubject(h) }
}

// aeadSealCalls generates the calls of AEAD.Seal with the inputs that its
// builder SealInputs returns.
var aeadSealCalls = prop.Composite(func(*prop.Case) func(context.Context, AEAD) error {
	return func(_ context.Context, subject AEAD) error {
		dst, nonce, plaintext, additionalData := SealInputs(subject)
		_ = subject.Seal(dst, nonce, plaintext, additionalData)
		return nil
	}
})

// aeadOpenCalls generates the calls of AEAD.Open with the inputs that its
// builder OpenInputs returns.
var aeadOpenCalls = prop.Composite(func(*prop.Case) func(context.Context, AEAD) error {
	return func(_ context.Context, subject AEAD) error {
		dst, nonce, ciphertext, additionalData := OpenInputs(subject)
		_, err := subject.Open(dst, nonce, ciphertext, additionalData)
		return err
	}
})

// aeadNonceSizeCalls generates the calls of AEAD.NonceSize with the inputs
// that each case draws.
var aeadNonceSizeCalls = prop.Composite(func(*prop.Case) func(context.Context, AEAD) error {
	return func(_ context.Context, subject AEAD) error {
		_ = subject.NonceSize()
		return nil
	}
})

// aeadOverheadCalls generates the calls of AEAD.Overhead with the inputs
// that each case draws.
var aeadOverheadCalls = prop.Composite(func(*prop.Case) func(context.Context, AEAD) error {
	return func(_ context.Context, subject AEAD) error {
		_ = subject.Overhead()
		return nil
	}
})

// aeadOverEmbedded is AEAD under the name of the field that aeadOver
// embeds.
type aeadOverEmbedded = AEAD

// aeadOver is a subject of AEAD whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type aeadOver struct {
	aeadOverEmbedded
	spend func()
}

// Seal calls spend and returns its zero results.
func (o aeadOver) Seal([]byte, []byte, []byte, []byte) (r0 []byte) {
	o.spend()
	return r0
}

// Open calls spend and returns its zero results.
func (o aeadOver) Open([]byte, []byte, []byte, []byte) (r0 []byte, r1 error) {
	o.spend()
	return r0, r1
}

// NonceSize calls spend and returns its zero results.
func (o aeadOver) NonceSize() (r0 int) {
	o.spend()
	return r0
}

// Overhead calls spend and returns its zero results.
func (o aeadOver) Overhead() (r0 int) {
	o.spend()
	return r0
}

// newAEADOver returns the planted defect of a ceiling of AEAD, whose
// callables call spend.
func newAEADOver(_ *AEADHandle, spend func()) AEAD {
	return aeadOver{spend: spend}
}

// aeadSealPanicsDefect is the planted defect of a AEAD whose Seal panics.
type aeadSealPanicsDefect struct {
	aeadOverEmbedded
}

// Seal is the override of the planted defect.
func (*aeadSealPanicsDefect) Seal(p0, p1, p2, p3 []byte) (r0 []byte) {
	panic("the planted defect of a AEAD whose Seal panics")
}

// aeadSealPanics returns the planted defect of a AEAD whose Seal panics
// over a fresh subject.
func aeadSealPanics(h *AEADHandle) AEAD {
	return &aeadSealPanicsDefect{aeadOverEmbedded: h.Fresh()}
}

// aeadSealDriftsDefect is the planted defect of a AEAD whose Seal returns
// another value on the repeat.
type aeadSealDriftsDefect struct {
	aeadOverEmbedded
	calls int
}

// Seal is the override of the planted defect.
func (d *aeadSealDriftsDefect) Seal(p0, p1, p2, p3 []byte) (r0 []byte) {
	d.calls++
	if d.calls > 1 {
		r0 = []byte{2}
	}
	return r0
}

// aeadSealDrifts returns the planted defect of a AEAD whose Seal returns
// another value on the repeat over a fresh subject.
func aeadSealDrifts(h *AEADHandle) AEAD {
	return &aeadSealDriftsDefect{aeadOverEmbedded: h.Fresh()}
}

// aeadOpenPanicsDefect is the planted defect of a AEAD whose Open panics.
type aeadOpenPanicsDefect struct {
	aeadOverEmbedded
}

// Open is the override of the planted defect.
func (*aeadOpenPanicsDefect) Open(p0, p1, p2, p3 []byte) (r0 []byte, r1 error) {
	panic("the planted defect of a AEAD whose Open panics")
}

// aeadOpenPanics returns the planted defect of a AEAD whose Open panics
// over a fresh subject.
func aeadOpenPanics(h *AEADHandle) AEAD {
	return &aeadOpenPanicsDefect{aeadOverEmbedded: h.Fresh()}
}

// aeadNonceSizePanicsDefect is the planted defect of a AEAD whose
// NonceSize panics.
type aeadNonceSizePanicsDefect struct {
	aeadOverEmbedded
}

// NonceSize is the override of the planted defect.
func (*aeadNonceSizePanicsDefect) NonceSize() (r0 int) {
	panic("the planted defect of a AEAD whose NonceSize panics")
}

// aeadNonceSizePanics returns the planted defect of a AEAD whose NonceSize
// panics over a fresh subject.
func aeadNonceSizePanics(h *AEADHandle) AEAD {
	return &aeadNonceSizePanicsDefect{aeadOverEmbedded: h.Fresh()}
}

// aeadNonceSizeDriftsDefect is the planted defect of a AEAD whose
// NonceSize returns another value on the repeat.
type aeadNonceSizeDriftsDefect struct {
	aeadOverEmbedded
	calls int
}

// NonceSize is the override of the planted defect.
func (d *aeadNonceSizeDriftsDefect) NonceSize() (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// aeadNonceSizeDrifts returns the planted defect of a AEAD whose NonceSize
// returns another value on the repeat over a fresh subject.
func aeadNonceSizeDrifts(h *AEADHandle) AEAD {
	return &aeadNonceSizeDriftsDefect{aeadOverEmbedded: h.Fresh()}
}

// aeadOverheadPanicsDefect is the planted defect of a AEAD whose Overhead
// panics.
type aeadOverheadPanicsDefect struct {
	aeadOverEmbedded
}

// Overhead is the override of the planted defect.
func (*aeadOverheadPanicsDefect) Overhead() (r0 int) {
	panic("the planted defect of a AEAD whose Overhead panics")
}

// aeadOverheadPanics returns the planted defect of a AEAD whose Overhead
// panics over a fresh subject.
func aeadOverheadPanics(h *AEADHandle) AEAD {
	return &aeadOverheadPanicsDefect{aeadOverEmbedded: h.Fresh()}
}

// aeadOverheadDriftsDefect is the planted defect of a AEAD whose Overhead
// returns another value on the repeat.
type aeadOverheadDriftsDefect struct {
	aeadOverEmbedded
	calls int
}

// Overhead is the override of the planted defect.
func (d *aeadOverheadDriftsDefect) Overhead() (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// aeadOverheadDrifts returns the planted defect of a AEAD whose Overhead
// returns another value on the repeat over a fresh subject.
func aeadOverheadDrifts(h *AEADHandle) AEAD {
	return &aeadOverheadDriftsDefect{aeadOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Metronome writes.
type (
	MetronomeHandle = conform.Handle[Metronome]
	MetronomeCheck  = conform.Check[Metronome]
	MetronomeOption = conform.Option[Metronome]
)

// MetronomeChecks contains the name of each callable of Metronome, with
// the identities of its checks.
var MetronomeChecks = struct {
	Tick metronomeTick
	Pace metronomePace
}{
	Tick: metronomeTick{"Metronome.Tick"},
	Pace: metronomePace{"Metronome.Pace"},
}

// metronomeTick is the name of Metronome.Tick, with the identities of its
// checks.
type metronomeTick struct {
	conform.Callable[Metronome]
}

// Smoke returns the identity of the check that Tick returns without a
// panic for derived inputs.
func (c metronomeTick) Smoke() conform.ID[Metronome] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Tick returns an error for
// a cancelled context.
func (c metronomeTick) Cancel() conform.ID[Metronome] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Tick returns the zero
// value beside an error.
func (c metronomeTick) ZeroOnError() conform.ID[Metronome] {
	return c.Check("zero-on-error")
}

// metronomePace is the name of Metronome.Pace, with the identities of its
// checks.
type metronomePace struct {
	conform.Callable[Metronome]
}

// Smoke returns the identity of the check that Pace returns without a
// panic for derived inputs.
func (c metronomePace) Smoke() conform.ID[Metronome] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Pace returns an error for
// a cancelled context.
func (c metronomePace) Cancel() conform.ID[Metronome] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Pace returns the zero
// value beside an error.
func (c metronomePace) ZeroOnError() conform.ID[Metronome] {
	return c.Check("zero-on-error")
}

// MetronomeSuite returns the checks of Metronome, and the calls and the
// budgets of its callables.
func MetronomeSuite() conform.Suite[Metronome] {
	return conform.Suite[Metronome]{
		Declaration: "Metronome",
		Package:     "example.com/shop/sealed",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Metronome]{
			{
				ID:          MetronomeChecks.Tick.Smoke(),
				Statement:   "Tick returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MetronomeHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtTick0 := Ticks()
						_, err := h.Subject().Tick(ctx, builtTick0)
						return err
					})
				},
				Defect: metronomeTickPanics,
			},
			{
				ID:          MetronomeChecks.Tick.Cancel(),
				Statement:   "Tick returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MetronomeHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						builtTick0 := Ticks()
						_, err := h.Subject().Tick(ctx, builtTick0)
						return err
					})
				},
				Defect: metronomeTickIgnoresCancel,
			},
			{
				ID:          MetronomeChecks.Tick.ZeroOnError(),
				Statement:   "Tick returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MetronomeHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						builtTick0 := Ticks()
						return h.Subject().Tick(ctx, builtTick0)
					}, int(0))
				},
				Defect: metronomeTickValueBeside,
			},
			{
				ID:          MetronomeChecks.Pace.Smoke(),
				Statement:   "Pace returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MetronomeHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtPace0 := Ticks()
						_, err := h.Subject().Pace(ctx, builtPace0)
						return err
					})
				},
				Defect: metronomePacePanics,
			},
			{
				ID:          MetronomeChecks.Pace.Cancel(),
				Statement:   "Pace returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MetronomeHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						builtPace0 := Ticks()
						_, err := h.Subject().Pace(ctx, builtPace0)
						return err
					})
				},
				Defect: metronomePaceIgnoresCancel,
			},
			{
				ID:          MetronomeChecks.Pace.ZeroOnError(),
				Statement:   "Pace returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MetronomeHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						builtPace0 := Ticks()
						return h.Subject().Pace(ctx, builtPace0)
					}, int(0))
				},
				Defect: metronomePaceValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "mixin:sample",
				Callable:       "Metronome.Tick",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
			{
				Classification: "mixin:sample",
				Callable:       "Metronome.Pace",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
		},
		Calls: []conform.Call[Metronome]{
			{
				Calls:    metronomeTickCalls,
				Over:     newMetronomeOver,
				Callable: MetronomeChecks.Tick.Callable,
			},
			{
				Calls:    metronomePaceCalls,
				Over:     newMetronomeOver,
				Callable: MetronomeChecks.Pace.Callable,
			},
		},
	}
}

// RunMetronomeSuite applies the checks of Metronome, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunMetronomeSuite[T Metronome](t *testing.T, newSubject func() T, opts ...MetronomeOption) {
	t.Helper()
	RunMetronomeSuiteFromHandle(t, func(*MetronomeHandle) T { return newSubject() }, opts...)
}

// RunMetronomeSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunMetronomeSuiteFromHandle[T Metronome](
	t *testing.T,
	newSubject func(h *MetronomeHandle) T,
	opts ...MetronomeOption,
) {
	t.Helper()
	conform.Run(t, MetronomeSuite(), metronomeSubjects(newSubject), opts...)
}

// FuzzMetronomeSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzMetronomeSuite[T Metronome](
	f *testing.F,
	callable conform.Callable[Metronome],
	newSubject func(h *MetronomeHandle) T,
	opts ...MetronomeOption,
) {
	f.Helper()
	conform.Fuzz(f, MetronomeSuite(), callable, metronomeSubjects(newSubject), opts...)
}

// RunMetronomeBudgets applies the allocation ceilings of Metronome and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunMetronomeBudgets[T Metronome](
	t *testing.T,
	newSubject func(h *MetronomeHandle) T,
	opts ...MetronomeOption,
) {
	t.Helper()
	conform.RunBudgets(t, MetronomeSuite(), metronomeSubjects(newSubject), opts...)
}

// BenchMetronomeSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchMetronomeSuite[T Metronome](
	b conform.B,
	callable conform.Callable[Metronome],
	newSubject func(h *MetronomeHandle) T,
	opts ...MetronomeOption,
) {
	b.Helper()
	conform.Bench(b, MetronomeSuite(), callable, metronomeSubjects(newSubject), opts...)
}

// metronomeSubjects returns the factory of the subjects of Metronome that
// calls newSubject with the handle of each check.
func metronomeSubjects[T Metronome](
	newSubject func(h *MetronomeHandle) T,
) func(h *MetronomeHandle) Metronome {
	return func(h *MetronomeHandle) Metronome { return newSubject(h) }
}

// metronomeTickCalls generates the calls of Metronome.Tick with the inputs
// that its builder Ticks returns.
var metronomeTickCalls = prop.Composite(func(*prop.Case) func(context.Context, Metronome) error {
	return func(ctx context.Context, subject Metronome) error {
		n := Ticks()
		_, err := subject.Tick(ctx, n)
		return err
	}
})

// metronomePaceCalls generates the calls of Metronome.Pace with the inputs
// that its builder Ticks returns.
var metronomePaceCalls = prop.Composite(func(*prop.Case) func(context.Context, Metronome) error {
	return func(ctx context.Context, subject Metronome) error {
		a1 := Ticks()
		_, err := subject.Pace(ctx, a1)
		return err
	}
})

// metronomeOverEmbedded is Metronome under the name of the field that
// metronomeOver embeds.
type metronomeOverEmbedded = Metronome

// metronomeOver is a subject of Metronome whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type metronomeOver struct {
	metronomeOverEmbedded
	spend func()
}

// Tick calls spend and returns its zero results.
func (o metronomeOver) Tick(context.Context, int) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// Pace calls spend and returns its zero results.
func (o metronomeOver) Pace(context.Context, int) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newMetronomeOver returns the planted defect of a ceiling of Metronome,
// whose callables call spend.
func newMetronomeOver(_ *MetronomeHandle, spend func()) Metronome {
	return metronomeOver{spend: spend}
}

// metronomeTickPanicsDefect is the planted defect of a Metronome whose
// Tick panics.
type metronomeTickPanicsDefect struct {
	metronomeOverEmbedded
}

// Tick is the override of the planted defect.
func (*metronomeTickPanicsDefect) Tick(p0 context.Context, p1 int) (r0 int, r1 error) {
	panic("the planted defect of a Metronome whose Tick panics")
}

// metronomeTickPanics returns the planted defect of a Metronome whose Tick
// panics over a fresh subject.
func metronomeTickPanics(h *MetronomeHandle) Metronome {
	return &metronomeTickPanicsDefect{metronomeOverEmbedded: h.Fresh()}
}

// metronomeTickIgnoresCancelDefect is the planted defect of a Metronome
// whose Tick ignores the context.
type metronomeTickIgnoresCancelDefect struct {
	metronomeOverEmbedded
}

// Tick is the override of the planted defect.
func (d *metronomeTickIgnoresCancelDefect) Tick(p0 context.Context, p1 int) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.metronomeOverEmbedded.Tick(p0, p1)
}

// metronomeTickIgnoresCancel returns the planted defect of a Metronome
// whose Tick ignores the context over a fresh subject.
func metronomeTickIgnoresCancel(h *MetronomeHandle) Metronome {
	return &metronomeTickIgnoresCancelDefect{metronomeOverEmbedded: h.Fresh()}
}

// metronomeTickValueBesideDefect is the planted defect of a Metronome
// whose Tick returns a value beside an error.
type metronomeTickValueBesideDefect struct {
	metronomeOverEmbedded
}

// Tick is the override of the planted defect.
func (d *metronomeTickValueBesideDefect) Tick(p0 context.Context, p1 int) (r0 int, r1 error) {
	r0, r1 = d.metronomeOverEmbedded.Tick(p0, p1)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// metronomeTickValueBeside returns the planted defect of a Metronome whose
// Tick returns a value beside an error over a fresh subject.
func metronomeTickValueBeside(h *MetronomeHandle) Metronome {
	return &metronomeTickValueBesideDefect{metronomeOverEmbedded: h.Fresh()}
}

// metronomePacePanicsDefect is the planted defect of a Metronome whose
// Pace panics.
type metronomePacePanicsDefect struct {
	metronomeOverEmbedded
}

// Pace is the override of the planted defect.
func (*metronomePacePanicsDefect) Pace(p0 context.Context, p1 int) (r0 int, r1 error) {
	panic("the planted defect of a Metronome whose Pace panics")
}

// metronomePacePanics returns the planted defect of a Metronome whose Pace
// panics over a fresh subject.
func metronomePacePanics(h *MetronomeHandle) Metronome {
	return &metronomePacePanicsDefect{metronomeOverEmbedded: h.Fresh()}
}

// metronomePaceIgnoresCancelDefect is the planted defect of a Metronome
// whose Pace ignores the context.
type metronomePaceIgnoresCancelDefect struct {
	metronomeOverEmbedded
}

// Pace is the override of the planted defect.
func (d *metronomePaceIgnoresCancelDefect) Pace(p0 context.Context, p1 int) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.metronomeOverEmbedded.Pace(p0, p1)
}

// metronomePaceIgnoresCancel returns the planted defect of a Metronome
// whose Pace ignores the context over a fresh subject.
func metronomePaceIgnoresCancel(h *MetronomeHandle) Metronome {
	return &metronomePaceIgnoresCancelDefect{metronomeOverEmbedded: h.Fresh()}
}

// metronomePaceValueBesideDefect is the planted defect of a Metronome
// whose Pace returns a value beside an error.
type metronomePaceValueBesideDefect struct {
	metronomeOverEmbedded
}

// Pace is the override of the planted defect.
func (d *metronomePaceValueBesideDefect) Pace(p0 context.Context, p1 int) (r0 int, r1 error) {
	r0, r1 = d.metronomeOverEmbedded.Pace(p0, p1)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// metronomePaceValueBeside returns the planted defect of a Metronome whose
// Pace returns a value beside an error over a fresh subject.
func metronomePaceValueBeside(h *MetronomeHandle) Metronome {
	return &metronomePaceValueBesideDefect{metronomeOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Locker writes.
type (
	LockerHandle = conform.Handle[Locker]
	LockerCheck  = conform.Check[Locker]
	LockerOption = conform.Option[Locker]
)

// LockerChecks contains the name of each callable of Locker, with the
// identities of its checks.
var LockerChecks = struct {
	Put lockerPut
	Get lockerGet
}{
	Put: lockerPut{"Locker.Put"},
	Get: lockerGet{"Locker.Get"},
}

// lockerPut is the name of Locker.Put, with the identities of its checks.
type lockerPut struct {
	conform.Callable[Locker]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c lockerPut) Smoke() conform.ID[Locker] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c lockerPut) Cancel() conform.ID[Locker] {
	return c.Check("cancel")
}

// SyncedWritesSurvive returns the identity of the check that Get returns
// each value that Put wrote with success before a crash at any persistence
// point.
func (c lockerPut) SyncedWritesSurvive() conform.ID[Locker] {
	return c.Check("synced-writes-survive")
}

// NoTornRecord returns the identity of the check that Get returns each
// value of Put whole or not at all after a crash during a torn write.
func (c lockerPut) NoTornRecord() conform.ID[Locker] {
	return c.Check("no-torn-record")
}

// FullMediumIsReported returns the identity of the check that Put fails on
// a full medium, and Get returns each value as before the write.
func (c lockerPut) FullMediumIsReported() conform.ID[Locker] {
	return c.Check("full-medium-is-reported")
}

// CompletesAfterRecovery returns the identity of the check that each call
// of Locker completes within the bound of the world once every fault ends
// and the subject opens again.
func (c lockerPut) CompletesAfterRecovery() conform.ID[Locker] {
	return "Locker/completes-after-recovery"
}

// lockerGet is the name of Locker.Get, with the identities of its checks.
type lockerGet struct {
	conform.Callable[Locker]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c lockerGet) Smoke() conform.ID[Locker] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c lockerGet) Cancel() conform.ID[Locker] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c lockerGet) ZeroOnError() conform.ID[Locker] {
	return c.Check("zero-on-error")
}

// LockerSuite returns the checks of Locker, and the calls and the budgets
// of its callables.
func LockerSuite() conform.Suite[Locker] {
	return conform.Suite[Locker]{
		Declaration: "Locker",
		Package:     "example.com/shop/sealed",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Locker]{
			{
				ID:          LockerChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtPut0, builtPut1 := Entry()
						return h.Subject().Put(ctx, builtPut0, builtPut1)
					})
				},
				Defect: lockerPutPanics,
			},
			{
				ID:          LockerChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						builtPut0, builtPut1 := Entry()
						return h.Subject().Put(ctx, builtPut0, builtPut1)
					})
				},
				Defect: lockerPutIgnoresCancel,
			},
			{
				ID:          LockerChecks.Put.SyncedWritesSurvive(),
				Statement:   "Get returns each value that Put wrote with success before a crash at any persistence point",
				Claim:       spec.SyncedWritesSurvive,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *LockerHandle) {
					probe.SurvivesEveryCrashPoint(h, probe.Durable[Locker, string]{
						Write: func(ctx context.Context, s Locker, i int) error {
							if i == 0 {
								builtPut0, builtPut1 := Entry()
								return s.Put(ctx, builtPut0, builtPut1)
							}
							builtPut0, builtPut1 := Entry()
							return s.Put(ctx, builtPut0, builtPut1)
						},
						Read: func(ctx context.Context, s Locker, i int) (string, error) {
							if i == 0 {
								builtPut0, _ := Entry()
								return s.Get(ctx, builtPut0)
							}
							builtPut0, _ := Entry()
							return s.Get(ctx, builtPut0)
						},
						Want: func(i int) string {
							if i == 0 {
								return func() string {
									_, v := Entry()
									return v
								}()
							}
							return func() string {
								_, v := Entry()
								return v
							}()
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{func() string {
									v, _ := Entry()
									return v
								}()}
							}
							return []any{func() string {
								v, _ := Entry()
								return v
							}()}
						},
					})
				},
				Defect: lockerPutAcknowledgesBeforeSync,
				Doors:  conform.Doors[Locker]{Open: lockerPutAcknowledgesBeforeSyncOpen},
			},
			{
				ID:          LockerChecks.Put.NoTornRecord(),
				Statement:   "Get returns each value of Put whole or not at all after a crash during a torn write",
				Claim:       spec.NoTornRecord,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *LockerHandle) {
					probe.WholeAtEveryCrashPoint(h, probe.Durable[Locker, string]{
						Write: func(ctx context.Context, s Locker, i int) error {
							if i == 0 {
								builtPut0, builtPut1 := Entry()
								return s.Put(ctx, builtPut0, builtPut1)
							}
							builtPut0, builtPut1 := Entry()
							return s.Put(ctx, builtPut0, builtPut1)
						},
						Read: func(ctx context.Context, s Locker, i int) (string, error) {
							if i == 0 {
								builtPut0, _ := Entry()
								return s.Get(ctx, builtPut0)
							}
							builtPut0, _ := Entry()
							return s.Get(ctx, builtPut0)
						},
						Want: func(i int) string {
							if i == 0 {
								return func() string {
									_, v := Entry()
									return v
								}()
							}
							return func() string {
								_, v := Entry()
								return v
							}()
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{func() string {
									v, _ := Entry()
									return v
								}()}
							}
							return []any{func() string {
								v, _ := Entry()
								return v
							}()}
						},
					})
				},
				Defect: lockerPutTrustsAPartialRecord,
				Doors:  conform.Doors[Locker]{Open: lockerPutTrustsAPartialRecordOpen},
			},
			{
				ID:          LockerChecks.Put.FullMediumIsReported(),
				Statement:   "Put fails on a full medium, and Get returns each value as before the write",
				Claim:       spec.FullMediumIsReported,
				Machineries: []spec.MachineryID{spec.KillableProcessOrBreakableMedium},
				Apply: func(h *LockerHandle) {
					probe.RefusesOnAFullMedium(h, probe.Durable[Locker, string]{
						Write: func(ctx context.Context, s Locker, i int) error {
							if i == 0 {
								builtPut0, builtPut1 := Entry()
								return s.Put(ctx, builtPut0, builtPut1)
							}
							builtPut0, builtPut1 := Entry()
							return s.Put(ctx, builtPut0, builtPut1)
						},
						Read: func(ctx context.Context, s Locker, i int) (string, error) {
							if i == 0 {
								builtPut0, _ := Entry()
								return s.Get(ctx, builtPut0)
							}
							builtPut0, _ := Entry()
							return s.Get(ctx, builtPut0)
						},
						Want: func(i int) string {
							if i == 0 {
								return func() string {
									_, v := Entry()
									return v
								}()
							}
							return func() string {
								_, v := Entry()
								return v
							}()
						},
						Key: func(i int) any {
							if i == 0 {
								return []any{func() string {
									v, _ := Entry()
									return v
								}()}
							}
							return []any{func() string {
								v, _ := Entry()
								return v
							}()}
						},
					})
				},
				Defect: lockerPutHidesAFullMedium,
				Doors:  conform.Doors[Locker]{Open: lockerPutHidesAFullMediumOpen},
			},
			{
				ID:        LockerChecks.Put.CompletesAfterRecovery(),
				Statement: "each call of Locker completes within the bound of the world once every fault ends and the subject opens again",
				Claim:     spec.CompletesAfterRecovery,
				Machineries: []spec.MachineryID{
					spec.Quiescence,
					spec.GeneratedSequences,
					spec.KillableProcessOrBreakableMedium,
				},
				Property: func(h *LockerHandle, c *prop.Case) {
					probe.WorldRecovery(h, c, lockerWorld)
				},
				Defect:  lockerBlocksAfterRecovery,
				Doors:   conform.Doors[Locker]{Host: lockerBlocksAfterRecoveryHost},
				Around:  conform.Bubble,
				Options: probe.WorldOptions[Locker],
			},
			{
				ID:          LockerChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LockerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: lockerGetPanics,
			},
			{
				ID:          LockerChecks.Get.Smoke().Generated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: lockerGetPanics,
			},
			{
				ID:          LockerChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: lockerGetIgnoresCancel,
			},
			{
				ID:          LockerChecks.Get.Cancel().Generated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: lockerGetIgnoresCancel,
			},
			{
				ID:          LockerChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LockerHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, string(""))
				},
				Defect: lockerGetValueBeside,
			},
			{
				ID:          LockerChecks.Get.ZeroOnError().Generated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LockerHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, string(""))
				},
				Defect: lockerGetValueBeside,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Locker.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "mixin:sample",
				Callable:       "Locker.Put",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
		},
		Calls: []conform.Call[Locker]{
			{
				Calls:    lockerPutCalls,
				Over:     newLockerOver,
				Callable: LockerChecks.Put.Callable,
			},
			{
				Calls:    lockerGetCalls,
				Over:     newLockerOver,
				Callable: LockerChecks.Get.Callable,
			},
		},
	}
}

// RunLockerSuite applies the checks of Locker, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLockerSuite[T Locker](t *testing.T, newSubject func() T, opts ...LockerOption) {
	t.Helper()
	RunLockerSuiteFromHandle(t, func(*LockerHandle) T { return newSubject() }, opts...)
}

// RunLockerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLockerSuiteFromHandle[T Locker](
	t *testing.T,
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	t.Helper()
	conform.Run(t, LockerSuite(), lockerSubjects(newSubject), opts...)
}

// FuzzLockerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLockerSuite[T Locker](
	f *testing.F,
	callable conform.Callable[Locker],
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	f.Helper()
	conform.Fuzz(f, LockerSuite(), callable, lockerSubjects(newSubject), opts...)
}

// RunLockerBudgets applies the allocation ceilings of Locker and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLockerBudgets[T Locker](
	t *testing.T,
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	t.Helper()
	conform.RunBudgets(t, LockerSuite(), lockerSubjects(newSubject), opts...)
}

// BenchLockerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLockerSuite[T Locker](
	b conform.B,
	callable conform.Callable[Locker],
	newSubject func(h *LockerHandle) T,
	opts ...LockerOption,
) {
	b.Helper()
	conform.Bench(b, LockerSuite(), callable, lockerSubjects(newSubject), opts...)
}

// lockerSubjects returns the factory of the subjects of Locker that calls
// newSubject with the handle of each check.
func lockerSubjects[T Locker](newSubject func(h *LockerHandle) T) func(h *LockerHandle) Locker {
	return func(h *LockerHandle) Locker { return newSubject(h) }
}

// lockerPutCalls generates the calls of Locker.Put with the inputs that
// its builder Entry returns.
var lockerPutCalls = prop.Composite(func(*prop.Case) func(context.Context, Locker) error {
	return func(ctx context.Context, subject Locker) error {
		key, value := Entry()
		return subject.Put(ctx, key, value)
	}
})

// lockerGetCalls generates the calls of Locker.Get with the inputs that
// each case draws.
var lockerGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Locker) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Locker) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// lockerOverEmbedded is Locker under the name of the field that lockerOver
// embeds.
type lockerOverEmbedded = Locker

// lockerOver is a subject of Locker whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type lockerOver struct {
	lockerOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o lockerOver) Put(context.Context, string, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o lockerOver) Get(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newLockerOver returns the planted defect of a ceiling of Locker, whose
// callables call spend.
func newLockerOver(_ *LockerHandle, spend func()) Locker {
	return lockerOver{spend: spend}
}

// lockerPutPanicsDefect is the planted defect of a Locker whose Put
// panics.
type lockerPutPanicsDefect struct {
	lockerOverEmbedded
}

// Put is the override of the planted defect.
func (*lockerPutPanicsDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	panic("the planted defect of a Locker whose Put panics")
}

// lockerPutPanics returns the planted defect of a Locker whose Put panics
// over a fresh subject.
func lockerPutPanics(h *LockerHandle) Locker {
	return &lockerPutPanicsDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerPutIgnoresCancelDefect is the planted defect of a Locker whose Put
// ignores the context.
type lockerPutIgnoresCancelDefect struct {
	lockerOverEmbedded
}

// Put is the override of the planted defect.
func (d *lockerPutIgnoresCancelDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.lockerOverEmbedded.Put(p0, p1, p2)
}

// lockerPutIgnoresCancel returns the planted defect of a Locker whose Put
// ignores the context over a fresh subject.
func lockerPutIgnoresCancel(h *LockerHandle) Locker {
	return &lockerPutIgnoresCancelDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerPutAcknowledgesBeforeSync returns the planted defect of a Locker
// over a medium whose syncs persist nothing. The door of its check opens
// the defect over a new medium in memory.
func lockerPutAcknowledgesBeforeSync(h *LockerHandle) Locker {
	return h.Open(medium.NewMemory())
}

// lockerPutAcknowledgesBeforeSyncOpen opens the planted defect of a Locker
// over a medium whose syncs persist nothing. The application opens the
// subject over the medium m that probe.Unsynced wraps.
func lockerPutAcknowledgesBeforeSyncOpen(h *LockerHandle, m medium.FS) Locker {
	return h.FreshOpen(probe.Unsynced(m))
}

// lockerPutTrustsAPartialRecordDefect is the planted defect of a Locker
// whose Put writes the zero value under the key before its first value.
type lockerPutTrustsAPartialRecordDefect struct {
	lockerOverEmbedded
	wrote bool
}

// Put is the override of the planted defect.
func (d *lockerPutTrustsAPartialRecordDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	if !d.wrote {
		d.wrote = true
		r0 = d.lockerOverEmbedded.Put(p0, p1, *new(string))
		if r0 != nil {
			return r0
		}
	}
	return d.lockerOverEmbedded.Put(p0, p1, p2)
}

// lockerPutTrustsAPartialRecord returns the planted defect of a Locker
// whose Put writes the zero value under the key before its first value.
// The door of its check opens the defect over a new medium in memory.
func lockerPutTrustsAPartialRecord(h *LockerHandle) Locker {
	return h.Open(medium.NewMemory())
}

// lockerPutTrustsAPartialRecordOpen opens the planted defect of a Locker
// whose Put writes the zero value under the key before its first value.
// The defect wraps a subject that the application opens over the medium m.
func lockerPutTrustsAPartialRecordOpen(h *LockerHandle, m medium.FS) Locker {
	return &lockerPutTrustsAPartialRecordDefect{lockerOverEmbedded: h.FreshOpen(m)}
}

// lockerPutHidesAFullMedium returns the planted defect of a Locker over a
// medium that reports a full write as written. The door of its check opens
// the defect over a new medium in memory.
func lockerPutHidesAFullMedium(h *LockerHandle) Locker {
	return h.Open(medium.NewMemory())
}

// lockerPutHidesAFullMediumOpen opens the planted defect of a Locker over
// a medium that reports a full write as written. The application opens the
// subject over the medium m that probe.Unfilled wraps.
func lockerPutHidesAFullMediumOpen(h *LockerHandle, m medium.FS) Locker {
	return h.FreshOpen(probe.Unfilled(m))
}

// lockerBlocksAfterRecoveryDefect is the planted defect of a Locker whose
// each call after the first restart blocks until its case ends.
type lockerBlocksAfterRecoveryDefect struct {
	lockerOverEmbedded
	lifetime context.Context
}

// Put is the override of the planted defect.
func (d *lockerBlocksAfterRecoveryDefect) Put(p0 context.Context, p1, p2 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	if d.lifetime != nil {
		r0 = host.Park(p0)
		return r0
	}
	return d.lockerOverEmbedded.Put(p0, p1, p2)
}

// Get is the override of the planted defect.
func (d *lockerBlocksAfterRecoveryDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	if d.lifetime != nil {
		r1 = host.Park(p0)
		return r0, r1
	}
	return d.lockerOverEmbedded.Get(p0, p1)
}

// lockerBlocksAfterRecovery returns the subject of the application in
// place of the planted defect of a Locker whose each call after the first
// restart blocks until its case ends. The door Host of its check opens the
// defect in a world.
func lockerBlocksAfterRecovery(h *LockerHandle) Locker {
	return h.Fresh()
}

// lockerBlocksAfterRecoveryHost opens the planted defect of a Locker whose
// each call after the first restart blocks until its case ends. The defect
// wraps the subject that the application opens over the environment env of
// an incarnation. From the second incarnation of the host of env on, each
// callable of the defect blocks until its context ends.
func lockerBlocksAfterRecoveryHost(h *LockerHandle, env *host.Env) Locker {
	d := &lockerBlocksAfterRecoveryDefect{lockerOverEmbedded: h.FreshHost(env)}
	if env.Incarnation() > 1 {
		d.lifetime = env.Lifetime()
	}
	return d
}

// lockerGetPanicsDefect is the planted defect of a Locker whose Get
// panics.
type lockerGetPanicsDefect struct {
	lockerOverEmbedded
}

// Get is the override of the planted defect.
func (*lockerGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	panic("the planted defect of a Locker whose Get panics")
}

// lockerGetPanics returns the planted defect of a Locker whose Get panics
// over a fresh subject.
func lockerGetPanics(h *LockerHandle) Locker {
	return &lockerGetPanicsDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerGetIgnoresCancelDefect is the planted defect of a Locker whose Get
// ignores the context.
type lockerGetIgnoresCancelDefect struct {
	lockerOverEmbedded
}

// Get is the override of the planted defect.
func (d *lockerGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.lockerOverEmbedded.Get(p0, p1)
}

// lockerGetIgnoresCancel returns the planted defect of a Locker whose Get
// ignores the context over a fresh subject.
func lockerGetIgnoresCancel(h *LockerHandle) Locker {
	return &lockerGetIgnoresCancelDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerGetValueBesideDefect is the planted defect of a Locker whose Get
// returns a value beside an error.
type lockerGetValueBesideDefect struct {
	lockerOverEmbedded
}

// Get is the override of the planted defect.
func (d *lockerGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 string, r1 error) {
	r0, r1 = d.lockerOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// lockerGetValueBeside returns the planted defect of a Locker whose Get
// returns a value beside an error over a fresh subject.
func lockerGetValueBeside(h *LockerHandle) Locker {
	return &lockerGetValueBesideDefect{lockerOverEmbedded: h.Fresh()}
}

// lockerPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func lockerPutCall(ctx context.Context, subject Locker, args []any) ([]any, error) {
	a0, a1 := Entry()
	return nil, subject.Put(ctx, a0, a1)
}

// lockerGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func lockerGetCall(ctx context.Context, subject Locker, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// lockerReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func lockerReference(r Locker) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Put":
			values, err = lockerPutCall(ctx, r, args)
		case "Get":
			values, err = lockerGetCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// lockerWorld is the world of the subject of Locker, which the checks of
// its world probes run.
var lockerWorld = probe.World[Locker]{
	Operations: lockerOperations,
	Reference:  lockerReference,
}

// lockerOperations returns the operations of the world of Locker in the
// case c, over the pools of the keys that c chooses.
func lockerOperations(c *prop.Case) []probe.Operation[Locker] {
	keysString := probe.NewPool(
		c,
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []probe.Operation[Locker]{
		{
			Name: "Put",
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
			Call: func(ctx context.Context, s Locker, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := lockerPutCall(ctx, s, args)
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
			Call: func(ctx context.Context, s Locker, input any) (conform.Output, error) {
				args, _ := input.([]any)
				values, err := lockerGetCall(ctx, s, args)
				return conform.OutputOf(err, values...), err
			},
		},
	}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Shelf writes.
type (
	ShelfHandle = conform.Handle[Shelf]
	ShelfCheck  = conform.Check[Shelf]
	ShelfOption = conform.Option[Shelf]
)

// ShelfChecks contains the name of each callable of Shelf, with the
// identities of its checks.
var ShelfChecks = struct {
	Put shelfPut
	Get shelfGet
}{
	Put: shelfPut{"Shelf.Put"},
	Get: shelfGet{"Shelf.Get"},
}

// shelfPut is the name of Shelf.Put, with the identities of its checks.
type shelfPut struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c shelfPut) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c shelfPut) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// shelfGet is the name of Shelf.Get, with the identities of its checks.
type shelfGet struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c shelfGet) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c shelfGet) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c shelfGet) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c shelfGet) ReadAfterWrite() conform.ID[Shelf] {
	return c.Check("read-after-write")
}

// ReadsEachLastWrite returns the identity of the check that Get returns
// what Put wrote right before it, over any sequence of calls.
func (c shelfGet) ReadsEachLastWrite() conform.ID[Shelf] {
	return c.Check("reads-each-last-write")
}

// ShelfSuite returns the checks of Shelf, and the calls and the budgets of
// its callables.
func ShelfSuite() conform.Suite[Shelf] {
	return conform.Suite[Shelf]{
		Declaration: "Shelf",
		Package:     "example.com/shop/sealed",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Shelf]{
			{
				ID:          ShelfChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtPut0, builtPut1 := Shelved()
						return h.Subject().Put(ctx, builtPut0, builtPut1)
					})
				},
				Defect: shelfPutPanics,
			},
			{
				ID:          ShelfChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						builtPut0, builtPut1 := Shelved()
						return h.Subject().Put(ctx, builtPut0, builtPut1)
					})
				},
				Defect: shelfPutIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: shelfGetPanics,
			},
			{
				ID:          ShelfChecks.Get.Smoke().Generated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: shelfGetPanics,
			},
			{
				ID:          ShelfChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: shelfGetIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Get.Cancel().Generated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: shelfGetIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, int(0))
				},
				Defect: shelfGetValueBeside,
			},
			{
				ID:          ShelfChecks.Get.ZeroOnError().Generated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, int(0))
				},
				Defect: shelfGetValueBeside,
			},
			{
				ID:          ShelfChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (int, error) {
						builtPut0, _ := Shelved()
						return h.Subject().Get(ctx, builtPut0)
					}, func(ctx context.Context) (int, error) {
						builtPut0, builtPut1 := Shelved()
						w := builtPut1
						err := h.Subject().Put(ctx, builtPut0, w)
						return w, err
					})
				},
				Defect: shelfGetReadsAnotherValue,
			},
			{
				ID:          ShelfChecks.Get.ReadsEachLastWrite(),
				Statement:   "Get returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Shelf]{ShelfChecks.Get.ReadAfterWrite()},
				Property: func(h *ShelfHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						shelfActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							values, err := shelfGetCall(ctx, h.Subject(), []any{args[0]})
							return conform.OutputOf(err, values...)
						},
						func(args []any) int {
							v, _ := args[1].(int)
							return v
						},
					)
				},
				Defect: shelfGetReadsAStaleValueAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:composite-writer",
				Callable:       "Shelf.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "mixin:sample",
				Callable:       "Shelf.Put",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
		},
		Calls: []conform.Call[Shelf]{
			{
				Calls:    shelfPutCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Put.Callable,
			},
			{
				Calls:    shelfGetCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Get.Callable,
			},
		},
	}
}

// RunShelfSuite applies the checks of Shelf, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunShelfSuite[T Shelf](t *testing.T, newSubject func() T, opts ...ShelfOption) {
	t.Helper()
	RunShelfSuiteFromHandle(t, func(*ShelfHandle) T { return newSubject() }, opts...)
}

// RunShelfSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunShelfSuiteFromHandle[T Shelf](
	t *testing.T,
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	t.Helper()
	conform.Run(t, ShelfSuite(), shelfSubjects(newSubject), opts...)
}

// FuzzShelfSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzShelfSuite[T Shelf](
	f *testing.F,
	callable conform.Callable[Shelf],
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	f.Helper()
	conform.Fuzz(f, ShelfSuite(), callable, shelfSubjects(newSubject), opts...)
}

// RunShelfBudgets applies the allocation ceilings of Shelf and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunShelfBudgets[T Shelf](
	t *testing.T,
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	t.Helper()
	conform.RunBudgets(t, ShelfSuite(), shelfSubjects(newSubject), opts...)
}

// BenchShelfSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchShelfSuite[T Shelf](
	b conform.B,
	callable conform.Callable[Shelf],
	newSubject func(h *ShelfHandle) T,
	opts ...ShelfOption,
) {
	b.Helper()
	conform.Bench(b, ShelfSuite(), callable, shelfSubjects(newSubject), opts...)
}

// shelfSubjects returns the factory of the subjects of Shelf that calls
// newSubject with the handle of each check.
func shelfSubjects[T Shelf](newSubject func(h *ShelfHandle) T) func(h *ShelfHandle) Shelf {
	return func(h *ShelfHandle) Shelf { return newSubject(h) }
}

// shelfPutCalls generates the calls of Shelf.Put with the inputs that its
// builder Shelved returns.
var shelfPutCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(ctx context.Context, subject Shelf) error {
		key, n := Shelved()
		return subject.Put(ctx, key, n)
	}
})

// shelfGetCalls generates the calls of Shelf.Get with the inputs that each
// case draws.
var shelfGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// shelfOverEmbedded is Shelf under the name of the field that shelfOver
// embeds.
type shelfOverEmbedded = Shelf

// shelfOver is a subject of Shelf whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type shelfOver struct {
	shelfOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o shelfOver) Put(context.Context, string, int) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o shelfOver) Get(context.Context, string) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newShelfOver returns the planted defect of a ceiling of Shelf, whose
// callables call spend.
func newShelfOver(_ *ShelfHandle, spend func()) Shelf {
	return shelfOver{spend: spend}
}

// shelfPutPanicsDefect is the planted defect of a Shelf whose Put panics.
type shelfPutPanicsDefect struct {
	shelfOverEmbedded
}

// Put is the override of the planted defect.
func (*shelfPutPanicsDefect) Put(p0 context.Context, p1 string, p2 int) (r0 error) {
	panic("the planted defect of a Shelf whose Put panics")
}

// shelfPutPanics returns the planted defect of a Shelf whose Put panics
// over a fresh subject.
func shelfPutPanics(h *ShelfHandle) Shelf {
	return &shelfPutPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfPutIgnoresCancelDefect is the planted defect of a Shelf whose Put
// ignores the context.
type shelfPutIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Put is the override of the planted defect.
func (d *shelfPutIgnoresCancelDefect) Put(p0 context.Context, p1 string, p2 int) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.shelfOverEmbedded.Put(p0, p1, p2)
}

// shelfPutIgnoresCancel returns the planted defect of a Shelf whose Put
// ignores the context over a fresh subject.
func shelfPutIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfPutIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetPanicsDefect is the planted defect of a Shelf whose Get panics.
type shelfGetPanicsDefect struct {
	shelfOverEmbedded
}

// Get is the override of the planted defect.
func (*shelfGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 int, r1 error) {
	panic("the planted defect of a Shelf whose Get panics")
}

// shelfGetPanics returns the planted defect of a Shelf whose Get panics
// over a fresh subject.
func shelfGetPanics(h *ShelfHandle) Shelf {
	return &shelfGetPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetIgnoresCancelDefect is the planted defect of a Shelf whose Get
// ignores the context.
type shelfGetIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Get is the override of the planted defect.
func (d *shelfGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.Get(p0, p1)
}

// shelfGetIgnoresCancel returns the planted defect of a Shelf whose Get
// ignores the context over a fresh subject.
func shelfGetIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfGetIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetValueBesideDefect is the planted defect of a Shelf whose Get
// returns a value beside an error.
type shelfGetValueBesideDefect struct {
	shelfOverEmbedded
}

// Get is the override of the planted defect.
func (d *shelfGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 int, r1 error) {
	r0, r1 = d.shelfOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// shelfGetValueBeside returns the planted defect of a Shelf whose Get
// returns a value beside an error over a fresh subject.
func shelfGetValueBeside(h *ShelfHandle) Shelf {
	return &shelfGetValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetReadsAnotherValueDefect is the planted defect of a Shelf whose
// Get returns a value that Put did not write.
type shelfGetReadsAnotherValueDefect struct {
	shelfOverEmbedded
}

// Get is the override of the planted defect.
func (d *shelfGetReadsAnotherValueDefect) Get(p0 context.Context, p1 string) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.shelfOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = int(7)
	}
	return r0, r1
}

// shelfGetReadsAnotherValue returns the planted defect of a Shelf whose
// Get returns a value that Put did not write over a fresh subject.
func shelfGetReadsAnotherValue(h *ShelfHandle) Shelf {
	return &shelfGetReadsAnotherValueDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetReadsAStaleValueAfterThreeDefect is the planted defect of a
// Shelf whose Get returns its first value again from its fourth call on.
type shelfGetReadsAStaleValueAfterThreeDefect struct {
	shelfOverEmbedded
	calls int
	first *int
}

// Get is the override of the planted defect.
func (d *shelfGetReadsAStaleValueAfterThreeDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.shelfOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// shelfGetReadsAStaleValueAfterThree returns the planted defect of a Shelf
// whose Get returns its first value again from its fourth call on over a
// fresh subject.
func shelfGetReadsAStaleValueAfterThree(h *ShelfHandle) Shelf {
	return &shelfGetReadsAStaleValueAfterThreeDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfActions returns the actions of the generated sequences of Shelf
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func shelfActions(c *prop.Case, subject func() Shelf) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Put", args)
				values, err := shelfPutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args, args[0])
				values, err := shelfGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// shelfPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfPutCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, a1 := Shelved()
	return nil, subject.Put(ctx, a0, a1)
}

// shelfGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfGetCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Mixer writes.
type (
	MixerHandle = conform.Handle[Mixer]
	MixerCheck  = conform.Check[Mixer]
	MixerOption = conform.Option[Mixer]
)

// MixerChecks contains the name of each callable of Mixer, with the
// identities of its checks.
var MixerChecks = struct {
	Mix mixerMix
}{
	Mix: mixerMix{"Mixer.Mix"},
}

// mixerMix is the name of Mixer.Mix, with the identities of its checks.
type mixerMix struct {
	conform.Callable[Mixer]
}

// Smoke returns the identity of the check that Mix returns without a panic
// for derived inputs.
func (c mixerMix) Smoke() conform.ID[Mixer] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Mix returns an error for a
// cancelled context.
func (c mixerMix) Cancel() conform.ID[Mixer] {
	return c.Check("cancel")
}

// MixerSuite returns the checks of Mixer, and the calls and the budgets of
// its callables.
func MixerSuite() conform.Suite[Mixer] {
	return conform.Suite[Mixer]{
		Declaration: "Mixer",
		Package:     "example.com/shop/sealed",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Mixer]{
			{
				ID:          MixerChecks.Mix.Smoke(),
				Statement:   "Mix returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MixerHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtMix0, builtMix1, builtMix2, builtMix3, builtMix4 := Mixed()
						return h.Subject().
							Mix(ctx, builtMix0, builtMix1, builtMix2, builtMix3, builtMix4)
					})
				},
				Defect: mixerMixPanics,
			},
			{
				ID:          MixerChecks.Mix.Cancel(),
				Statement:   "Mix returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MixerHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						builtMix0, builtMix1, builtMix2, builtMix3, builtMix4 := Mixed()
						return h.Subject().
							Mix(ctx, builtMix0, builtMix1, builtMix2, builtMix3, builtMix4)
					})
				},
				Defect: mixerMixIgnoresCancel,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:multi-arg-writer",
				Callable:       "Mixer.Mix",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "mixin:sample",
				Callable:       "Mixer.Mix",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
		},
		Calls: []conform.Call[Mixer]{
			{
				Calls:    mixerMixCalls,
				Over:     newMixerOver,
				Callable: MixerChecks.Mix.Callable,
			},
		},
	}
}

// RunMixerSuite applies the checks of Mixer, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunMixerSuite[T Mixer](t *testing.T, newSubject func() T, opts ...MixerOption) {
	t.Helper()
	RunMixerSuiteFromHandle(t, func(*MixerHandle) T { return newSubject() }, opts...)
}

// RunMixerSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunMixerSuiteFromHandle[T Mixer](
	t *testing.T,
	newSubject func(h *MixerHandle) T,
	opts ...MixerOption,
) {
	t.Helper()
	conform.Run(t, MixerSuite(), mixerSubjects(newSubject), opts...)
}

// FuzzMixerSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzMixerSuite[T Mixer](
	f *testing.F,
	callable conform.Callable[Mixer],
	newSubject func(h *MixerHandle) T,
	opts ...MixerOption,
) {
	f.Helper()
	conform.Fuzz(f, MixerSuite(), callable, mixerSubjects(newSubject), opts...)
}

// RunMixerBudgets applies the allocation ceilings of Mixer and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunMixerBudgets[T Mixer](
	t *testing.T,
	newSubject func(h *MixerHandle) T,
	opts ...MixerOption,
) {
	t.Helper()
	conform.RunBudgets(t, MixerSuite(), mixerSubjects(newSubject), opts...)
}

// BenchMixerSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchMixerSuite[T Mixer](
	b conform.B,
	callable conform.Callable[Mixer],
	newSubject func(h *MixerHandle) T,
	opts ...MixerOption,
) {
	b.Helper()
	conform.Bench(b, MixerSuite(), callable, mixerSubjects(newSubject), opts...)
}

// mixerSubjects returns the factory of the subjects of Mixer that calls
// newSubject with the handle of each check.
func mixerSubjects[T Mixer](newSubject func(h *MixerHandle) T) func(h *MixerHandle) Mixer {
	return func(h *MixerHandle) Mixer { return newSubject(h) }
}

// mixerMixCalls generates the calls of Mixer.Mix with the inputs that its
// builder Mixed returns.
var mixerMixCalls = prop.Composite(func(*prop.Case) func(context.Context, Mixer) error {
	return func(ctx context.Context, subject Mixer) error {
		a1, time0, subject0, subject00, track := Mixed()
		return subject.Mix(ctx, a1, time0, subject0, subject00, track)
	}
})

// mixerOverEmbedded is Mixer under the name of the field that mixerOver
// embeds.
type mixerOverEmbedded = Mixer

// mixerOver is a subject of Mixer whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type mixerOver struct {
	mixerOverEmbedded
	spend func()
}

// Mix calls spend and returns its zero results.
func (o mixerOver) Mix(context.Context, int, time.Duration, string, string, *Track) (r0 error) {
	o.spend()
	return r0
}

// newMixerOver returns the planted defect of a ceiling of Mixer, whose
// callables call spend.
func newMixerOver(_ *MixerHandle, spend func()) Mixer {
	return mixerOver{spend: spend}
}

// mixerMixPanicsDefect is the planted defect of a Mixer whose Mix panics.
type mixerMixPanicsDefect struct {
	mixerOverEmbedded
}

// Mix is the override of the planted defect.
func (*mixerMixPanicsDefect) Mix(
	p0 context.Context,
	p1 int,
	p2 time.Duration,
	p3, p4 string,
	p5 *Track,
) (r0 error) {
	panic("the planted defect of a Mixer whose Mix panics")
}

// mixerMixPanics returns the planted defect of a Mixer whose Mix panics
// over a fresh subject.
func mixerMixPanics(h *MixerHandle) Mixer {
	return &mixerMixPanicsDefect{mixerOverEmbedded: h.Fresh()}
}

// mixerMixIgnoresCancelDefect is the planted defect of a Mixer whose Mix
// ignores the context.
type mixerMixIgnoresCancelDefect struct {
	mixerOverEmbedded
}

// Mix is the override of the planted defect.
func (d *mixerMixIgnoresCancelDefect) Mix(
	p0 context.Context,
	p1 int,
	p2 time.Duration,
	p3, p4 string,
	p5 *Track,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.mixerOverEmbedded.Mix(p0, p1, p2, p3, p4, p5)
}

// mixerMixIgnoresCancel returns the planted defect of a Mixer whose Mix
// ignores the context over a fresh subject.
func mixerMixIgnoresCancel(h *MixerHandle) Mixer {
	return &mixerMixIgnoresCancelDefect{mixerOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Depot writes.
type (
	DepotHandle = conform.Handle[Depot]
	DepotCheck  = conform.Check[Depot]
	DepotOption = conform.Option[Depot]
)

// DepotChecks contains the name of each callable of Depot, with the
// identities of its checks.
var DepotChecks = struct {
	Stock depotStock
	Find  depotFind
}{
	Stock: depotStock{"Depot.Stock"},
	Find:  depotFind{"Depot.Find"},
}

// depotStock is the name of Depot.Stock, with the identities of its
// checks.
type depotStock struct {
	conform.Callable[Depot]
}

// Smoke returns the identity of the check that Stock returns without a
// panic for derived inputs.
func (c depotStock) Smoke() conform.ID[Depot] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Stock returns an error for
// a cancelled context.
func (c depotStock) Cancel() conform.ID[Depot] {
	return c.Check("cancel")
}

// depotFind is the name of Depot.Find, with the identities of its checks.
type depotFind struct {
	conform.Callable[Depot]
}

// Smoke returns the identity of the check that Find returns without a
// panic for derived inputs.
func (c depotFind) Smoke() conform.ID[Depot] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Find returns an error for
// a cancelled context.
func (c depotFind) Cancel() conform.ID[Depot] {
	return c.Check("cancel")
}

// ReadAfterWrite returns the identity of the check that Find returns what
// Stock stored before it.
func (c depotFind) ReadAfterWrite() conform.ID[Depot] {
	return c.Check("read-after-write")
}

// ReadsEachLastWrite returns the identity of the check that Find returns
// what Stock wrote right before it, over any sequence of calls.
func (c depotFind) ReadsEachLastWrite() conform.ID[Depot] {
	return c.Check("reads-each-last-write")
}

// DepotSuite returns the checks of Depot, and the calls and the budgets of
// its callables.
func DepotSuite() conform.Suite[Depot] {
	return conform.Suite[Depot]{
		Declaration: "Depot",
		Package:     "example.com/shop/sealed",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Depot]{
			{
				ID:          DepotChecks.Stock.Smoke(),
				Statement:   "Stock returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DepotHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtStock0 := Stocked()
						return h.Subject().Stock(ctx, builtStock0)
					})
				},
				Defect: depotStockPanics,
			},
			{
				ID:          DepotChecks.Stock.Cancel(),
				Statement:   "Stock returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DepotHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						builtStock0 := Stocked()
						return h.Subject().Stock(ctx, builtStock0)
					})
				},
				Defect: depotStockIgnoresCancel,
			},
			{
				ID:          DepotChecks.Find.Smoke(),
				Statement:   "Find returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DepotHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						builtFind0, builtFind1 := Found()
						_, err := h.Subject().Find(ctx, builtFind0, builtFind1)
						return err
					})
				},
				Defect: depotFindPanics,
			},
			{
				ID:          DepotChecks.Find.Cancel(),
				Statement:   "Find returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DepotHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						builtFind0, builtFind1 := Found()
						_, err := h.Subject().Find(ctx, builtFind0, builtFind1)
						return err
					})
				},
				Defect: depotFindIgnoresCancel,
			},
			{
				ID:          DepotChecks.Find.ReadAfterWrite(),
				Statement:   "Find returns what Stock stored before it",
				Claim:       spec.ReadAfterWrite,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DepotHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						builtStock0 := Stocked()
						_, builtFind1 := Found()
						return h.Subject().Find(ctx, builtStock0.SKU, builtFind1)
					}, func(ctx context.Context) (Item, error) {
						builtStock0 := Stocked()
						w := builtStock0
						err := h.Subject().Stock(ctx, w)
						return w, err
					})
				},
				Defect: depotFindReadsAnotherValue,
			},
			{
				ID:          DepotChecks.Find.ReadsEachLastWrite(),
				Statement:   "Find returns what Stock wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Depot]{DepotChecks.Find.ReadAfterWrite()},
				Property: func(h *DepotHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						depotActions,
						"Stock",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := depotFindCall(
								ctx,
								h.Subject(),
								[]any{w.SKU, func(Item) bool { return false }},
							)
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: depotFindReadsAStaleValueAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Depot.Stock",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "mixin:sample",
				Callable:       "Depot.Stock",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
			{
				Classification: "mixin:sample",
				Callable:       "Depot.Find",
				Stance:         spec.Silent,
				Note:           "The mixin names the builder of the inputs of each derived call of the callable.",
			},
		},
		Calls: []conform.Call[Depot]{
			{
				Calls:    depotStockCalls,
				Over:     newDepotOver,
				Callable: DepotChecks.Stock.Callable,
			},
			{
				Calls:    depotFindCalls,
				Over:     newDepotOver,
				Callable: DepotChecks.Find.Callable,
			},
		},
	}
}

// RunDepotSuite applies the checks of Depot, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunDepotSuite[T Depot](t *testing.T, newSubject func() T, opts ...DepotOption) {
	t.Helper()
	RunDepotSuiteFromHandle(t, func(*DepotHandle) T { return newSubject() }, opts...)
}

// RunDepotSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDepotSuiteFromHandle[T Depot](
	t *testing.T,
	newSubject func(h *DepotHandle) T,
	opts ...DepotOption,
) {
	t.Helper()
	conform.Run(t, DepotSuite(), depotSubjects(newSubject), opts...)
}

// FuzzDepotSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzDepotSuite[T Depot](
	f *testing.F,
	callable conform.Callable[Depot],
	newSubject func(h *DepotHandle) T,
	opts ...DepotOption,
) {
	f.Helper()
	conform.Fuzz(f, DepotSuite(), callable, depotSubjects(newSubject), opts...)
}

// RunDepotBudgets applies the allocation ceilings of Depot and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunDepotBudgets[T Depot](
	t *testing.T,
	newSubject func(h *DepotHandle) T,
	opts ...DepotOption,
) {
	t.Helper()
	conform.RunBudgets(t, DepotSuite(), depotSubjects(newSubject), opts...)
}

// BenchDepotSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchDepotSuite[T Depot](
	b conform.B,
	callable conform.Callable[Depot],
	newSubject func(h *DepotHandle) T,
	opts ...DepotOption,
) {
	b.Helper()
	conform.Bench(b, DepotSuite(), callable, depotSubjects(newSubject), opts...)
}

// depotSubjects returns the factory of the subjects of Depot that calls
// newSubject with the handle of each check.
func depotSubjects[T Depot](newSubject func(h *DepotHandle) T) func(h *DepotHandle) Depot {
	return func(h *DepotHandle) Depot { return newSubject(h) }
}

// depotStockCalls generates the calls of Depot.Stock with the inputs that
// its builder Stocked returns.
var depotStockCalls = prop.Composite(func(*prop.Case) func(context.Context, Depot) error {
	return func(ctx context.Context, subject Depot) error {
		item := Stocked()
		return subject.Stock(ctx, item)
	}
})

// depotFindCalls generates the calls of Depot.Find with the inputs that
// its builder Found returns.
var depotFindCalls = prop.Composite(func(*prop.Case) func(context.Context, Depot) error {
	return func(ctx context.Context, subject Depot) error {
		sku, match := Found()
		_, err := subject.Find(ctx, sku, match)
		return err
	}
})

// depotOverEmbedded is Depot under the name of the field that depotOver
// embeds.
type depotOverEmbedded = Depot

// depotOver is a subject of Depot whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type depotOver struct {
	depotOverEmbedded
	spend func()
}

// Stock calls spend and returns its zero results.
func (o depotOver) Stock(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Find calls spend and returns its zero results.
func (o depotOver) Find(context.Context, string, func(Item) bool) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newDepotOver returns the planted defect of a ceiling of Depot, whose
// callables call spend.
func newDepotOver(_ *DepotHandle, spend func()) Depot {
	return depotOver{spend: spend}
}

// depotStockPanicsDefect is the planted defect of a Depot whose Stock
// panics.
type depotStockPanicsDefect struct {
	depotOverEmbedded
}

// Stock is the override of the planted defect.
func (*depotStockPanicsDefect) Stock(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Depot whose Stock panics")
}

// depotStockPanics returns the planted defect of a Depot whose Stock
// panics over a fresh subject.
func depotStockPanics(h *DepotHandle) Depot {
	return &depotStockPanicsDefect{depotOverEmbedded: h.Fresh()}
}

// depotStockIgnoresCancelDefect is the planted defect of a Depot whose
// Stock ignores the context.
type depotStockIgnoresCancelDefect struct {
	depotOverEmbedded
}

// Stock is the override of the planted defect.
func (d *depotStockIgnoresCancelDefect) Stock(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.depotOverEmbedded.Stock(p0, p1)
}

// depotStockIgnoresCancel returns the planted defect of a Depot whose
// Stock ignores the context over a fresh subject.
func depotStockIgnoresCancel(h *DepotHandle) Depot {
	return &depotStockIgnoresCancelDefect{depotOverEmbedded: h.Fresh()}
}

// depotFindPanicsDefect is the planted defect of a Depot whose Find
// panics.
type depotFindPanicsDefect struct {
	depotOverEmbedded
}

// Find is the override of the planted defect.
func (*depotFindPanicsDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	panic("the planted defect of a Depot whose Find panics")
}

// depotFindPanics returns the planted defect of a Depot whose Find panics
// over a fresh subject.
func depotFindPanics(h *DepotHandle) Depot {
	return &depotFindPanicsDefect{depotOverEmbedded: h.Fresh()}
}

// depotFindIgnoresCancelDefect is the planted defect of a Depot whose Find
// ignores the context.
type depotFindIgnoresCancelDefect struct {
	depotOverEmbedded
}

// Find is the override of the planted defect.
func (d *depotFindIgnoresCancelDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.depotOverEmbedded.Find(p0, p1, p2)
}

// depotFindIgnoresCancel returns the planted defect of a Depot whose Find
// ignores the context over a fresh subject.
func depotFindIgnoresCancel(h *DepotHandle) Depot {
	return &depotFindIgnoresCancelDefect{depotOverEmbedded: h.Fresh()}
}

// depotFindReadsAnotherValueDefect is the planted defect of a Depot whose
// Find returns a value that Stock did not write.
type depotFindReadsAnotherValueDefect struct {
	depotOverEmbedded
}

// Find is the override of the planted defect.
func (d *depotFindReadsAnotherValueDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.depotOverEmbedded.Find(p0, p1, p2)
	if r1 == nil {
		r0 = Item{SKU: "SKU-b"}
	}
	return r0, r1
}

// depotFindReadsAnotherValue returns the planted defect of a Depot whose
// Find returns a value that Stock did not write over a fresh subject.
func depotFindReadsAnotherValue(h *DepotHandle) Depot {
	return &depotFindReadsAnotherValueDefect{depotOverEmbedded: h.Fresh()}
}

// depotFindReadsAStaleValueAfterThreeDefect is the planted defect of a
// Depot whose Find returns its first value again from its fourth call on.
type depotFindReadsAStaleValueAfterThreeDefect struct {
	depotOverEmbedded
	calls int
	first *Item
}

// Find is the override of the planted defect.
func (d *depotFindReadsAStaleValueAfterThreeDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.depotOverEmbedded.Find(p0, p1, p2)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// depotFindReadsAStaleValueAfterThree returns the planted defect of a
// Depot whose Find returns its first value again from its fourth call on
// over a fresh subject.
func depotFindReadsAStaleValueAfterThree(h *DepotHandle) Depot {
	return &depotFindReadsAStaleValueAfterThreeDefect{depotOverEmbedded: h.Fresh()}
}

// depotActions returns the actions of the generated sequences of Depot
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func depotActions(c *prop.Case, subject func() Depot) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Stock",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Stock", args)
				values, err := depotStockCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Find",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Find", args)
				values, err := depotFindCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// depotStockCall calls Stock on subject with the arguments args, and
// returns the values of the call beside its error.
func depotStockCall(ctx context.Context, subject Depot, args []any) ([]any, error) {
	a0 := Stocked()
	return nil, subject.Stock(ctx, a0)
}

// depotFindCall calls Find on subject with the arguments args, and returns
// the values of the call beside its error.
func depotFindCall(ctx context.Context, subject Depot, args []any) ([]any, error) {
	a0, a1 := Found()
	r0, err := subject.Find(ctx, a0, a1)
	return []any{r0}, err
}
