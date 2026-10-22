package session

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// sessionsSubject is a value of Sessions that only TestSessionsSuite
// builds.
type sessionsSubject struct{ Sessions }

// TestSessionsSuite checks the suite of Sessions and the two forms of its
// application.
func TestSessionsSuite(t *testing.T) {
	t.Parallel()

	t.Run("SessionsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, SessionsSuite())
	})

	t.Run("RunSessionsSuite", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		RunSessionsSuite(
			t,
			func() *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *SessionsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Sessions(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*SessionsHandle) Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				SessionsChecks.Get.KeepsTheLastPut(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSessionsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		RunSessionsSuiteFromHandle(
			t,
			func(*SessionsHandle) *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *SessionsHandle) {
					got := h.Subject()
					assert.Equal(h, got, Sessions(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*SessionsHandle) Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				SessionsChecks.Get.KeepsTheLastPut(),
				"the test applies the check of its option alone",
			),
		)
	})
}
