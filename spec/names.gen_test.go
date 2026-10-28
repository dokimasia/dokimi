// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec_test

import (
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/spec"
)

func TestNames(t *testing.T) {
	t.Parallel()

	t.Run("constants", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			give string
			want string
		}{
			{
				name: "spells a law in Pascal case",
				give: string(spec.ReadReturnsTheWriteBeforeIt),
				want: "read-returns-the-write-before-it",
			},
			{name: "spells a claim in Pascal case", give: string(spec.ZeroOnError), want: "zero-on-error"},
			{name: "spells a gate in Pascal case", give: string(spec.ReturnsBesideError), want: "returns-beside-error"},
			{name: "spells a probe in Pascal case", give: string(spec.WriteThenRead), want: "write-then-read"},
			{name: "spells a reason code in Pascal case", give: string(spec.UnobservedWrite), want: "unobserved-write"},
			{name: "spells a defect form in Pascal case", give: string(spec.ValueBeside), want: "value-beside"},
			{
				name: "spells a machinery in Pascal case",
				give: string(spec.ClockTheRunAdvances),
				want: "clock-the-run-advances",
			},
			{name: "spells a stance in Pascal case", give: string(spec.Silent), want: "silent"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, tt.give, tt.want, "the constant has the identifier of its entry")
			})
		}
	})
}
