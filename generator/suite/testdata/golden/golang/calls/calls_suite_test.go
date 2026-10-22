package calls

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// boardSubject is a value of Board that only TestBoardSuite builds.
type boardSubject struct{ Board }

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
