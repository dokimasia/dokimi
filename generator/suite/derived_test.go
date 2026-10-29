// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"

	"go.dokimi.dev/dokimi/generator/suite"
	"go.dokimi.dev/dokimi/internal/derive"
)

func TestDerived(t *testing.T) {
	t.Parallel()

	got := generate(t, goTree, derivedPackage).suites["Store"]

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("records a derived check for each check of the derivation", func(t *testing.T) {
			t.Parallel()

			assert.Length(t, got.Derived, len(got.Derivation.Checks), "each derived check has its record")
		})

		t.Run("names the constants of the spec of each derived check", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &got, "Get", "read-after-write")
			expect.Equal(t, d.Method, "ReadAfterWrite", "the method of the identity is the segment in Pascal case")
			expect.Equal(t, d.Claim, "ReadAfterWrite", "the claim is the constant of the claim")
			expect.Equal(t, d.Probe, "WriteThenRead", "the probe is the constant of the probe")
			expect.Equal(t, d.Reason, "UnobservedWrite", "the reason is the constant of the reason of the defect")
			expect.Equal(t, d.Law, "ReadReturnsTheWriteBeforeIt", "the law is the constant of the law")
			expect.Equal(t, d.Machineries, []string{"NothingBeyondACaller"}, "the machineries are their constants")
			expect.False(t, d.Clock, "the probe does not advance the clock")
		})

		t.Run("passes the read and the write that returns the written value to a probe of a pair", func(t *testing.T) {
			t.Parallel()

			d := derivedNamed(t, &got, "Get", "read-after-write")
			assert.Length(t, d.Args, 2, "the probe takes the read and the write")
			read, write := d.Args[0].Closure, d.Args[1].Closure
			assert.NotNil(t, read, "the read is a closure")
			assert.NotNil(t, write, "the write is a closure")
			expect.Equal(t, read.Method, "Get", "the read calls Get")
			expect.Equal(t, write.Method, "Put", "the write calls its partner Put")
			assert.NotNil(t, read.Args[1].Of, "the read selects a field of a value")
			expect.Equal(t, *read.Args[1].Of, write.Lets[0].Value, "the read selects the session of the write")
			expect.Equal(t, read.Args[1].Select, "ID", "the read takes the ID of the written session")
			expect.Equal(t, write.Args[1].Name, "w", "the write passes the written session")
			expect.Equal(t, write.Results[0].Name, "w", "the write returns the written session")
		})

		t.Run("shares one planted defect between the checks with one break of one callable", func(t *testing.T) {
			t.Parallel()

			names := make([]string, 0, len(got.Defects))
			for _, d := range got.Defects {
				names = append(names, d.Name)
			}
			assert.NoDuplicates(t, func() ([]string, error) { return names, nil }, "each break has one defect")
			smoke := derivedNamed(t, &got, "Get", "smoke")
			expect.Equal(t, smoke.Defect, "StoreGetPanics", "the defect of a panic of Get is StoreGetPanics")
			raw := derivedNamed(t, &got, "Get", "read-after-write")
			expect.Equal(
				t,
				raw.Defect,
				"StoreGetReadsAnotherValue",
				"the defect of read-after-write breaks the read Get",
			)
		})

		families := generate(t, goTree, familiesPackage).suites
		family := func(name string) *suite.Suite {
			s := families[name]
			return &s
		}
		reach := generate(t, goTree, wrapPackage, toolPackage, deepPackage).suites["Reach"]
		refusals := []struct {
			name     string
			give     *suite.Suite
			giveFrom string
			want     string
		}{
			{
				name: "refuses a check of a cursor whose handle is not an interface", give: family("Sessions"),
				giveFrom: "Open", want: "cursor-open",
			},
			{
				name: "refuses a check whose call passes a sample that the suite cannot import", give: &reach,
				giveFrom: "Store", want: "smoke",
			},
			{
				name: "refuses a check whose write declares a sample that the suite cannot import", give: &reach,
				giveFrom: "Current", want: "read-after-write",
			},
			{
				name: "refuses a check whose planted defect returns a sample that the suite cannot import",
				give: &reach, giveFrom: "Gear", want: "zero-on-error",
			},
		}
		for _, tt := range refusals {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				refused := slices.ContainsFunc(tt.give.Derivation.Refusals, func(r derive.Refusal) bool {
					return r.Callable == tt.giveFrom && string(r.Claim) == tt.want
				})
				expect.True(t, refused, "the suite refuses the claim "+tt.want+" of "+tt.giveFrom)
				derived := slices.ContainsFunc(tt.give.Derived, func(d suite.Derived) bool {
					return d.Field == tt.giveFrom && d.Segment == tt.want
				})
				expect.False(t, derived, "the suite derives no check of the refused claim")
			})
		}
	})
}
