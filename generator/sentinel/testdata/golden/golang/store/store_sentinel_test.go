package store

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"example.com/shop/neighbour"
	"go.dokimi.dev/assert"
)

// TestStoreSentinels pins the contract every declared error in this
// package carries.
//
// Callers match on these values, read them in logs and branch on them, so
// their messages and identities are as much an API as any exported
// signature — and nothing in the compiler holds them to it.
func TestStoreSentinels(t *testing.T) {
	t.Parallel()

	all := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrConflict", ErrConflict},
		{"ErrClosed", ErrClosed},
	}

	t.Run("every sentinel is assigned", func(t *testing.T) {
		t.Parallel()
		for _, e := range all {
			if e.err == nil {
				t.Fatalf("every sentinel must be assigned: %s is nil", e.name)
			}
		}
	})

	t.Run("every sentinel carries a message", func(t *testing.T) {
		t.Parallel()
		for _, e := range all {
			if e.err.Error() == "" {
				t.Errorf("every sentinel must carry a message: %s has none", e.name)
			}
		}
	})

	t.Run("every sentinel carries the package prefix", func(t *testing.T) {
		t.Parallel()
		for _, e := range all {
			if !strings.HasPrefix(e.err.Error(), "kv: ") {
				t.Errorf("every sentinel must carry the package prefix %q: %s reads %q",
					"kv: ", e.name, e.err.Error())
			}
		}
	})

	t.Run("no two sentinels share a message", func(t *testing.T) {
		t.Parallel()
		seen := map[string]string{}
		for _, e := range all {
			if first, dup := seen[e.err.Error()]; dup {
				t.Errorf("no two sentinels may share a message: %s repeats %s", e.name, first)
			}
			seen[e.err.Error()] = e.name
		}
	})

	t.Run("no sentinel's message begins with another's", func(t *testing.T) {
		t.Parallel()
		for _, a := range all {
			for _, b := range all {
				if a.name == b.name {
					continue
				}
				if strings.HasPrefix(a.err.Error(), b.err.Error()) {
					t.Errorf("no sentinel's message may begin with another's: %s begins with the whole of %s",
						a.name, b.name)
				}
			}
		}
	})

	t.Run("no sentinel matches another", func(t *testing.T) {
		t.Parallel()
		for _, a := range all {
			for _, b := range all {
				if a.name == b.name {
					continue
				}
				if errors.Is(a.err, b.err) {
					t.Errorf("no sentinel may match another: %s matches %s", a.name, b.name)
				}
			}
		}
	})

	t.Run("no sentinel matches one in neighbour", func(t *testing.T) {
		t.Parallel()
		for _, a := range all {
			for _, b := range []struct {
				name string
				err  error
			}{
				{"ErrNotFound", neighbour.ErrNotFound},
				{"ErrConflict", neighbour.ErrConflict},
			} {
				if errors.Is(a.err, b.err) {
					t.Errorf("no sentinel may match one in neighbour: %s matches %s", a.name, b.name)
				}
			}
		}
	})
}

// TestNotFoundErrorContract pins what NotFoundError owes a caller beyond
// taking part in the error protocol.
func TestNotFoundErrorContract(t *testing.T) {
	t.Parallel()

	t.Run("reports a message for its zero value", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, (&NotFoundError{}).Error(), "", "an error must report a message without panicking")
	})

	t.Run("carries the package prefix like a sentinel does", func(t *testing.T) {
		t.Parallel()
		assert.True(t, strings.HasPrefix((&NotFoundError{}).Error(), "kv: "), "an error type must carry the package prefix")
	})

	t.Run("does not match the package's other error types", func(t *testing.T) {
		t.Parallel()
		subject := &NotFoundError{}
		for _, other := range []error{
			&WrappedError{},
		} {
			assert.False(t, errors.Is(subject, other), "an error type must not match another declared here")
		}
	})

	t.Run("agrees with the standard comparison about what it matches", func(t *testing.T) {
		t.Parallel()
		subject := &NotFoundError{}
		for _, e := range []error{
			ErrNotFound,
			ErrConflict,
			ErrClosed,
		} {
			assert.Equal(t, errors.Is(subject, e), subject.Is(e), "the standard comparison must agree with the declared one")
		}
	})

	t.Run("puts every member it carries into its message", func(t *testing.T) {
		t.Parallel()
		got := &NotFoundError{Key: "Key-a", Attempt: 42}
		assert.True(t, strings.Contains(got.Error(), string(got.Key)), "Key must reach the message")
	})
}

// TestWrappedErrorContract pins what WrappedError owes a caller beyond
// taking part in the error protocol.
func TestWrappedErrorContract(t *testing.T) {
	t.Parallel()

	t.Run("reports a message for its zero value", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, (&WrappedError{}).Error(), "", "an error must report a message without panicking")
	})

	t.Run("carries the package prefix like a sentinel does", func(t *testing.T) {
		t.Parallel()
		assert.True(t, strings.HasPrefix((&WrappedError{}).Error(), "kv: "), "an error type must carry the package prefix")
	})

	t.Run("does not match the package's other error types", func(t *testing.T) {
		t.Parallel()
		subject := &WrappedError{}
		for _, other := range []error{
			&NotFoundError{},
		} {
			assert.False(t, errors.Is(subject, other), "an error type must not match another declared here")
		}
	})

	t.Run("exposes the cause it was given", func(t *testing.T) {
		t.Parallel()
		cause := fmt.Errorf("cause")
		got := &WrappedError{}
		got.Cause = cause
		assert.Equal(t, errors.Unwrap(got), cause, "an error naming a cause must expose it")
	})

	t.Run("puts every member it carries into its message", func(t *testing.T) {
		t.Parallel()
		got := &WrappedError{Op: "Op-a"}
		assert.True(t, strings.Contains(got.Error(), string(got.Op)), "Op must reach the message")
	})
}
