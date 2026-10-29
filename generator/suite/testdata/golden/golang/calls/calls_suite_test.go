package calls

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// boardSubjectEmbedded is Board under the name of the field that
// boardSubject embeds.
type boardSubjectEmbedded = Board

// boardSubject is a value of Board that only TestBoardSuite builds.
type boardSubject struct{ boardSubjectEmbedded }

// TestBoardSuite checks the suite of Board and the two forms of its
// application.
func TestBoardSuite(t *testing.T) {
	t.Parallel()

	t.Run("BoardSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, BoardSuite())
	})

	t.Run("RunBoardSuite", func(t *testing.T) {
		t.Parallel()
		var built *boardSubject
		RunBoardSuite(
			t,
			func() *boardSubject {
				built = &boardSubject{}
				return built
			},
			conform.With(BoardCheck{
				ID:        "Board/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *BoardHandle) {
					got := h.Subject()
					assert.Equal(h, got, Board(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*BoardHandle) Board { return &boardSubject{} },
			}),
			conform.Excuse(
				BoardChecks.Place.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Clash.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Clash.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Blank.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Index.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Index.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Label.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Wait.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Hold.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Swap.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Swap.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Mark.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Paint.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Repeat.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Keep.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Stock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Discard.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Twice.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Gather.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunBoardSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *boardSubject
		RunBoardSuiteFromHandle(
			t,
			func(*BoardHandle) *boardSubject {
				built = &boardSubject{}
				return built
			},
			conform.With(BoardCheck{
				ID:        "Board/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *BoardHandle) {
					got := h.Subject()
					assert.Equal(h, got, Board(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*BoardHandle) Board { return &boardSubject{} },
			}),
			conform.Excuse(
				BoardChecks.Place.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Clash.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Clash.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Blank.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Index.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Index.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Label.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Wait.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Hold.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Swap.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Swap.Deterministic(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Mark.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Paint.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Repeat.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Keep.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Stock.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Discard.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Twice.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				BoardChecks.Gather.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// boardIdle returns the planted defect of a ceiling of Board with a spend
// of nothing, a subject that answers each call at once.
func boardIdle(h *BoardHandle) Board {
	return BoardSuite().Calls[0].Over(h, func() {})
}

// TestBoardSuiteBudgets applies the allocation ceilings of the suite of
// Board to boardIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestBoardSuiteBudgets(t *testing.T) {
	RunBoardBudgets(
		t,
		boardIdle,
		conform.Ceilings(BoardChecks.Place.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Clash.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Blank.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Index.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Label.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Wait.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Hold.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Swap.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Mark.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Paint.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Repeat.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Keep.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Stock.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Discard.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Twice.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(BoardChecks.Gather.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuitePlace registers the fuzz target of Board.Place on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuitePlace(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Place.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteClash registers the fuzz target of Board.Clash on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteClash(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Clash.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteBlank registers the fuzz target of Board.Blank on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteBlank(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Blank.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteIndex registers the fuzz target of Board.Index on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteIndex(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Index.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteLabel registers the fuzz target of Board.Label on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteLabel(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Label.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteWait registers the fuzz target of Board.Wait on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteWait(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Wait.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteHold registers the fuzz target of Board.Hold on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteHold(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Hold.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteSwap registers the fuzz target of Board.Swap on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteSwap(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Swap.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteMark registers the fuzz target of Board.Mark on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteMark(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Mark.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuitePaint registers the fuzz target of Board.Paint on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuitePaint(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Paint.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteRepeat registers the fuzz target of Board.Repeat on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteRepeat(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Repeat.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteKeep registers the fuzz target of Board.Keep on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteKeep(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Keep.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteStock registers the fuzz target of Board.Stock on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteStock(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Stock.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteDiscard registers the fuzz target of Board.Discard on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteDiscard(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Discard.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteTwice registers the fuzz target of Board.Twice on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteTwice(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Twice.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteDrain registers the fuzz target of Board.Drain on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteDrain(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Drain.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzBoardSuiteGather registers the fuzz target of Board.Gather on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteGather(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Gather.Callable,
		boardIdle,
		conform.Excuse(
			BoardChecks.Place.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Clash.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Blank.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Index.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Label.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Wait.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Hold.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Swap.Deterministic(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Mark.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Paint.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Repeat.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Keep.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Stock.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Discard.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Twice.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			BoardChecks.Gather.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkBoardSuite measures each callable of the suite of Board on
// boardIdle.
func BenchmarkBoardSuite(b *testing.B) {
	b.Run("Place", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Place.Callable, boardIdle)
	})
	b.Run("Clash", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Clash.Callable, boardIdle)
	})
	b.Run("Blank", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Blank.Callable, boardIdle)
	})
	b.Run("Index", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Index.Callable, boardIdle)
	})
	b.Run("Label", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Label.Callable, boardIdle)
	})
	b.Run("Wait", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Wait.Callable, boardIdle)
	})
	b.Run("Hold", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Hold.Callable, boardIdle)
	})
	b.Run("Swap", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Swap.Callable, boardIdle)
	})
	b.Run("Mark", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Mark.Callable, boardIdle)
	})
	b.Run("Paint", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Paint.Callable, boardIdle)
	})
	b.Run("Repeat", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Repeat.Callable, boardIdle)
	})
	b.Run("Keep", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Keep.Callable, boardIdle)
	})
	b.Run("Stock", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Stock.Callable, boardIdle)
	})
	b.Run("Discard", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Discard.Callable, boardIdle)
	})
	b.Run("Twice", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Twice.Callable, boardIdle)
	})
	b.Run("Drain", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Drain.Callable, boardIdle)
	})
	b.Run("Gather", func(b *testing.B) {
		BenchBoardSuite(b, BoardChecks.Gather.Callable, boardIdle)
	})
}
