package durable

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Keeper writes.
type (
	KeeperHandle = conform.Handle[Keeper]
	KeeperCheck  = conform.Check[Keeper]
	KeeperOption = conform.Option[Keeper]
)

// KeeperChecks contains the name of each callable of Keeper, with the
// identities of its checks.
var KeeperChecks = struct {
	Put keeperPut
	Get keeperGet
}{
	Put: keeperPut{"Keeper.Put"},
	Get: keeperGet{"Keeper.Get"},
}

// keeperPut is the name of Keeper.Put, with the identities of its checks.
type keeperPut struct {
	conform.Callable[Keeper]
}

// LeavesTheWholeValue returns the identity of the check that a crash
// during Put leaves the whole value or no value.
func (c keeperPut) LeavesTheWholeValue() conform.ID[Keeper] {
	return c.Check("leaves-the-whole-value")
}

// keeperGet is the name of Keeper.Get, with the identities of its checks.
type keeperGet struct {
	conform.Callable[Keeper]
}

// KeepsEachPutAcrossARestart returns the identity of the check that Get
// returns each value that Put stored before a restart.
func (c keeperGet) KeepsEachPutAcrossARestart() conform.ID[Keeper] {
	return c.Check("keeps-each-put-across-a-restart")
}

// KeepsTheLastPut returns the identity of the check that Get returns the
// value of the last Put of a key.
func (c keeperGet) KeepsTheLastPut() conform.ID[Keeper] {
	return c.Check("keeps-the-last-put")
}

// KeeperSuite returns the checks of Keeper, and the calls and the budgets
// of its callables.
func KeeperSuite() conform.Suite[Keeper] {
	return conform.Suite[Keeper]{
		Declaration: "Keeper",
		Checks: []conform.Check[Keeper]{
			{
				ID:        KeeperChecks.Put.LeavesTheWholeValue(),
				Statement: "a crash during Put leaves the whole value or no value",
				Apply:     LeavesTheWholeValue,
				Defect:    tearsWrites,
				Doors:     conform.Doors[Keeper]{Reference: reference, Second: second, Crash: tear},
			},
			{
				ID:        KeeperChecks.Get.KeepsEachPutAcrossARestart(),
				Statement: "Get returns each value that Put stored before a restart",
				Apply:     KeepsEachPutAcrossARestart,
				Defect:    forgetsOnRestart,
				Doors:     conform.Doors[Keeper]{Restart: forget},
			},
			{
				ID:        KeeperChecks.Get.KeepsTheLastPut(),
				Statement: "Get returns the value of the last Put of a key",
				Property:  KeepsTheLastPut,
				Defect:    keepsTheFirstPut,
				Options:   putShares,
			},
		},
		Calls: []conform.Call[Keeper]{
			{
				Calls:    keeperPutCalls,
				Over:     newKeeperOver,
				Callable: KeeperChecks.Put.Callable,
			},
			{
				Calls:    keeperGetCalls,
				Over:     newKeeperOver,
				Callable: KeeperChecks.Get.Callable,
			},
		},
	}
}

// RunKeeperSuite applies the checks of Keeper, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunKeeperSuite[T Keeper](t *testing.T, newSubject func() T, opts ...KeeperOption) {
	t.Helper()
	RunKeeperSuiteFromHandle(t, func(*KeeperHandle) T { return newSubject() }, opts...)
}

// RunKeeperSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunKeeperSuiteFromHandle[T Keeper](
	t *testing.T,
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	t.Helper()
	conform.Run(t, KeeperSuite(), keeperSubjects(newSubject), opts...)
}

// FuzzKeeperSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzKeeperSuite[T Keeper](
	f *testing.F,
	callable conform.Callable[Keeper],
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	f.Helper()
	conform.Fuzz(f, KeeperSuite(), callable, keeperSubjects(newSubject), opts...)
}

// RunKeeperBudgets applies the allocation ceilings of Keeper and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunKeeperBudgets[T Keeper](
	t *testing.T,
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	t.Helper()
	conform.RunBudgets(t, KeeperSuite(), keeperSubjects(newSubject), opts...)
}

// BenchKeeperSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchKeeperSuite[T Keeper](
	b conform.B,
	callable conform.Callable[Keeper],
	newSubject func(h *KeeperHandle) T,
	opts ...KeeperOption,
) {
	b.Helper()
	conform.Bench(b, KeeperSuite(), callable, keeperSubjects(newSubject), opts...)
}

// keeperSubjects returns the factory of the subjects of Keeper that calls
// newSubject with the handle of each check.
func keeperSubjects[T Keeper](newSubject func(h *KeeperHandle) T) func(h *KeeperHandle) Keeper {
	return func(h *KeeperHandle) Keeper { return newSubject(h) }
}

// keeperPutCalls generates the calls of Keeper.Put with the inputs that
// each case draws.
var keeperPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Keeper) error {
	key := c.Draw(prop.Of[string](), "key")
	value := c.Draw(prop.Of[string](), "value")
	return func(ctx context.Context, subject Keeper) error {
		return subject.Put(ctx, key, value)
	}
})

// keeperGetCalls generates the calls of Keeper.Get with the inputs that
// each case draws.
var keeperGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Keeper) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Keeper) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// keeperOverEmbedded is Keeper under the name of the field that keeperOver
// embeds.
type keeperOverEmbedded = Keeper

// keeperOver is a subject of Keeper whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type keeperOver struct {
	keeperOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o keeperOver) Put(context.Context, string, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o keeperOver) Get(context.Context, string) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newKeeperOver returns the planted defect of a ceiling of Keeper, whose
// callables call spend.
func newKeeperOver(_ *KeeperHandle, spend func()) Keeper {
	return keeperOver{spend: spend}
}
