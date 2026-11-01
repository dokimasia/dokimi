package scoped_test

import (
	errors2 "errors"
	strings2 "strings"
	testing2 "testing"

	"example.com/shop/scoped"
	assert2 "go.dokimi.dev/assert"
	expect2 "go.dokimi.dev/assert/expect"
)

// TestErrMissing checks the sentinel error ErrMissing. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrMissing(t *testing2.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrStale", scoped.ErrStale},
	}

	t.Run("Error", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing2.T) {
			t.Parallel()
			assert2.NotNil(t, scoped.ErrMissing, "ErrMissing is assigned")
			assert2.NotEqual(t, scoped.ErrMissing.Error(), "", "ErrMissing has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing2.T) {
			t.Parallel()
			assert2.HasPrefix(
				t,
				scoped.ErrMissing.Error(),
				"scoped: ",
				"the message of ErrMissing opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing2.T) {
			t.Parallel()
			for _, other := range others {
				expect2.NotEqual(
					t,
					scoped.ErrMissing.Error(),
					other.err.Error(),
					"ErrMissing has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing2.T) {
			t.Parallel()
			for _, other := range others {
				expect2.False(t, strings2.HasPrefix(scoped.ErrMissing.Error(), other.err.Error()),
					"the message of ErrMissing opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing2.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing2.T) {
			t.Parallel()
			for _, other := range others {
				expect2.ErrorIsNot(
					t,
					scoped.ErrMissing,
					other.err,
					"ErrMissing does not match "+other.name,
				)
			}
		})
	})
}

// TestErrStale checks the sentinel error ErrStale. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrStale(t *testing2.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrMissing", scoped.ErrMissing},
	}

	t.Run("Error", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing2.T) {
			t.Parallel()
			assert2.NotNil(t, scoped.ErrStale, "ErrStale is assigned")
			assert2.NotEqual(t, scoped.ErrStale.Error(), "", "ErrStale has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing2.T) {
			t.Parallel()
			assert2.HasPrefix(
				t,
				scoped.ErrStale.Error(),
				"scoped: ",
				"the message of ErrStale opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing2.T) {
			t.Parallel()
			for _, other := range others {
				expect2.NotEqual(
					t,
					scoped.ErrStale.Error(),
					other.err.Error(),
					"ErrStale has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing2.T) {
			t.Parallel()
			for _, other := range others {
				expect2.False(t, strings2.HasPrefix(scoped.ErrStale.Error(), other.err.Error()),
					"the message of ErrStale opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing2.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing2.T) {
			t.Parallel()
			for _, other := range others {
				expect2.ErrorIsNot(
					t,
					scoped.ErrStale,
					other.err,
					"ErrStale does not match "+other.name,
				)
			}
		})
	})
}

// TestMissingError checks the error type MissingError.
func TestMissingError(t *testing2.T) {
	t.Parallel()

	t.Run("Error", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing2.T) {
			t.Parallel()
			assert2.NotEqual(
				t,
				(&scoped.MissingError{}).Error(),
				"",
				"the zero value has a message",
			)
		})

		t.Run("opens the message with the package prefix", func(t *testing2.T) {
			t.Parallel()
			assert2.HasPrefix(
				t,
				(&scoped.MissingError{}).Error(),
				"scoped: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing2.T) {
			t.Parallel()
			got := &scoped.MissingError{Key: "Key-a"}
			assert2.Contains(t, got.Error(), got.Key, "the message contains Key")
		})
	})

	t.Run("Is", func(t *testing2.T) {
		t.Parallel()

		t.Run("agrees with errors.Is on each sentinel", func(t *testing2.T) {
			t.Parallel()
			subject := &scoped.MissingError{}
			for _, e := range []error{
				scoped.ErrMissing,
				scoped.ErrStale,
			} {
				assert2.Equal(t, subject.Is(e), errors2.Is(subject, e), "Is agrees with errors.Is")
			}
		})
	})

	t.Run("errors.Is", func(t *testing2.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing2.T) {
			t.Parallel()
			subject := &scoped.MissingError{}
			for _, other := range []error{
				&scoped.WrappedError{},
			} {
				assert2.ErrorIsNot(
					t,
					subject,
					other,
					"MissingError does not match another error type of the package",
				)
			}
		})
	})
}

// TestWrappedError checks the error type WrappedError.
func TestWrappedError(t *testing2.T) {
	t.Parallel()

	t.Run("Error", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing2.T) {
			t.Parallel()
			assert2.NotEqual(
				t,
				(&scoped.WrappedError{}).Error(),
				"",
				"the zero value has a message",
			)
		})

		t.Run("opens the message with the package prefix", func(t *testing2.T) {
			t.Parallel()
			assert2.HasPrefix(
				t,
				(&scoped.WrappedError{}).Error(),
				"scoped: ",
				"the message opens with the package prefix",
			)
		})
	})

	t.Run("Unwrap", func(t *testing2.T) {
		t.Parallel()

		t.Run("returns the cause", func(t *testing2.T) {
			t.Parallel()
			cause := errors2.New("cause")
			got := &scoped.WrappedError{}
			got.Cause = cause
			assert2.Equal(t, errors2.Unwrap(got), cause, "errors.Unwrap returns the cause")
		})
	})

	t.Run("errors.Is", func(t *testing2.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing2.T) {
			t.Parallel()
			subject := &scoped.WrappedError{}
			for _, other := range []error{
				&scoped.MissingError{},
			} {
				assert2.ErrorIsNot(
					t,
					subject,
					other,
					"WrappedError does not match another error type of the package",
				)
			}
		})
	})
}
