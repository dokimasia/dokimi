// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/spec"

	"go.dokimi.dev/dokimi/internal/report"
)

// storeFile is the path of the test file of the Vitest reports of the specs.
const storeFile = "/corpus/iface/store.suite.test.ts"

// vitestTest is a test of the JSON report of Vitest.
type vitestTest struct {
	FullName        string         `json:"fullName"`
	Status          string         `json:"status"`
	FailureMessages []string       `json:"failureMessages"`
	Meta            map[string]any `json:"meta"`
}

// vitestFile is a test file of the JSON report of Vitest.
type vitestFile struct {
	Name             string       `json:"name"`
	Status           string       `json:"status"`
	Message          string       `json:"message"`
	AssertionResults []vitestTest `json:"assertionResults"`
}

// vitest returns the JSON report of Vitest of the files. A count comes before
// the key testResults, as in the report of the JSON reporter of Vitest.
func vitest(files ...vitestFile) string {
	data, _ := json.Marshal(struct {
		NumTotalTestSuites int          `json:"numTotalTestSuites"`
		TestResults        []vitestFile `json:"testResults"`
	}{len(files), files})
	return string(data)
}

// recorded returns the passed test of the block storeSuite whose meta has
// the records.
func recorded(records ...string) vitestTest {
	suite := make([]json.RawMessage, len(records))
	for i, r := range records {
		suite[i] = json.RawMessage(r)
	}
	return vitestTest{
		FullName: "storeSuite suite records", Status: "passed", Meta: map[string]any{"dokimi.suite": suite},
	}
}

func TestVitest(t *testing.T) {
	t.Parallel()

	falsified := record(spec.Falsification, check("Store.Read/zero-on-error", spec.Proven))
	judged := record(spec.OnSubject, check("Store.Read/zero-on-error", spec.Vacuous))

	t.Run("Read", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each suite record of the meta of each test with its file and its test", func(t *testing.T) {
			t.Parallel()

			got, err := report.Read(strings.NewReader(vitest(vitestFile{
				Name:   storeFile,
				Status: "passed",
				AssertionResults: []vitestTest{
					{FullName: "storeSuite store/read/zero-on-error", Status: "passed", Meta: map[string]any{}},
					recorded(falsified, judged),
				},
			})))

			assert.NoError(t, err, "the report reads")
			assert.Length(t, got, 2, "the meta of the last test has two suite records")
			assert.Equal(t, got[0].Package, storeFile, "the record has the path of its test file")
			assert.Equal(t, got[0].Test, "storeSuite suite records", "the record has the full name of its test")
			assert.Equal(t, got[0].Application, spec.Falsification, "the records are in the order of the meta")
			assert.Equal(t, got[1].Checks[0].Result, spec.Vacuous, "the second record follows the first")
		})

		t.Run("reads a report longer than a line of events", func(t *testing.T) {
			t.Parallel()
			checks := make([]spec.CheckRecord, 16384)
			for i := range checks {
				checks[i] = check("Store.Read/keeps-the-item", spec.Pass)
			}

			got, err := report.Read(strings.NewReader(vitest(vitestFile{
				Name:             storeFile,
				Status:           "passed",
				AssertionResults: []vitestTest{recorded(record(spec.OnSubject, checks...))},
			})))

			assert.NoError(t, err, "the report reads")
			assert.Length(t, got, 1, "the report has one record")
			assert.Length(t, got[0].Checks, 16384, "the record has each check")
		})

		t.Run("returns the records that validate beside the error of a record that does not validate",
			func(t *testing.T) {
				t.Parallel()

				got, err := report.Read(strings.NewReader(vitest(vitestFile{
					Name:             storeFile,
					Status:           "passed",
					AssertionResults: []vitestTest{recorded(falsified, `{"spec": "1.0.0"}`)},
				})))

				assert.HasError(t, err, "the second record does not validate")
				assert.Contains(t, err.Error(),
					storeFile+" storeSuite suite records dokimi.suite[1]: spec: the suite record does not validate",
					"the error names the file, the test and the index of the record")
				assert.Length(t, got, 1, "the first record validates")
			})

		t.Run("returns an error for a report whose meta has no array of records", func(t *testing.T) {
			t.Parallel()

			_, err := report.Read(strings.NewReader(vitest(vitestFile{
				Name:   storeFile,
				Status: "passed",
				AssertionResults: []vitestTest{
					{FullName: "storeSuite suite records", Status: "passed", Meta: map[string]any{"dokimi.suite": "x"}},
				},
			})))

			assert.HasError(t, err, "the report does not read")
			assert.Contains(t, err.Error(), "decode the Vitest report: json: cannot unmarshal string",
				"the error names the report")
		})
	})

	t.Run("Failed", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each test that failed with its messages and each file that failed without a failed test",
			func(t *testing.T) {
				t.Parallel()

				got, err := report.Failed(strings.NewReader(vitest(
					vitestFile{Name: storeFile, Status: "failed", AssertionResults: []vitestTest{
						{
							FullName: "storeSuite store/read/cancel",
							Status:   "failed",
							FailureMessages: []string{
								"AssertionFailed: the call ends\n    at store.suite.test.ts:9:3",
								"Error: x",
							},
						},
						{FullName: "storeSuite store/read/smoke", Status: "passed"},
					}},
					vitestFile{
						Name:    "/corpus/iface/log.suite.test.ts",
						Status:  "failed",
						Message: "Cannot find module './log.js'",
					},
					vitestFile{Name: "/corpus/iface/queue.suite.test.ts", Status: "failed"},
					vitestFile{Name: "/corpus/iface/cache.suite.test.ts", Status: "passed"},
				)))

				assert.NoError(t, err, "the report reads")
				assert.Equal(t, got, []report.Failure{
					{
						Package: storeFile, Test: "storeSuite store/read/cancel",
						Output: "AssertionFailed: the call ends\n    at store.suite.test.ts:9:3\nError: x\n",
					},
					{Package: "/corpus/iface/log.suite.test.ts", Output: "Cannot find module './log.js'\n"},
					{Package: "/corpus/iface/queue.suite.test.ts"},
				}, "the failures are the failed test and each file that failed without a failed test")
			})

		t.Run("returns an error for a report whose results do not decode", func(t *testing.T) {
			t.Parallel()

			_, err := report.Failed(strings.NewReader(`{"testResults": 7}`))

			assert.HasError(t, err, "the report does not read")
			assert.Contains(t, err.Error(), "decode the Vitest report: json: cannot unmarshal number",
				"the error names the report")
		})
	})
}
