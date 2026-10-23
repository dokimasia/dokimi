package session

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
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

// SessionsSuite returns the checks of Sessions, and the calls and the
// budgets of its callables.
func SessionsSuite() conform.Suite[Sessions] {
	return conform.Suite[Sessions]{
		Declaration: "Sessions",
		Checks: []conform.Check[Sessions]{
			{
				ID:        SessionsChecks.Get.KeepsTheLastPut(),
				Statement: "Get returns the session of the last Put of its ID",
				Property:  KeepsTheLastPut,
				Defect:    keepsTheFirstPut,
			},
		},
		Calls: []conform.Call[Sessions]{
			{
				Calls:    sessionsPutCalls,
				Over:     newSessionsOver,
				Callable: SessionsChecks.Put.Callable,
			},
			{
				Calls:    sessionsGetCalls,
				Over:     newSessionsOver,
				Callable: SessionsChecks.Get.Callable,
			},
		},
	}
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
	conform.Run(t, SessionsSuite(), sessionsSubjects(newSubject), opts...)
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
	conform.Fuzz(f, SessionsSuite(), callable, sessionsSubjects(newSubject), opts...)
}

// RunSessionsBudgets applies the allocation ceilings of Sessions, and the
// ceilings that opts add, to a fresh subject of newSubject for each case.
func RunSessionsBudgets[T Sessions](
	t *testing.T,
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	t.Helper()
	conform.RunBudgets(t, SessionsSuite(), sessionsSubjects(newSubject), opts...)
}

// BenchSessionsSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchSessionsSuite[T Sessions](
	b conform.B,
	callable conform.Callable[Sessions],
	newSubject func(h *SessionsHandle) T,
	opts ...SessionsOption,
) {
	b.Helper()
	conform.Bench(b, SessionsSuite(), callable, sessionsSubjects(newSubject), opts...)
}

// sessionsSubjects returns the factory of the subjects of Sessions that
// calls newSubject with the handle of each check.
func sessionsSubjects[T Sessions](
	newSubject func(h *SessionsHandle) T,
) func(h *SessionsHandle) Sessions {
	return func(h *SessionsHandle) Sessions { return newSubject(h) }
}

// sessionsPutCalls generates the calls of Sessions.Put with the inputs
// that each case draws.
var sessionsPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Sessions) error {
	s := c.Draw(prop.Of[Session](), "s")
	return func(ctx context.Context, subject Sessions) error {
		return subject.Put(ctx, s)
	}
})

// sessionsGetCalls generates the calls of Sessions.Get with the inputs
// that each case draws.
var sessionsGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Sessions) error {
	id := c.Draw(prop.Of[string](), "id")
	return func(ctx context.Context, subject Sessions) error {
		_, err := subject.Get(ctx, id)
		return err
	}
})

// sessionsOverEmbedded is Sessions under the name of the field that
// sessionsOver embeds.
type sessionsOverEmbedded = Sessions

// sessionsOver is a subject of Sessions whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type sessionsOver struct {
	sessionsOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o sessionsOver) Put(context.Context, Session) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o sessionsOver) Get(context.Context, string) (r0 Session, r1 error) {
	o.spend()
	return r0, r1
}

// newSessionsOver returns the planted defect of a ceiling of Sessions,
// whose callables call spend.
func newSessionsOver(_ *SessionsHandle, spend func()) Sessions {
	return sessionsOver{spend: spend}
}
