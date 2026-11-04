package isolated

import (
	"context"
	"database/sql/driver"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Lists writes.
type (
	ListsHandle = conform.Handle[Lists]
	ListsCheck  = conform.Check[Lists]
	ListsOption = conform.Option[Lists]
)

// ListsChecks contains the name of each callable of Lists, with the
// identities of its checks.
var ListsChecks = struct {
	Begin    listsBegin
	Commit   listsCommit
	Rollback listsRollback
	Append   listsAppend
	Get      listsGet
}{
	Begin:    listsBegin{"Lists.Begin"},
	Commit:   listsCommit{"Lists.Commit"},
	Rollback: listsRollback{"Lists.Rollback"},
	Append:   listsAppend{"Lists.Append"},
	Get:      listsGet{"Lists.Get"},
}

// listsBegin is the name of Lists.Begin, with the identities of its
// checks.
type listsBegin struct {
	conform.Callable[Lists]
}

// Smoke returns the identity of the check that Begin returns without a
// panic for derived inputs.
func (c listsBegin) Smoke() conform.ID[Lists] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Begin returns an error for
// a cancelled context.
func (c listsBegin) Cancel() conform.ID[Lists] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Begin returns the
// zero value beside an error.
func (c listsBegin) ZeroOnError() conform.ID[Lists] {
	return c.Check("zero-on-error")
}

// listsCommit is the name of Lists.Commit, with the identities of its
// checks.
type listsCommit struct {
	conform.Callable[Lists]
}

// Smoke returns the identity of the check that Commit returns without a
// panic for derived inputs.
func (c listsCommit) Smoke() conform.ID[Lists] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Commit returns an error
// for a cancelled context.
func (c listsCommit) Cancel() conform.ID[Lists] {
	return c.Check("cancel")
}

// listsRollback is the name of Lists.Rollback, with the identities of its
// checks.
type listsRollback struct {
	conform.Callable[Lists]
}

// Smoke returns the identity of the check that Rollback returns without a
// panic for derived inputs.
func (c listsRollback) Smoke() conform.ID[Lists] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Rollback returns an error
// for a cancelled context.
func (c listsRollback) Cancel() conform.ID[Lists] {
	return c.Check("cancel")
}

// listsAppend is the name of Lists.Append, with the identities of its
// checks.
type listsAppend struct {
	conform.Callable[Lists]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c listsAppend) Smoke() conform.ID[Lists] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c listsAppend) Cancel() conform.ID[Lists] {
	return c.Check("cancel")
}

// Serializable returns the identity of the check that every concurrent
// history of the transactions of Lists has a serial order in which each
// read returns the same result.
func (c listsAppend) Serializable() conform.ID[Lists] {
	return "Lists/serializable"
}

// listsGet is the name of Lists.Get, with the identities of its checks.
type listsGet struct {
	conform.Callable[Lists]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c listsGet) Smoke() conform.ID[Lists] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c listsGet) Cancel() conform.ID[Lists] {
	return c.Check("cancel")
}

// ListsSuite returns the checks of Lists, and the calls and the budgets of
// its callables.
func ListsSuite() conform.Suite[Lists] {
	return conform.Suite[Lists]{
		Declaration: "Lists",
		Package:     "example.com/shop/isolated",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Lists]{
			{
				ID:          ListsChecks.Begin.Smoke(),
				Statement:   "Begin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ListsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx, string("label-a"))
						return err
					})
				},
				Defect: listsBeginPanics,
			},
			{
				ID:          ListsChecks.Begin.Smoke().Generated(),
				Statement:   "Begin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnLabel := c.Draw(prop.Of[string](), "label")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx, drawnLabel)
						return err
					})
				},
				Defect: listsBeginPanics,
			},
			{
				ID:          ListsChecks.Begin.Cancel(),
				Statement:   "Begin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ListsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx, string("label-a"))
						return err
					})
				},
				Defect: listsBeginIgnoresCancel,
			},
			{
				ID:          ListsChecks.Begin.Cancel().Generated(),
				Statement:   "Begin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnLabel := c.Draw(prop.Of[string](), "label")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx, drawnLabel)
						return err
					})
				},
				Defect: listsBeginIgnoresCancel,
			},
			{
				ID:          ListsChecks.Begin.ZeroOnError(),
				Statement:   "Begin returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ListsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Tx, error) {
						return h.Subject().Begin(ctx, string("label-a"))
					}, Tx{})
				},
				Defect: listsBeginValueBeside,
			},
			{
				ID:          ListsChecks.Begin.ZeroOnError().Generated(),
				Statement:   "Begin returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnLabel := c.Draw(prop.Of[string](), "label")
					probe.ZeroBesideError(h, func(ctx context.Context) (Tx, error) {
						return h.Subject().Begin(ctx, drawnLabel)
					}, Tx{})
				},
				Defect: listsBeginValueBeside,
			},
			{
				ID:          ListsChecks.Commit.Smoke(),
				Statement:   "Commit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ListsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx, Tx{ID: 42})
					})
				},
				Defect: listsCommitPanics,
			},
			{
				ID:          ListsChecks.Commit.Smoke().Generated(),
				Statement:   "Commit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx, drawnTx)
					})
				},
				Defect: listsCommitPanics,
			},
			{
				ID:          ListsChecks.Commit.Cancel(),
				Statement:   "Commit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ListsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx, Tx{ID: 42})
					})
				},
				Defect: listsCommitIgnoresCancel,
			},
			{
				ID:          ListsChecks.Commit.Cancel().Generated(),
				Statement:   "Commit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx, drawnTx)
					})
				},
				Defect: listsCommitIgnoresCancel,
			},
			{
				ID:          ListsChecks.Rollback.Smoke(),
				Statement:   "Rollback returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ListsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx, Tx{ID: 42})
					})
				},
				Defect: listsRollbackPanics,
			},
			{
				ID:          ListsChecks.Rollback.Smoke().Generated(),
				Statement:   "Rollback returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx, drawnTx)
					})
				},
				Defect: listsRollbackPanics,
			},
			{
				ID:          ListsChecks.Rollback.Cancel(),
				Statement:   "Rollback returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ListsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx, Tx{ID: 42})
					})
				},
				Defect: listsRollbackIgnoresCancel,
			},
			{
				ID:          ListsChecks.Rollback.Cancel().Generated(),
				Statement:   "Rollback returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx, drawnTx)
					})
				},
				Defect: listsRollbackIgnoresCancel,
			},
			{
				ID:          ListsChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ListsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Tx{ID: 42}, string("key-a"), int(42))
					})
				},
				Defect: listsAppendPanics,
			},
			{
				ID:          ListsChecks.Append.Smoke().Generated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnV := c.Draw(prop.Of[int](), "v")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnTx, drawnKey, drawnV)
					})
				},
				Defect: listsAppendPanics,
			},
			{
				ID:          ListsChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ListsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Tx{ID: 42}, string("key-a"), int(42))
					})
				},
				Defect: listsAppendIgnoresCancel,
			},
			{
				ID:          ListsChecks.Append.Cancel().Generated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnV := c.Draw(prop.Of[int](), "v")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnTx, drawnKey, drawnV)
					})
				},
				Defect: listsAppendIgnoresCancel,
			},
			{
				ID:          ListsChecks.Append.Serializable(),
				Statement:   "every concurrent history of the transactions of Lists has a serial order in which each read returns the same result",
				Claim:       spec.Serializable,
				Machineries: []spec.MachineryID{spec.ConcurrentCallers, spec.GeneratedSequences},
				Property: func(h *ListsHandle, c *prop.Case) {
					probe.SerializableHistory(h, c, listsActions, 2)
				},
				Defect: listsReadsUncommitted,
			},
			{
				ID:          ListsChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ListsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, Tx{ID: 42}, string("key-a"))
						return err
					})
				},
				Defect: listsGetPanics,
			},
			{
				ID:          ListsChecks.Get.Smoke().Generated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnTx, drawnKey)
						return err
					})
				},
				Defect: listsGetPanics,
			},
			{
				ID:          ListsChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ListsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, Tx{ID: 42}, string("key-a"))
						return err
					})
				},
				Defect: listsGetIgnoresCancel,
			},
			{
				ID:          ListsChecks.Get.Cancel().Generated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ListsHandle, c *prop.Case) {
					drawnTx := c.Draw(prop.Of[Tx](), "tx")
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnTx, drawnKey)
						return err
					})
				},
				Defect: listsGetIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Tx,
				Callable: "Lists.Begin",
				Gate:     spec.PartnerArguments,
				Why:      "a parameter of the read of Begin selects neither a value that commit takes nor a field of the written value",
				Remedy:   "align the name of the parameter of the read with a parameter or a field of the write",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Lists.Commit",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Lists.Rollback",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:multi-arg-writer",
				Callable:       "Lists.Append",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Lists]{
			{
				Calls:    listsBeginCalls,
				Over:     newListsOver,
				Callable: ListsChecks.Begin.Callable,
			},
			{
				Calls:    listsCommitCalls,
				Over:     newListsOver,
				Callable: ListsChecks.Commit.Callable,
			},
			{
				Calls:    listsRollbackCalls,
				Over:     newListsOver,
				Callable: ListsChecks.Rollback.Callable,
			},
			{
				Calls:    listsAppendCalls,
				Over:     newListsOver,
				Callable: ListsChecks.Append.Callable,
			},
			{
				Calls:    listsGetCalls,
				Over:     newListsOver,
				Callable: ListsChecks.Get.Callable,
			},
		},
	}
}

// RunListsSuite applies the checks of Lists, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunListsSuite[T Lists](t *testing.T, newSubject func() T, opts ...ListsOption) {
	t.Helper()
	RunListsSuiteFromHandle(t, func(*ListsHandle) T { return newSubject() }, opts...)
}

// RunListsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunListsSuiteFromHandle[T Lists](
	t *testing.T,
	newSubject func(h *ListsHandle) T,
	opts ...ListsOption,
) {
	t.Helper()
	conform.Run(t, ListsSuite(), listsSubjects(newSubject), opts...)
}

// FuzzListsSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzListsSuite[T Lists](
	f *testing.F,
	callable conform.Callable[Lists],
	newSubject func(h *ListsHandle) T,
	opts ...ListsOption,
) {
	f.Helper()
	conform.Fuzz(f, ListsSuite(), callable, listsSubjects(newSubject), opts...)
}

// RunListsBudgets applies the allocation ceilings of Lists and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunListsBudgets[T Lists](
	t *testing.T,
	newSubject func(h *ListsHandle) T,
	opts ...ListsOption,
) {
	t.Helper()
	conform.RunBudgets(t, ListsSuite(), listsSubjects(newSubject), opts...)
}

// BenchListsSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchListsSuite[T Lists](
	b conform.B,
	callable conform.Callable[Lists],
	newSubject func(h *ListsHandle) T,
	opts ...ListsOption,
) {
	b.Helper()
	conform.Bench(b, ListsSuite(), callable, listsSubjects(newSubject), opts...)
}

// listsSubjects returns the factory of the subjects of Lists that calls
// newSubject with the handle of each check.
func listsSubjects[T Lists](newSubject func(h *ListsHandle) T) func(h *ListsHandle) Lists {
	return func(h *ListsHandle) Lists { return newSubject(h) }
}

// listsBeginCalls generates the calls of Lists.Begin with the inputs that
// each case draws.
var listsBeginCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lists) error {
	label := c.Draw(prop.Of[string](), "label")
	return func(ctx context.Context, subject Lists) error {
		_, err := subject.Begin(ctx, label)
		return err
	}
})

// listsCommitCalls generates the calls of Lists.Commit with the inputs
// that each case draws.
var listsCommitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lists) error {
	tx := c.Draw(prop.Of[Tx](), "tx")
	return func(ctx context.Context, subject Lists) error {
		return subject.Commit(ctx, tx)
	}
})

// listsRollbackCalls generates the calls of Lists.Rollback with the inputs
// that each case draws.
var listsRollbackCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lists) error {
	tx := c.Draw(prop.Of[Tx](), "tx")
	return func(ctx context.Context, subject Lists) error {
		return subject.Rollback(ctx, tx)
	}
})

// listsAppendCalls generates the calls of Lists.Append with the inputs
// that each case draws.
var listsAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lists) error {
	tx := c.Draw(prop.Of[Tx](), "tx")
	key := c.Draw(prop.Of[string](), "key")
	v := c.Draw(prop.Of[int](), "v")
	return func(ctx context.Context, subject Lists) error {
		return subject.Append(ctx, tx, key, v)
	}
})

// listsGetCalls generates the calls of Lists.Get with the inputs that each
// case draws.
var listsGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Lists) error {
	tx := c.Draw(prop.Of[Tx](), "tx")
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Lists) error {
		_, err := subject.Get(ctx, tx, key)
		return err
	}
})

// listsOverEmbedded is Lists under the name of the field that listsOver
// embeds.
type listsOverEmbedded = Lists

// listsOver is a subject of Lists whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type listsOver struct {
	listsOverEmbedded
	spend func()
}

// Begin calls spend and returns its zero results.
func (o listsOver) Begin(context.Context, string) (r0 Tx, r1 error) {
	o.spend()
	return r0, r1
}

// Commit calls spend and returns its zero results.
func (o listsOver) Commit(context.Context, Tx) (r0 error) {
	o.spend()
	return r0
}

// Rollback calls spend and returns its zero results.
func (o listsOver) Rollback(context.Context, Tx) (r0 error) {
	o.spend()
	return r0
}

// Append calls spend and returns its zero results.
func (o listsOver) Append(context.Context, Tx, string, int) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o listsOver) Get(context.Context, Tx, string) (r0 []int, r1 error) {
	o.spend()
	return r0, r1
}

// newListsOver returns the planted defect of a ceiling of Lists, whose
// callables call spend.
func newListsOver(_ *ListsHandle, spend func()) Lists {
	return listsOver{spend: spend}
}

// listsBeginPanicsDefect is the planted defect of a Lists whose Begin
// panics.
type listsBeginPanicsDefect struct {
	listsOverEmbedded
}

// Begin is the override of the planted defect.
func (*listsBeginPanicsDefect) Begin(p0 context.Context, p1 string) (r0 Tx, r1 error) {
	panic("the planted defect of a Lists whose Begin panics")
}

// listsBeginPanics returns the planted defect of a Lists whose Begin
// panics over a fresh subject.
func listsBeginPanics(h *ListsHandle) Lists {
	return &listsBeginPanicsDefect{listsOverEmbedded: h.Fresh()}
}

// listsBeginIgnoresCancelDefect is the planted defect of a Lists whose
// Begin ignores the context.
type listsBeginIgnoresCancelDefect struct {
	listsOverEmbedded
}

// Begin is the override of the planted defect.
func (d *listsBeginIgnoresCancelDefect) Begin(p0 context.Context, p1 string) (r0 Tx, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.listsOverEmbedded.Begin(p0, p1)
}

// listsBeginIgnoresCancel returns the planted defect of a Lists whose
// Begin ignores the context over a fresh subject.
func listsBeginIgnoresCancel(h *ListsHandle) Lists {
	return &listsBeginIgnoresCancelDefect{listsOverEmbedded: h.Fresh()}
}

// listsBeginValueBesideDefect is the planted defect of a Lists whose Begin
// returns a value beside an error.
type listsBeginValueBesideDefect struct {
	listsOverEmbedded
}

// Begin is the override of the planted defect.
func (d *listsBeginValueBesideDefect) Begin(p0 context.Context, p1 string) (r0 Tx, r1 error) {
	r0, r1 = d.listsOverEmbedded.Begin(p0, p1)
	if r1 != nil {
		r0 = Tx{ID: 7}
	}
	return r0, r1
}

// listsBeginValueBeside returns the planted defect of a Lists whose Begin
// returns a value beside an error over a fresh subject.
func listsBeginValueBeside(h *ListsHandle) Lists {
	return &listsBeginValueBesideDefect{listsOverEmbedded: h.Fresh()}
}

// listsCommitPanicsDefect is the planted defect of a Lists whose Commit
// panics.
type listsCommitPanicsDefect struct {
	listsOverEmbedded
}

// Commit is the override of the planted defect.
func (*listsCommitPanicsDefect) Commit(p0 context.Context, p1 Tx) (r0 error) {
	panic("the planted defect of a Lists whose Commit panics")
}

// listsCommitPanics returns the planted defect of a Lists whose Commit
// panics over a fresh subject.
func listsCommitPanics(h *ListsHandle) Lists {
	return &listsCommitPanicsDefect{listsOverEmbedded: h.Fresh()}
}

// listsCommitIgnoresCancelDefect is the planted defect of a Lists whose
// Commit ignores the context.
type listsCommitIgnoresCancelDefect struct {
	listsOverEmbedded
}

// Commit is the override of the planted defect.
func (d *listsCommitIgnoresCancelDefect) Commit(p0 context.Context, p1 Tx) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.listsOverEmbedded.Commit(p0, p1)
}

// listsCommitIgnoresCancel returns the planted defect of a Lists whose
// Commit ignores the context over a fresh subject.
func listsCommitIgnoresCancel(h *ListsHandle) Lists {
	return &listsCommitIgnoresCancelDefect{listsOverEmbedded: h.Fresh()}
}

// listsRollbackPanicsDefect is the planted defect of a Lists whose
// Rollback panics.
type listsRollbackPanicsDefect struct {
	listsOverEmbedded
}

// Rollback is the override of the planted defect.
func (*listsRollbackPanicsDefect) Rollback(p0 context.Context, p1 Tx) (r0 error) {
	panic("the planted defect of a Lists whose Rollback panics")
}

// listsRollbackPanics returns the planted defect of a Lists whose Rollback
// panics over a fresh subject.
func listsRollbackPanics(h *ListsHandle) Lists {
	return &listsRollbackPanicsDefect{listsOverEmbedded: h.Fresh()}
}

// listsRollbackIgnoresCancelDefect is the planted defect of a Lists whose
// Rollback ignores the context.
type listsRollbackIgnoresCancelDefect struct {
	listsOverEmbedded
}

// Rollback is the override of the planted defect.
func (d *listsRollbackIgnoresCancelDefect) Rollback(p0 context.Context, p1 Tx) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.listsOverEmbedded.Rollback(p0, p1)
}

// listsRollbackIgnoresCancel returns the planted defect of a Lists whose
// Rollback ignores the context over a fresh subject.
func listsRollbackIgnoresCancel(h *ListsHandle) Lists {
	return &listsRollbackIgnoresCancelDefect{listsOverEmbedded: h.Fresh()}
}

// listsAppendPanicsDefect is the planted defect of a Lists whose Append
// panics.
type listsAppendPanicsDefect struct {
	listsOverEmbedded
}

// Append is the override of the planted defect.
func (*listsAppendPanicsDefect) Append(p0 context.Context, p1 Tx, p2 string, p3 int) (r0 error) {
	panic("the planted defect of a Lists whose Append panics")
}

// listsAppendPanics returns the planted defect of a Lists whose Append
// panics over a fresh subject.
func listsAppendPanics(h *ListsHandle) Lists {
	return &listsAppendPanicsDefect{listsOverEmbedded: h.Fresh()}
}

// listsAppendIgnoresCancelDefect is the planted defect of a Lists whose
// Append ignores the context.
type listsAppendIgnoresCancelDefect struct {
	listsOverEmbedded
}

// Append is the override of the planted defect.
func (d *listsAppendIgnoresCancelDefect) Append(
	p0 context.Context,
	p1 Tx,
	p2 string,
	p3 int,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.listsOverEmbedded.Append(p0, p1, p2, p3)
}

// listsAppendIgnoresCancel returns the planted defect of a Lists whose
// Append ignores the context over a fresh subject.
func listsAppendIgnoresCancel(h *ListsHandle) Lists {
	return &listsAppendIgnoresCancelDefect{listsOverEmbedded: h.Fresh()}
}

// listsReadsUncommittedDefect is the planted defect of a Lists whose read
// returns the write of a transaction that has not committed.
type listsReadsUncommittedDefect struct {
	listsOverEmbedded
	h    *ListsHandle
	mu   sync.Mutex
	open map[any][2]any
}

// Append is the override of the planted defect.
func (d *listsReadsUncommittedDefect) Append(
	p0 context.Context,
	p1 Tx,
	p2 string,
	p3 int,
) (r0 error) {
	r0 = d.listsOverEmbedded.Append(p0, p1, p2, p3)
	if r0 == nil {
		d.mu.Lock()
		d.open[p2] = [2]any{p1, p3}
		d.mu.Unlock()
	}
	d.h.Yield()
	return r0
}

// Get is the override of the planted defect.
func (d *listsReadsUncommittedDefect) Get(
	p0 context.Context,
	p1 Tx,
	p2 string,
) (r0 []int, r1 error) {
	r0, r1 = d.listsOverEmbedded.Get(p0, p1, p2)
	d.mu.Lock()
	kept, found := d.open[p2]
	d.mu.Unlock()
	if found && kept[0] != any(p1) && r1 == nil {
		v, _ := kept[1].(int)
		r0 = append(r0[:len(r0):len(r0)], v)
	}
	d.h.Yield()
	return r0, r1
}

// Commit is the override of the planted defect.
func (d *listsReadsUncommittedDefect) Commit(p0 context.Context, p1 Tx) (r0 error) {
	r0 = d.listsOverEmbedded.Commit(p0, p1)
	d.mu.Lock()
	for key, kept := range d.open {
		if kept[0] == any(p1) {
			delete(d.open, key)
		}
	}
	d.mu.Unlock()
	d.h.Yield()
	return r0
}

// Rollback is the override of the planted defect.
func (d *listsReadsUncommittedDefect) Rollback(p0 context.Context, p1 Tx) (r0 error) {
	r0 = d.listsOverEmbedded.Rollback(p0, p1)
	d.mu.Lock()
	for key, kept := range d.open {
		if kept[0] == any(p1) {
			delete(d.open, key)
		}
	}
	d.mu.Unlock()
	d.h.Yield()
	return r0
}

// listsReadsUncommitted returns the planted defect of a Lists whose read
// returns the write of a transaction that has not committed over a fresh
// subject.
func listsReadsUncommitted(h *ListsHandle) Lists {
	return &listsReadsUncommittedDefect{listsOverEmbedded: h.Fresh(), h: h, open: map[any][2]any{}}
}

// listsGetPanicsDefect is the planted defect of a Lists whose Get panics.
type listsGetPanicsDefect struct {
	listsOverEmbedded
}

// Get is the override of the planted defect.
func (*listsGetPanicsDefect) Get(p0 context.Context, p1 Tx, p2 string) (r0 []int, r1 error) {
	panic("the planted defect of a Lists whose Get panics")
}

// listsGetPanics returns the planted defect of a Lists whose Get panics
// over a fresh subject.
func listsGetPanics(h *ListsHandle) Lists {
	return &listsGetPanicsDefect{listsOverEmbedded: h.Fresh()}
}

// listsGetIgnoresCancelDefect is the planted defect of a Lists whose Get
// ignores the context.
type listsGetIgnoresCancelDefect struct {
	listsOverEmbedded
}

// Get is the override of the planted defect.
func (d *listsGetIgnoresCancelDefect) Get(
	p0 context.Context,
	p1 Tx,
	p2 string,
) (r0 []int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.listsOverEmbedded.Get(p0, p1, p2)
}

// listsGetIgnoresCancel returns the planted defect of a Lists whose Get
// ignores the context over a fresh subject.
func listsGetIgnoresCancel(h *ListsHandle) Lists {
	return &listsGetIgnoresCancelDefect{listsOverEmbedded: h.Fresh()}
}

// listsActions returns the actions of the generated sequences of Lists
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func listsActions(c *prop.Case, subject func() Lists) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	var counter atomic.Int64
	return []stateful.Action[struct{}]{
		{
			Name: "Begin",
			Input: func(c *prop.Case, _ struct{}) any {
				label := c.Draw(prop.SampledFrom(keysString...), "begin.label")
				return []any{label}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Begin", args, args[0])
				values, err := listsBeginCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Commit",
			Input: func(c *prop.Case, _ struct{}) any {
				tx := c.Draw(prop.Of[Tx](), "commit.tx")
				return []any{tx}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Commit", args)
				values, err := listsCommitCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Rollback",
			Input: func(c *prop.Case, _ struct{}) any {
				tx := c.Draw(prop.Of[Tx](), "rollback.tx")
				return []any{tx}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Rollback", args)
				values, err := listsRollbackCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Append",
			Input: func(c *prop.Case, _ struct{}) any {
				tx := c.Draw(prop.Of[Tx](), "append.tx")
				key := c.Draw(prop.Of[string](), "append.key")
				v := c.Draw(prop.Of[int](), "append.v")
				return []any{tx, key, v}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Append", args)
				values, err := listsAppendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				tx := c.Draw(prop.Of[Tx](), "get.tx")
				key := c.Draw(prop.Of[string](), "get.key")
				return []any{tx, key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args)
				values, err := listsGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name:    "txn",
			Enabled: func(struct{}) bool { return false },
			Input: func(c *prop.Case, _ struct{}) any {
				ops := make([]any, c.Draw(prop.Integer(1, 4), "txn.ops"))
				for i := range ops {
					key := c.Draw(prop.SampledFrom(keysString...), "txn.key")
					if c.Draw(prop.Boolean(), "txn.read") {
						ops[i] = []any{"read", key, nil}
						continue
					}
					ops[i] = []any{"append", key, int(counter.Add(1))}
				}
				return []any{ops, c.Draw(prop.Integer(0, 9), "txn.commit") > 0}
			},
			Run: func(c *prop.Case, client int, input any) {
				in, _ := input.([]any)
				ops, _ := in[0].([]any)
				commit, _ := in[1].(bool)
				call := c.History().Invoke(client, "txn", ops)
				ctx, s := c.Context(), subject()
				end := func(err error) {
					if ctx.Err() != nil && errors.Is(err, ctx.Err()) {
						call.Unknown(err)
						return
					}
					call.Fail(err)
				}
				begun, err := listsBeginCall(ctx, s, []any{string("label-a")})
				if err != nil {
					end(err)
					return
				}
				tx := begun[0]
				rollback := func() { _, _ = listsRollbackCall(ctx, s, []any{tx}) }
				out := make([]any, len(ops))
				for i, op := range ops {
					m, _ := op.([]any)
					key, v := m[1], m[2]
					if m[0] == "read" {
						read, err := listsGetCall(ctx, s, []any{tx, key})
						if err != nil {
							rollback()
							end(err)
							return
						}
						out[i] = []any{"read", key, read[0]}
						continue
					}
					if _, err := listsAppendCall(ctx, s, []any{tx, key, v}); err != nil {
						rollback()
						end(err)
						return
					}
					out[i] = m
				}
				if !commit {
					rollback()
					call.Fail(conform.ErrGenerated)
					return
				}
				if _, err := listsCommitCall(ctx, s, []any{tx}); err != nil {
					end(err)
					return
				}
				call.OK(out)
			},
		},
	}
}

// listsBeginCall calls Begin on subject with the arguments args, and
// returns the values of the call beside its error.
func listsBeginCall(ctx context.Context, subject Lists, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Begin(ctx, a0)
	return []any{r0}, err
}

// listsCommitCall calls Commit on subject with the arguments args, and
// returns the values of the call beside its error.
func listsCommitCall(ctx context.Context, subject Lists, args []any) ([]any, error) {
	a0, _ := args[0].(Tx)
	return nil, subject.Commit(ctx, a0)
}

// listsRollbackCall calls Rollback on subject with the arguments args, and
// returns the values of the call beside its error.
func listsRollbackCall(ctx context.Context, subject Lists, args []any) ([]any, error) {
	a0, _ := args[0].(Tx)
	return nil, subject.Rollback(ctx, a0)
}

// listsAppendCall calls Append on subject with the arguments args, and
// returns the values of the call beside its error.
func listsAppendCall(ctx context.Context, subject Lists, args []any) ([]any, error) {
	a0, _ := args[0].(Tx)
	a1, _ := args[1].(string)
	a2, _ := args[2].(int)
	return nil, subject.Append(ctx, a0, a1, a2)
}

// listsGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func listsGetCall(ctx context.Context, subject Lists, args []any) ([]any, error) {
	a0, _ := args[0].(Tx)
	a1, _ := args[1].(string)
	r0, err := subject.Get(ctx, a0, a1)
	return []any{r0}, err
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Registers writes.
type (
	RegistersHandle = conform.Handle[Registers]
	RegistersCheck  = conform.Check[Registers]
	RegistersOption = conform.Option[Registers]
)

// RegistersChecks contains the name of each callable of Registers, with
// the identities of its checks.
var RegistersChecks = struct {
	Begin  registersBegin
	Latest registersLatest
}{
	Begin:  registersBegin{"Registers.Begin"},
	Latest: registersLatest{"Registers.Latest"},
}

// registersBegin is the name of Registers.Begin, with the identities of
// its checks.
type registersBegin struct {
	conform.Callable[Registers]
}

// Smoke returns the identity of the check that Begin returns without a
// panic for derived inputs.
func (c registersBegin) Smoke() conform.ID[Registers] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Begin returns an error for
// a cancelled context.
func (c registersBegin) Cancel() conform.ID[Registers] {
	return c.Check("cancel")
}

// SnapshotIsolation returns the identity of the check that every
// concurrent history of the transactions of Registers has no anomaly that
// snapshot isolation excludes.
func (c registersBegin) SnapshotIsolation() conform.ID[Registers] {
	return "Registers/snapshot-isolation"
}

// registersLatest is the name of Registers.Latest, with the identities of
// its checks.
type registersLatest struct {
	conform.Callable[Registers]
}

// Smoke returns the identity of the check that Latest returns without a
// panic for derived inputs.
func (c registersLatest) Smoke() conform.ID[Registers] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Latest returns an error
// for a cancelled context.
func (c registersLatest) Cancel() conform.ID[Registers] {
	return c.Check("cancel")
}

// RegistersSuite returns the checks of Registers, and the calls and the
// budgets of its callables.
func RegistersSuite() conform.Suite[Registers] {
	return conform.Suite[Registers]{
		Declaration: "Registers",
		Package:     "example.com/shop/isolated",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Registers]{
			{
				ID:          RegistersChecks.Begin.Smoke(),
				Statement:   "Begin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistersHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx)
						return err
					})
				},
				Defect: registersBeginPanics,
			},
			{
				ID:          RegistersChecks.Begin.Cancel(),
				Statement:   "Begin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistersHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx)
						return err
					})
				},
				Defect: registersBeginIgnoresCancel,
			},
			{
				ID:          RegistersChecks.Begin.SnapshotIsolation(),
				Statement:   "every concurrent history of the transactions of Registers has no anomaly that snapshot isolation excludes",
				Claim:       spec.SnapshotIsolation,
				Machineries: []spec.MachineryID{spec.ConcurrentCallers, spec.GeneratedSequences},
				Property: func(h *RegistersHandle, c *prop.Case) {
					probe.SnapshotHistory(h, c, registersActions, 2)
				},
				Defect: registersWritesWithoutASnapshot,
			},
			{
				ID:          RegistersChecks.Latest.Smoke(),
				Statement:   "Latest returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegistersHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Latest(ctx)
						return err
					})
				},
				Defect: registersLatestPanics,
			},
			{
				ID:          RegistersChecks.Latest.Cancel(),
				Statement:   "Latest returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegistersHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Latest(ctx)
						return err
					})
				},
				Defect: registersLatestIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ZeroOnError,
				Callable: "Registers.Begin",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Begin",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.Tx,
				Callable: "Registers.Begin",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Begin",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Registers.Latest",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Latest",
				Remedy:   "author a sample of the result type with the directive sample",
			},
		},
		Calls: []conform.Call[Registers]{
			{
				Calls:    registersBeginCalls,
				Over:     newRegistersOver,
				Callable: RegistersChecks.Begin.Callable,
			},
			{
				Calls:    registersLatestCalls,
				Over:     newRegistersOver,
				Callable: RegistersChecks.Latest.Callable,
			},
		},
	}
}

// RunRegistersSuite applies the checks of Registers, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunRegistersSuite[T Registers](t *testing.T, newSubject func() T, opts ...RegistersOption) {
	t.Helper()
	RunRegistersSuiteFromHandle(t, func(*RegistersHandle) T { return newSubject() }, opts...)
}

// RunRegistersSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRegistersSuiteFromHandle[T Registers](
	t *testing.T,
	newSubject func(h *RegistersHandle) T,
	opts ...RegistersOption,
) {
	t.Helper()
	conform.Run(t, RegistersSuite(), registersSubjects(newSubject), opts...)
}

// FuzzRegistersSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRegistersSuite[T Registers](
	f *testing.F,
	callable conform.Callable[Registers],
	newSubject func(h *RegistersHandle) T,
	opts ...RegistersOption,
) {
	f.Helper()
	conform.Fuzz(f, RegistersSuite(), callable, registersSubjects(newSubject), opts...)
}

// RunRegistersBudgets applies the allocation ceilings of Registers and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRegistersBudgets[T Registers](
	t *testing.T,
	newSubject func(h *RegistersHandle) T,
	opts ...RegistersOption,
) {
	t.Helper()
	conform.RunBudgets(t, RegistersSuite(), registersSubjects(newSubject), opts...)
}

// BenchRegistersSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchRegistersSuite[T Registers](
	b conform.B,
	callable conform.Callable[Registers],
	newSubject func(h *RegistersHandle) T,
	opts ...RegistersOption,
) {
	b.Helper()
	conform.Bench(b, RegistersSuite(), callable, registersSubjects(newSubject), opts...)
}

// registersSubjects returns the factory of the subjects of Registers that
// calls newSubject with the handle of each check.
func registersSubjects[T Registers](
	newSubject func(h *RegistersHandle) T,
) func(h *RegistersHandle) Registers {
	return func(h *RegistersHandle) Registers { return newSubject(h) }
}

// registersBeginCalls generates the calls of Registers.Begin with the
// inputs that each case draws.
var registersBeginCalls = prop.Composite(func(*prop.Case) func(context.Context, Registers) error {
	return func(ctx context.Context, subject Registers) error {
		_, err := subject.Begin(ctx)
		return err
	}
})

// registersLatestCalls generates the calls of Registers.Latest with the
// inputs that each case draws.
var registersLatestCalls = prop.Composite(func(*prop.Case) func(context.Context, Registers) error {
	return func(ctx context.Context, subject Registers) error {
		_, err := subject.Latest(ctx)
		return err
	}
})

// registersOverEmbedded is Registers under the name of the field that
// registersOver embeds.
type registersOverEmbedded = Registers

// registersOver is a subject of Registers whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type registersOver struct {
	registersOverEmbedded
	spend func()
}

// Begin calls spend and returns its zero results.
func (o registersOver) Begin(context.Context) (r0 Txn, r1 error) {
	o.spend()
	return r0, r1
}

// Latest calls spend and returns its zero results.
func (o registersOver) Latest(context.Context) (r0 Txn, r1 error) {
	o.spend()
	return r0, r1
}

// newRegistersOver returns the planted defect of a ceiling of Registers,
// whose callables call spend.
func newRegistersOver(_ *RegistersHandle, spend func()) Registers {
	return registersOver{spend: spend}
}

// registersBeginPanicsDefect is the planted defect of a Registers whose
// Begin panics.
type registersBeginPanicsDefect struct {
	registersOverEmbedded
}

// Begin is the override of the planted defect.
func (*registersBeginPanicsDefect) Begin(p0 context.Context) (r0 Txn, r1 error) {
	panic("the planted defect of a Registers whose Begin panics")
}

// registersBeginPanics returns the planted defect of a Registers whose
// Begin panics over a fresh subject.
func registersBeginPanics(h *RegistersHandle) Registers {
	return &registersBeginPanicsDefect{registersOverEmbedded: h.Fresh()}
}

// registersBeginIgnoresCancelDefect is the planted defect of a Registers
// whose Begin ignores the context.
type registersBeginIgnoresCancelDefect struct {
	registersOverEmbedded
}

// Begin is the override of the planted defect.
func (d *registersBeginIgnoresCancelDefect) Begin(p0 context.Context) (r0 Txn, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registersOverEmbedded.Begin(p0)
}

// registersBeginIgnoresCancel returns the planted defect of a Registers
// whose Begin ignores the context over a fresh subject.
func registersBeginIgnoresCancel(h *RegistersHandle) Registers {
	return &registersBeginIgnoresCancelDefect{registersOverEmbedded: h.Fresh()}
}

// registersWritesWithoutASnapshotDefect is the planted defect of a
// Registers whose commit applies the writes of its transaction without the
// snapshot of the transaction.
type registersWritesWithoutASnapshotDefect struct {
	registersOverEmbedded
	h      *RegistersHandle
	mu     sync.Mutex
	staged map[any][][2]any
}

// Begin is the override of the planted defect.
func (d *registersWritesWithoutASnapshotDefect) Begin(p0 context.Context) (r0 Txn, r1 error) {
	r0, r1 = d.registersOverEmbedded.Begin(p0)
	if r1 == nil {
		r0 = &registersWritesWithoutASnapshotTxn{Txn: r0, d: d}
	}
	d.h.Yield()
	return r0, r1
}

// registersWritesWithoutASnapshotTxn wraps a handle of a transaction of
// the planted defect of a Registers whose commit applies the writes of its
// transaction without the snapshot of the transaction.
type registersWritesWithoutASnapshotTxn struct {
	Txn
	d *registersWritesWithoutASnapshotDefect
}

// Set is the override of the planted defect.
func (t *registersWritesWithoutASnapshotTxn) Set(p0 context.Context, p1, p2 string) (r0 error) {
	t.d.mu.Lock()
	t.d.staged[t] = append(t.d.staged[t], [2]any{p1, p2})
	t.d.mu.Unlock()
	t.d.h.Yield()
	return r0
}

// Load is the override of the planted defect.
func (t *registersWritesWithoutASnapshotTxn) Load(
	p0 context.Context,
	p1 string,
) (r0 string, r1 error) {
	r0, r1 = t.Txn.Load(p0, p1)
	if r1 == nil {
		t.d.mu.Lock()
		for _, w := range t.d.staged[t] {
			if w[0] == any(p1) {
				v, _ := w[1].(string)
				r0 = v
			}
		}
		t.d.mu.Unlock()
	}
	t.d.h.Yield()
	return r0, r1
}

// Commit is the override of the planted defect.
func (t *registersWritesWithoutASnapshotTxn) Commit(p0 context.Context) (r0 error) {
	t.d.mu.Lock()
	staged := t.d.staged[t]
	delete(t.d.staged, t)
	t.d.mu.Unlock()
	r0 = t.Txn.Commit(p0)
	if r0 != nil || len(staged) == 0 {
		t.d.h.Yield()
		return r0
	}
	ctx := p0
	begun, err := registersBeginCall(ctx, t.d.registersOverEmbedded, []any{})
	if err != nil {
		r0 = err
		return r0
	}
	tx, _ := begun[0].(Txn)
	for _, w := range staged {
		key, v := w[0], w[1]
		if _, err := txnSetCall(ctx, tx, []any{key, v}); err != nil {
			_, _ = txnRollbackCall(ctx, tx, []any{})
			r0 = err
			return r0
		}
	}
	_, r0 = txnCommitCall(ctx, tx, []any{})
	t.d.h.Yield()
	return r0
}

// Rollback is the override of the planted defect.
func (t *registersWritesWithoutASnapshotTxn) Rollback(p0 context.Context) (r0 error) {
	t.d.mu.Lock()
	delete(t.d.staged, t)
	t.d.mu.Unlock()
	r0 = t.Txn.Rollback(p0)
	t.d.h.Yield()
	return r0
}

// registersWritesWithoutASnapshot returns the planted defect of a
// Registers whose commit applies the writes of its transaction without the
// snapshot of the transaction over a fresh subject.
func registersWritesWithoutASnapshot(h *RegistersHandle) Registers {
	return &registersWritesWithoutASnapshotDefect{
		registersOverEmbedded: h.Fresh(),
		h:                     h,
		staged:                map[any][][2]any{},
	}
}

// registersLatestPanicsDefect is the planted defect of a Registers whose
// Latest panics.
type registersLatestPanicsDefect struct {
	registersOverEmbedded
}

// Latest is the override of the planted defect.
func (*registersLatestPanicsDefect) Latest(p0 context.Context) (r0 Txn, r1 error) {
	panic("the planted defect of a Registers whose Latest panics")
}

// registersLatestPanics returns the planted defect of a Registers whose
// Latest panics over a fresh subject.
func registersLatestPanics(h *RegistersHandle) Registers {
	return &registersLatestPanicsDefect{registersOverEmbedded: h.Fresh()}
}

// registersLatestIgnoresCancelDefect is the planted defect of a Registers
// whose Latest ignores the context.
type registersLatestIgnoresCancelDefect struct {
	registersOverEmbedded
}

// Latest is the override of the planted defect.
func (d *registersLatestIgnoresCancelDefect) Latest(p0 context.Context) (r0 Txn, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registersOverEmbedded.Latest(p0)
}

// registersLatestIgnoresCancel returns the planted defect of a Registers
// whose Latest ignores the context over a fresh subject.
func registersLatestIgnoresCancel(h *RegistersHandle) Registers {
	return &registersLatestIgnoresCancelDefect{registersOverEmbedded: h.Fresh()}
}

// registersActions returns the actions of the generated sequences of
// Registers over the subject that subject returns. The actions of a case
// take their keys from the pools that it draws.
func registersActions(c *prop.Case, subject func() Registers) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	var counter atomic.Int64
	return []stateful.Action[struct{}]{
		{
			Name: "Begin",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Begin", args)
				values, err := registersBeginCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Latest",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Latest", args)
				values, err := registersLatestCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name:    "txn",
			Enabled: func(struct{}) bool { return false },
			Input: func(c *prop.Case, _ struct{}) any {
				ops := make([]any, c.Draw(prop.Integer(1, 4), "txn.ops"))
				for i := range ops {
					key := c.Draw(prop.SampledFrom(keysString...), "txn.key")
					if c.Draw(prop.Boolean(), "txn.read") {
						ops[i] = []any{"read", key, nil}
						continue
					}
					ops[i] = []any{"write", key, string(strconv.FormatInt(counter.Add(1), 10))}
				}
				return []any{ops, c.Draw(prop.Integer(0, 9), "txn.commit") > 0}
			},
			Run: func(c *prop.Case, client int, input any) {
				in, _ := input.([]any)
				ops, _ := in[0].([]any)
				commit, _ := in[1].(bool)
				call := c.History().Invoke(client, "txn", ops)
				ctx, s := c.Context(), subject()
				end := func(err error) {
					if ctx.Err() != nil && errors.Is(err, ctx.Err()) {
						call.Unknown(err)
						return
					}
					call.Fail(err)
				}
				begun, err := registersBeginCall(ctx, s, []any{})
				if err != nil {
					end(err)
					return
				}
				tx, _ := begun[0].(Txn)
				rollback := func() { _, _ = txnRollbackCall(ctx, tx, []any{}) }
				out := make([]any, len(ops))
				for i, op := range ops {
					m, _ := op.([]any)
					key, v := m[1], m[2]
					if m[0] == "read" {
						read, err := txnLoadCall(ctx, tx, []any{key})
						if err != nil {
							rollback()
							end(err)
							return
						}
						var value any
						if r := read[0]; r != any(*new(string)) {
							value = r
						}
						out[i] = []any{"read", key, value}
						continue
					}
					if _, err := txnSetCall(ctx, tx, []any{key, v}); err != nil {
						rollback()
						end(err)
						return
					}
					out[i] = m
				}
				if !commit {
					rollback()
					call.Fail(conform.ErrGenerated)
					return
				}
				if _, err := txnCommitCall(ctx, tx, []any{}); err != nil {
					end(err)
					return
				}
				call.OK(out)
			},
		},
	}
}

// registersBeginCall calls Begin on subject with the arguments args, and
// returns the values of the call beside its error.
func registersBeginCall(ctx context.Context, subject Registers, args []any) ([]any, error) {
	r0, err := subject.Begin(ctx)
	return []any{r0}, err
}

// registersLatestCall calls Latest on subject with the arguments args, and
// returns the values of the call beside its error.
func registersLatestCall(ctx context.Context, subject Registers, args []any) ([]any, error) {
	r0, err := subject.Latest(ctx)
	return []any{r0}, err
}

// txnSetCall calls Set on subject with the arguments args, and returns the
// values of the call beside its error.
func txnSetCall(ctx context.Context, subject Txn, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(string)
	return nil, subject.Set(ctx, a0, a1)
}

// txnLoadCall calls Load on subject with the arguments args, and returns
// the values of the call beside its error.
func txnLoadCall(ctx context.Context, subject Txn, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Load(ctx, a0)
	return []any{r0}, err
}

// txnCommitCall calls Commit on subject with the arguments args, and
// returns the values of the call beside its error.
func txnCommitCall(ctx context.Context, subject Txn, args []any) ([]any, error) {
	return nil, subject.Commit(ctx)
}

// txnRollbackCall calls Rollback on subject with the arguments args, and
// returns the values of the call beside its error.
func txnRollbackCall(ctx context.Context, subject Txn, args []any) ([]any, error) {
	return nil, subject.Rollback(ctx)
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Latch writes.
type (
	LatchHandle = conform.Handle[Latch]
	LatchCheck  = conform.Check[Latch]
	LatchOption = conform.Option[Latch]
)

// LatchChecks contains the name of each callable of Latch, with the
// identities of its checks.
var LatchChecks = struct {
	Begin    latchBegin
	Commit   latchCommit
	Rollback latchRollback
}{
	Begin:    latchBegin{"Latch.Begin"},
	Commit:   latchCommit{"Latch.Commit"},
	Rollback: latchRollback{"Latch.Rollback"},
}

// latchBegin is the name of Latch.Begin, with the identities of its
// checks.
type latchBegin struct {
	conform.Callable[Latch]
}

// Smoke returns the identity of the check that Begin returns without a
// panic for derived inputs.
func (c latchBegin) Smoke() conform.ID[Latch] {
	return c.Check("smoke")
}

// latchCommit is the name of Latch.Commit, with the identities of its
// checks.
type latchCommit struct {
	conform.Callable[Latch]
}

// Smoke returns the identity of the check that Commit returns without a
// panic for derived inputs.
func (c latchCommit) Smoke() conform.ID[Latch] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Commit returns an error
// for a cancelled context.
func (c latchCommit) Cancel() conform.ID[Latch] {
	return c.Check("cancel")
}

// latchRollback is the name of Latch.Rollback, with the identities of its
// checks.
type latchRollback struct {
	conform.Callable[Latch]
}

// Smoke returns the identity of the check that Rollback returns without a
// panic for derived inputs.
func (c latchRollback) Smoke() conform.ID[Latch] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Rollback returns an error
// for a cancelled context.
func (c latchRollback) Cancel() conform.ID[Latch] {
	return c.Check("cancel")
}

// LatchSuite returns the checks of Latch, and the calls and the budgets of
// its callables.
func LatchSuite() conform.Suite[Latch] {
	return conform.Suite[Latch]{
		Declaration: "Latch",
		Package:     "example.com/shop/isolated",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Latch]{
			{
				ID:          LatchChecks.Begin.Smoke(),
				Statement:   "Begin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LatchHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Begin(ctx)
						return nil
					})
				},
				Defect: latchBeginPanics,
			},
			{
				ID:          LatchChecks.Commit.Smoke(),
				Statement:   "Commit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LatchHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx)
					})
				},
				Defect: latchCommitPanics,
			},
			{
				ID:          LatchChecks.Commit.Cancel(),
				Statement:   "Commit returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LatchHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Commit(ctx)
					})
				},
				Defect: latchCommitIgnoresCancel,
			},
			{
				ID:          LatchChecks.Rollback.Smoke(),
				Statement:   "Rollback returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LatchHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx)
					})
				},
				Defect: latchRollbackPanics,
			},
			{
				ID:          LatchChecks.Rollback.Cancel(),
				Statement:   "Rollback returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LatchHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Rollback(ctx)
					})
				},
				Defect: latchRollbackIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Latch.Begin",
				Gate:     spec.Failure,
				Why:      "Begin cannot fail",
				Remedy:   "give Begin an error result",
			},
			{
				Claim:    spec.Tx,
				Callable: "Latch.Begin",
				Gate:     spec.OneValue,
				Why:      "Begin does not return exactly one value",
				Remedy:   "return one value from Begin",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Latch.Commit",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Latch.Rollback",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Latch]{
			{
				Calls:    latchBeginCalls,
				Over:     newLatchOver,
				Callable: LatchChecks.Begin.Callable,
			},
			{
				Calls:    latchCommitCalls,
				Over:     newLatchOver,
				Callable: LatchChecks.Commit.Callable,
			},
			{
				Calls:    latchRollbackCalls,
				Over:     newLatchOver,
				Callable: LatchChecks.Rollback.Callable,
			},
		},
	}
}

// RunLatchSuite applies the checks of Latch, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunLatchSuite[T Latch](t *testing.T, newSubject func() T, opts ...LatchOption) {
	t.Helper()
	RunLatchSuiteFromHandle(t, func(*LatchHandle) T { return newSubject() }, opts...)
}

// RunLatchSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLatchSuiteFromHandle[T Latch](
	t *testing.T,
	newSubject func(h *LatchHandle) T,
	opts ...LatchOption,
) {
	t.Helper()
	conform.Run(t, LatchSuite(), latchSubjects(newSubject), opts...)
}

// FuzzLatchSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLatchSuite[T Latch](
	f *testing.F,
	callable conform.Callable[Latch],
	newSubject func(h *LatchHandle) T,
	opts ...LatchOption,
) {
	f.Helper()
	conform.Fuzz(f, LatchSuite(), callable, latchSubjects(newSubject), opts...)
}

// RunLatchBudgets applies the allocation ceilings of Latch and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunLatchBudgets[T Latch](
	t *testing.T,
	newSubject func(h *LatchHandle) T,
	opts ...LatchOption,
) {
	t.Helper()
	conform.RunBudgets(t, LatchSuite(), latchSubjects(newSubject), opts...)
}

// BenchLatchSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLatchSuite[T Latch](
	b conform.B,
	callable conform.Callable[Latch],
	newSubject func(h *LatchHandle) T,
	opts ...LatchOption,
) {
	b.Helper()
	conform.Bench(b, LatchSuite(), callable, latchSubjects(newSubject), opts...)
}

// latchSubjects returns the factory of the subjects of Latch that calls
// newSubject with the handle of each check.
func latchSubjects[T Latch](newSubject func(h *LatchHandle) T) func(h *LatchHandle) Latch {
	return func(h *LatchHandle) Latch { return newSubject(h) }
}

// latchBeginCalls generates the calls of Latch.Begin with the inputs that
// each case draws.
var latchBeginCalls = prop.Composite(func(*prop.Case) func(context.Context, Latch) error {
	return func(ctx context.Context, subject Latch) error {
		subject.Begin(ctx)
		return nil
	}
})

// latchCommitCalls generates the calls of Latch.Commit with the inputs
// that each case draws.
var latchCommitCalls = prop.Composite(func(*prop.Case) func(context.Context, Latch) error {
	return func(ctx context.Context, subject Latch) error {
		return subject.Commit(ctx)
	}
})

// latchRollbackCalls generates the calls of Latch.Rollback with the inputs
// that each case draws.
var latchRollbackCalls = prop.Composite(func(*prop.Case) func(context.Context, Latch) error {
	return func(ctx context.Context, subject Latch) error {
		return subject.Rollback(ctx)
	}
})

// latchOverEmbedded is Latch under the name of the field that latchOver
// embeds.
type latchOverEmbedded = Latch

// latchOver is a subject of Latch whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type latchOver struct {
	latchOverEmbedded
	spend func()
}

// Begin calls spend and returns its zero results.
func (o latchOver) Begin(context.Context) {
	o.spend()
}

// Commit calls spend and returns its zero results.
func (o latchOver) Commit(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Rollback calls spend and returns its zero results.
func (o latchOver) Rollback(context.Context) (r0 error) {
	o.spend()
	return r0
}

// newLatchOver returns the planted defect of a ceiling of Latch, whose
// callables call spend.
func newLatchOver(_ *LatchHandle, spend func()) Latch {
	return latchOver{spend: spend}
}

// latchBeginPanicsDefect is the planted defect of a Latch whose Begin
// panics.
type latchBeginPanicsDefect struct {
	latchOverEmbedded
}

// Begin is the override of the planted defect.
func (*latchBeginPanicsDefect) Begin(p0 context.Context) {
	panic("the planted defect of a Latch whose Begin panics")
}

// latchBeginPanics returns the planted defect of a Latch whose Begin
// panics over a fresh subject.
func latchBeginPanics(h *LatchHandle) Latch {
	return &latchBeginPanicsDefect{latchOverEmbedded: h.Fresh()}
}

// latchCommitPanicsDefect is the planted defect of a Latch whose Commit
// panics.
type latchCommitPanicsDefect struct {
	latchOverEmbedded
}

// Commit is the override of the planted defect.
func (*latchCommitPanicsDefect) Commit(p0 context.Context) (r0 error) {
	panic("the planted defect of a Latch whose Commit panics")
}

// latchCommitPanics returns the planted defect of a Latch whose Commit
// panics over a fresh subject.
func latchCommitPanics(h *LatchHandle) Latch {
	return &latchCommitPanicsDefect{latchOverEmbedded: h.Fresh()}
}

// latchCommitIgnoresCancelDefect is the planted defect of a Latch whose
// Commit ignores the context.
type latchCommitIgnoresCancelDefect struct {
	latchOverEmbedded
}

// Commit is the override of the planted defect.
func (d *latchCommitIgnoresCancelDefect) Commit(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.latchOverEmbedded.Commit(p0)
}

// latchCommitIgnoresCancel returns the planted defect of a Latch whose
// Commit ignores the context over a fresh subject.
func latchCommitIgnoresCancel(h *LatchHandle) Latch {
	return &latchCommitIgnoresCancelDefect{latchOverEmbedded: h.Fresh()}
}

// latchRollbackPanicsDefect is the planted defect of a Latch whose
// Rollback panics.
type latchRollbackPanicsDefect struct {
	latchOverEmbedded
}

// Rollback is the override of the planted defect.
func (*latchRollbackPanicsDefect) Rollback(p0 context.Context) (r0 error) {
	panic("the planted defect of a Latch whose Rollback panics")
}

// latchRollbackPanics returns the planted defect of a Latch whose Rollback
// panics over a fresh subject.
func latchRollbackPanics(h *LatchHandle) Latch {
	return &latchRollbackPanicsDefect{latchOverEmbedded: h.Fresh()}
}

// latchRollbackIgnoresCancelDefect is the planted defect of a Latch whose
// Rollback ignores the context.
type latchRollbackIgnoresCancelDefect struct {
	latchOverEmbedded
}

// Rollback is the override of the planted defect.
func (d *latchRollbackIgnoresCancelDefect) Rollback(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.latchOverEmbedded.Rollback(p0)
}

// latchRollbackIgnoresCancel returns the planted defect of a Latch whose
// Rollback ignores the context over a fresh subject.
func latchRollbackIgnoresCancel(h *LatchHandle) Latch {
	return &latchRollbackIgnoresCancelDefect{latchOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Driver writes.
type (
	DriverHandle = conform.Handle[Driver]
	DriverCheck  = conform.Check[Driver]
	DriverOption = conform.Option[Driver]
)

// DriverChecks contains the name of each callable of Driver, with the
// identities of its checks.
var DriverChecks = struct {
	Begin    driverBegin
	Commit   driverCommit
	Rollback driverRollback
}{
	Begin:    driverBegin{"Driver.Begin"},
	Commit:   driverCommit{"Driver.Commit"},
	Rollback: driverRollback{"Driver.Rollback"},
}

// driverBegin is the name of Driver.Begin, with the identities of its
// checks.
type driverBegin struct {
	conform.Callable[Driver]
}

// Smoke returns the identity of the check that Begin returns without a
// panic for derived inputs.
func (c driverBegin) Smoke() conform.ID[Driver] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Begin returns an error for
// a cancelled context.
func (c driverBegin) Cancel() conform.ID[Driver] {
	return c.Check("cancel")
}

// driverCommit is the name of Driver.Commit, with the identities of its
// checks.
type driverCommit struct {
	conform.Callable[Driver]
}

// driverRollback is the name of Driver.Rollback, with the identities of
// its checks.
type driverRollback struct {
	conform.Callable[Driver]
}

// DriverSuite returns the checks of Driver, and the calls and the budgets
// of its callables.
func DriverSuite() conform.Suite[Driver] {
	return conform.Suite[Driver]{
		Declaration: "Driver",
		Package:     "example.com/shop/isolated",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Driver]{
			{
				ID:          DriverChecks.Begin.Smoke(),
				Statement:   "Begin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DriverHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx)
						return err
					})
				},
				Defect: driverBeginPanics,
			},
			{
				ID:          DriverChecks.Begin.Cancel(),
				Statement:   "Begin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DriverHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Begin(ctx)
						return err
					})
				},
				Defect: driverBeginIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ZeroOnError,
				Callable: "Driver.Begin",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Begin",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.Tx,
				Callable: "Driver.Begin",
				Gate:     spec.ValueSamples,
				Why:      "eidos derives no sample of a result of Begin",
				Remedy:   "author a sample of the result type with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Driver.Commit",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter tx of Commit",
				Remedy:   "author a sample of the parameter tx with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Driver.Rollback",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter tx of Rollback",
				Remedy:   "author a sample of the parameter tx with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Driver.Commit",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Driver.Rollback",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Driver]{
			{
				Calls:    driverBeginCalls,
				Over:     newDriverOver,
				Callable: DriverChecks.Begin.Callable,
			},
			{
				Calls:    driverCommitCalls,
				Over:     newDriverOver,
				Callable: DriverChecks.Commit.Callable,
			},
			{
				Calls:    driverRollbackCalls,
				Over:     newDriverOver,
				Callable: DriverChecks.Rollback.Callable,
			},
		},
	}
}

// RunDriverSuite applies the checks of Driver, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunDriverSuite[T Driver](t *testing.T, newSubject func() T, opts ...DriverOption) {
	t.Helper()
	RunDriverSuiteFromHandle(t, func(*DriverHandle) T { return newSubject() }, opts...)
}

// RunDriverSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDriverSuiteFromHandle[T Driver](
	t *testing.T,
	newSubject func(h *DriverHandle) T,
	opts ...DriverOption,
) {
	t.Helper()
	conform.Run(t, DriverSuite(), driverSubjects(newSubject), opts...)
}

// FuzzDriverSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzDriverSuite[T Driver](
	f *testing.F,
	callable conform.Callable[Driver],
	newSubject func(h *DriverHandle) T,
	opts ...DriverOption,
) {
	f.Helper()
	conform.Fuzz(f, DriverSuite(), callable, driverSubjects(newSubject), opts...)
}

// RunDriverBudgets applies the allocation ceilings of Driver and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunDriverBudgets[T Driver](
	t *testing.T,
	newSubject func(h *DriverHandle) T,
	opts ...DriverOption,
) {
	t.Helper()
	conform.RunBudgets(t, DriverSuite(), driverSubjects(newSubject), opts...)
}

// BenchDriverSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchDriverSuite[T Driver](
	b conform.B,
	callable conform.Callable[Driver],
	newSubject func(h *DriverHandle) T,
	opts ...DriverOption,
) {
	b.Helper()
	conform.Bench(b, DriverSuite(), callable, driverSubjects(newSubject), opts...)
}

// driverSubjects returns the factory of the subjects of Driver that calls
// newSubject with the handle of each check.
func driverSubjects[T Driver](newSubject func(h *DriverHandle) T) func(h *DriverHandle) Driver {
	return func(h *DriverHandle) Driver { return newSubject(h) }
}

// driverBeginCalls generates the calls of Driver.Begin with the inputs
// that each case draws.
var driverBeginCalls = prop.Composite(func(*prop.Case) func(context.Context, Driver) error {
	return func(ctx context.Context, subject Driver) error {
		_, err := subject.Begin(ctx)
		return err
	}
})

// driverCommitCalls generates the calls of Driver.Commit with the inputs
// that each case draws.
var driverCommitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Driver) error {
	txCommitResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"tx-commit",
	)
	var txCommitAt atomic.Int64
	txRollbackResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"tx-rollback",
	)
	var txRollbackAt atomic.Int64
	tx := &driverTxInput{
		onCommit: func() (r0 error) {
			i := int(txCommitAt.Add(1) - 1)
			if i >= len(txCommitResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return txCommitResults[i]
		},
		onRollback: func() (r0 error) {
			i := int(txRollbackAt.Add(1) - 1)
			if i >= len(txRollbackResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return txRollbackResults[i]
		},
	}
	return func(ctx context.Context, subject Driver) error {
		return subject.Commit(ctx, tx)
	}
})

// driverRollbackCalls generates the calls of Driver.Rollback with the
// inputs that each case draws.
var driverRollbackCalls = prop.Composite(func(c *prop.Case) func(context.Context, Driver) error {
	txCommitResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"tx-commit",
	)
	var txCommitAt atomic.Int64
	txRollbackResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"tx-rollback",
	)
	var txRollbackAt atomic.Int64
	tx := &driverTxInput{
		onCommit: func() (r0 error) {
			i := int(txCommitAt.Add(1) - 1)
			if i >= len(txCommitResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return txCommitResults[i]
		},
		onRollback: func() (r0 error) {
			i := int(txRollbackAt.Add(1) - 1)
			if i >= len(txRollbackResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return txRollbackResults[i]
		},
	}
	return func(ctx context.Context, subject Driver) error {
		return subject.Rollback(ctx, tx)
	}
})

// driverTxInputEmbedded is driver.Tx under the name of the field that
// driverTxInput embeds.
type driverTxInputEmbedded = driver.Tx

// driverTxInput implements driver.Tx. Its methods return the results that
// the generators of the suite of Driver draw for them, in turn.
type driverTxInput struct {
	driverTxInputEmbedded
	onCommit   func() error
	onRollback func() error
}

// Commit calls the function of Commit.
func (in driverTxInput) Commit() error {
	return in.onCommit()
}

// Rollback calls the function of Rollback.
func (in driverTxInput) Rollback() error {
	return in.onRollback()
}

// driverOverEmbedded is Driver under the name of the field that driverOver
// embeds.
type driverOverEmbedded = Driver

// driverOver is a subject of Driver whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type driverOver struct {
	driverOverEmbedded
	spend func()
}

// Begin calls spend and returns its zero results.
func (o driverOver) Begin(context.Context) (r0 driver.Tx, r1 error) {
	o.spend()
	return r0, r1
}

// Commit calls spend and returns its zero results.
func (o driverOver) Commit(context.Context, driver.Tx) (r0 error) {
	o.spend()
	return r0
}

// Rollback calls spend and returns its zero results.
func (o driverOver) Rollback(context.Context, driver.Tx) (r0 error) {
	o.spend()
	return r0
}

// newDriverOver returns the planted defect of a ceiling of Driver, whose
// callables call spend.
func newDriverOver(_ *DriverHandle, spend func()) Driver {
	return driverOver{spend: spend}
}

// driverBeginPanicsDefect is the planted defect of a Driver whose Begin
// panics.
type driverBeginPanicsDefect struct {
	driverOverEmbedded
}

// Begin is the override of the planted defect.
func (*driverBeginPanicsDefect) Begin(p0 context.Context) (r0 driver.Tx, r1 error) {
	panic("the planted defect of a Driver whose Begin panics")
}

// driverBeginPanics returns the planted defect of a Driver whose Begin
// panics over a fresh subject.
func driverBeginPanics(h *DriverHandle) Driver {
	return &driverBeginPanicsDefect{driverOverEmbedded: h.Fresh()}
}

// driverBeginIgnoresCancelDefect is the planted defect of a Driver whose
// Begin ignores the context.
type driverBeginIgnoresCancelDefect struct {
	driverOverEmbedded
}

// Begin is the override of the planted defect.
func (d *driverBeginIgnoresCancelDefect) Begin(p0 context.Context) (r0 driver.Tx, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.driverOverEmbedded.Begin(p0)
}

// driverBeginIgnoresCancel returns the planted defect of a Driver whose
// Begin ignores the context over a fresh subject.
func driverBeginIgnoresCancel(h *DriverHandle) Driver {
	return &driverBeginIgnoresCancelDefect{driverOverEmbedded: h.Fresh()}
}
