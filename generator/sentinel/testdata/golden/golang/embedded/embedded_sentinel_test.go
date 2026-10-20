package embedded

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
)

// TestEmbeddedSentinels pins the contract every declared error in this
// package carries.
//
// Callers match on these values, read them in logs and branch on them, so
// their messages and identities are as much an API as any exported
// signature — and nothing in the compiler holds them to it.
func TestEmbeddedSentinels(t *testing.T) {
	t.Parallel()

	all := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
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
			if !strings.HasPrefix(e.err.Error(), "embedded: ") {
				t.Errorf("every sentinel must carry the package prefix %q: %s reads %q",
					"embedded: ", e.name, e.err.Error())
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
}

// TestBaseErrorContract pins what BaseError owes a caller beyond taking
// part in the error protocol.
func TestBaseErrorContract(t *testing.T) {
	t.Parallel()

	t.Run("reports a message for its zero value", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, (BaseError{}).Error(), "", "an error must report a message without panicking")
	})

	t.Run("carries the package prefix like a sentinel does", func(t *testing.T) {
		t.Parallel()
		assert.True(t, strings.HasPrefix((BaseError{}).Error(), "embedded: "), "an error type must carry the package prefix")
	})

	t.Run("does not match the package's other error types", func(t *testing.T) {
		t.Parallel()
		subject := BaseError{}
		for _, other := range []error{
			&LostError{},
			NotFoundError{},
		} {
			assert.False(t, errors.Is(subject, other), "an error type must not match another declared here")
		}
	})

	t.Run("exposes the cause it was given", func(t *testing.T) {
		t.Parallel()
		cause := fmt.Errorf("cause")
		var got BaseError
		got.Cause = cause
		assert.Equal(t, errors.Unwrap(got), cause, "an error naming a cause must expose it")
	})

	t.Run("puts every member it carries into its message", func(t *testing.T) {
		t.Parallel()
		got := BaseError{Op: "Op-a"}
		assert.True(t, strings.Contains(got.Error(), string(got.Op)), "Op must reach the message")
	})
}

// TestLostErrorContract pins what LostError owes a caller beyond taking
// part in the error protocol.
func TestLostErrorContract(t *testing.T) {
	t.Parallel()

	t.Run("reports a message for its zero value", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, (&LostError{}).Error(), "", "an error must report a message without panicking")
	})

	t.Run("carries the package prefix like a sentinel does", func(t *testing.T) {
		t.Parallel()
		assert.True(t, strings.HasPrefix((&LostError{}).Error(), "embedded: "), "an error type must carry the package prefix")
	})

	t.Run("does not match the package's other error types", func(t *testing.T) {
		t.Parallel()
		subject := &LostError{}
		for _, other := range []error{
			BaseError{},
			NotFoundError{},
		} {
			assert.False(t, errors.Is(subject, other), "an error type must not match another declared here")
		}
	})
}

// TestNotFoundErrorContract pins what NotFoundError owes a caller beyond
// taking part in the error protocol.
func TestNotFoundErrorContract(t *testing.T) {
	t.Parallel()

	t.Run("reports a message for its zero value", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, (NotFoundError{}).Error(), "", "an error must report a message without panicking")
	})

	t.Run("carries the package prefix like a sentinel does", func(t *testing.T) {
		t.Parallel()
		assert.True(t, strings.HasPrefix((NotFoundError{}).Error(), "embedded: "), "an error type must carry the package prefix")
	})

	t.Run("does not match the package's other error types", func(t *testing.T) {
		t.Parallel()
		subject := NotFoundError{}
		for _, other := range []error{
			BaseError{},
			&LostError{},
		} {
			assert.False(t, errors.Is(subject, other), "an error type must not match another declared here")
		}
	})

	t.Run("exposes the cause it was given", func(t *testing.T) {
		t.Parallel()
		cause := fmt.Errorf("cause")
		var got NotFoundError
		got.Cause = cause
		assert.Equal(t, errors.Unwrap(got), cause, "an error naming a cause must expose it")
	})
}
