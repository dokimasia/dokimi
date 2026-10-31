// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/plugin"

	"go.dokimi.dev/dokimi/generator/suite"
)

// lockTemplate is the template of the family of the lock.
const lockTemplate = "lock.tmpl"

// linesOf returns the lines of every lock record of got, in the order of the
// records.
func linesOf(got generated) []string {
	var lines []string
	for _, l := range got.locks {
		lines = append(lines, l.Lines...)
	}
	return lines
}

// line returns the line of the lock of the fields, joined by tabs.
func line(fields ...string) string {
	return strings.Join(fields, "\t")
}

func TestLock(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		t.Run("declares the lock as a line set of one file per plan", func(t *testing.T) {
			t.Parallel()
			outputs, declares := suite.New().(plugin.OutputProvider)
			assert.True(t, declares, "the plugin declares its families")

			assert.Contains(t, outputs.Outputs(), plugin.Output{
				Tag: suite.LockTag, Per: plugin.PerPlan, Word: "lock", Name: suite.LockName, LineSet: true,
				Template: lockTemplate,
			}, "the lock is a data family of line sets")
		})

		t.Run("records a line for each derived check with its partners", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, derivedPackage)

			assert.Contains(t, linesOf(got),
				line(derivedPackage, "store.get/read-after-write", "read-after-write", "check", "write=store.put"),
				"the check has its claim and the key of its partner")
		})

		t.Run("records a line for each refusal", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, derivedPackage)

			derivation := got.suites["Store"].Derivation
			assert.NotEmpty(t, derivation.Refusals, "the fixture refuses claims")
			for _, r := range derivation.Refusals {
				expect.NotEqual(t, r.Key, "", "the refusal of "+string(r.Claim)+" has a key")
				expect.Contains(t, linesOf(got), line(derivedPackage, r.Key, string(r.Claim), "refused", "-"),
					"the refusal has its line")
			}
		})

		t.Run("records a line for each unclaimed classification", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, derivedPackage)

			derivation := got.suites["Store"].Derivation
			assert.NotEmpty(t, derivation.Unclaimed, "the fixture has unclaimed classifications")
			for _, u := range derivation.Unclaimed {
				expect.Contains(t, linesOf(got), line(derivedPackage, u.Key, u.Classification, string(u.Stance), "-"),
					"the unclaimed classification has its line with its stance")
			}
		})

		t.Run("records a line for each check of the consumer", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, storePackage)

			assert.Contains(t, linesOf(got),
				line(storePackage, "store.get/keeps-the-item", "keeps-the-item", "authored", "-"),
				"the check of the directive check has its segment")
		})

		t.Run("records a line for each ceiling of a budget", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, meterPackage)

			assert.Contains(t, linesOf(got),
				line(meterPackage, "meter.count/allocation-budget", "allocation-budget", "ceiling", "-"),
				"the ceiling has its segment")
		})

		t.Run("records the lines of a suite in order, each line once", func(t *testing.T) {
			t.Parallel()

			got := generate(t, goTree, derivedPackage, meterPackage)

			for _, l := range got.locks {
				expect.True(t, slices.IsSorted(l.Lines), "the lines of a suite are sorted")
				expect.Length(t, slices.Compact(slices.Clone(l.Lines)), len(l.Lines), "each line is there once")
			}
		})

		t.Run("records each lock in its family under a declaration without an origin", func(t *testing.T) {
			t.Parallel()
			p := suite.New()

			_, res := phases(t, fixture(t, goTree, p, derivedPackage, meterPackage), p)

			locks := 0
			for u := range res.Emit.Units() {
				if u.Tag != suite.LockTag {
					continue
				}
				locks++
				expect.Equal(t, u.Template, lockTemplate, "the template of the lock writes its file")
				for _, d := range u.Decls {
					origin, _ := emit.OriginOf(d)
					expect.True(t, origin.IsZero(), "a line of the lock stays in the lock of its plan")
				}
			}
			assert.Equal(t, locks, 1, "the plan has one lock")
		})

		tests := []struct {
			name       string
			giveLine   string
			wantReport bool
		}{
			{
				name: "reports a removed line of a derived check",
				giveLine: line(
					derivedPackage, "store.get/read-after-write", "read-after-write", "check", "write=store.put",
				),
				wantReport: true,
			},
			{
				name:       "reports a removed line of a check of the consumer",
				giveLine:   line(storePackage, "store.get/keeps-the-item", "keeps-the-item", "authored", "-"),
				wantReport: true,
			},
			{
				name:       "reports a removed line of a ceiling",
				giveLine:   line(meterPackage, "meter.count/allocation-budget", "allocation-budget", "ceiling", "-"),
				wantReport: true,
			},
			{
				name:     "leaves a removed line of a refusal to the diff",
				giveLine: line(derivedPackage, "store.find/zero-on-error", "zero-on-error", "refused", "-"),
			},
			{
				name:     "leaves a removed line of an unclaimed classification to the diff",
				giveLine: line(derivedPackage, "store.put", "shape:writer", "silent", "-"),
			},
			{
				name:     "leaves a removed line of the header to the diff",
				giveLine: "# dokimi lock 1",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				reader, reads := suite.New().(plugin.LineSetReader)
				assert.True(t, reads, "the plugin reads the removed lines of its lock")

				got, report := reader.Removed(plugin.RemovedLine{Family: suite.LockTag, Line: 3, Text: tt.giveLine})

				assert.Equal(t, report, tt.wantReport, "the plugin reports the removed lines of checks alone")
				if tt.wantReport {
					fields := strings.Split(tt.giveLine, "\t")
					assert.Equal(t, got.Code, suite.LostCheck, "the finding is LostCheck")
					assert.Equal(t, got.Msg, fields[1]+" ("+fields[2]+") of "+fields[0]+
						" is no longer generated as the lock states. Commit the regenerated lock to approve the change",
						"the finding names the key, the claim and the package")
				}
			})
		}
	})
}
