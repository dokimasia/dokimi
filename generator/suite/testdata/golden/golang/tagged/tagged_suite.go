package tagged

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Hasher writes.
type (
	HasherHandle = conform.Handle[Hasher]
	HasherCheck  = conform.Check[Hasher]
	HasherOption = conform.Option[Hasher]
)

// HasherChecks contains the name of each callable of Hasher, with the
// identities of its checks.
var HasherChecks = struct {
	HashTagged    hasherHashTagged
	CombineTagged hasherCombineTagged
	HashRoles     hasherHashRoles
	Truncate      hasherTruncate
	Seal          hasherSeal
}{
	HashTagged:    hasherHashTagged{"Hasher.HashTagged"},
	CombineTagged: hasherCombineTagged{"Hasher.CombineTagged"},
	HashRoles:     hasherHashRoles{"Hasher.HashRoles"},
	Truncate:      hasherTruncate{"Hasher.Truncate"},
	Seal:          hasherSeal{"Hasher.Seal"},
}

// hasherHashTagged is the name of Hasher.HashTagged, with the identities
// of its checks.
type hasherHashTagged struct {
	conform.Callable[Hasher]
}

// hasherCombineTagged is the name of Hasher.CombineTagged, with the
// identities of its checks.
type hasherCombineTagged struct {
	conform.Callable[Hasher]
}

// hasherHashRoles is the name of Hasher.HashRoles, with the identities of
// its checks.
type hasherHashRoles struct {
	conform.Callable[Hasher]
}

// hasherTruncate is the name of Hasher.Truncate, with the identities of
// its checks.
type hasherTruncate struct {
	conform.Callable[Hasher]
}

// hasherSeal is the name of Hasher.Seal, with the identities of its
// checks.
type hasherSeal struct {
	conform.Callable[Hasher]
}

// HasherSuite returns the checks of Hasher, and the calls and the budgets
// of its callables.
func HasherSuite() conform.Suite[Hasher] {
	return conform.Suite[Hasher]{
		Declaration: "Hasher",
		Calls: []conform.Call[Hasher]{
			{
				Calls:    hasherHashTaggedCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.HashTagged.Callable,
			},
			{
				Calls:    hasherCombineTaggedCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.CombineTagged.Callable,
			},
			{
				Calls:    hasherHashRolesCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.HashRoles.Callable,
			},
			{
				Calls:    hasherTruncateCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.Truncate.Callable,
			},
			{
				Calls:    hasherSealCalls,
				Over:     newHasherOver,
				Callable: HasherChecks.Seal.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(hasherSizeValues),
		},
	}
}

// RunHasherSuite applies the checks of Hasher, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunHasherSuite[T Hasher](t *testing.T, newSubject func() T, opts ...HasherOption) {
	t.Helper()
	RunHasherSuiteFromHandle(t, func(*HasherHandle) T { return newSubject() }, opts...)
}

// RunHasherSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunHasherSuiteFromHandle[T Hasher](
	t *testing.T,
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	t.Helper()
	conform.Run(t, HasherSuite(), hasherSubjects(newSubject), opts...)
}

// FuzzHasherSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzHasherSuite[T Hasher](
	f *testing.F,
	callable conform.Callable[Hasher],
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	f.Helper()
	conform.Fuzz(f, HasherSuite(), callable, hasherSubjects(newSubject), opts...)
}

// RunHasherBudgets applies the allocation ceilings of Hasher and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunHasherBudgets[T Hasher](
	t *testing.T,
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	t.Helper()
	conform.RunBudgets(t, HasherSuite(), hasherSubjects(newSubject), opts...)
}

// BenchHasherSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchHasherSuite[T Hasher](
	b conform.B,
	callable conform.Callable[Hasher],
	newSubject func(h *HasherHandle) T,
	opts ...HasherOption,
) {
	b.Helper()
	conform.Bench(b, HasherSuite(), callable, hasherSubjects(newSubject), opts...)
}

// hasherSubjects returns the factory of the subjects of Hasher that calls
// newSubject with the handle of each check.
func hasherSubjects[T Hasher](newSubject func(h *HasherHandle) T) func(h *HasherHandle) Hasher {
	return func(h *HasherHandle) Hasher { return newSubject(h) }
}

// hasherHashTaggedCalls generates the calls of Hasher.HashTagged with the
// inputs that each case draws.
var hasherHashTaggedCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hasher) error {
	r := Role(1)
	data := c.Draw(prop.Of[[]byte](), "data")
	return func(_ context.Context, subject Hasher) error {
		_ = subject.HashTagged(r, data)
		return nil
	}
})

// hasherCombineTaggedCalls generates the calls of Hasher.CombineTagged
// with the inputs that each case draws.
var hasherCombineTaggedCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, Hasher) error {
		r := Role(129)
		left := c.Draw(prop.Of[[32]byte](), "left")
		right := c.Draw(prop.Of[[32]byte](), "right")
		return func(_ context.Context, subject Hasher) error {
			_ = subject.CombineTagged(r, left, right)
			return nil
		}
	},
)

// hasherHashRolesCalls generates the calls of Hasher.HashRoles with the
// inputs that each case draws.
var hasherHashRolesCalls = prop.Composite(func(*prop.Case) func(context.Context, Hasher) error {
	roles := []Role{Role(2)}
	return func(_ context.Context, subject Hasher) error {
		_ = subject.HashRoles(roles...)
		return nil
	}
})

// hasherTruncateCalls generates the calls of Hasher.Truncate with the
// inputs that each case draws.
var hasherTruncateCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hasher) error {
	d := c.Draw(prop.Of[[32]byte](), "d")
	n := c.Draw(prop.Of[Size](), "n")
	return func(_ context.Context, subject Hasher) error {
		_ = subject.Truncate(d, n)
		return nil
	}
})

// hasherSealCalls generates the calls of Hasher.Seal with the inputs that
// each case draws.
var hasherSealCalls = prop.Composite(func(c *prop.Case) func(context.Context, Hasher) error {
	d := c.Draw(prop.Of[[32]byte](), "d")
	return func(ctx context.Context, subject Hasher) error {
		return subject.Seal(ctx, d)
	}
})

// hasherSizeValues generates the values of Size that keep or change its
// starting value.
var hasherSizeValues = prop.Composite(func(c *prop.Case) Size {
	changed := c.Draw(prop.Of[Size](), "size")
	if c.Draw(prop.Boolean(), "size.changed") {
		return changed
	}
	return 16
})

// hasherOverEmbedded is Hasher under the name of the field that hasherOver
// embeds.
type hasherOverEmbedded = Hasher

// hasherOver is a subject of Hasher whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type hasherOver struct {
	hasherOverEmbedded
	spend func()
}

// HashTagged calls spend and returns its zero results.
func (o hasherOver) HashTagged(Role, []byte) (r0 [32]byte) {
	o.spend()
	return r0
}

// CombineTagged calls spend and returns its zero results.
func (o hasherOver) CombineTagged(Role, [32]byte, [32]byte) (r0 [32]byte) {
	o.spend()
	return r0
}

// HashRoles calls spend and returns its zero results.
func (o hasherOver) HashRoles(...Role) (r0 [32]byte) {
	o.spend()
	return r0
}

// Truncate calls spend and returns its zero results.
func (o hasherOver) Truncate([32]byte, Size) (r0 []byte) {
	o.spend()
	return r0
}

// Seal calls spend and returns its zero results.
func (o hasherOver) Seal(context.Context, [32]byte) (r0 error) {
	o.spend()
	return r0
}

// newHasherOver returns the planted defect of a ceiling of Hasher, whose
// callables call spend.
func newHasherOver(_ *HasherHandle, spend func()) Hasher {
	return hasherOver{spend: spend}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Digest writes.
type (
	DigestHandle = conform.Handle[DigestFunc]
	DigestCheck  = conform.Check[DigestFunc]
	DigestOption = conform.Option[DigestFunc]
)

// DigestFunc is the type of the function Digest. Each subject of the suite
// of Digest has it.
type DigestFunc func(r Role, data []byte) [32]byte

// DigestChecks contains the name of each callable of Digest, with the
// identities of its checks.
var DigestChecks = digestCallable{"Digest"}

// digestCallable is the name of Digest, with the identities of its checks.
type digestCallable struct {
	conform.Callable[DigestFunc]
}

// DigestSuite returns the checks of Digest, and the calls and the budgets
// of its callables.
func DigestSuite() conform.Suite[DigestFunc] {
	return conform.Suite[DigestFunc]{
		Declaration: "Digest",
		Calls: []conform.Call[DigestFunc]{
			{
				Calls:    digestCalls,
				Over:     newDigestOver,
				Callable: DigestChecks.Callable,
			},
		},
	}
}

// RunDigestSuite applies the checks of Digest, and the checks that opts
// add, to subject.
func RunDigestSuite(t *testing.T, subject DigestFunc, opts ...DigestOption) {
	t.Helper()
	RunDigestSuiteFromHandle(t, digestSubjects(subject), opts...)
}

// RunDigestSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDigestSuiteFromHandle(
	t *testing.T,
	newSubject func(h *DigestHandle) DigestFunc,
	opts ...DigestOption,
) {
	t.Helper()
	conform.Run(t, DigestSuite(), newSubject, opts...)
}

// FuzzDigestSuite registers on f the fuzz target of the property checks of
// Digest on subject.
func FuzzDigestSuite(f *testing.F, subject DigestFunc, opts ...DigestOption) {
	f.Helper()
	conform.Fuzz(f, DigestSuite(), DigestChecks.Callable, digestSubjects(subject), opts...)
}

// RunDigestBudgets applies the allocation ceiling of Digest and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunDigestBudgets(t *testing.T, subject DigestFunc, opts ...DigestOption) {
	t.Helper()
	conform.RunBudgets(t, DigestSuite(), digestSubjects(subject), opts...)
}

// BenchDigestSuite measures subject under the ceilings of the budget of
// Digest.
func BenchDigestSuite(b conform.B, subject DigestFunc, opts ...DigestOption) {
	b.Helper()
	conform.Bench(b, DigestSuite(), DigestChecks.Callable, digestSubjects(subject), opts...)
}

// digestSubjects returns the factory of the subjects of Digest that
// returns subject for each handle.
func digestSubjects(subject DigestFunc) func(*DigestHandle) DigestFunc {
	return func(*DigestHandle) DigestFunc { return subject }
}

// digestCalls generates the calls of Digest with the inputs that each case
// draws.
var digestCalls = prop.Composite(func(c *prop.Case) func(context.Context, DigestFunc) error {
	r := Role(3)
	data := c.Draw(prop.Of[[]byte](), "data")
	return func(_ context.Context, subject DigestFunc) error {
		_ = subject(r, data)
		return nil
	}
})

// newDigestOver returns the planted defect of a ceiling of Digest, which
// calls spend and returns its zero results.
func newDigestOver(_ *DigestHandle, spend func()) DigestFunc {
	return func(Role, []byte) (r0 [32]byte) {
		spend()
		return r0
	}
}
