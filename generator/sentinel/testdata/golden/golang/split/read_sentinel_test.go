package split

import (
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
)

// TestErrRead checks the sentinel error ErrRead. A caller matches it with
// errors.Is and reads its message, so it has a message and an identity of
// its own.
func TestErrRead(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrWrite", ErrWrite},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrRead, "ErrRead is assigned")
			assert.NotEqual(t, ErrRead.Error(), "", "ErrRead has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				ErrRead.Error(),
				"split: ",
				"the message of ErrRead opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrRead.Error(),
					other.err.Error(),
					"ErrRead has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrRead.Error(), other.err.Error()),
					"the message of ErrRead opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(t, ErrRead, other.err, "ErrRead does not match "+other.name)
			}
		})
	})
}

// TestErrWrite checks the sentinel error ErrWrite. A caller matches it
// with errors.Is and reads its message, so it has a message and an
// identity of its own.
func TestErrWrite(t *testing.T) {
	t.Parallel()

	others := []struct {
		name string
		err  error
	}{
		{"ErrRead", ErrRead},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, ErrWrite, "ErrWrite is assigned")
			assert.NotEqual(t, ErrWrite.Error(), "", "ErrWrite has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				ErrWrite.Error(),
				"split: ",
				"the message of ErrWrite opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					ErrWrite.Error(),
					other.err.Error(),
					"ErrWrite has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(ErrWrite.Error(), other.err.Error()),
					"the message of ErrWrite opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(t, ErrWrite, other.err, "ErrWrite does not match "+other.name)
			}
		})
	})
}
