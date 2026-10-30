// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/internal/report"
	"go.dokimi.dev/dokimi/spec"
)

// record returns the JSON of a suite record of the declaration Store in the
// application, with the checks.
func record(application spec.Application, checks ...spec.CheckRecord) string {
	r := spec.SuiteRecord{
		Spec:        "1.0.0",
		Suite:       spec.RecordSuite{Language: "go", Package: "example.com/store", Declaration: "Store"},
		Generator:   spec.RecordTool{Name: "dokimi", Version: "0.1.0"},
		Runtime:     spec.RecordTool{Name: "go.dokimi.dev/conform", Version: "v0.9.0"},
		Seed:        "0x5f3a1c9e27d4b860",
		Application: application,
		Checks:      checks,
		Refusals: []spec.RefusalRecord{{
			Claim: spec.ZeroOnError, Callable: "Store.Scan", Why: "Scan cannot fail", Remedy: "return an error",
		}},
		Unclaimed: []spec.UnclaimedRecord{{
			Classification: "mixin:concurrent", Callable: "Store.Put", Stance: spec.Needs,
			Note: "needs concurrent callers",
		}},
	}
	data, _ := json.Marshal(r)
	return string(data)
}

// check returns the record of the check of the identity with the result.
func check(identity string, result spec.Result) spec.CheckRecord {
	return spec.CheckRecord{Identity: identity, Key: strings.ToLower(identity), Result: result}
}

// events returns the events of go test -json that write each record under the
// test TestStore, one attribute event for each record, split into parts of at
// most size bytes.
func events(size int, records ...string) string {
	var b strings.Builder
	b.WriteString(`{"Action":"start","Package":"example.com/store/storetest"}` + "\n")
	b.WriteString("a line that is no event\n")
	for i, r := range records {
		key := "dokimi.suite." + string(rune('1'+i))
		for len(r) > 0 {
			part := r[:min(size, len(r))]
			r = r[len(part):]
			line, _ := json.Marshal(map[string]string{
				"Action": "attr", "Package": "example.com/store/storetest", "Test": "TestStore", "Key": key,
				"Value": part,
			})
			b.Write(line)
			b.WriteByte('\n')
		}
	}
	b.WriteString(`{"Action":"attr","Package":"example.com/store/storetest","Test":"TestStore",` +
		`"Key":"dokimi.assert.1","Value":"{}"}` + "\n")
	return b.String()
}

func TestReport(t *testing.T) {
	t.Parallel()

	falsified := record(spec.Falsification, check("Store.Read/zero-on-error", spec.Proven))
	judged := record(spec.OnSubject, check("Store.Read/zero-on-error", spec.Vacuous),
		check("Store.Read/keeps-the-item", spec.Pass), check("Store.Put/smoke", spec.Excused))
	stopped := func(application spec.Application) string {
		return strings.Replace(record(application), `"checks":null`, `"ended":"the option Cases asks for 0 cases\n`+
			`the excusal of Store.Get/x has no reason","checks":null`, 1)
	}

	t.Run("Read", func(t *testing.T) {
		t.Parallel()

		t.Run("returns each suite record of the events with its package and its test", func(t *testing.T) {
			t.Parallel()

			got, err := report.Read(strings.NewReader(events(64, falsified, judged)))

			assert.NoError(t, err, "the events read")
			assert.Length(t, got, 2, "the events have two suite records")
			assert.Equal(t, got[0].Package, "example.com/store/storetest", "the record has the package of its test")
			assert.Equal(t, got[0].Test, "TestStore", "the record has its test")
			assert.Equal(t, got[0].Application, spec.Falsification, "the parts of a record join in their order")
			assert.Equal(t, got[1].Checks[0].Result, spec.Vacuous, "the second record follows the first")
		})

		t.Run("reads an event longer than the buffer of a scanner", func(t *testing.T) {
			t.Parallel()
			checks := make([]spec.CheckRecord, 2000)
			for i := range checks {
				checks[i] = check("Store.Read/keeps-the-item", spec.Pass)
			}

			got, err := report.Read(strings.NewReader(events(1<<20, record(spec.OnSubject, checks...))))

			assert.NoError(t, err, "the events read")
			assert.Length(t, got, 1, "the long event has one record")
			assert.Length(t, got[0].Checks, 2000, "the record has each check")
		})

		tests := []struct {
			name       string
			giveEvents string
			want       string
		}{
			{
				name:       "returns an error for a record that does not validate",
				giveEvents: events(4096, `{"spec": "1.0.0"}`),
				want:       "example.com/store/storetest TestStore dokimi.suite.1: spec: the suite record does not validate",
			},
			{
				name:       "returns an error that names the line of an event that does not decode",
				giveEvents: "\n\n" + `{"Key": 7}` + "\n",
				want:       "line 3: json: cannot unmarshal number into Go struct field",
			},
			{
				name:       "returns an error for a line longer than a mebibyte",
				giveEvents: "{" + strings.Repeat(" ", 1<<20) + "}\n",
				want:       "read the events: bufio.Scanner: token too long",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				_, err := report.Read(strings.NewReader(tt.giveEvents))

				assert.HasError(t, err, "the events do not read")
				assert.Contains(t, err.Error(), tt.want, "the error names the event")
			})
		}
	})

	t.Run("Suites", func(t *testing.T) {
		t.Parallel()

		t.Run("groups the records by declaration and lists the checks that are vacuous in each application",
			func(t *testing.T) {
				t.Parallel()
				engaged := record(spec.OnSubject, check("Store.Read/keeps-the-item", spec.Vacuous),
					check("Store.Read/zero-on-error", spec.Vacuous), check("Store.Put/smoke", spec.Excused))
				records, err := report.Read(strings.NewReader(events(4096, falsified, judged,
					record(spec.OnSubject, check("Store.Read/keeps-the-item", spec.Pass)), engaged)))
				assert.NoError(t, err, "the events read")

				got := report.Suites(records)

				assert.Length(t, got, 1, "the records have one suite")
				assert.Equal(t, got[0].Declaration, "Store", "the suite has its declaration")
				assert.Length(t, got[0].Records, 4, "the suite has each record")
				assert.Equal(t, got[0].Vacuous, []string{"Store.Read/zero-on-error"},
					"a check is vacuous when each application that judged it found it vacuous")
			})

		t.Run("leaves out an application that excused, skipped or left a check unchecked", func(t *testing.T) {
			t.Parallel()
			records, err := report.Read(strings.NewReader(events(
				4096,
				record(
					spec.OnSubject,
					check("Store.Read/excused", spec.Vacuous),
					check("Store.Read/skipped", spec.Vacuous),
					check("Store.Read/unchecked", spec.Vacuous),
				),
				record(
					spec.OnSubject,
					check("Store.Read/excused", spec.Excused),
					check("Store.Read/skipped", spec.Skipped),
					check("Store.Read/unchecked", spec.Unchecked),
				),
			)))
			assert.NoError(t, err, "the events read")

			got := report.Suites(records)

			assert.Equal(
				t,
				got[0].Vacuous,
				[]string{"Store.Read/excused", "Store.Read/skipped", "Store.Read/unchecked"},
				"an application without a verdict on a check does not count",
			)
		})

		t.Run("returns a suite for each package and each declaration", func(t *testing.T) {
			t.Parallel()
			logRecord := strings.Replace(falsified, `"declaration":"Store"`, `"declaration":"Log"`, 1)
			otherPackage := strings.Replace(
				falsified,
				`"package":"example.com/store"`,
				`"package":"example.com/other"`,
				1,
			)
			records, err := report.Read(strings.NewReader(events(4096, falsified, logRecord, otherPackage, judged)))
			assert.NoError(t, err, "the events read")

			got := report.Suites(records)

			assert.Length(t, got, 3, "the records have three suites")
			assert.Equal(t, []string{got[1].Declaration, got[2].Package}, []string{"Log", "example.com/other"},
				"the suites are in the order of their first records")
			assert.Length(t, got[0].Records, 2, "the first suite has its two records")
		})
	})

	t.Run("Failures", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name        string
			giveRecords []string
			want        []string
		}{
			{
				name: "returns nothing for proven checks and checks without a verdict",
				giveRecords: []string{record(spec.Falsification, check("Store.Read/zero-on-error", spec.Proven),
					check("Store.Read/allocation-budget", spec.Unchecked), check("Store.Put/smoke", spec.Skipped),
					check("Store.Put/cancel", spec.Excused))},
			},
			{
				name:        "returns a check whose falsification result is not proven",
				giveRecords: []string{record(spec.Falsification, check("Store.Read/zero-on-error", spec.WrongReason))},
				want:        []string{"Store.Read/zero-on-error is wrong-reason in falsification"},
			},
			{
				name:        "returns nothing for a check that fails on the subject",
				giveRecords: []string{record(spec.OnSubject, check("Store.Read/zero-on-error", spec.Fail))},
			},
			{
				name:        "returns each problem that ended an application once",
				giveRecords: []string{stopped(spec.Falsification), stopped(spec.OnSubject)},
				want: []string{
					"the application in TestStore ended before its checks: the option Cases asks for 0 cases",
					"the application in TestStore ended before its checks: the excusal of Store.Get/x has no reason",
				},
			},
			{
				name: "returns the problems of an application after an application of the same test that ran",
				giveRecords: []string{
					record(spec.OnSubject, check("Store.Read/zero-on-error", spec.Pass)),
					stopped(spec.OnSubject),
				},
				want: []string{
					"the application in TestStore ended before its checks: the option Cases asks for 0 cases",
					"the application in TestStore ended before its checks: the excusal of Store.Get/x has no reason",
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				records, err := report.Read(strings.NewReader(events(4096, tt.giveRecords...)))
				assert.NoError(t, err, "the events read")

				assert.Equal(t, report.Suites(records)[0].Failures(), tt.want, "the failures are the failed verdicts")
			})
		}
	})

	t.Run("Print", func(t *testing.T) {
		t.Parallel()

		t.Run("writes the claims, the results, the refusals, the unclaimed and the vacuous checks", func(t *testing.T) {
			t.Parallel()
			records, err := report.Read(strings.NewReader(events(4096, falsified, judged)))
			assert.NoError(t, err, "the events read")
			var b strings.Builder

			report.Print(&b, report.Suites(records))

			assert.Equal(t, b.String(), "example.com/store  Store\n"+
				"  claims       1   proven 1   wrong-reason 0   passed 0\n"+
				"  subject      3   pass 1   vacuous 1   fail 0   excused 1   skipped 0\n"+
				"  refused      1   Store.Scan/zero-on-error: Scan cannot fail\n"+
				"  unclaimed    1   mixin:concurrent on Store.Put: needs concurrent callers\n"+
				"  vacuous          Store.Read/zero-on-error\n"+
				"seed 0x5f3a1c9e27d4b860   spec 1.0.0   go.dokimi.dev/conform v0.9.0\n",
				"the report has one block for the suite")
		})

		t.Run("writes each other result that occurs and each refusal on a line of its own", func(t *testing.T) {
			t.Parallel()
			isolated := record(spec.Falsification, check("Store.Read/zero-on-error", spec.NotIsolated))
			isolated = strings.Replace(isolated, `"refusals":[`,
				`"refusals":[{"claim":"cancel","callable":"Store.Scan","why":"Scan takes no context","remedy":"x"},`, 1)
			isolated = strings.Replace(
				isolated,
				`"unclaimed":[{"classification":"mixin:concurrent","callable":"Store.Put",`+
					`"stance":"needs","note":"needs concurrent callers"}]`,
				`"unclaimed":[]`,
				1,
			)
			records, err := report.Read(strings.NewReader(events(4096, isolated)))
			assert.NoError(t, err, "the events read")
			var b strings.Builder

			report.Print(&b, report.Suites(records))

			assert.Contains(t, b.String(), "  claims       1   proven 0   wrong-reason 0   passed 0   not-isolated 1\n",
				"the count line writes the other result after the results that it always writes")
			assert.Contains(t, b.String(), "  refused      2   Store.Scan/cancel: Scan takes no context\n"+
				"                   Store.Scan/zero-on-error: Scan cannot fail\n",
				"each refusal after the first has a line of its own")
			assert.Contains(t, b.String(), "  unclaimed    0\n", "a suite without unclaimed classifications writes 0")
		})

		t.Run("writes each problem that ended an application", func(t *testing.T) {
			t.Parallel()
			records, err := report.Read(strings.NewReader(events(4096, stopped(spec.Falsification),
				stopped(spec.OnSubject))))
			assert.NoError(t, err, "the events read")
			var b strings.Builder

			report.Print(&b, report.Suites(records))

			assert.Contains(t, b.String(), "  ended            TestStore: the option Cases asks for 0 cases\n"+
				"  ended            TestStore: the excusal of Store.Get/x has no reason\n"+
				"seed ", "each problem has a line of its own before the seed")
		})
	})
}
