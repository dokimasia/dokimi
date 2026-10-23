package relay

import (
	"strings"
	"testing"

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
		{"ErrDropped", ErrDropped},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrNotFound, "ErrNotFound is assigned")
			assert.NotEqual(t, ErrNotFound.Error(), "", "ErrNotFound has a message")
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
	})
}

// TestErrDropped checks the sentinel error ErrDropped. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrDropped(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrDropped, "ErrDropped is assigned")
			assert.NotEqual(t, ErrDropped.Error(), "", "ErrDropped has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				ErrDropped.Error(),
				"relay: ",
				"the message of ErrDropped opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrDropped.Error(),
					other.err.Error(),
					"ErrDropped has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrDropped.Error(), other.err.Error()),
					"the message of ErrDropped opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(t, ErrDropped, other.err, "ErrDropped does not match "+other.name)
			}
		})
	})
}

// TestCallError checks the error type CallError.
func TestCallError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := &CallError{Message: "Message-a", Code: 42}
			assert.Contains(t, got.Error(), got.Message, "the message contains Message")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &CallError{}
			for _, other := range []error{
				&TimeoutError{},
			} {
				assert.ErrorIsNot(
					t,
					subject,
					other,
					"CallError does not match another error type of the package",
				)
			}
		})
	})
}

// TestTimeoutError checks the error type TimeoutError.
func TestTimeoutError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing.T) {
			t.Parallel()
			assert.NotEqual(t, (&TimeoutError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(&TimeoutError{}).Error(),
				"relay: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := &TimeoutError{Op: "Op-a"}
			assert.Contains(t, got.Error(), got.Op, "the message contains Op")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &TimeoutError{}
			for _, other := range []error{
				&CallError{},
			} {
				assert.ErrorIsNot(
					t,
					subject,
					other,
					"TimeoutError does not match another error type of the package",
				)
			}
		})
	})
}
