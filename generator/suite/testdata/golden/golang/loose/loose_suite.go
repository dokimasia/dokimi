package loose

import (
	"context"
	"math/big"
	"testing"
	"time"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Printer writes.
type (
	PrinterHandle = conform.Handle[Printer]
	PrinterCheck  = conform.Check[Printer]
	PrinterOption = conform.Option[Printer]
)

// PrinterChecks contains the name of each callable of Printer, with the
// identities of its checks.
var PrinterChecks = struct {
	Print printerPrint
}{
	Print: printerPrint{"Printer.Print"},
}

// printerPrint is the name of Printer.Print, with the identities of its
// checks.
type printerPrint struct {
	conform.Callable[Printer]
}

// Smoke returns the identity of the check that Print returns without a
// panic for derived inputs.
func (c printerPrint) Smoke() conform.ID[Printer] {
	return c.Check("smoke")
}

// PrinterSuite returns the checks of Printer, and the calls and the
// budgets of its callables.
func PrinterSuite() conform.Suite[Printer] {
	return conform.Suite[Printer]{
		Declaration: "Printer",
		Checks: []conform.Check[Printer]{
			{
				ID:          PrinterChecks.Print.Smoke(),
				Statement:   "Print returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PrinterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Print(string("msg-a"), any("values-a"))
						return nil
					})
				},
				Defect: printerPrintPanics,
			},
		},
		Calls: []conform.Call[Printer]{
			{
				Calls:    printerPrintCalls,
				Over:     newPrinterOver,
				Callable: PrinterChecks.Print.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) any {
				switch c.Draw(prop.Integer(0, 2), "any") {
				case 0:
					return c.Draw(prop.String(), "any.string")
				case 1:
					return c.Draw(prop.Of[int](), "any.int")
				}
				return c.Draw(prop.Boolean(), "any.bool")
			})),
		},
	}
}

// RunPrinterSuite applies the checks of Printer, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunPrinterSuite[T Printer](t *testing.T, newSubject func() T, opts ...PrinterOption) {
	t.Helper()
	RunPrinterSuiteFromHandle(t, func(*PrinterHandle) T { return newSubject() }, opts...)
}

// RunPrinterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPrinterSuiteFromHandle[T Printer](
	t *testing.T,
	newSubject func(h *PrinterHandle) T,
	opts ...PrinterOption,
) {
	t.Helper()
	conform.Run(t, PrinterSuite(), printerSubjects(newSubject), opts...)
}

// FuzzPrinterSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPrinterSuite[T Printer](
	f *testing.F,
	callable conform.Callable[Printer],
	newSubject func(h *PrinterHandle) T,
	opts ...PrinterOption,
) {
	f.Helper()
	conform.Fuzz(f, PrinterSuite(), callable, printerSubjects(newSubject), opts...)
}

// RunPrinterBudgets applies the allocation ceilings of Printer and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunPrinterBudgets[T Printer](
	t *testing.T,
	newSubject func(h *PrinterHandle) T,
	opts ...PrinterOption,
) {
	t.Helper()
	conform.RunBudgets(t, PrinterSuite(), printerSubjects(newSubject), opts...)
}

// BenchPrinterSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchPrinterSuite[T Printer](
	b conform.B,
	callable conform.Callable[Printer],
	newSubject func(h *PrinterHandle) T,
	opts ...PrinterOption,
) {
	b.Helper()
	conform.Bench(b, PrinterSuite(), callable, printerSubjects(newSubject), opts...)
}

// printerSubjects returns the factory of the subjects of Printer that
// calls newSubject with the handle of each check.
func printerSubjects[T Printer](
	newSubject func(h *PrinterHandle) T,
) func(h *PrinterHandle) Printer {
	return func(h *PrinterHandle) Printer { return newSubject(h) }
}

// printerPrintCalls generates the calls of Printer.Print with the inputs
// that each case draws.
var printerPrintCalls = prop.Composite(func(c *prop.Case) func(context.Context, Printer) error {
	msg := c.Draw(prop.Of[string](), "msg")
	values := c.Draw(prop.List(prop.Of[any]()), "values")
	return func(_ context.Context, subject Printer) error {
		subject.Print(msg, values...)
		return nil
	}
})

// printerOverEmbedded is Printer under the name of the field that
// printerOver embeds.
type printerOverEmbedded = Printer

// printerOver is a subject of Printer whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type printerOver struct {
	printerOverEmbedded
	spend func()
}

// Print calls spend and returns its zero results.
func (o printerOver) Print(string, ...any) {
	o.spend()
}

// newPrinterOver returns the planted defect of a ceiling of Printer, whose
// callables call spend.
func newPrinterOver(_ *PrinterHandle, spend func()) Printer {
	return printerOver{spend: spend}
}

// printerPrintPanicsDefect is the planted defect of a Printer whose Print
// panics.
type printerPrintPanicsDefect struct {
	printerOverEmbedded
}

// Print is the override of the planted defect.
func (*printerPrintPanicsDefect) Print(p0 string, p1 ...any) {
	panic("the planted defect of a Printer whose Print panics")
}

// printerPrintPanics returns the planted defect of a Printer whose Print
// panics over a fresh subject.
func printerPrintPanics(h *PrinterHandle) Printer {
	return &printerPrintPanicsDefect{printerOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Recorder writes.
type (
	RecorderHandle = conform.Handle[Recorder]
	RecorderCheck  = conform.Check[Recorder]
	RecorderOption = conform.Option[Recorder]
)

// RecorderChecks contains the name of each callable of Recorder, with the
// identities of its checks.
var RecorderChecks = struct {
	Field recorderField
}{
	Field: recorderField{"Recorder.Field"},
}

// recorderField is the name of Recorder.Field, with the identities of its
// checks.
type recorderField struct {
	conform.Callable[Recorder]
}

// Smoke returns the identity of the check that Field returns without a
// panic for derived inputs.
func (c recorderField) Smoke() conform.ID[Recorder] {
	return c.Check("smoke")
}

// RecorderSuite returns the checks of Recorder, and the calls and the
// budgets of its callables.
func RecorderSuite() conform.Suite[Recorder] {
	return conform.Suite[Recorder]{
		Declaration: "Recorder",
		Checks: []conform.Check[Recorder]{
			{
				ID:          RecorderChecks.Field.Smoke(),
				Statement:   "Field returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RecorderHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Field(string("name-a"), interface{}("value-a"))
						return nil
					})
				},
				Defect: recorderFieldPanics,
			},
		},
		Calls: []conform.Call[Recorder]{
			{
				Calls:    recorderFieldCalls,
				Over:     newRecorderOver,
				Callable: RecorderChecks.Field.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) any {
				switch c.Draw(prop.Integer(0, 2), "any") {
				case 0:
					return c.Draw(prop.String(), "any.string")
				case 1:
					return c.Draw(prop.Of[int](), "any.int")
				}
				return c.Draw(prop.Boolean(), "any.bool")
			})),
		},
	}
}

// RunRecorderSuite applies the checks of Recorder, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunRecorderSuite[T Recorder](t *testing.T, newSubject func() T, opts ...RecorderOption) {
	t.Helper()
	RunRecorderSuiteFromHandle(t, func(*RecorderHandle) T { return newSubject() }, opts...)
}

// RunRecorderSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRecorderSuiteFromHandle[T Recorder](
	t *testing.T,
	newSubject func(h *RecorderHandle) T,
	opts ...RecorderOption,
) {
	t.Helper()
	conform.Run(t, RecorderSuite(), recorderSubjects(newSubject), opts...)
}

// FuzzRecorderSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRecorderSuite[T Recorder](
	f *testing.F,
	callable conform.Callable[Recorder],
	newSubject func(h *RecorderHandle) T,
	opts ...RecorderOption,
) {
	f.Helper()
	conform.Fuzz(f, RecorderSuite(), callable, recorderSubjects(newSubject), opts...)
}

// RunRecorderBudgets applies the allocation ceilings of Recorder and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRecorderBudgets[T Recorder](
	t *testing.T,
	newSubject func(h *RecorderHandle) T,
	opts ...RecorderOption,
) {
	t.Helper()
	conform.RunBudgets(t, RecorderSuite(), recorderSubjects(newSubject), opts...)
}

// BenchRecorderSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchRecorderSuite[T Recorder](
	b conform.B,
	callable conform.Callable[Recorder],
	newSubject func(h *RecorderHandle) T,
	opts ...RecorderOption,
) {
	b.Helper()
	conform.Bench(b, RecorderSuite(), callable, recorderSubjects(newSubject), opts...)
}

// recorderSubjects returns the factory of the subjects of Recorder that
// calls newSubject with the handle of each check.
func recorderSubjects[T Recorder](
	newSubject func(h *RecorderHandle) T,
) func(h *RecorderHandle) Recorder {
	return func(h *RecorderHandle) Recorder { return newSubject(h) }
}

// recorderFieldCalls generates the calls of Recorder.Field with the inputs
// that each case draws.
var recorderFieldCalls = prop.Composite(func(c *prop.Case) func(context.Context, Recorder) error {
	name := c.Draw(prop.Of[string](), "name")
	value := c.Draw(prop.Of[interface{}](), "value")
	return func(_ context.Context, subject Recorder) error {
		subject.Field(name, value)
		return nil
	}
})

// recorderOverEmbedded is Recorder under the name of the field that
// recorderOver embeds.
type recorderOverEmbedded = Recorder

// recorderOver is a subject of Recorder whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type recorderOver struct {
	recorderOverEmbedded
	spend func()
}

// Field calls spend and returns its zero results.
func (o recorderOver) Field(string, interface{}) {
	o.spend()
}

// newRecorderOver returns the planted defect of a ceiling of Recorder,
// whose callables call spend.
func newRecorderOver(_ *RecorderHandle, spend func()) Recorder {
	return recorderOver{spend: spend}
}

// recorderFieldPanicsDefect is the planted defect of a Recorder whose
// Field panics.
type recorderFieldPanicsDefect struct {
	recorderOverEmbedded
}

// Field is the override of the planted defect.
func (*recorderFieldPanicsDefect) Field(p0 string, p1 interface{}) {
	panic("the planted defect of a Recorder whose Field panics")
}

// recorderFieldPanics returns the planted defect of a Recorder whose Field
// panics over a fresh subject.
func recorderFieldPanics(h *RecorderHandle) Recorder {
	return &recorderFieldPanicsDefect{recorderOverEmbedded: h.Fresh()}
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
	Put journalPut
}{
	Put: journalPut{"Journal.Put"},
}

// journalPut is the name of Journal.Put, with the identities of its
// checks.
type journalPut struct {
	conform.Callable[Journal]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c journalPut) Smoke() conform.ID[Journal] {
	return c.Check("smoke")
}

// JournalSuite returns the checks of Journal, and the calls and the
// budgets of its callables.
func JournalSuite() conform.Suite[Journal] {
	return conform.Suite[Journal]{
		Declaration: "Journal",
		Checks: []conform.Check[Journal]{
			{
				ID:          JournalChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *JournalHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Put(Entry{Value: "Value-a"})
						return nil
					})
				},
				Defect: journalPutPanics,
			},
		},
		Calls: []conform.Call[Journal]{
			{
				Calls:    journalPutCalls,
				Over:     newJournalOver,
				Callable: JournalChecks.Put.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) any {
				switch c.Draw(prop.Integer(0, 2), "any") {
				case 0:
					return c.Draw(prop.String(), "any.string")
				case 1:
					return c.Draw(prop.Of[int](), "any.int")
				}
				return c.Draw(prop.Boolean(), "any.bool")
			})),
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

// journalPutCalls generates the calls of Journal.Put with the inputs that
// each case draws.
var journalPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Journal) error {
	e := c.Draw(prop.Of[Entry](), "e")
	return func(_ context.Context, subject Journal) error {
		subject.Put(e)
		return nil
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

// Put calls spend and returns its zero results.
func (o journalOver) Put(Entry) {
	o.spend()
}

// newJournalOver returns the planted defect of a ceiling of Journal, whose
// callables call spend.
func newJournalOver(_ *JournalHandle, spend func()) Journal {
	return journalOver{spend: spend}
}

// journalPutPanicsDefect is the planted defect of a Journal whose Put
// panics.
type journalPutPanicsDefect struct {
	journalOverEmbedded
}

// Put is the override of the planted defect.
func (*journalPutPanicsDefect) Put(p0 Entry) {
	panic("the planted defect of a Journal whose Put panics")
}

// journalPutPanics returns the planted defect of a Journal whose Put
// panics over a fresh subject.
func journalPutPanics(h *JournalHandle) Journal {
	return &journalPutPanicsDefect{journalOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Plain writes.
type (
	PlainHandle = conform.Handle[Plain]
	PlainCheck  = conform.Check[Plain]
	PlainOption = conform.Option[Plain]
)

// PlainChecks contains the name of each callable of Plain, with the
// identities of its checks.
var PlainChecks = struct {
	Print plainPrint
}{
	Print: plainPrint{"Plain.Print"},
}

// plainPrint is the name of Plain.Print, with the identities of its
// checks.
type plainPrint struct {
	conform.Callable[Plain]
}

// Smoke returns the identity of the check that Print returns without a
// panic for derived inputs.
func (c plainPrint) Smoke() conform.ID[Plain] {
	return c.Check("smoke")
}

// PlainSuite returns the checks of Plain, and the calls and the budgets of
// its callables.
func PlainSuite() conform.Suite[Plain] {
	return conform.Suite[Plain]{
		Declaration: "Plain",
		Checks: []conform.Check[Plain]{
			{
				ID:          PlainChecks.Print.Smoke(),
				Statement:   "Print returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlainHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Print(string("msg-a"))
						return nil
					})
				},
				Defect: plainPrintPanics,
			},
		},
		Calls: []conform.Call[Plain]{
			{
				Calls:    plainPrintCalls,
				Over:     newPlainOver,
				Callable: PlainChecks.Print.Callable,
			},
		},
	}
}

// RunPlainSuite applies the checks of Plain, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunPlainSuite[T Plain](t *testing.T, newSubject func() T, opts ...PlainOption) {
	t.Helper()
	RunPlainSuiteFromHandle(t, func(*PlainHandle) T { return newSubject() }, opts...)
}

// RunPlainSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPlainSuiteFromHandle[T Plain](
	t *testing.T,
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	t.Helper()
	conform.Run(t, PlainSuite(), plainSubjects(newSubject), opts...)
}

// FuzzPlainSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPlainSuite[T Plain](
	f *testing.F,
	callable conform.Callable[Plain],
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	f.Helper()
	conform.Fuzz(f, PlainSuite(), callable, plainSubjects(newSubject), opts...)
}

// RunPlainBudgets applies the allocation ceilings of Plain and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunPlainBudgets[T Plain](
	t *testing.T,
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	t.Helper()
	conform.RunBudgets(t, PlainSuite(), plainSubjects(newSubject), opts...)
}

// BenchPlainSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPlainSuite[T Plain](
	b conform.B,
	callable conform.Callable[Plain],
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	b.Helper()
	conform.Bench(b, PlainSuite(), callable, plainSubjects(newSubject), opts...)
}

// plainSubjects returns the factory of the subjects of Plain that calls
// newSubject with the handle of each check.
func plainSubjects[T Plain](newSubject func(h *PlainHandle) T) func(h *PlainHandle) Plain {
	return func(h *PlainHandle) Plain { return newSubject(h) }
}

// plainPrintCalls generates the calls of Plain.Print with the inputs that
// each case draws.
var plainPrintCalls = prop.Composite(func(c *prop.Case) func(context.Context, Plain) error {
	msg := c.Draw(prop.Of[string](), "msg")
	return func(_ context.Context, subject Plain) error {
		subject.Print(msg)
		return nil
	}
})

// plainOverEmbedded is Plain under the name of the field that plainOver
// embeds.
type plainOverEmbedded = Plain

// plainOver is a subject of Plain whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type plainOver struct {
	plainOverEmbedded
	spend func()
}

// Print calls spend and returns its zero results.
func (o plainOver) Print(string) {
	o.spend()
}

// newPlainOver returns the planted defect of a ceiling of Plain, whose
// callables call spend.
func newPlainOver(_ *PlainHandle, spend func()) Plain {
	return plainOver{spend: spend}
}

// plainPrintPanicsDefect is the planted defect of a Plain whose Print
// panics.
type plainPrintPanicsDefect struct {
	plainOverEmbedded
}

// Print is the override of the planted defect.
func (*plainPrintPanicsDefect) Print(p0 string) {
	panic("the planted defect of a Plain whose Print panics")
}

// plainPrintPanics returns the planted defect of a Plain whose Print
// panics over a fresh subject.
func plainPrintPanics(h *PlainHandle) Plain {
	return &plainPrintPanicsDefect{plainOverEmbedded: h.Fresh()}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Zoned writes.
type (
	ZonedHandle = conform.Handle[Zoned]
	ZonedCheck  = conform.Check[Zoned]
	ZonedOption = conform.Option[Zoned]
)

// ZonedChecks contains the name of each callable of Zoned, with the
// identities of its checks.
var ZonedChecks = struct {
	In    zonedIn
	Scale zonedScale
	Ratio zonedRatio
}{
	In:    zonedIn{"Zoned.In"},
	Scale: zonedScale{"Zoned.Scale"},
	Ratio: zonedRatio{"Zoned.Ratio"},
}

// zonedIn is the name of Zoned.In, with the identities of its checks.
type zonedIn struct {
	conform.Callable[Zoned]
}

// zonedScale is the name of Zoned.Scale, with the identities of its
// checks.
type zonedScale struct {
	conform.Callable[Zoned]
}

// zonedRatio is the name of Zoned.Ratio, with the identities of its
// checks.
type zonedRatio struct {
	conform.Callable[Zoned]
}

// ZonedSuite returns the checks of Zoned, and the calls and the budgets of
// its callables.
func ZonedSuite() conform.Suite[Zoned] {
	return conform.Suite[Zoned]{
		Declaration: "Zoned",
		Calls: []conform.Call[Zoned]{
			{
				Calls:    zonedInCalls,
				Over:     newZonedOver,
				Callable: ZonedChecks.In.Callable,
			},
			{
				Calls:    zonedScaleCalls,
				Over:     newZonedOver,
				Callable: ZonedChecks.Scale.Callable,
			},
			{
				Calls:    zonedRatioCalls,
				Over:     newZonedOver,
				Callable: ZonedChecks.Ratio.Callable,
			},
		},
	}
}

// RunZonedSuite applies the checks of Zoned, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunZonedSuite[T Zoned](t *testing.T, newSubject func() T, opts ...ZonedOption) {
	t.Helper()
	RunZonedSuiteFromHandle(t, func(*ZonedHandle) T { return newSubject() }, opts...)
}

// RunZonedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunZonedSuiteFromHandle[T Zoned](
	t *testing.T,
	newSubject func(h *ZonedHandle) T,
	opts ...ZonedOption,
) {
	t.Helper()
	conform.Run(t, ZonedSuite(), zonedSubjects(newSubject), opts...)
}

// FuzzZonedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzZonedSuite[T Zoned](
	f *testing.F,
	callable conform.Callable[Zoned],
	newSubject func(h *ZonedHandle) T,
	opts ...ZonedOption,
) {
	f.Helper()
	conform.Fuzz(f, ZonedSuite(), callable, zonedSubjects(newSubject), opts...)
}

// RunZonedBudgets applies the allocation ceilings of Zoned and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunZonedBudgets[T Zoned](
	t *testing.T,
	newSubject func(h *ZonedHandle) T,
	opts ...ZonedOption,
) {
	t.Helper()
	conform.RunBudgets(t, ZonedSuite(), zonedSubjects(newSubject), opts...)
}

// BenchZonedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchZonedSuite[T Zoned](
	b conform.B,
	callable conform.Callable[Zoned],
	newSubject func(h *ZonedHandle) T,
	opts ...ZonedOption,
) {
	b.Helper()
	conform.Bench(b, ZonedSuite(), callable, zonedSubjects(newSubject), opts...)
}

// zonedSubjects returns the factory of the subjects of Zoned that calls
// newSubject with the handle of each check.
func zonedSubjects[T Zoned](newSubject func(h *ZonedHandle) T) func(h *ZonedHandle) Zoned {
	return func(h *ZonedHandle) Zoned { return newSubject(h) }
}

// zonedInCalls generates the calls of Zoned.In with the inputs that each
// case draws.
var zonedInCalls = prop.Composite(func(c *prop.Case) func(context.Context, Zoned) error {
	loc := c.Draw(prop.Of[*time.Location](), "loc")
	return func(_ context.Context, subject Zoned) error {
		subject.In(loc)
		return nil
	}
})

// zonedScaleCalls generates the calls of Zoned.Scale with the inputs that
// each case draws.
var zonedScaleCalls = prop.Composite(func(c *prop.Case) func(context.Context, Zoned) error {
	n := c.Draw(prop.Of[*big.Int](), "n")
	return func(_ context.Context, subject Zoned) error {
		subject.Scale(n)
		return nil
	}
})

// zonedRatioCalls generates the calls of Zoned.Ratio with the inputs that
// each case draws.
var zonedRatioCalls = prop.Composite(func(c *prop.Case) func(context.Context, Zoned) error {
	r := c.Draw(prop.Of[*big.Rat](), "r")
	return func(_ context.Context, subject Zoned) error {
		subject.Ratio(r)
		return nil
	}
})

// zonedOverEmbedded is Zoned under the name of the field that zonedOver
// embeds.
type zonedOverEmbedded = Zoned

// zonedOver is a subject of Zoned whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type zonedOver struct {
	zonedOverEmbedded
	spend func()
}

// In calls spend and returns its zero results.
func (o zonedOver) In(*time.Location) {
	o.spend()
}

// Scale calls spend and returns its zero results.
func (o zonedOver) Scale(*big.Int) {
	o.spend()
}

// Ratio calls spend and returns its zero results.
func (o zonedOver) Ratio(*big.Rat) {
	o.spend()
}

// newZonedOver returns the planted defect of a ceiling of Zoned, whose
// callables call spend.
func newZonedOver(_ *ZonedHandle, spend func()) Zoned {
	return zonedOver{spend: spend}
}
