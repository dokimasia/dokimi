package session_test

import (
	"testing"

	"example.com/shop/session"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// sessionsSubjectEmbedded is session.Sessions under the name of the field
// that sessionsSubject embeds.
type sessionsSubjectEmbedded = session.Sessions

// sessionsSubject is a value of Sessions that only TestSessionsSuite
// builds.
type sessionsSubject struct{ sessionsSubjectEmbedded }

// TestSessionsSuite checks the suite of Sessions and the two forms of its
// application.
func TestSessionsSuite(t *testing.T) {
	t.Parallel()

	t.Run("SessionsSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, session.SessionsSuite())
	})

	t.Run("RunSessionsSuite", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		session.RunSessionsSuite(
			t,
			func() *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(session.SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *session.SessionsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						session.Sessions(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*session.SessionsHandle) session.Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				session.SessionsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunSessionsSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *sessionsSubject
		session.RunSessionsSuiteFromHandle(
			t,
			func(*session.SessionsHandle) *sessionsSubject {
				built = &sessionsSubject{}
				return built
			},
			conform.With(session.SessionsCheck{
				ID:        "Sessions/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *session.SessionsHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						session.Sessions(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*session.SessionsHandle) session.Sessions { return &sessionsSubject{} },
			}),
			conform.Excuse(
				session.SessionsChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Put.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Get.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Get.Cancel(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				session.SessionsChecks.Get.ZeroOnError(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// sessionsIdle returns the planted defect of a ceiling of Sessions with a
// spend of nothing, a subject that answers each call at once.
func sessionsIdle(h *session.SessionsHandle) session.Sessions {
	return session.SessionsSuite().Calls[0].Over(h, func() {})
}

// TestSessionsSuiteBudgets applies the allocation ceilings of the suite of
// Sessions to sessionsIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestSessionsSuiteBudgets(t *testing.T) {
	session.RunSessionsBudgets(
		t,
		sessionsIdle,
		conform.Ceilings(
			session.SessionsChecks.Put.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Ceilings(
			session.SessionsChecks.Get.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			session.SessionsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSessionsSuitePut registers the fuzz target of Sessions.Put on
// sessionsIdle, with every check of the suite excused.
func FuzzSessionsSuitePut(f *testing.F) {
	session.FuzzSessionsSuite(
		f,
		session.SessionsChecks.Put.Callable,
		sessionsIdle,
		conform.Excuse(
			session.SessionsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzSessionsSuiteGet registers the fuzz target of Sessions.Get on
// sessionsIdle, with every check of the suite excused.
func FuzzSessionsSuiteGet(f *testing.F) {
	session.FuzzSessionsSuite(
		f,
		session.SessionsChecks.Get.Callable,
		sessionsIdle,
		conform.Excuse(
			session.SessionsChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Put.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.Cancel(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			session.SessionsChecks.Get.ZeroOnError(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkSessionsSuite measures each callable of the suite of Sessions
// on sessionsIdle.
func BenchmarkSessionsSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		session.BenchSessionsSuite(b, session.SessionsChecks.Put.Callable, sessionsIdle)
	})
	b.Run("Get", func(b *testing.B) {
		session.BenchSessionsSuite(b, session.SessionsChecks.Get.Callable, sessionsIdle)
	})
}
