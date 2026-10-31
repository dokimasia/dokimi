package store_test

import (
	"errors"
	"strings"
	"testing"

	"example.com/shop/neighbour"
	"example.com/shop/store"
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
		{"ErrConflict", store.ErrConflict},
		{"ErrClosed", store.ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, store.ErrNotFound, "ErrNotFound is assigned")
			assert.NotEqual(t, store.ErrNotFound.Error(), "", "ErrNotFound has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				store.ErrNotFound.Error(),
				"kv: ",
				"the message of ErrNotFound opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					store.ErrNotFound.Error(),
					other.err.Error(),
					"ErrNotFound has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(store.ErrNotFound.Error(), other.err.Error()),
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
					store.ErrNotFound,
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
					store.ErrNotFound,
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
		{"ErrNotFound", store.ErrNotFound},
		{"ErrClosed", store.ErrClosed},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, store.ErrConflict, "ErrConflict is assigned")
			assert.NotEqual(t, store.ErrConflict.Error(), "", "ErrConflict has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				store.ErrConflict.Error(),
				"kv: ",
				"the message of ErrConflict opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					store.ErrConflict.Error(),
					other.err.Error(),
					"ErrConflict has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(store.ErrConflict.Error(), other.err.Error()),
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
					store.ErrConflict,
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
					store.ErrConflict,
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
		{"ErrNotFound", store.ErrNotFound},
		{"ErrConflict", store.ErrConflict},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, store.ErrClosed, "ErrClosed is assigned")
			assert.NotEqual(t, store.ErrClosed.Error(), "", "ErrClosed has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				store.ErrClosed.Error(),
				"kv: ",
				"the message of ErrClosed opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					store.ErrClosed.Error(),
					other.err.Error(),
					"ErrClosed has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(store.ErrClosed.Error(), other.err.Error()),
					"the message of ErrClosed opens with the message of "+other.name)
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
					store.ErrClosed,
					other.err,
					"ErrClosed does not match "+other.name,
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
					store.ErrClosed,
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
			assert.NotEqual(t, (&store.NotFoundError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(&store.NotFoundError{}).Error(),
				"kv: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := &store.NotFoundError{Key: "Key-a", Attempt: 42}
			assert.Contains(t, got.Error(), got.Key, "the message contains Key")
		})
	})

	t.Run("Is", func(t *testing.T) {
		t.Parallel()

		t.Run("agrees with errors.Is on each sentinel", func(t *testing.T) {
			t.Parallel()
			subject := &store.NotFoundError{}
			for _, e := range []error{
				store.ErrNotFound,
				store.ErrConflict,
				store.ErrClosed,
			} {
				assert.Equal(t, subject.Is(e), errors.Is(subject, e), "Is agrees with errors.Is")
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &store.NotFoundError{}
			for _, other := range []error{
				&store.WrappedError{},
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
			assert.NotEqual(t, (&store.WrappedError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(&store.WrappedError{}).Error(),
				"kv: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := &store.WrappedError{Op: "Op-a"}
			assert.Contains(t, got.Error(), got.Op, "the message contains Op")
		})
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the cause", func(t *testing.T) {
			t.Parallel()
			cause := errors.New("cause")
			got := &store.WrappedError{}
			got.Cause = cause
			assert.Equal(t, errors.Unwrap(got), cause, "errors.Unwrap returns the cause")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &store.WrappedError{}
			for _, other := range []error{
				&store.NotFoundError{},
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
