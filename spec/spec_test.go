// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"os"
	"strings"
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/spec"
)

func TestSpec(t *testing.T) {
	t.Parallel()

	t.Run("Version", func(t *testing.T) {
		t.Parallel()

		t.Run("returns the content of VERSION without its newline", func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile("VERSION")
			assert.NoError(t, err, "VERSION is readable")

			assert.Equal(t, spec.Version(), strings.TrimSpace(string(data)), "the version is the content of VERSION")
		})
	})
}
