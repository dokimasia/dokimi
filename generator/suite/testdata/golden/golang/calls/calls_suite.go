package calls

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"example.com/shop/store"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

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

// boardClash is the name of Board.Clash, with the identities of its
// checks.
type boardClash struct {
	conform.Callable[Board]
}

// boardBlank is the name of Board.Blank, with the identities of its
// checks.
type boardBlank struct {
	conform.Callable[Board]
}

// boardIndex is the name of Board.Index, with the identities of its
// checks.
type boardIndex struct {
	conform.Callable[Board]
}

// boardLabel is the name of Board.Label, with the identities of its
// checks.
type boardLabel struct {
	conform.Callable[Board]
}

// boardWait is the name of Board.Wait, with the identities of its checks.
type boardWait struct {
	conform.Callable[Board]
}

// boardHold is the name of Board.Hold, with the identities of its checks.
type boardHold struct {
	conform.Callable[Board]
}

// boardSwap is the name of Board.Swap, with the identities of its checks.
type boardSwap struct {
	conform.Callable[Board]
}

// boardMark is the name of Board.Mark, with the identities of its checks.
type boardMark struct {
	conform.Callable[Board]
}

// boardPaint is the name of Board.Paint, with the identities of its
// checks.
type boardPaint struct {
	conform.Callable[Board]
}

// boardRepeat is the name of Board.Repeat, with the identities of its
// checks.
type boardRepeat struct {
	conform.Callable[Board]
}

// boardKeep is the name of Board.Keep, with the identities of its checks.
type boardKeep struct {
	conform.Callable[Board]
}

// boardStock is the name of Board.Stock, with the identities of its
// checks.
type boardStock struct {
	conform.Callable[Board]
}

// boardDiscard is the name of Board.Discard, with the identities of its
// checks.
type boardDiscard struct {
	conform.Callable[Board]
}

// boardTwice is the name of Board.Twice, with the identities of its
// checks.
type boardTwice struct {
	conform.Callable[Board]
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

// BoardSuite returns the checks of Board, and the calls and the budgets of
// its callables.
func BoardSuite() conform.Suite[Board] {
	return conform.Suite[Board]{
		Declaration: "Board",
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
	conform.Run(t, BoardSuite(), func(h *BoardHandle) Board { return newSubject(h) }, opts...)
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
	conform.Fuzz(
		f,
		BoardSuite(),
		callable,
		func(h *BoardHandle) Board { return newSubject(h) },
		opts...)
}

// RunBoardBudgets applies the allocation ceilings of Board, and the
// ceilings that opts add, to a fresh subject of newSubject for each case.
func RunBoardBudgets[T Board](
	t *testing.T,
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	t.Helper()
	conform.RunBudgets(
		t,
		BoardSuite(),
		func(h *BoardHandle) Board { return newSubject(h) },
		opts...)
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
	conform.Bench(
		b,
		BoardSuite(),
		callable,
		func(h *BoardHandle) Board { return newSubject(h) },
		opts...)
}

// BenchBoardSuiteParallel measures callable on a subject of newSubject on
// GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchBoardSuiteParallel[T Board](
	b *testing.B,
	callable conform.Callable[Board],
	newSubject func(h *BoardHandle) T,
	opts ...BoardOption,
) {
	b.Helper()
	conform.BenchParallel(
		b,
		BoardSuite(),
		callable,
		func(h *BoardHandle) Board { return newSubject(h) },
		opts...)
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

// boardOver is a subject of Board whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type boardOver struct {
	Board
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
