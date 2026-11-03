package logtest_test

import (
	"testing"

	"example.com/shop/log/logtest"
)

// TestLogFake applies the checks of the suite of Log to LogFake.
func TestLogFake(t *testing.T) {
	t.Parallel()

	logtest.RunLogSuite(t, logtest.NewLogFake)
}
