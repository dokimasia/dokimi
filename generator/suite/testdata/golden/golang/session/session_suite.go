package session

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Sessions:
//
// Sessions.Put/smoke: Put returns without a panic for derived inputs.
//
// Sessions.Put/cancel: Put returns an error for a cancelled context.
//
// Sessions.Get/smoke: Get returns without a panic for derived inputs.
//
// Sessions.Get/cancel: Get returns an error for a cancelled context.
//
// Sessions.Get/zero-on-error: Get returns the zero value beside an error.

// The derived checks of Sessions that witness their laws at the shallow
// end:
//
// Sessions.Put/smoke
//
// Sessions.Get/smoke

// The unclaimed classifications of Sessions:
//
// The shape:writer of Sessions.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

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

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c sessionsPut) Smoke() conform.ID[Sessions] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c sessionsPut) Cancel() conform.ID[Sessions] {
	return c.Check("cancel")
}

// sessionsGet is the name of Sessions.Get, with the identities of its
// checks.
type sessionsGet struct {
	conform.Callable[Sessions]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c sessionsGet) Smoke() conform.ID[Sessions] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c sessionsGet) Cancel() conform.ID[Sessions] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c sessionsGet) ZeroOnError() conform.ID[Sessions] {
	return c.Check("zero-on-error")
}

// SessionsSuite returns the checks of Sessions, and the calls and the
// budgets of its callables.
func SessionsSuite() conform.Suite[Sessions] {
	return conform.Suite[Sessions]{
		Declaration: "Sessions",
		Checks: []conform.Check[Sessions]{
			{
				ID:          SessionsChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SessionsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Session{ID: "ID-a"})
					})
				},
				Defect: sessionsPutPanics,
			},
			{
				ID:          SessionsChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SessionsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Session{ID: "ID-a"})
					})
				},
				Defect: sessionsPutIgnoresCancel,
			},
			{
				ID:          SessionsChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SessionsHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("id-a"))
						return err
					})
				},
				Defect: sessionsGetPanics,
			},
			{
				ID:          SessionsChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SessionsHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("id-a"))
						return err
					})
				},
				Defect: sessionsGetIgnoresCancel,
			},
			{
				ID:          SessionsChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SessionsHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Session, error) {
						return h.Subject().Get(ctx, string("id-a"))
					}, Session{})
				},
				Defect: sessionsGetValueBeside,
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

// RunSessionsBudgets applies the allocation ceilings of Sessions and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
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

// sessionsPutPanicsDefect is the planted defect of a Sessions whose Put
// panics.
type sessionsPutPanicsDefect struct {
	sessionsOverEmbedded
}

// Put is the override of the planted defect.
func (*sessionsPutPanicsDefect) Put(p0 context.Context, p1 Session) (r0 error) {
	panic("the planted defect of a Sessions whose Put panics")
}

// sessionsPutPanics returns the planted defect of a Sessions whose Put
// panics over a fresh subject.
func sessionsPutPanics(h *SessionsHandle) Sessions {
	return &sessionsPutPanicsDefect{sessionsOverEmbedded: h.Fresh()}
}

// sessionsPutIgnoresCancelDefect is the planted defect of a Sessions whose
// Put ignores the context.
type sessionsPutIgnoresCancelDefect struct {
	sessionsOverEmbedded
}

// Put is the override of the planted defect.
func (d *sessionsPutIgnoresCancelDefect) Put(p0 context.Context, p1 Session) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.sessionsOverEmbedded.Put(p0, p1)
}

// sessionsPutIgnoresCancel returns the planted defect of a Sessions whose
// Put ignores the context over a fresh subject.
func sessionsPutIgnoresCancel(h *SessionsHandle) Sessions {
	return &sessionsPutIgnoresCancelDefect{sessionsOverEmbedded: h.Fresh()}
}

// sessionsGetPanicsDefect is the planted defect of a Sessions whose Get
// panics.
type sessionsGetPanicsDefect struct {
	sessionsOverEmbedded
}

// Get is the override of the planted defect.
func (*sessionsGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Session, r1 error) {
	panic("the planted defect of a Sessions whose Get panics")
}

// sessionsGetPanics returns the planted defect of a Sessions whose Get
// panics over a fresh subject.
func sessionsGetPanics(h *SessionsHandle) Sessions {
	return &sessionsGetPanicsDefect{sessionsOverEmbedded: h.Fresh()}
}

// sessionsGetIgnoresCancelDefect is the planted defect of a Sessions whose
// Get ignores the context.
type sessionsGetIgnoresCancelDefect struct {
	sessionsOverEmbedded
}

// Get is the override of the planted defect.
func (d *sessionsGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Session, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.sessionsOverEmbedded.Get(p0, p1)
}

// sessionsGetIgnoresCancel returns the planted defect of a Sessions whose
// Get ignores the context over a fresh subject.
func sessionsGetIgnoresCancel(h *SessionsHandle) Sessions {
	return &sessionsGetIgnoresCancelDefect{sessionsOverEmbedded: h.Fresh()}
}

// sessionsGetValueBesideDefect is the planted defect of a Sessions whose
// Get returns a value beside an error.
type sessionsGetValueBesideDefect struct {
	sessionsOverEmbedded
}

// Get is the override of the planted defect.
func (d *sessionsGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Session, r1 error) {
	r0, r1 = d.sessionsOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Session{ID: "ID-b"}
	}
	return r0, r1
}

// sessionsGetValueBeside returns the planted defect of a Sessions whose
// Get returns a value beside an error over a fresh subject.
func sessionsGetValueBeside(h *SessionsHandle) Sessions {
	return &sessionsGetValueBesideDefect{sessionsOverEmbedded: h.Fresh()}
}
