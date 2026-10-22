package typed

import (
	"testing"

	"go.dokimi.dev/assert"
)

// TestLimitError checks the error type LimitError.
func TestLimitError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		t.Run("returns a message for the zero value", func(t *testing.T) {
			t.Parallel()
			assert.NotEqual(t, (LimitError{}).Error(), "", "the zero value has a message")
		})

		t.Run("opens the message with the package prefix", func(t *testing.T) {
			t.Parallel()
			assert.HasPrefix(
				t,
				(LimitError{}).Error(),
				"typed: ",
				"the message opens with the package prefix",
			)
		})

		t.Run("contains each text member", func(t *testing.T) {
			t.Parallel()
			got := LimitError{Name: "Name-a"}
			assert.Contains(t, got.Error(), got.Name, "the message contains Name")
		})
	})
}
