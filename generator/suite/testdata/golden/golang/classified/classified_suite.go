package classified

import (
	"context"
	"iter"
	"slices"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Shelf:
//
// Shelf.Put/smoke: Put returns without a panic for derived inputs.
//
// Shelf.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Shelf.Put/cancel: Put returns an error for a cancelled context.
//
// Shelf.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Shelf.GetMany/smoke: GetMany returns without a panic for derived inputs.
//
// Shelf.GetMany/smoke#generated: GetMany returns without a panic for
// derived inputs.
//
// Shelf.GetMany/cancel: GetMany returns an error for a cancelled context.
//
// Shelf.GetMany/cancel#generated: GetMany returns an error for a cancelled
// context.
//
// Shelf.GetMany/reads-each-key: GetMany returns the value that Put stored
// under each key.
//
// Shelf.GetMany/reads-each-key#generated: GetMany returns the value that
// Put stored under each key.
//
// Shelf.All/smoke: All returns without a panic for derived inputs.
//
// Shelf.All/streams-what-was-stored: All yields each value that Put
// stored.
//
// Shelf.All/streams-what-was-stored#generated: All yields each value that
// Put stored.
//
// Shelf.Items/smoke: Items returns without a panic for derived inputs.
//
// Shelf.Items/streams-what-was-stored: Items yields each value that Put
// stored.
//
// Shelf.Items/streams-what-was-stored#generated: Items yields each value
// that Put stored.
//
// Shelf.Scan/smoke: Scan returns without a panic for derived inputs.
//
// Shelf.Scan/cancel: Scan returns an error for a cancelled context.
//
// Shelf.Scan/over-match: Scan gains at least the elements that Put wrote
// between two of its calls.
//
// Shelf.Scan/over-match#generated: Scan gains at least the elements that
// Put wrote between two of its calls.
//
// Shelf.Scan/drains-keep-what-went-in: Scan yields the elements of the
// call before it and each element that Put wrote since, over any sequence
// of calls.
//
// Shelf.Keyed/smoke: Keyed returns without a panic for derived inputs.
//
// Shelf.Keyed/streams-what-was-stored: Keyed yields each value that Put
// stored.
//
// Shelf.Keyed/streams-what-was-stored#generated: Keyed yields each value
// that Put stored.
//
// Shelf.Snapshot/smoke: Snapshot returns without a panic for derived
// inputs.
//
// Shelf.Snapshot/cancel: Snapshot returns an error for a cancelled
// context.
//
// Shelf.Snapshot/zero-on-error: Snapshot returns the zero value beside an
// error.
//
// Shelf.Snapshot/over-match: Snapshot gains at least the elements that Put
// wrote between two of its calls.
//
// Shelf.Snapshot/over-match#generated: Snapshot gains at least the
// elements that Put wrote between two of its calls.
//
// Shelf.Snapshot/drains-keep-what-went-in: Snapshot yields the elements of
// the call before it and each element that Put wrote since, over any
// sequence of calls.
//
// Shelf.Len/smoke: Len returns without a panic for derived inputs.
//
// Shelf.Len/cancel: Len returns an error for a cancelled context.
//
// Shelf.Len/zero-on-error: Len returns the zero value beside an error.
//
// Shelf.At/smoke: At returns without a panic for derived inputs.
//
// Shelf.At/smoke#generated: At returns without a panic for derived inputs.
//
// Shelf.At/cancel: At returns an error for a cancelled context.
//
// Shelf.At/cancel#generated: At returns an error for a cancelled context.
//
// Shelf.At/zero-on-error: At returns the zero value beside an error.
//
// Shelf.At/zero-on-error#generated: At returns the zero value beside an
// error.
//
// Shelf.At/respects-its-bounds: At admits each position below the size
// that Len reports and refuses the size.
//
// Shelf.AtOffset/smoke: AtOffset returns without a panic for derived
// inputs.
//
// Shelf.AtOffset/smoke#generated: AtOffset returns without a panic for
// derived inputs.
//
// Shelf.AtOffset/cancel: AtOffset returns an error for a cancelled
// context.
//
// Shelf.AtOffset/cancel#generated: AtOffset returns an error for a
// cancelled context.
//
// Shelf.AtOffset/zero-on-error: AtOffset returns the zero value beside an
// error.
//
// Shelf.AtOffset/zero-on-error#generated: AtOffset returns the zero value
// beside an error.
//
// Shelf.AtOffset/respects-its-bounds: AtOffset admits each position below
// the size that Len reports and refuses the size.
//
// Shelf.Slot/smoke: Slot returns without a panic for derived inputs.
//
// Shelf.Slot/smoke#generated: Slot returns without a panic for derived
// inputs.
//
// Shelf.Slot/cancel: Slot returns an error for a cancelled context.
//
// Shelf.Slot/cancel#generated: Slot returns an error for a cancelled
// context.
//
// Shelf.Slot/respects-its-bounds: Slot admits each position below the size
// that Len reports and refuses the size.
//
// Shelf.Slot/respects-its-bounds#generated: Slot admits each position
// below the size that Len reports and refuses the size.
//
// Shelf.Get/smoke: Get returns without a panic for derived inputs.
//
// Shelf.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Shelf.Get/cancel: Get returns an error for a cancelled context.
//
// Shelf.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Shelf.Get/zero-on-error: Get returns the zero value beside an error.
//
// Shelf.Get/zero-on-error#generated: Get returns the zero value beside an
// error.
//
// Shelf.TracedGet/smoke: TracedGet returns without a panic for derived
// inputs.
//
// Shelf.TracedGet/smoke#generated: TracedGet returns without a panic for
// derived inputs.
//
// Shelf.TracedGet/cancel: TracedGet returns an error for a cancelled
// context.
//
// Shelf.TracedGet/cancel#generated: TracedGet returns an error for a
// cancelled context.
//
// Shelf.TracedGet/zero-on-error: TracedGet returns the zero value beside
// an error.
//
// Shelf.TracedGet/zero-on-error#generated: TracedGet returns the zero
// value beside an error.
//
// Shelf.TracedGet/agrees-with-its-sibling: TracedGet returns what Get
// returns for one input.
//
// Shelf.TracedGet/agrees-with-its-sibling#generated: TracedGet returns
// what Get returns for one input.
//
// Shelf.Err/smoke: Err returns without a panic for derived inputs.
//
// Shelf.Err/fresh-reports-no-poison: Err reports no error for a fresh
// receiver.

// The derived checks of Shelf that witness their laws at the shallow end:
//
// Shelf.Put/smoke
//
// Shelf.GetMany/smoke
//
// Shelf.All/smoke
//
// Shelf.All/streams-what-was-stored
//
// Shelf.All/streams-what-was-stored#generated
//
// Shelf.Items/smoke
//
// Shelf.Items/streams-what-was-stored
//
// Shelf.Items/streams-what-was-stored#generated
//
// Shelf.Scan/smoke
//
// Shelf.Scan/over-match
//
// Shelf.Scan/over-match#generated
//
// Shelf.Keyed/smoke
//
// Shelf.Keyed/streams-what-was-stored
//
// Shelf.Keyed/streams-what-was-stored#generated
//
// Shelf.Snapshot/smoke
//
// Shelf.Snapshot/over-match
//
// Shelf.Snapshot/over-match#generated
//
// Shelf.Len/smoke
//
// Shelf.At/smoke
//
// Shelf.AtOffset/smoke
//
// Shelf.Slot/smoke
//
// Shelf.Get/smoke
//
// Shelf.TracedGet/smoke
//
// Shelf.Err/smoke

// The refused claims of Shelf:
//
// Shelf.All/cancel: All cannot fail. The remedy is to give All an error
// result.
//
// Shelf.Items/cancel: Items cannot fail. The remedy is to give Items an
// error result.
//
// Shelf.Scan/zero-on-error: Scan does not return exactly one value. The
// remedy is to return one value from Scan.
//
// Shelf.Keyed/cancel: Keyed cannot fail. The remedy is to give Keyed an
// error result.

// The unclaimed classifications of Shelf:
//
// The shape:writer of Shelf.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

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
	Put       shelfPut
	GetMany   shelfGetMany
	All       shelfAll
	Items     shelfItems
	Scan      shelfScan
	Keyed     shelfKeyed
	Snapshot  shelfSnapshot
	Len       shelfLen
	At        shelfAt
	AtOffset  shelfAtOffset
	Slot      shelfSlot
	Get       shelfGet
	TracedGet shelfTracedGet
	Err       shelfErr
}{
	Put:       shelfPut{"Shelf.Put"},
	GetMany:   shelfGetMany{"Shelf.GetMany"},
	All:       shelfAll{"Shelf.All"},
	Items:     shelfItems{"Shelf.Items"},
	Scan:      shelfScan{"Shelf.Scan"},
	Keyed:     shelfKeyed{"Shelf.Keyed"},
	Snapshot:  shelfSnapshot{"Shelf.Snapshot"},
	Len:       shelfLen{"Shelf.Len"},
	At:        shelfAt{"Shelf.At"},
	AtOffset:  shelfAtOffset{"Shelf.AtOffset"},
	Slot:      shelfSlot{"Shelf.Slot"},
	Get:       shelfGet{"Shelf.Get"},
	TracedGet: shelfTracedGet{"Shelf.TracedGet"},
	Err:       shelfErr{"Shelf.Err"},
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

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c shelfPut) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c shelfPut) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c shelfPut) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// shelfGetMany is the name of Shelf.GetMany, with the identities of its
// checks.
type shelfGetMany struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that GetMany returns without a
// panic for derived inputs.
func (c shelfGetMany) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that GetMany returns
// without a panic for derived inputs.
func (c shelfGetMany) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that GetMany returns an error
// for a cancelled context.
func (c shelfGetMany) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that GetMany returns
// an error for a cancelled context.
func (c shelfGetMany) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// ReadsEachKey returns the identity of the check that GetMany returns the
// value that Put stored under each key.
func (c shelfGetMany) ReadsEachKey() conform.ID[Shelf] {
	return c.Check("reads-each-key")
}

// ReadsEachKeyGenerated returns the identity of the check that GetMany
// returns the value that Put stored under each key.
func (c shelfGetMany) ReadsEachKeyGenerated() conform.ID[Shelf] {
	return c.Check("reads-each-key#generated")
}

// shelfAll is the name of Shelf.All, with the identities of its checks.
type shelfAll struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that All returns without a panic
// for derived inputs.
func (c shelfAll) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// StreamsWhatWasStored returns the identity of the check that All yields
// each value that Put stored.
func (c shelfAll) StreamsWhatWasStored() conform.ID[Shelf] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that All
// yields each value that Put stored.
func (c shelfAll) StreamsWhatWasStoredGenerated() conform.ID[Shelf] {
	return c.Check("streams-what-was-stored#generated")
}

// shelfItems is the name of Shelf.Items, with the identities of its
// checks.
type shelfItems struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Items returns without a
// panic for derived inputs.
func (c shelfItems) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// StreamsWhatWasStored returns the identity of the check that Items yields
// each value that Put stored.
func (c shelfItems) StreamsWhatWasStored() conform.ID[Shelf] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Items yields each value that Put stored.
func (c shelfItems) StreamsWhatWasStoredGenerated() conform.ID[Shelf] {
	return c.Check("streams-what-was-stored#generated")
}

// shelfScan is the name of Shelf.Scan, with the identities of its checks.
type shelfScan struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Scan returns without a
// panic for derived inputs.
func (c shelfScan) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Scan returns an error for
// a cancelled context.
func (c shelfScan) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// OverMatch returns the identity of the check that Scan gains at least the
// elements that Put wrote between two of its calls.
func (c shelfScan) OverMatch() conform.ID[Shelf] {
	return c.Check("over-match")
}

// OverMatchGenerated returns the identity of the check that Scan gains at
// least the elements that Put wrote between two of its calls.
func (c shelfScan) OverMatchGenerated() conform.ID[Shelf] {
	return c.Check("over-match#generated")
}

// DrainsKeepWhatWentIn returns the identity of the check that Scan yields
// the elements of the call before it and each element that Put wrote
// since, over any sequence of calls.
func (c shelfScan) DrainsKeepWhatWentIn() conform.ID[Shelf] {
	return c.Check("drains-keep-what-went-in")
}

// shelfKeyed is the name of Shelf.Keyed, with the identities of its
// checks.
type shelfKeyed struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Keyed returns without a
// panic for derived inputs.
func (c shelfKeyed) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// StreamsWhatWasStored returns the identity of the check that Keyed yields
// each value that Put stored.
func (c shelfKeyed) StreamsWhatWasStored() conform.ID[Shelf] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Keyed yields each value that Put stored.
func (c shelfKeyed) StreamsWhatWasStoredGenerated() conform.ID[Shelf] {
	return c.Check("streams-what-was-stored#generated")
}

// shelfSnapshot is the name of Shelf.Snapshot, with the identities of its
// checks.
type shelfSnapshot struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Snapshot returns without a
// panic for derived inputs.
func (c shelfSnapshot) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Snapshot returns an error
// for a cancelled context.
func (c shelfSnapshot) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Snapshot returns the
// zero value beside an error.
func (c shelfSnapshot) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// OverMatch returns the identity of the check that Snapshot gains at least
// the elements that Put wrote between two of its calls.
func (c shelfSnapshot) OverMatch() conform.ID[Shelf] {
	return c.Check("over-match")
}

// OverMatchGenerated returns the identity of the check that Snapshot gains
// at least the elements that Put wrote between two of its calls.
func (c shelfSnapshot) OverMatchGenerated() conform.ID[Shelf] {
	return c.Check("over-match#generated")
}

// DrainsKeepWhatWentIn returns the identity of the check that Snapshot
// yields the elements of the call before it and each element that Put
// wrote since, over any sequence of calls.
func (c shelfSnapshot) DrainsKeepWhatWentIn() conform.ID[Shelf] {
	return c.Check("drains-keep-what-went-in")
}

// shelfLen is the name of Shelf.Len, with the identities of its checks.
type shelfLen struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Len returns without a panic
// for derived inputs.
func (c shelfLen) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Len returns an error for a
// cancelled context.
func (c shelfLen) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Len returns the zero
// value beside an error.
func (c shelfLen) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// shelfAt is the name of Shelf.At, with the identities of its checks.
type shelfAt struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that At returns without a panic
// for derived inputs.
func (c shelfAt) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that At returns without
// a panic for derived inputs.
func (c shelfAt) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that At returns an error for a
// cancelled context.
func (c shelfAt) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that At returns an
// error for a cancelled context.
func (c shelfAt) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that At returns the zero
// value beside an error.
func (c shelfAt) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that At returns
// the zero value beside an error.
func (c shelfAt) ZeroOnErrorGenerated() conform.ID[Shelf] {
	return c.Check("zero-on-error#generated")
}

// RespectsItsBounds returns the identity of the check that At admits each
// position below the size that Len reports and refuses the size.
func (c shelfAt) RespectsItsBounds() conform.ID[Shelf] {
	return c.Check("respects-its-bounds")
}

// shelfAtOffset is the name of Shelf.AtOffset, with the identities of its
// checks.
type shelfAtOffset struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that AtOffset returns without a
// panic for derived inputs.
func (c shelfAtOffset) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that AtOffset returns
// without a panic for derived inputs.
func (c shelfAtOffset) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that AtOffset returns an error
// for a cancelled context.
func (c shelfAtOffset) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that AtOffset returns
// an error for a cancelled context.
func (c shelfAtOffset) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that AtOffset returns the
// zero value beside an error.
func (c shelfAtOffset) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that AtOffset
// returns the zero value beside an error.
func (c shelfAtOffset) ZeroOnErrorGenerated() conform.ID[Shelf] {
	return c.Check("zero-on-error#generated")
}

// RespectsItsBounds returns the identity of the check that AtOffset admits
// each position below the size that Len reports and refuses the size.
func (c shelfAtOffset) RespectsItsBounds() conform.ID[Shelf] {
	return c.Check("respects-its-bounds")
}

// shelfSlot is the name of Shelf.Slot, with the identities of its checks.
type shelfSlot struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Slot returns without a
// panic for derived inputs.
func (c shelfSlot) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Slot returns
// without a panic for derived inputs.
func (c shelfSlot) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Slot returns an error for
// a cancelled context.
func (c shelfSlot) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Slot returns an
// error for a cancelled context.
func (c shelfSlot) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// RespectsItsBounds returns the identity of the check that Slot admits
// each position below the size that Len reports and refuses the size.
func (c shelfSlot) RespectsItsBounds() conform.ID[Shelf] {
	return c.Check("respects-its-bounds")
}

// RespectsItsBoundsGenerated returns the identity of the check that Slot
// admits each position below the size that Len reports and refuses the
// size.
func (c shelfSlot) RespectsItsBoundsGenerated() conform.ID[Shelf] {
	return c.Check("respects-its-bounds#generated")
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

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c shelfGet) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c shelfGet) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c shelfGet) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c shelfGet) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c shelfGet) ZeroOnErrorGenerated() conform.ID[Shelf] {
	return c.Check("zero-on-error#generated")
}

// shelfTracedGet is the name of Shelf.TracedGet, with the identities of
// its checks.
type shelfTracedGet struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that TracedGet returns without a
// panic for derived inputs.
func (c shelfTracedGet) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that TracedGet returns
// without a panic for derived inputs.
func (c shelfTracedGet) SmokeGenerated() conform.ID[Shelf] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that TracedGet returns an error
// for a cancelled context.
func (c shelfTracedGet) Cancel() conform.ID[Shelf] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that TracedGet returns
// an error for a cancelled context.
func (c shelfTracedGet) CancelGenerated() conform.ID[Shelf] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that TracedGet returns the
// zero value beside an error.
func (c shelfTracedGet) ZeroOnError() conform.ID[Shelf] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that TracedGet
// returns the zero value beside an error.
func (c shelfTracedGet) ZeroOnErrorGenerated() conform.ID[Shelf] {
	return c.Check("zero-on-error#generated")
}

// AgreesWithItsSibling returns the identity of the check that TracedGet
// returns what Get returns for one input.
func (c shelfTracedGet) AgreesWithItsSibling() conform.ID[Shelf] {
	return c.Check("agrees-with-its-sibling")
}

// AgreesWithItsSiblingGenerated returns the identity of the check that
// TracedGet returns what Get returns for one input.
func (c shelfTracedGet) AgreesWithItsSiblingGenerated() conform.ID[Shelf] {
	return c.Check("agrees-with-its-sibling#generated")
}

// shelfErr is the name of Shelf.Err, with the identities of its checks.
type shelfErr struct {
	conform.Callable[Shelf]
}

// Smoke returns the identity of the check that Err returns without a panic
// for derived inputs.
func (c shelfErr) Smoke() conform.ID[Shelf] {
	return c.Check("smoke")
}

// FreshReportsNoPoison returns the identity of the check that Err reports
// no error for a fresh receiver.
func (c shelfErr) FreshReportsNoPoison() conform.ID[Shelf] {
	return c.Check("fresh-reports-no-poison")
}

// ShelfSuite returns the checks of Shelf, and the calls and the budgets of
// its callables.
func ShelfSuite() conform.Suite[Shelf] {
	return conform.Suite[Shelf]{
		Declaration: "Shelf",
		Package:     "example.com/shop/classified",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Shelf]{
			{
				ID:          ShelfChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: shelfPutPanics,
			},
			{
				ID:          ShelfChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: shelfPutPanics,
			},
			{
				ID:          ShelfChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: shelfPutIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: shelfPutIgnoresCancel,
			},
			{
				ID:          ShelfChecks.GetMany.Smoke(),
				Statement:   "GetMany returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().GetMany(ctx, string("keys-a"))
						return err
					})
				},
				Defect: shelfGetManyPanics,
			},
			{
				ID:          ShelfChecks.GetMany.SmokeGenerated(),
				Statement:   "GetMany returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKeys := c.Draw(prop.Of[string](), "keys")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().GetMany(ctx, drawnKeys)
						return err
					})
				},
				Defect: shelfGetManyPanics,
			},
			{
				ID:          ShelfChecks.GetMany.Cancel(),
				Statement:   "GetMany returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().GetMany(ctx, string("keys-a"))
						return err
					})
				},
				Defect: shelfGetManyIgnoresCancel,
			},
			{
				ID:          ShelfChecks.GetMany.CancelGenerated(),
				Statement:   "GetMany returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKeys := c.Draw(prop.Of[string](), "keys")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().GetMany(ctx, drawnKeys)
						return err
					})
				},
				Defect: shelfGetManyIgnoresCancel,
			},
			{
				ID:          ShelfChecks.GetMany.ReadsEachKey(),
				Statement:   "GetMany returns the value that Put stored under each key",
				Claim:       spec.ReadsEachKey,
				Probe:       spec.ReadsWhatWasStored,
				Reason:      spec.LostAKey,
				Law:         spec.BatchReadReturnsEachStoredValue,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ReadsWhatWasStored(h, func(ctx context.Context) ([]Item, error) {
						return h.Subject().GetMany(ctx, string("Key-a"), string("Key-b"))
					}, func() func(ctx context.Context) (Item, error) {
						calls := 0
						return func(ctx context.Context) (Item, error) {
							calls++
							w := Item{Key: "Key-a"}
							w.Key = string("Key-a")
							if calls > 1 {
								w.Key = string("Key-b")
							}
							err := h.Subject().Put(ctx, w)
							return w, err
						}
					}())
				},
				Defect: shelfGetManyDropsTheLastKey,
			},
			{
				ID:          ShelfChecks.GetMany.ReadsEachKeyGenerated(),
				Statement:   "GetMany returns the value that Put stored under each key",
				Claim:       spec.ReadsEachKey,
				Probe:       spec.ReadsWhatWasStored,
				Reason:      spec.LostAKey,
				Law:         spec.BatchReadReturnsEachStoredValue,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnItKey := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"it.key",
					)
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.ReadsWhatWasStored(h, func(ctx context.Context) ([]Item, error) {
						return h.Subject().GetMany(ctx, drawnItKey[0], drawnItKey[1])
					}, func() func(ctx context.Context) (Item, error) {
						calls := 0
						return func(ctx context.Context) (Item, error) {
							calls++
							w := drawnIt
							w.Key = drawnItKey[0]
							if calls > 1 {
								w.Key = drawnItKey[1]
							}
							err := h.Subject().Put(ctx, w)
							return w, err
						}
					}())
				},
				Defect: shelfGetManyDropsTheLastKey,
			},
			{
				ID:          ShelfChecks.All.Smoke(),
				Statement:   "All returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().All(ctx)
						return nil
					})
				},
				Defect: shelfAllPanics,
			},
			{
				ID:          ShelfChecks.All.StreamsWhatWasStored(),
				Statement:   "All yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						r0 := h.Subject().All(ctx)
						var drained []Item
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Item{Key: "Key-a"}
							if calls > 1 {
								w = Item{Key: "Key-b"}
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfAllDrainsNothing,
			},
			{
				ID:          ShelfChecks.All.StreamsWhatWasStoredGenerated(),
				Statement:   "All yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnIt := c.Draw(
						prop.List(prop.Of[Item](), prop.Unique(), prop.MinSize(2), prop.MaxSize(2)),
						"it",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						r0 := h.Subject().All(ctx)
						var drained []Item
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnIt[0]
							if calls > 1 {
								w = drawnIt[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfAllDrainsNothing,
			},
			{
				ID:          ShelfChecks.Items.Smoke(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Items(ctx)
						return nil
					})
				},
				Defect: shelfItemsPanics,
			},
			{
				ID:          ShelfChecks.Items.StreamsWhatWasStored(),
				Statement:   "Items yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						r0 := h.Subject().Items(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Item{Key: "Key-a"}
							if calls > 1 {
								w = Item{Key: "Key-b"}
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfItemsDrainsNothing,
			},
			{
				ID:          ShelfChecks.Items.StreamsWhatWasStoredGenerated(),
				Statement:   "Items yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnIt := c.Draw(
						prop.List(prop.Of[Item](), prop.Unique(), prop.MinSize(2), prop.MaxSize(2)),
						"it",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						r0 := h.Subject().Items(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnIt[0]
							if calls > 1 {
								w = drawnIt[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfItemsDrainsNothing,
			},
			{
				ID:          ShelfChecks.Scan.Smoke(),
				Statement:   "Scan returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Scan(ctx)
						return err
					})
				},
				Defect: shelfScanPanics,
			},
			{
				ID:          ShelfChecks.Scan.Cancel(),
				Statement:   "Scan returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Scan(ctx)
						return err
					})
				},
				Defect: shelfScanIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Scan.OverMatch(),
				Statement:   "Scan gains at least the elements that Put wrote between two of its calls",
				Claim:       spec.OverMatch,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						r0, err := h.Subject().Scan(ctx)
						if err != nil {
							return nil, err
						}
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Item{Key: "Key-a"}
							if calls > 1 {
								w = Item{Key: "Key-b"}
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfScanDrainsNothing,
			},
			{
				ID:          ShelfChecks.Scan.OverMatchGenerated(),
				Statement:   "Scan gains at least the elements that Put wrote between two of its calls",
				Claim:       spec.OverMatch,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnIt := c.Draw(
						prop.List(prop.Of[Item](), prop.Unique(), prop.MinSize(2), prop.MaxSize(2)),
						"it",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						r0, err := h.Subject().Scan(ctx)
						if err != nil {
							return nil, err
						}
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnIt[0]
							if calls > 1 {
								w = drawnIt[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfScanDrainsNothing,
			},
			{
				ID:          ShelfChecks.Scan.DrainsKeepWhatWentIn(),
				Statement:   "Scan yields the elements of the call before it and each element that Put wrote since, over any sequence of calls",
				Claim:       spec.DrainsKeepWhatWentIn,
				Probe:       spec.DrainsContainTheWrites,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Shelf]{ShelfChecks.Scan.OverMatch()},
				Property: func(h *ShelfHandle, c *prop.Case) {
					probe.DrainsContainTheWrites(
						h,
						c,
						shelfActions,
						"Scan",
						"Put",
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: shelfScanDropsAfterThree,
			},
			{
				ID:          ShelfChecks.Keyed.Smoke(),
				Statement:   "Keyed returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Keyed(ctx)
						return nil
					})
				},
				Defect: shelfKeyedPanics,
			},
			{
				ID:          ShelfChecks.Keyed.StreamsWhatWasStored(),
				Statement:   "Keyed yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, Item], error) {
							r0 := h.Subject().Keyed(ctx)
							var drained []probe.Pair[string, Item]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, Item]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := Item{Key: "Key-a"}
								if calls > 1 {
									w = Item{Key: "Key-b"}
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: shelfKeyedDrainsNothing,
			},
			{
				ID:          ShelfChecks.Keyed.StreamsWhatWasStoredGenerated(),
				Statement:   "Keyed yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnIt := c.Draw(
						prop.List(prop.Of[Item](), prop.Unique(), prop.MinSize(2), prop.MaxSize(2)),
						"it",
					)
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, Item], error) {
							r0 := h.Subject().Keyed(ctx)
							var drained []probe.Pair[string, Item]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, Item]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := drawnIt[0]
								if calls > 1 {
									w = drawnIt[1]
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: shelfKeyedDrainsNothing,
			},
			{
				ID:          ShelfChecks.Snapshot.Smoke(),
				Statement:   "Snapshot returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Snapshot(ctx)
						return err
					})
				},
				Defect: shelfSnapshotPanics,
			},
			{
				ID:          ShelfChecks.Snapshot.Cancel(),
				Statement:   "Snapshot returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Snapshot(ctx)
						return err
					})
				},
				Defect: shelfSnapshotIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Snapshot.ZeroOnError(),
				Statement:   "Snapshot returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]Item, error) {
						return h.Subject().Snapshot(ctx)
					}, []Item(nil))
				},
				Defect: shelfSnapshotValueBeside,
			},
			{
				ID:          ShelfChecks.Snapshot.OverMatch(),
				Statement:   "Snapshot gains at least the elements that Put wrote between two of its calls",
				Claim:       spec.OverMatch,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						return h.Subject().Snapshot(ctx)
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Item{Key: "Key-a"}
							if calls > 1 {
								w = Item{Key: "Key-b"}
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfSnapshotDrainsNothing,
			},
			{
				ID:          ShelfChecks.Snapshot.OverMatchGenerated(),
				Statement:   "Snapshot gains at least the elements that Put wrote between two of its calls",
				Claim:       spec.OverMatch,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnIt := c.Draw(
						prop.List(prop.Of[Item](), prop.Unique(), prop.MinSize(2), prop.MaxSize(2)),
						"it",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]Item, error) {
						return h.Subject().Snapshot(ctx)
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnIt[0]
							if calls > 1 {
								w = drawnIt[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: shelfSnapshotDrainsNothing,
			},
			{
				ID:          ShelfChecks.Snapshot.DrainsKeepWhatWentIn(),
				Statement:   "Snapshot yields the elements of the call before it and each element that Put wrote since, over any sequence of calls",
				Claim:       spec.DrainsKeepWhatWentIn,
				Probe:       spec.DrainsContainTheWrites,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Shelf]{ShelfChecks.Snapshot.OverMatch()},
				Property: func(h *ShelfHandle, c *prop.Case) {
					probe.DrainsContainTheWrites(
						h,
						c,
						shelfActions,
						"Snapshot",
						"Put",
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: shelfSnapshotDropsAfterThree,
			},
			{
				ID:          ShelfChecks.Len.Smoke(),
				Statement:   "Len returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Len(ctx)
						return err
					})
				},
				Defect: shelfLenPanics,
			},
			{
				ID:          ShelfChecks.Len.Cancel(),
				Statement:   "Len returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Len(ctx)
						return err
					})
				},
				Defect: shelfLenIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Len.ZeroOnError(),
				Statement:   "Len returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Len(ctx)
					}, int(0))
				},
				Defect: shelfLenValueBeside,
			},
			{
				ID:          ShelfChecks.At.Smoke(),
				Statement:   "At returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().At(ctx, int(42))
						return err
					})
				},
				Defect: shelfAtPanics,
			},
			{
				ID:          ShelfChecks.At.SmokeGenerated(),
				Statement:   "At returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnI := c.Draw(prop.Of[int](), "i")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().At(ctx, drawnI)
						return err
					})
				},
				Defect: shelfAtPanics,
			},
			{
				ID:          ShelfChecks.At.Cancel(),
				Statement:   "At returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().At(ctx, int(42))
						return err
					})
				},
				Defect: shelfAtIgnoresCancel,
			},
			{
				ID:          ShelfChecks.At.CancelGenerated(),
				Statement:   "At returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnI := c.Draw(prop.Of[int](), "i")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().At(ctx, drawnI)
						return err
					})
				},
				Defect: shelfAtIgnoresCancel,
			},
			{
				ID:          ShelfChecks.At.ZeroOnError(),
				Statement:   "At returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().At(ctx, int(42))
					}, Item{})
				},
				Defect: shelfAtValueBeside,
			},
			{
				ID:          ShelfChecks.At.ZeroOnErrorGenerated(),
				Statement:   "At returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnI := c.Draw(prop.Of[int](), "i")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().At(ctx, drawnI)
					}, Item{})
				},
				Defect: shelfAtValueBeside,
			},
			{
				ID:          ShelfChecks.At.RespectsItsBounds(),
				Statement:   "At admits each position below the size that Len reports and refuses the size",
				Claim:       spec.RespectsItsBounds,
				Probe:       spec.BoundedByTheSize,
				Reason:      spec.OutOfBounds,
				Law:         spec.PositionsAreBelowTheSize,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.BoundedByTheSize(h, func(ctx context.Context, position int) error {
						_, err := h.Subject().At(ctx, position)
						return err
					}, func(ctx context.Context) (int, error) {
						return h.Subject().Len(ctx)
					})
				},
				Defect: shelfAtAdmitsTheSize,
			},
			{
				ID:          ShelfChecks.AtOffset.Smoke(),
				Statement:   "AtOffset returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().AtOffset(ctx, int64(42))
						return err
					})
				},
				Defect: shelfAtOffsetPanics,
			},
			{
				ID:          ShelfChecks.AtOffset.SmokeGenerated(),
				Statement:   "AtOffset returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnI := c.Draw(prop.Of[int64](), "i")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().AtOffset(ctx, drawnI)
						return err
					})
				},
				Defect: shelfAtOffsetPanics,
			},
			{
				ID:          ShelfChecks.AtOffset.Cancel(),
				Statement:   "AtOffset returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().AtOffset(ctx, int64(42))
						return err
					})
				},
				Defect: shelfAtOffsetIgnoresCancel,
			},
			{
				ID:          ShelfChecks.AtOffset.CancelGenerated(),
				Statement:   "AtOffset returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnI := c.Draw(prop.Of[int64](), "i")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().AtOffset(ctx, drawnI)
						return err
					})
				},
				Defect: shelfAtOffsetIgnoresCancel,
			},
			{
				ID:          ShelfChecks.AtOffset.ZeroOnError(),
				Statement:   "AtOffset returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().AtOffset(ctx, int64(42))
					}, Item{})
				},
				Defect: shelfAtOffsetValueBeside,
			},
			{
				ID:          ShelfChecks.AtOffset.ZeroOnErrorGenerated(),
				Statement:   "AtOffset returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnI := c.Draw(prop.Of[int64](), "i")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().AtOffset(ctx, drawnI)
					}, Item{})
				},
				Defect: shelfAtOffsetValueBeside,
			},
			{
				ID:          ShelfChecks.AtOffset.RespectsItsBounds(),
				Statement:   "AtOffset admits each position below the size that Len reports and refuses the size",
				Claim:       spec.RespectsItsBounds,
				Probe:       spec.BoundedByTheSize,
				Reason:      spec.OutOfBounds,
				Law:         spec.PositionsAreBelowTheSize,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.BoundedByTheSize(h, func(ctx context.Context, position int) error {
						_, err := h.Subject().AtOffset(ctx, int64(position))
						return err
					}, func(ctx context.Context) (int, error) {
						return h.Subject().Len(ctx)
					})
				},
				Defect: shelfAtOffsetAdmitsTheSize,
			},
			{
				ID:          ShelfChecks.Slot.Smoke(),
				Statement:   "Slot returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Slot(ctx, string("name-a"), uint(42))
						return err
					})
				},
				Defect: shelfSlotPanics,
			},
			{
				ID:          ShelfChecks.Slot.SmokeGenerated(),
				Statement:   "Slot returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnName := c.Draw(prop.Of[string](), "name")
					drawnI := c.Draw(prop.Of[uint](), "i")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Slot(ctx, drawnName, drawnI)
						return err
					})
				},
				Defect: shelfSlotPanics,
			},
			{
				ID:          ShelfChecks.Slot.Cancel(),
				Statement:   "Slot returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Slot(ctx, string("name-a"), uint(42))
						return err
					})
				},
				Defect: shelfSlotIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Slot.CancelGenerated(),
				Statement:   "Slot returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnName := c.Draw(prop.Of[string](), "name")
					drawnI := c.Draw(prop.Of[uint](), "i")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Slot(ctx, drawnName, drawnI)
						return err
					})
				},
				Defect: shelfSlotIgnoresCancel,
			},
			{
				ID:          ShelfChecks.Slot.RespectsItsBounds(),
				Statement:   "Slot admits each position below the size that Len reports and refuses the size",
				Claim:       spec.RespectsItsBounds,
				Probe:       spec.BoundedByTheSize,
				Reason:      spec.OutOfBounds,
				Law:         spec.PositionsAreBelowTheSize,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.BoundedByTheSize(h, func(ctx context.Context, position int) error {
						_, err := h.Subject().Slot(ctx, string("name-a"), uint(position))
						return err
					}, func(ctx context.Context) (int, error) {
						return h.Subject().Len(ctx)
					})
				},
				Defect: shelfSlotAdmitsTheSize,
			},
			{
				ID:          ShelfChecks.Slot.RespectsItsBoundsGenerated(),
				Statement:   "Slot admits each position below the size that Len reports and refuses the size",
				Claim:       spec.RespectsItsBounds,
				Probe:       spec.BoundedByTheSize,
				Reason:      spec.OutOfBounds,
				Law:         spec.PositionsAreBelowTheSize,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnName := c.Draw(prop.Of[string](), "name")
					probe.BoundedByTheSize(h, func(ctx context.Context, position int) error {
						_, err := h.Subject().Slot(ctx, drawnName, uint(position))
						return err
					}, func(ctx context.Context) (int, error) {
						return h.Subject().Len(ctx)
					})
				},
				Defect: shelfSlotAdmitsTheSize,
			},
			{
				ID:          ShelfChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
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
				ID:          ShelfChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
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
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
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
				ID:          ShelfChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
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
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: shelfGetValueBeside,
			},
			{
				ID:          ShelfChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: shelfGetValueBeside,
			},
			{
				ID:          ShelfChecks.TracedGet.Smoke(),
				Statement:   "TracedGet returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().TracedGet(ctx, string("key-a"))
						return err
					})
				},
				Defect: shelfTracedGetPanics,
			},
			{
				ID:          ShelfChecks.TracedGet.SmokeGenerated(),
				Statement:   "TracedGet returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().TracedGet(ctx, drawnKey)
						return err
					})
				},
				Defect: shelfTracedGetPanics,
			},
			{
				ID:          ShelfChecks.TracedGet.Cancel(),
				Statement:   "TracedGet returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().TracedGet(ctx, string("key-a"))
						return err
					})
				},
				Defect: shelfTracedGetIgnoresCancel,
			},
			{
				ID:          ShelfChecks.TracedGet.CancelGenerated(),
				Statement:   "TracedGet returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().TracedGet(ctx, drawnKey)
						return err
					})
				},
				Defect: shelfTracedGetIgnoresCancel,
			},
			{
				ID:          ShelfChecks.TracedGet.ZeroOnError(),
				Statement:   "TracedGet returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().TracedGet(ctx, string("key-a"))
					}, Item{})
				},
				Defect: shelfTracedGetValueBeside,
			},
			{
				ID:          ShelfChecks.TracedGet.ZeroOnErrorGenerated(),
				Statement:   "TracedGet returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().TracedGet(ctx, drawnKey)
					}, Item{})
				},
				Defect: shelfTracedGetValueBeside,
			},
			{
				ID:          ShelfChecks.TracedGet.AgreesWithItsSibling(),
				Statement:   "TracedGet returns what Get returns for one input",
				Claim:       spec.AgreesWithItsSibling,
				Probe:       spec.AgreesWithTheSibling,
				Reason:      spec.DivergedFromTheSibling,
				Law:         spec.WrapperReturnsWhatItsSiblingReturns,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.AgreesWithTheSibling(h, func(ctx context.Context) (Item, error) {
						return h.Subject().TracedGet(ctx, string("key-a"))
					}, func(ctx context.Context) (Item, error) {
						return h.New().Get(ctx, string("key-a"))
					})
				},
				Defect: shelfTracedGetStopsDelegating,
			},
			{
				ID:          ShelfChecks.TracedGet.AgreesWithItsSiblingGenerated(),
				Statement:   "TracedGet returns what Get returns for one input",
				Claim:       spec.AgreesWithItsSibling,
				Probe:       spec.AgreesWithTheSibling,
				Reason:      spec.DivergedFromTheSibling,
				Law:         spec.WrapperReturnsWhatItsSiblingReturns,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ShelfHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.AgreesWithTheSibling(h, func(ctx context.Context) (Item, error) {
						return h.Subject().TracedGet(ctx, drawnKey)
					}, func(ctx context.Context) (Item, error) {
						return h.New().Get(ctx, drawnKey)
					})
				},
				Defect: shelfTracedGetStopsDelegating,
			},
			{
				ID:          ShelfChecks.Err.Smoke(),
				Statement:   "Err returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ShelfHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Err()
					})
				},
				Defect: shelfErrPanics,
			},
			{
				ID:          ShelfChecks.Err.FreshReportsNoPoison(),
				Statement:   "Err reports no error for a fresh receiver",
				Claim:       spec.FreshReportsNoPoison,
				Probe:       spec.ReportsNoError,
				Reason:      spec.FalsePoison,
				Law:         spec.FreshReceiverIsHealthy,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ShelfHandle) {
					probe.ReportsNoError(h, func(ctx context.Context) error {
						return h.Subject().Err()
					})
				},
				Defect: shelfErrReportsPoison,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Shelf.All",
				Gate:     spec.Failure,
				Why:      "All cannot fail",
				Remedy:   "give All an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Shelf.Items",
				Gate:     spec.Failure,
				Why:      "Items cannot fail",
				Remedy:   "give Items an error result",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Shelf.Scan",
				Gate:     spec.OneValue,
				Why:      "Scan does not return exactly one value",
				Remedy:   "return one value from Scan",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Shelf.Keyed",
				Gate:     spec.Failure,
				Why:      "Keyed cannot fail",
				Remedy:   "give Keyed an error result",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Shelf.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Shelf]{
			{
				Calls:    shelfPutCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Put.Callable,
			},
			{
				Calls:    shelfGetManyCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.GetMany.Callable,
			},
			{
				Calls:    shelfAllCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.All.Callable,
			},
			{
				Calls:    shelfItemsCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Items.Callable,
			},
			{
				Calls:    shelfScanCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Scan.Callable,
			},
			{
				Calls:    shelfKeyedCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Keyed.Callable,
			},
			{
				Calls:    shelfSnapshotCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Snapshot.Callable,
			},
			{
				Calls:    shelfLenCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Len.Callable,
			},
			{
				Calls:    shelfAtCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.At.Callable,
			},
			{
				Calls:    shelfAtOffsetCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.AtOffset.Callable,
			},
			{
				Calls:    shelfSlotCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Slot.Callable,
			},
			{
				Calls:    shelfGetCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Get.Callable,
			},
			{
				Calls:    shelfTracedGetCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.TracedGet.Callable,
			},
			{
				Calls:    shelfErrCalls,
				Over:     newShelfOver,
				Callable: ShelfChecks.Err.Callable,
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

// shelfPutCalls generates the calls of Shelf.Put with the inputs that each
// case draws.
var shelfPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Shelf) error {
		return subject.Put(ctx, it)
	}
})

// shelfGetManyCalls generates the calls of Shelf.GetMany with the inputs
// that each case draws.
var shelfGetManyCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	keys := c.Draw(prop.List(prop.Of[string]()), "keys")
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.GetMany(ctx, keys...)
		return err
	}
})

// shelfAllCalls generates the calls of Shelf.All with the inputs that each
// case draws.
var shelfAllCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(ctx context.Context, subject Shelf) error {
		_ = subject.All(ctx)
		return nil
	}
})

// shelfItemsCalls generates the calls of Shelf.Items with the inputs that
// each case draws.
var shelfItemsCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(ctx context.Context, subject Shelf) error {
		_ = subject.Items(ctx)
		return nil
	}
})

// shelfScanCalls generates the calls of Shelf.Scan with the inputs that
// each case draws.
var shelfScanCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.Scan(ctx)
		return err
	}
})

// shelfKeyedCalls generates the calls of Shelf.Keyed with the inputs that
// each case draws.
var shelfKeyedCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(ctx context.Context, subject Shelf) error {
		_ = subject.Keyed(ctx)
		return nil
	}
})

// shelfSnapshotCalls generates the calls of Shelf.Snapshot with the inputs
// that each case draws.
var shelfSnapshotCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.Snapshot(ctx)
		return err
	}
})

// shelfLenCalls generates the calls of Shelf.Len with the inputs that each
// case draws.
var shelfLenCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.Len(ctx)
		return err
	}
})

// shelfAtCalls generates the calls of Shelf.At with the inputs that each
// case draws.
var shelfAtCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	i := c.Draw(prop.Of[int](), "i")
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.At(ctx, i)
		return err
	}
})

// shelfAtOffsetCalls generates the calls of Shelf.AtOffset with the inputs
// that each case draws.
var shelfAtOffsetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	i := c.Draw(prop.Of[int64](), "i")
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.AtOffset(ctx, i)
		return err
	}
})

// shelfSlotCalls generates the calls of Shelf.Slot with the inputs that
// each case draws.
var shelfSlotCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	name := c.Draw(prop.Of[string](), "name")
	i := c.Draw(prop.Of[uint](), "i")
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.Slot(ctx, name, i)
		return err
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

// shelfTracedGetCalls generates the calls of Shelf.TracedGet with the
// inputs that each case draws.
var shelfTracedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Shelf) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Shelf) error {
		_, err := subject.TracedGet(ctx, key)
		return err
	}
})

// shelfErrCalls generates the calls of Shelf.Err with the inputs that each
// case draws.
var shelfErrCalls = prop.Composite(func(*prop.Case) func(context.Context, Shelf) error {
	return func(_ context.Context, subject Shelf) error {
		return subject.Err()
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
func (o shelfOver) Put(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// GetMany calls spend and returns its zero results.
func (o shelfOver) GetMany(context.Context, ...string) (r0 []Item, r1 error) {
	o.spend()
	return r0, r1
}

// All calls spend and returns its zero results.
func (o shelfOver) All(context.Context) (r0 iter.Seq2[Item, error]) {
	o.spend()
	return r0
}

// Items calls spend and returns its zero results.
func (o shelfOver) Items(context.Context) (r0 Stream) {
	o.spend()
	return r0
}

// Scan calls spend and returns its zero results.
func (o shelfOver) Scan(context.Context) (r0 iter.Seq[Item], r1 error) {
	o.spend()
	return r0, r1
}

// Keyed calls spend and returns its zero results.
func (o shelfOver) Keyed(context.Context) (r0 iter.Seq2[string, Item]) {
	o.spend()
	return r0
}

// Snapshot calls spend and returns its zero results.
func (o shelfOver) Snapshot(context.Context) (r0 []Item, r1 error) {
	o.spend()
	return r0, r1
}

// Len calls spend and returns its zero results.
func (o shelfOver) Len(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// At calls spend and returns its zero results.
func (o shelfOver) At(context.Context, int) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// AtOffset calls spend and returns its zero results.
func (o shelfOver) AtOffset(context.Context, int64) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Slot calls spend and returns its zero results.
func (o shelfOver) Slot(context.Context, string, uint) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Get calls spend and returns its zero results.
func (o shelfOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// TracedGet calls spend and returns its zero results.
func (o shelfOver) TracedGet(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Err calls spend and returns its zero results.
func (o shelfOver) Err() (r0 error) {
	o.spend()
	return r0
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
func (*shelfPutPanicsDefect) Put(p0 context.Context, p1 Item) (r0 error) {
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
func (d *shelfPutIgnoresCancelDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.shelfOverEmbedded.Put(p0, p1)
}

// shelfPutIgnoresCancel returns the planted defect of a Shelf whose Put
// ignores the context over a fresh subject.
func shelfPutIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfPutIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetManyPanicsDefect is the planted defect of a Shelf whose GetMany
// panics.
type shelfGetManyPanicsDefect struct {
	shelfOverEmbedded
}

// GetMany is the override of the planted defect.
func (*shelfGetManyPanicsDefect) GetMany(p0 context.Context, p1 ...string) (r0 []Item, r1 error) {
	panic("the planted defect of a Shelf whose GetMany panics")
}

// shelfGetManyPanics returns the planted defect of a Shelf whose GetMany
// panics over a fresh subject.
func shelfGetManyPanics(h *ShelfHandle) Shelf {
	return &shelfGetManyPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetManyIgnoresCancelDefect is the planted defect of a Shelf whose
// GetMany ignores the context.
type shelfGetManyIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// GetMany is the override of the planted defect.
func (d *shelfGetManyIgnoresCancelDefect) GetMany(
	p0 context.Context,
	p1 ...string,
) (r0 []Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.GetMany(p0, p1...)
}

// shelfGetManyIgnoresCancel returns the planted defect of a Shelf whose
// GetMany ignores the context over a fresh subject.
func shelfGetManyIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfGetManyIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetManyDropsTheLastKeyDefect is the planted defect of a Shelf whose
// GetMany drops the value of its last key.
type shelfGetManyDropsTheLastKeyDefect struct {
	shelfOverEmbedded
}

// GetMany is the override of the planted defect.
func (d *shelfGetManyDropsTheLastKeyDefect) GetMany(
	p0 context.Context,
	p1 ...string,
) (r0 []Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.shelfOverEmbedded.GetMany(p0, p1...)
	if len(r0) > 0 {
		r0 = r0[:len(r0)-1]
	}
	return r0, r1
}

// shelfGetManyDropsTheLastKey returns the planted defect of a Shelf whose
// GetMany drops the value of its last key over a fresh subject.
func shelfGetManyDropsTheLastKey(h *ShelfHandle) Shelf {
	return &shelfGetManyDropsTheLastKeyDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAllPanicsDefect is the planted defect of a Shelf whose All panics.
type shelfAllPanicsDefect struct {
	shelfOverEmbedded
}

// All is the override of the planted defect.
func (*shelfAllPanicsDefect) All(p0 context.Context) (r0 iter.Seq2[Item, error]) {
	panic("the planted defect of a Shelf whose All panics")
}

// shelfAllPanics returns the planted defect of a Shelf whose All panics
// over a fresh subject.
func shelfAllPanics(h *ShelfHandle) Shelf {
	return &shelfAllPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAllDrainsNothingDefect is the planted defect of a Shelf whose All
// loses what Put stored.
type shelfAllDrainsNothingDefect struct {
	shelfOverEmbedded
}

// All is the override of the planted defect.
func (*shelfAllDrainsNothingDefect) All(p0 context.Context) (r0 iter.Seq2[Item, error]) {
	r0 = func(func(Item, error) bool) {}
	return r0
}

// shelfAllDrainsNothing returns the planted defect of a Shelf whose All
// loses what Put stored over a fresh subject.
func shelfAllDrainsNothing(h *ShelfHandle) Shelf {
	return &shelfAllDrainsNothingDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfItemsPanicsDefect is the planted defect of a Shelf whose Items
// panics.
type shelfItemsPanicsDefect struct {
	shelfOverEmbedded
}

// Items is the override of the planted defect.
func (*shelfItemsPanicsDefect) Items(p0 context.Context) (r0 Stream) {
	panic("the planted defect of a Shelf whose Items panics")
}

// shelfItemsPanics returns the planted defect of a Shelf whose Items
// panics over a fresh subject.
func shelfItemsPanics(h *ShelfHandle) Shelf {
	return &shelfItemsPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfItemsDrainsNothingDefect is the planted defect of a Shelf whose
// Items loses what Put stored.
type shelfItemsDrainsNothingDefect struct {
	shelfOverEmbedded
}

// Items is the override of the planted defect.
func (*shelfItemsDrainsNothingDefect) Items(p0 context.Context) (r0 Stream) {
	r0 = func(func(Item) bool) {}
	return r0
}

// shelfItemsDrainsNothing returns the planted defect of a Shelf whose
// Items loses what Put stored over a fresh subject.
func shelfItemsDrainsNothing(h *ShelfHandle) Shelf {
	return &shelfItemsDrainsNothingDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfScanPanicsDefect is the planted defect of a Shelf whose Scan
// panics.
type shelfScanPanicsDefect struct {
	shelfOverEmbedded
}

// Scan is the override of the planted defect.
func (*shelfScanPanicsDefect) Scan(p0 context.Context) (r0 iter.Seq[Item], r1 error) {
	panic("the planted defect of a Shelf whose Scan panics")
}

// shelfScanPanics returns the planted defect of a Shelf whose Scan panics
// over a fresh subject.
func shelfScanPanics(h *ShelfHandle) Shelf {
	return &shelfScanPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfScanIgnoresCancelDefect is the planted defect of a Shelf whose Scan
// ignores the context.
type shelfScanIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Scan is the override of the planted defect.
func (d *shelfScanIgnoresCancelDefect) Scan(p0 context.Context) (r0 iter.Seq[Item], r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.Scan(p0)
}

// shelfScanIgnoresCancel returns the planted defect of a Shelf whose Scan
// ignores the context over a fresh subject.
func shelfScanIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfScanIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfScanDrainsNothingDefect is the planted defect of a Shelf whose Scan
// loses what Put wrote.
type shelfScanDrainsNothingDefect struct {
	shelfOverEmbedded
}

// Scan is the override of the planted defect.
func (*shelfScanDrainsNothingDefect) Scan(p0 context.Context) (r0 iter.Seq[Item], r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = func(func(Item) bool) {}
	return r0, r1
}

// shelfScanDrainsNothing returns the planted defect of a Shelf whose Scan
// loses what Put wrote over a fresh subject.
func shelfScanDrainsNothing(h *ShelfHandle) Shelf {
	return &shelfScanDrainsNothingDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfScanDropsAfterThreeDefect is the planted defect of a Shelf whose
// Scan drops its last element from its fourth call on.
type shelfScanDropsAfterThreeDefect struct {
	shelfOverEmbedded
	calls int
}

// Scan is the override of the planted defect.
func (d *shelfScanDropsAfterThreeDefect) Scan(p0 context.Context) (r0 iter.Seq[Item], r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.shelfOverEmbedded.Scan(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = r0[:len(r0)-1]
	}
	return r0, r1
}

// shelfScanDropsAfterThree returns the planted defect of a Shelf whose
// Scan drops its last element from its fourth call on over a fresh
// subject.
func shelfScanDropsAfterThree(h *ShelfHandle) Shelf {
	return &shelfScanDropsAfterThreeDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfKeyedPanicsDefect is the planted defect of a Shelf whose Keyed
// panics.
type shelfKeyedPanicsDefect struct {
	shelfOverEmbedded
}

// Keyed is the override of the planted defect.
func (*shelfKeyedPanicsDefect) Keyed(p0 context.Context) (r0 iter.Seq2[string, Item]) {
	panic("the planted defect of a Shelf whose Keyed panics")
}

// shelfKeyedPanics returns the planted defect of a Shelf whose Keyed
// panics over a fresh subject.
func shelfKeyedPanics(h *ShelfHandle) Shelf {
	return &shelfKeyedPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfKeyedDrainsNothingDefect is the planted defect of a Shelf whose
// Keyed loses what Put stored.
type shelfKeyedDrainsNothingDefect struct {
	shelfOverEmbedded
}

// Keyed is the override of the planted defect.
func (*shelfKeyedDrainsNothingDefect) Keyed(p0 context.Context) (r0 iter.Seq2[string, Item]) {
	r0 = func(func(string, Item) bool) {}
	return r0
}

// shelfKeyedDrainsNothing returns the planted defect of a Shelf whose
// Keyed loses what Put stored over a fresh subject.
func shelfKeyedDrainsNothing(h *ShelfHandle) Shelf {
	return &shelfKeyedDrainsNothingDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSnapshotPanicsDefect is the planted defect of a Shelf whose
// Snapshot panics.
type shelfSnapshotPanicsDefect struct {
	shelfOverEmbedded
}

// Snapshot is the override of the planted defect.
func (*shelfSnapshotPanicsDefect) Snapshot(p0 context.Context) (r0 []Item, r1 error) {
	panic("the planted defect of a Shelf whose Snapshot panics")
}

// shelfSnapshotPanics returns the planted defect of a Shelf whose Snapshot
// panics over a fresh subject.
func shelfSnapshotPanics(h *ShelfHandle) Shelf {
	return &shelfSnapshotPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSnapshotIgnoresCancelDefect is the planted defect of a Shelf whose
// Snapshot ignores the context.
type shelfSnapshotIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Snapshot is the override of the planted defect.
func (d *shelfSnapshotIgnoresCancelDefect) Snapshot(p0 context.Context) (r0 []Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.Snapshot(p0)
}

// shelfSnapshotIgnoresCancel returns the planted defect of a Shelf whose
// Snapshot ignores the context over a fresh subject.
func shelfSnapshotIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfSnapshotIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSnapshotValueBesideDefect is the planted defect of a Shelf whose
// Snapshot returns a value beside an error.
type shelfSnapshotValueBesideDefect struct {
	shelfOverEmbedded
}

// Snapshot is the override of the planted defect.
func (d *shelfSnapshotValueBesideDefect) Snapshot(p0 context.Context) (r0 []Item, r1 error) {
	r0, r1 = d.shelfOverEmbedded.Snapshot(p0)
	if r1 != nil {
		r0 = []Item{{Key: "Key-b"}}
	}
	return r0, r1
}

// shelfSnapshotValueBeside returns the planted defect of a Shelf whose
// Snapshot returns a value beside an error over a fresh subject.
func shelfSnapshotValueBeside(h *ShelfHandle) Shelf {
	return &shelfSnapshotValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSnapshotDrainsNothingDefect is the planted defect of a Shelf whose
// Snapshot loses what Put wrote.
type shelfSnapshotDrainsNothingDefect struct {
	shelfOverEmbedded
}

// Snapshot is the override of the planted defect.
func (*shelfSnapshotDrainsNothingDefect) Snapshot(p0 context.Context) (r0 []Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	return r0, r1
}

// shelfSnapshotDrainsNothing returns the planted defect of a Shelf whose
// Snapshot loses what Put wrote over a fresh subject.
func shelfSnapshotDrainsNothing(h *ShelfHandle) Shelf {
	return &shelfSnapshotDrainsNothingDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSnapshotDropsAfterThreeDefect is the planted defect of a Shelf
// whose Snapshot drops its last element from its fourth call on.
type shelfSnapshotDropsAfterThreeDefect struct {
	shelfOverEmbedded
	calls int
}

// Snapshot is the override of the planted defect.
func (d *shelfSnapshotDropsAfterThreeDefect) Snapshot(p0 context.Context) (r0 []Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.shelfOverEmbedded.Snapshot(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = r0[:len(r0)-1]
	}
	return r0, r1
}

// shelfSnapshotDropsAfterThree returns the planted defect of a Shelf whose
// Snapshot drops its last element from its fourth call on over a fresh
// subject.
func shelfSnapshotDropsAfterThree(h *ShelfHandle) Shelf {
	return &shelfSnapshotDropsAfterThreeDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfLenPanicsDefect is the planted defect of a Shelf whose Len panics.
type shelfLenPanicsDefect struct {
	shelfOverEmbedded
}

// Len is the override of the planted defect.
func (*shelfLenPanicsDefect) Len(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Shelf whose Len panics")
}

// shelfLenPanics returns the planted defect of a Shelf whose Len panics
// over a fresh subject.
func shelfLenPanics(h *ShelfHandle) Shelf {
	return &shelfLenPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfLenIgnoresCancelDefect is the planted defect of a Shelf whose Len
// ignores the context.
type shelfLenIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Len is the override of the planted defect.
func (d *shelfLenIgnoresCancelDefect) Len(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.Len(p0)
}

// shelfLenIgnoresCancel returns the planted defect of a Shelf whose Len
// ignores the context over a fresh subject.
func shelfLenIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfLenIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfLenValueBesideDefect is the planted defect of a Shelf whose Len
// returns a value beside an error.
type shelfLenValueBesideDefect struct {
	shelfOverEmbedded
}

// Len is the override of the planted defect.
func (d *shelfLenValueBesideDefect) Len(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.shelfOverEmbedded.Len(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// shelfLenValueBeside returns the planted defect of a Shelf whose Len
// returns a value beside an error over a fresh subject.
func shelfLenValueBeside(h *ShelfHandle) Shelf {
	return &shelfLenValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtPanicsDefect is the planted defect of a Shelf whose At panics.
type shelfAtPanicsDefect struct {
	shelfOverEmbedded
}

// At is the override of the planted defect.
func (*shelfAtPanicsDefect) At(p0 context.Context, p1 int) (r0 Item, r1 error) {
	panic("the planted defect of a Shelf whose At panics")
}

// shelfAtPanics returns the planted defect of a Shelf whose At panics over
// a fresh subject.
func shelfAtPanics(h *ShelfHandle) Shelf {
	return &shelfAtPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtIgnoresCancelDefect is the planted defect of a Shelf whose At
// ignores the context.
type shelfAtIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// At is the override of the planted defect.
func (d *shelfAtIgnoresCancelDefect) At(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.At(p0, p1)
}

// shelfAtIgnoresCancel returns the planted defect of a Shelf whose At
// ignores the context over a fresh subject.
func shelfAtIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfAtIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtValueBesideDefect is the planted defect of a Shelf whose At
// returns a value beside an error.
type shelfAtValueBesideDefect struct {
	shelfOverEmbedded
}

// At is the override of the planted defect.
func (d *shelfAtValueBesideDefect) At(p0 context.Context, p1 int) (r0 Item, r1 error) {
	r0, r1 = d.shelfOverEmbedded.At(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// shelfAtValueBeside returns the planted defect of a Shelf whose At
// returns a value beside an error over a fresh subject.
func shelfAtValueBeside(h *ShelfHandle) Shelf {
	return &shelfAtValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtAdmitsTheSizeDefect is the planted defect of a Shelf whose At
// admits the position equal to the size.
type shelfAtAdmitsTheSizeDefect struct {
	shelfOverEmbedded
}

// At is the override of the planted defect.
func (*shelfAtAdmitsTheSizeDefect) At(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = Item{Key: "Key-b"}
	return r0, r1
}

// shelfAtAdmitsTheSize returns the planted defect of a Shelf whose At
// admits the position equal to the size over a fresh subject.
func shelfAtAdmitsTheSize(h *ShelfHandle) Shelf {
	return &shelfAtAdmitsTheSizeDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtOffsetPanicsDefect is the planted defect of a Shelf whose
// AtOffset panics.
type shelfAtOffsetPanicsDefect struct {
	shelfOverEmbedded
}

// AtOffset is the override of the planted defect.
func (*shelfAtOffsetPanicsDefect) AtOffset(p0 context.Context, p1 int64) (r0 Item, r1 error) {
	panic("the planted defect of a Shelf whose AtOffset panics")
}

// shelfAtOffsetPanics returns the planted defect of a Shelf whose AtOffset
// panics over a fresh subject.
func shelfAtOffsetPanics(h *ShelfHandle) Shelf {
	return &shelfAtOffsetPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtOffsetIgnoresCancelDefect is the planted defect of a Shelf whose
// AtOffset ignores the context.
type shelfAtOffsetIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// AtOffset is the override of the planted defect.
func (d *shelfAtOffsetIgnoresCancelDefect) AtOffset(
	p0 context.Context,
	p1 int64,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.AtOffset(p0, p1)
}

// shelfAtOffsetIgnoresCancel returns the planted defect of a Shelf whose
// AtOffset ignores the context over a fresh subject.
func shelfAtOffsetIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfAtOffsetIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtOffsetValueBesideDefect is the planted defect of a Shelf whose
// AtOffset returns a value beside an error.
type shelfAtOffsetValueBesideDefect struct {
	shelfOverEmbedded
}

// AtOffset is the override of the planted defect.
func (d *shelfAtOffsetValueBesideDefect) AtOffset(
	p0 context.Context,
	p1 int64,
) (r0 Item, r1 error) {
	r0, r1 = d.shelfOverEmbedded.AtOffset(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// shelfAtOffsetValueBeside returns the planted defect of a Shelf whose
// AtOffset returns a value beside an error over a fresh subject.
func shelfAtOffsetValueBeside(h *ShelfHandle) Shelf {
	return &shelfAtOffsetValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfAtOffsetAdmitsTheSizeDefect is the planted defect of a Shelf whose
// AtOffset admits the position equal to the size.
type shelfAtOffsetAdmitsTheSizeDefect struct {
	shelfOverEmbedded
}

// AtOffset is the override of the planted defect.
func (*shelfAtOffsetAdmitsTheSizeDefect) AtOffset(
	p0 context.Context,
	p1 int64,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = Item{Key: "Key-b"}
	return r0, r1
}

// shelfAtOffsetAdmitsTheSize returns the planted defect of a Shelf whose
// AtOffset admits the position equal to the size over a fresh subject.
func shelfAtOffsetAdmitsTheSize(h *ShelfHandle) Shelf {
	return &shelfAtOffsetAdmitsTheSizeDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSlotPanicsDefect is the planted defect of a Shelf whose Slot
// panics.
type shelfSlotPanicsDefect struct {
	shelfOverEmbedded
}

// Slot is the override of the planted defect.
func (*shelfSlotPanicsDefect) Slot(p0 context.Context, p1 string, p2 uint) (r0 Item, r1 error) {
	panic("the planted defect of a Shelf whose Slot panics")
}

// shelfSlotPanics returns the planted defect of a Shelf whose Slot panics
// over a fresh subject.
func shelfSlotPanics(h *ShelfHandle) Shelf {
	return &shelfSlotPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSlotIgnoresCancelDefect is the planted defect of a Shelf whose Slot
// ignores the context.
type shelfSlotIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// Slot is the override of the planted defect.
func (d *shelfSlotIgnoresCancelDefect) Slot(
	p0 context.Context,
	p1 string,
	p2 uint,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.Slot(p0, p1, p2)
}

// shelfSlotIgnoresCancel returns the planted defect of a Shelf whose Slot
// ignores the context over a fresh subject.
func shelfSlotIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfSlotIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfSlotAdmitsTheSizeDefect is the planted defect of a Shelf whose Slot
// admits the position equal to the size.
type shelfSlotAdmitsTheSizeDefect struct {
	shelfOverEmbedded
}

// Slot is the override of the planted defect.
func (*shelfSlotAdmitsTheSizeDefect) Slot(
	p0 context.Context,
	p1 string,
	p2 uint,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = Item{Key: "Key-b"}
	return r0, r1
}

// shelfSlotAdmitsTheSize returns the planted defect of a Shelf whose Slot
// admits the position equal to the size over a fresh subject.
func shelfSlotAdmitsTheSize(h *ShelfHandle) Shelf {
	return &shelfSlotAdmitsTheSizeDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfGetPanicsDefect is the planted defect of a Shelf whose Get panics.
type shelfGetPanicsDefect struct {
	shelfOverEmbedded
}

// Get is the override of the planted defect.
func (*shelfGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
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
func (d *shelfGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
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
func (d *shelfGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.shelfOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// shelfGetValueBeside returns the planted defect of a Shelf whose Get
// returns a value beside an error over a fresh subject.
func shelfGetValueBeside(h *ShelfHandle) Shelf {
	return &shelfGetValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfTracedGetPanicsDefect is the planted defect of a Shelf whose
// TracedGet panics.
type shelfTracedGetPanicsDefect struct {
	shelfOverEmbedded
}

// TracedGet is the override of the planted defect.
func (*shelfTracedGetPanicsDefect) TracedGet(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Shelf whose TracedGet panics")
}

// shelfTracedGetPanics returns the planted defect of a Shelf whose
// TracedGet panics over a fresh subject.
func shelfTracedGetPanics(h *ShelfHandle) Shelf {
	return &shelfTracedGetPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfTracedGetIgnoresCancelDefect is the planted defect of a Shelf whose
// TracedGet ignores the context.
type shelfTracedGetIgnoresCancelDefect struct {
	shelfOverEmbedded
}

// TracedGet is the override of the planted defect.
func (d *shelfTracedGetIgnoresCancelDefect) TracedGet(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.shelfOverEmbedded.TracedGet(p0, p1)
}

// shelfTracedGetIgnoresCancel returns the planted defect of a Shelf whose
// TracedGet ignores the context over a fresh subject.
func shelfTracedGetIgnoresCancel(h *ShelfHandle) Shelf {
	return &shelfTracedGetIgnoresCancelDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfTracedGetValueBesideDefect is the planted defect of a Shelf whose
// TracedGet returns a value beside an error.
type shelfTracedGetValueBesideDefect struct {
	shelfOverEmbedded
}

// TracedGet is the override of the planted defect.
func (d *shelfTracedGetValueBesideDefect) TracedGet(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	r0, r1 = d.shelfOverEmbedded.TracedGet(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// shelfTracedGetValueBeside returns the planted defect of a Shelf whose
// TracedGet returns a value beside an error over a fresh subject.
func shelfTracedGetValueBeside(h *ShelfHandle) Shelf {
	return &shelfTracedGetValueBesideDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfTracedGetStopsDelegatingDefect is the planted defect of a Shelf
// whose TracedGet stops calling Get.
type shelfTracedGetStopsDelegatingDefect struct {
	shelfOverEmbedded
}

// TracedGet is the override of the planted defect.
func (*shelfTracedGetStopsDelegatingDefect) TracedGet(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0 = Item{Key: "Key-b"}
	return r0, r1
}

// shelfTracedGetStopsDelegating returns the planted defect of a Shelf
// whose TracedGet stops calling Get over a fresh subject.
func shelfTracedGetStopsDelegating(h *ShelfHandle) Shelf {
	return &shelfTracedGetStopsDelegatingDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfErrPanicsDefect is the planted defect of a Shelf whose Err panics.
type shelfErrPanicsDefect struct {
	shelfOverEmbedded
}

// Err is the override of the planted defect.
func (*shelfErrPanicsDefect) Err() (r0 error) {
	panic("the planted defect of a Shelf whose Err panics")
}

// shelfErrPanics returns the planted defect of a Shelf whose Err panics
// over a fresh subject.
func shelfErrPanics(h *ShelfHandle) Shelf {
	return &shelfErrPanicsDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfErrReportsPoisonDefect is the planted defect of a Shelf whose Err
// reports an error for a fresh receiver.
type shelfErrReportsPoisonDefect struct {
	shelfOverEmbedded
}

// Err is the override of the planted defect.
func (*shelfErrReportsPoisonDefect) Err() (r0 error) {
	r0 = conform.ErrGenerated
	return r0
}

// shelfErrReportsPoison returns the planted defect of a Shelf whose Err
// reports an error for a fresh receiver over a fresh subject.
func shelfErrReportsPoison(h *ShelfHandle) Shelf {
	return &shelfErrReportsPoisonDefect{shelfOverEmbedded: h.Fresh()}
}

// shelfActions returns the actions of the generated sequences of Shelf
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func shelfActions(c *prop.Case, subject func() Shelf) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	keysInt := c.Draw(
		prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int",
	)
	keysInt64 := c.Draw(
		prop.List(prop.Of[int64](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int64",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "put.it")
				it.Key = c.Draw(prop.SampledFrom(keysString...), "put.it.key")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				key0, _ := args[0].(Item)
				call := c.History().Invoke(client, "Put", args, key0.Key)
				values, err := shelfPutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "GetMany",
			Input: func(c *prop.Case, _ struct{}) any {
				keys := c.Draw(prop.SampledFrom(keysString...), "get-many.keys")
				return []any{keys}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "GetMany", args, args[0])
				values, err := shelfGetManyCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "All",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "All", args)
				values, err := shelfAllCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Items",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Items", args)
				values, err := shelfItemsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Scan",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Scan", args)
				values, err := shelfScanCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Keyed",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Keyed", args)
				values, err := shelfKeyedCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Snapshot",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Snapshot", args)
				values, err := shelfSnapshotCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Len",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Len", args)
				values, err := shelfLenCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "At",
			Input: func(c *prop.Case, _ struct{}) any {
				i := c.Draw(prop.SampledFrom(keysInt...), "at.i")
				return []any{i}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "At", args, args[0])
				values, err := shelfAtCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "AtOffset",
			Input: func(c *prop.Case, _ struct{}) any {
				i := c.Draw(prop.SampledFrom(keysInt64...), "at-offset.i")
				return []any{i}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "AtOffset", args, args[0])
				values, err := shelfAtOffsetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Slot",
			Input: func(c *prop.Case, _ struct{}) any {
				name := c.Draw(prop.Of[string](), "slot.name")
				i := c.Draw(prop.Of[uint](), "slot.i")
				return []any{name, i}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Slot", args)
				values, err := shelfSlotCall(c.Context(), subject(), args)
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
		{
			Name: "TracedGet",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "traced-get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "TracedGet", args, args[0])
				values, err := shelfTracedGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Err",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Err", args)
				values, err := shelfErrCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// shelfPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfPutCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Put(ctx, a0)
}

// shelfGetManyCall calls GetMany on subject with the arguments args, and
// returns the values of the call beside its error.
func shelfGetManyCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].([]string)
	r0, err := subject.GetMany(ctx, a0...)
	return []any{r0}, err
}

// shelfAllCall calls All on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfAllCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	r0 := subject.All(ctx)
	var drained []Item
	for e, err := range r0 {
		if err != nil {
			return []any{drained}, err
		}
		drained = append(drained, e)
	}
	return []any{drained}, nil
}

// shelfItemsCall calls Items on subject with the arguments args, and
// returns the values of the call beside its error.
func shelfItemsCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	r0 := subject.Items(ctx)
	return []any{slices.Collect(r0)}, nil
}

// shelfScanCall calls Scan on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfScanCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	r0, err := subject.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return []any{slices.Collect(r0)}, nil
}

// shelfKeyedCall calls Keyed on subject with the arguments args, and
// returns the values of the call beside its error.
func shelfKeyedCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	r0 := subject.Keyed(ctx)
	var drained []probe.Pair[string, Item]
	for k, v := range r0 {
		drained = append(drained, probe.Pair[string, Item]{First: k, Second: v})
	}
	return []any{drained}, nil
}

// shelfSnapshotCall calls Snapshot on subject with the arguments args, and
// returns the values of the call beside its error.
func shelfSnapshotCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	r0, err := subject.Snapshot(ctx)
	return []any{r0}, err
}

// shelfLenCall calls Len on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfLenCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	r0, err := subject.Len(ctx)
	return []any{r0}, err
}

// shelfAtCall calls At on subject with the arguments args, and returns the
// values of the call beside its error.
func shelfAtCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(int)
	r0, err := subject.At(ctx, a0)
	return []any{r0}, err
}

// shelfAtOffsetCall calls AtOffset on subject with the arguments args, and
// returns the values of the call beside its error.
func shelfAtOffsetCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(int64)
	r0, err := subject.AtOffset(ctx, a0)
	return []any{r0}, err
}

// shelfSlotCall calls Slot on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfSlotCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(uint)
	r0, err := subject.Slot(ctx, a0, a1)
	return []any{r0}, err
}

// shelfGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfGetCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// shelfTracedGetCall calls TracedGet on subject with the arguments args,
// and returns the values of the call beside its error.
func shelfTracedGetCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.TracedGet(ctx, a0)
	return []any{r0}, err
}

// shelfErrCall calls Err on subject with the arguments args, and returns
// the values of the call beside its error.
func shelfErrCall(ctx context.Context, subject Shelf, args []any) ([]any, error) {
	return nil, subject.Err()
}

// The derived checks of Tags:
//
// Tags.Add/smoke: Add returns without a panic for derived inputs.
//
// Tags.Add/smoke#generated: Add returns without a panic for derived
// inputs.
//
// Tags.Add/cancel: Add returns an error for a cancelled context.
//
// Tags.Add/cancel#generated: Add returns an error for a cancelled context.
//
// Tags.Has/smoke: Has returns without a panic for derived inputs.
//
// Tags.Has/smoke#generated: Has returns without a panic for derived
// inputs.
//
// Tags.Has/cancel: Has returns an error for a cancelled context.
//
// Tags.Has/cancel#generated: Has returns an error for a cancelled context.
//
// Tags.Has/reads-each-key: Has returns the value that Add stored under
// each key.
//
// Tags.Has/reads-each-key#generated: Has returns the value that Add stored
// under each key.

// The derived checks of Tags that witness their laws at the shallow end:
//
// Tags.Add/smoke
//
// Tags.Has/smoke

// The unclaimed classifications of Tags:
//
// The shape:writer of Tags.Add has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Tags writes.
type (
	TagsHandle = conform.Handle[Tags]
	TagsCheck  = conform.Check[Tags]
	TagsOption = conform.Option[Tags]
)

// TagsChecks contains the name of each callable of Tags, with the
// identities of its checks.
var TagsChecks = struct {
	Add tagsAdd
	Has tagsHas
}{
	Add: tagsAdd{"Tags.Add"},
	Has: tagsHas{"Tags.Has"},
}

// tagsAdd is the name of Tags.Add, with the identities of its checks.
type tagsAdd struct {
	conform.Callable[Tags]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c tagsAdd) Smoke() conform.ID[Tags] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Add returns
// without a panic for derived inputs.
func (c tagsAdd) SmokeGenerated() conform.ID[Tags] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c tagsAdd) Cancel() conform.ID[Tags] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Add returns an
// error for a cancelled context.
func (c tagsAdd) CancelGenerated() conform.ID[Tags] {
	return c.Check("cancel#generated")
}

// tagsHas is the name of Tags.Has, with the identities of its checks.
type tagsHas struct {
	conform.Callable[Tags]
}

// Smoke returns the identity of the check that Has returns without a panic
// for derived inputs.
func (c tagsHas) Smoke() conform.ID[Tags] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Has returns
// without a panic for derived inputs.
func (c tagsHas) SmokeGenerated() conform.ID[Tags] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Has returns an error for a
// cancelled context.
func (c tagsHas) Cancel() conform.ID[Tags] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Has returns an
// error for a cancelled context.
func (c tagsHas) CancelGenerated() conform.ID[Tags] {
	return c.Check("cancel#generated")
}

// ReadsEachKey returns the identity of the check that Has returns the
// value that Add stored under each key.
func (c tagsHas) ReadsEachKey() conform.ID[Tags] {
	return c.Check("reads-each-key")
}

// ReadsEachKeyGenerated returns the identity of the check that Has returns
// the value that Add stored under each key.
func (c tagsHas) ReadsEachKeyGenerated() conform.ID[Tags] {
	return c.Check("reads-each-key#generated")
}

// TagsSuite returns the checks of Tags, and the calls and the budgets of
// its callables.
func TagsSuite() conform.Suite[Tags] {
	return conform.Suite[Tags]{
		Declaration: "Tags",
		Package:     "example.com/shop/classified",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Tags]{
			{
				ID:          TagsChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TagsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("tag-a"))
					})
				},
				Defect: tagsAddPanics,
			},
			{
				ID:          TagsChecks.Add.SmokeGenerated(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TagsHandle, c *prop.Case) {
					drawnTag := c.Draw(prop.Of[string](), "tag")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnTag)
					})
				},
				Defect: tagsAddPanics,
			},
			{
				ID:          TagsChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TagsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("tag-a"))
					})
				},
				Defect: tagsAddIgnoresCancel,
			},
			{
				ID:          TagsChecks.Add.CancelGenerated(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TagsHandle, c *prop.Case) {
					drawnTag := c.Draw(prop.Of[string](), "tag")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnTag)
					})
				},
				Defect: tagsAddIgnoresCancel,
			},
			{
				ID:          TagsChecks.Has.Smoke(),
				Statement:   "Has returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TagsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, string("tags-a"))
						return err
					})
				},
				Defect: tagsHasPanics,
			},
			{
				ID:          TagsChecks.Has.SmokeGenerated(),
				Statement:   "Has returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TagsHandle, c *prop.Case) {
					drawnTags := c.Draw(prop.Of[string](), "tags")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, drawnTags)
						return err
					})
				},
				Defect: tagsHasPanics,
			},
			{
				ID:          TagsChecks.Has.Cancel(),
				Statement:   "Has returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TagsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, string("tags-a"))
						return err
					})
				},
				Defect: tagsHasIgnoresCancel,
			},
			{
				ID:          TagsChecks.Has.CancelGenerated(),
				Statement:   "Has returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TagsHandle, c *prop.Case) {
					drawnTags := c.Draw(prop.Of[string](), "tags")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Has(ctx, drawnTags)
						return err
					})
				},
				Defect: tagsHasIgnoresCancel,
			},
			{
				ID:          TagsChecks.Has.ReadsEachKey(),
				Statement:   "Has returns the value that Add stored under each key",
				Claim:       spec.ReadsEachKey,
				Probe:       spec.ReadsWhatWasStored,
				Reason:      spec.LostAKey,
				Law:         spec.BatchReadReturnsEachStoredValue,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TagsHandle) {
					probe.ReadsWhatWasStored(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Has(ctx, string("tag-a"), string("tag-b"))
					}, func() func(ctx context.Context) (string, error) {
						calls := 0
						return func(ctx context.Context) (string, error) {
							calls++
							w := string("tag-a")
							if calls > 1 {
								w = string("tag-b")
							}
							err := h.Subject().Add(ctx, w)
							return w, err
						}
					}())
				},
				Defect: tagsHasDropsTheLastKey,
			},
			{
				ID:          TagsChecks.Has.ReadsEachKeyGenerated(),
				Statement:   "Has returns the value that Add stored under each key",
				Claim:       spec.ReadsEachKey,
				Probe:       spec.ReadsWhatWasStored,
				Reason:      spec.LostAKey,
				Law:         spec.BatchReadReturnsEachStoredValue,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *TagsHandle, c *prop.Case) {
					drawnTag := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"tag",
					)
					probe.ReadsWhatWasStored(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Has(ctx, drawnTag[0], drawnTag[1])
					}, func() func(ctx context.Context) (string, error) {
						calls := 0
						return func(ctx context.Context) (string, error) {
							calls++
							w := drawnTag[0]
							if calls > 1 {
								w = drawnTag[1]
							}
							err := h.Subject().Add(ctx, w)
							return w, err
						}
					}())
				},
				Defect: tagsHasDropsTheLastKey,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Tags.Add",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Tags]{
			{
				Calls:    tagsAddCalls,
				Over:     newTagsOver,
				Callable: TagsChecks.Add.Callable,
			},
			{
				Calls:    tagsHasCalls,
				Over:     newTagsOver,
				Callable: TagsChecks.Has.Callable,
			},
		},
	}
}

// RunTagsSuite applies the checks of Tags, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunTagsSuite[T Tags](t *testing.T, newSubject func() T, opts ...TagsOption) {
	t.Helper()
	RunTagsSuiteFromHandle(t, func(*TagsHandle) T { return newSubject() }, opts...)
}

// RunTagsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTagsSuiteFromHandle[T Tags](
	t *testing.T,
	newSubject func(h *TagsHandle) T,
	opts ...TagsOption,
) {
	t.Helper()
	conform.Run(t, TagsSuite(), tagsSubjects(newSubject), opts...)
}

// FuzzTagsSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTagsSuite[T Tags](
	f *testing.F,
	callable conform.Callable[Tags],
	newSubject func(h *TagsHandle) T,
	opts ...TagsOption,
) {
	f.Helper()
	conform.Fuzz(f, TagsSuite(), callable, tagsSubjects(newSubject), opts...)
}

// RunTagsBudgets applies the allocation ceilings of Tags and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunTagsBudgets[T Tags](t *testing.T, newSubject func(h *TagsHandle) T, opts ...TagsOption) {
	t.Helper()
	conform.RunBudgets(t, TagsSuite(), tagsSubjects(newSubject), opts...)
}

// BenchTagsSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchTagsSuite[T Tags](
	b conform.B,
	callable conform.Callable[Tags],
	newSubject func(h *TagsHandle) T,
	opts ...TagsOption,
) {
	b.Helper()
	conform.Bench(b, TagsSuite(), callable, tagsSubjects(newSubject), opts...)
}

// tagsSubjects returns the factory of the subjects of Tags that calls
// newSubject with the handle of each check.
func tagsSubjects[T Tags](newSubject func(h *TagsHandle) T) func(h *TagsHandle) Tags {
	return func(h *TagsHandle) Tags { return newSubject(h) }
}

// tagsAddCalls generates the calls of Tags.Add with the inputs that each
// case draws.
var tagsAddCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tags) error {
	tag := c.Draw(prop.Of[string](), "tag")
	return func(ctx context.Context, subject Tags) error {
		return subject.Add(ctx, tag)
	}
})

// tagsHasCalls generates the calls of Tags.Has with the inputs that each
// case draws.
var tagsHasCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tags) error {
	tags := c.Draw(prop.List(prop.Of[string]()), "tags")
	return func(ctx context.Context, subject Tags) error {
		_, err := subject.Has(ctx, tags...)
		return err
	}
})

// tagsOverEmbedded is Tags under the name of the field that tagsOver
// embeds.
type tagsOverEmbedded = Tags

// tagsOver is a subject of Tags whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type tagsOver struct {
	tagsOverEmbedded
	spend func()
}

// Add calls spend and returns its zero results.
func (o tagsOver) Add(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Has calls spend and returns its zero results.
func (o tagsOver) Has(context.Context, ...string) (r0 []string, r1 error) {
	o.spend()
	return r0, r1
}

// newTagsOver returns the planted defect of a ceiling of Tags, whose
// callables call spend.
func newTagsOver(_ *TagsHandle, spend func()) Tags {
	return tagsOver{spend: spend}
}

// tagsAddPanicsDefect is the planted defect of a Tags whose Add panics.
type tagsAddPanicsDefect struct {
	tagsOverEmbedded
}

// Add is the override of the planted defect.
func (*tagsAddPanicsDefect) Add(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Tags whose Add panics")
}

// tagsAddPanics returns the planted defect of a Tags whose Add panics over
// a fresh subject.
func tagsAddPanics(h *TagsHandle) Tags {
	return &tagsAddPanicsDefect{tagsOverEmbedded: h.Fresh()}
}

// tagsAddIgnoresCancelDefect is the planted defect of a Tags whose Add
// ignores the context.
type tagsAddIgnoresCancelDefect struct {
	tagsOverEmbedded
}

// Add is the override of the planted defect.
func (d *tagsAddIgnoresCancelDefect) Add(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.tagsOverEmbedded.Add(p0, p1)
}

// tagsAddIgnoresCancel returns the planted defect of a Tags whose Add
// ignores the context over a fresh subject.
func tagsAddIgnoresCancel(h *TagsHandle) Tags {
	return &tagsAddIgnoresCancelDefect{tagsOverEmbedded: h.Fresh()}
}

// tagsHasPanicsDefect is the planted defect of a Tags whose Has panics.
type tagsHasPanicsDefect struct {
	tagsOverEmbedded
}

// Has is the override of the planted defect.
func (*tagsHasPanicsDefect) Has(p0 context.Context, p1 ...string) (r0 []string, r1 error) {
	panic("the planted defect of a Tags whose Has panics")
}

// tagsHasPanics returns the planted defect of a Tags whose Has panics over
// a fresh subject.
func tagsHasPanics(h *TagsHandle) Tags {
	return &tagsHasPanicsDefect{tagsOverEmbedded: h.Fresh()}
}

// tagsHasIgnoresCancelDefect is the planted defect of a Tags whose Has
// ignores the context.
type tagsHasIgnoresCancelDefect struct {
	tagsOverEmbedded
}

// Has is the override of the planted defect.
func (d *tagsHasIgnoresCancelDefect) Has(p0 context.Context, p1 ...string) (r0 []string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.tagsOverEmbedded.Has(p0, p1...)
}

// tagsHasIgnoresCancel returns the planted defect of a Tags whose Has
// ignores the context over a fresh subject.
func tagsHasIgnoresCancel(h *TagsHandle) Tags {
	return &tagsHasIgnoresCancelDefect{tagsOverEmbedded: h.Fresh()}
}

// tagsHasDropsTheLastKeyDefect is the planted defect of a Tags whose Has
// drops the value of its last key.
type tagsHasDropsTheLastKeyDefect struct {
	tagsOverEmbedded
}

// Has is the override of the planted defect.
func (d *tagsHasDropsTheLastKeyDefect) Has(
	p0 context.Context,
	p1 ...string,
) (r0 []string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.tagsOverEmbedded.Has(p0, p1...)
	if len(r0) > 0 {
		r0 = r0[:len(r0)-1]
	}
	return r0, r1
}

// tagsHasDropsTheLastKey returns the planted defect of a Tags whose Has
// drops the value of its last key over a fresh subject.
func tagsHasDropsTheLastKey(h *TagsHandle) Tags {
	return &tagsHasDropsTheLastKeyDefect{tagsOverEmbedded: h.Fresh()}
}

// The derived checks of Vault:
//
// Vault.Append/smoke: Append returns without a panic for derived inputs.
//
// Vault.Append/smoke#generated: Append returns without a panic for derived
// inputs.
//
// Vault.Append/cancel: Append returns an error for a cancelled context.
//
// Vault.Append/cancel#generated: Append returns an error for a cancelled
// context.
//
// Vault.Append/chain: Append grows Replay by the entry that it admits.
//
// Vault.Append/chain#generated: Append grows Replay by the entry that it
// admits.
//
// Vault.Replay/smoke: Replay returns without a panic for derived inputs.
//
// Vault.Replay/cancel: Replay returns an error for a cancelled context.
//
// Vault.Replay/zero-on-error: Replay returns the zero value beside an
// error.
//
// Vault.Verify/smoke: Verify returns without a panic for derived inputs.
//
// Vault.Verify/cancel: Verify returns an error for a cancelled context.
//
// Vault.Verify/verifies-what-it-appended: Verify reports no error for a
// chain that Append produced.
//
// Vault.Verify/verifies-what-it-appended#generated: Verify reports no
// error for a chain that Append produced.
//
// Vault.Acquire/smoke: Acquire returns without a panic for derived inputs.
//
// Vault.Acquire/smoke#generated: Acquire returns without a panic for
// derived inputs.
//
// Vault.Acquire/cancel: Acquire returns an error for a cancelled context.
//
// Vault.Acquire/cancel#generated: Acquire returns an error for a cancelled
// context.
//
// Vault.Acquire/lease: a second Acquire for one key fails with the
// declared sentinel.
//
// Vault.Acquire/lease#generated: a second Acquire for one key fails with
// the declared sentinel.
//
// Vault.Acquire/release-frees-the-lease: Acquire grants a lease again
// after Release released it.
//
// Vault.Acquire/release-frees-the-lease#generated: Acquire grants a lease
// again after Release released it.
//
// Vault.Release/smoke: Release returns without a panic for derived inputs.
//
// Vault.Release/smoke#generated: Release returns without a panic for
// derived inputs.
//
// Vault.Release/cancel: Release returns an error for a cancelled context.
//
// Vault.Release/cancel#generated: Release returns an error for a cancelled
// context.
//
// Vault.Borrow/smoke: Borrow returns without a panic for derived inputs.
//
// Vault.Borrow/cancel: Borrow returns an error for a cancelled context.
//
// Vault.Borrow/zero-on-error: Borrow returns the zero value beside an
// error.
//
// Vault.Borrow/pool: Return takes back what Borrow lent.
//
// Vault.Return/smoke: Return returns without a panic for derived inputs.
//
// Vault.Return/smoke#generated: Return returns without a panic for derived
// inputs.
//
// Vault.Return/cancel: Return returns an error for a cancelled context.
//
// Vault.Return/cancel#generated: Return returns an error for a cancelled
// context.
//
// Vault.Stats/smoke: Stats returns without a panic for derived inputs.
//
// Vault.Stats/cancel: Stats returns an error for a cancelled context.
//
// Vault.Stats/zero-on-error: Stats returns the zero value beside an error.
//
// Vault.Stats/counts-each-loan: Stats reports another count after Borrow
// lends an object.

// The derived checks of Vault that witness their laws at the shallow end:
//
// Vault.Append/smoke
//
// Vault.Replay/smoke
//
// Vault.Verify/smoke
//
// Vault.Acquire/smoke
//
// Vault.Acquire/lease
//
// Vault.Acquire/lease#generated
//
// Vault.Release/smoke
//
// Vault.Borrow/smoke
//
// Vault.Return/smoke
//
// Vault.Stats/smoke

// The unclaimed classifications of Vault:
//
// The shape:writer of Vault.Append has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Vault.Verify has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The shape:writer of Vault.Acquire has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Vault.Release has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Vault.Return has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

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
	Append  vaultAppend
	Replay  vaultReplay
	Verify  vaultVerify
	Acquire vaultAcquire
	Release vaultRelease
	Borrow  vaultBorrow
	Return  vaultReturn
	Stats   vaultStats
}{
	Append:  vaultAppend{"Vault.Append"},
	Replay:  vaultReplay{"Vault.Replay"},
	Verify:  vaultVerify{"Vault.Verify"},
	Acquire: vaultAcquire{"Vault.Acquire"},
	Release: vaultRelease{"Vault.Release"},
	Borrow:  vaultBorrow{"Vault.Borrow"},
	Return:  vaultReturn{"Vault.Return"},
	Stats:   vaultStats{"Vault.Stats"},
}

// vaultAppend is the name of Vault.Append, with the identities of its
// checks.
type vaultAppend struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c vaultAppend) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Append returns
// without a panic for derived inputs.
func (c vaultAppend) SmokeGenerated() conform.ID[Vault] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c vaultAppend) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Append returns an
// error for a cancelled context.
func (c vaultAppend) CancelGenerated() conform.ID[Vault] {
	return c.Check("cancel#generated")
}

// Chain returns the identity of the check that Append grows Replay by the
// entry that it admits.
func (c vaultAppend) Chain() conform.ID[Vault] {
	return c.Check("chain")
}

// ChainGenerated returns the identity of the check that Append grows
// Replay by the entry that it admits.
func (c vaultAppend) ChainGenerated() conform.ID[Vault] {
	return c.Check("chain#generated")
}

// vaultReplay is the name of Vault.Replay, with the identities of its
// checks.
type vaultReplay struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Replay returns without a
// panic for derived inputs.
func (c vaultReplay) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Replay returns an error
// for a cancelled context.
func (c vaultReplay) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Replay returns the
// zero value beside an error.
func (c vaultReplay) ZeroOnError() conform.ID[Vault] {
	return c.Check("zero-on-error")
}

// vaultVerify is the name of Vault.Verify, with the identities of its
// checks.
type vaultVerify struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Verify returns without a
// panic for derived inputs.
func (c vaultVerify) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Verify returns an error
// for a cancelled context.
func (c vaultVerify) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// VerifiesWhatItAppended returns the identity of the check that Verify
// reports no error for a chain that Append produced.
func (c vaultVerify) VerifiesWhatItAppended() conform.ID[Vault] {
	return c.Check("verifies-what-it-appended")
}

// VerifiesWhatItAppendedGenerated returns the identity of the check that
// Verify reports no error for a chain that Append produced.
func (c vaultVerify) VerifiesWhatItAppendedGenerated() conform.ID[Vault] {
	return c.Check("verifies-what-it-appended#generated")
}

// vaultAcquire is the name of Vault.Acquire, with the identities of its
// checks.
type vaultAcquire struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Acquire returns without a
// panic for derived inputs.
func (c vaultAcquire) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Acquire returns
// without a panic for derived inputs.
func (c vaultAcquire) SmokeGenerated() conform.ID[Vault] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Acquire returns an error
// for a cancelled context.
func (c vaultAcquire) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Acquire returns
// an error for a cancelled context.
func (c vaultAcquire) CancelGenerated() conform.ID[Vault] {
	return c.Check("cancel#generated")
}

// Lease returns the identity of the check that a second Acquire for one
// key fails with the declared sentinel.
func (c vaultAcquire) Lease() conform.ID[Vault] {
	return c.Check("lease")
}

// LeaseGenerated returns the identity of the check that a second Acquire
// for one key fails with the declared sentinel.
func (c vaultAcquire) LeaseGenerated() conform.ID[Vault] {
	return c.Check("lease#generated")
}

// ReleaseFreesTheLease returns the identity of the check that Acquire
// grants a lease again after Release released it.
func (c vaultAcquire) ReleaseFreesTheLease() conform.ID[Vault] {
	return c.Check("release-frees-the-lease")
}

// ReleaseFreesTheLeaseGenerated returns the identity of the check that
// Acquire grants a lease again after Release released it.
func (c vaultAcquire) ReleaseFreesTheLeaseGenerated() conform.ID[Vault] {
	return c.Check("release-frees-the-lease#generated")
}

// vaultRelease is the name of Vault.Release, with the identities of its
// checks.
type vaultRelease struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Release returns without a
// panic for derived inputs.
func (c vaultRelease) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Release returns
// without a panic for derived inputs.
func (c vaultRelease) SmokeGenerated() conform.ID[Vault] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Release returns an error
// for a cancelled context.
func (c vaultRelease) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Release returns
// an error for a cancelled context.
func (c vaultRelease) CancelGenerated() conform.ID[Vault] {
	return c.Check("cancel#generated")
}

// vaultBorrow is the name of Vault.Borrow, with the identities of its
// checks.
type vaultBorrow struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Borrow returns without a
// panic for derived inputs.
func (c vaultBorrow) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Borrow returns an error
// for a cancelled context.
func (c vaultBorrow) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Borrow returns the
// zero value beside an error.
func (c vaultBorrow) ZeroOnError() conform.ID[Vault] {
	return c.Check("zero-on-error")
}

// Pool returns the identity of the check that Return takes back what
// Borrow lent.
func (c vaultBorrow) Pool() conform.ID[Vault] {
	return c.Check("pool")
}

// vaultReturn is the name of Vault.Return, with the identities of its
// checks.
type vaultReturn struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Return returns without a
// panic for derived inputs.
func (c vaultReturn) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Return returns
// without a panic for derived inputs.
func (c vaultReturn) SmokeGenerated() conform.ID[Vault] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Return returns an error
// for a cancelled context.
func (c vaultReturn) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Return returns an
// error for a cancelled context.
func (c vaultReturn) CancelGenerated() conform.ID[Vault] {
	return c.Check("cancel#generated")
}

// vaultStats is the name of Vault.Stats, with the identities of its
// checks.
type vaultStats struct {
	conform.Callable[Vault]
}

// Smoke returns the identity of the check that Stats returns without a
// panic for derived inputs.
func (c vaultStats) Smoke() conform.ID[Vault] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Stats returns an error for
// a cancelled context.
func (c vaultStats) Cancel() conform.ID[Vault] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Stats returns the
// zero value beside an error.
func (c vaultStats) ZeroOnError() conform.ID[Vault] {
	return c.Check("zero-on-error")
}

// CountsEachLoan returns the identity of the check that Stats reports
// another count after Borrow lends an object.
func (c vaultStats) CountsEachLoan() conform.ID[Vault] {
	return c.Check("counts-each-loan")
}

// VaultSuite returns the checks of Vault, and the calls and the budgets of
// its callables.
func VaultSuite() conform.Suite[Vault] {
	return conform.Suite[Vault]{
		Declaration: "Vault",
		Package:     "example.com/shop/classified",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Vault]{
			{
				ID:          VaultChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: vaultAppendPanics,
			},
			{
				ID:          VaultChecks.Append.SmokeGenerated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnIt)
					})
				},
				Defect: vaultAppendPanics,
			},
			{
				ID:          VaultChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: vaultAppendIgnoresCancel,
			},
			{
				ID:          VaultChecks.Append.CancelGenerated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnIt)
					})
				},
				Defect: vaultAppendIgnoresCancel,
			},
			{
				ID:          VaultChecks.Append.Chain(),
				Statement:   "Append grows Replay by the entry that it admits",
				Claim:       spec.Chain,
				Probe:       spec.ReplayGrows,
				Reason:      spec.ReplayNotGrown,
				Law:         spec.AdmittedEntryGrowsTheReplay,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.ReplayGrows(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, Item{Key: "Key-a"})
					}, func(ctx context.Context) ([]Item, error) {
						return h.Subject().Replay(ctx)
					})
				},
				Defect: vaultAppendSucceedsAlways,
			},
			{
				ID:          VaultChecks.Append.ChainGenerated(),
				Statement:   "Append grows Replay by the entry that it admits",
				Claim:       spec.Chain,
				Probe:       spec.ReplayGrows,
				Reason:      spec.ReplayNotGrown,
				Law:         spec.AdmittedEntryGrowsTheReplay,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.ReplayGrows(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnIt)
					}, func(ctx context.Context) ([]Item, error) {
						return h.Subject().Replay(ctx)
					})
				},
				Defect: vaultAppendSucceedsAlways,
			},
			{
				ID:          VaultChecks.Replay.Smoke(),
				Statement:   "Replay returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Replay(ctx)
						return err
					})
				},
				Defect: vaultReplayPanics,
			},
			{
				ID:          VaultChecks.Replay.Cancel(),
				Statement:   "Replay returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Replay(ctx)
						return err
					})
				},
				Defect: vaultReplayIgnoresCancel,
			},
			{
				ID:          VaultChecks.Replay.ZeroOnError(),
				Statement:   "Replay returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]Item, error) {
						return h.Subject().Replay(ctx)
					}, []Item(nil))
				},
				Defect: vaultReplayValueBeside,
			},
			{
				ID:          VaultChecks.Verify.Smoke(),
				Statement:   "Verify returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Verify(ctx)
					})
				},
				Defect: vaultVerifyPanics,
			},
			{
				ID:          VaultChecks.Verify.Cancel(),
				Statement:   "Verify returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Verify(ctx)
					})
				},
				Defect: vaultVerifyIgnoresCancel,
			},
			{
				ID:          VaultChecks.Verify.VerifiesWhatItAppended(),
				Statement:   "Verify reports no error for a chain that Append produced",
				Claim:       spec.VerifiesWhatItAppended,
				Probe:       spec.VerifiesAfterAppends,
				Reason:      spec.FalseTamper,
				Law:         spec.AppendedChainVerifies,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.VerifiesAfterAppends(h, func(ctx context.Context) error {
						return h.Subject().Verify(ctx)
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := Item{Key: "Key-a"}
							if calls > 1 {
								w = Item{Key: "Key-b"}
							}
							return h.Subject().Append(ctx, w)
						}
					}())
				},
				Defect: vaultVerifyReportsTamper,
			},
			{
				ID:          VaultChecks.Verify.VerifiesWhatItAppendedGenerated(),
				Statement:   "Verify reports no error for a chain that Append produced",
				Claim:       spec.VerifiesWhatItAppended,
				Probe:       spec.VerifiesAfterAppends,
				Reason:      spec.FalseTamper,
				Law:         spec.AppendedChainVerifies,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnIt := c.Draw(
						prop.List(prop.Of[Item](), prop.Unique(), prop.MinSize(2), prop.MaxSize(2)),
						"it",
					)
					probe.VerifiesAfterAppends(h, func(ctx context.Context) error {
						return h.Subject().Verify(ctx)
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnIt[0]
							if calls > 1 {
								w = drawnIt[1]
							}
							return h.Subject().Append(ctx, w)
						}
					}())
				},
				Defect: vaultVerifyReportsTamper,
			},
			{
				ID:          VaultChecks.Acquire.Smoke(),
				Statement:   "Acquire returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, string("key-a"))
					})
				},
				Defect: vaultAcquirePanics,
			},
			{
				ID:          VaultChecks.Acquire.SmokeGenerated(),
				Statement:   "Acquire returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, drawnKey)
					})
				},
				Defect: vaultAcquirePanics,
			},
			{
				ID:          VaultChecks.Acquire.Cancel(),
				Statement:   "Acquire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, string("key-a"))
					})
				},
				Defect: vaultAcquireIgnoresCancel,
			},
			{
				ID:          VaultChecks.Acquire.CancelGenerated(),
				Statement:   "Acquire returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, drawnKey)
					})
				},
				Defect: vaultAcquireIgnoresCancel,
			},
			{
				ID:          VaultChecks.Acquire.Lease(),
				Statement:   "a second Acquire for one key fails with the declared sentinel",
				Claim:       spec.Lease,
				Probe:       spec.RefusesTheSecondWrite,
				Reason:      spec.UnguardedWrite,
				Law:         spec.GrantedLeaseIsNotGrantedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.RefusesTheSecondWrite(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, string("key-a"))
					}, ErrHeld)
				},
				Defect: vaultAcquireSucceedsAlways,
			},
			{
				ID:          VaultChecks.Acquire.LeaseGenerated(),
				Statement:   "a second Acquire for one key fails with the declared sentinel",
				Claim:       spec.Lease,
				Probe:       spec.RefusesTheSecondWrite,
				Reason:      spec.UnguardedWrite,
				Law:         spec.GrantedLeaseIsNotGrantedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.RefusesTheSecondWrite(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, drawnKey)
					}, ErrHeld)
				},
				Defect: vaultAcquireSucceedsAlways,
			},
			{
				ID:          VaultChecks.Acquire.ReleaseFreesTheLease(),
				Statement:   "Acquire grants a lease again after Release released it",
				Claim:       spec.ReleaseFreesTheLease,
				Probe:       spec.GrantsAfterRelease,
				Reason:      spec.KeptAfterRelease,
				Law:         spec.ReleasedLeaseIsGrantedAgain,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.GrantsAfterRelease(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, string("key-a"))
					}, func(ctx context.Context) error {
						return h.Subject().Release(ctx, string("key-a"))
					})
				},
				Defect: vaultReleaseReleaseKeepsTheLease,
			},
			{
				ID:          VaultChecks.Acquire.ReleaseFreesTheLeaseGenerated(),
				Statement:   "Acquire grants a lease again after Release released it",
				Claim:       spec.ReleaseFreesTheLease,
				Probe:       spec.GrantsAfterRelease,
				Reason:      spec.KeptAfterRelease,
				Law:         spec.ReleasedLeaseIsGrantedAgain,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.GrantsAfterRelease(h, func(ctx context.Context) error {
						return h.Subject().Acquire(ctx, drawnKey)
					}, func(ctx context.Context) error {
						return h.Subject().Release(ctx, drawnKey)
					})
				},
				Defect: vaultReleaseReleaseKeepsTheLease,
			},
			{
				ID:          VaultChecks.Release.Smoke(),
				Statement:   "Release returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Release(ctx, string("key-a"))
					})
				},
				Defect: vaultReleasePanics,
			},
			{
				ID:          VaultChecks.Release.SmokeGenerated(),
				Statement:   "Release returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Release(ctx, drawnKey)
					})
				},
				Defect: vaultReleasePanics,
			},
			{
				ID:          VaultChecks.Release.Cancel(),
				Statement:   "Release returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Release(ctx, string("key-a"))
					})
				},
				Defect: vaultReleaseIgnoresCancel,
			},
			{
				ID:          VaultChecks.Release.CancelGenerated(),
				Statement:   "Release returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Release(ctx, drawnKey)
					})
				},
				Defect: vaultReleaseIgnoresCancel,
			},
			{
				ID:          VaultChecks.Borrow.Smoke(),
				Statement:   "Borrow returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Borrow(ctx)
						return err
					})
				},
				Defect: vaultBorrowPanics,
			},
			{
				ID:          VaultChecks.Borrow.Cancel(),
				Statement:   "Borrow returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Borrow(ctx)
						return err
					})
				},
				Defect: vaultBorrowIgnoresCancel,
			},
			{
				ID:          VaultChecks.Borrow.ZeroOnError(),
				Statement:   "Borrow returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Borrow(ctx)
					}, Item{})
				},
				Defect: vaultBorrowValueBeside,
			},
			{
				ID:          VaultChecks.Borrow.Pool(),
				Statement:   "Return takes back what Borrow lent",
				Claim:       spec.Pool,
				Probe:       spec.TakesBackWhatItLent,
				Reason:      spec.UnreturnedLoan,
				Law:         spec.BorrowedValueIsTakenBack,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.TakesBackWhatItLent(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Borrow(ctx)
					}, func(ctx context.Context, v Item) error {
						return h.Subject().Return(ctx, v)
					})
				},
				Defect: vaultReturnPutRefuses,
			},
			{
				ID:          VaultChecks.Return.Smoke(),
				Statement:   "Return returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Return(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: vaultReturnPanics,
			},
			{
				ID:          VaultChecks.Return.SmokeGenerated(),
				Statement:   "Return returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Return(ctx, drawnIt)
					})
				},
				Defect: vaultReturnPanics,
			},
			{
				ID:          VaultChecks.Return.Cancel(),
				Statement:   "Return returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Return(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: vaultReturnIgnoresCancel,
			},
			{
				ID:          VaultChecks.Return.CancelGenerated(),
				Statement:   "Return returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VaultHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Return(ctx, drawnIt)
					})
				},
				Defect: vaultReturnIgnoresCancel,
			},
			{
				ID:          VaultChecks.Stats.Smoke(),
				Statement:   "Stats returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VaultHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Stats(ctx)
						return err
					})
				},
				Defect: vaultStatsPanics,
			},
			{
				ID:          VaultChecks.Stats.Cancel(),
				Statement:   "Stats returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Stats(ctx)
						return err
					})
				},
				Defect: vaultStatsIgnoresCancel,
			},
			{
				ID:          VaultChecks.Stats.ZeroOnError(),
				Statement:   "Stats returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Usage, error) {
						return h.Subject().Stats(ctx)
					}, Usage{})
				},
				Defect: vaultStatsValueBeside,
			},
			{
				ID:          VaultChecks.Stats.CountsEachLoan(),
				Statement:   "Stats reports another count after Borrow lends an object",
				Claim:       spec.CountsEachLoan,
				Probe:       spec.MovesWithALoan,
				Reason:      spec.UncountedLoan,
				Law:         spec.StatsMoveWithALoan,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VaultHandle) {
					probe.MovesWithALoan(h, func(ctx context.Context) (Usage, error) {
						return h.Subject().Stats(ctx)
					}, func(ctx context.Context) error {
						_, err := h.Subject().Borrow(ctx)
						return err
					})
				},
				Defect: vaultStatsFrozenStats,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Vault.Append",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Vault.Verify",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Vault.Acquire",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Vault.Release",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Vault.Return",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Vault]{
			{
				Calls:    vaultAppendCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Append.Callable,
			},
			{
				Calls:    vaultReplayCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Replay.Callable,
			},
			{
				Calls:    vaultVerifyCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Verify.Callable,
			},
			{
				Calls:    vaultAcquireCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Acquire.Callable,
			},
			{
				Calls:    vaultReleaseCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Release.Callable,
			},
			{
				Calls:    vaultBorrowCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Borrow.Callable,
			},
			{
				Calls:    vaultReturnCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Return.Callable,
			},
			{
				Calls:    vaultStatsCalls,
				Over:     newVaultOver,
				Callable: VaultChecks.Stats.Callable,
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

// vaultAppendCalls generates the calls of Vault.Append with the inputs
// that each case draws.
var vaultAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Vault) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Vault) error {
		return subject.Append(ctx, it)
	}
})

// vaultReplayCalls generates the calls of Vault.Replay with the inputs
// that each case draws.
var vaultReplayCalls = prop.Composite(func(*prop.Case) func(context.Context, Vault) error {
	return func(ctx context.Context, subject Vault) error {
		_, err := subject.Replay(ctx)
		return err
	}
})

// vaultVerifyCalls generates the calls of Vault.Verify with the inputs
// that each case draws.
var vaultVerifyCalls = prop.Composite(func(*prop.Case) func(context.Context, Vault) error {
	return func(ctx context.Context, subject Vault) error {
		return subject.Verify(ctx)
	}
})

// vaultAcquireCalls generates the calls of Vault.Acquire with the inputs
// that each case draws.
var vaultAcquireCalls = prop.Composite(func(c *prop.Case) func(context.Context, Vault) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Vault) error {
		return subject.Acquire(ctx, key)
	}
})

// vaultReleaseCalls generates the calls of Vault.Release with the inputs
// that each case draws.
var vaultReleaseCalls = prop.Composite(func(c *prop.Case) func(context.Context, Vault) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Vault) error {
		return subject.Release(ctx, key)
	}
})

// vaultBorrowCalls generates the calls of Vault.Borrow with the inputs
// that each case draws.
var vaultBorrowCalls = prop.Composite(func(*prop.Case) func(context.Context, Vault) error {
	return func(ctx context.Context, subject Vault) error {
		_, err := subject.Borrow(ctx)
		return err
	}
})

// vaultReturnCalls generates the calls of Vault.Return with the inputs
// that each case draws.
var vaultReturnCalls = prop.Composite(func(c *prop.Case) func(context.Context, Vault) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Vault) error {
		return subject.Return(ctx, it)
	}
})

// vaultStatsCalls generates the calls of Vault.Stats with the inputs that
// each case draws.
var vaultStatsCalls = prop.Composite(func(*prop.Case) func(context.Context, Vault) error {
	return func(ctx context.Context, subject Vault) error {
		_, err := subject.Stats(ctx)
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

// Append calls spend and returns its zero results.
func (o vaultOver) Append(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Replay calls spend and returns its zero results.
func (o vaultOver) Replay(context.Context) (r0 []Item, r1 error) {
	o.spend()
	return r0, r1
}

// Verify calls spend and returns its zero results.
func (o vaultOver) Verify(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Acquire calls spend and returns its zero results.
func (o vaultOver) Acquire(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Release calls spend and returns its zero results.
func (o vaultOver) Release(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Borrow calls spend and returns its zero results.
func (o vaultOver) Borrow(context.Context) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Return calls spend and returns its zero results.
func (o vaultOver) Return(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Stats calls spend and returns its zero results.
func (o vaultOver) Stats(context.Context) (r0 Usage, r1 error) {
	o.spend()
	return r0, r1
}

// newVaultOver returns the planted defect of a ceiling of Vault, whose
// callables call spend.
func newVaultOver(_ *VaultHandle, spend func()) Vault {
	return vaultOver{spend: spend}
}

// vaultAppendPanicsDefect is the planted defect of a Vault whose Append
// panics.
type vaultAppendPanicsDefect struct {
	vaultOverEmbedded
}

// Append is the override of the planted defect.
func (*vaultAppendPanicsDefect) Append(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Vault whose Append panics")
}

// vaultAppendPanics returns the planted defect of a Vault whose Append
// panics over a fresh subject.
func vaultAppendPanics(h *VaultHandle) Vault {
	return &vaultAppendPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultAppendIgnoresCancelDefect is the planted defect of a Vault whose
// Append ignores the context.
type vaultAppendIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Append is the override of the planted defect.
func (d *vaultAppendIgnoresCancelDefect) Append(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.vaultOverEmbedded.Append(p0, p1)
}

// vaultAppendIgnoresCancel returns the planted defect of a Vault whose
// Append ignores the context over a fresh subject.
func vaultAppendIgnoresCancel(h *VaultHandle) Vault {
	return &vaultAppendIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultAppendSucceedsAlwaysDefect is the planted defect of a Vault whose
// Append returns success and appends nothing.
type vaultAppendSucceedsAlwaysDefect struct {
	vaultOverEmbedded
}

// Append is the override of the planted defect.
func (*vaultAppendSucceedsAlwaysDefect) Append(p0 context.Context, p1 Item) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// vaultAppendSucceedsAlways returns the planted defect of a Vault whose
// Append returns success and appends nothing over a fresh subject.
func vaultAppendSucceedsAlways(h *VaultHandle) Vault {
	return &vaultAppendSucceedsAlwaysDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReplayPanicsDefect is the planted defect of a Vault whose Replay
// panics.
type vaultReplayPanicsDefect struct {
	vaultOverEmbedded
}

// Replay is the override of the planted defect.
func (*vaultReplayPanicsDefect) Replay(p0 context.Context) (r0 []Item, r1 error) {
	panic("the planted defect of a Vault whose Replay panics")
}

// vaultReplayPanics returns the planted defect of a Vault whose Replay
// panics over a fresh subject.
func vaultReplayPanics(h *VaultHandle) Vault {
	return &vaultReplayPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReplayIgnoresCancelDefect is the planted defect of a Vault whose
// Replay ignores the context.
type vaultReplayIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Replay is the override of the planted defect.
func (d *vaultReplayIgnoresCancelDefect) Replay(p0 context.Context) (r0 []Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.vaultOverEmbedded.Replay(p0)
}

// vaultReplayIgnoresCancel returns the planted defect of a Vault whose
// Replay ignores the context over a fresh subject.
func vaultReplayIgnoresCancel(h *VaultHandle) Vault {
	return &vaultReplayIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReplayValueBesideDefect is the planted defect of a Vault whose
// Replay returns a value beside an error.
type vaultReplayValueBesideDefect struct {
	vaultOverEmbedded
}

// Replay is the override of the planted defect.
func (d *vaultReplayValueBesideDefect) Replay(p0 context.Context) (r0 []Item, r1 error) {
	r0, r1 = d.vaultOverEmbedded.Replay(p0)
	if r1 != nil {
		r0 = []Item{{Key: "Key-b"}}
	}
	return r0, r1
}

// vaultReplayValueBeside returns the planted defect of a Vault whose
// Replay returns a value beside an error over a fresh subject.
func vaultReplayValueBeside(h *VaultHandle) Vault {
	return &vaultReplayValueBesideDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultVerifyPanicsDefect is the planted defect of a Vault whose Verify
// panics.
type vaultVerifyPanicsDefect struct {
	vaultOverEmbedded
}

// Verify is the override of the planted defect.
func (*vaultVerifyPanicsDefect) Verify(p0 context.Context) (r0 error) {
	panic("the planted defect of a Vault whose Verify panics")
}

// vaultVerifyPanics returns the planted defect of a Vault whose Verify
// panics over a fresh subject.
func vaultVerifyPanics(h *VaultHandle) Vault {
	return &vaultVerifyPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultVerifyIgnoresCancelDefect is the planted defect of a Vault whose
// Verify ignores the context.
type vaultVerifyIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Verify is the override of the planted defect.
func (d *vaultVerifyIgnoresCancelDefect) Verify(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.vaultOverEmbedded.Verify(p0)
}

// vaultVerifyIgnoresCancel returns the planted defect of a Vault whose
// Verify ignores the context over a fresh subject.
func vaultVerifyIgnoresCancel(h *VaultHandle) Vault {
	return &vaultVerifyIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultVerifyReportsTamperDefect is the planted defect of a Vault whose
// Verify reports a change of a chain that Append produced.
type vaultVerifyReportsTamperDefect struct {
	vaultOverEmbedded
}

// Verify is the override of the planted defect.
func (*vaultVerifyReportsTamperDefect) Verify(p0 context.Context) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	r0 = conform.ErrGenerated
	return r0
}

// vaultVerifyReportsTamper returns the planted defect of a Vault whose
// Verify reports a change of a chain that Append produced over a fresh
// subject.
func vaultVerifyReportsTamper(h *VaultHandle) Vault {
	return &vaultVerifyReportsTamperDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultAcquirePanicsDefect is the planted defect of a Vault whose Acquire
// panics.
type vaultAcquirePanicsDefect struct {
	vaultOverEmbedded
}

// Acquire is the override of the planted defect.
func (*vaultAcquirePanicsDefect) Acquire(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Vault whose Acquire panics")
}

// vaultAcquirePanics returns the planted defect of a Vault whose Acquire
// panics over a fresh subject.
func vaultAcquirePanics(h *VaultHandle) Vault {
	return &vaultAcquirePanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultAcquireIgnoresCancelDefect is the planted defect of a Vault whose
// Acquire ignores the context.
type vaultAcquireIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Acquire is the override of the planted defect.
func (d *vaultAcquireIgnoresCancelDefect) Acquire(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.vaultOverEmbedded.Acquire(p0, p1)
}

// vaultAcquireIgnoresCancel returns the planted defect of a Vault whose
// Acquire ignores the context over a fresh subject.
func vaultAcquireIgnoresCancel(h *VaultHandle) Vault {
	return &vaultAcquireIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultAcquireSucceedsAlwaysDefect is the planted defect of a Vault whose
// Acquire grants a lease that it holds.
type vaultAcquireSucceedsAlwaysDefect struct {
	vaultOverEmbedded
}

// Acquire is the override of the planted defect.
func (*vaultAcquireSucceedsAlwaysDefect) Acquire(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// vaultAcquireSucceedsAlways returns the planted defect of a Vault whose
// Acquire grants a lease that it holds over a fresh subject.
func vaultAcquireSucceedsAlways(h *VaultHandle) Vault {
	return &vaultAcquireSucceedsAlwaysDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReleaseReleaseKeepsTheLeaseDefect is the planted defect of a Vault
// whose Release keeps the lease.
type vaultReleaseReleaseKeepsTheLeaseDefect struct {
	vaultOverEmbedded
}

// Release is the override of the planted defect.
func (*vaultReleaseReleaseKeepsTheLeaseDefect) Release(p0 context.Context, p1 string) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// vaultReleaseReleaseKeepsTheLease returns the planted defect of a Vault
// whose Release keeps the lease over a fresh subject.
func vaultReleaseReleaseKeepsTheLease(h *VaultHandle) Vault {
	return &vaultReleaseReleaseKeepsTheLeaseDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReleasePanicsDefect is the planted defect of a Vault whose Release
// panics.
type vaultReleasePanicsDefect struct {
	vaultOverEmbedded
}

// Release is the override of the planted defect.
func (*vaultReleasePanicsDefect) Release(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Vault whose Release panics")
}

// vaultReleasePanics returns the planted defect of a Vault whose Release
// panics over a fresh subject.
func vaultReleasePanics(h *VaultHandle) Vault {
	return &vaultReleasePanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReleaseIgnoresCancelDefect is the planted defect of a Vault whose
// Release ignores the context.
type vaultReleaseIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Release is the override of the planted defect.
func (d *vaultReleaseIgnoresCancelDefect) Release(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.vaultOverEmbedded.Release(p0, p1)
}

// vaultReleaseIgnoresCancel returns the planted defect of a Vault whose
// Release ignores the context over a fresh subject.
func vaultReleaseIgnoresCancel(h *VaultHandle) Vault {
	return &vaultReleaseIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultBorrowPanicsDefect is the planted defect of a Vault whose Borrow
// panics.
type vaultBorrowPanicsDefect struct {
	vaultOverEmbedded
}

// Borrow is the override of the planted defect.
func (*vaultBorrowPanicsDefect) Borrow(p0 context.Context) (r0 Item, r1 error) {
	panic("the planted defect of a Vault whose Borrow panics")
}

// vaultBorrowPanics returns the planted defect of a Vault whose Borrow
// panics over a fresh subject.
func vaultBorrowPanics(h *VaultHandle) Vault {
	return &vaultBorrowPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultBorrowIgnoresCancelDefect is the planted defect of a Vault whose
// Borrow ignores the context.
type vaultBorrowIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Borrow is the override of the planted defect.
func (d *vaultBorrowIgnoresCancelDefect) Borrow(p0 context.Context) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.vaultOverEmbedded.Borrow(p0)
}

// vaultBorrowIgnoresCancel returns the planted defect of a Vault whose
// Borrow ignores the context over a fresh subject.
func vaultBorrowIgnoresCancel(h *VaultHandle) Vault {
	return &vaultBorrowIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultBorrowValueBesideDefect is the planted defect of a Vault whose
// Borrow returns a value beside an error.
type vaultBorrowValueBesideDefect struct {
	vaultOverEmbedded
}

// Borrow is the override of the planted defect.
func (d *vaultBorrowValueBesideDefect) Borrow(p0 context.Context) (r0 Item, r1 error) {
	r0, r1 = d.vaultOverEmbedded.Borrow(p0)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// vaultBorrowValueBeside returns the planted defect of a Vault whose
// Borrow returns a value beside an error over a fresh subject.
func vaultBorrowValueBeside(h *VaultHandle) Vault {
	return &vaultBorrowValueBesideDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReturnPutRefusesDefect is the planted defect of a Vault whose
// Return refuses what Borrow lent.
type vaultReturnPutRefusesDefect struct {
	vaultOverEmbedded
}

// Return is the override of the planted defect.
func (*vaultReturnPutRefusesDefect) Return(p0 context.Context, p1 Item) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	r0 = conform.ErrGenerated
	return r0
}

// vaultReturnPutRefuses returns the planted defect of a Vault whose Return
// refuses what Borrow lent over a fresh subject.
func vaultReturnPutRefuses(h *VaultHandle) Vault {
	return &vaultReturnPutRefusesDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReturnPanicsDefect is the planted defect of a Vault whose Return
// panics.
type vaultReturnPanicsDefect struct {
	vaultOverEmbedded
}

// Return is the override of the planted defect.
func (*vaultReturnPanicsDefect) Return(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Vault whose Return panics")
}

// vaultReturnPanics returns the planted defect of a Vault whose Return
// panics over a fresh subject.
func vaultReturnPanics(h *VaultHandle) Vault {
	return &vaultReturnPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultReturnIgnoresCancelDefect is the planted defect of a Vault whose
// Return ignores the context.
type vaultReturnIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Return is the override of the planted defect.
func (d *vaultReturnIgnoresCancelDefect) Return(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.vaultOverEmbedded.Return(p0, p1)
}

// vaultReturnIgnoresCancel returns the planted defect of a Vault whose
// Return ignores the context over a fresh subject.
func vaultReturnIgnoresCancel(h *VaultHandle) Vault {
	return &vaultReturnIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultStatsPanicsDefect is the planted defect of a Vault whose Stats
// panics.
type vaultStatsPanicsDefect struct {
	vaultOverEmbedded
}

// Stats is the override of the planted defect.
func (*vaultStatsPanicsDefect) Stats(p0 context.Context) (r0 Usage, r1 error) {
	panic("the planted defect of a Vault whose Stats panics")
}

// vaultStatsPanics returns the planted defect of a Vault whose Stats
// panics over a fresh subject.
func vaultStatsPanics(h *VaultHandle) Vault {
	return &vaultStatsPanicsDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultStatsIgnoresCancelDefect is the planted defect of a Vault whose
// Stats ignores the context.
type vaultStatsIgnoresCancelDefect struct {
	vaultOverEmbedded
}

// Stats is the override of the planted defect.
func (d *vaultStatsIgnoresCancelDefect) Stats(p0 context.Context) (r0 Usage, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.vaultOverEmbedded.Stats(p0)
}

// vaultStatsIgnoresCancel returns the planted defect of a Vault whose
// Stats ignores the context over a fresh subject.
func vaultStatsIgnoresCancel(h *VaultHandle) Vault {
	return &vaultStatsIgnoresCancelDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultStatsValueBesideDefect is the planted defect of a Vault whose Stats
// returns a value beside an error.
type vaultStatsValueBesideDefect struct {
	vaultOverEmbedded
}

// Stats is the override of the planted defect.
func (d *vaultStatsValueBesideDefect) Stats(p0 context.Context) (r0 Usage, r1 error) {
	r0, r1 = d.vaultOverEmbedded.Stats(p0)
	if r1 != nil {
		r0 = Usage{Lent: 7}
	}
	return r0, r1
}

// vaultStatsValueBeside returns the planted defect of a Vault whose Stats
// returns a value beside an error over a fresh subject.
func vaultStatsValueBeside(h *VaultHandle) Vault {
	return &vaultStatsValueBesideDefect{vaultOverEmbedded: h.Fresh()}
}

// vaultStatsFrozenStatsDefect is the planted defect of a Vault whose Stats
// reports the same stats after Borrow lends an object.
type vaultStatsFrozenStatsDefect struct {
	vaultOverEmbedded
}

// Stats is the override of the planted defect.
func (*vaultStatsFrozenStatsDefect) Stats(p0 context.Context) (r0 Usage, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	return r0, r1
}

// vaultStatsFrozenStats returns the planted defect of a Vault whose Stats
// reports the same stats after Borrow lends an object over a fresh
// subject.
func vaultStatsFrozenStats(h *VaultHandle) Vault {
	return &vaultStatsFrozenStatsDefect{vaultOverEmbedded: h.Fresh()}
}
