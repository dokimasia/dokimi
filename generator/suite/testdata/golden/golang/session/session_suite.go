package session

import (
	"testing"

	"go.dokimi.dev/conform"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Sessions writes.
type (
	SessionsHandle = conform.Handle[Sessions]
	SessionsCheck  = conform.Check[Sessions]
	SessionsOption = conform.Option[Sessions]
)

// SessionsChecks contains the name of each callable of Sessions, with the
// identities of its checks.
var SessionsChecks = struct {
	Put sessionsPut
	Get sessionsGet
}{
	Put: sessionsPut{"Sessions.Put"},
	Get: sessionsGet{"Sessions.Get"},
}

// sessionsPut is the name of Sessions.Put, with the identities of its
// checks.
type sessionsPut struct {
	conform.Callable[Sessions]
}

// sessionsGet is the name of Sessions.Get, with the identities of its
// checks.
type sessionsGet struct {
	conform.Callable[Sessions]
}

// KeepsTheLastPut returns the identity of the check that Get returns the
// session of the last Put of its ID.
func (c sessionsGet) KeepsTheLastPut() conform.ID[Sessions] {
	return c.Check("keeps-the-last-put")
}

// SessionsSuite returns the checks of Sessions.
func SessionsSuite() conform.Suite[Sessions] {
	return conform.Suite[Sessions]{Declaration: "Sessions", Checks: []conform.Check[Sessions]{
		{
			ID:        SessionsChecks.Get.KeepsTheLastPut(),
			Statement: "Get returns the session of the last Put of its ID",
			Property:  KeepsTheLastPut,
			Defect:    keepsTheFirstPut,
		},
	}}
}

// RunSessionsSuite applies the checks of Sessions, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunSessionsSuite[T Sessions](t *testing.T, newSubject func() T, opts ...SessionsOption) {
	t.Helper()
	RunSessionsSuiteFromHandle(t, func(*SessionsHandle) T { return newSubject() }, opts...)
}

// RunSessionsSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunSessionsSuiteFromHandle[T Sessions](
	t *testing.T,
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	t.Helper()
	conform.Run(
		t,
		SessionsSuite(),
		func(h *SessionsHandle) Sessions { return newSubject(h) },
		opts...)
}

// FuzzSessionsSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzSessionsSuite[T Sessions](
	f *testing.F,
	callable conform.Callable[Sessions],
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	f.Helper()
	conform.Fuzz(
		f,
		SessionsSuite(),
		callable,
		func(h *SessionsHandle) Sessions { return newSubject(h) },
		opts...)
}
