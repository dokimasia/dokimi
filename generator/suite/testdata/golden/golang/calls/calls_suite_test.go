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
	)
}

// FuzzBoardSuitePlace registers the fuzz target of Board.Place on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuitePlace(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Place.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteClash registers the fuzz target of Board.Clash on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteClash(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Clash.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteBlank registers the fuzz target of Board.Blank on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteBlank(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Blank.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteIndex registers the fuzz target of Board.Index on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteIndex(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Index.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteLabel registers the fuzz target of Board.Label on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteLabel(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Label.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteWait registers the fuzz target of Board.Wait on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteWait(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Wait.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteHold registers the fuzz target of Board.Hold on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteHold(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Hold.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteSwap registers the fuzz target of Board.Swap on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteSwap(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Swap.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteMark registers the fuzz target of Board.Mark on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteMark(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Mark.Callable,
		boardIdle,
	)
}

// FuzzBoardSuitePaint registers the fuzz target of Board.Paint on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuitePaint(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Paint.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteRepeat registers the fuzz target of Board.Repeat on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteRepeat(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Repeat.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteKeep registers the fuzz target of Board.Keep on boardIdle,
// with every check of the suite excused.
func FuzzBoardSuiteKeep(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Keep.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteStock registers the fuzz target of Board.Stock on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteStock(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Stock.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteDiscard registers the fuzz target of Board.Discard on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteDiscard(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Discard.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteTwice registers the fuzz target of Board.Twice on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteTwice(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Twice.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteDrain registers the fuzz target of Board.Drain on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteDrain(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Drain.Callable,
		boardIdle,
	)
}

// FuzzBoardSuiteGather registers the fuzz target of Board.Gather on
// boardIdle, with every check of the suite excused.
func FuzzBoardSuiteGather(f *testing.F) {
	FuzzBoardSuite(
		f,
		BoardChecks.Gather.Callable,
		boardIdle,
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
