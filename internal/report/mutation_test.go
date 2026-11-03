// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package report_test

import (
	"slices"
	"strings"
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/internal/report"
	"go.dokimi.dev/dokimi/spec"
)

// target is the package of the tests of the suite records of events.
const target = "example.com/store/storetest"

// mutant returns a mutant of the file store.go on the lines from first to last
// with the verdict, the tests that killed it and the tests that covered it.
func mutant(key string, first, last int, verdict string, tests, coveredBy []string) report.Mutant {
	return report.Mutant{
		Key: key, File: "store.go", Start: report.Position{Line: first, Column: 2},
		End: report.Position{Line: last, Column: 9}, Verdict: verdict, Tests: tests,
		AllKillers: verdict == "killed", CoveredBy: coveredBy,
	}
}

// suites returns the suites of the suite records of events.
func suites(t *testing.T, records ...string) []report.Suite {
	t.Helper()
	read, err := report.Read(strings.NewReader(events(4096, records...)))
	assert.NoError(t, err, "the events read")
	return report.Suites(read)
}

func TestMutation(t *testing.T) {
	t.Parallel()

	judged := record(spec.OnSubject, check("Store.Read/keeps-the-item", spec.Pass),
		check("Store.Read/zero-on-error", spec.Pass), check("Store.Put/smoke", spec.Pass),
		check("Store.Close/smoke", spec.Pass))
	keeps, zero := "TestStore/Read/keeps-the-item", "TestStore/Read/zero-on-error"
	put, closed := "TestStore/Put/smoke", "TestStore/Close/smoke"
	full := report.Mutation{
		Target: report.MutationTarget{Language: "go", Name: target}, Killers: "all", Coverage: "leaves",
		Mutants: []report.Mutant{
			mutant("0000000000000001", 10, 12, "killed", []string{keeps, zero}, []string{keeps, zero, put}),
			mutant("0000000000000002", 14, 14, "killed", []string{keeps}, []string{keeps, zero}),
			mutant("0000000000000003", 20, 21, "survived", nil, []string{closed}),
			mutant("0000000000000004", 88, 90, "survived", nil, nil),
			mutant("0000000000000005", 92, 94, "survived", nil, []string{"TestOther/sorts"}),
			mutant("0000000000000006", 30, 30, "no-coverage", nil, nil),
			mutant("000000000000000c", 40, 41, "killed", []string{"TestOther/sorts"}, []string{"TestOther/sorts"}),
			mutant("000000000000000d", 50, 50, "timed-out", nil, []string{keeps}),
		},
	}

	t.Run("ReadMutations", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each record of a stream of records", func(t *testing.T) {
			t.Parallel()
			one := `{"record": "dokimi-mutate", "version": 2, "target": {"language": "go", "name": "a"},
  "killers": "all", "coverage": "leaves", "mutants": [{"key": "00000000000000aa", "file": "a.go",
  "start": {"line": 3, "column": 1}, "end": {"line": 3, "column": 9}, "verdict": "killed",
  "tests": ["TestA/x"], "allKillers": true, "coveredBy": ["TestA/x"]}]}`
			two := `{"record":"dokimi-mutate","version":2,"target":{"language":"go","name":"b"},"mutants":[]}`

			got, err := report.ReadMutations(strings.NewReader(one + "\n" + two + "\n"))

			assert.NoError(t, err, "the records read")
			assert.Equal(t, got, []report.Mutation{
				{
					Target: report.MutationTarget{Language: "go", Name: "a"}, Killers: "all", Coverage: "leaves",
					Mutants: []report.Mutant{{
						Key: "00000000000000aa", File: "a.go", Start: report.Position{Line: 3, Column: 1},
						End: report.Position{Line: 3, Column: 9}, Verdict: "killed", Tests: []string{"TestA/x"},
						AllKillers: true, CoveredBy: []string{"TestA/x"},
					}},
				},
				{Target: report.MutationTarget{Language: "go", Name: "b"}, Mutants: []report.Mutant{}},
			}, "the stream has two records")
		})

		tests := []struct {
			name       string
			giveRecord string
			want       string
		}{
			{
				name:       "returns an error for a record of another kind",
				giveRecord: `{"record": "stryker", "version": 2, "mutants": []}`,
				want:       `record 1: the record "stryker" of version 2 is no record dokimi-mutate of version 2`,
			},
			{
				name:       "returns an error for a record of another version",
				giveRecord: `{"record": "dokimi-mutate", "version": 1, "mutants": []}`,
				want:       `record 1: the record "dokimi-mutate" of version 1 is no record dokimi-mutate of version 2`,
			},
			{
				name:       "returns an error that names a record that does not decode",
				giveRecord: `{"record": "dokimi-mutate", "version": 2, "mutants": []}` + "\n" + `{"mutants": 7}`,
				want:       "record 2: json: cannot unmarshal number",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				_, err := report.ReadMutations(strings.NewReader(tt.giveRecord))

				assert.HasError(t, err, "the record does not read")
				assert.Contains(t, err.Error(), tt.want, "the error names the record")
			})
		}
	})

	t.Run("Attribute", func(t *testing.T) {
		t.Parallel()

		t.Run("counts the mutants that the checks of the suite covered, killed and let survive", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)

			report.Attribute(got, []report.Mutation{full})

			a := got[0].Mutation
			assert.Equal(t, []int{a.Covered, a.Killed, a.Survived}, []int{3, 2, 1},
				"the checks covered the two kills and the survivor of Close, and a timeout counts as neither")
			assert.Empty(t, a.Lacking, "a record of both modes lacks no mode")
		})

		t.Run("lists the survivors of no check in the suites that covered their target", func(t *testing.T) {
			t.Parallel()
			log := strings.Replace(record(spec.OnSubject, check("Log.Append/smoke", spec.Pass)),
				`"declaration":"Store"`, `"declaration":"Log"`, 1)
			got := suites(t, judged, log)

			report.Attribute(got, []report.Mutation{full})

			assert.Length(t, got[0].Mutation.Unattributed, 1, "the checks of Store covered mutants of the target")
			assert.Empty(t, got[1].Mutation.Unattributed, "no check of Log covered a mutant of the target")
		})

		t.Run("finds each check that covered mutants and killed none", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)

			report.Attribute(got, []report.Mutation{full})

			assert.Equal(t, got[0].Mutation.Dead, []report.DeadCheck{
				{Identity: "Store.Put/smoke", Covered: 1},
				{Identity: "Store.Close/smoke", Covered: 1},
			}, "Put and Close covered a mutant each and killed none")
		})

		t.Run("finds each check whose every kill another check of the suite also made", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)

			report.Attribute(got, []report.Mutation{full})

			assert.Equal(t, got[0].Mutation.Subsumed, []report.SubsumedCheck{
				{Identity: "Store.Read/zero-on-error", By: "Store.Read/keeps-the-item"},
			}, "keeps-the-item kills the mutant of zero-on-error and one more")
		})

		t.Run("leaves out a check whose kills are the kills of another check", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)
			equal := full
			equal.Mutants = []report.Mutant{
				mutant("0000000000000001", 10, 12, "killed", []string{keeps, zero}, []string{keeps, zero}),
			}

			report.Attribute(got, []report.Mutation{equal})

			assert.Empty(t, got[0].Mutation.Subsumed, "no check kills a mutant more than the other")
		})

		t.Run("leaves out a check whose kill a check with more kills did not make", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)
			apart := full
			apart.Mutants = []report.Mutant{
				mutant("0000000000000001", 10, 12, "killed", []string{zero}, []string{keeps, zero}),
				mutant("0000000000000002", 14, 14, "killed", []string{keeps}, []string{keeps}),
				mutant("0000000000000003", 16, 16, "killed", []string{keeps}, []string{keeps}),
			}

			report.Attribute(got, []report.Mutation{apart})

			assert.Empty(t, got[0].Mutation.Subsumed, "keeps-the-item did not kill the mutant of zero-on-error")
		})

		t.Run("joins the check of a function by its identity", func(t *testing.T) {
			t.Parallel()
			slug := strings.Replace(record(spec.OnSubject, check("Slug/keeps-the-case", spec.Pass)),
				`"declaration":"Store"`, `"declaration":"Slug"`, 1)
			got := suites(t, slug)
			function := full
			function.Mutants = []report.Mutant{
				mutant("0000000000000001", 3, 3, "survived", nil, []string{"TestStore/Slug/keeps-the-case"}),
			}

			report.Attribute(got, []report.Mutation{function})

			assert.Equal(t, got[0].Mutation.Dead, []report.DeadCheck{{Identity: "Slug/keeps-the-case", Covered: 1}},
				"the test of the check of a function has the identity of the check as its name")
		})

		t.Run("orders the survivors that no check covered by their file", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)
			files := full
			files.Mutants = append(slices.Clone(full.Mutants), mutant("000000000000000b", 7, 7, "survived", nil, nil))
			files.Mutants[len(files.Mutants)-1].File = "scan.go"

			report.Attribute(got, []report.Mutation{files})

			assert.Equal(t, got[0].Mutation.Unattributed, []report.Survivors{
				{Target: target, File: "scan.go", First: 7, Last: 7, Count: 1},
				{Target: target, File: "store.go", First: 88, Last: 94, Count: 2},
			}, "scan.go comes before store.go")
		})

		t.Run("groups the survivors that no check covered by their file", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)

			report.Attribute(got, []report.Mutation{full})

			assert.Equal(t, got[0].Mutation.Unattributed, []report.Survivors{
				{Target: target, File: "store.go", First: 88, Last: 94, Count: 2},
			}, "the survivor of no test and the survivor of a test that is no check are unattributed")
		})

		t.Run("joins a test of another package by its package and its name", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)
			other := report.Mutation{
				Target: report.MutationTarget{Language: "go", Name: "example.com/store"}, Killers: "all",
				Coverage: "leaves",
				Mutants: []report.Mutant{
					mutant("0000000000000007", 5, 5, "survived", nil, []string{target + ": " + closed}),
					mutant("0000000000000008", 6, 6, "survived", nil, []string{closed}),
				},
			}

			report.Attribute(got, []report.Mutation{other})

			a := got[0].Mutation
			assert.Equal(t, a.Dead, []report.DeadCheck{{Identity: "Store.Close/smoke", Covered: 1}},
				"the name with the package of the test joins the check")
			assert.Equal(t, a.Unattributed, []report.Survivors{
				{Target: "example.com/store", File: "store.go", First: 6, Last: 6, Count: 1},
			}, "the name without the package of the test joins no check of another package")
		})

		t.Run("leaves a kill without every killer out of the dead and subsumed checks", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)
			partial := full
			partial.Mutants = []report.Mutant{
				mutant("0000000000000001", 10, 12, "killed", []string{keeps}, []string{keeps, zero, put}),
				mutant("0000000000000002", 14, 14, "killed", []string{keeps}, []string{keeps, zero}),
			}
			partial.Mutants[0].AllKillers = false

			report.Attribute(got, []report.Mutation{partial})

			a := got[0].Mutation
			assert.Equal(t, a.Killed, 2, "the kill counts as killed")
			assert.Equal(t, a.Dead, []report.DeadCheck{{Identity: "Store.Read/zero-on-error", Covered: 1}},
				"the check that covered the kill of every killer and did not kill it is dead")
			assert.Empty(t, a.Subsumed, "zero-on-error kills no mutant that keeps-the-item could subsume")
		})

		t.Run("names each target whose record lacks a mode of dead and subsumed checks", func(t *testing.T) {
			t.Parallel()
			got := suites(t, judged)
			first := full
			first.Killers = ""
			first.Coverage = ""
			first.Mutants = []report.Mutant{
				mutant("0000000000000009", 3, 3, "survived", nil, []string{"TestStore"}),
				mutant("000000000000000a", 4, 4, "survived", nil, []string{"TestOther"}),
			}

			report.Attribute(got, []report.Mutation{first})

			a := got[0].Mutation
			assert.Equal(t, a.Lacking, []report.Lack{
				{Target: target, Modes: []string{"killers all", "coverage leaves"}},
			}, "the record of the first killer and of the top-level tests lacks both modes")
			assert.Equal(t, []int{a.Covered, a.Survived}, []int{1, 1},
				"the test of the application covers the checks of its record")
			assert.Empty(t, a.Dead, "a record without the modes gives no dead check")
			assert.Equal(t, a.Unattributed, []report.Survivors{
				{Target: target, File: "store.go", First: 4, Last: 4, Count: 1},
			}, "a survivor of a test of no application is unattributed")
		})

		{
			tests := []struct {
				name         string
				giveKillers  string
				giveCoverage string
				giveCovered  string
				want         []string
			}{
				{
					name:         "names the mode of every killer that a record of the leaf coverage lacks",
					giveCoverage: "leaves",
					giveCovered:  closed,
					want:         []string{"killers all"},
				},
				{
					name:        "names the mode of the leaf coverage that a record of every killer lacks",
					giveKillers: "all",
					giveCovered: "TestStore",
					want:        []string{"coverage leaves"},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()
					got := suites(t, judged)
					one := report.Mutation{
						Target: report.MutationTarget{Language: "go", Name: target}, Killers: tt.giveKillers,
						Coverage: tt.giveCoverage,
						Mutants: []report.Mutant{
							mutant("000000000000000e", 3, 3, "survived", nil, []string{tt.giveCovered}),
						},
					}

					report.Attribute(got, []report.Mutation{one})

					a := got[0].Mutation
					assert.Equal(t, a.Lacking, []report.Lack{{Target: target, Modes: tt.want}},
						"the record lacks the other mode")
					assert.Empty(t, a.Dead, "a record of one mode gives no dead check")
				})
			}
		}

		t.Run("leaves out the checks of falsification", func(t *testing.T) {
			t.Parallel()
			got := suites(t, record(spec.Falsification, check("Store.Close/smoke", spec.Proven)))

			report.Attribute(got, []report.Mutation{full})

			a := got[0].Mutation
			assert.Equal(t, []int{a.Covered, len(a.Dead)}, []int{0, 0},
				"a check of falsification runs on a planted defect and covers no code of the subject")
		})
	})
}
