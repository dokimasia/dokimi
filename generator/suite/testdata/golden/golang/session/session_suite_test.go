package session

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// sessionsSubjectEmbedded is Sessions under the name of the field that
// sessionsSubject embeds.
type sessionsSubjectEmbedded = Sessions

// sessionsSubject is a value of Sessions that only TestSessionsSuite
// builds.
type sessionsSubject struct{ sessionsSubjectEmbedded }

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

	t.Run("SessionsChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(SessionsChecks.Get.KeepsTheLastPut()),
			"Sessions.Get/keeps-the-last-put",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// sessionsIdle returns the planted defect of a ceiling of Sessions with a
// spend of nothing, a subject that answers each call at once.
func sessionsIdle(h *SessionsHandle) Sessions {
	return SessionsSuite().Calls[0].Over(h, func() {})
}

// TestSessionsSuiteBudgets applies the allocation ceilings of the suite of
// Sessions to sessionsIdle, and a ceiling of 0 to each callable without a
// budget whose calls build no input of their own.
func TestSessionsSuiteBudgets(t *testing.T) {
	RunSessionsBudgets(
		t,
		sessionsIdle,
		conform.Ceilings(SessionsChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(SessionsChecks.Get.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzSessionsSuitePut registers the fuzz target of Sessions.Put on
// sessionsIdle, with every check of the suite excused.
func FuzzSessionsSuitePut(f *testing.F) {
	FuzzSessionsSuite(
		f,
		SessionsChecks.Put.Callable,
		sessionsIdle,
		conform.Excuse(
			SessionsChecks.Get.KeepsTheLastPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSessionsSuiteGet registers the fuzz target of Sessions.Get on
// sessionsIdle, with every check of the suite excused.
func FuzzSessionsSuiteGet(f *testing.F) {
	FuzzSessionsSuite(
		f,
		SessionsChecks.Get.Callable,
		sessionsIdle,
		conform.Excuse(
			SessionsChecks.Get.KeepsTheLastPut(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSessionsSuite measures each callable of the suite of Sessions
// on sessionsIdle.
func BenchmarkSessionsSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchSessionsSuite(b, SessionsChecks.Put.Callable, sessionsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		BenchSessionsSuite(b, SessionsChecks.Get.Callable, sessionsIdle)
	})
}
