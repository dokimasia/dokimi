package noprefix

import (
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
)

// TestErrRetry checks the sentinel error ErrRetry. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrRetry(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrGiveUp", ErrGiveUp},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrRetry, "ErrRetry is assigned")
			assert.NotEqual(t, ErrRetry.Error(), "", "ErrRetry has a message")
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrRetry.Error(),
					other.err.Error(),
					"ErrRetry has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrRetry.Error(), other.err.Error()),
					"the message of ErrRetry opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(t, ErrRetry, other.err, "ErrRetry does not match "+other.name)
			}
		})
	})
}

// TestErrGiveUp checks the sentinel error ErrGiveUp. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrGiveUp(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrRetry", ErrRetry},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrGiveUp, "ErrGiveUp is assigned")
			assert.NotEqual(t, ErrGiveUp.Error(), "", "ErrGiveUp has a message")
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrGiveUp.Error(),
					other.err.Error(),
					"ErrGiveUp has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrGiveUp.Error(), other.err.Error()),
					"the message of ErrGiveUp opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(t, ErrGiveUp, other.err, "ErrGiveUp does not match "+other.name)
			}
		})
	})
}
