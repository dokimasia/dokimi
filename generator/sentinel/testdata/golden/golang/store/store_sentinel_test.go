package store

import (
	"errors"
	"strings"
	"testing"

	"example.com/shop/neighbour"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
)

// TestErrNotFound checks the sentinel error ErrNotFound. A caller matches
// it with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrNotFound(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrConflict", ErrConflict},
		{"ErrClosed", ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrNotFound, "ErrNotFound is assigned")
			assert.NotEqual(t, ErrNotFound.Error(), "", "ErrNotFound has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				ErrNotFound.Error(),
				"kv: ",
				"the message of ErrNotFound opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrNotFound.Error(),
					other.err.Error(),
					"ErrNotFound has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrNotFound.Error(), other.err.Error()),
					"the message of ErrNotFound opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					ErrNotFound,
					other.err,
					"ErrNotFound does not match "+other.name,
				)
			}
		})

		t.Run("matches no sentinel of neighbour", func(t *testing.T) {
			t.Parallel()
			for _, other := range []struct {
				name string
				err  error
			}{
				{"ErrNotFound", neighbour.ErrNotFound},
				{"ErrConflict", neighbour.ErrConflict},
			} {
				expect.ErrorIsNot(
					t,
					ErrNotFound,
					other.err,
					"ErrNotFound does not match "+other.name+" of neighbour",
				)
			}
		})
	})
}

// TestErrConflict checks the sentinel error ErrConflict. A caller matches
// it with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrConflict(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrClosed", ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrConflict, "ErrConflict is assigned")
			assert.NotEqual(t, ErrConflict.Error(), "", "ErrConflict has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				ErrConflict.Error(),
				"kv: ",
				"the message of ErrConflict opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrConflict.Error(),
					other.err.Error(),
					"ErrConflict has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrConflict.Error(), other.err.Error()),
					"the message of ErrConflict opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(
					t,
					ErrConflict,
					other.err,
					"ErrConflict does not match "+other.name,
				)
			}
		})

		t.Run("matches no sentinel of neighbour", func(t *testing.T) {
			t.Parallel()
			for _, other := range []struct {
				name string
				err  error
			}{
				{"ErrNotFound", neighbour.ErrNotFound},
				{"ErrConflict", neighbour.ErrConflict},
			} {
				expect.ErrorIsNot(
					t,
					ErrConflict,
					other.err,
					"ErrConflict does not match "+other.name+" of neighbour",
				)
			}
		})
	})
}

// TestErrClosed checks the sentinel error ErrClosed. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrClosed(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrConflict", ErrConflict},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrClosed, "ErrClosed is assigned")
			assert.NotEqual(t, ErrClosed.Error(), "", "ErrClosed has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				ErrClosed.Error(),
				"kv: ",
				"the message of ErrClosed opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrClosed.Error(),
					other.err.Error(),
					"ErrClosed has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrClosed.Error(), other.err.Error()),
					"the message of ErrClosed opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(t, ErrClosed, other.err, "ErrClosed does not match "+other.name)
			}
		})

		t.Run("matches no sentinel of neighbour", func(t *testing.T) {
			t.Parallel()
			for _, other := range []struct {
				name string
				err  error
			}{
				{"ErrNotFound", neighbour.ErrNotFound},
				{"ErrConflict", neighbour.ErrConflict},
			} {
				expect.ErrorIsNot(
					t,
					ErrClosed,
					other.err,
					"ErrClosed does not match "+other.name+" of neighbour",
				)
			}
		})
	})
}

// TestNotFoundError checks the error type NotFoundError.
func TestNotFoundError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing.T) {
			t.Parallel()
			assert.NotEqual(t, (&NotFoundError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(&NotFoundError{}).Error(),
				"kv: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := &NotFoundError{Key: "Key-a", Attempt: 42}
			assert.Contains(t, got.Error(), got.Key, "the message contains Key")
		})
	})

	t.Run("Is", func(t *testing.T) {
		t.Parallel()

		t.Run("agrees with errors.Is on each sentinel", func(t *testing.T) {
			t.Parallel()
			subject := &NotFoundError{}
			for _, e := range []error{
				ErrNotFound,
				ErrConflict,
				ErrClosed,
			} {
				assert.Equal(t, subject.Is(e), errors.Is(subject, e), "Is agrees with errors.Is")
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &NotFoundError{}
			for _, other := range []error{
				&WrappedError{},
			} {
				assert.ErrorIsNot(
					t,
					subject,
					other,
					"NotFoundError does not match another error type of the package",
				)
			}
		})
	})
}

// TestWrappedError checks the error type WrappedError.
func TestWrappedError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing.T) {
			t.Parallel()
			assert.NotEqual(t, (&WrappedError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(&WrappedError{}).Error(),
				"kv: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := &WrappedError{Op: "Op-a"}
			assert.Contains(t, got.Error(), got.Op, "the message contains Op")
		})
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the cause", func(t *testing.T) {
			t.Parallel()
			cause := errors.New("cause")
			got := &WrappedError{}
			got.Cause = cause
			assert.Equal(t, errors.Unwrap(got), cause, "errors.Unwrap returns the cause")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &WrappedError{}
			for _, other := range []error{
				&NotFoundError{},
			} {
				assert.ErrorIsNot(
					t,
					subject,
					other,
					"WrappedError does not match another error type of the package",
				)
			}
		})
	})
}
