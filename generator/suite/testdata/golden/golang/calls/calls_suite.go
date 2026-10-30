package calls

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"example.com/shop/store"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Board:
//
// Board.Place/smoke: Place returns without a panic for derived inputs.
//
// Board.Clash/smoke: Clash returns without a panic for derived inputs.
//
// Board.Clash/cancel: Clash returns an error for a cancelled context.
//
// Board.Blank/smoke: Blank returns without a panic for derived inputs.
//
// Board.Index/smoke: Index returns without a panic for derived inputs.
//
// Board.Index/deterministic: Index returns the same value for the same
// inputs.
//
// Board.Label/smoke: Label returns without a panic for derived inputs.
//
// Board.Wait/smoke: Wait returns without a panic for derived inputs.
//
// Board.Hold/smoke: Hold returns without a panic for derived inputs.
//
// Board.Swap/smoke: Swap returns without a panic for derived inputs.
//
// Board.Swap/deterministic: Swap returns the same value for the same
// inputs.
//
// Board.Mark/smoke: Mark returns without a panic for derived inputs.
//
// Board.Paint/smoke: Paint returns without a panic for derived inputs.
//
// Board.Repeat/smoke: Repeat returns without a panic for derived inputs.
//
// Board.Keep/smoke: Keep returns without a panic for derived inputs.
//
// Board.Stock/smoke: Stock returns without a panic for derived inputs.
//
// Board.Discard/smoke: Discard returns without a panic for derived inputs.
//
// Board.Twice/smoke: Twice returns without a panic for derived inputs.
//
// Board.Gather/smoke: Gather returns without a panic for derived inputs.

// The derived checks of Board that witness their laws at the shallow end:
//
// Board.Place/smoke
//
// Board.Clash/smoke
//
// Board.Blank/smoke
//
// Board.Index/smoke
//
// Board.Index/deterministic
//
// Board.Label/smoke
//
// Board.Wait/smoke
//
// Board.Hold/smoke
//
// Board.Swap/smoke
//
// Board.Swap/deterministic
//
// Board.Mark/smoke
//
// Board.Paint/smoke
//
// Board.Repeat/smoke
//
// Board.Keep/smoke
//
// Board.Stock/smoke
//
// Board.Discard/smoke
//
// Board.Twice/smoke
//
// Board.Gather/smoke

// The unclaimed classifications of Board:
//
// The shape:writer of Board.Place has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:multi-arg-writer of Board.Clash has the stance silent. A write
// is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Board.Label has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Board.Mark has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Board.Paint has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Board.Repeat has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Board.Discard has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Board writes.
type (
	BoardHandle = conform.Handle[Board]
	BoardCheck  = conform.Check[Board]
	BoardOption = conform.Option[Board]
)

// BoardChecks contains the name of each callable of Board, with the
// identities of its checks.
var BoardChecks = struct {
	Place   boardPlace
	Clash   boardClash
	Blank   boardBlank
	Index   boardIndex
	Label   boardLabel
	Wait    boardWait
	Hold    boardHold
	Swap    boardSwap
	Mark    boardMark
	Paint   boardPaint
	Repeat  boardRepeat
	Keep    boardKeep
	Stock   boardStock
	Discard boardDiscard
	Twice   boardTwice
	Drain   boardDrain
	Gather  boardGather
}{
	Place:   boardPlace{"Board.Place"},
	Clash:   boardClash{"Board.Clash"},
	Blank:   boardBlank{"Board.Blank"},
	Index:   boardIndex{"Board.Index"},
	Label:   boardLabel{"Board.Label"},
	Wait:    boardWait{"Board.Wait"},
	Hold:    boardHold{"Board.Hold"},
	Swap:    boardSwap{"Board.Swap"},
	Mark:    boardMark{"Board.Mark"},
	Paint:   boardPaint{"Board.Paint"},
	Repeat:  boardRepeat{"Board.Repeat"},
	Keep:    boardKeep{"Board.Keep"},
	Stock:   boardStock{"Board.Stock"},
	Discard: boardDiscard{"Board.Discard"},
	Twice:   boardTwice{"Board.Twice"},
	Drain:   boardDrain{"Board.Drain"},
	Gather:  boardGather{"Board.Gather"},
}

// boardPlace is the name of Board.Place, with the identities of its
// checks.
type boardPlace struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Place returns without a
// panic for derived inputs.
func (c boardPlace) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardClash is the name of Board.Clash, with the identities of its
// checks.
type boardClash struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Clash returns without a
// panic for derived inputs.
func (c boardClash) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Clash returns an error for
// a cancelled context.
func (c boardClash) Cancel() conform.ID[Board] {
	return c.Check("cancel")
}

// boardBlank is the name of Board.Blank, with the identities of its
// checks.
type boardBlank struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Blank returns without a
// panic for derived inputs.
func (c boardBlank) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardIndex is the name of Board.Index, with the identities of its
// checks.
type boardIndex struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Index returns without a
// panic for derived inputs.
func (c boardIndex) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Index returns the
// same value for the same inputs.
func (c boardIndex) Deterministic() conform.ID[Board] {
	return c.Check("deterministic")
}

// boardLabel is the name of Board.Label, with the identities of its
// checks.
type boardLabel struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Label returns without a
// panic for derived inputs.
func (c boardLabel) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardWait is the name of Board.Wait, with the identities of its checks.
type boardWait struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Wait returns without a
// panic for derived inputs.
func (c boardWait) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardHold is the name of Board.Hold, with the identities of its checks.
type boardHold struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Hold returns without a
// panic for derived inputs.
func (c boardHold) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardSwap is the name of Board.Swap, with the identities of its checks.
type boardSwap struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Swap returns without a
// panic for derived inputs.
func (c boardSwap) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Swap returns the
// same value for the same inputs.
func (c boardSwap) Deterministic() conform.ID[Board] {
	return c.Check("deterministic")
}

// boardMark is the name of Board.Mark, with the identities of its checks.
type boardMark struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Mark returns without a
// panic for derived inputs.
func (c boardMark) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardPaint is the name of Board.Paint, with the identities of its
// checks.
type boardPaint struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Paint returns without a
// panic for derived inputs.
func (c boardPaint) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardRepeat is the name of Board.Repeat, with the identities of its
// checks.
type boardRepeat struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Repeat returns without a
// panic for derived inputs.
func (c boardRepeat) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardKeep is the name of Board.Keep, with the identities of its checks.
type boardKeep struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Keep returns without a
// panic for derived inputs.
func (c boardKeep) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardStock is the name of Board.Stock, with the identities of its
// checks.
type boardStock struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Stock returns without a
// panic for derived inputs.
func (c boardStock) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardDiscard is the name of Board.Discard, with the identities of its
// checks.
type boardDiscard struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Discard returns without a
// panic for derived inputs.
func (c boardDiscard) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardTwice is the name of Board.Twice, with the identities of its
// checks.
type boardTwice struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Twice returns without a
// panic for derived inputs.
func (c boardTwice) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// boardDrain is the name of Board.Drain, with the identities of its
// checks.
type boardDrain struct {
	conform.Callable[Board]
}

// boardGather is the name of Board.Gather, with the identities of its
// checks.
type boardGather struct {
	conform.Callable[Board]
}

// Smoke returns the identity of the check that Gather returns without a
// panic for derived inputs.
func (c boardGather) Smoke() conform.ID[Board] {
	return c.Check("smoke")
}

// BoardSuite returns the checks of Board, and the calls and the budgets of
// its callables.
func BoardSuite() conform.Suite[Board] {
	return conform.Suite[Board]{
		Declaration: "Board",
		Package:     "example.com/shop/calls",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Board]{
			{
				ID:          BoardChecks.Place.Smoke(),
				Statement:   "Place returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Place(&Point{X: 42})
					})
				},
				Defect: boardPlacePanics,
			},
			{
				ID:          BoardChecks.Clash.Smoke(),
				Statement:   "Clash returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Clash(ctx, string("c-a"), string("subject-a"), string("err-a"), int(42))
					})
				},
				Defect: boardClashPanics,
			},
			{
				ID:          BoardChecks.Clash.Cancel(),
				Statement:   "Clash returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *BoardHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Clash(ctx, string("c-a"), string("subject-a"), string("err-a"), int(42))
					})
				},
				Defect: boardClashIgnoresCancel,
			},
			{
				ID:          BoardChecks.Blank.Smoke(),
				Statement:   "Blank returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Blank(string("sample-a"), int(42))
						return nil
					})
				},
				Defect: boardBlankPanics,
			},
			{
				ID:          BoardChecks.Index.Smoke(),
				Statement:   "Index returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Index(map[Shape][]Point{Shape(1): {{X: 42}}})
						return nil
					})
				},
				Defect: boardIndexPanics,
			},
			{
				ID:          BoardChecks.Index.Deterministic(),
				Statement:   "Index returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Board]{BoardChecks.Index.Smoke()},
				Apply: func(h *BoardHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Index(map[Shape][]Point{Shape(1): {{X: 42}}})
						return r0, nil
					})
				},
				Defect: boardIndexDrifts,
			},
			{
				ID:          BoardChecks.Label.Smoke(),
				Statement:   "Label returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Label(string("string-a"))
						return nil
					})
				},
				Defect: boardLabelPanics,
			},
			{
				ID:          BoardChecks.Wait.Smoke(),
				Statement:   "Wait returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Wait(int(42), []time.Duration{time.Duration(1)})
						return nil
					})
				},
				Defect: boardWaitPanics,
			},
			{
				ID:          BoardChecks.Hold.Smoke(),
				Statement:   "Hold returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Hold(int(42), Pair[time.Month]{First: time.Month(1)})
						return nil
					})
				},
				Defect: boardHoldPanics,
			},
			{
				ID:          BoardChecks.Swap.Smoke(),
				Statement:   "Swap returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Swap(Pair[Tint]{First: Tint(1)})
						return nil
					})
				},
				Defect: boardSwapPanics,
			},
			{
				ID:          BoardChecks.Swap.Deterministic(),
				Statement:   "Swap returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Board]{BoardChecks.Swap.Smoke()},
				Apply: func(h *BoardHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (Pair[Tint], error) {
						r0 := h.Subject().Swap(Pair[Tint]{First: Tint(1)})
						return r0, nil
					})
				},
				Defect: boardSwapDrifts,
			},
			{
				ID:          BoardChecks.Mark.Smoke(),
				Statement:   "Mark returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Mark(Cell{Note: "Note-a"})
					})
				},
				Defect: boardMarkPanics,
			},
			{
				ID:          BoardChecks.Paint.Smoke(),
				Statement:   "Paint returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Paint(Palette{Hue(1)})
						return nil
					})
				},
				Defect: boardPaintPanics,
			},
			{
				ID:          BoardChecks.Repeat.Smoke(),
				Statement:   "Repeat returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Repeat(int(42))
						return nil
					})
				},
				Defect: boardRepeatPanics,
			},
			{
				ID:          BoardChecks.Keep.Smoke(),
				Statement:   "Keep returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Keep(int(42), store.Item{ID: "ID-a"})
						return nil
					})
				},
				Defect: boardKeepPanics,
			},
			{
				ID:          BoardChecks.Stock.Smoke(),
				Statement:   "Stock returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Stock(int(42), store.Item{ID: "ID-a"})
						return nil
					})
				},
				Defect: boardStockPanics,
			},
			{
				ID:          BoardChecks.Discard.Smoke(),
				Statement:   "Discard returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Discard(string("_-a"))
						return nil
					})
				},
				Defect: boardDiscardPanics,
			},
			{
				ID:          BoardChecks.Twice.Smoke(),
				Statement:   "Twice returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Twice(string("c-a"), string("c0-a"))
						return nil
					})
				},
				Defect: boardTwicePanics,
			},
			{
				ID:          BoardChecks.Gather.Smoke(),
				Statement:   "Gather returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *BoardHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Gather(func() int { return 0 }, int(42), int(42))
						return nil
					})
				},
				Defect: boardGatherPanics,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Board.Place",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:multi-arg-writer",
				Callable:       "Board.Clash",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Board.Label",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Board.Mark",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Board.Paint",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Board.Repeat",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Board.Discard",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Board]{
			{
				Calls:    boardPlaceCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Place.Callable,
			},
			{
				Calls:    boardClashCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Clash.Callable,
			},
			{
				Calls:    boardBlankCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Blank.Callable,
			},
			{
				Calls:    boardIndexCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Index.Callable,
			},
			{
				Calls:    boardLabelCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Label.Callable,
			},
			{
				Calls:    boardWaitCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Wait.Callable,
			},
			{
				Calls:    boardHoldCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Hold.Callable,
			},
			{
				Calls:    boardSwapCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Swap.Callable,
			},
			{
				Calls:    boardMarkCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Mark.Callable,
			},
			{
				Calls:    boardPaintCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Paint.Callable,
			},
			{
				Calls:    boardRepeatCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Repeat.Callable,
			},
			{
				Calls:    boardKeepCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Keep.Callable,
			},
			{
				Calls:    boardStockCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Stock.Callable,
			},
			{
				Calls:    boardDiscardCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Discard.Callable,
			},
			{
				Calls:    boardTwiceCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Twice.Callable,
			},
			{
				Calls:    boardDrainCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Drain.Callable,
			},
			{
				Calls:    boardGatherCalls,
				Over:     newBoardOver,
				Callable: BoardChecks.Gather.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.SampledFrom(Round, Square)),
			prop.Using(prop.SampledFrom(Light, Dark)),
			prop.Using(prop.SampledFrom(Red, Blue)),
			prop.Using(prop.SampledFrom(Warm, Cool)),
		},
	}
}

// RunBoardSuite applies the checks of Board, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunBoardSuite[T Board](t *testing.T, newSubject func() T, opts ...BoardOption) {
	t.Helper()
	RunBoardSuiteFromHandle(t, func(*BoardHandle) T { return newSubject() }, opts...)
}

// RunBoardSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunBoardSuiteFromHandle[T Board](
	t *testing.T,
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	t.Helper()
	conform.Run(t, BoardSuite(), boardSubjects(newSubject), opts...)
}

// FuzzBoardSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzBoardSuite[T Board](
	f *testing.F,
	callable conform.Callable[Board],
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	f.Helper()
	conform.Fuzz(f, BoardSuite(), callable, boardSubjects(newSubject), opts...)
}

// RunBoardBudgets applies the allocation ceilings of Board and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunBoardBudgets[T Board](
	t *testing.T,
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	t.Helper()
	conform.RunBudgets(t, BoardSuite(), boardSubjects(newSubject), opts...)
}

// BenchBoardSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchBoardSuite[T Board](
	b conform.B,
	callable conform.Callable[Board],
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	b.Helper()
	conform.Bench(b, BoardSuite(), callable, boardSubjects(newSubject), opts...)
}

// boardSubjects returns the factory of the subjects of Board that calls
// newSubject with the handle of each check.
func boardSubjects[T Board](newSubject func(h *BoardHandle) T) func(h *BoardHandle) Board {
	return func(h *BoardHandle) Board { return newSubject(h) }
}

// boardPlaceCalls generates the calls of Board.Place with the inputs that
// each case draws.
var boardPlaceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	p := c.Draw(prop.Of[Point](), "p")
	return func(_ context.Context, subject Board) error {
		return subject.Place(&p)
	}
})

// boardClashCalls generates the calls of Board.Clash with the inputs that
// each case draws.
var boardClashCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	c0 := c.Draw(prop.Of[string](), "c")
	subject0 := c.Draw(prop.Of[string](), "subject")
	err0 := c.Draw(prop.Of[string](), "err")
	prop0 := c.Draw(prop.Of[int](), "prop")
	return func(ctx context.Context, subject Board) error {
		return subject.Clash(ctx, c0, subject0, err0, prop0)
	}
})

// boardBlankCalls generates the calls of Board.Blank with the inputs that
// each case draws.
var boardBlankCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	a0 := c.Draw(prop.Of[string](), "a0")
	a1 := c.Draw(prop.Of[int](), "a1")
	return func(_ context.Context, subject Board) error {
		subject.Blank(a0, a1)
		return nil
	}
})

// boardIndexCalls generates the calls of Board.Index with the inputs that
// each case draws.
var boardIndexCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	points := c.Draw(prop.Of[map[Shape][]Point](), "points")
	return func(_ context.Context, subject Board) error {
		_ = subject.Index(points)
		return nil
	}
})

// boardLabelCalls generates the calls of Board.Label with the inputs that
// each case draws.
var boardLabelCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	string0 := c.Draw(prop.Of[string](), "string")
	return func(_ context.Context, subject Board) error {
		subject.Label(string0)
		return nil
	}
})

// boardWaitCalls generates the calls of Board.Wait with the inputs that
// each case draws.
var boardWaitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	time0 := c.Draw(prop.Of[int](), "time")
	delays := c.Draw(prop.Of[[]time.Duration](), "delays")
	return func(_ context.Context, subject Board) error {
		subject.Wait(time0, delays)
		return nil
	}
})

// boardHoldCalls generates the calls of Board.Hold with the inputs that
// each case draws.
var boardHoldCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	time0 := c.Draw(prop.Of[int](), "time")
	pair := c.Draw(prop.Of[Pair[time.Month]](), "pair")
	return func(_ context.Context, subject Board) error {
		subject.Hold(time0, pair)
		return nil
	}
})

// boardSwapCalls generates the calls of Board.Swap with the inputs that
// each case draws.
var boardSwapCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	pair := c.Draw(prop.Of[Pair[Tint]](), "pair")
	return func(_ context.Context, subject Board) error {
		_ = subject.Swap(pair)
		return nil
	}
})

// boardMarkCalls generates the calls of Board.Mark with the inputs that
// each case draws.
var boardMarkCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	cell := c.Draw(prop.Of[Cell](), "cell")
	return func(_ context.Context, subject Board) error {
		return subject.Mark(cell)
	}
})

// boardPaintCalls generates the calls of Board.Paint with the inputs that
// each case draws.
var boardPaintCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	palette := c.Draw(prop.Of[Palette](), "palette")
	return func(_ context.Context, subject Board) error {
		subject.Paint(palette)
		return nil
	}
})

// boardRepeatCalls generates the calls of Board.Repeat with the inputs
// that each case draws.
var boardRepeatCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	calls0 := c.Draw(prop.Of[int](), "calls")
	return func(_ context.Context, subject Board) error {
		subject.Repeat(calls0)
		return nil
	}
})

// boardKeepCalls generates the calls of Board.Keep with the inputs that
// each case draws.
var boardKeepCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	store0 := c.Draw(prop.Of[int](), "store")
	item := c.Draw(prop.Of[store.Item](), "item")
	return func(_ context.Context, subject Board) error {
		subject.Keep(store0, item)
		return nil
	}
})

// boardStockCalls generates the calls of Board.Stock with the inputs that
// each case draws.
var boardStockCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	shop0 := c.Draw(prop.Of[int](), "shop")
	item := c.Draw(prop.Of[store.Item](), "item")
	return func(_ context.Context, subject Board) error {
		subject.Stock(shop0, item)
		return nil
	}
})

// boardDiscardCalls generates the calls of Board.Discard with the inputs
// that each case draws.
var boardDiscardCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	a0 := c.Draw(prop.Of[string](), "a0")
	return func(_ context.Context, subject Board) error {
		subject.Discard(a0)
		return nil
	}
})

// boardTwiceCalls generates the calls of Board.Twice with the inputs that
// each case draws.
var boardTwiceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	c0 := c.Draw(prop.Of[string](), "c")
	c00 := c.Draw(prop.Of[string](), "c0")
	return func(_ context.Context, subject Board) error {
		subject.Twice(c0, c00)
		return nil
	}
})

// boardDrainCalls generates the calls of Board.Drain with the inputs that
// each case draws.
var boardDrainCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	intsValues := c.Draw(prop.List(prop.Of[int]()), "ints")
	intsValues0 := c.Draw(prop.Of[int](), "ints-values")
	return func(_ context.Context, subject Board) error {
		subject.Drain(func() <-chan int {
			ch := make(chan int, len(intsValues))
			for _, v := range intsValues {
				ch <- v
			}
			close(ch)
			return ch
		}(), intsValues0)
		return nil
	}
})

// boardGatherCalls generates the calls of Board.Gather with the inputs
// that each case draws.
var boardGatherCalls = prop.Composite(func(c *prop.Case) func(context.Context, Board) error {
	fnResults := c.Draw(prop.List(prop.Of[int]()), "fn")
	var fnAt atomic.Int64
	fn := func() (r0 int) {
		i := int(fnAt.Add(1) - 1)
		if i >= len(fnResults) {
			return r0
		}
		return fnResults[i]
	}
	fnResults0 := c.Draw(prop.Of[int](), "fn-results")
	fnAt0 := c.Draw(prop.Of[int](), "fn-at")
	return func(_ context.Context, subject Board) error {
		subject.Gather(fn, fnResults0, fnAt0)
		return nil
	}
})

// boardOverEmbedded is Board under the name of the field that boardOver
// embeds.
type boardOverEmbedded = Board

// boardOver is a subject of Board whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type boardOver struct {
	boardOverEmbedded
	spend func()
}

// Place calls spend and returns its zero results.
func (o boardOver) Place(*Point) (r0 error) {
	o.spend()
	return r0
}

// Clash calls spend and returns its zero results.
func (o boardOver) Clash(context.Context, string, string, string, int) (r0 error) {
	o.spend()
	return r0
}

// Blank calls spend and returns its zero results.
func (o boardOver) Blank(string, int) {
	o.spend()
}

// Index calls spend and returns its zero results.
func (o boardOver) Index(map[Shape][]Point) (r0 int) {
	o.spend()
	return r0
}

// Label calls spend and returns its zero results.
func (o boardOver) Label(string) {
	o.spend()
}

// Wait calls spend and returns its zero results.
func (o boardOver) Wait(int, []time.Duration) {
	o.spend()
}

// Hold calls spend and returns its zero results.
func (o boardOver) Hold(int, Pair[time.Month]) {
	o.spend()
}

// Swap calls spend and returns its zero results.
func (o boardOver) Swap(Pair[Tint]) (r0 Pair[Tint]) {
	o.spend()
	return r0
}

// Mark calls spend and returns its zero results.
func (o boardOver) Mark(Cell) (r0 error) {
	o.spend()
	return r0
}

// Paint calls spend and returns its zero results.
func (o boardOver) Paint(Palette) {
	o.spend()
}

// Repeat calls spend and returns its zero results.
func (o boardOver) Repeat(int) {
	o.spend()
}

// Keep calls spend and returns its zero results.
func (o boardOver) Keep(int, store.Item) {
	o.spend()
}

// Stock calls spend and returns its zero results.
func (o boardOver) Stock(int, store.Item) {
	o.spend()
}

// Discard calls spend and returns its zero results.
func (o boardOver) Discard(string) {
	o.spend()
}

// Twice calls spend and returns its zero results.
func (o boardOver) Twice(string, string) {
	o.spend()
}

// Drain calls spend and returns its zero results.
func (o boardOver) Drain(<-chan int, int) {
	o.spend()
}

// Gather calls spend and returns its zero results.
func (o boardOver) Gather(func() int, int, int) {
	o.spend()
}

// newBoardOver returns the planted defect of a ceiling of Board, whose
// callables call spend.
func newBoardOver(_ *BoardHandle, spend func()) Board {
	return boardOver{spend: spend}
}

// boardPlacePanicsDefect is the planted defect of a Board whose Place
// panics.
type boardPlacePanicsDefect struct {
	boardOverEmbedded
}

// Place is the override of the planted defect.
func (*boardPlacePanicsDefect) Place(p0 *Point) (r0 error) {
	panic("the planted defect of a Board whose Place panics")
}

// boardPlacePanics returns the planted defect of a Board whose Place
// panics over a fresh subject.
func boardPlacePanics(h *BoardHandle) Board {
	return &boardPlacePanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardClashPanicsDefect is the planted defect of a Board whose Clash
// panics.
type boardClashPanicsDefect struct {
	boardOverEmbedded
}

// Clash is the override of the planted defect.
func (*boardClashPanicsDefect) Clash(p0 context.Context, p1, p2, p3 string, p4 int) (r0 error) {
	panic("the planted defect of a Board whose Clash panics")
}

// boardClashPanics returns the planted defect of a Board whose Clash
// panics over a fresh subject.
func boardClashPanics(h *BoardHandle) Board {
	return &boardClashPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardClashIgnoresCancelDefect is the planted defect of a Board whose
// Clash ignores the context.
type boardClashIgnoresCancelDefect struct {
	boardOverEmbedded
}

// Clash is the override of the planted defect.
func (d *boardClashIgnoresCancelDefect) Clash(
	p0 context.Context,
	p1, p2, p3 string,
	p4 int,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.boardOverEmbedded.Clash(p0, p1, p2, p3, p4)
}

// boardClashIgnoresCancel returns the planted defect of a Board whose
// Clash ignores the context over a fresh subject.
func boardClashIgnoresCancel(h *BoardHandle) Board {
	return &boardClashIgnoresCancelDefect{boardOverEmbedded: h.Fresh()}
}

// boardBlankPanicsDefect is the planted defect of a Board whose Blank
// panics.
type boardBlankPanicsDefect struct {
	boardOverEmbedded
}

// Blank is the override of the planted defect.
func (*boardBlankPanicsDefect) Blank(p0 string, p1 int) {
	panic("the planted defect of a Board whose Blank panics")
}

// boardBlankPanics returns the planted defect of a Board whose Blank
// panics over a fresh subject.
func boardBlankPanics(h *BoardHandle) Board {
	return &boardBlankPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardIndexPanicsDefect is the planted defect of a Board whose Index
// panics.
type boardIndexPanicsDefect struct {
	boardOverEmbedded
}

// Index is the override of the planted defect.
func (*boardIndexPanicsDefect) Index(p0 map[Shape][]Point) (r0 int) {
	panic("the planted defect of a Board whose Index panics")
}

// boardIndexPanics returns the planted defect of a Board whose Index
// panics over a fresh subject.
func boardIndexPanics(h *BoardHandle) Board {
	return &boardIndexPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardIndexDriftsDefect is the planted defect of a Board whose Index
// returns another value on the repeat.
type boardIndexDriftsDefect struct {
	boardOverEmbedded
	calls int
}

// Index is the override of the planted defect.
func (d *boardIndexDriftsDefect) Index(p0 map[Shape][]Point) (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// boardIndexDrifts returns the planted defect of a Board whose Index
// returns another value on the repeat over a fresh subject.
func boardIndexDrifts(h *BoardHandle) Board {
	return &boardIndexDriftsDefect{boardOverEmbedded: h.Fresh()}
}

// boardLabelPanicsDefect is the planted defect of a Board whose Label
// panics.
type boardLabelPanicsDefect struct {
	boardOverEmbedded
}

// Label is the override of the planted defect.
func (*boardLabelPanicsDefect) Label(p0 string) {
	panic("the planted defect of a Board whose Label panics")
}

// boardLabelPanics returns the planted defect of a Board whose Label
// panics over a fresh subject.
func boardLabelPanics(h *BoardHandle) Board {
	return &boardLabelPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardWaitPanicsDefect is the planted defect of a Board whose Wait
// panics.
type boardWaitPanicsDefect struct {
	boardOverEmbedded
}

// Wait is the override of the planted defect.
func (*boardWaitPanicsDefect) Wait(p0 int, p1 []time.Duration) {
	panic("the planted defect of a Board whose Wait panics")
}

// boardWaitPanics returns the planted defect of a Board whose Wait panics
// over a fresh subject.
func boardWaitPanics(h *BoardHandle) Board {
	return &boardWaitPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardHoldPanicsDefect is the planted defect of a Board whose Hold
// panics.
type boardHoldPanicsDefect struct {
	boardOverEmbedded
}

// Hold is the override of the planted defect.
func (*boardHoldPanicsDefect) Hold(p0 int, p1 Pair[time.Month]) {
	panic("the planted defect of a Board whose Hold panics")
}

// boardHoldPanics returns the planted defect of a Board whose Hold panics
// over a fresh subject.
func boardHoldPanics(h *BoardHandle) Board {
	return &boardHoldPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardSwapPanicsDefect is the planted defect of a Board whose Swap
// panics.
type boardSwapPanicsDefect struct {
	boardOverEmbedded
}

// Swap is the override of the planted defect.
func (*boardSwapPanicsDefect) Swap(p0 Pair[Tint]) (r0 Pair[Tint]) {
	panic("the planted defect of a Board whose Swap panics")
}

// boardSwapPanics returns the planted defect of a Board whose Swap panics
// over a fresh subject.
func boardSwapPanics(h *BoardHandle) Board {
	return &boardSwapPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardSwapDriftsDefect is the planted defect of a Board whose Swap
// returns another value on the repeat.
type boardSwapDriftsDefect struct {
	boardOverEmbedded
	calls int
}

// Swap is the override of the planted defect.
func (d *boardSwapDriftsDefect) Swap(p0 Pair[Tint]) (r0 Pair[Tint]) {
	d.calls++
	if d.calls > 1 {
		r0 = Pair[Tint]{First: Tint(2)}
	}
	return r0
}

// boardSwapDrifts returns the planted defect of a Board whose Swap returns
// another value on the repeat over a fresh subject.
func boardSwapDrifts(h *BoardHandle) Board {
	return &boardSwapDriftsDefect{boardOverEmbedded: h.Fresh()}
}

// boardMarkPanicsDefect is the planted defect of a Board whose Mark
// panics.
type boardMarkPanicsDefect struct {
	boardOverEmbedded
}

// Mark is the override of the planted defect.
func (*boardMarkPanicsDefect) Mark(p0 Cell) (r0 error) {
	panic("the planted defect of a Board whose Mark panics")
}

// boardMarkPanics returns the planted defect of a Board whose Mark panics
// over a fresh subject.
func boardMarkPanics(h *BoardHandle) Board {
	return &boardMarkPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardPaintPanicsDefect is the planted defect of a Board whose Paint
// panics.
type boardPaintPanicsDefect struct {
	boardOverEmbedded
}

// Paint is the override of the planted defect.
func (*boardPaintPanicsDefect) Paint(p0 Palette) {
	panic("the planted defect of a Board whose Paint panics")
}

// boardPaintPanics returns the planted defect of a Board whose Paint
// panics over a fresh subject.
func boardPaintPanics(h *BoardHandle) Board {
	return &boardPaintPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardRepeatPanicsDefect is the planted defect of a Board whose Repeat
// panics.
type boardRepeatPanicsDefect struct {
	boardOverEmbedded
}

// Repeat is the override of the planted defect.
func (*boardRepeatPanicsDefect) Repeat(p0 int) {
	panic("the planted defect of a Board whose Repeat panics")
}

// boardRepeatPanics returns the planted defect of a Board whose Repeat
// panics over a fresh subject.
func boardRepeatPanics(h *BoardHandle) Board {
	return &boardRepeatPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardKeepPanicsDefect is the planted defect of a Board whose Keep
// panics.
type boardKeepPanicsDefect struct {
	boardOverEmbedded
}

// Keep is the override of the planted defect.
func (*boardKeepPanicsDefect) Keep(p0 int, p1 store.Item) {
	panic("the planted defect of a Board whose Keep panics")
}

// boardKeepPanics returns the planted defect of a Board whose Keep panics
// over a fresh subject.
func boardKeepPanics(h *BoardHandle) Board {
	return &boardKeepPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardStockPanicsDefect is the planted defect of a Board whose Stock
// panics.
type boardStockPanicsDefect struct {
	boardOverEmbedded
}

// Stock is the override of the planted defect.
func (*boardStockPanicsDefect) Stock(p0 int, p1 store.Item) {
	panic("the planted defect of a Board whose Stock panics")
}

// boardStockPanics returns the planted defect of a Board whose Stock
// panics over a fresh subject.
func boardStockPanics(h *BoardHandle) Board {
	return &boardStockPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardDiscardPanicsDefect is the planted defect of a Board whose Discard
// panics.
type boardDiscardPanicsDefect struct {
	boardOverEmbedded
}

// Discard is the override of the planted defect.
func (*boardDiscardPanicsDefect) Discard(p0 string) {
	panic("the planted defect of a Board whose Discard panics")
}

// boardDiscardPanics returns the planted defect of a Board whose Discard
// panics over a fresh subject.
func boardDiscardPanics(h *BoardHandle) Board {
	return &boardDiscardPanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardTwicePanicsDefect is the planted defect of a Board whose Twice
// panics.
type boardTwicePanicsDefect struct {
	boardOverEmbedded
}

// Twice is the override of the planted defect.
func (*boardTwicePanicsDefect) Twice(p0, p1 string) {
	panic("the planted defect of a Board whose Twice panics")
}

// boardTwicePanics returns the planted defect of a Board whose Twice
// panics over a fresh subject.
func boardTwicePanics(h *BoardHandle) Board {
	return &boardTwicePanicsDefect{boardOverEmbedded: h.Fresh()}
}

// boardGatherPanicsDefect is the planted defect of a Board whose Gather
// panics.
type boardGatherPanicsDefect struct {
	boardOverEmbedded
}

// Gather is the override of the planted defect.
func (*boardGatherPanicsDefect) Gather(p0 func() int, p1, p2 int) {
	panic("the planted defect of a Board whose Gather panics")
}

// boardGatherPanics returns the planted defect of a Board whose Gather
// panics over a fresh subject.
func boardGatherPanics(h *BoardHandle) Board {
	return &boardGatherPanicsDefect{boardOverEmbedded: h.Fresh()}
}
