package embedded

import (
	"errors"
	"testing"

	"go.dokimi.dev/assert"
)

// TestErrNotFound checks the sentinel error ErrNotFound. A caller matches
// it with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrNotFound(t *testing.T) {
	t.Parallel()

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
				"embedded: ",
				"the message of ErrNotFound opens with the package prefix",
			)
		})
	})
}

// TestBaseError checks the error type BaseError.
func TestBaseError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing.T) {
			t.Parallel()
			assert.NotEqual(t, (BaseError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(BaseError{}).Error(),
				"embedded: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := BaseError{Op: "Op-a"}
			assert.Contains(t, got.Error(), got.Op, "the message contains Op")
		})
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the cause", func(t *testing.T) {
			t.Parallel()
			cause := errors.New("cause")
			var got BaseError
			got.Cause = cause
			assert.Equal(t, errors.Unwrap(got), cause, "errors.Unwrap returns the cause")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := BaseError{}
			for _, other := range []error{
				&LostError{},
				NotFoundError{},
			} {
				assert.ErrorIsNot(
					t,
					subject,
					other,
					"BaseError does not match another error type of the package",
				)
			}
		})
	})
}

// TestLostError checks the error type LostError.
func TestLostError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing.T) {
			t.Parallel()
			assert.NotEqual(t, (&LostError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(&LostError{}).Error(),
				"embedded: ",
				"the message opens with the package prefix",
			)
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &LostError{}
			for _, other := range []error{
				BaseError{},
				NotFoundError{},
			} {
				assert.ErrorIsNot(
					t,
					subject,
					other,
					"LostError does not match another error type of the package",
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
			assert.NotEqual(t, (NotFoundError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(NotFoundError{}).Error(),
				"embedded: ",
				"the message opens with the package prefix",
			)
		})
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the cause", func(t *testing.T) {
			t.Parallel()
			cause := errors.New("cause")
			var got NotFoundError
			got.Cause = cause
			assert.Equal(t, errors.Unwrap(got), cause, "errors.Unwrap returns the cause")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := NotFoundError{}
			for _, other := range []error{
				BaseError{},
				&LostError{},
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
