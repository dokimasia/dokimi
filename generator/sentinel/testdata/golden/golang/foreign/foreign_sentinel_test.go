package foreign_test

import (
	"strings"
	"testing"

	"example.com/shop/foreign"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
)

// TestErrMissing checks the sentinel error ErrMissing. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrMissing(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrStale", foreign.ErrStale},
		{"ErrNotFound", foreign.ErrNotFound},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, foreign.ErrMissing, "ErrMissing is assigned")
			assert.NotEqual(t, foreign.ErrMissing.Error(), "", "ErrMissing has a message")
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					foreign.ErrMissing.Error(),
					other.err.Error(),
					"ErrMissing has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(foreign.ErrMissing.Error(), other.err.Error()),
					"the message of ErrMissing opens with the message of "+other.name)
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
					foreign.ErrMissing,
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
func TestErrStale(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrMissing", foreign.ErrMissing},
		{"ErrNotFound", foreign.ErrNotFound},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, foreign.ErrStale, "ErrStale is assigned")
			assert.NotEqual(t, foreign.ErrStale.Error(), "", "ErrStale has a message")
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					foreign.ErrStale.Error(),
					other.err.Error(),
					"ErrStale has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(foreign.ErrStale.Error(), other.err.Error()),
					"the message of ErrStale opens with the message of "+other.name)
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
					foreign.ErrStale,
					other.err,
					"ErrStale does not match "+other.name,
				)
			}
		})
	})
}

// TestErrNotFound checks the sentinel error ErrNotFound. A caller matches
// it with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrNotFound(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrMissing", foreign.ErrMissing},
		{"ErrStale", foreign.ErrStale},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, foreign.ErrNotFound, "ErrNotFound is assigned")
			assert.NotEqual(t, foreign.ErrNotFound.Error(), "", "ErrNotFound has a message")
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					foreign.ErrNotFound.Error(),
					other.err.Error(),
					"ErrNotFound has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(foreign.ErrNotFound.Error(), other.err.Error()),
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
					foreign.ErrNotFound,
					other.err,
					"ErrNotFound does not match "+other.name,
				)
			}
		})
	})
}
