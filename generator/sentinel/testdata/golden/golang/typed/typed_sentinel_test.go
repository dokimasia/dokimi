package typed

import (
	"strings"
	"testing"

	"go.dokimi.dev/assert"
)

// TestLimitErrorContract pins what LimitError owes a caller beyond taking
// part in the error protocol.
func TestLimitErrorContract(t *testing.T) {
	t.Parallel()

	t.Run("reports a message for its zero value", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, (LimitError{}).Error(), "", "an error must report a message without panicking")
	})

	t.Run("carries the package prefix like a sentinel does", func(t *testing.T) {
		t.Parallel()
		assert.True(t, strings.HasPrefix((LimitError{}).Error(), "typed: "), "an error type must carry the package prefix")
	})

	t.Run("puts every member it carries into its message", func(t *testing.T) {
		t.Parallel()
		got := LimitError{Name: "Name-a"}
		assert.True(t, strings.Contains(got.Error(), string(got.Name)), "Name must reach the message")
	})
}
