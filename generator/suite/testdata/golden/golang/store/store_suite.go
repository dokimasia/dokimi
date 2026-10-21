package store

import (
	"testing"

	"go.dokimi.dev/dokimi-go/conform"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Store writes.
type (
	StoreHandle = conform.Handle[Store]
	StoreCheck  = conform.Check[Store]
	StoreOption = conform.Option[Store]
)

// StoreChecks contains the name of each callable of Store, with the
// identities of its checks.
var StoreChecks = struct {
	Get storeGet
	Put storePut
}{
	Get: storeGet{"Store.Get"},
	Put: storePut{"Store.Put"},
}

// storeGet is the name of Store.Get, with the identities of its checks.
type storeGet struct {
	conform.Callable[Store]
}

// storePut is the name of Store.Put, with the identities of its checks.
type storePut struct {
	conform.Callable[Store]
}

// StoreSuite returns the checks of Store.
func StoreSuite() conform.Suite[Store] {
	return conform.Suite[Store]{Declaration: "Store"}
}

// RunStore applies the checks of Store, and the checks that opts add, to a
// fresh subject of newSubject for each check.
func RunStore[T Store](t *testing.T, newSubject func() T, opts ...StoreOption) {
	t.Helper()
	conform.Verify(t, StoreSuite(), func(*StoreHandle) Store { return newSubject() }, opts...)
}

// RunStoreOn applies the same checks to a subject that newSubject builds
// on the handle of each check.
func RunStoreOn[T Store](t *testing.T, newSubject func(h *StoreHandle) T, opts ...StoreOption) {
	t.Helper()
	conform.Verify(t, StoreSuite(), func(h *StoreHandle) Store { return newSubject(h) }, opts...)
}

// ControlStore applies the same checks to a fresh control of newControl
// for each check. A control is a correct implementation that differs from
// the subject.
func ControlStore[T Store](t *testing.T, newControl func() T, opts ...StoreOption) {
	t.Helper()
	conform.Control(t, StoreSuite(), func(*StoreHandle) Store { return newControl() }, opts...)
}

// ControlStoreOn applies the same checks to a control that newControl
// builds on the handle of each check.
func ControlStoreOn[T Store](t *testing.T, newControl func(h *StoreHandle) T, opts ...StoreOption) {
	t.Helper()
	conform.Control(t, StoreSuite(), func(h *StoreHandle) Store { return newControl(h) }, opts...)
}
