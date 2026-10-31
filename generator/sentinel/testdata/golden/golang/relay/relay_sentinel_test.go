package relay_test

import (
	"strings"
	"testing"

	"example.com/shop/relay"
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
		{"ErrDropped", relay.ErrDropped},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, relay.ErrNotFound, "ErrNotFound is assigned")
			assert.NotEqual(t, relay.ErrNotFound.Error(), "", "ErrNotFound has a message")
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					relay.ErrNotFound.Error(),
					other.err.Error(),
					"ErrNotFound has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(relay.ErrNotFound.Error(), other.err.Error()),
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
					relay.ErrNotFound,
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
		{"ErrNotFound", relay.ErrNotFound},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, relay.ErrDropped, "ErrDropped is assigned")
			assert.NotEqual(t, relay.ErrDropped.Error(), "", "ErrDropped has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				relay.ErrDropped.Error(),
				"relay: ",
				"the message of ErrDropped opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					relay.ErrDropped.Error(),
					other.err.Error(),
					"ErrDropped has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(relay.ErrDropped.Error(), other.err.Error()),
					"the message of ErrDropped opens with the message of "+other.name)
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
					relay.ErrDropped,
					other.err,
					"ErrDropped does not match "+other.name,
				)
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
			got := &relay.CallError{Message: "Message-a", Code: 42}
			assert.Contains(t, got.Error(), got.Message, "the message contains Message")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &relay.CallError{}
			for _, other := range []error{
				&relay.TimeoutError{},
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
			assert.NotEqual(t, (&relay.TimeoutError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(&relay.TimeoutError{}).Error(),
				"relay: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := &relay.TimeoutError{Op: "Op-a"}
			assert.Contains(t, got.Error(), got.Op, "the message contains Op")
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other error type of the package", func(t *testing.T) {
			t.Parallel()
			subject := &relay.TimeoutError{}
			for _, other := range []error{
				&relay.CallError{},
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
