package neighbour

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
		{"ErrConflict", ErrConflict},
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
				"neighbour: ",
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
				"neighbour: ",
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
	})
}
