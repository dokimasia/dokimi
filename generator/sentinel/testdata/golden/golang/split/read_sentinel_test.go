package split_test

import (
	"strings"
	"testing"

	"example.com/shop/split"
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
		{"ErrWrite", split.ErrWrite},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, split.ErrRead, "ErrRead is assigned")
			assert.NotEqual(t, split.ErrRead.Error(), "", "ErrRead has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				split.ErrRead.Error(),
				"split: ",
				"the message of ErrRead opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					split.ErrRead.Error(),
					other.err.Error(),
					"ErrRead has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(split.ErrRead.Error(), other.err.Error()),
					"the message of ErrRead opens with the message of "+other.name)
			}
		})
	})

	t.Run("errors.Is", func(t *testing.T) {
		t.Parallel()

		t.Run("matches no other sentinel of the package", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.ErrorIsNot(t, split.ErrRead, other.err, "ErrRead does not match "+other.name)
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
		{"ErrRead", split.ErrRead},
	}

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message", func(t *testing.T) {
			t.Parallel()
			assert.NotNil(t, split.ErrWrite, "ErrWrite is assigned")
			assert.NotEqual(t, split.ErrWrite.Error(), "", "ErrWrite has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				split.ErrWrite.Error(),
				"split: ",
				"the message of ErrWrite opens with the package prefix",
			)
		})

		t.Run("returns a message of its own", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.NotEqual(
					t,
					split.ErrWrite.Error(),
					other.err.Error(),
					"ErrWrite has the message of "+other.name,
				)
			}
		})

		t.Run("opens with the message of no other sentinel", func(t *testing.T) {
			t.Parallel()
			for _, other := range others {
				expect.False(t, strings.HasPrefix(split.ErrWrite.Error(), other.err.Error()),
					"the message of ErrWrite opens with the message of "+other.name)
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
					split.ErrWrite,
					other.err,
					"ErrWrite does not match "+other.name,
				)
			}
		})
	})
}
