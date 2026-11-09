package buffered

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Records writes.
type (
	RecordsHandle = conform.Handle[Records]
	RecordsCheck  = conform.Check[Records]
	RecordsOption = conform.Option[Records]
)

// RecordsChecks contains the name of each callable of Records, with the
// identities of its checks.
var RecordsChecks = struct {
	Read recordsRead
	Fill recordsFill
	Peek recordsPeek
}{
	Read: recordsRead{"Records.Read"},
	Fill: recordsFill{"Records.Fill"},
	Peek: recordsPeek{"Records.Peek"},
}

// recordsRead is the name of Records.Read, with the identities of its
// checks.
type recordsRead struct {
	conform.Callable[Records]
}

// Smoke returns the identity of the check that Read returns without a
// panic for derived inputs.
func (c recordsRead) Smoke() conform.ID[Records] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Read returns an error for
// a cancelled context.
func (c recordsRead) Cancel() conform.ID[Records] {
	return c.Check("cancel")
}

// ReturnsOnError returns the identity of the check that Read returns the
// field of its request beside an error.
func (c recordsRead) ReturnsOnError() conform.ID[Records] {
	return c.Check("returns-on-error")
}

// recordsFill is the name of Records.Fill, with the identities of its
// checks.
type recordsFill struct {
	conform.Callable[Records]
}

// Smoke returns the identity of the check that Fill returns without a
// panic for derived inputs.
func (c recordsFill) Smoke() conform.ID[Records] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Fill returns an error for
// a cancelled context.
func (c recordsFill) Cancel() conform.ID[Records] {
	return c.Check("cancel")
}

// ReturnsOnError returns the identity of the check that Fill returns the
// field of its request beside an error.
func (c recordsFill) ReturnsOnError() conform.ID[Records] {
	return c.Check("returns-on-error")
}

// recordsPeek is the name of Records.Peek, with the identities of its
// checks.
type recordsPeek struct {
	conform.Callable[Records]
}

// Smoke returns the identity of the check that Peek returns without a
// panic for derived inputs.
func (c recordsPeek) Smoke() conform.ID[Records] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Peek returns an error for
// a cancelled context.
func (c recordsPeek) Cancel() conform.ID[Records] {
	return c.Check("cancel")
}

// RecordsSuite returns the checks of Records, and the calls and the
// budgets of its callables.
func RecordsSuite() conform.Suite[Records] {
	return conform.Suite[Records]{
		Declaration: "Records",
		Package:     "example.com/shop/buffered",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Records]{
			{
				ID:          RecordsChecks.Read.Smoke(),
				Statement:   "Read returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RecordsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, &Request{Key: "Key-a"})
						return err
					})
				},
				Defect: recordsReadPanics,
			},
			{
				ID:          RecordsChecks.Read.Smoke().Generated(),
				Statement:   "Read returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, &drawnReq)
						return err
					})
				},
				Defect: recordsReadPanics,
			},
			{
				ID:          RecordsChecks.Read.Cancel(),
				Statement:   "Read returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RecordsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, &Request{Key: "Key-a"})
						return err
					})
				},
				Defect: recordsReadIgnoresCancel,
			},
			{
				ID:          RecordsChecks.Read.Cancel().Generated(),
				Statement:   "Read returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Read(ctx, &drawnReq)
						return err
					})
				},
				Defect: recordsReadIgnoresCancel,
			},
			{
				ID:          RecordsChecks.Read.ReturnsOnError(),
				Statement:   "Read returns the field of its request beside an error",
				Claim:       spec.ReturnsOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RecordsHandle) {
					probe.FieldBesideError(h, func(ctx context.Context) (Result, error) {
						req := &Request{Key: "Key-a"}
						req.Dst = []byte{1}
						return h.Subject().Read(ctx, req)
					}, func(v Result) []byte { return v.Record }, []byte{1})
				},
				Defect: recordsReadDropsTheField,
			},
			{
				ID:          RecordsChecks.Read.ReturnsOnError().Generated(),
				Statement:   "Read returns the field of its request beside an error",
				Claim:       spec.ReturnsOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.FieldBesideError(h, func(ctx context.Context) (Result, error) {
						return h.Subject().Read(ctx, &drawnReq)
					}, func(v Result) []byte { return v.Record }, (&drawnReq).Dst)
				},
				Defect: recordsReadDropsTheField,
			},
			{
				ID:          RecordsChecks.Fill.Smoke(),
				Statement:   "Fill returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RecordsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fill(ctx, Request{Key: "Key-a"})
						return err
					})
				},
				Defect: recordsFillPanics,
			},
			{
				ID:          RecordsChecks.Fill.Smoke().Generated(),
				Statement:   "Fill returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fill(ctx, drawnReq)
						return err
					})
				},
				Defect: recordsFillPanics,
			},
			{
				ID:          RecordsChecks.Fill.Cancel(),
				Statement:   "Fill returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RecordsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fill(ctx, Request{Key: "Key-a"})
						return err
					})
				},
				Defect: recordsFillIgnoresCancel,
			},
			{
				ID:          RecordsChecks.Fill.Cancel().Generated(),
				Statement:   "Fill returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fill(ctx, drawnReq)
						return err
					})
				},
				Defect: recordsFillIgnoresCancel,
			},
			{
				ID:          RecordsChecks.Fill.ReturnsOnError(),
				Statement:   "Fill returns the field of its request beside an error",
				Claim:       spec.ReturnsOnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RecordsHandle) {
					probe.FieldBesideError(h, func(ctx context.Context) ([]byte, error) {
						req := Request{Key: "Key-a"}
						req.Dst = []byte{1}
						return h.Subject().Fill(ctx, req)
					}, func(v []byte) []byte { return v }, []byte{1})
				},
				Defect: recordsFillDropsTheField,
			},
			{
				ID:          RecordsChecks.Fill.ReturnsOnError().Generated(),
				Statement:   "Fill returns the field of its request beside an error",
				Claim:       spec.ReturnsOnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.FieldBesideError(h, func(ctx context.Context) ([]byte, error) {
						return h.Subject().Fill(ctx, drawnReq)
					}, func(v []byte) []byte { return v }, drawnReq.Dst)
				},
				Defect: recordsFillDropsTheField,
			},
			{
				ID:          RecordsChecks.Peek.Smoke(),
				Statement:   "Peek returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RecordsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, Request{Key: "Key-a"})
						return err
					})
				},
				Defect: recordsPeekPanics,
			},
			{
				ID:          RecordsChecks.Peek.Smoke().Generated(),
				Statement:   "Peek returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, drawnReq)
						return err
					})
				},
				Defect: recordsPeekPanics,
			},
			{
				ID:          RecordsChecks.Peek.Cancel(),
				Statement:   "Peek returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RecordsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, Request{Key: "Key-a"})
						return err
					})
				},
				Defect: recordsPeekIgnoresCancel,
			},
			{
				ID:          RecordsChecks.Peek.Cancel().Generated(),
				Statement:   "Peek returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RecordsHandle, c *prop.Case) {
					drawnReq := c.Draw(prop.Of[Request](), "req")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, drawnReq)
						return err
					})
				},
				Defect: recordsPeekIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.ZeroOnError,
				Callable: "Records.Read",
				Gate:     spec.ReturnsNoField,
				Why:      "Read returns a field of its request beside an error, as its mixin returns-on-error states",
				Remedy:   "remove the mixin returns-on-error where Read returns the zero value beside an error",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Records.Fill",
				Gate:     spec.ReturnsNoField,
				Why:      "Fill returns a field of its request beside an error, as its mixin returns-on-error states",
				Remedy:   "remove the mixin returns-on-error where Fill returns the zero value beside an error",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Records.Peek",
				Gate:     spec.ReturnsNoField,
				Why:      "Peek returns a field of its request beside an error, as its mixin returns-on-error states",
				Remedy:   "remove the mixin returns-on-error where Peek returns the zero value beside an error",
			},
			{
				Claim:    spec.ReturnsOnError,
				Callable: "Records.Peek",
				Gate:     spec.ClassificationParams,
				Why:      "the classification of Peek sets no param {param}",
				Remedy:   "set the param {param} in the directive of Peek",
			},
		},
		Calls: []conform.Call[Records]{
			{
				Calls:    recordsReadCalls,
				Over:     newRecordsOver,
				Callable: RecordsChecks.Read.Callable,
			},
			{
				Calls:    recordsFillCalls,
				Over:     newRecordsOver,
				Callable: RecordsChecks.Fill.Callable,
			},
			{
				Calls:    recordsPeekCalls,
				Over:     newRecordsOver,
				Callable: RecordsChecks.Peek.Callable,
			},
		},
	}
}

// RunRecordsSuite applies the checks of Records, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunRecordsSuite[T Records](t *testing.T, newSubject func() T, opts ...RecordsOption) {
	t.Helper()
	RunRecordsSuiteFromHandle(t, func(*RecordsHandle) T { return newSubject() }, opts...)
}

// RunRecordsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRecordsSuiteFromHandle[T Records](
	t *testing.T,
	newSubject func(h *RecordsHandle) T,
	opts ...RecordsOption,
) {
	t.Helper()
	conform.Run(t, RecordsSuite(), recordsSubjects(newSubject), opts...)
}

// FuzzRecordsSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRecordsSuite[T Records](
	f *testing.F,
	callable conform.Callable[Records],
	newSubject func(h *RecordsHandle) T,
	opts ...RecordsOption,
) {
	f.Helper()
	conform.Fuzz(f, RecordsSuite(), callable, recordsSubjects(newSubject), opts...)
}

// RunRecordsBudgets applies the allocation ceilings of Records and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRecordsBudgets[T Records](
	t *testing.T,
	newSubject func(h *RecordsHandle) T,
	opts ...RecordsOption,
) {
	t.Helper()
	conform.RunBudgets(t, RecordsSuite(), recordsSubjects(newSubject), opts...)
}

// BenchRecordsSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchRecordsSuite[T Records](
	b conform.B,
	callable conform.Callable[Records],
	newSubject func(h *RecordsHandle) T,
	opts ...RecordsOption,
) {
	b.Helper()
	conform.Bench(b, RecordsSuite(), callable, recordsSubjects(newSubject), opts...)
}

// recordsSubjects returns the factory of the subjects of Records that
// calls newSubject with the handle of each check.
func recordsSubjects[T Records](
	newSubject func(h *RecordsHandle) T,
) func(h *RecordsHandle) Records {
	return func(h *RecordsHandle) Records { return newSubject(h) }
}

// recordsReadCalls generates the calls of Records.Read with the inputs
// that each case draws.
var recordsReadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Records) error {
	req := c.Draw(prop.Of[Request](), "req")
	return func(ctx context.Context, subject Records) error {
		_, err := subject.Read(ctx, &req)
		return err
	}
})

// recordsFillCalls generates the calls of Records.Fill with the inputs
// that each case draws.
var recordsFillCalls = prop.Composite(func(c *prop.Case) func(context.Context, Records) error {
	req := c.Draw(prop.Of[Request](), "req")
	return func(ctx context.Context, subject Records) error {
		_, err := subject.Fill(ctx, req)
		return err
	}
})

// recordsPeekCalls generates the calls of Records.Peek with the inputs
// that each case draws.
var recordsPeekCalls = prop.Composite(func(c *prop.Case) func(context.Context, Records) error {
	req := c.Draw(prop.Of[Request](), "req")
	return func(ctx context.Context, subject Records) error {
		_, err := subject.Peek(ctx, req)
		return err
	}
})

// recordsOverEmbedded is Records under the name of the field that
// recordsOver embeds.
type recordsOverEmbedded = Records

// recordsOver is a subject of Records whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type recordsOver struct {
	recordsOverEmbedded
	spend func()
}

// Read calls spend and returns its zero results.
func (o recordsOver) Read(context.Context, *Request) (r0 Result, r1 error) {
	o.spend()
	return r0, r1
}

// Fill calls spend and returns its zero results.
func (o recordsOver) Fill(context.Context, Request) (r0 []byte, r1 error) {
	o.spend()
	return r0, r1
}

// Peek calls spend and returns its zero results.
func (o recordsOver) Peek(context.Context, Request) (r0 Result, r1 error) {
	o.spend()
	return r0, r1
}

// newRecordsOver returns the planted defect of a ceiling of Records, whose
// callables call spend.
func newRecordsOver(_ *RecordsHandle, spend func()) Records {
	return recordsOver{spend: spend}
}

// recordsReadPanicsDefect is the planted defect of a Records whose Read
// panics.
type recordsReadPanicsDefect struct {
	recordsOverEmbedded
}

// Read is the override of the planted defect.
func (*recordsReadPanicsDefect) Read(p0 context.Context, p1 *Request) (r0 Result, r1 error) {
	panic("the planted defect of a Records whose Read panics")
}

// recordsReadPanics returns the planted defect of a Records whose Read
// panics over a fresh subject.
func recordsReadPanics(h *RecordsHandle) Records {
	return &recordsReadPanicsDefect{recordsOverEmbedded: h.Fresh()}
}

// recordsReadIgnoresCancelDefect is the planted defect of a Records whose
// Read ignores the context.
type recordsReadIgnoresCancelDefect struct {
	recordsOverEmbedded
}

// Read is the override of the planted defect.
func (d *recordsReadIgnoresCancelDefect) Read(
	p0 context.Context,
	p1 *Request,
) (r0 Result, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.recordsOverEmbedded.Read(p0, p1)
}

// recordsReadIgnoresCancel returns the planted defect of a Records whose
// Read ignores the context over a fresh subject.
func recordsReadIgnoresCancel(h *RecordsHandle) Records {
	return &recordsReadIgnoresCancelDefect{recordsOverEmbedded: h.Fresh()}
}

// recordsReadDropsTheFieldDefect is the planted defect of a Records whose
// Read returns the zero value of its result in place of the field of its
// request.
type recordsReadDropsTheFieldDefect struct {
	recordsOverEmbedded
}

// Read is the override of the planted defect.
func (d *recordsReadDropsTheFieldDefect) Read(
	p0 context.Context,
	p1 *Request,
) (r0 Result, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.recordsOverEmbedded.Read(p0, p1)
	if r1 != nil {
		r0 = *new(Result)
	}
	return r0, r1
}

// recordsReadDropsTheField returns the planted defect of a Records whose
// Read returns the zero value of its result in place of the field of its
// request over a fresh subject.
func recordsReadDropsTheField(h *RecordsHandle) Records {
	return &recordsReadDropsTheFieldDefect{recordsOverEmbedded: h.Fresh()}
}

// recordsFillPanicsDefect is the planted defect of a Records whose Fill
// panics.
type recordsFillPanicsDefect struct {
	recordsOverEmbedded
}

// Fill is the override of the planted defect.
func (*recordsFillPanicsDefect) Fill(p0 context.Context, p1 Request) (r0 []byte, r1 error) {
	panic("the planted defect of a Records whose Fill panics")
}

// recordsFillPanics returns the planted defect of a Records whose Fill
// panics over a fresh subject.
func recordsFillPanics(h *RecordsHandle) Records {
	return &recordsFillPanicsDefect{recordsOverEmbedded: h.Fresh()}
}

// recordsFillIgnoresCancelDefect is the planted defect of a Records whose
// Fill ignores the context.
type recordsFillIgnoresCancelDefect struct {
	recordsOverEmbedded
}

// Fill is the override of the planted defect.
func (d *recordsFillIgnoresCancelDefect) Fill(
	p0 context.Context,
	p1 Request,
) (r0 []byte, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.recordsOverEmbedded.Fill(p0, p1)
}

// recordsFillIgnoresCancel returns the planted defect of a Records whose
// Fill ignores the context over a fresh subject.
func recordsFillIgnoresCancel(h *RecordsHandle) Records {
	return &recordsFillIgnoresCancelDefect{recordsOverEmbedded: h.Fresh()}
}

// recordsFillDropsTheFieldDefect is the planted defect of a Records whose
// Fill returns the zero value of its result in place of the field of its
// request.
type recordsFillDropsTheFieldDefect struct {
	recordsOverEmbedded
}

// Fill is the override of the planted defect.
func (d *recordsFillDropsTheFieldDefect) Fill(
	p0 context.Context,
	p1 Request,
) (r0 []byte, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.recordsOverEmbedded.Fill(p0, p1)
	if r1 != nil {
		r0 = *new([]byte)
	}
	return r0, r1
}

// recordsFillDropsTheField returns the planted defect of a Records whose
// Fill returns the zero value of its result in place of the field of its
// request over a fresh subject.
func recordsFillDropsTheField(h *RecordsHandle) Records {
	return &recordsFillDropsTheFieldDefect{recordsOverEmbedded: h.Fresh()}
}

// recordsPeekPanicsDefect is the planted defect of a Records whose Peek
// panics.
type recordsPeekPanicsDefect struct {
	recordsOverEmbedded
}

// Peek is the override of the planted defect.
func (*recordsPeekPanicsDefect) Peek(p0 context.Context, p1 Request) (r0 Result, r1 error) {
	panic("the planted defect of a Records whose Peek panics")
}

// recordsPeekPanics returns the planted defect of a Records whose Peek
// panics over a fresh subject.
func recordsPeekPanics(h *RecordsHandle) Records {
	return &recordsPeekPanicsDefect{recordsOverEmbedded: h.Fresh()}
}

// recordsPeekIgnoresCancelDefect is the planted defect of a Records whose
// Peek ignores the context.
type recordsPeekIgnoresCancelDefect struct {
	recordsOverEmbedded
}

// Peek is the override of the planted defect.
func (d *recordsPeekIgnoresCancelDefect) Peek(
	p0 context.Context,
	p1 Request,
) (r0 Result, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.recordsOverEmbedded.Peek(p0, p1)
}

// recordsPeekIgnoresCancel returns the planted defect of a Records whose
// Peek ignores the context over a fresh subject.
func recordsPeekIgnoresCancel(h *RecordsHandle) Records {
	return &recordsPeekIgnoresCancelDefect{recordsOverEmbedded: h.Fresh()}
}
