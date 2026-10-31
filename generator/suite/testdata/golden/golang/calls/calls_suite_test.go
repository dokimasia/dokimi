package calls_test

import (
	"testing"

	"example.com/shop/calls"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// boardSubjectEmbedded is calls.Board under the name of the field that
// boardSubject embeds.
type boardSubjectEmbedded = calls.Board

// boardSubject is a value of Board that only TestBoardSuite builds.
type boardSubject struct{ boardSubjectEmbedded }

// TestBoardSuite checks the suite of Board and the two forms of its
// application.
func TestBoardSuite(t *testing.T) {
	t.Parallel()

	t.Run("BoardSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, calls.BoardSuite())
	})

	t.Run("RunBoardSuite", func(t *testing.T) {
		t.Parallel()
		var built *boardSubject
		calls.RunBoardSuite(
			t,
			func() *boardSubject {
				built = &boardSubject{}
				return built
			},
			conform.With(calls.BoardCheck{
				ID:        "Board/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *calls.BoardHandle) {
					got := h.Subject()
					assert.Equal(h, got, calls.Board(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*calls.BoardHandle) calls.Board { return &boardSubject{} },
			}),
			conform.Excuse(
				calls.BoardChecks.Place.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Clash.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Clash.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Blank.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Index.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Index.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Label.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Wait.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Hold.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Swap.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Swap.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Mark.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Paint.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Repeat.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Keep.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Stock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Discard.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Twice.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Gather.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBoardSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *boardSubject
		calls.RunBoardSuiteFromHandle(
			t,
			func(*calls.BoardHandle) *boardSubject {
				built = &boardSubject{}
				return built
			},
			conform.With(calls.BoardCheck{
				ID:        "Board/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *calls.BoardHandle) {
					got := h.Subject()
					assert.Equal(h, got, calls.Board(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*calls.BoardHandle) calls.Board { return &boardSubject{} },
			}),
			conform.Excuse(
				calls.BoardChecks.Place.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Clash.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Clash.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Blank.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Index.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Index.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Label.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Wait.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Hold.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Swap.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Swap.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Mark.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Paint.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Repeat.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Keep.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Stock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Discard.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Twice.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				calls.BoardChecks.Gather.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// boardIdle returns the planted defect of a ceiling of Board with a spend
// of nothing, a subject that answers each call at once.
func boardIdle(h *calls.BoardHandle) calls.Board {
	return calls.BoardSuite().Calls[0].Over(h, func() {})
}

// TestBoardSuiteBudgets applies the allocation ceilings of the suite of
// Board to boardIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestBoardSuiteBudgets(t *testing.T) {
	calls.RunBoardBudgets(
		t,
		boardIdle,
		conform.Ceilings(calls.BoardChecks.Place.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Clash.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Blank.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Index.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Label.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Wait.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Hold.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Swap.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Mark.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Paint.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Repeat.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Keep.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Stock.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(
			calls.BoardChecks.Discard.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(calls.BoardChecks.Twice.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(calls.BoardChecks.Gather.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuitePlace registers the fuzz target of Board.Place on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuitePlace(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Place.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteClash registers the fuzz target of Board.Clash on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteClash(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Clash.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteBlank registers the fuzz target of Board.Blank on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteBlank(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Blank.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteIndex registers the fuzz target of Board.Index on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteIndex(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Index.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteLabel registers the fuzz target of Board.Label on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteLabel(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Label.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteWait registers the fuzz target of Board.Wait on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteWait(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Wait.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteHold registers the fuzz target of Board.Hold on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteHold(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Hold.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteSwap registers the fuzz target of Board.Swap on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteSwap(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Swap.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteMark registers the fuzz target of Board.Mark on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteMark(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Mark.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuitePaint registers the fuzz target of Board.Paint on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuitePaint(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Paint.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteRepeat registers the fuzz target of Board.Repeat on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteRepeat(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Repeat.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteKeep registers the fuzz target of Board.Keep on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteKeep(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Keep.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteStock registers the fuzz target of Board.Stock on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteStock(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Stock.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteDiscard registers the fuzz target of Board.Discard on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteDiscard(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Discard.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteTwice registers the fuzz target of Board.Twice on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteTwice(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Twice.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteDrain registers the fuzz target of Board.Drain on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteDrain(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Drain.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteGather registers the fuzz target of Board.Gather on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteGather(f *testing.F) {
	calls.FuzzBoardSuite(
		f,
		calls.BoardChecks.Gather.Callable,
		boardIdle,
		conform.Excuse(
			calls.BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			calls.BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBoardSuite measures each callable of the suite of Board on
// boardIdle.
func BenchmarkBoardSuite(b *testing.B) {
	b.Run("Place", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Place.Callable, boardIdle)
	})
	b.Run("Clash", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Clash.Callable, boardIdle)
	})
	b.Run("Blank", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Blank.Callable, boardIdle)
	})
	b.Run("Index", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Index.Callable, boardIdle)
	})
	b.Run("Label", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Label.Callable, boardIdle)
	})
	b.Run("Wait", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Wait.Callable, boardIdle)
	})
	b.Run("Hold", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Hold.Callable, boardIdle)
	})
	b.Run("Swap", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Swap.Callable, boardIdle)
	})
	b.Run("Mark", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Mark.Callable, boardIdle)
	})
	b.Run("Paint", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Paint.Callable, boardIdle)
	})
	b.Run("Repeat", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Repeat.Callable, boardIdle)
	})
	b.Run("Keep", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Keep.Callable, boardIdle)
	})
	b.Run("Stock", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Stock.Callable, boardIdle)
	})
	b.Run("Discard", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Discard.Callable, boardIdle)
	})
	b.Run("Twice", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Twice.Callable, boardIdle)
	})
	b.Run("Drain", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Drain.Callable, boardIdle)
	})
	b.Run("Gather", func(b *testing.B) {
		calls.BenchBoardSuite(b, calls.BoardChecks.Gather.Callable, boardIdle)
	})
}
